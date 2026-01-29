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

type Side int

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

type Intent struct {
	MarginAccountID common.Address
	IntentAccountID common.Address
	Hash            common.Hash
	Data            IntentData
	Signature       []byte
}

// FBA Specific Wrapper
type Order struct {
	Intent            Intent
	IntervalID        *big.Int
	RemainingQuantity *big.Int
}

// Batch Input
type Batch struct {
	BatchID         *big.Int
	ProductID       common.Hash
	StartIntervalID *big.Int
	EndIntervalID   *big.Int
	Bids            []Order // Sorted: Price DESC, Interval ASC
	Asks            []Order // Sorted: Price ASC, Interval ASC
}

// Output Results
type FillResult struct {
	IntentHash        common.Hash
	FillQuantity      *big.Int
	RemainingQuantity *big.Int
}

// Settlement struct for MAE Check (Matches IMarginAccount)
type Settlement struct {
	PositionId common.Hash
	Quantity   *big.Int // Signed int256
	Price      *big.Int
}
