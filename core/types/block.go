// Copyright 2014 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

// Package types contains data types related to Ethereum consensus.
package types

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/autonity/autonity/consensus/tendermint/bft"
	"github.com/autonity/autonity/crypto"
	"io"
	"math/big"
	"reflect"
	"sync/atomic"
	"time"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/common/hexutil"
	"github.com/autonity/autonity/crypto/blst"
	"github.com/autonity/autonity/rlp"
)

var (
	EmptyRootHash                = common.HexToHash("56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421")
	EmptyUncleHash               = rlpHash([]*Header(nil))
	errInvalidSignature          = errors.New("aggregate signature is invalid")
	ErrNonAggregatablePublicKeys = errors.New("provided public keys cannot be aggregated")
	errNoQuorum                  = errors.New("aggregate signature does not contain quorum voting power")
)

const (
	/*
	* We currently have in HeaderExtra:
	* - ProposerSeal (ECDSA signature) = 65 bytes
	* - Round (uint64) = 8 byte
	* - ActivityProofRound (uint64) = 8 byte
	* - QuorumCertificate (AggregateSignature) = BlsSignature + Signers struct = 96 bytes + n/4 bytes + 2n bytes = 9/4n + 96 bytes
	*	with n being the committee size. This is the absolute worst case scenario, which I believe will rarely happen
	* - Epoch = 3 *big.Int + committee = 3 * 32 + CommitteeMemberSize * n = 96 + (20+32+48) * n = 96 + 100n
	* - ActivityProof, same as QuorumCertificate = 9/4n + 96 bytes
	*
	* Assuming n=1000, the worst case scenario total size of extra data will be:
	* 65+8+8+96+9/4n+96+100n+9/4n+96 = 104.5n + 369 ~= 103kb
	 */
	maximumHeaderExtraLength = 110 * 1024 //110kb
	maximumDifficultyBitlen  = 80
	maximumBaseFeeBitlen     = 256
	maximumVotingPowerBitlen = 256
)

// A BlockNonce is a 64-bit hash which proves (combined with the
// mix-hash) that a sufficient amount of computation has been carried
// out on a block.
type BlockNonce [8]byte

// EncodeNonce converts the given integer to a block nonce.
func EncodeNonce(i uint64) BlockNonce {
	var n BlockNonce
	binary.BigEndian.PutUint64(n[:], i)
	return n
}

// Uint64 returns the integer value of a block nonce.
func (n BlockNonce) Uint64() uint64 {
	return binary.BigEndian.Uint64(n[:])
}

// MarshalText encodes n as a hex string with 0x prefix.
func (n BlockNonce) MarshalText() ([]byte, error) {
	return hexutil.Bytes(n[:]).MarshalText()
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (n *BlockNonce) UnmarshalText(input []byte) error {
	return hexutil.UnmarshalFixedText("BlockNonce", input, n[:])
}

//go:generate gencodec -type Header -field-override headerMarshaling -out gen_header_json.go

// Header represents a block header in the Autonity blockchain.
type Header struct {
	// NOTE: HeaderParentHashFromRLP relies on ParentHash being the first element of the struct. Do not move it.
	ParentHash  common.Hash    `json:"parentHash"       gencodec:"required"`
	UncleHash   common.Hash    `json:"sha3Uncles"       gencodec:"required"`
	Coinbase    common.Address `json:"miner"            gencodec:"required"`
	Root        common.Hash    `json:"stateRoot"        gencodec:"required"`
	TxHash      common.Hash    `json:"transactionsRoot" gencodec:"required"`
	ReceiptHash common.Hash    `json:"receiptsRoot"     gencodec:"required"`
	Bloom       Bloom          `json:"logsBloom"        gencodec:"required"`
	Difficulty  *big.Int       `json:"difficulty"       gencodec:"required"`
	Number      *big.Int       `json:"number"           gencodec:"required"`
	GasLimit    uint64         `json:"gasLimit"         gencodec:"required"`
	GasUsed     uint64         `json:"gasUsed"          gencodec:"required"`
	Time        uint64         `json:"timestamp"        gencodec:"required"`
	Extra       []byte         `json:"extraData"        gencodec:"required"`
	MixDigest   common.Hash    `json:"mixHash"`
	Nonce       BlockNonce     `json:"nonce"`
	BaseFee     *big.Int       `json:"baseFeePerGas" gencodec:"required"`

	// autonity custom fields
	ProposerSeal       []byte `json:"proposerSeal"        gencodec:"required"`
	Round              uint64 `json:"round"               gencodec:"required"`
	ActivityProofRound uint64 `json:"activityProofRound"  gencodec:"required"`
	// the following fields will be left nil if they were nil when encoded.
	// see the headerExtra struct for rlp:nil tags and the custom decodeRLP method for Header. This is where sanity checks are done.
	QuorumCertificate *AggregateSignature `json:"quorumCertificate"   gencodec:"optional"`
	Epoch             *Epoch              `json:"epoch"               gencodec:"optional"`
	ActivityProof     *AggregateSignature `json:"activityProof"       gencodec:"optional"`
}

type AggregateSignature struct {
	// leave these pointers nil if they were nil when encoded
	// this is because otherwise rlp creates a signature with new(blst.BlsSignature)
	// which causes all sorts of problem because the private inner signature s.s remains nil
	Signature *blst.BlsSignature `rlp:"nil"`
	Signers   *Signers           `rlp:"nil"`
}

func NewAggregateSignature(signature *blst.BlsSignature, signers *Signers) *AggregateSignature {
	return &AggregateSignature{Signature: signature, Signers: signers}
}

func (a *AggregateSignature) Copy() *AggregateSignature {
	return &AggregateSignature{Signature: a.Signature.Copy(), Signers: a.Signers.Copy()}
}

func (a *AggregateSignature) Malformed() bool {
	return a.Signature == nil || a.Signers == nil
}

// validates the aggregate signature. It does not modify any internal data structure nor does any caching
// returns map of signers and total power of the signers
func (a *AggregateSignature) Validate(message common.Hash, committee *Committee, checkQuorum bool) (map[common.Address]struct{}, *big.Int, error) {
	// validate signers information first
	distinctSigners, err := a.Signers.validate(committee.Len())
	if err != nil {
		return nil, nil, fmt.Errorf("invalid signers information: %w", err)
	}

	// verify signature
	flattenedIndexes := a.Signers.flatten(committee.Len())
	keys := make([]blst.PublicKey, len(flattenedIndexes))
	for i, index := range flattenedIndexes {
		keys[i] = committee.Members[index].ConsensusKey
	}
	aggregatedKey, err := blst.AggregatePublicKeys(keys)
	if err != nil {
		return nil, nil, errors.Join(ErrNonAggregatablePublicKeys, err)
	}
	valid := a.Signature.Verify(aggregatedKey, message[:])
	if !valid {
		return nil, nil, errInvalidSignature
	}

	// Total assembled voting power for the activity proof
	power := new(big.Int)
	signers := make(map[common.Address]struct{}, distinctSigners)
	for _, index := range a.Signers.flattenUniq(committee.Len()) {
		power.Add(power, committee.Members[index].VotingPower)
		signers[committee.Members[index].Address] = struct{}{}
	}

	if checkQuorum && power.Cmp(bft.Quorum(committee.TotalVotingPower())) < 0 {
		return nil, nil, errNoQuorum
	}

	return signers, power, nil
}

// originalHeader represents the ethereum blockchain header.
type originalHeader struct {
	ParentHash  common.Hash    `json:"parentHash"       gencodec:"required"`
	UncleHash   common.Hash    `json:"sha3Uncles"       gencodec:"required"`
	Coinbase    common.Address `json:"miner"            gencodec:"required"`
	Root        common.Hash    `json:"stateRoot"        gencodec:"required"`
	TxHash      common.Hash    `json:"transactionsRoot" gencodec:"required"`
	ReceiptHash common.Hash    `json:"receiptsRoot"     gencodec:"required"`
	Bloom       Bloom          `json:"logsBloom"        gencodec:"required"`
	Difficulty  *big.Int       `json:"difficulty"       gencodec:"required"`
	Number      *big.Int       `json:"number"           gencodec:"required"`
	GasLimit    uint64         `json:"gasLimit"         gencodec:"required"`
	GasUsed     uint64         `json:"gasUsed"          gencodec:"required"`
	Time        uint64         `json:"timestamp"        gencodec:"required"`
	Extra       []byte         `json:"extraData"        gencodec:"required"`
	MixDigest   common.Hash    `json:"mixHash"`
	Nonce       BlockNonce     `json:"nonce"`
	BaseFee     *big.Int       `json:"baseFeePerGas" gencodec:"required"`

	/*
		TODO (MariusVanDerWijden) Add this field once needed
		// Random was added during the merge and contains the BeaconState randomness
		Random common.Hash `json:"random" rlp:"optional"`
	*/
}

// if anything gets added here, check that the `maximumHeaderExtraLength` still makes sense
type headerExtra struct {
	ProposerSeal       []byte `json:"proposerSeal"        gencodec:"required"`
	Round              uint64 `json:"round"               gencodec:"required"`
	ActivityProofRound uint64 `json:"activityProofRound"  gencodec:"required"`

	QuorumCertificate *AggregateSignature `rlp:"nil" json:"quorumCertificate" gencodec:"required"`
	Epoch             *Epoch              `rlp:"nil" json:"epoch"             gencodec:"required"`
	ActivityProof     *AggregateSignature `rlp:"nil" json:"activityProof"     gencodec:"required"`
}

// headerMarshaling is used by gencodec (which can be invoked by running go
// generate in this package) and defines marshalling types for fields that
// would not marshal correctly to hex of their own accord. When modifying the
// structure of Header, this will likely need to be updated before running go
// generate to regenerate the json marshalling code.
type headerMarshaling struct {
	Difficulty *hexutil.Big
	Number     *hexutil.Big
	GasLimit   hexutil.Uint64
	GasUsed    hexutil.Uint64
	Time       hexutil.Uint64
	Extra      hexutil.Bytes
	BaseFee    *hexutil.Big
	Hash       common.Hash `json:"hash"` // adds call to Hash() in MarshalJSON
	/*
		PoS header fields type overrides
	*/
	ProposerSeal       hexutil.Bytes
	Round              hexutil.Uint64
	ActivityProofRound hexutil.Uint64
}

// Hash returns the block hash of the header, which is simply the keccak256 hash of its
// RLP encoding.
func (h *Header) Hash() common.Hash {
	// If the mix digest is equivalent to the predefined BFT digest, use BFT
	// specific hash calculation. This is always the case with tendermint consensus protocol.
	if h.MixDigest == BFTDigest {
		// Seal is reserved in extra-data. To prove block is signed by the proposer.
		if posHeader := BFTFilteredHeader(h, true); posHeader != nil {
			return rlpHash(posHeader)
		}
	}

	// If not using the BFT mixdigest then return the original ethereum block header hash, this
	// let Autonity to remain compatible with original go-ethereum tests.
	return rlpHash(h.original())
}

func (h *Header) IsGenesis() bool {
	return h.Number.Uint64() == 0
}

func (h *Header) IsEpochHeader() bool {
	return h.Epoch != nil
}

// returns true if the block contains a quorum certificate or a finalization round
func (h *Header) HasCommitInformation() bool {
	return h.QuorumCertificate != nil || h.Round != 0
}

var headerSize = common.StorageSize(reflect.TypeOf(Header{}).Size())

// Size returns the approximate memory used by all internal contents. It is used
// to approximate and limit the memory consumption of various caches.
func (h *Header) Size() common.StorageSize {
	return headerSize + common.StorageSize(len(h.Extra)+(h.Difficulty.BitLen()+h.Number.BitLen())/8)
}

// sanityCheck checks a few basic things -- these checks are way beyond what
// any 'sane' production values should hold, and can mainly be used to prevent
// that the unbounded fields are stuffed with junk data to add processing
// overhead
func (h *Header) sanityCheck() error {
	if !h.Number.IsUint64() {
		return fmt.Errorf("too large block number: bitlen %d", h.Number.BitLen())
	}
	if diffLen := h.Difficulty.BitLen(); diffLen > maximumDifficultyBitlen {
		return fmt.Errorf("too large block difficulty: bitlen %d", diffLen)
	}
	if eLen := len(h.Extra); eLen > maximumHeaderExtraLength {
		return fmt.Errorf("too large block extradata: size %d", eLen)
	}
	if bfLen := h.BaseFee.BitLen(); bfLen > maximumBaseFeeBitlen {
		return fmt.Errorf("too large base fee: bitlen %d", bfLen)
	}

	// check sanity of epoch info if the header is an epoch header.
	// assumes sanity nil checks have already been done
	if h.IsEpochHeader() {
		if !h.Epoch.PreviousEpochBlock.IsUint64() {
			return fmt.Errorf("too large previous epoch block number: bitlen %d", h.Epoch.PreviousEpochBlock.BitLen())
		}

		if !h.Epoch.NextEpochBlock.IsUint64() {
			return fmt.Errorf("too large next epoch block number: bitlen %d", h.Epoch.NextEpochBlock.BitLen())
		}

		if !h.Epoch.Delta.IsUint64() {
			return fmt.Errorf("too large next epoch delta: bitlen %d", h.Epoch.Delta.BitLen())
		}

		if h.Epoch.PreviousEpochBlock.Cmp(h.Number) > 0 {
			return fmt.Errorf("previous epoch block number %d is larger than current epoch block number %d", h.Epoch.PreviousEpochBlock.Uint64(), h.Number.Uint64())
		}

		if h.Epoch.PreviousEpochBlock.Cmp(h.Number) == 0 && !h.IsGenesis() { // genesis is allowed to have previousEpochBlock == epochBlock == common.Big0
			return fmt.Errorf("previous epoch block number %d is equal to current epoch block number %d", h.Epoch.PreviousEpochBlock.Uint64(), h.Number.Uint64())
		}

		if h.Number.Cmp(h.Epoch.NextEpochBlock) >= 0 {
			return fmt.Errorf("current epoch block number %d is larger or equal than next epoch block number %d", h.Number.Uint64(), h.Epoch.NextEpochBlock.Uint64())
		}
	}
	return nil
}

// DecodeRLP decodes the Ethereum
func (h *Header) DecodeRLP(s *rlp.Stream) error {
	origin := &originalHeader{}
	if err := s.Decode(origin); err != nil {
		return err
	}

	// *big.Int fields sanity check
	if origin.Number == nil {
		return fmt.Errorf("header number is nil")
	}
	if origin.Difficulty == nil {
		return fmt.Errorf("header difficulty is nil")
	}
	if origin.BaseFee == nil {
		return fmt.Errorf("header base fee is nil")
	}

	if origin.MixDigest == BFTDigest {
		hExtra := &headerExtra{}
		if err := rlp.DecodeBytes(origin.Extra, hExtra); err != nil {
			return err
		}

		// only genesis is allowed to have an empty proposer seal
		if origin.Number.Uint64() != 0 && len(hExtra.ProposerSeal) != crypto.SignatureLength {
			return fmt.Errorf("invalid proposer seal length: %d", len(hExtra.ProposerSeal))
		}

		// sanity checks
		if hExtra.QuorumCertificate != nil {
			if hExtra.QuorumCertificate.Malformed() {
				return fmt.Errorf("malformed header quorum certificate")
			}
			if err := hExtra.QuorumCertificate.Signers.SanityCheck(); err != nil {
				return fmt.Errorf("invalid quorum certificate signers info: %w", err)
			}
		}

		if hExtra.ActivityProof != nil {
			if hExtra.ActivityProof.Malformed() {
				return fmt.Errorf("malformed header activity proof")
			}
			if err := hExtra.ActivityProof.Signers.SanityCheck(); err != nil {
				return fmt.Errorf("invalid activity proof signers info: %w", err)
			}
		}

		if hExtra.ActivityProof == nil && hExtra.ActivityProofRound != 0 {
			return fmt.Errorf("activity proof round should be 0 if proof is empty")
		}

		if hExtra.Epoch != nil {
			if hExtra.Epoch.Committee == nil {
				return fmt.Errorf("committee should not be nil")
			}

			if len(hExtra.Epoch.Committee.Members) == 0 {
				return fmt.Errorf("no members in committee set")
			}

			if err := hExtra.Epoch.Committee.Enrich(); err != nil {
				return fmt.Errorf("error while deserializing consensus keys: %w", err)
			}

			// sanity check the voting power of the members
			for _, member := range hExtra.Epoch.Committee.Members {
				if member.VotingPower == nil {
					return fmt.Errorf("voting power should not be nil")
				}
				if member.VotingPower.BitLen() > maximumVotingPowerBitlen {
					return fmt.Errorf("voting power too large for committee member, bitlen: %d", member.VotingPower.BitLen())
				}
				if member.VotingPower.Cmp(common.Big0) <= 0 {
					return fmt.Errorf("invalid voting power for committee member (negative or equal to 0): %s", member.VotingPower.String())
				}
			}

			if hExtra.Epoch.PreviousEpochBlock == nil || hExtra.Epoch.NextEpochBlock == nil {
				return fmt.Errorf("invalid epoch boundary")
			}

			if hExtra.Epoch.Delta == nil {
				return fmt.Errorf("invalid epoch delta")
			}

			if hExtra.Epoch.Delta.Cmp(common.Big0) == 0 {
				return fmt.Errorf("epoch delta is zero")
			}
		}

		h.QuorumCertificate = hExtra.QuorumCertificate
		h.ActivityProof = hExtra.ActivityProof
		h.ActivityProofRound = hExtra.ActivityProofRound
		h.ProposerSeal = hExtra.ProposerSeal
		h.Round = hExtra.Round
		h.Epoch = hExtra.Epoch
	} else {
		h.Extra = origin.Extra
	}

	h.ParentHash = origin.ParentHash
	h.UncleHash = origin.UncleHash
	h.Coinbase = origin.Coinbase
	h.Root = origin.Root
	h.TxHash = origin.TxHash
	h.ReceiptHash = origin.ReceiptHash
	h.Bloom = origin.Bloom
	h.Difficulty = origin.Difficulty
	h.Number = origin.Number
	h.GasLimit = origin.GasLimit
	h.GasUsed = origin.GasUsed
	h.Time = origin.Time
	h.MixDigest = origin.MixDigest
	h.Nonce = origin.Nonce
	h.BaseFee = origin.BaseFee

	if err := h.sanityCheck(); err != nil {
		return fmt.Errorf("failed sanity check: %w", err)
	}

	return nil
}

// EncodeRLP serializes b into the Ethereum RLP block format.
//
// To maintain RLP compatibility with eth tooling we have to encode our
// additional header fields into the extra data field. RLP decoding expects the
// encoded data to have an exact number of fields of a certain type in a
// particular order, if there is a mismatch decoding fails. So to maintain
// compatibility with ethereum we encode all our additional header fields into
// the extra data field leaving us with just the original ethereum header
// fields. When we decode we repopulate our additional header fields from the
// extra data.
func (h *Header) EncodeRLP(w io.Writer) error {
	hExtra := headerExtra{
		ProposerSeal:       h.ProposerSeal,
		Round:              h.Round,
		QuorumCertificate:  h.QuorumCertificate,
		Epoch:              h.Epoch,
		ActivityProof:      h.ActivityProof,
		ActivityProofRound: h.ActivityProofRound,
	}

	original := h.original()
	if h.MixDigest == BFTDigest {
		extra, err := rlp.EncodeToBytes(hExtra)
		if err != nil {
			return err
		}
		original.Extra = extra
	} else {
		original.Extra = h.Extra
	}

	return rlp.Encode(w, *original)
}

func (h *Header) original() *originalHeader {
	return &originalHeader{
		ParentHash:  h.ParentHash,
		UncleHash:   h.UncleHash,
		Coinbase:    h.Coinbase,
		Root:        h.Root,
		TxHash:      h.TxHash,
		ReceiptHash: h.ReceiptHash,
		Bloom:       h.Bloom,
		Difficulty:  h.Difficulty,
		Number:      h.Number,
		GasLimit:    h.GasLimit,
		GasUsed:     h.GasUsed,
		BaseFee:     h.BaseFee,
		Time:        h.Time,
		Extra:       h.Extra,
		MixDigest:   h.MixDigest,
		Nonce:       h.Nonce,
	}
}

// EmptyBody returns true if there is no additional 'body' to complete the header
// that is: no transactions and no uncles.
func (h *Header) EmptyBody() bool {
	return h.TxHash == EmptyRootHash && h.UncleHash == EmptyUncleHash
}

// EmptyReceipts returns true if there are no receipts for this header/block.
func (h *Header) EmptyReceipts() bool {
	return h.ReceiptHash == EmptyRootHash
}

// Body is a simple (mutable, non-safe) data container for storing and moving
// a block's data contents (transactions and uncles) together.
type Body struct {
	Transactions []*Transaction
	Uncles       []*Header
}

// Block represents an entire block in the Ethereum blockchain.
type Block struct {
	header       *Header
	uncles       []*Header
	transactions Transactions

	// caches
	hash atomic.Value
	size atomic.Value

	// Td is used by package core to store the total difficulty
	// of the chain up to and including the block.
	td *big.Int

	// These fields are used by package eth to track
	// inter-peer block relay.
	ReceivedAt   time.Time
	ReceivedFrom interface{}
}

// "external" block encoding. used for eth protocol, etc.
type extblock struct {
	Header *Header
	Txs    []*Transaction
	Uncles []*Header
}

// NewBlock creates a new block. The input data is copied,
// changes to header and to the field values will not affect the
// block.
//
// The values of TxHash, UncleHash, ReceiptHash and Bloom in header
// are ignored and set to values derived from the given txs, uncles
// and receipts.
func NewBlock(header *Header, txs []*Transaction, uncles []*Header, receipts []*Receipt, hasher TrieHasher) *Block {
	b := &Block{header: CopyHeader(header), td: new(big.Int)}

	// TODO: panic if len(txs) != len(receipts)
	if len(txs) == 0 {
		b.header.TxHash = EmptyRootHash
	} else {
		b.header.TxHash = DeriveSha(Transactions(txs), hasher)
		b.transactions = make(Transactions, len(txs))
		copy(b.transactions, txs)
	}

	if len(receipts) == 0 {
		b.header.ReceiptHash = EmptyRootHash
	} else {
		b.header.ReceiptHash = DeriveSha(Receipts(receipts), hasher)
		b.header.Bloom = CreateBloom(receipts)
	}

	if len(uncles) == 0 {
		b.header.UncleHash = EmptyUncleHash
	} else {
		b.header.UncleHash = CalcUncleHash(uncles)
		b.uncles = make([]*Header, len(uncles))
		for i := range uncles {
			b.uncles[i] = CopyHeader(uncles[i])
		}
	}

	return b
}

// NewBlockWithHeader creates a block with the given header data. The
// header data is copied, changes to header and to the field values
// will not affect the block.
func NewBlockWithHeader(header *Header) *Block {
	return &Block{header: CopyHeader(header)}
}

// CopyHeader creates a deep copy of a block header to prevent side effects from
// modifying a header variable.
func CopyHeader(h *Header) *Header {
	difficulty := big.NewInt(0)
	if h.Difficulty != nil {
		difficulty.Set(h.Difficulty)
	}

	number := big.NewInt(0)
	if h.Number != nil {
		number.Set(h.Number)
	}

	var baseFee *big.Int
	if h.BaseFee != nil {
		baseFee = new(big.Int).Set(h.BaseFee)
	}

	extra := make([]byte, 0)
	if len(h.Extra) > 0 {
		extra = make([]byte, len(h.Extra))
		copy(extra, h.Extra)
	}

	/* PoS fields deep copy section*/
	proposerSeal := make([]byte, 0)
	if len(h.ProposerSeal) > 0 {
		proposerSeal = make([]byte, len(h.ProposerSeal))
		copy(proposerSeal, h.ProposerSeal)
	}

	var quorumCertificate *AggregateSignature
	if h.QuorumCertificate != nil {
		quorumCertificate = h.QuorumCertificate.Copy()
	}

	var activityProof *AggregateSignature
	if h.ActivityProof != nil {
		activityProof = h.ActivityProof.Copy()
	}

	var epoch *Epoch
	if h.Epoch != nil {
		epoch = h.Epoch.Copy()
	}

	cpy := &Header{
		ParentHash:         h.ParentHash,
		UncleHash:          h.UncleHash,
		Coinbase:           h.Coinbase,
		Root:               h.Root,
		TxHash:             h.TxHash,
		ReceiptHash:        h.ReceiptHash,
		Bloom:              h.Bloom,
		Difficulty:         difficulty,
		Number:             number,
		GasLimit:           h.GasLimit,
		GasUsed:            h.GasUsed,
		Time:               h.Time,
		Extra:              extra,
		MixDigest:          h.MixDigest,
		Nonce:              h.Nonce,
		ProposerSeal:       proposerSeal,
		BaseFee:            baseFee,
		Round:              h.Round,
		QuorumCertificate:  quorumCertificate,
		Epoch:              epoch,
		ActivityProof:      activityProof,
		ActivityProofRound: h.ActivityProofRound,
	}

	return cpy
}

// DecodeRLP decodes the Ethereum
func (b *Block) DecodeRLP(s *rlp.Stream) error {
	var eb extblock
	_, size, _ := s.Kind()
	if err := s.Decode(&eb); err != nil {
		return err
	}
	b.header, b.uncles, b.transactions = eb.Header, eb.Uncles, eb.Txs
	b.size.Store(common.StorageSize(rlp.ListSize(size)))
	return nil
}

// EncodeRLP serializes b into the Ethereum RLP block format.
func (b *Block) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, extblock{
		Header: b.header,
		Txs:    b.transactions,
		Uncles: b.uncles,
	})
}

// TODO: copies

func (b *Block) Uncles() []*Header                 { return b.uncles }
func (b *Block) Transactions() Transactions        { return b.transactions }
func (b *Block) SetTransactions(tr []*Transaction) { b.transactions = tr }
func (b *Block) SetHash(h atomic.Value)            { b.hash = h }

func (b *Block) Transaction(hash common.Hash) *Transaction {
	for _, transaction := range b.transactions {
		if transaction.Hash() == hash {
			return transaction
		}
	}
	return nil
}

func (b *Block) Number() *big.Int     { return new(big.Int).Set(b.header.Number) }
func (b *Block) GasLimit() uint64     { return b.header.GasLimit }
func (b *Block) GasUsed() uint64      { return b.header.GasUsed }
func (b *Block) Difficulty() *big.Int { return new(big.Int).Set(b.header.Difficulty) }
func (b *Block) Time() uint64         { return b.header.Time }

func (b *Block) NumberU64() uint64          { return b.header.Number.Uint64() }
func (b *Block) MixDigest() common.Hash     { return b.header.MixDigest }
func (b *Block) Nonce() uint64              { return binary.BigEndian.Uint64(b.header.Nonce[:]) }
func (b *Block) Bloom() Bloom               { return b.header.Bloom }
func (b *Block) Coinbase() common.Address   { return b.header.Coinbase }
func (b *Block) Root() common.Hash          { return b.header.Root }
func (b *Block) ParentHash() common.Hash    { return b.header.ParentHash }
func (b *Block) TxHash() common.Hash        { return b.header.TxHash }
func (b *Block) ReceiptHash() common.Hash   { return b.header.ReceiptHash }
func (b *Block) UncleHash() common.Hash     { return b.header.UncleHash }
func (b *Block) Extra() []byte              { return common.CopyBytes(b.header.Extra) }
func (b *Block) HasCommitInformation() bool { return b.header.HasCommitInformation() }

func (b *Block) BaseFee() *big.Int {
	if b.header.BaseFee == nil {
		return nil
	}
	return new(big.Int).Set(b.header.BaseFee)
}

func (b *Block) Header() *Header               { return CopyHeader(b.header) }
func (b *Block) SetHeaderNumber(hNum *big.Int) { b.header.Number = hNum }

// Body returns the non-header content of the block.
func (b *Block) Body() *Body { return &Body{b.transactions, b.uncles} }

// Size returns the true RLP encoded storage size of the block, either by encoding
// and returning it, or returning a previsouly cached value.
func (b *Block) Size() common.StorageSize {
	if size := b.size.Load(); size != nil {
		return size.(common.StorageSize)
	}
	c := writeCounter(0)
	rlp.Encode(&c, b)
	b.size.Store(common.StorageSize(c))
	return common.StorageSize(c)
}

func (b *Block) IsEpochHead() bool { return b.header.IsEpochHeader() }

type writeCounter common.StorageSize

func (c *writeCounter) Write(b []byte) (int, error) {
	*c += writeCounter(len(b))
	return len(b), nil
}

func CalcUncleHash(uncles []*Header) common.Hash {
	if len(uncles) == 0 {
		return EmptyUncleHash
	}
	// len(uncles) > 0 can only happen during tests.
	// We revert to the original structure to keep compatibility with hardcoded hash values.
	originalUncles := make([]*originalHeader, len(uncles))
	for i := range uncles {
		originalUncles[i] = uncles[i].original()
	}
	return rlpHash(originalUncles)
}

// WithSeal returns a new block with the data from b but the header replaced with
// the sealed one.
func (b *Block) WithSeal(header *Header) *Block {
	cpy := *header
	return &Block{
		header:       &cpy,
		transactions: b.transactions,
		uncles:       b.uncles,
	}
}

// WithBody returns a new block with the given transaction and uncle contents.
func (b *Block) WithBody(transactions []*Transaction, uncles []*Header) *Block {
	block := &Block{
		header:       CopyHeader(b.header),
		transactions: make([]*Transaction, len(transactions)),
		uncles:       make([]*Header, len(uncles)),
	}
	copy(block.transactions, transactions)
	for i := range uncles {
		block.uncles[i] = CopyHeader(uncles[i])
	}
	return block
}

// Hash returns the keccak256 hash of b's header.
// The hash is computed on the first call and cached thereafter.
func (b *Block) Hash() common.Hash {
	if hash := b.hash.Load(); hash != nil {
		return hash.(common.Hash)
	}
	v := b.header.Hash()
	b.hash.Store(v)
	return v
}

type Blocks []*Block

// HeaderParentHashFromRLP returns the parentHash of an RLP-encoded
// header. If 'header' is invalid, the zero hash is returned.
func HeaderParentHashFromRLP(header []byte) common.Hash {
	// parentHash is the first list element.
	listContent, _, err := rlp.SplitList(header)
	if err != nil {
		return common.Hash{}
	}
	parentHash, _, err := rlp.SplitString(listContent)
	if err != nil {
		return common.Hash{}
	}
	if len(parentHash) != 32 {
		return common.Hash{}
	}
	return common.BytesToHash(parentHash)
}
