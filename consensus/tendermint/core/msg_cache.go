package core

import (
	"fmt"

	common "github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus/tendermint/algorithm"
	types "github.com/clearmatics/autonity/core/types"
)

// Message cache caches messages

type messageCache struct {
	// msgHashes maps height, round, message type and address to message hash.
	msgHashes map[uint64]map[int64]map[algorithm.Step]map[common.Address]common.Hash
	// valid is a set containing message hashes for messages considered valid.
	valid map[common.Hash]struct{}
	// consensusMsgs maps message hash to consensus message.
	consensusMsgs map[common.Hash]*algorithm.ConsensusMessage
	// rawMessages maps message hash to raw message.
	rawMessages map[common.Hash]*Message
	// values maps value hash to value.
	values map[common.Hash]*types.Block

	// TODO use this to impose upper and lower limits on what can go in the
	// message cache.
	currentHeight uint64
}

func (m *messageCache) heightMessages(height uint64) []*Message {
	var messages []*Message
	for _, msgTypeMap := range m.msgHashes[height] {
		for _, addressMap := range msgTypeMap {
			for _, hash := range addressMap {
				messages = append(messages, m.rawMessages[hash])
			}
		}
	}
	return messages
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

func (m *messageCache) signatures(valueHash common.Hash, round int64, height uint64) [][]byte {
	var sigs [][]byte
	for _, msgHash := range m.msgHashes[height][round][algorithm.Step(msgPrecommit)] {
		if valueHash == m.rawMessages[msgHash].decodedMsg.ProposedValueHash() {
			sigs = append(sigs, m.rawMessages[msgHash].CommittedSeal)
		}
	}
	return sigs
}

func (m *messageCache) prevoteQuorum(valueHash *common.Hash, round int64, header *types.Header) bool {
	msgType := new(algorithm.Step)
	*msgType = algorithm.Step(msgPrevote)
	return m.votePower(valueHash, round, msgType, header) >= header.Committee.Quorum()
}

func (m *messageCache) precommitQuorum(valueHash *common.Hash, round int64, header *types.Header) bool {
	msgType := new(algorithm.Step)
	*msgType = algorithm.Step(msgPrecommit)
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
	msgType *algorithm.Step, // A nil value hash indicates that we match both prevote and precommit.
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
			if valueHash != nil && *valueHash != m.consensusMsgs[msgHash].Value {
				continue
			}
			// Now either value hash is nil (matches everything) or it actually matches the msg's value.
			power += header.CommitteeMember(address).VotingPower.Uint64()
		}
	}
	return power
}

func addMsgHash(
	hashes map[uint64]map[int64]map[consensusMessageType]map[common.Address]common.Hash,
	height uint64,
	round int64,
	msgType consensusMessageType,
	address common.Address,
	hash common.Hash,
) error {
	// todo check bounds
	roundMap, ok := hashes[height]
	if !ok {
		roundMap = make(map[int64]map[consensusMessageType]map[common.Address]common.Hash)
	}
	msgTypeMap, ok := roundMap[round]
	if !ok {
		msgTypeMap = make(map[consensusMessageType]map[common.Address]common.Hash)
	}
	addressMap, ok := msgTypeMap[msgType]
	if !ok {
		addressMap = make(map[common.Address]common.Hash)
	}
	addressMap[address] = hash // TODO check for duplicates here, accountablitiy
	return nil
}

func (m *messageCache) addMessage(msg *Message, cm *consensusMessage) error {
	err := addMsgHash(m.msgHashes, cm.Height, cm.Round, cm.MsgType, msg.Address, msg.Hash)
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

type messageProcessor func(cm *consensusMessage) error

func (m *messageCache) roundMessages(height uint64, round int64, p messageProcessor) error {
	for _, addressMap := range m.msgHashes[height][round] {
		for _, msgHash := range addressMap {
			if _, ok := m.valid[msgHash]; ok {
				err := p(m.consensusMsgs[msgHash])
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (m *messageCache) proposal(height uint64, round int64, proposer common.Address) *consensusMessage {
	return m.consensusMsgs[m.msgHashes[height][round][consensusMessageType(msgProposal)][proposer]]
}

func (m *messageCache) matchingProposal(cm *consensusMessage) *consensusMessage {
	if cm.MsgType == consensusMessageType(msgProposal) {
		return cm
	}
	for _, proposalHash := range m.msgHashes[cm.Height][cm.Round][consensusMessageType(msgProposal)] {
		proposal := m.consensusMsgs[proposalHash]
		if proposal.Value == cm.Value {
			return proposal
		}
	}
	return nil
}

// TODO cover the case where we receive multiple proposals for future heights and we don't know who the propose is?
