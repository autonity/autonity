package events

import (
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/autonity/autonity/accounts/abi"
	"github.com/autonity/autonity/autonity/tests"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core/vm/pdk/storage"
	"github.com/autonity/autonity/crypto"
)

type SimpleEvent struct {
	ID uint64
}

func TestEmit_SimpleEvent(t *testing.T) {
	r := tests.Setup(t, nil)
	id := uint64(10)
	event := SimpleEvent{ID: id}
	address := common.HexToAddress("0xaabbccddeeff00112233445566778899aabbccdd")
	st := storage.NewStorage(address, r.Evm.StateDB)
	err := Emit(st, event)
	require.NoError(t, err)

	// verify logs
	logs := r.Evm.StateDB.GetLogs(common.Hash{}, common.Hash{})
	require.Len(t, logs, 9)
	log := logs[8]
	require.Equal(t, address, log.Address)

	expectedSig := crypto.Keccak256Hash([]byte("SimpleEvent(uint64)"))
	require.Len(t, log.Topics, 1)
	require.Equal(t, expectedSig, log.Topics[0])

	require.Equal(t, uint64(32), uint64(len(log.Data)))
	fieldData := log.Data[24:32]
	retrieveID := binary.BigEndian.Uint64(fieldData)
	require.Equal(t, id, retrieveID)
}

type IndexedEvent struct {
	ID    uint64
	Title string `indexed:"true"`
	Size  uint32 `indexed:"true"`
}

func TestEmit_IndexedEvent(t *testing.T) {
	r := tests.Setup(t, nil)
	id := uint64(20)
	title := "HelloEvent"
	size := uint32(11)
	event := IndexedEvent{ID: id, Title: title, Size: size}
	address := common.HexToAddress("0x11223344556677889900aabbccddeeff11223344")
	st := storage.NewStorage(address, r.Evm.StateDB)
	err := Emit(st, event)
	require.NoError(t, err)

	// verify logs
	logs := r.Evm.StateDB.GetLogs(common.Hash{}, common.Hash{})
	require.Len(t, logs, 9)
	log := logs[8]
	require.Equal(t, address, log.Address)

	expectedSig := crypto.Keccak256Hash([]byte("IndexedEvent(uint64,string,uint32)"))
	require.Len(t, log.Topics, 3)
	require.Equal(t, expectedSig, log.Topics[0])

	// verify indexed field
	titleHash := crypto.Keccak256Hash([]byte(title))
	require.Equal(t, titleHash, log.Topics[1])

	// verify indexed field 2 (size), which is uint32, we can retrieve the value from topic[2]
	retrieveSize := log.Topics[2].Big().Uint64()
	require.Equal(t, size, uint32(retrieveSize))
	// verify data field
	require.Equal(t, uint64(32), uint64(len(log.Data)))
	fieldData := log.Data[24:32]
	retrieveID := binary.BigEndian.Uint64(fieldData)
	require.Equal(t, id, retrieveID)
}

type DynamicIndexedEvent struct {
	ID      uint64 `indexed:"true"`
	Payload []byte
}

func TestEmit_DynamicIndexedEvent(t *testing.T) {
	r := tests.Setup(t, nil)
	id := uint64(30)
	payload := []byte("DynamicPayloadData")
	event := DynamicIndexedEvent{ID: id, Payload: payload}
	address := common.HexToAddress("0xdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef")
	st := storage.NewStorage(address, r.Evm.StateDB)
	err := Emit(st, event)
	require.NoError(t, err)
	// verify logs
	logs := r.Evm.StateDB.GetLogs(common.Hash{}, common.Hash{})
	require.Len(t, logs, 9)
	log := logs[8]
	require.Equal(t, address, log.Address)
	expectedSig := crypto.Keccak256Hash([]byte("DynamicIndexedEvent(uint64,bytes)"))
	require.Len(t, log.Topics, 2)
	require.Equal(t, expectedSig, log.Topics[0])
	// verify indexed field (ID)
	retrieveID := log.Topics[1].Big().Uint64()
	require.Equal(t, id, retrieveID)
	// verify data field (Payload)
	bytesType, _ := abi.NewType("bytes", "", nil)
	args := abi.Arguments{{Name: "Data", Type: bytesType}}
	unpacked, err := args.Unpack(log.Data)
	require.NoError(t, err)

	require.Equal(t, payload, unpacked[0].([]byte))
}
