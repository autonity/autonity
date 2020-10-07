package core

import (
	"crypto/rand"
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

func newMessageStore() *messageCache {
	return &messageCache{
		msgHashes:     make(map[uint64]map[int64]map[algorithm.Step]map[common.Address]common.Hash),
		consensusMsgs: make(map[common.Hash]*algorithm.ConsensusMessage),
		rawMessages:   make(map[common.Hash]*Message),
		valid:         make(map[common.Hash]struct{}),
		values:        make(map[common.Hash]*types.Block),
	}

}

func (m *messageCache) Message(h common.Hash) *Message {
	return m.rawMessages[h]
}

func (m *messageCache) signatures(valueHash common.Hash, round int64, height uint64) [][]byte {
	var sigs [][]byte

	// println("signatures -----")
	for _, msgHash := range m.msgHashes[height][round][algorithm.Step(msgPrecommit)] {
		// spew.Dump(m.rawMessages[msgHash].decodedMsg)
		if valueHash == m.rawMessages[msgHash].decodedMsg.ProposedValueHash() {
			sigs = append(sigs, m.rawMessages[msgHash].CommittedSeal)
		}
	}
	// println("----------------")
	return sigs
}

func (m *messageCache) prevoteQuorum(valueHash *common.Hash, round int64, header *types.Header) bool {
	msgType := new(algorithm.Step)
	*msgType = algorithm.Prevote
	// println("prevote power --------")
	vp := m.votePower(valueHash, round, msgType, header) >= header.Committee.Quorum()
	// println(vp, "----------------")
	return vp
}

func (m *messageCache) precommitQuorum(valueHash *common.Hash, round int64, header *types.Header) bool {
	msgType := new(algorithm.Step)
	*msgType = algorithm.Precommit
	// // println("precommit power", m.votePower(valueHash, round, msgType, header))
	// println("precommit power --------")
	vp := m.votePower(valueHash, round, msgType, header) >= header.Committee.Quorum()
	// println(vp, "----------------")
	return vp
}

func (m *messageCache) fail(round int64, header *types.Header) bool {
	// println("fail power --------")
	vp := m.votePower(nil, round, nil, header) >= header.Committee.Quorum()
	// println(vp, "----------------")
	return vp
}

func (m *messageCache) votePower(
	valueHash *common.Hash, // A nil value hash indicates that we match any value.
	round int64,
	msgType *algorithm.Step, // A nil value hash indicates that we match both prevote and precommit.
	header *types.Header,
) uint64 {

	// Only prevotes and precommits impart vote power.
	if msgType != nil && !msgType.In(algorithm.Prevote, algorithm.Precommit) {
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
	for mType, addressMap := range m.msgHashes[header.Number.Uint64()+1][round] {
		// spew.Dump(addressMap)
		// Skip in the case that this is not a message type we are considering.
		if msgType != nil && *msgType != mType {
			continue
		}
		for address, msgHash := range addressMap {
			// Skip messages not considered valid
			_, ok := m.valid[msgHash]
			if !ok {
				// println("skippng not valid")
				continue
			}

			// Skip messages with differing values
			if valueHash != nil && *valueHash != common.Hash(m.consensusMsgs[msgHash].Value) {
				// // println("skipping mismatch value")
				continue
			}
			// spew.Dump(m.consensusMsgs[msgHash])
			// Now either value hash is nil (matches everything) or it actually matches the msg's value.
			power += header.CommitteeMember(address).VotingPower.Uint64()
		}
	}
	return power
}

func addMsgHash(
	hashes map[uint64]map[int64]map[algorithm.Step]map[common.Address]common.Hash,
	height uint64,
	round int64,
	msgType algorithm.Step,
	address common.Address,
	hash common.Hash,
) error {
	// todo check bounds
	roundMap, ok := hashes[height]
	if !ok {
		roundMap = make(map[int64]map[algorithm.Step]map[common.Address]common.Hash)
		hashes[height] = roundMap
	}
	msgTypeMap, ok := roundMap[round]
	if !ok {
		msgTypeMap = make(map[algorithm.Step]map[common.Address]common.Hash)
		roundMap[round] = msgTypeMap
	}
	addressMap, ok := msgTypeMap[msgType]
	if !ok {
		addressMap = make(map[common.Address]common.Hash)
		msgTypeMap[msgType] = addressMap
	}
	addressMap[address] = hash // TODO check for duplicates here, accountablitiy

	return nil
}

func (m *messageCache) addMessage(msg *Message, cm *algorithm.ConsensusMessage) error {
	r := make([]byte, 2)
	_, err := rand.Read(r)
	if err != nil {
		panic(err)
	}
	// id := hex.EncodeToString(r)

	// // println(id, "hashes len", len(m.msgHashes))
	// // println(id, spew.Sdump(m.msgHashes))
	// // println(id, "add message", cm.String(), msg.Hash.String())
	err = addMsgHash(m.msgHashes, cm.Height, cm.Round, cm.MsgType, msg.Address, msg.Hash)
	if err != nil {
		return err
	}
	// // println(id, "hashes len", len(m.msgHashes))
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

type messageProcessor func(cm *algorithm.ConsensusMessage) error

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

func (m *messageCache) proposal(height uint64, round int64, proposer common.Address) *algorithm.ConsensusMessage {
	return m.consensusMsgs[m.msgHashes[height][round][algorithm.Step(msgProposal)][proposer]]
}

func (m *messageCache) matchingProposal(cm *algorithm.ConsensusMessage) *algorithm.ConsensusMessage {
	if cm.MsgType == algorithm.Step(msgProposal) {
		return cm
	}
	// if cm.MsgType == algorithm.Precommit {
	// 	fmt.Printf("fetching proposal for: %s", cm.String())
	// }
	for _, proposalHash := range m.msgHashes[cm.Height][cm.Round][algorithm.Step(msgProposal)] {
		proposal := m.consensusMsgs[proposalHash]
		if proposal.Value == cm.Value {
			// fmt.Printf(" got: %s\n", proposal.String())
			return proposal
		}
	}
	// fmt.Printf(" no proposal\n")
	return nil
}

// TODO cover the case where we receive multiple proposals for future heights and we don't know who the propose is?
