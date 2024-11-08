package types

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"slices"
	"sync"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/common/hexutil"
	"github.com/autonity/autonity/consensus/tendermint/core/constants"
	"github.com/autonity/autonity/crypto"
	"github.com/autonity/autonity/crypto/blst"
)

var _ = (*Epoch)(nil)

// EpochInfo is used to cache the epoch info in the blockchain context for the protocol, it not only have the raw
// epoch but also include the epoch block number of an epoch, thus that we don't need to cache the entire epoch block
// headers.
type EpochInfo struct {
	Epoch
	EpochBlock *big.Int
}

// Epoch contains the previous epoch block, next epoch block and its committee of the current epoch.
// It is saved in the block header if current block is an epoch block. The epoch block number is not saved
// here since it is duplicated than the block number in the block header.
type Epoch struct {
	PreviousEpochBlock *big.Int   `rlp:"nil" json:"previousEpochBlock" gencodec:"required"`
	NextEpochBlock     *big.Int   `rlp:"nil" json:"nextEpochBlock" gencodec:"required"`
	Committee          *Committee `rlp:"nil" json:"committee" gencodec:"required"`
	Delta              *big.Int   `rlp:"nil" json:"delta" gencodec:"required"` // delta for omission failure
}

// MarshalJSON marshals as JSON.
func (e Epoch) MarshalJSON() ([]byte, error) {
	type epoch struct {
		PreviousEpochBlock *hexutil.Big `json:"previousEpochBlock" gencodec:"required"`
		NextEpochBlock     *hexutil.Big `json:"nextEpochBlock" gencodec:"required"`
		Committee          Committee    `json:"committee" gencodec:"required"`
		Delta              *hexutil.Big `json:"delta" gencodec:"required"`
	}
	var enc epoch
	enc.PreviousEpochBlock = (*hexutil.Big)(e.PreviousEpochBlock)
	enc.NextEpochBlock = (*hexutil.Big)(e.NextEpochBlock)
	enc.Committee = Committee{}
	if e.Committee != nil {
		enc.Committee = *e.Committee // nolint
	}
	enc.Delta = (*hexutil.Big)(e.Delta)
	return json.Marshal(&enc)
}

// UnmarshalJSON unmarshals from JSON.
func (e *Epoch) UnmarshalJSON(input []byte) error {
	type epoch struct {
		PreviousEpochBlock *hexutil.Big `json:"previousEpochBlock" gencodec:"required"`
		NextEpochBlock     *hexutil.Big `json:"nextEpochBlock" gencodec:"required"`
		Committee          *Committee   `json:"committee" gencodec:"required"`
		Delta              *hexutil.Big `json:"delta" gencodec:"required"`
	}

	var dec epoch
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}

	if dec.PreviousEpochBlock == nil {
		return errors.New("missing required field 'previousEpochBlock' for epoch")
	}
	e.PreviousEpochBlock = dec.PreviousEpochBlock.ToInt()

	if dec.NextEpochBlock == nil {
		return errors.New("missing required field 'nextEpochBlock' for epoch")
	}
	e.NextEpochBlock = dec.NextEpochBlock.ToInt()

	if dec.Committee == nil {
		return errors.New("missing required field 'committee' for epoch")
	}
	e.Committee = dec.Committee

	if dec.Delta == nil {
		return errors.New("missing required field 'delta' for epoch")
	}
	e.Delta = dec.Delta.ToInt()
	return nil
}

func (e *Epoch) Copy() *Epoch {
	cpy := &Epoch{}

	if e.PreviousEpochBlock != nil {
		cpy.PreviousEpochBlock = new(big.Int).Set(e.PreviousEpochBlock)
	}

	if e.NextEpochBlock != nil {
		cpy.NextEpochBlock = new(big.Int).Set(e.NextEpochBlock)
	}

	if e.Committee != nil {
		cpy.Committee = e.Committee.Copy()
	}

	if e.Delta != nil {
		cpy.Delta = new(big.Int).Set(e.Delta)
	}
	return cpy
}

func (e *Epoch) Equal(other *Epoch) bool {
	if e == nil && other == nil {
		return true
	}

	if e == nil || other == nil {
		return false
	}

	// as all the fields of epoch were checked on the rlp decoding phase,
	// thus they cannot be nil here.
	if e.PreviousEpochBlock.Cmp(other.PreviousEpochBlock) != 0 {
		return false
	}

	if e.NextEpochBlock.Cmp(other.NextEpochBlock) != 0 {
		return false
	}

	if !e.Committee.Equal(other.Committee) {
		return false
	}

	if e.Delta.Cmp(other.Delta) != 0 {
		return false
	}

	return true
}

//go:generate gencodec -type CommitteeMember -field-override committeeMemberMarshaling -out gen_member_json.go

type CommitteeMember struct {
	Address           common.Address `json:"address"            gencodec:"required"       abi:"addr"`
	VotingPower       *big.Int       `json:"votingPower"        gencodec:"required"`
	ConsensusKeyBytes []byte         `json:"consensusKey"       gencodec:"required"       abi:"consensusKey"`
	// this field is ignored when rlp/json encoding/decoding, it is computed locally from the bytes
	ConsensusKey blst.PublicKey `json:"-" rlp:"-"`
	Index        uint64         `json:"-" rlp:"-"` // index of this committee member in the committee array
}

func (c *CommitteeMember) Bytes() []byte {
	return c.Address.Bytes()
}

func (c *CommitteeMember) String() string {
	return c.Address.String()
}

type committeeMemberMarshaling struct {
	Address           common.Address
	VotingPower       *hexutil.Big
	ConsensusKeyBytes hexutil.Bytes
}

type Committee struct {
	// Current committee members.
	Members []CommitteeMember `json:"members"`
	// this field is ignored when rlp/json encoding/decoding, it is computed locally from the bytes
	// mutex to protect internal cached fields of committee from race condition.
	lock sync.RWMutex `json:"-" rlp:"-"`
	// cached total voting power.
	totalVotingPower *big.Int `json:"-" rlp:"-"`
	// cached indexing of committee for member lookup
	membersMap map[common.Address]*CommitteeMember `json:"-" rlp:"-"`
}

func (c *Committee) String() string {
	var ret string
	for _, val := range c.Members {
		ret += "[" + val.Address.String() + " - " + val.VotingPower.String() + "] "
	}
	return ret
}

func (c *Committee) Len() int {
	return len(c.Members)
}

func (c *Committee) Copy() *Committee {
	if c == nil {
		return nil
	}
	var clone = &Committee{}
	if c.Members != nil {
		clone.Members = make([]CommitteeMember, len(c.Members))
		for i, val := range c.Members {
			clone.Members[i] = CommitteeMember{
				Address:     val.Address,
				VotingPower: new(big.Int).Set(val.VotingPower),
				Index:       val.Index,
			}
			clone.Members[i].ConsensusKeyBytes = make([]byte, len(val.ConsensusKeyBytes))
			copy(clone.Members[i].ConsensusKeyBytes, val.ConsensusKeyBytes)
			if val.ConsensusKey != nil {
				clone.Members[i].ConsensusKey = val.ConsensusKey.Copy()
			}
		}
	}

	c.lock.RLock()
	defer c.lock.RUnlock()
	if c.totalVotingPower != nil {
		clone.totalVotingPower = new(big.Int).Set(c.totalVotingPower)
	}

	if c.membersMap != nil {
		clone.membersMap = make(map[common.Address]*CommitteeMember)
		for _, v := range clone.Members {
			member := v
			clone.membersMap[v.Address] = &member
		}
	}
	return clone
}

func (c *Committee) Equal(other *Committee) bool {
	if c == nil && other == nil {
		return true
	}
	if c == nil || other == nil {
		return false
	}

	if len(c.Members) != len(other.Members) {
		return false
	}
	for i := range c.Members {
		if c.Members[i].Address != other.Members[i].Address ||
			c.Members[i].VotingPower.Cmp(other.Members[i].VotingPower) != 0 ||
			!bytes.Equal(c.Members[i].ConsensusKeyBytes, other.Members[i].ConsensusKeyBytes) ||
			!bytes.Equal(c.Members[i].ConsensusKey.Marshal(), other.Members[i].ConsensusKey.Marshal()) ||
			c.Members[i].Index != other.Members[i].Index {
			return false
		}
	}
	return true
}

func (c *Committee) MemberByIndex(index int) *CommitteeMember {
	if index >= len(c.Members) || index < 0 {
		return nil
	}
	return &c.Members[index]
}

func (c *Committee) MemberByAddress(address common.Address) *CommitteeMember {
	c.lock.RLock()
	if c.membersMap != nil {
		defer c.lock.RUnlock()
		return c.membersMap[address]
	}
	c.lock.RUnlock() // Release read lock before acquiring write lock
	c.lock.Lock()
	defer c.lock.Unlock()
	// Double-check if the map was initialized while waiting for the lock
	if c.membersMap == nil {
		c.membersMap = make(map[common.Address]*CommitteeMember)
		for i := range c.Members {
			c.membersMap[c.Members[i].Address] = &c.Members[i]
		}
	}
	return c.membersMap[address]
}

/*
func (c *Committee) MemberByAddress(address common.Address) *CommitteeMember {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.membersMap == nil {
		c.membersMap = make(map[common.Address]*CommitteeMember)
		for _, member := range c.Members {
			m := member
			c.membersMap[member.Address] = &m
		}
	}
	return c.membersMap[address]
}*/

func (c *Committee) TotalVotingPower() *big.Int {
	c.lock.RLock() // Acquire read lock
	if c.totalVotingPower != nil {
		defer c.lock.RUnlock()                      // Release read lock
		return new(big.Int).Set(c.totalVotingPower) // Return a copy of the cached value
	}
	c.lock.RUnlock() // Release read lock before acquiring write lock
	c.lock.Lock()    // Acquire write lock
	defer c.lock.Unlock()

	// Double-check if the value was initialized while waiting for the lock
	if c.totalVotingPower == nil {
		total := new(big.Int)
		for _, m := range c.Members {
			total.Add(total, m.VotingPower)
		}
		c.totalVotingPower = total
	}

	// Return a copy of the cached value
	return new(big.Int).Set(c.totalVotingPower)
}

/*
// TotalVotingPower returns the total voting power contained in the committee.
func (c *Committee) TotalVotingPower() *big.Int {
	c.lock.Lock()
	defer c.lock.Unlock()
	// compute power only once, then returned cached value
	if c.totalVotingPower == nil {
		total := new(big.Int)
		for _, m := range c.Members {
			total.Add(total, m.VotingPower)
		}
		c.totalVotingPower = total
	}

	// return a copy of the cached value to prevent un-expected modification of the cached value.
	return new(big.Int).Set(c.totalVotingPower)
}*/

// Enrich adds some convenience information to the committee member structs
func (c *Committee) Enrich() error {
	for i := range c.Members {
		consensusKey, err := blst.PublicKeyFromBytes(c.Members[i].ConsensusKeyBytes)
		if err != nil {
			return fmt.Errorf("error when decoding bls key: index %d, address: %v, err: %w", i, c.Members[i].Address, err)
		}
		c.Members[i].ConsensusKey = consensusKey
		c.Members[i].Index = uint64(i)
	}
	return nil
}

// Proposer elects a proposer from the voting power sorted committee for the height and round, it implements same
// election algorithm as proposer election in AC contract. The underlying sorted committee is sourced from the AC contract.
func (c *Committee) Proposer(height uint64, round int64) common.Address {
	totalVotingPower := c.TotalVotingPower()
	seed := big.NewInt(constants.ProposerElectionCycleConstant)
	// for power weighted sampling, we distribute seed into a 256bits key-space, and compute the hit index.
	h := new(big.Int).SetUint64(height)
	r := new(big.Int).SetInt64(round)
	key := r.Add(r, h.Mul(h, seed))
	value := new(big.Int).SetBytes(crypto.Keccak256(key.Bytes()))
	index := value.Mod(value, totalVotingPower)
	// find the index hit which committee member which line up in the committee list.
	// we assume there is no 0 stake/power validators.
	counter := new(big.Int).SetUint64(0)
	c.lock.RLock()
	defer c.lock.RUnlock()
	for _, member := range c.Members {
		counter.Add(counter, member.VotingPower)
		if index.Cmp(counter) == -1 {
			return member.Address
		}
	}

	// otherwise, we elect with round-robin.
	return c.Members[round%int64(len(c.Members))].Address
}

// SortCommitteeMembers sort validators according to their voting power in descending order
// stable sort keeps the original order of equal elements
func SortCommitteeMembers(members []CommitteeMember) {
	slices.SortStableFunc(members, func(a, b CommitteeMember) int {
		return b.VotingPower.Cmp(a.VotingPower)
	})
}
