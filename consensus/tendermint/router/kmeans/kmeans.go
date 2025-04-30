// Package kmeans implements the k-means clustering algorithm
// See: https://en.wikipedia.org/wiki/K-means_clustering
package kmeans

import (
	"fmt"
	"math/rand"
)

const (
	DefaultDeltaThreshold     = 0.01
	DefaultIterationThreshold = 96
)

// Kmeans configuration/option struct
type Kmeans struct {
	// deltaThreshold (in percent between 0.0 and 0.1) aborts processing if
	// less than n% of data points shifted clusters in the last iteration
	deltaThreshold float64
	// iterationThreshold aborts processing when the specified amount of
	// algorithm iterations was reached
	iterationThreshold int
}

// NewKmeansWithOptions returns a Kmeans configuration struct with custom settings
func NewKmeansWithOptions(deltaThreshold float64, iterationThreshold int) (*Kmeans, error) {
	if deltaThreshold <= 0.0 || deltaThreshold >= 1.0 {
		return nil, fmt.Errorf("delta threshold is out of bounds (must be > 0.0 and < 1.0)")
	}

	if iterationThreshold < 0 {
		return nil, fmt.Errorf("iteration threshold is out of bounds (must be >= 0)")
	}

	return &Kmeans{
		deltaThreshold:     deltaThreshold,
		iterationThreshold: iterationThreshold,
	}, nil
}

// New returns a Kmeans configuration struct with default settings
func New() *Kmeans {
	km, _ := NewKmeansWithOptions(DefaultDeltaThreshold, DefaultIterationThreshold)
	return km
}

// Partition executes the k-means algorithm on the given dataset and
// partitions it into k clusters
func (m *Kmeans) Partition(dataset Observations, k int, seed int64) (Clusters, error) {
	if k > len(dataset) {
		return Clusters{}, fmt.Errorf("the size of the data set must at least equal to k (%d)", k)
	}

	cc, err := NewClusters(seed, k, dataset)
	if err != nil {
		return cc, err
	}

	points := make([]int, len(dataset))
	changes := 1

	for i := 0; changes > 0; i++ {
		changes = 0
		cc.Reset()

		for p, point := range dataset {
			ci := cc.Nearest(point)
			cc[ci].Append(point)
			if points[p] != ci {
				points[p] = ci
				changes++
			}
		}

		for ci := 0; ci < len(cc); ci++ {
			if len(cc[ci].Observations) == 0 {
				// During the iterations, if any of the cluster centers has no
				// data points associated with it, assign a random data point
				// to it.
				// Also see: http://user.ceng.metu.edu.tr/~tcan/ceng465_f1314/Schedule/KMeansEmpty.html
				var ri int
				for {
					// find a cluster with at least two data points, otherwise
					// we're just emptying one cluster to fill another
					ri = rand.Intn(len(dataset))
					if len(cc[points[ri]].Observations) > 1 {
						break
					}
				}
				cc[ci].Append(dataset[ri])
				points[ri] = ci

				// Ensure that we always see at least one more iteration after
				// randomly assigning a data point to a cluster
				changes = len(dataset)
			}
		}

		if changes > 0 {
			cc.Recenter()
		}
		if i == m.iterationThreshold ||
			changes < int(float64(len(dataset))*m.deltaThreshold) {
			// return Clusters{}, fmt.Errorf("iteration threshold '%d' reached", m.iterationThreshold)
			break
		}
	}

	return cc, nil
}
