package autonity

import (
	"bytes"
	"errors"
	"math/big"
	"strings"
	"sync"

	"github.com/autonity/autonity/accounts/abi"
	"github.com/autonity/autonity/accounts/abi/bind"
	"github.com/autonity/autonity/autonity/bindings"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/common/math"
	"github.com/autonity/autonity/core/rawdb"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/core/vm"
	"github.com/autonity/autonity/ethdb"
	"github.com/autonity/autonity/log"
	"github.com/autonity/autonity/params"
)

// "soft minimum" cache size for the proposer cache.
// "soft minimum" because we can have less entries in it, but once we reach `thresholdSize` we want to keep having at least `thresholdSize`.
// ensure it is >= accountability.accountabilityHeightRange to ensure reduced performance overhead in case of old proposal spamming
const thresholdSize = 256

var (
	AbiKey = []byte("ABISPEC")
)

// Errors
var (
	ErrNoAutonityConfig = errors.New("autonity section missing from genesis")
)

// EVMProvider provides a new evm. This allows us to decouple the contract from *core.blockchain which is required to build a new evm.
type EVMProvider func(header *types.Header, origin common.Address, statedb vm.StateDB) *vm.EVM

type evmContract struct {
	evmProvider EVMProvider
	contractABI *abi.ABI
	db          ethdb.Database
	chainConfig *params.ChainConfig

	sync.RWMutex
}

func NewEVMContract(evmProvider EVMProvider, contractABI *abi.ABI, db ethdb.Database, chainConfig *params.ChainConfig) *evmContract {
	return &evmContract{
		evmProvider: evmProvider,
		contractABI: contractABI,
		db:          db,
		chainConfig: chainConfig,
	}
}

func (c *evmContract) upgradeAbiCache(newAbi string) error {
	c.Lock()
	defer c.Unlock()
	newABI, err := abi.JSON(strings.NewReader(newAbi))
	if err != nil {
		return err
	}
	c.contractABI = &newABI
	return nil
}

// ABI returns the current autonity contract's ABI
func (c *evmContract) ABI() *abi.ABI {
	return c.contractABI
}

// callContractFunc creates an evm object, uses it to call the
// specified function of the contract at contractAddress with packedArgs and returns the
// packed result. If there is an error making the evm call it will be returned.
// Callers should use the contract ABI to pack and unpack the args and result.
// It returns the amount of gas used for the call
func (c *evmContract) callContractFunc(statedb vm.StateDB, header *types.Header, contractAddress common.Address, packedArgs []byte) ([]byte, uint64, error) {
	gas := uint64(math.MaxUint64)
	evm := c.evmProvider(header, params.DeployerAddress, statedb)
	packedResult, leftOverGas, err := evm.Call(vm.AccountRef(params.DeployerAddress), contractAddress, packedArgs, gas, new(big.Int))
	usedGas := gas - leftOverGas
	return packedResult, usedGas, err
}

func (c *evmContract) callContractFuncAs(statedb vm.StateDB, header *types.Header, contractAddress common.Address, origin common.Address, packedArgs []byte) ([]byte, error) {
	gas := uint64(math.MaxUint64)
	evm := c.evmProvider(header, origin, statedb)
	packedResult, _, err := evm.Call(vm.AccountRef(origin), contractAddress, packedArgs, gas, new(big.Int))
	return packedResult, err
}

//revive:disable:exported - Autonity is one of the contracts, so repetitive naming here is justified
type AutonityContract struct {
	evmContract
	*bindings.AutonityFilterer                                     // allows to watch for Autonity Contract events
	proposers                  map[uint64]map[int64]common.Address // map[height][round] --> proposer
}

type ProtocolContracts struct {
	*AutonityContract
	*bindings.Accountability
	*bindings.UpgradeManager1
}

func NewProtocolContracts(
	config *params.ChainConfig,
	db ethdb.Database,
	provider EVMProvider,
	contractBackend bind.ContractBackend,
) (*ProtocolContracts, error) {
	if config.AutonityContractConfig == nil {
		return nil, ErrNoAutonityConfig
	}
	ABI := config.AutonityContractConfig.ABI
	// Retrieve a potentially updated ABI from the database
	if raw, err := rawdb.GetKeyValue(db, AbiKey); err == nil || raw != nil {
		jsonABI, err := abi.JSON(bytes.NewReader(raw))
		if err != nil {
			return nil, err
		}
		ABI = &jsonABI
	}

	// create autonity EVM contract
	autonityFilterer, err := bindings.NewAutonityFilterer(params.AutonityContractAddress, contractBackend.(bind.ContractFilterer))
	if err != nil {
		return nil, err
	}
	autonityContract := &AutonityContract{
		evmContract: evmContract{
			evmProvider: provider,
			contractABI: ABI,
			db:          db,
			chainConfig: config,
		},
		AutonityFilterer: autonityFilterer,
		proposers:        make(map[uint64]map[int64]common.Address),
	}

	// bind to accountability contract
	accountabilityContract, err := bindings.NewAccountability(params.AccountabilityContractAddress, contractBackend)
	if err != nil {
		return nil, err
	}

	// bind to upgrade manager contract
	upgradeManagerContract, err := bindings.NewUpgradeManager1(params.UpgradeManagerContractAddress, contractBackend)
	if err != nil {
		return nil, err
	}

	contract := ProtocolContracts{
		AutonityContract: autonityContract,
		Accountability:   accountabilityContract,
		UpgradeManager1:  upgradeManagerContract,
	}

	return &contract, nil
}

// Proposer election is now computed by committee structure, it is on longer depends on AC contract.
func (c *AutonityContract) Proposer(committee *types.Committee, _ vm.StateDB, height uint64, round int64) (proposer common.Address) {
	c.Lock()
	defer c.Unlock()

	needsTrimming := false

	_, ok := c.proposers[height]
	if !ok {
		c.proposers[height] = make(map[int64]common.Address)
		needsTrimming = true // cannot trim here, we might delete the entry we just created if it is for a very old height
	}

	proposer, ok = c.proposers[height][round]
	if !ok {
		proposer = committee.Proposer(height, round)
		c.proposers[height][round] = proposer
	}

	if needsTrimming {
		// since the chain monotonically grows towards higher heights, we can use the latest buffered height to trim the cache
		// NOTE: this is not 100% accurate because we could be computing Proposer for a very old height. But it works anyways.
		c.trimProposerCache(height)
	}

	return proposer
}

/* the Proposer election function is called from core and from the fault detector.
*
* Core calls it only on current height messages, therefore we don't need to retain old (h,r) proposers for it
* note that it is still good to have a cache for current height, because core will call multiple times the Proposer function for the same (h,r)
*
* The fault detector calls it when checking a proposal and when validating an InvalidProposer fault. For the proposal checking, we will accepts messages that
* are in [currentHeight - accountability.accountabilityHeightRange, currentHeight]. However, I think that it will be rare to process old messages, unless there is
* a malicious peer which is spamming old height messages.
* For proof validation instead, there is no height boundary (could be a proof for a very old height).
 */
func (c *AutonityContract) trimProposerCache(height uint64) {
	// we trim the cache only when the size is double `thresholdSize`
	if height <= thresholdSize*2 || len(c.proposers) <= thresholdSize*2 {
		return
	}

	for h := range c.proposers {
		if h < height-thresholdSize {
			delete(c.proposers, h)
		}
	}
}

func (c *AutonityContract) FinalizeAndGetCommittee(header *types.Header, statedb vm.StateDB) (*types.Receipt, *types.Epoch, *types.ContractsConfig, error) {
	log.Debug(
		"Finalizing block",
		"balance",
		statedb.GetBalance(params.AutonityContractAddress),
		"block",
		header.Number.Uint64(),
	)
	upgradeContract, epochInfo, contractsConfig, err := c.callFinalize(statedb, header)
	if err != nil {
		return nil, nil, nil, err
	}

	// Create a new receipt for the finalize call
	receipt := types.NewReceipt(nil, false, 0)
	receipt.TxHash = common.ACHash(header.Number)
	receipt.GasUsed = 0
	receipt.Logs = statedb.GetLogs(receipt.TxHash, header.Hash())
	receipt.Bloom = types.CreateBloom(types.Receipts{receipt})
	receipt.BlockHash = header.Hash()
	receipt.BlockNumber = header.Number
	receipt.TransactionIndex = uint(statedb.TxIndex())

	if upgradeContract {
		// warning prints for failure rather than returning error to stuck engine.
		// in any failure, the state will be rollback to snapshot.
		if err = c.upgradeAutonityContract(statedb, header); err != nil {
			log.Warn("Autonity Contracts Upgrade Failed", "err", err)
		}
	}
	return receipt, epochInfo, contractsConfig, nil
}

func (c *AutonityContract) upgradeAutonityContract(statedb vm.StateDB, header *types.Header) error {
	log.Info("Initiating Autonity Contracts upgrade", "header", header.Number.Uint64())

	// get contract binary and abi set by system operator before.
	bytecode, newAbi, errContract := c.callRetrieveContract(statedb, header)
	if errContract != nil {
		return errContract
	}

	// take snapshot in case of roll back to former view.
	snapshot := statedb.Snapshot()
	if err := c.replaceAutonityBytecode(header, statedb, bytecode); err != nil {
		statedb.RevertToSnapshot(snapshot)
		return err
	}

	// save new abi in persistent, once node reset, it load from persistent level db.
	if err := rawdb.PutKeyValue(c.db, AbiKey, []byte(newAbi)); err != nil {
		statedb.RevertToSnapshot(snapshot)
		return err
	}

	// upgrade c.ContractStateStore too right after the contract upgrade successfully.
	if err := c.upgradeAbiCache(newAbi); err != nil {
		statedb.RevertToSnapshot(snapshot)
		return err
	}
	log.Info("Autonity Contract upgrade success")
	return nil
}
