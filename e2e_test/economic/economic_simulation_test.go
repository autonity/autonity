package economic

import (
	"github.com/autonity/autonity/log"
	"testing"
)

// TestDefaultConfig tests the economic model with different block-filling scenarios.
func TestDefaultConfig(t *testing.T) {
	tests := []struct {
		name   string
		blocks uint64
		packer TXNPacker
	}{
		{"dust filled blocks", 600, &dustFiller{params: &defSysParams}},
		{"1/3 filled blocks", 600, &oneThirdPacker{params: &defSysParams}},
		{"half filled blocks", 600, &halfPacker{params: &defSysParams}},
		{"2/3 filled blocks", 600, &twoThirdPacker{params: &defSysParams}},
		{"fully filled blocks", 600, &fullPacker{params: &defSysParams}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runSimulation(t, tt.name, tt.blocks, &defSysParams, tt.packer)
		})
	}
}

// runSimulation runs a simulation and logs the results.
func runSimulation(t *testing.T, name string, blocks uint64, params *systemParams, packer TXNPacker) {
	log.Info("Running test", "name", name, "blocks", blocks, "params", params)
	sim := newSimulator(blocks, params, packer)
	sim.start()
}

// todo, address the cost of spam.

// todo, add helpers to render data in a diagram.
