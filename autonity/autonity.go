package autonity

import (
	"bytes"
	"errors"
	"math/big"
	"strings"
	"sync"

	"github.com/autonity/autonity/accounts/abi"
	"github.com/autonity/autonity/accounts/abi/bind"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/consensus/tendermint/core/constants"
	"github.com/autonity/autonity/core/rawdb"
	"github.com/autonity/autonity/core/state"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/core/vm"
	"github.com/autonity/autonity/crypto"
	"github.com/autonity/autonity/ethdb"
	"github.com/autonity/autonity/log"
	"github.com/autonity/autonity/params"
)

var (
	DeployerAddress               = common.Address{}
	AutonityContractAddress       = crypto.CreateAddress(DeployerAddress, 0)
	AccountabilityContractAddress = crypto.CreateAddress(DeployerAddress, 1)
	OracleContractAddress         = crypto.CreateAddress(DeployerAddress, 2)
	ACUContractAddress            = crypto.CreateAddress(DeployerAddress, 3)
	SupplyControlContractAddress  = crypto.CreateAddress(DeployerAddress, 4)
	StabilizationContractAddress  = crypto.CreateAddress(DeployerAddress, 5)
	AutonityABIKey                = []byte("ABISPEC")
)

// Errors
var (
	ErrAutonityContract = errors.New("could not call Autonity contract")
	ErrWrongParameter   = errors.New("wrong parameter")
	ErrNoAutonityConfig = errors.New("autonity section missing from genesis")
)

// EVMProvider provides a new evm. This allows us to decouple the contract from *core.blockchain which is required to build a new evm.
type EVMProvider func(header *types.Header, origin common.Address, statedb *state.StateDB) *vm.EVM

type Contracts struct {
	evmProvider     EVMProvider
	contractABI     *abi.ABI
	db              ethdb.Database
	chainConfig     *params.ChainConfig
	contractBackend bind.ContractBackend
	// map[Height]map[Round]proposer_address
	proposers map[uint64]map[int64]common.Address

	*Accountability
	sync.RWMutex
}

func New(config *params.ChainConfig, db ethdb.Database, provider EVMProvider, contractBackend bind.ContractBackend) (*Contracts, error) {
	if config.AutonityContractConfig == nil {
		return nil, ErrNoAutonityConfig
	}
	ABI := config.AutonityContractConfig.ABI
	// Retrieve a potentially updated ABI from the database
	if raw, err := rawdb.GetKeyValue(db, AutonityABIKey); err == nil || raw != nil {
		jsonABI, err := abi.JSON(bytes.NewReader(raw))
		if err != nil {
			return nil, err
		}
		ABI = &jsonABI
	}
	accountabilityContract, _ := NewAccountability(AccountabilityContractAddress, contractBackend)

	contract := Contracts{
		evmProvider:     provider,
		contractBackend: contractBackend,
		chainConfig:     config,
		contractABI:     ABI,
		db:              db,
		proposers:       make(map[uint64]map[int64]common.Address),
		Accountability:  accountabilityContract,
	}
	return &contract, nil
}

func (c *Contracts) CommitteeEnodes(block *types.Block, db *state.StateDB) (*types.Nodes, error) {
	return c.callGetCommitteeEnodes(db, block.Header())
}

func (c *Contracts) MinimumBaseFee(block *types.Header, db *state.StateDB) (uint64, error) {
	if block.Number.Uint64() <= 1 {
		return c.chainConfig.AutonityContractConfig.MinBaseFee, nil
	}
	return c.callGetMinimumBaseFee(db, block)
}

func (c *Contracts) Proposer(header *types.Header, _ *state.StateDB, height uint64, round int64) common.Address {
	c.Lock()
	defer c.Unlock()
	roundMap, ok := c.proposers[height]
	if !ok {
		p := c.electProposer(header, height, round)
		roundMap = make(map[int64]common.Address)
		roundMap[round] = p
		c.proposers[height] = roundMap
		// since the growing of blockchain height, so most of the case we can use the latest buffered height to gc the buffer.
		c.trimProposerCache(height)
		return p
	}

	proposer, ok := roundMap[round]
	if !ok {
		p := c.electProposer(header, height, round)
		c.proposers[height][round] = p
		return p
	}

	return proposer
}

// electProposer is a part of consensus, that it elect proposer from parent header's committee list which was returned
// from autonity contract stable ordered by voting power in evm context.
func (c *Contracts) electProposer(parentHeader *types.Header, height uint64, round int64) common.Address {
	seed := big.NewInt(constants.MaxRound)
	totalVotingPower := big.NewInt(0)
	for _, c := range parentHeader.Committee {
		totalVotingPower.Add(totalVotingPower, c.VotingPower)
	}

	// for power weighted sampling, we distribute seed into a 256bits key-space, and compute the hit index.
	h := new(big.Int).SetUint64(height)
	r := new(big.Int).SetInt64(round)
	key := r.Add(r, h.Mul(h, seed))
	value := new(big.Int).SetBytes(crypto.Keccak256(key.Bytes()))
	index := value.Mod(value, totalVotingPower)

	// find the index hit which committee member which line up in the committee list.
	// we assume there is no 0 stake/power validators.
	counter := new(big.Int).SetUint64(0)
	for _, c := range parentHeader.Committee {
		counter.Add(counter, c.VotingPower)
		if index.Cmp(counter) == -1 {
			return c.Address
		}
	}

	// otherwise, we elect with round-robin.
	return parentHeader.Committee[round%int64(len(parentHeader.Committee))].Address
}

func (c *Contracts) trimProposerCache(height uint64) {
	if len(c.proposers) > consensus.AccountabilityHeightRange*2 && height > consensus.AccountabilityHeightRange*2 {
		gcFrom := height - consensus.AccountabilityHeightRange*2
		// keep at least consensus.AccountabilityHeightRange heights in the buffer.
		for i := uint64(0); i < consensus.AccountabilityHeightRange; i++ {
			gcHeight := gcFrom + i
			delete(c.proposers, gcHeight)
		}
	}
}

func (c *Contracts) FinalizeAndGetCommittee(header *types.Header, statedb *state.StateDB) (types.Committee, *types.Receipt, error) {
	if header.Number.Uint64() == 0 {
		return nil, nil, nil
	}

	log.Debug("Finalizing block",
		"balance", statedb.GetBalance(AutonityContractAddress),
		"block", header.Number.Uint64())

	upgradeContract, committee, err := c.callFinalize(statedb, header)
	if err != nil {
		return nil, nil, err
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
	return committee, receipt, nil
}

func (c *Contracts) upgradeAutonityContract(statedb *state.StateDB, header *types.Header) error {
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
	if err := rawdb.PutKeyValue(c.db, AutonityABIKey, []byte(newAbi)); err != nil {
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

func (c *Contracts) upgradeAbiCache(newAbi string) error {
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
func (c *Contracts) ABI() *abi.ABI {
	return c.contractABI
}
