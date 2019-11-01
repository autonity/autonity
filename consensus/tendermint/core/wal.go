package core

import (
	"fmt"

	"errors"

	"github.com/clearmatics/autonity/ethdb"
	"github.com/clearmatics/autonity/ethdb/leveldb"
	"github.com/clearmatics/autonity/log"
	ldberr "github.com/syndtr/goleveldb/leveldb/errors"

	"io"
	"math/big"
	"os"
	"path"
	"sync"

	"time"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus/tendermint/events"
	"github.com/clearmatics/autonity/event"
	"github.com/clearmatics/autonity/rlp"
)

type WALService interface {
	UpdateHeight(height *big.Int) error
	Start()
	Close()
}

type WALDB interface {
	ethdb.KeyValueReader
	ethdb.KeyValueWriter
	ethdb.Iteratee
	io.Closer
}

type WAL struct {
	db      WALDB
	baseDir string
	cache   map[string]struct{}
	height  *big.Int

	sub *event.TypeMuxSubscription

	m      sync.RWMutex
	stop   chan struct{}
	logger log.Logger
}

//type DataStore interface {
//	Key() []byte
//	Value() []byte
//}

var currentHeightKey = []byte("current_height")

const keysPrefix = "height-"

func NewWal(logger log.Logger, basedir string, sub *event.TypeMuxSubscription) *WAL {

	logger.Error("WAL initialised to", "dir", basedir)
	wal := &WAL{
		logger:  logger,
		baseDir: basedir,
		cache:   make(map[string]struct{}),
		sub:     sub,
		stop:    make(chan struct{}),
	}

	return wal
}

func getDb(basedir string, height fmt.Stringer) (*leveldb.Database, error) {
	return leveldb.New(getDir(basedir, height), 128, 1024, "")
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

func (wal *WAL) UpdateHeight(height *big.Int) error {

	oldHeight, err := wal.Height()
	if err != nil && err != ldberr.ErrNotFound {
		return err
	}

	if oldHeight.Cmp(height) == 0 {
		return nil
	}

	if oldHeight.Cmp(height) > 0 {
		wal.logger.Error("WAL: trying to set old height", "current", oldHeight.String(), "new", height.String())
		return nil
	}

	wal.logger.Info("WAL: get height", "base", wal.baseDir, "asked", height.String(), "current", wal.height.String())

	newDB, err := getDb(wal.baseDir, height)
	if err != nil {
		return err
	}

	// close old wal
	wal.m.Lock()
	if wal.db != nil {
		wal.db.Close()
	}

	wal.db = newDB
	wal.cache = make(map[string]struct{})
	wal.m.Unlock()

	err = wal.setHeight(height)
	if err != nil {
		return err
	}

	go func() {
		wal.logger.Warn("WAL: removing old wal", "path", getDir(wal.baseDir, oldHeight))
		err = os.RemoveAll(getDir(wal.baseDir, oldHeight))
		if err != nil {
			wal.logger.Error("WAL: cant remove old wal", "path", getDir(wal.baseDir, oldHeight))
		}
	}()

	return nil
}

func (wal *WAL) Start() {
	go wal.EventLoop()
}
func (wal *WAL) EventLoop() error {
	for {
		wal.m.RLock()
		if wal.db == nil {
			wal.m.RUnlock()
			time.Sleep(time.Second)
			continue
		} else {
			wal.m.RUnlock()
		}
		select {
		case ev, ok := <-wal.sub.Chan():
			if !ok {
				continue
			}
			// A real ev arrived, process interesting content
			switch e := ev.Data.(type) {
			case events.MessageEvent:
				height, round, m, err := parseMessageEvent(e)
				if err != nil {
					wal.logger.Error("parseMessageEvent", "err", err)
					continue
				}
				wal.m.RLock()

				switch height.Cmp(wal.height) {
				case 0:
					//wal.logger.Error("Store event", "height", height.Uint64(), "round", round.Uint64())
					wal.Store(height, round, m, e.Payload)
				case 1:
					wal.logger.Error("Wal got future event", "height", height, "wal", wal.height)
				case -1:
					wal.logger.Error("Wal got old event", "height", height, "wal", wal.height)

				}
				wal.m.RUnlock()

				//case backlogEvent:
				//	// No need to check signature for internal messages
				//	log.Debug("Started handling backlogEvent")
				//	p, err := e.msg.Payload()
				//	if err != nil {
				//		log.Debug("core.handleConsensusEvents Get message payload failed", "err", err)
				//		continue
				//	}
				//
				//	//wal.Store()
			}
		}
	}
}

func (wal *WAL) Close() {
	wal.m.Lock()
	wal.db.Close()
	wal.sub.Unsubscribe()
	wal.m.Unlock()
}

func (wal *WAL) Height() (*big.Int, error) {
	wal.m.RLock()
	defer wal.m.RUnlock()
	if wal.db == nil {
		return new(big.Int), nil
	}

	value, err := wal.db.Get(currentHeightKey)
	if err != nil {
		return nil, err
	}

	return big.NewInt(0).SetBytes(value), nil
}

func (wal *WAL) setHeight(height *big.Int) error {

	wal.m.Lock()
	wal.height = height
	wal.m.Unlock()

	wal.m.RLock()
	err := wal.db.Put(currentHeightKey, height.Bytes())
	if err != nil {
		wal.m.RUnlock()
		return err
	}
	wal.m.RUnlock()

	return nil
}

func (wal *WAL) Store(height, round *big.Int, msg Message, payload []byte) error {
	wal.m.RLock()
	defer wal.m.RUnlock()

	msgKey := keyPrefix(height, generateKey(round.Uint64(), msg.Code, msg.Address))

	if _, ok := wal.cache[string(msgKey)]; ok {
		// already stored
		return nil
	}

	val := payload
	err := wal.db.Put(msgKey, val)
	if err != nil {
		return err
	}

	wal.cache[string(msgKey)] = struct{}{}

	return nil
}

func (wal *WAL) Get(height *big.Int) ([]Message, error) {
	currentHeight, err := wal.Height()
	if err != nil {
		return nil, err
	}

	if currentHeight.Cmp(height) != 0 {
		return nil, fmt.Errorf("WAL: trying to get different height. Current %v. Given %v", currentHeight.String(), height.String())
	}

	wal.m.RLock()
	defer wal.m.RUnlock()
	iterator := wal.db.NewIteratorWithPrefix(keyPrefix(height, nil))

	var values []Message
	for iterator.Next() {
		m := Message{}
		err := rlp.DecodeBytes(iterator.Value(), &m)
		if err != nil {
			return nil, errors.New("corrupted message")
		}
		values = append(values, m)
	}

	iterator.Release()

	return values, nil
}

func keyPrefix(height *big.Int, key []byte) []byte {
	return append(append([]byte(keysPrefix), height.Bytes()...), key...)
}
func generateKey(round, code uint64, address common.Address) []byte {
	return []byte(fmt.Sprintf("%v-%v-%s", round, code, address))

}

type storedValue struct {
	value []byte
	key   []byte
}

func parseMessageEvent(e events.MessageEvent) (*big.Int, *big.Int, Message, error) {
	if len(e.Payload) == 0 {
		return nil, nil, Message{}, errors.New("empty event")
	}
	m := Message{}
	err := rlp.DecodeBytes(e.Payload, &m)
	if err != nil {
		return nil, nil, Message{}, errors.New("corrupted message: " + err.Error())
	}

	switch m.Code {
	case msgPrevote:
		v := Vote{}
		err := rlp.DecodeBytes(m.Msg, &v)
		if err != nil {
			return nil, nil, Message{}, errors.New("corrupted payload: " + err.Error())
		}
		return v.Height, v.Round, m, nil
	case msgPrecommit:
		v := Vote{}
		err := rlp.DecodeBytes(m.Msg, &v)
		if err != nil {
			return nil, nil, Message{}, errors.New("corrupted payload: " + err.Error())
		}
		return v.Height, v.Round, m, nil

	case msgProposal:
		v := Proposal{}
		err := rlp.DecodeBytes(m.Msg, &v)
		if err != nil {
			return nil, nil, Message{}, errors.New("corrupted proposal: " + err.Error())
		}
		return v.Height, v.Round, m, nil
	default:
		return nil, nil, Message{}, errors.New("unexpected message code")
	}
}
