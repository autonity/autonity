package accountability

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math"
	"math/big"
	"sort"
	"sync"
	"time"

	"github.com/autonity/autonity/accounts/abi/bind"
	"github.com/autonity/autonity/autonity"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/consensus/tendermint/bft"
	engineCore "github.com/autonity/autonity/consensus/tendermint/core"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/consensus/tendermint/events"
	"github.com/autonity/autonity/core"
	"github.com/autonity/autonity/core/state"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/event"
	"github.com/autonity/autonity/internal/ethapi"
	"github.com/autonity/autonity/log"
	"github.com/autonity/autonity/rlp"
)

type ChainContext interface {
	consensus.ChainReader
	CurrentBlock() *types.Block
	SubscribeChainEvent(ch chan<- core.ChainEvent) event.Subscription
	State() (*state.StateDB, error)
	ProtocolContracts() *autonity.ProtocolContracts
	StateAt(root common.Hash) (*state.StateDB, error)
	HasBadBlock(hash common.Hash) bool
	Validator() core.Validator
}

const (
	msgGCInterval                 = 60      // every 60 blocks to GC msg store.
	offChainAccusationProofWindow = 10      // the time window in block for one to provide off chain innocence proof before it is escalated on chain.
	maxNumOfInnocenceProofCached  = 120 * 4 // 120 blocks with 4 on each height that rule engine can produce totally over a height.
	maxAccusationRatePerHeight    = 4       // max number of accusation can be produced by rule engine over a height against to a validator.
	maxFutureHeightMsgs           = 1000    // max num of msg buffer for the future heights.
)

var (
	errDuplicatedMsg   = errors.New("duplicated msg")
	errEquivocation    = errors.New("equivocation")
	errFutureMsg       = errors.New("future height msg")
	errNotCommitteeMsg = errors.New("msg from none committee member")
	errProposer        = errors.New("proposal is not from proposer")
	errInvalidMessage  = errors.New("invalid consensus message")

	errNoEvidenceForPO  = errors.New("no proof of innocence found for rule PO")
	errNoEvidenceForPVN = errors.New("no proof of innocence found for rule PVN")
	errNoEvidenceForPVO = errors.New("no proof of innocence found for rule PVO")
	errNoEvidenceForC1  = errors.New("no proof of innocence found for rule C1")
	errUnprovableRule   = errors.New("unprovable rule")

	nilValue = common.Hash{}
)

// FaultDetector it subscribe chain event to trigger rule engine to apply patterns over
// msg store, it sends Proof of challenge if it detects any potential misbehavior, either it
// read state db on each new height to get the latest challenges from autonity contract's view,
// and to prove its innocent if there were any challenges on the suspicious node.
type FaultDetector struct {
	innocenceProofBuff *InnocenceProofBuffer
	protocolContracts  *autonity.ProtocolContracts
	rateLimiter        *AccusationRateLimiter

	wg               sync.WaitGroup
	tendermintMsgSub *event.TypeMuxSubscription

	txPool      *core.TxPool
	ethBackend  ethapi.Backend
	transactOps *bind.TransactOpts

	eventReporterCh chan *autonity.AccountabilityEvent
	// chain event subscriber for rule engine.
	ruleEngineBlockCh  chan core.ChainEvent
	ruleEngineBlockSub event.Subscription

	// on-chain accountability event
	accountabilityEventCh  chan *autonity.AccountabilityNewAccusation
	accountabilityEventSub event.Subscription

	blockchain ChainContext
	address    common.Address
	msgStore   *engineCore.MsgStore

	chainEventCh  chan core.ChainEvent
	chainEventSub event.Subscription

	misbehaviourProofCh chan *autonity.AccountabilityEvent
	futureMessages      map[uint64][]message.Msg        // map[blockHeight][]*tendermintMessages
	futureMessageCount  uint64                          // a counter to count the total cached future height msg.
	pendingEvents       []*autonity.AccountabilityEvent // accountability event buffer.

	offChainAccusationsMu sync.RWMutex
	offChainAccusations   []*Proof // off chain accusations list, ordered in chain height from low to high.
	broadcaster           consensus.Broadcaster

	logger log.Logger
}

// NewFaultDetector call by ethereum object to create fd instance.
func NewFaultDetector(
	chain ChainContext,
	nodeAddress common.Address,
	sub *event.TypeMuxSubscription,
	ms *engineCore.MsgStore,
	txPool *core.TxPool,
	ethBackend ethapi.Backend,
	nodeKey *ecdsa.PrivateKey,
	protocolContracts *autonity.ProtocolContracts,
	logger log.Logger) *FaultDetector {

	transactOps, err := bind.NewKeyedTransactorWithChainID(nodeKey, chain.Config().ChainID)
	if err != nil {
		logger.Crit("Critical error building transactor", "err", err)
	}
	transactOps.GasTipCap = common.Big0
	fd := &FaultDetector{
		innocenceProofBuff:    NewInnocenceProofBuffer(),
		protocolContracts:     protocolContracts,
		rateLimiter:           NewAccusationRateLimiter(),
		txPool:                txPool,
		ethBackend:            ethBackend,
		transactOps:           transactOps,
		tendermintMsgSub:      sub,
		ruleEngineBlockCh:     make(chan core.ChainEvent, 300),
		accountabilityEventCh: make(chan *autonity.AccountabilityNewAccusation),
		blockchain:            chain,
		address:               nodeAddress,
		msgStore:              ms,
		chainEventCh:          make(chan core.ChainEvent, 300),
		eventReporterCh:       make(chan *autonity.AccountabilityEvent, 10),
		misbehaviourProofCh:   make(chan *autonity.AccountabilityEvent, 100),
		futureMessages:        make(map[uint64][]message.Msg),
		futureMessageCount:    0,
		logger:                logger, // Todo(youssef): remove context
	}
	// todo(youssef): analyze chainEvent vs chainHeadEvent and very important: what to do during sync !
	fd.ruleEngineBlockSub = fd.blockchain.SubscribeChainEvent(fd.ruleEngineBlockCh)
	fd.chainEventSub = fd.blockchain.SubscribeChainEvent(fd.chainEventCh)

	fd.accountabilityEventSub, _ = protocolContracts.WatchNewAccusation(
		nil,
		fd.accountabilityEventCh,
		[]common.Address{nodeAddress},
	)
	return fd
}

// Start listen for new block events from blockchain, do the tasks like take challenge and provide Proof for innocent, the
// Fault Detector rule engine could also trigger from here to scan those msgs of msg store by applying rules.
func (fd *FaultDetector) Start() {
	fd.wg.Add(1)
	go fd.eventReporter()
	go fd.ruleEngine()
	go fd.consensusMsgHandlerLoop()
}

func (fd *FaultDetector) isMsgExpired(headHeight uint64, msgHeight uint64) bool {
	return headHeight > consensus.AccountabilityHeightRange && msgHeight < headHeight-consensus.AccountabilityHeightRange
}

func (fd *FaultDetector) SetBroadcaster(broadcaster consensus.Broadcaster) {
	fd.broadcaster = broadcaster
}

func (fd *FaultDetector) saveFutureHeightMsg(m message.Msg) {
	fd.futureMessages[m.H()] = append(fd.futureMessages[m.H()], m)
	fd.futureMessageCount++

	// buffer is full, remove the furthest away msg from buffer to prevent DoS attack.
	if fd.futureMessageCount >= maxFutureHeightMsgs {
		maxHeight := m.H()
		for h, msgs := range fd.futureMessages {
			if h > maxHeight && len(msgs) > 0 {
				maxHeight = h
			}
		}
		if len(fd.futureMessages[maxHeight]) > 1 {
			fd.futureMessages[maxHeight] = fd.futureMessages[maxHeight][:len(fd.futureMessages[maxHeight])-1]
		} else {
			delete(fd.futureMessages, maxHeight)
		}
		fd.futureMessageCount--
	}
}

func (fd *FaultDetector) deleteFutureHeightMsg(height uint64) {
	length := len(fd.futureMessages[height])
	fd.futureMessageCount = fd.futureMessageCount - uint64(length)
	delete(fd.futureMessages, height)
}

func preCheckMessage(m message.Msg, chain ChainContext) error {
	lastHeader := chain.GetHeaderByNumber(m.H() - 1)
	if lastHeader == nil {
		return errFutureMsg
	}
	return m.Validate(lastHeader.CommitteeMember)
}

func (fd *FaultDetector) consensusMsgHandlerLoop() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
tendermintMsgLoop:
	for {
		curHeight := fd.blockchain.CurrentBlock().Number().Uint64()
		select {
		case ev, ok := <-fd.tendermintMsgSub.Chan():
			if !ok {
				break tendermintMsgLoop
			}
			// handle consensus message or innocence proof messages
			switch e := ev.Data.(type) {
			case events.MessageEvent:
				if fd.isMsgExpired(curHeight, e.Message.H()) {
					fd.logger.Debug("Fault detector: discarding old message", "sender", e.Message.Sender())
					continue tendermintMsgLoop
				}
				if err := fd.processMsg(e.Message); err != nil && !errors.Is(err, errFutureMsg) {
					fd.logger.Warn("Detected faulty message", "return", err)
					continue tendermintMsgLoop
				}
			case events.AccountabilityEvent:
				err := fd.handleOffChainAccountabilityEvent(e.Payload, e.Sender)
				if err != nil {
					fd.logger.Info("Accountability: Dropping peer", "peer", e.Sender)
					// the errors return from handler could freeze the peer connection for 30 seconds by according to dev p2p protocol.
					select {
					case e.ErrCh <- err:
					default: // do nothing
					}
					continue tendermintMsgLoop
				}
			}
		case e, ok := <-fd.chainEventCh:
			if !ok {
				break tendermintMsgLoop
			}

			// on every 60 blocks, reset Peer Justified Accusations and height accusations counters.
			if e.Block.NumberU64()%msgGCInterval == 0 {
				fd.rateLimiter.resetHeightRateLimiter()
				fd.rateLimiter.resetPeerJustifiedAccusations()
			}
			/* THIS HAS BEEN DELETED TODO VERIFY
			height := e.block.NumberU64()
			if fd.isMsgExpired(curHeight, height) {
				fd.logger.Info("fault detector: discarding old height messages", "height", height)
				fd.deleteFutureHeightMsg(height)
				continue tendermintMsgLoop
			}
			*/

			for h, messages := range fd.futureMessages {
				if h <= curHeight {
					for _, m := range messages {
						if err := fd.processMsg(m); err != nil {
							fd.logger.Error("Fault detector: error while processing consensus msg", "err", err)
						}
					}
					// once messages are processed, delete it from buffer.
					fd.deleteFutureHeightMsg(h)
				}
			}
		case <-ticker.C:
			// on each 1 seconds, reset the rate limiter counters.
			fd.rateLimiter.resetRateLimiter()
		case err, ok := <-fd.chainEventSub.Err():
			if ok {
				// why crit? what can happen here?
				fd.logger.Crit("block subscription error", "err", err)
			}
			break tendermintMsgLoop
		}
	}
	close(fd.misbehaviourProofCh)
}

// check to GC msg store for those msgs out of buffering window on every 60 blocks.
// todo(youssef): this might tbe unsufficient and lead to a DDOS OOM attack
func (fd *FaultDetector) checkMsgStoreGC(height uint64) {
	if height > consensus.AccountabilityHeightRange && height%msgGCInterval == 0 {
		threshold := height - consensus.AccountabilityHeightRange
		fd.msgStore.DeleteOlds(threshold)
	}
}

func (fd *FaultDetector) ruleEngine() {
loop:
	for {
		select {
		// chain accusationEvent update, provide proof of innocent if one is on challenge, rule engine scanning is triggered also.
		case ev, ok := <-fd.ruleEngineBlockCh:
			if !ok {
				break loop
			}

			// try to escalate expired off chain accusation on chain.
			fd.escalateExpiredAccusations(ev.Block.NumberU64())

			// run rule engine over a specific height.
			if ev.Block.NumberU64() > uint64(consensus.DeltaBlocks) {
				checkpoint := ev.Block.NumberU64() - uint64(consensus.DeltaBlocks)
				if events := fd.runRuleEngine(checkpoint); len(events) > 0 {
					fd.pendingEvents = append(fd.pendingEvents, events...)
				}
				if len(fd.pendingEvents) != 0 && fd.canReport(checkpoint) {
					fd.pendingEvents = fd.reportEvents(fd.pendingEvents)
				}
			}
			// msg store delete msgs out of buffering window on every 60 blocks.
			fd.checkMsgStoreGC(ev.Block.NumberU64())
		case accusation := <-fd.accountabilityEventCh:
			fd.logger.Warn("Local node byzantine accusation!")
			accusationEvent, err := fd.protocolContracts.Events(nil, accusation.Id)
			if err != nil {
				// this should not happen
				fd.logger.Error("Can't retrieve accountability event", "id", accusation.Id.Uint64())
				break
			}
			decodedProof, err := decodeRawProof(accusationEvent.RawProof)
			if err != nil {
				fd.logger.Error("Can't decode accusation", "err", err)
				break
			}
			// The signatures must be valid at this stage, however we have to recover the original
			// senders, hence the following call.
			if err := verifyProofSignatures(fd.blockchain, decodedProof); err != nil {
				fd.logger.Error("Can't verify proof signatures", "err", err)
				break
			}
			innocenceProof, err := fd.innocenceProof(decodedProof)
			if err == nil && innocenceProof != nil {
				// send on chain innocence proof ASAP since the client is on challenge that requires the proof to be
				// provided before the client get slashed.
				fd.logger.Warn("Innocence proof found! reporting...")
				fd.eventReporterCh <- innocenceProof
			} else {
				fd.logger.Warn("************************** SLASHING EVENT **************************")
				fd.logger.Warn("Your local node has been accused of malicious behavior")
				fd.logger.Warn("A proof of innocence has not been found: the local node is at high risk of slashing")
				fd.logger.Warn("Reach out to Autonity social media channels for more informations")
				fd.logger.Warn("********************************************************************")
				if err != nil {
					fd.logger.Error("Could not handle accusation", "error", err)
				}
			}

		case m, ok := <-fd.misbehaviourProofCh:
			if !ok {
				break loop
			}
			fd.pendingEvents = append(fd.pendingEvents, m)
		case err, ok := <-fd.ruleEngineBlockSub.Err():
			if ok {
				// youssef: how can that happen?
				fd.logger.Crit("block subscription error", err.Error())
			}
			break loop
		}
	}
}

// canReport assign the validator a dedicated time-window to submit the accountability event
// todo(youssef): this needs to be thoroughly verified accounting for edge cases scenarios at
// the epoch limit. Also the contract side enforcement is missing.
func (fd *FaultDetector) canReport(height uint64) bool {

	committee := fd.blockchain.GetHeaderByNumber(height - 1).Committee

	// each reporting slot contains ReportingSlotPeriod block period that a unique and deterministic validator is asked to
	// be the reporter of that slot period, then at the end block of that slot, the reporter reports
	// available events. Thus, between each reporting slot, we have 5 block period to wait for
	// accountability events to be mined by network, and it is also disaster friendly that if the last
	// reporter fails, the next reporter will continue to report missing events.
	reporterIndex := (height / consensus.ReportingSlotPeriod) % uint64(len(committee))

	// if validator is the reporter of the slot period, and if checkpoint block is the end block of the
	// slot, then it is time to report the collected events by this validator.
	if height%consensus.ReportingSlotPeriod != 0 {
		return false
	}
	// todo(youssef): this seems like a non-committee member can't send a proof/ do we want that?
	return committee[reporterIndex].Address == fd.address
}

func (fd *FaultDetector) Stop() {
	fd.ruleEngineBlockSub.Unsubscribe()
	fd.chainEventSub.Unsubscribe()
	fd.tendermintMsgSub.Unsubscribe()
	fd.accountabilityEventSub.Unsubscribe()
	close(fd.eventReporterCh)
	fd.wg.Wait()
}

// convert the raw proofs into on-chain Proof which contains raw bytes of messages.
func (fd *FaultDetector) eventFromProof(p *Proof) *autonity.AccountabilityEvent {
	var ev = &autonity.AccountabilityEvent{
		EventType:      uint8(p.Type),
		Rule:           uint8(p.Rule),
		Reporter:       fd.address,
		Offender:       p.Message.Sender(),
		Block:          new(big.Int).SetUint64(p.Message.H()), // assigned contract-side
		ReportingBlock: common.Big0,                           // assigned contract-side
		Epoch:          common.Big0,                           // assigned contract-side
		MessageHash:    common.Big0,                           // assigned contract-side
	}
	// panic because encoding must not fail here
	rProof, err := rlp.EncodeToBytes(p)
	if err != nil {
		fd.logger.Crit("error encoding proof", err)
	}
	ev.RawProof = rProof
	return ev
}

// getInnocentProof is called by client who is on a challenge with a certain accusation, to get innocent proof from msg
// store.
func (fd *FaultDetector) innocenceProof(p *Proof) (*autonity.AccountabilityEvent, error) {
	// the protocol contains below provable accusations.
	switch p.Rule {
	case autonity.PO:
		return fd.innocenceProofPO(p)
	case autonity.PVN:
		return fd.innocenceProofPVN(p)
	case autonity.PVO:
		return fd.innocenceProofPVO(p)
	case autonity.C1:
		return fd.innocenceProofC1(p)
	default:
		return nil, errUnprovableRule
	}
}

// get innocent proof of accusation of rule C1 from msg store.
func (fd *FaultDetector) innocenceProofC1(c *Proof) (*autonity.AccountabilityEvent, error) {
	preCommit := c.Message
	height := preCommit.H()

	lastHeader := fd.blockchain.GetHeaderByNumber(height - 1)
	if lastHeader == nil {
		return nil, errNoParentHeader
	}
	quorum := bft.Quorum(lastHeader.TotalVotingPower())

	prevotesForV := fd.msgStore.Get(height, func(m message.Msg) bool {
		return m.Code() == message.PrevoteCode && m.Value() == preCommit.Value() && m.R() == preCommit.R()
	})

	overQuorumVotes := engineCore.OverQuorumVotes(prevotesForV, quorum)
	if overQuorumVotes == nil {
		return nil, errNoEvidenceForC1
	}

	p := fd.eventFromProof(&Proof{
		Type:      autonity.Innocence,
		Rule:      c.Rule,
		Message:   preCommit,
		Evidences: overQuorumVotes,
	})
	return p, nil
}

// get innocent proof of accusation of rule PO from msg store.
func (fd *FaultDetector) innocenceProofPO(c *Proof) (*autonity.AccountabilityEvent, error) {
	// PO: node propose an old value with an validRound, innocent onChainProof of it should be:
	// there are quorum num of prevote for that value at the validRound.
	liteProposal := c.Message
	height := liteProposal.H()
	validRound := liteProposal.(*message.LightProposal).ValidRound()
	lastHeader := fd.blockchain.GetHeaderByNumber(height - 1)
	if lastHeader == nil {
		return nil, errNoParentHeader
	}
	quorum := bft.Quorum(lastHeader.TotalVotingPower())

	prevotes := fd.msgStore.Get(height, func(m message.Msg) bool {
		return m.Code() == message.PrevoteCode && m.R() == validRound && m.Value() == liteProposal.Value()
	})

	overQuorumPreVotes := engineCore.OverQuorumVotes(prevotes, quorum)
	if overQuorumPreVotes == nil {
		// cannot onChainProof its innocent for PO, the on-chain contract will fine it latter once the
		// time window for onChainProof ends.
		return nil, errNoEvidenceForPO
	}
	p := fd.eventFromProof(&Proof{
		Type:      autonity.Innocence,
		Rule:      c.Rule,
		Message:   liteProposal,
		Evidences: overQuorumPreVotes,
	})
	return p, nil
}

// get innocent proof of accusation of rule PVN from msg store.
func (fd *FaultDetector) innocenceProofPVN(c *Proof) (*autonity.AccountabilityEvent, error) {
	// get innocent proofs for PVN, for a prevote that vote for a new value,
	// then there must be a proposal for this new value.
	prevote := c.Message
	height := prevote.H()
	// the only proof of innocence of PVN accusation is that there exist a corresponding proposal
	proposals := fd.msgStore.Get(height, func(m message.Msg) bool {
		return m.Code() == message.ProposalCode &&
			m.R() == prevote.R() &&
			m.Value() == prevote.Value() &&
			m.(*message.Propose).ValidRound() == -1
	})

	if len(proposals) != 0 {
		p := fd.eventFromProof(&Proof{
			Type:    autonity.Innocence,
			Rule:    c.Rule,
			Message: prevote,
			Evidences: []message.Msg{
				message.NewLightProposal(proposals[0].(*message.Propose)),
			},
		})
		return p, nil
	}
	return nil, errNoEvidenceForPVN
}

// get innocent proof of accusation of rule PVO from msg store, it collects quorum preVotes for the value voted at a valid round.
func (fd *FaultDetector) innocenceProofPVO(c *Proof) (*autonity.AccountabilityEvent, error) {
	// get innocent proofs for PVO, collect quorum preVotes at the valid round of the old proposal.
	oldProposal := c.Evidences[0]
	height := oldProposal.H()
	validRound := oldProposal.(*message.LightProposal).ValidRound()
	lastHeader := fd.blockchain.GetHeaderByNumber(height - 1)
	if lastHeader == nil {
		return nil, errNoParentHeader
	}
	quorum := bft.Quorum(lastHeader.TotalVotingPower())

	preVotes := fd.msgStore.Get(height, func(m message.Msg) bool {
		return m.Code() == message.PrevoteCode && m.Value() == oldProposal.Value() && m.R() == validRound
	})

	overQuorumVotes := engineCore.OverQuorumVotes(preVotes, quorum)

	if overQuorumVotes == nil {
		return nil, errNoEvidenceForPVO
	}

	p := fd.eventFromProof(&Proof{
		Type:      autonity.Innocence,
		Rule:      c.Rule,
		Message:   c.Message,
		Evidences: append(c.Evidences, overQuorumVotes...),
	})
	return p, nil
}

// processMsg, check and submit any auto-incriminating, equivocation challenges, and then only store checked msg in msg store.
func (fd *FaultDetector) processMsg(m message.Msg) error {
	// check if msg is from valid committee member
	if err := preCheckMessage(m, fd.blockchain); err != nil {
		if errors.Is(err, errFutureMsg) {
			fd.saveFutureHeightMsg(m)
		}
		return err
	}
	switch msg := m.(type) {
	case *message.Propose:
		if err := fd.checkSelfIncriminatingProposal(msg); err != nil {
			return err
		}
	case *message.Prevote, *message.Precommit:
		if err := fd.checkSelfIncriminatingVote(m); err != nil {
			return err
		}
	default:
		return errInvalidMessage
	}

	// msg pass the auto-incriminating checker, save it in msg store.
	fd.msgStore.Save(m)
	return nil
}

// run rule engine over the specific height of consensus msgs, return the accountable events in proofs.
func (fd *FaultDetector) runRuleEngine(height uint64) []*autonity.AccountabilityEvent {

	// To avoid none necessary accusations, we wait for delta blocks to start rule scan.
	// always skip the heights before first buffered height after the node start up, since it will rise lots of none
	// sense accusations due to the missing of messages during the startup phase, it cost un-necessary payments
	// for the committee member.
	if height <= fd.msgStore.FirstHeightBuffered() {
		return nil
	}
	lastHeader := fd.blockchain.GetHeaderByNumber(height - 1)
	if lastHeader == nil {
		// youssef: is that even possible?
		return nil
	}
	quorum := bft.Quorum(lastHeader.TotalVotingPower())
	proofs := fd.runRulesOverHeight(height, quorum)
	events := make([]*autonity.AccountabilityEvent, 0, len(proofs))

	for _, proof := range proofs {
		if proof.Message.Sender() == fd.address {
			// skip those misbehaviour or accusation against oneself.
			continue
		}
		// process accusation off chain first.
		if proof.Type == autonity.Accusation {
			// push task to accusation processing list, and send it to remote peer before it is escalated on chain.
			fd.addOffChainAccusation(proof)
			fd.sendOffChainAccusationMsg(proof)
			continue
		}

		p := fd.eventFromProof(proof)
		events = append(events, p)
	}

	return events
}

func (fd *FaultDetector) runRulesOverHeight(height uint64, quorum *big.Int) (proofs []*Proof) {
	// Rules read right to left (find  the right and look for the left)
	//
	// Rules should be evaluated such that we check all possible instances and if we can't find a single instance that
	// passes then we consider the rule failed.
	//
	// There are 2 types of provable misbehaviour.
	// 1. Conflicting messages from a single participant
	// 2. A message that conflicts with a quorum of prevotes.
	// (precommit for differing value in same round as the prevotes or proposal for an
	// old value where in each prior round we can see a quorum of precommits for a distinct value.)

	// We should be here at time t = timestamp(h+1) + delta

	proofs = append(proofs, fd.newProposalsAccountabilityCheck(height)...)
	proofs = append(proofs, fd.oldProposalsAccountabilityCheck(height, quorum)...)
	proofs = append(proofs, fd.prevotesAccountabilityCheck(height, quorum)...)
	proofs = append(proofs, fd.precommitsAccountabilityCheck(height, quorum)...)
	return proofs
}

func (fd *FaultDetector) newProposalsAccountabilityCheck(height uint64) (proofs []*Proof) {
	// ------------New Proposal------------
	// PN:  (Mr′<r,P C|pi)∗ <--- (Mr,P|pi)
	// PN1: [nil ∨ ⊥] <--- [V]
	//
	// Since the message pattern for PN includes only messages sent by pi, we cannot raise an accusation. We can only
	// raise a misbehaviour. To raise a misbehaviour for PN1 we need to have received all the precommits from pi for all
	// r' < r. If any of the precommits is for a non-nil value then we have proof of misbehaviour.

	proposalsNew := fd.msgStore.Get(height, func(m message.Msg) bool {
		return m.Code() == message.ProposalCode && m.(*message.Propose).ValidRound() == -1
	})

	for _, p := range proposalsNew {
		proposal := p

		// Skip if proposal is equivocated
		proposalsForR := fd.msgStore.Get(height, func(m message.Msg) bool {
			return m.Sender() == proposal.Sender() && m.Code() == message.ProposalCode && m.R() == proposal.R()

		})
		// Due to the for loop there must be at least one proposal
		if len(proposalsForR) > 1 {
			continue
		}

		//check all precommits for previous rounds from this sender are nil
		precommits := fd.msgStore.Get(height, func(m message.Msg) bool {
			return m.Sender() == proposal.Sender() && m.Code() == message.PrecommitCode && m.R() < proposal.R() && m.Value() != nilValue
		})
		if len(precommits) != 0 {
			proof := &Proof{
				Type:      autonity.Misbehaviour,
				Rule:      autonity.PN,
				Evidences: precommits[0:1],
				Message:   message.NewLightProposal(proposal.(*message.Propose)),
			}
			proofs = append(proofs, proof)
			fd.logger.Info("Misbehaviour detected", "rule", "PN", "incriminated", proposal.Sender())
		}
	}
	return proofs
}

func (fd *FaultDetector) oldProposalsAccountabilityCheck(height uint64, quorum *big.Int) (proofs []*Proof) {
	// ------------Old Proposal------------
	// PO: (Mr′<r,PV) ∧ (Mr′,PC|pi) ∧ (Mr′<r′′<r,P C|pi)∗ <--- (Mr,P|pi)
	// PO1: [#(Mr′,PV|V) ≥ 2f+ 1] ∧ [nil ∨ V ∨ ⊥] ∧ [nil ∨ ⊥] <--- [V]

	proposalsOld := fd.msgStore.Get(height, func(m message.Msg) bool {
		return m.Code() == message.ProposalCode && m.(*message.Propose).ValidRound() > -1
	})

oldProposalLoop:
	for _, p := range proposalsOld {
		proposal := p
		// Check that in the valid round we see a quorum of prevotes and that there is no precommit at all or a
		// precommit for v or nil.

		// Skip if proposal is equivocated
		proposalsForR := fd.msgStore.Get(height, func(m message.Msg) bool {
			return m.Sender() == proposal.Sender() && m.Code() == message.ProposalCode && m.R() == proposal.R()

		})
		// Due to the for loop there must be at least one proposal
		if len(proposalsForR) > 1 {
			continue oldProposalLoop
		}

		validRound := proposal.(*message.Propose).ValidRound()

		// Is there a precommit for a value other than nil or the proposed value by the current proposer in the valid
		// round? If there is, the proposer has proposed a value for which it is not locked on, thus a Proof of
		// misbehaviour can be generated.
		precommitsFromPiInVR := fd.msgStore.Get(height, func(m message.Msg) bool {
			return m.Code() == message.PrecommitCode && m.R() == validRound && m.Sender() == proposal.Sender() &&
				m.Value() != nilValue && m.Value() != proposal.Value()
		})
		if len(precommitsFromPiInVR) > 0 {
			proof := &Proof{
				Type:      autonity.Misbehaviour,
				Rule:      autonity.PO,
				Evidences: precommitsFromPiInVR[0:1],
				Message:   message.NewLightProposal(proposal.(*message.Propose)),
			}
			proofs = append(proofs, proof)
			fd.logger.Info("Misbehaviour detected", "rule", "PO", "incriminated", proposal.Sender())
			continue oldProposalLoop
		}

		// Is there a precommit for anything other than nil from the proposer between the valid round and the round of
		// the proposal? If there is then that implies the proposer saw 2f+1 prevotes in that round and hence it should
		// have set that round as the valid round.
		precommitsFromPiAfterVR := fd.msgStore.Get(height, func(m message.Msg) bool {
			return m.Code() == message.PrecommitCode && m.R() > validRound && m.R() < proposal.R() &&
				m.Sender() == proposal.Sender() && m.Value() != nilValue
		})
		if len(precommitsFromPiAfterVR) > 0 {
			proof := &Proof{
				Type:      autonity.Misbehaviour,
				Rule:      autonity.PO,
				Evidences: precommitsFromPiAfterVR[0:1],
				Message:   message.NewLightProposal(proposal.(*message.Propose)),
			}
			proofs = append(proofs, proof)
			fd.logger.Info("Misbehaviour detected", "rule", "PO", "incriminated", proposal.Sender())
			continue oldProposalLoop
		}

		// Do we see a quorum for a value other than the proposed value? If so, we have proof of misbehaviour.
		allPrevotesForValidRound := fd.msgStore.Get(height, func(m message.Msg) bool {
			return m.Code() == message.PrevoteCode && m.R() == validRound && m.Value() != proposal.Value()
		})

		prevotesMap := make(map[common.Hash][]message.Msg)
		for _, p := range allPrevotesForValidRound {
			prevotesMap[p.Value()] = append(prevotesMap[p.Value()], p)
		}

		for _, preVotes := range prevotesMap {
			// Here the assumption is that in a single round it is not possible to have 2 value which quorum votes,
			// this would imply at least quorum nodes are malicious which is much higher than our assumption.
			overQuorumVotes := engineCore.OverQuorumVotes(preVotes, quorum)
			if overQuorumVotes != nil {
				proof := &Proof{
					Type:      autonity.Misbehaviour,
					Rule:      autonity.PO,
					Evidences: overQuorumVotes,
					Message:   message.NewLightProposal(proposal.(*message.Propose)),
				}
				proofs = append(proofs, proof)
				fd.logger.Info("Misbehaviour detected", "rule", "PO", "incriminated", proposal.Sender())
				continue oldProposalLoop
			}
		}

		// Do we see a quorum of prevotes in the valid round, if not we can raise an accusation, since we cannot be sure
		// that these prevotes don't exist
		prevotes := fd.msgStore.Get(height, func(m message.Msg) bool {
			// since equivocation msgs are stored, we have to query those preVotes which has same value as the proposal.
			return m.Code() == message.PrevoteCode && m.R() == validRound && m.Value() == proposal.Value()
		})

		if engineCore.OverQuorumVotes(prevotes, quorum) == nil {
			accusation := &Proof{
				Type:    autonity.Accusation,
				Rule:    autonity.PO,
				Message: message.NewLightProposal(proposal.(*message.Propose)),
			}
			proofs = append(proofs, accusation)
			fd.logger.Info("🕵️ Suspicious behavior detected", "rule", "PO", "suspect", proposal.Sender())
		}
	}
	return proofs
}

func (fd *FaultDetector) prevotesAccountabilityCheck(height uint64, quorum *big.Int) (proofs []*Proof) {
	// ------------New and Old prevotes------------

	prevotes := fd.msgStore.Get(height, func(m message.Msg) bool {
		return m.Code() == message.PrevoteCode && m.Value() != nilValue
	})

prevotesLoop:
	for _, p := range prevotes {
		prevote := p

		// Skip if prevote is equivocated
		prevotesForR := fd.msgStore.Get(height, func(m message.Msg) bool {
			return m.Sender() == prevote.Sender() && m.Code() == message.PrevoteCode && m.R() == prevote.R()

		})
		// Due to the for loop there must be at least one preVote.
		if len(prevotesForR) > 1 {
			continue prevotesLoop
		}

		// We need to check whether we have proposals from the prevote's round
		correspondingProposals := fd.msgStore.Get(height, func(m message.Msg) bool {
			return m.Code() == message.ProposalCode && m.R() == prevote.R() && m.Value() == prevote.Value()
		})

		if len(correspondingProposals) == 0 {

			// if there are over quorum prevotes for this corresponding proposal's value, then it indicates current
			// peer just did not receive it. So we can skip the rising of such accusation.
			preVts := fd.msgStore.Get(height, func(m message.Msg) bool {
				return m.Code() == message.PrevoteCode && m.R() == prevote.R() && m.Value() == prevote.Value()
			})

			if engineCore.OverQuorumVotes(preVts, quorum) == nil {
				// The rule for this accusation could be PVO as well since we don't have the corresponding proposal, but
				// it does not mean it's incorrect. More over that, since over quorum prevotes at the round
				// of correspondingProposals are used as the innocence proof, rather than the correspondingProposals, thus
				// we don't worry that the correspondingProposals sender could lie on the proof providing phase.
				accusation := &Proof{
					Type:    autonity.Accusation,
					Rule:    autonity.PVN,
					Message: prevote,
				}
				proofs = append(proofs, accusation)
				continue prevotesLoop
			}
		}

		// We need to ensure that we keep all proposals in the message store, so that we have the maximum chance of
		// finding justification for prevotes. This is to account for equivocation where the proposer send 2 proposals
		// with the same value but different valid rounds to different nodes. We can't penalise the sender of prevote
		// since we can't tell which proposal they received. We just want to find a set of message which fit the rule.
		// Therefore, we need to check all the proposals to find a single one which shows the current prevote is
		// valid.
		var prevotesProofs []*Proof
		for _, cp := range correspondingProposals {
			proposal := cp.(*message.Propose)
			var proof *Proof
			if proposal.ValidRound() == -1 {
				proof = fd.newPrevotesAccountabilityCheck(height, prevote, proposal)
			} else {
				proof = fd.oldPrevotesAccountabilityCheck(height, quorum, proposal, prevote)
			}
			if proof != nil {
				prevotesProofs = append(prevotesProofs, proof)
			}
		}

		if len(prevotesProofs) > 0 {
			for _, proof := range prevotesProofs {
				// If there is any corresponding proposal for which no proof was returned then we know the current prevote
				// is valid.
				if proof == nil {
					continue prevotesLoop
				}
			}

			// There are no corresponding proposal for which the current prevote is valid. We prioritise misbehaviours over
			// accusation since they can be easily proved.
			for _, proof := range prevotesProofs {
				if proof.Type == autonity.Misbehaviour {
					proofs = append(proofs, proof)
					continue prevotesLoop
				}
			}

			// There were no misbehaviours for the current prevote, therefore, pick the first accusation
			proofs = append(proofs, prevotesProofs[0])
		}
	}
	return proofs
}

func (fd *FaultDetector) newPrevotesAccountabilityCheck(height uint64, prevote message.Msg, correspondingProposal *message.Propose) (proof *Proof) {
	// New Proposal, apply PVN rules

	// PVN: (Mr′<r,PC|pi)∧(Mr′<r′′<r,PC|pi)* ∧ (Mr,P|proposer(r)) <--- (Mr,PV|pi)

	// PVN2: [nil ∨ ⊥] ∧ [nil ∨ ⊥] ∧ [V:Valid(V)] <--- [V]: r′= 0,∀r′′< r:Mr′′,PC|pi=nil

	// PVN2, If there is a valid proposal V at round r, and pi never ever precommit(locked a value) before, then pi
	// should prevote for V or a nil in case of timeout at this round.

	// PVN3: [V] ∧ [nil ∨ ⊥] ∧ [V:Valid(V)] <--- [V]:∀r′< r′′<r,Mr′′,PC|pi=nil

	// There is no scope to raise an accusation for these rules since the only message in PVN that is not sent by pi is
	// the proposal, and you require the proposal before you can even attempt to apply the rule.

	// Since we cannot raise an accusation we can only create a proof of misbehaviour. To create a proof of misbehaviour
	// we need to have all the messages in the message pattern, otherwise, we cannot make any statement about the
	// message. We may not have enough information, and we don't want to accuse someone unnecessarily. To show a proof of
	// misbehaviour for PVN2 and PVN3 we need to collect all the precommits from pi and set the latest precommit round
	// as r' and we need to have all the precommit messages from r' to r for pi to be able to check for misbehaviour. If
	// the latest precommit is not for V, and we have all the precommits from r' to r which are nil, then we have proof
	// of misbehaviour.
	precommitsFromPi := fd.msgStore.Get(height, func(m message.Msg) bool {
		return m.Code() == message.PrecommitCode && prevote.Sender() == m.Sender() && m.R() < prevote.R()
	})

	// Check for missing messages. If there are gaps those missing message could be the one that proves pi acted
	// correctly however since we don't have information and enough time has passed we are just going to ignore and move
	// to the next prevote.
	if len(precommitsFromPi) > 0 {
		sort.SliceStable(precommitsFromPi, func(i, j int) bool {
			return precommitsFromPi[i].R() < precommitsFromPi[j].R()
		})
		r := prevote.R()
		rPrime := precommitsFromPi[len(precommitsFromPi)-1].R()
		// Check if the difference between the previous round and current round is more than 1 then exit and return nil
		for i := len(precommitsFromPi) - 1; i >= 0 && (r-rPrime) <= 1; i-- {
			if precommitsFromPi[i].Value() != nilValue {
				pc := precommitsFromPi[i]
				precommitsAtRPrime := fd.msgStore.Get(height, func(m message.Msg) bool {
					return m.Code() == message.PrecommitCode && pc.Sender() == m.Sender() && m.R() == pc.R()
				})

				// Check for equivocation, it is possible there are multiple precommit from pi for the same round.
				// If there are equivocated messages: do nothing. Since pi has already being punished for equivocation
				// round when the equivocated message was first received.
				if len(precommitsAtRPrime) == 1 {
					if precommitsAtRPrime[0].Value() != prevote.Value() {
						fd.logger.Info("Misbehaviour detected", "rule", "PVN", "incriminated", prevote.Sender())
						proof := &Proof{
							Type:    autonity.Misbehaviour,
							Rule:    autonity.PVN,
							Message: prevote,
						}
						// to guarantee this prevote is for a new proposal that is the PVN rule account for, otherwise in
						// prevote for an old proposal, it is valid for one to prevote it if lockedRound <= vr, thus the
						// round gump is valid. This prevents from rising a PVN misbehavior proof from a malicious fault
						// detector by using prevote for an old proposal to challenge an honest slow validator.
						proof.Evidences = append(proof.Evidences, message.NewLightProposal(correspondingProposal))
						proof.Evidences = append(proof.Evidences, precommitsFromPi[i:]...)
						return proof
					}
				}
			}
			if i > 0 {
				r = rPrime
				rPrime = precommitsFromPi[i-1].R()
			}
		}
	}
	return nil
}

func (fd *FaultDetector) oldPrevotesAccountabilityCheck(height uint64, quorum *big.Int, correspondingProposal *message.Propose, prevote message.Msg) (proof *Proof) {
	currentR := correspondingProposal.R()
	validRound := correspondingProposal.ValidRound()

	// If one prevotes for an invalid old proposal, then it should be a misbehaviour.
	if validRound >= currentR {
		fd.logger.Info("Misbehaviour detected", "rule", "PVO3", "incriminated", prevote.Sender())
		proof := &Proof{
			Type:    autonity.Misbehaviour,
			Rule:    autonity.PVO3,
			Message: prevote,
		}
		proof.Evidences = append(proof.Evidences, message.NewLightProposal(correspondingProposal))
		return proof
	}

	// If there is a prevote for an old proposal then pi can only vote for v or send nil (see line 28 and 29 of
	// tendermint pseudocode), therefore if in the valid round there is a quorum for a value other than v, we know pi
	// prevoted incorrectly. If the proposal was a bad proposal, then pi should not have voted for it, thus we do not
	// need to make sure whether the proposal is correct or not (which we would in the proposal checking rules, however,
	// a bad proposal will still exist in our message store, and it shouldn't have an impact on the checking of prevotes).

	allPrevotesForValidRound := fd.msgStore.Get(height, func(m message.Msg) bool {
		return m.Code() == message.PrevoteCode && m.R() == validRound && m.Value() != correspondingProposal.Value()
	})

	prevotesMap := make(map[common.Hash][]message.Msg)
	for _, p := range allPrevotesForValidRound {
		prevotesMap[p.Value()] = append(prevotesMap[p.Value()], p)
	}

	for _, preVotes := range prevotesMap {
		// Here the assumption is that in a single round it is not possible to have 2 value which quorum votes,
		// this would imply at least quorum nodes are malicious which is much higher than our assumption.
		overQuorumVotes := engineCore.OverQuorumVotes(preVotes, quorum)
		if overQuorumVotes != nil {
			fd.logger.Info("Misbehaviour detected", "rule", "PV0", "incriminated", prevote.Sender())
			proof := &Proof{
				Type:    autonity.Misbehaviour,
				Rule:    autonity.PVO,
				Message: prevote,
			}
			proof.Evidences = append(proof.Evidences, message.NewLightProposal(correspondingProposal))
			proof.Evidences = append(proof.Evidences, overQuorumVotes...)
			return proof
		}
	}

	prevotesForVFromValidRound := fd.msgStore.Get(height, func(m message.Msg) bool {
		return m.Code() == message.PrevoteCode && m.R() == validRound && m.Value() == correspondingProposal.Value()
	})

	overQuorumPrevotesForVFromValidRound := engineCore.OverQuorumVotes(prevotesForVFromValidRound, quorum)

	if overQuorumPrevotesForVFromValidRound != nil {
		// PVO: (Mr′′′<r,PV) ∧ (Mr′′′≤r′<r,PC|pi) ∧ (Mr′<r′′<r,PC|pi)∗ ∧ (Mr, P|proposer(r)) ⇐= (Mr,PV|pi)
		// PVO1: [#(V)≥2f+ 1] ∧ [V] ∧ [V ∨ nil ∨ ⊥] ∧ [ V: validRound(V) = r′′′] ⇐= [V]

		// if V is the proposed value at round r and pi did already precommit on V at round r′< r (it locked on it)
		// and did not precommit for other values in any round between r′and r then in round r either pi prevotes
		// for V or nil (in case of a timeout), Moreover, we expect to find 2f+ 1 prevotes for V issued at round
		// r′′′=validRound(V). Notice that, we can have other rounds in which there are 2f+ 1 prevotes for V, but it
		// must be the case at least for this round (as required by line 28).  Indeed, if pi precommitted for V a
		// round r′ != r′′′ then also at round r′we must have 2f+ 1 prevotes for V(will be checked by the precommit
		// rule C1). It follows that there is not relationship between the round r′′′ and r′,which must be set to
		// the last round (if multiple ones) in which pi precommitted for V.

		// Please note pi doesn't need to have precommite for V in valid round, since it could have timed out.
		// Rather we need to find the latest round for which pi committed for V and ensure any rounds after that pi
		// only precommitted for nil

		// PVO’:(Mr′<r, PV) ∧ (Mr′<r′′<r, PC|pi)∗ ∧ (Mr,P|proposer(r)) ⇐= (Mr,P V|pi)
		// PVO2: [#(V)≥2f+ 1] ∧ [V ∨ nil ∨⊥] ∧ [V:validRound(V) =r′] ⇐= [V];
		// if V is the proposed value at round r with validRound(V) =r′ then there must be 2f+ 1 prevotes
		// for V issued at round r′. If moreover, pi did not precommit for other values in any round between
		// r′and r(thus it can be either locked on some values or not) then in round r pi prevotes for V.

		// PVO1 and PVO2 can be merged together. We just need to fetch all precommits between (validRound, currentR)
		// check that we have no gaps and raise a misbehaviour if the last one is not for V.

		precommitsFromPi := fd.msgStore.Get(height, func(m message.Msg) bool {
			return m.Code() == message.PrecommitCode && m.R() > validRound && m.R() < currentR && m.Sender() == prevote.Sender()
		})

		if len(precommitsFromPi) > 0 {
			// sort by round ascending
			sort.SliceStable(precommitsFromPi, func(i, j int) bool {
				return precommitsFromPi[i].R() < precommitsFromPi[j].R()
			})

			// ensure there are no gaps
			if precommitsFromPi[0].R() != validRound+1 || precommitsFromPi[len(precommitsFromPi)-1].R() != currentR-1 {
				return nil
			}
			for i := 1; i < len(precommitsFromPi); i++ {
				prev, cur := precommitsFromPi[i-1].R(), precommitsFromPi[i].R()
				diff := math.Abs(float64(cur) - float64(prev))
				if diff > 1 {
					// at least one round's precommit is missing
					return nil
				}
			}

			// If the last precommit for notV is after the last one for V, raise misbehaviour
			// If all precommits are nil, do not raise misbehaviour. It is a valid correct scenario.
			lastRoundForV := int64(-1)
			lastRoundForNotV := int64(-1)
			for _, pc := range precommitsFromPi {
				if pc.Value() == prevote.Value() && pc.R() > lastRoundForV {
					lastRoundForV = pc.R()
				}

				if pc.Value() != prevote.Value() && pc.Value() != nilValue && pc.R() > lastRoundForNotV {
					lastRoundForNotV = pc.R()
				}
			}

			if lastRoundForNotV > lastRoundForV {
				fd.logger.Info("Misbehaviour detected", "rule", "PVO12", "incriminated", prevote.Sender())
				proof := &Proof{
					Type:    autonity.Misbehaviour,
					Rule:    autonity.PVO12,
					Message: prevote,
				}
				proof.Evidences = append(proof.Evidences, message.NewLightProposal(correspondingProposal))
				proof.Evidences = append(proof.Evidences, precommitsFromPi...)
				return proof
			}
		}
	}

	// if there is no misbehaviour of the prevote msg addressed, then we lastly check accusation.
	if overQuorumPrevotesForVFromValidRound == nil {
		// raise an accusation
		fd.logger.Info("🕵️ Suspicious behavior detected", "rule", "PVO", "suspect", prevote.Sender())
		return &Proof{
			Type:      autonity.Accusation,
			Rule:      autonity.PVO,
			Message:   prevote,
			Evidences: []message.Msg{message.NewLightProposal(correspondingProposal)},
		}
	}

	return nil
}

func (fd *FaultDetector) precommitsAccountabilityCheck(height uint64, quorum *big.Int) (proofs []*Proof) {
	// ------------precommits------------
	// C: [Mr,P|proposer(r)] ∧ [Mr,PV] <--- [Mr,PC|pi]
	// C1: [V:Valid(V)] ∧ [#(V) ≥ 2f+ 1] <--- [V]

	precommits := fd.msgStore.Get(height, func(m message.Msg) bool {
		return m.Code() == message.PrecommitCode && m.Value() != nilValue
	})

precommitLoop:
	for _, preC := range precommits {
		precommit := preC

		// Skip if preCommit is equivocated
		precommitsForR := fd.msgStore.Get(height, func(m message.Msg) bool {
			return m.Sender() == precommit.Sender() && m.Code() == message.PrecommitCode && m.R() == precommit.R()

		})
		// Due to the for loop there must be at least one preCommit.
		if len(precommitsForR) > 1 {
			continue precommitLoop
		}

		// Do we see a quorum for a value other than the proposed value? If so, we have proof of misbehaviour.
		allPrevotesForR := fd.msgStore.Get(height, func(m message.Msg) bool {
			return m.Code() == message.PrevoteCode && m.R() == precommit.R() && m.Value() != precommit.Value()
		})

		prevotesMap := make(map[common.Hash][]message.Msg)
		for _, p := range allPrevotesForR {
			prevotesMap[p.Value()] = append(prevotesMap[p.Value()], p)
		}

		for _, preVotes := range prevotesMap {
			// Here the assumption is that in a single round it is not possible to have 2 value which quorum votes,
			// this would imply at least quorum nodes are malicious which is much higher than our assumption.
			overQuorumVotes := engineCore.OverQuorumVotes(preVotes, quorum)
			if overQuorumVotes != nil {
				proof := &Proof{
					Type:      autonity.Misbehaviour,
					Rule:      autonity.C,
					Evidences: overQuorumVotes,
					Message:   precommit,
				}
				proofs = append(proofs, proof)
				fd.logger.Info("Misbehaviour detected", "rule", "C", "incriminated", precommit.Sender())
				continue precommitLoop
			}
		}

		// Do we see a quorum of prevotes in the same round, if not we can raise an accusation, since we cannot be sure
		// that these prevotes don't exist, this block also covers the Accusation of C since if over quorum prevotes for
		// V indicates that the corresponding proposal of V do exist, thus we don't need to rise accusation for the missing
		// proposal since over 2/3 member should all ready received it
		prevotes := fd.msgStore.Get(height, func(m message.Msg) bool {
			// since equivocation msgs are stored, we have to query those preVotes which has same value as the proposal.
			return m.Code() == message.PrevoteCode && m.R() == precommit.R() && m.Value() == precommit.Value()
		})

		if engineCore.OverQuorumVotes(prevotes, quorum) == nil {
			// We don't have a quorum of prevotes for this precommit to be justified
			accusation := &Proof{
				Type:    autonity.Accusation,
				Rule:    autonity.C1,
				Message: precommit,
			}
			proofs = append(proofs, accusation)

			fd.logger.Info("🕵️ Suspicious behavior detected", "rule", "C1", "suspect", precommit.Sender())
		}
	}
	return proofs
}

// submitMisbehavior takes proof of misbehavior, and error id to construct the on-chain accountability event, and
// send the event of misbehavior to event channel that is listened by ethereum object to sign the reporting TX.
func (fd *FaultDetector) submitMisbehavior(m message.Msg, evidence []message.Msg, err error) {
	rule, e := errorToRule(err)
	if e != nil {
		fd.logger.Warn("error to rule", "fault detector", e)
	}
	proof := fd.eventFromProof(&Proof{
		Type:      autonity.Misbehaviour,
		Rule:      rule,
		Message:   m,
		Evidences: evidence,
	})
	// submit misbehavior proof to buffer, it will be sent once aggregated.
	fd.misbehaviourProofCh <- proof
}

func (fd *FaultDetector) checkSelfIncriminatingProposal(proposal *message.Propose) error {
	// skip processing duplicated msg.
	duplicated := engineCore.GetStore(fd.msgStore, proposal.H(), func(p *message.Propose) bool {
		return p.R() == proposal.R() &&
			p.Sender() == proposal.Sender() &&
			p.Value() == proposal.Value() &&
			p.ValidRound() == proposal.ValidRound()
	})

	if len(duplicated) > 0 {
		return errDuplicatedMsg
	}

	// account for wrong proposer.
	if !isProposerValid(fd.blockchain, proposal) {
		fd.submitMisbehavior(message.NewLightProposal(proposal), nil, errProposer)
		return errProposer
	}

	// account for equivocation
	equivocated := fd.msgStore.Get(proposal.H(), func(msg message.Msg) bool {
		// todo(youssef) : again validValue missing here
		return msg.R() == proposal.R() && msg.Code() == message.ProposalCode && msg.Sender() == proposal.Sender() && msg.Value() != proposal.Value()
	})

	if len(equivocated) > 0 {
		var equivocatedMsgs = []message.Msg{
			message.NewLightProposal(equivocated[0].(*message.Propose)),
		}
		fd.submitMisbehavior(message.NewLightProposal(proposal), equivocatedMsgs, errEquivocation)
		// we allow the equivocated msg to be stored in msg store.
		fd.msgStore.Save(proposal)
		return errEquivocation
	}

	return nil
}

func (fd *FaultDetector) checkSelfIncriminatingVote(m message.Msg) error {
	// skip process duplicated for votes.
	duplicatedMsg := fd.msgStore.Get(m.H(), func(msg message.Msg) bool {
		return msg.R() == m.R() && msg.Code() == m.Code() && msg.Sender() == m.Sender() && msg.Value() == m.Value()
	})
	if len(duplicatedMsg) > 0 {
		return errDuplicatedMsg
	}
	// account for equivocation for votes.
	equivocatedMessages := fd.msgStore.Get(m.H(), func(msg message.Msg) bool {
		return msg.R() == m.R() && msg.Code() == m.Code() && msg.Sender() == m.Sender() && msg.Value() != m.Value()
	})
	if len(equivocatedMessages) > 0 {
		fd.submitMisbehavior(m, equivocatedMessages[:1], errEquivocation)
		// we allow store equivocated msg in msg store.
		fd.msgStore.Save(m)
		return errEquivocation
	}
	return nil
}

func errorToRule(err error) (autonity.Rule, error) {
	var rule autonity.Rule
	switch {
	case errors.Is(err, errEquivocation):
		rule = autonity.Equivocation
	case errors.Is(err, errProposer):
		rule = autonity.InvalidProposer
	default:
		return rule, fmt.Errorf("errors of not provable")
	}

	return rule, nil
}

func getProposer(chain ChainContext, h uint64, r int64) (common.Address, error) {
	parentHeader := chain.GetHeaderByNumber(h - 1)
	// to prevent the panic on node shutdown.
	if parentHeader == nil {
		return common.Address{}, fmt.Errorf("cannot find parent header")
	}
	statedb, err := chain.State()
	if err != nil {
		log.Crit("could not retrieve state")
		return common.Address{}, err
	}
	proposer := chain.ProtocolContracts().Proposer(parentHeader, statedb, parentHeader.Number.Uint64(), r)
	member := parentHeader.CommitteeMember(proposer)
	if member == nil {
		return common.Address{}, fmt.Errorf("cannot find correct proposer")
	}
	return proposer, nil
}

func isProposerValid(chain ChainContext, m message.Msg) bool {
	proposer, err := getProposer(chain, m.H(), m.R())
	if err != nil {
		log.Error("get proposer err", "err", err)
		return false
	}
	return m.Sender() == proposer
}
