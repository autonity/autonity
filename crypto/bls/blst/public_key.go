package blst

import (
	"encoding/hex"
	"fmt"
	"github.com/autonity/autonity/crypto/bls/common"
	lru "github.com/hashicorp/golang-lru"
	"github.com/pkg/errors"
	"reflect"
)

var maxKeys = 10000

var PubkeyCache, _ = lru.NewARC(maxKeys)

// PublicKey used in the BLS signature scheme.
type PublicKey struct {
	p *blstPublicKey
}

// PublicKeyBytesFromString decode a BLS public key from a hex format string.
func PublicKeyBytesFromString(pubKey string) ([]byte, error) {
	// Since a public key is 48 bytes, the hex string representing it must be 48*2 = 96 chars + "0x" prefix = 98 chars
	if len(pubKey) != common.BLSPubKeyHexStringLength {
		return nil, fmt.Errorf("invalid BLS key string")
	}
	return hex.DecodeString(pubKey[2:])
}

// PublicKeyFromBytes creates a BLS public key from a  BigEndian byte slice.
func PublicKeyFromBytes(pubKey []byte) (common.BLSPublicKey, error) {

	if len(pubKey) != common.BLSPubkeyLength {
		return nil, fmt.Errorf("public key must be %d bytes", common.BLSPubkeyLength)
	}

	if cv, ok := PubkeyCache.Get(string(pubKey)); ok {
		return cv.(*PublicKey).Copy(), nil
	}

	// Subgroup check NOT done when decompressing pubkey.
	p := new(blstPublicKey).Uncompress(pubKey)
	if p == nil {
		return nil, errors.New("could not unmarshal bytes into public key")
	}

	// Subgroup and infinity check
	if !p.KeyValidate() {
		// NOTE: the error is not quite accurate since it includes group check
		return nil, common.ErrInfinitePubKey
	}

	pubKeyObj := &PublicKey{p: p}
	copiedKey := pubKeyObj.Copy()

	PubkeyCache.Add(string(pubKey), copiedKey)
	return pubKeyObj, nil
}

// AggregatePublicKeys aggregates the provided raw public keys into a single key.
func AggregatePublicKeys(pubs [][]byte) (common.BLSPublicKey, error) {
	if len(pubs) == 0 {
		return nil, fmt.Errorf("empty pub-key set for key aggregation")
	}

	agg := new(blstAggregatePublicKey)
	mulP1 := make([]*blstPublicKey, 0, len(pubs))
	for _, pubkey := range pubs {
		pubKeyObj, err := PublicKeyFromBytes(pubkey)
		if err != nil {
			return nil, err
		}
		mulP1 = append(mulP1, pubKeyObj.(*PublicKey).p)
	}
	// No group check needed here since it is done in PublicKeyFromBytes
	// Note the checks could be moved from PublicKeyFromBytes into Aggregate
	// and take advantage of multi-threading.
	if agg.Aggregate(mulP1, false) {
		return &PublicKey{p: agg.ToAffine()}, nil
	}
	return nil, fmt.Errorf("cannot aggregate public keys")
}

// Marshal a public key into a LittleEndian byte slice.
func (p *PublicKey) Marshal() []byte {
	return p.p.Compress()
}

// Copy the public key to a new pointer reference.
func (p *PublicKey) Copy() common.BLSPublicKey {
	np := *p.p
	return &PublicKey{p: &np}
}

// IsInfinite checks if the public key is infinite.
func (p *PublicKey) IsInfinite() bool {
	zeroKey := new(blstPublicKey)
	return p.p.Equals(zeroKey)
}

// Aggregate two public keys.
func (p *PublicKey) Aggregate(p2 common.BLSPublicKey) (common.BLSPublicKey, error) {
	if reflect.TypeOf(p2) != reflect.TypeOf(p) {
		return nil, fmt.Errorf("wrong type of bls public key")
	}
	agg := new(blstAggregatePublicKey)
	// No group check here since it is checked at decompression time
	agg.Add(p.p, false)
	agg.Add(p2.(*PublicKey).p, false)
	p.p = agg.ToAffine()
	return p, nil
}

func (p *PublicKey) Hex() string {
	return "0x" + hex.EncodeToString(p.Marshal())
}
