// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package system_contracts

import (
	"errors"
	"math/big"
	"strings"

	"github.com/klaytn/klaytn"
	"github.com/klaytn/klaytn/accounts/abi/bind"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = klaytn.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// IRegistryMetaData contains all meta data concerning the IRegistry contract.
var IRegistryMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"name\":\"getActiveAddr\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"names\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"records\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"activation\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"e2693e3f": "getActiveAddr(string)",
		"4622ab03": "names(uint256)",
		"3b51650d": "records(string,uint256)",
	},
}

// IRegistryABI is the input ABI used to generate the binding from.
// Deprecated: Use IRegistryMetaData.ABI instead.
var IRegistryABI = IRegistryMetaData.ABI

// IRegistryBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const IRegistryBinRuntime = ``

// IRegistryFuncSigs maps the 4-byte function signature to its string representation.
// Deprecated: Use IRegistryMetaData.Sigs instead.
var IRegistryFuncSigs = IRegistryMetaData.Sigs

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
	parsed, err := IRegistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
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

// GetActiveAddr is a free data retrieval call binding the contract method 0xe2693e3f.
//
// Solidity: function getActiveAddr(string name) view returns(address)
func (_IRegistry *IRegistryCaller) GetActiveAddr(opts *bind.CallOpts, name string) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _IRegistry.contract.Call(opts, out, "getActiveAddr", name)
	return *ret0, err
}

// GetActiveAddr is a free data retrieval call binding the contract method 0xe2693e3f.
//
// Solidity: function getActiveAddr(string name) view returns(address)
func (_IRegistry *IRegistrySession) GetActiveAddr(name string) (common.Address, error) {
	return _IRegistry.Contract.GetActiveAddr(&_IRegistry.CallOpts, name)
}

// GetActiveAddr is a free data retrieval call binding the contract method 0xe2693e3f.
//
// Solidity: function getActiveAddr(string name) view returns(address)
func (_IRegistry *IRegistryCallerSession) GetActiveAddr(name string) (common.Address, error) {
	return _IRegistry.Contract.GetActiveAddr(&_IRegistry.CallOpts, name)
}

// Names is a free data retrieval call binding the contract method 0x4622ab03.
//
// Solidity: function names(uint256 ) view returns(string)
func (_IRegistry *IRegistryCaller) Names(opts *bind.CallOpts, arg0 *big.Int) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _IRegistry.contract.Call(opts, out, "names", arg0)
	return *ret0, err
}

// Names is a free data retrieval call binding the contract method 0x4622ab03.
//
// Solidity: function names(uint256 ) view returns(string)
func (_IRegistry *IRegistrySession) Names(arg0 *big.Int) (string, error) {
	return _IRegistry.Contract.Names(&_IRegistry.CallOpts, arg0)
}

// Names is a free data retrieval call binding the contract method 0x4622ab03.
//
// Solidity: function names(uint256 ) view returns(string)
func (_IRegistry *IRegistryCallerSession) Names(arg0 *big.Int) (string, error) {
	return _IRegistry.Contract.Names(&_IRegistry.CallOpts, arg0)
}

// Records is a free data retrieval call binding the contract method 0x3b51650d.
//
// Solidity: function records(string , uint256 ) view returns(address addr, uint256 activation)
func (_IRegistry *IRegistryCaller) Records(opts *bind.CallOpts, arg0 string, arg1 *big.Int) (struct {
	Addr       common.Address
	Activation *big.Int
}, error) {
	ret := new(struct {
		Addr       common.Address
		Activation *big.Int
	})
	out := ret
	err := _IRegistry.contract.Call(opts, out, "records", arg0, arg1)
	return *ret, err
}

// Records is a free data retrieval call binding the contract method 0x3b51650d.
//
// Solidity: function records(string , uint256 ) view returns(address addr, uint256 activation)
func (_IRegistry *IRegistrySession) Records(arg0 string, arg1 *big.Int) (struct {
	Addr       common.Address
	Activation *big.Int
}, error) {
	return _IRegistry.Contract.Records(&_IRegistry.CallOpts, arg0, arg1)
}

// Records is a free data retrieval call binding the contract method 0x3b51650d.
//
// Solidity: function records(string , uint256 ) view returns(address addr, uint256 activation)
func (_IRegistry *IRegistryCallerSession) Records(arg0 string, arg1 *big.Int) (struct {
	Addr       common.Address
	Activation *big.Int
}, error) {
	return _IRegistry.Contract.Records(&_IRegistry.CallOpts, arg0, arg1)
}

// RegistryMetaData contains all meta data concerning the Registry contract.
var RegistryMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"name\":\"getActiveAddr\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"names\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"records\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"activation\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"e2693e3f": "getActiveAddr(string)",
		"4622ab03": "names(uint256)",
		"3b51650d": "records(string,uint256)",
	},
	Bin: "0x608060405234801561001057600080fd5b5061040e806100206000396000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c80633b51650d146100465780634622ab031461007d578063e2693e3f1461009d575b600080fd5b6100596100543660046102ad565b6100c8565b604080516001600160a01b0390931683526020830191909152015b60405180910390f35b61009061008b3660046102f2565b61011d565b604051610074919061030b565b6100b06100ab366004610360565b6101c9565b6040516001600160a01b039091168152602001610074565b815160208184018101805160008252928201918501919091209190528054829081106100f357600080fd5b6000918252602090912060029091020180546001909101546001600160a01b039091169250905082565b6001818154811061012d57600080fd5b9060005260206000200160009150905080546101489061039d565b80601f01602080910402602001604051908101604052809291908181526020018280546101749061039d565b80156101c15780601f10610196576101008083540402835291602001916101c1565b820191906000526020600020905b8154815290600101906020018083116101a457829003601f168201915b505050505081565b60405162461bcd60e51b815260206004820152600f60248201526e1b9bdd081a5b5c1b195b595b9d1959608a1b604482015260009060640160405180910390fd5b634e487b7160e01b600052604160045260246000fd5b600082601f83011261023157600080fd5b813567ffffffffffffffff8082111561024c5761024c61020a565b604051601f8301601f19908116603f011681019082821181831017156102745761027461020a565b8160405283815286602085880101111561028d57600080fd5b836020870160208301376000602085830101528094505050505092915050565b600080604083850312156102c057600080fd5b823567ffffffffffffffff8111156102d757600080fd5b6102e385828601610220565b95602094909401359450505050565b60006020828403121561030457600080fd5b5035919050565b600060208083528351808285015260005b818110156103385785810183015185820160400152820161031c565b8181111561034a576000604083870101525b50601f01601f1916929092016040019392505050565b60006020828403121561037257600080fd5b813567ffffffffffffffff81111561038957600080fd5b61039584828501610220565b949350505050565b600181811c908216806103b157607f821691505b602082108114156103d257634e487b7160e01b600052602260045260246000fd5b5091905056fea2646970667358221220e26a575b63762d698eb6e82f6fd9cac682122d4874acca851d6e8e9f20ac556d64736f6c634300080b0033",
}

// RegistryABI is the input ABI used to generate the binding from.
// Deprecated: Use RegistryMetaData.ABI instead.
var RegistryABI = RegistryMetaData.ABI

// RegistryBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const RegistryBinRuntime = `608060405234801561001057600080fd5b50600436106100415760003560e01c80633b51650d146100465780634622ab031461007d578063e2693e3f1461009d575b600080fd5b6100596100543660046102ad565b6100c8565b604080516001600160a01b0390931683526020830191909152015b60405180910390f35b61009061008b3660046102f2565b61011d565b604051610074919061030b565b6100b06100ab366004610360565b6101c9565b6040516001600160a01b039091168152602001610074565b815160208184018101805160008252928201918501919091209190528054829081106100f357600080fd5b6000918252602090912060029091020180546001909101546001600160a01b039091169250905082565b6001818154811061012d57600080fd5b9060005260206000200160009150905080546101489061039d565b80601f01602080910402602001604051908101604052809291908181526020018280546101749061039d565b80156101c15780601f10610196576101008083540402835291602001916101c1565b820191906000526020600020905b8154815290600101906020018083116101a457829003601f168201915b505050505081565b60405162461bcd60e51b815260206004820152600f60248201526e1b9bdd081a5b5c1b195b595b9d1959608a1b604482015260009060640160405180910390fd5b634e487b7160e01b600052604160045260246000fd5b600082601f83011261023157600080fd5b813567ffffffffffffffff8082111561024c5761024c61020a565b604051601f8301601f19908116603f011681019082821181831017156102745761027461020a565b8160405283815286602085880101111561028d57600080fd5b836020870160208301376000602085830101528094505050505092915050565b600080604083850312156102c057600080fd5b823567ffffffffffffffff8111156102d757600080fd5b6102e385828601610220565b95602094909401359450505050565b60006020828403121561030457600080fd5b5035919050565b600060208083528351808285015260005b818110156103385785810183015185820160400152820161031c565b8181111561034a576000604083870101525b50601f01601f1916929092016040019392505050565b60006020828403121561037257600080fd5b813567ffffffffffffffff81111561038957600080fd5b61039584828501610220565b949350505050565b600181811c908216806103b157607f821691505b602082108114156103d257634e487b7160e01b600052602260045260246000fd5b5091905056fea2646970667358221220e26a575b63762d698eb6e82f6fd9cac682122d4874acca851d6e8e9f20ac556d64736f6c634300080b0033`

// RegistryFuncSigs maps the 4-byte function signature to its string representation.
// Deprecated: Use RegistryMetaData.Sigs instead.
var RegistryFuncSigs = RegistryMetaData.Sigs

// RegistryBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use RegistryMetaData.Bin instead.
var RegistryBin = RegistryMetaData.Bin

// DeployRegistry deploys a new Klaytn contract, binding an instance of Registry to it.
func DeployRegistry(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Registry, error) {
	parsed, err := RegistryMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(RegistryBin), backend)
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
	parsed, err := RegistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
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

// GetActiveAddr is a free data retrieval call binding the contract method 0xe2693e3f.
//
// Solidity: function getActiveAddr(string name) view returns(address)
func (_Registry *RegistryCaller) GetActiveAddr(opts *bind.CallOpts, name string) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Registry.contract.Call(opts, out, "getActiveAddr", name)
	return *ret0, err
}

// GetActiveAddr is a free data retrieval call binding the contract method 0xe2693e3f.
//
// Solidity: function getActiveAddr(string name) view returns(address)
func (_Registry *RegistrySession) GetActiveAddr(name string) (common.Address, error) {
	return _Registry.Contract.GetActiveAddr(&_Registry.CallOpts, name)
}

// GetActiveAddr is a free data retrieval call binding the contract method 0xe2693e3f.
//
// Solidity: function getActiveAddr(string name) view returns(address)
func (_Registry *RegistryCallerSession) GetActiveAddr(name string) (common.Address, error) {
	return _Registry.Contract.GetActiveAddr(&_Registry.CallOpts, name)
}

// Names is a free data retrieval call binding the contract method 0x4622ab03.
//
// Solidity: function names(uint256 ) view returns(string)
func (_Registry *RegistryCaller) Names(opts *bind.CallOpts, arg0 *big.Int) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _Registry.contract.Call(opts, out, "names", arg0)
	return *ret0, err
}

// Names is a free data retrieval call binding the contract method 0x4622ab03.
//
// Solidity: function names(uint256 ) view returns(string)
func (_Registry *RegistrySession) Names(arg0 *big.Int) (string, error) {
	return _Registry.Contract.Names(&_Registry.CallOpts, arg0)
}

// Names is a free data retrieval call binding the contract method 0x4622ab03.
//
// Solidity: function names(uint256 ) view returns(string)
func (_Registry *RegistryCallerSession) Names(arg0 *big.Int) (string, error) {
	return _Registry.Contract.Names(&_Registry.CallOpts, arg0)
}

// Records is a free data retrieval call binding the contract method 0x3b51650d.
//
// Solidity: function records(string , uint256 ) view returns(address addr, uint256 activation)
func (_Registry *RegistryCaller) Records(opts *bind.CallOpts, arg0 string, arg1 *big.Int) (struct {
	Addr       common.Address
	Activation *big.Int
}, error) {
	ret := new(struct {
		Addr       common.Address
		Activation *big.Int
	})
	out := ret
	err := _Registry.contract.Call(opts, out, "records", arg0, arg1)
	return *ret, err
}

// Records is a free data retrieval call binding the contract method 0x3b51650d.
//
// Solidity: function records(string , uint256 ) view returns(address addr, uint256 activation)
func (_Registry *RegistrySession) Records(arg0 string, arg1 *big.Int) (struct {
	Addr       common.Address
	Activation *big.Int
}, error) {
	return _Registry.Contract.Records(&_Registry.CallOpts, arg0, arg1)
}

// Records is a free data retrieval call binding the contract method 0x3b51650d.
//
// Solidity: function records(string , uint256 ) view returns(address addr, uint256 activation)
func (_Registry *RegistryCallerSession) Records(arg0 string, arg1 *big.Int) (struct {
	Addr       common.Address
	Activation *big.Int
}, error) {
	return _Registry.Contract.Records(&_Registry.CallOpts, arg0, arg1)
}

// RegistryMockMetaData contains all meta data concerning the RegistryMock contract.
var RegistryMockMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"name\":\"getActiveAddr\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"names\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"records\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"activation\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"activation\",\"type\":\"uint256\"}],\"name\":\"register\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"e2693e3f": "getActiveAddr(string)",
		"4622ab03": "names(uint256)",
		"3b51650d": "records(string,uint256)",
		"d393c871": "register(string,address,uint256)",
		"f2fde38b": "transferOwnership(address)",
	},
	Bin: "0x608060405234801561001057600080fd5b50610720806100206000396000f3fe608060405234801561001057600080fd5b50600436106100575760003560e01c80633b51650d1461005c5780634622ab0314610093578063d393c871146100b3578063e2693e3f146100c8578063f2fde38b146100f3575b600080fd5b61006f61006a3660046104cb565b610123565b604080516001600160a01b0390931683526020830191909152015b60405180910390f35b6100a66100a1366004610510565b610178565b60405161008a9190610559565b6100c66100c13660046105a8565b610224565b005b6100db6100d63660046105ff565b6102fb565b6040516001600160a01b03909116815260200161008a565b6100c661010136600461063c565b600280546001600160a01b0319166001600160a01b0392909216919091179055565b8151602081840181018051600082529282019185019190912091905280548290811061014e57600080fd5b6000918252602090912060029091020180546001909101546001600160a01b039091169250905082565b6001818154811061018857600080fd5b9060005260206000200160009150905080546101a39061065e565b80601f01602080910402602001604051908101604052809291908181526020018280546101cf9061065e565b801561021c5780601f106101f15761010080835404028352916020019161021c565b820191906000526020600020905b8154815290600101906020018083116101ff57829003601f168201915b505050505081565b6000836040516102349190610693565b9081526040519081900360200190205461028c57600180548082018255600091909152835161028a917fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf60190602086019061038f565b505b60008360405161029c9190610693565b90815260408051602092819003830181208183019092526001600160a01b039485168152828101938452815460018082018455600093845293909220905160029092020180546001600160a01b03191691909416178355905191015550565b60008060008360405161030e9190610693565b9081526040519081900360200190205490508061032e5750600092915050565b60008360405161033e9190610693565b9081526040519081900360200190206103586001836106af565b81548110610368576103686106d4565b60009182526020909120600290910201546001600160a01b03169392505050565b50919050565b82805461039b9061065e565b90600052602060002090601f0160209004810192826103bd5760008555610403565b82601f106103d657805160ff1916838001178555610403565b82800160010185558215610403579182015b828111156104035782518255916020019190600101906103e8565b5061040f929150610413565b5090565b5b8082111561040f5760008155600101610414565b634e487b7160e01b600052604160045260246000fd5b600082601f83011261044f57600080fd5b813567ffffffffffffffff8082111561046a5761046a610428565b604051601f8301601f19908116603f0116810190828211818310171561049257610492610428565b816040528381528660208588010111156104ab57600080fd5b836020870160208301376000602085830101528094505050505092915050565b600080604083850312156104de57600080fd5b823567ffffffffffffffff8111156104f557600080fd5b6105018582860161043e565b95602094909401359450505050565b60006020828403121561052257600080fd5b5035919050565b60005b8381101561054457818101518382015260200161052c565b83811115610553576000848401525b50505050565b6020815260008251806020840152610578816040850160208701610529565b601f01601f19169190910160400192915050565b80356001600160a01b03811681146105a357600080fd5b919050565b6000806000606084860312156105bd57600080fd5b833567ffffffffffffffff8111156105d457600080fd5b6105e08682870161043e565b9350506105ef6020850161058c565b9150604084013590509250925092565b60006020828403121561061157600080fd5b813567ffffffffffffffff81111561062857600080fd5b6106348482850161043e565b949350505050565b60006020828403121561064e57600080fd5b6106578261058c565b9392505050565b600181811c9082168061067257607f821691505b6020821081141561038957634e487b7160e01b600052602260045260246000fd5b600082516106a5818460208701610529565b9190910192915050565b6000828210156106cf57634e487b7160e01b600052601160045260246000fd5b500390565b634e487b7160e01b600052603260045260246000fdfea26469706673582212207c642130a3785d0c34f24c86b872c3271bfc8012eb6d21a3a27b0b1ec3b2660c64736f6c634300080b0033",
}

// RegistryMockABI is the input ABI used to generate the binding from.
// Deprecated: Use RegistryMockMetaData.ABI instead.
var RegistryMockABI = RegistryMockMetaData.ABI

// RegistryMockBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const RegistryMockBinRuntime = `608060405234801561001057600080fd5b50600436106100575760003560e01c80633b51650d1461005c5780634622ab0314610093578063d393c871146100b3578063e2693e3f146100c8578063f2fde38b146100f3575b600080fd5b61006f61006a3660046104cb565b610123565b604080516001600160a01b0390931683526020830191909152015b60405180910390f35b6100a66100a1366004610510565b610178565b60405161008a9190610559565b6100c66100c13660046105a8565b610224565b005b6100db6100d63660046105ff565b6102fb565b6040516001600160a01b03909116815260200161008a565b6100c661010136600461063c565b600280546001600160a01b0319166001600160a01b0392909216919091179055565b8151602081840181018051600082529282019185019190912091905280548290811061014e57600080fd5b6000918252602090912060029091020180546001909101546001600160a01b039091169250905082565b6001818154811061018857600080fd5b9060005260206000200160009150905080546101a39061065e565b80601f01602080910402602001604051908101604052809291908181526020018280546101cf9061065e565b801561021c5780601f106101f15761010080835404028352916020019161021c565b820191906000526020600020905b8154815290600101906020018083116101ff57829003601f168201915b505050505081565b6000836040516102349190610693565b9081526040519081900360200190205461028c57600180548082018255600091909152835161028a917fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf60190602086019061038f565b505b60008360405161029c9190610693565b90815260408051602092819003830181208183019092526001600160a01b039485168152828101938452815460018082018455600093845293909220905160029092020180546001600160a01b03191691909416178355905191015550565b60008060008360405161030e9190610693565b9081526040519081900360200190205490508061032e5750600092915050565b60008360405161033e9190610693565b9081526040519081900360200190206103586001836106af565b81548110610368576103686106d4565b60009182526020909120600290910201546001600160a01b03169392505050565b50919050565b82805461039b9061065e565b90600052602060002090601f0160209004810192826103bd5760008555610403565b82601f106103d657805160ff1916838001178555610403565b82800160010185558215610403579182015b828111156104035782518255916020019190600101906103e8565b5061040f929150610413565b5090565b5b8082111561040f5760008155600101610414565b634e487b7160e01b600052604160045260246000fd5b600082601f83011261044f57600080fd5b813567ffffffffffffffff8082111561046a5761046a610428565b604051601f8301601f19908116603f0116810190828211818310171561049257610492610428565b816040528381528660208588010111156104ab57600080fd5b836020870160208301376000602085830101528094505050505092915050565b600080604083850312156104de57600080fd5b823567ffffffffffffffff8111156104f557600080fd5b6105018582860161043e565b95602094909401359450505050565b60006020828403121561052257600080fd5b5035919050565b60005b8381101561054457818101518382015260200161052c565b83811115610553576000848401525b50505050565b6020815260008251806020840152610578816040850160208701610529565b601f01601f19169190910160400192915050565b80356001600160a01b03811681146105a357600080fd5b919050565b6000806000606084860312156105bd57600080fd5b833567ffffffffffffffff8111156105d457600080fd5b6105e08682870161043e565b9350506105ef6020850161058c565b9150604084013590509250925092565b60006020828403121561061157600080fd5b813567ffffffffffffffff81111561062857600080fd5b6106348482850161043e565b949350505050565b60006020828403121561064e57600080fd5b6106578261058c565b9392505050565b600181811c9082168061067257607f821691505b6020821081141561038957634e487b7160e01b600052602260045260246000fd5b600082516106a5818460208701610529565b9190910192915050565b6000828210156106cf57634e487b7160e01b600052601160045260246000fd5b500390565b634e487b7160e01b600052603260045260246000fdfea26469706673582212207c642130a3785d0c34f24c86b872c3271bfc8012eb6d21a3a27b0b1ec3b2660c64736f6c634300080b0033`

// RegistryMockFuncSigs maps the 4-byte function signature to its string representation.
// Deprecated: Use RegistryMockMetaData.Sigs instead.
var RegistryMockFuncSigs = RegistryMockMetaData.Sigs

// RegistryMockBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use RegistryMockMetaData.Bin instead.
var RegistryMockBin = RegistryMockMetaData.Bin

// DeployRegistryMock deploys a new Klaytn contract, binding an instance of RegistryMock to it.
func DeployRegistryMock(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *RegistryMock, error) {
	parsed, err := RegistryMockMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(RegistryMockBin), backend)
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
	parsed, err := RegistryMockMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
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

// GetActiveAddr is a free data retrieval call binding the contract method 0xe2693e3f.
//
// Solidity: function getActiveAddr(string name) view returns(address)
func (_RegistryMock *RegistryMockCaller) GetActiveAddr(opts *bind.CallOpts, name string) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _RegistryMock.contract.Call(opts, out, "getActiveAddr", name)
	return *ret0, err
}

// GetActiveAddr is a free data retrieval call binding the contract method 0xe2693e3f.
//
// Solidity: function getActiveAddr(string name) view returns(address)
func (_RegistryMock *RegistryMockSession) GetActiveAddr(name string) (common.Address, error) {
	return _RegistryMock.Contract.GetActiveAddr(&_RegistryMock.CallOpts, name)
}

// GetActiveAddr is a free data retrieval call binding the contract method 0xe2693e3f.
//
// Solidity: function getActiveAddr(string name) view returns(address)
func (_RegistryMock *RegistryMockCallerSession) GetActiveAddr(name string) (common.Address, error) {
	return _RegistryMock.Contract.GetActiveAddr(&_RegistryMock.CallOpts, name)
}

// Names is a free data retrieval call binding the contract method 0x4622ab03.
//
// Solidity: function names(uint256 ) view returns(string)
func (_RegistryMock *RegistryMockCaller) Names(opts *bind.CallOpts, arg0 *big.Int) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _RegistryMock.contract.Call(opts, out, "names", arg0)
	return *ret0, err
}

// Names is a free data retrieval call binding the contract method 0x4622ab03.
//
// Solidity: function names(uint256 ) view returns(string)
func (_RegistryMock *RegistryMockSession) Names(arg0 *big.Int) (string, error) {
	return _RegistryMock.Contract.Names(&_RegistryMock.CallOpts, arg0)
}

// Names is a free data retrieval call binding the contract method 0x4622ab03.
//
// Solidity: function names(uint256 ) view returns(string)
func (_RegistryMock *RegistryMockCallerSession) Names(arg0 *big.Int) (string, error) {
	return _RegistryMock.Contract.Names(&_RegistryMock.CallOpts, arg0)
}

// Records is a free data retrieval call binding the contract method 0x3b51650d.
//
// Solidity: function records(string , uint256 ) view returns(address addr, uint256 activation)
func (_RegistryMock *RegistryMockCaller) Records(opts *bind.CallOpts, arg0 string, arg1 *big.Int) (struct {
	Addr       common.Address
	Activation *big.Int
}, error) {
	ret := new(struct {
		Addr       common.Address
		Activation *big.Int
	})
	out := ret
	err := _RegistryMock.contract.Call(opts, out, "records", arg0, arg1)
	return *ret, err
}

// Records is a free data retrieval call binding the contract method 0x3b51650d.
//
// Solidity: function records(string , uint256 ) view returns(address addr, uint256 activation)
func (_RegistryMock *RegistryMockSession) Records(arg0 string, arg1 *big.Int) (struct {
	Addr       common.Address
	Activation *big.Int
}, error) {
	return _RegistryMock.Contract.Records(&_RegistryMock.CallOpts, arg0, arg1)
}

// Records is a free data retrieval call binding the contract method 0x3b51650d.
//
// Solidity: function records(string , uint256 ) view returns(address addr, uint256 activation)
func (_RegistryMock *RegistryMockCallerSession) Records(arg0 string, arg1 *big.Int) (struct {
	Addr       common.Address
	Activation *big.Int
}, error) {
	return _RegistryMock.Contract.Records(&_RegistryMock.CallOpts, arg0, arg1)
}

// Register is a paid mutator transaction binding the contract method 0xd393c871.
//
// Solidity: function register(string name, address addr, uint256 activation) returns()
func (_RegistryMock *RegistryMockTransactor) Register(opts *bind.TransactOpts, name string, addr common.Address, activation *big.Int) (*types.Transaction, error) {
	return _RegistryMock.contract.Transact(opts, "register", name, addr, activation)
}

// Register is a paid mutator transaction binding the contract method 0xd393c871.
//
// Solidity: function register(string name, address addr, uint256 activation) returns()
func (_RegistryMock *RegistryMockSession) Register(name string, addr common.Address, activation *big.Int) (*types.Transaction, error) {
	return _RegistryMock.Contract.Register(&_RegistryMock.TransactOpts, name, addr, activation)
}

// Register is a paid mutator transaction binding the contract method 0xd393c871.
//
// Solidity: function register(string name, address addr, uint256 activation) returns()
func (_RegistryMock *RegistryMockTransactorSession) Register(name string, addr common.Address, activation *big.Int) (*types.Transaction, error) {
	return _RegistryMock.Contract.Register(&_RegistryMock.TransactOpts, name, addr, activation)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_RegistryMock *RegistryMockTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _RegistryMock.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_RegistryMock *RegistryMockSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _RegistryMock.Contract.TransferOwnership(&_RegistryMock.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_RegistryMock *RegistryMockTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _RegistryMock.Contract.TransferOwnership(&_RegistryMock.TransactOpts, newOwner)
}
