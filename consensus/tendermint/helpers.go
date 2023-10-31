package tendermint

import (
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/crypto"
	lru "github.com/hashicorp/golang-lru"
)

var (
	signatureCacheSize = 5000 // Number of recent addresses from ecrecover
	signatureCache, _  = lru.NewARC(signatureCacheSize)
)

func SigToAddr(hash common.Hash, sig []byte) (common.Address, error) {
	var cacheKey [65 + 32]byte
	copy(cacheKey[:], hash[:])
	copy(cacheKey[32:], sig)

	if data, ok := signatureCache.Get(cacheKey); ok {
		return data.(common.Address), nil
	}
	addr, err := crypto.SigToAddr(hash[:], sig)
	if err != nil {
		return common.Address{}, err
	}
	signatureCache.Add(cacheKey, addr)
	return addr, nil
}
