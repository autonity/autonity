package autonity

import (
	"errors"
	"math/big"
	"sort"
	"strings"
	"sync"

	"github.com/clearmatics/autonity/accounts/abi"
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus"
	"github.com/clearmatics/autonity/core/state"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/core/vm"
	"github.com/clearmatics/autonity/crypto"
	"github.com/clearmatics/autonity/log"
)

var ErrAutonityContract = errors.New("could not call Autonity contract")
var ErrWrongParameter = errors.New("wrong parameter")
var deployer = common.Address{}
var ContractAddress = crypto.CreateAddress(deployer, 0)

const ABISPEC = "ABISPEC"

// EVMProvider provides a new evm. This allows us to decouple the contract from *params.ChainConfig which is required to build a new evm.
type EVMProvider interface {
	EVM(header *types.Header, origin common.Address, statedb *state.StateDB) *vm.EVM
}

func NewAutonityContract(
	bc Blockchainer,
	operator common.Address,
	minGasPrice uint64,
	ABI string,
	evmProvider EVMProvider,
) (*Contract, error) {

	contractABI, err := abi.JSON(strings.NewReader(ABI))
	if err != nil {
		return nil, err
	}
	return &Contract{
		stringContractABI:  ABI,
		contractABI:        &contractABI,
		operator:           operator,
		initialMinGasPrice: minGasPrice,
		bc:                 bc,
		evmProvider:        evmProvider,
	}, nil
}

type Blockchainer interface {
	UpdateEnodeWhitelist(newWhitelist *types.Nodes)
	ReadEnodeWhitelist() *types.Nodes

	PutKeyValue(key []byte, value []byte) error
}

type Contract struct {
	evmProvider        EVMProvider
	operator           common.Address
	initialMinGasPrice uint64
	contractABI        *abi.ABI
	stringContractABI  string
	bc                 Blockchainer
	metrics            EconomicMetrics

	sync.RWMutex
}

// measure metrics of user's meta data by regarding of network economic.
func (ac *Contract) MeasureMetricsOfNetworkEconomic(header *types.Header, stateDB *state.StateDB) {
	if header == nil || stateDB == nil || header.Number.Uint64() < 1 {
		return
	}

	// prepare abi and evm context
	gas := uint64(0xFFFFFFFF)
	evm := ac.evmProvider.EVM(header, deployer, stateDB)

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
	ret, _, vmerr := evm.Call(vm.AccountRef(deployer), ContractAddress, input, gas, value)
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

	ac.metrics.SubmitEconomicMetrics(&v, stateDB, header.Number.Uint64(), ac.operator)
}

func (ac *Contract) GetCommittee(chain consensus.ChainReader, header *types.Header, statedb *state.StateDB) (types.Committee, error) {
	// The Autonity Contract is not deployed yet at block #1, the committee is supposed to remains the same as genesis.
	if header.Number.Cmp(big.NewInt(1)) == 0 {
		return chain.GetHeaderByNumber(0).Committee, nil
	}

	var committeeSet types.Committee
	err := ac.AutonityContractCall(statedb, header, "getCommittee", &committeeSet)
	if err != nil {
		return nil, err
	}
	sort.Sort(committeeSet)
	return committeeSet, err
}

func (ac *Contract) UpdateEnodesWhitelist(state *state.StateDB, block *types.Block) error {
	newWhitelist, err := ac.GetWhitelist(block, state)
	if err != nil {
		log.Error("Could not call contract", "err", err)
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
		newWhitelist = ac.bc.ReadEnodeWhitelist()
	} else {
		// call retrieveWhitelist contract function
		newWhitelist, err = ac.callGetWhitelist(db, block.Header())
	}

	return newWhitelist, err
}

func (ac *Contract) GetMinimumGasPrice(block *types.Block, db *state.StateDB) (uint64, error) {
	if block.Number().Uint64() <= 1 {
		return ac.initialMinGasPrice, nil
	}

	return ac.callGetMinimumGasPrice(db, block.Header())
}

func (ac *Contract) SetMinimumGasPrice(block *types.Block, db *state.StateDB, price *big.Int) error {
	if block.Number().Uint64() <= 1 {
		return nil
	}

	return ac.callSetMinimumGasPrice(db, block.Header(), price)
}

func (ac *Contract) FinalizeAndGetCommittee(transactions types.Transactions, receipts types.Receipts, header *types.Header, statedb *state.StateDB) (types.Committee, *types.Receipt, error) {
	if header.Number.Uint64() == 0 {
		return nil, nil, nil
	}
	blockGas := new(big.Int)
	for i, tx := range transactions {
		blockGas.Add(blockGas, new(big.Int).Mul(tx.GasPrice(), new(big.Int).SetUint64(receipts[i].GasUsed)))
	}

	log.Info("ApplyFinalize",
		"balance", statedb.GetBalance(ContractAddress),
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
			log.Warn("Autonity Contract Upgrade Failed", "err", err)
		}
	}
	return committee, receipt, nil
}

func (ac *Contract) performContractUpgrade(statedb *state.StateDB, header *types.Header) error {
	log.Error("Initiating Autonity Contract upgrade", "header", header.Number.Uint64())

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
	statedb.CreateAccount(ContractAddress)

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

func (ac *Contract) abi() (*abi.ABI, error) {
	return ac.contractABI, nil
}

func (ac *Contract) upgradeAbiCache(newAbi string) error {
	ac.Lock()
	defer ac.Unlock()
	newABI, err := abi.JSON(strings.NewReader(newAbi))
	if err != nil {
		return err
	}

	ac.contractABI = &newABI
	return nil
}

func (ac *Contract) GetContractABI() string {
	return ac.stringContractABI
}
