// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package kip13

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

// InterfaceIdentifierABI is the input ABI used to generate the binding from.
const InterfaceIdentifierABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"interfaceID\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"

// InterfaceIdentifierBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const InterfaceIdentifierBinRuntime = ``

// InterfaceIdentifierFuncSigs maps the 4-byte function signature to its string representation.
var InterfaceIdentifierFuncSigs = map[string]string{
	"01ffc9a7": "supportsInterface(bytes4)",
}

// InterfaceIdentifier is an auto generated Go binding around a Klaytn contract.
type InterfaceIdentifier struct {
	InterfaceIdentifierCaller     // Read-only binding to the contract
	InterfaceIdentifierTransactor // Write-only binding to the contract
	InterfaceIdentifierFilterer   // Log filterer for contract events
}

// InterfaceIdentifierCaller is an auto generated read-only Go binding around a Klaytn contract.
type InterfaceIdentifierCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// InterfaceIdentifierTransactor is an auto generated write-only Go binding around a Klaytn contract.
type InterfaceIdentifierTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// InterfaceIdentifierFilterer is an auto generated log filtering Go binding around a Klaytn contract events.
type InterfaceIdentifierFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// InterfaceIdentifierSession is an auto generated Go binding around a Klaytn contract,
// with pre-set call and transact options.
type InterfaceIdentifierSession struct {
	Contract     *InterfaceIdentifier // Generic contract binding to set the session for
	CallOpts     bind.CallOpts        // Call options to use throughout this session
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// InterfaceIdentifierCallerSession is an auto generated read-only Go binding around a Klaytn contract,
// with pre-set call options.
type InterfaceIdentifierCallerSession struct {
	Contract *InterfaceIdentifierCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts              // Call options to use throughout this session
}

// InterfaceIdentifierTransactorSession is an auto generated write-only Go binding around a Klaytn contract,
// with pre-set transact options.
type InterfaceIdentifierTransactorSession struct {
	Contract     *InterfaceIdentifierTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts              // Transaction auth options to use throughout this session
}

// InterfaceIdentifierRaw is an auto generated low-level Go binding around a Klaytn contract.
type InterfaceIdentifierRaw struct {
	Contract *InterfaceIdentifier // Generic contract binding to access the raw methods on
}

// InterfaceIdentifierCallerRaw is an auto generated low-level read-only Go binding around a Klaytn contract.
type InterfaceIdentifierCallerRaw struct {
	Contract *InterfaceIdentifierCaller // Generic read-only contract binding to access the raw methods on
}

// InterfaceIdentifierTransactorRaw is an auto generated low-level write-only Go binding around a Klaytn contract.
type InterfaceIdentifierTransactorRaw struct {
	Contract *InterfaceIdentifierTransactor // Generic write-only contract binding to access the raw methods on
}

// NewInterfaceIdentifier creates a new instance of InterfaceIdentifier, bound to a specific deployed contract.
func NewInterfaceIdentifier(address common.Address, backend bind.ContractBackend) (*InterfaceIdentifier, error) {
	contract, err := bindInterfaceIdentifier(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &InterfaceIdentifier{InterfaceIdentifierCaller: InterfaceIdentifierCaller{contract: contract}, InterfaceIdentifierTransactor: InterfaceIdentifierTransactor{contract: contract}, InterfaceIdentifierFilterer: InterfaceIdentifierFilterer{contract: contract}}, nil
}

// NewInterfaceIdentifierCaller creates a new read-only instance of InterfaceIdentifier, bound to a specific deployed contract.
func NewInterfaceIdentifierCaller(address common.Address, caller bind.ContractCaller) (*InterfaceIdentifierCaller, error) {
	contract, err := bindInterfaceIdentifier(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &InterfaceIdentifierCaller{contract: contract}, nil
}

// NewInterfaceIdentifierTransactor creates a new write-only instance of InterfaceIdentifier, bound to a specific deployed contract.
func NewInterfaceIdentifierTransactor(address common.Address, transactor bind.ContractTransactor) (*InterfaceIdentifierTransactor, error) {
	contract, err := bindInterfaceIdentifier(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &InterfaceIdentifierTransactor{contract: contract}, nil
}

// NewInterfaceIdentifierFilterer creates a new log filterer instance of InterfaceIdentifier, bound to a specific deployed contract.
func NewInterfaceIdentifierFilterer(address common.Address, filterer bind.ContractFilterer) (*InterfaceIdentifierFilterer, error) {
	contract, err := bindInterfaceIdentifier(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &InterfaceIdentifierFilterer{contract: contract}, nil
}

// bindInterfaceIdentifier binds a generic wrapper to an already deployed contract.
func bindInterfaceIdentifier(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(InterfaceIdentifierABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_InterfaceIdentifier *InterfaceIdentifierRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _InterfaceIdentifier.Contract.InterfaceIdentifierCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_InterfaceIdentifier *InterfaceIdentifierRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _InterfaceIdentifier.Contract.InterfaceIdentifierTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_InterfaceIdentifier *InterfaceIdentifierRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _InterfaceIdentifier.Contract.InterfaceIdentifierTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_InterfaceIdentifier *InterfaceIdentifierCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _InterfaceIdentifier.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_InterfaceIdentifier *InterfaceIdentifierTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _InterfaceIdentifier.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_InterfaceIdentifier *InterfaceIdentifierTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _InterfaceIdentifier.Contract.contract.Transact(opts, method, params...)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceID) view returns(bool)
func (_InterfaceIdentifier *InterfaceIdentifierCaller) SupportsInterface(opts *bind.CallOpts, interfaceID [4]byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _InterfaceIdentifier.contract.Call(opts, out, "supportsInterface", interfaceID)
	return *ret0, err
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceID) view returns(bool)
func (_InterfaceIdentifier *InterfaceIdentifierSession) SupportsInterface(interfaceID [4]byte) (bool, error) {
	return _InterfaceIdentifier.Contract.SupportsInterface(&_InterfaceIdentifier.CallOpts, interfaceID)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceID) view returns(bool)
func (_InterfaceIdentifier *InterfaceIdentifierCallerSession) SupportsInterface(interfaceID [4]byte) (bool, error) {
	return _InterfaceIdentifier.Contract.SupportsInterface(&_InterfaceIdentifier.CallOpts, interfaceID)
}
