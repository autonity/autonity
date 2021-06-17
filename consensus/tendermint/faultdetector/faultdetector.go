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
	"math"
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
	PV
	PVN
	PVO
	PVO1
	PVO2
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

const (
	// todo: refine the window and buffer range in contract which can be tuned during run time.
	deltaBlocks          = 30                       // Wait until the GST + delta (30 blocks) to start rule scan.
	randomDelayRange     = 5000                     // (0, 5] seconds random delay range
	msgHeightBufferRange = uint64(deltaBlocks + 60) // buffer such range of msgs in height at msg store.
	tooFutureHeightRange = uint64(90)               // skip messages with height that is higher than current head over this value.
)

var (
	errDuplicatedMsg   = errors.New("duplicated msg")
	errEquivocation    = errors.New("equivocation happens")
	errFutureMsg       = errors.New("future height msg")
	errGarbageMsg      = errors.New("garbage msg")
	errNotCommitteeMsg = errors.New("msg from none committee member")
	errProposal        = errors.New("proposal have invalid values")
	errProposer        = errors.New("proposal is not from proposer")

	errNoEvidenceForPO  = errors.New("no evidence for innocence of rule PO")
	errNoEvidenceForPVN = errors.New("no evidence for innocence of rule PVN")
	errNoEvidenceForPVO = errors.New("no evidence for innocence of rule PVO")
	errNoEvidenceForC   = errors.New("no evidence for innocence of rule C")
	errNoEvidenceForC1  = errors.New("no evidence for innocence of rule C1")

	nilValue = common.Hash{}
)

// Proof is what to prove that one is misbehaving, one should be slashed when a valid Proof is rise.
type Proof struct {
	Type     uint8 // Misbehaviour, Accusation, Innocence.
	Rule     Rule
	Message  *tendermintCore.Message   // the msg to be considered as suspicious or misbehaved one
	Evidence []*tendermintCore.Message // the proofs of innocence or misbehaviour.
}

// FaultDetector it subscribe chain event to trigger rule engine to apply patterns over
// msg store, it send Proof of challenge if it detects any potential misbehavior, either it
// read state db on each new height to get latest challenges from autonity contract's view,
// and to prove its innocent if there were any challenges on the suspicious node.
type FaultDetector struct {
	proofWG           sync.WaitGroup
	faultDetectorFeed event.Feed

	tendermintMsgSub *event.TypeMuxSubscription
	blockCh          chan core.ChainEvent
	blockSub         event.Subscription

	blockchain BlockChainContext

	address  common.Address
	msgStore *MsgStore

	processFutureHeightMsgCh chan uint64
	misbehaviourProofsCh     chan *autonity.OnChainProof
	futureHeightMsgBuffer    map[uint64][]*tendermintCore.Message // map[blockHeight][]*tendermintMessages
	onChainProofsBuffer      []*autonity.OnChainProof             // buffer proofs to aggregate them into single TX.

	logger log.Logger
}

// call by ethereum object to create fd instance.
func NewFaultDetector(chain BlockChainContext, nodeAddress common.Address, sub *event.TypeMuxSubscription) *FaultDetector {
	fd := &FaultDetector{
		tendermintMsgSub:         sub,
		blockCh:                  make(chan core.ChainEvent, 300),
		blockchain:               chain,
		address:                  nodeAddress,
		msgStore:                 newMsgStore(),
		processFutureHeightMsgCh: make(chan uint64, deltaBlocks),
		misbehaviourProofsCh:     make(chan *autonity.OnChainProof, 100),
		futureHeightMsgBuffer:    make(map[uint64][]*tendermintCore.Message),
		logger:                   log.New("FaultDetector", nodeAddress),
	}

	fd.blockSub = fd.blockchain.SubscribeChainEvent(fd.blockCh)
	// register faultdetector contracts on evm's precompiled contract set.
	registerFaultDetectorContracts(chain)
	return fd
}

// listen for new block events from block-chain, do the tasks like take challenge and provide Proof for innocent, the
// Fault Detector rule engine could also triggered from here to scan those msgs of msg store by applying rules.
func (fd *FaultDetector) FaultDetectorEventLoop() {
	go fd.blockEventLoop(fd.blockCh, fd.misbehaviourProofsCh, fd.blockSub.Err())
	go fd.tendermintMsgEventLoop(fd.tendermintMsgSub.Chan(), fd.processFutureHeightMsgCh)
}

func (fd *FaultDetector) tendermintMsgEventLoop(tendermintMsgCh <-chan *event.TypeMuxEvent, futureHeightMsgsCh <-chan uint64) {
tendermintMsgLoop:
	for {
		curHeight := fd.blockchain.CurrentHeader().Number.Uint64()

		select {
		case ev, ok := <-tendermintMsgCh:
			if !ok {
				break tendermintMsgLoop
			}

			mv, ok := ev.Data.(events.MessageEvent)
			if !ok {
				fd.logger.Crit("programming error", "cannot cast message event to events.MessageEvent instead received ", ev.Data)
				break tendermintMsgLoop
			}

			msg := new(tendermintCore.Message)
			if err := msg.FromPayload(mv.Payload); err != nil {
				fd.logger.Error("fault detector: error while retrieving  payload", "err", err)
				continue tendermintMsgLoop
			}

			if curHeight > msgHeightBufferRange && msg.H() < curHeight-msgHeightBufferRange {
				fd.logger.Info("fault detector: discarding old message", "sender", msg.Sender())
				continue tendermintMsgLoop
			}

			if msg.H() > curHeight && msg.H()-curHeight > tooFutureHeightRange {
				fd.logger.Info("fault detector: discarding too future message", "sender", msg.Sender())
				continue tendermintMsgLoop
			}

			if err := fd.processMsg(msg); err != nil {
				fd.logger.Error("fault detector: error while processing consensus msg", "err", err)
				continue tendermintMsgLoop
			}

		case height, ok := <-futureHeightMsgsCh:
			if !ok {
				break tendermintMsgLoop
			}

			if curHeight > msgHeightBufferRange && height < curHeight-msgHeightBufferRange {
				fd.logger.Info("fault detector: discarding old height messages", "height", height)
				delete(fd.futureHeightMsgBuffer, height)
				continue tendermintMsgLoop
			}

			for h, msgs := range fd.futureHeightMsgBuffer {
				if h <= curHeight {
					for _, m := range msgs {
						if err := fd.processMsg(m); err != nil {
							fd.logger.Error("fault detector: error while processing consensus msg", "err", err)
						}
					}
					// once messages are processed, delete it from buffer.
					delete(fd.futureHeightMsgBuffer, h)
				}
			}
		}
	}
	close(fd.misbehaviourProofsCh)
}

func (fd *FaultDetector) blockEventLoop(blockCh <-chan core.ChainEvent, misbehaviourCh <-chan *autonity.OnChainProof, errCh <-chan error) {
blockChainLoop:
	for {
		select {
		// chain event update, provide proof of innocent if one is on challenge, rule engine scanning is triggered also.
		case ev, ok := <-blockCh:
			if !ok {
				break blockChainLoop
			}
			// before run rule engine over msg store, process any buffered msg.
			fd.processFutureHeightMsgCh <- ev.Block.NumberU64()

			// handle accusations and provide innocence Proof if there were any for a node.
			innocenceProofs, _ := fd.handleAccusations(ev.Block, ev.Block.Root())
			if innocenceProofs != nil {
				fd.onChainProofsBuffer = append(fd.onChainProofsBuffer, innocenceProofs...)
			}

			// run rule engine over a specific height.
			proofs := fd.runRuleEngine(ev.Block.NumberU64())
			if len(proofs) > 0 {
				fd.onChainProofsBuffer = append(fd.onChainProofsBuffer, proofs...)
			}

			// aggregate buffered proofs into single TX and send.
			if len(fd.onChainProofsBuffer) != 0 {
				copyOnChainProofs := make([]*autonity.OnChainProof, len(fd.onChainProofsBuffer))
				copy(copyOnChainProofs, fd.onChainProofsBuffer)
				fd.sendProofs(copyOnChainProofs)
				// release items from buffer
				fd.onChainProofsBuffer = fd.onChainProofsBuffer[:0]
			}

			// msg store delete msgs out of buffering window.
			fd.msgStore.DeleteMsgsAtHeight(ev.Block.NumberU64() - msgHeightBufferRange)
		case m, ok := <-misbehaviourCh:
			if !ok {
				break blockChainLoop
			}
			fd.onChainProofsBuffer = append(fd.onChainProofsBuffer, m)
		case err, ok := <-errCh:
			if ok {
				fd.logger.Crit("block subscription error", err.Error())
			}
			break blockChainLoop

		}
	}
	close(fd.processFutureHeightMsgCh)
}

func (fd *FaultDetector) Stop() {
	fd.blockSub.Unsubscribe()
	fd.tendermintMsgSub.Unsubscribe()
	fd.proofWG.Wait()
	unRegisterFaultDetectorContracts()
}

// call by ethereum object to subscribe proofs Events.
func (fd *FaultDetector) SubscribeFaultDetectorEvents(ch chan<- []*autonity.OnChainProof) event.Subscription {
	return fd.faultDetectorFeed.Subscribe(ch)
}

func (fd *FaultDetector) filterPresentedOnes(proofs []*autonity.OnChainProof) (result []*autonity.OnChainProof) {
	// get latest chain state.
	state, err := fd.blockchain.State()
	if err != nil {
		return nil
	}
	header := fd.blockchain.CurrentBlock().Header()

	proofsMap := make(map[common.Hash]*autonity.OnChainProof)
	for _, p := range proofs {
		proofsMap[p.Msghash] = p
	}

	contract := fd.blockchain.GetAutonityContract()
	presentedAccusation := contract.GetAccusations(header, state)
	presentedMisbehavior := contract.GetMisBehaviours(header, state)

	for _, p := range presentedAccusation {
		delete(proofsMap, p.Msghash)
	}

	for _, p := range presentedMisbehavior {
		delete(proofsMap, p.Msghash)
	}

	for _, p := range proofsMap {
		result = append(result, p)
	}

	return result
}

// convert the raw proofs into on-chain Proof which contains raw bytes of messages.
func (fd *FaultDetector) generateOnChainProof(p *Proof) (*autonity.OnChainProof, error) {
	var onChainProof = &autonity.OnChainProof{
		Type:    p.Type,
		Rule:    uint8(p.Rule),
		Sender:  p.Message.Address,
		Msghash: p.Message.MsgHash(),
	}

	rproof, err := rlp.EncodeToBytes(p)
	if err != nil {
		return nil, err
	}
	onChainProof.Rawproof = rproof
	return onChainProof, nil
}

// getInnocentProof called by client who is on challenge to get Proof of innocent from msg store.
func (fd *FaultDetector) getInnocentProof(c *Proof) (*autonity.OnChainProof, error) {
	var onChainProof *autonity.OnChainProof
	// rule engine have below provable accusation for the time being:
	switch c.Rule {
	case PO:
		return fd.getInnocentProofOfPO(c)
	case PVN:
		return fd.getInnocentProofOfPVN(c)
	case PVO:
		return fd.getInnocentProofOfPVO(c)
	case C:
		return fd.getInnocentProofOfC(c)
	case C1:
		return fd.getInnocentProofOfC1(c)
	default:
		return onChainProof, fmt.Errorf("not provable rule")
	}
}

// get Proof of innocent of C from msg store.
func (fd *FaultDetector) getInnocentProofOfC(c *Proof) (*autonity.OnChainProof, error) {
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
	p, err := fd.generateOnChainProof(&Proof{
		Type:     autonity.Innocence,
		Rule:     c.Rule,
		Message:  preCommit,
		Evidence: proposals,
	})
	return p, err
}

// get Proof of innocent of C1 from msg store.
func (fd *FaultDetector) getInnocentProofOfC1(c *Proof) (*autonity.OnChainProof, error) {
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

	p, err := fd.generateOnChainProof(&Proof{
		Type:     autonity.Innocence,
		Rule:     c.Rule,
		Message:  preCommit,
		Evidence: prevotesForV,
	})
	return p, err
}

// get Proof of innocent of PO from msg store.
func (fd *FaultDetector) getInnocentProofOfPO(c *Proof) (*autonity.OnChainProof, error) {
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

	p, err := fd.generateOnChainProof(&Proof{
		Type:     autonity.Innocence,
		Rule:     c.Rule,
		Message:  proposal,
		Evidence: prevotes,
	})
	return p, err
}

// get Proof of innocent of PVN from msg store.
func (fd *FaultDetector) getInnocentProofOfPVN(c *Proof) (*autonity.OnChainProof, error) {
	// get innocent proofs for PVN, for a prevote that vote for a new value,
	// then there must be a proposal for this new value.
	var onChainProof *autonity.OnChainProof
	prevote := c.Message
	height := prevote.H()

	correspondingProposals := fd.msgStore.Get(height, func(m *tendermintCore.Message) bool {
		return m.Type() == msgProposal && m.Value() == prevote.Value() && m.R() == prevote.R()
	})

	if len(correspondingProposals) == 0 {
		// cannot provide onChainProof for innocent of PVN, the on-chain contract will fine it latter once the
		// time window for onChainProof ends.
		return onChainProof, errNoEvidenceForPVN
	}

	p, err := fd.generateOnChainProof(&Proof{
		Type:     autonity.Innocence,
		Rule:     c.Rule,
		Message:  prevote,
		Evidence: correspondingProposals,
	})

	return p, err
}

// get Proof of innocent of PVO from msg store, it collects quorum preVotes for the value voted at a valid round.
func (fd *FaultDetector) getInnocentProofOfPVO(c *Proof) (*autonity.OnChainProof, error) {
	// get innocent proofs for PVO, collect quorum preVotes at the valid round of the old proposal.
	var onChainProof *autonity.OnChainProof
	oldProposal := c.Evidence[0]
	height := oldProposal.H()
	validRound := oldProposal.ValidRound()
	quorum := bft.Quorum(fd.blockchain.GetHeaderByNumber(height - 1).TotalVotingPower())

	preVotes := fd.msgStore.Get(height, func(m *tendermintCore.Message) bool {
		return m.Type() == msgPrevote && m.Value() == oldProposal.Value() && m.R() == validRound
	})

	if len(preVotes) == 0 {
		// cannot provide on-chain proof for accusation of PVO.
		return onChainProof, errNoEvidenceForPVO
	}

	votes := deEquivocatedMsgs(preVotes)
	if powerOfVotes(votes) < quorum {
		return onChainProof, errNoEvidenceForPVO
	}

	p, err := fd.generateOnChainProof(&Proof{
		Type:     autonity.Innocence,
		Rule:     c.Rule,
		Message:  c.Message,
		Evidence: append(c.Evidence, votes...),
	})
	return p, err
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

// processMsg, check and submit any auto-incriminating, equivocation challenges, and then only store checked msg into msg store.
func (fd *FaultDetector) processMsg(m *tendermintCore.Message) error {
	// pre-check if msg is from valid committee member
	err := checkMsgSignature(fd.blockchain, m)
	if err != nil {
		if err == errFutureMsg {
			fd.futureHeightMsgBuffer[m.H()] = append(fd.futureHeightMsgBuffer[m.H()], m)
		}
		return err
	}

	// decode consensus msg, and auto-incriminating msg is addressed here.
	err = checkAutoIncriminatingMsg(fd.blockchain, m)
	if err != nil {
		if err == errFutureMsg {
			fd.futureHeightMsgBuffer[m.H()] = append(fd.futureHeightMsgBuffer[m.H()], m)
		} else {
			proofs := []*tendermintCore.Message{m}
			fd.submitMisbehavior(m, proofs, err, fd.misbehaviourProofsCh)
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
		fd.submitMisbehavior(m, proofs, err, fd.misbehaviourProofsCh)
		return err
	}

	if err == errDuplicatedMsg {
		fd.logger.Warn("Consensus msg already processed by fault detector before", "msg sender", m.Sender())
		// todo: think about network layer's reputation slashing, when it exceed a threshold freeze the remote peer
		//  for period of time in case of DoS attack.
	}
	return nil
}

// run rule engine over latest msg store, if the return proofs is not empty, then rise challenge.
func (fd *FaultDetector) runRuleEngine(height uint64) []*autonity.OnChainProof {
	var onChainProofs []*autonity.OnChainProof
	// To avoid none necessary accusations, we wait for delta blocks to start rule scan.
	if height > uint64(deltaBlocks) {
		// run rule engine over the previous delta offset height.
		checkPointHeight := height - uint64(deltaBlocks)
		quorum := bft.Quorum(fd.blockchain.GetHeaderByNumber(checkPointHeight - 1).TotalVotingPower())
		proofs := fd.runRulesOverHeight(checkPointHeight, quorum)
		if len(proofs) > 0 {
			for i := 0; i < len(proofs); i++ {
				p, err := fd.generateOnChainProof(proofs[i])
				if err != nil {
					fd.logger.Warn("convert Proof to on-chain Proof", "faultdetector", err)
					continue
				}
				onChainProofs = append(onChainProofs, p)
			}
		}
	}
	return onChainProofs
}

func (fd *FaultDetector) runRulesOverHeight(height uint64, quorum uint64) (proofs []*Proof) {
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

	proposalsNew := fd.msgStore.Get(height, func(m *tendermintCore.Message) bool {
		return m.Type() == msgProposal && m.ValidRound() == -1
	})

	for _, p := range proposalsNew {
		proposal := p

		// Skip if proposal is equivocated
		proposalsForR := fd.msgStore.Get(height, func(m *tendermintCore.Message) bool {
			return m.Sender() == proposal.Sender() && m.Type() == msgProposal && m.R() == proposal.R()

		})
		// Due to the for loop there must be at least one proposal
		if len(proposalsForR) > 1 {
			continue
		}

		//check all precommits for previous rounds from this sender are nil
		precommits := fd.msgStore.Get(height, func(m *tendermintCore.Message) bool {
			return m.Sender() == proposal.Sender() && m.Type() == msgPrecommit && m.R() < proposal.R() && m.Value() != nilValue
		})
		if len(precommits) != 0 {
			proof := &Proof{
				Type:     autonity.Misbehaviour,
				Rule:     PN,
				Evidence: precommits,
				Message:  proposal,
			}
			proofs = append(proofs, proof)
			fd.logger.Info("Misbehaviour detected", "faultdetector", fd.address, "rulePN", PN, "sender", proposal.Sender())
		}
	}
	return proofs
}

func (fd *FaultDetector) oldProposalsAccountabilityCheck(height uint64, quorum uint64) (proofs []*Proof) {
	// ------------Old Proposal------------
	// PO: (Mr′<r,PV) ∧ (Mr′,PC|pi) ∧ (Mr′<r′′<r,P C|pi)∗ <--- (Mr,P|pi)
	// PO1: [#(Mr′,PV|V) ≥ 2f+ 1] ∧ [nil ∨ V ∨ ⊥] ∧ [nil ∨ ⊥] <--- [V]

	proposalsOld := fd.msgStore.Get(height, func(m *tendermintCore.Message) bool {
		return m.Type() == msgProposal && m.ValidRound() > -1
	})

oldProposalLoop:
	for _, p := range proposalsOld {
		proposal := p
		// Check that in the valid round we see a quorum of prevotes and that there is no precommit at all or a
		// precommit for v or nil.

		// Skip if proposal is equivocated
		proposalsForR := fd.msgStore.Get(height, func(m *tendermintCore.Message) bool {
			return m.Sender() == proposal.Sender() && m.Type() == msgProposal && m.R() == proposal.R()

		})
		// Due to the for loop there must be at least one proposal
		if len(proposalsForR) > 1 {
			continue oldProposalLoop
		}

		validRound := proposal.ValidRound()

		// Is there a precommit for a value other than nil or the proposed value by the current proposer in the valid
		// round? If there is, the proposer has proposed a value for which it is not locked on, thus a Proof of
		// misbehaviour can be generated.
		precommitsFromPiInVR := fd.msgStore.Get(height, func(m *tendermintCore.Message) bool {
			return m.Type() == msgPrecommit && m.R() == validRound && m.Sender() == proposal.Sender() &&
				m.Value() != nilValue && m.Value() != proposal.Value()
		})
		if len(precommitsFromPiInVR) > 0 {
			proof := &Proof{
				Type:     autonity.Misbehaviour,
				Rule:     PO,
				Evidence: precommitsFromPiInVR,
				Message:  proposal,
			}
			proofs = append(proofs, proof)
			fd.logger.Info("Misbehaviour detected", "faultdetector", fd.address, "rulePO", PO, "sender", proposal.Sender())
			continue oldProposalLoop
		}

		// Is there a precommit for anything other than nil from the proposer between the valid round and the round of
		// the proposal? If there is then that implies the proposer saw 2f+1 prevotes in that round and hence it should
		// have set that round as the valid round.
		preommitsFromPiAfterVR := fd.msgStore.Get(height, func(m *tendermintCore.Message) bool {
			return m.Type() == msgPrecommit && m.R() > validRound && m.R() < proposal.R() &&
				m.Sender() == proposal.Sender() && m.Value() != nilValue
		})
		if len(preommitsFromPiAfterVR) > 0 {
			proof := &Proof{
				Type:     autonity.Misbehaviour,
				Rule:     PO,
				Evidence: preommitsFromPiAfterVR,
				Message:  proposal,
			}
			proofs = append(proofs, proof)
			fd.logger.Info("Misbehaviour detected", "faultdetector", fd.address, "rulePO", PO, "sender", proposal.Sender())
			continue oldProposalLoop
		}

		// Do we see a quorum for a value other than the proposed value? If so, we have proof of misbehaviour.
		allPrevotesForValidRound := fd.msgStore.Get(height, func(m *tendermintCore.Message) bool {
			return m.Type() == msgPrevote && m.R() == validRound && m.Value() != nilValue && m.Value() != proposal.Value()
		})

		prevotesMap := make(map[common.Hash][]*tendermintCore.Message)
		for _, p := range allPrevotesForValidRound {
			prevotesMap[p.Value()] = append(prevotesMap[p.Value()], p)
		}

		for _, preVotes := range prevotesMap {
			// Here the assumption is that in a single round it is not possible to have 2 value which quorum votes,
			// this would imply at least quorum nodes are malicious which is much higher than our assumption.
			if powerOfVotes(deEquivocatedMsgs(preVotes)) >= quorum {
				proof := &Proof{
					Type:     autonity.Misbehaviour,
					Rule:     PO,
					Evidence: preVotes,
					Message:  proposal,
				}
				proofs = append(proofs, proof)
				fd.logger.Info("Misbehaviour detected", "faultdetector", fd.address, "rulePO", PO, "sender", proposal.Sender())
				continue oldProposalLoop
			}
		}

		// Do we see a quorum of prevotes in the valid round, if not we can raise an accusation, since we cannot be sure
		// that these prevotes don't exist
		prevotes := fd.msgStore.Get(height, func(m *tendermintCore.Message) bool {
			// since equivocation msgs are stored, we have to query those preVotes which has same value as the proposal.
			return m.Type() == msgPrevote && m.R() == validRound && m.Value() == proposal.Value()
		})

		if powerOfVotes(deEquivocatedMsgs(prevotes)) < quorum {
			accusation := &Proof{
				Type:    autonity.Accusation,
				Rule:    PO,
				Message: proposal,
			}
			proofs = append(proofs, accusation)
			fd.logger.Info("Accusation detected", "faultdetector", fd.address, "rulePO", PO, "sender", proposal.Sender())
		}
		continue oldProposalLoop
	}
	return proofs
}

func (fd *FaultDetector) prevotesAccountabilityCheck(height uint64, quorum uint64) (proofs []*Proof) {
	// ------------New and Old Prevotes------------

	//Todo: you will have to add labels to the code to make sure only one proof is raised per message
	// - Check the gap between the r and the last round in the preCommits array

	prevotes := fd.msgStore.Get(height, func(m *tendermintCore.Message) bool {
		return m.Type() == msgPrevote && m.Value() != nilValue
	})

prevotesLoop:
	for _, p := range prevotes {
		prevote := p

		// Skip if prevote is equivocated
		prevotesForR := fd.msgStore.Get(height, func(m *tendermintCore.Message) bool {
			return m.Sender() == prevote.Sender() && m.Type() == msgPrevote && m.R() == prevote.R()

		})
		// Due to the for loop there must be at least one proposal
		if len(prevotesForR) > 1 {
			continue prevotesLoop
		}

		// We need to check whether we have proposals from the prevote's round
		correspondingProposals := fd.msgStore.Get(height, func(m *tendermintCore.Message) bool {
			return m.Type() == msgProposal && m.R() == prevote.R()
		})

		//Todo: decide how to process PVN/PVO accusation since we cannot know which one it is unless we have the
		// corresponding proposal. We need to consider what to do when the prevote sender is also the proposer of the
		// the current round, we need to get the information on the proposer of the current round before we can create
		// an accusation. We can only create accusation on message pattern which involve messages from other users as
		// their signatures cannot be forged. However, in this case there is a caveat about what to do when the proposer
		// is also the sender of the prevote, do we raise an accusation? In such a case I don't think it would be wise
		// to create an accusation since they may just lie at the time of providing the proof.
		if len(correspondingProposals) == 0 {
			accusation := &Proof{
				Type: autonity.Accusation,
				Rule: PVN, //This could be PVO as well, however, we can't decide since there are no corresponding
				// proposal
				Message: prevote,
			}
			proofs = append(proofs, accusation)
			fd.logger.Info("Accusation detected", "faultdetector", fd.address, "rulePVN", PVN, "sender", prevote.Sender())
			continue prevotesLoop
		}

		// We need to ensure that we keep all proposals in the message store, so that we have the maximum chance of
		// finding justification for prevotes. This is to account for equivocation where the proposer send 2 proposals
		// with the same value but different valid rounds to different nodes. We can't penalise the sender of prevote
		// since we can't tell which proposal they received. We just want to find a set of message which fit the rule.
		// Therefore, we need to check all of the proposals to find a single one which shows the current prevote is
		// valid.
		var prevotesProofs []*Proof
		for _, cp := range correspondingProposals {
			correspondingProposal := cp
			if cp.Value() != prevote.Value() {
				misbehaviour := &Proof{
					Type:     autonity.Misbehaviour,
					Rule:     PV,
					Evidence: []*tendermintCore.Message{cp, prevote},
					Message:  prevote,
				}
				prevotesProofs = append(prevotesProofs, misbehaviour)
			} else {
				if correspondingProposal.ValidRound() == -1 {
					prevotesProofs = append(prevotesProofs, fd.newPrevotesAccountabilityCheck(height, prevote))
				} else {
					prevotesProofs = append(prevotesProofs, fd.oldPrevotesAccountabilityCheck(height, quorum, correspondingProposal, prevote))
				}
			}
		}

		// The current prevote is valid
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

func (fd *FaultDetector) newPrevotesAccountabilityCheck(height uint64, prevote *tendermintCore.Message) (proof *Proof) {
	// New Proposal, apply PVN rules

	// PVN: (Mr′<r,PC|pi)∧(Mr′<r′′<r,PC|pi)* ∧ (Mr,P|proposer(r)) <--- (Mr,PV|pi)

	// PVN2: [nil ∨ ⊥] ∧ [nil ∨ ⊥] ∧ [V:Valid(V)] <--- [V]: r′= 0,∀r′′< r:Mr′′,PC|pi=nil

	// PVN2, If there is a valid proposal V at round r, and pi never ever precommit(locked a value) before, then pi
	// should prevote for V or a nil in case of timeout at this round.

	// PVN3: [V] ∧ [nil ∨ ⊥] ∧ [V:Valid(V)] <--- [V]:∀r′< r′′<r,Mr′′,PC|pi=nil

	// There is no scope to raise an accusation for these rules since the only message in PVN that is not sent by pi is
	// the proposal and you require the proposal before you can even attempt to apply the rule.

	// Since we cannot raise an accusation we can only create a proof of misbehaviour. To create a proof of misbehaviour
	// we need to have all the messages in the message pattern, otherwise, we cannot make any statement about the
	// message. We may not have enough information and we don't want to accuse someone unnecessarily. To show a proof of
	// misbehaviour for PVN2 and PVN3 we need to collect all the precommits from pi and set the latest precommit round
	// as r' and we need to have all the precommit messages from r' to r for pi to be able to check for misbehaviour. If
	// the latest precommit is not for V and we have all the precommits from r' to r which are nil, then we have proof
	// of misbehaviour.
	precommitsFromPi := fd.msgStore.Get(height, func(m *tendermintCore.Message) bool {
		return m.Type() == msgPrecommit && prevote.Sender() == m.Sender() && m.R() < prevote.R()
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
		for i := len(precommitsFromPi) - 1; i >= 0 && math.Abs(float64(r)-float64(rPrime)) <= 1; i-- {
			if precommitsFromPi[i].Value() != nilValue {
				pc := precommitsFromPi[i]
				precommitsAtRPrime := fd.msgStore.Get(height, func(m *tendermintCore.Message) bool {
					return m.Type() == msgPrecommit && pc.Sender() == m.Sender() && m.R() == pc.R()
				})

				// Check for equivocation, it is possible there are multiple precommit from pi for the same round.
				// If there are equivocated messages: do nothing. Since pi has already being punished for equivocation
				// round when the equivocated message was first received.
				if len(precommitsAtRPrime) == 1 {
					if precommitsAtRPrime[0].Value() != prevote.Value() {
						// there is no equivocation
						fd.logger.Info("Misbehaviour detected", "faultdetector", fd.address, "rulePVN", PVN, "sender", prevote.Sender())
						return &Proof{
							Type:     autonity.Misbehaviour,
							Rule:     PVN,
							Evidence: precommitsFromPi[i:],
							Message:  prevote,
						}
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

func (fd *FaultDetector) oldPrevotesAccountabilityCheck(height uint64, quorum uint64, correspondingProposal *tendermintCore.Message, prevote *tendermintCore.Message) (proof *Proof) {
	currentR := correspondingProposal.R()
	validRound := correspondingProposal.ValidRound()

	// If there is a prevote for an old proposal then pi can only vote for v or send nil (see line 28 and 29 of
	// tendermint pseudocode). Therefore if in the valid round there is a quorum for a value other than v, we know pi
	// prevoted incorrectly. If the proposal was a bad proposal, then pi should not have voted for it. Thus we do not
	// need to make sure whether the proposal is correct or not (which we would in the proposal checking rules, however,
	// a bad proposal will still exist in our message store and it shouldn't have an impact on the checking of prevotes).
	allPrevotesForValidRound := fd.msgStore.Get(height, func(m *tendermintCore.Message) bool {
		return m.Type() == msgPrevote && m.R() == validRound && m.Value() != nilValue &&
			m.Value() != correspondingProposal.Value()
	})

	prevotesMap := make(map[common.Hash][]*tendermintCore.Message)
	for _, p := range allPrevotesForValidRound {
		prevotesMap[p.Value()] = append(prevotesMap[p.Value()], p)
	}

	for _, preVotes := range prevotesMap {
		// Here the assumption is that in a single round it is not possible to have 2 value which quorum votes,
		// this would imply at least quorum nodes are malicious which is much higher than our assumption.
		if powerOfVotes(deEquivocatedMsgs(preVotes)) >= quorum {
			fd.logger.Info("Misbehaviour detected", "faultdetector", fd.address, "rulePVO", PVO, "sender", prevote.Sender())
			proof := &Proof{
				Type:    autonity.Misbehaviour,
				Rule:    PVO,
				Message: prevote,
			}
			proof.Evidence = append(proof.Evidence, correspondingProposal)
			proof.Evidence = append(proof.Evidence, preVotes...)
			return proof
		}
	}

	prevotesForVFromValidRound := fd.msgStore.Get(height, func(m *tendermintCore.Message) bool {
		return m.Type() == msgPrevote && m.R() == validRound && m.Value() == correspondingProposal.Value()
	})

	// Check whether we have a quorum for v, if not raise an accusation
	if powerOfVotes(deEquivocatedMsgs(prevotesForVFromValidRound)) >= quorum {
		preCommitsForVFromPi := fd.msgStore.Get(height, func(m *tendermintCore.Message) bool {
			return m.Type() == msgPrecommit && m.R() >= validRound && m.R() < currentR &&
				m.Sender() == prevote.Sender() && m.Value() == prevote.Value()
		})

		if len(preCommitsForVFromPi) > 0 {
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

			sort.SliceStable(preCommitsForVFromPi, func(i, j int) bool {
				return preCommitsForVFromPi[i].R() < preCommitsForVFromPi[j].R()
			})
			// Get the round of the latest precommit for V from pi
			latestPrecommitForV := preCommitsForVFromPi[len(preCommitsForVFromPi)-1]
			preCommitsAfterLatestPrecommitForV := fd.msgStore.Get(height, func(m *tendermintCore.Message) bool {
				return m.Type() == msgPrecommit && m.R() > latestPrecommitForV.R() && m.R() < currentR &&
					m.Sender() == prevote.Sender()
			})

			// Due to equivocation we cannot use the range between latest round for V and current round to determine
			// whether precommits from all rounds between latest round for V and current round are present. We need to
			// ensure there are no gaps between
			for i := 1; i < len(preCommitsAfterLatestPrecommitForV); i++ {
				prev, cur := preCommitsAfterLatestPrecommitForV[i-1].R(), preCommitsAfterLatestPrecommitForV[i].R()
				diff := math.Abs(float64(cur) - float64(prev))
				if diff > 1 {
					// at least one round's precommit is missing
					return nil
				}
			}

			// We have precommits from all the round between latest round for V and current round. Thus, we can check
			// for misbehaviour.
			for _, v := range preCommitsAfterLatestPrecommitForV {
				if v.Value() != nilValue && v.Value() != prevote.Value() {
					fd.logger.Info("Misbehaviour detected", "faultdetector", fd.address, "rulePVO1", PVO1, "sender", prevote.Sender())
					proof := &Proof{
						Type:    autonity.Misbehaviour,
						Rule:    PVO1,
						Message: prevote,
					}
					proof.Evidence = append(proof.Evidence, correspondingProposal)
					proof.Evidence = append(proof.Evidence, latestPrecommitForV)
					proof.Evidence = append(proof.Evidence, preCommitsAfterLatestPrecommitForV...)
					return proof
				}
			}
		} else {
			// PVO’:(Mr′<r, PV) ∧ (Mr′<r′′<r, PC|pi)∗ ∧ (Mr,P|proposer(r)) ⇐= (Mr,P V|pi)
			// PVO2: [#(V)≥2f+ 1] ∧ [V ∨ nil ∨⊥] ∧ [V:validRound(V) =r′] ⇐= [V];
			// if V is the proposed value at round r with validRound(V) =r′ then there must be 2f+ 1 prevotes
			// for V issued at round r′. If moreover, pi did not precommit for other values in any round between
			// r′and r(thus it can be either locked on some values or not) then in round r pi prevotes for V.

			// We need to ensure that there are no precommits for V'. Since we already check for precommits for
			// V in the PVO1 rule we only need make sure that all the precommits are nil. Therefore, we don't
			// need to check for precommits for V, since the PVO1 block takes the latest V.

			precommitsFromPiAfterVR := fd.msgStore.Get(height, func(m *tendermintCore.Message) bool {
				return m.Type() == msgPrecommit && m.R() > validRound && m.R() < currentR && m.Sender() == prevote.Sender()
			})

			if len(precommitsFromPiAfterVR) > 0 {
				sort.SliceStable(precommitsFromPiAfterVR, func(i, j int) bool {
					return precommitsFromPiAfterVR[i].R() < precommitsFromPiAfterVR[j].R()
				})

				// Ensure there are no gaps
				for i := 1; i < len(precommitsFromPiAfterVR); i++ {
					prev, cur := precommitsFromPiAfterVR[i-1].R(), precommitsFromPiAfterVR[i].R()
					diff := math.Abs(float64(cur) - float64(prev))
					if diff > 1 {
						// at least one round's precommit is missing
						return nil
					}
				}

				// We have precommits from all the round between valid round and current round. Thus, we can check for
				// misbehaviour.
				for _, v := range precommitsFromPiAfterVR {
					if v.Value() != nilValue {
						fd.logger.Info("Misbehaviour detected", "faultdetector", fd.address, "rulePVO2", PVO2, "sender", prevote.Sender())
						proof := &Proof{
							Type:    autonity.Misbehaviour,
							Rule:    PVO2,
							Message: prevote,
						}
						proof.Evidence = append(proof.Evidence, correspondingProposal)
						proof.Evidence = append(proof.Evidence, precommitsFromPiAfterVR...)
						return proof
					}
				}

			}

		}
	} else {
		// raise an accusation
		fd.logger.Info("Accusation detected", "faultdetector", fd.address, "rulePVO", PVO, "sender", prevote.Sender())
		//Todo: We need to add more rules to distinguish between pvn accusations/misbehaviours from pvo
		// accusations/misbehaviours
		return &Proof{
			Type:     autonity.Accusation,
			Rule:     PVO,
			Message:  prevote,
			Evidence: []*tendermintCore.Message{correspondingProposal},
		}
	}
	return nil
}

func (fd *FaultDetector) precommitsAccountabilityCheck(height uint64, quorum uint64) (proofs []*Proof) {
	// ------------Precommits------------
	// C: [Mr,P|proposer(r)] ∧ [Mr,PV] <--- [Mr,PC|pi]
	// C1: [V:Valid(V)] ∧ [#(V) ≥ 2f+ 1] <--- [V]

	precommits := fd.msgStore.Get(height, func(m *tendermintCore.Message) bool {
		return m.Type() == msgPrecommit && m.Value() != nilValue
	})

	for _, preC := range precommits {
		precommit := preC

		// Skip if prevote is equivocated
		precommitsForR := fd.msgStore.Get(height, func(m *tendermintCore.Message) bool {
			return m.Sender() == precommit.Sender() && m.Type() == msgPrecommit && m.R() == precommit.R()

		})
		// Due to the for loop there must be at least one proposal
		if len(precommitsForR) > 1 {
			continue
		}

		proposals := fd.msgStore.Get(height, func(m *tendermintCore.Message) bool {
			return m.Type() == msgProposal && m.Value() == precommit.Value() && m.R() == precommit.R()
		})

		if len(proposals) == 0 {
			accusation := &Proof{
				Type:    autonity.Accusation,
				Rule:    C,
				Message: precommit,
			}
			proofs = append(proofs, accusation)
			fd.logger.Info("Accusation detected", "faultdetector", fd.address, "ruleC", C, "sender", precommit.Sender())
			continue
		}

		prevotesForNotV := fd.msgStore.Get(height, func(m *tendermintCore.Message) bool {
			return m.Type() == msgPrevote && m.Value() != precommit.Value() && m.R() == precommit.R()
		})
		prevotesForV := fd.msgStore.Get(height, func(m *tendermintCore.Message) bool {
			return m.Type() == msgPrevote && m.Value() == precommit.Value() && m.R() == precommit.R()
		})

		// even if we have equivocated preVotes for not V, we still assume that there are less f+1 malicious node in the
		// network, so the powerOfVotes of preVotesForNotV which was deEquivocated is still valid to prove that the
		// preCommit is a misbehaviour of rule C.
		deEquivocatedPreVotesForNotV := deEquivocatedMsgs(prevotesForNotV)
		if powerOfVotes(deEquivocatedPreVotesForNotV) >= quorum {
			// In this case there cannot be enough remaining prevotes
			// to justify a precommit for V.
			proof := &Proof{
				Type:     autonity.Misbehaviour,
				Rule:     C,
				Evidence: deEquivocatedPreVotesForNotV,
				Message:  precommit,
			}
			proofs = append(proofs, proof)
			fd.logger.Info("Misbehaviour detected", "faultdetector", fd.address, "ruleC", C, "sender", precommit.Sender())
		} else if powerOfVotes(prevotesForV) < quorum {
			// In this case we simply don't see enough prevotes to
			// justify the precommit.
			accusation := &Proof{
				Type:    autonity.Accusation,
				Rule:    C1,
				Message: precommit,
			}
			proofs = append(proofs, accusation)
			fd.logger.Info("Accusation detected", "faultdetector", fd.address, "ruleC1", C1, "sender", precommit.Sender())
		}
	}
	return proofs
}

// send proofs via event which will handled by ethereum object to signed the TX to send Proof.
func (fd *FaultDetector) sendProofs(proofs []*autonity.OnChainProof) {
	fd.proofWG.Add(1)
	go func() {
		defer fd.proofWG.Done()
		randomDelay()
		unPresented := fd.filterPresentedOnes(proofs)
		if len(unPresented) != 0 {
			fd.faultDetectorFeed.Send(unPresented)
		}
	}()
}

// submitMisbehavior takes proofs of misbehavior msg, and error id to form the on-chain proof, and
// send the proof of misbehavior to event channel.
func (fd *FaultDetector) submitMisbehavior(m *tendermintCore.Message, evidence []*tendermintCore.Message, err error, submitCh chan<- *autonity.OnChainProof) {
	rule, e := errorToRule(err)
	if e != nil {
		fd.logger.Warn("error to rule", "faultdetector", e)
	}
	proof, err := fd.generateOnChainProof(&Proof{
		Type:     autonity.Misbehaviour,
		Rule:     rule,
		Message:  m,
		Evidence: evidence,
	})
	if err != nil {
		fd.logger.Warn("generate misbehavior Proof", "faultdetector", err)
		return
	}

	// submit misbehavior proof to buffer, it will be sent once aggregated.
	submitCh <- proof
}

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
	header := chain.CurrentHeader()
	if m.H() > header.Number.Uint64()+1 {
		return errFutureMsg
	}

	lastHeader := chain.GetHeaderByNumber(m.H() - 1)
	if lastHeader == nil {
		return errFutureMsg
	}

	if _, err := m.Validate(crypto.CheckValidatorSignature, lastHeader); err != nil {
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
	proposer, err := getProposer(chain, m.H(), m.R())
	if err != nil {
		log.Error("get proposer err", "err", err)
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
	n := rand.Intn(randomDelayRange)
	// let min delay to be a single block period, to get latest state synced.
	if n < 1000 {
		n = 1000
	}
	time.Sleep(time.Duration(n) * time.Millisecond)
}

func sameVote(a *tendermintCore.Message, b *tendermintCore.Message) bool {
	return a.MsgHash() == b.MsgHash()
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
