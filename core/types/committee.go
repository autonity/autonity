package types

import (
	"bytes"
	"fmt"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/common/hexutil"
	"github.com/autonity/autonity/crypto/blst"
	"math/big"
	"slices"
	"sync"
)

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
	// todo: implement helpers for this new fields.
	// ParentEpochBlock points to the parent epoch block number for backward searching of epoch header.
	ParentEpochBlock *big.Int `json:"parentEpochBlock"`
	// NextEpochBlock points to the next epoch block number for forward searching of epoch header.
	NextEpochBlock *big.Int `json:"nextEpochBlock"`
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
	var clone = &Committee{}
	if c.Members != nil {
		clone.Members = make([]CommitteeMember, len(c.Members))
		for i, val := range c.Members {
			clone.Members[i] = CommitteeMember{
				Address:      val.Address,
				VotingPower:  new(big.Int).Set(val.VotingPower),
				Index:        val.Index,
				ConsensusKey: val.ConsensusKey.Copy(),
			}
			clone.Members[i].ConsensusKeyBytes = make([]byte, len(val.ConsensusKeyBytes))
			copy(clone.Members[i].ConsensusKeyBytes, val.ConsensusKeyBytes)
		}
	}

	if c.ParentEpochBlock != nil {
		clone.ParentEpochBlock = new(big.Int).SetUint64(c.ParentEpochBlock.Uint64())
	}
	if c.NextEpochBlock != nil {
		clone.NextEpochBlock = new(big.Int).SetUint64(c.NextEpochBlock.Uint64())
	}
	c.lock.RLock()
	defer c.lock.RUnlock()
	if c.totalVotingPower != nil {
		clone.totalVotingPower = new(big.Int).Set(c.totalVotingPower)
	}

	if c.membersMap != nil {
		clone.membersMap = make(map[common.Address]*CommitteeMember)
		for k, v := range c.membersMap {
			clone.membersMap[k] = v
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

	if (c.ParentEpochBlock == nil && other.ParentEpochBlock != nil) ||
		(c.ParentEpochBlock != nil && other.ParentEpochBlock == nil) ||
		(c.ParentEpochBlock != nil && other.ParentEpochBlock != nil && c.ParentEpochBlock.Cmp(other.ParentEpochBlock) != 0) {
		return false
	}

	if (c.NextEpochBlock == nil && other.NextEpochBlock != nil) ||
		(c.NextEpochBlock != nil && other.NextEpochBlock == nil) ||
		(c.NextEpochBlock != nil && other.NextEpochBlock != nil && c.NextEpochBlock.Cmp(other.NextEpochBlock) != 0) {
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
	c.lock.RLock()
	defer c.lock.RUnlock()
	if index >= len(c.Members) || index < 0 {
		return nil
	}
	return &c.Members[index]
}

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
}

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
}

func (c *Committee) Sort() {
	// Todo: (Jason) since the compute committee precompile contract did the sorting, thus we might not need this
	//  anymore, however some of the unit test need this sorting.
	if len(c.Members) != 0 {
		// sort validators according to their voting power in descending order
		// stable sort keeps the original order of equal elements
		slices.SortStableFunc(c.Members, func(a, b CommitteeMember) int {
			return b.VotingPower.Cmp(a.VotingPower)
		})
	}
}

// Enrich adds some convenience information to the committee member structs
func (c *Committee) Enrich() error {
	for i := range c.Members {
		consensusKey, err := blst.PublicKeyFromBytes(c.Members[i].ConsensusKeyBytes)
		if err != nil {
			return fmt.Errorf("Error when decoding bls key: index %d, address: %v, err: %w", i, c.Members[i].Address, err)
		}
		c.Members[i].ConsensusKey = consensusKey
		c.Members[i].Index = uint64(i)
	}
	return nil
}
