package faultdetector

import (
	"fmt"
	"github.com/clearmatics/autonity/autonity"
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus"
	"github.com/clearmatics/autonity/consensus/tendermint/bft"
	tendermintCore "github.com/clearmatics/autonity/consensus/tendermint/core"
	"github.com/clearmatics/autonity/consensus/tendermint/crypto"
	"github.com/clearmatics/autonity/consensus/tendermint/events"
	"github.com/clearmatics/autonity/core"
	"github.com/clearmatics/autonity/core/state"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/core/vm"
	"github.com/clearmatics/autonity/event"
	"github.com/clearmatics/autonity/log"
	"github.com/clearmatics/autonity/rlp"
	"github.com/syndtr/goleveldb/leveldb/errors"

	"math/rand"
	"sort"
	"sync"
	"time"
)

const (
	msgProposal uint8 = iota
	msgPrevote
	msgPrecommit
)

type Rule uint8

const (
	PN Rule = iota
	PO
	PVN
	PVO
	C
	C1

	GarbageMessage  // message was signed by valid member, but it cannot be decoded.
	InvalidProposal // The value proposed by proposer cannot pass the blockchain's validation.
	InvalidProposer // A proposal sent from none proposer nodes of the committee.
	Equivocation    // Multiple distinguish votes(proposal, prevote, precommit) sent by validator.
	UnknownRule
)

type BlockChainContext interface {
	consensus.ChainReader
	CurrentBlock() *types.Block
	SubscribeChainEvent(ch chan<- core.ChainEvent) event.Subscription
	State() (*state.StateDB, error)
	GetAutonityContract() *autonity.Contract
	StateAt(root common.Hash) (*state.StateDB, error)
	HasBadBlock(hash common.Hash) bool
	Validator() core.Validator
}

var (
	// todo: refine the window and buffer range in contract which can be tuned during run time.
	deltaToWaitForAccountability = 30                                // Wait until the GST + delta (30 blocks) to start rule scan.
	msgBufferInHeight            = deltaToWaitForAccountability + 60 // buffer such range of msgs in height at msg store.

	errEquivocation    = errors.New("equivocation happens")
	errFutureMsg       = errors.New("future height msg")
	errGarbageMsg      = errors.New("garbage msg")
	errNotCommitteeMsg = errors.New("msg from none committee member")
	errProposal        = errors.New("proposal have invalid values")
	errProposer        = errors.New("proposal is not from proposer")

	errNoEvidenceForPO  = errors.New("no evidence for innocence of rule PO")
	errNoEvidenceForPVN = errors.New("no evidence for innocence of rule PVN")
	errNoEvidenceForC   = errors.New("no evidence for innocence of rule C")
	errNoEvidenceForC1  = errors.New("no evidence for innocence of rule C1")

	nilValue          = common.Hash{}
	randomDelayWindow = 1000 * 5 // (0, 5] seconds random time window
)

// proof is what to prove that one is misbehaving, one should be slashed when a valid proof is rise.
type proof struct {
	Type     autonity.ProofType // Misbehaviour, Accusation, Innocence.
	Rule     Rule
	Message  *tendermintCore.Message   // the msg to be considered as suspicious or misbehaved one
	Evidence []*tendermintCore.Message // the proofs of innocence or misbehaviour.
}

// event to submit proofs via standard transaction.
type AccountabilityEvent struct {
	OnChainProofs []*autonity.OnChainProof
}

// FaultDetector it subscribe chain event to trigger rule engine to apply patterns over
// msg store, it send proof of challenge if it detects any potential misbehavior, either it
// read state db on each new height to get latest challenges from autonity contract's view,
// and to prove its innocent if there were any challenges on the suspicious node.
type FaultDetector struct {
	sync.RWMutex

	proofWG           sync.WaitGroup
	faultDetectorFeed event.Feed

	tendermintMsgSub *event.TypeMuxSubscription
	blockChan        chan core.ChainEvent
	blockSub         event.Subscription

	blockchain BlockChainContext

	address         common.Address
	msgStore        *MsgStore
	futureHeightMsg map[uint64][]*tendermintCore.Message // map[blockHeight][]*tendermintMessages

	onChainProofsBuffer []*autonity.OnChainProof // buffer proofs to aggregate them into single TX.

	logger log.Logger
}

// call by ethereum object to create fd instance.
func NewFaultDetector(chain BlockChainContext, nodeAddress common.Address, sub *event.TypeMuxSubscription) *FaultDetector {
	fd := &FaultDetector{
		RWMutex:          sync.RWMutex{},
		address:          nodeAddress,
		blockChan:        make(chan core.ChainEvent, 300),
		blockchain:       chain,
		msgStore:         newMsgStore(),
		logger:           log.New("FaultDetector", nodeAddress),
		tendermintMsgSub: sub,
		futureHeightMsg:  make(map[uint64][]*tendermintCore.Message),
	}

	// register faultdetector contracts on evm's precompiled contract set.
	registerFaultDetectorContracts(chain)
	return fd
}

// listen for new block events from block-chain, do the tasks like take challenge and provide proof for innocent, the
// Fault Detector rule engine could also triggered from here to scan those msgs of msg store by applying rules.
func (fd *FaultDetector) FaultDetectorEventLoop() {
	go fd.blockEventLoop()
	go fd.tendermintMsgEventLoop()
}

func (fd *FaultDetector) tendermintMsgEventLoop() {
	for {
		ev, ok := <-fd.tendermintMsgSub.Chan()
		if !ok {
			return
		}

		mv, ok := ev.Data.(events.MessageEvent)
		if !ok {
			fd.logger.Crit("programming error", "cannot cast message event to events.MessageEvent instead received ", ev.Data)
			return
		}

		msg := new(tendermintCore.Message)
		if err := msg.FromPayload(mv.Payload); err != nil {
			fd.logger.Error("invalid payload", "faultdetector", err)
			continue
		}

		// discard too old messages which is out of accountability buffering window.
		head := fd.blockchain.CurrentHeader().Number.Uint64()
		if head > uint64(msgBufferInHeight) && msg.H() < head-uint64(msgBufferInHeight) {
			fd.logger.Info("discard too old message for accountability", "faultdetector", msg.Sender())
			continue
		}

		if err := fd.processMsg(msg); err != nil {
			fd.logger.Warn("process consensus msg", "faultdetector", err)
			continue
		}
	}
}

func (fd *FaultDetector) blockEventLoop() {
	fd.blockSub = fd.blockchain.SubscribeChainEvent(fd.blockChan)

	for {
		select {
		// chain event update, provide proof of innocent if one is on challenge, rule engine scanning is triggered also.
		case ev := <-fd.blockChan:
			// before run rule engine over msg store, process any buffered msg.
			fd.processBufferedMsgs(ev.Block.NumberU64())

			// handle accusations and provide innocence proof if there were any for a node.
			innocenceProofs, _ := fd.handleAccusations(ev.Block, ev.Block.Root())
			if innocenceProofs != nil {
				fd.Lock()
				fd.onChainProofsBuffer = append(fd.onChainProofsBuffer, innocenceProofs...)
				fd.Unlock()
			}

			// run rule engine over a specific height.
			proofs := fd.runRuleEngine(ev.Block.NumberU64())
			if len(proofs) > 0 {
				fd.Lock()
				fd.onChainProofsBuffer = append(fd.onChainProofsBuffer, proofs...)
				fd.Unlock()
			}

			// aggregate buffered proofs into single TX and send.
			fd.sentProofs()

			// msg store delete msgs out of buffering window.
			fd.msgStore.DeleteMsgsAtHeight(ev.Block.NumberU64() - uint64(msgBufferInHeight))
		case err, ok := <-fd.blockSub.Err():
			if ok {
				fd.logger.Crit("block subscription error", err.Error())
			}
			return
		}
	}
}

func (fd *FaultDetector) Stop() {
	fd.blockSub.Unsubscribe()
	fd.tendermintMsgSub.Unsubscribe()
	fd.proofWG.Wait()
	unRegisterFaultDetectorContracts()
}

// call by ethereum object to subscribe proofs Events.
func (fd *FaultDetector) SubscribeFaultDetectorEvents(ch chan<- AccountabilityEvent) event.Subscription {
	return fd.faultDetectorFeed.Subscribe(ch)
}

// buffer Msg since local chain may not synced yet to verify if msg is from correct committee.
func (fd *FaultDetector) bufferMsg(m *tendermintCore.Message) {
	h, err := m.Height()
	if err != nil {
		return
	}

	fd.Lock()
	defer fd.Unlock()
	fd.futureHeightMsg[h.Uint64()] = append(fd.futureHeightMsg[h.Uint64()], m)
}

func (fd *FaultDetector) filterPresentedOnes(proofs []*autonity.OnChainProof) []*autonity.OnChainProof {
	// get latest chain state.
	var result []*autonity.OnChainProof
	state, err := fd.blockchain.State()
	if err != nil {
		return nil
	}
	header := fd.blockchain.CurrentBlock().Header()

	presentedAccusation := fd.blockchain.GetAutonityContract().GetAccusations(header, state)
	presentedMisbehavior := fd.blockchain.GetAutonityContract().GetMisBehaviours(header, state)

	for i := 0; i < len(proofs); i++ {
		present := false
		for j := 0; j < len(presentedAccusation); j++ {
			if proofs[i].Msghash == presentedAccusation[j].Msghash &&
				proofs[i].Type == autonity.Accusation {
				present = true
			}
		}

		for j := 0; j < len(presentedMisbehavior); j++ {
			if proofs[i].Msghash == presentedMisbehavior[j].Msghash &&
				proofs[i].Type == autonity.Misbehaviour {
				present = true
			}
		}

		if !present {
			result = append(result, proofs[i])
		}
	}

	return result
}

// convert the raw proofs into on-chain proof which contains raw bytes of messages.
func (fd *FaultDetector) generateOnChainProof(p *proof) (*autonity.OnChainProof, error) {
	var onChainProof = &autonity.OnChainProof{
		Type:    p.Type,
		Sender:  p.Message.Address,
		Msghash: types.RLPHash(p.Message),
	}

	rproof, err := rlp.EncodeToBytes(p)
	if err != nil {
		return nil, err
	}
	onChainProof.Rawproof = rproof
	return onChainProof, nil
}

// getInnocentProof called by client who is on challenge to get proof of innocent from msg store.
func (fd *FaultDetector) getInnocentProof(c *proof) (*autonity.OnChainProof, error) {
	var onChainProof *autonity.OnChainProof
	// rule engine have below provable accusation for the time being:
	switch c.Rule {
	case PO:
		return fd.getInnocentProofOfPO(c)
	case PVN:
		return fd.getInnocentProofOfPVN(c)
	case C:
		return fd.getInnocentProofOfC(c)
	case C1:
		return fd.getInnocentProofOfC1(c)
	default:
		return onChainProof, fmt.Errorf("not provable rule")
	}
}

// get proof of innocent of C from msg store.
func (fd *FaultDetector) getInnocentProofOfC(c *proof) (*autonity.OnChainProof, error) {
	var onChainProof *autonity.OnChainProof
	preCommit := c.Message
	height := preCommit.H()

	proposals := fd.msgStore.Get(height, func(m *tendermintCore.Message) bool {
		return m.Type() == msgProposal && m.Value() == preCommit.Value() && m.R() == preCommit.R()
	})

	if len(proposals) == 0 {
		// cannot onChainProof its innocent for PVN, the on-chain contract will fine it latter once the
		// time window for onChainProof ends.
		return onChainProof, errNoEvidenceForC
	}
	p, err := fd.generateOnChainProof(&proof{
		Type:     autonity.Innocence,
		Rule:     c.Rule,
		Message:  preCommit,
		Evidence: proposals,
	})
	if err != nil {
		return p, err
	}
	return p, nil
}

// get proof of innocent of C1 from msg store.
func (fd *FaultDetector) getInnocentProofOfC1(c *proof) (*autonity.OnChainProof, error) {
	var onChainProof *autonity.OnChainProof
	preCommit := c.Message
	height := preCommit.H()
	quorum := bft.Quorum(fd.blockchain.GetHeaderByNumber(height - 1).TotalVotingPower())

	prevotesForV := fd.msgStore.Get(height, func(m *tendermintCore.Message) bool {
		return m.Type() == msgPrevote && m.Value() == preCommit.Value() && m.R() == preCommit.R()
	})

	if powerOfVotes(prevotesForV) < quorum {
		// cannot onChainProof its innocent for PO for now, the on-chain contract will fine it latter once the
		// time window for onChainProof ends.
		return onChainProof, errNoEvidenceForC1
	}

	p, err := fd.generateOnChainProof(&proof{
		Type:     autonity.Innocence,
		Rule:     c.Rule,
		Message:  preCommit,
		Evidence: prevotesForV,
	})
	if err != nil {
		return p, err
	}

	return p, nil
}

// get proof of innocent of PO from msg store.
func (fd *FaultDetector) getInnocentProofOfPO(c *proof) (*autonity.OnChainProof, error) {
	// PO: node propose an old value with an validRound, innocent onChainProof of it should be:
	// there are quorum num of prevote for that value at the validRound.
	var onChainProof *autonity.OnChainProof
	proposal := c.Message
	height := proposal.H()
	validRound := proposal.ValidRound()
	quorum := bft.Quorum(fd.blockchain.GetHeaderByNumber(height - 1).TotalVotingPower())

	prevotes := fd.msgStore.Get(height, func(m *tendermintCore.Message) bool {
		return m.Type() == msgPrevote && m.R() == validRound && m.Value() == proposal.Value()
	})

	if powerOfVotes(prevotes) < quorum {
		// cannot onChainProof its innocent for PO, the on-chain contract will fine it latter once the
		// time window for onChainProof ends.
		return onChainProof, errNoEvidenceForPO
	}

	p, err := fd.generateOnChainProof(&proof{
		Type:     autonity.Innocence,
		Rule:     c.Rule,
		Message:  proposal,
		Evidence: prevotes,
	})
	if err != nil {
		return p, err
	}

	return p, nil
}

// get proof of innocent of PVN from msg store.
func (fd *FaultDetector) getInnocentProofOfPVN(c *proof) (*autonity.OnChainProof, error) {
	// get innocent proofs for PVN, for a prevote that vote for a new value,
	// then there must be a proposal for this new value.
	var onChainProof *autonity.OnChainProof
	prevote := c.Message
	height := prevote.H()

	correspondingProposals := fd.msgStore.Get(height, func(m *tendermintCore.Message) bool {
		return m.Type() == msgProposal && m.Value() == prevote.Value() && m.R() == prevote.R()
	})

	if len(correspondingProposals) == 0 {
		// cannot onChainProof its innocent for PVN, the on-chain contract will fine it latter once the
		// time window for onChainProof ends.
		return onChainProof, errNoEvidenceForPVN
	}

	p, err := fd.generateOnChainProof(&proof{
		Type:     autonity.Innocence,
		Rule:     c.Rule,
		Message:  prevote,
		Evidence: correspondingProposals,
	})
	if err != nil {
		return p, nil
	}

	return p, nil
}

// get accusations from chain via autonityContract calls, and provide innocent proofs if there were any challenge on node.
func (fd *FaultDetector) handleAccusations(block *types.Block, hash common.Hash) ([]*autonity.OnChainProof, error) {
	var innocentOnChainProofs []*autonity.OnChainProof
	state, err := fd.blockchain.StateAt(hash)
	if err != nil || state == nil {
		fd.logger.Error("handleAccusation", "faultdetector", err)
		return nil, err
	}

	contract := fd.blockchain.GetAutonityContract()
	if contract == nil {
		return nil, fmt.Errorf("cannot get contract instance")
	}

	accusations := contract.GetAccusations(block.Header(), state)
	for i := 0; i < len(accusations); i++ {
		if accusations[i].Sender == fd.address {
			c, err := decodeRawProof(accusations[i].Rawproof)
			if err != nil {
				continue
			}

			p, err := fd.getInnocentProof(c)
			if err != nil {
				continue
			}
			innocentOnChainProofs = append(innocentOnChainProofs, p)
		}
	}

	return innocentOnChainProofs, nil
}

// processBufferedMsgs, called on chain event update, it process msgs from the latest height buffered before.
func (fd *FaultDetector) processBufferedMsgs(height uint64) {
	fd.RLock()
	defer fd.RUnlock()
	for h, msgs := range fd.futureHeightMsg {
		if h <= height {
			for i := 0; i < len(msgs); i++ {
				if err := fd.processMsg(msgs[i]); err != nil {
					fd.logger.Warn("process consensus msg", "faultdetector", err)
					continue
				}
			}
			// once message processed, release it from buffer.
			delete(fd.futureHeightMsg, h)
		}
	}
}

// processMsg, check and submit any auto-incriminating, equivocation challenges, and then only store checked msg into msg store.
func (fd *FaultDetector) processMsg(m *tendermintCore.Message) error {
	// pre-check if msg is from valid committee member
	err := checkMsgSignature(fd.blockchain, m)
	if err != nil {
		if err == errFutureMsg {
			fd.bufferMsg(m)
		}
		return err
	}

	// decode consensus msg, and auto-incriminating msg is addressed here.
	err = checkAutoIncriminatingMsg(fd.blockchain, m)
	if err != nil {
		if err == errFutureMsg {
			fd.bufferMsg(m)
		} else {
			proofs := []*tendermintCore.Message{m}
			fd.submitMisbehavior(m, proofs, err)
			return err
		}
	}

	// store msg, if there is equivocation, msg store would then rise errEquivocation and proofs.
	msgs, err := fd.msgStore.Save(m)
	if err == errEquivocation && len(msgs) > 0 {
		var proofs []*tendermintCore.Message
		for i := 0; i < len(msgs); i++ {
			proofs = append(proofs, msgs[i])
		}
		fd.submitMisbehavior(m, proofs, err)
		return err
	}
	return nil
}

// run rule engine over latest msg store, if the return proofs is not empty, then rise challenge.
func (fd *FaultDetector) runRuleEngine(height uint64) []*autonity.OnChainProof {
	var onChainProofs []*autonity.OnChainProof
	// To avoid none necessary accusations, we wait for delta blocks to start rule scan.
	if height > uint64(deltaToWaitForAccountability) {
		// run rule engine over the previous delta offset height.
		checkPointHeight := height - uint64(deltaToWaitForAccountability)
		quorum := bft.Quorum(fd.blockchain.GetHeaderByNumber(checkPointHeight - 1).TotalVotingPower())
		proofs := fd.runRulesOverHeight(checkPointHeight, quorum)
		if len(proofs) > 0 {
			for i := 0; i < len(proofs); i++ {
				p, err := fd.generateOnChainProof(proofs[i])
				if err != nil {
					fd.logger.Warn("convert proof to on-chain proof", "faultdetector", err)
					continue
				}
				onChainProofs = append(onChainProofs, p)
			}
		}
	}
	return onChainProofs
}

func (fd *FaultDetector) runRulesOverHeight(height uint64, quorum uint64) (proofs []*proof) {
	// Rules read right to left (find  the right and look for the left)
	//
	// Rules should be evealuated such that we check all paossible instances and if
	// we can't find a single instance that passes then we consider the rule
	// failed.
	//
	// There are 2 types of provable misbehaviors.
	// 1. Conflicting messages from a single participant
	// 2. A message that conflicts with a quorum of prevotes.
	// (precommit for differing value in same round as the prevotes or proposal for an
	// old value where in each prior round we can see a quorum of precommits for a distinct value.)

	// We should be here at time t = timestamp(h+1) + delta

	// ------------New Proposal------------
	// PN:  (Mr′<r,P C|pi)∗ <--- (Mr,P|pi)
	// PN1: [nil ∨ ⊥] <--- [V]

	proposalsNew := fd.msgStore.Get(height, func(m *tendermintCore.Message) bool {
		return m.Type() == msgProposal && m.ValidRound() == -1
	})

	for _, proposal := range proposalsNew {
		//check all precommits for previous rounds from this sender are nil
		precommits := fd.msgStore.Get(height, func(m *tendermintCore.Message) bool {
			return m.Sender() == proposal.Sender() && m.Type() == msgPrecommit && m.R() < proposal.R() && m.Value() != nilValue // nolint: scopelint
		})
		if len(precommits) != 0 {
			proof := &proof{
				Type:     autonity.Misbehaviour,
				Rule:     PN,
				Evidence: precommits,
				Message:  proposal,
			}
			proofs = append(proofs, proof)
		}
	}

	// ------------Old Proposal------------
	// PO: (Mr′<r,PV) ∧ (Mr′,PC|pi) ∧ (Mr′<r′′<r,P C|pi)∗ <--- (Mr,P|pi)
	// PO1: [#(Mr′,PV|V) ≥ 2f+ 1] ∧ [nil ∨ V ∨ ⊥] ∧ [nil ∨ ⊥] <--- [V]

	proposalsOld := fd.msgStore.Get(height, func(m *tendermintCore.Message) bool {
		return m.Type() == msgProposal && m.ValidRound() > -1
	})

	for _, proposal := range proposalsOld {
		// Check that in the valid round we see a quorum of prevotes and that
		// there is no precommit at all or a precommit for v or nil.

		validRound := proposal.ValidRound()

		// Is there a precommit for a value other than nil or the proposed value
		// by the current proposer in the valid round? If there is the proposer
		// has proposed a value for which it is not locked on, thus a proof of
		// misbehaviour can be generated.
		precommits := fd.msgStore.Get(height, func(m *tendermintCore.Message) bool {
			return m.Type() == msgPrecommit && m.R() == validRound &&
				m.Sender() == proposal.Sender() && m.Value() != nilValue && m.Value() != proposal.Value() // nolint: scopelint
		})
		if len(precommits) > 0 {
			proof := &proof{
				Type:     autonity.Misbehaviour,
				Rule:     PO,
				Evidence: precommits,
				Message:  proposal,
			}
			proofs = append(proofs, proof)
		}

		// Is there a precommit for anything other than nil from the proposer
		// between the valid round and the round of the proposal? If there is
		// then that implies the proposer saw 2f+1 prevotes in that round and
		// hence it should have set that round as the valid round.
		precommits = fd.msgStore.Get(height, func(m *tendermintCore.Message) bool {
			return m.Type() == msgPrecommit &&
				m.R() > validRound && m.R() < proposal.R() && m.Sender() == proposal.Sender() && m.Value() != nilValue // nolint: scopelint
		})
		if len(precommits) > 0 {
			proof := &proof{
				Type:     autonity.Misbehaviour,
				Rule:     PO,
				Evidence: precommits,
				Message:  proposal,
			}
			proofs = append(proofs, proof)
		}

		// Do we see a quorum of prevotes in the valid round, if not we can
		// raise an accusation, since we cannot be sure that these prevotes
		// don't exist
		prevotes := fd.msgStore.Get(height, func(m *tendermintCore.Message) bool {
			// since equivocation msgs are stored, we have to query those preVotes which has same value as the proposal.
			return m.Type() == msgPrevote && m.R() == validRound && m.Value() == proposal.Value() // nolints: scopelint
		})

		if powerOfVotes(deEquivocatedMsgs(prevotes)) < quorum {
			accusation := &proof{
				Type:    autonity.Accusation,
				Rule:    PO,
				Message: proposal,
			}
			proofs = append(proofs, accusation)
		}
	}

	// ------------New and Old Prevotes------------

	prevotes := fd.msgStore.Get(height, func(m *tendermintCore.Message) bool {
		return m.Type() == msgPrevote && m.Value() != nilValue
	})

	for _, prevote := range prevotes {
		correspondingProposals := fd.msgStore.Get(height, func(m *tendermintCore.Message) bool {
			return m.Type() == msgProposal && m.Value() == prevote.Value() && m.R() == prevote.R() // nolint: scopelint
		})

		if len(correspondingProposals) == 0 {
			accusation := &proof{
				Type: autonity.Accusation,
				Rule: PVN, //This could be PVO as well, however, we can't decide since there are no corresponding
				// proposal
				Message: prevote,
			}
			proofs = append(proofs, accusation)
		}

		// We need to ensure that we keep all proposals in the message store, so that we have the maximum chance of
		// finding justification for prevotes. This is to account for equivocation where the proposer send 2 proposals
		// with the same value but different valid rounds to different nodes. We can't penalise the sender of prevote
		// since we can't tell which proposal they received. We just want to find a set of message which fit the rule.

		for _, correspondingProposal := range correspondingProposals {
			if correspondingProposal.ValidRound() == -1 {
				// New Proposal, apply PVN rules

				// PVN: (Mr′<r,PC|pi)∧(Mr′<r′′<r,PC|pi)* ∧ (Mr,P|proposer(r)) <--- (Mr,PV|pi)

				// PVN2: [nil ∨ ⊥] ∧ [nil ∨ ⊥] ∧ [V:Valid(V)] <--- [V]: r′= 0,∀r′′< r:Mr′′,PC|pi=nil

				// PVN2, If there is a valid proposal V at round r, and pi never
				// ever precommit(locked a value) before, then pi should prevote
				// for V or a nil in case of timeout at this round.

				// PVN3: [V] ∧ [nil ∨ ⊥] ∧ [V:Valid(V)] <--- [V]:∀r′< r′′<r,Mr′′,PC|pi=nil

				// We can check both PVN2 and PVN3 by simply searching for a
				// precommit for a value other than V or nil. This is a proof of
				// misbehaviour. There is no scope to raise an accusation for
				// these rules since the only message in PVN that is not sent by
				// pi is the proposal and you require the proposal before you
				// can even attempt to apply the rule.
				precommits := fd.msgStore.Get(height, func(m *tendermintCore.Message) bool {
					return m.Type() == msgPrecommit && m.Value() != nilValue &&
						m.Value() != prevote.Value() && prevote.Sender() == m.Sender() && m.R() < prevote.R() // nolint: scopelint
				})

				if len(precommits) > 0 {
					proof := &proof{
						Type:     autonity.Misbehaviour,
						Rule:     PVN,
						Evidence: precommits,
						Message:  prevote,
					}
					proofs = append(proofs, proof)
					break
				}

			} /*else {
				todo: missing PVO rules from D3
				// PVO:   (Mr′<r,PC|pi) ∧ (Mr′≤r′′′<r,PV) ∧ (Mr′<r′′<r,PC|pi)* ∧ (Mr,P|proposer(r)) <--- (Mr,PV|pi)

				// PVO1A: [V] ∧ [∗] ∧ [nil v ⊥] ∧ [V] <--- [V]:∀r′<r′′<r,Mr′′,PC|pi=nil <-- broken we need to see the prevotes for valid round

				// PVO2: [*] ∧ [#(V) ≥ 2f+1] ∧ [nil v ⊥] ∧ [V:validRound(V)=r′′′] <--- [V]:∀r′<r′′<r,Mr′′,PC|pi=nil ∧ ∃r′′′∈[r′,r−1],#(Mr′′′,PV|V) ≥ 2f+ 1

				// If pi previously precommitted for V and between this precommit and
				// the proposal precommitted for a different value V', then the prevote
				// is considered invalid.

				precommits := fd.msgStore.Get(height, func(m *tendermintCore.Message) bool {
					uint8()(return m.Type()) == msgPrecommit && prevote.Sender() == m.Sender() &&
						m.R() < prevote.R() && m.Value() != nilValue
				})
				//check most recent precommit if == V -> pass else --> fail

				// 2f+1 PV(V) round 2

				// round 4 p_i receiveds 2f+1 PV(V') Sends PC(V') and it sets its locked value and locked round=4

				// round 5 proposer proposes P(V, VR=2), so this would mean that p_i prevote nil even though there are 2f+1 prevotes for V in round 2

				// Aneeque's initials thoughts on PVO
				if len(precommits) > 0 {
					// PVO1a

					// sort according to round
					//sort.Sort(precommits)

					// Proof of misbehaviour:

					// Get the lastest precommit
					// Check the precommit value
					// if it precommit.Value() != prevote.Value
					// 		check all round from precommit to current round for 2f+1 prevotes
					// 		if even a single round doesn't have 2f+1 prevotes, raise an accusation
					//		else we have proof of misbehaviour if non of the 2f+1 prevotes are for precommit.Value()

					// if it precommit.Value() == prevote.Value
					// 		Check that if we 2f+1 prevotes for all rounds since precommit.Round() till current round,
					//      if yes, than non of them can be for value other than prevote.Value, otherwise we have proof of misbehaviour
					// 		if there are gaps then the condition passes

				} else {
					// PVO2

					// We don't have a precommit from the p_i
					// check that in valid round we have 2f+1 prevotes for V rule passes, otherwise raise an accustion
				}

				// PVO1B: [∗] ∧ [∗] ∧ [V:r′′=r−1] ∧ [V] <--- [V] -- not needed as it is a special case of PVO1A

				// PVO2: [*] ∧ [#(V) ≥ 2f+1] ∧ [nil v ⊥] ∧ [V:validRound(V)=r′′′] <--- [V]:∀r′<r′′<r,Mr′′,PC|pi=nil ∧ ∃r′′′∈[r′,r−1],#(Mr′′′,PV|V) ≥ 2f+ 1
				// If we can see an old proposal for V with valid round vr and
				// 2f+1 prevotes for the V in round vr, then pi could have also
				// seen them and hence be able to prevote for the old proposal.
			} */
		}
	}

	// ------------Precommits------------
	// C: [Mr,P|proposer(r)] ∧ [Mr,PV] <--- [Mr,PC|pi]
	// C1: [V:Valid(V)] ∧ [#(V) ≥ 2f+ 1] <--- [V]

	precommits := fd.msgStore.Get(height, func(m *tendermintCore.Message) bool {
		return m.Type() == msgPrecommit && m.Value() != nilValue
	})

	for _, precommit := range precommits {
		proposals := fd.msgStore.Get(height, func(m *tendermintCore.Message) bool {
			return m.Type() == msgProposal && m.Value() == precommit.Value() && m.R() == precommit.R() // nolint: scopelint
		})

		if len(proposals) == 0 {
			accusation := &proof{
				Type:    autonity.Accusation,
				Rule:    C,
				Message: precommit,
			}
			proofs = append(proofs, accusation)
			continue
		}

		prevotesForNotV := fd.msgStore.Get(height, func(m *tendermintCore.Message) bool {
			return m.Type() == msgPrevote && m.Value() != precommit.Value() && m.R() == precommit.R() // nolint: scopelint
		})
		prevotesForV := fd.msgStore.Get(height, func(m *tendermintCore.Message) bool {
			return m.Type() == msgPrevote && m.Value() == precommit.Value() && m.R() == precommit.R() // nolint: scopelint
		})

		// even if we have equivocated preVotes for not V, we still assume that there are less f+1 malicious node in the
		// network, so the powerOfVotes of preVotesForNotV which was deEquivocated is still valid to prove that the
		// preCommit is a misbehaviour of rule C.
		deEquivocatedPreVotesForNotV := deEquivocatedMsgs(prevotesForNotV)
		if powerOfVotes(deEquivocatedPreVotesForNotV) >= quorum {
			// In this case there cannot be enough remaining prevotes
			// to justify a precommit for V.
			proof := &proof{
				Type:     autonity.Misbehaviour,
				Rule:     C,
				Evidence: deEquivocatedPreVotesForNotV,
				Message:  precommit,
			}
			proofs = append(proofs, proof)

		} else if powerOfVotes(prevotesForV) < quorum {
			// In this case we simply don't see enough prevotes to
			// justify the precommit.
			accusation := &proof{
				Type:    autonity.Accusation,
				Rule:    C1,
				Message: precommit,
			}
			proofs = append(proofs, accusation)
		}
	}

	return proofs
}

// send proofs via event which will handled by ethereum object to signed the TX to send proof.
func (fd *FaultDetector) sendProofs(proofs []*autonity.OnChainProof) {
	fd.proofWG.Add(1)
	go func() {
		defer fd.proofWG.Done()
		randomDelay()
		unPresented := fd.filterPresentedOnes(proofs)
		if len(unPresented) != 0 {
			fd.faultDetectorFeed.Send(AccountabilityEvent{OnChainProofs: unPresented})
		}
	}()
}

func (fd *FaultDetector) sentProofs() {
	fd.Lock()
	defer fd.Unlock()

	if len(fd.onChainProofsBuffer) != 0 {
		copyOnChainProofs := make([]*autonity.OnChainProof, len(fd.onChainProofsBuffer))
		copy(copyOnChainProofs, fd.onChainProofsBuffer)
		fd.sendProofs(copyOnChainProofs)
		// release items from buffer
		fd.onChainProofsBuffer = fd.onChainProofsBuffer[:0]
	}
}

// submitMisbehavior takes proofs of misbehavior msg, and error id to form the on-chain proof, and
// send the proof of misbehavior to event channel.
func (fd *FaultDetector) submitMisbehavior(m *tendermintCore.Message, proofs []*tendermintCore.Message, err error) {
	rule, e := errorToRule(err)
	if e != nil {
		fd.logger.Warn("error to rule", "faultdetector", e)
	}
	proof, err := fd.generateOnChainProof(&proof{
		Type:     autonity.Misbehaviour,
		Rule:     rule,
		Message:  m,
		Evidence: proofs,
	})
	if err != nil {
		fd.logger.Warn("generate misbehavior proof", "faultdetector", err)
		return
	}

	// submit misbehavior proof to buffer, it will be sent once aggregated.
	fd.Lock()
	defer fd.Unlock()
	fd.onChainProofsBuffer = append(fd.onChainProofsBuffer, proof)
}

/////// common helper functions shared between faultdetector and precompiled contract to validate msgs.

// decode consensus msgs, address garbage msg and invalid proposal by returning error.
func checkAutoIncriminatingMsg(chain BlockChainContext, m *tendermintCore.Message) error {
	if m.Code == msgProposal {
		return checkProposal(chain, m)
	}

	if m.Code == msgPrevote || m.Code == msgPrecommit {
		return decodeVote(m)
	}

	return errors.New("unknown consensus msg")
}

func checkEquivocation(chain BlockChainContext, m *tendermintCore.Message, proof []*tendermintCore.Message) error {
	// decode msgs
	err := checkAutoIncriminatingMsg(chain, m)
	if err != nil {
		return err
	}

	for i := 0; i < len(proof); i++ {
		err := checkAutoIncriminatingMsg(chain, proof[i])
		if err != nil {
			return err
		}
	}
	// check equivocations.
	if !sameVote(m, proof[0]) {
		return errEquivocation
	}
	return nil
}

//checkMsgSignature, it check if msg is from valid member of the committee.
func checkMsgSignature(chain BlockChainContext, m *tendermintCore.Message) error {
	msgHeight, err := m.Height()
	if err != nil {
		return err
	}

	header := chain.CurrentHeader()
	if msgHeight.Uint64() > header.Number.Uint64()+1 {
		return errFutureMsg
	}

	lastHeader := chain.GetHeaderByNumber(msgHeight.Uint64() - 1)
	if lastHeader == nil {
		return errFutureMsg
	}

	if _, err = m.Validate(crypto.CheckValidatorSignature, lastHeader); err != nil {
		return errNotCommitteeMsg
	}
	return nil
}

// checkProposal, checks if proposal is valid and it's from correct proposer.
func checkProposal(chain BlockChainContext, m *tendermintCore.Message) error {
	var proposal tendermintCore.Proposal
	err := m.Decode(&proposal)
	if err != nil {
		return errGarbageMsg
	}
	if !isProposerMsg(chain, m) {
		return errProposer
	}

	err = verifyProposal(chain, *proposal.ProposalBlock)
	// due to timing issue, when Fault Detector validate a proposal, that proposal could already be
	// committed on the chain view. Since the msg sender were checked as the correct proposer, so we
	// consider this proposal as a valid proposal.
	if err == core.ErrKnownBlock {
		return nil
	}

	if err == consensus.ErrFutureBlock {
		return errFutureMsg
	}

	if err != nil {
		return errProposal
	}

	return nil
}

func decodeVote(m *tendermintCore.Message) error {
	var vote tendermintCore.Vote
	err := m.Decode(&vote)
	if err != nil {
		return errGarbageMsg
	}
	return nil
}

func deEquivocatedMsgs(msgs []*tendermintCore.Message) (deEquivocated []*tendermintCore.Message) {
	presented := make(map[common.Address]struct{})
	for _, v := range msgs {
		if _, ok := presented[v.Address]; ok {
			continue
		}
		deEquivocated = append(deEquivocated, v)
		presented[v.Address] = struct{}{}
	}
	return deEquivocated
}

func errorToRule(err error) (Rule, error) {
	rule := UnknownRule
	switch err {
	case errEquivocation:
		rule = Equivocation
	case errProposer:
		rule = InvalidProposer
	case errProposal:
		rule = InvalidProposal
	case errGarbageMsg:
		rule = GarbageMessage
	default:
		return rule, fmt.Errorf("errors of not provable")
	}

	return rule, nil
}

func getProposer(chain BlockChainContext, h uint64, r int64) (common.Address, error) {
	parentHeader := chain.GetHeaderByNumber(h - 1)
	if parentHeader.IsGenesis() {
		sort.Sort(parentHeader.Committee)
		return parentHeader.Committee[r%int64(len(parentHeader.Committee))].Address, nil
	}

	statedb, err := chain.StateAt(parentHeader.Root)
	if err != nil {
		return common.Address{}, err
	}

	proposer := chain.GetAutonityContract().GetProposerFromAC(parentHeader, statedb, parentHeader.Number.Uint64(), r)
	member := parentHeader.CommitteeMember(proposer)
	if member == nil {
		return common.Address{}, fmt.Errorf("cannot find correct proposer")
	}
	return proposer, nil
}

func isProposerMsg(chain BlockChainContext, m *tendermintCore.Message) bool {
	h, _ := m.Height()
	r, _ := m.Round()

	proposer, err := getProposer(chain, h.Uint64(), r)
	if err != nil {
		return false
	}

	return m.Address == proposer
}

func powerOfVotes(votes []*tendermintCore.Message) uint64 {
	counted := make(map[common.Address]struct{})
	power := uint64(0)
	for i := 0; i < len(votes); i++ {
		if votes[i].Type() == msgProposal {
			continue
		}

		if _, ok := counted[votes[i].Address]; ok {
			continue
		}

		power += votes[i].GetPower()
		counted[votes[i].Address] = struct{}{}
	}
	return power
}

func randomDelay() {
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(randomDelayWindow)
	time.Sleep(time.Duration(n) * time.Millisecond)
}

func sameVote(a *tendermintCore.Message, b *tendermintCore.Message) bool {
	ah, _ := a.Height()
	ar, _ := a.Round()
	bh, _ := b.Height()
	br, _ := b.Round()
	aHash := types.RLPHash(a)
	bHash := types.RLPHash(b)

	if ah == bh && ar == br && a.Code == b.Code && a.Address == b.Address && aHash == bHash {
		return true
	}
	return false
}

func verifyProposal(chain BlockChainContext, proposal types.Block) error {
	block := &proposal
	if chain.HasBadBlock(block.Hash()) {
		return core.ErrBlacklistedHash
	}

	err := chain.Engine().VerifyHeader(chain, block.Header(), false)
	if err == nil || err == types.ErrEmptyCommittedSeals {
		var (
			receipts types.Receipts
			usedGas  = new(uint64)
			gp       = new(core.GasPool).AddGas(block.GasLimit())
			header   = block.Header()
			parent   = chain.GetBlock(block.ParentHash(), block.NumberU64()-1)
		)

		// We need to process all of the transaction to get the latest state to get the latest committee
		state, stateErr := chain.StateAt(parent.Root())
		if stateErr != nil {
			return stateErr
		}

		// Validate the body of the proposal
		if err = chain.Validator().ValidateBody(block); err != nil {
			return err
		}

		// sb.chain.Processor().Process() was not called because it calls back Finalize() and would have modified the proposal
		// Instead only the transactions are applied to the copied state
		for i, tx := range block.Transactions() {
			state.Prepare(tx.Hash(), block.Hash(), i)
			vmConfig := vm.Config{
				EnablePreimageRecording: true,
				EWASMInterpreter:        "",
				EVMInterpreter:          "",
			}
			receipt, receiptErr := core.ApplyTransaction(chain.Config(), chain, nil, gp, state, header, tx, usedGas, vmConfig)
			if receiptErr != nil {
				return receiptErr
			}
			receipts = append(receipts, receipt)
		}

		state.Prepare(common.ACHash(block.Number()), block.Hash(), len(block.Transactions()))
		committeeSet, receipt, err := chain.Engine().Finalize(chain, header, state, block.Transactions(), nil, receipts)
		if err != nil {
			return err
		}
		receipts = append(receipts, receipt)
		//Validate the state of the proposal
		if err = chain.Validator().ValidateState(block, state, receipts, *usedGas); err != nil {
			return err
		}

		//Perform the actual comparison
		if len(header.Committee) != len(committeeSet) {
			return consensus.ErrInconsistentCommitteeSet
		}

		for i := range committeeSet {
			if header.Committee[i].Address != committeeSet[i].Address ||
				header.Committee[i].VotingPower.Cmp(committeeSet[i].VotingPower) != 0 {
				return consensus.ErrInconsistentCommitteeSet
			}
		}

		return nil
	}
	return err
}
