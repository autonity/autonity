package afd

import (
	"github.com/clearmatics/autonity/core/types"
	"math/big"
)

type MsgStore struct {
	//map[Height][Round][MsgType] []ConsensusMessage
	messages map[*big.Int]map[uint64]map[uint64][]types.ConsensusMessage
}