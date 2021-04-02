package eth

import (
	"github.com/clearmatics/autonity/autonity"
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/crypto"
	"github.com/clearmatics/autonity/faultdetector"
	"github.com/clearmatics/autonity/log"
)

func (s *Ethereum) sendAccountabilityTransaction(ev *faultdetector.AccountabilityEvent) {

	tx, err := s.generateAccountabilityTX("handleProofs", ev.Proofs)
	if err != nil {
		log.Error("Could not generate accountability transaction", "err", err)
		return
	}

	// todo: check returned error for the reason of sending failures, eg. oversize, try to break proofs into multiple TXs.
	//  we cannot handle a single proof which exceed 512 KB that is the latest size of a transaction. It happens when
	//  the proof contains a big proposal which contains the entire block.
	e := s.TxPool().AddLocal(tx)
	if e != nil {
		log.Error("Cound not add TX into TX pool", "err", e)
		return
	}
	log.Debug("Generate accountability transaction", "hash", tx.Hash())
}

func (s *Ethereum) generateAccountabilityTX(method string, proofs []autonity.OnChainProof) (*types.Transaction, error) {
	nonce := s.TxPool().Nonce(crypto.PubkeyToAddress(s.defaultKey.PublicKey))
	to := s.BlockChain().GetAutonityContract().Address()
	abi := s.BlockChain().GetAutonityContract().ABI()
	packedData, err := abi.Pack(method, proofs)
	if err != nil {
		log.Error("Cannot pack accountability transaction", "err", err)
		return nil, err
	}

	// todo: estimate a reasonable gasLimit, and check if TX size exceed 512 KB, need to break proofs into multiple TXs.
	return types.SignTx(
		types.NewTransaction(nonce, to, common.Big0, 210000000, s.gasPrice, packedData),
		types.HomesteadSigner{},
		s.defaultKey)
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
