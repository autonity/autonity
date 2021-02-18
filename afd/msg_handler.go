package afd

import (
	"fmt"
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus"
	"github.com/clearmatics/autonity/consensus/tendermint/crypto"
	"github.com/clearmatics/autonity/core"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/core/vm"
	"sort"
)

// decode consensus msgs, address garbage msg and invalid proposal by returning error.
func checkAutoIncriminatingMsg(chain *core.BlockChain, m *types.ConsensusMessage) error {
	if m.Code == types.MsgProposal {
		return checkProposal(chain, m)
	}

	if m.Code == types.MsgPrevote || m.Code == types.MsgPrecommit {
		return decodeVote(m)
	}

	return errUnknownMsg
}

func checkEquivocation(chain *core.BlockChain, m *types.ConsensusMessage, proof[]types.ConsensusMessage) error {
	// decode msgs
	err := checkAutoIncriminatingMsg(chain, m)
	if err != nil {
		return err
	}

	for i:= 0; i < len(proof); i++ {
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

func sameVote(a *types.ConsensusMessage, b *types.ConsensusMessage) bool {
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
func checkProposal(chain *core.BlockChain, m *types.ConsensusMessage) error {
	var proposal types.Proposal
	err := m.Decode(&proposal)
	if err != nil {
		return errGarbageMsg
	}

	if !isProposerMsg(chain, m) {
		return errProposer
	}

	err = verifyProposal(chain, *proposal.ProposalBlock)
	if err != nil {
		if err == consensus.ErrFutureBlock {
			return errFutureMsg
		} else {
			return errProposal
		}
	}

	return nil
}

//checkMsgSignature, it check if msg is from valid member of the committee.
func checkMsgSignature(chain *core.BlockChain, m *types.ConsensusMessage) error {
	msgHeight, err := m.Height()
	if err != nil {
		return err
	}

	header := chain.CurrentHeader()
	if msgHeight.Cmp(header.Number) > 1 {
		return errFutureMsg
	}

	lastHeader := chain.GetHeaderByNumber(msgHeight.Uint64() - 1)

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
			usedGas        = new(uint64)
			gp             = new(core.GasPool).AddGas(block.GasLimit())
			header         = block.Header()
			parent         = chain.GetBlock(block.ParentHash(), block.NumberU64()-1)
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
			// Might be vulnerable to DoS Attack depending on gaslimit
			// Todo : Double check
			// use default values for vmConfig.
			vmConfig := vm.Config{
				EnablePreimageRecording: true,
				EWASMInterpreter: "",
				EVMInterpreter: "",
			}
			receipt, receiptErr := core.ApplyTransaction(chain.Config(), chain, nil, gp, state, header, tx, usedGas, vmConfig)
			if receiptErr != nil {
				return receiptErr
			}
			receipts = append(receipts, receipt)
		}

		state.Prepare(common.ACHash(block.Number()), block.Hash(), len(block.Transactions()))
		committeeSet, receipt, err := chain.Engine().Finalize(chain, header, state, block.Transactions(), nil, receipts)
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

func isProposerMsg(chain *core.BlockChain, m *types.ConsensusMessage) bool {
	h, _ := m.Height()
	r, _ := m.Round()

	proposer, err := getProposer(chain, h.Uint64(), r)
	if err != nil {
		return false
	}

	return m.Address == proposer
}

func getProposer(chain *core.BlockChain, h uint64, r int64) (common.Address, error) {
	parentHeader := chain.GetHeaderByNumber(h-1)
	if parentHeader.IsGenesis() {
		sort.Sort(parentHeader.Committee)
		return parentHeader.Committee[r%int64(len(parentHeader.Committee))].Address, nil
	}

	statedb, err := chain.StateAt(parentHeader.Hash())
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

func decodeVote(m *types.ConsensusMessage) error {
	var vote types.Vote
	err := m.Decode(&vote)
	if err != nil {
		return errGarbageMsg
	}
	return nil
}