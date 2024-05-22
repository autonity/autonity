package message

import (
	"math/big"
	"sync"

	"github.com/autonity/autonity/common"
)

//TODO(lorenzo) refinements2, analyze more duplicated msgs and equivocation scenarios

//TODO(lorenzo) not sure this is the right place for this + add tests for it

// auxiliary data structure to take into account aggregated power of a set of signers
type PowerInfo struct {
	power   *big.Int
	signers *big.Int // used as bitmap, we do not care about coefficients here, only if a validator is present or not
}

func (p *PowerInfo) Set(index int, power *big.Int) {
	if p.signers.Bit(index) == 1 {
		return
	}

	p.signers.SetBit(p.signers, index, 1)
	p.power.Add(p.power, power)
}

func (p *PowerInfo) Pow() *big.Int {
	return p.power
}

func (p *PowerInfo) Copy() *PowerInfo {
	return &PowerInfo{power: new(big.Int).Set(p.power), signers: new(big.Int).Set(p.signers)}
}

func NewPowerInfo() *PowerInfo {
	return &PowerInfo{power: new(big.Int), signers: new(big.Int)}
}

type Set struct {
	// In some conditions we might receive prevotes or precommit before
	// receiving a proposal, so we must save received message with different proposed block hash.
	votes map[common.Hash][]Vote // map[proposedBlockHash][]vote

	/* we use PowerInfo because we cannot simply sum the voting power of the votes. This is because we might have:
	* 1. duplicated votes between different overlapping aggregates for the same value
	* 2. equivocated votes from the same validator across different values
	 */
	powers     map[common.Hash]*PowerInfo // cumulative voting power for each value
	totalPower *PowerInfo                 // total voting power of votes

	sync.RWMutex //TODO(lorenzo) refinements, do we need this lock since there is already one is round_messages?
}

func NewSet() *Set {
	return &Set{
		votes:      make(map[common.Hash][]Vote),
		powers:     make(map[common.Hash]*PowerInfo),
		totalPower: NewPowerInfo(),
	}
}

func (s *Set) Add(vote Vote) {
	s.Lock()
	defer s.Unlock()

	value := vote.Value()
	previousVotes, ok := s.votes[value]
	if !ok {
		s.votes[value] = make([]Vote, 1)
		s.powers[value] = NewPowerInfo()
	}

	// update total power and power for value
	for index, power := range vote.Signers().Powers() {
		s.totalPower.Set(index, power)
		s.powers[value].Set(index, power)
	}

	// check if we are adding the first vote
	if len(previousVotes) == 0 {
		s.votes[value][0] = vote
		return
	}

	// if not first vote, aggregate previous votes and new vote
	//TODO(lorenzo) performance, verify that this doesn't create too much memory
	switch vote.(type) {
	case *Prevote:
		aggregatedVotes := AggregatePrevotesSimple(append(previousVotes, vote))
		s.votes[value] = make([]Vote, len(aggregatedVotes))
		for i, aggregatedVote := range aggregatedVotes {
			s.votes[value][i] = aggregatedVote
		}
	case *Precommit:
		aggregatedVotes := AggregatePrecommitsSimple(append(previousVotes, vote))
		s.votes[value] = make([]Vote, len(aggregatedVotes))
		for i, aggregatedVote := range aggregatedVotes {
			s.votes[value][i] = aggregatedVote
		}
	default:
		panic("Trying to add a vote that is not Prevote nor Precommit")
	}
}

func (s *Set) Messages() []Msg {
	s.RLock()
	defer s.RUnlock()

	messages := make([]Msg, 0)
	for _, votes := range s.votes {
		for _, vote := range votes {
			messages = append(messages, vote.(Msg))
		}
	}
	return messages
}

func (s *Set) PowerFor(h common.Hash) *big.Int {
	s.RLock()
	defer s.RUnlock()
	_, ok := s.powers[h]
	if ok {
		return new(big.Int).Set(s.powers[h].power)
	}
	return new(big.Int)
}

func (s *Set) TotalPower() *big.Int {
	s.RLock()
	defer s.RUnlock()

	// NOTE: in case of equivocated messages, we count power only once
	// TODO(lorenzo) refinements, write a test for it
	return new(big.Int).Set(s.totalPower.power)
}

func (s *Set) VotesFor(blockHash common.Hash) []Vote {
	s.RLock()
	defer s.RUnlock()

	return s.votes[blockHash]
}
