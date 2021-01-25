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

	err := CheckChallenge(input, c.chainContext)
	if err != nil {
		return false32Byte, fmt.Errorf("invalid proof of challenge %v", err)
	}

	return true32Byte, nil
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

	err := CheckProof(input, c.chainContext)
	if err != nil {
		return false32Byte, fmt.Errorf("invalid proof of innocent %v", err)
	}

	return true32Byte, nil
}

// validate the proof is a valid challenge.
func validateChallenge(c *types.Proof, chain ChainContext) error {
	// get committee from block header
	h, err := c.Message.Height()
	if err != nil {
		return nil
	}

	header := chain.GetHeader(c.ParentHash, h.Uint64())

	// check if evidences senders are presented in committee.
	for i:=0; i < len(c.Evidence); i++ {
		member := header.CommitteeMember(c.Evidence[i].Address)
		if member == nil {
			return fmt.Errorf("invalid evidence for proof of susipicous message")
		}
	}

	// todo: check if the proof is a valid suspicion.
	return nil
}

// validate the innocent proof is valid.
func validateInnocentProof(in *types.Proof, chain ChainContext) error {
	// get committee from block header
	h, err := in.Message.Height()
	if err != nil {
		return nil
	}
	header := chain.GetHeader(in.ParentHash, h.Uint64())

	// check if evidences senders are presented in committee.
	for i:=0; i < len(in.Evidence); i++ {
		member := header.CommitteeMember(in.Evidence[i].Address)
		if member == nil {
			return fmt.Errorf("invalid evidence for proof of susipicous message")
		}
	}
	// todo: check if the proof is an innocent behavior.
	return nil
}

// Check the proof of innocent, it is called from precompiled contracts of EVM package.
func CheckProof(packedProof []byte, chain ChainContext) error {
	p := new(types.Proof)
	err := rlp.DecodeBytes(packedProof, p)
	if err != nil {
		return err
	}

	err = validateInnocentProof(p, chain)
	if err != nil {
		return err
	}

	return nil
}

// validate challenge, call from EVM package.
func CheckChallenge(packedProof []byte, chain ChainContext) error {
	p := new(types.Proof)
	err := rlp.DecodeBytes(packedProof, p)
	if err != nil {
		return err
	}

	err = validateChallenge(p, chain)
	if err != nil {
		return err
	}
	return nil
}