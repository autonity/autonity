# Autonity PDK (Precompile Development Kit)

> **Type-Safe EVM Precompiles in Pure Go.**

The Autonity PDK is a framework for building stateful, high-performance EVM precompiled contracts using idiomatic Go. It eliminates the need for manual ABI decoding and raw storage slot manipulation, providing a robust "ORM-like" layer for EVM state.

## 📚 Documentation

*   [**Architecture Guide**](ARCHITECTURE.md) - Learn how the Dispatcher, Storage, and Gas models work together.
*   [**Storage Engine Deep Dive**](STORAGE.md) - Master `Var`, `Map`, `Slice`, and data layout.
*   [**Cookbook & Examples**](COOKBOOK.md) - Copy-pasteable patterns for common use cases.

## 🚀 Quick Start

Define your contract as a standard Go struct. The PDK handles storage mapping, ABI dispatch, and gas metering.

```go
package main

import (
    "github.com/autonity/autonity/common"
    "github.com/autonity/autonity/core/vm"
    "github.com/autonity/autonity/core/vm/pdk"
    "github.com/autonity/autonity/core/vm/pdk/storage"
    "github.com/autonity/autonity/core/vm/pdk/gas"
)

// 1. Define State
type CounterContract struct {
    Count  storage.Var[uint64]
    Owners storage.Map[common.Address, bool]
}

// 2. Define Methods
func (c *CounterContract) Increment(evm *vm.EVM, caller common.Address, st *storage.Storage) error {
    current := c.Count.Get()
    c.Count.Set(current + 1)
    c.Owners.Get(caller).Set(true)
    return nil
}

func (c *CounterContract) GetCount(evm *vm.EVM, caller common.Address, st *storage.Storage) (uint64, error) {
    return c.Count.Get(), nil
}

// 3. Register
func Register(vm *vm.EVM) {
    contract := &CounterContract{}
    gasConfig := gas.NewConfig(21000) // Base cost
    
    // Register as a precompile (Address 0x...FF)
    pdk.AddToPrecompiles(
        common.HexToAddress("0x00000000000000000000000000000000000000FF"), 
        contract, 
        vm, 
        nil, 
        gasConfig,
    )
}
```

## 📦 Installation

Since the PDK is part of the `autonity/core` codebase, simply import it:

```go
import "github.com/autonity/autonity/core/vm/pdk"
```

## ⚠️ Requirements

*   **Go 1.18+** (Generics support required)
*   **Autonity Core** dependency

---
*Note: Diagrams in documentation require a Mermaid-compatible viewer.*