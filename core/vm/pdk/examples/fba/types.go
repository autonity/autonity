package fba

import (
	"math/big"

	"github.com/autonity/autonity/common"
)

// Enums
const (
	SideBid = 0
	SideAsk = 1
)

// Side represents the side of an order (Bid/Ask).
type Side int

// IntentData holds the core data of a trading intent.
type IntentData struct {
	Nonce             *big.Int
	TradingProtocolID common.Address
	ProductID         common.Hash
	LimitPrice        *big.Int
	Quantity          *big.Int
	MaxTradingFeeRate *big.Int
	GoodUntil         *big.Int
	Side              uint8 // 0: BID, 1: ASK
}

// Intent represents a user intent to trade.
type Intent struct {
	MarginAccountID common.Address
	IntentAccountID common.Address
	Hash            common.Hash
	Data            IntentData
	Signature       []byte
}

// Order represents an order in the order book.
type Order struct {
	Intent            Intent
	IntervalID        *big.Int
	RemainingQuantity *big.Int
}

// Batch represents a batch of orders to be processed.
type Batch struct {
	BatchID         *big.Int
	ProductID       common.Hash
	StartIntervalID *big.Int
	EndIntervalID   *big.Int
	Bids            []Order // Sorted: Price DESC, Interval ASC
	Asks            []Order // Sorted: Price ASC, Interval ASC
}

// FillResult represents the result of filling an order.
type FillResult struct {
	IntentHash        common.Hash
	FillQuantity      *big.Int
	RemainingQuantity *big.Int
}

// Settlement represents the settlement details for a position.
type Settlement struct {
	PositionID common.Hash
	Quantity   *big.Int // Signed int256
	Price      *big.Int
}
