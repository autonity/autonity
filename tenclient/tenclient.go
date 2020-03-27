package tenclient

import (
	"context"
	"math/big"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/rpc"
)

// TenClient provides an interface to interact with the tendermint namespace
// rpc api.
type TenClient struct {
	c *rpc.Client
}

// ContractAddresss returns the Autonity contract address.
func (tc *TenClient) ContractAddresss(ctx context.Context) (common.Address, error) {
	var result common.Address
	err := tc.c.CallContext(ctx, &result, "tendermint_getContractAddress")
	if err != nil {
		return common.Address{}, err
	}
	return result, err
}

// ContractABI returns the Autonity contract ABI.
func (tc *TenClient) ContractABI(ctx context.Context) (string, error) {
	var result string
	err := tc.c.CallContext(ctx, &result, "tendermint_getContractABI")
	if err != nil {
		return "", err
	}
	return result, err
}

// Whitelist returns the whitelist.
func (tc *TenClient) Whitelist(ctx context.Context) ([]string, error) {
	var result []string
	err := tc.c.CallContext(ctx, &result, "tendermint_getWhitelist")
	if err != nil {
		return nil, err
	}
	return result, err
}

// Committee returns the committee at the given blockNumber.
func (tc *TenClient) Committee(ctx context.Context, blockNumber *big.Int) (types.Committee, error) {
	var result types.Committee
	err := tc.c.CallContext(ctx, &result, "tendermint_getCommittee")
	if err != nil {
		return nil, err
	}
	return result, err
}

// Committee returns the committee at for the block identified by blockHash.
func (tc *TenClient) CommitteeAtHash(ctx context.Context, blockHash common.Hash) (types.Committee, error) {
	var result types.Committee
	err := tc.c.CallContext(ctx, &result, "tendermint_getCommitteeAtHash")
	if err != nil {
		return nil, err
	}
	return result, err
}
