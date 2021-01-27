package eth

import (
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/crypto"
	"github.com/clearmatics/autonity/log"
)

func (s *Ethereum) sendAccountabilityTransaction(e *types.SubmitProofEvent) {
	var method string
	if e.Type == types.InnocentProof {
		method = "resolveChallenge"
	}

	if e.Type == types.ChallengeProof {
		method = "addChallenge"
	}

	tx, err := s.generateAccountabilityTX(method, e.Proofs)
	if err != nil {
		log.Error("Could not generate accountability transaction", "err", err)
	}

	err = s.TxPool().AddLocal(tx)
	log.Debug("Generate accountability transaction", "hash", tx.Hash())
}

func (s *Ethereum) generateAccountabilityTX(method string, params ...interface{} ) (*types.Transaction, error) {
	nonce := s.TxPool().Nonce(crypto.PubkeyToAddress(s.defaultKey.PublicKey))
	to := s.BlockChain().GetAutonityContract().Address()
	abi := s.BlockChain().GetAutonityContract().ABI()
	packedData, err := abi.Pack(method, params)
	if err != nil {
		return nil, err
	}

	// todo: resolve a reasonable gasLimit.
	return types.SignTx(
		types.NewTransaction(nonce,	to,	common.Big0,210000000, s.gasPrice, packedData),
		types.HomesteadSigner{},
		s.defaultKey)
}

func (s *Ethereum) afdTXEventLoop() {
	for {
		select {
		case event := <-s.afdCh:
			//todo: send accountability proofs.
			s.sendAccountabilityTransaction(&event)
		// Err() channel will be closed when unsubscribing.
		case <-s.afdSub.Err():
			return
		}
	}
}