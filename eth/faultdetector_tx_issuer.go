package eth

import (
	"github.com/clearmatics/autonity/autonity"
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus/tendermint/faultdetector"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/crypto"
	"github.com/clearmatics/autonity/log"
)

func (s *Ethereum) sendAccountabilityTransaction(ev *faultdetector.AccountabilityEvent) {

	tx, err := s.generateAccountabilityTX("handleProofs", ev.Proofs)
	if err != nil {
		log.Error("Could not generate accountability transaction", "err", err)
		return
	}

	e := s.TxPool().AddLocal(tx)
	if e != nil {
		log.Error("Cound not add TX into TX pool", "err", e)
		return
	}
	log.Debug("Generate accountability transaction", "hash", tx.Hash())
}

func (s *Ethereum) generateAccountabilityTX(method string, proofs []autonity.OnChainProof) (*types.Transaction, error) {
	nonce := s.TxPool().Nonce(crypto.PubkeyToAddress(s.nodeKey.PublicKey))
	to := s.BlockChain().GetAutonityContract().Address()
	abi := s.BlockChain().GetAutonityContract().ABI()
	packedData, err := abi.Pack(method, proofs)
	if err != nil {
		return nil, err
	}

	// might to resolve a reasonable gas limit by weighting the bytes of TX.
	return types.SignTx(
		types.NewTransaction(nonce, to, common.Big0, 210000000, s.gasPrice, packedData),
		types.HomesteadSigner{},
		s.nodeKey)
}

func (s *Ethereum) faultDetectorTXEventLoop() {
	go func() {
		for {
			select {
			case event := <-s.faultDetectorCh:
				s.sendAccountabilityTransaction(&event)
			// Err() channel will be closed when unsubscribing.
			case <-s.faultDetectorSub.Err():
				return
			}
		}
	}()
}
