// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package bridge

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

// BridgeAnotherVersionABI is the input ABI used to generate the binding from.
const BridgeAnotherVersionABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"VERSION\",\"outputs\":[{\"name\":\"\",\"type\":\"uint64\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"}]"

// BridgeAnotherVersionBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const BridgeAnotherVersionBinRuntime = `6080604052348015600f57600080fd5b506004361060285760003560e01c8063ffa1ad7414602d575b600080fd5b60336050565b6040805167ffffffffffffffff9092168252519081900360200190f35b60005467ffffffffffffffff168156fea165627a7a72305820deda475869c31c7eea972dca13cb66d75fc5c8d12d64d0798ee640146c7288030029`

// BridgeAnotherVersionFuncSigs maps the 4-byte function signature to its string representation.
var BridgeAnotherVersionFuncSigs = map[string]string{
	"ffa1ad74": "VERSION()",
}

// BridgeAnotherVersionBin is the compiled bytecode used for deploying new contracts.
var BridgeAnotherVersionBin = "0x6080604052600080546001600160401b031916905534801561002057600080fd5b506040516020806100fd8339810180604052602081101561004057600080fd5b5051600080546001600160401b039092166001600160401b0319909216919091179055608c806100716000396000f3fe6080604052348015600f57600080fd5b506004361060285760003560e01c8063ffa1ad7414602d575b600080fd5b60336050565b6040805167ffffffffffffffff9092168252519081900360200190f35b60005467ffffffffffffffff168156fea165627a7a72305820deda475869c31c7eea972dca13cb66d75fc5c8d12d64d0798ee640146c7288030029"

// DeployBridgeAnotherVersion deploys a new Klaytn contract, binding an instance of BridgeAnotherVersion to it.
func DeployBridgeAnotherVersion(auth *bind.TransactOpts, backend bind.ContractBackend, version uint64) (common.Address, *types.Transaction, *BridgeAnotherVersion, error) {
	parsed, err := abi.JSON(strings.NewReader(BridgeAnotherVersionABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(BridgeAnotherVersionBin), backend, version)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &BridgeAnotherVersion{BridgeAnotherVersionCaller: BridgeAnotherVersionCaller{contract: contract}, BridgeAnotherVersionTransactor: BridgeAnotherVersionTransactor{contract: contract}, BridgeAnotherVersionFilterer: BridgeAnotherVersionFilterer{contract: contract}}, nil
}

// BridgeAnotherVersion is an auto generated Go binding around a Klaytn contract.
type BridgeAnotherVersion struct {
	BridgeAnotherVersionCaller     // Read-only binding to the contract
	BridgeAnotherVersionTransactor // Write-only binding to the contract
	BridgeAnotherVersionFilterer   // Log filterer for contract events
}

// BridgeAnotherVersionCaller is an auto generated read-only Go binding around a Klaytn contract.
type BridgeAnotherVersionCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BridgeAnotherVersionTransactor is an auto generated write-only Go binding around a Klaytn contract.
type BridgeAnotherVersionTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BridgeAnotherVersionFilterer is an auto generated log filtering Go binding around a Klaytn contract events.
type BridgeAnotherVersionFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BridgeAnotherVersionSession is an auto generated Go binding around a Klaytn contract,
// with pre-set call and transact options.
type BridgeAnotherVersionSession struct {
	Contract     *BridgeAnotherVersion // Generic contract binding to set the session for
	CallOpts     bind.CallOpts         // Call options to use throughout this session
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// BridgeAnotherVersionCallerSession is an auto generated read-only Go binding around a Klaytn contract,
// with pre-set call options.
type BridgeAnotherVersionCallerSession struct {
	Contract *BridgeAnotherVersionCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts               // Call options to use throughout this session
}

// BridgeAnotherVersionTransactorSession is an auto generated write-only Go binding around a Klaytn contract,
// with pre-set transact options.
type BridgeAnotherVersionTransactorSession struct {
	Contract     *BridgeAnotherVersionTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts               // Transaction auth options to use throughout this session
}

// BridgeAnotherVersionRaw is an auto generated low-level Go binding around a Klaytn contract.
type BridgeAnotherVersionRaw struct {
	Contract *BridgeAnotherVersion // Generic contract binding to access the raw methods on
}

// BridgeAnotherVersionCallerRaw is an auto generated low-level read-only Go binding around a Klaytn contract.
type BridgeAnotherVersionCallerRaw struct {
	Contract *BridgeAnotherVersionCaller // Generic read-only contract binding to access the raw methods on
}

// BridgeAnotherVersionTransactorRaw is an auto generated low-level write-only Go binding around a Klaytn contract.
type BridgeAnotherVersionTransactorRaw struct {
	Contract *BridgeAnotherVersionTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBridgeAnotherVersion creates a new instance of BridgeAnotherVersion, bound to a specific deployed contract.
func NewBridgeAnotherVersion(address common.Address, backend bind.ContractBackend) (*BridgeAnotherVersion, error) {
	contract, err := bindBridgeAnotherVersion(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &BridgeAnotherVersion{BridgeAnotherVersionCaller: BridgeAnotherVersionCaller{contract: contract}, BridgeAnotherVersionTransactor: BridgeAnotherVersionTransactor{contract: contract}, BridgeAnotherVersionFilterer: BridgeAnotherVersionFilterer{contract: contract}}, nil
}

// NewBridgeAnotherVersionCaller creates a new read-only instance of BridgeAnotherVersion, bound to a specific deployed contract.
func NewBridgeAnotherVersionCaller(address common.Address, caller bind.ContractCaller) (*BridgeAnotherVersionCaller, error) {
	contract, err := bindBridgeAnotherVersion(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BridgeAnotherVersionCaller{contract: contract}, nil
}

// NewBridgeAnotherVersionTransactor creates a new write-only instance of BridgeAnotherVersion, bound to a specific deployed contract.
func NewBridgeAnotherVersionTransactor(address common.Address, transactor bind.ContractTransactor) (*BridgeAnotherVersionTransactor, error) {
	contract, err := bindBridgeAnotherVersion(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BridgeAnotherVersionTransactor{contract: contract}, nil
}

// NewBridgeAnotherVersionFilterer creates a new log filterer instance of BridgeAnotherVersion, bound to a specific deployed contract.
func NewBridgeAnotherVersionFilterer(address common.Address, filterer bind.ContractFilterer) (*BridgeAnotherVersionFilterer, error) {
	contract, err := bindBridgeAnotherVersion(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BridgeAnotherVersionFilterer{contract: contract}, nil
}

// bindBridgeAnotherVersion binds a generic wrapper to an already deployed contract.
func bindBridgeAnotherVersion(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(BridgeAnotherVersionABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BridgeAnotherVersion *BridgeAnotherVersionRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _BridgeAnotherVersion.Contract.BridgeAnotherVersionCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BridgeAnotherVersion *BridgeAnotherVersionRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BridgeAnotherVersion.Contract.BridgeAnotherVersionTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BridgeAnotherVersion *BridgeAnotherVersionRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BridgeAnotherVersion.Contract.BridgeAnotherVersionTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BridgeAnotherVersion *BridgeAnotherVersionCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _BridgeAnotherVersion.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BridgeAnotherVersion *BridgeAnotherVersionTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BridgeAnotherVersion.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BridgeAnotherVersion *BridgeAnotherVersionTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BridgeAnotherVersion.Contract.contract.Transact(opts, method, params...)
}

// VERSION is a free data retrieval call binding the contract method 0xffa1ad74.
//
// Solidity: function VERSION() view returns(uint64)
func (_BridgeAnotherVersion *BridgeAnotherVersionCaller) VERSION(opts *bind.CallOpts) (uint64, error) {
	var (
		ret0 = new(uint64)
	)
	out := ret0
	err := _BridgeAnotherVersion.contract.Call(opts, out, "VERSION")
	return *ret0, err
}

// VERSION is a free data retrieval call binding the contract method 0xffa1ad74.
//
// Solidity: function VERSION() view returns(uint64)
func (_BridgeAnotherVersion *BridgeAnotherVersionSession) VERSION() (uint64, error) {
	return _BridgeAnotherVersion.Contract.VERSION(&_BridgeAnotherVersion.CallOpts)
}

// VERSION is a free data retrieval call binding the contract method 0xffa1ad74.
//
// Solidity: function VERSION() view returns(uint64)
func (_BridgeAnotherVersion *BridgeAnotherVersionCallerSession) VERSION() (uint64, error) {
	return _BridgeAnotherVersion.Contract.VERSION(&_BridgeAnotherVersion.CallOpts)
}

// BridgeVersionControllerABI is the input ABI used to generate the binding from.
const BridgeVersionControllerABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"VERSION\",\"outputs\":[{\"name\":\"\",\"type\":\"uint64\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"}]"

// BridgeVersionControllerBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const BridgeVersionControllerBinRuntime = `6080604052348015600f57600080fd5b506004361060285760003560e01c8063ffa1ad7414602d575b600080fd5b60336050565b6040805167ffffffffffffffff9092168252519081900360200190f35b60005467ffffffffffffffff168156fea165627a7a72305820103675040118bf16fec749270367b54bdfcfc8e40fb9174d60a7b20584aaa0240029`

// BridgeVersionControllerFuncSigs maps the 4-byte function signature to its string representation.
var BridgeVersionControllerFuncSigs = map[string]string{
	"ffa1ad74": "VERSION()",
}

// BridgeVersionControllerBin is the compiled bytecode used for deploying new contracts.
var BridgeVersionControllerBin = "0x6080604052600080546001600160401b031916905534801561002057600080fd5b506040516020806100fd8339810180604052602081101561004057600080fd5b5051600080546001600160401b039092166001600160401b0319909216919091179055608c806100716000396000f3fe6080604052348015600f57600080fd5b506004361060285760003560e01c8063ffa1ad7414602d575b600080fd5b60336050565b6040805167ffffffffffffffff9092168252519081900360200190f35b60005467ffffffffffffffff168156fea165627a7a72305820103675040118bf16fec749270367b54bdfcfc8e40fb9174d60a7b20584aaa0240029"

// DeployBridgeVersionController deploys a new Klaytn contract, binding an instance of BridgeVersionController to it.
func DeployBridgeVersionController(auth *bind.TransactOpts, backend bind.ContractBackend, version uint64) (common.Address, *types.Transaction, *BridgeVersionController, error) {
	parsed, err := abi.JSON(strings.NewReader(BridgeVersionControllerABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(BridgeVersionControllerBin), backend, version)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &BridgeVersionController{BridgeVersionControllerCaller: BridgeVersionControllerCaller{contract: contract}, BridgeVersionControllerTransactor: BridgeVersionControllerTransactor{contract: contract}, BridgeVersionControllerFilterer: BridgeVersionControllerFilterer{contract: contract}}, nil
}

// BridgeVersionController is an auto generated Go binding around a Klaytn contract.
type BridgeVersionController struct {
	BridgeVersionControllerCaller     // Read-only binding to the contract
	BridgeVersionControllerTransactor // Write-only binding to the contract
	BridgeVersionControllerFilterer   // Log filterer for contract events
}

// BridgeVersionControllerCaller is an auto generated read-only Go binding around a Klaytn contract.
type BridgeVersionControllerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BridgeVersionControllerTransactor is an auto generated write-only Go binding around a Klaytn contract.
type BridgeVersionControllerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BridgeVersionControllerFilterer is an auto generated log filtering Go binding around a Klaytn contract events.
type BridgeVersionControllerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BridgeVersionControllerSession is an auto generated Go binding around a Klaytn contract,
// with pre-set call and transact options.
type BridgeVersionControllerSession struct {
	Contract     *BridgeVersionController // Generic contract binding to set the session for
	CallOpts     bind.CallOpts            // Call options to use throughout this session
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// BridgeVersionControllerCallerSession is an auto generated read-only Go binding around a Klaytn contract,
// with pre-set call options.
type BridgeVersionControllerCallerSession struct {
	Contract *BridgeVersionControllerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                  // Call options to use throughout this session
}

// BridgeVersionControllerTransactorSession is an auto generated write-only Go binding around a Klaytn contract,
// with pre-set transact options.
type BridgeVersionControllerTransactorSession struct {
	Contract     *BridgeVersionControllerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                  // Transaction auth options to use throughout this session
}

// BridgeVersionControllerRaw is an auto generated low-level Go binding around a Klaytn contract.
type BridgeVersionControllerRaw struct {
	Contract *BridgeVersionController // Generic contract binding to access the raw methods on
}

// BridgeVersionControllerCallerRaw is an auto generated low-level read-only Go binding around a Klaytn contract.
type BridgeVersionControllerCallerRaw struct {
	Contract *BridgeVersionControllerCaller // Generic read-only contract binding to access the raw methods on
}

// BridgeVersionControllerTransactorRaw is an auto generated low-level write-only Go binding around a Klaytn contract.
type BridgeVersionControllerTransactorRaw struct {
	Contract *BridgeVersionControllerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBridgeVersionController creates a new instance of BridgeVersionController, bound to a specific deployed contract.
func NewBridgeVersionController(address common.Address, backend bind.ContractBackend) (*BridgeVersionController, error) {
	contract, err := bindBridgeVersionController(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &BridgeVersionController{BridgeVersionControllerCaller: BridgeVersionControllerCaller{contract: contract}, BridgeVersionControllerTransactor: BridgeVersionControllerTransactor{contract: contract}, BridgeVersionControllerFilterer: BridgeVersionControllerFilterer{contract: contract}}, nil
}

// NewBridgeVersionControllerCaller creates a new read-only instance of BridgeVersionController, bound to a specific deployed contract.
func NewBridgeVersionControllerCaller(address common.Address, caller bind.ContractCaller) (*BridgeVersionControllerCaller, error) {
	contract, err := bindBridgeVersionController(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BridgeVersionControllerCaller{contract: contract}, nil
}

// NewBridgeVersionControllerTransactor creates a new write-only instance of BridgeVersionController, bound to a specific deployed contract.
func NewBridgeVersionControllerTransactor(address common.Address, transactor bind.ContractTransactor) (*BridgeVersionControllerTransactor, error) {
	contract, err := bindBridgeVersionController(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BridgeVersionControllerTransactor{contract: contract}, nil
}

// NewBridgeVersionControllerFilterer creates a new log filterer instance of BridgeVersionController, bound to a specific deployed contract.
func NewBridgeVersionControllerFilterer(address common.Address, filterer bind.ContractFilterer) (*BridgeVersionControllerFilterer, error) {
	contract, err := bindBridgeVersionController(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BridgeVersionControllerFilterer{contract: contract}, nil
}

// bindBridgeVersionController binds a generic wrapper to an already deployed contract.
func bindBridgeVersionController(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(BridgeVersionControllerABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BridgeVersionController *BridgeVersionControllerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _BridgeVersionController.Contract.BridgeVersionControllerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BridgeVersionController *BridgeVersionControllerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BridgeVersionController.Contract.BridgeVersionControllerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BridgeVersionController *BridgeVersionControllerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BridgeVersionController.Contract.BridgeVersionControllerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BridgeVersionController *BridgeVersionControllerCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _BridgeVersionController.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BridgeVersionController *BridgeVersionControllerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BridgeVersionController.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BridgeVersionController *BridgeVersionControllerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BridgeVersionController.Contract.contract.Transact(opts, method, params...)
}

// VERSION is a free data retrieval call binding the contract method 0xffa1ad74.
//
// Solidity: function VERSION() view returns(uint64)
func (_BridgeVersionController *BridgeVersionControllerCaller) VERSION(opts *bind.CallOpts) (uint64, error) {
	var (
		ret0 = new(uint64)
	)
	out := ret0
	err := _BridgeVersionController.contract.Call(opts, out, "VERSION")
	return *ret0, err
}

// VERSION is a free data retrieval call binding the contract method 0xffa1ad74.
//
// Solidity: function VERSION() view returns(uint64)
func (_BridgeVersionController *BridgeVersionControllerSession) VERSION() (uint64, error) {
	return _BridgeVersionController.Contract.VERSION(&_BridgeVersionController.CallOpts)
}

// VERSION is a free data retrieval call binding the contract method 0xffa1ad74.
//
// Solidity: function VERSION() view returns(uint64)
func (_BridgeVersionController *BridgeVersionControllerCallerSession) VERSION() (uint64, error) {
	return _BridgeVersionController.Contract.VERSION(&_BridgeVersionController.CallOpts)
}
