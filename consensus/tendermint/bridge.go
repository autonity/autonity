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
	dlog := newDebugLog("address", address.String()[2:6])
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
		broadcaster:         broadcaster,
		syncer:              syncer,
		blockReader:         blockRetreiver,
		timeoutScheduler:    timeoutScheduler,

		eventChannel:     make(chan interface{}),
		commitChannel:    make(chan *types.Block),
		closeChannel:     make(chan struct{}),
		autonityContract: ac,
		wg:               &sync.WaitGroup{},
	}
	b.syncer.Start()
	b.currentBlockAwaiter.start()
	b.wg.Add(1)
	go b.mainEventLoop()
	return b
}

// Bridge acts as a intermediary between the tendermint algorithm and the go
// ethereum system. Internally it starts up one long running go-routine for
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
// tendermint algorithm and potentially emitting a message of our own if our
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
// sync is fairly basic. The sync request message contains no information, and
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
	// proposer is the address of the proposer for the current height and round
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
	broadcaster      Broadcaster
	syncer           Syncer
	blockBroadcaster consensus.Broadcaster

	// mutext protects the fields below
	mutex        sync.RWMutex
	stopped      bool
	closeChannel chan struct{}
	wg           *sync.WaitGroup
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

// a sentinal type to indicate that we have a new chain head
type newChainHead struct{}

// NewChainHead implements consensus.Handler.NewChainHead
func (b *Bridge) NewChainHead() error {
	b.postEvent(newChainHead{})
	return nil
}

// postEvent posts an event to the main handler if Bridge is started and has a
// broadcaster, otherwise the event is dropped. This is to prevent an event
// buildup when Bridge is stopped, since the ethereum code that passes messages
// to the Bridge seems to be unaware of whether the Bridge is in a position to
// handle them.
func (b *Bridge) postEvent(e interface{}) {
	b.mutex.RLock()
	defer b.mutex.RUnlock()
	if b.stopped {
		return // Drop event if stopped
	}

	// start := time.Now()
	// b.dlog.print("posting event", fmt.Sprintf("%T", e))
	b.wg.Add(1)
	go func() {
		defer b.wg.Done()
		// I'm seeing a buildup of events here, I guess because the main
		// routine is blocked waiting for a value and so its not
		// processing these message events.
		select {
		case b.eventChannel <- e:
			// since := time.Since(start)
			// if since > time.Second {
			// 	// b.dlog.print("eventCh send took", since, "event", fmt.Sprintf("%T", e))
			// }
		case <-b.closeChannel:
			// since := time.Since(start)
			// b.dlog.print("eventCh send, stopped by closeCh, took", since/time.Second, "seconds", "event", fmt.Sprintf("%T", e))
		}
	}()
}

// SetExtraComponents must be called before the ethereum service is started,
// this is not ideal but is the best I think we can do without re-writing the
// core of go-ethereum. We end up having to do this because go-etherum itself
// is quite tangled and there is no easy way to access just the functionality
// we need.  In this case the blockchain and broadcaster both need to be
// constructed with a reference to the bridge. So we build the bridge then
// build the blockchain and broadcaster and then call this.
func (b *Bridge) SetExtraComponents(blockchain *core.BlockChain, broadcaster consensus.Broadcaster) {
	b.blockBroadcaster = broadcaster
	b.blockchain = blockchain
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
			// b.dlog.print("commitCh send, stopped by closeCh", bid(block))
		}
	} else {
		b.blockBroadcaster.Enqueue("tendermint", block)
	}

	b.logger.Info("committed a block", "hash", block.Hash())
	return nil
}

var errStopped = errors.New("stopped")

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
		// println("msghash", common.BytesToHash(crypto.Keccak256(msg)).String()[2:6])

		// send to self
		go b.postEvent(msg)
		// Broadcast to peers
		b.broadcaster.Broadcast(msg)
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

				// println("handling current height message", m.consensusMessage.String())
				err = b.handleCurrentHeightMessage(m)
				if err == errStopped {
					return
				}
				if err != nil {
					b.logger.Debug("core.mainEventLoop problem processing message", "err", err)
					continue
				}
				// Re-broadcast the message if it is not a message from ourselves,
				// if it is a message from ourselves we will have already
				// broadcast it.
				if m.address != b.address {
					b.broadcaster.Broadcast(e)
				}
			case *algorithm.Timeout:
				var cm *algorithm.ConsensusMessage
				var rc *algorithm.RoundChange
				switch e.TimeoutType {
				case algorithm.Propose:
					// b.dlog.print("timeout propose", "height", e.Height, "round", e.Round)
					cm = b.algo.OnTimeoutPropose(e.Height, e.Round)
				case algorithm.Prevote:
					// b.dlog.print("timeout prevote", "height", e.Height, "round", e.Round)
					cm = b.algo.OnTimeoutPrevote(e.Height, e.Round)
				case algorithm.Precommit:
					// b.dlog.print("timeout precommit", "height", e.Height, "round", e.Round)
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
			return fmt.Errorf("received message from non proposer: %v", m)
		}
		// Proposal values are allowed to be invalid.
		_, err := b.Verifier.VerifyProposal(*b.msgStore.value(common.Hash(cm.Value)), b.blockchain, b.address.String())
		if err == nil {
			b.msgStore.setValid(common.Hash(cm.Value))
		}
	}

	rc, cm, to := b.algo.ReceiveMessage(cm)
	return b.handleResult(rc, cm, to)
}

func (b *Bridge) proposerAddr(previousHeader *types.Header, round int64) (common.Address, error) {
	state, err := b.blockReader.BlockState(previousHeader.Root)
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
	// log := append(d.prefix, info...)
	// fmt.Printf("%v %v", time.Now().Format(time.RFC3339Nano), fmt.Sprintln(log...))
}

func bid(b *types.Block) string {
	return fmt.Sprintf("hash: %v, number: %v", b.Hash().String()[2:8], b.Number().String())
}

// TimeoutScheduler is an interface that can be used to schedule actions after some delay.
type TimeoutScheduler interface {
	ScheduleTimeout(delay uint, f func())
}

type DefaultTimeoutScheduler struct{}

func (s *DefaultTimeoutScheduler) ScheduleTimeout(delay uint, f func()) {
	time.AfterFunc(time.Duration(delay)*time.Second, f)
}
