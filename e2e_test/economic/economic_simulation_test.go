package economic

import (
	"github.com/autonity/autonity/log"
	"github.com/shopspring/decimal"
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
		{"dust filled blocks", 1801, &defSysParams, &RatedPacker{ratio: decimal.NewFromFloatWithExponent(0.01, -2)}},
		{"1/10 filled blocks", 1801, &defSysParams, &RatedPacker{ratio: decimal.NewFromFloatWithExponent(0.1, -2)}},
		{"1/5 filled blocks", 1801, &defSysParams, &RatedPacker{ratio: decimal.NewFromFloatWithExponent(0.2, -2)}},
		{"1/4 filled blocks", 1801, &defSysParams, &RatedPacker{ratio: decimal.NewFromFloatWithExponent(0.25, -2)}},
		{"1/3 filled blocks", 1801, &defSysParams, &RatedPacker{ratio: decimal.NewFromFloatWithExponent(0.33, -2)}},
		{"1/2 filled blocks", 1801, &defSysParams, &RatedPacker{ratio: decimal.NewFromFloatWithExponent(0.5, -2)}},
		{"3/5 filled blocks", 1801, &defSysParams, &RatedPacker{ratio: decimal.NewFromFloatWithExponent(0.66, -2)}},
		{"3/4 filled blocks", 1801, &defSysParams, &RatedPacker{ratio: decimal.NewFromFloatWithExponent(0.75, -2)}},
		{"full filled blocks", 1801, &defSysParams, &RatedPacker{ratio: decimal.NewFromFloatWithExponent(1.0, -2)}},
		{"dynamic filled blocks", 1801, &defSysParams, &dynamicPacker{increasingInterval: 100, decreasingInterval: 500}},
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
