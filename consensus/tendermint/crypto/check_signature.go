package crypto

import (
	"errors"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus/tendermint/committee"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/log"
)

var ErrUnauthorizedAddress = errors.New("unauthorized address")

func CheckValidatorSignature(valSet committee.Set, data []byte, sig []byte) (common.Address, error) {
	// 1. Get signature address
	signer, err := types.GetSignatureAddress(data, sig)
	if err != nil {
		log.Error("Failed to get signer address", "err", err)
		return common.Address{}, err
	}

	// 2. Check validator
	_, val, err := valSet.GetByAddress(signer)
	if err != nil {
		return common.Address{}, ErrUnauthorizedAddress
	}

	return val.Address, nil
}
