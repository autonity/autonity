package accountability

import (
	"errors"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/core/constants"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/crypto/blst"
	"github.com/autonity/autonity/rlp"
	"io"
)

var (
	ErrSignatureInvalid     = errors.New("HighlyAggregatedPrecommit has invalid signature")
	ErrInvalidSignerIndex   = errors.New("HighlyAggregatedPrecommit has invalid signer index")
	ErrNoSigners            = errors.New("No signers found")
	ErrInvalidRound         = errors.New("Invalid round")
	ErrDuplicatedPrecommits = errors.New("Duplicated precommits")
)

// Signers is set that contains signers of the same message with the using of fastAggregate().
type Signers struct {
	Round   int64
	Value   common.Hash
	Signers []int // it could contain duplicated index.

	// computed fields
	aggregatedPublicKey blst.PublicKey   `rlp:"-"`
	hasSigners          map[int]struct{} `rlp:"-"`
	preValidated        bool             `rlp:"-"`
}

type extSigners struct {
	Round   uint64
	Value   common.Hash
	Signers []uint
}

func (r *Signers) EncodeRLP(w io.Writer) error {
	signers := make([]uint, len(r.Signers))
	for i, s := range r.Signers {
		signers[i] = uint(s)
	}

	ext := extSigners{
		Round:   uint64(r.Round),
		Value:   r.Value,
		Signers: signers,
	}

	return rlp.Encode(w, &ext)
}

func (r *Signers) DecodeRLP(stream *rlp.Stream) error {
	ext := extSigners{}
	if err := stream.Decode(&ext); err != nil {
		return err
	}

	if len(ext.Signers) == 0 {
		return ErrNoSigners
	}

	if int64(ext.Round) > constants.MaxRound {
		return ErrInvalidRound
	}

	r.Round = int64(ext.Round)
	r.Value = ext.Value
	signers := make([]int, len(ext.Signers))
	for i, s := range ext.Signers {
		signers[i] = int(s)
	}
	r.Signers = signers
	return nil
}

// PreValidate computes the aggregated public key and set the preValidated flag.
func (r *Signers) PreValidate(parentHeader *types.Header) error {
	committeeSize := parentHeader.Committee.Len()
	publicKeys := make([][]byte, len(r.Signers))
	r.hasSigners = make(map[int]struct{})
	for i, idx := range r.Signers {
		if idx >= committeeSize || idx < 0 {
			return ErrInvalidSignerIndex
		}
		publicKeys[i] = parentHeader.Committee[idx].ConsensusKeyBytes
		r.hasSigners[idx] = struct{}{}
	}

	aggKey, err := blst.AggregatePublicKeys(publicKeys)
	if err != nil {
		// should not happen, as the public key is query from the committee by their index.
		panic(err)
	}

	r.aggregatedPublicKey = aggKey
	r.preValidated = true
	return nil
}

func (r *Signers) Contains(index int) bool {
	if !r.preValidated {
		panic("Signers was not pre-validated yet")
	}
	_, ok := r.hasSigners[index]
	return ok
}

func (r *Signers) AggregatedPublicKey() blst.PublicKey {
	if !r.preValidated {
		panic("Signers was not pre-validated yet")
	}
	return r.aggregatedPublicKey
}

// HighlyAggregatedPrecommit is used only in the context of accountability to aggregate different precommit messages
// into a single one as a reasonable proof that could be carried by an accountability TXN.
type HighlyAggregatedPrecommit struct {
	// in the proof's context of accountability event, height is always common, thus we put it in the common part
	// to save the size of the msg, moreover that, as the consensus step/code is always be precommit, thus we omit it
	// as well to save size of TXN.
	Height uint64

	// Distinguish info of the precommit are grouped into a set, MsgSigners, each of them contains a fast aggregated
	// precommit over the same (h, r, code, value). Thus, there are multiple distinct sets to be aggregated and
	// to be aggregateVerified for the single highly aggregated precommit.
	MsgSigners []*Signers

	// a single highly aggregated signature.
	Signature []byte

	// computed fields from the validation phase.
	signature    blst.Signature `rlp:"-"`
	preValidated bool           `rlp:"-"`
	validated    bool           `rlp:"-"`
}

func (h *HighlyAggregatedPrecommit) Len() int {
	return len(h.MsgSigners)
}

// PreValidate checks if the index of each sub set are reasonable, and aggregate public keys for each sub set.
func (h *HighlyAggregatedPrecommit) PreValidate(parentHeader *types.Header) error {

	if h.Height-1 != parentHeader.Number.Uint64() {
		return errBadHeight
	}

	// check there are no duplicated precommits
	presentedMsgs := make(map[int64]map[common.Hash]struct{})
	for _, m := range h.MsgSigners {
		roundMap, ok := presentedMsgs[m.Round]
		if !ok {
			roundMap = make(map[common.Hash]struct{})
			presentedMsgs[m.Round] = roundMap
		}

		if _, ok = roundMap[m.Value]; ok {
			return ErrDuplicatedPrecommits
		}

		roundMap[m.Value] = struct{}{}

		if err := m.PreValidate(parentHeader); err != nil {
			return err
		}
	}

	signature, err := blst.SignatureFromBytes(h.Signature)
	if err != nil {
		return err
	}

	h.signature = signature
	h.preValidated = true
	return nil
}

// Validate validate the aggregated signature
func (h *HighlyAggregatedPrecommit) Validate() error {
	if !h.preValidated {
		panic("HighlyAggregatedPrecommit was not pre-validated yet")
	}

	publicKeys := make([]blst.PublicKey, len(h.MsgSigners))
	msgs := make([][32]byte, len(h.MsgSigners))
	for i, m := range h.MsgSigners {
		publicKeys[i] = m.AggregatedPublicKey()
		msgs[i] = message.VoteSignatureInput(h.Height, uint64(m.Round), message.PrecommitCode, m.Value)
	}

	if !h.signature.AggregateVerify(publicKeys, msgs) {
		return ErrSignatureInvalid
	}

	h.validated = true
	return nil
}

// It assumes that the input has multiple precommits and they are sorted by round ascending with same height
func AggregateDistinctPrecommits(precommits []*message.Precommit) HighlyAggregatedPrecommit {
	var precommitsToBeAggregated []*message.Precommit

	presentedMsgs := make(map[int64]map[common.Hash]struct{})

	precommitsToBeAggregated = append(precommitsToBeAggregated, precommits[0])
	height := precommits[0].H()

	for i := 1; i < len(precommits); i++ {
		// skip duplicated msg.
		if _, ok := presentedMsgs[precommits[i].R()]; !ok {
			presentedMsgs[precommits[i].R()] = make(map[common.Hash]struct{})
		}
		if _, ok := presentedMsgs[precommits[i].R()][precommits[i].Value()]; !ok {
			presentedMsgs[precommits[i].R()][precommits[i].Value()] = struct{}{}
			precommitsToBeAggregated = append(precommitsToBeAggregated, precommits[i])
		}
	}

	result := HighlyAggregatedPrecommit{}
	signatures := make([]blst.Signature, len(precommitsToBeAggregated))
	for i, m := range precommitsToBeAggregated {
		roundValueSigners := &Signers{
			Round:   m.R(),
			Value:   m.Value(),
			Signers: m.Signers().Flatten(),
		}
		result.MsgSigners = append(result.MsgSigners, roundValueSigners)
		signatures[i] = m.Signature()
	}
	result.Height = height
	result.Signature = blst.AggregateSignatures(signatures).Marshal()
	return result
}

// AggregateSamePrevotes assumes the votes are for the same msg, it does a BLS fast aggregate for the input signatures.
func AggregateSamePrevotes(prevotes []*message.Prevote) *message.Prevote {
	votes := make([]message.Vote, len(prevotes))
	for i, prevote := range prevotes {
		votes[i] = prevote
	}
	return message.AggregatePrevotes(votes)
}
