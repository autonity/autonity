package trade_engine

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math/big"
	"time"

	"github.com/holiman/uint256"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core/vm"
	"github.com/autonity/autonity/core/vm/pdk"
	"github.com/autonity/autonity/core/vm/pdk/storage"
	"github.com/autonity/autonity/crypto"
	"github.com/autonity/autonity/log"
)

var (
	ContractAddress = common.HexToAddress("0x23") // example address for precompiled contract
)

type TradingEngineContract struct {
	// storage layout
	Books        storage.Map[common.Hash, OrderBook] // per pair orderbook, e.g. NTN/USDC, ATN/USDC etc.
	Orders       storage.Map[common.Hash, Order]     // order to orderID mapping
	TradeHistory storage.Slice[Trade]
}

func SetupTradingEngineContract(vm *vm.EVM) *TradingEngineContract {
	c := &TradingEngineContract{}
	pdk.AddToPrecompiles(ContractAddress, c, vm, initialize)
	return c
}

func initialize(st *storage.Storage, bc *pdk.BaseContract) {
	pairHash := sha256.Sum256([]byte("NTN/USDC"))
	nextID := storage.NewUint256FromInt(1)

	te := bc.GetAppContract().(*TradingEngineContract)
	ob := te.Books.Get(pairHash)
	ob.NextID.Set(nextID)
	st.Commit()
}

func (te *TradingEngineContract) MatchOrders(evm *vm.EVM, caller common.Address, st *storage.Storage, pair string) error {
	// only autonity contract can call this function
	return nil
}

func (te *TradingEngineContract) SubmitOrder(evm *vm.EVM, caller common.Address, st *storage.Storage,
	pair string, side uint8, price, qty *big.Int) (common.Hash, error) {
	pairHash := sha256.Sum256([]byte(pair))
	ob := te.Books.Get(pairHash)

	ts := time.Now().Unix()
	newID := genOrderID(caller, ob.NextID.Get(), ts)
	newOrder := te.Orders.Get(newID)
	newOrder.ID.Set(newID)
	newOrder.Side.Set(side)
	newOrder.User.Set(caller)
	newOrder.Qty.Set(storage.NewUint256FromBig(qty))
	newOrder.Price.Set(storage.NewUint256FromBig(price))
	newOrder.Status.Set(0)
	newOrder.Timestamp.Set(ts)

	if side == 0 {
		ob.Bids.Append(func(l *Level) {
			l.Price.Set(storage.NewUint256FromBig(price))
			lq := l.TotalQty.Get()
			newQty := newOrder.Qty.Get()
			newQty.Add(&lq.Int, &newQty.Int)
			l.TotalQty.Set(newQty)

			l.OrderIDs.Append(func(idVar *storage.Var[common.Hash]) {
				idVar.Set(newOrder.ID.Get())
			})
		})
	} else {
		ob.Asks.Append(func(l *Level) {
			l.Price.Set(storage.NewUint256FromBig(price))
			lq := l.TotalQty.Get()
			newQty := newOrder.Qty.Get()
			newQty.Add(&lq.Int, &newQty.Int)
			l.TotalQty.Set(newQty)
			l.OrderIDs.Append(func(idVar *storage.Var[common.Hash]) {
				idVar.Set(newOrder.ID.Get())
			})
		})
	}

	nextId := ob.NextID.Get()
	nextId.Add(&nextId.Int, &uint256.Int{1})
	ob.NextID.Set(nextId)

	book := te.Books.Get(pairHash)
	log.Info("updated book", "value ", book.NextID.Get())
	return newID, nil
}

func (te *TradingEngineContract) CancelOrder(evm *vm.EVM, caller common.Address, st *storage.Storage, pair string, orderID common.Hash) error {
	ord := te.Orders.Get(orderID)
	if ord.User.Get() != caller || ord.Status.Get() != 0 {
		return fmt.Errorf("unauthorized: wrong user or already cancelled")
	}

	ord.Status.Set(2) // Cancelled

	// todo remove from order book
	return nil
}

// stub function to generate order ID
func genOrderID(user common.Address, next storage.Uint256, ts int64) common.Hash {
	bytes := append(user.Bytes(), append(next.Bytes(), make([]byte, 8)...)...)
	binary.BigEndian.PutUint64(bytes[len(bytes)-8:], uint64(ts))
	return common.BytesToHash(crypto.Keccak256(bytes))
}
