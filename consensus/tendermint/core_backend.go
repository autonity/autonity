package tendermint

import (
	"context"

	"github.com/clearmatics/autonity/consensus"
	"github.com/clearmatics/autonity/core"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/rpc"
)

type Tendermint interface {
	consensus.Handler
	Start(ctx context.Context, blockchain *core.BlockChain)
	Stop()
	Seal(chain consensus.ChainReader, block *types.Block, results chan<- *types.Block, stop <-chan struct{}) error
	APIs(chain consensus.ChainReader) []rpc.API
}
