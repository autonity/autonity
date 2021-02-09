package afd

import (
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/core/types"
)

// todo: integrate msg store in this file.
type MsgStore struct {
	//map[Height]map[Round]map[MsgType]map[common.address][]ConsensusMessage
	/*
	  height->round->MsgCode->MsgSender-> [msg, msg,,,]
	 */
	messages map[uint64]map[int64]map[uint64]map[common.Address][]types.ConsensusMessage
}

// store msg into msg store, it returns msgs that is equivocation than the input msg, and errEquivocation
func(ms *MsgStore) StoreMsg(m *types.ConsensusMessage) ([]types.ConsensusMessage, error) {
	return nil, nil
}

// clean those ancient msgs.
func(ms *MsgStore) cleanAncientMsg(fromHeight uint64) {
	//todo: clean up ancient msgs by block height.
}