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

package tendermint

import (
	"github.com/clearmatics/autonity/autonity"
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/log"
	"github.com/clearmatics/autonity/rpc"
)

// getCommittee retrieves the committee for the given header.
func getCommittee(header *types.Header, chain consensus.ChainReader) (types.Committee, error) {
	parent := chain.GetHeaderByHash(header.ParentHash)
	if parent == nil {
		return nil, errUnknownBlock
	}
	return parent.Committee, nil
}

// API is a user facing RPC API to dump BFT state
type API struct {
	chain            consensus.ChainReader
	autonityContract *autonity.Contract
	blockRetriever   *BlockReader
	getCommittee     func(header *types.Header, chain consensus.ChainReader) (types.Committee, error)
}

func NewAPI(chain consensus.ChainReader, ac *autonity.Contract, lbr *BlockReader) *API {
	return &API{
		chain:            chain,
		autonityContract: ac,
		blockRetriever:   lbr,
		getCommittee:     getCommittee,
	}
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

// Get Autonity contract Address
func (api *API) GetContractAddress() common.Address {
	return autonity.ContractAddress
}

func (api *API) GetContractABI() string {
	// after the contract is upgradable, call it from contract object rather than from conf.
	return api.autonityContract.StringABI()
}

func (api *API) GetAccusations() []autonity.OnChainProof {
	b, err := api.blockRetriever.LatestBlock()
	if err != nil {
		panic(err)
	}
	state, err := api.blockRetriever.BlockState(b.Root())
	if err != nil {
		panic(err)
	}

	proofs := api.autonityContract.GetAccusations(b.Header(), state)
	return proofs
}

func (api *API) GetMisbehaviours() []autonity.OnChainProof {
	b, err := api.blockRetriever.LatestBlock()
	if err != nil {
		panic(err)
	}
	state, err := api.blockRetriever.BlockState(b.Root())
	if err != nil {
		panic(err)
	}

	proofs := api.autonityContract.GetMisBehaviours(b.Header(), state)
	return proofs
}

// Whitelist for the current block
func (api *API) GetWhitelist() []string {
	// TODO this should really return errors
	b, err := api.blockRetriever.LatestBlock()
	if err != nil {
		panic(err)
	}
	state, err := api.blockRetriever.BlockState(b.Root())
	if err != nil {
		// TODO: Decide how to log errors correctly
		log.Error("Failed to get block white list", "err", err)
		return nil
	}

	enodes, err := api.autonityContract.GetWhitelist(b, state)
	if err != nil {
		// TODO: Decide how to log errors correctly
		log.Error("Failed to get block white list", "err", err)
		return nil
	}

	return enodes.StrList
}

// Get current tendermint's core state
// func (api *API) GetCoreState() TendermintState {
// 	return api.tendermint.CoreState()
// }