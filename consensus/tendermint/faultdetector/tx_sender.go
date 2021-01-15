package faultdetector

import (
	"crypto"
	"github.com/clearmatics/autonity/common"
)

type TXSender struct {
	contract common.Address // autonity contract address
	abi string // cotnract abi
	key crypto.PrivateKey // private key of etherbase
}

func (s *TXSender) SendSuspicion(suspicion *Suspicion) common.Hash {
	// todo send suspicoin and return a TX hash.
	return common.Hash{}
}

func (s *TXSender) SendInnocentProof(proof *InnocentProof) common.Hash {
	// todo send innocent proof and return a TX hash.
	return common.Hash{}
}
