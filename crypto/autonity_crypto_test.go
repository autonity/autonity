package crypto

import (
	"encoding/hex"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/common/hexutil"
	"github.com/autonity/autonity/crypto/blst"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestSaveNodeKey(t *testing.T) {
	f, err := os.CreateTemp("", "save_node_key_test.*.txt")
	if err != nil {
		t.Fatal(err)
	}
	file := f.Name()
	f.Close()
	defer os.Remove(file)

	key, err := GenerateKey()
	require.NoError(t, err)

	consensusKey, err := blst.RandKey()
	require.NoError(t, err)

	err = SaveAutonityKeys(file, key, consensusKey)
	require.NoError(t, err)

	loadedKey, loadedDerivedKey, err := LoadAutonityKeys(file)
	require.NoError(t, err)

	require.Equal(t, loadedKey, key)
	require.Equal(t, loadedDerivedKey, consensusKey)
}

func TestHexToNodeKey(t *testing.T) {
	key, err := GenerateKey()
	require.NoError(t, err)

	consensusKey, err := blst.RandKey()
	require.NoError(t, err)

	keyHex := hex.EncodeToString(FromECDSA(key))
	derivedKeyHex := hex.EncodeToString(consensusKey.Marshal())

	parsedKey, parsedConsensusKey, err := HexToAutonityKeys(keyHex + derivedKeyHex)
	require.NoError(t, err)

	require.Equal(t, key, parsedKey)
	require.Equal(t, true, key.Equal(parsedKey))
	require.Equal(t, consensusKey.Hex(), parsedConsensusKey.Hex())
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
		f, err := os.CreateTemp("", "load_bls_key_test.*.txt")
		if err != nil {
			t.Fatal(err)
		}
		filename := f.Name()
		f.WriteString(test.input)
		f.Close()

		_, _, err = LoadAutonityKeys(filename)
		switch {
		case err != nil && test.err == "":
			t.Fatalf("unexpected error for input %q:\n  %v", test.input, err)
		case err != nil && err.Error() != test.err:
			t.Fatalf("wrong error for input %q:\n  %v", test.input, err)
		case err == nil && test.err != "":
			t.Fatalf("LoadAutonityKeys did not return error for input %q", test.input)
		}
	}
}

func TestPOPVerifier(t *testing.T) {
	privKey, err := GenerateKey()
	require.NoError(t, err)
	treasury := PubkeyToAddress(privKey.PublicKey)

	consensusKey, err := blst.RandKey()
	require.NoError(t, err)

	proof, err := BLSPOPProof(consensusKey, treasury.Bytes())
	require.NoError(t, err)

	sig, err := blst.SignatureFromBytes(proof)
	require.NoError(t, err)

	require.Equal(t, false, sig.IsZero())

	err = BLSPOPVerify(consensusKey.PublicKey(), sig, treasury.Bytes())
	require.NoError(t, err)
}

func TestAutonityPOPProof(t *testing.T) {
	treasury, err := GenerateKey()
	require.NoError(t, err)
	nodeKey, err := GenerateKey()
	require.NoError(t, err)
	oracleKey, err := GenerateKey()
	require.NoError(t, err)
	consensusKey, err := blst.RandKey()
	require.NoError(t, err)

	msg := PubkeyToAddress(treasury.PublicKey).Hex()
	autonityPOP, err := AutonityPOPProof(nodeKey, oracleKey, msg, consensusKey)
	require.NoError(t, err)
	require.Equal(t, AutonityPOPLen, len(autonityPOP))

	err = autonityPOPVerify(autonityPOP, msg, PubkeyToAddress(nodeKey.PublicKey), PubkeyToAddress(oracleKey.PublicKey), consensusKey.PublicKey().Marshal())
	require.NoError(t, err)
}

func autonityPOPVerify(signatures []byte, treasuryHex string, nodeAddress, oracleAddress common.Address, consensusKey []byte) error {
	if len(signatures) != AutonityPOPLen {
		return ErrorInvalidPOP
	}

	msg, err := hexutil.Decode(treasuryHex)
	if err != nil {
		return err
	}

	hash := POPMsgHash(msg)
	if err = ecdsaPOPVerify(signatures[0:common.SealLength], hash, nodeAddress); err != nil {
		return err
	}

	blsSigOffset := common.SealLength * 2
	if err = ecdsaPOPVerify(signatures[common.SealLength:blsSigOffset], hash, oracleAddress); err != nil {
		return err
	}

	validatorSig, err := blst.SignatureFromBytes(signatures[blsSigOffset:])
	if err != nil {
		return err
	}

	// check zero signature.
	if validatorSig.IsZero() {
		return ErrorInvalidPOP
	}

	// check zero public key.
	blsPubKey, err := blst.PublicKeyFromBytes(consensusKey)
	if err != nil {
		return err
	}
	return BLSPOPVerify(blsPubKey, validatorSig, msg)
}

func ecdsaPOPVerify(sig []byte, hash common.Hash, expectedSigner common.Address) error {
	signer, err := SigToAddr(hash[:], sig)
	if err != nil {
		return err
	}

	if signer != expectedSigner {
		return ErrorInvalidSigner
	}

	return nil
}
