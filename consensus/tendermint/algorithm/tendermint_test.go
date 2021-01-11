package algorithm

import (
	"crypto/rand"
	"errors"
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

func newValue(t *testing.T) ValueID {
	var value ValueID
	_, err := rand.Read(value[:])
	require.NoError(t, err)
	return value
}

func TestStartRound(t *testing.T) {
	proposer := newNodeId(t)
	var round int64 = 0

	o := &mockOracle{
		proposer: func(round int64, nodeID NodeID) bool {
			return nodeID == proposer
		},
	}

	// We are proposer, expect propose message
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

	// We are proposer, and validValue has been set, expect propose with
	// validValue.
	algo = New(proposer, o)
	algo.validValue = newValue(t)
	cm, to, err = algo.StartRound(round)
	assert.Nil(t, err)
	assert.Nil(t, to)
	expected = &ConsensusMessage{
		MsgType:    Propose,
		Height:     o.Height(),
		Round:      round,
		Value:      algo.validValue,
		ValidRound: -1,
	}
	assert.Equal(t, expected, cm)

	// We are proposer but oracle value returns an error, expect the error
	o.valueError = errors.New("")
	algo = New(proposer, o)
	cm, to, err = algo.StartRound(round)
	assert.Nil(t, cm)
	assert.Nil(t, to)
	assert.Error(t, err)

	// We are not the proposer, expect timeout message
	algo = New(newNodeId(t), o)
	cm, to, err = algo.StartRound(round)
	assert.Nil(t, err)
	assert.Nil(t, cm)
	expectedTimeout := &Timeout{
		TimeoutType: Propose,
		Delay:       1,
		Height:      o.Height(),
		Round:       round,
	}
	assert.Equal(t, expectedTimeout, to)

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
