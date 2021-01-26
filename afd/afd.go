package afd

import (
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/core"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/event"
	"sync"
)

type FaultDetector struct {
	wg sync.WaitGroup
	afdFeed event.Feed
	scope event.SubscriptionScope

	blockChan chan core.ChainEvent
	blockSub event.Subscription
	blockchain *core.BlockChain
	address common.Address
}

func NewFaultDetector(chain *core.BlockChain, nodeAddress common.Address) *FaultDetector {
	fd := &FaultDetector{
		address: nodeAddress,
		blockChan:  make(chan core.ChainEvent, 300),
		blockchain: chain,
	}
	return fd
}

// listen for external events like new block be committed on the chain, to get latest view of on-chain challenges, if
// there are challenges rise to current validator, then take the challenge and provide the proof of innocent if possible.
func (fd *FaultDetector) Run() {
	fd.blockSub = fd.blockchain.SubscribeChainEvent(fd.blockChan)
	for {
		select {
		case ev := <-fd.blockChan:
			err := fd.takeMyChallenge(ev.Block, ev.Hash)
			if err != nil {
				// prints something.
			}
		}
	}
}

// get challenges from blockchain via blockchain.autonityContract calls.
func (fd *FaultDetector) takeMyChallenge(block *types.Block, hash common.Hash) error {
	//todo: get challenges from blockchain via blockchain.autonityContract go wrappers.

	//todo: get those challenges should take by myself.

	//todo: get proof for each challenge

	//todo: send proofs via SendOnChainProofs() function.
	return nil
}

func (fd *FaultDetector) Stop() {
	fd.scope.Close()
	fd.blockSub.Unsubscribe()
	fd.wg.Wait()
}

func (fd *FaultDetector) SubscribeAFDEvents(ch chan<- types.SubmitProofEvent) event.Subscription {
	return fd.scope.Track(fd.afdFeed.Subscribe(ch))
}

// afd send proof to ethereum object which will submit the on-chain proofs via transaction.
func (fd *FaultDetector) SendOnChainProofs(t types.ProofType,  proofs[]types.OnChainProof) {
	fd.wg.Add(1)
	go func() {
		defer fd.wg.Done()
		fd.afdFeed.Send(types.SubmitProofEvent{Proofs:proofs, Type:t})
	}()
}