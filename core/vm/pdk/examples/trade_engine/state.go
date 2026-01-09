package trade_engine

import (
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core/vm/pdk/storage"
)

type Order struct {
	ID        storage.Var[common.Hash]
	User      storage.Var[common.Address]
	Side      storage.Var[uint8]
	Price     storage.Var[storage.Uint256]
	Qty       storage.Var[storage.Uint256]
	Filled    storage.Var[storage.Uint256]
	Status    storage.Var[uint8]
	Timestamp storage.Var[int64]
}

type Level struct {
	Price    storage.Var[storage.Uint256]
	TotalQty storage.Var[storage.Uint256]
	OrderIDs storage.Slice[storage.Var[common.Hash]] // reference to all orders for this price
}

type OrderBook struct {
	Bids   storage.Slice[Level]
	Asks   storage.Slice[Level]
	NextID storage.Var[storage.Uint256] // simple sequence for orderIDs
}

type Trade struct {
	TakerID   storage.Var[common.Hash]
	MakerID   storage.Var[common.Hash]
	Price     storage.Var[storage.Uint256]
	Qty       storage.Var[storage.Uint256]
	Timestamp storage.Var[uint64]
}
