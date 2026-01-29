# PDK Cookbook

Common patterns and snippets for building precompiles.

## 1. Simple Token (ERC-20 style)

Manages a balance map.

```go
type TokenContract struct {
    Balances storage.Map[common.Address, storage.Uint256]
    TotalSupply storage.Var[storage.Uint256]
}

func (t *TokenContract) Transfer(evm *vm.EVM, caller common.Address, st *storage.Storage, to common.Address, amount *big.Int) error {
    // 1. Check Balance
    senderBal := t.Balances.Get(caller).Get()
    amt := storage.NewUint256FromBig(amount)
    
    if senderBal.Lt(&amt) {
        return fmt.Errorf("insufficient balance")
    }

    // 2. Debit Sender
    newSenderBal := senderBal.Sub(&senderBal, &amt)
    t.Balances.Get(caller).Set(newSenderBal)

    // 3. Credit Recipient
    receiverBal := t.Balances.Get(to).Get()
    newReceiverBal := receiverBal.Add(&receiverBal, &amt)
    t.Balances.Get(to).Set(newReceiverBal)

    return nil
}
```

## 2. Order Book (Limit Orders)

Shows usage of `Slice` inside `Map` (One-to-Many).

```go
type Order struct {
    Price storage.Var[uint64]
    Qty   storage.Var[uint64]
}

type Market struct {
    // Map: Price -> List of Orders at that price
    Bids storage.Map[uint64, storage.Slice[Order]]
}

func (m *Market) AddBid(evm *vm.EVM, caller common.Address, st *storage.Storage, price uint64, qty uint64) error {
    orders := m.Bids.Get(price)
    
    orders.Append(func(o *Order) {
        o.Price.Set(price)
        o.Qty.Set(qty)
    })
    return nil
}

func (m *Market) ClearPriceLevel(evm *vm.EVM, caller common.Address, st *storage.Storage, price uint64) error {
    // Removes all orders at this price level
    // Recursively cleans up the Slice and all Orders
    m.Bids.Delete(price) 
    return nil
}
```

## 3. Admin Control

Restricting access to methods.

```go
type AdminContract struct {
    Admin storage.Var[common.Address]
}

func (a *AdminContract) CriticalOp(evm *vm.EVM, caller common.Address, st *storage.Storage) error {
    if caller != a.Admin.Get() {
        return fmt.Errorf("unauthorized")
    }
    // Do critical stuff
    return nil
}
```
