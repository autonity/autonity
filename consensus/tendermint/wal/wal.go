package wal

import (
	"bytes"
	"fmt"
	"github.com/clearmatics/autonity/ethdb"
	"github.com/clearmatics/autonity/log"
	"github.com/syndtr/goleveldb/leveldb"
	"math/big"
	"os"
	"path"
	"sort"
	"sync"
)

type WAL struct {
	db      *ethdb.LDBDatabase
	baseDir string
	height  *big.Int
	cache   map[string]struct{}

	m sync.RWMutex
}

type DataStore interface {
	Key() []byte
	Value() []byte
}

var currentHeightKey = []byte("current_height")

const keysPrefix = "height-"

func New(basedir string, height *big.Int) *WAL {
	db, err := getDb(basedir, height)
	if err != nil {
		panic(err)
	}

	log.Warn("WAL initialised to", "dir", basedir)
	wal := &WAL{db: db, baseDir: basedir, cache: make(map[string]struct{})}

	err = wal.setHeight(height)
	if err != nil {
		panic(err)
	}

	return wal
}

func getDb(basedir string, height fmt.Stringer) (*ethdb.LDBDatabase, error) {
	return ethdb.NewLDBDatabase(getDir(basedir, height), 128, 1024)
}

func getDir(basedir string, height fmt.Stringer) string {

	p := path.Join(basedir, height.String())

	var err error
	if _, err = os.Stat(p); os.IsNotExist(err) {
		err = os.MkdirAll(p, os.ModePerm)
	}
	if err != nil {
		log.Error("can't create WAL directory", "path", p, "err", err)
	}

	return p
}

func (db *WAL) UpdateHeight(height *big.Int) error {

	oldHeight, err := db.Height()
	if err != nil && err != leveldb.ErrNotFound {
		return err
	}

	if oldHeight.Cmp(height) == 0 {
		return nil
	}

	if oldHeight.Cmp(height) > 0 {
		log.Error("WAL: trying to set old height", "current", oldHeight.String(), "new", height.String())
		return nil
	}

	log.Info("WAL: get height", "base", db.baseDir, "asked", height.String(), "current", db.height.String())

	newDB, err := getDb(db.baseDir, height)
	if err != nil {
		return err
	}

	// close old db
	db.m.Lock()
	db.db.Close()

	db.db = newDB
	db.cache = make(map[string]struct{})
	db.m.Unlock()

	err = db.setHeight(height)
	if err != nil {
		return err
	}

	go func() {
		log.Warn("WAL: removing old db", "path", getDir(db.baseDir, oldHeight))
		err = os.RemoveAll(getDir(db.baseDir, oldHeight))
		if err != nil {
			log.Error("WAL: cant remove old db", "path", getDir(db.baseDir, oldHeight))
		}
	}()

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

func (db *WAL) setHeight(height *big.Int) error {

	db.m.RLock()
	err := db.db.Put(currentHeightKey, height.Bytes())
	if err != nil {
		db.m.RUnlock()
		return err
	}
	db.m.RUnlock()

	db.m.Lock()
	db.height = height
	db.m.Unlock()

	return nil
}

func (db *WAL) Store(msg DataStore) error {
	db.m.RLock()
	defer db.m.RUnlock()

	msgKey := keyPrefix(db.height, msg.Key())

	if _, ok := db.cache[string(msgKey)]; ok {
		// already stored
		return nil
	}

	val := msg.Value()
	err := db.db.Put(msgKey, val)
	if err != nil {
		return err
	}

	db.cache[string(msgKey)] = struct{}{}

	return nil
}

func (db *WAL) Get(height *big.Int) ([][]byte, error) {
	currentHeight, err := db.Height()
	if err != nil {
		return nil, err
	}

	if currentHeight.Cmp(height) != 0 {
		return nil, fmt.Errorf("WAL: trying to get different height. Current %v. Given %v", currentHeight.String(), height.String())
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
