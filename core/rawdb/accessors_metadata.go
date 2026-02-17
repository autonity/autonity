// Copyright 2018 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package rawdb

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/autonity/autonity/core/types"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/ethdb"
	"github.com/autonity/autonity/log"
	"github.com/autonity/autonity/params"
	"github.com/autonity/autonity/rlp"
)

var (
	// NOTE: this prefix are not used for prefixing keys, but
	// rather for prefixing data. See WriteContractsConfig
	// the code assumes these prefixes will remain of 1 byte of length
	contractsConfigDataPrefix = []byte("c")
	blockNumberDataPrefix     = []byte("b")
)

// ReadDatabaseVersion retrieves the version number of the database.
func ReadDatabaseVersion(db ethdb.KeyValueReader) *uint64 {
	var version uint64

	enc, _ := db.Get(databaseVersionKey)
	if len(enc) == 0 {
		return nil
	}
	if err := rlp.DecodeBytes(enc, &version); err != nil {
		return nil
	}

	return &version
}

// WriteDatabaseVersion stores the version number of the database
func WriteDatabaseVersion(db ethdb.KeyValueWriter, version uint64) {
	enc, err := rlp.EncodeToBytes(version)
	if err != nil {
		log.Crit("Failed to encode database version", "err", err)
	}
	if err = db.Put(databaseVersionKey, enc); err != nil {
		log.Crit("Failed to store the database version", "err", err)
	}
}

// ReadChainConfig retrieves the consensus settings based on the given genesis hash.
func ReadChainConfig(db ethdb.KeyValueReader, hash common.Hash) *params.ChainConfig {
	data, _ := db.Get(configKey(hash))
	if len(data) == 0 {
		return nil
	}
	var config params.ChainConfig
	if err := json.Unmarshal(data, &config); err != nil {
		log.Error("Invalid chain config JSON", "hash", hash, "err", err)
		return nil
	}
	return &config
}

// WriteChainConfig writes the chain config settings to the database.
func WriteChainConfig(db ethdb.KeyValueWriter, hash common.Hash, cfg *params.ChainConfig) {
	if cfg == nil {
		return
	}
	data, err := json.Marshal(cfg)
	if err != nil {
		log.Crit("Failed to JSON encode chain config", "err", err)
	}
	log.Warn("Storing chain config", "hash", hash.String(), "genesis", cfg)
	if err := db.Put(configKey(hash), data); err != nil {
		log.Crit("Failed to store chain config", "err", err)
	}
}

// ReadGenesisStateSpec retrieves the genesis state specification based on the
// given genesis (block-)hash.
func ReadGenesisStateSpec(db ethdb.KeyValueReader, blockhash common.Hash) []byte {
	data, _ := db.Get(genesisStateSpecKey(blockhash))
	return data
}

// WriteGenesisStateSpec writes the genesis state specification into the disk.
func WriteGenesisStateSpec(db ethdb.KeyValueWriter, blockhash common.Hash, data []byte) {
	if err := db.Put(genesisStateSpecKey(blockhash), data); err != nil {
		log.Crit("Failed to store genesis state", "err", err)
	}
}

func isRlpEncodedUint64(b byte) bool {
	return b == blockNumberDataPrefix[0]
}

func isRlpEncodedConfig(b byte) bool {
	return b == contractsConfigDataPrefix[0]
}

func rlpEncodeUint64WithPrefix(number uint64) []byte {
	encoded, err := rlp.EncodeToBytes(number)
	if err != nil {
		panic("failed to encode RLP encoded number, err: " + err.Error())
	}
	return append(blockNumberDataPrefix, encoded...)
}

func rlpDecodeUint64WithPrefix(encoded []byte) (uint64, error) {
	if !isRlpEncodedUint64(encoded[0]) {
		panic("unexpected prefix")
	}
	var number uint64
	// discarding the first byte of encoded since it contains the prefix
	if err := rlp.DecodeBytes(encoded[1:], &number); err != nil {
		return 0, fmt.Errorf("failed to decode RLP encoded number: %w", err)
	}
	return number, nil
}

// returns the config at block `number` and the number at which it was stored in statedb
func ReadContractsConfig(db ethdb.KeyValueReader, number uint64) (*types.ContractsConfig, uint64) {
	requestedNumber := number
	data, _ := db.Get(contractsConfigKey(number))
	// NOTE: this case can happen for the blocks before the pivot block
	// when snap syncing. Those will have no configuration attached.
	if len(data) == 0 {
		return nil, 0
	}

	// if result is a block number, find the config in the respective block
	if isRlpEncodedUint64(data[0]) {
		var err error
		number, err = rlpDecodeUint64WithPrefix(data)
		if err != nil {
			log.Error("Detected corrupted contracts config db", "requestedNumber", requestedNumber, "data", data)
			return nil, 0
		}
		data, _ = db.Get(contractsConfigKey(number))
	}

	if len(data) == 0 || !isRlpEncodedConfig(data[0]) {
		log.Error("Detected corrupted contracts config db", "requestedNumber", requestedNumber, "number", number, "data", data)
		return nil, 0
	}

	config := &types.ContractsConfig{}
	// discarding the first byte of data since it contains the prefix
	err := rlp.DecodeBytes(data[1:], config)
	if err != nil {
		log.Error("Detected corrupted contracts config db", "requestedNumber", requestedNumber, "number", number, "data", data, "err", err)
		return nil, 0
	}
	return config, number
}

func WriteContractsConfig(db ethdb.KeyValueReaderWriter, targetNumber uint64, cfg *types.ContractsConfig) {
	// if writing genesis contracts config
	// no need to check previous ones.
	if targetNumber == 0 {
		writeContractsConfig(db, targetNumber, cfg)
		return
	}

	// check if something changed wrt to previous config.
	// previous config can be nil if:
	// - we are committing the config of the snap synced head
	// - db corruption of some sort
	previousConfig, number := ReadContractsConfig(db, targetNumber-1)
	if previousConfig != nil && previousConfig.Equal(cfg) {
		// nothing changed, just point to the previous config
		if err := db.Put(contractsConfigKey(targetNumber), rlpEncodeUint64WithPrefix(number)); err != nil {
			panic("Failed to store contracts config: " + err.Error()) //nolint:goconst
		}
	} else {
		// if we end up here because we failed to fetch the new config, notify operator
		if previousConfig == nil {
			log.Warn("failed to fetch contract config for target number, de-duplication not possible", "number", targetNumber-1)
		}
		// in any case, write the new config at the right spot
		writeContractsConfig(db, targetNumber, cfg)
	}
}

func writeContractsConfig(db ethdb.KeyValueWriter, number uint64, cfg *types.ContractsConfig) {
	data, err := rlp.EncodeToBytes(cfg)
	if err != nil {
		panic("Failed to RLP encode contracts config: " + err.Error())
	}
	if err := db.Put(contractsConfigKey(number), append(contractsConfigDataPrefix, data...)); err != nil {
		panic("Failed to store contracts config: " + err.Error()) //nolint:goconst
	}
}

// crashList is a list of unclean-shutdown-markers, for rlp-encoding to the
// database
type crashList struct {
	Discarded uint64   // how many ucs have we deleted
	Recent    []uint64 // unix timestamps of 10 latest unclean shutdowns
}

const crashesToKeep = 10

// PushUncleanShutdownMarker appends a new unclean shutdown marker and returns
// the previous data
// - a list of timestamps
// - a count of how many old unclean-shutdowns have been discarded
func PushUncleanShutdownMarker(db ethdb.KeyValueStore) ([]uint64, uint64, error) {
	var uncleanShutdowns crashList
	// Read old data
	if data, err := db.Get(uncleanShutdownKey); err != nil {
		log.Warn("Error reading unclean shutdown markers", "error", err)
	} else if err := rlp.DecodeBytes(data, &uncleanShutdowns); err != nil {
		return nil, 0, err
	}
	var discarded = uncleanShutdowns.Discarded
	var previous = make([]uint64, len(uncleanShutdowns.Recent))
	copy(previous, uncleanShutdowns.Recent)
	// Add a new (but cap it)
	uncleanShutdowns.Recent = append(uncleanShutdowns.Recent, uint64(time.Now().Unix()))
	if count := len(uncleanShutdowns.Recent); count > crashesToKeep+1 {
		numDel := count - (crashesToKeep + 1)
		uncleanShutdowns.Recent = uncleanShutdowns.Recent[numDel:]
		uncleanShutdowns.Discarded += uint64(numDel)
	}
	// And save it again
	data, _ := rlp.EncodeToBytes(uncleanShutdowns)
	if err := db.Put(uncleanShutdownKey, data); err != nil {
		log.Warn("Failed to write unclean-shutdown marker", "err", err)
		return nil, 0, err
	}
	return previous, discarded, nil
}

// PopUncleanShutdownMarker removes the last unclean shutdown marker
func PopUncleanShutdownMarker(db ethdb.KeyValueStore) {
	var uncleanShutdowns crashList
	// Read old data
	if data, err := db.Get(uncleanShutdownKey); err != nil {
		log.Warn("Error reading unclean shutdown markers", "error", err)
	} else if err := rlp.DecodeBytes(data, &uncleanShutdowns); err != nil {
		log.Error("Error decoding unclean shutdown markers", "error", err) // Should mos def _not_ happen
	}
	if l := len(uncleanShutdowns.Recent); l > 0 {
		uncleanShutdowns.Recent = uncleanShutdowns.Recent[:l-1]
	}
	data, _ := rlp.EncodeToBytes(uncleanShutdowns)
	if err := db.Put(uncleanShutdownKey, data); err != nil {
		log.Warn("Failed to clear unclean-shutdown marker", "err", err)
	}
}

// UpdateUncleanShutdownMarker updates the last marker's timestamp to now.
func UpdateUncleanShutdownMarker(db ethdb.KeyValueStore) {
	var uncleanShutdowns crashList
	// Read old data
	if data, err := db.Get(uncleanShutdownKey); err != nil {
		log.Warn("Error reading unclean shutdown markers", "error", err)
	} else if err := rlp.DecodeBytes(data, &uncleanShutdowns); err != nil {
		log.Warn("Error decoding unclean shutdown markers", "error", err)
	}
	// This shouldn't happen because we push a marker on Backend instantiation
	count := len(uncleanShutdowns.Recent)
	if count == 0 {
		log.Warn("No unclean shutdown marker to update")
		return
	}
	uncleanShutdowns.Recent[count-1] = uint64(time.Now().Unix())
	data, _ := rlp.EncodeToBytes(uncleanShutdowns)
	if err := db.Put(uncleanShutdownKey, data); err != nil {
		log.Warn("Failed to write unclean-shutdown marker", "err", err)
	}
}

// ReadTransitionStatus retrieves the eth2 transition status from the database
func ReadTransitionStatus(db ethdb.KeyValueReader) []byte {
	data, _ := db.Get(transitionStatusKey)
	return data
}

// WriteTransitionStatus stores the eth2 transition status to the database
func WriteTransitionStatus(db ethdb.KeyValueWriter, data []byte) {
	if err := db.Put(transitionStatusKey, data); err != nil {
		log.Crit("Failed to store the eth2 transition status", "err", err)
	}
}
