package trade_engine

import "github.com/autonity/autonity/common"

type Order struct {
	ID uint64
}

type OrderBook struct {
}

type Trade struct {
}

type TradingEngineState struct {
	Books map[common.Hash]OrderBook
}
