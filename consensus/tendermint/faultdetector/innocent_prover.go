package faultdetector

// validate the proof is a valid challenge.
func validateChallenge(c *Proof) error {
	// todo: check if messages are signed by correct committee member of its round. AFD should buffer N block headers at package level,
	// otherwise static precompiled contract cannot validate whether messages are from correct committee member.
	// todo: check if the suspicious message is proved by given evidence as a valid suspicion.
	return nil
}

// validate the innocent proof is valid.
func validateInnocentProof(i *Proof) error {
	// todo: check if messages are signed by correct committee member of its round. AFD should buffer N block headers at package level,
	// otherwise static precompiled contract cannot validate whether messages are from correct committee member.
	// todo: check if the suspicious message is proved by given evidence as an innocent behavior.
	return nil
}

// used by those who is on-challenge, to get innocent proof from msg store.
func resolveChallenge(c *Proof) (innocentProof *Proof, err error) {
	// todo: get innocent proof from msg store. Would need to distinguish the different rules.
	return nil, nil
}

// Check the proof of innocent, it is called from precompiled contracts of EVM package.
func CheckProof(packedProof []byte) error {
	innocentProof, err := UnpackProof(packedProof)
	if err != nil {
		return err
	}
	err = validateInnocentProof(innocentProof)
	if err != nil {
		return err
	}

	return nil
}

// validate challenge, call from EVM package.
func CheckChallenge(packedProof []byte) error {
	challenge, err := UnpackProof(packedProof)
	if err != nil {
		return err
	}

	err = validateChallenge(challenge)
	if err != nil {
		return err
	}
	return nil
}
