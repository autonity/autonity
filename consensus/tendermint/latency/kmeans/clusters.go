package kmeans

import (
	"fmt"
	"math/rand"
	"time"
)

// A Cluster which data points gravitate around
type Cluster struct {
	Center       Coordinates  `json:"center"`
	Observations Observations `json:"-"`
}

// Clusters is a slice of clusters
type Clusters []Cluster

// NewClusters sets up a new set of clusters and seeds their initial positions according
// to seed. If seed is 0 the seed is taken by the current time
func NewClusters(seed int64, k int, dataset Observations) (Clusters, error) {
	var c Clusters
	if len(dataset) == 0 || len(dataset[0].Coordinates()) == 0 {
		return c, fmt.Errorf("there must be at least one dimension in the data set")
	}
	if k == 0 {
		return c, fmt.Errorf("k must be greater than 0")
	}

	var r *rand.Rand
	if seed != 0 {
		r = rand.New(rand.NewSource(seed))
	} else {
		r = rand.New(rand.NewSource(time.Now().UnixNano()))
	}

	for i := 0; i < k; i++ {
		var p Coordinates
		for j := 0; j < len(dataset[0].Coordinates()); j++ {
			p = append(p, r.Float64())
		}

		c = append(c, Cluster{
			Center: p,
		})
	}
	return c, nil
}

// Append adds an observation to the Cluster
func (c *Cluster) Append(point Observation) {
	c.Observations = append(c.Observations, point)
}

// Nearest returns the index of the cluster nearest to point
func (c *Clusters) Nearest(point Observation) int {
	var ci int
	dist := -1.0

	// Find the nearest cluster for this data point
	for i, cluster := range *c {
		d := point.Distance(cluster.Center)
		if dist < 0 || d < dist {
			dist = d
			ci = i
		}
	}

	return ci
}

// Neighbour returns the neighbouring cluster of a point along with the average distance to its points
func (c *Clusters) Neighbour(point Observation, fromCluster int) (int, float64) {
	var d float64
	nc := -1

	for i, cluster := range *c {
		if i == fromCluster {
			continue
		}

		cd := AverageDistance(point, cluster.Observations)
		if nc < 0 || cd < d {
			nc = i
			d = cd
		}
	}

	return nc, d
}

// Recenter recenters a cluster
func (c *Cluster) Recenter() {
	center, err := c.Observations.Center()
	if err != nil {
		return
	}

	c.Center = center
}

// Recenter recenters all clusters
func (c Clusters) Recenter() {
	for i := 0; i < len(c); i++ {
		c[i].Recenter()
	}
}

// Reset clears all point assignments
func (c Clusters) Reset() {
	for i := 0; i < len(c); i++ {
		c[i].Observations = Observations{}
	}
}
