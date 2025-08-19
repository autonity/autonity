package types

import (
	"errors"
	"fmt"
	"math/big"
	"sync"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/crypto/blst"
)

var (
	MaxAllowedSigners    = 16384
	MaxAllowedSignersBig = new(big.Int).SetUint64(uint64(MaxAllowedSigners))

	ErrWrongSizeSigners    = errors.New("validator bitmap or coefficient array has incorrect size")
	ErrEmptySigners        = errors.New("signers information is empty")
	ErrWrongCoefficientLen = errors.New("coefficient array has incorrect length")
	ErrInvalidSingleSig    = errors.New("individual signature has coefficient != 1")
	ErrInvalidCoefficient  = errors.New("coefficient exceeds maximum boundary")
	ErrNotValidated        = errors.New("using un-validated signers information")
	ErrDifferentSize       = errors.New("comparing signers information with different committee size")
)

// represents the senders of an aggregated signature

type Signers struct {
	Bitmap       *Bitmap
	Coefficients []*big.Int

	// these fields are not serialized, but instead computed at preValidate steps
	committee     *Committee   `rlp:"-"`
	powerMu       sync.RWMutex `rlp:"-"`
	computedPower *big.Int     `rlp:"-"`

	length          int      `rlp:"-"` // number of distinct signers
	rightmostSigner int      `rlp:"-"` // last (LSB) > 0 index in Bitmap
	maxCoefficient  *big.Int `rlp:"-"`

	// auxiliary data structures flags
	// if validated = true -->
	// 1. Bitmap and Coefficients have been validated
	// 2. committeeSize, length, rightmostSigner and maxCoefficient are assigned
	validated bool `rlp:"-"`
}

func NewSigners(committee *Committee) *Signers {
	return &Signers{
		Bitmap:       NewBitmap(),
		Coefficients: make([]*big.Int, 0),
		committee:    committee,

		length:          0,
		rightmostSigner: committee.Len(), // means no signers
		maxCoefficient:  new(big.Int),
		validated:       true,
	}
}

func (s *Signers) Committee() *Committee {
	if !s.validated {
		panic("Using un-validated signers information")
	}
	return s.committee
}

// Correct full validation cannot be done until we know the committee size of this block,
// but we can already ensure that the fields have sane values
func (s *Signers) SanityCheck() error {
	if s.Bitmap == nil || s.Coefficients == nil {
		return errors.New("bitmap and coefficients cannot be nil")
	}
	if s.Bitmap.Sign() <= 0 {
		return fmt.Errorf("validator bitmap is invalid")
	}
	if s.Bitmap.Len() > MaxAllowedSigners {
		return fmt.Errorf("invalid Bitmap length: %d", s.Bitmap.Len())
	}
	if len(s.Coefficients) == 0 || len(s.Coefficients) > MaxAllowedSigners {
		return fmt.Errorf("invalid Coefficient length: %d", len(s.Coefficients))
	}
	for _, coefficient := range s.Coefficients {
		// coefficients cannot be nil
		if coefficient == nil {
			return fmt.Errorf("coefficient cannot be nil")
		}
		// every coefficient should be > 0
		if coefficient.Sign() <= 0 {
			return fmt.Errorf("invalid coefficient. Sign: %d", coefficient.Sign())
		}
		// coefficient cannot exceed QuorumVote bitsize in any situation
		if coefficient.BitLen() > common.QuorumCap {
			return fmt.Errorf("coefficient too big. BitLen: %d", coefficient.BitLen())
		}
	}
	return nil
}

// validates the sender info, used to ensure received aggregates have correctly sized buffers
func (s *Signers) Validate(committee *Committee) error {
	if committee == nil {
		return errors.New("committee cannot be nil")
	}
	distinctSigners, rightmostSigner, maxCoefficient, err := s.validate(committee.Len())
	if err != nil {
		return err
	}

	s.committee = committee
	s.length = distinctSigners
	s.rightmostSigner = rightmostSigner
	s.maxCoefficient = maxCoefficient
	s.validated = true
	return nil

}

// validates the signer information and returns:
// 1. the number of distinct signers
// 2. the rightmost signer index
// 3. the max coefficient
// 4. error
// NOTE: it does not mutate the signers state
func (s *Signers) validate(committeeSize int) (int, int, *big.Int, error) {
	// length safety check
	if len(s.Coefficients) > committeeSize || s.Bitmap.Len() > committeeSize {
		return 0, committeeSize, nil, ErrWrongSizeSigners
	}

	// gather data about signers bits
	countNonZero := s.Bitmap.Count()

	// there has to be at least a signer
	if countNonZero == 0 {
		return 0, committeeSize, nil, ErrEmptySigners
	}

	// len(s.Coefficients) should be the same as the signer length
	// because each validator occupies only 1 bit in `s.Bitmap`
	if len(s.Coefficients) != countNonZero {
		return 0, committeeSize, nil, ErrWrongCoefficientLen
	}

	// if individual signature, its coefficient should be one
	if countNonZero == 1 && s.Coefficients[0].Cmp(common.Big1) != 0 {
		return 0, committeeSize, nil, ErrInvalidSingleSig
	}

	maxCoefficient := new(big.Int)
	sum := new(big.Int)
	for _, coefficient := range s.Coefficients {
		sum.Add(sum, coefficient)
		maxCoefficient = common.Max(maxCoefficient, coefficient)
	}

	// `maxCoefficient` <= (2^(countNonZero-2))
	if maxCoefficient.Cmp(common.Pow2(countNonZero-2)) > 0 {
		return 0, committeeSize, nil, ErrInvalidCoefficient
	}

	// `sum` <= (2^(countNonZero-1))
	if sum.Cmp(common.Pow2(countNonZero-1)) > 0 {
		return 0, committeeSize, nil, ErrInvalidCoefficient
	}

	rightmostSigner := s.Bitmap.RightmostIndex()
	return countNonZero, rightmostSigner, new(big.Int).Set(maxCoefficient), nil
}

func (s *Signers) Contains(index int) bool {
	if !s.validated {
		panic("Trying to use not validated signer information")
	}
	if index >= s.committee.Len() {
		panic("trying to call contains on non-existent committee member")
	}
	return s.Bitmap.IsSet(index)
}

func safetyCheck(first *Signers, second *Signers) error {
	if !first.validated || !second.validated {
		return ErrNotValidated
	}
	if first.committee.Len() != second.committee.Len() {
		return ErrDifferentSize
	}
	return nil
}

func (s *Signers) RightmostSigner() int {
	if !s.validated {
		panic("Trying to use not validated signer information")
	}
	return s.rightmostSigner
}

// check if `other` adds any information in `s`
func (s *Signers) AddsInformation(other *Signers) bool {
	if err := safetyCheck(s, other); err != nil {
		panic(err.Error())
	}
	return !s.Bitmap.Contains(other.Bitmap)
}

// this function will add the `coefficient` in `s.Coefficient` and will update the `s.Bitmap` if needed
// the caller is responsible to check if the coefficient does not overflow
func (s *Signers) addOrUpdate(validatorIndex, coeffIndex int, coefficient *big.Int) {
	if s.Bitmap.IsSet(validatorIndex) {
		s.Coefficients[coeffIndex].Add(s.Coefficients[coeffIndex], coefficient)
		s.maxCoefficient = new(big.Int).Set(common.Max(s.Coefficients[coeffIndex], s.maxCoefficient))
		return
	}

	if validatorIndex < s.rightmostSigner {
		s.rightmostSigner = validatorIndex
	}
	s.length++
	s.Bitmap.Set(validatorIndex)

	if coeffIndex == len(s.Coefficients) {
		s.Coefficients = append(s.Coefficients, new(big.Int).Set(coefficient))
	} else {
		s.Coefficients = append(s.Coefficients[:coeffIndex+1], s.Coefficients[coeffIndex:]...)
		s.Coefficients[coeffIndex] = new(big.Int).Set(coefficient)
	}
	s.maxCoefficient = new(big.Int).Set(common.Max(s.Coefficients[coeffIndex], s.maxCoefficient))
	s.powerMu.Lock()
	s.computedPower = nil // reset computed power, as it is no longer valid
	s.powerMu.Unlock()
}

// this function adds `index` in signer `s`. It assumes the signer is prevalidated.
// The caller is responsible to check so the coefficient doesn't overflow.
func (s *Signers) increment(index int) {

	count := 0 // count of signers present before `index`

	for i := 0; i < index; i++ {
		if s.Bitmap.IsSet(i) {
			count++
		}
	}

	s.addOrUpdate(index, count, common.Big1)
}

func (s *Signers) PowerByIndex(index int) *big.Int {
	if !s.validated {
		panic("Power has not been assigned in signers information")
	}
	if index >= s.committee.Len() {
		return common.Big0
	}
	return s.committee.MemberByIndex(index).VotingPower
}

// This function adds the `member` in signer `s`. This function assumes that `member` is absent in signer `s`.
// The caller is responsible to check if `member` is already present in `s` or not.
func (s *Signers) AddSigner(memberIndex uint64) {
	if !s.validated {
		panic("Using un-validated signers information")
	}
	index := int(memberIndex) // nolint
	if index >= s.committee.Len() {
		panic("trying to increment signer information of non-existent committee member")
	}

	s.increment(index)
}

// Iterates over the signers present in `other` and merge it with `s`.
// It assumes that they are mergeable. The caller is responsible to check
// so that coefficient overflow does not happen and the merge increases the total power in `s`.
func (s *Signers) Merge(other *Signers) {
	if err := safetyCheck(s, other); err != nil {
		panic(err.Error())
	}

	count1 := 0 // count of signers in `s`
	count2 := 0 // count of signers in `other`

	for i := 0; i < s.committee.Len(); i++ {
		if other.Bitmap.IsSet(i) {
			s.addOrUpdate(i, count1, other.Coefficients[count2])
			count2++
		}
		if s.Bitmap.IsSet(i) {
			count1++
		}
	}
}

// returns aggregated power of all senders
func (s *Signers) Power() *big.Int {
	if !s.validated {
		panic(ErrNotValidated.Error())
	}
	s.powerMu.RLock()
	if s.computedPower != nil {
		power := new(big.Int).Set(s.computedPower)
		s.powerMu.RUnlock()
		return power
	}
	s.powerMu.RUnlock()

	s.powerMu.Lock() // for writing
	if s.computedPower == nil {
		s.computedPower = new(big.Int)
		it := s.NewIterator()
		for it.Next() {
			s.computedPower.Add(s.computedPower, s.PowerByIndex(it.index))
		}
	}
	power := new(big.Int).Set(s.computedPower)
	s.powerMu.Unlock()
	return power
}

func (s *Signers) RespectsBoundaries(other *Signers) bool {
	if err := safetyCheck(s, other); err != nil {
		panic(err.Error())
	}

	var firstCoefficient *big.Int
	var firstCount int
	var secondCoefficient *big.Int
	var secondCount int

	for i := 0; i < s.committee.Len(); i++ {
		firstCoefficient = common.Big0
		if s.Bitmap.IsSet(i) {
			firstCoefficient = s.Coefficients[firstCount]
			firstCount++
		}

		secondCoefficient = common.Big0
		if other.Bitmap.IsSet(i) {
			secondCoefficient = other.Coefficients[secondCount]
			secondCount++
		}

		// TODO: optimize
		sum := new(big.Int).Add(firstCoefficient, secondCoefficient)
		if sum.BitLen() > common.VoteCap {
			return false
		}
	}

	return true
}

func (s *Signers) Copy() *Signers {
	coefficients := make([]*big.Int, 0, len(s.Coefficients))
	for _, coefficient := range s.Coefficients {
		coefficients = append(coefficients, new(big.Int).Set(coefficient))
	}
	// note: skip copying computedPower as it will be recomputed on demand,
	// committee is not copied as well because it is not meant to be altered and
	// only a reference in signers object, The idea is to keep the copy a lean operation
	return &Signers{
		Bitmap:          s.Bitmap.Copy(),
		Coefficients:    coefficients,
		maxCoefficient:  s.maxCoefficient,
		length:          s.length,
		rightmostSigner: s.rightmostSigner,
		committee:       s.committee,
		validated:       s.validated,
	}
}

func (s *Signers) CopyCoefficients() []*big.Int {
	if !s.validated {
		panic("Using un-validated signers information")
	}
	coeffs := make([]*big.Int, 0, len(s.Coefficients))
	for _, c := range s.Coefficients {
		coeffs = append(coeffs, new(big.Int).Set(c))
	}
	return coeffs
}

func (s *Signers) FlattenUniq() []int {
	if !s.validated {
		panic("Using un-validated signers information")
	}
	return s.flattenUniq(s.length, s.committee.Len())
}

// ForEachDistinctSigner iterates over each signer and calls the callback function.
func (s *Signers) ForEachDistinctSigner(callback func(signerIndex int)) {
	if !s.validated {
		panic("Using un-validated signers information")
	}

	it := s.NewIterator()
	for it.Next() {
		callback(it.index)
	}
}

// it is responsibility of the caller to pass the correct committee size
func (s *Signers) flattenUniq(distinctSigners, committeeSize int) []int {
	indexes := make([]int, 0, distinctSigners)

	for i := 0; i < committeeSize; i++ {
		if s.Bitmap.IsSet(i) {
			indexes = append(indexes, i)
		}
	}
	return indexes
}

// returns number of distinct signers of the aggregate
func (s *Signers) Len() int {
	if !s.validated {
		panic("Using un-validated signers information")
	}
	return s.length
}

func (s *Signers) String() string {
	return fmt.Sprintf("Bitmap: %s, Coefficients: %v, validated: %v", s.Bitmap.String(), s.Coefficients, s.validated)
}

func (s *Signers) CommitteeSize() int {
	if !s.validated {
		panic("Using un-validated signers information")
	}
	return s.committee.Len()
}

func (s *Signers) MaxCoefficient() *big.Int {
	if !s.validated {
		panic("Using un-validated signers information")
	}
	return s.maxCoefficient
}

func (s *Signers) AggregatePublicKey(keys []blst.PublicKey) blst.PublicKey {
	if !s.validated {
		panic("Using un-validated signers information")
	}
	return s.aggregatePublicKey(keys, s.length)
}

func (s *Signers) aggregatePublicKey(keys []blst.PublicKey, length int) blst.PublicKey {
	if len(keys) != length {
		panic("invalid public key length")
	}

	return blst.AggregatePublicKeysMultScalars(keys, blst.ToScalars(s.Coefficients))
}

func (s *Signers) NewIterator() *BitmapIterator {
	return NewBitmapIterator(s.Bitmap)
}
