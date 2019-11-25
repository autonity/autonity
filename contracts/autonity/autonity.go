package autonity

import (
	"errors"
	"github.com/clearmatics/autonity/accounts/abi"
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus"
	"github.com/clearmatics/autonity/core/state"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/core/vm"
	"github.com/clearmatics/autonity/log"
	"github.com/clearmatics/autonity/params"
	"math"
	"math/big"
	"reflect"
	"sort"
	"strings"
	"sync"
)

var ErrAutonityContract = errors.New("could not call Autonity contract")

type ChainContext interface {
	// Engine retrieves the chain's consensus engine.
	Engine() consensus.Engine

	// GetHeader returns the hash corresponding to their hash.
	GetHeader(common.Hash, uint64) *types.Header
}
type Blockchainer interface {
	ChainContext
	GetVMConfig() *vm.Config
	Config() *params.ChainConfig

	UpdateEnodeWhitelist(newWhitelist *types.Nodes)
	ReadEnodeWhitelist(openNetwork bool) *types.Nodes
}

type Contract struct {
	address                  common.Address
	contractABI              *abi.ABI
	bc                       Blockchainer
	SavedValidatorsRetriever func(i uint64) ([]common.Address, error)

	canTransfer func(db vm.StateDB, addr common.Address, amount *big.Int) bool
	transfer    func(db vm.StateDB, sender, recipient common.Address, amount *big.Int)
	GetHashFn   func(ref *types.Header, chain ChainContext) func(n uint64) common.Hash
	sync.RWMutex
}

func NewAutonityContract(
	bc Blockchainer,
	canTransfer func(db vm.StateDB, addr common.Address, amount *big.Int) bool,
	transfer func(db vm.StateDB, sender, recipient common.Address, amount *big.Int),
	GetHashFn func(ref *types.Header, chain ChainContext) func(n uint64) common.Hash,
) *Contract {
	return &Contract{
		bc:          bc,
		canTransfer: canTransfer,
		transfer:    transfer,
		GetHashFn:   GetHashFn,
	}
}

func (ac *Contract) ContractGetValidators(chain consensus.ChainReader, header *types.Header, statedb *state.StateDB) ([]common.Address, error) {
	if header.Number.Cmp(big.NewInt(1)) == 0 && ac.SavedValidatorsRetriever != nil {
		return ac.SavedValidatorsRetriever(1)
	}
	var addresses []common.Address
	err := ac.AutonityContractCall(chain, statedb, header, "getValidators", &addresses)
	if err != nil {
		return nil, err
	}
	sortableAddresses := common.Addresses(addresses)
	sort.Sort(sortableAddresses)
	return sortableAddresses, nil
}

func (ac *Contract) UpdateEnodesWhitelist(state *state.StateDB, block *types.Block) error {
	newWhitelist, err := ac.GetWhitelist(block, state)
	if err != nil {
		log.Error("could not call contract", "err", err)
		return ErrAutonityContract
	}

	ac.bc.UpdateEnodeWhitelist(newWhitelist)
	return nil
}

func (ac *Contract) GetWhitelist(block *types.Block, db *state.StateDB) (*types.Nodes, error) {
	var (
		newWhitelist *types.Nodes
		err          error
	)

	if block.Number().Uint64() == 1 {
		// use genesis block whitelist
		newWhitelist = ac.bc.ReadEnodeWhitelist(false)
	} else {
		// call retrieveWhitelist contract function
		newWhitelist, err = ac.callGetWhitelist(db, block.Header())
	}

	return newWhitelist, err
}

//blockchain

func (ac *Contract) GetMinimumGasPrice(block *types.Block, db *state.StateDB) (uint64, error) {
	if block.Number().Uint64() <= 1 {
		return ac.bc.Config().AutonityContractConfig.MinGasPrice, nil
	}

	return ac.callGetMinimumGasPrice(db, block.Header())
}

func (ac *Contract) ApplyFinalize(transactions types.Transactions, receipts types.Receipts, header *types.Header, statedb *state.StateDB) error {
	log.Info("ApplyFinalize", "header", header.Number.Uint64())
	if header.Number.Cmp(big.NewInt(1)) < 1 {
		return nil
	}
	blockGas := new(big.Int)
	for i, tx := range transactions {
		blockGas.Add(blockGas, new(big.Int).Mul(tx.GasPrice(), new(big.Int).SetUint64(receipts[i].GasUsed)))
	}
	log.Info("execution start ApplyFinalize", "balance", statedb.GetBalance(ac.Address()), "block", header.Number.Uint64(), "gas", blockGas.Uint64())

	if header.Number.Uint64() <= 1 {
		return nil
	}

	upgradeContract, err := ac.callFinalize(statedb, header, blockGas)
	if err != nil {
		return err
	}
	if upgradeContract {
		return ac.performContractUpgrade(statedb, header)
	}

	return err
}

func (ac *Contract) performContractUpgrade(statedb *state.StateDB, header *types.Header) error {
	log.Info("performing Autonity Contract upgrade", "header", header.Number.Uint64())
	state, errState := ac.callRetrieveState(statedb, header)
	if errState != nil {
		return errState
	}
	bytecode, abi, errContract := ac.callRetrieveContract(statedb, header)
	if errContract != nil {
		return errContract
	}
	//we save a snapshot of the statedb in case of upgrade failure
	snapshot := statedb.Snapshot()

	return nil
}

func (ac *Contract) Address() common.Address {
	if reflect.DeepEqual(ac.address, common.Address{}) {
		addr, err := ac.bc.Config().AutonityContractConfig.GetContractAddress()
		if err != nil {
			log.Error("Cant get contract address", "err", err)
		}
		return addr
	}
	return ac.address
}

func (ac *Contract) abi() (*abi.ABI, error) {
	ac.Lock()
	defer ac.Unlock()
	if ac.contractABI != nil {
		return ac.contractABI, nil
	}
	ABI, err := abi.JSON(strings.NewReader(ac.bc.Config().AutonityContractConfig.ABI))
	if err != nil {
		return nil, err
	}
	ac.contractABI = &ABI
	return ac.contractABI, nil

}
