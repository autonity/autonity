package kmeans

import (
	"testing"
)

// Test Center function
func TestCenter(t *testing.T) {
	tests := []struct {
		observations Observations
		expected     Coordinates
		expectError  bool
	}{
		{
			observations: Observations{
				TestObservation{1.0, 2.0},
				TestObservation{3.0, 4.0},
				TestObservation{5.0, 6.0},
			},
			expected:    Coordinates{3.0, 4.0},
			expectError: false,
		},
		{
			observations: Observations{},
			expected:     nil,
			expectError:  true,
		},
	}

	for _, test := range tests {
		center, err := test.observations.Center()
		if (err != nil) != test.expectError {
			t.Errorf("Center() error = %v, expectError %v", err, test.expectError)
		}
		if !test.expectError && !equalCoordinates(center, test.expected) {
			t.Errorf("Center() = %v, want %v", center, test.expected)
		}
	}
}

// Test AverageDistance function
func TestAverageDistance(t *testing.T) {
	observations := Observations{
		TestObservation{1.0, 2.0},
		TestObservation{3.0, 4.0},
		TestObservation{5.0, 6.0},
	}
	cPoint := TestObservation{0.0, 0.0}

	// Calculate the average distance from observations to cPoint
	avgDistance := AverageDistance(cPoint, observations)

	// Calculate expected distances manually
	expectedDistances := []float64{
		cPoint.Distance(observations[0].Coordinates()), // Distance to (1.0, 2.0)
		cPoint.Distance(observations[1].Coordinates()), // Distance to (3.0, 4.0)
		cPoint.Distance(observations[2].Coordinates()), // Distance to (5.0, 6.0)
	}

	// Calculate the expected average distance
	var totalDistance float64
	for _, d := range expectedDistances {
		totalDistance += d
	}
	expectedAverageDistance := totalDistance / float64(len(observations))

	if avgDistance != expectedAverageDistance {
		t.Errorf("AverageDistance() = %v, want %v", avgDistance, expectedAverageDistance)
	}
}

// Helper function to check equality of Coordinates
func equalCoordinates(a, b Coordinates) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
