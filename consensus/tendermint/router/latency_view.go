package router

import (
	"github.com/autonity/autonity/common"
)

type latencyReports struct {
	latMat    [][]uint8
	committee []common.Address
}

func constructTransitional(oldLatMat [][]uint8, oldCommittee, newCommittee []common.Address) [][]uint8 {
	oldMapping := make(map[common.Address]map[common.Address]uint8)
	for i, addr := range oldCommittee {
		oldMapping[addr] = make(map[common.Address]uint8)
		for j, latency := range oldLatMat[i] {
			oldMapping[addr][oldCommittee[j]] = latency
		}
	}
	result := make([][]uint8, len(newCommittee))
	for i, addr := range newCommittee {
		result[i] = make([]uint8, len(newCommittee))
		for j, newAddr := range newCommittee {
			if i == j {
				result[i][j] = 0 // Self otherLatency is 0
				continue
			}
			if latency, ok := oldMapping[addr][newAddr]; ok {
				result[i][j] = latency
			} else if otherLatency, otherOk := oldMapping[newAddr][addr]; otherOk {
				result[i][j] = otherLatency // Default value if not found
			} else {
				result[i][j] = DefaultLatency // Default value if not found
			}
		}
	}
	return result
}

// Helper function to find the index of an address in a slice
func indexOf(slice []common.Address, addr common.Address) int {
	for i, a := range slice {
		if a == addr {
			return i
		}
	}
	return -1
}

func fillMissingLatencies(report latencyReports) [][]uint8 {
	// Create a new mapping to avoid modifying the original latMat
	result := make([][]uint8, len(report.committee))
	for i := range report.committee {
		result[i] = make([]uint8, len(report.committee))
		for j := range report.committee {
			if i == j {
				result[i][j] = 0 // Self latency is 0
				continue
			}
			if report.latMat[i][j] == 0 || report.latMat[i][j] == DefaultLatency {
				// try the other direction
				if report.latMat[j][i] != 0 && report.latMat[j][i] != DefaultLatency {
					result[i][j] = report.latMat[j][i]
				} else {
					result[i][j] = DefaultLatency // Default value if not found
				}
			}
		}
	}

	return result
}
