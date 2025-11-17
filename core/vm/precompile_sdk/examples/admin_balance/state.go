package admin_balance

import (
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core/vm/precompile_sdk/storage"
)

// AdminBalanceState defines state fields.
type AdminBalanceState struct {
	Admin   common.Address
	Balance storage.Uint256
}
