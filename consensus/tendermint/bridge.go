package tendermint

import (
	context "context"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"sync"
	time "time"

	common "github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus"
	"github.com/clearmatics/autonity/consensus/tendermint/algorithm"
	"github.com/clearmatics/autonity/consensus/tendermint/bft"
	"github.com/clearmatics/autonity/consensus/tendermint/config"
	"github.com/clearmatics/autonity/consensus/tendermint/events"
	autonity "github.com/clearmatics/autonity/contracts/autonity"
	core "github.com/clearmatics/autonity/core"
	"github.com/clearmatics/autonity/core/rawdb"
	"github.com/clearmatics/autonity/core/state"
	types "github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/ethdb"
	event "github.com/clearmatics/autonity/event"
	"github.com/clearmatics/autonity/log"
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
func New(backend Backend, config *config.Config, key *ecdsa.PrivateKey, broadcaster *Broadcaster, syncer *Syncer, address common.Address, latestBlockRetreiver *LatestBlockRetriever, statedb state.Database, verifier *Verifier) *bridge {
	logger := log.New("addr", address.String())
	c := &bridge{
		key:                  key,
		proposerPolicy:       config.ProposerPolicy,
		address:              address,
		logger:               logger,
		backend:              backend,
		valueSet:             sync.NewCond(&sync.Mutex{}),
		msgStore:             newMessageStore(),
		broadcaster:          broadcaster,
		syncer:               syncer,
		latestBlockRetreiver: latestBlockRetreiver,
		statedb:              statedb,
		verifier:             verifier,
	}
	o := &oracle{
		c:     c,
		store: c.msgStore,
	}
	c.ora = o
	return c
}

type bridge struct {
	key            *ecdsa.PrivateKey
	proposerPolicy config.ProposerPolicy
	address        common.Address
	logger         log.Logger

	backend Backend
	cancel  context.CancelFunc

	eventsSub    *event.TypeMuxSubscription
	syncEventSub *event.TypeMuxSubscription
	wg           *sync.WaitGroup

	msgStore  *messageStore
	syncTimer *time.Timer

	committee  committee
	lastHeader *types.Header

	autonityContract *autonity.Contract

	height *big.Int
	algo   *algorithm.OneShotTendermint
	ora    *oracle

	valueSet     *sync.Cond
	value        *types.Block
	currentBlock *types.Block

	broadcaster          *Broadcaster
	syncer               *Syncer
	latestBlockRetreiver *LatestBlockRetriever
	statedb              state.Database

	verifier *Verifier

	blockchain *core.BlockChain // TODO need to set this on start
}

func (c *bridge) SetValue(b *types.Block) {
	c.valueSet.L.Lock()
	defer c.valueSet.L.Unlock()
	if c.value == nil {
		c.valueSet.Signal()
	}
	c.value = b
	//println(addr(c.address), c.height, "setting value", c.value.Hash().String()[2:8], "value height", c.value.Number().String())
}

func (c *bridge) AwaitValue(ctx context.Context, height *big.Int) (*types.Block, error) {
	c.valueSet.L.Lock()
	defer c.valueSet.L.Unlock()

	for {
		select {
		case <-ctx.Done():
			return nil, errStopped
		default:
			if c.value == nil || c.value.Number().Cmp(height) != 0 {
				c.value = nil
				// if c.value == nil {
				// 	println(addr(c.address), c.height.String(), "awaiting vlaue", "valueisnil")
				// } else {
				// 	println(addr(c.address), c.height.String(), "awaiting vlaue", "value height", c.value.Number().String(), "awaited height", height.String())
				// }
				c.valueSet.Wait()
			} else {
				v := c.value
				//println(addr(c.address), c.height, "received awaited vlaue", c.value.Hash().String()[2:8], "value height", c.value.Number().String(), "awaited height", height.String())

				// We put the value in the store here since this is called from the main
				// thread of the algorithm, and so we don't end up needing to syncronise
				// the store.  TODO this is a potential memory leak. We are adding a value
				// without it being referenced by a message that is tied to a height, so it
				// may never be cleared.
				c.msgStore.addValue(v.Hash(), v)
				// We assume our own suggestions are valid
				c.msgStore.setValid(v.Hash())
				c.value = nil
				return v, nil
			}
		}
	}
}

func (c *bridge) Commit(proposal *algorithm.ConsensusMessage) (*types.Block, error) {
	committedSeals := c.msgStore.signatures(proposal.Value, proposal.Round, proposal.Height)
	message := c.msgStore.matchingProposal(proposal)
	// Sanity checks
	if message == nil || message.value == nil {
		return nil, fmt.Errorf("attempted to commit nil block")
	}
	if message.proposerSeal == nil {
		return nil, fmt.Errorf("attempted to commit block without proposer seal")
	}
	if proposal.Round < 0 {
		return nil, fmt.Errorf("attempted to commit a block in a negative round: %d", proposal.Round)
	}
	if len(committedSeals) == 0 {
		return nil, fmt.Errorf("attempted to commit block without any committed seals")
	}

	for _, seal := range committedSeals {
		if len(seal) != types.BFTExtraSeal {
			return nil, fmt.Errorf("attempted to commit block with a committed seal of invalid length: %s", hex.EncodeToString(seal))
		}
	}
	// Add the proposer seal coinbase and committed seals into the block.
	h := message.value.Header()
	h.CommittedSeals = committedSeals
	h.ProposerSeal = message.proposerSeal
	h.Coinbase = message.address
	h.Round = uint64(proposal.Round)
	block := message.value.WithSeal(h)
	c.backend.Commit(block, c.committee.GetProposer(proposal.Round).Address)

	c.logger.Info("commit a block", "hash", block.Hash())
	return block, nil
}

func (c *bridge) createCommittee(block *types.Block) committee {
	var committeeSet committee
	var err error
	var lastProposer common.Address
	header := block.Header()
	switch c.proposerPolicy {
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
		committeeSet = newWeightedRandomSamplingCommittee(block, c.autonityContract, c.statedb)
	default:
		panic(fmt.Sprintf("unrecognised proposer policy %q", c.proposerPolicy))
	}
	return committeeSet
}

var errStopped error = errors.New("stopped")

// Start implements core.Tendermint.Start
func (c *bridge) Start(ctx context.Context, contract *autonity.Contract, blockchain *core.BlockChain) {
	//println("starting")
	// Set the autonity contract and blockchain
	c.autonityContract = contract
	c.blockchain = blockchain
	ctx, c.cancel = context.WithCancel(ctx)

	// Subscribe
	c.eventsSub = c.backend.Subscribe(events.MessageEvent{}, &algorithm.Timeout{}, events.CommitEvent{})
	c.syncEventSub = c.backend.Subscribe(events.SyncEvent{})

	c.wg = &sync.WaitGroup{}

	// Tendermint Finite State Machine discrete event loop
	c.wg.Add(1)
	go c.mainEventLoop(ctx)
}

// stop implements core.Engine.stop
func (c *bridge) Stop() {
	//println(addr(c.address), c.height, "stopping")

	c.logger.Info("stopping tendermint.core", "addr", addr(c.address))

	c.cancel()

	// Signal to wake up await value if it is waiting.
	c.valueSet.L.Lock()
	c.valueSet.Signal()
	c.valueSet.L.Unlock()

	// Unsubscribe
	c.eventsSub.Unsubscribe()
	c.syncEventSub.Unsubscribe()

	//println(addr(c.address), c.height, "almost stopped")
	// Ensure all event handling go routines exit
	c.wg.Wait()
}

func (c *bridge) newHeight(ctx context.Context, height uint64) error {
	c.syncTimer = time.NewTimer(20 * time.Second)
	newHeight := new(big.Int).SetUint64(height)
	// set the new height
	c.height = newHeight
	var err error
	c.currentBlock, err = c.AwaitValue(ctx, newHeight)
	if err != nil {
		return err
	}
	prevBlock, err := c.latestBlockRetreiver.RetrieveLatestBlock()
	if err != nil {
		panic(err)
	}

	c.lastHeader = prevBlock.Header()
	committeeSet := c.createCommittee(prevBlock)
	c.committee = committeeSet

	// Update internals of oracle
	c.ora.lastHeader = c.lastHeader
	c.ora.committeeSet = committeeSet

	// Handle messages for the new height
	msg, timeout := c.algo.StartRound(newHeight.Uint64(), 0, algorithm.ValueID(c.currentBlock.Hash()))

	// If we are making a proposal, we need to ensure that we add the proposal
	// block to the msg store, so that it can be picked up in buildMessage.
	if msg != nil {
		//println(addr(c.address), "adding value", height, c.currentBlock.Hash().String())
		c.msgStore.addValue(c.currentBlock.Hash(), c.currentBlock)
	}

	// Note that we don't risk enterning an infinite loop here since
	// start round can only return results with brodcasts or schedules.
	// TODO actually don't return result from Start round.
	err = c.handleResult(ctx, nil, msg, timeout)
	if err != nil {
		return err
	}
	for _, msg := range c.msgStore.heightMessages(newHeight.Uint64()) {
		err := c.handleCurrentHeightMessage(ctx, msg)
		c.logger.Error("failed to handle current height message", "message", msg.String(), "err", err)
	}
	return nil
}

func (c *bridge) handleResult(ctx context.Context, rc *algorithm.RoundChange, cm *algorithm.ConsensusMessage, to *algorithm.Timeout) error {

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

			// This will ultimately lead to a commit event, which we will pick
			// up on but we will ignore it because instead we will wait here to
			// select the next value that matches this height.
			_, err := c.Commit(rc.Decision)
			if err != nil {
				panic(fmt.Sprintf("%s Failed to commit sr.Decision: %s err: %v", algorithm.NodeID(c.address).String(), spew.Sdump(rc.Decision), err))
			}
			err = c.newHeight(ctx, rc.Height)
			if err != nil {
				return err
			}

		} else {
			// sanity check
			currBlockNum := c.currentBlock.Number().Uint64()
			if currBlockNum != rc.Height {
				panic(fmt.Sprintf("current block number %d out of sync with  height %d", currBlockNum, rc.Height))
			}

			cm, to := c.algo.StartRound(rc.Height, rc.Round, algorithm.ValueID(c.currentBlock.Hash())) // nolint
			// Note that we don't risk enterning an infinite loop here since
			// start round can only return results with brodcasts or schedules.
			// TODO actually don't return result from Start round.
			err := c.handleResult(ctx, nil, cm, to)
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
		msg, err := encodeSignedMessage(cm, c.key, c.msgStore)
		if err != nil {
			panic(fmt.Sprintf(
				"%s We were unable to build a message, this indicates a programming error: %v",
				addr(c.address),
				err,
			))
		}

		// Broadcast in a new goroutine
		go func(committee types.Committee) {
			// send to self
			event := events.MessageEvent{
				Payload: msg,
			}
			c.backend.Post(event)
			// Broadcast to peers
			c.broadcaster.Broadcast(ctx, committee, msg)
		}(c.lastHeader.Committee)

	case to != nil:
		time.AfterFunc(time.Duration(to.Delay)*time.Second, func() {
			c.backend.Post(to)
		})

	}
	return nil
}

func (c *bridge) mainEventLoop(ctx context.Context) {
	defer c.wg.Done()
	// Start a new round from last height + 1
	c.algo = algorithm.New(algorithm.NodeID(c.address), c.ora)

	lastBlockMined, err := c.latestBlockRetreiver.RetrieveLatestBlock()
	if err != nil {
		panic(err)
	}
	err = c.newHeight(ctx, lastBlockMined.NumberU64()+1)
	if err != nil {
		//println(addr(c.address), c.height.Uint64(), "exiting main event loop", "err", err)
		return
	}

	// Ask for sync when the engine starts
	c.syncer.AskSync(c.lastHeader)

eventLoop:
	for {
		select {
		case <-c.syncTimer.C:
			c.syncer.AskSync(c.lastHeader)
			c.syncTimer = time.NewTimer(20 * time.Second)

		case ev, ok := <-c.syncEventSub.Chan():
			if !ok {
				break eventLoop
			}
			event := ev.Data.(events.SyncEvent)
			c.logger.Info("Processing sync message", "from", event.Addr)
			c.syncer.SyncPeer(event.Addr, c.msgStore.rawHeightMessages(c.height.Uint64()))
		case ev, ok := <-c.eventsSub.Chan():
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
				if c.msgStore.Message(m.hash) != nil {
					// Message was already processed
					continue
				}
				err = c.msgStore.addMessage(m, e.Payload)
				if err != nil {
					// could be multiple proposal messages from the same proposer
					continue
				}
				if m.consensusMessage.MsgType == algorithm.Propose {
					c.msgStore.addValue(m.value.Hash(), m.value)
				}

				// If this message is for a future height then we cannot validate it
				// because we lack the relevant header, we will process it when we reach
				// that height. If it is for a previous height then we are not intersted in
				// it. But it has been added to the message store in case other peers would
				// like to sync it.
				if m.consensusMessage.Height != c.height.Uint64() {
					// Nothing to do here
					continue
				}

				err = c.handleCurrentHeightMessage(ctx, m)
				if err == errStopped {
					return
				}
				if err != nil {
					c.logger.Debug("core.mainEventLoop problem processing message", "err", err)
					continue
				}
				c.broadcaster.Broadcast(ctx, c.lastHeader.Committee, e.Payload)
			case *algorithm.Timeout:
				var cm *algorithm.ConsensusMessage
				var rc *algorithm.RoundChange
				switch e.TimeoutType {
				case algorithm.Propose:
					//println(addr(c.address), "on timeout propose", e.Height, "round", e.Round)
					cm = c.algo.OnTimeoutPropose(e.Height, e.Round)
				case algorithm.Prevote:
					//println(addr(c.address), "on timeout prevote", e.Height, "round", e.Round)
					cm = c.algo.OnTimeoutPrevote(e.Height, e.Round)
				case algorithm.Precommit:
					//println(addr(c.address), "on timeout precommit", e.Height, "round", e.Round)
					rc = c.algo.OnTimeoutPrecommit(e.Height, e.Round)
				}
				// if cm != nil {
				// 	println("nonnil timeout")
				// }
				err := c.handleResult(ctx, rc, cm, nil)
				if err != nil {
					//println(addr(c.address), c.height.Uint64(), "exiting main event loop", "err", err)
					return
				}
			case events.CommitEvent:
				//println(addr(c.address), "commit event")
				c.logger.Debug("Received a final committed proposal")

				lastBlock, err := c.latestBlockRetreiver.RetrieveLatestBlock()
				if err != nil {
					panic(err)
				}

				height := new(big.Int).Add(lastBlock.Number(), common.Big1)
				if height.Cmp(c.height) == 0 {
					//println(addr(c.address), "Discarding event as core is at the same height", "height", c.height)
					c.logger.Debug("Discarding event as core is at the same height", "height", c.height)
				} else {
					//println(addr(c.address), "Received proposal is ahead", "height", c.height, "block_height", height.String())
					c.logger.Debug("Received proposal is ahead", "height", c.height, "block_height", height)
					err := c.newHeight(ctx, height.Uint64())
					if err != nil {
						//println(addr(c.address), c.height.Uint64(), "exiting main event loop", "err", err)
						return
					}
				}
			}
		case <-ctx.Done():
			c.logger.Info("mainEventLoop is stopped", "event", ctx.Err())
			break eventLoop
		}
	}

}

func (c *bridge) handleCurrentHeightMessage(ctx context.Context, m *message) error {
	//println(addr(c.address), c.height.String(), m.String(), "received")
	cm := m.consensusMessage
	/*
		Domain specific validity checks, now we know that we are at the same
		height as this message we can rely on lastHeader.
	*/

	// Check that the message came from a committee member, if not we ignore it.
	if c.lastHeader.CommitteeMember(m.address) == nil {
		// TODO turn this into an error type that can be checked for at a
		// higher level to close the connection to this peer.
		return fmt.Errorf("received message from non committee member: %v", m)
	}

	switch cm.MsgType {
	case algorithm.Propose:
		// We ignore proposals from non proposers
		if c.committee.GetProposer(cm.Round).Address != m.address {
			c.logger.Warn("Ignore proposal messages from non-proposer")
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
		if _, err := c.verifier.VerifyProposal(*c.msgStore.value(common.Hash(cm.Value)), c.blockchain, c.address.String()); err == nil {
			//println(addr(c.address), "valid", cm.Value.String())
			c.msgStore.setValid(common.Hash(cm.Value))
		}
	default:
		// All other messages that have reached this point are valid, but we
		// are not marking the value valid here, we are marking the message
		// valid.
		c.msgStore.setValid(m.hash)
	}

	rc, cm, to := c.algo.ReceiveMessage(cm)
	err := c.handleResult(ctx, rc, cm, to)
	if err != nil {
		return err
	}
	return nil
}

const (
	tendermintMsg     = 0x11
	tendermintSyncMsg = 0x12
)

type peerMessageMap interface {
	// knowsMessage returns true if the peer knows the current message
	knowsMessage(addr common.Address, hash common.Hash) bool
}

// TODO actually implement thit
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
func (b *Broadcaster) Broadcast(ctx context.Context, committee types.Committee, payload []byte) {
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
