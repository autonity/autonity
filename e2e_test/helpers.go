package test

import (
	"crypto/rand"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/core/messageutils"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/crypto"
	"math/big"
)

var autonityContractAddr = crypto.CreateAddress(common.Address{}, 0)

func MsgPropose(address common.Address, block *types.Block, h uint64, r int64, vr int64) *messageutils.Message {
	proposal := messageutils.NewProposal(r, new(big.Int).SetUint64(h), vr, block)
	v, err := messageutils.Encode(proposal)
	if err != nil {
		return nil
	}
	return &messageutils.Message{
		Code:          messageutils.MsgProposal,
		Msg:           v,
		Address:       address,
		CommittedSeal: []byte{},
	}
}

// GenerateRandomBytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}
