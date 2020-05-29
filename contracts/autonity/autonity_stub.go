package autonity

import (
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/core/state"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/params"
	"math/big"
)

// Minimal implementation of autonity.Contrat in order ot make miner tests passing.
type stubContract struct {
	committee types.Committee
}

func NewStubContract(config *params.AutonityContractGenesis) Contract {
	validators := config.GetValidatorUsers()
	committee := make(types.Committee, len(validators))
	for i, val := range validators {
		committee[i] = types.CommitteeMember{
			Address:     val.Address,
			VotingPower: new(big.Int).SetUint64(val.Stake),
		}
	}
	return &stubContract{committee: committee}
}

func (s *stubContract) GetCommittee(header *types.Header, statedb *state.StateDB) (types.Committee, error) {
	return s.committee, nil
}

func (s *stubContract) GetMinimumGasPrice(block *types.Block, db *state.StateDB) (uint64, error) {
	return 1000, nil
}

func (s *stubContract) FinalizeAndGetCommittee(txs types.Transactions, r types.Receipts, h *types.Header, db *state.StateDB) (types.Committee, *types.Receipt, error) {
	return s.committee, new(types.Receipt), nil
}

func (s *stubContract) MeasureMetricsOfNetworkEconomic(header *types.Header, stateDB *state.StateDB) {
	return
}

func (s *stubContract) UpdateEnodesWhitelist(state *state.StateDB, block *types.Block) error {
	return nil
}

func (s *stubContract) GetContractABI() string {
	return ""
}

func (s *stubContract) Address() common.Address {
	return common.Address{}
}

func (s *stubContract) DeployAutonityContract(chainConfig *params.ChainConfig, header *types.Header, statedb *state.StateDB) (common.Address, error) {
	return common.Address{}, nil
}

func (s *stubContract) GetWhitelist(block *types.Block, db *state.StateDB) (*types.Nodes, error) {
	return nil, nil
}
