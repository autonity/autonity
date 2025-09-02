package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBitmap(t *testing.T) {
	bitmap := NewBitmap()

	for i := 0; i < 10; i++ {
		require.Equal(t, false, bitmap.IsSet(i))
	}

	for i := 0; i < 10; i++ {
		bitmap.Set(i)
		require.Equal(t, true, bitmap.IsSet(i))
	}

	bitmap.Set(15)
	require.Equal(t, true, bitmap.IsSet(15))

	bitmap.Set(18)
	require.Equal(t, true, bitmap.IsSet(18))

	bitmap.Set(19)
	require.Equal(t, true, bitmap.IsSet(19))
}

func TestBitmapIterator(t *testing.T) {
	bitMap := NewBitmap()
	bi := NewBitmapIterator(bitMap)
	require.False(t, bi.Next())
	for i := 0; i < 4; i++ {
		bitMap.Set(i)
	}

	bi = NewBitmapIterator(bitMap)
	for i := 0; i < 4; i++ {
		require.True(t, bi.Next())
		require.Equal(t, i, bi.Index())
	}
	require.False(t, bi.Next())
	require.Equal(t, bitMap.Len(), bi.Index())
}
