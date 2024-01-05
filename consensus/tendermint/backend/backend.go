package backend

import (
	"crypto/ecdsa"
	"errors"
	"sync"
	"time"

	"github.com/autonity/autonity/consensus/tendermint/core/constants"
	"github.com/autonity/autonity/consensus/tendermint/core/interfaces"
	"github.com/autonity/autonity/consensus/tendermint/core/message"

	lru "github.com/hashicorp/golang-lru"
	ring "github.com/zfjagann/golang-ring"

	"github.com/autonity/autonity/accounts/abi"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/consensus/misc"
	tendermintCore "github.com/autonity/autonity/consensus/tendermint/core"
	"github.com/autonity/autonity/consensus/tendermint/events"
	"github.com/autonity/autonity/core"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/core/vm"
	"github.com/autonity/autonity/crypto"
	"github.com/autonity/autonity/event"
	"github.com/autonity/autonity/log"
)

const (
	// fetcherID is the ID indicates the block is from BFT engine
	fetcherID = "tendermint"
	// ring buffer to be able to handle at maximum 10 rounds, 20 committee and 3 messages types
	ringCapacity = 10 * 20 * 3
	// while asking sync for consensus messages, if we do not find any peers we try again after 10 ms
	retryPeriod = 10
)

var (
	// ErrStoppedEngine is returned if the engine is stopped
	ErrStoppedEngine = errors.New("stopped engine")
)

// New creates an Ethereum Backend for BFT core engine.
func New(privateKey *ecdsa.PrivateKey,
	vmConfig *vm.Config,
	services *interfaces.Services,
	evMux *event.TypeMux,
	ms *tendermintCore.MsgStore,
	log log.Logger) *Backend {

	recentMessages, _ := lru.NewARC(inmemoryPeers)
	knownMessages, _ := lru.NewARC(inmemoryMessages)

	backend := &Backend{
		eventMux:       event.NewTypeMuxSilent(evMux, log),
		privateKey:     privateKey,
		address:        crypto.PubkeyToAddress(privateKey.PublicKey),
		logger:         log,
		coreStarted:    false,
		recentMessages: recentMessages,
		knownMessages:  knownMessages,
		vmConfig:       vmConfig,
		MsgStore:       ms,
		jailed:         make(map[common.Address]uint64),
	}

	backend.pendingMessages.SetCapacity(ringCapacity)
	core := tendermintCore.New(backend, services, backend.address, log)

	backend.gossiper = NewGossiper(backend.recentMessages, backend.knownMessages, backend.address, backend.logger, backend.stopped)
	if services != nil {
		backend.gossiper = services.Gossiper(backend)
	}
	backend.core = core
	return backend
}

// ----------------------------------------------------------------------------

type Backend struct {
	eventMux     *event.TypeMuxSilent
	privateKey   *ecdsa.PrivateKey
	address      common.Address
	logger       log.Logger
	blockchain   *core.BlockChain
	currentBlock func() *types.Block
	hasBadBlock  func(hash common.Hash) bool

	// the channels for tendermint engine notifications
	commitCh          chan<- *types.Block
	proposedBlockHash common.Hash
	coreStarted       bool
	core              interfaces.Core
	stopped           chan struct{}
	wg                sync.WaitGroup
	coreMu            sync.RWMutex

	// we save the last received p2p.messages in the ring buffer
	pendingMessages ring.Ring

	// interface to enqueue blocks to fetcher and find peers
	Broadcaster consensus.Broadcaster

	// interface to enqueue blocks to fetcher and find peers
	Enqueuer consensus.Enqueuer
	// interface to gossip consensus messages
	gossiper interfaces.Gossiper

	//ARCCache is patented by IBM but it has expired https://patents.google.com/patent/US7167953B2/en
	recentMessages *lru.ARCCache // the cache of peer's messages
	knownMessages  *lru.ARCCache // the cache of self messages

	contractsMu sync.RWMutex //todo(youssef): is that necessary?
	vmConfig    *vm.Config

	MsgStore   *tendermintCore.MsgStore
	jailed     map[common.Address]uint64
	jailedLock sync.RWMutex
}

func (sb *Backend) BlockChain() *core.BlockChain {
	return sb.blockchain
}

// Address implements tendermint.Backend.Address
func (sb *Backend) Address() common.Address {
	return sb.address
}

// Broadcast implements tendermint.Backend.Broadcast
func (sb *Backend) Broadcast(committee types.Committee, message message.Msg) {
	// send to others
	sb.Gossip(committee, message)
	// send to self
	go sb.Post(events.MessageEvent{
		Message: message,
	})
}

func (sb *Backend) AskSync(header *types.Header) {
	sb.gossiper.AskSync(header)
}

// Gossip implements tendermint.Backend.Gossip
func (sb *Backend) Gossip(committee types.Committee, msg message.Msg) {
	sb.gossiper.Gossip(committee, msg)
}

// Gossip implements tendermint.Backend.Gossip
func (sb *Backend) UpdateStopChannel(stopCh chan struct{}) {
	sb.gossiper.UpdateStopChannel(stopCh)
}

// KnownMsgHash dumps the known messages in case of gossiping.
func (sb *Backend) KnownMsgHash() []common.Hash {
	m := make([]common.Hash, 0, sb.knownMessages.Len())
	for _, v := range sb.knownMessages.Keys() {
		m = append(m, v.(common.Hash))
	}
	return m
}

func (sb *Backend) Logger() log.Logger {
	return sb.logger
}

func (sb *Backend) Gossiper() interfaces.Gossiper {
	return sb.gossiper
}

// Commit implements tendermint.Backend.Commit
func (sb *Backend) Commit(proposal *types.Block, round int64, seals [][]byte) error {
	h := proposal.Header()
	// Append seals and round into extra-data
	if err := types.WriteCommittedSeals(h, seals); err != nil {
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
	sb.eventMux.Post(ev)
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
	// ignore errEmptyCommittedSeals error because we don't have the committed seals yet
	if err == nil || errors.Is(err, types.ErrEmptyCommittedSeals) {
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
		for i, tx := range proposal.Transactions() {
			state.Prepare(tx.Hash(), i)
			// Might be vulnerable to DoS Attack depending on gaslimit
			// Todo : Double check
			receipt, receiptErr := core.ApplyTransaction(sb.blockchain.Config(), sb.blockchain, nil, gp, state, header, tx, usedGas, *sb.vmConfig)
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
func (sb *Backend) Sign(data common.Hash) ([]byte, common.Address) {
	ret, err := crypto.Sign(data[:], sb.privateKey)
	if err != nil {
		// We panic here, it should never happen.
		sb.logger.Crit("Consensus signing failed")
	}
	return ret, sb.address
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
	targets := map[common.Address]struct{}{address: {}}
	ps := sb.Broadcaster.FindPeers(targets)
	p, connected := ps[address]
	if !connected {
		return
	}
	messages := sb.core.CurrentHeightMessages()
	for _, msg := range messages {
		//We do not save sync messages in the arc cache as recipient could not have been able to process some previous sent.
		go p.SendRaw(NetworkCodes[msg.Code()], msg.Payload()) //nolint
	}
}

func (sb *Backend) ResetPeerCache(address common.Address) {
	ms, ok := sb.recentMessages.Get(address)
	var m *lru.ARCCache
	if ok {
		m, _ = ms.(*lru.ARCCache)
		m.Purge()
	}
}

func (sb *Backend) RemoveMessageFromLocalCache(message message.Msg) {
	sb.knownMessages.Remove(message.Hash())
}
