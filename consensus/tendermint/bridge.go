package tendermint

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/clearmatics/autonity/autonity"
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus"
	"github.com/clearmatics/autonity/consensus/tendermint/algorithm"
	"github.com/clearmatics/autonity/consensus/tendermint/config"
	"github.com/clearmatics/autonity/core"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/crypto"
	"github.com/clearmatics/autonity/log"
	"github.com/clearmatics/autonity/p2p"
	"github.com/clearmatics/autonity/rpc"
	"github.com/davecgh/go-spew/spew"
)

const (
	// TendermintMsg is the p2p Message code assigned to tendermint algorithm
	// messages.
	TendermintMsg = 0x11
	// TendermintSyncMsg is the p2p Message code assigned to tendermint
	// algorithm sync requests.
	TendermintSyncMsg = 0x12
)

var (
	// errDecodeFailed is returned when decode message fails
	errDecodeFailed = errors.New("fail to decode tendermint message")
)

// New creates a new Bridge instance.
func New(
	config *config.Config,
	key *ecdsa.PrivateKey,
	broadcaster Broadcaster,
	syncer Syncer,
	verifier *Verifier,
	finalizer *DefaultFinalizer,
	blockRetreiver *BlockReader,
	ac *autonity.Contract,
	timeoutScheduler TimeoutScheduler,
) *Bridge {
	address := crypto.PubkeyToAddress(key.PublicKey)
	logger := log.New("addr", address.String())
	dlog := newDebugLog("Address", addr(address))
	messageBounds := &bounds{
		centre: 0,
		high:   5,
		low:    5,
	}
	b := &Bridge{
		Verifier:            verifier,
		DefaultFinalizer:    finalizer,
		key:                 key,
		blockPeriod:         config.BlockPeriod,
		address:             address,
		logger:              logger,
		dlog:                dlog,
		currentBlockAwaiter: newBlockAwaiter(dlog),
		msgStore:            newMessageStore(messageBounds),
		peerBroadcaster:     broadcaster,
		syncer:              syncer,
		blockReader:         blockRetreiver,
		timeoutScheduler:    timeoutScheduler,

		eventChannel:     make(chan interface{}),
		commitChannel:    make(chan *types.Block),
		closeChannel:     make(chan struct{}),
		autonityContract: ac,
		wg:               &sync.WaitGroup{},
	}
	b.currentBlockAwaiter.start()
	b.wg.Add(1)
	go b.mainEventLoop()
	return b
}

// Bridge acts as a intermediary between the tendermint algorithm and the go
// ethereum system. Internally it starts up one long running go-routine for the
// mainEventLoop. The various inputs to the bridge are serialised through a
// selection of channels such that the mainEventLoop can handle them in a
// single threaded manner. This allows us to have a straightforward
// implementation of the tendermint algorithm.
//
// The ethereum system interacts with the Bridge instance primarily through the
// following methods and objects. Stop, Seal, NewChainHead, HandleMsg and the
// methods of the embedded Finalizer and Verifier. Stop is used to close the
// bridge when a node is shutting down, Seal, NewChainHead and HandleMsg are
// how inputs arrive in the bridge, these inputs will end up being processed in
// the mainEventLoop. The finalizer and verifier define the block state
// transition function and block verification logic respectively, the ethereum
// system uses them in the fashion of a utility function, they are totally
// separate from the operation of the rest of the bridge and could live
// separately from it if it were not for the fact that ethereum bundles
// together all this functionality in 2 interfaces 'consensus.Engine' and
// 'consensus.Handler' which may as well be one interface because one is cast
// to the other.
//
// Seal -> provides blocks to come to agreement on. These are provided
// continuously by the miner, irrespective of whether we are mining or this
// instance is the proposer.
//
// NewChainHead -> is called by the miner when the miner has received a new
// block. This is what signals to the bridge that it should start a new height.
// Note, even though at the level of the bridge we know when we have committed
// a block, we are not in control of moving to the next height, we need to wait
// for the committed block to make its way back to the miner and for the miner
// to call us.
//
// HandleMsg -> is called by the ethereum protocol manager, when messages are
// received from other peers, we process these eventually passing them to the
// tendermint algorithm and potentially emitting a Message of our own if our
// state changed.
//
// The bridge interacts with the ethereum system through the Broadcaster,
// Syncer, consensus.Broadcaster, BlockReader and BlockChain. BlockReader and
// BlockChain are used to read information about the state of the ethereum
// system, Broadcaster, Syncer and consensus.Broadcaster are how the Bridge
// sends information to the ethereum system.
//
// BlockReader -> used to read blocks and block state.
//
// BlockChain -> only used for verification of headers.
//
// Broadcaster -> broadcasts consensus messages to the rest of the system.
//
// Syncer -> initiates and fulfills sync requests. Note sync in this sense is
// concerned with syncing the tendermint protocol messages, not blocks from the
// chain, that is handled in the 'eth/downloader' package. Note also that this
// sync is fairly basic. The sync request Message contains no information, and
// the sync response is always to send all the messages that a node has for the
// current height, unless nodes happen to be at the same height, the sync will
// be of no use and will simply clog up the network.
//
// consensus.Broadcaster -> broadcasts committed blocks to the rest of the system.
type Bridge struct {
	// These embedded fields provide utility functions for the ethereum system.
	*DefaultFinalizer
	*Verifier

	// These fields could be considered config
	blockPeriod uint64
	key         *ecdsa.PrivateKey
	address     common.Address
	logger      log.Logger
	dlog        *debugLog

	// These fields support the functioning of the mainEventLoop
	//
	// eventChannel serialises events into the mainEventLoop
	eventChannel chan interface{}
	// syncTimer instigates periodic sync operations
	syncTimer *time.Timer
	// msgStore keeps track of all sent and received messages within some
	// bounds defined as a range of blocks.
	msgStore *messageStore
	// lastHeader is the header of the latest confirmed block
	lastHeader *types.Header
	// proposer is the Address of the proposer for the current height and round
	proposer common.Address
	// autonityContract interfaces with the deployed autonity contract
	autonityContract *autonity.Contract
	// height is the current height
	height uint64
	// algo is the tendermint algorithm
	algo *algorithm.Algorithm
	// currentBlockAwaiter provides a mechanism to wait for a block provided by
	// the miner at a specific height.
	currentBlockAwaiter *blockAwaiter
	// blockReader provides functionality to read blocks.
	blockReader *BlockReader
	// blockchain is only used to be able to call the verifier internally to
	// verify proposals.
	blockchain *core.BlockChain
	// commitChannel propagates blocks to the results channel provided by the
	// miner on calls to Seal.
	commitChannel chan *types.Block
	// timeoutScheduler schedules timeout events.
	timeoutScheduler TimeoutScheduler

	// These 3 fields are used to communicate back to the ethereum system.
	peerBroadcaster  Broadcaster
	syncer           Syncer
	localBroadcaster consensus.Broadcaster

	// mutext protects the fields below
	mutex        sync.RWMutex
	stopped      bool
	closeChannel chan struct{}
	wg           *sync.WaitGroup
}

// Protocol implements consensus.Handler.Protocol
func (b *Bridge) Protocol() (protocolName string, extraMsgCodes uint64) {
	return "tendermint", 2 //nolint
}

// HandleMsg implements consensus.Handler.HandleMsg, this returns a byte slice to
// indicate whether p2p protocol the Message need to forward the consensus msg to fault detector,
// handled, if this function returns an error then the connection to the peer
// sending the Message will be dropped.
func (b *Bridge) HandleMsg(addr common.Address, msg p2p.Msg) ([]byte, error) {

	if msg.Code == TendermintSyncMsg {
		b.logger.Info("Received sync message", "from", addr)
		b.postEvent(addr)
	}

	if msg.Code == TendermintMsg {
		buff := new(bytes.Buffer)
		if _, err := io.Copy(buff, msg.Payload); err != nil {
			return nil, errDecodeFailed
		}
		copyPayload := make([]byte, len(buff.Bytes()))
		copy(copyPayload, buff.Bytes())

		var data []byte
		copyMsg := msg
		buf := new(bytes.Buffer)
		buf.Write(copyPayload)
		copyMsg.Payload = buf
		if err := copyMsg.Decode(&data); err != nil {
			return copyPayload, errDecodeFailed
		}

		b.postEvent(data)
		return copyPayload, nil
	}

	return nil, nil
}

// a sentinal type to indicate that we have a new chain head
type newChainHead struct{}

// NewChainHead implements consensus.Handler.NewChainHead
func (b *Bridge) NewChainHead() error {
	b.postEvent(newChainHead{})
	return nil
}

// APIs returns the RPC APIs this consensus engine provides.
func (b *Bridge) APIs(chain consensus.ChainReader) []rpc.API {
	return []rpc.API{{
		Namespace: "tendermint",
		Version:   "1.0",
		Service:   NewAPI(chain, b.autonityContract, b.blockReader),
		Public:    true,
	}}
}

// Seal implements consensus.Engine.Seal.
func (b *Bridge) Seal(chain consensus.ChainReader, block *types.Block, results chan<- *types.Block, stop <-chan struct{}) error {

	// Verify that the block looks correct and that we are part of the committee
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

	// Set up a goroutine to pass results back to the miner. We will only send
	// a block on the commitChannel if we are the proposer.
	//
	// Also note that the block received on the commitChannel may not be the
	// block that was last passed to seal. There is no guarantee that any block
	// passed to seal will actually be proposed. Since a prior block may have
	// already been proposed and this is what will be returned on the
	// commitChannel. Also note that in the lifetime of this goroutine multiple
	// blocks may be passed to the results chan. This can happen if we are the
	// proposer and we receive a block through seal which we propose, a member
	// of the network is able to commit the block before we do, and broadcasts
	// it to us. This will result in the miner closing the current stop channel
	// (thus closing our goroutine) and calling Seal with a block for the next
	// height. We set up a new goroutine for this block and now we commit the
	// previous block, which will be passed to the commitChannel, we then
	// propose and commit the next block that was passed on the most recent
	// invocation of Seal, this too will be sent on the commitChannel.

	// This protects us from calling this function after we have been stopped
	// and also from calling b.wg.Add after b.wg.Wait has been called in Close.
	b.mutex.RLock()
	if b.stopped {
		b.mutex.RUnlock()
		return nil
	}
	b.wg.Add(1)
	b.mutex.RUnlock()

	go func() {
		defer b.wg.Done()
		for {
			select {
			case committedBlock := <-b.commitChannel:
				select {
				case results <- committedBlock:
				case <-b.closeChannel:
					return
				}
			case <-stop:
				return
			case <-b.closeChannel:
				return
			}
		}
	}()

	// wait for the timestamp of header before passing on the block
	delay := time.Until(time.Unix(int64(block.Header().Time), 0))
	select {
	case <-time.After(delay):
		// nothing to do
	case <-stop:
		return nil
	case <-b.closeChannel:
		return nil
	}

	// Pass the block to the awaiter, the block will be available to to us when
	// we next need to propose.
	b.currentBlockAwaiter.setValue(block)
	return nil
}

// postEvent posts an event to the main handler if Bridge is started and has a
// peerBroadcaster, otherwise the event is dropped. This is to prevent an event
// buildup when Bridge is stopped, since the ethereum code that passes messages
// to the Bridge seems to be unaware of whether the Bridge is in a position to
// handle them.
func (b *Bridge) postEvent(e interface{}) {
	b.mutex.RLock()
	defer b.mutex.RUnlock()
	if b.stopped {
		return // Drop event if stopped
	}

	b.wg.Add(1)
	go func() {
		defer b.wg.Done()
		select {
		case b.eventChannel <- e:
		case <-b.closeChannel:
		}
	}()
}

// SetExtraComponents must be called before the ethereum service is started,
// this is not ideal but is the best I think we can do without re-writing the
// core of go-ethereum. We end up having to do this because go-etherum itself
// is quite tangled and there is no easy way to access just the functionality
// we need.  In this case the blockchain and peerBroadcaster both need to be
// constructed with a reference to the bridge. So we build the bridge then
// build the blockchain and peerBroadcaster and then call this.
func (b *Bridge) SetExtraComponents(blockchain *core.BlockChain, broadcaster consensus.Broadcaster) {
	b.localBroadcaster = broadcaster
	b.blockchain = blockchain
}

// commit takes a confirmed proposal and builds a corresponding block from that
// and either broadcasts it to the network if we are not the proposer or if we
// are sends it back to the miner via the commitChannel.
func (b *Bridge) commit(proposal *algorithm.ConsensusMessage) error {
	committedSeals := b.msgStore.signatures(proposal.Value, proposal.Round, proposal.Height)
	message := b.msgStore.matchingProposal(proposal)
	// Sanity checks
	if message == nil || message.Value == nil {
		return fmt.Errorf("attempted to commit nil block")
	}
	if message.ProposerSeal == nil {
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
	h := message.Value.Header()
	h.CommittedSeals = committedSeals
	h.ProposerSeal = message.ProposerSeal
	h.Coinbase = message.Address
	h.Round = uint64(proposal.Round)
	block := message.Value.WithSeal(h)

	if b.address == b.proposer {
		select {
		case b.commitChannel <- block:
		case <-b.closeChannel:
		}
	} else {
		b.localBroadcaster.Enqueue("tendermint", block)
	}

	b.logger.Info("committed a block", "Hash", block.Hash())
	return nil
}

// Close stops and waits for all goroutines started by the bridge to exit.
func (b *Bridge) Close() error {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	if b.stopped {
		return fmt.Errorf("bridge %s closed twice", b.address.String())
	}
	b.stopped = true

	close(b.closeChannel)
	b.syncer.Stop()
	b.currentBlockAwaiter.stop()
	// Ensure all event handling go routines exit
	b.wg.Wait()
	return nil
}

func (b *Bridge) newHeight(prevBlock *types.Block) error {
	b.syncTimer = time.NewTimer(20 * time.Second)
	b.lastHeader = prevBlock.Header()
	b.height = prevBlock.NumberU64() + 1
	proposeValue, err := b.updateProposer(b.lastHeader, 0)
	if err != nil {
		return fmt.Errorf("failed to update proposer: %v", err)
	}

	// Update the height in the Message store, this will clean out old messages.
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
			b.logger.Error("failed to handle current height Message", "Message", msg.String(), "err", err)
		}
	}
	return nil
}

// handle result handles the output of the tendermint algorithm which can
// either be nothing, in the case where no state change occurred in the
// algorithm or it can be one of round change consensus Message or timeout.
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
			proposalValue, err := b.updateProposer(b.lastHeader, rc.Round)
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
		// Broadcast the new Message to the network.

		// We must build Message here since buildMessage relies on accessing
		// the msg store, and since the Message store is not syncronised we
		// need to do it from the handler routine.
		msg, err := EncodeSignedMessage(cm, b.key, b.msgStore.value(common.Hash(cm.Value)))
		if err != nil {
			panic(fmt.Sprintf(
				"%s We were unable to build a Message, this indicates a programming error: %v",
				addr(b.address),
				err,
			))
		}
		b.dlog.print("sending Message", cm.String())

		// send to self, we process our own messages just as we process
		// messgaes from other network participants.
		go b.postEvent(msg)
		// send msg to local AFD for accountability.
		if b.localBroadcaster != nil {
			b.localBroadcaster.SendLocalMsgToAFD(msg)
		}

		// Broadcast to peers.
		//
		// Note the tests in bridge_test.go rely on calls to Broadcast
		// being done in the main handler routine.
		b.peerBroadcaster.Broadcast(msg)
	case to != nil:
		b.timeoutScheduler.ScheduleTimeout(to.Delay, func() {
			b.postEvent(to)
		})

	}
	return nil
}

func (b *Bridge) mainEventLoop() {
	defer b.wg.Done()

	lastBlockMined, err := b.blockReader.LatestBlock()
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
				b.logger.Info("Processing sync Message", "from", e)
				b.syncer.SyncPeer(e, b.msgStore.rawHeightMessages(b.height))
			case []byte:
				/*
					Basic validity checks
				*/

				m, err := DecodeSignedMessage(e)
				if err != nil {
					fmt.Printf("some error: %v\n", err)
					continue
				}
				err = b.msgStore.addMessage(m, e)
				if err != nil {
					// could be multiple proposal messages from the same proposer
					continue
				}
				if m.ConsensusMessage.MsgType == algorithm.Propose {
					b.msgStore.addValue(m.Value.Hash(), m.Value)
				}

				// If this Message is for a future height then we cannot validate it
				// because we lack the relevant header, we will process it when we reach
				// that height. If it is for a previous height then we are not intersted in
				// it. But it has been added to the Message store in case other peers would
				// like to sync it.
				if m.ConsensusMessage.Height != b.height {
					// Nothing to do here
					continue
				}

				// println("handling current height Message", m.ConsensusMessage.String())
				err = b.handleCurrentHeightMessage(m)
				if err == errStopped {
					return
				}
				if err != nil {
					b.logger.Debug("core.mainEventLoop problem processing Message", "err", err)
					continue
				}
				// Re-broadcast the Message if it is not a Message from ourselves,
				// if it is a Message from ourselves we will have already
				// broadcast it.
				if m.Address != b.address {
					// Note the tests in bridge_test.go rely on calls to Broadcast
					// being done in the main handler routine.
					b.peerBroadcaster.Broadcast(e)
				}
			case *algorithm.Timeout:
				var cm *algorithm.ConsensusMessage
				var rc *algorithm.RoundChange
				switch e.TimeoutType {
				case algorithm.Propose:
					cm = b.algo.OnTimeoutPropose(e.Height, e.Round)
				case algorithm.Prevote:
					cm = b.algo.OnTimeoutPrevote(e.Height, e.Round)
				case algorithm.Precommit:
					rc = b.algo.OnTimeoutPrecommit(e.Height, e.Round)
				}
				err := b.handleResult(rc, cm, nil)
				if err != nil {
					b.dlog.print("exiting main event loop", "height", e.Height, "round", e.Round, "err", err.Error())
					return
				}
			case newChainHead:

				lastBlock, err := b.blockReader.LatestBlock()
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

// checkFromCommittee checks that m is from a committee member and if not
// removes it from the store and returns an error.
func (b *Bridge) checkFromCommittee(m *Message) error {
	if b.lastHeader.CommitteeMember(m.Address) == nil {
		// We remove the Message from the store since it came from a non
		// validator.
		b.msgStore.removeMessage(m)

		// TODO turn this into an error type that can be checked for at a
		// higher level to close the connection to this peer.
		return fmt.Errorf("received Message from non committee member: %v", m)
	}
	return nil
}

// handleCurrentHeightMessage processes messages that are at the same height as
// the bridge, messages at a future height cannot be processed since we cannot
// know if they come from a committee member until we have committed the previous block.
func (b *Bridge) handleCurrentHeightMessage(m *Message) error {
	cm := m.ConsensusMessage

	/*
		Domain specific validity checks, now we know that we are at the same
		height as this Message we can rely on lastHeader.
	*/

	err := b.checkFromCommittee(m)
	if err != nil {
		return err
	}
	if cm.MsgType == algorithm.Propose {
		// We ignore proposals from non proposers
		if b.proposer != m.Address {
			return fmt.Errorf("received Message from non proposer: %v", m)
		}
		// Proposal values are allowed to be invalid.
		_, err := b.Verifier.VerifyProposal(*b.msgStore.value(common.Hash(cm.Value)), b.blockchain, b.address.String())
		if err == nil {
			b.msgStore.setValid(common.Hash(cm.Value))
		}
	}

	// let the algorithm receive the Message and handle the result.
	rc, cm, to := b.algo.ReceiveMessage(cm)
	return b.handleResult(rc, cm, to)
}

// proposerAddr gets the Address of the proposer given the previous header and round.
func (b *Bridge) proposerAddr(previousHeader *types.Header, round int64) (common.Address, error) {
	state, err := b.blockReader.BlockState(previousHeader.Root)
	if err != nil {
		return common.Address{}, fmt.Errorf("cannot load state from block chain: %v", err)
	}
	return b.autonityContract.GetProposerFromAC(previousHeader, state, round)
}

// updateProposer updates b.proposer and if we are the proposer waits for a
// proposal Value, stores it in the msgStore and returns an algorithm.ValueID
// representing the proposal Value. If we are not the proposer then it returns
// algorithm.NilValue.
func (b *Bridge) updateProposer(previousHeader *types.Header, round int64) (algorithm.ValueID, error) {
	var err error
	b.proposer, err = b.proposerAddr(previousHeader, round)
	if err != nil {
		return algorithm.NilValue, fmt.Errorf("cannot load state from block chain: %v", err)
	}
	// If we are not the proposer then return nil Value.
	if b.address != b.proposer {
		return algorithm.NilValue, nil
	}
	v, err := b.currentBlockAwaiter.value(b.height)
	if err != nil {
		return algorithm.NilValue, fmt.Errorf("failed to get Value: %v", err)
	}
	// Add the Value to the store, we do not mark it valid here since we
	// will validate it when whe process our own proposal.
	b.msgStore.addValue(v.Hash(), v)
	return algorithm.ValueID(v.Hash()), nil
}

// TODO need to clear this out, ideally when a peer disconnects and when we stop
// caring about the tracked messages. So really we need a notion of height to
// be worked in here.
type peerMessageMap interface {
	// knowsMessage returns true if the peer knows the current Message
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
			// Message, once successfully sent. previously we were updating
			// the map before trying to send the Message so if Message
			// sending failed we would not have tried again.
			go p.Send(TendermintMsg, payload) //nolint
		}
	}
}

// TimeoutScheduler is an interface that can be used to schedule actions after some delay.
type TimeoutScheduler interface {
	ScheduleTimeout(delay uint, f func())
}

// DefaultTimeoutScheduler schedules the action after 'delay' seconds.
type DefaultTimeoutScheduler struct{}

func (s *DefaultTimeoutScheduler) ScheduleTimeout(delay uint, f func()) {
	time.AfterFunc(time.Duration(delay)*time.Second, f)
}
