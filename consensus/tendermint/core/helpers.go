package core

import (
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/core/interfaces"
	mUtils "github.com/autonity/autonity/consensus/tendermint/core/messageutils"
	"math/big"
)

// OverQuorumVotes compute voting power out from a set of prevotes or precommits of a certain round and height, the caller
// should make sure that the votes belong to a certain round and height, it returns a set of votes that the corresponding
// voting power is over quorum, otherwise it returns nil.
func OverQuorumVotes(msgs []*mUtils.Message, quorum uint64) (overQuorumVotes []*mUtils.Message) {
	votingPower := new(big.Int).SetUint64(0)

	counted := make(map[common.Address]struct{})
	for _, v := range msgs {
		if _, ok := counted[v.Address]; ok {
			continue
		}

		counted[v.Address] = struct{}{}
		votingPower = votingPower.Add(votingPower, v.GetPower())
		overQuorumVotes = append(overQuorumVotes, v)
		if votingPower.Cmp(new(big.Int).SetUint64(quorum)) >= 0 {
			return overQuorumVotes
		}
	}

	return nil
}

func LiteProposalSignature(backend interfaces.Backend, height *big.Int, round int64, validRound int64, value common.Hash) ([]byte, error) {
	liteProposal := &mUtils.LiteProposal{
		Round:      round,
		Height:     height,
		ValidRound: validRound,
		Value:      value,
	}

	payload, err := liteProposal.PayloadNoSig()
	if err != nil {
		return nil, err
	}

	sigLiteProposal, err := backend.Sign(payload)
	if err != nil {
		return nil, err
	}
	return sigLiteProposal, nil
}
