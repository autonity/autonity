package accountability

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math"
	"math/big"
	"sort"
	"strconv"
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
	msgGCInterval                 = 60                           // every 60 blocks to GC msg store.
	offChainAccusationProofWindow = 10                           // the time window in block for one to provide off chain innocence proof before it is escalated on chain.
	maxAccusationPerHeight        = 4                            // max number of accusation allowed to be produced by rule engine over a height against a validator.
	maxNumOfInnocenceProofCached  = 120 * maxAccusationPerHeight // 120 blocks with 4 on each height that rule engine can produce totally over a height.
	reportingSlotPeriod           = 20                           // Each AFD reporting slot holds 20 blocks, each validator response for a slot.
	//NOTE: update to below constants might require a chain fork to upgrade clients, since they impact the Accountability Event execution result. They should be turned into protocol parameters https://github.com/autonity/autonity/issues/949
	HeightRange = 256 // Default msg buffer range for AFD.
	DeltaBlocks = 10  // Wait until the GST + delta blocks to start accounting.
)

var (
	errDuplicatedMsg      = errors.New("duplicated msg")
	errEquivocation       = errors.New("equivocation")
	errNotCommitteeMsg    = errors.New("msg from none committee member")
	errProposer           = errors.New("proposal is not from proposer")
	errFutureMsg          = errors.New("future message")
	errInvalidOffenderIdx = errors.New("invalid offender index")

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

	txPool     *core.TxPool
	ethBackend ethapi.Backend
	txOpts     *bind.TransactOpts // transactor options for accountability events

	eventReporterCh chan *autonity.AccountabilityEvent
	stopRetry       chan struct{}
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

	txOpts, err := bind.NewKeyedTransactorWithChainID(nodeKey, chain.Config().ChainID)
	if err != nil {
		logger.Crit("Critical error building transactor", "err", err)
	}
	// tip needs to be >=1, otherwise accountability tx will not be broadcasted due to the txpool logic (validateTx function)
	txOpts.GasTipCap = common.Big1

	fd := &FaultDetector{
		innocenceProofBuff:    NewInnocenceProofBuffer(),
		protocolContracts:     protocolContracts,
		rateLimiter:           NewAccusationRateLimiter(),
		txPool:                txPool,
		ethBackend:            ethBackend,
		txOpts:                txOpts,
		tendermintMsgSub:      sub,
		ruleEngineBlockCh:     make(chan core.ChainEvent, 300),
		accountabilityEventCh: make(chan *autonity.AccountabilityNewAccusation),
		blockchain:            chain,
		address:               nodeAddress,
		msgStore:              ms,
		chainEventCh:          make(chan core.ChainEvent, 300),
		eventReporterCh:       make(chan *autonity.AccountabilityEvent, 10),
		stopRetry:             make(chan struct{}),
		misbehaviourProofCh:   make(chan *autonity.AccountabilityEvent, 100),
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

func (fd *FaultDetector) isHeightExpired(headHeight uint64, height uint64) bool {
	return headHeight > HeightRange && height < headHeight-HeightRange
}

func (fd *FaultDetector) SetBroadcaster(broadcaster consensus.Broadcaster) {
	fd.broadcaster = broadcaster
}

func (fd *FaultDetector) consensusMsgHandlerLoop() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
tendermintMsgLoop:
	for {
		select {
		case ev, ok := <-fd.tendermintMsgSub.Chan():
			if !ok {
				break tendermintMsgLoop
			}

			currentHeight := fd.blockchain.CurrentBlock().NumberU64()
			// handle consensus message or innocence proof messages
			switch e := ev.Data.(type) {
			case events.MessageEvent:
				if fd.isHeightExpired(currentHeight, e.Message.H()) {
					fd.logger.Debug("Fault detector: discarding old message")
					continue tendermintMsgLoop
				}
				if err := fd.processMsg(e.Message); err != nil {
					fd.logger.Warn("Detected faulty message", "return", err)
					continue tendermintMsgLoop
				}
				//TODO(lorenzo) should we gossip old height messages?
				// might be useful for accountability, but might be exploitable for DoS
				// (Jason): No, gossiping does not happens here, it requires the underlying layer to do so.
			case events.AccountabilityEvent:
				err := fd.handleOffChainAccountabilityEvent(e.Payload, e.P2pSender)
				if err != nil {
					fd.logger.Info("Accountability: Dropping peer", "peer", e.P2pSender)
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
	if height > HeightRange && height%msgGCInterval == 0 {
		threshold := height - HeightRange
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
			if ev.Block.NumberU64() > uint64(DeltaBlocks) {
				checkpoint := ev.Block.NumberU64() - uint64(DeltaBlocks)
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
				// this should never happen
				fd.logger.Crit("Can't retrieve accountability event", "id", accusation.Id.Uint64())
			}
			decodedProof, err := decodeRawProof(accusationEvent.RawProof)
			if err != nil {
				fd.logger.Error("Can't decode accusation", "err", err)
				break
			}
			// The signatures must be valid at this stage, however we have to recover the original
			// senders, hence the following call.
			if _, err := verifyProofSignatures(fd.blockchain, decodedProof); err != nil {
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
	hdr := fd.blockchain.GetHeaderByNumber(height - 1)
	if hdr == nil {
		panic("must not happen, height:" + strconv.Itoa(int(height)))
	}
	committee := hdr.Committee

	// each reporting slot contains reportingSlotPeriod block period that a unique and deterministic validator is asked to
	// be the reporter of that slot period, then at the end block of that slot, the reporter reports
	// available events. Thus, between each reporting slot, we have 5 block period to wait for
	// accountability events to be mined by network, and it is also disaster friendly that if the last
	// reporter fails, the next reporter will continue to report missing events.
	reporterIndex := (height / reportingSlotPeriod) % uint64(len(committee))

	// if validator is the reporter of the slot period, and if checkpoint block is the end block of the
	// slot, then it is time to report the collected events by this validator.
	if height%reportingSlotPeriod != 0 {
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
	close(fd.stopRetry)
	close(fd.eventReporterCh)
	fd.wg.Wait()
}

// convert the raw proofs into on-chain Proof which contains raw bytes of messages.
func (fd *FaultDetector) eventFromProof(p *Proof) *autonity.AccountabilityEvent {
	var ev = &autonity.AccountabilityEvent{
		EventType: uint8(p.Type),
		Rule:      uint8(p.Rule),
		Reporter:  fd.address,
		Offender:  p.Offender,

		Id:             common.Big0,                           // assigned contract-side
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
		// TODO(lorenzo) apply
		// whether the accusation comes from off-chain or on-chain
		// it always gets verified before we try to fetch the innocence proof
		//panic("Trying to fetch innocence proof for invalid accusation")
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

	prevotesForV := fd.msgStore.Get(func(m message.Msg) bool {
		return m.Code() == message.PrevoteCode && m.Value() == preCommit.Value() && m.R() == preCommit.R()
	}, height)

	overQuorumVotes := engineCore.OverQuorumVotes(prevotesForV, quorum)
	if overQuorumVotes == nil {
		return nil, errNoEvidenceForC1
	}

	p := fd.eventFromProof(&Proof{
		Type:      autonity.Innocence,
		Rule:      c.Rule,
		Message:   preCommit,
		Evidences: overQuorumVotes,
		Offender:  c.Offender,
	})
	return p, nil
}

// get innocent proof of accusation of rule PO from msg store.
func (fd *FaultDetector) innocenceProofPO(c *Proof) (*autonity.AccountabilityEvent, error) {
	// PO: node propose an old value with an validRound, innocent onChainProof of it should be:
	// there are quorum voting power prevotes for that value at the validRound.
	liteProposal := c.Message
	height := liteProposal.H()
	validRound := liteProposal.(*message.LightProposal).ValidRound()
	lastHeader := fd.blockchain.GetHeaderByNumber(height - 1)
	if lastHeader == nil {
		return nil, errNoParentHeader
	}
	quorum := bft.Quorum(lastHeader.TotalVotingPower())

	prevotes := fd.msgStore.Get(func(m message.Msg) bool {
		return m.Code() == message.PrevoteCode && m.R() == validRound && m.Value() == liteProposal.Value()
	}, height)

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
		Offender:  c.Offender,
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
	proposals := fd.msgStore.Get(func(m message.Msg) bool {
		return m.Code() == message.ProposalCode &&
			m.R() == prevote.R() &&
			m.Value() == prevote.Value() &&
			m.(*message.Propose).ValidRound() == -1
	}, height)

	if len(proposals) != 0 {
		p := fd.eventFromProof(&Proof{
			Type:    autonity.Innocence,
			Rule:    c.Rule,
			Message: prevote,
			Evidences: []message.Msg{
				message.NewLightProposal(proposals[0].(*message.Propose)),
			},
			Offender: c.Offender,
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

	preVotes := fd.msgStore.Get(func(m message.Msg) bool {
		return m.Code() == message.PrevoteCode && m.Value() == oldProposal.Value() && m.R() == validRound
	}, height)

	overQuorumVotes := engineCore.OverQuorumVotes(preVotes, quorum)

	if overQuorumVotes == nil {
		return nil, errNoEvidenceForPVO
	}

	p := fd.eventFromProof(&Proof{
		Type:      autonity.Innocence,
		Rule:      c.Rule,
		Message:   c.Message,
		Evidences: append(c.Evidences, overQuorumVotes...),
		Offender:  c.Offender,
	})
	return p, nil
}

// processMsg, check and submit any auto-incriminating, equivocation challenges, and then only store checked msg in msg store.
func (fd *FaultDetector) processMsg(m message.Msg) error {
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
		panic("Wrong my code for accountability")
	}

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
	proofs := fd.runRulesOverHeight(height, quorum, lastHeader.Committee)
	events := make([]*autonity.AccountabilityEvent, 0, len(proofs))

	// used to enforce max accusation per committee member per height
	accused := make(map[common.Address]uint64)

	for _, proof := range proofs {
		offender := proof.Offender

		// skip misbehaviour or accusation against self
		if fd.address == offender {
			fd.logger.Warn("found accountability proof against local node. Something went wrong, please analyze your setup and reach out on our discord", "proof", proof)
			continue
		}

		// attempt off-chain accusation resolution before escalating on-chain
		if proof.Type == autonity.Accusation {
			if accused[offender] < maxAccusationPerHeight {
				fd.addOffChainAccusation(proof)
				fd.sendOffChainAccusationMsg(proof)
				accused[offender]++
			} else {
				fd.logger.Debug("Discarding accusation, maximum already reached for this height", "offender", offender)
			}
			continue
		}

		p := fd.eventFromProof(proof)
		events = append(events, p)
	}

	return events
}

func (fd *FaultDetector) runRulesOverHeight(height uint64, quorum *big.Int, committee types.Committee) (proofs []*Proof) {
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
	proofs = append(proofs, fd.prevotesAccountabilityCheck(height, quorum, committee)...)
	proofs = append(proofs, fd.precommitsAccountabilityCheck(height, quorum, committee)...)
	return proofs
}

func (fd *FaultDetector) newProposalsAccountabilityCheck(height uint64) (proofs []*Proof) {
	// ------------New Proposal------------
	// PN:  (Mr′<r,PC|pi)∗ <--- (Mr,P|pi)
	// PN1: [nil ∨ ⊥] <--- [V]
	//
	// Since the message pattern for PN includes only messages sent by pi, we cannot raise an accusation. pi could easily
	// forge nil precommits to use as innocence proof. We can only raise a misbehaviour. If any of the precommits sent by
	// pi in rounds r' < r is for a non-nil value then we have proof of misbehaviour.

	proposalsNew := fd.msgStore.Get(func(m message.Msg) bool {
		return m.Code() == message.ProposalCode && m.(*message.Propose).ValidRound() == -1
	}, height)

	for _, p := range proposalsNew {
		proposal := p

		signer := proposal.(*message.Propose).Signer()
		// Skip if proposal is equivocated
		proposalsForR := fd.msgStore.Get(func(m message.Msg) bool {
			return m.Code() == message.ProposalCode && m.R() == proposal.R()
		}, height, signer)

		// Due to the for loop there must be at least one proposal
		if len(proposalsForR) > 1 {
			continue
		}

		//check all precommits for previous rounds from this signer are nil
		precommits := fd.msgStore.Get(func(m message.Msg) bool {
			return m.Code() == message.PrecommitCode && m.R() < proposal.R() && m.Value() != nilValue
		}, height, signer)

		if len(precommits) != 0 {
			propose := proposal.(*message.Propose)
			proof := &Proof{
				Type:          autonity.Misbehaviour,
				Rule:          autonity.PN,
				Evidences:     precommits[0:1],
				Message:       message.NewLightProposal(propose),
				Offender:      signer,
				OffenderIndex: propose.SignerIndex(),
			}
			proofs = append(proofs, proof)
			fd.logger.Info("Misbehaviour detected", "rule", "PN", "incriminated", signer)
		}
	}
	return proofs
}

func (fd *FaultDetector) oldProposalsAccountabilityCheck(height uint64, quorum *big.Int) (proofs []*Proof) {
	// ------------Old Proposal------------
	// PO: (Mr′<r,PV) ∧ (Mr′,PC|pi) ∧ (Mr′<r′′<r,P C|pi)∗ <--- (Mr,P|pi)
	// PO1: [#(Mr′,PV|V) ≥ 2f+ 1] ∧ [nil ∨ V ∨ ⊥] ∧ [nil ∨ ⊥] <--- [V]

	proposalsOld := fd.msgStore.Get(func(m message.Msg) bool {
		return m.Code() == message.ProposalCode && m.(*message.Propose).ValidRound() > -1
	}, height)

oldProposalLoop:
	for _, p := range proposalsOld {
		proposal := p
		// Check that in the valid round we see a quorum of prevotes and that there is no precommit at all or a
		// precommit for v or nil.

		signer := proposal.(*message.Propose).Signer()
		// Skip if proposal is equivocated
		proposalsForR := fd.msgStore.Get(func(m message.Msg) bool {
			return m.Code() == message.ProposalCode && m.R() == proposal.R()
		}, height, signer)

		// Due to the for loop there must be at least one proposal
		if len(proposalsForR) > 1 {
			continue oldProposalLoop
		}

		validRound := proposal.(*message.Propose).ValidRound()

		// Is there a precommit for a value other than nil or the proposed value by the current proposer in the valid
		// round? If there is, the proposer has proposed a value for which it is not locked on, thus a Proof of
		// misbehaviour can be generated.
		precommitsFromPiInVR := fd.msgStore.Get(func(m message.Msg) bool {
			return m.Code() == message.PrecommitCode && m.R() == validRound &&
				m.Value() != nilValue && m.Value() != proposal.Value()
		}, height, signer)
		if len(precommitsFromPiInVR) > 0 {
			propose := proposal.(*message.Propose)
			proof := &Proof{
				Type:          autonity.Misbehaviour,
				Rule:          autonity.PO,
				Evidences:     precommitsFromPiInVR[0:1],
				Message:       message.NewLightProposal(propose),
				Offender:      signer,
				OffenderIndex: propose.SignerIndex(),
			}
			proofs = append(proofs, proof)
			fd.logger.Info("Misbehaviour detected", "rule", "PO", "incriminated", signer)
			continue oldProposalLoop
		}

		// Is there a precommit for anything other than nil from the proposer between the valid round and the round of
		// the proposal? If there is then that implies the proposer saw 2f+1 prevotes in that round and hence it should
		// have set that round as the valid round.
		precommitsFromPiAfterVR := fd.msgStore.Get(func(m message.Msg) bool {
			return m.Code() == message.PrecommitCode && m.R() > validRound && m.R() < proposal.R() &&
				m.Value() != nilValue
		}, height, signer)

		if len(precommitsFromPiAfterVR) > 0 {
			propose := proposal.(*message.Propose)
			proof := &Proof{
				Type:          autonity.Misbehaviour,
				Rule:          autonity.PO,
				Evidences:     precommitsFromPiAfterVR[0:1],
				Message:       message.NewLightProposal(propose),
				Offender:      signer,
				OffenderIndex: propose.SignerIndex(),
			}
			proofs = append(proofs, proof)
			fd.logger.Info("Misbehaviour detected", "rule", "PO", "incriminated", signer)
			continue oldProposalLoop
		}

		// Do we see a quorum for a value other than the proposed value? If so, we have proof of misbehaviour.
		allPrevotesForValidRound := fd.msgStore.Get(func(m message.Msg) bool {
			return m.Code() == message.PrevoteCode && m.R() == validRound && m.Value() != proposal.Value()
		}, height)

		prevotesMap := make(map[common.Hash][]message.Msg)
		for _, v := range allPrevotesForValidRound {
			vote := v
			prevotesMap[vote.Value()] = append(prevotesMap[vote.Value()], vote)
		}

		for _, preVotes := range prevotesMap {
			// Here the assumption is that in a single round it is not possible to have 2 value which quorum votes,
			// this would imply at least quorum nodes are malicious which is much higher than our assumption.
			overQuorumVotes := engineCore.OverQuorumVotes(preVotes, quorum)
			if overQuorumVotes != nil {
				propose := proposal.(*message.Propose)
				proof := &Proof{
					Type:          autonity.Misbehaviour,
					Rule:          autonity.PO,
					Evidences:     overQuorumVotes,
					Message:       message.NewLightProposal(propose),
					Offender:      signer,
					OffenderIndex: propose.SignerIndex(),
				}
				proofs = append(proofs, proof)
				fd.logger.Info("Misbehaviour detected", "rule", "PO", "incriminated", signer)
				continue oldProposalLoop
			}
		}

		// Do we see a quorum of prevotes in the valid round, if not we can raise an accusation, since we cannot be sure
		// that these prevotes exist
		prevotes := fd.msgStore.Get(func(m message.Msg) bool {
			// since equivocation msgs are stored, we have to query those preVotes which has same value as the proposal.
			return m.Code() == message.PrevoteCode && m.R() == validRound && m.Value() == proposal.Value()
		}, height)

		if engineCore.OverQuorumVotes(prevotes, quorum) == nil {
			/* We do not have a quorum of prevotes for valid round here.
			* However if the propose was for a value that got committed, we do not send the accusation.
			* NOTE: this is an effective way to reduce the number of accusations and prevent accusation spamming,
			* however we assume the risk of ignoring a potentially malicious committee member.
			* Indeed the fact that the same value got committed does not rule out the fact that the suspected
			* node was misbehaving. We can just infer that if he was misbehaving, he did so in line with the decision of the network.
			* The only way to rule out misbehaviour would be to check also that the value was committed at the propose round.
			* However the commit round is not deterministic between all nodes.
			 */
			if fd.blockchain.GetBlock(proposal.Value(), proposal.H()) == nil {
				propose := proposal.(*message.Propose)
				// proposal was not committed --> send accusation
				accusation := &Proof{
					Type:          autonity.Accusation,
					Rule:          autonity.PO,
					Message:       message.NewLightProposal(propose),
					Offender:      signer,
					OffenderIndex: propose.SignerIndex(),
				}
				proofs = append(proofs, accusation)
				fd.logger.Info("🕵️ Suspicious behavior detected", "rule", "PO", "suspect", signer)
			}
		}
	}
	return proofs
}

func (fd *FaultDetector) prevotesAccountabilityCheck(height uint64, quorum *big.Int, committee types.Committee) (proofs []*Proof) {
	// ------------New and Old prevotes------------

	prevotes := fd.msgStore.Get(func(m message.Msg) bool {
		return m.Code() == message.PrevoteCode && m.Value() != nilValue
	}, height)

prevotesLoop:
	for _, p := range prevotes {
		prevote := p

		// check each signer of prevote message.
		vote := prevote.(*message.Prevote)
	signersLoop:
		for _, signerIndex := range vote.Signers().FlattenUniq() {
			signer := committee[signerIndex].Address
			// Skip the prevotes that the signer addressed as equivocated
			prevotesForR := fd.msgStore.GetEquivocatedVotes(height, prevote.R(), message.PrevoteCode, signer, prevote.Value())
			// Due to the for loop there must be at least one preVote.
			if len(prevotesForR) > 0 {
				continue signersLoop
			}

			// We need to check whether we have proposals from the prevote's round
			correspondingProposals := fd.msgStore.Get(func(m message.Msg) bool {
				return m.Code() == message.ProposalCode && m.R() == prevote.R() && m.Value() == prevote.Value()
			}, height)

			if len(correspondingProposals) == 0 {
				// if there are over quorum prevotes for this corresponding proposal's value, then it indicates current
				// peer just did not receive it. So we can skip the rising of such accusation.
				preVts := fd.msgStore.Get(func(m message.Msg) bool {
					return m.Code() == message.PrevoteCode && m.R() == prevote.R() && m.Value() == prevote.Value()
				}, height)

				if engineCore.OverQuorumVotes(preVts, quorum) == nil {
					/* The rule for this accusation could be PVO as well since we don't have the corresponding proposal.
					* If the prevote was for a value that got committed, we do not send the accusation.
					* NOTE: this is an effective way to reduce the number of accusations and prevent accusation spamming,
					* however we assume the risk of ignoring a potentially malicious committee member.
					* Indeed the fact that the same value got committed does not rule out the fact that the suspected
					* node was misbehaving. We can just infer that if he was misbehaving, he did so in line with the decision of the network.
					* The only way to rule out misbehaviour would be to check also that the value was committed at the prevote round.
					* However the commit round is not deterministic between all nodes.
					 */
					if fd.blockchain.GetBlock(prevote.Value(), prevote.H()) == nil {
						accusation := &Proof{
							Type:          autonity.Accusation,
							Rule:          autonity.PVN,
							Message:       prevote,
							Offender:      signer,
							OffenderIndex: signerIndex,
						}
						proofs = append(proofs, accusation)
						fd.logger.Info("🕵️ Suspicious behavior detected", "rule", "PVN", "suspect", signer)
					}
				}
				continue prevotesLoop // we have no corresponding proposal, so we cannot check new and old prevote rules
			}

			// We need to ensure that we keep all proposals in the message store, so that we have the maximum chance of
			// finding justification for prevotes. This is to account for equivocation where the proposer send 2 proposals
			// with the same value but different valid rounds to different nodes. We can't penalise the signer of prevote
			// since we can't tell which proposal they received. We just want to find a set of message which fit the rule.
			// Therefore, we need to check all the proposals to find a single one which shows the current prevote is
			// valid.
			var prevotesProofs []*Proof
			for _, cp := range correspondingProposals {
				proposal := cp.(*message.Propose)
				var proof *Proof
				if proposal.ValidRound() == -1 {
					proof = fd.newPrevotesAccountabilityCheck(height, prevote, proposal, signer, signerIndex)
				} else {
					proof = fd.oldPrevotesAccountabilityCheck(height, quorum, proposal, prevote, signer, signerIndex)
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
						continue signersLoop
					}
				}

				// There are no corresponding proposal for which the current prevote is valid. We prioritise misbehaviours over
				// accusation since they can be easily proved.
				for _, proof := range prevotesProofs {
					if proof.Type == autonity.Misbehaviour {
						proofs = append(proofs, proof)
						continue signersLoop
					}
				}

				// There were no misbehaviours for the current prevote, therefore, pick the first accusation
				proofs = append(proofs, prevotesProofs[0])
			}
		}
	}
	return proofs
}

func (fd *FaultDetector) newPrevotesAccountabilityCheck(height uint64, prevote message.Msg,
	correspondingProposal *message.Propose, signer common.Address, signerIndex int) (proof *Proof) {
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
	precommitsFromPi := fd.msgStore.Get(func(m message.Msg) bool {
		return m.Code() == message.PrecommitCode && m.R() < prevote.R()
	}, height, signer)

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
				// we found the latest non-nil precommit and we don't have gaps in the following ones
				pc := precommitsFromPi[i]

				// check for equivocation. If present, bail out on the checking of this rule. Remote peer has already been punished for equivocation
				precommitsAtRPrime := fd.msgStore.GetEquivocatedVotes(height, pc.R(), message.PrecommitCode, signer, pc.Value())
				if len(precommitsAtRPrime) > 0 {
					break
				}

				// if precommit at r' is for V, then all good --> no misbehaviour
				if pc.Value() == prevote.Value() {
					break
				}

				// precommit at r' is not for V --> remote peer is malicious
				fd.logger.Info("Misbehaviour detected", "rule", "PVN", "incriminated", signer)
				proof := &Proof{
					Type:          autonity.Misbehaviour,
					Rule:          autonity.PVN,
					Message:       prevote,
					Offender:      signer,
					OffenderIndex: signerIndex,
				}
				// to guarantee this prevote is for a new proposal that is the PVN rule account for, otherwise in
				// prevote for an old proposal, it is valid for one to prevote it if lockedRound <= vr, thus the
				// round jump is valid. This prevents from rising a PVN misbehavior proof from a malicious fault
				// detector by using prevote for an old proposal to challenge an honest slow validator.
				proof.Evidences = append(proof.Evidences, message.NewLightProposal(correspondingProposal))
				proof.Evidences = append(proof.Evidences, precommitsFromPi[i:]...)
				return proof
			}
			if i > 0 {
				r = rPrime
				rPrime = precommitsFromPi[i-1].R()
			}
		}
	}
	/* we end up here if:
	* - pi never locked (i.e. precommitted) before sending this prevote
	* - pi always precommitted nil
	* - we have gaps in the precommits
	* - latest non-nil precommit (at r') is for V
	* - latest non-nil precommit is equivocated
	 */
	return nil
}

func (fd *FaultDetector) oldPrevotesAccountabilityCheck(height uint64, quorum *big.Int,
	correspondingProposal *message.Propose, prevote message.Msg, signer common.Address, signerIndex int) (proof *Proof) {
	currentR := correspondingProposal.R()
	validRound := correspondingProposal.ValidRound()

	// If there is a prevote for an old proposal then pi can only vote for v or send nil (see line 28 and 29 of
	// tendermint pseudocode), therefore if in the valid round there is a quorum for a value other than v, we know pi
	// prevoted incorrectly. If the proposal was a bad proposal, then pi should not have voted for it, thus we do not
	// need to make sure whether the proposal is correct or not (which we would in the proposal checking rules, however,
	// a bad proposal will still exist in our message store, and it shouldn't have an impact on the checking of prevotes).

	allPrevotesForValidRound := fd.msgStore.Get(func(m message.Msg) bool {
		return m.Code() == message.PrevoteCode && m.R() == validRound && m.Value() != correspondingProposal.Value()
	}, height)

	prevotesMap := make(map[common.Hash][]message.Msg)
	for _, p := range allPrevotesForValidRound {
		vote := p
		prevotesMap[vote.Value()] = append(prevotesMap[vote.Value()], vote)
	}

	for _, preVotes := range prevotesMap {
		// Here the assumption is that in a single round it is not possible to have 2 value which quorum votes,
		// this would imply at least quorum nodes are malicious which is much higher than our assumption.
		overQuorumVotes := engineCore.OverQuorumVotes(preVotes, quorum)
		if overQuorumVotes != nil {
			fd.logger.Info("Misbehaviour detected", "rule", "PV0", "incriminated", signer)
			proof := &Proof{
				Type:          autonity.Misbehaviour,
				Rule:          autonity.PVO,
				Message:       prevote,
				Offender:      signer,
				OffenderIndex: signerIndex,
			}
			proof.Evidences = append(proof.Evidences, message.NewLightProposal(correspondingProposal))
			proof.Evidences = append(proof.Evidences, overQuorumVotes...)
			return proof
		}
	}

	prevotesForVFromValidRound := fd.msgStore.Get(func(m message.Msg) bool {
		return m.Code() == message.PrevoteCode && m.R() == validRound && m.Value() == correspondingProposal.Value()
	}, height)

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

		precommitsFromPi := fd.msgStore.Get(func(m message.Msg) bool {
			return m.Code() == message.PrecommitCode && m.R() > validRound && m.R() < currentR
		}, height, signer)

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
				fd.logger.Info("Misbehaviour detected", "rule", "PVO12", "incriminated", signer)
				proof := &Proof{
					Type:          autonity.Misbehaviour,
					Rule:          autonity.PVO12,
					Message:       prevote,
					Offender:      signer,
					OffenderIndex: signerIndex,
				}
				proof.Evidences = append(proof.Evidences, message.NewLightProposal(correspondingProposal))
				proof.Evidences = append(proof.Evidences, precommitsFromPi...)
				return proof
			}
		}
	}

	// if there is no misbehaviour of the prevote msg addressed, then we lastly check accusation.
	if overQuorumPrevotesForVFromValidRound == nil {
		/* We do not have a quorum of prevotes for valid round here.
		* However if the prevote was for a value that got committed, we do not send the accusation.
		* NOTE: this is an effective way to reduce the number of accusations and prevent accusation spamming,
		* however we assume the risk of ignoring a potentially malicious committee member.
		* Indeed the fact that the same value got committed does not rule out the fact that the suspected
		* node was misbehaving. We can just infer that if he was misbehaving, he did so in line with the decision of the network.
		* The only way to rule out misbehaviour would be to check also that the value was committed at the prevote round.
		* However the commit round is not deterministic between all nodes.
		 */
		if fd.blockchain.GetBlock(prevote.Value(), prevote.H()) == nil {
			fd.logger.Info("🕵️ Suspicious behavior detected", "rule", "PVO", "suspect", signer)
			return &Proof{
				Type:          autonity.Accusation,
				Rule:          autonity.PVO,
				Message:       prevote,
				Evidences:     []message.Msg{message.NewLightProposal(correspondingProposal)},
				Offender:      signer,
				OffenderIndex: signerIndex,
			}
		}
	}

	return nil
}

func (fd *FaultDetector) precommitsAccountabilityCheck(height uint64, quorum *big.Int, committee types.Committee) (proofs []*Proof) {
	// ------------precommits------------
	// C: [Mr,P|proposer(r)] ∧ [Mr,PV] <--- [Mr,PC|pi]
	// C1: [V:Valid(V)] ∧ [#(V) ≥ 2f+ 1] <--- [V]

	precommits := fd.msgStore.Get(func(m message.Msg) bool {
		return m.Code() == message.PrecommitCode && m.Value() != nilValue
	}, height)

	for _, preC := range precommits {
		precommit := preC

		vote := precommit.(*message.Precommit)
	singersLoop:
		for _, signerIndex := range vote.Signers().FlattenUniq() {
			signer := committee[signerIndex].Address
			// Skip if preCommit is equivocated
			precommitsForR := fd.msgStore.GetEquivocatedVotes(height, precommit.R(), message.PrecommitCode, signer, precommit.Value())
			// Due to the for loop there must be at least one preCommit.
			if len(precommitsForR) > 0 {
				continue singersLoop
			}

			// Do we see a quorum for a value other than the proposed value? If so, we have proof of misbehaviour.
			allPrevotesForR := fd.msgStore.Get(func(m message.Msg) bool {
				return m.Code() == message.PrevoteCode && m.R() == precommit.R() && m.Value() != precommit.Value()
			}, height)

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
						Type:          autonity.Misbehaviour,
						Rule:          autonity.C,
						Evidences:     overQuorumVotes,
						Message:       precommit,
						Offender:      signer,
						OffenderIndex: signerIndex,
					}
					proofs = append(proofs, proof)
					fd.logger.Info("Misbehaviour detected", "rule", "C", "incriminated", signer)
					continue singersLoop
				}
			}

			// Do we see a quorum of prevotes in the same round? if not we can raise an accusation, since we cannot be sure
			// that these prevotes do exist, this block also covers the Accusation of C since if over quorum prevotes for
			// V indicates that the corresponding proposal of V do exist, thus we don't need to raise accusation for the missing
			// proposal since over 2/3 member should all ready received it
			prevotes := fd.msgStore.Get(func(m message.Msg) bool {
				// since equivocation msgs are stored, we have to query those preVotes which has same value as the proposal.
				return m.Code() == message.PrevoteCode && m.R() == precommit.R() && m.Value() == precommit.Value()
			}, height)

			if engineCore.OverQuorumVotes(prevotes, quorum) == nil {
				/* We do not have a quorum of prevotes for this precommit to be justified.
				* However if the precommit was for a value that got committed, we do not send the accusation.
				* NOTE: this is an effective way to reduce the number of accusations and prevent accusation spamming,
				* however we assume the risk of ignoring a potentially malicious committee member.
				* Indeed the fact that the same value got committed does not rule out the fact that the suspected
				* node was misbehaving. We can just infer that if he was misbehaving, he did so in line with the decision of the network.
				* The only way to rule out misbehaviour would be to check also that the value was committed at the precommit round.
				* However the commit round is not deterministic between all nodes.
				 */
				if fd.blockchain.GetBlock(precommit.Value(), precommit.H()) == nil {
					accusation := &Proof{
						Type:          autonity.Accusation,
						Rule:          autonity.C1,
						Message:       precommit,
						Offender:      signer,
						OffenderIndex: signerIndex,
					}
					proofs = append(proofs, accusation)
					fd.logger.Info("🕵️ Suspicious behavior detected", "rule", "C1", "suspect", signer)
				}
			}
		}
	}
	return proofs
}

// submitMisbehavior takes proof of misbehavior, and error id to construct the on-chain accountability event, and
// send the event of misbehavior to event channel that is listened by ethereum object to sign the reporting TX.
func (fd *FaultDetector) submitMisbehavior(m message.Msg, evidence []message.Msg, err error,
	offender common.Address, offenderIndex int) {
	rule := errorToRule(err)
	proof := fd.eventFromProof(&Proof{
		Type:          autonity.Misbehaviour,
		Rule:          rule,
		Message:       m,
		Evidences:     evidence,
		Offender:      offender,
		OffenderIndex: offenderIndex,
	})
	// submit misbehavior proof to buffer, it will be sent once aggregated.
	fd.misbehaviourProofCh <- proof
}

func (fd *FaultDetector) checkSelfIncriminatingProposal(proposal *message.Propose) error {
	// skip processing duplicated msg.
	duplicated := engineCore.GetStore(fd.msgStore, proposal.H(), func(p *message.Propose) bool {
		return p.R() == proposal.R() &&
			p.Signer() == proposal.Signer() &&
			p.Value() == proposal.Value() &&
			p.ValidRound() == proposal.ValidRound()
	})

	if len(duplicated) > 0 {
		return errDuplicatedMsg
	}

	// account for wrong proposer.
	if !isProposerValid(fd.blockchain, proposal) {
		fd.submitMisbehavior(message.NewLightProposal(proposal), nil, errProposer, proposal.Signer(), proposal.SignerIndex())
		return errProposer
	}

	// account for equivocation
	equivocated := engineCore.GetStore(fd.msgStore, proposal.H(), func(p *message.Propose) bool {
		// todo(youssef) : again validValue missing here
		return p.R() == proposal.R() &&
			p.Signer() == proposal.Signer() &&
			p.Value() != proposal.Value() &&
			p.ValidRound() == proposal.ValidRound()
	})

	if len(equivocated) > 0 {
		var equivocatedMsgs = []message.Msg{
			message.NewLightProposal(equivocated[0]),
		}
		fd.submitMisbehavior(message.NewLightProposal(proposal), equivocatedMsgs, errEquivocation, proposal.Signer(), proposal.SignerIndex())
		// we allow the equivocated msg to be stored in msg store.
		fd.msgStore.Save(proposal, nil)
		return errEquivocation
	}
	fd.msgStore.Save(proposal, nil)
	return nil
}

func (fd *FaultDetector) checkSelfIncriminatingVote(m message.Msg) error {
	// skip process duplicated for votes.
	duplicatedMsg := fd.msgStore.Get(func(msg message.Msg) bool {
		return msg.Hash() == m.Hash()
	}, m.H())
	if len(duplicatedMsg) > 0 {
		return errDuplicatedMsg
	}

	// account for equivocation for votes.
	signers := m.(message.Vote).Signers()
	var err error
	lastHeader := fd.blockchain.GetHeaderByNumber(m.H() - 1)
	for _, signerIndex := range signers.FlattenUniq() {
		signer := lastHeader.Committee[signerIndex].Address
		equivocatedMessages := fd.msgStore.GetEquivocatedVotes(m.H(), m.R(), m.Code(), signer, m.Value())
		if len(equivocatedMessages) > 0 {
			fd.submitMisbehavior(m, equivocatedMessages[:1], errEquivocation, signer, signerIndex)
			err = errEquivocation
		}
	}

	fd.msgStore.Save(m, lastHeader.Committee)
	return err
}

func errorToRule(err error) autonity.Rule {
	var rule autonity.Rule
	switch {
	case errors.Is(err, errEquivocation):
		rule = autonity.Equivocation
	case errors.Is(err, errProposer):
		rule = autonity.InvalidProposer
	default:
		panic("unknown error to accountability rule mapping")
	}

	return rule
}

// TODO: this function basically reimplements committee.GetProposer
// it would be better to use that function but it requires sharing the CommitteeSet between Core and the FD.
// It would reduce code repetition, however the proposer cache is already shared so it would not improve performance.
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
	signer := m.(*message.Propose).Signer()
	return signer == proposer
}
