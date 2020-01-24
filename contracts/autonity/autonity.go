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
	"github.com/clearmatics/autonity/consensus"
	"github.com/clearmatics/autonity/core/state"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/core/vm"
	"github.com/clearmatics/autonity/log"
	"github.com/clearmatics/autonity/params"
)

var ErrAutonityContract = errors.New("could not call Autonity contract")
var ErrWrongParameter = errors.New("wrong parameter")

const ABISPEC = "ABISPEC"

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

	PutKeyValue(key []byte, value []byte) error
	GetKeyValue(key []byte) ([]byte, error)
}

type Contract struct {
	address                 common.Address
	contractABI             *abi.ABI
	bc                      Blockchainer
	SavedCommitteeRetriever func(i uint64) (types.Committee, error)
	metrics                 EconomicMetrics

	canTransfer func(db vm.StateDB, addr common.Address, amount *big.Int) bool
	transfer    func(db vm.StateDB, sender, recipient common.Address, amount *big.Int)
	GetHashFn   func(ref *types.Header, chain ChainContext) func(n uint64) common.Hash
	sync.RWMutex
}

// measure metrics of user's meta data by regarding of network economic.
func (ac *Contract) MeasureMetricsOfNetworkEconomic(header *types.Header, stateDB *state.StateDB) {
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
		log.Warn("Cannot pack the method: ", err.Error())
		return
	}

	// call evm.
	value := new(big.Int).SetUint64(0x00)
	ret, _, vmerr := evm.Call(sender, ac.Address(), input, gas, value)
	log.Debug("bytes return from contract: ", ret)
	if vmerr != nil {
		log.Warn("Error Autonity Contract dumpNetworkEconomics")
		return
	}

	// marshal the data from bytes arrays into specified structure.
	v := EconomicMetaData{make([]common.Address, 32), make([]uint8, 32), make([]*big.Int, 32),
		make([]*big.Int, 32), new(big.Int), new(big.Int)}

	if err := ABI.Unpack(&v, "dumpEconomicsMetricData", ret); err != nil { // can't work with aliased types
		log.Warn("Could not unpack dumpNetworkEconomicsData returned value", "err", err, "header.num",
			header.Number.Uint64())
		return
	}

	ac.metrics.SubmitEconomicMetrics(&v, stateDB, header.Number.Uint64(), ac.bc.Config().AutonityContractConfig.Operator)
}

/*
//// Instantiates a new EVM object which is required when creating or calling a deployed contract
func (ac *Contract) getEVM(header *types.Header, origin common.Address, statedb *state.StateDB) *vm.EVM {
	coinbase, _ := types.Ecrecover(header)
	evmContext := vm.Context{
		CanTransfer: ac.canTransfer,
		Transfer:    ac.transfer,
		GetHash:     ac.GetHashFn(header, ac.bc),
		Origin:      origin,
		Coinbase:    coinbase,
		BlockNumber: header.Number,
		Time:        new(big.Int).SetUint64(header.Time),
		GasLimit:    header.GasLimit,
		Difficulty:  header.Difficulty,
		GasPrice:    new(big.Int).SetUint64(0x0),
	}
	vmConfig := *ac.bc.GetVMConfig()
	evm := vm.NewEVM(evmContext, statedb, ac.bc.Config(), vmConfig)
	return evm
}

// deployContract deploys the contract contained within the genesis field bytecode
func (ac *Contract) DeployAutonityContract(chain consensus.ChainReader, header *types.Header, statedb *state.StateDB) (common.Address, error) {
	// Convert the contract bytecode from hex into bytes
	contractBytecode := common.Hex2Bytes(chain.Config().AutonityContractConfig.Bytecode)
	evm := ac.getEVM(header, chain.Config().AutonityContractConfig.Deployer, statedb)
	sender := vm.AccountRef(chain.Config().AutonityContractConfig.Deployer)
	log.Error("Operator:", "addr", chain.Config().AutonityContractConfig.Operator.String())

	contractABI, err := ac.abi()
	if err != nil {
		log.Error("abi.JSON returns err", "err", err)
		return common.Address{}, err
	}

	ln := len(chain.Config().AutonityContractConfig.Users)
	users := make(common.Addresses, 0, ln)
	enodes := make([]string, 0, ln)
	accTypes := make([]*big.Int, 0, ln)
	participantStake := make([]*big.Int, 0, ln)
	for _, v := range chain.Config().AutonityContractConfig.Users {
		users = append(users, v.Address)
		enodes = append(enodes, v.Enode)
		accTypes = append(accTypes, big.NewInt(int64(v.Type.GetID())))
		participantStake = append(participantStake, big.NewInt(int64(v.Stake)))
	}

	//"" means contructor
	constructorParams, err := contractABI.Pack("",
		users,
		enodes,
		accTypes,
		participantStake,
		chain.Config().AutonityContractConfig.Operator,
		new(big.Int).SetUint64(chain.Config().AutonityContractConfig.MinGasPrice))
	if err != nil {
		log.Error("contractABI.Pack returns err", "err", err)
		return common.Address{}, err
	}

	//We need to append to data the constructor's parameters
	data := append(contractBytecode, constructorParams...)
	gas := uint64(0xFFFFFFFF)
	value := new(big.Int).SetUint64(0x00)

	// Deploy the Autonity contract
	_, contractAddress, _, vmerr := evm.Create(sender, data, gas, value)
	if vmerr != nil {
		log.Error("evm.Create returns err", "err", vmerr)
		return contractAddress, vmerr
	}
	ac.Lock()
	ac.address = contractAddress
	ac.Unlock()
	log.Info("Deployed Autonity Contract", "Address", contractAddress.String())

	return contractAddress, nil
}*/

func (ac *Contract) ContractGetCommittee(chain consensus.ChainReader, header *types.Header, statedb *state.StateDB) (types.Committee, error) {
	if header.Number.Cmp(big.NewInt(1)) == 0 && ac.SavedCommitteeRetriever != nil {
		return ac.SavedCommitteeRetriever(1)
	}

	var addresses []common.Address
	err := ac.AutonityContractCall(statedb, header, "getValidators", &addresses)

	if err != nil {
		return nil, err
	}

	sortableAddresses := common.Addresses(addresses)
	sort.Sort(sortableAddresses)

	var committee types.Committee
	for _, val := range sortableAddresses {
		committee = append(committee, types.CommitteeMember{Address: val, VotingPower: new(big.Int).SetUint64(1)})
	}
	return committee, nil
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
		newWhitelist = ac.bc.ReadEnodeWhitelist(false)
	} else {
		// call retrieveWhitelist contract function
		newWhitelist, err = ac.callGetWhitelist(db, block.Header())
	}

	return newWhitelist, err
}

func (ac *Contract) GetMinimumGasPrice(block *types.Block, db *state.StateDB) (uint64, error) {
	if block.Number().Uint64() <= 1 {
		return ac.bc.Config().AutonityContractConfig.MinGasPrice, nil
	}

	return ac.callGetMinimumGasPrice(db, block.Header())
}

func (ac *Contract) SetMinimumGasPrice(block *types.Block, db *state.StateDB, price *big.Int) error {
	if block.Number().Uint64() <= 1 {
		return nil
	}

	return ac.callSetMinimumGasPrice(db, block.Header(), price)
}

func (ac *Contract) ApplyFinalize(transactions types.Transactions, receipts types.Receipts, header *types.Header, statedb *state.StateDB) error {
	if header.Number.Cmp(big.NewInt(1)) < 1 {
		return nil
	}
	blockGas := new(big.Int)
	for i, tx := range transactions {
		blockGas.Add(blockGas, new(big.Int).Mul(tx.GasPrice(), new(big.Int).SetUint64(receipts[i].GasUsed)))
	}
	log.Info("ApplyFinalize", "balance", statedb.GetBalance(ac.Address()), "block", header.Number.Uint64(), "gas", blockGas.Uint64())

	if header.Number.Uint64() <= 1 {
		return nil
	}

	upgradeContract, err := ac.callFinalize(statedb, header, blockGas)
	if err != nil {
		return err
	}

	if upgradeContract {
		// warning prints for failure rather than returning error to stuck engine.
		// in any failure, the state will be rollback to snapshot.
		err = ac.performContractUpgrade(statedb, header)
		if err != nil {
			log.Warn("Autonity Contract Upgrade Failed")
		}
	}

	return nil
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

	//Create account will delete previous the AC stateobject and carry over the balance
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
	log.Info("Autonity Contract upgrade success 🙌")
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
