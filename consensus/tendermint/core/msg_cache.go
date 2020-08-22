package core

import (
	common "github.com/clearmatics/autonity/common"
	types "github.com/clearmatics/autonity/core/types"
)

// Message cache caches messages

type messageCache struct {
	proposalMsgHashes    map[uint64]map[int64]map[common.Address]common.Hash
	msgHashToProposal    map[common.Hash]*Proposal
	valueHashToProposals map[common.Hash]map[common.Address]*Proposal

	prevoteMsgHashes    map[uint64]map[int64]map[common.Address]common.Hash
	msgHashToPrevote    map[common.Hash]*Vote
	valueHashToPrevotes map[common.Hash]map[common.Address]*Vote

	precommitMsgHashes    map[uint64]map[int64]map[common.Address]common.Hash
	msgHashToPrecommit    map[common.Hash]*Vote
	valueHashToSignatures map[common.Hash]map[common.Address][]byte

	rawMessages map[common.Hash]*Message

	// TODO use this to impose upper and lower limits on what can go in the
	// message cache.
	currentHeight uint64
}

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

func (m *messageCache) addProposal(p *Proposal, msg *Message) error {
	err := addMsgHash(m.proposalMsgHashes, p.Height.Uint64(), p.Round, msg.Address, msg.Hash)
	if err != nil {
		return err
	}
	m.msgHashToProposal[msg.Hash] = p
	m.rawMessages[msg.Hash] = msg
	return nil
}

func (m *messageCache) addPrevote(p *Vote, msg *Message) error {
	err := addMsgHash(m.prevoteMsgHashes, p.Height.Uint64(), p.Round, msg.Address, msg.Hash)
	if err != nil {
		return err
	}
	m.msgHashToPrevote[msg.Hash] = p
	m.rawMessages[msg.Hash] = msg
	addVoteForValue(m.valueHashToPrevotes, p.ProposedBlockHash, msg.Address, p)
	return nil
}

func (m *messageCache) addPrecommit(p *Vote, msg *Message) error {
	err := addMsgHash(m.precommitMsgHashes, p.Height.Uint64(), p.Round, msg.Address, msg.Hash)
	if err != nil {
		return err
	}
	m.msgHashToPrecommit[msg.Hash] = p
	m.rawMessages[msg.Hash] = msg

	addressMap, ok := m.valueHashToSignatures[p.ProposedBlockHash]
	if !ok {
		addressMap = make(map[common.Address][]byte)
	}
	addressMap[msg.Address] = msg.CommittedSeal
	return nil
}

func (m *messageCache) Message(h common.Hash) *Message {
	return m.rawMessages[h]
}

func (m *messageCache) prevotePower(valueHash common.Hash, header *types.Header) uint64 {
	return votePower(m.valueHashToPrevotes, valueHash, header)
}

func (m *messageCache) totalPrevotePower(header *types.Header, round int64) uint64 {
	// get all messages for a round and height
	// iterate then and add the power
	return 0
}

func (m *messageCache) precommitPower(valueHash common.Hash, header *types.Header) uint64 {
	// Temp figure out what i need here
	// return votePower(m.valueHashToPrecommits, valueHash, header)

	return 0
}

func (m *messageCache) signatures(valueHash common.Hash) [][]byte {
	signaturesByAddress := m.valueHashToSignatures[valueHash]
	// Find all the signatures for this value
	sigs := make([][]byte, len(signaturesByAddress))
	for _, sig := range signaturesByAddress {
		sigs = append(sigs, sig)
	}
	return sigs
}

func votePower(valueHashToVotes map[common.Hash]map[common.Address]*Vote, valueHash common.Hash, header *types.Header) uint64 {
	votesByAddress := valueHashToVotes[valueHash]
	// Total the power of all votes for this value, failure to find a committee
	// member in the header indicates a programming error and an invalid memory
	// acccess panic will ensue.
	var power uint64
	for address := range votesByAddress {
		power += header.CommitteeMember(address).VotingPower.Uint64()
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

func addMsgHash(hashes map[uint64]map[int64]map[common.Address]common.Hash, height uint64, round int64, address common.Address, hash common.Hash) error {
	// todo check bounds
	roundMap, ok := hashes[height]
	if !ok {
		roundMap = make(map[int64]map[common.Address]common.Hash)
	}
	addressMap, ok := roundMap[round]
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
