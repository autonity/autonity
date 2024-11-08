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
	"fmt"

	"github.com/autonity/autonity/accounts/abi"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/consensus/tendermint/core/interfaces"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/params"
	"github.com/autonity/autonity/rpc"
)

// API is a user facing RPC API to dump BFT state
type API struct {
	chain      consensus.ChainReader
	tendermint *Backend
}

// GetCommittee retrieves the list of authorized committee at the specified block.
func (api *API) GetCommittee(number *rpc.BlockNumber) (*types.Committee, error) {
	if number == nil {
		return nil, fmt.Errorf("block number cannot be nil")
	}
	epoch, err := api.tendermint.EpochByHeight(uint64(*number))
	if err != nil {
		return nil, err
	}
	return epoch.Committee, nil
}

// GetCommitteeAtHash retrieves the state snapshot at a given block.
func (api *API) GetCommitteeAtHash(hash common.Hash) (*types.Committee, error) {
	header := api.chain.GetHeaderByHash(hash)
	if header == nil {
		return nil, errUnknownBlock
	}
	epoch, err := api.tendermint.EpochByHeight(header.Number.Uint64())
	if err != nil {
		return nil, err
	}
	return epoch.Committee, nil
}

// Get Autonity contract address
func (api *API) GetContractAddress() common.Address {
	return params.AutonityContractAddress
}

// Get Autonity contract ABI
func (api *API) GetContractABI() *abi.ABI {
	return api.tendermint.GetContractABI()
}

// Get current white list
func (api *API) GetCommitteeEnodes() []string {
	return api.tendermint.CommitteeEnodes()
}

// Get current tendermint's core state
func (api *API) GetCoreState() interfaces.CoreState {
	return api.tendermint.CoreState()
}
