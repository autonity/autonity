package badgerdb

import (
	"errors"
	"fmt"
	"github.com/autonity/autonity/ethdb"
	"github.com/dgraph-io/badger/v4"
)

type Database struct {
	dir string // directory name for database.
	db  *badger.DB
}

func New(dir string) (*Database, error) {
	db, err := badger.Open(badger.DefaultOptions(dir))
	if err != nil {
		return nil, err
	}

	return &Database{
		dir: dir,
		db:  db,
	}, nil
}

func (db *Database) Close() error {
	return db.db.Close()
}

// Has checks if the key exists in the database.
func (db *Database) Has(key []byte) (bool, error) {
	err := db.db.View(func(txn *badger.Txn) error {
		_, err := txn.Get(key)
		return err
	})

	if errors.Is(badger.ErrKeyNotFound, err) {
		return false, nil // Key does not exist
	} else if err != nil {
		return false, err // An error occurred
	}

	return true, nil // Key exists
}

// Get retrieves the given key if it's present in the key-value store.
func (db *Database) Get(key []byte) ([]byte, error) {
	var value []byte

	err := db.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err // Return the error if the key is not found or other error
		}

		// Get the value from the item
		return item.Value(func(val []byte) error {
			value = append([]byte{}, val...) // Copy the value to avoid reference issues
			return nil
		})
	})

	if err != nil {
		if errors.Is(err, badger.ErrKeyNotFound) {
			return nil, fmt.Errorf("key not found: %s", key)
		}
		return nil, err // Return any other error encountered
	}

	return value, nil // Return the retrieved value
}

// Put inserts the given value into the key-value store.
func (db *Database) Put(key []byte, value []byte) error {
	return db.db.Update(func(txn *badger.Txn) error {
		return txn.Set(key, value) // Insert the key-value pair
	})
}

// Delete removes the key from the key-value store.
func (db *Database) Delete(key []byte) error {
	return db.db.Update(func(txn *badger.Txn) error {
		return txn.Delete(key) // Remove the key
	})
}

// NewBatch creates a write-only key-value store that buffers changes to its host
// database until a final write is called.
func (db *Database) NewBatch() ethdb.Batch {
	return &batch{
		db:   db.db,
		b:    db.db.NewWriteBatch(),
		size: 0,
	}
}

// batch is a write-only badger DB batch that commits changes to its host database
// when Write is called. A batch cannot be used concurrently.
type batch struct {
	db   *badger.DB
	b    *badger.WriteBatch
	size int // size is a counter to count the accumulated size of values in this batch.
}

// Put inserts the given value into the batch for later committing.
func (b *batch) Put(key []byte, value []byte) error {
	err := b.b.Set(key, value)
	if err != nil {
		return err
	}
	b.size += len(key) + len(value) // Update the size counter
	return nil
}

// Delete inserts a key removal into the batch for later committing.
func (b *batch) Delete(key []byte) error {
	err := b.b.Delete(key)
	if err != nil {
		return err
	}
	b.size += len(key) // Update the size counter for the deleted key
	return nil
}

// ValueSize retrieves the amount of data queued up for writing.
func (b *batch) ValueSize() int {
	return b.size
}

// Write flushes any accumulated data to disk.
func (b *batch) Write() error {
	err := b.b.Flush()
	if err != nil {
		return err
	}
	b.size = 0 // Reset size after writing
	return nil
}

// Reset resets the batch for reuse, the last batch must be flushed or cancel before this call.
func (b *batch) Reset() {
	b.b.Cancel()
	b.b = b.db.NewWriteBatch()
	b.size = 0 // Reset size counter
}

// Replay replays the batch contents.
func (b *batch) Replay(_ ethdb.KeyValueWriter) error {
	return fmt.Errorf("replay not implemented yet")
}
