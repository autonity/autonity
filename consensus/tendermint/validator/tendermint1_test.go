package validator

/*
import (
	"bytes"
	"fmt"
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus/tendermint"
	"math"
	"strings"
	"testing"
	"testing/quick"

	"github.com/stretchr/testify/assert"
)

func TestValidatorSetBasic(t *testing.T) {
	// empty or nil validator lists are allowed,
	// but attempting to IncrementProposerPriority on them will panic.
	vset := newDefaultSet([]common.Address{}, tendermint.Tendermint)
	assert.Panics(t, func() { vset.IncrementProposerPriority(1) })

	vset = newDefaultSet(nil, tendermint.Tendermint)
	assert.Panics(t, func() { vset.IncrementProposerPriority(1) })

	assert.EqualValues(t, vset, vset.Copy())
	someAddress := &common.Address{}
	someAddress.SetBytes([]byte("some val"))
	idx, val := vset.GetByAddress(*someAddress)
	assert.Equal(t, -1, idx)
	assert.Nil(t, val)
	val = vset.GetByIndex(-100)
	assert.Nil(t, val)
	val = vset.GetByIndex(0)
	assert.Nil(t, val)
	val = vset.GetByIndex(100)
	assert.Nil(t, val)
	assert.Zero(t, vset.Size())
	assert.Equal(t, int64(0), vset.TotalVotingPower())
	assert.Nil(t, vset.GetProposer())

	// add

	val = randValidator_(vset.TotalVotingPower())
	assert.True(t, vset.Add(val))
	assert.True(t, vset.HasAddress(val.Address))
	idx, val2 := vset.GetByAddress(val.Address)
	assert.Equal(t, 0, idx)
	assert.Equal(t, val, val2)
	addr, val2 = vset.GetByIndex(0)
	assert.Equal(t, []byte(val.Address), addr)
	assert.Equal(t, val, val2)
	assert.Equal(t, 1, vset.Size())
	assert.Equal(t, val.VotingPower, vset.TotalVotingPower())
	assert.Equal(t, val, vset.GetProposer())
	assert.NotNil(t, vset.Hash())
	assert.NotPanics(t, func() { vset.IncrementProposerPriority(1) })

	// update
	assert.False(t, vset.Update(randValidator_(vset.TotalVotingPower())))
	_, val = vset.GetByAddress(val.Address)
	val.VotingPower += 100
	proposerPriority := val.ProposerPriority
	// Mimic update from types.PB2TM.ValidatorUpdates which does not know about ProposerPriority
	// and hence defaults to 0.
	val.ProposerPriority = 0
	assert.True(t, vset.Update(val))
	_, val = vset.GetByAddress(val.Address)
	assert.Equal(t, proposerPriority, val.ProposerPriority)

	// remove
	val2, removed := vset.Remove(randValidator_(vset.TotalVotingPower()).Address)
	assert.Nil(t, val2)
	assert.False(t, removed)
	val2, removed = vset.Remove(val.Address)
	assert.Equal(t, val.Address, val2.Address)
	assert.True(t, removed)
}

func TestCopy(t *testing.T) {
	vset := randValidatorSet(10)
	vsetHash := vset.Hash()
	if len(vsetHash) == 0 {
		t.Fatalf("ValidatorSet had unexpected zero hash")
	}

	vsetCopy := vset.Copy()
	vsetCopyHash := vsetCopy.Hash()

	if !bytes.Equal(vsetHash, vsetCopyHash) {
		t.Fatalf("ValidatorSet copy had wrong hash. Orig: %X, Copy: %X", vsetHash, vsetCopyHash)
	}
}

// Test that IncrementProposerPriority requires positive times.
func TestIncrementProposerPriorityPositiveTimes(t *testing.T) {
	vset := NewValidatorSet([]*Validator{
		newValidator([]byte("foo"), 1000),
		newValidator([]byte("bar"), 300),
		newValidator([]byte("baz"), 330),
	})

	assert.Panics(t, func() { vset.IncrementProposerPriority(-1) })
	assert.Panics(t, func() { vset.IncrementProposerPriority(0) })
	vset.IncrementProposerPriority(1)
}

func BenchmarkValidatorSetCopy(b *testing.B) {
	b.StopTimer()
	vset := NewValidatorSet([]*Validator{})
	for i := 0; i < 1000; i++ {
		privKey := ed25519.GenPrivKey()
		pubKey := privKey.PubKey()
		val := NewValidator(pubKey, 0)
		if !vset.Add(val) {
			panic("Failed to add validator")
		}
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		vset.Copy()
	}
}

//-------------------------------------------------------------------

func TestProposerSelection1(t *testing.T) {
	vset := NewValidatorSet([]*Validator{
		newValidator([]byte("foo"), 1000),
		newValidator([]byte("bar"), 300),
		newValidator([]byte("baz"), 330),
	})
	var proposers []string
	for i := 0; i < 99; i++ {
		val := vset.GetProposer()
		proposers = append(proposers, string(val.Address))
		vset.IncrementProposerPriority(1)
	}
	expected := `foo baz foo bar foo foo baz foo bar foo foo baz foo foo bar foo baz foo foo bar foo foo baz foo bar foo foo baz foo bar foo foo baz foo foo bar foo baz foo foo bar foo baz foo foo bar foo baz foo foo bar foo baz foo foo foo baz bar foo foo foo baz foo bar foo foo baz foo bar foo foo baz foo bar foo foo baz foo bar foo foo baz foo foo bar foo baz foo foo bar foo baz foo foo bar foo baz foo foo`
	if expected != strings.Join(proposers, " ") {
		t.Errorf("Expected sequence of proposers was\n%v\nbut got \n%v", expected, strings.Join(proposers, " "))
	}
}

func TestProposerSelection2(t *testing.T) {
	addr0 := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	addr1 := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}
	addr2 := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2}

	// when all voting power is same, we go in order of addresses
	val0, val1, val2 := newValidator(addr0, 100), newValidator(addr1, 100), newValidator(addr2, 100)
	valList := []*Validator{val0, val1, val2}
	vals := NewValidatorSet(valList)
	for i := 0; i < len(valList)*5; i++ {
		ii := (i) % len(valList)
		prop := vals.GetProposer()
		if !bytes.Equal(prop.Address, valList[ii].Address) {
			t.Fatalf("(%d): Expected %X. Got %X", i, valList[ii].Address, prop.Address)
		}
		vals.IncrementProposerPriority(1)
	}

	// One validator has more than the others, but not enough to propose twice in a row
	*val2 = *newValidator(addr2, 400)
	vals = NewValidatorSet(valList)
	// vals.IncrementProposerPriority(1)
	prop := vals.GetProposer()
	if !bytes.Equal(prop.Address, addr2) {
		t.Fatalf("Expected address with highest voting power to be first proposer. Got %X", prop.Address)
	}
	vals.IncrementProposerPriority(1)
	prop = vals.GetProposer()
	if !bytes.Equal(prop.Address, addr0) {
		t.Fatalf("Expected smallest address to be validator. Got %X", prop.Address)
	}

	// One validator has more than the others, and enough to be proposer twice in a row
	*val2 = *newValidator(addr2, 401)
	vals = NewValidatorSet(valList)
	prop = vals.GetProposer()
	if !bytes.Equal(prop.Address, addr2) {
		t.Fatalf("Expected address with highest voting power to be first proposer. Got %X", prop.Address)
	}
	vals.IncrementProposerPriority(1)
	prop = vals.GetProposer()
	if !bytes.Equal(prop.Address, addr2) {
		t.Fatalf("Expected address with highest voting power to be second proposer. Got %X", prop.Address)
	}
	vals.IncrementProposerPriority(1)
	prop = vals.GetProposer()
	if !bytes.Equal(prop.Address, addr0) {
		t.Fatalf("Expected smallest address to be validator. Got %X", prop.Address)
	}

	// each validator should be the proposer a proportional number of times
	val0, val1, val2 = newValidator(addr0, 4), newValidator(addr1, 5), newValidator(addr2, 3)
	valList = []*Validator{val0, val1, val2}
	propCount := make([]int, 3)
	vals = NewValidatorSet(valList)
	N := 1
	for i := 0; i < 120*N; i++ {
		prop := vals.GetProposer()
		ii := prop.Address[19]
		propCount[ii]++
		vals.IncrementProposerPriority(1)
	}

	if propCount[0] != 40*N {
		t.Fatalf("Expected prop count for validator with 4/12 of voting power to be %d/%d. Got %d/%d", 40*N, 120*N, propCount[0], 120*N)
	}
	if propCount[1] != 50*N {
		t.Fatalf("Expected prop count for validator with 5/12 of voting power to be %d/%d. Got %d/%d", 50*N, 120*N, propCount[1], 120*N)
	}
	if propCount[2] != 30*N {
		t.Fatalf("Expected prop count for validator with 3/12 of voting power to be %d/%d. Got %d/%d", 30*N, 120*N, propCount[2], 120*N)
	}
}

func TestProposerSelection3(t *testing.T) {
	vset := NewValidatorSet([]*Validator{
		newValidator([]byte("a"), 1),
		newValidator([]byte("b"), 1),
		newValidator([]byte("c"), 1),
		newValidator([]byte("d"), 1),
	})

	proposerOrder := make([]*Validator, 4)
	for i := 0; i < 4; i++ {
		proposerOrder[i] = vset.GetProposer()
		vset.IncrementProposerPriority(1)
	}

	// i for the loop
	// j for the times
	// we should go in order for ever, despite some IncrementProposerPriority with times > 1
	var i, j int
	for ; i < 10000; i++ {
		got := vset.GetProposer().Address
		expected := proposerOrder[j%4].Address
		if !bytes.Equal(got, expected) {
			t.Fatalf(fmt.Sprintf("vset.Proposer (%X) does not match expected proposer (%X) for (%d, %d)", got, expected, i, j))
		}

		// serialize, deserialize, check proposer
		b := vset.toBytes()
		vset.fromBytes(b)

		computed := vset.GetProposer() // findGetProposer()
		if i != 0 {
			if !bytes.Equal(got, computed.Address) {
				t.Fatalf(fmt.Sprintf("vset.Proposer (%X) does not match computed proposer (%X) for (%d, %d)", got, computed.Address, i, j))
			}
		}

		// times is usually 1
		times := 1
		mod := (cmn.RandInt() % 5) + 1
		if cmn.RandInt()%mod > 0 {
			// sometimes its up to 5
			times = (cmn.RandInt() % 4) + 1
		}
		vset.IncrementProposerPriority(times)

		j += times
	}
}

func newValidator(address []byte, power int64) *Validator {
	return &Validator{Address: address, VotingPower: power}
}

func randPubKey() crypto.PubKey {
	var pubKey [32]byte
	copy(pubKey[:], cmn.RandBytes(32))
	return ed25519.PubKeyEd25519(pubKey)
}

func randValidator_(totalVotingPower int64) *Validator {
	// this modulo limits the ProposerPriority/VotingPower to stay in the
	// bounds of MaxTotalVotingPower minus the already existing voting power:
	val := NewValidator(randPubKey(), cmn.RandInt64()%(MaxTotalVotingPower-totalVotingPower))
	val.ProposerPriority = cmn.RandInt64() % (MaxTotalVotingPower - totalVotingPower)
	return val
}

func randValidatorSet(numValidators int) *ValidatorSet {
	validators := make([]*Validator, numValidators)
	totalVotingPower := int64(0)
	for i := 0; i < numValidators; i++ {
		validators[i] = randValidator_(totalVotingPower)
		totalVotingPower += validators[i].VotingPower
	}
	return NewValidatorSet(validators)
}

func (valSet *ValidatorSet) toBytes() []byte {
	bz, err := cdc.MarshalBinaryLengthPrefixed(valSet)
	if err != nil {
		panic(err)
	}
	return bz
}

func (valSet *ValidatorSet) fromBytes(b []byte) {
	err := cdc.UnmarshalBinaryLengthPrefixed(b, &valSet)
	if err != nil {
		// DATA HAS BEEN CORRUPTED OR THE SPEC HAS CHANGED
		panic(err)
	}
}

//-------------------------------------------------------------------

*/