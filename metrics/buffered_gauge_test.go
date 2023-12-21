package metrics

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func BenchmarkBufferedGauge(b *testing.B) {
	g := NewBufferedGauge()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g.Add(int64(i)) // worst case scenario, we need to reallocate multiple times to increase capacity
	}
}

func TestBufferedGauge(t *testing.T) {
	t.Skip("")
	g := NewBufferedGauge()
	g.Add(int64(47))
	g.Add(int64(12))
	require.Equal(t, 2, g.Len())
	values := g.Values()
	require.Equal(t, int64(47), values[0].Value())
	require.Equal(t, int64(12), values[1].Value())
	g.Clear()
	require.Equal(t, 0, g.Len())
}

func TestBufferedGaugeSnapshot(t *testing.T) {
	t.Skip("")
	g := NewBufferedGauge()
	g.Add(int64(47))
	snapshot := g.Snapshot()
	g.Clear()
	require.Equal(t, 0, g.Len())
	require.Equal(t, 1, snapshot.Len())
	require.Equal(t, int64(47), snapshot.Values()[0].Value())
}

func TestBufferedGaugeSnapshotAndClear(t *testing.T) {
	t.Skip("")
	g := NewBufferedGauge()
	g.Add(int64(47))
	g.Add(int64(48))
	g.Add(int64(49))
	snapshot := g.SnapshotAndClear()
	require.Equal(t, 0, g.Len())
	require.Equal(t, 3, snapshot.Len())
	values := snapshot.Values()
	require.Equal(t, int64(47), values[0].Value())
	require.Equal(t, int64(48), values[1].Value())
	require.Equal(t, int64(49), values[2].Value())
}

func TestGetOrRegisterBufferedGauge(t *testing.T) {
	t.Skip("")
	r := NewRegistry()
	NewRegisteredBufferedGauge("foo", r).Add(int64(47))
	g := GetOrRegisterBufferedGauge("foo", r)
	require.Equal(t, 1, g.Len())
	require.Equal(t, int64(47), g.Values()[0].Value())
}
