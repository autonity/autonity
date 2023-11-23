package types

import (
	"bytes"
	"errors"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/crypto"
	"github.com/autonity/autonity/crypto/bls"
	"github.com/autonity/autonity/log"
	"github.com/autonity/autonity/rlp"
	lru "github.com/hashicorp/golang-lru"
	"golang.org/x/crypto/sha3"
	"io"
	"math/big"
	"sort"
	"sync"
)

var (
	// BFTDigest represents a hash of "Tendermint byzantine fault tolerance"
	// to identify whether the block is from BFT consensus engine
	//TODO differentiate the digest between IBFT and Tendermint
	BFTDigest = common.HexToHash("0x63746963616c2062797a616e74696e65206661756c7420746f6c6572616e6365")

	BFTExtraVanity = 32 // Fixed number of extra-data bytes reserved for validator vanity
	BFTExtraSeal   = 65 // Fixed number of extra-data bytes reserved for validator seal

	inmemoryAddresses  = 500 // Number of recent addresses from ecrecover
	recentAddresses, _ = lru.NewARC(inmemoryAddresses)

	// ErrInvalidBFTHeaderExtra is returned if the length of extra-data is less than 32 bytes
	ErrInvalidBFTHeaderExtra = errors.New("invalid pos header extra-data")
	// ErrInvalidSignature is returned when given signature is not signed by given address.
	ErrInvalidSignature = errors.New("invalid signature")
	// ErrInvalidCommittedSeals is returned if the committed seal is not signed by any of parent validators.
	ErrInvalidCommittedSeals = errors.New("invalid committed seals")
	// ErrEmptyCommittedSeals is returned if the field of committed seals is zero.
	ErrEmptyCommittedSeals = errors.New("zero committed seals")
	// ErrNegativeRound is returned if the round field is negative
	ErrNegativeRound = errors.New("negative round")
)

// BFTFilteredHeader returns a filtered header which some information (like seal, committed seals)
// are clean to fulfill the BFT hash rules.
func BFTFilteredHeader(h *Header, keepSeal bool) *Header {
	newHeader := CopyHeader(h)
	if !keepSeal {
		newHeader.ProposerSeal = []byte{}
	}
	newHeader.CommittedSeals = [][]byte{}
	newHeader.Round = 0
	newHeader.Extra = []byte{}
	return newHeader
}

// sigHash returns the hash which is used as input for the BFT
// signing. It is the hash of the entire header apart from the 65 byte signature
// contained at the end of the extra data.
//
// Note, the method requires the extra data to be at least 65 bytes, otherwise it
// panics. This is done to avoid accidentally using both forms (signature present
// or not), which could be abused to produce different hashes for the same header.
func SigHash(header *Header) (hash common.Hash) {
	hasher := sha3.NewLegacyKeccak256()

	// Clean seal is required for calculating proposer seal.
	err := rlp.Encode(hasher, BFTFilteredHeader(header, false))
	if err != nil {
		log.Error("can't hash the header", "err", err, "header", header)
		return common.Hash{}
	}
	hasher.Sum(hash[:0])
	return hash
}

// Ecrecover extracts the Ethereum account address from a signed header.
func Ecrecover(header *Header) (common.Address, error) {
	hash := header.Hash()
	if addr, ok := recentAddresses.Get(hash); ok {
		return addr.(common.Address), nil
	}

	addr, err := GetSignatureAddress(SigHash(header).Bytes(), header.ProposerSeal)
	if err != nil {
		return addr, err
	}
	recentAddresses.Add(hash, addr)
	return addr, nil
}

// WriteSeal writes the extra-data field of the given header with the given seals.
// suggest to rename to writeSeal.
func WriteSeal(h *Header, seal []byte) error {
	if len(seal) != BFTExtraSeal {
		return ErrInvalidSignature
	}
	h.ProposerSeal = make([]byte, len(seal))
	copy(h.ProposerSeal, seal)
	return nil
}

// WriteRound writes the round field of the block header.
func WriteRound(h *Header, round int64) error {
	if round < 0 {
		return ErrNegativeRound
	}
	h.Round = uint64(round)
	return nil
}

// WriteCommittedSeals writes the extra-data field of a block header with given committed seals.
func WriteCommittedSeals(h *Header, committedSeals [][]byte) error {
	if len(committedSeals) == 0 {
		return ErrInvalidCommittedSeals
	}

	for _, seal := range committedSeals {
		if len(seal) != BFTExtraSeal {
			return ErrInvalidCommittedSeals
		}
	}

	h.CommittedSeals = make([][]byte, len(committedSeals))
	for i, val := range committedSeals {
		h.CommittedSeals[i] = make([]byte, len(val))
		copy(h.CommittedSeals[i], val)
	}

	return nil
}

func RLPHash(v any) (h common.Hash) {
	hw := sha3.NewLegacyKeccak256()
	rlp.Encode(hw, v)
	hw.Sum(h[:0])
	return h
}

// GetSignatureAddress gets the signer address from the signature
// Youssef: This is too generic and merely a wrapper against ecrecover, that belongs to the crypto package.
func GetSignatureAddress(data []byte, sig []byte) (common.Address, error) {
	// 1. Keccak data
	hashData := crypto.Keccak256(data)
	// 2. Recover public key
	pubkey, err := crypto.SigToPub(hashData, sig)
	if err != nil {
		return common.Address{}, err
	}
	return crypto.PubkeyToAddress(*pubkey), nil
}

type CommitteeMember struct {
	Address      common.Address `json:"address"            gencodec:"required"       abi:"addr"`
	VotingPower  *big.Int       `json:"votingPower"        gencodec:"required"       abi:"votingPower"`
	ValidatorKey []byte         `json:"validatorKey"       gencodec:"required"       abi:"validatorKey"`
}

func (c *CommitteeMember) Bytes() []byte { //nolint
	return c.Address.Bytes()
}

func (c *CommitteeMember) String() string {
	return c.Address.String()
}

func (c *CommitteeMember) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []any{c.Address, c.VotingPower, c.ValidatorKey})
}

func (c *CommitteeMember) DecodeRLP(s *rlp.Stream) error {
	var cm struct {
		Address      common.Address
		VotingPower  *big.Int
		ValidatorKey []byte
	}

	if err := s.Decode(&cm); err != nil {
		return err
	}

	c.Address, c.VotingPower, c.ValidatorKey = cm.Address, cm.VotingPower, cm.ValidatorKey
	return nil
}

type Committee struct {
	Members          []*CommitteeMember `json:"members"`
	totalVotingPower *big.Int
	power            sync.Once
	// aggregated validator public key of all committee members
	aggValidatorKey []byte
	aggregation     sync.Once
	// used for committee member lookup
	membersMap map[common.Address]*CommitteeMember
	members    sync.Once
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

func (c *Committee) Less(i, j int) bool {
	return bytes.Compare(c.Members[i].Bytes(), c.Members[j].Bytes()) < 0
}

func (c *Committee) Swap(i, j int) {
	c.Members[i], c.Members[j] = c.Members[j], c.Members[i]
}

// EncodeRLP To make the header hash deterministic, unify the rlp hash of committee member by excluding the cached values.
func (c *Committee) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []any{c.Members})
}

// DecodeRLP decodes the Ethereum
func (c *Committee) DecodeRLP(s *rlp.Stream) error {
	var committee struct {
		Members []*CommitteeMember
	}

	if err := s.Decode(&committee); err != nil {
		return err
	}

	c.Members = committee.Members
	return nil
}

func (c *Committee) CommitteeMember(address common.Address) *CommitteeMember { // nolint
	// build the map only once
	c.members.Do(func() {
		log.Debug("Creating committeeMap")
		c.membersMap = make(map[common.Address]*CommitteeMember)
		for _, member := range c.Members {
			c.membersMap[member.Address] = member
		}
	})
	return c.membersMap[address]
}

func (c *Committee) AggregatedValidatorKey() []byte { // nolint
	// compute aggregated key only once, then returned cached value
	c.aggregation.Do(func() {
		// collect bls keys of committee members in bytes
		rawKeys := make([][]byte, len(c.Members))
		for index, member := range c.Members {
			rawKeys[index] = member.ValidatorKey
		}

		aggKey, err := bls.AggregatePublicKeys(rawKeys)
		if err != nil {
			// this shouldn't happen as all validator keys are validity checked before they are inserted into the AC
			log.Crit("invalid BLS key in header")
		}
		c.aggValidatorKey = aggKey.Marshal()
	})
	return c.aggValidatorKey
}

func (c *Committee) Sort() {
	if len(c.Members) != 0 {
		// sort the slice by address
		sort.SliceStable(c.Members, func(i, j int) bool {
			return bytes.Compare(c.Members[i].Address.Bytes(), c.Members[j].Address.Bytes()) < 0
		})
	}
}

// TotalVotingPower returns the total voting power contained in the committee.
func (c *Committee) TotalVotingPower() *big.Int {
	// compute power only once, then returned cached value
	c.power.Do(func() {
		total := new(big.Int)
		for _, m := range c.Members {
			total.Add(total, m.VotingPower)
		}
		c.totalVotingPower = total
	})

	// return a copy of the cached value to prevent un-expected modification of the cached value.
	return new(big.Int).Set(c.totalVotingPower)
}
