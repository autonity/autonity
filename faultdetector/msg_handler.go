package faultdetector

import (
	"fmt"
	"github.com/clearmatics/autonity/autonity"
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus"
	core2 "github.com/clearmatics/autonity/consensus/tendermint/core"
	"github.com/clearmatics/autonity/consensus/tendermint/crypto"
	"github.com/clearmatics/autonity/core"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/core/vm"
	"github.com/clearmatics/autonity/rlp"
	"sort"
)

// convert the raw proofs into on-chain proof which contains raw bytes of messages.
func (fd *FaultDetector) generateOnChainProof(m *core2.Message, proofs []core2.Message, rule Rule) (autonity.OnChainProof, error) {
	var challenge autonity.OnChainProof
	challenge.Sender = m.Address
	challenge.Msghash = types.RLPHash(m.Payload())

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
		return challenge, err
	}

	challenge.Rawproof = rp
	return challenge, nil
}

// submitMisbehavior takes proofs of misbehavior msg, and error id to form the on-chain proof, and
// send the proof of misbehavior to event channel.
func (fd *FaultDetector) submitMisbehavior(m *core2.Message, proofs []core2.Message, err error) {
	rule, e := errorToRule(err)
	if e != nil {
		fd.logger.Warn("error to rule", "faultdetector", e)
	}
	proof, err := fd.generateOnChainProof(m, proofs, rule)
	if err != nil {
		fd.logger.Warn("generate misbehavior proof", "faultdetector", err)
		return
	}
	ps := []autonity.OnChainProof{proof}

	fd.sendProofs(ChallengeProof, ps)
}

// processMsg, check and submit any auto-incriminating, equivocation challenges, and then only store checked msg into msg store.
func (fd *FaultDetector) processMsg(m *core2.Message) error {
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
			proofs := []core2.Message{*m}
			fd.submitMisbehavior(m, proofs, err)
			return err
		}
	}

	// store msg, if there is equivocation, msg store would then rise errEquivocation and proofs.
	p, err := fd.msgStore.Save(m)
	if err == errEquivocation && p != nil {
		proof := []core2.Message{*p}
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
func (fd *FaultDetector) bufferMsg(m *core2.Message) {
	h, err := m.Height()
	if err != nil {
		return
	}

	fd.futureMsgs[h.Uint64()] = append(fd.futureMsgs[h.Uint64()], m)
}

/////// common helper functions shared between faultdetector and precompiled contract to validate msgs.

// decode consensus msgs, address garbage msg and invalid proposal by returning error.
func checkAutoIncriminatingMsg(chain *core.BlockChain, m *core2.Message) error {
	if m.Code == msgProposal {
		return checkProposal(chain, m)
	}

	if m.Code == msgPrevote || m.Code == msgPrecommit {
		return decodeVote(m)
	}

	return errUnknownMsg
}

func checkEquivocation(chain *core.BlockChain, m *core2.Message, proof []core2.Message) error {
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

func sameVote(a *core2.Message, b *core2.Message) bool {
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
func checkProposal(chain *core.BlockChain, m *core2.Message) error {
	var proposal core2.Proposal
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
func checkMsgSignature(chain *core.BlockChain, m *core2.Message) error {
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

func isProposerMsg(chain *core.BlockChain, m *core2.Message) bool {
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

func decodeVote(m *core2.Message) error {
	var vote core2.Vote
	err := m.Decode(&vote)
	if err != nil {
		return errGarbageMsg
	}
	return nil
}
