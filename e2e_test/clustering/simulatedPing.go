package clustering

import (
	"context"
	"math/rand"
	"time"

	"github.com/autonity/autonity/consensus/tendermint/router/ping"
)

// this type is used for local e2e testing.
type simulatedPinger struct{}

func NewSimulatedPinger() ping.Pinger {
	return &simulatedPinger{}
}

func (*simulatedPinger) Ping(_ context.Context, _ ping.Target) ping.Result {
	// Generate random latency between 10ms and 500ms
	latency := time.Duration(10+rand.Int63n(491)) * time.Millisecond
	return ping.Result{Latency: latency}
}
