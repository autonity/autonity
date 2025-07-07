package message

import (
	"bytes"
	"errors"
	"fmt"
	"math/big"
	"math/rand"
	"reflect"
	"testing"
	"time"

	"github.com/influxdata/influxdb/pkg/deep"
	"github.com/stretchr/testify/require"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/bft"
	"github.com/autonity/autonity/consensus/tendermint/core/constants"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/crypto"
	"github.com/autonity/autonity/crypto/blst"
	"github.com/autonity/autonity/p2p"
	"github.com/autonity/autonity/rlp"
)

// locally created messages are considered as verified, we decode it to simulate a msgs arriving from the wire
func newUnverifiedPrevote(r int64, h uint64, value common.Hash, signer Signer, self *types.CommitteeMember, csize int) *Prevote { //nolint
	prevote := NewPrevote(r, h, value, signer, self, csize)
	unverifiedPrevote := &Prevote{}
	reader := bytes.NewReader(prevote.Payload())
	if err := rlp.Decode(reader, unverifiedPrevote); err != nil {
		panic("cannot decode prevote: " + err.Error())
	}
	return unverifiedPrevote
}

func newUnverifiedPrecommit(r int64, h uint64, value common.Hash, signer Signer, self *types.CommitteeMember, csize int) *Precommit {
	precommit := NewPrecommit(r, h, value, signer, self, csize)
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
	t.Run("evidence", func(t *testing.T) {
		vote := newVote[Prevote](1, 2, common.HexToHash("0x1227"), defaultSigner, testCommitteeMember, 1).ToEvidence()
		decoded := &EvidenceVote{}
		reader := bytes.NewReader(vote.Payload())
		if err := rlp.Decode(reader, decoded); err != nil {
			t.Fatalf("have %v, want nil", err)
		}
		require.Equal(t, vote.Code(), decoded.Code())
		require.Equal(t, vote.R(), decoded.R())
		require.Equal(t, vote.H(), decoded.H())
		require.Equal(t, vote.Value(), decoded.Value())
		require.Equal(t, vote.Signers().Bits, decoded.Signers().Bits)
		require.Equal(t, vote.Signers().Coefficients, decoded.Signers().Coefficients)
		require.Equal(t, vote.Signature(), decoded.Signature())
	})
	t.Run("prevote", func(t *testing.T) {
		vote := newVote[Prevote](1, 2, common.HexToHash("0x1227"), defaultSigner, testCommitteeMember, 1)
		decoded := &Prevote{}
		reader := bytes.NewReader(vote.Payload())
		if err := rlp.Decode(reader, decoded); err != nil {
			t.Fatalf("have %v, want nil", err)
		}
		require.Equal(t, vote.Code(), decoded.Code())
		require.Equal(t, vote.R(), decoded.R())
		require.Equal(t, vote.H(), decoded.H())
		require.Equal(t, vote.Value(), decoded.Value())
		//require.Equal(t, vote.Signers().Bits, decoded.Signers().Bits)
		//require.Equal(t, vote.Signers().Coefficients, decoded.Signers().Coefficients)
		require.Equal(t, vote.Signature(), decoded.Signature())
	})
	t.Run("precommit", func(t *testing.T) {
		vote := newVote[Precommit](1, 2, common.HexToHash("0x1227"), defaultSigner, testCommitteeMember, 1)
		decoded := &Precommit{}
		reader := bytes.NewReader(vote.Payload())
		if err := rlp.Decode(reader, decoded); err != nil {
			t.Fatalf("have %v, want nil", err)
		}
		require.Equal(t, vote.Code(), decoded.Code())
		require.Equal(t, vote.R(), decoded.R())
		require.Equal(t, vote.H(), decoded.H())
		require.Equal(t, vote.Value(), decoded.Value())
		//require.Equal(t, vote.Signers().Bits, decoded.Signers().Bits)
		//require.Equal(t, vote.Signers().Coefficients, decoded.Signers().Coefficients)
		require.Equal(t, vote.Signature(), decoded.Signature())
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
		require.Equal(t, proposal.Signature(), decoded.Signature())
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
	csize := committee.Len()

	header := newHeader(25, committee)

	t.Run("invalid signature, error returned", func(t *testing.T) {
		msg := newUnverifiedPrevote(1, 25, header.Hash(), func(hash common.Hash) blst.Signature {
			// tamper the hash to make signature invalid
			hash[0] = 0xca
			hash[1] = 0xfe
			return signer1(hash)
		}, member1, csize)
		err := msg.PreValidate(committee)
		require.NoError(t, err)
		err = msg.Validate()
		require.ErrorIs(t, err, ErrBadSignature)

		// if aggregate with a valid signature, resulting signature is still invalid
		msg2 := NewPrevote(1, 25, header.Hash(), signer2, member2, csize)
		aggregatedVote := aggregatePrevotes([]Vote{msg, msg2})

		// unverify the aggregated vote
		unverifiedPrevote := &Prevote{}
		reader := bytes.NewReader(aggregatedVote.Payload())
		if err := rlp.Decode(reader, unverifiedPrevote); err != nil {
			panic("cannot decode prevote: " + err.Error())
		}
		require.False(t, unverifiedPrevote.verified)

		err = unverifiedPrevote.PreValidate(committee)
		require.NoError(t, err)
		err = unverifiedPrevote.Validate()
		require.ErrorIs(t, err, ErrBadSignature)
	})
	t.Run("valid signature, no error returned", func(t *testing.T) {
		msg := newUnverifiedPrevote(1, 25, header.Hash(), signer1, member1, csize)
		err := msg.PreValidate(committee)
		require.NoError(t, err)
		err = msg.Validate()
		require.NoError(t, err)

		// if aggregated with a valid signature, resulting signature is still valid
		msg2 := NewPrevote(1, 25, header.Hash(), signer2, member2, csize)
		aggregatedVote := aggregatePrevotes([]Vote{msg, msg2})

		// unverify the aggregated vote
		unverifiedPrevote := &Prevote{}
		reader := bytes.NewReader(aggregatedVote.Payload())
		if err := rlp.Decode(reader, unverifiedPrevote); err != nil {
			panic("cannot decode prevote: " + err.Error())
		}
		require.False(t, unverifiedPrevote.verified)

		err = unverifiedPrevote.PreValidate(committee)
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
			err := message.PreValidate(committee)
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
			err := message.PreValidate(committee)
			require.NoError(t, err)
		}
	})

	t.Run("votes with correct signers information, no error", func(t *testing.T) {
		committee := new(types.Committee)
		committee.Members = []types.CommitteeMember{*testCommitteeMember}
		header := newHeader(25, committee)
		messages := []Msg{
			newUnverifiedPrevote(1, 25, header.Hash(), defaultSigner, testCommitteeMember, 1),
			newUnverifiedPrecommit(1, 25, header.Hash(), defaultSigner, testCommitteeMember, 1),
		}

		for _, message := range messages {
			err := message.PreValidate(committee)
			require.NoError(t, err)
		}
	})

	t.Run("votes with incorrect signers information, error is returned", func(t *testing.T) {
		committee := new(types.Committee)
		committee.Members = []types.CommitteeMember{*testCommitteeMember}
		header := newHeader(25, committee)
		messages := []Vote{
			newUnverifiedPrevote(1, 25, header.Hash(), defaultSigner, testCommitteeMember, 1),
			newUnverifiedPrecommit(1, 25, header.Hash(), defaultSigner, testCommitteeMember, 1),
		}

		// tamper with signers information
		//messages[0].Signers().Bits = make([]byte, 100)
		//messages[1].Signers().Bits = make([]byte, 0)

		for _, message := range messages {
			err := message.PreValidate(committee)
			require.Error(t, err)
		}
	})

}

func TestMessageEncodeDecode(t *testing.T) {
	committee := new(types.Committee)
	committee.Members = []types.CommitteeMember{*testCommitteeMember}
	header := newHeader(2, committee)
	messages := []Msg{
		NewPropose(1, 2, -1, types.NewBlockWithHeader(header), defaultSigner, testCommitteeMember),
		NewPrevote(1, 2, header.Hash(), defaultSigner, testCommitteeMember, 1),
		NewPrecommit(1, 2, header.Hash(), defaultSigner, testCommitteeMember, 1),
		NewPrevote(1, 2, header.Hash(), defaultSigner, testCommitteeMember, 1).ToEvidence(),
	}
	for i := range messages {
		buff := new(bytes.Buffer)
		err := rlp.Encode(buff, messages[i])
		require.NoError(t, err)
		decoded := reflect.New(reflect.TypeOf(messages[i]).Elem()).Interface().(Msg)
		err = rlp.Decode(buff, decoded)
		require.NoError(t, err)
		err = decoded.PreValidate(committee)
		require.NoError(t, err)
		err = decoded.Validate()
		require.NoError(t, err)
		if decoded.Value() != messages[i].Value() ||
			decoded.R() != messages[i].R() ||
			decoded.H() != messages[i].H() ||
			decoded.Code() != messages[i].Code() ||
			decoded.Value() != messages[i].Value() ||
			!deep.Equal(decoded.SignerKey().Marshal(), messages[i].SignerKey().Marshal()) ||
			decoded.Hash() != messages[i].Hash() ||
			!deep.Equal(decoded.Payload(), messages[i].Payload()) {
			t.Error("does not match", i)
		}
	}
}

// verify that aggregating same votes in different orders doesn't change the hash
func TestMessageHash(t *testing.T) {
	h := uint64(1)
	r := int64(0)
	v := common.HexToHash("0a5843ac1c1247324a23a23f23f742f89f431293123020912dade33149f4fffe")
	vr := int64(-1)
	csize := testCommittee.Len()
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

		vote1 := NewPrevote(r, h, v, makeSigner(key1), &testCommittee.Members[0], csize)
		vote2 := NewPrevote(r, h, v, makeSigner(key2), &testCommittee.Members[1], csize)
		vote3 := NewPrevote(r, h, v, makeSigner(key3), &testCommittee.Members[2], csize)

		aggregates1 := aggregatePrevotes([]Vote{vote1, vote2, vote3})
		aggregates2 := aggregatePrevotes([]Vote{vote2, vote1, vote3})
		aggregates3 := aggregatePrevotes([]Vote{vote3, vote1, vote2})

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

		vote := NewPrecommit(r, h, v, defaultSigner, &testCommittee.Members[0], csize)
		vote2 := NewPrecommit(r, h, v, makeSigner(key1), &testCommittee.Members[0], csize)

		require.NotEqual(t, vote.Hash(), vote2.Hash())
	})
	t.Run("Change in the signers Bits and Coefficients should cause change in hash", func(t *testing.T) {
		// change signer
		vote := NewPrecommit(r, h, v, defaultSigner, &testCommittee.Members[0], csize)
		vote2 := NewPrecommit(r, h, v, defaultSigner, &testCommittee.Members[1], csize)
		require.NotEqual(t, vote.Hash(), vote2.Hash())

		// change committee size
		vote = NewPrecommit(r, h, v, defaultSigner, &testCommittee.Members[0], csize)
		vote2 = NewPrecommit(r, h, v, defaultSigner, &testCommittee.Members[0], csize+10)
		require.NotEqual(t, vote.Hash(), vote2.Hash())
	})
	t.Run("Change in the signers auxiliary data structures should NOT cause change in hash", func(t *testing.T) {
		// change signer
		vote := NewPrecommit(r, h, v, defaultSigner, &testCommittee.Members[0], csize)
		vote2 := NewPrecommit(r, h, v, defaultSigner, &testCommittee.Members[0], csize)

		// tamper with internal signers data structures of vote2 and recompute hash
		signers := types.NewVoteSigners(csize)
		signers.AddMember(&testCommittee.Members[0])
		tamperedPower := make(map[int]*big.Int)
		tamperedPower[123] = big.NewInt(1234)
		signers.AssignPower(tamperedPower, big.NewInt(223423))

		payload, _ := rlp.EncodeToBytes(extVote[uint16]{
			Code:       PrecommitCode,
			Round:      uint64(r), // #nosec
			Height:     h,
			Value:      v,
			Signers:    signers.SignersBase,
			Originator: testCommittee.Members[0].Address,
			Signature:  vote2.Signature().(*blst.BlsSignature),
		})
		vote2.hash = crypto.Hash(payload)

		require.Equal(t, vote.Hash(), vote2.Hash())
	})
}

func FuzzFromPayload(f *testing.F) {
	msg := NewPrevote(1, 2, common.Hash{}, defaultSigner, testCommitteeMember, 1)
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
	csize := len(testCommittee.Members)
	err := testCommittee.Enrich()
	require.NoError(t, err)
	//header := &types.Header{Committee: testCommittee}
	t.Run("aggregate same message", func(t *testing.T) {
		var votes1 []Vote
		votes1 = append(votes1, NewPrevote(r, h, v1, defaultSigner, &testCommittee.Members[0], csize))
		votes1 = append(votes1, NewPrevote(r, h, v1, defaultSigner, &testCommittee.Members[1], csize))
		aggregate1 := aggregatePrevotes(votes1)
		require.Equal(t, 2, len(aggregate1.Signers().FlattenUniq()))

		var votes2 []Vote
		votes2 = append(votes2, NewPrevote(r, h, v1, defaultSigner, &testCommittee.Members[0], csize))
		votes2 = append(votes2, NewPrevote(r, h, v1, defaultSigner, &testCommittee.Members[2], csize))
		aggregate2 := aggregatePrevotes(votes2)
		require.Equal(t, 2, len(aggregate2.Signers().FlattenUniq()))

		aggregate3 := aggregatePrevotes([]Vote{aggregate1, aggregate2})
		require.Equal(t, 3, len(aggregate3.signers.FlattenUniq()))
		require.Equal(t, 3, len(aggregate3.signers.Coefficients))
		require.Equal(t, uint16(2), aggregate3.signers.Coefficients[0])
		for _, c := range aggregate3.signers.Coefficients[1:] {
			require.Equal(t, uint16(1), c)
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
	csize := testCommittee.Len()
	err := testCommittee.Enrich()
	require.NoError(t, err)

	var votes []Vote

	//NOTE: we can use whatever signer, the aggregation functions do not verify the signature

	votes = append(votes, NewPrevote(r, h, v, defaultSigner, &testCommittee.Members[0], csize))
	aggregate := aggregatePrevotes(votes)
	t.Log(aggregate.Signers().String())
	require.Equal(t, fmt.Sprintf("%08b", aggregate.Signers().Bits[0]), "00000001")
	require.NoError(t, aggregate.Signers().Validate(csize))

	votes = append(votes, NewPrevote(r, h, v, defaultSigner, &testCommittee.Members[0], csize))
	aggregate = aggregatePrevotes(votes)
	t.Log(aggregate.Signers().String())
	require.Equal(t, fmt.Sprintf("%08b", aggregate.Signers().Bits[0]), "00000001")
	require.NoError(t, aggregate.Signers().Validate(csize))

	votes = append(votes, NewPrevote(r, h, v, defaultSigner, &testCommittee.Members[1], csize))
	aggregate = aggregatePrevotes(votes)
	t.Log(aggregate.Signers().String())
	require.Equal(t, fmt.Sprintf("%08b", aggregate.Signers().Bits[0]), "00000011")
	require.NoError(t, aggregate.Signers().Validate(csize))

	votes = append(votes, NewPrevote(r, h, v, defaultSigner, &testCommittee.Members[2], csize))
	aggregate = aggregatePrevotes(votes)
	t.Log(aggregate.Signers().String())
	require.Equal(t, fmt.Sprintf("%08b", aggregate.Signers().Bits[0]), "00000111")
	require.NoError(t, aggregate.Signers().Validate(csize))

	votes = append(votes, NewPrevote(r, h, v, defaultSigner, &testCommittee.Members[3], csize))
	aggregate = aggregatePrevotes(votes)
	t.Log(aggregate.Signers().String())
	require.Equal(t, fmt.Sprintf("%08b", aggregate.Signers().Bits[0]), "00001111")
	require.NoError(t, aggregate.Signers().Validate(csize))

	aggregate2 := aggregatePrevotes([]Vote{NewPrevote(r, h, v, defaultSigner, &testCommittee.Members[0], csize), NewPrevote(r, h, v, defaultSigner, &testCommittee.Members[1], csize), NewPrevote(r, h, v, defaultSigner, &testCommittee.Members[2], csize)})
	aggregate3 := aggregatePrevotes([]Vote{aggregate, aggregate2})
	t.Log(aggregate3.Signers().String())
	require.Equal(t, fmt.Sprintf("%08b", aggregate3.Signers().Bits[0]), "00001111")
	require.NoError(t, aggregate3.Signers().Validate(csize))

	aggregate4 := aggregatePrevotes([]Vote{NewPrevote(r, h, v, defaultSigner, &testCommittee.Members[0], csize), NewPrevote(r, h, v, defaultSigner, &testCommittee.Members[1], csize)})
	aggregate5 := aggregatePrevotes([]Vote{NewPrevote(r, h, v, defaultSigner, &testCommittee.Members[0], csize), NewPrevote(r, h, v, defaultSigner, &testCommittee.Members[2], csize)})
	aggregate6 := aggregatePrevotes([]Vote{aggregate4, aggregate5})
	t.Log(aggregate6.Signers().String())
	require.Equal(t, fmt.Sprintf("%08b", aggregate6.Signers().Bits[0]), "00000111")
	require.NoError(t, aggregate6.Signers().Validate(csize))

	aggregate7 := aggregatePrevotes([]Vote{aggregate6, aggregate3})
	t.Log(aggregate7.Signers().String())
	require.Equal(t, fmt.Sprintf("%08b", aggregate7.Signers().Bits[0]), "00001111")
	require.NoError(t, aggregate7.Signers().Validate(csize))

	// aggregate with artificially inflated contribution from index = 0. validator 0 has already the maximum coefficient.
	inflatedAggregate := aggregatePrevotes([]Vote{NewPrevote(r, h, v, defaultSigner, &testCommittee.Members[0], csize), NewPrevote(r, h, v, defaultSigner, &testCommittee.Members[1], csize)})

	t.Log(inflatedAggregate.Signers().String())
	require.Equal(t, fmt.Sprintf("%08b", inflatedAggregate.Signers().Bits[0]), "00000011")
	require.NoError(t, inflatedAggregate.Signers().Validate(csize))

	// inflatedAggregate has quorum
	require.True(t, inflatedAggregate.Power().Cmp(bft.Quorum(testCommittee.TotalVotingPower())) >= 0)

	aggregate8 := aggregatePrevotes([]Vote{NewPrevote(r, h, v, defaultSigner, &testCommittee.Members[0], csize), NewPrevote(r, h, v, defaultSigner, &testCommittee.Members[2], csize), NewPrevote(r, h, v, defaultSigner, &testCommittee.Members[3], csize)})
	require.NoError(t, aggregate8.Signers().Validate(csize))

	aggregate9 := aggregatePrevotes([]Vote{aggregate8, inflatedAggregate})
	t.Log(aggregate9.Signers().String())
	require.Equal(t, fmt.Sprintf("%08b", aggregate9.Signers().Bits[0]), "00001111")
	require.NoError(t, aggregate9.Signers().Validate(csize))

}

func TestPower(t *testing.T) {
	r := int64(1)
	h := uint64(1)
	err := testCommittee.Enrich()
	require.NoError(t, err)
	header := newHeader(25, &testCommittee)
	block := types.NewBlockWithHeader(header)
	csize := testCommittee.Len()

	proposal := NewPropose(r, h, -1, block, defaultSigner, &testCommittee.Members[0])

	power := Power([]Msg{proposal})
	require.Equal(t, testCommittee.Members[0].VotingPower.Uint64(), power.Uint64())

	vote := NewPrevote(r, h, block.Hash(), defaultSigner, &testCommittee.Members[0], csize)

	power = Power([]Msg{proposal, vote})
	require.Equal(t, testCommittee.Members[0].VotingPower.Uint64(), power.Uint64())

	vote.Signers().AddMember(&testCommittee.Members[0])
	vote.Signers().AddMember(&testCommittee.Members[1])

	power = Power([]Msg{proposal, vote})
	require.Equal(t, testCommittee.Members[0].VotingPower.Uint64()+testCommittee.Members[1].VotingPower.Uint64(), power.Uint64())

	vote2 := NewPrecommit(r, h, block.Hash(), defaultSigner, &testCommittee.Members[0], csize)
	vote2.Signers().AddMember(&testCommittee.Members[3])

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
	csize := testCommittee.Len()
	quorum := bft.Quorum(testCommittee.TotalVotingPower())

	vote := NewPrevote(r, h, block.Hash(), defaultSigner, &testCommittee.Members[0], csize)
	require.Nil(t, OverQuorumVotes([]Msg{vote}, quorum))
	vote2 := NewPrevote(r, h, block.Hash(), defaultSigner, &testCommittee.Members[1], csize)
	result := OverQuorumVotes([]Msg{vote, vote2}, quorum)
	require.NotNil(t, result)
	require.Equal(t, 2, len(result))
	require.Equal(t, vote.Hash(), result[0].Hash())
	require.Equal(t, vote2.Hash(), result[1].Hash())

	vote.Signers().AddMember(&testCommittee.Members[1])
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
	prevote := NewPrevote(int64(15), uint64(123345), hash, defaultSigner, testCommitteeMember, 1)

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
		// without this re-initialization the payload gets discarded after the first iteration, making the decoding fail
		b.StopTimer()
		p2pPrevote.Payload = bytes.NewReader(payload)
		b.StartTimer()
	}
}

func generateTestCommittee(t *testing.B, size int) (*types.Committee, []func(hash common.Hash) blst.Signature) {
	members := make([]types.CommitteeMember, size)
	signers := make([]func(hash common.Hash) blst.Signature, size)
	for i := 0; i < size; i++ {
		key, err := blst.RandKey()
		require.NoError(t, err)
		members[i] = *makeCommitteeMemberWithKey(key, uint64(i))
		signers[i] = makeSigner(key)
	}
	committee := &types.Committee{Members: members}
	err := committee.Enrich()
	require.NoError(t, err)
	return committee, signers
}

func generatePrevotes(t *testing.B, numVotes int, committee *types.Committee, value common.Hash, signers []func(hash common.Hash) blst.Signature) []Vote {
	require.LessOrEqual(t, numVotes, committee.Len())

	votes := make([]Vote, numVotes)
	csize := committee.Len()
	h := uint64(1)
	r := int64(0)

	p := rand.New(rand.NewSource(time.Now().UnixNano())).Perm(committee.Len())

	for i := 0; i < numVotes; i++ {
		member := &committee.Members[p[i]]
		votes[i] = NewPrevote(r, h, value, signers[i], member, csize)
	}

	return votes
}

func BenchmarkSignersEncodingSize(b *testing.B) {
	committeeSizes := []int{100, 250, 500, 1000}

	for _, cSize := range committeeSizes {
		b.Run(fmt.Sprintf("CommitteeSize-%d", cSize), func(b *testing.B) {
			committee, _ := generateTestCommittee(b, cSize)
			numSigners := (cSize) / 10 // Simulate a 1/10 committee size

			// Create a Signers object and populate it
			signers := types.NewVoteSigners(cSize)
			p := rand.New(rand.NewSource(time.Now().UnixNano())).Perm(cSize)
			for i := 0; i < numSigners; i++ {
				member := &committee.Members[p[i]]
				signers.AddMember(member)
			}

			// We are not measuring time, but the size of the output.
			// The loop will run b.N times, but we only need to calculate the size once.
			// b.ReportMetric expects the value per operation.
			payload, err := rlp.EncodeToBytes(signers)
			require.NoError(b, err)

			b.ReportAllocs()
			b.ResetTimer()

			// Report the size of the RLP-encoded payload in bytes.
			b.ReportMetric(float64(len(payload)), "payload-bytes/op")
			// Report the number of signers as a non-benchmark metric for context.
			b.ReportMetric(float64(numSigners), "signers/op")
		})
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

func aggregatePrevotes(votes []Vote) *Prevote {
	aggregates := AggregatePrevotes(votes)
	if len(aggregates) > 1 {
		panic("complex aggregate more than 1")
	}
	return aggregates[0]
}
