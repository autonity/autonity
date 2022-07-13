package interfaces

import (
	"context"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/core/messageutils"
	"math/big"
)

type Precommiter interface {
	SendPrecommit(ctx context.Context, isNil bool)
	HandlePrecommit(ctx context.Context, msg *messageutils.Message) error
	VerifyCommittedSeal(addressMsg common.Address, committedSealMsg []byte, proposedBlockHash common.Hash, round int64, height *big.Int) error
	HandleCommit(ctx context.Context)
	LogPrecommitMessageEvent(message string, precommit messageutils.Vote, from, to string)
}
