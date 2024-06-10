package message

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/autonity/autonity/common"
)

func TestAggregatedPower(t *testing.T) {
	aggregatedPower := NewAggregatedPower()

	aggregatedPower.Set(0, common.Big1)
	require.Equal(t, common.Big1, aggregatedPower.Power())
	require.Equal(t, uint(1), aggregatedPower.Signers().Bit(0))

	aggregatedPower.Set(0, common.Big1)
	require.Equal(t, common.Big1, aggregatedPower.Power())
	require.Equal(t, uint(1), aggregatedPower.Signers().Bit(0))

	aggregatedPower.Set(1, common.Big2)
	require.Equal(t, common.Big3, aggregatedPower.Power())
	require.Equal(t, uint(1), aggregatedPower.Signers().Bit(0))
	require.Equal(t, uint(1), aggregatedPower.Signers().Bit(1))

	aggregatedPower.Set(5, common.Big3)
	require.Equal(t, big.NewInt(6), aggregatedPower.Power())
	require.Equal(t, uint(1), aggregatedPower.Signers().Bit(0))
	require.Equal(t, uint(1), aggregatedPower.Signers().Bit(1))
	require.Equal(t, uint(0), aggregatedPower.Signers().Bit(2))
	require.Equal(t, uint(0), aggregatedPower.Signers().Bit(3))
	require.Equal(t, uint(0), aggregatedPower.Signers().Bit(4))
	require.Equal(t, uint(1), aggregatedPower.Signers().Bit(5))
}

func TestContribution(t *testing.T) {
	aggregator := NewAggregatedPower()
	core := NewAggregatedPower()

	aggregator.Set(0, common.Big1)
	contribution := Contribution(aggregator.Signers(), core.Signers())
	require.Equal(t, uint(1), contribution.Bit(0))
	require.NotEqual(t, 0, contribution.Cmp(common.Big0))

	core.Set(0, common.Big1)
	contribution = Contribution(aggregator.Signers(), core.Signers())
	require.Equal(t, uint(0), contribution.Bit(0))
	require.Equal(t, 0, contribution.Cmp(common.Big0))

	aggregator.Set(1000000000, common.Big1)
	contribution = Contribution(aggregator.Signers(), core.Signers())
	require.Equal(t, uint(0), contribution.Bit(0))
	require.Equal(t, uint(0), contribution.Bit(100))
	require.Equal(t, uint(1), contribution.Bit(1000000000))
	require.NotEqual(t, 0, contribution.Cmp(common.Big0))

	core.Set(1000000001, common.Big1)
	contribution = Contribution(aggregator.Signers(), core.Signers())
	require.Equal(t, uint(0), contribution.Bit(0))
	require.Equal(t, uint(0), contribution.Bit(100))
	require.Equal(t, uint(1), contribution.Bit(1000000000))
	require.Equal(t, uint(0), contribution.Bit(1000000001))
	require.NotEqual(t, 0, contribution.Cmp(common.Big0))

	core.Set(1000000000, common.Big1)
	contribution = Contribution(aggregator.Signers(), core.Signers())
	require.Equal(t, uint(0), contribution.Bit(0))
	require.Equal(t, uint(0), contribution.Bit(100))
	require.Equal(t, uint(0), contribution.Bit(1000000000))
	require.Equal(t, uint(0), contribution.Bit(1000000001))
	require.Equal(t, 0, contribution.Cmp(common.Big0))

	aggregator.Set(1, common.Big1)
	core.Set(4, common.Big1)
	core.Set(107, common.Big1)
	core.Set(64, common.Big1)
	contribution = Contribution(aggregator.Signers(), core.Signers())
	require.Equal(t, uint(0), contribution.Bit(0))
	require.Equal(t, uint(1), contribution.Bit(1))
	require.Equal(t, uint(0), contribution.Bit(4))
	require.Equal(t, uint(0), contribution.Bit(64))
	require.Equal(t, uint(0), contribution.Bit(100))
	require.Equal(t, uint(0), contribution.Bit(107))
	require.Equal(t, uint(0), contribution.Bit(1000000000))
	require.Equal(t, uint(0), contribution.Bit(1000000001))
	require.NotEqual(t, 0, contribution.Cmp(common.Big0))

	core.Set(1, common.Big1)
	contribution = Contribution(aggregator.Signers(), core.Signers())
	require.Equal(t, uint(0), contribution.Bit(0))
	require.Equal(t, uint(0), contribution.Bit(1))
	require.Equal(t, uint(0), contribution.Bit(4))
	require.Equal(t, uint(0), contribution.Bit(64))
	require.Equal(t, uint(0), contribution.Bit(100))
	require.Equal(t, uint(0), contribution.Bit(107))
	require.Equal(t, uint(0), contribution.Bit(1000000000))
	require.Equal(t, uint(0), contribution.Bit(1000000001))
	require.Equal(t, 0, contribution.Cmp(common.Big0))
}
