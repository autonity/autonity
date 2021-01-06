package tendermint

import (
	"context"
	"math/big"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus"
	"github.com/clearmatics/autonity/core"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/rpc"
)

type Tendermint interface {
	consensus.Handler
	Seal(chain consensus.ChainReader, block *types.Block, results chan<- *types.Block, stop <-chan struct{}) error
	APIs(chain consensus.ChainReader) []rpc.API
	Prepare(chain consensus.ChainReader, header *types.Header) error
	CalcDifficulty(chain consensus.ChainReader, time uint64, parent *types.Header) *big.Int
	Author(header *types.Header) (common.Address, error)
	Start(ctx context.Context, blockchain *core.BlockChain) error
	Close() error
}
