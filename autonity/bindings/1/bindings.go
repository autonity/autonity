// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package bindings1

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

// ACCUSATIONCONTRACT is a free data retrieval call binding the contract method 0x4dc925d3.
//
// Solidity: function ACCUSATION_CONTRACT() view returns(address)
func (_Precompiled *PrecompiledCaller) ACCUSATIONCONTRACT(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Precompiled.contract.Call(opts, &out, "ACCUSATION_CONTRACT")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ACCUSATIONCONTRACT is a free data retrieval call binding the contract method 0x4dc925d3.
//
// Solidity: function ACCUSATION_CONTRACT() view returns(address)
func (_Precompiled *PrecompiledSession) ACCUSATIONCONTRACT() (common.Address, error) {
	return _Precompiled.Contract.ACCUSATIONCONTRACT(&_Precompiled.CallOpts)
}

// ACCUSATIONCONTRACT is a free data retrieval call binding the contract method 0x4dc925d3.
//
// Solidity: function ACCUSATION_CONTRACT() view returns(address)
func (_Precompiled *PrecompiledCallerSession) ACCUSATIONCONTRACT() (common.Address, error) {
	return _Precompiled.Contract.ACCUSATIONCONTRACT(&_Precompiled.CallOpts)
}

// ACTIVITYCONTRACT is a free data retrieval call binding the contract method 0x625fb940.
//
// Solidity: function ACTIVITY_CONTRACT() view returns(address)
func (_Precompiled *PrecompiledCaller) ACTIVITYCONTRACT(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Precompiled.contract.Call(opts, &out, "ACTIVITY_CONTRACT")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ACTIVITYCONTRACT is a free data retrieval call binding the contract method 0x625fb940.
//
// Solidity: function ACTIVITY_CONTRACT() view returns(address)
func (_Precompiled *PrecompiledSession) ACTIVITYCONTRACT() (common.Address, error) {
	return _Precompiled.Contract.ACTIVITYCONTRACT(&_Precompiled.CallOpts)
}

// ACTIVITYCONTRACT is a free data retrieval call binding the contract method 0x625fb940.
//
// Solidity: function ACTIVITY_CONTRACT() view returns(address)
func (_Precompiled *PrecompiledCallerSession) ACTIVITYCONTRACT() (common.Address, error) {
	return _Precompiled.Contract.ACTIVITYCONTRACT(&_Precompiled.CallOpts)
}

// COMPUTECOMMITTEECONTRACT is a free data retrieval call binding the contract method 0x2090a442.
//
// Solidity: function COMPUTE_COMMITTEE_CONTRACT() view returns(address)
func (_Precompiled *PrecompiledCaller) COMPUTECOMMITTEECONTRACT(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Precompiled.contract.Call(opts, &out, "COMPUTE_COMMITTEE_CONTRACT")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// COMPUTECOMMITTEECONTRACT is a free data retrieval call binding the contract method 0x2090a442.
//
// Solidity: function COMPUTE_COMMITTEE_CONTRACT() view returns(address)
func (_Precompiled *PrecompiledSession) COMPUTECOMMITTEECONTRACT() (common.Address, error) {
	return _Precompiled.Contract.COMPUTECOMMITTEECONTRACT(&_Precompiled.CallOpts)
}

// COMPUTECOMMITTEECONTRACT is a free data retrieval call binding the contract method 0x2090a442.
//
// Solidity: function COMPUTE_COMMITTEE_CONTRACT() view returns(address)
func (_Precompiled *PrecompiledCallerSession) COMPUTECOMMITTEECONTRACT() (common.Address, error) {
	return _Precompiled.Contract.COMPUTECOMMITTEECONTRACT(&_Precompiled.CallOpts)
}

// ENODEVERIFIERCONTRACT is a free data retrieval call binding the contract method 0xc13974e1.
//
// Solidity: function ENODE_VERIFIER_CONTRACT() view returns(address)
func (_Precompiled *PrecompiledCaller) ENODEVERIFIERCONTRACT(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Precompiled.contract.Call(opts, &out, "ENODE_VERIFIER_CONTRACT")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ENODEVERIFIERCONTRACT is a free data retrieval call binding the contract method 0xc13974e1.
//
// Solidity: function ENODE_VERIFIER_CONTRACT() view returns(address)
func (_Precompiled *PrecompiledSession) ENODEVERIFIERCONTRACT() (common.Address, error) {
	return _Precompiled.Contract.ENODEVERIFIERCONTRACT(&_Precompiled.CallOpts)
}

// ENODEVERIFIERCONTRACT is a free data retrieval call binding the contract method 0xc13974e1.
//
// Solidity: function ENODE_VERIFIER_CONTRACT() view returns(address)
func (_Precompiled *PrecompiledCallerSession) ENODEVERIFIERCONTRACT() (common.Address, error) {
	return _Precompiled.Contract.ENODEVERIFIERCONTRACT(&_Precompiled.CallOpts)
}

// INNOCENCECONTRACT is a free data retrieval call binding the contract method 0x8e153dc3.
//
// Solidity: function INNOCENCE_CONTRACT() view returns(address)
func (_Precompiled *PrecompiledCaller) INNOCENCECONTRACT(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Precompiled.contract.Call(opts, &out, "INNOCENCE_CONTRACT")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// INNOCENCECONTRACT is a free data retrieval call binding the contract method 0x8e153dc3.
//
// Solidity: function INNOCENCE_CONTRACT() view returns(address)
func (_Precompiled *PrecompiledSession) INNOCENCECONTRACT() (common.Address, error) {
	return _Precompiled.Contract.INNOCENCECONTRACT(&_Precompiled.CallOpts)
}

// INNOCENCECONTRACT is a free data retrieval call binding the contract method 0x8e153dc3.
//
// Solidity: function INNOCENCE_CONTRACT() view returns(address)
func (_Precompiled *PrecompiledCallerSession) INNOCENCECONTRACT() (common.Address, error) {
	return _Precompiled.Contract.INNOCENCECONTRACT(&_Precompiled.CallOpts)
}

// MISBEHAVIOURCONTRACT is a free data retrieval call binding the contract method 0x925c5492.
//
// Solidity: function MISBEHAVIOUR_CONTRACT() view returns(address)
func (_Precompiled *PrecompiledCaller) MISBEHAVIOURCONTRACT(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Precompiled.contract.Call(opts, &out, "MISBEHAVIOUR_CONTRACT")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// MISBEHAVIOURCONTRACT is a free data retrieval call binding the contract method 0x925c5492.
//
// Solidity: function MISBEHAVIOUR_CONTRACT() view returns(address)
func (_Precompiled *PrecompiledSession) MISBEHAVIOURCONTRACT() (common.Address, error) {
	return _Precompiled.Contract.MISBEHAVIOURCONTRACT(&_Precompiled.CallOpts)
}

// MISBEHAVIOURCONTRACT is a free data retrieval call binding the contract method 0x925c5492.
//
// Solidity: function MISBEHAVIOUR_CONTRACT() view returns(address)
func (_Precompiled *PrecompiledCallerSession) MISBEHAVIOURCONTRACT() (common.Address, error) {
	return _Precompiled.Contract.MISBEHAVIOURCONTRACT(&_Precompiled.CallOpts)
}

// POPVERIFIERCONTRACT is a free data retrieval call binding the contract method 0x50d93720.
//
// Solidity: function POP_VERIFIER_CONTRACT() view returns(address)
func (_Precompiled *PrecompiledCaller) POPVERIFIERCONTRACT(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Precompiled.contract.Call(opts, &out, "POP_VERIFIER_CONTRACT")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// POPVERIFIERCONTRACT is a free data retrieval call binding the contract method 0x50d93720.
//
// Solidity: function POP_VERIFIER_CONTRACT() view returns(address)
func (_Precompiled *PrecompiledSession) POPVERIFIERCONTRACT() (common.Address, error) {
	return _Precompiled.Contract.POPVERIFIERCONTRACT(&_Precompiled.CallOpts)
}

// POPVERIFIERCONTRACT is a free data retrieval call binding the contract method 0x50d93720.
//
// Solidity: function POP_VERIFIER_CONTRACT() view returns(address)
func (_Precompiled *PrecompiledCallerSession) POPVERIFIERCONTRACT() (common.Address, error) {
	return _Precompiled.Contract.POPVERIFIERCONTRACT(&_Precompiled.CallOpts)
}

// SUCCESS is a free data retrieval call binding the contract method 0xd0a6d1a6.
//
// Solidity: function SUCCESS() view returns(uint256)
func (_Precompiled *PrecompiledCaller) SUCCESS(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Precompiled.contract.Call(opts, &out, "SUCCESS")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// SUCCESS is a free data retrieval call binding the contract method 0xd0a6d1a6.
//
// Solidity: function SUCCESS() view returns(uint256)
func (_Precompiled *PrecompiledSession) SUCCESS() (*big.Int, error) {
	return _Precompiled.Contract.SUCCESS(&_Precompiled.CallOpts)
}

// SUCCESS is a free data retrieval call binding the contract method 0xd0a6d1a6.
//
// Solidity: function SUCCESS() view returns(uint256)
func (_Precompiled *PrecompiledCallerSession) SUCCESS() (*big.Int, error) {
	return _Precompiled.Contract.SUCCESS(&_Precompiled.CallOpts)
}

// UPGRADERCONTRACT is a free data retrieval call binding the contract method 0xa4ad5d91.
//
// Solidity: function UPGRADER_CONTRACT() view returns(address)
func (_Precompiled *PrecompiledCaller) UPGRADERCONTRACT(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Precompiled.contract.Call(opts, &out, "UPGRADER_CONTRACT")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// UPGRADERCONTRACT is a free data retrieval call binding the contract method 0xa4ad5d91.
//
// Solidity: function UPGRADER_CONTRACT() view returns(address)
func (_Precompiled *PrecompiledSession) UPGRADERCONTRACT() (common.Address, error) {
	return _Precompiled.Contract.UPGRADERCONTRACT(&_Precompiled.CallOpts)
}

// UPGRADERCONTRACT is a free data retrieval call binding the contract method 0xa4ad5d91.
//
// Solidity: function UPGRADER_CONTRACT() view returns(address)
func (_Precompiled *PrecompiledCallerSession) UPGRADERCONTRACT() (common.Address, error) {
	return _Precompiled.Contract.UPGRADERCONTRACT(&_Precompiled.CallOpts)
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
func DeployUpgradeManager(auth *bind.TransactOpts, backend bind.ContractBackend, _autonity common.Address, _operator common.Address) (common.Address, *types.Transaction, *UpgradeManager, error) {
	parsed, err := UpgradeManagerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(UpgradeManagerBin), backend, _autonity, _operator)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &UpgradeManager{UpgradeManagerCaller: UpgradeManagerCaller{contract: contract}, UpgradeManagerTransactor: UpgradeManagerTransactor{contract: contract}, UpgradeManagerFilterer: UpgradeManagerFilterer{contract: contract}}, nil
}

// UpgradeManager is an auto generated Go binding around an Ethereum contract.
type UpgradeManager struct {
	UpgradeManagerCaller     // Read-only binding to the contract
	UpgradeManagerTransactor // Write-only binding to the contract
	UpgradeManagerFilterer   // Log filterer for contract events
}

// UpgradeManagerCaller is an auto generated read-only Go binding around an Ethereum contract.
type UpgradeManagerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UpgradeManagerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type UpgradeManagerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UpgradeManagerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type UpgradeManagerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UpgradeManagerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type UpgradeManagerSession struct {
	Contract     *UpgradeManager   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// UpgradeManagerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type UpgradeManagerCallerSession struct {
	Contract *UpgradeManagerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// UpgradeManagerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type UpgradeManagerTransactorSession struct {
	Contract     *UpgradeManagerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// UpgradeManagerRaw is an auto generated low-level Go binding around an Ethereum contract.
type UpgradeManagerRaw struct {
	Contract *UpgradeManager // Generic contract binding to access the raw methods on
}

// UpgradeManagerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type UpgradeManagerCallerRaw struct {
	Contract *UpgradeManagerCaller // Generic read-only contract binding to access the raw methods on
}

// UpgradeManagerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type UpgradeManagerTransactorRaw struct {
	Contract *UpgradeManagerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewUpgradeManager creates a new instance of UpgradeManager, bound to a specific deployed contract.
func NewUpgradeManager(address common.Address, backend bind.ContractBackend) (*UpgradeManager, error) {
	contract, err := bindUpgradeManager(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &UpgradeManager{UpgradeManagerCaller: UpgradeManagerCaller{contract: contract}, UpgradeManagerTransactor: UpgradeManagerTransactor{contract: contract}, UpgradeManagerFilterer: UpgradeManagerFilterer{contract: contract}}, nil
}

// NewUpgradeManagerCaller creates a new read-only instance of UpgradeManager, bound to a specific deployed contract.
func NewUpgradeManagerCaller(address common.Address, caller bind.ContractCaller) (*UpgradeManagerCaller, error) {
	contract, err := bindUpgradeManager(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &UpgradeManagerCaller{contract: contract}, nil
}

// NewUpgradeManagerTransactor creates a new write-only instance of UpgradeManager, bound to a specific deployed contract.
func NewUpgradeManagerTransactor(address common.Address, transactor bind.ContractTransactor) (*UpgradeManagerTransactor, error) {
	contract, err := bindUpgradeManager(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &UpgradeManagerTransactor{contract: contract}, nil
}

// NewUpgradeManagerFilterer creates a new log filterer instance of UpgradeManager, bound to a specific deployed contract.
func NewUpgradeManagerFilterer(address common.Address, filterer bind.ContractFilterer) (*UpgradeManagerFilterer, error) {
	contract, err := bindUpgradeManager(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &UpgradeManagerFilterer{contract: contract}, nil
}

// bindUpgradeManager binds a generic wrapper to an already deployed contract.
func bindUpgradeManager(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(UpgradeManagerABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_UpgradeManager *UpgradeManagerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _UpgradeManager.Contract.UpgradeManagerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_UpgradeManager *UpgradeManagerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UpgradeManager.Contract.UpgradeManagerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_UpgradeManager *UpgradeManagerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _UpgradeManager.Contract.UpgradeManagerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_UpgradeManager *UpgradeManagerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _UpgradeManager.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_UpgradeManager *UpgradeManagerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UpgradeManager.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_UpgradeManager *UpgradeManagerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _UpgradeManager.Contract.contract.Transact(opts, method, params...)
}

// GetAutonity is a free data retrieval call binding the contract method 0x7c8ccebe.
//
// Solidity: function getAutonity() view returns(address)
func (_UpgradeManager *UpgradeManagerCaller) GetAutonity(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _UpgradeManager.contract.Call(opts, &out, "getAutonity")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetAutonity is a free data retrieval call binding the contract method 0x7c8ccebe.
//
// Solidity: function getAutonity() view returns(address)
func (_UpgradeManager *UpgradeManagerSession) GetAutonity() (common.Address, error) {
	return _UpgradeManager.Contract.GetAutonity(&_UpgradeManager.CallOpts)
}

// GetAutonity is a free data retrieval call binding the contract method 0x7c8ccebe.
//
// Solidity: function getAutonity() view returns(address)
func (_UpgradeManager *UpgradeManagerCallerSession) GetAutonity() (common.Address, error) {
	return _UpgradeManager.Contract.GetAutonity(&_UpgradeManager.CallOpts)
}

// GetOperator is a free data retrieval call binding the contract method 0xe7f43c68.
//
// Solidity: function getOperator() view returns(address)
func (_UpgradeManager *UpgradeManagerCaller) GetOperator(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _UpgradeManager.contract.Call(opts, &out, "getOperator")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetOperator is a free data retrieval call binding the contract method 0xe7f43c68.
//
// Solidity: function getOperator() view returns(address)
func (_UpgradeManager *UpgradeManagerSession) GetOperator() (common.Address, error) {
	return _UpgradeManager.Contract.GetOperator(&_UpgradeManager.CallOpts)
}

// GetOperator is a free data retrieval call binding the contract method 0xe7f43c68.
//
// Solidity: function getOperator() view returns(address)
func (_UpgradeManager *UpgradeManagerCallerSession) GetOperator() (common.Address, error) {
	return _UpgradeManager.Contract.GetOperator(&_UpgradeManager.CallOpts)
}

// SetOperator is a paid mutator transaction binding the contract method 0xb3ab15fb.
//
// Solidity: function setOperator(address _account) returns()
func (_UpgradeManager *UpgradeManagerTransactor) SetOperator(opts *bind.TransactOpts, _account common.Address) (*types.Transaction, error) {
	return _UpgradeManager.contract.Transact(opts, "setOperator", _account)
}

// SetOperator is a paid mutator transaction binding the contract method 0xb3ab15fb.
//
// Solidity: function setOperator(address _account) returns()
func (_UpgradeManager *UpgradeManagerSession) SetOperator(_account common.Address) (*types.Transaction, error) {
	return _UpgradeManager.Contract.SetOperator(&_UpgradeManager.TransactOpts, _account)
}

// SetOperator is a paid mutator transaction binding the contract method 0xb3ab15fb.
//
// Solidity: function setOperator(address _account) returns()
func (_UpgradeManager *UpgradeManagerTransactorSession) SetOperator(_account common.Address) (*types.Transaction, error) {
	return _UpgradeManager.Contract.SetOperator(&_UpgradeManager.TransactOpts, _account)
}

// Upgrade is a paid mutator transaction binding the contract method 0x6e3d9ff0.
//
// Solidity: function upgrade(address _target, string _data) returns()
func (_UpgradeManager *UpgradeManagerTransactor) Upgrade(opts *bind.TransactOpts, _target common.Address, _data string) (*types.Transaction, error) {
	return _UpgradeManager.contract.Transact(opts, "upgrade", _target, _data)
}

// Upgrade is a paid mutator transaction binding the contract method 0x6e3d9ff0.
//
// Solidity: function upgrade(address _target, string _data) returns()
func (_UpgradeManager *UpgradeManagerSession) Upgrade(_target common.Address, _data string) (*types.Transaction, error) {
	return _UpgradeManager.Contract.Upgrade(&_UpgradeManager.TransactOpts, _target, _data)
}

// Upgrade is a paid mutator transaction binding the contract method 0x6e3d9ff0.
//
// Solidity: function upgrade(address _target, string _data) returns()
func (_UpgradeManager *UpgradeManagerTransactorSession) Upgrade(_target common.Address, _data string) (*types.Transaction, error) {
	return _UpgradeManager.Contract.Upgrade(&_UpgradeManager.TransactOpts, _target, _data)
}

// UpgradeManagerConfigUpdateAddressIterator is returned from FilterConfigUpdateAddress and is used to iterate over the raw logs and unpacked data for ConfigUpdateAddress events raised by the UpgradeManager contract.
type UpgradeManagerConfigUpdateAddressIterator struct {
	Event *UpgradeManagerConfigUpdateAddress // Event containing the contract specifics and raw log

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
func (it *UpgradeManagerConfigUpdateAddressIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UpgradeManagerConfigUpdateAddress)
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
		it.Event = new(UpgradeManagerConfigUpdateAddress)
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
func (it *UpgradeManagerConfigUpdateAddressIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UpgradeManagerConfigUpdateAddressIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UpgradeManagerConfigUpdateAddress represents a ConfigUpdateAddress event raised by the UpgradeManager contract.
type UpgradeManagerConfigUpdateAddress struct {
	Name            string
	OldValue        common.Address
	NewValue        common.Address
	AppliesAtHeight *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterConfigUpdateAddress is a free log retrieval operation binding the contract event 0xe821ac8084a7329d09d00cf1380cba50edea5f54a7bd453d49a9c354f565d4a0.
//
// Solidity: event ConfigUpdateAddress(string name, address oldValue, address newValue, uint256 appliesAtHeight)
func (_UpgradeManager *UpgradeManagerFilterer) FilterConfigUpdateAddress(opts *bind.FilterOpts) (*UpgradeManagerConfigUpdateAddressIterator, error) {

	logs, sub, err := _UpgradeManager.contract.FilterLogs(opts, "ConfigUpdateAddress")
	if err != nil {
		return nil, err
	}
	return &UpgradeManagerConfigUpdateAddressIterator{contract: _UpgradeManager.contract, event: "ConfigUpdateAddress", logs: logs, sub: sub}, nil
}

// WatchConfigUpdateAddress is a free log subscription operation binding the contract event 0xe821ac8084a7329d09d00cf1380cba50edea5f54a7bd453d49a9c354f565d4a0.
//
// Solidity: event ConfigUpdateAddress(string name, address oldValue, address newValue, uint256 appliesAtHeight)
func (_UpgradeManager *UpgradeManagerFilterer) WatchConfigUpdateAddress(opts *bind.WatchOpts, sink chan<- *UpgradeManagerConfigUpdateAddress) (event.Subscription, error) {

	logs, sub, err := _UpgradeManager.contract.WatchLogs(opts, "ConfigUpdateAddress")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UpgradeManagerConfigUpdateAddress)
				if err := _UpgradeManager.contract.UnpackLog(event, "ConfigUpdateAddress", log); err != nil {
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
func (_UpgradeManager *UpgradeManagerFilterer) ParseConfigUpdateAddress(log types.Log) (*UpgradeManagerConfigUpdateAddress, error) {
	event := new(UpgradeManagerConfigUpdateAddress)
	if err := _UpgradeManager.contract.UnpackLog(event, "ConfigUpdateAddress", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// UpgradeManagerConfigUpdateBoolIterator is returned from FilterConfigUpdateBool and is used to iterate over the raw logs and unpacked data for ConfigUpdateBool events raised by the UpgradeManager contract.
type UpgradeManagerConfigUpdateBoolIterator struct {
	Event *UpgradeManagerConfigUpdateBool // Event containing the contract specifics and raw log

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
func (it *UpgradeManagerConfigUpdateBoolIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UpgradeManagerConfigUpdateBool)
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
		it.Event = new(UpgradeManagerConfigUpdateBool)
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
func (it *UpgradeManagerConfigUpdateBoolIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UpgradeManagerConfigUpdateBoolIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UpgradeManagerConfigUpdateBool represents a ConfigUpdateBool event raised by the UpgradeManager contract.
type UpgradeManagerConfigUpdateBool struct {
	Name            string
	OldValue        bool
	NewValue        bool
	AppliesAtHeight *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterConfigUpdateBool is a free log retrieval operation binding the contract event 0x5edb308c5eddc69bcd31b4e689c5eed2fbd3155ae57915d1cad05425f6c1a39b.
//
// Solidity: event ConfigUpdateBool(string name, bool oldValue, bool newValue, uint256 appliesAtHeight)
func (_UpgradeManager *UpgradeManagerFilterer) FilterConfigUpdateBool(opts *bind.FilterOpts) (*UpgradeManagerConfigUpdateBoolIterator, error) {

	logs, sub, err := _UpgradeManager.contract.FilterLogs(opts, "ConfigUpdateBool")
	if err != nil {
		return nil, err
	}
	return &UpgradeManagerConfigUpdateBoolIterator{contract: _UpgradeManager.contract, event: "ConfigUpdateBool", logs: logs, sub: sub}, nil
}

// WatchConfigUpdateBool is a free log subscription operation binding the contract event 0x5edb308c5eddc69bcd31b4e689c5eed2fbd3155ae57915d1cad05425f6c1a39b.
//
// Solidity: event ConfigUpdateBool(string name, bool oldValue, bool newValue, uint256 appliesAtHeight)
func (_UpgradeManager *UpgradeManagerFilterer) WatchConfigUpdateBool(opts *bind.WatchOpts, sink chan<- *UpgradeManagerConfigUpdateBool) (event.Subscription, error) {

	logs, sub, err := _UpgradeManager.contract.WatchLogs(opts, "ConfigUpdateBool")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UpgradeManagerConfigUpdateBool)
				if err := _UpgradeManager.contract.UnpackLog(event, "ConfigUpdateBool", log); err != nil {
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
func (_UpgradeManager *UpgradeManagerFilterer) ParseConfigUpdateBool(log types.Log) (*UpgradeManagerConfigUpdateBool, error) {
	event := new(UpgradeManagerConfigUpdateBool)
	if err := _UpgradeManager.contract.UnpackLog(event, "ConfigUpdateBool", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// UpgradeManagerConfigUpdateIntIterator is returned from FilterConfigUpdateInt and is used to iterate over the raw logs and unpacked data for ConfigUpdateInt events raised by the UpgradeManager contract.
type UpgradeManagerConfigUpdateIntIterator struct {
	Event *UpgradeManagerConfigUpdateInt // Event containing the contract specifics and raw log

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
func (it *UpgradeManagerConfigUpdateIntIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UpgradeManagerConfigUpdateInt)
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
		it.Event = new(UpgradeManagerConfigUpdateInt)
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
func (it *UpgradeManagerConfigUpdateIntIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UpgradeManagerConfigUpdateIntIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UpgradeManagerConfigUpdateInt represents a ConfigUpdateInt event raised by the UpgradeManager contract.
type UpgradeManagerConfigUpdateInt struct {
	Name            string
	OldValue        *big.Int
	NewValue        *big.Int
	AppliesAtHeight *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterConfigUpdateInt is a free log retrieval operation binding the contract event 0xb5114472b89d1126433287defb6308bbabdb95b6ce5dd949ab2c151ed7b1ff4c.
//
// Solidity: event ConfigUpdateInt(string name, int256 oldValue, int256 newValue, uint256 appliesAtHeight)
func (_UpgradeManager *UpgradeManagerFilterer) FilterConfigUpdateInt(opts *bind.FilterOpts) (*UpgradeManagerConfigUpdateIntIterator, error) {

	logs, sub, err := _UpgradeManager.contract.FilterLogs(opts, "ConfigUpdateInt")
	if err != nil {
		return nil, err
	}
	return &UpgradeManagerConfigUpdateIntIterator{contract: _UpgradeManager.contract, event: "ConfigUpdateInt", logs: logs, sub: sub}, nil
}

// WatchConfigUpdateInt is a free log subscription operation binding the contract event 0xb5114472b89d1126433287defb6308bbabdb95b6ce5dd949ab2c151ed7b1ff4c.
//
// Solidity: event ConfigUpdateInt(string name, int256 oldValue, int256 newValue, uint256 appliesAtHeight)
func (_UpgradeManager *UpgradeManagerFilterer) WatchConfigUpdateInt(opts *bind.WatchOpts, sink chan<- *UpgradeManagerConfigUpdateInt) (event.Subscription, error) {

	logs, sub, err := _UpgradeManager.contract.WatchLogs(opts, "ConfigUpdateInt")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UpgradeManagerConfigUpdateInt)
				if err := _UpgradeManager.contract.UnpackLog(event, "ConfigUpdateInt", log); err != nil {
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
func (_UpgradeManager *UpgradeManagerFilterer) ParseConfigUpdateInt(log types.Log) (*UpgradeManagerConfigUpdateInt, error) {
	event := new(UpgradeManagerConfigUpdateInt)
	if err := _UpgradeManager.contract.UnpackLog(event, "ConfigUpdateInt", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// UpgradeManagerConfigUpdateUintIterator is returned from FilterConfigUpdateUint and is used to iterate over the raw logs and unpacked data for ConfigUpdateUint events raised by the UpgradeManager contract.
type UpgradeManagerConfigUpdateUintIterator struct {
	Event *UpgradeManagerConfigUpdateUint // Event containing the contract specifics and raw log

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
func (it *UpgradeManagerConfigUpdateUintIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UpgradeManagerConfigUpdateUint)
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
		it.Event = new(UpgradeManagerConfigUpdateUint)
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
func (it *UpgradeManagerConfigUpdateUintIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UpgradeManagerConfigUpdateUintIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UpgradeManagerConfigUpdateUint represents a ConfigUpdateUint event raised by the UpgradeManager contract.
type UpgradeManagerConfigUpdateUint struct {
	Name            string
	OldValue        *big.Int
	NewValue        *big.Int
	AppliesAtHeight *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterConfigUpdateUint is a free log retrieval operation binding the contract event 0x207e45ce6f2191c3efffe5d0d91591cf81afd23cbf4bd746981d7a8a8dcfe1ba.
//
// Solidity: event ConfigUpdateUint(string name, uint256 oldValue, uint256 newValue, uint256 appliesAtHeight)
func (_UpgradeManager *UpgradeManagerFilterer) FilterConfigUpdateUint(opts *bind.FilterOpts) (*UpgradeManagerConfigUpdateUintIterator, error) {

	logs, sub, err := _UpgradeManager.contract.FilterLogs(opts, "ConfigUpdateUint")
	if err != nil {
		return nil, err
	}
	return &UpgradeManagerConfigUpdateUintIterator{contract: _UpgradeManager.contract, event: "ConfigUpdateUint", logs: logs, sub: sub}, nil
}

// WatchConfigUpdateUint is a free log subscription operation binding the contract event 0x207e45ce6f2191c3efffe5d0d91591cf81afd23cbf4bd746981d7a8a8dcfe1ba.
//
// Solidity: event ConfigUpdateUint(string name, uint256 oldValue, uint256 newValue, uint256 appliesAtHeight)
func (_UpgradeManager *UpgradeManagerFilterer) WatchConfigUpdateUint(opts *bind.WatchOpts, sink chan<- *UpgradeManagerConfigUpdateUint) (event.Subscription, error) {

	logs, sub, err := _UpgradeManager.contract.WatchLogs(opts, "ConfigUpdateUint")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UpgradeManagerConfigUpdateUint)
				if err := _UpgradeManager.contract.UnpackLog(event, "ConfigUpdateUint", log); err != nil {
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
func (_UpgradeManager *UpgradeManagerFilterer) ParseConfigUpdateUint(log types.Log) (*UpgradeManagerConfigUpdateUint, error) {
	event := new(UpgradeManagerConfigUpdateUint)
	if err := _UpgradeManager.contract.UnpackLog(event, "ConfigUpdateUint", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// UpgradeManagerUpgradeResultIterator is returned from FilterUpgradeResult and is used to iterate over the raw logs and unpacked data for UpgradeResult events raised by the UpgradeManager contract.
type UpgradeManagerUpgradeResultIterator struct {
	Event *UpgradeManagerUpgradeResult // Event containing the contract specifics and raw log

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
func (it *UpgradeManagerUpgradeResultIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UpgradeManagerUpgradeResult)
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
		it.Event = new(UpgradeManagerUpgradeResult)
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
func (it *UpgradeManagerUpgradeResultIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UpgradeManagerUpgradeResultIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UpgradeManagerUpgradeResult represents a UpgradeResult event raised by the UpgradeManager contract.
type UpgradeManagerUpgradeResult struct {
	ContractAddress common.Address
	Success         bool
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterUpgradeResult is a free log retrieval operation binding the contract event 0x852faebd7b6261599c399849128319e84cd07abf61cfdea131405ee1d3934798.
//
// Solidity: event UpgradeResult(address indexed contractAddress, bool success)
func (_UpgradeManager *UpgradeManagerFilterer) FilterUpgradeResult(opts *bind.FilterOpts, contractAddress []common.Address) (*UpgradeManagerUpgradeResultIterator, error) {

	var contractAddressRule []interface{}
	for _, contractAddressItem := range contractAddress {
		contractAddressRule = append(contractAddressRule, contractAddressItem)
	}

	logs, sub, err := _UpgradeManager.contract.FilterLogs(opts, "UpgradeResult", contractAddressRule)
	if err != nil {
		return nil, err
	}
	return &UpgradeManagerUpgradeResultIterator{contract: _UpgradeManager.contract, event: "UpgradeResult", logs: logs, sub: sub}, nil
}

// WatchUpgradeResult is a free log subscription operation binding the contract event 0x852faebd7b6261599c399849128319e84cd07abf61cfdea131405ee1d3934798.
//
// Solidity: event UpgradeResult(address indexed contractAddress, bool success)
func (_UpgradeManager *UpgradeManagerFilterer) WatchUpgradeResult(opts *bind.WatchOpts, sink chan<- *UpgradeManagerUpgradeResult, contractAddress []common.Address) (event.Subscription, error) {

	var contractAddressRule []interface{}
	for _, contractAddressItem := range contractAddress {
		contractAddressRule = append(contractAddressRule, contractAddressItem)
	}

	logs, sub, err := _UpgradeManager.contract.WatchLogs(opts, "UpgradeResult", contractAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UpgradeManagerUpgradeResult)
				if err := _UpgradeManager.contract.UnpackLog(event, "UpgradeResult", log); err != nil {
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

// ParseUpgradeResult is a log parse operation binding the contract event 0x852faebd7b6261599c399849128319e84cd07abf61cfdea131405ee1d3934798.
//
// Solidity: event UpgradeResult(address indexed contractAddress, bool success)
func (_UpgradeManager *UpgradeManagerFilterer) ParseUpgradeResult(log types.Log) (*UpgradeManagerUpgradeResult, error) {
	event := new(UpgradeManagerUpgradeResult)
	if err := _UpgradeManager.contract.UnpackLog(event, "UpgradeResult", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
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
func DeployUpgradeManager1(auth *bind.TransactOpts, backend bind.ContractBackend, _autonity common.Address, _operator common.Address, _hashes [][32]byte, _versions []string) (common.Address, *types.Transaction, *UpgradeManager1, error) {
	parsed, err := UpgradeManager1MetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(UpgradeManager1Bin), backend, _autonity, _operator, _hashes, _versions)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &UpgradeManager1{UpgradeManager1Caller: UpgradeManager1Caller{contract: contract}, UpgradeManager1Transactor: UpgradeManager1Transactor{contract: contract}, UpgradeManager1Filterer: UpgradeManager1Filterer{contract: contract}}, nil
}

// UpgradeManager1 is an auto generated Go binding around an Ethereum contract.
type UpgradeManager1 struct {
	UpgradeManager1Caller     // Read-only binding to the contract
	UpgradeManager1Transactor // Write-only binding to the contract
	UpgradeManager1Filterer   // Log filterer for contract events
}

// UpgradeManager1Caller is an auto generated read-only Go binding around an Ethereum contract.
type UpgradeManager1Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UpgradeManager1Transactor is an auto generated write-only Go binding around an Ethereum contract.
type UpgradeManager1Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UpgradeManager1Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type UpgradeManager1Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UpgradeManager1Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type UpgradeManager1Session struct {
	Contract     *UpgradeManager1  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// UpgradeManager1CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type UpgradeManager1CallerSession struct {
	Contract *UpgradeManager1Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// UpgradeManager1TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type UpgradeManager1TransactorSession struct {
	Contract     *UpgradeManager1Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// UpgradeManager1Raw is an auto generated low-level Go binding around an Ethereum contract.
type UpgradeManager1Raw struct {
	Contract *UpgradeManager1 // Generic contract binding to access the raw methods on
}

// UpgradeManager1CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type UpgradeManager1CallerRaw struct {
	Contract *UpgradeManager1Caller // Generic read-only contract binding to access the raw methods on
}

// UpgradeManager1TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type UpgradeManager1TransactorRaw struct {
	Contract *UpgradeManager1Transactor // Generic write-only contract binding to access the raw methods on
}

// NewUpgradeManager1 creates a new instance of UpgradeManager1, bound to a specific deployed contract.
func NewUpgradeManager1(address common.Address, backend bind.ContractBackend) (*UpgradeManager1, error) {
	contract, err := bindUpgradeManager1(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &UpgradeManager1{UpgradeManager1Caller: UpgradeManager1Caller{contract: contract}, UpgradeManager1Transactor: UpgradeManager1Transactor{contract: contract}, UpgradeManager1Filterer: UpgradeManager1Filterer{contract: contract}}, nil
}

// NewUpgradeManager1Caller creates a new read-only instance of UpgradeManager1, bound to a specific deployed contract.
func NewUpgradeManager1Caller(address common.Address, caller bind.ContractCaller) (*UpgradeManager1Caller, error) {
	contract, err := bindUpgradeManager1(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &UpgradeManager1Caller{contract: contract}, nil
}

// NewUpgradeManager1Transactor creates a new write-only instance of UpgradeManager1, bound to a specific deployed contract.
func NewUpgradeManager1Transactor(address common.Address, transactor bind.ContractTransactor) (*UpgradeManager1Transactor, error) {
	contract, err := bindUpgradeManager1(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &UpgradeManager1Transactor{contract: contract}, nil
}

// NewUpgradeManager1Filterer creates a new log filterer instance of UpgradeManager1, bound to a specific deployed contract.
func NewUpgradeManager1Filterer(address common.Address, filterer bind.ContractFilterer) (*UpgradeManager1Filterer, error) {
	contract, err := bindUpgradeManager1(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &UpgradeManager1Filterer{contract: contract}, nil
}

// bindUpgradeManager1 binds a generic wrapper to an already deployed contract.
func bindUpgradeManager1(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(UpgradeManager1ABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_UpgradeManager1 *UpgradeManager1Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _UpgradeManager1.Contract.UpgradeManager1Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_UpgradeManager1 *UpgradeManager1Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UpgradeManager1.Contract.UpgradeManager1Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_UpgradeManager1 *UpgradeManager1Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _UpgradeManager1.Contract.UpgradeManager1Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_UpgradeManager1 *UpgradeManager1CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _UpgradeManager1.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_UpgradeManager1 *UpgradeManager1TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UpgradeManager1.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_UpgradeManager1 *UpgradeManager1TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _UpgradeManager1.Contract.contract.Transact(opts, method, params...)
}

// GetAutonity is a free data retrieval call binding the contract method 0x7c8ccebe.
//
// Solidity: function getAutonity() view returns(address)
func (_UpgradeManager1 *UpgradeManager1Caller) GetAutonity(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _UpgradeManager1.contract.Call(opts, &out, "getAutonity")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetAutonity is a free data retrieval call binding the contract method 0x7c8ccebe.
//
// Solidity: function getAutonity() view returns(address)
func (_UpgradeManager1 *UpgradeManager1Session) GetAutonity() (common.Address, error) {
	return _UpgradeManager1.Contract.GetAutonity(&_UpgradeManager1.CallOpts)
}

// GetAutonity is a free data retrieval call binding the contract method 0x7c8ccebe.
//
// Solidity: function getAutonity() view returns(address)
func (_UpgradeManager1 *UpgradeManager1CallerSession) GetAutonity() (common.Address, error) {
	return _UpgradeManager1.Contract.GetAutonity(&_UpgradeManager1.CallOpts)
}

// GetOperator is a free data retrieval call binding the contract method 0xe7f43c68.
//
// Solidity: function getOperator() view returns(address)
func (_UpgradeManager1 *UpgradeManager1Caller) GetOperator(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _UpgradeManager1.contract.Call(opts, &out, "getOperator")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetOperator is a free data retrieval call binding the contract method 0xe7f43c68.
//
// Solidity: function getOperator() view returns(address)
func (_UpgradeManager1 *UpgradeManager1Session) GetOperator() (common.Address, error) {
	return _UpgradeManager1.Contract.GetOperator(&_UpgradeManager1.CallOpts)
}

// GetOperator is a free data retrieval call binding the contract method 0xe7f43c68.
//
// Solidity: function getOperator() view returns(address)
func (_UpgradeManager1 *UpgradeManager1CallerSession) GetOperator() (common.Address, error) {
	return _UpgradeManager1.Contract.GetOperator(&_UpgradeManager1.CallOpts)
}

// GetVersion is a free data retrieval call binding the contract method 0x9aaf9f08.
//
// Solidity: function getVersion(bytes32 _hash) view returns(string)
func (_UpgradeManager1 *UpgradeManager1Caller) GetVersion(opts *bind.CallOpts, _hash [32]byte) (string, error) {
	var out []interface{}
	err := _UpgradeManager1.contract.Call(opts, &out, "getVersion", _hash)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// GetVersion is a free data retrieval call binding the contract method 0x9aaf9f08.
//
// Solidity: function getVersion(bytes32 _hash) view returns(string)
func (_UpgradeManager1 *UpgradeManager1Session) GetVersion(_hash [32]byte) (string, error) {
	return _UpgradeManager1.Contract.GetVersion(&_UpgradeManager1.CallOpts, _hash)
}

// GetVersion is a free data retrieval call binding the contract method 0x9aaf9f08.
//
// Solidity: function getVersion(bytes32 _hash) view returns(string)
func (_UpgradeManager1 *UpgradeManager1CallerSession) GetVersion(_hash [32]byte) (string, error) {
	return _UpgradeManager1.Contract.GetVersion(&_UpgradeManager1.CallOpts, _hash)
}

// SetOperator is a paid mutator transaction binding the contract method 0xb3ab15fb.
//
// Solidity: function setOperator(address _account) returns()
func (_UpgradeManager1 *UpgradeManager1Transactor) SetOperator(opts *bind.TransactOpts, _account common.Address) (*types.Transaction, error) {
	return _UpgradeManager1.contract.Transact(opts, "setOperator", _account)
}

// SetOperator is a paid mutator transaction binding the contract method 0xb3ab15fb.
//
// Solidity: function setOperator(address _account) returns()
func (_UpgradeManager1 *UpgradeManager1Session) SetOperator(_account common.Address) (*types.Transaction, error) {
	return _UpgradeManager1.Contract.SetOperator(&_UpgradeManager1.TransactOpts, _account)
}

// SetOperator is a paid mutator transaction binding the contract method 0xb3ab15fb.
//
// Solidity: function setOperator(address _account) returns()
func (_UpgradeManager1 *UpgradeManager1TransactorSession) SetOperator(_account common.Address) (*types.Transaction, error) {
	return _UpgradeManager1.Contract.SetOperator(&_UpgradeManager1.TransactOpts, _account)
}

// SetVersion is a paid mutator transaction binding the contract method 0xa1d20653.
//
// Solidity: function setVersion(bytes32 _hash, string _version) returns()
func (_UpgradeManager1 *UpgradeManager1Transactor) SetVersion(opts *bind.TransactOpts, _hash [32]byte, _version string) (*types.Transaction, error) {
	return _UpgradeManager1.contract.Transact(opts, "setVersion", _hash, _version)
}

// SetVersion is a paid mutator transaction binding the contract method 0xa1d20653.
//
// Solidity: function setVersion(bytes32 _hash, string _version) returns()
func (_UpgradeManager1 *UpgradeManager1Session) SetVersion(_hash [32]byte, _version string) (*types.Transaction, error) {
	return _UpgradeManager1.Contract.SetVersion(&_UpgradeManager1.TransactOpts, _hash, _version)
}

// SetVersion is a paid mutator transaction binding the contract method 0xa1d20653.
//
// Solidity: function setVersion(bytes32 _hash, string _version) returns()
func (_UpgradeManager1 *UpgradeManager1TransactorSession) SetVersion(_hash [32]byte, _version string) (*types.Transaction, error) {
	return _UpgradeManager1.Contract.SetVersion(&_UpgradeManager1.TransactOpts, _hash, _version)
}

// Upgrade is a paid mutator transaction binding the contract method 0x6e3d9ff0.
//
// Solidity: function upgrade(address _target, string _data) returns()
func (_UpgradeManager1 *UpgradeManager1Transactor) Upgrade(opts *bind.TransactOpts, _target common.Address, _data string) (*types.Transaction, error) {
	return _UpgradeManager1.contract.Transact(opts, "upgrade", _target, _data)
}

// Upgrade is a paid mutator transaction binding the contract method 0x6e3d9ff0.
//
// Solidity: function upgrade(address _target, string _data) returns()
func (_UpgradeManager1 *UpgradeManager1Session) Upgrade(_target common.Address, _data string) (*types.Transaction, error) {
	return _UpgradeManager1.Contract.Upgrade(&_UpgradeManager1.TransactOpts, _target, _data)
}

// Upgrade is a paid mutator transaction binding the contract method 0x6e3d9ff0.
//
// Solidity: function upgrade(address _target, string _data) returns()
func (_UpgradeManager1 *UpgradeManager1TransactorSession) Upgrade(_target common.Address, _data string) (*types.Transaction, error) {
	return _UpgradeManager1.Contract.Upgrade(&_UpgradeManager1.TransactOpts, _target, _data)
}

// UpgradeManager1ConfigUpdateAddressIterator is returned from FilterConfigUpdateAddress and is used to iterate over the raw logs and unpacked data for ConfigUpdateAddress events raised by the UpgradeManager1 contract.
type UpgradeManager1ConfigUpdateAddressIterator struct {
	Event *UpgradeManager1ConfigUpdateAddress // Event containing the contract specifics and raw log

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
func (it *UpgradeManager1ConfigUpdateAddressIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UpgradeManager1ConfigUpdateAddress)
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
		it.Event = new(UpgradeManager1ConfigUpdateAddress)
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
func (it *UpgradeManager1ConfigUpdateAddressIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UpgradeManager1ConfigUpdateAddressIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UpgradeManager1ConfigUpdateAddress represents a ConfigUpdateAddress event raised by the UpgradeManager1 contract.
type UpgradeManager1ConfigUpdateAddress struct {
	Name            string
	OldValue        common.Address
	NewValue        common.Address
	AppliesAtHeight *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterConfigUpdateAddress is a free log retrieval operation binding the contract event 0xe821ac8084a7329d09d00cf1380cba50edea5f54a7bd453d49a9c354f565d4a0.
//
// Solidity: event ConfigUpdateAddress(string name, address oldValue, address newValue, uint256 appliesAtHeight)
func (_UpgradeManager1 *UpgradeManager1Filterer) FilterConfigUpdateAddress(opts *bind.FilterOpts) (*UpgradeManager1ConfigUpdateAddressIterator, error) {

	logs, sub, err := _UpgradeManager1.contract.FilterLogs(opts, "ConfigUpdateAddress")
	if err != nil {
		return nil, err
	}
	return &UpgradeManager1ConfigUpdateAddressIterator{contract: _UpgradeManager1.contract, event: "ConfigUpdateAddress", logs: logs, sub: sub}, nil
}

// WatchConfigUpdateAddress is a free log subscription operation binding the contract event 0xe821ac8084a7329d09d00cf1380cba50edea5f54a7bd453d49a9c354f565d4a0.
//
// Solidity: event ConfigUpdateAddress(string name, address oldValue, address newValue, uint256 appliesAtHeight)
func (_UpgradeManager1 *UpgradeManager1Filterer) WatchConfigUpdateAddress(opts *bind.WatchOpts, sink chan<- *UpgradeManager1ConfigUpdateAddress) (event.Subscription, error) {

	logs, sub, err := _UpgradeManager1.contract.WatchLogs(opts, "ConfigUpdateAddress")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UpgradeManager1ConfigUpdateAddress)
				if err := _UpgradeManager1.contract.UnpackLog(event, "ConfigUpdateAddress", log); err != nil {
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
func (_UpgradeManager1 *UpgradeManager1Filterer) ParseConfigUpdateAddress(log types.Log) (*UpgradeManager1ConfigUpdateAddress, error) {
	event := new(UpgradeManager1ConfigUpdateAddress)
	if err := _UpgradeManager1.contract.UnpackLog(event, "ConfigUpdateAddress", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// UpgradeManager1ConfigUpdateBoolIterator is returned from FilterConfigUpdateBool and is used to iterate over the raw logs and unpacked data for ConfigUpdateBool events raised by the UpgradeManager1 contract.
type UpgradeManager1ConfigUpdateBoolIterator struct {
	Event *UpgradeManager1ConfigUpdateBool // Event containing the contract specifics and raw log

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
func (it *UpgradeManager1ConfigUpdateBoolIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UpgradeManager1ConfigUpdateBool)
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
		it.Event = new(UpgradeManager1ConfigUpdateBool)
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
func (it *UpgradeManager1ConfigUpdateBoolIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UpgradeManager1ConfigUpdateBoolIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UpgradeManager1ConfigUpdateBool represents a ConfigUpdateBool event raised by the UpgradeManager1 contract.
type UpgradeManager1ConfigUpdateBool struct {
	Name            string
	OldValue        bool
	NewValue        bool
	AppliesAtHeight *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterConfigUpdateBool is a free log retrieval operation binding the contract event 0x5edb308c5eddc69bcd31b4e689c5eed2fbd3155ae57915d1cad05425f6c1a39b.
//
// Solidity: event ConfigUpdateBool(string name, bool oldValue, bool newValue, uint256 appliesAtHeight)
func (_UpgradeManager1 *UpgradeManager1Filterer) FilterConfigUpdateBool(opts *bind.FilterOpts) (*UpgradeManager1ConfigUpdateBoolIterator, error) {

	logs, sub, err := _UpgradeManager1.contract.FilterLogs(opts, "ConfigUpdateBool")
	if err != nil {
		return nil, err
	}
	return &UpgradeManager1ConfigUpdateBoolIterator{contract: _UpgradeManager1.contract, event: "ConfigUpdateBool", logs: logs, sub: sub}, nil
}

// WatchConfigUpdateBool is a free log subscription operation binding the contract event 0x5edb308c5eddc69bcd31b4e689c5eed2fbd3155ae57915d1cad05425f6c1a39b.
//
// Solidity: event ConfigUpdateBool(string name, bool oldValue, bool newValue, uint256 appliesAtHeight)
func (_UpgradeManager1 *UpgradeManager1Filterer) WatchConfigUpdateBool(opts *bind.WatchOpts, sink chan<- *UpgradeManager1ConfigUpdateBool) (event.Subscription, error) {

	logs, sub, err := _UpgradeManager1.contract.WatchLogs(opts, "ConfigUpdateBool")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UpgradeManager1ConfigUpdateBool)
				if err := _UpgradeManager1.contract.UnpackLog(event, "ConfigUpdateBool", log); err != nil {
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
func (_UpgradeManager1 *UpgradeManager1Filterer) ParseConfigUpdateBool(log types.Log) (*UpgradeManager1ConfigUpdateBool, error) {
	event := new(UpgradeManager1ConfigUpdateBool)
	if err := _UpgradeManager1.contract.UnpackLog(event, "ConfigUpdateBool", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// UpgradeManager1ConfigUpdateIntIterator is returned from FilterConfigUpdateInt and is used to iterate over the raw logs and unpacked data for ConfigUpdateInt events raised by the UpgradeManager1 contract.
type UpgradeManager1ConfigUpdateIntIterator struct {
	Event *UpgradeManager1ConfigUpdateInt // Event containing the contract specifics and raw log

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
func (it *UpgradeManager1ConfigUpdateIntIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UpgradeManager1ConfigUpdateInt)
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
		it.Event = new(UpgradeManager1ConfigUpdateInt)
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
func (it *UpgradeManager1ConfigUpdateIntIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UpgradeManager1ConfigUpdateIntIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UpgradeManager1ConfigUpdateInt represents a ConfigUpdateInt event raised by the UpgradeManager1 contract.
type UpgradeManager1ConfigUpdateInt struct {
	Name            string
	OldValue        *big.Int
	NewValue        *big.Int
	AppliesAtHeight *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterConfigUpdateInt is a free log retrieval operation binding the contract event 0xb5114472b89d1126433287defb6308bbabdb95b6ce5dd949ab2c151ed7b1ff4c.
//
// Solidity: event ConfigUpdateInt(string name, int256 oldValue, int256 newValue, uint256 appliesAtHeight)
func (_UpgradeManager1 *UpgradeManager1Filterer) FilterConfigUpdateInt(opts *bind.FilterOpts) (*UpgradeManager1ConfigUpdateIntIterator, error) {

	logs, sub, err := _UpgradeManager1.contract.FilterLogs(opts, "ConfigUpdateInt")
	if err != nil {
		return nil, err
	}
	return &UpgradeManager1ConfigUpdateIntIterator{contract: _UpgradeManager1.contract, event: "ConfigUpdateInt", logs: logs, sub: sub}, nil
}

// WatchConfigUpdateInt is a free log subscription operation binding the contract event 0xb5114472b89d1126433287defb6308bbabdb95b6ce5dd949ab2c151ed7b1ff4c.
//
// Solidity: event ConfigUpdateInt(string name, int256 oldValue, int256 newValue, uint256 appliesAtHeight)
func (_UpgradeManager1 *UpgradeManager1Filterer) WatchConfigUpdateInt(opts *bind.WatchOpts, sink chan<- *UpgradeManager1ConfigUpdateInt) (event.Subscription, error) {

	logs, sub, err := _UpgradeManager1.contract.WatchLogs(opts, "ConfigUpdateInt")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UpgradeManager1ConfigUpdateInt)
				if err := _UpgradeManager1.contract.UnpackLog(event, "ConfigUpdateInt", log); err != nil {
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
func (_UpgradeManager1 *UpgradeManager1Filterer) ParseConfigUpdateInt(log types.Log) (*UpgradeManager1ConfigUpdateInt, error) {
	event := new(UpgradeManager1ConfigUpdateInt)
	if err := _UpgradeManager1.contract.UnpackLog(event, "ConfigUpdateInt", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// UpgradeManager1ConfigUpdateUintIterator is returned from FilterConfigUpdateUint and is used to iterate over the raw logs and unpacked data for ConfigUpdateUint events raised by the UpgradeManager1 contract.
type UpgradeManager1ConfigUpdateUintIterator struct {
	Event *UpgradeManager1ConfigUpdateUint // Event containing the contract specifics and raw log

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
func (it *UpgradeManager1ConfigUpdateUintIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UpgradeManager1ConfigUpdateUint)
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
		it.Event = new(UpgradeManager1ConfigUpdateUint)
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
func (it *UpgradeManager1ConfigUpdateUintIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UpgradeManager1ConfigUpdateUintIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UpgradeManager1ConfigUpdateUint represents a ConfigUpdateUint event raised by the UpgradeManager1 contract.
type UpgradeManager1ConfigUpdateUint struct {
	Name            string
	OldValue        *big.Int
	NewValue        *big.Int
	AppliesAtHeight *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterConfigUpdateUint is a free log retrieval operation binding the contract event 0x207e45ce6f2191c3efffe5d0d91591cf81afd23cbf4bd746981d7a8a8dcfe1ba.
//
// Solidity: event ConfigUpdateUint(string name, uint256 oldValue, uint256 newValue, uint256 appliesAtHeight)
func (_UpgradeManager1 *UpgradeManager1Filterer) FilterConfigUpdateUint(opts *bind.FilterOpts) (*UpgradeManager1ConfigUpdateUintIterator, error) {

	logs, sub, err := _UpgradeManager1.contract.FilterLogs(opts, "ConfigUpdateUint")
	if err != nil {
		return nil, err
	}
	return &UpgradeManager1ConfigUpdateUintIterator{contract: _UpgradeManager1.contract, event: "ConfigUpdateUint", logs: logs, sub: sub}, nil
}

// WatchConfigUpdateUint is a free log subscription operation binding the contract event 0x207e45ce6f2191c3efffe5d0d91591cf81afd23cbf4bd746981d7a8a8dcfe1ba.
//
// Solidity: event ConfigUpdateUint(string name, uint256 oldValue, uint256 newValue, uint256 appliesAtHeight)
func (_UpgradeManager1 *UpgradeManager1Filterer) WatchConfigUpdateUint(opts *bind.WatchOpts, sink chan<- *UpgradeManager1ConfigUpdateUint) (event.Subscription, error) {

	logs, sub, err := _UpgradeManager1.contract.WatchLogs(opts, "ConfigUpdateUint")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UpgradeManager1ConfigUpdateUint)
				if err := _UpgradeManager1.contract.UnpackLog(event, "ConfigUpdateUint", log); err != nil {
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
func (_UpgradeManager1 *UpgradeManager1Filterer) ParseConfigUpdateUint(log types.Log) (*UpgradeManager1ConfigUpdateUint, error) {
	event := new(UpgradeManager1ConfigUpdateUint)
	if err := _UpgradeManager1.contract.UnpackLog(event, "ConfigUpdateUint", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// UpgradeManager1UpgradeResultIterator is returned from FilterUpgradeResult and is used to iterate over the raw logs and unpacked data for UpgradeResult events raised by the UpgradeManager1 contract.
type UpgradeManager1UpgradeResultIterator struct {
	Event *UpgradeManager1UpgradeResult // Event containing the contract specifics and raw log

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
func (it *UpgradeManager1UpgradeResultIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UpgradeManager1UpgradeResult)
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
		it.Event = new(UpgradeManager1UpgradeResult)
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
func (it *UpgradeManager1UpgradeResultIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UpgradeManager1UpgradeResultIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UpgradeManager1UpgradeResult represents a UpgradeResult event raised by the UpgradeManager1 contract.
type UpgradeManager1UpgradeResult struct {
	ContractAddress common.Address
	Success         bool
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterUpgradeResult is a free log retrieval operation binding the contract event 0x852faebd7b6261599c399849128319e84cd07abf61cfdea131405ee1d3934798.
//
// Solidity: event UpgradeResult(address indexed contractAddress, bool success)
func (_UpgradeManager1 *UpgradeManager1Filterer) FilterUpgradeResult(opts *bind.FilterOpts, contractAddress []common.Address) (*UpgradeManager1UpgradeResultIterator, error) {

	var contractAddressRule []interface{}
	for _, contractAddressItem := range contractAddress {
		contractAddressRule = append(contractAddressRule, contractAddressItem)
	}

	logs, sub, err := _UpgradeManager1.contract.FilterLogs(opts, "UpgradeResult", contractAddressRule)
	if err != nil {
		return nil, err
	}
	return &UpgradeManager1UpgradeResultIterator{contract: _UpgradeManager1.contract, event: "UpgradeResult", logs: logs, sub: sub}, nil
}

// WatchUpgradeResult is a free log subscription operation binding the contract event 0x852faebd7b6261599c399849128319e84cd07abf61cfdea131405ee1d3934798.
//
// Solidity: event UpgradeResult(address indexed contractAddress, bool success)
func (_UpgradeManager1 *UpgradeManager1Filterer) WatchUpgradeResult(opts *bind.WatchOpts, sink chan<- *UpgradeManager1UpgradeResult, contractAddress []common.Address) (event.Subscription, error) {

	var contractAddressRule []interface{}
	for _, contractAddressItem := range contractAddress {
		contractAddressRule = append(contractAddressRule, contractAddressItem)
	}

	logs, sub, err := _UpgradeManager1.contract.WatchLogs(opts, "UpgradeResult", contractAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UpgradeManager1UpgradeResult)
				if err := _UpgradeManager1.contract.UnpackLog(event, "UpgradeResult", log); err != nil {
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

// ParseUpgradeResult is a log parse operation binding the contract event 0x852faebd7b6261599c399849128319e84cd07abf61cfdea131405ee1d3934798.
//
// Solidity: event UpgradeResult(address indexed contractAddress, bool success)
func (_UpgradeManager1 *UpgradeManager1Filterer) ParseUpgradeResult(log types.Log) (*UpgradeManager1UpgradeResult, error) {
	event := new(UpgradeManager1UpgradeResult)
	if err := _UpgradeManager1.contract.UnpackLog(event, "UpgradeResult", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
