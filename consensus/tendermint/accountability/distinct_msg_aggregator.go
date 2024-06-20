package accountability

import (
	"errors"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/crypto/blst"
	"github.com/autonity/autonity/rlp"
	"io"
)

var (
	ErrSignatureInvalid   = errors.New("HighlyAggregatedPrecommit has invalid signature")
	ErrInvalidSignerIndex = errors.New("HighlyAggregatedPrecommit has invalid signer index")
)

// RoundValueSigners is set that contains signers of the same message with the using of fastAggregate().
type RoundValueSigners struct {
	Round   int64
	Value   common.Hash
	Signers []int // it could contain duplicated index.

	// computed fields
	aggregatedPublicKey blst.PublicKey   `rlp:"-"`
	hasSigners          map[int]struct{} `rlp:"-"`
	preValidated        bool             `rlp:"-"`
}

type extRoundValueSigners struct {
	Round   uint64
	Value   common.Hash
	Signers []uint
}

func (r *RoundValueSigners) EncodeRLP(w io.Writer) error {
	signers := make([]uint, len(r.Signers))
	for i, s := range r.Signers {
		signers[i] = uint(s)
	}

	ext := extRoundValueSigners{
		Round:   uint64(r.Round),
		Value:   r.Value,
		Signers: signers,
	}

	return rlp.Encode(w, &ext)
}

func (r *RoundValueSigners) DecodeRLP(stream *rlp.Stream) error {
	ext := extRoundValueSigners{}
	if err := stream.Decode(&ext); err != nil {
		return err
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
func (r *RoundValueSigners) PreValidate(parentHeader *types.Header) error {
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
		return err
	}

	r.aggregatedPublicKey = aggKey
	r.preValidated = true
	return nil
}

func (r *RoundValueSigners) Contains(index int) bool {
	_, ok := r.hasSigners[index]
	return ok
}

func (r *RoundValueSigners) AggregatedPublicKey() blst.PublicKey {
	if !r.preValidated {
		panic("RoundValueSigners was not pre-validated yet")
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

	// Distinguish info of the precommit are grouped into a set, RoundValueSigners, each of them contains a fast aggregated
	// precommit over the same (h, r, code, value). Thus, there are multiple distinct sets to be aggregated and
	// to be aggregateVerified for the single highly aggregated precommit.
	RoundValueSigners []*RoundValueSigners

	// a single highly aggregated signature.
	Signature []byte

	// computed fields from the validation phase.
	signature    blst.Signature `rlp:"-"`
	preValidated bool           `rlp:"-"`
	validated    bool           `rlp:"-"`
}

func (h *HighlyAggregatedPrecommit) Len() int {
	return len(h.RoundValueSigners)
}

// PreValidate checks if the index of each sub set are reasonable, and aggregate public keys for each sub set.
func (h *HighlyAggregatedPrecommit) PreValidate(parentHeader *types.Header) error {

	if h.Height-1 != parentHeader.Number.Uint64() {
		return errBadHeight
	}

	for _, m := range h.RoundValueSigners {
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

	publicKeys := make([]blst.PublicKey, len(h.RoundValueSigners))
	msgs := make([][32]byte, len(h.RoundValueSigners))
	for i, m := range h.RoundValueSigners {
		publicKeys[i] = m.AggregatedPublicKey()
		msgs[i] = message.VoteSignatureInput(h.Height, uint64(m.Round), message.PrecommitCode, m.Value)
	}

	if !h.signature.AggregateVerify(publicKeys, msgs) {
		return ErrSignatureInvalid
	}

	h.validated = true
	return nil
}

// AggregateDistinctPrecommits assumes that the input precommits are sorted by round ascending with same height, the
// input precommits could be fastAggregated it here are duplicated msg.
func AggregateDistinctPrecommits(precommits []*message.Precommit) HighlyAggregatedPrecommit {

	var fastAggregatedPrecommits []*message.Precommit

	// first we filter out fast-aggregatable votes with same msg, and fast aggregate them into single one.
	var fastAggregatablePrecommits []message.Vote
	fastAggregatablePrecommits = append(fastAggregatablePrecommits, precommits[0])
	height := precommits[0].H()

	for i := 1; i < len(precommits); i++ {

		if fastAggregatablePrecommits[0].R() == precommits[i].R() &&
			fastAggregatablePrecommits[0].Value() == precommits[i].Value() {

			fastAggregatablePrecommits = append(fastAggregatablePrecommits, precommits[i])

		} else {
			// if there are multiple fast aggregatable precommits:  aggregate them into single one and append the
			// aggregated one to the fastAggregatedPrecommits, then reset fastAggregatablePrecommits with current
			// precommit and continue.
			if len(fastAggregatablePrecommits) > 1 {
				aggregatedMsg := message.AggregatePrecommits(fastAggregatablePrecommits)
				fastAggregatedPrecommits = append(fastAggregatedPrecommits, aggregatedMsg)
				fastAggregatablePrecommits = []message.Vote{precommits[i]}
				continue
			}

			// if there is only one precommit in the fastAggregatablePrecommits: append it to the fastAggregatedPrecommits.
			// then reset fastAggregatablePrecommits with current precommit and continue.
			fastAggregatedPrecommits = append(fastAggregatedPrecommits, fastAggregatablePrecommits[0].(*message.Precommit))
			fastAggregatablePrecommits = []message.Vote{precommits[i]}
		}
	}

	// append the last patch of fastAggregatablePrecommits to fastAggregatedPrecommits.
	if len(fastAggregatablePrecommits) > 1 {
		aggregatedMsg := message.AggregatePrecommits(fastAggregatablePrecommits)
		fastAggregatedPrecommits = append(fastAggregatedPrecommits, aggregatedMsg)
	} else {
		fastAggregatedPrecommits = append(fastAggregatedPrecommits, fastAggregatablePrecommits[0].(*message.Precommit))
	}

	result := HighlyAggregatedPrecommit{}
	signatures := make([]blst.Signature, len(fastAggregatedPrecommits))
	for i, m := range fastAggregatedPrecommits {
		roundValueSigners := &RoundValueSigners{
			Round:   m.R(),
			Value:   m.Value(),
			Signers: m.Signers().Flatten(),
		}
		result.RoundValueSigners = append(result.RoundValueSigners, roundValueSigners)
		signatures[i] = m.Signature()
	}
	result.Height = height
	result.Signature = blst.AggregateSignatures(signatures).Marshal()
	return result
}

// FastAggregatePrevotes assumes the votes are for the same msg, it does a BLS fast aggregate for the input signatures.
func FastAggregatePrevotes(prevotes []*message.Prevote) *message.Prevote {
	votes := make([]message.Vote, len(prevotes))
	for i, prevote := range prevotes {
		votes[i] = prevote
	}
	return message.AggregatePrevotes(votes)
}
