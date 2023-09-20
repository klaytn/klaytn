// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package system_contracts

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

// IRegistryABI is the input ABI used to generate the binding from.
const IRegistryABI = "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"name\":\"getActiveContract\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"

// IRegistryBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const IRegistryBinRuntime = ``

// IRegistryFuncSigs maps the 4-byte function signature to its string representation.
var IRegistryFuncSigs = map[string]string{
	"080d5721": "getActiveContract(string)",
}

// IRegistry is an auto generated Go binding around a Klaytn contract.
type IRegistry struct {
	IRegistryCaller     // Read-only binding to the contract
	IRegistryTransactor // Write-only binding to the contract
	IRegistryFilterer   // Log filterer for contract events
}

// IRegistryCaller is an auto generated read-only Go binding around a Klaytn contract.
type IRegistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IRegistryTransactor is an auto generated write-only Go binding around a Klaytn contract.
type IRegistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IRegistryFilterer is an auto generated log filtering Go binding around a Klaytn contract events.
type IRegistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IRegistrySession is an auto generated Go binding around a Klaytn contract,
// with pre-set call and transact options.
type IRegistrySession struct {
	Contract     *IRegistry        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IRegistryCallerSession is an auto generated read-only Go binding around a Klaytn contract,
// with pre-set call options.
type IRegistryCallerSession struct {
	Contract *IRegistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// IRegistryTransactorSession is an auto generated write-only Go binding around a Klaytn contract,
// with pre-set transact options.
type IRegistryTransactorSession struct {
	Contract     *IRegistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// IRegistryRaw is an auto generated low-level Go binding around a Klaytn contract.
type IRegistryRaw struct {
	Contract *IRegistry // Generic contract binding to access the raw methods on
}

// IRegistryCallerRaw is an auto generated low-level read-only Go binding around a Klaytn contract.
type IRegistryCallerRaw struct {
	Contract *IRegistryCaller // Generic read-only contract binding to access the raw methods on
}

// IRegistryTransactorRaw is an auto generated low-level write-only Go binding around a Klaytn contract.
type IRegistryTransactorRaw struct {
	Contract *IRegistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIRegistry creates a new instance of IRegistry, bound to a specific deployed contract.
func NewIRegistry(address common.Address, backend bind.ContractBackend) (*IRegistry, error) {
	contract, err := bindIRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IRegistry{IRegistryCaller: IRegistryCaller{contract: contract}, IRegistryTransactor: IRegistryTransactor{contract: contract}, IRegistryFilterer: IRegistryFilterer{contract: contract}}, nil
}

// NewIRegistryCaller creates a new read-only instance of IRegistry, bound to a specific deployed contract.
func NewIRegistryCaller(address common.Address, caller bind.ContractCaller) (*IRegistryCaller, error) {
	contract, err := bindIRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IRegistryCaller{contract: contract}, nil
}

// NewIRegistryTransactor creates a new write-only instance of IRegistry, bound to a specific deployed contract.
func NewIRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*IRegistryTransactor, error) {
	contract, err := bindIRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IRegistryTransactor{contract: contract}, nil
}

// NewIRegistryFilterer creates a new log filterer instance of IRegistry, bound to a specific deployed contract.
func NewIRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*IRegistryFilterer, error) {
	contract, err := bindIRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IRegistryFilterer{contract: contract}, nil
}

// bindIRegistry binds a generic wrapper to an already deployed contract.
func bindIRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(IRegistryABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IRegistry *IRegistryRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _IRegistry.Contract.IRegistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IRegistry *IRegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IRegistry.Contract.IRegistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IRegistry *IRegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IRegistry.Contract.IRegistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IRegistry *IRegistryCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _IRegistry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IRegistry *IRegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IRegistry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IRegistry *IRegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IRegistry.Contract.contract.Transact(opts, method, params...)
}

// GetActiveContract is a free data retrieval call binding the contract method 0x080d5721.
//
// Solidity: function getActiveContract(string name) view returns(address)
func (_IRegistry *IRegistryCaller) GetActiveContract(opts *bind.CallOpts, name string) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _IRegistry.contract.Call(opts, out, "getActiveContract", name)
	return *ret0, err
}

// GetActiveContract is a free data retrieval call binding the contract method 0x080d5721.
//
// Solidity: function getActiveContract(string name) view returns(address)
func (_IRegistry *IRegistrySession) GetActiveContract(name string) (common.Address, error) {
	return _IRegistry.Contract.GetActiveContract(&_IRegistry.CallOpts, name)
}

// GetActiveContract is a free data retrieval call binding the contract method 0x080d5721.
//
// Solidity: function getActiveContract(string name) view returns(address)
func (_IRegistry *IRegistryCallerSession) GetActiveContract(name string) (common.Address, error) {
	return _IRegistry.Contract.GetActiveContract(&_IRegistry.CallOpts, name)
}

// RegistryABI is the input ABI used to generate the binding from.
const RegistryABI = "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"name\":\"getActiveContract\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"

// RegistryBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const RegistryBinRuntime = `608060405234801561001057600080fd5b506004361061002b5760003560e01c8063080d572114610030575b600080fd5b61004361003e3660046100b6565b61005f565b6040516001600160a01b03909116815260200160405180910390f35b60405162461bcd60e51b815260206004820152600f60248201526e1b9bdd081a5b5c1b195b595b9d1959608a1b604482015260009060640160405180910390fd5b634e487b7160e01b600052604160045260246000fd5b6000602082840312156100c857600080fd5b813567ffffffffffffffff808211156100e057600080fd5b818401915084601f8301126100f457600080fd5b813581811115610106576101066100a0565b604051601f8201601f19908116603f0116810190838211818310171561012e5761012e6100a0565b8160405282815287602084870101111561014757600080fd5b82602086016020830137600092810160200192909252509594505050505056fea26469706673582212201b2ec5ac3f64bc0b73fe2515f6c12d958b63dbd7387cae5d88affc9c69b01e5664736f6c63430008110033`

// RegistryFuncSigs maps the 4-byte function signature to its string representation.
var RegistryFuncSigs = map[string]string{
	"080d5721": "getActiveContract(string)",
}

// RegistryBin is the compiled bytecode used for deploying new contracts.
var RegistryBin = "0x608060405234801561001057600080fd5b5061019d806100206000396000f3fe608060405234801561001057600080fd5b506004361061002b5760003560e01c8063080d572114610030575b600080fd5b61004361003e3660046100b6565b61005f565b6040516001600160a01b03909116815260200160405180910390f35b60405162461bcd60e51b815260206004820152600f60248201526e1b9bdd081a5b5c1b195b595b9d1959608a1b604482015260009060640160405180910390fd5b634e487b7160e01b600052604160045260246000fd5b6000602082840312156100c857600080fd5b813567ffffffffffffffff808211156100e057600080fd5b818401915084601f8301126100f457600080fd5b813581811115610106576101066100a0565b604051601f8201601f19908116603f0116810190838211818310171561012e5761012e6100a0565b8160405282815287602084870101111561014757600080fd5b82602086016020830137600092810160200192909252509594505050505056fea26469706673582212201b2ec5ac3f64bc0b73fe2515f6c12d958b63dbd7387cae5d88affc9c69b01e5664736f6c63430008110033"

// DeployRegistry deploys a new Klaytn contract, binding an instance of Registry to it.
func DeployRegistry(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Registry, error) {
	parsed, err := abi.JSON(strings.NewReader(RegistryABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(RegistryBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Registry{RegistryCaller: RegistryCaller{contract: contract}, RegistryTransactor: RegistryTransactor{contract: contract}, RegistryFilterer: RegistryFilterer{contract: contract}}, nil
}

// Registry is an auto generated Go binding around a Klaytn contract.
type Registry struct {
	RegistryCaller     // Read-only binding to the contract
	RegistryTransactor // Write-only binding to the contract
	RegistryFilterer   // Log filterer for contract events
}

// RegistryCaller is an auto generated read-only Go binding around a Klaytn contract.
type RegistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RegistryTransactor is an auto generated write-only Go binding around a Klaytn contract.
type RegistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RegistryFilterer is an auto generated log filtering Go binding around a Klaytn contract events.
type RegistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RegistrySession is an auto generated Go binding around a Klaytn contract,
// with pre-set call and transact options.
type RegistrySession struct {
	Contract     *Registry         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// RegistryCallerSession is an auto generated read-only Go binding around a Klaytn contract,
// with pre-set call options.
type RegistryCallerSession struct {
	Contract *RegistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// RegistryTransactorSession is an auto generated write-only Go binding around a Klaytn contract,
// with pre-set transact options.
type RegistryTransactorSession struct {
	Contract     *RegistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// RegistryRaw is an auto generated low-level Go binding around a Klaytn contract.
type RegistryRaw struct {
	Contract *Registry // Generic contract binding to access the raw methods on
}

// RegistryCallerRaw is an auto generated low-level read-only Go binding around a Klaytn contract.
type RegistryCallerRaw struct {
	Contract *RegistryCaller // Generic read-only contract binding to access the raw methods on
}

// RegistryTransactorRaw is an auto generated low-level write-only Go binding around a Klaytn contract.
type RegistryTransactorRaw struct {
	Contract *RegistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewRegistry creates a new instance of Registry, bound to a specific deployed contract.
func NewRegistry(address common.Address, backend bind.ContractBackend) (*Registry, error) {
	contract, err := bindRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Registry{RegistryCaller: RegistryCaller{contract: contract}, RegistryTransactor: RegistryTransactor{contract: contract}, RegistryFilterer: RegistryFilterer{contract: contract}}, nil
}

// NewRegistryCaller creates a new read-only instance of Registry, bound to a specific deployed contract.
func NewRegistryCaller(address common.Address, caller bind.ContractCaller) (*RegistryCaller, error) {
	contract, err := bindRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &RegistryCaller{contract: contract}, nil
}

// NewRegistryTransactor creates a new write-only instance of Registry, bound to a specific deployed contract.
func NewRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*RegistryTransactor, error) {
	contract, err := bindRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &RegistryTransactor{contract: contract}, nil
}

// NewRegistryFilterer creates a new log filterer instance of Registry, bound to a specific deployed contract.
func NewRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*RegistryFilterer, error) {
	contract, err := bindRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &RegistryFilterer{contract: contract}, nil
}

// bindRegistry binds a generic wrapper to an already deployed contract.
func bindRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(RegistryABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Registry *RegistryRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Registry.Contract.RegistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Registry *RegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Registry.Contract.RegistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Registry *RegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Registry.Contract.RegistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Registry *RegistryCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Registry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Registry *RegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Registry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Registry *RegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Registry.Contract.contract.Transact(opts, method, params...)
}

// GetActiveContract is a free data retrieval call binding the contract method 0x080d5721.
//
// Solidity: function getActiveContract(string name) view returns(address)
func (_Registry *RegistryCaller) GetActiveContract(opts *bind.CallOpts, name string) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Registry.contract.Call(opts, out, "getActiveContract", name)
	return *ret0, err
}

// GetActiveContract is a free data retrieval call binding the contract method 0x080d5721.
//
// Solidity: function getActiveContract(string name) view returns(address)
func (_Registry *RegistrySession) GetActiveContract(name string) (common.Address, error) {
	return _Registry.Contract.GetActiveContract(&_Registry.CallOpts, name)
}

// GetActiveContract is a free data retrieval call binding the contract method 0x080d5721.
//
// Solidity: function getActiveContract(string name) view returns(address)
func (_Registry *RegistryCallerSession) GetActiveContract(name string) (common.Address, error) {
	return _Registry.Contract.GetActiveContract(&_Registry.CallOpts, name)
}

// RegistryMockABI is the input ABI used to generate the binding from.
const RegistryMockABI = "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"name\":\"getActiveContract\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"name\":\"records\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"register\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// RegistryMockBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const RegistryMockBinRuntime = `608060405234801561001057600080fd5b50600436106100415760003560e01c8063080d5721146100465780631e59c52914610075578063541e771d1461008a575b600080fd5b6100596100543660046101d5565b6100be565b6040516001600160a01b03909116815260200160405180910390f35b610088610083366004610212565b6100ee565b005b6100596100983660046101d5565b80516020818301810180516000825292820191909301209152546001600160a01b031681565b600080826040516100cf9190610270565b908152604051908190036020019020546001600160a01b031692915050565b806000836040516100ff9190610270565b90815260405190819003602001902080546001600160a01b03929092166001600160a01b03199092169190911790555050565b634e487b7160e01b600052604160045260246000fd5b600082601f83011261015957600080fd5b813567ffffffffffffffff8082111561017457610174610132565b604051601f8301601f19908116603f0116810190828211818310171561019c5761019c610132565b816040528381528660208588010111156101b557600080fd5b836020870160208301376000602085830101528094505050505092915050565b6000602082840312156101e757600080fd5b813567ffffffffffffffff8111156101fe57600080fd5b61020a84828501610148565b949350505050565b6000806040838503121561022557600080fd5b823567ffffffffffffffff81111561023c57600080fd5b61024885828601610148565b92505060208301356001600160a01b038116811461026557600080fd5b809150509250929050565b6000825160005b818110156102915760208186018101518583015201610277565b50600092019182525091905056fea26469706673582212200be7666a50b9065e81c071c025510f3e57b6e8294690ac5940a1c711cf00340b64736f6c63430008110033`

// RegistryMockFuncSigs maps the 4-byte function signature to its string representation.
var RegistryMockFuncSigs = map[string]string{
	"080d5721": "getActiveContract(string)",
	"541e771d": "records(string)",
	"1e59c529": "register(string,address)",
}

// RegistryMockBin is the compiled bytecode used for deploying new contracts.
var RegistryMockBin = "0x608060405234801561001057600080fd5b506102d5806100206000396000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c8063080d5721146100465780631e59c52914610075578063541e771d1461008a575b600080fd5b6100596100543660046101d5565b6100be565b6040516001600160a01b03909116815260200160405180910390f35b610088610083366004610212565b6100ee565b005b6100596100983660046101d5565b80516020818301810180516000825292820191909301209152546001600160a01b031681565b600080826040516100cf9190610270565b908152604051908190036020019020546001600160a01b031692915050565b806000836040516100ff9190610270565b90815260405190819003602001902080546001600160a01b03929092166001600160a01b03199092169190911790555050565b634e487b7160e01b600052604160045260246000fd5b600082601f83011261015957600080fd5b813567ffffffffffffffff8082111561017457610174610132565b604051601f8301601f19908116603f0116810190828211818310171561019c5761019c610132565b816040528381528660208588010111156101b557600080fd5b836020870160208301376000602085830101528094505050505092915050565b6000602082840312156101e757600080fd5b813567ffffffffffffffff8111156101fe57600080fd5b61020a84828501610148565b949350505050565b6000806040838503121561022557600080fd5b823567ffffffffffffffff81111561023c57600080fd5b61024885828601610148565b92505060208301356001600160a01b038116811461026557600080fd5b809150509250929050565b6000825160005b818110156102915760208186018101518583015201610277565b50600092019182525091905056fea26469706673582212200be7666a50b9065e81c071c025510f3e57b6e8294690ac5940a1c711cf00340b64736f6c63430008110033"

// DeployRegistryMock deploys a new Klaytn contract, binding an instance of RegistryMock to it.
func DeployRegistryMock(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *RegistryMock, error) {
	parsed, err := abi.JSON(strings.NewReader(RegistryMockABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(RegistryMockBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &RegistryMock{RegistryMockCaller: RegistryMockCaller{contract: contract}, RegistryMockTransactor: RegistryMockTransactor{contract: contract}, RegistryMockFilterer: RegistryMockFilterer{contract: contract}}, nil
}

// RegistryMock is an auto generated Go binding around a Klaytn contract.
type RegistryMock struct {
	RegistryMockCaller     // Read-only binding to the contract
	RegistryMockTransactor // Write-only binding to the contract
	RegistryMockFilterer   // Log filterer for contract events
}

// RegistryMockCaller is an auto generated read-only Go binding around a Klaytn contract.
type RegistryMockCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RegistryMockTransactor is an auto generated write-only Go binding around a Klaytn contract.
type RegistryMockTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RegistryMockFilterer is an auto generated log filtering Go binding around a Klaytn contract events.
type RegistryMockFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RegistryMockSession is an auto generated Go binding around a Klaytn contract,
// with pre-set call and transact options.
type RegistryMockSession struct {
	Contract     *RegistryMock     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// RegistryMockCallerSession is an auto generated read-only Go binding around a Klaytn contract,
// with pre-set call options.
type RegistryMockCallerSession struct {
	Contract *RegistryMockCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// RegistryMockTransactorSession is an auto generated write-only Go binding around a Klaytn contract,
// with pre-set transact options.
type RegistryMockTransactorSession struct {
	Contract     *RegistryMockTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// RegistryMockRaw is an auto generated low-level Go binding around a Klaytn contract.
type RegistryMockRaw struct {
	Contract *RegistryMock // Generic contract binding to access the raw methods on
}

// RegistryMockCallerRaw is an auto generated low-level read-only Go binding around a Klaytn contract.
type RegistryMockCallerRaw struct {
	Contract *RegistryMockCaller // Generic read-only contract binding to access the raw methods on
}

// RegistryMockTransactorRaw is an auto generated low-level write-only Go binding around a Klaytn contract.
type RegistryMockTransactorRaw struct {
	Contract *RegistryMockTransactor // Generic write-only contract binding to access the raw methods on
}

// NewRegistryMock creates a new instance of RegistryMock, bound to a specific deployed contract.
func NewRegistryMock(address common.Address, backend bind.ContractBackend) (*RegistryMock, error) {
	contract, err := bindRegistryMock(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &RegistryMock{RegistryMockCaller: RegistryMockCaller{contract: contract}, RegistryMockTransactor: RegistryMockTransactor{contract: contract}, RegistryMockFilterer: RegistryMockFilterer{contract: contract}}, nil
}

// NewRegistryMockCaller creates a new read-only instance of RegistryMock, bound to a specific deployed contract.
func NewRegistryMockCaller(address common.Address, caller bind.ContractCaller) (*RegistryMockCaller, error) {
	contract, err := bindRegistryMock(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &RegistryMockCaller{contract: contract}, nil
}

// NewRegistryMockTransactor creates a new write-only instance of RegistryMock, bound to a specific deployed contract.
func NewRegistryMockTransactor(address common.Address, transactor bind.ContractTransactor) (*RegistryMockTransactor, error) {
	contract, err := bindRegistryMock(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &RegistryMockTransactor{contract: contract}, nil
}

// NewRegistryMockFilterer creates a new log filterer instance of RegistryMock, bound to a specific deployed contract.
func NewRegistryMockFilterer(address common.Address, filterer bind.ContractFilterer) (*RegistryMockFilterer, error) {
	contract, err := bindRegistryMock(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &RegistryMockFilterer{contract: contract}, nil
}

// bindRegistryMock binds a generic wrapper to an already deployed contract.
func bindRegistryMock(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(RegistryMockABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_RegistryMock *RegistryMockRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _RegistryMock.Contract.RegistryMockCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_RegistryMock *RegistryMockRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RegistryMock.Contract.RegistryMockTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_RegistryMock *RegistryMockRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RegistryMock.Contract.RegistryMockTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_RegistryMock *RegistryMockCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _RegistryMock.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_RegistryMock *RegistryMockTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RegistryMock.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_RegistryMock *RegistryMockTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RegistryMock.Contract.contract.Transact(opts, method, params...)
}

// GetActiveContract is a free data retrieval call binding the contract method 0x080d5721.
//
// Solidity: function getActiveContract(string name) view returns(address)
func (_RegistryMock *RegistryMockCaller) GetActiveContract(opts *bind.CallOpts, name string) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _RegistryMock.contract.Call(opts, out, "getActiveContract", name)
	return *ret0, err
}

// GetActiveContract is a free data retrieval call binding the contract method 0x080d5721.
//
// Solidity: function getActiveContract(string name) view returns(address)
func (_RegistryMock *RegistryMockSession) GetActiveContract(name string) (common.Address, error) {
	return _RegistryMock.Contract.GetActiveContract(&_RegistryMock.CallOpts, name)
}

// GetActiveContract is a free data retrieval call binding the contract method 0x080d5721.
//
// Solidity: function getActiveContract(string name) view returns(address)
func (_RegistryMock *RegistryMockCallerSession) GetActiveContract(name string) (common.Address, error) {
	return _RegistryMock.Contract.GetActiveContract(&_RegistryMock.CallOpts, name)
}

// Records is a free data retrieval call binding the contract method 0x541e771d.
//
// Solidity: function records(string ) view returns(address)
func (_RegistryMock *RegistryMockCaller) Records(opts *bind.CallOpts, arg0 string) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _RegistryMock.contract.Call(opts, out, "records", arg0)
	return *ret0, err
}

// Records is a free data retrieval call binding the contract method 0x541e771d.
//
// Solidity: function records(string ) view returns(address)
func (_RegistryMock *RegistryMockSession) Records(arg0 string) (common.Address, error) {
	return _RegistryMock.Contract.Records(&_RegistryMock.CallOpts, arg0)
}

// Records is a free data retrieval call binding the contract method 0x541e771d.
//
// Solidity: function records(string ) view returns(address)
func (_RegistryMock *RegistryMockCallerSession) Records(arg0 string) (common.Address, error) {
	return _RegistryMock.Contract.Records(&_RegistryMock.CallOpts, arg0)
}

// Register is a paid mutator transaction binding the contract method 0x1e59c529.
//
// Solidity: function register(string name, address addr) returns()
func (_RegistryMock *RegistryMockTransactor) Register(opts *bind.TransactOpts, name string, addr common.Address) (*types.Transaction, error) {
	return _RegistryMock.contract.Transact(opts, "register", name, addr)
}

// Register is a paid mutator transaction binding the contract method 0x1e59c529.
//
// Solidity: function register(string name, address addr) returns()
func (_RegistryMock *RegistryMockSession) Register(name string, addr common.Address) (*types.Transaction, error) {
	return _RegistryMock.Contract.Register(&_RegistryMock.TransactOpts, name, addr)
}

// Register is a paid mutator transaction binding the contract method 0x1e59c529.
//
// Solidity: function register(string name, address addr) returns()
func (_RegistryMock *RegistryMockTransactorSession) Register(name string, addr common.Address) (*types.Transaction, error) {
	return _RegistryMock.Contract.Register(&_RegistryMock.TransactOpts, name, addr)
}
