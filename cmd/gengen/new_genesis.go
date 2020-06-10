package main

import (
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
	"net"
	"strconv"
	"strings"
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

type user struct {
	// initialEth defines the starting eth in wei (1 ETH is 10^18 wei).
	initialEth *big.Int
	// userType is one of participant, stakeholder or validator.
	userType params.UserType
	// stake defines the amount of stake token a user has in the system.
	stake uint64
	// nodeIP is the ip that this user's node can be reached at.
	nodeIP net.IP
	// nodePort is the port that this user's node can be reached at.
	nodePort int
}

// newGenesis parses the input from the commandline and uses it to generate a
// genesis struct. It returns the generated genesis and associated user keys,
// if user keys were provided they will be returned, otherwise a set of
// generated keys will be returned. See gengen command help for a description
// of userStrings and userKeys, userKeys must either have a key for each user
// or be nil.
func newGenesis(minGasPrice uint64, userStrings []string, userKeys []string) (*core.Genesis, []string, error) {
	if len(userStrings) < 1 {
		return nil, nil, fmt.Errorf("at least one user must be specified")
	}

	// Holds generated or loaded keys for users
	pubKeys := make([]*ecdsa.PublicKey, len(userStrings))
	privKeys := make([]*ecdsa.PrivateKey, len(userStrings))

	if userKeys == nil {
		// No keys provided so generate keys
		userKeys = make([]string, len(userStrings))
		for i := range privKeys {
			k, err := crypto.GenerateKey()
			if err != nil {
				return nil, nil, fmt.Errorf("failed to generate key for user #%d: %v", i, err)
			}
			privKeys[i] = k
			pubKeys[i] = &k.PublicKey
			userKeys[i] = "priv:" + hex.EncodeToString(crypto.FromECDSA(k))
		}

	} else {
		if len(userKeys) != len(userStrings) {
			return nil, nil, fmt.Errorf(
				"%d users specified on commandline but user keys has %d keys",
				len(userStrings),
				len(userKeys),
			)
		}

		invalidKeyEntryErr := func(entry string, index int) error {
			return fmt.Errorf("invalid key entry %q at index %d", entry, index)
		}
		// Expecting lines of the format
		//
		// priv:<hex encoded secp256k1 ecdsa private key>
		// pub:<hex encoded secp256k1 ecdsa public key>
		for i, ks := range userKeys {
			parts := strings.Split(ks, ":")
			if len(parts) != 2 {
				return nil, nil, invalidKeyEntryErr(ks, i)
			}
			b, err := hex.DecodeString(parts[1])
			if err != nil {
				return nil, nil, fmt.Errorf("%v: %v", invalidKeyEntryErr(ks, i), err)
			}

			switch parts[0] {
			case "priv":
				priv, err := crypto.ToECDSA(b)
				if err != nil {
					return nil, nil, fmt.Errorf("%v: %v", invalidKeyEntryErr(ks, i), err)
				}
				pubKeys[i] = &priv.PublicKey
			case "pub":
				pubKeys[i], err = crypto.UnmarshalPubkey(b)
				if err != nil {
					return nil, nil, fmt.Errorf("%v: %v", invalidKeyEntryErr(ks, i), err)
				}
			default:
				return nil, nil, invalidKeyEntryErr(ks, i)
			}
		}
	}
	users := make([]*user, len(userStrings))

	// Parse users
	for i, userString := range userStrings {
		user, err := parseUser(userString)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to parse user %q: %v", userString, err)
		}
		users[i] = user
	}

	operatorAddress, genesisUsers, genesisAlloc, err := generateUserState(users, pubKeys)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to construct initial user state: %v", err)
	}

	chainID, err := rand.Int(rand.Reader, math.MaxBig256)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate random chainID: %v", err)
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
			Tendermint: &config.Config{
				BlockPeriod: 1,
			},
			AutonityContractConfig: &params.AutonityContractGenesis{
				MinGasPrice: minGasPrice,
				Operator:    *operatorAddress,
				Users:       genesisUsers,
			},
		},
	}

	return genesis, userKeys, nil
}

func parseUser(u string) (*user, error) {
	fields := strings.Split(u, ",")
	if len(fields) != 4 {
		return nil, fmt.Errorf("user strings need 4 fields, invalid user string %q", u)
	}

	initialEth, err := ParseUint(fields[0])
	if err != nil {
		return nil, fmt.Errorf("failed to parse initial eth: %v", err)
	}

	var userType params.UserType
	switch fields[1] {
	case "p":
		userType = params.UserParticipant
	case "s":
		userType = params.UserStakeHolder
	case "v":
		userType = params.UserValidator
	default:
		return nil, fmt.Errorf("failed to parse user type %q, not one of u, s or p", fields[1])
	}

	bigStake, err := ParseUint(fields[2])
	if err != nil {
		return nil, fmt.Errorf("failed to parse stake: %v", err)
	}
	if !bigStake.IsInt64() {
		return nil, fmt.Errorf("stake %q is not an integer in the uint64 domain", fields[2])
	}
	stake := bigStake.Uint64()

	address := fields[3]
	ipString, portString, err := net.SplitHostPort(address)
	if err != nil {
		return nil, fmt.Errorf("failed to parse address: %v", err)
	}

	// Try to parse port number into a 16 bit unsigned int
	port, err := strconv.ParseUint(portString, 10, 16)
	if err != nil {
		return nil, fmt.Errorf("failed to parse port %q in address: %v", portString, err)
	}
	if ipString == "" {
		ipString = "127.0.0.1"
	}

	ip := net.ParseIP(ipString)
	if ip == nil {
		return nil, fmt.Errorf("failed to parse ip %q in address: %v", ipString, err)
	}

	user := &user{
		initialEth: initialEth,
		userType:   userType,
		stake:      stake,
		nodeIP:     ip,
		nodePort:   int(port),
	}

	return user, nil
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
func generateUserState(users []*user, keys []*ecdsa.PublicKey) (
	operatorAddress *common.Address,
	genesisUsers []params.User,
	genesisAlloc core.GenesisAlloc,
	err error,
) {
	genesisUsers = make([]params.User, len(users))
	genesisAlloc = make(core.GenesisAlloc, len(users))
	for i, u := range users {
		pk := keys[i]
		e := enode.NewV4(pk, u.nodeIP, u.nodePort, u.nodePort)
		gu := params.User{
			Enode: e.String(),
			Type:  u.userType,
			Stake: u.stake,
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
			Balance: u.initialEth,
		}
	}

	return operatorAddress, genesisUsers, genesisAlloc, nil
}
