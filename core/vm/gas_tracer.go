package vm

import (
	"encoding/hex"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/metrics"
	"github.com/autonity/autonity/params"
)

var (
	accountabilityFinalizeGas         = metrics.NewRegisteredBufferedGauge("autonity/accountability/finalize", nil, nil)
	omissionAccountabilityFinalizeGas = metrics.NewRegisteredBufferedGauge("autonity/omission_accountability/finalize", nil, nil)
	oracleFinalizeGas                 = metrics.NewRegisteredBufferedGauge("autonity/oracle/finalize", nil, nil)
	inflationCalculateDeltaGas        = metrics.NewRegisteredBufferedGauge("autonity/inflation/calculate_supply_delta", nil, nil)
	acuUpdateGas                      = metrics.NewRegisteredBufferedGauge("autonity/acu/update", nil, nil)
)

type Tracer map[string]struct{}

// TracerRegistry save the tracked functions' signature in 4 bytes as the method ID.
var TracerRegistry = map[common.Address]Tracer{
	// params.AutonityContractAddress:               {"4bb278f3": struct{}{}}, // finalize(), duplicated with legacy metrics.
	params.AccountabilityContractAddress:         {"6c9789b0": struct{}{}}, // finalize(bool)
	params.OmissionAccountabilityContractAddress: {"6c9789b0": struct{}{}}, // finalize(bool)
	params.OracleContractAddress:                 {"4bb278f3": struct{}{}}, // finalize()
	params.InflationControllerContractAddress:    {"92eff3cd": struct{}{}}, // calculateSupplyDelta(uint256,uint256,uint256,uint256)
	params.ACUContractAddress:                    {"a2e62045": struct{}{}}, // update()
}

func trace(contract common.Address, input []byte, initialGas, leftOverGas uint64) {
	if len(input) < 4 {
		return
	}

	methodID := input[:4]
	funcs, ok := TracerRegistry[contract]
	if !ok {
		return
	}

	id := hex.EncodeToString(methodID)
	if _, ok := funcs[id]; !ok {
		return
	}

	used := int64(initialGas - leftOverGas)
	switch contract {
	case params.AccountabilityContractAddress:
		accountabilityFinalizeGas.Add(used)
	case params.OmissionAccountabilityContractAddress:
		omissionAccountabilityFinalizeGas.Add(used)
	case params.OracleContractAddress:
		oracleFinalizeGas.Add(used)
	case params.InflationControllerContractAddress:
		inflationCalculateDeltaGas.Add(used)
	case params.ACUContractAddress:
		acuUpdateGas.Add(used)
	}
}
