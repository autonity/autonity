package faultdetector

import (
	"fmt"
	"github.com/clearmatics/autonity/core"
)

// validate the proof is a valid challenge.
func validateChallenge(c *Proof, chain core.ChainContext) error {
	// get committee from block header
	header := chain.GetHeader(c.parentHash, uint64(c.Message.Height()-1))

	// check if evidences senders are presented in committee.
	for i:=0; i < len(c.Evidence); i++ {
		member := header.CommitteeMember(c.Evidence[i].Sender())
		if member == nil {
			return fmt.Errorf("invalid evidence for proof of susipicous message")
		}
	}

	// todo: check if the suspicious message is proved by given evidence as a valid suspicion.
	return nil
}

// validate the innocent proof is valid.
func validateInnocentProof(in *Proof, chain core.ChainContext) error {
	// get committee from block header
	header := chain.GetHeader(in.parentHash, uint64(in.Message.Height()-1))

	// check if evidences senders are presented in committee.
	for i:=0; i < len(in.Evidence); i++ {
		member := header.CommitteeMember(in.Evidence[i].Sender())
		if member == nil {
			return fmt.Errorf("invalid evidence for proof of susipicous message")
		}
	}
	// todo: check if the suspicious message is proved by given evidence as an innocent behavior.
	return nil
}

// used by those who is on-challenge, to get innocent proof from msg store.
func resolveChallenge(c *Proof) (innocentProof *Proof, err error) {
	// todo: get innocent proof from msg store. Would need to distinguish the different rules.
	return nil, nil
}

// Check the proof of innocent, it is called from precompiled contracts of EVM package.
func CheckProof(packedProof []byte, chain core.ChainContext) error {
	innocentProof, err := UnpackProof(packedProof)
	if err != nil {
		return err
	}
	err = validateInnocentProof(innocentProof, chain)
	if err != nil {
		return err
	}

	return nil
}

// validate challenge, call from EVM package.
func CheckChallenge(packedProof []byte, chain core.ChainContext) error {
	challenge, err := UnpackProof(packedProof)
	if err != nil {
		return err
	}

	err = validateChallenge(challenge, chain)
	if err != nil {
		return err
	}
	return nil
}
