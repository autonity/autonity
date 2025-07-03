package types

import (
	"errors"
	"fmt"
	"math/big"
	"math/bits"

	"github.com/autonity/autonity/common/bitutil"
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

type VoteSigners struct {
	*SignersBase
	Coefficients   []uint16 // support up to 65535 committee members
	MaxCoefficient uint16
}

type QuorumSigners struct {
	*SignersBase
	Coefficients   []uint32 // support up to 65535 committee members
	MaxCoefficient uint32
}

type SignersBase struct {
	Bits validatorBitmap

	// these fields are not serialized, but instead computed at preValidate steps
	committeeSize int              `rlp:"-"`
	length        int              `rlp:"-"` // number of distinct signers
	powers        map[int]*big.Int `rlp:"-"`
	power         *big.Int         `rlp:"-"` // aggregated power of all senders

	// auxiliary data structures flags
	validated     bool `rlp:"-"` // if true --> Bits and Coefficients have correct length + committeeSize and length is assigned
	powerAssigned bool `rlp:"-"` // if true --> powers and power assigned
}

func NewSigners(committeeSize int) *VoteSigners {
	if committeeSize > maxUint16 {
		panic("Unsupported committee size")
	}
	return &VoteSigners{
		SignersBase: &SignersBase{
			Bits:          NewValidatorBitmap(committeeSize),
			committeeSize: committeeSize,

			length:        0,
			powers:        make(map[int]*big.Int),
			power:         new(big.Int),
			validated:     true,
			powerAssigned: true, // when we are locally creating a sender info, we are ok with power being 0 initially
		},
		Coefficients:   make([]uint16, 0),
		MaxCoefficient: 0,
	}
}

type validatorBitmap []byte

func NewValidatorBitmap(committeeSize int) validatorBitmap {
	byteLength := (committeeSize*bitsPerValidator + bitsInByte - 1) / bitsInByte
	return make(validatorBitmap, byteLength)
}

// ensures that the validator bitmap has the correct length compared to the committee size
// used to validate aggregate messages coming from other peers
func (vb validatorBitmap) Valid(committeeSize int) bool {
	expectedByteLength := (committeeSize*bitsPerValidator + bitsInByte - 1) / bitsInByte
	return len(vb) == expectedByteLength
}

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
	// the following line is the same as `validatorIndex / validatorsPerByte`
	// but the following works because `validatorsPerByte = 8 = 2^3`
	byteIndex := validatorIndex >> 3

	// the following line is the same as `validatorIndex % validatorsPerByte`
	// but the following works because `validatorsPerByte = 8`, which is a power of 2
	bitIndex := validatorIndex & (validatorsPerByte - 1)

	// because of the constant `bitsPerValidator = 1`, each validator takes a single bit
	// we just need to check if the bit is 1 or 0
	// note that the bits are numbered from LSB to MSB (in both `HasSigner` and `SetSigner`)
	result := vb[byteIndex] & (1 << bitIndex)
	return result > 0
}

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
func (vb validatorBitmap) SetSigner(validatorIndex int) {
	// the following line is the same as `validatorIndex % validatorsPerByte`
	// but the following works because `validatorsPerByte = 8`, which is a power of 2
	bitIndex := validatorIndex & (validatorsPerByte - 1)

	// the following line is the same as `validatorIndex / validatorsPerByte`
	// but the following works because `validatorsPerByte = 8 = 2^3`
	byteIndex := validatorIndex >> 3
	// because of the constant `bitsPerValidator = 1`, each validator takes a single bit
	// we just need to set 1 in `bitIndex`
	// note that the bits are numbered from LSB to MSB (in both `HasSigner` and `SetSigner`)
	vb[byteIndex] = vb[byteIndex] | (1 << bitIndex)
}

func (vb validatorBitmap) ToSingleBitmap(committeeSize int) []byte {
	oneBitmapLength := (committeeSize + bitsInByte - 1) / bitsInByte
	oneBitmap := make([]byte, oneBitmapLength)

	for i := 0; i < committeeSize; i++ {
		byteIndex := i / validatorsPerByte
		bitIndex := i % validatorsPerByte
		shift := (validatorsPerByte - 1 - bitIndex) * bitsPerValidator
		value := (vb[byteIndex] >> shift) & 0x03 // Extract 2 bits
		if value > 0 {
			newByteIndex := i / 8
			newBitIndex := i % 8
			oneBitmap[newByteIndex] |= 1 << (7 - newBitIndex)
		}
	}
	return oneBitmap
}

// Correct full validation cannot be done until we know the committee size of this block,
// but we can already ensure that the fields have sane values
func (s *VoteSigners) SanityCheck() error {
	if err := s.sanityCheck(); err != nil {
		return err
	}
	if len(s.Coefficients) > maxAllowedSigners { // max 32kb of data
		return fmt.Errorf("invalid Coefficient length: %d", len(s.Coefficients))
	}
	return nil
}

func (s *QuorumSigners) SanityCheck() error {
	if err := s.sanityCheck(); err != nil {
		return err
	}
	if len(s.Coefficients) > maxAllowedSigners { // max 64kb of data
		return fmt.Errorf("invalid Coefficient length: %d", len(s.Coefficients))
	}
	return nil
}

func (s *SignersBase) sanityCheck() error {
	if len(s.Bits) == 0 {
		return fmt.Errorf("validator bitmap is empty")
	}
	// 1 byte --> 8 validators
	if len(s.Bits) > maxAllowedSigners/validatorsPerByte { // max 2kb of data
		return fmt.Errorf("invalid Bits length: %d", len(s.Bits))
	}
	return nil
}

// validates the sender info, used to ensure received aggregates have correctly sized buffers
func (s *VoteSigners) Validate(committeeSize int) error {
	if err := s.validate(committeeSize); err != nil {
		return err
	}

	// whether locally created or received from wire, Coefficients are never nil
	if s.Coefficients == nil {
		return ErrNilSigners
	}

	// length safety check
	if len(s.Coefficients) > committeeSize {
		return ErrWrongSizeSigners
	}

	// len(s.Coefficients) should be the same as the signer length
	// because each validator occupies only 1 bit in `s.Bits`
	if len(s.Coefficients) != s.length {
		return ErrWrongCoefficientLen
	}

	// if individual signature, its coefficient should be one (01)
	if s.length == 1 && s.Coefficients[0] != 1 {
		return ErrInvalidSingleSig
	}

	var maxCoefficient uint16 = 0
	for _, coefficient := range s.Coefficients {
		maxCoefficient = max(maxCoefficient, coefficient)
	}

	// check that `maxCoefficient` respects the maximum allowed boundary (2^(s.length-2))
	// if `s.length >= 18`, maximum coefficient can be over `maxUint16`, so we can skip the check
	if s.length < 18 && s.length > 1 {
		if maxCoefficient > (1 << (s.length - 2)) {
			return ErrInvalidCoefficient
		}
	}

	s.MaxCoefficient = maxCoefficient
	s.validated = true
	return nil

}

// validates the signer information and returns the number of distinct signers
// it does not mutate the signers state
func (s *SignersBase) validate(committeeSize int) error {
	// whether locally created or received from wire, Bits are never nil
	if s.Bits == nil {
		return ErrNilSigners
	}

	// length safety check
	if !s.Bits.Valid(committeeSize) {
		return ErrWrongSizeSigners
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
		return ErrEmptySigners
	}
	if countNonZero > committeeSize {
		return ErrWrongSizeSigners
	}

	s.committeeSize = committeeSize
	s.length = countNonZero
	return nil
}

func (s *SignersBase) Contains(index int) bool {
	if !s.validated {
		panic("Trying to use not validated signer information")
	}
	if index >= s.committeeSize {
		panic("trying to call contains on non-existent committee member")
	}
	return s.Bits.HasSigner(index)
}

func safetyCheck(first *SignersBase, second *SignersBase) error {
	if !first.validated || !second.validated {
		return ErrNotValidated
	}
	if first.committeeSize != second.committeeSize {
		return ErrDifferentSize
	}
	return nil
}

func (s *VoteSigners) canAddCoefficient(other *VoteSigners, position1, position2 int) bool {
	// check if `s.Coefficients[position1] + other.Coefficients[position2] <= maxUint16`
	// checking if `s.Coefficients[position1] <= (maxUint16 - other.Coefficients[position2])`

	// `maxUint16 = (1<<16) - 1`, so all the 16 positions in it has 1
	// so `maxUint16 - n` equals to `maxUint16 ^ n` if `n <= maxUint16`
	return s.Coefficients[position1] <= (maxUint16 ^ other.Coefficients[position2])
}

// checks that the resulting aggregate still respects the `committeeSize` boundary
func (s *VoteSigners) RespectsBoundaries(other *VoteSigners) bool {
	if err := safetyCheck(s.SignersBase, other.SignersBase); err != nil {
		panic(err.Error())
	}

	// check that addition of coefficient does not overflow uint16
	if s.MaxCoefficient <= minSafeCoefficient && other.MaxCoefficient <= minSafeCoefficient {
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

			signerAtLowerBits1 := bitutil.OnesCountLowerLSB8(b, commonBits) + count1
			signerAtLowerBits2 := bitutil.OnesCountLowerLSB8(other.Bits[i], commonBits) + count2

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
func (s *SignersBase) AddsInformation(other *SignersBase) bool {
	if err := safetyCheck(s, other); err != nil {
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

// this function adds `index` in signer `s`. It assumes the signer is prevalidated.
// The caller is responsible to check so the coefficient doesn't overflow.
// returns `true` if the increment contributes to the total power
func (s *VoteSigners) increment(index int, coefficient uint16) bool {

	count := 0 // count of signers present before `index`

	// we have `validatorByteInde = index / validatorsPerByte`
	// but because the constant `validatorsPerByte = 8 = 2^3`,
	// the following line is the same as `index / validatorsPerByte`
	validatorByteIndex := index >> 3
	// count the number of signers present before `validatorByteIndex`
	for i := 0; i < validatorByteIndex; i++ {
		count += bits.OnesCount8(s.Bits[i])
	}

	// the following line is the same as `index % validatorsPerByte`
	// but the following works because `validatorsPerByte = 8`, which is a power of 2
	validatorBitIndex := index & (validatorsPerByte - 1)

	// now count the number of signers in `validatorByteIndex` before `validatorBitIndex`
	count += bitutil.OnesCountBeforeIndex8(s.Bits[validatorByteIndex], uint8(validatorBitIndex))

	if s.Bits.HasSigner(index) {
		s.Coefficients[count] += coefficient
		return false
	}

	s.length++
	s.Bits.SetSigner(index)
	if count == len(s.Coefficients) {
		s.Coefficients = append(s.Coefficients, coefficient)
	} else {
		s.Coefficients = append(s.Coefficients[:count+1], s.Coefficients[count:]...)
		s.Coefficients[count] = coefficient
	}
	return true
}

// this function adds `member` in signer `s`. It assumes the signer is prevalidated.
// The caller is responsible to check so the coefficient doesn't overflow.
func (s *VoteSigners) addMember(member *CommitteeMember, count uint16) {
	index := int(member.Index)
	if s.increment(index, count) {
		s.powers[index] = member.VotingPower
		s.power.Add(s.power, member.VotingPower)
	}
}

// This function adds the `member` in signer `s`. This function assumes that `member` is absent in signer `s`.
// The caller is responsible to check if `member` is already present in `s` or not.
func (s *VoteSigners) AddMember(member *CommitteeMember) {
	if !s.validated {
		panic("Using un-validated signers information")
	}
	if !s.powerAssigned {
		panic("Power has not been assigned in signers information")
	}
	if int(member.Index) >= s.committeeSize {
		panic("trying to increment signer information of non-existent committee member")
	}

	s.addMember(member, 1)
}

// Merges `other` with `s`. It assumes that they are mergeable. The caller is responsible to check
// so that coefficient overflow does not happen and the merge increases the total power in `s`.
func (s *VoteSigners) Merge(other *VoteSigners) {
	if err := safetyCheck(s.SignersBase, other.SignersBase); err != nil {
		panic(err.Error())
	}
	if !s.powerAssigned || !other.powerAssigned {
		panic("Power has not been assigned in signers information")
	}

	otherCount := 0
Loop:
	for i := 0; i < other.committeeSize; i++ {
		value := other.Bits.Get(i)
		switch value {
		case noSignature:
			continue Loop
		case oneSignature:
			s.increment(i)
		case twoSignatures:
			s.increment(i)
			s.increment(i)
		case multipleSignatures:
			// update s without using increment to save CPU power
			previousValue := s.Bits.Get(i)
			innerCount := 0
			switch previousValue {
			case noSignature:
				// we are adding a new signer, update the length cache
				s.length++
				fallthrough
			case oneSignature:
				fallthrough
			case twoSignatures:
				// add a new uint16 into the Coefficients array
				for j := 0; j < i; j++ {
					if s.Bits.Get(j) == multipleSignatures {
						innerCount++
					}
				}
				if innerCount == len(s.Coefficients) {
					s.Coefficients = append(s.Coefficients, uint16(previousValue))
				} else {
					s.Coefficients = append(s.Coefficients[:innerCount+1], s.Coefficients[innerCount:]...)
					s.Coefficients[innerCount] = uint16(previousValue)
				}
			case multipleSignatures:
				for j := 0; j < i; j++ {
					if s.Bits.Get(j) == multipleSignatures {
						innerCount++
					}
				}
			}
			// max allowed coefficient for a single validator is committeeSize
			if int(s.Coefficients[innerCount])+int(other.Coefficients[otherCount]) > s.committeeSize {
				panic("Aggregate signature coefficients exceeds allowed boundaries")
			}
			s.Bits.Set(i, multipleSignatures)
			s.Coefficients[innerCount] += other.Coefficients[otherCount]
			otherCount++
		}
		// update powers
		_, alreadyPresent := s.powers[i]
		if !alreadyPresent {
			s.powers[i] = other.powers[i]
			s.power.Add(s.power, other.powers[i])
		}
	}
}

// returns aggregated power of all senders
func (s *SignersBase) Power() *big.Int {
	if !s.powerAssigned {
		panic("Power has not been assigned in signers information")
	}
	return s.power
}

func (s *SignersBase) AssignPower(powers map[int]*big.Int, power *big.Int) {
	s.powers = powers
	s.power = power
	s.powerAssigned = true
}

func (s *QuorumSigners) Copy() *QuorumSigners {
	return &QuorumSigners{
		SignersBase:  s.SignersBase.Copy(),
		Coefficients: append(s.Coefficients[:0:0], s.Coefficients...),
	}
}

func (s *VoteSigners) Copy() *VoteSigners {
	return &VoteSigners{
		SignersBase:  s.SignersBase.Copy(),
		Coefficients: append(s.Coefficients[:0:0], s.Coefficients...),
	}
}

func (s *SignersBase) Copy() *SignersBase {
	var powers map[int]*big.Int
	if s.powers != nil {
		powers = make(map[int]*big.Int, len(s.powers))
		for index, power := range s.powers {
			powers[index] = new(big.Int).Set(power)
		}
	}
	return &SignersBase{
		Bits:          append(s.Bits[:0:0], s.Bits...),
		committeeSize: s.committeeSize,
		length:        s.length,
		powers:        powers,
		power:         s.power,
		validated:     s.validated,
		powerAssigned: s.powerAssigned,
	}
}

// returns list of indexes of validators that signed
// e.g. for bitmap [0 1 2 1 0] will return [ 1 2 2 3 ]
// the index 2 is repeated because we need to aggregate two times his key
func (s *VoteSigners) Flatten() []int {
	if !s.validated {
		panic("Using un-validated signers information")
	}
	return s.flatten(s.committeeSize)
}

// it is responsibility of the caller to pass the correct committee size
func (s *VoteSigners) flatten(committeeSize int) []int {
	var indexes []int
	count := 0
	var indexes []int
Loop:
	for i := 0; i < committeeSize; i++ {
		value := s.Bits.Get(i)
		switch value {
		case noSignature:
			continue Loop
		case oneSignature:
			indexes = append(indexes, i)
		case twoSignatures:
			indexes = append(indexes, i)
			indexes = append(indexes, i)
		case multipleSignatures:
			for j := 0; j < int(s.Coefficients[count]); j++ {
				indexes = append(indexes, i)
			}
			count++
		}
	}
	return indexes
}

// ForEveryRepeatedSigner iterates over each signer and calls the callback function.
func (s *Signers) ForEveryRepeatedSigner(callback func(signerIndex int), committeeLen int) {
	count := 0
	for i := 0; i < committeeLen; i++ {
		value := s.Bits.Get(i)
		switch value {
		case noSignature:
			continue
		case oneSignature:
			callback(i)
		case twoSignatures:
			callback(i)
			callback(i)
		case multipleSignatures:
			for j := 0; j < int(s.Coefficients[count]); j++ {
				callback(i)
			}
			count++
		}
	}
}

// ForEachDistinctSigner iterates over each signer and calls the callback function.
func (s *Signers) ForEachDistinctSigner(callback func(signerIndex int), committeeLen int) {
	for i := 0; i < committeeLen; i++ {
		if s.Bits.Get(i) > noSignature {
			callback(i)
		}
	}
}

// same as before, but repeated indexes are returned only once
func (s *SignersBase) FlattenUniq() []int {
	if !s.validated {
		panic("Using un-validated signers information")
	}
	return s.flattenUniq(s.committeeSize)
}

// it is responsibility of the caller to pass the correct committee size
func (s *SignersBase) flattenUniq(committeeSize int) []int {
	var indexes []int
	for i := 0; i < committeeSize; i++ {
		if s.Bits.Get(i) > noSignature {
			indexes = append(indexes, i)
		}
	}
	return indexes
}

// returns number of distinct signers of the aggregate
func (s *SignersBase) Len() int {
	if !s.validated {
		panic("Using un-validated signers information")
	}
	return s.length
}

// returns whether an aggregate is:
//   - a simple aggregate (all coefficients are 0 or 1)
//   - a complex aggregate (at least one coefficient is > 1)
func (s *VoteSigners) IsComplex() bool {
	if !s.validated {
		panic("Using un-validated signers information")
	}
	for i := 0; i < s.committeeSize; i++ {
		if s.Bits.Get(i) > oneSignature {
			return true
		}
	}
	return false
}

func (s *QuorumSigners) String() string {
	return fmt.Sprintf("Bits: %08b, Coefficients: %v, power: %v, validated: %v, powerAssigned: %v", s.Bits, s.Coefficients, s.power, s.validated, s.powerAssigned)
}

func (s *VoteSigners) String() string {
	return fmt.Sprintf("Bits: %08b, Coefficients: %v, power: %v, validated: %v, powerAssigned: %v", s.Bits, s.Coefficients, s.power, s.validated, s.powerAssigned)
}

func (s *SignersBase) Powers() map[int]*big.Int {
	if !s.powerAssigned {
		panic("Power has not been assigned in signers information")
	}
	return s.powers
}

func (s *SignersBase) CommitteeSize() int {
	if !s.validated {
		panic("Using un-validated signers information")
	}
	return s.committeeSize
}
