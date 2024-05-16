package types

import (
	"errors"
	"fmt"
	"math"
	"math/big"
)

var (
	bitsPerValidator  = 2 //NOTE: if this gets changed, major refactoring will be needed for this file. Proceed with caution.
	bitsInByte        = 8
	validatorsPerByte = bitsInByte / bitsPerValidator
	maxValue          = byte(3) // max value that can be fit in 2 bits
	// 11000000 00110000 00001100 00000011
	getMasks = []byte{0xC0, 0x30, 0x0C, 0x03}
	// 00111111 11001111 11110011 11111100
	setMasks  = []byte{0x3F, 0xCF, 0xF3, 0xFC}
	maxUint16 = (1 << 16) - 1

	ErrNilSigners          = errors.New("validator bitmap or coefficient array is nil")
	ErrOversizedSigners    = errors.New("validator bitmap or coefficient array is oversized")
	ErrEmptySigners        = errors.New("signers information is empty")
	ErrWrongCoefficientLen = errors.New("coefficient array has incorrect length")
	ErrInvalidSingleSig    = errors.New("individual signature has coefficient != 1")
	ErrInvalidCoefficient  = errors.New("coefficient exceeds maximum boundary (committee size)")

	ErrNotValidated  = errors.New("Using un-validated signers information")
	ErrDifferentSize = errors.New("Comparing signers information with different committee size")
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
	powers        map[int]*big.Int `rlp:"-"` // TODO(lorenzo) performance, turn these into arrays? we will have sparse arrays but maybe it is fine. Other option is using functions
	power         *big.Int         `rlp:"-"` // aggregated power of all senders

	// auxiliary data structures flags
	validated     bool `rlp:"-"` // if true --> Bits and Coefficients have correct length + committeeSize is assigned
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
		powers:        make(map[int]*big.Int),
		power:         new(big.Int),
		validated:     true,
		powerAssigned: true, // when we are locally creating a sender info, we are ok with power being 0 initially
	}
}

type validatorBitmap []byte

func NewValidatorBitmap(committeeSize int) validatorBitmap {
	bitLength := committeeSize * bitsPerValidator
	byteLength := int(math.Ceil(float64(bitLength) / float64(bitsInByte)))
	return make(validatorBitmap, byteLength)
}

// ensures that the validator bitmap has the correct length compared to the committee size
// used to validate aggregate messages coming from other peers
func (vb validatorBitmap) Valid(committeeSize int) bool {
	expectedBitLength := committeeSize * bitsPerValidator
	expectedByteLength := int(math.Ceil(float64(expectedBitLength) / float64(bitsInByte)))
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

func (vb validatorBitmap) Set(validatorIndex int, value byte) {
	if value > maxValue {
		panic("Trying to set value that cannot fit into 2 bits")
	}
	bitIndex := validatorIndex % validatorsPerByte
	shift := validatorsPerByte - 1 - bitIndex
	valueShifted := value << (shift * bitsPerValidator)

	byteIndex := validatorIndex / validatorsPerByte
	vb[byteIndex] = vb[byteIndex] & setMasks[bitIndex]
	vb[byteIndex] = vb[byteIndex] | valueShifted
}

// validates the sender info, used to ensure received aggregates have correctly sized buffers
func (s *Signers) Valid(committeeSize int) error {
	// whether locally created or received from wire, Bits and Coefficients are never nil
	if s.Bits == nil || s.Coefficients == nil {
		return ErrNilSigners
	}

	// length safety check
	if !s.Bits.Valid(committeeSize) || len(s.Coefficients) > committeeSize {
		return ErrOversizedSigners
	}

	// gather data about senders bits
	countNonZero := 0
	countLong := 0
	sum := 0
	for i := 0; i < committeeSize; i++ {
		value := s.Bits.Get(i)
		if value > 0 { // 01 10 11
			countNonZero++
		}
		if value == 3 { // 11
			countLong++
		}
		sum += int(value)
	}

	// there has to be at least a sender
	if sum == 0 {
		return ErrEmptySigners
	}

	// len(s.Coefficients) should be the same as the number of elements with value 11 in s.Bits
	if len(s.Coefficients) != countLong {
		return ErrWrongCoefficientLen
	}

	// if individual signature, its coefficient should be one (01)
	if countNonZero == 1 && sum != 1 {
		return ErrInvalidSingleSig
	}

	// check that all coefficients respect the maximum allowed boundary (committeeSize)
	for _, coefficient := range s.Coefficients {
		if int(coefficient) > committeeSize {
			return ErrInvalidCoefficient
		}
	}

	s.committeeSize = committeeSize
	s.validated = true
	return nil
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
		if firstCoefficient == 3 {
			firstCoefficient = int(s.Coefficients[count])
			count++
		}

		secondCoefficient = int(other.Bits.Get(i))
		if secondCoefficient == 3 {
			secondCoefficient = int(other.Coefficients[count])
			count2++
		}

		if firstCoefficient+secondCoefficient > s.committeeSize {
			return false
		}
	}
	return true
}

// TODO(lorenzo) refinements, maybe I can do this more efficiently using bitwise operations
// however it is not trivial since we use two bits per validators
func (s *Signers) AddsInformation(other *Signers) bool {
	if err := safetyCheck(s, other); err != nil {
		panic(err.Error())
	}

	for i := 0; i < s.committeeSize; i++ {
		if s.Bits.Get(i) == 0 && other.Bits.Get(i) != 0 {
			return true
		}
	}
	return false
}

// TODO(lorenzo) refinements, maybe I can do this more efficiently using bitwise operations
// however it is not trivial since we use two bits per validators
func (s *Signers) CanMergeSimple(other *Signers) bool {
	if err := safetyCheck(s, other); err != nil {
		panic(err.Error())
	}

	for i := 0; i < s.committeeSize; i++ {
		if s.Bits.Get(i)+other.Bits.Get(i) > 1 {
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
	case 0:
		value = 1 // 01
	case 1:
		value = 2 // 10
	case 2:
		value = 3 // 11
		// add a new uint16 into the Coefficients array
		count := 0
		for i := 0; i < index; i++ {
			if s.Bits.Get(i) == byte(3) {
				count++
			}
		}
		if count == len(s.Coefficients) {
			s.Coefficients = append(s.Coefficients, uint16(3))
		} else {
			s.Coefficients = append(s.Coefficients[:count+1], s.Coefficients[count:]...)
			s.Coefficients[count] = uint16(3)
		}
	case 3:
		value = 3 // 11
		// update uint16 into the Coefficients array
		count := 0
		for i := 0; i < index; i++ {
			if s.Bits.Get(i) == 3 {
				count++
			}
		}
		// max allowed coefficient for a single validator is committeeSize
		//TODO(lorenzo) write test that ends up here, to verify that this actually never happens when aggregating votes
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
		panic("trying to increment signer information of non-existant committee member")
	}
	s.increment(index)

	//TODO(lorenzo) can senders info be updated concurrently?

	_, alreadyPresent := s.powers[index]
	if !alreadyPresent {
		s.powers[index] = new(big.Int).Set(member.VotingPower)
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

	count := 0
Loop:
	for i := 0; i < other.committeeSize; i++ {
		value := other.Bits.Get(i)
		switch value {
		case 0:
			continue Loop
		case 1:
			s.increment(i)
		case 2:
			s.increment(i)
			s.increment(i)
		case 3:
			//TODO(lorenzo) refinements, instaed of looping just sum the other uint16
			for j := 0; j < int(other.Coefficients[count]); j++ {
				s.increment(i)
			}
			count++
		}
		_, alreadyPresent := s.powers[i]
		if !alreadyPresent {
			s.powers[i] = new(big.Int).Set(other.powers[i])
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
	powers := make(map[int]*big.Int)
	for index, power := range s.powers {
		powers[index] = new(big.Int).Set(power)
	}
	return &Signers{
		Bits:          append(s.Bits[:0:0], s.Bits...),
		Coefficients:  append(s.Coefficients[:0:0], s.Coefficients...),
		committeeSize: s.committeeSize,
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
	var indexes []int
	count := 0
Loop:
	for i := 0; i < s.committeeSize; i++ {
		value := s.Bits.Get(i)
		switch value {
		case 0:
			continue Loop
		case 1:
			indexes = append(indexes, i)
		case 2:
			indexes = append(indexes, i)
			indexes = append(indexes, i)
		case 3:
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
	var indexes []int
	for i := 0; i < s.committeeSize; i++ {
		if s.Bits.Get(i) > 0 {
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
	count := 0
	for i := 0; i < s.committeeSize; i++ {
		if s.Bits.Get(i) > 0 {
			count++
		}
	}
	return count
}

// returns whether an aggregate is:
//   - a simple aggregate (all coefficients are 0 or 1)
//   - a complex aggregate (at least one coefficient is > 1)
func (s *Signers) IsComplex() bool {
	if !s.validated {
		panic("Using un-validated signers information")
	}
	for i := 0; i < s.committeeSize; i++ {
		if s.Bits.Get(i) > 1 {
			return true
		}
	}
	return false
}

func (s *Signers) String() string {
	return fmt.Sprintf("Bits: %08b, Coefficients: %v", s.Bits, s.Coefficients)
}

func (s *Signers) Powers() map[int]*big.Int {
	if !s.powerAssigned {
		panic("Power has not been assigned in signers information")
	}
	return s.powers
}

// Len returns how many committee members the Signers contains
func (s *Signers) CommitteeSize() int {
	if !s.validated {
		panic("Using un-validated signers information")
	}
	return s.committeeSize
}
