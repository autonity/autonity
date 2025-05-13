package router

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/internal/testrand"
)

func TestClusterRotation(t *testing.T) {
	committeeLen := 180
	var committee []common.Address
	for i := 0; i < committeeLen; i++ {
		committee = append(committee, testrand.Address())
	}
	latMap := make(map[common.Address]uint)
	for _, addr := range committee {
		latMap[addr] = uint(rand.Intn(256))
	}

	// Create a new cluster
	t.Run("Test cluster initialization in default clusters", func(t *testing.T) {
		self := committee[0]
		// default clusters have a latMat of nil
		clusters, err := NewClusters(committee, latMap, nil, self)
		require.NoError(t, err)

		// check clusters are correctly assigned
		for i, committeeMember := range committee {
			if committeeMember == self {
				continue
			}
			clusterID := clusters.clusterContaining(committeeMember)
			require.Equal(t, clusterID, i%len(clusters.base))
			found := false
			foundCluster := 0
			for clusterId := range clusters.base {
				for _, member := range clusters.base[clusterId].Members {
					if member.Addr == committeeMember {
						found = true
						foundCluster = clusterId + 1
						break
					}
				}
			}
			require.True(t, found, "cluster member not found in cluster")
			require.Equal(t, foundCluster-1, clusterID, "cluster ID mismatch")
		}

		// check latencies are correctly assigned
		for _, cluster := range clusters.base {
			for _, member := range cluster.Members {
				if member.Addr == self {
					continue
				}
				latency := latMap[member.Addr]
				require.Equal(t, latency, member.Lat, "latency mismatch")
			}
		}

		// check clusters are sorted
		for _, cluster := range clusters.base {
			for i := 0; i < len(cluster.Members)-1; i++ {
				require.LessOrEqual(t, cluster.Members[i].Lat, cluster.Members[i+1].Lat, "clusters not sorted")
			}
		}
	})

	t.Run("Test cluster initialization in kmeans clusters", func(t *testing.T) {
		self := committee[1]
		latMat := make([][]uint8, committeeLen)
		for i := 0; i < committeeLen; i++ {
			latMat[i] = make([]uint8, committeeLen)
			for j := 0; j < i; j++ {
				if committee[i] == self {
					latMat[i][j] = uint8(latMap[committee[j]])
				} else {
					latMat[i][j] = uint8(rand.Intn(256))
					latMat[j][i] = latMat[i][j]
				}
			}
		}

		// Create a new cluster
		clusters, err := NewClusters(committee, latMap, latMat, self)
		require.NoError(t, err)

		// check clusters are correctly assigned
		for _, committeeMember := range committee {
			if committeeMember == self {
				continue
			}
			clusterID := clusters.clusterContaining(committeeMember)
			found := false
			foundCluster := 0
			for clusterId := range clusters.base {
				for _, member := range clusters.base[clusterId].Members {
					if member.Addr == committeeMember {
						found = true
						foundCluster = clusterId + 1
						break
					}
				}
			}
			require.True(t, found, "cluster member not found in cluster")
			require.Equal(t, foundCluster-1, clusterID, "cluster ID mismatch")
		}

		// check latencies are correctly assigned
		for _, cluster := range clusters.base {
			for _, member := range cluster.Members {
				if member.Addr == self {
					continue
				}
				latency := latMap[member.Addr]
				require.Equal(t, latency, member.Lat, "latency mismatch")
			}
		}

		// check clusters are sorted
		for _, cluster := range clusters.base {
			for i := 0; i < len(cluster.Members)-1; i++ {
				require.LessOrEqual(t, cluster.Members[i].Lat, cluster.Members[i+1].Lat, "clusters not sorted")
			}
		}
	})

	t.Run("Test cluster latency update", func(t *testing.T) {
		self := committee[2]
		latMat := make([][]uint8, committeeLen)
		for i := 0; i < committeeLen; i++ {
			latMat[i] = make([]uint8, committeeLen)
			for j := 0; j < i; j++ {
				if committee[i] == self {
					latMat[i][j] = uint8(latMap[committee[j]])
				} else {
					latMat[i][j] = uint8(rand.Intn(256))
					latMat[j][i] = latMat[i][j]
				}
			}
		}

		clusters, err := NewClusters(committee, latMap, latMat, self)
		require.NoError(t, err)

		clusterMap := make(map[common.Address]int)
		for _, addr := range committee {
			cid := clusters.clusterContaining(addr)
			require.NotEqual(t, cid, -1, "cluster ID should not be -1")
			clusterMap[addr] = cid
		}

		newLatMap := make(map[common.Address]uint)
		for _, addr := range committee {
			newLatMap[addr] = uint(rand.Intn(256))
		}
		clusters = UpdateClusterLatencies(clusters, newLatMap, self)
		for clusterId, cluster := range clusters.base {
			for _, member := range cluster.Members {
				if member.Addr == self {
					continue
				}
				newLatency := newLatMap[member.Addr]
				require.NotEqual(t, member.Lat, latMap[member.Addr], "latency should be updated")
				require.Equal(t, newLatency, member.Lat, "latency mismatch after update")
				require.Equal(t, clusterId, clusterMap[member.Addr], "cluster ID mismatch after update")
			}
		}
	})

	t.Run("Test should return the correct cluster", func(t *testing.T) {
		firstEpoch := 0
		secondEpoch := 180
		thirdEpoch := 360
		fourthEpoch := 540

		firstEpochLockIn := firstEpoch + 15
		secondEpochLockIn := secondEpoch + 12
		thirdEpochLockIn := thirdEpoch + 10

		firstEpochCommittee := generateCommittee(committeeLen)
		secondEpochCommittee := rotateCommittee(firstEpochCommittee, 10)
		thirdEpochCommittee := rotateCommittee(secondEpochCommittee, 20)

		self := firstEpochCommittee[0]
		firstEpochLatMap := generateLatMap(firstEpochCommittee)
		secondEpochLatMap := generateLatMap(secondEpochCommittee)
		thirdEpochLatMap := generateLatMap(thirdEpochCommittee)

		firstEpochLatMat := generateLatMat(firstEpochCommittee, firstEpochLatMap, self)
		secondEpochLatMat := generateLatMat(secondEpochCommittee, secondEpochLatMap, self)
		thirdEpochLatMat := generateLatMat(thirdEpochCommittee, thirdEpochLatMap, self)

		firstEpochTransitional, err := NewClusters(firstEpochCommittee, firstEpochLatMap, nil, self)
		require.NoError(t, err)

		firstEpochLocked, err := NewClusters(firstEpochCommittee, firstEpochLatMap, firstEpochLatMat, self)
		require.NoError(t, err)

		secondEpochTransitional, err := NewClusters(secondEpochCommittee, secondEpochLatMap, nil, self)
		require.NoError(t, err)
		secondEpochLocked, err := NewClusters(secondEpochCommittee, secondEpochLatMap, secondEpochLatMat, self)
		require.NoError(t, err)

		thirdEpochTransitional, err := NewClusters(thirdEpochCommittee, thirdEpochLatMap, nil, self)
		require.NoError(t, err)
		thirdEpochLocked, err := NewClusters(thirdEpochCommittee, thirdEpochLatMap, thirdEpochLatMat, self)
		require.NoError(t, err)

		clusterRotation := &ClusterRotation{
			previousEpochClusters: Clusters{},
			transitionalClusters:  firstEpochTransitional,
			latestEpochClusters:   Clusters{},
		}

		for i := 0; i < firstEpochLockIn; i++ {
			requireSameCluster(t, firstEpochTransitional, clusterRotation.GetClusters(uint64(i)))
		}
		clusterRotation.LockIn(uint64(firstEpochLockIn), firstEpochLocked)
		for i := 0; i <= secondEpoch; i++ {
			if i < firstEpochLockIn {
				requireSameCluster(t, firstEpochTransitional, clusterRotation.GetClusters(uint64(i)))
			} else {
				requireSameCluster(t, firstEpochLocked, clusterRotation.GetClusters(uint64(i)))
			}
		}

		clusterRotation.EpochStart(uint64(secondEpoch), secondEpochTransitional)
		for i := firstEpochLockIn; i < secondEpochLockIn; i++ {
			if i <= secondEpoch {
				requireSameCluster(t, firstEpochLocked, clusterRotation.GetClusters(uint64(i)))
			} else {
				requireSameCluster(t, secondEpochTransitional, clusterRotation.GetClusters(uint64(i)))
			}
		}
		clusterRotation.LockIn(uint64(secondEpochLockIn), secondEpochLocked)
		for i := secondEpoch + 1; i <= thirdEpoch; i++ {
			if i < secondEpochLockIn {
				requireSameCluster(t, secondEpochTransitional, clusterRotation.GetClusters(uint64(i)))
			} else {
				requireSameCluster(t, secondEpochLocked, clusterRotation.GetClusters(uint64(i)))
			}
		}
		clusterRotation.EpochStart(uint64(thirdEpoch), thirdEpochTransitional)
		for i := secondEpochLockIn; i < thirdEpochLockIn; i++ {
			if i <= thirdEpoch {
				requireSameCluster(t, secondEpochLocked, clusterRotation.GetClusters(uint64(i)))
			} else {
				requireSameCluster(t, thirdEpochTransitional, clusterRotation.GetClusters(uint64(i)))
			}
		}
		clusterRotation.LockIn(uint64(thirdEpochLockIn), thirdEpochLocked)
		for i := thirdEpoch + 1; i <= fourthEpoch; i++ {
			if i < thirdEpochLockIn {
				requireSameCluster(t, thirdEpochTransitional, clusterRotation.GetClusters(uint64(i)))
			} else {
				requireSameCluster(t, thirdEpochLocked, clusterRotation.GetClusters(uint64(i)))
			}
		}
	})
}

func requireSameCluster(t *testing.T, clusterA, clusterB Clusters) {
	require.Equal(t, len(clusterA.base), len(clusterB.base), "num clusters mismatch")

	for clusterId := 0; clusterId < len(clusterA.base); clusterId++ {
		require.Equal(t, len(clusterA.base[clusterId].Members), len(clusterB.base[clusterId].Members), "cluster members length mismatch")
		for j := 0; j < len(clusterA.base[clusterId].Members); j++ {
			require.Equal(t, clusterA.base[clusterId].Members[j].Addr, clusterB.base[clusterId].Members[j].Addr, "cluster member address mismatch")
			require.Equal(t, clusterA.base[clusterId].Members[j].Lat, clusterB.base[clusterId].Members[j].Lat, "cluster member latency mismatch")
		}
	}
}

func generateLatMap(committee []common.Address) map[common.Address]uint {
	latMap := make(map[common.Address]uint)
	for _, addr := range committee {
		latMap[addr] = uint(rand.Intn(256))
	}
	return latMap
}

func generateLatMat(committee []common.Address, latMap map[common.Address]uint, self common.Address) [][]uint8 {
	committeeLen := len(committee)
	latMat := make([][]uint8, committeeLen)
	for i := 0; i < committeeLen; i++ {
		latMat[i] = make([]uint8, committeeLen)
		for j := 0; j < i; j++ {
			if committee[i] == self {
				latMat[i][j] = uint8(latMap[committee[j]])
			} else {
				latMat[i][j] = uint8(rand.Intn(256))
				latMat[j][i] = latMat[i][j]
			}
		}
	}
	return latMat
}

func generateCommittee(size int) []common.Address {
	committee := make([]common.Address, size)
	for i := 0; i < size; i++ {
		committee[i] = testrand.Address()
	}
	return committee
}

func rotateCommittee(committee []common.Address, rotation int) []common.Address {
	rotated := make([]common.Address, len(committee))
	for i := 0; i < len(committee); i++ {
		if i >= len(committee)-rotation {
			rotated[i] = testrand.Address()
		} else {
			rotated[i] = committee[i+rotation]
		}
	}
	return rotated
}
