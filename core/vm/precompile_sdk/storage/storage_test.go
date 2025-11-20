package storage

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/autonity/autonity/autonity/tests"
	"github.com/autonity/autonity/common"
)

type Config struct {
	TimeoutSeconds uint32
	IsEnabled      bool
	Padding        [3]byte // to fill the complete 8 bytes
}

type Profile struct {
	NameHash common.Hash
	Level    uint64
	Config   Config
}

type DeepCommission struct {
	MaxRate      uint64
	LastModified int64
	Info         Config
}

// test struct with various field types
type TestState struct {
	// primitive types
	Version uint8
	Address Address
	Balance Uint256

	// nested struct
	AdminProfile Profile

	// dynamic array/slice
	History []common.Hash

	// map with struct
	UserStats map[Address]DeepCommission

	// map with primitive arrays
	Scores map[Address][]uint32
}

var (
	testContractAddr = common.HexToAddress("0xCAFECAFECAFECAFECAFECAFECAFECAFECAFECAFE")
	testUserAddr     = NewAddressFromBytes(common.HexToAddress("0xUSERUSERUSERUSERUSERUSERUSERUSERUSERUSER").Bytes())
	testValidator    = NewAddressFromBytes(common.HexToAddress("0xVALIDATORVALIDATORVALIDATORVALIDATORVALIDATOR").Bytes())
)

func newStorage(t *testing.T) *Storage {
	r := tests.Setup(t, nil)
	slots := AssignSlots(reflect.TypeOf(TestState{}))
	return NewStorage(testContractAddr, r.Evm.StateDB, slots)
}

func TestStoragePrimitives(t *testing.T) {
	st := newStorage(t)
	err := st.SetUint8("Version", 42)
	require.NoError(t, err)
	v, err := st.GetUint8("Version")
	require.NoError(t, err)
	require.Equal(t, uint8(42), v)

	err = st.SetAddress("Address", testUserAddr)
	require.NoError(t, err)
}
