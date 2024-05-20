package message

import (
	"bytes"

	"crypto/rand"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/p2p"
	"github.com/autonity/autonity/rlp"
)

// locally created messages are considered as verified, we decode it to simulate a msgs arriving from the wire
func newUnverifiedPrevote(r int64, h uint64, value common.Hash, signer Signer, self *types.CommitteeMember, csize int) *Prevote {
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

/* //TODO(lorenzo) restore
func TestMessageDecode(t *testing.T) {
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
		require.Equal(t, vote.Signers().Bits, decoded.Signers().Bits)
		require.Equal(t, vote.Signers().Coefficients, decoded.Signers().Coefficients)
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
		require.Equal(t, vote.Signers().Bits, decoded.Signers().Bits)
		require.Equal(t, vote.Signers().Coefficients, decoded.Signers().Coefficients)
		require.Equal(t, vote.Signature(), decoded.Signature())
	})
	t.Run("propose", func(t *testing.T) {
		header := &types.Header{Number: common.Big2}
		block := types.NewBlockWithHeader(header)
		proposal := NewPropose(1, 2, -1, block, defaultSigner, testCommitteeMember)
		decoded := &Propose{}
		reader := bytes.NewReader(proposal.Payload())
		if err := rlp.Decode(reader, decoded); err != nil {
			t.Fatalf("have %v, want nil", err)
		}
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
	t.Run("invalid signature, error returned", func(t *testing.T) {
		header := &types.Header{Number: new(big.Int).SetUint64(25)}
		msg := newUnverifiedPrevote(1, 25, header.Hash(), func(hash common.Hash) blst.Signature {
			// tamper the hash to make signature invalid
			hash[0] = 0xca
			hash[1] = 0xfe
			return defaultSigner(hash)
		}, testCommitteeMember, 1)
		err := msg.PreValidate(testHeader)
		require.NoError(t, err)
		err = msg.Validate()
		require.ErrorIs(t, err, ErrBadSignature)
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
		header := &types.Header{Number: new(big.Int).SetUint64(25), Committee: types.Committee{otherCommitteeMember}}
		messages := []Msg{
			newUnverifiedPropose(1, 25, -1, types.NewBlockWithHeader(header), defaultSigner, testCommitteeMember),
			newUnverifiedLightPropose(1, 25, -1, types.NewBlockWithHeader(header), defaultSigner, testCommitteeMember),
		}

		for _, message := range messages {
			err := message.PreValidate(header)
			require.ErrorIs(t, err, ErrUnauthorizedAddress)
		}
	})
	t.Run("proposals from a committee member, no error", func(t *testing.T) {
		header := &types.Header{Number: new(big.Int).SetUint64(25), Committee: []types.CommitteeMember{*testCommitteeMember}}
		messages := []Msg{
			newUnverifiedPropose(1, 25, -1, types.NewBlockWithHeader(header), defaultSigner, testCommitteeMember),
			newUnverifiedLightPropose(1, 25, -1, types.NewBlockWithHeader(header), defaultSigner, testCommitteeMember),
		}

		for _, message := range messages {
			err := message.PreValidate(header)
			require.NoError(t, err)
		}
	})
	t.Run("votes with correct signers information, no error", func(t *testing.T) {
		header := &types.Header{Number: new(big.Int).SetUint64(25), Committee: []types.CommitteeMember{*testCommitteeMember}}
		messages := []Msg{
			newUnverifiedPrevote(1, 25, header.Hash(), defaultSigner, testCommitteeMember, 1),
			newUnverifiedPrecommit(1, 25, header.Hash(), defaultSigner, testCommitteeMember, 1),
		}

		for _, message := range messages {
			err := message.PreValidate(header)
			require.NoError(t, err)
		}
	})
	t.Run("votes with incorrect signers information, error is returned", func(t *testing.T) {
		header := &types.Header{Number: new(big.Int).SetUint64(25), Committee: []types.CommitteeMember{*testCommitteeMember}}
		messages := []Vote{
			newUnverifiedPrevote(1, 25, header.Hash(), defaultSigner, testCommitteeMember, 1),
			newUnverifiedPrecommit(1, 25, header.Hash(), defaultSigner, testCommitteeMember, 1),
		}

		// tamper with signers information
		messages[0].Signers().Bits = make([]byte, 100)
		messages[1].Signers().Bits = make([]byte, 0)

		for _, message := range messages {
			err := message.PreValidate(header)
			require.Error(t, err)
		}
	})
	t.Run("votes is complex aggregate but does not carry quorum, error is returned", func(t *testing.T) {
		header := &types.Header{Number: new(big.Int).SetUint64(25), Committee: []types.CommitteeMember{*testCommitteeMember, *testCommitteeMember, *testCommitteeMember, *testCommitteeMember, *testCommitteeMember}}
		vote := newUnverifiedPrevote(1, 25, header.Hash(), defaultSigner, testCommitteeMember, 5)

		// let's make this vote complex by tweaking the signers (NOTE: this will not pass validate since the signature doesn't actually match the signers)
		vote.Signers().Bits.Set(0, 2)
		vote.Signers().Bits.Set(1, 1)

		err := vote.PreValidate(header)
		require.True(t, errors.Is(err, ErrInvalidComplexAggregate))
	})

}

/* //TODO(lorenzo) restore
func TestMessageEncodeDecode(t *testing.T) {
	header := &types.Header{Number: new(big.Int).SetUint64(2), Committee: []types.CommitteeMember{*testCommitteeMember}}
	messages := []Msg{
		NewPropose(1, 2, -1, types.NewBlockWithHeader(header), defaultSigner, testCommitteeMember),
		NewPrevote(1, 2, header.Hash(), defaultSigner, testCommitteeMember, 1),
		NewPrecommit(1, 2, header.Hash(), defaultSigner, testCommitteeMember, 1),
	}
	for i := range messages {
		buff := new(bytes.Buffer)
		err := rlp.Encode(buff, messages[i])
		require.NoError(t, err)
		decoded := reflect.New(reflect.TypeOf(messages[i]).Elem()).Interface().(Msg)
		err = rlp.Decode(buff, decoded)
		require.NoError(t, err)
		err = decoded.PreValidate(header)
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

func FuzzFromPayload(f *testing.F) {
	msg := NewPrevote(1, 2, common.Hash{}, defaultSigner, testCommitteeMember, 1)
	f.Add(msg.Payload())
	f.Fuzz(func(t *testing.T, seed []byte) {
		var p Prevote
		rlp.Decode(bytes.NewReader(seed), &p)
	})
}

func TestAggregateVotes(t *testing.T) {
	// Rules:
	// 1. votes are ordered by decreasing number of distinct signers
	// 2. a vote is aggregated to the previous ones only if it adds information (i.e. adds a new signer) and the resulting aggregate respects the boundary (max coefficient = committee size)
	h := uint64(1)
	r := int64(0)
	v := common.HexToHash("0a5843ac1c1247324a23a23f23f742f89f431293123020912dade33149f4fffe")
	csize := len(testCommittee)
	err := testCommittee.Enrich()
	require.NoError(t, err)

	var votes []Vote

	//NOTE: we can use whatever signer, the aggregation functions do not verify the signature

	votes = append(votes, NewPrevote(r, h, v, defaultSigner, &testCommittee[0], csize))
	aggregate := AggregatePrevotes(votes)
	t.Log(aggregate.Signers().String())
	require.Equal(t, fmt.Sprintf("%08b", aggregate.Signers().Bits[0]), "01000000")
	require.NoError(t, aggregate.Signers().Validate(csize))

	votes = append(votes, NewPrevote(r, h, v, defaultSigner, &testCommittee[0], csize))
	aggregate = AggregatePrevotes(votes)
	t.Log(aggregate.Signers().String())
	require.Equal(t, fmt.Sprintf("%08b", aggregate.Signers().Bits[0]), "01000000")
	require.NoError(t, aggregate.Signers().Validate(csize))

	votes = append(votes, NewPrevote(r, h, v, defaultSigner, &testCommittee[1], csize))
	aggregate = AggregatePrevotes(votes)
	t.Log(aggregate.Signers().String())
	require.Equal(t, fmt.Sprintf("%08b", aggregate.Signers().Bits[0]), "01010000")
	require.NoError(t, aggregate.Signers().Validate(csize))

	votes = append(votes, NewPrevote(r, h, v, defaultSigner, &testCommittee[2], csize))
	aggregate = AggregatePrevotes(votes)
	t.Log(aggregate.Signers().String())
	require.Equal(t, fmt.Sprintf("%08b", aggregate.Signers().Bits[0]), "01010100")
	require.NoError(t, aggregate.Signers().Validate(csize))

	votes = append(votes, NewPrevote(r, h, v, defaultSigner, &testCommittee[3], csize))
	aggregate = AggregatePrevotes(votes)
	t.Log(aggregate.Signers().String())
	require.Equal(t, fmt.Sprintf("%08b", aggregate.Signers().Bits[0]), "01010101")
	require.NoError(t, aggregate.Signers().Validate(csize))

	aggregate2 := AggregatePrevotes([]Vote{NewPrevote(r, h, v, defaultSigner, &testCommittee[0], csize), NewPrevote(r, h, v, defaultSigner, &testCommittee[1], csize), NewPrevote(r, h, v, defaultSigner, &testCommittee[2], csize)})
	aggregate3 := AggregatePrevotes([]Vote{aggregate, aggregate2})
	t.Log(aggregate3.Signers().String())
	require.Equal(t, fmt.Sprintf("%08b", aggregate3.Signers().Bits[0]), "01010101")
	require.NoError(t, aggregate3.Signers().Validate(csize))

	aggregate4 := AggregatePrevotes([]Vote{NewPrevote(r, h, v, defaultSigner, &testCommittee[0], csize), NewPrevote(r, h, v, defaultSigner, &testCommittee[1], csize)})
	aggregate5 := AggregatePrevotes([]Vote{NewPrevote(r, h, v, defaultSigner, &testCommittee[0], csize), NewPrevote(r, h, v, defaultSigner, &testCommittee[2], csize)})
	aggregate6 := AggregatePrevotes([]Vote{aggregate4, aggregate5})
	t.Log(aggregate6.Signers().String())
	require.Equal(t, fmt.Sprintf("%08b", aggregate6.Signers().Bits[0]), "10010100")
	require.NoError(t, aggregate6.Signers().Validate(csize))

	aggregate7 := AggregatePrevotes([]Vote{aggregate6, aggregate3})
	t.Log(aggregate7.Signers().String())
	require.Equal(t, fmt.Sprintf("%08b", aggregate7.Signers().Bits[0]), "01010101")
	require.NoError(t, aggregate7.Signers().Validate(csize))

	// TODO(lorenzo) add test to verify that we never exceed committeeSize boundary, even when merging two complex aggregates
}

func TestAggregateVotesSimple(t *testing.T) {
	// Rules:
	// 1. votes are ordered by decreasing number of distinct signers
	// 2. a vote is aggregated to the previous ones only if it adds information (i.e. adds a new signer)
	// 3. a vote is aggregated to the previous one only if it does not produce a complex signature
	h := uint64(1)
	r := int64(0)
	v := common.HexToHash("0a5843ac1c1247324a23a23f23f742f89f431293123020912dade33149f4fffe")
	csize := len(testCommittee)
	err := testCommittee.Enrich()
	require.NoError(t, err)

	var votes []Vote

	votes = append(votes, NewPrevote(r, h, v, defaultSigner, &testCommittee[0], csize))
	aggregates := AggregatePrevotesSimple(votes)
	for _, aggregate := range aggregates {
		t.Log(aggregate.Signers().String())
	}
	require.Equal(t, 1, len(aggregates))
	require.Equal(t, fmt.Sprintf("%08b", aggregates[0].Signers().Bits[0]), "01000000")
	require.NoError(t, aggregates[0].Signers().Validate(csize))

	votes = append(votes, NewPrevote(r, h, v, defaultSigner, &testCommittee[0], csize))
	aggregates = AggregatePrevotesSimple(votes)
	for _, aggregate := range aggregates {
		t.Log(aggregate.Signers().String())
	}
	require.Equal(t, 1, len(aggregates))
	require.Equal(t, fmt.Sprintf("%08b", aggregates[0].Signers().Bits[0]), "01000000")
	require.NoError(t, aggregates[0].Signers().Validate(csize))

	votes = append(votes, NewPrevote(r, h, v, defaultSigner, &testCommittee[1], csize))
	aggregates = AggregatePrevotesSimple(votes)
	for _, aggregate := range aggregates {
		t.Log(aggregate.Signers().String())
	}
	require.Equal(t, 1, len(aggregates))
	require.Equal(t, fmt.Sprintf("%08b", aggregates[0].Signers().Bits[0]), "01010000")
	require.NoError(t, aggregates[0].Signers().Validate(csize))

	votes = append(votes, NewPrevote(r, h, v, defaultSigner, &testCommittee[2], csize))
	aggregates = AggregatePrevotesSimple(votes)
	for _, aggregate := range aggregates {
		t.Log(aggregate.Signers().String())
	}
	require.Equal(t, 1, len(aggregates))
	require.Equal(t, fmt.Sprintf("%08b", aggregates[0].Signers().Bits[0]), "01010100")
	require.NoError(t, aggregates[0].Signers().Validate(csize))

	votes = append(votes, NewPrevote(r, h, v, defaultSigner, &testCommittee[3], csize))
	aggregates = AggregatePrevotesSimple(votes)
	for _, aggregate := range aggregates {
		t.Log(aggregate.Signers().String())
	}
	require.Equal(t, 1, len(aggregates))
	require.Equal(t, fmt.Sprintf("%08b", aggregates[0].Signers().Bits[0]), "01010101")
	require.NoError(t, aggregates[0].Signers().Validate(csize))

	// aggregate overlaps, should not get merged

	aggregates2 := AggregatePrevotesSimple([]Vote{NewPrevote(r, h, v, defaultSigner, &testCommittee[0], csize), NewPrevote(r, h, v, defaultSigner, &testCommittee[1], csize)})
	aggregates3 := AggregatePrevotesSimple([]Vote{NewPrevote(r, h, v, defaultSigner, &testCommittee[0], csize), NewPrevote(r, h, v, defaultSigner, &testCommittee[2], csize)})
	aggregates4 := AggregatePrevotesSimple([]Vote{aggregates2[0], aggregates3[0]})
	for _, aggregate := range aggregates4 {
		t.Log(aggregate.Signers().String())
	}
	require.Equal(t, 2, len(aggregates4))
	require.Equal(t, fmt.Sprintf("%08b", aggregates4[0].Signers().Bits[0]), "01010000")
	require.Equal(t, fmt.Sprintf("%08b", aggregates4[1].Signers().Bits[0]), "01000100")
	require.NoError(t, aggregates4[0].Signers().Validate(csize))
	require.NoError(t, aggregates4[1].Signers().Validate(csize))

	vote := NewPrevote(r, h, v, defaultSigner, &testCommittee[3], csize)
	aggregates5 := AggregatePrevotesSimple([]Vote{aggregates2[0], aggregates3[0], vote})
	for _, aggregate := range aggregates5 {
		t.Log(aggregate.Signers().String())
	}
	require.Equal(t, 2, len(aggregates5))
	require.Equal(t, fmt.Sprintf("%08b", aggregates5[0].Signers().Bits[0]), "01010001")
	require.Equal(t, fmt.Sprintf("%08b", aggregates5[1].Signers().Bits[0]), "01000100")
	require.NoError(t, aggregates5[0].Signers().Validate(csize))
	require.NoError(t, aggregates5[1].Signers().Validate(csize))

	//TODO(Lorenzo) absolutely needs more test cases here
	// - add test where we feed a complex aggregate + other messages to aggregatesimple
	// - ??
}

//TODO(lorenzo) add tests for message.Power()

/* TODO: to fix in FD PR
func TestOverQuorumVotes(t *testing.T) {
  t.Run("with duplicated votes, it returns none duplicated votes of just quorum ones", func(t *testing.T) {
    seats := 10
    committee, _ := GenerateCommittee(seats)
    quorum := bft.Quorum(big.NewInt(int64(seats)))
    height := uint64(1)
    round := uint64(0)
    notNilValue := common.Hash{0x1}
    var preVotes []message.Msg
    for i := 0; i < len(committee); i++ {
      preVote := message.NewFakePrevote(message.Fake{
        FakeSender: committee[i].Address,
        FakeRound:  round,
        FakeHeight: height,
        FakeValue:  notNilValue,
        FakePower:  common.Big1,
      })
      preVotes = append(preVotes, preVote)
    }

    // let duplicated msg happens, the counting should skip duplicated ones.
    preVotes = append(preVotes, preVotes...)

    overQuorumVotes := OverQuorumVotes(preVotes, quorum)
    require.Equal(t, quorum.Uint64(), uint64(len(overQuorumVotes)))
  })

  t.Run("with less quorum votes, it returns no votes", func(t *testing.T) {
    seats := 10
    committee, _ := GenerateCommittee(seats)
    quorum := bft.Quorum(new(big.Int).SetInt64(int64(seats)))
    height := uint64(1)
    round := uint64(0)
    noneNilValue := common.Hash{0x1}
    var preVotes []message.Msg
    for i := 0; i < int(quorum.Uint64()-1); i++ {
      preVote := message.NewFakePrevote(message.Fake{
        FakeRound:  round,
        FakeHeight: height,
        FakeValue:  noneNilValue,
        FakeSender: committee[i].Address,
        FakePower:  common.Big1,
      })
      preVotes = append(preVotes, preVote)
    }

    // let duplicated msg happens, the counting should skip duplicated ones.
    preVotes = append(preVotes, preVotes...)

    overQuorumVotes := OverQuorumVotes(preVotes, quorum)
    require.Nil(t, overQuorumVotes)
  })
}
*/

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

func TestAggregateVotesSimple(t *testing.T) {
	// Rules:
	// 1. votes are ordered by decreasing number of distinct signers
	// 2. a vote is aggregated to the previous ones only if it adds information (i.e. adds a new signer)
	// 3. a vote is aggregated to the previous one only if it does not produce a complex signature
	h := uint64(1)
	r := int64(0)
	v := common.HexToHash("0a5843ac1c1247324a23a23f23f742f89f431293123020912dade33149f4fffe")
	csize := len(header.Committee)
	err := header.Committee.Enrich()
	require.NoError(t, err)

	//TODO(lorenzo) refactor to reduce code duplication and add more cases
	// - add test where we feed a complex aggregate + other messages to aggregatesimple

	var votes []Vote

	votes = append(votes, NewPrevote(r, h, v, makeSigner(0), header))
	aggregates := AggregatePrevotesSimple(votes)
	for _, aggregate := range aggregates {
		fmt.Println(aggregate.Signers().String())
	}
	require.Equal(t, 1, len(aggregates))
	require.Equal(t, fmt.Sprintf("%08b", aggregates[0].Signers().Bits[0]), "01000000")
	require.True(t, aggregates[0].Signers().Validate(csize))

	votes = append(votes, NewPrevote(r, h, v, makeSigner(0), header))
	aggregates = AggregatePrevotesSimple(votes)
	for _, aggregate := range aggregates {
		fmt.Println(aggregate.Signers().String())
	}
	require.Equal(t, 1, len(aggregates))
	require.Equal(t, fmt.Sprintf("%08b", aggregates[0].Signers().Bits[0]), "01000000")
	require.True(t, aggregates[0].Signers().Validate(csize))

	votes = append(votes, NewPrevote(r, h, v, makeSigner(1), header))
	aggregates = AggregatePrevotesSimple(votes)
	for _, aggregate := range aggregates {
		fmt.Println(aggregate.Signers().String())
	}
	require.Equal(t, 1, len(aggregates))
	require.Equal(t, fmt.Sprintf("%08b", aggregates[0].Signers().Bits[0]), "01010000")
	require.True(t, aggregates[0].Signers().Validate(csize))

	votes = append(votes, NewPrevote(r, h, v, makeSigner(2), header))
	aggregates = AggregatePrevotesSimple(votes)
	for _, aggregate := range aggregates {
		fmt.Println(aggregate.Signers().String())
	}
	require.Equal(t, 1, len(aggregates))
	require.Equal(t, fmt.Sprintf("%08b", aggregates[0].Signers().Bits[0]), "01010100")
	require.True(t, aggregates[0].Signers().Validate(csize))

	votes = append(votes, NewPrevote(r, h, v, makeSigner(3), header))
	aggregates = AggregatePrevotesSimple(votes)
	for _, aggregate := range aggregates {
		fmt.Println(aggregate.Signers().String())
	}
	require.Equal(t, 1, len(aggregates))
	require.Equal(t, fmt.Sprintf("%08b", aggregates[0].Signers().Bits[0]), "01010101")
	require.True(t, aggregates[0].Signers().Validate(csize))

	// aggregate overlaps, should not get merged

	aggregates2 := AggregatePrevotesSimple([]Vote{NewPrevote(r, h, v, makeSigner(0), header), NewPrevote(r, h, v, makeSigner(1), header)})
	aggregates3 := AggregatePrevotesSimple([]Vote{NewPrevote(r, h, v, makeSigner(0), header), NewPrevote(r, h, v, makeSigner(2), header)})
	aggregates4 := AggregatePrevotesSimple([]Vote{aggregates2[0], aggregates3[0]})
	for _, aggregate := range aggregates4 {
		fmt.Println(aggregate.Signers().String())
	}
	require.Equal(t, 2, len(aggregates4))
	require.Equal(t, fmt.Sprintf("%08b", aggregates4[0].Signers().Bits[0]), "01010000")
	require.Equal(t, fmt.Sprintf("%08b", aggregates4[1].Signers().Bits[0]), "01000100")
	require.True(t, aggregates4[0].Signers().Validate(csize))
	require.True(t, aggregates4[1].Signers().Validate(csize))

	vote := NewPrevote(r, h, v, makeSigner(3), header)
	aggregates5 := AggregatePrevotesSimple([]Vote{aggregates2[0], aggregates3[0], vote})
	for _, aggregate := range aggregates5 {
		fmt.Println(aggregate.Signers().String())
	}
	require.Equal(t, 2, len(aggregates5))
	require.Equal(t, fmt.Sprintf("%08b", aggregates5[0].Signers().Bits[0]), "01010001")
	require.Equal(t, fmt.Sprintf("%08b", aggregates5[1].Signers().Bits[0]), "01000100")
	require.True(t, aggregates5[0].Signers().Validate(csize))
	require.True(t, aggregates5[1].Signers().Validate(csize))

	//TODO(Lorenzo) absolutely needs more test cases here
}

/*
// TODO(lorenzo) refinements, fix this benchmark, I think the p2p MSg payload gets discarded after a couple iteration of decoding
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
}*/
