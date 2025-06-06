package blst

import (
	"encoding/hex"
	"fmt"
	farmhash "github.com/dgryski/go-farm"
	"github.com/pkg/errors"

	"github.com/autonity/autonity/common/fixsizecache"
)

func hashKey(key string) uint {
	return uint(farmhash.Hash64([]byte(key)))
}

// parameters set assuming max committee size = 1000
var PubkeyCache = fixsizecache.New[string, PublicKey](997, 5, hashKey)

const HexPrefix = "0x"

// BlsPublicKey used in the BLS signature scheme.
type BlsPublicKey struct {
	p *blstPublicKey
}

// PublicKeyFromBytes creates a BLS public key from a  BigEndian byte slice.
func PublicKeyFromBytes(pubKey []byte) (PublicKey, error) {
	if len(pubKey) != BLSPubkeyLength {
		return nil, fmt.Errorf("public key must be %d bytes", BLSPubkeyLength)
	}

	if cachedPublicKey, ok := PubkeyCache.Get(string(pubKey)); ok {
		return cachedPublicKey.(*BlsPublicKey).Copy(), nil
	}

	// Subgroup check NOT done when decompressing pubkey.
	p := new(blstPublicKey).Uncompress(pubKey)
	if p == nil {
		return nil, errors.New("could not unmarshal bytes into public key")
	}

	// Subgroup and infinity check
	if !p.KeyValidate() {
		// NOTE: the error is not quite accurate since it includes group check
		return nil, ErrInfinitePubKey
	}

	pubKeyObj := &BlsPublicKey{p: p}
	copiedKey := pubKeyObj.Copy()

	PubkeyCache.Add(string(pubKey), copiedKey)
	return pubKeyObj, nil
}

// AggregatePublicKeys aggregates the provided deserialized public keys into a single key.
func AggregatePublicKeys(pubs []PublicKey) (PublicKey, error) {
	if len(pubs) == 0 {
		return nil, fmt.Errorf("empty pub-key set for key aggregation")
	}

	mulP1 := make([]*blstPublicKey, 0, len(pubs))
	for _, pubkey := range pubs {
		mulP1 = append(mulP1, pubkey.(*BlsPublicKey).p)
	}

	agg := new(blstAggregatePublicKey)
	// No group check needed here since it is done at deserialization (`PublicKeyFromBytes`)
	if agg.Aggregate(mulP1, false) {
		return &BlsPublicKey{p: agg.ToAffine()}, nil
	}
	return nil, fmt.Errorf("cannot aggregate public keys")
}

// does group check and infinity check on the public key
func (p *BlsPublicKey) Validate() bool {
	return p.p.KeyValidate()
}

// Marshal a public key into a LittleEndian byte slice.
func (p *BlsPublicKey) Marshal() []byte {
	return p.p.Compress()
}

// Copy the public key to a new pointer reference.
func (p *BlsPublicKey) Copy() PublicKey {
	np := *p.p
	return &BlsPublicKey{p: &np}
}

func (p *BlsPublicKey) Hex() string {
	return HexPrefix + hex.EncodeToString(p.Marshal())
}
