package interfaces

import (
	"context"
	"github.com/autonity/autonity/consensus/tendermint/core/messageutils"
)

type Prevoter interface {
	SendPrevote(ctx context.Context, isNil bool, badProposal *messageutils.BadProposalInfo)
	HandlePrevote(ctx context.Context, msg *messageutils.Message) error
	LogPrevoteMessageEvent(message string, prevote messageutils.Vote, from, to string)
}
