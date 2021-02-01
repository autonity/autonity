package afd

import (
	"github.com/clearmatics/autonity/core/types"
	"math/big"
)

type MsgStore struct {
	//map[Height][Round][MsgType] []ConsensusMessage
	messages map[*big.Int]map[int64]map[uint64][]types.ConsensusMessage
}

func (ms *MsgStore) StoreMsg(m types.ConsensusMessage) error {
	/*
	h, err := m.Height()
	if err != nil {
		return err
	}

	r, err := m.Round()
	if err != nil {
		return err
	}
	*/
	// todo add msg into msg store.
	return nil
}