package faultdetector

import "github.com/clearmatics/autonity/common"

type InnocentProver struct {
	node common.Address // node address of this client
	txSender TXSender
}

// validate the proof is a valid challenge.
func (p *InnocentProver) validateChallenge(c *Proof) error {
	return nil
}

// validate the innocent proof is valid.
func (p *InnocentProver) validateInnocentProof(i *Proof) error {
	return nil
}

func (p *InnocentProver) resolveChallenge(c *Proof) (innocentProof *Proof, err error) {
	// todo: get innocent proof from msg store. Would need to distinguish the different rules.
	return nil, nil
}

// Check the proof of innocent, it is called from precompiled contracts of EVM package.
func (p *InnocentProver) CheckProof(packedProof []byte) error {
	innocentProof := unpackProof(packedProof)
	err := p.validateInnocentProof(innocentProof)
	if err != nil {
		return err
	}

	return nil
}

// validate challenge, and send proof via transaction if current client is on challenge, call from EVM package.
func (p *InnocentProver) TakeChallenge(packedProof []byte) error {
	challenge := unpackProof(packedProof)

	err := p.validateChallenge(challenge)
	if err != nil {
		return err
	}

	if challenge.Message.Sender() != p.node {
		return nil
	}

	// the suspicion to current node, try to resolve it.
	proof, err := p.resolveChallenge(challenge)
	if err != nil {
		return err
	}

	// send proof via transaction.
	if proof != nil {
		p.txSender.SendInnocentProof(proof)
	}

	return nil
}
