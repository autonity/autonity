package clustering

import (
	"github.com/autonity/autonity/consensus/tendermint/router/ping"
	"math/rand"
	"time"
)

// this type is used for local e2e testing.
type simulatedPinger struct{}

func NewSimulatedPinger() ping.Pinger {
	return &simulatedPinger{}
}

func (*simulatedPinger) Ping(_ ping.Target, resultCh chan<- time.Duration) {
	// Generate random latency between 10ms and 500ms
	latency := time.Duration(10+rand.Int63n(491)) * time.Millisecond
	resultCh <- latency
}
