package autonity

import (
	"errors"
	"math/big"
	"reflect"
	"sort"
	"strings"
	"sync"

	"github.com/clearmatics/autonity/accounts/abi"
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/core/state"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/core/vm"
	"github.com/clearmatics/autonity/log"
	"github.com/clearmatics/autonity/params"
)

var ErrAutonityContract = errors.New("could not call Autonity contract")
var ErrWrongParameter = errors.New("wrong parameter")

const ABISPEC = "ABISPEC"

type ChainContext interface {
	// GetHeader returns the hash corresponding to their hash.
	GetHeader(common.Hash, uint64) *types.Header
}

type Blockchainer interface {
	ChainContext
	GetVMConfig() *vm.Config
	Config() *params.ChainConfig

	UpdateEnodeWhitelist(newWhitelist *types.Nodes)
	ReadEnodeWhitelist() *types.Nodes

	PutKeyValue(key []byte, value []byte) error
	GetKeyValue(key []byte) ([]byte, error)
}

type Contract interface {
	// GetCommittee returns the current block consensus committee.
	GetCommittee(header *types.Header, statedb *state.StateDB) (types.Committee, error)

	// GetMinimumGasPrice returns the current block minimum gas price.
	GetMinimumGasPrice(block *types.Block, db *state.StateDB) (uint64, error)

	// FinalizeAndGetCommittee calls the contract's finalize function that normally perform reward redistribution.
	// This method returns the next block committee to avoid a further EVM call.
	FinalizeAndGetCommittee(txs types.Transactions, r types.Receipts, h *types.Header, db *state.StateDB) (types.Committee, *types.Receipt, error)

	MeasureMetricsOfNetworkEconomic(header *types.Header, stateDB *state.StateDB)

	UpdateEnodesWhitelist(state *state.StateDB, block *types.Block) error

	GetContractABI() string

	Address() common.Address

	DeployAutonityContract(chainConfig *params.ChainConfig, header *types.Header, statedb *state.StateDB) (common.Address, error)

	GetWhitelist(block *types.Block, db *state.StateDB) (*types.Nodes, error)
}

type evmContract struct {
	address     common.Address
	contractABI *abi.ABI
	bc          Blockchainer
	metrics     EconomicMetrics

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
) *evmContract {
	return &evmContract{
		bc:          bc,
		canTransfer: canTransfer,
		transfer:    transfer,
		GetHashFn:   GetHashFn,
	}
}

// measure metrics of user's meta data by regarding of network economic.
func (ac *evmContract) MeasureMetricsOfNetworkEconomic(header *types.Header, stateDB *state.StateDB) {
	if header == nil || stateDB == nil || header.Number.Uint64() < 1 {
		return
	}

	// prepare abi and evm context
	deployer := ac.bc.Config().AutonityContractConfig.Deployer
	sender := vm.AccountRef(deployer)
	gas := uint64(0xFFFFFFFF)
	evm := ac.getEVM(header, deployer, stateDB)

	ABI, err := ac.abi()
	if err != nil {
		return
	}

	// pack the function which dump the data from contract.
	input, err := ABI.Pack("dumpEconomicsMetricData")
	if err != nil {
		log.Warn("Cannot pack the method: ", "err", err.Error())
		return
	}

	// call evm.
	value := new(big.Int).SetUint64(0x00)
	ret, _, vmerr := evm.Call(sender, ac.Address(), input, gas, value)
	if vmerr != nil {
		log.Warn("Error Autonity Contract dumpNetworkEconomics", err, vmerr)
		return
	}

	// marshal the data from bytes arrays into specified structure.
	v := EconomicMetaData{make([]common.Address, 32), make([]uint8, 32), make([]*big.Int, 32),
		make([]*big.Int, 32), new(big.Int), new(big.Int)}

	if err := ABI.Unpack(&v, "dumpEconomicsMetricData", ret); err != nil {
		// can't work with aliased types
		log.Warn("Could not unpack dumpNetworkEconomicsData returned value",
			"err", err,
			"header.num", header.Number.Uint64())
		return
	}

	ac.metrics.SubmitEconomicMetrics(&v, stateDB, header.Number.Uint64(), ac.bc.Config().AutonityContractConfig.Operator)
}

func (ac *evmContract) GetCommittee(header *types.Header, statedb *state.StateDB) (types.Committee, error) {
	// The Autonity Contract is not deployed yet at block #1, we return an error if this
	// function is called at this height. In a past version we were returning the genesis committee field
	// but this was at the cost of having a parameter causing circular imports.
	if header.Number.Uint64() <= 1 {
		return nil, errors.New("calling GetCommittee for block #1 or #0")
	}

	var committeeSet types.Committee
	err := ac.AutonityContractCall(statedb, header, "getCommittee", &committeeSet)
	if err != nil {
		return nil, err
	}
	sort.Sort(committeeSet)
	return committeeSet, err
}

func (ac *evmContract) UpdateEnodesWhitelist(state *state.StateDB, block *types.Block) error {
	newWhitelist, err := ac.GetWhitelist(block, state)
	if err != nil {
		log.Error("Could not call contract", "err", err)
		return ErrAutonityContract
	}

	ac.bc.UpdateEnodeWhitelist(newWhitelist)
	return nil
}

func (ac *evmContract) GetWhitelist(block *types.Block, db *state.StateDB) (*types.Nodes, error) {
	var (
		newWhitelist *types.Nodes
		err          error
	)

	if block.Number().Uint64() == 1 {
		// use genesis block whitelist
		newWhitelist = ac.bc.ReadEnodeWhitelist()
	} else {
		// call retrieveWhitelist contract function
		newWhitelist, err = ac.callGetWhitelist(db, block.Header())
	}

	return newWhitelist, err
}

func (ac *evmContract) GetMinimumGasPrice(block *types.Block, db *state.StateDB) (uint64, error) {
	if block.Number().Uint64() <= 1 {
		return ac.bc.Config().AutonityContractConfig.MinGasPrice, nil
	}

	return ac.callGetMinimumGasPrice(db, block.Header())
}

func (ac *evmContract) SetMinimumGasPrice(block *types.Block, db *state.StateDB, price *big.Int) error {
	if block.Number().Uint64() <= 1 {
		return nil
	}

	return ac.callSetMinimumGasPrice(db, block.Header(), price)
}

func (ac *evmContract) FinalizeAndGetCommittee(transactions types.Transactions, receipts types.Receipts, header *types.Header, statedb *state.StateDB) (types.Committee, *types.Receipt, error) {
	if header.Number.Uint64() == 0 {
		return nil, nil, nil
	}
	blockGas := new(big.Int)
	for i, tx := range transactions {
		blockGas.Add(blockGas, new(big.Int).Mul(tx.GasPrice(), new(big.Int).SetUint64(receipts[i].GasUsed)))
	}

	log.Info("ApplyFinalize",
		"balance", statedb.GetBalance(ac.Address()),
		"block", header.Number.Uint64(),
		"gas", blockGas.Uint64())

	upgradeContract, committee, err := ac.callFinalize(statedb, header, blockGas)
	if err != nil {
		return nil, nil, err
	}

	// Create a new receipt for the finalize call
	receipt := types.NewReceipt(nil, false, 0)
	receipt.TxHash = common.ACHash(header.Number)
	receipt.GasUsed = 0
	receipt.Logs = statedb.GetLogs(receipt.TxHash)
	receipt.Bloom = types.CreateBloom(types.Receipts{receipt})
	receipt.BlockHash = statedb.BlockHash()
	receipt.BlockNumber = header.Number
	receipt.TransactionIndex = uint(statedb.TxIndex())

	log.Info("ApplyFinalize", "upgradeContract", upgradeContract)

	if upgradeContract {
		// warning prints for failure rather than returning error to stuck engine.
		// in any failure, the state will be rollback to snapshot.
		err = ac.performContractUpgrade(statedb, header)
		if err != nil {
			log.Warn("Autonity evmContract Upgrade Failed", "err", err)
		}
	}
	return committee, receipt, nil
}

func (ac *evmContract) performContractUpgrade(statedb *state.StateDB, header *types.Header) error {
	log.Error("Initiating Autonity evmContract upgrade", "header", header.Number.Uint64())

	// dump contract stateBefore first.
	stateBefore, errState := ac.callRetrieveState(statedb, header)
	if errState != nil {
		return errState
	}

	// get contract binary and abi set by system operator before.
	bytecode, newAbi, errContract := ac.callRetrieveContract(statedb, header)
	if errContract != nil {
		return errContract
	}

	// take snapshot in case of roll back to former view.
	snapshot := statedb.Snapshot()

	// Create account will delete previous the AC stateobject and carry over the balance
	statedb.CreateAccount(ac.Address())

	if err := ac.UpdateAutonityContract(header, statedb, bytecode, newAbi, stateBefore); err != nil {
		statedb.RevertToSnapshot(snapshot)
		return err
	}

	// save new abi in persistent, once node reset, it load from persistent level db.
	if err := ac.bc.PutKeyValue([]byte(ABISPEC), []byte(newAbi)); err != nil {
		statedb.RevertToSnapshot(snapshot)
		return err
	}

	// upgrade ac.ContractStateStore too right after the contract upgrade successfully.
	if err := ac.upgradeAbiCache(newAbi); err != nil {
		statedb.RevertToSnapshot(snapshot)
		return err
	}
	log.Info("Autonity Contract upgrade success")
	return nil
}

func (ac *evmContract) Address() common.Address {
	if reflect.DeepEqual(ac.address, common.Address{}) {
		addr, err := ac.bc.Config().AutonityContractConfig.GetContractAddress()
		if err != nil {
			log.Error("Cant get contract address", "err", err)
		}
		return addr
	}
	return ac.address
}

func (ac *evmContract) abi() (*abi.ABI, error) {
	ac.Lock()
	defer ac.Unlock()
	if ac.contractABI != nil {
		return ac.contractABI, nil
	}
	var JSONString = ac.bc.Config().AutonityContractConfig.ABI

	bytes, err := ac.bc.GetKeyValue([]byte(ABISPEC))
	if err == nil || bytes != nil {
		JSONString = string(bytes)
	}

	ABI, err := abi.JSON(strings.NewReader(JSONString))
	if err != nil {
		return nil, err
	}
	ac.contractABI = &ABI
	return ac.contractABI, nil
}

func (ac *evmContract) upgradeAbiCache(newAbi string) error {
	ac.Lock()
	defer ac.Unlock()
	newABI, err := abi.JSON(strings.NewReader(newAbi))
	if err != nil {
		return err
	}

	ac.contractABI = &newABI
	return nil
}

func (ac *evmContract) GetContractABI() string {
	ac.Lock()
	defer ac.Unlock()

	var JSONString = ac.bc.Config().AutonityContractConfig.ABI
	bytes, err := ac.bc.GetKeyValue([]byte(ABISPEC))
	if err == nil || bytes != nil {
		JSONString = string(bytes)
	}

	if err != nil {
		log.Warn("can't get the contract ABI", "err", err)
	}

	return JSONString
}
