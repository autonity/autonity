package core

import (
	"encoding/binary"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/ethdb"
	"github.com/autonity/autonity/log"
	"github.com/autonity/autonity/metrics"
	"github.com/autonity/autonity/rlp"
	"time"
)

// WAL db schemas
var (
	//todo: Jason, do we need to keep states of last height?

	lastTendermintStateKey = []byte("LastTendermintState")

	// Record the max msg ID as a helper garbage collection. It is OK for overwriting a key-value pair in level DB since
	// the storage engine handles the garbage collection for the overwritten value, thus we save the max msg ID in the
	// store to do a periodical garbage collection: the key-values of range (lastMsgID, maxMsgID] are the only required
	// items to be deleted.
	maxMessageIDKey = []byte("MaxMessageID")

	// Msg ID starts from 0 for each consensus instance, it increases by one for per message
	lastTBFTInstanceMsgIDKey = []byte("TBFTInstanceMsgID")

	messagePrefix = []byte("Message") // messagePrefix + MsgID -> consensus message
)

type TendermintStateDB struct {
	db ethdb.Database // WAL db which share the same key-value DB of blockchain DB.

	maxMsgID           uint64 // cache for max MSG ID
	lastConsensusMsgID uint64 // cache for last consensus msg ID of a height.

	rsRLPMeter    metrics.Meter // Meter for measuring the size of rs RLP-encoded data
	rsRLPEncTimer metrics.Timer // Timer measuring time required for rs RLP encoding
	rsDbSaveTimer metrics.Timer // Timer measuring rs DB write latency

	msgRLPMeter    metrics.Meter // Meter for measuring the size of received consensus messages RLP-encoded data
	msgRLPEncTimer metrics.Timer // Timer measuring time required for received consensus messages to be RLP encoded
	msgDbSaveTimer metrics.Timer // Timer measuring DB write latency for received consensus messages

	logger log.Logger
}

// newTendermintStateDB create the context of WAL database, it shares the same key-value store of blockchain DB.
func newTendermintStateDB(db ethdb.Database) *TendermintStateDB {
	logger := log.New("newTendermintStateDB", "type", "TendermintStateDB")
	rsdb := &TendermintStateDB{
		db:     db,
		logger: logger,
	}

	rsdb.rsRLPMeter = metrics.NewRegisteredMeter("wal/rs/rlp/encoding/size", nil)
	rsdb.rsRLPEncTimer = metrics.NewRegisteredTimer("wal/rs/rlp/encoding/duration", nil)
	rsdb.rsDbSaveTimer = metrics.NewRegisteredTimer("wal/rs/db/save/time", nil)

	rsdb.msgRLPMeter = metrics.NewRegisteredMeter("wal/message/rlp/encoding/size", nil)
	rsdb.msgRLPEncTimer = metrics.NewRegisteredTimer("wal/message/rlp/encoding/duration", nil)
	rsdb.msgDbSaveTimer = metrics.NewRegisteredTimer("wal/message/db/save/time", nil)

	lastMsgID, err := rsdb.GetMsgID(lastTBFTInstanceMsgIDKey)
	if err != nil {
		panic(err)
	}
	maxMsgID, err := rsdb.GetMsgID(maxMessageIDKey)
	if err != nil {
		panic(err)
	}

	rsdb.lastConsensusMsgID = lastMsgID
	rsdb.maxMsgID = maxMsgID

	return rsdb
}

// ExtTendermintState is used for RLP encoding and decoding, it replaces the locked value and valid value with hash to
// improve the performance, on the state recovery end, the values are retrieved from the rounds proposal of the specific
// round.
type ExtTendermintState struct {
	Height   uint64
	Round    uint64
	Step     Step
	Decision common.Hash

	LockedRound      uint64
	IsLockedRoundNil bool
	ValidRound       uint64
	IsValidRoundNil  bool
	LockedValue      common.Hash
	ValidValue       common.Hash

	// extra helper states base on our implementation.
	SentProposal          bool
	SentPrevote           bool
	SentPrecommit         bool
	SetValidRoundAndValue bool
}

// UpdateLastRoundState stores the latest tendermint state in DB, in case of a start of a new height, it also does
// garbage collection and reset MSG ID for a new height.
func (rsdb *TendermintStateDB) UpdateLastRoundState(rs TendermintState, startNewHeight bool) error {
	logger := rsdb.logger.New("func", "UpdateLastRoundState")
	viewKey := lastTendermintStateKey
	// todo: jason, check if we have other db engine options.
	// todo: jason, check if we have other options for encoding.
	extRoundState := ExtTendermintState{
		Height:                rs.height.Uint64(),
		Round:                 uint64(rs.round),
		Step:                  rs.step,
		LockedRound:           uint64(rs.lockedRound),
		ValidRound:            uint64(rs.validRound),
		SentProposal:          rs.sentProposal,
		SentPrevote:           rs.sentPrevote,
		SentPrecommit:         rs.sentPrecommit,
		SetValidRoundAndValue: rs.setValidRoundAndValue,
	}

	if rs.lockedRound == -1 {
		extRoundState.IsLockedRoundNil = true
	}

	if rs.validRound == -1 {
		extRoundState.IsValidRoundNil = true
	}

	if rs.decision != nil {
		extRoundState.Decision = rs.decision.Hash()
	}

	if rs.lockedValue != nil {
		extRoundState.LockedValue = rs.lockedValue.Hash()
	}
	if rs.validValue != nil {
		extRoundState.ValidValue = rs.validValue.Hash()
	}

	before := time.Now()
	entryBytes, err := rlp.EncodeToBytes(&extRoundState)
	rsdb.rsRLPEncTimer.UpdateSince(before)
	if err != nil {
		logger.Error("Failed to save roundState", "reason", "rlp encoding", "err", err)
		return err
	}

	rsdb.rsRLPMeter.Mark(int64(len(entryBytes)))

	before = time.Now()
	batch := rsdb.db.NewBatch()

	// in case of height rotation, check if it is time for garbage collection and reset msg ID for new height.
	if startNewHeight {
		if rs.height.Uint64()%garbageCollectionInterval == 0 {
			rsdb.GarbageCollection()
		}

		rsdb.lastConsensusMsgID = 0
		msgIDBytes, err := rlp.EncodeToBytes(rsdb.lastConsensusMsgID)
		if err != nil {
			rsdb.logger.Error("Failed to reset msg id", "reason", "rlp encoding", "err", err)
			return err
		}
		if err = batch.Put(lastTBFTInstanceMsgIDKey, msgIDBytes); err != nil {
			rsdb.logger.Error("Failed to reset msg id", "reason", "level db put", "err", err)
			return err
		}
	}

	batch.Put(viewKey, entryBytes)
	err = batch.Write()
	rsdb.rsDbSaveTimer.UpdateSince(before)
	if err != nil {
		logger.Error("Failed to save roundState", "reason", "levelDB write", "err", err, "func")
	}
	return err
}

// GetLastTendermintState will return tendermint state from DB, it will return an initial state if there was no state flushed.
// This function is called once on node start up.
func (rsdb *TendermintStateDB) GetLastTendermintState() ExtTendermintState {
	// set default states.
	var entry = ExtTendermintState{}
	viewKey := lastTendermintStateKey
	rawEntry, err := rsdb.db.Get(viewKey)
	if err != nil {
		rsdb.logger.Warn("failed to read tendermint state from WAL", "error", err)
		return entry
	}

	if err = rlp.DecodeBytes(rawEntry, &entry); err != nil {
		rsdb.logger.Warn("failed to read tendermint state from WAL", "error", err)
		return entry
	}
	return entry
}

type ExtMsg struct {
	Msg      message.Msg
	Verified bool
}

// AddMsg inserts a successfully applied consensus message of tendermint state engine into WAL. The inserting messages
// are ordered by a message ID managed in WAL for retrieval and garbage collection.
func (rsdb *TendermintStateDB) AddMsg(msg message.Msg, verified bool) error {
	nextMsgID := rsdb.lastConsensusMsgID + 1
	msgKey := messageKey(nextMsgID)
	before := time.Now()

	msgIDBytes, err := rlp.EncodeToBytes(nextMsgID)
	if err != nil {
		rsdb.logger.Error("Failed to save msg id", "reason", "rlp encoding", "err", err)
		return err
	}

	m := ExtMsg{Msg: msg, Verified: verified}
	msgBytes, err := rlp.EncodeToBytes(&m)
	rsdb.msgRLPEncTimer.UpdateSince(before)
	if err != nil {
		rsdb.logger.Error("Failed to save msg", "reason", "rlp encoding", "err", err)
		return err
	}

	rsdb.msgRLPMeter.Mark(int64(len(msgBytes)))

	before = time.Now()
	batch := rsdb.db.NewBatch()

	// increase max msg ID if current msg ID is greater than it.
	if nextMsgID > rsdb.maxMsgID {
		rsdb.maxMsgID = nextMsgID
		batch.Put(maxMessageIDKey, msgIDBytes)
	}

	batch.Put(lastTBFTInstanceMsgIDKey, msgIDBytes)
	batch.Put(msgKey, msgBytes)
	err = batch.Write()
	rsdb.rsDbSaveTimer.UpdateSince(before)
	if err != nil {
		rsdb.logger.Error("Failed to save roundState", "reason", "levelDB write", "err", err, "func")
		return err
	}

	rsdb.lastConsensusMsgID = nextMsgID
	return nil
}

// GetMsg retrieves the msg with its corresponding ID.
func (rsdb *TendermintStateDB) GetMsg(msgID uint64) (message.Msg, bool, error) {
	msgKey := messageKey(msgID)
	rawEntry, err := rsdb.db.Get(msgKey)
	if err != nil {
		return nil, false, err
	}
	var entry ExtMsg
	if err = rlp.DecodeBytes(rawEntry, &entry); err != nil {
		return nil, false, err
	}
	return entry.Msg, entry.Verified, nil
}

// GetMsgID retrieves the managed MSG ID from DB of the specific key.
func (rsdb *TendermintStateDB) GetMsgID(key []byte) (uint64, error) {
	has, err := rsdb.db.Has(key)
	if err != nil {
		return 0, err
	}

	// new node, return default.
	if !has {
		return 0, nil
	}

	var id uint64
	enc, _ := rsdb.db.Get(key)
	if len(enc) == 0 {
		return 0, nil
	}
	if err := rlp.DecodeBytes(enc, &id); err != nil {
		return 0, err
	}
	return id, nil
}

// RoundMsgsFromDB retrieves the entire round messages of the last consensus view flushed in the WAL. This function is
// called once at the node start up to rebuild the tendermint state.
func (rsdb *TendermintStateDB) RoundMsgsFromDB() *message.Map {
	roundMsgs := message.NewMap()
	if rsdb.lastConsensusMsgID == 0 {
		return roundMsgs
	}

	for id := uint64(1); id <= rsdb.lastConsensusMsgID; id++ {
		msg, verified, err := rsdb.GetMsg(id)
		if err != nil {
			rsdb.logger.Warn("failed to read WAL for msg", "error", err)
			continue
		}
		switch {
		case msg.Code() == message.ProposalCode:
			roundMsgs.GetOrCreate(msg.R()).SetProposal(msg.(*message.Propose), verified)
		case msg.Code() == message.PrevoteCode:
			roundMsgs.GetOrCreate(msg.R()).AddPrevote(msg.(*message.Prevote))
		case msg.Code() == message.PrecommitCode:
			roundMsgs.GetOrCreate(msg.R()).AddPrecommit(msg.(*message.Precommit))
		}
	}
	return roundMsgs
}

// GarbageCollection do a fast garbage collection: It is OK for overwriting a key-value pair in level DB since
// the storage engine handles the garbage collection for the overwritten value, thus we save the max msg ID in the
// store to do a periodical garbage collection: the key-values of range (lastMsgID, maxMsgID] are the only required
// items to be deleted.
func (rsdb *TendermintStateDB) GarbageCollection() {
	if rsdb.maxMsgID <= rsdb.lastConsensusMsgID {
		return
	}

	batch := rsdb.db.NewBatch()
	for i := rsdb.lastConsensusMsgID + 1; i <= rsdb.maxMsgID; i++ {
		msgKey := messageKey(i)
		if err := batch.Delete(msgKey); err != nil {
			rsdb.logger.Warn("delete msg from WAL failed", "err", err)
		}
	}
	if err := batch.Write(); err != nil {
		rsdb.logger.Warn("delete msg from WAL failed", "err", err)
		return
	}
	rsdb.maxMsgID = rsdb.lastConsensusMsgID
}

// messageKey = messagePrefix + MsgID (uint64 big endian)
func messageKey(msgID uint64) []byte {
	enc := make([]byte, 8)
	binary.BigEndian.PutUint64(enc, msgID)
	return append(messagePrefix, enc...)
}
