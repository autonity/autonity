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
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/core/vm"
	"github.com/clearmatics/autonity/event"
	"github.com/clearmatics/autonity/log"
	"github.com/clearmatics/autonity/rlp"
	"github.com/syndtr/goleveldb/leveldb/errors"
	"math/big"
	"math/rand"
	"sort"
	"sync"
	"time"
)

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

// wrap chain context calls to make unit test easier
type ProposerGetter func(chain *core.BlockChain, h uint64, r int64) (common.Address, error)
type ProposalChecker func(chain *core.BlockChain, proposal types.Block) error

// Fault detector, it subscribe chain event to trigger rule engine to apply patterns over
// msg store, it send proof of challenge if it detects any potential misbehavior, either it
// read state db on each new height to get latest challenges from autonity contract's view,
// and to prove its innocent if there were any challenges on the suspicious node.
type FaultDetector struct {
	wg                sync.WaitGroup
	faultDetectorFeed event.Feed
	tendermintMsgSub  *event.TypeMuxSubscription

	blockChan  chan core.ChainEvent
	blockSub   event.Subscription
	blockchain *core.BlockChain

	address common.Address

	msgStore   *MsgStore
	futureMsgs map[uint64][]*tendermintCore.Message

	// buffer quorum for blocks.
	totalPowers map[uint64]uint64

	// buffer those proofs, aggregate them into single TX to send with latest nonce of account.
	bufferedProofs []autonity.OnChainProof

	logger log.Logger
}

// call by ethereum object to create fd instance.
func NewFaultDetector(chain *core.BlockChain, nodeAddress common.Address, sub *event.TypeMuxSubscription) *FaultDetector {
	fd := &FaultDetector{
		address:          nodeAddress,
		blockChan:        make(chan core.ChainEvent, 300),
		blockchain:       chain,
		msgStore:         newMsgStore(),
		logger:           log.New("FaultDetector", nodeAddress),
		tendermintMsgSub: sub,
		futureMsgs:       make(map[uint64][]*tendermintCore.Message),
		totalPowers:      make(map[uint64]uint64),
	}

	// register faultdetector contracts on evm's precompiled contract set.
	registerFaultDetectorContracts(chain)
	return fd
}

// listen for new block events from block-chain, do the tasks like take challenge and provide proof for innocent, the
// Fault Detector rule engine could also triggered from here to scan those msgs of msg store by applying rules.
func (fd *FaultDetector) FaultDetectorEventLoop(ctx context.Context) {
	fd.blockSub = fd.blockchain.SubscribeChainEvent(fd.blockChan)

	for {
		select {
		// chain event update, provide proof of innocent if one is on challenge, rule engine scanning is triggered also.
		case ev := <-fd.blockChan:
			fd.totalPowers[ev.Block.Number().Uint64()] = ev.Block.Header().TotalVotingPower()

			// before run rule engine over msg store, process any buffered msg.
			fd.processBufferedMsgs(ev.Block.NumberU64())

			// handle accusations and provide innocence proof if there were any for a node.
			innocenceProofs, _ := fd.handleAccusations(ev.Block, ev.Block.Root())
			if innocenceProofs != nil {
				fd.bufferedProofs = append(fd.bufferedProofs, innocenceProofs...)
			}

			// run rule engine over a specific height.
			proofs := fd.runRuleEngine(ev.Block.NumberU64())
			if proofs != nil {
				fd.bufferedProofs = append(fd.bufferedProofs, proofs...)
			}

			// aggregate buffered proofs into single TX and send.
			fd.sentProofs()

			// msg store delete msgs out of buffering window.
			fd.msgStore.DeleteMsgsAtHeight(ev.Block.NumberU64() - uint64(msgBufferInHeight))

			// delete power out of buffering window.
			delete(fd.totalPowers, ev.Block.NumberU64()-uint64(msgBufferInHeight))
		// to handle consensus msg from p2p layer.
		case ev, ok := <-fd.tendermintMsgSub.Chan():
			if !ok {
				break eventLoop
			}
			switch e := ev.Data.(type) {
			case events.MessageEvent:
				msg := new(tendermintCore.Message)
				if err := msg.FromPayload(e.Payload); err != nil {
					fd.logger.Error("invalid payload", "faultdetector", err)
					continue
				}

				// discard too old messages which is out of accountability buffering window.
				head := fd.blockchain.CurrentHeader().Number.Uint64()
				if msg.H() < head-uint64(msgBufferInHeight) {
					fd.logger.Info("discard too old message for accountability", "faultdetector", msg.Sender())
					continue
				}

				if err := fd.processMsg(msg); err != nil {
					fd.logger.Warn("process consensus msg", "faultdetector", err)
					continue
				}
			}

		case <-fd.blockSub.Err():
			break eventLoop
		case <-ctx.Done():
			fd.logger.Info("FaultDetectorEventLoop is stopped", "event", ctx.Err())
			break eventLoop
		}
	}

	fd.stopped <- struct{}{}
}

func (fd *FaultDetector) Stop() {
	fd.blockSub.Unsubscribe()
	fd.tendermintMsgSub.Unsubscribe()
	fd.wg.Wait()
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

	fd.futureMsgs[h.Uint64()] = append(fd.futureMsgs[h.Uint64()], m)
}

func (fd *FaultDetector) filterPresentedOnes(proofs *[]autonity.OnChainProof) []autonity.OnChainProof {
	// get latest chain state.
	var result []autonity.OnChainProof
	state, err := fd.blockchain.State()
	if err != nil {
		return nil
	}
	header := fd.blockchain.CurrentBlock().Header()

	presentedAccusation := fd.blockchain.GetAutonityContract().GetAccusations(header, state)
	presentedMisbehavior := fd.blockchain.GetAutonityContract().GetMisBehaviours(header, state)

	for i := 0; i < len(*proofs); i++ {
		present := false
		for j := 0; j < len(presentedAccusation); j++ {
			if (*proofs)[i].Msghash == presentedAccusation[j].Msghash &&
				(*proofs)[i].Type.Cmp(new(big.Int).SetUint64(uint64(Accusation))) == 0 {
				present = true
			}
		}

		for j := 0; j < len(presentedMisbehavior); j++ {
			if (*proofs)[i].Msghash == presentedMisbehavior[j].Msghash &&
				(*proofs)[i].Type.Cmp(new(big.Int).SetUint64(uint64(Misbehaviour))) == 0 {
				present = true
			}
		}

		if !present {
			result = append(result, (*proofs)[i])
		}
	}

	return result
}

// convert the raw proofs into on-chain proof which contains raw bytes of messages.
func (fd *FaultDetector) generateOnChainProof(m *tendermintCore.Message, proofs []tendermintCore.Message, rule Rule, t ProofType) (autonity.OnChainProof, error) {
	var proof autonity.OnChainProof
	proof.Sender = m.Address
	proof.Msghash = types.RLPHash(m.Payload())
	proof.Type = new(big.Int).SetUint64(uint64(t))

	var rawProof RawProof
	rawProof.Rule = rule
	// generate raw bytes encoded in rlp, it is by passed into precompiled contracts.
	rawProof.Message = m.Payload()
	for i := 0; i < len(proofs); i++ {
		rawProof.Evidence = append(rawProof.Evidence, proofs[i].Payload())
	}

	rp, err := rlp.EncodeToBytes(&rawProof)
	if err != nil {
		fd.logger.Warn("fail to rlp encode raw proof", "faultdetector", err)
		return proof, err
	}

	proof.Rawproof = rp
	return proof, nil
}

// getInnocentProof called by client who is on challenge to get proof of innocent from msg store.
func (fd *FaultDetector) getInnocentProof(c *Proof) (autonity.OnChainProof, error) {
	var proof autonity.OnChainProof
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
		return proof, fmt.Errorf("not provable rule")
	}
}

// get proof of innocent of C from msg store.
func (fd *FaultDetector) getInnocentProofOfC(c *Proof) (autonity.OnChainProof, error) {
	var proof autonity.OnChainProof
	preCommit := c.Message
	height := preCommit.H()

	proposals := fd.msgStore.Get(height, func(m *tendermintCore.Message) bool {
		return m.Type() == msgProposal && m.Value() == preCommit.Value() && m.R() == preCommit.R()
	})

	if len(proposals) == 0 {
		// cannot proof its innocent for PVN, the on-chain contract will fine it latter once the
		// time window for proof ends.
		return proof, errNoEvidenceForC
	}
	p, err := fd.generateOnChainProof(&preCommit, proposals, c.Rule, Innocence)
	if err != nil {
		return p, err
	}
	return p, nil
}

// get proof of innocent of C1 from msg store.
func (fd *FaultDetector) getInnocentProofOfC1(c *Proof) (autonity.OnChainProof, error) {
	var proof autonity.OnChainProof
	preCommit := c.Message
	height := preCommit.H()
	quorum := fd.quorum(height - 1)

	prevotesForV := fd.msgStore.Get(height, func(m *tendermintCore.Message) bool {
		return m.Type() == msgPrevote && m.Value() == preCommit.Value() && m.R() == preCommit.R()
	})

	if powerOfVotes(prevotesForV) < quorum {
		// cannot proof its innocent for PO for now, the on-chain contract will fine it latter once the
		// time window for proof ends.
		return proof, errNoEvidenceForC1
	}

	p, err := fd.generateOnChainProof(&preCommit, prevotesForV, c.Rule, Innocence)
	if err != nil {
		return p, err
	}

	return p, nil
}

// get proof of innocent of PO from msg store.
func (fd *FaultDetector) getInnocentProofOfPO(c *Proof) (autonity.OnChainProof, error) {
	// PO: node propose an old value with an validRound, innocent proof of it should be:
	// there are quorum num of prevote for that value at the validRound.
	var proof autonity.OnChainProof
	proposal := c.Message
	height := proposal.H()
	validRound := proposal.ValidRound()
	quorum := fd.quorum(height - 1)

	prevotes := fd.msgStore.Get(height, func(m *tendermintCore.Message) bool {
		return m.Type() == msgPrevote && m.R() == validRound && m.Value() == proposal.Value()
	})

	if powerOfVotes(prevotes) < quorum {
		// cannot proof its innocent for PO, the on-chain contract will fine it latter once the
		// time window for proof ends.
		return proof, errNoEvidenceForPO
	}

	p, err := fd.generateOnChainProof(&proposal, prevotes, c.Rule, Innocence)
	if err != nil {
		return p, err
	}

	return p, nil
}

// get proof of innocent of PVN from msg store.
func (fd *FaultDetector) getInnocentProofOfPVN(c *Proof) (autonity.OnChainProof, error) {
	// get innocent proofs for PVN, for a prevote that vote for a new value,
	// then there must be a proposal for this new value.
	var proof autonity.OnChainProof
	prevote := c.Message
	height := prevote.H()

	correspondingProposals := fd.msgStore.Get(height, func(m *tendermintCore.Message) bool {
		return m.Type() == msgProposal && m.Value() == prevote.Value() && m.R() == prevote.R()
	})

	if len(correspondingProposals) == 0 {
		// cannot proof its innocent for PVN, the on-chain contract will fine it latter once the
		// time window for proof ends.
		return proof, errNoEvidenceForPVN
	}

	p, err := fd.generateOnChainProof(&prevote, correspondingProposals, c.Rule, Innocence)
	if err != nil {
		return p, nil
	}

	return p, nil
}

// get accusations from chain via autonityContract calls, and provide innocent proofs if there were any challenge on node.
func (fd *FaultDetector) handleAccusations(block *types.Block, hash common.Hash) ([]autonity.OnChainProof, error) {
	var innocentProofs []autonity.OnChainProof
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
			c, err := decodeProof(accusations[i].Rawproof)
			if err != nil {
				continue
			}

			p, err := fd.getInnocentProof(c)
			if err != nil {
				continue
			}
			innocentProofs = append(innocentProofs, p)
		}
	}

	return innocentProofs, nil
}

// processBufferedMsgs, called on chain event update, it process msgs from the latest height buffered before.
func (fd *FaultDetector) processBufferedMsgs(height uint64) {
	for h, msgs := range fd.futureMsgs {
		if h <= height {
			for i := 0; i < len(msgs); i++ {
				if err := fd.processMsg(msgs[i]); err != nil {
					fd.logger.Warn("process consensus msg", "faultdetector", err)
					continue
				}
			}
		}
	}
}

// processMsg, check and submit any auto-incriminating, equivocation challenges, and then only store checked msg into msg store.
func (fd *FaultDetector) processMsg(m *tendermintCore.Message) error {
	// pre-check if msg is from valid committee member
	err := checkMsgSignature(fd.blockchain, m, getHeader, currentHeader)
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
			proofs := []tendermintCore.Message{*m}
			fd.submitMisbehavior(m, proofs, err)
			return err
		}
	}

	// store msg, if there is equivocation, msg store would then rise errEquivocation and proofs.
	msgs, err := fd.msgStore.Save(m)
	if err == errEquivocation && msgs != nil {
		var proofs []tendermintCore.Message
		for i := 0; i < len(msgs); i++ {
			proofs = append(proofs, *msgs[i])
		}
		fd.submitMisbehavior(m, proofs, err)
		return err
	}
	return nil
}

func (fd *FaultDetector) quorum(h uint64) uint64 {
	power, ok := fd.totalPowers[h]
	if ok {
		return bft.Quorum(power)
	}

	return bft.Quorum(fd.blockchain.GetHeaderByNumber(h).TotalVotingPower())
}

// run rule engine over latest msg store, if the return proofs is not empty, then rise challenge.
func (fd *FaultDetector) runRuleEngine(height uint64) []autonity.OnChainProof {
	var onChainProofs []autonity.OnChainProof
	// To avoid none necessary accusations, we wait for delta blocks to start rule scan.
	if height > uint64(deltaToWaitForAccountability) {
		// run rule engine over the previous delta offset height.
		lastDeltaHeight := height - uint64(deltaToWaitForAccountability)
		quorum := fd.quorum(lastDeltaHeight - 1)
		proofs := fd.runRulesOverHeight(lastDeltaHeight, quorum)
		if len(proofs) > 0 {
			for i := 0; i < len(proofs); i++ {
				p, err := fd.generateOnChainProof(&proofs[i].Message, proofs[i].Evidence, proofs[i].Rule, proofs[i].Type)
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

func (fd *FaultDetector) runRulesOverHeight(height uint64, quorum uint64) (proofs []Proof) {
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
			proof := Proof{
				Type:     Misbehaviour,
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
			proof := Proof{
				Type:     Misbehaviour,
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
			proof := Proof{
				Type:     Misbehaviour,
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
			return m.Type() == msgPrevote && m.R() == validRound
		})

		if powerOfVotes(deEquivocatedMsgs(prevotes)) < quorum {
			accusation := Proof{
				Type:    Accusation,
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
			accusation := Proof{
				Type: Accusation,
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
					proof := Proof{
						Type:     Misbehaviour,
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
					return m.Type() == msgPrecommit && prevote.Sender() == m.Sender() &&
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
			accusation := Proof{
				Type:    Accusation,
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

		if powerOfVotes(prevotesForNotV) >= quorum {
			// In this case there cannot be enough remaining prevotes
			// to justify a precommit for V.
			proof := Proof{
				Type:     Misbehaviour,
				Rule:     C,
				Evidence: deEquivocatedMsgs(prevotesForNotV),
				Message:  precommit,
			}
			proofs = append(proofs, proof)

		} else if powerOfVotes(prevotesForV) < quorum {
			// In this case we simply don't see enough prevotes to
			// justify the precommit.
			accusation := Proof{
				Type:    Accusation,
				Rule:    C1,
				Message: precommit,
			}
			proofs = append(proofs, accusation)
		}
	}

	return proofs
}

// send proofs via event which will handled by ethereum object to signed the TX to send proof.
func (fd *FaultDetector) sendProofs(proofs []autonity.OnChainProof) {
	fd.wg.Add(1)
	go func() {
		defer fd.wg.Done()
		randomDelay()
		unPresented := fd.filterPresentedOnes(&proofs)
		if len(unPresented) != 0 {
			fd.faultDetectorFeed.Send(AccountabilityEvent{Proofs: unPresented})
		}
	}()
}

func (fd *FaultDetector) sentProofs() {
	// todo: weight proofs before deliver it to pool since the max size of a TX is limited to 512 KB.
	//  consider to break down into multiples if it cannot fit in.
	if len(fd.bufferedProofs) != 0 {
		copyProofs := make([]autonity.OnChainProof, len(fd.bufferedProofs))
		copy(copyProofs, fd.bufferedProofs)
		fd.sendProofs(copyProofs)
		// release items from buffer
		fd.bufferedProofs = fd.bufferedProofs[:0]
	}
}

// submitMisbehavior takes proofs of misbehavior msg, and error id to form the on-chain proof, and
// send the proof of misbehavior to event channel.
func (fd *FaultDetector) submitMisbehavior(m *tendermintCore.Message, proofs []tendermintCore.Message, err error) {
	rule, e := errorToRule(err)
	if e != nil {
		fd.logger.Warn("error to rule", "faultdetector", e)
	}
	proof, err := fd.generateOnChainProof(m, proofs, rule, Misbehaviour)
	if err != nil {
		fd.logger.Warn("generate misbehavior proof", "faultdetector", err)
		return
	}

	// submit misbehavior proof to buffer, it will be sent once aggregated.
	fd.bufferedProofs = append(fd.bufferedProofs, proof)
}

/////// common helper functions shared between faultdetector and precompiled contract to validate msgs.

// decode consensus msgs, address garbage msg and invalid proposal by returning error.
func checkAutoIncriminatingMsg(chain *core.BlockChain, m *tendermintCore.Message) error {
	if m.Code == msgProposal {
		return checkProposal(chain, m, verifyProposal)
	}

	if m.Code == msgPrevote || m.Code == msgPrecommit {
		return decodeVote(m)
	}

	return errors.New("unknown consensus msg")
}

func checkEquivocation(chain *core.BlockChain, m *tendermintCore.Message, proof []tendermintCore.Message) error {
	// decode msgs
	err := checkAutoIncriminatingMsg(chain, m)
	if err != nil {
		return err
	}

	for i := 0; i < len(proof); i++ {
		err := checkAutoIncriminatingMsg(chain, &proof[i])
		if err != nil {
			return err
		}
	}
	// check equivocations.
	if !sameVote(m, &proof[0]) {
		return errEquivocation
	}
	return nil
}

//checkMsgSignature, it check if msg is from valid member of the committee.
func checkMsgSignature(chain *core.BlockChain, m *tendermintCore.Message, getHeader HeaderGetter, currentHeader CurrentHeaderGetter) error {
	msgHeight, err := m.Height()
	if err != nil {
		return err
	}

	header := currentHeader(chain)
	if msgHeight.Uint64() > header.Number.Uint64()+1 {
		return errFutureMsg
	}

	lastHeader := getHeader(chain, msgHeight.Uint64()-1)
	if lastHeader == nil {
		return errFutureMsg
	}

	if _, err = m.Validate(crypto.CheckValidatorSignature, lastHeader); err != nil {
		return errNotCommitteeMsg
	}
	return nil
}

// checkProposal, checks if proposal is valid and it's from correct proposer.
func checkProposal(chain *core.BlockChain, m *tendermintCore.Message, validateProposal ProposalChecker) error {
	var proposal tendermintCore.Proposal
	err := m.Decode(&proposal)
	if err != nil {
		return errGarbageMsg
	}
	if !isProposerMsg(chain, m, getProposer) {
		return errProposer
	}

	err = validateProposal(chain, *proposal.ProposalBlock)
	// due to network delay or timing issue, when Fault Detector validate a proposal, that proposal could already be
	// committed on the chain view.
	// since the msg sender were checked with correct proposer, so we consider to take it as a valid proposal.
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

func deEquivocatedMsgs(msgs []tendermintCore.Message) (deEquivocated []tendermintCore.Message) {
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

func getProposer(chain *core.BlockChain, h uint64, r int64) (common.Address, error) {
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

func isProposerMsg(chain *core.BlockChain, m *tendermintCore.Message, proposerGetter ProposerGetter) bool {
	h, _ := m.Height()
	r, _ := m.Round()

	proposer, err := proposerGetter(chain, h.Uint64(), r)
	if err != nil {
		return false
	}

	return m.Address == proposer
}

func powerOfVotes(votes []tendermintCore.Message) uint64 {
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
	// wait for random milliseconds (under the range of 10 seconds) to check if need to rise challenge.
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(randomDelayWindow)
	time.Sleep(time.Duration(n) * time.Millisecond)
}

func sameVote(a *tendermintCore.Message, b *tendermintCore.Message) bool {
	ah, _ := a.Height()
	ar, _ := a.Round()
	bh, _ := b.Height()
	br, _ := b.Round()
	aHash := types.RLPHash(a.Payload())
	bHash := types.RLPHash(b.Payload())

	if ah == bh && ar == br && a.Code == b.Code && a.Address == b.Address && aHash == bHash {
		return true
	}
	return false
}

func verifyProposal(chain *core.BlockChain, proposal types.Block) error {
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
