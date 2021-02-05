package afd

import (
	"fmt"
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus"
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
	"math/rand"
	"sort"
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

	// buffer for proposer of rounds rather to get it by lifting evm again and again.
	// map[height]map[round]common.address
	proposersMap map[uint64]map[int64]common.Address
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
		case ev := <-fd.blockChan:
			// take my challenge from latest state DB, and provide innocent proof if there are any.
			err := fd.handleMyChallenges(ev.Block, ev.Hash)
			if err != nil {
				fd.logger.Warn("handle challenge","afd", err)
			}

			// todo: tell rule engine to run patterns over msg store on each new height.
			fd.ruleEngine.run()
		case ev, ok := <-fd.tendermintMsgSub.Chan():
			// take consensus msg from p2p protocol manager.
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

func (fd *FaultDetector) processAutoIncriminatingMsg(m *types.ConsensusMessage) error {
	// msg is checked, then do auto-incriminating checking
	switch m.Code {
	case types.MsgProposal:
		return fd.processProposal(m)
	case types.MsgPrevote:
		return fd.processPrevote(m)
	case types.MsgPrecommit:
		return fd.processPrecommit(m)
	default:
		fd.logger.Error("Invalid message", "afd", m)
	}

	return nil
}

func (fd *FaultDetector) generateOnChainProof(m *types.ConsensusMessage, proofs []types.ConsensusMessage, err error) (types.OnChainProof, error) {
	var challenge types.OnChainProof
	switch err {
	case errEquivocation:
		challenge.Rule = uint8(types.Equivocation)
	case errProposer:
		challenge.Rule = uint8(types.InvalidProposer)
	case errProposal:
		challenge.Rule = uint8(types.InvalidProposal)
	case errGarbageMsg:
		challenge.Rule = uint8(types.GarbageMessage)
	default:
		return challenge, fmt.Errorf("errors of not provable")
	}
	h, _ := m.Height()
	r, _ := m.Round()
	challenge.Height = h
	challenge.Round = uint64(r)
	challenge.MsgType = m.Code
	challenge.Sender = m.Address

	// generate raw bytes encoded in rlp, it is by passed into precompiled contracts.
	var rawProof types.RawProof
	rawProof.Rule = challenge.Rule
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

// processMisbehavior takes proofs of misbehavior msg, and error id to form the on-chain proof, and
// send the proof of misbehavior to TX pool.
func (fd *FaultDetector) processMisbehavior(m *types.ConsensusMessage, proofs []types.ConsensusMessage, err error) {

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
	err := fd.preProcessMsg(m)
	if err != nil {
		if err == errFutureMsg {
			fd.bufferMsg(m)
		}
		return err
	}

	// test auto incriminating msg.
	err = fd.processAutoIncriminatingMsg(m)
	if err != nil {
		proofs := []types.ConsensusMessage{*m}
		fd.processMisbehavior(m, proofs, err)
		return err
	}

	// store msg, if there is equivocation then rise errEquivocation and return proofs.
	equivocationProof, err := fd.msgStore.StoreMsg(m)
	if err == errEquivocation {
		fd.processMisbehavior(m, equivocationProof, err)
		return err
	}
	return nil
}

func (fd *FaultDetector) getProposer(h uint64, r int64) (common.Address, error) {
	// todo: before lifting evm again and again, let's buffer proposers in afd.
	parentHeader := fd.blockchain.GetHeaderByNumber(h-1)
	if parentHeader.IsGenesis() {
		sort.Sort(parentHeader.Committee)
		return parentHeader.Committee[r%int64(len(parentHeader.Committee))].Address, nil
	}

	statedb, err := fd.blockchain.StateAt(parentHeader.Hash())
	if err != nil {
		return common.Address{}, err
	}

	proposer := fd.blockchain.GetAutonityContract().GetProposerFromAC(parentHeader, statedb, parentHeader.Number.Uint64(), r)
	member := parentHeader.CommitteeMember(proposer)
	if member == nil {
		return common.Address{}, fmt.Errorf("cannot find correct proposer")
	}
	return proposer, nil
}

func (fd *FaultDetector) isProposerMsg(m *types.ConsensusMessage) bool {
	h, _ := m.Height()
	r, _ := m.Round()

	proposer, err := fd.getProposer(h.Uint64(), r)
	if err != nil {
		return false
	}

	return m.Address == proposer
}

func (fd *FaultDetector) verifyProposal(proposal types.Block) error {
	block := &proposal
	if fd.blockchain.HasBadBlock(block.Hash()) {
		return core.ErrBlacklistedHash
	}

	err := fd.blockchain.Engine().VerifyHeader(fd.blockchain, block.Header(), false)
	if err == nil || err == types.ErrEmptyCommittedSeals {
		var (
			receipts types.Receipts
			usedGas        = new(uint64)
			gp             = new(core.GasPool).AddGas(block.GasLimit())
			header         = block.Header()
			proposalNumber = header.Number.Uint64()
			parent         = fd.blockchain.GetBlock(block.ParentHash(), block.NumberU64()-1)
		)

		// We need to process all of the transaction to get the latest state to get the latest committee
		state, stateErr := fd.blockchain.StateAt(parent.Root())
		if stateErr != nil {
			return stateErr
		}

		// Validate the body of the proposal
		if err = fd.blockchain.Validator().ValidateBody(block); err != nil {
			return err
		}

		// sb.blockchain.Processor().Process() was not called because it calls back Finalize() and would have modified the proposal
		// Instead only the transactions are applied to the copied state
		for i, tx := range block.Transactions() {
			state.Prepare(tx.Hash(), block.Hash(), i)
			// Might be vulnerable to DoS Attack depending on gaslimit
			// Todo : Double check
			// use default values for vmConfig.
			vmConfig := vm.Config{
				EnablePreimageRecording: true,
				EWASMInterpreter: "",
				EVMInterpreter: "",
			}
			receipt, receiptErr := core.ApplyTransaction(fd.blockchain.Config(), fd.blockchain, nil, gp, state, header, tx, usedGas, vmConfig)
			if receiptErr != nil {
				return receiptErr
			}
			receipts = append(receipts, receipt)
		}

		state.Prepare(common.ACHash(block.Number()), block.Hash(), len(block.Transactions()))
		committeeSet, receipt, err := fd.blockchain.Engine().Finalize(fd.blockchain, header, state, block.Transactions(), nil, receipts)
		receipts = append(receipts, receipt)
		//Validate the state of the proposal
		if err = fd.blockchain.Validator().ValidateState(block, state, receipts, *usedGas); err != nil {
			return err
		}

		//Perform the actual comparison
		if len(header.Committee) != len(committeeSet) {
			fd.logger.Error("wrong committee set",
				"proposalNumber", proposalNumber,
				"extraLen", len(header.Committee),
				"currentLen", len(committeeSet),
				"committee", header.Committee,
				"current", committeeSet,
			)
			return consensus.ErrInconsistentCommitteeSet
		}

		for i := range committeeSet {
			if header.Committee[i].Address != committeeSet[i].Address ||
				header.Committee[i].VotingPower.Cmp(committeeSet[i].VotingPower) != 0 {
				fd.logger.Error("wrong committee member in the set",
					"index", i,
					"currentVerifier", fd.address.String(),
					"proposalNumber", proposalNumber,
					"headerCommittee", header.Committee[i],
					"computedCommittee", committeeSet[i],
					"fullHeader", header.Committee,
					"fullComputed", committeeSet,
				)
				return consensus.ErrInconsistentCommitteeSet
			}
		}

		return nil
	}
	return err
}

// buffer Msg since node are not synced to verify it.
func (fd *FaultDetector) bufferMsg(m *types.ConsensusMessage) {
	// todo: buffer the msg.
}

// processProposal, checks if proposal is valid (no garbage msg, no invalid tx ),
// it's from correct proposer.
func (fd *FaultDetector) processProposal(m *types.ConsensusMessage) error {
	var proposal types.Proposal
	err := m.Decode(&proposal)
	if err != nil {
		return errGarbageMsg
	}

	if !fd.isProposerMsg(m) {
		return errProposer
	}

	err = fd.verifyProposal(*proposal.ProposalBlock)
	if err != nil {
		if err == consensus.ErrFutureBlock {
			fd.bufferMsg(m)
		} else {
			return errProposal
		}
	}

	return nil
}

func (fd *FaultDetector) processPrevote(m *types.ConsensusMessage) error {
	var preVote types.Vote
	err := m.Decode(&preVote)
	if err != nil {
		return errGarbageMsg
	}
	return nil
}

func (fd *FaultDetector) processPrecommit(m *types.ConsensusMessage) error {
	var preCommit types.Vote
	err := m.Decode(&preCommit)
	if err != nil {
		return errGarbageMsg
	}
	return nil
}

//pre-process msg, it check if msg is from valid member of the committee, it return
func (fd *FaultDetector) preProcessMsg(m *types.ConsensusMessage) error {
	msgHeight, err := m.Height()
	if err != nil {
		return err
	}

	header := fd.blockchain.CurrentHeader()
	if msgHeight.Cmp(header.Number) > 1 {
		return errFutureMsg
	}

	lastHeader := fd.blockchain.GetHeaderByNumber(msgHeight.Uint64() - 1)

	if _, err = m.Validate(crypto.CheckValidatorSignature, lastHeader); err != nil {
		fd.logger.Error("Msg is not from committee member", "err", err)
		return errNotCommitteeMsg
	}
	return nil
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
		if challenges[i].Sender == fd.address {
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
			unPresented := fd.filterUnPresentedChallenges(&proofs)
			if len(unPresented) != 0 {
				fd.afdFeed.Send(types.SubmitProofEvent{Proofs:unPresented, Type:t})
			}
		}
	}()
}

func (fd *FaultDetector) filterUnPresentedChallenges(proofs *[]types.OnChainProof) []types.OnChainProof {
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
			if (*proofs)[i].Sender == challenges[i].Sender && (*proofs)[i].Height == challenges[i].Height &&
				(*proofs)[i].Round == challenges[i].Round && (*proofs)[i].Rule == challenges[i].Rule &&
				(*proofs)[i].MsgType == challenges[i].MsgType {
				present = true
			}
		}
		if !present {
			result = append(result, (*proofs)[i])
		}
	}

	return result
}
