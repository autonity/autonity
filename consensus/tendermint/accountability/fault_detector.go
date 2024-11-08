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
	CommitteeByHeight(height uint64) (*types.Committee, error)
}

const (
	msgGCInterval                 = 60                           // every 60 blocks to GC msg store.
	offChainAccusationProofWindow = 10                           // the time window in block for one to provide off chain innocence proof before it is escalated on chain.
	maxAccusationPerHeight        = 4                            // max number of accusation allowed to be produced by rule engine over a height against a validator.
	maxNumOfInnocenceProofCached  = 120 * maxAccusationPerHeight // 120 blocks with 4 on each height that rule engine can produce totally over a height.
	reportingSlotPeriod           = 20                           // Each AFD reporting slot holds 20 blocks, each validator response for a slot.
	//NOTE: update to below constants might require a chain fork to upgrade clients, since they impact the Accountability Event execution result. They should be turned into protocol parameters https://github.com/autonity/autonity/issues/949
	HeightRange = 256 // Default msg buffer range for AFD.
	DeltaBlocks = 10  // Wait until the GST + delta blocks to start accounting
)

var (
	errDuplicatedMsg      = errors.New("duplicated msg")
	errEquivocation       = errors.New("equivocation")
	errNotCommitteeMsg    = errors.New("msg from none committee member")
	errProposer           = errors.New("proposal is not from proposer")
	errInvalidOffenderIdx = errors.New("invalid offender index")

	errNoEvidenceForPO  = errors.New("no proof of innocence found for rule PO")
	errNoEvidenceForPVN = errors.New("no proof of innocence found for rule PVN")
	errNoEvidenceForPVO = errors.New("no proof of innocence found for rule PVO")
	errNoEvidenceForC1  = errors.New("no proof of innocence found for rule C1")

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
					if !errors.Is(err, errDuplicatedMsg) {
						fd.logger.Warn("Detected faulty message", "err", err)
					} else {
						// duplicated messages can arrive here if we receive an aggregate from a remote peer
						// and at the same time we computed the same aggregate locally.
						// No need to raise a warning level log.
						fd.logger.Debug("Detected faulty message", "err", err)
					}
					continue tendermintMsgLoop
				}
			case events.OldMessageEvent:
				if fd.isHeightExpired(currentHeight, e.Message.H()) {
					fd.logger.Debug("Fault detector: discarding old message")
					continue tendermintMsgLoop
				}
				if err := fd.processMsg(e.Message); err != nil {
					if !errors.Is(err, errDuplicatedMsg) {
						fd.logger.Warn("Detected faulty message", "err", err)
					} else {
						// duplicated messages can arrive here if we receive an aggregate from a remote peer
						// and at the same time we computed the same aggregate locally.
						// No need to raise a warning level log.
						fd.logger.Debug("Detected faulty message", "err", err)
					}
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

			h := decodedProof.Message.H()
			committee, err := fd.blockchain.CommitteeByHeight(h)
			if err != nil {
				fd.logger.Error("Can't retrieve committee for message", "err", err, "height", h)
				break
			}

			// The signatures must be valid at this stage, however we have to recover the original
			// senders, hence the following call.
			if err = verifyProofSignatures(committee, decodedProof); err != nil {
				fd.logger.Error("Can't verify proof signatures", "err", err)
				break
			}

			innocenceProof, err := fd.innocenceProof(decodedProof, committee)
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
	committee, err := fd.blockchain.CommitteeByHeight(height)
	if err != nil {
		fd.logger.Crit("Can't retrieve committee for message", "err", err, "height", height)
	}

	// each reporting slot contains reportingSlotPeriod block period that a unique and deterministic validator is asked to
	// be the reporter of that slot period, then at the end block of that slot, the reporter reports
	// available events. Thus, between each reporting slot, we have 5 block period to wait for
	// accountability events to be mined by network, and it is also disaster friendly that if the last
	// reporter fails, the next reporter will continue to report missing events.
	reporterIndex := (height / reportingSlotPeriod) % uint64(committee.Len())

	// if validator is the reporter of the slot period, and if checkpoint block is the end block of the
	// slot, then it is time to report the collected events by this validator.
	if height%reportingSlotPeriod != 0 {
		return false
	}
	// todo(youssef): this seems like a non-committee member can't send a proof/ do we want that?
	return committee.Members[reporterIndex].Address == fd.address
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
func (fd *FaultDetector) eventFromProof(p *Proof, offender common.Address) *autonity.AccountabilityEvent {
	var ev = &autonity.AccountabilityEvent{
		EventType: uint8(p.Type),
		Rule:      uint8(p.Rule),
		Reporter:  fd.address,
		Offender:  offender,

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
func (fd *FaultDetector) innocenceProof(p *Proof, committee *types.Committee) (*autonity.AccountabilityEvent, error) {
	// the protocol contains below provable accusations.
	switch p.Rule {
	case autonity.PO:
		return fd.innocenceProofPO(p, committee)
	case autonity.PVN:
		return fd.innocenceProofPVN(p, committee)
	case autonity.PVO:
		return fd.innocenceProofPVO(p, committee)
	case autonity.C1:
		return fd.innocenceProofC1(p, committee)
	default:
		// whether the accusation comes from off-chain or on-chain
		// it always gets verified before we try to fetch the innocence proof
		panic("Trying to fetch innocence proof for invalid accusation")
	}
}

// get innocent proof of accusation of rule C1 from msg store.
func (fd *FaultDetector) innocenceProofC1(c *Proof, committee *types.Committee) (*autonity.AccountabilityEvent, error) {
	precommit := c.Message
	height := precommit.H()

	// compute quorum
	quorum := bft.Quorum(committee.TotalVotingPower())

	// check if we have quorum voting power for V
	if fd.msgStore.PrevotesPowerFor(height, precommit.R(), precommit.Value()).Cmp(quorum) < 0 {
		// we don't have quorum, cannot defend the accusation!
		return nil, errNoEvidenceForC1
	}

	prevotesForV := fd.msgStore.GetPrevotes(height, func(m *message.Prevote) bool {
		return m.Value() == precommit.Value() && m.R() == precommit.R()
	})

	// although we checked over quorum prevotes for V at the round of precommit, however due to
	// off-chain accusation handler can run in different routine, thus GC of msg store could potentially
	// delete the prevotes.
	if len(prevotesForV) == 0 {
		return nil, errNoEvidenceForC1
	}

	// fast aggregate quorum prevotes for V into single one.
	evidences := make([]message.Msg, 1)
	evidences[0] = prevotesForV[0]
	if len(prevotesForV) > 1 {
		evidences[0] = AggregateSamePrevotes(prevotesForV)
	}
	p := fd.eventFromProof(&Proof{
		Type:          autonity.Innocence,
		Rule:          c.Rule,
		Message:       precommit,
		Evidences:     evidences,
		OffenderIndex: c.OffenderIndex,
	}, committee.Members[c.OffenderIndex].Address)
	return p, nil
}

// get innocent proof of accusation of rule PO from msg store.
func (fd *FaultDetector) innocenceProofPO(c *Proof, committee *types.Committee) (*autonity.AccountabilityEvent, error) {
	// PO: node propose an old value with an validRound, innocent onChainProof of it should be:
	// there are quorum voting power prevotes for that value at the validRound.
	liteProposal := c.Message
	height := liteProposal.H()
	validRound := liteProposal.(*message.LightProposal).ValidRound()

	// compute quorum
	quorum := bft.Quorum(committee.TotalVotingPower())

	// check if we have quorum voting power for V at validRound
	if fd.msgStore.PrevotesPowerFor(height, validRound, liteProposal.Value()).Cmp(quorum) < 0 {
		// we don't have quorum, cannot defend the accusation!
		return nil, errNoEvidenceForPO
	}

	prevotes := fd.msgStore.GetPrevotes(height, func(m *message.Prevote) bool {
		return m.R() == validRound && m.Value() == liteProposal.Value()
	})

	// although we checked over quorum prevotes for the lite proposal at vr, however due to
	// off-chain accusation handler runs in different routine, thus GC of msg store could potentially
	// delete the prevotes.
	if len(prevotes) == 0 {
		return nil, errNoEvidenceForPO
	}

	// fast aggregate quorum prevotes into single one.
	evidences := make([]message.Msg, 1)
	evidences[0] = prevotes[0]
	if len(prevotes) > 1 {
		evidences[0] = AggregateSamePrevotes(prevotes)
	}

	p := fd.eventFromProof(&Proof{
		Type:          autonity.Innocence,
		Rule:          c.Rule,
		Message:       liteProposal,
		Evidences:     evidences,
		OffenderIndex: c.OffenderIndex,
	}, committee.Members[c.OffenderIndex].Address)
	return p, nil
}

// get innocent proof of accusation of rule PVN from msg store.
func (fd *FaultDetector) innocenceProofPVN(c *Proof, committee *types.Committee) (*autonity.AccountabilityEvent, error) {
	// get innocent proofs for PVN, for a prevote that vote for a new value,
	// then there must be a proposal for this new value.
	prevote := c.Message
	height := prevote.H()
	// the only proof of innocence of PVN accusation is that there exist a corresponding proposal
	proposals := fd.msgStore.GetProposals(height, func(m *message.Propose) bool {
		return m.R() == prevote.R() &&
			m.Value() == prevote.Value() &&
			m.ValidRound() == -1
	})

	if len(proposals) != 0 {
		p := fd.eventFromProof(&Proof{
			Type:    autonity.Innocence,
			Rule:    c.Rule,
			Message: prevote,
			Evidences: []message.Msg{
				message.NewLightProposal(proposals[0]),
			},
			OffenderIndex: c.OffenderIndex,
		}, committee.Members[c.OffenderIndex].Address)
		return p, nil
	}
	return nil, errNoEvidenceForPVN
}

// get innocent proof of accusation of rule PVO from msg store, it collects quorum preVotes for the value voted at a valid round.
func (fd *FaultDetector) innocenceProofPVO(c *Proof, committee *types.Committee) (*autonity.AccountabilityEvent, error) {
	// get innocent proofs for PVO, collect quorum preVotes at the valid round of the old proposal.
	oldProposal := c.Evidences[0]
	height := oldProposal.H()
	validRound := oldProposal.(*message.LightProposal).ValidRound()

	// compute quorum
	quorum := bft.Quorum(committee.TotalVotingPower())

	// check if we have quorum voting power for V at validRound
	if fd.msgStore.PrevotesPowerFor(height, validRound, oldProposal.Value()).Cmp(quorum) < 0 {
		// we don't have quorum, cannot defend the accusation!
		return nil, errNoEvidenceForPVO
	}

	prevotes := fd.msgStore.GetPrevotes(height, func(m *message.Prevote) bool {
		return m.Value() == oldProposal.Value() && m.R() == validRound
	})

	// although we checked over quorum prevotes for old proposal at vr, however due to
	// off-chain accusation handler runs in different routine, thus GC of msg store could potentially
	// delete the prevotes.
	if len(prevotes) == 0 {
		return nil, errNoEvidenceForPVO
	}

	// fast aggregate quorum prevotes into single one.
	evidences := make([]message.Msg, 1)
	evidences[0] = prevotes[0]
	if len(prevotes) > 1 {
		evidences[0] = AggregateSamePrevotes(prevotes)
	}

	p := fd.eventFromProof(&Proof{
		Type:          autonity.Innocence,
		Rule:          c.Rule,
		Message:       c.Message,
		Evidences:     append(c.Evidences, evidences...),
		OffenderIndex: c.OffenderIndex,
	}, committee.Members[c.OffenderIndex].Address)
	return p, nil
}

// processMsg, check and submit any auto-incriminating, equivocation challenges, and then only store checked msg in msg store.
func (fd *FaultDetector) processMsg(m message.Msg) error {
	switch msg := m.(type) {
	case *message.Propose:
		if err := fd.checkSelfIncriminatingProposal(msg); err != nil {
			return err
		}
	case *message.Prevote:
		if err := fd.checkSelfIncriminatingPrevote(msg); err != nil {
			return err
		}
	case *message.Precommit:
		if err := fd.checkSelfIncriminatingPrecommit(msg); err != nil {
			return err
		}
	default:
		panic("Wrong msg code for accountability")
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

	committee, err := fd.blockchain.CommitteeByHeight(height)
	if err != nil {
		fd.logger.Crit("cannot find committee for height", "err", err, "height", height)
	}
	quorum := bft.Quorum(committee.TotalVotingPower())
	proofs := fd.runRulesOverHeight(height, quorum, committee)
	events := make([]*autonity.AccountabilityEvent, 0, len(proofs))

	// used to enforce max accusation per committee member per height
	accused := make(map[common.Address]uint64)

	for _, proof := range proofs {
		offender := committee.Members[proof.OffenderIndex].Address

		// skip misbehaviour or accusation against self
		if fd.address == offender {
			fd.logger.Warn("found accountability proof against local node. Something went wrong, please analyze your setup and reach out on our discord", "proof", proof)
			continue
		}

		// attempt off-chain accusation resolution before escalating on-chain
		if proof.Type == autonity.Accusation {
			if accused[offender] < maxAccusationPerHeight {
				fd.addOffChainAccusation(proof)
				fd.sendOffChainAccusationMsg(proof, committee)
				accused[offender]++
			} else {
				fd.logger.Debug("Discarding accusation, maximum already reached for this height", "offender", offender)
			}
			continue
		}

		p := fd.eventFromProof(proof, committee.Members[proof.OffenderIndex].Address)
		events = append(events, p)
	}

	return events
}

func (fd *FaultDetector) runRulesOverHeight(height uint64, quorum *big.Int, committee *types.Committee) (proofs []*Proof) {
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
	// In this rule engine context, the symbol `pi` stands for a consensus participant with unique identity `i`.

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

	proposalsNew := fd.msgStore.GetProposals(height, func(m *message.Propose) bool {
		return m.ValidRound() == -1
	})

	for _, proposal := range proposalsNew {
		signerIndex := proposal.SignerIndex()

		// Skip if proposal is equivocated
		proposalsForR := fd.msgStore.GetProposals(height, func(m *message.Propose) bool {
			return m.R() == proposal.R() && m.SignerIndex() == signerIndex && (m.Value() != proposal.Value() || m.ValidRound() != proposal.ValidRound())
		})
		if len(proposalsForR) > 0 {
			continue
		}

		//check all precommits for previous rounds from this signer are nil
		precommits := fd.msgStore.GetPrecommits(height, func(m *message.Precommit) bool {
			return m.R() < proposal.R() && m.Value() != nilValue && m.Signers().Contains(signerIndex)
		})

		if len(precommits) != 0 {
			proof := &Proof{
				Type:          autonity.Misbehaviour,
				Rule:          autonity.PN,
				Evidences:     []message.Msg{precommits[0]},
				Message:       message.NewLightProposal(proposal),
				OffenderIndex: signerIndex,
			}
			proofs = append(proofs, proof)
			fd.logger.Info("Misbehaviour detected", "rule", "PN", "incriminated", proposal.Signer())
		}
	}
	return proofs
}

func (fd *FaultDetector) oldProposalsAccountabilityCheck(height uint64, quorum *big.Int) (proofs []*Proof) {
	// ------------Old Proposal------------
	// PO: (Mr′<r,PV) ∧ (Mr′,PC|pi) ∧ (Mr′<r′′<r,P C|pi)∗ <--- (Mr,P|pi)
	// PO1: [#(Mr′,PV|V) ≥ 2f+ 1] ∧ [nil ∨ V ∨ ⊥] ∧ [nil ∨ ⊥] <--- [V]

	proposalsOld := fd.msgStore.GetProposals(height, func(m *message.Propose) bool {
		return m.ValidRound() > -1
	})

oldProposalLoop:
	for _, proposal := range proposalsOld {
		// Check that in the valid round we see a quorum of prevotes and that there is no precommit at all or a
		// precommit for v or nil.

		signer := proposal.Signer()
		signerIndex := proposal.SignerIndex()
		validRound := proposal.ValidRound()

		// Skip if proposal is equivocated
		proposalsForR := fd.msgStore.GetProposals(height, func(m *message.Propose) bool {
			return m.R() == proposal.R() && m.SignerIndex() == signerIndex && (m.Value() != proposal.Value() || m.ValidRound() != validRound)
		})
		if len(proposalsForR) > 0 {
			continue oldProposalLoop
		}

		// Is there a precommit for a value other than nil or the proposed value by the current proposer in the valid
		// round? If there is, the proposer has proposed a value for which it is not locked on, thus a Proof of
		// misbehaviour can be generated.
		precommitsFromPiInVR := fd.msgStore.GetPrecommits(height, func(m *message.Precommit) bool {
			return m.R() == validRound && m.Value() != nilValue && m.Value() != proposal.Value() && m.Signers().Contains(signerIndex)
		})
		if len(precommitsFromPiInVR) > 0 {
			proof := &Proof{
				Type:          autonity.Misbehaviour,
				Rule:          autonity.PO,
				Evidences:     []message.Msg{precommitsFromPiInVR[0]},
				Message:       message.NewLightProposal(proposal),
				OffenderIndex: signerIndex,
			}
			proofs = append(proofs, proof)
			fd.logger.Info("Misbehaviour detected", "rule", "PO", "incriminated", signer)
			continue oldProposalLoop
		}

		// Is there a precommit for anything other than nil from the proposer between the valid round and the round of
		// the proposal? If there is then that implies the proposer saw 2f+1 prevotes in that round and hence it should
		// have set that round as the valid round.
		precommitsFromPiAfterVR := fd.msgStore.GetPrecommits(height, func(m *message.Precommit) bool {
			return m.R() > validRound && m.R() < proposal.R() && m.Value() != nilValue && m.Signers().Contains(signerIndex)
		})

		if len(precommitsFromPiAfterVR) > 0 {
			proof := &Proof{
				Type:          autonity.Misbehaviour,
				Rule:          autonity.PO,
				Evidences:     []message.Msg{precommitsFromPiAfterVR[0]},
				Message:       message.NewLightProposal(proposal),
				OffenderIndex: signerIndex,
			}
			proofs = append(proofs, proof)
			fd.logger.Info("Misbehaviour detected", "rule", "PO", "incriminated", signer)
			continue oldProposalLoop
		}

		// Do we see a quorum for a value other than the proposed value? If so, we have proof of misbehaviour.
		alternativeQuorum := fd.msgStore.SearchQuorum(height, validRound, proposal.Value(), quorum)
		// Here the assumption is that in a single round it is not possible to have 2 value which quorum votes,
		// this would imply at least quorum nodes are malicious which is much higher than our assumption.
		if len(alternativeQuorum) > 0 {
			// fast aggregate quorum prevotes into single one.
			evidences := make([]message.Msg, 1)
			evidences[0] = alternativeQuorum[0]
			if len(alternativeQuorum) > 1 {
				evidences[0] = AggregateSamePrevotes(alternativeQuorum)
			}

			proof := &Proof{
				Type:          autonity.Misbehaviour,
				Rule:          autonity.PO,
				Evidences:     evidences,
				Message:       message.NewLightProposal(proposal),
				OffenderIndex: signerIndex,
			}
			proofs = append(proofs, proof)
			fd.logger.Info("Misbehaviour detected", "rule", "PO", "incriminated", signer)
			continue oldProposalLoop
		}

		// Do we see a quorum of prevotes in the valid round? if not we can raise an accusation, since we cannot be sure
		// that these prevotes exist
		if fd.msgStore.PrevotesPowerFor(height, validRound, proposal.Value()).Cmp(quorum) < 0 {
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
				// proposal was not committed --> send accusation
				accusation := &Proof{
					Type:          autonity.Accusation,
					Rule:          autonity.PO,
					Message:       message.NewLightProposal(proposal),
					OffenderIndex: signerIndex,
				}
				proofs = append(proofs, accusation)
				fd.logger.Info("🕵️ Suspicious behavior detected", "rule", "PO", "suspect", signer)
			}
		}
	}
	return proofs
}

func (fd *FaultDetector) prevotesAccountabilityCheck(height uint64, quorum *big.Int, committee *types.Committee) (proofs []*Proof) {
	// ------------New and Old prevotes------------

	prevotes := fd.msgStore.GetPrevotes(height, func(m *message.Prevote) bool {
		return m.Value() != nilValue
	})

	for _, prevote := range prevotes {
	signersLoop:
		for _, signerIndex := range prevote.Signers().FlattenUniq() {
			signer := committee.Members[signerIndex].Address
			// Skip the prevotes that the signer addressed as equivocated
			prevotesForR := fd.msgStore.GetPrevotes(height, func(m *message.Prevote) bool {
				return m.R() == prevote.R() && m.Signers().Contains(signerIndex) && m.Value() != prevote.Value()
			})
			if len(prevotesForR) > 0 {
				continue signersLoop
			}

			// We need to check whether we have proposals from the prevote's round
			correspondingProposals := fd.msgStore.GetProposals(height, func(m *message.Propose) bool {
				return m.R() == prevote.R() && m.Value() == prevote.Value()
			})

			if len(correspondingProposals) == 0 {
				// if there are over quorum prevotes for this corresponding proposal's value, then it indicates current
				// peer just did not receive it. So we can skip the rising of such accusation.
				if fd.msgStore.PrevotesPowerFor(height, prevote.R(), prevote.Value()).Cmp(quorum) < 0 {
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
							OffenderIndex: signerIndex,
						}
						proofs = append(proofs, accusation)
						fd.logger.Info("🕵️ Suspicious behavior detected", "rule", "PVN", "suspect", signer)
					}
				}
				continue signersLoop // we have no corresponding proposal, so we cannot check new and old prevote rules
			}

			// We need to ensure that we keep all proposals in the message store, so that we have the maximum chance of
			// finding justification for prevotes. This is to account for equivocation where the proposer send 2 proposals
			// with the same value but different valid rounds to different nodes. We can't penalise the signer of prevote
			// since we can't tell which proposal they received. We just want to find a set of message which fit the rule.
			// Therefore, we need to check all the proposals to find a single one which shows the current prevote is
			// valid.
			var prevotesProofs []*Proof
			for _, proposal := range correspondingProposals {
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
	precommitsFromPi := fd.msgStore.GetPrecommits(height, func(m *message.Precommit) bool {
		return m.R() < prevote.R() && m.Signers().Contains(signerIndex)
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
				// we found the latest non-nil precommit and we don't have gaps in the following ones
				pc := precommitsFromPi[i]

				// check for equivocation. If present, bail out on the checking of this rule. Remote peer has already been punished for equivocation
				precommitsAtRPrime := fd.msgStore.GetPrecommits(height, func(m *message.Precommit) bool {
					return m.R() == pc.R() && m.Signers().Contains(signerIndex) && m.Value() != pc.Value()
				})
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
					OffenderIndex: signerIndex,
				}
				// to guarantee this prevote is for a new proposal that is the PVN rule account for, otherwise in
				// prevote for an old proposal, it is valid for one to prevote it if lockedRound <= vr, thus the
				// round jump is valid. This prevents from rising a PVN misbehavior proof from a malicious fault
				// detector by using prevote for an old proposal to challenge an honest slow validator.
				proof.Evidences = append(proof.Evidences, message.NewLightProposal(correspondingProposal))
				precommits := precommitsFromPi[i:]

				// we won't do aggregation for a single prommit.
				if len(precommits) == 1 {
					proof.Evidences = append(proof.Evidences, message.Msg(precommits[0]))
					return proof
				}

				// we have multiple precommits, do the distinct msg aggregation
				proof.DistinctPrecommits = AggregateDistinctPrecommits(precommits)
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

	alternativeQuorum := fd.msgStore.SearchQuorum(height, validRound, correspondingProposal.Value(), quorum)
	if len(alternativeQuorum) > 0 {
		fd.logger.Info("Misbehaviour detected", "rule", "PV0", "incriminated", signer)
		// fast aggregate quorum prevotes into single one.
		evidences := make([]message.Msg, 1)
		evidences[0] = alternativeQuorum[0]
		if len(alternativeQuorum) > 1 {
			evidences[0] = AggregateSamePrevotes(alternativeQuorum)
		}

		proof := &Proof{
			Type:          autonity.Misbehaviour,
			Rule:          autonity.PVO,
			Message:       prevote,
			OffenderIndex: signerIndex,
		}
		proof.Evidences = append(proof.Evidences, message.NewLightProposal(correspondingProposal))
		proof.Evidences = append(proof.Evidences, evidences...)
		return proof
	}

	overQuorumPrevotesForVFromValidRound := fd.msgStore.PrevotesPowerFor(height, validRound, correspondingProposal.Value()).Cmp(quorum) >= 0

	if overQuorumPrevotesForVFromValidRound {
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

		precommitsFromPi := fd.msgStore.GetPrecommits(height, func(m *message.Precommit) bool {
			return m.R() > validRound && m.R() < currentR && m.Signers().Contains(signerIndex)
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
				fd.logger.Info("Misbehaviour detected", "rule", "PVO12", "incriminated", signer)
				proof := &Proof{
					Type:          autonity.Misbehaviour,
					Rule:          autonity.PVO12,
					Message:       prevote,
					OffenderIndex: signerIndex,
				}
				proof.Evidences = append(proof.Evidences, message.NewLightProposal(correspondingProposal))
				// we won't do aggregation for a single prommit.
				if len(precommitsFromPi) == 1 {
					proof.Evidences = append(proof.Evidences, message.Msg(precommitsFromPi[0]))
					return proof
				}
				// aggregate distinct precommits
				proof.DistinctPrecommits = AggregateDistinctPrecommits(precommitsFromPi)
				return proof
			}
		}
	}

	// if there is no misbehaviour of the prevote msg addressed, then we lastly check accusation.
	if !overQuorumPrevotesForVFromValidRound {
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
				OffenderIndex: signerIndex,
			}
		}
	}

	return nil
}

func (fd *FaultDetector) precommitsAccountabilityCheck(height uint64, quorum *big.Int, committee *types.Committee) (proofs []*Proof) {
	// ------------precommits------------
	// C: [Mr,P|proposer(r)] ∧ [Mr,PV] <--- [Mr,PC|pi]
	// C1: [V:Valid(V)] ∧ [#(V) ≥ 2f+ 1] <--- [V]

	precommits := fd.msgStore.GetPrecommits(height, func(m *message.Precommit) bool {
		return m.Value() != nilValue
	})

	for _, precommit := range precommits {
	signersLoop:
		for _, signerIndex := range precommit.Signers().FlattenUniq() {
			signer := committee.Members[signerIndex].Address

			// Skip if preCommit is equivocated
			precommitsForR := fd.msgStore.GetPrecommits(height, func(m *message.Precommit) bool {
				return m.R() == precommit.R() && m.Signers().Contains(signerIndex) && m.Value() != precommit.Value()
			})
			if len(precommitsForR) > 0 {
				continue signersLoop
			}

			// Do we see a quorum for a value other than the proposed value? If so, we have proof of misbehaviour.
			alternativeQuorum := fd.msgStore.SearchQuorum(height, precommit.R(), precommit.Value(), quorum)
			// Here the assumption is that in a single round it is not possible to have 2 value which quorum votes,
			// this would imply at least quorum nodes are malicious which is much higher than our assumption.
			if len(alternativeQuorum) > 0 {
				// fast aggregate quorum prevotes into single one.
				evidences := make([]message.Msg, 1)
				evidences[0] = alternativeQuorum[0]
				if len(alternativeQuorum) > 1 {
					evidences[0] = AggregateSamePrevotes(alternativeQuorum)
				}

				proof := &Proof{
					Type:          autonity.Misbehaviour,
					Rule:          autonity.C,
					Evidences:     evidences,
					Message:       precommit,
					OffenderIndex: signerIndex,
				}
				proofs = append(proofs, proof)
				fd.logger.Info("Misbehaviour detected", "rule", "C", "incriminated", signer)
				continue signersLoop
			}

			// Do we see a quorum of prevotes in the same round? if not we can raise an accusation, since we cannot be sure
			// that these prevotes do exist, this block also covers the Accusation of C since if over quorum prevotes for
			// V indicates that the corresponding proposal of V do exist, thus we don't need to raise accusation for the missing
			// proposal since over 2/3 member should all ready received it
			if fd.msgStore.PrevotesPowerFor(height, precommit.R(), precommit.Value()).Cmp(quorum) < 0 {
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
	offenderIndex int, offender common.Address) {
	rule := errorToRule(err)
	proof := fd.eventFromProof(&Proof{
		Type:          autonity.Misbehaviour,
		Rule:          rule,
		Message:       m,
		Evidences:     evidence,
		OffenderIndex: offenderIndex,
	}, offender)
	// submit misbehavior proof to buffer, it will be sent once aggregated.
	fd.misbehaviourProofCh <- proof
}

func (fd *FaultDetector) checkSelfIncriminatingProposal(proposal *message.Propose) error {
	// skip processing duplicated msg.
	// Cannot use .Hash() here because some fields of the block are not taken into account by block.Hash()
	// i.e. we can have two proposals from the same signer with:
	// - same height and round
	// - same value and validRound
	// BUT different payload hash
	duplicated := fd.msgStore.GetProposals(proposal.H(), func(p *message.Propose) bool {
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
		fd.submitMisbehavior(message.NewLightProposal(proposal), nil, errProposer, proposal.SignerIndex(), proposal.Signer())
		return errProposer
	}

	// account for equivocation
	equivocated := fd.msgStore.GetProposals(proposal.H(), func(p *message.Propose) bool {
		return p.R() == proposal.R() &&
			p.Signer() == proposal.Signer() &&
			(p.Value() != proposal.Value() || p.ValidRound() != proposal.ValidRound())
	})

	if len(equivocated) > 0 {
		var equivocatedMsgs = []message.Msg{
			message.NewLightProposal(equivocated[0]),
		}
		fd.submitMisbehavior(message.NewLightProposal(proposal), equivocatedMsgs, errEquivocation, proposal.SignerIndex(), proposal.Signer())
		// we allow the equivocated msg to be stored in msg store.
		fd.msgStore.Save(proposal)
		return errEquivocation
	}
	fd.msgStore.Save(proposal)
	return nil
}

func (fd *FaultDetector) checkSelfIncriminatingPrevote(m *message.Prevote) error {
	// skip process duplicated for votes.
	duplicatedMsg := fd.msgStore.GetPrevotes(m.H(), func(msg *message.Prevote) bool {
		return msg.Hash() == m.Hash()
	})

	if len(duplicatedMsg) > 0 {
		return errDuplicatedMsg
	}

	// account for equivocation for votes.
	var err error
	committee, err := fd.blockchain.CommitteeByHeight(m.H())
	if err != nil {
		panic(fmt.Sprintf("cannot find committee for height %d", m.H()))
	}

	for _, signerIndex := range m.Signers().FlattenUniq() {
		signer := committee.Members[signerIndex].Address
		equivocatedMessages := fd.msgStore.GetPrevotes(m.H(), func(msg *message.Prevote) bool {
			return msg.R() == m.R() && msg.Signers().Contains(signerIndex) && msg.Value() != m.Value()
		})
		if len(equivocatedMessages) > 0 {
			fd.submitMisbehavior(m, []message.Msg{equivocatedMessages[0]}, errEquivocation, signerIndex, signer)
			err = errEquivocation
		}
	}

	fd.msgStore.Save(m)
	return err
}

func (fd *FaultDetector) checkSelfIncriminatingPrecommit(m *message.Precommit) error {
	// skip process duplicated for votes.
	duplicatedMsg := fd.msgStore.GetPrecommits(m.H(), func(msg *message.Precommit) bool {
		return msg.Hash() == m.Hash()
	})

	if len(duplicatedMsg) > 0 {
		return errDuplicatedMsg
	}

	// account for equivocation for votes.
	var err error
	committee, err := fd.blockchain.CommitteeByHeight(m.H())
	if err != nil {
		panic(fmt.Sprintf("cannot get committee of height: %d", m.H()))
	}
	for _, signerIndex := range m.Signers().FlattenUniq() {
		signer := committee.Members[signerIndex].Address
		equivocatedMessages := fd.msgStore.GetPrecommits(m.H(), func(msg *message.Precommit) bool {
			return msg.R() == m.R() && msg.Signers().Contains(signerIndex) && msg.Value() != m.Value()
		})
		if len(equivocatedMessages) > 0 {
			fd.submitMisbehavior(m, []message.Msg{equivocatedMessages[0]}, errEquivocation, signerIndex, signer)
			err = errEquivocation
		}
	}

	fd.msgStore.Save(m)
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
		// these 2 errors are the only ones which can be raised by a self-incriminating msg.
		// if something else arrives here, it is a programming error.
		// there should also be 'InvalidProposal', however we do not currently make them accountable (due to oversized proof).
		panic("unknown error to accountability rule mapping")
	}

	return rule
}

func isProposerValid(chain ChainContext, m message.Msg) bool {
	committee, err := chain.CommitteeByHeight(m.H())
	if err != nil {
		panic(fmt.Sprintf("cannot get committee of height: %d", m.H()))
	}
	proposer := chain.ProtocolContracts().Proposer(committee, nil, m.H()-1, m.R())
	signer := m.(*message.Propose).Signer()
	return signer == proposer
}
