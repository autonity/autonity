package eth

import (
	"github.com/clearmatics/autonity/afd"
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/crypto"
	"github.com/clearmatics/autonity/log"
)

func (s *Ethereum) sendAccountabilityTransaction(e *afd.SubmitProofEvent) {
	var method string
	if e.Type == afd.InnocentProof {
		method = "resolveAccusation"
	}

	if e.Type == afd.ChallengeProof {
		method = "addChallenge"
	}

	if e.Type == afd.AccusationProof {
		method = "addAccusation"
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

	// might to resolve a reasonable gas limit by weighting the bytes of TX.
	return types.SignTx(
		types.NewTransaction(nonce,	to,	common.Big0,210000000, s.gasPrice, packedData),
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