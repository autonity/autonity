package gengen

import (
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"math/big"
	"net"
	"strconv"
	"time"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/core/interfaces"
	"github.com/autonity/autonity/core"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/crypto"
	"github.com/autonity/autonity/crypto/blst"
	"github.com/autonity/autonity/p2p/enode"
	"github.com/autonity/autonity/params"
)

// User holds the parameters that constitute a participant's initial state in
// the genesis file.
type Validator struct {
	// InitialEth defines the starting eth in wei (1 ETH is 10^18 wei).
	InitialEth *big.Int
	// Stake defines the amount of Stake token a user has in the system.
	// Note this is the sum of tokens delegates and SelfBondedStake
	Stake uint64
	// SelfBondedStake defines the amount of penalty absorbing stake
	SelfBondedStake uint64
	// NodeIP is the ip that this user's node can be reached at.
	NodeIP net.IP
	// NodePort is the port that this user's node can be reached at.
	NodePort int
	// NodeKey is a private key for the validator node.
	NodeKey *ecdsa.PrivateKey
	// OracleKey is a private key for the oracle node.
	OracleKey *ecdsa.PrivateKey
	// AcnIP is the ip that this user's consensus channel is running at
	AcnIP net.IP
	// AcnPort is the port that this user's consensus channel is listening at
	AcnPort int
	// Key is either a public or private key for the treasury account.
	TreasuryKey *ecdsa.PrivateKey
	// ConsensusKey is the BLS key for validator who participate in consensus.
	ConsensusKey blst.SecretKey
	// NodeKeyPath is the file path at which the node key is stored.
	NodeKeyPath string
	// OracleKeyPath is the file path at which the oracle key is stored.
	OracleKeyPath string
	// TreasuryKeyPath is the file path at which the treasury account key is stored.
	TreasuryKeyPath string
	// TendermintServices is an optional familly of consensus hooks used for testing purposes.
	TendermintServices *interfaces.Services
}

type GenesisOption func(*core.Genesis)

func copyConfig(original *params.ChainConfig) *params.ChainConfig {
	jsonBytes, err := json.Marshal(original)
	if err != nil {
		panic("cannot marshal genesis config: " + err.Error())
	}
	genesisCopy := &params.ChainConfig{}
	err = json.Unmarshal(jsonBytes, genesisCopy)
	if err != nil {
		panic("cannot unmarshal genesis config: " + err.Error())
	}
	// deep copying through json marshaling kinda messes up the ABIs (we have `{}` instead of `nil` for empty elements)
	// so let's just copy them from the original, we do not modify them in tests anyways for now
	// TODO: find a better solution to deep copy also the ABIs.
	// I suspect it is not straightforward due to how the generated ABIs are computed.
	// See gen-contract target in Makefile and `TestAPIGetContractABI`
	genesisCopy.AutonityContractConfig.ABI = original.AutonityContractConfig.ABI
	genesisCopy.OracleContractConfig.ABI = original.OracleContractConfig.ABI
	return genesisCopy
}

// NewGenesis parses the input from the commandline and uses it to generate a
// genesis struct. It returns the generated genesis and associated user keys,
// if user keys were provided they will be returned, otherwise a set of
// generated keys will be returned. See gengen command help for a description
// of userStrings and userKeys, userKeys must either have a key for each user
// or be nil.
func NewGenesis(validators []*Validator, options ...GenesisOption) (*core.Genesis, error) {
	if len(validators) < 1 {
		return nil, fmt.Errorf("at least one user must be specified")
	}

	operatorAddress, genesisValidators, genesisAlloc, err := generateValidatorState(validators)
	if err != nil {
		return nil, fmt.Errorf("failed to construct initial user state: %v", err)
	}

	config := copyConfig(params.TestChainConfig)
	config.AutonityContractConfig.Operator = *operatorAddress
	config.AutonityContractConfig.Validators = genesisValidators
	config.Ethash = nil
	genesis := &core.Genesis{

		Timestamp: uint64(time.Now().Unix()),
		Mixhash:   types.BFTDigest,

		// We need to set this to an empty slice to maintain consistency
		// between the genesis we return here and the same genesis encoded and
		// decoded from JSON. This is beause go-ethereum marshals []byte using
		// hexutil.Bytes which marshals nil to '0x' and '0x' subsequently
		// unmarshals to an empty slice.
		ExtraData: []byte{},

		GasLimit: 1_000_000_000, // gas limit setting in piccadilly network.
		//GasLimit: 10000000000, // gas limit in e2e test framework.

		BaseFee: new(big.Int).SetUint64(params.TestAutonityContractConfig.MinBaseFee),

		// Autonity relies on the difficulty always being 1 so that we can
		// compare chain length by comparing total difficulty during peer
		// connection handshake.
		Difficulty: big.NewInt(0),

		Alloc: genesisAlloc,

		Config: config,
	}

	for _, genesisOption := range options {
		genesisOption(genesis)
	}

	return genesis, nil
}

// ParseUint provides support for parsing large numbers in base 10 using
// scientific notation with either e or E used as a prefix for the exponent.
// E.g 10.456e34 or 1234E5. If the resultant number cannot be parsed, is not an
// exact integer or does not fit into 256 bits an error is returned.
func ParseUint(str string) (*big.Int, error) {

	// We parse the string as a base 10 number and allow for 64 bit precision.
	// Providing the base means that the function will fail to parse anything
	// with a prefix, e.g 0x... 0b... 0X... 0B... .
	f, _, err := big.ParseFloat(str, 10, 256, big.ToZero)
	if err != nil {
		return nil, fmt.Errorf("could not parse %s: %v", str, err)

	}

	// Check to see the number is positive.
	if f.Cmp(big.NewFloat(0)) == -1 {
		return nil, fmt.Errorf("the nunber defined by %q is less than 0", str)
	}

	// Check to see if the number is an integer.
	result, accuracy := f.Int(nil)
	if accuracy != big.Exact {
		return nil, fmt.Errorf("the nunber defined by %q is not an integer", str)
	}
	return result, nil
}

// Generates a slice of params.User along with a corresponding
// core.GenesisAlloc. Also returns the address of the first user in users as
// the operatorAddress.
func generateValidatorState(validators []*Validator) (
	operatorAddress *common.Address,
	genesisValidators []*params.Validator,
	genesisAlloc core.GenesisAlloc,
	err error,
) {
	genesisValidators = make([]*params.Validator, len(validators))
	genesisAlloc = make(core.GenesisAlloc, len(validators))
	for i, u := range validators {
		if u.SelfBondedStake > u.Stake {
			return nil, nil, nil, fmt.Errorf("selfBondedStake (%d) cannot be higher than total stake (%d)", u.SelfBondedStake, u.Stake)
		}

		if u.NodeKey == nil || u.ConsensusKey == nil || u.OracleKey == nil || u.TreasuryKey == nil {
			return nil, nil, nil, fmt.Errorf("validator had nil key")
		}

		e := enode.NewV4(&u.NodeKey.PublicKey, u.NodeIP, u.NodePort, u.NodePort)
		ens := enode.AppendConsensusEndpoint(u.AcnIP.String(), strconv.Itoa(u.AcnPort), e.String())

		treasuryAddress := crypto.PubkeyToAddress(u.TreasuryKey.PublicKey)
		oracleAddress := crypto.PubkeyToAddress(u.OracleKey.PublicKey)

		gu := params.Validator{
			OracleAddress:   oracleAddress,
			Enode:           ens,
			Treasury:        treasuryAddress, // rewards goes here
			BondedStake:     new(big.Int).SetUint64(u.Stake),
			SelfBondedStake: new(big.Int).SetUint64(u.SelfBondedStake),
			ConsensusKey:    u.ConsensusKey.PublicKey().Marshal(),
		}
		err = gu.Validate()
		if err != nil {
			return nil, nil, nil, fmt.Errorf("invalid user: %v", err)
		}

		genesisValidators[i] = &gu
		userAddress := crypto.PubkeyToAddress(u.NodeKey.PublicKey)
		if i == 0 {
			operatorAddress = &userAddress
		}
		genesisAlloc[treasuryAddress] = core.GenesisAccount{
			Balance: u.InitialEth,
		}
		genesisAlloc[userAddress] = core.GenesisAccount{
			Balance: u.InitialEth,
		}
	}

	return operatorAddress, genesisValidators, genesisAlloc, nil
}
