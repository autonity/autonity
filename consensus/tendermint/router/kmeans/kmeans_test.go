package kmeans

import (
	"math"
	"testing"
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

	clusters, err := km.Partition(dataset, k)
	if err != nil {
		t.Fatalf("Partition failed: %v", err)
	}

	if len(clusters) != k {
		t.Fatalf("Expected %d clusters, got %d", k, len(clusters))
	}

	// Additional checks can be added here to verify cluster contents
}

// Test Partition with invalid k
func TestPartitionInvalidK(t *testing.T) {
	km := New()
	dataset := Observations{
		TestObservation{1.0, 2.0},
		TestObservation{1.5, 1.8},
	}

	k := 3 // Invalid k

	_, err := km.Partition(dataset, k)
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
