package tendermint

import (
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/clearmatics/autonity/autonity"
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus"
	"github.com/clearmatics/autonity/consensus/tendermint/algorithm"
	"github.com/clearmatics/autonity/consensus/tendermint/bft"
	"github.com/clearmatics/autonity/consensus/tendermint/config"
	"github.com/clearmatics/autonity/core"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/crypto"
	"github.com/clearmatics/autonity/log"
	"github.com/clearmatics/autonity/p2p"
	"github.com/clearmatics/autonity/rpc"
	"github.com/davecgh/go-spew/spew"
)

var (
	// errNotFromProposer is returned when received message is supposed to be from
	// proposer.
	errNotFromProposer = errors.New("message does not come from proposer")
)

func addr(a common.Address) string {
	return hex.EncodeToString(a[:3])
}

// New creates an Tendermint consensus core
func New(
	config *config.Config,
	key *ecdsa.PrivateKey,
	broadcaster Broadcaster,
	syncer Syncer,
	verifier *Verifier,
	finalizer *DefaultFinalizer,
	blockRetreiver *BlockReader,
	ac *autonity.Contract,
) *Bridge {
	address := crypto.PubkeyToAddress(key.PublicKey)
	logger := log.New("addr", address.String())
	dlog := newDebugLog("address", address.String()[2:6])
	messageBounds := &bounds{
		centre: 0,
		high:   5,
		low:    5,
	}
	c := &Bridge{
		Verifier:             verifier,
		DefaultFinalizer:     finalizer,
		key:                  key,
		blockPeriod:          config.BlockPeriod,
		address:              address,
		logger:               logger,
		dlog:                 dlog,
		currentBlockAwaiter:  newBlockAwaiter(dlog),
		msgStore:             newMessageStore(messageBounds),
		broadcaster:          broadcaster,
		syncer:               syncer,
		latestBlockRetriever: blockRetreiver,
		verifier:             verifier,

		eventChannel:     make(chan interface{}),
		commitChannel:    make(chan *types.Block),
		autonityContract: ac,
		wg:               &sync.WaitGroup{},
	}
	return c
}

type Bridge struct {
	*DefaultFinalizer
	*Verifier

	key         *ecdsa.PrivateKey
	blockPeriod uint64
	address     common.Address
	logger      log.Logger

	eventChannel chan interface{}
	wg           *sync.WaitGroup

	msgStore  *messageStore
	syncTimer *time.Timer

	lastHeader *types.Header
	proposer   common.Address

	autonityContract *autonity.Contract

	height uint64
	round  uint64
	algo   *algorithm.Algorithm

	currentBlockAwaiter *blockAwaiter

	broadcaster          Broadcaster
	syncer               Syncer
	latestBlockRetriever *BlockReader

	verifier *Verifier

	blockchain *core.BlockChain

	blockBroadcaster consensus.Broadcaster

	mutex   sync.RWMutex
	started bool

	// Used to propagate blocks to the results channel provided by the miner on
	// calls to Seal.
	commitChannel chan *types.Block
	closeChannel  chan struct{}

	dlog *debugLog
}

func (b *Bridge) SealHash(header *types.Header) common.Hash {
	return types.SigHash(header)
}

// Author retrieves the Ethereum address of the account that minted the given
// block, which may be different from the header's coinbase if a consensus
// engine is based on signatures.
func (b *Bridge) Author(header *types.Header) (common.Address, error) {
	return types.Ecrecover(header)
}

// CalcDifficulty is the difficulty adjustment algorithm. It returns the difficulty
// that a new block should have based on the previous blocks in the blockchain and the
// current signer.
func (b *Bridge) CalcDifficulty(chain consensus.ChainHeaderReader, time uint64, parent *types.Header) *big.Int {
	return big.NewInt(1)
}

// Prepare initializes the consensus fields of a block header according to the
// rules of a particular engine. The changes are executed inline.
func (b *Bridge) Prepare(chain consensus.ChainHeaderReader, header *types.Header) error {
	// unused fields, force to set to empty
	header.Coinbase = b.address
	header.Nonce = emptyNonce
	header.MixDigest = types.BFTDigest

	// copy the parent extra data as the header extra data
	number := header.Number.Uint64()
	parent := chain.GetHeader(header.ParentHash, number-1)
	if parent == nil {
		return consensus.ErrUnknownAncestor
	}
	// use the same difficulty for all blocks
	header.Difficulty = defaultDifficulty

	// set header's timestamp
	header.Time = new(big.Int).Add(big.NewInt(int64(parent.Time)), new(big.Int).SetUint64(b.blockPeriod)).Uint64()
	if int64(header.Time) < time.Now().Unix() {
		header.Time = uint64(time.Now().Unix())
	}
	return nil
}

// APIs returns the RPC APIs this consensus engine provides.
func (b *Bridge) APIs(chain consensus.ChainReader) []rpc.API {
	return []rpc.API{{
		Namespace: "tendermint",
		Version:   "1.0",
		Service:   NewAPI(chain, b.autonityContract, b.latestBlockRetriever),
		Public:    true,
	}}
}

// So this method is meant to allow interrupting of mining a block to start on
// a new block, it doesn't make sense for autonity though because if we are not
// the proposer then we don't need this unsigned block, and if we are the
// proposer we only want the one unsigned block per round since we can't send
// multiple differing proposals.
//
// So we want to have just the latest block available to be taken from here when this node becomes the proposer.
//
// The miner only has one results channel for its lifetime and we will only
// have one miner so we can capture the results channel on the first call and
// then not worry about it after that.
//
// We can't build the bridge with the results chan since the worker will need
// the bridge to be constructed. We could create the results chan before
// building either and pass it to both. But lets save that for later.
func (b *Bridge) Seal(chain consensus.ChainReader, block *types.Block, results chan<- *types.Block, stop <-chan struct{}) error {

	// Check if we are handling the results and if not set up a goroutine to
	// pass results back to the miner. We will only send a block on the
	// commitChannel if we are the proposer.
	//
	// Ok I think I'm understanding the problem here better now. We can be in a
	// situation where the provided block has been proposed and is undergoing
	// agreement,and the miner can interrupt mining of that block to provide
	// another block at the same height, we may never propose that block if the
	// currently proposed block achieves agreement. And then when the currently
	// proposed block does reach agreement we will receive it on the commit
	// channel and it will not match the block most recently passed to Seal.
	//
	// In fact the interrupting block does not even need to be at the same
	// height, because some other network participant may have agreed the
	// currently proposed block before us and as such the miner may have
	// received a NewChainHead event and called Seal with a block for the next
	// Height. In this scenario we will want to pass the currently proposed
	// block back to the miner when it is committed even though the last
	// request to seal was for the next height. If we happen to also be
	// proposer for the next height we will also want to pass that block back
	// to the proposer when it is committed.
	//
	// The result of this will be that we can receive a block on the
	// commitChannel that does not match the block most recently passed to Seal
	// and also that we may pass multiple blocks back to the miner in the
	// lifetime of the goroutine that reads the commitChannel. So we must not
	// exit the commitChannel when sending a block, instead we must wait for
	// the miner's signal to stop or exit when we close the bridge.
	b.wg.Add(1)
	go func() {
		defer b.wg.Done()
		for {
			// b.dlog.print("commitCh receive start, block", bid(block))
			select {
			case committedBlock := <-b.commitChannel:
				// b.dlog.print("commitCh receive done", bid(committedBlock))
				results <- committedBlock
				// stop will be closed whenever eth is shutdouwn or a new
				// sealing task is provided.
			case <-stop:
				b.dlog.print("commitCh receive, stopped by miner", bid(block))
				return
			case <-b.closeChannel:
				b.dlog.print("commitCh receive, stopped by closeCh", bid(block))
				return
			}
		}
	}()

	// update the block header and signature and propose the block to core engine
	header := block.Header()

	parent := chain.GetHeader(header.ParentHash, header.Number.Uint64()-1)
	if parent == nil {
		b.logger.Error("Error ancestor")
		return consensus.ErrUnknownAncestor
	}
	nodeAddress := b.address
	if parent.CommitteeMember(nodeAddress) == nil {
		b.logger.Error("error validator errUnauthorized", "addr", b.address)
		return errUnauthorized
	}

	// wait for the timestamp of header, use this to adjust the block period
	delay := time.Until(time.Unix(int64(block.Header().Time), 0))
	select {
	case <-time.After(delay):
		// nothing to do
	case <-stop:
		return nil
	case <-b.closeChannel:
		return nil
	}

	// b.dlog.print("setting value", bid(block), "current height", b.height.String())
	b.currentBlockAwaiter.setValue(block)
	return nil
}

// Methods for consensus.Handler: This interface was introduced by the istanbul
// BFT fork, so we don't need to keep it to maintain some level of parity
// between Autonity and go-ethereum.

// Protocol implements consensus.Handler.Protocol
func (b *Bridge) Protocol() (protocolName string, extraMsgCodes uint64) {
	return "tendermint", 2 //nolint
}

// HandleMsg implements consensus.Handler.HandleMsg, this returns a bool to
// indicate whether the message was handled, if we return false then the
// message will be passed on by the caller to be handled by the default eth
// handler. If this function returns an error then the connection to the peer
// sending the message will be dropped.
func (b *Bridge) HandleMsg(addr common.Address, msg p2p.Msg) (bool, error) {
	switch msg.Code {
	case tendermintMsg:
		var data []byte
		if err := msg.Decode(&data); err != nil {
			return true, fmt.Errorf("failed to decode tendermint message: %v", err)
		}
		b.postEvent(data)
		return true, nil
	case tendermintSyncMsg:
		b.postEvent(addr)
		return true, nil
	default:
		return false, nil
	}
}

// postEvent posts an event to the main handler if Bridge is started and has a
// broadcaster, otherwise the event is dropped. This is to prevent an event
// buildup when Bridge is stopped, since the ethereum code that passes messages
// to the Bridge seems to be unaware of whether the Bridge is in a position to
// handle them.
func (b *Bridge) postEvent(e interface{}) {
	b.mutex.RLock()
	if !b.started {
		b.mutex.RUnlock()
		return // Drop event if not ready
	}
	b.mutex.RUnlock()

	start := time.Now()
	// b.dlog.print("posting event", fmt.Sprintf("%T", e))
	b.wg.Add(1)
	go func() {
		defer b.wg.Done()
		// I'm seeing a buildup of events here, I guess because the main
		// routine is blocked waiting for a value and so its not
		// processing these message events.
		select {
		case b.eventChannel <- e:
			since := time.Since(start)
			if since > time.Second {
				b.dlog.print("eventCh send took", since, "event", fmt.Sprintf("%T", e))
			}
		case <-b.closeChannel:
			since := time.Since(start)
			b.dlog.print("eventCh send, stopped by closeCh, took", since/time.Second, "seconds", "event", fmt.Sprintf("%T", e))
		}
	}()
}

// SetExtraComponents must be called before Start, this is not ideal but is the best I think
// we can do without re-writing the core of go-ethereum. We end up having to do
// this because go-etherum itself is quite tangled and there is no easy way to
// access just the functionality we need.
func (b *Bridge) SetExtraComponents(blockchain *core.BlockChain, broadcaster consensus.Broadcaster) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	b.blockBroadcaster = broadcaster
	b.blockchain = blockchain
}

type commitEvent struct{}

// NewChainHead implements consensus.Handler.NewChainHead
func (b *Bridge) NewChainHead() error {
	b.postEvent(commitEvent{})
	return nil
}

func (b *Bridge) commit(proposal *algorithm.ConsensusMessage) error {
	committedSeals := b.msgStore.signatures(proposal.Value, proposal.Round, proposal.Height)
	message := b.msgStore.matchingProposal(proposal)
	// Sanity checks
	if message == nil || message.value == nil {
		return fmt.Errorf("attempted to commit nil block")
	}
	if message.proposerSeal == nil {
		return fmt.Errorf("attempted to commit block without proposer seal")
	}
	if proposal.Round < 0 {
		return fmt.Errorf("attempted to commit a block in a negative round: %d", proposal.Round)
	}
	if len(committedSeals) == 0 {
		return fmt.Errorf("attempted to commit block without any committed seals")
	}

	for _, seal := range committedSeals {
		if len(seal) != types.BFTExtraSeal {
			return fmt.Errorf("attempted to commit block with a committed seal of invalid length: %s", hex.EncodeToString(seal))
		}
	}
	// Add the proposer seal coinbase and committed seals into the block.
	h := message.value.Header()
	h.CommittedSeals = committedSeals
	h.ProposerSeal = message.proposerSeal
	h.Coinbase = message.address
	h.Round = uint64(proposal.Round)
	block := message.value.WithSeal(h)

	if b.address == b.proposer {
		// b.dlog.print("commitCh send start", bid(block))
		select {
		case b.commitChannel <- block:
			// b.dlog.print("commitCh send done", bid(block))
		// Close channel must exist at this point (there is no way to reach
		// this without calling Start) no need for mutex.
		case <-b.closeChannel:
			b.dlog.print("commitCh send, stopped by closeCh", bid(block))
		}
	} else {
		b.blockBroadcaster.Enqueue("tendermint", block)
	}

	b.logger.Info("committed a block", "hash", block.Hash())
	return nil
}

var errStopped = errors.New("stopped")

// Start implements core.Tendermint.Start
func (b *Bridge) Start() error {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	if b.started {
		return fmt.Errorf("bridge %s started twice", b.address.String())
	}
	b.started = true
	b.closeChannel = make(chan struct{})

	b.syncer.Start()
	b.currentBlockAwaiter.start()
	// Tendermint Finite State Machine discrete event loop
	b.wg.Add(1)
	go b.mainEventLoop()
	return nil
}

func (b *Bridge) Close() error {
	b.mutex.Lock()
	if !b.started {
		b.mutex.Unlock()
		return fmt.Errorf("bridge %s closed twice", b.address.String())
	}
	b.started = false

	close(b.closeChannel)
	// println(addr(b.address), b.height, "stopping")

	// b.logger.Info("closing tendermint.Bridge", "addr", addr(b.address))

	b.syncer.Stop()
	// stop the block awaiter if it is waiting
	b.currentBlockAwaiter.stop()
	// println(addr(c.address), c.height, "almost stopped")
	// Ensure all event handling go routines exit
	b.mutex.Unlock()
	b.wg.Wait()
	return nil
}

func (b *Bridge) newHeight(prevBlock *types.Block) error {
	b.syncTimer = time.NewTimer(20 * time.Second)
	b.lastHeader = prevBlock.Header()
	b.height = prevBlock.NumberU64() + 1
	proposeValue, err := b.UpdateProposer(b.lastHeader, 0)
	if err != nil {
		return fmt.Errorf("failed to update proposer: %v", err)
	}

	// Update the height in the message store, this will clean out old messages.
	b.msgStore.setHeight(b.height)

	// Create new oracle and algorithm
	b.algo = algorithm.New(algorithm.NodeID(b.address), newOracle(b.lastHeader, b.msgStore, b.currentBlockAwaiter))

	// Handle messages for the new height
	msg, timeout := b.algo.StartRound(proposeValue, 0)

	// Note that we don't risk entering an infinite loop here since
	// start round can only return results with broadcasts or timeouts.
	err = b.handleResult(nil, msg, timeout)
	if err != nil {
		return err
	}
	// First we need to filter out messages from non committee members. This is
	// so that they do not interfere with voting calculations.
	// checkFromCommittee will remove non committee member messages from the
	// store.
	for _, msg := range b.msgStore.heightMessages(b.height) {
		err := b.checkFromCommittee(msg)
		if err != nil {
			b.logger.Error(err.Error())
		}
	}
	// Now we process the remaining messages
	for _, msg := range b.msgStore.heightMessages(b.height) {
		err := b.handleCurrentHeightMessage(msg)
		if err != nil {
			b.logger.Error("failed to handle current height message", "message", msg.String(), "err", err)
		}
	}
	return nil
}

func (b *Bridge) handleResult(rc *algorithm.RoundChange, cm *algorithm.ConsensusMessage, to *algorithm.Timeout) error {

	switch {
	case rc == nil && cm == nil && to == nil:
		return nil
	case rc != nil:
		if rc.Round == 0 && rc.Decision == nil {
			panic("round changes of 0 must be accompanied with a decision")
		}
		if rc.Decision != nil {
			// A decision has been reached
			//println(addr(c.address), "decided on block", rc.Decision.Height,common.Hash(rc.Decision.Value).String())

			// This will ultimately lead to a commit event, which we will pick up on in the mainEventLoop and start a
			// move to the new height by calling newHeight().
			err := b.commit(rc.Decision)
			if err != nil {
				panic(fmt.Sprintf("%s Failed to commit sr.Decision: %s err: %v", algorithm.NodeID(b.address).String(), spew.Sdump(rc.Decision), err))
			}
		} else {
			// We are just changing round
			var err error
			// Update the proposer
			proposalValue, err := b.UpdateProposer(b.lastHeader, rc.Round)
			if err != nil {
				return fmt.Errorf("failed to update proposer: %v", err)
			}
			startCM, startTO := b.algo.StartRound(proposalValue, rc.Round)
			// Note that we don't risk entering an infinite loop here since
			// start round can only return results with broadcasts or timeouts.
			err = b.handleResult(nil, startCM, startTO)
			if err != nil {
				return err
			}
		}
	case cm != nil:
		//println(addr(c.address), c.height.String(), cm.String(), "sending")
		// Broadcasting ends with the message reaching us eventually

		// We must build message here since buildMessage relies on accessing
		// the msg store, and since the message store is not syncronised we
		// need to do it from the handler routine.
		msg, err := encodeSignedMessage(cm, b.key, b.msgStore.value(common.Hash(cm.Value)))
		if err != nil {
			panic(fmt.Sprintf(
				"%s We were unable to build a message, this indicates a programming error: %v",
				addr(b.address),
				err,
			))
		}
		b.dlog.print("sending message", cm.String())
		println("msghash", common.BytesToHash(crypto.Keccak256(msg)).String()[2:6])

		// Broadcast in a new goroutine
		b.wg.Add(1)
		go func() {
			defer b.wg.Done()
			// send to self
			b.postEvent(msg)
			// Broadcast to peers
			b.broadcaster.Broadcast(msg)
		}()
	case to != nil:
		time.AfterFunc(time.Duration(to.Delay)*time.Second, func() {
			b.postEvent(to)
		})

	}
	return nil
}

func (b *Bridge) mainEventLoop() {
	defer b.wg.Done()

	lastBlockMined, err := b.latestBlockRetriever.LatestBlock()
	if err != nil {
		panic(err)
	}
	err = b.newHeight(lastBlockMined)
	if err != nil {
		return
	}

	// Ask for sync when the engine starts
	b.syncer.AskSync(b.lastHeader)

	lastHeight := b.height

eventLoop:
	for {
		select {
		case <-b.syncTimer.C:
			if lastHeight == b.height {
				b.dlog.print("syncing")
				b.syncer.AskSync(b.lastHeader)
			}
			lastHeight = b.height
			b.syncTimer = time.NewTimer(20 * time.Second)

		case ev := <-b.eventChannel:
			switch e := ev.(type) {
			case common.Address:
				b.logger.Info("Processing sync message", "from", e)
				b.syncer.SyncPeer(e, b.msgStore.rawHeightMessages(b.height))
			case []byte:
				/*
					Basic validity checks
				*/

				m, err := decodeSignedMessage(e)
				if err != nil {
					fmt.Printf("some error: %v\n", err)
					continue
				}
				err = b.msgStore.addMessage(m, e)
				if err != nil {
					// could be multiple proposal messages from the same proposer
					continue
				}
				if m.consensusMessage.MsgType == algorithm.Propose {
					b.msgStore.addValue(m.value.Hash(), m.value)
				}

				// If this message is for a future height then we cannot validate it
				// because we lack the relevant header, we will process it when we reach
				// that height. If it is for a previous height then we are not intersted in
				// it. But it has been added to the message store in case other peers would
				// like to sync it.
				if m.consensusMessage.Height != b.height {
					// Nothing to do here
					continue
				}

				println("handling current height message", m.consensusMessage.String())
				err = b.handleCurrentHeightMessage(m)
				if err == errStopped {
					return
				}
				if err != nil {
					b.logger.Debug("core.mainEventLoop problem processing message", "err", err)
					continue
				}
			case *algorithm.Timeout:
				var cm *algorithm.ConsensusMessage
				var rc *algorithm.RoundChange
				switch e.TimeoutType {
				case algorithm.Propose:
					b.dlog.print("timeout propose", "height", e.Height, "round", e.Round)
					cm = b.algo.OnTimeoutPropose(e.Height, e.Round)
				case algorithm.Prevote:
					b.dlog.print("timeout prevote", "height", e.Height, "round", e.Round)
					cm = b.algo.OnTimeoutPrevote(e.Height, e.Round)
				case algorithm.Precommit:
					b.dlog.print("timeout precommit", "height", e.Height, "round", e.Round)
					rc = b.algo.OnTimeoutPrecommit(e.Height, e.Round)
				}
				err := b.handleResult(rc, cm, nil)
				if err != nil {
					b.dlog.print("exiting main event loop", "height", e.Height, "round", e.Round, "err", err.Error())
					return
				}
			case commitEvent:
				b.logger.Debug("Received a final committed proposal")

				lastBlock, err := b.latestBlockRetriever.LatestBlock()
				if err != nil {
					panic(err)
				}
				b.dlog.print("commit event for block", bid(lastBlock))
				err = b.newHeight(lastBlock)
				if err != nil {
					return
				}
			}
		case <-b.closeChannel:
			b.dlog.print("mainEventLoop, stopped by closeCh, current height", b.height)
			b.logger.Info("Bridge closed, exiting mainEventLoop")
			break eventLoop
		}
	}

}

func (b *Bridge) checkFromCommittee(m *message) error {
	// Check that the message came from a committee member, if not remove it from the store and return an error.
	if b.lastHeader.CommitteeMember(m.address) == nil {
		// We remove the message from the store since it came from a non
		// validator.
		b.msgStore.removeMessage(m)

		// TODO turn this into an error type that can be checked for at a
		// higher level to close the connection to this peer.
		return fmt.Errorf("received message from non committee member: %v", m)
	}
	return nil
}

func (b *Bridge) handleCurrentHeightMessage(m *message) error {
	cm := m.consensusMessage
	/*
		Domain specific validity checks, now we know that we are at the same
		height as this message we can rely on lastHeader.
	*/

	err := b.checkFromCommittee(m)
	if err != nil {
		return err
	}
	if cm.MsgType == algorithm.Propose {
		// We ignore proposals from non proposers
		if b.proposer != m.address {
			b.logger.Warn("Ignore proposal messages from non-proposer")
			return errNotFromProposer
		}
		// Proposal values are allowed to be invalid.
		_, err := b.verifier.VerifyProposal(*b.msgStore.value(common.Hash(cm.Value)), b.blockchain, b.address.String())
		if err == nil {
			b.msgStore.setValid(common.Hash(cm.Value))
		} else {
			println("not valid", err.Error())
		}
	}

	rc, cm, to := b.algo.ReceiveMessage(cm)
	return b.handleResult(rc, cm, to)
}

func (b *Bridge) proposerAddr(previousHeader *types.Header, round int64) (common.Address, error) {
	state, err := b.latestBlockRetriever.BlockState(previousHeader.Root)
	if err != nil {
		return common.Address{}, fmt.Errorf("cannot load state from block chain: %v", err)
	}
	return b.autonityContract.GetProposerFromAC(previousHeader, state, round)
}

// UpdateProposer updates b.proposer and if we are the proposer waits for a
// proposal value and returns an algorithm.ValueID representing the proposal
// value. If we are not the proposer then it returns algorithm.NilValue.
func (b *Bridge) UpdateProposer(previousHeader *types.Header, round int64) (algorithm.ValueID, error) {
	var err error
	b.proposer, err = b.proposerAddr(previousHeader, round)
	if err != nil {
		return algorithm.NilValue, fmt.Errorf("cannot load state from block chain: %v", err)
	}
	// If we are not the proposer then return nil value.
	if b.address != b.proposer {
		return algorithm.NilValue, nil
	}
	v, err := b.currentBlockAwaiter.value(b.height)
	if err != nil {
		return algorithm.NilValue, fmt.Errorf("failed to get value: %v", err)
	}
	// Add the value to the store, we do not mark it valid here since we
	// will validate it when whe process our own proposal.
	b.msgStore.addValue(v.Hash(), v)
	return algorithm.ValueID(v.Hash()), nil
}

const (
	tendermintMsg     = 0x11
	tendermintSyncMsg = 0x12
)

// TODO need to clear this out, ideally when a peer disconnects and when we stop
// caring about the tracked messages. So really we need a notion of height to
// be worked in here.
type peerMessageMap interface {
	// knowsMessage returns true if the peer knows the current message
	knowsMessage(addr common.Address, hash common.Hash) bool
}

// TODO actually implement this
type degeneratePeerMessageMap struct {
}

func (p *degeneratePeerMessageMap) knowsMessage(_ common.Address, _ common.Hash) bool {
	return false
}

type Broadcaster interface {
	Broadcast(payload []byte)
}

type DefaultBroadcaster struct {
	address common.Address
	pmm     peerMessageMap
	peers   consensus.Peers
}

func NewBroadcaster(address common.Address, peers consensus.Peers) *DefaultBroadcaster {
	return &DefaultBroadcaster{
		address: address,
		peers:   peers,
		pmm:     &degeneratePeerMessageMap{},
	}
}

// Broadcast implements tendermint.Backend.Broadcast
func (b *DefaultBroadcaster) Broadcast(payload []byte) {
	hash := types.RLPHash(payload)

	for _, p := range b.peers.Peers() {
		if !b.pmm.knowsMessage(p.Address(), hash) {
			// TODO make sure we update the peerMessageMap with the sent
			// message, once successfully sent. previously we were updating
			// the map before trying to send the message so if message
			// sending failed we would not have tried again.
			go p.Send(tendermintMsg, payload) //nolint
		}
	}
}

type Syncer interface {
	Start()
	Stop()
	AskSync(lastestHeader *types.Header)
	SyncPeer(peerAddr common.Address, messages [][]byte)
}

type DefaultSyncer struct {
	address common.Address
	peers   consensus.Peers
	stopped chan struct{}
	mu      sync.Mutex
}

func NewSyncer(peers consensus.Peers, address common.Address) *DefaultSyncer {
	return &DefaultSyncer{
		peers:   peers,
		address: address,
	}
}

func (s *DefaultSyncer) Start() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.stopped = make(chan struct{})
}
func (s *DefaultSyncer) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()
	close(s.stopped)
}

func (s *DefaultSyncer) AskSync(latestHeader *types.Header) {
	var count uint64

	// Determine if there should be any other peers
	potentialPeerCount := len(latestHeader.Committee)
	if latestHeader.CommitteeMember(s.address) != nil {
		// Remove ourselves from the other potential peers
		potentialPeerCount--
	}
	// Exit if there are no other peers
	if potentialPeerCount == 0 {
		return
	}
	peers := s.peers.Peers()
	// Wait for there to be peers
	for len(peers) == 0 {
		t := time.NewTimer(10 * time.Millisecond)
		select {
		case <-t.C:
			peers = s.peers.Peers()
			continue
		case <-s.stopped:
			return
		}
	}

	// Ask for sync to peers
	for _, p := range peers {
		//ask to a quorum nodes to sync, 1 must then be honest and updated
		if count >= bft.Quorum(latestHeader.TotalVotingPower()) {
			break
		}
		go p.Send(tendermintSyncMsg, []byte{}) //nolint
		member := latestHeader.CommitteeMember(p.Address())
		if member == nil {
			continue
		}
		count += member.VotingPower.Uint64()
	}
}

// Synchronize new connected peer with current height state
func (s *DefaultSyncer) SyncPeer(address common.Address, messages [][]byte) {
	for _, p := range s.peers.Peers() {
		if address == p.Address() {
			for _, msg := range messages {
				go p.Send(tendermintMsg, msg) //nolint
			}
			break
		}
	}
}

type debugLog struct {
	prefix []interface{}
}

func newDebugLog(prefix ...interface{}) *debugLog {
	return &debugLog{
		prefix: prefix,
	}
}

func (d *debugLog) print(info ...interface{}) {
	log := append(d.prefix, info...)
	fmt.Printf("%v %v", time.Now().Format(time.RFC3339Nano), fmt.Sprintln(log...))
}

func bid(b *types.Block) string {
	return fmt.Sprintf("hash: %v, number: %v", b.Hash().String()[2:8], b.Number().String())
}
