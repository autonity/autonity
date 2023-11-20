package crypto

import (
	"github.com/autonity/autonity/crypto/bls"
	"github.com/stretchr/testify/require"
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
