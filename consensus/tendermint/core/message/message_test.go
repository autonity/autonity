package message

import (
	"bytes"
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/common/hexutil"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/crypto/blst"
)

/* //TODO(lorenzo) restore
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
*/

var header = &types.Header{
	ParentHash:  common.HexToHash("0a5843ac1cb04865017cb35a57b50b07084e5fcee39b5acadade33149f4fff9e"),
	UncleHash:   common.HexToHash("0a5843ac1c124732472342342387423897431293123020912dade33149f4fffe"),
	Coinbase:    common.HexToAddress("8888f1f195afa192cfee860698584c030f4c9db1"),
	Root:        common.HexToHash("0a5843ac1cb0486345235234564778768967856745645654649f4fff9321321e"),
	TxHash:      common.HexToHash("0a58213121cb0486345235234564778768967856745645654649f4fff932132e"),
	ReceiptHash: common.HexToHash("9a58213121cb0486345235234564778768967856745645654649f4fff932132e"),
	Bloom:       types.BytesToBloom(bytes.Repeat([]byte("a"), 128)),
	Difficulty:  big.NewInt(199),
	Number:      big.NewInt(239),
	GasLimit:    uint64(1000),
	GasUsed:     uint64(400),
	Time:        uint64(12343),
	MixDigest:   common.HexToHash("0a58213121cb0486345235234564778768967853123645654649f4fff932132e"),
	Nonce:       [8]byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07},
	BaseFee:     big.NewInt(20000),
	Committee: []types.CommitteeMember{{
		Address:           common.HexToAddress("0x76a685e4bf8cbcd25d7d3b6c342f64a30b503380"),
		ConsensusKeyBytes: hexutil.MustDecode("0x951f3f7ab473eb0d00eaaa569ba1a0be2877b794e29e0cbf504b7f00cb879a824b0b913397e0071a87cebaae2740002b"),
		VotingPower:       hexutil.MustDecodeBig("0x3029"),
	}, {
		Address:           common.HexToAddress("0xc44276975a6c2d12e62e18d814b507c38fc3646f"),
		ConsensusKeyBytes: hexutil.MustDecode("0x8bddc21fca7f3a920064729547605c73e55c17e20917eddc8788b97990c0d7e9420e51a97ea400fb58a5c28fa63984eb"),
		VotingPower:       hexutil.MustDecodeBig("0x3139"),
	}, {
		Address:           common.HexToAddress("0x1a72cb9d17c9e7acad03b4d3505f160e3782f2d5"),
		ConsensusKeyBytes: hexutil.MustDecode("0x9679c8ebd47d18b93acd90cd380debdcfdb140f38eca207c61463a47be85398ec3082a66f7f30635c11470f5c8e5cf6b"),
		VotingPower:       hexutil.MustDecodeBig("0x3056"),
	}, {
		Address:           common.HexToAddress("0xbaa58a01e5ca81dc288e2c46a8a467776bdb81c6"),
		ConsensusKeyBytes: hexutil.MustDecode("0xa460c204c407b6272f7731b0d15daca8f2564cf7ace301769e3b42de2482fc3bf8116dd13c0545e806441d074d02dcc2"),
		VotingPower:       hexutil.MustDecodeBig("0x39"),
	}},
	ProposerSeal:   bytes.Repeat([]byte("c"), 65),
	Round:          uint64(3),
	CommittedSeals: new(types.AggregateSignature),
}

// creates a dummy signer. We don't care that the signature doesn't actually match the consensus key of the validator
func makeSigner(i int) func(hash common.Hash) (signature blst.Signature, address common.Address) {
	return func(hash common.Hash) (signature blst.Signature, address common.Address) {
		sk, _ := blst.RandKey()
		return sk.Sign(hash[:]), header.Committee[i].Address
	}
}

func TestAggregateVotes(t *testing.T) {
	// Rules:
	// 1. votes are ordered by decreasing number of distinct senders
	// 2. a vote is aggregated to the previous ones only if it adds information (i.e. adds a new sender)
	h := uint64(1)
	r := int64(0)
	v := common.HexToHash("0a5843ac1c1247324a23a23f23f742f89f431293123020912dade33149f4fffe")
	csize := len(header.Committee)
	err := header.Committee.Enrich()
	require.NoError(t, err)

	//TODO(lorenzo) refactor to reduce code duplication and add more cases

	var votes []Vote

	votes = append(votes, NewPrevote(r, h, v, makeSigner(0), header))
	aggregate := AggregatePrevotes(votes)
	fmt.Println(aggregate.Senders().String())
	require.Equal(t, fmt.Sprintf("%08b", aggregate.Senders().Bits[0]), "01000000")
	require.True(t, aggregate.Senders().Valid(csize))

	votes = append(votes, NewPrevote(r, h, v, makeSigner(0), header))
	aggregate = AggregatePrevotes(votes)
	fmt.Println(aggregate.Senders().String())
	require.Equal(t, fmt.Sprintf("%08b", aggregate.Senders().Bits[0]), "01000000")
	require.True(t, aggregate.Senders().Valid(csize))

	votes = append(votes, NewPrevote(r, h, v, makeSigner(1), header))
	aggregate = AggregatePrevotes(votes)
	fmt.Println(aggregate.Senders().String())
	require.Equal(t, fmt.Sprintf("%08b", aggregate.Senders().Bits[0]), "01010000")
	require.True(t, aggregate.Senders().Valid(csize))

	votes = append(votes, NewPrevote(r, h, v, makeSigner(2), header))
	aggregate = AggregatePrevotes(votes)
	fmt.Println(aggregate.Senders().String())
	require.Equal(t, fmt.Sprintf("%08b", aggregate.Senders().Bits[0]), "01010100")
	require.True(t, aggregate.Senders().Valid(csize))

	votes = append(votes, NewPrevote(r, h, v, makeSigner(3), header))
	aggregate = AggregatePrevotes(votes)
	fmt.Println(aggregate.Senders().String())
	require.Equal(t, fmt.Sprintf("%08b", aggregate.Senders().Bits[0]), "01010101")
	require.True(t, aggregate.Senders().Valid(csize))

	aggregate2 := AggregatePrevotes([]Vote{NewPrevote(r, h, v, makeSigner(0), header), NewPrevote(r, h, v, makeSigner(1), header), NewPrevote(r, h, v, makeSigner(2), header)})
	aggregate3 := AggregatePrevotes([]Vote{aggregate, aggregate2})
	fmt.Println(aggregate3.Senders().String())
	require.Equal(t, fmt.Sprintf("%08b", aggregate3.Senders().Bits[0]), "01010101")
	require.True(t, aggregate3.Senders().Valid(csize))

	aggregate4 := AggregatePrevotes([]Vote{NewPrevote(r, h, v, makeSigner(0), header), NewPrevote(r, h, v, makeSigner(1), header)})
	aggregate5 := AggregatePrevotes([]Vote{NewPrevote(r, h, v, makeSigner(0), header), NewPrevote(r, h, v, makeSigner(2), header)})
	aggregate6 := AggregatePrevotes([]Vote{aggregate4, aggregate5})
	fmt.Println(aggregate6.Senders().String())
	require.Equal(t, fmt.Sprintf("%08b", aggregate6.Senders().Bits[0]), "10010100")
	require.True(t, aggregate6.Senders().Valid(csize))

	aggregate7 := AggregatePrevotes([]Vote{aggregate6, aggregate3})
	fmt.Println(aggregate7.Senders().String())
	require.Equal(t, fmt.Sprintf("%08b", aggregate7.Senders().Bits[0]), "01010101")
	require.True(t, aggregate7.Senders().Valid(csize))

}

func TestAggregateVotesSimple(t *testing.T) {
	// Rules:
	// 1. votes are ordered by decreasing number of distinct senders
	// 2. a vote is aggregated to the previous ones only if it adds information (i.e. adds a new sender)
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
		fmt.Println(aggregate.Senders().String())
	}
	require.Equal(t, 1, len(aggregates))
	require.Equal(t, fmt.Sprintf("%08b", aggregates[0].Senders().Bits[0]), "01000000")
	require.True(t, aggregates[0].Senders().Valid(csize))

	votes = append(votes, NewPrevote(r, h, v, makeSigner(0), header))
	aggregates = AggregatePrevotesSimple(votes)
	for _, aggregate := range aggregates {
		fmt.Println(aggregate.Senders().String())
	}
	require.Equal(t, 1, len(aggregates))
	require.Equal(t, fmt.Sprintf("%08b", aggregates[0].Senders().Bits[0]), "01000000")
	require.True(t, aggregates[0].Senders().Valid(csize))

	votes = append(votes, NewPrevote(r, h, v, makeSigner(1), header))
	aggregates = AggregatePrevotesSimple(votes)
	for _, aggregate := range aggregates {
		fmt.Println(aggregate.Senders().String())
	}
	require.Equal(t, 1, len(aggregates))
	require.Equal(t, fmt.Sprintf("%08b", aggregates[0].Senders().Bits[0]), "01010000")
	require.True(t, aggregates[0].Senders().Valid(csize))

	votes = append(votes, NewPrevote(r, h, v, makeSigner(2), header))
	aggregates = AggregatePrevotesSimple(votes)
	for _, aggregate := range aggregates {
		fmt.Println(aggregate.Senders().String())
	}
	require.Equal(t, 1, len(aggregates))
	require.Equal(t, fmt.Sprintf("%08b", aggregates[0].Senders().Bits[0]), "01010100")
	require.True(t, aggregates[0].Senders().Valid(csize))

	votes = append(votes, NewPrevote(r, h, v, makeSigner(3), header))
	aggregates = AggregatePrevotesSimple(votes)
	for _, aggregate := range aggregates {
		fmt.Println(aggregate.Senders().String())
	}
	require.Equal(t, 1, len(aggregates))
	require.Equal(t, fmt.Sprintf("%08b", aggregates[0].Senders().Bits[0]), "01010101")
	require.True(t, aggregates[0].Senders().Valid(csize))

	// aggregate overlaps, should not get merged

	aggregates2 := AggregatePrevotesSimple([]Vote{NewPrevote(r, h, v, makeSigner(0), header), NewPrevote(r, h, v, makeSigner(1), header)})
	aggregates3 := AggregatePrevotesSimple([]Vote{NewPrevote(r, h, v, makeSigner(0), header), NewPrevote(r, h, v, makeSigner(2), header)})
	aggregates4 := AggregatePrevotesSimple([]Vote{aggregates2[0], aggregates3[0]})
	for _, aggregate := range aggregates4 {
		fmt.Println(aggregate.Senders().String())
	}
	require.Equal(t, 2, len(aggregates4))
	require.Equal(t, fmt.Sprintf("%08b", aggregates4[0].Senders().Bits[0]), "01010000")
	require.Equal(t, fmt.Sprintf("%08b", aggregates4[1].Senders().Bits[0]), "01000100")
	require.True(t, aggregates4[0].Senders().Valid(csize))
	require.True(t, aggregates4[1].Senders().Valid(csize))

	vote := NewPrevote(r, h, v, makeSigner(3), header)
	aggregates5 := AggregatePrevotesSimple([]Vote{aggregates2[0], aggregates3[0], vote})
	for _, aggregate := range aggregates5 {
		fmt.Println(aggregate.Senders().String())
	}
	require.Equal(t, 2, len(aggregates5))
	require.Equal(t, fmt.Sprintf("%08b", aggregates5[0].Senders().Bits[0]), "01010001")
	require.Equal(t, fmt.Sprintf("%08b", aggregates5[1].Senders().Bits[0]), "01000100")
	require.True(t, aggregates5[0].Senders().Valid(csize))
	require.True(t, aggregates5[1].Senders().Valid(csize))

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
