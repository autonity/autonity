// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package main

import (
	"math/big"
	"strings"

	ethereum "github.com/clearmatics/autonity"
	"github.com/clearmatics/autonity/accounts/abi"
	"github.com/clearmatics/autonity/accounts/abi/bind"
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = abi.U256
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// Struct1 is an auto generated low-level Go binding around an user-defined struct.
type Struct1 struct {
	Addr           common.Address
	UserType       uint8
	Stake          *big.Int
	Enode          string
	CommissionRate *big.Int
}

// Struct0 is an auto generated low-level Go binding around an user-defined struct.
type Struct0 struct {
	Accounts        []common.Address
	Usertypes       []uint8
	Stakes          []*big.Int
	Commissionrates []*big.Int
	Mingasprice     *big.Int
	Stakesupply     *big.Int
}

// Struct2 is an auto generated low-level Go binding around an user-defined struct.
type Struct2 struct {
	Result          bool
	Stakeholders    []common.Address
	Rewardfractions []*big.Int
	Amount          *big.Int
}

// Autonity is an auto generated Go binding around an Ethereum contract.
type Autonity struct {
	AutonityCaller     // Read-only binding to the contract
	AutonityTransactor // Write-only binding to the contract
	AutonityFilterer   // Log filterer for contract events
}

// AutonityCaller is an auto generated read-only Go binding around an Ethereum contract.
type AutonityCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AutonityTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AutonityTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AutonityFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AutonityFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AutonitySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AutonitySession struct {
	Contract     *Autonity         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// AutonityCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AutonityCallerSession struct {
	Contract *AutonityCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// AutonityTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AutonityTransactorSession struct {
	Contract     *AutonityTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// AutonityRaw is an auto generated low-level Go binding around an Ethereum contract.
type AutonityRaw struct {
	Contract *Autonity // Generic contract binding to access the raw methods on
}

// AutonityCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AutonityCallerRaw struct {
	Contract *AutonityCaller // Generic read-only contract binding to access the raw methods on
}

// AutonityTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AutonityTransactorRaw struct {
	Contract *AutonityTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAutonity creates a new instance of Autonity, bound to a specific deployed contract.
func NewAutonity(address common.Address, backend bind.ContractBackend) (*Autonity, error) {
	contract, err := bindAutonity(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Autonity{AutonityCaller: AutonityCaller{contract: contract}, AutonityTransactor: AutonityTransactor{contract: contract}, AutonityFilterer: AutonityFilterer{contract: contract}}, nil
}

// NewAutonityCaller creates a new read-only instance of Autonity, bound to a specific deployed contract.
func NewAutonityCaller(address common.Address, caller bind.ContractCaller) (*AutonityCaller, error) {
	contract, err := bindAutonity(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AutonityCaller{contract: contract}, nil
}

// NewAutonityTransactor creates a new write-only instance of Autonity, bound to a specific deployed contract.
func NewAutonityTransactor(address common.Address, transactor bind.ContractTransactor) (*AutonityTransactor, error) {
	contract, err := bindAutonity(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AutonityTransactor{contract: contract}, nil
}

// NewAutonityFilterer creates a new log filterer instance of Autonity, bound to a specific deployed contract.
func NewAutonityFilterer(address common.Address, filterer bind.ContractFilterer) (*AutonityFilterer, error) {
	contract, err := bindAutonity(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AutonityFilterer{contract: contract}, nil
}

// bindAutonity binds a generic wrapper to an already deployed contract.
func bindAutonity(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(AutonityABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Autonity *AutonityRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Autonity.Contract.AutonityCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Autonity *AutonityRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Autonity.Contract.AutonityTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Autonity *AutonityRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Autonity.Contract.AutonityTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Autonity *AutonityCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Autonity.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Autonity *AutonityTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Autonity.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Autonity *AutonityTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Autonity.Contract.contract.Transact(opts, method, params...)
}

// BondingPeriod is a free data retrieval call binding the contract method 0xc31c6fb9.
//
// Solidity: function bondingPeriod() constant returns(uint256)
func (_Autonity *AutonityCaller) BondingPeriod(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Autonity.contract.Call(opts, out, "bondingPeriod")
	return *ret0, err
}

// BondingPeriod is a free data retrieval call binding the contract method 0xc31c6fb9.
//
// Solidity: function bondingPeriod() constant returns(uint256)
func (_Autonity *AutonitySession) BondingPeriod() (*big.Int, error) {
	return _Autonity.Contract.BondingPeriod(&_Autonity.CallOpts)
}

// BondingPeriod is a free data retrieval call binding the contract method 0xc31c6fb9.
//
// Solidity: function bondingPeriod() constant returns(uint256)
func (_Autonity *AutonityCallerSession) BondingPeriod() (*big.Int, error) {
	return _Autonity.Contract.BondingPeriod(&_Autonity.CallOpts)
}

// CheckMember is a free data retrieval call binding the contract method 0xaaf2e5d8.
//
// Solidity: function checkMember(address _account) constant returns(bool)
func (_Autonity *AutonityCaller) CheckMember(opts *bind.CallOpts, _account common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Autonity.contract.Call(opts, out, "checkMember", _account)
	return *ret0, err
}

// CheckMember is a free data retrieval call binding the contract method 0xaaf2e5d8.
//
// Solidity: function checkMember(address _account) constant returns(bool)
func (_Autonity *AutonitySession) CheckMember(_account common.Address) (bool, error) {
	return _Autonity.Contract.CheckMember(&_Autonity.CallOpts, _account)
}

// CheckMember is a free data retrieval call binding the contract method 0xaaf2e5d8.
//
// Solidity: function checkMember(address _account) constant returns(bool)
func (_Autonity *AutonityCallerSession) CheckMember(_account common.Address) (bool, error) {
	return _Autonity.Contract.CheckMember(&_Autonity.CallOpts, _account)
}

// Committee is a free data retrieval call binding the contract method 0xafe7fcf4.
//
// Solidity: function committee(uint256 ) constant returns(address addr, uint8 userType, uint256 stake, string enode, uint256 commissionRate)
func (_Autonity *AutonityCaller) Committee(opts *bind.CallOpts, arg0 *big.Int) (struct {
	Addr           common.Address
	UserType       uint8
	Stake          *big.Int
	Enode          string
	CommissionRate *big.Int
}, error) {
	ret := new(struct {
		Addr           common.Address
		UserType       uint8
		Stake          *big.Int
		Enode          string
		CommissionRate *big.Int
	})
	out := ret
	err := _Autonity.contract.Call(opts, out, "committee", arg0)
	return *ret, err
}

// Committee is a free data retrieval call binding the contract method 0xafe7fcf4.
//
// Solidity: function committee(uint256 ) constant returns(address addr, uint8 userType, uint256 stake, string enode, uint256 commissionRate)
func (_Autonity *AutonitySession) Committee(arg0 *big.Int) (struct {
	Addr           common.Address
	UserType       uint8
	Stake          *big.Int
	Enode          string
	CommissionRate *big.Int
}, error) {
	return _Autonity.Contract.Committee(&_Autonity.CallOpts, arg0)
}

// Committee is a free data retrieval call binding the contract method 0xafe7fcf4.
//
// Solidity: function committee(uint256 ) constant returns(address addr, uint8 userType, uint256 stake, string enode, uint256 commissionRate)
func (_Autonity *AutonityCallerSession) Committee(arg0 *big.Int) (struct {
	Addr           common.Address
	UserType       uint8
	Stake          *big.Int
	Enode          string
	CommissionRate *big.Int
}, error) {
	return _Autonity.Contract.Committee(&_Autonity.CallOpts, arg0)
}

// CommitteeSize is a free data retrieval call binding the contract method 0x9cf4364b.
//
// Solidity: function committeeSize() constant returns(uint256)
func (_Autonity *AutonityCaller) CommitteeSize(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Autonity.contract.Call(opts, out, "committeeSize")
	return *ret0, err
}

// CommitteeSize is a free data retrieval call binding the contract method 0x9cf4364b.
//
// Solidity: function committeeSize() constant returns(uint256)
func (_Autonity *AutonitySession) CommitteeSize() (*big.Int, error) {
	return _Autonity.Contract.CommitteeSize(&_Autonity.CallOpts)
}

// CommitteeSize is a free data retrieval call binding the contract method 0x9cf4364b.
//
// Solidity: function committeeSize() constant returns(uint256)
func (_Autonity *AutonityCallerSession) CommitteeSize() (*big.Int, error) {
	return _Autonity.Contract.CommitteeSize(&_Autonity.CallOpts)
}

// ContractVersion is a free data retrieval call binding the contract method 0xa0a8e460.
//
// Solidity: function contractVersion() constant returns(string)
func (_Autonity *AutonityCaller) ContractVersion(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _Autonity.contract.Call(opts, out, "contractVersion")
	return *ret0, err
}

// ContractVersion is a free data retrieval call binding the contract method 0xa0a8e460.
//
// Solidity: function contractVersion() constant returns(string)
func (_Autonity *AutonitySession) ContractVersion() (string, error) {
	return _Autonity.Contract.ContractVersion(&_Autonity.CallOpts)
}

// ContractVersion is a free data retrieval call binding the contract method 0xa0a8e460.
//
// Solidity: function contractVersion() constant returns(string)
func (_Autonity *AutonityCallerSession) ContractVersion() (string, error) {
	return _Autonity.Contract.ContractVersion(&_Autonity.CallOpts)
}

// Deployer is a free data retrieval call binding the contract method 0xd5f39488.
//
// Solidity: function deployer() constant returns(address)
func (_Autonity *AutonityCaller) Deployer(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Autonity.contract.Call(opts, out, "deployer")
	return *ret0, err
}

// Deployer is a free data retrieval call binding the contract method 0xd5f39488.
//
// Solidity: function deployer() constant returns(address)
func (_Autonity *AutonitySession) Deployer() (common.Address, error) {
	return _Autonity.Contract.Deployer(&_Autonity.CallOpts)
}

// Deployer is a free data retrieval call binding the contract method 0xd5f39488.
//
// Solidity: function deployer() constant returns(address)
func (_Autonity *AutonityCallerSession) Deployer() (common.Address, error) {
	return _Autonity.Contract.Deployer(&_Autonity.CallOpts)
}

// DumpEconomicsMetricData is a free data retrieval call binding the contract method 0x0f4f1176.
//
// Solidity: function dumpEconomicsMetricData() constant returns(Struct0 economics)
func (_Autonity *AutonityCaller) DumpEconomicsMetricData(opts *bind.CallOpts) (Struct0, error) {
	var (
		ret0 = new(Struct0)
	)
	out := ret0
	err := _Autonity.contract.Call(opts, out, "dumpEconomicsMetricData")
	return *ret0, err
}

// DumpEconomicsMetricData is a free data retrieval call binding the contract method 0x0f4f1176.
//
// Solidity: function dumpEconomicsMetricData() constant returns(Struct0 economics)
func (_Autonity *AutonitySession) DumpEconomicsMetricData() (Struct0, error) {
	return _Autonity.Contract.DumpEconomicsMetricData(&_Autonity.CallOpts)
}

// DumpEconomicsMetricData is a free data retrieval call binding the contract method 0x0f4f1176.
//
// Solidity: function dumpEconomicsMetricData() constant returns(Struct0 economics)
func (_Autonity *AutonityCallerSession) DumpEconomicsMetricData() (Struct0, error) {
	return _Autonity.Contract.DumpEconomicsMetricData(&_Autonity.CallOpts)
}

// EnodesWhitelist is a free data retrieval call binding the contract method 0xa7b05df5.
//
// Solidity: function enodesWhitelist(uint256 ) constant returns(string)
func (_Autonity *AutonityCaller) EnodesWhitelist(opts *bind.CallOpts, arg0 *big.Int) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _Autonity.contract.Call(opts, out, "enodesWhitelist", arg0)
	return *ret0, err
}

// EnodesWhitelist is a free data retrieval call binding the contract method 0xa7b05df5.
//
// Solidity: function enodesWhitelist(uint256 ) constant returns(string)
func (_Autonity *AutonitySession) EnodesWhitelist(arg0 *big.Int) (string, error) {
	return _Autonity.Contract.EnodesWhitelist(&_Autonity.CallOpts, arg0)
}

// EnodesWhitelist is a free data retrieval call binding the contract method 0xa7b05df5.
//
// Solidity: function enodesWhitelist(uint256 ) constant returns(string)
func (_Autonity *AutonityCallerSession) EnodesWhitelist(arg0 *big.Int) (string, error) {
	return _Autonity.Contract.EnodesWhitelist(&_Autonity.CallOpts, arg0)
}

// GetAccountStake is a free data retrieval call binding the contract method 0x5e30913f.
//
// Solidity: function getAccountStake(address _account) constant returns(uint256)
func (_Autonity *AutonityCaller) GetAccountStake(opts *bind.CallOpts, _account common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Autonity.contract.Call(opts, out, "getAccountStake", _account)
	return *ret0, err
}

// GetAccountStake is a free data retrieval call binding the contract method 0x5e30913f.
//
// Solidity: function getAccountStake(address _account) constant returns(uint256)
func (_Autonity *AutonitySession) GetAccountStake(_account common.Address) (*big.Int, error) {
	return _Autonity.Contract.GetAccountStake(&_Autonity.CallOpts, _account)
}

// GetAccountStake is a free data retrieval call binding the contract method 0x5e30913f.
//
// Solidity: function getAccountStake(address _account) constant returns(uint256)
func (_Autonity *AutonityCallerSession) GetAccountStake(_account common.Address) (*big.Int, error) {
	return _Autonity.Contract.GetAccountStake(&_Autonity.CallOpts, _account)
}

// GetCommittee is a free data retrieval call binding the contract method 0xab8f6ffe.
//
// Solidity: function getCommittee() constant returns([]Struct1)
func (_Autonity *AutonityCaller) GetCommittee(opts *bind.CallOpts) ([]Struct1, error) {
	var (
		ret0 = new([]Struct1)
	)
	out := ret0
	err := _Autonity.contract.Call(opts, out, "getCommittee")
	return *ret0, err
}

// GetCommittee is a free data retrieval call binding the contract method 0xab8f6ffe.
//
// Solidity: function getCommittee() constant returns([]Struct1)
func (_Autonity *AutonitySession) GetCommittee() ([]Struct1, error) {
	return _Autonity.Contract.GetCommittee(&_Autonity.CallOpts)
}

// GetCommittee is a free data retrieval call binding the contract method 0xab8f6ffe.
//
// Solidity: function getCommittee() constant returns([]Struct1)
func (_Autonity *AutonityCallerSession) GetCommittee() ([]Struct1, error) {
	return _Autonity.Contract.GetCommittee(&_Autonity.CallOpts)
}

// GetCurrentCommiteeSize is a free data retrieval call binding the contract method 0xfec1830f.
//
// Solidity: function getCurrentCommiteeSize() constant returns(uint256)
func (_Autonity *AutonityCaller) GetCurrentCommiteeSize(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Autonity.contract.Call(opts, out, "getCurrentCommiteeSize")
	return *ret0, err
}

// GetCurrentCommiteeSize is a free data retrieval call binding the contract method 0xfec1830f.
//
// Solidity: function getCurrentCommiteeSize() constant returns(uint256)
func (_Autonity *AutonitySession) GetCurrentCommiteeSize() (*big.Int, error) {
	return _Autonity.Contract.GetCurrentCommiteeSize(&_Autonity.CallOpts)
}

// GetCurrentCommiteeSize is a free data retrieval call binding the contract method 0xfec1830f.
//
// Solidity: function getCurrentCommiteeSize() constant returns(uint256)
func (_Autonity *AutonityCallerSession) GetCurrentCommiteeSize() (*big.Int, error) {
	return _Autonity.Contract.GetCurrentCommiteeSize(&_Autonity.CallOpts)
}

// GetMaxCommitteeSize is a free data retrieval call binding the contract method 0x819b6463.
//
// Solidity: function getMaxCommitteeSize() constant returns(uint256)
func (_Autonity *AutonityCaller) GetMaxCommitteeSize(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Autonity.contract.Call(opts, out, "getMaxCommitteeSize")
	return *ret0, err
}

// GetMaxCommitteeSize is a free data retrieval call binding the contract method 0x819b6463.
//
// Solidity: function getMaxCommitteeSize() constant returns(uint256)
func (_Autonity *AutonitySession) GetMaxCommitteeSize() (*big.Int, error) {
	return _Autonity.Contract.GetMaxCommitteeSize(&_Autonity.CallOpts)
}

// GetMaxCommitteeSize is a free data retrieval call binding the contract method 0x819b6463.
//
// Solidity: function getMaxCommitteeSize() constant returns(uint256)
func (_Autonity *AutonityCallerSession) GetMaxCommitteeSize() (*big.Int, error) {
	return _Autonity.Contract.GetMaxCommitteeSize(&_Autonity.CallOpts)
}

// GetMinimumGasPrice is a free data retrieval call binding the contract method 0xf918379a.
//
// Solidity: function getMinimumGasPrice() constant returns(uint256)
func (_Autonity *AutonityCaller) GetMinimumGasPrice(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Autonity.contract.Call(opts, out, "getMinimumGasPrice")
	return *ret0, err
}

// GetMinimumGasPrice is a free data retrieval call binding the contract method 0xf918379a.
//
// Solidity: function getMinimumGasPrice() constant returns(uint256)
func (_Autonity *AutonitySession) GetMinimumGasPrice() (*big.Int, error) {
	return _Autonity.Contract.GetMinimumGasPrice(&_Autonity.CallOpts)
}

// GetMinimumGasPrice is a free data retrieval call binding the contract method 0xf918379a.
//
// Solidity: function getMinimumGasPrice() constant returns(uint256)
func (_Autonity *AutonityCallerSession) GetMinimumGasPrice() (*big.Int, error) {
	return _Autonity.Contract.GetMinimumGasPrice(&_Autonity.CallOpts)
}

// GetRate is a free data retrieval call binding the contract method 0x37cef791.
//
// Solidity: function getRate(address _account) constant returns(uint256)
func (_Autonity *AutonityCaller) GetRate(opts *bind.CallOpts, _account common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Autonity.contract.Call(opts, out, "getRate", _account)
	return *ret0, err
}

// GetRate is a free data retrieval call binding the contract method 0x37cef791.
//
// Solidity: function getRate(address _account) constant returns(uint256)
func (_Autonity *AutonitySession) GetRate(_account common.Address) (*big.Int, error) {
	return _Autonity.Contract.GetRate(&_Autonity.CallOpts, _account)
}

// GetRate is a free data retrieval call binding the contract method 0x37cef791.
//
// Solidity: function getRate(address _account) constant returns(uint256)
func (_Autonity *AutonityCallerSession) GetRate(_account common.Address) (*big.Int, error) {
	return _Autonity.Contract.GetRate(&_Autonity.CallOpts, _account)
}

// GetStake is a free data retrieval call binding the contract method 0xfc0e3d90.
//
// Solidity: function getStake() constant returns(uint256)
func (_Autonity *AutonityCaller) GetStake(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Autonity.contract.Call(opts, out, "getStake")
	return *ret0, err
}

// GetStake is a free data retrieval call binding the contract method 0xfc0e3d90.
//
// Solidity: function getStake() constant returns(uint256)
func (_Autonity *AutonitySession) GetStake() (*big.Int, error) {
	return _Autonity.Contract.GetStake(&_Autonity.CallOpts)
}

// GetStake is a free data retrieval call binding the contract method 0xfc0e3d90.
//
// Solidity: function getStake() constant returns(uint256)
func (_Autonity *AutonityCallerSession) GetStake() (*big.Int, error) {
	return _Autonity.Contract.GetStake(&_Autonity.CallOpts)
}

// GetStakeholders is a free data retrieval call binding the contract method 0xb6992247.
//
// Solidity: function getStakeholders() constant returns(address[])
func (_Autonity *AutonityCaller) GetStakeholders(opts *bind.CallOpts) ([]common.Address, error) {
	var (
		ret0 = new([]common.Address)
	)
	out := ret0
	err := _Autonity.contract.Call(opts, out, "getStakeholders")
	return *ret0, err
}

// GetStakeholders is a free data retrieval call binding the contract method 0xb6992247.
//
// Solidity: function getStakeholders() constant returns(address[])
func (_Autonity *AutonitySession) GetStakeholders() ([]common.Address, error) {
	return _Autonity.Contract.GetStakeholders(&_Autonity.CallOpts)
}

// GetStakeholders is a free data retrieval call binding the contract method 0xb6992247.
//
// Solidity: function getStakeholders() constant returns(address[])
func (_Autonity *AutonityCallerSession) GetStakeholders() ([]common.Address, error) {
	return _Autonity.Contract.GetStakeholders(&_Autonity.CallOpts)
}

// GetValidators is a free data retrieval call binding the contract method 0xb7ab4db5.
//
// Solidity: function getValidators() constant returns(address[])
func (_Autonity *AutonityCaller) GetValidators(opts *bind.CallOpts) ([]common.Address, error) {
	var (
		ret0 = new([]common.Address)
	)
	out := ret0
	err := _Autonity.contract.Call(opts, out, "getValidators")
	return *ret0, err
}

// GetValidators is a free data retrieval call binding the contract method 0xb7ab4db5.
//
// Solidity: function getValidators() constant returns(address[])
func (_Autonity *AutonitySession) GetValidators() ([]common.Address, error) {
	return _Autonity.Contract.GetValidators(&_Autonity.CallOpts)
}

// GetValidators is a free data retrieval call binding the contract method 0xb7ab4db5.
//
// Solidity: function getValidators() constant returns(address[])
func (_Autonity *AutonityCallerSession) GetValidators() ([]common.Address, error) {
	return _Autonity.Contract.GetValidators(&_Autonity.CallOpts)
}

// GetVersion is a free data retrieval call binding the contract method 0x0d8e6e2c.
//
// Solidity: function getVersion() constant returns(string)
func (_Autonity *AutonityCaller) GetVersion(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _Autonity.contract.Call(opts, out, "getVersion")
	return *ret0, err
}

// GetVersion is a free data retrieval call binding the contract method 0x0d8e6e2c.
//
// Solidity: function getVersion() constant returns(string)
func (_Autonity *AutonitySession) GetVersion() (string, error) {
	return _Autonity.Contract.GetVersion(&_Autonity.CallOpts)
}

// GetVersion is a free data retrieval call binding the contract method 0x0d8e6e2c.
//
// Solidity: function getVersion() constant returns(string)
func (_Autonity *AutonityCallerSession) GetVersion() (string, error) {
	return _Autonity.Contract.GetVersion(&_Autonity.CallOpts)
}

// GetWhitelist is a free data retrieval call binding the contract method 0xd01f63f5.
//
// Solidity: function getWhitelist() constant returns(string[])
func (_Autonity *AutonityCaller) GetWhitelist(opts *bind.CallOpts) ([]string, error) {
	var (
		ret0 = new([]string)
	)
	out := ret0
	err := _Autonity.contract.Call(opts, out, "getWhitelist")
	return *ret0, err
}

// GetWhitelist is a free data retrieval call binding the contract method 0xd01f63f5.
//
// Solidity: function getWhitelist() constant returns(string[])
func (_Autonity *AutonitySession) GetWhitelist() ([]string, error) {
	return _Autonity.Contract.GetWhitelist(&_Autonity.CallOpts)
}

// GetWhitelist is a free data retrieval call binding the contract method 0xd01f63f5.
//
// Solidity: function getWhitelist() constant returns(string[])
func (_Autonity *AutonityCallerSession) GetWhitelist() ([]string, error) {
	return _Autonity.Contract.GetWhitelist(&_Autonity.CallOpts)
}

// MyUserType is a free data retrieval call binding the contract method 0x74be8f02.
//
// Solidity: function myUserType() constant returns(uint8)
func (_Autonity *AutonityCaller) MyUserType(opts *bind.CallOpts) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _Autonity.contract.Call(opts, out, "myUserType")
	return *ret0, err
}

// MyUserType is a free data retrieval call binding the contract method 0x74be8f02.
//
// Solidity: function myUserType() constant returns(uint8)
func (_Autonity *AutonitySession) MyUserType() (uint8, error) {
	return _Autonity.Contract.MyUserType(&_Autonity.CallOpts)
}

// MyUserType is a free data retrieval call binding the contract method 0x74be8f02.
//
// Solidity: function myUserType() constant returns(uint8)
func (_Autonity *AutonityCallerSession) MyUserType() (uint8, error) {
	return _Autonity.Contract.MyUserType(&_Autonity.CallOpts)
}

// OperatorAccount is a free data retrieval call binding the contract method 0x2801643d.
//
// Solidity: function operatorAccount() constant returns(address)
func (_Autonity *AutonityCaller) OperatorAccount(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Autonity.contract.Call(opts, out, "operatorAccount")
	return *ret0, err
}

// OperatorAccount is a free data retrieval call binding the contract method 0x2801643d.
//
// Solidity: function operatorAccount() constant returns(address)
func (_Autonity *AutonitySession) OperatorAccount() (common.Address, error) {
	return _Autonity.Contract.OperatorAccount(&_Autonity.CallOpts)
}

// OperatorAccount is a free data retrieval call binding the contract method 0x2801643d.
//
// Solidity: function operatorAccount() constant returns(address)
func (_Autonity *AutonityCallerSession) OperatorAccount() (common.Address, error) {
	return _Autonity.Contract.OperatorAccount(&_Autonity.CallOpts)
}

// RetrieveContract is a free data retrieval call binding the contract method 0x61d9d615.
//
// Solidity: function retrieveContract() constant returns(string, string)
func (_Autonity *AutonityCaller) RetrieveContract(opts *bind.CallOpts) (string, string, error) {
	var (
		ret0 = new(string)
		ret1 = new(string)
	)
	out := &[]interface{}{
		ret0,
		ret1,
	}
	err := _Autonity.contract.Call(opts, out, "retrieveContract")
	return *ret0, *ret1, err
}

// RetrieveContract is a free data retrieval call binding the contract method 0x61d9d615.
//
// Solidity: function retrieveContract() constant returns(string, string)
func (_Autonity *AutonitySession) RetrieveContract() (string, string, error) {
	return _Autonity.Contract.RetrieveContract(&_Autonity.CallOpts)
}

// RetrieveContract is a free data retrieval call binding the contract method 0x61d9d615.
//
// Solidity: function retrieveContract() constant returns(string, string)
func (_Autonity *AutonityCallerSession) RetrieveContract() (string, string, error) {
	return _Autonity.Contract.RetrieveContract(&_Autonity.CallOpts)
}

// RetrieveState is a free data retrieval call binding the contract method 0x11879449.
//
// Solidity: function retrieveState() constant returns(address[], string[], uint256[], uint256[], uint256[], address, uint256, uint256, uint256, string)
func (_Autonity *AutonityCaller) RetrieveState(opts *bind.CallOpts) ([]common.Address, []string, []*big.Int, []*big.Int, []*big.Int, common.Address, *big.Int, *big.Int, *big.Int, string, error) {
	var (
		ret0 = new([]common.Address)
		ret1 = new([]string)
		ret2 = new([]*big.Int)
		ret3 = new([]*big.Int)
		ret4 = new([]*big.Int)
		ret5 = new(common.Address)
		ret6 = new(*big.Int)
		ret7 = new(*big.Int)
		ret8 = new(*big.Int)
		ret9 = new(string)
	)
	out := &[]interface{}{
		ret0,
		ret1,
		ret2,
		ret3,
		ret4,
		ret5,
		ret6,
		ret7,
		ret8,
		ret9,
	}
	err := _Autonity.contract.Call(opts, out, "retrieveState")
	return *ret0, *ret1, *ret2, *ret3, *ret4, *ret5, *ret6, *ret7, *ret8, *ret9, err
}

// RetrieveState is a free data retrieval call binding the contract method 0x11879449.
//
// Solidity: function retrieveState() constant returns(address[], string[], uint256[], uint256[], uint256[], address, uint256, uint256, uint256, string)
func (_Autonity *AutonitySession) RetrieveState() ([]common.Address, []string, []*big.Int, []*big.Int, []*big.Int, common.Address, *big.Int, *big.Int, *big.Int, string, error) {
	return _Autonity.Contract.RetrieveState(&_Autonity.CallOpts)
}

// RetrieveState is a free data retrieval call binding the contract method 0x11879449.
//
// Solidity: function retrieveState() constant returns(address[], string[], uint256[], uint256[], uint256[], address, uint256, uint256, uint256, string)
func (_Autonity *AutonityCallerSession) RetrieveState() ([]common.Address, []string, []*big.Int, []*big.Int, []*big.Int, common.Address, *big.Int, *big.Int, *big.Int, string, error) {
	return _Autonity.Contract.RetrieveState(&_Autonity.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_Autonity *AutonityCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Autonity.contract.Call(opts, out, "totalSupply")
	return *ret0, err
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_Autonity *AutonitySession) TotalSupply() (*big.Int, error) {
	return _Autonity.Contract.TotalSupply(&_Autonity.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_Autonity *AutonityCallerSession) TotalSupply() (*big.Int, error) {
	return _Autonity.Contract.TotalSupply(&_Autonity.CallOpts)
}

// Validators is a free data retrieval call binding the contract method 0x35aa2e44.
//
// Solidity: function validators(uint256 ) constant returns(address)
func (_Autonity *AutonityCaller) Validators(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Autonity.contract.Call(opts, out, "validators", arg0)
	return *ret0, err
}

// Validators is a free data retrieval call binding the contract method 0x35aa2e44.
//
// Solidity: function validators(uint256 ) constant returns(address)
func (_Autonity *AutonitySession) Validators(arg0 *big.Int) (common.Address, error) {
	return _Autonity.Contract.Validators(&_Autonity.CallOpts, arg0)
}

// Validators is a free data retrieval call binding the contract method 0x35aa2e44.
//
// Solidity: function validators(uint256 ) constant returns(address)
func (_Autonity *AutonityCallerSession) Validators(arg0 *big.Int) (common.Address, error) {
	return _Autonity.Contract.Validators(&_Autonity.CallOpts, arg0)
}

// AddParticipant is a paid mutator transaction binding the contract method 0xb68feb84.
//
// Solidity: function addParticipant(address _address, string _enode) returns()
func (_Autonity *AutonityTransactor) AddParticipant(opts *bind.TransactOpts, _address common.Address, _enode string) (*types.Transaction, error) {
	return _Autonity.contract.Transact(opts, "addParticipant", _address, _enode)
}

// AddParticipant is a paid mutator transaction binding the contract method 0xb68feb84.
//
// Solidity: function addParticipant(address _address, string _enode) returns()
func (_Autonity *AutonitySession) AddParticipant(_address common.Address, _enode string) (*types.Transaction, error) {
	return _Autonity.Contract.AddParticipant(&_Autonity.TransactOpts, _address, _enode)
}

// AddParticipant is a paid mutator transaction binding the contract method 0xb68feb84.
//
// Solidity: function addParticipant(address _address, string _enode) returns()
func (_Autonity *AutonityTransactorSession) AddParticipant(_address common.Address, _enode string) (*types.Transaction, error) {
	return _Autonity.Contract.AddParticipant(&_Autonity.TransactOpts, _address, _enode)
}

// AddStakeholder is a paid mutator transaction binding the contract method 0x27e06247.
//
// Solidity: function addStakeholder(address _address, string _enode, uint256 _stake) returns()
func (_Autonity *AutonityTransactor) AddStakeholder(opts *bind.TransactOpts, _address common.Address, _enode string, _stake *big.Int) (*types.Transaction, error) {
	return _Autonity.contract.Transact(opts, "addStakeholder", _address, _enode, _stake)
}

// AddStakeholder is a paid mutator transaction binding the contract method 0x27e06247.
//
// Solidity: function addStakeholder(address _address, string _enode, uint256 _stake) returns()
func (_Autonity *AutonitySession) AddStakeholder(_address common.Address, _enode string, _stake *big.Int) (*types.Transaction, error) {
	return _Autonity.Contract.AddStakeholder(&_Autonity.TransactOpts, _address, _enode, _stake)
}

// AddStakeholder is a paid mutator transaction binding the contract method 0x27e06247.
//
// Solidity: function addStakeholder(address _address, string _enode, uint256 _stake) returns()
func (_Autonity *AutonityTransactorSession) AddStakeholder(_address common.Address, _enode string, _stake *big.Int) (*types.Transaction, error) {
	return _Autonity.Contract.AddStakeholder(&_Autonity.TransactOpts, _address, _enode, _stake)
}

// AddValidator is a paid mutator transaction binding the contract method 0x01736c35.
//
// Solidity: function addValidator(address _address, uint256 _stake, string _enode) returns()
func (_Autonity *AutonityTransactor) AddValidator(opts *bind.TransactOpts, _address common.Address, _stake *big.Int, _enode string) (*types.Transaction, error) {
	return _Autonity.contract.Transact(opts, "addValidator", _address, _stake, _enode)
}

// AddValidator is a paid mutator transaction binding the contract method 0x01736c35.
//
// Solidity: function addValidator(address _address, uint256 _stake, string _enode) returns()
func (_Autonity *AutonitySession) AddValidator(_address common.Address, _stake *big.Int, _enode string) (*types.Transaction, error) {
	return _Autonity.Contract.AddValidator(&_Autonity.TransactOpts, _address, _stake, _enode)
}

// AddValidator is a paid mutator transaction binding the contract method 0x01736c35.
//
// Solidity: function addValidator(address _address, uint256 _stake, string _enode) returns()
func (_Autonity *AutonityTransactorSession) AddValidator(_address common.Address, _stake *big.Int, _enode string) (*types.Transaction, error) {
	return _Autonity.Contract.AddValidator(&_Autonity.TransactOpts, _address, _stake, _enode)
}

// ChangeUserType is a paid mutator transaction binding the contract method 0x766f1fa6.
//
// Solidity: function changeUserType(address _address, uint8 newUserType) returns()
func (_Autonity *AutonityTransactor) ChangeUserType(opts *bind.TransactOpts, _address common.Address, newUserType uint8) (*types.Transaction, error) {
	return _Autonity.contract.Transact(opts, "changeUserType", _address, newUserType)
}

// ChangeUserType is a paid mutator transaction binding the contract method 0x766f1fa6.
//
// Solidity: function changeUserType(address _address, uint8 newUserType) returns()
func (_Autonity *AutonitySession) ChangeUserType(_address common.Address, newUserType uint8) (*types.Transaction, error) {
	return _Autonity.Contract.ChangeUserType(&_Autonity.TransactOpts, _address, newUserType)
}

// ChangeUserType is a paid mutator transaction binding the contract method 0x766f1fa6.
//
// Solidity: function changeUserType(address _address, uint8 newUserType) returns()
func (_Autonity *AutonityTransactorSession) ChangeUserType(_address common.Address, newUserType uint8) (*types.Transaction, error) {
	return _Autonity.Contract.ChangeUserType(&_Autonity.TransactOpts, _address, newUserType)
}

// Finalize is a paid mutator transaction binding the contract method 0x05261aea.
//
// Solidity: function finalize(uint256 _amount) returns(Struct2 rewarddistribution)
func (_Autonity *AutonityTransactor) Finalize(opts *bind.TransactOpts, _amount *big.Int) (*types.Transaction, error) {
	return _Autonity.contract.Transact(opts, "finalize", _amount)
}

// Finalize is a paid mutator transaction binding the contract method 0x05261aea.
//
// Solidity: function finalize(uint256 _amount) returns(Struct2 rewarddistribution)
func (_Autonity *AutonitySession) Finalize(_amount *big.Int) (*types.Transaction, error) {
	return _Autonity.Contract.Finalize(&_Autonity.TransactOpts, _amount)
}

// Finalize is a paid mutator transaction binding the contract method 0x05261aea.
//
// Solidity: function finalize(uint256 _amount) returns(Struct2 rewarddistribution)
func (_Autonity *AutonityTransactorSession) Finalize(_amount *big.Int) (*types.Transaction, error) {
	return _Autonity.Contract.Finalize(&_Autonity.TransactOpts, _amount)
}

// MintStake is a paid mutator transaction binding the contract method 0xca43c38f.
//
// Solidity: function mintStake(address _account, uint256 _amount) returns()
func (_Autonity *AutonityTransactor) MintStake(opts *bind.TransactOpts, _account common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Autonity.contract.Transact(opts, "mintStake", _account, _amount)
}

// MintStake is a paid mutator transaction binding the contract method 0xca43c38f.
//
// Solidity: function mintStake(address _account, uint256 _amount) returns()
func (_Autonity *AutonitySession) MintStake(_account common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Autonity.Contract.MintStake(&_Autonity.TransactOpts, _account, _amount)
}

// MintStake is a paid mutator transaction binding the contract method 0xca43c38f.
//
// Solidity: function mintStake(address _account, uint256 _amount) returns()
func (_Autonity *AutonityTransactorSession) MintStake(_account common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Autonity.Contract.MintStake(&_Autonity.TransactOpts, _account, _amount)
}

// RedeemStake is a paid mutator transaction binding the contract method 0xdfa6bd46.
//
// Solidity: function redeemStake(address _account, uint256 _amount) returns()
func (_Autonity *AutonityTransactor) RedeemStake(opts *bind.TransactOpts, _account common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Autonity.contract.Transact(opts, "redeemStake", _account, _amount)
}

// RedeemStake is a paid mutator transaction binding the contract method 0xdfa6bd46.
//
// Solidity: function redeemStake(address _account, uint256 _amount) returns()
func (_Autonity *AutonitySession) RedeemStake(_account common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Autonity.Contract.RedeemStake(&_Autonity.TransactOpts, _account, _amount)
}

// RedeemStake is a paid mutator transaction binding the contract method 0xdfa6bd46.
//
// Solidity: function redeemStake(address _account, uint256 _amount) returns()
func (_Autonity *AutonityTransactorSession) RedeemStake(_account common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Autonity.Contract.RedeemStake(&_Autonity.TransactOpts, _account, _amount)
}

// RemoveUser is a paid mutator transaction binding the contract method 0x98575188.
//
// Solidity: function removeUser(address _address) returns()
func (_Autonity *AutonityTransactor) RemoveUser(opts *bind.TransactOpts, _address common.Address) (*types.Transaction, error) {
	return _Autonity.contract.Transact(opts, "removeUser", _address)
}

// RemoveUser is a paid mutator transaction binding the contract method 0x98575188.
//
// Solidity: function removeUser(address _address) returns()
func (_Autonity *AutonitySession) RemoveUser(_address common.Address) (*types.Transaction, error) {
	return _Autonity.Contract.RemoveUser(&_Autonity.TransactOpts, _address)
}

// RemoveUser is a paid mutator transaction binding the contract method 0x98575188.
//
// Solidity: function removeUser(address _address) returns()
func (_Autonity *AutonityTransactorSession) RemoveUser(_address common.Address) (*types.Transaction, error) {
	return _Autonity.Contract.RemoveUser(&_Autonity.TransactOpts, _address)
}

// Send is a paid mutator transaction binding the contract method 0xd0679d34.
//
// Solidity: function send(address _recipient, uint256 _amount) returns(bool)
func (_Autonity *AutonityTransactor) Send(opts *bind.TransactOpts, _recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Autonity.contract.Transact(opts, "send", _recipient, _amount)
}

// Send is a paid mutator transaction binding the contract method 0xd0679d34.
//
// Solidity: function send(address _recipient, uint256 _amount) returns(bool)
func (_Autonity *AutonitySession) Send(_recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Autonity.Contract.Send(&_Autonity.TransactOpts, _recipient, _amount)
}

// Send is a paid mutator transaction binding the contract method 0xd0679d34.
//
// Solidity: function send(address _recipient, uint256 _amount) returns(bool)
func (_Autonity *AutonityTransactorSession) Send(_recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Autonity.Contract.Send(&_Autonity.TransactOpts, _recipient, _amount)
}

// SetCommissionRate is a paid mutator transaction binding the contract method 0x19fac8fd.
//
// Solidity: function setCommissionRate(uint256 rate) returns(bool)
func (_Autonity *AutonityTransactor) SetCommissionRate(opts *bind.TransactOpts, rate *big.Int) (*types.Transaction, error) {
	return _Autonity.contract.Transact(opts, "setCommissionRate", rate)
}

// SetCommissionRate is a paid mutator transaction binding the contract method 0x19fac8fd.
//
// Solidity: function setCommissionRate(uint256 rate) returns(bool)
func (_Autonity *AutonitySession) SetCommissionRate(rate *big.Int) (*types.Transaction, error) {
	return _Autonity.Contract.SetCommissionRate(&_Autonity.TransactOpts, rate)
}

// SetCommissionRate is a paid mutator transaction binding the contract method 0x19fac8fd.
//
// Solidity: function setCommissionRate(uint256 rate) returns(bool)
func (_Autonity *AutonityTransactorSession) SetCommissionRate(rate *big.Int) (*types.Transaction, error) {
	return _Autonity.Contract.SetCommissionRate(&_Autonity.TransactOpts, rate)
}

// SetCommittee is a paid mutator transaction binding the contract method 0xf611d7c9.
//
// Solidity: function setCommittee() returns([]Struct1)
func (_Autonity *AutonityTransactor) SetCommittee(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Autonity.contract.Transact(opts, "setCommittee")
}

// SetCommittee is a paid mutator transaction binding the contract method 0xf611d7c9.
//
// Solidity: function setCommittee() returns([]Struct1)
func (_Autonity *AutonitySession) SetCommittee() (*types.Transaction, error) {
	return _Autonity.Contract.SetCommittee(&_Autonity.TransactOpts)
}

// SetCommittee is a paid mutator transaction binding the contract method 0xf611d7c9.
//
// Solidity: function setCommittee() returns([]Struct1)
func (_Autonity *AutonityTransactorSession) SetCommittee() (*types.Transaction, error) {
	return _Autonity.Contract.SetCommittee(&_Autonity.TransactOpts)
}

// SetCommitteeSize is a paid mutator transaction binding the contract method 0x8bac7dad.
//
// Solidity: function setCommitteeSize(uint256 _size) returns()
func (_Autonity *AutonityTransactor) SetCommitteeSize(opts *bind.TransactOpts, _size *big.Int) (*types.Transaction, error) {
	return _Autonity.contract.Transact(opts, "setCommitteeSize", _size)
}

// SetCommitteeSize is a paid mutator transaction binding the contract method 0x8bac7dad.
//
// Solidity: function setCommitteeSize(uint256 _size) returns()
func (_Autonity *AutonitySession) SetCommitteeSize(_size *big.Int) (*types.Transaction, error) {
	return _Autonity.Contract.SetCommitteeSize(&_Autonity.TransactOpts, _size)
}

// SetCommitteeSize is a paid mutator transaction binding the contract method 0x8bac7dad.
//
// Solidity: function setCommitteeSize(uint256 _size) returns()
func (_Autonity *AutonityTransactorSession) SetCommitteeSize(_size *big.Int) (*types.Transaction, error) {
	return _Autonity.Contract.SetCommitteeSize(&_Autonity.TransactOpts, _size)
}

// SetMinimumGasPrice is a paid mutator transaction binding the contract method 0xd249b31c.
//
// Solidity: function setMinimumGasPrice(uint256 _value) returns()
func (_Autonity *AutonityTransactor) SetMinimumGasPrice(opts *bind.TransactOpts, _value *big.Int) (*types.Transaction, error) {
	return _Autonity.contract.Transact(opts, "setMinimumGasPrice", _value)
}

// SetMinimumGasPrice is a paid mutator transaction binding the contract method 0xd249b31c.
//
// Solidity: function setMinimumGasPrice(uint256 _value) returns()
func (_Autonity *AutonitySession) SetMinimumGasPrice(_value *big.Int) (*types.Transaction, error) {
	return _Autonity.Contract.SetMinimumGasPrice(&_Autonity.TransactOpts, _value)
}

// SetMinimumGasPrice is a paid mutator transaction binding the contract method 0xd249b31c.
//
// Solidity: function setMinimumGasPrice(uint256 _value) returns()
func (_Autonity *AutonityTransactorSession) SetMinimumGasPrice(_value *big.Int) (*types.Transaction, error) {
	return _Autonity.Contract.SetMinimumGasPrice(&_Autonity.TransactOpts, _value)
}

// UpgradeContract is a paid mutator transaction binding the contract method 0xf072929d.
//
// Solidity: function upgradeContract(string _bytecode, string _abi, string _version) returns(bool)
func (_Autonity *AutonityTransactor) UpgradeContract(opts *bind.TransactOpts, _bytecode string, _abi string, _version string) (*types.Transaction, error) {
	return _Autonity.contract.Transact(opts, "upgradeContract", _bytecode, _abi, _version)
}

// UpgradeContract is a paid mutator transaction binding the contract method 0xf072929d.
//
// Solidity: function upgradeContract(string _bytecode, string _abi, string _version) returns(bool)
func (_Autonity *AutonitySession) UpgradeContract(_bytecode string, _abi string, _version string) (*types.Transaction, error) {
	return _Autonity.Contract.UpgradeContract(&_Autonity.TransactOpts, _bytecode, _abi, _version)
}

// UpgradeContract is a paid mutator transaction binding the contract method 0xf072929d.
//
// Solidity: function upgradeContract(string _bytecode, string _abi, string _version) returns(bool)
func (_Autonity *AutonityTransactorSession) UpgradeContract(_bytecode string, _abi string, _version string) (*types.Transaction, error) {
	return _Autonity.Contract.UpgradeContract(&_Autonity.TransactOpts, _bytecode, _abi, _version)
}

// AutonityAddParticipantIterator is returned from FilterAddParticipant and is used to iterate over the raw logs and unpacked data for AddParticipant events raised by the Autonity contract.
type AutonityAddParticipantIterator struct {
	Event *AutonityAddParticipant // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AutonityAddParticipantIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutonityAddParticipant)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AutonityAddParticipant)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AutonityAddParticipantIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AutonityAddParticipantIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AutonityAddParticipant represents a AddParticipant event raised by the Autonity contract.
type AutonityAddParticipant struct {
	Address common.Address
	Stake   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterAddParticipant is a free log retrieval operation binding the contract event 0x9a3241a61899aa3b76752287aeacbe5298c70570fac9796bbf4716964d1a0147.
//
// Solidity: event AddParticipant(address _address, uint256 _stake)
func (_Autonity *AutonityFilterer) FilterAddParticipant(opts *bind.FilterOpts) (*AutonityAddParticipantIterator, error) {

	logs, sub, err := _Autonity.contract.FilterLogs(opts, "AddParticipant")
	if err != nil {
		return nil, err
	}
	return &AutonityAddParticipantIterator{contract: _Autonity.contract, event: "AddParticipant", logs: logs, sub: sub}, nil
}

// WatchAddParticipant is a free log subscription operation binding the contract event 0x9a3241a61899aa3b76752287aeacbe5298c70570fac9796bbf4716964d1a0147.
//
// Solidity: event AddParticipant(address _address, uint256 _stake)
func (_Autonity *AutonityFilterer) WatchAddParticipant(opts *bind.WatchOpts, sink chan<- *AutonityAddParticipant) (event.Subscription, error) {

	logs, sub, err := _Autonity.contract.WatchLogs(opts, "AddParticipant")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AutonityAddParticipant)
				if err := _Autonity.contract.UnpackLog(event, "AddParticipant", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseAddParticipant is a log parse operation binding the contract event 0x9a3241a61899aa3b76752287aeacbe5298c70570fac9796bbf4716964d1a0147.
//
// Solidity: event AddParticipant(address _address, uint256 _stake)
func (_Autonity *AutonityFilterer) ParseAddParticipant(log types.Log) (*AutonityAddParticipant, error) {
	event := new(AutonityAddParticipant)
	if err := _Autonity.contract.UnpackLog(event, "AddParticipant", log); err != nil {
		return nil, err
	}
	return event, nil
}

// AutonityAddStakeholderIterator is returned from FilterAddStakeholder and is used to iterate over the raw logs and unpacked data for AddStakeholder events raised by the Autonity contract.
type AutonityAddStakeholderIterator struct {
	Event *AutonityAddStakeholder // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AutonityAddStakeholderIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutonityAddStakeholder)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AutonityAddStakeholder)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AutonityAddStakeholderIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AutonityAddStakeholderIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AutonityAddStakeholder represents a AddStakeholder event raised by the Autonity contract.
type AutonityAddStakeholder struct {
	Address common.Address
	Stake   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterAddStakeholder is a free log retrieval operation binding the contract event 0xd08cf8a1921ddc51bc560b9f60369fe04e20c696b01c7cf4e8a49c692ee83ed4.
//
// Solidity: event AddStakeholder(address _address, uint256 _stake)
func (_Autonity *AutonityFilterer) FilterAddStakeholder(opts *bind.FilterOpts) (*AutonityAddStakeholderIterator, error) {

	logs, sub, err := _Autonity.contract.FilterLogs(opts, "AddStakeholder")
	if err != nil {
		return nil, err
	}
	return &AutonityAddStakeholderIterator{contract: _Autonity.contract, event: "AddStakeholder", logs: logs, sub: sub}, nil
}

// WatchAddStakeholder is a free log subscription operation binding the contract event 0xd08cf8a1921ddc51bc560b9f60369fe04e20c696b01c7cf4e8a49c692ee83ed4.
//
// Solidity: event AddStakeholder(address _address, uint256 _stake)
func (_Autonity *AutonityFilterer) WatchAddStakeholder(opts *bind.WatchOpts, sink chan<- *AutonityAddStakeholder) (event.Subscription, error) {

	logs, sub, err := _Autonity.contract.WatchLogs(opts, "AddStakeholder")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AutonityAddStakeholder)
				if err := _Autonity.contract.UnpackLog(event, "AddStakeholder", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseAddStakeholder is a log parse operation binding the contract event 0xd08cf8a1921ddc51bc560b9f60369fe04e20c696b01c7cf4e8a49c692ee83ed4.
//
// Solidity: event AddStakeholder(address _address, uint256 _stake)
func (_Autonity *AutonityFilterer) ParseAddStakeholder(log types.Log) (*AutonityAddStakeholder, error) {
	event := new(AutonityAddStakeholder)
	if err := _Autonity.contract.UnpackLog(event, "AddStakeholder", log); err != nil {
		return nil, err
	}
	return event, nil
}

// AutonityAddValidatorIterator is returned from FilterAddValidator and is used to iterate over the raw logs and unpacked data for AddValidator events raised by the Autonity contract.
type AutonityAddValidatorIterator struct {
	Event *AutonityAddValidator // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AutonityAddValidatorIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutonityAddValidator)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AutonityAddValidator)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AutonityAddValidatorIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AutonityAddValidatorIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AutonityAddValidator represents a AddValidator event raised by the Autonity contract.
type AutonityAddValidator struct {
	Address common.Address
	Stake   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterAddValidator is a free log retrieval operation binding the contract event 0x228a1437a402e19b16880154e2c1f2edc5600a20524c05d21f880e2efefe54ae.
//
// Solidity: event AddValidator(address _address, uint256 _stake)
func (_Autonity *AutonityFilterer) FilterAddValidator(opts *bind.FilterOpts) (*AutonityAddValidatorIterator, error) {

	logs, sub, err := _Autonity.contract.FilterLogs(opts, "AddValidator")
	if err != nil {
		return nil, err
	}
	return &AutonityAddValidatorIterator{contract: _Autonity.contract, event: "AddValidator", logs: logs, sub: sub}, nil
}

// WatchAddValidator is a free log subscription operation binding the contract event 0x228a1437a402e19b16880154e2c1f2edc5600a20524c05d21f880e2efefe54ae.
//
// Solidity: event AddValidator(address _address, uint256 _stake)
func (_Autonity *AutonityFilterer) WatchAddValidator(opts *bind.WatchOpts, sink chan<- *AutonityAddValidator) (event.Subscription, error) {

	logs, sub, err := _Autonity.contract.WatchLogs(opts, "AddValidator")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AutonityAddValidator)
				if err := _Autonity.contract.UnpackLog(event, "AddValidator", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseAddValidator is a log parse operation binding the contract event 0x228a1437a402e19b16880154e2c1f2edc5600a20524c05d21f880e2efefe54ae.
//
// Solidity: event AddValidator(address _address, uint256 _stake)
func (_Autonity *AutonityFilterer) ParseAddValidator(log types.Log) (*AutonityAddValidator, error) {
	event := new(AutonityAddValidator)
	if err := _Autonity.contract.UnpackLog(event, "AddValidator", log); err != nil {
		return nil, err
	}
	return event, nil
}

// AutonityChangeUserTypeIterator is returned from FilterChangeUserType and is used to iterate over the raw logs and unpacked data for ChangeUserType events raised by the Autonity contract.
type AutonityChangeUserTypeIterator struct {
	Event *AutonityChangeUserType // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AutonityChangeUserTypeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutonityChangeUserType)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AutonityChangeUserType)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AutonityChangeUserTypeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AutonityChangeUserTypeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AutonityChangeUserType represents a ChangeUserType event raised by the Autonity contract.
type AutonityChangeUserType struct {
	Address common.Address
	OldType uint8
	NewType uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterChangeUserType is a free log retrieval operation binding the contract event 0xd0130a9bb3215941e8f5c45cbc58d41998d4d14fbe446b04f70ab24d78c847ff.
//
// Solidity: event ChangeUserType(address _address, uint8 _oldType, uint8 _newType)
func (_Autonity *AutonityFilterer) FilterChangeUserType(opts *bind.FilterOpts) (*AutonityChangeUserTypeIterator, error) {

	logs, sub, err := _Autonity.contract.FilterLogs(opts, "ChangeUserType")
	if err != nil {
		return nil, err
	}
	return &AutonityChangeUserTypeIterator{contract: _Autonity.contract, event: "ChangeUserType", logs: logs, sub: sub}, nil
}

// WatchChangeUserType is a free log subscription operation binding the contract event 0xd0130a9bb3215941e8f5c45cbc58d41998d4d14fbe446b04f70ab24d78c847ff.
//
// Solidity: event ChangeUserType(address _address, uint8 _oldType, uint8 _newType)
func (_Autonity *AutonityFilterer) WatchChangeUserType(opts *bind.WatchOpts, sink chan<- *AutonityChangeUserType) (event.Subscription, error) {

	logs, sub, err := _Autonity.contract.WatchLogs(opts, "ChangeUserType")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AutonityChangeUserType)
				if err := _Autonity.contract.UnpackLog(event, "ChangeUserType", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseChangeUserType is a log parse operation binding the contract event 0xd0130a9bb3215941e8f5c45cbc58d41998d4d14fbe446b04f70ab24d78c847ff.
//
// Solidity: event ChangeUserType(address _address, uint8 _oldType, uint8 _newType)
func (_Autonity *AutonityFilterer) ParseChangeUserType(log types.Log) (*AutonityChangeUserType, error) {
	event := new(AutonityChangeUserType)
	if err := _Autonity.contract.UnpackLog(event, "ChangeUserType", log); err != nil {
		return nil, err
	}
	return event, nil
}

// AutonityMintStakeIterator is returned from FilterMintStake and is used to iterate over the raw logs and unpacked data for MintStake events raised by the Autonity contract.
type AutonityMintStakeIterator struct {
	Event *AutonityMintStake // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AutonityMintStakeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutonityMintStake)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AutonityMintStake)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AutonityMintStakeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AutonityMintStakeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AutonityMintStake represents a MintStake event raised by the Autonity contract.
type AutonityMintStake struct {
	Address common.Address
	Amount  *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterMintStake is a free log retrieval operation binding the contract event 0x96a9a8981a322aeae183999165c1fa2610a0c066a01fe86ae3194afade9b4968.
//
// Solidity: event MintStake(address _address, uint256 _amount)
func (_Autonity *AutonityFilterer) FilterMintStake(opts *bind.FilterOpts) (*AutonityMintStakeIterator, error) {

	logs, sub, err := _Autonity.contract.FilterLogs(opts, "MintStake")
	if err != nil {
		return nil, err
	}
	return &AutonityMintStakeIterator{contract: _Autonity.contract, event: "MintStake", logs: logs, sub: sub}, nil
}

// WatchMintStake is a free log subscription operation binding the contract event 0x96a9a8981a322aeae183999165c1fa2610a0c066a01fe86ae3194afade9b4968.
//
// Solidity: event MintStake(address _address, uint256 _amount)
func (_Autonity *AutonityFilterer) WatchMintStake(opts *bind.WatchOpts, sink chan<- *AutonityMintStake) (event.Subscription, error) {

	logs, sub, err := _Autonity.contract.WatchLogs(opts, "MintStake")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AutonityMintStake)
				if err := _Autonity.contract.UnpackLog(event, "MintStake", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseMintStake is a log parse operation binding the contract event 0x96a9a8981a322aeae183999165c1fa2610a0c066a01fe86ae3194afade9b4968.
//
// Solidity: event MintStake(address _address, uint256 _amount)
func (_Autonity *AutonityFilterer) ParseMintStake(log types.Log) (*AutonityMintStake, error) {
	event := new(AutonityMintStake)
	if err := _Autonity.contract.UnpackLog(event, "MintStake", log); err != nil {
		return nil, err
	}
	return event, nil
}

// AutonityRedeemStakeIterator is returned from FilterRedeemStake and is used to iterate over the raw logs and unpacked data for RedeemStake events raised by the Autonity contract.
type AutonityRedeemStakeIterator struct {
	Event *AutonityRedeemStake // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AutonityRedeemStakeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutonityRedeemStake)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AutonityRedeemStake)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AutonityRedeemStakeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AutonityRedeemStakeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AutonityRedeemStake represents a RedeemStake event raised by the Autonity contract.
type AutonityRedeemStake struct {
	Address common.Address
	Amount  *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRedeemStake is a free log retrieval operation binding the contract event 0x4258db2358b464608335ef14dc2734bb42b15a6d03279d5cf12cb066af068f9c.
//
// Solidity: event RedeemStake(address _address, uint256 _amount)
func (_Autonity *AutonityFilterer) FilterRedeemStake(opts *bind.FilterOpts) (*AutonityRedeemStakeIterator, error) {

	logs, sub, err := _Autonity.contract.FilterLogs(opts, "RedeemStake")
	if err != nil {
		return nil, err
	}
	return &AutonityRedeemStakeIterator{contract: _Autonity.contract, event: "RedeemStake", logs: logs, sub: sub}, nil
}

// WatchRedeemStake is a free log subscription operation binding the contract event 0x4258db2358b464608335ef14dc2734bb42b15a6d03279d5cf12cb066af068f9c.
//
// Solidity: event RedeemStake(address _address, uint256 _amount)
func (_Autonity *AutonityFilterer) WatchRedeemStake(opts *bind.WatchOpts, sink chan<- *AutonityRedeemStake) (event.Subscription, error) {

	logs, sub, err := _Autonity.contract.WatchLogs(opts, "RedeemStake")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AutonityRedeemStake)
				if err := _Autonity.contract.UnpackLog(event, "RedeemStake", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRedeemStake is a log parse operation binding the contract event 0x4258db2358b464608335ef14dc2734bb42b15a6d03279d5cf12cb066af068f9c.
//
// Solidity: event RedeemStake(address _address, uint256 _amount)
func (_Autonity *AutonityFilterer) ParseRedeemStake(log types.Log) (*AutonityRedeemStake, error) {
	event := new(AutonityRedeemStake)
	if err := _Autonity.contract.UnpackLog(event, "RedeemStake", log); err != nil {
		return nil, err
	}
	return event, nil
}

// AutonityRemoveUserIterator is returned from FilterRemoveUser and is used to iterate over the raw logs and unpacked data for RemoveUser events raised by the Autonity contract.
type AutonityRemoveUserIterator struct {
	Event *AutonityRemoveUser // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AutonityRemoveUserIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutonityRemoveUser)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AutonityRemoveUser)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AutonityRemoveUserIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AutonityRemoveUserIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AutonityRemoveUser represents a RemoveUser event raised by the Autonity contract.
type AutonityRemoveUser struct {
	Address common.Address
	Type    uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRemoveUser is a free log retrieval operation binding the contract event 0x0a9b5000d97f68a05b3d86a812e2d8e403fc40244cff1942ccc94fb4b96757d9.
//
// Solidity: event RemoveUser(address _address, uint8 _type)
func (_Autonity *AutonityFilterer) FilterRemoveUser(opts *bind.FilterOpts) (*AutonityRemoveUserIterator, error) {

	logs, sub, err := _Autonity.contract.FilterLogs(opts, "RemoveUser")
	if err != nil {
		return nil, err
	}
	return &AutonityRemoveUserIterator{contract: _Autonity.contract, event: "RemoveUser", logs: logs, sub: sub}, nil
}

// WatchRemoveUser is a free log subscription operation binding the contract event 0x0a9b5000d97f68a05b3d86a812e2d8e403fc40244cff1942ccc94fb4b96757d9.
//
// Solidity: event RemoveUser(address _address, uint8 _type)
func (_Autonity *AutonityFilterer) WatchRemoveUser(opts *bind.WatchOpts, sink chan<- *AutonityRemoveUser) (event.Subscription, error) {

	logs, sub, err := _Autonity.contract.WatchLogs(opts, "RemoveUser")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AutonityRemoveUser)
				if err := _Autonity.contract.UnpackLog(event, "RemoveUser", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRemoveUser is a log parse operation binding the contract event 0x0a9b5000d97f68a05b3d86a812e2d8e403fc40244cff1942ccc94fb4b96757d9.
//
// Solidity: event RemoveUser(address _address, uint8 _type)
func (_Autonity *AutonityFilterer) ParseRemoveUser(log types.Log) (*AutonityRemoveUser, error) {
	event := new(AutonityRemoveUser)
	if err := _Autonity.contract.UnpackLog(event, "RemoveUser", log); err != nil {
		return nil, err
	}
	return event, nil
}

// AutonitySetCommissionRateIterator is returned from FilterSetCommissionRate and is used to iterate over the raw logs and unpacked data for SetCommissionRate events raised by the Autonity contract.
type AutonitySetCommissionRateIterator struct {
	Event *AutonitySetCommissionRate // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AutonitySetCommissionRateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutonitySetCommissionRate)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AutonitySetCommissionRate)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AutonitySetCommissionRateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AutonitySetCommissionRateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AutonitySetCommissionRate represents a SetCommissionRate event raised by the Autonity contract.
type AutonitySetCommissionRate struct {
	Address common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterSetCommissionRate is a free log retrieval operation binding the contract event 0xfb621a017bb038be49d13b22e821cbca1b2f153f0a4933795e7a363aa47fdf88.
//
// Solidity: event SetCommissionRate(address _address, uint256 _value)
func (_Autonity *AutonityFilterer) FilterSetCommissionRate(opts *bind.FilterOpts) (*AutonitySetCommissionRateIterator, error) {

	logs, sub, err := _Autonity.contract.FilterLogs(opts, "SetCommissionRate")
	if err != nil {
		return nil, err
	}
	return &AutonitySetCommissionRateIterator{contract: _Autonity.contract, event: "SetCommissionRate", logs: logs, sub: sub}, nil
}

// WatchSetCommissionRate is a free log subscription operation binding the contract event 0xfb621a017bb038be49d13b22e821cbca1b2f153f0a4933795e7a363aa47fdf88.
//
// Solidity: event SetCommissionRate(address _address, uint256 _value)
func (_Autonity *AutonityFilterer) WatchSetCommissionRate(opts *bind.WatchOpts, sink chan<- *AutonitySetCommissionRate) (event.Subscription, error) {

	logs, sub, err := _Autonity.contract.WatchLogs(opts, "SetCommissionRate")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AutonitySetCommissionRate)
				if err := _Autonity.contract.UnpackLog(event, "SetCommissionRate", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseSetCommissionRate is a log parse operation binding the contract event 0xfb621a017bb038be49d13b22e821cbca1b2f153f0a4933795e7a363aa47fdf88.
//
// Solidity: event SetCommissionRate(address _address, uint256 _value)
func (_Autonity *AutonityFilterer) ParseSetCommissionRate(log types.Log) (*AutonitySetCommissionRate, error) {
	event := new(AutonitySetCommissionRate)
	if err := _Autonity.contract.UnpackLog(event, "SetCommissionRate", log); err != nil {
		return nil, err
	}
	return event, nil
}

// AutonitySetMinimumGasPriceIterator is returned from FilterSetMinimumGasPrice and is used to iterate over the raw logs and unpacked data for SetMinimumGasPrice events raised by the Autonity contract.
type AutonitySetMinimumGasPriceIterator struct {
	Event *AutonitySetMinimumGasPrice // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AutonitySetMinimumGasPriceIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutonitySetMinimumGasPrice)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AutonitySetMinimumGasPrice)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AutonitySetMinimumGasPriceIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AutonitySetMinimumGasPriceIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AutonitySetMinimumGasPrice represents a SetMinimumGasPrice event raised by the Autonity contract.
type AutonitySetMinimumGasPrice struct {
	GasPrice *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterSetMinimumGasPrice is a free log retrieval operation binding the contract event 0xb58ce08a43dbde3538e0851b84afb70f6ffe3ecfbc4d8383e9e92d552f9b41bb.
//
// Solidity: event SetMinimumGasPrice(uint256 _gasPrice)
func (_Autonity *AutonityFilterer) FilterSetMinimumGasPrice(opts *bind.FilterOpts) (*AutonitySetMinimumGasPriceIterator, error) {

	logs, sub, err := _Autonity.contract.FilterLogs(opts, "SetMinimumGasPrice")
	if err != nil {
		return nil, err
	}
	return &AutonitySetMinimumGasPriceIterator{contract: _Autonity.contract, event: "SetMinimumGasPrice", logs: logs, sub: sub}, nil
}

// WatchSetMinimumGasPrice is a free log subscription operation binding the contract event 0xb58ce08a43dbde3538e0851b84afb70f6ffe3ecfbc4d8383e9e92d552f9b41bb.
//
// Solidity: event SetMinimumGasPrice(uint256 _gasPrice)
func (_Autonity *AutonityFilterer) WatchSetMinimumGasPrice(opts *bind.WatchOpts, sink chan<- *AutonitySetMinimumGasPrice) (event.Subscription, error) {

	logs, sub, err := _Autonity.contract.WatchLogs(opts, "SetMinimumGasPrice")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AutonitySetMinimumGasPrice)
				if err := _Autonity.contract.UnpackLog(event, "SetMinimumGasPrice", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseSetMinimumGasPrice is a log parse operation binding the contract event 0xb58ce08a43dbde3538e0851b84afb70f6ffe3ecfbc4d8383e9e92d552f9b41bb.
//
// Solidity: event SetMinimumGasPrice(uint256 _gasPrice)
func (_Autonity *AutonityFilterer) ParseSetMinimumGasPrice(log types.Log) (*AutonitySetMinimumGasPrice, error) {
	event := new(AutonitySetMinimumGasPrice)
	if err := _Autonity.contract.UnpackLog(event, "SetMinimumGasPrice", log); err != nil {
		return nil, err
	}
	return event, nil
}

// AutonityTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the Autonity contract.
type AutonityTransferIterator struct {
	Event *AutonityTransfer // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AutonityTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutonityTransfer)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AutonityTransfer)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AutonityTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AutonityTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AutonityTransfer represents a Transfer event raised by the Autonity contract.
type AutonityTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_Autonity *AutonityFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*AutonityTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Autonity.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &AutonityTransferIterator{contract: _Autonity.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_Autonity *AutonityFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *AutonityTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Autonity.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AutonityTransfer)
				if err := _Autonity.contract.UnpackLog(event, "Transfer", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_Autonity *AutonityFilterer) ParseTransfer(log types.Log) (*AutonityTransfer, error) {
	event := new(AutonityTransfer)
	if err := _Autonity.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	return event, nil
}

// AutonityVersionIterator is returned from FilterVersion and is used to iterate over the raw logs and unpacked data for Version events raised by the Autonity contract.
type AutonityVersionIterator struct {
	Event *AutonityVersion // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AutonityVersionIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutonityVersion)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AutonityVersion)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AutonityVersionIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AutonityVersionIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AutonityVersion represents a Version event raised by the Autonity contract.
type AutonityVersion struct {
	Version string
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterVersion is a free log retrieval operation binding the contract event 0x2d8e32500d5075f9103fe815a436a455ffe264e3b6a008c88d40a21b21b1204d.
//
// Solidity: event Version(string version)
func (_Autonity *AutonityFilterer) FilterVersion(opts *bind.FilterOpts) (*AutonityVersionIterator, error) {

	logs, sub, err := _Autonity.contract.FilterLogs(opts, "Version")
	if err != nil {
		return nil, err
	}
	return &AutonityVersionIterator{contract: _Autonity.contract, event: "Version", logs: logs, sub: sub}, nil
}

// WatchVersion is a free log subscription operation binding the contract event 0x2d8e32500d5075f9103fe815a436a455ffe264e3b6a008c88d40a21b21b1204d.
//
// Solidity: event Version(string version)
func (_Autonity *AutonityFilterer) WatchVersion(opts *bind.WatchOpts, sink chan<- *AutonityVersion) (event.Subscription, error) {

	logs, sub, err := _Autonity.contract.WatchLogs(opts, "Version")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AutonityVersion)
				if err := _Autonity.contract.UnpackLog(event, "Version", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseVersion is a log parse operation binding the contract event 0x2d8e32500d5075f9103fe815a436a455ffe264e3b6a008c88d40a21b21b1204d.
//
// Solidity: event Version(string version)
func (_Autonity *AutonityFilterer) ParseVersion(log types.Log) (*AutonityVersion, error) {
	event := new(AutonityVersion)
	if err := _Autonity.contract.UnpackLog(event, "Version", log); err != nil {
		return nil, err
	}
	return event, nil
}
