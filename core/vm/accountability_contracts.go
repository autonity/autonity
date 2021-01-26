package vm

import (
	"fmt"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/params"
	"github.com/clearmatics/autonity/rlp"
)

// checkChallenge implemented as a native contract to take an on-chain challenge.
type checkChallenge struct{
	chainContext ChainContext
}
func (c *checkChallenge) InitChainContext(chain ChainContext) {
	c.chainContext = chain
	return
}
func (c *checkChallenge) RequiredGas(_ []byte) uint64 {
	return params.TakeChallengeGas
}

// checkChallenge, take challenge from AC by copy the packed byte array, decode and
// validate it, the on challenge client should send the proof of innocent via a transaction.
func (c *checkChallenge) Run(input []byte) ([]byte, error) {
	if len(input) == 0 {
		panic(fmt.Errorf("invalid proof of innocent - empty"))
	}

	err := c.CheckChallenge(input)
	if err != nil {
		return false32Byte, fmt.Errorf("invalid proof of challenge %v", err)
	}

	return true32Byte, nil
}

// validate the proof is a valid challenge.
func (c *checkChallenge) validateChallenge(p *types.Proof) error {
	// check if evidence msgs are from committee members of that height.
	h, err := p.Message.Height()
	if err != nil {
		return err
	}

	header := c.chainContext.GetHeader(p.ParentHash, h.Uint64())

	for i:=0; i < len(p.Evidence); i++ {
		m := header.CommitteeMember(p.Evidence[i].Address)
		if m == nil {
			return fmt.Errorf("evidence msg was not sent by committee member")
		}
	}

	// todo: check if the proof is a valid suspicion.
	return nil
}

// validate challenge, call from EVM package.
func (c *checkChallenge) CheckChallenge(packedProof []byte) error {
	p, err := decodeProof(packedProof)
	if err != nil {
		return err
	}

	err = c.validateChallenge(p)
	if err != nil {
		return err
	}
	return nil
}



// checkProof implemented as a native contract to validate an on-chain innocent proof.
type checkProof struct{
	chainContext ChainContext
}
func (c *checkProof) InitChainContext(chain ChainContext) {
	c.chainContext = chain
	return
}
func (c *checkProof) RequiredGas(_ []byte) uint64 {
	return params.CheckInnocentGas
}

// checkProof, take proof from AC by copy the packed byte array, decode and validate it.
func (c *checkProof) Run(input []byte) ([]byte, error) {
	// take an on-chain innocent proof, tell the results of the checking
	if len(input) == 0 {
		panic(fmt.Errorf("invalid proof of innocent - empty"))
	}

	err := c.CheckProof(input)
	if err != nil {
		return false32Byte, fmt.Errorf("invalid proof of innocent %v", err)
	}

	return true32Byte, nil
}

// Check the proof of innocent, it is called from precompiled contracts of EVM package.
func (c *checkProof) CheckProof(packedProof []byte) error {

	p, err := decodeProof(packedProof)
	if err != nil {
		return err
	}

	err = c.validateInnocentProof(p)
	if err != nil {
		return err
	}

	return nil
}

// validate the innocent proof is valid.
func (c *checkProof) validateInnocentProof(in *types.Proof) error {
	// check if evidence msgs are from committee members of that height.
	h, err := in.Message.Height()
	if err != nil {
		return err
	}

	header := c.chainContext.GetHeader(in.ParentHash, h.Uint64())

	for i:=0; i < len(in.Evidence); i++ {
		m := header.CommitteeMember(in.Evidence[i].Address)
		if m == nil {
			return fmt.Errorf("evidence msg was not sent by committee member")
		}
	}

	// todo: check if the proof is an innocent behavior.
	return nil
}

func decodeProof(proof []byte) (*types.Proof, error) {
	p := new(types.RawProof)
	err := rlp.DecodeBytes(proof, p)
	if err != nil {
		return nil, err
	}

	decodedP := new(types.Proof)
	decodedP.ParentHash = p.ParentHash
	decodedP.Rule = p.Rule

	// decode consensus message which is rlp encoded.
	msg := new(types.ConsensusMessage)
	if err := msg.FromPayload(p.Message); err != nil {
		return nil, err
	}
	decodedP.Message = *msg

	for i:= 0; i < len(p.Evidence); i++ {
		m := new(types.ConsensusMessage)
		if err := m.FromPayload(p.Evidence[i]); err != nil {
			return nil, fmt.Errorf("msg cannot be decoded")
		}
		decodedP.Evidence = append(decodedP.Evidence, *m)
	}
	return decodedP, nil
}