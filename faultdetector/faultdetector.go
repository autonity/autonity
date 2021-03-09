package faultdetector

import (
	"fmt"
	"github.com/clearmatics/autonity/autonity"
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus"
	tendermintBackend "github.com/clearmatics/autonity/consensus/tendermint/backend"
	"github.com/clearmatics/autonity/consensus/tendermint/bft"
	tendermintCore "github.com/clearmatics/autonity/consensus/tendermint/core"
	"github.com/clearmatics/autonity/consensus/tendermint/crypto"
	"github.com/clearmatics/autonity/consensus/tendermint/events"
	"github.com/clearmatics/autonity/core"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/core/vm"
	"github.com/clearmatics/autonity/event"
	"github.com/clearmatics/autonity/log"
	"github.com/clearmatics/autonity/p2p"
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
	randomDelayWindow            = 1000 * 5                          // (0, 5] seconds random time window
	deltaToWaitForAccountability = 30                                // Wait until the GST + delta (30 blocks) to start rule scan.
	msgBufferInHeight            = deltaToWaitForAccountability + 60 // buffer such range of msgs in height at msg store.
	errFutureMsg                 = errors.New("future height msg")
	errGarbageMsg                = errors.New("garbage msg")
	errNotCommitteeMsg           = errors.New("msg from none committee member")
	errProposer                  = errors.New("proposal is not from proposer")
	errProposal                  = errors.New("proposal have invalid values")
	errEquivocation              = errors.New("equivocation happens")
	errUnknownMsg                = errors.New("unknown consensus msg")
)

// Fault detector, it subscribe chain event to trigger rule engine to apply patterns over
// msg store, it send proof of challenge if it detects any potential misbehavior, either it
// read state db on each new height to get latest challenges from autonity contract's view,
// and to prove its innocent if there were any challenges on the suspicious node.
type FaultDetector struct {
	// use below 3 members to send proof via transaction issuing.
	wg      sync.WaitGroup
	afdFeed event.Feed
	scope   event.SubscriptionScope

	// use below 2 members to forward consensus msg from protocol manager to faultdetector.
	tendermintMsgSub *event.TypeMuxSubscription
	tendermintMsgMux *event.TypeMuxSilent

	// below 2 members subscribe block event to trigger execution
	// of rule engine and make proof of innocent.
	blockChan chan core.ChainEvent
	blockSub  event.Subscription

	// chain context to validate consensus msgs.
	blockchain *core.BlockChain

	// node address
	address common.Address

	// msg store
	msgStore *MsgStore

	// future height msg buffer
	futureMsgs map[uint64][]*tendermintCore.Message

	// buffer quorum for blocks.
	totalPowers map[uint64]uint64

	// buffer those proofs, aggregate them into single TX to send with latest nonce of account.
	bufferedProofs []autonity.OnChainProof

	logger log.Logger
}

// call by ethereum object to create fd instance.
func NewFaultDetector(chain *core.BlockChain, nodeAddress common.Address) *FaultDetector {
	logger := log.New("faultdetector", nodeAddress)
	fd := &FaultDetector{
		address:          nodeAddress,
		blockChan:        make(chan core.ChainEvent, 300),
		blockchain:       chain,
		msgStore:         newMsgStore(),
		logger:           logger,
		tendermintMsgMux: event.NewTypeMuxSilent(logger),
		futureMsgs:       make(map[uint64][]*tendermintCore.Message),
		totalPowers:      make(map[uint64]uint64),
	}

	// register faultdetector contracts on evm's precompiled contract set.
	registerAFDContracts(chain)

	// subscribe tendermint msg
	s := fd.tendermintMsgMux.Subscribe(events.MessageEvent{})
	fd.tendermintMsgSub = s
	return fd
}

func (fd *FaultDetector) quorum(h uint64) uint64 {
	power, ok := fd.totalPowers[h]
	if ok {
		return bft.Quorum(power)
	}

	return bft.Quorum(fd.blockchain.GetHeaderByNumber(h).TotalVotingPower())
}

func (fd *FaultDetector) savePower(h uint64, power uint64) {
	fd.totalPowers[h] = power
}

func (fd *FaultDetector) deletePower(h uint64) {
	delete(fd.totalPowers, h)
}

// listen for new block events from block-chain, do the tasks like take challenge and provide proof for innocent, the
// AFD rule engine could also triggered from here to scan those msgs of msg store by applying rules.
func (fd *FaultDetector) FaultDetectorEventLoop() {
	fd.blockSub = fd.blockchain.SubscribeChainEvent(fd.blockChan)

	for {
		select {
		// chain event update, provide proof of innocent if one is on challenge, rule engine scanning is triggered also.
		case ev := <-fd.blockChan:
			fd.savePower(ev.Block.Number().Uint64(), ev.Block.Header().TotalVotingPower())

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
			fd.deletePower(ev.Block.NumberU64() - uint64(msgBufferInHeight))
		// to handle consensus msg from p2p layer.
		case ev, ok := <-fd.tendermintMsgSub.Chan():
			if !ok {
				return
			}
			switch e := ev.Data.(type) {
			case events.MessageEvent:
				msg := new(tendermintCore.Message)
				if err := msg.FromPayload(e.Payload); err != nil {
					fd.logger.Error("invalid payload", "faultdetector", err)
					continue
				}
				if err := fd.processMsg(msg); err != nil {
					fd.logger.Warn("process consensus msg", "faultdetector", err)
					continue
				}
			}

		case <-fd.blockSub.Err():
			return
		}
	}
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

// HandleMsg is called by p2p protocol manager to deliver the consensus msg to faultdetector.
func (fd *FaultDetector) HandleMsg(addr common.Address, msg p2p.Msg) {
	if msg.Code != tendermintBackend.TendermintMsg {
		return
	}

	var data []byte
	if err := msg.Decode(&data); err != nil {
		log.Error("cannot decode consensus msg", "from", addr)
		return
	}

	// post consensus event to event loop.
	fd.tendermintMsgMux.Post(events.MessageEvent{Payload: data})
}

// since tendermint gossip only send to remote peer, so to handle self msgs called by protocol manager.
func (fd *FaultDetector) HandleSelfMsg(payload []byte) {
	fd.tendermintMsgMux.Post(events.MessageEvent{Payload: payload})
}

func (fd *FaultDetector) Stop() {
	fd.scope.Close()
	fd.blockSub.Unsubscribe()
	fd.tendermintMsgSub.Unsubscribe()
	fd.wg.Wait()
	unRegisterAFDContracts()
}

// call by ethereum object to subscribe proofs Events.
func (fd *FaultDetector) SubscribeAFDEvents(ch chan<- AccountabilityEvent) event.Subscription {
	return fd.scope.Track(fd.afdFeed.Subscribe(ch))
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

func (fd *FaultDetector) randomDelay() {
	// wait for random milliseconds (under the range of 10 seconds) to check if need to rise challenge.
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(randomDelayWindow)
	time.Sleep(time.Duration(n) * time.Millisecond)
}

// send proofs via event which will handled by ethereum object to signed the TX to send proof.
func (fd *FaultDetector) sendProofs(proofs []autonity.OnChainProof) {
	fd.wg.Add(1)
	go func() {
		defer fd.wg.Done()
		fd.randomDelay()
		unPresented := fd.filterPresentedOnes(&proofs)
		if len(unPresented) != 0 {
			fd.afdFeed.Send(AccountabilityEvent{Proofs: unPresented})
		}
	}()
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

// --------------------Functions from msg_handler.go-------------------

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
			proofs := []tendermintCore.Message{*m}
			fd.submitMisbehavior(m, proofs, err)
			return err
		}
	}

	// store msg, if there is equivocation, msg store would then rise errEquivocation and proofs.
	p, err := fd.msgStore.Save(m)
	if err == errEquivocation && p != nil {
		proof := []tendermintCore.Message{*p}
		fd.submitMisbehavior(m, proof, err)
		return err
	}
	return nil
}

// processBufferedMsgs, called on chain event update, it process msgs from the latest height buffered before.
func (fd *FaultDetector) processBufferedMsgs(height uint64) {
	for h, msgs := range fd.futureMsgs {
		if h <= height {
			for i := 0; i < len(msgs); i++ {
				if err := fd.processMsg(msgs[i]); err != nil {
					fd.logger.Error("process consensus msg", "faultdetector", err)
					continue
				}
			}
		}
	}
}

// buffer Msg since local chain may not synced yet to verify if msg is from correct committee.
func (fd *FaultDetector) bufferMsg(m *tendermintCore.Message) {
	h, err := m.Height()
	if err != nil {
		return
	}

	fd.futureMsgs[h.Uint64()] = append(fd.futureMsgs[h.Uint64()], m)
}

/////// common helper functions shared between faultdetector and precompiled contract to validate msgs.

// decode consensus msgs, address garbage msg and invalid proposal by returning error.
func checkAutoIncriminatingMsg(chain *core.BlockChain, m *tendermintCore.Message) error {
	if m.Code == msgProposal {
		return checkProposal(chain, m)
	}

	if m.Code == msgPrevote || m.Code == msgPrecommit {
		return decodeVote(m)
	}

	return errUnknownMsg
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

// checkProposal, checks if proposal is valid and it's from correct proposer.
func checkProposal(chain *core.BlockChain, m *tendermintCore.Message) error {
	var proposal tendermintCore.Proposal
	err := m.Decode(&proposal)
	if err != nil {
		return errGarbageMsg
	}
	if !isProposerMsg(chain, m) {
		return errProposer
	}

	err = verifyProposal(chain, *proposal.ProposalBlock)
	// due to network delay or timing issue, when AFD validate a proposal, that proposal could already be committed on the chain view.
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

//checkMsgSignature, it check if msg is from valid member of the committee.
func checkMsgSignature(chain *core.BlockChain, m *tendermintCore.Message) error {
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

func isProposerMsg(chain *core.BlockChain, m *tendermintCore.Message) bool {
	h, _ := m.Height()
	r, _ := m.Round()

	proposer, err := getProposer(chain, h.Uint64(), r)
	if err != nil {
		return false
	}

	return m.Address == proposer
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

func decodeVote(m *tendermintCore.Message) error {
	var vote tendermintCore.Vote
	err := m.Decode(&vote)
	if err != nil {
		return errGarbageMsg
	}
	return nil
}
