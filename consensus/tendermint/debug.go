package tendermint

import (
	"encoding/hex"

	"github.com/clearmatics/autonity/common"
)

func addr(a common.Address) string {
	return hex.EncodeToString(a[:3])
}
