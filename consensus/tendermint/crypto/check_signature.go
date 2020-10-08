package crypto

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/binary"
	"errors"
	"math/big"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/crypto"
	"github.com/clearmatics/autonity/log"
)

var ErrUnauthorizedAddress = errors.New("unauthorized address")

// SignHeader signs the given header with the given private key.
func SignHeader(h *types.Header, priv *ecdsa.PrivateKey) error {
	hashData := crypto.Keccak256(types.SigHash(h).Bytes())
	signature, err := crypto.Sign(hashData, priv)
	if err != nil {
		return err
	}
	err = types.WriteSeal(h, signature)
	if err != nil {
		return err
	}
	return nil
}

func CheckValidatorSignature(previousHeader *types.Header, data []byte, sig []byte) (common.Address, error) {
	// 1. Get signature address
	signer, err := types.GetSignatureAddress(data, sig)
	if err != nil {
		log.Error("Failed to get signer address", "err", err)
		return common.Address{}, err
	}

	// 2. Check validator
	val := previousHeader.CommitteeMember(signer)
	if val == nil {
		return common.Address{}, ErrUnauthorizedAddress
	}

	return val.Address, nil
}

//  BuildCommitment returns byte representation of the valueHash, height and
//  round.
func BuildCommitment(valueHash common.Hash, height *big.Int, round int64) []byte {
	var buf bytes.Buffer
	roundBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(roundBytes, uint64(round))
	buf.Write(roundBytes)
	buf.Write(height.Bytes())
	buf.Write(valueHash.Bytes())
	return buf.Bytes()
}
