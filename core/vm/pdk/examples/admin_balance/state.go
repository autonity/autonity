package admin_balance

import (
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core/vm/pdk/storage"
)

// AdminBalanceState defines state fields.
type AdminBalanceState struct {
	Admin   common.Address
	Balance storage.Uint256
}
