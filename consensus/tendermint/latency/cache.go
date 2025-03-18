package latency

import (
	"sync"
)

type clusterCache struct {
	mu sync.RWMutex
	//clusters map[uint64]*Clusters
	// We keep recent two epoches clusters, thus that on epoch rotation old message can be handled.
	// Each epoch has two clustering views, the 1st one is the default clustering before the optimized one.
	// The 2nd one is the optimized one that is activated on a specific height within that epoch.
	lastEpochClusters    []*Clusters
	currentEpochClusters []*Clusters
}

func newClusterCache() *clusterCache {
	return &clusterCache{}
}

func (c *clusterCache) insertClustering(cluster *Clusters) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// if current epoch clusters haven't being fully filled. Fill them.
	if c.currentEpochClusters == nil || len(c.currentEpochClusters) == 1 {
		c.lastEpochClusters = append(c.lastEpochClusters, cluster)
		return
	}

	// otherwise, this is an epoch rotation, reset last epoch clusters with current one, and add new epoch's cluster.
	c.lastEpochClusters = make([]*Clusters, len(c.currentEpochClusters))
	copy(c.lastEpochClusters, c.currentEpochClusters)
	c.currentEpochClusters = []*Clusters{cluster}
}

// clustersAt returns the clusters at a given block height, which will be the clusters at
// the last block height before the given block.
func (c *clusterCache) clustersAt(block uint64) (*Clusters, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// if the height belong to current epoch.
	if len(c.currentEpochClusters) > 0 {
		// the height is in the latest view
		latestClusters := c.currentEpochClusters[len(c.currentEpochClusters)-1]
		if block >= latestClusters.activatedHeight && block < latestClusters.nextEpochHeight {
			return latestClusters, true
		}

		// the height is in the default clustering range.
		if len(c.currentEpochClusters) == 2 {
			defaultClusters := c.currentEpochClusters[0]
			if block >= defaultClusters.activatedHeight && block < defaultClusters.nextEpochHeight {
				return defaultClusters, true
			}
		}
	}

	// if the height belong to last epoch, we just check from the last cluster of last epoch.
	if len(c.lastEpochClusters) > 0 {
		lastEpochClusters := c.lastEpochClusters[len(c.lastEpochClusters)-1]
		if block >= lastEpochClusters.activatedHeight && block < lastEpochClusters.nextEpochHeight {
			return lastEpochClusters, true
		}
	}

	// in any other case, we return nil and false.
	return nil, false
}
