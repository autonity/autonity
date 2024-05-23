package metrics

import (
	"sync"
	"time"
)

const (
	BufferedGaugeDefaultCapacity = 16
)

// value with its associated timestamp
type GaugeValue struct {
	value     int64
	timestamp time.Time
}

func (v GaugeValue) Value() int64         { return v.value }
func (v GaugeValue) Timestamp() time.Time { return v.timestamp }

// BufferedGauges holds a slice of int64 values. Used to get instantaneous metrics.
type BufferedGauge interface {
	Snapshot() BufferedGauge         // read only copy of the bufferedGauge
	SnapshotAndClear() BufferedGauge // read only copy of the bufferedGauge + clears underlying values
	Add(int64)                       // adds a new element to the values slice, timestamping it with the time of adding.
	Len() int                        // number of added values
	Clear()                          // clears the stored values. It is responsibility of the caller to clear the values once they are not needed anymore
	Values() []GaugeValue            // returns the values of the BufferedGauge
}

// GetOrRegisterBufferedGauge returns an existing BufferedGauge or constructs and registers a new StandardBufferedGauge.
func GetOrRegisterBufferedGauge(name string, r Registry) BufferedGauge {
	if nil == r {
		r = DefaultRegistry
	}
	return r.GetOrRegister(name, NewBufferedGauge(nil)).(BufferedGauge)
}

// NewBufferedGauge constructs a new BufferedGauge.
func NewBufferedGauge(capacity *int) BufferedGauge {
	if !Enabled {
		return NilBufferedGauge{}
	}
	var c int
	if capacity == nil {
		c = BufferedGaugeDefaultCapacity
		capacity = &c
	} else if capacity != nil && *capacity < 1 {
		c = 1 // minimum capacity
		capacity = &c
	}

	return &StandardBufferedGauge{values: make([]GaugeValue, 0, *capacity), capacity: *capacity}
}

// NewRegisteredBufferedGauge constructs and registers a new StandardBufferedGauge.
func NewRegisteredBufferedGauge(name string, r Registry, capacity *int) BufferedGauge {
	c := NewBufferedGauge(capacity)
	if nil == r {
		r = DefaultRegistry
	}
	r.Register(name, c)
	return c
}

// BufferedGaugeSnapshot is a read-only copy of another BufferedGauge.
type BufferedGaugeSnapshot struct {
	values []GaugeValue
}

// snapshot of a snapshot, it's the snasphot itself
func (g BufferedGaugeSnapshot) Snapshot() BufferedGauge { return g }

func (g BufferedGaugeSnapshot) SnapshotAndClear() BufferedGauge {
	panic("SnapshotAndClear called on a BufferedGaugeSnapshot")
}

func (g BufferedGaugeSnapshot) Add(_ int64) {
	panic("Add called on a BufferedGaugeSnapshot")
}

// length of the values at the time the snapshot was taken
func (g BufferedGaugeSnapshot) Len() int { return len(g.values) }

func (g BufferedGaugeSnapshot) Clear() {
	panic("Clear called on a BufferedGaugeSnapshot")
}

// Values returns the values at the time the snapshot was taken.
func (g BufferedGaugeSnapshot) Values() []GaugeValue { return g.values }

// NilBufferedGauge is a no-op Gauge.
type NilBufferedGauge struct{}

// Snapshot is a no-op.
func (NilBufferedGauge) Snapshot() BufferedGauge { return NilBufferedGauge{} }

// SnapshotAndClear is a no-op.
func (NilBufferedGauge) SnapshotAndClear() BufferedGauge { return NilBufferedGauge{} }

// Add is a no-op.
func (NilBufferedGauge) Add(_ int64) {}

// Len is a no-op.
func (NilBufferedGauge) Len() int { return 0 }

// Clear is a no-op.
func (NilBufferedGauge) Clear() {}

// Values is a no-op.
func (NilBufferedGauge) Values() []GaugeValue { return nil }

// StandardBufferedGauge is the standard implementation of a BufferedGauge
type StandardBufferedGauge struct {
	values   []GaugeValue
	capacity int
	rIndex   int
	sync.RWMutex
}

// Snapshot returns a read-only copy of the bufferedGauge.
func (g *StandardBufferedGauge) Snapshot() BufferedGauge {
	g.RLock()
	defer g.RUnlock()
	return BufferedGaugeSnapshot{values: g.values}
}

// SnapshotAndClear returns a read-only copy of the bufferedGauge and clears the underlying values
func (g *StandardBufferedGauge) SnapshotAndClear() BufferedGauge {
	g.Lock()
	defer g.Unlock()
	snapshot := BufferedGaugeSnapshot{values: g.values}
	g.values = g.values[:0]
	return snapshot
}

// Append a new value into the buffer
func (g *StandardBufferedGauge) Add(v int64) {
	g.Lock()
	defer g.Unlock()

	if len(g.values) < g.capacity {
		g.values = append(g.values, GaugeValue{value: v, timestamp: time.Now()})
	} else {
		g.rIndex = g.rIndex % g.capacity
		g.values[g.rIndex] = GaugeValue{value: v, timestamp: time.Now()}
		g.rIndex++
	}
}

func (g *StandardBufferedGauge) Len() int {
	g.RLock()
	defer g.RUnlock()
	return len(g.values)
}

func (g *StandardBufferedGauge) Clear() {
	g.Lock()
	defer g.Unlock()
	// we mantain the same capacity. Assuming that values will be read and cleared periodically, this approach will optimize reallocations
	g.values = g.values[:0]
}

// Value returns the gauge's current value.
func (g *StandardBufferedGauge) Values() []GaugeValue {
	g.Lock()
	defer g.Unlock()
	return g.values
}
