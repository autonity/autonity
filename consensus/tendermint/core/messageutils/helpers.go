package messageutils

import (
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/core/types"
	"math/big"
	"testing"
)

func CreatePrevote(t *testing.T, proposalHash common.Hash, round int64, height *big.Int, member types.CommitteeMember) *Message {
	var preVote = Vote{
		Round:             round,
		Height:            height,
		ProposedBlockHash: proposalHash,
	}

	encodedVote, err := Encode(&preVote)
	if err != nil {
		t.Fatalf("Expected nil, got %v", err)
		return nil
	}

	expectedMsg := &Message{
		Code:          consensus.MsgPrevote,
		TbftMsgBytes:  encodedVote,
		Address:       member.Address,
		CommittedSeal: []byte{},
		Signature:     []byte{0x1},
		Power:         member.VotingPower,
	}
	return expectedMsg
}
