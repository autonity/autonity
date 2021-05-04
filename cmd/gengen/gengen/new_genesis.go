package gengen

import (
	"crypto/ecdsa"
	"crypto/rand"
	"fmt"
	"math/big"
	"net"
	"time"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/common/math"
	"github.com/clearmatics/autonity/consensus/tendermint/config"
	"github.com/clearmatics/autonity/core"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/crypto"
	"github.com/clearmatics/autonity/p2p/enode"
	"github.com/clearmatics/autonity/params"
)

// User holds the parameters that constitute a participant's initial state in
// the genesis file.
type User struct {
	// InitialEth defines the starting eth in wei (1 ETH is 10^18 wei).
	InitialEth *big.Int
	// UserType is one of participant, stakeholder or validator.
	UserType params.UserType
	// Stake defines the amount of Stake token a user has in the system.
	Stake uint64
	// NodeIP is the ip that this user's node can be reached at.
	NodeIP net.IP
	// NodePort is the port that this user's node can be reached at.
	NodePort int
	// Key is either a public or private key for the user.
	Key interface{}
	// KeyPath is the file path at which the key is stored.
	KeyPath string
}

// NewGenesis parses the input from the commandline and uses it to generate a
// genesis struct. It returns the generated genesis and associated user keys,
// if user keys were provided they will be returned, otherwise a set of
// generated keys will be returned. See gengen command help for a description
// of userStrings and userKeys, userKeys must either have a key for each user
// or be nil.
func NewGenesis(minGasPrice uint64, users []*User) (*core.Genesis, error) {
	if len(users) < 1 {
		return nil, fmt.Errorf("at least one user must be specified")
	}

	operatorAddress, genesisUsers, genesisAlloc, err := generateUserState(users)
	if err != nil {
		return nil, fmt.Errorf("failed to construct initial user state: %v", err)
	}

	chainID, err := rand.Int(rand.Reader, math.MaxBig256)
	if err != nil {
		return nil, fmt.Errorf("failed to generate random chainID: %v", err)
	}

	genesis := &core.Genesis{

		Timestamp: uint64(time.Now().Unix()),
		Mixhash:   types.BFTDigest,

		// We need to set this to an empty slice to maintain consistency
		// between the genesis we return here and the same genesis encoded and
		// decoded from JSON. This is beause go-ethereum marshals []byte using
		// hexutil.Bytes which marshals nil to '0x' and '0x' subsequently
		// unmarshals to an empty slice.
		ExtraData: []byte{},

		GasLimit: math.MaxUint64,

		// Autonity relies on the difficulty always being 1 so that we can
		// compare chain length by comparing total difficulty during peer
		// connection handshake.
		Difficulty: big.NewInt(1),

		Alloc: genesisAlloc,

		Config: &params.ChainConfig{
			ChainID:             chainID,
			HomesteadBlock:      big.NewInt(0),
			EIP150Block:         big.NewInt(0),
			EIP155Block:         big.NewInt(0),
			EIP158Block:         big.NewInt(0),
			ByzantiumBlock:      big.NewInt(0),
			ConstantinopleBlock: big.NewInt(0),
			PetersburgBlock:     big.NewInt(0),
			IstanbulBlock:       big.NewInt(0),
			AutonityContractConfig: &params.AutonityContractGenesis{
				MinGasPrice: minGasPrice,
				Operator:    *operatorAddress,
				Validators:  genesisUsers,
			},
		},
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
func generateUserState(users []*User) (
	operatorAddress *common.Address,
	genesisUsers []params.User,
	genesisAlloc core.GenesisAlloc,
	err error,
) {
	genesisUsers = make([]params.User, len(users))
	genesisAlloc = make(core.GenesisAlloc, len(users))
	for i, u := range users {
		var pk *ecdsa.PublicKey
		switch k := u.Key.(type) {
		case *ecdsa.PublicKey:
			pk = k
		case *ecdsa.PrivateKey:
			pk = &k.PublicKey
		default:
			return nil, nil, nil, fmt.Errorf("expecting ecdsa public or private key, instead got %T", u.Key)
		}
		e := enode.NewV4(pk, u.NodeIP, u.NodePort, u.NodePort)
		gu := params.User{
			Enode: e.String(),
			Type:  u.UserType,
			Stake: u.Stake,
		}
		err := gu.Validate()
		if err != nil {
			return nil, nil, nil, fmt.Errorf("invalid user: %v", err)
		}
		genesisUsers[i] = gu

		userAddress := crypto.PubkeyToAddress(*pk)
		if i == 0 {
			operatorAddress = &userAddress
		}
		genesisAlloc[userAddress] = core.GenesisAccount{
			Balance: u.InitialEth,
		}
	}

	return operatorAddress, genesisUsers, genesisAlloc, nil
}
