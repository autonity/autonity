package common

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestContains(t *testing.T) {
	// slices.Contains compares using ==
	// which is fine for some data types
	// but not for *big.Int

	list := []*big.Int{Big1, Big2, Big3}

	require.True(t, Contains(list, Big2))
	// different pointer, but still contained
	require.True(t, Contains(list, new(big.Int).SetUint64(2)))
	require.False(t, Contains(list, new(big.Int).SetUint64(56)))
}
