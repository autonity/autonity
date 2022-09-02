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
	Addr              common.Address
	Enode             string
	CommissionRate    *big.Int
	BondedStake       *big.Int
	TotalSlashed      *big.Int
	LiquidContract    common.Address
	LiquidSupply      *big.Int
	RegistrationBlock *big.Int
	State             uint8
}

// AutonityMetaData contains all meta data concerning the Autonity contract.
var AutonityMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"components\":[{\"internalType\":\"addresspayable\",\"name\":\"treasury\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"enode\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"commissionRate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"bondedStake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalSlashed\",\"type\":\"uint256\"},{\"internalType\":\"contractLiquid\",\"name\":\"liquidContract\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"liquidSupply\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"registrationBlock\",\"type\":\"uint256\"},{\"internalType\":\"enumAutonity.ValidatorState\",\"name\":\"state\",\"type\":\"uint8\"}],\"internalType\":\"structAutonity.Validator[]\",\"name\":\"_validators\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"operatorAccount\",\"type\":\"address\"},{\"internalType\":\"addresspayable\",\"name\":\"treasuryAccount\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"treasuryFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"minBaseFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"delegationRate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"epochPeriod\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"unbondingPeriod\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"committeeSize\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"contractVersion\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"blockPeriod\",\"type\":\"uint256\"}],\"internalType\":\"structAutonity.Config\",\"name\":\"_config\",\"type\":\"tuple\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"BurnedStake\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasPrice\",\"type\":\"uint256\"}],\"name\":\"MinimumBaseFeeUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"MintedStake\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"treasury\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"effectiveBlock\",\"type\":\"uint256\"}],\"name\":\"PausedValidator\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"treasury\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"enode\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"liquidContract\",\"type\":\"address\"}],\"name\":\"RegisteredValidator\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Rewarded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"inputs\":[],\"name\":\"COMMISSION_RATE_PRECISION\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_address\",\"type\":\"address\"}],\"name\":\"activateValidator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_validator\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"bond\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"completeContractUpgrade\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"computeCommittee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"config\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"operatorAccount\",\"type\":\"address\"},{\"internalType\":\"addresspayable\",\"name\":\"treasuryAccount\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"treasuryFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"minBaseFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"delegationRate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"epochPeriod\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"unbondingPeriod\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"committeeSize\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"contractVersion\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"blockPeriod\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"deployer\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"epochID\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"epochReward\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"epochTotalBondedStake\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"finalize\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"votingPower\",\"type\":\"uint256\"}],\"internalType\":\"structAutonity.CommitteeMember[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"startId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"lastId\",\"type\":\"uint256\"}],\"name\":\"getBondingReq\",\"outputs\":[{\"components\":[{\"internalType\":\"addresspayable\",\"name\":\"delegator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"delegatee\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"startBlock\",\"type\":\"uint256\"}],\"internalType\":\"structAutonity.Staking[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCommittee\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"votingPower\",\"type\":\"uint256\"}],\"internalType\":\"structAutonity.CommitteeMember[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCommitteeEnodes\",\"outputs\":[{\"internalType\":\"string[]\",\"name\":\"\",\"type\":\"string[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLastEpochBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getMaxCommitteeSize\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getMinimumBaseFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getNewContract\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"},{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getOperator\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"height\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"round\",\"type\":\"uint256\"}],\"name\":\"getProposer\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"startId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"lastId\",\"type\":\"uint256\"}],\"name\":\"getUnbondingReq\",\"outputs\":[{\"components\":[{\"internalType\":\"addresspayable\",\"name\":\"delegator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"delegatee\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"startBlock\",\"type\":\"uint256\"}],\"internalType\":\"structAutonity.Staking[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"}],\"name\":\"getValidator\",\"outputs\":[{\"components\":[{\"internalType\":\"addresspayable\",\"name\":\"treasury\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"enode\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"commissionRate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"bondedStake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalSlashed\",\"type\":\"uint256\"},{\"internalType\":\"contractLiquid\",\"name\":\"liquidContract\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"liquidSupply\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"registrationBlock\",\"type\":\"uint256\"},{\"internalType\":\"enumAutonity.ValidatorState\",\"name\":\"state\",\"type\":\"uint8\"}],\"internalType\":\"structAutonity.Validator\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getValidators\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getVersion\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"headBondingID\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"headUnbondingID\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastEpochBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"mint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_address\",\"type\":\"address\"}],\"name\":\"pauseValidator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_enode\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"_proof\",\"type\":\"bytes\"}],\"name\":\"registerValidator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"resetContractUpgrade\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_size\",\"type\":\"uint256\"}],\"name\":\"setCommitteeSize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_period\",\"type\":\"uint256\"}],\"name\":\"setEpochPeriod\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_price\",\"type\":\"uint256\"}],\"name\":\"setMinimumBaseFee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_account\",\"type\":\"address\"}],\"name\":\"setOperatorAccount\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"_account\",\"type\":\"address\"}],\"name\":\"setTreasuryAccount\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_treasuryFee\",\"type\":\"uint256\"}],\"name\":\"setTreasuryFee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_period\",\"type\":\"uint256\"}],\"name\":\"setUnbondingPeriod\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"tailBondingID\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"tailUnbondingID\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalRedistributed\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_validator\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"unbond\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"_bytecode\",\"type\":\"bytes\"},{\"internalType\":\"string\",\"name\":\"_abi\",\"type\":\"string\"}],\"name\":\"upgradeContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Sigs: map[string]string{
		"2f2c3f2e": "COMMISSION_RATE_PRECISION()",
		"b46e5520": "activateValidator(address)",
		"dd62ed3e": "allowance(address,address)",
		"095ea7b3": "approve(address,uint256)",
		"70a08231": "balanceOf(address)",
		"a515366a": "bond(address,uint256)",
		"9dc29fac": "burn(address,uint256)",
		"872cf059": "completeContractUpgrade()",
		"ae1f5fa0": "computeCommittee()",
		"79502c55": "config()",
		"d5f39488": "deployer()",
		"c9d97af4": "epochID()",
		"1604e416": "epochReward()",
		"9c98e471": "epochTotalBondedStake()",
		"05261aea": "finalize(uint256)",
		"e485c6fb": "getBondingReq(uint256,uint256)",
		"ab8f6ffe": "getCommittee()",
		"a8b2216e": "getCommitteeEnodes()",
		"731b3a03": "getLastEpochBlock()",
		"819b6463": "getMaxCommitteeSize()",
		"11220633": "getMinimumBaseFee()",
		"b66b3e79": "getNewContract()",
		"e7f43c68": "getOperator()",
		"5f7d3949": "getProposer(uint256,uint256)",
		"55230e93": "getUnbondingReq(uint256,uint256)",
		"1904bb2e": "getValidator(address)",
		"b7ab4db5": "getValidators()",
		"0d8e6e2c": "getVersion()",
		"44697221": "headBondingID()",
		"4b0dff63": "headUnbondingID()",
		"c2362dd5": "lastEpochBlock()",
		"40c10f19": "mint(address,uint256)",
		"06fdde03": "name()",
		"0ae65e7a": "pauseValidator(address)",
		"5333d404": "registerValidator(string,bytes)",
		"cf9c5719": "resetContractUpgrade()",
		"8bac7dad": "setCommitteeSize(uint256)",
		"6b5f444c": "setEpochPeriod(uint256)",
		"cb696f54": "setMinimumBaseFee(uint256)",
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
	Bin: "0x60806040523480156200001157600080fd5b5060405162007bf138038062007bf1833981016040819052620000349162001829565b600b546200005a57601f80546001600160a01b031916331790556200005a828262000062565b505062001bec565b8051600380546001600160a01b039283166001600160a01b03199182161790915560208301516004805491909316911617905560408101516005556060810151600655608081015160075560a081015160085560c081015160095560e0810151600a55610100810151600b55610120810151600c5560005b8251811015620003ce576000838281518110620000fb57620000fb620019ce565b60200260200101516080015190506000848381518110620001205762000120620019ce565b602002602001015160e00181815250506000848381518110620001475762000147620019ce565b602002602001015160c001906001600160a01b031690816001600160a01b0316815250506000848381518110620001825762000182620019ce565b602002602001015160800181815250506000848381518110620001a957620001a9620019ce565b6020026020010151610100018181525050600360040154848381518110620001d557620001d5620019ce565b602002602001015160600181815250506000848381518110620001fc57620001fc620019ce565b6020026020010151610120019060018111156200021d576200021d620019e4565b90816001811115620002335762000233620019e4565b9052506710000000000000008110620002a45760405162461bcd60e51b815260206004820152602860248201527f697373756564204e6577746f6e2063616e277420626520677265617465722074604482015267068616e20325e36360c41b60648201526084015b60405180910390fd5b620002d1848381518110620002bd57620002bd620019ce565b6020026020010151620003e760201b60201c565b80601c6000868581518110620002eb57620002eb620019ce565b6020026020010151600001516001600160a01b03166001600160a01b03168152602001908152602001600020600082825462000328919062001a10565b9250508190555080601e600082825462000343919062001a10565b9250508190555080601060008282546200035e919062001a10565b92505081905550620003b88483815181106200037e576200037e620019ce565b60200260200101516020015182868581518110620003a057620003a0620019ce565b6020026020010151600001516200040060201b60201c565b5080620003c58162001a2b565b915050620000da565b50620003d96200058c565b620003e36200063b565b5050565b620003f28162000d6f565b620003fd8162000ea7565b50565b600082116200045e5760405162461bcd60e51b815260206004820152602360248201527f616d6f756e74206e65656420746f206265207374726963746c7920706f73697460448201526269766560e81b60648201526084016200029b565b6001600160a01b0381166000908152601c6020526040902054821115620004c85760405162461bcd60e51b815260206004820152601b60248201527f696e73756666696369656e74204e6577746f6e2062616c616e6365000000000060448201526064016200029b565b6001600160a01b0381166000908152601c602052604081208054849290620004f290849062001a49565b9091555050604080516080810182526001600160a01b038084168252858116602080840191825283850187815243606086019081526018805460009081526016909452968320865181549087166001600160a01b031991821617825594516001820180549190971695169490941790945551600283015591516003909101558254919290620005818362001a2b565b919050555050505050565b6017545b601854811015620005bb57620005a68162001049565b80620005b28162001a2b565b91505062000590565b50601854601755601a54805b601b5481101562000635576009546000828152601960205260409020600301544391620005f49162001a10565b116200061a57620006058162001162565b6200061260018362001a10565b915062000620565b62000635565b806200062c8162001a2b565b915050620005c7565b50601a55565b601f546001600160a01b03163314620006a35760405162461bcd60e51b815260206004820152602360248201527f66756e6374696f6e207265737472696374656420746f207468652070726f746f60448201526218dbdb60ea1b60648201526084016200029b565b600d54620006f45760405162461bcd60e51b815260206004820152601860248201527f5468657265206d7573742062652076616c696461746f7273000000000000000060448201526064016200029b565b6000805b600d54811015620007d8576000601d6000600d84815481106200071f576200071f620019ce565b60009182526020808320909101546001600160a01b0316835282019290925260400190206009015460ff1660018111156200075e576200075e620019e4565b148015620007ad57506000601d6000600d8481548110620007835762000783620019ce565b60009182526020808320909101546001600160a01b03168352820192909252604001902060040154115b15620007c35781620007bf8162001a2b565b9250505b80620007cf8162001a2b565b915050620006f8565b50600a54818110620007e75750805b6000826001600160401b0381111562000804576200080462001658565b6040519080825280602002602001820160405280156200084157816020015b6200082d6200143a565b815260200190600190039081620008235790505b5090506000826001600160401b0381111562000861576200086162001658565b6040519080825280602002602001820160405280156200089e57816020015b6200088a6200143a565b815260200190600190039081620008805790505b5090506000805b600d5481101562000b2a576000601d6000600d8481548110620008cc57620008cc620019ce565b60009182526020808320909101546001600160a01b0316835282019290925260400190206009015460ff1660018111156200090b576200090b620019e4565b1480156200095a57506000601d6000600d8481548110620009305762000930620019ce565b60009182526020808320909101546001600160a01b03168352820192909252604001902060040154115b1562000b15576000601d6000600d84815481106200097c576200097c620019ce565b60009182526020808320909101546001600160a01b0390811684528382019490945260409283019091208251610140810184528154851681526001820154909416918401919091526002810180549192840191620009da9062001a63565b80601f016020809104026020016040519081016040528092919081815260200182805462000a089062001a63565b801562000a595780601f1062000a2d5761010080835404028352916020019162000a59565b820191906000526020600020905b81548152906001019060200180831162000a3b57829003601f168201915b505050918352505060038201546020820152600482015460408201526005820154606082015260068201546001600160a01b03166080820152600782015460a0820152600882015460c0820152600982015460e09091019060ff16600181111562000ac85762000ac8620019e4565b600181111562000adc5762000adc620019e4565b8152505090508085848151811062000af85762000af8620019ce565b6020026020010181905250828062000b109062001a2b565b935050505b8062000b218162001a2b565b915050620008a5565b50600a548351111562000baa5762000b42836200121d565b60005b600a5481101562000ba35783818151811062000b655762000b65620019ce565b602002602001015183828151811062000b825762000b82620019ce565b6020026020010181905250808062000b9a9062001a2b565b91505062000b45565b5062000bae565b8291505b62000bbc60116000620014bb565b62000bca60146000620014de565b600060108190555b8481101562000d67576000604051806040016040528085848151811062000bfd5762000bfd620019ce565b6020026020010151602001516001600160a01b0316815260200185848151811062000c2c5762000c2c620019ce565b602090810291909101810151608001519091526011805460018101825560009190915282517f31ecc21a745e3968a04e9570e4425bc18fa8019c68028196b546d1669c200c68600290920291820180546001600160a01b0319166001600160a01b03909216919091179055908201517f31ecc21a745e3968a04e9570e4425bc18fa8019c68028196b546d1669c200c6990910155845190915060149085908490811062000cdd5762000cdd620019ce565b602090810291909101810151604001518254600181018455600093845292829020815162000d159491909101929190910190620014fe565b5083828151811062000d2b5762000d2b620019ce565b6020026020010151608001516010600082825462000d4a919062001a10565b9091555082915062000d5e90508162001a2b565b91505062000bd2565b505050505050565b600062000d8b82604001516200123a60201b620025a11760201c565b6001600160a01b0390911660208401529050801562000ddb5760405162461bcd60e51b815260206004820152600b60248201526a32b737b2329032b93937b960a91b60448201526064016200029b565b6020808301516001600160a01b039081166000908152601d909252604090912060010154161562000e4f5760405162461bcd60e51b815260206004820152601c60248201527f76616c696461746f7220616c726561647920726567697374657265640000000060448201526064016200029b565b61271082606001511115620003e35760405162461bcd60e51b815260206004820152601760248201527f696e76616c696420636f6d6d697373696f6e207261746500000000000000000060448201526064016200029b565b60c08101516001600160a01b031662000f1f5760208101518151606083015160405162000ed4906200158d565b6001600160a01b0393841681529290911660208301526040820152606001604051809103906000f08015801562000f0f573d6000803e3d6000fd5b506001600160a01b031660c08201525b60208082018051600d80546001818101835560009283527fd7b6990105719101dabeb77144f2a3385c8033acd3af97e9423a695e81ad1eb590910180546001600160a01b039485166001600160a01b031991821617909155845184168352601d86526040928390208751815490861690831617815594519185018054929094169116179091558301518051849362000fbf926002850192910190620014fe565b50606082015160038201556080820151600482015560a0820151600582015560c08201516006820180546001600160a01b0319166001600160a01b0390921691909117905560e08201516007820155610100820151600882015561012082015160098201805460ff1916600183818111156200103f576200103f620019e4565b0217905550505050565b600081815260166020908152604080832060018101546001600160a01b03168452601d90925282206004810154919290916200108b57506002820154620010b5565b816004015483600201548360070154620010a6919062001aa0565b620010b2919062001ad8565b90505b600682015483546040516340c10f1960e01b81526001600160a01b039182166004820152602481018490529116906340c10f1990604401600060405180830381600087803b1580156200110757600080fd5b505af11580156200111c573d6000803e3d6000fd5b5050505082600201548260040160008282546200113a919062001a10565b925050819055508082600701600082825462001157919062001a10565b909155505050505050565b600081815260196020908152604080832060018101546001600160a01b03168452601d909252822060078101546004820154600284015493949293620011a9919062001aa0565b620011b5919062001ad8565b905080826004016000828254620011cd919062001a49565b90915550506002830154600783018054600090620011ed90849062001a49565b909155505082546001600160a01b03166000908152601c6020526040812080548392906200115790849062001a10565b620003fd8160006001845162001234919062001a49565b62001283565b600080620012476200159b565b600060408286516020880160ff5afa6200126057600080fd5b5080516020909101516c0100000000000000000000000090910494909350915050565b81818082141562001295575050505050565b6000856002620012a6878762001aef565b620012b2919062001b34565b620012be908762001b68565b81518110620012d157620012d1620019ce565b60200260200101516080015190505b8183136200140e575b80868481518110620012ff57620012ff620019ce565b60200260200101516080015111156200132757826200131e8162001baf565b935050620012e9565b8582815181106200133c576200133c620019ce565b6020026020010151608001518111156200136557816200135c8162001bcb565b92505062001327565b8183136200140857858281518110620013825762001382620019ce565b60200260200101518684815181106200139f576200139f620019ce565b6020026020010151878581518110620013bc57620013bc620019ce565b60200260200101888581518110620013d857620013d8620019ce565b6020026020010182905282905250508280620013f49062001baf565b9350508180620014049062001bcb565b9250505b620012e0565b8185121562001424576200142486868462001283565b8383121562000d675762000d6786848662001283565b60405180610140016040528060006001600160a01b0316815260200160006001600160a01b031681526020016060815260200160008152602001600081526020016000815260200160006001600160a01b03168152602001600081526020016000815260200160006001811115620014b657620014b6620019e4565b905290565b5080546000825560020290600052602060002090810190620003fd9190620015b9565b5080546000825590600052602060002090810190620003fd9190620015e1565b8280546200150c9062001a63565b90600052602060002090601f0160209004810192826200153057600085556200157b565b82601f106200154b57805160ff19168380011785556200157b565b828001600101855582156200157b579182015b828111156200157b5782518255916020019190600101906200155e565b506200158992915062001602565b5090565b610cae8062006f4383390190565b60405180604001604052806002906020820280368337509192915050565b5b80821115620015895780546001600160a01b031916815560006001820155600201620015ba565b8082111562001589576000620015f8828262001619565b50600101620015e1565b5b8082111562001589576000815560010162001603565b508054620016279062001a63565b6000825580601f1062001638575050565b601f016020900490600052602060002090810190620003fd919062001602565b634e487b7160e01b600052604160045260246000fd5b60405161014081016001600160401b038111828210171562001694576200169462001658565b60405290565b604051601f8201601f191681016001600160401b0381118282101715620016c557620016c562001658565b604052919050565b80516001600160a01b0381168114620016e557600080fd5b919050565b600082601f830112620016fc57600080fd5b81516001600160401b0381111562001718576200171862001658565b60206200172e601f8301601f191682016200169a565b82815285828487010111156200174357600080fd5b60005b838110156200176357858101830151828201840152820162001746565b83811115620017755760008385840101525b5095945050505050565b805160028110620016e557600080fd5b60006101408284031215620017a357600080fd5b620017ad6200166e565b9050620017ba82620016cd565b8152620017ca60208301620016cd565b602082015260408201516040820152606082015160608201526080820151608082015260a082015160a082015260c082015160c082015260e082015160e082015261010080830151818301525061012080830151818301525092915050565b60008061016083850312156200183e57600080fd5b82516001600160401b03808211156200185657600080fd5b818501915085601f8301126200186b57600080fd5b815160208282111562001882576200188262001658565b8160051b620018938282016200169a565b928352848101820192828101908a851115620018ae57600080fd5b83870192505b84831015620019ad57825186811115620018cd57600080fd5b8701610140818d03601f1901811315620018e657600080fd5b620018f06200166e565b620018fd878401620016cd565b81526200190d60408401620016cd565b87820152606083015189811115620019255760008081fd5b620019358f8983870101620016ea565b604083015250608080840151606083015260a0808501518284015260c0915081850151818401525060e06200196c818601620016cd565b8284015261010091508185015181840152506101208085015182840152620019968486016200177f565b9083015250845250509183019190830190620018b4565b809850505050620019c1888289016200178f565b9450505050509250929050565b634e487b7160e01b600052603260045260246000fd5b634e487b7160e01b600052602160045260246000fd5b634e487b7160e01b600052601160045260246000fd5b6000821982111562001a265762001a26620019fa565b500190565b600060001982141562001a425762001a42620019fa565b5060010190565b60008282101562001a5e5762001a5e620019fa565b500390565b600181811c9082168062001a7857607f821691505b6020821081141562001a9a57634e487b7160e01b600052602260045260246000fd5b50919050565b600081600019048311821515161562001abd5762001abd620019fa565b500290565b634e487b7160e01b600052601260045260246000fd5b60008262001aea5762001aea62001ac2565b500490565b60008083128015600160ff1b85018412161562001b105762001b10620019fa565b6001600160ff1b038401831381161562001b2e5762001b2e620019fa565b50500390565b60008262001b465762001b4662001ac2565b600160ff1b82146000198414161562001b635762001b63620019fa565b500590565b600080821280156001600160ff1b038490038513161562001b8d5762001b8d620019fa565b600160ff1b839003841281161562001ba95762001ba9620019fa565b50500190565b60006001600160ff1b0382141562001a425762001a42620019fa565b6000600160ff1b82141562001be45762001be4620019fa565b506000190190565b6153478062001bfc6000396000f3fe6080604052600436106200031f5760003560e01c8063819b6463116200019f578063b2ea9adb11620000e5578063cb696f54116200009b578063d886f8a21162000075578063d886f8a21462000a49578063dd62ed3e1462000a6e578063e485c6fb1462000ab8578063e7f43c681462000add57005b8063cb696f5414620009ea578063cf9c57191462000a0f578063d5f394881462000a2757005b8063b2ea9adb1462000921578063b46e55201462000946578063b66b3e79146200096b578063b7ab4db51462000993578063c2362dd514620009ba578063c9d97af414620009d257005b80639dc29fac1162000155578063a8b2216e116200012f578063a8b2216e1462000896578063a9059cbb14620008bd578063ab8f6ffe14620008e2578063ae1f5fa0146200090957005b80639dc29fac1462000827578063a515366a146200084c578063a5d059ca146200087157005b8063819b64631462000775578063872cf059146200078c5780638bac7dad14620007a457806395d89b4114620007c95780639bb851c014620007f75780639c98e471146200080f57005b8063446972211162000265578063662cd7f4116200021b578063731b3a0311620001f5578063731b3a03146200068257806377e741c71462000699578063787a243314620006be57806379502c5514620006d657005b8063662cd7f4146200060b5780636b5f444c146200062357806370a08231146200064857005b806344697221146200051f5780634b0dff631462000537578063520fdbbc146200054f5780635333d404146200057457806355230e9314620005995780635f7d394914620005cd57005b8063114eaf5511620002d55780631904bb2e11620002af5780631904bb2e146200048957806323b872dd14620004bd5780632f2c3f2e14620004e257806340c10f1914620004fa57005b8063114eaf5514620004355780631604e416146200045a57806318160ddd146200047257005b806305261aea146200032957806306fdde031462000367578063095ea7b314620003a25780630ae65e7a14620003d85780630d8e6e2c14620003fd57806311220633146200041e57005b366200032757005b005b3480156200033657600080fd5b506200034e6200034836600462003d1b565b62000afd565b6040516200035e92919062003d86565b60405180910390f35b3480156200037457600080fd5b506040805180820190915260068152652732bbba37b760d11b60208201525b6040516200035e919062003e04565b348015620003af57600080fd5b50620003c7620003c136600462003e36565b62000c36565b60405190151581526020016200035e565b348015620003e557600080fd5b5062000327620003f736600462003e65565b62000c4f565b3480156200040a57600080fd5b50600b545b6040519081526020016200035e565b3480156200042b57600080fd5b506006546200040f565b3480156200044257600080fd5b50620003276200045436600462003d1b565b62000d08565b3480156200046757600080fd5b506200040f60135481565b3480156200047f57600080fd5b50601e546200040f565b3480156200049657600080fd5b50620004ae620004a836600462003e65565b62000d3a565b6040516200035e919062003ebe565b348015620004ca57600080fd5b50620003c7620004dc36600462003f88565b62000ea0565b348015620004ef57600080fd5b506200040f61271081565b3480156200050757600080fd5b50620003276200051936600462003e36565b62000efa565b3480156200052c57600080fd5b506200040f60185481565b3480156200054457600080fd5b506200040f601b5481565b3480156200055c57600080fd5b50620003276200056e36600462003e65565b62001022565b3480156200058157600080fd5b50620003276200059336600462004079565b62001071565b348015620005a657600080fd5b50620005be620005b8366004620040e4565b6200111a565b6040516200035e919062004107565b348015620005da57600080fd5b50620005f2620005ec366004620040e4565b62001238565b6040516001600160a01b0390911681526020016200035e565b3480156200061857600080fd5b506200040f601a5481565b3480156200063057600080fd5b50620003276200064236600462003d1b565b6200143e565b3480156200065557600080fd5b506200040f6200066736600462003e65565b6001600160a01b03166000908152601c602052604090205490565b3480156200068f57600080fd5b50600f546200040f565b348015620006a657600080fd5b5062000327620006b836600462003d1b565b62001470565b348015620006cb57600080fd5b506200040f60175481565b348015620006e357600080fd5b50600354600454600554600654600754600854600954600a54600b54600c546200071f996001600160a01b03908116991697969594939291908a565b604080516001600160a01b039b8c1681529a90991660208b0152978901969096526060880194909452608087019290925260a086015260c085015260e0840152610100830152610120820152610140016200035e565b3480156200078257600080fd5b50600a546200040f565b3480156200079957600080fd5b5062000327620014a2565b348015620007b157600080fd5b5062000327620007c336600462003d1b565b620014de565b348015620007d657600080fd5b50604080518082019091526003815262272a2760e91b602082015262000393565b3480156200080457600080fd5b506200040f60125481565b3480156200081c57600080fd5b506200040f60105481565b3480156200083457600080fd5b50620003276200084636600462003e36565b62001562565b3480156200085957600080fd5b50620003276200086b36600462003e36565b6200167c565b3480156200087e57600080fd5b50620003276200089036600462003e36565b62001779565b348015620008a357600080fd5b50620008ae620017f0565b6040516200035e919062004179565b348015620008ca57600080fd5b50620003c7620008dc36600462003e36565b620018d3565b348015620008ef57600080fd5b50620008fa620018e2565b6040516200035e9190620041df565b3480156200091657600080fd5b506200032762001950565b3480156200092e57600080fd5b50620003276200094036600462004079565b6200204b565b3480156200095357600080fd5b50620003276200096536600462003e65565b62002092565b3480156200097857600080fd5b5062000983620021e5565b6040516200035e929190620041f4565b348015620009a057600080fd5b50620009ab6200231c565b6040516200035e919062004226565b348015620009c757600080fd5b506200040f600f5481565b348015620009df57600080fd5b506200040f600e5481565b348015620009f757600080fd5b506200032762000a0936600462003d1b565b62002380565b34801562000a1c57600080fd5b5062000327620023e8565b34801562000a3457600080fd5b50601f54620005f2906001600160a01b031681565b34801562000a5657600080fd5b506200032762000a6836600462003e65565b6200243c565b34801562000a7b57600080fd5b506200040f62000a8d36600462004275565b6001600160a01b03918216600090815260156020908152604080832093909416825291909152205490565b34801562000ac557600080fd5b50620005be62000ad7366004620040e4565b6200248b565b34801562000aea57600080fd5b506003546001600160a01b0316620005f2565b601f546000906060906001600160a01b0316331462000b395760405162461bcd60e51b815260040162000b3090620042b3565b60405180910390fd5b826013600082825462000b4d91906200430c565b9091555050600854600f54439162000b65916200430c565b141562000bb45762000b79601354620025e1565b600060135562000b8862002845565b62000b9262001950565b43600f819055506001600e600082825462000bae91906200430c565b90915550505b600254601180546040805160208084028201810190925282815260ff9094169391839160009084015b8282101562000c27576000848152602090819020604080518082019091526002850290910180546001600160a01b0316825260019081015482840152908352909201910162000bdd565b50505050905091509150915091565b600062000c45338484620028f4565b5060015b92915050565b6001600160a01b038082166000818152601d60205260409020600101549091161462000cbe5760405162461bcd60e51b815260206004820152601c60248201527f76616c696461746f72206d757374206265207265676973746572656400000000604482015260640162000b30565b6001600160a01b038181166000908152601d602052604090205416331462000cfa5760405162461bcd60e51b815260040162000b309062004327565b62000d058162002a1d565b50565b6003546001600160a01b0316331462000d355760405162461bcd60e51b815260040162000b309062004373565b600955565b62000d4462003ac3565b6001600160a01b038083166000908152601d6020908152604091829020825161014081018452815485168152600182015490941691840191909152600281018054919284019162000d9590620043aa565b80601f016020809104026020016040519081016040528092919081815260200182805462000dc390620043aa565b801562000e145780601f1062000de85761010080835404028352916020019162000e14565b820191906000526020600020905b81548152906001019060200180831162000df657829003601f168201915b505050918352505060038201546020820152600482015460408201526005820154606082015260068201546001600160a01b03166080820152600782015460a0820152600882015460c0820152600982015460e09091019060ff16600181111562000e835762000e8362003e85565b600181111562000e975762000e9762003e85565b90525092915050565b600062000eaf84848462002b19565b6001600160a01b038416600090815260156020908152604080832033845290915281205462000ee0908490620043e7565b905062000eef853383620028f4565b506001949350505050565b6003546001600160a01b0316331462000f275760405162461bcd60e51b815260040162000b309062004373565b671000000000000000811062000f915760405162461bcd60e51b815260206004820152602860248201527f697373756564204e6577746f6e2063616e277420626520677265617465722074604482015267068616e20325e36360c41b606482015260840162000b30565b6001600160a01b0382166000908152601c60205260408120805483929062000fbb9084906200430c565b9250508190555080601e600082825462000fd691906200430c565b9091555050604080516001600160a01b0384168152602081018390527f48490b4407bb949b708ec5f514b4167f08f4969baaf78d53b05028adf369bfcf91015b60405180910390a15050565b6003546001600160a01b031633146200104f5760405162461bcd60e51b815260040162000b309062004373565b600380546001600160a01b0319166001600160a01b0392909216919091179055565b604080516101408101825233815260006020820181905291810184905260075460608201526080810182905260a0810182905260c0810182905260e0810182905243610100820152610120810191909152620010ce818362002c22565b7f6921859367aca5023ddf910758cd0cda74261b2e5c8425c253e9d03b62b950b8338260200151858460c001516040516200110d949392919062004401565b60405180910390a1505050565b606060006200112a8484620043e7565b67ffffffffffffffff81111562001145576200114562003fce565b6040519080825280602002602001820160405280156200118257816020015b6200116e62003b44565b815260200190600190039081620011645790505b50905060005b620011948585620043e7565b811015620012305760196000620011ac83886200430c565b81526020808201929092526040908101600020815160808101835281546001600160a01b0390811682526001830154169381019390935260028101549183019190915260030154606082015282518390839081106200120f576200120f62004441565b60200260200101819052508080620012279062004457565b91505062001188565b509392505050565b600080805b6011548110156200129457601181815481106200125e576200125e62004441565b906000526020600020906002020160010154826200127d91906200430c565b9150806200128b8162004457565b9150506200123d565b5080620012e45760405162461bcd60e51b815260206004820152601c60248201527f54686520636f6d6d6974746565206973206e6f74207374616b696e6700000000604482015260640162000b30565b6000620012f284866200430c565b90506000816040516020016200130a91815260200190565b60408051601f198184030181529190528051602090910120905060006200133284836200448b565b90506000805b601154811015620013e2576011818154811062001359576200135962004441565b906000526020600020906002020160010154826200137891906200430c565b915062001387600183620043e7565b8311620013cd5760118181548110620013a457620013a462004441565b60009182526020909120600290910201546001600160a01b0316965062000c4995505050505050565b80620013d98162004457565b91505062001338565b5060405162461bcd60e51b815260206004820152602960248201527f5468657265206973206e6f2076616c696461746f72206c65667420696e20746860448201526865206e6574776f726b60b81b606482015260840162000b30565b6003546001600160a01b031633146200146b5760405162461bcd60e51b815260040162000b309062004373565b600855565b6003546001600160a01b031633146200149d5760405162461bcd60e51b815260040162000b309062004373565b600555565b6003546001600160a01b03163314620014cf5760405162461bcd60e51b815260040162000b309062004373565b6002805460ff19166001179055565b6003546001600160a01b031633146200150b5760405162461bcd60e51b815260040162000b309062004373565b600081116200155d5760405162461bcd60e51b815260206004820152601960248201527f636f6d6d69747465652073697a652063616e2774206265203000000000000000604482015260640162000b30565b600a55565b6003546001600160a01b031633146200158f5760405162461bcd60e51b815260040162000b309062004373565b6001600160a01b0382166000908152601c6020526040902054811115620015f25760405162461bcd60e51b8152602060048201526016602482015275416d6f756e7420657863656564732062616c616e636560501b604482015260640162000b30565b6001600160a01b0382166000908152601c6020526040812080548392906200161c908490620043e7565b9250508190555080601e6000828254620016379190620043e7565b9091555050604080516001600160a01b0384168152602081018390527f5024dbeedf0c06664c9bd7be836915730c955e936972c020683dadf11d5488a3910162001016565b6001600160a01b038083166000818152601d602052604090206001015490911614620016e65760405162461bcd60e51b81526020600482015260186024820152771d985b1a59185d1bdc881b9bdd081c9959da5cdd195c995960421b604482015260640162000b30565b6001600160a01b0382166000908152601d602052604081206009015460ff16600181111562001719576200171962003e85565b14620017685760405162461bcd60e51b815260206004820152601b60248201527f76616c696461746f72206e65656420746f206265206163746976650000000000604482015260640162000b30565b6200177582823362002dde565b5050565b6001600160a01b038083166000818152601d602052604090206001015490911614620017e35760405162461bcd60e51b81526020600482015260186024820152771d985b1a59185d1bdc881b9bdd081c9959da5cdd195c995960421b604482015260640162000b30565b6200177582823362002f6a565b60606014805480602002602001604051908101604052809291908181526020016000905b82821015620018ca5783829060005260206000200180546200183690620043aa565b80601f01602080910402602001604051908101604052809291908181526020018280546200186490620043aa565b8015620018b55780601f106200188957610100808354040283529160200191620018b5565b820191906000526020600020905b8154815290600101906020018083116200189757829003601f168201915b50505050508152602001906001019062001814565b50505050905090565b600062000c4533848462002b19565b60606011805480602002602001604051908101604052809291908181526020016000905b82821015620018ca576000848152602090819020604080518082019091526002850290910180546001600160a01b0316825260019081015482840152908352909201910162001906565b601f546001600160a01b031633146200197d5760405162461bcd60e51b815260040162000b3090620042b3565b600d54620019ce5760405162461bcd60e51b815260206004820152601860248201527f5468657265206d7573742062652076616c696461746f72730000000000000000604482015260640162000b30565b6000805b600d5481101562001ab2576000601d6000600d8481548110620019f957620019f962004441565b60009182526020808320909101546001600160a01b0316835282019290925260400190206009015460ff16600181111562001a385762001a3862003e85565b14801562001a8757506000601d6000600d848154811062001a5d5762001a5d62004441565b60009182526020808320909101546001600160a01b03168352820192909252604001902060040154115b1562001a9d578162001a998162004457565b9250505b8062001aa98162004457565b915050620019d2565b50600a5481811062001ac15750805b60008267ffffffffffffffff81111562001adf5762001adf62003fce565b60405190808252806020026020018201604052801562001b1c57816020015b62001b0862003ac3565b81526020019060019003908162001afe5790505b50905060008267ffffffffffffffff81111562001b3d5762001b3d62003fce565b60405190808252806020026020018201604052801562001b7a57816020015b62001b6662003ac3565b81526020019060019003908162001b5c5790505b5090506000805b600d5481101562001e06576000601d6000600d848154811062001ba85762001ba862004441565b60009182526020808320909101546001600160a01b0316835282019290925260400190206009015460ff16600181111562001be75762001be762003e85565b14801562001c3657506000601d6000600d848154811062001c0c5762001c0c62004441565b60009182526020808320909101546001600160a01b03168352820192909252604001902060040154115b1562001df1576000601d6000600d848154811062001c585762001c5862004441565b60009182526020808320909101546001600160a01b039081168452838201949094526040928301909120825161014081018452815485168152600182015490941691840191909152600281018054919284019162001cb690620043aa565b80601f016020809104026020016040519081016040528092919081815260200182805462001ce490620043aa565b801562001d355780601f1062001d095761010080835404028352916020019162001d35565b820191906000526020600020905b81548152906001019060200180831162001d1757829003601f168201915b505050918352505060038201546020820152600482015460408201526005820154606082015260068201546001600160a01b03166080820152600782015460a0820152600882015460c0820152600982015460e09091019060ff16600181111562001da45762001da462003e85565b600181111562001db85762001db862003e85565b8152505090508085848151811062001dd45762001dd462004441565b6020026020010181905250828062001dec9062004457565b935050505b8062001dfd8162004457565b91505062001b81565b50600a548351111562001e865762001e1e836200315d565b60005b600a5481101562001e7f5783818151811062001e415762001e4162004441565b602002602001015183828151811062001e5e5762001e5e62004441565b6020026020010181905250808062001e769062004457565b91505062001e21565b5062001e8a565b8291505b62001e986011600062003b7e565b62001ea66014600062003ba1565b600060108190555b8481101562002043576000604051806040016040528085848151811062001ed95762001ed962004441565b6020026020010151602001516001600160a01b0316815260200185848151811062001f085762001f0862004441565b602090810291909101810151608001519091526011805460018101825560009190915282517f31ecc21a745e3968a04e9570e4425bc18fa8019c68028196b546d1669c200c68600290920291820180546001600160a01b0319166001600160a01b03909216919091179055908201517f31ecc21a745e3968a04e9570e4425bc18fa8019c68028196b546d1669c200c6990910155845190915060149085908490811062001fb95762001fb962004441565b602090810291909101810151604001518254600181018455600093845292829020815162001ff1949190910192919091019062003bc1565b5083828151811062002007576200200762004441565b602002602001015160800151601060008282546200202691906200430c565b909155508291506200203a90508162004457565b91505062001eae565b505050505050565b6003546001600160a01b03163314620020785760405162461bcd60e51b815260040162000b309062004373565b620020856000836200317a565b620017756001826200317a565b6001600160a01b038082166000818152601d602052604090206001015490911614620021015760405162461bcd60e51b815260206004820152601c60248201527f76616c696461746f72206d757374206265207265676973746572656400000000604482015260640162000b30565b6001600160a01b038181166000908152601d60205260409020541633146200213d5760405162461bcd60e51b815260040162000b309062004327565b60016001600160a01b0382166000908152601d602052604090206009015460ff16600181111562002172576200217262003e85565b14620021c15760405162461bcd60e51b815260206004820152601860248201527f76616c696461746f72206d757374206265207061757365640000000000000000604482015260640162000b30565b6001600160a01b03166000908152601d60205260409020600901805460ff19169055565b60608060006001818054620021fa90620043aa565b80601f01602080910402602001604051908101604052809291908181526020018280546200222890620043aa565b8015620022795780601f106200224d5761010080835404028352916020019162002279565b820191906000526020600020905b8154815290600101906020018083116200225b57829003601f168201915b505050505091508080546200228e90620043aa565b80601f0160208091040260200160405190810160405280929190818152602001828054620022bc90620043aa565b80156200230d5780601f10620022e1576101008083540402835291602001916200230d565b820191906000526020600020905b815481529060010190602001808311620022ef57829003601f168201915b50505050509050915091509091565b6060600d8054806020026020016040519081016040528092919081815260200182805480156200237657602002820191906000526020600020905b81546001600160a01b0316815260019091019060200180831162002357575b5050505050905090565b6003546001600160a01b03163314620023ad5760405162461bcd60e51b815260040162000b309062004373565b60068190556040518181527f1f4d2fc7529047a5bd96d3229bfea127fd18b7748f13586e097c69fccd3891289060200160405180910390a150565b6003546001600160a01b03163314620024155760405162461bcd60e51b815260040162000b309062004373565b6200242260008062003c50565b620024306001600062003c50565b6002805460ff19169055565b6003546001600160a01b03163314620024695760405162461bcd60e51b815260040162000b309062004373565b600480546001600160a01b0319166001600160a01b0392909216919091179055565b606060006200249b8484620043e7565b67ffffffffffffffff811115620024b657620024b662003fce565b604051908082528060200260200182016040528015620024f357816020015b620024df62003b44565b815260200190600190039081620024d55790505b50905060005b620025058585620043e7565b8110156200123057601660006200251d83886200430c565b81526020808201929092526040908101600020815160808101835281546001600160a01b03908116825260018301541693810193909352600281015491830191909152600301546060820152825183908390811062002580576200258062004441565b60200260200101819052508080620025989062004457565b915050620024f9565b600080620025ae62003c8f565b600060408286516020880160ff5afa620025c757600080fd5b508051602090910151600160601b90910494909350915050565b80471015620026465760405162461bcd60e51b815260206004820152602a60248201527f6e6f7420656e6f7567682066756e647320746f20706572666f726d207265646960448201526939ba3934b13aba34b7b760b11b606482015260840162000b30565b6000670de0b6b3a764000082600360020154620026649190620044a2565b620026709190620044c4565b90508015620026c4576004546040516001600160a01b039091169082156108fc029083906000818181858888f19350505050158015620026b4573d6000803e3d6000fd5b50620026c18183620043e7565b91505b8160126000828254620026d891906200430c565b90915550600090505b60115481101562002840576000601d60006011848154811062002708576200270862004441565b600091825260208083206002909202909101546001600160a01b0316835282019290925260400181206010546004820154919350906200274a908790620044a2565b620027569190620044c4565b90508015620027e0578160060160009054906101000a90046001600160a01b03166001600160a01b031663fb489a7b826040518263ffffffff1660e01b815260040160206040518083038185885af1158015620027b7573d6000803e3d6000fd5b50505050506040513d601f19601f82011682018060405250810190620027de9190620044db565b505b6001820154604080516001600160a01b039092168252602082018390527fb3b7a071186534c03b40695710096f289fd4ed6c1a374aff0bb648955e4fe563910160405180910390a150508080620028379062004457565b915050620026e1565b505050565b6017545b60185481101562002874576200285f81620032cd565b806200286b8162004457565b91505062002849565b50601854601755601a54805b601b54811015620028ee576009546000828152601960205260409020600301544391620028ad916200430c565b11620028d357620028be81620033e6565b620028cb6001836200430c565b9150620028d9565b620028ee565b80620028e58162004457565b91505062002880565b50601a55565b6001600160a01b038316620029585760405162461bcd60e51b8152602060048201526024808201527f45524332303a20617070726f76652066726f6d20746865207a65726f206164646044820152637265737360e01b606482015260840162000b30565b6001600160a01b038216620029bb5760405162461bcd60e51b815260206004820152602260248201527f45524332303a20617070726f766520746f20746865207a65726f206164647265604482015261737360f01b606482015260840162000b30565b6001600160a01b0383811660008181526015602090815260408083209487168084529482529182902085905590518481527f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b92591015b60405180910390a3505050565b6001600160a01b0381166000908152601d6020526040812090600982015460ff16600181111562002a525762002a5262003e85565b1462002aa15760405162461bcd60e51b815260206004820152601960248201527f76616c696461746f72206d75737420626520656e61626c656400000000000000604482015260640162000b30565b60098101805460ff191660011790558054600854600f547f75bdcdbe540758778e669d108fbcb7ede734f27f46e4e5525eeb8ecf91849a9c926001600160a01b031691859162002af291906200430c565b604080516001600160a01b0394851681529390921660208401529082015260600162001016565b6001600160a01b0383166000908152601c602052604090205481111562002b7c5760405162461bcd60e51b8152602060048201526016602482015275616d6f756e7420657863656564732062616c616e636560501b604482015260640162000b30565b6001600160a01b0383166000908152601c60205260408120805483929062002ba6908490620043e7565b90915550506001600160a01b0382166000908152601c60205260408120805483929062002bd59084906200430c565b92505081905550816001600160a01b0316836001600160a01b03167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef8360405162002a1091815260200190565b62002c2d82620034a1565b600080600062002c3d84620035ce565b604080518082018252601a81527f19457468657265756d205369676e6564204d6573736167653a0a0000000000006020808301919091528a51925160609390931b6bffffffffffffffffffffffff1916908301529396509194509250600090603401604051602081830303815290604052905060008262002cbf83516200364c565b8360405160200162002cd493929190620044f5565b60408051601f198184030181528282528051602080830191909120600080865291850180855281905260ff891693850193909352606084018a905260808401899052909350909160019060a0016020604051602081039080840390855afa15801562002d44573d6000803e3d6000fd5b50505060206040510351905089602001516001600160a01b0316816001600160a01b03161462002dc75760405162461bcd60e51b815260206004820152602760248201527f496e76616c69642070726f6f662070726f766964656420666f722072656769736044820152663a3930ba34b7b760c91b606482015260840162000b30565b62002dd28a6200376a565b50505050505050505050565b6000821162002e3c5760405162461bcd60e51b815260206004820152602360248201527f616d6f756e74206e65656420746f206265207374726963746c7920706f73697460448201526269766560e81b606482015260840162000b30565b6001600160a01b0381166000908152601c602052604090205482111562002ea65760405162461bcd60e51b815260206004820152601b60248201527f696e73756666696369656e74204e6577746f6e2062616c616e63650000000000604482015260640162000b30565b6001600160a01b0381166000908152601c60205260408120805484929062002ed0908490620043e7565b9091555050604080516080810182526001600160a01b038084168252858116602080840191825283850187815243606086019081526018805460009081526016909452968320865181549087166001600160a01b03199182161782559451600182018054919097169516949094179094555160028301559151600390910155825491929062002f5f8362004457565b919050555050505050565b6001600160a01b038381166000908152601d60205260408082206006015490516370a0823160e01b81528484166004820152919216906370a0823190602401602060405180830381865afa15801562002fc7573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062002fed9190620044db565b9050828110156200304c5760405162461bcd60e51b815260206004820152602260248201527f696e73756666696369656e74204c6971756964204e6577746f6e2062616c616e604482015261636560f01b606482015260840162000b30565b6001600160a01b038481166000908152601d602052604090819020600601549051632770a7eb60e21b8152848316600482015260248101869052911690639dc29fac90604401600060405180830381600087803b158015620030ad57600080fd5b505af1158015620030c2573d6000803e3d6000fd5b5050604080516080810182526001600160a01b03808716825288811660208084019182528385018a81524360608601908152601b805460009081526019909452968320865181549087166001600160a01b031991821617825594516001820180549190971695169490941790945551600283015591516003909101558254919450909250620031518362004457565b91905055505050505050565b62000d0581600060018451620031749190620043e7565b6200390c565b8154600260018083161561010002038216048251808201602081106020841001600281146200322957600181146200324f578660005260208404602060002001600160028402018855602085068060200390508088018589016001836101000a0392508282511684540184556001840193506020820191505b80821015620032125781518455600184019350602082019150620031f3565b815191036101000a908190040290915550620032c4565b60028302826020036101000a846020036101000a602089015104020185018755620032c4565b8660005260208404602060002001600160028402018855846020038088018589016001836101000a0392508282511660ff198a160184556020820191506001840193505b80821015620032b2578151845560018401935060208201915062003293565b815191036101000a9081900402909155505b50505050505050565b600081815260166020908152604080832060018101546001600160a01b03168452601d90925282206004810154919290916200330f5750600282015462003339565b8160040154836002015483600701546200332a9190620044a2565b620033369190620044c4565b90505b600682015483546040516340c10f1960e01b81526001600160a01b039182166004820152602481018490529116906340c10f1990604401600060405180830381600087803b1580156200338b57600080fd5b505af1158015620033a0573d6000803e3d6000fd5b505050508260020154826004016000828254620033be91906200430c565b9250508190555080826007016000828254620033db91906200430c565b909155505050505050565b600081815260196020908152604080832060018101546001600160a01b03168452601d9092528220600781015460048201546002840154939492936200342d9190620044a2565b620034399190620044c4565b905080826004016000828254620034519190620043e7565b9091555050600283015460078301805460009062003471908490620043e7565b909155505082546001600160a01b03166000908152601c602052604081208054839290620033db9084906200430c565b6000620034b28260400151620025a1565b6001600160a01b03909116602084015290508015620035025760405162461bcd60e51b815260206004820152600b60248201526a32b737b2329032b93937b960a91b604482015260640162000b30565b6020808301516001600160a01b039081166000908152601d9092526040909120600101541615620035765760405162461bcd60e51b815260206004820152601c60248201527f76616c696461746f7220616c7265616479207265676973746572656400000000604482015260640162000b30565b61271082606001511115620017755760405162461bcd60e51b815260206004820152601760248201527f696e76616c696420636f6d6d697373696f6e2072617465000000000000000000604482015260640162000b30565b60008060008351604114620036165760405162461bcd60e51b815260206004820152600d60248201526c34b73b30b634b210383937b7b360991b604482015260640162000b30565b50505060208101516040820151606083015160001a601b811015620036455762003642601b826200453e565b90505b9193909250565b606081620036715750506040805180820190915260018152600360fc1b602082015290565b8160005b8115620036a15780620036888162004457565b9150620036999050600a83620044c4565b915062003675565b60008167ffffffffffffffff811115620036bf57620036bf62003fce565b6040519080825280601f01601f191660200182016040528015620036ea576020820181803683370190505b5090505b8415620037625762003702600183620043e7565b915062003711600a866200448b565b6200371e9060306200430c565b60f81b81838151811062003736576200373662004441565b60200101906001600160f81b031916908160001a9053506200375a600a86620044c4565b9450620036ee565b949350505050565b60c08101516001600160a01b0316620037e257602081015181516060830151604051620037979062003cad565b6001600160a01b0393841681529290911660208301526040820152606001604051809103906000f080158015620037d2573d6000803e3d6000fd5b506001600160a01b031660c08201525b60208082018051600d80546001818101835560009283527fd7b6990105719101dabeb77144f2a3385c8033acd3af97e9423a695e81ad1eb590910180546001600160a01b039485166001600160a01b031991821617909155845184168352601d8652604092839020875181549086169083161781559451918501805492909416911617909155830151805184936200388292600285019291019062003bc1565b50606082015160038201556080820151600482015560a0820151600582015560c08201516006820180546001600160a01b0319166001600160a01b0390921691909117905560e08201516007820155610100820151600882015561012082015160098201805460ff19166001838181111562003902576200390262003e85565b0217905550505050565b8181808214156200391e575050505050565b60008560026200392f878762004566565b6200393b9190620045ab565b620039479087620045df565b815181106200395a576200395a62004441565b60200260200101516080015190505b81831362003a97575b8086848151811062003988576200398862004441565b6020026020010151608001511115620039b05782620039a78162004626565b93505062003972565b858281518110620039c557620039c562004441565b602002602001015160800151811115620039ee5781620039e58162004642565b925050620039b0565b81831362003a915785828151811062003a0b5762003a0b62004441565b602002602001015186848151811062003a285762003a2862004441565b602002602001015187858151811062003a455762003a4562004441565b6020026020010188858151811062003a615762003a6162004441565b602002602001018290528290525050828062003a7d9062004626565b935050818062003a8d9062004642565b9250505b62003969565b8185121562003aad5762003aad8686846200390c565b838312156200204357620020438684866200390c565b60405180610140016040528060006001600160a01b0316815260200160006001600160a01b031681526020016060815260200160008152602001600081526020016000815260200160006001600160a01b0316815260200160008152602001600081526020016000600181111562003b3f5762003b3f62003e85565b905290565b604051806080016040528060006001600160a01b0316815260200160006001600160a01b0316815260200160008152602001600081525090565b508054600082556002029060005260206000209081019062000d05919062003cbb565b508054600082559060005260206000209081019062000d05919062003ce3565b82805462003bcf90620043aa565b90600052602060002090601f01602090048101928262003bf3576000855562003c3e565b82601f1062003c0e57805160ff191683800117855562003c3e565b8280016001018555821562003c3e579182015b8281111562003c3e57825182559160200191906001019062003c21565b5062003c4c92915062003d04565b5090565b50805462003c5e90620043aa565b6000825580601f1062003c6f575050565b601f01602090049060005260206000209081019062000d05919062003d04565b60405180604001604052806002906020820280368337509192915050565b610cae806200466483390190565b5b8082111562003c4c5780546001600160a01b03191681556000600182015560020162003cbc565b8082111562003c4c57600062003cfa828262003c50565b5060010162003ce3565b5b8082111562003c4c576000815560010162003d05565b60006020828403121562003d2e57600080fd5b5035919050565b600081518084526020808501945080840160005b8381101562003d7b57815180516001600160a01b03168852830151838801526040909601959082019060010162003d49565b509495945050505050565b821515815260406020820152600062003762604083018462003d35565b60005b8381101562003dc057818101518382015260200162003da6565b8381111562003dd0576000848401525b50505050565b6000815180845262003df081602086016020860162003da3565b601f01601f19169290920160200192915050565b60208152600062003e19602083018462003dd6565b9392505050565b6001600160a01b038116811462000d0557600080fd5b6000806040838503121562003e4a57600080fd5b823562003e578162003e20565b946020939093013593505050565b60006020828403121562003e7857600080fd5b813562003e198162003e20565b634e487b7160e01b600052602160045260246000fd5b6002811062003eba57634e487b7160e01b600052602160045260246000fd5b9052565b6020815262003ed96020820183516001600160a01b03169052565b6000602083015162003ef660408401826001600160a01b03169052565b50604083015161014080606085015262003f1561016085018362003dd6565b915060608501516080850152608085015160a085015260a085015160c085015260c085015162003f5060e08601826001600160a01b03169052565b5060e0850151610100858101919091528501516101208086019190915285015162003f7e8286018262003e9b565b5090949350505050565b60008060006060848603121562003f9e57600080fd5b833562003fab8162003e20565b9250602084013562003fbd8162003e20565b929592945050506040919091013590565b634e487b7160e01b600052604160045260246000fd5b600082601f83011262003ff657600080fd5b813567ffffffffffffffff8082111562004014576200401462003fce565b604051601f8301601f19908116603f011681019082821181831017156200403f576200403f62003fce565b816040528381528660208588010111156200405957600080fd5b836020870160208301376000602085830101528094505050505092915050565b600080604083850312156200408d57600080fd5b823567ffffffffffffffff80821115620040a657600080fd5b620040b48683870162003fe4565b93506020850135915080821115620040cb57600080fd5b50620040da8582860162003fe4565b9150509250929050565b60008060408385031215620040f857600080fd5b50508035926020909101359150565b602080825282518282018190526000919060409081850190868401855b828110156200416c57815180516001600160a01b0390811686528782015116878601528581015186860152606090810151908501526080909301929085019060010162004124565b5091979650505050505050565b6000602080830181845280855180835260408601915060408160051b870101925083870160005b82811015620041d257603f19888603018452620041bf85835162003dd6565b94509285019290850190600101620041a0565b5092979650505050505050565b60208152600062003e19602083018462003d35565b60408152600062004209604083018562003dd6565b82810360208401526200421d818562003dd6565b95945050505050565b6020808252825182820181905260009190848201906040850190845b81811015620042695783516001600160a01b03168352928401929184019160010162004242565b50909695505050505050565b600080604083850312156200428957600080fd5b8235620042968162003e20565b91506020830135620042a88162003e20565b809150509250929050565b60208082526023908201527f66756e6374696f6e207265737472696374656420746f207468652070726f746f60408201526218dbdb60ea1b606082015260800190565b634e487b7160e01b600052601160045260246000fd5b60008219821115620043225762004322620042f6565b500190565b6020808252602c908201527f726571756972652063616c6c657220746f2062652076616c696461746f72206160408201526b191b5a5b881858d8dbdd5b9d60a21b606082015260800190565b6020808252601a908201527f63616c6c6572206973206e6f7420746865206f70657261746f72000000000000604082015260600190565b600181811c90821680620043bf57607f821691505b60208210811415620043e157634e487b7160e01b600052602260045260246000fd5b50919050565b600082821015620043fc57620043fc620042f6565b500390565b600060018060a01b0380871683528086166020840152608060408401526200442d608084018662003dd6565b915080841660608401525095945050505050565b634e487b7160e01b600052603260045260246000fd5b60006000198214156200446e576200446e620042f6565b5060010190565b634e487b7160e01b600052601260045260246000fd5b6000826200449d576200449d62004475565b500690565b6000816000190483118215151615620044bf57620044bf620042f6565b500290565b600082620044d657620044d662004475565b500490565b600060208284031215620044ee57600080fd5b5051919050565b600084516200450981846020890162003da3565b8451908301906200451f81836020890162003da3565b84519101906200453481836020880162003da3565b0195945050505050565b600060ff821660ff84168060ff038211156200455e576200455e620042f6565b019392505050565b60008083128015600160ff1b850184121615620045875762004587620042f6565b6001600160ff1b0384018313811615620045a557620045a5620042f6565b50500390565b600082620045bd57620045bd62004475565b600160ff1b821460001984141615620045da57620045da620042f6565b500590565b600080821280156001600160ff1b0384900385131615620046045762004604620042f6565b600160ff1b8390038412811615620046205762004620620042f6565b50500190565b60006001600160ff1b038214156200446e576200446e620042f6565b6000600160ff1b8214156200465b576200465b620042f6565b50600019019056fe608060405234801561001057600080fd5b50604051610cae380380610cae83398101604081905261002f9161009d565b61271081111561003e57600080fd5b600780546001600160a01b039485166001600160a01b03199182161790915560088054939094169281169290921790925560099190915560008054909116331790556100e0565b6001600160a01b038116811461009a57600080fd5b50565b6000806000606084860312156100b257600080fd5b83516100bd81610085565b60208501519093506100ce81610085565b80925050604084015190509250925092565b610bbf806100ef6000396000f3fe6080604052600436106100c25760003560e01c806340c10f191161007f5780639dc29fac116100595780639dc29fac146101f6578063a9059cbb14610216578063dd62ed3e14610236578063fb489a7b1461027c57600080fd5b806340c10f191461018057806370a08231146101a0578063949813b8146101d657600080fd5b8063095ea7b3146100c757806318160ddd146100fc578063187cf4d71461011b57806323b872dd146101335780632f2c3f2e14610153578063372500ab14610169575b600080fd5b3480156100d357600080fd5b506100e76100e23660046109e0565b610284565b60405190151581526020015b60405180910390f35b34801561010857600080fd5b506003545b6040519081526020016100f3565b34801561012757600080fd5b5061010d633b9aca0081565b34801561013f57600080fd5b506100e761014e366004610a0a565b61029a565b34801561015f57600080fd5b5061010d61271081565b34801561017557600080fd5b5061017e610392565b005b34801561018c57600080fd5b5061017e61019b3660046109e0565b610440565b3480156101ac57600080fd5b5061010d6101bb366004610a46565b6001600160a01b031660009081526001602052604090205490565b3480156101e257600080fd5b5061010d6101f1366004610a46565b6104a8565b34801561020257600080fd5b5061017e6102113660046109e0565b6104dc565b34801561022257600080fd5b506100e76102313660046109e0565b61053c565b34801561024257600080fd5b5061010d610251366004610a68565b6001600160a01b03918216600090815260026020908152604080832093909416825291909152205490565b61010d610589565b60006102913384846106d1565b50600192915050565b6001600160a01b0383166000908152600260209081526040808320338452909152812054828110156103245760405162461bcd60e51b815260206004820152602860248201527f45524332303a207472616e7366657220616d6f756e74206578636565647320616044820152676c6c6f77616e636560c01b60648201526084015b60405180910390fd5b61033885336103338685610ab1565b6106d1565b61034285846107f5565b61034c8484610898565b836001600160a01b0316856001600160a01b0316600080516020610b6a8339815191528560405161037f91815260200190565b60405180910390a3506001949350505050565b600061039d336108ec565b33600081815260046020526040808220829055519293509183908381818185875af1925050503d80600081146103ef576040519150601f19603f3d011682016040523d82523d6000602084013e6103f4565b606091505b505090508061043c5760405162461bcd60e51b81526020600482015260146024820152732330b4b632b2103a379039b2b7321022ba3432b960611b604482015260640161031b565b5050565b6000546001600160a01b0316331461046a5760405162461bcd60e51b815260040161031b90610ac8565b6104748282610898565b6040518181526001600160a01b03831690600090600080516020610b6a833981519152906020015b60405180910390a35050565b60006104b382610951565b6001600160a01b0383166000908152600460205260409020546104d69190610b10565b92915050565b6000546001600160a01b031633146105065760405162461bcd60e51b815260040161031b90610ac8565b61051082826107f5565b6040518181526000906001600160a01b03841690600080516020610b6a8339815191529060200161049c565b600061054833836107f5565b6105528383610898565b6040518281526001600160a01b038416903390600080516020610b6a8339815191529060200160405180910390a350600192915050565b600080546001600160a01b031633146105b45760405162461bcd60e51b815260040161031b90610ac8565b6009543490600090612710906105ca9084610b28565b6105d49190610b47565b90508181106106255760405162461bcd60e51b815260206004820152601860248201527f696e76616c69642076616c696461746f72207265776172640000000000000000604482015260640161031b565b61062f8183610ab1565b6008546040519193506001600160a01b03169082156108fc029083906000818181858888f1935050505015801561066a573d6000803e3d6000fd5b5060035460009061067f633b9aca0085610b28565b6106899190610b47565b9050806006546106999190610b10565b600655600354600090633b9aca00906106b29084610b28565b6106bc9190610b47565b90506106c88184610b10565b94505050505090565b6001600160a01b0383166107335760405162461bcd60e51b8152602060048201526024808201527f45524332303a20617070726f76652066726f6d20746865207a65726f206164646044820152637265737360e01b606482015260840161031b565b6001600160a01b0382166107945760405162461bcd60e51b815260206004820152602260248201527f45524332303a20617070726f766520746f20746865207a65726f206164647265604482015261737360f01b606482015260840161031b565b6001600160a01b0383811660008181526002602090815260408083209487168084529482529182902085905590518481527f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925910160405180910390a3505050565b6107fe826108ec565b506001600160a01b0382166000908152600160205260409020548082111561082557600080fd5b80821015610855576108378282610ab1565b6001600160a01b03841660009081526001602052604090205561087c565b6001600160a01b038316600090815260016020908152604080832083905560059091528120555b816003600082825461088e9190610ab1565b9091555050505050565b6108a1826108ec565b506001600160a01b038216600090815260016020526040812080548392906108ca908490610b10565b9250508190555080600360008282546108e39190610b10565b90915550505050565b6000806108f883610951565b6001600160a01b03841660009081526004602052604090205490915061091f908290610b10565b6001600160a01b0390931660009081526004602090815260408083208690556006546005909252909120555090919050565b6001600160a01b038116600090815260016020526040812054806109785750600092915050565b6001600160a01b03831660009081526005602052604081205460065461099e9190610ab1565b90506000633b9aca006109b18484610b28565b6109bb9190610b47565b95945050505050565b80356001600160a01b03811681146109db57600080fd5b919050565b600080604083850312156109f357600080fd5b6109fc836109c4565b946020939093013593505050565b600080600060608486031215610a1f57600080fd5b610a28846109c4565b9250610a36602085016109c4565b9150604084013590509250925092565b600060208284031215610a5857600080fd5b610a61826109c4565b9392505050565b60008060408385031215610a7b57600080fd5b610a84836109c4565b9150610a92602084016109c4565b90509250929050565b634e487b7160e01b600052601160045260246000fd5b600082821015610ac357610ac3610a9b565b500390565b60208082526028908201527f43616c6c207265737472696374656420746f20746865204175746f6e6974792060408201526710dbdb9d1c9858dd60c21b606082015260800190565b60008219821115610b2357610b23610a9b565b500190565b6000816000190483118215151615610b4257610b42610a9b565b500290565b600082610b6457634e487b7160e01b600052601260045260246000fd5b50049056feddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3efa2646970667358221220f9c95181bc1e859c0df200205d8951558d0b06d33a168b9f9240fc448aaea28064736f6c634300080b0033a26469706673582212208b635da6038f2e9d995dcb543237a2188370fe5d9ceafe1f535222bde39f69cc64736f6c634300080b0033608060405234801561001057600080fd5b50604051610cae380380610cae83398101604081905261002f9161009d565b61271081111561003e57600080fd5b600780546001600160a01b039485166001600160a01b03199182161790915560088054939094169281169290921790925560099190915560008054909116331790556100e0565b6001600160a01b038116811461009a57600080fd5b50565b6000806000606084860312156100b257600080fd5b83516100bd81610085565b60208501519093506100ce81610085565b80925050604084015190509250925092565b610bbf806100ef6000396000f3fe6080604052600436106100c25760003560e01c806340c10f191161007f5780639dc29fac116100595780639dc29fac146101f6578063a9059cbb14610216578063dd62ed3e14610236578063fb489a7b1461027c57600080fd5b806340c10f191461018057806370a08231146101a0578063949813b8146101d657600080fd5b8063095ea7b3146100c757806318160ddd146100fc578063187cf4d71461011b57806323b872dd146101335780632f2c3f2e14610153578063372500ab14610169575b600080fd5b3480156100d357600080fd5b506100e76100e23660046109e0565b610284565b60405190151581526020015b60405180910390f35b34801561010857600080fd5b506003545b6040519081526020016100f3565b34801561012757600080fd5b5061010d633b9aca0081565b34801561013f57600080fd5b506100e761014e366004610a0a565b61029a565b34801561015f57600080fd5b5061010d61271081565b34801561017557600080fd5b5061017e610392565b005b34801561018c57600080fd5b5061017e61019b3660046109e0565b610440565b3480156101ac57600080fd5b5061010d6101bb366004610a46565b6001600160a01b031660009081526001602052604090205490565b3480156101e257600080fd5b5061010d6101f1366004610a46565b6104a8565b34801561020257600080fd5b5061017e6102113660046109e0565b6104dc565b34801561022257600080fd5b506100e76102313660046109e0565b61053c565b34801561024257600080fd5b5061010d610251366004610a68565b6001600160a01b03918216600090815260026020908152604080832093909416825291909152205490565b61010d610589565b60006102913384846106d1565b50600192915050565b6001600160a01b0383166000908152600260209081526040808320338452909152812054828110156103245760405162461bcd60e51b815260206004820152602860248201527f45524332303a207472616e7366657220616d6f756e74206578636565647320616044820152676c6c6f77616e636560c01b60648201526084015b60405180910390fd5b61033885336103338685610ab1565b6106d1565b61034285846107f5565b61034c8484610898565b836001600160a01b0316856001600160a01b0316600080516020610b6a8339815191528560405161037f91815260200190565b60405180910390a3506001949350505050565b600061039d336108ec565b33600081815260046020526040808220829055519293509183908381818185875af1925050503d80600081146103ef576040519150601f19603f3d011682016040523d82523d6000602084013e6103f4565b606091505b505090508061043c5760405162461bcd60e51b81526020600482015260146024820152732330b4b632b2103a379039b2b7321022ba3432b960611b604482015260640161031b565b5050565b6000546001600160a01b0316331461046a5760405162461bcd60e51b815260040161031b90610ac8565b6104748282610898565b6040518181526001600160a01b03831690600090600080516020610b6a833981519152906020015b60405180910390a35050565b60006104b382610951565b6001600160a01b0383166000908152600460205260409020546104d69190610b10565b92915050565b6000546001600160a01b031633146105065760405162461bcd60e51b815260040161031b90610ac8565b61051082826107f5565b6040518181526000906001600160a01b03841690600080516020610b6a8339815191529060200161049c565b600061054833836107f5565b6105528383610898565b6040518281526001600160a01b038416903390600080516020610b6a8339815191529060200160405180910390a350600192915050565b600080546001600160a01b031633146105b45760405162461bcd60e51b815260040161031b90610ac8565b6009543490600090612710906105ca9084610b28565b6105d49190610b47565b90508181106106255760405162461bcd60e51b815260206004820152601860248201527f696e76616c69642076616c696461746f72207265776172640000000000000000604482015260640161031b565b61062f8183610ab1565b6008546040519193506001600160a01b03169082156108fc029083906000818181858888f1935050505015801561066a573d6000803e3d6000fd5b5060035460009061067f633b9aca0085610b28565b6106899190610b47565b9050806006546106999190610b10565b600655600354600090633b9aca00906106b29084610b28565b6106bc9190610b47565b90506106c88184610b10565b94505050505090565b6001600160a01b0383166107335760405162461bcd60e51b8152602060048201526024808201527f45524332303a20617070726f76652066726f6d20746865207a65726f206164646044820152637265737360e01b606482015260840161031b565b6001600160a01b0382166107945760405162461bcd60e51b815260206004820152602260248201527f45524332303a20617070726f766520746f20746865207a65726f206164647265604482015261737360f01b606482015260840161031b565b6001600160a01b0383811660008181526002602090815260408083209487168084529482529182902085905590518481527f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925910160405180910390a3505050565b6107fe826108ec565b506001600160a01b0382166000908152600160205260409020548082111561082557600080fd5b80821015610855576108378282610ab1565b6001600160a01b03841660009081526001602052604090205561087c565b6001600160a01b038316600090815260016020908152604080832083905560059091528120555b816003600082825461088e9190610ab1565b9091555050505050565b6108a1826108ec565b506001600160a01b038216600090815260016020526040812080548392906108ca908490610b10565b9250508190555080600360008282546108e39190610b10565b90915550505050565b6000806108f883610951565b6001600160a01b03841660009081526004602052604090205490915061091f908290610b10565b6001600160a01b0390931660009081526004602090815260408083208690556006546005909252909120555090919050565b6001600160a01b038116600090815260016020526040812054806109785750600092915050565b6001600160a01b03831660009081526005602052604081205460065461099e9190610ab1565b90506000633b9aca006109b18484610b28565b6109bb9190610b47565b95945050505050565b80356001600160a01b03811681146109db57600080fd5b919050565b600080604083850312156109f357600080fd5b6109fc836109c4565b946020939093013593505050565b600080600060608486031215610a1f57600080fd5b610a28846109c4565b9250610a36602085016109c4565b9150604084013590509250925092565b600060208284031215610a5857600080fd5b610a61826109c4565b9392505050565b60008060408385031215610a7b57600080fd5b610a84836109c4565b9150610a92602084016109c4565b90509250929050565b634e487b7160e01b600052601160045260246000fd5b600082821015610ac357610ac3610a9b565b500390565b60208082526028908201527f43616c6c207265737472696374656420746f20746865204175746f6e6974792060408201526710dbdb9d1c9858dd60c21b606082015260800190565b60008219821115610b2357610b23610a9b565b500190565b6000816000190483118215151615610b4257610b42610a9b565b500290565b600082610b6457634e487b7160e01b600052601260045260246000fd5b50049056feddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3efa2646970667358221220f9c95181bc1e859c0df200205d8951558d0b06d33a168b9f9240fc448aaea28064736f6c634300080b0033",
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

// EpochReward is a free data retrieval call binding the contract method 0x1604e416.
//
// Solidity: function epochReward() view returns(uint256)
func (_Autonity *AutonityCaller) EpochReward(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Autonity.contract.Call(opts, &out, "epochReward")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// EpochReward is a free data retrieval call binding the contract method 0x1604e416.
//
// Solidity: function epochReward() view returns(uint256)
func (_Autonity *AutonitySession) EpochReward() (*big.Int, error) {
	return _Autonity.Contract.EpochReward(&_Autonity.CallOpts)
}

// EpochReward is a free data retrieval call binding the contract method 0x1604e416.
//
// Solidity: function epochReward() view returns(uint256)
func (_Autonity *AutonityCallerSession) EpochReward() (*big.Int, error) {
	return _Autonity.Contract.EpochReward(&_Autonity.CallOpts)
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

// GetProposer is a free data retrieval call binding the contract method 0x5f7d3949.
//
// Solidity: function getProposer(uint256 height, uint256 round) view returns(address)
func (_Autonity *AutonityCaller) GetProposer(opts *bind.CallOpts, height *big.Int, round *big.Int) (common.Address, error) {
	var out []interface{}
	err := _Autonity.contract.Call(opts, &out, "getProposer", height, round)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetProposer is a free data retrieval call binding the contract method 0x5f7d3949.
//
// Solidity: function getProposer(uint256 height, uint256 round) view returns(address)
func (_Autonity *AutonitySession) GetProposer(height *big.Int, round *big.Int) (common.Address, error) {
	return _Autonity.Contract.GetProposer(&_Autonity.CallOpts, height, round)
}

// GetProposer is a free data retrieval call binding the contract method 0x5f7d3949.
//
// Solidity: function getProposer(uint256 height, uint256 round) view returns(address)
func (_Autonity *AutonityCallerSession) GetProposer(height *big.Int, round *big.Int) (common.Address, error) {
	return _Autonity.Contract.GetProposer(&_Autonity.CallOpts, height, round)
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
// Solidity: function getValidator(address _addr) view returns((address,address,string,uint256,uint256,uint256,address,uint256,uint256,uint8))
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
// Solidity: function getValidator(address _addr) view returns((address,address,string,uint256,uint256,uint256,address,uint256,uint256,uint8))
func (_Autonity *AutonitySession) GetValidator(_addr common.Address) (AutonityValidator, error) {
	return _Autonity.Contract.GetValidator(&_Autonity.CallOpts, _addr)
}

// GetValidator is a free data retrieval call binding the contract method 0x1904bb2e.
//
// Solidity: function getValidator(address _addr) view returns((address,address,string,uint256,uint256,uint256,address,uint256,uint256,uint8))
func (_Autonity *AutonityCallerSession) GetValidator(_addr common.Address) (AutonityValidator, error) {
	return _Autonity.Contract.GetValidator(&_Autonity.CallOpts, _addr)
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
// Solidity: function computeCommittee() returns()
func (_Autonity *AutonityTransactor) ComputeCommittee(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Autonity.contract.Transact(opts, "computeCommittee")
}

// ComputeCommittee is a paid mutator transaction binding the contract method 0xae1f5fa0.
//
// Solidity: function computeCommittee() returns()
func (_Autonity *AutonitySession) ComputeCommittee() (*types.Transaction, error) {
	return _Autonity.Contract.ComputeCommittee(&_Autonity.TransactOpts)
}

// ComputeCommittee is a paid mutator transaction binding the contract method 0xae1f5fa0.
//
// Solidity: function computeCommittee() returns()
func (_Autonity *AutonityTransactorSession) ComputeCommittee() (*types.Transaction, error) {
	return _Autonity.Contract.ComputeCommittee(&_Autonity.TransactOpts)
}

// Finalize is a paid mutator transaction binding the contract method 0x05261aea.
//
// Solidity: function finalize(uint256 amount) returns(bool, (address,uint256)[])
func (_Autonity *AutonityTransactor) Finalize(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _Autonity.contract.Transact(opts, "finalize", amount)
}

// Finalize is a paid mutator transaction binding the contract method 0x05261aea.
//
// Solidity: function finalize(uint256 amount) returns(bool, (address,uint256)[])
func (_Autonity *AutonitySession) Finalize(amount *big.Int) (*types.Transaction, error) {
	return _Autonity.Contract.Finalize(&_Autonity.TransactOpts, amount)
}

// Finalize is a paid mutator transaction binding the contract method 0x05261aea.
//
// Solidity: function finalize(uint256 amount) returns(bool, (address,uint256)[])
func (_Autonity *AutonityTransactorSession) Finalize(amount *big.Int) (*types.Transaction, error) {
	return _Autonity.Contract.Finalize(&_Autonity.TransactOpts, amount)
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

// RegisterValidator is a paid mutator transaction binding the contract method 0x5333d404.
//
// Solidity: function registerValidator(string _enode, bytes _proof) returns()
func (_Autonity *AutonityTransactor) RegisterValidator(opts *bind.TransactOpts, _enode string, _proof []byte) (*types.Transaction, error) {
	return _Autonity.contract.Transact(opts, "registerValidator", _enode, _proof)
}

// RegisterValidator is a paid mutator transaction binding the contract method 0x5333d404.
//
// Solidity: function registerValidator(string _enode, bytes _proof) returns()
func (_Autonity *AutonitySession) RegisterValidator(_enode string, _proof []byte) (*types.Transaction, error) {
	return _Autonity.Contract.RegisterValidator(&_Autonity.TransactOpts, _enode, _proof)
}

// RegisterValidator is a paid mutator transaction binding the contract method 0x5333d404.
//
// Solidity: function registerValidator(string _enode, bytes _proof) returns()
func (_Autonity *AutonityTransactorSession) RegisterValidator(_enode string, _proof []byte) (*types.Transaction, error) {
	return _Autonity.Contract.RegisterValidator(&_Autonity.TransactOpts, _enode, _proof)
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
	Enode          string
	LiquidContract common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterRegisteredValidator is a free log retrieval operation binding the contract event 0x6921859367aca5023ddf910758cd0cda74261b2e5c8425c253e9d03b62b950b8.
//
// Solidity: event RegisteredValidator(address treasury, address addr, string enode, address liquidContract)
func (_Autonity *AutonityFilterer) FilterRegisteredValidator(opts *bind.FilterOpts) (*AutonityRegisteredValidatorIterator, error) {

	logs, sub, err := _Autonity.contract.FilterLogs(opts, "RegisteredValidator")
	if err != nil {
		return nil, err
	}
	return &AutonityRegisteredValidatorIterator{contract: _Autonity.contract, event: "RegisteredValidator", logs: logs, sub: sub}, nil
}

// WatchRegisteredValidator is a free log subscription operation binding the contract event 0x6921859367aca5023ddf910758cd0cda74261b2e5c8425c253e9d03b62b950b8.
//
// Solidity: event RegisteredValidator(address treasury, address addr, string enode, address liquidContract)
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

// ParseRegisteredValidator is a log parse operation binding the contract event 0x6921859367aca5023ddf910758cd0cda74261b2e5c8425c253e9d03b62b950b8.
//
// Solidity: event RegisteredValidator(address treasury, address addr, string enode, address liquidContract)
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
	Bin: "0x60566037600b82828239805160001a607314602a57634e487b7160e01b600052600060045260246000fd5b30600052607381538281f3fe73000000000000000000000000000000000000000030146080604052600080fdfea264697066735822122014b844e19155670018f49d9420c48c85a5c41488dd53e1a3602e6a5d749f0d4164736f6c634300080b0033",
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
	Bin: "0x60566037600b82828239805160001a607314602a57634e487b7160e01b600052600060045260246000fd5b30600052607381538281f3fe73000000000000000000000000000000000000000030146080604052600080fdfea2646970667358221220595c45515cd78469fd85a5d9618ec2cf7da13afff35ffba5d9c101e2443fd4e064736f6c634300080b0033",
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

// LiquidMetaData contains all meta data concerning the Liquid contract.
var LiquidMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_validator\",\"type\":\"address\"},{\"internalType\":\"addresspayable\",\"name\":\"_treasury\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_commissionRate\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"COMMISSION_RATE_PRECISION\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"FEE_FACTOR_UNIT_RECIP\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_delegator\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"claimRewards\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"mint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"redistribute\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"_success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"_success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_account\",\"type\":\"address\"}],\"name\":\"unclaimedRewards\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"2f2c3f2e": "COMMISSION_RATE_PRECISION()",
		"187cf4d7": "FEE_FACTOR_UNIT_RECIP()",
		"dd62ed3e": "allowance(address,address)",
		"095ea7b3": "approve(address,uint256)",
		"70a08231": "balanceOf(address)",
		"9dc29fac": "burn(address,uint256)",
		"372500ab": "claimRewards()",
		"40c10f19": "mint(address,uint256)",
		"fb489a7b": "redistribute()",
		"18160ddd": "totalSupply()",
		"a9059cbb": "transfer(address,uint256)",
		"23b872dd": "transferFrom(address,address,uint256)",
		"949813b8": "unclaimedRewards(address)",
	},
	Bin: "0x608060405234801561001057600080fd5b50604051610cae380380610cae83398101604081905261002f9161009d565b61271081111561003e57600080fd5b600780546001600160a01b039485166001600160a01b03199182161790915560088054939094169281169290921790925560099190915560008054909116331790556100e0565b6001600160a01b038116811461009a57600080fd5b50565b6000806000606084860312156100b257600080fd5b83516100bd81610085565b60208501519093506100ce81610085565b80925050604084015190509250925092565b610bbf806100ef6000396000f3fe6080604052600436106100c25760003560e01c806340c10f191161007f5780639dc29fac116100595780639dc29fac146101f6578063a9059cbb14610216578063dd62ed3e14610236578063fb489a7b1461027c57600080fd5b806340c10f191461018057806370a08231146101a0578063949813b8146101d657600080fd5b8063095ea7b3146100c757806318160ddd146100fc578063187cf4d71461011b57806323b872dd146101335780632f2c3f2e14610153578063372500ab14610169575b600080fd5b3480156100d357600080fd5b506100e76100e23660046109e0565b610284565b60405190151581526020015b60405180910390f35b34801561010857600080fd5b506003545b6040519081526020016100f3565b34801561012757600080fd5b5061010d633b9aca0081565b34801561013f57600080fd5b506100e761014e366004610a0a565b61029a565b34801561015f57600080fd5b5061010d61271081565b34801561017557600080fd5b5061017e610392565b005b34801561018c57600080fd5b5061017e61019b3660046109e0565b610440565b3480156101ac57600080fd5b5061010d6101bb366004610a46565b6001600160a01b031660009081526001602052604090205490565b3480156101e257600080fd5b5061010d6101f1366004610a46565b6104a8565b34801561020257600080fd5b5061017e6102113660046109e0565b6104dc565b34801561022257600080fd5b506100e76102313660046109e0565b61053c565b34801561024257600080fd5b5061010d610251366004610a68565b6001600160a01b03918216600090815260026020908152604080832093909416825291909152205490565b61010d610589565b60006102913384846106d1565b50600192915050565b6001600160a01b0383166000908152600260209081526040808320338452909152812054828110156103245760405162461bcd60e51b815260206004820152602860248201527f45524332303a207472616e7366657220616d6f756e74206578636565647320616044820152676c6c6f77616e636560c01b60648201526084015b60405180910390fd5b61033885336103338685610ab1565b6106d1565b61034285846107f5565b61034c8484610898565b836001600160a01b0316856001600160a01b0316600080516020610b6a8339815191528560405161037f91815260200190565b60405180910390a3506001949350505050565b600061039d336108ec565b33600081815260046020526040808220829055519293509183908381818185875af1925050503d80600081146103ef576040519150601f19603f3d011682016040523d82523d6000602084013e6103f4565b606091505b505090508061043c5760405162461bcd60e51b81526020600482015260146024820152732330b4b632b2103a379039b2b7321022ba3432b960611b604482015260640161031b565b5050565b6000546001600160a01b0316331461046a5760405162461bcd60e51b815260040161031b90610ac8565b6104748282610898565b6040518181526001600160a01b03831690600090600080516020610b6a833981519152906020015b60405180910390a35050565b60006104b382610951565b6001600160a01b0383166000908152600460205260409020546104d69190610b10565b92915050565b6000546001600160a01b031633146105065760405162461bcd60e51b815260040161031b90610ac8565b61051082826107f5565b6040518181526000906001600160a01b03841690600080516020610b6a8339815191529060200161049c565b600061054833836107f5565b6105528383610898565b6040518281526001600160a01b038416903390600080516020610b6a8339815191529060200160405180910390a350600192915050565b600080546001600160a01b031633146105b45760405162461bcd60e51b815260040161031b90610ac8565b6009543490600090612710906105ca9084610b28565b6105d49190610b47565b90508181106106255760405162461bcd60e51b815260206004820152601860248201527f696e76616c69642076616c696461746f72207265776172640000000000000000604482015260640161031b565b61062f8183610ab1565b6008546040519193506001600160a01b03169082156108fc029083906000818181858888f1935050505015801561066a573d6000803e3d6000fd5b5060035460009061067f633b9aca0085610b28565b6106899190610b47565b9050806006546106999190610b10565b600655600354600090633b9aca00906106b29084610b28565b6106bc9190610b47565b90506106c88184610b10565b94505050505090565b6001600160a01b0383166107335760405162461bcd60e51b8152602060048201526024808201527f45524332303a20617070726f76652066726f6d20746865207a65726f206164646044820152637265737360e01b606482015260840161031b565b6001600160a01b0382166107945760405162461bcd60e51b815260206004820152602260248201527f45524332303a20617070726f766520746f20746865207a65726f206164647265604482015261737360f01b606482015260840161031b565b6001600160a01b0383811660008181526002602090815260408083209487168084529482529182902085905590518481527f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925910160405180910390a3505050565b6107fe826108ec565b506001600160a01b0382166000908152600160205260409020548082111561082557600080fd5b80821015610855576108378282610ab1565b6001600160a01b03841660009081526001602052604090205561087c565b6001600160a01b038316600090815260016020908152604080832083905560059091528120555b816003600082825461088e9190610ab1565b9091555050505050565b6108a1826108ec565b506001600160a01b038216600090815260016020526040812080548392906108ca908490610b10565b9250508190555080600360008282546108e39190610b10565b90915550505050565b6000806108f883610951565b6001600160a01b03841660009081526004602052604090205490915061091f908290610b10565b6001600160a01b0390931660009081526004602090815260408083208690556006546005909252909120555090919050565b6001600160a01b038116600090815260016020526040812054806109785750600092915050565b6001600160a01b03831660009081526005602052604081205460065461099e9190610ab1565b90506000633b9aca006109b18484610b28565b6109bb9190610b47565b95945050505050565b80356001600160a01b03811681146109db57600080fd5b919050565b600080604083850312156109f357600080fd5b6109fc836109c4565b946020939093013593505050565b600080600060608486031215610a1f57600080fd5b610a28846109c4565b9250610a36602085016109c4565b9150604084013590509250925092565b600060208284031215610a5857600080fd5b610a61826109c4565b9392505050565b60008060408385031215610a7b57600080fd5b610a84836109c4565b9150610a92602084016109c4565b90509250929050565b634e487b7160e01b600052601160045260246000fd5b600082821015610ac357610ac3610a9b565b500390565b60208082526028908201527f43616c6c207265737472696374656420746f20746865204175746f6e6974792060408201526710dbdb9d1c9858dd60c21b606082015260800190565b60008219821115610b2357610b23610a9b565b500190565b6000816000190483118215151615610b4257610b42610a9b565b500290565b600082610b6457634e487b7160e01b600052601260045260246000fd5b50049056feddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3efa2646970667358221220f9c95181bc1e859c0df200205d8951558d0b06d33a168b9f9240fc448aaea28064736f6c634300080b0033",
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
func DeployLiquid(auth *bind.TransactOpts, backend bind.ContractBackend, _validator common.Address, _treasury common.Address, _commissionRate *big.Int) (common.Address, *types.Transaction, *Liquid, error) {
	parsed, err := LiquidMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(LiquidBin), backend, _validator, _treasury, _commissionRate)
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

// PrecompiledMetaData contains all meta data concerning the Precompiled contract.
var PrecompiledMetaData = &bind.MetaData{
	ABI: "[]",
	Bin: "0x60566037600b82828239805160001a607314602a57634e487b7160e01b600052600060045260246000fd5b30600052607381538281f3fe73000000000000000000000000000000000000000030146080604052600080fdfea2646970667358221220fa3ed77bc986a365bb469fa3a76acc85e3eb5668655c67f068e7478d75d5415464736f6c634300080b0033",
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
