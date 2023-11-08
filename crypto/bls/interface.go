package bls

import (
	"github.com/autonity/autonity/crypto/bls/common"
)

// PublicKey represents a BLS public key.
type PublicKey = common.BLSPublicKey

// SecretKey represents a BLS secret or private key.
type SecretKey = common.BLSSecretKey

// Signature represents a BLS signature.
type Signature = common.BLSSignature
