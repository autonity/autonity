package metrics

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func BenchmarkBufferedGauge(b *testing.B) {
	capacity := 20
	g := NewBufferedGauge(&capacity)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g.Add(int64(i)) // worst case scenario, we need to reallocate multiple times to increase capacity
	}
}

func TestBufferedGauge(t *testing.T) {
	capacity := 20
	g := NewBufferedGauge(&capacity)
	g.Add(int64(47))
	g.Add(int64(12))
	require.Equal(t, 2, g.Len())
	values := g.Values()
	require.Equal(t, int64(47), values[0].Value())
	require.Equal(t, int64(12), values[1].Value())
	g.Clear()
	require.Equal(t, 0, g.Len())
}

func TestLowCapacityBufferedGauge(t *testing.T) {
	t.Run("0 capacity", func(t *testing.T) {
		capacity := 0
		g := NewBufferedGauge(&capacity)
		g.Add(int64(47))
		require.Equal(t, 1, g.Len())
		values := g.Values()
		require.Equal(t, int64(47), values[0].Value())
		g.Clear()
		require.Equal(t, 0, g.Len())
	})

	t.Run("-ve capacity", func(t *testing.T) {
		capacity := -1
		g := NewBufferedGauge(&capacity)
		g.Add(int64(47))
		require.Equal(t, 1, g.Len())
		values := g.Values()
		require.Equal(t, int64(47), values[0].Value())
		g.Clear()
		require.Equal(t, 0, g.Len())
	})

	t.Run("1 capacity", func(t *testing.T) {
		capacity := 1
		g := NewBufferedGauge(&capacity)
		g.Add(int64(47))
		require.Equal(t, 1, g.Len())
		values := g.Values()
		require.Equal(t, int64(47), values[0].Value())
		g.Clear()
		require.Equal(t, 0, g.Len())
	})
}

func TestBufferedGaugeOversized(t *testing.T) {
	capacity := 5
	g := NewBufferedGauge(&capacity)
	g.Add(int64(47))
	g.Add(int64(12))
	g.Add(int64(13))
	g.Add(int64(14))
	g.Add(int64(15))
	require.Equal(t, 5, g.Len())
	values := g.Values()
	require.Equal(t, int64(47), values[0].Value())
	require.Equal(t, int64(12), values[1].Value())
	require.Equal(t, int64(15), values[4].Value())
	g.Add(int64(16))
	g.Add(int64(17))
	g.Add(int64(18))
	g.Add(int64(19))
	g.Add(int64(20))
	require.Equal(t, int64(16), values[0].Value())
	require.Equal(t, int64(20), values[4].Value())
	g.Add(int64(21))
	require.Equal(t, int64(21), values[0].Value())
	g.Clear()
	require.Equal(t, 0, g.Len())
	g.Add(int64(23))
	g.Add(int64(24))
	g.Add(int64(18))
	g.Add(int64(19))
	g.Add(int64(29))
	require.Equal(t, int64(23), values[0].Value())
	require.Equal(t, int64(29), values[4].Value())
}

func TestBufferedGaugeSnapshot(t *testing.T) {
	capacity := 20
	g := NewBufferedGauge(&capacity)
	g.Add(int64(47))
	snapshot := g.Snapshot()
	g.Clear()
	require.Equal(t, 0, g.Len())
	require.Equal(t, 1, snapshot.Len())
	require.Equal(t, int64(47), snapshot.Values()[0].Value())
}

func TestBufferedGaugeSnapshotAndClear(t *testing.T) {
	capacity := 20
	g := NewBufferedGauge(&capacity)
	g.Add(int64(47))
	g.Add(int64(48))
	g.Add(int64(49))
	snapshot := g.SnapshotAndClear()
	values := snapshot.Values()
	require.Equal(t, int64(47), values[0].Value())
	require.Equal(t, int64(48), values[1].Value())
	require.Equal(t, int64(49), values[2].Value())
}

func TestGetOrRegisterBufferedGauge(t *testing.T) {
	r := NewRegistry()
	NewRegisteredBufferedGauge("foo", r, nil).Add(int64(47))
	g := GetOrRegisterBufferedGauge("foo", r)
	require.Equal(t, 1, g.Len())
	require.Equal(t, int64(47), g.Values()[0].Value())
}
