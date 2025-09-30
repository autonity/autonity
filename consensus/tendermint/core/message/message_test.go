package message

import (
	"bytes"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"reflect"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/bft"
	"github.com/autonity/autonity/consensus/tendermint/core/constants"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/crypto"
	"github.com/autonity/autonity/crypto/blst"
	"github.com/autonity/autonity/rlp"
)

// locally created messages are considered as verified, we decode it to simulate a msgs arriving from the wire
func newUnverifiedPrevote(r int64, h uint64, value common.Hash, signer Signer, self *types.CommitteeMember, committee *types.Committee) *Prevote { //nolint
	prevote := NewPrevote(r, h, value, signer, self, committee)
	unverifiedPrevote := &Prevote{}
	reader := bytes.NewReader(prevote.Payload())
	if err := rlp.Decode(reader, unverifiedPrevote); err != nil {
		panic("cannot decode prevote: " + err.Error())
	}
	return unverifiedPrevote
}

func newUnverifiedPrecommit(r int64, h uint64, value common.Hash, signer Signer, self *types.CommitteeMember, committee *types.Committee) *Precommit {
	precommit := NewPrecommit(r, h, value, signer, self, committee)
	unverifiedPrecommit := &Precommit{}
	reader := bytes.NewReader(precommit.Payload())
	if err := rlp.Decode(reader, unverifiedPrecommit); err != nil {
		panic("cannot decode precommit: " + err.Error())
	}
	return unverifiedPrecommit
}

func newUnverifiedPropose(r int64, h uint64, vr int64, block *types.Block, signer Signer, self *types.CommitteeMember) *Propose {
	propose := NewPropose(r, h, vr, block, signer, self)
	unverifiedPropose := &Propose{}
	reader := bytes.NewReader(propose.Payload())
	if err := rlp.Decode(reader, unverifiedPropose); err != nil {
		panic("cannot decode propose: " + err.Error())
	}
	return unverifiedPropose
}

func newUnverifiedLightPropose(r int64, h uint64, vr int64, block *types.Block, signer Signer, self *types.CommitteeMember) *LightProposal {
	propose := NewPropose(r, h, vr, block, signer, self).ToLight()
	unverifiedPropose := &LightProposal{}
	reader := bytes.NewReader(propose.Payload())
	if err := rlp.Decode(reader, unverifiedPropose); err != nil {
		panic("cannot decode light proposal: " + err.Error())
	}
	return unverifiedPropose
}

// don't care about power and address when dealing with signature verification
func makeCommitteeMemberWithKey(key blst.SecretKey, index uint64) *types.CommitteeMember {
	return &types.CommitteeMember{
		Address:           testAddr,
		VotingPower:       common.Big1,
		ConsensusKey:      key.PublicKey(),
		ConsensusKeyBytes: key.PublicKey().Marshal(),
		Index:             index,
	}
}

func makeSigner(key blst.SecretKey) func(common.Hash) blst.Signature {
	return func(h common.Hash) blst.Signature {
		return key.Sign(h[:])
	}
}

func TestMessageDecode(t *testing.T) {
	t.Run("prevote", func(t *testing.T) {
		vote := newVote[Prevote](1, 2, common.HexToHash("0x1227"), defaultSigner, testCommitteeMember, &testCommittee)
		decoded := &Prevote{}
		reader := bytes.NewReader(vote.Payload())
		if err := rlp.Decode(reader, decoded); err != nil {
			t.Fatalf("have %v, want nil", err)
		}
		require.Equal(t, vote.Code(), decoded.Code())
		require.Equal(t, vote.R(), decoded.R())
		require.Equal(t, vote.H(), decoded.H())
		require.Equal(t, vote.Value(), decoded.Value())
		require.NoError(t, decoded.PreValidate(&testCommittee, false))
		require.NoError(t, decoded.Signers().Validate(&testCommittee))
		require.Equal(t, vote.Signers().Bitmap, decoded.Signers().Bitmap)
		require.Equal(t, vote.Signers().Coefficients, decoded.Signers().Coefficients)
		voteSignature, _ := vote.Signature()
		decodedSignature, _ := decoded.Signature()
		require.Equal(t, voteSignature, decodedSignature)
	})
	t.Run("precommit", func(t *testing.T) {
		vote := newVote[Precommit](1, 2, common.HexToHash("0x1227"), defaultSigner, testCommitteeMember, &testCommittee)
		decoded := &Precommit{}
		reader := bytes.NewReader(vote.Payload())
		if err := rlp.Decode(reader, decoded); err != nil {
			t.Fatalf("have %v, want nil", err)
		}
		require.Equal(t, vote.Code(), decoded.Code())
		require.Equal(t, vote.R(), decoded.R())
		require.Equal(t, vote.H(), decoded.H())
		require.Equal(t, vote.Value(), decoded.Value())
		require.NoError(t, decoded.PreValidate(&testCommittee, false))
		require.NoError(t, decoded.Signers().Validate(&testCommittee))
		require.Equal(t, vote.Signers().Bitmap, decoded.Signers().Bitmap)
		require.Equal(t, vote.Signers().Coefficients, decoded.Signers().Coefficients)
		voteSignature, _ := vote.Signature()
		decodedSignature, _ := decoded.Signature()
		require.Equal(t, voteSignature, decodedSignature)
	})
	t.Run("propose", func(t *testing.T) {
		header := &types.Header{Number: common.Big2}
		block := types.NewBlockWithHeader(header)
		proposal := NewPropose(1, 2, -1, block, defaultSigner, testCommitteeMember)
		decoded := &Propose{}
		reader := bytes.NewReader(proposal.Payload())
		err := rlp.Decode(reader, decoded)
		require.NoError(t, err)
		require.Equal(t, proposal.Code(), decoded.Code())
		require.Equal(t, proposal.R(), decoded.R())
		require.Equal(t, proposal.H(), decoded.H())
		require.Equal(t, proposal.Value(), decoded.Value())
		require.Equal(t, proposal.ValidRound(), decoded.ValidRound())
		require.Equal(t, proposal.Signer(), decoded.Signer())
		propsalSignature, _ := proposal.Signature()
		decodedSignature, _ := decoded.Signature()
		require.Equal(t, propsalSignature, decodedSignature)
	})
	t.Run("invalid propose with vr > r", func(t *testing.T) {
		header := &types.Header{Number: common.Big2}
		block := types.NewBlockWithHeader(header)
		proposal := NewPropose(1, 2, 57, block, defaultSigner, testCommitteeMember)
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
		proposal := NewPropose(1, 4, 57, block, defaultSigner, testCommitteeMember)
		decoded := &Propose{}
		reader := bytes.NewReader(proposal.Payload())
		err := rlp.Decode(reader, decoded)
		if !errors.Is(err, constants.ErrInvalidMessage) {
			t.Error("Decoding should have failed")
		}
	})
}

func TestValidate(t *testing.T) {
	key1, err := blst.RandKey()
	require.NoError(t, err)
	key2, err := blst.RandKey()
	require.NoError(t, err)

	signer1 := makeSigner(key1)
	signer2 := makeSigner(key2)

	member1 := makeCommitteeMemberWithKey(key1, 0)
	member2 := makeCommitteeMemberWithKey(key2, 1)

	committee := new(types.Committee)
	committee.Members = []types.CommitteeMember{*member1, *member2}

	header := newHeader(25, committee)

	t.Run("invalid signature, error returned", func(t *testing.T) {
		msg := newUnverifiedPrevote(1, 25, header.Hash(), func(hash common.Hash) blst.Signature {
			// tamper the hash to make signature invalid
			hash[0] = 0xca
			hash[1] = 0xfe
			return signer1(hash)
		}, member1, committee)
		err := msg.PreValidate(committee, true)
		require.NoError(t, err)
		err = msg.Validate()
		require.ErrorIs(t, err, ErrBadSignature)

		// if aggregate with a valid signature, resulting signature is still invalid
		msg2 := NewPrevote(1, 25, header.Hash(), signer2, member2, committee)
		aggregatedVote := AggregatePrevotesSingle([]Vote{msg, msg2})

		// unverify the aggregated vote
		unverifiedPrevote := &Prevote{}
		reader := bytes.NewReader(aggregatedVote.Payload())
		if err := rlp.Decode(reader, unverifiedPrevote); err != nil {
			panic("cannot decode prevote: " + err.Error())
		}
		require.False(t, unverifiedPrevote.verified)

		err = unverifiedPrevote.PreValidate(committee, true)
		require.NoError(t, err)
		err = unverifiedPrevote.Validate()
		require.ErrorIs(t, err, ErrBadSignature)
	})
	t.Run("valid signature, no error returned", func(t *testing.T) {
		msg := newUnverifiedPrevote(1, 25, header.Hash(), signer1, member1, committee)
		err := msg.PreValidate(committee, true)
		require.NoError(t, err)
		err = msg.Validate()
		require.NoError(t, err)

		// if aggregated with a valid signature, resulting signature is still valid
		msg2 := NewPrevote(1, 25, header.Hash(), signer2, member2, committee)
		aggregatedVote := AggregatePrevotesSingle([]Vote{msg, msg2})

		// unverify the aggregated vote
		unverifiedPrevote := &Prevote{}
		reader := bytes.NewReader(aggregatedVote.Payload())
		if err := rlp.Decode(reader, unverifiedPrevote); err != nil {
			panic("cannot decode prevote: " + err.Error())
		}
		require.False(t, unverifiedPrevote.verified)

		err = unverifiedPrevote.PreValidate(committee, true)
		require.NoError(t, err)
		err = unverifiedPrevote.Validate()
		require.NoError(t, err)
	})
}

func TestPreValidate(t *testing.T) {
	t.Run("proposals from not a committee member, error returned", func(t *testing.T) {
		blsKey, err := blst.RandKey()
		require.NoError(t, err)
		ecdsaKey, err := crypto.GenerateKey()
		require.NoError(t, err)
		otherCommitteeMember := types.CommitteeMember{
			Address:           crypto.PubkeyToAddress(ecdsaKey.PublicKey),
			VotingPower:       common.Big1,
			ConsensusKeyBytes: blsKey.PublicKey().Marshal(),
			ConsensusKey:      blsKey.PublicKey(),
			Index:             1,
		}
		committee := new(types.Committee)
		committee.Members = []types.CommitteeMember{otherCommitteeMember}

		header := newHeader(25, committee)
		messages := []Msg{
			newUnverifiedPropose(1, 25, -1, types.NewBlockWithHeader(header), defaultSigner, testCommitteeMember),
			newUnverifiedLightPropose(1, 25, -1, types.NewBlockWithHeader(header), defaultSigner, testCommitteeMember),
		}

		for _, message := range messages {
			err := message.PreValidate(committee, true)
			require.ErrorIs(t, err, ErrUnauthorizedAddress)
		}
	})
	t.Run("proposals from a committee member, no error", func(t *testing.T) {
		committee := new(types.Committee)
		committee.Members = []types.CommitteeMember{*testCommitteeMember}
		header := newHeader(25, committee)
		messages := []Msg{
			newUnverifiedPropose(1, 25, -1, types.NewBlockWithHeader(header), defaultSigner, testCommitteeMember),
			newUnverifiedLightPropose(1, 25, -1, types.NewBlockWithHeader(header), defaultSigner, testCommitteeMember),
		}

		for _, message := range messages {
			err := message.PreValidate(committee, true)
			require.NoError(t, err)
		}
	})

	t.Run("votes with correct signers information, no error", func(t *testing.T) {
		committee := new(types.Committee)
		committee.Members = []types.CommitteeMember{*testCommitteeMember}
		header := newHeader(25, committee)
		messages := []Msg{
			newUnverifiedPrevote(1, 25, header.Hash(), defaultSigner, testCommitteeMember, committee),
			newUnverifiedPrecommit(1, 25, header.Hash(), defaultSigner, testCommitteeMember, committee),
		}

		for _, message := range messages {
			err := message.PreValidate(committee, true)
			require.NoError(t, err)
		}
	})

	t.Run("votes with incorrect signers information, error is returned", func(t *testing.T) {
		committee := new(types.Committee)
		committee.Members = []types.CommitteeMember{*testCommitteeMember}
		header := newHeader(25, committee)
		messages := []Vote{
			newUnverifiedPrevote(1, 25, header.Hash(), defaultSigner, testCommitteeMember, committee),
			newUnverifiedPrecommit(1, 25, header.Hash(), defaultSigner, testCommitteeMember, committee),
		}

		// tamper with signers information
		messages[0].Signers().Bitmap = (*types.Bitmap)(new(big.Int).SetUint64(100))
		messages[1].Signers().Bitmap = types.NewBitmap()

		for _, message := range messages {
			err := message.PreValidate(committee, true)
			t.Log(err)
			require.Error(t, err)
		}
	})
	t.Run("votes has coefficient with bitlen > 16 and is a network vote, error is returned", func(t *testing.T) {
		n := 18
		committee := new(types.Committee)
		for i := 0; i < n; i++ {
			committee.Members = append(committee.Members, *testCommitteeMember)
		}
		header := newHeader(25, committee)
		vote := newUnverifiedPrevote(1, 25, header.Hash(), defaultSigner, testCommitteeMember, committee)

		// make this vote bloated by tweaking the signers (NOTE: this will not pass validate since the signature doesn't actually match the signers)
		vote.Signers().Coefficients = make([]*big.Int, n)
		for i := 0; i < n-1; i++ {
			vote.Signers().Bitmap.Set(i)
			vote.Signers().Coefficients[i] = new(big.Int).SetUint64(1)
		}
		vote.Signers().Bitmap.Set(n - 1)
		vote.Signers().Coefficients[n-1] = new(big.Int).SetUint64(1 << (n - 2))

		t.Log(vote.Signers().String())
		err := vote.PreValidate(committee, true)
		t.Logf("prevalidate failed with error: %v", err)
		require.True(t, errors.Is(err, types.ErrInvalidCoefficient))

		// should not error out if cap is ignored
		err = vote.PreValidate(committee, false)
		t.Logf("prevalidate failed with error: %v", err)
		require.NoError(t, err)
	})

}

func TestMessageEncodeDecode(t *testing.T) {
	committee := new(types.Committee)
	committee.Members = []types.CommitteeMember{*testCommitteeMember}
	header := newHeader(2, committee)
	messages := []Msg{
		NewPropose(1, 2, -1, types.NewBlockWithHeader(header), defaultSigner, testCommitteeMember),
		NewPrevote(1, 2, header.Hash(), defaultSigner, testCommitteeMember, committee),
		NewPrecommit(1, 2, header.Hash(), defaultSigner, testCommitteeMember, committee),
	}
	for i := range messages {
		buff := new(bytes.Buffer)
		err := rlp.Encode(buff, messages[i])
		require.NoError(t, err)
		decoded := reflect.New(reflect.TypeOf(messages[i]).Elem()).Interface().(Msg)
		err = rlp.Decode(buff, decoded)
		require.NoError(t, err)
		err = decoded.PreValidate(committee, true)
		require.NoError(t, err)
		err = decoded.Validate()
		require.NoError(t, err)
		if decoded.Value() != messages[i].Value() ||
			decoded.R() != messages[i].R() ||
			decoded.H() != messages[i].H() ||
			decoded.Code() != messages[i].Code() ||
			decoded.Value() != messages[i].Value() ||
			!bytes.Equal(decoded.SignerKey().Marshal(), messages[i].SignerKey().Marshal()) ||
			decoded.Hash() != messages[i].Hash() ||
			!bytes.Equal(decoded.Payload(), messages[i].Payload()) {
			t.Error("does not match", i)
		}
	}
}

func TestPrevoteDecodeRLP(t *testing.T) {
	key, err := blst.RandKey()
	require.NoError(t, err)
	signer := makeSigner(key)

	signatureInputPayload, err := rlp.EncodeToBytes([]any{PrevoteCode, uint64(1), uint64(2), common.HexToHash("0xdeadbeef")})
	require.NoError(t, err)
	signatureInputHash := crypto.Hash(signatureInputPayload)
	signature := signer(signatureInputHash)

	validVote := extVote{
		Code:   PrevoteCode,
		Round:  1,
		Height: 2,
		Value:  common.HexToHash("0xdeadbeef"),
		Signers: &types.Signers{
			Bitmap:       (*types.Bitmap)(big.NewInt(1)),
			Coefficients: []*big.Int{big.NewInt(1)},
		},
		Signature: signature.(*blst.BlsSignature),
	}
	validPayload, err := rlp.EncodeToBytes(&validVote)
	require.NoError(t, err)

	testCases := []struct {
		name          string
		payload       []byte
		expectedError error
	}{
		{
			name:          "valid vote",
			payload:       validPayload,
			expectedError: nil,
		},
		{
			name: "extra data at the end of list",
			payload: func() []byte {
				var list []interface{}
				require.NoError(t, rlp.DecodeBytes(validPayload, &list))
				list = append(list, "extra_field")
				payload, err := rlp.EncodeToBytes(list)
				require.NoError(t, err)
				return payload
			}(),
			expectedError: constants.ErrInvalidMessage,
		},
		{
			name: "too few items in list",
			payload: func() []byte {
				payload, err := rlp.EncodeToBytes([]interface{}{
					validVote.Code,
					validVote.Round,
					validVote.Height,
				})
				require.NoError(t, err)
				return payload
			}(),
			expectedError: constants.ErrInvalidMessage,
		},
		{
			name: "wrong data type for height",
			payload: func() []byte {
				payload, err := rlp.EncodeToBytes([]interface{}{
					validVote.Code,
					validVote.Round,
					"this should be a number",
					validVote.Value,
					validVote.Signers,
					validVote.Signature,
				})
				require.NoError(t, err)
				return payload
			}(),
			expectedError: constants.ErrInvalidMessage,
		},
		{
			name: "list where item expected",
			payload: func() []byte {
				payload, err := rlp.EncodeToBytes([]interface{}{
					validVote.Code,
					validVote.Round,
					[]interface{}{validVote.Height}, // Height is now a list
					validVote.Value,
					validVote.Signers,
					validVote.Signature,
				})
				require.NoError(t, err)
				return payload
			}(),
			expectedError: constants.ErrInvalidMessage,
		},
		// Enhanced cases for signers list
		{
			name: "signers list with too few sub-items",
			payload: func() []byte {
				payload, err := rlp.EncodeToBytes([]interface{}{
					validVote.Code,
					validVote.Round,
					validVote.Height,
					validVote.Value,
					[]interface{}{big.NewInt(1)}, // Only bitmap, missing coefficients list
					validVote.Signature,
				})
				require.NoError(t, err)
				return payload
			}(),
			expectedError: constants.ErrInvalidMessage,
		},
		{
			name: "signers list with extra sub-items",
			payload: func() []byte {
				payload, err := rlp.EncodeToBytes([]interface{}{
					validVote.Code,
					validVote.Round,
					validVote.Height,
					validVote.Value,
					[]interface{}{big.NewInt(1), []interface{}{big.NewInt(1)}, "extra_item"},
					validVote.Signature,
				})
				require.NoError(t, err)
				return payload
			}(),
			expectedError: constants.ErrInvalidMessage,
		},
		{
			name: "valid coefficients list encoded int size bigger than 1 byte",
			payload: func() []byte {
				expectedCoefficients := []*big.Int{big.NewInt(5), big.NewInt(1000)}
				fullList := []interface{}{
					validVote.Code,
					validVote.Round,
					validVote.Height,
					validVote.Value,
					[]interface{}{big.NewInt(7), expectedCoefficients},
					validVote.Signature,
				}
				payload, err := rlp.EncodeToBytes(fullList)
				require.NoError(t, err)
				return payload
			}(),
			expectedError: nil,
		},
		{
			name: "coefficients list with zero value",
			payload: func() []byte {
				payload, err := rlp.EncodeToBytes([]interface{}{
					validVote.Code,
					validVote.Round,
					validVote.Height,
					validVote.Value,
					[]interface{}{big.NewInt(1), []interface{}{big.NewInt(0)}},
					validVote.Signature,
				})
				require.NoError(t, err)
				return payload
			}(),
			expectedError: constants.ErrInvalidMessage,
		},
		{
			name: "coefficients list with oversized BitLen",
			payload: func() []byte {
				oversized := new(big.Int).Lsh(big.NewInt(1), common.QuorumCap+1)
				payload, err := rlp.EncodeToBytes([]interface{}{
					validVote.Code,
					validVote.Round,
					validVote.Height,
					validVote.Value,
					[]interface{}{big.NewInt(1), []interface{}{oversized}},
					validVote.Signature,
				})
				require.NoError(t, err)
				return payload
			}(),
			expectedError: constants.ErrInvalidMessage,
		},
		{
			name: "coefficients list mismatched with bitmap len",
			payload: func() []byte {
				payload, err := rlp.EncodeToBytes([]interface{}{
					validVote.Code,
					validVote.Round,
					validVote.Height,
					validVote.Value,
					[]interface{}{big.NewInt(1), []interface{}{big.NewInt(1), big.NewInt(1)}},
					validVote.Signature,
				})
				require.NoError(t, err)
				return payload
			}(),
			expectedError: constants.ErrInvalidMessage,
		},
		{
			name: "coefficients list with extra items",
			payload: func() []byte {
				payload, err := rlp.EncodeToBytes([]interface{}{
					validVote.Code,
					validVote.Round,
					validVote.Height,
					validVote.Value,
					[]interface{}{big.NewInt(1), []interface{}{big.NewInt(1), "extra_coeff"}},
					validVote.Signature,
				})
				require.NoError(t, err)
				return payload
			}(),
			expectedError: constants.ErrInvalidMessage,
		},
		{
			name: "coefficients list with wrong type (e.g., string instead of big.Int)",
			payload: func() []byte {
				payload, err := rlp.EncodeToBytes([]interface{}{
					validVote.Code,
					validVote.Round,
					validVote.Height,
					validVote.Value,
					[]interface{}{big.NewInt(1), []interface{}{"not_a_bigint"}},
					validVote.Signature,
				})
				require.NoError(t, err)
				return payload
			}(),
			expectedError: constants.ErrInvalidMessage,
		},
		{
			name: "coefficients list oversized ( > MaxAllowedSigners )",
			payload: func() []byte {
				coeffs := make([]interface{}, types.MaxAllowedSigners+1)
				for i := range coeffs {
					coeffs[i] = big.NewInt(1)
				}
				bitmap := new(big.Int).SetBit(new(big.Int), types.MaxAllowedSigners, 1) // Oversized bitmap too
				payload, err := rlp.EncodeToBytes([]interface{}{
					validVote.Code,
					validVote.Round,
					validVote.Height,
					validVote.Value,
					[]interface{}{bitmap, coeffs},
					validVote.Signature,
				})
				require.NoError(t, err)
				return payload
			}(),
			expectedError: constants.ErrInvalidMessage,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			prevote := &Prevote{}
			err := prevote.DecodeRLPPayload(tc.payload, crypto.Hash(tc.payload))
			if tc.expectedError == nil {
				require.NoError(t, err)
			} else {
				require.True(t, errors.Is(err, tc.expectedError), "Expected error: %v, got: %v", tc.expectedError, err)
			}
		})
	}
}

func TestPrevoteDecodeRLPPayload_Concurrent(t *testing.T) {
	key, err := blst.RandKey()
	require.NoError(t, err)
	signer := makeSigner(key)

	signatureInputPayload, err := rlp.EncodeToBytes([]any{PrevoteCode, uint64(1), uint64(2), common.HexToHash("0xdeadbeef")})
	require.NoError(t, err)
	signatureInputHash := crypto.Hash(signatureInputPayload)
	signature := signer(signatureInputHash)

	validVote := extVote{
		Code:   PrevoteCode,
		Round:  1,
		Height: 2,
		Value:  common.HexToHash("0xdeadbeef"),
		Signers: &types.Signers{
			Bitmap:       (*types.Bitmap)(big.NewInt(1)),
			Coefficients: []*big.Int{big.NewInt(1)},
		},
		Signature: signature.(*blst.BlsSignature),
	}
	validPayload, err := rlp.EncodeToBytes(&validVote)
	require.NoError(t, err)
	validHash := crypto.Hash(validPayload)

	const numGoroutines = 1000
	var wg sync.WaitGroup
	errCh := make(chan error, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			prevote := &Prevote{}
			if err := prevote.DecodeRLPPayload(validPayload, validHash); err != nil {
				errCh <- err
				return
			}
			// Basic validation to ensure no corruption
			if prevote.Code() != PrevoteCode || prevote.R() != 1 || prevote.H() != 2 {
				errCh <- errors.New("decoded fields mismatch")
				return
			}
		}()
	}

	wg.Wait()
	close(errCh)

	for err := range errCh {
		require.NoError(t, err, "Concurrent decoding failed")
	}
}

// verify that aggregating same votes in different orders doesn't change the hash
func TestMessageHash(t *testing.T) {
	h := uint64(1)
	r := int64(0)
	v := common.HexToHash("0a5843ac1c1247324a23a23f23f742f89f431293123020912dade33149f4fffe")
	vr := int64(-1)
	err := testCommittee.Enrich()
	require.NoError(t, err)
	header := newHeader(25, &testCommittee)
	block := types.NewBlockWithHeader(header)

	t.Run("Aggregating same votes in different orders yields same hash", func(t *testing.T) {
		key1, err := blst.RandKey()
		require.NoError(t, err)
		key2, err := blst.RandKey()
		require.NoError(t, err)
		key3, err := blst.RandKey()
		require.NoError(t, err)

		vote1 := NewPrevote(r, h, v, makeSigner(key1), &testCommittee.Members[0], &testCommittee)
		vote2 := NewPrevote(r, h, v, makeSigner(key2), &testCommittee.Members[1], &testCommittee)
		vote3 := NewPrevote(r, h, v, makeSigner(key3), &testCommittee.Members[2], &testCommittee)

		aggregates1 := AggregatePrevotesSingle([]Vote{vote1, vote2, vote3})
		aggregates2 := AggregatePrevotesSingle([]Vote{vote2, vote1, vote3})
		aggregates3 := AggregatePrevotesSingle([]Vote{vote3, vote1, vote2})

		hash := aggregates1.Hash()
		require.Equal(t, hash, aggregates2.Hash())
		require.Equal(t, hash, aggregates3.Hash())
	})
	t.Run("Change in the signature causes change in hash", func(t *testing.T) {
		key1, err := blst.RandKey()
		require.NoError(t, err)
		proposal := NewPropose(r, h, vr, block, defaultSigner, &testCommittee.Members[0])
		proposal2 := NewPropose(r, h, vr, block, makeSigner(key1), &testCommittee.Members[0])

		require.NotEqual(t, proposal.Hash(), proposal2.Hash())

		vote := NewPrecommit(r, h, v, defaultSigner, &testCommittee.Members[0], &testCommittee)
		vote2 := NewPrecommit(r, h, v, makeSigner(key1), &testCommittee.Members[0], &testCommittee)

		require.NotEqual(t, vote.Hash(), vote2.Hash())
	})
	t.Run("Change in the signers Bitmap and Coefficients should cause change in hash", func(t *testing.T) {
		// change signer
		vote := NewPrecommit(r, h, v, defaultSigner, &testCommittee.Members[0], &testCommittee)
		vote2 := NewPrecommit(r, h, v, defaultSigner, &testCommittee.Members[1], &testCommittee)
		require.NotEqual(t, vote.Hash(), vote2.Hash())
	})
	t.Run("internal modifications of *big.Int in signers should not cause message hash change", func(t *testing.T) {
		vote := NewPrecommit(r, h, v, defaultSigner, &testCommittee.Members[0], &testCommittee)
		originalHash := vote.Hash()
		originalIndexes := vote.Signers().Bitmap.Indexes()
		t.Logf("original indexes: %v", originalIndexes)

		// modify internals via SetBytes()
		bitmap := vote.Signers().Bitmap
		bigBitmap := (*big.Int)(bitmap)
		underlyingBytes := bigBitmap.Bytes()
		t.Logf("underlying bytes: %v", underlyingBytes)
		underlyingBytes = make([]byte, 3)
		underlyingBytes[2] = 1 // setBytes expects big endian
		t.Logf("underlying bytes modified: %v", underlyingBytes)
		bigBitmap.SetBytes(underlyingBytes)

		// signers indexes and hash should remain the same
		modifiedIndexes := bitmap.Indexes()
		t.Logf("modified indexes: %v", modifiedIndexes)
		require.Equal(t, originalIndexes, modifiedIndexes)
		require.Equal(t, originalHash, recomputeHash(vote))

		// modify internals via SetBits()
		underlyingBits := bigBitmap.Bits()
		t.Logf("underlying bits: %v", underlyingBits)
		underlyingBits = make([]big.Word, 3)
		underlyingBits[0] = 1 // set bits expects small endian
		t.Logf("underlying bits modified: %v", underlyingBits)
		bigBitmap.SetBits(underlyingBits)

		// signers indexes and hash should remain the same
		modifiedIndexes = bitmap.Indexes()
		t.Logf("modified indexes: %v", modifiedIndexes)
		require.Equal(t, originalIndexes, modifiedIndexes)
		require.Equal(t, originalHash, recomputeHash(vote))

	})
}

// helper to recompute hash since in production code it is computed only once at decoding and cached
// therefore modifying internals of a vote doesn't lead to the cached hash to change
func recomputeHash(vote Vote) common.Hash {
	extvote := extVote{
		Code:    vote.Code(),
		Round:   uint64(vote.R()), // #nosec
		Height:  vote.H(),
		Value:   vote.Value(),
		Signers: vote.Signers(),
	}
	sig, _ := vote.Signature()
	extvote.Signature = sig.(*blst.BlsSignature)
	payload, _ := rlp.EncodeToBytes(extvote)
	return crypto.Hash(payload)
}

func FuzzFromPayload(f *testing.F) {
	msg := NewPrevote(1, 2, common.Hash{}, defaultSigner, testCommitteeMember, &testCommittee)
	f.Add(msg.Payload())
	f.Fuzz(func(t *testing.T, seed []byte) {
		var p Prevote
		rlp.Decode(bytes.NewReader(seed), &p)
	})
}

func TestComplexAggregation(t *testing.T) {
	h := uint64(1)
	r := int64(0)
	v1 := common.HexToHash("0a5843ac1c1247324a23a23f23f742f89f431293123020912dade33149f4fffe")
	err := testCommittee.Enrich()
	require.NoError(t, err)
	//header := &types.Header{Committee: testCommittee}
	t.Run("aggregate same message", func(t *testing.T) {
		var votes1 []Vote
		votes1 = append(votes1, NewPrevote(r, h, v1, defaultSigner, &testCommittee.Members[0], &testCommittee))
		votes1 = append(votes1, NewPrevote(r, h, v1, defaultSigner, &testCommittee.Members[1], &testCommittee))
		aggregate1 := AggregatePrevotesSingle(votes1)
		require.Equal(t, 2, len(aggregate1.Signers().FlattenUniq()))

		var votes2 []Vote
		votes2 = append(votes2, NewPrevote(r, h, v1, defaultSigner, &testCommittee.Members[0], &testCommittee))
		votes2 = append(votes2, NewPrevote(r, h, v1, defaultSigner, &testCommittee.Members[2], &testCommittee))
		aggregate2 := AggregatePrevotesSingle(votes2)
		require.Equal(t, 2, len(aggregate2.Signers().FlattenUniq()))

		aggregate3 := AggregatePrevotesSingle([]Vote{aggregate1, aggregate2})
		require.Equal(t, 3, len(aggregate3.signers.FlattenUniq()))
		require.Equal(t, 3, len(aggregate3.signers.Coefficients))
		require.Equal(t, common.Big2.String(), aggregate3.signers.Coefficients[0].String())
		for _, c := range aggregate3.signers.Coefficients[1:] {
			require.Equal(t, common.Big1.String(), c.String())
		}
	})
}

func TestAggregateVotes(t *testing.T) {
	// Rules:
	// 1. votes are ordered by decreasing power
	// 2. a vote is aggregated to the previous ones only if it adds information (i.e. adds a new signer) and the resulting aggregate respects the boundary (max coefficient = committee size)
	h := uint64(1)
	r := int64(0)
	v := common.HexToHash("0a5843ac1c1247324a23a23f23f742f89f431293123020912dade33149f4fffe")
	err := testCommittee.Enrich()
	require.NoError(t, err)

	var votes []Vote

	//NOTE: we can use whatever signer, the aggregation functions do not verify the signature

	toBigInt := func(b *types.Bitmap) *big.Int { return (*big.Int)(b) }

	votes = append(votes, NewPrevote(r, h, v, defaultSigner, &testCommittee.Members[0], &testCommittee))
	aggregate := AggregatePrevotesSingle(votes)
	t.Log(aggregate.Signers().String())
	require.Equal(t, fmt.Sprintf("%08b", toBigInt(aggregate.Signers().Bitmap).Bytes()[0]), "00000001")
	require.NoError(t, aggregate.Signers().Validate(&testCommittee))
	require.Equal(t, votes[0].Hash(), aggregate.Hash()) // currently aggregating a single vote

	votes = append(votes, NewPrevote(r, h, v, defaultSigner, &testCommittee.Members[0], &testCommittee))
	aggregate = AggregatePrevotesSingle(votes)
	t.Log(aggregate.Signers().String())
	require.Equal(t, fmt.Sprintf("%08b", toBigInt(aggregate.Signers().Bitmap).Bytes()[0]), "00000001")
	require.NoError(t, aggregate.Signers().Validate(&testCommittee))

	votes = append(votes, NewPrevote(r, h, v, defaultSigner, &testCommittee.Members[1], &testCommittee))
	aggregate = AggregatePrevotesSingle(votes)
	t.Log(aggregate.Signers().String())
	require.Equal(t, fmt.Sprintf("%08b", toBigInt(aggregate.Signers().Bitmap).Bytes()[0]), "00000011")
	require.NoError(t, aggregate.Signers().Validate(&testCommittee))

	votes = append(votes, NewPrevote(r, h, v, defaultSigner, &testCommittee.Members[2], &testCommittee))
	aggregate = AggregatePrevotesSingle(votes)
	t.Log(aggregate.Signers().String())
	require.Equal(t, fmt.Sprintf("%08b", toBigInt(aggregate.Signers().Bitmap).Bytes()[0]), "00000111")
	require.NoError(t, aggregate.Signers().Validate(&testCommittee))

	votes = append(votes, NewPrevote(r, h, v, defaultSigner, &testCommittee.Members[3], &testCommittee))
	aggregate = AggregatePrevotesSingle(votes)
	t.Log(aggregate.Signers().String())
	require.Equal(t, fmt.Sprintf("%08b", toBigInt(aggregate.Signers().Bitmap).Bytes()[0]), "00001111")
	require.NoError(t, aggregate.Signers().Validate(&testCommittee))

	aggregate2 := AggregatePrevotesSingle([]Vote{NewPrevote(r, h, v, defaultSigner, &testCommittee.Members[0], &testCommittee), NewPrevote(r, h, v, defaultSigner, &testCommittee.Members[1], &testCommittee), NewPrevote(r, h, v, defaultSigner, &testCommittee.Members[2], &testCommittee)})
	aggregate3 := AggregatePrevotesSingle([]Vote{aggregate, aggregate2})
	t.Log(aggregate3.Signers().String())
	require.Equal(t, fmt.Sprintf("%08b", toBigInt(aggregate3.Signers().Bitmap).Bytes()[0]), "00001111")
	require.NoError(t, aggregate3.Signers().Validate(&testCommittee))

	aggregate4 := AggregatePrevotesSingle([]Vote{NewPrevote(r, h, v, defaultSigner, &testCommittee.Members[0], &testCommittee), NewPrevote(r, h, v, defaultSigner, &testCommittee.Members[1], &testCommittee)})
	aggregate5 := AggregatePrevotesSingle([]Vote{NewPrevote(r, h, v, defaultSigner, &testCommittee.Members[0], &testCommittee), NewPrevote(r, h, v, defaultSigner, &testCommittee.Members[2], &testCommittee)})
	aggregate6 := AggregatePrevotesSingle([]Vote{aggregate4, aggregate5})
	t.Log(aggregate6.Signers().String())
	require.Equal(t, fmt.Sprintf("%08b", toBigInt(aggregate6.Signers().Bitmap).Bytes()[0]), "00000111")
	require.NoError(t, aggregate6.Signers().Validate(&testCommittee))

	aggregate7 := AggregatePrevotesSingle([]Vote{aggregate6, aggregate3})
	t.Log(aggregate7.Signers().String())
	require.Equal(t, fmt.Sprintf("%08b", toBigInt(aggregate7.Signers().Bitmap).Bytes()[0]), "00001111")
	require.NoError(t, aggregate7.Signers().Validate(&testCommittee))

	// aggregate with artificially inflated contribution from index = 0. validator 0 has already the maximum coefficient.
	inflatedAggregate := AggregatePrevotesSingle([]Vote{NewPrevote(r, h, v, defaultSigner, &testCommittee.Members[0], &testCommittee), NewPrevote(r, h, v, defaultSigner, &testCommittee.Members[1], &testCommittee)})

	t.Log(inflatedAggregate.Signers().String())
	require.Equal(t, fmt.Sprintf("%08b", toBigInt(inflatedAggregate.Signers().Bitmap).Bytes()[0]), "00000011")
	require.NoError(t, inflatedAggregate.Signers().Validate(&testCommittee))

	// inflatedAggregate has quorum
	require.True(t, inflatedAggregate.Power().Cmp(bft.Quorum(testCommittee.TotalVotingPower())) >= 0)

	aggregate8 := AggregatePrevotesSingle([]Vote{NewPrevote(r, h, v, defaultSigner, &testCommittee.Members[0], &testCommittee), NewPrevote(r, h, v, defaultSigner, &testCommittee.Members[2], &testCommittee), NewPrevote(r, h, v, defaultSigner, &testCommittee.Members[3], &testCommittee)})
	require.NoError(t, aggregate8.Signers().Validate(&testCommittee))

	aggregate9 := AggregatePrevotesSingle([]Vote{aggregate8, inflatedAggregate})
	t.Log(aggregate9.Signers().String())
	require.Equal(t, fmt.Sprintf("%08b", toBigInt(aggregate9.Signers().Bitmap).Bytes()[0]), "00001111")
	require.NoError(t, aggregate9.Signers().Validate(&testCommittee))

}

func TestPower(t *testing.T) {
	r := int64(1)
	h := uint64(1)
	err := testCommittee.Enrich()
	require.NoError(t, err)
	header := newHeader(25, &testCommittee)
	block := types.NewBlockWithHeader(header)

	proposal := NewPropose(r, h, -1, block, defaultSigner, &testCommittee.Members[0])

	power := Power([]Msg{proposal})
	require.Equal(t, testCommittee.Members[0].VotingPower.Uint64(), power.Uint64())

	vote := NewPrevote(r, h, block.Hash(), defaultSigner, &testCommittee.Members[0], &testCommittee)

	power = Power([]Msg{proposal, vote})
	require.Equal(t, testCommittee.Members[0].VotingPower.Uint64(), power.Uint64())

	vote.Signers().AddSigner(0)
	vote.Signers().AddSigner(1)

	power = Power([]Msg{proposal, vote})
	require.Equal(t, testCommittee.Members[0].VotingPower.Uint64()+testCommittee.Members[1].VotingPower.Uint64(), power.Uint64())

	vote2 := NewPrecommit(r, h, block.Hash(), defaultSigner, &testCommittee.Members[0], &testCommittee)
	vote2.Signers().AddSigner(3)

	power = Power([]Msg{proposal, vote, vote2})
	require.Equal(t, testCommittee.Members[0].VotingPower.Uint64()+testCommittee.Members[1].VotingPower.Uint64()+testCommittee.Members[3].VotingPower.Uint64(), power.Uint64())

	proposal2 := NewPropose(r, h, -1, block, defaultSigner, &testCommittee.Members[3])

	power = Power([]Msg{proposal, vote, vote2, proposal2})
	require.Equal(t, testCommittee.Members[0].VotingPower.Uint64()+testCommittee.Members[1].VotingPower.Uint64()+testCommittee.Members[3].VotingPower.Uint64(), power.Uint64())

	proposal3 := NewPropose(r, h, -1, block, defaultSigner, &testCommittee.Members[2])

	power = Power([]Msg{proposal, vote, vote2, proposal2, proposal3})
	require.Equal(t, testCommittee.TotalVotingPower().Uint64(), power.Uint64())
}

func TestOverQuorumVotes(t *testing.T) {
	r := int64(1)
	h := uint64(1)
	err := testCommittee.Enrich()
	require.NoError(t, err)
	header := newHeader(25, &testCommittee)
	block := types.NewBlockWithHeader(header)
	quorum := bft.Quorum(testCommittee.TotalVotingPower())

	vote := NewPrevote(r, h, block.Hash(), defaultSigner, &testCommittee.Members[0], &testCommittee)
	require.Nil(t, OverQuorumVotes([]Msg{vote}, quorum))
	vote2 := NewPrevote(r, h, block.Hash(), defaultSigner, &testCommittee.Members[1], &testCommittee)
	result := OverQuorumVotes([]Msg{vote, vote2}, quorum)
	require.NotNil(t, result)
	require.Equal(t, 2, len(result))
	require.Equal(t, vote.Hash(), result[0].Hash())
	require.Equal(t, vote2.Hash(), result[1].Hash())

	vote.Signers().AddSigner(1)
	result = OverQuorumVotes([]Msg{vote}, quorum)
	require.NotNil(t, result)
	require.Equal(t, 1, len(result))
	require.Equal(t, vote.Hash(), result[0].Hash())

}

func BenchmarkDecodeVote(b *testing.B) {
	// setup vote
	hashBytes := make([]byte, 32)
	_, err := rand.Read(hashBytes)
	if err != nil {
		b.Fatal("failed to generate random bytes: ", err)
	}
	hash := common.BytesToHash(hashBytes)
	prevote := NewPrevote(int64(15), uint64(123345), hash, defaultSigner, testCommitteeMember, &testCommittee)
	payload := prevote.Payload()
	payloadHash := crypto.Hash(payload)

	// start the actual benchmarking
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		prevoteDec := new(Prevote)
		if err := prevoteDec.DecodeRLPPayload(payload, payloadHash); err != nil {
			b.Fatal("failed prevote decoding: ", err)
		}
	}
}

func newHeader(number uint64, c *types.Committee) *types.Header {
	header := &types.Header{Number: new(big.Int).SetUint64(number)}

	var epoch types.Epoch
	epoch.Committee = c
	epoch.PreviousEpochBlock = common.Big0
	epoch.NextEpochBlock = new(big.Int).SetUint64(number + 30)
	header.Epoch = &epoch
	if err := header.Epoch.Committee.Enrich(); err != nil {
		panic(err)
	}

	return header
}
