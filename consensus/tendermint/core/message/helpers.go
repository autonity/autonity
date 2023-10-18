package message

import (
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core/types"
	"math/big"
	"testing"
)

func CreatePrevote(t *testing.T, proposalHash common.Hash, round int64, height *big.Int, member types.CommitteeMember) Message {
	expectedMsg := &Prevote{
		Value: proposalHash,
		baseMessage: baseMessage{
			Round:     round,
			Height:    height,
			Signature: nil,
			payload:   nil,
			power:     nil,
			sender:    member.Address,
		},
	}
	return expectedMsg
}
