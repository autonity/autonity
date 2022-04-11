package gengen

import (
	"crypto/ecdsa"
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/autonity/autonity/core"
	"github.com/autonity/autonity/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	validValidators = []string{"1e12,v,1,:6789,key1", "1e12,s,1,:6799,key2", "1e12,p,1,:6780,key3"}
)

// This test runs the command and checks that the output is can be json
// unmarshaled into a core.Genesis instance.
func TestGenesisCommand(t *testing.T) {
	// We make temp files for the paths.
	out, cleanup := tempFile(t)
	defer cleanup()

	keyfile1, cleanup := tempFile(t)
	// We delete this file immediately so that a key is genrated for this validator,
	// but we use the temp path as the destination.
	cleanup()
	defer cleanup()

	// We fill the following 2 key files with a public and private key.
	keyfile2, cleanup := tempFile(t)
	defer cleanup()
	k, err := crypto.GenerateKey()
	require.NoError(t, err)
	err = ioutil.WriteFile(keyfile2, crypto.PrivECDSAToHex(k), os.ModePerm)
	require.NoError(t, err)

	keyfile3, cleanup := tempFile(t)
	defer cleanup()
	err = ioutil.WriteFile(keyfile3, crypto.PubECDSAToHex(&k.PublicKey), os.ModePerm)
	require.NoError(t, err)

	validator1 := "1e12,v,1,:6789," + keyfile1
	validator2 := "1e12,s,1,:6799," + keyfile2
	validator3 := "1e12,p,1,:6780," + keyfile3

	args := []string{
		"",
		"--" + minGasPriceFlag,
		"10",
		"--" + validatorFlag,
		validator1,
		"--" + validatorFlag,
		validator2,
		"--" + validatorFlag,
		validator3,
		"--" + outFileFlag,
		out,
	}

	c := NewCmd()
	c.SetArgs(args)
	err = c.Execute()
	require.NoError(t, err)

	// Now try to load the genesis from disk
	data, err := ioutil.ReadFile(out)
	require.NoError(t, err)
	decoded := &core.Genesis{}
	err = json.Unmarshal(data, decoded)
	assert.NoError(t, err)
}

// This test checks that a generated *core.Genesis instance file is consistent
// with an instance obtained by JSON encoding and decoding it.
func TestEncodeDecodeConsistency(t *testing.T) {
	validators, err := parseValidators(validValidators)
	require.NoError(t, err)
	// Set one of the validators to have a publick key, just to cover more code
	// branches.
	k, ok := validators[0].Key.(*ecdsa.PrivateKey)
	require.True(t, ok, "key should be an *ecdsa.PrivateKey")
	validators[0].Key = &k.PublicKey
	g, err := NewGenesis(10, validators)
	require.NoError(t, err)
	encoded, err := json.Marshal(g)
	require.NoError(t, err)
	decoded := &core.Genesis{}
	err = json.Unmarshal(encoded, decoded)
	require.NoError(t, err)

	assert.Equal(t, g, decoded)
}

func TestGenesisCreationErrors(t *testing.T) {
	// nil validators
	_, err := NewGenesis(10, nil)
	assert.Error(t, err, "no validators provided")

	// Validator with nil key
	validators, err := parseValidators(validValidators)
	require.NoError(t, err)
	validators[0].Key = nil

	_, err = NewGenesis(10, validators)
	assert.Error(t, err, "validator had nil key")

	// Validator with key of invalid type
	validators, err = parseValidators(validValidators)
	require.NoError(t, err)
	validators[0].Key = "I am not a key"

	_, err = NewGenesis(10, validators)
	assert.Error(t, err, "validator had invalid type of key")

	// Invalid validator type
	validators, err = parseValidators(validValidators)
	require.NoError(t, err)
}

// Checks that errors are thrown appropriately in the case of invalid validators.
func TestValidatorParsingErrors(t *testing.T) {

	validator := ""
	_, err := ParseValidator(validator)
	assert.Error(t, err, "empty validator")

	validator = "1e12,v,:6789,key"
	_, err = ParseValidator(validator)
	assert.Error(t, err, "missing field")

	validator = "1e12zz,v,1,:6789,key"
	_, err = ParseValidator(validator)
	assert.Error(t, err, "invalid initial eth")

	validator = "456.789,v,1,:6789,key"
	_, err = ParseValidator(validator)
	assert.Error(t, err, "fractional initial eth")

	validator = "1e12,v,1.8446744e20,:6789,key"
	_, err = ParseValidator(validator)
	assert.Error(t, err, "stake out of uint64 range")

	validator = "1e12,v,stake,:6789,key"
	_, err = ParseValidator(validator)
	assert.Error(t, err, "invalid stake")

	validator = "1e12,v,-1,:6789,key"
	_, err = ParseValidator(validator)
	assert.Error(t, err, "invalid stake")

	validator = "1e12,v,1.2,:6789,key"
	_, err = ParseValidator(validator)
	assert.Error(t, err, "fractional stake")

	validator = "1e12,v,1,:6789999,key"
	_, err = ParseValidator(validator)
	assert.Error(t, err, "invalid port")

	validator = "1e12,v,1,:-1,key"
	_, err = ParseValidator(validator)
	assert.Error(t, err, "invalid port")

	validator = "1e12zz,v,1,:port,key"
	_, err = ParseValidator(validator)
	assert.Error(t, err, "invalid port")

	validator = "1e12,v,1,lll:6789,key"
	_, err = ParseValidator(validator)
	assert.Error(t, err, "invalid ip")

	validator = "1e12,v,1,akakak,key"
	_, err = ParseValidator(validator)
	assert.Error(t, err, "invalid address")

	validator = "1e12,v,1,lll:6789," + string(byte(0))
	_, err = ParseValidator(validator)
	assert.Error(t, err, "invalid key file name")

}

// Checks that when there is no file provided the keys are generated for validators.
func TestKeyRandomGeneration(t *testing.T) {
	validator := "1e12,v,1,:6789,key"

	u, err := ParseValidator(validator)
	require.NoError(t, err)

	// We expect key to have been generated because the file 'key' does not
	// exist.
	key1, ok := u.Key.(*ecdsa.PrivateKey)
	require.True(t, ok, "expecting key of type *ecdsa.PrivateKey")

	u, err = ParseValidator(validator)
	require.NoError(t, err)
	key2, ok := u.Key.(*ecdsa.PrivateKey)
	require.True(t, ok, "expecting key of type *ecdsa.PrivateKey")

	// We expect subsequent runs to generate a different (random) key.
	assert.NotEqual(t, key1, key2)
}

// Checks that if a file with a key is provided, the key is loaded from the file.
func TestKeysLoadedFromFile(t *testing.T) {

	// Make temp files for keys
	keyFile1, cleanup := tempFile(t)
	defer cleanup()

	keyFile2, cleanup := tempFile(t)
	defer cleanup()

	// Store keys to files
	key1, err := crypto.GenerateKey()
	require.NoError(t, err)
	key2, err := crypto.GenerateKey()
	require.NoError(t, err)

	// Store private key in key1File
	err = ioutil.WriteFile(keyFile1, crypto.PrivECDSAToHex(key1), os.ModePerm)
	require.NoError(t, err)

	// Store public key in key2File
	err = ioutil.WriteFile(keyFile2, crypto.PubECDSAToHex(&key2.PublicKey), os.ModePerm)
	require.NoError(t, err)

	// Check private key loaded from file
	validator := "1e12,v,1,:6789," + keyFile1
	u, err := ParseValidator(validator)
	assert.NoError(t, err)
	assert.Equal(t, key1, u.Key)

	// Check public key loaded from file
	validator = "1e12,v,1,:6789," + keyFile2
	u, err = ParseValidator(validator)
	assert.NoError(t, err)
	assert.Equal(t, &key2.PublicKey, u.Key)
}

// Checks that errors are thrown appropriately in the case of invalid validator
// keys.
func TestKeyParsingErrors(t *testing.T) {

	// Make temp file for keys
	keyFile, cleanup := tempFile(t)
	defer cleanup()

	// Write incorrect length garbage to file
	err := ioutil.WriteFile(keyFile, []byte("kjcld"), os.ModePerm)
	require.NoError(t, err)

	_, err = readKey(keyFile)
	assert.Error(t, err, "garbage provided in key file")

	// Write a private key missing the last hex char
	k, err := crypto.GenerateKey()
	require.NoError(t, err)
	data := crypto.PrivECDSAToHex(k)
	err = ioutil.WriteFile(keyFile, data[:len(data)-1], os.ModePerm)
	require.NoError(t, err)

	_, err = readKey(keyFile)
	assert.Error(t, err, "invalid key provided in key file")

	// Write a public key missing the last hex char
	data = crypto.PubECDSAToHex(&k.PublicKey)
	err = ioutil.WriteFile(keyFile, data[:len(data)-1], os.ModePerm)
	require.NoError(t, err)

	_, err = readKey(keyFile)
	assert.Error(t, err, "invalid key provided in key file")
}

func tempFile(t *testing.T) (name string, cleanup func()) {
	f, err := ioutil.TempFile("", "")
	require.NoError(t, err)
	err = f.Close()
	require.NoError(t, err)
	return f.Name(), func() { os.Remove(f.Name()) }
}
