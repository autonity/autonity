package types

import (
	"fmt"
	"math"
	"math/big"
	"sort"

	"github.com/autonity/autonity/common"
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

type SendersInfo struct {
	Bits         validatorBitmap
	Coefficients []uint16 // support up to 65535 committee members

	// these fields are not serialized, but instead computed at preValidate steps
	committeeSize int                    `rlp:"-"`
	addresses     map[int]common.Address `rlp:"-"` // TODO(lorenzo) turn these into arrays?
	powers        map[int]*big.Int       `rlp:"-"`
}

func NewSendersInfo(committeeSize int) *SendersInfo {
	if committeeSize > maxUint16 {
		panic("Unsupported committee size")
	}
	return &SendersInfo{
		Bits:          NewValidatorBitmap(committeeSize),
		Coefficients:  make([]uint16, 0),
		committeeSize: committeeSize,
		addresses:     make(map[int]common.Address),
		powers:        make(map[int]*big.Int),
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

// TODO(lorenzo) return error instead of bool
// validates the sender info, used to ensure received aggregates have correctly sized buffers
func (s *SendersInfo) Valid(committeeSize int) bool {
	// whether locally created or received from wire, Bits and Coefficients are never nil
	if s.Bits == nil || s.Coefficients == nil {
		return false
	}

	// length safety check
	if !s.Bits.Valid(committeeSize) || len(s.Coefficients) > committeeSize {
		return false
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
		return false
	}

	// len(s.Coefficients) should be the same as the number of elements with value 11 in s.Bits
	if len(s.Coefficients) != countLong {
		return false
	}

	// if individual signature, its coefficient should be one (01)
	if countNonZero == 1 && sum != 1 {
		return false
	}

	// shortcircuit validation here if single sig with coefficient one
	if countNonZero == 1 {
		return true
	}

	//TODO(lorenzo) do we still need this? I think it might still be correct, but is probably overkill
	// to be reviewed after I change the aggregation

	/* The following rules are enforced when aggregating votes:
	 	*   1. votes are ordered by decreasing number of distinct senders before aggregation
	  *   2. a vote is aggregated to the previous ones only if it adds useful information (i.e. adds a new sender)
		*
		* 	Example: AggregateVotes(AB,ABC,D) = ABCD. AB is discarded because it doesn't add any meaningful information and it just bloats the aggregate
		* 					 See the AggregateVotes function in message.go for more details
		*
		*  These aggregation rules cause that all aggregates needs to respect the following polynomial form:
		*    Assuming N distinct senders, let's define N' = N - 2 for convenience
		*    aggregate: A * 2^N' + B * 2^(N' - 1) + C * 2^(N' - 2) + ... + X * 2^(N' - N') + Y
		*
		*    NOTE:
		*      - This polynomial form defines the maximum coefficient that senders in an aggregate can have. However it is possible for the senders to have lower coefficients. See examples below.
		*      - there are two elements with coefficient 1 at the end (X and Y)
		*
		*  Examples:
		* 		- A + B + C + D   ---> VALID
		*     - 255A + B        ---> INVALID
		*     - 4A + 2B + C + D ---> VALID
		*     - 4A + 3B + C + D ---> INVALID
		*     - 5A + 2B + C + D ---> INVALID (A coefficient too high)
		*     - 4A + 2B + C     ---> INVALID (missing term)
		*
		*  Therefore we can infer the following validation rules. Each coefficient needs to be:
		* 	1. less or equal to 2^(N' - i)
		*   2. less or equal to the sum of all the remaining lower order coefficients
	*/

	// compute coefficients
	var coefficients []int
	count := 0
	for i := 0; i < committeeSize; i++ {
		value := s.Bits.Get(i)
		if value == 0 {
			continue // ignore 0 coefficients
		}
		if value == 1 || value == 2 {
			coefficients = append(coefficients, int(value))
		}
		if value == 3 { // look into s.Coefficients
			coefficients = append(coefficients, int(s.Coefficients[count]))
			count++
		}
	}

	// sort in descending order
	sort.Slice(coefficients, func(i, j int) bool {
		return coefficients[i] > coefficients[j]
	})

	N := len(coefficients) // N >= 2 here
	Nprime := N - 2        // Nprime >= 0
	for i, coefficient := range coefficients {
		if i < N-1 {
			// we are not dealing with the last element of the array

			// check rule 1
			limit := int(math.Pow(2, float64(Nprime-i))) //TODO(lorenzo) what if doesn't fit
			if coefficient > limit {
				return false
			}

			// check rule 2
			sum := 0
			for j := i + 1; j < N; j++ {
				sum += coefficients[j]
			}
			if coefficient > sum {
				return false
			}
		} else {
			// last element, just check that coefficient == 1
			if coefficient != 1 {
				return false
			}
		}
	}

	return true
}

// TODO(lorenzo) refinements, maybe I can do this more efficiently using bitwise operations
// however it is not trivial since we use two bits per validators
func (s *SendersInfo) AddsInformation(other *SendersInfo) bool {
	//TODO(lorenzo) refinements, add check on length of the two senderinfo?

	for i := 0; i < s.committeeSize; i++ {
		if s.Bits.Get(i) == 0 && other.Bits.Get(i) != 0 {
			return true
		}
	}
	return false
}

// TODO(lorenzo) refinements, maybe I can do this more efficiently using bitwise operations
// however it is not trivial since we use two bits per validators
func (s *SendersInfo) CanMergeSimple(other *SendersInfo) bool {
	//TODO(lorenzo) refinements, add check on length of the two senderinfo?

	for i := 0; i < s.committeeSize; i++ {
		if s.Bits.Get(i)+other.Bits.Get(i) > 1 {
			return false
		}
	}
	return true
}

func (s *SendersInfo) increment(index int) {
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
		// max allowed coefficient for a single validator is committeeSize - 1
		//TODO(lorenzo) write test that ends up here, to verify that this actually never happens when aggregating votes
		if int(s.Coefficients[count]) >= (s.committeeSize - 1) {
			panic("Aggregate signature coefficients exceeds allowed boundaries")
		}
		s.Coefficients[count]++
	}
	s.Bits.Set(index, value)
}

func (s *SendersInfo) Increment(member *CommitteeMember) {
	index := int(member.Index)
	if index >= s.committeeSize {
		panic("trying to increment sender information of non-existant committee member")
	}
	s.increment(index)

	s.addresses[index] = member.Address
	s.powers[index] = new(big.Int).Set(member.VotingPower)
}

func (s *SendersInfo) Merge(other *SendersInfo) {
	//TODO(lorenzo) refinements, maybe we can skip this since the Valid() check at preValidate()
	if len(s.Bits) != len(other.Bits) || s.committeeSize != other.committeeSize {
		// should always merge for the same height --> same committee size --> same legnth
		panic("trying to merge sender information with different length")
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
		s.addresses[i] = other.addresses[i]
		s.powers[i] = new(big.Int).Set(other.powers[i])
	}
}

// returns aggregated power of all senders
func (s *SendersInfo) Power() *big.Int {
	power := new(big.Int)
	for i := 0; i < s.committeeSize; i++ {
		// regardless of the value, we sum power only once (to prevent counting duplicated messages power twice)
		if s.Bits.Get(i) > 0 {
			power.Add(power, s.powers[i])
		}
	}
	return power
}

func (s *SendersInfo) Copy() *SendersInfo {
	addresses := make(map[int]common.Address)
	for index, address := range s.addresses {
		addresses[index] = address
	}
	powers := make(map[int]*big.Int)
	for index, power := range s.powers {
		powers[index] = new(big.Int).Set(power)
	}
	return &SendersInfo{
		Bits:          append(s.Bits[:0:0], s.Bits...),
		Coefficients:  append(s.Coefficients[:0:0], s.Coefficients...),
		committeeSize: s.committeeSize,
		addresses:     addresses,
		powers:        powers,
	}
}

// returns list of indexes of validators that signed
// e.g. for bitmap [0 1 2 1 0] will return [ 1 2 2 3 ]
// the index 2 is repeated because we need to aggregate two times his key
func (s *SendersInfo) Flatten() []int {
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
func (s *SendersInfo) FlattenUniq() []int {
	var indexes []int
	for i := 0; i < s.committeeSize; i++ {
		if s.Bits.Get(i) > 0 {
			indexes = append(indexes, i)
		}
	}
	return indexes
}

// returns number of distinct signers of the aggregate
func (s *SendersInfo) Len() int {
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
func (s *SendersInfo) IsComplex() bool {
	for i := 0; i < s.committeeSize; i++ {
		if s.Bits.Get(i) > 1 {
			return true
		}
	}
	return false
}

func (s *SendersInfo) String() string {
	return fmt.Sprintf("Bits: %08b, Coefficients: %v", s.Bits, s.Coefficients)
}

func (s *SendersInfo) Addresses() map[int]common.Address {
	return s.addresses
}
func (s *SendersInfo) SetAddresses(addresses map[int]common.Address) {
	s.addresses = addresses
}

func (s *SendersInfo) Powers() map[int]*big.Int {
	return s.powers
}

func (s *SendersInfo) SetPowers(powers map[int]*big.Int) {
	s.powers = powers
}

// Len returns how many committee members the sendersInfo contains
func (s *SendersInfo) CommitteeSize() int {
	return s.committeeSize
}

func (s *SendersInfo) SetCommitteeSize(committeeSize int) {
	s.committeeSize = committeeSize
}
