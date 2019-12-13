package autonity

import (
	"errors"
	"github.com/clearmatics/autonity/ethdb"
	"github.com/clearmatics/autonity/ethdb/leveldb"
	"github.com/clearmatics/autonity/log"
	"io"
	"os"
	"path"
	"sync"
)

const DirAutonityContractState = "contract_state_store"
const NamespaceContractState = "contract_state"
var ErrCreateContractStateStore = errors.New("cannot create contract state store")
var ErrNotFoundFromDB = errors.New("not found from state db")

var (
	DefaultStore = NewStore()
)

type ContractStateDB interface {
	ethdb.KeyValueReader
	ethdb.KeyValueWriter
	io.Closer
}

type ContractStateStore struct {
	db      ContractStateDB
	baseDir string
	subDir  string
	cache   map[string]string        // key value cache
	m       sync.RWMutex
}

func NewStore() *ContractStateStore {
	store := &ContractStateStore{
		subDir:  DirAutonityContractState,
		cache:   make(map[string]string),
		m:       sync.RWMutex{},
		db:      nil,
	}
	return store
}

func getDir(basedir string, subDir string) string {

	p := path.Join(basedir, subDir)

	var err error
	if _, err = os.Stat(p); os.IsNotExist(err) {
		err = os.MkdirAll(p, os.ModePerm)
	}
	if err != nil {
		log.Error("can't create ContractStateStore directory", "path", p, "err", err)
	}

	return p
}

func (st *ContractStateStore) InitDB(baseDir string) error {

	if len(baseDir) <= 0 || baseDir == "" {
		return ErrCreateContractStateStore
	}

	st.m.Lock()
	defer st.m.Unlock()

	if st.db == nil {
		dir := getDir(baseDir, st.subDir)
		newDB, err := leveldb.New(dir, 128, 1024, NamespaceContractState)
		if err != nil {
			return ErrCreateContractStateStore
		}
		st.baseDir = baseDir
		st.db = newDB
	}
	return nil
}

func (st *ContractStateStore) Put(key []byte, value []byte) error {
	if len(key) <= 0 || len(value) <= 0 {
		return ErrWrongParameter
	}

	st.m.Lock()
	defer st.m.Unlock()

	err := st.db.Put(key, value)
	if err != nil {
		return err
	}

	st.cache[string(key)] = string(value)
	return nil
}

func (st *ContractStateStore) Get(key []byte) ([]byte, error) {
	st.m.Lock()
	defer st.m.Unlock()

	if value, ok := st.cache[string(key)]; ok {
		return []byte(value), nil
	}

	//to do db reading.
	bytes, err := st.db.Get(key)
	if err != nil {
		return nil, err
	}

	if len(bytes) <= 0 {
		return nil, ErrNotFoundFromDB
	}

	if value, ok := st.cache[string(key)]; !ok {
		st.cache[string(key)] = value
	}

	return bytes, nil
}

func (st *ContractStateStore) Close() {
	st.m.Lock()
	defer st.m.Unlock()
	if st.db != nil{
		err := st.db.Close()
		if err != nil {
			log.Error("close ContractStateStore failed: ", err)
		}
	}
}
