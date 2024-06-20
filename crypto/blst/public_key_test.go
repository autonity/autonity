package blst

import (
	"bytes"
	"encoding/hex"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPublicKeyFromBytes(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		err   error
	}{
		{
			name: "Nil",
			err:  errors.New("public key must be 48 bytes"),
		},
		{
			name:  "Empty",
			input: []byte{},
			err:   errors.New("public key must be 48 bytes"),
		},
		{
			name:  "Short",
			input: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			err:   errors.New("public key must be 48 bytes"),
		},
		{
			name:  "Long",
			input: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			err:   errors.New("public key must be 48 bytes"),
		},
		{
			name:  "Bad",
			input: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			err:   errors.New("could not unmarshal bytes into public key"),
		},
		{
			name:  "Good",
			input: []byte{0xa9, 0x9a, 0x76, 0xed, 0x77, 0x96, 0xf7, 0xbe, 0x22, 0xd5, 0xb7, 0xe8, 0x5d, 0xee, 0xb7, 0xc5, 0x67, 0x7e, 0x88, 0xe5, 0x11, 0xe0, 0xb3, 0x37, 0x61, 0x8f, 0x8c, 0x4e, 0xb6, 0x13, 0x49, 0xb4, 0xbf, 0x2d, 0x15, 0x3f, 0x64, 0x9f, 0x7b, 0x53, 0x35, 0x9f, 0xe8, 0xb9, 0x4a, 0x38, 0xe4, 0x4c},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			res, err := PublicKeyFromBytes(test.input)
			if test.err != nil {
				require.NotEqual(t, nil, err, "No error returned")
				require.Error(t, test.err, err.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, test.input, res.Marshal(), "out put key is not match with input")
			}
		})
	}
}

func TestPublicKey_Copy(t *testing.T) {
	priv, err := RandKey()
	require.NoError(t, err)
	pubkeyA := priv.PublicKey()
	pubkeyBytes := pubkeyA.Marshal()

	pubkeyB := pubkeyA.Copy()
	pubkeyBBytes := pubkeyB.Marshal()

	require.Equal(t, true, bytes.Equal(pubkeyBytes, pubkeyBBytes))
}

func TestPublicKey_MarshalUnMarshal(t *testing.T) {
	priv, err := RandKey()
	require.NoError(t, err)
	pubkey := priv.PublicKey()

	pubkeyBytes := pubkey.Marshal()

	pKeyFromBytes, err := PublicKeyFromBytes(pubkeyBytes)
	require.NoError(t, err)

	require.Equal(t, pubkey, pKeyFromBytes)
}

func TestPublicKey_Hex(t *testing.T) {
	priv, err := RandKey()
	require.NoError(t, err)
	pubkey := priv.PublicKey()

	str := pubkey.Hex()

	b, err := hex.DecodeString(str[2:])
	require.NoError(t, err)

	pb, err := PublicKeyFromBytes(b)
	require.NoError(t, err)

	require.Equal(t, pubkey, pb)
}

func TestPublicKeyBytesFromString(t *testing.T) {
	priv, err := RandKey()
	require.NoError(t, err)

	sig := priv.Sign([]byte{0})

	pubKeyStr := priv.PublicKey().Hex()
	pubKeyByte, err := PublicKeyBytesFromString(pubKeyStr)
	require.NoError(t, err)

	pubKey, err := PublicKeyFromBytes(pubKeyByte)
	require.NoError(t, err)
	require.Equal(t, true, sig.Verify(pubKey, []byte{0}))
}

func TestAggregatePublicKey(t *testing.T) {
	priv1, err := RandKey()
	require.NoError(t, err)

	priv2, err := RandKey()
	require.NoError(t, err)

	_, err = priv1.PublicKey().Aggregate(priv2.PublicKey())
	require.NoError(t, err)
}

func TestAggregateRawPublicKeys(t *testing.T) {
	priv1, err := RandKey()
	require.NoError(t, err)

	priv2, err := RandKey()
	require.NoError(t, err)

	expectedKey, err := priv1.PublicKey().Aggregate(priv2.PublicKey())
	require.NoError(t, err)

	aggregatedKey, err := AggregateRawPublicKeys([][]byte{priv1.PublicKey().Marshal(), priv2.PublicKey().Marshal()})
	require.NoError(t, err)

	require.Equal(t, expectedKey.Marshal(), aggregatedKey.Marshal())
}

func TestAggregatePublicKeys(t *testing.T) {
	priv1, err := RandKey()
	require.NoError(t, err)
	pub1 := priv1.PublicKey()
	pub1Bytes := pub1.Marshal()

	priv2, err := RandKey()
	require.NoError(t, err)
	pub2 := priv2.PublicKey()
	pub2Bytes := pub2.Marshal()

	expectedKey, err := priv1.PublicKey().Aggregate(pub2)
	require.NoError(t, err)

	aggregatedKey, err := AggregatePublicKeys([]PublicKey{pub1, pub2})
	require.NoError(t, err)

	require.Equal(t, expectedKey.Marshal(), aggregatedKey.Marshal())

	// check no changes to the original pubkeys
	require.Equal(t, pub1Bytes, pub1.Marshal())
	require.Equal(t, pub2Bytes, pub2.Marshal())
}
