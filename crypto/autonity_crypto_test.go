package crypto

import (
	"encoding/hex"
	"github.com/autonity/autonity/crypto/bls"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"testing"
)

func TestValidateValidatorKeyProof(t *testing.T) {
	privKey, err := GenerateKey()
	require.NoError(t, err)
	address := PubkeyToAddress(privKey.PublicKey)

	validatorKey, err := bls.SecretKeyFromECDSAKey(privKey)
	require.NoError(t, err)

	proof, err := GenerateValidatorKeyProof(validatorKey, address.Bytes())
	require.NoError(t, err)

	sig, err := bls.SignatureFromBytes(proof)
	require.NoError(t, err)

	err = ValidateValidatorKeyProof(validatorKey.PublicKey(), sig, address.Bytes())
	require.NoError(t, err)
}

func TestSaveNodeKey(t *testing.T) {
	f, err := ioutil.TempFile("", "save_node_key_test.*.txt")
	if err != nil {
		t.Fatal(err)
	}
	file := f.Name()
	f.Close()
	defer os.Remove(file)

	key, err := GenerateKey()
	require.NoError(t, err)

	derivedKey, err := bls.SecretKeyFromECDSAKey(key)
	require.NoError(t, err)

	err = SaveNodeKey(file, key, derivedKey)
	require.NoError(t, err)

	loadedKey, loadedDerivedKey, err := LoadNodeKey(file)
	require.NoError(t, err)

	require.Equal(t, loadedKey, key)
	require.Equal(t, loadedDerivedKey, derivedKey)
}

func TestHexToNodeKey(t *testing.T) {
	key, err := GenerateKey()
	require.NoError(t, err)

	derivedKey, err := bls.SecretKeyFromECDSAKey(key)
	require.NoError(t, err)

	keyHex := hex.EncodeToString(FromECDSA(key))
	derivedKeyHex := hex.EncodeToString(derivedKey.Marshal())

	parsedKey, parsedValidatorKey, err := HexToNodeKey(keyHex + derivedKeyHex)
	require.NoError(t, err)

	require.Equal(t, key, parsedKey)
	require.Equal(t, true, key.Equal(parsedKey))
	require.Equal(t, derivedKey.Hex(), parsedValidatorKey.Hex())
}

func TestLoadNodeKey(t *testing.T) {
	tests := []struct {
		input string
		err   string
	}{
		// good
		{input: "e17cdf649e8f53e0a6afb5f032a162219ab698bf21267b1592e3d5ee110e86721465269f128ab506bed287f36a993916409f4383b6ee1306f7e00f739c18a16f"},
		{input: "e17cdf649e8f53e0a6afb5f032a162219ab698bf21267b1592e3d5ee110e86721465269f128ab506bed287f36a993916409f4383b6ee1306f7e00f739c18a16f\n"},
		{input: "e17cdf649e8f53e0a6afb5f032a162219ab698bf21267b1592e3d5ee110e86721465269f128ab506bed287f36a993916409f4383b6ee1306f7e00f739c18a16f\n\r"},
		{input: "e17cdf649e8f53e0a6afb5f032a162219ab698bf21267b1592e3d5ee110e86721465269f128ab506bed287f36a993916409f4383b6ee1306f7e00f739c18a16f\r\n"},
		{input: "e17cdf649e8f53e0a6afb5f032a162219ab698bf21267b1592e3d5ee110e86721465269f128ab506bed287f36a993916409f4383b6ee1306f7e00f739c18a16f\n\n"},
		{input: "e17cdf649e8f53e0a6afb5f032a162219ab698bf21267b1592e3d5ee110e86721465269f128ab506bed287f36a993916409f4383b6ee1306f7e00f739c18a16f\n\r"},
		// bad
		{
			input: "e17cdf649e8f53e0a6afb5f032a162219ab698bf21267b1592e3d5ee110e8672",
			err:   "key file too short, want 128 hex characters",
		},
		{
			input: "e17cdf649e8f53e0a6afb5f032a162219ab698bf21267b1592e3d5ee110e8672\n",
			err:   "key file too short, want 128 hex characters",
		},
		{
			input: "e17cdf649e8f53e0a6afb5f032a162219ab698bf21267b1592e3d5ee110e86721465269f128ab506bed287f36a993916409f4383b6ee1306f7e00f739c18a16fX",
			err:   "invalid character 'X' at end of key file",
		},
		{
			input: "e17cdf649e8f53e0a6afb5f032a162219ab698bf21267b1592e3d5ee110e86721465269f128ab506bed287f36a993916409f4383b6ee1306f7e00f739c18a16f\n\n\n",
			err:   "key file too long, want 128 hex characters",
		},
	}

	for _, test := range tests {
		f, err := ioutil.TempFile("", "load_bls_key_test.*.txt")
		if err != nil {
			t.Fatal(err)
		}
		filename := f.Name()
		f.WriteString(test.input)
		f.Close()

		_, _, err = LoadNodeKey(filename)
		switch {
		case err != nil && test.err == "":
			t.Fatalf("unexpected error for input %q:\n  %v", test.input, err)
		case err != nil && err.Error() != test.err:
			t.Fatalf("wrong error for input %q:\n  %v", test.input, err)
		case err == nil && test.err != "":
			t.Fatalf("LoadNodeKey did not return error for input %q", test.input)
		}
	}
}
