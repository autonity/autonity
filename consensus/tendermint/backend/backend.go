package backend

import (
	"crypto/ecdsa"
	"errors"
	"math"
	"sync"
	"sync/atomic"
	"time"

	ring "github.com/zfjagann/golang-ring"

	"github.com/autonity/autonity/accounts/abi"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/common/fixsizecache"
	"github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/consensus/misc"
	tendermintCore "github.com/autonity/autonity/consensus/tendermint/core"
	"github.com/autonity/autonity/consensus/tendermint/core/constants"
	"github.com/autonity/autonity/consensus/tendermint/core/interfaces"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/consensus/tendermint/events"
	"github.com/autonity/autonity/core"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/core/vm"
	"github.com/autonity/autonity/crypto"
	"github.com/autonity/autonity/crypto/blst"
	"github.com/autonity/autonity/event"
	"github.com/autonity/autonity/log"
)

const (
	// fetcherID is the ID indicates the block is from BFT engine
	fetcherID = "tendermint"
	// ring buffer to be able to handle at maximum 10 rounds, 100 committee and 3 messages types
	ringCapacity = 10 * 100 * 3
	// maximum number of future height messages
	maxFutureMsgs = 10 * 100 * 3
	// while asking sync for consensus messages, if we do not find any peers we try again after 10 ms
	retryPeriod = 10
	// number of buckets to allocate in the fixed cache
	numBuckets = 499
	// max number of entries in each packet
	numEntries = 10
)

var (
	// ErrStoppedEngine is returned if the engine is stopped
	ErrStoppedEngine = errors.New("stopped engine")
)

// New creates an Ethereum Backend for BFT core engine.
func New(nodeKey *ecdsa.PrivateKey,
	consensusKey blst.SecretKey,
	vmConfig *vm.Config,
	services *interfaces.Services,
	evMux *event.TypeMux,
	ms *tendermintCore.MsgStore,
	log log.Logger, noGossip bool) *Backend {

	knownMessages := fixsizecache.New[common.Hash, bool](numBuckets, numEntries, fixsizecache.HashKey[common.Hash])

	backend := &Backend{
		eventMux:        event.NewTypeMuxSilent(evMux, log),
		nodeKey:         nodeKey,
		consensusKey:    consensusKey,
		address:         crypto.PubkeyToAddress(nodeKey.PublicKey),
		logger:          log,
		knownMessages:   knownMessages,
		vmConfig:        vmConfig,
		MsgStore:        ms, //TODO: we use this only in tests, to easily reach the msg store when having a reference to the backend. It would be better to just have the `accountability` module as a part of the backend object.
		messageCh:       make(chan events.UnverifiedMessageEvent, 1000),
		jailed:          make(map[common.Address]uint64),
		future:          make(map[uint64][]*events.UnverifiedMessageEvent),
		futureMinHeight: math.MaxUint64,
	}

	backend.pendingMessages.SetCapacity(ringCapacity)

	backend.gossiper = NewGossiper(backend.knownMessages, backend.address, backend.logger, backend.stopped)
	if services != nil {
		backend.gossiper = services.Gossiper(backend)
	}

	core := tendermintCore.New(backend, services, backend.address, log, noGossip)
	backend.core = core
	backend.evDispatcher = core

	backend.aggregator = newAggregator(backend, core, log)

	return backend
}

// ----------------------------------------------------------------------------

type Backend struct {
	eventMux     *event.TypeMuxSilent
	nodeKey      *ecdsa.PrivateKey
	consensusKey blst.SecretKey
	address      common.Address
	logger       log.Logger
	blockchain   *core.BlockChain
	currentBlock func() *types.Block
	hasBadBlock  func(hash common.Hash) bool

	// the channels for tendermint engine notifications
	commitCh          chan<- *types.Block
	messageCh         chan events.UnverifiedMessageEvent // to send events to the aggregator
	proposedBlockHash common.Hash
	coreStarting      atomic.Bool
	coreRunning       atomic.Bool
	core              interfaces.Core
	evDispatcher      interfaces.EventDispatcher
	stopped           chan struct{}
	wg                sync.WaitGroup

	// used to save consensus messages while core is stopped
	pendingMessages ring.Ring

	// interface to find peers
	Broadcaster consensus.Broadcaster
	// interface to enqueue blocks to fetcher
	Enqueuer consensus.Enqueuer
	// interface to gossip consensus messages
	gossiper interfaces.Gossiper

	knownMessages *fixsizecache.Cache[common.Hash, bool] // the cache of self messages

	contractsMu sync.RWMutex //todo(youssef): is that necessary?
	vmConfig    *vm.Config

	MsgStore   *tendermintCore.MsgStore //TODO: we use this only in tests, to easily reach the msg store when having a reference to the backend. It would be better to just have the `accountability` module as a part of the backend object.
	jailed     map[common.Address]uint64
	jailedLock sync.RWMutex

	aggregator *aggregator

	// buffer for future height events and related metadata
	// TODO(lorenzo) refinements, wrap this stuff into a separate struct?
	future          map[uint64][]*events.UnverifiedMessageEvent // UnverifiedMessageEvent is used slightly inappropriately here, as the future height messages still need to pass the checks in `handleDecodedMsg` before being posted to the aggregator.
	futureMinHeight uint64
	futureMaxHeight uint64
	futureSize      uint64
	futureLock      sync.RWMutex
}

func (sb *Backend) BlockChain() *core.BlockChain {
	return sb.blockchain
}

func (sb *Backend) MessageCh() <-chan events.UnverifiedMessageEvent {
	return sb.messageCh
}

// Address implements tendermint.Backend.Address
func (sb *Backend) Address() common.Address {
	return sb.address
}

// Broadcast implements tendermint.Backend.Broadcast
func (sb *Backend) Broadcast(committee types.Committee, message message.Msg) {
	// send to others
	sb.Gossip(committee, message)
	// send to self (directly to Core and FD, no need to verify local messages)
	go sb.Post(events.MessageEvent{
		Message: message,
		ErrCh:   nil,
	})
}

func (sb *Backend) AskSync(header *types.Header) {
	sb.gossiper.AskSync(header)
}

// Gossip implements tendermint.Backend.Gossip
func (sb *Backend) Gossip(committee types.Committee, msg message.Msg) {
	sb.gossiper.Gossip(committee, msg)
}

// UpdateStopChannel implements tendermint.Backend.Gossip
func (sb *Backend) UpdateStopChannel(stopCh chan struct{}) {
	sb.gossiper.UpdateStopChannel(stopCh)
}

// KnownMsgHash dumps the known messages in case of gossiping.
func (sb *Backend) KnownMsgHash() []common.Hash {
	return sb.knownMessages.Keys()
}

func (sb *Backend) Logger() log.Logger {
	return sb.logger
}

func (sb *Backend) Gossiper() interfaces.Gossiper {
	return sb.gossiper
}

// Commit implements tendermint.Backend.Commit
func (sb *Backend) Commit(proposal *types.Block, round int64, quorumCertificate types.AggregateSignature) error {
	h := proposal.Header()
	// Append quorum certificate and round into extra-data
	if err := types.WriteQuorumCertificate(h, quorumCertificate); err != nil {
		return err
	}
	if err := types.WriteRound(h, round); err != nil {
		return err
	}
	// update block's header
	proposal = proposal.WithSeal(h)
	sb.logger.Info("Quorum of Precommits received", "proposal", proposal.Hash(), "round", round, "height", proposal.Number().Uint64())
	// - if the proposed and committed blocks are the same, send the proposed hash
	//   to resultCh channel, which is being watched inside the worker.ResultLoop() function.
	// - otherwise, we try to insert the block.
	// -- if success, the ChainHeadEvent event will be broadcasted, try to build
	//    the next block and the previous Seal() will be stopped.
	// -- otherwise, a error will be returned and a round change event will be fired.
	if sb.proposedBlockHash == proposal.Hash() && !sb.isResultChanNil() {
		sb.sendResultChan(proposal)
		return nil
	}

	if sb.Enqueuer != nil {
		sb.Enqueuer.Enqueue(fetcherID, proposal)
	}
	return nil
}

func (sb *Backend) Post(ev any) {
	switch ev := ev.(type) {
	case events.CommitEvent:
		sb.evDispatcher.Post(ev)
	case events.NewCandidateBlockEvent:
		sb.evDispatcher.Post(ev)
	case events.UnverifiedMessageEvent:
		sb.messageCh <- ev
	default:
		sb.eventMux.Post(ev)
	}
}

func (sb *Backend) Subscribe(types ...any) *event.TypeMuxSubscription {
	return sb.eventMux.Subscribe(types...)
}

// VerifyProposal implements tendermint.Backend.VerifyProposal and verifiy if the proposal is valid
func (sb *Backend) VerifyProposal(proposal *types.Block) (time.Duration, error) {
	// TODO: fix always false statement and check for non nil
	// TODO: use interface instead of type

	if sb.HasBadProposal(proposal.Hash()) {
		return 0, core.ErrBannedHash
	}

	// verify if the proposal block is already included in the node's local chain.
	// This scenario can happen when we are processing a proposal, but in the meantime other peers already reached quorum on it,
	// therefore we already received the finalized block through p2p block propagation.
	// NOTE: this function execution is not atomic, the block could be not included at the time of this check
	// and become included right after we passed the check.
	if sb.blockchain.HasHeader(proposal.Hash(), proposal.NumberU64()) {
		return 0, constants.ErrAlreadyHaveBlock
	}

	// verify the header of proposed proposal
	err := sb.VerifyHeader(sb.blockchain, proposal.Header(), false)
	// ignore errEmptyQuorumCertificate error because we don't have the quorum certificate yet
	if err == nil || errors.Is(err, types.ErrEmptyQuorumCertificate) {
		var (
			receipts types.Receipts

			usedGas        = new(uint64)
			gp             = new(core.GasPool).AddGas(proposal.GasLimit())
			header         = proposal.Header()
			proposalNumber = header.Number.Uint64()
			parent         = sb.blockchain.GetBlock(proposal.ParentHash(), proposal.NumberU64()-1)
		)

		// Verify London hard fork attributes including min base fee
		if err := misc.VerifyEip1559Header(sb.blockchain.Config(), sb.blockchain, parent.Header(), header); err != nil {
			// Verify the header's EIP-1559 attributes.
			return 0, err
		}
		// We need to process all the transaction to get the latest state to get the latest committee
		state, stateErr := sb.blockchain.StateAt(parent.Root())
		if stateErr != nil {
			return 0, stateErr
		}

		// Validate the body of the proposal
		if err = sb.blockchain.Validator().ValidateBody(proposal); err != nil {
			return 0, err
		}

		// sb.blockchain.Processor().Process() was not called because it calls back Finalize() and would have modified the proposal
		// Instead only the transactions are applied to the copied state
		config := sb.blockchain.Config()
		signer := types.MakeSigner(config, header.Number)
		// Create a new context to be used in the EVM environment
		blockContext := core.NewEVMBlockContext(header, sb.BlockChain(), nil)
		vmenv := vm.NewEVM(blockContext, vm.TxContext{}, state, config, *sb.vmConfig)
		for i, tx := range proposal.Transactions() {
			state.Prepare(tx.Hash(), i)
			// Might be vulnerable to DoS Attack depending on gaslimit
			// Todo : Double check
			receipt, receiptErr := core.ApplyTransactionWithContext(signer, config, sb.blockchain, nil, gp, state, header, tx, usedGas, vmenv)
			if receiptErr != nil {
				return 0, receiptErr
			}
			receipts = append(receipts, receipt)
		}

		state.Prepare(common.ACHash(proposal.Number()), len(proposal.Transactions()))
		committee, receipt, err := sb.Finalize(sb.blockchain, header, state, proposal.Transactions(), nil, receipts)
		if err != nil {
			return 0, err
		}
		receipts = append(receipts, receipt)
		//Validate the state of the proposal
		if err = sb.blockchain.Validator().ValidateState(proposal, state, receipts, *usedGas); err != nil {
			sb.logger.Error("proposal proposed, bad root state", err)
			return 0, err
		}

		//Perform the actual comparison
		if len(header.Committee) != len(committee) {
			sb.logger.Error("wrong committee set",
				"proposalNumber", proposalNumber,
				"extraLen", len(header.Committee),
				"currentLen", len(committee),
				"committee", header.Committee,
				"current", committee,
			)
			return 0, consensus.ErrInconsistentCommitteeSet
		}

		for i := range committee {
			if header.Committee[i].Address != committee[i].Address ||
				header.Committee[i].VotingPower.Cmp(committee[i].VotingPower) != 0 {
				sb.logger.Error("wrong committee member in the set",
					"index", i,
					"currentVerifier", sb.address.String(),
					"proposalNumber", proposalNumber,
					"headerCommittee", header.Committee[i],
					"computedCommittee", committee[i],
					"fullHeader", header.Committee,
					"fullComputed", committee,
				)
				return 0, consensus.ErrInconsistentCommitteeSet
			}
		}
		// At this stage committee field is consistent with the validator list returned by Soma-contract

		return 0, nil
	} else if errors.Is(err, consensus.ErrFutureTimestampBlock) {
		return time.Unix(int64(proposal.Header().Time), 0).Sub(now()), consensus.ErrFutureTimestampBlock
	}

	// Here we are considering this proposal invalid because we pruned the parent's state
	// however this is our local node fault, not the remote proposer fault.
	if errors.Is(err, consensus.ErrPrunedAncestor) {
		sb.logger.Error("Rejecting a proposal because local node has pruned parent's state")
		sb.logger.Error("Please check your pruning settings")
	}
	return 0, err
}

// Sign implements tendermint.Backend.Sign
func (sb *Backend) Sign(data common.Hash) blst.Signature {
	signature := sb.consensusKey.Sign(data[:])
	return signature
}

func (sb *Backend) HeadBlock() *types.Block {
	return sb.currentBlock()
}

func (sb *Backend) HasBadProposal(hash common.Hash) bool {
	if sb.hasBadBlock == nil {
		return false
	}
	return sb.hasBadBlock(hash)
}

func (sb *Backend) GetContractABI() *abi.ABI {
	// after the contract is upgradable, call it from contract object rather than from conf.
	return sb.blockchain.ProtocolContracts().ABI()
}

func (sb *Backend) CoreState() interfaces.CoreState {
	return sb.core.CoreState()
}

// CommitteeEnodes retrieve the list of validators enodes for the current block
func (sb *Backend) CommitteeEnodes() []string {
	db, err := sb.blockchain.State()
	if err != nil {
		sb.logger.Error("Failed to get state", "err", err)
		return nil
	}
	enodes, err := sb.blockchain.ProtocolContracts().CommitteeEnodes(sb.blockchain.CurrentBlock(), db, false)
	if err != nil {
		sb.logger.Error("Failed to get block committee", "err", err)
		return nil
	}
	return enodes.StrList
}

// SyncPeer Synchronize new connected peer with current height messages
func (sb *Backend) SyncPeer(address common.Address) {
	if sb.Broadcaster == nil {
		return
	}
	sb.logger.Debug("Syncing", "peer", address)
	peer, ok := sb.Broadcaster.FindPeer(address)
	if !ok {
		return
	}
	messages := sb.core.CurrentHeightMessages()
	sb.logger.Debug("sent current height messages", "peer", address, "n", len(messages), "msgs", messages)
	for _, msg := range messages {
		//We do not save sync messages in the arc cache as recipient could not have been able to process some previous sent.
		go peer.SendRaw(NetworkCodes[msg.Code()], msg.Payload()) //nolint
	}
}

// called by tendermint core to dump core state
func (sb *Backend) FutureMsgs() []message.Msg {
	sb.futureLock.RLock()
	defer sb.futureLock.RUnlock()

	var msgs []message.Msg
	for _, evs := range sb.future {
		for _, ev := range evs {
			msgs = append(msgs, ev.Message)
		}
	}

	return msgs
}
