package core

import (
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/core/types"
	"math/big"
)

type coreStateRequestEvent struct {
	stateChan chan TendermintState
}

// VoteState save the prevote or precommit voting status for a specific value.
type VoteState struct {
	Value            common.Hash
	ProposalVerified bool
	VotePower        uint64
}

// RoundState save the voting status for a specific round.
type RoundState struct {
	Round          int64
	Proposal       common.Hash
	PrevoteState   []VoteState
	PrecommitState []VoteState
}

// MsgWithHash save the msg and extra field to be marshal to JSON.
type MsgForDump struct {
	Message
	Hash   common.Hash
	Power  uint64
	Height *big.Int
	Round  int64
}

// TendermintState save an instant status for the tendermint consensus engine.
type TendermintState struct {
	// validator address
	Client common.Address

	// core state of tendermint
	Height      *big.Int
	Round       int64
	Step        uint64
	Proposal    *common.Hash
	LockedValue *common.Hash
	LockedRound int64
	ValidValue  *common.Hash
	ValidRound  int64

	// committee state
	Committee       types.Committee
	Proposer        common.Address
	IsProposer      bool
	QuorumVotePower uint64
	RoundStates     []RoundState
	ProposerPolicy  uint64

	// extra state
	SentProposal          bool
	SentPrevote           bool
	SentPrecommit         bool
	SetValidRoundAndValue bool

	// timer state
	BlockPeriod           uint64
	ProposeTimerStarted   bool
	PrevoteTimerStarted   bool
	PrecommitTimerStarted bool

	// current height messages.
	CurHeightMessages []*MsgForDump
	// backlog msgs
	BacklogMessages []*MsgForDump
	// backlog unchecked msgs.
	UncheckedMsgs []*MsgForDump
	// Known msg of gossip.
	KnownMsgHash []common.Hash
}

func (c *core) CoreState() TendermintState {
	// send state dump request.
	var e = coreStateRequestEvent{
		stateChan: make(chan TendermintState),
	}
	go c.sendEvent(e)
	return <-e.stateChan
}

// State Dump is handled in the main loop triggered by an event rather than using RLOCK mutex.
func (c *core) handleStateDump(e coreStateRequestEvent) {
	state := TendermintState{
		Client:            c.address,
		BlockPeriod:       c.blockPeriod,
		CurHeightMessages: msgForDump(c.GetCurrentHeightMessages()),
		BacklogMessages:   getBacklogMsgs(c),
		UncheckedMsgs:     getBacklogUncheckedMsgs(c),
		// tendermint core state:
		Height:      c.Height(),
		Round:       c.Round(),
		Step:        uint64(c.step),
		Proposal:    getProposal(c, c.Round()),
		LockedValue: getHash(c.lockedValue),
		LockedRound: c.lockedRound,
		ValidValue:  getHash(c.validValue),
		ValidRound:  c.validRound,

		// committee state
		Committee:       c.committeeSet().Committee(),
		Proposer:        c.committeeSet().GetProposer(c.Round()).Address,
		IsProposer:      c.isProposer(),
		QuorumVotePower: c.committeeSet().Quorum(),
		RoundStates:     getRoundState(c),
		// extra state
		SentProposal:          c.sentProposal,
		SentPrevote:           c.sentPrevote,
		SentPrecommit:         c.sentPrecommit,
		SetValidRoundAndValue: c.setValidRoundAndValue,
		// timer state
		ProposeTimerStarted:   c.proposeTimeout.timerStarted(),
		PrevoteTimerStarted:   c.prevoteTimeout.timerStarted(),
		PrecommitTimerStarted: c.precommitTimeout.timerStarted(),
		// known msgs in case of gossiping.
		KnownMsgHash: c.backend.KnownMsgHash(),
	}

	// for none blocking send state.
	c.logger.Debug("sending core state msg")
	e.stateChan <- state
	// let sender to close channel.
	close(e.stateChan)
}

func getBacklogUncheckedMsgs(c *core) []*MsgForDump {
	result := make([]*MsgForDump, 0)
	for _, ms := range c.backlogUnchecked {
		result = append(result, msgForDump(ms)...)
	}

	return result
}

// getBacklogUncheckedMsgs and getBacklogMsgs are kind of redundant code,
// don't know how to write it via golang like template in C++, since the only
// difference is the type of the data it operate on.
func getBacklogMsgs(c *core) []*MsgForDump {
	result := make([]*MsgForDump, 0)
	for _, ms := range c.backlogs {
		result = append(result, msgForDump(ms)...)
	}

	return result
}

func msgForDump(msgs []*Message) []*MsgForDump {
	result := make([]*MsgForDump, 0, len(msgs))
	for _, m := range msgs {
		msg := new(MsgForDump)
		msg.Message = *m
		msg.Power = m.GetPower()
		msg.Hash = types.RLPHash(m.payload)

		// in case of haven't decode msg yet, set round and height as -1.
		msg.Round = -1
		msg.Height = big.NewInt(-1)
		msg.Round, _ = m.Round()
		msg.Height, _ = m.Height()
		result = append(result, msg)
	}
	return result
}

func getProposal(c *core, round int64) *common.Hash {
	if c.messages.getOrCreate(round).proposal != nil && c.messages.getOrCreate(round).proposal.ProposalBlock != nil {
		v := c.messages.getOrCreate(round).proposal.ProposalBlock.Hash()
		return &v
	}
	return nil
}

func getHash(b *types.Block) *common.Hash {
	if b != nil {
		v := b.Hash()
		return &v
	}
	return nil
}

func getRoundState(c *core) []RoundState {
	rounds := c.messages.getRounds()
	states := make([]RoundState, 0, len(rounds))

	for _, r := range rounds {
		proposal, prevoteState, preCommitState := getVoteState(&c.messages, r)
		state := RoundState{
			Round:          r,
			Proposal:       proposal,
			PrevoteState:   prevoteState,
			PrecommitState: preCommitState,
		}
		states = append(states, state)
	}
	return states
}

func blockHashes(messages map[common.Hash]map[common.Address]Message) []common.Hash {
	blockHashes := make([]common.Hash, 0, len(messages))
	for key := range messages {
		blockHashes = append(blockHashes, key)
	}
	return blockHashes
}

func getVoteState(s *messagesMap, round int64) (common.Hash, []VoteState, []VoteState) {
	p := common.Hash{}

	if s.getOrCreate(round).Proposal() != nil && s.getOrCreate(round).Proposal().ProposalBlock != nil {
		p = s.getOrCreate(round).Proposal().ProposalBlock.Hash()
	}

	preVoteValues := blockHashes(s.getOrCreate(round).prevotes.votes)
	preCommitValues := blockHashes(s.getOrCreate(round).precommits.votes)
	prevoteState := make([]VoteState, 0, len(preVoteValues))
	precommitState := make([]VoteState, 0, len(preCommitValues))

	for _, v := range preVoteValues {
		var s = VoteState{
			Value:            v,
			ProposalVerified: s.getOrCreate(round).isProposalVerified(),
			VotePower:        s.getOrCreate(round).PrevotesPower(v),
		}
		prevoteState = append(prevoteState, s)
	}

	for _, v := range preCommitValues {
		var s = VoteState{
			Value:            v,
			ProposalVerified: s.getOrCreate(round).isProposalVerified(),
			VotePower:        s.getOrCreate(round).PrecommitsPower(v),
		}
		precommitState = append(precommitState, s)
	}

	return p, prevoteState, precommitState
}
