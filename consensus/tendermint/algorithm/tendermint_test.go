package algorithm

import (
	"crypto/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newNodeId(t *testing.T) NodeID {
	var nodeID NodeID
	_, err := rand.Read(nodeID[:])
	require.NoError(t, err)
	return nodeID
}

func TestStartRound(t *testing.T) {
	proposer := newNodeId(t)
	// node := newNodeId(t)
	var round int64 = 0
	o := &mockOracle{
		proposer: func(round int64, nodeID NodeID) bool {
			return nodeID == proposer
		},
	}
	algo := New(proposer, o)
	cm, to, err := algo.StartRound(round)
	assert.Nil(t, err)
	assert.Nil(t, to)
	expected := &ConsensusMessage{
		MsgType:    Propose,
		Height:     o.Height(),
		Round:      round,
		Value:      o.value,
		ValidRound: -1,
	}
	assert.Equal(t, expected, cm)

}

type mockOracle struct {
	valid            func(v ValueID) bool
	matchingProposal func(cm *ConsensusMessage) *ConsensusMessage
	prevoteQThresh   func(round int64, value *ValueID) bool
	precommitQThresh func(round int64, value *ValueID) bool
	fThresh          func(round int64) bool
	proposer         func(round int64, nodeID NodeID) bool
	height           uint64
	value            ValueID
	valueError       error
}

func (m *mockOracle) Valid(value ValueID) bool {
	return m.valid(value)
}

func (m *mockOracle) MatchingProposal(cm *ConsensusMessage) *ConsensusMessage {
	return m.matchingProposal(cm)
}

func (m *mockOracle) PrevoteQThresh(round int64, value *ValueID) bool {
	return m.prevoteQThresh(round, value)
}

func (m *mockOracle) PrecommitQThresh(round int64, value *ValueID) bool {
	return m.precommitQThresh(round, value)
}

func (m *mockOracle) FThresh(round int64) bool {
	return m.fThresh(round)
}

func (m *mockOracle) Proposer(round int64, nodeID NodeID) bool {
	return m.proposer(round, nodeID)
}

func (m *mockOracle) Height() uint64 {
	return m.height
}

func (m *mockOracle) Value() (ValueID, error) {
	return m.value, m.valueError
}
