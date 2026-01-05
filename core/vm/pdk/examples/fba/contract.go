package fba

import (
	"math/big"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core/vm"
	"github.com/autonity/autonity/core/vm/pdk/storage"
)

var (
	// todo
	clearingAddress = common.HexToAddress("0xFBACLeaRing")
)

type Side int

const (
	SideAsk = iota + 1
	SideBid
)

type BatchHandler struct {
}

func (bh *BatchHandler) ExecuteBatch(
	evm *vm.EVM,
	caller common.Address,
	st *storage.Storage,
	intents []OrderIntent,
) error {
	// implementation of batch execution logic
	bh.RunFBA(intents)
	return nil
}

func (bh *BatchHandler) ValidateIntents(intents []OrderIntent) error {
	// implementation of OrderIntent validation logic
	// mae check
	//
	return nil
}

func (bh *BatchHandler) RunFBA(intents []OrderIntent) {
	// implementation of FBA logic
	// filter order
	// verify using allowlist
	// since we trust our off chain engine, it can sign the whole batch and we verify once
	for _, intent := range intents {
		_ = intent.Quantity
	}
}

func (bh *BatchHandler) ComputeClearingPrice(intents []OrderIntent) *big.Int {
	// implementation of clearing price computation logic
	var asks, bids []OrderIntent
	for _, intent := range intents {
		if intent.Side == SideAsk {
			asks = append(asks, intent)
		} else if intent.Side == SideBid {
			bids = append(bids, intent)
		}
	}
	/* assume we have sorted asks and bids
		- find intersection point, price p, where best bid >= best ask
		- build demand and supply curves
	 		- cumulative quantity at each price level,
			- Demand = sum abs( bid quantities at price >= p)     1000 P  900
			- Supply = sum abs( ask quantities at price <= p)
		- allocate quantity ()
			-

	*/
	for i, j := 0, 0; i < len(asks) && j < len(bids); {
		ask := asks[i]
		bid := bids[j]
		if bid.Price.Cmp(ask.Price) >= 0 {
			// match found
			// determine matched quantity
			if bid.Quantity.Cmp(ask.Quantity) >= 0 {
				// bid can fully satisfy ask
				bid.Quantity.Sub(bid.Quantity, ask.Quantity)
				ask.Quantity.SetInt64(0)
				i++
			} else {
				// ask can fully satisfy bid
				ask.Quantity.Sub(ask.Quantity, bid.Quantity)
				bid.Quantity.SetInt64(0)
				j++
			}
			// record the trade at ask.Price
		} else {
			// no match
			i++
		}
	}

	return big.NewInt(0)
}
