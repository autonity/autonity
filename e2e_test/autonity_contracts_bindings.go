// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package test

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

// AutonityAccountabilityEvent is an auto generated low-level Go binding around an user-defined struct.
type AutonityAccountabilityEvent struct {
	Chunks   uint8
	ChunkID  uint8
	Type     uint8
	Rule     uint8
	Reporter common.Address
	Sender   common.Address
	MsgHash  [32]byte
	RawProof []byte
}

// AutonityCommitteeMember is an auto generated low-level Go binding around an user-defined struct.
type AutonityCommitteeMember struct {
	Addr        common.Address
	VotingPower *big.Int
}

// AutonityConfig is an auto generated low-level Go binding around an user-defined struct.
type AutonityConfig struct {
	OperatorAccount common.Address
	TreasuryAccount common.Address
	TreasuryFee     *big.Int
	MinBaseFee      *big.Int
	DelegationRate  *big.Int
	EpochPeriod     *big.Int
	UnbondingPeriod *big.Int
	CommitteeSize   *big.Int
	ContractVersion *big.Int
	BlockPeriod     *big.Int
}

// AutonityStaking is an auto generated low-level Go binding around an user-defined struct.
type AutonityStaking struct {
	Delegator  common.Address
	Delegatee  common.Address
	Amount     *big.Int
	StartBlock *big.Int
}

// AutonityValidator is an auto generated low-level Go binding around an user-defined struct.
type AutonityValidator struct {
	Treasury          common.Address
	NodeAddress       common.Address
	OracleAddress     common.Address
	Enode             string
	CommissionRate    *big.Int
	BondedStake       *big.Int
	TotalSlashed      *big.Int
	LiquidContract    common.Address
	LiquidSupply      *big.Int
	RegistrationBlock *big.Int
	State             uint8
}

// IOracleRoundData is an auto generated low-level Go binding around an user-defined struct.
type IOracleRoundData struct {
	Round     *big.Int
	Price     *big.Int
	Timestamp *big.Int
	Status    *big.Int
}

// AutonityMetaData contains all meta data concerning the Autonity contract.
var AutonityMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"components\":[{\"internalType\":\"addresspayable\",\"name\":\"treasury\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"nodeAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"oracleAddress\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"enode\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"commissionRate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"bondedStake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalSlashed\",\"type\":\"uint256\"},{\"internalType\":\"contractLiquid\",\"name\":\"liquidContract\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"liquidSupply\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"registrationBlock\",\"type\":\"uint256\"},{\"internalType\":\"enumAutonity.ValidatorState\",\"name\":\"state\",\"type\":\"uint8\"}],\"internalType\":\"structAutonity.Validator[]\",\"name\":\"_validators\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"operatorAccount\",\"type\":\"address\"},{\"internalType\":\"addresspayable\",\"name\":\"treasuryAccount\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"treasuryFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"minBaseFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"delegationRate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"epochPeriod\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"unbondingPeriod\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"committeeSize\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"contractVersion\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"blockPeriod\",\"type\":\"uint256\"}],\"internalType\":\"structAutonity.Config\",\"name\":\"_config\",\"type\":\"tuple\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"uint8\",\"name\":\"Chunks\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"ChunkID\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"Type\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"Rule\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"Reporter\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"Sender\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"MsgHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"RawProof\",\"type\":\"bytes\"}],\"indexed\":false,\"internalType\":\"structAutonity.AccountabilityEvent\",\"name\":\"ev\",\"type\":\"tuple\"}],\"name\":\"AccusationAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"uint8\",\"name\":\"Chunks\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"ChunkID\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"Type\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"Rule\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"Reporter\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"Sender\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"MsgHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"RawProof\",\"type\":\"bytes\"}],\"indexed\":false,\"internalType\":\"structAutonity.AccountabilityEvent\",\"name\":\"ev\",\"type\":\"tuple\"}],\"name\":\"AccusationRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"BurnedStake\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"rate\",\"type\":\"uint256\"}],\"name\":\"CommissionRateChange\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"period\",\"type\":\"uint256\"}],\"name\":\"EpochPeriodUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasPrice\",\"type\":\"uint256\"}],\"name\":\"MinimumBaseFeeUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"MintedStake\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"uint8\",\"name\":\"Chunks\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"ChunkID\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"Type\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"Rule\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"Reporter\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"Sender\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"MsgHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"RawProof\",\"type\":\"bytes\"}],\"indexed\":false,\"internalType\":\"structAutonity.AccountabilityEvent\",\"name\":\"ev\",\"type\":\"tuple\"}],\"name\":\"MisbehaviourAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"penalty\",\"type\":\"uint256\"}],\"name\":\"MisbehaviourPenaltyUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"penalty\",\"type\":\"uint256\"}],\"name\":\"NodeSlashed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"treasury\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"effectiveBlock\",\"type\":\"uint256\"}],\"name\":\"PausedValidator\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"treasury\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"oracleAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"enode\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"liquidContract\",\"type\":\"address\"}],\"name\":\"RegisteredValidator\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Rewarded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"uint8\",\"name\":\"Chunks\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"ChunkID\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"Type\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"Rule\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"Reporter\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"Sender\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"MsgHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"RawProof\",\"type\":\"bytes\"}],\"indexed\":false,\"internalType\":\"structAutonity.AccountabilityEvent\",\"name\":\"ev\",\"type\":\"tuple\"}],\"name\":\"SubmitGuiltyAccusation\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"inputs\":[],\"name\":\"COMMISSION_RATE_PRECISION\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_msgHash\",\"type\":\"bytes32\"}],\"name\":\"accusationProcessed\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_address\",\"type\":\"address\"}],\"name\":\"activateValidator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_validator\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"bond\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_validator\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_rate\",\"type\":\"uint256\"}],\"name\":\"changeCommissionRate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"completeContractUpgrade\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"computeCommittee\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"config\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"operatorAccount\",\"type\":\"address\"},{\"internalType\":\"addresspayable\",\"name\":\"treasuryAccount\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"treasuryFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"minBaseFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"delegationRate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"epochPeriod\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"unbondingPeriod\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"committeeSize\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"contractVersion\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"blockPeriod\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"deployer\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"epochID\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"epochTotalBondedStake\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"finalize\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"votingPower\",\"type\":\"uint256\"}],\"internalType\":\"structAutonity.CommitteeMember[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_msgHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint8\",\"name\":\"_type\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"_rule\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"_reporter\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"_chunkID\",\"type\":\"uint8\"}],\"name\":\"getAccountabilityEventChunk\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBlockPeriod\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"startId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"lastId\",\"type\":\"uint256\"}],\"name\":\"getBondingReq\",\"outputs\":[{\"components\":[{\"internalType\":\"addresspayable\",\"name\":\"delegator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"delegatee\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"startBlock\",\"type\":\"uint256\"}],\"internalType\":\"structAutonity.Staking[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCommittee\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"votingPower\",\"type\":\"uint256\"}],\"internalType\":\"structAutonity.CommitteeMember[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCommitteeEnodes\",\"outputs\":[{\"internalType\":\"string[]\",\"name\":\"\",\"type\":\"string[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getEpochPeriod\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLastEpochBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getMaxCommitteeSize\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getMinimumBaseFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getNewContract\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"},{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getOperator\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPenalty\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"}],\"name\":\"getSlashedStake\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTreasuryAccount\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTreasuryFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getUnbondingPeriod\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"startId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"lastId\",\"type\":\"uint256\"}],\"name\":\"getUnbondingReq\",\"outputs\":[{\"components\":[{\"internalType\":\"addresspayable\",\"name\":\"delegator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"delegatee\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"startBlock\",\"type\":\"uint256\"}],\"internalType\":\"structAutonity.Staking[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"}],\"name\":\"getValidator\",\"outputs\":[{\"components\":[{\"internalType\":\"addresspayable\",\"name\":\"treasury\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"nodeAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"oracleAddress\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"enode\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"commissionRate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"bondedStake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalSlashed\",\"type\":\"uint256\"},{\"internalType\":\"contractLiquid\",\"name\":\"liquidContract\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"liquidSupply\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"registrationBlock\",\"type\":\"uint256\"},{\"internalType\":\"enumAutonity.ValidatorState\",\"name\":\"state\",\"type\":\"uint8\"}],\"internalType\":\"structAutonity.Validator\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"}],\"name\":\"getValidatorRecentAccusations\",\"outputs\":[{\"components\":[{\"internalType\":\"uint8\",\"name\":\"Chunks\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"ChunkID\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"Type\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"Rule\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"Reporter\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"Sender\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"MsgHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"RawProof\",\"type\":\"bytes\"}],\"internalType\":\"structAutonity.AccountabilityEvent[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"}],\"name\":\"getValidatorRecentMisbehaviours\",\"outputs\":[{\"components\":[{\"internalType\":\"uint8\",\"name\":\"Chunks\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"ChunkID\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"Type\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"Rule\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"Reporter\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"Sender\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"MsgHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"RawProof\",\"type\":\"bytes\"}],\"internalType\":\"structAutonity.AccountabilityEvent[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getValidators\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getVersion\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint8\",\"name\":\"Chunks\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"ChunkID\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"Type\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"Rule\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"Reporter\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"Sender\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"MsgHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"RawProof\",\"type\":\"bytes\"}],\"internalType\":\"structAutonity.AccountabilityEvent[]\",\"name\":\"_events\",\"type\":\"tuple[]\"}],\"name\":\"handleAccountabilityEvents\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"headBondingID\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"headUnbondingID\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastEpochBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"mint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_msgHash\",\"type\":\"bytes32\"}],\"name\":\"misbehaviourProcessed\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_address\",\"type\":\"address\"}],\"name\":\"pauseValidator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_enode\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"_oracleAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"_multisig\",\"type\":\"bytes\"}],\"name\":\"registerValidator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"resetContractUpgrade\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_size\",\"type\":\"uint256\"}],\"name\":\"setCommitteeSize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_period\",\"type\":\"uint256\"}],\"name\":\"setEpochPeriod\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_price\",\"type\":\"uint256\"}],\"name\":\"setMinimumBaseFee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_newPenalty\",\"type\":\"uint256\"}],\"name\":\"setMisbehaviourPenalty\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_account\",\"type\":\"address\"}],\"name\":\"setOperatorAccount\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"_account\",\"type\":\"address\"}],\"name\":\"setTreasuryAccount\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_treasuryFee\",\"type\":\"uint256\"}],\"name\":\"setTreasuryFee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_period\",\"type\":\"uint256\"}],\"name\":\"setUnbondingPeriod\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"tailBondingID\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"tailUnbondingID\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalRedistributed\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_validator\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"unbond\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"_bytecode\",\"type\":\"bytes\"},{\"internalType\":\"string\",\"name\":\"_abi\",\"type\":\"string\"}],\"name\":\"upgradeContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Sigs: map[string]string{
		"2f2c3f2e": "COMMISSION_RATE_PRECISION()",
		"ffd9d914": "accusationProcessed(bytes32)",
		"b46e5520": "activateValidator(address)",
		"dd62ed3e": "allowance(address,address)",
		"095ea7b3": "approve(address,uint256)",
		"70a08231": "balanceOf(address)",
		"a515366a": "bond(address,uint256)",
		"9dc29fac": "burn(address,uint256)",
		"852c4849": "changeCommissionRate(address,uint256)",
		"872cf059": "completeContractUpgrade()",
		"ae1f5fa0": "computeCommittee()",
		"79502c55": "config()",
		"313ce567": "decimals()",
		"d5f39488": "deployer()",
		"c9d97af4": "epochID()",
		"9c98e471": "epochTotalBondedStake()",
		"4bb278f3": "finalize()",
		"f446c557": "getAccountabilityEventChunk(bytes32,uint8,uint8,address,uint8)",
		"43645969": "getBlockPeriod()",
		"e485c6fb": "getBondingReq(uint256,uint256)",
		"ab8f6ffe": "getCommittee()",
		"a8b2216e": "getCommitteeEnodes()",
		"dfb1a4d2": "getEpochPeriod()",
		"731b3a03": "getLastEpochBlock()",
		"819b6463": "getMaxCommitteeSize()",
		"11220633": "getMinimumBaseFee()",
		"b66b3e79": "getNewContract()",
		"e7f43c68": "getOperator()",
		"e56e56db": "getPenalty()",
		"fe44c7f5": "getSlashedStake(address)",
		"f7866ee3": "getTreasuryAccount()",
		"29070c6d": "getTreasuryFee()",
		"6fd2c80b": "getUnbondingPeriod()",
		"55230e93": "getUnbondingReq(uint256,uint256)",
		"1904bb2e": "getValidator(address)",
		"1bd38702": "getValidatorRecentAccusations(address)",
		"ac306841": "getValidatorRecentMisbehaviours(address)",
		"b7ab4db5": "getValidators()",
		"0d8e6e2c": "getVersion()",
		"3a17914f": "handleAccountabilityEvents((uint8,uint8,uint8,uint8,address,address,bytes32,bytes)[])",
		"44697221": "headBondingID()",
		"4b0dff63": "headUnbondingID()",
		"c2362dd5": "lastEpochBlock()",
		"40c10f19": "mint(address,uint256)",
		"8a7b7f62": "misbehaviourProcessed(bytes32)",
		"06fdde03": "name()",
		"0ae65e7a": "pauseValidator(address)",
		"ad722d4d": "registerValidator(string,address,bytes)",
		"cf9c5719": "resetContractUpgrade()",
		"8bac7dad": "setCommitteeSize(uint256)",
		"6b5f444c": "setEpochPeriod(uint256)",
		"cb696f54": "setMinimumBaseFee(uint256)",
		"8f5d0fcb": "setMisbehaviourPenalty(uint256)",
		"520fdbbc": "setOperatorAccount(address)",
		"d886f8a2": "setTreasuryAccount(address)",
		"77e741c7": "setTreasuryFee(uint256)",
		"114eaf55": "setUnbondingPeriod(uint256)",
		"95d89b41": "symbol()",
		"787a2433": "tailBondingID()",
		"662cd7f4": "tailUnbondingID()",
		"9bb851c0": "totalRedistributed()",
		"18160ddd": "totalSupply()",
		"a9059cbb": "transfer(address,uint256)",
		"23b872dd": "transferFrom(address,address,uint256)",
		"a5d059ca": "unbond(address,uint256)",
		"b2ea9adb": "upgradeContract(bytes,string)",
	},
	Bin: "0x60806040526001600355600480546001600160a01b031916735a443704dd4b594b382c22a083e2bd3090a6fef3179055600060068190556007553480156200004657600080fd5b506040516200bdfc3803806200bdfc833981016040819052620000699162001a5a565b6010546200008f57602b80546001600160a01b031916331790556200008f828262000097565b505062001e9c565b8051600880546001600160a01b039283166001600160a01b0319918216179091556020830151600980549190931691161790556040810151600a556060810151600b556080810151600c5560a0810151600d5560c0810151600e5560e0810151600f5561010081015160105561012081015160115560005b82518110156200039757600083828151811062000130576200013062001c11565b602002602001015160a001519050600084838151811062000155576200015562001c11565b602002602001015161010001818152505060008483815181106200017d576200017d62001c11565b602002602001015160e001906001600160a01b031690816001600160a01b0316815250506000848381518110620001b857620001b862001c11565b602002602001015160a00181815250506000848381518110620001df57620001df62001c11565b60200260200101516101200181815250506008600401548483815181106200020b576200020b62001c11565b60200260200101516080018181525050600084838151811062000232576200023262001c11565b60200260200101516101400190600181111562000253576200025362001c27565b9081600181111562000269576200026962001c27565b815250506200029a84838151811062000286576200028662001c11565b6020026020010151620003b160201b60201c565b8060286000868581518110620002b457620002b462001c11565b6020026020010151600001516001600160a01b03166001600160a01b031681526020019081526020016000206000828254620002f1919062001c53565b9250508190555080602a60008282546200030c919062001c53565b9250508190555080601d600082825462000327919062001c53565b925050819055506200038184838151811062000347576200034762001c11565b6020026020010151602001518286858151811062000369576200036962001c11565b602002602001015160000151620003ca60201b60201c565b50806200038e8162001c6e565b9150506200010f565b50620003a26200055a565b620003ac62000609565b505050565b620003bc8162000e9a565b620003c78162000fd6565b50565b600082116200042c5760405162461bcd60e51b815260206004820152602360248201527f616d6f756e74206e65656420746f206265207374726963746c7920706f73697460448201526269766560e81b60648201526084015b60405180910390fd5b6001600160a01b038116600090815260286020526040902054821115620004965760405162461bcd60e51b815260206004820152601b60248201527f696e73756666696369656e74204e6577746f6e2062616c616e63650000000000604482015260640162000423565b6001600160a01b03811660009081526028602052604081208054849290620004c090849062001c8c565b9091555050604080516080810182526001600160a01b038084168252858116602080840191825283850187815243606086019081526024805460009081526022909452968320865181549087166001600160a01b0319918216178255945160018201805491909716951694909417909455516002830155915160039091015582549192906200054f8362001c6e565b919050555050505050565b6023545b60245481101562000589576200057481620011a6565b80620005808162001c6e565b9150506200055e565b50602454602355602654805b6027548110156200060357600e546000828152602560205260409020600301544391620005c29162001c53565b11620005e857620005d381620012bf565b620005e060018362001c53565b9150620005ee565b62000603565b80620005fa8162001c6e565b91505062000595565b50602655565b602b546060906001600160a01b03163314620006745760405162461bcd60e51b815260206004820152602360248201527f66756e6374696f6e207265737472696374656420746f207468652070726f746f60448201526218dbdb60ea1b606482015260840162000423565b601254620006c55760405162461bcd60e51b815260206004820152601860248201527f5468657265206d7573742062652076616c696461746f72730000000000000000604482015260640162000423565b6000805b601254811015620007a95760006029600060128481548110620006f057620006f062001c11565b60009182526020808320909101546001600160a01b031683528201929092526040019020600a015460ff1660018111156200072f576200072f62001c27565b1480156200077e57506000602960006012848154811062000754576200075462001c11565b60009182526020808320909101546001600160a01b03168352820192909252604001902060050154115b15620007945781620007908162001c6e565b9250505b80620007a08162001c6e565b915050620006c9565b50600f54818110620007b85750805b6000826001600160401b03811115620007d557620007d562001859565b6040519080825280602002602001820160405280156200086557816020015b6200085160408051610160810182526000808252602082018190529181018290526060808201526080810182905260a0810182905260c0810182905260e08101829052610100810182905261012081018290529061014082015290565b815260200190600190039081620007f45790505b5090506000826001600160401b0381111562000885576200088562001859565b6040519080825280602002602001820160405280156200091557816020015b6200090160408051610160810182526000808252602082018190529181018290526060808201526080810182905260a0810182905260c0810182905260e08101829052610100810182905261012081018290529061014082015290565b815260200190600190039081620008a45790505b5090506000836001600160401b0381111562000935576200093562001859565b6040519080825280602002602001820160405280156200095f578160200160208202803683370190505b5090506000805b60125481101562000bfd57600060296000601284815481106200098d576200098d62001c11565b60009182526020808320909101546001600160a01b031683528201929092526040019020600a015460ff166001811115620009cc57620009cc62001c27565b14801562000a1b575060006029600060128481548110620009f157620009f162001c11565b60009182526020808320909101546001600160a01b03168352820192909252604001902060050154115b1562000be8576000602960006012848154811062000a3d5762000a3d62001c11565b60009182526020808320909101546001600160a01b039081168452838201949094526040928301909120825161016081018452815485168152600182015485169281019290925260028101549093169181019190915260038201805491929160608401919062000aad9062001ca6565b80601f016020809104026020016040519081016040528092919081815260200182805462000adb9062001ca6565b801562000b2c5780601f1062000b005761010080835404028352916020019162000b2c565b820191906000526020600020905b81548152906001019060200180831162000b0e57829003601f168201915b505050918352505060048201546020820152600582015460408201526006820154606082015260078201546001600160a01b03166080820152600882015460a0820152600982015460c0820152600a82015460e09091019060ff16600181111562000b9b5762000b9b62001c27565b600181111562000baf5762000baf62001c27565b8152505090508086848151811062000bcb5762000bcb62001c11565b6020026020010181905250828062000be39062001c6e565b935050505b8062000bf48162001c6e565b91505062000966565b50600f548451111562000c7d5762000c15846200137a565b60005b600f5481101562000c765784818151811062000c385762000c3862001c11565b602002602001015184828151811062000c555762000c5562001c11565b6020026020010181905250808062000c6d9062001c6e565b91505062000c18565b5062000c81565b8392505b62000c8f601e6000620016bc565b62000c9d60206000620016df565b6000601d8190555b8581101562000e8e576000604051806040016040528086848151811062000cd05762000cd062001c11565b6020026020010151602001516001600160a01b0316815260200186848151811062000cff5762000cff62001c11565b60209081029190910181015160a00151909152601e805460018101825560009190915282517f50bb669a95c7b50b7e8a6f09454034b2b14cf2b85c730dca9a539ca82cb6e350600290920291820180546001600160a01b0319166001600160a01b03909216919091179055828201517f50bb669a95c7b50b7e8a6f09454034b2b14cf2b85c730dca9a539ca82cb6e3519091015586519192509086908490811062000dae5762000dae62001c11565b602090810291909101810151606001518254600181018455600093845292829020815162000de69491909101929190910190620016ff565b5084828151811062000dfc5762000dfc62001c11565b60200260200101516040015184838151811062000e1d5762000e1d62001c11565b60200260200101906001600160a01b031690816001600160a01b03168152505084828151811062000e525762000e5262001c11565b602002602001015160a00151601d600082825462000e71919062001c53565b9091555082915062000e8590508162001c6e565b91505062000ca5565b50909550505050505090565b600062000eb682606001516200139760201b620037521760201c565b6001600160a01b0390911660208401529050801562000f065760405162461bcd60e51b815260206004820152600b60248201526a32b737b2329032b93937b960a91b604482015260640162000423565b6020808301516001600160a01b0390811660009081526029909252604090912060010154161562000f7a5760405162461bcd60e51b815260206004820152601c60248201527f76616c696461746f7220616c7265616479207265676973746572656400000000604482015260640162000423565b6127108260800151111562000fd25760405162461bcd60e51b815260206004820152601760248201527f696e76616c696420636f6d6d697373696f6e2072617465000000000000000000604482015260640162000423565b5050565b60e08101516001600160a01b03166200106357600062001006601280549050620013e060201b620037921760201c565b90508160200151826000015183608001518360405162001026906200178e565b62001035949392919062001ce3565b604051809103906000f08015801562001052573d6000803e3d6000fd5b506001600160a01b031660e0830152505b60208082018051601280546001818101835560009283527fbb8a6a4669ba250d26cd7a459eca9d215f8307e33aebe50379bc5a3617ec344490910180546001600160a01b039485166001600160a01b03199182161790915584518416835260298652604092839020875181549086169083161781559451918501805492851692821692909217909155908501516002840180549190931691161790556060830151805184936200111b926003850192910190620016ff565b506080820151600482015560a0820151600582015560c0820151600682015560e08201516007820180546001600160a01b0319166001600160a01b0390921691909117905561010082015160088201556101208201516009820155610140820151600a8201805460ff1916600183818111156200119c576200119c62001c27565b0217905550505050565b600081815260226020908152604080832060018101546001600160a01b0316845260299092528220600581015491929091620011e85750600282015462001212565b81600501548360020154836008015462001203919062001d39565b6200120f919062001d71565b90505b600782015483546040516340c10f1960e01b81526001600160a01b039182166004820152602481018490529116906340c10f1990604401600060405180830381600087803b1580156200126457600080fd5b505af115801562001279573d6000803e3d6000fd5b50505050826002015482600501600082825462001297919062001c53565b9250508190555080826008016000828254620012b4919062001c53565b909155505050505050565b600081815260256020908152604080832060018101546001600160a01b03168452602990925282206008810154600582015460028401549394929362001306919062001d39565b62001312919062001d71565b9050808260050160008282546200132a919062001c8c565b909155505060028301546008830180546000906200134a90849062001c8c565b909155505082546001600160a01b031660009081526028602052604081208054839290620012b490849062001c53565b620003c78160006001845162001391919062001c8c565b620014fd565b600080620013a46200179c565b600060408286516020880160ff5afa620013bd57600080fd5b5080516020909101516c0100000000000000000000000090910494909350915050565b606081620014055750506040805180820190915260018152600360fc1b602082015290565b8160005b81156200143557806200141c8162001c6e565b91506200142d9050600a8362001d71565b915062001409565b6000816001600160401b0381111562001452576200145262001859565b6040519080825280601f01601f1916602001820160405280156200147d576020820181803683370190505b5090505b8415620014f5576200149560018362001c8c565b9150620014a4600a8662001d88565b620014b190603062001c53565b60f81b818381518110620014c957620014c962001c11565b60200101906001600160f81b031916908160001a905350620014ed600a8662001d71565b945062001481565b949350505050565b8181808214156200150f575050505050565b600085600262001520878762001d9f565b6200152c919062001de4565b62001538908762001e18565b815181106200154b576200154b62001c11565b602002602001015160a0015190505b81831362001688575b8086848151811062001579576200157962001c11565b602002602001015160a001511115620015a15782620015988162001e5f565b93505062001563565b858281518110620015b657620015b662001c11565b602002602001015160a00151811115620015df5781620015d68162001e7b565b925050620015a1565b8183136200168257858281518110620015fc57620015fc62001c11565b602002602001015186848151811062001619576200161962001c11565b602002602001015187858151811062001636576200163662001c11565b6020026020010188858151811062001652576200165262001c11565b60200260200101829052829052505082806200166e9062001e5f565b93505081806200167e9062001e7b565b9250505b6200155a565b818512156200169e576200169e868684620014fd565b83831215620016b457620016b4868486620014fd565b505050505050565b5080546000825560020290600052602060002090810190620003c79190620017ba565b5080546000825590600052602060002090810190620003c79190620017e2565b8280546200170d9062001ca6565b90600052602060002090601f0160209004810192826200173157600085556200177c565b82601f106200174c57805160ff19168380011785556200177c565b828001600101855582156200177c579182015b828111156200177c5782518255916020019190600101906200175f565b506200178a92915062001803565b5090565b61116a806200ac9283390190565b60405180604001604052806002906020820280368337509192915050565b5b808211156200178a5780546001600160a01b031916815560006001820155600201620017bb565b808211156200178a576000620017f982826200181a565b50600101620017e2565b5b808211156200178a576000815560010162001804565b508054620018289062001ca6565b6000825580601f1062001839575050565b601f016020900490600052602060002090810190620003c7919062001803565b634e487b7160e01b600052604160045260246000fd5b60405161014081016001600160401b038111828210171562001895576200189562001859565b60405290565b60405161016081016001600160401b038111828210171562001895576200189562001859565b604051601f8201601f191681016001600160401b0381118282101715620018ec57620018ec62001859565b604052919050565b80516001600160a01b03811681146200190c57600080fd5b919050565b60005b838110156200192e57818101518382015260200162001914565b838111156200193e576000848401525b50505050565b600082601f8301126200195657600080fd5b81516001600160401b0381111562001972576200197262001859565b62001987601f8201601f1916602001620018c1565b8181528460208386010111156200199d57600080fd5b620014f582602083016020870162001911565b8051600281106200190c57600080fd5b60006101408284031215620019d457600080fd5b620019de6200186f565b9050620019eb82620018f4565b8152620019fb60208301620018f4565b602082015260408201516040820152606082015160608201526080820151608082015260a082015160a082015260c082015160c082015260e082015160e082015261010080830151818301525061012080830151818301525092915050565b60008061016080848603121562001a7057600080fd5b83516001600160401b038082111562001a8857600080fd5b818601915086601f83011262001a9d57600080fd5b815160208282111562001ab45762001ab462001859565b8160051b62001ac5828201620018c1565b928352848101820192828101908b85111562001ae057600080fd5b83870192505b8483101562001bef5782518681111562001aff57600080fd5b8701808d03601f190189131562001b1557600080fd5b62001b1f6200189b565b62001b2c868301620018f4565b815262001b3c60408301620018f4565b8682015262001b4e60608301620018f4565b604082015260808201518881111562001b675760008081fd5b62001b778f888386010162001944565b60608301525060a080830151608083015260c0808401518284015260e0915081840151818401525061010062001baf818501620018f4565b828401526101209150818401518184015250610140808401518284015262001bd98c8501620019b0565b9083015250835250918301919083019062001ae6565b80995050505062001c0389828a01620019c0565b955050505050509250929050565b634e487b7160e01b600052603260045260246000fd5b634e487b7160e01b600052602160045260246000fd5b634e487b7160e01b600052601160045260246000fd5b6000821982111562001c695762001c6962001c3d565b500190565b600060001982141562001c855762001c8562001c3d565b5060010190565b60008282101562001ca15762001ca162001c3d565b500390565b600181811c9082168062001cbb57607f821691505b6020821081141562001cdd57634e487b7160e01b600052602260045260246000fd5b50919050565b600060018060a01b03808716835280861660208401525083604083015260806060830152825180608084015262001d228160a085016020870162001911565b601f01601f19169190910160a00195945050505050565b600081600019048311821515161562001d565762001d5662001c3d565b500290565b634e487b7160e01b600052601260045260246000fd5b60008262001d835762001d8362001d5b565b500490565b60008262001d9a5762001d9a62001d5b565b500690565b60008083128015600160ff1b85018412161562001dc05762001dc062001c3d565b6001600160ff1b038401831381161562001dde5762001dde62001c3d565b50500390565b60008262001df65762001df662001d5b565b600160ff1b82146000198414161562001e135762001e1362001c3d565b500590565b600080821280156001600160ff1b038490038513161562001e3d5762001e3d62001c3d565b600160ff1b839003841281161562001e595762001e5962001c3d565b50500190565b60006001600160ff1b0382141562001c855762001c8562001c3d565b6000600160ff1b82141562001e945762001e9462001c3d565b506000190190565b618de68062001eac6000396000f3fe608060405260043610620003ff5760003560e01c80638bac7dad116200020f578063b66b3e791162000123578063dd62ed3e11620000b3578063e7f43c681162000081578063e7f43c681462000cff578063f446c5571462000d1f578063f7866ee31462000d44578063fe44c7f51462000d64578063ffd9d9141462000d8957005b8063dd62ed3e1462000c62578063dfb1a4d21462000cac578063e485c6fb1462000cc3578063e56e56db1462000ce857005b8063cb696f5411620000f1578063cb696f541462000bc5578063cf9c57191462000bea578063d5f394881462000c02578063d886f8a21462000c3d57005b8063b66b3e791462000b55578063b7ab4db51462000b7d578063c2362dd51462000b95578063c9d97af41462000bad57005b8063a8b2216e116200019f578063ad722d4d116200016d578063ad722d4d1462000abf578063ae1f5fa01462000ae4578063b2ea9adb1462000b0b578063b46e55201462000b3057005b8063a8b2216e1462000a27578063a9059cbb1462000a4e578063ab8f6ffe1462000a73578063ac3068411462000a9a57005b80639c98e47111620001dd5780639c98e47114620009a05780639dc29fac14620009b8578063a515366a14620009dd578063a5d059ca1462000a0257005b80638bac7dad14620009105780638f5d0fcb146200093557806395d89b41146200095a5780639bb851c0146200098857005b806344697221116200031357806370a0823111620002a357806379502c55116200027157806379502c5514620007ea578063819b64631462000889578063852c484914620008a0578063872cf05914620008c55780638a7b7f6214620008dd57005b806370a08231146200075c578063731b3a03146200079657806377e741c714620007ad578063787a243314620007d257005b806355230e9311620002e157806355230e9314620006d4578063662cd7f414620007085780636b5f444c14620007205780636fd2c80b146200074557005b80634469722114620006575780634b0dff63146200066f5780634bb278f31462000687578063520fdbbc14620006af57005b80631bd38702116200038f578063313ce567116200035d578063313ce56714620005d85780633a17914f14620005f657806340c10f19146200061b57806343645969146200064057005b80631bd38702146200055057806323b872dd146200058457806329070c6d14620005a95780632f2c3f2e14620005c057005b80631122063311620003cd5780631122063314620004c9578063114eaf5514620004e057806318160ddd14620005055780631904bb2e146200051c57005b806306fdde031462000409578063095ea7b3146200044d5780630ae65e7a14620004835780630d8e6e2c14620004a857005b366200040757005b005b3480156200041657600080fd5b506040805180820190915260068152652732bbba37b760d11b60208201525b60405162000444919062006e99565b60405180910390f35b3480156200045a57600080fd5b50620004726200046c36600462006edd565b62000dbc565b604051901515815260200162000444565b3480156200049057600080fd5b5062000407620004a236600462006f0c565b62000dd4565b348015620004b557600080fd5b506010545b60405190815260200162000444565b348015620004d657600080fd5b50600b54620004ba565b348015620004ed57600080fd5b5062000407620004ff36600462006f2c565b62000e67565b3480156200051257600080fd5b50602a54620004ba565b3480156200052957600080fd5b50620005416200053b36600462006f0c565b62000e99565b60405162000444919062006f7f565b3480156200055d57600080fd5b50620005756200056f36600462006f0c565b62001051565b604051620004449190620070e3565b3480156200059157600080fd5b5062000472620005a336600462007149565b62001450565b348015620005b657600080fd5b50600a54620004ba565b348015620005cd57600080fd5b50620004ba61271081565b348015620005e557600080fd5b506040516012815260200162000444565b3480156200060357600080fd5b5062000407620006153660046200728c565b620014aa565b3480156200062857600080fd5b50620004076200063a36600462006edd565b6200180a565b3480156200064d57600080fd5b50601154620004ba565b3480156200066457600080fd5b50620004ba60245481565b3480156200067c57600080fd5b50620004ba60275481565b3480156200069457600080fd5b506200069f620018c8565b6040516200044492919062007469565b348015620006bc57600080fd5b5062000407620006ce36600462006f0c565b62001ac5565b348015620006e157600080fd5b50620006f9620006f336600462007486565b62001b6c565b604051620004449190620074a9565b3480156200071557600080fd5b50620004ba60265481565b3480156200072d57600080fd5b50620004076200073f36600462006f2c565b62001c81565b3480156200075257600080fd5b50600e54620004ba565b3480156200076957600080fd5b50620004ba6200077b36600462006f0c565b6001600160a01b031660009081526028602052604090205490565b348015620007a357600080fd5b50601c54620004ba565b348015620007ba57600080fd5b5062000407620007cc36600462006f2c565b62001daf565b348015620007df57600080fd5b50620004ba60235481565b348015620007f757600080fd5b50600854600954600a54600b54600c54600d54600e54600f5460105460115462000833996001600160a01b03908116991697969594939291908a565b604080516001600160a01b039b8c1681529a90991660208b0152978901969096526060880194909452608087019290925260a086015260c085015260e08401526101008301526101208201526101400162000444565b3480156200089657600080fd5b50600f54620004ba565b348015620008ad57600080fd5b5062000407620008bf36600462006edd565b62001de1565b348015620008d257600080fd5b506200040762001f76565b348015620008ea57600080fd5b5062000472620008fc36600462006f2c565b600090815260166020526040902054151590565b3480156200091d57600080fd5b50620004076200092f36600462006f2c565b62001fb2565b3480156200094257600080fd5b50620004076200095436600462006f2c565b62002036565b3480156200096757600080fd5b50604080518082019091526003815262272a2760e91b602082015262000435565b3480156200099557600080fd5b50620004ba601f5481565b348015620009ad57600080fd5b50620004ba601d5481565b348015620009c557600080fd5b5062000407620009d736600462006edd565b620020a5565b348015620009ea57600080fd5b5062000407620009fc36600462006edd565b620021bf565b34801562000a0f57600080fd5b506200040762000a2136600462006edd565b6200228e565b34801562000a3457600080fd5b5062000a3f620022db565b6040516200044491906200751b565b34801562000a5b57600080fd5b506200047262000a6d36600462006edd565b620023be565b34801562000a8057600080fd5b5062000a8b620023cd565b60405162000444919062007574565b34801562000aa757600080fd5b506200057562000ab936600462006f0c565b6200243b565b34801562000acc57600080fd5b506200040762000ade36600462007589565b62002827565b34801562000af157600080fd5b5062000afc620028e5565b60405162000444919062007609565b34801562000b1857600080fd5b506200040762000b2a36600462007658565b62003095565b34801562000b3d57600080fd5b506200040762000b4f36600462006f0c565b620030dc565b34801562000b6257600080fd5b5062000b6d62003200565b60405162000444929190620076c2565b34801562000b8a57600080fd5b5062000afc62003337565b34801562000ba257600080fd5b50620004ba601c5481565b34801562000bba57600080fd5b50620004ba601b5481565b34801562000bd257600080fd5b506200040762000be436600462006f2c565b6200339b565b34801562000bf757600080fd5b5062000407620033fe565b34801562000c0f57600080fd5b50602b5462000c24906001600160a01b031681565b6040516001600160a01b03909116815260200162000444565b34801562000c4a57600080fd5b506200040762000c5c36600462006f0c565b62003452565b34801562000c6f57600080fd5b50620004ba62000c81366004620076eb565b6001600160a01b03918216600090815260216020908152604080832093909416825291909152205490565b34801562000cb957600080fd5b50600d54620004ba565b34801562000cd057600080fd5b50620006f962000ce236600462007486565b620034a1565b34801562000cf557600080fd5b50600354620004ba565b34801562000d0c57600080fd5b506008546001600160a01b031662000c24565b34801562000d2c57600080fd5b506200043562000d3e36600462007729565b620035b6565b34801562000d5157600080fd5b506009546001600160a01b031662000c24565b34801562000d7157600080fd5b50620004ba62000d8336600462006f0c565b620036f6565b34801562000d9657600080fd5b506200047262000da836600462006f2c565b600090815260176020526040902054151590565b600062000dcb338484620038af565b50600192915050565b6001600160a01b038082166000818152602960205260409020600101549091161462000e1d5760405162461bcd60e51b815260040162000e149062007792565b60405180910390fd5b6001600160a01b0381811660009081526029602052604090205416331462000e595760405162461bcd60e51b815260040162000e1490620077c9565b62000e6481620039d8565b50565b6008546001600160a01b0316331462000e945760405162461bcd60e51b815260040162000e149062007815565b600e55565b62000ea362006ab6565b6001600160a01b038083166000818152602960205260409020600101549091161462000ee35760405162461bcd60e51b815260040162000e14906200784c565b6001600160a01b03808316600090815260296020908152604091829020825161016081018452815485168152600182015485169281019290925260028101549093169181019190915260038201805491929160608401919062000f469062007883565b80601f016020809104026020016040519081016040528092919081815260200182805462000f749062007883565b801562000fc55780601f1062000f995761010080835404028352916020019162000fc5565b820191906000526020600020905b81548152906001019060200180831162000fa757829003601f168201915b505050918352505060048201546020820152600582015460408201526006820154606082015260078201546001600160a01b03166080820152600882015460a0820152600982015460c0820152600a82015460e09091019060ff16600181111562001034576200103462006f46565b600181111562001048576200104862006f46565b90525092915050565b6001600160a01b038082166000818152602960205260409020600101546060921614620010925760405162461bcd60e51b815260040162000e149062007792565b6001600160a01b0382166000908152601560205260409020546101001015620012d9576001600160a01b038216600090815260156020526040812054620010dd9061010090620078d6565b60408051610100808252612020820190925291925060009190816020015b6200110562006b0f565b815260200190600190039081620010fb579050509050815b6001600160a01b038516600090815260156020526040902054811015620012d1576001600160a01b03851660009081526015602052604081208054839081106200116b576200116b620078f0565b60009182526020918290206040805161010080820183526004909402909201805460ff80821685529481048516958401959095526201000085048416918301919091526301000000840490921660608201526001600160a01b03600160201b90930483166080820152600182015490921660a0830152600281015460c083015260038101805460e084019190620012029062007883565b80601f0160208091040260200160405190810160405280929190818152602001828054620012309062007883565b8015620012815780601f10620012555761010080835404028352916020019162001281565b820191906000526020600020905b8154815290600101906020018083116200126357829003601f168201915b5050505050815250509050808385846200129c9190620078d6565b81518110620012af57620012af620078f0565b6020026020010181905250508080620012c89062007906565b9150506200111d565b509392505050565b6001600160a01b038216600090815260156020908152604080832080548251818502810185019093528083529193909284015b828210156200144557600084815260209081902060408051610100808201835260048702909301805460ff8082168452948104851695830195909552620100008504841692820192909252630100000084049092166060830152600160201b9092046001600160a01b03908116608083015260018301541660a0820152600282015460c082015260038201805491929160e084019190620013ad9062007883565b80601f0160208091040260200160405190810160405280929190818152602001828054620013db9062007883565b80156200142c5780601f1062001400576101008083540402835291602001916200142c565b820191906000526020600020905b8154815290600101906020018083116200140e57829003601f168201915b505050505081525050815260200190600101906200130c565b505050509050919050565b60006200145f84848462003ad4565b6001600160a01b038416600090815260216020908152604080832033845290915281205462001490908490620078d6565b90506200149f853383620038af565b506001949350505050565b336000818152602960205260409020600101546001600160a01b031614620015215760405162461bcd60e51b8152602060048201526024808201527f66756e6374696f6e207265737472696374656420746f207468652076616c696460448201526330ba37b960e11b606482015260840162000e14565b60005b81518110156200180657336001600160a01b03168282815181106200154d576200154d620078f0565b6020026020010151608001516001600160a01b0316146200156e57620017f1565b818181518110620015835762001583620078f0565b60200260200101516000015160ff16600014620015c757620015c1828281518110620015b357620015b3620078f0565b602002602001015162003bdd565b620017f1565b6000828281518110620015de57620015de620078f0565b60200260200101516040015160ff16600281111562001601576200160162006f46565b600281111562001615576200161562006f46565b1415620016805760166000838381518110620016355762001635620078f0565b602002602001015160c00151815260200190815260200160002054600014156200168057620015c1828281518110620016725762001672620078f0565b6020026020010151620040a4565b6001828281518110620016975762001697620078f0565b60200260200101516040015160ff166002811115620016ba57620016ba62006f46565b6002811115620016ce57620016ce62006f46565b1415620017395760176000838381518110620016ee57620016ee620078f0565b602002602001015160c00151815260200190815260200160002054600014156200173957620015c18282815181106200172b576200172b620078f0565b60200260200101516200440a565b6002828281518110620017505762001750620078f0565b60200260200101516040015160ff16600281111562001773576200177362006f46565b600281111562001787576200178762006f46565b1415620017f15760176000838381518110620017a757620017a7620078f0565b602002602001015160c00151815260200190815260200160002054600014620017f157620015c1828281518110620017e357620017e3620078f0565b6020026020010151620045d4565b80620017fd8162007906565b91505062001524565b5050565b6008546001600160a01b03163314620018375760405162461bcd60e51b815260040162000e149062007815565b6001600160a01b038216600090815260286020526040812080548392906200186190849062007924565b9250508190555080602a60008282546200187c919062007924565b9091555050604080516001600160a01b0384168152602081018390527f48490b4407bb949b708ec5f514b4167f08f4969baaf78d53b05028adf369bfcf91015b60405180910390a15050565b602b546000906060906001600160a01b03163314620018fb5760405162461bcd60e51b815260040162000e14906200793f565b62001905620046ad565b6200190f62004a5e565b600d54601c544391620019229162007924565b1415620019e6576200193362005051565b6200193d620050c6565b62001947620051a3565b6200195162005252565b60006200195d620028e5565b6004805460405163422811f960e11b81529293506001600160a01b03169163845023f2916200198f9185910162007609565b600060405180830381600087803b158015620019aa57600080fd5b505af1158015620019bf573d6000803e3d6000fd5b5050505043601c819055506001601b6000828254620019df919062007924565b9091555050505b6004805460408051634bb278f360e01b815290516001600160a01b0390921692634bb278f392828201926000929082900301818387803b15801562001a2a57600080fd5b505af115801562001a3f573d6000803e3d6000fd5b5050600254601e80546040805160208084028201810190925282815260ff9094169550919350839160009084015b8282101562001ab7576000848152602090819020604080518082019091526002850290910180546001600160a01b0316825260019081015482840152908352909201910162001a6d565b505050509050915091509091565b6008546001600160a01b0316331462001af25760405162461bcd60e51b815260040162000e149062007815565b600880546001600160a01b0319166001600160a01b038381169182179092556004805460405163b3ab15fb60e01b81529182019290925291169063b3ab15fb90602401600060405180830381600087803b15801562001b5057600080fd5b505af115801562001b65573d6000803e3d6000fd5b5050505050565b6060600062001b7c8484620078d6565b6001600160401b0381111562001b965762001b966200718f565b60405190808252806020026020018201604052801562001bd357816020015b62001bbf62006b55565b81526020019060019003908162001bb55790505b50905060005b62001be58585620078d6565b811015620012d1576025600062001bfd838862007924565b81526020808201929092526040908101600020815160808101835281546001600160a01b03908116825260018301541693810193909352600281015491830191909152600301546060820152825183908390811062001c605762001c60620078f0565b6020026020010181905250808062001c789062007906565b91505062001bd9565b6008546001600160a01b0316331462001cae5760405162461bcd60e51b815260040162000e149062007815565b600d5481141562001cbc5750565b600d5481101562001d735780601c5462001cd7919062007924565b431062001d735760405162461bcd60e51b815260206004820152605760248201527f63757272656e7420636861696e2068656164206578636565642074686520776960448201527f6e646f773a206c617374426c6f636b45706f6368202b205f6e6577506572696f60648201527f642c2074727920616761696e206c6174746572206f6e2e000000000000000000608482015260a40162000e14565b600d8190556040518181527fd7f1279ded354dbf22a69fcc2fd661763a6e2956a5d2891af9410af880fa5f81906020015b60405180910390a150565b6008546001600160a01b0316331462001ddc5760405162461bcd60e51b815260040162000e149062007815565b600a55565b6001600160a01b038083166000818152602960205260409020600101549091161462001e215760405162461bcd60e51b815260040162000e149062007792565b6001600160a01b0382811660009081526029602052604090205416331462001e5d5760405162461bcd60e51b815260040162000e1490620077c9565b61271081111562001eb15760405162461bcd60e51b815260206004820152601f60248201527f7265717569726520636f727265637420636f6d6d697373696f6e207261746500604482015260640162000e14565b604080516060810182526001600160a01b0384811682524360208084019182528385018681526007805460009081526005909352958220855181546001600160a01b03191695169490941784559151600180850191909155915160029093019290925583549293909290919062001f2a90849062007924565b9091555050604080516001600160a01b0385168152602081018490527f4fba51c92fa3d6ad8374d394f6cd5766857552e153d7384a8f23aa4ce9a8a7cf910160405180910390a1505050565b6008546001600160a01b0316331462001fa35760405162461bcd60e51b815260040162000e149062007815565b6002805460ff19166001179055565b6008546001600160a01b0316331462001fdf5760405162461bcd60e51b815260040162000e149062007815565b60008111620020315760405162461bcd60e51b815260206004820152601960248201527f636f6d6d69747465652073697a652063616e2774206265203000000000000000604482015260640162000e14565b600f55565b6008546001600160a01b03163314620020635760405162461bcd60e51b815260040162000e149062007815565b600081116200206f5750565b60038190556040518181527f3e4df1a42f35d79ea7cc3833604cd6377005ec5985514c94c229fb83f36507039060200162001da4565b6008546001600160a01b03163314620020d25760405162461bcd60e51b815260040162000e149062007815565b6001600160a01b038216600090815260286020526040902054811115620021355760405162461bcd60e51b8152602060048201526016602482015275416d6f756e7420657863656564732062616c616e636560501b604482015260640162000e14565b6001600160a01b038216600090815260286020526040812080548392906200215f908490620078d6565b9250508190555080602a60008282546200217a9190620078d6565b9091555050604080516001600160a01b0384168152602081018390527f5024dbeedf0c06664c9bd7be836915730c955e936972c020683dadf11d5488a39101620018bc565b6001600160a01b0380831660008181526029602052604090206001015490911614620021ff5760405162461bcd60e51b815260040162000e14906200784c565b6001600160a01b0382166000908152602960205260408120600a015460ff16600181111562002232576200223262006f46565b14620022815760405162461bcd60e51b815260206004820152601b60248201527f76616c696461746f72206e65656420746f206265206163746976650000000000604482015260640162000e14565b620018068282336200536f565b6001600160a01b0380831660008181526029602052604090206001015490911614620022ce5760405162461bcd60e51b815260040162000e14906200784c565b62001806828233620054fb565b60606020805480602002602001604051908101604052809291908181526020016000905b82821015620023b5578382906000526020600020018054620023219062007883565b80601f01602080910402602001604051908101604052809291908181526020018280546200234f9062007883565b8015620023a05780601f106200237457610100808354040283529160200191620023a0565b820191906000526020600020905b8154815290600101906020018083116200238257829003601f168201915b505050505081526020019060010190620022ff565b50505050905090565b600062000dcb33848462003ad4565b6060601e805480602002602001604051908101604052809291908181526020016000905b82821015620023b5576000848152602090819020604080518082019091526002850290910180546001600160a01b03168252600190810154828401529083529092019101620023f1565b6001600160a01b0380821660008181526029602052604090206001015460609216146200247c5760405162461bcd60e51b815260040162000e149062007792565b6001600160a01b0382166000908152601460205260409020546101001015620026bb576001600160a01b038216600090815260146020526040812054620024c79061010090620078d6565b60408051610100808252612020820190925291925060009190816020015b620024ef62006b0f565b815260200190600190039081620024e5579050509050815b6001600160a01b038516600090815260146020526040902054811015620012d1576001600160a01b0385166000908152601460205260408120805483908110620025555762002555620078f0565b60009182526020918290206040805161010080820183526004909402909201805460ff80821685529481048516958401959095526201000085048416918301919091526301000000840490921660608201526001600160a01b03600160201b90930483166080820152600182015490921660a0830152600281015460c083015260038101805460e084019190620025ec9062007883565b80601f01602080910402602001604051908101604052809291908181526020018280546200261a9062007883565b80156200266b5780601f106200263f576101008083540402835291602001916200266b565b820191906000526020600020905b8154815290600101906020018083116200264d57829003601f168201915b505050505081525050905080838584620026869190620078d6565b81518110620026995762002699620078f0565b6020026020010181905250508080620026b29062007906565b91505062002507565b6001600160a01b038216600090815260146020908152604080832080548251818502810185019093528083529193909284015b828210156200144557600084815260209081902060408051610100808201835260048702909301805460ff8082168452948104851695830195909552620100008504841692820192909252630100000084049092166060830152600160201b9092046001600160a01b03908116608083015260018301541660a0820152600282015460c082015260038201805491929160e0840191906200278f9062007883565b80601f0160208091040260200160405190810160405280929190818152602001828054620027bd9062007883565b80156200280e5780601f10620027e2576101008083540402835291602001916200280e565b820191906000526020600020905b815481529060010190602001808311620027f057829003601f168201915b50505050508152505081526020019060010190620026ee565b60408051610160810182523381526000602082018190526001600160a01b0385169282019290925260608101859052600c54608082015260a0810182905260c0810182905260e08101829052610100810182905243610120820152610140810191909152620028978183620057e9565b602081015160e08201516040517f8ad8bd2eb6950e5f332fd3a6dca48cb358ecfe3057848902b98cbdfe455c915c92620028d792339288918a9162007982565b60405180910390a150505050565b602b546060906001600160a01b03163314620029155760405162461bcd60e51b815260040162000e14906200793f565b601254620029665760405162461bcd60e51b815260206004820152601860248201527f5468657265206d7573742062652076616c696461746f72730000000000000000604482015260640162000e14565b6000805b60125481101562002a4a5760006029600060128481548110620029915762002991620078f0565b60009182526020808320909101546001600160a01b031683528201929092526040019020600a015460ff166001811115620029d057620029d062006f46565b14801562002a1f575060006029600060128481548110620029f557620029f5620078f0565b60009182526020808320909101546001600160a01b03168352820192909252604001902060050154115b1562002a35578162002a318162007906565b9250505b8062002a418162007906565b9150506200296a565b50600f5481811062002a595750805b6000826001600160401b0381111562002a765762002a766200718f565b60405190808252806020026020018201604052801562002ab357816020015b62002a9f62006ab6565b81526020019060019003908162002a955790505b5090506000826001600160401b0381111562002ad35762002ad36200718f565b60405190808252806020026020018201604052801562002b1057816020015b62002afc62006ab6565b81526020019060019003908162002af25790505b5090506000836001600160401b0381111562002b305762002b306200718f565b60405190808252806020026020018201604052801562002b5a578160200160208202803683370190505b5090506000805b60125481101562002df8576000602960006012848154811062002b885762002b88620078f0565b60009182526020808320909101546001600160a01b031683528201929092526040019020600a015460ff16600181111562002bc75762002bc762006f46565b14801562002c1657506000602960006012848154811062002bec5762002bec620078f0565b60009182526020808320909101546001600160a01b03168352820192909252604001902060050154115b1562002de3576000602960006012848154811062002c385762002c38620078f0565b60009182526020808320909101546001600160a01b039081168452838201949094526040928301909120825161016081018452815485168152600182015485169281019290925260028101549093169181019190915260038201805491929160608401919062002ca89062007883565b80601f016020809104026020016040519081016040528092919081815260200182805462002cd69062007883565b801562002d275780601f1062002cfb5761010080835404028352916020019162002d27565b820191906000526020600020905b81548152906001019060200180831162002d0957829003601f168201915b505050918352505060048201546020820152600582015460408201526006820154606082015260078201546001600160a01b03166080820152600882015460a0820152600982015460c0820152600a82015460e09091019060ff16600181111562002d965762002d9662006f46565b600181111562002daa5762002daa62006f46565b8152505090508086848151811062002dc65762002dc6620078f0565b6020026020010181905250828062002dde9062007906565b935050505b8062002def8162007906565b91505062002b61565b50600f548451111562002e785762002e108462005b27565b60005b600f5481101562002e715784818151811062002e335762002e33620078f0565b602002602001015184828151811062002e505762002e50620078f0565b6020026020010181905250808062002e689062007906565b91505062002e13565b5062002e7c565b8392505b62002e8a601e600062006b8f565b62002e986020600062006bb2565b6000601d8190555b8581101562003089576000604051806040016040528086848151811062002ecb5762002ecb620078f0565b6020026020010151602001516001600160a01b0316815260200186848151811062002efa5762002efa620078f0565b60209081029190910181015160a00151909152601e805460018101825560009190915282517f50bb669a95c7b50b7e8a6f09454034b2b14cf2b85c730dca9a539ca82cb6e350600290920291820180546001600160a01b0319166001600160a01b03909216919091179055828201517f50bb669a95c7b50b7e8a6f09454034b2b14cf2b85c730dca9a539ca82cb6e3519091015586519192509086908490811062002fa95762002fa9620078f0565b602090810291909101810151606001518254600181018455600093845292829020815162002fe1949190910192919091019062006bd2565b5084828151811062002ff75762002ff7620078f0565b602002602001015160400151848381518110620030185762003018620078f0565b60200260200101906001600160a01b031690816001600160a01b0316815250508482815181106200304d576200304d620078f0565b602002602001015160a00151601d60008282546200306c919062007924565b909155508291506200308090508162007906565b91505062002ea0565b50909550505050505090565b6008546001600160a01b03163314620030c25760405162461bcd60e51b815260040162000e149062007815565b620030cf60008362005b44565b6200180660018262005b44565b6001600160a01b03808216600081815260296020526040902060010154909116146200311c5760405162461bcd60e51b815260040162000e149062007792565b6001600160a01b03818116600090815260296020526040902054163314620031585760405162461bcd60e51b815260040162000e1490620077c9565b60016001600160a01b0382166000908152602960205260409020600a015460ff1660018111156200318d576200318d62006f46565b14620031dc5760405162461bcd60e51b815260206004820152601860248201527f76616c696461746f72206d757374206265207061757365640000000000000000604482015260640162000e14565b6001600160a01b03166000908152602960205260409020600a01805460ff19169055565b60608060006001818054620032159062007883565b80601f0160208091040260200160405190810160405280929190818152602001828054620032439062007883565b8015620032945780601f10620032685761010080835404028352916020019162003294565b820191906000526020600020905b8154815290600101906020018083116200327657829003601f168201915b50505050509150808054620032a99062007883565b80601f0160208091040260200160405190810160405280929190818152602001828054620032d79062007883565b8015620033285780601f10620032fc5761010080835404028352916020019162003328565b820191906000526020600020905b8154815290600101906020018083116200330a57829003601f168201915b50505050509050915091509091565b606060128054806020026020016040519081016040528092919081815260200182805480156200339157602002820191906000526020600020905b81546001600160a01b0316815260019091019060200180831162003372575b5050505050905090565b6008546001600160a01b03163314620033c85760405162461bcd60e51b815260040162000e149062007815565b600b8190556040518181527f1f4d2fc7529047a5bd96d3229bfea127fd18b7748f13586e097c69fccd3891289060200162001da4565b6008546001600160a01b031633146200342b5760405162461bcd60e51b815260040162000e149062007815565b6200343860008062006c61565b620034466001600062006c61565b6002805460ff19169055565b6008546001600160a01b031633146200347f5760405162461bcd60e51b815260040162000e149062007815565b600980546001600160a01b0319166001600160a01b0392909216919091179055565b60606000620034b18484620078d6565b6001600160401b03811115620034cb57620034cb6200718f565b6040519080825280602002602001820160405280156200350857816020015b620034f462006b55565b815260200190600190039081620034ea5790505b50905060005b6200351a8585620078d6565b811015620012d1576022600062003532838862007924565b81526020808201929092526040908101600020815160808101835281546001600160a01b039081168252600183015416938101939093526002810154918301919091526003015460608201528251839083908110620035955762003595620078f0565b60200260200101819052508080620035ad9062007906565b9150506200350e565b600085815260136020908152604080832060ff8089168552908352818420878216855283528184206001600160a01b0387168552835281842090851684529091529020600201546060908190871415620036ea57600087815260136020908152604080832060ff808b168552908352818420898216855283528184206001600160a01b038916855283528184209087168452909152902060030180546200365d9062007883565b80601f01602080910402602001604051908101604052809291908181526020018280546200368b9062007883565b8015620036dc5780601f10620036b057610100808354040283529160200191620036dc565b820191906000526020600020905b815481529060010190602001808311620036be57829003601f168201915b5050505050915050620036ed565b90505b95945050505050565b6001600160a01b0380821660008181526029602052604081206001015490921614620037365760405162461bcd60e51b815260040162000e149062007792565b506001600160a01b03166000908152601a602052604090205490565b6000806200375f62006ca0565b600060408286516020880160ff5afa6200377857600080fd5b508051602090910151600160601b90910494909350915050565b606081620037b75750506040805180820190915260018152600360fc1b602082015290565b8160005b8115620037e75780620037ce8162007906565b9150620037df9050600a83620079e1565b9150620037bb565b6000816001600160401b038111156200380457620038046200718f565b6040519080825280601f01601f1916602001820160405280156200382f576020820181803683370190505b5090505b8415620038a75762003847600183620078d6565b915062003856600a86620079f8565b6200386390603062007924565b60f81b8183815181106200387b576200387b620078f0565b60200101906001600160f81b031916908160001a9053506200389f600a86620079e1565b945062003833565b949350505050565b6001600160a01b038316620039135760405162461bcd60e51b8152602060048201526024808201527f45524332303a20617070726f76652066726f6d20746865207a65726f206164646044820152637265737360e01b606482015260840162000e14565b6001600160a01b038216620039765760405162461bcd60e51b815260206004820152602260248201527f45524332303a20617070726f766520746f20746865207a65726f206164647265604482015261737360f01b606482015260840162000e14565b6001600160a01b0383811660008181526021602090815260408083209487168084529482529182902085905590518481527f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b92591015b60405180910390a3505050565b6001600160a01b038116600090815260296020526040812090600a82015460ff16600181111562003a0d5762003a0d62006f46565b1462003a5c5760405162461bcd60e51b815260206004820152601960248201527f76616c696461746f72206d75737420626520656e61626c656400000000000000604482015260640162000e14565b600a8101805460ff191660011790558054600d54601c547f75bdcdbe540758778e669d108fbcb7ede734f27f46e4e5525eeb8ecf91849a9c926001600160a01b031691859162003aad919062007924565b604080516001600160a01b03948516815293909216602084015290820152606001620018bc565b6001600160a01b03831660009081526028602052604090205481111562003b375760405162461bcd60e51b8152602060048201526016602482015275616d6f756e7420657863656564732062616c616e636560501b604482015260640162000e14565b6001600160a01b0383166000908152602860205260408120805483929062003b61908490620078d6565b90915550506001600160a01b0382166000908152602860205260408120805483929062003b9090849062007924565b92505081905550816001600160a01b0316836001600160a01b03167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef83604051620039cb91815260200190565b6000816040015160ff16600281111562003bfb5762003bfb62006f46565b600281111562003c0f5762003c0f62006f46565b14801562003c32575060c081015160009081526016602052604090205415156001145b1562003c3b5750565b6001816040015160ff16600281111562003c595762003c5962006f46565b600281111562003c6d5762003c6d62006f46565b14801562003c90575060c081015160009081526017602052604090205415156001145b1562003c995750565b60c0810180516000908152601360209081526040808320818601805160ff9081168652918452828520606088018051841687529085528386206080890180516001600160a01b039081168952918752858820878b01805187168a52908852959097208951815496519451935198518316600160201b02600160201b600160c01b031999871663010000000263ff0000001995881662010000029590951663ffff0000199688166101000261ffff199099169290971691909117969096179390931693909317179490941691909117835560a0850151600184018054919092166001600160a01b03199091161790559151600282015560e08301518051849362003daa92600385019291019062006bd2565b5090505060005b816000015160ff168160ff16101562003e445760c082015160009081526013602090815260408083208186015160ff9081168552908352818420606087015182168552835281842060808701516001600160a01b0390811686529084528285209186168552925290912054600160201b900416331462003e2f575050565b8062003e3b8162007a0f565b91505062003db1565b5062003e4f62006b0f565b6080808301516001600160a01b0390811691830191825260c080850151908401908152845160ff9081168552606080870151821690860190815260208088015183168188019081526040808a0151851690890190815260a0808b01518816908a01908152601880546001810182556000919091528a5160049091027fb13d2d76d1f4b7be834882e410b3e3a8afaf69f83600ae24db354391d2378d2e810180549551945197519b518b16600160201b02600160201b600160c01b03199c8a1663010000000263ff00000019998b1662010000029990991663ffff000019968b166101000261ffff1990981694909a1693909317959095179390931696909617949094179790971693909317835590517fb13d2d76d1f4b7be834882e410b3e3a8afaf69f83600ae24db354391d2378d2f86018054919095166001600160a01b03199091161790935590517fb13d2d76d1f4b7be834882e410b3e3a8afaf69f83600ae24db354391d2378d3084015560e084015180518594929362003ff8937fb13d2d76d1f4b7be834882e410b3e3a8afaf69f83600ae24db354391d2378d310192019062006bd2565b5060009150620040059050565b826040015160ff16600281111562004021576200402162006f46565b600281111562004035576200403562006f46565b1415620040535760c082015160009081526016602052604090204390555b6001826040015160ff16600281111562004071576200407162006f46565b600281111562004085576200408562006f46565b141562001806575060c001516000908152601760205260409020439055565b600080600080620040bb60fe8660e0015162005c97565b93509350935093508460c0015183141580620040ed57508460a001516001600160a01b0316846001600160a01b031614155b80620040f7575081155b806200410a5750846060015160ff168114155b1562004117575050505050565b845160ff16156200412a57606060e08601525b60a08501516001600160a01b03908116600081815260296020526040902060010154909116141562001b655760a0850180516001600160a01b039081166000908152601460209081526040808320805460018082018355918552938390208b516004909502018054848d0151938d015160608e015160808f015160ff98891661ffff1990941693909317610100968916969096029590951763ffff00001916620100009188169190910263ff0000001916176301000000969094169590950292909217600160201b600160c01b031916600160201b94861694909402939093178155935191840180546001600160a01b031916929093169190911790915560c0870151600283015560e087015180518893926200424f92600385019291019062006bd2565b50506019805460018101825560009190915286517f944998273e477b495144fb8794c914197f3ccb46be2900f4698fd0ef743c9695600490920291820180546020808b015160408c015160608d015160808e015160ff97881661ffff1990961695909517610100938816939093029290921763ffff00001916620100009187169190910263ff0000001916176301000000959091169490940293909317600160201b600160c01b031916600160201b6001600160a01b039283160217825560a08a01517f944998273e477b495144fb8794c914197f3ccb46be2900f4698fd0ef743c9696850180546001600160a01b0319169190921617905560c08901517f944998273e477b495144fb8794c914197f3ccb46be2900f4698fd0ef743c969784015560e089015180518a95509193620043b0937f944998273e477b495144fb8794c914197f3ccb46be2900f4698fd0ef743c96989091019291019062006bd2565b50505060c085015160009081526016602052604090819020439055517fe9b2e40b11e32b8729ed1bfd4c1ae17d2bcdc9af959564da14b39ca570607e3f90620043fb90879062007a32565b60405180910390a15050505050565b6000806000806200442160fc8660e0015162005c97565b93509350935093508460c00151831415806200445357508460a001516001600160a01b0316846001600160a01b031614155b806200445d575081155b80620044705750846060015160ff168114155b156200447d575050505050565b845160ff16156200449057606060e08601525b60a0850180516001600160a01b039081166000908152601560209081526040808320805460018082018355918552938390208b516004909502018054848d0151938d015160608e015160808f015160ff98891661ffff1990941693909317610100968916969096029590951763ffff00001916620100009188169190910263ff0000001916176301000000969094169590950292909217600160201b600160c01b031916600160201b94861694909402939093178155935191840180546001600160a01b031916929093169190911790915560c0870151600283015560e087015180518893926200458992600385019291019062006bd2565b50505060c085015160009081526017602052604090819020439055517f244ffefead78aaef5913a3abac1c8477dec686bf017a343f1679f9c8b6a77f1190620043fb90879062007a32565b600080600080620045eb60fd8660e0015162005c97565b93509350935093508460c00151831415806200461d57508460a001516001600160a01b0316846001600160a01b031614155b8062004627575081155b806200463a5750846060015160ff168114155b1562004647575050505050565b845160ff16156200465a57606060e08601525b620046658562005cf0565b60c085015160009081526017602052604080822091909155517f663327acde77befae0ec3fed52a32b993673702ddd832552e402d1afbd32158c90620043fb90879062007a32565b60005b60185481101562004a4d57620046c562006b0f565b60188281548110620046db57620046db620078f0565b6000918252602090912060049091020154600160201b90046001600160a01b031660808201526018805483908110620047185762004718620078f0565b9060005260206000209060040201600201548160c001818152505060188281548110620047495762004749620078f0565b600091825260209091206004909102015460ff1681526018805483908110620047765762004776620078f0565b60009182526020909120600490910201546301000000900460ff1660608201526018805483908110620047ad57620047ad620078f0565b60009182526020918290206004909102015460ff61010090910416908201526018805483908110620047e357620047e3620078f0565b600091825260209091206004909102015462010000900460ff1660408201526018805483908110620048195762004819620078f0565b600091825260208220600160049092020101546001600160a01b031660a08301525b816000015160ff168160ff161015620049625760e082015160c083015160009081526013602090815260408083208187015160ff9081168552908352818420606088015182168552835281842060808801516001600160a01b03168552835281842090861684529091529020600301805462004948929190620048be9062007883565b80601f0160208091040260200160405190810160405280929190818152602001828054620048ec9062007883565b80156200493d5780601f1062004911576101008083540402835291602001916200493d565b820191906000526020600020905b8154815290600101906020018083116200491f57829003601f168201915b505050505062005f82565b60e083015280620049598162007a0f565b9150506200483b565b506000816040015160ff16600281111562004981576200498162006f46565b600281111562004995576200499562006f46565b1415620049ae57620049a781620040a4565b5062004a38565b6001816040015160ff166002811115620049cc57620049cc62006f46565b6002811115620049e057620049e062006f46565b1415620049f257620049a7816200440a565b6002816040015160ff16600281111562004a105762004a1062006f46565b600281111562004a245762004a2462006f46565b141562004a3657620049a781620045d4565b505b8062004a448162007906565b915050620046b0565b5062004a5c6018600062006cbe565b565b60005b601e5481101562000e6457600060156000601e848154811062004a885762004a88620078f0565b600091825260208083206002909202909101546001600160a01b031683528201929092526040018120549150816001600160401b0381111562004acf5762004acf6200718f565b60405190808252806020026020018201604052801562004b0c57816020015b62004af862006b0f565b81526020019060019003908162004aee5790505b5090506000805b8381101562004fee57600060156000601e888154811062004b385762004b38620078f0565b600091825260208083206002909202909101546001600160a01b03168352820192909252604001902080548390811062004b765762004b76620078f0565b60009182526020918290206040805161010080820183526004909402909201805460ff80821685529481048516958401959095526201000085048416918301919091526301000000840490921660608201526001600160a01b03600160201b90930483166080820152600182015490921660a0830152600281015460c083015260038101805460e08401919062004c0d9062007883565b80601f016020809104026020016040519081016040528092919081815260200182805462004c3b9062007883565b801562004c8c5780601f1062004c605761010080835404028352916020019162004c8c565b820191906000526020600020905b81548152906001019060200180831162004c6e57829003601f168201915b5050509190925250505060c08101516000908152601760205260408120549192509062004cba9043620078d6565b60c08301516000908152601760205260409020549091501580159062004ce05750603c81115b1562004fd65760146000601e898154811062004d005762004d00620078f0565b60009182526020808320600292830201546001600160a01b0390811685528482019590955260409384018320805460018082018355918552938290208851600490950201805489840151968a015160608b015160808c015160ff98891661ffff1990941693909317610100998916999099029890981763ffff00001916620100009188169190910263ff0000001916176301000000969097169590950295909517600160201b600160c01b031916600160201b9487169490940293909317845560a087015192840180546001600160a01b031916939095169290921790935560c08501519282019290925560e08401518051859362004e0792600385019291019062006bd2565b50506019805460018101825560009190915283517f944998273e477b495144fb8794c914197f3ccb46be2900f4698fd0ef743c969560049092029182018054602080880151604089015160608a015160808b015160ff97881661ffff1990961695909517610100938816939093029290921763ffff00001916620100009187169190910263ff0000001916176301000000959091169490940293909317600160201b600160c01b031916600160201b6001600160a01b039283160217825560a08701517f944998273e477b495144fb8794c914197f3ccb46be2900f4698fd0ef743c9696850180546001600160a01b0319169190921617905560c08601517f944998273e477b495144fb8794c914197f3ccb46be2900f4698fd0ef743c969784015560e08601518051879550919362004f68937f944998273e477b495144fb8794c914197f3ccb46be2900f4698fd0ef743c96989091019291019062006bd2565b5050507f550841d9b29e92358159fe7d9bda9bfef1f1bc478fd9160241cff35168cfc7128260405162004f9c919062007a32565b60405180910390a18185858151811062004fba5762004fba620078f0565b6020026020010181905250838062004fd29062007906565b9450505b5050808062004fe59062007906565b91505062004b13565b5060005b81811015620050375762005022838281518110620050145762005014620078f0565b602002602001015162005cf0565b806200502e8162007906565b91505062004ff2565b505050508080620050489062007906565b91505062004a61565b60005b601954811015620050b757600060198281548110620050775762005077620078f0565b60009182526020909120600160049092020101546001600160a01b03169050620050a18162006003565b5080620050ae8162007906565b91505062005054565b5062004a5c6019600062006cbe565b47620050ce57565b600a544790600090670de0b6b3a764000090620050ed90849062007a47565b620050f99190620079e1565b905080156200517e5760095460405160009182916001600160a01b039091169084908381818185875af1925050503d806000811462005155576040519150601f19603f3d011682016040523d82523d6000602084013e6200515a565b606091505b509092509050600182151514156200517b57620051788385620078d6565b93505b50505b81601f600082825462005192919062007924565b909155506200180690508262006178565b6023545b602454811015620051d257620051bd816200638a565b80620051c98162007906565b915050620051a7565b50602454602355602654805b6027548110156200524c57600e5460008281526025602052604090206003015443916200520b9162007924565b1162005231576200521c81620064a3565b6200522960018362007924565b915062005237565b6200524c565b80620052438162007906565b915050620051de565b50602655565b600754600654101562004a5c576006546000908152600560205260409020600e5460018201544391620052859162007924565b11156200528f5750565b600281015481546001600160a01b039081166000908152602960205260408082206004908101859055855484168352918190206007015490516319fac8fd60e01b81529216926319fac8fd92620052ea920190815260200190565b600060405180830381600087803b1580156200530557600080fd5b505af11580156200531a573d6000803e3d6000fd5b505060068054600090815260056020526040812080546001600160a01b0319168155600180820183905560029091018290558254909450919250906200536290849062007924565b9091555062005252915050565b60008211620053cd5760405162461bcd60e51b815260206004820152602360248201527f616d6f756e74206e65656420746f206265207374726963746c7920706f73697460448201526269766560e81b606482015260840162000e14565b6001600160a01b038116600090815260286020526040902054821115620054375760405162461bcd60e51b815260206004820152601b60248201527f696e73756666696369656e74204e6577746f6e2062616c616e63650000000000604482015260640162000e14565b6001600160a01b0381166000908152602860205260408120805484929062005461908490620078d6565b9091555050604080516080810182526001600160a01b038084168252858116602080840191825283850187815243606086019081526024805460009081526022909452968320865181549087166001600160a01b031991821617825594516001820180549190971695169490941790945551600283015591516003909101558254919290620054f08362007906565b919050555050505050565b6001600160a01b038381166000908152602960205260408082206007015490516370a0823160e01b81528484166004820152919216906370a0823190602401602060405180830381865afa15801562005558573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200557e919062007a69565b905082811015620055dd5760405162461bcd60e51b815260206004820152602260248201527f696e73756666696369656e74204c6971756964204e6577746f6e2062616c616e604482015261636560f01b606482015260840162000e14565b6001600160a01b0380851660009081526029602090815260408083206007015481516318160ddd60e01b81529151939416926318160ddd926004808401939192918290030181865afa15801562005638573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200565e919062007a69565b90506200566b856200655e565b80156200567757508084145b15620056d75760405162461bcd60e51b815260206004820152602860248201527f63616e2774206861766520636f6d6d6974746565206d656d626572207769746860448201526737baba1026272a2760c11b606482015260840162000e14565b6001600160a01b0385811660009081526029602052604090819020600701549051632770a7eb60e21b8152858316600482015260248101879052911690639dc29fac90604401600060405180830381600087803b1580156200573857600080fd5b505af11580156200574d573d6000803e3d6000fd5b5050604080516080810182526001600160a01b03808816825289811660208084019182528385018b815243606086019081526027805460009081526025909452968320865181549087166001600160a01b031991821617825594516001820180549190971695169490941790945551600283015591516003909101558254919450909250620057dc8362007906565b9190505550505050505050565b8051608214620058335760405162461bcd60e51b8152602060048201526014602482015273092dcecc2d8d2c840e0e4dedecc40d8cadccee8d60631b604482015260640162000e14565b6200583e82620065cf565b604080518082018252601a81527f19457468657265756d205369676e6564204d6573736167653a0a000000000000602080830191909152845192519192600092620058a1920160609190911b6bffffffffffffffffffffffff1916815260140190565b6040516020818303038152906040529050600082620058c1835162003792565b83604051602001620058d69392919062007a83565b60408051601f1981840301815282825280516020918201206002808552606085018452909450600093929091830190803683370190505090506000808060205b8851811015620059e2576200592c8982620066fc565b6040805160008152602081018083528b905260ff8316918101919091526060810184905260808101839052929650909450925060019060a0016020604051602081039080840390855afa15801562005988573d6000803e3d6000fd5b5050604051601f190151905085620059a2604184620079e1565b81518110620059b557620059b5620078f0565b6001600160a01b0390921660209283029190910190910152620059da60418262007924565b905062005916565b5088602001516001600160a01b03168460008151811062005a075762005a07620078f0565b60200260200101516001600160a01b03161462005a795760405162461bcd60e51b815260206004820152602960248201527f496e76616c6964206e6f6465206b6579206f776e6572736869702070726f6f66604482015268081c1c9bdd9a59195960ba1b606482015260840162000e14565b88604001516001600160a01b03168460018151811062005a9d5762005a9d620078f0565b60200260200101516001600160a01b03161462005b115760405162461bcd60e51b815260206004820152602b60248201527f496e76616c6964206f7261636c65206b6579206f776e6572736869702070726f60448201526a1bd9881c1c9bdd9a59195960aa1b606482015260840162000e14565b62005b1c8962006733565b505050505050505050565b62000e648160006001845162005b3e9190620078d6565b620068f7565b81546002600180831615610100020382160482518082016020811060208410016002811462005bf3576001811462005c19578660005260208404602060002001600160028402018855602085068060200390508088018589016001836101000a0392508282511684540184556001840193506020820191505b8082101562005bdc578151845560018401935060208201915062005bbd565b815191036101000a90819004029091555062005c8e565b60028302826020036101000a846020036101000a60208901510402018501875562005c8e565b8660005260208404602060002001600160028402018855846020038088018589016001836101000a0392508282511660ff198a160184556020820191506001840193505b8082101562005c7c578151845560018401935060208201915062005c5d565b815191036101000a9081900402909155505b50505050505050565b60008060008060008551602062005caf919062007924565b905062005cbb62006ce1565b60808183898b5afa62005ccd57600080fd5b805160208201516040830151606090930151919a90995091975095509350505050565b60a08101516001600160a01b0316600090815260156020526040812054905b8181101562005f7d5760c083015160a08401516001600160a01b0316600090815260156020526040902080548390811062005d4e5762005d4e620078f0565b906000526020600020906004020160020154141562005f685760a08301516001600160a01b0316600090815260156020526040902062005d90600184620078d6565b8154811062005da35762005da3620078f0565b9060005260206000209060040201601560008560a001516001600160a01b03166001600160a01b03168152602001908152602001600020828154811062005dee5762005dee620078f0565b600091825260209091208254600490920201805460ff19811660ff938416908117835584546101009081900485160261ffff1990921617178082558354620100009081900484160262ff000019821681178355845463010000009081900490941690930263ff0000001990931663ffff000019909116179190911780825582546001600160a01b03600160201b918290048116909102600160201b600160c01b031990921691909117825560018084015490830180546001600160a01b0319169190921617905560028083015490820155600380830180549183019162005ed59062007883565b62005ee292919062006cff565b50505060a08301516001600160a01b0316600090815260156020526040902080548062005f135762005f1362007acc565b60008281526020812060046000199093019283020180546001600160c01b03191681556001810180546001600160a01b0319169055600281018290559062005f5f600383018262006c61565b50509055505050565b8062005f748162007906565b91505062005d0f565b505050565b6060806040519050835180825260208201818101602087015b8183101562005fb557805183526020928301920162005f9b565b50855184518101855292509050808201602086015b8183101562005fe457805183526020928301920162005fca565b508651929092011591909101601f01601f191660405250905092915050565b6003546001600160a01b0382166000908152602960205260409020600501546001106200602e575050565b6001600160a01b03821660009081526029602052604090206005015481106200607e576001600160a01b0382166000908152602960205260409020600501546200607b90600190620078d6565b90505b6001600160a01b03821660009081526029602052604081206005018054839290620060ab908490620078d6565b90915550503060009081526028602052604081208054839290620060d190849062007924565b90915550506001600160a01b03821660009081526029602052604081206006018054600192906200610490849062007924565b90915550506001600160a01b0382166000908152601a6020526040812080548392906200613390849062007924565b9091555050604080516001600160a01b0384168152602081018390527f51cf713376ddb1e5f5828bb6aa39d99de812176d62c3d3550bdc4e0b5e86e1a59101620018bc565b60008111620061845750565b6000805b601e548110156200620557600060296000601e8481548110620061af57620061af620078f0565b600091825260208083206002909202909101546001600160a01b0316835282019290925260400190206005810154909150620061ec908462007924565b9250508080620061fc9062007906565b91505062006188565b508062006210575050565b6000805b601e548110156200638457600060296000601e84815481106200623b576200623b620078f0565b600091825260208083206002909202909101546001600160a01b031683528201929092526040018120600581015490925085906200627b90889062007a47565b620062879190620079e1565b905080156200636c5760008260070160009054906101000a90046001600160a01b03166001600160a01b031663fb489a7b836040518263ffffffff1660e01b815260040160206040518083038185885af1158015620062ea573d6000803e3d6000fd5b50505050506040513d601f19601f8201168201806040525081019062006311919062007a69565b90506200631f818662007924565b6001840154604080516001600160a01b039092168252602082018490529196507fb3b7a071186534c03b40695710096f289fd4ed6c1a374aff0bb648955e4fe563910160405180910390a1505b505080806200637b9062007906565b91505062006214565b50505050565b600081815260226020908152604080832060018101546001600160a01b0316845260299092528220600581015491929091620063cc57506002820154620063f6565b816005015483600201548360080154620063e7919062007a47565b620063f39190620079e1565b90505b600782015483546040516340c10f1960e01b81526001600160a01b039182166004820152602481018490529116906340c10f1990604401600060405180830381600087803b1580156200644857600080fd5b505af11580156200645d573d6000803e3d6000fd5b5050505082600201548260050160008282546200647b919062007924565b925050819055508082600801600082825462006498919062007924565b909155505050505050565b600081815260256020908152604080832060018101546001600160a01b031684526029909252822060088101546005820154600284015493949293620064ea919062007a47565b620064f69190620079e1565b9050808260050160008282546200650e9190620078d6565b909155505060028301546008830180546000906200652e908490620078d6565b909155505082546001600160a01b0316600090815260286020526040812080548392906200649890849062007924565b6000805b601e54811015620065c657601e8181548110620065835762006583620078f0565b60009182526020909120600290910201546001600160a01b0384811691161415620065b15750600192915050565b80620065bd8162007906565b91505062006562565b50600092915050565b6000620065e0826060015162003752565b6001600160a01b03909116602084015290508015620066305760405162461bcd60e51b815260206004820152600b60248201526a32b737b2329032b93937b960a91b604482015260640162000e14565b6020808301516001600160a01b03908116600090815260299092526040909120600101541615620066a45760405162461bcd60e51b815260206004820152601c60248201527f76616c696461746f7220616c7265616479207265676973746572656400000000604482015260640162000e14565b61271082608001511115620018065760405162461bcd60e51b815260206004820152601760248201527f696e76616c696420636f6d6d697373696f6e2072617465000000000000000000604482015260640162000e14565b8181018051602082015160409092015190919060001a601b8110156200672c5762006729601b8262007ae2565b90505b9250925092565b60e08101516001600160a01b0316620067b457601254600090620067579062003792565b905081602001518260000151836080015183604051620067779062006d83565b62006786949392919062007b0a565b604051809103906000f080158015620067a3573d6000803e3d6000fd5b506001600160a01b031660e0830152505b60208082018051601280546001818101835560009283527fbb8a6a4669ba250d26cd7a459eca9d215f8307e33aebe50379bc5a3617ec344490910180546001600160a01b039485166001600160a01b03199182161790915584518416835260298652604092839020875181549086169083161781559451918501805492851692821692909217909155908501516002840180549190931691161790556060830151805184936200686c92600385019291019062006bd2565b506080820151600482015560a0820151600582015560c0820151600682015560e08201516007820180546001600160a01b0319166001600160a01b0390921691909117905561010082015160088201556101208201516009820155610140820151600a8201805460ff191660018381811115620068ed57620068ed62006f46565b0217905550505050565b81818082141562006909575050505050565b60008560026200691a878762007b49565b62006926919062007b8e565b62006932908762007bc2565b81518110620069455762006945620078f0565b602002602001015160a0015190505b81831362006a82575b80868481518110620069735762006973620078f0565b602002602001015160a0015111156200699b5782620069928162007c09565b9350506200695d565b858281518110620069b057620069b0620078f0565b602002602001015160a00151811115620069d95781620069d08162007c25565b9250506200699b565b81831362006a7c57858281518110620069f657620069f6620078f0565b602002602001015186848151811062006a135762006a13620078f0565b602002602001015187858151811062006a305762006a30620078f0565b6020026020010188858151811062006a4c5762006a4c620078f0565b602002602001018290528290525050828062006a689062007c09565b935050818062006a789062007c25565b9250505b62006954565b8185121562006a985762006a98868684620068f7565b8383121562006aae5762006aae868486620068f7565b505050505050565b60408051610160810182526000808252602082018190529181018290526060808201526080810182905260a0810182905260c0810182905260e08101829052610100810182905261012081018290529061014082015290565b604080516101008101825260008082526020820181905291810182905260608082018390526080820183905260a0820183905260c082019290925260e081019190915290565b604051806080016040528060006001600160a01b0316815260200160006001600160a01b0316815260200160008152602001600081525090565b508054600082556002029060005260206000209081019062000e64919062006d91565b508054600082559060005260206000209081019062000e64919062006db9565b82805462006be09062007883565b90600052602060002090601f01602090048101928262006c04576000855562006c4f565b82601f1062006c1f57805160ff191683800117855562006c4f565b8280016001018555821562006c4f579182015b8281111562006c4f57825182559160200191906001019062006c32565b5062006c5d92915062006dda565b5090565b50805462006c6f9062007883565b6000825580601f1062006c80575050565b601f01602090049060005260206000209081019062000e64919062006dda565b60405180604001604052806002906020820280368337509192915050565b508054600082556004029060005260206000209081019062000e64919062006df1565b60405180608001604052806004906020820280368337509192915050565b82805462006d0d9062007883565b90600052602060002090601f01602090048101928262006d31576000855562006c4f565b82601f1062006d44578054855562006c4f565b8280016001018555821562006c4f57600052602060002091601f016020900482015b8281111562006c4f57825482559160010191906001019062006d66565b61116a8062007c4783390190565b5b8082111562006c5d5780546001600160a01b03191681556000600182015560020162006d92565b8082111562006c5d57600062006dd0828262006c61565b5060010162006db9565b5b8082111562006c5d576000815560010162006ddb565b8082111562006c5d5780546001600160c01b03191681556001810180546001600160a01b031916905560006002820181905562006e32600383018262006c61565b5060040162006df1565b60005b8381101562006e5957818101518382015260200162006e3f565b83811115620063845750506000910152565b6000815180845262006e8581602086016020860162006e3c565b601f01601f19169290920160200192915050565b60208152600062006eae602083018462006e6b565b9392505050565b6001600160a01b038116811462000e6457600080fd5b803562006ed88162006eb5565b919050565b6000806040838503121562006ef157600080fd5b823562006efe8162006eb5565b946020939093013593505050565b60006020828403121562006f1f57600080fd5b813562006eae8162006eb5565b60006020828403121562006f3f57600080fd5b5035919050565b634e487b7160e01b600052602160045260246000fd5b6002811062006f7b57634e487b7160e01b600052602160045260246000fd5b9052565b6020815262006f9a6020820183516001600160a01b03169052565b6000602083015162006fb760408401826001600160a01b03169052565b5060408301516001600160a01b038116606084015250606083015161016080608085015262006feb61018085018362006e6b565b9150608085015160a085015260a085015160c085015260c085015160e085015260e085015161010062007028818701836001600160a01b03169052565b86015161012086810191909152860151610140808701919091528601519050620070558286018262006f5c565b5090949350505050565b600061010060ff835116845260ff602084015116602085015260ff604084015116604085015260ff606084015116606085015260018060a01b03608084015116608085015260a0830151620070bf60a08601826001600160a01b03169052565b5060c083015160c085015260e08301518160e0860152620036ed8286018262006e6b565b6000602080830181845280855180835260408601915060408160051b870101925083870160005b828110156200713c57603f19888603018452620071298583516200705f565b945092850192908501906001016200710a565b5092979650505050505050565b6000806000606084860312156200715f57600080fd5b83356200716c8162006eb5565b925060208401356200717e8162006eb5565b929592945050506040919091013590565b634e487b7160e01b600052604160045260246000fd5b60405161010081016001600160401b0381118282101715620071cb57620071cb6200718f565b60405290565b604051601f8201601f191681016001600160401b0381118282101715620071fc57620071fc6200718f565b604052919050565b803560ff8116811462006ed857600080fd5b600082601f8301126200722857600080fd5b81356001600160401b038111156200724457620072446200718f565b62007259601f8201601f1916602001620071d1565b8181528460208386010111156200726f57600080fd5b816020850160208301376000918101602001919091529392505050565b60006020808385031215620072a057600080fd5b82356001600160401b0380821115620072b857600080fd5b818501915085601f830112620072cd57600080fd5b813581811115620072e257620072e26200718f565b8060051b620072f3858201620071d1565b91825283810185019185810190898411156200730e57600080fd5b86860192505b838310156200740b578235858111156200732e5760008081fd5b8601610100818c03601f1901811315620073485760008081fd5b62007352620071a5565b6200735f8a840162007204565b815260406200737081850162007204565b8b83015260606200738381860162007204565b82840152608091506200739882860162007204565b9083015260a0620073ab85820162006ecb565b8284015260c09150620073c082860162006ecb565b9083015260e08481013582840152928401359289841115620073e457600091508182fd5b620073f48f8d8688010162007216565b908301525084525050918601919086019062007314565b9998505050505050505050565b600081518084526020808501945080840160005b838110156200745e57815180516001600160a01b0316885283015183880152604090960195908201906001016200742c565b509495945050505050565b8215158152604060208201526000620038a7604083018462007418565b600080604083850312156200749a57600080fd5b50508035926020909101359150565b602080825282518282018190526000919060409081850190868401855b828110156200750e57815180516001600160a01b03908116865287820151168786015285810151868601526060908101519085015260809093019290850190600101620074c6565b5091979650505050505050565b6000602080830181845280855180835260408601915060408160051b870101925083870160005b828110156200713c57603f198886030184526200756185835162006e6b565b9450928501929085019060010162007542565b60208152600062006eae602083018462007418565b6000806000606084860312156200759f57600080fd5b83356001600160401b0380821115620075b757600080fd5b620075c58783880162007216565b945060208601359150620075d98262006eb5565b90925060408501359080821115620075f057600080fd5b50620075ff8682870162007216565b9150509250925092565b6020808252825182820181905260009190848201906040850190845b818110156200764c5783516001600160a01b03168352928401929184019160010162007625565b50909695505050505050565b600080604083850312156200766c57600080fd5b82356001600160401b03808211156200768457600080fd5b620076928683870162007216565b93506020850135915080821115620076a957600080fd5b50620076b88582860162007216565b9150509250929050565b604081526000620076d7604083018562006e6b565b8281036020840152620036ed818562006e6b565b60008060408385031215620076ff57600080fd5b82356200770c8162006eb5565b915060208301356200771e8162006eb5565b809150509250929050565b600080600080600060a086880312156200774257600080fd5b85359450620077546020870162007204565b9350620077646040870162007204565b92506060860135620077768162006eb5565b9150620077866080870162007204565b90509295509295909350565b6020808252601c908201527f76616c696461746f72206d757374206265207265676973746572656400000000604082015260600190565b6020808252602c908201527f726571756972652063616c6c657220746f2062652076616c696461746f72206160408201526b191b5a5b881858d8dbdd5b9d60a21b606082015260800190565b6020808252601a908201527f63616c6c6572206973206e6f7420746865206f70657261746f72000000000000604082015260600190565b60208082526018908201527f76616c696461746f72206e6f7420726567697374657265640000000000000000604082015260600190565b600181811c908216806200789857607f821691505b60208210811415620078ba57634e487b7160e01b600052602260045260246000fd5b50919050565b634e487b7160e01b600052601160045260246000fd5b600082821015620078eb57620078eb620078c0565b500390565b634e487b7160e01b600052603260045260246000fd5b60006000198214156200791d576200791d620078c0565b5060010190565b600082198211156200793a576200793a620078c0565b500190565b60208082526023908201527f66756e6374696f6e207265737472696374656420746f207468652070726f746f60408201526218dbdb60ea1b606082015260800190565b600060018060a01b0380881683528087166020840152808616604084015260a06060840152620079b660a084018662006e6b565b91508084166080840152509695505050505050565b634e487b7160e01b600052601260045260246000fd5b600082620079f357620079f3620079cb565b500490565b60008262007a0a5762007a0a620079cb565b500690565b600060ff821660ff81141562007a295762007a29620078c0565b60010192915050565b60208152600062006eae60208301846200705f565b600081600019048311821515161562007a645762007a64620078c0565b500290565b60006020828403121562007a7c57600080fd5b5051919050565b6000845162007a9781846020890162006e3c565b84519083019062007aad81836020890162006e3c565b845191019062007ac281836020880162006e3c565b0195945050505050565b634e487b7160e01b600052603160045260246000fd5b600060ff821660ff84168060ff0382111562007b025762007b02620078c0565b019392505050565b6001600160a01b038581168252841660208201526040810183905260806060820181905260009062007b3f9083018462006e6b565b9695505050505050565b60008083128015600160ff1b85018412161562007b6a5762007b6a620078c0565b6001600160ff1b038401831381161562007b885762007b88620078c0565b50500390565b60008262007ba05762007ba0620079cb565b600160ff1b82146000198414161562007bbd5762007bbd620078c0565b500590565b600080821280156001600160ff1b038490038513161562007be75762007be7620078c0565b600160ff1b839003841281161562007c035762007c03620078c0565b50500190565b60006001600160ff1b038214156200791d576200791d620078c0565b6000600160ff1b82141562007c3e5762007c3e620078c0565b50600019019056fe60806040523480156200001157600080fd5b506040516200116a3803806200116a833981016040819052620000349162000212565b6127108211156200004457600080fd5b600980546001600160a01b038087166001600160a01b031992831617909255600a805492861692909116919091179055600b8290556040516200008c908290602001620002ff565b60405160208183030381529060405260079080519060200190620000b29291906200010a565b5080604051602001620000c69190620002ff565b60405160208183030381529060405260089080519060200190620000ec9291906200010a565b5050600080546001600160a01b03191633179055506200036b915050565b82805462000118906200032e565b90600052602060002090601f0160209004810192826200013c576000855562000187565b82601f106200015757805160ff191683800117855562000187565b8280016001018555821562000187579182015b82811115620001875782518255916020019190600101906200016a565b506200019592915062000199565b5090565b5b808211156200019557600081556001016200019a565b6001600160a01b0381168114620001c657600080fd5b50565b634e487b7160e01b600052604160045260246000fd5b60005b83811015620001fc578181015183820152602001620001e2565b838111156200020c576000848401525b50505050565b600080600080608085870312156200022957600080fd5b84516200023681620001b0565b60208601519094506200024981620001b0565b6040860151606087015191945092506001600160401b03808211156200026e57600080fd5b818701915087601f8301126200028357600080fd5b815181811115620002985762000298620001c9565b604051601f8201601f19908116603f01168101908382118183101715620002c357620002c3620001c9565b816040528281528a6020848701011115620002dd57600080fd5b620002f0836020830160208801620001df565b979a9699509497505050505050565b644c4e544e2d60d81b81526000825162000321816005850160208701620001df565b9190910160050192915050565b600181811c908216806200034357607f821691505b602082108114156200036557634e487b7160e01b600052602260045260246000fd5b50919050565b610def806200037b6000396000f3fe6080604052600436106100fe5760003560e01c8063372500ab1161009557806395d89b411161006457806395d89b41146102945780639dc29fac146102a9578063a9059cbb146102c9578063dd62ed3e146102e9578063fb489a7b1461032f57600080fd5b8063372500ab1461020957806340c10f191461021e57806370a082311461023e578063949813b81461027457600080fd5b806319fac8fd116100d157806319fac8fd1461019557806323b872dd146101b75780632f2c3f2e146101d7578063313ce567146101ed57600080fd5b806306fdde0314610103578063095ea7b31461012e57806318160ddd1461015e578063187cf4d71461017d575b600080fd5b34801561010f57600080fd5b50610118610337565b6040516101259190610b4b565b60405180910390f35b34801561013a57600080fd5b5061014e610149366004610bbc565b6103c9565b6040519015158152602001610125565b34801561016a57600080fd5b506003545b604051908152602001610125565b34801561018957600080fd5b5061016f633b9aca0081565b3480156101a157600080fd5b506101b56101b0366004610be6565b6103df565b005b3480156101c357600080fd5b5061014e6101d2366004610bff565b610417565b3480156101e357600080fd5b5061016f61271081565b3480156101f957600080fd5b5060405160128152602001610125565b34801561021557600080fd5b506101b561050a565b34801561022a57600080fd5b506101b5610239366004610bbc565b6105b8565b34801561024a57600080fd5b5061016f610259366004610c3b565b6001600160a01b031660009081526001602052604090205490565b34801561028057600080fd5b5061016f61028f366004610c3b565b610620565b3480156102a057600080fd5b50610118610654565b3480156102b557600080fd5b506101b56102c4366004610bbc565b610663565b3480156102d557600080fd5b5061014e6102e4366004610bbc565b6106c3565b3480156102f557600080fd5b5061016f610304366004610c5d565b6001600160a01b03918216600090815260026020908152604080832093909416825291909152205490565b61016f610710565b60606007805461034690610c90565b80601f016020809104026020016040519081016040528092919081815260200182805461037290610c90565b80156103bf5780601f10610394576101008083540402835291602001916103bf565b820191906000526020600020905b8154815290600101906020018083116103a257829003601f168201915b5050505050905090565b60006103d6338484610858565b50600192915050565b6000546001600160a01b031633146104125760405162461bcd60e51b815260040161040990610ccb565b60405180910390fd5b600b55565b6001600160a01b03831660009081526002602090815260408083203384529091528120548281101561049c5760405162461bcd60e51b815260206004820152602860248201527f45524332303a207472616e7366657220616d6f756e74206578636565647320616044820152676c6c6f77616e636560c01b6064820152608401610409565b6104b085336104ab8685610d29565b610858565b6104ba858461097c565b6104c48484610a1f565b836001600160a01b0316856001600160a01b0316600080516020610d9a833981519152856040516104f791815260200190565b60405180910390a3506001949350505050565b600061051533610a73565b33600081815260046020526040808220829055519293509183908381818185875af1925050503d8060008114610567576040519150601f19603f3d011682016040523d82523d6000602084013e61056c565b606091505b50509050806105b45760405162461bcd60e51b81526020600482015260146024820152732330b4b632b2103a379039b2b7321022ba3432b960611b6044820152606401610409565b5050565b6000546001600160a01b031633146105e25760405162461bcd60e51b815260040161040990610ccb565b6105ec8282610a1f565b6040518181526001600160a01b03831690600090600080516020610d9a833981519152906020015b60405180910390a35050565b600061062b82610ad8565b6001600160a01b03831660009081526004602052604090205461064e9190610d40565b92915050565b60606008805461034690610c90565b6000546001600160a01b0316331461068d5760405162461bcd60e51b815260040161040990610ccb565b610697828261097c565b6040518181526000906001600160a01b03841690600080516020610d9a83398151915290602001610614565b60006106cf338361097c565b6106d98383610a1f565b6040518281526001600160a01b038416903390600080516020610d9a8339815191529060200160405180910390a350600192915050565b600080546001600160a01b0316331461073b5760405162461bcd60e51b815260040161040990610ccb565b600b543490600090612710906107519084610d58565b61075b9190610d77565b90508181106107ac5760405162461bcd60e51b815260206004820152601860248201527f696e76616c69642076616c696461746f722072657761726400000000000000006044820152606401610409565b6107b68183610d29565b600a546040519193506001600160a01b03169082156108fc029083906000818181858888f193505050501580156107f1573d6000803e3d6000fd5b50600354600090610806633b9aca0085610d58565b6108109190610d77565b9050806006546108209190610d40565b600655600354600090633b9aca00906108399084610d58565b6108439190610d77565b905061084f8184610d40565b94505050505090565b6001600160a01b0383166108ba5760405162461bcd60e51b8152602060048201526024808201527f45524332303a20617070726f76652066726f6d20746865207a65726f206164646044820152637265737360e01b6064820152608401610409565b6001600160a01b03821661091b5760405162461bcd60e51b815260206004820152602260248201527f45524332303a20617070726f766520746f20746865207a65726f206164647265604482015261737360f01b6064820152608401610409565b6001600160a01b0383811660008181526002602090815260408083209487168084529482529182902085905590518481527f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925910160405180910390a3505050565b61098582610a73565b506001600160a01b038216600090815260016020526040902054808211156109ac57600080fd5b808210156109dc576109be8282610d29565b6001600160a01b038416600090815260016020526040902055610a03565b6001600160a01b038316600090815260016020908152604080832083905560059091528120555b8160036000828254610a159190610d29565b9091555050505050565b610a2882610a73565b506001600160a01b03821660009081526001602052604081208054839290610a51908490610d40565b925050819055508060036000828254610a6a9190610d40565b90915550505050565b600080610a7f83610ad8565b6001600160a01b038416600090815260046020526040902054909150610aa6908290610d40565b6001600160a01b0390931660009081526004602090815260408083208690556006546005909252909120555090919050565b6001600160a01b03811660009081526001602052604081205480610aff5750600092915050565b6001600160a01b038316600090815260056020526040812054600654610b259190610d29565b90506000633b9aca00610b388484610d58565b610b429190610d77565b95945050505050565b600060208083528351808285015260005b81811015610b7857858101830151858201604001528201610b5c565b81811115610b8a576000604083870101525b50601f01601f1916929092016040019392505050565b80356001600160a01b0381168114610bb757600080fd5b919050565b60008060408385031215610bcf57600080fd5b610bd883610ba0565b946020939093013593505050565b600060208284031215610bf857600080fd5b5035919050565b600080600060608486031215610c1457600080fd5b610c1d84610ba0565b9250610c2b60208501610ba0565b9150604084013590509250925092565b600060208284031215610c4d57600080fd5b610c5682610ba0565b9392505050565b60008060408385031215610c7057600080fd5b610c7983610ba0565b9150610c8760208401610ba0565b90509250929050565b600181811c90821680610ca457607f821691505b60208210811415610cc557634e487b7160e01b600052602260045260246000fd5b50919050565b60208082526028908201527f43616c6c207265737472696374656420746f20746865204175746f6e6974792060408201526710dbdb9d1c9858dd60c21b606082015260800190565b634e487b7160e01b600052601160045260246000fd5b600082821015610d3b57610d3b610d13565b500390565b60008219821115610d5357610d53610d13565b500190565b6000816000190483118215151615610d7257610d72610d13565b500290565b600082610d9457634e487b7160e01b600052601260045260246000fd5b50049056feddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3efa26469706673582212205311cbb54f78618267a028af290c7924ffb605213b8e8d196c0da082fc228cb264736f6c634300080c0033a26469706673582212204e32d007cb61aa0c27408ae1e76684ef49acf7ac86407c0c13449b3f8688356464736f6c634300080c003360806040523480156200001157600080fd5b506040516200116a3803806200116a833981016040819052620000349162000212565b6127108211156200004457600080fd5b600980546001600160a01b038087166001600160a01b031992831617909255600a805492861692909116919091179055600b8290556040516200008c908290602001620002ff565b60405160208183030381529060405260079080519060200190620000b29291906200010a565b5080604051602001620000c69190620002ff565b60405160208183030381529060405260089080519060200190620000ec9291906200010a565b5050600080546001600160a01b03191633179055506200036b915050565b82805462000118906200032e565b90600052602060002090601f0160209004810192826200013c576000855562000187565b82601f106200015757805160ff191683800117855562000187565b8280016001018555821562000187579182015b82811115620001875782518255916020019190600101906200016a565b506200019592915062000199565b5090565b5b808211156200019557600081556001016200019a565b6001600160a01b0381168114620001c657600080fd5b50565b634e487b7160e01b600052604160045260246000fd5b60005b83811015620001fc578181015183820152602001620001e2565b838111156200020c576000848401525b50505050565b600080600080608085870312156200022957600080fd5b84516200023681620001b0565b60208601519094506200024981620001b0565b6040860151606087015191945092506001600160401b03808211156200026e57600080fd5b818701915087601f8301126200028357600080fd5b815181811115620002985762000298620001c9565b604051601f8201601f19908116603f01168101908382118183101715620002c357620002c3620001c9565b816040528281528a6020848701011115620002dd57600080fd5b620002f0836020830160208801620001df565b979a9699509497505050505050565b644c4e544e2d60d81b81526000825162000321816005850160208701620001df565b9190910160050192915050565b600181811c908216806200034357607f821691505b602082108114156200036557634e487b7160e01b600052602260045260246000fd5b50919050565b610def806200037b6000396000f3fe6080604052600436106100fe5760003560e01c8063372500ab1161009557806395d89b411161006457806395d89b41146102945780639dc29fac146102a9578063a9059cbb146102c9578063dd62ed3e146102e9578063fb489a7b1461032f57600080fd5b8063372500ab1461020957806340c10f191461021e57806370a082311461023e578063949813b81461027457600080fd5b806319fac8fd116100d157806319fac8fd1461019557806323b872dd146101b75780632f2c3f2e146101d7578063313ce567146101ed57600080fd5b806306fdde0314610103578063095ea7b31461012e57806318160ddd1461015e578063187cf4d71461017d575b600080fd5b34801561010f57600080fd5b50610118610337565b6040516101259190610b4b565b60405180910390f35b34801561013a57600080fd5b5061014e610149366004610bbc565b6103c9565b6040519015158152602001610125565b34801561016a57600080fd5b506003545b604051908152602001610125565b34801561018957600080fd5b5061016f633b9aca0081565b3480156101a157600080fd5b506101b56101b0366004610be6565b6103df565b005b3480156101c357600080fd5b5061014e6101d2366004610bff565b610417565b3480156101e357600080fd5b5061016f61271081565b3480156101f957600080fd5b5060405160128152602001610125565b34801561021557600080fd5b506101b561050a565b34801561022a57600080fd5b506101b5610239366004610bbc565b6105b8565b34801561024a57600080fd5b5061016f610259366004610c3b565b6001600160a01b031660009081526001602052604090205490565b34801561028057600080fd5b5061016f61028f366004610c3b565b610620565b3480156102a057600080fd5b50610118610654565b3480156102b557600080fd5b506101b56102c4366004610bbc565b610663565b3480156102d557600080fd5b5061014e6102e4366004610bbc565b6106c3565b3480156102f557600080fd5b5061016f610304366004610c5d565b6001600160a01b03918216600090815260026020908152604080832093909416825291909152205490565b61016f610710565b60606007805461034690610c90565b80601f016020809104026020016040519081016040528092919081815260200182805461037290610c90565b80156103bf5780601f10610394576101008083540402835291602001916103bf565b820191906000526020600020905b8154815290600101906020018083116103a257829003601f168201915b5050505050905090565b60006103d6338484610858565b50600192915050565b6000546001600160a01b031633146104125760405162461bcd60e51b815260040161040990610ccb565b60405180910390fd5b600b55565b6001600160a01b03831660009081526002602090815260408083203384529091528120548281101561049c5760405162461bcd60e51b815260206004820152602860248201527f45524332303a207472616e7366657220616d6f756e74206578636565647320616044820152676c6c6f77616e636560c01b6064820152608401610409565b6104b085336104ab8685610d29565b610858565b6104ba858461097c565b6104c48484610a1f565b836001600160a01b0316856001600160a01b0316600080516020610d9a833981519152856040516104f791815260200190565b60405180910390a3506001949350505050565b600061051533610a73565b33600081815260046020526040808220829055519293509183908381818185875af1925050503d8060008114610567576040519150601f19603f3d011682016040523d82523d6000602084013e61056c565b606091505b50509050806105b45760405162461bcd60e51b81526020600482015260146024820152732330b4b632b2103a379039b2b7321022ba3432b960611b6044820152606401610409565b5050565b6000546001600160a01b031633146105e25760405162461bcd60e51b815260040161040990610ccb565b6105ec8282610a1f565b6040518181526001600160a01b03831690600090600080516020610d9a833981519152906020015b60405180910390a35050565b600061062b82610ad8565b6001600160a01b03831660009081526004602052604090205461064e9190610d40565b92915050565b60606008805461034690610c90565b6000546001600160a01b0316331461068d5760405162461bcd60e51b815260040161040990610ccb565b610697828261097c565b6040518181526000906001600160a01b03841690600080516020610d9a83398151915290602001610614565b60006106cf338361097c565b6106d98383610a1f565b6040518281526001600160a01b038416903390600080516020610d9a8339815191529060200160405180910390a350600192915050565b600080546001600160a01b0316331461073b5760405162461bcd60e51b815260040161040990610ccb565b600b543490600090612710906107519084610d58565b61075b9190610d77565b90508181106107ac5760405162461bcd60e51b815260206004820152601860248201527f696e76616c69642076616c696461746f722072657761726400000000000000006044820152606401610409565b6107b68183610d29565b600a546040519193506001600160a01b03169082156108fc029083906000818181858888f193505050501580156107f1573d6000803e3d6000fd5b50600354600090610806633b9aca0085610d58565b6108109190610d77565b9050806006546108209190610d40565b600655600354600090633b9aca00906108399084610d58565b6108439190610d77565b905061084f8184610d40565b94505050505090565b6001600160a01b0383166108ba5760405162461bcd60e51b8152602060048201526024808201527f45524332303a20617070726f76652066726f6d20746865207a65726f206164646044820152637265737360e01b6064820152608401610409565b6001600160a01b03821661091b5760405162461bcd60e51b815260206004820152602260248201527f45524332303a20617070726f766520746f20746865207a65726f206164647265604482015261737360f01b6064820152608401610409565b6001600160a01b0383811660008181526002602090815260408083209487168084529482529182902085905590518481527f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925910160405180910390a3505050565b61098582610a73565b506001600160a01b038216600090815260016020526040902054808211156109ac57600080fd5b808210156109dc576109be8282610d29565b6001600160a01b038416600090815260016020526040902055610a03565b6001600160a01b038316600090815260016020908152604080832083905560059091528120555b8160036000828254610a159190610d29565b9091555050505050565b610a2882610a73565b506001600160a01b03821660009081526001602052604081208054839290610a51908490610d40565b925050819055508060036000828254610a6a9190610d40565b90915550505050565b600080610a7f83610ad8565b6001600160a01b038416600090815260046020526040902054909150610aa6908290610d40565b6001600160a01b0390931660009081526004602090815260408083208690556006546005909252909120555090919050565b6001600160a01b03811660009081526001602052604081205480610aff5750600092915050565b6001600160a01b038316600090815260056020526040812054600654610b259190610d29565b90506000633b9aca00610b388484610d58565b610b429190610d77565b95945050505050565b600060208083528351808285015260005b81811015610b7857858101830151858201604001528201610b5c565b81811115610b8a576000604083870101525b50601f01601f1916929092016040019392505050565b80356001600160a01b0381168114610bb757600080fd5b919050565b60008060408385031215610bcf57600080fd5b610bd883610ba0565b946020939093013593505050565b600060208284031215610bf857600080fd5b5035919050565b600080600060608486031215610c1457600080fd5b610c1d84610ba0565b9250610c2b60208501610ba0565b9150604084013590509250925092565b600060208284031215610c4d57600080fd5b610c5682610ba0565b9392505050565b60008060408385031215610c7057600080fd5b610c7983610ba0565b9150610c8760208401610ba0565b90509250929050565b600181811c90821680610ca457607f821691505b60208210811415610cc557634e487b7160e01b600052602260045260246000fd5b50919050565b60208082526028908201527f43616c6c207265737472696374656420746f20746865204175746f6e6974792060408201526710dbdb9d1c9858dd60c21b606082015260800190565b634e487b7160e01b600052601160045260246000fd5b600082821015610d3b57610d3b610d13565b500390565b60008219821115610d5357610d53610d13565b500190565b6000816000190483118215151615610d7257610d72610d13565b500290565b600082610d9457634e487b7160e01b600052601260045260246000fd5b50049056feddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3efa26469706673582212205311cbb54f78618267a028af290c7924ffb605213b8e8d196c0da082fc228cb264736f6c634300080c0033",
}

// AutonityABI is the input ABI used to generate the binding from.
// Deprecated: Use AutonityMetaData.ABI instead.
var AutonityABI = AutonityMetaData.ABI

// Deprecated: Use AutonityMetaData.Sigs instead.
// AutonityFuncSigs maps the 4-byte function signature to its string representation.
var AutonityFuncSigs = AutonityMetaData.Sigs

// AutonityBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use AutonityMetaData.Bin instead.
var AutonityBin = AutonityMetaData.Bin

// DeployAutonity deploys a new Ethereum contract, binding an instance of Autonity to it.
func DeployAutonity(auth *bind.TransactOpts, backend bind.ContractBackend, _validators []AutonityValidator, _config AutonityConfig) (common.Address, *types.Transaction, *Autonity, error) {
	parsed, err := AutonityMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(AutonityBin), backend, _validators, _config)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Autonity{AutonityCaller: AutonityCaller{contract: contract}, AutonityTransactor: AutonityTransactor{contract: contract}, AutonityFilterer: AutonityFilterer{contract: contract}}, nil
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
func (_Autonity *AutonityRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
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
func (_Autonity *AutonityCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
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

// COMMISSIONRATEPRECISION is a free data retrieval call binding the contract method 0x2f2c3f2e.
//
// Solidity: function COMMISSION_RATE_PRECISION() view returns(uint256)
func (_Autonity *AutonityCaller) COMMISSIONRATEPRECISION(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Autonity.contract.Call(opts, &out, "COMMISSION_RATE_PRECISION")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// COMMISSIONRATEPRECISION is a free data retrieval call binding the contract method 0x2f2c3f2e.
//
// Solidity: function COMMISSION_RATE_PRECISION() view returns(uint256)
func (_Autonity *AutonitySession) COMMISSIONRATEPRECISION() (*big.Int, error) {
	return _Autonity.Contract.COMMISSIONRATEPRECISION(&_Autonity.CallOpts)
}

// COMMISSIONRATEPRECISION is a free data retrieval call binding the contract method 0x2f2c3f2e.
//
// Solidity: function COMMISSION_RATE_PRECISION() view returns(uint256)
func (_Autonity *AutonityCallerSession) COMMISSIONRATEPRECISION() (*big.Int, error) {
	return _Autonity.Contract.COMMISSIONRATEPRECISION(&_Autonity.CallOpts)
}

// AccusationProcessed is a free data retrieval call binding the contract method 0xffd9d914.
//
// Solidity: function accusationProcessed(bytes32 _msgHash) view returns(bool)
func (_Autonity *AutonityCaller) AccusationProcessed(opts *bind.CallOpts, _msgHash [32]byte) (bool, error) {
	var out []interface{}
	err := _Autonity.contract.Call(opts, &out, "accusationProcessed", _msgHash)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// AccusationProcessed is a free data retrieval call binding the contract method 0xffd9d914.
//
// Solidity: function accusationProcessed(bytes32 _msgHash) view returns(bool)
func (_Autonity *AutonitySession) AccusationProcessed(_msgHash [32]byte) (bool, error) {
	return _Autonity.Contract.AccusationProcessed(&_Autonity.CallOpts, _msgHash)
}

// AccusationProcessed is a free data retrieval call binding the contract method 0xffd9d914.
//
// Solidity: function accusationProcessed(bytes32 _msgHash) view returns(bool)
func (_Autonity *AutonityCallerSession) AccusationProcessed(_msgHash [32]byte) (bool, error) {
	return _Autonity.Contract.AccusationProcessed(&_Autonity.CallOpts, _msgHash)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_Autonity *AutonityCaller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Autonity.contract.Call(opts, &out, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_Autonity *AutonitySession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _Autonity.Contract.Allowance(&_Autonity.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_Autonity *AutonityCallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _Autonity.Contract.Allowance(&_Autonity.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address _addr) view returns(uint256)
func (_Autonity *AutonityCaller) BalanceOf(opts *bind.CallOpts, _addr common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Autonity.contract.Call(opts, &out, "balanceOf", _addr)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address _addr) view returns(uint256)
func (_Autonity *AutonitySession) BalanceOf(_addr common.Address) (*big.Int, error) {
	return _Autonity.Contract.BalanceOf(&_Autonity.CallOpts, _addr)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address _addr) view returns(uint256)
func (_Autonity *AutonityCallerSession) BalanceOf(_addr common.Address) (*big.Int, error) {
	return _Autonity.Contract.BalanceOf(&_Autonity.CallOpts, _addr)
}

// Config is a free data retrieval call binding the contract method 0x79502c55.
//
// Solidity: function config() view returns(address operatorAccount, address treasuryAccount, uint256 treasuryFee, uint256 minBaseFee, uint256 delegationRate, uint256 epochPeriod, uint256 unbondingPeriod, uint256 committeeSize, uint256 contractVersion, uint256 blockPeriod)
func (_Autonity *AutonityCaller) Config(opts *bind.CallOpts) (struct {
	OperatorAccount common.Address
	TreasuryAccount common.Address
	TreasuryFee     *big.Int
	MinBaseFee      *big.Int
	DelegationRate  *big.Int
	EpochPeriod     *big.Int
	UnbondingPeriod *big.Int
	CommitteeSize   *big.Int
	ContractVersion *big.Int
	BlockPeriod     *big.Int
}, error) {
	var out []interface{}
	err := _Autonity.contract.Call(opts, &out, "config")

	outstruct := new(struct {
		OperatorAccount common.Address
		TreasuryAccount common.Address
		TreasuryFee     *big.Int
		MinBaseFee      *big.Int
		DelegationRate  *big.Int
		EpochPeriod     *big.Int
		UnbondingPeriod *big.Int
		CommitteeSize   *big.Int
		ContractVersion *big.Int
		BlockPeriod     *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.OperatorAccount = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.TreasuryAccount = *abi.ConvertType(out[1], new(common.Address)).(*common.Address)
	outstruct.TreasuryFee = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.MinBaseFee = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.DelegationRate = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)
	outstruct.EpochPeriod = *abi.ConvertType(out[5], new(*big.Int)).(**big.Int)
	outstruct.UnbondingPeriod = *abi.ConvertType(out[6], new(*big.Int)).(**big.Int)
	outstruct.CommitteeSize = *abi.ConvertType(out[7], new(*big.Int)).(**big.Int)
	outstruct.ContractVersion = *abi.ConvertType(out[8], new(*big.Int)).(**big.Int)
	outstruct.BlockPeriod = *abi.ConvertType(out[9], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// Config is a free data retrieval call binding the contract method 0x79502c55.
//
// Solidity: function config() view returns(address operatorAccount, address treasuryAccount, uint256 treasuryFee, uint256 minBaseFee, uint256 delegationRate, uint256 epochPeriod, uint256 unbondingPeriod, uint256 committeeSize, uint256 contractVersion, uint256 blockPeriod)
func (_Autonity *AutonitySession) Config() (struct {
	OperatorAccount common.Address
	TreasuryAccount common.Address
	TreasuryFee     *big.Int
	MinBaseFee      *big.Int
	DelegationRate  *big.Int
	EpochPeriod     *big.Int
	UnbondingPeriod *big.Int
	CommitteeSize   *big.Int
	ContractVersion *big.Int
	BlockPeriod     *big.Int
}, error) {
	return _Autonity.Contract.Config(&_Autonity.CallOpts)
}

// Config is a free data retrieval call binding the contract method 0x79502c55.
//
// Solidity: function config() view returns(address operatorAccount, address treasuryAccount, uint256 treasuryFee, uint256 minBaseFee, uint256 delegationRate, uint256 epochPeriod, uint256 unbondingPeriod, uint256 committeeSize, uint256 contractVersion, uint256 blockPeriod)
func (_Autonity *AutonityCallerSession) Config() (struct {
	OperatorAccount common.Address
	TreasuryAccount common.Address
	TreasuryFee     *big.Int
	MinBaseFee      *big.Int
	DelegationRate  *big.Int
	EpochPeriod     *big.Int
	UnbondingPeriod *big.Int
	CommitteeSize   *big.Int
	ContractVersion *big.Int
	BlockPeriod     *big.Int
}, error) {
	return _Autonity.Contract.Config(&_Autonity.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() pure returns(uint8)
func (_Autonity *AutonityCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _Autonity.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() pure returns(uint8)
func (_Autonity *AutonitySession) Decimals() (uint8, error) {
	return _Autonity.Contract.Decimals(&_Autonity.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() pure returns(uint8)
func (_Autonity *AutonityCallerSession) Decimals() (uint8, error) {
	return _Autonity.Contract.Decimals(&_Autonity.CallOpts)
}

// Deployer is a free data retrieval call binding the contract method 0xd5f39488.
//
// Solidity: function deployer() view returns(address)
func (_Autonity *AutonityCaller) Deployer(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Autonity.contract.Call(opts, &out, "deployer")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Deployer is a free data retrieval call binding the contract method 0xd5f39488.
//
// Solidity: function deployer() view returns(address)
func (_Autonity *AutonitySession) Deployer() (common.Address, error) {
	return _Autonity.Contract.Deployer(&_Autonity.CallOpts)
}

// Deployer is a free data retrieval call binding the contract method 0xd5f39488.
//
// Solidity: function deployer() view returns(address)
func (_Autonity *AutonityCallerSession) Deployer() (common.Address, error) {
	return _Autonity.Contract.Deployer(&_Autonity.CallOpts)
}

// EpochID is a free data retrieval call binding the contract method 0xc9d97af4.
//
// Solidity: function epochID() view returns(uint256)
func (_Autonity *AutonityCaller) EpochID(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Autonity.contract.Call(opts, &out, "epochID")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// EpochID is a free data retrieval call binding the contract method 0xc9d97af4.
//
// Solidity: function epochID() view returns(uint256)
func (_Autonity *AutonitySession) EpochID() (*big.Int, error) {
	return _Autonity.Contract.EpochID(&_Autonity.CallOpts)
}

// EpochID is a free data retrieval call binding the contract method 0xc9d97af4.
//
// Solidity: function epochID() view returns(uint256)
func (_Autonity *AutonityCallerSession) EpochID() (*big.Int, error) {
	return _Autonity.Contract.EpochID(&_Autonity.CallOpts)
}

// EpochTotalBondedStake is a free data retrieval call binding the contract method 0x9c98e471.
//
// Solidity: function epochTotalBondedStake() view returns(uint256)
func (_Autonity *AutonityCaller) EpochTotalBondedStake(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Autonity.contract.Call(opts, &out, "epochTotalBondedStake")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// EpochTotalBondedStake is a free data retrieval call binding the contract method 0x9c98e471.
//
// Solidity: function epochTotalBondedStake() view returns(uint256)
func (_Autonity *AutonitySession) EpochTotalBondedStake() (*big.Int, error) {
	return _Autonity.Contract.EpochTotalBondedStake(&_Autonity.CallOpts)
}

// EpochTotalBondedStake is a free data retrieval call binding the contract method 0x9c98e471.
//
// Solidity: function epochTotalBondedStake() view returns(uint256)
func (_Autonity *AutonityCallerSession) EpochTotalBondedStake() (*big.Int, error) {
	return _Autonity.Contract.EpochTotalBondedStake(&_Autonity.CallOpts)
}

// GetAccountabilityEventChunk is a free data retrieval call binding the contract method 0xf446c557.
//
// Solidity: function getAccountabilityEventChunk(bytes32 _msgHash, uint8 _type, uint8 _rule, address _reporter, uint8 _chunkID) view returns(bytes)
func (_Autonity *AutonityCaller) GetAccountabilityEventChunk(opts *bind.CallOpts, _msgHash [32]byte, _type uint8, _rule uint8, _reporter common.Address, _chunkID uint8) ([]byte, error) {
	var out []interface{}
	err := _Autonity.contract.Call(opts, &out, "getAccountabilityEventChunk", _msgHash, _type, _rule, _reporter, _chunkID)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetAccountabilityEventChunk is a free data retrieval call binding the contract method 0xf446c557.
//
// Solidity: function getAccountabilityEventChunk(bytes32 _msgHash, uint8 _type, uint8 _rule, address _reporter, uint8 _chunkID) view returns(bytes)
func (_Autonity *AutonitySession) GetAccountabilityEventChunk(_msgHash [32]byte, _type uint8, _rule uint8, _reporter common.Address, _chunkID uint8) ([]byte, error) {
	return _Autonity.Contract.GetAccountabilityEventChunk(&_Autonity.CallOpts, _msgHash, _type, _rule, _reporter, _chunkID)
}

// GetAccountabilityEventChunk is a free data retrieval call binding the contract method 0xf446c557.
//
// Solidity: function getAccountabilityEventChunk(bytes32 _msgHash, uint8 _type, uint8 _rule, address _reporter, uint8 _chunkID) view returns(bytes)
func (_Autonity *AutonityCallerSession) GetAccountabilityEventChunk(_msgHash [32]byte, _type uint8, _rule uint8, _reporter common.Address, _chunkID uint8) ([]byte, error) {
	return _Autonity.Contract.GetAccountabilityEventChunk(&_Autonity.CallOpts, _msgHash, _type, _rule, _reporter, _chunkID)
}

// GetBlockPeriod is a free data retrieval call binding the contract method 0x43645969.
//
// Solidity: function getBlockPeriod() view returns(uint256)
func (_Autonity *AutonityCaller) GetBlockPeriod(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Autonity.contract.Call(opts, &out, "getBlockPeriod")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetBlockPeriod is a free data retrieval call binding the contract method 0x43645969.
//
// Solidity: function getBlockPeriod() view returns(uint256)
func (_Autonity *AutonitySession) GetBlockPeriod() (*big.Int, error) {
	return _Autonity.Contract.GetBlockPeriod(&_Autonity.CallOpts)
}

// GetBlockPeriod is a free data retrieval call binding the contract method 0x43645969.
//
// Solidity: function getBlockPeriod() view returns(uint256)
func (_Autonity *AutonityCallerSession) GetBlockPeriod() (*big.Int, error) {
	return _Autonity.Contract.GetBlockPeriod(&_Autonity.CallOpts)
}

// GetBondingReq is a free data retrieval call binding the contract method 0xe485c6fb.
//
// Solidity: function getBondingReq(uint256 startId, uint256 lastId) view returns((address,address,uint256,uint256)[])
func (_Autonity *AutonityCaller) GetBondingReq(opts *bind.CallOpts, startId *big.Int, lastId *big.Int) ([]AutonityStaking, error) {
	var out []interface{}
	err := _Autonity.contract.Call(opts, &out, "getBondingReq", startId, lastId)

	if err != nil {
		return *new([]AutonityStaking), err
	}

	out0 := *abi.ConvertType(out[0], new([]AutonityStaking)).(*[]AutonityStaking)

	return out0, err

}

// GetBondingReq is a free data retrieval call binding the contract method 0xe485c6fb.
//
// Solidity: function getBondingReq(uint256 startId, uint256 lastId) view returns((address,address,uint256,uint256)[])
func (_Autonity *AutonitySession) GetBondingReq(startId *big.Int, lastId *big.Int) ([]AutonityStaking, error) {
	return _Autonity.Contract.GetBondingReq(&_Autonity.CallOpts, startId, lastId)
}

// GetBondingReq is a free data retrieval call binding the contract method 0xe485c6fb.
//
// Solidity: function getBondingReq(uint256 startId, uint256 lastId) view returns((address,address,uint256,uint256)[])
func (_Autonity *AutonityCallerSession) GetBondingReq(startId *big.Int, lastId *big.Int) ([]AutonityStaking, error) {
	return _Autonity.Contract.GetBondingReq(&_Autonity.CallOpts, startId, lastId)
}

// GetCommittee is a free data retrieval call binding the contract method 0xab8f6ffe.
//
// Solidity: function getCommittee() view returns((address,uint256)[])
func (_Autonity *AutonityCaller) GetCommittee(opts *bind.CallOpts) ([]AutonityCommitteeMember, error) {
	var out []interface{}
	err := _Autonity.contract.Call(opts, &out, "getCommittee")

	if err != nil {
		return *new([]AutonityCommitteeMember), err
	}

	out0 := *abi.ConvertType(out[0], new([]AutonityCommitteeMember)).(*[]AutonityCommitteeMember)

	return out0, err

}

// GetCommittee is a free data retrieval call binding the contract method 0xab8f6ffe.
//
// Solidity: function getCommittee() view returns((address,uint256)[])
func (_Autonity *AutonitySession) GetCommittee() ([]AutonityCommitteeMember, error) {
	return _Autonity.Contract.GetCommittee(&_Autonity.CallOpts)
}

// GetCommittee is a free data retrieval call binding the contract method 0xab8f6ffe.
//
// Solidity: function getCommittee() view returns((address,uint256)[])
func (_Autonity *AutonityCallerSession) GetCommittee() ([]AutonityCommitteeMember, error) {
	return _Autonity.Contract.GetCommittee(&_Autonity.CallOpts)
}

// GetCommitteeEnodes is a free data retrieval call binding the contract method 0xa8b2216e.
//
// Solidity: function getCommitteeEnodes() view returns(string[])
func (_Autonity *AutonityCaller) GetCommitteeEnodes(opts *bind.CallOpts) ([]string, error) {
	var out []interface{}
	err := _Autonity.contract.Call(opts, &out, "getCommitteeEnodes")

	if err != nil {
		return *new([]string), err
	}

	out0 := *abi.ConvertType(out[0], new([]string)).(*[]string)

	return out0, err

}

// GetCommitteeEnodes is a free data retrieval call binding the contract method 0xa8b2216e.
//
// Solidity: function getCommitteeEnodes() view returns(string[])
func (_Autonity *AutonitySession) GetCommitteeEnodes() ([]string, error) {
	return _Autonity.Contract.GetCommitteeEnodes(&_Autonity.CallOpts)
}

// GetCommitteeEnodes is a free data retrieval call binding the contract method 0xa8b2216e.
//
// Solidity: function getCommitteeEnodes() view returns(string[])
func (_Autonity *AutonityCallerSession) GetCommitteeEnodes() ([]string, error) {
	return _Autonity.Contract.GetCommitteeEnodes(&_Autonity.CallOpts)
}

// GetEpochPeriod is a free data retrieval call binding the contract method 0xdfb1a4d2.
//
// Solidity: function getEpochPeriod() view returns(uint256)
func (_Autonity *AutonityCaller) GetEpochPeriod(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Autonity.contract.Call(opts, &out, "getEpochPeriod")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetEpochPeriod is a free data retrieval call binding the contract method 0xdfb1a4d2.
//
// Solidity: function getEpochPeriod() view returns(uint256)
func (_Autonity *AutonitySession) GetEpochPeriod() (*big.Int, error) {
	return _Autonity.Contract.GetEpochPeriod(&_Autonity.CallOpts)
}

// GetEpochPeriod is a free data retrieval call binding the contract method 0xdfb1a4d2.
//
// Solidity: function getEpochPeriod() view returns(uint256)
func (_Autonity *AutonityCallerSession) GetEpochPeriod() (*big.Int, error) {
	return _Autonity.Contract.GetEpochPeriod(&_Autonity.CallOpts)
}

// GetLastEpochBlock is a free data retrieval call binding the contract method 0x731b3a03.
//
// Solidity: function getLastEpochBlock() view returns(uint256)
func (_Autonity *AutonityCaller) GetLastEpochBlock(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Autonity.contract.Call(opts, &out, "getLastEpochBlock")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetLastEpochBlock is a free data retrieval call binding the contract method 0x731b3a03.
//
// Solidity: function getLastEpochBlock() view returns(uint256)
func (_Autonity *AutonitySession) GetLastEpochBlock() (*big.Int, error) {
	return _Autonity.Contract.GetLastEpochBlock(&_Autonity.CallOpts)
}

// GetLastEpochBlock is a free data retrieval call binding the contract method 0x731b3a03.
//
// Solidity: function getLastEpochBlock() view returns(uint256)
func (_Autonity *AutonityCallerSession) GetLastEpochBlock() (*big.Int, error) {
	return _Autonity.Contract.GetLastEpochBlock(&_Autonity.CallOpts)
}

// GetMaxCommitteeSize is a free data retrieval call binding the contract method 0x819b6463.
//
// Solidity: function getMaxCommitteeSize() view returns(uint256)
func (_Autonity *AutonityCaller) GetMaxCommitteeSize(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Autonity.contract.Call(opts, &out, "getMaxCommitteeSize")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetMaxCommitteeSize is a free data retrieval call binding the contract method 0x819b6463.
//
// Solidity: function getMaxCommitteeSize() view returns(uint256)
func (_Autonity *AutonitySession) GetMaxCommitteeSize() (*big.Int, error) {
	return _Autonity.Contract.GetMaxCommitteeSize(&_Autonity.CallOpts)
}

// GetMaxCommitteeSize is a free data retrieval call binding the contract method 0x819b6463.
//
// Solidity: function getMaxCommitteeSize() view returns(uint256)
func (_Autonity *AutonityCallerSession) GetMaxCommitteeSize() (*big.Int, error) {
	return _Autonity.Contract.GetMaxCommitteeSize(&_Autonity.CallOpts)
}

// GetMinimumBaseFee is a free data retrieval call binding the contract method 0x11220633.
//
// Solidity: function getMinimumBaseFee() view returns(uint256)
func (_Autonity *AutonityCaller) GetMinimumBaseFee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Autonity.contract.Call(opts, &out, "getMinimumBaseFee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetMinimumBaseFee is a free data retrieval call binding the contract method 0x11220633.
//
// Solidity: function getMinimumBaseFee() view returns(uint256)
func (_Autonity *AutonitySession) GetMinimumBaseFee() (*big.Int, error) {
	return _Autonity.Contract.GetMinimumBaseFee(&_Autonity.CallOpts)
}

// GetMinimumBaseFee is a free data retrieval call binding the contract method 0x11220633.
//
// Solidity: function getMinimumBaseFee() view returns(uint256)
func (_Autonity *AutonityCallerSession) GetMinimumBaseFee() (*big.Int, error) {
	return _Autonity.Contract.GetMinimumBaseFee(&_Autonity.CallOpts)
}

// GetNewContract is a free data retrieval call binding the contract method 0xb66b3e79.
//
// Solidity: function getNewContract() view returns(bytes, string)
func (_Autonity *AutonityCaller) GetNewContract(opts *bind.CallOpts) ([]byte, string, error) {
	var out []interface{}
	err := _Autonity.contract.Call(opts, &out, "getNewContract")

	if err != nil {
		return *new([]byte), *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)
	out1 := *abi.ConvertType(out[1], new(string)).(*string)

	return out0, out1, err

}

// GetNewContract is a free data retrieval call binding the contract method 0xb66b3e79.
//
// Solidity: function getNewContract() view returns(bytes, string)
func (_Autonity *AutonitySession) GetNewContract() ([]byte, string, error) {
	return _Autonity.Contract.GetNewContract(&_Autonity.CallOpts)
}

// GetNewContract is a free data retrieval call binding the contract method 0xb66b3e79.
//
// Solidity: function getNewContract() view returns(bytes, string)
func (_Autonity *AutonityCallerSession) GetNewContract() ([]byte, string, error) {
	return _Autonity.Contract.GetNewContract(&_Autonity.CallOpts)
}

// GetOperator is a free data retrieval call binding the contract method 0xe7f43c68.
//
// Solidity: function getOperator() view returns(address)
func (_Autonity *AutonityCaller) GetOperator(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Autonity.contract.Call(opts, &out, "getOperator")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetOperator is a free data retrieval call binding the contract method 0xe7f43c68.
//
// Solidity: function getOperator() view returns(address)
func (_Autonity *AutonitySession) GetOperator() (common.Address, error) {
	return _Autonity.Contract.GetOperator(&_Autonity.CallOpts)
}

// GetOperator is a free data retrieval call binding the contract method 0xe7f43c68.
//
// Solidity: function getOperator() view returns(address)
func (_Autonity *AutonityCallerSession) GetOperator() (common.Address, error) {
	return _Autonity.Contract.GetOperator(&_Autonity.CallOpts)
}

// GetPenalty is a free data retrieval call binding the contract method 0xe56e56db.
//
// Solidity: function getPenalty() view returns(uint256)
func (_Autonity *AutonityCaller) GetPenalty(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Autonity.contract.Call(opts, &out, "getPenalty")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetPenalty is a free data retrieval call binding the contract method 0xe56e56db.
//
// Solidity: function getPenalty() view returns(uint256)
func (_Autonity *AutonitySession) GetPenalty() (*big.Int, error) {
	return _Autonity.Contract.GetPenalty(&_Autonity.CallOpts)
}

// GetPenalty is a free data retrieval call binding the contract method 0xe56e56db.
//
// Solidity: function getPenalty() view returns(uint256)
func (_Autonity *AutonityCallerSession) GetPenalty() (*big.Int, error) {
	return _Autonity.Contract.GetPenalty(&_Autonity.CallOpts)
}

// GetSlashedStake is a free data retrieval call binding the contract method 0xfe44c7f5.
//
// Solidity: function getSlashedStake(address _addr) view returns(uint256)
func (_Autonity *AutonityCaller) GetSlashedStake(opts *bind.CallOpts, _addr common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Autonity.contract.Call(opts, &out, "getSlashedStake", _addr)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetSlashedStake is a free data retrieval call binding the contract method 0xfe44c7f5.
//
// Solidity: function getSlashedStake(address _addr) view returns(uint256)
func (_Autonity *AutonitySession) GetSlashedStake(_addr common.Address) (*big.Int, error) {
	return _Autonity.Contract.GetSlashedStake(&_Autonity.CallOpts, _addr)
}

// GetSlashedStake is a free data retrieval call binding the contract method 0xfe44c7f5.
//
// Solidity: function getSlashedStake(address _addr) view returns(uint256)
func (_Autonity *AutonityCallerSession) GetSlashedStake(_addr common.Address) (*big.Int, error) {
	return _Autonity.Contract.GetSlashedStake(&_Autonity.CallOpts, _addr)
}

// GetTreasuryAccount is a free data retrieval call binding the contract method 0xf7866ee3.
//
// Solidity: function getTreasuryAccount() view returns(address)
func (_Autonity *AutonityCaller) GetTreasuryAccount(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Autonity.contract.Call(opts, &out, "getTreasuryAccount")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetTreasuryAccount is a free data retrieval call binding the contract method 0xf7866ee3.
//
// Solidity: function getTreasuryAccount() view returns(address)
func (_Autonity *AutonitySession) GetTreasuryAccount() (common.Address, error) {
	return _Autonity.Contract.GetTreasuryAccount(&_Autonity.CallOpts)
}

// GetTreasuryAccount is a free data retrieval call binding the contract method 0xf7866ee3.
//
// Solidity: function getTreasuryAccount() view returns(address)
func (_Autonity *AutonityCallerSession) GetTreasuryAccount() (common.Address, error) {
	return _Autonity.Contract.GetTreasuryAccount(&_Autonity.CallOpts)
}

// GetTreasuryFee is a free data retrieval call binding the contract method 0x29070c6d.
//
// Solidity: function getTreasuryFee() view returns(uint256)
func (_Autonity *AutonityCaller) GetTreasuryFee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Autonity.contract.Call(opts, &out, "getTreasuryFee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetTreasuryFee is a free data retrieval call binding the contract method 0x29070c6d.
//
// Solidity: function getTreasuryFee() view returns(uint256)
func (_Autonity *AutonitySession) GetTreasuryFee() (*big.Int, error) {
	return _Autonity.Contract.GetTreasuryFee(&_Autonity.CallOpts)
}

// GetTreasuryFee is a free data retrieval call binding the contract method 0x29070c6d.
//
// Solidity: function getTreasuryFee() view returns(uint256)
func (_Autonity *AutonityCallerSession) GetTreasuryFee() (*big.Int, error) {
	return _Autonity.Contract.GetTreasuryFee(&_Autonity.CallOpts)
}

// GetUnbondingPeriod is a free data retrieval call binding the contract method 0x6fd2c80b.
//
// Solidity: function getUnbondingPeriod() view returns(uint256)
func (_Autonity *AutonityCaller) GetUnbondingPeriod(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Autonity.contract.Call(opts, &out, "getUnbondingPeriod")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetUnbondingPeriod is a free data retrieval call binding the contract method 0x6fd2c80b.
//
// Solidity: function getUnbondingPeriod() view returns(uint256)
func (_Autonity *AutonitySession) GetUnbondingPeriod() (*big.Int, error) {
	return _Autonity.Contract.GetUnbondingPeriod(&_Autonity.CallOpts)
}

// GetUnbondingPeriod is a free data retrieval call binding the contract method 0x6fd2c80b.
//
// Solidity: function getUnbondingPeriod() view returns(uint256)
func (_Autonity *AutonityCallerSession) GetUnbondingPeriod() (*big.Int, error) {
	return _Autonity.Contract.GetUnbondingPeriod(&_Autonity.CallOpts)
}

// GetUnbondingReq is a free data retrieval call binding the contract method 0x55230e93.
//
// Solidity: function getUnbondingReq(uint256 startId, uint256 lastId) view returns((address,address,uint256,uint256)[])
func (_Autonity *AutonityCaller) GetUnbondingReq(opts *bind.CallOpts, startId *big.Int, lastId *big.Int) ([]AutonityStaking, error) {
	var out []interface{}
	err := _Autonity.contract.Call(opts, &out, "getUnbondingReq", startId, lastId)

	if err != nil {
		return *new([]AutonityStaking), err
	}

	out0 := *abi.ConvertType(out[0], new([]AutonityStaking)).(*[]AutonityStaking)

	return out0, err

}

// GetUnbondingReq is a free data retrieval call binding the contract method 0x55230e93.
//
// Solidity: function getUnbondingReq(uint256 startId, uint256 lastId) view returns((address,address,uint256,uint256)[])
func (_Autonity *AutonitySession) GetUnbondingReq(startId *big.Int, lastId *big.Int) ([]AutonityStaking, error) {
	return _Autonity.Contract.GetUnbondingReq(&_Autonity.CallOpts, startId, lastId)
}

// GetUnbondingReq is a free data retrieval call binding the contract method 0x55230e93.
//
// Solidity: function getUnbondingReq(uint256 startId, uint256 lastId) view returns((address,address,uint256,uint256)[])
func (_Autonity *AutonityCallerSession) GetUnbondingReq(startId *big.Int, lastId *big.Int) ([]AutonityStaking, error) {
	return _Autonity.Contract.GetUnbondingReq(&_Autonity.CallOpts, startId, lastId)
}

// GetValidator is a free data retrieval call binding the contract method 0x1904bb2e.
//
// Solidity: function getValidator(address _addr) view returns((address,address,address,string,uint256,uint256,uint256,address,uint256,uint256,uint8))
func (_Autonity *AutonityCaller) GetValidator(opts *bind.CallOpts, _addr common.Address) (AutonityValidator, error) {
	var out []interface{}
	err := _Autonity.contract.Call(opts, &out, "getValidator", _addr)

	if err != nil {
		return *new(AutonityValidator), err
	}

	out0 := *abi.ConvertType(out[0], new(AutonityValidator)).(*AutonityValidator)

	return out0, err

}

// GetValidator is a free data retrieval call binding the contract method 0x1904bb2e.
//
// Solidity: function getValidator(address _addr) view returns((address,address,address,string,uint256,uint256,uint256,address,uint256,uint256,uint8))
func (_Autonity *AutonitySession) GetValidator(_addr common.Address) (AutonityValidator, error) {
	return _Autonity.Contract.GetValidator(&_Autonity.CallOpts, _addr)
}

// GetValidator is a free data retrieval call binding the contract method 0x1904bb2e.
//
// Solidity: function getValidator(address _addr) view returns((address,address,address,string,uint256,uint256,uint256,address,uint256,uint256,uint8))
func (_Autonity *AutonityCallerSession) GetValidator(_addr common.Address) (AutonityValidator, error) {
	return _Autonity.Contract.GetValidator(&_Autonity.CallOpts, _addr)
}

// GetValidatorRecentAccusations is a free data retrieval call binding the contract method 0x1bd38702.
//
// Solidity: function getValidatorRecentAccusations(address _addr) view returns((uint8,uint8,uint8,uint8,address,address,bytes32,bytes)[])
func (_Autonity *AutonityCaller) GetValidatorRecentAccusations(opts *bind.CallOpts, _addr common.Address) ([]AutonityAccountabilityEvent, error) {
	var out []interface{}
	err := _Autonity.contract.Call(opts, &out, "getValidatorRecentAccusations", _addr)

	if err != nil {
		return *new([]AutonityAccountabilityEvent), err
	}

	out0 := *abi.ConvertType(out[0], new([]AutonityAccountabilityEvent)).(*[]AutonityAccountabilityEvent)

	return out0, err

}

// GetValidatorRecentAccusations is a free data retrieval call binding the contract method 0x1bd38702.
//
// Solidity: function getValidatorRecentAccusations(address _addr) view returns((uint8,uint8,uint8,uint8,address,address,bytes32,bytes)[])
func (_Autonity *AutonitySession) GetValidatorRecentAccusations(_addr common.Address) ([]AutonityAccountabilityEvent, error) {
	return _Autonity.Contract.GetValidatorRecentAccusations(&_Autonity.CallOpts, _addr)
}

// GetValidatorRecentAccusations is a free data retrieval call binding the contract method 0x1bd38702.
//
// Solidity: function getValidatorRecentAccusations(address _addr) view returns((uint8,uint8,uint8,uint8,address,address,bytes32,bytes)[])
func (_Autonity *AutonityCallerSession) GetValidatorRecentAccusations(_addr common.Address) ([]AutonityAccountabilityEvent, error) {
	return _Autonity.Contract.GetValidatorRecentAccusations(&_Autonity.CallOpts, _addr)
}

// GetValidatorRecentMisbehaviours is a free data retrieval call binding the contract method 0xac306841.
//
// Solidity: function getValidatorRecentMisbehaviours(address _addr) view returns((uint8,uint8,uint8,uint8,address,address,bytes32,bytes)[])
func (_Autonity *AutonityCaller) GetValidatorRecentMisbehaviours(opts *bind.CallOpts, _addr common.Address) ([]AutonityAccountabilityEvent, error) {
	var out []interface{}
	err := _Autonity.contract.Call(opts, &out, "getValidatorRecentMisbehaviours", _addr)

	if err != nil {
		return *new([]AutonityAccountabilityEvent), err
	}

	out0 := *abi.ConvertType(out[0], new([]AutonityAccountabilityEvent)).(*[]AutonityAccountabilityEvent)

	return out0, err

}

// GetValidatorRecentMisbehaviours is a free data retrieval call binding the contract method 0xac306841.
//
// Solidity: function getValidatorRecentMisbehaviours(address _addr) view returns((uint8,uint8,uint8,uint8,address,address,bytes32,bytes)[])
func (_Autonity *AutonitySession) GetValidatorRecentMisbehaviours(_addr common.Address) ([]AutonityAccountabilityEvent, error) {
	return _Autonity.Contract.GetValidatorRecentMisbehaviours(&_Autonity.CallOpts, _addr)
}

// GetValidatorRecentMisbehaviours is a free data retrieval call binding the contract method 0xac306841.
//
// Solidity: function getValidatorRecentMisbehaviours(address _addr) view returns((uint8,uint8,uint8,uint8,address,address,bytes32,bytes)[])
func (_Autonity *AutonityCallerSession) GetValidatorRecentMisbehaviours(_addr common.Address) ([]AutonityAccountabilityEvent, error) {
	return _Autonity.Contract.GetValidatorRecentMisbehaviours(&_Autonity.CallOpts, _addr)
}

// GetValidators is a free data retrieval call binding the contract method 0xb7ab4db5.
//
// Solidity: function getValidators() view returns(address[])
func (_Autonity *AutonityCaller) GetValidators(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _Autonity.contract.Call(opts, &out, "getValidators")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// GetValidators is a free data retrieval call binding the contract method 0xb7ab4db5.
//
// Solidity: function getValidators() view returns(address[])
func (_Autonity *AutonitySession) GetValidators() ([]common.Address, error) {
	return _Autonity.Contract.GetValidators(&_Autonity.CallOpts)
}

// GetValidators is a free data retrieval call binding the contract method 0xb7ab4db5.
//
// Solidity: function getValidators() view returns(address[])
func (_Autonity *AutonityCallerSession) GetValidators() ([]common.Address, error) {
	return _Autonity.Contract.GetValidators(&_Autonity.CallOpts)
}

// GetVersion is a free data retrieval call binding the contract method 0x0d8e6e2c.
//
// Solidity: function getVersion() view returns(uint256)
func (_Autonity *AutonityCaller) GetVersion(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Autonity.contract.Call(opts, &out, "getVersion")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetVersion is a free data retrieval call binding the contract method 0x0d8e6e2c.
//
// Solidity: function getVersion() view returns(uint256)
func (_Autonity *AutonitySession) GetVersion() (*big.Int, error) {
	return _Autonity.Contract.GetVersion(&_Autonity.CallOpts)
}

// GetVersion is a free data retrieval call binding the contract method 0x0d8e6e2c.
//
// Solidity: function getVersion() view returns(uint256)
func (_Autonity *AutonityCallerSession) GetVersion() (*big.Int, error) {
	return _Autonity.Contract.GetVersion(&_Autonity.CallOpts)
}

// HeadBondingID is a free data retrieval call binding the contract method 0x44697221.
//
// Solidity: function headBondingID() view returns(uint256)
func (_Autonity *AutonityCaller) HeadBondingID(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Autonity.contract.Call(opts, &out, "headBondingID")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// HeadBondingID is a free data retrieval call binding the contract method 0x44697221.
//
// Solidity: function headBondingID() view returns(uint256)
func (_Autonity *AutonitySession) HeadBondingID() (*big.Int, error) {
	return _Autonity.Contract.HeadBondingID(&_Autonity.CallOpts)
}

// HeadBondingID is a free data retrieval call binding the contract method 0x44697221.
//
// Solidity: function headBondingID() view returns(uint256)
func (_Autonity *AutonityCallerSession) HeadBondingID() (*big.Int, error) {
	return _Autonity.Contract.HeadBondingID(&_Autonity.CallOpts)
}

// HeadUnbondingID is a free data retrieval call binding the contract method 0x4b0dff63.
//
// Solidity: function headUnbondingID() view returns(uint256)
func (_Autonity *AutonityCaller) HeadUnbondingID(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Autonity.contract.Call(opts, &out, "headUnbondingID")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// HeadUnbondingID is a free data retrieval call binding the contract method 0x4b0dff63.
//
// Solidity: function headUnbondingID() view returns(uint256)
func (_Autonity *AutonitySession) HeadUnbondingID() (*big.Int, error) {
	return _Autonity.Contract.HeadUnbondingID(&_Autonity.CallOpts)
}

// HeadUnbondingID is a free data retrieval call binding the contract method 0x4b0dff63.
//
// Solidity: function headUnbondingID() view returns(uint256)
func (_Autonity *AutonityCallerSession) HeadUnbondingID() (*big.Int, error) {
	return _Autonity.Contract.HeadUnbondingID(&_Autonity.CallOpts)
}

// LastEpochBlock is a free data retrieval call binding the contract method 0xc2362dd5.
//
// Solidity: function lastEpochBlock() view returns(uint256)
func (_Autonity *AutonityCaller) LastEpochBlock(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Autonity.contract.Call(opts, &out, "lastEpochBlock")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LastEpochBlock is a free data retrieval call binding the contract method 0xc2362dd5.
//
// Solidity: function lastEpochBlock() view returns(uint256)
func (_Autonity *AutonitySession) LastEpochBlock() (*big.Int, error) {
	return _Autonity.Contract.LastEpochBlock(&_Autonity.CallOpts)
}

// LastEpochBlock is a free data retrieval call binding the contract method 0xc2362dd5.
//
// Solidity: function lastEpochBlock() view returns(uint256)
func (_Autonity *AutonityCallerSession) LastEpochBlock() (*big.Int, error) {
	return _Autonity.Contract.LastEpochBlock(&_Autonity.CallOpts)
}

// MisbehaviourProcessed is a free data retrieval call binding the contract method 0x8a7b7f62.
//
// Solidity: function misbehaviourProcessed(bytes32 _msgHash) view returns(bool)
func (_Autonity *AutonityCaller) MisbehaviourProcessed(opts *bind.CallOpts, _msgHash [32]byte) (bool, error) {
	var out []interface{}
	err := _Autonity.contract.Call(opts, &out, "misbehaviourProcessed", _msgHash)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// MisbehaviourProcessed is a free data retrieval call binding the contract method 0x8a7b7f62.
//
// Solidity: function misbehaviourProcessed(bytes32 _msgHash) view returns(bool)
func (_Autonity *AutonitySession) MisbehaviourProcessed(_msgHash [32]byte) (bool, error) {
	return _Autonity.Contract.MisbehaviourProcessed(&_Autonity.CallOpts, _msgHash)
}

// MisbehaviourProcessed is a free data retrieval call binding the contract method 0x8a7b7f62.
//
// Solidity: function misbehaviourProcessed(bytes32 _msgHash) view returns(bool)
func (_Autonity *AutonityCallerSession) MisbehaviourProcessed(_msgHash [32]byte) (bool, error) {
	return _Autonity.Contract.MisbehaviourProcessed(&_Autonity.CallOpts, _msgHash)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() pure returns(string)
func (_Autonity *AutonityCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Autonity.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() pure returns(string)
func (_Autonity *AutonitySession) Name() (string, error) {
	return _Autonity.Contract.Name(&_Autonity.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() pure returns(string)
func (_Autonity *AutonityCallerSession) Name() (string, error) {
	return _Autonity.Contract.Name(&_Autonity.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() pure returns(string)
func (_Autonity *AutonityCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Autonity.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() pure returns(string)
func (_Autonity *AutonitySession) Symbol() (string, error) {
	return _Autonity.Contract.Symbol(&_Autonity.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() pure returns(string)
func (_Autonity *AutonityCallerSession) Symbol() (string, error) {
	return _Autonity.Contract.Symbol(&_Autonity.CallOpts)
}

// TailBondingID is a free data retrieval call binding the contract method 0x787a2433.
//
// Solidity: function tailBondingID() view returns(uint256)
func (_Autonity *AutonityCaller) TailBondingID(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Autonity.contract.Call(opts, &out, "tailBondingID")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TailBondingID is a free data retrieval call binding the contract method 0x787a2433.
//
// Solidity: function tailBondingID() view returns(uint256)
func (_Autonity *AutonitySession) TailBondingID() (*big.Int, error) {
	return _Autonity.Contract.TailBondingID(&_Autonity.CallOpts)
}

// TailBondingID is a free data retrieval call binding the contract method 0x787a2433.
//
// Solidity: function tailBondingID() view returns(uint256)
func (_Autonity *AutonityCallerSession) TailBondingID() (*big.Int, error) {
	return _Autonity.Contract.TailBondingID(&_Autonity.CallOpts)
}

// TailUnbondingID is a free data retrieval call binding the contract method 0x662cd7f4.
//
// Solidity: function tailUnbondingID() view returns(uint256)
func (_Autonity *AutonityCaller) TailUnbondingID(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Autonity.contract.Call(opts, &out, "tailUnbondingID")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TailUnbondingID is a free data retrieval call binding the contract method 0x662cd7f4.
//
// Solidity: function tailUnbondingID() view returns(uint256)
func (_Autonity *AutonitySession) TailUnbondingID() (*big.Int, error) {
	return _Autonity.Contract.TailUnbondingID(&_Autonity.CallOpts)
}

// TailUnbondingID is a free data retrieval call binding the contract method 0x662cd7f4.
//
// Solidity: function tailUnbondingID() view returns(uint256)
func (_Autonity *AutonityCallerSession) TailUnbondingID() (*big.Int, error) {
	return _Autonity.Contract.TailUnbondingID(&_Autonity.CallOpts)
}

// TotalRedistributed is a free data retrieval call binding the contract method 0x9bb851c0.
//
// Solidity: function totalRedistributed() view returns(uint256)
func (_Autonity *AutonityCaller) TotalRedistributed(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Autonity.contract.Call(opts, &out, "totalRedistributed")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalRedistributed is a free data retrieval call binding the contract method 0x9bb851c0.
//
// Solidity: function totalRedistributed() view returns(uint256)
func (_Autonity *AutonitySession) TotalRedistributed() (*big.Int, error) {
	return _Autonity.Contract.TotalRedistributed(&_Autonity.CallOpts)
}

// TotalRedistributed is a free data retrieval call binding the contract method 0x9bb851c0.
//
// Solidity: function totalRedistributed() view returns(uint256)
func (_Autonity *AutonityCallerSession) TotalRedistributed() (*big.Int, error) {
	return _Autonity.Contract.TotalRedistributed(&_Autonity.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_Autonity *AutonityCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Autonity.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_Autonity *AutonitySession) TotalSupply() (*big.Int, error) {
	return _Autonity.Contract.TotalSupply(&_Autonity.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_Autonity *AutonityCallerSession) TotalSupply() (*big.Int, error) {
	return _Autonity.Contract.TotalSupply(&_Autonity.CallOpts)
}

// ActivateValidator is a paid mutator transaction binding the contract method 0xb46e5520.
//
// Solidity: function activateValidator(address _address) returns()
func (_Autonity *AutonityTransactor) ActivateValidator(opts *bind.TransactOpts, _address common.Address) (*types.Transaction, error) {
	return _Autonity.contract.Transact(opts, "activateValidator", _address)
}

// ActivateValidator is a paid mutator transaction binding the contract method 0xb46e5520.
//
// Solidity: function activateValidator(address _address) returns()
func (_Autonity *AutonitySession) ActivateValidator(_address common.Address) (*types.Transaction, error) {
	return _Autonity.Contract.ActivateValidator(&_Autonity.TransactOpts, _address)
}

// ActivateValidator is a paid mutator transaction binding the contract method 0xb46e5520.
//
// Solidity: function activateValidator(address _address) returns()
func (_Autonity *AutonityTransactorSession) ActivateValidator(_address common.Address) (*types.Transaction, error) {
	return _Autonity.Contract.ActivateValidator(&_Autonity.TransactOpts, _address)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_Autonity *AutonityTransactor) Approve(opts *bind.TransactOpts, spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Autonity.contract.Transact(opts, "approve", spender, amount)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_Autonity *AutonitySession) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Autonity.Contract.Approve(&_Autonity.TransactOpts, spender, amount)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_Autonity *AutonityTransactorSession) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Autonity.Contract.Approve(&_Autonity.TransactOpts, spender, amount)
}

// Bond is a paid mutator transaction binding the contract method 0xa515366a.
//
// Solidity: function bond(address _validator, uint256 _amount) returns()
func (_Autonity *AutonityTransactor) Bond(opts *bind.TransactOpts, _validator common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Autonity.contract.Transact(opts, "bond", _validator, _amount)
}

// Bond is a paid mutator transaction binding the contract method 0xa515366a.
//
// Solidity: function bond(address _validator, uint256 _amount) returns()
func (_Autonity *AutonitySession) Bond(_validator common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Autonity.Contract.Bond(&_Autonity.TransactOpts, _validator, _amount)
}

// Bond is a paid mutator transaction binding the contract method 0xa515366a.
//
// Solidity: function bond(address _validator, uint256 _amount) returns()
func (_Autonity *AutonityTransactorSession) Bond(_validator common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Autonity.Contract.Bond(&_Autonity.TransactOpts, _validator, _amount)
}

// Burn is a paid mutator transaction binding the contract method 0x9dc29fac.
//
// Solidity: function burn(address _addr, uint256 _amount) returns()
func (_Autonity *AutonityTransactor) Burn(opts *bind.TransactOpts, _addr common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Autonity.contract.Transact(opts, "burn", _addr, _amount)
}

// Burn is a paid mutator transaction binding the contract method 0x9dc29fac.
//
// Solidity: function burn(address _addr, uint256 _amount) returns()
func (_Autonity *AutonitySession) Burn(_addr common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Autonity.Contract.Burn(&_Autonity.TransactOpts, _addr, _amount)
}

// Burn is a paid mutator transaction binding the contract method 0x9dc29fac.
//
// Solidity: function burn(address _addr, uint256 _amount) returns()
func (_Autonity *AutonityTransactorSession) Burn(_addr common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Autonity.Contract.Burn(&_Autonity.TransactOpts, _addr, _amount)
}

// ChangeCommissionRate is a paid mutator transaction binding the contract method 0x852c4849.
//
// Solidity: function changeCommissionRate(address _validator, uint256 _rate) returns()
func (_Autonity *AutonityTransactor) ChangeCommissionRate(opts *bind.TransactOpts, _validator common.Address, _rate *big.Int) (*types.Transaction, error) {
	return _Autonity.contract.Transact(opts, "changeCommissionRate", _validator, _rate)
}

// ChangeCommissionRate is a paid mutator transaction binding the contract method 0x852c4849.
//
// Solidity: function changeCommissionRate(address _validator, uint256 _rate) returns()
func (_Autonity *AutonitySession) ChangeCommissionRate(_validator common.Address, _rate *big.Int) (*types.Transaction, error) {
	return _Autonity.Contract.ChangeCommissionRate(&_Autonity.TransactOpts, _validator, _rate)
}

// ChangeCommissionRate is a paid mutator transaction binding the contract method 0x852c4849.
//
// Solidity: function changeCommissionRate(address _validator, uint256 _rate) returns()
func (_Autonity *AutonityTransactorSession) ChangeCommissionRate(_validator common.Address, _rate *big.Int) (*types.Transaction, error) {
	return _Autonity.Contract.ChangeCommissionRate(&_Autonity.TransactOpts, _validator, _rate)
}

// CompleteContractUpgrade is a paid mutator transaction binding the contract method 0x872cf059.
//
// Solidity: function completeContractUpgrade() returns()
func (_Autonity *AutonityTransactor) CompleteContractUpgrade(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Autonity.contract.Transact(opts, "completeContractUpgrade")
}

// CompleteContractUpgrade is a paid mutator transaction binding the contract method 0x872cf059.
//
// Solidity: function completeContractUpgrade() returns()
func (_Autonity *AutonitySession) CompleteContractUpgrade() (*types.Transaction, error) {
	return _Autonity.Contract.CompleteContractUpgrade(&_Autonity.TransactOpts)
}

// CompleteContractUpgrade is a paid mutator transaction binding the contract method 0x872cf059.
//
// Solidity: function completeContractUpgrade() returns()
func (_Autonity *AutonityTransactorSession) CompleteContractUpgrade() (*types.Transaction, error) {
	return _Autonity.Contract.CompleteContractUpgrade(&_Autonity.TransactOpts)
}

// ComputeCommittee is a paid mutator transaction binding the contract method 0xae1f5fa0.
//
// Solidity: function computeCommittee() returns(address[])
func (_Autonity *AutonityTransactor) ComputeCommittee(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Autonity.contract.Transact(opts, "computeCommittee")
}

// ComputeCommittee is a paid mutator transaction binding the contract method 0xae1f5fa0.
//
// Solidity: function computeCommittee() returns(address[])
func (_Autonity *AutonitySession) ComputeCommittee() (*types.Transaction, error) {
	return _Autonity.Contract.ComputeCommittee(&_Autonity.TransactOpts)
}

// ComputeCommittee is a paid mutator transaction binding the contract method 0xae1f5fa0.
//
// Solidity: function computeCommittee() returns(address[])
func (_Autonity *AutonityTransactorSession) ComputeCommittee() (*types.Transaction, error) {
	return _Autonity.Contract.ComputeCommittee(&_Autonity.TransactOpts)
}

// Finalize is a paid mutator transaction binding the contract method 0x4bb278f3.
//
// Solidity: function finalize() returns(bool, (address,uint256)[])
func (_Autonity *AutonityTransactor) Finalize(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Autonity.contract.Transact(opts, "finalize")
}

// Finalize is a paid mutator transaction binding the contract method 0x4bb278f3.
//
// Solidity: function finalize() returns(bool, (address,uint256)[])
func (_Autonity *AutonitySession) Finalize() (*types.Transaction, error) {
	return _Autonity.Contract.Finalize(&_Autonity.TransactOpts)
}

// Finalize is a paid mutator transaction binding the contract method 0x4bb278f3.
//
// Solidity: function finalize() returns(bool, (address,uint256)[])
func (_Autonity *AutonityTransactorSession) Finalize() (*types.Transaction, error) {
	return _Autonity.Contract.Finalize(&_Autonity.TransactOpts)
}

// HandleAccountabilityEvents is a paid mutator transaction binding the contract method 0x3a17914f.
//
// Solidity: function handleAccountabilityEvents((uint8,uint8,uint8,uint8,address,address,bytes32,bytes)[] _events) returns()
func (_Autonity *AutonityTransactor) HandleAccountabilityEvents(opts *bind.TransactOpts, _events []AutonityAccountabilityEvent) (*types.Transaction, error) {
	return _Autonity.contract.Transact(opts, "handleAccountabilityEvents", _events)
}

// HandleAccountabilityEvents is a paid mutator transaction binding the contract method 0x3a17914f.
//
// Solidity: function handleAccountabilityEvents((uint8,uint8,uint8,uint8,address,address,bytes32,bytes)[] _events) returns()
func (_Autonity *AutonitySession) HandleAccountabilityEvents(_events []AutonityAccountabilityEvent) (*types.Transaction, error) {
	return _Autonity.Contract.HandleAccountabilityEvents(&_Autonity.TransactOpts, _events)
}

// HandleAccountabilityEvents is a paid mutator transaction binding the contract method 0x3a17914f.
//
// Solidity: function handleAccountabilityEvents((uint8,uint8,uint8,uint8,address,address,bytes32,bytes)[] _events) returns()
func (_Autonity *AutonityTransactorSession) HandleAccountabilityEvents(_events []AutonityAccountabilityEvent) (*types.Transaction, error) {
	return _Autonity.Contract.HandleAccountabilityEvents(&_Autonity.TransactOpts, _events)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address _addr, uint256 _amount) returns()
func (_Autonity *AutonityTransactor) Mint(opts *bind.TransactOpts, _addr common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Autonity.contract.Transact(opts, "mint", _addr, _amount)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address _addr, uint256 _amount) returns()
func (_Autonity *AutonitySession) Mint(_addr common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Autonity.Contract.Mint(&_Autonity.TransactOpts, _addr, _amount)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address _addr, uint256 _amount) returns()
func (_Autonity *AutonityTransactorSession) Mint(_addr common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Autonity.Contract.Mint(&_Autonity.TransactOpts, _addr, _amount)
}

// PauseValidator is a paid mutator transaction binding the contract method 0x0ae65e7a.
//
// Solidity: function pauseValidator(address _address) returns()
func (_Autonity *AutonityTransactor) PauseValidator(opts *bind.TransactOpts, _address common.Address) (*types.Transaction, error) {
	return _Autonity.contract.Transact(opts, "pauseValidator", _address)
}

// PauseValidator is a paid mutator transaction binding the contract method 0x0ae65e7a.
//
// Solidity: function pauseValidator(address _address) returns()
func (_Autonity *AutonitySession) PauseValidator(_address common.Address) (*types.Transaction, error) {
	return _Autonity.Contract.PauseValidator(&_Autonity.TransactOpts, _address)
}

// PauseValidator is a paid mutator transaction binding the contract method 0x0ae65e7a.
//
// Solidity: function pauseValidator(address _address) returns()
func (_Autonity *AutonityTransactorSession) PauseValidator(_address common.Address) (*types.Transaction, error) {
	return _Autonity.Contract.PauseValidator(&_Autonity.TransactOpts, _address)
}

// RegisterValidator is a paid mutator transaction binding the contract method 0xad722d4d.
//
// Solidity: function registerValidator(string _enode, address _oracleAddress, bytes _multisig) returns()
func (_Autonity *AutonityTransactor) RegisterValidator(opts *bind.TransactOpts, _enode string, _oracleAddress common.Address, _multisig []byte) (*types.Transaction, error) {
	return _Autonity.contract.Transact(opts, "registerValidator", _enode, _oracleAddress, _multisig)
}

// RegisterValidator is a paid mutator transaction binding the contract method 0xad722d4d.
//
// Solidity: function registerValidator(string _enode, address _oracleAddress, bytes _multisig) returns()
func (_Autonity *AutonitySession) RegisterValidator(_enode string, _oracleAddress common.Address, _multisig []byte) (*types.Transaction, error) {
	return _Autonity.Contract.RegisterValidator(&_Autonity.TransactOpts, _enode, _oracleAddress, _multisig)
}

// RegisterValidator is a paid mutator transaction binding the contract method 0xad722d4d.
//
// Solidity: function registerValidator(string _enode, address _oracleAddress, bytes _multisig) returns()
func (_Autonity *AutonityTransactorSession) RegisterValidator(_enode string, _oracleAddress common.Address, _multisig []byte) (*types.Transaction, error) {
	return _Autonity.Contract.RegisterValidator(&_Autonity.TransactOpts, _enode, _oracleAddress, _multisig)
}

// ResetContractUpgrade is a paid mutator transaction binding the contract method 0xcf9c5719.
//
// Solidity: function resetContractUpgrade() returns()
func (_Autonity *AutonityTransactor) ResetContractUpgrade(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Autonity.contract.Transact(opts, "resetContractUpgrade")
}

// ResetContractUpgrade is a paid mutator transaction binding the contract method 0xcf9c5719.
//
// Solidity: function resetContractUpgrade() returns()
func (_Autonity *AutonitySession) ResetContractUpgrade() (*types.Transaction, error) {
	return _Autonity.Contract.ResetContractUpgrade(&_Autonity.TransactOpts)
}

// ResetContractUpgrade is a paid mutator transaction binding the contract method 0xcf9c5719.
//
// Solidity: function resetContractUpgrade() returns()
func (_Autonity *AutonityTransactorSession) ResetContractUpgrade() (*types.Transaction, error) {
	return _Autonity.Contract.ResetContractUpgrade(&_Autonity.TransactOpts)
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

// SetEpochPeriod is a paid mutator transaction binding the contract method 0x6b5f444c.
//
// Solidity: function setEpochPeriod(uint256 _period) returns()
func (_Autonity *AutonityTransactor) SetEpochPeriod(opts *bind.TransactOpts, _period *big.Int) (*types.Transaction, error) {
	return _Autonity.contract.Transact(opts, "setEpochPeriod", _period)
}

// SetEpochPeriod is a paid mutator transaction binding the contract method 0x6b5f444c.
//
// Solidity: function setEpochPeriod(uint256 _period) returns()
func (_Autonity *AutonitySession) SetEpochPeriod(_period *big.Int) (*types.Transaction, error) {
	return _Autonity.Contract.SetEpochPeriod(&_Autonity.TransactOpts, _period)
}

// SetEpochPeriod is a paid mutator transaction binding the contract method 0x6b5f444c.
//
// Solidity: function setEpochPeriod(uint256 _period) returns()
func (_Autonity *AutonityTransactorSession) SetEpochPeriod(_period *big.Int) (*types.Transaction, error) {
	return _Autonity.Contract.SetEpochPeriod(&_Autonity.TransactOpts, _period)
}

// SetMinimumBaseFee is a paid mutator transaction binding the contract method 0xcb696f54.
//
// Solidity: function setMinimumBaseFee(uint256 _price) returns()
func (_Autonity *AutonityTransactor) SetMinimumBaseFee(opts *bind.TransactOpts, _price *big.Int) (*types.Transaction, error) {
	return _Autonity.contract.Transact(opts, "setMinimumBaseFee", _price)
}

// SetMinimumBaseFee is a paid mutator transaction binding the contract method 0xcb696f54.
//
// Solidity: function setMinimumBaseFee(uint256 _price) returns()
func (_Autonity *AutonitySession) SetMinimumBaseFee(_price *big.Int) (*types.Transaction, error) {
	return _Autonity.Contract.SetMinimumBaseFee(&_Autonity.TransactOpts, _price)
}

// SetMinimumBaseFee is a paid mutator transaction binding the contract method 0xcb696f54.
//
// Solidity: function setMinimumBaseFee(uint256 _price) returns()
func (_Autonity *AutonityTransactorSession) SetMinimumBaseFee(_price *big.Int) (*types.Transaction, error) {
	return _Autonity.Contract.SetMinimumBaseFee(&_Autonity.TransactOpts, _price)
}

// SetMisbehaviourPenalty is a paid mutator transaction binding the contract method 0x8f5d0fcb.
//
// Solidity: function setMisbehaviourPenalty(uint256 _newPenalty) returns()
func (_Autonity *AutonityTransactor) SetMisbehaviourPenalty(opts *bind.TransactOpts, _newPenalty *big.Int) (*types.Transaction, error) {
	return _Autonity.contract.Transact(opts, "setMisbehaviourPenalty", _newPenalty)
}

// SetMisbehaviourPenalty is a paid mutator transaction binding the contract method 0x8f5d0fcb.
//
// Solidity: function setMisbehaviourPenalty(uint256 _newPenalty) returns()
func (_Autonity *AutonitySession) SetMisbehaviourPenalty(_newPenalty *big.Int) (*types.Transaction, error) {
	return _Autonity.Contract.SetMisbehaviourPenalty(&_Autonity.TransactOpts, _newPenalty)
}

// SetMisbehaviourPenalty is a paid mutator transaction binding the contract method 0x8f5d0fcb.
//
// Solidity: function setMisbehaviourPenalty(uint256 _newPenalty) returns()
func (_Autonity *AutonityTransactorSession) SetMisbehaviourPenalty(_newPenalty *big.Int) (*types.Transaction, error) {
	return _Autonity.Contract.SetMisbehaviourPenalty(&_Autonity.TransactOpts, _newPenalty)
}

// SetOperatorAccount is a paid mutator transaction binding the contract method 0x520fdbbc.
//
// Solidity: function setOperatorAccount(address _account) returns()
func (_Autonity *AutonityTransactor) SetOperatorAccount(opts *bind.TransactOpts, _account common.Address) (*types.Transaction, error) {
	return _Autonity.contract.Transact(opts, "setOperatorAccount", _account)
}

// SetOperatorAccount is a paid mutator transaction binding the contract method 0x520fdbbc.
//
// Solidity: function setOperatorAccount(address _account) returns()
func (_Autonity *AutonitySession) SetOperatorAccount(_account common.Address) (*types.Transaction, error) {
	return _Autonity.Contract.SetOperatorAccount(&_Autonity.TransactOpts, _account)
}

// SetOperatorAccount is a paid mutator transaction binding the contract method 0x520fdbbc.
//
// Solidity: function setOperatorAccount(address _account) returns()
func (_Autonity *AutonityTransactorSession) SetOperatorAccount(_account common.Address) (*types.Transaction, error) {
	return _Autonity.Contract.SetOperatorAccount(&_Autonity.TransactOpts, _account)
}

// SetTreasuryAccount is a paid mutator transaction binding the contract method 0xd886f8a2.
//
// Solidity: function setTreasuryAccount(address _account) returns()
func (_Autonity *AutonityTransactor) SetTreasuryAccount(opts *bind.TransactOpts, _account common.Address) (*types.Transaction, error) {
	return _Autonity.contract.Transact(opts, "setTreasuryAccount", _account)
}

// SetTreasuryAccount is a paid mutator transaction binding the contract method 0xd886f8a2.
//
// Solidity: function setTreasuryAccount(address _account) returns()
func (_Autonity *AutonitySession) SetTreasuryAccount(_account common.Address) (*types.Transaction, error) {
	return _Autonity.Contract.SetTreasuryAccount(&_Autonity.TransactOpts, _account)
}

// SetTreasuryAccount is a paid mutator transaction binding the contract method 0xd886f8a2.
//
// Solidity: function setTreasuryAccount(address _account) returns()
func (_Autonity *AutonityTransactorSession) SetTreasuryAccount(_account common.Address) (*types.Transaction, error) {
	return _Autonity.Contract.SetTreasuryAccount(&_Autonity.TransactOpts, _account)
}

// SetTreasuryFee is a paid mutator transaction binding the contract method 0x77e741c7.
//
// Solidity: function setTreasuryFee(uint256 _treasuryFee) returns()
func (_Autonity *AutonityTransactor) SetTreasuryFee(opts *bind.TransactOpts, _treasuryFee *big.Int) (*types.Transaction, error) {
	return _Autonity.contract.Transact(opts, "setTreasuryFee", _treasuryFee)
}

// SetTreasuryFee is a paid mutator transaction binding the contract method 0x77e741c7.
//
// Solidity: function setTreasuryFee(uint256 _treasuryFee) returns()
func (_Autonity *AutonitySession) SetTreasuryFee(_treasuryFee *big.Int) (*types.Transaction, error) {
	return _Autonity.Contract.SetTreasuryFee(&_Autonity.TransactOpts, _treasuryFee)
}

// SetTreasuryFee is a paid mutator transaction binding the contract method 0x77e741c7.
//
// Solidity: function setTreasuryFee(uint256 _treasuryFee) returns()
func (_Autonity *AutonityTransactorSession) SetTreasuryFee(_treasuryFee *big.Int) (*types.Transaction, error) {
	return _Autonity.Contract.SetTreasuryFee(&_Autonity.TransactOpts, _treasuryFee)
}

// SetUnbondingPeriod is a paid mutator transaction binding the contract method 0x114eaf55.
//
// Solidity: function setUnbondingPeriod(uint256 _period) returns()
func (_Autonity *AutonityTransactor) SetUnbondingPeriod(opts *bind.TransactOpts, _period *big.Int) (*types.Transaction, error) {
	return _Autonity.contract.Transact(opts, "setUnbondingPeriod", _period)
}

// SetUnbondingPeriod is a paid mutator transaction binding the contract method 0x114eaf55.
//
// Solidity: function setUnbondingPeriod(uint256 _period) returns()
func (_Autonity *AutonitySession) SetUnbondingPeriod(_period *big.Int) (*types.Transaction, error) {
	return _Autonity.Contract.SetUnbondingPeriod(&_Autonity.TransactOpts, _period)
}

// SetUnbondingPeriod is a paid mutator transaction binding the contract method 0x114eaf55.
//
// Solidity: function setUnbondingPeriod(uint256 _period) returns()
func (_Autonity *AutonityTransactorSession) SetUnbondingPeriod(_period *big.Int) (*types.Transaction, error) {
	return _Autonity.Contract.SetUnbondingPeriod(&_Autonity.TransactOpts, _period)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address _recipient, uint256 _amount) returns(bool)
func (_Autonity *AutonityTransactor) Transfer(opts *bind.TransactOpts, _recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Autonity.contract.Transact(opts, "transfer", _recipient, _amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address _recipient, uint256 _amount) returns(bool)
func (_Autonity *AutonitySession) Transfer(_recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Autonity.Contract.Transfer(&_Autonity.TransactOpts, _recipient, _amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address _recipient, uint256 _amount) returns(bool)
func (_Autonity *AutonityTransactorSession) Transfer(_recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Autonity.Contract.Transfer(&_Autonity.TransactOpts, _recipient, _amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_Autonity *AutonityTransactor) TransferFrom(opts *bind.TransactOpts, sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Autonity.contract.Transact(opts, "transferFrom", sender, recipient, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_Autonity *AutonitySession) TransferFrom(sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Autonity.Contract.TransferFrom(&_Autonity.TransactOpts, sender, recipient, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_Autonity *AutonityTransactorSession) TransferFrom(sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Autonity.Contract.TransferFrom(&_Autonity.TransactOpts, sender, recipient, amount)
}

// Unbond is a paid mutator transaction binding the contract method 0xa5d059ca.
//
// Solidity: function unbond(address _validator, uint256 _amount) returns()
func (_Autonity *AutonityTransactor) Unbond(opts *bind.TransactOpts, _validator common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Autonity.contract.Transact(opts, "unbond", _validator, _amount)
}

// Unbond is a paid mutator transaction binding the contract method 0xa5d059ca.
//
// Solidity: function unbond(address _validator, uint256 _amount) returns()
func (_Autonity *AutonitySession) Unbond(_validator common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Autonity.Contract.Unbond(&_Autonity.TransactOpts, _validator, _amount)
}

// Unbond is a paid mutator transaction binding the contract method 0xa5d059ca.
//
// Solidity: function unbond(address _validator, uint256 _amount) returns()
func (_Autonity *AutonityTransactorSession) Unbond(_validator common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Autonity.Contract.Unbond(&_Autonity.TransactOpts, _validator, _amount)
}

// UpgradeContract is a paid mutator transaction binding the contract method 0xb2ea9adb.
//
// Solidity: function upgradeContract(bytes _bytecode, string _abi) returns()
func (_Autonity *AutonityTransactor) UpgradeContract(opts *bind.TransactOpts, _bytecode []byte, _abi string) (*types.Transaction, error) {
	return _Autonity.contract.Transact(opts, "upgradeContract", _bytecode, _abi)
}

// UpgradeContract is a paid mutator transaction binding the contract method 0xb2ea9adb.
//
// Solidity: function upgradeContract(bytes _bytecode, string _abi) returns()
func (_Autonity *AutonitySession) UpgradeContract(_bytecode []byte, _abi string) (*types.Transaction, error) {
	return _Autonity.Contract.UpgradeContract(&_Autonity.TransactOpts, _bytecode, _abi)
}

// UpgradeContract is a paid mutator transaction binding the contract method 0xb2ea9adb.
//
// Solidity: function upgradeContract(bytes _bytecode, string _abi) returns()
func (_Autonity *AutonityTransactorSession) UpgradeContract(_bytecode []byte, _abi string) (*types.Transaction, error) {
	return _Autonity.Contract.UpgradeContract(&_Autonity.TransactOpts, _bytecode, _abi)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Autonity *AutonityTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _Autonity.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Autonity *AutonitySession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _Autonity.Contract.Fallback(&_Autonity.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Autonity *AutonityTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _Autonity.Contract.Fallback(&_Autonity.TransactOpts, calldata)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Autonity *AutonityTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Autonity.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Autonity *AutonitySession) Receive() (*types.Transaction, error) {
	return _Autonity.Contract.Receive(&_Autonity.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Autonity *AutonityTransactorSession) Receive() (*types.Transaction, error) {
	return _Autonity.Contract.Receive(&_Autonity.TransactOpts)
}

// AutonityAccusationAddedIterator is returned from FilterAccusationAdded and is used to iterate over the raw logs and unpacked data for AccusationAdded events raised by the Autonity contract.
type AutonityAccusationAddedIterator struct {
	Event *AutonityAccusationAdded // Event containing the contract specifics and raw log

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
func (it *AutonityAccusationAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutonityAccusationAdded)
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
		it.Event = new(AutonityAccusationAdded)
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
func (it *AutonityAccusationAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AutonityAccusationAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AutonityAccusationAdded represents a AccusationAdded event raised by the Autonity contract.
type AutonityAccusationAdded struct {
	Ev  AutonityAccountabilityEvent
	Raw types.Log // Blockchain specific contextual infos
}

// FilterAccusationAdded is a free log retrieval operation binding the contract event 0x244ffefead78aaef5913a3abac1c8477dec686bf017a343f1679f9c8b6a77f11.
//
// Solidity: event AccusationAdded((uint8,uint8,uint8,uint8,address,address,bytes32,bytes) ev)
func (_Autonity *AutonityFilterer) FilterAccusationAdded(opts *bind.FilterOpts) (*AutonityAccusationAddedIterator, error) {

	logs, sub, err := _Autonity.contract.FilterLogs(opts, "AccusationAdded")
	if err != nil {
		return nil, err
	}
	return &AutonityAccusationAddedIterator{contract: _Autonity.contract, event: "AccusationAdded", logs: logs, sub: sub}, nil
}

// WatchAccusationAdded is a free log subscription operation binding the contract event 0x244ffefead78aaef5913a3abac1c8477dec686bf017a343f1679f9c8b6a77f11.
//
// Solidity: event AccusationAdded((uint8,uint8,uint8,uint8,address,address,bytes32,bytes) ev)
func (_Autonity *AutonityFilterer) WatchAccusationAdded(opts *bind.WatchOpts, sink chan<- *AutonityAccusationAdded) (event.Subscription, error) {

	logs, sub, err := _Autonity.contract.WatchLogs(opts, "AccusationAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AutonityAccusationAdded)
				if err := _Autonity.contract.UnpackLog(event, "AccusationAdded", log); err != nil {
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

// ParseAccusationAdded is a log parse operation binding the contract event 0x244ffefead78aaef5913a3abac1c8477dec686bf017a343f1679f9c8b6a77f11.
//
// Solidity: event AccusationAdded((uint8,uint8,uint8,uint8,address,address,bytes32,bytes) ev)
func (_Autonity *AutonityFilterer) ParseAccusationAdded(log types.Log) (*AutonityAccusationAdded, error) {
	event := new(AutonityAccusationAdded)
	if err := _Autonity.contract.UnpackLog(event, "AccusationAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AutonityAccusationRemovedIterator is returned from FilterAccusationRemoved and is used to iterate over the raw logs and unpacked data for AccusationRemoved events raised by the Autonity contract.
type AutonityAccusationRemovedIterator struct {
	Event *AutonityAccusationRemoved // Event containing the contract specifics and raw log

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
func (it *AutonityAccusationRemovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutonityAccusationRemoved)
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
		it.Event = new(AutonityAccusationRemoved)
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
func (it *AutonityAccusationRemovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AutonityAccusationRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AutonityAccusationRemoved represents a AccusationRemoved event raised by the Autonity contract.
type AutonityAccusationRemoved struct {
	Ev  AutonityAccountabilityEvent
	Raw types.Log // Blockchain specific contextual infos
}

// FilterAccusationRemoved is a free log retrieval operation binding the contract event 0x663327acde77befae0ec3fed52a32b993673702ddd832552e402d1afbd32158c.
//
// Solidity: event AccusationRemoved((uint8,uint8,uint8,uint8,address,address,bytes32,bytes) ev)
func (_Autonity *AutonityFilterer) FilterAccusationRemoved(opts *bind.FilterOpts) (*AutonityAccusationRemovedIterator, error) {

	logs, sub, err := _Autonity.contract.FilterLogs(opts, "AccusationRemoved")
	if err != nil {
		return nil, err
	}
	return &AutonityAccusationRemovedIterator{contract: _Autonity.contract, event: "AccusationRemoved", logs: logs, sub: sub}, nil
}

// WatchAccusationRemoved is a free log subscription operation binding the contract event 0x663327acde77befae0ec3fed52a32b993673702ddd832552e402d1afbd32158c.
//
// Solidity: event AccusationRemoved((uint8,uint8,uint8,uint8,address,address,bytes32,bytes) ev)
func (_Autonity *AutonityFilterer) WatchAccusationRemoved(opts *bind.WatchOpts, sink chan<- *AutonityAccusationRemoved) (event.Subscription, error) {

	logs, sub, err := _Autonity.contract.WatchLogs(opts, "AccusationRemoved")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AutonityAccusationRemoved)
				if err := _Autonity.contract.UnpackLog(event, "AccusationRemoved", log); err != nil {
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

// ParseAccusationRemoved is a log parse operation binding the contract event 0x663327acde77befae0ec3fed52a32b993673702ddd832552e402d1afbd32158c.
//
// Solidity: event AccusationRemoved((uint8,uint8,uint8,uint8,address,address,bytes32,bytes) ev)
func (_Autonity *AutonityFilterer) ParseAccusationRemoved(log types.Log) (*AutonityAccusationRemoved, error) {
	event := new(AutonityAccusationRemoved)
	if err := _Autonity.contract.UnpackLog(event, "AccusationRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AutonityApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the Autonity contract.
type AutonityApprovalIterator struct {
	Event *AutonityApproval // Event containing the contract specifics and raw log

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
func (it *AutonityApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutonityApproval)
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
		it.Event = new(AutonityApproval)
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
func (it *AutonityApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AutonityApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AutonityApproval represents a Approval event raised by the Autonity contract.
type AutonityApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_Autonity *AutonityFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*AutonityApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _Autonity.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &AutonityApprovalIterator{contract: _Autonity.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_Autonity *AutonityFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *AutonityApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _Autonity.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AutonityApproval)
				if err := _Autonity.contract.UnpackLog(event, "Approval", log); err != nil {
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
func (_Autonity *AutonityFilterer) ParseApproval(log types.Log) (*AutonityApproval, error) {
	event := new(AutonityApproval)
	if err := _Autonity.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AutonityBurnedStakeIterator is returned from FilterBurnedStake and is used to iterate over the raw logs and unpacked data for BurnedStake events raised by the Autonity contract.
type AutonityBurnedStakeIterator struct {
	Event *AutonityBurnedStake // Event containing the contract specifics and raw log

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
func (it *AutonityBurnedStakeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutonityBurnedStake)
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
		it.Event = new(AutonityBurnedStake)
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
func (it *AutonityBurnedStakeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AutonityBurnedStakeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AutonityBurnedStake represents a BurnedStake event raised by the Autonity contract.
type AutonityBurnedStake struct {
	Addr   common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterBurnedStake is a free log retrieval operation binding the contract event 0x5024dbeedf0c06664c9bd7be836915730c955e936972c020683dadf11d5488a3.
//
// Solidity: event BurnedStake(address addr, uint256 amount)
func (_Autonity *AutonityFilterer) FilterBurnedStake(opts *bind.FilterOpts) (*AutonityBurnedStakeIterator, error) {

	logs, sub, err := _Autonity.contract.FilterLogs(opts, "BurnedStake")
	if err != nil {
		return nil, err
	}
	return &AutonityBurnedStakeIterator{contract: _Autonity.contract, event: "BurnedStake", logs: logs, sub: sub}, nil
}

// WatchBurnedStake is a free log subscription operation binding the contract event 0x5024dbeedf0c06664c9bd7be836915730c955e936972c020683dadf11d5488a3.
//
// Solidity: event BurnedStake(address addr, uint256 amount)
func (_Autonity *AutonityFilterer) WatchBurnedStake(opts *bind.WatchOpts, sink chan<- *AutonityBurnedStake) (event.Subscription, error) {

	logs, sub, err := _Autonity.contract.WatchLogs(opts, "BurnedStake")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AutonityBurnedStake)
				if err := _Autonity.contract.UnpackLog(event, "BurnedStake", log); err != nil {
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
// Solidity: event BurnedStake(address addr, uint256 amount)
func (_Autonity *AutonityFilterer) ParseBurnedStake(log types.Log) (*AutonityBurnedStake, error) {
	event := new(AutonityBurnedStake)
	if err := _Autonity.contract.UnpackLog(event, "BurnedStake", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AutonityCommissionRateChangeIterator is returned from FilterCommissionRateChange and is used to iterate over the raw logs and unpacked data for CommissionRateChange events raised by the Autonity contract.
type AutonityCommissionRateChangeIterator struct {
	Event *AutonityCommissionRateChange // Event containing the contract specifics and raw log

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
func (it *AutonityCommissionRateChangeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutonityCommissionRateChange)
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
		it.Event = new(AutonityCommissionRateChange)
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
func (it *AutonityCommissionRateChangeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AutonityCommissionRateChangeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AutonityCommissionRateChange represents a CommissionRateChange event raised by the Autonity contract.
type AutonityCommissionRateChange struct {
	Validator common.Address
	Rate      *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterCommissionRateChange is a free log retrieval operation binding the contract event 0x4fba51c92fa3d6ad8374d394f6cd5766857552e153d7384a8f23aa4ce9a8a7cf.
//
// Solidity: event CommissionRateChange(address validator, uint256 rate)
func (_Autonity *AutonityFilterer) FilterCommissionRateChange(opts *bind.FilterOpts) (*AutonityCommissionRateChangeIterator, error) {

	logs, sub, err := _Autonity.contract.FilterLogs(opts, "CommissionRateChange")
	if err != nil {
		return nil, err
	}
	return &AutonityCommissionRateChangeIterator{contract: _Autonity.contract, event: "CommissionRateChange", logs: logs, sub: sub}, nil
}

// WatchCommissionRateChange is a free log subscription operation binding the contract event 0x4fba51c92fa3d6ad8374d394f6cd5766857552e153d7384a8f23aa4ce9a8a7cf.
//
// Solidity: event CommissionRateChange(address validator, uint256 rate)
func (_Autonity *AutonityFilterer) WatchCommissionRateChange(opts *bind.WatchOpts, sink chan<- *AutonityCommissionRateChange) (event.Subscription, error) {

	logs, sub, err := _Autonity.contract.WatchLogs(opts, "CommissionRateChange")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AutonityCommissionRateChange)
				if err := _Autonity.contract.UnpackLog(event, "CommissionRateChange", log); err != nil {
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
// Solidity: event CommissionRateChange(address validator, uint256 rate)
func (_Autonity *AutonityFilterer) ParseCommissionRateChange(log types.Log) (*AutonityCommissionRateChange, error) {
	event := new(AutonityCommissionRateChange)
	if err := _Autonity.contract.UnpackLog(event, "CommissionRateChange", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AutonityEpochPeriodUpdatedIterator is returned from FilterEpochPeriodUpdated and is used to iterate over the raw logs and unpacked data for EpochPeriodUpdated events raised by the Autonity contract.
type AutonityEpochPeriodUpdatedIterator struct {
	Event *AutonityEpochPeriodUpdated // Event containing the contract specifics and raw log

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
func (it *AutonityEpochPeriodUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutonityEpochPeriodUpdated)
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
		it.Event = new(AutonityEpochPeriodUpdated)
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
func (it *AutonityEpochPeriodUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AutonityEpochPeriodUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AutonityEpochPeriodUpdated represents a EpochPeriodUpdated event raised by the Autonity contract.
type AutonityEpochPeriodUpdated struct {
	Period *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterEpochPeriodUpdated is a free log retrieval operation binding the contract event 0xd7f1279ded354dbf22a69fcc2fd661763a6e2956a5d2891af9410af880fa5f81.
//
// Solidity: event EpochPeriodUpdated(uint256 period)
func (_Autonity *AutonityFilterer) FilterEpochPeriodUpdated(opts *bind.FilterOpts) (*AutonityEpochPeriodUpdatedIterator, error) {

	logs, sub, err := _Autonity.contract.FilterLogs(opts, "EpochPeriodUpdated")
	if err != nil {
		return nil, err
	}
	return &AutonityEpochPeriodUpdatedIterator{contract: _Autonity.contract, event: "EpochPeriodUpdated", logs: logs, sub: sub}, nil
}

// WatchEpochPeriodUpdated is a free log subscription operation binding the contract event 0xd7f1279ded354dbf22a69fcc2fd661763a6e2956a5d2891af9410af880fa5f81.
//
// Solidity: event EpochPeriodUpdated(uint256 period)
func (_Autonity *AutonityFilterer) WatchEpochPeriodUpdated(opts *bind.WatchOpts, sink chan<- *AutonityEpochPeriodUpdated) (event.Subscription, error) {

	logs, sub, err := _Autonity.contract.WatchLogs(opts, "EpochPeriodUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AutonityEpochPeriodUpdated)
				if err := _Autonity.contract.UnpackLog(event, "EpochPeriodUpdated", log); err != nil {
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

// ParseEpochPeriodUpdated is a log parse operation binding the contract event 0xd7f1279ded354dbf22a69fcc2fd661763a6e2956a5d2891af9410af880fa5f81.
//
// Solidity: event EpochPeriodUpdated(uint256 period)
func (_Autonity *AutonityFilterer) ParseEpochPeriodUpdated(log types.Log) (*AutonityEpochPeriodUpdated, error) {
	event := new(AutonityEpochPeriodUpdated)
	if err := _Autonity.contract.UnpackLog(event, "EpochPeriodUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AutonityMinimumBaseFeeUpdatedIterator is returned from FilterMinimumBaseFeeUpdated and is used to iterate over the raw logs and unpacked data for MinimumBaseFeeUpdated events raised by the Autonity contract.
type AutonityMinimumBaseFeeUpdatedIterator struct {
	Event *AutonityMinimumBaseFeeUpdated // Event containing the contract specifics and raw log

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
func (it *AutonityMinimumBaseFeeUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutonityMinimumBaseFeeUpdated)
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
		it.Event = new(AutonityMinimumBaseFeeUpdated)
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
func (it *AutonityMinimumBaseFeeUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AutonityMinimumBaseFeeUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AutonityMinimumBaseFeeUpdated represents a MinimumBaseFeeUpdated event raised by the Autonity contract.
type AutonityMinimumBaseFeeUpdated struct {
	GasPrice *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterMinimumBaseFeeUpdated is a free log retrieval operation binding the contract event 0x1f4d2fc7529047a5bd96d3229bfea127fd18b7748f13586e097c69fccd389128.
//
// Solidity: event MinimumBaseFeeUpdated(uint256 gasPrice)
func (_Autonity *AutonityFilterer) FilterMinimumBaseFeeUpdated(opts *bind.FilterOpts) (*AutonityMinimumBaseFeeUpdatedIterator, error) {

	logs, sub, err := _Autonity.contract.FilterLogs(opts, "MinimumBaseFeeUpdated")
	if err != nil {
		return nil, err
	}
	return &AutonityMinimumBaseFeeUpdatedIterator{contract: _Autonity.contract, event: "MinimumBaseFeeUpdated", logs: logs, sub: sub}, nil
}

// WatchMinimumBaseFeeUpdated is a free log subscription operation binding the contract event 0x1f4d2fc7529047a5bd96d3229bfea127fd18b7748f13586e097c69fccd389128.
//
// Solidity: event MinimumBaseFeeUpdated(uint256 gasPrice)
func (_Autonity *AutonityFilterer) WatchMinimumBaseFeeUpdated(opts *bind.WatchOpts, sink chan<- *AutonityMinimumBaseFeeUpdated) (event.Subscription, error) {

	logs, sub, err := _Autonity.contract.WatchLogs(opts, "MinimumBaseFeeUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AutonityMinimumBaseFeeUpdated)
				if err := _Autonity.contract.UnpackLog(event, "MinimumBaseFeeUpdated", log); err != nil {
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

// ParseMinimumBaseFeeUpdated is a log parse operation binding the contract event 0x1f4d2fc7529047a5bd96d3229bfea127fd18b7748f13586e097c69fccd389128.
//
// Solidity: event MinimumBaseFeeUpdated(uint256 gasPrice)
func (_Autonity *AutonityFilterer) ParseMinimumBaseFeeUpdated(log types.Log) (*AutonityMinimumBaseFeeUpdated, error) {
	event := new(AutonityMinimumBaseFeeUpdated)
	if err := _Autonity.contract.UnpackLog(event, "MinimumBaseFeeUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AutonityMintedStakeIterator is returned from FilterMintedStake and is used to iterate over the raw logs and unpacked data for MintedStake events raised by the Autonity contract.
type AutonityMintedStakeIterator struct {
	Event *AutonityMintedStake // Event containing the contract specifics and raw log

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
func (it *AutonityMintedStakeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutonityMintedStake)
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
		it.Event = new(AutonityMintedStake)
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
func (it *AutonityMintedStakeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AutonityMintedStakeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AutonityMintedStake represents a MintedStake event raised by the Autonity contract.
type AutonityMintedStake struct {
	Addr   common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterMintedStake is a free log retrieval operation binding the contract event 0x48490b4407bb949b708ec5f514b4167f08f4969baaf78d53b05028adf369bfcf.
//
// Solidity: event MintedStake(address addr, uint256 amount)
func (_Autonity *AutonityFilterer) FilterMintedStake(opts *bind.FilterOpts) (*AutonityMintedStakeIterator, error) {

	logs, sub, err := _Autonity.contract.FilterLogs(opts, "MintedStake")
	if err != nil {
		return nil, err
	}
	return &AutonityMintedStakeIterator{contract: _Autonity.contract, event: "MintedStake", logs: logs, sub: sub}, nil
}

// WatchMintedStake is a free log subscription operation binding the contract event 0x48490b4407bb949b708ec5f514b4167f08f4969baaf78d53b05028adf369bfcf.
//
// Solidity: event MintedStake(address addr, uint256 amount)
func (_Autonity *AutonityFilterer) WatchMintedStake(opts *bind.WatchOpts, sink chan<- *AutonityMintedStake) (event.Subscription, error) {

	logs, sub, err := _Autonity.contract.WatchLogs(opts, "MintedStake")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AutonityMintedStake)
				if err := _Autonity.contract.UnpackLog(event, "MintedStake", log); err != nil {
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
// Solidity: event MintedStake(address addr, uint256 amount)
func (_Autonity *AutonityFilterer) ParseMintedStake(log types.Log) (*AutonityMintedStake, error) {
	event := new(AutonityMintedStake)
	if err := _Autonity.contract.UnpackLog(event, "MintedStake", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AutonityMisbehaviourAddedIterator is returned from FilterMisbehaviourAdded and is used to iterate over the raw logs and unpacked data for MisbehaviourAdded events raised by the Autonity contract.
type AutonityMisbehaviourAddedIterator struct {
	Event *AutonityMisbehaviourAdded // Event containing the contract specifics and raw log

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
func (it *AutonityMisbehaviourAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutonityMisbehaviourAdded)
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
		it.Event = new(AutonityMisbehaviourAdded)
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
func (it *AutonityMisbehaviourAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AutonityMisbehaviourAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AutonityMisbehaviourAdded represents a MisbehaviourAdded event raised by the Autonity contract.
type AutonityMisbehaviourAdded struct {
	Ev  AutonityAccountabilityEvent
	Raw types.Log // Blockchain specific contextual infos
}

// FilterMisbehaviourAdded is a free log retrieval operation binding the contract event 0xe9b2e40b11e32b8729ed1bfd4c1ae17d2bcdc9af959564da14b39ca570607e3f.
//
// Solidity: event MisbehaviourAdded((uint8,uint8,uint8,uint8,address,address,bytes32,bytes) ev)
func (_Autonity *AutonityFilterer) FilterMisbehaviourAdded(opts *bind.FilterOpts) (*AutonityMisbehaviourAddedIterator, error) {

	logs, sub, err := _Autonity.contract.FilterLogs(opts, "MisbehaviourAdded")
	if err != nil {
		return nil, err
	}
	return &AutonityMisbehaviourAddedIterator{contract: _Autonity.contract, event: "MisbehaviourAdded", logs: logs, sub: sub}, nil
}

// WatchMisbehaviourAdded is a free log subscription operation binding the contract event 0xe9b2e40b11e32b8729ed1bfd4c1ae17d2bcdc9af959564da14b39ca570607e3f.
//
// Solidity: event MisbehaviourAdded((uint8,uint8,uint8,uint8,address,address,bytes32,bytes) ev)
func (_Autonity *AutonityFilterer) WatchMisbehaviourAdded(opts *bind.WatchOpts, sink chan<- *AutonityMisbehaviourAdded) (event.Subscription, error) {

	logs, sub, err := _Autonity.contract.WatchLogs(opts, "MisbehaviourAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AutonityMisbehaviourAdded)
				if err := _Autonity.contract.UnpackLog(event, "MisbehaviourAdded", log); err != nil {
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

// ParseMisbehaviourAdded is a log parse operation binding the contract event 0xe9b2e40b11e32b8729ed1bfd4c1ae17d2bcdc9af959564da14b39ca570607e3f.
//
// Solidity: event MisbehaviourAdded((uint8,uint8,uint8,uint8,address,address,bytes32,bytes) ev)
func (_Autonity *AutonityFilterer) ParseMisbehaviourAdded(log types.Log) (*AutonityMisbehaviourAdded, error) {
	event := new(AutonityMisbehaviourAdded)
	if err := _Autonity.contract.UnpackLog(event, "MisbehaviourAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AutonityMisbehaviourPenaltyUpdatedIterator is returned from FilterMisbehaviourPenaltyUpdated and is used to iterate over the raw logs and unpacked data for MisbehaviourPenaltyUpdated events raised by the Autonity contract.
type AutonityMisbehaviourPenaltyUpdatedIterator struct {
	Event *AutonityMisbehaviourPenaltyUpdated // Event containing the contract specifics and raw log

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
func (it *AutonityMisbehaviourPenaltyUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutonityMisbehaviourPenaltyUpdated)
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
		it.Event = new(AutonityMisbehaviourPenaltyUpdated)
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
func (it *AutonityMisbehaviourPenaltyUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AutonityMisbehaviourPenaltyUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AutonityMisbehaviourPenaltyUpdated represents a MisbehaviourPenaltyUpdated event raised by the Autonity contract.
type AutonityMisbehaviourPenaltyUpdated struct {
	Penalty *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterMisbehaviourPenaltyUpdated is a free log retrieval operation binding the contract event 0x3e4df1a42f35d79ea7cc3833604cd6377005ec5985514c94c229fb83f3650703.
//
// Solidity: event MisbehaviourPenaltyUpdated(uint256 penalty)
func (_Autonity *AutonityFilterer) FilterMisbehaviourPenaltyUpdated(opts *bind.FilterOpts) (*AutonityMisbehaviourPenaltyUpdatedIterator, error) {

	logs, sub, err := _Autonity.contract.FilterLogs(opts, "MisbehaviourPenaltyUpdated")
	if err != nil {
		return nil, err
	}
	return &AutonityMisbehaviourPenaltyUpdatedIterator{contract: _Autonity.contract, event: "MisbehaviourPenaltyUpdated", logs: logs, sub: sub}, nil
}

// WatchMisbehaviourPenaltyUpdated is a free log subscription operation binding the contract event 0x3e4df1a42f35d79ea7cc3833604cd6377005ec5985514c94c229fb83f3650703.
//
// Solidity: event MisbehaviourPenaltyUpdated(uint256 penalty)
func (_Autonity *AutonityFilterer) WatchMisbehaviourPenaltyUpdated(opts *bind.WatchOpts, sink chan<- *AutonityMisbehaviourPenaltyUpdated) (event.Subscription, error) {

	logs, sub, err := _Autonity.contract.WatchLogs(opts, "MisbehaviourPenaltyUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AutonityMisbehaviourPenaltyUpdated)
				if err := _Autonity.contract.UnpackLog(event, "MisbehaviourPenaltyUpdated", log); err != nil {
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

// ParseMisbehaviourPenaltyUpdated is a log parse operation binding the contract event 0x3e4df1a42f35d79ea7cc3833604cd6377005ec5985514c94c229fb83f3650703.
//
// Solidity: event MisbehaviourPenaltyUpdated(uint256 penalty)
func (_Autonity *AutonityFilterer) ParseMisbehaviourPenaltyUpdated(log types.Log) (*AutonityMisbehaviourPenaltyUpdated, error) {
	event := new(AutonityMisbehaviourPenaltyUpdated)
	if err := _Autonity.contract.UnpackLog(event, "MisbehaviourPenaltyUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AutonityNodeSlashedIterator is returned from FilterNodeSlashed and is used to iterate over the raw logs and unpacked data for NodeSlashed events raised by the Autonity contract.
type AutonityNodeSlashedIterator struct {
	Event *AutonityNodeSlashed // Event containing the contract specifics and raw log

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
func (it *AutonityNodeSlashedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutonityNodeSlashed)
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
		it.Event = new(AutonityNodeSlashed)
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
func (it *AutonityNodeSlashedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AutonityNodeSlashedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AutonityNodeSlashed represents a NodeSlashed event raised by the Autonity contract.
type AutonityNodeSlashed struct {
	Validator common.Address
	Penalty   *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterNodeSlashed is a free log retrieval operation binding the contract event 0x51cf713376ddb1e5f5828bb6aa39d99de812176d62c3d3550bdc4e0b5e86e1a5.
//
// Solidity: event NodeSlashed(address validator, uint256 penalty)
func (_Autonity *AutonityFilterer) FilterNodeSlashed(opts *bind.FilterOpts) (*AutonityNodeSlashedIterator, error) {

	logs, sub, err := _Autonity.contract.FilterLogs(opts, "NodeSlashed")
	if err != nil {
		return nil, err
	}
	return &AutonityNodeSlashedIterator{contract: _Autonity.contract, event: "NodeSlashed", logs: logs, sub: sub}, nil
}

// WatchNodeSlashed is a free log subscription operation binding the contract event 0x51cf713376ddb1e5f5828bb6aa39d99de812176d62c3d3550bdc4e0b5e86e1a5.
//
// Solidity: event NodeSlashed(address validator, uint256 penalty)
func (_Autonity *AutonityFilterer) WatchNodeSlashed(opts *bind.WatchOpts, sink chan<- *AutonityNodeSlashed) (event.Subscription, error) {

	logs, sub, err := _Autonity.contract.WatchLogs(opts, "NodeSlashed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AutonityNodeSlashed)
				if err := _Autonity.contract.UnpackLog(event, "NodeSlashed", log); err != nil {
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

// ParseNodeSlashed is a log parse operation binding the contract event 0x51cf713376ddb1e5f5828bb6aa39d99de812176d62c3d3550bdc4e0b5e86e1a5.
//
// Solidity: event NodeSlashed(address validator, uint256 penalty)
func (_Autonity *AutonityFilterer) ParseNodeSlashed(log types.Log) (*AutonityNodeSlashed, error) {
	event := new(AutonityNodeSlashed)
	if err := _Autonity.contract.UnpackLog(event, "NodeSlashed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AutonityPausedValidatorIterator is returned from FilterPausedValidator and is used to iterate over the raw logs and unpacked data for PausedValidator events raised by the Autonity contract.
type AutonityPausedValidatorIterator struct {
	Event *AutonityPausedValidator // Event containing the contract specifics and raw log

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
func (it *AutonityPausedValidatorIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutonityPausedValidator)
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
		it.Event = new(AutonityPausedValidator)
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
func (it *AutonityPausedValidatorIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AutonityPausedValidatorIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AutonityPausedValidator represents a PausedValidator event raised by the Autonity contract.
type AutonityPausedValidator struct {
	Treasury       common.Address
	Addr           common.Address
	EffectiveBlock *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterPausedValidator is a free log retrieval operation binding the contract event 0x75bdcdbe540758778e669d108fbcb7ede734f27f46e4e5525eeb8ecf91849a9c.
//
// Solidity: event PausedValidator(address treasury, address addr, uint256 effectiveBlock)
func (_Autonity *AutonityFilterer) FilterPausedValidator(opts *bind.FilterOpts) (*AutonityPausedValidatorIterator, error) {

	logs, sub, err := _Autonity.contract.FilterLogs(opts, "PausedValidator")
	if err != nil {
		return nil, err
	}
	return &AutonityPausedValidatorIterator{contract: _Autonity.contract, event: "PausedValidator", logs: logs, sub: sub}, nil
}

// WatchPausedValidator is a free log subscription operation binding the contract event 0x75bdcdbe540758778e669d108fbcb7ede734f27f46e4e5525eeb8ecf91849a9c.
//
// Solidity: event PausedValidator(address treasury, address addr, uint256 effectiveBlock)
func (_Autonity *AutonityFilterer) WatchPausedValidator(opts *bind.WatchOpts, sink chan<- *AutonityPausedValidator) (event.Subscription, error) {

	logs, sub, err := _Autonity.contract.WatchLogs(opts, "PausedValidator")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AutonityPausedValidator)
				if err := _Autonity.contract.UnpackLog(event, "PausedValidator", log); err != nil {
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
// Solidity: event PausedValidator(address treasury, address addr, uint256 effectiveBlock)
func (_Autonity *AutonityFilterer) ParsePausedValidator(log types.Log) (*AutonityPausedValidator, error) {
	event := new(AutonityPausedValidator)
	if err := _Autonity.contract.UnpackLog(event, "PausedValidator", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AutonityRegisteredValidatorIterator is returned from FilterRegisteredValidator and is used to iterate over the raw logs and unpacked data for RegisteredValidator events raised by the Autonity contract.
type AutonityRegisteredValidatorIterator struct {
	Event *AutonityRegisteredValidator // Event containing the contract specifics and raw log

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
func (it *AutonityRegisteredValidatorIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutonityRegisteredValidator)
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
		it.Event = new(AutonityRegisteredValidator)
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
func (it *AutonityRegisteredValidatorIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AutonityRegisteredValidatorIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AutonityRegisteredValidator represents a RegisteredValidator event raised by the Autonity contract.
type AutonityRegisteredValidator struct {
	Treasury       common.Address
	Addr           common.Address
	OracleAddress  common.Address
	Enode          string
	LiquidContract common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterRegisteredValidator is a free log retrieval operation binding the contract event 0x8ad8bd2eb6950e5f332fd3a6dca48cb358ecfe3057848902b98cbdfe455c915c.
//
// Solidity: event RegisteredValidator(address treasury, address addr, address oracleAddress, string enode, address liquidContract)
func (_Autonity *AutonityFilterer) FilterRegisteredValidator(opts *bind.FilterOpts) (*AutonityRegisteredValidatorIterator, error) {

	logs, sub, err := _Autonity.contract.FilterLogs(opts, "RegisteredValidator")
	if err != nil {
		return nil, err
	}
	return &AutonityRegisteredValidatorIterator{contract: _Autonity.contract, event: "RegisteredValidator", logs: logs, sub: sub}, nil
}

// WatchRegisteredValidator is a free log subscription operation binding the contract event 0x8ad8bd2eb6950e5f332fd3a6dca48cb358ecfe3057848902b98cbdfe455c915c.
//
// Solidity: event RegisteredValidator(address treasury, address addr, address oracleAddress, string enode, address liquidContract)
func (_Autonity *AutonityFilterer) WatchRegisteredValidator(opts *bind.WatchOpts, sink chan<- *AutonityRegisteredValidator) (event.Subscription, error) {

	logs, sub, err := _Autonity.contract.WatchLogs(opts, "RegisteredValidator")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AutonityRegisteredValidator)
				if err := _Autonity.contract.UnpackLog(event, "RegisteredValidator", log); err != nil {
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
// Solidity: event RegisteredValidator(address treasury, address addr, address oracleAddress, string enode, address liquidContract)
func (_Autonity *AutonityFilterer) ParseRegisteredValidator(log types.Log) (*AutonityRegisteredValidator, error) {
	event := new(AutonityRegisteredValidator)
	if err := _Autonity.contract.UnpackLog(event, "RegisteredValidator", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AutonityRewardedIterator is returned from FilterRewarded and is used to iterate over the raw logs and unpacked data for Rewarded events raised by the Autonity contract.
type AutonityRewardedIterator struct {
	Event *AutonityRewarded // Event containing the contract specifics and raw log

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
func (it *AutonityRewardedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutonityRewarded)
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
		it.Event = new(AutonityRewarded)
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
func (it *AutonityRewardedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AutonityRewardedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AutonityRewarded represents a Rewarded event raised by the Autonity contract.
type AutonityRewarded struct {
	Addr   common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterRewarded is a free log retrieval operation binding the contract event 0xb3b7a071186534c03b40695710096f289fd4ed6c1a374aff0bb648955e4fe563.
//
// Solidity: event Rewarded(address addr, uint256 amount)
func (_Autonity *AutonityFilterer) FilterRewarded(opts *bind.FilterOpts) (*AutonityRewardedIterator, error) {

	logs, sub, err := _Autonity.contract.FilterLogs(opts, "Rewarded")
	if err != nil {
		return nil, err
	}
	return &AutonityRewardedIterator{contract: _Autonity.contract, event: "Rewarded", logs: logs, sub: sub}, nil
}

// WatchRewarded is a free log subscription operation binding the contract event 0xb3b7a071186534c03b40695710096f289fd4ed6c1a374aff0bb648955e4fe563.
//
// Solidity: event Rewarded(address addr, uint256 amount)
func (_Autonity *AutonityFilterer) WatchRewarded(opts *bind.WatchOpts, sink chan<- *AutonityRewarded) (event.Subscription, error) {

	logs, sub, err := _Autonity.contract.WatchLogs(opts, "Rewarded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AutonityRewarded)
				if err := _Autonity.contract.UnpackLog(event, "Rewarded", log); err != nil {
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

// ParseRewarded is a log parse operation binding the contract event 0xb3b7a071186534c03b40695710096f289fd4ed6c1a374aff0bb648955e4fe563.
//
// Solidity: event Rewarded(address addr, uint256 amount)
func (_Autonity *AutonityFilterer) ParseRewarded(log types.Log) (*AutonityRewarded, error) {
	event := new(AutonityRewarded)
	if err := _Autonity.contract.UnpackLog(event, "Rewarded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AutonitySubmitGuiltyAccusationIterator is returned from FilterSubmitGuiltyAccusation and is used to iterate over the raw logs and unpacked data for SubmitGuiltyAccusation events raised by the Autonity contract.
type AutonitySubmitGuiltyAccusationIterator struct {
	Event *AutonitySubmitGuiltyAccusation // Event containing the contract specifics and raw log

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
func (it *AutonitySubmitGuiltyAccusationIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutonitySubmitGuiltyAccusation)
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
		it.Event = new(AutonitySubmitGuiltyAccusation)
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
func (it *AutonitySubmitGuiltyAccusationIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AutonitySubmitGuiltyAccusationIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AutonitySubmitGuiltyAccusation represents a SubmitGuiltyAccusation event raised by the Autonity contract.
type AutonitySubmitGuiltyAccusation struct {
	Ev  AutonityAccountabilityEvent
	Raw types.Log // Blockchain specific contextual infos
}

// FilterSubmitGuiltyAccusation is a free log retrieval operation binding the contract event 0x550841d9b29e92358159fe7d9bda9bfef1f1bc478fd9160241cff35168cfc712.
//
// Solidity: event SubmitGuiltyAccusation((uint8,uint8,uint8,uint8,address,address,bytes32,bytes) ev)
func (_Autonity *AutonityFilterer) FilterSubmitGuiltyAccusation(opts *bind.FilterOpts) (*AutonitySubmitGuiltyAccusationIterator, error) {

	logs, sub, err := _Autonity.contract.FilterLogs(opts, "SubmitGuiltyAccusation")
	if err != nil {
		return nil, err
	}
	return &AutonitySubmitGuiltyAccusationIterator{contract: _Autonity.contract, event: "SubmitGuiltyAccusation", logs: logs, sub: sub}, nil
}

// WatchSubmitGuiltyAccusation is a free log subscription operation binding the contract event 0x550841d9b29e92358159fe7d9bda9bfef1f1bc478fd9160241cff35168cfc712.
//
// Solidity: event SubmitGuiltyAccusation((uint8,uint8,uint8,uint8,address,address,bytes32,bytes) ev)
func (_Autonity *AutonityFilterer) WatchSubmitGuiltyAccusation(opts *bind.WatchOpts, sink chan<- *AutonitySubmitGuiltyAccusation) (event.Subscription, error) {

	logs, sub, err := _Autonity.contract.WatchLogs(opts, "SubmitGuiltyAccusation")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AutonitySubmitGuiltyAccusation)
				if err := _Autonity.contract.UnpackLog(event, "SubmitGuiltyAccusation", log); err != nil {
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

// ParseSubmitGuiltyAccusation is a log parse operation binding the contract event 0x550841d9b29e92358159fe7d9bda9bfef1f1bc478fd9160241cff35168cfc712.
//
// Solidity: event SubmitGuiltyAccusation((uint8,uint8,uint8,uint8,address,address,bytes32,bytes) ev)
func (_Autonity *AutonityFilterer) ParseSubmitGuiltyAccusation(log types.Log) (*AutonitySubmitGuiltyAccusation, error) {
	event := new(AutonitySubmitGuiltyAccusation)
	if err := _Autonity.contract.UnpackLog(event, "SubmitGuiltyAccusation", log); err != nil {
		return nil, err
	}
	event.Raw = log
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
	event.Raw = log
	return event, nil
}

// BytesLibMetaData contains all meta data concerning the BytesLib contract.
var BytesLibMetaData = &bind.MetaData{
	ABI: "[]",
	Bin: "0x60566037600b82828239805160001a607314602a57634e487b7160e01b600052600060045260246000fd5b30600052607381538281f3fe73000000000000000000000000000000000000000030146080604052600080fdfea2646970667358221220a630712f88b90175d230def65da0aef488d82a7db1344dcf40042170c1414a2364736f6c634300080c0033",
}

// BytesLibABI is the input ABI used to generate the binding from.
// Deprecated: Use BytesLibMetaData.ABI instead.
var BytesLibABI = BytesLibMetaData.ABI

// BytesLibBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use BytesLibMetaData.Bin instead.
var BytesLibBin = BytesLibMetaData.Bin

// DeployBytesLib deploys a new Ethereum contract, binding an instance of BytesLib to it.
func DeployBytesLib(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *BytesLib, error) {
	parsed, err := BytesLibMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(BytesLibBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &BytesLib{BytesLibCaller: BytesLibCaller{contract: contract}, BytesLibTransactor: BytesLibTransactor{contract: contract}, BytesLibFilterer: BytesLibFilterer{contract: contract}}, nil
}

// BytesLib is an auto generated Go binding around an Ethereum contract.
type BytesLib struct {
	BytesLibCaller     // Read-only binding to the contract
	BytesLibTransactor // Write-only binding to the contract
	BytesLibFilterer   // Log filterer for contract events
}

// BytesLibCaller is an auto generated read-only Go binding around an Ethereum contract.
type BytesLibCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BytesLibTransactor is an auto generated write-only Go binding around an Ethereum contract.
type BytesLibTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BytesLibFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type BytesLibFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BytesLibSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type BytesLibSession struct {
	Contract     *BytesLib         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// BytesLibCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type BytesLibCallerSession struct {
	Contract *BytesLibCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// BytesLibTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type BytesLibTransactorSession struct {
	Contract     *BytesLibTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// BytesLibRaw is an auto generated low-level Go binding around an Ethereum contract.
type BytesLibRaw struct {
	Contract *BytesLib // Generic contract binding to access the raw methods on
}

// BytesLibCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type BytesLibCallerRaw struct {
	Contract *BytesLibCaller // Generic read-only contract binding to access the raw methods on
}

// BytesLibTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type BytesLibTransactorRaw struct {
	Contract *BytesLibTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBytesLib creates a new instance of BytesLib, bound to a specific deployed contract.
func NewBytesLib(address common.Address, backend bind.ContractBackend) (*BytesLib, error) {
	contract, err := bindBytesLib(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &BytesLib{BytesLibCaller: BytesLibCaller{contract: contract}, BytesLibTransactor: BytesLibTransactor{contract: contract}, BytesLibFilterer: BytesLibFilterer{contract: contract}}, nil
}

// NewBytesLibCaller creates a new read-only instance of BytesLib, bound to a specific deployed contract.
func NewBytesLibCaller(address common.Address, caller bind.ContractCaller) (*BytesLibCaller, error) {
	contract, err := bindBytesLib(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BytesLibCaller{contract: contract}, nil
}

// NewBytesLibTransactor creates a new write-only instance of BytesLib, bound to a specific deployed contract.
func NewBytesLibTransactor(address common.Address, transactor bind.ContractTransactor) (*BytesLibTransactor, error) {
	contract, err := bindBytesLib(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BytesLibTransactor{contract: contract}, nil
}

// NewBytesLibFilterer creates a new log filterer instance of BytesLib, bound to a specific deployed contract.
func NewBytesLibFilterer(address common.Address, filterer bind.ContractFilterer) (*BytesLibFilterer, error) {
	contract, err := bindBytesLib(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BytesLibFilterer{contract: contract}, nil
}

// bindBytesLib binds a generic wrapper to an already deployed contract.
func bindBytesLib(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(BytesLibABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BytesLib *BytesLibRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BytesLib.Contract.BytesLibCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BytesLib *BytesLibRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BytesLib.Contract.BytesLibTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BytesLib *BytesLibRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BytesLib.Contract.BytesLibTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BytesLib *BytesLibCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BytesLib.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BytesLib *BytesLibTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BytesLib.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BytesLib *BytesLibTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BytesLib.Contract.contract.Transact(opts, method, params...)
}

// HelpersMetaData contains all meta data concerning the Helpers contract.
var HelpersMetaData = &bind.MetaData{
	ABI: "[]",
	Bin: "0x60566037600b82828239805160001a607314602a57634e487b7160e01b600052600060045260246000fd5b30600052607381538281f3fe73000000000000000000000000000000000000000030146080604052600080fdfea26469706673582212201de81a26632f66746e9b76448f58c71ccc491ebbe7568c53108f7b224d29703564736f6c634300080c0033",
}

// HelpersABI is the input ABI used to generate the binding from.
// Deprecated: Use HelpersMetaData.ABI instead.
var HelpersABI = HelpersMetaData.ABI

// HelpersBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use HelpersMetaData.Bin instead.
var HelpersBin = HelpersMetaData.Bin

// DeployHelpers deploys a new Ethereum contract, binding an instance of Helpers to it.
func DeployHelpers(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Helpers, error) {
	parsed, err := HelpersMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(HelpersBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Helpers{HelpersCaller: HelpersCaller{contract: contract}, HelpersTransactor: HelpersTransactor{contract: contract}, HelpersFilterer: HelpersFilterer{contract: contract}}, nil
}

// Helpers is an auto generated Go binding around an Ethereum contract.
type Helpers struct {
	HelpersCaller     // Read-only binding to the contract
	HelpersTransactor // Write-only binding to the contract
	HelpersFilterer   // Log filterer for contract events
}

// HelpersCaller is an auto generated read-only Go binding around an Ethereum contract.
type HelpersCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// HelpersTransactor is an auto generated write-only Go binding around an Ethereum contract.
type HelpersTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// HelpersFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type HelpersFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// HelpersSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type HelpersSession struct {
	Contract     *Helpers          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// HelpersCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type HelpersCallerSession struct {
	Contract *HelpersCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// HelpersTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type HelpersTransactorSession struct {
	Contract     *HelpersTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// HelpersRaw is an auto generated low-level Go binding around an Ethereum contract.
type HelpersRaw struct {
	Contract *Helpers // Generic contract binding to access the raw methods on
}

// HelpersCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type HelpersCallerRaw struct {
	Contract *HelpersCaller // Generic read-only contract binding to access the raw methods on
}

// HelpersTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type HelpersTransactorRaw struct {
	Contract *HelpersTransactor // Generic write-only contract binding to access the raw methods on
}

// NewHelpers creates a new instance of Helpers, bound to a specific deployed contract.
func NewHelpers(address common.Address, backend bind.ContractBackend) (*Helpers, error) {
	contract, err := bindHelpers(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Helpers{HelpersCaller: HelpersCaller{contract: contract}, HelpersTransactor: HelpersTransactor{contract: contract}, HelpersFilterer: HelpersFilterer{contract: contract}}, nil
}

// NewHelpersCaller creates a new read-only instance of Helpers, bound to a specific deployed contract.
func NewHelpersCaller(address common.Address, caller bind.ContractCaller) (*HelpersCaller, error) {
	contract, err := bindHelpers(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &HelpersCaller{contract: contract}, nil
}

// NewHelpersTransactor creates a new write-only instance of Helpers, bound to a specific deployed contract.
func NewHelpersTransactor(address common.Address, transactor bind.ContractTransactor) (*HelpersTransactor, error) {
	contract, err := bindHelpers(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &HelpersTransactor{contract: contract}, nil
}

// NewHelpersFilterer creates a new log filterer instance of Helpers, bound to a specific deployed contract.
func NewHelpersFilterer(address common.Address, filterer bind.ContractFilterer) (*HelpersFilterer, error) {
	contract, err := bindHelpers(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &HelpersFilterer{contract: contract}, nil
}

// bindHelpers binds a generic wrapper to an already deployed contract.
func bindHelpers(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(HelpersABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Helpers *HelpersRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Helpers.Contract.HelpersCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Helpers *HelpersRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Helpers.Contract.HelpersTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Helpers *HelpersRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Helpers.Contract.HelpersTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Helpers *HelpersCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Helpers.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Helpers *HelpersTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Helpers.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Helpers *HelpersTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Helpers.Contract.contract.Transact(opts, method, params...)
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

// IOracleMetaData contains all meta data concerning the IOracle contract.
var IOracleMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_round\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_height\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_timestamp\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_votePeriod\",\"type\":\"uint256\"}],\"name\":\"NewRound\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string[]\",\"name\":\"_symbols\",\"type\":\"string[]\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_round\",\"type\":\"uint256\"}],\"name\":\"NewSymbols\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_voter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"int256[]\",\"name\":\"_votes\",\"type\":\"int256[]\"}],\"name\":\"Voted\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"getPrecision\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRound\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_round\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"_symbol\",\"type\":\"string\"}],\"name\":\"getRoundData\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"round\",\"type\":\"uint256\"},{\"internalType\":\"int256\",\"name\":\"price\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"status\",\"type\":\"uint256\"}],\"internalType\":\"structIOracle.RoundData\",\"name\":\"data\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getSymbols\",\"outputs\":[{\"internalType\":\"string[]\",\"name\":\"_symbols\",\"type\":\"string[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getVotePeriod\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getVoters\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_symbol\",\"type\":\"string\"}],\"name\":\"latestRoundData\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"round\",\"type\":\"uint256\"},{\"internalType\":\"int256\",\"name\":\"price\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"status\",\"type\":\"uint256\"}],\"internalType\":\"structIOracle.RoundData\",\"name\":\"data\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string[]\",\"name\":\"_symbols\",\"type\":\"string[]\"}],\"name\":\"setSymbols\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_commit\",\"type\":\"uint256\"},{\"internalType\":\"int256[]\",\"name\":\"_reports\",\"type\":\"int256[]\"},{\"internalType\":\"uint256\",\"name\":\"_salt\",\"type\":\"uint256\"}],\"name\":\"vote\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"9670c0bc": "getPrecision()",
		"9f8743f7": "getRound()",
		"3c8510fd": "getRoundData(uint256,string)",
		"df7f710e": "getSymbols()",
		"b78dec52": "getVotePeriod()",
		"cdd72253": "getVoters()",
		"33f98c77": "latestRoundData(string)",
		"8d4f75d2": "setSymbols(string[])",
		"307de9b6": "vote(uint256,int256[],uint256)",
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

// GetPrecision is a free data retrieval call binding the contract method 0x9670c0bc.
//
// Solidity: function getPrecision() view returns(uint256)
func (_IOracle *IOracleCaller) GetPrecision(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IOracle.contract.Call(opts, &out, "getPrecision")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetPrecision is a free data retrieval call binding the contract method 0x9670c0bc.
//
// Solidity: function getPrecision() view returns(uint256)
func (_IOracle *IOracleSession) GetPrecision() (*big.Int, error) {
	return _IOracle.Contract.GetPrecision(&_IOracle.CallOpts)
}

// GetPrecision is a free data retrieval call binding the contract method 0x9670c0bc.
//
// Solidity: function getPrecision() view returns(uint256)
func (_IOracle *IOracleCallerSession) GetPrecision() (*big.Int, error) {
	return _IOracle.Contract.GetPrecision(&_IOracle.CallOpts)
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
// Solidity: function getRoundData(uint256 _round, string _symbol) view returns((uint256,int256,uint256,uint256) data)
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
// Solidity: function getRoundData(uint256 _round, string _symbol) view returns((uint256,int256,uint256,uint256) data)
func (_IOracle *IOracleSession) GetRoundData(_round *big.Int, _symbol string) (IOracleRoundData, error) {
	return _IOracle.Contract.GetRoundData(&_IOracle.CallOpts, _round, _symbol)
}

// GetRoundData is a free data retrieval call binding the contract method 0x3c8510fd.
//
// Solidity: function getRoundData(uint256 _round, string _symbol) view returns((uint256,int256,uint256,uint256) data)
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
// Solidity: function latestRoundData(string _symbol) view returns((uint256,int256,uint256,uint256) data)
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
// Solidity: function latestRoundData(string _symbol) view returns((uint256,int256,uint256,uint256) data)
func (_IOracle *IOracleSession) LatestRoundData(_symbol string) (IOracleRoundData, error) {
	return _IOracle.Contract.LatestRoundData(&_IOracle.CallOpts, _symbol)
}

// LatestRoundData is a free data retrieval call binding the contract method 0x33f98c77.
//
// Solidity: function latestRoundData(string _symbol) view returns((uint256,int256,uint256,uint256) data)
func (_IOracle *IOracleCallerSession) LatestRoundData(_symbol string) (IOracleRoundData, error) {
	return _IOracle.Contract.LatestRoundData(&_IOracle.CallOpts, _symbol)
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

// Vote is a paid mutator transaction binding the contract method 0x307de9b6.
//
// Solidity: function vote(uint256 _commit, int256[] _reports, uint256 _salt) returns()
func (_IOracle *IOracleTransactor) Vote(opts *bind.TransactOpts, _commit *big.Int, _reports []*big.Int, _salt *big.Int) (*types.Transaction, error) {
	return _IOracle.contract.Transact(opts, "vote", _commit, _reports, _salt)
}

// Vote is a paid mutator transaction binding the contract method 0x307de9b6.
//
// Solidity: function vote(uint256 _commit, int256[] _reports, uint256 _salt) returns()
func (_IOracle *IOracleSession) Vote(_commit *big.Int, _reports []*big.Int, _salt *big.Int) (*types.Transaction, error) {
	return _IOracle.Contract.Vote(&_IOracle.TransactOpts, _commit, _reports, _salt)
}

// Vote is a paid mutator transaction binding the contract method 0x307de9b6.
//
// Solidity: function vote(uint256 _commit, int256[] _reports, uint256 _salt) returns()
func (_IOracle *IOracleTransactorSession) Vote(_commit *big.Int, _reports []*big.Int, _salt *big.Int) (*types.Transaction, error) {
	return _IOracle.Contract.Vote(&_IOracle.TransactOpts, _commit, _reports, _salt)
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
	Height     *big.Int
	Timestamp  *big.Int
	VotePeriod *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterNewRound is a free log retrieval operation binding the contract event 0xb5d8636ab45e6cac7a4a61cb7c77f77f61a454d73aa2e6139ff8dcaf463537e5.
//
// Solidity: event NewRound(uint256 _round, uint256 _height, uint256 _timestamp, uint256 _votePeriod)
func (_IOracle *IOracleFilterer) FilterNewRound(opts *bind.FilterOpts) (*IOracleNewRoundIterator, error) {

	logs, sub, err := _IOracle.contract.FilterLogs(opts, "NewRound")
	if err != nil {
		return nil, err
	}
	return &IOracleNewRoundIterator{contract: _IOracle.contract, event: "NewRound", logs: logs, sub: sub}, nil
}

// WatchNewRound is a free log subscription operation binding the contract event 0xb5d8636ab45e6cac7a4a61cb7c77f77f61a454d73aa2e6139ff8dcaf463537e5.
//
// Solidity: event NewRound(uint256 _round, uint256 _height, uint256 _timestamp, uint256 _votePeriod)
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

// ParseNewRound is a log parse operation binding the contract event 0xb5d8636ab45e6cac7a4a61cb7c77f77f61a454d73aa2e6139ff8dcaf463537e5.
//
// Solidity: event NewRound(uint256 _round, uint256 _height, uint256 _timestamp, uint256 _votePeriod)
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

// IOracleVotedIterator is returned from FilterVoted and is used to iterate over the raw logs and unpacked data for Voted events raised by the IOracle contract.
type IOracleVotedIterator struct {
	Event *IOracleVoted // Event containing the contract specifics and raw log

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
func (it *IOracleVotedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IOracleVoted)
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
		it.Event = new(IOracleVoted)
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
func (it *IOracleVotedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IOracleVotedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IOracleVoted represents a Voted event raised by the IOracle contract.
type IOracleVoted struct {
	Voter common.Address
	Votes []*big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterVoted is a free log retrieval operation binding the contract event 0xd0d8560f1076ac6b216b1091a2571d6f9bc3e0889f4dbdbe1c7d1be7136714d3.
//
// Solidity: event Voted(address indexed _voter, int256[] _votes)
func (_IOracle *IOracleFilterer) FilterVoted(opts *bind.FilterOpts, _voter []common.Address) (*IOracleVotedIterator, error) {

	var _voterRule []interface{}
	for _, _voterItem := range _voter {
		_voterRule = append(_voterRule, _voterItem)
	}

	logs, sub, err := _IOracle.contract.FilterLogs(opts, "Voted", _voterRule)
	if err != nil {
		return nil, err
	}
	return &IOracleVotedIterator{contract: _IOracle.contract, event: "Voted", logs: logs, sub: sub}, nil
}

// WatchVoted is a free log subscription operation binding the contract event 0xd0d8560f1076ac6b216b1091a2571d6f9bc3e0889f4dbdbe1c7d1be7136714d3.
//
// Solidity: event Voted(address indexed _voter, int256[] _votes)
func (_IOracle *IOracleFilterer) WatchVoted(opts *bind.WatchOpts, sink chan<- *IOracleVoted, _voter []common.Address) (event.Subscription, error) {

	var _voterRule []interface{}
	for _, _voterItem := range _voter {
		_voterRule = append(_voterRule, _voterItem)
	}

	logs, sub, err := _IOracle.contract.WatchLogs(opts, "Voted", _voterRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IOracleVoted)
				if err := _IOracle.contract.UnpackLog(event, "Voted", log); err != nil {
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

// ParseVoted is a log parse operation binding the contract event 0xd0d8560f1076ac6b216b1091a2571d6f9bc3e0889f4dbdbe1c7d1be7136714d3.
//
// Solidity: event Voted(address indexed _voter, int256[] _votes)
func (_IOracle *IOracleFilterer) ParseVoted(log types.Log) (*IOracleVoted, error) {
	event := new(IOracleVoted)
	if err := _IOracle.contract.UnpackLog(event, "Voted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LiquidMetaData contains all meta data concerning the Liquid contract.
var LiquidMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_validator\",\"type\":\"address\"},{\"internalType\":\"addresspayable\",\"name\":\"_treasury\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_commissionRate\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"_index\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"COMMISSION_RATE_PRECISION\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"FEE_FACTOR_UNIT_RECIP\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_delegator\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"claimRewards\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"mint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"redistribute\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_rate\",\"type\":\"uint256\"}],\"name\":\"setCommissionRate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"_success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"_success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_account\",\"type\":\"address\"}],\"name\":\"unclaimedRewards\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"2f2c3f2e": "COMMISSION_RATE_PRECISION()",
		"187cf4d7": "FEE_FACTOR_UNIT_RECIP()",
		"dd62ed3e": "allowance(address,address)",
		"095ea7b3": "approve(address,uint256)",
		"70a08231": "balanceOf(address)",
		"9dc29fac": "burn(address,uint256)",
		"372500ab": "claimRewards()",
		"313ce567": "decimals()",
		"40c10f19": "mint(address,uint256)",
		"06fdde03": "name()",
		"fb489a7b": "redistribute()",
		"19fac8fd": "setCommissionRate(uint256)",
		"95d89b41": "symbol()",
		"18160ddd": "totalSupply()",
		"a9059cbb": "transfer(address,uint256)",
		"23b872dd": "transferFrom(address,address,uint256)",
		"949813b8": "unclaimedRewards(address)",
	},
	Bin: "0x60806040523480156200001157600080fd5b506040516200116a3803806200116a833981016040819052620000349162000212565b6127108211156200004457600080fd5b600980546001600160a01b038087166001600160a01b031992831617909255600a805492861692909116919091179055600b8290556040516200008c908290602001620002ff565b60405160208183030381529060405260079080519060200190620000b29291906200010a565b5080604051602001620000c69190620002ff565b60405160208183030381529060405260089080519060200190620000ec9291906200010a565b5050600080546001600160a01b03191633179055506200036b915050565b82805462000118906200032e565b90600052602060002090601f0160209004810192826200013c576000855562000187565b82601f106200015757805160ff191683800117855562000187565b8280016001018555821562000187579182015b82811115620001875782518255916020019190600101906200016a565b506200019592915062000199565b5090565b5b808211156200019557600081556001016200019a565b6001600160a01b0381168114620001c657600080fd5b50565b634e487b7160e01b600052604160045260246000fd5b60005b83811015620001fc578181015183820152602001620001e2565b838111156200020c576000848401525b50505050565b600080600080608085870312156200022957600080fd5b84516200023681620001b0565b60208601519094506200024981620001b0565b6040860151606087015191945092506001600160401b03808211156200026e57600080fd5b818701915087601f8301126200028357600080fd5b815181811115620002985762000298620001c9565b604051601f8201601f19908116603f01168101908382118183101715620002c357620002c3620001c9565b816040528281528a6020848701011115620002dd57600080fd5b620002f0836020830160208801620001df565b979a9699509497505050505050565b644c4e544e2d60d81b81526000825162000321816005850160208701620001df565b9190910160050192915050565b600181811c908216806200034357607f821691505b602082108114156200036557634e487b7160e01b600052602260045260246000fd5b50919050565b610def806200037b6000396000f3fe6080604052600436106100fe5760003560e01c8063372500ab1161009557806395d89b411161006457806395d89b41146102945780639dc29fac146102a9578063a9059cbb146102c9578063dd62ed3e146102e9578063fb489a7b1461032f57600080fd5b8063372500ab1461020957806340c10f191461021e57806370a082311461023e578063949813b81461027457600080fd5b806319fac8fd116100d157806319fac8fd1461019557806323b872dd146101b75780632f2c3f2e146101d7578063313ce567146101ed57600080fd5b806306fdde0314610103578063095ea7b31461012e57806318160ddd1461015e578063187cf4d71461017d575b600080fd5b34801561010f57600080fd5b50610118610337565b6040516101259190610b4b565b60405180910390f35b34801561013a57600080fd5b5061014e610149366004610bbc565b6103c9565b6040519015158152602001610125565b34801561016a57600080fd5b506003545b604051908152602001610125565b34801561018957600080fd5b5061016f633b9aca0081565b3480156101a157600080fd5b506101b56101b0366004610be6565b6103df565b005b3480156101c357600080fd5b5061014e6101d2366004610bff565b610417565b3480156101e357600080fd5b5061016f61271081565b3480156101f957600080fd5b5060405160128152602001610125565b34801561021557600080fd5b506101b561050a565b34801561022a57600080fd5b506101b5610239366004610bbc565b6105b8565b34801561024a57600080fd5b5061016f610259366004610c3b565b6001600160a01b031660009081526001602052604090205490565b34801561028057600080fd5b5061016f61028f366004610c3b565b610620565b3480156102a057600080fd5b50610118610654565b3480156102b557600080fd5b506101b56102c4366004610bbc565b610663565b3480156102d557600080fd5b5061014e6102e4366004610bbc565b6106c3565b3480156102f557600080fd5b5061016f610304366004610c5d565b6001600160a01b03918216600090815260026020908152604080832093909416825291909152205490565b61016f610710565b60606007805461034690610c90565b80601f016020809104026020016040519081016040528092919081815260200182805461037290610c90565b80156103bf5780601f10610394576101008083540402835291602001916103bf565b820191906000526020600020905b8154815290600101906020018083116103a257829003601f168201915b5050505050905090565b60006103d6338484610858565b50600192915050565b6000546001600160a01b031633146104125760405162461bcd60e51b815260040161040990610ccb565b60405180910390fd5b600b55565b6001600160a01b03831660009081526002602090815260408083203384529091528120548281101561049c5760405162461bcd60e51b815260206004820152602860248201527f45524332303a207472616e7366657220616d6f756e74206578636565647320616044820152676c6c6f77616e636560c01b6064820152608401610409565b6104b085336104ab8685610d29565b610858565b6104ba858461097c565b6104c48484610a1f565b836001600160a01b0316856001600160a01b0316600080516020610d9a833981519152856040516104f791815260200190565b60405180910390a3506001949350505050565b600061051533610a73565b33600081815260046020526040808220829055519293509183908381818185875af1925050503d8060008114610567576040519150601f19603f3d011682016040523d82523d6000602084013e61056c565b606091505b50509050806105b45760405162461bcd60e51b81526020600482015260146024820152732330b4b632b2103a379039b2b7321022ba3432b960611b6044820152606401610409565b5050565b6000546001600160a01b031633146105e25760405162461bcd60e51b815260040161040990610ccb565b6105ec8282610a1f565b6040518181526001600160a01b03831690600090600080516020610d9a833981519152906020015b60405180910390a35050565b600061062b82610ad8565b6001600160a01b03831660009081526004602052604090205461064e9190610d40565b92915050565b60606008805461034690610c90565b6000546001600160a01b0316331461068d5760405162461bcd60e51b815260040161040990610ccb565b610697828261097c565b6040518181526000906001600160a01b03841690600080516020610d9a83398151915290602001610614565b60006106cf338361097c565b6106d98383610a1f565b6040518281526001600160a01b038416903390600080516020610d9a8339815191529060200160405180910390a350600192915050565b600080546001600160a01b0316331461073b5760405162461bcd60e51b815260040161040990610ccb565b600b543490600090612710906107519084610d58565b61075b9190610d77565b90508181106107ac5760405162461bcd60e51b815260206004820152601860248201527f696e76616c69642076616c696461746f722072657761726400000000000000006044820152606401610409565b6107b68183610d29565b600a546040519193506001600160a01b03169082156108fc029083906000818181858888f193505050501580156107f1573d6000803e3d6000fd5b50600354600090610806633b9aca0085610d58565b6108109190610d77565b9050806006546108209190610d40565b600655600354600090633b9aca00906108399084610d58565b6108439190610d77565b905061084f8184610d40565b94505050505090565b6001600160a01b0383166108ba5760405162461bcd60e51b8152602060048201526024808201527f45524332303a20617070726f76652066726f6d20746865207a65726f206164646044820152637265737360e01b6064820152608401610409565b6001600160a01b03821661091b5760405162461bcd60e51b815260206004820152602260248201527f45524332303a20617070726f766520746f20746865207a65726f206164647265604482015261737360f01b6064820152608401610409565b6001600160a01b0383811660008181526002602090815260408083209487168084529482529182902085905590518481527f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925910160405180910390a3505050565b61098582610a73565b506001600160a01b038216600090815260016020526040902054808211156109ac57600080fd5b808210156109dc576109be8282610d29565b6001600160a01b038416600090815260016020526040902055610a03565b6001600160a01b038316600090815260016020908152604080832083905560059091528120555b8160036000828254610a159190610d29565b9091555050505050565b610a2882610a73565b506001600160a01b03821660009081526001602052604081208054839290610a51908490610d40565b925050819055508060036000828254610a6a9190610d40565b90915550505050565b600080610a7f83610ad8565b6001600160a01b038416600090815260046020526040902054909150610aa6908290610d40565b6001600160a01b0390931660009081526004602090815260408083208690556006546005909252909120555090919050565b6001600160a01b03811660009081526001602052604081205480610aff5750600092915050565b6001600160a01b038316600090815260056020526040812054600654610b259190610d29565b90506000633b9aca00610b388484610d58565b610b429190610d77565b95945050505050565b600060208083528351808285015260005b81811015610b7857858101830151858201604001528201610b5c565b81811115610b8a576000604083870101525b50601f01601f1916929092016040019392505050565b80356001600160a01b0381168114610bb757600080fd5b919050565b60008060408385031215610bcf57600080fd5b610bd883610ba0565b946020939093013593505050565b600060208284031215610bf857600080fd5b5035919050565b600080600060608486031215610c1457600080fd5b610c1d84610ba0565b9250610c2b60208501610ba0565b9150604084013590509250925092565b600060208284031215610c4d57600080fd5b610c5682610ba0565b9392505050565b60008060408385031215610c7057600080fd5b610c7983610ba0565b9150610c8760208401610ba0565b90509250929050565b600181811c90821680610ca457607f821691505b60208210811415610cc557634e487b7160e01b600052602260045260246000fd5b50919050565b60208082526028908201527f43616c6c207265737472696374656420746f20746865204175746f6e6974792060408201526710dbdb9d1c9858dd60c21b606082015260800190565b634e487b7160e01b600052601160045260246000fd5b600082821015610d3b57610d3b610d13565b500390565b60008219821115610d5357610d53610d13565b500190565b6000816000190483118215151615610d7257610d72610d13565b500290565b600082610d9457634e487b7160e01b600052601260045260246000fd5b50049056feddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3efa26469706673582212205311cbb54f78618267a028af290c7924ffb605213b8e8d196c0da082fc228cb264736f6c634300080c0033",
}

// LiquidABI is the input ABI used to generate the binding from.
// Deprecated: Use LiquidMetaData.ABI instead.
var LiquidABI = LiquidMetaData.ABI

// Deprecated: Use LiquidMetaData.Sigs instead.
// LiquidFuncSigs maps the 4-byte function signature to its string representation.
var LiquidFuncSigs = LiquidMetaData.Sigs

// LiquidBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use LiquidMetaData.Bin instead.
var LiquidBin = LiquidMetaData.Bin

// DeployLiquid deploys a new Ethereum contract, binding an instance of Liquid to it.
func DeployLiquid(auth *bind.TransactOpts, backend bind.ContractBackend, _validator common.Address, _treasury common.Address, _commissionRate *big.Int, _index string) (common.Address, *types.Transaction, *Liquid, error) {
	parsed, err := LiquidMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(LiquidBin), backend, _validator, _treasury, _commissionRate, _index)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Liquid{LiquidCaller: LiquidCaller{contract: contract}, LiquidTransactor: LiquidTransactor{contract: contract}, LiquidFilterer: LiquidFilterer{contract: contract}}, nil
}

// Liquid is an auto generated Go binding around an Ethereum contract.
type Liquid struct {
	LiquidCaller     // Read-only binding to the contract
	LiquidTransactor // Write-only binding to the contract
	LiquidFilterer   // Log filterer for contract events
}

// LiquidCaller is an auto generated read-only Go binding around an Ethereum contract.
type LiquidCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LiquidTransactor is an auto generated write-only Go binding around an Ethereum contract.
type LiquidTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LiquidFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type LiquidFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LiquidSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type LiquidSession struct {
	Contract     *Liquid           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// LiquidCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type LiquidCallerSession struct {
	Contract *LiquidCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// LiquidTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type LiquidTransactorSession struct {
	Contract     *LiquidTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// LiquidRaw is an auto generated low-level Go binding around an Ethereum contract.
type LiquidRaw struct {
	Contract *Liquid // Generic contract binding to access the raw methods on
}

// LiquidCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type LiquidCallerRaw struct {
	Contract *LiquidCaller // Generic read-only contract binding to access the raw methods on
}

// LiquidTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type LiquidTransactorRaw struct {
	Contract *LiquidTransactor // Generic write-only contract binding to access the raw methods on
}

// NewLiquid creates a new instance of Liquid, bound to a specific deployed contract.
func NewLiquid(address common.Address, backend bind.ContractBackend) (*Liquid, error) {
	contract, err := bindLiquid(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Liquid{LiquidCaller: LiquidCaller{contract: contract}, LiquidTransactor: LiquidTransactor{contract: contract}, LiquidFilterer: LiquidFilterer{contract: contract}}, nil
}

// NewLiquidCaller creates a new read-only instance of Liquid, bound to a specific deployed contract.
func NewLiquidCaller(address common.Address, caller bind.ContractCaller) (*LiquidCaller, error) {
	contract, err := bindLiquid(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &LiquidCaller{contract: contract}, nil
}

// NewLiquidTransactor creates a new write-only instance of Liquid, bound to a specific deployed contract.
func NewLiquidTransactor(address common.Address, transactor bind.ContractTransactor) (*LiquidTransactor, error) {
	contract, err := bindLiquid(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &LiquidTransactor{contract: contract}, nil
}

// NewLiquidFilterer creates a new log filterer instance of Liquid, bound to a specific deployed contract.
func NewLiquidFilterer(address common.Address, filterer bind.ContractFilterer) (*LiquidFilterer, error) {
	contract, err := bindLiquid(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &LiquidFilterer{contract: contract}, nil
}

// bindLiquid binds a generic wrapper to an already deployed contract.
func bindLiquid(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(LiquidABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Liquid *LiquidRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Liquid.Contract.LiquidCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Liquid *LiquidRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Liquid.Contract.LiquidTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Liquid *LiquidRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Liquid.Contract.LiquidTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Liquid *LiquidCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Liquid.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Liquid *LiquidTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Liquid.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Liquid *LiquidTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Liquid.Contract.contract.Transact(opts, method, params...)
}

// COMMISSIONRATEPRECISION is a free data retrieval call binding the contract method 0x2f2c3f2e.
//
// Solidity: function COMMISSION_RATE_PRECISION() view returns(uint256)
func (_Liquid *LiquidCaller) COMMISSIONRATEPRECISION(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Liquid.contract.Call(opts, &out, "COMMISSION_RATE_PRECISION")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// COMMISSIONRATEPRECISION is a free data retrieval call binding the contract method 0x2f2c3f2e.
//
// Solidity: function COMMISSION_RATE_PRECISION() view returns(uint256)
func (_Liquid *LiquidSession) COMMISSIONRATEPRECISION() (*big.Int, error) {
	return _Liquid.Contract.COMMISSIONRATEPRECISION(&_Liquid.CallOpts)
}

// COMMISSIONRATEPRECISION is a free data retrieval call binding the contract method 0x2f2c3f2e.
//
// Solidity: function COMMISSION_RATE_PRECISION() view returns(uint256)
func (_Liquid *LiquidCallerSession) COMMISSIONRATEPRECISION() (*big.Int, error) {
	return _Liquid.Contract.COMMISSIONRATEPRECISION(&_Liquid.CallOpts)
}

// FEEFACTORUNITRECIP is a free data retrieval call binding the contract method 0x187cf4d7.
//
// Solidity: function FEE_FACTOR_UNIT_RECIP() view returns(uint256)
func (_Liquid *LiquidCaller) FEEFACTORUNITRECIP(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Liquid.contract.Call(opts, &out, "FEE_FACTOR_UNIT_RECIP")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// FEEFACTORUNITRECIP is a free data retrieval call binding the contract method 0x187cf4d7.
//
// Solidity: function FEE_FACTOR_UNIT_RECIP() view returns(uint256)
func (_Liquid *LiquidSession) FEEFACTORUNITRECIP() (*big.Int, error) {
	return _Liquid.Contract.FEEFACTORUNITRECIP(&_Liquid.CallOpts)
}

// FEEFACTORUNITRECIP is a free data retrieval call binding the contract method 0x187cf4d7.
//
// Solidity: function FEE_FACTOR_UNIT_RECIP() view returns(uint256)
func (_Liquid *LiquidCallerSession) FEEFACTORUNITRECIP() (*big.Int, error) {
	return _Liquid.Contract.FEEFACTORUNITRECIP(&_Liquid.CallOpts)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address _owner, address _spender) view returns(uint256)
func (_Liquid *LiquidCaller) Allowance(opts *bind.CallOpts, _owner common.Address, _spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Liquid.contract.Call(opts, &out, "allowance", _owner, _spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address _owner, address _spender) view returns(uint256)
func (_Liquid *LiquidSession) Allowance(_owner common.Address, _spender common.Address) (*big.Int, error) {
	return _Liquid.Contract.Allowance(&_Liquid.CallOpts, _owner, _spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address _owner, address _spender) view returns(uint256)
func (_Liquid *LiquidCallerSession) Allowance(_owner common.Address, _spender common.Address) (*big.Int, error) {
	return _Liquid.Contract.Allowance(&_Liquid.CallOpts, _owner, _spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address _delegator) view returns(uint256)
func (_Liquid *LiquidCaller) BalanceOf(opts *bind.CallOpts, _delegator common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Liquid.contract.Call(opts, &out, "balanceOf", _delegator)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address _delegator) view returns(uint256)
func (_Liquid *LiquidSession) BalanceOf(_delegator common.Address) (*big.Int, error) {
	return _Liquid.Contract.BalanceOf(&_Liquid.CallOpts, _delegator)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address _delegator) view returns(uint256)
func (_Liquid *LiquidCallerSession) BalanceOf(_delegator common.Address) (*big.Int, error) {
	return _Liquid.Contract.BalanceOf(&_Liquid.CallOpts, _delegator)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() pure returns(uint8)
func (_Liquid *LiquidCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _Liquid.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() pure returns(uint8)
func (_Liquid *LiquidSession) Decimals() (uint8, error) {
	return _Liquid.Contract.Decimals(&_Liquid.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() pure returns(uint8)
func (_Liquid *LiquidCallerSession) Decimals() (uint8, error) {
	return _Liquid.Contract.Decimals(&_Liquid.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_Liquid *LiquidCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Liquid.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_Liquid *LiquidSession) Name() (string, error) {
	return _Liquid.Contract.Name(&_Liquid.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_Liquid *LiquidCallerSession) Name() (string, error) {
	return _Liquid.Contract.Name(&_Liquid.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_Liquid *LiquidCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Liquid.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_Liquid *LiquidSession) Symbol() (string, error) {
	return _Liquid.Contract.Symbol(&_Liquid.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_Liquid *LiquidCallerSession) Symbol() (string, error) {
	return _Liquid.Contract.Symbol(&_Liquid.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_Liquid *LiquidCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Liquid.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_Liquid *LiquidSession) TotalSupply() (*big.Int, error) {
	return _Liquid.Contract.TotalSupply(&_Liquid.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_Liquid *LiquidCallerSession) TotalSupply() (*big.Int, error) {
	return _Liquid.Contract.TotalSupply(&_Liquid.CallOpts)
}

// UnclaimedRewards is a free data retrieval call binding the contract method 0x949813b8.
//
// Solidity: function unclaimedRewards(address _account) view returns(uint256)
func (_Liquid *LiquidCaller) UnclaimedRewards(opts *bind.CallOpts, _account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Liquid.contract.Call(opts, &out, "unclaimedRewards", _account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// UnclaimedRewards is a free data retrieval call binding the contract method 0x949813b8.
//
// Solidity: function unclaimedRewards(address _account) view returns(uint256)
func (_Liquid *LiquidSession) UnclaimedRewards(_account common.Address) (*big.Int, error) {
	return _Liquid.Contract.UnclaimedRewards(&_Liquid.CallOpts, _account)
}

// UnclaimedRewards is a free data retrieval call binding the contract method 0x949813b8.
//
// Solidity: function unclaimedRewards(address _account) view returns(uint256)
func (_Liquid *LiquidCallerSession) UnclaimedRewards(_account common.Address) (*big.Int, error) {
	return _Liquid.Contract.UnclaimedRewards(&_Liquid.CallOpts, _account)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address _spender, uint256 _amount) returns(bool)
func (_Liquid *LiquidTransactor) Approve(opts *bind.TransactOpts, _spender common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Liquid.contract.Transact(opts, "approve", _spender, _amount)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address _spender, uint256 _amount) returns(bool)
func (_Liquid *LiquidSession) Approve(_spender common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Liquid.Contract.Approve(&_Liquid.TransactOpts, _spender, _amount)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address _spender, uint256 _amount) returns(bool)
func (_Liquid *LiquidTransactorSession) Approve(_spender common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Liquid.Contract.Approve(&_Liquid.TransactOpts, _spender, _amount)
}

// Burn is a paid mutator transaction binding the contract method 0x9dc29fac.
//
// Solidity: function burn(address _account, uint256 _amount) returns()
func (_Liquid *LiquidTransactor) Burn(opts *bind.TransactOpts, _account common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Liquid.contract.Transact(opts, "burn", _account, _amount)
}

// Burn is a paid mutator transaction binding the contract method 0x9dc29fac.
//
// Solidity: function burn(address _account, uint256 _amount) returns()
func (_Liquid *LiquidSession) Burn(_account common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Liquid.Contract.Burn(&_Liquid.TransactOpts, _account, _amount)
}

// Burn is a paid mutator transaction binding the contract method 0x9dc29fac.
//
// Solidity: function burn(address _account, uint256 _amount) returns()
func (_Liquid *LiquidTransactorSession) Burn(_account common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Liquid.Contract.Burn(&_Liquid.TransactOpts, _account, _amount)
}

// ClaimRewards is a paid mutator transaction binding the contract method 0x372500ab.
//
// Solidity: function claimRewards() returns()
func (_Liquid *LiquidTransactor) ClaimRewards(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Liquid.contract.Transact(opts, "claimRewards")
}

// ClaimRewards is a paid mutator transaction binding the contract method 0x372500ab.
//
// Solidity: function claimRewards() returns()
func (_Liquid *LiquidSession) ClaimRewards() (*types.Transaction, error) {
	return _Liquid.Contract.ClaimRewards(&_Liquid.TransactOpts)
}

// ClaimRewards is a paid mutator transaction binding the contract method 0x372500ab.
//
// Solidity: function claimRewards() returns()
func (_Liquid *LiquidTransactorSession) ClaimRewards() (*types.Transaction, error) {
	return _Liquid.Contract.ClaimRewards(&_Liquid.TransactOpts)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address _account, uint256 _amount) returns()
func (_Liquid *LiquidTransactor) Mint(opts *bind.TransactOpts, _account common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Liquid.contract.Transact(opts, "mint", _account, _amount)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address _account, uint256 _amount) returns()
func (_Liquid *LiquidSession) Mint(_account common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Liquid.Contract.Mint(&_Liquid.TransactOpts, _account, _amount)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address _account, uint256 _amount) returns()
func (_Liquid *LiquidTransactorSession) Mint(_account common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Liquid.Contract.Mint(&_Liquid.TransactOpts, _account, _amount)
}

// Redistribute is a paid mutator transaction binding the contract method 0xfb489a7b.
//
// Solidity: function redistribute() payable returns(uint256)
func (_Liquid *LiquidTransactor) Redistribute(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Liquid.contract.Transact(opts, "redistribute")
}

// Redistribute is a paid mutator transaction binding the contract method 0xfb489a7b.
//
// Solidity: function redistribute() payable returns(uint256)
func (_Liquid *LiquidSession) Redistribute() (*types.Transaction, error) {
	return _Liquid.Contract.Redistribute(&_Liquid.TransactOpts)
}

// Redistribute is a paid mutator transaction binding the contract method 0xfb489a7b.
//
// Solidity: function redistribute() payable returns(uint256)
func (_Liquid *LiquidTransactorSession) Redistribute() (*types.Transaction, error) {
	return _Liquid.Contract.Redistribute(&_Liquid.TransactOpts)
}

// SetCommissionRate is a paid mutator transaction binding the contract method 0x19fac8fd.
//
// Solidity: function setCommissionRate(uint256 _rate) returns()
func (_Liquid *LiquidTransactor) SetCommissionRate(opts *bind.TransactOpts, _rate *big.Int) (*types.Transaction, error) {
	return _Liquid.contract.Transact(opts, "setCommissionRate", _rate)
}

// SetCommissionRate is a paid mutator transaction binding the contract method 0x19fac8fd.
//
// Solidity: function setCommissionRate(uint256 _rate) returns()
func (_Liquid *LiquidSession) SetCommissionRate(_rate *big.Int) (*types.Transaction, error) {
	return _Liquid.Contract.SetCommissionRate(&_Liquid.TransactOpts, _rate)
}

// SetCommissionRate is a paid mutator transaction binding the contract method 0x19fac8fd.
//
// Solidity: function setCommissionRate(uint256 _rate) returns()
func (_Liquid *LiquidTransactorSession) SetCommissionRate(_rate *big.Int) (*types.Transaction, error) {
	return _Liquid.Contract.SetCommissionRate(&_Liquid.TransactOpts, _rate)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address _to, uint256 _amount) returns(bool _success)
func (_Liquid *LiquidTransactor) Transfer(opts *bind.TransactOpts, _to common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Liquid.contract.Transact(opts, "transfer", _to, _amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address _to, uint256 _amount) returns(bool _success)
func (_Liquid *LiquidSession) Transfer(_to common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Liquid.Contract.Transfer(&_Liquid.TransactOpts, _to, _amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address _to, uint256 _amount) returns(bool _success)
func (_Liquid *LiquidTransactorSession) Transfer(_to common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Liquid.Contract.Transfer(&_Liquid.TransactOpts, _to, _amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address _sender, address _recipient, uint256 _amount) returns(bool _success)
func (_Liquid *LiquidTransactor) TransferFrom(opts *bind.TransactOpts, _sender common.Address, _recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Liquid.contract.Transact(opts, "transferFrom", _sender, _recipient, _amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address _sender, address _recipient, uint256 _amount) returns(bool _success)
func (_Liquid *LiquidSession) TransferFrom(_sender common.Address, _recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Liquid.Contract.TransferFrom(&_Liquid.TransactOpts, _sender, _recipient, _amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address _sender, address _recipient, uint256 _amount) returns(bool _success)
func (_Liquid *LiquidTransactorSession) TransferFrom(_sender common.Address, _recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Liquid.Contract.TransferFrom(&_Liquid.TransactOpts, _sender, _recipient, _amount)
}

// LiquidApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the Liquid contract.
type LiquidApprovalIterator struct {
	Event *LiquidApproval // Event containing the contract specifics and raw log

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
func (it *LiquidApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LiquidApproval)
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
		it.Event = new(LiquidApproval)
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
func (it *LiquidApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LiquidApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LiquidApproval represents a Approval event raised by the Liquid contract.
type LiquidApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_Liquid *LiquidFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*LiquidApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _Liquid.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &LiquidApprovalIterator{contract: _Liquid.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_Liquid *LiquidFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *LiquidApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _Liquid.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LiquidApproval)
				if err := _Liquid.contract.UnpackLog(event, "Approval", log); err != nil {
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
func (_Liquid *LiquidFilterer) ParseApproval(log types.Log) (*LiquidApproval, error) {
	event := new(LiquidApproval)
	if err := _Liquid.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LiquidTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the Liquid contract.
type LiquidTransferIterator struct {
	Event *LiquidTransfer // Event containing the contract specifics and raw log

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
func (it *LiquidTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LiquidTransfer)
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
		it.Event = new(LiquidTransfer)
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
func (it *LiquidTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LiquidTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LiquidTransfer represents a Transfer event raised by the Liquid contract.
type LiquidTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_Liquid *LiquidFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*LiquidTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Liquid.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &LiquidTransferIterator{contract: _Liquid.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_Liquid *LiquidFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *LiquidTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Liquid.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LiquidTransfer)
				if err := _Liquid.contract.UnpackLog(event, "Transfer", log); err != nil {
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
func (_Liquid *LiquidFilterer) ParseTransfer(log types.Log) (*LiquidTransfer, error) {
	event := new(LiquidTransfer)
	if err := _Liquid.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OracleMetaData contains all meta data concerning the Oracle contract.
var OracleMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_voters\",\"type\":\"address[]\"},{\"internalType\":\"address\",\"name\":\"_autonity\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_operator\",\"type\":\"address\"},{\"internalType\":\"string[]\",\"name\":\"_symbols\",\"type\":\"string[]\"},{\"internalType\":\"uint256\",\"name\":\"_votePeriod\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_round\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_height\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_timestamp\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_votePeriod\",\"type\":\"uint256\"}],\"name\":\"NewRound\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string[]\",\"name\":\"_symbols\",\"type\":\"string[]\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_round\",\"type\":\"uint256\"}],\"name\":\"NewSymbols\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_voter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"int256[]\",\"name\":\"_votes\",\"type\":\"int256[]\"}],\"name\":\"Voted\",\"type\":\"event\"},{\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"inputs\":[],\"name\":\"finalize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPrecision\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRound\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_round\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"_symbol\",\"type\":\"string\"}],\"name\":\"getRoundData\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"round\",\"type\":\"uint256\"},{\"internalType\":\"int256\",\"name\":\"price\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"status\",\"type\":\"uint256\"}],\"internalType\":\"structIOracle.RoundData\",\"name\":\"data\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getSymbols\",\"outputs\":[{\"internalType\":\"string[]\",\"name\":\"\",\"type\":\"string[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getVotePeriod\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getVoters\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastRoundBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastVoterUpdateRound\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_symbol\",\"type\":\"string\"}],\"name\":\"latestRoundData\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"round\",\"type\":\"uint256\"},{\"internalType\":\"int256\",\"name\":\"price\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"status\",\"type\":\"uint256\"}],\"internalType\":\"structIOracle.RoundData\",\"name\":\"data\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"newSymbols\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"reports\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"round\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_operator\",\"type\":\"address\"}],\"name\":\"setOperator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string[]\",\"name\":\"_symbols\",\"type\":\"string[]\"}],\"name\":\"setSymbols\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_newVoters\",\"type\":\"address[]\"}],\"name\":\"setVoters\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbolUpdatedRound\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"symbols\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_commit\",\"type\":\"uint256\"},{\"internalType\":\"int256[]\",\"name\":\"_reports\",\"type\":\"int256[]\"},{\"internalType\":\"uint256\",\"name\":\"_salt\",\"type\":\"uint256\"}],\"name\":\"vote\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"votePeriod\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"votingInfo\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"round\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"commit\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"isVoter\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Sigs: map[string]string{
		"4bb278f3": "finalize()",
		"9670c0bc": "getPrecision()",
		"9f8743f7": "getRound()",
		"3c8510fd": "getRoundData(uint256,string)",
		"df7f710e": "getSymbols()",
		"b78dec52": "getVotePeriod()",
		"cdd72253": "getVoters()",
		"e6a02a28": "lastRoundBlock()",
		"aa2f89b5": "lastVoterUpdateRound()",
		"33f98c77": "latestRoundData(string)",
		"5281b5c6": "newSymbols(uint256)",
		"4c56ea56": "reports(string,address)",
		"146ca531": "round()",
		"b3ab15fb": "setOperator(address)",
		"8d4f75d2": "setSymbols(string[])",
		"845023f2": "setVoters(address[])",
		"08f21ff5": "symbolUpdatedRound()",
		"ccce413b": "symbols(uint256)",
		"307de9b6": "vote(uint256,int256[],uint256)",
		"a7813587": "votePeriod()",
		"5412b3ae": "votingInfo(address)",
	},
	Bin: "0x6080604052600160ff1b600755600160ff1b6008553480156200002157600080fd5b5060405162002d1b38038062002d1b8339810160408190526200004491620006cb565b600280546001600160a01b038087166001600160a01b031992831617909255600380549286169290911691909117905581516200008990600090602085019062000362565b5081516200009f90600190602085019062000362565b5080600981905550620000c485600060018851620000be9190620007e0565b62000181565b8451620000d9906004906020880190620003c6565b508451620000ef906005906020880190620003c6565b5060016006819055600d8054909101815560009081525b855181101562000175576001600b60008884815181106200012b576200012b620007fa565b6020908102919091018101516001600160a01b03168252810191909152604001600020600201805460ff1916911515919091179055806200016c8162000810565b91505062000106565b50505050505062000974565b81818082141562000193575050505050565b6000856002620001a487876200082e565b620001b0919062000873565b620001bc9087620008b3565b81518110620001cf57620001cf620007fa565b602002602001015190505b8183136200032e575b806001600160a01b0316868481518110620002025762000202620007fa565b60200260200101516001600160a01b031611156200022f57826200022681620008fa565b935050620001e3565b858281518110620002445762000244620007fa565b60200260200101516001600160a01b0316816001600160a01b031611156200027b5781620002728162000916565b9250506200022f565b8183136200032857858281518110620002985762000298620007fa565b6020026020010151868481518110620002b557620002b5620007fa565b6020026020010151878581518110620002d257620002d2620007fa565b60200260200101888581518110620002ee57620002ee620007fa565b6001600160a01b03938416602091820292909201015291169052826200031481620008fa565b9350508180620003249062000916565b9250505b620001da565b8185121562000344576200034486868462000181565b838312156200035a576200035a86848662000181565b505050505050565b828054828255906000526020600020908101928215620003b4579160200282015b82811115620003b45782518051620003a39184916020909101906200042c565b509160200191906001019062000383565b50620003c2929150620004a9565b5090565b8280548282559060005260206000209081019282156200041e579160200282015b828111156200041e57825182546001600160a01b0319166001600160a01b03909116178255602090920191600190910190620003e7565b50620003c2929150620004ca565b8280546200043a9062000937565b90600052602060002090601f0160209004810192826200045e57600085556200041e565b82601f106200047957805160ff19168380011785556200041e565b828001600101855582156200041e579182015b828111156200041e5782518255916020019190600101906200048c565b80821115620003c2576000620004c08282620004e1565b50600101620004a9565b5b80821115620003c25760008155600101620004cb565b508054620004ef9062000937565b6000825580601f1062000500575050565b601f016020900490600052602060002090810190620005209190620004ca565b50565b634e487b7160e01b600052604160045260246000fd5b604051601f8201601f191681016001600160401b038111828210171562000564576200056462000523565b604052919050565b60006001600160401b0382111562000588576200058862000523565b5060051b60200190565b80516001600160a01b0381168114620005aa57600080fd5b919050565b6000601f8381840112620005c257600080fd5b82516020620005db620005d5836200056c565b62000539565b82815260059290921b85018101918181019087841115620005fb57600080fd5b8287015b84811015620006bf5780516001600160401b0380821115620006215760008081fd5b818a0191508a603f830112620006375760008081fd5b85820151818111156200064e576200064e62000523565b62000661818a01601f1916880162000539565b915080825260408c818386010111156200067b5760008081fd5b60005b828110156200069b578481018201518482018a015288016200067e565b82811115620006ad5760008984860101525b505050845250918301918301620005ff565b50979650505050505050565b600080600080600060a08688031215620006e457600080fd5b85516001600160401b0380821115620006fc57600080fd5b818801915088601f8301126200071157600080fd5b8151602062000724620005d5836200056c565b82815260059290921b8401810191818101908c8411156200074457600080fd5b948201945b838610156200076d576200075d8662000592565b8252948201949082019062000749565b99506200077e90508a820162000592565b97505050620007906040890162000592565b94506060880151915080821115620007a757600080fd5b50620007b688828901620005af565b925050608086015190509295509295909350565b634e487b7160e01b600052601160045260246000fd5b600082821015620007f557620007f5620007ca565b500390565b634e487b7160e01b600052603260045260246000fd5b6000600019821415620008275762000827620007ca565b5060010190565b60008083128015600160ff1b8501841216156200084f576200084f620007ca565b6001600160ff1b03840183138116156200086d576200086d620007ca565b50500390565b6000826200089157634e487b7160e01b600052601260045260246000fd5b600160ff1b821460001984141615620008ae57620008ae620007ca565b500590565b600080821280156001600160ff1b0384900385131615620008d857620008d8620007ca565b600160ff1b8390038412811615620008f457620008f4620007ca565b50500190565b60006001600160ff1b03821415620008275762000827620007ca565b6000600160ff1b8214156200092f576200092f620007ca565b506000190190565b600181811c908216806200094c57607f821691505b602082108114156200096e57634e487b7160e01b600052602260045260246000fd5b50919050565b61239780620009846000396000f3fe6080604052600436106101225760003560e01c80638d4f75d2116100a5578063b3ab15fb1161006c578063b3ab15fb1461037a578063b78dec521461039a578063ccce413b146103af578063cdd72253146103cf578063df7f710e146103f1578063e6a02a281461041357005b80638d4f75d2146103035780639670c0bc146103235780639f8743f714610339578063a78135871461034e578063aa2f89b51461036457005b80634bb278f3116100e95780634bb278f3146101fd5780634c56ea56146102125780635281b5c61461025a5780635412b3ae14610287578063845023f2146102e357005b806308f21ff51461012b578063146ca53114610154578063307de9b61461016a57806333f98c771461018a5780633c8510fd146101dd57005b3661012957005b005b34801561013757600080fd5b5061014160085481565b6040519081526020015b60405180910390f35b34801561016057600080fd5b5061014160065481565b34801561017657600080fd5b50610129610185366004611b20565b610429565b34801561019657600080fd5b506101aa6101a5366004611c5c565b61066f565b60405161014b91908151815260208083015190820152604080830151908201526060918201519181019190915260800190565b3480156101e957600080fd5b506101aa6101f8366004611c91565b610792565b34801561020957600080fd5b5061012961089c565b34801561021e57600080fd5b5061014161022d366004611cf4565b8151602081840181018051600c825292820194820194909420919093529091526000908152604090205481565b34801561026657600080fd5b5061027a610275366004611d42565b610a36565b60405161014b9190611db3565b34801561029357600080fd5b506102c66102a2366004611dcd565b600b6020526000908152604090208054600182015460029092015490919060ff1683565b60408051938452602084019290925215159082015260600161014b565b3480156102ef57600080fd5b506101296102fe366004611e0c565b610ae2565b34801561030f57600080fd5b5061012961031e366004611ea9565b610b88565b34801561032f57600080fd5b5062989680610141565b34801561034557600080fd5b50600654610141565b34801561035a57600080fd5b5061014160095481565b34801561037057600080fd5b5061014160075481565b34801561038657600080fd5b50610129610395366004611dcd565b610cf5565b3480156103a657600080fd5b50600954610141565b3480156103bb57600080fd5b5061027a6103ca366004611d42565b610d41565b3480156103db57600080fd5b506103e4610d51565b60405161014b9190611f5a565b3480156103fd57600080fd5b50610406610db3565b60405161014b9190611ffc565b34801561041f57600080fd5b50610141600a5481565b336000908152600b602052604090206002015460ff166104905760405162461bcd60e51b815260206004820152601960248201527f7265737472696374656420746f206f6e6c7920766f746572730000000000000060448201526064015b60405180910390fd5b600654336000908152600b602052604090205414156104e15760405162461bcd60e51b815260206004820152600d60248201526c185b1c9958591e481d9bdd1959609a1b6044820152606401610487565b336000908152600b60205260409020600181018054908690558154600654909255908061050f575050610669565b600054841461051f575050610669565b600160065461052e9190612025565b8114158061056b57508484843360405160200161054e949392919061203c565b6040516020818303038152906040528051906020012060001c8214155b156105e75760005b6000548110156105df576001600160ff1b03600c6000838154811061059a5761059a612087565b906000526020600020016040516105b191906120d8565b90815260408051602092819003830190203360009081529252902055806105d781612174565b915050610573565b505050610669565b60005b848110156106655785858281811061060457610604612087565b90506020020135600c6000838154811061062057610620612087565b9060005260206000200160405161063791906120d8565b908152604080516020928190038301902033600090815292529020558061065d81612174565b9150506105ea565b5050505b50505050565b61069a6040518060800160405280600081526020016000815260200160008152602001600081525090565b6000600d60016006546106ad9190612025565b815481106106bd576106bd612087565b90600052602060002001836040516106d5919061218f565b908152602001604051809103902060405180606001604052908160008201548152602001600182015481526020016002820160009054906101000a900460ff166001811115610726576107266121ab565b6001811115610737576107376121ab565b8152505090506000604051806080016040528060016006546107599190612025565b8152602001836000015181526020018360200151815260200183604001516001811115610788576107886121ab565b9052949350505050565b6107bd6040518060800160405280600081526020016000815260200160008152602001600081525090565b6000600d84815481106107d2576107d2612087565b90600052602060002001836040516107ea919061218f565b908152602001604051809103902060405180606001604052908160008201548152602001600182015481526020016002820160009054906101000a900460ff16600181111561083b5761083b6121ab565b600181111561084c5761084c6121ab565b8152505090506000604051806080016040528086815260200183600001518152602001836020015181526020018360400151600181111561088f5761088f6121ab565b9052925050505b92915050565b6002546001600160a01b031633146108c65760405162461bcd60e51b8152600401610487906121c1565b600954600a546108d69190612204565b4310610a345760005b600054811015610904576108f281610f72565b6108fd600182612204565b90506108df565b5060065460075414156109825760005b600554811015610980576001600b60006005848154811061093757610937612087565b6000918252602080832091909101546001600160a01b031683528201929092526040019020600201805460ff19169115159190911790558061097881612174565b915050610914565b505b60065460075461099390600161221c565b14156109a1576109a16112a7565b43600a819055506001600660008282546109bb9190612204565b90915550506008546109ce90600261221c565b60065414156109e957600180546109e79160009161186f565b505b60065460095460408051928352436020840152429083015260608201527fb5d8636ab45e6cac7a4a61cb7c77f77f61a454d73aa2e6139ff8dcaf463537e59060800160405180910390a15b565b60018181548110610a4657600080fd5b906000526020600020016000915090508054610a619061209d565b80601f0160208091040260200160405190810160405280929190818152602001828054610a8d9061209d565b8015610ada5780601f10610aaf57610100808354040283529160200191610ada565b820191906000526020600020905b815481529060010190602001808311610abd57829003601f168201915b505050505081565b6002546001600160a01b03163314610b0c5760405162461bcd60e51b8152600401610487906121c1565b8051610b525760405162461bcd60e51b8152602060048201526015602482015274566f746572732063616e277420626520656d70747960581b6044820152606401610487565b610b6b81600060018451610b669190612025565b611485565b8051610b7e9060059060208401906118d5565b5050600654600755565b6003546001600160a01b03163314610bdb5760405162461bcd60e51b81526020600482015260166024820152753932b9ba3934b1ba32b2103a379037b832b930ba37b960511b6044820152606401610487565b8051610c225760405162461bcd60e51b815260206004820152601660248201527573796d626f6c732063616e277420626520656d70747960501b6044820152606401610487565b600654600854610c3390600161221c565b14158015610c45575060065460085414155b610c915760405162461bcd60e51b815260206004820152601e60248201527f63616e2774206265207570646174656420696e207468697320726f756e6400006044820152606401610487565b8051610ca4906001906020840190611936565b5060065460088190557faa278e424da680ce5dad66510415760e78e0bd87d45c786c6e88bdde82f9342d908290610cdc906001612204565b604051610cea92919061225d565b60405180910390a150565b6002546001600160a01b03163314610d1f5760405162461bcd60e51b8152600401610487906121c1565b600380546001600160a01b0319166001600160a01b0392909216919091179055565b60008181548110610a4657600080fd5b60606005805480602002602001604051908101604052809291908181526020018280548015610da957602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311610d8b575b5050505050905090565b60606006546008546001610dc7919061221c565b1415610ea4576001805480602002602001604051908101604052809291908181526020016000905b82821015610e9b578382906000526020600020018054610e0e9061209d565b80601f0160208091040260200160405190810160405280929190818152602001828054610e3a9061209d565b8015610e875780601f10610e5c57610100808354040283529160200191610e87565b820191906000526020600020905b815481529060010190602001808311610e6a57829003601f168201915b505050505081526020019060010190610def565b50505050905090565b6000805480602002602001604051908101604052809291908181526020016000905b82821015610e9b578382906000526020600020018054610ee59061209d565b80601f0160208091040260200160405190810160405280929190818152602001828054610f119061209d565b8015610f5e5780601f10610f3357610100808354040283529160200191610f5e565b820191906000526020600020905b815481529060010190602001808311610f4157829003601f168201915b505050505081526020019060010190610ec6565b6000808281548110610f8657610f86612087565b906000526020600020018054610f9b9061209d565b80601f0160208091040260200160405190810160405280929190818152602001828054610fc79061209d565b80156110145780601f10610fe957610100808354040283529160200191611014565b820191906000526020600020905b815481529060010190602001808311610ff757829003601f168201915b50505050509050600060048054905067ffffffffffffffff81111561103b5761103b611ba5565b604051908082528060200260200182016040528015611064578160200160208202803683370190505b5090506000805b60045481101561117b5760006004828154811061108a5761108a612087565b60009182526020808320909101546006546001600160a01b03909116808452600b9092526040909220549092501415806110fe57506001600160ff1b03600c866040516110d7919061218f565b90815260408051602092819003830190206001600160a01b03851660009081529252902054145b156111095750611169565b600c85604051611119919061218f565b90815260408051602092819003830190206001600160a01b03841660009081529252902054848461114981612174565b95508151811061115b5761115b612087565b602002602001018181525050505b8061117381612174565b91505061106b565b506000600d600160065461118f9190612025565b8154811061119f5761119f612087565b90600052602060002001846040516111b7919061218f565b908152604051908190036020019020549050600182156111e2576111db8484611635565b9150600090505b600d805460019081018255600091909152604080516060810182528481524260208201529190820190839081111561121c5761121c6121ab565b815250600d6006548154811061123457611234612087565b906000526020600020018660405161124c919061218f565b9081526020016040518091039020600082015181600001556020820151816001015560408201518160020160006101000a81548160ff02191690836001811115611298576112986121ab565b02179055505050505050505050565b6000805b600454821080156112bd575060055481105b1561140357600581815481106112d5576112d5612087565b600091825260209091200154600480546001600160a01b03909216918490811061130157611301612087565b6000918252602090912001546001600160a01b0316141561133c578161132681612174565b925050808061133490612174565b9150506112ab565b6005818154811061134f5761134f612087565b600091825260209091200154600480546001600160a01b03909216918490811061137b5761137b612087565b6000918252602090912001546001600160a01b031610156113f957600b6000600484815481106113ad576113ad612087565b60009182526020808320909101546001600160a01b0316835282019290925260400181208181556001810191909155600201805460ff19169055816113f181612174565b9250506112ab565b8061133481612174565b60045482101561147057600b60006004848154811061142457611424612087565b60009182526020808320909101546001600160a01b0316835282019290925260400181208181556001810191909155600201805460ff191690558161146881612174565b925050611403565b6005805461148091600491611983565b505050565b818180821415611496575050505050565b60008560026114a5878761227f565b6114af91906122d4565b6114b9908761221c565b815181106114c9576114c9612087565b602002602001015190505b818313611607575b806001600160a01b03168684815181106114f8576114f8612087565b60200260200101516001600160a01b03161115611521578261151981612302565b9350506114dc565b85828151811061153357611533612087565b60200260200101516001600160a01b0316816001600160a01b03161115611566578161155e8161231b565b925050611521565b8183136116025785828151811061157f5761157f612087565b602002602001015186848151811061159957611599612087565b60200260200101518785815181106115b3576115b3612087565b602002602001018885815181106115cc576115cc612087565b6001600160a01b03938416602091820292909201015291169052826115f081612302565b93505081806115fe9061231b565b9250505b6114d4565b8185121561161a5761161a868684611485565b8383121561162d5761162d868486611485565b505050505050565b60008161164457506000610896565b61165a836000611655600186612025565b6116f6565b6000611667600284612339565b905061167460028461234d565b156116985783818151811061168b5761168b612087565b60200260200101516116ee565b60028482815181106116ac576116ac612087565b6020026020010151856001846116c29190612025565b815181106116d2576116d2612087565b60200260200101516116e4919061221c565b6116ee91906122d4565b949350505050565b818180821415611707575050505050565b6000856002611716878761227f565b61172091906122d4565b61172a908761221c565b8151811061173a5761173a612087565b602002602001015190505b818313611849575b8086848151811061176057611760612087565b60200260200101511215611780578261177881612302565b93505061174d565b85828151811061179257611792612087565b60200260200101518112156117b357816117ab8161231b565b925050611780565b818313611844578582815181106117cc576117cc612087565b60200260200101518684815181106117e6576117e6612087565b602002602001015187858151811061180057611800612087565b6020026020010188858151811061181957611819612087565b602090810291909101019190915252816118328161231b565b925050828061184090612302565b9350505b611745565b8185121561185c5761185c8686846116f6565b8383121561162d5761162d8684866116f6565b8280548282559060005260206000209081019282156118c55760005260206000209182015b828111156118c55782829080546118aa9061209d565b6118b59291906119c3565b5091600101919060010190611894565b506118d1929150611a3d565b5090565b82805482825590600052602060002090810192821561192a579160200282015b8281111561192a57825182546001600160a01b0319166001600160a01b039091161782556020909201916001909101906118f5565b506118d1929150611a5a565b8280548282559060005260206000209081019282156118c5579160200282015b828111156118c55782518051611973918491602090910190611a6f565b5091602001919060010190611956565b82805482825590600052602060002090810192821561192a5760005260206000209182015b8281111561192a5782548255916001019190600101906119a8565b8280546119cf9061209d565b90600052602060002090601f0160209004810192826119f1576000855561192a565b82601f10611a02578054855561192a565b8280016001018555821561192a57600052602060002091601f016020900482018281111561192a5782548255916001019190600101906119a8565b808211156118d1576000611a518282611ae3565b50600101611a3d565b5b808211156118d15760008155600101611a5b565b828054611a7b9061209d565b90600052602060002090601f016020900481019282611a9d576000855561192a565b82601f10611ab657805160ff191683800117855561192a565b8280016001018555821561192a579182015b8281111561192a578251825591602001919060010190611ac8565b508054611aef9061209d565b6000825580601f10611aff575050565b601f016020900490600052602060002090810190611b1d9190611a5a565b50565b60008060008060608587031215611b3657600080fd5b84359350602085013567ffffffffffffffff80821115611b5557600080fd5b818701915087601f830112611b6957600080fd5b813581811115611b7857600080fd5b8860208260051b8501011115611b8d57600080fd5b95986020929092019750949560400135945092505050565b634e487b7160e01b600052604160045260246000fd5b604051601f8201601f1916810167ffffffffffffffff81118282101715611be457611be4611ba5565b604052919050565b600082601f830112611bfd57600080fd5b813567ffffffffffffffff811115611c1757611c17611ba5565b611c2a601f8201601f1916602001611bbb565b818152846020838601011115611c3f57600080fd5b816020850160208301376000918101602001919091529392505050565b600060208284031215611c6e57600080fd5b813567ffffffffffffffff811115611c8557600080fd5b6116ee84828501611bec565b60008060408385031215611ca457600080fd5b82359150602083013567ffffffffffffffff811115611cc257600080fd5b611cce85828601611bec565b9150509250929050565b80356001600160a01b0381168114611cef57600080fd5b919050565b60008060408385031215611d0757600080fd5b823567ffffffffffffffff811115611d1e57600080fd5b611d2a85828601611bec565b925050611d3960208401611cd8565b90509250929050565b600060208284031215611d5457600080fd5b5035919050565b60005b83811015611d76578181015183820152602001611d5e565b838111156106695750506000910152565b60008151808452611d9f816020860160208601611d5b565b601f01601f19169290920160200192915050565b602081526000611dc66020830184611d87565b9392505050565b600060208284031215611ddf57600080fd5b611dc682611cd8565b600067ffffffffffffffff821115611e0257611e02611ba5565b5060051b60200190565b60006020808385031215611e1f57600080fd5b823567ffffffffffffffff811115611e3657600080fd5b8301601f81018513611e4757600080fd5b8035611e5a611e5582611de8565b611bbb565b81815260059190911b82018301908381019087831115611e7957600080fd5b928401925b82841015611e9e57611e8f84611cd8565b82529284019290840190611e7e565b979650505050505050565b60006020808385031215611ebc57600080fd5b823567ffffffffffffffff80821115611ed457600080fd5b818501915085601f830112611ee857600080fd5b8135611ef6611e5582611de8565b81815260059190911b83018401908481019088831115611f1557600080fd5b8585015b83811015611f4d57803585811115611f315760008081fd5b611f3f8b89838a0101611bec565b845250918601918601611f19565b5098975050505050505050565b6020808252825182820181905260009190848201906040850190845b81811015611f9b5783516001600160a01b031683529284019291840191600101611f76565b50909695505050505050565b600081518084526020808501808196508360051b8101915082860160005b85811015611fef578284038952611fdd848351611d87565b98850198935090840190600101611fc5565b5091979650505050505050565b602081526000611dc66020830184611fa7565b634e487b7160e01b600052601160045260246000fd5b6000828210156120375761203761200f565b500390565b60008186825b87811015612060578135835260209283019290910190600101612042565b5050938452505060601b6bffffffffffffffffffffffff1916602082015260340192915050565b634e487b7160e01b600052603260045260246000fd5b600181811c908216806120b157607f821691505b602082108114156120d257634e487b7160e01b600052602260045260246000fd5b50919050565b600080835481600182811c9150808316806120f457607f831692505b602080841082141561211457634e487b7160e01b86526022600452602486fd5b818015612128576001811461213957612166565b60ff19861689528489019650612166565b60008a81526020902060005b8681101561215e5781548b820152908501908301612145565b505084890196505b509498975050505050505050565b60006000198214156121885761218861200f565b5060010190565b600082516121a1818460208701611d5b565b9190910192915050565b634e487b7160e01b600052602160045260246000fd5b60208082526023908201527f7265737472696374656420746f20746865206175746f6e69747920636f6e74726040820152621858dd60ea1b606082015260800190565b600082198211156122175761221761200f565b500190565b600080821280156001600160ff1b038490038513161561223e5761223e61200f565b600160ff1b83900384128116156122575761225761200f565b50500190565b6040815260006122706040830185611fa7565b90508260208301529392505050565b60008083128015600160ff1b85018412161561229d5761229d61200f565b6001600160ff1b03840183138116156122b8576122b861200f565b50500390565b634e487b7160e01b600052601260045260246000fd5b6000826122e3576122e36122be565b600160ff1b8214600019841416156122fd576122fd61200f565b500590565b60006001600160ff1b038214156121885761218861200f565b6000600160ff1b8214156123315761233161200f565b506000190190565b600082612348576123486122be565b500490565b60008261235c5761235c6122be565b50069056fea26469706673582212203caecc2a661f75ffc5867e884fd0135d92c51861bd74855fe134d60741adefea64736f6c634300080c0033",
}

// OracleABI is the input ABI used to generate the binding from.
// Deprecated: Use OracleMetaData.ABI instead.
var OracleABI = OracleMetaData.ABI

// Deprecated: Use OracleMetaData.Sigs instead.
// OracleFuncSigs maps the 4-byte function signature to its string representation.
var OracleFuncSigs = OracleMetaData.Sigs

// OracleBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use OracleMetaData.Bin instead.
var OracleBin = OracleMetaData.Bin

// DeployOracle deploys a new Ethereum contract, binding an instance of Oracle to it.
func DeployOracle(auth *bind.TransactOpts, backend bind.ContractBackend, _voters []common.Address, _autonity common.Address, _operator common.Address, _symbols []string, _votePeriod *big.Int) (common.Address, *types.Transaction, *Oracle, error) {
	parsed, err := OracleMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(OracleBin), backend, _voters, _autonity, _operator, _symbols, _votePeriod)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Oracle{OracleCaller: OracleCaller{contract: contract}, OracleTransactor: OracleTransactor{contract: contract}, OracleFilterer: OracleFilterer{contract: contract}}, nil
}

// Oracle is an auto generated Go binding around an Ethereum contract.
type Oracle struct {
	OracleCaller     // Read-only binding to the contract
	OracleTransactor // Write-only binding to the contract
	OracleFilterer   // Log filterer for contract events
}

// OracleCaller is an auto generated read-only Go binding around an Ethereum contract.
type OracleCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OracleTransactor is an auto generated write-only Go binding around an Ethereum contract.
type OracleTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OracleFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type OracleFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OracleSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type OracleSession struct {
	Contract     *Oracle           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// OracleCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type OracleCallerSession struct {
	Contract *OracleCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// OracleTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type OracleTransactorSession struct {
	Contract     *OracleTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// OracleRaw is an auto generated low-level Go binding around an Ethereum contract.
type OracleRaw struct {
	Contract *Oracle // Generic contract binding to access the raw methods on
}

// OracleCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type OracleCallerRaw struct {
	Contract *OracleCaller // Generic read-only contract binding to access the raw methods on
}

// OracleTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type OracleTransactorRaw struct {
	Contract *OracleTransactor // Generic write-only contract binding to access the raw methods on
}

// NewOracle creates a new instance of Oracle, bound to a specific deployed contract.
func NewOracle(address common.Address, backend bind.ContractBackend) (*Oracle, error) {
	contract, err := bindOracle(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Oracle{OracleCaller: OracleCaller{contract: contract}, OracleTransactor: OracleTransactor{contract: contract}, OracleFilterer: OracleFilterer{contract: contract}}, nil
}

// NewOracleCaller creates a new read-only instance of Oracle, bound to a specific deployed contract.
func NewOracleCaller(address common.Address, caller bind.ContractCaller) (*OracleCaller, error) {
	contract, err := bindOracle(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OracleCaller{contract: contract}, nil
}

// NewOracleTransactor creates a new write-only instance of Oracle, bound to a specific deployed contract.
func NewOracleTransactor(address common.Address, transactor bind.ContractTransactor) (*OracleTransactor, error) {
	contract, err := bindOracle(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OracleTransactor{contract: contract}, nil
}

// NewOracleFilterer creates a new log filterer instance of Oracle, bound to a specific deployed contract.
func NewOracleFilterer(address common.Address, filterer bind.ContractFilterer) (*OracleFilterer, error) {
	contract, err := bindOracle(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OracleFilterer{contract: contract}, nil
}

// bindOracle binds a generic wrapper to an already deployed contract.
func bindOracle(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(OracleABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Oracle *OracleRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Oracle.Contract.OracleCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Oracle *OracleRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Oracle.Contract.OracleTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Oracle *OracleRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Oracle.Contract.OracleTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Oracle *OracleCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Oracle.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Oracle *OracleTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Oracle.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Oracle *OracleTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Oracle.Contract.contract.Transact(opts, method, params...)
}

// GetPrecision is a free data retrieval call binding the contract method 0x9670c0bc.
//
// Solidity: function getPrecision() pure returns(uint256)
func (_Oracle *OracleCaller) GetPrecision(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Oracle.contract.Call(opts, &out, "getPrecision")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetPrecision is a free data retrieval call binding the contract method 0x9670c0bc.
//
// Solidity: function getPrecision() pure returns(uint256)
func (_Oracle *OracleSession) GetPrecision() (*big.Int, error) {
	return _Oracle.Contract.GetPrecision(&_Oracle.CallOpts)
}

// GetPrecision is a free data retrieval call binding the contract method 0x9670c0bc.
//
// Solidity: function getPrecision() pure returns(uint256)
func (_Oracle *OracleCallerSession) GetPrecision() (*big.Int, error) {
	return _Oracle.Contract.GetPrecision(&_Oracle.CallOpts)
}

// GetRound is a free data retrieval call binding the contract method 0x9f8743f7.
//
// Solidity: function getRound() view returns(uint256)
func (_Oracle *OracleCaller) GetRound(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Oracle.contract.Call(opts, &out, "getRound")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetRound is a free data retrieval call binding the contract method 0x9f8743f7.
//
// Solidity: function getRound() view returns(uint256)
func (_Oracle *OracleSession) GetRound() (*big.Int, error) {
	return _Oracle.Contract.GetRound(&_Oracle.CallOpts)
}

// GetRound is a free data retrieval call binding the contract method 0x9f8743f7.
//
// Solidity: function getRound() view returns(uint256)
func (_Oracle *OracleCallerSession) GetRound() (*big.Int, error) {
	return _Oracle.Contract.GetRound(&_Oracle.CallOpts)
}

// GetRoundData is a free data retrieval call binding the contract method 0x3c8510fd.
//
// Solidity: function getRoundData(uint256 _round, string _symbol) view returns((uint256,int256,uint256,uint256) data)
func (_Oracle *OracleCaller) GetRoundData(opts *bind.CallOpts, _round *big.Int, _symbol string) (IOracleRoundData, error) {
	var out []interface{}
	err := _Oracle.contract.Call(opts, &out, "getRoundData", _round, _symbol)

	if err != nil {
		return *new(IOracleRoundData), err
	}

	out0 := *abi.ConvertType(out[0], new(IOracleRoundData)).(*IOracleRoundData)

	return out0, err

}

// GetRoundData is a free data retrieval call binding the contract method 0x3c8510fd.
//
// Solidity: function getRoundData(uint256 _round, string _symbol) view returns((uint256,int256,uint256,uint256) data)
func (_Oracle *OracleSession) GetRoundData(_round *big.Int, _symbol string) (IOracleRoundData, error) {
	return _Oracle.Contract.GetRoundData(&_Oracle.CallOpts, _round, _symbol)
}

// GetRoundData is a free data retrieval call binding the contract method 0x3c8510fd.
//
// Solidity: function getRoundData(uint256 _round, string _symbol) view returns((uint256,int256,uint256,uint256) data)
func (_Oracle *OracleCallerSession) GetRoundData(_round *big.Int, _symbol string) (IOracleRoundData, error) {
	return _Oracle.Contract.GetRoundData(&_Oracle.CallOpts, _round, _symbol)
}

// GetSymbols is a free data retrieval call binding the contract method 0xdf7f710e.
//
// Solidity: function getSymbols() view returns(string[])
func (_Oracle *OracleCaller) GetSymbols(opts *bind.CallOpts) ([]string, error) {
	var out []interface{}
	err := _Oracle.contract.Call(opts, &out, "getSymbols")

	if err != nil {
		return *new([]string), err
	}

	out0 := *abi.ConvertType(out[0], new([]string)).(*[]string)

	return out0, err

}

// GetSymbols is a free data retrieval call binding the contract method 0xdf7f710e.
//
// Solidity: function getSymbols() view returns(string[])
func (_Oracle *OracleSession) GetSymbols() ([]string, error) {
	return _Oracle.Contract.GetSymbols(&_Oracle.CallOpts)
}

// GetSymbols is a free data retrieval call binding the contract method 0xdf7f710e.
//
// Solidity: function getSymbols() view returns(string[])
func (_Oracle *OracleCallerSession) GetSymbols() ([]string, error) {
	return _Oracle.Contract.GetSymbols(&_Oracle.CallOpts)
}

// GetVotePeriod is a free data retrieval call binding the contract method 0xb78dec52.
//
// Solidity: function getVotePeriod() view returns(uint256)
func (_Oracle *OracleCaller) GetVotePeriod(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Oracle.contract.Call(opts, &out, "getVotePeriod")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetVotePeriod is a free data retrieval call binding the contract method 0xb78dec52.
//
// Solidity: function getVotePeriod() view returns(uint256)
func (_Oracle *OracleSession) GetVotePeriod() (*big.Int, error) {
	return _Oracle.Contract.GetVotePeriod(&_Oracle.CallOpts)
}

// GetVotePeriod is a free data retrieval call binding the contract method 0xb78dec52.
//
// Solidity: function getVotePeriod() view returns(uint256)
func (_Oracle *OracleCallerSession) GetVotePeriod() (*big.Int, error) {
	return _Oracle.Contract.GetVotePeriod(&_Oracle.CallOpts)
}

// GetVoters is a free data retrieval call binding the contract method 0xcdd72253.
//
// Solidity: function getVoters() view returns(address[])
func (_Oracle *OracleCaller) GetVoters(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _Oracle.contract.Call(opts, &out, "getVoters")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// GetVoters is a free data retrieval call binding the contract method 0xcdd72253.
//
// Solidity: function getVoters() view returns(address[])
func (_Oracle *OracleSession) GetVoters() ([]common.Address, error) {
	return _Oracle.Contract.GetVoters(&_Oracle.CallOpts)
}

// GetVoters is a free data retrieval call binding the contract method 0xcdd72253.
//
// Solidity: function getVoters() view returns(address[])
func (_Oracle *OracleCallerSession) GetVoters() ([]common.Address, error) {
	return _Oracle.Contract.GetVoters(&_Oracle.CallOpts)
}

// LastRoundBlock is a free data retrieval call binding the contract method 0xe6a02a28.
//
// Solidity: function lastRoundBlock() view returns(uint256)
func (_Oracle *OracleCaller) LastRoundBlock(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Oracle.contract.Call(opts, &out, "lastRoundBlock")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LastRoundBlock is a free data retrieval call binding the contract method 0xe6a02a28.
//
// Solidity: function lastRoundBlock() view returns(uint256)
func (_Oracle *OracleSession) LastRoundBlock() (*big.Int, error) {
	return _Oracle.Contract.LastRoundBlock(&_Oracle.CallOpts)
}

// LastRoundBlock is a free data retrieval call binding the contract method 0xe6a02a28.
//
// Solidity: function lastRoundBlock() view returns(uint256)
func (_Oracle *OracleCallerSession) LastRoundBlock() (*big.Int, error) {
	return _Oracle.Contract.LastRoundBlock(&_Oracle.CallOpts)
}

// LastVoterUpdateRound is a free data retrieval call binding the contract method 0xaa2f89b5.
//
// Solidity: function lastVoterUpdateRound() view returns(int256)
func (_Oracle *OracleCaller) LastVoterUpdateRound(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Oracle.contract.Call(opts, &out, "lastVoterUpdateRound")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LastVoterUpdateRound is a free data retrieval call binding the contract method 0xaa2f89b5.
//
// Solidity: function lastVoterUpdateRound() view returns(int256)
func (_Oracle *OracleSession) LastVoterUpdateRound() (*big.Int, error) {
	return _Oracle.Contract.LastVoterUpdateRound(&_Oracle.CallOpts)
}

// LastVoterUpdateRound is a free data retrieval call binding the contract method 0xaa2f89b5.
//
// Solidity: function lastVoterUpdateRound() view returns(int256)
func (_Oracle *OracleCallerSession) LastVoterUpdateRound() (*big.Int, error) {
	return _Oracle.Contract.LastVoterUpdateRound(&_Oracle.CallOpts)
}

// LatestRoundData is a free data retrieval call binding the contract method 0x33f98c77.
//
// Solidity: function latestRoundData(string _symbol) view returns((uint256,int256,uint256,uint256) data)
func (_Oracle *OracleCaller) LatestRoundData(opts *bind.CallOpts, _symbol string) (IOracleRoundData, error) {
	var out []interface{}
	err := _Oracle.contract.Call(opts, &out, "latestRoundData", _symbol)

	if err != nil {
		return *new(IOracleRoundData), err
	}

	out0 := *abi.ConvertType(out[0], new(IOracleRoundData)).(*IOracleRoundData)

	return out0, err

}

// LatestRoundData is a free data retrieval call binding the contract method 0x33f98c77.
//
// Solidity: function latestRoundData(string _symbol) view returns((uint256,int256,uint256,uint256) data)
func (_Oracle *OracleSession) LatestRoundData(_symbol string) (IOracleRoundData, error) {
	return _Oracle.Contract.LatestRoundData(&_Oracle.CallOpts, _symbol)
}

// LatestRoundData is a free data retrieval call binding the contract method 0x33f98c77.
//
// Solidity: function latestRoundData(string _symbol) view returns((uint256,int256,uint256,uint256) data)
func (_Oracle *OracleCallerSession) LatestRoundData(_symbol string) (IOracleRoundData, error) {
	return _Oracle.Contract.LatestRoundData(&_Oracle.CallOpts, _symbol)
}

// NewSymbols is a free data retrieval call binding the contract method 0x5281b5c6.
//
// Solidity: function newSymbols(uint256 ) view returns(string)
func (_Oracle *OracleCaller) NewSymbols(opts *bind.CallOpts, arg0 *big.Int) (string, error) {
	var out []interface{}
	err := _Oracle.contract.Call(opts, &out, "newSymbols", arg0)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// NewSymbols is a free data retrieval call binding the contract method 0x5281b5c6.
//
// Solidity: function newSymbols(uint256 ) view returns(string)
func (_Oracle *OracleSession) NewSymbols(arg0 *big.Int) (string, error) {
	return _Oracle.Contract.NewSymbols(&_Oracle.CallOpts, arg0)
}

// NewSymbols is a free data retrieval call binding the contract method 0x5281b5c6.
//
// Solidity: function newSymbols(uint256 ) view returns(string)
func (_Oracle *OracleCallerSession) NewSymbols(arg0 *big.Int) (string, error) {
	return _Oracle.Contract.NewSymbols(&_Oracle.CallOpts, arg0)
}

// Reports is a free data retrieval call binding the contract method 0x4c56ea56.
//
// Solidity: function reports(string , address ) view returns(int256)
func (_Oracle *OracleCaller) Reports(opts *bind.CallOpts, arg0 string, arg1 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Oracle.contract.Call(opts, &out, "reports", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Reports is a free data retrieval call binding the contract method 0x4c56ea56.
//
// Solidity: function reports(string , address ) view returns(int256)
func (_Oracle *OracleSession) Reports(arg0 string, arg1 common.Address) (*big.Int, error) {
	return _Oracle.Contract.Reports(&_Oracle.CallOpts, arg0, arg1)
}

// Reports is a free data retrieval call binding the contract method 0x4c56ea56.
//
// Solidity: function reports(string , address ) view returns(int256)
func (_Oracle *OracleCallerSession) Reports(arg0 string, arg1 common.Address) (*big.Int, error) {
	return _Oracle.Contract.Reports(&_Oracle.CallOpts, arg0, arg1)
}

// Round is a free data retrieval call binding the contract method 0x146ca531.
//
// Solidity: function round() view returns(uint256)
func (_Oracle *OracleCaller) Round(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Oracle.contract.Call(opts, &out, "round")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Round is a free data retrieval call binding the contract method 0x146ca531.
//
// Solidity: function round() view returns(uint256)
func (_Oracle *OracleSession) Round() (*big.Int, error) {
	return _Oracle.Contract.Round(&_Oracle.CallOpts)
}

// Round is a free data retrieval call binding the contract method 0x146ca531.
//
// Solidity: function round() view returns(uint256)
func (_Oracle *OracleCallerSession) Round() (*big.Int, error) {
	return _Oracle.Contract.Round(&_Oracle.CallOpts)
}

// SymbolUpdatedRound is a free data retrieval call binding the contract method 0x08f21ff5.
//
// Solidity: function symbolUpdatedRound() view returns(int256)
func (_Oracle *OracleCaller) SymbolUpdatedRound(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Oracle.contract.Call(opts, &out, "symbolUpdatedRound")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// SymbolUpdatedRound is a free data retrieval call binding the contract method 0x08f21ff5.
//
// Solidity: function symbolUpdatedRound() view returns(int256)
func (_Oracle *OracleSession) SymbolUpdatedRound() (*big.Int, error) {
	return _Oracle.Contract.SymbolUpdatedRound(&_Oracle.CallOpts)
}

// SymbolUpdatedRound is a free data retrieval call binding the contract method 0x08f21ff5.
//
// Solidity: function symbolUpdatedRound() view returns(int256)
func (_Oracle *OracleCallerSession) SymbolUpdatedRound() (*big.Int, error) {
	return _Oracle.Contract.SymbolUpdatedRound(&_Oracle.CallOpts)
}

// Symbols is a free data retrieval call binding the contract method 0xccce413b.
//
// Solidity: function symbols(uint256 ) view returns(string)
func (_Oracle *OracleCaller) Symbols(opts *bind.CallOpts, arg0 *big.Int) (string, error) {
	var out []interface{}
	err := _Oracle.contract.Call(opts, &out, "symbols", arg0)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbols is a free data retrieval call binding the contract method 0xccce413b.
//
// Solidity: function symbols(uint256 ) view returns(string)
func (_Oracle *OracleSession) Symbols(arg0 *big.Int) (string, error) {
	return _Oracle.Contract.Symbols(&_Oracle.CallOpts, arg0)
}

// Symbols is a free data retrieval call binding the contract method 0xccce413b.
//
// Solidity: function symbols(uint256 ) view returns(string)
func (_Oracle *OracleCallerSession) Symbols(arg0 *big.Int) (string, error) {
	return _Oracle.Contract.Symbols(&_Oracle.CallOpts, arg0)
}

// VotePeriod is a free data retrieval call binding the contract method 0xa7813587.
//
// Solidity: function votePeriod() view returns(uint256)
func (_Oracle *OracleCaller) VotePeriod(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Oracle.contract.Call(opts, &out, "votePeriod")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// VotePeriod is a free data retrieval call binding the contract method 0xa7813587.
//
// Solidity: function votePeriod() view returns(uint256)
func (_Oracle *OracleSession) VotePeriod() (*big.Int, error) {
	return _Oracle.Contract.VotePeriod(&_Oracle.CallOpts)
}

// VotePeriod is a free data retrieval call binding the contract method 0xa7813587.
//
// Solidity: function votePeriod() view returns(uint256)
func (_Oracle *OracleCallerSession) VotePeriod() (*big.Int, error) {
	return _Oracle.Contract.VotePeriod(&_Oracle.CallOpts)
}

// VotingInfo is a free data retrieval call binding the contract method 0x5412b3ae.
//
// Solidity: function votingInfo(address ) view returns(uint256 round, uint256 commit, bool isVoter)
func (_Oracle *OracleCaller) VotingInfo(opts *bind.CallOpts, arg0 common.Address) (struct {
	Round   *big.Int
	Commit  *big.Int
	IsVoter bool
}, error) {
	var out []interface{}
	err := _Oracle.contract.Call(opts, &out, "votingInfo", arg0)

	outstruct := new(struct {
		Round   *big.Int
		Commit  *big.Int
		IsVoter bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Round = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Commit = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.IsVoter = *abi.ConvertType(out[2], new(bool)).(*bool)

	return *outstruct, err

}

// VotingInfo is a free data retrieval call binding the contract method 0x5412b3ae.
//
// Solidity: function votingInfo(address ) view returns(uint256 round, uint256 commit, bool isVoter)
func (_Oracle *OracleSession) VotingInfo(arg0 common.Address) (struct {
	Round   *big.Int
	Commit  *big.Int
	IsVoter bool
}, error) {
	return _Oracle.Contract.VotingInfo(&_Oracle.CallOpts, arg0)
}

// VotingInfo is a free data retrieval call binding the contract method 0x5412b3ae.
//
// Solidity: function votingInfo(address ) view returns(uint256 round, uint256 commit, bool isVoter)
func (_Oracle *OracleCallerSession) VotingInfo(arg0 common.Address) (struct {
	Round   *big.Int
	Commit  *big.Int
	IsVoter bool
}, error) {
	return _Oracle.Contract.VotingInfo(&_Oracle.CallOpts, arg0)
}

// Finalize is a paid mutator transaction binding the contract method 0x4bb278f3.
//
// Solidity: function finalize() returns()
func (_Oracle *OracleTransactor) Finalize(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Oracle.contract.Transact(opts, "finalize")
}

// Finalize is a paid mutator transaction binding the contract method 0x4bb278f3.
//
// Solidity: function finalize() returns()
func (_Oracle *OracleSession) Finalize() (*types.Transaction, error) {
	return _Oracle.Contract.Finalize(&_Oracle.TransactOpts)
}

// Finalize is a paid mutator transaction binding the contract method 0x4bb278f3.
//
// Solidity: function finalize() returns()
func (_Oracle *OracleTransactorSession) Finalize() (*types.Transaction, error) {
	return _Oracle.Contract.Finalize(&_Oracle.TransactOpts)
}

// SetOperator is a paid mutator transaction binding the contract method 0xb3ab15fb.
//
// Solidity: function setOperator(address _operator) returns()
func (_Oracle *OracleTransactor) SetOperator(opts *bind.TransactOpts, _operator common.Address) (*types.Transaction, error) {
	return _Oracle.contract.Transact(opts, "setOperator", _operator)
}

// SetOperator is a paid mutator transaction binding the contract method 0xb3ab15fb.
//
// Solidity: function setOperator(address _operator) returns()
func (_Oracle *OracleSession) SetOperator(_operator common.Address) (*types.Transaction, error) {
	return _Oracle.Contract.SetOperator(&_Oracle.TransactOpts, _operator)
}

// SetOperator is a paid mutator transaction binding the contract method 0xb3ab15fb.
//
// Solidity: function setOperator(address _operator) returns()
func (_Oracle *OracleTransactorSession) SetOperator(_operator common.Address) (*types.Transaction, error) {
	return _Oracle.Contract.SetOperator(&_Oracle.TransactOpts, _operator)
}

// SetSymbols is a paid mutator transaction binding the contract method 0x8d4f75d2.
//
// Solidity: function setSymbols(string[] _symbols) returns()
func (_Oracle *OracleTransactor) SetSymbols(opts *bind.TransactOpts, _symbols []string) (*types.Transaction, error) {
	return _Oracle.contract.Transact(opts, "setSymbols", _symbols)
}

// SetSymbols is a paid mutator transaction binding the contract method 0x8d4f75d2.
//
// Solidity: function setSymbols(string[] _symbols) returns()
func (_Oracle *OracleSession) SetSymbols(_symbols []string) (*types.Transaction, error) {
	return _Oracle.Contract.SetSymbols(&_Oracle.TransactOpts, _symbols)
}

// SetSymbols is a paid mutator transaction binding the contract method 0x8d4f75d2.
//
// Solidity: function setSymbols(string[] _symbols) returns()
func (_Oracle *OracleTransactorSession) SetSymbols(_symbols []string) (*types.Transaction, error) {
	return _Oracle.Contract.SetSymbols(&_Oracle.TransactOpts, _symbols)
}

// SetVoters is a paid mutator transaction binding the contract method 0x845023f2.
//
// Solidity: function setVoters(address[] _newVoters) returns()
func (_Oracle *OracleTransactor) SetVoters(opts *bind.TransactOpts, _newVoters []common.Address) (*types.Transaction, error) {
	return _Oracle.contract.Transact(opts, "setVoters", _newVoters)
}

// SetVoters is a paid mutator transaction binding the contract method 0x845023f2.
//
// Solidity: function setVoters(address[] _newVoters) returns()
func (_Oracle *OracleSession) SetVoters(_newVoters []common.Address) (*types.Transaction, error) {
	return _Oracle.Contract.SetVoters(&_Oracle.TransactOpts, _newVoters)
}

// SetVoters is a paid mutator transaction binding the contract method 0x845023f2.
//
// Solidity: function setVoters(address[] _newVoters) returns()
func (_Oracle *OracleTransactorSession) SetVoters(_newVoters []common.Address) (*types.Transaction, error) {
	return _Oracle.Contract.SetVoters(&_Oracle.TransactOpts, _newVoters)
}

// Vote is a paid mutator transaction binding the contract method 0x307de9b6.
//
// Solidity: function vote(uint256 _commit, int256[] _reports, uint256 _salt) returns()
func (_Oracle *OracleTransactor) Vote(opts *bind.TransactOpts, _commit *big.Int, _reports []*big.Int, _salt *big.Int) (*types.Transaction, error) {
	return _Oracle.contract.Transact(opts, "vote", _commit, _reports, _salt)
}

// Vote is a paid mutator transaction binding the contract method 0x307de9b6.
//
// Solidity: function vote(uint256 _commit, int256[] _reports, uint256 _salt) returns()
func (_Oracle *OracleSession) Vote(_commit *big.Int, _reports []*big.Int, _salt *big.Int) (*types.Transaction, error) {
	return _Oracle.Contract.Vote(&_Oracle.TransactOpts, _commit, _reports, _salt)
}

// Vote is a paid mutator transaction binding the contract method 0x307de9b6.
//
// Solidity: function vote(uint256 _commit, int256[] _reports, uint256 _salt) returns()
func (_Oracle *OracleTransactorSession) Vote(_commit *big.Int, _reports []*big.Int, _salt *big.Int) (*types.Transaction, error) {
	return _Oracle.Contract.Vote(&_Oracle.TransactOpts, _commit, _reports, _salt)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Oracle *OracleTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _Oracle.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Oracle *OracleSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _Oracle.Contract.Fallback(&_Oracle.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Oracle *OracleTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _Oracle.Contract.Fallback(&_Oracle.TransactOpts, calldata)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Oracle *OracleTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Oracle.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Oracle *OracleSession) Receive() (*types.Transaction, error) {
	return _Oracle.Contract.Receive(&_Oracle.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Oracle *OracleTransactorSession) Receive() (*types.Transaction, error) {
	return _Oracle.Contract.Receive(&_Oracle.TransactOpts)
}

// OracleNewRoundIterator is returned from FilterNewRound and is used to iterate over the raw logs and unpacked data for NewRound events raised by the Oracle contract.
type OracleNewRoundIterator struct {
	Event *OracleNewRound // Event containing the contract specifics and raw log

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
func (it *OracleNewRoundIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OracleNewRound)
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
		it.Event = new(OracleNewRound)
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
func (it *OracleNewRoundIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OracleNewRoundIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OracleNewRound represents a NewRound event raised by the Oracle contract.
type OracleNewRound struct {
	Round      *big.Int
	Height     *big.Int
	Timestamp  *big.Int
	VotePeriod *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterNewRound is a free log retrieval operation binding the contract event 0xb5d8636ab45e6cac7a4a61cb7c77f77f61a454d73aa2e6139ff8dcaf463537e5.
//
// Solidity: event NewRound(uint256 _round, uint256 _height, uint256 _timestamp, uint256 _votePeriod)
func (_Oracle *OracleFilterer) FilterNewRound(opts *bind.FilterOpts) (*OracleNewRoundIterator, error) {

	logs, sub, err := _Oracle.contract.FilterLogs(opts, "NewRound")
	if err != nil {
		return nil, err
	}
	return &OracleNewRoundIterator{contract: _Oracle.contract, event: "NewRound", logs: logs, sub: sub}, nil
}

// WatchNewRound is a free log subscription operation binding the contract event 0xb5d8636ab45e6cac7a4a61cb7c77f77f61a454d73aa2e6139ff8dcaf463537e5.
//
// Solidity: event NewRound(uint256 _round, uint256 _height, uint256 _timestamp, uint256 _votePeriod)
func (_Oracle *OracleFilterer) WatchNewRound(opts *bind.WatchOpts, sink chan<- *OracleNewRound) (event.Subscription, error) {

	logs, sub, err := _Oracle.contract.WatchLogs(opts, "NewRound")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OracleNewRound)
				if err := _Oracle.contract.UnpackLog(event, "NewRound", log); err != nil {
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

// ParseNewRound is a log parse operation binding the contract event 0xb5d8636ab45e6cac7a4a61cb7c77f77f61a454d73aa2e6139ff8dcaf463537e5.
//
// Solidity: event NewRound(uint256 _round, uint256 _height, uint256 _timestamp, uint256 _votePeriod)
func (_Oracle *OracleFilterer) ParseNewRound(log types.Log) (*OracleNewRound, error) {
	event := new(OracleNewRound)
	if err := _Oracle.contract.UnpackLog(event, "NewRound", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OracleNewSymbolsIterator is returned from FilterNewSymbols and is used to iterate over the raw logs and unpacked data for NewSymbols events raised by the Oracle contract.
type OracleNewSymbolsIterator struct {
	Event *OracleNewSymbols // Event containing the contract specifics and raw log

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
func (it *OracleNewSymbolsIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OracleNewSymbols)
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
		it.Event = new(OracleNewSymbols)
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
func (it *OracleNewSymbolsIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OracleNewSymbolsIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OracleNewSymbols represents a NewSymbols event raised by the Oracle contract.
type OracleNewSymbols struct {
	Symbols []string
	Round   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterNewSymbols is a free log retrieval operation binding the contract event 0xaa278e424da680ce5dad66510415760e78e0bd87d45c786c6e88bdde82f9342d.
//
// Solidity: event NewSymbols(string[] _symbols, uint256 _round)
func (_Oracle *OracleFilterer) FilterNewSymbols(opts *bind.FilterOpts) (*OracleNewSymbolsIterator, error) {

	logs, sub, err := _Oracle.contract.FilterLogs(opts, "NewSymbols")
	if err != nil {
		return nil, err
	}
	return &OracleNewSymbolsIterator{contract: _Oracle.contract, event: "NewSymbols", logs: logs, sub: sub}, nil
}

// WatchNewSymbols is a free log subscription operation binding the contract event 0xaa278e424da680ce5dad66510415760e78e0bd87d45c786c6e88bdde82f9342d.
//
// Solidity: event NewSymbols(string[] _symbols, uint256 _round)
func (_Oracle *OracleFilterer) WatchNewSymbols(opts *bind.WatchOpts, sink chan<- *OracleNewSymbols) (event.Subscription, error) {

	logs, sub, err := _Oracle.contract.WatchLogs(opts, "NewSymbols")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OracleNewSymbols)
				if err := _Oracle.contract.UnpackLog(event, "NewSymbols", log); err != nil {
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
func (_Oracle *OracleFilterer) ParseNewSymbols(log types.Log) (*OracleNewSymbols, error) {
	event := new(OracleNewSymbols)
	if err := _Oracle.contract.UnpackLog(event, "NewSymbols", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OracleVotedIterator is returned from FilterVoted and is used to iterate over the raw logs and unpacked data for Voted events raised by the Oracle contract.
type OracleVotedIterator struct {
	Event *OracleVoted // Event containing the contract specifics and raw log

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
func (it *OracleVotedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OracleVoted)
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
		it.Event = new(OracleVoted)
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
func (it *OracleVotedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OracleVotedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OracleVoted represents a Voted event raised by the Oracle contract.
type OracleVoted struct {
	Voter common.Address
	Votes []*big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterVoted is a free log retrieval operation binding the contract event 0xd0d8560f1076ac6b216b1091a2571d6f9bc3e0889f4dbdbe1c7d1be7136714d3.
//
// Solidity: event Voted(address indexed _voter, int256[] _votes)
func (_Oracle *OracleFilterer) FilterVoted(opts *bind.FilterOpts, _voter []common.Address) (*OracleVotedIterator, error) {

	var _voterRule []interface{}
	for _, _voterItem := range _voter {
		_voterRule = append(_voterRule, _voterItem)
	}

	logs, sub, err := _Oracle.contract.FilterLogs(opts, "Voted", _voterRule)
	if err != nil {
		return nil, err
	}
	return &OracleVotedIterator{contract: _Oracle.contract, event: "Voted", logs: logs, sub: sub}, nil
}

// WatchVoted is a free log subscription operation binding the contract event 0xd0d8560f1076ac6b216b1091a2571d6f9bc3e0889f4dbdbe1c7d1be7136714d3.
//
// Solidity: event Voted(address indexed _voter, int256[] _votes)
func (_Oracle *OracleFilterer) WatchVoted(opts *bind.WatchOpts, sink chan<- *OracleVoted, _voter []common.Address) (event.Subscription, error) {

	var _voterRule []interface{}
	for _, _voterItem := range _voter {
		_voterRule = append(_voterRule, _voterItem)
	}

	logs, sub, err := _Oracle.contract.WatchLogs(opts, "Voted", _voterRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OracleVoted)
				if err := _Oracle.contract.UnpackLog(event, "Voted", log); err != nil {
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

// ParseVoted is a log parse operation binding the contract event 0xd0d8560f1076ac6b216b1091a2571d6f9bc3e0889f4dbdbe1c7d1be7136714d3.
//
// Solidity: event Voted(address indexed _voter, int256[] _votes)
func (_Oracle *OracleFilterer) ParseVoted(log types.Log) (*OracleVoted, error) {
	event := new(OracleVoted)
	if err := _Oracle.contract.UnpackLog(event, "Voted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PrecompiledMetaData contains all meta data concerning the Precompiled contract.
var PrecompiledMetaData = &bind.MetaData{
	ABI: "[]",
	Bin: "0x60566037600b82828239805160001a607314602a57634e487b7160e01b600052600060045260246000fd5b30600052607381538281f3fe73000000000000000000000000000000000000000030146080604052600080fdfea26469706673582212207ae1da95cda5388d9e69c43836037b675fc61f0ffd041262f569b69c4bd4d92d64736f6c634300080c0033",
}

// PrecompiledABI is the input ABI used to generate the binding from.
// Deprecated: Use PrecompiledMetaData.ABI instead.
var PrecompiledABI = PrecompiledMetaData.ABI

// PrecompiledBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use PrecompiledMetaData.Bin instead.
var PrecompiledBin = PrecompiledMetaData.Bin

// DeployPrecompiled deploys a new Ethereum contract, binding an instance of Precompiled to it.
func DeployPrecompiled(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Precompiled, error) {
	parsed, err := PrecompiledMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(PrecompiledBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Precompiled{PrecompiledCaller: PrecompiledCaller{contract: contract}, PrecompiledTransactor: PrecompiledTransactor{contract: contract}, PrecompiledFilterer: PrecompiledFilterer{contract: contract}}, nil
}

// Precompiled is an auto generated Go binding around an Ethereum contract.
type Precompiled struct {
	PrecompiledCaller     // Read-only binding to the contract
	PrecompiledTransactor // Write-only binding to the contract
	PrecompiledFilterer   // Log filterer for contract events
}

// PrecompiledCaller is an auto generated read-only Go binding around an Ethereum contract.
type PrecompiledCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PrecompiledTransactor is an auto generated write-only Go binding around an Ethereum contract.
type PrecompiledTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PrecompiledFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type PrecompiledFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PrecompiledSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type PrecompiledSession struct {
	Contract     *Precompiled      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// PrecompiledCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type PrecompiledCallerSession struct {
	Contract *PrecompiledCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// PrecompiledTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type PrecompiledTransactorSession struct {
	Contract     *PrecompiledTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// PrecompiledRaw is an auto generated low-level Go binding around an Ethereum contract.
type PrecompiledRaw struct {
	Contract *Precompiled // Generic contract binding to access the raw methods on
}

// PrecompiledCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type PrecompiledCallerRaw struct {
	Contract *PrecompiledCaller // Generic read-only contract binding to access the raw methods on
}

// PrecompiledTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type PrecompiledTransactorRaw struct {
	Contract *PrecompiledTransactor // Generic write-only contract binding to access the raw methods on
}

// NewPrecompiled creates a new instance of Precompiled, bound to a specific deployed contract.
func NewPrecompiled(address common.Address, backend bind.ContractBackend) (*Precompiled, error) {
	contract, err := bindPrecompiled(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Precompiled{PrecompiledCaller: PrecompiledCaller{contract: contract}, PrecompiledTransactor: PrecompiledTransactor{contract: contract}, PrecompiledFilterer: PrecompiledFilterer{contract: contract}}, nil
}

// NewPrecompiledCaller creates a new read-only instance of Precompiled, bound to a specific deployed contract.
func NewPrecompiledCaller(address common.Address, caller bind.ContractCaller) (*PrecompiledCaller, error) {
	contract, err := bindPrecompiled(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &PrecompiledCaller{contract: contract}, nil
}

// NewPrecompiledTransactor creates a new write-only instance of Precompiled, bound to a specific deployed contract.
func NewPrecompiledTransactor(address common.Address, transactor bind.ContractTransactor) (*PrecompiledTransactor, error) {
	contract, err := bindPrecompiled(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &PrecompiledTransactor{contract: contract}, nil
}

// NewPrecompiledFilterer creates a new log filterer instance of Precompiled, bound to a specific deployed contract.
func NewPrecompiledFilterer(address common.Address, filterer bind.ContractFilterer) (*PrecompiledFilterer, error) {
	contract, err := bindPrecompiled(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &PrecompiledFilterer{contract: contract}, nil
}

// bindPrecompiled binds a generic wrapper to an already deployed contract.
func bindPrecompiled(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(PrecompiledABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Precompiled *PrecompiledRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Precompiled.Contract.PrecompiledCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Precompiled *PrecompiledRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Precompiled.Contract.PrecompiledTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Precompiled *PrecompiledRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Precompiled.Contract.PrecompiledTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Precompiled *PrecompiledCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Precompiled.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Precompiled *PrecompiledTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Precompiled.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Precompiled *PrecompiledTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Precompiled.Contract.contract.Transact(opts, method, params...)
}

// SafeMathMetaData contains all meta data concerning the SafeMath contract.
var SafeMathMetaData = &bind.MetaData{
	ABI: "[]",
	Bin: "0x60566037600b82828239805160001a607314602a57634e487b7160e01b600052600060045260246000fd5b30600052607381538281f3fe73000000000000000000000000000000000000000030146080604052600080fdfea2646970667358221220c86b0948850cfc3789c47065bffd40b3df54e6cc3c8e16585962908f821174c264736f6c634300080c0033",
}

// SafeMathABI is the input ABI used to generate the binding from.
// Deprecated: Use SafeMathMetaData.ABI instead.
var SafeMathABI = SafeMathMetaData.ABI

// SafeMathBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use SafeMathMetaData.Bin instead.
var SafeMathBin = SafeMathMetaData.Bin

// DeploySafeMath deploys a new Ethereum contract, binding an instance of SafeMath to it.
func DeploySafeMath(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *SafeMath, error) {
	parsed, err := SafeMathMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(SafeMathBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &SafeMath{SafeMathCaller: SafeMathCaller{contract: contract}, SafeMathTransactor: SafeMathTransactor{contract: contract}, SafeMathFilterer: SafeMathFilterer{contract: contract}}, nil
}

// SafeMath is an auto generated Go binding around an Ethereum contract.
type SafeMath struct {
	SafeMathCaller     // Read-only binding to the contract
	SafeMathTransactor // Write-only binding to the contract
	SafeMathFilterer   // Log filterer for contract events
}

// SafeMathCaller is an auto generated read-only Go binding around an Ethereum contract.
type SafeMathCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SafeMathTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SafeMathTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SafeMathFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SafeMathFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SafeMathSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SafeMathSession struct {
	Contract     *SafeMath         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SafeMathCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SafeMathCallerSession struct {
	Contract *SafeMathCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// SafeMathTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SafeMathTransactorSession struct {
	Contract     *SafeMathTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// SafeMathRaw is an auto generated low-level Go binding around an Ethereum contract.
type SafeMathRaw struct {
	Contract *SafeMath // Generic contract binding to access the raw methods on
}

// SafeMathCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SafeMathCallerRaw struct {
	Contract *SafeMathCaller // Generic read-only contract binding to access the raw methods on
}

// SafeMathTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SafeMathTransactorRaw struct {
	Contract *SafeMathTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSafeMath creates a new instance of SafeMath, bound to a specific deployed contract.
func NewSafeMath(address common.Address, backend bind.ContractBackend) (*SafeMath, error) {
	contract, err := bindSafeMath(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SafeMath{SafeMathCaller: SafeMathCaller{contract: contract}, SafeMathTransactor: SafeMathTransactor{contract: contract}, SafeMathFilterer: SafeMathFilterer{contract: contract}}, nil
}

// NewSafeMathCaller creates a new read-only instance of SafeMath, bound to a specific deployed contract.
func NewSafeMathCaller(address common.Address, caller bind.ContractCaller) (*SafeMathCaller, error) {
	contract, err := bindSafeMath(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SafeMathCaller{contract: contract}, nil
}

// NewSafeMathTransactor creates a new write-only instance of SafeMath, bound to a specific deployed contract.
func NewSafeMathTransactor(address common.Address, transactor bind.ContractTransactor) (*SafeMathTransactor, error) {
	contract, err := bindSafeMath(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SafeMathTransactor{contract: contract}, nil
}

// NewSafeMathFilterer creates a new log filterer instance of SafeMath, bound to a specific deployed contract.
func NewSafeMathFilterer(address common.Address, filterer bind.ContractFilterer) (*SafeMathFilterer, error) {
	contract, err := bindSafeMath(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SafeMathFilterer{contract: contract}, nil
}

// bindSafeMath binds a generic wrapper to an already deployed contract.
func bindSafeMath(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(SafeMathABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SafeMath *SafeMathRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SafeMath.Contract.SafeMathCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SafeMath *SafeMathRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SafeMath.Contract.SafeMathTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SafeMath *SafeMathRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SafeMath.Contract.SafeMathTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SafeMath *SafeMathCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SafeMath.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SafeMath *SafeMathTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SafeMath.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SafeMath *SafeMathTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SafeMath.Contract.contract.Transact(opts, method, params...)
}

// UpgradeableMetaData contains all meta data concerning the Upgradeable contract.
var UpgradeableMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"completeContractUpgrade\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getNewContract\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"},{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"resetContractUpgrade\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"_bytecode\",\"type\":\"bytes\"},{\"internalType\":\"string\",\"name\":\"_abi\",\"type\":\"string\"}],\"name\":\"upgradeContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"872cf059": "completeContractUpgrade()",
		"b66b3e79": "getNewContract()",
		"cf9c5719": "resetContractUpgrade()",
		"b2ea9adb": "upgradeContract(bytes,string)",
	},
}

// UpgradeableABI is the input ABI used to generate the binding from.
// Deprecated: Use UpgradeableMetaData.ABI instead.
var UpgradeableABI = UpgradeableMetaData.ABI

// Deprecated: Use UpgradeableMetaData.Sigs instead.
// UpgradeableFuncSigs maps the 4-byte function signature to its string representation.
var UpgradeableFuncSigs = UpgradeableMetaData.Sigs

// Upgradeable is an auto generated Go binding around an Ethereum contract.
type Upgradeable struct {
	UpgradeableCaller     // Read-only binding to the contract
	UpgradeableTransactor // Write-only binding to the contract
	UpgradeableFilterer   // Log filterer for contract events
}

// UpgradeableCaller is an auto generated read-only Go binding around an Ethereum contract.
type UpgradeableCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UpgradeableTransactor is an auto generated write-only Go binding around an Ethereum contract.
type UpgradeableTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UpgradeableFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type UpgradeableFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UpgradeableSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type UpgradeableSession struct {
	Contract     *Upgradeable      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// UpgradeableCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type UpgradeableCallerSession struct {
	Contract *UpgradeableCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// UpgradeableTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type UpgradeableTransactorSession struct {
	Contract     *UpgradeableTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// UpgradeableRaw is an auto generated low-level Go binding around an Ethereum contract.
type UpgradeableRaw struct {
	Contract *Upgradeable // Generic contract binding to access the raw methods on
}

// UpgradeableCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type UpgradeableCallerRaw struct {
	Contract *UpgradeableCaller // Generic read-only contract binding to access the raw methods on
}

// UpgradeableTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type UpgradeableTransactorRaw struct {
	Contract *UpgradeableTransactor // Generic write-only contract binding to access the raw methods on
}

// NewUpgradeable creates a new instance of Upgradeable, bound to a specific deployed contract.
func NewUpgradeable(address common.Address, backend bind.ContractBackend) (*Upgradeable, error) {
	contract, err := bindUpgradeable(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Upgradeable{UpgradeableCaller: UpgradeableCaller{contract: contract}, UpgradeableTransactor: UpgradeableTransactor{contract: contract}, UpgradeableFilterer: UpgradeableFilterer{contract: contract}}, nil
}

// NewUpgradeableCaller creates a new read-only instance of Upgradeable, bound to a specific deployed contract.
func NewUpgradeableCaller(address common.Address, caller bind.ContractCaller) (*UpgradeableCaller, error) {
	contract, err := bindUpgradeable(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &UpgradeableCaller{contract: contract}, nil
}

// NewUpgradeableTransactor creates a new write-only instance of Upgradeable, bound to a specific deployed contract.
func NewUpgradeableTransactor(address common.Address, transactor bind.ContractTransactor) (*UpgradeableTransactor, error) {
	contract, err := bindUpgradeable(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &UpgradeableTransactor{contract: contract}, nil
}

// NewUpgradeableFilterer creates a new log filterer instance of Upgradeable, bound to a specific deployed contract.
func NewUpgradeableFilterer(address common.Address, filterer bind.ContractFilterer) (*UpgradeableFilterer, error) {
	contract, err := bindUpgradeable(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &UpgradeableFilterer{contract: contract}, nil
}

// bindUpgradeable binds a generic wrapper to an already deployed contract.
func bindUpgradeable(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(UpgradeableABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Upgradeable *UpgradeableRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Upgradeable.Contract.UpgradeableCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Upgradeable *UpgradeableRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Upgradeable.Contract.UpgradeableTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Upgradeable *UpgradeableRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Upgradeable.Contract.UpgradeableTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Upgradeable *UpgradeableCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Upgradeable.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Upgradeable *UpgradeableTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Upgradeable.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Upgradeable *UpgradeableTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Upgradeable.Contract.contract.Transact(opts, method, params...)
}

// GetNewContract is a free data retrieval call binding the contract method 0xb66b3e79.
//
// Solidity: function getNewContract() view returns(bytes, string)
func (_Upgradeable *UpgradeableCaller) GetNewContract(opts *bind.CallOpts) ([]byte, string, error) {
	var out []interface{}
	err := _Upgradeable.contract.Call(opts, &out, "getNewContract")

	if err != nil {
		return *new([]byte), *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)
	out1 := *abi.ConvertType(out[1], new(string)).(*string)

	return out0, out1, err

}

// GetNewContract is a free data retrieval call binding the contract method 0xb66b3e79.
//
// Solidity: function getNewContract() view returns(bytes, string)
func (_Upgradeable *UpgradeableSession) GetNewContract() ([]byte, string, error) {
	return _Upgradeable.Contract.GetNewContract(&_Upgradeable.CallOpts)
}

// GetNewContract is a free data retrieval call binding the contract method 0xb66b3e79.
//
// Solidity: function getNewContract() view returns(bytes, string)
func (_Upgradeable *UpgradeableCallerSession) GetNewContract() ([]byte, string, error) {
	return _Upgradeable.Contract.GetNewContract(&_Upgradeable.CallOpts)
}

// CompleteContractUpgrade is a paid mutator transaction binding the contract method 0x872cf059.
//
// Solidity: function completeContractUpgrade() returns()
func (_Upgradeable *UpgradeableTransactor) CompleteContractUpgrade(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Upgradeable.contract.Transact(opts, "completeContractUpgrade")
}

// CompleteContractUpgrade is a paid mutator transaction binding the contract method 0x872cf059.
//
// Solidity: function completeContractUpgrade() returns()
func (_Upgradeable *UpgradeableSession) CompleteContractUpgrade() (*types.Transaction, error) {
	return _Upgradeable.Contract.CompleteContractUpgrade(&_Upgradeable.TransactOpts)
}

// CompleteContractUpgrade is a paid mutator transaction binding the contract method 0x872cf059.
//
// Solidity: function completeContractUpgrade() returns()
func (_Upgradeable *UpgradeableTransactorSession) CompleteContractUpgrade() (*types.Transaction, error) {
	return _Upgradeable.Contract.CompleteContractUpgrade(&_Upgradeable.TransactOpts)
}

// ResetContractUpgrade is a paid mutator transaction binding the contract method 0xcf9c5719.
//
// Solidity: function resetContractUpgrade() returns()
func (_Upgradeable *UpgradeableTransactor) ResetContractUpgrade(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Upgradeable.contract.Transact(opts, "resetContractUpgrade")
}

// ResetContractUpgrade is a paid mutator transaction binding the contract method 0xcf9c5719.
//
// Solidity: function resetContractUpgrade() returns()
func (_Upgradeable *UpgradeableSession) ResetContractUpgrade() (*types.Transaction, error) {
	return _Upgradeable.Contract.ResetContractUpgrade(&_Upgradeable.TransactOpts)
}

// ResetContractUpgrade is a paid mutator transaction binding the contract method 0xcf9c5719.
//
// Solidity: function resetContractUpgrade() returns()
func (_Upgradeable *UpgradeableTransactorSession) ResetContractUpgrade() (*types.Transaction, error) {
	return _Upgradeable.Contract.ResetContractUpgrade(&_Upgradeable.TransactOpts)
}

// UpgradeContract is a paid mutator transaction binding the contract method 0xb2ea9adb.
//
// Solidity: function upgradeContract(bytes _bytecode, string _abi) returns()
func (_Upgradeable *UpgradeableTransactor) UpgradeContract(opts *bind.TransactOpts, _bytecode []byte, _abi string) (*types.Transaction, error) {
	return _Upgradeable.contract.Transact(opts, "upgradeContract", _bytecode, _abi)
}

// UpgradeContract is a paid mutator transaction binding the contract method 0xb2ea9adb.
//
// Solidity: function upgradeContract(bytes _bytecode, string _abi) returns()
func (_Upgradeable *UpgradeableSession) UpgradeContract(_bytecode []byte, _abi string) (*types.Transaction, error) {
	return _Upgradeable.Contract.UpgradeContract(&_Upgradeable.TransactOpts, _bytecode, _abi)
}

// UpgradeContract is a paid mutator transaction binding the contract method 0xb2ea9adb.
//
// Solidity: function upgradeContract(bytes _bytecode, string _abi) returns()
func (_Upgradeable *UpgradeableTransactorSession) UpgradeContract(_bytecode []byte, _abi string) (*types.Transaction, error) {
	return _Upgradeable.Contract.UpgradeContract(&_Upgradeable.TransactOpts, _bytecode, _abi)
}
