package acn

import (
	"context"

	"github.com/autonity/autonity/core"
	"github.com/autonity/autonity/core/types"
)

// TODO: watch for epoch rotation event instead
func (acn *ACN) watchCommittee(ctx context.Context) {
	acn.wg.Add(1)

	chainHeadCh := make(chan core.ChainHeadEvent)
	chainHeadSub := acn.chain.SubscribeChainHeadEvent(chainHeadCh)

	updateConsensusEnodes := func(block *types.Block) {
		state, err := acn.chain.StateAt(block.Header().Root)
		if err != nil {
			acn.log.Error("Could not retrieve state at head block", "err", err)
			return
		}
		enodesList, err := acn.chain.ProtocolContracts().CommitteeEnodes(block, state, true)
		if err != nil {
			acn.log.Error("Could not retrieve consensus whitelist at head block", "err", err)
			return
		}
		acn.server.UpdateConsensusEnodes(enodesList.List, enodesList.List)
	}

	wasValidating := false
	currentBlock := acn.chain.CurrentBlock()
	if currentBlock.Header().CommitteeMember(acn.address) != nil {
		updateConsensusEnodes(currentBlock)
		wasValidating = true
	}

	go func() {
		defer acn.wg.Done()
		defer chainHeadSub.Unsubscribe()
		for {
			select {
			case ev := <-chainHeadCh:
				header := ev.Block.Header()
				// check if the local node belongs to the consensus committee.
				if header.CommitteeMember(acn.address) == nil {
					// if the local node was part of the committee set for the previous block
					// there is no longer the need to retain the full connections and the
					// consensus engine enabled.
					if wasValidating {
						acn.server.UpdateConsensusEnodes(nil, nil)
						wasValidating = false
					}
					continue
				}
				updateConsensusEnodes(ev.Block)
				wasValidating = true
			// Err() channel will be closed when unsubscribing.
			case <-chainHeadSub.Err():
				return
			case <-ctx.Done():
				return
			}
		}
	}()
}
