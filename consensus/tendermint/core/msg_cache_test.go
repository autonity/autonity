package core

import (
	"math/big"
	"testing"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus/tendermint/algorithm"
	"github.com/clearmatics/autonity/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProposal(t *testing.T) {
	m := newMessageStore()

	b := types.NewBlockWithHeader(&types.Header{
		Number: big.NewInt(1),
	})

	h := b.Hash()
	m.addValue(h, b)

	p := &algorithm.ConsensusMessage{
		MsgType:    algorithm.Propose,
		Height:     b.NumberU64(),
		Round:      0,
		Value:      algorithm.ValueID(h),
		ValidRound: -1,
	}

	internalMessage := NewProposal(p.Round, new(big.Int).SetUint64(p.Height), p.ValidRound, m.value(common.Hash(p.Value)))

	a := common.Address{}
	marshalledInternalMessage, err := Encode(internalMessage)
	require.NoError(t, err)
	msg := &Message{
		Code:          msgProposal,
		Address:       a,
		CommittedSeal: []byte{}, // Not sure why this is empty but it seems to be set like this everywhere.
		Msg:           marshalledInternalMessage,
	}
	m.addMessage(msg, p)

	pv := &algorithm.ConsensusMessage{
		MsgType: algorithm.Prevote,
		Height:  1,
		Round:   0,
		Value:   algorithm.ValueID(h),
	}
	retrieved := m.matchingProposal(pv)
	assert.Equal(t, p, retrieved)
}
