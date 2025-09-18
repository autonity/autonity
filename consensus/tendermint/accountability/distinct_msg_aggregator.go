package accountability

import (
	"errors"
	"fmt"
	"io"
	"math/big"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/core/constants"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/crypto/blst"
	"github.com/autonity/autonity/log"
	"github.com/autonity/autonity/rlp"
)

var (
	ErrSignatureInvalid           = errors.New("HighlyAggregatedPrecommit has invalid signature")
	ErrInvalidIndividualSignature = errors.New("HighlyAggregatedPrecommit has invalid individual signature")
	ErrInvalidSignerIndex         = errors.New("HighlyAggregatedPrecommit has invalid signer index")
	ErrInvalidSignerCoeff         = errors.New("HighlyAggregatedPrecommit has invalid signer coefficients")
	ErrInvalidSigners             = errors.New("number of signers is not valid")
	ErrInvalidRound               = errors.New("invalid round")
	ErrDuplicatedPrecommits       = errors.New("duplicated precommits")
)

// Signers is a set that contains signers of the same message with the using of fastAggregate().
type Signers struct {
	Round        int64
	Value        common.Hash
	SignersIndex []int
	SignersCoeff []*big.Int

	// computed fields
	aggregatedPublicKey blst.PublicKey   `rlp:"-"`
	hasSigners          map[int]struct{} `rlp:"-"`
	preValidated        bool             `rlp:"-"`
}

type extSigners struct {
	Round        uint64
	Value        common.Hash
	SignersIndex []*big.Int
	SignersCoeff []*big.Int
}

func (r *Signers) EncodeRLP(w io.Writer) error {
	signersIndex := make([]*big.Int, len(r.SignersIndex))
	for i, s := range r.SignersIndex {
		signersIndex[i] = big.NewInt(int64(s))
	}

	ext := extSigners{
		Round:        uint64(r.Round), //nolint:gosec
		Value:        r.Value,
		SignersIndex: signersIndex,
		SignersCoeff: r.SignersCoeff,
	}

	return rlp.Encode(w, &ext)
}

func (r *Signers) DecodeRLP(stream *rlp.Stream) error {
	ext := extSigners{}
	if err := stream.Decode(&ext); err != nil {
		return err
	}

	if int64(ext.Round) > constants.MaxRound {
		return ErrInvalidRound
	}

	if len(ext.SignersIndex) == 0 || len(ext.SignersIndex) > types.MaxAllowedSigners {
		return ErrInvalidSigners
	}
	if len(ext.SignersCoeff) != len(ext.SignersIndex) {
		return ErrInvalidSignerCoeff
	}
	for i, coefficient := range ext.SignersCoeff {
		// coefficients cannot be nil
		if coefficient == nil {
			return fmt.Errorf("coefficient cannot be nil")
		}
		// every coefficient should be > 0
		if coefficient.Sign() <= 0 {
			return fmt.Errorf("invalid coefficient. Sign: %d", coefficient.Sign())
		}
		// coefficient cannot exceed VoteCap bitsize as they are "normal" network messages (not aggregated further)
		if coefficient.BitLen() > common.VoteCap {
			return fmt.Errorf("coefficient too big. BitLen: %d", coefficient.BitLen())
		}
		index := ext.SignersIndex[i]
		// indexes cannot be nil
		if index == nil {
			return fmt.Errorf("index cannot be nil")
		}
		// every index should be >= 0
		if index.Sign() < 0 {
			return fmt.Errorf("invalid index. Sign: %d", index.Sign())
		}
		// index cannot be >= MaxAllowedSigners
		if !index.IsInt64() || index.Cmp(types.MaxAllowedSignersBig) >= 0 {
			return fmt.Errorf("index too big: %s", index.String())
		}
	}

	r.Round = int64(ext.Round)
	r.Value = ext.Value
	signersIndex := make([]int, len(ext.SignersIndex))
	for i, s := range ext.SignersIndex {
		signersIndex[i] = int(s.Int64())
	}
	r.SignersIndex = signersIndex
	r.SignersCoeff = ext.SignersCoeff
	return nil
}

// PreValidate computes the aggregated public key and set the preValidated flag.
func (r *Signers) PreValidate(committee *types.Committee) error {
	committeeSize := committee.Len()
	// early return, as signer indexes are distinct
	if len(r.SignersIndex) > committeeSize {
		return ErrInvalidSigners
	}

	publicKeys := make([]blst.PublicKey, len(r.SignersIndex))
	r.hasSigners = make(map[int]struct{})
	maxCoefficient := new(big.Int)

	for i, idx := range r.SignersIndex {
		if idx >= committeeSize || idx < 0 {
			return ErrInvalidSignerIndex
		}
		publicKeys[i] = committee.Members[idx].ConsensusKey

		// distinct signers
		if _, ok := r.hasSigners[idx]; ok {
			return ErrInvalidSigners
		}
		r.hasSigners[idx] = struct{}{}

		maxCoefficient = common.Max(maxCoefficient, r.SignersCoeff[i])
	}

	var aggregatedKey blst.PublicKey

	// if it is an individual signatures, it should have coefficient 1 and aggregation can be skipped
	if len(r.SignersCoeff) == 1 {
		if r.SignersCoeff[0].Cmp(common.Big1) != 0 {
			return ErrInvalidIndividualSignature
		}
		aggregatedKey = publicKeys[0]
	} else {
		// len(publicKeys) > 1

		// if the maximum coefficient is 1 --> all coefficients are 1 --> no need for scalar multiplication, just sum the pubkeys
		if maxCoefficient.Cmp(common.Big1) == 0 {
			var err error
			aggregatedKey, err = blst.AggregatePublicKeys(publicKeys)
			if err != nil {
				// should not happen since all public keys are validated at validator registration
				panic("cannot aggregate public keys from committee: " + err.Error())
			}
		} else {
			// otherwise, do the scalar multiplication
			aggregatedKey = blst.AggregatePublicKeysMultScalars(publicKeys, blst.ToScalars(r.SignersCoeff))
		}
	}

	r.aggregatedPublicKey = aggregatedKey
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
	height := precommits[0].H()

	// skip precommits for same r and value but different signers set. This function assumes all precommits include
	// the offender/accused, therefore a single one including it is enough to prove its behaviour.
	for _, m := range precommits {
		roundMap, ok := presentedMsgs[m.R()]
		if !ok {
			roundMap = make(map[common.Hash]struct{})
			presentedMsgs[m.R()] = roundMap
		}

		if _, ok := roundMap[m.Value()]; !ok {
			roundMap[m.Value()] = struct{}{}
			precommitsToBeAggregated = append(precommitsToBeAggregated, m)
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

// aggregateSamePrevotes assumes the votes are for the same msg. It is a convenience wrapper whose
// main function is to convert from []*message.Prevote to []message.Vote
func aggregateSamePrevotes(prevotes []*message.Prevote) *message.Prevote {
	// shortcircuit if we have a single prevote, no need to get into aggregation
	if len(prevotes) == 1 {
		return prevotes[0]
	}
	votes := make([]message.Vote, len(prevotes))
	for i, prevote := range prevotes {
		votes[i] = prevote
	}
	return message.AggregatePrevotesSingle(votes)
}
