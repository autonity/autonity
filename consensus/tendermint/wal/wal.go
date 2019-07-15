package wal

import (
	"bytes"
	"fmt"
	"github.com/clearmatics/autonity/ethdb"
	"github.com/clearmatics/autonity/log"
	"math/big"
	"os"
	"path"
	"sync"
)

type wal struct {
	db *ethdb.LDBDatabase
	dbDir string
	baseDir string
	m sync.RWMutex
	height *big.Int
}

type Value interface {
	Key() []byte
	Value() []byte
}

var currentHeightKey = []byte("current_height")
const keysPrefix = "height-"

func New(basedir string, height *big.Int) *wal {
	db, err := getDb(basedir, height)
	if err != nil {
		panic(err)
	}

	return &wal{db: db, dbDir: basedir, height: height}
}

func getDb(basedir string, height *big.Int) (*ethdb.LDBDatabase, error) {
	return ethdb.NewLDBDatabase(getDir(basedir, height), 128, 1024)
}

func getDir(basedir string, height *big.Int) string {
	return path.Join(basedir, height.String())
}

func (db *wal) UpdateHeight(height *big.Int) error {
	db.m.Lock()
	defer db.m.Unlock()

	oldHeight, err := db.Height()
	if err != nil {
		return err
	}

	if oldHeight.Cmp(height) == 0 {
		return nil
	}

	if oldHeight.Cmp(height) > 0 {
		log.Error("WAL: trying to set old height", "current", oldHeight.String(), "new", height.String())
		return nil
	}

	newDB, err := getDb(db.baseDir, height)
	if err != nil {
		return err
	}
	err = db.setHeight(height)
	if err != nil {
		return err
	}

	// close old db
	db.db.Close()

	db.db = newDB

	err = os.RemoveAll(getDir(db.baseDir, oldHeight))
	if err != nil {
		log.Error("WAL: cant remove old db", "path", getDir(db.baseDir, oldHeight))
	}

	return nil
}

func (db *wal) Close() {
	db.m.Lock()
	db.db.Close()
	db.m.Unlock()
}

func (db *wal) Height() (*big.Int, error) {
	db.m.RLock()
	defer db.m.RUnlock()

	value, err := db.db.Get(currentHeightKey)
	if err != nil {
		return nil, err
	}

	return big.NewInt(0).SetBytes(value), nil
}

// thread unsafe
func (db *wal) setHeight(height *big.Int) error {
	err := db.db.Put(height.Bytes(), currentHeightKey)
	if err != nil {
		return err
	}

	db.height = height

	return nil
}

func (db *wal) Store(msg Value) error {
	db.m.Lock()
	defer db.m.Unlock()

	return db.db.Put(keyPrefix(db.height, msg.Key()), msg.Value())
}

func (db *wal) Get(height *big.Int) error {
	currentHeight, err := db.Height()
	if err != nil {
		return nil
	}

	if currentHeight.Cmp(height) != 0 {
		return fmt.Errorf("WAL: trying to set old height. Current %v. Given %v", currentHeight.String(), height.String())
	}


	db.m.Lock()
	defer db.m.Unlock()

	iterator := db.db.NewIteratorWithPrefix(keyPrefix(height, nil))

	var values []storedValue
	for iterator.Next() {
		values = append(values, storedValue{iterator.Value(), iterator.Key()})
	}

	iterator.Release()

	return db.db.Put(msg.Key(), msg.Value())
}

func keyPrefix(height *big.Int, key []byte) []byte {
	return append(append([]byte(keysPrefix), height.Bytes()...), key...)
}

type storedValue struct {
	value []byte
	key []byte
}

func sort() {}