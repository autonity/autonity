package gengen

import (
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/clearmatics/autonity/core"
	"github.com/clearmatics/autonity/crypto"
	"github.com/spf13/cobra"
)

const (
	validatorFlag = "validator"
	outFileFlag   = "out-file"
)

var (
	validators []string
	outFile    string

	// Note in order to achieve a consistent output formatting for the flag
	// descriptions we need to de-indent lines such that they have no leading
	// whitespace. To ensure that all lines can be easily wrapped to the same
	// length we make the first line empty. Even if a parameter description is
	// not multi-line we still use an initial blank line since this provides a
	// consistent look in the outputted help. These descriptions have line
	// wrapping applied at 80 chars.

	helpDescription = `
Help for gengen`

	validatorDescription = `
Specifies the parameters for a user. Can be specified multiple times, once for
each user. The first user to be defined on the command line will become the
operator (a user with special privileges and the ability to administer the
system). A user is defined by 5 parameters each separated by a comma with no
spaces.

Initial eth - The starting eth in wei (1 ETH is 10^18 wei) for this user, a non
negative 256 bit integer value. This can be specified with decimal, scientific
or hex notation, the max value specified in scientific notation would be
1.1579209e77.

User type - One of p (participant) or v (validator).

Stake - The stake this user holds in the system, a non negative 64 bit integer.
This can be specified with decimal, scientific or hex notation, the max value
specified using scientific notation would be 1.8446744e19. It is an error to
set a stake value other than zero for participants. The total stake in the
system is sum of the stake allocated for each user but cannot exceed 2^256-1.

Address - The address is specified as an optional IP followed by a colon and
the port number. e.g. 234.334.455.566:3030 or simply :3030 when no IP is
specified a default local IP of 127.0.0.1 is used.

Keyfile - The path to the file holding this users key. The content can be
either a private or public hex encoded secp256k1 ecdsa key, if the file does
not exist then a private secp256k1 ecdsa key will be generated and stored with
hex encoding at the given path.

So an example of a user string could be "10,v,1,:3030,key1" or
"8,v,1,133.433.654.777:3030,key2". An example of an invalid user string would
be "10,p,1,:7070,key3" since the stake value is 1 and participants cannot have
stake.`

	outFileDescription = `
Specifies the path at which the generated genesis file will be stored. If a
file exists at this path it will be overwritten`

	gengenDescription = `
gengen is a commandline tool to generate genesis files. It allows you to set
the configurable parameters of an Autonity network and takes care to correctly
set all other genesis parameters that have only one valid value or those that
should be randomised. It has a number of flags which are all required except
the help flag. The user flag can be specified multiple times to define multiple
users. If provided with a set of keys gengen will use those keys for the users,
otherwise it will generate and store a key for each user.`
)

// NewCmd returns a new gengen command.
func NewCmd() *cobra.Command {

	rootCmd := &cobra.Command{
		Use:     "gengen",
		Short:   gengenDescription,
		Example: `./gengen --validator 1e12,v,1,:6789,key1 --validator 1e12,v,1,:6799,key2 --out-file genesis.json`,
		RunE:    generateGenesis,
	}

	// Set up persistent flags
	flags := rootCmd.PersistentFlags()

	// Override default help flag to set message formatted with a preceding
	// newline.
	flags.BoolP("help", "h", false, helpDescription)

	// We panic on making these flags required since the error returned
	// indicates a programming error.
	flags.StringArrayVar(&validators, validatorFlag, nil, validatorDescription)
	err := rootCmd.MarkPersistentFlagRequired("validator")
	if err != nil {
		panic(err)
	}

	flags.StringVar(&outFile, outFileFlag, "", outFileDescription)
	err = rootCmd.MarkPersistentFlagRequired("out-file")
	if err != nil {
		panic(err)
	}

	return rootCmd

}

func generateGenesis(cmd *cobra.Command, args []string) error {

	parsed, err := parseValidators(validators)
	if err != nil {
		return err
	}
	genesis, err := NewGenesis(parsed)
	if err != nil {
		return fmt.Errorf("failed to generate genesis: %v", err)
	}
	err = writeKeys(parsed)
	if err != nil {
		return fmt.Errorf("failed to write validator keys: %v", err)
	}
	err = writeGenesis(outFile, genesis)
	if err != nil {
		return fmt.Errorf("failed to write validator keys: %v", err)
	}
	return nil
}

// readKey reads the file and returns a slice of strings one per line.
func readKey(keyFile string) (interface{}, error) {

	_, err := os.Stat(keyFile)
	if os.IsNotExist(err) {
		// Generate non existent key
		k, err := crypto.GenerateKey()
		if err != nil {
			return nil, fmt.Errorf("failed to generate key for file %q : %v", keyFile, err)
		}
		return k, nil
	}

	// Read the keys from the file
	content, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return nil, err
	}

	var k interface{}
	switch len(content) {
	case 64: // Private key
		k, err = crypto.PrivECDSAFromHex(content)
	case 130: // Public key
		k, err = crypto.PubECDSAFromHex(content)
	default:
		return nil, fmt.Errorf("unexpected key length, expecting 64 or 130 hex chars, insted got %d", len(content))
	}
	return k, err
}

// writeKeys writes the keys to file at path, one per line.
func writeKeys(validators []*Validator) error {
	writeErr := func(file string, err error) error {
		return fmt.Errorf("failed to write key to %q: %v", file, err)
	}
	for _, u := range validators {
		switch k := u.Key.(type) {
		case *ecdsa.PrivateKey:
			err := ioutil.WriteFile(u.KeyPath, crypto.PrivECDSAToHex(k), os.ModePerm)
			if err != nil {
				return writeErr(u.KeyPath, err)
			}

		case *ecdsa.PublicKey:
			err := ioutil.WriteFile(u.KeyPath, crypto.PubECDSAToHex(k), os.ModePerm)
			if err != nil {
				return writeErr(u.KeyPath, err)
			}
		}
	}

	return nil
}

// writeGenesis writes a json encoded representation of genesis to path.
func writeGenesis(path string, genesis *core.Genesis) (err error) {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create genesis file: %v", err)
	}

	defer func() {
		cerr := f.Close()
		if err == nil {
			// if there was no error in execution update err
			// by return value from close
			err = cerr
		}
	}()

	e := json.NewEncoder(f)
	e.SetIndent("", "  ")
	err = e.Encode(genesis)
	if err != nil {
		return fmt.Errorf("failed to marshal genesis to json: %v", err)
	}
	return nil
}

func parseValidators(validators []string) ([]*Validator, error) {
	parsed := make([]*Validator, len(validators))
	for i, u := range validators {
		p, err := ParseValidator(u)
		if err != nil {
			return nil, err
		}
		parsed[i] = p
	}
	return parsed, nil
}

func ParseValidator(u string) (*Validator, error) {
	fields := strings.Split(u, ",")
	if len(fields) != 5 {
		return nil, fmt.Errorf("validator strings need 5 fields, invalid validator string %q", u)
	}

	initialEth, err := ParseUint(fields[0])
	if err != nil {
		return nil, fmt.Errorf("failed to parse initial eth: %v", err)
	}

	bigStake, err := ParseUint(fields[2])
	if err != nil {
		return nil, fmt.Errorf("failed to parse stake: %v", err)
	}
	if !bigStake.IsUint64() {
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

	keyPath := fields[4]
	key, err := readKey(keyPath)
	if err != nil {
		return nil, err
	}

	validator := &Validator{
		InitialEth: initialEth,
		Stake:      stake,
		NodeIP:     ip,
		NodePort:   int(port),
		KeyPath:    keyPath,
		Key:        key,
	}

	return validator, nil
}
