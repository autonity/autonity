package kmeans

import (
	"math"
	"math/rand"
	"testing"
	"time"
)

// Define a concrete type for testing that implements the Observation interface
type TestObservation Coordinates

// Ensure TestObservation implements the Observation interface
func (o TestObservation) Coordinates() Coordinates {
	return Coordinates(o)
}

func (o TestObservation) Distance(point Coordinates) float64 {
	return math.Sqrt(Coordinates(o).Distance(point))
}

// Test NewKmeansWithOptions function
func TestNewKmeansWithOptions(t *testing.T) {
	tests := []struct {
		deltaThreshold     float64
		iterationThreshold int
		expectError        bool
	}{
		{0.05, 100, false},
		{1.1, 100, true}, // Invalid delta
		{0.05, -1, true}, // Invalid iteration
	}

	for _, test := range tests {
		_, err := NewKmeansWithOptions(test.deltaThreshold, test.iterationThreshold)
		if (err != nil) != test.expectError {
			t.Errorf("NewKmeansWithOptions(%v, %v) error = %v, expectError %v", test.deltaThreshold, test.iterationThreshold, err, test.expectError)
		}
	}
}

// Test Partition function with a simple dataset
func TestPartition(t *testing.T) {
	km := New()
	dataset := Observations{
		TestObservation{1.0, 2.0},
		TestObservation{1.5, 1.8},
		TestObservation{5.0, 8.0},
		TestObservation{8.0, 8.0},
	}
	k := 2
	seed := int64(42)

	clusters, err := km.Partition(dataset, k, seed)
	if err != nil {
		t.Fatalf("Partition failed: %v", err)
	}

	if len(clusters) != k {
		t.Fatalf("Expected %d clusters, got %d", k, len(clusters))
	}

	// Additional checks can be added here to verify cluster contents
}

// Benchmark partition function with different scale of nodes
func BenchmarkPartition(t *testing.B) {
	tests := []struct {
		name     string
		numNodes int
	}{
		{"100 nodes", 100},
		{"200 nodes", 200},
		{"400 nodes", 400},
		{"800 nodes", 800},
		{"1000 nodes", 1000},
		{"1200 nodes", 1200},
		{"1400 nodes", 1400},
		{"1600 nodes", 1600},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.B) {
			km := New()
			dataset := make(Observations, tt.numNodes)
			// Populate dataset with random observations
			for i := 0; i < tt.numNodes; i++ {
				dataset[i] = TestObservation{float64(rand.Intn(100)), float64(rand.Intn(100))}
			}
			k := int(math.Sqrt(float64(tt.numNodes))) // Set k to the square root of the size of the observations
			seed := int64(42)

			// Measure the time taken for the benchmark
			start := time.Now() // Start the timer
			for i := 0; i < t.N; i++ {
				_, err := km.Partition(dataset, k, seed)
				if err != nil {
					t.Fatalf("Partition failed: %v", err)
				}
			}
			duration := time.Since(start) // Calculate duration

			// Print the average time per operation in milliseconds
			avgTime := duration.Milliseconds() / int64(t.N)
			t.Logf("Average time for %s: %d ms", tt.name, avgTime)
		})
	}
}

// Test Partition with invalid k
func TestPartitionInvalidK(t *testing.T) {
	km := New()
	dataset := Observations{
		TestObservation{1.0, 2.0},
		TestObservation{1.5, 1.8},
	}

	k := 3 // Invalid k
	seed := int64(42)

	_, err := km.Partition(dataset, k, seed)
	if err == nil {
		t.Fatal("Expected error for invalid k, got nil")
	}
}

// Test default settings
func TestNew(t *testing.T) {
	km := New()
	if km.deltaThreshold != DefaultDeltaThreshold {
		t.Errorf("Expected deltaThreshold %v, got %v", DefaultDeltaThreshold, km.deltaThreshold)
	}
	if km.iterationThreshold != DefaultIterationThreshold {
		t.Errorf("Expected iterationThreshold %v, got %v", DefaultIterationThreshold, km.iterationThreshold)
	}
}
