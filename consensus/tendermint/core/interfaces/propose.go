package interfaces

import (
	"context"
	"github.com/autonity/autonity/consensus/tendermint/core/messageutils"
	"github.com/autonity/autonity/core/types"
)

type Proposer interface {
	SendProposal(ctx context.Context, p *types.Block)
	HandleProposal(ctx context.Context, msg *messageutils.Message) error
	StopFutureProposalTimer()
	LogProposalMessageEvent(message string, proposal messageutils.Proposal, from, to string)
	HandleNewCandidateBlockMsg(ctx context.Context, candidateBlock *types.Block)
}
