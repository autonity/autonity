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
	m       sync.RWMutex
	height  *big.Int

	cache map[string]struct{}
}

type Value interface {
	Key() []byte
	Value() ([]byte, error)
}

var currentHeightKey = []byte("current_height")

const keysPrefix = "height-"

func New(basedir string, height *big.Int) *WAL {
	fmt.Println("=== New 1")
	defer fmt.Println("=== New 2")

	basedir = path.Join("autonity-data", "wal")

	db, err := getDb(basedir, height)
	if err != nil {
		panic(err)
	}

	log.Warn("WAL initialised to", "dir", basedir)
	wal := &WAL{db: db, baseDir: basedir, cache: make(map[string]struct{})}

	err = wal.setHeight(height)
	if err != nil {
		fmt.Println("=== UpdateHeight 6")
		panic(err)
	}

	return wal
}

func getDb(basedir string, height fmt.Stringer) (*ethdb.LDBDatabase, error) {
	fmt.Println("=== getDb 1")
	defer fmt.Println("=== getDb 2")
	return ethdb.NewLDBDatabase(getDir(basedir, height), 128, 1024)
}

func getDir(basedir string, height fmt.Stringer) string {
	fmt.Println("=== getDir 1")
	defer fmt.Println("=== getDir 2")

	path := path.Join(basedir, height.String())

	var err error
	if _, err = os.Stat(path); os.IsNotExist(err) {
		err = os.MkdirAll(path, os.ModePerm)
	}
	if err != nil {
		log.Error("can't create WAL directory", "path", path, "err", err)
	}

	return path
}

func (db *WAL) UpdateHeight(height *big.Int) error {
	fmt.Println("=== UpdateHeight 1")
	defer fmt.Println("=== UpdateHeight 2")

	oldHeight, err := db.Height()
	if err != nil && err != leveldb.ErrNotFound {
		fmt.Println("=== UpdateHeight X.2")
		fmt.Println("=== UpdateHeight 3", "err", err)
		return err
	}

	fmt.Println("=== UpdateHeight X.3")
	if oldHeight.Cmp(height) == 0 {
		return nil
	}

	if oldHeight.Cmp(height) > 0 {
		fmt.Println("=== UpdateHeight 4")
		log.Error("WAL: trying to set old height", "current", oldHeight.String(), "new", height.String())
		return nil
	}

	log.Error("WAL: get height", "base", db.baseDir, "asked", height.String(), "current", db.height.String())

	newDB, err := getDb(db.baseDir, height)
	if err != nil {
		fmt.Println("=== UpdateHeight 5", err.Error(), height.String(), oldHeight.String())
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
		fmt.Println("=== UpdateHeight 6")
		return err
	}

	go func() {
		log.Warn("WAL: removing old db", "path", getDir(db.baseDir, oldHeight))
		err = os.RemoveAll(getDir(db.baseDir, oldHeight))
		if err != nil {
			log.Error("WAL: cant remove old db", "path", getDir(db.baseDir, oldHeight))
		}
	}()

	fmt.Println("=== UpdateHeight 7")
	return nil
}

func (db *WAL) Close() {
	fmt.Println("=== Close 1")
	defer fmt.Println("=== Close 2")
	db.m.Lock()
	db.db.Close()
	db.m.Unlock()
}

func (db *WAL) Height() (*big.Int, error) {
	fmt.Println("=== Height 1")
	defer fmt.Println("=== Height 2")
	db.m.RLock()
	defer db.m.RUnlock()

	value, err := db.db.Get(currentHeightKey)
	fmt.Println("=== 1", len(value))
	if err != nil {
		fmt.Println("=== 1.1", err)
		return nil, err
	}

	fmt.Println("=== 2", len(value))
	return big.NewInt(0).SetBytes(value), nil
}

func (db *WAL) setHeight(height *big.Int) error {
	fmt.Println("=== setHeight 1")
	defer fmt.Println("=== setHeight 2")

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

func (db *WAL) Store(msg Value) error {
	fmt.Println("=== Store 1")
	defer fmt.Println("=== Store 2")

	log.Error("WAL: Get store")
	defer log.Error("WAL: store end")

	db.m.RLock()
	defer db.m.RUnlock()

	msgKey := keyPrefix(db.height, msg.Key())

	fmt.Println("**** storing", string(msgKey))

	if _, ok := db.cache[string(msgKey)]; ok {
		// already stored
		return nil
	}

	val, err := msg.Value()
	if err != nil {
		return err
	}

	err = db.db.Put(msgKey, val)
	if err != nil {
		return err
	}

	db.cache[string(msgKey)] = struct{}{}
	fmt.Println("**** stored", string(msgKey))

	return nil
}

func (db *WAL) Get(height *big.Int) ([][]byte, error) {
	fmt.Println("msg.m.String()", db.height.String())
	defer fmt.Println("=== Get 2")

	log.Error("WAL: Get start", "h", height.String())
	defer log.Error("WAL: Get end", "h", height.String())

	currentHeight, err := db.Height()
	if err != nil {
		fmt.Println("=== ERROR 1", err.Error())
		return nil, nil
	}
	fmt.Println("=== 3", currentHeight.String())

	if currentHeight.Cmp(height) != 0 {
		fmt.Println("=== ERROR 2")
		return nil, fmt.Errorf("WAL: trying to set old height. Current %v. Given %v", currentHeight.String(), height.String())
	}

	db.m.RLock()
	defer db.m.RUnlock()
	iterator := db.db.NewIteratorWithPrefix(keyPrefix(height, nil))

	fmt.Println("=== 4", currentHeight.String(), height.String())

	var values []storedValue
	for iterator.Next() {
		fmt.Println("=== 5 in next")
		fmt.Println("**** got", string(iterator.Key()))
		values = append(values, storedValue{iterator.Value(), iterator.Key()})
	}
	fmt.Println("=== 6 after next")

	iterator.Release()

	fmt.Println("=== 7 Release")

	sort.Slice(values, func(i, j int) bool { return bytes.Compare(values[i].key, values[j].key) == -1 })

	fmt.Println("=== 8 Sort")

	vs := make([][]byte, 0, len(values))
	for _, v := range values {
		vs = append(vs, v.value)
	}

	fmt.Println("=== 9 return")

	return vs, nil
}

func keyPrefix(height *big.Int, key []byte) []byte {
	return append(append([]byte(keysPrefix), height.Bytes()...), key...)
}

type storedValue struct {
	value []byte
	key   []byte
}
