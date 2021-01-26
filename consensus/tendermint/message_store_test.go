package tendermint

import (
	"crypto/rand"
	"crypto/sha256"
	"math/big"
	"testing"

	"github.com/clearmatics/autonity/consensus/tendermint/algorithm"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetMatchingProposal(t *testing.T) {
	m := newMessageStore(&bounds{
		centre: 0,
		high:   5,
		low:    5,
	})

	k, err := crypto.GenerateKey()
	require.NoError(t, err)
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

	msgBytes, err := encodeSignedMessage(p, k, m)
	require.NoError(t, err)

	msg, err := decodeSignedMessage(msgBytes)
	require.NoError(t, err)

	err = m.addMessage(msg, msgBytes)
	require.NoError(t, err)

	pv := &algorithm.ConsensusMessage{
		MsgType: algorithm.Prevote,
		Height:  1,
		Round:   0,
		Value:   algorithm.ValueID(h),
	}
	retrieved := m.matchingProposal(pv)
	assert.Equal(t, msg, retrieved)
}

func BenchmarkKeccak256(b *testing.B) {
	data := make([]byte, 100)
	for i := 0; i < b.N; i++ {
		_, err := rand.Read(data)
		require.NoError(b, err)
		b := crypto.Keccak256(data)
		if 0 == 1 {
			println(b)
		}
	}
}

func BenchmarkSha256(b *testing.B) {
	data := make([]byte, 100)
	for i := 0; i < b.N; i++ {
		_, err := rand.Read(data)
		require.NoError(b, err)
		b := sha256.Sum256(data)
		if 0 == 1 {
			println(b)
		}
	}
}
