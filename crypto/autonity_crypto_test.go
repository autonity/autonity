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

	validatorKey, err := blst.RandKey()
	require.NoError(t, err)

	proof, err := POPProof(validatorKey, address.Bytes())
	require.NoError(t, err)

	sig, err := blst.SignatureFromBytes(proof)
	require.NoError(t, err)

	err = POPVerify(validatorKey.PublicKey(), sig, address.Bytes())
	require.NoError(t, err)
}
