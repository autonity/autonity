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
	Books        map[common.Hash]OrderBook // per pair orderbook, e.g. NTN/USDC, ATN/USDC etc.
	Orders       map[common.Hash]Order     // order to orderID mapping
	TradeHistory []Trade                   // todo:
}

func SetupTradingEngineContract(vm *vm.EVM) *TradingEngineContract {
	// todo: define address
	// address := common.HexToAddress("0x23")
	c := &TradingEngineContract{}
	pdk.AddToPrecompiles(ContractAddress, c, vm, initialize)
	return c
}

func initialize(st *storage.Storage) {
	// base contract assignment in add to precompiles
	repo := NewPrecompileRepository(st)
	// set order book
	pairHash := sha256.Sum256([]byte("NTN/USDC"))
	emptyBook := OrderBook{NextID: storage.NewUint256FromInt(1)}
	err := repo.SetBook(pairHash, emptyBook)
	if err != nil {
		panic("failed to set initial order book: " + err.Error())
	}
	st.Commit()
}

func (te *TradingEngineContract) MatchOrders(evm *vm.EVM, caller common.Address, st *storage.Storage, pair string) error {
	// only autonity contract can call this function
	return nil
}

func (te *TradingEngineContract) SubmitOrder(evm *vm.EVM, caller common.Address, st *storage.Storage,
	pair string, side uint8, price, qty *big.Int) (common.Hash, error) {
	pairHash := sha256.Sum256([]byte(pair))
	repo := NewPrecompileRepository(st)
	ob, err := repo.GetBook(pairHash)
	if err != nil {
		return common.Hash{}, err
	}

	ts := time.Now().Unix()
	newID := genOrderID(caller, &ob.NextID, ts)
	newOrder := Order{
		ID:        newID,
		User:      caller,
		Side:      side,
		Price:     storage.NewUint256FromBig(price),
		Qty:       storage.NewUint256FromBig(qty),
		Status:    0, // open
		Timestamp: ts,
	}
	err = repo.SetOrder(newID, newOrder)
	if err != nil {
		return common.Hash{}, err
	}
	updatedOrder, err := repo.GetOrder(newID)
	if err != nil {
		return common.Hash{}, err
	}
	log.Info("updated order", "value ", updatedOrder)

	if side == 0 {
		if len(ob.Bids) == 0 {
			ob.Bids = append(ob.Bids, Level{})
		}
		ob.Bids[0].Price = storage.NewUint256FromBig(price)
		ob.Bids[0].TotalQty.Add(&ob.Bids[0].TotalQty.Int, &newOrder.Qty.Int)
		ob.Bids[0].OrderIDs = append(ob.Bids[0].OrderIDs, newOrder.ID)
	} else {
		if len(ob.Asks) == 0 {
			ob.Asks = append(ob.Asks, Level{})
		}
		ob.Asks[0].Price = storage.NewUint256FromBig(price)
		ob.Asks[0].TotalQty.Add(&ob.Bids[0].TotalQty.Int, &newOrder.Qty.Int)
		ob.Asks[0].OrderIDs = append(ob.Asks[0].OrderIDs, newOrder.ID)
	}

	ob.NextID.Add(&ob.NextID.Int, &uint256.Int{1})
	err = repo.SetBook(pairHash, ob)
	book, err := repo.GetBook(pairHash)
	if err != nil {
		return [32]byte{}, err
	}
	log.Info("updated book", "value ", book)
	return newID, err
}

func (te *TradingEngineContract) CancelOrder(evm *vm.EVM, caller common.Address, st *storage.Storage, pair string, orderID common.Hash) error {
	repo := NewPrecompileRepository(st)
	ord, err := repo.GetOrder(orderID)
	if err != nil {
		return err
	}
	if ord.User != caller || ord.Status != 0 {
		return fmt.Errorf("unauthorized: wrong user or already cancelled")
	}

	ord.Status = 2 // Cancelled
	if err := repo.SetOrder(orderID, ord); err != nil {
		return err
	}

	pairHash := sha256.Sum256([]byte(pair))
	// todo(optimize): pass order here
	if err := repo.RemoveFromOrderBook(pairHash, orderID); err != nil {
		return err
	}
	return nil
}

// stub function to generate order ID
func genOrderID(user common.Address, next *storage.Uint256, ts int64) common.Hash {
	bytes := append(user.Bytes(), append(next.Bytes(), make([]byte, 8)...)...)
	binary.BigEndian.PutUint64(bytes[len(bytes)-8:], uint64(ts))
	return common.BytesToHash(crypto.Keccak256(bytes))
}
