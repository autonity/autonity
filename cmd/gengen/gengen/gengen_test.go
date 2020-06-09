package gengen

import (
	"encoding/json"
	"testing"

	"github.com/clearmatics/autonity/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	validUsers = []string{"1e12,v,1,:6789", "1e12,v,1,:6799"}
)

// This test checks that a generated *core.Genesis instance file is consistent
// with an instance obtained by JSON encoding and decoding it.
func TestEncodeDecodeConsistency(t *testing.T) {
	// We want to use the temp name, we dont actually want the file to exist when calling newGenesis.
	g, _, err := newGenesis(10, validUsers, nil)
	require.NoError(t, err)
	encoded, err := json.Marshal(g)
	require.NoError(t, err)
	decoded := &core.Genesis{}
	err = json.Unmarshal(encoded, decoded)
	require.NoError(t, err)

	assert.Equal(t, g, decoded)
}

// The gengen tool should only generate users (params.User) with enodes
// specified. It does this because it is redundant to specify the address as
// well since that can be derived from the enode, and it is an error to specify
// conflicting enodes and addresses, so by not specifying address we avoid this
// case.
func TestUsersAddressIsNil(t *testing.T) {
	g, _, err := newGenesis(10, validUsers, nil)
	require.NoError(t, err)

	assert.Nil(t, g.Config.AutonityContractConfig.Users[0].Address)
	assert.Nil(t, g.Config.AutonityContractConfig.Users[1].Address)
}

// Checks that errors are thrown appropriately in the case of invalid users.
func TestUserParsingErrors(t *testing.T) {

	_, _, err := newGenesis(10, nil, nil)
	assert.Error(t, err, "no users provided")

	user := ""
	_, _, err = newGenesis(10, []string{user}, nil)
	assert.Error(t, err, "empty user")

	user = "1e12,v,:6789"
	_, _, err = newGenesis(10, []string{user}, nil)
	assert.Error(t, err, "missing field")

	user = "1e12zz,v,1,:6789"
	_, _, err = newGenesis(10, []string{user}, nil)
	assert.Error(t, err, "invalid initial eth")

	user = "1e12,q,1,:6789"
	_, _, err = newGenesis(10, []string{user}, nil)
	assert.Error(t, err, "invalid user type")

	user = "1e12,v,stake,:6789"
	_, _, err = newGenesis(10, []string{user}, nil)
	assert.Error(t, err, "invalid stake")

	user = "1e12,v,-1,:6789"
	_, _, err = newGenesis(10, []string{user}, nil)
	assert.Error(t, err, "invalid stake")

	user = "1e12,v,1,:6789999"
	_, _, err = newGenesis(10, []string{user}, nil)
	assert.Error(t, err, "invalid port")

	user = "1e12,v,1,:-1"
	_, _, err = newGenesis(10, []string{user}, nil)
	assert.Error(t, err, "invalid port")

	user = "1e12zz,v,1,:port"
	_, _, err = newGenesis(10, []string{user}, nil)
	assert.Error(t, err, "invalid port")

	user = "1e12,v,1,lll:6789"
	_, _, err = newGenesis(10, []string{user}, nil)
	assert.Error(t, err, "invalid ip")

	user = "1e12,p,1,:6789"
	_, _, err = newGenesis(10, []string{user}, nil)
	assert.Error(t, err, "invalid user type and stake combination")
}

func TestKeysProcessing(t *testing.T) {
	_, generatedKeys, err := newGenesis(10, validUsers, nil)
	require.NoError(t, err)

	// Check keys were generated for users.
	require.Equal(t, len(validUsers), len(generatedKeys))

	_, keys, err := newGenesis(10, validUsers, generatedKeys)
	require.NoError(t, err)

	// Check that when keys are provided the same keys are returned.
	assert.Equal(t, generatedKeys, keys)
}

// Checks that errors are thrown appropriately in the case of invalid user keys.
func TestKeyParsingErrors(t *testing.T) {

	//  Generate a valid set of keys
	_, keys, err := newGenesis(10, validUsers, nil)
	require.NoError(t, err)

	_, _, err = newGenesis(10, validUsers, []string{})
	assert.Error(t, err, "no keys provided")

	_, _, err = newGenesis(10, validUsers, []string{keys[0]})
	assert.Error(t, err, "insufficient keys provided")

	_, _, err = newGenesis(10, validUsers, []string{keys[0] + "x", keys[1]})
	assert.Error(t, err, "invalid hex encoded key")

	_, _, err = newGenesis(10, validUsers, []string{keys[0] + "ff", keys[1]})
	assert.Error(t, err, "invalid key")
}
