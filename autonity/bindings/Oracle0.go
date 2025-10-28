// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package bindings

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/autonity/autonity"
	"github.com/autonity/autonity/accounts/abi"
	"github.com/autonity/autonity/accounts/abi/bind"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// IAccountabilityBaseSlashingRates is an auto generated low-level Go binding around an user-defined struct.
type IAccountabilityBaseSlashingRates struct {
	Low  *big.Int
	Mid  *big.Int
	High *big.Int
}

// IAccountabilityConfig is an auto generated low-level Go binding around an user-defined struct.
type IAccountabilityConfig struct {
	InnocenceProofSubmissionWindow *big.Int
	Delta                          *big.Int
	Range                          *big.Int
	BaseSlashingRates              IAccountabilityBaseSlashingRates
	Factors                        IAccountabilityFactors
}

// IAccountabilityFactors is an auto generated low-level Go binding around an user-defined struct.
type IAccountabilityFactors struct {
	Collusion *big.Int
	History   *big.Int
	Jail      *big.Int
}

// IAutonityAccountability is an auto generated low-level Go binding around an user-defined struct.
type IAutonityAccountability struct {
	Range       *big.Int
	Delta       *big.Int
	GracePeriod *big.Int
}

// IAutonityBondingRequest is an auto generated low-level Go binding around an user-defined struct.
type IAutonityBondingRequest struct {
	Delegator    common.Address
	Delegatee    common.Address
	Amount       *big.Int
	RequestBlock *big.Int
}

// IAutonityClientAwareConfig is an auto generated low-level Go binding around an user-defined struct.
type IAutonityClientAwareConfig struct {
	EpochPeriod         *big.Int
	BlockPeriod         *big.Int
	GasLimit            *big.Int
	ClusteringThreshold *big.Int
	Accountability      IAutonityAccountability
	Eip1559             IAutonityEip1559
}

// IAutonityCommitteeMember is an auto generated low-level Go binding around an user-defined struct.
type IAutonityCommitteeMember struct {
	Addr         common.Address
	VotingPower  *big.Int
	ConsensusKey []byte
}

// IAutonityConfig is an auto generated low-level Go binding around an user-defined struct.
type IAutonityConfig struct {
	Policy          IAutonityPolicy
	Contracts       IAutonityContracts
	Protocol        IAutonityProtocol
	ContractVersion *big.Int
}

// IAutonityContracts is an auto generated low-level Go binding around an user-defined struct.
type IAutonityContracts struct {
	AccountabilityContract         common.Address
	OracleContract                 common.Address
	AcuContract                    common.Address
	SupplyControlContract          common.Address
	StabilizationContract          common.Address
	UpgradeManagerContract         common.Address
	InflationControllerContract    common.Address
	OmissionAccountabilityContract common.Address
	AuctioneerContract             common.Address
}

// IAutonityEip1559 is an auto generated low-level Go binding around an user-defined struct.
type IAutonityEip1559 struct {
	MinBaseFee               *big.Int
	BaseFeeChangeDenominator *big.Int
	ElasticityMultiplier     *big.Int
	GasLimitBoundDivisor     *big.Int
}

// IAutonityEpochInfo is an auto generated low-level Go binding around an user-defined struct.
type IAutonityEpochInfo struct {
	Committee          []IAutonityCommitteeMember
	PreviousEpochBlock *big.Int
	EpochBlock         *big.Int
	NextEpochBlock     *big.Int
	OmissionDelta      *big.Int
	Eip1559            IAutonityEip1559
}

// IAutonityPolicy is an auto generated low-level Go binding around an user-defined struct.
type IAutonityPolicy struct {
	TreasuryFee              *big.Int
	MinBaseFee               *big.Int
	DelegationRate           *big.Int
	UnbondingPeriod          *big.Int
	InitialInflationReserve  *big.Int
	WithholdingThreshold     *big.Int
	ProposerRewardRate       *big.Int
	OracleRewardRate         *big.Int
	WithheldRewardsPool      common.Address
	TreasuryAccount          common.Address
	BaseFeeChangeDenominator *big.Int
	ElasticityMultiplier     *big.Int
}

// IAutonityProtocol is an auto generated low-level Go binding around an user-defined struct.
type IAutonityProtocol struct {
	OperatorAccount      common.Address
	EpochPeriod          *big.Int
	BlockPeriod          *big.Int
	CommitteeSize        *big.Int
	MaxScheduleDuration  *big.Int
	GasLimit             *big.Int
	GasLimitBoundDivisor *big.Int
	ClusteringThreshold  *big.Int
}

// IAutonityUnbondingRequest is an auto generated low-level Go binding around an user-defined struct.
type IAutonityUnbondingRequest struct {
	Delegator      common.Address
	Delegatee      common.Address
	Amount         *big.Int
	UnbondingShare *big.Int
	RequestBlock   *big.Int
	Unlocked       bool
	Released       bool
	SelfDelegation bool
}

// IAutonityValidator is an auto generated low-level Go binding around an user-defined struct.
type IAutonityValidator struct {
	Treasury                 common.Address
	NodeAddress              common.Address
	OracleAddress            common.Address
	Enode                    string
	CommissionRate           *big.Int
	BondedStake              *big.Int
	UnbondingStake           *big.Int
	UnbondingShares          *big.Int
	SelfBondedStake          *big.Int
	SelfUnbondingStake       *big.Int
	SelfUnbondingShares      *big.Int
	SelfUnbondingStakeLocked *big.Int
	LiquidStateContract      common.Address
	LiquidSupply             *big.Int
	RegistrationBlock        *big.Int
	TotalSlashed             *big.Int
	JailReleaseBlock         *big.Int
	ConsensusKey             []byte
	State                    uint8
	ConversionRatio          *big.Int
}

// IOracleReport is an auto generated low-level Go binding around an user-defined struct.
type IOracleReport struct {
	Price      *big.Int
	Confidence uint8
}

// IOracleRoundData is an auto generated low-level Go binding around an user-defined struct.
type IOracleRoundData struct {
	Round     *big.Int
	Price     *big.Int
	Timestamp *big.Int
	Success   bool
}

// IScheduleControllerSchedule is an auto generated low-level Go binding around an user-defined struct.
type IScheduleControllerSchedule struct {
	TotalAmount    *big.Int
	UnlockedAmount *big.Int
	Start          *big.Int
	TotalDuration  *big.Int
	LastUnlockTime *big.Int
}

// IStabilizationCDP is an auto generated low-level Go binding around an user-defined struct.
type IStabilizationCDP struct {
	Timestamp                      *big.Int
	Collateral                     *big.Int
	Principal                      *big.Int
	Interest                       *big.Int
	LastAggregatedInterestExponent *big.Int
}

// IStabilizationConfig is an auto generated low-level Go binding around an user-defined struct.
type IStabilizationConfig struct {
	BorrowInterestRate        *big.Int
	AnnouncementWindow        *big.Int
	LiquidationRatio          *big.Int
	MinCollateralizationRatio *big.Int
	MinDebtRequirement        *big.Int
	TargetPrice               *big.Int
	DefaultNTNATNPrice        *big.Int
	DefaultNTNUSDPrice        *big.Int
	DefaultACUUSDPrice        *big.Int
}

// IStabilizationLastUpdated is an auto generated low-level Go binding around an user-defined struct.
type IStabilizationLastUpdated struct {
	BorrowInterestRateTimestamp        *big.Int
	AnnouncementWindowTimestamp        *big.Int
	LiquidationRatioTimestamp          *big.Int
	MinCollateralizationRatioTimestamp *big.Int
}

// Oracle0Config is an auto generated low-level Go binding around an user-defined struct.
type Oracle0Config struct {
	Autonity                  common.Address
	Operator                  common.Address
	VotePeriod                *big.Int
	OutlierDetectionThreshold *big.Int
	OutlierSlashingThreshold  *big.Int
	BaseSlashingRate          *big.Int
	NonRevealThreshold        *big.Int
	RevealResetInterval       *big.Int
	SlashingRateCap           *big.Int
}

// Oracle0VoterInfo is an auto generated low-level Go binding around an user-defined struct.
type Oracle0VoterInfo struct {
	Round           *big.Int
	Commit          *big.Int
	Performance     *big.Int
	NonRevealCount  *big.Int
	IsVoter         bool
	ReportAvailable bool
}

// EnumerableSetMetaData contains all meta data concerning the EnumerableSet contract.
var EnumerableSetMetaData = &bind.MetaData{
	ABI: "[]",
	Bin: "0x60566037600b82828239805160001a607314602a57634e487b7160e01b600052600060045260246000fd5b30600052607381538281f3fe73000000000000000000000000000000000000000030146080604052600080fdfea264697066735822122090b2b6a56cc83a8aeac938298e2604999288403d1f4a15a5d50dddeda211624564736f6c634300081e0033",
}

// EnumerableSetABI is the input ABI used to generate the binding from.
// Deprecated: Use EnumerableSetMetaData.ABI instead.
var EnumerableSetABI = EnumerableSetMetaData.ABI

// EnumerableSetBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use EnumerableSetMetaData.Bin instead.
var EnumerableSetBin = EnumerableSetMetaData.Bin

// DeployEnumerableSet deploys a new Ethereum contract, binding an instance of EnumerableSet to it.
func DeployEnumerableSet(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *EnumerableSet, error) {
	parsed, err := EnumerableSetMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(EnumerableSetBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &EnumerableSet{EnumerableSetCaller: EnumerableSetCaller{contract: contract}, EnumerableSetTransactor: EnumerableSetTransactor{contract: contract}, EnumerableSetFilterer: EnumerableSetFilterer{contract: contract}}, nil
}

// EnumerableSet is an auto generated Go binding around an Ethereum contract.
type EnumerableSet struct {
	EnumerableSetCaller     // Read-only binding to the contract
	EnumerableSetTransactor // Write-only binding to the contract
	EnumerableSetFilterer   // Log filterer for contract events
}

// EnumerableSetCaller is an auto generated read-only Go binding around an Ethereum contract.
type EnumerableSetCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EnumerableSetTransactor is an auto generated write-only Go binding around an Ethereum contract.
type EnumerableSetTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EnumerableSetFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type EnumerableSetFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EnumerableSetSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type EnumerableSetSession struct {
	Contract     *EnumerableSet    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// EnumerableSetCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type EnumerableSetCallerSession struct {
	Contract *EnumerableSetCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// EnumerableSetTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type EnumerableSetTransactorSession struct {
	Contract     *EnumerableSetTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// EnumerableSetRaw is an auto generated low-level Go binding around an Ethereum contract.
type EnumerableSetRaw struct {
	Contract *EnumerableSet // Generic contract binding to access the raw methods on
}

// EnumerableSetCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type EnumerableSetCallerRaw struct {
	Contract *EnumerableSetCaller // Generic read-only contract binding to access the raw methods on
}

// EnumerableSetTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type EnumerableSetTransactorRaw struct {
	Contract *EnumerableSetTransactor // Generic write-only contract binding to access the raw methods on
}

// NewEnumerableSet creates a new instance of EnumerableSet, bound to a specific deployed contract.
func NewEnumerableSet(address common.Address, backend bind.ContractBackend) (*EnumerableSet, error) {
	contract, err := bindEnumerableSet(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &EnumerableSet{EnumerableSetCaller: EnumerableSetCaller{contract: contract}, EnumerableSetTransactor: EnumerableSetTransactor{contract: contract}, EnumerableSetFilterer: EnumerableSetFilterer{contract: contract}}, nil
}

// NewEnumerableSetCaller creates a new read-only instance of EnumerableSet, bound to a specific deployed contract.
func NewEnumerableSetCaller(address common.Address, caller bind.ContractCaller) (*EnumerableSetCaller, error) {
	contract, err := bindEnumerableSet(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &EnumerableSetCaller{contract: contract}, nil
}

// NewEnumerableSetTransactor creates a new write-only instance of EnumerableSet, bound to a specific deployed contract.
func NewEnumerableSetTransactor(address common.Address, transactor bind.ContractTransactor) (*EnumerableSetTransactor, error) {
	contract, err := bindEnumerableSet(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &EnumerableSetTransactor{contract: contract}, nil
}

// NewEnumerableSetFilterer creates a new log filterer instance of EnumerableSet, bound to a specific deployed contract.
func NewEnumerableSetFilterer(address common.Address, filterer bind.ContractFilterer) (*EnumerableSetFilterer, error) {
	contract, err := bindEnumerableSet(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &EnumerableSetFilterer{contract: contract}, nil
}

// bindEnumerableSet binds a generic wrapper to an already deployed contract.
func bindEnumerableSet(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(EnumerableSetABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_EnumerableSet *EnumerableSetRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _EnumerableSet.Contract.EnumerableSetCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_EnumerableSet *EnumerableSetRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EnumerableSet.Contract.EnumerableSetTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_EnumerableSet *EnumerableSetRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _EnumerableSet.Contract.EnumerableSetTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_EnumerableSet *EnumerableSetCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _EnumerableSet.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_EnumerableSet *EnumerableSetTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EnumerableSet.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_EnumerableSet *EnumerableSetTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _EnumerableSet.Contract.contract.Transact(opts, method, params...)
}

// IACUMetaData contains all meta data concerning the IACU contract.
var IACUMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"getScaleFactor\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"setOperator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"oracle\",\"type\":\"address\"}],\"name\":\"setOracle\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"update\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"status\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"value\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"7f5e2f11": "getScaleFactor()",
		"b3ab15fb": "setOperator(address)",
		"7adbf973": "setOracle(address)",
		"a2e62045": "update()",
		"3fa4f245": "value()",
	},
}

// IACUABI is the input ABI used to generate the binding from.
// Deprecated: Use IACUMetaData.ABI instead.
var IACUABI = IACUMetaData.ABI

// Deprecated: Use IACUMetaData.Sigs instead.
// IACUFuncSigs maps the 4-byte function signature to its string representation.
var IACUFuncSigs = IACUMetaData.Sigs

// IACU is an auto generated Go binding around an Ethereum contract.
type IACU struct {
	IACUCaller     // Read-only binding to the contract
	IACUTransactor // Write-only binding to the contract
	IACUFilterer   // Log filterer for contract events
}

// IACUCaller is an auto generated read-only Go binding around an Ethereum contract.
type IACUCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IACUTransactor is an auto generated write-only Go binding around an Ethereum contract.
type IACUTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IACUFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type IACUFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IACUSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type IACUSession struct {
	Contract     *IACU             // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IACUCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type IACUCallerSession struct {
	Contract *IACUCaller   // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// IACUTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type IACUTransactorSession struct {
	Contract     *IACUTransactor   // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IACURaw is an auto generated low-level Go binding around an Ethereum contract.
type IACURaw struct {
	Contract *IACU // Generic contract binding to access the raw methods on
}

// IACUCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type IACUCallerRaw struct {
	Contract *IACUCaller // Generic read-only contract binding to access the raw methods on
}

// IACUTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type IACUTransactorRaw struct {
	Contract *IACUTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIACU creates a new instance of IACU, bound to a specific deployed contract.
func NewIACU(address common.Address, backend bind.ContractBackend) (*IACU, error) {
	contract, err := bindIACU(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IACU{IACUCaller: IACUCaller{contract: contract}, IACUTransactor: IACUTransactor{contract: contract}, IACUFilterer: IACUFilterer{contract: contract}}, nil
}

// NewIACUCaller creates a new read-only instance of IACU, bound to a specific deployed contract.
func NewIACUCaller(address common.Address, caller bind.ContractCaller) (*IACUCaller, error) {
	contract, err := bindIACU(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IACUCaller{contract: contract}, nil
}

// NewIACUTransactor creates a new write-only instance of IACU, bound to a specific deployed contract.
func NewIACUTransactor(address common.Address, transactor bind.ContractTransactor) (*IACUTransactor, error) {
	contract, err := bindIACU(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IACUTransactor{contract: contract}, nil
}

// NewIACUFilterer creates a new log filterer instance of IACU, bound to a specific deployed contract.
func NewIACUFilterer(address common.Address, filterer bind.ContractFilterer) (*IACUFilterer, error) {
	contract, err := bindIACU(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IACUFilterer{contract: contract}, nil
}

// bindIACU binds a generic wrapper to an already deployed contract.
func bindIACU(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(IACUABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IACU *IACURaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IACU.Contract.IACUCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IACU *IACURaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IACU.Contract.IACUTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IACU *IACURaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IACU.Contract.IACUTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IACU *IACUCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IACU.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IACU *IACUTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IACU.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IACU *IACUTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IACU.Contract.contract.Transact(opts, method, params...)
}

// GetScaleFactor is a free data retrieval call binding the contract method 0x7f5e2f11.
//
// Solidity: function getScaleFactor() view returns(uint256)
func (_IACU *IACUCaller) GetScaleFactor(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IACU.contract.Call(opts, &out, "getScaleFactor")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetScaleFactor is a free data retrieval call binding the contract method 0x7f5e2f11.
//
// Solidity: function getScaleFactor() view returns(uint256)
func (_IACU *IACUSession) GetScaleFactor() (*big.Int, error) {
	return _IACU.Contract.GetScaleFactor(&_IACU.CallOpts)
}

// GetScaleFactor is a free data retrieval call binding the contract method 0x7f5e2f11.
//
// Solidity: function getScaleFactor() view returns(uint256)
func (_IACU *IACUCallerSession) GetScaleFactor() (*big.Int, error) {
	return _IACU.Contract.GetScaleFactor(&_IACU.CallOpts)
}

// Value is a free data retrieval call binding the contract method 0x3fa4f245.
//
// Solidity: function value() view returns(uint256)
func (_IACU *IACUCaller) Value(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IACU.contract.Call(opts, &out, "value")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Value is a free data retrieval call binding the contract method 0x3fa4f245.
//
// Solidity: function value() view returns(uint256)
func (_IACU *IACUSession) Value() (*big.Int, error) {
	return _IACU.Contract.Value(&_IACU.CallOpts)
}

// Value is a free data retrieval call binding the contract method 0x3fa4f245.
//
// Solidity: function value() view returns(uint256)
func (_IACU *IACUCallerSession) Value() (*big.Int, error) {
	return _IACU.Contract.Value(&_IACU.CallOpts)
}

// SetOperator is a paid mutator transaction binding the contract method 0xb3ab15fb.
//
// Solidity: function setOperator(address operator) returns()
func (_IACU *IACUTransactor) SetOperator(opts *bind.TransactOpts, operator common.Address) (*types.Transaction, error) {
	return _IACU.contract.Transact(opts, "setOperator", operator)
}

// SetOperator is a paid mutator transaction binding the contract method 0xb3ab15fb.
//
// Solidity: function setOperator(address operator) returns()
func (_IACU *IACUSession) SetOperator(operator common.Address) (*types.Transaction, error) {
	return _IACU.Contract.SetOperator(&_IACU.TransactOpts, operator)
}

// SetOperator is a paid mutator transaction binding the contract method 0xb3ab15fb.
//
// Solidity: function setOperator(address operator) returns()
func (_IACU *IACUTransactorSession) SetOperator(operator common.Address) (*types.Transaction, error) {
	return _IACU.Contract.SetOperator(&_IACU.TransactOpts, operator)
}

// SetOracle is a paid mutator transaction binding the contract method 0x7adbf973.
//
// Solidity: function setOracle(address oracle) returns()
func (_IACU *IACUTransactor) SetOracle(opts *bind.TransactOpts, oracle common.Address) (*types.Transaction, error) {
	return _IACU.contract.Transact(opts, "setOracle", oracle)
}

// SetOracle is a paid mutator transaction binding the contract method 0x7adbf973.
//
// Solidity: function setOracle(address oracle) returns()
func (_IACU *IACUSession) SetOracle(oracle common.Address) (*types.Transaction, error) {
	return _IACU.Contract.SetOracle(&_IACU.TransactOpts, oracle)
}

// SetOracle is a paid mutator transaction binding the contract method 0x7adbf973.
//
// Solidity: function setOracle(address oracle) returns()
func (_IACU *IACUTransactorSession) SetOracle(oracle common.Address) (*types.Transaction, error) {
	return _IACU.Contract.SetOracle(&_IACU.TransactOpts, oracle)
}

// Update is a paid mutator transaction binding the contract method 0xa2e62045.
//
// Solidity: function update() returns(bool status)
func (_IACU *IACUTransactor) Update(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IACU.contract.Transact(opts, "update")
}

// Update is a paid mutator transaction binding the contract method 0xa2e62045.
//
// Solidity: function update() returns(bool status)
func (_IACU *IACUSession) Update() (*types.Transaction, error) {
	return _IACU.Contract.Update(&_IACU.TransactOpts)
}

// Update is a paid mutator transaction binding the contract method 0xa2e62045.
//
// Solidity: function update() returns(bool status)
func (_IACU *IACUTransactorSession) Update() (*types.Transaction, error) {
	return _IACU.Contract.Update(&_IACU.TransactOpts)
}

// IAccountabilityMetaData contains all meta data concerning the IAccountability contract.
var IAccountabilityMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"collusion\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"history\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"jail\",\"type\":\"uint256\"}],\"indexed\":false,\"internalType\":\"structIAccountability.Factors\",\"name\":\"oldFactors\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"collusion\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"history\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"jail\",\"type\":\"uint256\"}],\"indexed\":false,\"internalType\":\"structIAccountability.Factors\",\"name\":\"newFactors\",\"type\":\"tuple\"}],\"name\":\"AccountabilityFactorsUpdate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"low\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"mid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"high\",\"type\":\"uint256\"}],\"indexed\":false,\"internalType\":\"structIAccountability.BaseSlashingRates\",\"name\":\"oldRates\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"low\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"mid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"high\",\"type\":\"uint256\"}],\"indexed\":false,\"internalType\":\"structIAccountability.BaseSlashingRates\",\"name\":\"newRates\",\"type\":\"tuple\"}],\"name\":\"BaseSlashingRateUpdate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_offender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_id\",\"type\":\"uint256\"}],\"name\":\"InnocenceProven\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_offender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_severity\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_id\",\"type\":\"uint256\"}],\"name\":\"NewAccusation\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_offender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_severity\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_epoch\",\"type\":\"uint256\"}],\"name\":\"NewFaultProof\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"_reporter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_offender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_ntnReward\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_atnReward\",\"type\":\"uint256\"}],\"name\":\"ReporterRewarded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"releaseBlock\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"isJailbound\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"eventId\",\"type\":\"uint256\"}],\"name\":\"SlashingEvent\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_validator\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_ntnReward\",\"type\":\"uint256\"}],\"name\":\"distributeRewards\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"_epochEnd\",\"type\":\"bool\"}],\"name\":\"finalize\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getConfig\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"innocenceProofSubmissionWindow\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"delta\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"range\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"low\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"mid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"high\",\"type\":\"uint256\"}],\"internalType\":\"structIAccountability.BaseSlashingRates\",\"name\":\"baseSlashingRates\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"collusion\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"history\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"jail\",\"type\":\"uint256\"}],\"internalType\":\"structIAccountability.Factors\",\"name\":\"factors\",\"type\":\"tuple\"}],\"internalType\":\"structIAccountability.Config\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getGracePeriod\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_committee\",\"type\":\"address[]\"}],\"name\":\"setCommittee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"a8031a1d": "distributeRewards(address,uint256)",
		"6c9789b0": "finalize(bool)",
		"c3f909d4": "getConfig()",
		"dbd18388": "getGracePeriod()",
		"e08b14ed": "setCommittee(address[])",
	},
}

// IAccountabilityABI is the input ABI used to generate the binding from.
// Deprecated: Use IAccountabilityMetaData.ABI instead.
var IAccountabilityABI = IAccountabilityMetaData.ABI

// Deprecated: Use IAccountabilityMetaData.Sigs instead.
// IAccountabilityFuncSigs maps the 4-byte function signature to its string representation.
var IAccountabilityFuncSigs = IAccountabilityMetaData.Sigs

// IAccountability is an auto generated Go binding around an Ethereum contract.
type IAccountability struct {
	IAccountabilityCaller     // Read-only binding to the contract
	IAccountabilityTransactor // Write-only binding to the contract
	IAccountabilityFilterer   // Log filterer for contract events
}

// IAccountabilityCaller is an auto generated read-only Go binding around an Ethereum contract.
type IAccountabilityCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IAccountabilityTransactor is an auto generated write-only Go binding around an Ethereum contract.
type IAccountabilityTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IAccountabilityFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type IAccountabilityFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IAccountabilitySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type IAccountabilitySession struct {
	Contract     *IAccountability  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IAccountabilityCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type IAccountabilityCallerSession struct {
	Contract *IAccountabilityCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// IAccountabilityTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type IAccountabilityTransactorSession struct {
	Contract     *IAccountabilityTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// IAccountabilityRaw is an auto generated low-level Go binding around an Ethereum contract.
type IAccountabilityRaw struct {
	Contract *IAccountability // Generic contract binding to access the raw methods on
}

// IAccountabilityCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type IAccountabilityCallerRaw struct {
	Contract *IAccountabilityCaller // Generic read-only contract binding to access the raw methods on
}

// IAccountabilityTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type IAccountabilityTransactorRaw struct {
	Contract *IAccountabilityTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIAccountability creates a new instance of IAccountability, bound to a specific deployed contract.
func NewIAccountability(address common.Address, backend bind.ContractBackend) (*IAccountability, error) {
	contract, err := bindIAccountability(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IAccountability{IAccountabilityCaller: IAccountabilityCaller{contract: contract}, IAccountabilityTransactor: IAccountabilityTransactor{contract: contract}, IAccountabilityFilterer: IAccountabilityFilterer{contract: contract}}, nil
}

// NewIAccountabilityCaller creates a new read-only instance of IAccountability, bound to a specific deployed contract.
func NewIAccountabilityCaller(address common.Address, caller bind.ContractCaller) (*IAccountabilityCaller, error) {
	contract, err := bindIAccountability(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IAccountabilityCaller{contract: contract}, nil
}

// NewIAccountabilityTransactor creates a new write-only instance of IAccountability, bound to a specific deployed contract.
func NewIAccountabilityTransactor(address common.Address, transactor bind.ContractTransactor) (*IAccountabilityTransactor, error) {
	contract, err := bindIAccountability(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IAccountabilityTransactor{contract: contract}, nil
}

// NewIAccountabilityFilterer creates a new log filterer instance of IAccountability, bound to a specific deployed contract.
func NewIAccountabilityFilterer(address common.Address, filterer bind.ContractFilterer) (*IAccountabilityFilterer, error) {
	contract, err := bindIAccountability(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IAccountabilityFilterer{contract: contract}, nil
}

// bindIAccountability binds a generic wrapper to an already deployed contract.
func bindIAccountability(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(IAccountabilityABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IAccountability *IAccountabilityRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IAccountability.Contract.IAccountabilityCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IAccountability *IAccountabilityRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IAccountability.Contract.IAccountabilityTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IAccountability *IAccountabilityRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IAccountability.Contract.IAccountabilityTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IAccountability *IAccountabilityCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IAccountability.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IAccountability *IAccountabilityTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IAccountability.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IAccountability *IAccountabilityTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IAccountability.Contract.contract.Transact(opts, method, params...)
}

// GetConfig is a free data retrieval call binding the contract method 0xc3f909d4.
//
// Solidity: function getConfig() view returns((uint256,uint256,uint256,(uint256,uint256,uint256),(uint256,uint256,uint256)))
func (_IAccountability *IAccountabilityCaller) GetConfig(opts *bind.CallOpts) (IAccountabilityConfig, error) {
	var out []interface{}
	err := _IAccountability.contract.Call(opts, &out, "getConfig")

	if err != nil {
		return *new(IAccountabilityConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(IAccountabilityConfig)).(*IAccountabilityConfig)

	return out0, err

}

// GetConfig is a free data retrieval call binding the contract method 0xc3f909d4.
//
// Solidity: function getConfig() view returns((uint256,uint256,uint256,(uint256,uint256,uint256),(uint256,uint256,uint256)))
func (_IAccountability *IAccountabilitySession) GetConfig() (IAccountabilityConfig, error) {
	return _IAccountability.Contract.GetConfig(&_IAccountability.CallOpts)
}

// GetConfig is a free data retrieval call binding the contract method 0xc3f909d4.
//
// Solidity: function getConfig() view returns((uint256,uint256,uint256,(uint256,uint256,uint256),(uint256,uint256,uint256)))
func (_IAccountability *IAccountabilityCallerSession) GetConfig() (IAccountabilityConfig, error) {
	return _IAccountability.Contract.GetConfig(&_IAccountability.CallOpts)
}

// GetGracePeriod is a free data retrieval call binding the contract method 0xdbd18388.
//
// Solidity: function getGracePeriod() view returns(uint256)
func (_IAccountability *IAccountabilityCaller) GetGracePeriod(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IAccountability.contract.Call(opts, &out, "getGracePeriod")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetGracePeriod is a free data retrieval call binding the contract method 0xdbd18388.
//
// Solidity: function getGracePeriod() view returns(uint256)
func (_IAccountability *IAccountabilitySession) GetGracePeriod() (*big.Int, error) {
	return _IAccountability.Contract.GetGracePeriod(&_IAccountability.CallOpts)
}

// GetGracePeriod is a free data retrieval call binding the contract method 0xdbd18388.
//
// Solidity: function getGracePeriod() view returns(uint256)
func (_IAccountability *IAccountabilityCallerSession) GetGracePeriod() (*big.Int, error) {
	return _IAccountability.Contract.GetGracePeriod(&_IAccountability.CallOpts)
}

// DistributeRewards is a paid mutator transaction binding the contract method 0xa8031a1d.
//
// Solidity: function distributeRewards(address _validator, uint256 _ntnReward) payable returns()
func (_IAccountability *IAccountabilityTransactor) DistributeRewards(opts *bind.TransactOpts, _validator common.Address, _ntnReward *big.Int) (*types.Transaction, error) {
	return _IAccountability.contract.Transact(opts, "distributeRewards", _validator, _ntnReward)
}

// DistributeRewards is a paid mutator transaction binding the contract method 0xa8031a1d.
//
// Solidity: function distributeRewards(address _validator, uint256 _ntnReward) payable returns()
func (_IAccountability *IAccountabilitySession) DistributeRewards(_validator common.Address, _ntnReward *big.Int) (*types.Transaction, error) {
	return _IAccountability.Contract.DistributeRewards(&_IAccountability.TransactOpts, _validator, _ntnReward)
}

// DistributeRewards is a paid mutator transaction binding the contract method 0xa8031a1d.
//
// Solidity: function distributeRewards(address _validator, uint256 _ntnReward) payable returns()
func (_IAccountability *IAccountabilityTransactorSession) DistributeRewards(_validator common.Address, _ntnReward *big.Int) (*types.Transaction, error) {
	return _IAccountability.Contract.DistributeRewards(&_IAccountability.TransactOpts, _validator, _ntnReward)
}

// Finalize is a paid mutator transaction binding the contract method 0x6c9789b0.
//
// Solidity: function finalize(bool _epochEnd) returns(uint256, uint256, uint256)
func (_IAccountability *IAccountabilityTransactor) Finalize(opts *bind.TransactOpts, _epochEnd bool) (*types.Transaction, error) {
	return _IAccountability.contract.Transact(opts, "finalize", _epochEnd)
}

// Finalize is a paid mutator transaction binding the contract method 0x6c9789b0.
//
// Solidity: function finalize(bool _epochEnd) returns(uint256, uint256, uint256)
func (_IAccountability *IAccountabilitySession) Finalize(_epochEnd bool) (*types.Transaction, error) {
	return _IAccountability.Contract.Finalize(&_IAccountability.TransactOpts, _epochEnd)
}

// Finalize is a paid mutator transaction binding the contract method 0x6c9789b0.
//
// Solidity: function finalize(bool _epochEnd) returns(uint256, uint256, uint256)
func (_IAccountability *IAccountabilityTransactorSession) Finalize(_epochEnd bool) (*types.Transaction, error) {
	return _IAccountability.Contract.Finalize(&_IAccountability.TransactOpts, _epochEnd)
}

// SetCommittee is a paid mutator transaction binding the contract method 0xe08b14ed.
//
// Solidity: function setCommittee(address[] _committee) returns()
func (_IAccountability *IAccountabilityTransactor) SetCommittee(opts *bind.TransactOpts, _committee []common.Address) (*types.Transaction, error) {
	return _IAccountability.contract.Transact(opts, "setCommittee", _committee)
}

// SetCommittee is a paid mutator transaction binding the contract method 0xe08b14ed.
//
// Solidity: function setCommittee(address[] _committee) returns()
func (_IAccountability *IAccountabilitySession) SetCommittee(_committee []common.Address) (*types.Transaction, error) {
	return _IAccountability.Contract.SetCommittee(&_IAccountability.TransactOpts, _committee)
}

// SetCommittee is a paid mutator transaction binding the contract method 0xe08b14ed.
//
// Solidity: function setCommittee(address[] _committee) returns()
func (_IAccountability *IAccountabilityTransactorSession) SetCommittee(_committee []common.Address) (*types.Transaction, error) {
	return _IAccountability.Contract.SetCommittee(&_IAccountability.TransactOpts, _committee)
}

// IAccountabilityAccountabilityFactorsUpdateIterator is returned from FilterAccountabilityFactorsUpdate and is used to iterate over the raw logs and unpacked data for AccountabilityFactorsUpdate events raised by the IAccountability contract.
type IAccountabilityAccountabilityFactorsUpdateIterator struct {
	Event *IAccountabilityAccountabilityFactorsUpdate // Event containing the contract specifics and raw log

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
func (it *IAccountabilityAccountabilityFactorsUpdateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAccountabilityAccountabilityFactorsUpdate)
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
		it.Event = new(IAccountabilityAccountabilityFactorsUpdate)
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
func (it *IAccountabilityAccountabilityFactorsUpdateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IAccountabilityAccountabilityFactorsUpdateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IAccountabilityAccountabilityFactorsUpdate represents a AccountabilityFactorsUpdate event raised by the IAccountability contract.
type IAccountabilityAccountabilityFactorsUpdate struct {
	OldFactors IAccountabilityFactors
	NewFactors IAccountabilityFactors
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterAccountabilityFactorsUpdate is a free log retrieval operation binding the contract event 0xeda1502b5625ce16270e5a5f5561c68dfecaef2c5eacc1cc9521fad59c34d82f.
//
// Solidity: event AccountabilityFactorsUpdate((uint256,uint256,uint256) oldFactors, (uint256,uint256,uint256) newFactors)
func (_IAccountability *IAccountabilityFilterer) FilterAccountabilityFactorsUpdate(opts *bind.FilterOpts) (*IAccountabilityAccountabilityFactorsUpdateIterator, error) {

	logs, sub, err := _IAccountability.contract.FilterLogs(opts, "AccountabilityFactorsUpdate")
	if err != nil {
		return nil, err
	}
	return &IAccountabilityAccountabilityFactorsUpdateIterator{contract: _IAccountability.contract, event: "AccountabilityFactorsUpdate", logs: logs, sub: sub}, nil
}

// WatchAccountabilityFactorsUpdate is a free log subscription operation binding the contract event 0xeda1502b5625ce16270e5a5f5561c68dfecaef2c5eacc1cc9521fad59c34d82f.
//
// Solidity: event AccountabilityFactorsUpdate((uint256,uint256,uint256) oldFactors, (uint256,uint256,uint256) newFactors)
func (_IAccountability *IAccountabilityFilterer) WatchAccountabilityFactorsUpdate(opts *bind.WatchOpts, sink chan<- *IAccountabilityAccountabilityFactorsUpdate) (event.Subscription, error) {

	logs, sub, err := _IAccountability.contract.WatchLogs(opts, "AccountabilityFactorsUpdate")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IAccountabilityAccountabilityFactorsUpdate)
				if err := _IAccountability.contract.UnpackLog(event, "AccountabilityFactorsUpdate", log); err != nil {
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

// ParseAccountabilityFactorsUpdate is a log parse operation binding the contract event 0xeda1502b5625ce16270e5a5f5561c68dfecaef2c5eacc1cc9521fad59c34d82f.
//
// Solidity: event AccountabilityFactorsUpdate((uint256,uint256,uint256) oldFactors, (uint256,uint256,uint256) newFactors)
func (_IAccountability *IAccountabilityFilterer) ParseAccountabilityFactorsUpdate(log types.Log) (*IAccountabilityAccountabilityFactorsUpdate, error) {
	event := new(IAccountabilityAccountabilityFactorsUpdate)
	if err := _IAccountability.contract.UnpackLog(event, "AccountabilityFactorsUpdate", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IAccountabilityBaseSlashingRateUpdateIterator is returned from FilterBaseSlashingRateUpdate and is used to iterate over the raw logs and unpacked data for BaseSlashingRateUpdate events raised by the IAccountability contract.
type IAccountabilityBaseSlashingRateUpdateIterator struct {
	Event *IAccountabilityBaseSlashingRateUpdate // Event containing the contract specifics and raw log

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
func (it *IAccountabilityBaseSlashingRateUpdateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAccountabilityBaseSlashingRateUpdate)
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
		it.Event = new(IAccountabilityBaseSlashingRateUpdate)
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
func (it *IAccountabilityBaseSlashingRateUpdateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IAccountabilityBaseSlashingRateUpdateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IAccountabilityBaseSlashingRateUpdate represents a BaseSlashingRateUpdate event raised by the IAccountability contract.
type IAccountabilityBaseSlashingRateUpdate struct {
	OldRates IAccountabilityBaseSlashingRates
	NewRates IAccountabilityBaseSlashingRates
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterBaseSlashingRateUpdate is a free log retrieval operation binding the contract event 0x6fea73568c60fe418d5e4d1a9fde6dce66d8816ffcebe91449e9bf1098a5ce1a.
//
// Solidity: event BaseSlashingRateUpdate((uint256,uint256,uint256) oldRates, (uint256,uint256,uint256) newRates)
func (_IAccountability *IAccountabilityFilterer) FilterBaseSlashingRateUpdate(opts *bind.FilterOpts) (*IAccountabilityBaseSlashingRateUpdateIterator, error) {

	logs, sub, err := _IAccountability.contract.FilterLogs(opts, "BaseSlashingRateUpdate")
	if err != nil {
		return nil, err
	}
	return &IAccountabilityBaseSlashingRateUpdateIterator{contract: _IAccountability.contract, event: "BaseSlashingRateUpdate", logs: logs, sub: sub}, nil
}

// WatchBaseSlashingRateUpdate is a free log subscription operation binding the contract event 0x6fea73568c60fe418d5e4d1a9fde6dce66d8816ffcebe91449e9bf1098a5ce1a.
//
// Solidity: event BaseSlashingRateUpdate((uint256,uint256,uint256) oldRates, (uint256,uint256,uint256) newRates)
func (_IAccountability *IAccountabilityFilterer) WatchBaseSlashingRateUpdate(opts *bind.WatchOpts, sink chan<- *IAccountabilityBaseSlashingRateUpdate) (event.Subscription, error) {

	logs, sub, err := _IAccountability.contract.WatchLogs(opts, "BaseSlashingRateUpdate")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IAccountabilityBaseSlashingRateUpdate)
				if err := _IAccountability.contract.UnpackLog(event, "BaseSlashingRateUpdate", log); err != nil {
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

// ParseBaseSlashingRateUpdate is a log parse operation binding the contract event 0x6fea73568c60fe418d5e4d1a9fde6dce66d8816ffcebe91449e9bf1098a5ce1a.
//
// Solidity: event BaseSlashingRateUpdate((uint256,uint256,uint256) oldRates, (uint256,uint256,uint256) newRates)
func (_IAccountability *IAccountabilityFilterer) ParseBaseSlashingRateUpdate(log types.Log) (*IAccountabilityBaseSlashingRateUpdate, error) {
	event := new(IAccountabilityBaseSlashingRateUpdate)
	if err := _IAccountability.contract.UnpackLog(event, "BaseSlashingRateUpdate", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IAccountabilityInnocenceProvenIterator is returned from FilterInnocenceProven and is used to iterate over the raw logs and unpacked data for InnocenceProven events raised by the IAccountability contract.
type IAccountabilityInnocenceProvenIterator struct {
	Event *IAccountabilityInnocenceProven // Event containing the contract specifics and raw log

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
func (it *IAccountabilityInnocenceProvenIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAccountabilityInnocenceProven)
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
		it.Event = new(IAccountabilityInnocenceProven)
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
func (it *IAccountabilityInnocenceProvenIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IAccountabilityInnocenceProvenIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IAccountabilityInnocenceProven represents a InnocenceProven event raised by the IAccountability contract.
type IAccountabilityInnocenceProven struct {
	Offender common.Address
	Id       *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterInnocenceProven is a free log retrieval operation binding the contract event 0x1fa96beb8dddcb7d4484dd00c4059e872439f7a474a2ecf49c430fc6e86c9e1f.
//
// Solidity: event InnocenceProven(address indexed _offender, uint256 _id)
func (_IAccountability *IAccountabilityFilterer) FilterInnocenceProven(opts *bind.FilterOpts, _offender []common.Address) (*IAccountabilityInnocenceProvenIterator, error) {

	var _offenderRule []interface{}
	for _, _offenderItem := range _offender {
		_offenderRule = append(_offenderRule, _offenderItem)
	}

	logs, sub, err := _IAccountability.contract.FilterLogs(opts, "InnocenceProven", _offenderRule)
	if err != nil {
		return nil, err
	}
	return &IAccountabilityInnocenceProvenIterator{contract: _IAccountability.contract, event: "InnocenceProven", logs: logs, sub: sub}, nil
}

// WatchInnocenceProven is a free log subscription operation binding the contract event 0x1fa96beb8dddcb7d4484dd00c4059e872439f7a474a2ecf49c430fc6e86c9e1f.
//
// Solidity: event InnocenceProven(address indexed _offender, uint256 _id)
func (_IAccountability *IAccountabilityFilterer) WatchInnocenceProven(opts *bind.WatchOpts, sink chan<- *IAccountabilityInnocenceProven, _offender []common.Address) (event.Subscription, error) {

	var _offenderRule []interface{}
	for _, _offenderItem := range _offender {
		_offenderRule = append(_offenderRule, _offenderItem)
	}

	logs, sub, err := _IAccountability.contract.WatchLogs(opts, "InnocenceProven", _offenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IAccountabilityInnocenceProven)
				if err := _IAccountability.contract.UnpackLog(event, "InnocenceProven", log); err != nil {
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

// ParseInnocenceProven is a log parse operation binding the contract event 0x1fa96beb8dddcb7d4484dd00c4059e872439f7a474a2ecf49c430fc6e86c9e1f.
//
// Solidity: event InnocenceProven(address indexed _offender, uint256 _id)
func (_IAccountability *IAccountabilityFilterer) ParseInnocenceProven(log types.Log) (*IAccountabilityInnocenceProven, error) {
	event := new(IAccountabilityInnocenceProven)
	if err := _IAccountability.contract.UnpackLog(event, "InnocenceProven", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IAccountabilityNewAccusationIterator is returned from FilterNewAccusation and is used to iterate over the raw logs and unpacked data for NewAccusation events raised by the IAccountability contract.
type IAccountabilityNewAccusationIterator struct {
	Event *IAccountabilityNewAccusation // Event containing the contract specifics and raw log

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
func (it *IAccountabilityNewAccusationIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAccountabilityNewAccusation)
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
		it.Event = new(IAccountabilityNewAccusation)
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
func (it *IAccountabilityNewAccusationIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IAccountabilityNewAccusationIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IAccountabilityNewAccusation represents a NewAccusation event raised by the IAccountability contract.
type IAccountabilityNewAccusation struct {
	Offender common.Address
	Severity *big.Int
	Id       *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterNewAccusation is a free log retrieval operation binding the contract event 0x2e8e354b41470731dafa7c3df150e9498a8d5b9c51ff0259fbf77f721ba40351.
//
// Solidity: event NewAccusation(address indexed _offender, uint256 _severity, uint256 _id)
func (_IAccountability *IAccountabilityFilterer) FilterNewAccusation(opts *bind.FilterOpts, _offender []common.Address) (*IAccountabilityNewAccusationIterator, error) {

	var _offenderRule []interface{}
	for _, _offenderItem := range _offender {
		_offenderRule = append(_offenderRule, _offenderItem)
	}

	logs, sub, err := _IAccountability.contract.FilterLogs(opts, "NewAccusation", _offenderRule)
	if err != nil {
		return nil, err
	}
	return &IAccountabilityNewAccusationIterator{contract: _IAccountability.contract, event: "NewAccusation", logs: logs, sub: sub}, nil
}

// WatchNewAccusation is a free log subscription operation binding the contract event 0x2e8e354b41470731dafa7c3df150e9498a8d5b9c51ff0259fbf77f721ba40351.
//
// Solidity: event NewAccusation(address indexed _offender, uint256 _severity, uint256 _id)
func (_IAccountability *IAccountabilityFilterer) WatchNewAccusation(opts *bind.WatchOpts, sink chan<- *IAccountabilityNewAccusation, _offender []common.Address) (event.Subscription, error) {

	var _offenderRule []interface{}
	for _, _offenderItem := range _offender {
		_offenderRule = append(_offenderRule, _offenderItem)
	}

	logs, sub, err := _IAccountability.contract.WatchLogs(opts, "NewAccusation", _offenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IAccountabilityNewAccusation)
				if err := _IAccountability.contract.UnpackLog(event, "NewAccusation", log); err != nil {
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

// ParseNewAccusation is a log parse operation binding the contract event 0x2e8e354b41470731dafa7c3df150e9498a8d5b9c51ff0259fbf77f721ba40351.
//
// Solidity: event NewAccusation(address indexed _offender, uint256 _severity, uint256 _id)
func (_IAccountability *IAccountabilityFilterer) ParseNewAccusation(log types.Log) (*IAccountabilityNewAccusation, error) {
	event := new(IAccountabilityNewAccusation)
	if err := _IAccountability.contract.UnpackLog(event, "NewAccusation", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IAccountabilityNewFaultProofIterator is returned from FilterNewFaultProof and is used to iterate over the raw logs and unpacked data for NewFaultProof events raised by the IAccountability contract.
type IAccountabilityNewFaultProofIterator struct {
	Event *IAccountabilityNewFaultProof // Event containing the contract specifics and raw log

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
func (it *IAccountabilityNewFaultProofIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAccountabilityNewFaultProof)
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
		it.Event = new(IAccountabilityNewFaultProof)
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
func (it *IAccountabilityNewFaultProofIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IAccountabilityNewFaultProofIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IAccountabilityNewFaultProof represents a NewFaultProof event raised by the IAccountability contract.
type IAccountabilityNewFaultProof struct {
	Offender common.Address
	Severity *big.Int
	Id       *big.Int
	Epoch    *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterNewFaultProof is a free log retrieval operation binding the contract event 0x5fd9605880705541b88eb5df56222ffa3c3e6884010fcf26d4c2c372917d98d7.
//
// Solidity: event NewFaultProof(address indexed _offender, uint256 _severity, uint256 _id, uint256 _epoch)
func (_IAccountability *IAccountabilityFilterer) FilterNewFaultProof(opts *bind.FilterOpts, _offender []common.Address) (*IAccountabilityNewFaultProofIterator, error) {

	var _offenderRule []interface{}
	for _, _offenderItem := range _offender {
		_offenderRule = append(_offenderRule, _offenderItem)
	}

	logs, sub, err := _IAccountability.contract.FilterLogs(opts, "NewFaultProof", _offenderRule)
	if err != nil {
		return nil, err
	}
	return &IAccountabilityNewFaultProofIterator{contract: _IAccountability.contract, event: "NewFaultProof", logs: logs, sub: sub}, nil
}

// WatchNewFaultProof is a free log subscription operation binding the contract event 0x5fd9605880705541b88eb5df56222ffa3c3e6884010fcf26d4c2c372917d98d7.
//
// Solidity: event NewFaultProof(address indexed _offender, uint256 _severity, uint256 _id, uint256 _epoch)
func (_IAccountability *IAccountabilityFilterer) WatchNewFaultProof(opts *bind.WatchOpts, sink chan<- *IAccountabilityNewFaultProof, _offender []common.Address) (event.Subscription, error) {

	var _offenderRule []interface{}
	for _, _offenderItem := range _offender {
		_offenderRule = append(_offenderRule, _offenderItem)
	}

	logs, sub, err := _IAccountability.contract.WatchLogs(opts, "NewFaultProof", _offenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IAccountabilityNewFaultProof)
				if err := _IAccountability.contract.UnpackLog(event, "NewFaultProof", log); err != nil {
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

// ParseNewFaultProof is a log parse operation binding the contract event 0x5fd9605880705541b88eb5df56222ffa3c3e6884010fcf26d4c2c372917d98d7.
//
// Solidity: event NewFaultProof(address indexed _offender, uint256 _severity, uint256 _id, uint256 _epoch)
func (_IAccountability *IAccountabilityFilterer) ParseNewFaultProof(log types.Log) (*IAccountabilityNewFaultProof, error) {
	event := new(IAccountabilityNewFaultProof)
	if err := _IAccountability.contract.UnpackLog(event, "NewFaultProof", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IAccountabilityReporterRewardedIterator is returned from FilterReporterRewarded and is used to iterate over the raw logs and unpacked data for ReporterRewarded events raised by the IAccountability contract.
type IAccountabilityReporterRewardedIterator struct {
	Event *IAccountabilityReporterRewarded // Event containing the contract specifics and raw log

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
func (it *IAccountabilityReporterRewardedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAccountabilityReporterRewarded)
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
		it.Event = new(IAccountabilityReporterRewarded)
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
func (it *IAccountabilityReporterRewardedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IAccountabilityReporterRewardedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IAccountabilityReporterRewarded represents a ReporterRewarded event raised by the IAccountability contract.
type IAccountabilityReporterRewarded struct {
	Reporter  common.Address
	Offender  common.Address
	NtnReward *big.Int
	AtnReward *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterReporterRewarded is a free log retrieval operation binding the contract event 0xd7bd322f319b32cbf0c40184f51d4204b4e08b1c95ca79f470747483881840e9.
//
// Solidity: event ReporterRewarded(address _reporter, address indexed _offender, uint256 _ntnReward, uint256 _atnReward)
func (_IAccountability *IAccountabilityFilterer) FilterReporterRewarded(opts *bind.FilterOpts, _offender []common.Address) (*IAccountabilityReporterRewardedIterator, error) {

	var _offenderRule []interface{}
	for _, _offenderItem := range _offender {
		_offenderRule = append(_offenderRule, _offenderItem)
	}

	logs, sub, err := _IAccountability.contract.FilterLogs(opts, "ReporterRewarded", _offenderRule)
	if err != nil {
		return nil, err
	}
	return &IAccountabilityReporterRewardedIterator{contract: _IAccountability.contract, event: "ReporterRewarded", logs: logs, sub: sub}, nil
}

// WatchReporterRewarded is a free log subscription operation binding the contract event 0xd7bd322f319b32cbf0c40184f51d4204b4e08b1c95ca79f470747483881840e9.
//
// Solidity: event ReporterRewarded(address _reporter, address indexed _offender, uint256 _ntnReward, uint256 _atnReward)
func (_IAccountability *IAccountabilityFilterer) WatchReporterRewarded(opts *bind.WatchOpts, sink chan<- *IAccountabilityReporterRewarded, _offender []common.Address) (event.Subscription, error) {

	var _offenderRule []interface{}
	for _, _offenderItem := range _offender {
		_offenderRule = append(_offenderRule, _offenderItem)
	}

	logs, sub, err := _IAccountability.contract.WatchLogs(opts, "ReporterRewarded", _offenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IAccountabilityReporterRewarded)
				if err := _IAccountability.contract.UnpackLog(event, "ReporterRewarded", log); err != nil {
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

// ParseReporterRewarded is a log parse operation binding the contract event 0xd7bd322f319b32cbf0c40184f51d4204b4e08b1c95ca79f470747483881840e9.
//
// Solidity: event ReporterRewarded(address _reporter, address indexed _offender, uint256 _ntnReward, uint256 _atnReward)
func (_IAccountability *IAccountabilityFilterer) ParseReporterRewarded(log types.Log) (*IAccountabilityReporterRewarded, error) {
	event := new(IAccountabilityReporterRewarded)
	if err := _IAccountability.contract.UnpackLog(event, "ReporterRewarded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IAccountabilitySlashingEventIterator is returned from FilterSlashingEvent and is used to iterate over the raw logs and unpacked data for SlashingEvent events raised by the IAccountability contract.
type IAccountabilitySlashingEventIterator struct {
	Event *IAccountabilitySlashingEvent // Event containing the contract specifics and raw log

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
func (it *IAccountabilitySlashingEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAccountabilitySlashingEvent)
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
		it.Event = new(IAccountabilitySlashingEvent)
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
func (it *IAccountabilitySlashingEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IAccountabilitySlashingEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IAccountabilitySlashingEvent represents a SlashingEvent event raised by the IAccountability contract.
type IAccountabilitySlashingEvent struct {
	Validator    common.Address
	Amount       *big.Int
	ReleaseBlock *big.Int
	IsJailbound  bool
	EventId      *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterSlashingEvent is a free log retrieval operation binding the contract event 0x6617e612ea2d01b5a235997fa4963b56b1097df6f968a82972433e9ff852e0f9.
//
// Solidity: event SlashingEvent(address validator, uint256 amount, uint256 releaseBlock, bool isJailbound, uint256 eventId)
func (_IAccountability *IAccountabilityFilterer) FilterSlashingEvent(opts *bind.FilterOpts) (*IAccountabilitySlashingEventIterator, error) {

	logs, sub, err := _IAccountability.contract.FilterLogs(opts, "SlashingEvent")
	if err != nil {
		return nil, err
	}
	return &IAccountabilitySlashingEventIterator{contract: _IAccountability.contract, event: "SlashingEvent", logs: logs, sub: sub}, nil
}

// WatchSlashingEvent is a free log subscription operation binding the contract event 0x6617e612ea2d01b5a235997fa4963b56b1097df6f968a82972433e9ff852e0f9.
//
// Solidity: event SlashingEvent(address validator, uint256 amount, uint256 releaseBlock, bool isJailbound, uint256 eventId)
func (_IAccountability *IAccountabilityFilterer) WatchSlashingEvent(opts *bind.WatchOpts, sink chan<- *IAccountabilitySlashingEvent) (event.Subscription, error) {

	logs, sub, err := _IAccountability.contract.WatchLogs(opts, "SlashingEvent")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IAccountabilitySlashingEvent)
				if err := _IAccountability.contract.UnpackLog(event, "SlashingEvent", log); err != nil {
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

// ParseSlashingEvent is a log parse operation binding the contract event 0x6617e612ea2d01b5a235997fa4963b56b1097df6f968a82972433e9ff852e0f9.
//
// Solidity: event SlashingEvent(address validator, uint256 amount, uint256 releaseBlock, bool isJailbound, uint256 eventId)
func (_IAccountability *IAccountabilityFilterer) ParseSlashingEvent(log types.Log) (*IAccountabilitySlashingEvent, error) {
	event := new(IAccountabilitySlashingEvent)
	if err := _IAccountability.contract.UnpackLog(event, "SlashingEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IAuctioneerMetaData contains all meta data concerning the IAuctioneer contract.
var IAuctioneerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"paidInterest\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"setOperator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"oracle\",\"type\":\"address\"}],\"name\":\"setOracle\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"stabilization\",\"type\":\"address\"}],\"name\":\"setStabilization\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"96e4547c": "paidInterest()",
		"b3ab15fb": "setOperator(address)",
		"7adbf973": "setOracle(address)",
		"4f505895": "setStabilization(address)",
	},
}

// IAuctioneerABI is the input ABI used to generate the binding from.
// Deprecated: Use IAuctioneerMetaData.ABI instead.
var IAuctioneerABI = IAuctioneerMetaData.ABI

// Deprecated: Use IAuctioneerMetaData.Sigs instead.
// IAuctioneerFuncSigs maps the 4-byte function signature to its string representation.
var IAuctioneerFuncSigs = IAuctioneerMetaData.Sigs

// IAuctioneer is an auto generated Go binding around an Ethereum contract.
type IAuctioneer struct {
	IAuctioneerCaller     // Read-only binding to the contract
	IAuctioneerTransactor // Write-only binding to the contract
	IAuctioneerFilterer   // Log filterer for contract events
}

// IAuctioneerCaller is an auto generated read-only Go binding around an Ethereum contract.
type IAuctioneerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IAuctioneerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type IAuctioneerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IAuctioneerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type IAuctioneerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IAuctioneerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type IAuctioneerSession struct {
	Contract     *IAuctioneer      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IAuctioneerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type IAuctioneerCallerSession struct {
	Contract *IAuctioneerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// IAuctioneerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type IAuctioneerTransactorSession struct {
	Contract     *IAuctioneerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// IAuctioneerRaw is an auto generated low-level Go binding around an Ethereum contract.
type IAuctioneerRaw struct {
	Contract *IAuctioneer // Generic contract binding to access the raw methods on
}

// IAuctioneerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type IAuctioneerCallerRaw struct {
	Contract *IAuctioneerCaller // Generic read-only contract binding to access the raw methods on
}

// IAuctioneerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type IAuctioneerTransactorRaw struct {
	Contract *IAuctioneerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIAuctioneer creates a new instance of IAuctioneer, bound to a specific deployed contract.
func NewIAuctioneer(address common.Address, backend bind.ContractBackend) (*IAuctioneer, error) {
	contract, err := bindIAuctioneer(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IAuctioneer{IAuctioneerCaller: IAuctioneerCaller{contract: contract}, IAuctioneerTransactor: IAuctioneerTransactor{contract: contract}, IAuctioneerFilterer: IAuctioneerFilterer{contract: contract}}, nil
}

// NewIAuctioneerCaller creates a new read-only instance of IAuctioneer, bound to a specific deployed contract.
func NewIAuctioneerCaller(address common.Address, caller bind.ContractCaller) (*IAuctioneerCaller, error) {
	contract, err := bindIAuctioneer(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IAuctioneerCaller{contract: contract}, nil
}

// NewIAuctioneerTransactor creates a new write-only instance of IAuctioneer, bound to a specific deployed contract.
func NewIAuctioneerTransactor(address common.Address, transactor bind.ContractTransactor) (*IAuctioneerTransactor, error) {
	contract, err := bindIAuctioneer(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IAuctioneerTransactor{contract: contract}, nil
}

// NewIAuctioneerFilterer creates a new log filterer instance of IAuctioneer, bound to a specific deployed contract.
func NewIAuctioneerFilterer(address common.Address, filterer bind.ContractFilterer) (*IAuctioneerFilterer, error) {
	contract, err := bindIAuctioneer(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IAuctioneerFilterer{contract: contract}, nil
}

// bindIAuctioneer binds a generic wrapper to an already deployed contract.
func bindIAuctioneer(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(IAuctioneerABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IAuctioneer *IAuctioneerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IAuctioneer.Contract.IAuctioneerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IAuctioneer *IAuctioneerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IAuctioneer.Contract.IAuctioneerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IAuctioneer *IAuctioneerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IAuctioneer.Contract.IAuctioneerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IAuctioneer *IAuctioneerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IAuctioneer.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IAuctioneer *IAuctioneerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IAuctioneer.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IAuctioneer *IAuctioneerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IAuctioneer.Contract.contract.Transact(opts, method, params...)
}

// PaidInterest is a paid mutator transaction binding the contract method 0x96e4547c.
//
// Solidity: function paidInterest() payable returns()
func (_IAuctioneer *IAuctioneerTransactor) PaidInterest(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IAuctioneer.contract.Transact(opts, "paidInterest")
}

// PaidInterest is a paid mutator transaction binding the contract method 0x96e4547c.
//
// Solidity: function paidInterest() payable returns()
func (_IAuctioneer *IAuctioneerSession) PaidInterest() (*types.Transaction, error) {
	return _IAuctioneer.Contract.PaidInterest(&_IAuctioneer.TransactOpts)
}

// PaidInterest is a paid mutator transaction binding the contract method 0x96e4547c.
//
// Solidity: function paidInterest() payable returns()
func (_IAuctioneer *IAuctioneerTransactorSession) PaidInterest() (*types.Transaction, error) {
	return _IAuctioneer.Contract.PaidInterest(&_IAuctioneer.TransactOpts)
}

// SetOperator is a paid mutator transaction binding the contract method 0xb3ab15fb.
//
// Solidity: function setOperator(address operator) returns()
func (_IAuctioneer *IAuctioneerTransactor) SetOperator(opts *bind.TransactOpts, operator common.Address) (*types.Transaction, error) {
	return _IAuctioneer.contract.Transact(opts, "setOperator", operator)
}

// SetOperator is a paid mutator transaction binding the contract method 0xb3ab15fb.
//
// Solidity: function setOperator(address operator) returns()
func (_IAuctioneer *IAuctioneerSession) SetOperator(operator common.Address) (*types.Transaction, error) {
	return _IAuctioneer.Contract.SetOperator(&_IAuctioneer.TransactOpts, operator)
}

// SetOperator is a paid mutator transaction binding the contract method 0xb3ab15fb.
//
// Solidity: function setOperator(address operator) returns()
func (_IAuctioneer *IAuctioneerTransactorSession) SetOperator(operator common.Address) (*types.Transaction, error) {
	return _IAuctioneer.Contract.SetOperator(&_IAuctioneer.TransactOpts, operator)
}

// SetOracle is a paid mutator transaction binding the contract method 0x7adbf973.
//
// Solidity: function setOracle(address oracle) returns()
func (_IAuctioneer *IAuctioneerTransactor) SetOracle(opts *bind.TransactOpts, oracle common.Address) (*types.Transaction, error) {
	return _IAuctioneer.contract.Transact(opts, "setOracle", oracle)
}

// SetOracle is a paid mutator transaction binding the contract method 0x7adbf973.
//
// Solidity: function setOracle(address oracle) returns()
func (_IAuctioneer *IAuctioneerSession) SetOracle(oracle common.Address) (*types.Transaction, error) {
	return _IAuctioneer.Contract.SetOracle(&_IAuctioneer.TransactOpts, oracle)
}

// SetOracle is a paid mutator transaction binding the contract method 0x7adbf973.
//
// Solidity: function setOracle(address oracle) returns()
func (_IAuctioneer *IAuctioneerTransactorSession) SetOracle(oracle common.Address) (*types.Transaction, error) {
	return _IAuctioneer.Contract.SetOracle(&_IAuctioneer.TransactOpts, oracle)
}

// SetStabilization is a paid mutator transaction binding the contract method 0x4f505895.
//
// Solidity: function setStabilization(address stabilization) returns()
func (_IAuctioneer *IAuctioneerTransactor) SetStabilization(opts *bind.TransactOpts, stabilization common.Address) (*types.Transaction, error) {
	return _IAuctioneer.contract.Transact(opts, "setStabilization", stabilization)
}

// SetStabilization is a paid mutator transaction binding the contract method 0x4f505895.
//
// Solidity: function setStabilization(address stabilization) returns()
func (_IAuctioneer *IAuctioneerSession) SetStabilization(stabilization common.Address) (*types.Transaction, error) {
	return _IAuctioneer.Contract.SetStabilization(&_IAuctioneer.TransactOpts, stabilization)
}

// SetStabilization is a paid mutator transaction binding the contract method 0x4f505895.
//
// Solidity: function setStabilization(address stabilization) returns()
func (_IAuctioneer *IAuctioneerTransactorSession) SetStabilization(stabilization common.Address) (*types.Transaction, error) {
	return _IAuctioneer.Contract.SetStabilization(&_IAuctioneer.TransactOpts, stabilization)
}

// IAutonityMetaData contains all meta data concerning the IAutonity contract.
var IAutonityMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"treasury\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"effectiveBlock\",\"type\":\"uint256\"}],\"name\":\"ActivatedValidator\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_caller\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"BondingApproval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"enumIAutonity.ValidatorState\",\"name\":\"state\",\"type\":\"uint8\"}],\"name\":\"BondingRejected\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"BurnedStake\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"methodSignature\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"returnData\",\"type\":\"bytes\"}],\"name\":\"CallFailed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"rate\",\"type\":\"uint256\"}],\"name\":\"CommissionRateChange\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"minBaseFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"baseFeeChangeDenominator\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"elasticityMultiplier\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasLimitBoundDivisor\",\"type\":\"uint256\"}],\"indexed\":false,\"internalType\":\"structIAutonity.Eip1559\",\"name\":\"oldParams\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"minBaseFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"baseFeeChangeDenominator\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"elasticityMultiplier\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasLimitBoundDivisor\",\"type\":\"uint256\"}],\"indexed\":false,\"internalType\":\"structIAutonity.Eip1559\",\"name\":\"newParams\",\"type\":\"tuple\"}],\"name\":\"Eip1559ParamsUpdate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"oldEnode\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"newEnode\",\"type\":\"string\"}],\"name\":\"EnodeUpdate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"period\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"appliedAtBlock\",\"type\":\"uint256\"}],\"name\":\"EpochPeriodUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"MintedStake\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"caller\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"selfBonded\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"headBondingID\",\"type\":\"uint256\"}],\"name\":\"NewBondingRequest\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"epoch\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"inflationReserve\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"stakeCirculating\",\"type\":\"uint256\"}],\"name\":\"NewEpoch\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"caller\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"selfBonded\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"headUnbondingID\",\"type\":\"uint256\"}],\"name\":\"NewUnbondingRequest\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"treasury\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"effectiveBlock\",\"type\":\"uint256\"}],\"name\":\"PausedValidator\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"treasury\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"oracleAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"enode\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"liquidStateContract\",\"type\":\"address\"}],\"name\":\"RegisteredValidator\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"atnSelfAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"atnDelegatedAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"ntnSelfAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"ntnDelegatedAmount\",\"type\":\"uint256\"}],\"name\":\"Rewarded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_address\",\"type\":\"address\"}],\"name\":\"activateValidator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_caller\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"approveBonding\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_validator\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_selfBond\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_delegated\",\"type\":\"uint256\"}],\"name\":\"autobond\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_validator\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"bond\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_account\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_validator\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"bondFrom\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_caller\",\"type\":\"address\"}],\"name\":\"bondingAllowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_validator\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_rate\",\"type\":\"uint256\"}],\"name\":\"changeCommissionRate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"circulatingSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBlockPeriod\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_id\",\"type\":\"uint256\"}],\"name\":\"getBondingRequestByID\",\"outputs\":[{\"components\":[{\"internalType\":\"addresspayable\",\"name\":\"delegator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"delegatee\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"requestBlock\",\"type\":\"uint256\"}],\"internalType\":\"structIAutonity.BondingRequest\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getClientConfig\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"epochPeriod\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"blockPeriod\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"clusteringThreshold\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"range\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"delta\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gracePeriod\",\"type\":\"uint256\"}],\"internalType\":\"structIAutonity.Accountability\",\"name\":\"accountability\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"minBaseFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"baseFeeChangeDenominator\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"elasticityMultiplier\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasLimitBoundDivisor\",\"type\":\"uint256\"}],\"internalType\":\"structIAutonity.Eip1559\",\"name\":\"eip1559\",\"type\":\"tuple\"}],\"internalType\":\"structIAutonity.ClientAwareConfig\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCommittee\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"votingPower\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"consensusKey\",\"type\":\"bytes\"}],\"internalType\":\"structIAutonity.CommitteeMember[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCommitteeEnodes\",\"outputs\":[{\"internalType\":\"string[]\",\"name\":\"\",\"type\":\"string[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getConfig\",\"outputs\":[{\"components\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"treasuryFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"minBaseFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"delegationRate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"unbondingPeriod\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"initialInflationReserve\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"withholdingThreshold\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"proposerRewardRate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"oracleRewardRate\",\"type\":\"uint256\"},{\"internalType\":\"addresspayable\",\"name\":\"withheldRewardsPool\",\"type\":\"address\"},{\"internalType\":\"addresspayable\",\"name\":\"treasuryAccount\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"baseFeeChangeDenominator\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"elasticityMultiplier\",\"type\":\"uint256\"}],\"internalType\":\"structIAutonity.Policy\",\"name\":\"policy\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"contractIAccountability\",\"name\":\"accountabilityContract\",\"type\":\"address\"},{\"internalType\":\"contractIOracle\",\"name\":\"oracleContract\",\"type\":\"address\"},{\"internalType\":\"contractIACU\",\"name\":\"acuContract\",\"type\":\"address\"},{\"internalType\":\"contractISupplyControl\",\"name\":\"supplyControlContract\",\"type\":\"address\"},{\"internalType\":\"contractIStabilization\",\"name\":\"stabilizationContract\",\"type\":\"address\"},{\"internalType\":\"contractIUpgradeManager\",\"name\":\"upgradeManagerContract\",\"type\":\"address\"},{\"internalType\":\"contractIInflationController\",\"name\":\"inflationControllerContract\",\"type\":\"address\"},{\"internalType\":\"contractIOmissionAccountability\",\"name\":\"omissionAccountabilityContract\",\"type\":\"address\"},{\"internalType\":\"contractIAuctioneer\",\"name\":\"auctioneerContract\",\"type\":\"address\"}],\"internalType\":\"structIAutonity.Contracts\",\"name\":\"contracts\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"operatorAccount\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"epochPeriod\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"blockPeriod\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"committeeSize\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxScheduleDuration\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasLimitBoundDivisor\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"clusteringThreshold\",\"type\":\"uint256\"}],\"internalType\":\"structIAutonity.Protocol\",\"name\":\"protocol\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"contractVersion\",\"type\":\"uint256\"}],\"internalType\":\"structIAutonity.Config\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCurrentCommitteeSize\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCurrentEpochPeriod\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_height\",\"type\":\"uint256\"}],\"name\":\"getEpochByHeight\",\"outputs\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"votingPower\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"consensusKey\",\"type\":\"bytes\"}],\"internalType\":\"structIAutonity.CommitteeMember[]\",\"name\":\"committee\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"previousEpochBlock\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"epochBlock\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"nextEpochBlock\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"omissionDelta\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"minBaseFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"baseFeeChangeDenominator\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"elasticityMultiplier\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasLimitBoundDivisor\",\"type\":\"uint256\"}],\"internalType\":\"structIAutonity.Eip1559\",\"name\":\"eip1559\",\"type\":\"tuple\"}],\"internalType\":\"structIAutonity.EpochInfo\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_block\",\"type\":\"uint256\"}],\"name\":\"getEpochFromBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getEpochID\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getEpochInfo\",\"outputs\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"votingPower\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"consensusKey\",\"type\":\"bytes\"}],\"internalType\":\"structIAutonity.CommitteeMember[]\",\"name\":\"committee\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"previousEpochBlock\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"epochBlock\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"nextEpochBlock\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"omissionDelta\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"minBaseFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"baseFeeChangeDenominator\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"elasticityMultiplier\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasLimitBoundDivisor\",\"type\":\"uint256\"}],\"internalType\":\"structIAutonity.Eip1559\",\"name\":\"eip1559\",\"type\":\"tuple\"}],\"internalType\":\"structIAutonity.EpochInfo\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getEpochPeriod\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getEpochTotalBondedStake\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getInflationReserve\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLastEpochBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLastEpochTime\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLiquidLogicContract\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getMaxCommitteeSize\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getMaxScheduleDuration\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getMinimumBaseFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getNextEpochBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getOperator\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getOracle\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_vault\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_id\",\"type\":\"uint256\"}],\"name\":\"getSchedule\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"totalAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"unlockedAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"start\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalDuration\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"lastUnlockTime\",\"type\":\"uint256\"}],\"internalType\":\"structIScheduleController.Schedule\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_vault\",\"type\":\"address\"}],\"name\":\"getTotalSchedules\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTreasuryAccount\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTreasuryFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getUnbondingPeriod\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_id\",\"type\":\"uint256\"}],\"name\":\"getUnbondingRequestByID\",\"outputs\":[{\"components\":[{\"internalType\":\"addresspayable\",\"name\":\"delegator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"delegatee\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"unbondingShare\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"requestBlock\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"unlocked\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"released\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"selfDelegation\",\"type\":\"bool\"}],\"internalType\":\"structIAutonity.UnbondingRequest\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_unbondingID\",\"type\":\"uint256\"}],\"name\":\"getUnbondingShare\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"}],\"name\":\"getValidator\",\"outputs\":[{\"components\":[{\"internalType\":\"addresspayable\",\"name\":\"treasury\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"nodeAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"oracleAddress\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"enode\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"commissionRate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"bondedStake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"unbondingStake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"unbondingShares\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"selfBondedStake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"selfUnbondingStake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"selfUnbondingShares\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"selfUnbondingStakeLocked\",\"type\":\"uint256\"},{\"internalType\":\"contractILiquid\",\"name\":\"liquidStateContract\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"liquidSupply\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"registrationBlock\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalSlashed\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"jailReleaseBlock\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"consensusKey\",\"type\":\"bytes\"},{\"internalType\":\"enumIAutonity.ValidatorState\",\"name\":\"state\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"conversionRatio\",\"type\":\"uint256\"}],\"internalType\":\"structIAutonity.Validator\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"}],\"name\":\"getValidatorState\",\"outputs\":[{\"internalType\":\"enumIAutonity.ValidatorState\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getValidators\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getVersion\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_unbondingID\",\"type\":\"uint256\"}],\"name\":\"isUnbondingReleased\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_nodeAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_jailtime\",\"type\":\"uint256\"},{\"internalType\":\"enumIAutonity.ValidatorState\",\"name\":\"_newJailedState\",\"type\":\"uint8\"}],\"name\":\"jail\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_nodeAddress\",\"type\":\"address\"},{\"internalType\":\"enumIAutonity.ValidatorState\",\"name\":\"_newJailboundState\",\"type\":\"uint8\"}],\"name\":\"jailbound\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_address\",\"type\":\"address\"}],\"name\":\"pauseValidator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_enode\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"_oracleAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"_consensusKey\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"_signatures\",\"type\":\"bytes\"}],\"name\":\"registerValidator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_nodeAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_slashingRate\",\"type\":\"uint256\"}],\"name\":\"slash\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"slashingAmount\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_nodeAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_slashingRate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_jailtime\",\"type\":\"uint256\"},{\"internalType\":\"enumIAutonity.ValidatorState\",\"name\":\"_newJailedState\",\"type\":\"uint8\"},{\"internalType\":\"enumIAutonity.ValidatorState\",\"name\":\"_newJailboundState\",\"type\":\"uint8\"}],\"name\":\"slashAndJail\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"slashingAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"jailReleaseBlock\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"isJailbound\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_validator\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"unbond\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_account\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_validator\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"unbondFrom\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_nodeAddress\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"_enode\",\"type\":\"string\"}],\"name\":\"updateEnode\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"b46e5520": "activateValidator(address)",
		"dd62ed3e": "allowance(address,address)",
		"095ea7b3": "approve(address,uint256)",
		"50492571": "approveBonding(address,uint256)",
		"f7fcc510": "autobond(address,uint256,uint256)",
		"70a08231": "balanceOf(address)",
		"a515366a": "bond(address,uint256)",
		"41de7400": "bondFrom(address,address,uint256)",
		"e0e01d54": "bondingAllowance(address,address)",
		"852c4849": "changeCommissionRate(address,uint256)",
		"9358928b": "circulatingSupply()",
		"43645969": "getBlockPeriod()",
		"8ebb48b7": "getBondingRequestByID(uint256)",
		"f3f759c1": "getClientConfig()",
		"ab8f6ffe": "getCommittee()",
		"a8b2216e": "getCommitteeEnodes()",
		"c3f909d4": "getConfig()",
		"2b56feac": "getCurrentCommitteeSize()",
		"0aac2da1": "getCurrentEpochPeriod()",
		"affb1cf1": "getEpochByHeight(uint256)",
		"96b477cb": "getEpochFromBlock(uint256)",
		"6fc53515": "getEpochID()",
		"a9fd1a8f": "getEpochInfo()",
		"dfb1a4d2": "getEpochPeriod()",
		"4efcd15f": "getEpochTotalBondedStake()",
		"4651e07b": "getInflationReserve()",
		"731b3a03": "getLastEpochBlock()",
		"ba522458": "getLastEpochTime()",
		"4c1f1c77": "getLiquidLogicContract()",
		"819b6463": "getMaxCommitteeSize()",
		"fed76a56": "getMaxScheduleDuration()",
		"11220633": "getMinimumBaseFee()",
		"25ce1bb9": "getNextEpochBlock()",
		"e7f43c68": "getOperator()",
		"833b1fce": "getOracle()",
		"7264c4da": "getSchedule(address,uint256)",
		"088566e9": "getTotalSchedules(address)",
		"f7866ee3": "getTreasuryAccount()",
		"29070c6d": "getTreasuryFee()",
		"6fd2c80b": "getUnbondingPeriod()",
		"4bfe23f1": "getUnbondingRequestByID(uint256)",
		"8d347287": "getUnbondingShare(uint256)",
		"1904bb2e": "getValidator(address)",
		"5b7d6c36": "getValidatorState(address)",
		"b7ab4db5": "getValidators()",
		"0d8e6e2c": "getVersion()",
		"e294df7c": "isUnbondingReleased(uint256)",
		"154d76d7": "jail(address,uint256,uint8)",
		"8ef8c2fd": "jailbound(address,uint8)",
		"0ae65e7a": "pauseValidator(address)",
		"84467fdb": "registerValidator(string,address,bytes,bytes)",
		"02fb4d85": "slash(address,uint256)",
		"122b4122": "slashAndJail(address,uint256,uint256,uint8,uint8)",
		"18160ddd": "totalSupply()",
		"a9059cbb": "transfer(address,uint256)",
		"23b872dd": "transferFrom(address,address,uint256)",
		"a5d059ca": "unbond(address,uint256)",
		"a9a7d7c9": "unbondFrom(address,address,uint256)",
		"784304b5": "updateEnode(address,string)",
	},
}

// IAutonityABI is the input ABI used to generate the binding from.
// Deprecated: Use IAutonityMetaData.ABI instead.
var IAutonityABI = IAutonityMetaData.ABI

// Deprecated: Use IAutonityMetaData.Sigs instead.
// IAutonityFuncSigs maps the 4-byte function signature to its string representation.
var IAutonityFuncSigs = IAutonityMetaData.Sigs

// IAutonity is an auto generated Go binding around an Ethereum contract.
type IAutonity struct {
	IAutonityCaller     // Read-only binding to the contract
	IAutonityTransactor // Write-only binding to the contract
	IAutonityFilterer   // Log filterer for contract events
}

// IAutonityCaller is an auto generated read-only Go binding around an Ethereum contract.
type IAutonityCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IAutonityTransactor is an auto generated write-only Go binding around an Ethereum contract.
type IAutonityTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IAutonityFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type IAutonityFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IAutonitySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type IAutonitySession struct {
	Contract     *IAutonity        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IAutonityCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type IAutonityCallerSession struct {
	Contract *IAutonityCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// IAutonityTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type IAutonityTransactorSession struct {
	Contract     *IAutonityTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// IAutonityRaw is an auto generated low-level Go binding around an Ethereum contract.
type IAutonityRaw struct {
	Contract *IAutonity // Generic contract binding to access the raw methods on
}

// IAutonityCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type IAutonityCallerRaw struct {
	Contract *IAutonityCaller // Generic read-only contract binding to access the raw methods on
}

// IAutonityTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type IAutonityTransactorRaw struct {
	Contract *IAutonityTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIAutonity creates a new instance of IAutonity, bound to a specific deployed contract.
func NewIAutonity(address common.Address, backend bind.ContractBackend) (*IAutonity, error) {
	contract, err := bindIAutonity(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IAutonity{IAutonityCaller: IAutonityCaller{contract: contract}, IAutonityTransactor: IAutonityTransactor{contract: contract}, IAutonityFilterer: IAutonityFilterer{contract: contract}}, nil
}

// NewIAutonityCaller creates a new read-only instance of IAutonity, bound to a specific deployed contract.
func NewIAutonityCaller(address common.Address, caller bind.ContractCaller) (*IAutonityCaller, error) {
	contract, err := bindIAutonity(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IAutonityCaller{contract: contract}, nil
}

// NewIAutonityTransactor creates a new write-only instance of IAutonity, bound to a specific deployed contract.
func NewIAutonityTransactor(address common.Address, transactor bind.ContractTransactor) (*IAutonityTransactor, error) {
	contract, err := bindIAutonity(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IAutonityTransactor{contract: contract}, nil
}

// NewIAutonityFilterer creates a new log filterer instance of IAutonity, bound to a specific deployed contract.
func NewIAutonityFilterer(address common.Address, filterer bind.ContractFilterer) (*IAutonityFilterer, error) {
	contract, err := bindIAutonity(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IAutonityFilterer{contract: contract}, nil
}

// bindIAutonity binds a generic wrapper to an already deployed contract.
func bindIAutonity(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(IAutonityABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IAutonity *IAutonityRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IAutonity.Contract.IAutonityCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IAutonity *IAutonityRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IAutonity.Contract.IAutonityTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IAutonity *IAutonityRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IAutonity.Contract.IAutonityTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IAutonity *IAutonityCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IAutonity.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IAutonity *IAutonityTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IAutonity.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IAutonity *IAutonityTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IAutonity.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_IAutonity *IAutonityCaller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _IAutonity.contract.Call(opts, &out, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_IAutonity *IAutonitySession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _IAutonity.Contract.Allowance(&_IAutonity.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_IAutonity *IAutonityCallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _IAutonity.Contract.Allowance(&_IAutonity.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_IAutonity *IAutonityCaller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _IAutonity.contract.Call(opts, &out, "balanceOf", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_IAutonity *IAutonitySession) BalanceOf(account common.Address) (*big.Int, error) {
	return _IAutonity.Contract.BalanceOf(&_IAutonity.CallOpts, account)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_IAutonity *IAutonityCallerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _IAutonity.Contract.BalanceOf(&_IAutonity.CallOpts, account)
}

// BondingAllowance is a free data retrieval call binding the contract method 0xe0e01d54.
//
// Solidity: function bondingAllowance(address _owner, address _caller) view returns(uint256)
func (_IAutonity *IAutonityCaller) BondingAllowance(opts *bind.CallOpts, _owner common.Address, _caller common.Address) (*big.Int, error) {
	var out []interface{}
	err := _IAutonity.contract.Call(opts, &out, "bondingAllowance", _owner, _caller)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BondingAllowance is a free data retrieval call binding the contract method 0xe0e01d54.
//
// Solidity: function bondingAllowance(address _owner, address _caller) view returns(uint256)
func (_IAutonity *IAutonitySession) BondingAllowance(_owner common.Address, _caller common.Address) (*big.Int, error) {
	return _IAutonity.Contract.BondingAllowance(&_IAutonity.CallOpts, _owner, _caller)
}

// BondingAllowance is a free data retrieval call binding the contract method 0xe0e01d54.
//
// Solidity: function bondingAllowance(address _owner, address _caller) view returns(uint256)
func (_IAutonity *IAutonityCallerSession) BondingAllowance(_owner common.Address, _caller common.Address) (*big.Int, error) {
	return _IAutonity.Contract.BondingAllowance(&_IAutonity.CallOpts, _owner, _caller)
}

// CirculatingSupply is a free data retrieval call binding the contract method 0x9358928b.
//
// Solidity: function circulatingSupply() view returns(uint256)
func (_IAutonity *IAutonityCaller) CirculatingSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IAutonity.contract.Call(opts, &out, "circulatingSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CirculatingSupply is a free data retrieval call binding the contract method 0x9358928b.
//
// Solidity: function circulatingSupply() view returns(uint256)
func (_IAutonity *IAutonitySession) CirculatingSupply() (*big.Int, error) {
	return _IAutonity.Contract.CirculatingSupply(&_IAutonity.CallOpts)
}

// CirculatingSupply is a free data retrieval call binding the contract method 0x9358928b.
//
// Solidity: function circulatingSupply() view returns(uint256)
func (_IAutonity *IAutonityCallerSession) CirculatingSupply() (*big.Int, error) {
	return _IAutonity.Contract.CirculatingSupply(&_IAutonity.CallOpts)
}

// GetBlockPeriod is a free data retrieval call binding the contract method 0x43645969.
//
// Solidity: function getBlockPeriod() view returns(uint256)
func (_IAutonity *IAutonityCaller) GetBlockPeriod(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IAutonity.contract.Call(opts, &out, "getBlockPeriod")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetBlockPeriod is a free data retrieval call binding the contract method 0x43645969.
//
// Solidity: function getBlockPeriod() view returns(uint256)
func (_IAutonity *IAutonitySession) GetBlockPeriod() (*big.Int, error) {
	return _IAutonity.Contract.GetBlockPeriod(&_IAutonity.CallOpts)
}

// GetBlockPeriod is a free data retrieval call binding the contract method 0x43645969.
//
// Solidity: function getBlockPeriod() view returns(uint256)
func (_IAutonity *IAutonityCallerSession) GetBlockPeriod() (*big.Int, error) {
	return _IAutonity.Contract.GetBlockPeriod(&_IAutonity.CallOpts)
}

// GetBondingRequestByID is a free data retrieval call binding the contract method 0x8ebb48b7.
//
// Solidity: function getBondingRequestByID(uint256 _id) view returns((address,address,uint256,uint256))
func (_IAutonity *IAutonityCaller) GetBondingRequestByID(opts *bind.CallOpts, _id *big.Int) (IAutonityBondingRequest, error) {
	var out []interface{}
	err := _IAutonity.contract.Call(opts, &out, "getBondingRequestByID", _id)

	if err != nil {
		return *new(IAutonityBondingRequest), err
	}

	out0 := *abi.ConvertType(out[0], new(IAutonityBondingRequest)).(*IAutonityBondingRequest)

	return out0, err

}

// GetBondingRequestByID is a free data retrieval call binding the contract method 0x8ebb48b7.
//
// Solidity: function getBondingRequestByID(uint256 _id) view returns((address,address,uint256,uint256))
func (_IAutonity *IAutonitySession) GetBondingRequestByID(_id *big.Int) (IAutonityBondingRequest, error) {
	return _IAutonity.Contract.GetBondingRequestByID(&_IAutonity.CallOpts, _id)
}

// GetBondingRequestByID is a free data retrieval call binding the contract method 0x8ebb48b7.
//
// Solidity: function getBondingRequestByID(uint256 _id) view returns((address,address,uint256,uint256))
func (_IAutonity *IAutonityCallerSession) GetBondingRequestByID(_id *big.Int) (IAutonityBondingRequest, error) {
	return _IAutonity.Contract.GetBondingRequestByID(&_IAutonity.CallOpts, _id)
}

// GetClientConfig is a free data retrieval call binding the contract method 0xf3f759c1.
//
// Solidity: function getClientConfig() view returns((uint256,uint256,uint256,uint256,(uint256,uint256,uint256),(uint256,uint256,uint256,uint256)))
func (_IAutonity *IAutonityCaller) GetClientConfig(opts *bind.CallOpts) (IAutonityClientAwareConfig, error) {
	var out []interface{}
	err := _IAutonity.contract.Call(opts, &out, "getClientConfig")

	if err != nil {
		return *new(IAutonityClientAwareConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(IAutonityClientAwareConfig)).(*IAutonityClientAwareConfig)

	return out0, err

}

// GetClientConfig is a free data retrieval call binding the contract method 0xf3f759c1.
//
// Solidity: function getClientConfig() view returns((uint256,uint256,uint256,uint256,(uint256,uint256,uint256),(uint256,uint256,uint256,uint256)))
func (_IAutonity *IAutonitySession) GetClientConfig() (IAutonityClientAwareConfig, error) {
	return _IAutonity.Contract.GetClientConfig(&_IAutonity.CallOpts)
}

// GetClientConfig is a free data retrieval call binding the contract method 0xf3f759c1.
//
// Solidity: function getClientConfig() view returns((uint256,uint256,uint256,uint256,(uint256,uint256,uint256),(uint256,uint256,uint256,uint256)))
func (_IAutonity *IAutonityCallerSession) GetClientConfig() (IAutonityClientAwareConfig, error) {
	return _IAutonity.Contract.GetClientConfig(&_IAutonity.CallOpts)
}

// GetCommittee is a free data retrieval call binding the contract method 0xab8f6ffe.
//
// Solidity: function getCommittee() view returns((address,uint256,bytes)[])
func (_IAutonity *IAutonityCaller) GetCommittee(opts *bind.CallOpts) ([]IAutonityCommitteeMember, error) {
	var out []interface{}
	err := _IAutonity.contract.Call(opts, &out, "getCommittee")

	if err != nil {
		return *new([]IAutonityCommitteeMember), err
	}

	out0 := *abi.ConvertType(out[0], new([]IAutonityCommitteeMember)).(*[]IAutonityCommitteeMember)

	return out0, err

}

// GetCommittee is a free data retrieval call binding the contract method 0xab8f6ffe.
//
// Solidity: function getCommittee() view returns((address,uint256,bytes)[])
func (_IAutonity *IAutonitySession) GetCommittee() ([]IAutonityCommitteeMember, error) {
	return _IAutonity.Contract.GetCommittee(&_IAutonity.CallOpts)
}

// GetCommittee is a free data retrieval call binding the contract method 0xab8f6ffe.
//
// Solidity: function getCommittee() view returns((address,uint256,bytes)[])
func (_IAutonity *IAutonityCallerSession) GetCommittee() ([]IAutonityCommitteeMember, error) {
	return _IAutonity.Contract.GetCommittee(&_IAutonity.CallOpts)
}

// GetCommitteeEnodes is a free data retrieval call binding the contract method 0xa8b2216e.
//
// Solidity: function getCommitteeEnodes() view returns(string[])
func (_IAutonity *IAutonityCaller) GetCommitteeEnodes(opts *bind.CallOpts) ([]string, error) {
	var out []interface{}
	err := _IAutonity.contract.Call(opts, &out, "getCommitteeEnodes")

	if err != nil {
		return *new([]string), err
	}

	out0 := *abi.ConvertType(out[0], new([]string)).(*[]string)

	return out0, err

}

// GetCommitteeEnodes is a free data retrieval call binding the contract method 0xa8b2216e.
//
// Solidity: function getCommitteeEnodes() view returns(string[])
func (_IAutonity *IAutonitySession) GetCommitteeEnodes() ([]string, error) {
	return _IAutonity.Contract.GetCommitteeEnodes(&_IAutonity.CallOpts)
}

// GetCommitteeEnodes is a free data retrieval call binding the contract method 0xa8b2216e.
//
// Solidity: function getCommitteeEnodes() view returns(string[])
func (_IAutonity *IAutonityCallerSession) GetCommitteeEnodes() ([]string, error) {
	return _IAutonity.Contract.GetCommitteeEnodes(&_IAutonity.CallOpts)
}

// GetConfig is a free data retrieval call binding the contract method 0xc3f909d4.
//
// Solidity: function getConfig() view returns(((uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256,address,address,uint256,uint256),(address,address,address,address,address,address,address,address,address),(address,uint256,uint256,uint256,uint256,uint256,uint256,uint256),uint256))
func (_IAutonity *IAutonityCaller) GetConfig(opts *bind.CallOpts) (IAutonityConfig, error) {
	var out []interface{}
	err := _IAutonity.contract.Call(opts, &out, "getConfig")

	if err != nil {
		return *new(IAutonityConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(IAutonityConfig)).(*IAutonityConfig)

	return out0, err

}

// GetConfig is a free data retrieval call binding the contract method 0xc3f909d4.
//
// Solidity: function getConfig() view returns(((uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256,address,address,uint256,uint256),(address,address,address,address,address,address,address,address,address),(address,uint256,uint256,uint256,uint256,uint256,uint256,uint256),uint256))
func (_IAutonity *IAutonitySession) GetConfig() (IAutonityConfig, error) {
	return _IAutonity.Contract.GetConfig(&_IAutonity.CallOpts)
}

// GetConfig is a free data retrieval call binding the contract method 0xc3f909d4.
//
// Solidity: function getConfig() view returns(((uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256,address,address,uint256,uint256),(address,address,address,address,address,address,address,address,address),(address,uint256,uint256,uint256,uint256,uint256,uint256,uint256),uint256))
func (_IAutonity *IAutonityCallerSession) GetConfig() (IAutonityConfig, error) {
	return _IAutonity.Contract.GetConfig(&_IAutonity.CallOpts)
}

// GetCurrentCommitteeSize is a free data retrieval call binding the contract method 0x2b56feac.
//
// Solidity: function getCurrentCommitteeSize() view returns(uint256)
func (_IAutonity *IAutonityCaller) GetCurrentCommitteeSize(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IAutonity.contract.Call(opts, &out, "getCurrentCommitteeSize")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetCurrentCommitteeSize is a free data retrieval call binding the contract method 0x2b56feac.
//
// Solidity: function getCurrentCommitteeSize() view returns(uint256)
func (_IAutonity *IAutonitySession) GetCurrentCommitteeSize() (*big.Int, error) {
	return _IAutonity.Contract.GetCurrentCommitteeSize(&_IAutonity.CallOpts)
}

// GetCurrentCommitteeSize is a free data retrieval call binding the contract method 0x2b56feac.
//
// Solidity: function getCurrentCommitteeSize() view returns(uint256)
func (_IAutonity *IAutonityCallerSession) GetCurrentCommitteeSize() (*big.Int, error) {
	return _IAutonity.Contract.GetCurrentCommitteeSize(&_IAutonity.CallOpts)
}

// GetCurrentEpochPeriod is a free data retrieval call binding the contract method 0x0aac2da1.
//
// Solidity: function getCurrentEpochPeriod() view returns(uint256)
func (_IAutonity *IAutonityCaller) GetCurrentEpochPeriod(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IAutonity.contract.Call(opts, &out, "getCurrentEpochPeriod")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetCurrentEpochPeriod is a free data retrieval call binding the contract method 0x0aac2da1.
//
// Solidity: function getCurrentEpochPeriod() view returns(uint256)
func (_IAutonity *IAutonitySession) GetCurrentEpochPeriod() (*big.Int, error) {
	return _IAutonity.Contract.GetCurrentEpochPeriod(&_IAutonity.CallOpts)
}

// GetCurrentEpochPeriod is a free data retrieval call binding the contract method 0x0aac2da1.
//
// Solidity: function getCurrentEpochPeriod() view returns(uint256)
func (_IAutonity *IAutonityCallerSession) GetCurrentEpochPeriod() (*big.Int, error) {
	return _IAutonity.Contract.GetCurrentEpochPeriod(&_IAutonity.CallOpts)
}

// GetEpochByHeight is a free data retrieval call binding the contract method 0xaffb1cf1.
//
// Solidity: function getEpochByHeight(uint256 _height) view returns(((address,uint256,bytes)[],uint256,uint256,uint256,uint256,(uint256,uint256,uint256,uint256)))
func (_IAutonity *IAutonityCaller) GetEpochByHeight(opts *bind.CallOpts, _height *big.Int) (IAutonityEpochInfo, error) {
	var out []interface{}
	err := _IAutonity.contract.Call(opts, &out, "getEpochByHeight", _height)

	if err != nil {
		return *new(IAutonityEpochInfo), err
	}

	out0 := *abi.ConvertType(out[0], new(IAutonityEpochInfo)).(*IAutonityEpochInfo)

	return out0, err

}

// GetEpochByHeight is a free data retrieval call binding the contract method 0xaffb1cf1.
//
// Solidity: function getEpochByHeight(uint256 _height) view returns(((address,uint256,bytes)[],uint256,uint256,uint256,uint256,(uint256,uint256,uint256,uint256)))
func (_IAutonity *IAutonitySession) GetEpochByHeight(_height *big.Int) (IAutonityEpochInfo, error) {
	return _IAutonity.Contract.GetEpochByHeight(&_IAutonity.CallOpts, _height)
}

// GetEpochByHeight is a free data retrieval call binding the contract method 0xaffb1cf1.
//
// Solidity: function getEpochByHeight(uint256 _height) view returns(((address,uint256,bytes)[],uint256,uint256,uint256,uint256,(uint256,uint256,uint256,uint256)))
func (_IAutonity *IAutonityCallerSession) GetEpochByHeight(_height *big.Int) (IAutonityEpochInfo, error) {
	return _IAutonity.Contract.GetEpochByHeight(&_IAutonity.CallOpts, _height)
}

// GetEpochFromBlock is a free data retrieval call binding the contract method 0x96b477cb.
//
// Solidity: function getEpochFromBlock(uint256 _block) view returns(uint256)
func (_IAutonity *IAutonityCaller) GetEpochFromBlock(opts *bind.CallOpts, _block *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _IAutonity.contract.Call(opts, &out, "getEpochFromBlock", _block)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetEpochFromBlock is a free data retrieval call binding the contract method 0x96b477cb.
//
// Solidity: function getEpochFromBlock(uint256 _block) view returns(uint256)
func (_IAutonity *IAutonitySession) GetEpochFromBlock(_block *big.Int) (*big.Int, error) {
	return _IAutonity.Contract.GetEpochFromBlock(&_IAutonity.CallOpts, _block)
}

// GetEpochFromBlock is a free data retrieval call binding the contract method 0x96b477cb.
//
// Solidity: function getEpochFromBlock(uint256 _block) view returns(uint256)
func (_IAutonity *IAutonityCallerSession) GetEpochFromBlock(_block *big.Int) (*big.Int, error) {
	return _IAutonity.Contract.GetEpochFromBlock(&_IAutonity.CallOpts, _block)
}

// GetEpochID is a free data retrieval call binding the contract method 0x6fc53515.
//
// Solidity: function getEpochID() view returns(uint256)
func (_IAutonity *IAutonityCaller) GetEpochID(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IAutonity.contract.Call(opts, &out, "getEpochID")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetEpochID is a free data retrieval call binding the contract method 0x6fc53515.
//
// Solidity: function getEpochID() view returns(uint256)
func (_IAutonity *IAutonitySession) GetEpochID() (*big.Int, error) {
	return _IAutonity.Contract.GetEpochID(&_IAutonity.CallOpts)
}

// GetEpochID is a free data retrieval call binding the contract method 0x6fc53515.
//
// Solidity: function getEpochID() view returns(uint256)
func (_IAutonity *IAutonityCallerSession) GetEpochID() (*big.Int, error) {
	return _IAutonity.Contract.GetEpochID(&_IAutonity.CallOpts)
}

// GetEpochInfo is a free data retrieval call binding the contract method 0xa9fd1a8f.
//
// Solidity: function getEpochInfo() view returns(((address,uint256,bytes)[],uint256,uint256,uint256,uint256,(uint256,uint256,uint256,uint256)))
func (_IAutonity *IAutonityCaller) GetEpochInfo(opts *bind.CallOpts) (IAutonityEpochInfo, error) {
	var out []interface{}
	err := _IAutonity.contract.Call(opts, &out, "getEpochInfo")

	if err != nil {
		return *new(IAutonityEpochInfo), err
	}

	out0 := *abi.ConvertType(out[0], new(IAutonityEpochInfo)).(*IAutonityEpochInfo)

	return out0, err

}

// GetEpochInfo is a free data retrieval call binding the contract method 0xa9fd1a8f.
//
// Solidity: function getEpochInfo() view returns(((address,uint256,bytes)[],uint256,uint256,uint256,uint256,(uint256,uint256,uint256,uint256)))
func (_IAutonity *IAutonitySession) GetEpochInfo() (IAutonityEpochInfo, error) {
	return _IAutonity.Contract.GetEpochInfo(&_IAutonity.CallOpts)
}

// GetEpochInfo is a free data retrieval call binding the contract method 0xa9fd1a8f.
//
// Solidity: function getEpochInfo() view returns(((address,uint256,bytes)[],uint256,uint256,uint256,uint256,(uint256,uint256,uint256,uint256)))
func (_IAutonity *IAutonityCallerSession) GetEpochInfo() (IAutonityEpochInfo, error) {
	return _IAutonity.Contract.GetEpochInfo(&_IAutonity.CallOpts)
}

// GetEpochPeriod is a free data retrieval call binding the contract method 0xdfb1a4d2.
//
// Solidity: function getEpochPeriod() view returns(uint256)
func (_IAutonity *IAutonityCaller) GetEpochPeriod(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IAutonity.contract.Call(opts, &out, "getEpochPeriod")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetEpochPeriod is a free data retrieval call binding the contract method 0xdfb1a4d2.
//
// Solidity: function getEpochPeriod() view returns(uint256)
func (_IAutonity *IAutonitySession) GetEpochPeriod() (*big.Int, error) {
	return _IAutonity.Contract.GetEpochPeriod(&_IAutonity.CallOpts)
}

// GetEpochPeriod is a free data retrieval call binding the contract method 0xdfb1a4d2.
//
// Solidity: function getEpochPeriod() view returns(uint256)
func (_IAutonity *IAutonityCallerSession) GetEpochPeriod() (*big.Int, error) {
	return _IAutonity.Contract.GetEpochPeriod(&_IAutonity.CallOpts)
}

// GetEpochTotalBondedStake is a free data retrieval call binding the contract method 0x4efcd15f.
//
// Solidity: function getEpochTotalBondedStake() view returns(uint256)
func (_IAutonity *IAutonityCaller) GetEpochTotalBondedStake(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IAutonity.contract.Call(opts, &out, "getEpochTotalBondedStake")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetEpochTotalBondedStake is a free data retrieval call binding the contract method 0x4efcd15f.
//
// Solidity: function getEpochTotalBondedStake() view returns(uint256)
func (_IAutonity *IAutonitySession) GetEpochTotalBondedStake() (*big.Int, error) {
	return _IAutonity.Contract.GetEpochTotalBondedStake(&_IAutonity.CallOpts)
}

// GetEpochTotalBondedStake is a free data retrieval call binding the contract method 0x4efcd15f.
//
// Solidity: function getEpochTotalBondedStake() view returns(uint256)
func (_IAutonity *IAutonityCallerSession) GetEpochTotalBondedStake() (*big.Int, error) {
	return _IAutonity.Contract.GetEpochTotalBondedStake(&_IAutonity.CallOpts)
}

// GetInflationReserve is a free data retrieval call binding the contract method 0x4651e07b.
//
// Solidity: function getInflationReserve() view returns(uint256)
func (_IAutonity *IAutonityCaller) GetInflationReserve(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IAutonity.contract.Call(opts, &out, "getInflationReserve")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetInflationReserve is a free data retrieval call binding the contract method 0x4651e07b.
//
// Solidity: function getInflationReserve() view returns(uint256)
func (_IAutonity *IAutonitySession) GetInflationReserve() (*big.Int, error) {
	return _IAutonity.Contract.GetInflationReserve(&_IAutonity.CallOpts)
}

// GetInflationReserve is a free data retrieval call binding the contract method 0x4651e07b.
//
// Solidity: function getInflationReserve() view returns(uint256)
func (_IAutonity *IAutonityCallerSession) GetInflationReserve() (*big.Int, error) {
	return _IAutonity.Contract.GetInflationReserve(&_IAutonity.CallOpts)
}

// GetLastEpochBlock is a free data retrieval call binding the contract method 0x731b3a03.
//
// Solidity: function getLastEpochBlock() view returns(uint256)
func (_IAutonity *IAutonityCaller) GetLastEpochBlock(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IAutonity.contract.Call(opts, &out, "getLastEpochBlock")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetLastEpochBlock is a free data retrieval call binding the contract method 0x731b3a03.
//
// Solidity: function getLastEpochBlock() view returns(uint256)
func (_IAutonity *IAutonitySession) GetLastEpochBlock() (*big.Int, error) {
	return _IAutonity.Contract.GetLastEpochBlock(&_IAutonity.CallOpts)
}

// GetLastEpochBlock is a free data retrieval call binding the contract method 0x731b3a03.
//
// Solidity: function getLastEpochBlock() view returns(uint256)
func (_IAutonity *IAutonityCallerSession) GetLastEpochBlock() (*big.Int, error) {
	return _IAutonity.Contract.GetLastEpochBlock(&_IAutonity.CallOpts)
}

// GetLastEpochTime is a free data retrieval call binding the contract method 0xba522458.
//
// Solidity: function getLastEpochTime() view returns(uint256)
func (_IAutonity *IAutonityCaller) GetLastEpochTime(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IAutonity.contract.Call(opts, &out, "getLastEpochTime")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetLastEpochTime is a free data retrieval call binding the contract method 0xba522458.
//
// Solidity: function getLastEpochTime() view returns(uint256)
func (_IAutonity *IAutonitySession) GetLastEpochTime() (*big.Int, error) {
	return _IAutonity.Contract.GetLastEpochTime(&_IAutonity.CallOpts)
}

// GetLastEpochTime is a free data retrieval call binding the contract method 0xba522458.
//
// Solidity: function getLastEpochTime() view returns(uint256)
func (_IAutonity *IAutonityCallerSession) GetLastEpochTime() (*big.Int, error) {
	return _IAutonity.Contract.GetLastEpochTime(&_IAutonity.CallOpts)
}

// GetLiquidLogicContract is a free data retrieval call binding the contract method 0x4c1f1c77.
//
// Solidity: function getLiquidLogicContract() view returns(address)
func (_IAutonity *IAutonityCaller) GetLiquidLogicContract(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IAutonity.contract.Call(opts, &out, "getLiquidLogicContract")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetLiquidLogicContract is a free data retrieval call binding the contract method 0x4c1f1c77.
//
// Solidity: function getLiquidLogicContract() view returns(address)
func (_IAutonity *IAutonitySession) GetLiquidLogicContract() (common.Address, error) {
	return _IAutonity.Contract.GetLiquidLogicContract(&_IAutonity.CallOpts)
}

// GetLiquidLogicContract is a free data retrieval call binding the contract method 0x4c1f1c77.
//
// Solidity: function getLiquidLogicContract() view returns(address)
func (_IAutonity *IAutonityCallerSession) GetLiquidLogicContract() (common.Address, error) {
	return _IAutonity.Contract.GetLiquidLogicContract(&_IAutonity.CallOpts)
}

// GetMaxCommitteeSize is a free data retrieval call binding the contract method 0x819b6463.
//
// Solidity: function getMaxCommitteeSize() view returns(uint256)
func (_IAutonity *IAutonityCaller) GetMaxCommitteeSize(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IAutonity.contract.Call(opts, &out, "getMaxCommitteeSize")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetMaxCommitteeSize is a free data retrieval call binding the contract method 0x819b6463.
//
// Solidity: function getMaxCommitteeSize() view returns(uint256)
func (_IAutonity *IAutonitySession) GetMaxCommitteeSize() (*big.Int, error) {
	return _IAutonity.Contract.GetMaxCommitteeSize(&_IAutonity.CallOpts)
}

// GetMaxCommitteeSize is a free data retrieval call binding the contract method 0x819b6463.
//
// Solidity: function getMaxCommitteeSize() view returns(uint256)
func (_IAutonity *IAutonityCallerSession) GetMaxCommitteeSize() (*big.Int, error) {
	return _IAutonity.Contract.GetMaxCommitteeSize(&_IAutonity.CallOpts)
}

// GetMaxScheduleDuration is a free data retrieval call binding the contract method 0xfed76a56.
//
// Solidity: function getMaxScheduleDuration() view returns(uint256)
func (_IAutonity *IAutonityCaller) GetMaxScheduleDuration(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IAutonity.contract.Call(opts, &out, "getMaxScheduleDuration")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetMaxScheduleDuration is a free data retrieval call binding the contract method 0xfed76a56.
//
// Solidity: function getMaxScheduleDuration() view returns(uint256)
func (_IAutonity *IAutonitySession) GetMaxScheduleDuration() (*big.Int, error) {
	return _IAutonity.Contract.GetMaxScheduleDuration(&_IAutonity.CallOpts)
}

// GetMaxScheduleDuration is a free data retrieval call binding the contract method 0xfed76a56.
//
// Solidity: function getMaxScheduleDuration() view returns(uint256)
func (_IAutonity *IAutonityCallerSession) GetMaxScheduleDuration() (*big.Int, error) {
	return _IAutonity.Contract.GetMaxScheduleDuration(&_IAutonity.CallOpts)
}

// GetMinimumBaseFee is a free data retrieval call binding the contract method 0x11220633.
//
// Solidity: function getMinimumBaseFee() view returns(uint256)
func (_IAutonity *IAutonityCaller) GetMinimumBaseFee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IAutonity.contract.Call(opts, &out, "getMinimumBaseFee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetMinimumBaseFee is a free data retrieval call binding the contract method 0x11220633.
//
// Solidity: function getMinimumBaseFee() view returns(uint256)
func (_IAutonity *IAutonitySession) GetMinimumBaseFee() (*big.Int, error) {
	return _IAutonity.Contract.GetMinimumBaseFee(&_IAutonity.CallOpts)
}

// GetMinimumBaseFee is a free data retrieval call binding the contract method 0x11220633.
//
// Solidity: function getMinimumBaseFee() view returns(uint256)
func (_IAutonity *IAutonityCallerSession) GetMinimumBaseFee() (*big.Int, error) {
	return _IAutonity.Contract.GetMinimumBaseFee(&_IAutonity.CallOpts)
}

// GetNextEpochBlock is a free data retrieval call binding the contract method 0x25ce1bb9.
//
// Solidity: function getNextEpochBlock() view returns(uint256)
func (_IAutonity *IAutonityCaller) GetNextEpochBlock(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IAutonity.contract.Call(opts, &out, "getNextEpochBlock")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetNextEpochBlock is a free data retrieval call binding the contract method 0x25ce1bb9.
//
// Solidity: function getNextEpochBlock() view returns(uint256)
func (_IAutonity *IAutonitySession) GetNextEpochBlock() (*big.Int, error) {
	return _IAutonity.Contract.GetNextEpochBlock(&_IAutonity.CallOpts)
}

// GetNextEpochBlock is a free data retrieval call binding the contract method 0x25ce1bb9.
//
// Solidity: function getNextEpochBlock() view returns(uint256)
func (_IAutonity *IAutonityCallerSession) GetNextEpochBlock() (*big.Int, error) {
	return _IAutonity.Contract.GetNextEpochBlock(&_IAutonity.CallOpts)
}

// GetOperator is a free data retrieval call binding the contract method 0xe7f43c68.
//
// Solidity: function getOperator() view returns(address)
func (_IAutonity *IAutonityCaller) GetOperator(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IAutonity.contract.Call(opts, &out, "getOperator")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetOperator is a free data retrieval call binding the contract method 0xe7f43c68.
//
// Solidity: function getOperator() view returns(address)
func (_IAutonity *IAutonitySession) GetOperator() (common.Address, error) {
	return _IAutonity.Contract.GetOperator(&_IAutonity.CallOpts)
}

// GetOperator is a free data retrieval call binding the contract method 0xe7f43c68.
//
// Solidity: function getOperator() view returns(address)
func (_IAutonity *IAutonityCallerSession) GetOperator() (common.Address, error) {
	return _IAutonity.Contract.GetOperator(&_IAutonity.CallOpts)
}

// GetOracle is a free data retrieval call binding the contract method 0x833b1fce.
//
// Solidity: function getOracle() view returns(address)
func (_IAutonity *IAutonityCaller) GetOracle(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IAutonity.contract.Call(opts, &out, "getOracle")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetOracle is a free data retrieval call binding the contract method 0x833b1fce.
//
// Solidity: function getOracle() view returns(address)
func (_IAutonity *IAutonitySession) GetOracle() (common.Address, error) {
	return _IAutonity.Contract.GetOracle(&_IAutonity.CallOpts)
}

// GetOracle is a free data retrieval call binding the contract method 0x833b1fce.
//
// Solidity: function getOracle() view returns(address)
func (_IAutonity *IAutonityCallerSession) GetOracle() (common.Address, error) {
	return _IAutonity.Contract.GetOracle(&_IAutonity.CallOpts)
}

// GetSchedule is a free data retrieval call binding the contract method 0x7264c4da.
//
// Solidity: function getSchedule(address _vault, uint256 _id) view returns((uint256,uint256,uint256,uint256,uint256))
func (_IAutonity *IAutonityCaller) GetSchedule(opts *bind.CallOpts, _vault common.Address, _id *big.Int) (IScheduleControllerSchedule, error) {
	var out []interface{}
	err := _IAutonity.contract.Call(opts, &out, "getSchedule", _vault, _id)

	if err != nil {
		return *new(IScheduleControllerSchedule), err
	}

	out0 := *abi.ConvertType(out[0], new(IScheduleControllerSchedule)).(*IScheduleControllerSchedule)

	return out0, err

}

// GetSchedule is a free data retrieval call binding the contract method 0x7264c4da.
//
// Solidity: function getSchedule(address _vault, uint256 _id) view returns((uint256,uint256,uint256,uint256,uint256))
func (_IAutonity *IAutonitySession) GetSchedule(_vault common.Address, _id *big.Int) (IScheduleControllerSchedule, error) {
	return _IAutonity.Contract.GetSchedule(&_IAutonity.CallOpts, _vault, _id)
}

// GetSchedule is a free data retrieval call binding the contract method 0x7264c4da.
//
// Solidity: function getSchedule(address _vault, uint256 _id) view returns((uint256,uint256,uint256,uint256,uint256))
func (_IAutonity *IAutonityCallerSession) GetSchedule(_vault common.Address, _id *big.Int) (IScheduleControllerSchedule, error) {
	return _IAutonity.Contract.GetSchedule(&_IAutonity.CallOpts, _vault, _id)
}

// GetTotalSchedules is a free data retrieval call binding the contract method 0x088566e9.
//
// Solidity: function getTotalSchedules(address _vault) view returns(uint256)
func (_IAutonity *IAutonityCaller) GetTotalSchedules(opts *bind.CallOpts, _vault common.Address) (*big.Int, error) {
	var out []interface{}
	err := _IAutonity.contract.Call(opts, &out, "getTotalSchedules", _vault)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetTotalSchedules is a free data retrieval call binding the contract method 0x088566e9.
//
// Solidity: function getTotalSchedules(address _vault) view returns(uint256)
func (_IAutonity *IAutonitySession) GetTotalSchedules(_vault common.Address) (*big.Int, error) {
	return _IAutonity.Contract.GetTotalSchedules(&_IAutonity.CallOpts, _vault)
}

// GetTotalSchedules is a free data retrieval call binding the contract method 0x088566e9.
//
// Solidity: function getTotalSchedules(address _vault) view returns(uint256)
func (_IAutonity *IAutonityCallerSession) GetTotalSchedules(_vault common.Address) (*big.Int, error) {
	return _IAutonity.Contract.GetTotalSchedules(&_IAutonity.CallOpts, _vault)
}

// GetTreasuryAccount is a free data retrieval call binding the contract method 0xf7866ee3.
//
// Solidity: function getTreasuryAccount() view returns(address)
func (_IAutonity *IAutonityCaller) GetTreasuryAccount(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IAutonity.contract.Call(opts, &out, "getTreasuryAccount")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetTreasuryAccount is a free data retrieval call binding the contract method 0xf7866ee3.
//
// Solidity: function getTreasuryAccount() view returns(address)
func (_IAutonity *IAutonitySession) GetTreasuryAccount() (common.Address, error) {
	return _IAutonity.Contract.GetTreasuryAccount(&_IAutonity.CallOpts)
}

// GetTreasuryAccount is a free data retrieval call binding the contract method 0xf7866ee3.
//
// Solidity: function getTreasuryAccount() view returns(address)
func (_IAutonity *IAutonityCallerSession) GetTreasuryAccount() (common.Address, error) {
	return _IAutonity.Contract.GetTreasuryAccount(&_IAutonity.CallOpts)
}

// GetTreasuryFee is a free data retrieval call binding the contract method 0x29070c6d.
//
// Solidity: function getTreasuryFee() view returns(uint256)
func (_IAutonity *IAutonityCaller) GetTreasuryFee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IAutonity.contract.Call(opts, &out, "getTreasuryFee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetTreasuryFee is a free data retrieval call binding the contract method 0x29070c6d.
//
// Solidity: function getTreasuryFee() view returns(uint256)
func (_IAutonity *IAutonitySession) GetTreasuryFee() (*big.Int, error) {
	return _IAutonity.Contract.GetTreasuryFee(&_IAutonity.CallOpts)
}

// GetTreasuryFee is a free data retrieval call binding the contract method 0x29070c6d.
//
// Solidity: function getTreasuryFee() view returns(uint256)
func (_IAutonity *IAutonityCallerSession) GetTreasuryFee() (*big.Int, error) {
	return _IAutonity.Contract.GetTreasuryFee(&_IAutonity.CallOpts)
}

// GetUnbondingPeriod is a free data retrieval call binding the contract method 0x6fd2c80b.
//
// Solidity: function getUnbondingPeriod() view returns(uint256)
func (_IAutonity *IAutonityCaller) GetUnbondingPeriod(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IAutonity.contract.Call(opts, &out, "getUnbondingPeriod")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetUnbondingPeriod is a free data retrieval call binding the contract method 0x6fd2c80b.
//
// Solidity: function getUnbondingPeriod() view returns(uint256)
func (_IAutonity *IAutonitySession) GetUnbondingPeriod() (*big.Int, error) {
	return _IAutonity.Contract.GetUnbondingPeriod(&_IAutonity.CallOpts)
}

// GetUnbondingPeriod is a free data retrieval call binding the contract method 0x6fd2c80b.
//
// Solidity: function getUnbondingPeriod() view returns(uint256)
func (_IAutonity *IAutonityCallerSession) GetUnbondingPeriod() (*big.Int, error) {
	return _IAutonity.Contract.GetUnbondingPeriod(&_IAutonity.CallOpts)
}

// GetUnbondingRequestByID is a free data retrieval call binding the contract method 0x4bfe23f1.
//
// Solidity: function getUnbondingRequestByID(uint256 _id) view returns((address,address,uint256,uint256,uint256,bool,bool,bool))
func (_IAutonity *IAutonityCaller) GetUnbondingRequestByID(opts *bind.CallOpts, _id *big.Int) (IAutonityUnbondingRequest, error) {
	var out []interface{}
	err := _IAutonity.contract.Call(opts, &out, "getUnbondingRequestByID", _id)

	if err != nil {
		return *new(IAutonityUnbondingRequest), err
	}

	out0 := *abi.ConvertType(out[0], new(IAutonityUnbondingRequest)).(*IAutonityUnbondingRequest)

	return out0, err

}

// GetUnbondingRequestByID is a free data retrieval call binding the contract method 0x4bfe23f1.
//
// Solidity: function getUnbondingRequestByID(uint256 _id) view returns((address,address,uint256,uint256,uint256,bool,bool,bool))
func (_IAutonity *IAutonitySession) GetUnbondingRequestByID(_id *big.Int) (IAutonityUnbondingRequest, error) {
	return _IAutonity.Contract.GetUnbondingRequestByID(&_IAutonity.CallOpts, _id)
}

// GetUnbondingRequestByID is a free data retrieval call binding the contract method 0x4bfe23f1.
//
// Solidity: function getUnbondingRequestByID(uint256 _id) view returns((address,address,uint256,uint256,uint256,bool,bool,bool))
func (_IAutonity *IAutonityCallerSession) GetUnbondingRequestByID(_id *big.Int) (IAutonityUnbondingRequest, error) {
	return _IAutonity.Contract.GetUnbondingRequestByID(&_IAutonity.CallOpts, _id)
}

// GetUnbondingShare is a free data retrieval call binding the contract method 0x8d347287.
//
// Solidity: function getUnbondingShare(uint256 _unbondingID) view returns(uint256)
func (_IAutonity *IAutonityCaller) GetUnbondingShare(opts *bind.CallOpts, _unbondingID *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _IAutonity.contract.Call(opts, &out, "getUnbondingShare", _unbondingID)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetUnbondingShare is a free data retrieval call binding the contract method 0x8d347287.
//
// Solidity: function getUnbondingShare(uint256 _unbondingID) view returns(uint256)
func (_IAutonity *IAutonitySession) GetUnbondingShare(_unbondingID *big.Int) (*big.Int, error) {
	return _IAutonity.Contract.GetUnbondingShare(&_IAutonity.CallOpts, _unbondingID)
}

// GetUnbondingShare is a free data retrieval call binding the contract method 0x8d347287.
//
// Solidity: function getUnbondingShare(uint256 _unbondingID) view returns(uint256)
func (_IAutonity *IAutonityCallerSession) GetUnbondingShare(_unbondingID *big.Int) (*big.Int, error) {
	return _IAutonity.Contract.GetUnbondingShare(&_IAutonity.CallOpts, _unbondingID)
}

// GetValidator is a free data retrieval call binding the contract method 0x1904bb2e.
//
// Solidity: function getValidator(address _addr) view returns((address,address,address,string,uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256,address,uint256,uint256,uint256,uint256,bytes,uint8,uint256))
func (_IAutonity *IAutonityCaller) GetValidator(opts *bind.CallOpts, _addr common.Address) (IAutonityValidator, error) {
	var out []interface{}
	err := _IAutonity.contract.Call(opts, &out, "getValidator", _addr)

	if err != nil {
		return *new(IAutonityValidator), err
	}

	out0 := *abi.ConvertType(out[0], new(IAutonityValidator)).(*IAutonityValidator)

	return out0, err

}

// GetValidator is a free data retrieval call binding the contract method 0x1904bb2e.
//
// Solidity: function getValidator(address _addr) view returns((address,address,address,string,uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256,address,uint256,uint256,uint256,uint256,bytes,uint8,uint256))
func (_IAutonity *IAutonitySession) GetValidator(_addr common.Address) (IAutonityValidator, error) {
	return _IAutonity.Contract.GetValidator(&_IAutonity.CallOpts, _addr)
}

// GetValidator is a free data retrieval call binding the contract method 0x1904bb2e.
//
// Solidity: function getValidator(address _addr) view returns((address,address,address,string,uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256,address,uint256,uint256,uint256,uint256,bytes,uint8,uint256))
func (_IAutonity *IAutonityCallerSession) GetValidator(_addr common.Address) (IAutonityValidator, error) {
	return _IAutonity.Contract.GetValidator(&_IAutonity.CallOpts, _addr)
}

// GetValidatorState is a free data retrieval call binding the contract method 0x5b7d6c36.
//
// Solidity: function getValidatorState(address _addr) view returns(uint8)
func (_IAutonity *IAutonityCaller) GetValidatorState(opts *bind.CallOpts, _addr common.Address) (uint8, error) {
	var out []interface{}
	err := _IAutonity.contract.Call(opts, &out, "getValidatorState", _addr)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// GetValidatorState is a free data retrieval call binding the contract method 0x5b7d6c36.
//
// Solidity: function getValidatorState(address _addr) view returns(uint8)
func (_IAutonity *IAutonitySession) GetValidatorState(_addr common.Address) (uint8, error) {
	return _IAutonity.Contract.GetValidatorState(&_IAutonity.CallOpts, _addr)
}

// GetValidatorState is a free data retrieval call binding the contract method 0x5b7d6c36.
//
// Solidity: function getValidatorState(address _addr) view returns(uint8)
func (_IAutonity *IAutonityCallerSession) GetValidatorState(_addr common.Address) (uint8, error) {
	return _IAutonity.Contract.GetValidatorState(&_IAutonity.CallOpts, _addr)
}

// GetValidators is a free data retrieval call binding the contract method 0xb7ab4db5.
//
// Solidity: function getValidators() view returns(address[])
func (_IAutonity *IAutonityCaller) GetValidators(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _IAutonity.contract.Call(opts, &out, "getValidators")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// GetValidators is a free data retrieval call binding the contract method 0xb7ab4db5.
//
// Solidity: function getValidators() view returns(address[])
func (_IAutonity *IAutonitySession) GetValidators() ([]common.Address, error) {
	return _IAutonity.Contract.GetValidators(&_IAutonity.CallOpts)
}

// GetValidators is a free data retrieval call binding the contract method 0xb7ab4db5.
//
// Solidity: function getValidators() view returns(address[])
func (_IAutonity *IAutonityCallerSession) GetValidators() ([]common.Address, error) {
	return _IAutonity.Contract.GetValidators(&_IAutonity.CallOpts)
}

// GetVersion is a free data retrieval call binding the contract method 0x0d8e6e2c.
//
// Solidity: function getVersion() view returns(uint256)
func (_IAutonity *IAutonityCaller) GetVersion(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IAutonity.contract.Call(opts, &out, "getVersion")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetVersion is a free data retrieval call binding the contract method 0x0d8e6e2c.
//
// Solidity: function getVersion() view returns(uint256)
func (_IAutonity *IAutonitySession) GetVersion() (*big.Int, error) {
	return _IAutonity.Contract.GetVersion(&_IAutonity.CallOpts)
}

// GetVersion is a free data retrieval call binding the contract method 0x0d8e6e2c.
//
// Solidity: function getVersion() view returns(uint256)
func (_IAutonity *IAutonityCallerSession) GetVersion() (*big.Int, error) {
	return _IAutonity.Contract.GetVersion(&_IAutonity.CallOpts)
}

// IsUnbondingReleased is a free data retrieval call binding the contract method 0xe294df7c.
//
// Solidity: function isUnbondingReleased(uint256 _unbondingID) view returns(bool)
func (_IAutonity *IAutonityCaller) IsUnbondingReleased(opts *bind.CallOpts, _unbondingID *big.Int) (bool, error) {
	var out []interface{}
	err := _IAutonity.contract.Call(opts, &out, "isUnbondingReleased", _unbondingID)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsUnbondingReleased is a free data retrieval call binding the contract method 0xe294df7c.
//
// Solidity: function isUnbondingReleased(uint256 _unbondingID) view returns(bool)
func (_IAutonity *IAutonitySession) IsUnbondingReleased(_unbondingID *big.Int) (bool, error) {
	return _IAutonity.Contract.IsUnbondingReleased(&_IAutonity.CallOpts, _unbondingID)
}

// IsUnbondingReleased is a free data retrieval call binding the contract method 0xe294df7c.
//
// Solidity: function isUnbondingReleased(uint256 _unbondingID) view returns(bool)
func (_IAutonity *IAutonityCallerSession) IsUnbondingReleased(_unbondingID *big.Int) (bool, error) {
	return _IAutonity.Contract.IsUnbondingReleased(&_IAutonity.CallOpts, _unbondingID)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_IAutonity *IAutonityCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IAutonity.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_IAutonity *IAutonitySession) TotalSupply() (*big.Int, error) {
	return _IAutonity.Contract.TotalSupply(&_IAutonity.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_IAutonity *IAutonityCallerSession) TotalSupply() (*big.Int, error) {
	return _IAutonity.Contract.TotalSupply(&_IAutonity.CallOpts)
}

// ActivateValidator is a paid mutator transaction binding the contract method 0xb46e5520.
//
// Solidity: function activateValidator(address _address) returns()
func (_IAutonity *IAutonityTransactor) ActivateValidator(opts *bind.TransactOpts, _address common.Address) (*types.Transaction, error) {
	return _IAutonity.contract.Transact(opts, "activateValidator", _address)
}

// ActivateValidator is a paid mutator transaction binding the contract method 0xb46e5520.
//
// Solidity: function activateValidator(address _address) returns()
func (_IAutonity *IAutonitySession) ActivateValidator(_address common.Address) (*types.Transaction, error) {
	return _IAutonity.Contract.ActivateValidator(&_IAutonity.TransactOpts, _address)
}

// ActivateValidator is a paid mutator transaction binding the contract method 0xb46e5520.
//
// Solidity: function activateValidator(address _address) returns()
func (_IAutonity *IAutonityTransactorSession) ActivateValidator(_address common.Address) (*types.Transaction, error) {
	return _IAutonity.Contract.ActivateValidator(&_IAutonity.TransactOpts, _address)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_IAutonity *IAutonityTransactor) Approve(opts *bind.TransactOpts, spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _IAutonity.contract.Transact(opts, "approve", spender, amount)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_IAutonity *IAutonitySession) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _IAutonity.Contract.Approve(&_IAutonity.TransactOpts, spender, amount)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_IAutonity *IAutonityTransactorSession) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _IAutonity.Contract.Approve(&_IAutonity.TransactOpts, spender, amount)
}

// ApproveBonding is a paid mutator transaction binding the contract method 0x50492571.
//
// Solidity: function approveBonding(address _caller, uint256 _amount) returns(bool)
func (_IAutonity *IAutonityTransactor) ApproveBonding(opts *bind.TransactOpts, _caller common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _IAutonity.contract.Transact(opts, "approveBonding", _caller, _amount)
}

// ApproveBonding is a paid mutator transaction binding the contract method 0x50492571.
//
// Solidity: function approveBonding(address _caller, uint256 _amount) returns(bool)
func (_IAutonity *IAutonitySession) ApproveBonding(_caller common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _IAutonity.Contract.ApproveBonding(&_IAutonity.TransactOpts, _caller, _amount)
}

// ApproveBonding is a paid mutator transaction binding the contract method 0x50492571.
//
// Solidity: function approveBonding(address _caller, uint256 _amount) returns(bool)
func (_IAutonity *IAutonityTransactorSession) ApproveBonding(_caller common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _IAutonity.Contract.ApproveBonding(&_IAutonity.TransactOpts, _caller, _amount)
}

// Autobond is a paid mutator transaction binding the contract method 0xf7fcc510.
//
// Solidity: function autobond(address _validator, uint256 _selfBond, uint256 _delegated) returns()
func (_IAutonity *IAutonityTransactor) Autobond(opts *bind.TransactOpts, _validator common.Address, _selfBond *big.Int, _delegated *big.Int) (*types.Transaction, error) {
	return _IAutonity.contract.Transact(opts, "autobond", _validator, _selfBond, _delegated)
}

// Autobond is a paid mutator transaction binding the contract method 0xf7fcc510.
//
// Solidity: function autobond(address _validator, uint256 _selfBond, uint256 _delegated) returns()
func (_IAutonity *IAutonitySession) Autobond(_validator common.Address, _selfBond *big.Int, _delegated *big.Int) (*types.Transaction, error) {
	return _IAutonity.Contract.Autobond(&_IAutonity.TransactOpts, _validator, _selfBond, _delegated)
}

// Autobond is a paid mutator transaction binding the contract method 0xf7fcc510.
//
// Solidity: function autobond(address _validator, uint256 _selfBond, uint256 _delegated) returns()
func (_IAutonity *IAutonityTransactorSession) Autobond(_validator common.Address, _selfBond *big.Int, _delegated *big.Int) (*types.Transaction, error) {
	return _IAutonity.Contract.Autobond(&_IAutonity.TransactOpts, _validator, _selfBond, _delegated)
}

// Bond is a paid mutator transaction binding the contract method 0xa515366a.
//
// Solidity: function bond(address _validator, uint256 _amount) returns(uint256)
func (_IAutonity *IAutonityTransactor) Bond(opts *bind.TransactOpts, _validator common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _IAutonity.contract.Transact(opts, "bond", _validator, _amount)
}

// Bond is a paid mutator transaction binding the contract method 0xa515366a.
//
// Solidity: function bond(address _validator, uint256 _amount) returns(uint256)
func (_IAutonity *IAutonitySession) Bond(_validator common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _IAutonity.Contract.Bond(&_IAutonity.TransactOpts, _validator, _amount)
}

// Bond is a paid mutator transaction binding the contract method 0xa515366a.
//
// Solidity: function bond(address _validator, uint256 _amount) returns(uint256)
func (_IAutonity *IAutonityTransactorSession) Bond(_validator common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _IAutonity.Contract.Bond(&_IAutonity.TransactOpts, _validator, _amount)
}

// BondFrom is a paid mutator transaction binding the contract method 0x41de7400.
//
// Solidity: function bondFrom(address _account, address _validator, uint256 _amount) returns(uint256)
func (_IAutonity *IAutonityTransactor) BondFrom(opts *bind.TransactOpts, _account common.Address, _validator common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _IAutonity.contract.Transact(opts, "bondFrom", _account, _validator, _amount)
}

// BondFrom is a paid mutator transaction binding the contract method 0x41de7400.
//
// Solidity: function bondFrom(address _account, address _validator, uint256 _amount) returns(uint256)
func (_IAutonity *IAutonitySession) BondFrom(_account common.Address, _validator common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _IAutonity.Contract.BondFrom(&_IAutonity.TransactOpts, _account, _validator, _amount)
}

// BondFrom is a paid mutator transaction binding the contract method 0x41de7400.
//
// Solidity: function bondFrom(address _account, address _validator, uint256 _amount) returns(uint256)
func (_IAutonity *IAutonityTransactorSession) BondFrom(_account common.Address, _validator common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _IAutonity.Contract.BondFrom(&_IAutonity.TransactOpts, _account, _validator, _amount)
}

// ChangeCommissionRate is a paid mutator transaction binding the contract method 0x852c4849.
//
// Solidity: function changeCommissionRate(address _validator, uint256 _rate) returns()
func (_IAutonity *IAutonityTransactor) ChangeCommissionRate(opts *bind.TransactOpts, _validator common.Address, _rate *big.Int) (*types.Transaction, error) {
	return _IAutonity.contract.Transact(opts, "changeCommissionRate", _validator, _rate)
}

// ChangeCommissionRate is a paid mutator transaction binding the contract method 0x852c4849.
//
// Solidity: function changeCommissionRate(address _validator, uint256 _rate) returns()
func (_IAutonity *IAutonitySession) ChangeCommissionRate(_validator common.Address, _rate *big.Int) (*types.Transaction, error) {
	return _IAutonity.Contract.ChangeCommissionRate(&_IAutonity.TransactOpts, _validator, _rate)
}

// ChangeCommissionRate is a paid mutator transaction binding the contract method 0x852c4849.
//
// Solidity: function changeCommissionRate(address _validator, uint256 _rate) returns()
func (_IAutonity *IAutonityTransactorSession) ChangeCommissionRate(_validator common.Address, _rate *big.Int) (*types.Transaction, error) {
	return _IAutonity.Contract.ChangeCommissionRate(&_IAutonity.TransactOpts, _validator, _rate)
}

// Jail is a paid mutator transaction binding the contract method 0x154d76d7.
//
// Solidity: function jail(address _nodeAddress, uint256 _jailtime, uint8 _newJailedState) returns(uint256)
func (_IAutonity *IAutonityTransactor) Jail(opts *bind.TransactOpts, _nodeAddress common.Address, _jailtime *big.Int, _newJailedState uint8) (*types.Transaction, error) {
	return _IAutonity.contract.Transact(opts, "jail", _nodeAddress, _jailtime, _newJailedState)
}

// Jail is a paid mutator transaction binding the contract method 0x154d76d7.
//
// Solidity: function jail(address _nodeAddress, uint256 _jailtime, uint8 _newJailedState) returns(uint256)
func (_IAutonity *IAutonitySession) Jail(_nodeAddress common.Address, _jailtime *big.Int, _newJailedState uint8) (*types.Transaction, error) {
	return _IAutonity.Contract.Jail(&_IAutonity.TransactOpts, _nodeAddress, _jailtime, _newJailedState)
}

// Jail is a paid mutator transaction binding the contract method 0x154d76d7.
//
// Solidity: function jail(address _nodeAddress, uint256 _jailtime, uint8 _newJailedState) returns(uint256)
func (_IAutonity *IAutonityTransactorSession) Jail(_nodeAddress common.Address, _jailtime *big.Int, _newJailedState uint8) (*types.Transaction, error) {
	return _IAutonity.Contract.Jail(&_IAutonity.TransactOpts, _nodeAddress, _jailtime, _newJailedState)
}

// Jailbound is a paid mutator transaction binding the contract method 0x8ef8c2fd.
//
// Solidity: function jailbound(address _nodeAddress, uint8 _newJailboundState) returns()
func (_IAutonity *IAutonityTransactor) Jailbound(opts *bind.TransactOpts, _nodeAddress common.Address, _newJailboundState uint8) (*types.Transaction, error) {
	return _IAutonity.contract.Transact(opts, "jailbound", _nodeAddress, _newJailboundState)
}

// Jailbound is a paid mutator transaction binding the contract method 0x8ef8c2fd.
//
// Solidity: function jailbound(address _nodeAddress, uint8 _newJailboundState) returns()
func (_IAutonity *IAutonitySession) Jailbound(_nodeAddress common.Address, _newJailboundState uint8) (*types.Transaction, error) {
	return _IAutonity.Contract.Jailbound(&_IAutonity.TransactOpts, _nodeAddress, _newJailboundState)
}

// Jailbound is a paid mutator transaction binding the contract method 0x8ef8c2fd.
//
// Solidity: function jailbound(address _nodeAddress, uint8 _newJailboundState) returns()
func (_IAutonity *IAutonityTransactorSession) Jailbound(_nodeAddress common.Address, _newJailboundState uint8) (*types.Transaction, error) {
	return _IAutonity.Contract.Jailbound(&_IAutonity.TransactOpts, _nodeAddress, _newJailboundState)
}

// PauseValidator is a paid mutator transaction binding the contract method 0x0ae65e7a.
//
// Solidity: function pauseValidator(address _address) returns()
func (_IAutonity *IAutonityTransactor) PauseValidator(opts *bind.TransactOpts, _address common.Address) (*types.Transaction, error) {
	return _IAutonity.contract.Transact(opts, "pauseValidator", _address)
}

// PauseValidator is a paid mutator transaction binding the contract method 0x0ae65e7a.
//
// Solidity: function pauseValidator(address _address) returns()
func (_IAutonity *IAutonitySession) PauseValidator(_address common.Address) (*types.Transaction, error) {
	return _IAutonity.Contract.PauseValidator(&_IAutonity.TransactOpts, _address)
}

// PauseValidator is a paid mutator transaction binding the contract method 0x0ae65e7a.
//
// Solidity: function pauseValidator(address _address) returns()
func (_IAutonity *IAutonityTransactorSession) PauseValidator(_address common.Address) (*types.Transaction, error) {
	return _IAutonity.Contract.PauseValidator(&_IAutonity.TransactOpts, _address)
}

// RegisterValidator is a paid mutator transaction binding the contract method 0x84467fdb.
//
// Solidity: function registerValidator(string _enode, address _oracleAddress, bytes _consensusKey, bytes _signatures) returns()
func (_IAutonity *IAutonityTransactor) RegisterValidator(opts *bind.TransactOpts, _enode string, _oracleAddress common.Address, _consensusKey []byte, _signatures []byte) (*types.Transaction, error) {
	return _IAutonity.contract.Transact(opts, "registerValidator", _enode, _oracleAddress, _consensusKey, _signatures)
}

// RegisterValidator is a paid mutator transaction binding the contract method 0x84467fdb.
//
// Solidity: function registerValidator(string _enode, address _oracleAddress, bytes _consensusKey, bytes _signatures) returns()
func (_IAutonity *IAutonitySession) RegisterValidator(_enode string, _oracleAddress common.Address, _consensusKey []byte, _signatures []byte) (*types.Transaction, error) {
	return _IAutonity.Contract.RegisterValidator(&_IAutonity.TransactOpts, _enode, _oracleAddress, _consensusKey, _signatures)
}

// RegisterValidator is a paid mutator transaction binding the contract method 0x84467fdb.
//
// Solidity: function registerValidator(string _enode, address _oracleAddress, bytes _consensusKey, bytes _signatures) returns()
func (_IAutonity *IAutonityTransactorSession) RegisterValidator(_enode string, _oracleAddress common.Address, _consensusKey []byte, _signatures []byte) (*types.Transaction, error) {
	return _IAutonity.Contract.RegisterValidator(&_IAutonity.TransactOpts, _enode, _oracleAddress, _consensusKey, _signatures)
}

// Slash is a paid mutator transaction binding the contract method 0x02fb4d85.
//
// Solidity: function slash(address _nodeAddress, uint256 _slashingRate) returns(uint256 slashingAmount)
func (_IAutonity *IAutonityTransactor) Slash(opts *bind.TransactOpts, _nodeAddress common.Address, _slashingRate *big.Int) (*types.Transaction, error) {
	return _IAutonity.contract.Transact(opts, "slash", _nodeAddress, _slashingRate)
}

// Slash is a paid mutator transaction binding the contract method 0x02fb4d85.
//
// Solidity: function slash(address _nodeAddress, uint256 _slashingRate) returns(uint256 slashingAmount)
func (_IAutonity *IAutonitySession) Slash(_nodeAddress common.Address, _slashingRate *big.Int) (*types.Transaction, error) {
	return _IAutonity.Contract.Slash(&_IAutonity.TransactOpts, _nodeAddress, _slashingRate)
}

// Slash is a paid mutator transaction binding the contract method 0x02fb4d85.
//
// Solidity: function slash(address _nodeAddress, uint256 _slashingRate) returns(uint256 slashingAmount)
func (_IAutonity *IAutonityTransactorSession) Slash(_nodeAddress common.Address, _slashingRate *big.Int) (*types.Transaction, error) {
	return _IAutonity.Contract.Slash(&_IAutonity.TransactOpts, _nodeAddress, _slashingRate)
}

// SlashAndJail is a paid mutator transaction binding the contract method 0x122b4122.
//
// Solidity: function slashAndJail(address _nodeAddress, uint256 _slashingRate, uint256 _jailtime, uint8 _newJailedState, uint8 _newJailboundState) returns(uint256 slashingAmount, uint256 jailReleaseBlock, bool isJailbound)
func (_IAutonity *IAutonityTransactor) SlashAndJail(opts *bind.TransactOpts, _nodeAddress common.Address, _slashingRate *big.Int, _jailtime *big.Int, _newJailedState uint8, _newJailboundState uint8) (*types.Transaction, error) {
	return _IAutonity.contract.Transact(opts, "slashAndJail", _nodeAddress, _slashingRate, _jailtime, _newJailedState, _newJailboundState)
}

// SlashAndJail is a paid mutator transaction binding the contract method 0x122b4122.
//
// Solidity: function slashAndJail(address _nodeAddress, uint256 _slashingRate, uint256 _jailtime, uint8 _newJailedState, uint8 _newJailboundState) returns(uint256 slashingAmount, uint256 jailReleaseBlock, bool isJailbound)
func (_IAutonity *IAutonitySession) SlashAndJail(_nodeAddress common.Address, _slashingRate *big.Int, _jailtime *big.Int, _newJailedState uint8, _newJailboundState uint8) (*types.Transaction, error) {
	return _IAutonity.Contract.SlashAndJail(&_IAutonity.TransactOpts, _nodeAddress, _slashingRate, _jailtime, _newJailedState, _newJailboundState)
}

// SlashAndJail is a paid mutator transaction binding the contract method 0x122b4122.
//
// Solidity: function slashAndJail(address _nodeAddress, uint256 _slashingRate, uint256 _jailtime, uint8 _newJailedState, uint8 _newJailboundState) returns(uint256 slashingAmount, uint256 jailReleaseBlock, bool isJailbound)
func (_IAutonity *IAutonityTransactorSession) SlashAndJail(_nodeAddress common.Address, _slashingRate *big.Int, _jailtime *big.Int, _newJailedState uint8, _newJailboundState uint8) (*types.Transaction, error) {
	return _IAutonity.Contract.SlashAndJail(&_IAutonity.TransactOpts, _nodeAddress, _slashingRate, _jailtime, _newJailedState, _newJailboundState)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_IAutonity *IAutonityTransactor) Transfer(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _IAutonity.contract.Transact(opts, "transfer", recipient, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_IAutonity *IAutonitySession) Transfer(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _IAutonity.Contract.Transfer(&_IAutonity.TransactOpts, recipient, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_IAutonity *IAutonityTransactorSession) Transfer(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _IAutonity.Contract.Transfer(&_IAutonity.TransactOpts, recipient, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_IAutonity *IAutonityTransactor) TransferFrom(opts *bind.TransactOpts, sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _IAutonity.contract.Transact(opts, "transferFrom", sender, recipient, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_IAutonity *IAutonitySession) TransferFrom(sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _IAutonity.Contract.TransferFrom(&_IAutonity.TransactOpts, sender, recipient, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_IAutonity *IAutonityTransactorSession) TransferFrom(sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _IAutonity.Contract.TransferFrom(&_IAutonity.TransactOpts, sender, recipient, amount)
}

// Unbond is a paid mutator transaction binding the contract method 0xa5d059ca.
//
// Solidity: function unbond(address _validator, uint256 _amount) returns(uint256)
func (_IAutonity *IAutonityTransactor) Unbond(opts *bind.TransactOpts, _validator common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _IAutonity.contract.Transact(opts, "unbond", _validator, _amount)
}

// Unbond is a paid mutator transaction binding the contract method 0xa5d059ca.
//
// Solidity: function unbond(address _validator, uint256 _amount) returns(uint256)
func (_IAutonity *IAutonitySession) Unbond(_validator common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _IAutonity.Contract.Unbond(&_IAutonity.TransactOpts, _validator, _amount)
}

// Unbond is a paid mutator transaction binding the contract method 0xa5d059ca.
//
// Solidity: function unbond(address _validator, uint256 _amount) returns(uint256)
func (_IAutonity *IAutonityTransactorSession) Unbond(_validator common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _IAutonity.Contract.Unbond(&_IAutonity.TransactOpts, _validator, _amount)
}

// UnbondFrom is a paid mutator transaction binding the contract method 0xa9a7d7c9.
//
// Solidity: function unbondFrom(address _account, address _validator, uint256 _amount) returns(uint256)
func (_IAutonity *IAutonityTransactor) UnbondFrom(opts *bind.TransactOpts, _account common.Address, _validator common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _IAutonity.contract.Transact(opts, "unbondFrom", _account, _validator, _amount)
}

// UnbondFrom is a paid mutator transaction binding the contract method 0xa9a7d7c9.
//
// Solidity: function unbondFrom(address _account, address _validator, uint256 _amount) returns(uint256)
func (_IAutonity *IAutonitySession) UnbondFrom(_account common.Address, _validator common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _IAutonity.Contract.UnbondFrom(&_IAutonity.TransactOpts, _account, _validator, _amount)
}

// UnbondFrom is a paid mutator transaction binding the contract method 0xa9a7d7c9.
//
// Solidity: function unbondFrom(address _account, address _validator, uint256 _amount) returns(uint256)
func (_IAutonity *IAutonityTransactorSession) UnbondFrom(_account common.Address, _validator common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _IAutonity.Contract.UnbondFrom(&_IAutonity.TransactOpts, _account, _validator, _amount)
}

// UpdateEnode is a paid mutator transaction binding the contract method 0x784304b5.
//
// Solidity: function updateEnode(address _nodeAddress, string _enode) returns()
func (_IAutonity *IAutonityTransactor) UpdateEnode(opts *bind.TransactOpts, _nodeAddress common.Address, _enode string) (*types.Transaction, error) {
	return _IAutonity.contract.Transact(opts, "updateEnode", _nodeAddress, _enode)
}

// UpdateEnode is a paid mutator transaction binding the contract method 0x784304b5.
//
// Solidity: function updateEnode(address _nodeAddress, string _enode) returns()
func (_IAutonity *IAutonitySession) UpdateEnode(_nodeAddress common.Address, _enode string) (*types.Transaction, error) {
	return _IAutonity.Contract.UpdateEnode(&_IAutonity.TransactOpts, _nodeAddress, _enode)
}

// UpdateEnode is a paid mutator transaction binding the contract method 0x784304b5.
//
// Solidity: function updateEnode(address _nodeAddress, string _enode) returns()
func (_IAutonity *IAutonityTransactorSession) UpdateEnode(_nodeAddress common.Address, _enode string) (*types.Transaction, error) {
	return _IAutonity.Contract.UpdateEnode(&_IAutonity.TransactOpts, _nodeAddress, _enode)
}

// IAutonityActivatedValidatorIterator is returned from FilterActivatedValidator and is used to iterate over the raw logs and unpacked data for ActivatedValidator events raised by the IAutonity contract.
type IAutonityActivatedValidatorIterator struct {
	Event *IAutonityActivatedValidator // Event containing the contract specifics and raw log

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
func (it *IAutonityActivatedValidatorIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutonityActivatedValidator)
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
		it.Event = new(IAutonityActivatedValidator)
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
func (it *IAutonityActivatedValidatorIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IAutonityActivatedValidatorIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IAutonityActivatedValidator represents a ActivatedValidator event raised by the IAutonity contract.
type IAutonityActivatedValidator struct {
	Treasury       common.Address
	Addr           common.Address
	EffectiveBlock *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterActivatedValidator is a free log retrieval operation binding the contract event 0x60fcbf2d07dc712a93e59fb28f1edb626d7c2497c57ba71a8c0b3999ecb9a3b5.
//
// Solidity: event ActivatedValidator(address indexed treasury, address indexed addr, uint256 effectiveBlock)
func (_IAutonity *IAutonityFilterer) FilterActivatedValidator(opts *bind.FilterOpts, treasury []common.Address, addr []common.Address) (*IAutonityActivatedValidatorIterator, error) {

	var treasuryRule []interface{}
	for _, treasuryItem := range treasury {
		treasuryRule = append(treasuryRule, treasuryItem)
	}
	var addrRule []interface{}
	for _, addrItem := range addr {
		addrRule = append(addrRule, addrItem)
	}

	logs, sub, err := _IAutonity.contract.FilterLogs(opts, "ActivatedValidator", treasuryRule, addrRule)
	if err != nil {
		return nil, err
	}
	return &IAutonityActivatedValidatorIterator{contract: _IAutonity.contract, event: "ActivatedValidator", logs: logs, sub: sub}, nil
}

// WatchActivatedValidator is a free log subscription operation binding the contract event 0x60fcbf2d07dc712a93e59fb28f1edb626d7c2497c57ba71a8c0b3999ecb9a3b5.
//
// Solidity: event ActivatedValidator(address indexed treasury, address indexed addr, uint256 effectiveBlock)
func (_IAutonity *IAutonityFilterer) WatchActivatedValidator(opts *bind.WatchOpts, sink chan<- *IAutonityActivatedValidator, treasury []common.Address, addr []common.Address) (event.Subscription, error) {

	var treasuryRule []interface{}
	for _, treasuryItem := range treasury {
		treasuryRule = append(treasuryRule, treasuryItem)
	}
	var addrRule []interface{}
	for _, addrItem := range addr {
		addrRule = append(addrRule, addrItem)
	}

	logs, sub, err := _IAutonity.contract.WatchLogs(opts, "ActivatedValidator", treasuryRule, addrRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IAutonityActivatedValidator)
				if err := _IAutonity.contract.UnpackLog(event, "ActivatedValidator", log); err != nil {
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

// ParseActivatedValidator is a log parse operation binding the contract event 0x60fcbf2d07dc712a93e59fb28f1edb626d7c2497c57ba71a8c0b3999ecb9a3b5.
//
// Solidity: event ActivatedValidator(address indexed treasury, address indexed addr, uint256 effectiveBlock)
func (_IAutonity *IAutonityFilterer) ParseActivatedValidator(log types.Log) (*IAutonityActivatedValidator, error) {
	event := new(IAutonityActivatedValidator)
	if err := _IAutonity.contract.UnpackLog(event, "ActivatedValidator", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IAutonityApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the IAutonity contract.
type IAutonityApprovalIterator struct {
	Event *IAutonityApproval // Event containing the contract specifics and raw log

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
func (it *IAutonityApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutonityApproval)
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
		it.Event = new(IAutonityApproval)
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
func (it *IAutonityApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IAutonityApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IAutonityApproval represents a Approval event raised by the IAutonity contract.
type IAutonityApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_IAutonity *IAutonityFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*IAutonityApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _IAutonity.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &IAutonityApprovalIterator{contract: _IAutonity.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_IAutonity *IAutonityFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *IAutonityApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _IAutonity.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IAutonityApproval)
				if err := _IAutonity.contract.UnpackLog(event, "Approval", log); err != nil {
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

// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_IAutonity *IAutonityFilterer) ParseApproval(log types.Log) (*IAutonityApproval, error) {
	event := new(IAutonityApproval)
	if err := _IAutonity.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IAutonityBondingApprovalIterator is returned from FilterBondingApproval and is used to iterate over the raw logs and unpacked data for BondingApproval events raised by the IAutonity contract.
type IAutonityBondingApprovalIterator struct {
	Event *IAutonityBondingApproval // Event containing the contract specifics and raw log

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
func (it *IAutonityBondingApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutonityBondingApproval)
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
		it.Event = new(IAutonityBondingApproval)
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
func (it *IAutonityBondingApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IAutonityBondingApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IAutonityBondingApproval represents a BondingApproval event raised by the IAutonity contract.
type IAutonityBondingApproval struct {
	Owner  common.Address
	Caller common.Address
	Value  *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterBondingApproval is a free log retrieval operation binding the contract event 0x78225b0c38a851deb4e94d5318d250e0add756beaec6e85252429c41ed55f8ea.
//
// Solidity: event BondingApproval(address indexed _owner, address indexed _caller, uint256 _value)
func (_IAutonity *IAutonityFilterer) FilterBondingApproval(opts *bind.FilterOpts, _owner []common.Address, _caller []common.Address) (*IAutonityBondingApprovalIterator, error) {

	var _ownerRule []interface{}
	for _, _ownerItem := range _owner {
		_ownerRule = append(_ownerRule, _ownerItem)
	}
	var _callerRule []interface{}
	for _, _callerItem := range _caller {
		_callerRule = append(_callerRule, _callerItem)
	}

	logs, sub, err := _IAutonity.contract.FilterLogs(opts, "BondingApproval", _ownerRule, _callerRule)
	if err != nil {
		return nil, err
	}
	return &IAutonityBondingApprovalIterator{contract: _IAutonity.contract, event: "BondingApproval", logs: logs, sub: sub}, nil
}

// WatchBondingApproval is a free log subscription operation binding the contract event 0x78225b0c38a851deb4e94d5318d250e0add756beaec6e85252429c41ed55f8ea.
//
// Solidity: event BondingApproval(address indexed _owner, address indexed _caller, uint256 _value)
func (_IAutonity *IAutonityFilterer) WatchBondingApproval(opts *bind.WatchOpts, sink chan<- *IAutonityBondingApproval, _owner []common.Address, _caller []common.Address) (event.Subscription, error) {

	var _ownerRule []interface{}
	for _, _ownerItem := range _owner {
		_ownerRule = append(_ownerRule, _ownerItem)
	}
	var _callerRule []interface{}
	for _, _callerItem := range _caller {
		_callerRule = append(_callerRule, _callerItem)
	}

	logs, sub, err := _IAutonity.contract.WatchLogs(opts, "BondingApproval", _ownerRule, _callerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IAutonityBondingApproval)
				if err := _IAutonity.contract.UnpackLog(event, "BondingApproval", log); err != nil {
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

// ParseBondingApproval is a log parse operation binding the contract event 0x78225b0c38a851deb4e94d5318d250e0add756beaec6e85252429c41ed55f8ea.
//
// Solidity: event BondingApproval(address indexed _owner, address indexed _caller, uint256 _value)
func (_IAutonity *IAutonityFilterer) ParseBondingApproval(log types.Log) (*IAutonityBondingApproval, error) {
	event := new(IAutonityBondingApproval)
	if err := _IAutonity.contract.UnpackLog(event, "BondingApproval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IAutonityBondingRejectedIterator is returned from FilterBondingRejected and is used to iterate over the raw logs and unpacked data for BondingRejected events raised by the IAutonity contract.
type IAutonityBondingRejectedIterator struct {
	Event *IAutonityBondingRejected // Event containing the contract specifics and raw log

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
func (it *IAutonityBondingRejectedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutonityBondingRejected)
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
		it.Event = new(IAutonityBondingRejected)
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
func (it *IAutonityBondingRejectedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IAutonityBondingRejectedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IAutonityBondingRejected represents a BondingRejected event raised by the IAutonity contract.
type IAutonityBondingRejected struct {
	Validator common.Address
	Delegator common.Address
	Amount    *big.Int
	State     uint8
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterBondingRejected is a free log retrieval operation binding the contract event 0x1ff2b052afa4bb37ce30d9aaccde416a700b97e632d089111749af937f878342.
//
// Solidity: event BondingRejected(address indexed validator, address indexed delegator, uint256 amount, uint8 state)
func (_IAutonity *IAutonityFilterer) FilterBondingRejected(opts *bind.FilterOpts, validator []common.Address, delegator []common.Address) (*IAutonityBondingRejectedIterator, error) {

	var validatorRule []interface{}
	for _, validatorItem := range validator {
		validatorRule = append(validatorRule, validatorItem)
	}
	var delegatorRule []interface{}
	for _, delegatorItem := range delegator {
		delegatorRule = append(delegatorRule, delegatorItem)
	}

	logs, sub, err := _IAutonity.contract.FilterLogs(opts, "BondingRejected", validatorRule, delegatorRule)
	if err != nil {
		return nil, err
	}
	return &IAutonityBondingRejectedIterator{contract: _IAutonity.contract, event: "BondingRejected", logs: logs, sub: sub}, nil
}

// WatchBondingRejected is a free log subscription operation binding the contract event 0x1ff2b052afa4bb37ce30d9aaccde416a700b97e632d089111749af937f878342.
//
// Solidity: event BondingRejected(address indexed validator, address indexed delegator, uint256 amount, uint8 state)
func (_IAutonity *IAutonityFilterer) WatchBondingRejected(opts *bind.WatchOpts, sink chan<- *IAutonityBondingRejected, validator []common.Address, delegator []common.Address) (event.Subscription, error) {

	var validatorRule []interface{}
	for _, validatorItem := range validator {
		validatorRule = append(validatorRule, validatorItem)
	}
	var delegatorRule []interface{}
	for _, delegatorItem := range delegator {
		delegatorRule = append(delegatorRule, delegatorItem)
	}

	logs, sub, err := _IAutonity.contract.WatchLogs(opts, "BondingRejected", validatorRule, delegatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IAutonityBondingRejected)
				if err := _IAutonity.contract.UnpackLog(event, "BondingRejected", log); err != nil {
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

// ParseBondingRejected is a log parse operation binding the contract event 0x1ff2b052afa4bb37ce30d9aaccde416a700b97e632d089111749af937f878342.
//
// Solidity: event BondingRejected(address indexed validator, address indexed delegator, uint256 amount, uint8 state)
func (_IAutonity *IAutonityFilterer) ParseBondingRejected(log types.Log) (*IAutonityBondingRejected, error) {
	event := new(IAutonityBondingRejected)
	if err := _IAutonity.contract.UnpackLog(event, "BondingRejected", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IAutonityBurnedStakeIterator is returned from FilterBurnedStake and is used to iterate over the raw logs and unpacked data for BurnedStake events raised by the IAutonity contract.
type IAutonityBurnedStakeIterator struct {
	Event *IAutonityBurnedStake // Event containing the contract specifics and raw log

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
func (it *IAutonityBurnedStakeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutonityBurnedStake)
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
		it.Event = new(IAutonityBurnedStake)
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
func (it *IAutonityBurnedStakeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IAutonityBurnedStakeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IAutonityBurnedStake represents a BurnedStake event raised by the IAutonity contract.
type IAutonityBurnedStake struct {
	Addr   common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterBurnedStake is a free log retrieval operation binding the contract event 0x5024dbeedf0c06664c9bd7be836915730c955e936972c020683dadf11d5488a3.
//
// Solidity: event BurnedStake(address indexed addr, uint256 amount)
func (_IAutonity *IAutonityFilterer) FilterBurnedStake(opts *bind.FilterOpts, addr []common.Address) (*IAutonityBurnedStakeIterator, error) {

	var addrRule []interface{}
	for _, addrItem := range addr {
		addrRule = append(addrRule, addrItem)
	}

	logs, sub, err := _IAutonity.contract.FilterLogs(opts, "BurnedStake", addrRule)
	if err != nil {
		return nil, err
	}
	return &IAutonityBurnedStakeIterator{contract: _IAutonity.contract, event: "BurnedStake", logs: logs, sub: sub}, nil
}

// WatchBurnedStake is a free log subscription operation binding the contract event 0x5024dbeedf0c06664c9bd7be836915730c955e936972c020683dadf11d5488a3.
//
// Solidity: event BurnedStake(address indexed addr, uint256 amount)
func (_IAutonity *IAutonityFilterer) WatchBurnedStake(opts *bind.WatchOpts, sink chan<- *IAutonityBurnedStake, addr []common.Address) (event.Subscription, error) {

	var addrRule []interface{}
	for _, addrItem := range addr {
		addrRule = append(addrRule, addrItem)
	}

	logs, sub, err := _IAutonity.contract.WatchLogs(opts, "BurnedStake", addrRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IAutonityBurnedStake)
				if err := _IAutonity.contract.UnpackLog(event, "BurnedStake", log); err != nil {
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

// ParseBurnedStake is a log parse operation binding the contract event 0x5024dbeedf0c06664c9bd7be836915730c955e936972c020683dadf11d5488a3.
//
// Solidity: event BurnedStake(address indexed addr, uint256 amount)
func (_IAutonity *IAutonityFilterer) ParseBurnedStake(log types.Log) (*IAutonityBurnedStake, error) {
	event := new(IAutonityBurnedStake)
	if err := _IAutonity.contract.UnpackLog(event, "BurnedStake", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IAutonityCallFailedIterator is returned from FilterCallFailed and is used to iterate over the raw logs and unpacked data for CallFailed events raised by the IAutonity contract.
type IAutonityCallFailedIterator struct {
	Event *IAutonityCallFailed // Event containing the contract specifics and raw log

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
func (it *IAutonityCallFailedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutonityCallFailed)
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
		it.Event = new(IAutonityCallFailed)
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
func (it *IAutonityCallFailedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IAutonityCallFailedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IAutonityCallFailed represents a CallFailed event raised by the IAutonity contract.
type IAutonityCallFailed struct {
	To              common.Address
	MethodSignature string
	ReturnData      []byte
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterCallFailed is a free log retrieval operation binding the contract event 0x1137d8c966ce69b9630fb2294be011f3d64cc56e91fad7d375f0662568e9d352.
//
// Solidity: event CallFailed(address to, string methodSignature, bytes returnData)
func (_IAutonity *IAutonityFilterer) FilterCallFailed(opts *bind.FilterOpts) (*IAutonityCallFailedIterator, error) {

	logs, sub, err := _IAutonity.contract.FilterLogs(opts, "CallFailed")
	if err != nil {
		return nil, err
	}
	return &IAutonityCallFailedIterator{contract: _IAutonity.contract, event: "CallFailed", logs: logs, sub: sub}, nil
}

// WatchCallFailed is a free log subscription operation binding the contract event 0x1137d8c966ce69b9630fb2294be011f3d64cc56e91fad7d375f0662568e9d352.
//
// Solidity: event CallFailed(address to, string methodSignature, bytes returnData)
func (_IAutonity *IAutonityFilterer) WatchCallFailed(opts *bind.WatchOpts, sink chan<- *IAutonityCallFailed) (event.Subscription, error) {

	logs, sub, err := _IAutonity.contract.WatchLogs(opts, "CallFailed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IAutonityCallFailed)
				if err := _IAutonity.contract.UnpackLog(event, "CallFailed", log); err != nil {
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

// ParseCallFailed is a log parse operation binding the contract event 0x1137d8c966ce69b9630fb2294be011f3d64cc56e91fad7d375f0662568e9d352.
//
// Solidity: event CallFailed(address to, string methodSignature, bytes returnData)
func (_IAutonity *IAutonityFilterer) ParseCallFailed(log types.Log) (*IAutonityCallFailed, error) {
	event := new(IAutonityCallFailed)
	if err := _IAutonity.contract.UnpackLog(event, "CallFailed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IAutonityCommissionRateChangeIterator is returned from FilterCommissionRateChange and is used to iterate over the raw logs and unpacked data for CommissionRateChange events raised by the IAutonity contract.
type IAutonityCommissionRateChangeIterator struct {
	Event *IAutonityCommissionRateChange // Event containing the contract specifics and raw log

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
func (it *IAutonityCommissionRateChangeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutonityCommissionRateChange)
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
		it.Event = new(IAutonityCommissionRateChange)
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
func (it *IAutonityCommissionRateChangeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IAutonityCommissionRateChangeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IAutonityCommissionRateChange represents a CommissionRateChange event raised by the IAutonity contract.
type IAutonityCommissionRateChange struct {
	Validator common.Address
	Rate      *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterCommissionRateChange is a free log retrieval operation binding the contract event 0x4fba51c92fa3d6ad8374d394f6cd5766857552e153d7384a8f23aa4ce9a8a7cf.
//
// Solidity: event CommissionRateChange(address indexed validator, uint256 rate)
func (_IAutonity *IAutonityFilterer) FilterCommissionRateChange(opts *bind.FilterOpts, validator []common.Address) (*IAutonityCommissionRateChangeIterator, error) {

	var validatorRule []interface{}
	for _, validatorItem := range validator {
		validatorRule = append(validatorRule, validatorItem)
	}

	logs, sub, err := _IAutonity.contract.FilterLogs(opts, "CommissionRateChange", validatorRule)
	if err != nil {
		return nil, err
	}
	return &IAutonityCommissionRateChangeIterator{contract: _IAutonity.contract, event: "CommissionRateChange", logs: logs, sub: sub}, nil
}

// WatchCommissionRateChange is a free log subscription operation binding the contract event 0x4fba51c92fa3d6ad8374d394f6cd5766857552e153d7384a8f23aa4ce9a8a7cf.
//
// Solidity: event CommissionRateChange(address indexed validator, uint256 rate)
func (_IAutonity *IAutonityFilterer) WatchCommissionRateChange(opts *bind.WatchOpts, sink chan<- *IAutonityCommissionRateChange, validator []common.Address) (event.Subscription, error) {

	var validatorRule []interface{}
	for _, validatorItem := range validator {
		validatorRule = append(validatorRule, validatorItem)
	}

	logs, sub, err := _IAutonity.contract.WatchLogs(opts, "CommissionRateChange", validatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IAutonityCommissionRateChange)
				if err := _IAutonity.contract.UnpackLog(event, "CommissionRateChange", log); err != nil {
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

// ParseCommissionRateChange is a log parse operation binding the contract event 0x4fba51c92fa3d6ad8374d394f6cd5766857552e153d7384a8f23aa4ce9a8a7cf.
//
// Solidity: event CommissionRateChange(address indexed validator, uint256 rate)
func (_IAutonity *IAutonityFilterer) ParseCommissionRateChange(log types.Log) (*IAutonityCommissionRateChange, error) {
	event := new(IAutonityCommissionRateChange)
	if err := _IAutonity.contract.UnpackLog(event, "CommissionRateChange", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IAutonityEip1559ParamsUpdateIterator is returned from FilterEip1559ParamsUpdate and is used to iterate over the raw logs and unpacked data for Eip1559ParamsUpdate events raised by the IAutonity contract.
type IAutonityEip1559ParamsUpdateIterator struct {
	Event *IAutonityEip1559ParamsUpdate // Event containing the contract specifics and raw log

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
func (it *IAutonityEip1559ParamsUpdateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutonityEip1559ParamsUpdate)
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
		it.Event = new(IAutonityEip1559ParamsUpdate)
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
func (it *IAutonityEip1559ParamsUpdateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IAutonityEip1559ParamsUpdateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IAutonityEip1559ParamsUpdate represents a Eip1559ParamsUpdate event raised by the IAutonity contract.
type IAutonityEip1559ParamsUpdate struct {
	OldParams IAutonityEip1559
	NewParams IAutonityEip1559
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterEip1559ParamsUpdate is a free log retrieval operation binding the contract event 0xa1a63f0900be3abc55a58106fb99ccb53c1cdeaee2488038528f520d142563cc.
//
// Solidity: event Eip1559ParamsUpdate((uint256,uint256,uint256,uint256) oldParams, (uint256,uint256,uint256,uint256) newParams)
func (_IAutonity *IAutonityFilterer) FilterEip1559ParamsUpdate(opts *bind.FilterOpts) (*IAutonityEip1559ParamsUpdateIterator, error) {

	logs, sub, err := _IAutonity.contract.FilterLogs(opts, "Eip1559ParamsUpdate")
	if err != nil {
		return nil, err
	}
	return &IAutonityEip1559ParamsUpdateIterator{contract: _IAutonity.contract, event: "Eip1559ParamsUpdate", logs: logs, sub: sub}, nil
}

// WatchEip1559ParamsUpdate is a free log subscription operation binding the contract event 0xa1a63f0900be3abc55a58106fb99ccb53c1cdeaee2488038528f520d142563cc.
//
// Solidity: event Eip1559ParamsUpdate((uint256,uint256,uint256,uint256) oldParams, (uint256,uint256,uint256,uint256) newParams)
func (_IAutonity *IAutonityFilterer) WatchEip1559ParamsUpdate(opts *bind.WatchOpts, sink chan<- *IAutonityEip1559ParamsUpdate) (event.Subscription, error) {

	logs, sub, err := _IAutonity.contract.WatchLogs(opts, "Eip1559ParamsUpdate")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IAutonityEip1559ParamsUpdate)
				if err := _IAutonity.contract.UnpackLog(event, "Eip1559ParamsUpdate", log); err != nil {
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

// ParseEip1559ParamsUpdate is a log parse operation binding the contract event 0xa1a63f0900be3abc55a58106fb99ccb53c1cdeaee2488038528f520d142563cc.
//
// Solidity: event Eip1559ParamsUpdate((uint256,uint256,uint256,uint256) oldParams, (uint256,uint256,uint256,uint256) newParams)
func (_IAutonity *IAutonityFilterer) ParseEip1559ParamsUpdate(log types.Log) (*IAutonityEip1559ParamsUpdate, error) {
	event := new(IAutonityEip1559ParamsUpdate)
	if err := _IAutonity.contract.UnpackLog(event, "Eip1559ParamsUpdate", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IAutonityEnodeUpdateIterator is returned from FilterEnodeUpdate and is used to iterate over the raw logs and unpacked data for EnodeUpdate events raised by the IAutonity contract.
type IAutonityEnodeUpdateIterator struct {
	Event *IAutonityEnodeUpdate // Event containing the contract specifics and raw log

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
func (it *IAutonityEnodeUpdateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutonityEnodeUpdate)
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
		it.Event = new(IAutonityEnodeUpdate)
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
func (it *IAutonityEnodeUpdateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IAutonityEnodeUpdateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IAutonityEnodeUpdate represents a EnodeUpdate event raised by the IAutonity contract.
type IAutonityEnodeUpdate struct {
	Validator common.Address
	OldEnode  string
	NewEnode  string
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterEnodeUpdate is a free log retrieval operation binding the contract event 0x0693af84e99c5d55723f139a6badcfffa980bb354cf9c32bd5da745e9edaa46b.
//
// Solidity: event EnodeUpdate(address validator, string oldEnode, string newEnode)
func (_IAutonity *IAutonityFilterer) FilterEnodeUpdate(opts *bind.FilterOpts) (*IAutonityEnodeUpdateIterator, error) {

	logs, sub, err := _IAutonity.contract.FilterLogs(opts, "EnodeUpdate")
	if err != nil {
		return nil, err
	}
	return &IAutonityEnodeUpdateIterator{contract: _IAutonity.contract, event: "EnodeUpdate", logs: logs, sub: sub}, nil
}

// WatchEnodeUpdate is a free log subscription operation binding the contract event 0x0693af84e99c5d55723f139a6badcfffa980bb354cf9c32bd5da745e9edaa46b.
//
// Solidity: event EnodeUpdate(address validator, string oldEnode, string newEnode)
func (_IAutonity *IAutonityFilterer) WatchEnodeUpdate(opts *bind.WatchOpts, sink chan<- *IAutonityEnodeUpdate) (event.Subscription, error) {

	logs, sub, err := _IAutonity.contract.WatchLogs(opts, "EnodeUpdate")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IAutonityEnodeUpdate)
				if err := _IAutonity.contract.UnpackLog(event, "EnodeUpdate", log); err != nil {
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

// ParseEnodeUpdate is a log parse operation binding the contract event 0x0693af84e99c5d55723f139a6badcfffa980bb354cf9c32bd5da745e9edaa46b.
//
// Solidity: event EnodeUpdate(address validator, string oldEnode, string newEnode)
func (_IAutonity *IAutonityFilterer) ParseEnodeUpdate(log types.Log) (*IAutonityEnodeUpdate, error) {
	event := new(IAutonityEnodeUpdate)
	if err := _IAutonity.contract.UnpackLog(event, "EnodeUpdate", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IAutonityEpochPeriodUpdatedIterator is returned from FilterEpochPeriodUpdated and is used to iterate over the raw logs and unpacked data for EpochPeriodUpdated events raised by the IAutonity contract.
type IAutonityEpochPeriodUpdatedIterator struct {
	Event *IAutonityEpochPeriodUpdated // Event containing the contract specifics and raw log

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
func (it *IAutonityEpochPeriodUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutonityEpochPeriodUpdated)
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
		it.Event = new(IAutonityEpochPeriodUpdated)
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
func (it *IAutonityEpochPeriodUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IAutonityEpochPeriodUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IAutonityEpochPeriodUpdated represents a EpochPeriodUpdated event raised by the IAutonity contract.
type IAutonityEpochPeriodUpdated struct {
	Period         *big.Int
	AppliedAtBlock *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterEpochPeriodUpdated is a free log retrieval operation binding the contract event 0x2eea6438d890c8603d4df81ad1bad2a4ea45c02b4837165f461ff3c81603abc7.
//
// Solidity: event EpochPeriodUpdated(uint256 period, uint256 appliedAtBlock)
func (_IAutonity *IAutonityFilterer) FilterEpochPeriodUpdated(opts *bind.FilterOpts) (*IAutonityEpochPeriodUpdatedIterator, error) {

	logs, sub, err := _IAutonity.contract.FilterLogs(opts, "EpochPeriodUpdated")
	if err != nil {
		return nil, err
	}
	return &IAutonityEpochPeriodUpdatedIterator{contract: _IAutonity.contract, event: "EpochPeriodUpdated", logs: logs, sub: sub}, nil
}

// WatchEpochPeriodUpdated is a free log subscription operation binding the contract event 0x2eea6438d890c8603d4df81ad1bad2a4ea45c02b4837165f461ff3c81603abc7.
//
// Solidity: event EpochPeriodUpdated(uint256 period, uint256 appliedAtBlock)
func (_IAutonity *IAutonityFilterer) WatchEpochPeriodUpdated(opts *bind.WatchOpts, sink chan<- *IAutonityEpochPeriodUpdated) (event.Subscription, error) {

	logs, sub, err := _IAutonity.contract.WatchLogs(opts, "EpochPeriodUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IAutonityEpochPeriodUpdated)
				if err := _IAutonity.contract.UnpackLog(event, "EpochPeriodUpdated", log); err != nil {
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

// ParseEpochPeriodUpdated is a log parse operation binding the contract event 0x2eea6438d890c8603d4df81ad1bad2a4ea45c02b4837165f461ff3c81603abc7.
//
// Solidity: event EpochPeriodUpdated(uint256 period, uint256 appliedAtBlock)
func (_IAutonity *IAutonityFilterer) ParseEpochPeriodUpdated(log types.Log) (*IAutonityEpochPeriodUpdated, error) {
	event := new(IAutonityEpochPeriodUpdated)
	if err := _IAutonity.contract.UnpackLog(event, "EpochPeriodUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IAutonityMintedStakeIterator is returned from FilterMintedStake and is used to iterate over the raw logs and unpacked data for MintedStake events raised by the IAutonity contract.
type IAutonityMintedStakeIterator struct {
	Event *IAutonityMintedStake // Event containing the contract specifics and raw log

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
func (it *IAutonityMintedStakeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutonityMintedStake)
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
		it.Event = new(IAutonityMintedStake)
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
func (it *IAutonityMintedStakeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IAutonityMintedStakeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IAutonityMintedStake represents a MintedStake event raised by the IAutonity contract.
type IAutonityMintedStake struct {
	Addr   common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterMintedStake is a free log retrieval operation binding the contract event 0x48490b4407bb949b708ec5f514b4167f08f4969baaf78d53b05028adf369bfcf.
//
// Solidity: event MintedStake(address indexed addr, uint256 amount)
func (_IAutonity *IAutonityFilterer) FilterMintedStake(opts *bind.FilterOpts, addr []common.Address) (*IAutonityMintedStakeIterator, error) {

	var addrRule []interface{}
	for _, addrItem := range addr {
		addrRule = append(addrRule, addrItem)
	}

	logs, sub, err := _IAutonity.contract.FilterLogs(opts, "MintedStake", addrRule)
	if err != nil {
		return nil, err
	}
	return &IAutonityMintedStakeIterator{contract: _IAutonity.contract, event: "MintedStake", logs: logs, sub: sub}, nil
}

// WatchMintedStake is a free log subscription operation binding the contract event 0x48490b4407bb949b708ec5f514b4167f08f4969baaf78d53b05028adf369bfcf.
//
// Solidity: event MintedStake(address indexed addr, uint256 amount)
func (_IAutonity *IAutonityFilterer) WatchMintedStake(opts *bind.WatchOpts, sink chan<- *IAutonityMintedStake, addr []common.Address) (event.Subscription, error) {

	var addrRule []interface{}
	for _, addrItem := range addr {
		addrRule = append(addrRule, addrItem)
	}

	logs, sub, err := _IAutonity.contract.WatchLogs(opts, "MintedStake", addrRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IAutonityMintedStake)
				if err := _IAutonity.contract.UnpackLog(event, "MintedStake", log); err != nil {
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

// ParseMintedStake is a log parse operation binding the contract event 0x48490b4407bb949b708ec5f514b4167f08f4969baaf78d53b05028adf369bfcf.
//
// Solidity: event MintedStake(address indexed addr, uint256 amount)
func (_IAutonity *IAutonityFilterer) ParseMintedStake(log types.Log) (*IAutonityMintedStake, error) {
	event := new(IAutonityMintedStake)
	if err := _IAutonity.contract.UnpackLog(event, "MintedStake", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IAutonityNewBondingRequestIterator is returned from FilterNewBondingRequest and is used to iterate over the raw logs and unpacked data for NewBondingRequest events raised by the IAutonity contract.
type IAutonityNewBondingRequestIterator struct {
	Event *IAutonityNewBondingRequest // Event containing the contract specifics and raw log

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
func (it *IAutonityNewBondingRequestIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutonityNewBondingRequest)
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
		it.Event = new(IAutonityNewBondingRequest)
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
func (it *IAutonityNewBondingRequestIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IAutonityNewBondingRequestIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IAutonityNewBondingRequest represents a NewBondingRequest event raised by the IAutonity contract.
type IAutonityNewBondingRequest struct {
	Validator     common.Address
	Delegator     common.Address
	Caller        common.Address
	SelfBonded    bool
	Amount        *big.Int
	HeadBondingID *big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterNewBondingRequest is a free log retrieval operation binding the contract event 0xcedee1b81e707f05d3f5f6b965053c24ce7759214e4886d87c02e0261367d1f2.
//
// Solidity: event NewBondingRequest(address indexed validator, address indexed delegator, address indexed caller, bool selfBonded, uint256 amount, uint256 headBondingID)
func (_IAutonity *IAutonityFilterer) FilterNewBondingRequest(opts *bind.FilterOpts, validator []common.Address, delegator []common.Address, caller []common.Address) (*IAutonityNewBondingRequestIterator, error) {

	var validatorRule []interface{}
	for _, validatorItem := range validator {
		validatorRule = append(validatorRule, validatorItem)
	}
	var delegatorRule []interface{}
	for _, delegatorItem := range delegator {
		delegatorRule = append(delegatorRule, delegatorItem)
	}
	var callerRule []interface{}
	for _, callerItem := range caller {
		callerRule = append(callerRule, callerItem)
	}

	logs, sub, err := _IAutonity.contract.FilterLogs(opts, "NewBondingRequest", validatorRule, delegatorRule, callerRule)
	if err != nil {
		return nil, err
	}
	return &IAutonityNewBondingRequestIterator{contract: _IAutonity.contract, event: "NewBondingRequest", logs: logs, sub: sub}, nil
}

// WatchNewBondingRequest is a free log subscription operation binding the contract event 0xcedee1b81e707f05d3f5f6b965053c24ce7759214e4886d87c02e0261367d1f2.
//
// Solidity: event NewBondingRequest(address indexed validator, address indexed delegator, address indexed caller, bool selfBonded, uint256 amount, uint256 headBondingID)
func (_IAutonity *IAutonityFilterer) WatchNewBondingRequest(opts *bind.WatchOpts, sink chan<- *IAutonityNewBondingRequest, validator []common.Address, delegator []common.Address, caller []common.Address) (event.Subscription, error) {

	var validatorRule []interface{}
	for _, validatorItem := range validator {
		validatorRule = append(validatorRule, validatorItem)
	}
	var delegatorRule []interface{}
	for _, delegatorItem := range delegator {
		delegatorRule = append(delegatorRule, delegatorItem)
	}
	var callerRule []interface{}
	for _, callerItem := range caller {
		callerRule = append(callerRule, callerItem)
	}

	logs, sub, err := _IAutonity.contract.WatchLogs(opts, "NewBondingRequest", validatorRule, delegatorRule, callerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IAutonityNewBondingRequest)
				if err := _IAutonity.contract.UnpackLog(event, "NewBondingRequest", log); err != nil {
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

// ParseNewBondingRequest is a log parse operation binding the contract event 0xcedee1b81e707f05d3f5f6b965053c24ce7759214e4886d87c02e0261367d1f2.
//
// Solidity: event NewBondingRequest(address indexed validator, address indexed delegator, address indexed caller, bool selfBonded, uint256 amount, uint256 headBondingID)
func (_IAutonity *IAutonityFilterer) ParseNewBondingRequest(log types.Log) (*IAutonityNewBondingRequest, error) {
	event := new(IAutonityNewBondingRequest)
	if err := _IAutonity.contract.UnpackLog(event, "NewBondingRequest", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IAutonityNewEpochIterator is returned from FilterNewEpoch and is used to iterate over the raw logs and unpacked data for NewEpoch events raised by the IAutonity contract.
type IAutonityNewEpochIterator struct {
	Event *IAutonityNewEpoch // Event containing the contract specifics and raw log

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
func (it *IAutonityNewEpochIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutonityNewEpoch)
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
		it.Event = new(IAutonityNewEpoch)
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
func (it *IAutonityNewEpochIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IAutonityNewEpochIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IAutonityNewEpoch represents a NewEpoch event raised by the IAutonity contract.
type IAutonityNewEpoch struct {
	Epoch            *big.Int
	InflationReserve *big.Int
	StakeCirculating *big.Int
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterNewEpoch is a free log retrieval operation binding the contract event 0x3bb7b347508b7c148ec2094ac60d2e3d8b7595421025643f08b45cb78b326b58.
//
// Solidity: event NewEpoch(uint256 epoch, uint256 inflationReserve, uint256 stakeCirculating)
func (_IAutonity *IAutonityFilterer) FilterNewEpoch(opts *bind.FilterOpts) (*IAutonityNewEpochIterator, error) {

	logs, sub, err := _IAutonity.contract.FilterLogs(opts, "NewEpoch")
	if err != nil {
		return nil, err
	}
	return &IAutonityNewEpochIterator{contract: _IAutonity.contract, event: "NewEpoch", logs: logs, sub: sub}, nil
}

// WatchNewEpoch is a free log subscription operation binding the contract event 0x3bb7b347508b7c148ec2094ac60d2e3d8b7595421025643f08b45cb78b326b58.
//
// Solidity: event NewEpoch(uint256 epoch, uint256 inflationReserve, uint256 stakeCirculating)
func (_IAutonity *IAutonityFilterer) WatchNewEpoch(opts *bind.WatchOpts, sink chan<- *IAutonityNewEpoch) (event.Subscription, error) {

	logs, sub, err := _IAutonity.contract.WatchLogs(opts, "NewEpoch")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IAutonityNewEpoch)
				if err := _IAutonity.contract.UnpackLog(event, "NewEpoch", log); err != nil {
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

// ParseNewEpoch is a log parse operation binding the contract event 0x3bb7b347508b7c148ec2094ac60d2e3d8b7595421025643f08b45cb78b326b58.
//
// Solidity: event NewEpoch(uint256 epoch, uint256 inflationReserve, uint256 stakeCirculating)
func (_IAutonity *IAutonityFilterer) ParseNewEpoch(log types.Log) (*IAutonityNewEpoch, error) {
	event := new(IAutonityNewEpoch)
	if err := _IAutonity.contract.UnpackLog(event, "NewEpoch", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IAutonityNewUnbondingRequestIterator is returned from FilterNewUnbondingRequest and is used to iterate over the raw logs and unpacked data for NewUnbondingRequest events raised by the IAutonity contract.
type IAutonityNewUnbondingRequestIterator struct {
	Event *IAutonityNewUnbondingRequest // Event containing the contract specifics and raw log

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
func (it *IAutonityNewUnbondingRequestIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutonityNewUnbondingRequest)
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
		it.Event = new(IAutonityNewUnbondingRequest)
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
func (it *IAutonityNewUnbondingRequestIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IAutonityNewUnbondingRequestIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IAutonityNewUnbondingRequest represents a NewUnbondingRequest event raised by the IAutonity contract.
type IAutonityNewUnbondingRequest struct {
	Validator       common.Address
	Delegator       common.Address
	Caller          common.Address
	SelfBonded      bool
	Amount          *big.Int
	HeadUnbondingID *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterNewUnbondingRequest is a free log retrieval operation binding the contract event 0x58cac3f318ed4be30304429f66e0d6794a95041b92ebe8d35ec0d0abcf1b593c.
//
// Solidity: event NewUnbondingRequest(address indexed validator, address indexed delegator, address indexed caller, bool selfBonded, uint256 amount, uint256 headUnbondingID)
func (_IAutonity *IAutonityFilterer) FilterNewUnbondingRequest(opts *bind.FilterOpts, validator []common.Address, delegator []common.Address, caller []common.Address) (*IAutonityNewUnbondingRequestIterator, error) {

	var validatorRule []interface{}
	for _, validatorItem := range validator {
		validatorRule = append(validatorRule, validatorItem)
	}
	var delegatorRule []interface{}
	for _, delegatorItem := range delegator {
		delegatorRule = append(delegatorRule, delegatorItem)
	}
	var callerRule []interface{}
	for _, callerItem := range caller {
		callerRule = append(callerRule, callerItem)
	}

	logs, sub, err := _IAutonity.contract.FilterLogs(opts, "NewUnbondingRequest", validatorRule, delegatorRule, callerRule)
	if err != nil {
		return nil, err
	}
	return &IAutonityNewUnbondingRequestIterator{contract: _IAutonity.contract, event: "NewUnbondingRequest", logs: logs, sub: sub}, nil
}

// WatchNewUnbondingRequest is a free log subscription operation binding the contract event 0x58cac3f318ed4be30304429f66e0d6794a95041b92ebe8d35ec0d0abcf1b593c.
//
// Solidity: event NewUnbondingRequest(address indexed validator, address indexed delegator, address indexed caller, bool selfBonded, uint256 amount, uint256 headUnbondingID)
func (_IAutonity *IAutonityFilterer) WatchNewUnbondingRequest(opts *bind.WatchOpts, sink chan<- *IAutonityNewUnbondingRequest, validator []common.Address, delegator []common.Address, caller []common.Address) (event.Subscription, error) {

	var validatorRule []interface{}
	for _, validatorItem := range validator {
		validatorRule = append(validatorRule, validatorItem)
	}
	var delegatorRule []interface{}
	for _, delegatorItem := range delegator {
		delegatorRule = append(delegatorRule, delegatorItem)
	}
	var callerRule []interface{}
	for _, callerItem := range caller {
		callerRule = append(callerRule, callerItem)
	}

	logs, sub, err := _IAutonity.contract.WatchLogs(opts, "NewUnbondingRequest", validatorRule, delegatorRule, callerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IAutonityNewUnbondingRequest)
				if err := _IAutonity.contract.UnpackLog(event, "NewUnbondingRequest", log); err != nil {
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

// ParseNewUnbondingRequest is a log parse operation binding the contract event 0x58cac3f318ed4be30304429f66e0d6794a95041b92ebe8d35ec0d0abcf1b593c.
//
// Solidity: event NewUnbondingRequest(address indexed validator, address indexed delegator, address indexed caller, bool selfBonded, uint256 amount, uint256 headUnbondingID)
func (_IAutonity *IAutonityFilterer) ParseNewUnbondingRequest(log types.Log) (*IAutonityNewUnbondingRequest, error) {
	event := new(IAutonityNewUnbondingRequest)
	if err := _IAutonity.contract.UnpackLog(event, "NewUnbondingRequest", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IAutonityPausedValidatorIterator is returned from FilterPausedValidator and is used to iterate over the raw logs and unpacked data for PausedValidator events raised by the IAutonity contract.
type IAutonityPausedValidatorIterator struct {
	Event *IAutonityPausedValidator // Event containing the contract specifics and raw log

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
func (it *IAutonityPausedValidatorIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutonityPausedValidator)
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
		it.Event = new(IAutonityPausedValidator)
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
func (it *IAutonityPausedValidatorIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IAutonityPausedValidatorIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IAutonityPausedValidator represents a PausedValidator event raised by the IAutonity contract.
type IAutonityPausedValidator struct {
	Treasury       common.Address
	Addr           common.Address
	EffectiveBlock *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterPausedValidator is a free log retrieval operation binding the contract event 0x75bdcdbe540758778e669d108fbcb7ede734f27f46e4e5525eeb8ecf91849a9c.
//
// Solidity: event PausedValidator(address indexed treasury, address indexed addr, uint256 effectiveBlock)
func (_IAutonity *IAutonityFilterer) FilterPausedValidator(opts *bind.FilterOpts, treasury []common.Address, addr []common.Address) (*IAutonityPausedValidatorIterator, error) {

	var treasuryRule []interface{}
	for _, treasuryItem := range treasury {
		treasuryRule = append(treasuryRule, treasuryItem)
	}
	var addrRule []interface{}
	for _, addrItem := range addr {
		addrRule = append(addrRule, addrItem)
	}

	logs, sub, err := _IAutonity.contract.FilterLogs(opts, "PausedValidator", treasuryRule, addrRule)
	if err != nil {
		return nil, err
	}
	return &IAutonityPausedValidatorIterator{contract: _IAutonity.contract, event: "PausedValidator", logs: logs, sub: sub}, nil
}

// WatchPausedValidator is a free log subscription operation binding the contract event 0x75bdcdbe540758778e669d108fbcb7ede734f27f46e4e5525eeb8ecf91849a9c.
//
// Solidity: event PausedValidator(address indexed treasury, address indexed addr, uint256 effectiveBlock)
func (_IAutonity *IAutonityFilterer) WatchPausedValidator(opts *bind.WatchOpts, sink chan<- *IAutonityPausedValidator, treasury []common.Address, addr []common.Address) (event.Subscription, error) {

	var treasuryRule []interface{}
	for _, treasuryItem := range treasury {
		treasuryRule = append(treasuryRule, treasuryItem)
	}
	var addrRule []interface{}
	for _, addrItem := range addr {
		addrRule = append(addrRule, addrItem)
	}

	logs, sub, err := _IAutonity.contract.WatchLogs(opts, "PausedValidator", treasuryRule, addrRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IAutonityPausedValidator)
				if err := _IAutonity.contract.UnpackLog(event, "PausedValidator", log); err != nil {
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

// ParsePausedValidator is a log parse operation binding the contract event 0x75bdcdbe540758778e669d108fbcb7ede734f27f46e4e5525eeb8ecf91849a9c.
//
// Solidity: event PausedValidator(address indexed treasury, address indexed addr, uint256 effectiveBlock)
func (_IAutonity *IAutonityFilterer) ParsePausedValidator(log types.Log) (*IAutonityPausedValidator, error) {
	event := new(IAutonityPausedValidator)
	if err := _IAutonity.contract.UnpackLog(event, "PausedValidator", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IAutonityRegisteredValidatorIterator is returned from FilterRegisteredValidator and is used to iterate over the raw logs and unpacked data for RegisteredValidator events raised by the IAutonity contract.
type IAutonityRegisteredValidatorIterator struct {
	Event *IAutonityRegisteredValidator // Event containing the contract specifics and raw log

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
func (it *IAutonityRegisteredValidatorIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutonityRegisteredValidator)
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
		it.Event = new(IAutonityRegisteredValidator)
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
func (it *IAutonityRegisteredValidatorIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IAutonityRegisteredValidatorIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IAutonityRegisteredValidator represents a RegisteredValidator event raised by the IAutonity contract.
type IAutonityRegisteredValidator struct {
	Treasury            common.Address
	Addr                common.Address
	OracleAddress       common.Address
	Enode               string
	LiquidStateContract common.Address
	Raw                 types.Log // Blockchain specific contextual infos
}

// FilterRegisteredValidator is a free log retrieval operation binding the contract event 0x8ad8bd2eb6950e5f332fd3a6dca48cb358ecfe3057848902b98cbdfe455c915c.
//
// Solidity: event RegisteredValidator(address treasury, address addr, address oracleAddress, string enode, address liquidStateContract)
func (_IAutonity *IAutonityFilterer) FilterRegisteredValidator(opts *bind.FilterOpts) (*IAutonityRegisteredValidatorIterator, error) {

	logs, sub, err := _IAutonity.contract.FilterLogs(opts, "RegisteredValidator")
	if err != nil {
		return nil, err
	}
	return &IAutonityRegisteredValidatorIterator{contract: _IAutonity.contract, event: "RegisteredValidator", logs: logs, sub: sub}, nil
}

// WatchRegisteredValidator is a free log subscription operation binding the contract event 0x8ad8bd2eb6950e5f332fd3a6dca48cb358ecfe3057848902b98cbdfe455c915c.
//
// Solidity: event RegisteredValidator(address treasury, address addr, address oracleAddress, string enode, address liquidStateContract)
func (_IAutonity *IAutonityFilterer) WatchRegisteredValidator(opts *bind.WatchOpts, sink chan<- *IAutonityRegisteredValidator) (event.Subscription, error) {

	logs, sub, err := _IAutonity.contract.WatchLogs(opts, "RegisteredValidator")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IAutonityRegisteredValidator)
				if err := _IAutonity.contract.UnpackLog(event, "RegisteredValidator", log); err != nil {
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

// ParseRegisteredValidator is a log parse operation binding the contract event 0x8ad8bd2eb6950e5f332fd3a6dca48cb358ecfe3057848902b98cbdfe455c915c.
//
// Solidity: event RegisteredValidator(address treasury, address addr, address oracleAddress, string enode, address liquidStateContract)
func (_IAutonity *IAutonityFilterer) ParseRegisteredValidator(log types.Log) (*IAutonityRegisteredValidator, error) {
	event := new(IAutonityRegisteredValidator)
	if err := _IAutonity.contract.UnpackLog(event, "RegisteredValidator", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IAutonityRewardedIterator is returned from FilterRewarded and is used to iterate over the raw logs and unpacked data for Rewarded events raised by the IAutonity contract.
type IAutonityRewardedIterator struct {
	Event *IAutonityRewarded // Event containing the contract specifics and raw log

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
func (it *IAutonityRewardedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutonityRewarded)
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
		it.Event = new(IAutonityRewarded)
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
func (it *IAutonityRewardedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IAutonityRewardedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IAutonityRewarded represents a Rewarded event raised by the IAutonity contract.
type IAutonityRewarded struct {
	Addr               common.Address
	AtnSelfAmount      *big.Int
	AtnDelegatedAmount *big.Int
	NtnSelfAmount      *big.Int
	NtnDelegatedAmount *big.Int
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterRewarded is a free log retrieval operation binding the contract event 0x5d4e9c803c111cc5e507b1cd7ffb638d0084c4d0696bb89a39434637ce31b4d5.
//
// Solidity: event Rewarded(address indexed addr, uint256 atnSelfAmount, uint256 atnDelegatedAmount, uint256 ntnSelfAmount, uint256 ntnDelegatedAmount)
func (_IAutonity *IAutonityFilterer) FilterRewarded(opts *bind.FilterOpts, addr []common.Address) (*IAutonityRewardedIterator, error) {

	var addrRule []interface{}
	for _, addrItem := range addr {
		addrRule = append(addrRule, addrItem)
	}

	logs, sub, err := _IAutonity.contract.FilterLogs(opts, "Rewarded", addrRule)
	if err != nil {
		return nil, err
	}
	return &IAutonityRewardedIterator{contract: _IAutonity.contract, event: "Rewarded", logs: logs, sub: sub}, nil
}

// WatchRewarded is a free log subscription operation binding the contract event 0x5d4e9c803c111cc5e507b1cd7ffb638d0084c4d0696bb89a39434637ce31b4d5.
//
// Solidity: event Rewarded(address indexed addr, uint256 atnSelfAmount, uint256 atnDelegatedAmount, uint256 ntnSelfAmount, uint256 ntnDelegatedAmount)
func (_IAutonity *IAutonityFilterer) WatchRewarded(opts *bind.WatchOpts, sink chan<- *IAutonityRewarded, addr []common.Address) (event.Subscription, error) {

	var addrRule []interface{}
	for _, addrItem := range addr {
		addrRule = append(addrRule, addrItem)
	}

	logs, sub, err := _IAutonity.contract.WatchLogs(opts, "Rewarded", addrRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IAutonityRewarded)
				if err := _IAutonity.contract.UnpackLog(event, "Rewarded", log); err != nil {
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

// ParseRewarded is a log parse operation binding the contract event 0x5d4e9c803c111cc5e507b1cd7ffb638d0084c4d0696bb89a39434637ce31b4d5.
//
// Solidity: event Rewarded(address indexed addr, uint256 atnSelfAmount, uint256 atnDelegatedAmount, uint256 ntnSelfAmount, uint256 ntnDelegatedAmount)
func (_IAutonity *IAutonityFilterer) ParseRewarded(log types.Log) (*IAutonityRewarded, error) {
	event := new(IAutonityRewarded)
	if err := _IAutonity.contract.UnpackLog(event, "Rewarded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IAutonityTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the IAutonity contract.
type IAutonityTransferIterator struct {
	Event *IAutonityTransfer // Event containing the contract specifics and raw log

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
func (it *IAutonityTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutonityTransfer)
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
		it.Event = new(IAutonityTransfer)
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
func (it *IAutonityTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IAutonityTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IAutonityTransfer represents a Transfer event raised by the IAutonity contract.
type IAutonityTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_IAutonity *IAutonityFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*IAutonityTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _IAutonity.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &IAutonityTransferIterator{contract: _IAutonity.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_IAutonity *IAutonityFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *IAutonityTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _IAutonity.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IAutonityTransfer)
				if err := _IAutonity.contract.UnpackLog(event, "Transfer", log); err != nil {
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
func (_IAutonity *IAutonityFilterer) ParseTransfer(log types.Log) (*IAutonityTransfer, error) {
	event := new(IAutonityTransfer)
	if err := _IAutonity.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IConfigEventsMetaData contains all meta data concerning the IConfigEvents contract.
var IConfigEventsMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"oldValue\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newValue\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"appliesAtHeight\",\"type\":\"uint256\"}],\"name\":\"ConfigUpdateAddress\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"oldValue\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"newValue\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"appliesAtHeight\",\"type\":\"uint256\"}],\"name\":\"ConfigUpdateBool\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"int256\",\"name\":\"oldValue\",\"type\":\"int256\"},{\"indexed\":false,\"internalType\":\"int256\",\"name\":\"newValue\",\"type\":\"int256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"appliesAtHeight\",\"type\":\"uint256\"}],\"name\":\"ConfigUpdateInt\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldValue\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newValue\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"appliesAtHeight\",\"type\":\"uint256\"}],\"name\":\"ConfigUpdateUint\",\"type\":\"event\"}]",
}

// IConfigEventsABI is the input ABI used to generate the binding from.
// Deprecated: Use IConfigEventsMetaData.ABI instead.
var IConfigEventsABI = IConfigEventsMetaData.ABI

// IConfigEvents is an auto generated Go binding around an Ethereum contract.
type IConfigEvents struct {
	IConfigEventsCaller     // Read-only binding to the contract
	IConfigEventsTransactor // Write-only binding to the contract
	IConfigEventsFilterer   // Log filterer for contract events
}

// IConfigEventsCaller is an auto generated read-only Go binding around an Ethereum contract.
type IConfigEventsCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IConfigEventsTransactor is an auto generated write-only Go binding around an Ethereum contract.
type IConfigEventsTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IConfigEventsFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type IConfigEventsFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IConfigEventsSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type IConfigEventsSession struct {
	Contract     *IConfigEvents    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IConfigEventsCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type IConfigEventsCallerSession struct {
	Contract *IConfigEventsCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// IConfigEventsTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type IConfigEventsTransactorSession struct {
	Contract     *IConfigEventsTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// IConfigEventsRaw is an auto generated low-level Go binding around an Ethereum contract.
type IConfigEventsRaw struct {
	Contract *IConfigEvents // Generic contract binding to access the raw methods on
}

// IConfigEventsCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type IConfigEventsCallerRaw struct {
	Contract *IConfigEventsCaller // Generic read-only contract binding to access the raw methods on
}

// IConfigEventsTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type IConfigEventsTransactorRaw struct {
	Contract *IConfigEventsTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIConfigEvents creates a new instance of IConfigEvents, bound to a specific deployed contract.
func NewIConfigEvents(address common.Address, backend bind.ContractBackend) (*IConfigEvents, error) {
	contract, err := bindIConfigEvents(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IConfigEvents{IConfigEventsCaller: IConfigEventsCaller{contract: contract}, IConfigEventsTransactor: IConfigEventsTransactor{contract: contract}, IConfigEventsFilterer: IConfigEventsFilterer{contract: contract}}, nil
}

// NewIConfigEventsCaller creates a new read-only instance of IConfigEvents, bound to a specific deployed contract.
func NewIConfigEventsCaller(address common.Address, caller bind.ContractCaller) (*IConfigEventsCaller, error) {
	contract, err := bindIConfigEvents(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IConfigEventsCaller{contract: contract}, nil
}

// NewIConfigEventsTransactor creates a new write-only instance of IConfigEvents, bound to a specific deployed contract.
func NewIConfigEventsTransactor(address common.Address, transactor bind.ContractTransactor) (*IConfigEventsTransactor, error) {
	contract, err := bindIConfigEvents(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IConfigEventsTransactor{contract: contract}, nil
}

// NewIConfigEventsFilterer creates a new log filterer instance of IConfigEvents, bound to a specific deployed contract.
func NewIConfigEventsFilterer(address common.Address, filterer bind.ContractFilterer) (*IConfigEventsFilterer, error) {
	contract, err := bindIConfigEvents(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IConfigEventsFilterer{contract: contract}, nil
}

// bindIConfigEvents binds a generic wrapper to an already deployed contract.
func bindIConfigEvents(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(IConfigEventsABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IConfigEvents *IConfigEventsRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IConfigEvents.Contract.IConfigEventsCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IConfigEvents *IConfigEventsRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IConfigEvents.Contract.IConfigEventsTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IConfigEvents *IConfigEventsRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IConfigEvents.Contract.IConfigEventsTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IConfigEvents *IConfigEventsCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IConfigEvents.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IConfigEvents *IConfigEventsTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IConfigEvents.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IConfigEvents *IConfigEventsTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IConfigEvents.Contract.contract.Transact(opts, method, params...)
}

// IConfigEventsConfigUpdateAddressIterator is returned from FilterConfigUpdateAddress and is used to iterate over the raw logs and unpacked data for ConfigUpdateAddress events raised by the IConfigEvents contract.
type IConfigEventsConfigUpdateAddressIterator struct {
	Event *IConfigEventsConfigUpdateAddress // Event containing the contract specifics and raw log

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
func (it *IConfigEventsConfigUpdateAddressIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IConfigEventsConfigUpdateAddress)
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
		it.Event = new(IConfigEventsConfigUpdateAddress)
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
func (it *IConfigEventsConfigUpdateAddressIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IConfigEventsConfigUpdateAddressIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IConfigEventsConfigUpdateAddress represents a ConfigUpdateAddress event raised by the IConfigEvents contract.
type IConfigEventsConfigUpdateAddress struct {
	Name            string
	OldValue        common.Address
	NewValue        common.Address
	AppliesAtHeight *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterConfigUpdateAddress is a free log retrieval operation binding the contract event 0xe821ac8084a7329d09d00cf1380cba50edea5f54a7bd453d49a9c354f565d4a0.
//
// Solidity: event ConfigUpdateAddress(string name, address oldValue, address newValue, uint256 appliesAtHeight)
func (_IConfigEvents *IConfigEventsFilterer) FilterConfigUpdateAddress(opts *bind.FilterOpts) (*IConfigEventsConfigUpdateAddressIterator, error) {

	logs, sub, err := _IConfigEvents.contract.FilterLogs(opts, "ConfigUpdateAddress")
	if err != nil {
		return nil, err
	}
	return &IConfigEventsConfigUpdateAddressIterator{contract: _IConfigEvents.contract, event: "ConfigUpdateAddress", logs: logs, sub: sub}, nil
}

// WatchConfigUpdateAddress is a free log subscription operation binding the contract event 0xe821ac8084a7329d09d00cf1380cba50edea5f54a7bd453d49a9c354f565d4a0.
//
// Solidity: event ConfigUpdateAddress(string name, address oldValue, address newValue, uint256 appliesAtHeight)
func (_IConfigEvents *IConfigEventsFilterer) WatchConfigUpdateAddress(opts *bind.WatchOpts, sink chan<- *IConfigEventsConfigUpdateAddress) (event.Subscription, error) {

	logs, sub, err := _IConfigEvents.contract.WatchLogs(opts, "ConfigUpdateAddress")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IConfigEventsConfigUpdateAddress)
				if err := _IConfigEvents.contract.UnpackLog(event, "ConfigUpdateAddress", log); err != nil {
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

// ParseConfigUpdateAddress is a log parse operation binding the contract event 0xe821ac8084a7329d09d00cf1380cba50edea5f54a7bd453d49a9c354f565d4a0.
//
// Solidity: event ConfigUpdateAddress(string name, address oldValue, address newValue, uint256 appliesAtHeight)
func (_IConfigEvents *IConfigEventsFilterer) ParseConfigUpdateAddress(log types.Log) (*IConfigEventsConfigUpdateAddress, error) {
	event := new(IConfigEventsConfigUpdateAddress)
	if err := _IConfigEvents.contract.UnpackLog(event, "ConfigUpdateAddress", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IConfigEventsConfigUpdateBoolIterator is returned from FilterConfigUpdateBool and is used to iterate over the raw logs and unpacked data for ConfigUpdateBool events raised by the IConfigEvents contract.
type IConfigEventsConfigUpdateBoolIterator struct {
	Event *IConfigEventsConfigUpdateBool // Event containing the contract specifics and raw log

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
func (it *IConfigEventsConfigUpdateBoolIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IConfigEventsConfigUpdateBool)
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
		it.Event = new(IConfigEventsConfigUpdateBool)
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
func (it *IConfigEventsConfigUpdateBoolIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IConfigEventsConfigUpdateBoolIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IConfigEventsConfigUpdateBool represents a ConfigUpdateBool event raised by the IConfigEvents contract.
type IConfigEventsConfigUpdateBool struct {
	Name            string
	OldValue        bool
	NewValue        bool
	AppliesAtHeight *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterConfigUpdateBool is a free log retrieval operation binding the contract event 0x5edb308c5eddc69bcd31b4e689c5eed2fbd3155ae57915d1cad05425f6c1a39b.
//
// Solidity: event ConfigUpdateBool(string name, bool oldValue, bool newValue, uint256 appliesAtHeight)
func (_IConfigEvents *IConfigEventsFilterer) FilterConfigUpdateBool(opts *bind.FilterOpts) (*IConfigEventsConfigUpdateBoolIterator, error) {

	logs, sub, err := _IConfigEvents.contract.FilterLogs(opts, "ConfigUpdateBool")
	if err != nil {
		return nil, err
	}
	return &IConfigEventsConfigUpdateBoolIterator{contract: _IConfigEvents.contract, event: "ConfigUpdateBool", logs: logs, sub: sub}, nil
}

// WatchConfigUpdateBool is a free log subscription operation binding the contract event 0x5edb308c5eddc69bcd31b4e689c5eed2fbd3155ae57915d1cad05425f6c1a39b.
//
// Solidity: event ConfigUpdateBool(string name, bool oldValue, bool newValue, uint256 appliesAtHeight)
func (_IConfigEvents *IConfigEventsFilterer) WatchConfigUpdateBool(opts *bind.WatchOpts, sink chan<- *IConfigEventsConfigUpdateBool) (event.Subscription, error) {

	logs, sub, err := _IConfigEvents.contract.WatchLogs(opts, "ConfigUpdateBool")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IConfigEventsConfigUpdateBool)
				if err := _IConfigEvents.contract.UnpackLog(event, "ConfigUpdateBool", log); err != nil {
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

// ParseConfigUpdateBool is a log parse operation binding the contract event 0x5edb308c5eddc69bcd31b4e689c5eed2fbd3155ae57915d1cad05425f6c1a39b.
//
// Solidity: event ConfigUpdateBool(string name, bool oldValue, bool newValue, uint256 appliesAtHeight)
func (_IConfigEvents *IConfigEventsFilterer) ParseConfigUpdateBool(log types.Log) (*IConfigEventsConfigUpdateBool, error) {
	event := new(IConfigEventsConfigUpdateBool)
	if err := _IConfigEvents.contract.UnpackLog(event, "ConfigUpdateBool", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IConfigEventsConfigUpdateIntIterator is returned from FilterConfigUpdateInt and is used to iterate over the raw logs and unpacked data for ConfigUpdateInt events raised by the IConfigEvents contract.
type IConfigEventsConfigUpdateIntIterator struct {
	Event *IConfigEventsConfigUpdateInt // Event containing the contract specifics and raw log

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
func (it *IConfigEventsConfigUpdateIntIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IConfigEventsConfigUpdateInt)
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
		it.Event = new(IConfigEventsConfigUpdateInt)
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
func (it *IConfigEventsConfigUpdateIntIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IConfigEventsConfigUpdateIntIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IConfigEventsConfigUpdateInt represents a ConfigUpdateInt event raised by the IConfigEvents contract.
type IConfigEventsConfigUpdateInt struct {
	Name            string
	OldValue        *big.Int
	NewValue        *big.Int
	AppliesAtHeight *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterConfigUpdateInt is a free log retrieval operation binding the contract event 0xb5114472b89d1126433287defb6308bbabdb95b6ce5dd949ab2c151ed7b1ff4c.
//
// Solidity: event ConfigUpdateInt(string name, int256 oldValue, int256 newValue, uint256 appliesAtHeight)
func (_IConfigEvents *IConfigEventsFilterer) FilterConfigUpdateInt(opts *bind.FilterOpts) (*IConfigEventsConfigUpdateIntIterator, error) {

	logs, sub, err := _IConfigEvents.contract.FilterLogs(opts, "ConfigUpdateInt")
	if err != nil {
		return nil, err
	}
	return &IConfigEventsConfigUpdateIntIterator{contract: _IConfigEvents.contract, event: "ConfigUpdateInt", logs: logs, sub: sub}, nil
}

// WatchConfigUpdateInt is a free log subscription operation binding the contract event 0xb5114472b89d1126433287defb6308bbabdb95b6ce5dd949ab2c151ed7b1ff4c.
//
// Solidity: event ConfigUpdateInt(string name, int256 oldValue, int256 newValue, uint256 appliesAtHeight)
func (_IConfigEvents *IConfigEventsFilterer) WatchConfigUpdateInt(opts *bind.WatchOpts, sink chan<- *IConfigEventsConfigUpdateInt) (event.Subscription, error) {

	logs, sub, err := _IConfigEvents.contract.WatchLogs(opts, "ConfigUpdateInt")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IConfigEventsConfigUpdateInt)
				if err := _IConfigEvents.contract.UnpackLog(event, "ConfigUpdateInt", log); err != nil {
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

// ParseConfigUpdateInt is a log parse operation binding the contract event 0xb5114472b89d1126433287defb6308bbabdb95b6ce5dd949ab2c151ed7b1ff4c.
//
// Solidity: event ConfigUpdateInt(string name, int256 oldValue, int256 newValue, uint256 appliesAtHeight)
func (_IConfigEvents *IConfigEventsFilterer) ParseConfigUpdateInt(log types.Log) (*IConfigEventsConfigUpdateInt, error) {
	event := new(IConfigEventsConfigUpdateInt)
	if err := _IConfigEvents.contract.UnpackLog(event, "ConfigUpdateInt", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IConfigEventsConfigUpdateUintIterator is returned from FilterConfigUpdateUint and is used to iterate over the raw logs and unpacked data for ConfigUpdateUint events raised by the IConfigEvents contract.
type IConfigEventsConfigUpdateUintIterator struct {
	Event *IConfigEventsConfigUpdateUint // Event containing the contract specifics and raw log

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
func (it *IConfigEventsConfigUpdateUintIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IConfigEventsConfigUpdateUint)
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
		it.Event = new(IConfigEventsConfigUpdateUint)
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
func (it *IConfigEventsConfigUpdateUintIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IConfigEventsConfigUpdateUintIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IConfigEventsConfigUpdateUint represents a ConfigUpdateUint event raised by the IConfigEvents contract.
type IConfigEventsConfigUpdateUint struct {
	Name            string
	OldValue        *big.Int
	NewValue        *big.Int
	AppliesAtHeight *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterConfigUpdateUint is a free log retrieval operation binding the contract event 0x207e45ce6f2191c3efffe5d0d91591cf81afd23cbf4bd746981d7a8a8dcfe1ba.
//
// Solidity: event ConfigUpdateUint(string name, uint256 oldValue, uint256 newValue, uint256 appliesAtHeight)
func (_IConfigEvents *IConfigEventsFilterer) FilterConfigUpdateUint(opts *bind.FilterOpts) (*IConfigEventsConfigUpdateUintIterator, error) {

	logs, sub, err := _IConfigEvents.contract.FilterLogs(opts, "ConfigUpdateUint")
	if err != nil {
		return nil, err
	}
	return &IConfigEventsConfigUpdateUintIterator{contract: _IConfigEvents.contract, event: "ConfigUpdateUint", logs: logs, sub: sub}, nil
}

// WatchConfigUpdateUint is a free log subscription operation binding the contract event 0x207e45ce6f2191c3efffe5d0d91591cf81afd23cbf4bd746981d7a8a8dcfe1ba.
//
// Solidity: event ConfigUpdateUint(string name, uint256 oldValue, uint256 newValue, uint256 appliesAtHeight)
func (_IConfigEvents *IConfigEventsFilterer) WatchConfigUpdateUint(opts *bind.WatchOpts, sink chan<- *IConfigEventsConfigUpdateUint) (event.Subscription, error) {

	logs, sub, err := _IConfigEvents.contract.WatchLogs(opts, "ConfigUpdateUint")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IConfigEventsConfigUpdateUint)
				if err := _IConfigEvents.contract.UnpackLog(event, "ConfigUpdateUint", log); err != nil {
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

// ParseConfigUpdateUint is a log parse operation binding the contract event 0x207e45ce6f2191c3efffe5d0d91591cf81afd23cbf4bd746981d7a8a8dcfe1ba.
//
// Solidity: event ConfigUpdateUint(string name, uint256 oldValue, uint256 newValue, uint256 appliesAtHeight)
func (_IConfigEvents *IConfigEventsFilterer) ParseConfigUpdateUint(log types.Log) (*IConfigEventsConfigUpdateUint, error) {
	event := new(IConfigEventsConfigUpdateUint)
	if err := _IConfigEvents.contract.UnpackLog(event, "ConfigUpdateUint", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IERC20MetaData contains all meta data concerning the IERC20 contract.
var IERC20MetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"dd62ed3e": "allowance(address,address)",
		"095ea7b3": "approve(address,uint256)",
		"70a08231": "balanceOf(address)",
		"18160ddd": "totalSupply()",
		"a9059cbb": "transfer(address,uint256)",
		"23b872dd": "transferFrom(address,address,uint256)",
	},
}

// IERC20ABI is the input ABI used to generate the binding from.
// Deprecated: Use IERC20MetaData.ABI instead.
var IERC20ABI = IERC20MetaData.ABI

// Deprecated: Use IERC20MetaData.Sigs instead.
// IERC20FuncSigs maps the 4-byte function signature to its string representation.
var IERC20FuncSigs = IERC20MetaData.Sigs

// IERC20 is an auto generated Go binding around an Ethereum contract.
type IERC20 struct {
	IERC20Caller     // Read-only binding to the contract
	IERC20Transactor // Write-only binding to the contract
	IERC20Filterer   // Log filterer for contract events
}

// IERC20Caller is an auto generated read-only Go binding around an Ethereum contract.
type IERC20Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IERC20Transactor is an auto generated write-only Go binding around an Ethereum contract.
type IERC20Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IERC20Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type IERC20Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IERC20Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type IERC20Session struct {
	Contract     *IERC20           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IERC20CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type IERC20CallerSession struct {
	Contract *IERC20Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// IERC20TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type IERC20TransactorSession struct {
	Contract     *IERC20Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IERC20Raw is an auto generated low-level Go binding around an Ethereum contract.
type IERC20Raw struct {
	Contract *IERC20 // Generic contract binding to access the raw methods on
}

// IERC20CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type IERC20CallerRaw struct {
	Contract *IERC20Caller // Generic read-only contract binding to access the raw methods on
}

// IERC20TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type IERC20TransactorRaw struct {
	Contract *IERC20Transactor // Generic write-only contract binding to access the raw methods on
}

// NewIERC20 creates a new instance of IERC20, bound to a specific deployed contract.
func NewIERC20(address common.Address, backend bind.ContractBackend) (*IERC20, error) {
	contract, err := bindIERC20(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IERC20{IERC20Caller: IERC20Caller{contract: contract}, IERC20Transactor: IERC20Transactor{contract: contract}, IERC20Filterer: IERC20Filterer{contract: contract}}, nil
}

// NewIERC20Caller creates a new read-only instance of IERC20, bound to a specific deployed contract.
func NewIERC20Caller(address common.Address, caller bind.ContractCaller) (*IERC20Caller, error) {
	contract, err := bindIERC20(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IERC20Caller{contract: contract}, nil
}

// NewIERC20Transactor creates a new write-only instance of IERC20, bound to a specific deployed contract.
func NewIERC20Transactor(address common.Address, transactor bind.ContractTransactor) (*IERC20Transactor, error) {
	contract, err := bindIERC20(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IERC20Transactor{contract: contract}, nil
}

// NewIERC20Filterer creates a new log filterer instance of IERC20, bound to a specific deployed contract.
func NewIERC20Filterer(address common.Address, filterer bind.ContractFilterer) (*IERC20Filterer, error) {
	contract, err := bindIERC20(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IERC20Filterer{contract: contract}, nil
}

// bindIERC20 binds a generic wrapper to an already deployed contract.
func bindIERC20(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(IERC20ABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IERC20 *IERC20Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IERC20.Contract.IERC20Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IERC20 *IERC20Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IERC20.Contract.IERC20Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IERC20 *IERC20Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IERC20.Contract.IERC20Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IERC20 *IERC20CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IERC20.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IERC20 *IERC20TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IERC20.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IERC20 *IERC20TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IERC20.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_IERC20 *IERC20Caller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _IERC20.contract.Call(opts, &out, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_IERC20 *IERC20Session) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _IERC20.Contract.Allowance(&_IERC20.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_IERC20 *IERC20CallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _IERC20.Contract.Allowance(&_IERC20.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_IERC20 *IERC20Caller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _IERC20.contract.Call(opts, &out, "balanceOf", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_IERC20 *IERC20Session) BalanceOf(account common.Address) (*big.Int, error) {
	return _IERC20.Contract.BalanceOf(&_IERC20.CallOpts, account)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_IERC20 *IERC20CallerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _IERC20.Contract.BalanceOf(&_IERC20.CallOpts, account)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_IERC20 *IERC20Caller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IERC20.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_IERC20 *IERC20Session) TotalSupply() (*big.Int, error) {
	return _IERC20.Contract.TotalSupply(&_IERC20.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_IERC20 *IERC20CallerSession) TotalSupply() (*big.Int, error) {
	return _IERC20.Contract.TotalSupply(&_IERC20.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_IERC20 *IERC20Transactor) Approve(opts *bind.TransactOpts, spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _IERC20.contract.Transact(opts, "approve", spender, amount)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_IERC20 *IERC20Session) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _IERC20.Contract.Approve(&_IERC20.TransactOpts, spender, amount)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_IERC20 *IERC20TransactorSession) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _IERC20.Contract.Approve(&_IERC20.TransactOpts, spender, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_IERC20 *IERC20Transactor) Transfer(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _IERC20.contract.Transact(opts, "transfer", recipient, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_IERC20 *IERC20Session) Transfer(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _IERC20.Contract.Transfer(&_IERC20.TransactOpts, recipient, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_IERC20 *IERC20TransactorSession) Transfer(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _IERC20.Contract.Transfer(&_IERC20.TransactOpts, recipient, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_IERC20 *IERC20Transactor) TransferFrom(opts *bind.TransactOpts, sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _IERC20.contract.Transact(opts, "transferFrom", sender, recipient, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_IERC20 *IERC20Session) TransferFrom(sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _IERC20.Contract.TransferFrom(&_IERC20.TransactOpts, sender, recipient, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_IERC20 *IERC20TransactorSession) TransferFrom(sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _IERC20.Contract.TransferFrom(&_IERC20.TransactOpts, sender, recipient, amount)
}

// IERC20ApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the IERC20 contract.
type IERC20ApprovalIterator struct {
	Event *IERC20Approval // Event containing the contract specifics and raw log

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
func (it *IERC20ApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IERC20Approval)
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
		it.Event = new(IERC20Approval)
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
func (it *IERC20ApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IERC20ApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IERC20Approval represents a Approval event raised by the IERC20 contract.
type IERC20Approval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_IERC20 *IERC20Filterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*IERC20ApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _IERC20.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &IERC20ApprovalIterator{contract: _IERC20.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_IERC20 *IERC20Filterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *IERC20Approval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _IERC20.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IERC20Approval)
				if err := _IERC20.contract.UnpackLog(event, "Approval", log); err != nil {
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

// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_IERC20 *IERC20Filterer) ParseApproval(log types.Log) (*IERC20Approval, error) {
	event := new(IERC20Approval)
	if err := _IERC20.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IERC20TransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the IERC20 contract.
type IERC20TransferIterator struct {
	Event *IERC20Transfer // Event containing the contract specifics and raw log

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
func (it *IERC20TransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IERC20Transfer)
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
		it.Event = new(IERC20Transfer)
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
func (it *IERC20TransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IERC20TransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IERC20Transfer represents a Transfer event raised by the IERC20 contract.
type IERC20Transfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_IERC20 *IERC20Filterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*IERC20TransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _IERC20.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &IERC20TransferIterator{contract: _IERC20.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_IERC20 *IERC20Filterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *IERC20Transfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _IERC20.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IERC20Transfer)
				if err := _IERC20.contract.UnpackLog(event, "Transfer", log); err != nil {
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
func (_IERC20 *IERC20Filterer) ParseTransfer(log types.Log) (*IERC20Transfer, error) {
	event := new(IERC20Transfer)
	if err := _IERC20.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IInflationControllerMetaData contains all meta data concerning the IInflationController contract.
var IInflationControllerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_currentSupply\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_inflationReserve\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_lastEpochTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_currentEpochTime\",\"type\":\"uint256\"}],\"name\":\"calculateSupplyDelta\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"92eff3cd": "calculateSupplyDelta(uint256,uint256,uint256,uint256)",
	},
}

// IInflationControllerABI is the input ABI used to generate the binding from.
// Deprecated: Use IInflationControllerMetaData.ABI instead.
var IInflationControllerABI = IInflationControllerMetaData.ABI

// Deprecated: Use IInflationControllerMetaData.Sigs instead.
// IInflationControllerFuncSigs maps the 4-byte function signature to its string representation.
var IInflationControllerFuncSigs = IInflationControllerMetaData.Sigs

// IInflationController is an auto generated Go binding around an Ethereum contract.
type IInflationController struct {
	IInflationControllerCaller     // Read-only binding to the contract
	IInflationControllerTransactor // Write-only binding to the contract
	IInflationControllerFilterer   // Log filterer for contract events
}

// IInflationControllerCaller is an auto generated read-only Go binding around an Ethereum contract.
type IInflationControllerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IInflationControllerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type IInflationControllerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IInflationControllerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type IInflationControllerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IInflationControllerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type IInflationControllerSession struct {
	Contract     *IInflationController // Generic contract binding to set the session for
	CallOpts     bind.CallOpts         // Call options to use throughout this session
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// IInflationControllerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type IInflationControllerCallerSession struct {
	Contract *IInflationControllerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts               // Call options to use throughout this session
}

// IInflationControllerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type IInflationControllerTransactorSession struct {
	Contract     *IInflationControllerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts               // Transaction auth options to use throughout this session
}

// IInflationControllerRaw is an auto generated low-level Go binding around an Ethereum contract.
type IInflationControllerRaw struct {
	Contract *IInflationController // Generic contract binding to access the raw methods on
}

// IInflationControllerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type IInflationControllerCallerRaw struct {
	Contract *IInflationControllerCaller // Generic read-only contract binding to access the raw methods on
}

// IInflationControllerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type IInflationControllerTransactorRaw struct {
	Contract *IInflationControllerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIInflationController creates a new instance of IInflationController, bound to a specific deployed contract.
func NewIInflationController(address common.Address, backend bind.ContractBackend) (*IInflationController, error) {
	contract, err := bindIInflationController(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IInflationController{IInflationControllerCaller: IInflationControllerCaller{contract: contract}, IInflationControllerTransactor: IInflationControllerTransactor{contract: contract}, IInflationControllerFilterer: IInflationControllerFilterer{contract: contract}}, nil
}

// NewIInflationControllerCaller creates a new read-only instance of IInflationController, bound to a specific deployed contract.
func NewIInflationControllerCaller(address common.Address, caller bind.ContractCaller) (*IInflationControllerCaller, error) {
	contract, err := bindIInflationController(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IInflationControllerCaller{contract: contract}, nil
}

// NewIInflationControllerTransactor creates a new write-only instance of IInflationController, bound to a specific deployed contract.
func NewIInflationControllerTransactor(address common.Address, transactor bind.ContractTransactor) (*IInflationControllerTransactor, error) {
	contract, err := bindIInflationController(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IInflationControllerTransactor{contract: contract}, nil
}

// NewIInflationControllerFilterer creates a new log filterer instance of IInflationController, bound to a specific deployed contract.
func NewIInflationControllerFilterer(address common.Address, filterer bind.ContractFilterer) (*IInflationControllerFilterer, error) {
	contract, err := bindIInflationController(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IInflationControllerFilterer{contract: contract}, nil
}

// bindIInflationController binds a generic wrapper to an already deployed contract.
func bindIInflationController(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(IInflationControllerABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IInflationController *IInflationControllerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IInflationController.Contract.IInflationControllerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IInflationController *IInflationControllerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IInflationController.Contract.IInflationControllerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IInflationController *IInflationControllerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IInflationController.Contract.IInflationControllerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IInflationController *IInflationControllerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IInflationController.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IInflationController *IInflationControllerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IInflationController.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IInflationController *IInflationControllerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IInflationController.Contract.contract.Transact(opts, method, params...)
}

// CalculateSupplyDelta is a free data retrieval call binding the contract method 0x92eff3cd.
//
// Solidity: function calculateSupplyDelta(uint256 _currentSupply, uint256 _inflationReserve, uint256 _lastEpochTime, uint256 _currentEpochTime) view returns(uint256)
func (_IInflationController *IInflationControllerCaller) CalculateSupplyDelta(opts *bind.CallOpts, _currentSupply *big.Int, _inflationReserve *big.Int, _lastEpochTime *big.Int, _currentEpochTime *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _IInflationController.contract.Call(opts, &out, "calculateSupplyDelta", _currentSupply, _inflationReserve, _lastEpochTime, _currentEpochTime)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CalculateSupplyDelta is a free data retrieval call binding the contract method 0x92eff3cd.
//
// Solidity: function calculateSupplyDelta(uint256 _currentSupply, uint256 _inflationReserve, uint256 _lastEpochTime, uint256 _currentEpochTime) view returns(uint256)
func (_IInflationController *IInflationControllerSession) CalculateSupplyDelta(_currentSupply *big.Int, _inflationReserve *big.Int, _lastEpochTime *big.Int, _currentEpochTime *big.Int) (*big.Int, error) {
	return _IInflationController.Contract.CalculateSupplyDelta(&_IInflationController.CallOpts, _currentSupply, _inflationReserve, _lastEpochTime, _currentEpochTime)
}

// CalculateSupplyDelta is a free data retrieval call binding the contract method 0x92eff3cd.
//
// Solidity: function calculateSupplyDelta(uint256 _currentSupply, uint256 _inflationReserve, uint256 _lastEpochTime, uint256 _currentEpochTime) view returns(uint256)
func (_IInflationController *IInflationControllerCallerSession) CalculateSupplyDelta(_currentSupply *big.Int, _inflationReserve *big.Int, _lastEpochTime *big.Int, _currentEpochTime *big.Int) (*big.Int, error) {
	return _IInflationController.Contract.CalculateSupplyDelta(&_IInflationController.CallOpts, _currentSupply, _inflationReserve, _lastEpochTime, _currentEpochTime)
}

// ILiquidMetaData contains all meta data concerning the ILiquid contract.
var ILiquidMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_caller\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"UnbondingApproval\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_caller\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"approveUnbonding\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"claimRewards\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"claimTreasuryATN\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCommissionRate\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTreasury\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTreasuryUnclaimedATN\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getValidator\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"lock\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_account\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_caller\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"lockFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_delegator\",\"type\":\"address\"}],\"name\":\"lockedBalanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"mint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_ntnReward\",\"type\":\"uint256\"}],\"name\":\"redistribute\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_rate\",\"type\":\"uint256\"}],\"name\":\"setCommissionRate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_caller\",\"type\":\"address\"}],\"name\":\"unbondingAllowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_account\",\"type\":\"address\"}],\"name\":\"unclaimedRewards\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"unlock\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_delegator\",\"type\":\"address\"}],\"name\":\"unlockedBalanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"dd62ed3e": "allowance(address,address)",
		"095ea7b3": "approve(address,uint256)",
		"bf99b73e": "approveUnbonding(address,uint256)",
		"70a08231": "balanceOf(address)",
		"9dc29fac": "burn(address,uint256)",
		"372500ab": "claimRewards()",
		"bd96102f": "claimTreasuryATN()",
		"313ce567": "decimals()",
		"3e4eb36c": "getCommissionRate()",
		"3b19e84a": "getTreasury()",
		"1eeffad0": "getTreasuryUnclaimedATN()",
		"1195e07e": "getValidator()",
		"282d3fdf": "lock(address,uint256)",
		"708e91e5": "lockFrom(address,address,uint256)",
		"59355736": "lockedBalanceOf(address)",
		"40c10f19": "mint(address,uint256)",
		"06fdde03": "name()",
		"a0ce552d": "redistribute(uint256)",
		"19fac8fd": "setCommissionRate(uint256)",
		"95d89b41": "symbol()",
		"18160ddd": "totalSupply()",
		"a9059cbb": "transfer(address,uint256)",
		"23b872dd": "transferFrom(address,address,uint256)",
		"d768d578": "unbondingAllowance(address,address)",
		"949813b8": "unclaimedRewards(address)",
		"7eee288d": "unlock(address,uint256)",
		"84955c88": "unlockedBalanceOf(address)",
	},
}

// ILiquidABI is the input ABI used to generate the binding from.
// Deprecated: Use ILiquidMetaData.ABI instead.
var ILiquidABI = ILiquidMetaData.ABI

// Deprecated: Use ILiquidMetaData.Sigs instead.
// ILiquidFuncSigs maps the 4-byte function signature to its string representation.
var ILiquidFuncSigs = ILiquidMetaData.Sigs

// ILiquid is an auto generated Go binding around an Ethereum contract.
type ILiquid struct {
	ILiquidCaller     // Read-only binding to the contract
	ILiquidTransactor // Write-only binding to the contract
	ILiquidFilterer   // Log filterer for contract events
}

// ILiquidCaller is an auto generated read-only Go binding around an Ethereum contract.
type ILiquidCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ILiquidTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ILiquidTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ILiquidFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ILiquidFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ILiquidSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ILiquidSession struct {
	Contract     *ILiquid          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ILiquidCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ILiquidCallerSession struct {
	Contract *ILiquidCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// ILiquidTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ILiquidTransactorSession struct {
	Contract     *ILiquidTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// ILiquidRaw is an auto generated low-level Go binding around an Ethereum contract.
type ILiquidRaw struct {
	Contract *ILiquid // Generic contract binding to access the raw methods on
}

// ILiquidCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ILiquidCallerRaw struct {
	Contract *ILiquidCaller // Generic read-only contract binding to access the raw methods on
}

// ILiquidTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ILiquidTransactorRaw struct {
	Contract *ILiquidTransactor // Generic write-only contract binding to access the raw methods on
}

// NewILiquid creates a new instance of ILiquid, bound to a specific deployed contract.
func NewILiquid(address common.Address, backend bind.ContractBackend) (*ILiquid, error) {
	contract, err := bindILiquid(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ILiquid{ILiquidCaller: ILiquidCaller{contract: contract}, ILiquidTransactor: ILiquidTransactor{contract: contract}, ILiquidFilterer: ILiquidFilterer{contract: contract}}, nil
}

// NewILiquidCaller creates a new read-only instance of ILiquid, bound to a specific deployed contract.
func NewILiquidCaller(address common.Address, caller bind.ContractCaller) (*ILiquidCaller, error) {
	contract, err := bindILiquid(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ILiquidCaller{contract: contract}, nil
}

// NewILiquidTransactor creates a new write-only instance of ILiquid, bound to a specific deployed contract.
func NewILiquidTransactor(address common.Address, transactor bind.ContractTransactor) (*ILiquidTransactor, error) {
	contract, err := bindILiquid(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ILiquidTransactor{contract: contract}, nil
}

// NewILiquidFilterer creates a new log filterer instance of ILiquid, bound to a specific deployed contract.
func NewILiquidFilterer(address common.Address, filterer bind.ContractFilterer) (*ILiquidFilterer, error) {
	contract, err := bindILiquid(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ILiquidFilterer{contract: contract}, nil
}

// bindILiquid binds a generic wrapper to an already deployed contract.
func bindILiquid(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ILiquidABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ILiquid *ILiquidRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ILiquid.Contract.ILiquidCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ILiquid *ILiquidRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ILiquid.Contract.ILiquidTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ILiquid *ILiquidRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ILiquid.Contract.ILiquidTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ILiquid *ILiquidCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ILiquid.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ILiquid *ILiquidTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ILiquid.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ILiquid *ILiquidTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ILiquid.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_ILiquid *ILiquidCaller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ILiquid.contract.Call(opts, &out, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_ILiquid *ILiquidSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _ILiquid.Contract.Allowance(&_ILiquid.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_ILiquid *ILiquidCallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _ILiquid.Contract.Allowance(&_ILiquid.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_ILiquid *ILiquidCaller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ILiquid.contract.Call(opts, &out, "balanceOf", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_ILiquid *ILiquidSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _ILiquid.Contract.BalanceOf(&_ILiquid.CallOpts, account)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_ILiquid *ILiquidCallerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _ILiquid.Contract.BalanceOf(&_ILiquid.CallOpts, account)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() pure returns(uint8)
func (_ILiquid *ILiquidCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _ILiquid.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() pure returns(uint8)
func (_ILiquid *ILiquidSession) Decimals() (uint8, error) {
	return _ILiquid.Contract.Decimals(&_ILiquid.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() pure returns(uint8)
func (_ILiquid *ILiquidCallerSession) Decimals() (uint8, error) {
	return _ILiquid.Contract.Decimals(&_ILiquid.CallOpts)
}

// GetCommissionRate is a free data retrieval call binding the contract method 0x3e4eb36c.
//
// Solidity: function getCommissionRate() view returns(uint256)
func (_ILiquid *ILiquidCaller) GetCommissionRate(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ILiquid.contract.Call(opts, &out, "getCommissionRate")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetCommissionRate is a free data retrieval call binding the contract method 0x3e4eb36c.
//
// Solidity: function getCommissionRate() view returns(uint256)
func (_ILiquid *ILiquidSession) GetCommissionRate() (*big.Int, error) {
	return _ILiquid.Contract.GetCommissionRate(&_ILiquid.CallOpts)
}

// GetCommissionRate is a free data retrieval call binding the contract method 0x3e4eb36c.
//
// Solidity: function getCommissionRate() view returns(uint256)
func (_ILiquid *ILiquidCallerSession) GetCommissionRate() (*big.Int, error) {
	return _ILiquid.Contract.GetCommissionRate(&_ILiquid.CallOpts)
}

// GetTreasury is a free data retrieval call binding the contract method 0x3b19e84a.
//
// Solidity: function getTreasury() view returns(address)
func (_ILiquid *ILiquidCaller) GetTreasury(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ILiquid.contract.Call(opts, &out, "getTreasury")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetTreasury is a free data retrieval call binding the contract method 0x3b19e84a.
//
// Solidity: function getTreasury() view returns(address)
func (_ILiquid *ILiquidSession) GetTreasury() (common.Address, error) {
	return _ILiquid.Contract.GetTreasury(&_ILiquid.CallOpts)
}

// GetTreasury is a free data retrieval call binding the contract method 0x3b19e84a.
//
// Solidity: function getTreasury() view returns(address)
func (_ILiquid *ILiquidCallerSession) GetTreasury() (common.Address, error) {
	return _ILiquid.Contract.GetTreasury(&_ILiquid.CallOpts)
}

// GetTreasuryUnclaimedATN is a free data retrieval call binding the contract method 0x1eeffad0.
//
// Solidity: function getTreasuryUnclaimedATN() view returns(uint256)
func (_ILiquid *ILiquidCaller) GetTreasuryUnclaimedATN(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ILiquid.contract.Call(opts, &out, "getTreasuryUnclaimedATN")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetTreasuryUnclaimedATN is a free data retrieval call binding the contract method 0x1eeffad0.
//
// Solidity: function getTreasuryUnclaimedATN() view returns(uint256)
func (_ILiquid *ILiquidSession) GetTreasuryUnclaimedATN() (*big.Int, error) {
	return _ILiquid.Contract.GetTreasuryUnclaimedATN(&_ILiquid.CallOpts)
}

// GetTreasuryUnclaimedATN is a free data retrieval call binding the contract method 0x1eeffad0.
//
// Solidity: function getTreasuryUnclaimedATN() view returns(uint256)
func (_ILiquid *ILiquidCallerSession) GetTreasuryUnclaimedATN() (*big.Int, error) {
	return _ILiquid.Contract.GetTreasuryUnclaimedATN(&_ILiquid.CallOpts)
}

// GetValidator is a free data retrieval call binding the contract method 0x1195e07e.
//
// Solidity: function getValidator() view returns(address)
func (_ILiquid *ILiquidCaller) GetValidator(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ILiquid.contract.Call(opts, &out, "getValidator")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetValidator is a free data retrieval call binding the contract method 0x1195e07e.
//
// Solidity: function getValidator() view returns(address)
func (_ILiquid *ILiquidSession) GetValidator() (common.Address, error) {
	return _ILiquid.Contract.GetValidator(&_ILiquid.CallOpts)
}

// GetValidator is a free data retrieval call binding the contract method 0x1195e07e.
//
// Solidity: function getValidator() view returns(address)
func (_ILiquid *ILiquidCallerSession) GetValidator() (common.Address, error) {
	return _ILiquid.Contract.GetValidator(&_ILiquid.CallOpts)
}

// LockedBalanceOf is a free data retrieval call binding the contract method 0x59355736.
//
// Solidity: function lockedBalanceOf(address _delegator) view returns(uint256)
func (_ILiquid *ILiquidCaller) LockedBalanceOf(opts *bind.CallOpts, _delegator common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ILiquid.contract.Call(opts, &out, "lockedBalanceOf", _delegator)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LockedBalanceOf is a free data retrieval call binding the contract method 0x59355736.
//
// Solidity: function lockedBalanceOf(address _delegator) view returns(uint256)
func (_ILiquid *ILiquidSession) LockedBalanceOf(_delegator common.Address) (*big.Int, error) {
	return _ILiquid.Contract.LockedBalanceOf(&_ILiquid.CallOpts, _delegator)
}

// LockedBalanceOf is a free data retrieval call binding the contract method 0x59355736.
//
// Solidity: function lockedBalanceOf(address _delegator) view returns(uint256)
func (_ILiquid *ILiquidCallerSession) LockedBalanceOf(_delegator common.Address) (*big.Int, error) {
	return _ILiquid.Contract.LockedBalanceOf(&_ILiquid.CallOpts, _delegator)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_ILiquid *ILiquidCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _ILiquid.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_ILiquid *ILiquidSession) Name() (string, error) {
	return _ILiquid.Contract.Name(&_ILiquid.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_ILiquid *ILiquidCallerSession) Name() (string, error) {
	return _ILiquid.Contract.Name(&_ILiquid.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_ILiquid *ILiquidCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _ILiquid.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_ILiquid *ILiquidSession) Symbol() (string, error) {
	return _ILiquid.Contract.Symbol(&_ILiquid.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_ILiquid *ILiquidCallerSession) Symbol() (string, error) {
	return _ILiquid.Contract.Symbol(&_ILiquid.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_ILiquid *ILiquidCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ILiquid.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_ILiquid *ILiquidSession) TotalSupply() (*big.Int, error) {
	return _ILiquid.Contract.TotalSupply(&_ILiquid.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_ILiquid *ILiquidCallerSession) TotalSupply() (*big.Int, error) {
	return _ILiquid.Contract.TotalSupply(&_ILiquid.CallOpts)
}

// UnbondingAllowance is a free data retrieval call binding the contract method 0xd768d578.
//
// Solidity: function unbondingAllowance(address _owner, address _caller) view returns(uint256)
func (_ILiquid *ILiquidCaller) UnbondingAllowance(opts *bind.CallOpts, _owner common.Address, _caller common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ILiquid.contract.Call(opts, &out, "unbondingAllowance", _owner, _caller)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// UnbondingAllowance is a free data retrieval call binding the contract method 0xd768d578.
//
// Solidity: function unbondingAllowance(address _owner, address _caller) view returns(uint256)
func (_ILiquid *ILiquidSession) UnbondingAllowance(_owner common.Address, _caller common.Address) (*big.Int, error) {
	return _ILiquid.Contract.UnbondingAllowance(&_ILiquid.CallOpts, _owner, _caller)
}

// UnbondingAllowance is a free data retrieval call binding the contract method 0xd768d578.
//
// Solidity: function unbondingAllowance(address _owner, address _caller) view returns(uint256)
func (_ILiquid *ILiquidCallerSession) UnbondingAllowance(_owner common.Address, _caller common.Address) (*big.Int, error) {
	return _ILiquid.Contract.UnbondingAllowance(&_ILiquid.CallOpts, _owner, _caller)
}

// UnclaimedRewards is a free data retrieval call binding the contract method 0x949813b8.
//
// Solidity: function unclaimedRewards(address _account) view returns(uint256)
func (_ILiquid *ILiquidCaller) UnclaimedRewards(opts *bind.CallOpts, _account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ILiquid.contract.Call(opts, &out, "unclaimedRewards", _account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// UnclaimedRewards is a free data retrieval call binding the contract method 0x949813b8.
//
// Solidity: function unclaimedRewards(address _account) view returns(uint256)
func (_ILiquid *ILiquidSession) UnclaimedRewards(_account common.Address) (*big.Int, error) {
	return _ILiquid.Contract.UnclaimedRewards(&_ILiquid.CallOpts, _account)
}

// UnclaimedRewards is a free data retrieval call binding the contract method 0x949813b8.
//
// Solidity: function unclaimedRewards(address _account) view returns(uint256)
func (_ILiquid *ILiquidCallerSession) UnclaimedRewards(_account common.Address) (*big.Int, error) {
	return _ILiquid.Contract.UnclaimedRewards(&_ILiquid.CallOpts, _account)
}

// UnlockedBalanceOf is a free data retrieval call binding the contract method 0x84955c88.
//
// Solidity: function unlockedBalanceOf(address _delegator) view returns(uint256)
func (_ILiquid *ILiquidCaller) UnlockedBalanceOf(opts *bind.CallOpts, _delegator common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ILiquid.contract.Call(opts, &out, "unlockedBalanceOf", _delegator)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// UnlockedBalanceOf is a free data retrieval call binding the contract method 0x84955c88.
//
// Solidity: function unlockedBalanceOf(address _delegator) view returns(uint256)
func (_ILiquid *ILiquidSession) UnlockedBalanceOf(_delegator common.Address) (*big.Int, error) {
	return _ILiquid.Contract.UnlockedBalanceOf(&_ILiquid.CallOpts, _delegator)
}

// UnlockedBalanceOf is a free data retrieval call binding the contract method 0x84955c88.
//
// Solidity: function unlockedBalanceOf(address _delegator) view returns(uint256)
func (_ILiquid *ILiquidCallerSession) UnlockedBalanceOf(_delegator common.Address) (*big.Int, error) {
	return _ILiquid.Contract.UnlockedBalanceOf(&_ILiquid.CallOpts, _delegator)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_ILiquid *ILiquidTransactor) Approve(opts *bind.TransactOpts, spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ILiquid.contract.Transact(opts, "approve", spender, amount)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_ILiquid *ILiquidSession) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ILiquid.Contract.Approve(&_ILiquid.TransactOpts, spender, amount)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_ILiquid *ILiquidTransactorSession) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ILiquid.Contract.Approve(&_ILiquid.TransactOpts, spender, amount)
}

// ApproveUnbonding is a paid mutator transaction binding the contract method 0xbf99b73e.
//
// Solidity: function approveUnbonding(address _caller, uint256 _amount) returns(bool)
func (_ILiquid *ILiquidTransactor) ApproveUnbonding(opts *bind.TransactOpts, _caller common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _ILiquid.contract.Transact(opts, "approveUnbonding", _caller, _amount)
}

// ApproveUnbonding is a paid mutator transaction binding the contract method 0xbf99b73e.
//
// Solidity: function approveUnbonding(address _caller, uint256 _amount) returns(bool)
func (_ILiquid *ILiquidSession) ApproveUnbonding(_caller common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _ILiquid.Contract.ApproveUnbonding(&_ILiquid.TransactOpts, _caller, _amount)
}

// ApproveUnbonding is a paid mutator transaction binding the contract method 0xbf99b73e.
//
// Solidity: function approveUnbonding(address _caller, uint256 _amount) returns(bool)
func (_ILiquid *ILiquidTransactorSession) ApproveUnbonding(_caller common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _ILiquid.Contract.ApproveUnbonding(&_ILiquid.TransactOpts, _caller, _amount)
}

// Burn is a paid mutator transaction binding the contract method 0x9dc29fac.
//
// Solidity: function burn(address _account, uint256 _amount) returns()
func (_ILiquid *ILiquidTransactor) Burn(opts *bind.TransactOpts, _account common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _ILiquid.contract.Transact(opts, "burn", _account, _amount)
}

// Burn is a paid mutator transaction binding the contract method 0x9dc29fac.
//
// Solidity: function burn(address _account, uint256 _amount) returns()
func (_ILiquid *ILiquidSession) Burn(_account common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _ILiquid.Contract.Burn(&_ILiquid.TransactOpts, _account, _amount)
}

// Burn is a paid mutator transaction binding the contract method 0x9dc29fac.
//
// Solidity: function burn(address _account, uint256 _amount) returns()
func (_ILiquid *ILiquidTransactorSession) Burn(_account common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _ILiquid.Contract.Burn(&_ILiquid.TransactOpts, _account, _amount)
}

// ClaimRewards is a paid mutator transaction binding the contract method 0x372500ab.
//
// Solidity: function claimRewards() returns()
func (_ILiquid *ILiquidTransactor) ClaimRewards(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ILiquid.contract.Transact(opts, "claimRewards")
}

// ClaimRewards is a paid mutator transaction binding the contract method 0x372500ab.
//
// Solidity: function claimRewards() returns()
func (_ILiquid *ILiquidSession) ClaimRewards() (*types.Transaction, error) {
	return _ILiquid.Contract.ClaimRewards(&_ILiquid.TransactOpts)
}

// ClaimRewards is a paid mutator transaction binding the contract method 0x372500ab.
//
// Solidity: function claimRewards() returns()
func (_ILiquid *ILiquidTransactorSession) ClaimRewards() (*types.Transaction, error) {
	return _ILiquid.Contract.ClaimRewards(&_ILiquid.TransactOpts)
}

// ClaimTreasuryATN is a paid mutator transaction binding the contract method 0xbd96102f.
//
// Solidity: function claimTreasuryATN() returns()
func (_ILiquid *ILiquidTransactor) ClaimTreasuryATN(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ILiquid.contract.Transact(opts, "claimTreasuryATN")
}

// ClaimTreasuryATN is a paid mutator transaction binding the contract method 0xbd96102f.
//
// Solidity: function claimTreasuryATN() returns()
func (_ILiquid *ILiquidSession) ClaimTreasuryATN() (*types.Transaction, error) {
	return _ILiquid.Contract.ClaimTreasuryATN(&_ILiquid.TransactOpts)
}

// ClaimTreasuryATN is a paid mutator transaction binding the contract method 0xbd96102f.
//
// Solidity: function claimTreasuryATN() returns()
func (_ILiquid *ILiquidTransactorSession) ClaimTreasuryATN() (*types.Transaction, error) {
	return _ILiquid.Contract.ClaimTreasuryATN(&_ILiquid.TransactOpts)
}

// Lock is a paid mutator transaction binding the contract method 0x282d3fdf.
//
// Solidity: function lock(address _account, uint256 _amount) returns()
func (_ILiquid *ILiquidTransactor) Lock(opts *bind.TransactOpts, _account common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _ILiquid.contract.Transact(opts, "lock", _account, _amount)
}

// Lock is a paid mutator transaction binding the contract method 0x282d3fdf.
//
// Solidity: function lock(address _account, uint256 _amount) returns()
func (_ILiquid *ILiquidSession) Lock(_account common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _ILiquid.Contract.Lock(&_ILiquid.TransactOpts, _account, _amount)
}

// Lock is a paid mutator transaction binding the contract method 0x282d3fdf.
//
// Solidity: function lock(address _account, uint256 _amount) returns()
func (_ILiquid *ILiquidTransactorSession) Lock(_account common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _ILiquid.Contract.Lock(&_ILiquid.TransactOpts, _account, _amount)
}

// LockFrom is a paid mutator transaction binding the contract method 0x708e91e5.
//
// Solidity: function lockFrom(address _account, address _caller, uint256 _amount) returns()
func (_ILiquid *ILiquidTransactor) LockFrom(opts *bind.TransactOpts, _account common.Address, _caller common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _ILiquid.contract.Transact(opts, "lockFrom", _account, _caller, _amount)
}

// LockFrom is a paid mutator transaction binding the contract method 0x708e91e5.
//
// Solidity: function lockFrom(address _account, address _caller, uint256 _amount) returns()
func (_ILiquid *ILiquidSession) LockFrom(_account common.Address, _caller common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _ILiquid.Contract.LockFrom(&_ILiquid.TransactOpts, _account, _caller, _amount)
}

// LockFrom is a paid mutator transaction binding the contract method 0x708e91e5.
//
// Solidity: function lockFrom(address _account, address _caller, uint256 _amount) returns()
func (_ILiquid *ILiquidTransactorSession) LockFrom(_account common.Address, _caller common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _ILiquid.Contract.LockFrom(&_ILiquid.TransactOpts, _account, _caller, _amount)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address _account, uint256 _amount) returns()
func (_ILiquid *ILiquidTransactor) Mint(opts *bind.TransactOpts, _account common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _ILiquid.contract.Transact(opts, "mint", _account, _amount)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address _account, uint256 _amount) returns()
func (_ILiquid *ILiquidSession) Mint(_account common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _ILiquid.Contract.Mint(&_ILiquid.TransactOpts, _account, _amount)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address _account, uint256 _amount) returns()
func (_ILiquid *ILiquidTransactorSession) Mint(_account common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _ILiquid.Contract.Mint(&_ILiquid.TransactOpts, _account, _amount)
}

// Redistribute is a paid mutator transaction binding the contract method 0xa0ce552d.
//
// Solidity: function redistribute(uint256 _ntnReward) payable returns(uint256)
func (_ILiquid *ILiquidTransactor) Redistribute(opts *bind.TransactOpts, _ntnReward *big.Int) (*types.Transaction, error) {
	return _ILiquid.contract.Transact(opts, "redistribute", _ntnReward)
}

// Redistribute is a paid mutator transaction binding the contract method 0xa0ce552d.
//
// Solidity: function redistribute(uint256 _ntnReward) payable returns(uint256)
func (_ILiquid *ILiquidSession) Redistribute(_ntnReward *big.Int) (*types.Transaction, error) {
	return _ILiquid.Contract.Redistribute(&_ILiquid.TransactOpts, _ntnReward)
}

// Redistribute is a paid mutator transaction binding the contract method 0xa0ce552d.
//
// Solidity: function redistribute(uint256 _ntnReward) payable returns(uint256)
func (_ILiquid *ILiquidTransactorSession) Redistribute(_ntnReward *big.Int) (*types.Transaction, error) {
	return _ILiquid.Contract.Redistribute(&_ILiquid.TransactOpts, _ntnReward)
}

// SetCommissionRate is a paid mutator transaction binding the contract method 0x19fac8fd.
//
// Solidity: function setCommissionRate(uint256 _rate) returns()
func (_ILiquid *ILiquidTransactor) SetCommissionRate(opts *bind.TransactOpts, _rate *big.Int) (*types.Transaction, error) {
	return _ILiquid.contract.Transact(opts, "setCommissionRate", _rate)
}

// SetCommissionRate is a paid mutator transaction binding the contract method 0x19fac8fd.
//
// Solidity: function setCommissionRate(uint256 _rate) returns()
func (_ILiquid *ILiquidSession) SetCommissionRate(_rate *big.Int) (*types.Transaction, error) {
	return _ILiquid.Contract.SetCommissionRate(&_ILiquid.TransactOpts, _rate)
}

// SetCommissionRate is a paid mutator transaction binding the contract method 0x19fac8fd.
//
// Solidity: function setCommissionRate(uint256 _rate) returns()
func (_ILiquid *ILiquidTransactorSession) SetCommissionRate(_rate *big.Int) (*types.Transaction, error) {
	return _ILiquid.Contract.SetCommissionRate(&_ILiquid.TransactOpts, _rate)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_ILiquid *ILiquidTransactor) Transfer(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ILiquid.contract.Transact(opts, "transfer", recipient, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_ILiquid *ILiquidSession) Transfer(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ILiquid.Contract.Transfer(&_ILiquid.TransactOpts, recipient, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_ILiquid *ILiquidTransactorSession) Transfer(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ILiquid.Contract.Transfer(&_ILiquid.TransactOpts, recipient, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_ILiquid *ILiquidTransactor) TransferFrom(opts *bind.TransactOpts, sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ILiquid.contract.Transact(opts, "transferFrom", sender, recipient, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_ILiquid *ILiquidSession) TransferFrom(sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ILiquid.Contract.TransferFrom(&_ILiquid.TransactOpts, sender, recipient, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_ILiquid *ILiquidTransactorSession) TransferFrom(sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ILiquid.Contract.TransferFrom(&_ILiquid.TransactOpts, sender, recipient, amount)
}

// Unlock is a paid mutator transaction binding the contract method 0x7eee288d.
//
// Solidity: function unlock(address _account, uint256 _amount) returns()
func (_ILiquid *ILiquidTransactor) Unlock(opts *bind.TransactOpts, _account common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _ILiquid.contract.Transact(opts, "unlock", _account, _amount)
}

// Unlock is a paid mutator transaction binding the contract method 0x7eee288d.
//
// Solidity: function unlock(address _account, uint256 _amount) returns()
func (_ILiquid *ILiquidSession) Unlock(_account common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _ILiquid.Contract.Unlock(&_ILiquid.TransactOpts, _account, _amount)
}

// Unlock is a paid mutator transaction binding the contract method 0x7eee288d.
//
// Solidity: function unlock(address _account, uint256 _amount) returns()
func (_ILiquid *ILiquidTransactorSession) Unlock(_account common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _ILiquid.Contract.Unlock(&_ILiquid.TransactOpts, _account, _amount)
}

// ILiquidApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the ILiquid contract.
type ILiquidApprovalIterator struct {
	Event *ILiquidApproval // Event containing the contract specifics and raw log

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
func (it *ILiquidApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ILiquidApproval)
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
		it.Event = new(ILiquidApproval)
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
func (it *ILiquidApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ILiquidApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ILiquidApproval represents a Approval event raised by the ILiquid contract.
type ILiquidApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_ILiquid *ILiquidFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*ILiquidApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _ILiquid.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &ILiquidApprovalIterator{contract: _ILiquid.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_ILiquid *ILiquidFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *ILiquidApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _ILiquid.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ILiquidApproval)
				if err := _ILiquid.contract.UnpackLog(event, "Approval", log); err != nil {
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

// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_ILiquid *ILiquidFilterer) ParseApproval(log types.Log) (*ILiquidApproval, error) {
	event := new(ILiquidApproval)
	if err := _ILiquid.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ILiquidTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the ILiquid contract.
type ILiquidTransferIterator struct {
	Event *ILiquidTransfer // Event containing the contract specifics and raw log

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
func (it *ILiquidTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ILiquidTransfer)
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
		it.Event = new(ILiquidTransfer)
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
func (it *ILiquidTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ILiquidTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ILiquidTransfer represents a Transfer event raised by the ILiquid contract.
type ILiquidTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_ILiquid *ILiquidFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*ILiquidTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ILiquid.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &ILiquidTransferIterator{contract: _ILiquid.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_ILiquid *ILiquidFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *ILiquidTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ILiquid.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ILiquidTransfer)
				if err := _ILiquid.contract.UnpackLog(event, "Transfer", log); err != nil {
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
func (_ILiquid *ILiquidFilterer) ParseTransfer(log types.Log) (*ILiquidTransfer, error) {
	event := new(ILiquidTransfer)
	if err := _ILiquid.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ILiquidUnbondingApprovalIterator is returned from FilterUnbondingApproval and is used to iterate over the raw logs and unpacked data for UnbondingApproval events raised by the ILiquid contract.
type ILiquidUnbondingApprovalIterator struct {
	Event *ILiquidUnbondingApproval // Event containing the contract specifics and raw log

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
func (it *ILiquidUnbondingApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ILiquidUnbondingApproval)
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
		it.Event = new(ILiquidUnbondingApproval)
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
func (it *ILiquidUnbondingApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ILiquidUnbondingApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ILiquidUnbondingApproval represents a UnbondingApproval event raised by the ILiquid contract.
type ILiquidUnbondingApproval struct {
	Owner  common.Address
	Caller common.Address
	Value  *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterUnbondingApproval is a free log retrieval operation binding the contract event 0x03530f06053c088596fe8c12632b52218b0685204ad1b757662178f87e27713b.
//
// Solidity: event UnbondingApproval(address indexed _owner, address indexed _caller, uint256 _value)
func (_ILiquid *ILiquidFilterer) FilterUnbondingApproval(opts *bind.FilterOpts, _owner []common.Address, _caller []common.Address) (*ILiquidUnbondingApprovalIterator, error) {

	var _ownerRule []interface{}
	for _, _ownerItem := range _owner {
		_ownerRule = append(_ownerRule, _ownerItem)
	}
	var _callerRule []interface{}
	for _, _callerItem := range _caller {
		_callerRule = append(_callerRule, _callerItem)
	}

	logs, sub, err := _ILiquid.contract.FilterLogs(opts, "UnbondingApproval", _ownerRule, _callerRule)
	if err != nil {
		return nil, err
	}
	return &ILiquidUnbondingApprovalIterator{contract: _ILiquid.contract, event: "UnbondingApproval", logs: logs, sub: sub}, nil
}

// WatchUnbondingApproval is a free log subscription operation binding the contract event 0x03530f06053c088596fe8c12632b52218b0685204ad1b757662178f87e27713b.
//
// Solidity: event UnbondingApproval(address indexed _owner, address indexed _caller, uint256 _value)
func (_ILiquid *ILiquidFilterer) WatchUnbondingApproval(opts *bind.WatchOpts, sink chan<- *ILiquidUnbondingApproval, _owner []common.Address, _caller []common.Address) (event.Subscription, error) {

	var _ownerRule []interface{}
	for _, _ownerItem := range _owner {
		_ownerRule = append(_ownerRule, _ownerItem)
	}
	var _callerRule []interface{}
	for _, _callerItem := range _caller {
		_callerRule = append(_callerRule, _callerItem)
	}

	logs, sub, err := _ILiquid.contract.WatchLogs(opts, "UnbondingApproval", _ownerRule, _callerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ILiquidUnbondingApproval)
				if err := _ILiquid.contract.UnpackLog(event, "UnbondingApproval", log); err != nil {
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

// ParseUnbondingApproval is a log parse operation binding the contract event 0x03530f06053c088596fe8c12632b52218b0685204ad1b757662178f87e27713b.
//
// Solidity: event UnbondingApproval(address indexed _owner, address indexed _caller, uint256 _value)
func (_ILiquid *ILiquidFilterer) ParseUnbondingApproval(log types.Log) (*ILiquidUnbondingApproval, error) {
	event := new(ILiquidUnbondingApproval)
	if err := _ILiquid.contract.UnpackLog(event, "UnbondingApproval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IOmissionAccountabilityMetaData contains all meta data concerning the IOmissionAccountability contract.
var IOmissionAccountabilityMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_ntnReward\",\"type\":\"uint256\"}],\"name\":\"distributeProposerRewards\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"_epochEnded\",\"type\":\"bool\"}],\"name\":\"finalize\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getDelta\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_validator\",\"type\":\"address\"}],\"name\":\"getInactivityScore\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLookbackWindow\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getScaleFactor\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTotalEffort\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"votingPower\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"consensusKey\",\"type\":\"bytes\"}],\"internalType\":\"structIAutonity.CommitteeMember[]\",\"name\":\"_committee\",\"type\":\"tuple[]\"},{\"internalType\":\"address[]\",\"name\":\"_treasuries\",\"type\":\"address[]\"}],\"name\":\"setCommittee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_epochBlock\",\"type\":\"uint256\"}],\"name\":\"setEpochBlock\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_operator\",\"type\":\"address\"}],\"name\":\"setOperator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"eeb92233": "distributeProposerRewards(uint256)",
		"6c9789b0": "finalize(bool)",
		"c549176e": "getDelta()",
		"9a11e0e6": "getInactivityScore(address)",
		"5ca1809c": "getLookbackWindow()",
		"7f5e2f11": "getScaleFactor()",
		"53b1821b": "getTotalEffort()",
		"e3deef9c": "setCommittee((address,uint256,bytes)[],address[])",
		"c024cc2c": "setEpochBlock(uint256)",
		"b3ab15fb": "setOperator(address)",
	},
}

// IOmissionAccountabilityABI is the input ABI used to generate the binding from.
// Deprecated: Use IOmissionAccountabilityMetaData.ABI instead.
var IOmissionAccountabilityABI = IOmissionAccountabilityMetaData.ABI

// Deprecated: Use IOmissionAccountabilityMetaData.Sigs instead.
// IOmissionAccountabilityFuncSigs maps the 4-byte function signature to its string representation.
var IOmissionAccountabilityFuncSigs = IOmissionAccountabilityMetaData.Sigs

// IOmissionAccountability is an auto generated Go binding around an Ethereum contract.
type IOmissionAccountability struct {
	IOmissionAccountabilityCaller     // Read-only binding to the contract
	IOmissionAccountabilityTransactor // Write-only binding to the contract
	IOmissionAccountabilityFilterer   // Log filterer for contract events
}

// IOmissionAccountabilityCaller is an auto generated read-only Go binding around an Ethereum contract.
type IOmissionAccountabilityCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IOmissionAccountabilityTransactor is an auto generated write-only Go binding around an Ethereum contract.
type IOmissionAccountabilityTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IOmissionAccountabilityFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type IOmissionAccountabilityFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IOmissionAccountabilitySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type IOmissionAccountabilitySession struct {
	Contract     *IOmissionAccountability // Generic contract binding to set the session for
	CallOpts     bind.CallOpts            // Call options to use throughout this session
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// IOmissionAccountabilityCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type IOmissionAccountabilityCallerSession struct {
	Contract *IOmissionAccountabilityCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                  // Call options to use throughout this session
}

// IOmissionAccountabilityTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type IOmissionAccountabilityTransactorSession struct {
	Contract     *IOmissionAccountabilityTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                  // Transaction auth options to use throughout this session
}

// IOmissionAccountabilityRaw is an auto generated low-level Go binding around an Ethereum contract.
type IOmissionAccountabilityRaw struct {
	Contract *IOmissionAccountability // Generic contract binding to access the raw methods on
}

// IOmissionAccountabilityCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type IOmissionAccountabilityCallerRaw struct {
	Contract *IOmissionAccountabilityCaller // Generic read-only contract binding to access the raw methods on
}

// IOmissionAccountabilityTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type IOmissionAccountabilityTransactorRaw struct {
	Contract *IOmissionAccountabilityTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIOmissionAccountability creates a new instance of IOmissionAccountability, bound to a specific deployed contract.
func NewIOmissionAccountability(address common.Address, backend bind.ContractBackend) (*IOmissionAccountability, error) {
	contract, err := bindIOmissionAccountability(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IOmissionAccountability{IOmissionAccountabilityCaller: IOmissionAccountabilityCaller{contract: contract}, IOmissionAccountabilityTransactor: IOmissionAccountabilityTransactor{contract: contract}, IOmissionAccountabilityFilterer: IOmissionAccountabilityFilterer{contract: contract}}, nil
}

// NewIOmissionAccountabilityCaller creates a new read-only instance of IOmissionAccountability, bound to a specific deployed contract.
func NewIOmissionAccountabilityCaller(address common.Address, caller bind.ContractCaller) (*IOmissionAccountabilityCaller, error) {
	contract, err := bindIOmissionAccountability(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IOmissionAccountabilityCaller{contract: contract}, nil
}

// NewIOmissionAccountabilityTransactor creates a new write-only instance of IOmissionAccountability, bound to a specific deployed contract.
func NewIOmissionAccountabilityTransactor(address common.Address, transactor bind.ContractTransactor) (*IOmissionAccountabilityTransactor, error) {
	contract, err := bindIOmissionAccountability(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IOmissionAccountabilityTransactor{contract: contract}, nil
}

// NewIOmissionAccountabilityFilterer creates a new log filterer instance of IOmissionAccountability, bound to a specific deployed contract.
func NewIOmissionAccountabilityFilterer(address common.Address, filterer bind.ContractFilterer) (*IOmissionAccountabilityFilterer, error) {
	contract, err := bindIOmissionAccountability(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IOmissionAccountabilityFilterer{contract: contract}, nil
}

// bindIOmissionAccountability binds a generic wrapper to an already deployed contract.
func bindIOmissionAccountability(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(IOmissionAccountabilityABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IOmissionAccountability *IOmissionAccountabilityRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IOmissionAccountability.Contract.IOmissionAccountabilityCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IOmissionAccountability *IOmissionAccountabilityRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IOmissionAccountability.Contract.IOmissionAccountabilityTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IOmissionAccountability *IOmissionAccountabilityRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IOmissionAccountability.Contract.IOmissionAccountabilityTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IOmissionAccountability *IOmissionAccountabilityCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IOmissionAccountability.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IOmissionAccountability *IOmissionAccountabilityTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IOmissionAccountability.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IOmissionAccountability *IOmissionAccountabilityTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IOmissionAccountability.Contract.contract.Transact(opts, method, params...)
}

// GetDelta is a free data retrieval call binding the contract method 0xc549176e.
//
// Solidity: function getDelta() view returns(uint256)
func (_IOmissionAccountability *IOmissionAccountabilityCaller) GetDelta(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IOmissionAccountability.contract.Call(opts, &out, "getDelta")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetDelta is a free data retrieval call binding the contract method 0xc549176e.
//
// Solidity: function getDelta() view returns(uint256)
func (_IOmissionAccountability *IOmissionAccountabilitySession) GetDelta() (*big.Int, error) {
	return _IOmissionAccountability.Contract.GetDelta(&_IOmissionAccountability.CallOpts)
}

// GetDelta is a free data retrieval call binding the contract method 0xc549176e.
//
// Solidity: function getDelta() view returns(uint256)
func (_IOmissionAccountability *IOmissionAccountabilityCallerSession) GetDelta() (*big.Int, error) {
	return _IOmissionAccountability.Contract.GetDelta(&_IOmissionAccountability.CallOpts)
}

// GetInactivityScore is a free data retrieval call binding the contract method 0x9a11e0e6.
//
// Solidity: function getInactivityScore(address _validator) view returns(uint256)
func (_IOmissionAccountability *IOmissionAccountabilityCaller) GetInactivityScore(opts *bind.CallOpts, _validator common.Address) (*big.Int, error) {
	var out []interface{}
	err := _IOmissionAccountability.contract.Call(opts, &out, "getInactivityScore", _validator)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetInactivityScore is a free data retrieval call binding the contract method 0x9a11e0e6.
//
// Solidity: function getInactivityScore(address _validator) view returns(uint256)
func (_IOmissionAccountability *IOmissionAccountabilitySession) GetInactivityScore(_validator common.Address) (*big.Int, error) {
	return _IOmissionAccountability.Contract.GetInactivityScore(&_IOmissionAccountability.CallOpts, _validator)
}

// GetInactivityScore is a free data retrieval call binding the contract method 0x9a11e0e6.
//
// Solidity: function getInactivityScore(address _validator) view returns(uint256)
func (_IOmissionAccountability *IOmissionAccountabilityCallerSession) GetInactivityScore(_validator common.Address) (*big.Int, error) {
	return _IOmissionAccountability.Contract.GetInactivityScore(&_IOmissionAccountability.CallOpts, _validator)
}

// GetLookbackWindow is a free data retrieval call binding the contract method 0x5ca1809c.
//
// Solidity: function getLookbackWindow() view returns(uint256)
func (_IOmissionAccountability *IOmissionAccountabilityCaller) GetLookbackWindow(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IOmissionAccountability.contract.Call(opts, &out, "getLookbackWindow")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetLookbackWindow is a free data retrieval call binding the contract method 0x5ca1809c.
//
// Solidity: function getLookbackWindow() view returns(uint256)
func (_IOmissionAccountability *IOmissionAccountabilitySession) GetLookbackWindow() (*big.Int, error) {
	return _IOmissionAccountability.Contract.GetLookbackWindow(&_IOmissionAccountability.CallOpts)
}

// GetLookbackWindow is a free data retrieval call binding the contract method 0x5ca1809c.
//
// Solidity: function getLookbackWindow() view returns(uint256)
func (_IOmissionAccountability *IOmissionAccountabilityCallerSession) GetLookbackWindow() (*big.Int, error) {
	return _IOmissionAccountability.Contract.GetLookbackWindow(&_IOmissionAccountability.CallOpts)
}

// GetScaleFactor is a free data retrieval call binding the contract method 0x7f5e2f11.
//
// Solidity: function getScaleFactor() pure returns(uint256)
func (_IOmissionAccountability *IOmissionAccountabilityCaller) GetScaleFactor(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IOmissionAccountability.contract.Call(opts, &out, "getScaleFactor")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetScaleFactor is a free data retrieval call binding the contract method 0x7f5e2f11.
//
// Solidity: function getScaleFactor() pure returns(uint256)
func (_IOmissionAccountability *IOmissionAccountabilitySession) GetScaleFactor() (*big.Int, error) {
	return _IOmissionAccountability.Contract.GetScaleFactor(&_IOmissionAccountability.CallOpts)
}

// GetScaleFactor is a free data retrieval call binding the contract method 0x7f5e2f11.
//
// Solidity: function getScaleFactor() pure returns(uint256)
func (_IOmissionAccountability *IOmissionAccountabilityCallerSession) GetScaleFactor() (*big.Int, error) {
	return _IOmissionAccountability.Contract.GetScaleFactor(&_IOmissionAccountability.CallOpts)
}

// GetTotalEffort is a free data retrieval call binding the contract method 0x53b1821b.
//
// Solidity: function getTotalEffort() view returns(uint256)
func (_IOmissionAccountability *IOmissionAccountabilityCaller) GetTotalEffort(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IOmissionAccountability.contract.Call(opts, &out, "getTotalEffort")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetTotalEffort is a free data retrieval call binding the contract method 0x53b1821b.
//
// Solidity: function getTotalEffort() view returns(uint256)
func (_IOmissionAccountability *IOmissionAccountabilitySession) GetTotalEffort() (*big.Int, error) {
	return _IOmissionAccountability.Contract.GetTotalEffort(&_IOmissionAccountability.CallOpts)
}

// GetTotalEffort is a free data retrieval call binding the contract method 0x53b1821b.
//
// Solidity: function getTotalEffort() view returns(uint256)
func (_IOmissionAccountability *IOmissionAccountabilityCallerSession) GetTotalEffort() (*big.Int, error) {
	return _IOmissionAccountability.Contract.GetTotalEffort(&_IOmissionAccountability.CallOpts)
}

// DistributeProposerRewards is a paid mutator transaction binding the contract method 0xeeb92233.
//
// Solidity: function distributeProposerRewards(uint256 _ntnReward) payable returns()
func (_IOmissionAccountability *IOmissionAccountabilityTransactor) DistributeProposerRewards(opts *bind.TransactOpts, _ntnReward *big.Int) (*types.Transaction, error) {
	return _IOmissionAccountability.contract.Transact(opts, "distributeProposerRewards", _ntnReward)
}

// DistributeProposerRewards is a paid mutator transaction binding the contract method 0xeeb92233.
//
// Solidity: function distributeProposerRewards(uint256 _ntnReward) payable returns()
func (_IOmissionAccountability *IOmissionAccountabilitySession) DistributeProposerRewards(_ntnReward *big.Int) (*types.Transaction, error) {
	return _IOmissionAccountability.Contract.DistributeProposerRewards(&_IOmissionAccountability.TransactOpts, _ntnReward)
}

// DistributeProposerRewards is a paid mutator transaction binding the contract method 0xeeb92233.
//
// Solidity: function distributeProposerRewards(uint256 _ntnReward) payable returns()
func (_IOmissionAccountability *IOmissionAccountabilityTransactorSession) DistributeProposerRewards(_ntnReward *big.Int) (*types.Transaction, error) {
	return _IOmissionAccountability.Contract.DistributeProposerRewards(&_IOmissionAccountability.TransactOpts, _ntnReward)
}

// Finalize is a paid mutator transaction binding the contract method 0x6c9789b0.
//
// Solidity: function finalize(bool _epochEnded) returns(uint256)
func (_IOmissionAccountability *IOmissionAccountabilityTransactor) Finalize(opts *bind.TransactOpts, _epochEnded bool) (*types.Transaction, error) {
	return _IOmissionAccountability.contract.Transact(opts, "finalize", _epochEnded)
}

// Finalize is a paid mutator transaction binding the contract method 0x6c9789b0.
//
// Solidity: function finalize(bool _epochEnded) returns(uint256)
func (_IOmissionAccountability *IOmissionAccountabilitySession) Finalize(_epochEnded bool) (*types.Transaction, error) {
	return _IOmissionAccountability.Contract.Finalize(&_IOmissionAccountability.TransactOpts, _epochEnded)
}

// Finalize is a paid mutator transaction binding the contract method 0x6c9789b0.
//
// Solidity: function finalize(bool _epochEnded) returns(uint256)
func (_IOmissionAccountability *IOmissionAccountabilityTransactorSession) Finalize(_epochEnded bool) (*types.Transaction, error) {
	return _IOmissionAccountability.Contract.Finalize(&_IOmissionAccountability.TransactOpts, _epochEnded)
}

// SetCommittee is a paid mutator transaction binding the contract method 0xe3deef9c.
//
// Solidity: function setCommittee((address,uint256,bytes)[] _committee, address[] _treasuries) returns()
func (_IOmissionAccountability *IOmissionAccountabilityTransactor) SetCommittee(opts *bind.TransactOpts, _committee []IAutonityCommitteeMember, _treasuries []common.Address) (*types.Transaction, error) {
	return _IOmissionAccountability.contract.Transact(opts, "setCommittee", _committee, _treasuries)
}

// SetCommittee is a paid mutator transaction binding the contract method 0xe3deef9c.
//
// Solidity: function setCommittee((address,uint256,bytes)[] _committee, address[] _treasuries) returns()
func (_IOmissionAccountability *IOmissionAccountabilitySession) SetCommittee(_committee []IAutonityCommitteeMember, _treasuries []common.Address) (*types.Transaction, error) {
	return _IOmissionAccountability.Contract.SetCommittee(&_IOmissionAccountability.TransactOpts, _committee, _treasuries)
}

// SetCommittee is a paid mutator transaction binding the contract method 0xe3deef9c.
//
// Solidity: function setCommittee((address,uint256,bytes)[] _committee, address[] _treasuries) returns()
func (_IOmissionAccountability *IOmissionAccountabilityTransactorSession) SetCommittee(_committee []IAutonityCommitteeMember, _treasuries []common.Address) (*types.Transaction, error) {
	return _IOmissionAccountability.Contract.SetCommittee(&_IOmissionAccountability.TransactOpts, _committee, _treasuries)
}

// SetEpochBlock is a paid mutator transaction binding the contract method 0xc024cc2c.
//
// Solidity: function setEpochBlock(uint256 _epochBlock) returns()
func (_IOmissionAccountability *IOmissionAccountabilityTransactor) SetEpochBlock(opts *bind.TransactOpts, _epochBlock *big.Int) (*types.Transaction, error) {
	return _IOmissionAccountability.contract.Transact(opts, "setEpochBlock", _epochBlock)
}

// SetEpochBlock is a paid mutator transaction binding the contract method 0xc024cc2c.
//
// Solidity: function setEpochBlock(uint256 _epochBlock) returns()
func (_IOmissionAccountability *IOmissionAccountabilitySession) SetEpochBlock(_epochBlock *big.Int) (*types.Transaction, error) {
	return _IOmissionAccountability.Contract.SetEpochBlock(&_IOmissionAccountability.TransactOpts, _epochBlock)
}

// SetEpochBlock is a paid mutator transaction binding the contract method 0xc024cc2c.
//
// Solidity: function setEpochBlock(uint256 _epochBlock) returns()
func (_IOmissionAccountability *IOmissionAccountabilityTransactorSession) SetEpochBlock(_epochBlock *big.Int) (*types.Transaction, error) {
	return _IOmissionAccountability.Contract.SetEpochBlock(&_IOmissionAccountability.TransactOpts, _epochBlock)
}

// SetOperator is a paid mutator transaction binding the contract method 0xb3ab15fb.
//
// Solidity: function setOperator(address _operator) returns()
func (_IOmissionAccountability *IOmissionAccountabilityTransactor) SetOperator(opts *bind.TransactOpts, _operator common.Address) (*types.Transaction, error) {
	return _IOmissionAccountability.contract.Transact(opts, "setOperator", _operator)
}

// SetOperator is a paid mutator transaction binding the contract method 0xb3ab15fb.
//
// Solidity: function setOperator(address _operator) returns()
func (_IOmissionAccountability *IOmissionAccountabilitySession) SetOperator(_operator common.Address) (*types.Transaction, error) {
	return _IOmissionAccountability.Contract.SetOperator(&_IOmissionAccountability.TransactOpts, _operator)
}

// SetOperator is a paid mutator transaction binding the contract method 0xb3ab15fb.
//
// Solidity: function setOperator(address _operator) returns()
func (_IOmissionAccountability *IOmissionAccountabilityTransactorSession) SetOperator(_operator common.Address) (*types.Transaction, error) {
	return _IOmissionAccountability.Contract.SetOperator(&_IOmissionAccountability.TransactOpts, _operator)
}

// IOracleMetaData contains all meta data concerning the IOracle contract.
var IOracleMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_voter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_round\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_nonRevealCount\",\"type\":\"uint256\"}],\"name\":\"CommitRevealMissed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"cause\",\"type\":\"string\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"reporter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"expValue\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"actualValue\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"extra\",\"type\":\"uint8\"}],\"name\":\"InvalidVote\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_round\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_timestamp\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_votePeriod\",\"type\":\"uint256\"}],\"name\":\"NewRound\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string[]\",\"name\":\"_symbols\",\"type\":\"string[]\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_round\",\"type\":\"uint256\"}],\"name\":\"NewSymbols\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"reporter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"extra\",\"type\":\"uint8\"}],\"name\":\"NewVoter\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_voter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_round\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_missedReveal\",\"type\":\"uint256\"}],\"name\":\"NoRevealPenalty\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_participant\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_slashingAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"_symbol\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"int256\",\"name\":\"_median\",\"type\":\"int256\"},{\"indexed\":false,\"internalType\":\"uint120\",\"name\":\"_reported\",\"type\":\"uint120\"}],\"name\":\"Penalized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"price\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"round\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"status\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"}],\"name\":\"PriceUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"reporter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"extra\",\"type\":\"uint8\"}],\"name\":\"SuccessfulVote\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"ntnReward\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"atnReward\",\"type\":\"uint256\"}],\"name\":\"TotalOracleRewards\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_ntnRewards\",\"type\":\"uint256\"}],\"name\":\"distributeRewards\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"finalize\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getDecimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getNewVotePeriod\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getNewVoters\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getNonRevealThreshold\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRound\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_round\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"_symbol\",\"type\":\"string\"}],\"name\":\"getRoundData\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"round\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"price\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"internalType\":\"structIOracle.RoundData\",\"name\":\"data\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getSymbols\",\"outputs\":[{\"internalType\":\"string[]\",\"name\":\"_symbols\",\"type\":\"string[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getVotePeriod\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getVoters\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_symbol\",\"type\":\"string\"}],\"name\":\"latestRoundData\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"round\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"price\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"internalType\":\"structIOracle.RoundData\",\"name\":\"data\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_threshold\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_resetInterval\",\"type\":\"uint256\"}],\"name\":\"setCommitRevealConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_operator\",\"type\":\"address\"}],\"name\":\"setOperator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"_outlierSlashingThreshold\",\"type\":\"int256\"},{\"internalType\":\"int256\",\"name\":\"_outlierDetectionThreshold\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"_baseSlashingRate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_slashingRateCap\",\"type\":\"uint256\"}],\"name\":\"setSlashingConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string[]\",\"name\":\"_symbols\",\"type\":\"string[]\"}],\"name\":\"setSymbols\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_newVoters\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"_treasury\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"_validator\",\"type\":\"address[]\"}],\"name\":\"setVoters\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"updateVotersAndSymbol\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_commit\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"uint120\",\"name\":\"price\",\"type\":\"uint120\"},{\"internalType\":\"uint8\",\"name\":\"confidence\",\"type\":\"uint8\"}],\"internalType\":\"structIOracle.Report[]\",\"name\":\"_reports\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"_salt\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"_extra\",\"type\":\"uint8\"}],\"name\":\"vote\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"59974e38": "distributeRewards(uint256)",
		"4bb278f3": "finalize()",
		"f0141d84": "getDecimals()",
		"57eba759": "getNewVotePeriod()",
		"077945d3": "getNewVoters()",
		"ed78349d": "getNonRevealThreshold()",
		"9f8743f7": "getRound()",
		"3c8510fd": "getRoundData(uint256,string)",
		"df7f710e": "getSymbols()",
		"b78dec52": "getVotePeriod()",
		"cdd72253": "getVoters()",
		"33f98c77": "latestRoundData(string)",
		"3f422ef3": "setCommitRevealConfig(uint256,uint256)",
		"b3ab15fb": "setOperator(address)",
		"da39fbfe": "setSlashingConfig(int256,int256,uint256,uint256)",
		"8d4f75d2": "setSymbols(string[])",
		"da78110e": "setVoters(address[],address[],address[])",
		"0f65875c": "updateVotersAndSymbol()",
		"56833ebe": "vote(uint256,(uint120,uint8)[],uint256,uint8)",
	},
}

// IOracleABI is the input ABI used to generate the binding from.
// Deprecated: Use IOracleMetaData.ABI instead.
var IOracleABI = IOracleMetaData.ABI

// Deprecated: Use IOracleMetaData.Sigs instead.
// IOracleFuncSigs maps the 4-byte function signature to its string representation.
var IOracleFuncSigs = IOracleMetaData.Sigs

// IOracle is an auto generated Go binding around an Ethereum contract.
type IOracle struct {
	IOracleCaller     // Read-only binding to the contract
	IOracleTransactor // Write-only binding to the contract
	IOracleFilterer   // Log filterer for contract events
}

// IOracleCaller is an auto generated read-only Go binding around an Ethereum contract.
type IOracleCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IOracleTransactor is an auto generated write-only Go binding around an Ethereum contract.
type IOracleTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IOracleFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type IOracleFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IOracleSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type IOracleSession struct {
	Contract     *IOracle          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IOracleCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type IOracleCallerSession struct {
	Contract *IOracleCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// IOracleTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type IOracleTransactorSession struct {
	Contract     *IOracleTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// IOracleRaw is an auto generated low-level Go binding around an Ethereum contract.
type IOracleRaw struct {
	Contract *IOracle // Generic contract binding to access the raw methods on
}

// IOracleCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type IOracleCallerRaw struct {
	Contract *IOracleCaller // Generic read-only contract binding to access the raw methods on
}

// IOracleTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type IOracleTransactorRaw struct {
	Contract *IOracleTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIOracle creates a new instance of IOracle, bound to a specific deployed contract.
func NewIOracle(address common.Address, backend bind.ContractBackend) (*IOracle, error) {
	contract, err := bindIOracle(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IOracle{IOracleCaller: IOracleCaller{contract: contract}, IOracleTransactor: IOracleTransactor{contract: contract}, IOracleFilterer: IOracleFilterer{contract: contract}}, nil
}

// NewIOracleCaller creates a new read-only instance of IOracle, bound to a specific deployed contract.
func NewIOracleCaller(address common.Address, caller bind.ContractCaller) (*IOracleCaller, error) {
	contract, err := bindIOracle(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IOracleCaller{contract: contract}, nil
}

// NewIOracleTransactor creates a new write-only instance of IOracle, bound to a specific deployed contract.
func NewIOracleTransactor(address common.Address, transactor bind.ContractTransactor) (*IOracleTransactor, error) {
	contract, err := bindIOracle(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IOracleTransactor{contract: contract}, nil
}

// NewIOracleFilterer creates a new log filterer instance of IOracle, bound to a specific deployed contract.
func NewIOracleFilterer(address common.Address, filterer bind.ContractFilterer) (*IOracleFilterer, error) {
	contract, err := bindIOracle(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IOracleFilterer{contract: contract}, nil
}

// bindIOracle binds a generic wrapper to an already deployed contract.
func bindIOracle(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(IOracleABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IOracle *IOracleRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IOracle.Contract.IOracleCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IOracle *IOracleRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IOracle.Contract.IOracleTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IOracle *IOracleRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IOracle.Contract.IOracleTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IOracle *IOracleCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IOracle.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IOracle *IOracleTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IOracle.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IOracle *IOracleTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IOracle.Contract.contract.Transact(opts, method, params...)
}

// GetDecimals is a free data retrieval call binding the contract method 0xf0141d84.
//
// Solidity: function getDecimals() view returns(uint8)
func (_IOracle *IOracleCaller) GetDecimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _IOracle.contract.Call(opts, &out, "getDecimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// GetDecimals is a free data retrieval call binding the contract method 0xf0141d84.
//
// Solidity: function getDecimals() view returns(uint8)
func (_IOracle *IOracleSession) GetDecimals() (uint8, error) {
	return _IOracle.Contract.GetDecimals(&_IOracle.CallOpts)
}

// GetDecimals is a free data retrieval call binding the contract method 0xf0141d84.
//
// Solidity: function getDecimals() view returns(uint8)
func (_IOracle *IOracleCallerSession) GetDecimals() (uint8, error) {
	return _IOracle.Contract.GetDecimals(&_IOracle.CallOpts)
}

// GetNewVotePeriod is a free data retrieval call binding the contract method 0x57eba759.
//
// Solidity: function getNewVotePeriod() view returns(uint256)
func (_IOracle *IOracleCaller) GetNewVotePeriod(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IOracle.contract.Call(opts, &out, "getNewVotePeriod")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetNewVotePeriod is a free data retrieval call binding the contract method 0x57eba759.
//
// Solidity: function getNewVotePeriod() view returns(uint256)
func (_IOracle *IOracleSession) GetNewVotePeriod() (*big.Int, error) {
	return _IOracle.Contract.GetNewVotePeriod(&_IOracle.CallOpts)
}

// GetNewVotePeriod is a free data retrieval call binding the contract method 0x57eba759.
//
// Solidity: function getNewVotePeriod() view returns(uint256)
func (_IOracle *IOracleCallerSession) GetNewVotePeriod() (*big.Int, error) {
	return _IOracle.Contract.GetNewVotePeriod(&_IOracle.CallOpts)
}

// GetNewVoters is a free data retrieval call binding the contract method 0x077945d3.
//
// Solidity: function getNewVoters() view returns(address[])
func (_IOracle *IOracleCaller) GetNewVoters(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _IOracle.contract.Call(opts, &out, "getNewVoters")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// GetNewVoters is a free data retrieval call binding the contract method 0x077945d3.
//
// Solidity: function getNewVoters() view returns(address[])
func (_IOracle *IOracleSession) GetNewVoters() ([]common.Address, error) {
	return _IOracle.Contract.GetNewVoters(&_IOracle.CallOpts)
}

// GetNewVoters is a free data retrieval call binding the contract method 0x077945d3.
//
// Solidity: function getNewVoters() view returns(address[])
func (_IOracle *IOracleCallerSession) GetNewVoters() ([]common.Address, error) {
	return _IOracle.Contract.GetNewVoters(&_IOracle.CallOpts)
}

// GetNonRevealThreshold is a free data retrieval call binding the contract method 0xed78349d.
//
// Solidity: function getNonRevealThreshold() view returns(uint256)
func (_IOracle *IOracleCaller) GetNonRevealThreshold(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IOracle.contract.Call(opts, &out, "getNonRevealThreshold")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetNonRevealThreshold is a free data retrieval call binding the contract method 0xed78349d.
//
// Solidity: function getNonRevealThreshold() view returns(uint256)
func (_IOracle *IOracleSession) GetNonRevealThreshold() (*big.Int, error) {
	return _IOracle.Contract.GetNonRevealThreshold(&_IOracle.CallOpts)
}

// GetNonRevealThreshold is a free data retrieval call binding the contract method 0xed78349d.
//
// Solidity: function getNonRevealThreshold() view returns(uint256)
func (_IOracle *IOracleCallerSession) GetNonRevealThreshold() (*big.Int, error) {
	return _IOracle.Contract.GetNonRevealThreshold(&_IOracle.CallOpts)
}

// GetRound is a free data retrieval call binding the contract method 0x9f8743f7.
//
// Solidity: function getRound() view returns(uint256)
func (_IOracle *IOracleCaller) GetRound(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IOracle.contract.Call(opts, &out, "getRound")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetRound is a free data retrieval call binding the contract method 0x9f8743f7.
//
// Solidity: function getRound() view returns(uint256)
func (_IOracle *IOracleSession) GetRound() (*big.Int, error) {
	return _IOracle.Contract.GetRound(&_IOracle.CallOpts)
}

// GetRound is a free data retrieval call binding the contract method 0x9f8743f7.
//
// Solidity: function getRound() view returns(uint256)
func (_IOracle *IOracleCallerSession) GetRound() (*big.Int, error) {
	return _IOracle.Contract.GetRound(&_IOracle.CallOpts)
}

// GetRoundData is a free data retrieval call binding the contract method 0x3c8510fd.
//
// Solidity: function getRoundData(uint256 _round, string _symbol) view returns((uint256,uint256,uint256,bool) data)
func (_IOracle *IOracleCaller) GetRoundData(opts *bind.CallOpts, _round *big.Int, _symbol string) (IOracleRoundData, error) {
	var out []interface{}
	err := _IOracle.contract.Call(opts, &out, "getRoundData", _round, _symbol)

	if err != nil {
		return *new(IOracleRoundData), err
	}

	out0 := *abi.ConvertType(out[0], new(IOracleRoundData)).(*IOracleRoundData)

	return out0, err

}

// GetRoundData is a free data retrieval call binding the contract method 0x3c8510fd.
//
// Solidity: function getRoundData(uint256 _round, string _symbol) view returns((uint256,uint256,uint256,bool) data)
func (_IOracle *IOracleSession) GetRoundData(_round *big.Int, _symbol string) (IOracleRoundData, error) {
	return _IOracle.Contract.GetRoundData(&_IOracle.CallOpts, _round, _symbol)
}

// GetRoundData is a free data retrieval call binding the contract method 0x3c8510fd.
//
// Solidity: function getRoundData(uint256 _round, string _symbol) view returns((uint256,uint256,uint256,bool) data)
func (_IOracle *IOracleCallerSession) GetRoundData(_round *big.Int, _symbol string) (IOracleRoundData, error) {
	return _IOracle.Contract.GetRoundData(&_IOracle.CallOpts, _round, _symbol)
}

// GetSymbols is a free data retrieval call binding the contract method 0xdf7f710e.
//
// Solidity: function getSymbols() view returns(string[] _symbols)
func (_IOracle *IOracleCaller) GetSymbols(opts *bind.CallOpts) ([]string, error) {
	var out []interface{}
	err := _IOracle.contract.Call(opts, &out, "getSymbols")

	if err != nil {
		return *new([]string), err
	}

	out0 := *abi.ConvertType(out[0], new([]string)).(*[]string)

	return out0, err

}

// GetSymbols is a free data retrieval call binding the contract method 0xdf7f710e.
//
// Solidity: function getSymbols() view returns(string[] _symbols)
func (_IOracle *IOracleSession) GetSymbols() ([]string, error) {
	return _IOracle.Contract.GetSymbols(&_IOracle.CallOpts)
}

// GetSymbols is a free data retrieval call binding the contract method 0xdf7f710e.
//
// Solidity: function getSymbols() view returns(string[] _symbols)
func (_IOracle *IOracleCallerSession) GetSymbols() ([]string, error) {
	return _IOracle.Contract.GetSymbols(&_IOracle.CallOpts)
}

// GetVotePeriod is a free data retrieval call binding the contract method 0xb78dec52.
//
// Solidity: function getVotePeriod() view returns(uint256)
func (_IOracle *IOracleCaller) GetVotePeriod(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IOracle.contract.Call(opts, &out, "getVotePeriod")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetVotePeriod is a free data retrieval call binding the contract method 0xb78dec52.
//
// Solidity: function getVotePeriod() view returns(uint256)
func (_IOracle *IOracleSession) GetVotePeriod() (*big.Int, error) {
	return _IOracle.Contract.GetVotePeriod(&_IOracle.CallOpts)
}

// GetVotePeriod is a free data retrieval call binding the contract method 0xb78dec52.
//
// Solidity: function getVotePeriod() view returns(uint256)
func (_IOracle *IOracleCallerSession) GetVotePeriod() (*big.Int, error) {
	return _IOracle.Contract.GetVotePeriod(&_IOracle.CallOpts)
}

// GetVoters is a free data retrieval call binding the contract method 0xcdd72253.
//
// Solidity: function getVoters() view returns(address[])
func (_IOracle *IOracleCaller) GetVoters(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _IOracle.contract.Call(opts, &out, "getVoters")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// GetVoters is a free data retrieval call binding the contract method 0xcdd72253.
//
// Solidity: function getVoters() view returns(address[])
func (_IOracle *IOracleSession) GetVoters() ([]common.Address, error) {
	return _IOracle.Contract.GetVoters(&_IOracle.CallOpts)
}

// GetVoters is a free data retrieval call binding the contract method 0xcdd72253.
//
// Solidity: function getVoters() view returns(address[])
func (_IOracle *IOracleCallerSession) GetVoters() ([]common.Address, error) {
	return _IOracle.Contract.GetVoters(&_IOracle.CallOpts)
}

// LatestRoundData is a free data retrieval call binding the contract method 0x33f98c77.
//
// Solidity: function latestRoundData(string _symbol) view returns((uint256,uint256,uint256,bool) data)
func (_IOracle *IOracleCaller) LatestRoundData(opts *bind.CallOpts, _symbol string) (IOracleRoundData, error) {
	var out []interface{}
	err := _IOracle.contract.Call(opts, &out, "latestRoundData", _symbol)

	if err != nil {
		return *new(IOracleRoundData), err
	}

	out0 := *abi.ConvertType(out[0], new(IOracleRoundData)).(*IOracleRoundData)

	return out0, err

}

// LatestRoundData is a free data retrieval call binding the contract method 0x33f98c77.
//
// Solidity: function latestRoundData(string _symbol) view returns((uint256,uint256,uint256,bool) data)
func (_IOracle *IOracleSession) LatestRoundData(_symbol string) (IOracleRoundData, error) {
	return _IOracle.Contract.LatestRoundData(&_IOracle.CallOpts, _symbol)
}

// LatestRoundData is a free data retrieval call binding the contract method 0x33f98c77.
//
// Solidity: function latestRoundData(string _symbol) view returns((uint256,uint256,uint256,bool) data)
func (_IOracle *IOracleCallerSession) LatestRoundData(_symbol string) (IOracleRoundData, error) {
	return _IOracle.Contract.LatestRoundData(&_IOracle.CallOpts, _symbol)
}

// DistributeRewards is a paid mutator transaction binding the contract method 0x59974e38.
//
// Solidity: function distributeRewards(uint256 _ntnRewards) payable returns()
func (_IOracle *IOracleTransactor) DistributeRewards(opts *bind.TransactOpts, _ntnRewards *big.Int) (*types.Transaction, error) {
	return _IOracle.contract.Transact(opts, "distributeRewards", _ntnRewards)
}

// DistributeRewards is a paid mutator transaction binding the contract method 0x59974e38.
//
// Solidity: function distributeRewards(uint256 _ntnRewards) payable returns()
func (_IOracle *IOracleSession) DistributeRewards(_ntnRewards *big.Int) (*types.Transaction, error) {
	return _IOracle.Contract.DistributeRewards(&_IOracle.TransactOpts, _ntnRewards)
}

// DistributeRewards is a paid mutator transaction binding the contract method 0x59974e38.
//
// Solidity: function distributeRewards(uint256 _ntnRewards) payable returns()
func (_IOracle *IOracleTransactorSession) DistributeRewards(_ntnRewards *big.Int) (*types.Transaction, error) {
	return _IOracle.Contract.DistributeRewards(&_IOracle.TransactOpts, _ntnRewards)
}

// Finalize is a paid mutator transaction binding the contract method 0x4bb278f3.
//
// Solidity: function finalize() returns(bool)
func (_IOracle *IOracleTransactor) Finalize(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IOracle.contract.Transact(opts, "finalize")
}

// Finalize is a paid mutator transaction binding the contract method 0x4bb278f3.
//
// Solidity: function finalize() returns(bool)
func (_IOracle *IOracleSession) Finalize() (*types.Transaction, error) {
	return _IOracle.Contract.Finalize(&_IOracle.TransactOpts)
}

// Finalize is a paid mutator transaction binding the contract method 0x4bb278f3.
//
// Solidity: function finalize() returns(bool)
func (_IOracle *IOracleTransactorSession) Finalize() (*types.Transaction, error) {
	return _IOracle.Contract.Finalize(&_IOracle.TransactOpts)
}

// SetCommitRevealConfig is a paid mutator transaction binding the contract method 0x3f422ef3.
//
// Solidity: function setCommitRevealConfig(uint256 _threshold, uint256 _resetInterval) returns()
func (_IOracle *IOracleTransactor) SetCommitRevealConfig(opts *bind.TransactOpts, _threshold *big.Int, _resetInterval *big.Int) (*types.Transaction, error) {
	return _IOracle.contract.Transact(opts, "setCommitRevealConfig", _threshold, _resetInterval)
}

// SetCommitRevealConfig is a paid mutator transaction binding the contract method 0x3f422ef3.
//
// Solidity: function setCommitRevealConfig(uint256 _threshold, uint256 _resetInterval) returns()
func (_IOracle *IOracleSession) SetCommitRevealConfig(_threshold *big.Int, _resetInterval *big.Int) (*types.Transaction, error) {
	return _IOracle.Contract.SetCommitRevealConfig(&_IOracle.TransactOpts, _threshold, _resetInterval)
}

// SetCommitRevealConfig is a paid mutator transaction binding the contract method 0x3f422ef3.
//
// Solidity: function setCommitRevealConfig(uint256 _threshold, uint256 _resetInterval) returns()
func (_IOracle *IOracleTransactorSession) SetCommitRevealConfig(_threshold *big.Int, _resetInterval *big.Int) (*types.Transaction, error) {
	return _IOracle.Contract.SetCommitRevealConfig(&_IOracle.TransactOpts, _threshold, _resetInterval)
}

// SetOperator is a paid mutator transaction binding the contract method 0xb3ab15fb.
//
// Solidity: function setOperator(address _operator) returns()
func (_IOracle *IOracleTransactor) SetOperator(opts *bind.TransactOpts, _operator common.Address) (*types.Transaction, error) {
	return _IOracle.contract.Transact(opts, "setOperator", _operator)
}

// SetOperator is a paid mutator transaction binding the contract method 0xb3ab15fb.
//
// Solidity: function setOperator(address _operator) returns()
func (_IOracle *IOracleSession) SetOperator(_operator common.Address) (*types.Transaction, error) {
	return _IOracle.Contract.SetOperator(&_IOracle.TransactOpts, _operator)
}

// SetOperator is a paid mutator transaction binding the contract method 0xb3ab15fb.
//
// Solidity: function setOperator(address _operator) returns()
func (_IOracle *IOracleTransactorSession) SetOperator(_operator common.Address) (*types.Transaction, error) {
	return _IOracle.Contract.SetOperator(&_IOracle.TransactOpts, _operator)
}

// SetSlashingConfig is a paid mutator transaction binding the contract method 0xda39fbfe.
//
// Solidity: function setSlashingConfig(int256 _outlierSlashingThreshold, int256 _outlierDetectionThreshold, uint256 _baseSlashingRate, uint256 _slashingRateCap) returns()
func (_IOracle *IOracleTransactor) SetSlashingConfig(opts *bind.TransactOpts, _outlierSlashingThreshold *big.Int, _outlierDetectionThreshold *big.Int, _baseSlashingRate *big.Int, _slashingRateCap *big.Int) (*types.Transaction, error) {
	return _IOracle.contract.Transact(opts, "setSlashingConfig", _outlierSlashingThreshold, _outlierDetectionThreshold, _baseSlashingRate, _slashingRateCap)
}

// SetSlashingConfig is a paid mutator transaction binding the contract method 0xda39fbfe.
//
// Solidity: function setSlashingConfig(int256 _outlierSlashingThreshold, int256 _outlierDetectionThreshold, uint256 _baseSlashingRate, uint256 _slashingRateCap) returns()
func (_IOracle *IOracleSession) SetSlashingConfig(_outlierSlashingThreshold *big.Int, _outlierDetectionThreshold *big.Int, _baseSlashingRate *big.Int, _slashingRateCap *big.Int) (*types.Transaction, error) {
	return _IOracle.Contract.SetSlashingConfig(&_IOracle.TransactOpts, _outlierSlashingThreshold, _outlierDetectionThreshold, _baseSlashingRate, _slashingRateCap)
}

// SetSlashingConfig is a paid mutator transaction binding the contract method 0xda39fbfe.
//
// Solidity: function setSlashingConfig(int256 _outlierSlashingThreshold, int256 _outlierDetectionThreshold, uint256 _baseSlashingRate, uint256 _slashingRateCap) returns()
func (_IOracle *IOracleTransactorSession) SetSlashingConfig(_outlierSlashingThreshold *big.Int, _outlierDetectionThreshold *big.Int, _baseSlashingRate *big.Int, _slashingRateCap *big.Int) (*types.Transaction, error) {
	return _IOracle.Contract.SetSlashingConfig(&_IOracle.TransactOpts, _outlierSlashingThreshold, _outlierDetectionThreshold, _baseSlashingRate, _slashingRateCap)
}

// SetSymbols is a paid mutator transaction binding the contract method 0x8d4f75d2.
//
// Solidity: function setSymbols(string[] _symbols) returns()
func (_IOracle *IOracleTransactor) SetSymbols(opts *bind.TransactOpts, _symbols []string) (*types.Transaction, error) {
	return _IOracle.contract.Transact(opts, "setSymbols", _symbols)
}

// SetSymbols is a paid mutator transaction binding the contract method 0x8d4f75d2.
//
// Solidity: function setSymbols(string[] _symbols) returns()
func (_IOracle *IOracleSession) SetSymbols(_symbols []string) (*types.Transaction, error) {
	return _IOracle.Contract.SetSymbols(&_IOracle.TransactOpts, _symbols)
}

// SetSymbols is a paid mutator transaction binding the contract method 0x8d4f75d2.
//
// Solidity: function setSymbols(string[] _symbols) returns()
func (_IOracle *IOracleTransactorSession) SetSymbols(_symbols []string) (*types.Transaction, error) {
	return _IOracle.Contract.SetSymbols(&_IOracle.TransactOpts, _symbols)
}

// SetVoters is a paid mutator transaction binding the contract method 0xda78110e.
//
// Solidity: function setVoters(address[] _newVoters, address[] _treasury, address[] _validator) returns()
func (_IOracle *IOracleTransactor) SetVoters(opts *bind.TransactOpts, _newVoters []common.Address, _treasury []common.Address, _validator []common.Address) (*types.Transaction, error) {
	return _IOracle.contract.Transact(opts, "setVoters", _newVoters, _treasury, _validator)
}

// SetVoters is a paid mutator transaction binding the contract method 0xda78110e.
//
// Solidity: function setVoters(address[] _newVoters, address[] _treasury, address[] _validator) returns()
func (_IOracle *IOracleSession) SetVoters(_newVoters []common.Address, _treasury []common.Address, _validator []common.Address) (*types.Transaction, error) {
	return _IOracle.Contract.SetVoters(&_IOracle.TransactOpts, _newVoters, _treasury, _validator)
}

// SetVoters is a paid mutator transaction binding the contract method 0xda78110e.
//
// Solidity: function setVoters(address[] _newVoters, address[] _treasury, address[] _validator) returns()
func (_IOracle *IOracleTransactorSession) SetVoters(_newVoters []common.Address, _treasury []common.Address, _validator []common.Address) (*types.Transaction, error) {
	return _IOracle.Contract.SetVoters(&_IOracle.TransactOpts, _newVoters, _treasury, _validator)
}

// UpdateVotersAndSymbol is a paid mutator transaction binding the contract method 0x0f65875c.
//
// Solidity: function updateVotersAndSymbol() returns()
func (_IOracle *IOracleTransactor) UpdateVotersAndSymbol(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IOracle.contract.Transact(opts, "updateVotersAndSymbol")
}

// UpdateVotersAndSymbol is a paid mutator transaction binding the contract method 0x0f65875c.
//
// Solidity: function updateVotersAndSymbol() returns()
func (_IOracle *IOracleSession) UpdateVotersAndSymbol() (*types.Transaction, error) {
	return _IOracle.Contract.UpdateVotersAndSymbol(&_IOracle.TransactOpts)
}

// UpdateVotersAndSymbol is a paid mutator transaction binding the contract method 0x0f65875c.
//
// Solidity: function updateVotersAndSymbol() returns()
func (_IOracle *IOracleTransactorSession) UpdateVotersAndSymbol() (*types.Transaction, error) {
	return _IOracle.Contract.UpdateVotersAndSymbol(&_IOracle.TransactOpts)
}

// Vote is a paid mutator transaction binding the contract method 0x56833ebe.
//
// Solidity: function vote(uint256 _commit, (uint120,uint8)[] _reports, uint256 _salt, uint8 _extra) returns()
func (_IOracle *IOracleTransactor) Vote(opts *bind.TransactOpts, _commit *big.Int, _reports []IOracleReport, _salt *big.Int, _extra uint8) (*types.Transaction, error) {
	return _IOracle.contract.Transact(opts, "vote", _commit, _reports, _salt, _extra)
}

// Vote is a paid mutator transaction binding the contract method 0x56833ebe.
//
// Solidity: function vote(uint256 _commit, (uint120,uint8)[] _reports, uint256 _salt, uint8 _extra) returns()
func (_IOracle *IOracleSession) Vote(_commit *big.Int, _reports []IOracleReport, _salt *big.Int, _extra uint8) (*types.Transaction, error) {
	return _IOracle.Contract.Vote(&_IOracle.TransactOpts, _commit, _reports, _salt, _extra)
}

// Vote is a paid mutator transaction binding the contract method 0x56833ebe.
//
// Solidity: function vote(uint256 _commit, (uint120,uint8)[] _reports, uint256 _salt, uint8 _extra) returns()
func (_IOracle *IOracleTransactorSession) Vote(_commit *big.Int, _reports []IOracleReport, _salt *big.Int, _extra uint8) (*types.Transaction, error) {
	return _IOracle.Contract.Vote(&_IOracle.TransactOpts, _commit, _reports, _salt, _extra)
}

// IOracleCommitRevealMissedIterator is returned from FilterCommitRevealMissed and is used to iterate over the raw logs and unpacked data for CommitRevealMissed events raised by the IOracle contract.
type IOracleCommitRevealMissedIterator struct {
	Event *IOracleCommitRevealMissed // Event containing the contract specifics and raw log

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
func (it *IOracleCommitRevealMissedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IOracleCommitRevealMissed)
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
		it.Event = new(IOracleCommitRevealMissed)
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
func (it *IOracleCommitRevealMissedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IOracleCommitRevealMissedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IOracleCommitRevealMissed represents a CommitRevealMissed event raised by the IOracle contract.
type IOracleCommitRevealMissed struct {
	Voter          common.Address
	Round          *big.Int
	NonRevealCount *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterCommitRevealMissed is a free log retrieval operation binding the contract event 0x176956a4e941f6737f81a3c9a09d8571dd0438d86e25a432beb2013aced43092.
//
// Solidity: event CommitRevealMissed(address indexed _voter, uint256 _round, uint256 _nonRevealCount)
func (_IOracle *IOracleFilterer) FilterCommitRevealMissed(opts *bind.FilterOpts, _voter []common.Address) (*IOracleCommitRevealMissedIterator, error) {

	var _voterRule []interface{}
	for _, _voterItem := range _voter {
		_voterRule = append(_voterRule, _voterItem)
	}

	logs, sub, err := _IOracle.contract.FilterLogs(opts, "CommitRevealMissed", _voterRule)
	if err != nil {
		return nil, err
	}
	return &IOracleCommitRevealMissedIterator{contract: _IOracle.contract, event: "CommitRevealMissed", logs: logs, sub: sub}, nil
}

// WatchCommitRevealMissed is a free log subscription operation binding the contract event 0x176956a4e941f6737f81a3c9a09d8571dd0438d86e25a432beb2013aced43092.
//
// Solidity: event CommitRevealMissed(address indexed _voter, uint256 _round, uint256 _nonRevealCount)
func (_IOracle *IOracleFilterer) WatchCommitRevealMissed(opts *bind.WatchOpts, sink chan<- *IOracleCommitRevealMissed, _voter []common.Address) (event.Subscription, error) {

	var _voterRule []interface{}
	for _, _voterItem := range _voter {
		_voterRule = append(_voterRule, _voterItem)
	}

	logs, sub, err := _IOracle.contract.WatchLogs(opts, "CommitRevealMissed", _voterRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IOracleCommitRevealMissed)
				if err := _IOracle.contract.UnpackLog(event, "CommitRevealMissed", log); err != nil {
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

// ParseCommitRevealMissed is a log parse operation binding the contract event 0x176956a4e941f6737f81a3c9a09d8571dd0438d86e25a432beb2013aced43092.
//
// Solidity: event CommitRevealMissed(address indexed _voter, uint256 _round, uint256 _nonRevealCount)
func (_IOracle *IOracleFilterer) ParseCommitRevealMissed(log types.Log) (*IOracleCommitRevealMissed, error) {
	event := new(IOracleCommitRevealMissed)
	if err := _IOracle.contract.UnpackLog(event, "CommitRevealMissed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IOracleInvalidVoteIterator is returned from FilterInvalidVote and is used to iterate over the raw logs and unpacked data for InvalidVote events raised by the IOracle contract.
type IOracleInvalidVoteIterator struct {
	Event *IOracleInvalidVote // Event containing the contract specifics and raw log

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
func (it *IOracleInvalidVoteIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IOracleInvalidVote)
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
		it.Event = new(IOracleInvalidVote)
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
func (it *IOracleInvalidVoteIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IOracleInvalidVoteIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IOracleInvalidVote represents a InvalidVote event raised by the IOracle contract.
type IOracleInvalidVote struct {
	Cause       string
	Reporter    common.Address
	ExpValue    *big.Int
	ActualValue *big.Int
	Extra       uint8
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterInvalidVote is a free log retrieval operation binding the contract event 0x04ca4e0efda95f8b780c116574d1521309010b38d8f7b75705495703a0f570b1.
//
// Solidity: event InvalidVote(string cause, address indexed reporter, uint256 expValue, uint256 actualValue, uint8 extra)
func (_IOracle *IOracleFilterer) FilterInvalidVote(opts *bind.FilterOpts, reporter []common.Address) (*IOracleInvalidVoteIterator, error) {

	var reporterRule []interface{}
	for _, reporterItem := range reporter {
		reporterRule = append(reporterRule, reporterItem)
	}

	logs, sub, err := _IOracle.contract.FilterLogs(opts, "InvalidVote", reporterRule)
	if err != nil {
		return nil, err
	}
	return &IOracleInvalidVoteIterator{contract: _IOracle.contract, event: "InvalidVote", logs: logs, sub: sub}, nil
}

// WatchInvalidVote is a free log subscription operation binding the contract event 0x04ca4e0efda95f8b780c116574d1521309010b38d8f7b75705495703a0f570b1.
//
// Solidity: event InvalidVote(string cause, address indexed reporter, uint256 expValue, uint256 actualValue, uint8 extra)
func (_IOracle *IOracleFilterer) WatchInvalidVote(opts *bind.WatchOpts, sink chan<- *IOracleInvalidVote, reporter []common.Address) (event.Subscription, error) {

	var reporterRule []interface{}
	for _, reporterItem := range reporter {
		reporterRule = append(reporterRule, reporterItem)
	}

	logs, sub, err := _IOracle.contract.WatchLogs(opts, "InvalidVote", reporterRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IOracleInvalidVote)
				if err := _IOracle.contract.UnpackLog(event, "InvalidVote", log); err != nil {
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

// ParseInvalidVote is a log parse operation binding the contract event 0x04ca4e0efda95f8b780c116574d1521309010b38d8f7b75705495703a0f570b1.
//
// Solidity: event InvalidVote(string cause, address indexed reporter, uint256 expValue, uint256 actualValue, uint8 extra)
func (_IOracle *IOracleFilterer) ParseInvalidVote(log types.Log) (*IOracleInvalidVote, error) {
	event := new(IOracleInvalidVote)
	if err := _IOracle.contract.UnpackLog(event, "InvalidVote", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IOracleNewRoundIterator is returned from FilterNewRound and is used to iterate over the raw logs and unpacked data for NewRound events raised by the IOracle contract.
type IOracleNewRoundIterator struct {
	Event *IOracleNewRound // Event containing the contract specifics and raw log

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
func (it *IOracleNewRoundIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IOracleNewRound)
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
		it.Event = new(IOracleNewRound)
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
func (it *IOracleNewRoundIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IOracleNewRoundIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IOracleNewRound represents a NewRound event raised by the IOracle contract.
type IOracleNewRound struct {
	Round      *big.Int
	Timestamp  *big.Int
	VotePeriod *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterNewRound is a free log retrieval operation binding the contract event 0x5aec57d81928b24d30b1a2aec0d23d693412c37d7ec106b5d8259413716bb1f4.
//
// Solidity: event NewRound(uint256 _round, uint256 _timestamp, uint256 _votePeriod)
func (_IOracle *IOracleFilterer) FilterNewRound(opts *bind.FilterOpts) (*IOracleNewRoundIterator, error) {

	logs, sub, err := _IOracle.contract.FilterLogs(opts, "NewRound")
	if err != nil {
		return nil, err
	}
	return &IOracleNewRoundIterator{contract: _IOracle.contract, event: "NewRound", logs: logs, sub: sub}, nil
}

// WatchNewRound is a free log subscription operation binding the contract event 0x5aec57d81928b24d30b1a2aec0d23d693412c37d7ec106b5d8259413716bb1f4.
//
// Solidity: event NewRound(uint256 _round, uint256 _timestamp, uint256 _votePeriod)
func (_IOracle *IOracleFilterer) WatchNewRound(opts *bind.WatchOpts, sink chan<- *IOracleNewRound) (event.Subscription, error) {

	logs, sub, err := _IOracle.contract.WatchLogs(opts, "NewRound")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IOracleNewRound)
				if err := _IOracle.contract.UnpackLog(event, "NewRound", log); err != nil {
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

// ParseNewRound is a log parse operation binding the contract event 0x5aec57d81928b24d30b1a2aec0d23d693412c37d7ec106b5d8259413716bb1f4.
//
// Solidity: event NewRound(uint256 _round, uint256 _timestamp, uint256 _votePeriod)
func (_IOracle *IOracleFilterer) ParseNewRound(log types.Log) (*IOracleNewRound, error) {
	event := new(IOracleNewRound)
	if err := _IOracle.contract.UnpackLog(event, "NewRound", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IOracleNewSymbolsIterator is returned from FilterNewSymbols and is used to iterate over the raw logs and unpacked data for NewSymbols events raised by the IOracle contract.
type IOracleNewSymbolsIterator struct {
	Event *IOracleNewSymbols // Event containing the contract specifics and raw log

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
func (it *IOracleNewSymbolsIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IOracleNewSymbols)
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
		it.Event = new(IOracleNewSymbols)
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
func (it *IOracleNewSymbolsIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IOracleNewSymbolsIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IOracleNewSymbols represents a NewSymbols event raised by the IOracle contract.
type IOracleNewSymbols struct {
	Symbols []string
	Round   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterNewSymbols is a free log retrieval operation binding the contract event 0xaa278e424da680ce5dad66510415760e78e0bd87d45c786c6e88bdde82f9342d.
//
// Solidity: event NewSymbols(string[] _symbols, uint256 _round)
func (_IOracle *IOracleFilterer) FilterNewSymbols(opts *bind.FilterOpts) (*IOracleNewSymbolsIterator, error) {

	logs, sub, err := _IOracle.contract.FilterLogs(opts, "NewSymbols")
	if err != nil {
		return nil, err
	}
	return &IOracleNewSymbolsIterator{contract: _IOracle.contract, event: "NewSymbols", logs: logs, sub: sub}, nil
}

// WatchNewSymbols is a free log subscription operation binding the contract event 0xaa278e424da680ce5dad66510415760e78e0bd87d45c786c6e88bdde82f9342d.
//
// Solidity: event NewSymbols(string[] _symbols, uint256 _round)
func (_IOracle *IOracleFilterer) WatchNewSymbols(opts *bind.WatchOpts, sink chan<- *IOracleNewSymbols) (event.Subscription, error) {

	logs, sub, err := _IOracle.contract.WatchLogs(opts, "NewSymbols")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IOracleNewSymbols)
				if err := _IOracle.contract.UnpackLog(event, "NewSymbols", log); err != nil {
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

// ParseNewSymbols is a log parse operation binding the contract event 0xaa278e424da680ce5dad66510415760e78e0bd87d45c786c6e88bdde82f9342d.
//
// Solidity: event NewSymbols(string[] _symbols, uint256 _round)
func (_IOracle *IOracleFilterer) ParseNewSymbols(log types.Log) (*IOracleNewSymbols, error) {
	event := new(IOracleNewSymbols)
	if err := _IOracle.contract.UnpackLog(event, "NewSymbols", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IOracleNewVoterIterator is returned from FilterNewVoter and is used to iterate over the raw logs and unpacked data for NewVoter events raised by the IOracle contract.
type IOracleNewVoterIterator struct {
	Event *IOracleNewVoter // Event containing the contract specifics and raw log

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
func (it *IOracleNewVoterIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IOracleNewVoter)
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
		it.Event = new(IOracleNewVoter)
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
func (it *IOracleNewVoterIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IOracleNewVoterIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IOracleNewVoter represents a NewVoter event raised by the IOracle contract.
type IOracleNewVoter struct {
	Reporter common.Address
	Extra    uint8
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterNewVoter is a free log retrieval operation binding the contract event 0xd2ec8e890a03083998d3e16f98044fd3dd13fe3e61b7bc2e58ee6da43b50af73.
//
// Solidity: event NewVoter(address reporter, uint8 extra)
func (_IOracle *IOracleFilterer) FilterNewVoter(opts *bind.FilterOpts) (*IOracleNewVoterIterator, error) {

	logs, sub, err := _IOracle.contract.FilterLogs(opts, "NewVoter")
	if err != nil {
		return nil, err
	}
	return &IOracleNewVoterIterator{contract: _IOracle.contract, event: "NewVoter", logs: logs, sub: sub}, nil
}

// WatchNewVoter is a free log subscription operation binding the contract event 0xd2ec8e890a03083998d3e16f98044fd3dd13fe3e61b7bc2e58ee6da43b50af73.
//
// Solidity: event NewVoter(address reporter, uint8 extra)
func (_IOracle *IOracleFilterer) WatchNewVoter(opts *bind.WatchOpts, sink chan<- *IOracleNewVoter) (event.Subscription, error) {

	logs, sub, err := _IOracle.contract.WatchLogs(opts, "NewVoter")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IOracleNewVoter)
				if err := _IOracle.contract.UnpackLog(event, "NewVoter", log); err != nil {
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

// ParseNewVoter is a log parse operation binding the contract event 0xd2ec8e890a03083998d3e16f98044fd3dd13fe3e61b7bc2e58ee6da43b50af73.
//
// Solidity: event NewVoter(address reporter, uint8 extra)
func (_IOracle *IOracleFilterer) ParseNewVoter(log types.Log) (*IOracleNewVoter, error) {
	event := new(IOracleNewVoter)
	if err := _IOracle.contract.UnpackLog(event, "NewVoter", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IOracleNoRevealPenaltyIterator is returned from FilterNoRevealPenalty and is used to iterate over the raw logs and unpacked data for NoRevealPenalty events raised by the IOracle contract.
type IOracleNoRevealPenaltyIterator struct {
	Event *IOracleNoRevealPenalty // Event containing the contract specifics and raw log

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
func (it *IOracleNoRevealPenaltyIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IOracleNoRevealPenalty)
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
		it.Event = new(IOracleNoRevealPenalty)
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
func (it *IOracleNoRevealPenaltyIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IOracleNoRevealPenaltyIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IOracleNoRevealPenalty represents a NoRevealPenalty event raised by the IOracle contract.
type IOracleNoRevealPenalty struct {
	Voter        common.Address
	Round        *big.Int
	MissedReveal *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterNoRevealPenalty is a free log retrieval operation binding the contract event 0x9e6b40f10c60d1ad09594f3b6ed7043d0e978f584d354ace6e1f6025660c42b1.
//
// Solidity: event NoRevealPenalty(address indexed _voter, uint256 _round, uint256 _missedReveal)
func (_IOracle *IOracleFilterer) FilterNoRevealPenalty(opts *bind.FilterOpts, _voter []common.Address) (*IOracleNoRevealPenaltyIterator, error) {

	var _voterRule []interface{}
	for _, _voterItem := range _voter {
		_voterRule = append(_voterRule, _voterItem)
	}

	logs, sub, err := _IOracle.contract.FilterLogs(opts, "NoRevealPenalty", _voterRule)
	if err != nil {
		return nil, err
	}
	return &IOracleNoRevealPenaltyIterator{contract: _IOracle.contract, event: "NoRevealPenalty", logs: logs, sub: sub}, nil
}

// WatchNoRevealPenalty is a free log subscription operation binding the contract event 0x9e6b40f10c60d1ad09594f3b6ed7043d0e978f584d354ace6e1f6025660c42b1.
//
// Solidity: event NoRevealPenalty(address indexed _voter, uint256 _round, uint256 _missedReveal)
func (_IOracle *IOracleFilterer) WatchNoRevealPenalty(opts *bind.WatchOpts, sink chan<- *IOracleNoRevealPenalty, _voter []common.Address) (event.Subscription, error) {

	var _voterRule []interface{}
	for _, _voterItem := range _voter {
		_voterRule = append(_voterRule, _voterItem)
	}

	logs, sub, err := _IOracle.contract.WatchLogs(opts, "NoRevealPenalty", _voterRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IOracleNoRevealPenalty)
				if err := _IOracle.contract.UnpackLog(event, "NoRevealPenalty", log); err != nil {
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

// ParseNoRevealPenalty is a log parse operation binding the contract event 0x9e6b40f10c60d1ad09594f3b6ed7043d0e978f584d354ace6e1f6025660c42b1.
//
// Solidity: event NoRevealPenalty(address indexed _voter, uint256 _round, uint256 _missedReveal)
func (_IOracle *IOracleFilterer) ParseNoRevealPenalty(log types.Log) (*IOracleNoRevealPenalty, error) {
	event := new(IOracleNoRevealPenalty)
	if err := _IOracle.contract.UnpackLog(event, "NoRevealPenalty", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IOraclePenalizedIterator is returned from FilterPenalized and is used to iterate over the raw logs and unpacked data for Penalized events raised by the IOracle contract.
type IOraclePenalizedIterator struct {
	Event *IOraclePenalized // Event containing the contract specifics and raw log

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
func (it *IOraclePenalizedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IOraclePenalized)
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
		it.Event = new(IOraclePenalized)
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
func (it *IOraclePenalizedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IOraclePenalizedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IOraclePenalized represents a Penalized event raised by the IOracle contract.
type IOraclePenalized struct {
	Participant    common.Address
	SlashingAmount *big.Int
	Symbol         string
	Median         *big.Int
	Reported       *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterPenalized is a free log retrieval operation binding the contract event 0x372858b237c8bd0714183e8351a461d6c3cb1ef83806181b36bf5943711f4f57.
//
// Solidity: event Penalized(address indexed _participant, uint256 _slashingAmount, string _symbol, int256 _median, uint120 _reported)
func (_IOracle *IOracleFilterer) FilterPenalized(opts *bind.FilterOpts, _participant []common.Address) (*IOraclePenalizedIterator, error) {

	var _participantRule []interface{}
	for _, _participantItem := range _participant {
		_participantRule = append(_participantRule, _participantItem)
	}

	logs, sub, err := _IOracle.contract.FilterLogs(opts, "Penalized", _participantRule)
	if err != nil {
		return nil, err
	}
	return &IOraclePenalizedIterator{contract: _IOracle.contract, event: "Penalized", logs: logs, sub: sub}, nil
}

// WatchPenalized is a free log subscription operation binding the contract event 0x372858b237c8bd0714183e8351a461d6c3cb1ef83806181b36bf5943711f4f57.
//
// Solidity: event Penalized(address indexed _participant, uint256 _slashingAmount, string _symbol, int256 _median, uint120 _reported)
func (_IOracle *IOracleFilterer) WatchPenalized(opts *bind.WatchOpts, sink chan<- *IOraclePenalized, _participant []common.Address) (event.Subscription, error) {

	var _participantRule []interface{}
	for _, _participantItem := range _participant {
		_participantRule = append(_participantRule, _participantItem)
	}

	logs, sub, err := _IOracle.contract.WatchLogs(opts, "Penalized", _participantRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IOraclePenalized)
				if err := _IOracle.contract.UnpackLog(event, "Penalized", log); err != nil {
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

// ParsePenalized is a log parse operation binding the contract event 0x372858b237c8bd0714183e8351a461d6c3cb1ef83806181b36bf5943711f4f57.
//
// Solidity: event Penalized(address indexed _participant, uint256 _slashingAmount, string _symbol, int256 _median, uint120 _reported)
func (_IOracle *IOracleFilterer) ParsePenalized(log types.Log) (*IOraclePenalized, error) {
	event := new(IOraclePenalized)
	if err := _IOracle.contract.UnpackLog(event, "Penalized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IOraclePriceUpdatedIterator is returned from FilterPriceUpdated and is used to iterate over the raw logs and unpacked data for PriceUpdated events raised by the IOracle contract.
type IOraclePriceUpdatedIterator struct {
	Event *IOraclePriceUpdated // Event containing the contract specifics and raw log

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
func (it *IOraclePriceUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IOraclePriceUpdated)
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
		it.Event = new(IOraclePriceUpdated)
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
func (it *IOraclePriceUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IOraclePriceUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IOraclePriceUpdated represents a PriceUpdated event raised by the IOracle contract.
type IOraclePriceUpdated struct {
	Price     *big.Int
	Round     *big.Int
	Symbol    common.Hash
	Status    bool
	Timestamp *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterPriceUpdated is a free log retrieval operation binding the contract event 0x5f2aa51aa7889ad71d9318fa7fd83c8ff3277434249bd06073f15986e197911c.
//
// Solidity: event PriceUpdated(uint256 price, uint256 round, string indexed symbol, bool status, uint256 timestamp)
func (_IOracle *IOracleFilterer) FilterPriceUpdated(opts *bind.FilterOpts, symbol []string) (*IOraclePriceUpdatedIterator, error) {

	var symbolRule []interface{}
	for _, symbolItem := range symbol {
		symbolRule = append(symbolRule, symbolItem)
	}

	logs, sub, err := _IOracle.contract.FilterLogs(opts, "PriceUpdated", symbolRule)
	if err != nil {
		return nil, err
	}
	return &IOraclePriceUpdatedIterator{contract: _IOracle.contract, event: "PriceUpdated", logs: logs, sub: sub}, nil
}

// WatchPriceUpdated is a free log subscription operation binding the contract event 0x5f2aa51aa7889ad71d9318fa7fd83c8ff3277434249bd06073f15986e197911c.
//
// Solidity: event PriceUpdated(uint256 price, uint256 round, string indexed symbol, bool status, uint256 timestamp)
func (_IOracle *IOracleFilterer) WatchPriceUpdated(opts *bind.WatchOpts, sink chan<- *IOraclePriceUpdated, symbol []string) (event.Subscription, error) {

	var symbolRule []interface{}
	for _, symbolItem := range symbol {
		symbolRule = append(symbolRule, symbolItem)
	}

	logs, sub, err := _IOracle.contract.WatchLogs(opts, "PriceUpdated", symbolRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IOraclePriceUpdated)
				if err := _IOracle.contract.UnpackLog(event, "PriceUpdated", log); err != nil {
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

// ParsePriceUpdated is a log parse operation binding the contract event 0x5f2aa51aa7889ad71d9318fa7fd83c8ff3277434249bd06073f15986e197911c.
//
// Solidity: event PriceUpdated(uint256 price, uint256 round, string indexed symbol, bool status, uint256 timestamp)
func (_IOracle *IOracleFilterer) ParsePriceUpdated(log types.Log) (*IOraclePriceUpdated, error) {
	event := new(IOraclePriceUpdated)
	if err := _IOracle.contract.UnpackLog(event, "PriceUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IOracleSuccessfulVoteIterator is returned from FilterSuccessfulVote and is used to iterate over the raw logs and unpacked data for SuccessfulVote events raised by the IOracle contract.
type IOracleSuccessfulVoteIterator struct {
	Event *IOracleSuccessfulVote // Event containing the contract specifics and raw log

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
func (it *IOracleSuccessfulVoteIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IOracleSuccessfulVote)
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
		it.Event = new(IOracleSuccessfulVote)
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
func (it *IOracleSuccessfulVoteIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IOracleSuccessfulVoteIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IOracleSuccessfulVote represents a SuccessfulVote event raised by the IOracle contract.
type IOracleSuccessfulVote struct {
	Reporter common.Address
	Extra    uint8
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterSuccessfulVote is a free log retrieval operation binding the contract event 0x8bdddd7f2f2c74679ffa6beb8f86aa18bfa5baf1bfaf534d0b66596babc53f08.
//
// Solidity: event SuccessfulVote(address indexed reporter, uint8 extra)
func (_IOracle *IOracleFilterer) FilterSuccessfulVote(opts *bind.FilterOpts, reporter []common.Address) (*IOracleSuccessfulVoteIterator, error) {

	var reporterRule []interface{}
	for _, reporterItem := range reporter {
		reporterRule = append(reporterRule, reporterItem)
	}

	logs, sub, err := _IOracle.contract.FilterLogs(opts, "SuccessfulVote", reporterRule)
	if err != nil {
		return nil, err
	}
	return &IOracleSuccessfulVoteIterator{contract: _IOracle.contract, event: "SuccessfulVote", logs: logs, sub: sub}, nil
}

// WatchSuccessfulVote is a free log subscription operation binding the contract event 0x8bdddd7f2f2c74679ffa6beb8f86aa18bfa5baf1bfaf534d0b66596babc53f08.
//
// Solidity: event SuccessfulVote(address indexed reporter, uint8 extra)
func (_IOracle *IOracleFilterer) WatchSuccessfulVote(opts *bind.WatchOpts, sink chan<- *IOracleSuccessfulVote, reporter []common.Address) (event.Subscription, error) {

	var reporterRule []interface{}
	for _, reporterItem := range reporter {
		reporterRule = append(reporterRule, reporterItem)
	}

	logs, sub, err := _IOracle.contract.WatchLogs(opts, "SuccessfulVote", reporterRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IOracleSuccessfulVote)
				if err := _IOracle.contract.UnpackLog(event, "SuccessfulVote", log); err != nil {
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

// ParseSuccessfulVote is a log parse operation binding the contract event 0x8bdddd7f2f2c74679ffa6beb8f86aa18bfa5baf1bfaf534d0b66596babc53f08.
//
// Solidity: event SuccessfulVote(address indexed reporter, uint8 extra)
func (_IOracle *IOracleFilterer) ParseSuccessfulVote(log types.Log) (*IOracleSuccessfulVote, error) {
	event := new(IOracleSuccessfulVote)
	if err := _IOracle.contract.UnpackLog(event, "SuccessfulVote", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IOracleTotalOracleRewardsIterator is returned from FilterTotalOracleRewards and is used to iterate over the raw logs and unpacked data for TotalOracleRewards events raised by the IOracle contract.
type IOracleTotalOracleRewardsIterator struct {
	Event *IOracleTotalOracleRewards // Event containing the contract specifics and raw log

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
func (it *IOracleTotalOracleRewardsIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IOracleTotalOracleRewards)
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
		it.Event = new(IOracleTotalOracleRewards)
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
func (it *IOracleTotalOracleRewardsIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IOracleTotalOracleRewardsIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IOracleTotalOracleRewards represents a TotalOracleRewards event raised by the IOracle contract.
type IOracleTotalOracleRewards struct {
	NtnReward *big.Int
	AtnReward *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterTotalOracleRewards is a free log retrieval operation binding the contract event 0x3e5aaff9e8fd4293ae18127809c2d4069d87fe10c7de92aa39557a1edbd48fec.
//
// Solidity: event TotalOracleRewards(uint256 ntnReward, uint256 atnReward)
func (_IOracle *IOracleFilterer) FilterTotalOracleRewards(opts *bind.FilterOpts) (*IOracleTotalOracleRewardsIterator, error) {

	logs, sub, err := _IOracle.contract.FilterLogs(opts, "TotalOracleRewards")
	if err != nil {
		return nil, err
	}
	return &IOracleTotalOracleRewardsIterator{contract: _IOracle.contract, event: "TotalOracleRewards", logs: logs, sub: sub}, nil
}

// WatchTotalOracleRewards is a free log subscription operation binding the contract event 0x3e5aaff9e8fd4293ae18127809c2d4069d87fe10c7de92aa39557a1edbd48fec.
//
// Solidity: event TotalOracleRewards(uint256 ntnReward, uint256 atnReward)
func (_IOracle *IOracleFilterer) WatchTotalOracleRewards(opts *bind.WatchOpts, sink chan<- *IOracleTotalOracleRewards) (event.Subscription, error) {

	logs, sub, err := _IOracle.contract.WatchLogs(opts, "TotalOracleRewards")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IOracleTotalOracleRewards)
				if err := _IOracle.contract.UnpackLog(event, "TotalOracleRewards", log); err != nil {
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

// ParseTotalOracleRewards is a log parse operation binding the contract event 0x3e5aaff9e8fd4293ae18127809c2d4069d87fe10c7de92aa39557a1edbd48fec.
//
// Solidity: event TotalOracleRewards(uint256 ntnReward, uint256 atnReward)
func (_IOracle *IOracleFilterer) ParseTotalOracleRewards(log types.Log) (*IOracleTotalOracleRewards, error) {
	event := new(IOracleTotalOracleRewards)
	if err := _IOracle.contract.UnpackLog(event, "TotalOracleRewards", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IScheduleControllerMetaData contains all meta data concerning the IScheduleController contract.
var IScheduleControllerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_vault\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_id\",\"type\":\"uint256\"}],\"name\":\"getSchedule\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"totalAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"unlockedAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"start\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalDuration\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"lastUnlockTime\",\"type\":\"uint256\"}],\"internalType\":\"structIScheduleController.Schedule\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_vault\",\"type\":\"address\"}],\"name\":\"getTotalSchedules\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"7264c4da": "getSchedule(address,uint256)",
		"088566e9": "getTotalSchedules(address)",
	},
}

// IScheduleControllerABI is the input ABI used to generate the binding from.
// Deprecated: Use IScheduleControllerMetaData.ABI instead.
var IScheduleControllerABI = IScheduleControllerMetaData.ABI

// Deprecated: Use IScheduleControllerMetaData.Sigs instead.
// IScheduleControllerFuncSigs maps the 4-byte function signature to its string representation.
var IScheduleControllerFuncSigs = IScheduleControllerMetaData.Sigs

// IScheduleController is an auto generated Go binding around an Ethereum contract.
type IScheduleController struct {
	IScheduleControllerCaller     // Read-only binding to the contract
	IScheduleControllerTransactor // Write-only binding to the contract
	IScheduleControllerFilterer   // Log filterer for contract events
}

// IScheduleControllerCaller is an auto generated read-only Go binding around an Ethereum contract.
type IScheduleControllerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IScheduleControllerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type IScheduleControllerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IScheduleControllerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type IScheduleControllerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IScheduleControllerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type IScheduleControllerSession struct {
	Contract     *IScheduleController // Generic contract binding to set the session for
	CallOpts     bind.CallOpts        // Call options to use throughout this session
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// IScheduleControllerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type IScheduleControllerCallerSession struct {
	Contract *IScheduleControllerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts              // Call options to use throughout this session
}

// IScheduleControllerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type IScheduleControllerTransactorSession struct {
	Contract     *IScheduleControllerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts              // Transaction auth options to use throughout this session
}

// IScheduleControllerRaw is an auto generated low-level Go binding around an Ethereum contract.
type IScheduleControllerRaw struct {
	Contract *IScheduleController // Generic contract binding to access the raw methods on
}

// IScheduleControllerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type IScheduleControllerCallerRaw struct {
	Contract *IScheduleControllerCaller // Generic read-only contract binding to access the raw methods on
}

// IScheduleControllerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type IScheduleControllerTransactorRaw struct {
	Contract *IScheduleControllerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIScheduleController creates a new instance of IScheduleController, bound to a specific deployed contract.
func NewIScheduleController(address common.Address, backend bind.ContractBackend) (*IScheduleController, error) {
	contract, err := bindIScheduleController(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IScheduleController{IScheduleControllerCaller: IScheduleControllerCaller{contract: contract}, IScheduleControllerTransactor: IScheduleControllerTransactor{contract: contract}, IScheduleControllerFilterer: IScheduleControllerFilterer{contract: contract}}, nil
}

// NewIScheduleControllerCaller creates a new read-only instance of IScheduleController, bound to a specific deployed contract.
func NewIScheduleControllerCaller(address common.Address, caller bind.ContractCaller) (*IScheduleControllerCaller, error) {
	contract, err := bindIScheduleController(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IScheduleControllerCaller{contract: contract}, nil
}

// NewIScheduleControllerTransactor creates a new write-only instance of IScheduleController, bound to a specific deployed contract.
func NewIScheduleControllerTransactor(address common.Address, transactor bind.ContractTransactor) (*IScheduleControllerTransactor, error) {
	contract, err := bindIScheduleController(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IScheduleControllerTransactor{contract: contract}, nil
}

// NewIScheduleControllerFilterer creates a new log filterer instance of IScheduleController, bound to a specific deployed contract.
func NewIScheduleControllerFilterer(address common.Address, filterer bind.ContractFilterer) (*IScheduleControllerFilterer, error) {
	contract, err := bindIScheduleController(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IScheduleControllerFilterer{contract: contract}, nil
}

// bindIScheduleController binds a generic wrapper to an already deployed contract.
func bindIScheduleController(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(IScheduleControllerABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IScheduleController *IScheduleControllerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IScheduleController.Contract.IScheduleControllerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IScheduleController *IScheduleControllerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IScheduleController.Contract.IScheduleControllerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IScheduleController *IScheduleControllerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IScheduleController.Contract.IScheduleControllerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IScheduleController *IScheduleControllerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IScheduleController.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IScheduleController *IScheduleControllerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IScheduleController.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IScheduleController *IScheduleControllerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IScheduleController.Contract.contract.Transact(opts, method, params...)
}

// GetSchedule is a free data retrieval call binding the contract method 0x7264c4da.
//
// Solidity: function getSchedule(address _vault, uint256 _id) view returns((uint256,uint256,uint256,uint256,uint256))
func (_IScheduleController *IScheduleControllerCaller) GetSchedule(opts *bind.CallOpts, _vault common.Address, _id *big.Int) (IScheduleControllerSchedule, error) {
	var out []interface{}
	err := _IScheduleController.contract.Call(opts, &out, "getSchedule", _vault, _id)

	if err != nil {
		return *new(IScheduleControllerSchedule), err
	}

	out0 := *abi.ConvertType(out[0], new(IScheduleControllerSchedule)).(*IScheduleControllerSchedule)

	return out0, err

}

// GetSchedule is a free data retrieval call binding the contract method 0x7264c4da.
//
// Solidity: function getSchedule(address _vault, uint256 _id) view returns((uint256,uint256,uint256,uint256,uint256))
func (_IScheduleController *IScheduleControllerSession) GetSchedule(_vault common.Address, _id *big.Int) (IScheduleControllerSchedule, error) {
	return _IScheduleController.Contract.GetSchedule(&_IScheduleController.CallOpts, _vault, _id)
}

// GetSchedule is a free data retrieval call binding the contract method 0x7264c4da.
//
// Solidity: function getSchedule(address _vault, uint256 _id) view returns((uint256,uint256,uint256,uint256,uint256))
func (_IScheduleController *IScheduleControllerCallerSession) GetSchedule(_vault common.Address, _id *big.Int) (IScheduleControllerSchedule, error) {
	return _IScheduleController.Contract.GetSchedule(&_IScheduleController.CallOpts, _vault, _id)
}

// GetTotalSchedules is a free data retrieval call binding the contract method 0x088566e9.
//
// Solidity: function getTotalSchedules(address _vault) view returns(uint256)
func (_IScheduleController *IScheduleControllerCaller) GetTotalSchedules(opts *bind.CallOpts, _vault common.Address) (*big.Int, error) {
	var out []interface{}
	err := _IScheduleController.contract.Call(opts, &out, "getTotalSchedules", _vault)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetTotalSchedules is a free data retrieval call binding the contract method 0x088566e9.
//
// Solidity: function getTotalSchedules(address _vault) view returns(uint256)
func (_IScheduleController *IScheduleControllerSession) GetTotalSchedules(_vault common.Address) (*big.Int, error) {
	return _IScheduleController.Contract.GetTotalSchedules(&_IScheduleController.CallOpts, _vault)
}

// GetTotalSchedules is a free data retrieval call binding the contract method 0x088566e9.
//
// Solidity: function getTotalSchedules(address _vault) view returns(uint256)
func (_IScheduleController *IScheduleControllerCallerSession) GetTotalSchedules(_vault common.Address) (*big.Int, error) {
	return _IScheduleController.Contract.GetTotalSchedules(&_IScheduleController.CallOpts, _vault)
}

// IStabilizationMetaData contains all meta data concerning the IStabilization contract.
var IStabilizationMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"cdps\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"collateral\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"principal\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"interest\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"lastAggregatedInterestExponent\",\"type\":\"uint256\"}],\"internalType\":\"structIStabilization.CDP\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"collateralPrice\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"price\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"config\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"borrowInterestRate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"announcementWindow\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"liquidationRatio\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"minCollateralizationRatio\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"minDebtRequirement\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"targetPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"defaultNTNATNPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"defaultNTNUSDPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"defaultACUUSDPrice\",\"type\":\"uint256\"}],\"internalType\":\"structIStabilization.Config\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"}],\"name\":\"debtAmountAtTime\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastUpdated\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"borrowInterestRateTimestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"announcementWindowTimestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"liquidationRatioTimestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"minCollateralizationRatioTimestamp\",\"type\":\"uint256\"}],\"internalType\":\"structIStabilization.LastUpdated\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"collateralSold\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"bidder\",\"type\":\"address\"}],\"name\":\"liquidate\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"acu\",\"type\":\"address\"}],\"name\":\"setACU\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"auctioneer\",\"type\":\"address\"}],\"name\":\"setAuctioneer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"setOperator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"oracle\",\"type\":\"address\"}],\"name\":\"setOracle\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"supplyControl\",\"type\":\"address\"}],\"name\":\"setSupplyControl\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"840c7e24": "cdps(address)",
		"5891de72": "collateralPrice()",
		"79502c55": "config()",
		"8ffba9dd": "debtAmountAtTime(address,uint256)",
		"d0b06f5d": "lastUpdated()",
		"4914c008": "liquidate(address,uint256,address)",
		"4b8ef943": "setACU(address)",
		"00ede7e4": "setAuctioneer(address)",
		"b3ab15fb": "setOperator(address)",
		"7adbf973": "setOracle(address)",
		"52e5a050": "setSupplyControl(address)",
	},
}

// IStabilizationABI is the input ABI used to generate the binding from.
// Deprecated: Use IStabilizationMetaData.ABI instead.
var IStabilizationABI = IStabilizationMetaData.ABI

// Deprecated: Use IStabilizationMetaData.Sigs instead.
// IStabilizationFuncSigs maps the 4-byte function signature to its string representation.
var IStabilizationFuncSigs = IStabilizationMetaData.Sigs

// IStabilization is an auto generated Go binding around an Ethereum contract.
type IStabilization struct {
	IStabilizationCaller     // Read-only binding to the contract
	IStabilizationTransactor // Write-only binding to the contract
	IStabilizationFilterer   // Log filterer for contract events
}

// IStabilizationCaller is an auto generated read-only Go binding around an Ethereum contract.
type IStabilizationCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IStabilizationTransactor is an auto generated write-only Go binding around an Ethereum contract.
type IStabilizationTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IStabilizationFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type IStabilizationFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IStabilizationSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type IStabilizationSession struct {
	Contract     *IStabilization   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IStabilizationCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type IStabilizationCallerSession struct {
	Contract *IStabilizationCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// IStabilizationTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type IStabilizationTransactorSession struct {
	Contract     *IStabilizationTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// IStabilizationRaw is an auto generated low-level Go binding around an Ethereum contract.
type IStabilizationRaw struct {
	Contract *IStabilization // Generic contract binding to access the raw methods on
}

// IStabilizationCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type IStabilizationCallerRaw struct {
	Contract *IStabilizationCaller // Generic read-only contract binding to access the raw methods on
}

// IStabilizationTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type IStabilizationTransactorRaw struct {
	Contract *IStabilizationTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIStabilization creates a new instance of IStabilization, bound to a specific deployed contract.
func NewIStabilization(address common.Address, backend bind.ContractBackend) (*IStabilization, error) {
	contract, err := bindIStabilization(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IStabilization{IStabilizationCaller: IStabilizationCaller{contract: contract}, IStabilizationTransactor: IStabilizationTransactor{contract: contract}, IStabilizationFilterer: IStabilizationFilterer{contract: contract}}, nil
}

// NewIStabilizationCaller creates a new read-only instance of IStabilization, bound to a specific deployed contract.
func NewIStabilizationCaller(address common.Address, caller bind.ContractCaller) (*IStabilizationCaller, error) {
	contract, err := bindIStabilization(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IStabilizationCaller{contract: contract}, nil
}

// NewIStabilizationTransactor creates a new write-only instance of IStabilization, bound to a specific deployed contract.
func NewIStabilizationTransactor(address common.Address, transactor bind.ContractTransactor) (*IStabilizationTransactor, error) {
	contract, err := bindIStabilization(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IStabilizationTransactor{contract: contract}, nil
}

// NewIStabilizationFilterer creates a new log filterer instance of IStabilization, bound to a specific deployed contract.
func NewIStabilizationFilterer(address common.Address, filterer bind.ContractFilterer) (*IStabilizationFilterer, error) {
	contract, err := bindIStabilization(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IStabilizationFilterer{contract: contract}, nil
}

// bindIStabilization binds a generic wrapper to an already deployed contract.
func bindIStabilization(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(IStabilizationABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IStabilization *IStabilizationRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IStabilization.Contract.IStabilizationCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IStabilization *IStabilizationRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IStabilization.Contract.IStabilizationTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IStabilization *IStabilizationRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IStabilization.Contract.IStabilizationTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IStabilization *IStabilizationCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IStabilization.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IStabilization *IStabilizationTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IStabilization.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IStabilization *IStabilizationTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IStabilization.Contract.contract.Transact(opts, method, params...)
}

// Cdps is a free data retrieval call binding the contract method 0x840c7e24.
//
// Solidity: function cdps(address owner) view returns((uint256,uint256,uint256,uint256,uint256))
func (_IStabilization *IStabilizationCaller) Cdps(opts *bind.CallOpts, owner common.Address) (IStabilizationCDP, error) {
	var out []interface{}
	err := _IStabilization.contract.Call(opts, &out, "cdps", owner)

	if err != nil {
		return *new(IStabilizationCDP), err
	}

	out0 := *abi.ConvertType(out[0], new(IStabilizationCDP)).(*IStabilizationCDP)

	return out0, err

}

// Cdps is a free data retrieval call binding the contract method 0x840c7e24.
//
// Solidity: function cdps(address owner) view returns((uint256,uint256,uint256,uint256,uint256))
func (_IStabilization *IStabilizationSession) Cdps(owner common.Address) (IStabilizationCDP, error) {
	return _IStabilization.Contract.Cdps(&_IStabilization.CallOpts, owner)
}

// Cdps is a free data retrieval call binding the contract method 0x840c7e24.
//
// Solidity: function cdps(address owner) view returns((uint256,uint256,uint256,uint256,uint256))
func (_IStabilization *IStabilizationCallerSession) Cdps(owner common.Address) (IStabilizationCDP, error) {
	return _IStabilization.Contract.Cdps(&_IStabilization.CallOpts, owner)
}

// CollateralPrice is a free data retrieval call binding the contract method 0x5891de72.
//
// Solidity: function collateralPrice() view returns(uint256 price)
func (_IStabilization *IStabilizationCaller) CollateralPrice(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IStabilization.contract.Call(opts, &out, "collateralPrice")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CollateralPrice is a free data retrieval call binding the contract method 0x5891de72.
//
// Solidity: function collateralPrice() view returns(uint256 price)
func (_IStabilization *IStabilizationSession) CollateralPrice() (*big.Int, error) {
	return _IStabilization.Contract.CollateralPrice(&_IStabilization.CallOpts)
}

// CollateralPrice is a free data retrieval call binding the contract method 0x5891de72.
//
// Solidity: function collateralPrice() view returns(uint256 price)
func (_IStabilization *IStabilizationCallerSession) CollateralPrice() (*big.Int, error) {
	return _IStabilization.Contract.CollateralPrice(&_IStabilization.CallOpts)
}

// Config is a free data retrieval call binding the contract method 0x79502c55.
//
// Solidity: function config() view returns((uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256))
func (_IStabilization *IStabilizationCaller) Config(opts *bind.CallOpts) (IStabilizationConfig, error) {
	var out []interface{}
	err := _IStabilization.contract.Call(opts, &out, "config")

	if err != nil {
		return *new(IStabilizationConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(IStabilizationConfig)).(*IStabilizationConfig)

	return out0, err

}

// Config is a free data retrieval call binding the contract method 0x79502c55.
//
// Solidity: function config() view returns((uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256))
func (_IStabilization *IStabilizationSession) Config() (IStabilizationConfig, error) {
	return _IStabilization.Contract.Config(&_IStabilization.CallOpts)
}

// Config is a free data retrieval call binding the contract method 0x79502c55.
//
// Solidity: function config() view returns((uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256))
func (_IStabilization *IStabilizationCallerSession) Config() (IStabilizationConfig, error) {
	return _IStabilization.Contract.Config(&_IStabilization.CallOpts)
}

// DebtAmountAtTime is a free data retrieval call binding the contract method 0x8ffba9dd.
//
// Solidity: function debtAmountAtTime(address account, uint256 timestamp) view returns(uint256)
func (_IStabilization *IStabilizationCaller) DebtAmountAtTime(opts *bind.CallOpts, account common.Address, timestamp *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _IStabilization.contract.Call(opts, &out, "debtAmountAtTime", account, timestamp)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// DebtAmountAtTime is a free data retrieval call binding the contract method 0x8ffba9dd.
//
// Solidity: function debtAmountAtTime(address account, uint256 timestamp) view returns(uint256)
func (_IStabilization *IStabilizationSession) DebtAmountAtTime(account common.Address, timestamp *big.Int) (*big.Int, error) {
	return _IStabilization.Contract.DebtAmountAtTime(&_IStabilization.CallOpts, account, timestamp)
}

// DebtAmountAtTime is a free data retrieval call binding the contract method 0x8ffba9dd.
//
// Solidity: function debtAmountAtTime(address account, uint256 timestamp) view returns(uint256)
func (_IStabilization *IStabilizationCallerSession) DebtAmountAtTime(account common.Address, timestamp *big.Int) (*big.Int, error) {
	return _IStabilization.Contract.DebtAmountAtTime(&_IStabilization.CallOpts, account, timestamp)
}

// LastUpdated is a free data retrieval call binding the contract method 0xd0b06f5d.
//
// Solidity: function lastUpdated() view returns((uint256,uint256,uint256,uint256))
func (_IStabilization *IStabilizationCaller) LastUpdated(opts *bind.CallOpts) (IStabilizationLastUpdated, error) {
	var out []interface{}
	err := _IStabilization.contract.Call(opts, &out, "lastUpdated")

	if err != nil {
		return *new(IStabilizationLastUpdated), err
	}

	out0 := *abi.ConvertType(out[0], new(IStabilizationLastUpdated)).(*IStabilizationLastUpdated)

	return out0, err

}

// LastUpdated is a free data retrieval call binding the contract method 0xd0b06f5d.
//
// Solidity: function lastUpdated() view returns((uint256,uint256,uint256,uint256))
func (_IStabilization *IStabilizationSession) LastUpdated() (IStabilizationLastUpdated, error) {
	return _IStabilization.Contract.LastUpdated(&_IStabilization.CallOpts)
}

// LastUpdated is a free data retrieval call binding the contract method 0xd0b06f5d.
//
// Solidity: function lastUpdated() view returns((uint256,uint256,uint256,uint256))
func (_IStabilization *IStabilizationCallerSession) LastUpdated() (IStabilizationLastUpdated, error) {
	return _IStabilization.Contract.LastUpdated(&_IStabilization.CallOpts)
}

// Liquidate is a paid mutator transaction binding the contract method 0x4914c008.
//
// Solidity: function liquidate(address account, uint256 collateralSold, address bidder) payable returns()
func (_IStabilization *IStabilizationTransactor) Liquidate(opts *bind.TransactOpts, account common.Address, collateralSold *big.Int, bidder common.Address) (*types.Transaction, error) {
	return _IStabilization.contract.Transact(opts, "liquidate", account, collateralSold, bidder)
}

// Liquidate is a paid mutator transaction binding the contract method 0x4914c008.
//
// Solidity: function liquidate(address account, uint256 collateralSold, address bidder) payable returns()
func (_IStabilization *IStabilizationSession) Liquidate(account common.Address, collateralSold *big.Int, bidder common.Address) (*types.Transaction, error) {
	return _IStabilization.Contract.Liquidate(&_IStabilization.TransactOpts, account, collateralSold, bidder)
}

// Liquidate is a paid mutator transaction binding the contract method 0x4914c008.
//
// Solidity: function liquidate(address account, uint256 collateralSold, address bidder) payable returns()
func (_IStabilization *IStabilizationTransactorSession) Liquidate(account common.Address, collateralSold *big.Int, bidder common.Address) (*types.Transaction, error) {
	return _IStabilization.Contract.Liquidate(&_IStabilization.TransactOpts, account, collateralSold, bidder)
}

// SetACU is a paid mutator transaction binding the contract method 0x4b8ef943.
//
// Solidity: function setACU(address acu) returns()
func (_IStabilization *IStabilizationTransactor) SetACU(opts *bind.TransactOpts, acu common.Address) (*types.Transaction, error) {
	return _IStabilization.contract.Transact(opts, "setACU", acu)
}

// SetACU is a paid mutator transaction binding the contract method 0x4b8ef943.
//
// Solidity: function setACU(address acu) returns()
func (_IStabilization *IStabilizationSession) SetACU(acu common.Address) (*types.Transaction, error) {
	return _IStabilization.Contract.SetACU(&_IStabilization.TransactOpts, acu)
}

// SetACU is a paid mutator transaction binding the contract method 0x4b8ef943.
//
// Solidity: function setACU(address acu) returns()
func (_IStabilization *IStabilizationTransactorSession) SetACU(acu common.Address) (*types.Transaction, error) {
	return _IStabilization.Contract.SetACU(&_IStabilization.TransactOpts, acu)
}

// SetAuctioneer is a paid mutator transaction binding the contract method 0x00ede7e4.
//
// Solidity: function setAuctioneer(address auctioneer) returns()
func (_IStabilization *IStabilizationTransactor) SetAuctioneer(opts *bind.TransactOpts, auctioneer common.Address) (*types.Transaction, error) {
	return _IStabilization.contract.Transact(opts, "setAuctioneer", auctioneer)
}

// SetAuctioneer is a paid mutator transaction binding the contract method 0x00ede7e4.
//
// Solidity: function setAuctioneer(address auctioneer) returns()
func (_IStabilization *IStabilizationSession) SetAuctioneer(auctioneer common.Address) (*types.Transaction, error) {
	return _IStabilization.Contract.SetAuctioneer(&_IStabilization.TransactOpts, auctioneer)
}

// SetAuctioneer is a paid mutator transaction binding the contract method 0x00ede7e4.
//
// Solidity: function setAuctioneer(address auctioneer) returns()
func (_IStabilization *IStabilizationTransactorSession) SetAuctioneer(auctioneer common.Address) (*types.Transaction, error) {
	return _IStabilization.Contract.SetAuctioneer(&_IStabilization.TransactOpts, auctioneer)
}

// SetOperator is a paid mutator transaction binding the contract method 0xb3ab15fb.
//
// Solidity: function setOperator(address operator) returns()
func (_IStabilization *IStabilizationTransactor) SetOperator(opts *bind.TransactOpts, operator common.Address) (*types.Transaction, error) {
	return _IStabilization.contract.Transact(opts, "setOperator", operator)
}

// SetOperator is a paid mutator transaction binding the contract method 0xb3ab15fb.
//
// Solidity: function setOperator(address operator) returns()
func (_IStabilization *IStabilizationSession) SetOperator(operator common.Address) (*types.Transaction, error) {
	return _IStabilization.Contract.SetOperator(&_IStabilization.TransactOpts, operator)
}

// SetOperator is a paid mutator transaction binding the contract method 0xb3ab15fb.
//
// Solidity: function setOperator(address operator) returns()
func (_IStabilization *IStabilizationTransactorSession) SetOperator(operator common.Address) (*types.Transaction, error) {
	return _IStabilization.Contract.SetOperator(&_IStabilization.TransactOpts, operator)
}

// SetOracle is a paid mutator transaction binding the contract method 0x7adbf973.
//
// Solidity: function setOracle(address oracle) returns()
func (_IStabilization *IStabilizationTransactor) SetOracle(opts *bind.TransactOpts, oracle common.Address) (*types.Transaction, error) {
	return _IStabilization.contract.Transact(opts, "setOracle", oracle)
}

// SetOracle is a paid mutator transaction binding the contract method 0x7adbf973.
//
// Solidity: function setOracle(address oracle) returns()
func (_IStabilization *IStabilizationSession) SetOracle(oracle common.Address) (*types.Transaction, error) {
	return _IStabilization.Contract.SetOracle(&_IStabilization.TransactOpts, oracle)
}

// SetOracle is a paid mutator transaction binding the contract method 0x7adbf973.
//
// Solidity: function setOracle(address oracle) returns()
func (_IStabilization *IStabilizationTransactorSession) SetOracle(oracle common.Address) (*types.Transaction, error) {
	return _IStabilization.Contract.SetOracle(&_IStabilization.TransactOpts, oracle)
}

// SetSupplyControl is a paid mutator transaction binding the contract method 0x52e5a050.
//
// Solidity: function setSupplyControl(address supplyControl) returns()
func (_IStabilization *IStabilizationTransactor) SetSupplyControl(opts *bind.TransactOpts, supplyControl common.Address) (*types.Transaction, error) {
	return _IStabilization.contract.Transact(opts, "setSupplyControl", supplyControl)
}

// SetSupplyControl is a paid mutator transaction binding the contract method 0x52e5a050.
//
// Solidity: function setSupplyControl(address supplyControl) returns()
func (_IStabilization *IStabilizationSession) SetSupplyControl(supplyControl common.Address) (*types.Transaction, error) {
	return _IStabilization.Contract.SetSupplyControl(&_IStabilization.TransactOpts, supplyControl)
}

// SetSupplyControl is a paid mutator transaction binding the contract method 0x52e5a050.
//
// Solidity: function setSupplyControl(address supplyControl) returns()
func (_IStabilization *IStabilizationTransactorSession) SetSupplyControl(supplyControl common.Address) (*types.Transaction, error) {
	return _IStabilization.Contract.SetSupplyControl(&_IStabilization.TransactOpts, supplyControl)
}

// ISupplyControlMetaData contains all meta data concerning the ISupplyControl contract.
var ISupplyControlMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Burn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Mint\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"availableSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"burn\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getStabilizer\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTotalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"mint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"setOperator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"stabilizer_\",\"type\":\"address\"}],\"name\":\"setStabilizer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"7ecc2b56": "availableSupply()",
		"44df8e70": "burn()",
		"80af1799": "getStabilizer()",
		"c4e41b22": "getTotalSupply()",
		"40c10f19": "mint(address,uint256)",
		"b3ab15fb": "setOperator(address)",
		"db7f521a": "setStabilizer(address)",
	},
}

// ISupplyControlABI is the input ABI used to generate the binding from.
// Deprecated: Use ISupplyControlMetaData.ABI instead.
var ISupplyControlABI = ISupplyControlMetaData.ABI

// Deprecated: Use ISupplyControlMetaData.Sigs instead.
// ISupplyControlFuncSigs maps the 4-byte function signature to its string representation.
var ISupplyControlFuncSigs = ISupplyControlMetaData.Sigs

// ISupplyControl is an auto generated Go binding around an Ethereum contract.
type ISupplyControl struct {
	ISupplyControlCaller     // Read-only binding to the contract
	ISupplyControlTransactor // Write-only binding to the contract
	ISupplyControlFilterer   // Log filterer for contract events
}

// ISupplyControlCaller is an auto generated read-only Go binding around an Ethereum contract.
type ISupplyControlCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ISupplyControlTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ISupplyControlTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ISupplyControlFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ISupplyControlFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ISupplyControlSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ISupplyControlSession struct {
	Contract     *ISupplyControl   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ISupplyControlCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ISupplyControlCallerSession struct {
	Contract *ISupplyControlCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// ISupplyControlTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ISupplyControlTransactorSession struct {
	Contract     *ISupplyControlTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// ISupplyControlRaw is an auto generated low-level Go binding around an Ethereum contract.
type ISupplyControlRaw struct {
	Contract *ISupplyControl // Generic contract binding to access the raw methods on
}

// ISupplyControlCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ISupplyControlCallerRaw struct {
	Contract *ISupplyControlCaller // Generic read-only contract binding to access the raw methods on
}

// ISupplyControlTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ISupplyControlTransactorRaw struct {
	Contract *ISupplyControlTransactor // Generic write-only contract binding to access the raw methods on
}

// NewISupplyControl creates a new instance of ISupplyControl, bound to a specific deployed contract.
func NewISupplyControl(address common.Address, backend bind.ContractBackend) (*ISupplyControl, error) {
	contract, err := bindISupplyControl(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ISupplyControl{ISupplyControlCaller: ISupplyControlCaller{contract: contract}, ISupplyControlTransactor: ISupplyControlTransactor{contract: contract}, ISupplyControlFilterer: ISupplyControlFilterer{contract: contract}}, nil
}

// NewISupplyControlCaller creates a new read-only instance of ISupplyControl, bound to a specific deployed contract.
func NewISupplyControlCaller(address common.Address, caller bind.ContractCaller) (*ISupplyControlCaller, error) {
	contract, err := bindISupplyControl(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ISupplyControlCaller{contract: contract}, nil
}

// NewISupplyControlTransactor creates a new write-only instance of ISupplyControl, bound to a specific deployed contract.
func NewISupplyControlTransactor(address common.Address, transactor bind.ContractTransactor) (*ISupplyControlTransactor, error) {
	contract, err := bindISupplyControl(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ISupplyControlTransactor{contract: contract}, nil
}

// NewISupplyControlFilterer creates a new log filterer instance of ISupplyControl, bound to a specific deployed contract.
func NewISupplyControlFilterer(address common.Address, filterer bind.ContractFilterer) (*ISupplyControlFilterer, error) {
	contract, err := bindISupplyControl(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ISupplyControlFilterer{contract: contract}, nil
}

// bindISupplyControl binds a generic wrapper to an already deployed contract.
func bindISupplyControl(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ISupplyControlABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ISupplyControl *ISupplyControlRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ISupplyControl.Contract.ISupplyControlCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ISupplyControl *ISupplyControlRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ISupplyControl.Contract.ISupplyControlTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ISupplyControl *ISupplyControlRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ISupplyControl.Contract.ISupplyControlTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ISupplyControl *ISupplyControlCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ISupplyControl.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ISupplyControl *ISupplyControlTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ISupplyControl.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ISupplyControl *ISupplyControlTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ISupplyControl.Contract.contract.Transact(opts, method, params...)
}

// AvailableSupply is a free data retrieval call binding the contract method 0x7ecc2b56.
//
// Solidity: function availableSupply() view returns(uint256)
func (_ISupplyControl *ISupplyControlCaller) AvailableSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ISupplyControl.contract.Call(opts, &out, "availableSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// AvailableSupply is a free data retrieval call binding the contract method 0x7ecc2b56.
//
// Solidity: function availableSupply() view returns(uint256)
func (_ISupplyControl *ISupplyControlSession) AvailableSupply() (*big.Int, error) {
	return _ISupplyControl.Contract.AvailableSupply(&_ISupplyControl.CallOpts)
}

// AvailableSupply is a free data retrieval call binding the contract method 0x7ecc2b56.
//
// Solidity: function availableSupply() view returns(uint256)
func (_ISupplyControl *ISupplyControlCallerSession) AvailableSupply() (*big.Int, error) {
	return _ISupplyControl.Contract.AvailableSupply(&_ISupplyControl.CallOpts)
}

// GetStabilizer is a free data retrieval call binding the contract method 0x80af1799.
//
// Solidity: function getStabilizer() view returns(address)
func (_ISupplyControl *ISupplyControlCaller) GetStabilizer(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ISupplyControl.contract.Call(opts, &out, "getStabilizer")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetStabilizer is a free data retrieval call binding the contract method 0x80af1799.
//
// Solidity: function getStabilizer() view returns(address)
func (_ISupplyControl *ISupplyControlSession) GetStabilizer() (common.Address, error) {
	return _ISupplyControl.Contract.GetStabilizer(&_ISupplyControl.CallOpts)
}

// GetStabilizer is a free data retrieval call binding the contract method 0x80af1799.
//
// Solidity: function getStabilizer() view returns(address)
func (_ISupplyControl *ISupplyControlCallerSession) GetStabilizer() (common.Address, error) {
	return _ISupplyControl.Contract.GetStabilizer(&_ISupplyControl.CallOpts)
}

// GetTotalSupply is a free data retrieval call binding the contract method 0xc4e41b22.
//
// Solidity: function getTotalSupply() view returns(uint256)
func (_ISupplyControl *ISupplyControlCaller) GetTotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ISupplyControl.contract.Call(opts, &out, "getTotalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetTotalSupply is a free data retrieval call binding the contract method 0xc4e41b22.
//
// Solidity: function getTotalSupply() view returns(uint256)
func (_ISupplyControl *ISupplyControlSession) GetTotalSupply() (*big.Int, error) {
	return _ISupplyControl.Contract.GetTotalSupply(&_ISupplyControl.CallOpts)
}

// GetTotalSupply is a free data retrieval call binding the contract method 0xc4e41b22.
//
// Solidity: function getTotalSupply() view returns(uint256)
func (_ISupplyControl *ISupplyControlCallerSession) GetTotalSupply() (*big.Int, error) {
	return _ISupplyControl.Contract.GetTotalSupply(&_ISupplyControl.CallOpts)
}

// Burn is a paid mutator transaction binding the contract method 0x44df8e70.
//
// Solidity: function burn() payable returns()
func (_ISupplyControl *ISupplyControlTransactor) Burn(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ISupplyControl.contract.Transact(opts, "burn")
}

// Burn is a paid mutator transaction binding the contract method 0x44df8e70.
//
// Solidity: function burn() payable returns()
func (_ISupplyControl *ISupplyControlSession) Burn() (*types.Transaction, error) {
	return _ISupplyControl.Contract.Burn(&_ISupplyControl.TransactOpts)
}

// Burn is a paid mutator transaction binding the contract method 0x44df8e70.
//
// Solidity: function burn() payable returns()
func (_ISupplyControl *ISupplyControlTransactorSession) Burn() (*types.Transaction, error) {
	return _ISupplyControl.Contract.Burn(&_ISupplyControl.TransactOpts)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address recipient, uint256 amount) returns()
func (_ISupplyControl *ISupplyControlTransactor) Mint(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ISupplyControl.contract.Transact(opts, "mint", recipient, amount)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address recipient, uint256 amount) returns()
func (_ISupplyControl *ISupplyControlSession) Mint(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ISupplyControl.Contract.Mint(&_ISupplyControl.TransactOpts, recipient, amount)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address recipient, uint256 amount) returns()
func (_ISupplyControl *ISupplyControlTransactorSession) Mint(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ISupplyControl.Contract.Mint(&_ISupplyControl.TransactOpts, recipient, amount)
}

// SetOperator is a paid mutator transaction binding the contract method 0xb3ab15fb.
//
// Solidity: function setOperator(address operator) returns()
func (_ISupplyControl *ISupplyControlTransactor) SetOperator(opts *bind.TransactOpts, operator common.Address) (*types.Transaction, error) {
	return _ISupplyControl.contract.Transact(opts, "setOperator", operator)
}

// SetOperator is a paid mutator transaction binding the contract method 0xb3ab15fb.
//
// Solidity: function setOperator(address operator) returns()
func (_ISupplyControl *ISupplyControlSession) SetOperator(operator common.Address) (*types.Transaction, error) {
	return _ISupplyControl.Contract.SetOperator(&_ISupplyControl.TransactOpts, operator)
}

// SetOperator is a paid mutator transaction binding the contract method 0xb3ab15fb.
//
// Solidity: function setOperator(address operator) returns()
func (_ISupplyControl *ISupplyControlTransactorSession) SetOperator(operator common.Address) (*types.Transaction, error) {
	return _ISupplyControl.Contract.SetOperator(&_ISupplyControl.TransactOpts, operator)
}

// SetStabilizer is a paid mutator transaction binding the contract method 0xdb7f521a.
//
// Solidity: function setStabilizer(address stabilizer_) returns()
func (_ISupplyControl *ISupplyControlTransactor) SetStabilizer(opts *bind.TransactOpts, stabilizer_ common.Address) (*types.Transaction, error) {
	return _ISupplyControl.contract.Transact(opts, "setStabilizer", stabilizer_)
}

// SetStabilizer is a paid mutator transaction binding the contract method 0xdb7f521a.
//
// Solidity: function setStabilizer(address stabilizer_) returns()
func (_ISupplyControl *ISupplyControlSession) SetStabilizer(stabilizer_ common.Address) (*types.Transaction, error) {
	return _ISupplyControl.Contract.SetStabilizer(&_ISupplyControl.TransactOpts, stabilizer_)
}

// SetStabilizer is a paid mutator transaction binding the contract method 0xdb7f521a.
//
// Solidity: function setStabilizer(address stabilizer_) returns()
func (_ISupplyControl *ISupplyControlTransactorSession) SetStabilizer(stabilizer_ common.Address) (*types.Transaction, error) {
	return _ISupplyControl.Contract.SetStabilizer(&_ISupplyControl.TransactOpts, stabilizer_)
}

// ISupplyControlBurnIterator is returned from FilterBurn and is used to iterate over the raw logs and unpacked data for Burn events raised by the ISupplyControl contract.
type ISupplyControlBurnIterator struct {
	Event *ISupplyControlBurn // Event containing the contract specifics and raw log

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
func (it *ISupplyControlBurnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ISupplyControlBurn)
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
		it.Event = new(ISupplyControlBurn)
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
func (it *ISupplyControlBurnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ISupplyControlBurnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ISupplyControlBurn represents a Burn event raised by the ISupplyControl contract.
type ISupplyControlBurn struct {
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterBurn is a free log retrieval operation binding the contract event 0xb90306ad06b2a6ff86ddc9327db583062895ef6540e62dc50add009db5b356eb.
//
// Solidity: event Burn(uint256 amount)
func (_ISupplyControl *ISupplyControlFilterer) FilterBurn(opts *bind.FilterOpts) (*ISupplyControlBurnIterator, error) {

	logs, sub, err := _ISupplyControl.contract.FilterLogs(opts, "Burn")
	if err != nil {
		return nil, err
	}
	return &ISupplyControlBurnIterator{contract: _ISupplyControl.contract, event: "Burn", logs: logs, sub: sub}, nil
}

// WatchBurn is a free log subscription operation binding the contract event 0xb90306ad06b2a6ff86ddc9327db583062895ef6540e62dc50add009db5b356eb.
//
// Solidity: event Burn(uint256 amount)
func (_ISupplyControl *ISupplyControlFilterer) WatchBurn(opts *bind.WatchOpts, sink chan<- *ISupplyControlBurn) (event.Subscription, error) {

	logs, sub, err := _ISupplyControl.contract.WatchLogs(opts, "Burn")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ISupplyControlBurn)
				if err := _ISupplyControl.contract.UnpackLog(event, "Burn", log); err != nil {
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

// ParseBurn is a log parse operation binding the contract event 0xb90306ad06b2a6ff86ddc9327db583062895ef6540e62dc50add009db5b356eb.
//
// Solidity: event Burn(uint256 amount)
func (_ISupplyControl *ISupplyControlFilterer) ParseBurn(log types.Log) (*ISupplyControlBurn, error) {
	event := new(ISupplyControlBurn)
	if err := _ISupplyControl.contract.UnpackLog(event, "Burn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ISupplyControlMintIterator is returned from FilterMint and is used to iterate over the raw logs and unpacked data for Mint events raised by the ISupplyControl contract.
type ISupplyControlMintIterator struct {
	Event *ISupplyControlMint // Event containing the contract specifics and raw log

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
func (it *ISupplyControlMintIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ISupplyControlMint)
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
		it.Event = new(ISupplyControlMint)
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
func (it *ISupplyControlMintIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ISupplyControlMintIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ISupplyControlMint represents a Mint event raised by the ISupplyControl contract.
type ISupplyControlMint struct {
	Recipient common.Address
	Amount    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterMint is a free log retrieval operation binding the contract event 0x0f6798a560793a54c3bcfe86a93cde1e73087d944c0ea20544137d4121396885.
//
// Solidity: event Mint(address recipient, uint256 amount)
func (_ISupplyControl *ISupplyControlFilterer) FilterMint(opts *bind.FilterOpts) (*ISupplyControlMintIterator, error) {

	logs, sub, err := _ISupplyControl.contract.FilterLogs(opts, "Mint")
	if err != nil {
		return nil, err
	}
	return &ISupplyControlMintIterator{contract: _ISupplyControl.contract, event: "Mint", logs: logs, sub: sub}, nil
}

// WatchMint is a free log subscription operation binding the contract event 0x0f6798a560793a54c3bcfe86a93cde1e73087d944c0ea20544137d4121396885.
//
// Solidity: event Mint(address recipient, uint256 amount)
func (_ISupplyControl *ISupplyControlFilterer) WatchMint(opts *bind.WatchOpts, sink chan<- *ISupplyControlMint) (event.Subscription, error) {

	logs, sub, err := _ISupplyControl.contract.WatchLogs(opts, "Mint")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ISupplyControlMint)
				if err := _ISupplyControl.contract.UnpackLog(event, "Mint", log); err != nil {
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

// ParseMint is a log parse operation binding the contract event 0x0f6798a560793a54c3bcfe86a93cde1e73087d944c0ea20544137d4121396885.
//
// Solidity: event Mint(address recipient, uint256 amount)
func (_ISupplyControl *ISupplyControlFilterer) ParseMint(log types.Log) (*ISupplyControlMint, error) {
	event := new(ISupplyControlMint)
	if err := _ISupplyControl.contract.UnpackLog(event, "Mint", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IUpgradeManagerMetaData contains all meta data concerning the IUpgradeManager contract.
var IUpgradeManagerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_account\",\"type\":\"address\"}],\"name\":\"setOperator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"b3ab15fb": "setOperator(address)",
	},
}

// IUpgradeManagerABI is the input ABI used to generate the binding from.
// Deprecated: Use IUpgradeManagerMetaData.ABI instead.
var IUpgradeManagerABI = IUpgradeManagerMetaData.ABI

// Deprecated: Use IUpgradeManagerMetaData.Sigs instead.
// IUpgradeManagerFuncSigs maps the 4-byte function signature to its string representation.
var IUpgradeManagerFuncSigs = IUpgradeManagerMetaData.Sigs

// IUpgradeManager is an auto generated Go binding around an Ethereum contract.
type IUpgradeManager struct {
	IUpgradeManagerCaller     // Read-only binding to the contract
	IUpgradeManagerTransactor // Write-only binding to the contract
	IUpgradeManagerFilterer   // Log filterer for contract events
}

// IUpgradeManagerCaller is an auto generated read-only Go binding around an Ethereum contract.
type IUpgradeManagerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IUpgradeManagerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type IUpgradeManagerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IUpgradeManagerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type IUpgradeManagerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IUpgradeManagerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type IUpgradeManagerSession struct {
	Contract     *IUpgradeManager  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IUpgradeManagerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type IUpgradeManagerCallerSession struct {
	Contract *IUpgradeManagerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// IUpgradeManagerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type IUpgradeManagerTransactorSession struct {
	Contract     *IUpgradeManagerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// IUpgradeManagerRaw is an auto generated low-level Go binding around an Ethereum contract.
type IUpgradeManagerRaw struct {
	Contract *IUpgradeManager // Generic contract binding to access the raw methods on
}

// IUpgradeManagerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type IUpgradeManagerCallerRaw struct {
	Contract *IUpgradeManagerCaller // Generic read-only contract binding to access the raw methods on
}

// IUpgradeManagerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type IUpgradeManagerTransactorRaw struct {
	Contract *IUpgradeManagerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIUpgradeManager creates a new instance of IUpgradeManager, bound to a specific deployed contract.
func NewIUpgradeManager(address common.Address, backend bind.ContractBackend) (*IUpgradeManager, error) {
	contract, err := bindIUpgradeManager(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IUpgradeManager{IUpgradeManagerCaller: IUpgradeManagerCaller{contract: contract}, IUpgradeManagerTransactor: IUpgradeManagerTransactor{contract: contract}, IUpgradeManagerFilterer: IUpgradeManagerFilterer{contract: contract}}, nil
}

// NewIUpgradeManagerCaller creates a new read-only instance of IUpgradeManager, bound to a specific deployed contract.
func NewIUpgradeManagerCaller(address common.Address, caller bind.ContractCaller) (*IUpgradeManagerCaller, error) {
	contract, err := bindIUpgradeManager(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IUpgradeManagerCaller{contract: contract}, nil
}

// NewIUpgradeManagerTransactor creates a new write-only instance of IUpgradeManager, bound to a specific deployed contract.
func NewIUpgradeManagerTransactor(address common.Address, transactor bind.ContractTransactor) (*IUpgradeManagerTransactor, error) {
	contract, err := bindIUpgradeManager(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IUpgradeManagerTransactor{contract: contract}, nil
}

// NewIUpgradeManagerFilterer creates a new log filterer instance of IUpgradeManager, bound to a specific deployed contract.
func NewIUpgradeManagerFilterer(address common.Address, filterer bind.ContractFilterer) (*IUpgradeManagerFilterer, error) {
	contract, err := bindIUpgradeManager(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IUpgradeManagerFilterer{contract: contract}, nil
}

// bindIUpgradeManager binds a generic wrapper to an already deployed contract.
func bindIUpgradeManager(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(IUpgradeManagerABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IUpgradeManager *IUpgradeManagerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IUpgradeManager.Contract.IUpgradeManagerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IUpgradeManager *IUpgradeManagerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IUpgradeManager.Contract.IUpgradeManagerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IUpgradeManager *IUpgradeManagerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IUpgradeManager.Contract.IUpgradeManagerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IUpgradeManager *IUpgradeManagerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IUpgradeManager.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IUpgradeManager *IUpgradeManagerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IUpgradeManager.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IUpgradeManager *IUpgradeManagerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IUpgradeManager.Contract.contract.Transact(opts, method, params...)
}

// SetOperator is a paid mutator transaction binding the contract method 0xb3ab15fb.
//
// Solidity: function setOperator(address _account) returns()
func (_IUpgradeManager *IUpgradeManagerTransactor) SetOperator(opts *bind.TransactOpts, _account common.Address) (*types.Transaction, error) {
	return _IUpgradeManager.contract.Transact(opts, "setOperator", _account)
}

// SetOperator is a paid mutator transaction binding the contract method 0xb3ab15fb.
//
// Solidity: function setOperator(address _account) returns()
func (_IUpgradeManager *IUpgradeManagerSession) SetOperator(_account common.Address) (*types.Transaction, error) {
	return _IUpgradeManager.Contract.SetOperator(&_IUpgradeManager.TransactOpts, _account)
}

// SetOperator is a paid mutator transaction binding the contract method 0xb3ab15fb.
//
// Solidity: function setOperator(address _account) returns()
func (_IUpgradeManager *IUpgradeManagerTransactorSession) SetOperator(_account common.Address) (*types.Transaction, error) {
	return _IUpgradeManager.Contract.SetOperator(&_IUpgradeManager.TransactOpts, _account)
}

// Oracle0MetaData contains all meta data concerning the Oracle0 contract.
var Oracle0MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"methodSignature\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"returnData\",\"type\":\"bytes\"}],\"name\":\"CallFailed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_voter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_round\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_nonRevealCount\",\"type\":\"uint256\"}],\"name\":\"CommitRevealMissed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"oldValue\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newValue\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"appliesAtHeight\",\"type\":\"uint256\"}],\"name\":\"ConfigUpdateAddress\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"oldValue\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"newValue\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"appliesAtHeight\",\"type\":\"uint256\"}],\"name\":\"ConfigUpdateBool\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"int256\",\"name\":\"oldValue\",\"type\":\"int256\"},{\"indexed\":false,\"internalType\":\"int256\",\"name\":\"newValue\",\"type\":\"int256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"appliesAtHeight\",\"type\":\"uint256\"}],\"name\":\"ConfigUpdateInt\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldValue\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newValue\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"appliesAtHeight\",\"type\":\"uint256\"}],\"name\":\"ConfigUpdateUint\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"cause\",\"type\":\"string\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"reporter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"expValue\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"actualValue\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"extra\",\"type\":\"uint8\"}],\"name\":\"InvalidVote\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_round\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_timestamp\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_votePeriod\",\"type\":\"uint256\"}],\"name\":\"NewRound\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string[]\",\"name\":\"_symbols\",\"type\":\"string[]\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_round\",\"type\":\"uint256\"}],\"name\":\"NewSymbols\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"reporter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"extra\",\"type\":\"uint8\"}],\"name\":\"NewVoter\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_voter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_round\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_missedReveal\",\"type\":\"uint256\"}],\"name\":\"NoRevealPenalty\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_participant\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_slashingAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"_symbol\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"int256\",\"name\":\"_median\",\"type\":\"int256\"},{\"indexed\":false,\"internalType\":\"uint120\",\"name\":\"_reported\",\"type\":\"uint120\"}],\"name\":\"Penalized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"price\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"round\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"status\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"}],\"name\":\"PriceUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"reporter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"extra\",\"type\":\"uint8\"}],\"name\":\"SuccessfulVote\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"ntnReward\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"atnReward\",\"type\":\"uint256\"}],\"name\":\"TotalOracleRewards\",\"type\":\"event\"},{\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_ntn\",\"type\":\"uint256\"}],\"name\":\"distributeRewards\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"finalize\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getConfig\",\"outputs\":[{\"components\":[{\"internalType\":\"contractIAutonity\",\"name\":\"autonity\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"votePeriod\",\"type\":\"uint256\"},{\"internalType\":\"int256\",\"name\":\"outlierDetectionThreshold\",\"type\":\"int256\"},{\"internalType\":\"int256\",\"name\":\"outlierSlashingThreshold\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"baseSlashingRate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"nonRevealThreshold\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"revealResetInterval\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"slashingRateCap\",\"type\":\"uint256\"}],\"internalType\":\"structOracle0.Config\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getDecimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLastRoundBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getNewVotePeriod\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getNewVoters\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getNonRevealThreshold\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_symbol\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"_voter\",\"type\":\"address\"}],\"name\":\"getReports\",\"outputs\":[{\"components\":[{\"internalType\":\"uint120\",\"name\":\"price\",\"type\":\"uint120\"},{\"internalType\":\"uint8\",\"name\":\"confidence\",\"type\":\"uint8\"}],\"internalType\":\"structIOracle.Report\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_voter\",\"type\":\"address\"}],\"name\":\"getRewardPeriodPerformance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRound\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_round\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"_symbol\",\"type\":\"string\"}],\"name\":\"getRoundData\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"round\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"price\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"internalType\":\"structIOracle.RoundData\",\"name\":\"data\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getSymbolUpdatedRound\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getSymbols\",\"outputs\":[{\"internalType\":\"string[]\",\"name\":\"\",\"type\":\"string[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getVotePeriod\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_voter\",\"type\":\"address\"}],\"name\":\"getVoterInfo\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"round\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"commit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"performance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"nonRevealCount\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"isVoter\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"reportAvailable\",\"type\":\"bool\"}],\"internalType\":\"structOracle0.VoterInfo\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_oracleAddress\",\"type\":\"address\"}],\"name\":\"getVoterTreasuries\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_oracleAddress\",\"type\":\"address\"}],\"name\":\"getVoterValidators\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getVoters\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_symbol\",\"type\":\"string\"}],\"name\":\"latestRoundData\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"round\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"price\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"internalType\":\"structIOracle.RoundData\",\"name\":\"data\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_threshold\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_resetInterval\",\"type\":\"uint256\"}],\"name\":\"setCommitRevealConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_operator\",\"type\":\"address\"}],\"name\":\"setOperator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"_outlierSlashingThreshold\",\"type\":\"int256\"},{\"internalType\":\"int256\",\"name\":\"_outlierDetectionThreshold\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"_baseSlashingRate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_slashingRateCap\",\"type\":\"uint256\"}],\"name\":\"setSlashingConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string[]\",\"name\":\"_symbols\",\"type\":\"string[]\"}],\"name\":\"setSymbols\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_votePeriod\",\"type\":\"uint256\"}],\"name\":\"setVotePeriod\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_newVoters\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"_treasury\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"_validator\",\"type\":\"address[]\"}],\"name\":\"setVoters\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"updateVotersAndSymbol\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_commit\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"uint120\",\"name\":\"price\",\"type\":\"uint120\"},{\"internalType\":\"uint8\",\"name\":\"confidence\",\"type\":\"uint8\"}],\"internalType\":\"structIOracle.Report[]\",\"name\":\"_reports\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"_salt\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"_extra\",\"type\":\"uint8\"}],\"name\":\"vote\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Sigs: map[string]string{
		"59974e38": "distributeRewards(uint256)",
		"4bb278f3": "finalize()",
		"c3f909d4": "getConfig()",
		"f0141d84": "getDecimals()",
		"5a4d3a27": "getLastRoundBlock()",
		"57eba759": "getNewVotePeriod()",
		"077945d3": "getNewVoters()",
		"ed78349d": "getNonRevealThreshold()",
		"fb09917e": "getReports(string,address)",
		"33d16293": "getRewardPeriodPerformance(address)",
		"9f8743f7": "getRound()",
		"3c8510fd": "getRoundData(uint256,string)",
		"99b0014b": "getSymbolUpdatedRound()",
		"df7f710e": "getSymbols()",
		"b78dec52": "getVotePeriod()",
		"9ed1f255": "getVoterInfo(address)",
		"ef5cc4d1": "getVoterTreasuries(address)",
		"2d35d158": "getVoterValidators(address)",
		"cdd72253": "getVoters()",
		"33f98c77": "latestRoundData(string)",
		"3f422ef3": "setCommitRevealConfig(uint256,uint256)",
		"b3ab15fb": "setOperator(address)",
		"da39fbfe": "setSlashingConfig(int256,int256,uint256,uint256)",
		"8d4f75d2": "setSymbols(string[])",
		"67b11630": "setVotePeriod(uint256)",
		"da78110e": "setVoters(address[],address[],address[])",
		"0f65875c": "updateVotersAndSymbol()",
		"56833ebe": "vote(uint256,(uint120,uint8)[],uint256,uint8)",
	},
	Bin: "0x6080604052600160ff1b600a55348015601757600080fd5b50615351806100276000396000f3fe6080604052600436106101a25760003560e01c806399b0014b116100e0578063da39fbfe11610084578063ed78349d11610061578063ed78349d14610509578063ef5cc4d11461051e578063f0141d841461053e578063fb09917e1461055a57005b8063da39fbfe146104a7578063da78110e146104c7578063df7f710e146104e757005b8063b3ab15fb116100bd578063b3ab15fb1461043b578063b78dec521461045b578063c3f909d414610470578063cdd722531461049257005b806399b0014b146103a15780639ed1f255146103b65780639f8743f71461042657005b80634bb278f31161014757806359974e381161012457806359974e38146103395780635a4d3a271461034c57806367b11630146103615780638d4f75d21461038157005b80634bb278f3146102df57806356833ebe1461030457806357eba7591461032457005b806333d162931161018057806333d162931461021c57806333f98c771461024a5780633c8510fd1461029f5780633f422ef3146102bf57005b8063077945d3146101a45780630f65875c146101cf5780632d35d158146101e4575b005b3480156101b057600080fd5b506101b96105a9565b6040516101c691906144d2565b60405180910390f35b3480156101db57600080fd5b506101a261066a565b3480156101f057600080fd5b506102046101ff366004614535565b610888565b6040516001600160a01b0390911681526020016101c6565b34801561022857600080fd5b5061023c610237366004614535565b6108fd565b6040519081526020016101c6565b34801561025657600080fd5b5061026a610265366004614620565b61096c565b6040516101c6919081518152602080830151908201526040808301519082015260609182015115159181019190915260800190565b3480156102ab57600080fd5b5061026a6102ba366004614655565b610aa8565b3480156102cb57600080fd5b506101a26102da36600461469c565b610bbb565b3480156102eb57600080fd5b506102f4610d6b565b60405190151581526020016101c6565b34801561031057600080fd5b506101a261031f3660046146d8565b610f0a565b34801561033057600080fd5b5061023c611454565b6101a2610347366004614774565b6114ae565b34801561035857600080fd5b5061023c61154b565b34801561036d57600080fd5b506101a261037c366004614774565b6115a5565b34801561038d57600080fd5b506101a261039c3660046147b1565b61169a565b3480156103ad57600080fd5b50600a5461023c565b3480156103c257600080fd5b506103d66103d1366004614535565b61180d565b6040516101c69190600060c0820190508251825260208301516020830152604083015160408301526060830151606083015260808301511515608083015260a0830151151560a083015292915050565b34801561043257600080fd5b5061023c611909565b34801561044757600080fd5b506101a2610456366004614535565b611963565b34801561046757600080fd5b5061023c611a1d565b34801561047c57600080fd5b50610485611a77565b6040516101c6919061486d565b34801561049e57600080fd5b506101b9611b8d565b3480156104b357600080fd5b506101a26104c23660046148e9565b611c46565b3480156104d357600080fd5b506101a26104e236600461498a565b611e97565b3480156104f357600080fd5b506104fc6120cd565b6040516101c69190614aca565b34801561051557600080fd5b5060075461023c565b34801561052a57600080fd5b50610204610539366004614535565b61228b565b34801561054a57600080fd5b50604051601281526020016101c6565b34801561056657600080fd5b5061057a610575366004614add565b6122fd565b6040805182516effffffffffffffffffffffffffffff16815260209283015160ff1692810192909252016101c6565b60606105b760005460011490565b156106095760405162461bcd60e51b815260206004820152601d60248201527f72656164206f6e6c79207265656e7472616e637920646574656374656400000060448201526064015b60405180910390fd5b601280548060200260200160405190810160405280929190818152602001828054801561065f57602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311610641575b505050505090505b90565b6001546001600160a01b031633146106ea5760405162461bcd60e51b815260206004820152602360248201527f7265737472696374656420746f20746865206175746f6e69747920636f6e747260448201527f61637400000000000000000000000000000000000000000000000000000000006064820152608401610600565b60135460ff1615156001036107905760005b60125481101561075e576001600c60006012848154811061071f5761071f614b2b565b6000918252602080832091909101546001600160a01b031683528201929092526040019020600401805460ff19169115159190911790556001016106fc565b50601380547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000166101001790556107d5565b601354610100900460ff1615156001036107d5576107ac6123d7565b601380547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff1690555b600a546107e3906002614b89565b6014540361088657600f80546107fb91600e9161430f565b5060005b601154811015610884576000600c60006011848154811061082257610822614b2b565b60009182526020808320909101546001600160a01b0316835282019290925260400190206004018054911515610100027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff9092169190911790556001016107ff565b505b565b600080546001036108db5760405162461bcd60e51b815260206004820152601d60248201527f72656164206f6e6c79207265656e7472616e63792064657465637465640000006044820152606401610600565b506001600160a01b03808216600090815260196020526040902054165b919050565b600080546001036109505760405162461bcd60e51b815260206004820152601d60248201527f72656164206f6e6c79207265656e7472616e63792064657465637465640000006044820152606401610600565b506001600160a01b03166000908152601a602052604090205490565b61099960405180608001604052806000815260200160008152602001600081526020016000151581525090565b6000546001036109eb5760405162461bcd60e51b815260206004820152601d60248201527f72656164206f6e6c79207265656e7472616e63792064657465637465640000006044820152606401610600565b6000601560016014546109fe9190614bb1565b81548110610a0e57610a0e614b2b565b9060005260206000200183604051610a269190614bc4565b90815260408051918290036020908101832060608401835280548452600180820154928501929092526002015460ff161515838301528151608081019092526014549293506000928291610a7991614bb1565b815260200183600001518152602001836020015181526020018360400151151581525090508092505050919050565b610ad560405180608001604052806000815260200160008152602001600081526020016000151581525090565b600054600103610b275760405162461bcd60e51b815260206004820152601d60248201527f72656164206f6e6c79207265656e7472616e63792064657465637465640000006044820152606401610600565b600060158481548110610b3c57610b3c614b2b565b9060005260206000200183604051610b549190614bc4565b9081526040805191829003602090810183206060808501845281548552600182015485840190815260029092015460ff16151585850190815284516080810186528a8152955193860193909352905192840192909252511515908201529150505b92915050565b6002546001600160a01b03163314610c155760405162461bcd60e51b815260206004820152601660248201527f7265737472696374656420746f206f70657261746f72000000000000000000006044820152606401610600565b8082108015610c245750600081115b610c705760405162461bcd60e51b815260206004820152600e60248201527f696e76616c696420636f6e6669670000000000000000000000000000000000006044820152606401610600565b6008546040805160808082526013908201527f72657665616c5265736574496e74657276616c0000000000000000000000000060a0820152602081019290925281018290524360608201527f207e45ce6f2191c3efffe5d0d91591cf81afd23cbf4bd746981d7a8a8dcfe1ba9060c00160405180910390a160088190556007546040805160808082526012908201527f6e6f6e52657665616c5468726573686f6c64000000000000000000000000000060a0820152602081019290925281018390524360608201527f207e45ce6f2191c3efffe5d0d91591cf81afd23cbf4bd746981d7a8a8dcfe1ba9060c00160405180910390a150600755565b6001546000906001600160a01b03163314610dee5760405162461bcd60e51b815260206004820152602360248201527f7265737472696374656420746f20746865206175746f6e69747920636f6e747260448201527f61637400000000000000000000000000000000000000000000000000000000006064820152608401610600565b610df6612607565b600354600b54610e069190614be0565b431015610e1557506000610f01565b6000610e1f61265e565b600854601454919250610e3191614c22565b600003610e4057610e4061285d565b60158054600101815560009081525b600e54811015610e7557610e6381836128b3565b610e6e600182614be0565b9050610e4f565b50610e7e613073565b43600b81905550600160146000828254610e989190614be0565b909155505060105460035414610eaf576010546003555b601454600354604080519283524260208401528201527f5aec57d81928b24d30b1a2aec0d23d693412c37d7ec106b5d8259413716bb1f49060600160405180910390a1610efb81613167565b60019150505b61066760008055565b336000908152600c602052604090206004015460ff16610f6c5760405162461bcd60e51b815260206004820152601960248201527f7265737472696374656420746f206f6e6c7920766f74657273000000000000006044820152606401610600565b610f74612607565b336000908152600c60205260409020601454815403610fd55760405162461bcd60e51b815260206004820152600d60248201527f616c726561647920766f746564000000000000000000000000000000000000006044820152606401610600565b60018101805490879055815460145483556000819003611032576040805133815260ff861660208201527fd2ec8e890a03083998d3e16f98044fd3dd13fe3e61b7bc2e58ee6da43b50af73910160405180910390a1505050611444565b60016014546110419190614bb1565b81146110e057336001600160a01b03167f04ca4e0efda95f8b780c116574d1521309010b38d8f7b75705495703a0f570b160016014546110819190614bb1565b6040805160808082526016908201527f4c617374566f746564526f756e644d69736d617463680000000000000000000060a08201526020810192909252810184905260ff8716606082015260c0015b60405180910390a2505050611444565b868686336040516020016110f79493929190614c53565b6040516020818303038152906040528051906020012060001c975087821461119a576111233384613213565b604080516080808252600e908201527f436f6d6d69744d69736d6174636800000000000000000000000000000000000060a08201526020810184905290810189905260ff8516606082015233907f04ca4e0efda95f8b780c116574d1521309010b38d8f7b75705495703a0f570b19060c0016110d0565b600e54861461121e57600e546040805160808082526014908201527f5265706f72744c656e6774684d69736d6174636800000000000000000000000060a0820152602081018990529081019190915260ff8516606082015233907f04ca4e0efda95f8b780c116574d1521309010b38d8f7b75705495703a0f570b19060c0016110d0565b60005b868110156113d957606488888381811061123d5761123d614b2b565b90506040020160200160208101906112559190614cd7565b60ff1611156112a65760405162461bcd60e51b815260206004820152601860248201527f696e76616c696420636f6e666964656e63652073636f726500000000000000006044820152606401610600565b60008888838181106112ba576112ba614b2b565b6112d09260206040909202019081019150614cf4565b6effffffffffffffffffffffffffffff1611801561131a575060008888838181106112fd576112fd614b2b565b90506040020160200160208101906113159190614cd7565b60ff16115b6113665760405162461bcd60e51b815260206004820152601660248201527f636f6e666964656e63652f7072696365206572726f72000000000000000000006044820152606401610600565b87878281811061137857611378614b2b565b905060400201600d600e838154811061139357611393614b2b565b906000526020600020016040516113aa9190614d5e565b9081526040805160209281900383019020336000908152925290206113cf8282614dd3565b5050600101611221565b506004830180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff1661010017905560405160ff8516815233907f8bdddd7f2f2c74679ffa6beb8f86aa18bfa5baf1bfaf534d0b66596babc53f089060200160405180910390a25050505b61144d60008055565b5050505050565b600080546001036114a75760405162461bcd60e51b815260206004820152601d60248201527f72656164206f6e6c79207265656e7472616e63792064657465637465640000006044820152606401610600565b5060105490565b6001546001600160a01b0316331461152e5760405162461bcd60e51b815260206004820152602360248201527f7265737472696374656420746f20746865206175746f6e69747920636f6e747260448201527f61637400000000000000000000000000000000000000000000000000000000006064820152608401610600565b611536612607565b476115418183613280565b5061088460008055565b6000805460010361159e5760405162461bcd60e51b815260206004820152601d60248201527f72656164206f6e6c79207265656e7472616e63792064657465637465640000006044820152606401610600565b50600b5490565b6002546001600160a01b031633146115ff5760405162461bcd60e51b815260206004820152601660248201527f7265737472696374656420746f206f70657261746f72000000000000000000006044820152606401610600565b6116088161352a565b6010819055600354600b547f207e45ce6f2191c3efffe5d0d91591cf81afd23cbf4bd746981d7a8a8dcfe1ba91908390611643908390614be0565b604080516080808252600a908201527f766f7465506572696f640000000000000000000000000000000000000000000060a08201526020810194909452830191909152606082015260c0015b60405180910390a150565b6002546001600160a01b031633146116f45760405162461bcd60e51b815260206004820152601660248201527f7265737472696374656420746f206f70657261746f72000000000000000000006044820152606401610600565b80516000036117455760405162461bcd60e51b815260206004820152601660248201527f73796d626f6c732063616e277420626520656d707479000000000000000000006044820152606401610600565b601454600a54611756906001614b89565b141580156117685750601454600a5414155b6117b45760405162461bcd60e51b815260206004820152601e60248201527f63616e2774206265207570646174656420696e207468697320726f756e6400006044820152606401610600565b80516117c790600f906020840190614367565b50601454600a8190557faa278e424da680ce5dad66510415760e78e0bd87d45c786c6e88bdde82f9342d9082906117ff906001614be0565b60405161168f929190614e6e565b61184a6040518060c00160405280600081526020016000815260200160008152602001600081526020016000151581526020016000151581525090565b60005460010361189c5760405162461bcd60e51b815260206004820152601d60248201527f72656164206f6e6c79207265656e7472616e63792064657465637465640000006044820152606401610600565b506001600160a01b03166000908152600c6020908152604091829020825160c081018452815481526001820154928101929092526002810154928201929092526003820154606082015260049091015460ff8082161515608084015261010090910416151560a082015290565b6000805460010361195c5760405162461bcd60e51b815260206004820152601d60248201527f72656164206f6e6c79207265656e7472616e63792064657465637465640000006044820152606401610600565b5060145490565b6001546001600160a01b031633146119e35760405162461bcd60e51b815260206004820152602360248201527f7265737472696374656420746f20746865206175746f6e69747920636f6e747260448201527f61637400000000000000000000000000000000000000000000000000000000006064820152608401610600565b600280547fffffffffffffffffffffffff0000000000000000000000000000000000000000166001600160a01b0392909216919091179055565b60008054600103611a705760405162461bcd60e51b815260206004820152601d60248201527f72656164206f6e6c79207265656e7472616e63792064657465637465640000006044820152606401610600565b5060035490565b611ad860405180610120016040528060006001600160a01b0316815260200160006001600160a01b03168152602001600081526020016000815260200160008152602001600081526020016000815260200160008152602001600081525090565b600054600103611b2a5760405162461bcd60e51b815260206004820152601d60248201527f72656164206f6e6c79207265656e7472616e63792064657465637465640000006044820152606401610600565b5060408051610120810182526001546001600160a01b039081168252600254166020820152600354918101919091526004546060820152600554608082015260065460a082015260075460c082015260085460e082015260095461010082015290565b6060611b9b60005460011490565b15611be85760405162461bcd60e51b815260206004820152601d60248201527f72656164206f6e6c79207265656e7472616e63792064657465637465640000006044820152606401610600565b601180548060200260200160405190810160405280929190818152602001828054801561065f576020028201919060005260206000209081546001600160a01b03168152600190910190602001808311610641575050505050905090565b6002546001600160a01b03163314611ca05760405162461bcd60e51b815260206004820152601660248201527f7265737472696374656420746f206f70657261746f72000000000000000000006044820152606401610600565b6005546040805160808082526018908201527f6f75746c696572536c617368696e675468726573686f6c64000000000000000060a0820152602081019290925281018590524360608201527fb5114472b89d1126433287defb6308bbabdb95b6ce5dd949ab2c151ed7b1ff4c9060c00160405180910390a160058490556004546040805160808082526019908201527f6f75746c696572446574656374696f6e5468726573686f6c640000000000000060a0820152602081019290925281018490524360608201527fb5114472b89d1126433287defb6308bbabdb95b6ce5dd949ab2c151ed7b1ff4c9060c00160405180910390a160048390556006546040805160808082526010908201527f62617365536c617368696e67526174650000000000000000000000000000000060a0820152602081019290925281018390524360608201527f207e45ce6f2191c3efffe5d0d91591cf81afd23cbf4bd746981d7a8a8dcfe1ba9060c00160405180910390a16006829055600954604080516080808252600f908201527f736c617368696e6752617465436170000000000000000000000000000000000060a0820152602081019290925281018290524360608201527f207e45ce6f2191c3efffe5d0d91591cf81afd23cbf4bd746981d7a8a8dcfe1ba9060c00160405180910390a1600955505050565b6001546001600160a01b03163314611f175760405162461bcd60e51b815260206004820152602360248201527f7265737472696374656420746f20746865206175746f6e69747920636f6e747260448201527f61637400000000000000000000000000000000000000000000000000000000006064820152608401610600565b611f1f612607565b8251600003611f705760405162461bcd60e51b815260206004820152601560248201527f566f746572732063616e277420626520656d70747900000000000000000000006044820152606401610600565b60005b835181101561208457828181518110611f8e57611f8e614b2b565b602002602001015160186000868481518110611fac57611fac614b2b565b60200260200101516001600160a01b03166001600160a01b0316815260200190815260200160002060006101000a8154816001600160a01b0302191690836001600160a01b0316021790555081818151811061200a5761200a614b2b565b60200260200101516019600086848151811061202857612028614b2b565b6020908102919091018101516001600160a01b0390811683529082019290925260400160002080547fffffffffffffffffffffffff00000000000000000000000000000000000000001692909116919091179055600101611f73565b5061209e836000600186516120999190614bb1565b6136f0565b82516120b19060129060208601906143ad565b506013805460ff191660011790556120c860008055565b505050565b6060601454600a5460016120e19190614b89565b036121bd57600f805480602002602001604051908101604052809291908181526020016000905b828210156121b457838290600052602060002001805461212790614d11565b80601f016020809104026020016040519081016040528092919081815260200182805461215390614d11565b80156121a05780601f10612175576101008083540402835291602001916121a0565b820191906000526020600020905b81548152906001019060200180831161218357829003601f168201915b505050505081526020019060010190612108565b50505050905090565b600e805480602002602001604051908101604052809291908181526020016000905b828210156121b45783829060005260206000200180546121fe90614d11565b80601f016020809104026020016040519081016040528092919081815260200182805461222a90614d11565b80156122775780601f1061224c57610100808354040283529160200191612277565b820191906000526020600020905b81548152906001019060200180831161225a57829003601f168201915b5050505050815260200190600101906121df565b600080546001036122de5760405162461bcd60e51b815260206004820152601d60248201527f72656164206f6e6c79207265656e7472616e63792064657465637465640000006044820152606401610600565b506001600160a01b039081166000908152601860205260409020541690565b60408051808201909152600080825260208201526000546001036123635760405162461bcd60e51b815260206004820152601d60248201527f72656164206f6e6c79207265656e7472616e63792064657465637465640000006044820152606401610600565b600d836040516123739190614bc4565b90815260408051602092819003830181206001600160a01b0395909516600090815294835293819020848201909152546effffffffffffffffffffffffffffff811684526f01000000000000000000000000000000900460ff169083015250919050565b6000805b601154821080156123ed575060125481105b1561255e576012818154811061240557612405614b2b565b600091825260209091200154601180546001600160a01b03909216918490811061243157612431614b2b565b6000918252602090912001546001600160a01b03160361246b578161245581614e90565b925050808061246390614e90565b9150506123db565b6012818154811061247e5761247e614b2b565b600091825260209091200154601180546001600160a01b0390921691849081106124aa576124aa614b2b565b6000918252602090912001546001600160a01b0316101561255457600c6000601184815481106124dc576124dc614b2b565b60009182526020808320909101546001600160a01b0316835282019290925260400181208181556001810182905560028101829055600381019190915560040180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00001690558161254c81614e90565b9250506123db565b8061246381614e90565b6011548210156125f757600c60006011848154811061257f5761257f614b2b565b60009182526020808320909101546001600160a01b0316835282019290925260400181208181556001810182905560028101829055600381019190915560040180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000169055816125ef81614e90565b92505061255e565b601280546120c891601191614426565b600054156126575760405162461bcd60e51b815260206004820152601360248201527f7265656e7472616e6379206465746563746564000000000000000000000000006044820152606401610600565b6001600055565b60115460609060009067ffffffffffffffff81111561267f5761267f614550565b6040519080825280602002602001820160405280156126a8578160200160208202803683370190505b50905060005b601154811015612857576000601182815481106126cd576126cd614b2b565b60009182526020808320909101546001600160a01b0316808352600c90915260409091206014549192509061270490600190614bb1565b8154148015612717575060008160010154115b15612726576127268282613213565b6007546003820154111561284d57600184848151811061274857612748614b2b565b911515602092830291909101820152601454600383015460408051928352928201526001600160a01b038416917f9e6b40f10c60d1ad09594f3b6ed7043d0e978f584d354ace6e1f6025660c42b1910160405180910390a26000600382018190556001546001600160a01b03848116835260196020526040928390205460095493517f02fb4d850000000000000000000000000000000000000000000000000000000081529082166004820152602481019390935216906302fb4d85906044016020604051808303816000875af1158015612827573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061284b9190614eaa565b505b50506001016126ae565b50919050565b60005b601154811015610884576000600c60006011848154811061288357612883614b2b565b60009182526020808320909101546001600160a01b03168352820192909252604001902060030155600101612860565b6000600e83815481106128c8576128c8614b2b565b9060005260206000200180546128dd90614d11565b80601f016020809104026020016040519081016040528092919081815260200182805461290990614d11565b80156129565780601f1061292b57610100808354040283529160200191612956565b820191906000526020600020905b81548152906001019060200180831161293957829003601f168201915b50505050509050600060118054905067ffffffffffffffff81111561297d5761297d614550565b6040519080825280602002602001820160405280156129c257816020015b604080518082019091526000808252602082015281526020019060019003908161299b5790505b5090506000805b601154811015612ac6576000601182815481106129e8576129e8614b2b565b60009182526020808320909101546001600160a01b0316808352600c90915260409091206004015490915060ff61010090910416612a265750612abe565b600d85604051612a369190614bc4565b9081526040805191829003602090810183206001600160a01b038516600090815290825282902083830190925290546effffffffffffffffffffffffffffff8116835260ff6f0100000000000000000000000000000090910416908201528484612a9f81614e90565b955081518110612ab157612ab1614b2b565b6020026020010181905250505b6001016129c9565b508015612f41576000612ad9838361389d565b6effffffffffffffffffffffffffffff1690506000612af8828661396d565b606081015190915015612e0e5760005b8160200151811015612d16576000612bf083600001518381518110612b2f57612b2f614b2b565b602002602001015185600d8a604051612b489190614bc4565b90815260200160405180910390206000601188600001518881518110612b7057612b70614b2b565b602002602001015181548110612b8857612b88614b2b565b60009182526020808320909101546001600160a01b0316835282810193909352604091820190208151808301909252546effffffffffffffffffffffffffffff8116825260ff6f0100000000000000000000000000000090910416918101919091528b613cb6565b9050601183600001518381518110612c0a57612c0a614b2b565b602002602001015181548110612c2257612c22614b2b565b9060005260206000200160009054906101000a90046001600160a01b03166001600160a01b03167f372858b237c8bd0714183e8351a461d6c3cb1ef83806181b36bf5943711f4f57828987600d8c604051612c7d9190614bc4565b9081526020016040518091039020600060118a600001518a81518110612ca557612ca5614b2b565b602002602001015181548110612cbd57612cbd614b2b565b60009182526020808320909101546001600160a01b031683528201929092526040908101909120549051612d05949392916effffffffffffffffffffffffffffff1690614ec3565b60405180910390a250600101612b08565b506000612d2b82604001518360600151613ea5565b9050604051806060016040528082815260200142815260200160011515815250601560145481548110612d6057612d60614b2b565b9060005260206000200187604051612d789190614bc4565b9081526040805160209281900383018120845181559284015160018401559201516002909101805460ff1916911515919091179055612db8908790614bc4565b604080519182900382206014548484526020840152600183830152426060840152905190917f5f2aa51aa7889ad71d9318fa7fd83c8ff3277434249bd06073f15986e197911c919081900360800190a250612f3a565b600060156001601454612e219190614bb1565b81548110612e3157612e31614b2b565b9060005260206000200186604051612e499190614bc4565b9081526020016040518091039020600001549050604051806060016040528082815260200142815260200160001515815250601560145481548110612e9057612e90614b2b565b9060005260206000200187604051612ea89190614bc4565b9081526040805160209281900383018120845181559284015160018401559201516002909101805460ff1916911515919091179055612ee8908790614bc4565b604080519182900382206014548484526020840152600083830152426060840152905190917f5f2aa51aa7889ad71d9318fa7fd83c8ff3277434249bd06073f15986e197911c919081900360800190a2505b505061144d565b600060156001601454612f549190614bb1565b81548110612f6457612f64614b2b565b9060005260206000200184604051612f7c9190614bc4565b9081526020016040518091039020600001549050604051806060016040528082815260200142815260200160001515815250601560145481548110612fc357612fc3614b2b565b9060005260206000200185604051612fdb9190614bc4565b9081526040805160209281900383018120845181559284015160018401559201516002909101805460ff191691151591909117905561301b908590614bc4565b604080519182900382206014548484526020840152600083830152426060840152905190917f5f2aa51aa7889ad71d9318fa7fd83c8ff3277434249bd06073f15986e197911c919081900360800190a2505050505050565b60005b6011548110156108845760006011828154811061309557613095614b2b565b60009182526020808320909101546001600160a01b0316808352600c9091526040909120600201549091501561315e576130d0601682613f63565b506001600160a01b0381166000908152600c6020908152604080832060020154601a9092528220805491929091613108908490614be0565b90915550506001600160a01b0381166000908152600c6020526040812060020154601b80549192909161313c908490614be0565b90915550506001600160a01b0381166000908152600c60205260408120600201555b50600101613076565b60005b60115481101561320f5781818151811061318657613186614b2b565b602002602001015115613207576000600c6000601184815481106131ac576131ac614b2b565b60009182526020808320909101546001600160a01b0316835282019290925260400190206004018054911515610100027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff9092169190911790555b60010161316a565b5050565b60038101805490600061322583614e90565b9190505550816001600160a01b03167f176956a4e941f6737f81a3c9a09d8571dd0438d86e25a432beb2013aced430926014548360030154604051613274929190918252602082015260400190565b60405180910390a25050565b601b5460000361328e575050565b600061329a6016613f7f565b905060005b81518110156134e65760008282815181106132bc576132bc614b2b565b602002602001015190506000601b54601a6000846001600160a01b03166001600160a01b0316815260200190815260200160002054876132fc9190614f04565b6133069190614f1b565b601b546001600160a01b0384166000908152601a6020526040812054929350916133309088614f04565b61333a9190614f1b565b90508115613413576001600160a01b03838116600090815260186020526040808220549051919283929116906108fc90869084818181858888f193505050503d80600081146133a5576040519150601f19603f3d011682016040523d82523d6000602084013e6133aa565b606091505b509092509050811515600003613410576001600160a01b03808616600090815260186020526040908190205490517f1137d8c966ce69b9630fb2294be011f3d64cc56e91fad7d375f0662568e9d352926134079216908490614f2f565b60405180910390a15b50505b80156134b2576001546001600160a01b038481166000908152601960205260408082205490517ff7fcc510000000000000000000000000000000000000000000000000000000008152908316600482015260248101859052604481019190915291169063f7fcc51090606401600060405180830381600087803b15801561349957600080fd5b505af11580156134ad573d6000803e3d6000fd5b505050505b6001600160a01b0383166000908152601a60205260408120556134d6601684613f8c565b50506001909201915061329f9050565b5060408051838152602081018590527f3e5aaff9e8fd4293ae18127809c2d4069d87fe10c7de92aa39557a1edbd48fec910160405180910390a150506000601b5550565b600154604080517f0aac2da100000000000000000000000000000000000000000000000000000000815290516000926001600160a01b031691630aac2da19160048083019260209291908290030181865afa15801561358d573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906135b19190614eaa565b9050806135bf836002614f04565b111561360d5760405162461bcd60e51b815260206004820152601660248201527f766f746520706572696f6420697320746f6f20626967000000000000000000006044820152606401610600565b600154604080517fdfb1a4d200000000000000000000000000000000000000000000000000000000815290516001600160a01b039092169163dfb1a4d2916004808201926020929091908290030181865afa158015613670573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906136949190614eaa565b9050806136a2836002614f04565b111561320f5760405162461bcd60e51b815260206004820152601660248201527f766f746520706572696f6420697320746f6f20626967000000000000000000006044820152606401610600565b8082126136fc57505050565b8181600085600261370d8585614f5f565b6137179190614f7f565b6137219087614b89565b8151811061373157613731614b2b565b602002602001015190505b81831361386f575b806001600160a01b031686848151811061376057613760614b2b565b60200260200101516001600160a01b03161015613789578261378181614fc9565b935050613744565b806001600160a01b03168683815181106137a5576137a5614b2b565b60200260200101516001600160a01b031611156137ce57816137c681614ffa565b925050613789565b81831361386a578582815181106137e7576137e7614b2b565b602002602001015186848151811061380157613801614b2b565b602002602001015187858151811061381b5761381b614b2b565b6020026020010188858151811061383457613834614b2b565b6001600160a01b039384166020918202929092010152911690528261385881614fc9565b935050818061386690614ffa565b9250505b61373c565b81851215613882576138828686846136f0565b83831215613895576138958684866136f0565b505050505050565b6000816000036138af57506000610bb5565b6138c58360006138c0600186614bb1565b613fa1565b60006138d2600284614f1b565b90506138df600284614c22565b15613907578381815181106138f6576138f6614b2b565b602002602001015160000151613965565b600284828151811061391b5761391b614b2b565b602002602001015160000151856001846139359190614bb1565b8151811061394557613945614b2b565b60200260200101516000015161395b9190615033565b613965919061505a565b949350505050565b6139986040518060800160405280606081526020016000815260200160608152602001600081525090565b6139c36040518060800160405280606081526020016000815260200160608152602001600081525090565b60115467ffffffffffffffff8111156139de576139de614550565b604051908082528060200260200182016040528015613a2357816020015b60408051808201909152600080825260208201528152602001906001900390816139fc5790505b50604082015260115467ffffffffffffffff811115613a4457613a44614550565b604051908082528060200260200182016040528015613a6d578160200160208202803683370190505b50815260005b601154811015613cae57600060118281548110613a9257613a92614b2b565b60009182526020808320909101546001600160a01b0316808352600c90915260409091206004015490915060ff61010090910416613ad05750613ca6565b600086600d87604051613ae39190614bc4565b90815260408051602092819003830190206001600160a01b03861660009081529252902054613b23906effffffffffffffffffffffffffffff1689614f5f565b613b2e906064615098565b613b389190614f7f565b6004549091508113801590613b5a5750600454613b5782600019615098565b13155b15613c7057600d86604051613b6f9190614bc4565b9081526040805191829003602090810183206001600160a01b03861660009081529082528290208383018352546effffffffffffffffffffffffffffff8116845260ff6f01000000000000000000000000000000909104169083015285015160608601805190613bde82614e90565b905281518110613bf057613bf0614b2b565b6020026020010181905250600d86604051613c0b9190614bc4565b90815260408051602092819003830190206001600160a01b038516600090815290835281812054600c90935290812060020180546f0100000000000000000000000000000090930460ff1692909190613c65908490614be0565b90915550613ca39050565b8351602085018051859291613c8482614e90565b905281518110613c9657613c96614b2b565b6020026020010181815250505b50505b600101613a73565b509392505050565b60008060118681548110613ccc57613ccc614b2b565b60009182526020808320909101546001600160a01b0316808352600c909152604090912060040180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff1690558351909150839087908110613d3057613d30614b2b565b602002602001015115613d47576000915050613965565b6000858686600001516effffffffffffffffffffffffffffff16613d6b9190614f5f565b613d76906064615098565b613d809190614f7f565b9050613d8c8180615098565b6005549091508113613da357600092505050613965565b6000612710600160050154876020015160ff1660016004015485613dc79190614f5f565b613dd19190614f04565b613ddb9190614f04565b613de59190614f1b565b600954909150811115613df757506009545b6001546001600160a01b03848116600090815260196020526040908190205490517f02fb4d850000000000000000000000000000000000000000000000000000000081529082166004820152602481018490529116906302fb4d85906044016020604051808303816000875af1158015613e75573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613e999190614eaa565b98975050505050505050565b60008080805b84811015613f4f57858181518110613ec557613ec5614b2b565b60200260200101516020015160ff16868281518110613ee657613ee6614b2b565b602002602001015160000151613efc91906150e4565b613f17906effffffffffffffffffffffffffffff1683614be0565b9150858181518110613f2b57613f2b614b2b565b60200260200101516020015160ff1683613f459190614be0565b9250600101613eab565b50613f5a8282614f1b565b95945050505050565b6000613f78836001600160a01b03841661416a565b9392505050565b60606000613f78836141b9565b6000613f78836001600160a01b038416614215565b8181808203613fb1575050505050565b6000856002613fc08787614f5f565b613fca9190614f7f565b613fd49087614b89565b81518110613fe457613fe4614b2b565b60200260200101516000015190505b818313614144575b806effffffffffffffffffffffffffffff1686848151811061401f5761401f614b2b565b6020026020010151600001516effffffffffffffffffffffffffffff161015614054578261404c81614fc9565b935050613ffb565b85828151811061406657614066614b2b565b6020026020010151600001516effffffffffffffffffffffffffffff16816effffffffffffffffffffffffffffff1610156140ad57816140a581614ffa565b925050614054565b81831361413f578582815181106140c6576140c6614b2b565b60200260200101518684815181106140e0576140e0614b2b565b60200260200101518785815181106140fa576140fa614b2b565b6020026020010188858151811061411357614113614b2b565b602002602001018290528290525050818061412d90614ffa565b925050828061413b90614fc9565b9350505b613ff3565b8185121561415757614157868684613fa1565b8383121561389557613895868486613fa1565b60008181526001830160205260408120546141b157508154600181810184556000848152602080822090930184905584548482528286019093526040902091909155610bb5565b506000610bb5565b60608160000180548060200260200160405190810160405280929190818152602001828054801561420957602002820191906000526020600020905b8154815260200190600101908083116141f5575b50505050509050919050565b600081815260018301602052604081205480156142fe576000614239600183614bb1565b855490915060009061424d90600190614bb1565b90508082146142b257600086600001828154811061426d5761426d614b2b565b906000526020600020015490508087600001848154811061429057614290614b2b565b6000918252602080832090910192909255918252600188019052604090208390555b85548690806142c3576142c361510e565b600190038181906000526020600020016000905590558560010160008681526020019081526020016000206000905560019350505050610bb5565b6000915050610bb5565b5092915050565b8280548282559060005260206000209081019282156143575760005260206000209182015b8281111561435757816143478482615184565b5091600101919060010190614334565b50614363929150614466565b5090565b828054828255906000526020600020908101928215614357579160200282015b82811115614357578251829061439d9082615264565b5091602001919060010190614387565b82805482825590600052602060002090810192821561441a579160200282015b8281111561441a57825182547fffffffffffffffffffffffff0000000000000000000000000000000000000000166001600160a01b039091161782556020909201916001909101906143cd565b50614363929150614483565b82805482825590600052602060002090810192821561441a5760005260206000209182015b8281111561441a57825482559160010191906001019061444b565b8082111561436357600061447a8282614498565b50600101614466565b5b808211156143635760008155600101614484565b5080546144a490614d11565b6000825580601f106144b4575050565b601f0160209004906000526020600020908101906108849190614483565b602080825282518282018190526000918401906040840190835b818110156145135783516001600160a01b03168352602093840193909201916001016144ec565b509095945050505050565b80356001600160a01b03811681146108f857600080fd5b60006020828403121561454757600080fd5b613f788261451e565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051601f8201601f1916810167ffffffffffffffff811182821017156145a8576145a8614550565b604052919050565b600082601f8301126145c157600080fd5b813567ffffffffffffffff8111156145db576145db614550565b6145ee6020601f19601f8401160161457f565b81815284602083860101111561460357600080fd5b816020850160208301376000918101602001919091529392505050565b60006020828403121561463257600080fd5b813567ffffffffffffffff81111561464957600080fd5b613965848285016145b0565b6000806040838503121561466857600080fd5b82359150602083013567ffffffffffffffff81111561468657600080fd5b614692858286016145b0565b9150509250929050565b600080604083850312156146af57600080fd5b50508035926020909101359150565b60ff8116811461088457600080fd5b80356108f8816146be565b6000806000806000608086880312156146f057600080fd5b85359450602086013567ffffffffffffffff81111561470e57600080fd5b8601601f8101881361471f57600080fd5b803567ffffffffffffffff81111561473657600080fd5b8860208260061b840101111561474b57600080fd5b6020919091019450925060408601359150614768606087016146cd565b90509295509295909350565b60006020828403121561478657600080fd5b5035919050565b600067ffffffffffffffff8211156147a7576147a7614550565b5060051b60200190565b6000602082840312156147c357600080fd5b813567ffffffffffffffff8111156147da57600080fd5b8201601f810184136147eb57600080fd5b80356147fe6147f98261478d565b61457f565b8082825260208201915060208360051b85010192508683111561482057600080fd5b602084015b8381101561486257803567ffffffffffffffff81111561484457600080fd5b614853896020838901016145b0565b84525060209283019201614825565b509695505050505050565b81516001600160a01b0316815260208083015161012083019161489a908401826001600160a01b03169052565b5060408301516040830152606083015160608301526080830151608083015260a083015160a083015260c083015160c083015260e083015160e083015261010083015161010083015292915050565b600080600080608085870312156148ff57600080fd5b5050823594602084013594506040840135936060013592509050565b600082601f83011261492c57600080fd5b813561493a6147f98261478d565b8082825260208201915060208360051b86010192508583111561495c57600080fd5b602085015b83811015614980576149728161451e565b835260209283019201614961565b5095945050505050565b60008060006060848603121561499f57600080fd5b833567ffffffffffffffff8111156149b657600080fd5b6149c28682870161491b565b935050602084013567ffffffffffffffff8111156149df57600080fd5b6149eb8682870161491b565b925050604084013567ffffffffffffffff811115614a0857600080fd5b614a148682870161491b565b9150509250925092565b60005b83811015614a39578181015183820152602001614a21565b50506000910152565b60008151808452614a5a816020860160208601614a1e565b601f01601f19169290920160200192915050565b600082825180855260208501945060208160051b8301016020850160005b83811015614abe57601f19858403018852614aa8838351614a42565b6020988901989093509190910190600101614a8c565b50909695505050505050565b602081526000613f786020830184614a6e565b60008060408385031215614af057600080fd5b823567ffffffffffffffff811115614b0757600080fd5b614b13858286016145b0565b925050614b226020840161451e565b90509250929050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b8082018281126000831280158216821582161715614ba957614ba9614b5a565b505092915050565b81810381811115610bb557610bb5614b5a565b60008251614bd6818460208701614a1e565b9190910192915050565b80820180821115610bb557610bb5614b5a565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b600082614c3157614c31614bf3565b500690565b6effffffffffffffffffffffffffffff8116811461088457600080fd5b6060808252810184905260008560808301825b87811015614cb4578235614c7981614c36565b6effffffffffffffffffffffffffffff1682526020830135614c9a816146be565b60ff16602083015260409283019290910190600101614c66565b50602084019590955250506001600160a01b039190911660409091015292915050565b600060208284031215614ce957600080fd5b8135613f78816146be565b600060208284031215614d0657600080fd5b8135613f7881614c36565b600181811c90821680614d2557607f821691505b602082108103612857577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b6000808354614d6c81614d11565b600182168015614d835760018114614d9857614dc8565b60ff1983168652811515820286019350614dc8565b86600052602060002060005b83811015614dc057815488820152600190910190602001614da4565b505081860193505b509195945050505050565b8135614dde81614c36565b6effffffffffffffffffffffffffffff811690508154817fffffffffffffffffffffffffffffffffff00000000000000000000000000000082161783556020840135614e29816146be565b6fff0000000000000000000000000000008160781b16837fffffffffffffffffffffffffffffffff000000000000000000000000000000008416171784555050505050565b604081526000614e816040830185614a6e565b90508260208301529392505050565b60006000198203614ea357614ea3614b5a565b5060010190565b600060208284031215614ebc57600080fd5b5051919050565b848152608060208201526000614edc6080830186614a42565b90508360408301526effffffffffffffffffffffffffffff8316606083015295945050505050565b8082028115828204841417610bb557610bb5614b5a565b600082614f2a57614f2a614bf3565b500490565b6001600160a01b038316815260606020820152600060608201526080604082015260006139656080830184614a42565b818103600083128015838313168383128216171561430857614308614b5a565b600082614f8e57614f8e614bf3565b60001983147f800000000000000000000000000000000000000000000000000000000000000083141615614fc457614fc4614b5a565b500590565b60007f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8203614ea357614ea3614b5a565b60007f8000000000000000000000000000000000000000000000000000000000000000820361502b5761502b614b5a565b506000190190565b6effffffffffffffffffffffffffffff8181168382160190811115610bb557610bb5614b5a565b60006effffffffffffffffffffffffffffff83168061507b5761507b614bf3565b806effffffffffffffffffffffffffffff84160491505092915050565b808202600082127f8000000000000000000000000000000000000000000000000000000000000000841416156150d0576150d0614b5a565b8181058314821517610bb557610bb5614b5a565b6effffffffffffffffffffffffffffff818116838216029081169081811461430857614308614b5a565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fd5b601f8211156120c857806000526020600020601f840160051c810160208510156151645750805b601f840160051c820191505b8181101561144d5760008155600101615170565b81810361518f575050565b6151998254614d11565b67ffffffffffffffff8111156151b1576151b1614550565b6151c5816151bf8454614d11565b8461513d565b6000601f8211600181146151fc57600083156151e15750848201545b600184901b600019600386901b1c198216175b85555061144d565b600085815260209020601f19841690600086815260209020845b838110156152365782860154825560019586019590910190602001615216565b50858310156152545781850154600019600388901b60f8161c191681555b5050505050600190811b01905550565b815167ffffffffffffffff81111561527e5761527e614550565b61528c816151bf8454614d11565b6020601f8211600181146152be57600083156151e1575081850151600184901b600019600386901b1c198216176151f4565b600084815260208120601f198516915b828110156152ee57878501518255602094850194600190920191016152ce565b508482101561530c5786840151600019600387901b60f8161c191681555b50505050600190811b0190555056fea26469706673582212209332f2eb7a6ca90dad8644e835fd01558f41287c53fb4f634a9a437d9ac6c96a64736f6c634300081e0033",
}

// Oracle0ABI is the input ABI used to generate the binding from.
// Deprecated: Use Oracle0MetaData.ABI instead.
var Oracle0ABI = Oracle0MetaData.ABI

// Deprecated: Use Oracle0MetaData.Sigs instead.
// Oracle0FuncSigs maps the 4-byte function signature to its string representation.
var Oracle0FuncSigs = Oracle0MetaData.Sigs

// Oracle0Bin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use Oracle0MetaData.Bin instead.
var Oracle0Bin = Oracle0MetaData.Bin

// DeployOracle0 deploys a new Ethereum contract, binding an instance of Oracle0 to it.
func DeployOracle0(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Oracle0, error) {
	parsed, err := Oracle0MetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(Oracle0Bin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Oracle0{Oracle0Caller: Oracle0Caller{contract: contract}, Oracle0Transactor: Oracle0Transactor{contract: contract}, Oracle0Filterer: Oracle0Filterer{contract: contract}}, nil
}

// Oracle0 is an auto generated Go binding around an Ethereum contract.
type Oracle0 struct {
	Oracle0Caller     // Read-only binding to the contract
	Oracle0Transactor // Write-only binding to the contract
	Oracle0Filterer   // Log filterer for contract events
}

// Oracle0Caller is an auto generated read-only Go binding around an Ethereum contract.
type Oracle0Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Oracle0Transactor is an auto generated write-only Go binding around an Ethereum contract.
type Oracle0Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Oracle0Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type Oracle0Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Oracle0Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type Oracle0Session struct {
	Contract     *Oracle0          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// Oracle0CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type Oracle0CallerSession struct {
	Contract *Oracle0Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// Oracle0TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type Oracle0TransactorSession struct {
	Contract     *Oracle0Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// Oracle0Raw is an auto generated low-level Go binding around an Ethereum contract.
type Oracle0Raw struct {
	Contract *Oracle0 // Generic contract binding to access the raw methods on
}

// Oracle0CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type Oracle0CallerRaw struct {
	Contract *Oracle0Caller // Generic read-only contract binding to access the raw methods on
}

// Oracle0TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type Oracle0TransactorRaw struct {
	Contract *Oracle0Transactor // Generic write-only contract binding to access the raw methods on
}

// NewOracle0 creates a new instance of Oracle0, bound to a specific deployed contract.
func NewOracle0(address common.Address, backend bind.ContractBackend) (*Oracle0, error) {
	contract, err := bindOracle0(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Oracle0{Oracle0Caller: Oracle0Caller{contract: contract}, Oracle0Transactor: Oracle0Transactor{contract: contract}, Oracle0Filterer: Oracle0Filterer{contract: contract}}, nil
}

// NewOracle0Caller creates a new read-only instance of Oracle0, bound to a specific deployed contract.
func NewOracle0Caller(address common.Address, caller bind.ContractCaller) (*Oracle0Caller, error) {
	contract, err := bindOracle0(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &Oracle0Caller{contract: contract}, nil
}

// NewOracle0Transactor creates a new write-only instance of Oracle0, bound to a specific deployed contract.
func NewOracle0Transactor(address common.Address, transactor bind.ContractTransactor) (*Oracle0Transactor, error) {
	contract, err := bindOracle0(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &Oracle0Transactor{contract: contract}, nil
}

// NewOracle0Filterer creates a new log filterer instance of Oracle0, bound to a specific deployed contract.
func NewOracle0Filterer(address common.Address, filterer bind.ContractFilterer) (*Oracle0Filterer, error) {
	contract, err := bindOracle0(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &Oracle0Filterer{contract: contract}, nil
}

// bindOracle0 binds a generic wrapper to an already deployed contract.
func bindOracle0(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(Oracle0ABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Oracle0 *Oracle0Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Oracle0.Contract.Oracle0Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Oracle0 *Oracle0Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Oracle0.Contract.Oracle0Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Oracle0 *Oracle0Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Oracle0.Contract.Oracle0Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Oracle0 *Oracle0CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Oracle0.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Oracle0 *Oracle0TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Oracle0.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Oracle0 *Oracle0TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Oracle0.Contract.contract.Transact(opts, method, params...)
}

// GetConfig is a free data retrieval call binding the contract method 0xc3f909d4.
//
// Solidity: function getConfig() view returns((address,address,uint256,int256,int256,uint256,uint256,uint256,uint256))
func (_Oracle0 *Oracle0Caller) GetConfig(opts *bind.CallOpts) (Oracle0Config, error) {
	var out []interface{}
	err := _Oracle0.contract.Call(opts, &out, "getConfig")

	if err != nil {
		return *new(Oracle0Config), err
	}

	out0 := *abi.ConvertType(out[0], new(Oracle0Config)).(*Oracle0Config)

	return out0, err

}

// GetConfig is a free data retrieval call binding the contract method 0xc3f909d4.
//
// Solidity: function getConfig() view returns((address,address,uint256,int256,int256,uint256,uint256,uint256,uint256))
func (_Oracle0 *Oracle0Session) GetConfig() (Oracle0Config, error) {
	return _Oracle0.Contract.GetConfig(&_Oracle0.CallOpts)
}

// GetConfig is a free data retrieval call binding the contract method 0xc3f909d4.
//
// Solidity: function getConfig() view returns((address,address,uint256,int256,int256,uint256,uint256,uint256,uint256))
func (_Oracle0 *Oracle0CallerSession) GetConfig() (Oracle0Config, error) {
	return _Oracle0.Contract.GetConfig(&_Oracle0.CallOpts)
}

// GetDecimals is a free data retrieval call binding the contract method 0xf0141d84.
//
// Solidity: function getDecimals() pure returns(uint8)
func (_Oracle0 *Oracle0Caller) GetDecimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _Oracle0.contract.Call(opts, &out, "getDecimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// GetDecimals is a free data retrieval call binding the contract method 0xf0141d84.
//
// Solidity: function getDecimals() pure returns(uint8)
func (_Oracle0 *Oracle0Session) GetDecimals() (uint8, error) {
	return _Oracle0.Contract.GetDecimals(&_Oracle0.CallOpts)
}

// GetDecimals is a free data retrieval call binding the contract method 0xf0141d84.
//
// Solidity: function getDecimals() pure returns(uint8)
func (_Oracle0 *Oracle0CallerSession) GetDecimals() (uint8, error) {
	return _Oracle0.Contract.GetDecimals(&_Oracle0.CallOpts)
}

// GetLastRoundBlock is a free data retrieval call binding the contract method 0x5a4d3a27.
//
// Solidity: function getLastRoundBlock() view returns(uint256)
func (_Oracle0 *Oracle0Caller) GetLastRoundBlock(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Oracle0.contract.Call(opts, &out, "getLastRoundBlock")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetLastRoundBlock is a free data retrieval call binding the contract method 0x5a4d3a27.
//
// Solidity: function getLastRoundBlock() view returns(uint256)
func (_Oracle0 *Oracle0Session) GetLastRoundBlock() (*big.Int, error) {
	return _Oracle0.Contract.GetLastRoundBlock(&_Oracle0.CallOpts)
}

// GetLastRoundBlock is a free data retrieval call binding the contract method 0x5a4d3a27.
//
// Solidity: function getLastRoundBlock() view returns(uint256)
func (_Oracle0 *Oracle0CallerSession) GetLastRoundBlock() (*big.Int, error) {
	return _Oracle0.Contract.GetLastRoundBlock(&_Oracle0.CallOpts)
}

// GetNewVotePeriod is a free data retrieval call binding the contract method 0x57eba759.
//
// Solidity: function getNewVotePeriod() view returns(uint256)
func (_Oracle0 *Oracle0Caller) GetNewVotePeriod(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Oracle0.contract.Call(opts, &out, "getNewVotePeriod")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetNewVotePeriod is a free data retrieval call binding the contract method 0x57eba759.
//
// Solidity: function getNewVotePeriod() view returns(uint256)
func (_Oracle0 *Oracle0Session) GetNewVotePeriod() (*big.Int, error) {
	return _Oracle0.Contract.GetNewVotePeriod(&_Oracle0.CallOpts)
}

// GetNewVotePeriod is a free data retrieval call binding the contract method 0x57eba759.
//
// Solidity: function getNewVotePeriod() view returns(uint256)
func (_Oracle0 *Oracle0CallerSession) GetNewVotePeriod() (*big.Int, error) {
	return _Oracle0.Contract.GetNewVotePeriod(&_Oracle0.CallOpts)
}

// GetNewVoters is a free data retrieval call binding the contract method 0x077945d3.
//
// Solidity: function getNewVoters() view returns(address[])
func (_Oracle0 *Oracle0Caller) GetNewVoters(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _Oracle0.contract.Call(opts, &out, "getNewVoters")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// GetNewVoters is a free data retrieval call binding the contract method 0x077945d3.
//
// Solidity: function getNewVoters() view returns(address[])
func (_Oracle0 *Oracle0Session) GetNewVoters() ([]common.Address, error) {
	return _Oracle0.Contract.GetNewVoters(&_Oracle0.CallOpts)
}

// GetNewVoters is a free data retrieval call binding the contract method 0x077945d3.
//
// Solidity: function getNewVoters() view returns(address[])
func (_Oracle0 *Oracle0CallerSession) GetNewVoters() ([]common.Address, error) {
	return _Oracle0.Contract.GetNewVoters(&_Oracle0.CallOpts)
}

// GetNonRevealThreshold is a free data retrieval call binding the contract method 0xed78349d.
//
// Solidity: function getNonRevealThreshold() view returns(uint256)
func (_Oracle0 *Oracle0Caller) GetNonRevealThreshold(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Oracle0.contract.Call(opts, &out, "getNonRevealThreshold")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetNonRevealThreshold is a free data retrieval call binding the contract method 0xed78349d.
//
// Solidity: function getNonRevealThreshold() view returns(uint256)
func (_Oracle0 *Oracle0Session) GetNonRevealThreshold() (*big.Int, error) {
	return _Oracle0.Contract.GetNonRevealThreshold(&_Oracle0.CallOpts)
}

// GetNonRevealThreshold is a free data retrieval call binding the contract method 0xed78349d.
//
// Solidity: function getNonRevealThreshold() view returns(uint256)
func (_Oracle0 *Oracle0CallerSession) GetNonRevealThreshold() (*big.Int, error) {
	return _Oracle0.Contract.GetNonRevealThreshold(&_Oracle0.CallOpts)
}

// GetReports is a free data retrieval call binding the contract method 0xfb09917e.
//
// Solidity: function getReports(string _symbol, address _voter) view returns((uint120,uint8))
func (_Oracle0 *Oracle0Caller) GetReports(opts *bind.CallOpts, _symbol string, _voter common.Address) (IOracleReport, error) {
	var out []interface{}
	err := _Oracle0.contract.Call(opts, &out, "getReports", _symbol, _voter)

	if err != nil {
		return *new(IOracleReport), err
	}

	out0 := *abi.ConvertType(out[0], new(IOracleReport)).(*IOracleReport)

	return out0, err

}

// GetReports is a free data retrieval call binding the contract method 0xfb09917e.
//
// Solidity: function getReports(string _symbol, address _voter) view returns((uint120,uint8))
func (_Oracle0 *Oracle0Session) GetReports(_symbol string, _voter common.Address) (IOracleReport, error) {
	return _Oracle0.Contract.GetReports(&_Oracle0.CallOpts, _symbol, _voter)
}

// GetReports is a free data retrieval call binding the contract method 0xfb09917e.
//
// Solidity: function getReports(string _symbol, address _voter) view returns((uint120,uint8))
func (_Oracle0 *Oracle0CallerSession) GetReports(_symbol string, _voter common.Address) (IOracleReport, error) {
	return _Oracle0.Contract.GetReports(&_Oracle0.CallOpts, _symbol, _voter)
}

// GetRewardPeriodPerformance is a free data retrieval call binding the contract method 0x33d16293.
//
// Solidity: function getRewardPeriodPerformance(address _voter) view returns(uint256)
func (_Oracle0 *Oracle0Caller) GetRewardPeriodPerformance(opts *bind.CallOpts, _voter common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Oracle0.contract.Call(opts, &out, "getRewardPeriodPerformance", _voter)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetRewardPeriodPerformance is a free data retrieval call binding the contract method 0x33d16293.
//
// Solidity: function getRewardPeriodPerformance(address _voter) view returns(uint256)
func (_Oracle0 *Oracle0Session) GetRewardPeriodPerformance(_voter common.Address) (*big.Int, error) {
	return _Oracle0.Contract.GetRewardPeriodPerformance(&_Oracle0.CallOpts, _voter)
}

// GetRewardPeriodPerformance is a free data retrieval call binding the contract method 0x33d16293.
//
// Solidity: function getRewardPeriodPerformance(address _voter) view returns(uint256)
func (_Oracle0 *Oracle0CallerSession) GetRewardPeriodPerformance(_voter common.Address) (*big.Int, error) {
	return _Oracle0.Contract.GetRewardPeriodPerformance(&_Oracle0.CallOpts, _voter)
}

// GetRound is a free data retrieval call binding the contract method 0x9f8743f7.
//
// Solidity: function getRound() view returns(uint256)
func (_Oracle0 *Oracle0Caller) GetRound(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Oracle0.contract.Call(opts, &out, "getRound")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetRound is a free data retrieval call binding the contract method 0x9f8743f7.
//
// Solidity: function getRound() view returns(uint256)
func (_Oracle0 *Oracle0Session) GetRound() (*big.Int, error) {
	return _Oracle0.Contract.GetRound(&_Oracle0.CallOpts)
}

// GetRound is a free data retrieval call binding the contract method 0x9f8743f7.
//
// Solidity: function getRound() view returns(uint256)
func (_Oracle0 *Oracle0CallerSession) GetRound() (*big.Int, error) {
	return _Oracle0.Contract.GetRound(&_Oracle0.CallOpts)
}

// GetRoundData is a free data retrieval call binding the contract method 0x3c8510fd.
//
// Solidity: function getRoundData(uint256 _round, string _symbol) view returns((uint256,uint256,uint256,bool) data)
func (_Oracle0 *Oracle0Caller) GetRoundData(opts *bind.CallOpts, _round *big.Int, _symbol string) (IOracleRoundData, error) {
	var out []interface{}
	err := _Oracle0.contract.Call(opts, &out, "getRoundData", _round, _symbol)

	if err != nil {
		return *new(IOracleRoundData), err
	}

	out0 := *abi.ConvertType(out[0], new(IOracleRoundData)).(*IOracleRoundData)

	return out0, err

}

// GetRoundData is a free data retrieval call binding the contract method 0x3c8510fd.
//
// Solidity: function getRoundData(uint256 _round, string _symbol) view returns((uint256,uint256,uint256,bool) data)
func (_Oracle0 *Oracle0Session) GetRoundData(_round *big.Int, _symbol string) (IOracleRoundData, error) {
	return _Oracle0.Contract.GetRoundData(&_Oracle0.CallOpts, _round, _symbol)
}

// GetRoundData is a free data retrieval call binding the contract method 0x3c8510fd.
//
// Solidity: function getRoundData(uint256 _round, string _symbol) view returns((uint256,uint256,uint256,bool) data)
func (_Oracle0 *Oracle0CallerSession) GetRoundData(_round *big.Int, _symbol string) (IOracleRoundData, error) {
	return _Oracle0.Contract.GetRoundData(&_Oracle0.CallOpts, _round, _symbol)
}

// GetSymbolUpdatedRound is a free data retrieval call binding the contract method 0x99b0014b.
//
// Solidity: function getSymbolUpdatedRound() view returns(int256)
func (_Oracle0 *Oracle0Caller) GetSymbolUpdatedRound(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Oracle0.contract.Call(opts, &out, "getSymbolUpdatedRound")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetSymbolUpdatedRound is a free data retrieval call binding the contract method 0x99b0014b.
//
// Solidity: function getSymbolUpdatedRound() view returns(int256)
func (_Oracle0 *Oracle0Session) GetSymbolUpdatedRound() (*big.Int, error) {
	return _Oracle0.Contract.GetSymbolUpdatedRound(&_Oracle0.CallOpts)
}

// GetSymbolUpdatedRound is a free data retrieval call binding the contract method 0x99b0014b.
//
// Solidity: function getSymbolUpdatedRound() view returns(int256)
func (_Oracle0 *Oracle0CallerSession) GetSymbolUpdatedRound() (*big.Int, error) {
	return _Oracle0.Contract.GetSymbolUpdatedRound(&_Oracle0.CallOpts)
}

// GetSymbols is a free data retrieval call binding the contract method 0xdf7f710e.
//
// Solidity: function getSymbols() view returns(string[])
func (_Oracle0 *Oracle0Caller) GetSymbols(opts *bind.CallOpts) ([]string, error) {
	var out []interface{}
	err := _Oracle0.contract.Call(opts, &out, "getSymbols")

	if err != nil {
		return *new([]string), err
	}

	out0 := *abi.ConvertType(out[0], new([]string)).(*[]string)

	return out0, err

}

// GetSymbols is a free data retrieval call binding the contract method 0xdf7f710e.
//
// Solidity: function getSymbols() view returns(string[])
func (_Oracle0 *Oracle0Session) GetSymbols() ([]string, error) {
	return _Oracle0.Contract.GetSymbols(&_Oracle0.CallOpts)
}

// GetSymbols is a free data retrieval call binding the contract method 0xdf7f710e.
//
// Solidity: function getSymbols() view returns(string[])
func (_Oracle0 *Oracle0CallerSession) GetSymbols() ([]string, error) {
	return _Oracle0.Contract.GetSymbols(&_Oracle0.CallOpts)
}

// GetVotePeriod is a free data retrieval call binding the contract method 0xb78dec52.
//
// Solidity: function getVotePeriod() view returns(uint256)
func (_Oracle0 *Oracle0Caller) GetVotePeriod(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Oracle0.contract.Call(opts, &out, "getVotePeriod")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetVotePeriod is a free data retrieval call binding the contract method 0xb78dec52.
//
// Solidity: function getVotePeriod() view returns(uint256)
func (_Oracle0 *Oracle0Session) GetVotePeriod() (*big.Int, error) {
	return _Oracle0.Contract.GetVotePeriod(&_Oracle0.CallOpts)
}

// GetVotePeriod is a free data retrieval call binding the contract method 0xb78dec52.
//
// Solidity: function getVotePeriod() view returns(uint256)
func (_Oracle0 *Oracle0CallerSession) GetVotePeriod() (*big.Int, error) {
	return _Oracle0.Contract.GetVotePeriod(&_Oracle0.CallOpts)
}

// GetVoterInfo is a free data retrieval call binding the contract method 0x9ed1f255.
//
// Solidity: function getVoterInfo(address _voter) view returns((uint256,uint256,uint256,uint256,bool,bool))
func (_Oracle0 *Oracle0Caller) GetVoterInfo(opts *bind.CallOpts, _voter common.Address) (Oracle0VoterInfo, error) {
	var out []interface{}
	err := _Oracle0.contract.Call(opts, &out, "getVoterInfo", _voter)

	if err != nil {
		return *new(Oracle0VoterInfo), err
	}

	out0 := *abi.ConvertType(out[0], new(Oracle0VoterInfo)).(*Oracle0VoterInfo)

	return out0, err

}

// GetVoterInfo is a free data retrieval call binding the contract method 0x9ed1f255.
//
// Solidity: function getVoterInfo(address _voter) view returns((uint256,uint256,uint256,uint256,bool,bool))
func (_Oracle0 *Oracle0Session) GetVoterInfo(_voter common.Address) (Oracle0VoterInfo, error) {
	return _Oracle0.Contract.GetVoterInfo(&_Oracle0.CallOpts, _voter)
}

// GetVoterInfo is a free data retrieval call binding the contract method 0x9ed1f255.
//
// Solidity: function getVoterInfo(address _voter) view returns((uint256,uint256,uint256,uint256,bool,bool))
func (_Oracle0 *Oracle0CallerSession) GetVoterInfo(_voter common.Address) (Oracle0VoterInfo, error) {
	return _Oracle0.Contract.GetVoterInfo(&_Oracle0.CallOpts, _voter)
}

// GetVoterTreasuries is a free data retrieval call binding the contract method 0xef5cc4d1.
//
// Solidity: function getVoterTreasuries(address _oracleAddress) view returns(address)
func (_Oracle0 *Oracle0Caller) GetVoterTreasuries(opts *bind.CallOpts, _oracleAddress common.Address) (common.Address, error) {
	var out []interface{}
	err := _Oracle0.contract.Call(opts, &out, "getVoterTreasuries", _oracleAddress)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetVoterTreasuries is a free data retrieval call binding the contract method 0xef5cc4d1.
//
// Solidity: function getVoterTreasuries(address _oracleAddress) view returns(address)
func (_Oracle0 *Oracle0Session) GetVoterTreasuries(_oracleAddress common.Address) (common.Address, error) {
	return _Oracle0.Contract.GetVoterTreasuries(&_Oracle0.CallOpts, _oracleAddress)
}

// GetVoterTreasuries is a free data retrieval call binding the contract method 0xef5cc4d1.
//
// Solidity: function getVoterTreasuries(address _oracleAddress) view returns(address)
func (_Oracle0 *Oracle0CallerSession) GetVoterTreasuries(_oracleAddress common.Address) (common.Address, error) {
	return _Oracle0.Contract.GetVoterTreasuries(&_Oracle0.CallOpts, _oracleAddress)
}

// GetVoterValidators is a free data retrieval call binding the contract method 0x2d35d158.
//
// Solidity: function getVoterValidators(address _oracleAddress) view returns(address)
func (_Oracle0 *Oracle0Caller) GetVoterValidators(opts *bind.CallOpts, _oracleAddress common.Address) (common.Address, error) {
	var out []interface{}
	err := _Oracle0.contract.Call(opts, &out, "getVoterValidators", _oracleAddress)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetVoterValidators is a free data retrieval call binding the contract method 0x2d35d158.
//
// Solidity: function getVoterValidators(address _oracleAddress) view returns(address)
func (_Oracle0 *Oracle0Session) GetVoterValidators(_oracleAddress common.Address) (common.Address, error) {
	return _Oracle0.Contract.GetVoterValidators(&_Oracle0.CallOpts, _oracleAddress)
}

// GetVoterValidators is a free data retrieval call binding the contract method 0x2d35d158.
//
// Solidity: function getVoterValidators(address _oracleAddress) view returns(address)
func (_Oracle0 *Oracle0CallerSession) GetVoterValidators(_oracleAddress common.Address) (common.Address, error) {
	return _Oracle0.Contract.GetVoterValidators(&_Oracle0.CallOpts, _oracleAddress)
}

// GetVoters is a free data retrieval call binding the contract method 0xcdd72253.
//
// Solidity: function getVoters() view returns(address[])
func (_Oracle0 *Oracle0Caller) GetVoters(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _Oracle0.contract.Call(opts, &out, "getVoters")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// GetVoters is a free data retrieval call binding the contract method 0xcdd72253.
//
// Solidity: function getVoters() view returns(address[])
func (_Oracle0 *Oracle0Session) GetVoters() ([]common.Address, error) {
	return _Oracle0.Contract.GetVoters(&_Oracle0.CallOpts)
}

// GetVoters is a free data retrieval call binding the contract method 0xcdd72253.
//
// Solidity: function getVoters() view returns(address[])
func (_Oracle0 *Oracle0CallerSession) GetVoters() ([]common.Address, error) {
	return _Oracle0.Contract.GetVoters(&_Oracle0.CallOpts)
}

// LatestRoundData is a free data retrieval call binding the contract method 0x33f98c77.
//
// Solidity: function latestRoundData(string _symbol) view returns((uint256,uint256,uint256,bool) data)
func (_Oracle0 *Oracle0Caller) LatestRoundData(opts *bind.CallOpts, _symbol string) (IOracleRoundData, error) {
	var out []interface{}
	err := _Oracle0.contract.Call(opts, &out, "latestRoundData", _symbol)

	if err != nil {
		return *new(IOracleRoundData), err
	}

	out0 := *abi.ConvertType(out[0], new(IOracleRoundData)).(*IOracleRoundData)

	return out0, err

}

// LatestRoundData is a free data retrieval call binding the contract method 0x33f98c77.
//
// Solidity: function latestRoundData(string _symbol) view returns((uint256,uint256,uint256,bool) data)
func (_Oracle0 *Oracle0Session) LatestRoundData(_symbol string) (IOracleRoundData, error) {
	return _Oracle0.Contract.LatestRoundData(&_Oracle0.CallOpts, _symbol)
}

// LatestRoundData is a free data retrieval call binding the contract method 0x33f98c77.
//
// Solidity: function latestRoundData(string _symbol) view returns((uint256,uint256,uint256,bool) data)
func (_Oracle0 *Oracle0CallerSession) LatestRoundData(_symbol string) (IOracleRoundData, error) {
	return _Oracle0.Contract.LatestRoundData(&_Oracle0.CallOpts, _symbol)
}

// DistributeRewards is a paid mutator transaction binding the contract method 0x59974e38.
//
// Solidity: function distributeRewards(uint256 _ntn) payable returns()
func (_Oracle0 *Oracle0Transactor) DistributeRewards(opts *bind.TransactOpts, _ntn *big.Int) (*types.Transaction, error) {
	return _Oracle0.contract.Transact(opts, "distributeRewards", _ntn)
}

// DistributeRewards is a paid mutator transaction binding the contract method 0x59974e38.
//
// Solidity: function distributeRewards(uint256 _ntn) payable returns()
func (_Oracle0 *Oracle0Session) DistributeRewards(_ntn *big.Int) (*types.Transaction, error) {
	return _Oracle0.Contract.DistributeRewards(&_Oracle0.TransactOpts, _ntn)
}

// DistributeRewards is a paid mutator transaction binding the contract method 0x59974e38.
//
// Solidity: function distributeRewards(uint256 _ntn) payable returns()
func (_Oracle0 *Oracle0TransactorSession) DistributeRewards(_ntn *big.Int) (*types.Transaction, error) {
	return _Oracle0.Contract.DistributeRewards(&_Oracle0.TransactOpts, _ntn)
}

// Finalize is a paid mutator transaction binding the contract method 0x4bb278f3.
//
// Solidity: function finalize() returns(bool)
func (_Oracle0 *Oracle0Transactor) Finalize(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Oracle0.contract.Transact(opts, "finalize")
}

// Finalize is a paid mutator transaction binding the contract method 0x4bb278f3.
//
// Solidity: function finalize() returns(bool)
func (_Oracle0 *Oracle0Session) Finalize() (*types.Transaction, error) {
	return _Oracle0.Contract.Finalize(&_Oracle0.TransactOpts)
}

// Finalize is a paid mutator transaction binding the contract method 0x4bb278f3.
//
// Solidity: function finalize() returns(bool)
func (_Oracle0 *Oracle0TransactorSession) Finalize() (*types.Transaction, error) {
	return _Oracle0.Contract.Finalize(&_Oracle0.TransactOpts)
}

// SetCommitRevealConfig is a paid mutator transaction binding the contract method 0x3f422ef3.
//
// Solidity: function setCommitRevealConfig(uint256 _threshold, uint256 _resetInterval) returns()
func (_Oracle0 *Oracle0Transactor) SetCommitRevealConfig(opts *bind.TransactOpts, _threshold *big.Int, _resetInterval *big.Int) (*types.Transaction, error) {
	return _Oracle0.contract.Transact(opts, "setCommitRevealConfig", _threshold, _resetInterval)
}

// SetCommitRevealConfig is a paid mutator transaction binding the contract method 0x3f422ef3.
//
// Solidity: function setCommitRevealConfig(uint256 _threshold, uint256 _resetInterval) returns()
func (_Oracle0 *Oracle0Session) SetCommitRevealConfig(_threshold *big.Int, _resetInterval *big.Int) (*types.Transaction, error) {
	return _Oracle0.Contract.SetCommitRevealConfig(&_Oracle0.TransactOpts, _threshold, _resetInterval)
}

// SetCommitRevealConfig is a paid mutator transaction binding the contract method 0x3f422ef3.
//
// Solidity: function setCommitRevealConfig(uint256 _threshold, uint256 _resetInterval) returns()
func (_Oracle0 *Oracle0TransactorSession) SetCommitRevealConfig(_threshold *big.Int, _resetInterval *big.Int) (*types.Transaction, error) {
	return _Oracle0.Contract.SetCommitRevealConfig(&_Oracle0.TransactOpts, _threshold, _resetInterval)
}

// SetOperator is a paid mutator transaction binding the contract method 0xb3ab15fb.
//
// Solidity: function setOperator(address _operator) returns()
func (_Oracle0 *Oracle0Transactor) SetOperator(opts *bind.TransactOpts, _operator common.Address) (*types.Transaction, error) {
	return _Oracle0.contract.Transact(opts, "setOperator", _operator)
}

// SetOperator is a paid mutator transaction binding the contract method 0xb3ab15fb.
//
// Solidity: function setOperator(address _operator) returns()
func (_Oracle0 *Oracle0Session) SetOperator(_operator common.Address) (*types.Transaction, error) {
	return _Oracle0.Contract.SetOperator(&_Oracle0.TransactOpts, _operator)
}

// SetOperator is a paid mutator transaction binding the contract method 0xb3ab15fb.
//
// Solidity: function setOperator(address _operator) returns()
func (_Oracle0 *Oracle0TransactorSession) SetOperator(_operator common.Address) (*types.Transaction, error) {
	return _Oracle0.Contract.SetOperator(&_Oracle0.TransactOpts, _operator)
}

// SetSlashingConfig is a paid mutator transaction binding the contract method 0xda39fbfe.
//
// Solidity: function setSlashingConfig(int256 _outlierSlashingThreshold, int256 _outlierDetectionThreshold, uint256 _baseSlashingRate, uint256 _slashingRateCap) returns()
func (_Oracle0 *Oracle0Transactor) SetSlashingConfig(opts *bind.TransactOpts, _outlierSlashingThreshold *big.Int, _outlierDetectionThreshold *big.Int, _baseSlashingRate *big.Int, _slashingRateCap *big.Int) (*types.Transaction, error) {
	return _Oracle0.contract.Transact(opts, "setSlashingConfig", _outlierSlashingThreshold, _outlierDetectionThreshold, _baseSlashingRate, _slashingRateCap)
}

// SetSlashingConfig is a paid mutator transaction binding the contract method 0xda39fbfe.
//
// Solidity: function setSlashingConfig(int256 _outlierSlashingThreshold, int256 _outlierDetectionThreshold, uint256 _baseSlashingRate, uint256 _slashingRateCap) returns()
func (_Oracle0 *Oracle0Session) SetSlashingConfig(_outlierSlashingThreshold *big.Int, _outlierDetectionThreshold *big.Int, _baseSlashingRate *big.Int, _slashingRateCap *big.Int) (*types.Transaction, error) {
	return _Oracle0.Contract.SetSlashingConfig(&_Oracle0.TransactOpts, _outlierSlashingThreshold, _outlierDetectionThreshold, _baseSlashingRate, _slashingRateCap)
}

// SetSlashingConfig is a paid mutator transaction binding the contract method 0xda39fbfe.
//
// Solidity: function setSlashingConfig(int256 _outlierSlashingThreshold, int256 _outlierDetectionThreshold, uint256 _baseSlashingRate, uint256 _slashingRateCap) returns()
func (_Oracle0 *Oracle0TransactorSession) SetSlashingConfig(_outlierSlashingThreshold *big.Int, _outlierDetectionThreshold *big.Int, _baseSlashingRate *big.Int, _slashingRateCap *big.Int) (*types.Transaction, error) {
	return _Oracle0.Contract.SetSlashingConfig(&_Oracle0.TransactOpts, _outlierSlashingThreshold, _outlierDetectionThreshold, _baseSlashingRate, _slashingRateCap)
}

// SetSymbols is a paid mutator transaction binding the contract method 0x8d4f75d2.
//
// Solidity: function setSymbols(string[] _symbols) returns()
func (_Oracle0 *Oracle0Transactor) SetSymbols(opts *bind.TransactOpts, _symbols []string) (*types.Transaction, error) {
	return _Oracle0.contract.Transact(opts, "setSymbols", _symbols)
}

// SetSymbols is a paid mutator transaction binding the contract method 0x8d4f75d2.
//
// Solidity: function setSymbols(string[] _symbols) returns()
func (_Oracle0 *Oracle0Session) SetSymbols(_symbols []string) (*types.Transaction, error) {
	return _Oracle0.Contract.SetSymbols(&_Oracle0.TransactOpts, _symbols)
}

// SetSymbols is a paid mutator transaction binding the contract method 0x8d4f75d2.
//
// Solidity: function setSymbols(string[] _symbols) returns()
func (_Oracle0 *Oracle0TransactorSession) SetSymbols(_symbols []string) (*types.Transaction, error) {
	return _Oracle0.Contract.SetSymbols(&_Oracle0.TransactOpts, _symbols)
}

// SetVotePeriod is a paid mutator transaction binding the contract method 0x67b11630.
//
// Solidity: function setVotePeriod(uint256 _votePeriod) returns()
func (_Oracle0 *Oracle0Transactor) SetVotePeriod(opts *bind.TransactOpts, _votePeriod *big.Int) (*types.Transaction, error) {
	return _Oracle0.contract.Transact(opts, "setVotePeriod", _votePeriod)
}

// SetVotePeriod is a paid mutator transaction binding the contract method 0x67b11630.
//
// Solidity: function setVotePeriod(uint256 _votePeriod) returns()
func (_Oracle0 *Oracle0Session) SetVotePeriod(_votePeriod *big.Int) (*types.Transaction, error) {
	return _Oracle0.Contract.SetVotePeriod(&_Oracle0.TransactOpts, _votePeriod)
}

// SetVotePeriod is a paid mutator transaction binding the contract method 0x67b11630.
//
// Solidity: function setVotePeriod(uint256 _votePeriod) returns()
func (_Oracle0 *Oracle0TransactorSession) SetVotePeriod(_votePeriod *big.Int) (*types.Transaction, error) {
	return _Oracle0.Contract.SetVotePeriod(&_Oracle0.TransactOpts, _votePeriod)
}

// SetVoters is a paid mutator transaction binding the contract method 0xda78110e.
//
// Solidity: function setVoters(address[] _newVoters, address[] _treasury, address[] _validator) returns()
func (_Oracle0 *Oracle0Transactor) SetVoters(opts *bind.TransactOpts, _newVoters []common.Address, _treasury []common.Address, _validator []common.Address) (*types.Transaction, error) {
	return _Oracle0.contract.Transact(opts, "setVoters", _newVoters, _treasury, _validator)
}

// SetVoters is a paid mutator transaction binding the contract method 0xda78110e.
//
// Solidity: function setVoters(address[] _newVoters, address[] _treasury, address[] _validator) returns()
func (_Oracle0 *Oracle0Session) SetVoters(_newVoters []common.Address, _treasury []common.Address, _validator []common.Address) (*types.Transaction, error) {
	return _Oracle0.Contract.SetVoters(&_Oracle0.TransactOpts, _newVoters, _treasury, _validator)
}

// SetVoters is a paid mutator transaction binding the contract method 0xda78110e.
//
// Solidity: function setVoters(address[] _newVoters, address[] _treasury, address[] _validator) returns()
func (_Oracle0 *Oracle0TransactorSession) SetVoters(_newVoters []common.Address, _treasury []common.Address, _validator []common.Address) (*types.Transaction, error) {
	return _Oracle0.Contract.SetVoters(&_Oracle0.TransactOpts, _newVoters, _treasury, _validator)
}

// UpdateVotersAndSymbol is a paid mutator transaction binding the contract method 0x0f65875c.
//
// Solidity: function updateVotersAndSymbol() returns()
func (_Oracle0 *Oracle0Transactor) UpdateVotersAndSymbol(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Oracle0.contract.Transact(opts, "updateVotersAndSymbol")
}

// UpdateVotersAndSymbol is a paid mutator transaction binding the contract method 0x0f65875c.
//
// Solidity: function updateVotersAndSymbol() returns()
func (_Oracle0 *Oracle0Session) UpdateVotersAndSymbol() (*types.Transaction, error) {
	return _Oracle0.Contract.UpdateVotersAndSymbol(&_Oracle0.TransactOpts)
}

// UpdateVotersAndSymbol is a paid mutator transaction binding the contract method 0x0f65875c.
//
// Solidity: function updateVotersAndSymbol() returns()
func (_Oracle0 *Oracle0TransactorSession) UpdateVotersAndSymbol() (*types.Transaction, error) {
	return _Oracle0.Contract.UpdateVotersAndSymbol(&_Oracle0.TransactOpts)
}

// Vote is a paid mutator transaction binding the contract method 0x56833ebe.
//
// Solidity: function vote(uint256 _commit, (uint120,uint8)[] _reports, uint256 _salt, uint8 _extra) returns()
func (_Oracle0 *Oracle0Transactor) Vote(opts *bind.TransactOpts, _commit *big.Int, _reports []IOracleReport, _salt *big.Int, _extra uint8) (*types.Transaction, error) {
	return _Oracle0.contract.Transact(opts, "vote", _commit, _reports, _salt, _extra)
}

// Vote is a paid mutator transaction binding the contract method 0x56833ebe.
//
// Solidity: function vote(uint256 _commit, (uint120,uint8)[] _reports, uint256 _salt, uint8 _extra) returns()
func (_Oracle0 *Oracle0Session) Vote(_commit *big.Int, _reports []IOracleReport, _salt *big.Int, _extra uint8) (*types.Transaction, error) {
	return _Oracle0.Contract.Vote(&_Oracle0.TransactOpts, _commit, _reports, _salt, _extra)
}

// Vote is a paid mutator transaction binding the contract method 0x56833ebe.
//
// Solidity: function vote(uint256 _commit, (uint120,uint8)[] _reports, uint256 _salt, uint8 _extra) returns()
func (_Oracle0 *Oracle0TransactorSession) Vote(_commit *big.Int, _reports []IOracleReport, _salt *big.Int, _extra uint8) (*types.Transaction, error) {
	return _Oracle0.Contract.Vote(&_Oracle0.TransactOpts, _commit, _reports, _salt, _extra)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Oracle0 *Oracle0Transactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _Oracle0.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Oracle0 *Oracle0Session) Fallback(calldata []byte) (*types.Transaction, error) {
	return _Oracle0.Contract.Fallback(&_Oracle0.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Oracle0 *Oracle0TransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _Oracle0.Contract.Fallback(&_Oracle0.TransactOpts, calldata)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Oracle0 *Oracle0Transactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Oracle0.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Oracle0 *Oracle0Session) Receive() (*types.Transaction, error) {
	return _Oracle0.Contract.Receive(&_Oracle0.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Oracle0 *Oracle0TransactorSession) Receive() (*types.Transaction, error) {
	return _Oracle0.Contract.Receive(&_Oracle0.TransactOpts)
}

// Oracle0CallFailedIterator is returned from FilterCallFailed and is used to iterate over the raw logs and unpacked data for CallFailed events raised by the Oracle0 contract.
type Oracle0CallFailedIterator struct {
	Event *Oracle0CallFailed // Event containing the contract specifics and raw log

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
func (it *Oracle0CallFailedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Oracle0CallFailed)
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
		it.Event = new(Oracle0CallFailed)
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
func (it *Oracle0CallFailedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Oracle0CallFailedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Oracle0CallFailed represents a CallFailed event raised by the Oracle0 contract.
type Oracle0CallFailed struct {
	To              common.Address
	MethodSignature string
	ReturnData      []byte
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterCallFailed is a free log retrieval operation binding the contract event 0x1137d8c966ce69b9630fb2294be011f3d64cc56e91fad7d375f0662568e9d352.
//
// Solidity: event CallFailed(address to, string methodSignature, bytes returnData)
func (_Oracle0 *Oracle0Filterer) FilterCallFailed(opts *bind.FilterOpts) (*Oracle0CallFailedIterator, error) {

	logs, sub, err := _Oracle0.contract.FilterLogs(opts, "CallFailed")
	if err != nil {
		return nil, err
	}
	return &Oracle0CallFailedIterator{contract: _Oracle0.contract, event: "CallFailed", logs: logs, sub: sub}, nil
}

// WatchCallFailed is a free log subscription operation binding the contract event 0x1137d8c966ce69b9630fb2294be011f3d64cc56e91fad7d375f0662568e9d352.
//
// Solidity: event CallFailed(address to, string methodSignature, bytes returnData)
func (_Oracle0 *Oracle0Filterer) WatchCallFailed(opts *bind.WatchOpts, sink chan<- *Oracle0CallFailed) (event.Subscription, error) {

	logs, sub, err := _Oracle0.contract.WatchLogs(opts, "CallFailed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Oracle0CallFailed)
				if err := _Oracle0.contract.UnpackLog(event, "CallFailed", log); err != nil {
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

// ParseCallFailed is a log parse operation binding the contract event 0x1137d8c966ce69b9630fb2294be011f3d64cc56e91fad7d375f0662568e9d352.
//
// Solidity: event CallFailed(address to, string methodSignature, bytes returnData)
func (_Oracle0 *Oracle0Filterer) ParseCallFailed(log types.Log) (*Oracle0CallFailed, error) {
	event := new(Oracle0CallFailed)
	if err := _Oracle0.contract.UnpackLog(event, "CallFailed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Oracle0CommitRevealMissedIterator is returned from FilterCommitRevealMissed and is used to iterate over the raw logs and unpacked data for CommitRevealMissed events raised by the Oracle0 contract.
type Oracle0CommitRevealMissedIterator struct {
	Event *Oracle0CommitRevealMissed // Event containing the contract specifics and raw log

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
func (it *Oracle0CommitRevealMissedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Oracle0CommitRevealMissed)
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
		it.Event = new(Oracle0CommitRevealMissed)
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
func (it *Oracle0CommitRevealMissedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Oracle0CommitRevealMissedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Oracle0CommitRevealMissed represents a CommitRevealMissed event raised by the Oracle0 contract.
type Oracle0CommitRevealMissed struct {
	Voter          common.Address
	Round          *big.Int
	NonRevealCount *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterCommitRevealMissed is a free log retrieval operation binding the contract event 0x176956a4e941f6737f81a3c9a09d8571dd0438d86e25a432beb2013aced43092.
//
// Solidity: event CommitRevealMissed(address indexed _voter, uint256 _round, uint256 _nonRevealCount)
func (_Oracle0 *Oracle0Filterer) FilterCommitRevealMissed(opts *bind.FilterOpts, _voter []common.Address) (*Oracle0CommitRevealMissedIterator, error) {

	var _voterRule []interface{}
	for _, _voterItem := range _voter {
		_voterRule = append(_voterRule, _voterItem)
	}

	logs, sub, err := _Oracle0.contract.FilterLogs(opts, "CommitRevealMissed", _voterRule)
	if err != nil {
		return nil, err
	}
	return &Oracle0CommitRevealMissedIterator{contract: _Oracle0.contract, event: "CommitRevealMissed", logs: logs, sub: sub}, nil
}

// WatchCommitRevealMissed is a free log subscription operation binding the contract event 0x176956a4e941f6737f81a3c9a09d8571dd0438d86e25a432beb2013aced43092.
//
// Solidity: event CommitRevealMissed(address indexed _voter, uint256 _round, uint256 _nonRevealCount)
func (_Oracle0 *Oracle0Filterer) WatchCommitRevealMissed(opts *bind.WatchOpts, sink chan<- *Oracle0CommitRevealMissed, _voter []common.Address) (event.Subscription, error) {

	var _voterRule []interface{}
	for _, _voterItem := range _voter {
		_voterRule = append(_voterRule, _voterItem)
	}

	logs, sub, err := _Oracle0.contract.WatchLogs(opts, "CommitRevealMissed", _voterRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Oracle0CommitRevealMissed)
				if err := _Oracle0.contract.UnpackLog(event, "CommitRevealMissed", log); err != nil {
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

// ParseCommitRevealMissed is a log parse operation binding the contract event 0x176956a4e941f6737f81a3c9a09d8571dd0438d86e25a432beb2013aced43092.
//
// Solidity: event CommitRevealMissed(address indexed _voter, uint256 _round, uint256 _nonRevealCount)
func (_Oracle0 *Oracle0Filterer) ParseCommitRevealMissed(log types.Log) (*Oracle0CommitRevealMissed, error) {
	event := new(Oracle0CommitRevealMissed)
	if err := _Oracle0.contract.UnpackLog(event, "CommitRevealMissed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Oracle0ConfigUpdateAddressIterator is returned from FilterConfigUpdateAddress and is used to iterate over the raw logs and unpacked data for ConfigUpdateAddress events raised by the Oracle0 contract.
type Oracle0ConfigUpdateAddressIterator struct {
	Event *Oracle0ConfigUpdateAddress // Event containing the contract specifics and raw log

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
func (it *Oracle0ConfigUpdateAddressIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Oracle0ConfigUpdateAddress)
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
		it.Event = new(Oracle0ConfigUpdateAddress)
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
func (it *Oracle0ConfigUpdateAddressIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Oracle0ConfigUpdateAddressIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Oracle0ConfigUpdateAddress represents a ConfigUpdateAddress event raised by the Oracle0 contract.
type Oracle0ConfigUpdateAddress struct {
	Name            string
	OldValue        common.Address
	NewValue        common.Address
	AppliesAtHeight *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterConfigUpdateAddress is a free log retrieval operation binding the contract event 0xe821ac8084a7329d09d00cf1380cba50edea5f54a7bd453d49a9c354f565d4a0.
//
// Solidity: event ConfigUpdateAddress(string name, address oldValue, address newValue, uint256 appliesAtHeight)
func (_Oracle0 *Oracle0Filterer) FilterConfigUpdateAddress(opts *bind.FilterOpts) (*Oracle0ConfigUpdateAddressIterator, error) {

	logs, sub, err := _Oracle0.contract.FilterLogs(opts, "ConfigUpdateAddress")
	if err != nil {
		return nil, err
	}
	return &Oracle0ConfigUpdateAddressIterator{contract: _Oracle0.contract, event: "ConfigUpdateAddress", logs: logs, sub: sub}, nil
}

// WatchConfigUpdateAddress is a free log subscription operation binding the contract event 0xe821ac8084a7329d09d00cf1380cba50edea5f54a7bd453d49a9c354f565d4a0.
//
// Solidity: event ConfigUpdateAddress(string name, address oldValue, address newValue, uint256 appliesAtHeight)
func (_Oracle0 *Oracle0Filterer) WatchConfigUpdateAddress(opts *bind.WatchOpts, sink chan<- *Oracle0ConfigUpdateAddress) (event.Subscription, error) {

	logs, sub, err := _Oracle0.contract.WatchLogs(opts, "ConfigUpdateAddress")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Oracle0ConfigUpdateAddress)
				if err := _Oracle0.contract.UnpackLog(event, "ConfigUpdateAddress", log); err != nil {
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

// ParseConfigUpdateAddress is a log parse operation binding the contract event 0xe821ac8084a7329d09d00cf1380cba50edea5f54a7bd453d49a9c354f565d4a0.
//
// Solidity: event ConfigUpdateAddress(string name, address oldValue, address newValue, uint256 appliesAtHeight)
func (_Oracle0 *Oracle0Filterer) ParseConfigUpdateAddress(log types.Log) (*Oracle0ConfigUpdateAddress, error) {
	event := new(Oracle0ConfigUpdateAddress)
	if err := _Oracle0.contract.UnpackLog(event, "ConfigUpdateAddress", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Oracle0ConfigUpdateBoolIterator is returned from FilterConfigUpdateBool and is used to iterate over the raw logs and unpacked data for ConfigUpdateBool events raised by the Oracle0 contract.
type Oracle0ConfigUpdateBoolIterator struct {
	Event *Oracle0ConfigUpdateBool // Event containing the contract specifics and raw log

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
func (it *Oracle0ConfigUpdateBoolIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Oracle0ConfigUpdateBool)
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
		it.Event = new(Oracle0ConfigUpdateBool)
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
func (it *Oracle0ConfigUpdateBoolIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Oracle0ConfigUpdateBoolIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Oracle0ConfigUpdateBool represents a ConfigUpdateBool event raised by the Oracle0 contract.
type Oracle0ConfigUpdateBool struct {
	Name            string
	OldValue        bool
	NewValue        bool
	AppliesAtHeight *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterConfigUpdateBool is a free log retrieval operation binding the contract event 0x5edb308c5eddc69bcd31b4e689c5eed2fbd3155ae57915d1cad05425f6c1a39b.
//
// Solidity: event ConfigUpdateBool(string name, bool oldValue, bool newValue, uint256 appliesAtHeight)
func (_Oracle0 *Oracle0Filterer) FilterConfigUpdateBool(opts *bind.FilterOpts) (*Oracle0ConfigUpdateBoolIterator, error) {

	logs, sub, err := _Oracle0.contract.FilterLogs(opts, "ConfigUpdateBool")
	if err != nil {
		return nil, err
	}
	return &Oracle0ConfigUpdateBoolIterator{contract: _Oracle0.contract, event: "ConfigUpdateBool", logs: logs, sub: sub}, nil
}

// WatchConfigUpdateBool is a free log subscription operation binding the contract event 0x5edb308c5eddc69bcd31b4e689c5eed2fbd3155ae57915d1cad05425f6c1a39b.
//
// Solidity: event ConfigUpdateBool(string name, bool oldValue, bool newValue, uint256 appliesAtHeight)
func (_Oracle0 *Oracle0Filterer) WatchConfigUpdateBool(opts *bind.WatchOpts, sink chan<- *Oracle0ConfigUpdateBool) (event.Subscription, error) {

	logs, sub, err := _Oracle0.contract.WatchLogs(opts, "ConfigUpdateBool")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Oracle0ConfigUpdateBool)
				if err := _Oracle0.contract.UnpackLog(event, "ConfigUpdateBool", log); err != nil {
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

// ParseConfigUpdateBool is a log parse operation binding the contract event 0x5edb308c5eddc69bcd31b4e689c5eed2fbd3155ae57915d1cad05425f6c1a39b.
//
// Solidity: event ConfigUpdateBool(string name, bool oldValue, bool newValue, uint256 appliesAtHeight)
func (_Oracle0 *Oracle0Filterer) ParseConfigUpdateBool(log types.Log) (*Oracle0ConfigUpdateBool, error) {
	event := new(Oracle0ConfigUpdateBool)
	if err := _Oracle0.contract.UnpackLog(event, "ConfigUpdateBool", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Oracle0ConfigUpdateIntIterator is returned from FilterConfigUpdateInt and is used to iterate over the raw logs and unpacked data for ConfigUpdateInt events raised by the Oracle0 contract.
type Oracle0ConfigUpdateIntIterator struct {
	Event *Oracle0ConfigUpdateInt // Event containing the contract specifics and raw log

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
func (it *Oracle0ConfigUpdateIntIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Oracle0ConfigUpdateInt)
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
		it.Event = new(Oracle0ConfigUpdateInt)
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
func (it *Oracle0ConfigUpdateIntIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Oracle0ConfigUpdateIntIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Oracle0ConfigUpdateInt represents a ConfigUpdateInt event raised by the Oracle0 contract.
type Oracle0ConfigUpdateInt struct {
	Name            string
	OldValue        *big.Int
	NewValue        *big.Int
	AppliesAtHeight *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterConfigUpdateInt is a free log retrieval operation binding the contract event 0xb5114472b89d1126433287defb6308bbabdb95b6ce5dd949ab2c151ed7b1ff4c.
//
// Solidity: event ConfigUpdateInt(string name, int256 oldValue, int256 newValue, uint256 appliesAtHeight)
func (_Oracle0 *Oracle0Filterer) FilterConfigUpdateInt(opts *bind.FilterOpts) (*Oracle0ConfigUpdateIntIterator, error) {

	logs, sub, err := _Oracle0.contract.FilterLogs(opts, "ConfigUpdateInt")
	if err != nil {
		return nil, err
	}
	return &Oracle0ConfigUpdateIntIterator{contract: _Oracle0.contract, event: "ConfigUpdateInt", logs: logs, sub: sub}, nil
}

// WatchConfigUpdateInt is a free log subscription operation binding the contract event 0xb5114472b89d1126433287defb6308bbabdb95b6ce5dd949ab2c151ed7b1ff4c.
//
// Solidity: event ConfigUpdateInt(string name, int256 oldValue, int256 newValue, uint256 appliesAtHeight)
func (_Oracle0 *Oracle0Filterer) WatchConfigUpdateInt(opts *bind.WatchOpts, sink chan<- *Oracle0ConfigUpdateInt) (event.Subscription, error) {

	logs, sub, err := _Oracle0.contract.WatchLogs(opts, "ConfigUpdateInt")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Oracle0ConfigUpdateInt)
				if err := _Oracle0.contract.UnpackLog(event, "ConfigUpdateInt", log); err != nil {
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

// ParseConfigUpdateInt is a log parse operation binding the contract event 0xb5114472b89d1126433287defb6308bbabdb95b6ce5dd949ab2c151ed7b1ff4c.
//
// Solidity: event ConfigUpdateInt(string name, int256 oldValue, int256 newValue, uint256 appliesAtHeight)
func (_Oracle0 *Oracle0Filterer) ParseConfigUpdateInt(log types.Log) (*Oracle0ConfigUpdateInt, error) {
	event := new(Oracle0ConfigUpdateInt)
	if err := _Oracle0.contract.UnpackLog(event, "ConfigUpdateInt", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Oracle0ConfigUpdateUintIterator is returned from FilterConfigUpdateUint and is used to iterate over the raw logs and unpacked data for ConfigUpdateUint events raised by the Oracle0 contract.
type Oracle0ConfigUpdateUintIterator struct {
	Event *Oracle0ConfigUpdateUint // Event containing the contract specifics and raw log

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
func (it *Oracle0ConfigUpdateUintIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Oracle0ConfigUpdateUint)
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
		it.Event = new(Oracle0ConfigUpdateUint)
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
func (it *Oracle0ConfigUpdateUintIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Oracle0ConfigUpdateUintIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Oracle0ConfigUpdateUint represents a ConfigUpdateUint event raised by the Oracle0 contract.
type Oracle0ConfigUpdateUint struct {
	Name            string
	OldValue        *big.Int
	NewValue        *big.Int
	AppliesAtHeight *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterConfigUpdateUint is a free log retrieval operation binding the contract event 0x207e45ce6f2191c3efffe5d0d91591cf81afd23cbf4bd746981d7a8a8dcfe1ba.
//
// Solidity: event ConfigUpdateUint(string name, uint256 oldValue, uint256 newValue, uint256 appliesAtHeight)
func (_Oracle0 *Oracle0Filterer) FilterConfigUpdateUint(opts *bind.FilterOpts) (*Oracle0ConfigUpdateUintIterator, error) {

	logs, sub, err := _Oracle0.contract.FilterLogs(opts, "ConfigUpdateUint")
	if err != nil {
		return nil, err
	}
	return &Oracle0ConfigUpdateUintIterator{contract: _Oracle0.contract, event: "ConfigUpdateUint", logs: logs, sub: sub}, nil
}

// WatchConfigUpdateUint is a free log subscription operation binding the contract event 0x207e45ce6f2191c3efffe5d0d91591cf81afd23cbf4bd746981d7a8a8dcfe1ba.
//
// Solidity: event ConfigUpdateUint(string name, uint256 oldValue, uint256 newValue, uint256 appliesAtHeight)
func (_Oracle0 *Oracle0Filterer) WatchConfigUpdateUint(opts *bind.WatchOpts, sink chan<- *Oracle0ConfigUpdateUint) (event.Subscription, error) {

	logs, sub, err := _Oracle0.contract.WatchLogs(opts, "ConfigUpdateUint")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Oracle0ConfigUpdateUint)
				if err := _Oracle0.contract.UnpackLog(event, "ConfigUpdateUint", log); err != nil {
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

// ParseConfigUpdateUint is a log parse operation binding the contract event 0x207e45ce6f2191c3efffe5d0d91591cf81afd23cbf4bd746981d7a8a8dcfe1ba.
//
// Solidity: event ConfigUpdateUint(string name, uint256 oldValue, uint256 newValue, uint256 appliesAtHeight)
func (_Oracle0 *Oracle0Filterer) ParseConfigUpdateUint(log types.Log) (*Oracle0ConfigUpdateUint, error) {
	event := new(Oracle0ConfigUpdateUint)
	if err := _Oracle0.contract.UnpackLog(event, "ConfigUpdateUint", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Oracle0InvalidVoteIterator is returned from FilterInvalidVote and is used to iterate over the raw logs and unpacked data for InvalidVote events raised by the Oracle0 contract.
type Oracle0InvalidVoteIterator struct {
	Event *Oracle0InvalidVote // Event containing the contract specifics and raw log

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
func (it *Oracle0InvalidVoteIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Oracle0InvalidVote)
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
		it.Event = new(Oracle0InvalidVote)
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
func (it *Oracle0InvalidVoteIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Oracle0InvalidVoteIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Oracle0InvalidVote represents a InvalidVote event raised by the Oracle0 contract.
type Oracle0InvalidVote struct {
	Cause       string
	Reporter    common.Address
	ExpValue    *big.Int
	ActualValue *big.Int
	Extra       uint8
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterInvalidVote is a free log retrieval operation binding the contract event 0x04ca4e0efda95f8b780c116574d1521309010b38d8f7b75705495703a0f570b1.
//
// Solidity: event InvalidVote(string cause, address indexed reporter, uint256 expValue, uint256 actualValue, uint8 extra)
func (_Oracle0 *Oracle0Filterer) FilterInvalidVote(opts *bind.FilterOpts, reporter []common.Address) (*Oracle0InvalidVoteIterator, error) {

	var reporterRule []interface{}
	for _, reporterItem := range reporter {
		reporterRule = append(reporterRule, reporterItem)
	}

	logs, sub, err := _Oracle0.contract.FilterLogs(opts, "InvalidVote", reporterRule)
	if err != nil {
		return nil, err
	}
	return &Oracle0InvalidVoteIterator{contract: _Oracle0.contract, event: "InvalidVote", logs: logs, sub: sub}, nil
}

// WatchInvalidVote is a free log subscription operation binding the contract event 0x04ca4e0efda95f8b780c116574d1521309010b38d8f7b75705495703a0f570b1.
//
// Solidity: event InvalidVote(string cause, address indexed reporter, uint256 expValue, uint256 actualValue, uint8 extra)
func (_Oracle0 *Oracle0Filterer) WatchInvalidVote(opts *bind.WatchOpts, sink chan<- *Oracle0InvalidVote, reporter []common.Address) (event.Subscription, error) {

	var reporterRule []interface{}
	for _, reporterItem := range reporter {
		reporterRule = append(reporterRule, reporterItem)
	}

	logs, sub, err := _Oracle0.contract.WatchLogs(opts, "InvalidVote", reporterRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Oracle0InvalidVote)
				if err := _Oracle0.contract.UnpackLog(event, "InvalidVote", log); err != nil {
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

// ParseInvalidVote is a log parse operation binding the contract event 0x04ca4e0efda95f8b780c116574d1521309010b38d8f7b75705495703a0f570b1.
//
// Solidity: event InvalidVote(string cause, address indexed reporter, uint256 expValue, uint256 actualValue, uint8 extra)
func (_Oracle0 *Oracle0Filterer) ParseInvalidVote(log types.Log) (*Oracle0InvalidVote, error) {
	event := new(Oracle0InvalidVote)
	if err := _Oracle0.contract.UnpackLog(event, "InvalidVote", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Oracle0NewRoundIterator is returned from FilterNewRound and is used to iterate over the raw logs and unpacked data for NewRound events raised by the Oracle0 contract.
type Oracle0NewRoundIterator struct {
	Event *Oracle0NewRound // Event containing the contract specifics and raw log

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
func (it *Oracle0NewRoundIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Oracle0NewRound)
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
		it.Event = new(Oracle0NewRound)
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
func (it *Oracle0NewRoundIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Oracle0NewRoundIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Oracle0NewRound represents a NewRound event raised by the Oracle0 contract.
type Oracle0NewRound struct {
	Round      *big.Int
	Timestamp  *big.Int
	VotePeriod *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterNewRound is a free log retrieval operation binding the contract event 0x5aec57d81928b24d30b1a2aec0d23d693412c37d7ec106b5d8259413716bb1f4.
//
// Solidity: event NewRound(uint256 _round, uint256 _timestamp, uint256 _votePeriod)
func (_Oracle0 *Oracle0Filterer) FilterNewRound(opts *bind.FilterOpts) (*Oracle0NewRoundIterator, error) {

	logs, sub, err := _Oracle0.contract.FilterLogs(opts, "NewRound")
	if err != nil {
		return nil, err
	}
	return &Oracle0NewRoundIterator{contract: _Oracle0.contract, event: "NewRound", logs: logs, sub: sub}, nil
}

// WatchNewRound is a free log subscription operation binding the contract event 0x5aec57d81928b24d30b1a2aec0d23d693412c37d7ec106b5d8259413716bb1f4.
//
// Solidity: event NewRound(uint256 _round, uint256 _timestamp, uint256 _votePeriod)
func (_Oracle0 *Oracle0Filterer) WatchNewRound(opts *bind.WatchOpts, sink chan<- *Oracle0NewRound) (event.Subscription, error) {

	logs, sub, err := _Oracle0.contract.WatchLogs(opts, "NewRound")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Oracle0NewRound)
				if err := _Oracle0.contract.UnpackLog(event, "NewRound", log); err != nil {
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

// ParseNewRound is a log parse operation binding the contract event 0x5aec57d81928b24d30b1a2aec0d23d693412c37d7ec106b5d8259413716bb1f4.
//
// Solidity: event NewRound(uint256 _round, uint256 _timestamp, uint256 _votePeriod)
func (_Oracle0 *Oracle0Filterer) ParseNewRound(log types.Log) (*Oracle0NewRound, error) {
	event := new(Oracle0NewRound)
	if err := _Oracle0.contract.UnpackLog(event, "NewRound", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Oracle0NewSymbolsIterator is returned from FilterNewSymbols and is used to iterate over the raw logs and unpacked data for NewSymbols events raised by the Oracle0 contract.
type Oracle0NewSymbolsIterator struct {
	Event *Oracle0NewSymbols // Event containing the contract specifics and raw log

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
func (it *Oracle0NewSymbolsIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Oracle0NewSymbols)
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
		it.Event = new(Oracle0NewSymbols)
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
func (it *Oracle0NewSymbolsIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Oracle0NewSymbolsIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Oracle0NewSymbols represents a NewSymbols event raised by the Oracle0 contract.
type Oracle0NewSymbols struct {
	Symbols []string
	Round   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterNewSymbols is a free log retrieval operation binding the contract event 0xaa278e424da680ce5dad66510415760e78e0bd87d45c786c6e88bdde82f9342d.
//
// Solidity: event NewSymbols(string[] _symbols, uint256 _round)
func (_Oracle0 *Oracle0Filterer) FilterNewSymbols(opts *bind.FilterOpts) (*Oracle0NewSymbolsIterator, error) {

	logs, sub, err := _Oracle0.contract.FilterLogs(opts, "NewSymbols")
	if err != nil {
		return nil, err
	}
	return &Oracle0NewSymbolsIterator{contract: _Oracle0.contract, event: "NewSymbols", logs: logs, sub: sub}, nil
}

// WatchNewSymbols is a free log subscription operation binding the contract event 0xaa278e424da680ce5dad66510415760e78e0bd87d45c786c6e88bdde82f9342d.
//
// Solidity: event NewSymbols(string[] _symbols, uint256 _round)
func (_Oracle0 *Oracle0Filterer) WatchNewSymbols(opts *bind.WatchOpts, sink chan<- *Oracle0NewSymbols) (event.Subscription, error) {

	logs, sub, err := _Oracle0.contract.WatchLogs(opts, "NewSymbols")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Oracle0NewSymbols)
				if err := _Oracle0.contract.UnpackLog(event, "NewSymbols", log); err != nil {
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

// ParseNewSymbols is a log parse operation binding the contract event 0xaa278e424da680ce5dad66510415760e78e0bd87d45c786c6e88bdde82f9342d.
//
// Solidity: event NewSymbols(string[] _symbols, uint256 _round)
func (_Oracle0 *Oracle0Filterer) ParseNewSymbols(log types.Log) (*Oracle0NewSymbols, error) {
	event := new(Oracle0NewSymbols)
	if err := _Oracle0.contract.UnpackLog(event, "NewSymbols", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Oracle0NewVoterIterator is returned from FilterNewVoter and is used to iterate over the raw logs and unpacked data for NewVoter events raised by the Oracle0 contract.
type Oracle0NewVoterIterator struct {
	Event *Oracle0NewVoter // Event containing the contract specifics and raw log

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
func (it *Oracle0NewVoterIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Oracle0NewVoter)
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
		it.Event = new(Oracle0NewVoter)
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
func (it *Oracle0NewVoterIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Oracle0NewVoterIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Oracle0NewVoter represents a NewVoter event raised by the Oracle0 contract.
type Oracle0NewVoter struct {
	Reporter common.Address
	Extra    uint8
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterNewVoter is a free log retrieval operation binding the contract event 0xd2ec8e890a03083998d3e16f98044fd3dd13fe3e61b7bc2e58ee6da43b50af73.
//
// Solidity: event NewVoter(address reporter, uint8 extra)
func (_Oracle0 *Oracle0Filterer) FilterNewVoter(opts *bind.FilterOpts) (*Oracle0NewVoterIterator, error) {

	logs, sub, err := _Oracle0.contract.FilterLogs(opts, "NewVoter")
	if err != nil {
		return nil, err
	}
	return &Oracle0NewVoterIterator{contract: _Oracle0.contract, event: "NewVoter", logs: logs, sub: sub}, nil
}

// WatchNewVoter is a free log subscription operation binding the contract event 0xd2ec8e890a03083998d3e16f98044fd3dd13fe3e61b7bc2e58ee6da43b50af73.
//
// Solidity: event NewVoter(address reporter, uint8 extra)
func (_Oracle0 *Oracle0Filterer) WatchNewVoter(opts *bind.WatchOpts, sink chan<- *Oracle0NewVoter) (event.Subscription, error) {

	logs, sub, err := _Oracle0.contract.WatchLogs(opts, "NewVoter")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Oracle0NewVoter)
				if err := _Oracle0.contract.UnpackLog(event, "NewVoter", log); err != nil {
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

// ParseNewVoter is a log parse operation binding the contract event 0xd2ec8e890a03083998d3e16f98044fd3dd13fe3e61b7bc2e58ee6da43b50af73.
//
// Solidity: event NewVoter(address reporter, uint8 extra)
func (_Oracle0 *Oracle0Filterer) ParseNewVoter(log types.Log) (*Oracle0NewVoter, error) {
	event := new(Oracle0NewVoter)
	if err := _Oracle0.contract.UnpackLog(event, "NewVoter", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Oracle0NoRevealPenaltyIterator is returned from FilterNoRevealPenalty and is used to iterate over the raw logs and unpacked data for NoRevealPenalty events raised by the Oracle0 contract.
type Oracle0NoRevealPenaltyIterator struct {
	Event *Oracle0NoRevealPenalty // Event containing the contract specifics and raw log

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
func (it *Oracle0NoRevealPenaltyIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Oracle0NoRevealPenalty)
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
		it.Event = new(Oracle0NoRevealPenalty)
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
func (it *Oracle0NoRevealPenaltyIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Oracle0NoRevealPenaltyIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Oracle0NoRevealPenalty represents a NoRevealPenalty event raised by the Oracle0 contract.
type Oracle0NoRevealPenalty struct {
	Voter        common.Address
	Round        *big.Int
	MissedReveal *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterNoRevealPenalty is a free log retrieval operation binding the contract event 0x9e6b40f10c60d1ad09594f3b6ed7043d0e978f584d354ace6e1f6025660c42b1.
//
// Solidity: event NoRevealPenalty(address indexed _voter, uint256 _round, uint256 _missedReveal)
func (_Oracle0 *Oracle0Filterer) FilterNoRevealPenalty(opts *bind.FilterOpts, _voter []common.Address) (*Oracle0NoRevealPenaltyIterator, error) {

	var _voterRule []interface{}
	for _, _voterItem := range _voter {
		_voterRule = append(_voterRule, _voterItem)
	}

	logs, sub, err := _Oracle0.contract.FilterLogs(opts, "NoRevealPenalty", _voterRule)
	if err != nil {
		return nil, err
	}
	return &Oracle0NoRevealPenaltyIterator{contract: _Oracle0.contract, event: "NoRevealPenalty", logs: logs, sub: sub}, nil
}

// WatchNoRevealPenalty is a free log subscription operation binding the contract event 0x9e6b40f10c60d1ad09594f3b6ed7043d0e978f584d354ace6e1f6025660c42b1.
//
// Solidity: event NoRevealPenalty(address indexed _voter, uint256 _round, uint256 _missedReveal)
func (_Oracle0 *Oracle0Filterer) WatchNoRevealPenalty(opts *bind.WatchOpts, sink chan<- *Oracle0NoRevealPenalty, _voter []common.Address) (event.Subscription, error) {

	var _voterRule []interface{}
	for _, _voterItem := range _voter {
		_voterRule = append(_voterRule, _voterItem)
	}

	logs, sub, err := _Oracle0.contract.WatchLogs(opts, "NoRevealPenalty", _voterRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Oracle0NoRevealPenalty)
				if err := _Oracle0.contract.UnpackLog(event, "NoRevealPenalty", log); err != nil {
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

// ParseNoRevealPenalty is a log parse operation binding the contract event 0x9e6b40f10c60d1ad09594f3b6ed7043d0e978f584d354ace6e1f6025660c42b1.
//
// Solidity: event NoRevealPenalty(address indexed _voter, uint256 _round, uint256 _missedReveal)
func (_Oracle0 *Oracle0Filterer) ParseNoRevealPenalty(log types.Log) (*Oracle0NoRevealPenalty, error) {
	event := new(Oracle0NoRevealPenalty)
	if err := _Oracle0.contract.UnpackLog(event, "NoRevealPenalty", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Oracle0PenalizedIterator is returned from FilterPenalized and is used to iterate over the raw logs and unpacked data for Penalized events raised by the Oracle0 contract.
type Oracle0PenalizedIterator struct {
	Event *Oracle0Penalized // Event containing the contract specifics and raw log

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
func (it *Oracle0PenalizedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Oracle0Penalized)
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
		it.Event = new(Oracle0Penalized)
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
func (it *Oracle0PenalizedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Oracle0PenalizedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Oracle0Penalized represents a Penalized event raised by the Oracle0 contract.
type Oracle0Penalized struct {
	Participant    common.Address
	SlashingAmount *big.Int
	Symbol         string
	Median         *big.Int
	Reported       *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterPenalized is a free log retrieval operation binding the contract event 0x372858b237c8bd0714183e8351a461d6c3cb1ef83806181b36bf5943711f4f57.
//
// Solidity: event Penalized(address indexed _participant, uint256 _slashingAmount, string _symbol, int256 _median, uint120 _reported)
func (_Oracle0 *Oracle0Filterer) FilterPenalized(opts *bind.FilterOpts, _participant []common.Address) (*Oracle0PenalizedIterator, error) {

	var _participantRule []interface{}
	for _, _participantItem := range _participant {
		_participantRule = append(_participantRule, _participantItem)
	}

	logs, sub, err := _Oracle0.contract.FilterLogs(opts, "Penalized", _participantRule)
	if err != nil {
		return nil, err
	}
	return &Oracle0PenalizedIterator{contract: _Oracle0.contract, event: "Penalized", logs: logs, sub: sub}, nil
}

// WatchPenalized is a free log subscription operation binding the contract event 0x372858b237c8bd0714183e8351a461d6c3cb1ef83806181b36bf5943711f4f57.
//
// Solidity: event Penalized(address indexed _participant, uint256 _slashingAmount, string _symbol, int256 _median, uint120 _reported)
func (_Oracle0 *Oracle0Filterer) WatchPenalized(opts *bind.WatchOpts, sink chan<- *Oracle0Penalized, _participant []common.Address) (event.Subscription, error) {

	var _participantRule []interface{}
	for _, _participantItem := range _participant {
		_participantRule = append(_participantRule, _participantItem)
	}

	logs, sub, err := _Oracle0.contract.WatchLogs(opts, "Penalized", _participantRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Oracle0Penalized)
				if err := _Oracle0.contract.UnpackLog(event, "Penalized", log); err != nil {
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

// ParsePenalized is a log parse operation binding the contract event 0x372858b237c8bd0714183e8351a461d6c3cb1ef83806181b36bf5943711f4f57.
//
// Solidity: event Penalized(address indexed _participant, uint256 _slashingAmount, string _symbol, int256 _median, uint120 _reported)
func (_Oracle0 *Oracle0Filterer) ParsePenalized(log types.Log) (*Oracle0Penalized, error) {
	event := new(Oracle0Penalized)
	if err := _Oracle0.contract.UnpackLog(event, "Penalized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Oracle0PriceUpdatedIterator is returned from FilterPriceUpdated and is used to iterate over the raw logs and unpacked data for PriceUpdated events raised by the Oracle0 contract.
type Oracle0PriceUpdatedIterator struct {
	Event *Oracle0PriceUpdated // Event containing the contract specifics and raw log

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
func (it *Oracle0PriceUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Oracle0PriceUpdated)
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
		it.Event = new(Oracle0PriceUpdated)
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
func (it *Oracle0PriceUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Oracle0PriceUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Oracle0PriceUpdated represents a PriceUpdated event raised by the Oracle0 contract.
type Oracle0PriceUpdated struct {
	Price     *big.Int
	Round     *big.Int
	Symbol    common.Hash
	Status    bool
	Timestamp *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterPriceUpdated is a free log retrieval operation binding the contract event 0x5f2aa51aa7889ad71d9318fa7fd83c8ff3277434249bd06073f15986e197911c.
//
// Solidity: event PriceUpdated(uint256 price, uint256 round, string indexed symbol, bool status, uint256 timestamp)
func (_Oracle0 *Oracle0Filterer) FilterPriceUpdated(opts *bind.FilterOpts, symbol []string) (*Oracle0PriceUpdatedIterator, error) {

	var symbolRule []interface{}
	for _, symbolItem := range symbol {
		symbolRule = append(symbolRule, symbolItem)
	}

	logs, sub, err := _Oracle0.contract.FilterLogs(opts, "PriceUpdated", symbolRule)
	if err != nil {
		return nil, err
	}
	return &Oracle0PriceUpdatedIterator{contract: _Oracle0.contract, event: "PriceUpdated", logs: logs, sub: sub}, nil
}

// WatchPriceUpdated is a free log subscription operation binding the contract event 0x5f2aa51aa7889ad71d9318fa7fd83c8ff3277434249bd06073f15986e197911c.
//
// Solidity: event PriceUpdated(uint256 price, uint256 round, string indexed symbol, bool status, uint256 timestamp)
func (_Oracle0 *Oracle0Filterer) WatchPriceUpdated(opts *bind.WatchOpts, sink chan<- *Oracle0PriceUpdated, symbol []string) (event.Subscription, error) {

	var symbolRule []interface{}
	for _, symbolItem := range symbol {
		symbolRule = append(symbolRule, symbolItem)
	}

	logs, sub, err := _Oracle0.contract.WatchLogs(opts, "PriceUpdated", symbolRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Oracle0PriceUpdated)
				if err := _Oracle0.contract.UnpackLog(event, "PriceUpdated", log); err != nil {
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

// ParsePriceUpdated is a log parse operation binding the contract event 0x5f2aa51aa7889ad71d9318fa7fd83c8ff3277434249bd06073f15986e197911c.
//
// Solidity: event PriceUpdated(uint256 price, uint256 round, string indexed symbol, bool status, uint256 timestamp)
func (_Oracle0 *Oracle0Filterer) ParsePriceUpdated(log types.Log) (*Oracle0PriceUpdated, error) {
	event := new(Oracle0PriceUpdated)
	if err := _Oracle0.contract.UnpackLog(event, "PriceUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Oracle0SuccessfulVoteIterator is returned from FilterSuccessfulVote and is used to iterate over the raw logs and unpacked data for SuccessfulVote events raised by the Oracle0 contract.
type Oracle0SuccessfulVoteIterator struct {
	Event *Oracle0SuccessfulVote // Event containing the contract specifics and raw log

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
func (it *Oracle0SuccessfulVoteIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Oracle0SuccessfulVote)
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
		it.Event = new(Oracle0SuccessfulVote)
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
func (it *Oracle0SuccessfulVoteIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Oracle0SuccessfulVoteIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Oracle0SuccessfulVote represents a SuccessfulVote event raised by the Oracle0 contract.
type Oracle0SuccessfulVote struct {
	Reporter common.Address
	Extra    uint8
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterSuccessfulVote is a free log retrieval operation binding the contract event 0x8bdddd7f2f2c74679ffa6beb8f86aa18bfa5baf1bfaf534d0b66596babc53f08.
//
// Solidity: event SuccessfulVote(address indexed reporter, uint8 extra)
func (_Oracle0 *Oracle0Filterer) FilterSuccessfulVote(opts *bind.FilterOpts, reporter []common.Address) (*Oracle0SuccessfulVoteIterator, error) {

	var reporterRule []interface{}
	for _, reporterItem := range reporter {
		reporterRule = append(reporterRule, reporterItem)
	}

	logs, sub, err := _Oracle0.contract.FilterLogs(opts, "SuccessfulVote", reporterRule)
	if err != nil {
		return nil, err
	}
	return &Oracle0SuccessfulVoteIterator{contract: _Oracle0.contract, event: "SuccessfulVote", logs: logs, sub: sub}, nil
}

// WatchSuccessfulVote is a free log subscription operation binding the contract event 0x8bdddd7f2f2c74679ffa6beb8f86aa18bfa5baf1bfaf534d0b66596babc53f08.
//
// Solidity: event SuccessfulVote(address indexed reporter, uint8 extra)
func (_Oracle0 *Oracle0Filterer) WatchSuccessfulVote(opts *bind.WatchOpts, sink chan<- *Oracle0SuccessfulVote, reporter []common.Address) (event.Subscription, error) {

	var reporterRule []interface{}
	for _, reporterItem := range reporter {
		reporterRule = append(reporterRule, reporterItem)
	}

	logs, sub, err := _Oracle0.contract.WatchLogs(opts, "SuccessfulVote", reporterRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Oracle0SuccessfulVote)
				if err := _Oracle0.contract.UnpackLog(event, "SuccessfulVote", log); err != nil {
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

// ParseSuccessfulVote is a log parse operation binding the contract event 0x8bdddd7f2f2c74679ffa6beb8f86aa18bfa5baf1bfaf534d0b66596babc53f08.
//
// Solidity: event SuccessfulVote(address indexed reporter, uint8 extra)
func (_Oracle0 *Oracle0Filterer) ParseSuccessfulVote(log types.Log) (*Oracle0SuccessfulVote, error) {
	event := new(Oracle0SuccessfulVote)
	if err := _Oracle0.contract.UnpackLog(event, "SuccessfulVote", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Oracle0TotalOracleRewardsIterator is returned from FilterTotalOracleRewards and is used to iterate over the raw logs and unpacked data for TotalOracleRewards events raised by the Oracle0 contract.
type Oracle0TotalOracleRewardsIterator struct {
	Event *Oracle0TotalOracleRewards // Event containing the contract specifics and raw log

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
func (it *Oracle0TotalOracleRewardsIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Oracle0TotalOracleRewards)
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
		it.Event = new(Oracle0TotalOracleRewards)
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
func (it *Oracle0TotalOracleRewardsIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Oracle0TotalOracleRewardsIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Oracle0TotalOracleRewards represents a TotalOracleRewards event raised by the Oracle0 contract.
type Oracle0TotalOracleRewards struct {
	NtnReward *big.Int
	AtnReward *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterTotalOracleRewards is a free log retrieval operation binding the contract event 0x3e5aaff9e8fd4293ae18127809c2d4069d87fe10c7de92aa39557a1edbd48fec.
//
// Solidity: event TotalOracleRewards(uint256 ntnReward, uint256 atnReward)
func (_Oracle0 *Oracle0Filterer) FilterTotalOracleRewards(opts *bind.FilterOpts) (*Oracle0TotalOracleRewardsIterator, error) {

	logs, sub, err := _Oracle0.contract.FilterLogs(opts, "TotalOracleRewards")
	if err != nil {
		return nil, err
	}
	return &Oracle0TotalOracleRewardsIterator{contract: _Oracle0.contract, event: "TotalOracleRewards", logs: logs, sub: sub}, nil
}

// WatchTotalOracleRewards is a free log subscription operation binding the contract event 0x3e5aaff9e8fd4293ae18127809c2d4069d87fe10c7de92aa39557a1edbd48fec.
//
// Solidity: event TotalOracleRewards(uint256 ntnReward, uint256 atnReward)
func (_Oracle0 *Oracle0Filterer) WatchTotalOracleRewards(opts *bind.WatchOpts, sink chan<- *Oracle0TotalOracleRewards) (event.Subscription, error) {

	logs, sub, err := _Oracle0.contract.WatchLogs(opts, "TotalOracleRewards")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Oracle0TotalOracleRewards)
				if err := _Oracle0.contract.UnpackLog(event, "TotalOracleRewards", log); err != nil {
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

// ParseTotalOracleRewards is a log parse operation binding the contract event 0x3e5aaff9e8fd4293ae18127809c2d4069d87fe10c7de92aa39557a1edbd48fec.
//
// Solidity: event TotalOracleRewards(uint256 ntnReward, uint256 atnReward)
func (_Oracle0 *Oracle0Filterer) ParseTotalOracleRewards(log types.Log) (*Oracle0TotalOracleRewards, error) {
	event := new(Oracle0TotalOracleRewards)
	if err := _Oracle0.contract.UnpackLog(event, "TotalOracleRewards", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ReentrancyGuardMetaData contains all meta data concerning the ReentrancyGuard contract.
var ReentrancyGuardMetaData = &bind.MetaData{
	ABI: "[]",
}

// ReentrancyGuardABI is the input ABI used to generate the binding from.
// Deprecated: Use ReentrancyGuardMetaData.ABI instead.
var ReentrancyGuardABI = ReentrancyGuardMetaData.ABI

// ReentrancyGuard is an auto generated Go binding around an Ethereum contract.
type ReentrancyGuard struct {
	ReentrancyGuardCaller     // Read-only binding to the contract
	ReentrancyGuardTransactor // Write-only binding to the contract
	ReentrancyGuardFilterer   // Log filterer for contract events
}

// ReentrancyGuardCaller is an auto generated read-only Go binding around an Ethereum contract.
type ReentrancyGuardCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ReentrancyGuardTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ReentrancyGuardTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ReentrancyGuardFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ReentrancyGuardFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ReentrancyGuardSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ReentrancyGuardSession struct {
	Contract     *ReentrancyGuard  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ReentrancyGuardCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ReentrancyGuardCallerSession struct {
	Contract *ReentrancyGuardCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// ReentrancyGuardTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ReentrancyGuardTransactorSession struct {
	Contract     *ReentrancyGuardTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// ReentrancyGuardRaw is an auto generated low-level Go binding around an Ethereum contract.
type ReentrancyGuardRaw struct {
	Contract *ReentrancyGuard // Generic contract binding to access the raw methods on
}

// ReentrancyGuardCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ReentrancyGuardCallerRaw struct {
	Contract *ReentrancyGuardCaller // Generic read-only contract binding to access the raw methods on
}

// ReentrancyGuardTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ReentrancyGuardTransactorRaw struct {
	Contract *ReentrancyGuardTransactor // Generic write-only contract binding to access the raw methods on
}

// NewReentrancyGuard creates a new instance of ReentrancyGuard, bound to a specific deployed contract.
func NewReentrancyGuard(address common.Address, backend bind.ContractBackend) (*ReentrancyGuard, error) {
	contract, err := bindReentrancyGuard(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ReentrancyGuard{ReentrancyGuardCaller: ReentrancyGuardCaller{contract: contract}, ReentrancyGuardTransactor: ReentrancyGuardTransactor{contract: contract}, ReentrancyGuardFilterer: ReentrancyGuardFilterer{contract: contract}}, nil
}

// NewReentrancyGuardCaller creates a new read-only instance of ReentrancyGuard, bound to a specific deployed contract.
func NewReentrancyGuardCaller(address common.Address, caller bind.ContractCaller) (*ReentrancyGuardCaller, error) {
	contract, err := bindReentrancyGuard(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ReentrancyGuardCaller{contract: contract}, nil
}

// NewReentrancyGuardTransactor creates a new write-only instance of ReentrancyGuard, bound to a specific deployed contract.
func NewReentrancyGuardTransactor(address common.Address, transactor bind.ContractTransactor) (*ReentrancyGuardTransactor, error) {
	contract, err := bindReentrancyGuard(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ReentrancyGuardTransactor{contract: contract}, nil
}

// NewReentrancyGuardFilterer creates a new log filterer instance of ReentrancyGuard, bound to a specific deployed contract.
func NewReentrancyGuardFilterer(address common.Address, filterer bind.ContractFilterer) (*ReentrancyGuardFilterer, error) {
	contract, err := bindReentrancyGuard(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ReentrancyGuardFilterer{contract: contract}, nil
}

// bindReentrancyGuard binds a generic wrapper to an already deployed contract.
func bindReentrancyGuard(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ReentrancyGuardABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ReentrancyGuard *ReentrancyGuardRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ReentrancyGuard.Contract.ReentrancyGuardCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ReentrancyGuard *ReentrancyGuardRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ReentrancyGuard.Contract.ReentrancyGuardTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ReentrancyGuard *ReentrancyGuardRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ReentrancyGuard.Contract.ReentrancyGuardTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ReentrancyGuard *ReentrancyGuardCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ReentrancyGuard.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ReentrancyGuard *ReentrancyGuardTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ReentrancyGuard.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ReentrancyGuard *ReentrancyGuardTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ReentrancyGuard.Contract.contract.Transact(opts, method, params...)
}
