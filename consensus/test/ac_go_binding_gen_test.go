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
	SelfBondedStake   *big.Int
	TotalSlashed      *big.Int
	LiquidContract    common.Address
	LiquidSupply      *big.Int
	RegistrationBlock *big.Int
	State             uint8
}

// AutonityMetaData contains all meta data concerning the Autonity contract.
var AutonityMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"components\":[{\"internalType\":\"addresspayable\",\"name\":\"treasury\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"enode\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"commissionRate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"bondedStake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"selfBondedStake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalSlashed\",\"type\":\"uint256\"},{\"internalType\":\"contractLiquid\",\"name\":\"liquidContract\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"liquidSupply\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"registrationBlock\",\"type\":\"uint256\"},{\"internalType\":\"enumAutonity.ValidatorState\",\"name\":\"state\",\"type\":\"uint8\"}],\"internalType\":\"structAutonity.Validator[]\",\"name\":\"_validators\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"operatorAccount\",\"type\":\"address\"},{\"internalType\":\"addresspayable\",\"name\":\"treasuryAccount\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"treasuryFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"minBaseFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"delegationRate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"epochPeriod\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"unbondingPeriod\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"committeeSize\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"contractVersion\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"blockPeriod\",\"type\":\"uint256\"}],\"internalType\":\"structAutonity.Config\",\"name\":\"_config\",\"type\":\"tuple\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"BurnedStake\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasPrice\",\"type\":\"uint256\"}],\"name\":\"MinimumBaseFeeUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"MintedStake\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"treasury\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"effectiveBlock\",\"type\":\"uint256\"}],\"name\":\"PausedValidator\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"treasury\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"enode\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"liquidContract\",\"type\":\"address\"}],\"name\":\"RegisteredValidator\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Rewarded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"inputs\":[],\"name\":\"COMMISSION_RATE_PRECISION\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_address\",\"type\":\"address\"}],\"name\":\"activateValidator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_validator\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"bond\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"completeContractUpgrade\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"computeCommittee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"config\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"operatorAccount\",\"type\":\"address\"},{\"internalType\":\"addresspayable\",\"name\":\"treasuryAccount\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"treasuryFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"minBaseFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"delegationRate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"epochPeriod\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"unbondingPeriod\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"committeeSize\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"contractVersion\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"blockPeriod\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"deployer\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"epochID\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"epochReward\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"epochTotalBondedStake\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"finalize\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"votingPower\",\"type\":\"uint256\"}],\"internalType\":\"structAutonity.CommitteeMember[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"startId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"lastId\",\"type\":\"uint256\"}],\"name\":\"getBondingReq\",\"outputs\":[{\"components\":[{\"internalType\":\"addresspayable\",\"name\":\"delegator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"delegatee\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"startBlock\",\"type\":\"uint256\"}],\"internalType\":\"structAutonity.Staking[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCommittee\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"votingPower\",\"type\":\"uint256\"}],\"internalType\":\"structAutonity.CommitteeMember[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCommitteeEnodes\",\"outputs\":[{\"internalType\":\"string[]\",\"name\":\"\",\"type\":\"string[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLastEpochBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getMaxCommitteeSize\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getMinimumBaseFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getNewContract\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"},{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getOperator\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"height\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"round\",\"type\":\"uint256\"}],\"name\":\"getProposer\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"startId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"lastId\",\"type\":\"uint256\"}],\"name\":\"getUnbondingReq\",\"outputs\":[{\"components\":[{\"internalType\":\"addresspayable\",\"name\":\"delegator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"delegatee\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"startBlock\",\"type\":\"uint256\"}],\"internalType\":\"structAutonity.Staking[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"}],\"name\":\"getValidator\",\"outputs\":[{\"components\":[{\"internalType\":\"addresspayable\",\"name\":\"treasury\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"enode\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"commissionRate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"bondedStake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"selfBondedStake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalSlashed\",\"type\":\"uint256\"},{\"internalType\":\"contractLiquid\",\"name\":\"liquidContract\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"liquidSupply\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"registrationBlock\",\"type\":\"uint256\"},{\"internalType\":\"enumAutonity.ValidatorState\",\"name\":\"state\",\"type\":\"uint8\"}],\"internalType\":\"structAutonity.Validator\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getValidators\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getVersion\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"headBondingID\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"headUnbondingID\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastEpochBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"mint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_address\",\"type\":\"address\"}],\"name\":\"pauseValidator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_enode\",\"type\":\"string\"}],\"name\":\"registerValidator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"resetContractUpgrade\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_size\",\"type\":\"uint256\"}],\"name\":\"setCommitteeSize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_period\",\"type\":\"uint256\"}],\"name\":\"setEpochPeriod\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_price\",\"type\":\"uint256\"}],\"name\":\"setMinimumBaseFee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_account\",\"type\":\"address\"}],\"name\":\"setOperatorAccount\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"_account\",\"type\":\"address\"}],\"name\":\"setTreasuryAccount\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_treasuryFee\",\"type\":\"uint256\"}],\"name\":\"setTreasuryFee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_period\",\"type\":\"uint256\"}],\"name\":\"setUnbondingPeriod\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"tailBondingID\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"tailUnbondingID\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalRedistributed\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_validator\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"unbond\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"_bytecode\",\"type\":\"bytes\"},{\"internalType\":\"string\",\"name\":\"_abi\",\"type\":\"string\"}],\"name\":\"upgradeContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]=======autonity/solidity/contracts/Liquid.sol:Liquid=======Binary:60806040523480156200001157600080fd5b5060405162001a4538038062001a45833981810160405281019062000037919062000203565b6127108111156200004757600080fd5b82600760006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555081600860006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555080600981905550336000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055505050506200025f565b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b60006200014b826200011e565b9050919050565b6200015d816200013e565b81146200016957600080fd5b50565b6000815190506200017d8162000152565b92915050565b600062000190826200011e565b9050919050565b620001a28162000183565b8114620001ae57600080fd5b50565b600081519050620001c28162000197565b92915050565b6000819050919050565b620001dd81620001c8565b8114620001e957600080fd5b50565b600081519050620001fd81620001d2565b92915050565b6000806000606084860312156200021f576200021e62000119565b5b60006200022f868287016200016c565b93505060206200024286828701620001b1565b92505060406200025586828701620001ec565b9150509250925092565b6117d6806200026f6000396000f3fe6080604052600436106100c25760003560e01c806340c10f191161007f5780639dc29fac116100595780639dc29fac1461027c578063a9059cbb146102a5578063dd62ed3e146102e2578063fb489a7b1461031f576100c2565b806340c10f19146101d957806370a0823114610202578063949813b81461023f576100c2565b8063095ea7b3146100c757806318160ddd14610104578063187cf4d71461012f57806323b872dd1461015a5780632f2c3f2e14610197578063372500ab146101c2575b600080fd5b3480156100d357600080fd5b506100ee60048036038101906100e99190611156565b61033d565b6040516100fb91906111b1565b60405180910390f35b34801561011057600080fd5b50610119610354565b60405161012691906111db565b60405180910390f35b34801561013b57600080fd5b5061014461035e565b60405161015191906111db565b60405180910390f35b34801561016657600080fd5b50610181600480360381019061017c91906111f6565b610366565b60405161018e91906111b1565b60405180910390f35b3480156101a357600080fd5b506101ac6104c6565b6040516101b991906111db565b60405180910390f35b3480156101ce57600080fd5b506101d76104cc565b005b3480156101e557600080fd5b5061020060048036038101906101fb9190611156565b6105cc565b005b34801561020e57600080fd5b5061022960048036038101906102249190611249565b6106ce565b60405161023691906111db565b60405180910390f35b34801561024b57600080fd5b5061026660048036038101906102619190611249565b610717565b60405161027391906111db565b60405180910390f35b34801561028857600080fd5b506102a3600480360381019061029e9190611156565b610773565b005b3480156102b157600080fd5b506102cc60048036038101906102c79190611156565b610875565b6040516102d991906111b1565b60405180910390f35b3480156102ee57600080fd5b5061030960048036038101906103049190611276565b6108fa565b60405161031691906111db565b60405180910390f35b610327610981565b60405161033491906111db565b60405180910390f35b600061034a338484610b5a565b6001905092915050565b6000600354905090565b633b9aca0081565b600080600260008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205490508281101561042b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161042290611339565b60405180910390fd5b6104418533858461043c9190611388565b610b5a565b61044b8584610d25565b6104558484610e81565b8373ffffffffffffffffffffffffffffffffffffffff168573ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef856040516104b291906111db565b60405180910390a360019150509392505050565b61271081565b60006104d733610efe565b9050600460003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000905560003373ffffffffffffffffffffffffffffffffffffffff1682604051610542906113ed565b60006040518083038185875af1925050503d806000811461057f576040519150601f19603f3d011682016040523d82523d6000602084013e610584565b606091505b50509050806105c8576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016105bf9061144e565b60405180910390fd5b5050565b60008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161461065a576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610651906114e0565b60405180910390fd5b6106648282610e81565b8173ffffffffffffffffffffffffffffffffffffffff16600073ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef836040516106c291906111db565b60405180910390a35050565b6000600160008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020549050919050565b600061072282610fe9565b600460008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205461076c9190611500565b9050919050565b60008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614610801576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016107f8906114e0565b60405180910390fd5b61080b8282610d25565b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef8360405161086991906111db565b60405180910390a35050565b60006108813383610d25565b61088b8383610e81565b8273ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef846040516108e891906111db565b60405180910390a36001905092915050565b6000600260008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054905092915050565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614610a12576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610a09906114e0565b60405180910390fd5b6000349050600061271060095483610a2a9190611556565b610a3491906115df565b9050818110610a78576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610a6f9061165c565b60405180910390fd5b8082610a849190611388565b9150600860009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166108fc829081150290604051600060405180830381858888f19350505050158015610aee573d6000803e3d6000fd5b506000600354633b9aca0084610b049190611556565b610b0e91906115df565b905080600654610b1e9190611500565b6006819055506000633b9aca0060035483610b399190611556565b610b4391906115df565b90508083610b519190611500565b94505050505090565b600073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff161415610bca576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610bc1906116ee565b60405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff161415610c3a576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610c3190611780565b60405180910390fd5b80600260008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055508173ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b92583604051610d1891906111db565b60405180910390a3505050565b610d2e82610efe565b506000600160008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054905080821115610d8057600080fd5b80821015610ddc578181610d949190611388565b600160008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002081905550610e63565b600160008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009055600560008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600090555b8160036000828254610e759190611388565b92505081905550505050565b610e8a82610efe565b5080600160008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000828254610eda9190611500565b925050819055508060036000828254610ef39190611500565b925050819055505050565b600080610f0a83610fe9565b905080600460008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054610f579190611500565b915081600460008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002081905550600654600560008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000208190555050919050565b600080600160008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054905060008114156110415760009150506110b8565b6000600560008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020546006546110909190611388565b90506000633b9aca0083836110a59190611556565b6110af91906115df565b90508093505050505b919050565b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b60006110ed826110c2565b9050919050565b6110fd816110e2565b811461110857600080fd5b50565b60008135905061111a816110f4565b92915050565b6000819050919050565b61113381611120565b811461113e57600080fd5b50565b6000813590506111508161112a565b92915050565b6000806040838503121561116d5761116c6110bd565b5b600061117b8582860161110b565b925050602061118c85828601611141565b9150509250929050565b60008115159050919050565b6111ab81611196565b82525050565b60006020820190506111c660008301846111a2565b92915050565b6111d581611120565b82525050565b60006020820190506111f060008301846111cc565b92915050565b60008060006060848603121561120f5761120e6110bd565b5b600061121d8682870161110b565b935050602061122e8682870161110b565b925050604061123f86828701611141565b9150509250925092565b60006020828403121561125f5761125e6110bd565b5b600061126d8482850161110b565b91505092915050565b6000806040838503121561128d5761128c6110bd565b5b600061129b8582860161110b565b92505060206112ac8582860161110b565b9150509250929050565b600082825260208201905092915050565b7f45524332303a207472616e7366657220616d6f756e742065786365656473206160008201527f6c6c6f77616e6365000000000000000000000000000000000000000000000000602082015250565b60006113236028836112b6565b915061132e826112c7565b604082019050919050565b6000602082019050818103600083015261135281611316565b9050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b600061139382611120565b915061139e83611120565b9250828210156113b1576113b0611359565b5b828203905092915050565b600081905092915050565b50565b60006113d76000836113bc565b91506113e2826113c7565b600082019050919050565b60006113f8826113ca565b9150819050919050565b7f4661696c656420746f2073656e64204574686572000000000000000000000000600082015250565b60006114386014836112b6565b915061144382611402565b602082019050919050565b600060208201905081810360008301526114678161142b565b9050919050565b7f43616c6c207265737472696374656420746f20746865204175746f6e6974792060008201527f436f6e7472616374000000000000000000000000000000000000000000000000602082015250565b60006114ca6028836112b6565b91506114d58261146e565b604082019050919050565b600060208201905081810360008301526114f9816114bd565b9050919050565b600061150b82611120565b915061151683611120565b9250827fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0382111561154b5761154a611359565b5b828201905092915050565b600061156182611120565b915061156c83611120565b9250817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff04831182151516156115a5576115a4611359565b5b828202905092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b60006115ea82611120565b91506115f583611120565b925082611605576116046115b0565b5b828204905092915050565b7f696e76616c69642076616c696461746f72207265776172640000000000000000600082015250565b60006116466018836112b6565b915061165182611610565b602082019050919050565b6000602082019050818103600083015261167581611639565b9050919050565b7f45524332303a20617070726f76652066726f6d20746865207a65726f2061646460008201527f7265737300000000000000000000000000000000000000000000000000000000602082015250565b60006116d86024836112b6565b91506116e38261167c565b604082019050919050565b60006020820190508181036000830152611707816116cb565b9050919050565b7f45524332303a20617070726f766520746f20746865207a65726f20616464726560008201527f7373000000000000000000000000000000000000000000000000000000000000602082015250565b600061176a6022836112b6565b91506117758261170e565b604082019050919050565b600060208201905081810360008301526117998161175d565b905091905056fea26469706673582212208e86d53132513cc17541e03ceb564543376a3a09974b59a56103d20dd278ce2864736f6c634300080a0033ContractJSONABI[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_validator\",\"type\":\"address\"},{\"internalType\":\"addresspayable\",\"name\":\"_treasury\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_commissionRate\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"COMMISSION_RATE_PRECISION\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"FEE_FACTOR_UNIT_RECIP\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_delegator\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"claimRewards\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"mint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"redistribute\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"_success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"_success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_account\",\"type\":\"address\"}],\"name\":\"unclaimedRewards\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// AutonityABI is the input ABI used to generate the binding from.
// Deprecated: Use AutonityMetaData.ABI instead.
var AutonityABI = AutonityMetaData.ABI

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
// Solidity: function getValidator(address _addr) view returns((address,address,string,uint256,uint256,uint256,uint256,address,uint256,uint256,uint8))
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
// Solidity: function getValidator(address _addr) view returns((address,address,string,uint256,uint256,uint256,uint256,address,uint256,uint256,uint8))
func (_Autonity *AutonitySession) GetValidator(_addr common.Address) (AutonityValidator, error) {
	return _Autonity.Contract.GetValidator(&_Autonity.CallOpts, _addr)
}

// GetValidator is a free data retrieval call binding the contract method 0x1904bb2e.
//
// Solidity: function getValidator(address _addr) view returns((address,address,string,uint256,uint256,uint256,uint256,address,uint256,uint256,uint8))
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

// RegisterValidator is a paid mutator transaction binding the contract method 0xa416aa06.
//
// Solidity: function registerValidator(string _enode) returns()
func (_Autonity *AutonityTransactor) RegisterValidator(opts *bind.TransactOpts, _enode string) (*types.Transaction, error) {
	return _Autonity.contract.Transact(opts, "registerValidator", _enode)
}

// RegisterValidator is a paid mutator transaction binding the contract method 0xa416aa06.
//
// Solidity: function registerValidator(string _enode) returns()
func (_Autonity *AutonitySession) RegisterValidator(_enode string) (*types.Transaction, error) {
	return _Autonity.Contract.RegisterValidator(&_Autonity.TransactOpts, _enode)
}

// RegisterValidator is a paid mutator transaction binding the contract method 0xa416aa06.
//
// Solidity: function registerValidator(string _enode) returns()
func (_Autonity *AutonityTransactorSession) RegisterValidator(_enode string) (*types.Transaction, error) {
	return _Autonity.Contract.RegisterValidator(&_Autonity.TransactOpts, _enode)
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
