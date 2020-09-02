package core

import (
	"fmt"

	common "github.com/clearmatics/autonity/common"
	types "github.com/clearmatics/autonity/core/types"
)

// Message cache caches messages

type messageCache struct {
	proposalMsgHashes map[uint64]map[int64]map[common.Address]common.Hash
	msgHashToProposal map[common.Hash]*Proposal

	prevoteMsgHashes map[uint64]map[int64]map[common.Address]common.Hash
	msgHashToPrevote map[common.Hash]*Vote

	precommitMsgHashes map[uint64]map[int64]map[common.Address]common.Hash
	msgHashToPrecommit map[common.Hash]*Vote

	// msgHashes maps height, round, message type and address to message hash.
	msgHashes map[uint64]map[int64]map[consensusMessageType]map[common.Address]common.Hash
	// valid is a set containing message hashes for messages considered valid.
	valid map[common.Hash]struct{}
	// consensusMsgs maps message hash to consensus message.
	consensusMsgs map[common.Hash]*consensusMessage
	// rawMessages maps message hash to raw message.
	rawMessages map[common.Hash]*Message
	// values maps value hash to value.
	values map[common.Hash]*types.Block

	// TODO use this to impose upper and lower limits on what can go in the
	// message cache.
	currentHeight uint64
}

// func roundMap(msgHashes map[uint64]map[int64]map[common.Address]common.Hash) map[int64]map[common.Address]common.Hash {
// }

// func (m *messageCache) addMessage(msg *Message) error {
// 	// TODO check bounds and reject message
// 	switch msg.Code {
// 	case msgProposal:
// 		var p Proposal
// 		err := msg.Decode(&p)
// 		if err != nil {
// 			return err
// 		}
// 		roundMap, ok := m.proposalHashes[p.Height.Uint64()]
// 		if !ok {
// 			roundMap = make(map[int64]map[common.Address]*Proposal)
// 		}
// 		msgMap, ok := roundMap[p.Round]
// 		if !ok {
// 			msgMap = make(map[common.Address]*Proposal)
// 		}
// 		msgMap[msg.Address] = &p
// 	case msgPrevote:
// 		var preV Vote
// 		err := msg.Decode(&preV)
// 		if err != nil {
// 			return err
// 		}
// 		addHash(m.prevoteHashes, &preV, msg.Address)
// 	case msgPrecommit:
// 		var preC Vote
// 		err := msg.Decode(&preC)
// 		if err != nil {
// 			return err
// 		}
// 		addHash(m.prevoteHashes, &preC, msg.Address)
// 	default:
// 		panic(fmt.Sprintf("Unrecognised message code %d for message: %s", msg.Code, spew.Sdump(msg)))
// 	}
// 	return nil
// }

// func (m *messageCache) addProposal(p *Proposal, msg *Message) error {
// 	err := addMsgHash(m.proposalMsgHashes, p.Height.Uint64(), p.Round, msg.Address, msg.Hash)
// 	if err != nil {
// 		return err
// 	}
// 	m.msgHashToProposal[msg.Hash] = p
// 	m.rawMessages[msg.Hash] = msg
// 	return nil
// }

// func (m *messageCache) addPrevote(p *Vote, msg *Message) error {
// 	err := addMsgHash(m.prevoteMsgHashes, p.Height.Uint64(), p.Round, msg.Address, msg.Hash)
// 	if err != nil {
// 		return err
// 	}
// 	m.msgHashToPrevote[msg.Hash] = p
// 	m.rawMessages[msg.Hash] = msg
// 	return nil
// }

// func (m *messageCache) addPrecommit(p *Vote, msg *Message) error {
// 	err := addMsgHash(m.precommitMsgHashes, p.Height.Uint64(), p.Round, msg.Address, msg.Hash)
// 	if err != nil {
// 		return err
// 	}
// 	m.msgHashToPrecommit[msg.Hash] = p
// 	m.rawMessages[msg.Hash] = msg

// 	return nil
// }

func (m *messageCache) Message(h common.Hash) *Message {
	return m.rawMessages[h]
}

func (m *messageCache) prevotePower(valueHash common.Hash, round int64, header *types.Header) uint64 {
	return votePower(m.prevoteMsgHashes, m.msgHashToPrevote, valueHash, round, header)
}

func (m *messageCache) precommitPower(valueHash common.Hash, round int64, header *types.Header) uint64 {
	return votePower(m.precommitMsgHashes, m.msgHashToPrecommit, valueHash, round, header)
}

func (m *messageCache) totalPrevotePower(round int64, header *types.Header) uint64 {
	return totalVotePower(m.prevoteMsgHashes, round, header)
}

func (m *messageCache) totalPrecommitPower(round int64, header *types.Header) uint64 {
	return totalVotePower(m.precommitMsgHashes, round, header)
}

func totalVotePower(voteMsgHashes map[uint64]map[int64]map[common.Address]common.Hash, round int64, header *types.Header) uint64 {
	var total uint64
	// Iterate all prevotes for the round and total their voting power.
	for address := range voteMsgHashes[header.Number.Uint64()][round] {
		total += header.CommitteeMember(address).VotingPower.Uint64()
	}
	return total
}

func (m *messageCache) signatures(valueHash common.Hash, round int64, height uint64) [][]byte {
	var sigs [][]byte
	for _, msgHash := range m.precommitMsgHashes[height][round] {
		if valueHash == m.msgHashToPrecommit[msgHash].ProposedBlockHash {
			sigs = append(sigs, m.rawMessages[msgHash].CommittedSeal)
		}
	}
	return sigs
}
func (m *messageCache) prevoteQuorum(valueHash *common.Hash, round int64, header *types.Header) bool {
	msgType := &uint8(msgPrevote)
	return m.votePower(valueHash, round, msgType, header) >= header.Committee.Quorum()
}

func (m *messageCache) precommitQuorum(valueHash *common.Hash, round int64, header *types.Header) bool {
	msgType := &uint8(msgPrecommit)
	return m.votePower(valueHash, round, msgType, header) >= header.Committee.Quorum()
}

func (m *messageCache) fail(round int64, header *types.Header) bool {
	return m.votePower(nil, round, nil, header) >= header.Committee.Quorum()
}

// func (m *messageCache) futureRoundFail(round int64, header *types.Header) bool {
// 	// Only prevotes and precommits impart vote power.
// 	if uint8(msgType) != msgPrevote || uint8(msgType) != msgPrecommit {
// 		panic(fmt.Sprintf(
// 			"Unexpected msgType %d, expecting either %d or %d",
// 			msgType,
// 			msgPrevote,
// 			msgPrecommit,
// 		))
// 	}

// 	// Total the power of all votes in this height and round for this value,
// 	// failure to find a committee member in the header indicates a programming
// 	// error and an invalid memory acccess panic will ensue.
// 	var power uint64
// 	// For all messages at the given height in the given round of the given type ...
// 	for address, msgHash := range m.msgHashes[header.Number.Uint64()][round][msgType] {
// 		// Skip messages not considered valid
// 		_, ok := m.valid[msgHash]
// 		if !ok {
// 			continue
// 		}

// 		// Skip messages with differing values
// 		if valueHash != nil && *valueHash != m.consensusMsgs[msgHash].value {
// 			continue
// 		}
// 		// Now either value hash is nil (matches everything) or it actually matches the msg's value.
// 		power += header.CommitteeMember(address).VotingPower.Uint64()
// 	}
// 	return power
// }

// func (m *messageCache) hasQuorum(
// 	valueHash *common.Hash,
// 	round int64,
// 	msgType uint8,
// 	header *types.Header,

// ) bool {
// 	return m.votePower(valueHash, round, msgType, header) >= header.Committee.Quorum()
// }

// type messageMatcher func(valueHash *common.Hash, round int64, msgType uint8)

// var defaultMatcher messageMatcher = func(m messageCache, valueHash *common.Hash, round int64, msgType uint8) bool {
// 	// Skip messages not considered valid
// 	_, ok := m.valid[msgHash]
// 	if !ok {
// 		continue
// 	}

// 	// Skip messages with differing values
// 	if valueHash != nil && *valueHash != m.consensusMsgs[msgHash].value {
// 		continue
// 	}
// 	return false
// }

// // This is required in order to be able to support aggregating votes from multiple rounds. See the upon condition from line 55.
// type roundIterator interface {
// 	// Next returns the next msgTypeMap from the round map provided, a return
// 	// value of nil indicates that the iteration has finished.
// 	next(roundMap map[int64]map[uint8]map[common.Address]common.Hash) (msgTypeMap map[uint8]map[common.Address]common.Hash)
// }

// // singleRoundIterator returns the msgTypeMap for a single round only.
// type singleRoundIterator struct {
// 	round int64
// }

// func (sr *singleRoundIterator) next(roundMap map[int64]map[uint8]map[common.Address]common.Hash) map[uint8]map[common.Address]common.Hash {
// 	return roundMap[sr.round]
// }

// // higerRoundIterator returns the msgTypeMap for a single round only.
// type higerRoundIterator struct {
// 	round int64
// }

// func (hr *higerRoundIterator) next(roundMap map[int64]map[uint8]map[common.Address]common.Hash) map[uint8]map[common.Address]common.Hash {
// 	for round, msgTypeMap := range roundMap {
// 		if
// 		// return roundMap[sr.round]

// 	}
// }

func (m *messageCache) votePower(
	valueHash *common.Hash, // A nil value hash indicates that we match any value.
	round int64,
	msgType *consensusMessageType, // A nil value hash indicates that we match both prevote and precommit.
	header *types.Header,
) uint64 {

	// Only prevotes and precommits impart vote power.
	if msgType != nil && !msgType.in(msgPrevote, msgPrecommit) {
		panic(fmt.Sprintf(
			"Unexpected msgType %d, expecting either %d or %d",
			*msgType,
			msgPrevote,
			msgPrecommit,
		))
	}

	// Total the power of all votes in this height and round for this value,
	// failure to find a committee member in the header indicates a programming
	// error and an invalid memory acccess panic will ensue.
	var power uint64
	// For all messages at the given height in the given round ...
	for mType, addressMap := range m.msgHashes[header.Number.Uint64()][round] {
		// Skip in the case that this is not a message type we are considering.
		if msgType != nil && *msgType != mType {
			continue
		}
		for address, msgHash := range addressMap {
			// Skip messages not considered valid
			_, ok := m.valid[msgHash]
			if !ok {
				continue
			}

			// Skip messages with differing values
			if valueHash != nil && *valueHash != m.consensusMsgs[msgHash].value {
				continue
			}
			// Now either value hash is nil (matches everything) or it actually matches the msg's value.
			power += header.CommitteeMember(address).VotingPower.Uint64()
		}
	}
	return power
}

func (m *messageCache) proposal(height uint64, round int64, proposerAddress common.Address) *Proposal {
	roundMap, ok := m.proposalMsgHashes[height]
	if !ok {
		return nil
	}
	addressMap, ok := roundMap[round]
	if !ok {
		return nil
	}
	msgHash, ok := addressMap[proposerAddress]
	if !ok {
		return nil
	}
	return m.msgHashToProposal[msgHash]
}

func (m *messageCache) proposalVerified(proposalHash common.Hash) bool {
	return false
}

// func (m *messageCache) getproposals(height uint64, round int64) map[common.Address]*Proposal {
// 	roundMap, ok := m.proposalHashes[height]
// 	if !ok {
// 		return nil
// 	}
// 	return roundMap[round]
// }

// func (m *messageCache) getprevotes(height uint64, round int64) map[common.Address]*Vote {
// 	return getVotes(height, round, m.prevoteHashes)
// }

// func (m *messageCache) getprecommits(height uint64, round int64) map[common.Address]*Vote {
// 	return getVotes(height, round, m.precommitHashes)
// }

func getVotes(height uint64, round int64, votes map[uint64]map[int64]map[common.Address]*Vote) map[common.Address]*Vote {
	roundMap, ok := votes[height]
	if !ok {
		return nil
	}
	return roundMap[round]
}

func addMsgHash(
	hashes map[uint64]map[int64]map[common.Address]common.Hash,
	height uint64,
	round int64,
	msgType uint8,
	address common.Address,
	hash common.Hash,
) error {
	// todo check bounds
	roundMap, ok := hashes[height]
	if !ok {
		roundMap = make(map[int64]map[uint8]map[common.Address]common.Hash)
	}
	msgTypeMap, ok := roundMap[round]
	if !ok {
		msgTypeMap = make(map[uint8]map[common.Address]common.Hash)
	}
	addressMap, ok := msgTypeMap[msgType]
	if !ok {
		addressMap = make(map[common.Address]common.Hash)
	}
	addressMap[address] = hash // TODO check for duplicates here, accountablitiy
	return nil
}

// No need to worry about duplicates and accountability here since we check that in addMsgHash.
func addVoteForValue(votes map[common.Hash]map[common.Address]*Vote, valueHash common.Hash, address common.Address, vote *Vote) {
	addressMap, ok := votes[valueHash]
	if !ok {
		addressMap = make(map[common.Address]*Vote)
	}
	addressMap[address] = vote
}

func addProposeForValue(proposals map[common.Hash]map[common.Address]*Proposal, valueHash common.Hash, address common.Address, vote *Proposal) {
	addressMap, ok := proposals[valueHash]
	if !ok {
		addressMap = make(map[common.Address]*Proposal)
	}
	addressMap[address] = vote
}

func (m *messageCache) heightMessages(height uint64) []*Message {
	var messages []*Message
	accumulateMessagesForHeight(m.proposalMsgHashes, m.rawMessages, height, messages)
	accumulateMessagesForHeight(m.prevoteMsgHashes, m.rawMessages, height, messages)
	accumulateMessagesForHeight(m.precommitMsgHashes, m.rawMessages, height, messages)
	return messages
}

func accumulateMessagesForHeight(msgHashes map[uint64]map[int64]map[common.Address]common.Hash, msgHashToMsg map[common.Hash]*Message, height uint64, accumulator []*Message) {
	// Accumuate all messages for all rounds at the given height
	for _, addressMap := range msgHashes[height] {
		for _, hash := range addressMap {
			accumulator = append(accumulator, msgHashToMsg[hash])
		}
	}
}

func (m *messageCache) addMessage(msg *Message, cm *consensusMessage) error {
	err := addMsgHash(m.msgHashes, cm.height, cm.round, cm.msgType, msg.Address, msg.Hash)
	if err != nil {
		return err
	}
	m.consensusMsgs[msg.Hash] = cm
	m.rawMessages[msg.Hash] = msg

	return nil
}
func (m *messageCache) addValue(valueHash common.Hash, value *types.Block) {
	m.values[valueHash] = value
}
func (m *messageCache) value(valueHash common.Hash) *types.Block {
	return m.values[valueHash]
}

// Mark the hash of something valid, it could be a message hash or a value hash
func (m *messageCache) setValid(itemHash common.Hash) {
	m.valid[itemHash] = struct{}{}
}
func (m *messageCache) isValid(itemHash common.Hash) bool {
	_, ok := m.valid[itemHash]
	return ok
}
