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
	"github.com/clearmatics/autonity/contracts/autonity"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/rpc"
)

// API is a user facing RPC API to dump BFT state
type API struct {
	chain        consensus.ChainReader
	tendermint   *Backend
	getCommittee func(header *types.Header, chain consensus.ChainReader) (types.Committee, error)
}

// GetCommittee retrieves the list of authorized committee at the specified block.
func (api *API) GetCommittee(number *rpc.BlockNumber) (types.Committee, error) {
	header := api.chain.GetHeaderByNumber(uint64(*number))
	if header == nil {
		return nil, errUnknownBlock
	}
	committee, err := api.getCommittee(header, api.chain)
	if err != nil {
		return nil, err
	}
	return committee, nil
}

// GetCommitteeAtHash retrieves the state snapshot at a given block.
func (api *API) GetCommitteeAtHash(hash common.Hash) (types.Committee, error) {
	header := api.chain.GetHeaderByHash(hash)
	if header == nil {
		return nil, errUnknownBlock
	}
	committee, err := api.getCommittee(header, api.chain)
	if err != nil {
		return nil, err
	}
	return committee, nil
}

// Get Autonity contract address
func (api *API) GetContractAddress() common.Address {
	return autonity.ContractAddress
}

// Get Autonity contract ABI
func (api *API) GetContractABI() string {
	return api.tendermint.GetContractABI()
}

// Get current white list
func (api *API) GetWhitelist() []string {
	return api.tendermint.WhiteList()
}

// Get current tendermint's core state
func (api *API) GetCoreState() types.TendermintState {
	return api.tendermint.GetCoreState()
}
