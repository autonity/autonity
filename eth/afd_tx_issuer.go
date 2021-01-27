package eth

import "github.com/clearmatics/autonity/core/types"

func (s *Ethereum) sendAccountabilityTransaction(e types.SubmitProofEvent) {
	if e.Type == types.InnocentProof {
		//todo: pack innocent proof transaction with abi and send it.
	}

	if e.Type == types.ChallengeProof {
		//todo: pack challenge proof transaction with abi and send it.
	}
}

func (s *Ethereum) afdTXEventLoop() {
	for {
		select {
		case event := <-s.afdCh:
			//todo: send accountability proofs.
			s.sendAccountabilityTransaction(event)
		// Err() channel will be closed when unsubscribing.
		case <-s.afdSub.Err():
			return
		}
	}
}