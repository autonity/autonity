package message

import (
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/rlp"
	"math/big"
	"testing"
)

func CreatePrevote(t *testing.T, proposalHash common.Hash, round uint64, height *big.Int, member types.CommitteeMember) *Message {
	var preVote = Vote{
		Round:             round,
		Height:            height,
		ProposedBlockHash: proposalHash,
	}

	encodedVote, err := rlp.EncodeToBytes(&preVote)
	if err != nil {
		t.Fatalf("Expected nil, got %v", err)
		return nil
	}

	expectedMsg := &Message{
		Code:          consensus.MsgPrevote,
		Payload:       encodedVote,
		Address:       member.Address,
		CommittedSeal: []byte{},
		Signature:     []byte{0x1},
		ConsensusMsg:  ConsensusMsg(&preVote),
		Power:         member.VotingPower,
	}
	return expectedMsg
}
