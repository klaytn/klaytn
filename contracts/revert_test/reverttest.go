// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package revertcontract

import (
	"math/big"
	"strings"

	"github.com/klaytn/klaytn"
	"github.com/klaytn/klaytn/accounts/abi"
	"github.com/klaytn/klaytn/accounts/abi/bind"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = klaytn.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// RevertContractABI is the input ABI used to generate the binding from.
const RevertContractABI = "[{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"fallback\"}]"

// RevertContractBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const RevertContractBinRuntime = `60806040819052600160e51b62461bcd0281526020608452601260a4527f72657665727420696e2066616c6c6261636b000000000000000000000000000060c452606490fdfea165627a7a7230582096379605b267a16a9ce05dcb93552bd99e9597695737798d500eaf84bbf56bce0029`

// RevertContractBin is the compiled bytecode used for deploying new contracts.
var RevertContractBin = "0x6080604052348015600f57600080fd5b50607180601d6000396000f3fe60806040819052600160e51b62461bcd0281526020608452601260a4527f72657665727420696e2066616c6c6261636b000000000000000000000000000060c452606490fdfea165627a7a7230582096379605b267a16a9ce05dcb93552bd99e9597695737798d500eaf84bbf56bce0029"

// DeployRevertContract deploys a new Klaytn contract, binding an instance of RevertContract to it.
func DeployRevertContract(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *RevertContract, error) {
	parsed, err := abi.JSON(strings.NewReader(RevertContractABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(RevertContractBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &RevertContract{RevertContractCaller: RevertContractCaller{contract: contract}, RevertContractTransactor: RevertContractTransactor{contract: contract}, RevertContractFilterer: RevertContractFilterer{contract: contract}}, nil
}

// RevertContract is an auto generated Go binding around a Klaytn contract.
type RevertContract struct {
	RevertContractCaller     // Read-only binding to the contract
	RevertContractTransactor // Write-only binding to the contract
	RevertContractFilterer   // Log filterer for contract events
}

// RevertContractCaller is an auto generated read-only Go binding around a Klaytn contract.
type RevertContractCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RevertContractTransactor is an auto generated write-only Go binding around a Klaytn contract.
type RevertContractTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RevertContractFilterer is an auto generated log filtering Go binding around a Klaytn contract events.
type RevertContractFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RevertContractSession is an auto generated Go binding around a Klaytn contract,
// with pre-set call and transact options.
type RevertContractSession struct {
	Contract     *RevertContract   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// RevertContractCallerSession is an auto generated read-only Go binding around a Klaytn contract,
// with pre-set call options.
type RevertContractCallerSession struct {
	Contract *RevertContractCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// RevertContractTransactorSession is an auto generated write-only Go binding around a Klaytn contract,
// with pre-set transact options.
type RevertContractTransactorSession struct {
	Contract     *RevertContractTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// RevertContractRaw is an auto generated low-level Go binding around a Klaytn contract.
type RevertContractRaw struct {
	Contract *RevertContract // Generic contract binding to access the raw methods on
}

// RevertContractCallerRaw is an auto generated low-level read-only Go binding around a Klaytn contract.
type RevertContractCallerRaw struct {
	Contract *RevertContractCaller // Generic read-only contract binding to access the raw methods on
}

// RevertContractTransactorRaw is an auto generated low-level write-only Go binding around a Klaytn contract.
type RevertContractTransactorRaw struct {
	Contract *RevertContractTransactor // Generic write-only contract binding to access the raw methods on
}

// NewRevertContract creates a new instance of RevertContract, bound to a specific deployed contract.
func NewRevertContract(address common.Address, backend bind.ContractBackend) (*RevertContract, error) {
	contract, err := bindRevertContract(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &RevertContract{RevertContractCaller: RevertContractCaller{contract: contract}, RevertContractTransactor: RevertContractTransactor{contract: contract}, RevertContractFilterer: RevertContractFilterer{contract: contract}}, nil
}

// NewRevertContractCaller creates a new read-only instance of RevertContract, bound to a specific deployed contract.
func NewRevertContractCaller(address common.Address, caller bind.ContractCaller) (*RevertContractCaller, error) {
	contract, err := bindRevertContract(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &RevertContractCaller{contract: contract}, nil
}

// NewRevertContractTransactor creates a new write-only instance of RevertContract, bound to a specific deployed contract.
func NewRevertContractTransactor(address common.Address, transactor bind.ContractTransactor) (*RevertContractTransactor, error) {
	contract, err := bindRevertContract(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &RevertContractTransactor{contract: contract}, nil
}

// NewRevertContractFilterer creates a new log filterer instance of RevertContract, bound to a specific deployed contract.
func NewRevertContractFilterer(address common.Address, filterer bind.ContractFilterer) (*RevertContractFilterer, error) {
	contract, err := bindRevertContract(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &RevertContractFilterer{contract: contract}, nil
}

// bindRevertContract binds a generic wrapper to an already deployed contract.
func bindRevertContract(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(RevertContractABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_RevertContract *RevertContractRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _RevertContract.Contract.RevertContractCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_RevertContract *RevertContractRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RevertContract.Contract.RevertContractTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_RevertContract *RevertContractRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RevertContract.Contract.RevertContractTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_RevertContract *RevertContractCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _RevertContract.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_RevertContract *RevertContractTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RevertContract.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_RevertContract *RevertContractTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RevertContract.Contract.contract.Transact(opts, method, params...)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_RevertContract *RevertContractTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _RevertContract.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_RevertContract *RevertContractSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _RevertContract.Contract.Fallback(&_RevertContract.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_RevertContract *RevertContractTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _RevertContract.Contract.Fallback(&_RevertContract.TransactOpts, calldata)
}
