// // Code generated - DO NOT EDIT.
// // This file is a generated binding and any manual changes will be lost.

package contract

// import (
// 	"errors"
// 	"math/big"
// 	"strings"

// 	ethereum "github.com/ethereum/go-ethereum"
// 	"github.com/ethereum/go-ethereum/accounts/abi"
// 	"github.com/ethereum/go-ethereum/accounts/abi/bind"
// 	"github.com/ethereum/go-ethereum/common"
// 	"github.com/ethereum/go-ethereum/core/types"
// 	"github.com/ethereum/go-ethereum/event"
// )

// // Reference imports to suppress errors if they are not otherwise used.
// var (
// 	_ = errors.New
// 	_ = big.NewInt
// 	_ = strings.NewReader
// 	_ = ethereum.NotFound
// 	_ = bind.Bind
// 	_ = common.Big1
// 	_ = types.BloomLookup
// 	_ = event.NewSubscription
// 	_ = abi.ConvertType
// )

// // ContractMetaData contains all meta data concerning the Contract contract.
// var ContractMetaData = &bind.MetaData{
// 	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"cid\",\"type\":\"string\"}],\"name\":\"CIDSaved\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"key\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"cid\",\"type\":\"string\"}],\"name\":\"saveCID\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"key\",\"type\":\"string\"}],\"name\":\"getCID\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
// }

// // ContractABI is the input ABI used to generate the binding from.
// // Deprecated: Use ContractMetaData.ABI instead.
// var ContractABI = ContractMetaData.ABI

// // Contract is an auto generated Go binding around an Ethereum contract.
// type Contract struct {
// 	ContractCaller     // Read-only binding to the contract
// 	ContractTransactor // Write-only binding to the contract
// 	ContractFilterer   // Log filterer for contract events
// }

// // ContractCaller is an auto generated read-only Go binding around an Ethereum contract.
// type ContractCaller struct {
// 	contract *bind.BoundContract // Generic contract wrapper for the low level calls
// }

// // ContractTransactor is an auto generated write-only Go binding around an Ethereum contract.
// type ContractTransactor struct {
// 	contract *bind.BoundContract // Generic contract wrapper for the low level calls
// }

// // ContractFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
// type ContractFilterer struct {
// 	contract *bind.BoundContract // Generic contract wrapper for the low level calls
// }

// // ContractSession is an auto generated Go binding around an Ethereum contract,
// // with pre-set call and transact options.
// type ContractSession struct {
// 	Contract     *Contract         // Generic contract binding to set the session for
// 	CallOpts     bind.CallOpts     // Call options to use throughout this session
// 	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
// }

// // ContractCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// // with pre-set call options.
// type ContractCallerSession struct {
// 	Contract *ContractCaller // Generic contract caller binding to set the session for
// 	CallOpts bind.CallOpts   // Call options to use throughout this session
// }

// // ContractTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// // with pre-set transact options.
// type ContractTransactorSession struct {
// 	Contract     *ContractTransactor // Generic contract transactor binding to set the session for
// 	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
// }

// // ContractRaw is an auto generated low-level Go binding around an Ethereum contract.
// type ContractRaw struct {
// 	Contract *Contract // Generic contract binding to access the raw methods on
// }

// // ContractCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
// type ContractCallerRaw struct {
// 	Contract *ContractCaller // Generic read-only contract binding to access the raw methods on
// }

// // ContractTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
// type ContractTransactorRaw struct {
// 	Contract *ContractTransactor // Generic write-only contract binding to access the raw methods on
// }

// // NewContract creates a new instance of Contract, bound to a specific deployed contract.
// func NewContract(address common.Address, backend bind.ContractBackend) (*Contract, error) {
// 	contract, err := bindContract(address, backend, backend, backend)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &Contract{ContractCaller: ContractCaller{contract: contract}, ContractTransactor: ContractTransactor{contract: contract}, ContractFilterer: ContractFilterer{contract: contract}}, nil
// }

// // NewContractCaller creates a new read-only instance of Contract, bound to a specific deployed contract.
// func NewContractCaller(address common.Address, caller bind.ContractCaller) (*ContractCaller, error) {
// 	contract, err := bindContract(address, caller, nil, nil)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &ContractCaller{contract: contract}, nil
// }

// // NewContractTransactor creates a new write-only instance of Contract, bound to a specific deployed contract.
// func NewContractTransactor(address common.Address, transactor bind.ContractTransactor) (*ContractTransactor, error) {
// 	contract, err := bindContract(address, nil, transactor, nil)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &ContractTransactor{contract: contract}, nil
// }

// // NewContractFilterer creates a new log filterer instance of Contract, bound to a specific deployed contract.
// func NewContractFilterer(address common.Address, filterer bind.ContractFilterer) (*ContractFilterer, error) {
// 	contract, err := bindContract(address, nil, nil, filterer)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &ContractFilterer{contract: contract}, nil
// }

// // bindContract binds a generic wrapper to an already deployed contract.
// func bindContract(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
// 	parsed, err := ContractMetaData.GetAbi()
// 	if err != nil {
// 		return nil, err
// 	}
// 	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
// }

// // Call invokes the (constant) contract method with params as input values and
// // sets the output to result. The result type might be a single field for simple
// // returns, a slice of interfaces for anonymous returns and a struct for named
// // returns.
// func (_Contract *ContractRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
// 	return _Contract.Contract.ContractCaller.contract.Call(opts, result, method, params...)
// }

// // Transfer initiates a plain transaction to move funds to the contract, calling
// // its default method if one is available.
// func (_Contract *ContractRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
// 	return _Contract.Contract.ContractTransactor.contract.Transfer(opts)
// }

// // Transact invokes the (paid) contract method with params as input values.
// func (_Contract *ContractRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
// 	return _Contract.Contract.ContractTransactor.contract.Transact(opts, method, params...)
// }

// // Call invokes the (constant) contract method with params as input values and
// // sets the output to result. The result type might be a single field for simple
// // returns, a slice of interfaces for anonymous returns and a struct for named
// // returns.
// func (_Contract *ContractCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
// 	return _Contract.Contract.contract.Call(opts, result, method, params...)
// }

// // Transfer initiates a plain transaction to move funds to the contract, calling
// // its default method if one is available.
// func (_Contract *ContractTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
// 	return _Contract.Contract.contract.Transfer(opts)
// }

// // Transact invokes the (paid) contract method with params as input values.
// func (_Contract *ContractTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
// 	return _Contract.Contract.contract.Transact(opts, method, params...)
// }

// // GetCID is a free data retrieval call binding the contract method 0x1c461ee1.
// //
// // Solidity: function getCID(string key) view returns(string)
// func (_Contract *ContractCaller) GetCID(opts *bind.CallOpts, key string) (string, error) {
// 	var out []interface{}
// 	err := _Contract.contract.Call(opts, &out, "getCID", key)

// 	if err != nil {
// 		return *new(string), err
// 	}

// 	out0 := *abi.ConvertType(out[0], new(string)).(*string)

// 	return out0, err

// }

// // GetCID is a free data retrieval call binding the contract method 0x1c461ee1.
// //
// // Solidity: function getCID(string key) view returns(string)
// func (_Contract *ContractSession) GetCID(key string) (string, error) {
// 	return _Contract.Contract.GetCID(&_Contract.CallOpts, key)
// }

// // GetCID is a free data retrieval call binding the contract method 0x1c461ee1.
// //
// // Solidity: function getCID(string key) view returns(string)
// func (_Contract *ContractCallerSession) GetCID(key string) (string, error) {
// 	return _Contract.Contract.GetCID(&_Contract.CallOpts, key)
// }

// // SaveCID is a paid mutator transaction binding the contract method 0x620403ad.
// //
// // Solidity: function saveCID(string key, string cid) returns()
// func (_Contract *ContractTransactor) SaveCID(opts *bind.TransactOpts, key string, cid string) (*types.Transaction, error) {
// 	return _Contract.contract.Transact(opts, "saveCID", key, cid)
// }

// // SaveCID is a paid mutator transaction binding the contract method 0x620403ad.
// //
// // Solidity: function saveCID(string key, string cid) returns()
// func (_Contract *ContractSession) SaveCID(key string, cid string) (*types.Transaction, error) {
// 	return _Contract.Contract.SaveCID(&_Contract.TransactOpts, key, cid)
// }

// // SaveCID is a paid mutator transaction binding the contract method 0x620403ad.
// //
// // Solidity: function saveCID(string key, string cid) returns()
// func (_Contract *ContractTransactorSession) SaveCID(key string, cid string) (*types.Transaction, error) {
// 	return _Contract.Contract.SaveCID(&_Contract.TransactOpts, key, cid)
// }

// // ContractCIDSavedIterator is returned from FilterCIDSaved and is used to iterate over the raw logs and unpacked data for CIDSaved events raised by the Contract contract.
// type ContractCIDSavedIterator struct {
// 	Event *ContractCIDSaved // Event containing the contract specifics and raw log

// 	contract *bind.BoundContract // Generic contract to use for unpacking event data
// 	event    string              // Event name to use for unpacking event data

// 	logs chan types.Log        // Log channel receiving the found contract events
// 	sub  ethereum.Subscription // Subscription for errors, completion and termination
// 	done bool                  // Whether the subscription completed delivering logs
// 	fail error                 // Occurred error to stop iteration
// }

// // Next advances the iterator to the subsequent event, returning whether there
// // are any more events found. In case of a retrieval or parsing error, false is
// // returned and Error() can be queried for the exact failure.
// func (it *ContractCIDSavedIterator) Next() bool {
// 	// If the iterator failed, stop iterating
// 	if it.fail != nil {
// 		return false
// 	}
// 	// If the iterator completed, deliver directly whatever's available
// 	if it.done {
// 		select {
// 		case log := <-it.logs:
// 			it.Event = new(ContractCIDSaved)
// 			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
// 				it.fail = err
// 				return false
// 			}
// 			it.Event.Raw = log
// 			return true

// 		default:
// 			return false
// 		}
// 	}
// 	// Iterator still in progress, wait for either a data or an error event
// 	select {
// 	case log := <-it.logs:
// 		it.Event = new(ContractCIDSaved)
// 		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
// 			it.fail = err
// 			return false
// 		}
// 		it.Event.Raw = log
// 		return true

// 	case err := <-it.sub.Err():
// 		it.done = true
// 		it.fail = err
// 		return it.Next()
// 	}
// }

// // Error returns any retrieval or parsing error occurred during filtering.
// func (it *ContractCIDSavedIterator) Error() error {
// 	return it.fail
// }

// // Close terminates the iteration process, releasing any pending underlying
// // resources.
// func (it *ContractCIDSavedIterator) Close() error {
// 	it.sub.Unsubscribe()
// 	return nil
// }

// // ContractCIDSaved represents a CIDSaved event raised by the Contract contract.
// type ContractCIDSaved struct {
// 	Cid string
// 	Raw types.Log // Blockchain specific contextual infos
// }

// // FilterCIDSaved is a free log retrieval operation binding the contract event 0x05cbad985f9ea427f160237a3491e243b4f0a10a8969b4e630859eda4aeb66c3.
// //
// // Solidity: event CIDSaved(string cid)
// func (_Contract *ContractFilterer) FilterCIDSaved(opts *bind.FilterOpts) (*ContractCIDSavedIterator, error) {

// 	logs, sub, err := _Contract.contract.FilterLogs(opts, "CIDSaved")
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &ContractCIDSavedIterator{contract: _Contract.contract, event: "CIDSaved", logs: logs, sub: sub}, nil
// }

// // WatchCIDSaved is a free log subscription operation binding the contract event 0x05cbad985f9ea427f160237a3491e243b4f0a10a8969b4e630859eda4aeb66c3.
// //
// // Solidity: event CIDSaved(string cid)
// func (_Contract *ContractFilterer) WatchCIDSaved(opts *bind.WatchOpts, sink chan<- *ContractCIDSaved) (event.Subscription, error) {

// 	logs, sub, err := _Contract.contract.WatchLogs(opts, "CIDSaved")
// 	if err != nil {
// 		return nil, err
// 	}
// 	return event.NewSubscription(func(quit <-chan struct{}) error {
// 		defer sub.Unsubscribe()
// 		for {
// 			select {
// 			case log := <-logs:
// 				// New log arrived, parse the event and forward to the user
// 				event := new(ContractCIDSaved)
// 				if err := _Contract.contract.UnpackLog(event, "CIDSaved", log); err != nil {
// 					return err
// 				}
// 				event.Raw = log

// 				select {
// 				case sink <- event:
// 				case err := <-sub.Err():
// 					return err
// 				case <-quit:
// 					return nil
// 				}
// 			case err := <-sub.Err():
// 				return err
// 			case <-quit:
// 				return nil
// 			}
// 		}
// 	}), nil
// }

// // ParseCIDSaved is a log parse operation binding the contract event 0x05cbad985f9ea427f160237a3491e243b4f0a10a8969b4e630859eda4aeb66c3.
// //
// // Solidity: event CIDSaved(string cid)
// func (_Contract *ContractFilterer) ParseCIDSaved(log types.Log) (*ContractCIDSaved, error) {
// 	event := new(ContractCIDSaved)
// 	if err := _Contract.contract.UnpackLog(event, "CIDSaved", log); err != nil {
// 		return nil, err
// 	}
// 	event.Raw = log
// 	return event, nil
// }