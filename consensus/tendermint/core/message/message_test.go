package message

import (
	"bytes"
	"errors"
	"math/big"
	"math/rand"
	"reflect"
	"testing"

	"github.com/influxdata/influxdb/pkg/deep"
	"github.com/stretchr/testify/require"

	"github.com/autonity/autonity/consensus/tendermint/core/constants"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/crypto/blst"
	"github.com/autonity/autonity/p2p"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/rlp"
)

func TestMessageDecode(t *testing.T) {
	t.Run("prevote", func(t *testing.T) {
		vote := newVote[Prevote](1, 2, common.HexToHash("0x1227"), defaultSigner)
		decoded := &Prevote{}
		reader := bytes.NewReader(vote.Payload())
		if err := rlp.Decode(reader, decoded); err != nil {
			t.Fatalf("have %v, want nil", err)
		}
		if decoded.Value() != vote.Value() {
			t.Errorf("values are not the same: have %v, want %v", decoded, vote)
		}
		if decoded.H() != vote.H() {
			t.Errorf("values are not the same: have %v, want %v", decoded, vote)
		}
		if decoded.R() != vote.R() {
			t.Errorf("values are not the same: have %v, want %v", decoded, vote)
		}
	})
	t.Run("precommit", func(t *testing.T) {
		vote := newVote[Precommit](1, 2, common.HexToHash("0x1227"), defaultSigner)
		decoded := &Precommit{}
		reader := bytes.NewReader(vote.Payload())
		if err := rlp.Decode(reader, decoded); err != nil {
			t.Fatalf("have %v, want nil", err)
		}
		if decoded.Value() != vote.Value() {
			t.Errorf("values are not the same: have %v, want %v", decoded, vote)
		}
		if decoded.H() != vote.H() {
			t.Errorf("values are not the same: have %v, want %v", decoded, vote)
		}
		if decoded.R() != vote.R() {
			t.Errorf("values are not the same: have %v, want %v", decoded, vote)
		}
	})
	t.Run("propose", func(t *testing.T) {
		header := &types.Header{Number: common.Big2}
		block := types.NewBlockWithHeader(header)
		proposal := NewPropose(1, 2, -1, block, defaultSigner)
		decoded := &Propose{}
		reader := bytes.NewReader(proposal.Payload())
		if err := rlp.Decode(reader, decoded); err != nil {
			t.Fatalf("have %v, want nil", err)
		}
		if decoded.Value() != proposal.Value() {
			t.Errorf("values are not the same: have %v, want %v", decoded, proposal)
		}
		if decoded.H() != proposal.H() {
			t.Errorf("values are not the same: have %v, want %v", decoded, proposal)
		}
		if decoded.R() != proposal.R() {
			t.Errorf("values are not the same: have %v, want %v", decoded, proposal)
		}
		if decoded.ValidRound() != proposal.ValidRound() {
			t.Errorf("values are not the same: have %v, want %v", decoded, proposal)
		}
	})
	t.Run("invalid propose with vr > r", func(t *testing.T) {
		header := &types.Header{Number: common.Big2}
		block := types.NewBlockWithHeader(header)
		proposal := NewPropose(1, 2, 57, block, defaultSigner)
		decoded := &Propose{}
		reader := bytes.NewReader(proposal.Payload())
		err := rlp.Decode(reader, decoded)
		if !errors.Is(err, constants.ErrInvalidMessage) {
			t.Error("Decoding should have failed")
		}
	})
	t.Run("invalid propose with proposal height != block number", func(t *testing.T) {
		header := &types.Header{Number: common.Big2}
		block := types.NewBlockWithHeader(header)
		proposal := NewPropose(1, 4, 57, block, defaultSigner)
		decoded := &Propose{}
		reader := bytes.NewReader(proposal.Payload())
		err := rlp.Decode(reader, decoded)
		if !errors.Is(err, constants.ErrInvalidMessage) {
			t.Error("Decoding should have failed")
		}
	})
}

func TestValidate(t *testing.T) {
	t.Run("invalid signature, error returned", func(t *testing.T) {
		lastHeader := &types.Header{Number: new(big.Int).SetUint64(25)}
		msg := newVote[Prevote](1, 25, lastHeader.Hash(), func(hash common.Hash) (signature blst.Signature, address common.Address) {
			// tamper the hash to make signature invalid
			hash[0] = 0xca
			hash[1] = 0xfe
			signature, addr := defaultSigner(hash)
			return signature, addr
		})
		err := msg.Validate(func(_ common.Address) *types.CommitteeMember {
			return &types.CommitteeMember{
				Address:      testAddr,
				VotingPower:  common.Big1,
				ConsensusKey: testConsensusKey.PublicKey(),
			}
		})
		require.ErrorIs(t, err, ErrBadSignature)
	})

	t.Run("not a committee member, error returned", func(t *testing.T) {
		lastHeader := &types.Header{Number: new(big.Int).SetUint64(25), Committee: []types.CommitteeMember{
			{
				Address:           testAddr,
				VotingPower:       big.NewInt(2),
				ConsensusKey:      testConsensusKey.PublicKey(),
				ConsensusKeyBytes: testConsensusKey.PublicKey().Marshal(),
			},
		}}
		messages := []Msg{
			newVote[Prevote](1, 25, lastHeader.Hash(), defaultSigner),
			newVote[Precommit](1, 25, lastHeader.Hash(), defaultSigner),
			NewPropose(1, 25, 2, types.NewBlockWithHeader(lastHeader), defaultSigner),
		}

		validateFn := func(address common.Address) *types.CommitteeMember { //nolint
			return nil
		}

		for i := range messages {
			err := messages[i].Validate(validateFn)
			require.ErrorIs(t, err, ErrUnauthorizedAddress)
		}
	})

	t.Run("valid params given, valid validator returned", func(t *testing.T) {
		validator := &types.CommitteeMember{
			Address:           testAddr,
			VotingPower:       big.NewInt(2),
			ConsensusKey:      testConsensusKey.PublicKey(),
			ConsensusKeyBytes: testConsensusKey.PublicKey().Marshal(),
		}
		lastHeader := &types.Header{Number: new(big.Int).SetUint64(25), Committee: []types.CommitteeMember{*validator}}
		messages := []Msg{
			newVote[Prevote](1, 25, lastHeader.Hash(), defaultSigner),
			newVote[Precommit](1, 25, lastHeader.Hash(), defaultSigner),
			NewPropose(1, 25, 2, types.NewBlockWithHeader(lastHeader), defaultSigner),
		}

		validateFn := func(address common.Address) *types.CommitteeMember { //nolint
			return validator
		}

		for i := range messages {
			err := messages[i].Validate(validateFn)
			require.NoError(t, err)
		}
	})
}

func TestMessageEncodeDecode(t *testing.T) {
	validator := &types.CommitteeMember{
		Address:           testAddr,
		VotingPower:       big.NewInt(2),
		ConsensusKey:      testConsensusKey.PublicKey(),
		ConsensusKeyBytes: testConsensusKey.PublicKey().Marshal(),
	}
	lastHeader := &types.Header{Number: new(big.Int).SetUint64(2), Committee: []types.CommitteeMember{*validator}}
	messages := []Msg{
		NewPropose(1, 2, -1, types.NewBlockWithHeader(lastHeader), defaultSigner).MustVerify(stubVerifier),
		NewPrevote(1, 2, lastHeader.Hash(), defaultSigner).MustVerify(stubVerifier),
		NewPrecommit(1, 2, lastHeader.Hash(), defaultSigner).MustVerify(stubVerifier),
	}
	for i := range messages {
		buff := new(bytes.Buffer)
		err := rlp.Encode(buff, messages[i])
		require.NoError(t, err)
		decoded := reflect.New(reflect.TypeOf(messages[i]).Elem()).Interface().(Msg)
		err = rlp.Decode(buff, decoded)
		require.NoError(t, err)
		decoded.Validate(stubVerifier)
		if decoded.Value() != messages[i].Value() ||
			decoded.R() != messages[i].R() ||
			decoded.H() != messages[i].H() ||
			decoded.Hash() != messages[i].Hash() ||
			!deep.Equal(decoded.Payload(), messages[i].Payload()) {
			t.Error("does not match", i)
		}
	}
}

func FuzzFromPayload(f *testing.F) {
	msg := NewPrevote(1, 2, common.Hash{}, defaultSigner).MustVerify(stubVerifier)
	f.Add(msg.Payload())
	f.Fuzz(func(t *testing.T, seed []byte) {
		var p Prevote
		rlp.Decode(bytes.NewReader(seed), &p)
	})
}

// TODO(lorenzo) fix this benchmark, I think the p2p MSg payload gets discarded after a couple iteration of decoding
func BenchmarkDecodeVote(b *testing.B) {
	// setup vote
	hashBytes := make([]byte, 32)
	_, err := rand.Read(hashBytes)
	if err != nil {
		b.Fatal("failed to generate random bytes: ", err)
	}
	prevote := NewPrevote(int64(15), uint64(123345), common.BytesToHash(hashBytes), defaultSigner)

	// create p2p prevote
	payload := prevote.Payload()
	r := bytes.NewReader(payload)
	size := len(payload)
	p2pPrevote := p2p.Msg{Code: 0x12, Size: uint32(size), Payload: r}

	// start the actual benchmarking
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		prevoteDec := new(Prevote)
		if err := p2pPrevote.Decode(prevoteDec); err != nil {
			b.Fatal("failed prevote decoding: ", err)
		}
	}
}
