package crypto

import (
	"github.com/autonity/autonity/crypto/blst"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestVerifyOwnershipProof(t *testing.T) {
	privKey, err := GenerateKey()
	require.NoError(t, err)
	address := PubkeyToAddress(privKey.PublicKey)

	validatorKey, err := blst.SecretKeyFromECDSAKey(privKey.D.Bytes())
	require.NoError(t, err)

	proof, err := PopProof(validatorKey, address.Bytes())
	require.NoError(t, err)

	sig, err := blst.SignatureFromBytes(proof)
	require.NoError(t, err)

	err = PopVerify(validatorKey.PublicKey(), sig, address.Bytes())
	require.NoError(t, err)
}
