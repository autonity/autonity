package crypto

import (
	"crypto/ecdsa"
	"github.com/clearmatics/autonity/core/types"
	"math/big"
	"sort"
	"strings"
	"testing"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus/tendermint/config"
	"github.com/clearmatics/autonity/consensus/tendermint/validator"
	"github.com/clearmatics/autonity/crypto"
)

func TestCheckValidatorSignature(t *testing.T) {
	vset, keys := newTestValidatorSet(5)

	// 1. Positive test: sign with validator's key should succeed
	data := []byte("dummy data")
	hashData := crypto.Keccak256(data)
	for i, k := range keys {
		// Sign
		sig, err := crypto.Sign(hashData, k)
		if err != nil {
			t.Errorf("error mismatch: have %v, want nil", err)
		}
		// CheckValidatorSignature should succeed
		addr, err := CheckValidatorSignature(vset, data, sig)
		if err != nil {
			t.Errorf("error mismatch: have %v, want nil", err)
		}
		val := vset.GetByIndex(uint64(i))
		if addr != val.GetAddress() {
			t.Errorf("validator address mismatch: have %v, want %v", addr, val.GetAddress())
		}
	}

	// 2. Negative test: sign with any key other than validator's key should return error
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Errorf("error mismatch: have %v, want nil", err)
	}
	// Sign
	sig, err := crypto.Sign(hashData, key)
	if err != nil {
		t.Errorf("error mismatch: have %v, want nil", err)
	}

	// CheckValidatorSignature should return ErrUnauthorizedAddress
	addr, err := CheckValidatorSignature(vset, data, sig)
	if err.Error() != ErrUnauthorizedAddress.Error() {
		t.Errorf("error mismatch: have %v, want %v", err, ErrUnauthorizedAddress)
	}

	emptyAddr := common.Address{}
	if addr != emptyAddr {
		t.Errorf("address mismatch: have %v, want %v", addr, emptyAddr)
	}
}

func TestCheckValidatorSignatureInvalid(t *testing.T) {
	vset, keys := newTestValidatorSet(5)

	// 1. Positive test: sign with validator's key should succeed
	data := []byte("dummy data")
	hashData := crypto.Keccak256(data)
	for i, k := range keys {
		// Sign
		sig, err := crypto.Sign(hashData, k)
		if err != nil {
			t.Errorf("sign error mismatch: have %v, want nil", err)
		}

		sig = sig[1:]

		// CheckValidatorSignature should succeed
		addr, err := CheckValidatorSignature(vset, data, sig)
		if err.Error() != "invalid signature length" {
			t.Errorf("check error mismatch: have %v, want ErrUnauthorizedAddress", err)
		}

		val := vset.GetByIndex(uint64(i))
		if addr == val.GetAddress() {
			t.Errorf("validator address match: have %v, want != %v", addr, val.GetAddress())
		}
	}

	// 2. Negative test: sign with any key other than validator's key should return error
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Errorf("error mismatch: have %v, want nil", err)
	}
	// Sign
	sig, err := crypto.Sign(hashData, key)
	if err != nil {
		t.Errorf("error mismatch: have %v, want nil", err)
	}

	// CheckValidatorSignature should return ErrUnauthorizedAddress
	addr, err := CheckValidatorSignature(vset, data, sig)
	if err.Error() != ErrUnauthorizedAddress.Error() {
		t.Errorf("error mismatch: have %v, want %v", err, ErrUnauthorizedAddress)
	}

	emptyAddr := common.Address{}
	if addr != emptyAddr {
		t.Errorf("address mismatch: have %v, want %v", addr, emptyAddr)
	}
}

func TestCheckValidatorUnauthorizedAddress(t *testing.T) {
	vset, keys := newTestValidatorSet(5)

	// 1. Positive test: sign with validator's key should succeed
	data := []byte("dummy data")
	hashData := crypto.Keccak256(data)
	for i, k := range keys {
		// Sign
		if hashData[0] != 0 {
			hashData[0] = 0
		} else {
			hashData[0] = 1
		}

		sig, err := crypto.Sign(hashData, k)
		if err != nil {
			t.Errorf("sign error mismatch: have %v, want nil", err)
		}

		// CheckValidatorSignature should succeed
		addr, err := CheckValidatorSignature(vset, data, sig)
		if err != ErrUnauthorizedAddress {
			t.Errorf("check error mismatch: have %v, want ErrUnauthorizedAddress", err)
		}

		val := vset.GetByIndex(uint64(i))
		if addr == val.GetAddress() {
			t.Errorf("validator address match: have %v, want != %v", addr, val.GetAddress())
		}
	}

	// 2. Negative test: sign with any key other than validator's key should return error
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Errorf("error mismatch: have %v, want nil", err)
	}
	// Sign
	sig, err := crypto.Sign(hashData, key)
	if err != nil {
		t.Errorf("error mismatch: have %v, want nil", err)
	}

	// CheckValidatorSignature should return ErrUnauthorizedAddress
	addr, err := CheckValidatorSignature(vset, data, sig)
	if err.Error() != ErrUnauthorizedAddress.Error() {
		t.Errorf("error mismatch: have %v, want %v", err, ErrUnauthorizedAddress)
	}

	emptyAddr := common.Address{}
	if addr != emptyAddr {
		t.Errorf("address mismatch: have %v, want %v", addr, emptyAddr)
	}
}

func newTestValidatorSet(n int) (validator.Set, []*ecdsa.PrivateKey) {
	// generate validators
	keys := make(Keys, n)
	addrs := make(types.Committee, n)
	for i := 0; i < n; i++ {
		privateKey, _ := crypto.GenerateKey()
		keys[i] = privateKey
		addrs[i] = types.CommitteeMember{
			Address:     crypto.PubkeyToAddress(privateKey.PublicKey),
			VotingPower: new(big.Int).SetUint64(1),
		}
	}
	vset := validator.NewSet(addrs, config.RoundRobin)
	sort.Sort(keys) //Keys need to be sorted by its public key address
	return vset, keys
}

type Keys []*ecdsa.PrivateKey

func (slice Keys) Len() int {
	return len(slice)
}

func (slice Keys) Less(i, j int) bool {
	return strings.Compare(crypto.PubkeyToAddress(slice[i].PublicKey).String(), crypto.PubkeyToAddress(slice[j].PublicKey).String()) < 0
}

func (slice Keys) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}
