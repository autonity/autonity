package tendermint

import (
	"encoding/hex"
	"fmt"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/core/types"
)

func addr(a common.Address) string {
	return hex.EncodeToString(a[2:6])
}

func bid(b *types.Block) string {
	return fmt.Sprintf("hash: %v, number: %v", b.Hash().String()[2:8], b.Number().String())
}

type debugLog struct {
	prefix []interface{}
}

func newDebugLog(prefix ...interface{}) *debugLog {
	return &debugLog{
		prefix: prefix,
	}
}

func (d *debugLog) print(info ...interface{}) {
	// log := append(d.prefix, info...)
	// fmt.Printf("%v %v", time.Now().Format(time.RFC3339Nano), fmt.Sprintln(log...))
}
