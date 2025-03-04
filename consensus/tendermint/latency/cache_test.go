package latency

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/internal/testrand"
)

func TestLatencyCache(t *testing.T) {
	t.Run("Test cache should insert matrix line", func(t *testing.T) {
		cache := newLatencyCache()
		validators := []common.Address{
			testrand.Address(),
			testrand.Address(),
			testrand.Address(),
		}

		validator1LatencyVec := []uint8{0, 2, 3}
		cache.insertMatrixLine(validators[0], validators, validator1LatencyVec)

		mat := cache.readMatrix(validators)

		// own latency should be set
		require.Equal(t, uint8(0), mat[validators[0]][0])
		require.Equal(t, uint8(2), mat[validators[0]][1])
		require.Equal(t, uint8(3), mat[validators[0]][2])
	})

	t.Run("Test should return missing measurements", func(t *testing.T) {
		cache := newLatencyCache()
		oldCommittee := []common.Address{
			testrand.Address(),
			testrand.Address(),
			testrand.Address(),
		}

		measurements := map[common.Address]uint8{
			oldCommittee[0]: 0,
			oldCommittee[1]: 1,
			oldCommittee[2]: 2,
		}
		cache.insertMeasurements(oldCommittee, measurements)

		newCommittee := append(oldCommittee, testrand.Address())

		missing := cache.missingMeasurements(newCommittee)
		require.Len(t, missing, 1)
		require.Equal(t, newCommittee[3], missing[0])
	})
}

func TestClusterCache(t *testing.T) {
	t.Run("Test cache should return cluster containing address", func(t *testing.T) {
		cache := newClusterCache()
		cluster := [][]common.Address{
			{testrand.Address(), testrand.Address()},
			{testrand.Address()},
			{testrand.Address()},
		}

		cache.insertClustering(1, &Clusters{cluster, nil})
		c, ok := cache.clustersAt(2)
		require.True(t, ok)

		require.Equal(t, c.clusterContaining(cluster[2][0]), 2)
		require.Equal(t, c.clusterContaining(cluster[0][0]), 0)
	})

	t.Run("Test cache race condition", func(t *testing.T) {
		wg := sync.WaitGroup{}
		cache := newClusterCache()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 1000; i++ {
				select {
				case <-ctx.Done():
					t.Fail()
				default:
					cluster := [][]common.Address{
						{testrand.Address(), testrand.Address()},
						{testrand.Address()},
						{testrand.Address()},
					}
					cache.insertClustering(uint64(i), &Clusters{cluster, nil})
				}
			}
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 1000; i++ {
				select {
				case <-ctx.Done():
					t.Fail()
				default:
					cache.clustersAt(uint64(i))
				}
			}
		}()

		wg.Wait()
	})
}
