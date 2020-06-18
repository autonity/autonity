package messages

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/rlp"
	"github.com/stretchr/testify/require"
)

type ValidatableStub types.Block

func (v ValidatableStub) Valid() bool { return true }

func (v ValidatableStub) EncodeRLP(w io.Writer) error {
	return nil
}
func (v ValidatableStub) DecodeRLP(s *rlp.Stream) error {
	return nil
}

func TestEncodeDecodeProposal(t *testing.T) {
	buf := &bytes.Buffer{}

	id := make([]byte, rand.Intn(32))
	_, err := rand.Read(id)
	require.NoError(t, err)

	proposal := Proposal{
		Round:      int64(rand.Intn(100)),
		Height:     uint64(rand.Intn(100)),
		ValidRound: int64(rand.Intn(100)),
		NodeId:     id,
		Value:      ValidatableStub{},
	}
	fmt.Printf("type: %v\n", reflect.TypeOf(proposal.Value).Kind())
	err = rlp.Encode(buf, &proposal)
	assert.NoError(t, err)

	decodedProposal := Proposal{
		Value: ValidatableStub{},
	}
	err = rlp.Decode(buf, &decodedProposal)
	assert.NoError(t, err)

	assert.Equal(t, proposal, decodedProposal)
}
