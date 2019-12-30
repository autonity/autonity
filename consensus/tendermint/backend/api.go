// Copyright 2017 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package backend

import (
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus"
	"github.com/clearmatics/autonity/consensus/tendermint/core"
	"github.com/clearmatics/autonity/rpc"
)

// API is a user facing RPC API to dump BFT state
type API struct {
	chain      consensus.ChainReader
	tendermint core.Backend
}

// GetValidators retrieves the list of authorized validators at the specified block.
func (api *API) GetValidators(number *rpc.BlockNumber) ([]common.Address, error) {
	validators := api.tendermint.Validators(uint64(*number)).List()
	addresses := make([]common.Address, len(validators))
	for i, validator := range validators {
		addresses[i] = validator.GetAddress()
	}
	return addresses, nil
}

// GetValidatorsAtHash retrieves the state snapshot at a given block.
func (api *API) GetValidatorsAtHash(hash common.Hash) ([]common.Address, error) {
	header := api.chain.GetHeaderByHash(hash)
	if header == nil {
		return nil, errUnknownBlock
	}

	validators := api.tendermint.Validators(header.Number.Uint64()).List()
	addresses := make([]common.Address, len(validators))
	for i, validator := range validators {
		addresses[i] = validator.GetAddress()
	}
	return addresses, nil
}

// Get Autonity contract address
func (api *API) GetContractAddress() common.Address {
	return api.tendermint.GetContractAddress()
}

// Get Autonity contract ABI
func (api *API) GetContractABI() string {
	return api.tendermint.GetContractABI()
}

// Get current white list
func (api *API) GetWhitelist() []string {
	return api.tendermint.WhiteList()
}
