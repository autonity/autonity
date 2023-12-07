package blst

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMarshalUnmarshal(t *testing.T) {
	priv, err := RandKey()
	require.NoError(t, err)
	b := priv.Marshal()
	b32 := b
	pk, err := SecretKeyFromBytes(b32[:])
	require.NoError(t, err)
	pk2, err := SecretKeyFromBytes(b32[:])
	require.NoError(t, err)
	require.Equal(t, pk.Marshal(), priv.Marshal(), "Keys not equal")
	require.Equal(t, pk.Marshal(), pk2.Marshal(), "Keys not equal")
}

func TestSecretKeyFromBytes(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		err   error
	}{
		{
			name: "Nil",
			err:  ErrInvalidSecreteKey,
		},
		{
			name:  "Empty",
			input: []byte{},
			err:   ErrInvalidSecreteKey,
		},
		{
			name:  "Short",
			input: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			err:   ErrInvalidSecreteKey,
		},
		{
			name:  "Long",
			input: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			err:   ErrInvalidSecreteKey,
		},
		{
			name:  "Bad",
			input: []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
			err:   ErrSecretUnmarshal,
		},
		{
			name:  "Good",
			input: []byte{0x25, 0x29, 0x5f, 0x0d, 0x1d, 0x59, 0x2a, 0x90, 0xb3, 0x33, 0xe2, 0x6e, 0x85, 0x14, 0x97, 0x08, 0x20, 0x8e, 0x9f, 0x8e, 0x8b, 0xc1, 0x8f, 0x6c, 0x77, 0xbd, 0x62, 0xf8, 0xad, 0x7a, 0x68, 0x66},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			res, err := SecretKeyFromBytes(test.input)
			if test.err != nil {
				require.NotEqual(t, nil, err, "No error returned")
				require.Error(t, test.err, err.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, 0, bytes.Compare(res.Marshal(), test.input))
			}
		})
	}
}

func TestSerialize(t *testing.T) {
	rk, err := RandKey()
	require.NoError(t, err)
	b := rk.Marshal()

	_, err = SecretKeyFromBytes(b)
	require.NoError(t, err)
}

func TestSecretKey_Hex(t *testing.T) {
	sKey, err := RandKey()
	require.NoError(t, err)

	str := sKey.Hex()
	b, err := hex.DecodeString(str[2:])
	require.NoError(t, err)

	sk, err := SecretKeyFromBytes(b)
	require.NoError(t, err)

	require.Equal(t, sKey, sk)
}
