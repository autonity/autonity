package trade_engine

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math/big"
	"reflect"
	"time"

	"github.com/holiman/uint256"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core/vm"
	"github.com/autonity/autonity/core/vm/precompile_sdk"
	"github.com/autonity/autonity/core/vm/precompile_sdk/storage"
	"github.com/autonity/autonity/crypto"
)

type TradingEngineContract struct {
	*precompile_sdk.BaseContract
	repo Repository
}

func NewTradingEngineContract(vm *vm.EVM, address common.Address) *TradingEngineContract {
	c := &TradingEngineContract{}
	c.BaseContract = precompile_sdk.SetupContract(c, reflect.TypeOf(TradingEngineState{}), address)

	st := storage.NewStorage(c.Address, vm.StateDB, c.Slots)
	c.repo = NewPrecompileRepository(st)
	// set order book
	pairHash := sha256.Sum256([]byte("NTN/USDC"))
	emptyBook := OrderBook{NextID: storage.NewUint256FromInt(1)}
	err := c.repo.SetBook(pairHash, emptyBook)
	if err != nil {
		panic("failed to set initial order book: " + err.Error())
	}
	return c
}

func (te *TradingEngineContract) SubmitOrder(evm *vm.EVM, caller common.Address, st *storage.Storage,
	pair string, side uint8, price, qty *big.Int) (common.Hash, error) {
	pairHash := sha256.Sum256([]byte(pair))
	ob, err := te.repo.GetBook(pairHash)
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
	err = te.repo.SetOrder(newID, newOrder)
	if err != nil {
		return common.Hash{}, err
	}

	err = te.repo.InsertToOrderBook(pairHash, newOrder)
	if err != nil {
		return common.Hash{}, err
	}

	if side == 0 {
		if len(ob.Bids) == 0 {
			ob.Bids = append(ob.Bids, Level{})
		}
		ob.Bids[0].TotalQty.Add(&ob.Bids[0].TotalQty.Int, &newOrder.Qty.Int)
	} else {
		if len(ob.Asks) == 0 {
			ob.Asks = append(ob.Asks, Level{})
		}
		ob.Asks[0].TotalQty.Add(&ob.Bids[0].TotalQty.Int, &newOrder.Qty.Int)
	}

	ob.NextID.Add(&ob.NextID.Int, &uint256.Int{1})
	err = te.repo.SetBook(pairHash, ob)
	return newID, err
}

func (te *TradingEngineContract) CancelOrder(evm *vm.EVM, caller common.Address, st *storage.Storage, pair string, orderID common.Hash) error {
	ord, err := te.repo.GetOrder(orderID)
	if err != nil {
		return err
	}
	if ord.User != caller || ord.Status != 0 {
		return fmt.Errorf("unauthorized: wrong user or already cancelled")
	}

	ord.Status = 2 // Cancelled
	if err := te.repo.SetOrder(orderID, ord); err != nil {
		return err
	}

	pairHash := sha256.Sum256([]byte(pair))
	// todo(optimize): pass order here
	if err := te.repo.RemoveFromOrderBook(pairHash, orderID); err != nil {
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
