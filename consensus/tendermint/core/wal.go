package core

import (
	"math/big"

	"github.com/clearmatics/autonity/consensus/tendermint/wal"
)

type WAL interface {
	UpdateHeight(height *big.Int) error
	Store(msg wal.Value) error
	Get(height *big.Int) ([][]byte, error)
	Close()
}
