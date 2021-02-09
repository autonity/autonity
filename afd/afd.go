package afd

import (
	"fmt"
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus/tendermint/events"
	"github.com/clearmatics/autonity/core"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/event"
	"github.com/clearmatics/autonity/log"
	"github.com/clearmatics/autonity/p2p"
	"github.com/clearmatics/autonity/rlp"
	"github.com/syndtr/goleveldb/leveldb/errors"
	"math/rand"
	"sync"
	"time"
)

var (
	errFutureMsg = errors.New("future height msg")
	errGarbageMsg = errors.New("garbage msg")
	errNotCommitteeMsg = errors.New("msg from none committee member")
	errProposer = errors.New("proposal is not from proposer")
	errProposal = errors.New("proposal have invalid values")
	errEquivocation = errors.New("equivocation happens")
	errUnknownMsg = errors.New("unknown consensus msg")
	errInvalidChallenge = errors.New("invalid challenge")
	errNoEvidence = errors.New("no evidence")
)

// Fault detector, it subscribe chain event to trigger rule engine to apply patterns over
// msg store, it send proof of challenge if it detects any potential misbehavior, either it
// read state db on each new height to get latest challenges from autonity contract's view,
// and to prove its innocent if there were any challenges on the suspicious node.
type FaultDetector struct {
	// use below 3 members to send proof via transaction issuing.
	wg sync.WaitGroup
	afdFeed event.Feed
	scope event.SubscriptionScope

	// use below 2 members to forward consensus msg from protocol manager to afd.
	tendermintMsgSub *event.TypeMuxSubscription
	tendermintMsgMux *event.TypeMuxSilent

	// below 2 members subscribe block event to trigger execution
	// of rule engine and make proof of innocent.
	blockChan chan core.ChainEvent
	blockSub event.Subscription

	// chain context to validate consensus msgs.
	blockchain *core.BlockChain

	// node address
	address common.Address

	// msg store
	msgStore *MsgStore

	// rule engine
	ruleEngine *RuleEngine

	// future height msg buffer
	futureMsgs map[uint64][]*types.ConsensusMessage

	logger log.Logger
}

var(
	// todo: to be configured at genesis.
	randomDelayWindow = 10000
)

// call by ethereum object to create fd instance.
func NewFaultDetector(chain *core.BlockChain, nodeAddress common.Address) *FaultDetector {
	logger := log.New("afd", nodeAddress)
	fd := &FaultDetector{
		address: nodeAddress,
		blockChan:  make(chan core.ChainEvent, 300),
		blockchain: chain,
		msgStore: new(MsgStore),
		ruleEngine: new(RuleEngine),
		logger:logger,
		tendermintMsgMux:  event.NewTypeMuxSilent(logger),
		proposersMap: map[uint64]map[int64]common.Address{},
		futureMsgs: make(map[uint64][]*types.ConsensusMessage),
	}

	// init accountability precompiled contracts.
	initAccountabilityContracts(chain)

	// subscribe tendermint msg
	s := fd.tendermintMsgMux.Subscribe(events.MessageEvent{})
	fd.tendermintMsgSub = s
	return fd
}

// listen for new block events from block-chain, do the tasks like take challenge and provide proof for innocent, the
// AFD rule engine could also triggered from here to scan those msgs of msg store by applying rules.
func (fd *FaultDetector) FaultDetectorEventLoop() {
	fd.blockSub = fd.blockchain.SubscribeChainEvent(fd.blockChan)

	for {
		select {
		// chain event update, provide proof of innocent if one is on challenge, rule engine scanning is triggered also.
		case ev := <-fd.blockChan:
			// take my challenge from latest state DB, and provide innocent proof if there are any.
			err := fd.handleMyChallenges(ev.Block, ev.Hash)
			if err != nil {
				fd.logger.Warn("handle challenge","afd", err)
			}

			// before run rule engine over msg store, check to process any buffered msg.
			fd.processBufferedMsgs(ev.Block.NumberU64())

			// run rule engine over msg store on each height update.
			fd.runRuleEngine(ev.Block.NumberU64())
		// to handle consensus msg from p2p layer.	
		case ev, ok := <-fd.tendermintMsgSub.Chan():
			if !ok {
				return
			}
			switch e := ev.Data.(type) {
			case events.MessageEvent:
				msg := new(types.ConsensusMessage)
				if err := msg.FromPayload(e.Payload); err != nil {
					fd.logger.Error("invalid payload", "afd", err)
					continue
				}
				if err := fd.processMsg(msg); err != nil {
					fd.logger.Error("process consensus msg", "afd", err)
					continue
				}
			}

		case <-fd.blockSub.Err():
			return
		}
	}
}

// HandleMsg is called by p2p protocol manager to deliver the consensus msg to afd.
func (fd *FaultDetector) HandleMsg(addr common.Address, msg p2p.Msg) {
	if msg.Code != types.TendermintMsg {
		return
	}

	var data []byte
	if err := msg.Decode(&data); err != nil {
		log.Error("cannot decode consensus msg", "from", addr)
		return
	}

	// post consensus event to event loop.
	fd.tendermintMsgMux.Post(events.MessageEvent{Payload:data})
	return
}

func (fd *FaultDetector) Stop() {
	fd.scope.Close()
	fd.blockSub.Unsubscribe()
	fd.tendermintMsgSub.Unsubscribe()
	fd.wg.Wait()
	cleanContracts()
}

// call by ethereum object to subscribe proofs Events.
func (fd *FaultDetector) SubscribeAFDEvents(ch chan<- types.SubmitProofEvent) event.Subscription {
	return fd.scope.Track(fd.afdFeed.Subscribe(ch))
}

// run rule engine over latest msg store, if the return proofs is not empty, then rise challenge.
func (fd *FaultDetector) runRuleEngine(headHeight uint64) {
	proofs := fd.ruleEngine.run(fd.msgStore, headHeight)
	if len(proofs) > 0 {
		fd.sendProofs(types.ChallengeProof, proofs)
	}
}

func (fd *FaultDetector) generateOnChainProof(m *types.ConsensusMessage, proofs []types.ConsensusMessage, err error) (types.OnChainProof, error) {
	var challenge types.OnChainProof
	challenge.SenderHash = types.RLPHash(m.Address)
	challenge.MsgHash = types.RLPHash(m.Payload())

	var rawProof types.RawProof
	switch err {
	case errEquivocation:
		rawProof.Rule = uint8(types.Equivocation)
	case errProposer:
		rawProof.Rule = uint8(types.InvalidProposer)
	case errProposal:
		rawProof.Rule = uint8(types.InvalidProposal)
	case errGarbageMsg:
		rawProof.Rule = uint8(types.GarbageMessage)
	default:
		return challenge, fmt.Errorf("errors of not provable")
	}
	// generate raw bytes encoded in rlp, it is by passed into precompiled contracts.
	rawProof.Message = m.Payload()
	for i:= 0; i < len(proofs); i++ {
		rawProof.Evidence = append(rawProof.Evidence, proofs[i].Payload())
	}

	rp, err := rlp.EncodeToBytes(rawProof)
	if err != nil {
		fd.logger.Warn("fail to rlp encode raw proof", "afd", err)
		return challenge, err
	}

	challenge.RawProofBytes = rp
	return challenge, nil
}

// submitMisbehavior takes proofs of misbehavior msg, and error id to form the on-chain proof, and
// send the proof of misbehavior to TX pool.
func (fd *FaultDetector) submitMisbehavior(m *types.ConsensusMessage, proofs []types.ConsensusMessage, err error) {

	proof, err := fd.generateOnChainProof(m, proofs, err)
	if err != nil {
		fd.logger.Warn("generate misbehavior proof", "afd", err)
	}
	ps := []types.OnChainProof{proof}

	fd.sendProofs(types.ChallengeProof, ps)
}

// processMsg it decode consensus msg, apply auto-incriminating, equivocation rules to it,
// and store it to msg store.
func (fd *FaultDetector) processMsg(m *types.ConsensusMessage) error {
	// pre-check if msg is from valid committee member
	err := checkMsgSignature(fd.blockchain, m)
	if err != nil {
		if err == errFutureMsg {
			fd.bufferMsg(m)
		}
		return err
	}

	// decode consensus msg, and auto-incriminating msg is addressed here.
	err = preProcessConsensusMsg(fd.blockchain, m)
	if err != nil {
		if err == errFutureMsg {
			fd.bufferMsg(m)
		} else {
			proofs := []types.ConsensusMessage{*m}
			fd.submitMisbehavior(m, proofs, err)
			return err
		}
	}

	// store msg, if there is equivocation then rise errEquivocation and return proofs.
	equivocationProof, err := fd.msgStore.StoreMsg(m)
	if err == errEquivocation {
		fd.submitMisbehavior(m, equivocationProof, err)
		return err
	}
	return nil
}

// processBufferedMsgs, called on chain event update, it take msgs from the buffered future height msgs, process them.
func (fd *FaultDetector) processBufferedMsgs(headHeight uint64) {
	for height, msgs := range fd.futureMsgs {
		if height <= headHeight {
			for i:= 0; i < len(msgs); i++ {
				if err := fd.processMsg(msgs[i]); err != nil {
					fd.logger.Error("process consensus msg", "afd", err)
					continue
				}
			}
		}
	}
}

// buffer Msg since node are not synced to verify it.
func (fd *FaultDetector) bufferMsg(m *types.ConsensusMessage) {
	h, err := m.Height()
	if err != nil {
		return
	}

	fd.futureMsgs[h.Uint64()] = append(fd.futureMsgs[h.Uint64()], m)
}

// get challenges from blockchain via autonityContract calls.
func (fd *FaultDetector) handleMyChallenges(block *types.Block, hash common.Hash) error {
	var innocentProofs []types.OnChainProof
	state, err := fd.blockchain.StateAt(hash)
	if err != nil {
		return err
	}

	challenges := fd.blockchain.GetAutonityContract().GetChallenges(block.Header(), state)
	for i:=0; i < len(challenges); i++ {
		if challenges[i].SenderHash == types.RLPHash(fd.address) {
			p, err := fd.proveInnocent(challenges[i])
			if err != nil {
				continue
			}
			innocentProofs = append(innocentProofs, p)
		}
	}

	// send proofs via standard transaction.
	fd.sendProofs(types.InnocentProof, innocentProofs)
	return nil
}

// get proof of innocent over msg store.
func (fd *FaultDetector) proveInnocent(challenge types.OnChainProof) (types.OnChainProof, error) {
	// todo: get proof from msg store over the rule.
	var proof types.OnChainProof
	return proof, nil
}

// send proofs via event which will handled by ethereum object to signed the TX to send proof.
func (fd *FaultDetector) sendProofs(t types.ProofType,  proofs[]types.OnChainProof) {
	fd.wg.Add(1)
	go func() {
		defer fd.wg.Done()
		if t == types.InnocentProof {
			fd.afdFeed.Send(types.SubmitProofEvent{Proofs:proofs, Type:t})
		}

		if t == types.ChallengeProof {
			// wait for random milliseconds (under the range of 10 seconds) to check if need to rise challenge.
			rand.Seed(time.Now().UnixNano())
			n := rand.Intn(randomDelayWindow)
			time.Sleep(time.Duration(n)*time.Millisecond)
			unPresented := fd.filterMissingChallenges(&proofs)
			if len(unPresented) != 0 {
				fd.afdFeed.Send(types.SubmitProofEvent{Proofs:unPresented, Type:t})
			}
		}
	}()
}

func (fd *FaultDetector) filterMissingChallenges(proofs *[]types.OnChainProof) []types.OnChainProof {
	// get latest chain state.
	var result []types.OnChainProof
	state, err := fd.blockchain.State()
	if err != nil {
		return nil
	}

	header := fd.blockchain.CurrentBlock().Header()
	challenges := fd.blockchain.GetAutonityContract().GetChallenges(header, state)

	for i:=0; i < len(*proofs); i++ {
		present := false
		for i:=0; i < len(challenges); i++ {
			if (*proofs)[i].MsgHash == challenges[i].MsgHash {
				present = true
			}
		}
		if !present {
			result = append(result, (*proofs)[i])
		}
	}

	return result
}
