package core

import (
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"math/big"
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

	var views []*message.RoundMsgView

	for r, roundMsgs := range c.futureRound {
		roundView := &message.RoundMsgView{}
		roundView.Round = uint64(r)

		preVoteSigners := make(map[common.Hash]*message.AggregatedPower)
		preCommitSigners := make(map[common.Hash]*message.AggregatedPower)

		for _, m := range roundMsgs {
			if m.Code() == message.ProposalCode {
				roundView.Proposal = m.Value()
			}

			if m.Code() == message.PrevoteCode {
				value := m.Value()
				_, ok := preVoteSigners[value]
				if !ok {
					preVoteSigners[value] = message.NewAggregatedPower()
				}
				for index, _ := range m.(message.Vote).Signers().Powers() {
					preVoteSigners[value].SetBit(index)
				}
			}

			if m.Code() == message.PrecommitCode {
				value := m.Value()
				_, ok := preCommitSigners[value]
				if !ok {
					preCommitSigners[value] = message.NewAggregatedPower()
				}
				for index, _ := range m.(message.Vote).Signers().Powers() {
					preCommitSigners[value].SetBit(index)
				}
			}
		}

		var preVotes []common.Hash
		var preVoteSigns []*big.Int
		for v, s := range preVoteSigners {
			preVotes = append(preVotes, v)
			preVoteSigns = append(preVoteSigns, s.Signers())
		}

		roundView.Prevotes = preVotes
		roundView.PrevotesSigners = preVoteSigns

		var preCommits []common.Hash
		var preCommitSigns []*big.Int
		for v, s := range preCommitSigners {
			preCommits = append(preCommits, v)
			preCommitSigns = append(preCommitSigns, s.Signers())
		}
		roundView.Precommits = preCommits
		roundView.PrecommitsSigners = preCommitSigns

		views = append(views, roundView)
	}

	return views
}
