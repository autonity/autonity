package gengen

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/clearmatics/autonity/core"
	"github.com/spf13/cobra"
)

var (
	minGasPrice  uint64
	users        []string
	userKeysFile string
	outFile      string

	// Note in order to achieve a consistent output formatting for the flag
	// descriptions we need to de-indent lines such that they have no leading
	// whitespace. To ensure that all lines can be easily wrapped to the same
	// length we make the first line empty. Even if a parameter description is
	// not multi-line we still use an initial blank line since this provides a
	// consistent look in the outputted help. These descriptions have line
	// wrapping applied at 80 chars.

	helpDescription = `
Help for gengen`

	minGasPriceDescription = `
A 64 bit non negative integer that sets the minimum gas price in wei (1 ETH is
10^18 wei) for submitting transactions to the network. Transactions submitted
with a gas price lower than this will be ignored. This is what will determine
the minimum cost for users to interact with the network. And consequently how
much gas validators will receive for running the network. Depending on how
GasLimit is set there may or may not be an incentive to submit transactions
with higher gas prices than the minimum.`

	userDescription = `
Specifies the parameters for a user.Can be specified multiple times, once for
each user. The first user to be defined on the command line will become the
operator (a user with special privileges and the ability to administer the
system). A user is defined by 4 parameters separated by a comma with no spaces.

Initial eth - The starting eth in wei (1 ETH is 10^18 wei) for this user, a non
negative 256 bit integer value. This can be specified with scientific notation,
the max value specified in scientific notation would be 1.1579209e77.

User type - One of p (participant), s (stakeholder) or v (validator).

Stake - The stake this user holds in the system, a non negative 64 bit integer.
This can be specified with scientific notation, max value specified using
scientific notation would be 1.8446744e19. It is an error to set a stake value
other than zero for participants. The total stake in the system is sum of the
stake allocated for each user but cannot exceed 2^256-1.

Address - The address is specified as an optional IP followed by a colon and
the port number. e.g. 234.334.455.566:3030 or simply :3030 when no IP is
specified a default local IP of 127.0.0.1 is used.

So an example of a user string could be 10,v,1,:3030 or
8,s,1,133.433.654.777:3030.  An example of an invalid user string would be
10,p,1,:7070 since the stake value is 1 and participants cannot have stake.`

	userKeysDescription = `
Specifies the path at which the user keys are stored. Each line holds one key.
The format of a line can be either:

priv:<hex encoded secp256k1 ecdsa private key>
or
pub:<hex encoded secp256k1 ecdsa public key>

If no file is found at this path, then keys will be generated for each user
specified with the '-user' flag and the private keys stored in a new file at
this path.

If a file containing keys is found at this path then the keys in the file are
used as the keys for the users specified by the '-user' flag, in their
respective order. This approach can be used when setting up a production
network, participants can send their public keys to someone who can generate
the genesis on the behalf of the network. It is an error to have a differing
number of keys and users.`

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
		Example: `./gengen --min-gas-price 10 --user 1e12,v,1,:6789 --user 1e12,v,1,:6799 --user-keys userkeys`,
		RunE:    generateGenesis,
	}

	// Set up persistent flags
	flags := rootCmd.PersistentFlags()

	// Override default help flag to set message formatted with a preceding
	// newline.
	flags.BoolP("help", "h", false, helpDescription)

	// We panic on making these flags required since the error returned
	// indicates a programming error.
	flags.Uint64Var(&minGasPrice, "min-gas-price", 0, minGasPriceDescription)
	err := rootCmd.MarkPersistentFlagRequired("min-gas-price")
	if err != nil {
		panic(err)
	}

	flags.StringArrayVar(&users, "user", nil, userDescription)
	err = rootCmd.MarkPersistentFlagRequired("user")
	if err != nil {
		panic(err)
	}

	flags.StringVar(&userKeysFile, "user-keys", "", userKeysDescription)
	err = rootCmd.MarkPersistentFlagRequired("user-keys")
	if err != nil {
		panic(err)
	}

	flags.StringVar(&outFile, "out-file", "", outFileDescription)
	err = rootCmd.MarkPersistentFlagRequired("out-file")
	if err != nil {
		panic(err)
	}

	return rootCmd

}

func generateGenesis(cmd *cobra.Command, args []string) error {

	var userKeys []string
	// If no file exists then we don't read the keys
	_, err := os.Stat(userKeysFile)
	if err == nil {
		userKeys, err = readKeys(userKeysFile)
		if err != nil {
			return err
		}
	}

	genesis, userKeys, err := newGenesis(minGasPrice, users, userKeys)
	if err != nil {
		return fmt.Errorf("failed to generate genesis: %v", err)
	}
	err = writeKeys(userKeysFile, userKeys)
	if err != nil {
		return fmt.Errorf("failed to write user keys: %v", err)
	}
	err = writeGenesis(outFile, genesis)
	if err != nil {
		return fmt.Errorf("failed to write user keys: %v", err)
	}
	return nil
}

// readKeys reads the file and returns a slice of strings one per line.
func readKeys(keyFile string) ([]string, error) {

	// Read the keys from the file
	content, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return nil, err
	}

	// We trim the whitespace here otherwise we can end up with a 0 length trailing '\r' string.
	keyStrings := strings.Split(strings.TrimSpace(string(content)), "\n")
	return keyStrings, nil
}

// writeKeys writes the keys to file at path, one per line.
func writeKeys(path string, userKeys []string) error {
	writeErr := func(file string, err error) error {
		return fmt.Errorf("failed to write keys to %q: %v", file, err)
	}
	f, err := os.Create(path)
	if err != nil {
		return writeErr(path, err)
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	for _, k := range userKeys {
		_, err := w.WriteString(k + "\n")
		if err != nil {
			return writeErr(path, err)
		}
	}
	err = w.Flush()
	if err != nil {
		return writeErr(path, err)
	}
	return nil
}

// writeGenesis writes a json encoded representation of genesis to path.
func writeGenesis(path string, genesis *core.Genesis) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create genesis file: %v", err)
	}
	defer f.Close()
	e := json.NewEncoder(f)
	e.SetIndent("", "  ")
	err = e.Encode(genesis)
	if err != nil {
		return fmt.Errorf("failed to marshal genesis to json: %v", err)
	}
	return nil
}
