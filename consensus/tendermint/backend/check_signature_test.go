package backend

import (
	"crypto/ecdsa"
	"math/big"
	"strings"
	"testing"

	"github.com/autonity/autonity/core/types"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/crypto"
)

var ErrUnauthorizedAddress = errors.New("unauthorized address")

func CheckValidatorSignature(previousHeader *types.Header, data []byte, sig []byte) (common.Address, error) {
	// 1. Get signature address
	signer, err := types.GetSignatureAddress(data, sig)
	if err != nil {
		log.Error("Failed to get signer address", "err", err)
		return common.Address{}, err
	}

	// 2. Check validator
	val := previousHeader.CommitteeMember(signer)
	if val == nil {
		return common.Address{}, ErrUnauthorizedAddress
	}

	return val.Address, nil
}

func TestCheckValidatorSignature(t *testing.T) {
	header, keys := newTestHeader(5)

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
		addr, err := CheckValidatorSignature(header, data, sig)
		if err != nil {
			t.Errorf("error mismatch: have %v, want nil", err)
		}
		val := header.Committee[i]
		if addr != val.Address {
			t.Errorf("validator address mismatch: have %v, want %v", addr, val.Address)
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
	addr, err := CheckValidatorSignature(header, data, sig)
	if err.Error() != ErrUnauthorizedAddress.Error() {
		t.Errorf("error mismatch: have %v, want %v", err, ErrUnauthorizedAddress)
	}

	emptyAddr := common.Address{}
	if addr != emptyAddr {
		t.Errorf("address mismatch: have %v, want %v", addr, emptyAddr)
	}
}

func TestCheckValidatorSignatureInvalid(t *testing.T) {
	header, keys := newTestHeader(5)

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
		addr, err := CheckValidatorSignature(header, data, sig)
		if err.Error() != "invalid signature length" {
			t.Errorf("check error mismatch: have %v, want ErrUnauthorizedAddress", err)
		}

		val := header.Committee[i]
		if addr == val.Address {
			t.Errorf("validator address match: have %v, want != %v", addr, val.Address)
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
	addr, err := CheckValidatorSignature(header, data, sig)
	if err.Error() != ErrUnauthorizedAddress.Error() {
		t.Errorf("error mismatch: have %v, want %v", err, ErrUnauthorizedAddress)
	}

	emptyAddr := common.Address{}
	if addr != emptyAddr {
		t.Errorf("address mismatch: have %v, want %v", addr, emptyAddr)
	}
}

func TestCheckValidatorUnauthorizedAddress(t *testing.T) {
	header, keys := newTestHeader(5)

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
		addr, err := CheckValidatorSignature(header, data, sig)
		if err != ErrUnauthorizedAddress {
			t.Errorf("check error mismatch: have %v, want ErrUnauthorizedAddress", err)
		}

		val := header.Committee[i]
		if addr == val.Address {
			t.Errorf("validator address match: have %v, want != %v", addr, val.Address)
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
	addr, err := CheckValidatorSignature(header, data, sig)
	if err.Error() != ErrUnauthorizedAddress.Error() {
		t.Errorf("error mismatch: have %v, want %v", err, ErrUnauthorizedAddress)
	}

	emptyAddr := common.Address{}
	if addr != emptyAddr {
		t.Errorf("address mismatch: have %v, want %v", addr, emptyAddr)
	}
}

func newTestHeader(n int) (*types.Header, []*ecdsa.PrivateKey) {
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
	h := &types.Header{
		Committee: addrs,
	}
	return h, keys
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
