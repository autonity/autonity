package autonity

import (
	"errors"
	"github.com/clearmatics/autonity/ethdb"
	"github.com/clearmatics/autonity/ethdb/leveldb"
	"github.com/clearmatics/autonity/log"
	"io"
	"math/big"
	"os"
	"path"
	"sync"
)


const datadirAutonityContractABI = "abistore"       // Path within the datadir to store the autonity abi spec.

var ErrCannotCreateABIStore = errors.New("cannot create abi store")
var ErrNoABIFoundFromDB = errors.New("not abi found from db")
var ABISpecKey = "ABI_SPEC"

type ABIDB interface {
	ethdb.KeyValueReader
	ethdb.KeyValueWriter
	io.Closer
}

type ABIStore struct {
	db      ABIDB
	baseDir string
	subDir  string
	cache   string
	m       sync.RWMutex
	logger  log.Logger
}

func NewABIStore(logger log.Logger, baseDir string, subDir string) *ABIStore {
	abi := &ABIStore{
		baseDir: baseDir,
		subDir:  subDir,
		cache:   "",
		m:       sync.RWMutex{},
		logger:  logger,
		db:      nil,
	}
	return abi
}

func getDir(basedir string, subDir string) string {

	p := path.Join(basedir, subDir)

	var err error
	if _, err = os.Stat(p); os.IsNotExist(err) {
		err = os.MkdirAll(p, os.ModePerm)
	}
	if err != nil {
		log.Error("can't create ABIStore directory", "path", p, "err", err)
	}

	return p
}

func (abi *ABIStore) prepareDB() bool {
	if abi.db == nil {
		dir := getDir(abi.baseDir, abi.subDir)
		newDB, err := leveldb.New(dir, 128, 1024, "")
		if err != nil {
			return false
		}
		abi.db = newDB
	}
	return true
}

func (abi *ABIStore) UpdateABI(height *big.Int, abiSpec string) error {
	if height == nil || abiSpec == "" {
		return ErrWrongParameter
	}

	abi.m.Lock()
	defer abi.m.Unlock()

	if !abi.prepareDB() {
		return ErrCannotCreateABIStore
	}

	//to do db update.
	err := abi.db.Put([]byte(ABISpecKey), []byte(abiSpec))
	if err != nil {
		return err
	}

	abi.cache = abiSpec
	return nil
}

func (abi *ABIStore) GetABI(height *big.Int) (string, error) {
	if height == nil {
		return "", ErrWrongParameter
	}
	abi.m.Lock()
	defer abi.m.Unlock()

	if abi.cache != "" {
		return abi.cache, nil
	}

	if !abi.prepareDB() {
		return "", ErrCannotCreateABIStore
	}

	//to do db reading.
	bytes, err := abi.db.Get([]byte(ABISpecKey))
	if err != nil {
		return "", err
	}

	if len(bytes) <= 0 {
		return "", ErrNoABIFoundFromDB
	}

	if abi.cache == "" {
		abi.cache = string(bytes)
	}

	return abi.cache, nil
}

func (abi *ABIStore) Close() {
	abi.m.Lock()
	defer abi.m.Unlock()
	if abi.db != nil{
		err := abi.db.Close()
		if err != nil {
			log.Error("close ABIStore failed: ", err)
		}
	}
}
