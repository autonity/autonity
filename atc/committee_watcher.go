package atc

import (
	"context"

	"github.com/autonity/autonity/core"
	"github.com/autonity/autonity/core/types"
)

func (atc *ATC) watchCommittee(ctx context.Context) {
	atc.wg.Add(1)

	chainHeadCh := make(chan core.ChainHeadEvent)
	chainHeadSub := atc.chain.SubscribeChainHeadEvent(chainHeadCh)

	updateConsensusEnodes := func(block *types.Block) {
		state, err := atc.chain.StateAt(block.Header().Root)
		if err != nil {
			atc.log.Error("Could not retrieve state at head block", "err", err)
			return
		}
		enodesList, err := atc.chain.ProtocolContracts().CommitteeEnodes(block, state, true)
		if err != nil {
			atc.log.Error("Could not retrieve consensus whitelist at head block", "err", err)
			return
		}

		atc.server.UpdateConsensusEnodes(enodesList.List)
	}

	wasValidating := false
	currentBlock := atc.chain.CurrentBlock()
	if currentBlock.Header().CommitteeMember(atc.address) != nil {
		updateConsensusEnodes(currentBlock)
		wasValidating = true
	}

	go func() {
		defer atc.wg.Done()
		defer chainHeadSub.Unsubscribe()
		for {
			select {
			case ev := <-chainHeadCh:
				header := ev.Block.Header()
				// check if the local node belongs to the consensus committee.
				if header.CommitteeMember(atc.address) == nil {
					// if the local node was part of the committee set for the previous block
					// there is no longer the need to retain the full connections and the
					// consensus engine enabled.
					if wasValidating {
						atc.server.UpdateConsensusEnodes(nil)
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
