package misbehaviourdetector

import (
	"bytes"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math"
	"sort"
	"sync"
	"time"

	"github.com/autonity/autonity/autonity"
	"github.com/autonity/autonity/common"
	proto "github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/consensus/tendermint/bft"
	engineCore "github.com/autonity/autonity/consensus/tendermint/core"
	"github.com/autonity/autonity/consensus/tendermint/core/constants"
	mUtils "github.com/autonity/autonity/consensus/tendermint/core/messageutils"
	"github.com/autonity/autonity/consensus/tendermint/events"
	"github.com/autonity/autonity/core"
	"github.com/autonity/autonity/core/state"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/event"
	"github.com/autonity/autonity/internal/ethapi"
	"github.com/autonity/autonity/log"
	"github.com/autonity/autonity/rlp"
)

type BlockChainContext interface {
	proto.ChainReader
	CurrentBlock() *types.Block
	SubscribeChainEvent(ch chan<- core.ChainEvent) event.Subscription
	State() (*state.StateDB, error)
	GetAutonityContract() *autonity.Contract
	StateAt(root common.Hash) (*state.StateDB, error)
	HasBadBlock(hash common.Hash) bool
	Validator() core.Validator
}

const (
	msgGCInterval                 = 60      // every 60 blocks to GC msg store.
	offChainAccusationProofWindow = 10      // the time window in block for one to provide off chain innocence proof before it is escalated on chain.
	maxNumOfInnocenceProofCached  = 120 * 4 // 120 blocks with 4 on each height that rule engine can produce totally over a height.
	maxAccusationRatePerHeight    = 4       // max number of accusation can be produced by rule engine over a height against to a validator.
	maxMsgsCachedForFutureHeight  = 1000    // max num of msg buffer for the future heights.
)

var (
	errWrongSignatureMsg     = errors.New("invalid signature of message")
	errAccountableGarbageMsg = errors.New("accountable garbage message")
	errInvalidRound          = errors.New("invalid round")
	errWrongValidRound       = errors.New("wrong valid-round")
	errDuplicatedMsg         = errors.New("duplicated msg")
	errEquivocation          = errors.New("equivocation happens")
	errFutureMsg             = errors.New("future height msg")
	errNotCommitteeMsg       = errors.New("msg from none committee member")
	errProposer              = errors.New("proposal is not from proposer")

	errNoEvidenceForPO  = errors.New("no evidence for innocence of rule PO")
	errNoEvidenceForPVN = errors.New("no evidence for innocence of rule PVN")
	errNoEvidenceForPVO = errors.New("no evidence for innocence of rule PVO")
	errNoEvidenceForC1  = errors.New("no evidence for innocence of rule C1")

	nilValue = common.Hash{}
)

// AccountabilityProof is what to prove that one is misbehaving, one should be slashed when a valid AccountabilityProof is rise.
type AccountabilityProof struct {
	Type     autonity.AccountabilityEventType // Accountability event types: Misbehaviour, Accusation, Innocence.
	Rule     autonity.Rule                    // Rule ID defined in AFD rule engine.
	Message  *mUtils.Message                  // the consensus message which is accountable.
	Evidence []*mUtils.Message                // the proofs of the accountability event.
}

// FaultDetector it subscribe chain event to trigger rule engine to apply patterns over
// msg store, it sends AccountabilityProof of challenge if it detects any potential misbehavior, either it
// read state db on each new height to get the latest challenges from autonity contract's view,
// and to prove its innocent if there were any challenges on the suspicious node.
type FaultDetector struct {
	innocenceProofBuff *InnocenceProofBuffer
	rateLimiter        *OffChainAccusationRateLimiter

	proofWG           sync.WaitGroup
	faultDetectorFeed event.Feed

	tendermintMsgSub *event.TypeMuxSubscription

	txPool                   *core.TxPool
	ethBackend               ethapi.Backend
	nodeKey                  *ecdsa.PrivateKey
	accountabilityTXCh       chan []*autonity.AccountabilityEvent
	accountabilityTXEventSub event.Subscription

	// chain event subscriber for rule engine.
	ruleEngineBlockCh  chan core.ChainEvent
	ruleEngineBlockSub event.Subscription

	blockchain BlockChainContext
	address    common.Address
	msgStore   *engineCore.MsgStore

	// chain event subscriber for msg handler.
	msgHandlerBlockCh  chan core.ChainEvent
	msgHandlerBlockSub event.Subscription

	misbehaviourProofsCh      chan *autonity.AccountabilityEvent
	futureHeightMsgBuffer     map[uint64][]*mUtils.Message    // map[blockHeight][]*tendermintMessages
	totalBufferedFutureMsg    uint64                          // a counter to count the total cached future height msg.
	accountabilityEventBuffer []*autonity.AccountabilityEvent // accountability event buffer.

	offChainAccusationsMu sync.RWMutex
	//offChainAccusations   []*AccountabilityProof // off chain accusations list, ordered in chain height from low to high.
	offChainAccusations map[common.Hash]*AccountabilityProof // off chain accusation map, index by the proof's rlp hash.
	broadcaster         proto.Broadcaster

	logger log.Logger
}

// NewFaultDetector call by ethereum object to create fd instance.
func NewFaultDetector(chain BlockChainContext, nodeAddress common.Address, sub *event.TypeMuxSubscription,
	ms *engineCore.MsgStore, txPool *core.TxPool, ethBackend ethapi.Backend, nodeKey *ecdsa.PrivateKey) *FaultDetector {
	fd := &FaultDetector{
		innocenceProofBuff:     NewInnocenceProofBuffer(),
		rateLimiter:            NewOffChainAccusationRateLimiter(),
		txPool:                 txPool,
		ethBackend:             ethBackend,
		nodeKey:                nodeKey,
		tendermintMsgSub:       sub,
		ruleEngineBlockCh:      make(chan core.ChainEvent, 300),
		blockchain:             chain,
		address:                nodeAddress,
		msgStore:               ms,
		msgHandlerBlockCh:      make(chan core.ChainEvent, 300),
		accountabilityTXCh:     make(chan []*autonity.AccountabilityEvent),
		misbehaviourProofsCh:   make(chan *autonity.AccountabilityEvent, 100),
		futureHeightMsgBuffer:  make(map[uint64][]*mUtils.Message),
		offChainAccusations:    make(map[common.Hash]*AccountabilityProof),
		totalBufferedFutureMsg: 0,
		logger:                 log.New("FaultDetector", nodeAddress),
	}

	fd.ruleEngineBlockSub = fd.blockchain.SubscribeChainEvent(fd.ruleEngineBlockCh)
	fd.msgHandlerBlockSub = fd.blockchain.SubscribeChainEvent(fd.msgHandlerBlockCh)
	fd.accountabilityTXEventSub = fd.faultDetectorFeed.Subscribe(fd.accountabilityTXCh)
	return fd
}

// FaultDetectorEventLoop listen for new block events from blockchain, do the tasks like take challenge and provide AccountabilityProof for innocent, the
// Fault Detector rule engine could also trigger from here to scan those msgs of msg store by applying rules.
func (fd *FaultDetector) FaultDetectorEventLoop() {
	go fd.faultDetectorTXEventLoop()
	go fd.ruleEngineLoop()
	go fd.msgHandlerLoop()
}

func (fd *FaultDetector) tooOldHeightMsg(headHeight uint64, msgHeight uint64) bool {
	return headHeight > proto.AccountabilityHeightRange && msgHeight < headHeight-proto.AccountabilityHeightRange
}

func (fd *FaultDetector) SetBroadcaster(broadcaster proto.Broadcaster) {
	fd.broadcaster = broadcaster
}

func (fd *FaultDetector) bufferFutureHeightMsg(m *mUtils.Message) {

	fd.futureHeightMsgBuffer[m.H()] = append(fd.futureHeightMsgBuffer[m.H()], m)
	fd.totalBufferedFutureMsg++

	// buffer is full, remove the furthest away msg from buffer to prevent DoS attack.
	if fd.totalBufferedFutureMsg >= maxMsgsCachedForFutureHeight {
		maxHeight := m.H()
		for h, msgs := range fd.futureHeightMsgBuffer {
			if h > maxHeight && len(msgs) > 0 {
				maxHeight = h
			}
		}
		if len(fd.futureHeightMsgBuffer[maxHeight]) > 1 {
			fd.futureHeightMsgBuffer[maxHeight] = fd.futureHeightMsgBuffer[maxHeight][:len(fd.futureHeightMsgBuffer[maxHeight])-1]
		} else {
			delete(fd.futureHeightMsgBuffer, maxHeight)
		}

		fd.totalBufferedFutureMsg--
	}
}

func (fd *FaultDetector) deleteFutureHeightMsg(height uint64) {
	length := len(fd.futureHeightMsgBuffer[height])
	fd.totalBufferedFutureMsg = fd.totalBufferedFutureMsg - uint64(length)
	delete(fd.futureHeightMsgBuffer, height)
}

// decodeLiteProposalMsg it decodes the m from the TbftMsgBytes in rlp encoded lite proposal, then check the signature.
func decodeLiteProposalMsg(m *mUtils.Message) error {
	var liteProposal mUtils.LiteProposal
	if err := m.Decode(&liteProposal); err != nil {
		return err
	}

	// check signature,
	if err := liteProposal.ValidSignature(m.Address); err != nil {
		return err
	}

	return nil
}

func decodeConsensusMsg(m *mUtils.Message) error {

	//if msg is a lite proposal that is for accountability , decode in lite proposal.
	if m.Code == proto.MsgLiteProposal {
		return decodeLiteProposalMsg(m)
	}

	// verify if the msg is signed by sender.
	var payload []byte
	payload, err := m.PayloadNoSig()
	if err != nil {
		return err
	}

	signer, err := types.GetSignatureAddress(payload, m.Signature)
	if err != nil {
		return err
	}

	if !bytes.Equal(m.Address.Bytes(), signer.Bytes()) {
		return errWrongSignatureMsg
	}

	// message is signed by the signer, now it is accountable for msg of proposal and votes.
	if m.Code > proto.MsgPrecommit {
		return errAccountableGarbageMsg
	}

	// then try to decode the tendermint msg bytes to construct msg height, round, step, etc...
	if m.Code == proto.MsgProposal {
		var proposal mUtils.Proposal
		err = m.Decode(&proposal)
		if err != nil {
			return errAccountableGarbageMsg
		}
	} else if m.Code == proto.MsgPrevote || m.Code == proto.MsgPrecommit {
		var vote mUtils.Vote
		err := m.Decode(&vote)
		if err != nil {
			return errAccountableGarbageMsg
		}
	}

	return nil
}

func preCheckMessage(m *mUtils.Message, chain BlockChainContext) error {
	lastHeader := chain.GetHeaderByNumber(m.H() - 1)
	if lastHeader == nil {
		return errFutureMsg
	}

	v := lastHeader.CommitteeMember(m.Address)
	if v == nil {
		return errNotCommitteeMsg
	}

	m.Power = v.VotingPower
	return nil
}

func (fd *FaultDetector) msgHandlerLoop() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
tendermintMsgLoop:
	for {
		curHeight := fd.blockchain.CurrentBlock().Number().Uint64()
		curHeader := fd.blockchain.CurrentHeader()
		select {
		case ev, ok := <-fd.tendermintMsgSub.Chan():
			if !ok {
				break tendermintMsgLoop
			}

			// handle consensus msg or innocence proof msgs
			switch e := ev.Data.(type) {
			case events.MessageEvent:
				// decode msg from payload to construct msg code, tendermint msg bytes, sender address, committed seal and signature.
				msg := new(mUtils.Message)
				msg.Bytes = e.Payload
				err := rlp.DecodeBytes(msg.Bytes, msg)
				if err != nil {
					continue tendermintMsgLoop
				}

				err = decodeConsensusMsg(msg)
				if err != nil {
					// make this fault accountable only for committee members, otherwise validators might pay fees to
					// report lots of none sense proof which is a vector of attack as well.
					if err == errAccountableGarbageMsg && curHeader.CommitteeMember(msg.Address) != nil {
						fd.submitMisbehavior(msg, nil, errAccountableGarbageMsg, fd.misbehaviourProofsCh)
					}
					continue tendermintMsgLoop
				}

				if fd.tooOldHeightMsg(curHeight, msg.H()) {
					fd.logger.Info("fault detector: discarding old message", "sender", msg.Sender())
					continue tendermintMsgLoop
				}

				if err := fd.processMsg(msg); err != nil {
					fd.logger.Info("fault detector: error while processing consensus msg", "err", err)
					continue tendermintMsgLoop
				}
			case events.AccountabilityOffChainEvent:
				err := fd.handleOffChainAccountabilityEvent(e.Payload, e.Sender)
				if err != nil {
					fd.logger.Info("going to drop peer", "peer", e.Sender)
					// the errors return from handler could freeze the peer connection for 30 seconds by according to dev p2p protocol.
					select {
					case e.ErrCh <- err:
					default: // do nothing
					}
					continue tendermintMsgLoop
				}
			}
		case e, ok := <-fd.msgHandlerBlockCh:
			if !ok {
				break tendermintMsgLoop
			}

			// on every 60 blocks, reset Peer Justified Accusations and height accusations counters.
			if e.Block.NumberU64()%msgGCInterval == 0 {
				fd.rateLimiter.resetHeightRateLimiter()
				fd.rateLimiter.resetPeerJustifiedAccusations()
			}

			for h, msgs := range fd.futureHeightMsgBuffer {
				if h <= curHeight {
					for _, m := range msgs {
						if err := fd.processMsg(m); err != nil {
							fd.logger.Error("fault detector: error while processing consensus msg", "err", err)
						}
					}
					// once messages are processed, delete it from buffer.
					fd.deleteFutureHeightMsg(h)
				}
			}
		case <-ticker.C:
			// on each 1 seconds, reset the rate limiter counters.
			fd.rateLimiter.resetRateLimiter()
		case err, ok := <-fd.msgHandlerBlockSub.Err():
			if ok {
				fd.logger.Crit("block subscription error", err.Error())
			}
			break tendermintMsgLoop
		}
	}
	close(fd.misbehaviourProofsCh)
}

// check to GC msg store for those msgs out of buffering window on every 60 blocks.
func (fd *FaultDetector) checkMsgStoreGC(height uint64) {

	if height > proto.AccountabilityHeightRange {
		if height%msgGCInterval == 0 {
			gcHeight := height - proto.AccountabilityHeightRange
			fd.msgStore.DeleteMsgsBeforeHeight(gcHeight)
		}
	}
}

func (fd *FaultDetector) ruleEngineLoop() {
blockChainLoop:
	for {
		select {
		// chain event update, provide proof of innocent if one is on challenge, rule engine scanning is triggered also.
		case ev, ok := <-fd.ruleEngineBlockCh:
			if !ok {
				break blockChainLoop
			}

			// handle accusations and provide innocence AccountabilityProof if there were any for a node.
			innocenceProofs, err := fd.handleAccusations(ev.Block)
			if err != nil {
				fd.logger.Crit("handleAccusation", "fault detector", err)
			}

			if innocenceProofs != nil {
				// send on chain innocence proof ASAP since the client is on challenge that requires the proof to be
				// provided before the client get slashed.
				fd.sendAccountabilityEvents(innocenceProofs, false)
			}

			// try to escalate expired off chain accusation on chain.
			fd.escalateExpiredOffChainAccusation(ev.Block.NumberU64())

			// run rule engine over a specific height.
			if ev.Block.NumberU64() > uint64(proto.DeltaBlocks) {
				checkPointHeight := ev.Block.NumberU64() - uint64(proto.DeltaBlocks)
				proofs := fd.runRuleEngine(checkPointHeight)
				if len(proofs) > 0 {
					fd.accountabilityEventBuffer = append(fd.accountabilityEventBuffer, proofs...)
				}

				// send accountability event base on reporting slots.
				if len(fd.accountabilityEventBuffer) != 0 {
					fd.slotBasedAccountabilityEventReporting(checkPointHeight)
				}
			}

			// msg store delete msgs out of buffering window on every 60 blocks.
			fd.checkMsgStoreGC(ev.Block.NumberU64())

		case m, ok := <-fd.misbehaviourProofsCh:
			if !ok {
				break blockChainLoop
			}
			fd.accountabilityEventBuffer = append(fd.accountabilityEventBuffer, m)
		case err, ok := <-fd.ruleEngineBlockSub.Err():
			if ok {
				fd.logger.Crit("block subscription error", err.Error())
			}
			break blockChainLoop
		}
	}
}

func (fd *FaultDetector) slotBasedAccountabilityEventReporting(checkPointHeight uint64) {
	// due to the committee rotation of per epoch, the slot management should base on the committee
	// falls behind delta block of current chain height.
	lastHeader := fd.blockchain.GetHeaderByNumber(checkPointHeight - 1)
	if lastHeader == nil {
		return
	}
	committee := lastHeader.Committee

	// each reporting slot contains 5 block period that a unique and deterministic validator is asked to
	// be the reporter of that slot period, then at the end block of that slot, the reporter reports
	// available events. Thus, between each reporting slot, we have 5 block period to wait for
	// accountability events to be mined by network, and it is also disaster friendly that if the last
	// reporter fails, the next reporter will continue to report missing events.
	reporterIndex := (checkPointHeight / proto.ReportingSlotPeriod) % uint64(len(committee))

	// if validator is the reporter of the slot period, and if checkpoint block is the end block of the
	// slot, then it is time to report the collected events by this validator.
	if committee[reporterIndex].Address == fd.address && (checkPointHeight+1)%proto.ReportingSlotPeriod == 0 {
		fd.logger.Info("----> On the slot to report accountability event", "number of events", len(fd.accountabilityEventBuffer))
		copyOnChainProofs := make([]*autonity.AccountabilityEvent, len(fd.accountabilityEventBuffer))
		copy(copyOnChainProofs, fd.accountabilityEventBuffer)
		fd.sendAccountabilityEvents(copyOnChainProofs, true)
		// release events once events were sent.
		fd.accountabilityEventBuffer = nil
	}
}

func (fd *FaultDetector) Stop() {
	fd.ruleEngineBlockSub.Unsubscribe()
	fd.msgHandlerBlockSub.Unsubscribe()
	fd.accountabilityTXEventSub.Unsubscribe()
	fd.tendermintMsgSub.Unsubscribe()
	fd.proofWG.Wait()
	unRegisterFaultDetectorContracts()
}

func (fd *FaultDetector) filterProcessedOnes(proofs []*autonity.AccountabilityEvent) (result []*autonity.AccountabilityEvent) {
	// get latest chain state.
	stateDB, err := fd.blockchain.State()
	if err != nil {
		return nil
	}
	header := fd.blockchain.CurrentBlock().Header()
	contract := fd.blockchain.GetAutonityContract()

	for _, e := range proofs {
		if e.Type == uint8(autonity.Misbehaviour) {
			if contract.MisbehaviourProcessed(header, stateDB, e.MsgHash) {
				continue
			}
		} else if e.Type == uint8(autonity.Accusation) {
			if contract.AccusationProcessed(header, stateDB, e.MsgHash) {
				continue
			}
		}
		result = append(result, e)
	}
	return result
}

// convert the raw proofs into on-chain AccountabilityProof which contains raw bytes of messages.
func (fd *FaultDetector) generateOnChainProof(p *AccountabilityProof) (*autonity.AccountabilityEvent, error) {
	var ev = &autonity.AccountabilityEvent{
		Type:     uint8(p.Type),
		Rule:     uint8(p.Rule),
		Reporter: fd.address,
		Sender:   p.Message.Address,
		MsgHash:  p.Message.MsgHash(),
	}

	rProof, err := rlp.EncodeToBytes(p)
	if err != nil {
		return nil, err
	}
	fd.logger.Info("gen accountability proof", "type", p.Type, "rule", p.Rule, "msg sender",
		p.Message.Address)
	ev.RawProof = rProof
	return ev, nil
}

// getInnocentProof is called by client who is on a challenge with a certain accusation, to get innocent proof from msg
// store.
func (fd *FaultDetector) getInnocentProof(acProof *AccountabilityProof) (*autonity.AccountabilityEvent, error) {
	var onChainEvent *autonity.AccountabilityEvent
	// the protocol contains below provable accusations.
	switch acProof.Rule {
	case autonity.PO:
		return fd.getInnocentProofOfPO(acProof)
	case autonity.PVN:
		return fd.getInnocentProofOfPVN(acProof)
	case autonity.PVO:
		return fd.getInnocentProofOfPVO(acProof)
	case autonity.C1:
		return fd.getInnocentProofOfC1(acProof)
	default:
		return onChainEvent, fmt.Errorf("not provable rule")
	}
}

// get innocent proof of accusation of rule C1 from msg store.
func (fd *FaultDetector) getInnocentProofOfC1(c *AccountabilityProof) (*autonity.AccountabilityEvent, error) {
	preCommit := c.Message
	height := preCommit.H()

	lastHeader := fd.blockchain.GetHeaderByNumber(height - 1)
	if lastHeader == nil {
		return nil, errNoParentHeader
	}
	quorum := bft.Quorum(lastHeader.TotalVotingPower())

	prevotesForV := fd.msgStore.Get(height, func(m *mUtils.Message) bool {
		return m.Type() == proto.MsgPrevote && m.Value() == preCommit.Value() && m.R() == preCommit.R()
	})

	overQuorumVotes := engineCore.OverQuorumVotes(prevotesForV, quorum.Uint64())
	if overQuorumVotes == nil {
		return nil, errNoEvidenceForC1
	}

	p, err := fd.generateOnChainProof(&AccountabilityProof{
		Type:     autonity.Innocence,
		Rule:     c.Rule,
		Message:  preCommit,
		Evidence: overQuorumVotes,
	})
	return p, err
}

// get innocent proof of accusation of rule PO from msg store.
func (fd *FaultDetector) getInnocentProofOfPO(ac *AccountabilityProof) (*autonity.AccountabilityEvent, error) {
	// PO: node propose an old value with an validRound, innocent onChainProof of it should be:
	// there are quorum num of prevote for that value at the validRound.
	liteProposal := ac.Message
	height := liteProposal.H()
	validRound := liteProposal.ValidRound()
	lastHeader := fd.blockchain.GetHeaderByNumber(height - 1)
	if lastHeader == nil {
		return nil, errNoParentHeader
	}
	quorum := bft.Quorum(lastHeader.TotalVotingPower())

	prevotes := fd.msgStore.Get(height, func(m *mUtils.Message) bool {
		return m.Type() == proto.MsgPrevote && m.R() == validRound && m.Value() == liteProposal.Value()
	})

	overQuorumPreVotes := engineCore.OverQuorumVotes(prevotes, quorum.Uint64())
	if overQuorumPreVotes == nil {
		// cannot onChainProof its innocent for PO, the on-chain contract will fine it latter once the
		// time window for onChainProof ends.
		return nil, errNoEvidenceForPO
	}

	p, err := fd.generateOnChainProof(&AccountabilityProof{
		Type:     autonity.Innocence,
		Rule:     ac.Rule,
		Message:  liteProposal,
		Evidence: overQuorumPreVotes,
	})
	return p, err
}

// get innocent proof of accusation of rule PVN from msg store.
func (fd *FaultDetector) getInnocentProofOfPVN(c *AccountabilityProof) (*autonity.AccountabilityEvent, error) {
	// get innocent proofs for PVN, for a prevote that vote for a new value,
	// then there must be a proposal for this new value.
	prevote := c.Message
	height := prevote.H()

	// the only proof of innocence of PVN accusation is that there exist a corresponding proposal
	proposals := fd.msgStore.Get(height, func(m *mUtils.Message) bool {
		return m.Type() == proto.MsgProposal && m.R() == prevote.R() && m.Value() == prevote.Value() && m.ValidRound() == -1
	})

	var ev []*mUtils.Message
	if len(proposals) != 0 {
		ev = []*mUtils.Message{proposals[0].ToLiteProposal()}
		p, err := fd.generateOnChainProof(&AccountabilityProof{
			Type:     autonity.Innocence,
			Rule:     c.Rule,
			Message:  prevote,
			Evidence: ev,
		})
		return p, err
	}
	return nil, errNoEvidenceForPVN
}

// get innocent proof of accusation of rule PVO from msg store, it collects quorum preVotes for the value voted at a valid round.
func (fd *FaultDetector) getInnocentProofOfPVO(c *AccountabilityProof) (*autonity.AccountabilityEvent, error) {
	// get innocent proofs for PVO, collect quorum preVotes at the valid round of the old proposal.
	oldProposal := c.Evidence[0]
	height := oldProposal.H()
	validRound := oldProposal.ValidRound()
	lastHeader := fd.blockchain.GetHeaderByNumber(height - 1)
	if lastHeader == nil {
		return nil, errNoParentHeader
	}
	quorum := bft.Quorum(lastHeader.TotalVotingPower())

	preVotes := fd.msgStore.Get(height, func(m *mUtils.Message) bool {
		return m.Type() == proto.MsgPrevote && m.Value() == oldProposal.Value() && m.R() == validRound
	})

	overQuorumVotes := engineCore.OverQuorumVotes(preVotes, quorum.Uint64())

	if overQuorumVotes == nil {
		return nil, errNoEvidenceForPVO
	}

	p, err := fd.generateOnChainProof(&AccountabilityProof{
		Type:     autonity.Innocence,
		Rule:     c.Rule,
		Message:  c.Message,
		Evidence: append(c.Evidence, overQuorumVotes...),
	})
	return p, err
}

// handle on-chain challenges, and provide innocent proof for accusations if there exist proof in msg store.
func (fd *FaultDetector) handleAccusations(block *types.Block) ([]*autonity.AccountabilityEvent, error) {
	var innocentProofs []*autonity.AccountabilityEvent // nolint
	stateDB, err := fd.blockchain.StateAt(block.Root())
	if err != nil || stateDB == nil {
		return nil, err
	}

	contract := fd.blockchain.GetAutonityContract()
	if contract == nil {
		return nil, fmt.Errorf("cannot get contract instance")
	}

	accusations := contract.GetValidatorAccusations(block.Header(), stateDB, fd.address)
	for _, a := range accusations {
		b := a
		var acProof *AccountabilityProof
		if b.Chunks != 0 {
			acProof, err = fd.decodeChunkedAccountabilityEvent(contract, stateDB, block.Header(), &b)
			if err != nil {
				continue
			}
		} else {
			acProof, err = decodeRawProof(b.RawProof)
			if err != nil {
				continue
			}
		}
		p, err := fd.getInnocentProof(acProof)
		if err != nil {
			continue
		}
		innocentProofs = append(innocentProofs, p)
	}

	return innocentProofs, nil
}

func (fd *FaultDetector) decodeChunkedAccountabilityEvent(contract *autonity.Contract, state *state.StateDB, header *types.Header, ev *autonity.AccountabilityEvent) (*AccountabilityProof, error) {

	var constructedBytes []byte
	for chunkID := uint8(0); chunkID < ev.Chunks; chunkID++ {
		chunk, err := contract.GetAccountabilityEventChunk(header, state, ev.MsgHash, ev.Type, ev.Rule, ev.Reporter, chunkID)
		if err != nil {
			return nil, err
		}
		constructedBytes = append(constructedBytes, chunk...)
	}

	ac, err := decodeRawProof(constructedBytes)
	if err != nil {
		return nil, err
	}

	return ac, nil
}

// processMsg, check and submit any auto-incriminating, equivocation challenges, and then only store checked msg in msg store.
func (fd *FaultDetector) processMsg(m *mUtils.Message) error {
	// check if msg is from valid committee member
	err := preCheckMessage(m, fd.blockchain)
	if err != nil {
		if err == errFutureMsg {
			fd.bufferFutureHeightMsg(m)
		}
		return err
	}

	switch m.Code {
	case proto.MsgProposal:
		err := fd.accountForAutoIncriminatingProposal(m)
		if err != nil {
			return err
		}
	case proto.MsgPrevote:
		fallthrough
	case proto.MsgPrecommit:
		err := fd.accountForAutoIncriminatingVote(m)
		if err != nil {
			return err
		}
	default:
		// shouldn't happen since the consensus msg decoder address this as an accountable garbage msg.
		return fmt.Errorf("unknown msg code")
	}

	// msg pass the auto-incriminating checker, save it in msg store.
	fd.msgStore.Save(m)
	return nil
}

// run rule engine over the specific height of consensus msgs, return the accountable events in proofs.
func (fd *FaultDetector) runRuleEngine(checkPointHeight uint64) []*autonity.AccountabilityEvent {
	var onChainProofs []*autonity.AccountabilityEvent
	// To avoid none necessary accusations, we wait for delta blocks to start rule scan.
	// always skip the heights before first buffered height after the node start up, since it will rise lots of none
	// sense accusations due to the missing of messages during the startup phase, it cost un-necessary payments
	// for the committee member.
	if checkPointHeight <= fd.msgStore.FirstHeightBuffered() {
		return nil
	}
	lastHeader := fd.blockchain.GetHeaderByNumber(checkPointHeight - 1)
	if lastHeader == nil {
		return nil
	}
	quorum := bft.Quorum(lastHeader.TotalVotingPower())
	proofs := fd.runRulesOverHeight(checkPointHeight, quorum.Uint64())
	if len(proofs) > 0 {
		for _, proof := range proofs {
			if proof.Message.Address == fd.address {
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

			p, err := fd.generateOnChainProof(proof)
			if err != nil {
				fd.logger.Warn("convert AccountabilityProof to on-chain AccountabilityProof", "fault detector", err)
				continue
			}
			onChainProofs = append(onChainProofs, p)
		}
	}
	return onChainProofs
}

func (fd *FaultDetector) runRulesOverHeight(height uint64, quorum uint64) (proofs []*AccountabilityProof) {
	// Rules read right to left (find the right and look for the left)
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
	proofs = append(proofs, fd.invalidProposalAccountabilityCheck(height, quorum)...)
	proofs = append(proofs, fd.newProposalsAccountabilityCheck(height)...)
	proofs = append(proofs, fd.oldProposalsAccountabilityCheck(height, quorum)...)
	proofs = append(proofs, fd.prevotesAccountabilityCheck(height, quorum)...)
	proofs = append(proofs, fd.precommitsAccountabilityCheck(height, quorum)...)
	return proofs
}

// invalidProposalAccountabilityCheck checks the proposal without by executing it from evm, it is depends on if there are
// quorum prevotes targets the proposal as a bad proposal.
func (fd *FaultDetector) invalidProposalAccountabilityCheck(height uint64, quorum uint64) (proofs []*AccountabilityProof) {
	proposals := fd.msgStore.Get(height, func(m *mUtils.Message) bool {
		return m.Type() == proto.MsgProposal
	})

	for _, p := range proposals {
		proposal := p

		// Skip if proposal is equivocated
		proposalsForR := fd.msgStore.Get(height, func(m *mUtils.Message) bool {
			return m.Sender() == proposal.Sender() && m.Type() == proto.MsgProposal && m.R() == proposal.R()

		})
		// Due to the for loop there must be at least one proposal
		if len(proposalsForR) > 1 {
			continue
		}

		badProposalVotes := fd.msgStore.Get(height, func(m *mUtils.Message) bool {
			return m.Type() == proto.MsgPrevote && m.Value() == nilValue && m.BadValue() == proposal.Value() &&
				m.BadProposer() == proposal.Sender() && m.R() == proposal.R()
		})

		quorumBadVotes := engineCore.OverQuorumVotes(badProposalVotes, quorum)
		if quorumBadVotes == nil {
			continue
		}

		proof := &AccountabilityProof{
			Type:     autonity.Misbehaviour,
			Rule:     autonity.InvalidProposal,
			Evidence: quorumBadVotes,
			Message:  proposal.ToLiteProposal(),
		}
		proofs = append(proofs, proof)
	}
	return proofs
}

func (fd *FaultDetector) newProposalsAccountabilityCheck(height uint64) (proofs []*AccountabilityProof) {
	// ------------New Proposal------------
	// PN:  (Mr′<r,P C|pi)∗ <--- (Mr,P|pi)
	// PN1: [nil ∨ ⊥] <--- [V]
	//
	// Since the message pattern for PN includes only messages sent by pi, we cannot raise an accusation. We can only
	// raise a misbehaviour. To raise a misbehaviour for PN1 we need to have received all the precommits from pi for all
	// r' < r. If any of the precommits is for a non-nil value then we have proof of misbehaviour.

	proposalsNew := fd.msgStore.Get(height, func(m *mUtils.Message) bool {
		return m.Type() == proto.MsgProposal && m.ValidRound() == -1
	})

	for _, p := range proposalsNew {
		proposal := p

		// Skip if proposal is equivocated
		proposalsForR := fd.msgStore.Get(height, func(m *mUtils.Message) bool {
			return m.Sender() == proposal.Sender() && m.Type() == proto.MsgProposal && m.R() == proposal.R()

		})
		// Due to the for loop there must be at least one proposal
		if len(proposalsForR) > 1 {
			continue
		}

		//check all precommits for previous rounds from this sender are nil
		precommits := fd.msgStore.Get(height, func(m *mUtils.Message) bool {
			return m.Sender() == proposal.Sender() && m.Type() == proto.MsgPrecommit && m.R() < proposal.R() && m.Value() != nilValue
		})
		if len(precommits) != 0 {
			proof := &AccountabilityProof{
				Type:     autonity.Misbehaviour,
				Rule:     autonity.PN,
				Evidence: precommits[0:1],
				Message:  proposal.ToLiteProposal(),
			}
			proofs = append(proofs, proof)
			fd.logger.Info("Misbehaviour detected", "fault detector", fd.address, "rulePN", autonity.PN, "sender", proposal.Sender())
		}
	}
	return proofs
}

func (fd *FaultDetector) oldProposalsAccountabilityCheck(height uint64, quorum uint64) (proofs []*AccountabilityProof) {
	// ------------Old Proposal------------
	// PO: (Mr′<r,PV) ∧ (Mr′,PC|pi) ∧ (Mr′<r′′<r,P C|pi)∗ <--- (Mr,P|pi)
	// PO1: [#(Mr′,PV|V) ≥ 2f+ 1] ∧ [nil ∨ V ∨ ⊥] ∧ [nil ∨ ⊥] <--- [V]

	proposalsOld := fd.msgStore.Get(height, func(m *mUtils.Message) bool {
		return m.Type() == proto.MsgProposal && m.ValidRound() > -1
	})

oldProposalLoop:
	for _, p := range proposalsOld {
		proposal := p
		// Check that in the valid round we see a quorum of prevotes and that there is no precommit at all or a
		// precommit for v or nil.

		// Skip if proposal is equivocated
		proposalsForR := fd.msgStore.Get(height, func(m *mUtils.Message) bool {
			return m.Sender() == proposal.Sender() && m.Type() == proto.MsgProposal && m.R() == proposal.R()

		})
		// Due to the for loop there must be at least one proposal
		if len(proposalsForR) > 1 {
			continue oldProposalLoop
		}

		validRound := proposal.ValidRound()

		// Is there a precommit for a value other than nil or the proposed value by the current proposer in the valid
		// round? If there is, the proposer has proposed a value for which it is not locked on, thus a AccountabilityProof of
		// misbehaviour can be generated.
		precommitsFromPiInVR := fd.msgStore.Get(height, func(m *mUtils.Message) bool {
			return m.Type() == proto.MsgPrecommit && m.R() == validRound && m.Sender() == proposal.Sender() &&
				m.Value() != nilValue && m.Value() != proposal.Value()
		})
		if len(precommitsFromPiInVR) > 0 {
			proof := &AccountabilityProof{
				Type:     autonity.Misbehaviour,
				Rule:     autonity.PO,
				Evidence: precommitsFromPiInVR[0:1],
				Message:  proposal.ToLiteProposal(),
			}
			proofs = append(proofs, proof)
			fd.logger.Info("Misbehaviour detected", "fault detector", fd.address, "rulePO", autonity.PO, "sender", proposal.Sender())
			continue oldProposalLoop
		}

		// Is there a precommit for anything other than nil from the proposer between the valid round and the round of
		// the proposal? If there is then that implies the proposer saw 2f+1 prevotes in that round and hence it should
		// have set that round as the valid round.
		precommitsFromPiAfterVR := fd.msgStore.Get(height, func(m *mUtils.Message) bool {
			return m.Type() == proto.MsgPrecommit && m.R() > validRound && m.R() < proposal.R() &&
				m.Sender() == proposal.Sender() && m.Value() != nilValue
		})
		if len(precommitsFromPiAfterVR) > 0 {
			proof := &AccountabilityProof{
				Type:     autonity.Misbehaviour,
				Rule:     autonity.PO,
				Evidence: precommitsFromPiAfterVR[0:1],
				Message:  proposal.ToLiteProposal(),
			}
			proofs = append(proofs, proof)
			fd.logger.Info("Misbehaviour detected", "fault detector", fd.address, "rulePO", autonity.PO, "sender", proposal.Sender())
			continue oldProposalLoop
		}

		// Do we see a quorum for a value other than the proposed value? If so, we have proof of misbehaviour.
		allPrevotesForValidRound := fd.msgStore.Get(height, func(m *mUtils.Message) bool {
			return m.Type() == proto.MsgPrevote && m.R() == validRound && m.Value() != proposal.Value()
		})

		prevotesMap := make(map[common.Hash][]*mUtils.Message)
		for _, p := range allPrevotesForValidRound {
			prevotesMap[p.Value()] = append(prevotesMap[p.Value()], p)
		}

		for _, preVotes := range prevotesMap {
			// Here the assumption is that in a single round it is not possible to have 2 value which quorum votes,
			// this would imply at least quorum nodes are malicious which is much higher than our assumption.
			overQuorumVotes := engineCore.OverQuorumVotes(preVotes, quorum)
			if overQuorumVotes != nil {
				proof := &AccountabilityProof{
					Type:     autonity.Misbehaviour,
					Rule:     autonity.PO,
					Evidence: overQuorumVotes,
					Message:  proposal.ToLiteProposal(),
				}
				proofs = append(proofs, proof)
				fd.logger.Info("Misbehaviour detected", "fault detector", fd.address, "rulePO", autonity.PO, "sender", proposal.Sender())
				continue oldProposalLoop
			}
		}

		// Do we see a quorum of prevotes in the valid round, if not we can raise an accusation, since we cannot be sure
		// that these prevotes don't exist
		prevotes := fd.msgStore.Get(height, func(m *mUtils.Message) bool {
			// since equivocation msgs are stored, we have to query those preVotes which has same value as the proposal.
			return m.Type() == proto.MsgPrevote && m.R() == validRound && m.Value() == proposal.Value()
		})

		if engineCore.OverQuorumVotes(prevotes, quorum) == nil {
			accusation := &AccountabilityProof{
				Type:    autonity.Accusation,
				Rule:    autonity.PO,
				Message: proposal.ToLiteProposal(),
			}
			proofs = append(proofs, accusation)
			fd.logger.Info("Accusation detected", "fault detector", fd.address, "rulePO", autonity.PO, "sender", proposal.Sender())
		}
	}
	return proofs
}

func (fd *FaultDetector) prevotesAccountabilityCheck(height uint64, quorum uint64) (proofs []*AccountabilityProof) {
	// ------------New and Old Prevotes------------

	prevotes := fd.msgStore.Get(height, func(m *mUtils.Message) bool {
		return m.Type() == proto.MsgPrevote && m.Value() != nilValue
	})

prevotesLoop:
	for _, p := range prevotes {
		prevote := p

		// Skip if prevote is equivocated
		prevotesForR := fd.msgStore.Get(height, func(m *mUtils.Message) bool {
			return m.Sender() == prevote.Sender() && m.Type() == proto.MsgPrevote && m.R() == prevote.R()

		})
		// Due to the for loop there must be at least one preVote.
		if len(prevotesForR) > 1 {
			continue prevotesLoop
		}

		// We need to check whether we have proposals from the prevote's round
		correspondingProposals := fd.msgStore.Get(height, func(m *mUtils.Message) bool {
			return m.Type() == proto.MsgProposal && m.R() == prevote.R() && m.Value() == prevote.Value()
		})

		if len(correspondingProposals) == 0 {

			// if there are over quorum prevotes for this corresponding proposal's value, then it indicates current
			// peer just did not receive it. So we can skip the rising of such accusation.
			preVts := fd.msgStore.Get(height, func(m *mUtils.Message) bool {
				return m.Type() == proto.MsgPrevote && m.R() == prevote.R() && m.Value() == prevote.Value()
			})

			if engineCore.OverQuorumVotes(preVts, quorum) == nil {
				// The rule for this accusation could be PVO as well since we don't have the corresponding proposal, but
				// it does not mean it's incorrect. More over that, since over quorum prevotes at the round
				// of correspondingProposals are used as the innocence proof, rather than the correspondingProposals, thus
				// we don't worry that the correspondingProposals sender could lie on the proof providing phase.
				accusation := &AccountabilityProof{
					Type:    autonity.Accusation,
					Rule:    autonity.PVN,
					Message: prevote,
				}
				proofs = append(proofs, accusation)
				fd.logger.Info("Accusation detected", "fault detector", fd.address, "rulePVN", autonity.PVN, "sender", prevote.Sender())
				continue prevotesLoop
			}
		}

		// We need to ensure that we keep all proposals in the message store, so that we have the maximum chance of
		// finding justification for prevotes. This is to account for equivocation where the proposer send 2 proposals
		// with the same value but different valid rounds to different nodes. We can't penalise the sender of prevote
		// since we can't tell which proposal they received. We just want to find a set of message which fit the rule.
		// Therefore, we need to check all the proposals to find a single one which shows the current prevote is
		// valid.
		var prevotesProofs []*AccountabilityProof
		for _, cp := range correspondingProposals {
			correspondingProposal := cp
			if correspondingProposal.ValidRound() == -1 {
				prevotesProofs = append(prevotesProofs, fd.newPrevotesAccountabilityCheck(height, prevote, correspondingProposal))
			} else {
				prevotesProofs = append(prevotesProofs, fd.oldPrevotesAccountabilityCheck(height, quorum, correspondingProposal, prevote))
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

func (fd *FaultDetector) newPrevotesAccountabilityCheck(height uint64, prevote *mUtils.Message, correspondingProposal *mUtils.Message) (proof *AccountabilityProof) {
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
	precommitsFromPi := fd.msgStore.Get(height, func(m *mUtils.Message) bool {
		return m.Type() == proto.MsgPrecommit && prevote.Sender() == m.Sender() && m.R() < prevote.R()
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
				precommitsAtRPrime := fd.msgStore.Get(height, func(m *mUtils.Message) bool {
					return m.Type() == proto.MsgPrecommit && pc.Sender() == m.Sender() && m.R() == pc.R()
				})

				// Check for equivocation, it is possible there are multiple precommit from pi for the same round.
				// If there are equivocated messages: do nothing. Since pi has already being punished for equivocation
				// round when the equivocated message was first received.
				if len(precommitsAtRPrime) == 1 {
					if precommitsAtRPrime[0].Value() != prevote.Value() {
						fd.logger.Info("Misbehaviour detected", "fault detector", fd.address, "rulePVN", autonity.PVN, "sender", prevote.Sender())
						proof := &AccountabilityProof{
							Type:    autonity.Misbehaviour,
							Rule:    autonity.PVN,
							Message: prevote,
						}
						// to guarantee this prevote is for a new proposal that is the PVN rule account for, otherwise in
						// prevote for an old proposal, it is valid for one to prevote it if lockedRound <= vr, thus the
						// round gump is valid. This prevents from rising a PVN misbehavior proof from a malicious fault
						// detector by using prevote for an old proposal to challenge an honest slow validator.
						proof.Evidence = append(proof.Evidence, correspondingProposal.ToLiteProposal())
						proof.Evidence = append(proof.Evidence, precommitsFromPi[i:]...)
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

func (fd *FaultDetector) oldPrevotesAccountabilityCheck(height uint64, quorum uint64, correspondingProposal *mUtils.Message, prevote *mUtils.Message) (proof *AccountabilityProof) {
	currentR := correspondingProposal.R()
	validRound := correspondingProposal.ValidRound()

	// If one prevotes for an invalid old proposal, then it should be a misbehaviour.
	if validRound >= currentR {
		fd.logger.Info("Misbehaviour detected", "fault detector", fd.address, "rulePVO3", autonity.PVO3, "sender", prevote.Sender())
		proof := &AccountabilityProof{
			Type:    autonity.Misbehaviour,
			Rule:    autonity.PVO3,
			Message: prevote,
		}
		proof.Evidence = append(proof.Evidence, correspondingProposal.ToLiteProposal())
		return proof
	}

	// If there is a prevote for an old proposal then pi can only vote for v or send nil (see line 28 and 29 of
	// tendermint pseudocode), therefore if in the valid round there is a quorum for a value other than v, we know pi
	// prevoted incorrectly. If the proposal was a bad proposal, then pi should not have voted for it, thus we do not
	// need to make sure whether the proposal is correct or not (which we would in the proposal checking rules, however,
	// a bad proposal will still exist in our message store, and it shouldn't have an impact on the checking of prevotes).

	allPrevotesForValidRound := fd.msgStore.Get(height, func(m *mUtils.Message) bool {
		return m.Type() == proto.MsgPrevote && m.R() == validRound && m.Value() != correspondingProposal.Value()
	})

	prevotesMap := make(map[common.Hash][]*mUtils.Message)
	for _, p := range allPrevotesForValidRound {
		prevotesMap[p.Value()] = append(prevotesMap[p.Value()], p)
	}

	for _, preVotes := range prevotesMap {
		// Here the assumption is that in a single round it is not possible to have 2 value which quorum votes,
		// this would imply at least quorum nodes are malicious which is much higher than our assumption.
		overQuorumVotes := engineCore.OverQuorumVotes(preVotes, quorum)
		if overQuorumVotes != nil {
			fd.logger.Info("Misbehaviour detected", "fault detector", fd.address, "rulePVO", autonity.PVO, "sender", prevote.Sender())
			proof := &AccountabilityProof{
				Type:    autonity.Misbehaviour,
				Rule:    autonity.PVO,
				Message: prevote,
			}
			proof.Evidence = append(proof.Evidence, correspondingProposal.ToLiteProposal())
			proof.Evidence = append(proof.Evidence, overQuorumVotes...)
			return proof
		}
	}

	prevotesForVFromValidRound := fd.msgStore.Get(height, func(m *mUtils.Message) bool {
		return m.Type() == proto.MsgPrevote && m.R() == validRound && m.Value() == correspondingProposal.Value()
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

		precommitsFromPi := fd.msgStore.Get(height, func(m *mUtils.Message) bool {
			return m.Type() == proto.MsgPrecommit && m.R() > validRound && m.R() < currentR && m.Sender() == prevote.Sender()
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
				fd.logger.Info("Misbehaviour detected", "fault detector", fd.address, "rulePVO12", autonity.PVO12, "sender", prevote.Sender())
				proof := &AccountabilityProof{
					Type:    autonity.Misbehaviour,
					Rule:    autonity.PVO12,
					Message: prevote,
				}
				proof.Evidence = append(proof.Evidence, correspondingProposal.ToLiteProposal())
				proof.Evidence = append(proof.Evidence, precommitsFromPi...)
				return proof
			}
		}
	}

	// if there is no misbehaviour of the prevote msg addressed, then we lastly check accusation.
	if overQuorumPrevotesForVFromValidRound == nil {
		// raise an accusation
		fd.logger.Info("Accusation detected", "fault detector", fd.address, "rulePVO", autonity.PVO, "sender", prevote.Sender())
		return &AccountabilityProof{
			Type:     autonity.Accusation,
			Rule:     autonity.PVO,
			Message:  prevote,
			Evidence: []*mUtils.Message{correspondingProposal.ToLiteProposal()},
		}
	}

	return nil
}

func (fd *FaultDetector) precommitsAccountabilityCheck(height uint64, quorum uint64) (proofs []*AccountabilityProof) {
	// ------------Precommits------------
	// C: [Mr,P|proposer(r)] ∧ [Mr,PV] <--- [Mr,PC|pi]
	// C1: [V:Valid(V)] ∧ [#(V) ≥ 2f+ 1] <--- [V]

	precommits := fd.msgStore.Get(height, func(m *mUtils.Message) bool {
		return m.Type() == proto.MsgPrecommit && m.Value() != nilValue
	})

precommitLoop:
	for _, preC := range precommits {
		precommit := preC

		// Skip if preCommit is equivocated
		precommitsForR := fd.msgStore.Get(height, func(m *mUtils.Message) bool {
			return m.Sender() == precommit.Sender() && m.Type() == proto.MsgPrecommit && m.R() == precommit.R()

		})
		// Due to the for loop there must be at least one preCommit.
		if len(precommitsForR) > 1 {
			continue precommitLoop
		}

		// Do we see a quorum for a value other than the proposed value? If so, we have proof of misbehaviour.
		allPrevotesForR := fd.msgStore.Get(height, func(m *mUtils.Message) bool {
			return m.Type() == proto.MsgPrevote && m.R() == precommit.R() && m.Value() != precommit.Value()
		})

		prevotesMap := make(map[common.Hash][]*mUtils.Message)
		for _, p := range allPrevotesForR {
			prevotesMap[p.Value()] = append(prevotesMap[p.Value()], p)
		}

		for _, preVotes := range prevotesMap {
			// Here the assumption is that in a single round it is not possible to have 2 value which quorum votes,
			// this would imply at least quorum nodes are malicious which is much higher than our assumption.
			overQuorumVotes := engineCore.OverQuorumVotes(preVotes, quorum)
			if overQuorumVotes != nil {
				proof := &AccountabilityProof{
					Type:     autonity.Misbehaviour,
					Rule:     autonity.C,
					Evidence: overQuorumVotes,
					Message:  precommit,
				}
				proofs = append(proofs, proof)
				fd.logger.Info("Misbehaviour detected", "fault detector", fd.address, "ruleC", autonity.C, "sender", precommit.Sender())
				continue precommitLoop
			}
		}

		// Do we see a quorum of prevotes in the same round, if not we can raise an accusation, since we cannot be sure
		// that these prevotes don't exist, this block also covers the Accusation of C since if over quorum prevotes for
		// V indicates that the corresponding proposal of V do exist, thus we don't need to rise accusation for the missing
		// proposal since over 2/3 member should all ready received it
		prevotes := fd.msgStore.Get(height, func(m *mUtils.Message) bool {
			// since equivocation msgs are stored, we have to query those preVotes which has same value as the proposal.
			return m.Type() == proto.MsgPrevote && m.R() == precommit.R() && m.Value() == precommit.Value()
		})

		if engineCore.OverQuorumVotes(prevotes, quorum) == nil {
			accusation := &AccountabilityProof{
				Type:    autonity.Accusation,
				Rule:    autonity.C1,
				Message: precommit,
			}
			proofs = append(proofs, accusation)
			fd.logger.Info("Accusation detected", "fault detector", fd.address, "ruleC1", autonity.C1, "sender", precommit.Sender())
		}
	}
	return proofs
}

// sendAccountabilityEvents before sending the accountability events, the protocol apply a random delay to make it
// more economic scalable.
func (fd *FaultDetector) sendAccountabilityEvents(proofs []*autonity.AccountabilityEvent, checkExists bool) {
	fd.proofWG.Add(1)
	go func() {
		defer fd.proofWG.Done()

		// only the one on challenge is managed to send innocence proof, so we don't need to check exists if one not send it yet.
		if !checkExists {
			fd.faultDetectorFeed.Send(proofs)
			return
		}

		// since we have slotted reporting with each slot have 10 block period, so there is no need for random delay here.
		unProcessedOnes := fd.filterProcessedOnes(proofs)
		fd.logger.Info("sendAccountabilityEvents with un-processed ones", "number of un-processed ones", len(unProcessedOnes))
		if len(unProcessedOnes) != 0 {
			fd.faultDetectorFeed.Send(unProcessedOnes)
		}
	}()
}

// submitMisbehavior takes proof of misbehavior, and error id to construct the on-chain accountability event, and
// send the event of misbehavior to event channel that is listened by ethereum object to sign the reporting TX.
func (fd *FaultDetector) submitMisbehavior(m *mUtils.Message, evidence []*mUtils.Message, err error, submitCh chan<- *autonity.AccountabilityEvent) {
	rule, e := errorToRule(err)
	if e != nil {
		fd.logger.Warn("error to rule", "fault detector", e)
	}
	proof, err := fd.generateOnChainProof(&AccountabilityProof{
		Type:     autonity.Misbehaviour,
		Rule:     rule,
		Message:  m,
		Evidence: evidence,
	})
	if err != nil {
		fd.logger.Warn("generate misbehavior AccountabilityProof", "fault detector", err)
		return
	}

	// submit misbehavior proof to buffer, it will be sent once aggregated.
	submitCh <- proof
}

func (fd *FaultDetector) accountForAutoIncriminatingProposal(m *mUtils.Message) error {
	if m.R() > constants.MaxRound {
		fd.submitMisbehavior(m.ToLiteProposal(), nil, errInvalidRound, fd.misbehaviourProofsCh)
		return errInvalidRound
	}

	// skip process duplicated msg.
	duplicatedMsg := fd.msgStore.Get(m.H(), func(msg *mUtils.Message) bool {
		return msg.R() == m.R() && msg.Type() == proto.MsgProposal && msg.Sender() == m.Sender() && msg.Value() == m.Value()
	})
	if len(duplicatedMsg) > 0 {
		return errDuplicatedMsg
	}

	// account for wrong proposer.
	if !isProposerMsg(fd.blockchain, m) {
		fd.submitMisbehavior(m.ToLiteProposal(), nil, errProposer, fd.misbehaviourProofsCh)
		return errProposer
	}

	// account for wrong valid round.
	if m.ValidRound() >= m.R() {
		fd.submitMisbehavior(m.ToLiteProposal(), nil, errWrongValidRound, fd.misbehaviourProofsCh)
		return errWrongValidRound
	}

	// account for equivocation
	equivocated := fd.msgStore.Get(m.H(), func(msg *mUtils.Message) bool {
		return msg.R() == m.R() && msg.Type() == proto.MsgProposal && msg.Sender() == m.Sender() && msg.Value() != m.Value()
	})

	if len(equivocated) > 0 {
		var equivocatedMsgs = []*mUtils.Message{equivocated[0].ToLiteProposal()}
		fd.submitMisbehavior(m.ToLiteProposal(), equivocatedMsgs, errEquivocation, fd.misbehaviourProofsCh)
		// we allow the equivocated msg to be stored in msg store.
		fd.msgStore.Save(m)
		return errEquivocation
	}

	// accounting for wrong proposal is no longer depends on the trie db state, it is delayed by delta blocks for if
	// there are quorum prevotes for a bad proposal.
	return nil
}

func (fd *FaultDetector) accountForAutoIncriminatingVote(m *mUtils.Message) error {
	if m.R() > constants.MaxRound {
		fd.submitMisbehavior(m, nil, errInvalidRound, fd.misbehaviourProofsCh)
		return errInvalidRound
	}

	// skip process duplicated for votes.
	duplicatedMsg := fd.msgStore.Get(m.H(), func(msg *mUtils.Message) bool {
		return msg.R() == m.R() && msg.Type() == m.Type() && msg.Sender() == m.Sender() && msg.Value() == m.Value()
	})

	if len(duplicatedMsg) > 0 {
		return errDuplicatedMsg
	}

	// account for equivocation for votes.
	equivocatedMsgs := fd.msgStore.Get(m.H(), func(msg *mUtils.Message) bool {
		return msg.R() == m.R() && msg.Type() == m.Type() && msg.Sender() == m.Sender() && msg.Value() != m.Value()
	})

	if len(equivocatedMsgs) > 0 {
		fd.submitMisbehavior(m, equivocatedMsgs[:1], errEquivocation, fd.misbehaviourProofsCh)
		// we allow store equivocated msg in msg store.
		fd.msgStore.Save(m)
		return errEquivocation
	}
	return nil
}

func errorToRule(err error) (autonity.Rule, error) {
	var rule autonity.Rule
	switch err {
	case errWrongValidRound:
		rule = autonity.WrongValidRound
	case errInvalidRound:
		rule = autonity.InvalidRound
	case errEquivocation:
		rule = autonity.Equivocation
	case errProposer:
		rule = autonity.InvalidProposer
	case errAccountableGarbageMsg:
		rule = autonity.AccountableGarbageMessage
	default:
		return rule, fmt.Errorf("errors of not provable")
	}

	return rule, nil
}

func getProposer(chain BlockChainContext, h uint64, r int64) (common.Address, error) {
	parentHeader := chain.GetHeaderByNumber(h - 1)
	// to prevent the panic on node shutdown.
	if parentHeader == nil {
		return common.Address{}, fmt.Errorf("cannot find parent header")
	}

	if parentHeader.IsGenesis() {
		sort.Sort(parentHeader.Committee)
		return parentHeader.Committee[r%int64(len(parentHeader.Committee))].Address, nil
	}

	proposer := chain.GetAutonityContract().GetProposer(parentHeader, parentHeader.Number.Uint64(), r)
	member := parentHeader.CommitteeMember(proposer)
	if member == nil {
		return common.Address{}, fmt.Errorf("cannot find correct proposer")
	}
	return proposer, nil
}

func isProposerMsg(chain BlockChainContext, m *mUtils.Message) bool {
	proposer, err := getProposer(chain, m.H(), m.R())
	if err != nil {
		log.Error("get proposer err", "err", err)
		return false
	}
	return m.Address == proposer
}

func sameConsensusMsg(a *mUtils.Message, b *mUtils.Message) bool {
	return a.MsgHash() == b.MsgHash()
}
