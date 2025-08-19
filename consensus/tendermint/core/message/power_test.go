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

func TestSetReturnValue(t *testing.T) {
	p := NewAggregatedPower()

	require.True(t, p.Set(0, common.Big1))
	require.True(t, p.Set(1, common.Big4))

	require.False(t, p.Set(0, common.Big5))
	require.False(t, p.Set(1, common.Big3))

	require.True(t, p.Set(3, common.Big4))

}
