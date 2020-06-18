package messages

import (
	"bytes"
	"io"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/rlp"
	"github.com/stretchr/testify/require"
)

type ValidatableStub types.Block

func (v ValidatableStub) Valid() bool { return true }

func (v ValidatableStub) EncodeRLP(w io.Writer)   {}
func (v ValidatableStub) DecodeRLP(s *rlp.Stream) {}

func TestEncodeDecodeProposal(t *testing.T) {
	buf := &bytes.Buffer{}

	id := make([]byte, rand.Intn(32))
	_, err := rand.Read(id)
	require.NoError(t, err)

	proposal := Proposal{
		Round:      int64(rand.Intn(100)),
		Height:     uint64(rand.Intn(100)),
		ValidRound: int64(rand.Intn(100) - 1),
		NodeId:     id,
		Value:      ValidatableStub{},
	}
	err = rlp.Encode(buf, &proposal)
	assert.NoError(t, err)

	decodedProposal := Proposal{
		Value: ValidatableStub{},
	}
	err = rlp.Decode(buf, &decodedProposal)
	assert.NoError(t, err)

	assert.Equal(t, proposal, decodedProposal)
}
