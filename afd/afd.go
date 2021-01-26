package afd

import (
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/event"
	"sync"
)

type FaultDetector struct {
	wg sync.WaitGroup
	afdFeed event.Feed
	scope event.SubscriptionScope
}

func (fd *FaultDetector) SubscribeAFDEvents(ch chan<- types.SubmitProofEvent) event.Subscription {
	return fd.scope.Track(fd.afdFeed.Subscribe(ch))
}

func (fd *FaultDetector) SendOnChainProofs(t types.ProofType,  proofs[]types.OnChainProof) {
	fd.wg.Add(1)
	go func() {
		defer fd.wg.Done()
		fd.afdFeed.Send(types.SubmitProofEvent{Proofs:proofs, Type:t})
	}()
}