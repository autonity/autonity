package core

import (
	"fmt"
	"testing"
)

func TestLatencyPing(t *testing.T) {
	results := PingFixedNTP()
	fmt.Printf("Pinged servers: %d\n", len(results))

	for i, result := range results {
		if result.AvgRtt == 0 {
			fmt.Printf("Expected non-zero RTT, got 0, IP: %s INDEX: %d\n", NtpServers[i], i)
			t.Errorf("Expected non-zero RTT, got 0")
		}
	}

}
