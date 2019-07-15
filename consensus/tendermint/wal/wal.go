package wal

import (
	"bytes"
	"fmt"
	"github.com/clearmatics/autonity/ethdb"
	"github.com/clearmatics/autonity/log"
	"math/big"
	"os"
	"path"
	"sort"
	"sync"
)

type WAL struct {
	db      *ethdb.LDBDatabase
	dbDir   string
	baseDir string
	m       sync.RWMutex
	height  *big.Int
}

type Value interface {
	Key() []byte
	Value() ([]byte, error)
}

var currentHeightKey = []byte("current_height")

const keysPrefix = "height-"

func New(basedir string, height *big.Int) *WAL {
	db, err := getDb(basedir, height)
	if err != nil {
		panic(err)
	}

	log.Warn("WAL inited to", "dir", basedir)
	return &WAL{db: db, dbDir: basedir, height: height}
}

func getDb(basedir string, height fmt.Stringer) (*ethdb.LDBDatabase, error) {
	return ethdb.NewLDBDatabase(getDir(basedir, height), 128, 1024)
}

func getDir(basedir string, height fmt.Stringer) string {
	return path.Join(basedir, height.String())
}

func (db *WAL) UpdateHeight(height *big.Int) error {
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

func (db *WAL) Close() {
	db.m.Lock()
	db.db.Close()
	db.m.Unlock()
}

func (db *WAL) Height() (*big.Int, error) {
	db.m.RLock()
	defer db.m.RUnlock()

	value, err := db.db.Get(currentHeightKey)
	if err != nil {
		return nil, err
	}

	return big.NewInt(0).SetBytes(value), nil
}

// thread unsafe
func (db *WAL) setHeight(height *big.Int) error {
	err := db.db.Put(height.Bytes(), currentHeightKey)
	if err != nil {
		return err
	}

	db.height = height

	return nil
}

func (db *WAL) Store(msg Value) error {
	log.Error("WAL: Get store")
	defer log.Error("WAL: store end")

	db.m.Lock()
	defer db.m.Unlock()

	val, err := msg.Value()
	if err != nil {
		return err
	}

	return db.db.Put(keyPrefix(db.height, msg.Key()), val)
}

func (db *WAL) Get(height *big.Int) ([][]byte, error) {
	log.Error("WAL: Get start", "h", height.String())
	defer log.Error("WAL: Get end", "h", height.String())

	currentHeight, err := db.Height()
	if err != nil {
		return nil, nil
	}

	if currentHeight.Cmp(height) != 0 {
		return nil, fmt.Errorf("WAL: trying to set old height. Current %v. Given %v", currentHeight.String(), height.String())
	}

	db.m.RLock()
	defer db.m.RUnlock()
	iterator := db.db.NewIteratorWithPrefix(keyPrefix(height, nil))

	var values []storedValue
	for iterator.Next() {
		values = append(values, storedValue{iterator.Value(), iterator.Key()})
	}

	iterator.Release()

	sort.Slice(values, func(i, j int) bool { return bytes.Compare(values[i].key, values[j].key) == -1 })

	vs := make([][]byte, 0, len(values))
	for _, v := range values {
		vs = append(vs, v.value)
	}

	return vs, nil
}

func keyPrefix(height *big.Int, key []byte) []byte {
	return append(append([]byte(keysPrefix), height.Bytes()...), key...)
}

type storedValue struct {
	value []byte
	key   []byte
}
