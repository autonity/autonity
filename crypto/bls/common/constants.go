package common

const (
	BLSSecretKeyLength       = 32
	BLSPubkeyLength          = 48
	BLSSignatureLength       = 96
	BLSPubKeyHexStringLength = 98
)

// ZeroSecretKey represents a zero secret key.
var ZeroSecretKey = [32]byte{}

// InfinitePublicKey represents an infinite public key.
var InfinitePublicKey = [48]byte{0xC0}
