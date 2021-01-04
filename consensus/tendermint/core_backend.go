package tendermint

import (
	"context"

	"github.com/clearmatics/autonity/consensus"
	"github.com/clearmatics/autonity/core"
	"github.com/clearmatics/autonity/core/types"
)

type Tendermint interface {
	consensus.Handler
	Start(ctx context.Context, blockchain *core.BlockChain)
	Stop()
	Seal(chain consensus.ChainReader, block *types.Block, results chan<- *types.Block, stop <-chan struct{}) error
}
