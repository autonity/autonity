package backend

import (
	"math/big"

	"github.com/clearmatics/autonity/consensus/tendermint/wal"
)

type walStub struct{}

func (walStub) UpdateHeight(height *big.Int) error {
	return nil
}

func (walStub) Close() {}

func (walStub) Store(msg wal.Value) error {
	return nil
}

func (walStub) Get(height *big.Int) ([][]byte, error) {
	return nil, nil
}
