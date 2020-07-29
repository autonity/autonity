package rawdb

import (
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"testing"

	"github.com/clearmatics/autonity/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRawdbAccessReturnsNilWhenDatabaseClosed(t *testing.T) {
	temp, err := ioutil.TempDir("", "")
	require.NoError(t, err)
	defer os.RemoveAll(temp)

	root := filepath.Join(temp, "chaindata")
	freezer := filepath.Join(root, "ancient")

	db, err := NewLevelDBDatabaseWithFreezer(
		root,
		256,
		256,
		freezer,
		"eth/db/chaindata/",
	)

	require.NoError(t, err)

	hash := common.Hash{}
	var num uint64 = 1
	td := big.NewInt(10)
	WriteTd(db, hash, num, td)
	retrieved := ReadTd(db, hash, num)
	assert.Equal(t, retrieved, td)
	err = db.Close()
	require.NoError(t, err)
	retrieved = ReadTd(db, hash, num)
	var nilBigInt *big.Int
	assert.Equal(t, retrieved, nilBigInt)
}
