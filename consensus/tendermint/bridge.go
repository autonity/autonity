package tendermint

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"sync/atomic"
	"time"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus"
	"github.com/clearmatics/autonity/consensus/tendermint/algorithm"
	"github.com/clearmatics/autonity/consensus/tendermint/bft"
	"github.com/clearmatics/autonity/consensus/tendermint/config"
	"github.com/clearmatics/autonity/consensus/tendermint/events"
	"github.com/clearmatics/autonity/contracts/autonity"
	"github.com/clearmatics/autonity/core"
	"github.com/clearmatics/autonity/core/rawdb"
	"github.com/clearmatics/autonity/core/state"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/ethdb"
	"github.com/clearmatics/autonity/event"
	"github.com/clearmatics/autonity/log"
	"github.com/clearmatics/autonity/p2p"
	"github.com/davecgh/go-spew/spew"
)

var (
	// errNotFromProposer is returned when received message is supposed to be from
	// proposer.
	errNotFromProposer = errors.New("message does not come from proposer")
	// errInvalidMessage is returned when the message is malformed.
	errInvalidMessage = errors.New("invalid message")
)

const (
	MaxRound = 99 // consequence of backlog priority
)

func addr(a common.Address) string {
	return hex.EncodeToString(a[:3])
}

// New creates an Tendermint consensus core
func New(config *config.Config, key *ecdsa.PrivateKey, broadcaster *Broadcaster, syncer *Syncer, address common.Address, latestBlockRetreiver *LatestBlockRetriever, statedb state.Database, verifier *Verifier) *bridge {
	logger := log.New("addr", address.String())
	c := &bridge{
		key:                  key,
		proposerPolicy:       config.ProposerPolicy,
		address:              address,
		logger:               logger,
		currentBlockAwaiter:  newBlockAwaiter(),
		msgStore:             newMessageStore(),
		broadcaster:          broadcaster,
		syncer:               syncer,
		latestBlockRetreiver: latestBlockRetreiver,
		statedb:              statedb,
		verifier:             verifier,
		eventMux:             event.NewTypeMuxSilent(logger),
		commitChannel:        make(chan *types.Block),
	}
	return c
}

type bridge struct {
	key            *ecdsa.PrivateKey
	proposerPolicy config.ProposerPolicy
	address        common.Address
	logger         log.Logger

	cancel context.CancelFunc

	eventsSub    *event.TypeMuxSubscription
	syncEventSub *event.TypeMuxSubscription
	wg           *sync.WaitGroup

	msgStore  *messageStore
	syncTimer *time.Timer

	committee  committee
	lastHeader *types.Header

	autonityContract *autonity.Contract

	height *big.Int
	algo   *algorithm.Algorithm

	currentBlockAwaiter *blockAwaiter

	broadcaster          *Broadcaster
	syncer               *Syncer
	latestBlockRetreiver *LatestBlockRetriever
	statedb              state.Database

	verifier *Verifier

	blockchain *core.BlockChain // TODO need to set this on start

	eventMux         *event.TypeMuxSilent
	blockBroadcaster consensus.Broadcaster

	// 1 means started, 0 means stopped.
	started int32
	// 1 means set, 0 means not set.
	broadcasterSet int32

	// Used to propagate blocks to the results channel provided by the miner on
	// calls to Seal.
	commitChannel chan *types.Block
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
func (b *bridge) Seal(chain consensus.ChainReader, block *types.Block, results chan<- *types.Block, stop <-chan struct{}) error {

	// Check if we are handling the results and if not set up a goroutine to
	// pass results back to the miner. We will only send a block on the
	// commitChannel if we are the proposer.
	//
	// TODO I think there is a problem here that if we are the proposer and we
	// receive a future block from a peer before we have committed the block,
	// then we may end this goroutine because stop is closed before we read the
	// committed block from the commitChannel. The result of this would be that
	// we receive a committed block from the previous sealing operation on the
	// commitChannel in the current seal operation. For now we will skip blocks
	// that do not match.
	go func() {
		for {
			select {
			case <-stop:
				return
			case committedBlock := <-b.commitChannel:
				// Check that we are committing the block we were asked to seal.
				if committedBlock.Hash() != block.Hash() {
					continue
				}
				results <- committedBlock
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
	}

	b.currentBlockAwaiter.setValue(block)
	return nil
}

// readyForMessages returns true if the bridge is in a state to handle a
// message from a peer.
func (b *bridge) readyForMessages() bool {
	return atomic.LoadInt32(&b.broadcasterSet) == 1 && atomic.LoadInt32(&b.started) == 1
}

// Methods for consensus.Handler: This interface was introduced by the istanbul
// BFT fork, so we don't need to keep it to maintain some level of parity
// between Autonity and go-ethereum.

// Protocol implements consensus.Handler.Protocol
func (b *bridge) Protocol() (protocolName string, extraMsgCodes uint64) {
	return "tendermint", 2 //nolint
}

// HandleMsg implements consensus.Handler.HandleMsg, this returns a bool to
// indicate whether the message was handled, if we return false then the
// message will be passed on by the caller to be handled by the default eth
// handler.
func (b *bridge) HandleMsg(addr common.Address, msg p2p.Msg) (bool, error) {
	switch msg.Code {
	case tendermintMsg:
		if !b.readyForMessages() {
			return true, nil
		}
		var data []byte
		if err := msg.Decode(&data); err != nil {
			return true, fmt.Errorf("failed to decode tendermint message: %v", err)
		}
		go b.eventMux.Post(events.MessageEvent{
			Payload: data,
		})
	case tendermintSyncMsg:
		if !b.readyForMessages() {
			return true, nil
		}
		go b.eventMux.Post(events.SyncEvent{Addr: addr})
	default:
		return false, nil
	}

	return true, nil
}

// SetBroadcaster implements consensus.Handler.SetBroadcaster
func (b *bridge) SetBroadcaster(broadcaster consensus.Broadcaster) {
	atomic.StoreInt32(&b.broadcasterSet, 1)
	b.blockBroadcaster = broadcaster
}

// NewChainHead implements consensus.Handler.NewChainHead
func (b *bridge) NewChainHead() error {
	go b.eventMux.Post(events.CommitEvent{})
	return nil
}

func (b *bridge) Commit(proposal *algorithm.ConsensusMessage) error {
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

	// If we are the proposer, send the block to the  commit channel
	if b.address == b.committee.GetProposer(proposal.Round).Address {
		b.commitChannel <- block
	} else {
		b.blockBroadcaster.Enqueue("tendermint", block)
	}

	b.logger.Info("committed a block", "hash", block.Hash())
	return nil
}

func (b *bridge) createCommittee(block *types.Block) committee {
	var committeeSet committee
	var err error
	var lastProposer common.Address
	header := block.Header()
	switch b.proposerPolicy {
	case config.RoundRobin:
		if !header.IsGenesis() {
			lastProposer, err = types.Ecrecover(header)
			if err != nil {
				panic(fmt.Sprintf("unable to recover proposer address from header %q: %v", header, err))
			}
		}
		committeeSet, err = newRoundRobinSet(header.Committee, lastProposer)
		if err != nil {
			panic(fmt.Sprintf("failed to construct committee %v", err))
		}
	case config.WeightedRandomSampling:
		committeeSet = newWeightedRandomSamplingCommittee(block, b.autonityContract, b.statedb)
	default:
		panic(fmt.Sprintf("unrecognised proposer policy %q", b.proposerPolicy))
	}
	return committeeSet
}

var errStopped = errors.New("stopped")

// Start implements core.Tendermint.Start
func (b *bridge) Start(ctx context.Context, contract *autonity.Contract, blockchain *core.BlockChain) {
	atomic.StoreInt32(&b.started, 1)
	//println("starting")
	// Set the autonity contract and blockchain
	b.autonityContract = contract
	b.blockchain = blockchain
	ctx, b.cancel = context.WithCancel(ctx)

	// Subscribe
	b.eventsSub = b.eventMux.Subscribe(events.MessageEvent{}, &algorithm.Timeout{}, events.CommitEvent{})
	b.syncEventSub = b.eventMux.Subscribe(events.SyncEvent{})

	b.wg = &sync.WaitGroup{}

	// Tendermint Finite State Machine discrete event loop
	b.wg.Add(1)
	go b.mainEventLoop(ctx)
}

func (b *bridge) Stop() {
	atomic.StoreInt32(&b.started, 0)
	//println(addr(c.address), c.height, "stopping")

	b.logger.Info("stopping tendermint.core", "addr", addr(b.address))

	b.cancel()

	// stop the block awaiter if it is waiting
	b.currentBlockAwaiter.stop()

	// Unsubscribe
	b.eventsSub.Unsubscribe()
	b.syncEventSub.Unsubscribe()

	//println(addr(c.address), c.height, "almost stopped")
	// Ensure all event handling go routines exit
	b.wg.Wait()
}

func (b *bridge) newHeight(prevBlock *types.Block) error {
	b.syncTimer = time.NewTimer(20 * time.Second)
	b.lastHeader = prevBlock.Header()
	b.height = new(big.Int).SetUint64(prevBlock.NumberU64() + 1)
	b.committee = b.createCommittee(prevBlock)

	// Create new oracle and algorithm
	b.algo = algorithm.New(algorithm.NodeID(b.address), newOracle(b.lastHeader, b.msgStore, b.committee, b.currentBlockAwaiter))

	// Handle messages for the new height
	msg, timeout, err := b.algo.StartRound(0)
	if err != nil {
		return err
	}

	// Note that we don't risk entering an infinite loop here since
	// start round can only return results with broadcasts or timeouts.
	err = b.handleResult(nil, msg, timeout)
	if err != nil {
		return err
	}
	for _, msg := range b.msgStore.heightMessages(b.height.Uint64()) {
		err := b.handleCurrentHeightMessage(msg)
		b.logger.Error("failed to handle current height message", "message", msg.String(), "err", err)
	}
	return nil
}

func (b *bridge) handleResult(rc *algorithm.RoundChange, cm *algorithm.ConsensusMessage, to *algorithm.Timeout) error {

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
			err := b.Commit(rc.Decision)
			if err != nil {
				panic(fmt.Sprintf("%s Failed to commit sr.Decision: %s err: %v", algorithm.NodeID(b.address).String(), spew.Sdump(rc.Decision), err))
			}
		} else {
			cm, to, err := b.algo.StartRound(rc.Round) // nolint
			if err != nil {
				return err
			}
			// Note that we don't risk entering an infinite loop here since
			// start round can only return results with broadcasts or timeouts.
			err = b.handleResult(nil, cm, to)
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
		msg, err := encodeSignedMessage(cm, b.key, b.msgStore)
		if err != nil {
			panic(fmt.Sprintf(
				"%s We were unable to build a message, this indicates a programming error: %v",
				addr(b.address),
				err,
			))
		}

		// Broadcast in a new goroutine
		go func() {
			// send to self
			messageEvent := events.MessageEvent{
				Payload: msg,
			}
			b.eventMux.Post(messageEvent)
			// Broadcast to peers
			b.broadcaster.Broadcast(msg)
		}()

	case to != nil:
		time.AfterFunc(time.Duration(to.Delay)*time.Second, func() {
			b.eventMux.Post(to)
		})

	}
	return nil
}

func (b *bridge) mainEventLoop(ctx context.Context) {
	defer b.wg.Done()

	lastBlockMined, err := b.latestBlockRetreiver.RetrieveLatestBlock()
	if err != nil {
		panic(err)
	}
	err = b.newHeight(lastBlockMined)
	if err != nil {
		//println(addr(c.address), c.height.Uint64(), "exiting main event loop", "err", err)
		return
	}

	// Ask for sync when the engine starts
	b.syncer.AskSync(b.lastHeader)

eventLoop:
	for {
		select {
		case <-b.syncTimer.C:
			b.syncer.AskSync(b.lastHeader)
			b.syncTimer = time.NewTimer(20 * time.Second)

		case ev, ok := <-b.syncEventSub.Chan():
			if !ok {
				break eventLoop
			}
			syncEvent := ev.Data.(events.SyncEvent)
			b.logger.Info("Processing sync message", "from", syncEvent.Addr)
			b.syncer.SyncPeer(syncEvent.Addr, b.msgStore.rawHeightMessages(b.height.Uint64()))
		case ev, ok := <-b.eventsSub.Chan():
			if !ok {
				break eventLoop
			}
			// A real ev arrived, process interesting content
			switch e := ev.Data.(type) {
			case events.MessageEvent:
				//println("got a message")
				/*
					Basic validity checks
				*/

				m, err := decodeSignedMessage(e.Payload)
				if err != nil {
					fmt.Printf("some error: %v\n", err)
					continue
				}
				// Check we haven't already processed this message
				if b.msgStore.Message(m.hash) != nil {
					// Message was already processed
					continue
				}
				err = b.msgStore.addMessage(m, e.Payload)
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
				if m.consensusMessage.Height != b.height.Uint64() {
					// Nothing to do here
					continue
				}

				err = b.handleCurrentHeightMessage(m)
				if err == errStopped {
					return
				}
				if err != nil {
					b.logger.Debug("core.mainEventLoop problem processing message", "err", err)
					continue
				}
				b.broadcaster.Broadcast(e.Payload)
			case *algorithm.Timeout:
				var cm *algorithm.ConsensusMessage
				var rc *algorithm.RoundChange
				switch e.TimeoutType {
				case algorithm.Propose:
					//println(addr(c.address), "on timeout propose", e.Height, "round", e.Round)
					cm = b.algo.OnTimeoutPropose(e.Height, e.Round)
				case algorithm.Prevote:
					//println(addr(c.address), "on timeout prevote", e.Height, "round", e.Round)
					cm = b.algo.OnTimeoutPrevote(e.Height, e.Round)
				case algorithm.Precommit:
					//println(addr(c.address), "on timeout precommit", e.Height, "round", e.Round)
					rc = b.algo.OnTimeoutPrecommit(e.Height, e.Round)
				}
				// if cm != nil {
				// 	println("nonnil timeout")
				// }
				err := b.handleResult(rc, cm, nil)
				if err != nil {
					//println(addr(c.address), c.height.Uint64(), "exiting main event loop", "err", err)
					return
				}
			case events.CommitEvent:
				println(addr(b.address), "commit event")
				b.logger.Debug("Received a final committed proposal")

				lastBlock, err := b.latestBlockRetreiver.RetrieveLatestBlock()
				if err != nil {
					panic(err)
				}
				err = b.newHeight(lastBlock)
				if err != nil {
					//println(addr(c.address), c.height.Uint64(), "exiting main event loop", "err", err)
					return
				}
			}

		case <-ctx.Done():
			b.logger.Info("mainEventLoop is stopped", "event", ctx.Err())
			break eventLoop
		}
	}

}

func (b *bridge) handleCurrentHeightMessage(m *message) error {
	//println(addr(c.address), c.height.String(), m.String(), "received")
	cm := m.consensusMessage
	/*
		Domain specific validity checks, now we know that we are at the same
		height as this message we can rely on lastHeader.
	*/

	// Check that the message came from a committee member, if not we ignore it.
	if b.lastHeader.CommitteeMember(m.address) == nil {
		// TODO turn this into an error type that can be checked for at a
		// higher level to close the connection to this peer.
		return fmt.Errorf("received message from non committee member: %v", m)
	}

	switch cm.MsgType {
	case algorithm.Propose:
		// We ignore proposals from non proposers
		if b.committee.GetProposer(cm.Round).Address != m.address {
			b.logger.Warn("Ignore proposal messages from non-proposer")
			return errNotFromProposer

			// TODO verify proposal here.
			//
			// If we are introducing time into the mix then what we are saying
			// is that we don't expect different participants' clocks to drift
			// out of sync more than some delta. And if they do then we don't
			// expect consensus to work.
			//
			// So in the case that clocks drift too far out of sync and say a
			// node considers a proposal invalid that 2f+1 other nodes
			// precommit for that node becomes stuck and can only continue in
			// consensus by re-syncing the blocks.
			//
			// So in verifying the proposal wrt time we should verify once
			// within reasonable clock sync bounds and then set the validity
			// based on that and never re-process the message again.

		}
		// Proposals values are allowed to be invalid.
		if _, err := b.verifier.VerifyProposal(*b.msgStore.value(common.Hash(cm.Value)), b.blockchain, b.address.String()); err == nil {
			//println(addr(c.address), "valid", cm.Value.String())
			b.msgStore.setValid(common.Hash(cm.Value))
		}
	default:
		// All other messages that have reached this point are valid, but we
		// are not marking the value valid here, we are marking the message
		// valid.
		b.msgStore.setValid(m.hash)
	}

	rc, cm, to := b.algo.ReceiveMessage(cm)
	err := b.handleResult(rc, cm, to)
	if err != nil {
		return err
	}
	return nil
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

type Broadcaster struct {
	address common.Address
	pmm     peerMessageMap
	peers   consensus.Peers
}

func NewBroadcaster(address common.Address, peers consensus.Peers) *Broadcaster {
	return &Broadcaster{
		address: address,
		peers:   peers,
		pmm:     &degeneratePeerMessageMap{},
	}
}

// Broadcast implements tendermint.Backend.Broadcast
func (b *Broadcaster) Broadcast(payload []byte) {
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

type Syncer struct {
	peers consensus.Peers
}

func NewSyncer(peers consensus.Peers) *Syncer {
	return &Syncer{
		peers: peers,
	}
}

func (s *Syncer) AskSync(header *types.Header) {
	var count uint64
	for _, p := range s.peers.Peers() {
		//ask to a quorum nodes to sync, 1 must then be honest and updated
		if count >= bft.Quorum(header.TotalVotingPower()) {
			break
		}
		go p.Send(tendermintSyncMsg, []byte{}) //nolint

		member := header.CommitteeMember(p.Address())
		if member == nil {
			continue
		}
		count += member.VotingPower.Uint64()
	}
}

// Synchronize new connected peer with current height state
func (s *Syncer) SyncPeer(address common.Address, messages [][]byte) {
	for _, p := range s.peers.Peers() {
		if address == p.Address() {
			for _, msg := range messages {
				//We do not save sync messages in the arc cache as recipient could not have been able to process some previous sent.
				go p.Send(tendermintMsg, msg) //nolint
			}
			break
		}
	}
}

type LatestBlockRetriever struct {
	db      ethdb.Database
	statedb state.Database
}

func NewLatestBlockRetriever(db ethdb.Database, state state.Database) *LatestBlockRetriever {
	return &LatestBlockRetriever{
		db:      db,
		statedb: state,
		// Here we use the value of 256 which is the
		// eth.DefaultConfig.TrieCleanCache value which is value assigned to
		// cacheConfig.TrieCleanLimit which is what is then used in
		// eth.BlockChain to initialise the state database.
		// statedb: state.NewDatabase(db),
	}
}
func (l *LatestBlockRetriever) RetrieveLatestBlock() (*types.Block, error) {
	hash := rawdb.ReadHeadBlockHash(l.db)
	if hash == (common.Hash{}) {
		return nil, fmt.Errorf("empty database")
	}

	number := rawdb.ReadHeaderNumber(l.db, hash)
	if number == nil {
		return nil, fmt.Errorf("failed to find number for block hash %s", hash.String())
	}

	block := rawdb.ReadBlock(l.db, hash, *number)
	if block == nil {
		return nil, fmt.Errorf("failed to read block content for block number %d with hash %s", *number, hash.String())
	}

	_, err := l.statedb.OpenTrie(block.Root())
	if err != nil {
		return nil, fmt.Errorf("missing state for block number %d with hash %s err: %v", *number, hash.String(), err)
	}
	return block, nil
}

func (l *LatestBlockRetriever) RetrieveBlockState(block *types.Block) (*state.StateDB, error) {
	return state.New(block.Root(), l.statedb, nil)
}
