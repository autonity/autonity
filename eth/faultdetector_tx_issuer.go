package eth

import (
	"github.com/clearmatics/autonity/autonity"
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/core"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/crypto"
	"github.com/clearmatics/autonity/faultdetector"
	"github.com/clearmatics/autonity/log"
	"github.com/syndtr/goleveldb/leveldb/errors"
)

var(
	errOverSizedEvent = errors.New("oversized accountability event")
)

func (s *Ethereum) sendAccountabilityTransaction(ev *faultdetector.AccountabilityEvent) {

	txs, err := s.generateAccountabilityTXs("handleProofs", ev.Proofs)
	if err != nil {
		log.Error("Could not generate accountability transaction", "err", err)
		return
	}

	for _, tx := range txs {
		e := s.TxPool().AddLocal(tx)
		if e != nil {
			log.Error("Cound not add TX into TX pool", "err", e)
			continue
		}
		log.Info("Generate accountability transaction", "hash", tx.Hash())
	}
}

// generate on-chain events for accountability, it take the proofs and pack them into the accountability contract
// interface, since max transaction size was limited into 512 KB, so we need to estimate the size of the event, and
// consider to break them into pieces once the proofs exceed 512 KB.
func (s *Ethereum) generateAccountabilityTXs(method string, proofs []autonity.OnChainProof) (txs []*types.Transaction, e error) {
	nonce := s.TxPool().Nonce(crypto.PubkeyToAddress(s.defaultKey.PublicKey))
	// try to generate a single event to contain all the proofs.
	tx, err := s.genAccountabilityEvent(nonce, method, proofs)
	if err == nil {
		return append(txs, tx), nil
	}

	// accountability events exceed 512 KB, break the events into pieces.
	if err == errOverSizedEvent {
		if len(proofs) == 1 {
			log.Error("over-sized accountability event", "err", "cannot pack over-sized proof")
			return nil, errOverSizedEvent
		}

		// try to pack as much events as possible until TX exceed 512 KB.
		start := 0
		for i:=1; i<=len(proofs) && start < len(proofs); i++ {
			tx, err := s.genAccountabilityEvent(nonce, method, proofs[start:i])
			// exceed 512 KB, try to break it into pieces.
			if err == errOverSizedEvent {
				if len(proofs[start:i]) == 1 {
					//single event exceed 512 KB, skip it.
					start++
					log.Error("skip over-sized accountability event", "err")
					continue
				}

				// break sub piece of events
				p, err := s.genAccountabilityEvent(nonce, method, proofs[start:i-1])
				if err == nil {
					start = i-1
					i = start
					nonce++
					txs = append(txs, p)
				}
			}

			// append the last piece of events
			if err == nil && i == len(proofs){
				txs = append(txs, tx)
			}
		}

		return txs, nil
	}

	return nil, err
}

func (s *Ethereum) genAccountabilityEvent(nonce uint64, method string, proofs []autonity.OnChainProof) (*types.Transaction, error) {
	to := s.BlockChain().GetAutonityContract().Address()
	abi := s.BlockChain().GetAutonityContract().ABI()
	packedData, err := abi.Pack(method, proofs)
	if err != nil {
		log.Error("Cannot pack accountability transaction", "err", err)
		return nil, err
	}

	tx, err := types.SignTx(types.NewTransaction(nonce, to, common.Big0, 210000000, s.gasPrice, packedData), types.HomesteadSigner{}, s.defaultKey)
	if err != nil {
		return nil, err
	}

	if uint64(tx.Size()) > core.TxMaxSize {
		return nil, errOverSizedEvent
	}
	return tx, nil
}

func (s *Ethereum) afdTXEventLoop() {
	for {
		select {
		case event := <-s.afdCh:
			s.sendAccountabilityTransaction(&event)
		// Err() channel will be closed when unsubscribing.
		case <-s.afdSub.Err():
			return
		}
	}
}
