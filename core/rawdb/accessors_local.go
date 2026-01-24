package rawdb

import (
	"encoding/binary"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/ethdb"
	"github.com/autonity/autonity/log"
)

// ReadJailedCount retrieves the count of jailed validators for epoch = `epochID`.
func ReadJailedCount(db ethdb.Reader, epochID uint64) uint64 {
	data, _ := db.Get(jailedCountKey(epochID))
	if len(data) == 0 {
		return 0
	}
	return binary.BigEndian.Uint64(data)
}

// WriteJailedCount stores the count of jailed validators for epoch = `epochID`.
func WriteJailedCount(db ethdb.KeyValueWriter, epochID, count uint64) {
	if err := db.Put(jailedCountKey(epochID), encodeNumber(count)); err != nil {
		log.Crit("Failed to store jailed validators count", "epochID", epochID, "count", count, "err", err)
	}
}

func DeleteJailedCount(db ethdb.KeyValueWriter, epochID uint64) {
	if err := db.Delete(jailedCountKey(epochID)); err != nil {
		log.Crit("Failed to delete jailed validators count", "epochID", epochID, "err", err)
	}
}

// ReadJailedAddress retrieves the address of jailed validator `index`.
func ReadJailedAddress(db ethdb.Reader, epochID, index uint64) common.Address {
	data, _ := db.Get(jailedAddressKey(epochID, index))
	if len(data) == 0 {
		return common.Address{}
	}
	return common.BytesToAddress(data)
}

// WriteJailedAddress stores the address of jailed validator `index`.
func WriteJailedAddress(db ethdb.KeyValueWriter, epochID, index uint64, address common.Address) {
	if err := db.Put(jailedAddressKey(epochID, index), address[:]); err != nil {
		log.Crit("Failed to store jailed validator address", "err", err)
	}
}

func DeleteJailedAddress(db ethdb.KeyValueWriter, epochID, index uint64) {
	if err := db.Delete(jailedAddressKey(epochID, index)); err != nil {
		log.Crit("Failed to delete jailed validator address", "err", err)
	}
}
