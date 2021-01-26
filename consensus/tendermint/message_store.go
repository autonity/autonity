package tendermint

import (
	"fmt"

	common "github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus/tendermint/algorithm"
	types "github.com/clearmatics/autonity/core/types"
)

// messageStore stores messages
type messageStore struct {
	// msgHashes maps height, round, message type and address to message hash.
	msgHashes map[uint64]map[int64]map[algorithm.Step]map[common.Address]common.Hash
	// valid is a set containing message hashes for messages considered valid.
	valid map[common.Hash]struct{}
	// messages maps message hash to message.
	messages map[common.Hash]*message
	// rawMessages maps message hash to raw message bytes.
	rawMessages map[common.Hash][]byte
	// values maps value hash to value.
	values map[common.Hash]*types.Block
	// messageBounds keeps track of the heights for which the messageStore retains messages.
	messageBounds *bounds
}

func (m *messageStore) heightMessages(height uint64) []*message {
	var messages []*message
	for _, msgTypeMap := range m.msgHashes[height] {
		for _, addressMap := range msgTypeMap {
			for _, hash := range addressMap {
				messages = append(messages, m.messages[hash])
			}
		}
	}
	return messages
}

func (m *messageStore) rawHeightMessages(height uint64) [][]byte {
	var messages [][]byte
	for _, msgTypeMap := range m.msgHashes[height] {
		for _, addressMap := range msgTypeMap {
			for _, hash := range addressMap {
				messages = append(messages, m.rawMessages[hash])
			}
		}
	}
	return messages
}

func newMessageStore(messageBounds *bounds) *messageStore {
	return &messageStore{
		msgHashes:     make(map[uint64]map[int64]map[algorithm.Step]map[common.Address]common.Hash),
		rawMessages:   make(map[common.Hash][]byte),
		messages:      make(map[common.Hash]*message),
		valid:         make(map[common.Hash]struct{}),
		values:        make(map[common.Hash]*types.Block),
		messageBounds: messageBounds,
	}

}

func (m *messageStore) Message(h common.Hash) *message {
	return m.messages[h]
}

func (m *messageStore) signatures(value algorithm.ValueID, round int64, height uint64) [][]byte {
	var sigs [][]byte

	// println("signatures -----")
	for _, msgHash := range m.msgHashes[height][round][algorithm.Precommit] {
		// spew.Dump(m.rawMessages[msgHash].decodedMsg)
		if value == m.messages[msgHash].consensusMessage.Value {
			sigs = append(sigs, m.messages[msgHash].signature)
		}
	}
	// println("----------------")
	return sigs
}

func (m *messageStore) prevoteQuorum(valueHash *common.Hash, round int64, header *types.Header) bool {
	msgType := new(algorithm.Step)
	*msgType = algorithm.Prevote
	// println("prevote power --------")
	vp := m.votePower(valueHash, round, msgType, header) >= header.Committee.Quorum()
	// println(vp, "----------------")
	return vp
}

func (m *messageStore) precommitQuorum(valueHash *common.Hash, round int64, header *types.Header) bool {
	msgType := new(algorithm.Step)
	*msgType = algorithm.Precommit
	// // println("precommit power", m.votePower(valueHash, round, msgType, header))
	// println("precommit power --------")
	vp := m.votePower(valueHash, round, msgType, header) >= header.Committee.Quorum()
	// println(vp, "----------------")
	return vp
}

func (m *messageStore) fail(round int64, header *types.Header) bool {
	// println("fail power --------")
	vp := m.votePower(nil, round, nil, header) >= header.Committee.Quorum()
	// println(vp, "----------------")
	return vp
}

func (m *messageStore) votePower(
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
			algorithm.Prevote,
			algorithm.Precommit,
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
			if valueHash != nil && *valueHash != common.Hash(m.messages[msgHash].consensusMessage.Value) {
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

func (m *messageStore) addMessage(msg *message, rawMsg []byte) error {
	// Check message is in bounds
	if !m.messageBounds.in(msg.consensusMessage.Height) {
		return fmt.Errorf("message %v out of bounds", msg.String())
	}
	// Check we haven't already processed this message
	if m.Message(msg.hash) != nil {
		// Message was already processed
		return fmt.Errorf("message %v already processed", msg.String())
	}
	err := addMsgHash(m.msgHashes, msg.consensusMessage.Height, msg.consensusMessage.Round, msg.consensusMessage.MsgType, msg.address, msg.hash)
	if err != nil {
		return err
	}
	// // println(id, "hashes len", len(m.msgHashes))
	m.messages[msg.hash] = msg
	m.rawMessages[msg.hash] = rawMsg
	return nil
}

// removeMessage removes a single message, it does not handle deleting empty
// maps after message removal, that will be handled when deleting whole heights
// due to height changes.
func (m *messageStore) removeMessage(msg *message) {
	// Delete entry in hashes
	delete(m.msgHashes[msg.consensusMessage.Height][msg.consensusMessage.Round][msg.consensusMessage.MsgType], msg.address)
	// Delete entry in messages
	delete(m.messages, msg.hash)
	// Delete entry in rawMessages
	delete(m.rawMessages, msg.hash)
	delete(m.valid, msg.hash)

	if msg.consensusMessage.MsgType == algorithm.Propose {
		valueHash := msg.value.Hash()
		delete(m.values, valueHash)
		delete(m.valid, valueHash)
	} else {
		delete(m.valid, msg.hash)
	}
}

// setHeight updates the height in the bounds and removes messages that are now out of bounds.
func (m *messageStore) setHeight(height uint64) {
	low, high := m.messageBounds.setCentre(height)
	for i := low; i < high; i++ {
		m.removeMessagesAtHeight(i)
	}
}

func (m *messageStore) removeMessagesAtHeight(height uint64) {
	// Remove all messgages at this height
	for _, msgTypeMap := range m.msgHashes[height] {
		for _, addressMap := range msgTypeMap {
			for _, hash := range addressMap {
				m.removeMessage(m.messages[hash])
			}
		}
	}
	// Delete map entry for this height
	delete(m.msgHashes, height)
}

func (m *messageStore) addValue(valueHash common.Hash, value *types.Block) {
	m.values[valueHash] = value
}
func (m *messageStore) value(valueHash common.Hash) *types.Block {
	return m.values[valueHash]
}

// Mark the hash of something valid, it could be a message hash or a value hash
func (m *messageStore) setValid(itemHash common.Hash) {
	m.valid[itemHash] = struct{}{}
}
func (m *messageStore) isValid(itemHash common.Hash) bool {
	_, ok := m.valid[itemHash]
	return ok
}

func (m *messageStore) matchingProposal(cm *algorithm.ConsensusMessage) *message {
	for _, proposalHash := range m.msgHashes[cm.Height][cm.Round][algorithm.Propose] {
		proposal := m.messages[proposalHash]
		if proposal.consensusMessage.Value == cm.Value {
			// fmt.Printf(" got: %s\n", proposal.String())
			return proposal
		}
	}
	// fmt.Printf(" no proposal\n")
	return nil
}

type bounds struct {
	centre uint64
	high   uint64
	low    uint64
}

// setCentre sets the centre value for the bounds, it retruns a range of all
// values that were in the range and are now not in the range [a,b) with a
// inclusive and b exclusive.
func (b *bounds) setCentre(v uint64) (uint64, uint64) {
	l := b.lower()
	h := b.upper()

	b.centre = v

	newl := b.lower()
	newh := b.upper()

	switch {
	case h < newl:
		return l, h + 1
	case l > newh:
		return l, h + 1
	case h < newh:
		return l, newl
	case h > newh:
		return newh, h
	default: // The bounds have not changed
		return 0, 0
	}
}

func (b *bounds) in(v uint64) bool {
	return v >= b.lower() && v <= b.upper()
}

func (b *bounds) upper() uint64 {
	return b.centre + b.high
}

func (b *bounds) lower() uint64 {
	if b.low >= b.centre {
		return 0
	}
	return b.centre - b.low
}
