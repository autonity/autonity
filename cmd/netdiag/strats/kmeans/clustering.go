package kmeans

import (
	"time"

	"github.com/parallelo-ai/kmeans"
)

func AssignClusters(latencyMat [][]time.Duration, k int) ([][]int, error) {
	nodes := fromLatencyMatrix(latencyMat)

	var obs []kmeans.Observation
	for _, n := range nodes {
		obs = append(obs, n)
	}

	km := kmeans.New()
	cstrs, err := km.Partition(obs, k, 12345)
	if err != nil {
		return nil, err
	}

	clusterSlice := make([][]int, k)
	for i, c := range cstrs {
		clusterSlice[i] = make([]int, len(c.Observations))
		for j, o := range c.Observations {
			clusterSlice[i][j] = o.(*Node).Id
		}
	}
	return clusterSlice, nil
}
