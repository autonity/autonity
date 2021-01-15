package faultdetector

import "github.com/clearmatics/autonity/common"

type InnocentProver struct {
	node common.Address // node address of this client
	txSender TXSender
}

type InnocentProof struct {
	Rule Rule
	Message message
	RawMessages [][]byte // use raw msg payload as proofs to be verified on precompile contract
}

func (p *InnocentProver) proveInnocent(suspicion Proof) (innocentProof *InnocentProof, err error) {
	// todo: get innocent proof from msg store. Would need to distinguish the different rules.
	return nil, nil
}

// prepare innocent proof for node's challenge, and send proof via transaction.
func (p *InnocentProver) ResolveSuspicion(suspicion Proof) error {
	if suspicion.Message.Sender() != p.node {
		return nil
	}

	// the suspicion to current node, try to resolve it.
	proof, err := p.proveInnocent(suspicion)
	if err != nil {
		return err
	}

	// send proof via transaction.
	if proof != nil {
		p.txSender.SendInnocentProof(proof)
	}

	return nil
}
