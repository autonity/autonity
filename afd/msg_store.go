package afd

import (
	"github.com/clearmatics/autonity/core/types"
	"math/big"
)

type MsgStore struct {
	//map[Height][Round][MsgType] []ConsensusMessage
	messages map[*big.Int]map[uint64]map[uint64][]types.ConsensusMessage
}

// store msg into msg store, it returns proofs of equivocation, and errEquivocation
func(ms *MsgStore) StoreMsg(m *types.ConsensusMessage) ([]types.ConsensusMessage, error) {
	return nil, nil
}