# Storage Engine Deep Dive

The PDK Storage Engine acts as an Object-Relational Mapper (ORM) for the EVM's key-value store. It maps Go types to 256-bit storage slots automatically.

## Core Wrappers

### 1. `storage.Var[T]`
A single variable. Supports primitive types (`uint64`, `bool`, `common.Address`) and fixed-size arrays.

*   **Packing:** The PDK packs variables into slots if they fit.
    *   `Var[uint64]` (8 bytes) + `Var[address]` (20 bytes) = 28 bytes. These will share **Slot 0**.
*   **Usage:**
    ```go
    v.Get()      // Returns T
    v.Set(val)   // Writes val
    v.Clear()    // Resets to zero (Triggering Gas Refund)
    ```

### 2. `storage.Map[K, V]`
A cryptographic hash map.
*   **Slot calculation:** `keccak256(key . slot)`
*   **Usage:**
    ```go
    m.Get(key)   // Returns *V (Pointer to value wrapper)
    m.Delete(key) // Recursively clears storage
    ```
    > **Note:** `Get` never returns nil. It returns a bound wrapper pointing to the slot. If the slot is empty, `Get().Get()` returns the zero value.

### 3. `storage.Slice[T]`
A dynamic array.
*   **Storage Layout:**
    *   **Header Slot:** Stores the length (`uint64`).
    *   **Data Slots:** `keccak256(HeaderSlot) + index`.
*   **Usage:**
    ```go
    s.Len()
    s.Get(i)     // Returns *T
    s.Append(func(v *T) { ... })
    s.Clear()    // Clears all elements and resets length
    ```

## 🔐 Cleanup & Gas Refunds

One of the most critical features of the PDK is its handling of state cleanup.

### The "Orphan Data" Problem
In raw Solidity/EVM, if you delete a mapping entry that contains a dynamic array, you must manually delete the array elements first. Failing to do so leaves "orphaned" data in storage that costs gas but isn't accessible.

### The PDK Solution: `RecursiveClear`
The PDK implements a recursive cleanup strategy.

*   **`Map.Delete(key)`**: 
    1.  Calculates the slot for the key.
    2.  Instantiates a temporary binding for the value type `V`.
    3.  Calls `RecursiveClear` on the value.
    4.  If `V` is a struct, it iterates all fields.
    5.  If a field is a `Slice` or another `Map`, it recurses.
    6.  Finally, it writes zero to the slots.

**Example:**
```go
// Nested structure
type User struct {
    Name    storage.Var[[]byte]      // Dynamic bytes
    Orders  storage.Slice[uint64]    // Dynamic array
}

// Map
Users storage.Map[address, User]

// Deletion
func RemoveUser(addr common.Address) {
    // This single call will:
    // 1. Clear the Name length header AND all name data chunks.
    // 2. Clear the Orders length header AND all order items.
    // 3. Trigger EVM gas refunds for every cleared slot.
    c.Users.Delete(addr)
}
```

## Storage Layout Visualization

```mermaid
graph LR
    subgraph "Contract Base"
       Slot0[Slot 0]
       Slot1[Slot 1]
    end

    subgraph "Packed Vars (Slot 0)"
       ID[Var: ID (uint64)]
       Active[Var: Active (bool)]
    end
    Slot0 --> ID
    Slot0 --> Active

    subgraph "Map (Slot 1)"
       MapHead[Map Base Slot]
       MapKey[Key] -->|Hash| MapSlot[Data Slot]
    end
    Slot1 --> MapHead
```
