// TODO: Copyrights template here
package tm

import (
	"math/big"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core/types"
)

const (
	// This is the target size for the packs of transactions or announcements. A
	// pack can get larger than this if a single transactions exceeds this size.
	maxTxPacketSize = 100 * 1024
)

// blockPropagation is a block propagation event, waiting for its turn in the
// broadcast queue.
type blockPropagation struct {
	block *types.Block
	td    *big.Int
}

// broadcastBlocks is a write loop that multiplexes blocks and block accouncements
// to the remote peer. The goal is to have an async writer that does not lock up
// node internals and at the same time rate limits queued data.
func (p *Peer) broadcastBlocks() {
	for {
		select {
		case prop := <-p.queuedBlocks:
			if err := p.SendNewBlock(prop.block, prop.td); err != nil {
				return
			}
			p.Log().Trace("Propagated block", "number", prop.block.Number(), "hash", prop.block.Hash(), "td", prop.td)

		case block := <-p.queuedBlockAnns:
			if err := p.SendNewBlockHashes([]common.Hash{block.Hash()}, []uint64{block.NumberU64()}); err != nil {
				return
			}
			p.Log().Trace("Announced block", "number", block.Number(), "hash", block.Hash())

		case <-p.term:
			return
		}
	}
}
