// Code generated for internal testing purposes only - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package tests1

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

// PrecompiledMetaData contains all meta data concerning the Precompiled contract.
var PrecompiledMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"ACCUSATION_CONTRACT\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"ACTIVITY_CONTRACT\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"COMPUTE_COMMITTEE_CONTRACT\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"ENODE_VERIFIER_CONTRACT\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"INNOCENCE_CONTRACT\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MISBEHAVIOUR_CONTRACT\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"POP_VERIFIER_CONTRACT\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"SUCCESS\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"UPGRADER_CONTRACT\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"4dc925d3": "ACCUSATION_CONTRACT()",
		"625fb940": "ACTIVITY_CONTRACT()",
		"2090a442": "COMPUTE_COMMITTEE_CONTRACT()",
		"c13974e1": "ENODE_VERIFIER_CONTRACT()",
		"8e153dc3": "INNOCENCE_CONTRACT()",
		"925c5492": "MISBEHAVIOUR_CONTRACT()",
		"50d93720": "POP_VERIFIER_CONTRACT()",
		"d0a6d1a6": "SUCCESS()",
		"a4ad5d91": "UPGRADER_CONTRACT()",
	},
	Bin: "0x610168610039600b82828239805160001a607314602c57634e487b7160e01b600052600060045260246000fd5b30600052607381538281f3fe73000000000000000000000000000000000000000030146080604052600436106100ad5760003560e01c80638e153dc311610080578063a4ad5d9111610065578063a4ad5d911461010c578063c13974e114610114578063d0a6d1a61461011c57600080fd5b80638e153dc3146100fc578063925c54921461010457600080fd5b80632090a442146100b25780634dc925d3146100e457806350d93720146100ec578063625fb940146100f4575b600080fd5b6100ba60fa81565b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020015b60405180910390f35b6100ba60fc81565b6100ba60fb81565b6100ba60f881565b6100ba60fd81565b6100ba60fe81565b6100ba60f981565b6100ba60ff81565b610124600181565b6040519081526020016100db56fea2646970667358221220ad1b147e0db6b5cc8b6ee04a14d20cf08a9e327bf70ab8857ba5df504e85a10864736f6c634300081e0033",
}

// PrecompiledABI is the input ABI used to generate the binding from.
// Deprecated: Use PrecompiledMetaData.ABI instead.
var PrecompiledABI = PrecompiledMetaData.ABI

// Deprecated: Use PrecompiledMetaData.Sigs instead.
// PrecompiledFuncSigs maps the 4-byte function signature to its string representation.
var PrecompiledFuncSigs = PrecompiledMetaData.Sigs

// PrecompiledBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use PrecompiledMetaData.Bin instead.
var PrecompiledBin = PrecompiledMetaData.Bin

// DeployPrecompiled deploys a new Ethereum contract, binding an instance of Precompiled to it.
func (r *Runner) DeployPrecompiled(opts *tests.RunOptions) (common.Address, uint64, *Precompiled, error) {
	parsed, err := PrecompiledMetaData.GetAbi()
	if err != nil {
		return common.Address{}, 0, nil, err
	}
	if parsed == nil {
		return common.Address{}, 0, nil, errors.New("GetABI returned nil")
	}

	address, gasConsumed, c, data, err := (*tests.Runner)(r).DeployContract(opts, parsed, common.FromHex(PrecompiledBin))
	if err != nil {
		return common.Address{}, 0, nil, (&Precompiled{Contract: c}).DecodeError(data, err)
	}
	return address, gasConsumed, &Precompiled{Contract: c}, nil
}

// Precompiled is an auto generated Go binding around an Ethereum contract.
type Precompiled struct {
	*tests.Contract
}

// ACCUSATIONCONTRACT is a free data retrieval call binding the contract method 0x4dc925d3.
//
// Solidity: function ACCUSATION_CONTRACT() view returns(address)
func (_Precompiled *Precompiled) ACCUSATIONCONTRACT(opts *tests.RunOptions) (common.Address, uint64, error) {
	data, consumed, err := _Precompiled.Call(opts, "ACCUSATION_CONTRACT")

	if err != nil {
		return *new(common.Address), consumed, _Precompiled.DecodeError(data, err)
	}
	out, err := _Precompiled.Contract.Abi().Unpack("ACCUSATION_CONTRACT", data)
	if err != nil {
		return *new(common.Address), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, consumed, err

}

// ACTIVITYCONTRACT is a free data retrieval call binding the contract method 0x625fb940.
//
// Solidity: function ACTIVITY_CONTRACT() view returns(address)
func (_Precompiled *Precompiled) ACTIVITYCONTRACT(opts *tests.RunOptions) (common.Address, uint64, error) {
	data, consumed, err := _Precompiled.Call(opts, "ACTIVITY_CONTRACT")

	if err != nil {
		return *new(common.Address), consumed, _Precompiled.DecodeError(data, err)
	}
	out, err := _Precompiled.Contract.Abi().Unpack("ACTIVITY_CONTRACT", data)
	if err != nil {
		return *new(common.Address), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, consumed, err

}

// COMPUTECOMMITTEECONTRACT is a free data retrieval call binding the contract method 0x2090a442.
//
// Solidity: function COMPUTE_COMMITTEE_CONTRACT() view returns(address)
func (_Precompiled *Precompiled) COMPUTECOMMITTEECONTRACT(opts *tests.RunOptions) (common.Address, uint64, error) {
	data, consumed, err := _Precompiled.Call(opts, "COMPUTE_COMMITTEE_CONTRACT")

	if err != nil {
		return *new(common.Address), consumed, _Precompiled.DecodeError(data, err)
	}
	out, err := _Precompiled.Contract.Abi().Unpack("COMPUTE_COMMITTEE_CONTRACT", data)
	if err != nil {
		return *new(common.Address), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, consumed, err

}

// ENODEVERIFIERCONTRACT is a free data retrieval call binding the contract method 0xc13974e1.
//
// Solidity: function ENODE_VERIFIER_CONTRACT() view returns(address)
func (_Precompiled *Precompiled) ENODEVERIFIERCONTRACT(opts *tests.RunOptions) (common.Address, uint64, error) {
	data, consumed, err := _Precompiled.Call(opts, "ENODE_VERIFIER_CONTRACT")

	if err != nil {
		return *new(common.Address), consumed, _Precompiled.DecodeError(data, err)
	}
	out, err := _Precompiled.Contract.Abi().Unpack("ENODE_VERIFIER_CONTRACT", data)
	if err != nil {
		return *new(common.Address), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, consumed, err

}

// INNOCENCECONTRACT is a free data retrieval call binding the contract method 0x8e153dc3.
//
// Solidity: function INNOCENCE_CONTRACT() view returns(address)
func (_Precompiled *Precompiled) INNOCENCECONTRACT(opts *tests.RunOptions) (common.Address, uint64, error) {
	data, consumed, err := _Precompiled.Call(opts, "INNOCENCE_CONTRACT")

	if err != nil {
		return *new(common.Address), consumed, _Precompiled.DecodeError(data, err)
	}
	out, err := _Precompiled.Contract.Abi().Unpack("INNOCENCE_CONTRACT", data)
	if err != nil {
		return *new(common.Address), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, consumed, err

}

// MISBEHAVIOURCONTRACT is a free data retrieval call binding the contract method 0x925c5492.
//
// Solidity: function MISBEHAVIOUR_CONTRACT() view returns(address)
func (_Precompiled *Precompiled) MISBEHAVIOURCONTRACT(opts *tests.RunOptions) (common.Address, uint64, error) {
	data, consumed, err := _Precompiled.Call(opts, "MISBEHAVIOUR_CONTRACT")

	if err != nil {
		return *new(common.Address), consumed, _Precompiled.DecodeError(data, err)
	}
	out, err := _Precompiled.Contract.Abi().Unpack("MISBEHAVIOUR_CONTRACT", data)
	if err != nil {
		return *new(common.Address), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, consumed, err

}

// POPVERIFIERCONTRACT is a free data retrieval call binding the contract method 0x50d93720.
//
// Solidity: function POP_VERIFIER_CONTRACT() view returns(address)
func (_Precompiled *Precompiled) POPVERIFIERCONTRACT(opts *tests.RunOptions) (common.Address, uint64, error) {
	data, consumed, err := _Precompiled.Call(opts, "POP_VERIFIER_CONTRACT")

	if err != nil {
		return *new(common.Address), consumed, _Precompiled.DecodeError(data, err)
	}
	out, err := _Precompiled.Contract.Abi().Unpack("POP_VERIFIER_CONTRACT", data)
	if err != nil {
		return *new(common.Address), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, consumed, err

}

// SUCCESS is a free data retrieval call binding the contract method 0xd0a6d1a6.
//
// Solidity: function SUCCESS() view returns(uint256)
func (_Precompiled *Precompiled) SUCCESS(opts *tests.RunOptions) (*big.Int, uint64, error) {
	data, consumed, err := _Precompiled.Call(opts, "SUCCESS")

	if err != nil {
		return *new(*big.Int), consumed, _Precompiled.DecodeError(data, err)
	}
	out, err := _Precompiled.Contract.Abi().Unpack("SUCCESS", data)
	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// UPGRADERCONTRACT is a free data retrieval call binding the contract method 0xa4ad5d91.
//
// Solidity: function UPGRADER_CONTRACT() view returns(address)
func (_Precompiled *Precompiled) UPGRADERCONTRACT(opts *tests.RunOptions) (common.Address, uint64, error) {
	data, consumed, err := _Precompiled.Call(opts, "UPGRADER_CONTRACT")

	if err != nil {
		return *new(common.Address), consumed, _Precompiled.DecodeError(data, err)
	}
	out, err := _Precompiled.Contract.Abi().Unpack("UPGRADER_CONTRACT", data)
	if err != nil {
		return *new(common.Address), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, consumed, err

}

func (_Precompiled *Precompiled) DecodeError(data []byte, err error) error {
	if err == nil {
		return nil
	}

	reason, _ := abi.UnpackRevert(data)
	return fmt.Errorf("%w: %s", err, reason)

}

// UpgradeManagerMetaData contains all meta data concerning the UpgradeManager contract.
var UpgradeManagerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_autonity\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_operator\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"oldValue\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newValue\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"appliesAtHeight\",\"type\":\"uint256\"}],\"name\":\"ConfigUpdateAddress\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"oldValue\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"newValue\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"appliesAtHeight\",\"type\":\"uint256\"}],\"name\":\"ConfigUpdateBool\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"int256\",\"name\":\"oldValue\",\"type\":\"int256\"},{\"indexed\":false,\"internalType\":\"int256\",\"name\":\"newValue\",\"type\":\"int256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"appliesAtHeight\",\"type\":\"uint256\"}],\"name\":\"ConfigUpdateInt\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldValue\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newValue\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"appliesAtHeight\",\"type\":\"uint256\"}],\"name\":\"ConfigUpdateUint\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"contractAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"name\":\"UpgradeResult\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"getAutonity\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getOperator\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_account\",\"type\":\"address\"}],\"name\":\"setOperator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_target\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"_data\",\"type\":\"string\"}],\"name\":\"upgrade\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"7c8ccebe": "getAutonity()",
		"e7f43c68": "getOperator()",
		"b3ab15fb": "setOperator(address)",
		"6e3d9ff0": "upgrade(address,string)",
	},
	Bin: "0x6080604052348015600f57600080fd5b50604051610664380380610664833981016040819052602c916077565b600080546001600160a01b039384166001600160a01b0319918216179091556001805492909316911617905560a5565b80516001600160a01b0381168114607257600080fd5b919050565b60008060408385031215608957600080fd5b609083605c565b9150609c60208401605c565b90509250929050565b6105b0806100b46000396000f3fe608060405234801561001057600080fd5b506004361061004c5760003560e01c80636e3d9ff0146100515780637c8ccebe14610066578063b3ab15fb146100a9578063e7f43c68146100bc575b600080fd5b61006461005f3660046103f3565b6100da565b005b60005473ffffffffffffffffffffffffffffffffffffffff165b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200160405180910390f35b6100646100b73660046104fa565b610220565b60015473ffffffffffffffffffffffffffffffffffffffff16610080565b60015473ffffffffffffffffffffffffffffffffffffffff163314610160576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601a60248201527f63616c6c6572206973206e6f7420746865206f70657261746f7200000000000060448201526064015b60405180910390fd5b60405160f990600090610179908590859060200161051c565b60405160208183030381529060405290506000806060600080855160208701885af4809350503d91506040519050602082018101604052818152816000602083013e8673ffffffffffffffffffffffffffffffffffffffff167f852faebd7b6261599c399849128319e84cd07abf61cfdea131405ee1d393479884604051610205911515815260200190565b60405180910390a282610219578160208201fd5b8160208201f35b60005473ffffffffffffffffffffffffffffffffffffffff1633146102c7576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602360248201527f63616c6c6572206973206e6f7420746865204175746f6e69747920636f6e747260448201527f61637400000000000000000000000000000000000000000000000000000000006064820152608401610157565b6001546040805160808082526008908201527f6f70657261746f7200000000000000000000000000000000000000000000000060a082015273ffffffffffffffffffffffffffffffffffffffff928316602082015291831682820152436060830152517fe821ac8084a7329d09d00cf1380cba50edea5f54a7bd453d49a9c354f565d4a09181900360c00190a1600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b803573ffffffffffffffffffffffffffffffffffffffff811681146103bf57600080fd5b919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6000806040838503121561040657600080fd5b61040f8361039b565b9150602083013567ffffffffffffffff81111561042b57600080fd5b8301601f8101851361043c57600080fd5b803567ffffffffffffffff811115610456576104566103c4565b6040517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0603f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8501160116810181811067ffffffffffffffff821117156104c2576104c26103c4565b6040528181528282016020018710156104da57600080fd5b816020840160208301376000602083830101528093505050509250929050565b60006020828403121561050c57600080fd5b6105158261039b565b9392505050565b7fffffffffffffffffffffffffffffffffffffffff0000000000000000000000008360601b1681526000825160005b81811015610568576020818601810151601486840101520161054b565b5060009201601401918252509291505056fea264697066735822122066eb1c1500246869fcc3df2d5d5c68c919de0b235533cd2a514f16c9f7459a7764736f6c634300081e0033",
}

// UpgradeManagerABI is the input ABI used to generate the binding from.
// Deprecated: Use UpgradeManagerMetaData.ABI instead.
var UpgradeManagerABI = UpgradeManagerMetaData.ABI

// Deprecated: Use UpgradeManagerMetaData.Sigs instead.
// UpgradeManagerFuncSigs maps the 4-byte function signature to its string representation.
var UpgradeManagerFuncSigs = UpgradeManagerMetaData.Sigs

// UpgradeManagerBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use UpgradeManagerMetaData.Bin instead.
var UpgradeManagerBin = UpgradeManagerMetaData.Bin

// DeployUpgradeManager deploys a new Ethereum contract, binding an instance of UpgradeManager to it.
func (r *Runner) DeployUpgradeManager(opts *tests.RunOptions, _autonity common.Address, _operator common.Address) (common.Address, uint64, *UpgradeManager, error) {
	parsed, err := UpgradeManagerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, 0, nil, err
	}
	if parsed == nil {
		return common.Address{}, 0, nil, errors.New("GetABI returned nil")
	}

	address, gasConsumed, c, data, err := (*tests.Runner)(r).DeployContract(opts, parsed, common.FromHex(UpgradeManagerBin), _autonity, _operator)
	if err != nil {
		return common.Address{}, 0, nil, (&UpgradeManager{Contract: c}).DecodeError(data, err)
	}
	return address, gasConsumed, &UpgradeManager{Contract: c}, nil
}

// UpgradeManager is an auto generated Go binding around an Ethereum contract.
type UpgradeManager struct {
	*tests.Contract
}

// GetAutonity is a free data retrieval call binding the contract method 0x7c8ccebe.
//
// Solidity: function getAutonity() view returns(address)
func (_UpgradeManager *UpgradeManager) GetAutonity(opts *tests.RunOptions) (common.Address, uint64, error) {
	data, consumed, err := _UpgradeManager.Call(opts, "getAutonity")

	if err != nil {
		return *new(common.Address), consumed, _UpgradeManager.DecodeError(data, err)
	}
	out, err := _UpgradeManager.Contract.Abi().Unpack("getAutonity", data)
	if err != nil {
		return *new(common.Address), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, consumed, err

}

// GetOperator is a free data retrieval call binding the contract method 0xe7f43c68.
//
// Solidity: function getOperator() view returns(address)
func (_UpgradeManager *UpgradeManager) GetOperator(opts *tests.RunOptions) (common.Address, uint64, error) {
	data, consumed, err := _UpgradeManager.Call(opts, "getOperator")

	if err != nil {
		return *new(common.Address), consumed, _UpgradeManager.DecodeError(data, err)
	}
	out, err := _UpgradeManager.Contract.Abi().Unpack("getOperator", data)
	if err != nil {
		return *new(common.Address), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, consumed, err

}

// SetOperator is a free data retrieval call for a paid mutator transaction binding the contract method 0xb3ab15fb.
// Similar to eth_call from rpc calls or function.call from truffle, it reverts the state after the call and returns the output. The output is extracted
// the same way as done above for view only functions.
// Solidity: function setOperator(address _account) returns()
func (_UpgradeManager *UpgradeManager) CallSetOperator(r *tests.Runner, opts *tests.RunOptions, _account common.Address) (uint64, error) {
	snap := r.Snapshot()

	data, consumed, err := _UpgradeManager.Call(opts, "setOperator", _account)
	r.RevertSnapshot(snap)
	return consumed, _UpgradeManager.DecodeError(data, err)

}

// Upgrade is a free data retrieval call for a paid mutator transaction binding the contract method 0x6e3d9ff0.
// Similar to eth_call from rpc calls or function.call from truffle, it reverts the state after the call and returns the output. The output is extracted
// the same way as done above for view only functions.
// Solidity: function upgrade(address _target, string _data) returns()
func (_UpgradeManager *UpgradeManager) CallUpgrade(r *tests.Runner, opts *tests.RunOptions, _target common.Address, _data string) (uint64, error) {
	snap := r.Snapshot()

	data, consumed, err := _UpgradeManager.Call(opts, "upgrade", _target, _data)
	r.RevertSnapshot(snap)
	return consumed, _UpgradeManager.DecodeError(data, err)

}

// SetOperator is a paid mutator transaction binding the contract method 0xb3ab15fb.
//
// Solidity: function setOperator(address _account) returns()
func (_UpgradeManager *UpgradeManager) SetOperator(opts *tests.RunOptions, _account common.Address) (uint64, error) {
	data, consumed, err := _UpgradeManager.Call(opts, "setOperator", _account)
	return consumed, _UpgradeManager.DecodeError(data, err)
}

// Upgrade is a paid mutator transaction binding the contract method 0x6e3d9ff0.
//
// Solidity: function upgrade(address _target, string _data) returns()
func (_UpgradeManager *UpgradeManager) Upgrade(opts *tests.RunOptions, _target common.Address, _data string) (uint64, error) {
	data, consumed, err := _UpgradeManager.Call(opts, "upgrade", _target, _data)
	return consumed, _UpgradeManager.DecodeError(data, err)
}

func (_UpgradeManager *UpgradeManager) DecodeError(data []byte, err error) error {
	if err == nil {
		return nil
	}

	reason, _ := abi.UnpackRevert(data)
	return fmt.Errorf("%w: %s", err, reason)

}

// UpgradeManager1MetaData contains all meta data concerning the UpgradeManager1 contract.
var UpgradeManager1MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_autonity\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_operator\",\"type\":\"address\"},{\"internalType\":\"bytes32[]\",\"name\":\"_hashes\",\"type\":\"bytes32[]\"},{\"internalType\":\"string[]\",\"name\":\"_versions\",\"type\":\"string[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"oldValue\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newValue\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"appliesAtHeight\",\"type\":\"uint256\"}],\"name\":\"ConfigUpdateAddress\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"oldValue\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"newValue\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"appliesAtHeight\",\"type\":\"uint256\"}],\"name\":\"ConfigUpdateBool\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"int256\",\"name\":\"oldValue\",\"type\":\"int256\"},{\"indexed\":false,\"internalType\":\"int256\",\"name\":\"newValue\",\"type\":\"int256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"appliesAtHeight\",\"type\":\"uint256\"}],\"name\":\"ConfigUpdateInt\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldValue\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newValue\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"appliesAtHeight\",\"type\":\"uint256\"}],\"name\":\"ConfigUpdateUint\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"contractAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"name\":\"UpgradeResult\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"getAutonity\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getOperator\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_hash\",\"type\":\"bytes32\"}],\"name\":\"getVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_account\",\"type\":\"address\"}],\"name\":\"setOperator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_hash\",\"type\":\"bytes32\"},{\"internalType\":\"string\",\"name\":\"_version\",\"type\":\"string\"}],\"name\":\"setVersion\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_target\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"_data\",\"type\":\"string\"}],\"name\":\"upgrade\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"7c8ccebe": "getAutonity()",
		"e7f43c68": "getOperator()",
		"9aaf9f08": "getVersion(bytes32)",
		"b3ab15fb": "setOperator(address)",
		"a1d20653": "setVersion(bytes32,string)",
		"6e3d9ff0": "upgrade(address,string)",
	},
	Bin: "0x608060405234801561001057600080fd5b50604051610edb380380610edb83398101604081905261002f916102c8565b600080546001600160a01b038087166001600160a01b031992831617909255600180549286169290911691909117905580518251146100c65760405162461bcd60e51b815260206004820152602960248201527f68617368657320616e642076657273696f6e73206861766520646966666572656044820152680dce840d8cadccee8d60bb1b606482015260840160405180910390fd5b60005b825181101561012d578181815181106100e4576100e46103ac565b602002602001015160026000858481518110610102576101026103ac565b602002602001015181526020019081526020016000209081610124919061044b565b506001016100c9565b5050505050610509565b80516001600160a01b038116811461014e57600080fd5b919050565b634e487b7160e01b600052604160045260246000fd5b604051601f8201601f191681016001600160401b038111828210171561019157610191610153565b604052919050565b60006001600160401b038211156101b2576101b2610153565b5060051b60200190565b600082601f8301126101cd57600080fd5b81516101e06101db82610199565b610169565b8082825260208201915060208360051b86010192508583111561020257600080fd5b602085015b838110156102be5780516001600160401b0381111561022557600080fd5b8601603f8101881361023657600080fd5b60208101516001600160401b0381111561025257610252610153565b610265601f8201601f1916602001610169565b8181526040838301018a101561027a57600080fd5b60005b8281101561029d578084016040015160208383018101919091520161027d565b50600060208383010152808652505050602083019250602081019050610207565b5095945050505050565b600080600080608085870312156102de57600080fd5b6102e785610137565b93506102f560208601610137565b60408601519093506001600160401b0381111561031157600080fd5b8501601f8101871361032257600080fd5b80516103306101db82610199565b8082825260208201915060208360051b85010192508983111561035257600080fd5b6020840193505b82841015610374578351825260209384019390910190610359565b6060890151909550925050506001600160401b0381111561039457600080fd5b6103a0878288016101bc565b91505092959194509250565b634e487b7160e01b600052603260045260246000fd5b600181811c908216806103d657607f821691505b6020821081036103f657634e487b7160e01b600052602260045260246000fd5b50919050565b601f82111561044657806000526020600020601f840160051c810160208510156104235750805b601f840160051c820191505b81811015610443576000815560010161042f565b50505b505050565b81516001600160401b0381111561046457610464610153565b6104788161047284546103c2565b846103fc565b6020601f8211600181146104ac57600083156104945750848201515b600019600385901b1c1916600184901b178455610443565b600084815260208120601f198516915b828110156104dc57878501518255602094850194600190920191016104bc565b50848210156104fa5786840151600019600387901b60f8161c191681555b50505050600190811b01905550565b6109c3806105186000396000f3fe608060405234801561001057600080fd5b50600436106100725760003560e01c8063a1d2065311610050578063a1d20653146100f0578063b3ab15fb14610103578063e7f43c681461011657600080fd5b80636e3d9ff0146100775780637c8ccebe1461008c5780639aaf9f08146100d0575b600080fd5b61008a610085366004610659565b610134565b005b60005473ffffffffffffffffffffffffffffffffffffffff165b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020015b60405180910390f35b6100e36100de3660046106a7565b61027a565b6040516100c791906106e4565b61008a6100fe366004610735565b61031c565b61008a610111366004610766565b6103ba565b60015473ffffffffffffffffffffffffffffffffffffffff166100a6565b60015473ffffffffffffffffffffffffffffffffffffffff1633146101ba576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601a60248201527f63616c6c6572206973206e6f7420746865206f70657261746f7200000000000060448201526064015b60405180910390fd5b60405160f9906000906101d39085908590602001610788565b60405160208183030381529060405290506000806060600080855160208701885af4809350503d91506040519050602082018101604052818152816000602083013e8673ffffffffffffffffffffffffffffffffffffffff167f852faebd7b6261599c399849128319e84cd07abf61cfdea131405ee1d39347988460405161025f911515815260200190565b60405180910390a282610273578160208201fd5b8160208201f35b6000818152600260205260409020805460609190610297906107d3565b80601f01602080910402602001604051908101604052809291908181526020018280546102c3906107d3565b80156103105780601f106102e557610100808354040283529160200191610310565b820191906000526020600020905b8154815290600101906020018083116102f357829003601f168201915b50505050509050919050565b60015473ffffffffffffffffffffffffffffffffffffffff16331461039d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601a60248201527f63616c6c6572206973206e6f7420746865206f70657261746f7200000000000060448201526064016101b1565b60008281526002602052604090206103b58282610874565b505050565b60005473ffffffffffffffffffffffffffffffffffffffff163314610461576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602360248201527f63616c6c6572206973206e6f7420746865204175746f6e69747920636f6e747260448201527f616374000000000000000000000000000000000000000000000000000000000060648201526084016101b1565b6001546040805160808082526008908201527f6f70657261746f7200000000000000000000000000000000000000000000000060a082015273ffffffffffffffffffffffffffffffffffffffff928316602082015291831682820152436060830152517fe821ac8084a7329d09d00cf1380cba50edea5f54a7bd453d49a9c354f565d4a09181900360c00190a1600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b803573ffffffffffffffffffffffffffffffffffffffff8116811461055957600080fd5b919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b600082601f83011261059e57600080fd5b813567ffffffffffffffff8111156105b8576105b861055e565b6040517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0603f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8501160116810181811067ffffffffffffffff821117156106245761062461055e565b60405281815283820160200185101561063c57600080fd5b816020850160208301376000918101602001919091529392505050565b6000806040838503121561066c57600080fd5b61067583610535565b9150602083013567ffffffffffffffff81111561069157600080fd5b61069d8582860161058d565b9150509250929050565b6000602082840312156106b957600080fd5b5035919050565b60005b838110156106db5781810151838201526020016106c3565b50506000910152565b60208152600082518060208401526107038160408501602087016106c0565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169190910160400192915050565b6000806040838503121561074857600080fd5b82359150602083013567ffffffffffffffff81111561069157600080fd5b60006020828403121561077857600080fd5b61078182610535565b9392505050565b7fffffffffffffffffffffffffffffffffffffffff0000000000000000000000008360601b168152600082516107c58160148501602087016106c0565b919091016014019392505050565b600181811c908216806107e757607f821691505b602082108103610820577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b601f8211156103b557806000526020600020601f840160051c8101602085101561084d5750805b601f840160051c820191505b8181101561086d5760008155600101610859565b5050505050565b815167ffffffffffffffff81111561088e5761088e61055e565b6108a28161089c84546107d3565b84610826565b6020601f8211600181146108f457600083156108be5750848201515b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600385901b1c1916600184901b17845561086d565b6000848152602081207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08516915b828110156109425787850151825560209485019460019092019101610922565b508482101561097e57868401517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600387901b60f8161c191681555b50505050600190811b0190555056fea26469706673582212207b90f4e2fd139cd2a2b80e5fae9d3fae8bf84fb8540995b491decd702dc2fe1664736f6c634300081e0033",
}

// UpgradeManager1ABI is the input ABI used to generate the binding from.
// Deprecated: Use UpgradeManager1MetaData.ABI instead.
var UpgradeManager1ABI = UpgradeManager1MetaData.ABI

// Deprecated: Use UpgradeManager1MetaData.Sigs instead.
// UpgradeManager1FuncSigs maps the 4-byte function signature to its string representation.
var UpgradeManager1FuncSigs = UpgradeManager1MetaData.Sigs

// UpgradeManager1Bin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use UpgradeManager1MetaData.Bin instead.
var UpgradeManager1Bin = UpgradeManager1MetaData.Bin

// DeployUpgradeManager1 deploys a new Ethereum contract, binding an instance of UpgradeManager1 to it.
func (r *Runner) DeployUpgradeManager1(opts *tests.RunOptions, _autonity common.Address, _operator common.Address, _hashes [][32]byte, _versions []string) (common.Address, uint64, *UpgradeManager1, error) {
	parsed, err := UpgradeManager1MetaData.GetAbi()
	if err != nil {
		return common.Address{}, 0, nil, err
	}
	if parsed == nil {
		return common.Address{}, 0, nil, errors.New("GetABI returned nil")
	}

	address, gasConsumed, c, data, err := (*tests.Runner)(r).DeployContract(opts, parsed, common.FromHex(UpgradeManager1Bin), _autonity, _operator, _hashes, _versions)
	if err != nil {
		return common.Address{}, 0, nil, (&UpgradeManager1{Contract: c}).DecodeError(data, err)
	}
	return address, gasConsumed, &UpgradeManager1{Contract: c}, nil
}

// UpgradeManager1 is an auto generated Go binding around an Ethereum contract.
type UpgradeManager1 struct {
	*tests.Contract
}

// GetAutonity is a free data retrieval call binding the contract method 0x7c8ccebe.
//
// Solidity: function getAutonity() view returns(address)
func (_UpgradeManager1 *UpgradeManager1) GetAutonity(opts *tests.RunOptions) (common.Address, uint64, error) {
	data, consumed, err := _UpgradeManager1.Call(opts, "getAutonity")

	if err != nil {
		return *new(common.Address), consumed, _UpgradeManager1.DecodeError(data, err)
	}
	out, err := _UpgradeManager1.Contract.Abi().Unpack("getAutonity", data)
	if err != nil {
		return *new(common.Address), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, consumed, err

}

// GetOperator is a free data retrieval call binding the contract method 0xe7f43c68.
//
// Solidity: function getOperator() view returns(address)
func (_UpgradeManager1 *UpgradeManager1) GetOperator(opts *tests.RunOptions) (common.Address, uint64, error) {
	data, consumed, err := _UpgradeManager1.Call(opts, "getOperator")

	if err != nil {
		return *new(common.Address), consumed, _UpgradeManager1.DecodeError(data, err)
	}
	out, err := _UpgradeManager1.Contract.Abi().Unpack("getOperator", data)
	if err != nil {
		return *new(common.Address), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, consumed, err

}

// GetVersion is a free data retrieval call binding the contract method 0x9aaf9f08.
//
// Solidity: function getVersion(bytes32 _hash) view returns(string)
func (_UpgradeManager1 *UpgradeManager1) GetVersion(opts *tests.RunOptions, _hash [32]byte) (string, uint64, error) {
	data, consumed, err := _UpgradeManager1.Call(opts, "getVersion", _hash)

	if err != nil {
		return *new(string), consumed, _UpgradeManager1.DecodeError(data, err)
	}
	out, err := _UpgradeManager1.Contract.Abi().Unpack("getVersion", data)
	if err != nil {
		return *new(string), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)
	return out0, consumed, err

}

// SetOperator is a free data retrieval call for a paid mutator transaction binding the contract method 0xb3ab15fb.
// Similar to eth_call from rpc calls or function.call from truffle, it reverts the state after the call and returns the output. The output is extracted
// the same way as done above for view only functions.
// Solidity: function setOperator(address _account) returns()
func (_UpgradeManager1 *UpgradeManager1) CallSetOperator(r *tests.Runner, opts *tests.RunOptions, _account common.Address) (uint64, error) {
	snap := r.Snapshot()

	data, consumed, err := _UpgradeManager1.Call(opts, "setOperator", _account)
	r.RevertSnapshot(snap)
	return consumed, _UpgradeManager1.DecodeError(data, err)

}

// SetVersion is a free data retrieval call for a paid mutator transaction binding the contract method 0xa1d20653.
// Similar to eth_call from rpc calls or function.call from truffle, it reverts the state after the call and returns the output. The output is extracted
// the same way as done above for view only functions.
// Solidity: function setVersion(bytes32 _hash, string _version) returns()
func (_UpgradeManager1 *UpgradeManager1) CallSetVersion(r *tests.Runner, opts *tests.RunOptions, _hash [32]byte, _version string) (uint64, error) {
	snap := r.Snapshot()

	data, consumed, err := _UpgradeManager1.Call(opts, "setVersion", _hash, _version)
	r.RevertSnapshot(snap)
	return consumed, _UpgradeManager1.DecodeError(data, err)

}

// Upgrade is a free data retrieval call for a paid mutator transaction binding the contract method 0x6e3d9ff0.
// Similar to eth_call from rpc calls or function.call from truffle, it reverts the state after the call and returns the output. The output is extracted
// the same way as done above for view only functions.
// Solidity: function upgrade(address _target, string _data) returns()
func (_UpgradeManager1 *UpgradeManager1) CallUpgrade(r *tests.Runner, opts *tests.RunOptions, _target common.Address, _data string) (uint64, error) {
	snap := r.Snapshot()

	data, consumed, err := _UpgradeManager1.Call(opts, "upgrade", _target, _data)
	r.RevertSnapshot(snap)
	return consumed, _UpgradeManager1.DecodeError(data, err)

}

// SetOperator is a paid mutator transaction binding the contract method 0xb3ab15fb.
//
// Solidity: function setOperator(address _account) returns()
func (_UpgradeManager1 *UpgradeManager1) SetOperator(opts *tests.RunOptions, _account common.Address) (uint64, error) {
	data, consumed, err := _UpgradeManager1.Call(opts, "setOperator", _account)
	return consumed, _UpgradeManager1.DecodeError(data, err)
}

// SetVersion is a paid mutator transaction binding the contract method 0xa1d20653.
//
// Solidity: function setVersion(bytes32 _hash, string _version) returns()
func (_UpgradeManager1 *UpgradeManager1) SetVersion(opts *tests.RunOptions, _hash [32]byte, _version string) (uint64, error) {
	data, consumed, err := _UpgradeManager1.Call(opts, "setVersion", _hash, _version)
	return consumed, _UpgradeManager1.DecodeError(data, err)
}

// Upgrade is a paid mutator transaction binding the contract method 0x6e3d9ff0.
//
// Solidity: function upgrade(address _target, string _data) returns()
func (_UpgradeManager1 *UpgradeManager1) Upgrade(opts *tests.RunOptions, _target common.Address, _data string) (uint64, error) {
	data, consumed, err := _UpgradeManager1.Call(opts, "upgrade", _target, _data)
	return consumed, _UpgradeManager1.DecodeError(data, err)
}

func (_UpgradeManager1 *UpgradeManager1) DecodeError(data []byte, err error) error {
	if err == nil {
		return nil
	}

	reason, _ := abi.UnpackRevert(data)
	return fmt.Errorf("%w: %s", err, reason)

}
