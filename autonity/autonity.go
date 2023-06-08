package autonity

import (
	"github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/consensus/tendermint/core/constants"
	"math/big"
	"sort"
	"strings"
	"sync"

	"github.com/autonity/autonity/accounts/abi"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core/state"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/core/vm"
	"github.com/autonity/autonity/crypto"
	"github.com/autonity/autonity/log"
)

var Deployer = common.Address{}
var ContractAddress = crypto.CreateAddress(Deployer, 0)

const ABISPEC = "ABISPEC"

// EVMProvider provides a new evm. This allows us to decouple the contract from *params.ChainConfig which is required to build a new evm.
type EVMProvider interface {
	EVM(header *types.Header, origin common.Address, statedb *state.StateDB) *vm.EVM
}

type Blockchainer interface {
	PutKeyValue(key []byte, value []byte) error
}

type Contract struct {
	evmProvider       EVMProvider
	operator          common.Address
	initialMinBaseFee uint64
	contractABI       *abi.ABI
	stringABI         string
	bc                Blockchainer

	// map[Height]map[Round]proposer_address
	proposers map[uint64]map[int64]common.Address
	sync.RWMutex
}

func NewAutonityContract(
	bc Blockchainer,
	operator common.Address,
	minBaseFee uint64,
	ABI string,
	evmProvider EVMProvider,
) (*Contract, error) {
	contract := Contract{
		stringABI:         ABI,
		operator:          operator,
		initialMinBaseFee: minBaseFee,
		bc:                bc,
		evmProvider:       evmProvider,
		proposers:         make(map[uint64]map[int64]common.Address),
	}
	err := contract.upgradeAbiCache(ABI)
	return &contract, err
}

func (ac *Contract) GetCommittee(header *types.Header, statedb *state.StateDB) (types.Committee, error) {
	var committeeSet types.Committee
	err := ac.AutonityContractCall(statedb, header, "getCommittee", &committeeSet)
	if err != nil {
		return nil, err
	}
	sort.Sort(committeeSet)
	return committeeSet, err
}

func (ac *Contract) GetCommitteeEnodes(block *types.Block, db *state.StateDB) (*types.Nodes, error) {
	return ac.callGetCommitteeEnodes(db, block.Header())
}

func (ac *Contract) GetMinimumBaseFee(block *types.Header, db *state.StateDB) (uint64, error) {
	if block.Number.Uint64() <= 1 {
		return ac.initialMinBaseFee, nil
	}

	return ac.callGetMinimumBaseFee(db, block)
}

func (ac *Contract) GetProposer(header *types.Header, height uint64, round int64) common.Address {
	ac.Lock()
	defer ac.Unlock()
	roundMap, ok := ac.proposers[height]
	if !ok {
		p := ac.electProposer(header, height, round)
		roundMap = make(map[int64]common.Address)
		roundMap[round] = p
		ac.proposers[height] = roundMap
		// since the growing of blockchain height, so most of the case we can use the latest buffered height to gc the buffer.
		ac.proposerBufferGC(height)
		return p
	}

	proposer, ok := roundMap[round]
	if !ok {
		p := ac.electProposer(header, height, round)
		ac.proposers[height][round] = p
		return p
	}

	return proposer
}

// electProposer is a part of consensus, that it elect proposer from parent header's committee list which was returned
// from autonity contract stable ordered by voting power in evm context.
func (ac *Contract) electProposer(parentHeader *types.Header, height uint64, round int64) common.Address {
	h := new(big.Int).SetUint64(height)
	r := new(big.Int).SetInt64(round)
	seed := new(big.Int).SetInt64(constants.MaxRound)
	totalVotingPower := new(big.Int).SetUint64(0)
	for _, c := range parentHeader.Committee {
		totalVotingPower.Add(totalVotingPower, c.VotingPower)
	}

	// this shouldn't be happen, to keep the liveness, we elect with round-robin.
	if totalVotingPower.Cmp(common.Big0) == 0 {
		return parentHeader.Committee[round%int64(len(parentHeader.Committee))].Address
	}

	// for power weighted sampling, we distribute seed into a 256bits key-space, and compute the hit index.
	key := r.Add(r, h.Mul(h, seed))
	value := new(big.Int).SetBytes(crypto.Keccak256(key.Bytes()))
	index := value.Mod(value, totalVotingPower)

	// find the index hit which committee member which line up in the committee list.
	// we assume there is no 0 stake/power validators.
	counter := new(big.Int).SetUint64(0)
	for _, c := range parentHeader.Committee {
		counter.Add(counter, c.VotingPower)
		if index.Cmp(counter) == -1 {
			log.Info("Elect stake weighted proposer", "Proposer", c.Address, "height", height, "round", round)
			return c.Address
		}
	}

	// otherwise, we elect with round-robin.
	return parentHeader.Committee[round%int64(len(parentHeader.Committee))].Address
}

func (ac *Contract) proposerBufferGC(height uint64) {
	if len(ac.proposers) > consensus.AccountabilityHeightRange*2 && height > consensus.AccountabilityHeightRange*2 {
		gcFrom := height - consensus.AccountabilityHeightRange*2
		// keep at least consensus.AccountabilityHeightRange heights in the buffer.
		for i := uint64(0); i < consensus.AccountabilityHeightRange; i++ {
			gcHeight := gcFrom + i
			delete(ac.proposers, gcHeight)
		}
	}
}

func (ac *Contract) SetMinimumBaseFee(block *types.Block, db *state.StateDB, price *big.Int) error {
	if block.Number().Uint64() <= 1 {
		return nil
	}

	return ac.callSetMinimumBaseFee(db, block.Header(), price)
}

func (ac *Contract) FinalizeAndGetCommittee(header *types.Header, statedb *state.StateDB) (types.Committee, *types.Receipt, error) {
	if header.Number.Uint64() == 0 {
		return nil, nil, nil
	}

	log.Debug("Finalizing block",
		"balance", statedb.GetBalance(ContractAddress),
		"block", header.Number.Uint64())

	upgradeContract, committee, err := ac.callFinalize(statedb, header)
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

	log.Debug("Finalizing block", "contract upgrade", upgradeContract)

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

	// get contract binary and abi set by system operator before.
	bytecode, newAbi, errContract := ac.callRetrieveContract(statedb, header)
	if errContract != nil {
		return errContract
	}

	// take snapshot in case of roll back to former view.
	snapshot := statedb.Snapshot()

	if err := ac.updateAutonityContract(header, statedb, bytecode); err != nil {
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

// Address return the current autonity contract address
func (ac *Contract) Address() common.Address {
	return ContractAddress
}

func (ac *Contract) MisbehaviourProcessed(header *types.Header, db *state.StateDB, msgHash common.Hash) bool {
	return ac.callAccountabilityEventProcessed(db, header, "misbehaviourProcessed", msgHash)
}

func (ac *Contract) AccusationProcessed(header *types.Header, db *state.StateDB, msgHash common.Hash) bool {
	return ac.callAccountabilityEventProcessed(db, header, "accusationProcessed", msgHash)
}

func (ac *Contract) GetValidatorAccusations(header *types.Header, db *state.StateDB, validator common.Address) []AccountabilityEvent {
	return ac.callGetValidatorAccusations(db, header, validator)
}

func (ac *Contract) GetAccountabilityEventChunk(header *types.Header, db *state.StateDB, msgHash common.Hash, tp uint8, rule uint8, reporter common.Address, chunkID uint8) ([]byte, error) {
	return ac.callGetAccountabilityEventChunk(db, header, msgHash, tp, rule, reporter, chunkID)
}

func (ac *Contract) GetEpochPeriod(header *types.Header, db *state.StateDB) (uint64, error) {
	return ac.callGetEpochPeriod(db, header)
}

func (ac *Contract) GetLastEpochBlock(header *types.Header, db *state.StateDB) (uint64, error) {
	return ac.callGetLastEpochBlock(db, header)
}
