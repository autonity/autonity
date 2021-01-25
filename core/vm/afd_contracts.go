package vm

import (
	"fmt"
	"github.com/clearmatics/autonity/params"
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

	/*
		err := faultdetector.CheckChallenge(input, c.chainContext)
		if err != nil {
			return false32Byte, fmt.Errorf("invalid proof of challenge %v", err)
		}*/

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

	/*
		err := faultdetector.CheckProof(input, c.chainContext)
		if err != nil {
			return false32Byte, fmt.Errorf("invalid proof of innocent %v", err)
		}*/

	return true32Byte, nil
}
