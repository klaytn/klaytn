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
const BridgeAnotherVersionABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"VERSION\",\"outputs\":[{\"name\":\"\",\"type\":\"uint64\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"}]"

// BridgeAnotherVersionBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const BridgeAnotherVersionBinRuntime = `6080604052348015600f57600080fd5b506004361060285760003560e01c8063ffa1ad7414602d575b600080fd5b60336050565b6040805167ffffffffffffffff9092168252519081900360200190f35b60028156fea165627a7a72305820da45156b665fabcb5c1572fd0c9a83cb13075cf1b9d0af72d3ad657151d7ead00029`

// BridgeAnotherVersionFuncSigs maps the 4-byte function signature to its string representation.
var BridgeAnotherVersionFuncSigs = map[string]string{
	"ffa1ad74": "VERSION()",
}

// BridgeAnotherVersionBin is the compiled bytecode used for deploying new contracts.
var BridgeAnotherVersionBin = "0x6080604052348015600f57600080fd5b5060818061001e6000396000f3fe6080604052348015600f57600080fd5b506004361060285760003560e01c8063ffa1ad7414602d575b600080fd5b60336050565b6040805167ffffffffffffffff9092168252519081900360200190f35b60028156fea165627a7a72305820da45156b665fabcb5c1572fd0c9a83cb13075cf1b9d0af72d3ad657151d7ead00029"

// DeployBridgeAnotherVersion deploys a new Klaytn contract, binding an instance of BridgeAnotherVersion to it.
func DeployBridgeAnotherVersion(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *BridgeAnotherVersion, error) {
	parsed, err := abi.JSON(strings.NewReader(BridgeAnotherVersionABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(BridgeAnotherVersionBin), backend)
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
