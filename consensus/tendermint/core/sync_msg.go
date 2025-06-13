package core

import (
	"math/big"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
)

func (c *Core) createSyncMsg() *message.AskSyncMsg {
	msg := &message.AskSyncMsg{}
	msg.Height = c.Height().Uint64()
	msg.KnownMessages = c.Messages().DumpMsgView()

	// and future rounds
	future := c.futureRoundMsgView()
	if len(future) > 0 {
		msg.KnownMessages = append(msg.KnownMessages, future...)
	}

	return msg
}

func (c *Core) futureRoundMsgView() []*message.RoundMsgView {
	c.futureRoundLock.RLock()
	defer c.futureRoundLock.RUnlock()

	views := make([]*message.RoundMsgView, 0, 32)

	for r, roundMsgs := range c.futureRound {
		roundView := &message.RoundMsgView{}
		roundView.Round = uint64(r)

		preVoteSigners := make(map[common.Hash]*big.Int)
		preCommitSigners := make(map[common.Hash]*big.Int)

		for _, m := range roundMsgs {
			if m.Code() == message.ProposalCode {
				roundView.HaveProposal = true
			}

			if m.Code() == message.PrevoteCode {
				value := m.Value()
				_, ok := preVoteSigners[value]
				if !ok {
					preVoteSigners[value] = new(big.Int)
				}
				for _, index := range m.(message.Vote).Signers().FlattenUniq() {
					preVoteSigners[value].SetBit(preVoteSigners[value], index, 1)
				}
			}

			if m.Code() == message.PrecommitCode {
				value := m.Value()
				_, ok := preCommitSigners[value]
				if !ok {
					preCommitSigners[value] = new(big.Int)
				}
				for _, index := range m.(message.Vote).Signers().FlattenUniq() {
					preCommitSigners[value].SetBit(preCommitSigners[value], index, 1)
				}
			}
		}

		var preVotes []common.Hash
		var preVoteSigns []*big.Int
		for v, s := range preVoteSigners {
			preVotes = append(preVotes, v)
			preVoteSigns = append(preVoteSigns, s)
		}

		roundView.Prevotes = preVotes
		roundView.PrevotesSigners = preVoteSigns

		var preCommits []common.Hash
		var preCommitSigns []*big.Int
		for v, s := range preCommitSigners {
			preCommits = append(preCommits, v)
			preCommitSigns = append(preCommitSigns, s)
		}
		roundView.Precommits = preCommits
		roundView.PrecommitsSigners = preCommitSigns

		views = append(views, roundView)
	}

	return views
}
