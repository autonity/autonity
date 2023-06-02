package message

import (
	"bytes"
	"math/big"
	"reflect"
	"testing"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/crypto"
	"github.com/autonity/autonity/rlp"
	"github.com/influxdata/influxdb/pkg/deep"
	"github.com/stretchr/testify/require"
)

var (
	key, _  = crypto.GenerateKey()
	address = crypto.PubkeyToAddress(key.PublicKey)
	signer  = func(hash common.Hash) ([]byte, common.Address) {
		out, _ := crypto.Sign(hash[:], key)
		return out, address
	}
)

func TestMessageDecode(t *testing.T) {
	t.Run("prevote", func(t *testing.T) {
		vote := newVote[Prevote](1, 2, common.HexToHash("0x1227"), defaultSigner)
		decrypted := &Prevote{}
		reader := bytes.NewReader(vote.Payload())
		if err := rlp.Decode(reader, decrypted); err != nil {
			t.Fatalf("have %v, want nil", err)
		}
		if decrypted.Value() != vote.Value() {
			t.Errorf("values are not the same: have %v, want %v", decrypted, vote)
		}
		if decrypted.H() != vote.H() {
			t.Errorf("values are not the same: have %v, want %v", decrypted, vote)
		}
		if decrypted.R() != vote.R() {
			t.Errorf("values are not the same: have %v, want %v", decrypted, vote)
		}
	})
	t.Run("precommit", func(t *testing.T) {
		vote := newVote[Precommit](1, 2, common.HexToHash("0x1227"), defaultSigner)
		decrypted := &Precommit{}
		reader := bytes.NewReader(vote.Payload())
		if err := rlp.Decode(reader, decrypted); err != nil {
			t.Fatalf("have %v, want nil", err)
		}
		if decrypted.Value() != vote.Value() {
			t.Errorf("values are not the same: have %v, want %v", decrypted, vote)
		}
		if decrypted.H() != vote.H() {
			t.Errorf("values are not the same: have %v, want %v", decrypted, vote)
		}
		if decrypted.R() != vote.R() {
			t.Errorf("values are not the same: have %v, want %v", decrypted, vote)
		}
	})
	t.Run("propose", func(t *testing.T) {
		//todo: ...
	})
}

func TestValidate(t *testing.T) {
	t.Run("invalid signature, error returned", func(t *testing.T) {
		lastHeader := &types.Header{Number: new(big.Int).SetUint64(25)}
		msg := newVote[Prevote](1, 25, lastHeader.Hash(), func(hash common.Hash) (signature []byte, address common.Address) {
			out, addr := defaultSigner(hash)
			out = append(out, 1)
			return out, addr
		})
		err := msg.Validate(func(_ common.Address) *types.CommitteeMember {
			return nil
		})
		require.ErrorIs(t, err, ErrBadSignature)
	})

	t.Run("not a committee member, error returned", func(t *testing.T) {
		lastHeader := &types.Header{Number: new(big.Int).SetUint64(25), Committee: []types.CommitteeMember{
			{
				Address:     address,
				VotingPower: big.NewInt(2),
			},
		}}
		messages := []Msg{
			newVote[Prevote](1, 25, lastHeader.Hash(), signer),
			newVote[Precommit](1, 25, lastHeader.Hash(), signer),
			NewPropose(1, 25, 2, types.NewBlockWithHeader(lastHeader), signer),
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
			Address:     address,
			VotingPower: big.NewInt(2),
		}
		lastHeader := &types.Header{Number: new(big.Int).SetUint64(25), Committee: []types.CommitteeMember{*validator}}
		messages := []Msg{
			newVote[Prevote](1, 25, lastHeader.Hash(), signer),
			newVote[Precommit](1, 25, lastHeader.Hash(), signer),
			NewPropose(1, 25, 2, types.NewBlockWithHeader(lastHeader), signer),
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
		Address:     address,
		VotingPower: big.NewInt(2),
	}
	lastHeader := &types.Header{Number: new(big.Int).SetUint64(25), Committee: []types.CommitteeMember{*validator}}
	messages := []Msg{
		NewPropose(1, 2, -1, types.NewBlockWithHeader(lastHeader), signer).MustVerify(stubVerifier),
		NewPrevote(1, 2, lastHeader.Hash(), signer).MustVerify(stubVerifier),
		NewPrecommit(1, 2, lastHeader.Hash(), signer).MustVerify(stubVerifier),
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
	msg := NewPrevote(1, 2, common.Hash{}, signer).MustVerify(stubVerifier)
	f.Add(msg.Payload())
	f.Fuzz(func(t *testing.T, seed []byte) {
		var p Prevote
		rlp.Decode(bytes.NewReader(seed), &p)
	})

}
