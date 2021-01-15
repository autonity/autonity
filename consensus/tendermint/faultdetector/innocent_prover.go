package faultdetector

import "github.com/clearmatics/autonity/common"

type InnocentProver struct {
	node common.Address // node address of this client
	txSender TXSender
}

type Suspicion struct {
	Rule Rule
	Message message
	Proof [][]byte
}

type InnocentProof struct {
	Rule Rule
	Message message
	RawMessages [][]byte
}

func (p *InnocentProver) resolveChallenge(suspicion *Suspicion) (innocentProof *InnocentProof, err error) {
	// todo: get innocent proof from msg store. Would need to distinguish the different rules.
	return nil, nil
}

// Check the on-chain proof of innocent.
func (p *InnocentProver) CheckProof(proof *InnocentProof) bool {
	return true
}

// prepare innocent proof for node's challenge, and send proof via transaction.
func (p *InnocentProver) TakeChallenge(suspicion *Suspicion) error {
	if suspicion.Message.Sender() != p.node {
		return nil
	}

	// the suspicion to current node, try to resolve it.
	proof, err := p.resolveChallenge(suspicion)
	if err != nil {
		return err
	}

	// send proof via transaction.
	if proof != nil {
		p.txSender.SendInnocentProof(proof)
	}

	return nil
}
