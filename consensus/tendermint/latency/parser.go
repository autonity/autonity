package latency

import (
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/log"
	"github.com/autonity/autonity/params"
	"github.com/autonity/autonity/params/generated"
)

type report struct {
	reporter  common.Address
	latencies []uint8
	height    uint64
}

func parseReports(block *types.Block) ([]report, error) {
	var reports []report
	for _, tx := range block.Transactions() {
		if tx.To() == nil || *tx.To() != params.LatencyContractAddress {
			continue
		}
		reportMethod, err := generated.LatencyAbi.MethodById(tx.Data())
		if err != nil {
			log.Error("LatencyReport: error fetching method by ID", "error", err, "block", block.NumberU64())
			continue
		}
		reportData, err := reportMethod.Inputs.Unpack(tx.Data()[4:])
		if err != nil {
			log.Error("LatencyReport: unable to unpack vote method ", "error", err)
			continue
		}
		signer, err := types.NewLondonSigner(tx.ChainId()).Sender(tx)
		if err != nil {
			log.Error("LatencyReport: unable to get signer", "error", err)
			continue
		}
		latencies := reportData[0].([]uint8)
		r := report{
			reporter:  signer,
			latencies: latencies,
			height:    block.NumberU64(),
		}
		reports = append(reports, r)
	}
	return reports, nil
}
