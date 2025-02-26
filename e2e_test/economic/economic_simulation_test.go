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
		params *systemParams
		packer TXNPacker
	}{
		// 0 - 5 run with default params
		{"dust filled blocks", 600, &defSysParams, &dustPacker{}},
		{"1/3 filled blocks", 600, &defSysParams, &oneThirdPacker{}},
		{"half filled blocks", 600, &defSysParams, &halfPacker{}},
		{"2/3 filled blocks", 600, &defSysParams, &twoThirdPacker{}},
		{"fully filled blocks", 600, &defSysParams, &fullPacker{}},
		{"dynamic filled blocks", 3600, &defSysParams, &dynamicPacker{increasingInterval: 100, decreasingInterval: 500}},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runSimulation(i, tt.name, tt.blocks, tt.params, tt.packer)
		})
	}
}

// runSimulation runs a simulation and logs the results.
func runSimulation(index int, name string, blocks uint64, params *systemParams, packer TXNPacker) {
	log.Info("Running test", "name", name, "blocks", blocks, "params", params)
	sim := newSimulator(name, blocks, params, packer)
	sim.start()
	sim.assembleData(index)
}
