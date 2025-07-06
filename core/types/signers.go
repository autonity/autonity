package types

import (
	"errors"
	"fmt"
	"math/big"
	"math/bits" //nolint
	"reflect"

	blstbind "github.com/supranational/blst/bindings/go"

	"github.com/autonity/autonity/common/bitutil"
	"github.com/autonity/autonity/crypto/blst"
)

const (
	bitsPerValidator   = 1 //NOTE: if this gets changed, major refactoring will be needed for this file. Proceed with caution.
	bitsInByte         = 8
	validatorsPerByte  = bitsInByte / bitsPerValidator
	maxUint16          = (1 << 16) - 1
	minSafeCoefficient = maxUint16 / 2 // because adding two coefficient with value `minSafeCoefficient` will not overflow uint16
)

var (
	// // 11000000 00110000 00001100 00000011
	// getMasks = []byte{0xC0, 0x30, 0x0C, 0x03}
	// // 00111111 11001111 11110011 11111100
	// setMasks  = []byte{0x3F, 0xCF, 0xF3, 0xFC}
	maxAllowedSigners = 16384

	ErrNilSigners          = errors.New("validator bitmap or coefficient array is nil")
	ErrWrongSizeSigners    = errors.New("validator bitmap or coefficient array has incorrect size")
	ErrEmptySigners        = errors.New("signers information is empty")
	ErrWrongCoefficientLen = errors.New("coefficient array has incorrect length")
	ErrInvalidSingleSig    = errors.New("individual signature has coefficient != 1")
	ErrInvalidCoefficient  = errors.New("coefficient exceeds maximum boundary (committee size)")

	ErrNotValidated  = errors.New("using un-validated signers information")
	ErrDifferentSize = errors.New("comparing signers information with different committee size")
)

// represents the senders of an aggregated signature

/*
* Two bits for each validator. The meaning is:
* 00 --> no signature from the validator
* 01 --> 1 signature
* 10 --> 2 signatures
* 11 --> look at the number of signatures in the `Coefficients` array
*
* example:
* indexes 			 0  1  2  3  4  5
* bits    			 00 11 10 01 00 11
* coefficients   17 170
* #sigs          0  17 2  1  0  170
 */

type validatorBitmap []byte

func NewValidatorBitmap(committeeSize int) validatorBitmap { //nolint
	byteLength := (committeeSize*bitsPerValidator + bitsInByte - 1) / bitsInByte
	return make(validatorBitmap, byteLength)
}

// ensures that the validator bitmap has the correct length compared to the committee size
// used to validate aggregate messages coming from other peers
func (vb validatorBitmap) Valid(committeeSize int) bool {
	expectedByteLength := (committeeSize*bitsPerValidator + bitsInByte - 1) / bitsInByte
	return len(vb) == expectedByteLength
}

// Do not call this function on the signers' bitmap. Create a new `validatorBitmap`, `vb`
// to call `vb.Merge`. `other` can be from some signers' bitmap as it's not modified.
// returns true if the `other` object contributes to the `vb` object
func (vb validatorBitmap) Merge(other validatorBitmap) bool {
	if len(vb) != len(other) {
		panic(ErrDifferentSize.Error())
	}
	contributed := false
	for i, b := range other {
		if (vb[i] & b) != b {
			contributed = true
		}
		vb[i] |= b
	}
	return contributed
}

// TODO: remove
// func (vb validatorBitmap) Get(validatorIndex int) byte {
// 	byteIndex := validatorIndex / validatorsPerByte
// 	bitIndex := validatorIndex % validatorsPerByte
// 	b := vb[byteIndex]

// 	result := b & getMasks[bitIndex]
// 	shift := validatorsPerByte - 1 - bitIndex
// 	result = result >> (shift * bitsPerValidator)
// 	return result
// }

// this function needs to be modified if the constant `validatorsPerByte` is changed
func (vb validatorBitmap) HasSigner(validatorIndex int) bool {
	byteIndex, bitIndex := indexToBitMapPosition(validatorIndex)

	// because of the constant `bitsPerValidator = 1`, each validator takes a single bit
	// we just need to check if the bit is 1 or 0
	// note that the bits are numbered from LSB to MSB (in both `HasSigner` and `setSigner`)
	return vb.hasBit(byteIndex, bitIndex)
}

func (vb validatorBitmap) hasBit(byteIndex, bitIndex int) bool {
	return (vb[byteIndex] & (1 << bitIndex)) > 0
}

// TODO: remove
// // NOTE: be careful when calling directly this function without passing through the `increment` function.
// // this function will not invalidate any cache, it is just a naive setter
// func (vb validatorBitmap) Set(validatorIndex int, value byte) {
// 	if value > multipleSignatures {
// 		panic("Trying to set value that cannot fit into 2 bits")
// 	}
// 	bitIndex := validatorIndex % validatorsPerByte
// 	shift := validatorsPerByte - 1 - bitIndex
// 	valueShifted := value << (shift * bitsPerValidator)

// 	byteIndex := validatorIndex / validatorsPerByte
// 	vb[byteIndex] = vb[byteIndex] & setMasks[bitIndex]
// 	vb[byteIndex] = vb[byteIndex] | valueShifted
// }

// NOTE: be careful when calling directly this function without passing through the `increment` function.
// this function will not invalidate any cache, it is just a naive setter
// this function needs to be modified if the constant `validatorsPerByte` is changed
func (vb validatorBitmap) setSigner(validatorIndex int) {
	byteIndex, bitIndex := indexToBitMapPosition(validatorIndex)

	// because of the constant `bitsPerValidator = 1`, each validator takes a single bit
	// we just need to set 1 in `bitIndex`
	// note that the bits are numbered from LSB to MSB (in both `HasSigner` and `setSigner`)
	vb.setBit(byteIndex, bitIndex)
}

func (vb validatorBitmap) setBit(byteIndex, bitIndex int) {
	vb[byteIndex] = vb[byteIndex] | (1 << bitIndex)
}

// it should support `VoteSigners` for `T = uint16` and `QuorumSigners` for `T = uint32`
type SignersBase[T uint16 | uint32] struct {
	Bits         validatorBitmap
	Coefficients []T // support up to 65535 committee members

	// these fields are not serialized, but instead computed at preValidate steps
	committeeSize  int              `rlp:"-"`
	length         int              `rlp:"-"` // number of distinct signers
	powers         map[int]*big.Int `rlp:"-"`
	power          *big.Int         `rlp:"-"` // aggregated power of all senders
	maxCoefficient T                `rlp:"-"`

	// auxiliary data structures flags
	validated     bool `rlp:"-"` // if true --> Bits and Coefficients have correct length + committeeSize and length is assigned
	powerAssigned bool `rlp:"-"` // if true --> powers and power assigned
}

func NewSigners[T uint16 | uint32](committeeSize int) *SignersBase[T] {
	if committeeSize > maxUint16 {
		panic("Unsupported committee size")
	}
	return &SignersBase[T]{
		Bits:          NewValidatorBitmap(committeeSize),
		committeeSize: committeeSize,

		length:        0,
		powers:        make(map[int]*big.Int),
		power:         new(big.Int),
		validated:     true,
		powerAssigned: true, // when we are locally creating a sender info, we are ok with power being 0 initially

		Coefficients:   make([]T, 0),
		maxCoefficient: 0,
	}
}

// Correct full validation cannot be done until we know the committee size of this block,
// but we can already ensure that the fields have sane values
func (s *SignersBase[T]) SanityCheck() error {
	if len(s.Bits) == 0 {
		return fmt.Errorf("validator bitmap is empty")
	}
	// 1 byte --> 8 validators
	if len(s.Bits) > maxAllowedSigners/validatorsPerByte { // max 2kb of data
		return fmt.Errorf("invalid Bits length: %d", len(s.Bits))
	}
	if len(s.Coefficients) > maxAllowedSigners { // max 32kb or 64kb of data
		return fmt.Errorf("invalid Coefficient length: %d", len(s.Coefficients))
	}
	return nil
}

// validates the sender info, used to ensure received aggregates have correctly sized buffers
func (s *SignersBase[T]) Validate(committeeSize int) error {
	countNonZero, maxCoefficient, err := s.validate(committeeSize)
	if err != nil {
		return err
	}

	s.committeeSize = committeeSize
	s.length = countNonZero
	s.maxCoefficient = maxCoefficient
	s.validated = true
	return nil

}

// validates the signer information and returns the number of distinct signers
// it does not mutate the signers state
func (s *SignersBase[T]) validate(committeeSize int) (int, T, error) {
	// whether locally created or received from wire, Bits are never nil
	if s.Bits == nil || s.Coefficients == nil {
		return 0, 0, ErrNilSigners
	}

	// length safety check
	if len(s.Coefficients) > committeeSize || !s.Bits.Valid(committeeSize) {
		return 0, 0, ErrWrongSizeSigners
	}

	// gather data about signers bits
	countNonZero := 0
	for _, b := range s.Bits {
		if b == 0 {
			continue
		}
		countNonZero += bits.OnesCount8(b)
	}

	// there has to be at least a signer
	if countNonZero == 0 {
		return 0, 0, ErrEmptySigners
	}
	if countNonZero > committeeSize {
		return 0, 0, ErrWrongSizeSigners
	}

	// len(s.Coefficients) should be the same as the signer length
	// because each validator occupies only 1 bit in `s.Bits`
	if len(s.Coefficients) != countNonZero {
		return 0, 0, ErrWrongCoefficientLen
	}

	// if individual signature, its coefficient should be one (01)
	if countNonZero == 1 && s.Coefficients[0] != 1 {
		return 0, 0, ErrInvalidSingleSig
	}

	var maxCoefficient T
	for _, coefficient := range s.Coefficients {
		if coefficient == 0 {
			return 0, 0, ErrInvalidCoefficient
		}
		maxCoefficient = max(maxCoefficient, coefficient)
	}

	// the following check is needed only for `VoteSigners`

	if reflect.TypeOf(maxCoefficient) == reflect.TypeOf(uint16(0)) {
		// check that `maxCoefficient` respects the maximum allowed boundary (2^(s.length-2))
		// if `s.length >= 18`, maximum coefficient can be over `maxUint16`, so we can skip the check
		if s.length < 18 && s.length > 1 {
			if maxCoefficient > (1 << (s.length - 2)) {
				return 0, 0, ErrInvalidCoefficient
			}
		}
	}
	return countNonZero, maxCoefficient, nil
}

func (s *SignersBase[T]) Contains(index int) bool {
	if !s.validated {
		panic("Trying to use not validated signer information")
	}
	if index >= s.committeeSize {
		panic("trying to call contains on non-existent committee member")
	}
	return s.Bits.HasSigner(index)
}

func (s *SignersBase[T]) safetyCheck(other *SignersBase[T]) error {
	if !s.validated || !other.validated {
		return ErrNotValidated
	}
	if s.committeeSize != other.committeeSize {
		return ErrDifferentSize
	}
	return nil
}

func (s *SignersBase[uint16]) canAddCoefficient(other *SignersBase[uint16], position1, position2 int) bool {
	// check if `s.Coefficients[position1] + other.Coefficients[position2] <= maxUint16`
	// checking if `s.Coefficients[position1] <= (maxUint16 - other.Coefficients[position2])`

	// `maxUint16 = (1<<16) - 1`, so all the 16 positions in it has 1
	// so `maxUint16 - n` equals to `maxUint16 ^ n` if `n <= maxUint16`
	return s.Coefficients[position1] <= (maxUint16 ^ other.Coefficients[position2])
}

// checks that the resulting aggregate still respects the logical boundary for `VoteSigners`
// we don't need this method for `QuorumSigners`
func (s *SignersBase[uint16]) RespectsBoundaries(other *SignersBase[uint16]) bool {
	if err := s.safetyCheck(other); err != nil {
		panic(err.Error())
	}

	// check that addition of coefficient does not overflow uint16
	if s.maxCoefficient <= minSafeCoefficient && other.maxCoefficient <= minSafeCoefficient {
		// adding two coefficient in these signers cannot overflow uint16
		return true
	}
	// one of the signers has coefficient > minSafeCoefficient

	// We need to check if some signer has coefficient in both `s` and `other`
	// and if their sum exceeds `maxUint16`

	// if some signer coefficient <= minSafeCoefficient in both `s` and `other`,
	// then no need to check the sum as it `minSafeCoefficient * 2 < maxUint16`

	// because constant `bitsPerValidator = 1`, we only count the number of 1 bits to count signers
	// in case constant `bitsPerValidator` is changed, this block of codes needs refactoring
	count1 := 0 // number of signers in `s.Bits`
	count2 := 0 // number of signers in `other.Bits`
	for i, b := range s.Bits {
		commonBits := b & other.Bits[i]
		for commonBits > 0 {
			// we have some signers present in both `s` and `other` which are present in `commonBits`

			// we will take the signer at the LSB of `commonBits`
			// so we need to know how many signers are present in `s.Bits[i]`
			// and `other.Bits[i]` separately in lower bits than the LSB of `commonBits`

			signerAtLowerBits1 := s.signerCountBeforeLSB(i, commonBits) + count1
			signerAtLowerBits2 := other.signerCountBeforeLSB(i, commonBits) + count2

			// if both coefficients are <= minSafeCoefficient than no need to check
			if !s.canAddCoefficient(other, signerAtLowerBits1, signerAtLowerBits2) {
				return false
			}

			// then we remove the LSB
			commonBits = commonBits & (commonBits - 1)

			// the above idea is inspired from Brian Kernighan's Algorithm
			// see more here: https://how.dev/answers/what-is-kernighans-algorithm
		}
		// the total complexity of the above loop is `O(number_of_common_signers)`

		// update `count1`, `count2`
		if b > 0 {
			count1 += bits.OnesCount8(b)
		}
		if other.Bits[i] > 0 {
			count2 += bits.OnesCount8(other.Bits[i])
		}
	}
	// the total complexity of the nested loops is `O(len(s.bits)) + O(number_of_common_signers)`
	return true
}

// check if `other` adds any information in `s`
func (s *SignersBase[T]) AddsInformation(other *SignersBase[T]) bool {
	if err := s.safetyCheck(other); err != nil {
		panic(err.Error())
	}

	// can we use `bitutil.ANDBytes` instead?
	for i, b := range other.Bits {
		// check if `other.Bits[i]` is a subset of `s.Bits[i]`
		if (b & s.Bits[i]) != b {
			// `b` is not a subset of `s.Bits[i]`
			return true
		}
	}
	return false
}

// returns the count of signers in `s` at `s.Bits[byteIndex]` before `bitIndex`
func (s *SignersBase[T]) signerCountBeforeIndex(byteIndex, bitIndex int) int {
	return bitutil.ByteOnesCountLowerIndex(s.Bits[byteIndex], bitIndex)
}

// returns the count of signers in `s` at `s.Bits[byteIndex]` before the
// LSB of `b`
func (s *SignersBase[T]) signerCountBeforeLSB(byteIndex int, b byte) int {
	return bitutil.ByteOnesCountLowerLSB(s.Bits[byteIndex], b)
}

// this function will add the `coefficient` in `s.Coefficient` and will update the `s.Bits` if needed
// the caller is responsible to check if the coefficient does not overflow
func (s *SignersBase[T]) addOrUpdate(
	byteIndex, bitIndex, coeffIndex, validatorIndex int,
	coefficient T,
	votingPower *big.Int,
) {
	if s.Bits.hasBit(byteIndex, bitIndex) {
		s.Coefficients[coeffIndex] += coefficient
		s.maxCoefficient = max(s.maxCoefficient, s.Coefficients[coeffIndex])
		return
	}

	s.length++
	s.Bits.setBit(byteIndex, bitIndex)

	if coeffIndex == len(s.Coefficients) {
		s.Coefficients = append(s.Coefficients, coefficient)
	} else {
		s.Coefficients = append(s.Coefficients[:coeffIndex+1], s.Coefficients[coeffIndex:]...)
		s.Coefficients[coeffIndex] = coefficient
	}
	s.maxCoefficient = max(s.maxCoefficient, s.Coefficients[coeffIndex])

	s.powers[validatorIndex] = votingPower
	s.power.Add(s.power, votingPower)
}

// this function adds `index` in signer `s`. It assumes the signer is prevalidated.
// The caller is responsible to check so the coefficient doesn't overflow.
func (s *SignersBase[T]) increment(index int, coefficient T, votingPower *big.Int) {

	count := 0 // count of signers present before `index`

	byteIndex, bitIndex := indexToBitMapPosition(index)
	// count the number of signers present before `byteIndex`
	for i := 0; i < byteIndex; i++ {
		count += bits.OnesCount8(s.Bits[i])
	}

	// now count the number of signers in `byteIndex` before `bitIndex`
	count += s.signerCountBeforeIndex(byteIndex, bitIndex)

	s.addOrUpdate(byteIndex, bitIndex, count, index, coefficient, votingPower)
}

// This function adds the `member` in signer `s`. This function assumes that `member` is absent in signer `s`.
// The caller is responsible to check if `member` is already present in `s` or not.
func (s *SignersBase[T]) AddMember(member *CommitteeMember) {
	if !s.validated {
		panic("Using un-validated signers information")
	}
	if !s.powerAssigned {
		panic("Power has not been assigned in signers information")
	}
	if int(member.Index) >= s.committeeSize {
		panic("trying to increment signer information of non-existent committee member")
	}

	s.increment(int(member.Index), 1, member.VotingPower)
}

// Iterates over the signers present in `other` and merge it with `s`.
// It assumes that they are mergeable. The caller is responsible to check
// so that coefficient overflow does not happen and the merge increases the total power in `s`.
func (s *SignersBase[T]) Merge(other *SignersBase[T]) {
	if err := s.safetyCheck(other); err != nil {
		panic(err.Error())
	}
	if !s.powerAssigned || !other.powerAssigned {
		panic("Power has not been assigned in signers information")
	}

	count1 := 0 // count of signers in `s`
	count2 := 0 // count of signers in `other`
	for i, b := range other.Bits {
		for b > 0 {
			// take the `bitIndex` of the LSB of `b`
			bitIndex := bitutil.ByteLSBPosition(b)
			signerAtLowerBits1 := s.signerCountBeforeIndex(i, bitIndex) + count1
			signerAtLowerBits2 := other.signerCountBeforeIndex(i, bitIndex) + count2

			validatorIndex := bitMapPositionToIndex(i, bitIndex)

			s.addOrUpdate(
				i, bitIndex, signerAtLowerBits1, validatorIndex,
				other.Coefficients[signerAtLowerBits2],
				other.powers[validatorIndex],
			)

			// remove the LSB of `b`
			b = b & (b - 1)
			// the above idea is inspired from Brian Kernighan's Algorithm
			// see more here: https://how.dev/answers/what-is-kernighans-algorithm
		}

		// update `count1`, `count2`
		count1 += bits.OnesCount8(s.Bits[i])
		count2 += bits.OnesCount8(other.Bits[i])
	}
}

// returns aggregated power of all senders
func (s *SignersBase[T]) Power() *big.Int {
	if !s.powerAssigned {
		panic("Power has not been assigned in signers information")
	}
	return s.power
}

func (s *SignersBase[T]) AssignPower(powers map[int]*big.Int, power *big.Int) {
	s.powers = powers
	s.power = power
	s.powerAssigned = true
}

func (s *SignersBase[T]) Copy() *SignersBase[T] {
	var powers map[int]*big.Int
	if s.powers != nil {
		powers = make(map[int]*big.Int, len(s.powers))
		for index, power := range s.powers {
			powers[index] = new(big.Int).Set(power)
		}
	}
	var power *big.Int
	if s.power != nil {
		power = new(big.Int).Set(s.power)
	}
	return &SignersBase[T]{
		Bits:           append(s.Bits[:0:0], s.Bits...),
		Coefficients:   append(s.Coefficients[:0:0], s.Coefficients...),
		maxCoefficient: s.maxCoefficient,
		committeeSize:  s.committeeSize,
		length:         s.length,
		powers:         powers,
		power:          power,
		validated:      s.validated,
		powerAssigned:  s.powerAssigned,
	}
}

func (s *SignersBase[T]) CopyCoefficients() []T {
	if !s.validated {
		panic("Using un-validated signers information")
	}
	coeffs := make([]T, len(s.Coefficients))
	copy(coeffs, s.Coefficients)
	return coeffs
}

// same as before, but repeated indexes are returned only once
func (s *SignersBase[T]) FlattenUniq() []int {
	if !s.validated {
		panic("Using un-validated signers information")
	}
	return s.flattenUniq()
}

// it is responsibility of the caller to pass the correct committee size
func (s *SignersBase[T]) flattenUniq() []int {
	indexes := make([]int, 0, s.length)
	for i, b := range s.Bits {
		for b > 0 {
			// get the `bitIndex` of the LSB of `b`
			bitIndex := bitutil.ByteLSBPosition(b)
			indexes = append(indexes, bitMapPositionToIndex(i, bitIndex))

			// remove the LSB of `b`
			b = b & (b - 1)
			// the above idea is inspired from Brian Kernighan's Algorithm
			// see more here: https://how.dev/answers/what-is-kernighans-algorithm
		}
	}
	return indexes
}

// returns number of distinct signers of the aggregate
func (s *SignersBase[T]) Len() int {
	if !s.validated {
		panic("Using un-validated signers information")
	}
	return s.length
}

// returns whether an aggregate is:
//   - a simple aggregate (all coefficients are 0 or 1)
//   - a complex aggregate (at least one coefficient is > 1)
func (s *SignersBase[T]) IsComplex() bool {
	if !s.validated {
		panic("Using un-validated signers information")
	}
	return s.maxCoefficient > 1
}

func (s *SignersBase[T]) String() string {
	return fmt.Sprintf("Bits: %08b, Coefficients: %v, power: %v, validated: %v, powerAssigned: %v", s.Bits, s.Coefficients, s.power, s.validated, s.powerAssigned)
}

func (s *SignersBase[T]) Powers() map[int]*big.Int {
	if !s.powerAssigned {
		panic("Power has not been assigned in signers information")
	}
	return s.powers
}

func (s *SignersBase[T]) CommitteeSize() int {
	if !s.validated {
		panic("Using un-validated signers information")
	}
	return s.committeeSize
}

func (s *SignersBase[T]) MaxCoefficient() T {
	if !s.validated {
		panic("Using un-validated signers information")
	}
	return s.maxCoefficient
}

func (s *SignersBase[T]) AggregatePublicKey(keys []blst.PublicKey) blst.PublicKey {
	if !s.validated {
		panic("Using un-validated signers information")
	}

	return s.aggregatePublicKey(keys, s.maxCoefficient, s.length)
}

func (s *SignersBase[T]) aggregatePublicKey(keys []blst.PublicKey, maxCoefficient T, length int) blst.PublicKey {
	if len(keys) != length {
		panic("invalid public key length")
	}

	var bitsEntropy int
	if reflect.TypeOf(maxCoefficient) == reflect.TypeOf(uint16(0)) {
		bitsEntropy = bitutil.Uint16MSBPosition(uint16(maxCoefficient)) + 1
	} else {
		bitsEntropy = bitutil.Uint32MSBPosition(uint32(maxCoefficient)) + 1
	}

	return blst.AggregatePublicKeysMultScalars(
		keys,
		s.toBlstScalars(),
		bitsEntropy,
	)
}

func (s *SignersBase[T]) ToBlstScalars() []*blstbind.Scalar {
	if !s.validated {
		panic("Using un-validated signers information")
	}
	return s.toBlstScalars()
}

func (s *SignersBase[T]) toBlstScalars() []*blstbind.Scalar {
	return blst.ToBlstScalars(s.Coefficients)
}

func (s *SignersBase[T]) ToQuorumSigners() *QuorumSigners {
	if !s.validated {
		panic("Trying to use not validated signer information")
	}
	coefficient := make([]uint32, s.length)
	for i, c := range s.Coefficients {
		coefficient[i] = uint32(c)
	}
	duplicate := s.Copy()
	return &QuorumSigners{
		SignersBase: &SignersBase[uint32]{
			Bits:          duplicate.Bits,
			committeeSize: duplicate.committeeSize,
			length:        duplicate.length,
			powers:        duplicate.powers,
			power:         duplicate.power,
			validated:     duplicate.validated,
			powerAssigned: duplicate.powerAssigned,

			Coefficients:   coefficient,
			maxCoefficient: uint32(duplicate.maxCoefficient),
		},
	}
}

type VoteSigners struct {
	*SignersBase[uint16]
}

func NewVoteSigners(committeeSize int) *VoteSigners {
	return &VoteSigners{
		SignersBase: NewSigners[uint16](committeeSize),
	}
}

type QuorumSigners struct {
	*SignersBase[uint32]
}

func NewQuorumSigners(committeeSize int) *QuorumSigners {
	return &QuorumSigners{
		SignersBase: NewSigners[uint32](committeeSize),
	}
}

func indexToBitMapPosition(index int) (int, int) {
	// the following line is the same as `index / validatorsPerByte`
	// but the following works because `validatorsPerByte = 8 = 2^3`
	byteIndex := index >> 3

	// the following line is the same as `index % validatorsPerByte`
	// but the following works because `validatorsPerByte = 8`, which is a power of 2
	bitIndex := index & (validatorsPerByte - 1)
	return byteIndex, bitIndex
}

func bitMapPositionToIndex(byteIndex, bitIndex int) int {
	// the following line is the same as `byteIndex * validatorsPerByte + bitIndex`
	// but the following works because `validatorsPerByte = 8 = 2^3` and `bitIndex < validatorsPerByte`
	// `bitIndex < validatorsPerByte` should always be true
	return (byteIndex << 3) | bitIndex
}
