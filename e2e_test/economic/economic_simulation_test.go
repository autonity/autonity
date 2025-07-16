package economic

import (
	"github.com/autonity/autonity/log"
	"github.com/shopspring/decimal"
	"testing"
)

// TestDefaultConfig tests the economic model with different block-filling scenarios.
func TestDefaultConfig(t *testing.T) {
	//tenMins := uint64(600)
	thirtyMins := uint64(1800)
	tests := []struct {
		name   string
		blocks uint64
		params *systemParams
		packer TXNPacker
	}{
		// under block gas target tests
		{"1% filled blocks", 3600 * 24 * 30 * 12, &defSysParams, &RatedPacker{ratio: decimal.NewFromFloatWithExponent(0.01, -2)}},
		{"5% filled blocks", 3600 * 24 * 30, &defSysParams, &RatedPacker{ratio: decimal.NewFromFloatWithExponent(0.05, -2)}},

		{"1/10 filled blocks", 3600, &defSysParams, &RatedPacker{ratio: decimal.NewFromFloatWithExponent(0.1, -2)}},
		{"1/5 filled blocks", 3600 * 24 * 7, &defSysParams, &RatedPacker{ratio: decimal.NewFromFloatWithExponent(0.2, -2)}},
		{"1/4 filled blocks", 3600 * 24 * 7, &defSysParams, &RatedPacker{ratio: decimal.NewFromFloatWithExponent(0.25, -2)}},
		{"1/3 filled blocks", 3600 * 24 * 7, &defSysParams, &RatedPacker{ratio: decimal.NewFromFloatWithExponent(0.33, -2)}},
		{"1/2 filled blocks", 3600 * 24 * 7, &defSysParams, &RatedPacker{ratio: decimal.NewFromFloatWithExponent(0.5, -2)}},
		// exceeding block gas target tests, the baseFee adjustment will start for below tests, the TXN fee will increase fast.

		{"over target 10%", thirtyMins, &defSysParams, &RatedPacker{ratio: decimal.NewFromFloatWithExponent(0.55, -2)}},

		{"over target 20%", thirtyMins, &defSysParams, &RatedPacker{ratio: decimal.NewFromFloatWithExponent(0.60, -2)}},

		{"over target 30%", thirtyMins, &defSysParams, &RatedPacker{ratio: decimal.NewFromFloatWithExponent(0.65, -2)}},

		{"over target 40%", thirtyMins, &defSysParams, &RatedPacker{ratio: decimal.NewFromFloatWithExponent(0.70, -2)}},

		{"over target 50%", thirtyMins, &defSysParams, &RatedPacker{ratio: decimal.NewFromFloatWithExponent(0.75, -2)}},

		{"over target 100%", thirtyMins, &defSysParams, &RatedPacker{ratio: decimal.NewFromFloatWithExponent(1.0, -2)}},

		// dynamic traffic
		//{"dynamic filled blocks", 181, &defSysParams, &dynamicPacker{increasingInterval: 50, decreasingInterval: 130}},
	}

	for i, tt := range tests {
		if i != 0 {
			continue
		}

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
