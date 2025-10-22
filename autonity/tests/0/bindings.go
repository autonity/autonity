// Code generated for internal testing purposes only - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package tests0

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	ethereum "github.com/autonity/autonity"
	"github.com/autonity/autonity/accounts/abi"
	"github.com/autonity/autonity/accounts/abi/bind"
	"github.com/autonity/autonity/autonity/tests"
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

// type alias Runner, otherwise functions cannot be defined on it
type Runner tests.Runner

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
func (r *Runner) DeployEnumerableSet(opts *tests.RunOptions) (common.Address, uint64, *EnumerableSet, error) {
	parsed, err := EnumerableSetMetaData.GetAbi()
	if err != nil {
		return common.Address{}, 0, nil, err
	}
	if parsed == nil {
		return common.Address{}, 0, nil, errors.New("GetABI returned nil")
	}

	address, gasConsumed, c, data, err := (*tests.Runner)(r).DeployContract(opts, parsed, common.FromHex(EnumerableSetBin))
	if err != nil {
		return common.Address{}, 0, nil, (&EnumerableSet{Contract: c}).DecodeError(data, err)
	}
	return address, gasConsumed, &EnumerableSet{Contract: c}, nil
}

// EnumerableSet is an auto generated Go binding around an Ethereum contract.
type EnumerableSet struct {
	*tests.Contract
}

func (_EnumerableSet *EnumerableSet) DecodeError(data []byte, err error) error {
	if err == nil {
		return nil
	}

	reason, _ := abi.UnpackRevert(data)
	return fmt.Errorf("%w: %s", err, reason)

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
	*tests.Contract
}

// GetScaleFactor is a free data retrieval call binding the contract method 0x7f5e2f11.
//
// Solidity: function getScaleFactor() view returns(uint256)
func (_IACU *IACU) GetScaleFactor(opts *tests.RunOptions) (*big.Int, uint64, error) {
	data, consumed, err := _IACU.Call(opts, "getScaleFactor")

	if err != nil {
		return *new(*big.Int), consumed, _IACU.DecodeError(data, err)
	}
	out, err := _IACU.Contract.Abi().Unpack("getScaleFactor", data)
	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// Value is a free data retrieval call binding the contract method 0x3fa4f245.
//
// Solidity: function value() view returns(uint256)
func (_IACU *IACU) Value(opts *tests.RunOptions) (*big.Int, uint64, error) {
	data, consumed, err := _IACU.Call(opts, "value")

	if err != nil {
		return *new(*big.Int), consumed, _IACU.DecodeError(data, err)
	}
	out, err := _IACU.Contract.Abi().Unpack("value", data)
	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// SetOperator is a free data retrieval call for a paid mutator transaction binding the contract method 0xb3ab15fb.
// Similar to eth_call from rpc calls or function.call from truffle, it reverts the state after the call and returns the output. The output is extracted
// the same way as done above for view only functions.
// Solidity: function setOperator(address operator) returns()
func (_IACU *IACU) CallSetOperator(r *tests.Runner, opts *tests.RunOptions, operator common.Address) (uint64, error) {
	snap := r.Snapshot()

	data, consumed, err := _IACU.Call(opts, "setOperator", operator)
	r.RevertSnapshot(snap)
	return consumed, _IACU.DecodeError(data, err)

}

// SetOracle is a free data retrieval call for a paid mutator transaction binding the contract method 0x7adbf973.
// Similar to eth_call from rpc calls or function.call from truffle, it reverts the state after the call and returns the output. The output is extracted
// the same way as done above for view only functions.
// Solidity: function setOracle(address oracle) returns()
func (_IACU *IACU) CallSetOracle(r *tests.Runner, opts *tests.RunOptions, oracle common.Address) (uint64, error) {
	snap := r.Snapshot()

	data, consumed, err := _IACU.Call(opts, "setOracle", oracle)
	r.RevertSnapshot(snap)
	return consumed, _IACU.DecodeError(data, err)

}

// Update is a free data retrieval call for a paid mutator transaction binding the contract method 0xa2e62045.
// Similar to eth_call from rpc calls or function.call from truffle, it reverts the state after the call and returns the output. The output is extracted
// the same way as done above for view only functions.
// Solidity: function update() returns(bool status)
func (_IACU *IACU) CallUpdate(r *tests.Runner, opts *tests.RunOptions) (bool, uint64, error) {
	snap := r.Snapshot()

	data, consumed, err := _IACU.Call(opts, "update")
	r.RevertSnapshot(snap)

	if err != nil {
		return *new(bool), consumed, _IACU.DecodeError(data, err)
	}
	out, err := _IACU.Contract.Abi().Unpack("update", data)
	if err != nil {
		return *new(bool), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, consumed, err

}

// SetOperator is a paid mutator transaction binding the contract method 0xb3ab15fb.
//
// Solidity: function setOperator(address operator) returns()
func (_IACU *IACU) SetOperator(opts *tests.RunOptions, operator common.Address) (uint64, error) {
	data, consumed, err := _IACU.Call(opts, "setOperator", operator)
	return consumed, _IACU.DecodeError(data, err)
}

// SetOracle is a paid mutator transaction binding the contract method 0x7adbf973.
//
// Solidity: function setOracle(address oracle) returns()
func (_IACU *IACU) SetOracle(opts *tests.RunOptions, oracle common.Address) (uint64, error) {
	data, consumed, err := _IACU.Call(opts, "setOracle", oracle)
	return consumed, _IACU.DecodeError(data, err)
}

// Update is a paid mutator transaction binding the contract method 0xa2e62045.
//
// Solidity: function update() returns(bool status)
func (_IACU *IACU) Update(opts *tests.RunOptions) (uint64, error) {
	data, consumed, err := _IACU.Call(opts, "update")
	return consumed, _IACU.DecodeError(data, err)
}

func (_IACU *IACU) DecodeError(data []byte, err error) error {
	if err == nil {
		return nil
	}

	reason, _ := abi.UnpackRevert(data)
	return fmt.Errorf("%w: %s", err, reason)

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
	*tests.Contract
}

// GetConfig is a free data retrieval call binding the contract method 0xc3f909d4.
//
// Solidity: function getConfig() view returns((uint256,uint256,uint256,(uint256,uint256,uint256),(uint256,uint256,uint256)))
func (_IAccountability *IAccountability) GetConfig(opts *tests.RunOptions) (IAccountabilityConfig, uint64, error) {
	data, consumed, err := _IAccountability.Call(opts, "getConfig")

	if err != nil {
		return *new(IAccountabilityConfig), consumed, _IAccountability.DecodeError(data, err)
	}
	out, err := _IAccountability.Contract.Abi().Unpack("getConfig", data)
	if err != nil {
		return *new(IAccountabilityConfig), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(IAccountabilityConfig)).(*IAccountabilityConfig)
	return out0, consumed, err

}

// GetGracePeriod is a free data retrieval call binding the contract method 0xdbd18388.
//
// Solidity: function getGracePeriod() view returns(uint256)
func (_IAccountability *IAccountability) GetGracePeriod(opts *tests.RunOptions) (*big.Int, uint64, error) {
	data, consumed, err := _IAccountability.Call(opts, "getGracePeriod")

	if err != nil {
		return *new(*big.Int), consumed, _IAccountability.DecodeError(data, err)
	}
	out, err := _IAccountability.Contract.Abi().Unpack("getGracePeriod", data)
	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// DistributeRewards is a free data retrieval call for a paid mutator transaction binding the contract method 0xa8031a1d.
// Similar to eth_call from rpc calls or function.call from truffle, it reverts the state after the call and returns the output. The output is extracted
// the same way as done above for view only functions.
// Solidity: function distributeRewards(address _validator, uint256 _ntnReward) payable returns()
func (_IAccountability *IAccountability) CallDistributeRewards(r *tests.Runner, opts *tests.RunOptions, _validator common.Address, _ntnReward *big.Int) (uint64, error) {
	snap := r.Snapshot()

	data, consumed, err := _IAccountability.Call(opts, "distributeRewards", _validator, _ntnReward)
	r.RevertSnapshot(snap)
	return consumed, _IAccountability.DecodeError(data, err)

}

// Finalize is a free data retrieval call for a paid mutator transaction binding the contract method 0x6c9789b0.
// Similar to eth_call from rpc calls or function.call from truffle, it reverts the state after the call and returns the output. The output is extracted
// the same way as done above for view only functions.
// Solidity: function finalize(bool _epochEnd) returns(uint256, uint256, uint256)
func (_IAccountability *IAccountability) CallFinalize(r *tests.Runner, opts *tests.RunOptions, _epochEnd bool) (*big.Int, *big.Int, *big.Int, uint64, error) {
	snap := r.Snapshot()

	data, consumed, err := _IAccountability.Call(opts, "finalize", _epochEnd)
	r.RevertSnapshot(snap)

	if err != nil {
		return *new(*big.Int), *new(*big.Int), *new(*big.Int), consumed, _IAccountability.DecodeError(data, err)
	}
	out, err := _IAccountability.Contract.Abi().Unpack("finalize", data)
	if err != nil {
		return *new(*big.Int), *new(*big.Int), *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	out1 := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	out2 := *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	return out0, out1, out2, consumed, err

}

// SetCommittee is a free data retrieval call for a paid mutator transaction binding the contract method 0xe08b14ed.
// Similar to eth_call from rpc calls or function.call from truffle, it reverts the state after the call and returns the output. The output is extracted
// the same way as done above for view only functions.
// Solidity: function setCommittee(address[] _committee) returns()
func (_IAccountability *IAccountability) CallSetCommittee(r *tests.Runner, opts *tests.RunOptions, _committee []common.Address) (uint64, error) {
	snap := r.Snapshot()

	data, consumed, err := _IAccountability.Call(opts, "setCommittee", _committee)
	r.RevertSnapshot(snap)
	return consumed, _IAccountability.DecodeError(data, err)

}

// DistributeRewards is a paid mutator transaction binding the contract method 0xa8031a1d.
//
// Solidity: function distributeRewards(address _validator, uint256 _ntnReward) payable returns()
func (_IAccountability *IAccountability) DistributeRewards(opts *tests.RunOptions, _validator common.Address, _ntnReward *big.Int) (uint64, error) {
	data, consumed, err := _IAccountability.Call(opts, "distributeRewards", _validator, _ntnReward)
	return consumed, _IAccountability.DecodeError(data, err)
}

// Finalize is a paid mutator transaction binding the contract method 0x6c9789b0.
//
// Solidity: function finalize(bool _epochEnd) returns(uint256, uint256, uint256)
func (_IAccountability *IAccountability) Finalize(opts *tests.RunOptions, _epochEnd bool) (uint64, error) {
	data, consumed, err := _IAccountability.Call(opts, "finalize", _epochEnd)
	return consumed, _IAccountability.DecodeError(data, err)
}

// SetCommittee is a paid mutator transaction binding the contract method 0xe08b14ed.
//
// Solidity: function setCommittee(address[] _committee) returns()
func (_IAccountability *IAccountability) SetCommittee(opts *tests.RunOptions, _committee []common.Address) (uint64, error) {
	data, consumed, err := _IAccountability.Call(opts, "setCommittee", _committee)
	return consumed, _IAccountability.DecodeError(data, err)
}

func (_IAccountability *IAccountability) DecodeError(data []byte, err error) error {
	if err == nil {
		return nil
	}

	reason, _ := abi.UnpackRevert(data)
	return fmt.Errorf("%w: %s", err, reason)

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
	*tests.Contract
}

// PaidInterest is a free data retrieval call for a paid mutator transaction binding the contract method 0x96e4547c.
// Similar to eth_call from rpc calls or function.call from truffle, it reverts the state after the call and returns the output. The output is extracted
// the same way as done above for view only functions.
// Solidity: function paidInterest() payable returns()
func (_IAuctioneer *IAuctioneer) CallPaidInterest(r *tests.Runner, opts *tests.RunOptions) (uint64, error) {
	snap := r.Snapshot()

	data, consumed, err := _IAuctioneer.Call(opts, "paidInterest")
	r.RevertSnapshot(snap)
	return consumed, _IAuctioneer.DecodeError(data, err)

}

// SetOperator is a free data retrieval call for a paid mutator transaction binding the contract method 0xb3ab15fb.
// Similar to eth_call from rpc calls or function.call from truffle, it reverts the state after the call and returns the output. The output is extracted
// the same way as done above for view only functions.
// Solidity: function setOperator(address operator) returns()
func (_IAuctioneer *IAuctioneer) CallSetOperator(r *tests.Runner, opts *tests.RunOptions, operator common.Address) (uint64, error) {
	snap := r.Snapshot()

	data, consumed, err := _IAuctioneer.Call(opts, "setOperator", operator)
	r.RevertSnapshot(snap)
	return consumed, _IAuctioneer.DecodeError(data, err)

}

// SetOracle is a free data retrieval call for a paid mutator transaction binding the contract method 0x7adbf973.
// Similar to eth_call from rpc calls or function.call from truffle, it reverts the state after the call and returns the output. The output is extracted
// the same way as done above for view only functions.
// Solidity: function setOracle(address oracle) returns()
func (_IAuctioneer *IAuctioneer) CallSetOracle(r *tests.Runner, opts *tests.RunOptions, oracle common.Address) (uint64, error) {
	snap := r.Snapshot()

	data, consumed, err := _IAuctioneer.Call(opts, "setOracle", oracle)
	r.RevertSnapshot(snap)
	return consumed, _IAuctioneer.DecodeError(data, err)

}

// SetStabilization is a free data retrieval call for a paid mutator transaction binding the contract method 0x4f505895.
// Similar to eth_call from rpc calls or function.call from truffle, it reverts the state after the call and returns the output. The output is extracted
// the same way as done above for view only functions.
// Solidity: function setStabilization(address stabilization) returns()
func (_IAuctioneer *IAuctioneer) CallSetStabilization(r *tests.Runner, opts *tests.RunOptions, stabilization common.Address) (uint64, error) {
	snap := r.Snapshot()

	data, consumed, err := _IAuctioneer.Call(opts, "setStabilization", stabilization)
	r.RevertSnapshot(snap)
	return consumed, _IAuctioneer.DecodeError(data, err)

}

// PaidInterest is a paid mutator transaction binding the contract method 0x96e4547c.
//
// Solidity: function paidInterest() payable returns()
func (_IAuctioneer *IAuctioneer) PaidInterest(opts *tests.RunOptions) (uint64, error) {
	data, consumed, err := _IAuctioneer.Call(opts, "paidInterest")
	return consumed, _IAuctioneer.DecodeError(data, err)
}

// SetOperator is a paid mutator transaction binding the contract method 0xb3ab15fb.
//
// Solidity: function setOperator(address operator) returns()
func (_IAuctioneer *IAuctioneer) SetOperator(opts *tests.RunOptions, operator common.Address) (uint64, error) {
	data, consumed, err := _IAuctioneer.Call(opts, "setOperator", operator)
	return consumed, _IAuctioneer.DecodeError(data, err)
}

// SetOracle is a paid mutator transaction binding the contract method 0x7adbf973.
//
// Solidity: function setOracle(address oracle) returns()
func (_IAuctioneer *IAuctioneer) SetOracle(opts *tests.RunOptions, oracle common.Address) (uint64, error) {
	data, consumed, err := _IAuctioneer.Call(opts, "setOracle", oracle)
	return consumed, _IAuctioneer.DecodeError(data, err)
}

// SetStabilization is a paid mutator transaction binding the contract method 0x4f505895.
//
// Solidity: function setStabilization(address stabilization) returns()
func (_IAuctioneer *IAuctioneer) SetStabilization(opts *tests.RunOptions, stabilization common.Address) (uint64, error) {
	data, consumed, err := _IAuctioneer.Call(opts, "setStabilization", stabilization)
	return consumed, _IAuctioneer.DecodeError(data, err)
}

func (_IAuctioneer *IAuctioneer) DecodeError(data []byte, err error) error {
	if err == nil {
		return nil
	}

	reason, _ := abi.UnpackRevert(data)
	return fmt.Errorf("%w: %s", err, reason)

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
	*tests.Contract
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_IAutonity *IAutonity) Allowance(opts *tests.RunOptions, owner common.Address, spender common.Address) (*big.Int, uint64, error) {
	data, consumed, err := _IAutonity.Call(opts, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), consumed, _IAutonity.DecodeError(data, err)
	}
	out, err := _IAutonity.Contract.Abi().Unpack("allowance", data)
	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_IAutonity *IAutonity) BalanceOf(opts *tests.RunOptions, account common.Address) (*big.Int, uint64, error) {
	data, consumed, err := _IAutonity.Call(opts, "balanceOf", account)

	if err != nil {
		return *new(*big.Int), consumed, _IAutonity.DecodeError(data, err)
	}
	out, err := _IAutonity.Contract.Abi().Unpack("balanceOf", data)
	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// BondingAllowance is a free data retrieval call binding the contract method 0xe0e01d54.
//
// Solidity: function bondingAllowance(address _owner, address _caller) view returns(uint256)
func (_IAutonity *IAutonity) BondingAllowance(opts *tests.RunOptions, _owner common.Address, _caller common.Address) (*big.Int, uint64, error) {
	data, consumed, err := _IAutonity.Call(opts, "bondingAllowance", _owner, _caller)

	if err != nil {
		return *new(*big.Int), consumed, _IAutonity.DecodeError(data, err)
	}
	out, err := _IAutonity.Contract.Abi().Unpack("bondingAllowance", data)
	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// CirculatingSupply is a free data retrieval call binding the contract method 0x9358928b.
//
// Solidity: function circulatingSupply() view returns(uint256)
func (_IAutonity *IAutonity) CirculatingSupply(opts *tests.RunOptions) (*big.Int, uint64, error) {
	data, consumed, err := _IAutonity.Call(opts, "circulatingSupply")

	if err != nil {
		return *new(*big.Int), consumed, _IAutonity.DecodeError(data, err)
	}
	out, err := _IAutonity.Contract.Abi().Unpack("circulatingSupply", data)
	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// GetBlockPeriod is a free data retrieval call binding the contract method 0x43645969.
//
// Solidity: function getBlockPeriod() view returns(uint256)
func (_IAutonity *IAutonity) GetBlockPeriod(opts *tests.RunOptions) (*big.Int, uint64, error) {
	data, consumed, err := _IAutonity.Call(opts, "getBlockPeriod")

	if err != nil {
		return *new(*big.Int), consumed, _IAutonity.DecodeError(data, err)
	}
	out, err := _IAutonity.Contract.Abi().Unpack("getBlockPeriod", data)
	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// GetBondingRequestByID is a free data retrieval call binding the contract method 0x8ebb48b7.
//
// Solidity: function getBondingRequestByID(uint256 _id) view returns((address,address,uint256,uint256))
func (_IAutonity *IAutonity) GetBondingRequestByID(opts *tests.RunOptions, _id *big.Int) (IAutonityBondingRequest, uint64, error) {
	data, consumed, err := _IAutonity.Call(opts, "getBondingRequestByID", _id)

	if err != nil {
		return *new(IAutonityBondingRequest), consumed, _IAutonity.DecodeError(data, err)
	}
	out, err := _IAutonity.Contract.Abi().Unpack("getBondingRequestByID", data)
	if err != nil {
		return *new(IAutonityBondingRequest), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(IAutonityBondingRequest)).(*IAutonityBondingRequest)
	return out0, consumed, err

}

// GetClientConfig is a free data retrieval call binding the contract method 0xf3f759c1.
//
// Solidity: function getClientConfig() view returns((uint256,uint256,uint256,uint256,(uint256,uint256,uint256),(uint256,uint256,uint256,uint256)))
func (_IAutonity *IAutonity) GetClientConfig(opts *tests.RunOptions) (IAutonityClientAwareConfig, uint64, error) {
	data, consumed, err := _IAutonity.Call(opts, "getClientConfig")

	if err != nil {
		return *new(IAutonityClientAwareConfig), consumed, _IAutonity.DecodeError(data, err)
	}
	out, err := _IAutonity.Contract.Abi().Unpack("getClientConfig", data)
	if err != nil {
		return *new(IAutonityClientAwareConfig), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(IAutonityClientAwareConfig)).(*IAutonityClientAwareConfig)
	return out0, consumed, err

}

// GetCommittee is a free data retrieval call binding the contract method 0xab8f6ffe.
//
// Solidity: function getCommittee() view returns((address,uint256,bytes)[])
func (_IAutonity *IAutonity) GetCommittee(opts *tests.RunOptions) ([]IAutonityCommitteeMember, uint64, error) {
	data, consumed, err := _IAutonity.Call(opts, "getCommittee")

	if err != nil {
		return *new([]IAutonityCommitteeMember), consumed, _IAutonity.DecodeError(data, err)
	}
	out, err := _IAutonity.Contract.Abi().Unpack("getCommittee", data)
	if err != nil {
		return *new([]IAutonityCommitteeMember), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new([]IAutonityCommitteeMember)).(*[]IAutonityCommitteeMember)
	return out0, consumed, err

}

// GetCommitteeEnodes is a free data retrieval call binding the contract method 0xa8b2216e.
//
// Solidity: function getCommitteeEnodes() view returns(string[])
func (_IAutonity *IAutonity) GetCommitteeEnodes(opts *tests.RunOptions) ([]string, uint64, error) {
	data, consumed, err := _IAutonity.Call(opts, "getCommitteeEnodes")

	if err != nil {
		return *new([]string), consumed, _IAutonity.DecodeError(data, err)
	}
	out, err := _IAutonity.Contract.Abi().Unpack("getCommitteeEnodes", data)
	if err != nil {
		return *new([]string), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new([]string)).(*[]string)
	return out0, consumed, err

}

// GetConfig is a free data retrieval call binding the contract method 0xc3f909d4.
//
// Solidity: function getConfig() view returns(((uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256,address,address,uint256,uint256),(address,address,address,address,address,address,address,address,address),(address,uint256,uint256,uint256,uint256,uint256,uint256,uint256),uint256))
func (_IAutonity *IAutonity) GetConfig(opts *tests.RunOptions) (IAutonityConfig, uint64, error) {
	data, consumed, err := _IAutonity.Call(opts, "getConfig")

	if err != nil {
		return *new(IAutonityConfig), consumed, _IAutonity.DecodeError(data, err)
	}
	out, err := _IAutonity.Contract.Abi().Unpack("getConfig", data)
	if err != nil {
		return *new(IAutonityConfig), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(IAutonityConfig)).(*IAutonityConfig)
	return out0, consumed, err

}

// GetCurrentCommitteeSize is a free data retrieval call binding the contract method 0x2b56feac.
//
// Solidity: function getCurrentCommitteeSize() view returns(uint256)
func (_IAutonity *IAutonity) GetCurrentCommitteeSize(opts *tests.RunOptions) (*big.Int, uint64, error) {
	data, consumed, err := _IAutonity.Call(opts, "getCurrentCommitteeSize")

	if err != nil {
		return *new(*big.Int), consumed, _IAutonity.DecodeError(data, err)
	}
	out, err := _IAutonity.Contract.Abi().Unpack("getCurrentCommitteeSize", data)
	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// GetCurrentEpochPeriod is a free data retrieval call binding the contract method 0x0aac2da1.
//
// Solidity: function getCurrentEpochPeriod() view returns(uint256)
func (_IAutonity *IAutonity) GetCurrentEpochPeriod(opts *tests.RunOptions) (*big.Int, uint64, error) {
	data, consumed, err := _IAutonity.Call(opts, "getCurrentEpochPeriod")

	if err != nil {
		return *new(*big.Int), consumed, _IAutonity.DecodeError(data, err)
	}
	out, err := _IAutonity.Contract.Abi().Unpack("getCurrentEpochPeriod", data)
	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// GetEpochByHeight is a free data retrieval call binding the contract method 0xaffb1cf1.
//
// Solidity: function getEpochByHeight(uint256 _height) view returns(((address,uint256,bytes)[],uint256,uint256,uint256,uint256,(uint256,uint256,uint256,uint256)))
func (_IAutonity *IAutonity) GetEpochByHeight(opts *tests.RunOptions, _height *big.Int) (IAutonityEpochInfo, uint64, error) {
	data, consumed, err := _IAutonity.Call(opts, "getEpochByHeight", _height)

	if err != nil {
		return *new(IAutonityEpochInfo), consumed, _IAutonity.DecodeError(data, err)
	}
	out, err := _IAutonity.Contract.Abi().Unpack("getEpochByHeight", data)
	if err != nil {
		return *new(IAutonityEpochInfo), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(IAutonityEpochInfo)).(*IAutonityEpochInfo)
	return out0, consumed, err

}

// GetEpochFromBlock is a free data retrieval call binding the contract method 0x96b477cb.
//
// Solidity: function getEpochFromBlock(uint256 _block) view returns(uint256)
func (_IAutonity *IAutonity) GetEpochFromBlock(opts *tests.RunOptions, _block *big.Int) (*big.Int, uint64, error) {
	data, consumed, err := _IAutonity.Call(opts, "getEpochFromBlock", _block)

	if err != nil {
		return *new(*big.Int), consumed, _IAutonity.DecodeError(data, err)
	}
	out, err := _IAutonity.Contract.Abi().Unpack("getEpochFromBlock", data)
	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// GetEpochID is a free data retrieval call binding the contract method 0x6fc53515.
//
// Solidity: function getEpochID() view returns(uint256)
func (_IAutonity *IAutonity) GetEpochID(opts *tests.RunOptions) (*big.Int, uint64, error) {
	data, consumed, err := _IAutonity.Call(opts, "getEpochID")

	if err != nil {
		return *new(*big.Int), consumed, _IAutonity.DecodeError(data, err)
	}
	out, err := _IAutonity.Contract.Abi().Unpack("getEpochID", data)
	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// GetEpochInfo is a free data retrieval call binding the contract method 0xa9fd1a8f.
//
// Solidity: function getEpochInfo() view returns(((address,uint256,bytes)[],uint256,uint256,uint256,uint256,(uint256,uint256,uint256,uint256)))
func (_IAutonity *IAutonity) GetEpochInfo(opts *tests.RunOptions) (IAutonityEpochInfo, uint64, error) {
	data, consumed, err := _IAutonity.Call(opts, "getEpochInfo")

	if err != nil {
		return *new(IAutonityEpochInfo), consumed, _IAutonity.DecodeError(data, err)
	}
	out, err := _IAutonity.Contract.Abi().Unpack("getEpochInfo", data)
	if err != nil {
		return *new(IAutonityEpochInfo), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(IAutonityEpochInfo)).(*IAutonityEpochInfo)
	return out0, consumed, err

}

// GetEpochPeriod is a free data retrieval call binding the contract method 0xdfb1a4d2.
//
// Solidity: function getEpochPeriod() view returns(uint256)
func (_IAutonity *IAutonity) GetEpochPeriod(opts *tests.RunOptions) (*big.Int, uint64, error) {
	data, consumed, err := _IAutonity.Call(opts, "getEpochPeriod")

	if err != nil {
		return *new(*big.Int), consumed, _IAutonity.DecodeError(data, err)
	}
	out, err := _IAutonity.Contract.Abi().Unpack("getEpochPeriod", data)
	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// GetEpochTotalBondedStake is a free data retrieval call binding the contract method 0x4efcd15f.
//
// Solidity: function getEpochTotalBondedStake() view returns(uint256)
func (_IAutonity *IAutonity) GetEpochTotalBondedStake(opts *tests.RunOptions) (*big.Int, uint64, error) {
	data, consumed, err := _IAutonity.Call(opts, "getEpochTotalBondedStake")

	if err != nil {
		return *new(*big.Int), consumed, _IAutonity.DecodeError(data, err)
	}
	out, err := _IAutonity.Contract.Abi().Unpack("getEpochTotalBondedStake", data)
	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// GetInflationReserve is a free data retrieval call binding the contract method 0x4651e07b.
//
// Solidity: function getInflationReserve() view returns(uint256)
func (_IAutonity *IAutonity) GetInflationReserve(opts *tests.RunOptions) (*big.Int, uint64, error) {
	data, consumed, err := _IAutonity.Call(opts, "getInflationReserve")

	if err != nil {
		return *new(*big.Int), consumed, _IAutonity.DecodeError(data, err)
	}
	out, err := _IAutonity.Contract.Abi().Unpack("getInflationReserve", data)
	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// GetLastEpochBlock is a free data retrieval call binding the contract method 0x731b3a03.
//
// Solidity: function getLastEpochBlock() view returns(uint256)
func (_IAutonity *IAutonity) GetLastEpochBlock(opts *tests.RunOptions) (*big.Int, uint64, error) {
	data, consumed, err := _IAutonity.Call(opts, "getLastEpochBlock")

	if err != nil {
		return *new(*big.Int), consumed, _IAutonity.DecodeError(data, err)
	}
	out, err := _IAutonity.Contract.Abi().Unpack("getLastEpochBlock", data)
	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// GetLastEpochTime is a free data retrieval call binding the contract method 0xba522458.
//
// Solidity: function getLastEpochTime() view returns(uint256)
func (_IAutonity *IAutonity) GetLastEpochTime(opts *tests.RunOptions) (*big.Int, uint64, error) {
	data, consumed, err := _IAutonity.Call(opts, "getLastEpochTime")

	if err != nil {
		return *new(*big.Int), consumed, _IAutonity.DecodeError(data, err)
	}
	out, err := _IAutonity.Contract.Abi().Unpack("getLastEpochTime", data)
	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// GetLiquidLogicContract is a free data retrieval call binding the contract method 0x4c1f1c77.
//
// Solidity: function getLiquidLogicContract() view returns(address)
func (_IAutonity *IAutonity) GetLiquidLogicContract(opts *tests.RunOptions) (common.Address, uint64, error) {
	data, consumed, err := _IAutonity.Call(opts, "getLiquidLogicContract")

	if err != nil {
		return *new(common.Address), consumed, _IAutonity.DecodeError(data, err)
	}
	out, err := _IAutonity.Contract.Abi().Unpack("getLiquidLogicContract", data)
	if err != nil {
		return *new(common.Address), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, consumed, err

}

// GetMaxCommitteeSize is a free data retrieval call binding the contract method 0x819b6463.
//
// Solidity: function getMaxCommitteeSize() view returns(uint256)
func (_IAutonity *IAutonity) GetMaxCommitteeSize(opts *tests.RunOptions) (*big.Int, uint64, error) {
	data, consumed, err := _IAutonity.Call(opts, "getMaxCommitteeSize")

	if err != nil {
		return *new(*big.Int), consumed, _IAutonity.DecodeError(data, err)
	}
	out, err := _IAutonity.Contract.Abi().Unpack("getMaxCommitteeSize", data)
	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// GetMaxScheduleDuration is a free data retrieval call binding the contract method 0xfed76a56.
//
// Solidity: function getMaxScheduleDuration() view returns(uint256)
func (_IAutonity *IAutonity) GetMaxScheduleDuration(opts *tests.RunOptions) (*big.Int, uint64, error) {
	data, consumed, err := _IAutonity.Call(opts, "getMaxScheduleDuration")

	if err != nil {
		return *new(*big.Int), consumed, _IAutonity.DecodeError(data, err)
	}
	out, err := _IAutonity.Contract.Abi().Unpack("getMaxScheduleDuration", data)
	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// GetMinimumBaseFee is a free data retrieval call binding the contract method 0x11220633.
//
// Solidity: function getMinimumBaseFee() view returns(uint256)
func (_IAutonity *IAutonity) GetMinimumBaseFee(opts *tests.RunOptions) (*big.Int, uint64, error) {
	data, consumed, err := _IAutonity.Call(opts, "getMinimumBaseFee")

	if err != nil {
		return *new(*big.Int), consumed, _IAutonity.DecodeError(data, err)
	}
	out, err := _IAutonity.Contract.Abi().Unpack("getMinimumBaseFee", data)
	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// GetNextEpochBlock is a free data retrieval call binding the contract method 0x25ce1bb9.
//
// Solidity: function getNextEpochBlock() view returns(uint256)
func (_IAutonity *IAutonity) GetNextEpochBlock(opts *tests.RunOptions) (*big.Int, uint64, error) {
	data, consumed, err := _IAutonity.Call(opts, "getNextEpochBlock")

	if err != nil {
		return *new(*big.Int), consumed, _IAutonity.DecodeError(data, err)
	}
	out, err := _IAutonity.Contract.Abi().Unpack("getNextEpochBlock", data)
	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// GetOperator is a free data retrieval call binding the contract method 0xe7f43c68.
//
// Solidity: function getOperator() view returns(address)
func (_IAutonity *IAutonity) GetOperator(opts *tests.RunOptions) (common.Address, uint64, error) {
	data, consumed, err := _IAutonity.Call(opts, "getOperator")

	if err != nil {
		return *new(common.Address), consumed, _IAutonity.DecodeError(data, err)
	}
	out, err := _IAutonity.Contract.Abi().Unpack("getOperator", data)
	if err != nil {
		return *new(common.Address), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, consumed, err

}

// GetOracle is a free data retrieval call binding the contract method 0x833b1fce.
//
// Solidity: function getOracle() view returns(address)
func (_IAutonity *IAutonity) GetOracle(opts *tests.RunOptions) (common.Address, uint64, error) {
	data, consumed, err := _IAutonity.Call(opts, "getOracle")

	if err != nil {
		return *new(common.Address), consumed, _IAutonity.DecodeError(data, err)
	}
	out, err := _IAutonity.Contract.Abi().Unpack("getOracle", data)
	if err != nil {
		return *new(common.Address), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, consumed, err

}

// GetSchedule is a free data retrieval call binding the contract method 0x7264c4da.
//
// Solidity: function getSchedule(address _vault, uint256 _id) view returns((uint256,uint256,uint256,uint256,uint256))
func (_IAutonity *IAutonity) GetSchedule(opts *tests.RunOptions, _vault common.Address, _id *big.Int) (IScheduleControllerSchedule, uint64, error) {
	data, consumed, err := _IAutonity.Call(opts, "getSchedule", _vault, _id)

	if err != nil {
		return *new(IScheduleControllerSchedule), consumed, _IAutonity.DecodeError(data, err)
	}
	out, err := _IAutonity.Contract.Abi().Unpack("getSchedule", data)
	if err != nil {
		return *new(IScheduleControllerSchedule), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(IScheduleControllerSchedule)).(*IScheduleControllerSchedule)
	return out0, consumed, err

}

// GetTotalSchedules is a free data retrieval call binding the contract method 0x088566e9.
//
// Solidity: function getTotalSchedules(address _vault) view returns(uint256)
func (_IAutonity *IAutonity) GetTotalSchedules(opts *tests.RunOptions, _vault common.Address) (*big.Int, uint64, error) {
	data, consumed, err := _IAutonity.Call(opts, "getTotalSchedules", _vault)

	if err != nil {
		return *new(*big.Int), consumed, _IAutonity.DecodeError(data, err)
	}
	out, err := _IAutonity.Contract.Abi().Unpack("getTotalSchedules", data)
	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// GetTreasuryAccount is a free data retrieval call binding the contract method 0xf7866ee3.
//
// Solidity: function getTreasuryAccount() view returns(address)
func (_IAutonity *IAutonity) GetTreasuryAccount(opts *tests.RunOptions) (common.Address, uint64, error) {
	data, consumed, err := _IAutonity.Call(opts, "getTreasuryAccount")

	if err != nil {
		return *new(common.Address), consumed, _IAutonity.DecodeError(data, err)
	}
	out, err := _IAutonity.Contract.Abi().Unpack("getTreasuryAccount", data)
	if err != nil {
		return *new(common.Address), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, consumed, err

}

// GetTreasuryFee is a free data retrieval call binding the contract method 0x29070c6d.
//
// Solidity: function getTreasuryFee() view returns(uint256)
func (_IAutonity *IAutonity) GetTreasuryFee(opts *tests.RunOptions) (*big.Int, uint64, error) {
	data, consumed, err := _IAutonity.Call(opts, "getTreasuryFee")

	if err != nil {
		return *new(*big.Int), consumed, _IAutonity.DecodeError(data, err)
	}
	out, err := _IAutonity.Contract.Abi().Unpack("getTreasuryFee", data)
	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// GetUnbondingPeriod is a free data retrieval call binding the contract method 0x6fd2c80b.
//
// Solidity: function getUnbondingPeriod() view returns(uint256)
func (_IAutonity *IAutonity) GetUnbondingPeriod(opts *tests.RunOptions) (*big.Int, uint64, error) {
	data, consumed, err := _IAutonity.Call(opts, "getUnbondingPeriod")

	if err != nil {
		return *new(*big.Int), consumed, _IAutonity.DecodeError(data, err)
	}
	out, err := _IAutonity.Contract.Abi().Unpack("getUnbondingPeriod", data)
	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// GetUnbondingRequestByID is a free data retrieval call binding the contract method 0x4bfe23f1.
//
// Solidity: function getUnbondingRequestByID(uint256 _id) view returns((address,address,uint256,uint256,uint256,bool,bool,bool))
func (_IAutonity *IAutonity) GetUnbondingRequestByID(opts *tests.RunOptions, _id *big.Int) (IAutonityUnbondingRequest, uint64, error) {
	data, consumed, err := _IAutonity.Call(opts, "getUnbondingRequestByID", _id)

	if err != nil {
		return *new(IAutonityUnbondingRequest), consumed, _IAutonity.DecodeError(data, err)
	}
	out, err := _IAutonity.Contract.Abi().Unpack("getUnbondingRequestByID", data)
	if err != nil {
		return *new(IAutonityUnbondingRequest), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(IAutonityUnbondingRequest)).(*IAutonityUnbondingRequest)
	return out0, consumed, err

}

// GetUnbondingShare is a free data retrieval call binding the contract method 0x8d347287.
//
// Solidity: function getUnbondingShare(uint256 _unbondingID) view returns(uint256)
func (_IAutonity *IAutonity) GetUnbondingShare(opts *tests.RunOptions, _unbondingID *big.Int) (*big.Int, uint64, error) {
	data, consumed, err := _IAutonity.Call(opts, "getUnbondingShare", _unbondingID)

	if err != nil {
		return *new(*big.Int), consumed, _IAutonity.DecodeError(data, err)
	}
	out, err := _IAutonity.Contract.Abi().Unpack("getUnbondingShare", data)
	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// GetValidator is a free data retrieval call binding the contract method 0x1904bb2e.
//
// Solidity: function getValidator(address _addr) view returns((address,address,address,string,uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256,address,uint256,uint256,uint256,uint256,bytes,uint8,uint256))
func (_IAutonity *IAutonity) GetValidator(opts *tests.RunOptions, _addr common.Address) (IAutonityValidator, uint64, error) {
	data, consumed, err := _IAutonity.Call(opts, "getValidator", _addr)

	if err != nil {
		return *new(IAutonityValidator), consumed, _IAutonity.DecodeError(data, err)
	}
	out, err := _IAutonity.Contract.Abi().Unpack("getValidator", data)
	if err != nil {
		return *new(IAutonityValidator), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(IAutonityValidator)).(*IAutonityValidator)
	return out0, consumed, err

}

// GetValidatorState is a free data retrieval call binding the contract method 0x5b7d6c36.
//
// Solidity: function getValidatorState(address _addr) view returns(uint8)
func (_IAutonity *IAutonity) GetValidatorState(opts *tests.RunOptions, _addr common.Address) (uint8, uint64, error) {
	data, consumed, err := _IAutonity.Call(opts, "getValidatorState", _addr)

	if err != nil {
		return *new(uint8), consumed, _IAutonity.DecodeError(data, err)
	}
	out, err := _IAutonity.Contract.Abi().Unpack("getValidatorState", data)
	if err != nil {
		return *new(uint8), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)
	return out0, consumed, err

}

// GetValidators is a free data retrieval call binding the contract method 0xb7ab4db5.
//
// Solidity: function getValidators() view returns(address[])
func (_IAutonity *IAutonity) GetValidators(opts *tests.RunOptions) ([]common.Address, uint64, error) {
	data, consumed, err := _IAutonity.Call(opts, "getValidators")

	if err != nil {
		return *new([]common.Address), consumed, _IAutonity.DecodeError(data, err)
	}
	out, err := _IAutonity.Contract.Abi().Unpack("getValidators", data)
	if err != nil {
		return *new([]common.Address), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)
	return out0, consumed, err

}

// GetVersion is a free data retrieval call binding the contract method 0x0d8e6e2c.
//
// Solidity: function getVersion() view returns(uint256)
func (_IAutonity *IAutonity) GetVersion(opts *tests.RunOptions) (*big.Int, uint64, error) {
	data, consumed, err := _IAutonity.Call(opts, "getVersion")

	if err != nil {
		return *new(*big.Int), consumed, _IAutonity.DecodeError(data, err)
	}
	out, err := _IAutonity.Contract.Abi().Unpack("getVersion", data)
	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// IsUnbondingReleased is a free data retrieval call binding the contract method 0xe294df7c.
//
// Solidity: function isUnbondingReleased(uint256 _unbondingID) view returns(bool)
func (_IAutonity *IAutonity) IsUnbondingReleased(opts *tests.RunOptions, _unbondingID *big.Int) (bool, uint64, error) {
	data, consumed, err := _IAutonity.Call(opts, "isUnbondingReleased", _unbondingID)

	if err != nil {
		return *new(bool), consumed, _IAutonity.DecodeError(data, err)
	}
	out, err := _IAutonity.Contract.Abi().Unpack("isUnbondingReleased", data)
	if err != nil {
		return *new(bool), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, consumed, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_IAutonity *IAutonity) TotalSupply(opts *tests.RunOptions) (*big.Int, uint64, error) {
	data, consumed, err := _IAutonity.Call(opts, "totalSupply")

	if err != nil {
		return *new(*big.Int), consumed, _IAutonity.DecodeError(data, err)
	}
	out, err := _IAutonity.Contract.Abi().Unpack("totalSupply", data)
	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// ActivateValidator is a free data retrieval call for a paid mutator transaction binding the contract method 0xb46e5520.
// Similar to eth_call from rpc calls or function.call from truffle, it reverts the state after the call and returns the output. The output is extracted
// the same way as done above for view only functions.
// Solidity: function activateValidator(address _address) returns()
func (_IAutonity *IAutonity) CallActivateValidator(r *tests.Runner, opts *tests.RunOptions, _address common.Address) (uint64, error) {
	snap := r.Snapshot()

	data, consumed, err := _IAutonity.Call(opts, "activateValidator", _address)
	r.RevertSnapshot(snap)
	return consumed, _IAutonity.DecodeError(data, err)

}

// Approve is a free data retrieval call for a paid mutator transaction binding the contract method 0x095ea7b3.
// Similar to eth_call from rpc calls or function.call from truffle, it reverts the state after the call and returns the output. The output is extracted
// the same way as done above for view only functions.
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_IAutonity *IAutonity) CallApprove(r *tests.Runner, opts *tests.RunOptions, spender common.Address, amount *big.Int) (bool, uint64, error) {
	snap := r.Snapshot()

	data, consumed, err := _IAutonity.Call(opts, "approve", spender, amount)
	r.RevertSnapshot(snap)

	if err != nil {
		return *new(bool), consumed, _IAutonity.DecodeError(data, err)
	}
	out, err := _IAutonity.Contract.Abi().Unpack("approve", data)
	if err != nil {
		return *new(bool), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, consumed, err

}

// ApproveBonding is a free data retrieval call for a paid mutator transaction binding the contract method 0x50492571.
// Similar to eth_call from rpc calls or function.call from truffle, it reverts the state after the call and returns the output. The output is extracted
// the same way as done above for view only functions.
// Solidity: function approveBonding(address _caller, uint256 _amount) returns(bool)
func (_IAutonity *IAutonity) CallApproveBonding(r *tests.Runner, opts *tests.RunOptions, _caller common.Address, _amount *big.Int) (bool, uint64, error) {
	snap := r.Snapshot()

	data, consumed, err := _IAutonity.Call(opts, "approveBonding", _caller, _amount)
	r.RevertSnapshot(snap)

	if err != nil {
		return *new(bool), consumed, _IAutonity.DecodeError(data, err)
	}
	out, err := _IAutonity.Contract.Abi().Unpack("approveBonding", data)
	if err != nil {
		return *new(bool), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, consumed, err

}

// Autobond is a free data retrieval call for a paid mutator transaction binding the contract method 0xf7fcc510.
// Similar to eth_call from rpc calls or function.call from truffle, it reverts the state after the call and returns the output. The output is extracted
// the same way as done above for view only functions.
// Solidity: function autobond(address _validator, uint256 _selfBond, uint256 _delegated) returns()
func (_IAutonity *IAutonity) CallAutobond(r *tests.Runner, opts *tests.RunOptions, _validator common.Address, _selfBond *big.Int, _delegated *big.Int) (uint64, error) {
	snap := r.Snapshot()

	data, consumed, err := _IAutonity.Call(opts, "autobond", _validator, _selfBond, _delegated)
	r.RevertSnapshot(snap)
	return consumed, _IAutonity.DecodeError(data, err)

}

// Bond is a free data retrieval call for a paid mutator transaction binding the contract method 0xa515366a.
// Similar to eth_call from rpc calls or function.call from truffle, it reverts the state after the call and returns the output. The output is extracted
// the same way as done above for view only functions.
// Solidity: function bond(address _validator, uint256 _amount) returns(uint256)
func (_IAutonity *IAutonity) CallBond(r *tests.Runner, opts *tests.RunOptions, _validator common.Address, _amount *big.Int) (*big.Int, uint64, error) {
	snap := r.Snapshot()

	data, consumed, err := _IAutonity.Call(opts, "bond", _validator, _amount)
	r.RevertSnapshot(snap)

	if err != nil {
		return *new(*big.Int), consumed, _IAutonity.DecodeError(data, err)
	}
	out, err := _IAutonity.Contract.Abi().Unpack("bond", data)
	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// BondFrom is a free data retrieval call for a paid mutator transaction binding the contract method 0x41de7400.
// Similar to eth_call from rpc calls or function.call from truffle, it reverts the state after the call and returns the output. The output is extracted
// the same way as done above for view only functions.
// Solidity: function bondFrom(address _account, address _validator, uint256 _amount) returns(uint256)
func (_IAutonity *IAutonity) CallBondFrom(r *tests.Runner, opts *tests.RunOptions, _account common.Address, _validator common.Address, _amount *big.Int) (*big.Int, uint64, error) {
	snap := r.Snapshot()

	data, consumed, err := _IAutonity.Call(opts, "bondFrom", _account, _validator, _amount)
	r.RevertSnapshot(snap)

	if err != nil {
		return *new(*big.Int), consumed, _IAutonity.DecodeError(data, err)
	}
	out, err := _IAutonity.Contract.Abi().Unpack("bondFrom", data)
	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// ChangeCommissionRate is a free data retrieval call for a paid mutator transaction binding the contract method 0x852c4849.
// Similar to eth_call from rpc calls or function.call from truffle, it reverts the state after the call and returns the output. The output is extracted
// the same way as done above for view only functions.
// Solidity: function changeCommissionRate(address _validator, uint256 _rate) returns()
func (_IAutonity *IAutonity) CallChangeCommissionRate(r *tests.Runner, opts *tests.RunOptions, _validator common.Address, _rate *big.Int) (uint64, error) {
	snap := r.Snapshot()

	data, consumed, err := _IAutonity.Call(opts, "changeCommissionRate", _validator, _rate)
	r.RevertSnapshot(snap)
	return consumed, _IAutonity.DecodeError(data, err)

}

// Jail is a free data retrieval call for a paid mutator transaction binding the contract method 0x154d76d7.
// Similar to eth_call from rpc calls or function.call from truffle, it reverts the state after the call and returns the output. The output is extracted
// the same way as done above for view only functions.
// Solidity: function jail(address _nodeAddress, uint256 _jailtime, uint8 _newJailedState) returns(uint256)
func (_IAutonity *IAutonity) CallJail(r *tests.Runner, opts *tests.RunOptions, _nodeAddress common.Address, _jailtime *big.Int, _newJailedState uint8) (*big.Int, uint64, error) {
	snap := r.Snapshot()

	data, consumed, err := _IAutonity.Call(opts, "jail", _nodeAddress, _jailtime, _newJailedState)
	r.RevertSnapshot(snap)

	if err != nil {
		return *new(*big.Int), consumed, _IAutonity.DecodeError(data, err)
	}
	out, err := _IAutonity.Contract.Abi().Unpack("jail", data)
	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// Jailbound is a free data retrieval call for a paid mutator transaction binding the contract method 0x8ef8c2fd.
// Similar to eth_call from rpc calls or function.call from truffle, it reverts the state after the call and returns the output. The output is extracted
// the same way as done above for view only functions.
// Solidity: function jailbound(address _nodeAddress, uint8 _newJailboundState) returns()
func (_IAutonity *IAutonity) CallJailbound(r *tests.Runner, opts *tests.RunOptions, _nodeAddress common.Address, _newJailboundState uint8) (uint64, error) {
	snap := r.Snapshot()

	data, consumed, err := _IAutonity.Call(opts, "jailbound", _nodeAddress, _newJailboundState)
	r.RevertSnapshot(snap)
	return consumed, _IAutonity.DecodeError(data, err)

}

// PauseValidator is a free data retrieval call for a paid mutator transaction binding the contract method 0x0ae65e7a.
// Similar to eth_call from rpc calls or function.call from truffle, it reverts the state after the call and returns the output. The output is extracted
// the same way as done above for view only functions.
// Solidity: function pauseValidator(address _address) returns()
func (_IAutonity *IAutonity) CallPauseValidator(r *tests.Runner, opts *tests.RunOptions, _address common.Address) (uint64, error) {
	snap := r.Snapshot()

	data, consumed, err := _IAutonity.Call(opts, "pauseValidator", _address)
	r.RevertSnapshot(snap)
	return consumed, _IAutonity.DecodeError(data, err)

}

// RegisterValidator is a free data retrieval call for a paid mutator transaction binding the contract method 0x84467fdb.
// Similar to eth_call from rpc calls or function.call from truffle, it reverts the state after the call and returns the output. The output is extracted
// the same way as done above for view only functions.
// Solidity: function registerValidator(string _enode, address _oracleAddress, bytes _consensusKey, bytes _signatures) returns()
func (_IAutonity *IAutonity) CallRegisterValidator(r *tests.Runner, opts *tests.RunOptions, _enode string, _oracleAddress common.Address, _consensusKey []byte, _signatures []byte) (uint64, error) {
	snap := r.Snapshot()

	data, consumed, err := _IAutonity.Call(opts, "registerValidator", _enode, _oracleAddress, _consensusKey, _signatures)
	r.RevertSnapshot(snap)
	return consumed, _IAutonity.DecodeError(data, err)

}

// Slash is a free data retrieval call for a paid mutator transaction binding the contract method 0x02fb4d85.
// Similar to eth_call from rpc calls or function.call from truffle, it reverts the state after the call and returns the output. The output is extracted
// the same way as done above for view only functions.
// Solidity: function slash(address _nodeAddress, uint256 _slashingRate) returns(uint256 slashingAmount)
func (_IAutonity *IAutonity) CallSlash(r *tests.Runner, opts *tests.RunOptions, _nodeAddress common.Address, _slashingRate *big.Int) (*big.Int, uint64, error) {
	snap := r.Snapshot()

	data, consumed, err := _IAutonity.Call(opts, "slash", _nodeAddress, _slashingRate)
	r.RevertSnapshot(snap)

	if err != nil {
		return *new(*big.Int), consumed, _IAutonity.DecodeError(data, err)
	}
	out, err := _IAutonity.Contract.Abi().Unpack("slash", data)
	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// SlashAndJail is a free data retrieval call for a paid mutator transaction binding the contract method 0x122b4122.
// Similar to eth_call from rpc calls or function.call from truffle, it reverts the state after the call and returns the output. The output is extracted
// the same way as done above for view only functions.
// Solidity: function slashAndJail(address _nodeAddress, uint256 _slashingRate, uint256 _jailtime, uint8 _newJailedState, uint8 _newJailboundState) returns(uint256 slashingAmount, uint256 jailReleaseBlock, bool isJailbound)
func (_IAutonity *IAutonity) CallSlashAndJail(r *tests.Runner, opts *tests.RunOptions, _nodeAddress common.Address, _slashingRate *big.Int, _jailtime *big.Int, _newJailedState uint8, _newJailboundState uint8) (struct {
	SlashingAmount   *big.Int
	JailReleaseBlock *big.Int
	IsJailbound      bool
}, uint64, error) {
	snap := r.Snapshot()

	data, consumed, err := _IAutonity.Call(opts, "slashAndJail", _nodeAddress, _slashingRate, _jailtime, _newJailedState, _newJailboundState)
	r.RevertSnapshot(snap)

	outstruct := new(struct {
		SlashingAmount   *big.Int
		JailReleaseBlock *big.Int
		IsJailbound      bool
	})
	if err != nil {
		return *outstruct, consumed, _IAutonity.DecodeError(data, err)
	}
	out, err := _IAutonity.Contract.Abi().Unpack("slashAndJail", data)
	if err != nil {
		return *outstruct, consumed, err
	}

	outstruct.SlashingAmount = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.JailReleaseBlock = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.IsJailbound = *abi.ConvertType(out[2], new(bool)).(*bool)
	return *outstruct, consumed, err

}

// Transfer is a free data retrieval call for a paid mutator transaction binding the contract method 0xa9059cbb.
// Similar to eth_call from rpc calls or function.call from truffle, it reverts the state after the call and returns the output. The output is extracted
// the same way as done above for view only functions.
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_IAutonity *IAutonity) CallTransfer(r *tests.Runner, opts *tests.RunOptions, recipient common.Address, amount *big.Int) (bool, uint64, error) {
	snap := r.Snapshot()

	data, consumed, err := _IAutonity.Call(opts, "transfer", recipient, amount)
	r.RevertSnapshot(snap)

	if err != nil {
		return *new(bool), consumed, _IAutonity.DecodeError(data, err)
	}
	out, err := _IAutonity.Contract.Abi().Unpack("transfer", data)
	if err != nil {
		return *new(bool), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, consumed, err

}

// TransferFrom is a free data retrieval call for a paid mutator transaction binding the contract method 0x23b872dd.
// Similar to eth_call from rpc calls or function.call from truffle, it reverts the state after the call and returns the output. The output is extracted
// the same way as done above for view only functions.
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_IAutonity *IAutonity) CallTransferFrom(r *tests.Runner, opts *tests.RunOptions, sender common.Address, recipient common.Address, amount *big.Int) (bool, uint64, error) {
	snap := r.Snapshot()

	data, consumed, err := _IAutonity.Call(opts, "transferFrom", sender, recipient, amount)
	r.RevertSnapshot(snap)

	if err != nil {
		return *new(bool), consumed, _IAutonity.DecodeError(data, err)
	}
	out, err := _IAutonity.Contract.Abi().Unpack("transferFrom", data)
	if err != nil {
		return *new(bool), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, consumed, err

}

// Unbond is a free data retrieval call for a paid mutator transaction binding the contract method 0xa5d059ca.
// Similar to eth_call from rpc calls or function.call from truffle, it reverts the state after the call and returns the output. The output is extracted
// the same way as done above for view only functions.
// Solidity: function unbond(address _validator, uint256 _amount) returns(uint256)
func (_IAutonity *IAutonity) CallUnbond(r *tests.Runner, opts *tests.RunOptions, _validator common.Address, _amount *big.Int) (*big.Int, uint64, error) {
	snap := r.Snapshot()

	data, consumed, err := _IAutonity.Call(opts, "unbond", _validator, _amount)
	r.RevertSnapshot(snap)

	if err != nil {
		return *new(*big.Int), consumed, _IAutonity.DecodeError(data, err)
	}
	out, err := _IAutonity.Contract.Abi().Unpack("unbond", data)
	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// UnbondFrom is a free data retrieval call for a paid mutator transaction binding the contract method 0xa9a7d7c9.
// Similar to eth_call from rpc calls or function.call from truffle, it reverts the state after the call and returns the output. The output is extracted
// the same way as done above for view only functions.
// Solidity: function unbondFrom(address _account, address _validator, uint256 _amount) returns(uint256)
func (_IAutonity *IAutonity) CallUnbondFrom(r *tests.Runner, opts *tests.RunOptions, _account common.Address, _validator common.Address, _amount *big.Int) (*big.Int, uint64, error) {
	snap := r.Snapshot()

	data, consumed, err := _IAutonity.Call(opts, "unbondFrom", _account, _validator, _amount)
	r.RevertSnapshot(snap)

	if err != nil {
		return *new(*big.Int), consumed, _IAutonity.DecodeError(data, err)
	}
	out, err := _IAutonity.Contract.Abi().Unpack("unbondFrom", data)
	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// UpdateEnode is a free data retrieval call for a paid mutator transaction binding the contract method 0x784304b5.
// Similar to eth_call from rpc calls or function.call from truffle, it reverts the state after the call and returns the output. The output is extracted
// the same way as done above for view only functions.
// Solidity: function updateEnode(address _nodeAddress, string _enode) returns()
func (_IAutonity *IAutonity) CallUpdateEnode(r *tests.Runner, opts *tests.RunOptions, _nodeAddress common.Address, _enode string) (uint64, error) {
	snap := r.Snapshot()

	data, consumed, err := _IAutonity.Call(opts, "updateEnode", _nodeAddress, _enode)
	r.RevertSnapshot(snap)
	return consumed, _IAutonity.DecodeError(data, err)

}

// ActivateValidator is a paid mutator transaction binding the contract method 0xb46e5520.
//
// Solidity: function activateValidator(address _address) returns()
func (_IAutonity *IAutonity) ActivateValidator(opts *tests.RunOptions, _address common.Address) (uint64, error) {
	data, consumed, err := _IAutonity.Call(opts, "activateValidator", _address)
	return consumed, _IAutonity.DecodeError(data, err)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_IAutonity *IAutonity) Approve(opts *tests.RunOptions, spender common.Address, amount *big.Int) (uint64, error) {
	data, consumed, err := _IAutonity.Call(opts, "approve", spender, amount)
	return consumed, _IAutonity.DecodeError(data, err)
}

// ApproveBonding is a paid mutator transaction binding the contract method 0x50492571.
//
// Solidity: function approveBonding(address _caller, uint256 _amount) returns(bool)
func (_IAutonity *IAutonity) ApproveBonding(opts *tests.RunOptions, _caller common.Address, _amount *big.Int) (uint64, error) {
	data, consumed, err := _IAutonity.Call(opts, "approveBonding", _caller, _amount)
	return consumed, _IAutonity.DecodeError(data, err)
}

// Autobond is a paid mutator transaction binding the contract method 0xf7fcc510.
//
// Solidity: function autobond(address _validator, uint256 _selfBond, uint256 _delegated) returns()
func (_IAutonity *IAutonity) Autobond(opts *tests.RunOptions, _validator common.Address, _selfBond *big.Int, _delegated *big.Int) (uint64, error) {
	data, consumed, err := _IAutonity.Call(opts, "autobond", _validator, _selfBond, _delegated)
	return consumed, _IAutonity.DecodeError(data, err)
}

// Bond is a paid mutator transaction binding the contract method 0xa515366a.
//
// Solidity: function bond(address _validator, uint256 _amount) returns(uint256)
func (_IAutonity *IAutonity) Bond(opts *tests.RunOptions, _validator common.Address, _amount *big.Int) (uint64, error) {
	data, consumed, err := _IAutonity.Call(opts, "bond", _validator, _amount)
	return consumed, _IAutonity.DecodeError(data, err)
}

// BondFrom is a paid mutator transaction binding the contract method 0x41de7400.
//
// Solidity: function bondFrom(address _account, address _validator, uint256 _amount) returns(uint256)
func (_IAutonity *IAutonity) BondFrom(opts *tests.RunOptions, _account common.Address, _validator common.Address, _amount *big.Int) (uint64, error) {
	data, consumed, err := _IAutonity.Call(opts, "bondFrom", _account, _validator, _amount)
	return consumed, _IAutonity.DecodeError(data, err)
}

// ChangeCommissionRate is a paid mutator transaction binding the contract method 0x852c4849.
//
// Solidity: function changeCommissionRate(address _validator, uint256 _rate) returns()
func (_IAutonity *IAutonity) ChangeCommissionRate(opts *tests.RunOptions, _validator common.Address, _rate *big.Int) (uint64, error) {
	data, consumed, err := _IAutonity.Call(opts, "changeCommissionRate", _validator, _rate)
	return consumed, _IAutonity.DecodeError(data, err)
}

// Jail is a paid mutator transaction binding the contract method 0x154d76d7.
//
// Solidity: function jail(address _nodeAddress, uint256 _jailtime, uint8 _newJailedState) returns(uint256)
func (_IAutonity *IAutonity) Jail(opts *tests.RunOptions, _nodeAddress common.Address, _jailtime *big.Int, _newJailedState uint8) (uint64, error) {
	data, consumed, err := _IAutonity.Call(opts, "jail", _nodeAddress, _jailtime, _newJailedState)
	return consumed, _IAutonity.DecodeError(data, err)
}

// Jailbound is a paid mutator transaction binding the contract method 0x8ef8c2fd.
//
// Solidity: function jailbound(address _nodeAddress, uint8 _newJailboundState) returns()
func (_IAutonity *IAutonity) Jailbound(opts *tests.RunOptions, _nodeAddress common.Address, _newJailboundState uint8) (uint64, error) {
	data, consumed, err := _IAutonity.Call(opts, "jailbound", _nodeAddress, _newJailboundState)
	return consumed, _IAutonity.DecodeError(data, err)
}

// PauseValidator is a paid mutator transaction binding the contract method 0x0ae65e7a.
//
// Solidity: function pauseValidator(address _address) returns()
func (_IAutonity *IAutonity) PauseValidator(opts *tests.RunOptions, _address common.Address) (uint64, error) {
	data, consumed, err := _IAutonity.Call(opts, "pauseValidator", _address)
	return consumed, _IAutonity.DecodeError(data, err)
}

// RegisterValidator is a paid mutator transaction binding the contract method 0x84467fdb.
//
// Solidity: function registerValidator(string _enode, address _oracleAddress, bytes _consensusKey, bytes _signatures) returns()
func (_IAutonity *IAutonity) RegisterValidator(opts *tests.RunOptions, _enode string, _oracleAddress common.Address, _consensusKey []byte, _signatures []byte) (uint64, error) {
	data, consumed, err := _IAutonity.Call(opts, "registerValidator", _enode, _oracleAddress, _consensusKey, _signatures)
	return consumed, _IAutonity.DecodeError(data, err)
}

// Slash is a paid mutator transaction binding the contract method 0x02fb4d85.
//
// Solidity: function slash(address _nodeAddress, uint256 _slashingRate) returns(uint256 slashingAmount)
func (_IAutonity *IAutonity) Slash(opts *tests.RunOptions, _nodeAddress common.Address, _slashingRate *big.Int) (uint64, error) {
	data, consumed, err := _IAutonity.Call(opts, "slash", _nodeAddress, _slashingRate)
	return consumed, _IAutonity.DecodeError(data, err)
}

// SlashAndJail is a paid mutator transaction binding the contract method 0x122b4122.
//
// Solidity: function slashAndJail(address _nodeAddress, uint256 _slashingRate, uint256 _jailtime, uint8 _newJailedState, uint8 _newJailboundState) returns(uint256 slashingAmount, uint256 jailReleaseBlock, bool isJailbound)
func (_IAutonity *IAutonity) SlashAndJail(opts *tests.RunOptions, _nodeAddress common.Address, _slashingRate *big.Int, _jailtime *big.Int, _newJailedState uint8, _newJailboundState uint8) (uint64, error) {
	data, consumed, err := _IAutonity.Call(opts, "slashAndJail", _nodeAddress, _slashingRate, _jailtime, _newJailedState, _newJailboundState)
	return consumed, _IAutonity.DecodeError(data, err)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_IAutonity *IAutonity) Transfer(opts *tests.RunOptions, recipient common.Address, amount *big.Int) (uint64, error) {
	data, consumed, err := _IAutonity.Call(opts, "transfer", recipient, amount)
	return consumed, _IAutonity.DecodeError(data, err)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_IAutonity *IAutonity) TransferFrom(opts *tests.RunOptions, sender common.Address, recipient common.Address, amount *big.Int) (uint64, error) {
	data, consumed, err := _IAutonity.Call(opts, "transferFrom", sender, recipient, amount)
	return consumed, _IAutonity.DecodeError(data, err)
}

// Unbond is a paid mutator transaction binding the contract method 0xa5d059ca.
//
// Solidity: function unbond(address _validator, uint256 _amount) returns(uint256)
func (_IAutonity *IAutonity) Unbond(opts *tests.RunOptions, _validator common.Address, _amount *big.Int) (uint64, error) {
	data, consumed, err := _IAutonity.Call(opts, "unbond", _validator, _amount)
	return consumed, _IAutonity.DecodeError(data, err)
}

// UnbondFrom is a paid mutator transaction binding the contract method 0xa9a7d7c9.
//
// Solidity: function unbondFrom(address _account, address _validator, uint256 _amount) returns(uint256)
func (_IAutonity *IAutonity) UnbondFrom(opts *tests.RunOptions, _account common.Address, _validator common.Address, _amount *big.Int) (uint64, error) {
	data, consumed, err := _IAutonity.Call(opts, "unbondFrom", _account, _validator, _amount)
	return consumed, _IAutonity.DecodeError(data, err)
}

// UpdateEnode is a paid mutator transaction binding the contract method 0x784304b5.
//
// Solidity: function updateEnode(address _nodeAddress, string _enode) returns()
func (_IAutonity *IAutonity) UpdateEnode(opts *tests.RunOptions, _nodeAddress common.Address, _enode string) (uint64, error) {
	data, consumed, err := _IAutonity.Call(opts, "updateEnode", _nodeAddress, _enode)
	return consumed, _IAutonity.DecodeError(data, err)
}

func (_IAutonity *IAutonity) DecodeError(data []byte, err error) error {
	if err == nil {
		return nil
	}

	reason, _ := abi.UnpackRevert(data)
	return fmt.Errorf("%w: %s", err, reason)

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
	*tests.Contract
}

func (_IConfigEvents *IConfigEvents) DecodeError(data []byte, err error) error {
	if err == nil {
		return nil
	}

	reason, _ := abi.UnpackRevert(data)
	return fmt.Errorf("%w: %s", err, reason)

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
	*tests.Contract
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_IERC20 *IERC20) Allowance(opts *tests.RunOptions, owner common.Address, spender common.Address) (*big.Int, uint64, error) {
	data, consumed, err := _IERC20.Call(opts, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), consumed, _IERC20.DecodeError(data, err)
	}
	out, err := _IERC20.Contract.Abi().Unpack("allowance", data)
	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_IERC20 *IERC20) BalanceOf(opts *tests.RunOptions, account common.Address) (*big.Int, uint64, error) {
	data, consumed, err := _IERC20.Call(opts, "balanceOf", account)

	if err != nil {
		return *new(*big.Int), consumed, _IERC20.DecodeError(data, err)
	}
	out, err := _IERC20.Contract.Abi().Unpack("balanceOf", data)
	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_IERC20 *IERC20) TotalSupply(opts *tests.RunOptions) (*big.Int, uint64, error) {
	data, consumed, err := _IERC20.Call(opts, "totalSupply")

	if err != nil {
		return *new(*big.Int), consumed, _IERC20.DecodeError(data, err)
	}
	out, err := _IERC20.Contract.Abi().Unpack("totalSupply", data)
	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// Approve is a free data retrieval call for a paid mutator transaction binding the contract method 0x095ea7b3.
// Similar to eth_call from rpc calls or function.call from truffle, it reverts the state after the call and returns the output. The output is extracted
// the same way as done above for view only functions.
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_IERC20 *IERC20) CallApprove(r *tests.Runner, opts *tests.RunOptions, spender common.Address, amount *big.Int) (bool, uint64, error) {
	snap := r.Snapshot()

	data, consumed, err := _IERC20.Call(opts, "approve", spender, amount)
	r.RevertSnapshot(snap)

	if err != nil {
		return *new(bool), consumed, _IERC20.DecodeError(data, err)
	}
	out, err := _IERC20.Contract.Abi().Unpack("approve", data)
	if err != nil {
		return *new(bool), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, consumed, err

}

// Transfer is a free data retrieval call for a paid mutator transaction binding the contract method 0xa9059cbb.
// Similar to eth_call from rpc calls or function.call from truffle, it reverts the state after the call and returns the output. The output is extracted
// the same way as done above for view only functions.
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_IERC20 *IERC20) CallTransfer(r *tests.Runner, opts *tests.RunOptions, recipient common.Address, amount *big.Int) (bool, uint64, error) {
	snap := r.Snapshot()

	data, consumed, err := _IERC20.Call(opts, "transfer", recipient, amount)
	r.RevertSnapshot(snap)

	if err != nil {
		return *new(bool), consumed, _IERC20.DecodeError(data, err)
	}
	out, err := _IERC20.Contract.Abi().Unpack("transfer", data)
	if err != nil {
		return *new(bool), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, consumed, err

}

// TransferFrom is a free data retrieval call for a paid mutator transaction binding the contract method 0x23b872dd.
// Similar to eth_call from rpc calls or function.call from truffle, it reverts the state after the call and returns the output. The output is extracted
// the same way as done above for view only functions.
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_IERC20 *IERC20) CallTransferFrom(r *tests.Runner, opts *tests.RunOptions, sender common.Address, recipient common.Address, amount *big.Int) (bool, uint64, error) {
	snap := r.Snapshot()

	data, consumed, err := _IERC20.Call(opts, "transferFrom", sender, recipient, amount)
	r.RevertSnapshot(snap)

	if err != nil {
		return *new(bool), consumed, _IERC20.DecodeError(data, err)
	}
	out, err := _IERC20.Contract.Abi().Unpack("transferFrom", data)
	if err != nil {
		return *new(bool), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, consumed, err

}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_IERC20 *IERC20) Approve(opts *tests.RunOptions, spender common.Address, amount *big.Int) (uint64, error) {
	data, consumed, err := _IERC20.Call(opts, "approve", spender, amount)
	return consumed, _IERC20.DecodeError(data, err)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_IERC20 *IERC20) Transfer(opts *tests.RunOptions, recipient common.Address, amount *big.Int) (uint64, error) {
	data, consumed, err := _IERC20.Call(opts, "transfer", recipient, amount)
	return consumed, _IERC20.DecodeError(data, err)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_IERC20 *IERC20) TransferFrom(opts *tests.RunOptions, sender common.Address, recipient common.Address, amount *big.Int) (uint64, error) {
	data, consumed, err := _IERC20.Call(opts, "transferFrom", sender, recipient, amount)
	return consumed, _IERC20.DecodeError(data, err)
}

func (_IERC20 *IERC20) DecodeError(data []byte, err error) error {
	if err == nil {
		return nil
	}

	reason, _ := abi.UnpackRevert(data)
	return fmt.Errorf("%w: %s", err, reason)

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
	*tests.Contract
}

// CalculateSupplyDelta is a free data retrieval call binding the contract method 0x92eff3cd.
//
// Solidity: function calculateSupplyDelta(uint256 _currentSupply, uint256 _inflationReserve, uint256 _lastEpochTime, uint256 _currentEpochTime) view returns(uint256)
func (_IInflationController *IInflationController) CalculateSupplyDelta(opts *tests.RunOptions, _currentSupply *big.Int, _inflationReserve *big.Int, _lastEpochTime *big.Int, _currentEpochTime *big.Int) (*big.Int, uint64, error) {
	data, consumed, err := _IInflationController.Call(opts, "calculateSupplyDelta", _currentSupply, _inflationReserve, _lastEpochTime, _currentEpochTime)

	if err != nil {
		return *new(*big.Int), consumed, _IInflationController.DecodeError(data, err)
	}
	out, err := _IInflationController.Contract.Abi().Unpack("calculateSupplyDelta", data)
	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

func (_IInflationController *IInflationController) DecodeError(data []byte, err error) error {
	if err == nil {
		return nil
	}

	reason, _ := abi.UnpackRevert(data)
	return fmt.Errorf("%w: %s", err, reason)

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
	*tests.Contract
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_ILiquid *ILiquid) Allowance(opts *tests.RunOptions, owner common.Address, spender common.Address) (*big.Int, uint64, error) {
	data, consumed, err := _ILiquid.Call(opts, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), consumed, _ILiquid.DecodeError(data, err)
	}
	out, err := _ILiquid.Contract.Abi().Unpack("allowance", data)
	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_ILiquid *ILiquid) BalanceOf(opts *tests.RunOptions, account common.Address) (*big.Int, uint64, error) {
	data, consumed, err := _ILiquid.Call(opts, "balanceOf", account)

	if err != nil {
		return *new(*big.Int), consumed, _ILiquid.DecodeError(data, err)
	}
	out, err := _ILiquid.Contract.Abi().Unpack("balanceOf", data)
	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() pure returns(uint8)
func (_ILiquid *ILiquid) Decimals(opts *tests.RunOptions) (uint8, uint64, error) {
	data, consumed, err := _ILiquid.Call(opts, "decimals")

	if err != nil {
		return *new(uint8), consumed, _ILiquid.DecodeError(data, err)
	}
	out, err := _ILiquid.Contract.Abi().Unpack("decimals", data)
	if err != nil {
		return *new(uint8), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)
	return out0, consumed, err

}

// GetCommissionRate is a free data retrieval call binding the contract method 0x3e4eb36c.
//
// Solidity: function getCommissionRate() view returns(uint256)
func (_ILiquid *ILiquid) GetCommissionRate(opts *tests.RunOptions) (*big.Int, uint64, error) {
	data, consumed, err := _ILiquid.Call(opts, "getCommissionRate")

	if err != nil {
		return *new(*big.Int), consumed, _ILiquid.DecodeError(data, err)
	}
	out, err := _ILiquid.Contract.Abi().Unpack("getCommissionRate", data)
	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// GetTreasury is a free data retrieval call binding the contract method 0x3b19e84a.
//
// Solidity: function getTreasury() view returns(address)
func (_ILiquid *ILiquid) GetTreasury(opts *tests.RunOptions) (common.Address, uint64, error) {
	data, consumed, err := _ILiquid.Call(opts, "getTreasury")

	if err != nil {
		return *new(common.Address), consumed, _ILiquid.DecodeError(data, err)
	}
	out, err := _ILiquid.Contract.Abi().Unpack("getTreasury", data)
	if err != nil {
		return *new(common.Address), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, consumed, err

}

// GetTreasuryUnclaimedATN is a free data retrieval call binding the contract method 0x1eeffad0.
//
// Solidity: function getTreasuryUnclaimedATN() view returns(uint256)
func (_ILiquid *ILiquid) GetTreasuryUnclaimedATN(opts *tests.RunOptions) (*big.Int, uint64, error) {
	data, consumed, err := _ILiquid.Call(opts, "getTreasuryUnclaimedATN")

	if err != nil {
		return *new(*big.Int), consumed, _ILiquid.DecodeError(data, err)
	}
	out, err := _ILiquid.Contract.Abi().Unpack("getTreasuryUnclaimedATN", data)
	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// GetValidator is a free data retrieval call binding the contract method 0x1195e07e.
//
// Solidity: function getValidator() view returns(address)
func (_ILiquid *ILiquid) GetValidator(opts *tests.RunOptions) (common.Address, uint64, error) {
	data, consumed, err := _ILiquid.Call(opts, "getValidator")

	if err != nil {
		return *new(common.Address), consumed, _ILiquid.DecodeError(data, err)
	}
	out, err := _ILiquid.Contract.Abi().Unpack("getValidator", data)
	if err != nil {
		return *new(common.Address), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, consumed, err

}

// LockedBalanceOf is a free data retrieval call binding the contract method 0x59355736.
//
// Solidity: function lockedBalanceOf(address _delegator) view returns(uint256)
func (_ILiquid *ILiquid) LockedBalanceOf(opts *tests.RunOptions, _delegator common.Address) (*big.Int, uint64, error) {
	data, consumed, err := _ILiquid.Call(opts, "lockedBalanceOf", _delegator)

	if err != nil {
		return *new(*big.Int), consumed, _ILiquid.DecodeError(data, err)
	}
	out, err := _ILiquid.Contract.Abi().Unpack("lockedBalanceOf", data)
	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_ILiquid *ILiquid) Name(opts *tests.RunOptions) (string, uint64, error) {
	data, consumed, err := _ILiquid.Call(opts, "name")

	if err != nil {
		return *new(string), consumed, _ILiquid.DecodeError(data, err)
	}
	out, err := _ILiquid.Contract.Abi().Unpack("name", data)
	if err != nil {
		return *new(string), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)
	return out0, consumed, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_ILiquid *ILiquid) Symbol(opts *tests.RunOptions) (string, uint64, error) {
	data, consumed, err := _ILiquid.Call(opts, "symbol")

	if err != nil {
		return *new(string), consumed, _ILiquid.DecodeError(data, err)
	}
	out, err := _ILiquid.Contract.Abi().Unpack("symbol", data)
	if err != nil {
		return *new(string), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)
	return out0, consumed, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_ILiquid *ILiquid) TotalSupply(opts *tests.RunOptions) (*big.Int, uint64, error) {
	data, consumed, err := _ILiquid.Call(opts, "totalSupply")

	if err != nil {
		return *new(*big.Int), consumed, _ILiquid.DecodeError(data, err)
	}
	out, err := _ILiquid.Contract.Abi().Unpack("totalSupply", data)
	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// UnbondingAllowance is a free data retrieval call binding the contract method 0xd768d578.
//
// Solidity: function unbondingAllowance(address _owner, address _caller) view returns(uint256)
func (_ILiquid *ILiquid) UnbondingAllowance(opts *tests.RunOptions, _owner common.Address, _caller common.Address) (*big.Int, uint64, error) {
	data, consumed, err := _ILiquid.Call(opts, "unbondingAllowance", _owner, _caller)

	if err != nil {
		return *new(*big.Int), consumed, _ILiquid.DecodeError(data, err)
	}
	out, err := _ILiquid.Contract.Abi().Unpack("unbondingAllowance", data)
	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// UnclaimedRewards is a free data retrieval call binding the contract method 0x949813b8.
//
// Solidity: function unclaimedRewards(address _account) view returns(uint256)
func (_ILiquid *ILiquid) UnclaimedRewards(opts *tests.RunOptions, _account common.Address) (*big.Int, uint64, error) {
	data, consumed, err := _ILiquid.Call(opts, "unclaimedRewards", _account)

	if err != nil {
		return *new(*big.Int), consumed, _ILiquid.DecodeError(data, err)
	}
	out, err := _ILiquid.Contract.Abi().Unpack("unclaimedRewards", data)
	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// UnlockedBalanceOf is a free data retrieval call binding the contract method 0x84955c88.
//
// Solidity: function unlockedBalanceOf(address _delegator) view returns(uint256)
func (_ILiquid *ILiquid) UnlockedBalanceOf(opts *tests.RunOptions, _delegator common.Address) (*big.Int, uint64, error) {
	data, consumed, err := _ILiquid.Call(opts, "unlockedBalanceOf", _delegator)

	if err != nil {
		return *new(*big.Int), consumed, _ILiquid.DecodeError(data, err)
	}
	out, err := _ILiquid.Contract.Abi().Unpack("unlockedBalanceOf", data)
	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// Approve is a free data retrieval call for a paid mutator transaction binding the contract method 0x095ea7b3.
// Similar to eth_call from rpc calls or function.call from truffle, it reverts the state after the call and returns the output. The output is extracted
// the same way as done above for view only functions.
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_ILiquid *ILiquid) CallApprove(r *tests.Runner, opts *tests.RunOptions, spender common.Address, amount *big.Int) (bool, uint64, error) {
	snap := r.Snapshot()

	data, consumed, err := _ILiquid.Call(opts, "approve", spender, amount)
	r.RevertSnapshot(snap)

	if err != nil {
		return *new(bool), consumed, _ILiquid.DecodeError(data, err)
	}
	out, err := _ILiquid.Contract.Abi().Unpack("approve", data)
	if err != nil {
		return *new(bool), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, consumed, err

}

// ApproveUnbonding is a free data retrieval call for a paid mutator transaction binding the contract method 0xbf99b73e.
// Similar to eth_call from rpc calls or function.call from truffle, it reverts the state after the call and returns the output. The output is extracted
// the same way as done above for view only functions.
// Solidity: function approveUnbonding(address _caller, uint256 _amount) returns(bool)
func (_ILiquid *ILiquid) CallApproveUnbonding(r *tests.Runner, opts *tests.RunOptions, _caller common.Address, _amount *big.Int) (bool, uint64, error) {
	snap := r.Snapshot()

	data, consumed, err := _ILiquid.Call(opts, "approveUnbonding", _caller, _amount)
	r.RevertSnapshot(snap)

	if err != nil {
		return *new(bool), consumed, _ILiquid.DecodeError(data, err)
	}
	out, err := _ILiquid.Contract.Abi().Unpack("approveUnbonding", data)
	if err != nil {
		return *new(bool), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, consumed, err

}

// Burn is a free data retrieval call for a paid mutator transaction binding the contract method 0x9dc29fac.
// Similar to eth_call from rpc calls or function.call from truffle, it reverts the state after the call and returns the output. The output is extracted
// the same way as done above for view only functions.
// Solidity: function burn(address _account, uint256 _amount) returns()
func (_ILiquid *ILiquid) CallBurn(r *tests.Runner, opts *tests.RunOptions, _account common.Address, _amount *big.Int) (uint64, error) {
	snap := r.Snapshot()

	data, consumed, err := _ILiquid.Call(opts, "burn", _account, _amount)
	r.RevertSnapshot(snap)
	return consumed, _ILiquid.DecodeError(data, err)

}

// ClaimRewards is a free data retrieval call for a paid mutator transaction binding the contract method 0x372500ab.
// Similar to eth_call from rpc calls or function.call from truffle, it reverts the state after the call and returns the output. The output is extracted
// the same way as done above for view only functions.
// Solidity: function claimRewards() returns()
func (_ILiquid *ILiquid) CallClaimRewards(r *tests.Runner, opts *tests.RunOptions) (uint64, error) {
	snap := r.Snapshot()

	data, consumed, err := _ILiquid.Call(opts, "claimRewards")
	r.RevertSnapshot(snap)
	return consumed, _ILiquid.DecodeError(data, err)

}

// ClaimTreasuryATN is a free data retrieval call for a paid mutator transaction binding the contract method 0xbd96102f.
// Similar to eth_call from rpc calls or function.call from truffle, it reverts the state after the call and returns the output. The output is extracted
// the same way as done above for view only functions.
// Solidity: function claimTreasuryATN() returns()
func (_ILiquid *ILiquid) CallClaimTreasuryATN(r *tests.Runner, opts *tests.RunOptions) (uint64, error) {
	snap := r.Snapshot()

	data, consumed, err := _ILiquid.Call(opts, "claimTreasuryATN")
	r.RevertSnapshot(snap)
	return consumed, _ILiquid.DecodeError(data, err)

}

// Lock is a free data retrieval call for a paid mutator transaction binding the contract method 0x282d3fdf.
// Similar to eth_call from rpc calls or function.call from truffle, it reverts the state after the call and returns the output. The output is extracted
// the same way as done above for view only functions.
// Solidity: function lock(address _account, uint256 _amount) returns()
func (_ILiquid *ILiquid) CallLock(r *tests.Runner, opts *tests.RunOptions, _account common.Address, _amount *big.Int) (uint64, error) {
	snap := r.Snapshot()

	data, consumed, err := _ILiquid.Call(opts, "lock", _account, _amount)
	r.RevertSnapshot(snap)
	return consumed, _ILiquid.DecodeError(data, err)

}

// LockFrom is a free data retrieval call for a paid mutator transaction binding the contract method 0x708e91e5.
// Similar to eth_call from rpc calls or function.call from truffle, it reverts the state after the call and returns the output. The output is extracted
// the same way as done above for view only functions.
// Solidity: function lockFrom(address _account, address _caller, uint256 _amount) returns()
func (_ILiquid *ILiquid) CallLockFrom(r *tests.Runner, opts *tests.RunOptions, _account common.Address, _caller common.Address, _amount *big.Int) (uint64, error) {
	snap := r.Snapshot()

	data, consumed, err := _ILiquid.Call(opts, "lockFrom", _account, _caller, _amount)
	r.RevertSnapshot(snap)
	return consumed, _ILiquid.DecodeError(data, err)

}

// Mint is a free data retrieval call for a paid mutator transaction binding the contract method 0x40c10f19.
// Similar to eth_call from rpc calls or function.call from truffle, it reverts the state after the call and returns the output. The output is extracted
// the same way as done above for view only functions.
// Solidity: function mint(address _account, uint256 _amount) returns()
func (_ILiquid *ILiquid) CallMint(r *tests.Runner, opts *tests.RunOptions, _account common.Address, _amount *big.Int) (uint64, error) {
	snap := r.Snapshot()

	data, consumed, err := _ILiquid.Call(opts, "mint", _account, _amount)
	r.RevertSnapshot(snap)
	return consumed, _ILiquid.DecodeError(data, err)

}

// Redistribute is a free data retrieval call for a paid mutator transaction binding the contract method 0xa0ce552d.
// Similar to eth_call from rpc calls or function.call from truffle, it reverts the state after the call and returns the output. The output is extracted
// the same way as done above for view only functions.
// Solidity: function redistribute(uint256 _ntnReward) payable returns(uint256)
func (_ILiquid *ILiquid) CallRedistribute(r *tests.Runner, opts *tests.RunOptions, _ntnReward *big.Int) (*big.Int, uint64, error) {
	snap := r.Snapshot()

	data, consumed, err := _ILiquid.Call(opts, "redistribute", _ntnReward)
	r.RevertSnapshot(snap)

	if err != nil {
		return *new(*big.Int), consumed, _ILiquid.DecodeError(data, err)
	}
	out, err := _ILiquid.Contract.Abi().Unpack("redistribute", data)
	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// SetCommissionRate is a free data retrieval call for a paid mutator transaction binding the contract method 0x19fac8fd.
// Similar to eth_call from rpc calls or function.call from truffle, it reverts the state after the call and returns the output. The output is extracted
// the same way as done above for view only functions.
// Solidity: function setCommissionRate(uint256 _rate) returns()
func (_ILiquid *ILiquid) CallSetCommissionRate(r *tests.Runner, opts *tests.RunOptions, _rate *big.Int) (uint64, error) {
	snap := r.Snapshot()

	data, consumed, err := _ILiquid.Call(opts, "setCommissionRate", _rate)
	r.RevertSnapshot(snap)
	return consumed, _ILiquid.DecodeError(data, err)

}

// Transfer is a free data retrieval call for a paid mutator transaction binding the contract method 0xa9059cbb.
// Similar to eth_call from rpc calls or function.call from truffle, it reverts the state after the call and returns the output. The output is extracted
// the same way as done above for view only functions.
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_ILiquid *ILiquid) CallTransfer(r *tests.Runner, opts *tests.RunOptions, recipient common.Address, amount *big.Int) (bool, uint64, error) {
	snap := r.Snapshot()

	data, consumed, err := _ILiquid.Call(opts, "transfer", recipient, amount)
	r.RevertSnapshot(snap)

	if err != nil {
		return *new(bool), consumed, _ILiquid.DecodeError(data, err)
	}
	out, err := _ILiquid.Contract.Abi().Unpack("transfer", data)
	if err != nil {
		return *new(bool), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, consumed, err

}

// TransferFrom is a free data retrieval call for a paid mutator transaction binding the contract method 0x23b872dd.
// Similar to eth_call from rpc calls or function.call from truffle, it reverts the state after the call and returns the output. The output is extracted
// the same way as done above for view only functions.
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_ILiquid *ILiquid) CallTransferFrom(r *tests.Runner, opts *tests.RunOptions, sender common.Address, recipient common.Address, amount *big.Int) (bool, uint64, error) {
	snap := r.Snapshot()

	data, consumed, err := _ILiquid.Call(opts, "transferFrom", sender, recipient, amount)
	r.RevertSnapshot(snap)

	if err != nil {
		return *new(bool), consumed, _ILiquid.DecodeError(data, err)
	}
	out, err := _ILiquid.Contract.Abi().Unpack("transferFrom", data)
	if err != nil {
		return *new(bool), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, consumed, err

}

// Unlock is a free data retrieval call for a paid mutator transaction binding the contract method 0x7eee288d.
// Similar to eth_call from rpc calls or function.call from truffle, it reverts the state after the call and returns the output. The output is extracted
// the same way as done above for view only functions.
// Solidity: function unlock(address _account, uint256 _amount) returns()
func (_ILiquid *ILiquid) CallUnlock(r *tests.Runner, opts *tests.RunOptions, _account common.Address, _amount *big.Int) (uint64, error) {
	snap := r.Snapshot()

	data, consumed, err := _ILiquid.Call(opts, "unlock", _account, _amount)
	r.RevertSnapshot(snap)
	return consumed, _ILiquid.DecodeError(data, err)

}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_ILiquid *ILiquid) Approve(opts *tests.RunOptions, spender common.Address, amount *big.Int) (uint64, error) {
	data, consumed, err := _ILiquid.Call(opts, "approve", spender, amount)
	return consumed, _ILiquid.DecodeError(data, err)
}

// ApproveUnbonding is a paid mutator transaction binding the contract method 0xbf99b73e.
//
// Solidity: function approveUnbonding(address _caller, uint256 _amount) returns(bool)
func (_ILiquid *ILiquid) ApproveUnbonding(opts *tests.RunOptions, _caller common.Address, _amount *big.Int) (uint64, error) {
	data, consumed, err := _ILiquid.Call(opts, "approveUnbonding", _caller, _amount)
	return consumed, _ILiquid.DecodeError(data, err)
}

// Burn is a paid mutator transaction binding the contract method 0x9dc29fac.
//
// Solidity: function burn(address _account, uint256 _amount) returns()
func (_ILiquid *ILiquid) Burn(opts *tests.RunOptions, _account common.Address, _amount *big.Int) (uint64, error) {
	data, consumed, err := _ILiquid.Call(opts, "burn", _account, _amount)
	return consumed, _ILiquid.DecodeError(data, err)
}

// ClaimRewards is a paid mutator transaction binding the contract method 0x372500ab.
//
// Solidity: function claimRewards() returns()
func (_ILiquid *ILiquid) ClaimRewards(opts *tests.RunOptions) (uint64, error) {
	data, consumed, err := _ILiquid.Call(opts, "claimRewards")
	return consumed, _ILiquid.DecodeError(data, err)
}

// ClaimTreasuryATN is a paid mutator transaction binding the contract method 0xbd96102f.
//
// Solidity: function claimTreasuryATN() returns()
func (_ILiquid *ILiquid) ClaimTreasuryATN(opts *tests.RunOptions) (uint64, error) {
	data, consumed, err := _ILiquid.Call(opts, "claimTreasuryATN")
	return consumed, _ILiquid.DecodeError(data, err)
}

// Lock is a paid mutator transaction binding the contract method 0x282d3fdf.
//
// Solidity: function lock(address _account, uint256 _amount) returns()
func (_ILiquid *ILiquid) Lock(opts *tests.RunOptions, _account common.Address, _amount *big.Int) (uint64, error) {
	data, consumed, err := _ILiquid.Call(opts, "lock", _account, _amount)
	return consumed, _ILiquid.DecodeError(data, err)
}

// LockFrom is a paid mutator transaction binding the contract method 0x708e91e5.
//
// Solidity: function lockFrom(address _account, address _caller, uint256 _amount) returns()
func (_ILiquid *ILiquid) LockFrom(opts *tests.RunOptions, _account common.Address, _caller common.Address, _amount *big.Int) (uint64, error) {
	data, consumed, err := _ILiquid.Call(opts, "lockFrom", _account, _caller, _amount)
	return consumed, _ILiquid.DecodeError(data, err)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address _account, uint256 _amount) returns()
func (_ILiquid *ILiquid) Mint(opts *tests.RunOptions, _account common.Address, _amount *big.Int) (uint64, error) {
	data, consumed, err := _ILiquid.Call(opts, "mint", _account, _amount)
	return consumed, _ILiquid.DecodeError(data, err)
}

// Redistribute is a paid mutator transaction binding the contract method 0xa0ce552d.
//
// Solidity: function redistribute(uint256 _ntnReward) payable returns(uint256)
func (_ILiquid *ILiquid) Redistribute(opts *tests.RunOptions, _ntnReward *big.Int) (uint64, error) {
	data, consumed, err := _ILiquid.Call(opts, "redistribute", _ntnReward)
	return consumed, _ILiquid.DecodeError(data, err)
}

// SetCommissionRate is a paid mutator transaction binding the contract method 0x19fac8fd.
//
// Solidity: function setCommissionRate(uint256 _rate) returns()
func (_ILiquid *ILiquid) SetCommissionRate(opts *tests.RunOptions, _rate *big.Int) (uint64, error) {
	data, consumed, err := _ILiquid.Call(opts, "setCommissionRate", _rate)
	return consumed, _ILiquid.DecodeError(data, err)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_ILiquid *ILiquid) Transfer(opts *tests.RunOptions, recipient common.Address, amount *big.Int) (uint64, error) {
	data, consumed, err := _ILiquid.Call(opts, "transfer", recipient, amount)
	return consumed, _ILiquid.DecodeError(data, err)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_ILiquid *ILiquid) TransferFrom(opts *tests.RunOptions, sender common.Address, recipient common.Address, amount *big.Int) (uint64, error) {
	data, consumed, err := _ILiquid.Call(opts, "transferFrom", sender, recipient, amount)
	return consumed, _ILiquid.DecodeError(data, err)
}

// Unlock is a paid mutator transaction binding the contract method 0x7eee288d.
//
// Solidity: function unlock(address _account, uint256 _amount) returns()
func (_ILiquid *ILiquid) Unlock(opts *tests.RunOptions, _account common.Address, _amount *big.Int) (uint64, error) {
	data, consumed, err := _ILiquid.Call(opts, "unlock", _account, _amount)
	return consumed, _ILiquid.DecodeError(data, err)
}

func (_ILiquid *ILiquid) DecodeError(data []byte, err error) error {
	if err == nil {
		return nil
	}

	reason, _ := abi.UnpackRevert(data)
	return fmt.Errorf("%w: %s", err, reason)

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
	*tests.Contract
}

// GetDelta is a free data retrieval call binding the contract method 0xc549176e.
//
// Solidity: function getDelta() view returns(uint256)
func (_IOmissionAccountability *IOmissionAccountability) GetDelta(opts *tests.RunOptions) (*big.Int, uint64, error) {
	data, consumed, err := _IOmissionAccountability.Call(opts, "getDelta")

	if err != nil {
		return *new(*big.Int), consumed, _IOmissionAccountability.DecodeError(data, err)
	}
	out, err := _IOmissionAccountability.Contract.Abi().Unpack("getDelta", data)
	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// GetInactivityScore is a free data retrieval call binding the contract method 0x9a11e0e6.
//
// Solidity: function getInactivityScore(address _validator) view returns(uint256)
func (_IOmissionAccountability *IOmissionAccountability) GetInactivityScore(opts *tests.RunOptions, _validator common.Address) (*big.Int, uint64, error) {
	data, consumed, err := _IOmissionAccountability.Call(opts, "getInactivityScore", _validator)

	if err != nil {
		return *new(*big.Int), consumed, _IOmissionAccountability.DecodeError(data, err)
	}
	out, err := _IOmissionAccountability.Contract.Abi().Unpack("getInactivityScore", data)
	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// GetLookbackWindow is a free data retrieval call binding the contract method 0x5ca1809c.
//
// Solidity: function getLookbackWindow() view returns(uint256)
func (_IOmissionAccountability *IOmissionAccountability) GetLookbackWindow(opts *tests.RunOptions) (*big.Int, uint64, error) {
	data, consumed, err := _IOmissionAccountability.Call(opts, "getLookbackWindow")

	if err != nil {
		return *new(*big.Int), consumed, _IOmissionAccountability.DecodeError(data, err)
	}
	out, err := _IOmissionAccountability.Contract.Abi().Unpack("getLookbackWindow", data)
	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// GetScaleFactor is a free data retrieval call binding the contract method 0x7f5e2f11.
//
// Solidity: function getScaleFactor() pure returns(uint256)
func (_IOmissionAccountability *IOmissionAccountability) GetScaleFactor(opts *tests.RunOptions) (*big.Int, uint64, error) {
	data, consumed, err := _IOmissionAccountability.Call(opts, "getScaleFactor")

	if err != nil {
		return *new(*big.Int), consumed, _IOmissionAccountability.DecodeError(data, err)
	}
	out, err := _IOmissionAccountability.Contract.Abi().Unpack("getScaleFactor", data)
	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// GetTotalEffort is a free data retrieval call binding the contract method 0x53b1821b.
//
// Solidity: function getTotalEffort() view returns(uint256)
func (_IOmissionAccountability *IOmissionAccountability) GetTotalEffort(opts *tests.RunOptions) (*big.Int, uint64, error) {
	data, consumed, err := _IOmissionAccountability.Call(opts, "getTotalEffort")

	if err != nil {
		return *new(*big.Int), consumed, _IOmissionAccountability.DecodeError(data, err)
	}
	out, err := _IOmissionAccountability.Contract.Abi().Unpack("getTotalEffort", data)
	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// DistributeProposerRewards is a free data retrieval call for a paid mutator transaction binding the contract method 0xeeb92233.
// Similar to eth_call from rpc calls or function.call from truffle, it reverts the state after the call and returns the output. The output is extracted
// the same way as done above for view only functions.
// Solidity: function distributeProposerRewards(uint256 _ntnReward) payable returns()
func (_IOmissionAccountability *IOmissionAccountability) CallDistributeProposerRewards(r *tests.Runner, opts *tests.RunOptions, _ntnReward *big.Int) (uint64, error) {
	snap := r.Snapshot()

	data, consumed, err := _IOmissionAccountability.Call(opts, "distributeProposerRewards", _ntnReward)
	r.RevertSnapshot(snap)
	return consumed, _IOmissionAccountability.DecodeError(data, err)

}

// Finalize is a free data retrieval call for a paid mutator transaction binding the contract method 0x6c9789b0.
// Similar to eth_call from rpc calls or function.call from truffle, it reverts the state after the call and returns the output. The output is extracted
// the same way as done above for view only functions.
// Solidity: function finalize(bool _epochEnded) returns(uint256)
func (_IOmissionAccountability *IOmissionAccountability) CallFinalize(r *tests.Runner, opts *tests.RunOptions, _epochEnded bool) (*big.Int, uint64, error) {
	snap := r.Snapshot()

	data, consumed, err := _IOmissionAccountability.Call(opts, "finalize", _epochEnded)
	r.RevertSnapshot(snap)

	if err != nil {
		return *new(*big.Int), consumed, _IOmissionAccountability.DecodeError(data, err)
	}
	out, err := _IOmissionAccountability.Contract.Abi().Unpack("finalize", data)
	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// SetCommittee is a free data retrieval call for a paid mutator transaction binding the contract method 0xe3deef9c.
// Similar to eth_call from rpc calls or function.call from truffle, it reverts the state after the call and returns the output. The output is extracted
// the same way as done above for view only functions.
// Solidity: function setCommittee((address,uint256,bytes)[] _committee, address[] _treasuries) returns()
func (_IOmissionAccountability *IOmissionAccountability) CallSetCommittee(r *tests.Runner, opts *tests.RunOptions, _committee []IAutonityCommitteeMember, _treasuries []common.Address) (uint64, error) {
	snap := r.Snapshot()

	data, consumed, err := _IOmissionAccountability.Call(opts, "setCommittee", _committee, _treasuries)
	r.RevertSnapshot(snap)
	return consumed, _IOmissionAccountability.DecodeError(data, err)

}

// SetEpochBlock is a free data retrieval call for a paid mutator transaction binding the contract method 0xc024cc2c.
// Similar to eth_call from rpc calls or function.call from truffle, it reverts the state after the call and returns the output. The output is extracted
// the same way as done above for view only functions.
// Solidity: function setEpochBlock(uint256 _epochBlock) returns()
func (_IOmissionAccountability *IOmissionAccountability) CallSetEpochBlock(r *tests.Runner, opts *tests.RunOptions, _epochBlock *big.Int) (uint64, error) {
	snap := r.Snapshot()

	data, consumed, err := _IOmissionAccountability.Call(opts, "setEpochBlock", _epochBlock)
	r.RevertSnapshot(snap)
	return consumed, _IOmissionAccountability.DecodeError(data, err)

}

// SetOperator is a free data retrieval call for a paid mutator transaction binding the contract method 0xb3ab15fb.
// Similar to eth_call from rpc calls or function.call from truffle, it reverts the state after the call and returns the output. The output is extracted
// the same way as done above for view only functions.
// Solidity: function setOperator(address _operator) returns()
func (_IOmissionAccountability *IOmissionAccountability) CallSetOperator(r *tests.Runner, opts *tests.RunOptions, _operator common.Address) (uint64, error) {
	snap := r.Snapshot()

	data, consumed, err := _IOmissionAccountability.Call(opts, "setOperator", _operator)
	r.RevertSnapshot(snap)
	return consumed, _IOmissionAccountability.DecodeError(data, err)

}

// DistributeProposerRewards is a paid mutator transaction binding the contract method 0xeeb92233.
//
// Solidity: function distributeProposerRewards(uint256 _ntnReward) payable returns()
func (_IOmissionAccountability *IOmissionAccountability) DistributeProposerRewards(opts *tests.RunOptions, _ntnReward *big.Int) (uint64, error) {
	data, consumed, err := _IOmissionAccountability.Call(opts, "distributeProposerRewards", _ntnReward)
	return consumed, _IOmissionAccountability.DecodeError(data, err)
}

// Finalize is a paid mutator transaction binding the contract method 0x6c9789b0.
//
// Solidity: function finalize(bool _epochEnded) returns(uint256)
func (_IOmissionAccountability *IOmissionAccountability) Finalize(opts *tests.RunOptions, _epochEnded bool) (uint64, error) {
	data, consumed, err := _IOmissionAccountability.Call(opts, "finalize", _epochEnded)
	return consumed, _IOmissionAccountability.DecodeError(data, err)
}

// SetCommittee is a paid mutator transaction binding the contract method 0xe3deef9c.
//
// Solidity: function setCommittee((address,uint256,bytes)[] _committee, address[] _treasuries) returns()
func (_IOmissionAccountability *IOmissionAccountability) SetCommittee(opts *tests.RunOptions, _committee []IAutonityCommitteeMember, _treasuries []common.Address) (uint64, error) {
	data, consumed, err := _IOmissionAccountability.Call(opts, "setCommittee", _committee, _treasuries)
	return consumed, _IOmissionAccountability.DecodeError(data, err)
}

// SetEpochBlock is a paid mutator transaction binding the contract method 0xc024cc2c.
//
// Solidity: function setEpochBlock(uint256 _epochBlock) returns()
func (_IOmissionAccountability *IOmissionAccountability) SetEpochBlock(opts *tests.RunOptions, _epochBlock *big.Int) (uint64, error) {
	data, consumed, err := _IOmissionAccountability.Call(opts, "setEpochBlock", _epochBlock)
	return consumed, _IOmissionAccountability.DecodeError(data, err)
}

// SetOperator is a paid mutator transaction binding the contract method 0xb3ab15fb.
//
// Solidity: function setOperator(address _operator) returns()
func (_IOmissionAccountability *IOmissionAccountability) SetOperator(opts *tests.RunOptions, _operator common.Address) (uint64, error) {
	data, consumed, err := _IOmissionAccountability.Call(opts, "setOperator", _operator)
	return consumed, _IOmissionAccountability.DecodeError(data, err)
}

func (_IOmissionAccountability *IOmissionAccountability) DecodeError(data []byte, err error) error {
	if err == nil {
		return nil
	}

	reason, _ := abi.UnpackRevert(data)
	return fmt.Errorf("%w: %s", err, reason)

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
	*tests.Contract
}

// GetDecimals is a free data retrieval call binding the contract method 0xf0141d84.
//
// Solidity: function getDecimals() view returns(uint8)
func (_IOracle *IOracle) GetDecimals(opts *tests.RunOptions) (uint8, uint64, error) {
	data, consumed, err := _IOracle.Call(opts, "getDecimals")

	if err != nil {
		return *new(uint8), consumed, _IOracle.DecodeError(data, err)
	}
	out, err := _IOracle.Contract.Abi().Unpack("getDecimals", data)
	if err != nil {
		return *new(uint8), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)
	return out0, consumed, err

}

// GetNewVotePeriod is a free data retrieval call binding the contract method 0x57eba759.
//
// Solidity: function getNewVotePeriod() view returns(uint256)
func (_IOracle *IOracle) GetNewVotePeriod(opts *tests.RunOptions) (*big.Int, uint64, error) {
	data, consumed, err := _IOracle.Call(opts, "getNewVotePeriod")

	if err != nil {
		return *new(*big.Int), consumed, _IOracle.DecodeError(data, err)
	}
	out, err := _IOracle.Contract.Abi().Unpack("getNewVotePeriod", data)
	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// GetNewVoters is a free data retrieval call binding the contract method 0x077945d3.
//
// Solidity: function getNewVoters() view returns(address[])
func (_IOracle *IOracle) GetNewVoters(opts *tests.RunOptions) ([]common.Address, uint64, error) {
	data, consumed, err := _IOracle.Call(opts, "getNewVoters")

	if err != nil {
		return *new([]common.Address), consumed, _IOracle.DecodeError(data, err)
	}
	out, err := _IOracle.Contract.Abi().Unpack("getNewVoters", data)
	if err != nil {
		return *new([]common.Address), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)
	return out0, consumed, err

}

// GetNonRevealThreshold is a free data retrieval call binding the contract method 0xed78349d.
//
// Solidity: function getNonRevealThreshold() view returns(uint256)
func (_IOracle *IOracle) GetNonRevealThreshold(opts *tests.RunOptions) (*big.Int, uint64, error) {
	data, consumed, err := _IOracle.Call(opts, "getNonRevealThreshold")

	if err != nil {
		return *new(*big.Int), consumed, _IOracle.DecodeError(data, err)
	}
	out, err := _IOracle.Contract.Abi().Unpack("getNonRevealThreshold", data)
	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// GetRound is a free data retrieval call binding the contract method 0x9f8743f7.
//
// Solidity: function getRound() view returns(uint256)
func (_IOracle *IOracle) GetRound(opts *tests.RunOptions) (*big.Int, uint64, error) {
	data, consumed, err := _IOracle.Call(opts, "getRound")

	if err != nil {
		return *new(*big.Int), consumed, _IOracle.DecodeError(data, err)
	}
	out, err := _IOracle.Contract.Abi().Unpack("getRound", data)
	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// GetRoundData is a free data retrieval call binding the contract method 0x3c8510fd.
//
// Solidity: function getRoundData(uint256 _round, string _symbol) view returns((uint256,uint256,uint256,bool) data)
func (_IOracle *IOracle) GetRoundData(opts *tests.RunOptions, _round *big.Int, _symbol string) (IOracleRoundData, uint64, error) {
	data, consumed, err := _IOracle.Call(opts, "getRoundData", _round, _symbol)

	if err != nil {
		return *new(IOracleRoundData), consumed, _IOracle.DecodeError(data, err)
	}
	out, err := _IOracle.Contract.Abi().Unpack("getRoundData", data)
	if err != nil {
		return *new(IOracleRoundData), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(IOracleRoundData)).(*IOracleRoundData)
	return out0, consumed, err

}

// GetSymbols is a free data retrieval call binding the contract method 0xdf7f710e.
//
// Solidity: function getSymbols() view returns(string[] _symbols)
func (_IOracle *IOracle) GetSymbols(opts *tests.RunOptions) ([]string, uint64, error) {
	data, consumed, err := _IOracle.Call(opts, "getSymbols")

	if err != nil {
		return *new([]string), consumed, _IOracle.DecodeError(data, err)
	}
	out, err := _IOracle.Contract.Abi().Unpack("getSymbols", data)
	if err != nil {
		return *new([]string), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new([]string)).(*[]string)
	return out0, consumed, err

}

// GetVotePeriod is a free data retrieval call binding the contract method 0xb78dec52.
//
// Solidity: function getVotePeriod() view returns(uint256)
func (_IOracle *IOracle) GetVotePeriod(opts *tests.RunOptions) (*big.Int, uint64, error) {
	data, consumed, err := _IOracle.Call(opts, "getVotePeriod")

	if err != nil {
		return *new(*big.Int), consumed, _IOracle.DecodeError(data, err)
	}
	out, err := _IOracle.Contract.Abi().Unpack("getVotePeriod", data)
	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// GetVoters is a free data retrieval call binding the contract method 0xcdd72253.
//
// Solidity: function getVoters() view returns(address[])
func (_IOracle *IOracle) GetVoters(opts *tests.RunOptions) ([]common.Address, uint64, error) {
	data, consumed, err := _IOracle.Call(opts, "getVoters")

	if err != nil {
		return *new([]common.Address), consumed, _IOracle.DecodeError(data, err)
	}
	out, err := _IOracle.Contract.Abi().Unpack("getVoters", data)
	if err != nil {
		return *new([]common.Address), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)
	return out0, consumed, err

}

// LatestRoundData is a free data retrieval call binding the contract method 0x33f98c77.
//
// Solidity: function latestRoundData(string _symbol) view returns((uint256,uint256,uint256,bool) data)
func (_IOracle *IOracle) LatestRoundData(opts *tests.RunOptions, _symbol string) (IOracleRoundData, uint64, error) {
	data, consumed, err := _IOracle.Call(opts, "latestRoundData", _symbol)

	if err != nil {
		return *new(IOracleRoundData), consumed, _IOracle.DecodeError(data, err)
	}
	out, err := _IOracle.Contract.Abi().Unpack("latestRoundData", data)
	if err != nil {
		return *new(IOracleRoundData), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(IOracleRoundData)).(*IOracleRoundData)
	return out0, consumed, err

}

// DistributeRewards is a free data retrieval call for a paid mutator transaction binding the contract method 0x59974e38.
// Similar to eth_call from rpc calls or function.call from truffle, it reverts the state after the call and returns the output. The output is extracted
// the same way as done above for view only functions.
// Solidity: function distributeRewards(uint256 _ntnRewards) payable returns()
func (_IOracle *IOracle) CallDistributeRewards(r *tests.Runner, opts *tests.RunOptions, _ntnRewards *big.Int) (uint64, error) {
	snap := r.Snapshot()

	data, consumed, err := _IOracle.Call(opts, "distributeRewards", _ntnRewards)
	r.RevertSnapshot(snap)
	return consumed, _IOracle.DecodeError(data, err)

}

// Finalize is a free data retrieval call for a paid mutator transaction binding the contract method 0x4bb278f3.
// Similar to eth_call from rpc calls or function.call from truffle, it reverts the state after the call and returns the output. The output is extracted
// the same way as done above for view only functions.
// Solidity: function finalize() returns(bool)
func (_IOracle *IOracle) CallFinalize(r *tests.Runner, opts *tests.RunOptions) (bool, uint64, error) {
	snap := r.Snapshot()

	data, consumed, err := _IOracle.Call(opts, "finalize")
	r.RevertSnapshot(snap)

	if err != nil {
		return *new(bool), consumed, _IOracle.DecodeError(data, err)
	}
	out, err := _IOracle.Contract.Abi().Unpack("finalize", data)
	if err != nil {
		return *new(bool), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, consumed, err

}

// SetCommitRevealConfig is a free data retrieval call for a paid mutator transaction binding the contract method 0x3f422ef3.
// Similar to eth_call from rpc calls or function.call from truffle, it reverts the state after the call and returns the output. The output is extracted
// the same way as done above for view only functions.
// Solidity: function setCommitRevealConfig(uint256 _threshold, uint256 _resetInterval) returns()
func (_IOracle *IOracle) CallSetCommitRevealConfig(r *tests.Runner, opts *tests.RunOptions, _threshold *big.Int, _resetInterval *big.Int) (uint64, error) {
	snap := r.Snapshot()

	data, consumed, err := _IOracle.Call(opts, "setCommitRevealConfig", _threshold, _resetInterval)
	r.RevertSnapshot(snap)
	return consumed, _IOracle.DecodeError(data, err)

}

// SetOperator is a free data retrieval call for a paid mutator transaction binding the contract method 0xb3ab15fb.
// Similar to eth_call from rpc calls or function.call from truffle, it reverts the state after the call and returns the output. The output is extracted
// the same way as done above for view only functions.
// Solidity: function setOperator(address _operator) returns()
func (_IOracle *IOracle) CallSetOperator(r *tests.Runner, opts *tests.RunOptions, _operator common.Address) (uint64, error) {
	snap := r.Snapshot()

	data, consumed, err := _IOracle.Call(opts, "setOperator", _operator)
	r.RevertSnapshot(snap)
	return consumed, _IOracle.DecodeError(data, err)

}

// SetSlashingConfig is a free data retrieval call for a paid mutator transaction binding the contract method 0xda39fbfe.
// Similar to eth_call from rpc calls or function.call from truffle, it reverts the state after the call and returns the output. The output is extracted
// the same way as done above for view only functions.
// Solidity: function setSlashingConfig(int256 _outlierSlashingThreshold, int256 _outlierDetectionThreshold, uint256 _baseSlashingRate, uint256 _slashingRateCap) returns()
func (_IOracle *IOracle) CallSetSlashingConfig(r *tests.Runner, opts *tests.RunOptions, _outlierSlashingThreshold *big.Int, _outlierDetectionThreshold *big.Int, _baseSlashingRate *big.Int, _slashingRateCap *big.Int) (uint64, error) {
	snap := r.Snapshot()

	data, consumed, err := _IOracle.Call(opts, "setSlashingConfig", _outlierSlashingThreshold, _outlierDetectionThreshold, _baseSlashingRate, _slashingRateCap)
	r.RevertSnapshot(snap)
	return consumed, _IOracle.DecodeError(data, err)

}

// SetSymbols is a free data retrieval call for a paid mutator transaction binding the contract method 0x8d4f75d2.
// Similar to eth_call from rpc calls or function.call from truffle, it reverts the state after the call and returns the output. The output is extracted
// the same way as done above for view only functions.
// Solidity: function setSymbols(string[] _symbols) returns()
func (_IOracle *IOracle) CallSetSymbols(r *tests.Runner, opts *tests.RunOptions, _symbols []string) (uint64, error) {
	snap := r.Snapshot()

	data, consumed, err := _IOracle.Call(opts, "setSymbols", _symbols)
	r.RevertSnapshot(snap)
	return consumed, _IOracle.DecodeError(data, err)

}

// SetVoters is a free data retrieval call for a paid mutator transaction binding the contract method 0xda78110e.
// Similar to eth_call from rpc calls or function.call from truffle, it reverts the state after the call and returns the output. The output is extracted
// the same way as done above for view only functions.
// Solidity: function setVoters(address[] _newVoters, address[] _treasury, address[] _validator) returns()
func (_IOracle *IOracle) CallSetVoters(r *tests.Runner, opts *tests.RunOptions, _newVoters []common.Address, _treasury []common.Address, _validator []common.Address) (uint64, error) {
	snap := r.Snapshot()

	data, consumed, err := _IOracle.Call(opts, "setVoters", _newVoters, _treasury, _validator)
	r.RevertSnapshot(snap)
	return consumed, _IOracle.DecodeError(data, err)

}

// UpdateVotersAndSymbol is a free data retrieval call for a paid mutator transaction binding the contract method 0x0f65875c.
// Similar to eth_call from rpc calls or function.call from truffle, it reverts the state after the call and returns the output. The output is extracted
// the same way as done above for view only functions.
// Solidity: function updateVotersAndSymbol() returns()
func (_IOracle *IOracle) CallUpdateVotersAndSymbol(r *tests.Runner, opts *tests.RunOptions) (uint64, error) {
	snap := r.Snapshot()

	data, consumed, err := _IOracle.Call(opts, "updateVotersAndSymbol")
	r.RevertSnapshot(snap)
	return consumed, _IOracle.DecodeError(data, err)

}

// Vote is a free data retrieval call for a paid mutator transaction binding the contract method 0x56833ebe.
// Similar to eth_call from rpc calls or function.call from truffle, it reverts the state after the call and returns the output. The output is extracted
// the same way as done above for view only functions.
// Solidity: function vote(uint256 _commit, (uint120,uint8)[] _reports, uint256 _salt, uint8 _extra) returns()
func (_IOracle *IOracle) CallVote(r *tests.Runner, opts *tests.RunOptions, _commit *big.Int, _reports []IOracleReport, _salt *big.Int, _extra uint8) (uint64, error) {
	snap := r.Snapshot()

	data, consumed, err := _IOracle.Call(opts, "vote", _commit, _reports, _salt, _extra)
	r.RevertSnapshot(snap)
	return consumed, _IOracle.DecodeError(data, err)

}

// DistributeRewards is a paid mutator transaction binding the contract method 0x59974e38.
//
// Solidity: function distributeRewards(uint256 _ntnRewards) payable returns()
func (_IOracle *IOracle) DistributeRewards(opts *tests.RunOptions, _ntnRewards *big.Int) (uint64, error) {
	data, consumed, err := _IOracle.Call(opts, "distributeRewards", _ntnRewards)
	return consumed, _IOracle.DecodeError(data, err)
}

// Finalize is a paid mutator transaction binding the contract method 0x4bb278f3.
//
// Solidity: function finalize() returns(bool)
func (_IOracle *IOracle) Finalize(opts *tests.RunOptions) (uint64, error) {
	data, consumed, err := _IOracle.Call(opts, "finalize")
	return consumed, _IOracle.DecodeError(data, err)
}

// SetCommitRevealConfig is a paid mutator transaction binding the contract method 0x3f422ef3.
//
// Solidity: function setCommitRevealConfig(uint256 _threshold, uint256 _resetInterval) returns()
func (_IOracle *IOracle) SetCommitRevealConfig(opts *tests.RunOptions, _threshold *big.Int, _resetInterval *big.Int) (uint64, error) {
	data, consumed, err := _IOracle.Call(opts, "setCommitRevealConfig", _threshold, _resetInterval)
	return consumed, _IOracle.DecodeError(data, err)
}

// SetOperator is a paid mutator transaction binding the contract method 0xb3ab15fb.
//
// Solidity: function setOperator(address _operator) returns()
func (_IOracle *IOracle) SetOperator(opts *tests.RunOptions, _operator common.Address) (uint64, error) {
	data, consumed, err := _IOracle.Call(opts, "setOperator", _operator)
	return consumed, _IOracle.DecodeError(data, err)
}

// SetSlashingConfig is a paid mutator transaction binding the contract method 0xda39fbfe.
//
// Solidity: function setSlashingConfig(int256 _outlierSlashingThreshold, int256 _outlierDetectionThreshold, uint256 _baseSlashingRate, uint256 _slashingRateCap) returns()
func (_IOracle *IOracle) SetSlashingConfig(opts *tests.RunOptions, _outlierSlashingThreshold *big.Int, _outlierDetectionThreshold *big.Int, _baseSlashingRate *big.Int, _slashingRateCap *big.Int) (uint64, error) {
	data, consumed, err := _IOracle.Call(opts, "setSlashingConfig", _outlierSlashingThreshold, _outlierDetectionThreshold, _baseSlashingRate, _slashingRateCap)
	return consumed, _IOracle.DecodeError(data, err)
}

// SetSymbols is a paid mutator transaction binding the contract method 0x8d4f75d2.
//
// Solidity: function setSymbols(string[] _symbols) returns()
func (_IOracle *IOracle) SetSymbols(opts *tests.RunOptions, _symbols []string) (uint64, error) {
	data, consumed, err := _IOracle.Call(opts, "setSymbols", _symbols)
	return consumed, _IOracle.DecodeError(data, err)
}

// SetVoters is a paid mutator transaction binding the contract method 0xda78110e.
//
// Solidity: function setVoters(address[] _newVoters, address[] _treasury, address[] _validator) returns()
func (_IOracle *IOracle) SetVoters(opts *tests.RunOptions, _newVoters []common.Address, _treasury []common.Address, _validator []common.Address) (uint64, error) {
	data, consumed, err := _IOracle.Call(opts, "setVoters", _newVoters, _treasury, _validator)
	return consumed, _IOracle.DecodeError(data, err)
}

// UpdateVotersAndSymbol is a paid mutator transaction binding the contract method 0x0f65875c.
//
// Solidity: function updateVotersAndSymbol() returns()
func (_IOracle *IOracle) UpdateVotersAndSymbol(opts *tests.RunOptions) (uint64, error) {
	data, consumed, err := _IOracle.Call(opts, "updateVotersAndSymbol")
	return consumed, _IOracle.DecodeError(data, err)
}

// Vote is a paid mutator transaction binding the contract method 0x56833ebe.
//
// Solidity: function vote(uint256 _commit, (uint120,uint8)[] _reports, uint256 _salt, uint8 _extra) returns()
func (_IOracle *IOracle) Vote(opts *tests.RunOptions, _commit *big.Int, _reports []IOracleReport, _salt *big.Int, _extra uint8) (uint64, error) {
	data, consumed, err := _IOracle.Call(opts, "vote", _commit, _reports, _salt, _extra)
	return consumed, _IOracle.DecodeError(data, err)
}

func (_IOracle *IOracle) DecodeError(data []byte, err error) error {
	if err == nil {
		return nil
	}

	reason, _ := abi.UnpackRevert(data)
	return fmt.Errorf("%w: %s", err, reason)

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
	*tests.Contract
}

// GetSchedule is a free data retrieval call binding the contract method 0x7264c4da.
//
// Solidity: function getSchedule(address _vault, uint256 _id) view returns((uint256,uint256,uint256,uint256,uint256))
func (_IScheduleController *IScheduleController) GetSchedule(opts *tests.RunOptions, _vault common.Address, _id *big.Int) (IScheduleControllerSchedule, uint64, error) {
	data, consumed, err := _IScheduleController.Call(opts, "getSchedule", _vault, _id)

	if err != nil {
		return *new(IScheduleControllerSchedule), consumed, _IScheduleController.DecodeError(data, err)
	}
	out, err := _IScheduleController.Contract.Abi().Unpack("getSchedule", data)
	if err != nil {
		return *new(IScheduleControllerSchedule), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(IScheduleControllerSchedule)).(*IScheduleControllerSchedule)
	return out0, consumed, err

}

// GetTotalSchedules is a free data retrieval call binding the contract method 0x088566e9.
//
// Solidity: function getTotalSchedules(address _vault) view returns(uint256)
func (_IScheduleController *IScheduleController) GetTotalSchedules(opts *tests.RunOptions, _vault common.Address) (*big.Int, uint64, error) {
	data, consumed, err := _IScheduleController.Call(opts, "getTotalSchedules", _vault)

	if err != nil {
		return *new(*big.Int), consumed, _IScheduleController.DecodeError(data, err)
	}
	out, err := _IScheduleController.Contract.Abi().Unpack("getTotalSchedules", data)
	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

func (_IScheduleController *IScheduleController) DecodeError(data []byte, err error) error {
	if err == nil {
		return nil
	}

	reason, _ := abi.UnpackRevert(data)
	return fmt.Errorf("%w: %s", err, reason)

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
	*tests.Contract
}

// Cdps is a free data retrieval call binding the contract method 0x840c7e24.
//
// Solidity: function cdps(address owner) view returns((uint256,uint256,uint256,uint256,uint256))
func (_IStabilization *IStabilization) Cdps(opts *tests.RunOptions, owner common.Address) (IStabilizationCDP, uint64, error) {
	data, consumed, err := _IStabilization.Call(opts, "cdps", owner)

	if err != nil {
		return *new(IStabilizationCDP), consumed, _IStabilization.DecodeError(data, err)
	}
	out, err := _IStabilization.Contract.Abi().Unpack("cdps", data)
	if err != nil {
		return *new(IStabilizationCDP), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(IStabilizationCDP)).(*IStabilizationCDP)
	return out0, consumed, err

}

// CollateralPrice is a free data retrieval call binding the contract method 0x5891de72.
//
// Solidity: function collateralPrice() view returns(uint256 price)
func (_IStabilization *IStabilization) CollateralPrice(opts *tests.RunOptions) (*big.Int, uint64, error) {
	data, consumed, err := _IStabilization.Call(opts, "collateralPrice")

	if err != nil {
		return *new(*big.Int), consumed, _IStabilization.DecodeError(data, err)
	}
	out, err := _IStabilization.Contract.Abi().Unpack("collateralPrice", data)
	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// Config is a free data retrieval call binding the contract method 0x79502c55.
//
// Solidity: function config() view returns((uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256))
func (_IStabilization *IStabilization) Config(opts *tests.RunOptions) (IStabilizationConfig, uint64, error) {
	data, consumed, err := _IStabilization.Call(opts, "config")

	if err != nil {
		return *new(IStabilizationConfig), consumed, _IStabilization.DecodeError(data, err)
	}
	out, err := _IStabilization.Contract.Abi().Unpack("config", data)
	if err != nil {
		return *new(IStabilizationConfig), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(IStabilizationConfig)).(*IStabilizationConfig)
	return out0, consumed, err

}

// DebtAmountAtTime is a free data retrieval call binding the contract method 0x8ffba9dd.
//
// Solidity: function debtAmountAtTime(address account, uint256 timestamp) view returns(uint256)
func (_IStabilization *IStabilization) DebtAmountAtTime(opts *tests.RunOptions, account common.Address, timestamp *big.Int) (*big.Int, uint64, error) {
	data, consumed, err := _IStabilization.Call(opts, "debtAmountAtTime", account, timestamp)

	if err != nil {
		return *new(*big.Int), consumed, _IStabilization.DecodeError(data, err)
	}
	out, err := _IStabilization.Contract.Abi().Unpack("debtAmountAtTime", data)
	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// LastUpdated is a free data retrieval call binding the contract method 0xd0b06f5d.
//
// Solidity: function lastUpdated() view returns((uint256,uint256,uint256,uint256))
func (_IStabilization *IStabilization) LastUpdated(opts *tests.RunOptions) (IStabilizationLastUpdated, uint64, error) {
	data, consumed, err := _IStabilization.Call(opts, "lastUpdated")

	if err != nil {
		return *new(IStabilizationLastUpdated), consumed, _IStabilization.DecodeError(data, err)
	}
	out, err := _IStabilization.Contract.Abi().Unpack("lastUpdated", data)
	if err != nil {
		return *new(IStabilizationLastUpdated), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(IStabilizationLastUpdated)).(*IStabilizationLastUpdated)
	return out0, consumed, err

}

// Liquidate is a free data retrieval call for a paid mutator transaction binding the contract method 0x4914c008.
// Similar to eth_call from rpc calls or function.call from truffle, it reverts the state after the call and returns the output. The output is extracted
// the same way as done above for view only functions.
// Solidity: function liquidate(address account, uint256 collateralSold, address bidder) payable returns()
func (_IStabilization *IStabilization) CallLiquidate(r *tests.Runner, opts *tests.RunOptions, account common.Address, collateralSold *big.Int, bidder common.Address) (uint64, error) {
	snap := r.Snapshot()

	data, consumed, err := _IStabilization.Call(opts, "liquidate", account, collateralSold, bidder)
	r.RevertSnapshot(snap)
	return consumed, _IStabilization.DecodeError(data, err)

}

// SetACU is a free data retrieval call for a paid mutator transaction binding the contract method 0x4b8ef943.
// Similar to eth_call from rpc calls or function.call from truffle, it reverts the state after the call and returns the output. The output is extracted
// the same way as done above for view only functions.
// Solidity: function setACU(address acu) returns()
func (_IStabilization *IStabilization) CallSetACU(r *tests.Runner, opts *tests.RunOptions, acu common.Address) (uint64, error) {
	snap := r.Snapshot()

	data, consumed, err := _IStabilization.Call(opts, "setACU", acu)
	r.RevertSnapshot(snap)
	return consumed, _IStabilization.DecodeError(data, err)

}

// SetAuctioneer is a free data retrieval call for a paid mutator transaction binding the contract method 0x00ede7e4.
// Similar to eth_call from rpc calls or function.call from truffle, it reverts the state after the call and returns the output. The output is extracted
// the same way as done above for view only functions.
// Solidity: function setAuctioneer(address auctioneer) returns()
func (_IStabilization *IStabilization) CallSetAuctioneer(r *tests.Runner, opts *tests.RunOptions, auctioneer common.Address) (uint64, error) {
	snap := r.Snapshot()

	data, consumed, err := _IStabilization.Call(opts, "setAuctioneer", auctioneer)
	r.RevertSnapshot(snap)
	return consumed, _IStabilization.DecodeError(data, err)

}

// SetOperator is a free data retrieval call for a paid mutator transaction binding the contract method 0xb3ab15fb.
// Similar to eth_call from rpc calls or function.call from truffle, it reverts the state after the call and returns the output. The output is extracted
// the same way as done above for view only functions.
// Solidity: function setOperator(address operator) returns()
func (_IStabilization *IStabilization) CallSetOperator(r *tests.Runner, opts *tests.RunOptions, operator common.Address) (uint64, error) {
	snap := r.Snapshot()

	data, consumed, err := _IStabilization.Call(opts, "setOperator", operator)
	r.RevertSnapshot(snap)
	return consumed, _IStabilization.DecodeError(data, err)

}

// SetOracle is a free data retrieval call for a paid mutator transaction binding the contract method 0x7adbf973.
// Similar to eth_call from rpc calls or function.call from truffle, it reverts the state after the call and returns the output. The output is extracted
// the same way as done above for view only functions.
// Solidity: function setOracle(address oracle) returns()
func (_IStabilization *IStabilization) CallSetOracle(r *tests.Runner, opts *tests.RunOptions, oracle common.Address) (uint64, error) {
	snap := r.Snapshot()

	data, consumed, err := _IStabilization.Call(opts, "setOracle", oracle)
	r.RevertSnapshot(snap)
	return consumed, _IStabilization.DecodeError(data, err)

}

// SetSupplyControl is a free data retrieval call for a paid mutator transaction binding the contract method 0x52e5a050.
// Similar to eth_call from rpc calls or function.call from truffle, it reverts the state after the call and returns the output. The output is extracted
// the same way as done above for view only functions.
// Solidity: function setSupplyControl(address supplyControl) returns()
func (_IStabilization *IStabilization) CallSetSupplyControl(r *tests.Runner, opts *tests.RunOptions, supplyControl common.Address) (uint64, error) {
	snap := r.Snapshot()

	data, consumed, err := _IStabilization.Call(opts, "setSupplyControl", supplyControl)
	r.RevertSnapshot(snap)
	return consumed, _IStabilization.DecodeError(data, err)

}

// Liquidate is a paid mutator transaction binding the contract method 0x4914c008.
//
// Solidity: function liquidate(address account, uint256 collateralSold, address bidder) payable returns()
func (_IStabilization *IStabilization) Liquidate(opts *tests.RunOptions, account common.Address, collateralSold *big.Int, bidder common.Address) (uint64, error) {
	data, consumed, err := _IStabilization.Call(opts, "liquidate", account, collateralSold, bidder)
	return consumed, _IStabilization.DecodeError(data, err)
}

// SetACU is a paid mutator transaction binding the contract method 0x4b8ef943.
//
// Solidity: function setACU(address acu) returns()
func (_IStabilization *IStabilization) SetACU(opts *tests.RunOptions, acu common.Address) (uint64, error) {
	data, consumed, err := _IStabilization.Call(opts, "setACU", acu)
	return consumed, _IStabilization.DecodeError(data, err)
}

// SetAuctioneer is a paid mutator transaction binding the contract method 0x00ede7e4.
//
// Solidity: function setAuctioneer(address auctioneer) returns()
func (_IStabilization *IStabilization) SetAuctioneer(opts *tests.RunOptions, auctioneer common.Address) (uint64, error) {
	data, consumed, err := _IStabilization.Call(opts, "setAuctioneer", auctioneer)
	return consumed, _IStabilization.DecodeError(data, err)
}

// SetOperator is a paid mutator transaction binding the contract method 0xb3ab15fb.
//
// Solidity: function setOperator(address operator) returns()
func (_IStabilization *IStabilization) SetOperator(opts *tests.RunOptions, operator common.Address) (uint64, error) {
	data, consumed, err := _IStabilization.Call(opts, "setOperator", operator)
	return consumed, _IStabilization.DecodeError(data, err)
}

// SetOracle is a paid mutator transaction binding the contract method 0x7adbf973.
//
// Solidity: function setOracle(address oracle) returns()
func (_IStabilization *IStabilization) SetOracle(opts *tests.RunOptions, oracle common.Address) (uint64, error) {
	data, consumed, err := _IStabilization.Call(opts, "setOracle", oracle)
	return consumed, _IStabilization.DecodeError(data, err)
}

// SetSupplyControl is a paid mutator transaction binding the contract method 0x52e5a050.
//
// Solidity: function setSupplyControl(address supplyControl) returns()
func (_IStabilization *IStabilization) SetSupplyControl(opts *tests.RunOptions, supplyControl common.Address) (uint64, error) {
	data, consumed, err := _IStabilization.Call(opts, "setSupplyControl", supplyControl)
	return consumed, _IStabilization.DecodeError(data, err)
}

func (_IStabilization *IStabilization) DecodeError(data []byte, err error) error {
	if err == nil {
		return nil
	}

	reason, _ := abi.UnpackRevert(data)
	return fmt.Errorf("%w: %s", err, reason)

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
	*tests.Contract
}

// AvailableSupply is a free data retrieval call binding the contract method 0x7ecc2b56.
//
// Solidity: function availableSupply() view returns(uint256)
func (_ISupplyControl *ISupplyControl) AvailableSupply(opts *tests.RunOptions) (*big.Int, uint64, error) {
	data, consumed, err := _ISupplyControl.Call(opts, "availableSupply")

	if err != nil {
		return *new(*big.Int), consumed, _ISupplyControl.DecodeError(data, err)
	}
	out, err := _ISupplyControl.Contract.Abi().Unpack("availableSupply", data)
	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// GetStabilizer is a free data retrieval call binding the contract method 0x80af1799.
//
// Solidity: function getStabilizer() view returns(address)
func (_ISupplyControl *ISupplyControl) GetStabilizer(opts *tests.RunOptions) (common.Address, uint64, error) {
	data, consumed, err := _ISupplyControl.Call(opts, "getStabilizer")

	if err != nil {
		return *new(common.Address), consumed, _ISupplyControl.DecodeError(data, err)
	}
	out, err := _ISupplyControl.Contract.Abi().Unpack("getStabilizer", data)
	if err != nil {
		return *new(common.Address), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, consumed, err

}

// GetTotalSupply is a free data retrieval call binding the contract method 0xc4e41b22.
//
// Solidity: function getTotalSupply() view returns(uint256)
func (_ISupplyControl *ISupplyControl) GetTotalSupply(opts *tests.RunOptions) (*big.Int, uint64, error) {
	data, consumed, err := _ISupplyControl.Call(opts, "getTotalSupply")

	if err != nil {
		return *new(*big.Int), consumed, _ISupplyControl.DecodeError(data, err)
	}
	out, err := _ISupplyControl.Contract.Abi().Unpack("getTotalSupply", data)
	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// Burn is a free data retrieval call for a paid mutator transaction binding the contract method 0x44df8e70.
// Similar to eth_call from rpc calls or function.call from truffle, it reverts the state after the call and returns the output. The output is extracted
// the same way as done above for view only functions.
// Solidity: function burn() payable returns()
func (_ISupplyControl *ISupplyControl) CallBurn(r *tests.Runner, opts *tests.RunOptions) (uint64, error) {
	snap := r.Snapshot()

	data, consumed, err := _ISupplyControl.Call(opts, "burn")
	r.RevertSnapshot(snap)
	return consumed, _ISupplyControl.DecodeError(data, err)

}

// Mint is a free data retrieval call for a paid mutator transaction binding the contract method 0x40c10f19.
// Similar to eth_call from rpc calls or function.call from truffle, it reverts the state after the call and returns the output. The output is extracted
// the same way as done above for view only functions.
// Solidity: function mint(address recipient, uint256 amount) returns()
func (_ISupplyControl *ISupplyControl) CallMint(r *tests.Runner, opts *tests.RunOptions, recipient common.Address, amount *big.Int) (uint64, error) {
	snap := r.Snapshot()

	data, consumed, err := _ISupplyControl.Call(opts, "mint", recipient, amount)
	r.RevertSnapshot(snap)
	return consumed, _ISupplyControl.DecodeError(data, err)

}

// SetOperator is a free data retrieval call for a paid mutator transaction binding the contract method 0xb3ab15fb.
// Similar to eth_call from rpc calls or function.call from truffle, it reverts the state after the call and returns the output. The output is extracted
// the same way as done above for view only functions.
// Solidity: function setOperator(address operator) returns()
func (_ISupplyControl *ISupplyControl) CallSetOperator(r *tests.Runner, opts *tests.RunOptions, operator common.Address) (uint64, error) {
	snap := r.Snapshot()

	data, consumed, err := _ISupplyControl.Call(opts, "setOperator", operator)
	r.RevertSnapshot(snap)
	return consumed, _ISupplyControl.DecodeError(data, err)

}

// SetStabilizer is a free data retrieval call for a paid mutator transaction binding the contract method 0xdb7f521a.
// Similar to eth_call from rpc calls or function.call from truffle, it reverts the state after the call and returns the output. The output is extracted
// the same way as done above for view only functions.
// Solidity: function setStabilizer(address stabilizer_) returns()
func (_ISupplyControl *ISupplyControl) CallSetStabilizer(r *tests.Runner, opts *tests.RunOptions, stabilizer_ common.Address) (uint64, error) {
	snap := r.Snapshot()

	data, consumed, err := _ISupplyControl.Call(opts, "setStabilizer", stabilizer_)
	r.RevertSnapshot(snap)
	return consumed, _ISupplyControl.DecodeError(data, err)

}

// Burn is a paid mutator transaction binding the contract method 0x44df8e70.
//
// Solidity: function burn() payable returns()
func (_ISupplyControl *ISupplyControl) Burn(opts *tests.RunOptions) (uint64, error) {
	data, consumed, err := _ISupplyControl.Call(opts, "burn")
	return consumed, _ISupplyControl.DecodeError(data, err)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address recipient, uint256 amount) returns()
func (_ISupplyControl *ISupplyControl) Mint(opts *tests.RunOptions, recipient common.Address, amount *big.Int) (uint64, error) {
	data, consumed, err := _ISupplyControl.Call(opts, "mint", recipient, amount)
	return consumed, _ISupplyControl.DecodeError(data, err)
}

// SetOperator is a paid mutator transaction binding the contract method 0xb3ab15fb.
//
// Solidity: function setOperator(address operator) returns()
func (_ISupplyControl *ISupplyControl) SetOperator(opts *tests.RunOptions, operator common.Address) (uint64, error) {
	data, consumed, err := _ISupplyControl.Call(opts, "setOperator", operator)
	return consumed, _ISupplyControl.DecodeError(data, err)
}

// SetStabilizer is a paid mutator transaction binding the contract method 0xdb7f521a.
//
// Solidity: function setStabilizer(address stabilizer_) returns()
func (_ISupplyControl *ISupplyControl) SetStabilizer(opts *tests.RunOptions, stabilizer_ common.Address) (uint64, error) {
	data, consumed, err := _ISupplyControl.Call(opts, "setStabilizer", stabilizer_)
	return consumed, _ISupplyControl.DecodeError(data, err)
}

func (_ISupplyControl *ISupplyControl) DecodeError(data []byte, err error) error {
	if err == nil {
		return nil
	}

	reason, _ := abi.UnpackRevert(data)
	return fmt.Errorf("%w: %s", err, reason)

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
	*tests.Contract
}

// SetOperator is a free data retrieval call for a paid mutator transaction binding the contract method 0xb3ab15fb.
// Similar to eth_call from rpc calls or function.call from truffle, it reverts the state after the call and returns the output. The output is extracted
// the same way as done above for view only functions.
// Solidity: function setOperator(address _account) returns()
func (_IUpgradeManager *IUpgradeManager) CallSetOperator(r *tests.Runner, opts *tests.RunOptions, _account common.Address) (uint64, error) {
	snap := r.Snapshot()

	data, consumed, err := _IUpgradeManager.Call(opts, "setOperator", _account)
	r.RevertSnapshot(snap)
	return consumed, _IUpgradeManager.DecodeError(data, err)

}

// SetOperator is a paid mutator transaction binding the contract method 0xb3ab15fb.
//
// Solidity: function setOperator(address _account) returns()
func (_IUpgradeManager *IUpgradeManager) SetOperator(opts *tests.RunOptions, _account common.Address) (uint64, error) {
	data, consumed, err := _IUpgradeManager.Call(opts, "setOperator", _account)
	return consumed, _IUpgradeManager.DecodeError(data, err)
}

func (_IUpgradeManager *IUpgradeManager) DecodeError(data []byte, err error) error {
	if err == nil {
		return nil
	}

	reason, _ := abi.UnpackRevert(data)
	return fmt.Errorf("%w: %s", err, reason)

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
	Bin: "0x6080604052600160ff1b600a55348015601757600080fd5b50615351806100276000396000f3fe6080604052600436106101a25760003560e01c806399b0014b116100e0578063da39fbfe11610084578063ed78349d11610061578063ed78349d14610509578063ef5cc4d11461051e578063f0141d841461053e578063fb09917e1461055a57005b8063da39fbfe146104a7578063da78110e146104c7578063df7f710e146104e757005b8063b3ab15fb116100bd578063b3ab15fb1461043b578063b78dec521461045b578063c3f909d414610470578063cdd722531461049257005b806399b0014b146103a15780639ed1f255146103b65780639f8743f71461042657005b80634bb278f31161014757806359974e381161012457806359974e38146103395780635a4d3a271461034c57806367b11630146103615780638d4f75d21461038157005b80634bb278f3146102df57806356833ebe1461030457806357eba7591461032457005b806333d162931161018057806333d162931461021c57806333f98c771461024a5780633c8510fd1461029f5780633f422ef3146102bf57005b8063077945d3146101a45780630f65875c146101cf5780632d35d158146101e4575b005b3480156101b057600080fd5b506101b96105a9565b6040516101c691906144d2565b60405180910390f35b3480156101db57600080fd5b506101a261066a565b3480156101f057600080fd5b506102046101ff366004614535565b610888565b6040516001600160a01b0390911681526020016101c6565b34801561022857600080fd5b5061023c610237366004614535565b6108fd565b6040519081526020016101c6565b34801561025657600080fd5b5061026a610265366004614620565b61096c565b6040516101c6919081518152602080830151908201526040808301519082015260609182015115159181019190915260800190565b3480156102ab57600080fd5b5061026a6102ba366004614655565b610aa8565b3480156102cb57600080fd5b506101a26102da36600461469c565b610bbb565b3480156102eb57600080fd5b506102f4610d6b565b60405190151581526020016101c6565b34801561031057600080fd5b506101a261031f3660046146d8565b610f0a565b34801561033057600080fd5b5061023c611454565b6101a2610347366004614774565b6114ae565b34801561035857600080fd5b5061023c61154b565b34801561036d57600080fd5b506101a261037c366004614774565b6115a5565b34801561038d57600080fd5b506101a261039c3660046147b1565b61169a565b3480156103ad57600080fd5b50600a5461023c565b3480156103c257600080fd5b506103d66103d1366004614535565b61180d565b6040516101c69190600060c0820190508251825260208301516020830152604083015160408301526060830151606083015260808301511515608083015260a0830151151560a083015292915050565b34801561043257600080fd5b5061023c611909565b34801561044757600080fd5b506101a2610456366004614535565b611963565b34801561046757600080fd5b5061023c611a1d565b34801561047c57600080fd5b50610485611a77565b6040516101c6919061486d565b34801561049e57600080fd5b506101b9611b8d565b3480156104b357600080fd5b506101a26104c23660046148e9565b611c46565b3480156104d357600080fd5b506101a26104e236600461498a565b611e97565b3480156104f357600080fd5b506104fc6120cd565b6040516101c69190614aca565b34801561051557600080fd5b5060075461023c565b34801561052a57600080fd5b50610204610539366004614535565b61228b565b34801561054a57600080fd5b50604051601281526020016101c6565b34801561056657600080fd5b5061057a610575366004614add565b6122fd565b6040805182516effffffffffffffffffffffffffffff16815260209283015160ff1692810192909252016101c6565b60606105b760005460011490565b156106095760405162461bcd60e51b815260206004820152601d60248201527f72656164206f6e6c79207265656e7472616e637920646574656374656400000060448201526064015b60405180910390fd5b601280548060200260200160405190810160405280929190818152602001828054801561065f57602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311610641575b505050505090505b90565b6001546001600160a01b031633146106ea5760405162461bcd60e51b815260206004820152602360248201527f7265737472696374656420746f20746865206175746f6e69747920636f6e747260448201527f61637400000000000000000000000000000000000000000000000000000000006064820152608401610600565b60135460ff1615156001036107905760005b60125481101561075e576001600c60006012848154811061071f5761071f614b2b565b6000918252602080832091909101546001600160a01b031683528201929092526040019020600401805460ff19169115159190911790556001016106fc565b50601380547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000166101001790556107d5565b601354610100900460ff1615156001036107d5576107ac6123d7565b601380547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff1690555b600a546107e3906002614b89565b6014540361088657600f80546107fb91600e9161430f565b5060005b601154811015610884576000600c60006011848154811061082257610822614b2b565b60009182526020808320909101546001600160a01b0316835282019290925260400190206004018054911515610100027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff9092169190911790556001016107ff565b505b565b600080546001036108db5760405162461bcd60e51b815260206004820152601d60248201527f72656164206f6e6c79207265656e7472616e63792064657465637465640000006044820152606401610600565b506001600160a01b03808216600090815260196020526040902054165b919050565b600080546001036109505760405162461bcd60e51b815260206004820152601d60248201527f72656164206f6e6c79207265656e7472616e63792064657465637465640000006044820152606401610600565b506001600160a01b03166000908152601a602052604090205490565b61099960405180608001604052806000815260200160008152602001600081526020016000151581525090565b6000546001036109eb5760405162461bcd60e51b815260206004820152601d60248201527f72656164206f6e6c79207265656e7472616e63792064657465637465640000006044820152606401610600565b6000601560016014546109fe9190614bb1565b81548110610a0e57610a0e614b2b565b9060005260206000200183604051610a269190614bc4565b90815260408051918290036020908101832060608401835280548452600180820154928501929092526002015460ff161515838301528151608081019092526014549293506000928291610a7991614bb1565b815260200183600001518152602001836020015181526020018360400151151581525090508092505050919050565b610ad560405180608001604052806000815260200160008152602001600081526020016000151581525090565b600054600103610b275760405162461bcd60e51b815260206004820152601d60248201527f72656164206f6e6c79207265656e7472616e63792064657465637465640000006044820152606401610600565b600060158481548110610b3c57610b3c614b2b565b9060005260206000200183604051610b549190614bc4565b9081526040805191829003602090810183206060808501845281548552600182015485840190815260029092015460ff16151585850190815284516080810186528a8152955193860193909352905192840192909252511515908201529150505b92915050565b6002546001600160a01b03163314610c155760405162461bcd60e51b815260206004820152601660248201527f7265737472696374656420746f206f70657261746f72000000000000000000006044820152606401610600565b8082108015610c245750600081115b610c705760405162461bcd60e51b815260206004820152600e60248201527f696e76616c696420636f6e6669670000000000000000000000000000000000006044820152606401610600565b6008546040805160808082526013908201527f72657665616c5265736574496e74657276616c0000000000000000000000000060a0820152602081019290925281018290524360608201527f207e45ce6f2191c3efffe5d0d91591cf81afd23cbf4bd746981d7a8a8dcfe1ba9060c00160405180910390a160088190556007546040805160808082526012908201527f6e6f6e52657665616c5468726573686f6c64000000000000000000000000000060a0820152602081019290925281018390524360608201527f207e45ce6f2191c3efffe5d0d91591cf81afd23cbf4bd746981d7a8a8dcfe1ba9060c00160405180910390a150600755565b6001546000906001600160a01b03163314610dee5760405162461bcd60e51b815260206004820152602360248201527f7265737472696374656420746f20746865206175746f6e69747920636f6e747260448201527f61637400000000000000000000000000000000000000000000000000000000006064820152608401610600565b610df6612607565b600354600b54610e069190614be0565b431015610e1557506000610f01565b6000610e1f61265e565b600854601454919250610e3191614c22565b600003610e4057610e4061285d565b60158054600101815560009081525b600e54811015610e7557610e6381836128b3565b610e6e600182614be0565b9050610e4f565b50610e7e613073565b43600b81905550600160146000828254610e989190614be0565b909155505060105460035414610eaf576010546003555b601454600354604080519283524260208401528201527f5aec57d81928b24d30b1a2aec0d23d693412c37d7ec106b5d8259413716bb1f49060600160405180910390a1610efb81613167565b60019150505b61066760008055565b336000908152600c602052604090206004015460ff16610f6c5760405162461bcd60e51b815260206004820152601960248201527f7265737472696374656420746f206f6e6c7920766f74657273000000000000006044820152606401610600565b610f74612607565b336000908152600c60205260409020601454815403610fd55760405162461bcd60e51b815260206004820152600d60248201527f616c726561647920766f746564000000000000000000000000000000000000006044820152606401610600565b60018101805490879055815460145483556000819003611032576040805133815260ff861660208201527fd2ec8e890a03083998d3e16f98044fd3dd13fe3e61b7bc2e58ee6da43b50af73910160405180910390a1505050611444565b60016014546110419190614bb1565b81146110e057336001600160a01b03167f04ca4e0efda95f8b780c116574d1521309010b38d8f7b75705495703a0f570b160016014546110819190614bb1565b6040805160808082526016908201527f4c617374566f746564526f756e644d69736d617463680000000000000000000060a08201526020810192909252810184905260ff8716606082015260c0015b60405180910390a2505050611444565b868686336040516020016110f79493929190614c53565b6040516020818303038152906040528051906020012060001c975087821461119a576111233384613213565b604080516080808252600e908201527f436f6d6d69744d69736d6174636800000000000000000000000000000000000060a08201526020810184905290810189905260ff8516606082015233907f04ca4e0efda95f8b780c116574d1521309010b38d8f7b75705495703a0f570b19060c0016110d0565b600e54861461121e57600e546040805160808082526014908201527f5265706f72744c656e6774684d69736d6174636800000000000000000000000060a0820152602081018990529081019190915260ff8516606082015233907f04ca4e0efda95f8b780c116574d1521309010b38d8f7b75705495703a0f570b19060c0016110d0565b60005b868110156113d957606488888381811061123d5761123d614b2b565b90506040020160200160208101906112559190614cd7565b60ff1611156112a65760405162461bcd60e51b815260206004820152601860248201527f696e76616c696420636f6e666964656e63652073636f726500000000000000006044820152606401610600565b60008888838181106112ba576112ba614b2b565b6112d09260206040909202019081019150614cf4565b6effffffffffffffffffffffffffffff1611801561131a575060008888838181106112fd576112fd614b2b565b90506040020160200160208101906113159190614cd7565b60ff16115b6113665760405162461bcd60e51b815260206004820152601660248201527f636f6e666964656e63652f7072696365206572726f72000000000000000000006044820152606401610600565b87878281811061137857611378614b2b565b905060400201600d600e838154811061139357611393614b2b565b906000526020600020016040516113aa9190614d5e565b9081526040805160209281900383019020336000908152925290206113cf8282614dd3565b5050600101611221565b506004830180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff1661010017905560405160ff8516815233907f8bdddd7f2f2c74679ffa6beb8f86aa18bfa5baf1bfaf534d0b66596babc53f089060200160405180910390a25050505b61144d60008055565b5050505050565b600080546001036114a75760405162461bcd60e51b815260206004820152601d60248201527f72656164206f6e6c79207265656e7472616e63792064657465637465640000006044820152606401610600565b5060105490565b6001546001600160a01b0316331461152e5760405162461bcd60e51b815260206004820152602360248201527f7265737472696374656420746f20746865206175746f6e69747920636f6e747260448201527f61637400000000000000000000000000000000000000000000000000000000006064820152608401610600565b611536612607565b476115418183613280565b5061088460008055565b6000805460010361159e5760405162461bcd60e51b815260206004820152601d60248201527f72656164206f6e6c79207265656e7472616e63792064657465637465640000006044820152606401610600565b50600b5490565b6002546001600160a01b031633146115ff5760405162461bcd60e51b815260206004820152601660248201527f7265737472696374656420746f206f70657261746f72000000000000000000006044820152606401610600565b6116088161352a565b6010819055600354600b547f207e45ce6f2191c3efffe5d0d91591cf81afd23cbf4bd746981d7a8a8dcfe1ba91908390611643908390614be0565b604080516080808252600a908201527f766f7465506572696f640000000000000000000000000000000000000000000060a08201526020810194909452830191909152606082015260c0015b60405180910390a150565b6002546001600160a01b031633146116f45760405162461bcd60e51b815260206004820152601660248201527f7265737472696374656420746f206f70657261746f72000000000000000000006044820152606401610600565b80516000036117455760405162461bcd60e51b815260206004820152601660248201527f73796d626f6c732063616e277420626520656d707479000000000000000000006044820152606401610600565b601454600a54611756906001614b89565b141580156117685750601454600a5414155b6117b45760405162461bcd60e51b815260206004820152601e60248201527f63616e2774206265207570646174656420696e207468697320726f756e6400006044820152606401610600565b80516117c790600f906020840190614367565b50601454600a8190557faa278e424da680ce5dad66510415760e78e0bd87d45c786c6e88bdde82f9342d9082906117ff906001614be0565b60405161168f929190614e6e565b61184a6040518060c00160405280600081526020016000815260200160008152602001600081526020016000151581526020016000151581525090565b60005460010361189c5760405162461bcd60e51b815260206004820152601d60248201527f72656164206f6e6c79207265656e7472616e63792064657465637465640000006044820152606401610600565b506001600160a01b03166000908152600c6020908152604091829020825160c081018452815481526001820154928101929092526002810154928201929092526003820154606082015260049091015460ff8082161515608084015261010090910416151560a082015290565b6000805460010361195c5760405162461bcd60e51b815260206004820152601d60248201527f72656164206f6e6c79207265656e7472616e63792064657465637465640000006044820152606401610600565b5060145490565b6001546001600160a01b031633146119e35760405162461bcd60e51b815260206004820152602360248201527f7265737472696374656420746f20746865206175746f6e69747920636f6e747260448201527f61637400000000000000000000000000000000000000000000000000000000006064820152608401610600565b600280547fffffffffffffffffffffffff0000000000000000000000000000000000000000166001600160a01b0392909216919091179055565b60008054600103611a705760405162461bcd60e51b815260206004820152601d60248201527f72656164206f6e6c79207265656e7472616e63792064657465637465640000006044820152606401610600565b5060035490565b611ad860405180610120016040528060006001600160a01b0316815260200160006001600160a01b03168152602001600081526020016000815260200160008152602001600081526020016000815260200160008152602001600081525090565b600054600103611b2a5760405162461bcd60e51b815260206004820152601d60248201527f72656164206f6e6c79207265656e7472616e63792064657465637465640000006044820152606401610600565b5060408051610120810182526001546001600160a01b039081168252600254166020820152600354918101919091526004546060820152600554608082015260065460a082015260075460c082015260085460e082015260095461010082015290565b6060611b9b60005460011490565b15611be85760405162461bcd60e51b815260206004820152601d60248201527f72656164206f6e6c79207265656e7472616e63792064657465637465640000006044820152606401610600565b601180548060200260200160405190810160405280929190818152602001828054801561065f576020028201919060005260206000209081546001600160a01b03168152600190910190602001808311610641575050505050905090565b6002546001600160a01b03163314611ca05760405162461bcd60e51b815260206004820152601660248201527f7265737472696374656420746f206f70657261746f72000000000000000000006044820152606401610600565b6005546040805160808082526018908201527f6f75746c696572536c617368696e675468726573686f6c64000000000000000060a0820152602081019290925281018590524360608201527fb5114472b89d1126433287defb6308bbabdb95b6ce5dd949ab2c151ed7b1ff4c9060c00160405180910390a160058490556004546040805160808082526019908201527f6f75746c696572446574656374696f6e5468726573686f6c640000000000000060a0820152602081019290925281018490524360608201527fb5114472b89d1126433287defb6308bbabdb95b6ce5dd949ab2c151ed7b1ff4c9060c00160405180910390a160048390556006546040805160808082526010908201527f62617365536c617368696e67526174650000000000000000000000000000000060a0820152602081019290925281018390524360608201527f207e45ce6f2191c3efffe5d0d91591cf81afd23cbf4bd746981d7a8a8dcfe1ba9060c00160405180910390a16006829055600954604080516080808252600f908201527f736c617368696e6752617465436170000000000000000000000000000000000060a0820152602081019290925281018290524360608201527f207e45ce6f2191c3efffe5d0d91591cf81afd23cbf4bd746981d7a8a8dcfe1ba9060c00160405180910390a1600955505050565b6001546001600160a01b03163314611f175760405162461bcd60e51b815260206004820152602360248201527f7265737472696374656420746f20746865206175746f6e69747920636f6e747260448201527f61637400000000000000000000000000000000000000000000000000000000006064820152608401610600565b611f1f612607565b8251600003611f705760405162461bcd60e51b815260206004820152601560248201527f566f746572732063616e277420626520656d70747900000000000000000000006044820152606401610600565b60005b835181101561208457828181518110611f8e57611f8e614b2b565b602002602001015160186000868481518110611fac57611fac614b2b565b60200260200101516001600160a01b03166001600160a01b0316815260200190815260200160002060006101000a8154816001600160a01b0302191690836001600160a01b0316021790555081818151811061200a5761200a614b2b565b60200260200101516019600086848151811061202857612028614b2b565b6020908102919091018101516001600160a01b0390811683529082019290925260400160002080547fffffffffffffffffffffffff00000000000000000000000000000000000000001692909116919091179055600101611f73565b5061209e836000600186516120999190614bb1565b6136f0565b82516120b19060129060208601906143ad565b506013805460ff191660011790556120c860008055565b505050565b6060601454600a5460016120e19190614b89565b036121bd57600f805480602002602001604051908101604052809291908181526020016000905b828210156121b457838290600052602060002001805461212790614d11565b80601f016020809104026020016040519081016040528092919081815260200182805461215390614d11565b80156121a05780601f10612175576101008083540402835291602001916121a0565b820191906000526020600020905b81548152906001019060200180831161218357829003601f168201915b505050505081526020019060010190612108565b50505050905090565b600e805480602002602001604051908101604052809291908181526020016000905b828210156121b45783829060005260206000200180546121fe90614d11565b80601f016020809104026020016040519081016040528092919081815260200182805461222a90614d11565b80156122775780601f1061224c57610100808354040283529160200191612277565b820191906000526020600020905b81548152906001019060200180831161225a57829003601f168201915b5050505050815260200190600101906121df565b600080546001036122de5760405162461bcd60e51b815260206004820152601d60248201527f72656164206f6e6c79207265656e7472616e63792064657465637465640000006044820152606401610600565b506001600160a01b039081166000908152601860205260409020541690565b60408051808201909152600080825260208201526000546001036123635760405162461bcd60e51b815260206004820152601d60248201527f72656164206f6e6c79207265656e7472616e63792064657465637465640000006044820152606401610600565b600d836040516123739190614bc4565b90815260408051602092819003830181206001600160a01b0395909516600090815294835293819020848201909152546effffffffffffffffffffffffffffff811684526f01000000000000000000000000000000900460ff169083015250919050565b6000805b601154821080156123ed575060125481105b1561255e576012818154811061240557612405614b2b565b600091825260209091200154601180546001600160a01b03909216918490811061243157612431614b2b565b6000918252602090912001546001600160a01b03160361246b578161245581614e90565b925050808061246390614e90565b9150506123db565b6012818154811061247e5761247e614b2b565b600091825260209091200154601180546001600160a01b0390921691849081106124aa576124aa614b2b565b6000918252602090912001546001600160a01b0316101561255457600c6000601184815481106124dc576124dc614b2b565b60009182526020808320909101546001600160a01b0316835282019290925260400181208181556001810182905560028101829055600381019190915560040180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00001690558161254c81614e90565b9250506123db565b8061246381614e90565b6011548210156125f757600c60006011848154811061257f5761257f614b2b565b60009182526020808320909101546001600160a01b0316835282019290925260400181208181556001810182905560028101829055600381019190915560040180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000169055816125ef81614e90565b92505061255e565b601280546120c891601191614426565b600054156126575760405162461bcd60e51b815260206004820152601360248201527f7265656e7472616e6379206465746563746564000000000000000000000000006044820152606401610600565b6001600055565b60115460609060009067ffffffffffffffff81111561267f5761267f614550565b6040519080825280602002602001820160405280156126a8578160200160208202803683370190505b50905060005b601154811015612857576000601182815481106126cd576126cd614b2b565b60009182526020808320909101546001600160a01b0316808352600c90915260409091206014549192509061270490600190614bb1565b8154148015612717575060008160010154115b15612726576127268282613213565b6007546003820154111561284d57600184848151811061274857612748614b2b565b911515602092830291909101820152601454600383015460408051928352928201526001600160a01b038416917f9e6b40f10c60d1ad09594f3b6ed7043d0e978f584d354ace6e1f6025660c42b1910160405180910390a26000600382018190556001546001600160a01b03848116835260196020526040928390205460095493517f02fb4d850000000000000000000000000000000000000000000000000000000081529082166004820152602481019390935216906302fb4d85906044016020604051808303816000875af1158015612827573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061284b9190614eaa565b505b50506001016126ae565b50919050565b60005b601154811015610884576000600c60006011848154811061288357612883614b2b565b60009182526020808320909101546001600160a01b03168352820192909252604001902060030155600101612860565b6000600e83815481106128c8576128c8614b2b565b9060005260206000200180546128dd90614d11565b80601f016020809104026020016040519081016040528092919081815260200182805461290990614d11565b80156129565780601f1061292b57610100808354040283529160200191612956565b820191906000526020600020905b81548152906001019060200180831161293957829003601f168201915b50505050509050600060118054905067ffffffffffffffff81111561297d5761297d614550565b6040519080825280602002602001820160405280156129c257816020015b604080518082019091526000808252602082015281526020019060019003908161299b5790505b5090506000805b601154811015612ac6576000601182815481106129e8576129e8614b2b565b60009182526020808320909101546001600160a01b0316808352600c90915260409091206004015490915060ff61010090910416612a265750612abe565b600d85604051612a369190614bc4565b9081526040805191829003602090810183206001600160a01b038516600090815290825282902083830190925290546effffffffffffffffffffffffffffff8116835260ff6f0100000000000000000000000000000090910416908201528484612a9f81614e90565b955081518110612ab157612ab1614b2b565b6020026020010181905250505b6001016129c9565b508015612f41576000612ad9838361389d565b6effffffffffffffffffffffffffffff1690506000612af8828661396d565b606081015190915015612e0e5760005b8160200151811015612d16576000612bf083600001518381518110612b2f57612b2f614b2b565b602002602001015185600d8a604051612b489190614bc4565b90815260200160405180910390206000601188600001518881518110612b7057612b70614b2b565b602002602001015181548110612b8857612b88614b2b565b60009182526020808320909101546001600160a01b0316835282810193909352604091820190208151808301909252546effffffffffffffffffffffffffffff8116825260ff6f0100000000000000000000000000000090910416918101919091528b613cb6565b9050601183600001518381518110612c0a57612c0a614b2b565b602002602001015181548110612c2257612c22614b2b565b9060005260206000200160009054906101000a90046001600160a01b03166001600160a01b03167f372858b237c8bd0714183e8351a461d6c3cb1ef83806181b36bf5943711f4f57828987600d8c604051612c7d9190614bc4565b9081526020016040518091039020600060118a600001518a81518110612ca557612ca5614b2b565b602002602001015181548110612cbd57612cbd614b2b565b60009182526020808320909101546001600160a01b031683528201929092526040908101909120549051612d05949392916effffffffffffffffffffffffffffff1690614ec3565b60405180910390a250600101612b08565b506000612d2b82604001518360600151613ea5565b9050604051806060016040528082815260200142815260200160011515815250601560145481548110612d6057612d60614b2b565b9060005260206000200187604051612d789190614bc4565b9081526040805160209281900383018120845181559284015160018401559201516002909101805460ff1916911515919091179055612db8908790614bc4565b604080519182900382206014548484526020840152600183830152426060840152905190917f5f2aa51aa7889ad71d9318fa7fd83c8ff3277434249bd06073f15986e197911c919081900360800190a250612f3a565b600060156001601454612e219190614bb1565b81548110612e3157612e31614b2b565b9060005260206000200186604051612e499190614bc4565b9081526020016040518091039020600001549050604051806060016040528082815260200142815260200160001515815250601560145481548110612e9057612e90614b2b565b9060005260206000200187604051612ea89190614bc4565b9081526040805160209281900383018120845181559284015160018401559201516002909101805460ff1916911515919091179055612ee8908790614bc4565b604080519182900382206014548484526020840152600083830152426060840152905190917f5f2aa51aa7889ad71d9318fa7fd83c8ff3277434249bd06073f15986e197911c919081900360800190a2505b505061144d565b600060156001601454612f549190614bb1565b81548110612f6457612f64614b2b565b9060005260206000200184604051612f7c9190614bc4565b9081526020016040518091039020600001549050604051806060016040528082815260200142815260200160001515815250601560145481548110612fc357612fc3614b2b565b9060005260206000200185604051612fdb9190614bc4565b9081526040805160209281900383018120845181559284015160018401559201516002909101805460ff191691151591909117905561301b908590614bc4565b604080519182900382206014548484526020840152600083830152426060840152905190917f5f2aa51aa7889ad71d9318fa7fd83c8ff3277434249bd06073f15986e197911c919081900360800190a2505050505050565b60005b6011548110156108845760006011828154811061309557613095614b2b565b60009182526020808320909101546001600160a01b0316808352600c9091526040909120600201549091501561315e576130d0601682613f63565b506001600160a01b0381166000908152600c6020908152604080832060020154601a9092528220805491929091613108908490614be0565b90915550506001600160a01b0381166000908152600c6020526040812060020154601b80549192909161313c908490614be0565b90915550506001600160a01b0381166000908152600c60205260408120600201555b50600101613076565b60005b60115481101561320f5781818151811061318657613186614b2b565b602002602001015115613207576000600c6000601184815481106131ac576131ac614b2b565b60009182526020808320909101546001600160a01b0316835282019290925260400190206004018054911515610100027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff9092169190911790555b60010161316a565b5050565b60038101805490600061322583614e90565b9190505550816001600160a01b03167f176956a4e941f6737f81a3c9a09d8571dd0438d86e25a432beb2013aced430926014548360030154604051613274929190918252602082015260400190565b60405180910390a25050565b601b5460000361328e575050565b600061329a6016613f7f565b905060005b81518110156134e65760008282815181106132bc576132bc614b2b565b602002602001015190506000601b54601a6000846001600160a01b03166001600160a01b0316815260200190815260200160002054876132fc9190614f04565b6133069190614f1b565b601b546001600160a01b0384166000908152601a6020526040812054929350916133309088614f04565b61333a9190614f1b565b90508115613413576001600160a01b03838116600090815260186020526040808220549051919283929116906108fc90869084818181858888f193505050503d80600081146133a5576040519150601f19603f3d011682016040523d82523d6000602084013e6133aa565b606091505b509092509050811515600003613410576001600160a01b03808616600090815260186020526040908190205490517f1137d8c966ce69b9630fb2294be011f3d64cc56e91fad7d375f0662568e9d352926134079216908490614f2f565b60405180910390a15b50505b80156134b2576001546001600160a01b038481166000908152601960205260408082205490517ff7fcc510000000000000000000000000000000000000000000000000000000008152908316600482015260248101859052604481019190915291169063f7fcc51090606401600060405180830381600087803b15801561349957600080fd5b505af11580156134ad573d6000803e3d6000fd5b505050505b6001600160a01b0383166000908152601a60205260408120556134d6601684613f8c565b50506001909201915061329f9050565b5060408051838152602081018590527f3e5aaff9e8fd4293ae18127809c2d4069d87fe10c7de92aa39557a1edbd48fec910160405180910390a150506000601b5550565b600154604080517f0aac2da100000000000000000000000000000000000000000000000000000000815290516000926001600160a01b031691630aac2da19160048083019260209291908290030181865afa15801561358d573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906135b19190614eaa565b9050806135bf836002614f04565b111561360d5760405162461bcd60e51b815260206004820152601660248201527f766f746520706572696f6420697320746f6f20626967000000000000000000006044820152606401610600565b600154604080517fdfb1a4d200000000000000000000000000000000000000000000000000000000815290516001600160a01b039092169163dfb1a4d2916004808201926020929091908290030181865afa158015613670573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906136949190614eaa565b9050806136a2836002614f04565b111561320f5760405162461bcd60e51b815260206004820152601660248201527f766f746520706572696f6420697320746f6f20626967000000000000000000006044820152606401610600565b8082126136fc57505050565b8181600085600261370d8585614f5f565b6137179190614f7f565b6137219087614b89565b8151811061373157613731614b2b565b602002602001015190505b81831361386f575b806001600160a01b031686848151811061376057613760614b2b565b60200260200101516001600160a01b03161015613789578261378181614fc9565b935050613744565b806001600160a01b03168683815181106137a5576137a5614b2b565b60200260200101516001600160a01b031611156137ce57816137c681614ffa565b925050613789565b81831361386a578582815181106137e7576137e7614b2b565b602002602001015186848151811061380157613801614b2b565b602002602001015187858151811061381b5761381b614b2b565b6020026020010188858151811061383457613834614b2b565b6001600160a01b039384166020918202929092010152911690528261385881614fc9565b935050818061386690614ffa565b9250505b61373c565b81851215613882576138828686846136f0565b83831215613895576138958684866136f0565b505050505050565b6000816000036138af57506000610bb5565b6138c58360006138c0600186614bb1565b613fa1565b60006138d2600284614f1b565b90506138df600284614c22565b15613907578381815181106138f6576138f6614b2b565b602002602001015160000151613965565b600284828151811061391b5761391b614b2b565b602002602001015160000151856001846139359190614bb1565b8151811061394557613945614b2b565b60200260200101516000015161395b9190615033565b613965919061505a565b949350505050565b6139986040518060800160405280606081526020016000815260200160608152602001600081525090565b6139c36040518060800160405280606081526020016000815260200160608152602001600081525090565b60115467ffffffffffffffff8111156139de576139de614550565b604051908082528060200260200182016040528015613a2357816020015b60408051808201909152600080825260208201528152602001906001900390816139fc5790505b50604082015260115467ffffffffffffffff811115613a4457613a44614550565b604051908082528060200260200182016040528015613a6d578160200160208202803683370190505b50815260005b601154811015613cae57600060118281548110613a9257613a92614b2b565b60009182526020808320909101546001600160a01b0316808352600c90915260409091206004015490915060ff61010090910416613ad05750613ca6565b600086600d87604051613ae39190614bc4565b90815260408051602092819003830190206001600160a01b03861660009081529252902054613b23906effffffffffffffffffffffffffffff1689614f5f565b613b2e906064615098565b613b389190614f7f565b6004549091508113801590613b5a5750600454613b5782600019615098565b13155b15613c7057600d86604051613b6f9190614bc4565b9081526040805191829003602090810183206001600160a01b03861660009081529082528290208383018352546effffffffffffffffffffffffffffff8116845260ff6f01000000000000000000000000000000909104169083015285015160608601805190613bde82614e90565b905281518110613bf057613bf0614b2b565b6020026020010181905250600d86604051613c0b9190614bc4565b90815260408051602092819003830190206001600160a01b038516600090815290835281812054600c90935290812060020180546f0100000000000000000000000000000090930460ff1692909190613c65908490614be0565b90915550613ca39050565b8351602085018051859291613c8482614e90565b905281518110613c9657613c96614b2b565b6020026020010181815250505b50505b600101613a73565b509392505050565b60008060118681548110613ccc57613ccc614b2b565b60009182526020808320909101546001600160a01b0316808352600c909152604090912060040180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff1690558351909150839087908110613d3057613d30614b2b565b602002602001015115613d47576000915050613965565b6000858686600001516effffffffffffffffffffffffffffff16613d6b9190614f5f565b613d76906064615098565b613d809190614f7f565b9050613d8c8180615098565b6005549091508113613da357600092505050613965565b6000612710600160050154876020015160ff1660016004015485613dc79190614f5f565b613dd19190614f04565b613ddb9190614f04565b613de59190614f1b565b600954909150811115613df757506009545b6001546001600160a01b03848116600090815260196020526040908190205490517f02fb4d850000000000000000000000000000000000000000000000000000000081529082166004820152602481018490529116906302fb4d85906044016020604051808303816000875af1158015613e75573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613e999190614eaa565b98975050505050505050565b60008080805b84811015613f4f57858181518110613ec557613ec5614b2b565b60200260200101516020015160ff16868281518110613ee657613ee6614b2b565b602002602001015160000151613efc91906150e4565b613f17906effffffffffffffffffffffffffffff1683614be0565b9150858181518110613f2b57613f2b614b2b565b60200260200101516020015160ff1683613f459190614be0565b9250600101613eab565b50613f5a8282614f1b565b95945050505050565b6000613f78836001600160a01b03841661416a565b9392505050565b60606000613f78836141b9565b6000613f78836001600160a01b038416614215565b8181808203613fb1575050505050565b6000856002613fc08787614f5f565b613fca9190614f7f565b613fd49087614b89565b81518110613fe457613fe4614b2b565b60200260200101516000015190505b818313614144575b806effffffffffffffffffffffffffffff1686848151811061401f5761401f614b2b565b6020026020010151600001516effffffffffffffffffffffffffffff161015614054578261404c81614fc9565b935050613ffb565b85828151811061406657614066614b2b565b6020026020010151600001516effffffffffffffffffffffffffffff16816effffffffffffffffffffffffffffff1610156140ad57816140a581614ffa565b925050614054565b81831361413f578582815181106140c6576140c6614b2b565b60200260200101518684815181106140e0576140e0614b2b565b60200260200101518785815181106140fa576140fa614b2b565b6020026020010188858151811061411357614113614b2b565b602002602001018290528290525050818061412d90614ffa565b925050828061413b90614fc9565b9350505b613ff3565b8185121561415757614157868684613fa1565b8383121561389557613895868486613fa1565b60008181526001830160205260408120546141b157508154600181810184556000848152602080822090930184905584548482528286019093526040902091909155610bb5565b506000610bb5565b60608160000180548060200260200160405190810160405280929190818152602001828054801561420957602002820191906000526020600020905b8154815260200190600101908083116141f5575b50505050509050919050565b600081815260018301602052604081205480156142fe576000614239600183614bb1565b855490915060009061424d90600190614bb1565b90508082146142b257600086600001828154811061426d5761426d614b2b565b906000526020600020015490508087600001848154811061429057614290614b2b565b6000918252602080832090910192909255918252600188019052604090208390555b85548690806142c3576142c361510e565b600190038181906000526020600020016000905590558560010160008681526020019081526020016000206000905560019350505050610bb5565b6000915050610bb5565b5092915050565b8280548282559060005260206000209081019282156143575760005260206000209182015b8281111561435757816143478482615184565b5091600101919060010190614334565b50614363929150614466565b5090565b828054828255906000526020600020908101928215614357579160200282015b82811115614357578251829061439d9082615264565b5091602001919060010190614387565b82805482825590600052602060002090810192821561441a579160200282015b8281111561441a57825182547fffffffffffffffffffffffff0000000000000000000000000000000000000000166001600160a01b039091161782556020909201916001909101906143cd565b50614363929150614483565b82805482825590600052602060002090810192821561441a5760005260206000209182015b8281111561441a57825482559160010191906001019061444b565b8082111561436357600061447a8282614498565b50600101614466565b5b808211156143635760008155600101614484565b5080546144a490614d11565b6000825580601f106144b4575050565b601f0160209004906000526020600020908101906108849190614483565b602080825282518282018190526000918401906040840190835b818110156145135783516001600160a01b03168352602093840193909201916001016144ec565b509095945050505050565b80356001600160a01b03811681146108f857600080fd5b60006020828403121561454757600080fd5b613f788261451e565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051601f8201601f1916810167ffffffffffffffff811182821017156145a8576145a8614550565b604052919050565b600082601f8301126145c157600080fd5b813567ffffffffffffffff8111156145db576145db614550565b6145ee6020601f19601f8401160161457f565b81815284602083860101111561460357600080fd5b816020850160208301376000918101602001919091529392505050565b60006020828403121561463257600080fd5b813567ffffffffffffffff81111561464957600080fd5b613965848285016145b0565b6000806040838503121561466857600080fd5b82359150602083013567ffffffffffffffff81111561468657600080fd5b614692858286016145b0565b9150509250929050565b600080604083850312156146af57600080fd5b50508035926020909101359150565b60ff8116811461088457600080fd5b80356108f8816146be565b6000806000806000608086880312156146f057600080fd5b85359450602086013567ffffffffffffffff81111561470e57600080fd5b8601601f8101881361471f57600080fd5b803567ffffffffffffffff81111561473657600080fd5b8860208260061b840101111561474b57600080fd5b6020919091019450925060408601359150614768606087016146cd565b90509295509295909350565b60006020828403121561478657600080fd5b5035919050565b600067ffffffffffffffff8211156147a7576147a7614550565b5060051b60200190565b6000602082840312156147c357600080fd5b813567ffffffffffffffff8111156147da57600080fd5b8201601f810184136147eb57600080fd5b80356147fe6147f98261478d565b61457f565b8082825260208201915060208360051b85010192508683111561482057600080fd5b602084015b8381101561486257803567ffffffffffffffff81111561484457600080fd5b614853896020838901016145b0565b84525060209283019201614825565b509695505050505050565b81516001600160a01b0316815260208083015161012083019161489a908401826001600160a01b03169052565b5060408301516040830152606083015160608301526080830151608083015260a083015160a083015260c083015160c083015260e083015160e083015261010083015161010083015292915050565b600080600080608085870312156148ff57600080fd5b5050823594602084013594506040840135936060013592509050565b600082601f83011261492c57600080fd5b813561493a6147f98261478d565b8082825260208201915060208360051b86010192508583111561495c57600080fd5b602085015b83811015614980576149728161451e565b835260209283019201614961565b5095945050505050565b60008060006060848603121561499f57600080fd5b833567ffffffffffffffff8111156149b657600080fd5b6149c28682870161491b565b935050602084013567ffffffffffffffff8111156149df57600080fd5b6149eb8682870161491b565b925050604084013567ffffffffffffffff811115614a0857600080fd5b614a148682870161491b565b9150509250925092565b60005b83811015614a39578181015183820152602001614a21565b50506000910152565b60008151808452614a5a816020860160208601614a1e565b601f01601f19169290920160200192915050565b600082825180855260208501945060208160051b8301016020850160005b83811015614abe57601f19858403018852614aa8838351614a42565b6020988901989093509190910190600101614a8c565b50909695505050505050565b602081526000613f786020830184614a6e565b60008060408385031215614af057600080fd5b823567ffffffffffffffff811115614b0757600080fd5b614b13858286016145b0565b925050614b226020840161451e565b90509250929050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b8082018281126000831280158216821582161715614ba957614ba9614b5a565b505092915050565b81810381811115610bb557610bb5614b5a565b60008251614bd6818460208701614a1e565b9190910192915050565b80820180821115610bb557610bb5614b5a565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b600082614c3157614c31614bf3565b500690565b6effffffffffffffffffffffffffffff8116811461088457600080fd5b6060808252810184905260008560808301825b87811015614cb4578235614c7981614c36565b6effffffffffffffffffffffffffffff1682526020830135614c9a816146be565b60ff16602083015260409283019290910190600101614c66565b50602084019590955250506001600160a01b039190911660409091015292915050565b600060208284031215614ce957600080fd5b8135613f78816146be565b600060208284031215614d0657600080fd5b8135613f7881614c36565b600181811c90821680614d2557607f821691505b602082108103612857577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b6000808354614d6c81614d11565b600182168015614d835760018114614d9857614dc8565b60ff1983168652811515820286019350614dc8565b86600052602060002060005b83811015614dc057815488820152600190910190602001614da4565b505081860193505b509195945050505050565b8135614dde81614c36565b6effffffffffffffffffffffffffffff811690508154817fffffffffffffffffffffffffffffffffff00000000000000000000000000000082161783556020840135614e29816146be565b6fff0000000000000000000000000000008160781b16837fffffffffffffffffffffffffffffffff000000000000000000000000000000008416171784555050505050565b604081526000614e816040830185614a6e565b90508260208301529392505050565b60006000198203614ea357614ea3614b5a565b5060010190565b600060208284031215614ebc57600080fd5b5051919050565b848152608060208201526000614edc6080830186614a42565b90508360408301526effffffffffffffffffffffffffffff8316606083015295945050505050565b8082028115828204841417610bb557610bb5614b5a565b600082614f2a57614f2a614bf3565b500490565b6001600160a01b038316815260606020820152600060608201526080604082015260006139656080830184614a42565b818103600083128015838313168383128216171561430857614308614b5a565b600082614f8e57614f8e614bf3565b60001983147f800000000000000000000000000000000000000000000000000000000000000083141615614fc457614fc4614b5a565b500590565b60007f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8203614ea357614ea3614b5a565b60007f8000000000000000000000000000000000000000000000000000000000000000820361502b5761502b614b5a565b506000190190565b6effffffffffffffffffffffffffffff8181168382160190811115610bb557610bb5614b5a565b60006effffffffffffffffffffffffffffff83168061507b5761507b614bf3565b806effffffffffffffffffffffffffffff84160491505092915050565b808202600082127f8000000000000000000000000000000000000000000000000000000000000000841416156150d0576150d0614b5a565b8181058314821517610bb557610bb5614b5a565b6effffffffffffffffffffffffffffff818116838216029081169081811461430857614308614b5a565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fd5b601f8211156120c857806000526020600020601f840160051c810160208510156151645750805b601f840160051c820191505b8181101561144d5760008155600101615170565b81810361518f575050565b6151998254614d11565b67ffffffffffffffff8111156151b1576151b1614550565b6151c5816151bf8454614d11565b8461513d565b6000601f8211600181146151fc57600083156151e15750848201545b600184901b600019600386901b1c198216175b85555061144d565b600085815260209020601f19841690600086815260209020845b838110156152365782860154825560019586019590910190602001615216565b50858310156152545781850154600019600388901b60f8161c191681555b5050505050600190811b01905550565b815167ffffffffffffffff81111561527e5761527e614550565b61528c816151bf8454614d11565b6020601f8211600181146152be57600083156151e1575081850151600184901b600019600386901b1c198216176151f4565b600084815260208120601f198516915b828110156152ee57878501518255602094850194600190920191016152ce565b508482101561530c5786840151600019600387901b60f8161c191681555b50505050600190811b0190555056fea2646970667358221220397c7e11019699c95916f85b08a3150696523f8159ab4f3bc820cc82275c34bc64736f6c634300081e0033",
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
func (r *Runner) DeployOracle0(opts *tests.RunOptions) (common.Address, uint64, *Oracle0, error) {
	parsed, err := Oracle0MetaData.GetAbi()
	if err != nil {
		return common.Address{}, 0, nil, err
	}
	if parsed == nil {
		return common.Address{}, 0, nil, errors.New("GetABI returned nil")
	}

	address, gasConsumed, c, data, err := (*tests.Runner)(r).DeployContract(opts, parsed, common.FromHex(Oracle0Bin))
	if err != nil {
		return common.Address{}, 0, nil, (&Oracle0{Contract: c}).DecodeError(data, err)
	}
	return address, gasConsumed, &Oracle0{Contract: c}, nil
}

// Oracle0 is an auto generated Go binding around an Ethereum contract.
type Oracle0 struct {
	*tests.Contract
}

// GetConfig is a free data retrieval call binding the contract method 0xc3f909d4.
//
// Solidity: function getConfig() view returns((address,address,uint256,int256,int256,uint256,uint256,uint256,uint256))
func (_Oracle0 *Oracle0) GetConfig(opts *tests.RunOptions) (Oracle0Config, uint64, error) {
	data, consumed, err := _Oracle0.Call(opts, "getConfig")

	if err != nil {
		return *new(Oracle0Config), consumed, _Oracle0.DecodeError(data, err)
	}
	out, err := _Oracle0.Contract.Abi().Unpack("getConfig", data)
	if err != nil {
		return *new(Oracle0Config), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(Oracle0Config)).(*Oracle0Config)
	return out0, consumed, err

}

// GetDecimals is a free data retrieval call binding the contract method 0xf0141d84.
//
// Solidity: function getDecimals() pure returns(uint8)
func (_Oracle0 *Oracle0) GetDecimals(opts *tests.RunOptions) (uint8, uint64, error) {
	data, consumed, err := _Oracle0.Call(opts, "getDecimals")

	if err != nil {
		return *new(uint8), consumed, _Oracle0.DecodeError(data, err)
	}
	out, err := _Oracle0.Contract.Abi().Unpack("getDecimals", data)
	if err != nil {
		return *new(uint8), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)
	return out0, consumed, err

}

// GetLastRoundBlock is a free data retrieval call binding the contract method 0x5a4d3a27.
//
// Solidity: function getLastRoundBlock() view returns(uint256)
func (_Oracle0 *Oracle0) GetLastRoundBlock(opts *tests.RunOptions) (*big.Int, uint64, error) {
	data, consumed, err := _Oracle0.Call(opts, "getLastRoundBlock")

	if err != nil {
		return *new(*big.Int), consumed, _Oracle0.DecodeError(data, err)
	}
	out, err := _Oracle0.Contract.Abi().Unpack("getLastRoundBlock", data)
	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// GetNewVotePeriod is a free data retrieval call binding the contract method 0x57eba759.
//
// Solidity: function getNewVotePeriod() view returns(uint256)
func (_Oracle0 *Oracle0) GetNewVotePeriod(opts *tests.RunOptions) (*big.Int, uint64, error) {
	data, consumed, err := _Oracle0.Call(opts, "getNewVotePeriod")

	if err != nil {
		return *new(*big.Int), consumed, _Oracle0.DecodeError(data, err)
	}
	out, err := _Oracle0.Contract.Abi().Unpack("getNewVotePeriod", data)
	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// GetNewVoters is a free data retrieval call binding the contract method 0x077945d3.
//
// Solidity: function getNewVoters() view returns(address[])
func (_Oracle0 *Oracle0) GetNewVoters(opts *tests.RunOptions) ([]common.Address, uint64, error) {
	data, consumed, err := _Oracle0.Call(opts, "getNewVoters")

	if err != nil {
		return *new([]common.Address), consumed, _Oracle0.DecodeError(data, err)
	}
	out, err := _Oracle0.Contract.Abi().Unpack("getNewVoters", data)
	if err != nil {
		return *new([]common.Address), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)
	return out0, consumed, err

}

// GetNonRevealThreshold is a free data retrieval call binding the contract method 0xed78349d.
//
// Solidity: function getNonRevealThreshold() view returns(uint256)
func (_Oracle0 *Oracle0) GetNonRevealThreshold(opts *tests.RunOptions) (*big.Int, uint64, error) {
	data, consumed, err := _Oracle0.Call(opts, "getNonRevealThreshold")

	if err != nil {
		return *new(*big.Int), consumed, _Oracle0.DecodeError(data, err)
	}
	out, err := _Oracle0.Contract.Abi().Unpack("getNonRevealThreshold", data)
	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// GetReports is a free data retrieval call binding the contract method 0xfb09917e.
//
// Solidity: function getReports(string _symbol, address _voter) view returns((uint120,uint8))
func (_Oracle0 *Oracle0) GetReports(opts *tests.RunOptions, _symbol string, _voter common.Address) (IOracleReport, uint64, error) {
	data, consumed, err := _Oracle0.Call(opts, "getReports", _symbol, _voter)

	if err != nil {
		return *new(IOracleReport), consumed, _Oracle0.DecodeError(data, err)
	}
	out, err := _Oracle0.Contract.Abi().Unpack("getReports", data)
	if err != nil {
		return *new(IOracleReport), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(IOracleReport)).(*IOracleReport)
	return out0, consumed, err

}

// GetRewardPeriodPerformance is a free data retrieval call binding the contract method 0x33d16293.
//
// Solidity: function getRewardPeriodPerformance(address _voter) view returns(uint256)
func (_Oracle0 *Oracle0) GetRewardPeriodPerformance(opts *tests.RunOptions, _voter common.Address) (*big.Int, uint64, error) {
	data, consumed, err := _Oracle0.Call(opts, "getRewardPeriodPerformance", _voter)

	if err != nil {
		return *new(*big.Int), consumed, _Oracle0.DecodeError(data, err)
	}
	out, err := _Oracle0.Contract.Abi().Unpack("getRewardPeriodPerformance", data)
	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// GetRound is a free data retrieval call binding the contract method 0x9f8743f7.
//
// Solidity: function getRound() view returns(uint256)
func (_Oracle0 *Oracle0) GetRound(opts *tests.RunOptions) (*big.Int, uint64, error) {
	data, consumed, err := _Oracle0.Call(opts, "getRound")

	if err != nil {
		return *new(*big.Int), consumed, _Oracle0.DecodeError(data, err)
	}
	out, err := _Oracle0.Contract.Abi().Unpack("getRound", data)
	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// GetRoundData is a free data retrieval call binding the contract method 0x3c8510fd.
//
// Solidity: function getRoundData(uint256 _round, string _symbol) view returns((uint256,uint256,uint256,bool) data)
func (_Oracle0 *Oracle0) GetRoundData(opts *tests.RunOptions, _round *big.Int, _symbol string) (IOracleRoundData, uint64, error) {
	data, consumed, err := _Oracle0.Call(opts, "getRoundData", _round, _symbol)

	if err != nil {
		return *new(IOracleRoundData), consumed, _Oracle0.DecodeError(data, err)
	}
	out, err := _Oracle0.Contract.Abi().Unpack("getRoundData", data)
	if err != nil {
		return *new(IOracleRoundData), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(IOracleRoundData)).(*IOracleRoundData)
	return out0, consumed, err

}

// GetSymbolUpdatedRound is a free data retrieval call binding the contract method 0x99b0014b.
//
// Solidity: function getSymbolUpdatedRound() view returns(int256)
func (_Oracle0 *Oracle0) GetSymbolUpdatedRound(opts *tests.RunOptions) (*big.Int, uint64, error) {
	data, consumed, err := _Oracle0.Call(opts, "getSymbolUpdatedRound")

	if err != nil {
		return *new(*big.Int), consumed, _Oracle0.DecodeError(data, err)
	}
	out, err := _Oracle0.Contract.Abi().Unpack("getSymbolUpdatedRound", data)
	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// GetSymbols is a free data retrieval call binding the contract method 0xdf7f710e.
//
// Solidity: function getSymbols() view returns(string[])
func (_Oracle0 *Oracle0) GetSymbols(opts *tests.RunOptions) ([]string, uint64, error) {
	data, consumed, err := _Oracle0.Call(opts, "getSymbols")

	if err != nil {
		return *new([]string), consumed, _Oracle0.DecodeError(data, err)
	}
	out, err := _Oracle0.Contract.Abi().Unpack("getSymbols", data)
	if err != nil {
		return *new([]string), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new([]string)).(*[]string)
	return out0, consumed, err

}

// GetVotePeriod is a free data retrieval call binding the contract method 0xb78dec52.
//
// Solidity: function getVotePeriod() view returns(uint256)
func (_Oracle0 *Oracle0) GetVotePeriod(opts *tests.RunOptions) (*big.Int, uint64, error) {
	data, consumed, err := _Oracle0.Call(opts, "getVotePeriod")

	if err != nil {
		return *new(*big.Int), consumed, _Oracle0.DecodeError(data, err)
	}
	out, err := _Oracle0.Contract.Abi().Unpack("getVotePeriod", data)
	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// GetVoterInfo is a free data retrieval call binding the contract method 0x9ed1f255.
//
// Solidity: function getVoterInfo(address _voter) view returns((uint256,uint256,uint256,uint256,bool,bool))
func (_Oracle0 *Oracle0) GetVoterInfo(opts *tests.RunOptions, _voter common.Address) (Oracle0VoterInfo, uint64, error) {
	data, consumed, err := _Oracle0.Call(opts, "getVoterInfo", _voter)

	if err != nil {
		return *new(Oracle0VoterInfo), consumed, _Oracle0.DecodeError(data, err)
	}
	out, err := _Oracle0.Contract.Abi().Unpack("getVoterInfo", data)
	if err != nil {
		return *new(Oracle0VoterInfo), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(Oracle0VoterInfo)).(*Oracle0VoterInfo)
	return out0, consumed, err

}

// GetVoterTreasuries is a free data retrieval call binding the contract method 0xef5cc4d1.
//
// Solidity: function getVoterTreasuries(address _oracleAddress) view returns(address)
func (_Oracle0 *Oracle0) GetVoterTreasuries(opts *tests.RunOptions, _oracleAddress common.Address) (common.Address, uint64, error) {
	data, consumed, err := _Oracle0.Call(opts, "getVoterTreasuries", _oracleAddress)

	if err != nil {
		return *new(common.Address), consumed, _Oracle0.DecodeError(data, err)
	}
	out, err := _Oracle0.Contract.Abi().Unpack("getVoterTreasuries", data)
	if err != nil {
		return *new(common.Address), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, consumed, err

}

// GetVoterValidators is a free data retrieval call binding the contract method 0x2d35d158.
//
// Solidity: function getVoterValidators(address _oracleAddress) view returns(address)
func (_Oracle0 *Oracle0) GetVoterValidators(opts *tests.RunOptions, _oracleAddress common.Address) (common.Address, uint64, error) {
	data, consumed, err := _Oracle0.Call(opts, "getVoterValidators", _oracleAddress)

	if err != nil {
		return *new(common.Address), consumed, _Oracle0.DecodeError(data, err)
	}
	out, err := _Oracle0.Contract.Abi().Unpack("getVoterValidators", data)
	if err != nil {
		return *new(common.Address), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, consumed, err

}

// GetVoters is a free data retrieval call binding the contract method 0xcdd72253.
//
// Solidity: function getVoters() view returns(address[])
func (_Oracle0 *Oracle0) GetVoters(opts *tests.RunOptions) ([]common.Address, uint64, error) {
	data, consumed, err := _Oracle0.Call(opts, "getVoters")

	if err != nil {
		return *new([]common.Address), consumed, _Oracle0.DecodeError(data, err)
	}
	out, err := _Oracle0.Contract.Abi().Unpack("getVoters", data)
	if err != nil {
		return *new([]common.Address), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)
	return out0, consumed, err

}

// LatestRoundData is a free data retrieval call binding the contract method 0x33f98c77.
//
// Solidity: function latestRoundData(string _symbol) view returns((uint256,uint256,uint256,bool) data)
func (_Oracle0 *Oracle0) LatestRoundData(opts *tests.RunOptions, _symbol string) (IOracleRoundData, uint64, error) {
	data, consumed, err := _Oracle0.Call(opts, "latestRoundData", _symbol)

	if err != nil {
		return *new(IOracleRoundData), consumed, _Oracle0.DecodeError(data, err)
	}
	out, err := _Oracle0.Contract.Abi().Unpack("latestRoundData", data)
	if err != nil {
		return *new(IOracleRoundData), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(IOracleRoundData)).(*IOracleRoundData)
	return out0, consumed, err

}

// DistributeRewards is a free data retrieval call for a paid mutator transaction binding the contract method 0x59974e38.
// Similar to eth_call from rpc calls or function.call from truffle, it reverts the state after the call and returns the output. The output is extracted
// the same way as done above for view only functions.
// Solidity: function distributeRewards(uint256 _ntn) payable returns()
func (_Oracle0 *Oracle0) CallDistributeRewards(r *tests.Runner, opts *tests.RunOptions, _ntn *big.Int) (uint64, error) {
	snap := r.Snapshot()

	data, consumed, err := _Oracle0.Call(opts, "distributeRewards", _ntn)
	r.RevertSnapshot(snap)
	return consumed, _Oracle0.DecodeError(data, err)

}

// Finalize is a free data retrieval call for a paid mutator transaction binding the contract method 0x4bb278f3.
// Similar to eth_call from rpc calls or function.call from truffle, it reverts the state after the call and returns the output. The output is extracted
// the same way as done above for view only functions.
// Solidity: function finalize() returns(bool)
func (_Oracle0 *Oracle0) CallFinalize(r *tests.Runner, opts *tests.RunOptions) (bool, uint64, error) {
	snap := r.Snapshot()

	data, consumed, err := _Oracle0.Call(opts, "finalize")
	r.RevertSnapshot(snap)

	if err != nil {
		return *new(bool), consumed, _Oracle0.DecodeError(data, err)
	}
	out, err := _Oracle0.Contract.Abi().Unpack("finalize", data)
	if err != nil {
		return *new(bool), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, consumed, err

}

// SetCommitRevealConfig is a free data retrieval call for a paid mutator transaction binding the contract method 0x3f422ef3.
// Similar to eth_call from rpc calls or function.call from truffle, it reverts the state after the call and returns the output. The output is extracted
// the same way as done above for view only functions.
// Solidity: function setCommitRevealConfig(uint256 _threshold, uint256 _resetInterval) returns()
func (_Oracle0 *Oracle0) CallSetCommitRevealConfig(r *tests.Runner, opts *tests.RunOptions, _threshold *big.Int, _resetInterval *big.Int) (uint64, error) {
	snap := r.Snapshot()

	data, consumed, err := _Oracle0.Call(opts, "setCommitRevealConfig", _threshold, _resetInterval)
	r.RevertSnapshot(snap)
	return consumed, _Oracle0.DecodeError(data, err)

}

// SetOperator is a free data retrieval call for a paid mutator transaction binding the contract method 0xb3ab15fb.
// Similar to eth_call from rpc calls or function.call from truffle, it reverts the state after the call and returns the output. The output is extracted
// the same way as done above for view only functions.
// Solidity: function setOperator(address _operator) returns()
func (_Oracle0 *Oracle0) CallSetOperator(r *tests.Runner, opts *tests.RunOptions, _operator common.Address) (uint64, error) {
	snap := r.Snapshot()

	data, consumed, err := _Oracle0.Call(opts, "setOperator", _operator)
	r.RevertSnapshot(snap)
	return consumed, _Oracle0.DecodeError(data, err)

}

// SetSlashingConfig is a free data retrieval call for a paid mutator transaction binding the contract method 0xda39fbfe.
// Similar to eth_call from rpc calls or function.call from truffle, it reverts the state after the call and returns the output. The output is extracted
// the same way as done above for view only functions.
// Solidity: function setSlashingConfig(int256 _outlierSlashingThreshold, int256 _outlierDetectionThreshold, uint256 _baseSlashingRate, uint256 _slashingRateCap) returns()
func (_Oracle0 *Oracle0) CallSetSlashingConfig(r *tests.Runner, opts *tests.RunOptions, _outlierSlashingThreshold *big.Int, _outlierDetectionThreshold *big.Int, _baseSlashingRate *big.Int, _slashingRateCap *big.Int) (uint64, error) {
	snap := r.Snapshot()

	data, consumed, err := _Oracle0.Call(opts, "setSlashingConfig", _outlierSlashingThreshold, _outlierDetectionThreshold, _baseSlashingRate, _slashingRateCap)
	r.RevertSnapshot(snap)
	return consumed, _Oracle0.DecodeError(data, err)

}

// SetSymbols is a free data retrieval call for a paid mutator transaction binding the contract method 0x8d4f75d2.
// Similar to eth_call from rpc calls or function.call from truffle, it reverts the state after the call and returns the output. The output is extracted
// the same way as done above for view only functions.
// Solidity: function setSymbols(string[] _symbols) returns()
func (_Oracle0 *Oracle0) CallSetSymbols(r *tests.Runner, opts *tests.RunOptions, _symbols []string) (uint64, error) {
	snap := r.Snapshot()

	data, consumed, err := _Oracle0.Call(opts, "setSymbols", _symbols)
	r.RevertSnapshot(snap)
	return consumed, _Oracle0.DecodeError(data, err)

}

// SetVotePeriod is a free data retrieval call for a paid mutator transaction binding the contract method 0x67b11630.
// Similar to eth_call from rpc calls or function.call from truffle, it reverts the state after the call and returns the output. The output is extracted
// the same way as done above for view only functions.
// Solidity: function setVotePeriod(uint256 _votePeriod) returns()
func (_Oracle0 *Oracle0) CallSetVotePeriod(r *tests.Runner, opts *tests.RunOptions, _votePeriod *big.Int) (uint64, error) {
	snap := r.Snapshot()

	data, consumed, err := _Oracle0.Call(opts, "setVotePeriod", _votePeriod)
	r.RevertSnapshot(snap)
	return consumed, _Oracle0.DecodeError(data, err)

}

// SetVoters is a free data retrieval call for a paid mutator transaction binding the contract method 0xda78110e.
// Similar to eth_call from rpc calls or function.call from truffle, it reverts the state after the call and returns the output. The output is extracted
// the same way as done above for view only functions.
// Solidity: function setVoters(address[] _newVoters, address[] _treasury, address[] _validator) returns()
func (_Oracle0 *Oracle0) CallSetVoters(r *tests.Runner, opts *tests.RunOptions, _newVoters []common.Address, _treasury []common.Address, _validator []common.Address) (uint64, error) {
	snap := r.Snapshot()

	data, consumed, err := _Oracle0.Call(opts, "setVoters", _newVoters, _treasury, _validator)
	r.RevertSnapshot(snap)
	return consumed, _Oracle0.DecodeError(data, err)

}

// UpdateVotersAndSymbol is a free data retrieval call for a paid mutator transaction binding the contract method 0x0f65875c.
// Similar to eth_call from rpc calls or function.call from truffle, it reverts the state after the call and returns the output. The output is extracted
// the same way as done above for view only functions.
// Solidity: function updateVotersAndSymbol() returns()
func (_Oracle0 *Oracle0) CallUpdateVotersAndSymbol(r *tests.Runner, opts *tests.RunOptions) (uint64, error) {
	snap := r.Snapshot()

	data, consumed, err := _Oracle0.Call(opts, "updateVotersAndSymbol")
	r.RevertSnapshot(snap)
	return consumed, _Oracle0.DecodeError(data, err)

}

// Vote is a free data retrieval call for a paid mutator transaction binding the contract method 0x56833ebe.
// Similar to eth_call from rpc calls or function.call from truffle, it reverts the state after the call and returns the output. The output is extracted
// the same way as done above for view only functions.
// Solidity: function vote(uint256 _commit, (uint120,uint8)[] _reports, uint256 _salt, uint8 _extra) returns()
func (_Oracle0 *Oracle0) CallVote(r *tests.Runner, opts *tests.RunOptions, _commit *big.Int, _reports []IOracleReport, _salt *big.Int, _extra uint8) (uint64, error) {
	snap := r.Snapshot()

	data, consumed, err := _Oracle0.Call(opts, "vote", _commit, _reports, _salt, _extra)
	r.RevertSnapshot(snap)
	return consumed, _Oracle0.DecodeError(data, err)

}

// DistributeRewards is a paid mutator transaction binding the contract method 0x59974e38.
//
// Solidity: function distributeRewards(uint256 _ntn) payable returns()
func (_Oracle0 *Oracle0) DistributeRewards(opts *tests.RunOptions, _ntn *big.Int) (uint64, error) {
	data, consumed, err := _Oracle0.Call(opts, "distributeRewards", _ntn)
	return consumed, _Oracle0.DecodeError(data, err)
}

// Finalize is a paid mutator transaction binding the contract method 0x4bb278f3.
//
// Solidity: function finalize() returns(bool)
func (_Oracle0 *Oracle0) Finalize(opts *tests.RunOptions) (uint64, error) {
	data, consumed, err := _Oracle0.Call(opts, "finalize")
	return consumed, _Oracle0.DecodeError(data, err)
}

// SetCommitRevealConfig is a paid mutator transaction binding the contract method 0x3f422ef3.
//
// Solidity: function setCommitRevealConfig(uint256 _threshold, uint256 _resetInterval) returns()
func (_Oracle0 *Oracle0) SetCommitRevealConfig(opts *tests.RunOptions, _threshold *big.Int, _resetInterval *big.Int) (uint64, error) {
	data, consumed, err := _Oracle0.Call(opts, "setCommitRevealConfig", _threshold, _resetInterval)
	return consumed, _Oracle0.DecodeError(data, err)
}

// SetOperator is a paid mutator transaction binding the contract method 0xb3ab15fb.
//
// Solidity: function setOperator(address _operator) returns()
func (_Oracle0 *Oracle0) SetOperator(opts *tests.RunOptions, _operator common.Address) (uint64, error) {
	data, consumed, err := _Oracle0.Call(opts, "setOperator", _operator)
	return consumed, _Oracle0.DecodeError(data, err)
}

// SetSlashingConfig is a paid mutator transaction binding the contract method 0xda39fbfe.
//
// Solidity: function setSlashingConfig(int256 _outlierSlashingThreshold, int256 _outlierDetectionThreshold, uint256 _baseSlashingRate, uint256 _slashingRateCap) returns()
func (_Oracle0 *Oracle0) SetSlashingConfig(opts *tests.RunOptions, _outlierSlashingThreshold *big.Int, _outlierDetectionThreshold *big.Int, _baseSlashingRate *big.Int, _slashingRateCap *big.Int) (uint64, error) {
	data, consumed, err := _Oracle0.Call(opts, "setSlashingConfig", _outlierSlashingThreshold, _outlierDetectionThreshold, _baseSlashingRate, _slashingRateCap)
	return consumed, _Oracle0.DecodeError(data, err)
}

// SetSymbols is a paid mutator transaction binding the contract method 0x8d4f75d2.
//
// Solidity: function setSymbols(string[] _symbols) returns()
func (_Oracle0 *Oracle0) SetSymbols(opts *tests.RunOptions, _symbols []string) (uint64, error) {
	data, consumed, err := _Oracle0.Call(opts, "setSymbols", _symbols)
	return consumed, _Oracle0.DecodeError(data, err)
}

// SetVotePeriod is a paid mutator transaction binding the contract method 0x67b11630.
//
// Solidity: function setVotePeriod(uint256 _votePeriod) returns()
func (_Oracle0 *Oracle0) SetVotePeriod(opts *tests.RunOptions, _votePeriod *big.Int) (uint64, error) {
	data, consumed, err := _Oracle0.Call(opts, "setVotePeriod", _votePeriod)
	return consumed, _Oracle0.DecodeError(data, err)
}

// SetVoters is a paid mutator transaction binding the contract method 0xda78110e.
//
// Solidity: function setVoters(address[] _newVoters, address[] _treasury, address[] _validator) returns()
func (_Oracle0 *Oracle0) SetVoters(opts *tests.RunOptions, _newVoters []common.Address, _treasury []common.Address, _validator []common.Address) (uint64, error) {
	data, consumed, err := _Oracle0.Call(opts, "setVoters", _newVoters, _treasury, _validator)
	return consumed, _Oracle0.DecodeError(data, err)
}

// UpdateVotersAndSymbol is a paid mutator transaction binding the contract method 0x0f65875c.
//
// Solidity: function updateVotersAndSymbol() returns()
func (_Oracle0 *Oracle0) UpdateVotersAndSymbol(opts *tests.RunOptions) (uint64, error) {
	data, consumed, err := _Oracle0.Call(opts, "updateVotersAndSymbol")
	return consumed, _Oracle0.DecodeError(data, err)
}

// Vote is a paid mutator transaction binding the contract method 0x56833ebe.
//
// Solidity: function vote(uint256 _commit, (uint120,uint8)[] _reports, uint256 _salt, uint8 _extra) returns()
func (_Oracle0 *Oracle0) Vote(opts *tests.RunOptions, _commit *big.Int, _reports []IOracleReport, _salt *big.Int, _extra uint8) (uint64, error) {
	data, consumed, err := _Oracle0.Call(opts, "vote", _commit, _reports, _salt, _extra)
	return consumed, _Oracle0.DecodeError(data, err)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
// WARNING! UNTESTED
// Solidity: fallback() payable returns()
func (_Oracle0 *Oracle0) Fallback(opts *tests.RunOptions, calldata []byte) (uint64, error) {
	out, consumed, err := _Oracle0.Call(opts, "", calldata)
	return consumed, _Oracle0.DecodeError(out, err)
}

// Receive is a paid mutator transaction binding the contract receive function.
// WARNING! UNTESTED
// Solidity: receive() payable returns()
func (_Oracle0 *Oracle0) Receive(opts *tests.RunOptions) (uint64, error) {
	out, consumed, err := _Oracle0.Call(opts, "")
	return consumed, _Oracle0.DecodeError(out, err)
}

func (_Oracle0 *Oracle0) DecodeError(data []byte, err error) error {
	if err == nil {
		return nil
	}

	reason, _ := abi.UnpackRevert(data)
	return fmt.Errorf("%w: %s", err, reason)

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
	*tests.Contract
}

func (_ReentrancyGuard *ReentrancyGuard) DecodeError(data []byte, err error) error {
	if err == nil {
		return nil
	}

	reason, _ := abi.UnpackRevert(data)
	return fmt.Errorf("%w: %s", err, reason)

}
