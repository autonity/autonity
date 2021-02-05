package afd

import "github.com/clearmatics/autonity/core/types"

type RuleEngine struct {

}

// run apply rule pattern over current msg store, it returns a set of proofs of challenge, or
// accusation to be sent for on-chain management.
func (r *RuleEngine) run(ms *MsgStore) []types.OnChainProof {
	return nil
}