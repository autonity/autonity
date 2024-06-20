package accountability

import (
	cr "crypto/rand"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/crypto"
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
	aggregatedPrecommit := randomHighlyAggregatedPrecomits(h, r)
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

func randomHighlyAggregatedPrecomits(height uint64, round int64) HighlyAggregatedPrecommit {
	numOfFastAggPrecommits := 10 + rand.Intn(cSize)
	precommits := make([]*message.Precommit, numOfFastAggPrecommits)
	for n := 0; n < numOfFastAggPrecommits; n++ {
		value := values[n%len(values)]
		precommits[n] = fastAggregatedPrecommit(height, round+int64(n), value, randomSigners(cSize))
	}
	return AggregateDistinctPrecommits(precommits)
}
