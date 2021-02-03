package algorithm

import (
	"crypto/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newNodeID(t *testing.T) NodeID {
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
	var round int64 = 0
	value := newValue(t)

	o := &mockOracle{}

	// We are proposer, expect propose message
	algo := New(newNodeID(t), o)
	expected := &ConsensusMessage{
		MsgType:    Propose,
		Height:     o.Height(),
		Round:      round,
		Value:      value,
		ValidRound: -1,
	}
	cm, to := algo.StartRound(value, round)
	assert.Nil(t, to)
	assert.Equal(t, expected, cm)

	// We are proposer, and validValue has been set, expect propose with
	// validValue.
	algo = New(newNodeID(t), o)
	algo.validValue = newValue(t)
	expected = &ConsensusMessage{
		MsgType:    Propose,
		Height:     o.Height(),
		Round:      round,
		Value:      algo.validValue,
		ValidRound: -1,
	}
	cm, to = algo.StartRound(value, round)
	assert.Nil(t, to)
	assert.Equal(t, expected, cm)

	// We are not the proposer, expect timeout message
	algo = New(newNodeID(t), o)
	expectedTimeout := &Timeout{
		TimeoutType: Propose,
		Delay:       1,
		Height:      o.Height(),
		Round:       round,
	}
	cm, to = algo.StartRound(NilValue, round)
	assert.Nil(t, cm)
	assert.Equal(t, expectedTimeout, to)

}

func TestOnTimeoutPropose(t *testing.T) {
	o := &mockOracle{
		height: 1,
	}
	algo := New(newNodeID(t), o)
	algo.round = 1
	expected := &ConsensusMessage{
		MsgType:    Prevote,
		Height:     o.Height(),
		Round:      algo.round,
		Value:      NilValue,
		ValidRound: 0,
	}
	// The timeout round and height match the current round and height, expect
	// prevote for nilValue and algorithm step set to Prevote
	cm := algo.OnTimeoutPropose(o.Height(), algo.round)
	assert.Equal(t, expected, cm)
	assert.Equal(t, Prevote, algo.step)

	// The timeout height does not match, expect nil
	cm = algo.OnTimeoutPropose(o.Height()-1, algo.round)
	assert.Nil(t, cm)

	// The timeout round does not match, expect nil
	cm = algo.OnTimeoutPropose(o.Height(), algo.round-1)
	assert.Nil(t, cm)
}

func TestOnTimeoutPrevote(t *testing.T) {
	o := &mockOracle{
		height: 1,
	}
	algo := New(newNodeID(t), o)
	algo.round = 1
	algo.step = Prevote
	expected := &ConsensusMessage{
		MsgType:    Precommit,
		Height:     o.Height(),
		Round:      algo.round,
		Value:      NilValue,
		ValidRound: 0,
	}

	// The timeout round and height match the current round and height, expect
	// precommit for nilValue and algorithm step set to Precommit
	cm := algo.OnTimeoutPrevote(o.Height(), algo.round)
	assert.Equal(t, expected, cm)
	assert.Equal(t, Precommit, algo.step)

	// The timeout height does not match, expect nil
	cm = algo.OnTimeoutPrevote(o.Height()-1, algo.round)
	assert.Nil(t, cm)

	// The timeout round does not match, expect nil
	cm = algo.OnTimeoutPrevote(o.Height(), algo.round-1)
	assert.Nil(t, cm)
}

func TestOnTimeoutPrecommit(t *testing.T) {
	o := &mockOracle{
		height: 1,
	}
	algo := New(newNodeID(t), o)
	algo.round = 1
	expected := &RoundChange{
		Round: algo.round + 1,
	}

	// The timeout round and height match the current round and height, expect
	// round change message for next round with nil decision.
	cm := algo.OnTimeoutPrecommit(o.Height(), algo.round)
	assert.Equal(t, expected, cm)

	// The timeout height does not match, expect nil
	cm = algo.OnTimeoutPrecommit(o.Height()-1, algo.round)
	assert.Nil(t, cm)

	// The timeout round does not match, expect nil
	cm = algo.OnTimeoutPrecommit(o.Height(), algo.round-1)
	assert.Nil(t, cm)
}

// Handling a proposal message for a new value
func TestReceiveMessageLine22(t *testing.T) {
	// proposer := newNodeID(t)
	var height uint64 = 1
	var round int64 = 1
	newValueProposal := &ConsensusMessage{
		MsgType:    Propose,
		Height:     height,
		Round:      round,
		Value:      newValue(t),
		ValidRound: -1,
	}
	o := &mockOracle{
		height: height,
		matchingProposal: func(cm *ConsensusMessage) *ConsensusMessage {
			return newValueProposal
		},
		valid: func(v ValueID) bool {
			return v == newValueProposal.Value
		},
	}

	algo := New(newNodeID(t), o)
	algo.round = round
	// We haven't locked a round or a value, so we expect to prevote for the
	// proposal.
	rc, cm, to := algo.ReceiveMessage(newValueProposal)
	assert.Nil(t, rc)
	assert.Nil(t, to)

	expected := &ConsensusMessage{
		MsgType:    Prevote,
		Height:     height,
		Round:      round,
		Value:      newValueProposal.Value,
		ValidRound: 0,
	}

	assert.Equal(t, expected, cm)
}

type mockOracle struct {
	valid            func(v ValueID) bool
	matchingProposal func(cm *ConsensusMessage) *ConsensusMessage
	prevoteQThresh   func(round int64, value *ValueID) bool
	precommitQThresh func(round int64, value *ValueID) bool
	fThresh          func(round int64) bool
	height           uint64
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

func (m *mockOracle) Height() uint64 {
	return m.height
}
