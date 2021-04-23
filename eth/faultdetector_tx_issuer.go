package eth

import (
	"errors"
	"fmt"

	"github.com/clearmatics/autonity/autonity"
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/core"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/crypto"
	"github.com/clearmatics/autonity/log"
)

var (
	errOverSizedEvent = errors.New("oversized accountability event")
)

func (s *Ethereum) sendAccountabilityTXs(onChainProofs []*autonity.OnChainProof) {

	txs, err := s.generateAccountabilityTXs("handleProofs", onChainProofs)
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
func (s *Ethereum) generateAccountabilityTXs(method string, onChainProofs []*autonity.OnChainProof) (txs []*types.Transaction, e error) {
	nonce := s.TxPool().Nonce(crypto.PubkeyToAddress(s.nodeKey.PublicKey))
	// try to generate a single event to contain all the onChainProofs.
	tx, err := s.generateAccountabilityTX(nonce, method, onChainProofs)
	if err == nil {
		return append(txs, tx), nil
	}

	// accountability events exceed 512 KB, break the events into pieces.
	if err == errOverSizedEvent {
		if len(onChainProofs) == 1 {
			log.Error("over-sized accountability event", "err", "cannot pack over-sized proof")
			return nil, errOverSizedEvent
		}

		// try to pack as much events as possible until TX exceed 512 KB.
		start := 0
		for i := 1; i <= len(onChainProofs) && start < len(onChainProofs); i++ {
			tx, err := s.generateAccountabilityTX(nonce, method, onChainProofs[start:i])
			// exceed 512 KB, try to break it into pieces.
			if err == errOverSizedEvent {
				if len(onChainProofs[start:i]) == 1 {
					//single event exceed 512 KB, skip it.
					start++
					log.Error("skip over-sized accountability event", "err")
					continue
				}

				// break sub piece of events
				p, err := s.generateAccountabilityTX(nonce, method, onChainProofs[start:i-1])
				if err == nil {
					start = i - 1
					i = start
					nonce++
					txs = append(txs, p)
				}
			}

			// append the last piece of events
			if err == nil && i == len(onChainProofs) {
				txs = append(txs, tx)
			}
		}

		return txs, nil
	}

	return nil, err
}

func (s *Ethereum) generateAccountabilityTX(nonce uint64, method string, onChainProofs []*autonity.OnChainProof) (*types.Transaction, error) {
	to := autonity.ContractAddress
	abi := s.BlockChain().GetAutonityContract().ABI()

	var proofs = make([]autonity.OnChainProof, len(onChainProofs))
	for i, p := range onChainProofs {
		proofs[i] = *p
	}

	packedData, err := abi.Pack(method, proofs)
	if err != nil {
		log.Error("Cannot pack accountability transaction", "err", err)
		return nil, err
	}

	tx, err := types.SignTx(types.NewTransaction(nonce, to, common.Big0, 210000000, s.gasPrice, packedData), types.HomesteadSigner{}, s.nodeKey)
	if err != nil {
		return nil, err
	}

	if uint64(tx.Size()) > core.TxMaxSize {
		return nil, errOverSizedEvent
	}
	return tx, nil
}

func (s *Ethereum) faultDetectorTXEventLoop() {
	go func() {
		for {
			select {
			case onChainProofs := <-s.faultDetectorCh:
				s.sendAccountabilityTXs(onChainProofs)
			case err, ok := <-s.faultDetectorSub.Err():
				if ok {
					panic(fmt.Sprintf("faultDetectorSub error: %v", err.Error()))
				}
				return
			}
		}
	}()
}
