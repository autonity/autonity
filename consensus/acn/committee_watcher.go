package acn

import (
	"context"

	"github.com/autonity/autonity/core"
	"github.com/autonity/autonity/core/types"
)

func (acn *ACN) watchCommittee(ctx context.Context) {
	acn.wg.Add(1)

	epochHeadCh := make(chan core.EpochHeadEvent)
	epochHeadSub := acn.chain.SubscribeEpochHeadEvent(epochHeadCh)

	updateConsensusEnodes := func(header *types.Header) {
		state, err := acn.chain.StateAt(header.Root)
		if err != nil {
			acn.log.Error("Could not retrieve state at head block", "err", err)
			return
		}
		enodesList, err := acn.chain.ProtocolContracts().CommitteeEnodes(header, state, true)
		if err != nil {
			acn.log.Error("Could not retrieve consensus whitelist at head block", "err", err)
			return
		}
		acn.server.UpdateConsensusEnodes(enodesList.List, enodesList.List)
	}

	wasValidating := false
	committee, currentHead := acn.chain.ChainHeadAndCommittee()
	if committee.CommitteeMember(acn.address) != nil {
		updateConsensusEnodes(currentHead)
		wasValidating = true
	}

	go func() {
		defer acn.wg.Done()
		defer epochHeadSub.Unsubscribe()
		for {
			select {
			case ev := <-epochHeadCh:
				header := ev.Header
				// check if the local node belongs to the consensus committee.
				if header.Committee.CommitteeMember(acn.address) == nil {
					// if the local node was part of the committee set for the previous block
					// there is no longer the need to retain the full connections and the
					// consensus engine enabled.
					if wasValidating {
						acn.server.UpdateConsensusEnodes(nil, nil)
						wasValidating = false
					}
					continue
				}
				updateConsensusEnodes(ev.Header)
				wasValidating = true
			// Err() channel will be closed when unsubscribing.
			case <-epochHeadSub.Err():
				return
			case <-ctx.Done():
				return
			}
		}
	}()
}
