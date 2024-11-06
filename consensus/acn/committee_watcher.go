package acn

import (
	"context"

	"github.com/autonity/autonity/core"
	"github.com/autonity/autonity/core/types"
)

func (acn *ACN) watchCommittee(ctx context.Context) {
	acn.wg.Add(1)

	chainHeadCh := make(chan core.ChainHeadEvent)
	chainHeadSub := acn.chain.SubscribeChainHeadEvent(chainHeadCh)

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

	// read the committee base on latest state.
	currentHead := acn.chain.CurrentBlock().Header()
	currentState, err := acn.chain.StateAt(currentHead.Root)
	if err != nil {
		panic(err)
	}
	epoch, err := acn.chain.ProtocolContracts().EpochByHeight(currentHead, currentState, currentHead.Number)
	if err != nil {
		panic(err)
	}

	committee := epoch.Committee
	if committee.MemberByAddress(acn.address) != nil { //nolint
		updateConsensusEnodes(currentHead)
		wasValidating = true
	}

	go func() {
		defer acn.wg.Done()
		defer epochHeadSub.Unsubscribe()
		defer chainHeadSub.Unsubscribe()
		for {
			select {
			case ev := <-chainHeadCh:
				acn.server.SetCurrentBlockNumber(ev.Block.NumberU64())
			case ev := <-epochHeadCh:
				committee = ev.Header.Epoch.Committee
				// check if the local node belongs to the consensus committee.
				if committee.MemberByAddress(acn.address) == nil {
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
			case <-chainHeadSub.Err():
				return
			case <-ctx.Done():
				return
			}
		}
	}()
}
