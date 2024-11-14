package types

import (
	"errors"
	"fmt"
	"math/big"
)

var (
	bitsPerValidator  = 2 //NOTE: if this gets changed, major refactoring will be needed for this file. Proceed with caution.
	bitsInByte        = 8
	validatorsPerByte = bitsInByte / bitsPerValidator
	// 11000000 00110000 00001100 00000011
	getMasks = []byte{0xC0, 0x30, 0x0C, 0x03}
	// 00111111 11001111 11110011 11111100
	setMasks  = []byte{0x3F, 0xCF, 0xF3, 0xFC}
	maxUint16 = (1 << 16) - 1

	// possible values of 2 bits
	noSignature        = byte(0)
	oneSignature       = byte(1)
	twoSignatures      = byte(2)
	multipleSignatures = byte(3)

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

type Signers struct {
	Bits         validatorBitmap
	Coefficients []uint16 // support up to 65535 committee members

	// these fields are not serialized, but instead computed at preValidate steps
	committeeSize int              `rlp:"-"`
	length        int              `rlp:"-"` // number of distinct signers
	powers        map[int]*big.Int `rlp:"-"`
	power         *big.Int         `rlp:"-"` // aggregated power of all senders

	// auxiliary data structures flags
	validated     bool `rlp:"-"` // if true --> Bits and Coefficients have correct length + committeeSize and length is assigned
	powerAssigned bool `rlp:"-"` // if true --> powers and power assigned
}

func NewSigners(committeeSize int) *Signers {
	if committeeSize > maxUint16 {
		panic("Unsupported committee size")
	}
	return &Signers{
		Bits:          NewValidatorBitmap(committeeSize),
		Coefficients:  make([]uint16, 0),
		committeeSize: committeeSize,

		length:        0,
		powers:        make(map[int]*big.Int),
		power:         new(big.Int),
		validated:     true,
		powerAssigned: true, // when we are locally creating a sender info, we are ok with power being 0 initially
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

func (vb validatorBitmap) Get(validatorIndex int) byte {
	byteIndex := validatorIndex / validatorsPerByte
	bitIndex := validatorIndex % validatorsPerByte
	b := vb[byteIndex]

	result := b & getMasks[bitIndex]
	shift := validatorsPerByte - 1 - bitIndex
	result = result >> (shift * bitsPerValidator)
	return result
}

// NOTE: be careful when calling directly this function without passing through the `increment` function.
// this function will not invalidate any cache, it is just a naive setter
func (vb validatorBitmap) Set(validatorIndex int, value byte) {
	if value > multipleSignatures {
		panic("Trying to set value that cannot fit into 2 bits")
	}
	bitIndex := validatorIndex % validatorsPerByte
	shift := validatorsPerByte - 1 - bitIndex
	valueShifted := value << (shift * bitsPerValidator)

	byteIndex := validatorIndex / validatorsPerByte
	vb[byteIndex] = vb[byteIndex] & setMasks[bitIndex]
	vb[byteIndex] = vb[byteIndex] | valueShifted
}

// Correct full validation cannot be done until we know the committee size of this block,
// but we can already ensure that the fields have sane values
func (s *Signers) SanityCheck() error {
	if len(s.Bits) == 0 {
		return fmt.Errorf("validator bitmap is empty")
	}
	// 1 byte --> 4 validators
	if len(s.Bits) > maxAllowedSigners/validatorsPerByte { // max 4kb of data
		return fmt.Errorf("invalid Bits length: %d", len(s.Bits))
	}
	if len(s.Coefficients) > maxAllowedSigners { // max 32kb of data
		return fmt.Errorf("invalid Coefficient length: %d", len(s.Coefficients))
	}
	return nil
}

// validates the sender info, used to ensure received aggregates have correctly sized buffers
func (s *Signers) Validate(committeeSize int) error {
	distinctSigners, err := s.validate(committeeSize)
	if err != nil {
		return err
	}

	s.committeeSize = committeeSize
	s.length = distinctSigners
	s.validated = true
	return nil
}

// validates the signer information and returns the number of distinct signers
// it does not mutate the signers state
func (s *Signers) validate(committeeSize int) (int, error) {
	// whether locally created or received from wire, Bits and Coefficients are never nil
	if s.Bits == nil || s.Coefficients == nil {
		return 0, ErrNilSigners
	}

	// length safety check
	if !s.Bits.Valid(committeeSize) || len(s.Coefficients) > committeeSize {
		return 0, ErrWrongSizeSigners
	}

	// gather data about signers bits
	countNonZero := 0
	countLong := 0
	sum := 0
	for i := 0; i < committeeSize; i++ {
		value := s.Bits.Get(i)
		if value > noSignature { // 01 10 11
			countNonZero++
		}
		if value == multipleSignatures { // 11
			countLong++
		}
		sum += int(value)
	}

	// there has to be at least a signer
	if sum == 0 {
		return 0, ErrEmptySigners
	}

	// len(s.Coefficients) should be the same as the number of elements with value 11 in s.Bits
	if len(s.Coefficients) != countLong {
		return 0, ErrWrongCoefficientLen
	}

	// if individual signature, its coefficient should be one (01)
	if countNonZero == 1 && sum != 1 {
		return 0, ErrInvalidSingleSig
	}

	// check that all coefficients respect the maximum allowed boundary (committeeSize)
	for _, coefficient := range s.Coefficients {
		if int(coefficient) > committeeSize {
			return 0, ErrInvalidCoefficient
		}
	}
	return countNonZero, nil
}

func (s *Signers) Contains(index int) bool {
	if !s.validated {
		panic("Trying to use not validated signer information")
	}
	if index >= s.committeeSize {
		panic("trying to call contains on non-existent committee member")
	}
	return s.Bits.Get(index) > noSignature
}

func safetyCheck(first *Signers, second *Signers) error {
	if !first.validated || !second.validated {
		return ErrNotValidated
	}
	if first.committeeSize != second.committeeSize {
		return ErrDifferentSize
	}
	return nil
}

// checks that the resulting aggregate still respects the `committeeSize` boundary
func (s *Signers) RespectsBoundaries(other *Signers) bool {
	if err := safetyCheck(s, other); err != nil {
		panic(err.Error())
	}

	var firstCoefficient int
	var count int
	var secondCoefficient int
	var count2 int

	for i := 0; i < s.committeeSize; i++ {
		firstCoefficient = int(s.Bits.Get(i))
		if firstCoefficient == int(multipleSignatures) {
			firstCoefficient = int(s.Coefficients[count])
			count++
		}

		secondCoefficient = int(other.Bits.Get(i))
		if secondCoefficient == int(multipleSignatures) {
			secondCoefficient = int(other.Coefficients[count2])
			count2++
		}

		if firstCoefficient+secondCoefficient > s.committeeSize {
			return false
		}
	}
	return true
}

// TODO: maybe this can done more efficiently using bitwise operations.
// however it is not trivial since we use two bits per validators
func (s *Signers) AddsInformation(other *Signers) bool {
	if err := safetyCheck(s, other); err != nil {
		panic(err.Error())
	}

	for i := 0; i < s.committeeSize; i++ {
		if s.Bits.Get(i) == noSignature && other.Bits.Get(i) != noSignature {
			return true
		}
	}
	return false
}

// TODO: maybe this can done more efficiently using bitwise operations.
// however it is not trivial since we use two bits per validators
func (s *Signers) CanMergeSimple(other *Signers) bool {
	if err := safetyCheck(s, other); err != nil {
		panic(err.Error())
	}

	for i := 0; i < s.committeeSize; i++ {
		if s.Bits.Get(i)+other.Bits.Get(i) > oneSignature {
			return false
		}
	}
	return true
}

func (s *Signers) increment(index int) {
	if !s.validated {
		panic("Using un-validated signers information")
	}
	previousValue := s.Bits.Get(index)
	var value byte
	switch previousValue {
	case noSignature:
		value = oneSignature // 01
		// we are adding a new signer, update the length cache
		s.length++
	case oneSignature:
		value = twoSignatures // 10
	case twoSignatures:
		value = multipleSignatures // 11
		// add a new uint16 into the Coefficients array
		count := 0
		for i := 0; i < index; i++ {
			if s.Bits.Get(i) == multipleSignatures {
				count++
			}
		}
		if count == len(s.Coefficients) {
			s.Coefficients = append(s.Coefficients, uint16(multipleSignatures))
		} else {
			s.Coefficients = append(s.Coefficients[:count+1], s.Coefficients[count:]...)
			s.Coefficients[count] = uint16(multipleSignatures)
		}
	case multipleSignatures:
		value = multipleSignatures // 11
		// update uint16 into the Coefficients array
		count := 0
		for i := 0; i < index; i++ {
			if s.Bits.Get(i) == multipleSignatures {
				count++
			}
		}
		// max allowed coefficient for a single validator is committeeSize
		if int(s.Coefficients[count]) >= s.committeeSize {
			panic("Aggregate signature coefficients exceeds allowed boundaries")
		}
		s.Coefficients[count]++
	}
	s.Bits.Set(index, value)
}

func (s *Signers) Increment(member *CommitteeMember) {
	if !s.validated {
		panic("Using un-validated signers information")
	}
	if !s.powerAssigned {
		panic("Power has not been assigned in signers information")
	}

	index := int(member.Index)
	if index >= s.committeeSize {
		panic("trying to increment signer information of non-existent committee member")
	}
	s.increment(index)

	_, alreadyPresent := s.powers[index]
	if !alreadyPresent {
		s.powers[index] = member.VotingPower
		s.power.Add(s.power, member.VotingPower)
	}
}

func (s *Signers) Merge(other *Signers) {
	if err := safetyCheck(s, other); err != nil {
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
func (s *Signers) Power() *big.Int {
	if !s.powerAssigned {
		panic("Power has not been assigned in signers information")
	}
	return s.power
}

func (s *Signers) AssignPower(powers map[int]*big.Int, power *big.Int) {
	s.powers = powers
	s.power = power
	s.powerAssigned = true
}

func (s *Signers) Copy() *Signers {
	var powers map[int]*big.Int
	if s.powers != nil {
		powers = make(map[int]*big.Int, len(s.powers))
		for index, power := range s.powers {
			powers[index] = new(big.Int).Set(power)
		}
	}
	return &Signers{
		Bits:          append(s.Bits[:0:0], s.Bits...),
		Coefficients:  append(s.Coefficients[:0:0], s.Coefficients...),
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
func (s *Signers) Flatten() []int {
	if !s.validated {
		panic("Using un-validated signers information")
	}
	return s.flatten(s.committeeSize)
}

// it is responsibility of the caller to pass the correct committee size
func (s *Signers) flatten(committeeSize int) []int {
	var indexes []int
	count := 0
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

// same as before, but repeated indexes are returned only once
func (s *Signers) FlattenUniq() []int {
	if !s.validated {
		panic("Using un-validated signers information")
	}
	return s.flattenUniq(s.committeeSize)
}

// it is responsibility of the caller to pass the correct committee size
func (s *Signers) flattenUniq(committeeSize int) []int {
	var indexes []int
	for i := 0; i < committeeSize; i++ {
		if s.Bits.Get(i) > noSignature {
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

// returns whether an aggregate is:
//   - a simple aggregate (all coefficients are 0 or 1)
//   - a complex aggregate (at least one coefficient is > 1)
func (s *Signers) IsComplex() bool {
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

func (s *Signers) String() string {
	return fmt.Sprintf("Bits: %08b, Coefficients: %v, power: %v, validated: %v, powerAssigned: %v", s.Bits, s.Coefficients, s.power, s.validated, s.powerAssigned)
}

func (s *Signers) Powers() map[int]*big.Int {
	if !s.powerAssigned {
		panic("Power has not been assigned in signers information")
	}
	return s.powers
}

func (s *Signers) CommitteeSize() int {
	if !s.validated {
		panic("Using un-validated signers information")
	}
	return s.committeeSize
}
