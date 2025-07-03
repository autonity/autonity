package accountability

import (
	cr "crypto/rand"
	"math/big"
	"math/rand"
	"testing"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/core/constants"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/crypto"
	"github.com/autonity/autonity/crypto/blst"
	"github.com/autonity/autonity/rlp"
	"github.com/stretchr/testify/require"
)

var (
	h            = rand.Uint64()
	r            = rand.Int63n(constants.MaxRound / 2)
	randomBytes1 = make([]byte, 32)
	randomBytes2 = make([]byte, 32)
	_, _         = cr.Read(randomBytes1)
	_, _         = cr.Read(randomBytes2)
	values       = []common.Hash{crypto.Hash(randomBytes1), crypto.Hash(randomBytes2), common.NilValue}
	parentHeader = newBlockHeader(h-1, committee)
)

func TestRLPEncodingDecoding(t *testing.T) {
	rvs := Signers{
		Round:               r,
		Value:               common.Hash{},
		SignersIndex:        []int{0},
		SignersCoeff:        []uint16{1},
		aggregatedPublicKey: nil,
		hasSigners:          nil,
		preValidated:        false,
	}
	p, err := rlp.EncodeToBytes(&rvs)
	require.NoError(t, err)

	decodeRVS := &Signers{}
	err = rlp.DecodeBytes(p, decodeRVS)
	require.NoError(t, err)
	require.Equal(t, r, decodeRVS.Round)

	rvs.Round = constants.MaxRound + 1
	p, err = rlp.EncodeToBytes(&rvs)
	require.NoError(t, err)
	decodeRVS = &Signers{}
	err = rlp.DecodeBytes(p, decodeRVS)
	require.NotNil(t, err)

	rvs.Round = constants.MaxRound
	rvs.SignersIndex = []int{}
	rvs.SignersCoeff = []uint16{}
	p, err = rlp.EncodeToBytes(&rvs)
	require.NoError(t, err)
	decodeRVS = &Signers{}
	err = rlp.DecodeBytes(p, decodeRVS)
	require.NotNil(t, err)
}

func TestHighlyAggregatedPrecommit(t *testing.T) {
	aggregatedPrecommit := randomHighlyAggregatedPrecommits(h, r)
	payload, err := rlp.EncodeToBytes(aggregatedPrecommit)
	require.NoError(t, err)

	decodedPrecomit := &HighlyAggregatedPrecommit{}
	err = rlp.DecodeBytes(payload, decodedPrecomit)
	require.NoError(t, err)

	err = decodedPrecomit.PreValidate(parentHeader.Epoch.Committee, h)
	require.NoError(t, err)
	require.Equal(t, true, decodedPrecomit.preValidated)
	err = decodedPrecomit.Validate()
	require.NoError(t, err)
	require.Equal(t, true, decodedPrecomit.validated)
}

func TestVerifyMaliciousAggregatedPrecommits(t *testing.T) {
	numOfFastAggPrecommits := 10
	var precommits []*message.Precommit
	for n := 0; n < numOfFastAggPrecommits; n++ {
		value := values[n%len(values)]
		precommits = append(precommits, aggregatedPrecommit(height, int64(n), value, randomSigners(cSize), committee, keys))
	}

	t.Run("with wrong height", func(t *testing.T) {
		wrongHeight := height + 1
		aggPrecommits := maliciousAggregatePrecommits(precommits, &wrongHeight, nil, nil, nil)
		payload, err := rlp.EncodeToBytes(aggPrecommits)
		require.NoError(t, err)

		decodedPrecomit := &HighlyAggregatedPrecommit{}
		err = rlp.DecodeBytes(payload, decodedPrecomit)
		require.NoError(t, err)
		header := newBlockHeader(height-1, committee)
		err = decodedPrecomit.PreValidate(header.Epoch.Committee, height)
		require.Error(t, errBadHeight, err)
		require.Equal(t, false, decodedPrecomit.preValidated)
	})

	t.Run("with wrong round", func(t *testing.T) {
		wrongRound := int64(19)
		aggPrecommits := maliciousAggregatePrecommits(precommits, nil, &wrongRound, nil, nil)
		payload, err := rlp.EncodeToBytes(aggPrecommits)
		require.NoError(t, err)

		decodedPrecomit := &HighlyAggregatedPrecommit{}
		err = rlp.DecodeBytes(payload, decodedPrecomit)
		require.NoError(t, err)
		header := newBlockHeader(height-1, committee)
		err = decodedPrecomit.PreValidate(header.Epoch.Committee, height)
		require.NoError(t, err)
		require.Equal(t, true, decodedPrecomit.preValidated)
		err = decodedPrecomit.Validate()
		require.NotNil(t, err)
		require.Equal(t, false, decodedPrecomit.validated)
	})

	t.Run("with wrong value", func(t *testing.T) {
		wrongValue := common.Hash{0xff}
		aggPrecommits := maliciousAggregatePrecommits(precommits, nil, nil, &wrongValue, nil)
		payload, err := rlp.EncodeToBytes(aggPrecommits)
		require.NoError(t, err)

		decodedPrecomit := &HighlyAggregatedPrecommit{}
		err = rlp.DecodeBytes(payload, decodedPrecomit)
		require.NoError(t, err)
		header := newBlockHeader(height-1, committee)
		err = decodedPrecomit.PreValidate(header.Epoch.Committee, height)
		require.NoError(t, err)
		require.Equal(t, true, decodedPrecomit.preValidated)
		err = decodedPrecomit.Validate()
		require.NotNil(t, err)
		require.Equal(t, false, decodedPrecomit.validated)
	})

	t.Run("with wrong signers", func(t *testing.T) {
		wrongSigners := precommits[0].Signers().FlattenUniq()
		if len(wrongSigners) > 1 {
			wrongSigners[0], wrongSigners[1] = wrongSigners[1], wrongSigners[0]
		} else {
			wrongSigners[0] = (wrongSigners[0] + 1) % cSize
		}
		aggPrecommits := maliciousAggregatePrecommits(precommits, nil, nil, nil, wrongSigners)
		payload, err := rlp.EncodeToBytes(aggPrecommits)
		require.NoError(t, err)

		decodedPrecomit := &HighlyAggregatedPrecommit{}
		err = rlp.DecodeBytes(payload, decodedPrecomit)
		require.NoError(t, err)
		header := newBlockHeader(height-1, committee)
		err = decodedPrecomit.PreValidate(header.Epoch.Committee, height)
		require.NoError(t, err)
		require.Equal(t, true, decodedPrecomit.preValidated)
		err = decodedPrecomit.Validate()
		require.NotNil(t, err)
		require.Equal(t, false, decodedPrecomit.validated)
	})
}

func aggregatedPrecommit(h uint64, r int64, v common.Hash, signers []int, committee *types.Committee, keys []blst.SecretKey) *message.Precommit {
	precommits := make([]message.Vote, len(signers))
	for i, s := range signers {
		precommits[i] = newValidatedPrecommit(r, h, v, makeSigner(keys[s]), &committee.Members[s], committee.Len())
	}
	return aggregatePrecommits(precommits)
}

// randomSigners generate a set of signer's index, it could have duplicated index.
func randomSigners(committeeSize int) []int {
	size := rand.Intn(committeeSize) + 1
	result := make([]int, size)
	for i := 0; i < size; i++ {
		result[i] = rand.Intn(committeeSize)
	}
	return result
}

func randomHighlyAggregatedPrecommits(height uint64, round int64) HighlyAggregatedPrecommit {
	numOfFastAggPrecommits := 10

	var precommits []*message.Precommit
	for n := 0; n < numOfFastAggPrecommits; n++ {
		value := values[n%len(values)]
		// add duplicated msg but with different signers, thus the aggregation need to do a further fast aggregate.
		precommits = append(precommits, aggregatedPrecommit(height, round+int64(n), value, randomSigners(cSize), committee, keys))
	}
	return AggregateDistinctPrecommits(precommits)
}

// maliciousAggregatePrecommits aggregate the precommits in a wrong way by modifying meta-data.
func maliciousAggregatePrecommits(precommits []*message.Precommit, wrongHeight *uint64, wrongRound *int64,
	wrongValue *common.Hash, wrongSigners []int) HighlyAggregatedPrecommit {

	defaultHeight := precommits[0].H()
	if wrongHeight != nil {
		defaultHeight = *wrongHeight
	}

	presentedMsgs := make(map[int64]map[common.Hash]struct{})

	var precommitsToBeAggregated []*message.Precommit
	precommitsToBeAggregated = append(precommitsToBeAggregated, precommits[0])

	for i := 1; i < len(precommits); i++ {
		// skip duplicated msg.
		if _, ok := presentedMsgs[precommits[i].R()]; !ok {
			presentedMsgs[precommits[i].R()] = make(map[common.Hash]struct{})
		}
		if _, ok := presentedMsgs[precommits[i].R()][precommits[i].Value()]; !ok {
			presentedMsgs[precommits[i].R()][precommits[i].Value()] = struct{}{}
			precommitsToBeAggregated = append(precommitsToBeAggregated, precommits[i])
		}
	}

	result := HighlyAggregatedPrecommit{}
	signatures := make([]blst.Signature, len(precommitsToBeAggregated))
	for i, m := range precommitsToBeAggregated {
		defaultRound := m.R()
		defaultValue := m.Value()
		defaultSingers := m.Signers().FlattenUniq()
		coeffs := m.Signers().Coefficients
		if wrongRound != nil {
			defaultRound += *wrongRound
		}
		if wrongValue != nil {
			defaultValue = *wrongValue
		}

		if len(wrongSigners) > 0 {
			defaultSingers = wrongSigners
			for len(coeffs) < len(defaultSingers) {
				coeffs = append(coeffs, 1)
			}
			if len(coeffs) > len(defaultSingers) {
				coeffs = coeffs[:len(defaultSingers)]
			}
		}

		roundValueSigners := &Signers{
			Round:        defaultRound,
			Value:        defaultValue,
			SignersIndex: defaultSingers,
			SignersCoeff: coeffs,
		}
		result.MsgSigners = append(result.MsgSigners, roundValueSigners)
		signatures[i] = m.Signature()
	}
	result.Height = defaultHeight
	result.Signature = blst.AggregateSignatures(signatures).Marshal()
	return result
}

func TestDistinctPrecommitsWithZeroes(t *testing.T) {
	tweakedCommittee, blsKeys, _ := generateCommittee()

	// add some keys that sum to 0 to the committee
	x1Hex := "15d1a39e3e11a76ba764d153b1f5c28ace584e11a34d91f5bffb98e553e41079"
	x2Hex := "5e1c03b4eb8bd5dc8bd506b457ac157a856555f15cb0ca0940046719ac1bef88"

	x1, err := blst.SecretKeyFromHex(x1Hex)
	require.NoError(t, err)

	x2, err := blst.SecretKeyFromHex(x2Hex)
	require.NoError(t, err)

	X1 := x1.PublicKey()
	X2 := x2.PublicKey()

	aggX, err := blst.AggregatePublicKeys([]blst.PublicKey{X1, X2})
	require.NoError(t, err)

	// aggregated key is 0
	require.False(t, aggX.Validate())

	csize := tweakedCommittee.Len()
	privateKey1, _ := crypto.GenerateKey()
	committeeMember1 := types.CommitteeMember{
		Address:           crypto.PubkeyToAddress(privateKey1.PublicKey),
		VotingPower:       new(big.Int).SetUint64(1),
		ConsensusKey:      X1,
		ConsensusKeyBytes: X1.Marshal(),
		Index:             uint64(csize),
	}
	csize++
	tweakedCommittee.Members = append(tweakedCommittee.Members, committeeMember1)
	blsKeys = append(blsKeys, x1)

	privateKey2, _ := crypto.GenerateKey()
	committeeMember2 := types.CommitteeMember{
		Address:           crypto.PubkeyToAddress(privateKey2.PublicKey),
		VotingPower:       new(big.Int).SetUint64(1),
		ConsensusKey:      X2,
		ConsensusKeyBytes: X2.Marshal(),
		Index:             uint64(csize),
	}
	csize++
	tweakedCommittee.Members = append(tweakedCommittee.Members, committeeMember2)
	blsKeys = append(blsKeys, x2)

	t.Run("1 zero key/sig in the set", func(t *testing.T) {
		var precommits []*message.Precommit
		precommits = append(precommits, aggregatedPrecommit(h, r, values[0], randomSigners(cSize), tweakedCommittee, blsKeys))
		precommits = append(precommits, aggregatedPrecommit(h, r, values[1], randomSigners(cSize), tweakedCommittee, blsKeys))

		// add a precommit with zero signature
		precommits = append(precommits, aggregatedPrecommit(h, r, values[2], []int{csize - 1, csize - 2}, tweakedCommittee, blsKeys))

		// signature should be valid
		highlyAggregatedPrecommit := AggregateDistinctPrecommits(precommits)
		require.NoError(t, highlyAggregatedPrecommit.PreValidate(tweakedCommittee, h))
		require.NoError(t, highlyAggregatedPrecommit.Validate())
	})
	t.Run("whole sig/key is zero", func(t *testing.T) {
		var precommits []*message.Precommit
		precommits = append(precommits, aggregatedPrecommit(h, r, values[0], []int{csize - 1}, tweakedCommittee, blsKeys))
		precommits = append(precommits, aggregatedPrecommit(h, r, values[1], []int{csize - 2}, tweakedCommittee, blsKeys))

		// here key is 0 but signature is not 0
		highlyAggregatedPrecommit := AggregateDistinctPrecommits(precommits)
		sig, err := blst.SignatureFromBytes(highlyAggregatedPrecommit.Signature)
		require.NoError(t, err)
		require.False(t, sig.IsZero())
		require.NoError(t, highlyAggregatedPrecommit.PreValidate(tweakedCommittee, h))
		err = highlyAggregatedPrecommit.Validate()
		t.Log(err)
		require.Error(t, highlyAggregatedPrecommit.Validate())

		// substitute signature with a zero one
		m := common.Hash{0xca, 0xfe}
		sig1 := x1.Sign(m[:])
		sig2 := x2.Sign(m[:])
		zeroSig := blst.AggregateSignatures([]blst.Signature{sig1, sig2})
		require.True(t, zeroSig.IsZero())
		highlyAggregatedPrecommit.Signature = zeroSig.Marshal()

		// 0 sig with 0 key should be valid
		require.NoError(t, highlyAggregatedPrecommit.PreValidate(tweakedCommittee, h))
		require.NoError(t, highlyAggregatedPrecommit.Validate())

	})
	t.Run("all sigs are 0", func(t *testing.T) {
		var precommits []*message.Precommit

		precommits = append(precommits, aggregatedPrecommit(h, r, values[0], []int{csize - 1, csize - 2}, tweakedCommittee, blsKeys))
		precommits = append(precommits, aggregatedPrecommit(h, r, values[1], []int{csize - 1, csize - 2}, tweakedCommittee, blsKeys))

		highlyAggregatedPrecommit := AggregateDistinctPrecommits(precommits)
		require.NoError(t, highlyAggregatedPrecommit.PreValidate(tweakedCommittee, h))
		require.NoError(t, highlyAggregatedPrecommit.Validate())

		// put not 0 signature, verification should fail
		sig := x1.Sign(values[0][:])
		highlyAggregatedPrecommit.Signature = sig.Marshal()
		require.NoError(t, highlyAggregatedPrecommit.PreValidate(tweakedCommittee, h))
		err = highlyAggregatedPrecommit.Validate()
		t.Log(err)
		require.Error(t, err)
	})

}
