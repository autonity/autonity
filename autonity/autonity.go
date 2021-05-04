package autonity

import (
	"errors"
	"math/big"
	"sort"
	"strings"
	"sync"

	"github.com/clearmatics/autonity/accounts/abi"
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/core/state"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/core/vm"
	"github.com/clearmatics/autonity/crypto"
	"github.com/clearmatics/autonity/log"
)

var ErrAutonityContract = errors.New("could not call Autonity contract")
var ErrWrongParameter = errors.New("wrong parameter")
var Deployer = common.Address{}
var ContractAddress = crypto.CreateAddress(Deployer, 0)

const ABISPEC = "ABISPEC"

// EVMProvider provides a new evm. This allows us to decouple the contract from *params.ChainConfig which is required to build a new evm.
type EVMProvider interface {
	EVM(header *types.Header, origin common.Address, statedb *state.StateDB) *vm.EVM
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
	stringABI          string
	bc                 Blockchainer
	metrics            EconomicMetrics

	sync.RWMutex
}

func NewAutonityContract(
	bc Blockchainer,
	operator common.Address,
	minGasPrice uint64,
	ABI string,
	evmProvider EVMProvider,
) (*Contract, error) {
	contract := Contract{
		stringABI:          ABI,
		operator:           operator,
		initialMinGasPrice: minGasPrice,
		bc:                 bc,
		evmProvider:        evmProvider,
	}
	err := contract.upgradeAbiCache(ABI)
	return &contract, err
}

func (ac *Contract) GetCommittee(header *types.Header, statedb *state.StateDB) (types.Committee, error) {
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

func (ac *Contract) GetProposerFromAC(header *types.Header, db *state.StateDB, height uint64, round int64) common.Address {
	return ac.callGetProposer(db, header, height, round)
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

	log.Debug("ApplyFinalize",
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

	log.Debug("ApplyFinalize", "upgradeContract", upgradeContract)

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
	log.Warn("Initiating Autonity Contract upgrade", "header", header.Number.Uint64())

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

	if err := ac.updateAutonityContract(header, statedb, bytecode, stateBefore); err != nil {
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

func (ac *Contract) upgradeAbiCache(newAbi string) error {
	ac.Lock()
	defer ac.Unlock()
	newABI, err := abi.JSON(strings.NewReader(newAbi))
	if err != nil {
		return err
	}

	ac.contractABI = &newABI
	ac.stringABI = newAbi
	return nil
}

// StringABI returns the current autonity contract ABI in string format
func (ac *Contract) StringABI() string {
	return ac.stringABI
}

// ABI returns the current autonity contract's ABI
func (ac *Contract) ABI() *abi.ABI {
	return ac.contractABI
}
