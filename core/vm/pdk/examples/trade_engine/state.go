package trade_engine

import (
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core/vm/pdk/storage"
)

type Order struct {
	ID        common.Hash
	User      common.Address
	Side      uint8
	Price     storage.Uint256
	Qty       storage.Uint256
	Filled    storage.Uint256
	Status    uint8
	Timestamp int64
}

type Level struct {
	Price    storage.Uint256
	TotalQty storage.Uint256
	OrderIDs []common.Hash // reference to all orders for this price
}

type OrderBook struct {
	Bids   []Level
	Asks   []Level
	NextID storage.Uint256 // simple sequence for orderIDs
}

type Trade struct {
	TakerID   common.Hash
	MakerID   common.Hash
	Price     storage.Uint256
	Qty       storage.Uint256
	Timestamp uint64
}
