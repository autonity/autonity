package accountability

import (
	"errors"
	"io"
	"math/bits"

	blstbind "github.com/supranational/blst/bindings/go"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/core/constants"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/crypto/blst"
	"github.com/autonity/autonity/log"
	"github.com/autonity/autonity/rlp"
)

var (
	ErrSignatureInvalid     = errors.New("HighlyAggregatedPrecommit has invalid signature")
	ErrInvalidSignerIndex   = errors.New("HighlyAggregatedPrecommit has invalid signer index")
	ErrInvalidSignerCoeff   = errors.New("HighlyAggregatedPrecommit has invalid signer coefficient")
	ErrNoSigners            = errors.New("no signers found")
	ErrInvalidRound         = errors.New("invalid round")
	ErrDuplicatedPrecommits = errors.New("duplicated precommits")
)

// Signers is a set that contains signers of the same message with the using of fastAggregate().
type Signers struct {
	Round        int64
	Value        common.Hash
	SignersIndex []int
	SignersCoeff []uint16

	// computed fields
	aggregatedPublicKey blst.PublicKey   `rlp:"-"`
	hasSigners          map[int]struct{} `rlp:"-"`
	preValidated        bool             `rlp:"-"`
}

type extSigners struct {
	Round        uint64
	Value        common.Hash
	SignersIndex []uint
	SignersCoeff []uint16
}

func (r *Signers) EncodeRLP(w io.Writer) error {
	signersIndex := make([]uint, len(r.SignersIndex))
	signersCoeff := make([]uint16, len(r.SignersCoeff))
	for i, s := range r.SignersIndex {
		signersIndex[i] = uint(s)
	}
	copy(signersCoeff, r.SignersCoeff)

	ext := extSigners{
		Round:        uint64(r.Round),
		Value:        r.Value,
		SignersIndex: signersIndex,
		SignersCoeff: signersCoeff,
	}

	return rlp.Encode(w, &ext)
}

func (r *Signers) DecodeRLP(stream *rlp.Stream) error {
	ext := extSigners{}
	if err := stream.Decode(&ext); err != nil {
		return err
	}

	if len(ext.SignersIndex) == 0 {
		return ErrNoSigners
	}
	if len(ext.SignersCoeff) != len(ext.SignersIndex) {
		return ErrInvalidSignerCoeff
	}

	if int64(ext.Round) > constants.MaxRound {
		return ErrInvalidRound
	}

	r.Round = int64(ext.Round)
	r.Value = ext.Value
	signersIndex := make([]int, len(ext.SignersIndex))
	signersCoeff := make([]uint16, len(ext.SignersCoeff))
	for i, s := range ext.SignersIndex {
		signersIndex[i] = int(s)
	}
	copy(signersCoeff, ext.SignersCoeff)
	r.SignersIndex = signersIndex
	r.SignersCoeff = signersCoeff
	return nil
}

// PreValidate computes the aggregated public key and set the preValidated flag.
func (r *Signers) PreValidate(committee *types.Committee) error {
	if len(r.SignersIndex) == 0 || len(r.SignersCoeff) != len(r.SignersIndex) {
		return ErrSignatureInvalid
	}

	committeeSize := committee.Len()
	// early return, as signer indexes are distinct
	if len(r.SignersIndex) > committeeSize {
		return ErrSignatureInvalid
	}

	publicKeys := make([]blst.PublicKey, len(r.SignersIndex))
	r.hasSigners = make(map[int]struct{})
	var maxCoeff uint16

	for i, idx := range r.SignersIndex {
		if idx >= committeeSize || idx < 0 {
			return ErrInvalidSignerIndex
		}
		if r.SignersCoeff[i] == 0 {
			return ErrInvalidSignerCoeff
		}
		maxCoeff = max(maxCoeff, r.SignersCoeff[i])
		publicKeys[i] = committee.Members[idx].ConsensusKey

		// distinct signers
		if _, ok := r.hasSigners[idx]; ok {
			return ErrSignatureInvalid
		}
		r.hasSigners[idx] = struct{}{}
	}

	if len(r.SignersCoeff) == 1 && maxCoeff != 1 {
		return ErrSignatureInvalid
	}

	if len(r.SignersCoeff) > 1 && len(r.SignersCoeff) < 18 {
		if maxCoeff > (1 << (len(r.SignersCoeff) - 2)) {
			return ErrInvalidSignerCoeff
		}
	}

	aggKey := blst.AggregatePublicKeysMultScalars(
		publicKeys,
		r.toBlstScalars(),
		bits.Len16(maxCoeff),
	)

	r.aggregatedPublicKey = aggKey
	r.preValidated = true
	return nil
}

func (r *Signers) toBlstScalars() []*blstbind.Scalar {
	return blst.ToBlstScalars(r.SignersCoeff)
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
func (h *HighlyAggregatedPrecommit) PreValidate(committee *types.Committee, eventHeight uint64) error {

	if h.Height != eventHeight {
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

		if err := m.PreValidate(committee); err != nil {
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

	publicKeys := make([]blst.PublicKey, 0, len(h.MsgSigners))
	msgs := make([][32]byte, 0, len(h.MsgSigners))
	for _, m := range h.MsgSigners {
		pubKey := m.AggregatedPublicKey()
		// if the public key is infinite, skip it. Signature should be valid regardless if legit.
		if !pubKey.Validate() {
			log.Warn("detected infinite public key", "signers", m)
			continue
		}
		publicKeys = append(publicKeys, pubKey)
		msgs = append(msgs, message.VoteSignatureInput(h.Height, uint64(m.Round), message.PrecommitCode, m.Value))
	}

	var aggregatedPublicKey blst.PublicKey
	if len(publicKeys) > 0 {
		var err error
		aggregatedPublicKey, err = blst.AggregatePublicKeys(publicKeys)
		if err != nil {
			// should not happen since all public keys are validated at validator registration
			panic("cannot aggregate public keys from committee: " + err.Error())
		}
	}

	// if the aggregated public key is not 0, validate signature
	// otherwise only check that also the signature is 0
	if len(publicKeys) > 0 && aggregatedPublicKey.Validate() {
		if !h.signature.AggregateVerify(publicKeys, msgs) {
			return ErrSignatureInvalid
		}
	} else {
		if !h.signature.IsZero() {
			return ErrSignatureInvalid
		}
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

	// skip duplicated msg
	for i := 1; i < len(precommits); i++ {
		roundMap, ok := presentedMsgs[precommits[i].R()]
		if !ok {
			roundMap = make(map[common.Hash]struct{})
			presentedMsgs[precommits[i].R()] = roundMap
		}

		if _, ok := roundMap[precommits[i].Value()]; !ok {
			roundMap[precommits[i].Value()] = struct{}{}
			precommitsToBeAggregated = append(precommitsToBeAggregated, precommits[i])
		}
	}

	result := HighlyAggregatedPrecommit{}
	signatures := make([]blst.Signature, len(precommitsToBeAggregated))
	for i, m := range precommitsToBeAggregated {
		roundValueSigners := &Signers{
			Round:        m.R(),
			Value:        m.Value(),
			SignersIndex: m.Signers().FlattenUniq(),
			SignersCoeff: m.Signers().CopyCoefficients(),
		}
		result.MsgSigners = append(result.MsgSigners, roundValueSigners)
		signatures[i] = m.Signature()
	}
	result.Height = height
	result.Signature = blst.AggregateSignatures(signatures).Marshal()
	return result
}

// AggregateSamePrevotes assumes the votes are for the same msg, it does a BLS fast aggregate for the input signatures.
func AggregateSamePrevotes(prevotes []*message.Prevote) *message.EvidenceVote {
	votes := make([]message.Vote, len(prevotes))
	for i, prevote := range prevotes {
		votes[i] = prevote
	}
	return message.AggregatePrevotesToEvidence(votes)
}
