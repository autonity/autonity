// Package kmeans implements the k-means clustering algorithm
// See: https://en.wikipedia.org/wiki/K-means_clustering
package kmeans

import (
	"fmt"
	"sort"
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
// partitions it into k clusters, those small clusters will be merged into their nearest cluster.
func (m *Kmeans) Partition(dataset Observations, k int) (Clusters, error) {
	if k > len(dataset) {
		return Clusters{}, fmt.Errorf("the size of the data set must at least equal to k (%d)", k)
	}

	cc, err := NewClusters(k, dataset)
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
			// find the nearest point rather than a random point to make the algorithm deterministic.
			if len(cc[ci].Observations) == 0 {
				minDist := -1.0
				nearestPointIndex := -1
				for p, point := range dataset {
					dist := cc[ci].Center.Distance(point.Coordinates())
					if minDist < 0 || dist < minDist {
						minDist = dist
						nearestPointIndex = p
					}
				}
				if nearestPointIndex != -1 {
					cc[ci].Append(dataset[nearestPointIndex])
					points[nearestPointIndex] = ci
					changes = len(dataset)
				}
			}
		}

		if changes > 0 {
			cc.Recenter()
		}
		if i == m.iterationThreshold ||
			changes < int(float64(len(dataset))*m.deltaThreshold) {
			break
		}
	}

	// resolve params for re-clustering
	optimalSize := k
	smallClusterThreshold := optimalSize / 3
	if smallClusterThreshold < 3 {
		smallClusterThreshold = 3
	}

	// merge small clusters into their nearest cluster.
	cc = mergeSmallClusters(cc, smallClusterThreshold)

	return cc, nil
	//return equalize(cc, dataset), nil
}

// equalize the clusters to have the same size
func equalize(cc Clusters, dataset Observations) Clusters {
	clusterSize := len(dataset) / len(cc)
	results := make(Clusters, len(cc))
	for i := 0; i < len(cc); i++ {
		results[i].Center = cc[i].Center
	}

	type distance struct {
		cluster  int
		distance float64
	}
	distances := func(point Observation) []distance {
		d := make([]distance, len(results))
		for i, cluster := range results {
			d[i] = distance{
				cluster:  i,
				distance: point.Distance(cluster.Center),
			}
		}
		sort.Slice(d, func(i, j int) bool {
			return d[i].distance < d[j].distance
		})
		return d
	}
	full := make(map[int]bool)

	for _, point := range dataset {
		target := 0
		dists := distances(point)
		for i := 0; i < len(dists); i++ {
			if !full[dists[i].cluster] {
				target = dists[i].cluster
				break
			}
			target = dists[0].cluster
		}
		results[target].Observations = append(results[target].Observations, point)
		if len(results[target].Observations) >= clusterSize {
			full[target] = true
		}
	}

	return results
}

// mergeSmallClusters merge small clusters into their nearest cluster.
func mergeSmallClusters(cc Clusters, smallClusterThreshold int) Clusters {
	// merge small ones into their nearest cluster.
	for i := 0; i < len(cc); i++ {
		if len(cc[i].Observations) < smallClusterThreshold {
			nearest := findNearestNonSmallCluster(cc, i, smallClusterThreshold)
			if nearest != -1 {
				cc[nearest].Observations = append(cc[nearest].Observations, cc[i].Observations...)
				cc[i].Observations = nil
			}
		}
	}

	// remove empty clusters
	newCC := make(Clusters, 0)
	for _, c := range cc {
		if len(c.Observations) > 0 {
			newCC = append(newCC, c)
		}
	}

	return newCC
}

// findNearestNonSmallCluster
func findNearestNonSmallCluster(cc Clusters, index, smallClusterThreshold int) int {
	minDist := -1.0
	nearest := -1
	for i := 0; i < len(cc); i++ {
		if i != index && len(cc[i].Observations) >= smallClusterThreshold {
			dist := cc[index].Center.Distance(cc[i].Center.Coordinates())
			if minDist < 0 || dist < minDist {
				minDist = dist
				nearest = i
			}
		}
	}
	return nearest
}
