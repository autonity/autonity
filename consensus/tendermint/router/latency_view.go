package router

import (
	"github.com/autonity/autonity/common"
)

type latencyReports struct {
	latMat    [][]uint8
	committee []common.Address
}

func getForCommittee(lr latencyReports, committee []common.Address) [][]uint8 {
	latestReport := lr
	latestLatMat := fillMissingLatencies(latestReport)

	// Prepare the result map
	result := make([][]uint8, len(committee))

	for i, addr := range committee {
		result[i] = make([]uint8, len(committee))
		if idx := indexOf(latestReport.committee, addr); idx != -1 {
			for j, addr2 := range committee {
				if idx2 := indexOf(latestReport.committee, addr2); idx2 != -1 {
					result[i][j] = latestLatMat[idx][idx2]
				} else {
					result[i][j] = DefaultLatency
				}
			}
		} else {
			for j := range committee {
				result[i][j] = DefaultLatency
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
