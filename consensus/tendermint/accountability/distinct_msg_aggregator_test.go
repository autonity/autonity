package accountability

import (
	cr "crypto/rand"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/crypto"
	"github.com/autonity/autonity/crypto/blst"
	"github.com/autonity/autonity/rlp"
	"github.com/stretchr/testify/require"
	"math/rand"
	"testing"
)

var (
	h            = rand.Uint64()
	r            = rand.Int63()
	randomBytes1 = make([]byte, 32)
	randomBytes2 = make([]byte, 32)
	_, _         = cr.Read(randomBytes1)
	_, _         = cr.Read(randomBytes2)
	values       = []common.Hash{crypto.Hash(randomBytes1), crypto.Hash(randomBytes2), nilValue}
	parentHeader = newBlockHeader(h-1, committee)
)

func TestRLPEncodingDeconding(t *testing.T) {
	rvs := RoundValueSigners{
		Round:               r,
		Value:               common.Hash{},
		Signers:             nil,
		aggregatedPublicKey: nil,
		hasSigners:          nil,
		preValidated:        false,
	}
	p, err := rlp.EncodeToBytes(&rvs)
	require.NoError(t, err)

	decodeRVS := &RoundValueSigners{}
	err = rlp.DecodeBytes(p, decodeRVS)
	require.NoError(t, err)
	require.Equal(t, r, decodeRVS.Round)
}

func TestHighlyAggregatedPrecommit(t *testing.T) {
	aggregatedPrecommit := randomHighlyAggregatedPrecommits(h, r)
	payload, err := rlp.EncodeToBytes(aggregatedPrecommit)
	require.NoError(t, err)

	decodedPrecomit := &HighlyAggregatedPrecommit{}
	err = rlp.DecodeBytes(payload, decodedPrecomit)
	require.NoError(t, err)

	err = decodedPrecomit.PreValidate(parentHeader)
	require.NoError(t, err)
	require.Equal(t, true, decodedPrecomit.preValidated)
	err = decodedPrecomit.Validate()
	require.NoError(t, err)
	require.Equal(t, true, decodedPrecomit.validated)
}

func TestVerifyMaliciousAggregatedPrecommits(t *testing.T) {
	numOfFastAggPrecommits := 10 + rand.Intn(cSize)
	round := int64(0)
	var precommits []*message.Precommit
	for n := 0; n < numOfFastAggPrecommits; n++ {
		value := values[n%len(values)]
		// add duplicated msg but with different signers, thus the aggregation need to do a further fast aggregate.
		precommits = append(precommits, fastAggregatedPrecommit(height, round+int64(n), value, randomSigners(cSize)))
		precommits = append(precommits, fastAggregatedPrecommit(height, round+int64(n), value, randomSigners(cSize)))
	}
	precommits = append(precommits, fastAggregatedPrecommit(height, round+int64(numOfFastAggPrecommits), values[0], randomSigners(cSize)))

	t.Run("with wrong height", func(t *testing.T) {
		wrongHeight := height + 1
		aggPrecommits := maliciousAggregatePrecommits(precommits, &wrongHeight, nil, nil, nil)
		payload, err := rlp.EncodeToBytes(aggPrecommits)
		require.NoError(t, err)

		decodedPrecomit := &HighlyAggregatedPrecommit{}
		err = rlp.DecodeBytes(payload, decodedPrecomit)
		require.NoError(t, err)
		header := newBlockHeader(height-1, committee)
		err = decodedPrecomit.PreValidate(header)
		require.Error(t, errBadHeight, err)
		require.Equal(t, false, decodedPrecomit.preValidated)
	})

	t.Run("with wrong round", func(t *testing.T) {
		wrongRound := int64(0)
		aggPrecommits := maliciousAggregatePrecommits(precommits, nil, &wrongRound, nil, nil)
		payload, err := rlp.EncodeToBytes(aggPrecommits)
		require.NoError(t, err)

		decodedPrecomit := &HighlyAggregatedPrecommit{}
		err = rlp.DecodeBytes(payload, decodedPrecomit)
		require.NoError(t, err)
		header := newBlockHeader(height-1, committee)
		err = decodedPrecomit.PreValidate(header)
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
		err = decodedPrecomit.PreValidate(header)
		require.NoError(t, err)
		require.Equal(t, true, decodedPrecomit.preValidated)
		err = decodedPrecomit.Validate()
		require.NotNil(t, err)
		require.Equal(t, false, decodedPrecomit.validated)
	})

	t.Run("with wrong signers", func(t *testing.T) {
		wrongSigners := randomSigners(cSize)
		aggPrecommits := maliciousAggregatePrecommits(precommits, nil, nil, nil, wrongSigners)
		payload, err := rlp.EncodeToBytes(aggPrecommits)
		require.NoError(t, err)

		decodedPrecomit := &HighlyAggregatedPrecommit{}
		err = rlp.DecodeBytes(payload, decodedPrecomit)
		require.NoError(t, err)
		header := newBlockHeader(height-1, committee)
		err = decodedPrecomit.PreValidate(header)
		require.NoError(t, err)
		require.Equal(t, true, decodedPrecomit.preValidated)
		err = decodedPrecomit.Validate()
		require.NotNil(t, err)
		require.Equal(t, false, decodedPrecomit.validated)
	})
}

func fastAggregatedPrecommit(h uint64, r int64, v common.Hash, signers []int) *message.Precommit {
	precommits := make([]message.Vote, len(signers))
	for i, s := range signers {
		precommits[i] = newValidatedPrecommit(r, h, v, makeSigner(keys[s]), &committee[s], cSize)
	}
	return message.AggregatePrecommits(precommits)
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
	numOfFastAggPrecommits := 10 + rand.Intn(cSize)

	var precommits []*message.Precommit
	for n := 0; n < numOfFastAggPrecommits; n++ {
		value := values[n%len(values)]
		// add duplicated msg but with different signers, thus the aggregation need to do a further fast aggregate.
		precommits = append(precommits, fastAggregatedPrecommit(height, round+int64(n), value, randomSigners(cSize)))
		precommits = append(precommits, fastAggregatedPrecommit(height, round+int64(n), value, randomSigners(cSize)))
	}
	precommits = append(precommits, fastAggregatedPrecommit(height, round+int64(numOfFastAggPrecommits), values[0], randomSigners(cSize)))
	return AggregateDistinctPrecommits(precommits)
}

// maliciousAggregatePrecommits aggregate the precomits in a wrong way by modifying meta-data.
func maliciousAggregatePrecommits(precommits []*message.Precommit, wrongHeight *uint64, wrongRound *int64,
	wrongValue *common.Hash, wrongSigners []int) HighlyAggregatedPrecommit {

	defaultHeight := precommits[0].H()
	if wrongHeight != nil {
		defaultHeight = *wrongHeight
	}

	var fastAggregatedPrecommits []*message.Precommit
	// first we filter out fast-aggregatable votes with same msg, and fast aggregate them into single one.
	var fastAggregatablePrecommits []message.Vote
	fastAggregatablePrecommits = append(fastAggregatablePrecommits, precommits[0])

	for i := 1; i < len(precommits); i++ {

		if fastAggregatablePrecommits[0].R() == precommits[i].R() &&
			fastAggregatablePrecommits[0].Value() == precommits[i].Value() {

			fastAggregatablePrecommits = append(fastAggregatablePrecommits, precommits[i])

		} else {
			// if there are multiple fast aggregatable precommits:  aggregate them into single one and append the
			// aggregated one to the fastAggregatedPrecommits, then reset fastAggregatablePrecommits with current
			// precommit and continue.
			if len(fastAggregatablePrecommits) > 1 {
				aggregatedMsg := message.AggregatePrecommits(fastAggregatablePrecommits)
				fastAggregatedPrecommits = append(fastAggregatedPrecommits, aggregatedMsg)
				fastAggregatablePrecommits = []message.Vote{precommits[i]}
				continue
			}

			// if there is only one precommit in the fastAggregatablePrecommits: append it to the fastAggregatedPrecommits.
			// then reset fastAggregatablePrecommits with current precommit and continue.
			fastAggregatedPrecommits = append(fastAggregatedPrecommits, fastAggregatablePrecommits[0].(*message.Precommit))
			fastAggregatablePrecommits = []message.Vote{precommits[i]}
		}
	}

	// append the last patch of fastAggregatablePrecommits to fastAggregatedPrecommits.
	if len(fastAggregatablePrecommits) > 1 {
		aggregatedMsg := message.AggregatePrecommits(fastAggregatablePrecommits)
		fastAggregatedPrecommits = append(fastAggregatedPrecommits, aggregatedMsg)
	} else {
		fastAggregatedPrecommits = append(fastAggregatedPrecommits, fastAggregatablePrecommits[0].(*message.Precommit))
	}

	result := HighlyAggregatedPrecommit{}
	signatures := make([]blst.Signature, len(fastAggregatedPrecommits))
	for i, m := range fastAggregatedPrecommits {
		defaultRound := m.R()
		defaultValue := m.Value()
		defaultSingers := m.Signers().Flatten()
		if wrongRound != nil {
			defaultRound = *wrongRound
		}
		if wrongValue != nil {
			defaultValue = *wrongValue
		}

		if len(wrongSigners) > 0 {
			defaultSingers = wrongSigners
		}

		roundValueSigners := &RoundValueSigners{
			Round:   defaultRound,
			Value:   defaultValue,
			Signers: defaultSingers,
		}
		result.RoundValueSigners = append(result.RoundValueSigners, roundValueSigners)
		signatures[i] = m.Signature()
	}
	result.Height = defaultHeight
	result.Signature = blst.AggregateSignatures(signatures).Marshal()
	return result
}
