// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package system_contracts

import (
	"errors"
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
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = klaytn.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// IKIP113BlsPublicKeyInfo is an auto generated low-level Go binding around an user-defined struct.
type IKIP113BlsPublicKeyInfo struct {
	PublicKey []byte
	Pop       []byte
}

// IRegistryRecord is an auto generated low-level Go binding around an user-defined struct.
type IRegistryRecord struct {
	Addr       common.Address
	Activation *big.Int
}

// AddressMetaData contains all meta data concerning the Address contract.
var AddressMetaData = &bind.MetaData{
	ABI: "[]",
	Bin: "0x60566037600b82828239805160001a607314602a57634e487b7160e01b600052600060045260246000fd5b30600052607381538281f3fe73000000000000000000000000000000000000000030146080604052600080fdfea264697066735822122053101f7ec578e063c340643045a810e3dbc9e3321e7a485ecbf68e707cee458364736f6c63430008130033",
}

// AddressABI is the input ABI used to generate the binding from.
// Deprecated: Use AddressMetaData.ABI instead.
var AddressABI = AddressMetaData.ABI

// AddressBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const AddressBinRuntime = `73000000000000000000000000000000000000000030146080604052600080fdfea264697066735822122053101f7ec578e063c340643045a810e3dbc9e3321e7a485ecbf68e707cee458364736f6c63430008130033`

// AddressBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use AddressMetaData.Bin instead.
var AddressBin = AddressMetaData.Bin

// DeployAddress deploys a new Klaytn contract, binding an instance of Address to it.
func DeployAddress(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Address, error) {
	parsed, err := AddressMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(AddressBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Address{AddressCaller: AddressCaller{contract: contract}, AddressTransactor: AddressTransactor{contract: contract}, AddressFilterer: AddressFilterer{contract: contract}}, nil
}

// Address is an auto generated Go binding around a Klaytn contract.
type Address struct {
	AddressCaller     // Read-only binding to the contract
	AddressTransactor // Write-only binding to the contract
	AddressFilterer   // Log filterer for contract events
}

// AddressCaller is an auto generated read-only Go binding around a Klaytn contract.
type AddressCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AddressTransactor is an auto generated write-only Go binding around a Klaytn contract.
type AddressTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AddressFilterer is an auto generated log filtering Go binding around a Klaytn contract events.
type AddressFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AddressSession is an auto generated Go binding around a Klaytn contract,
// with pre-set call and transact options.
type AddressSession struct {
	Contract     *Address          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// AddressCallerSession is an auto generated read-only Go binding around a Klaytn contract,
// with pre-set call options.
type AddressCallerSession struct {
	Contract *AddressCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// AddressTransactorSession is an auto generated write-only Go binding around a Klaytn contract,
// with pre-set transact options.
type AddressTransactorSession struct {
	Contract     *AddressTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// AddressRaw is an auto generated low-level Go binding around a Klaytn contract.
type AddressRaw struct {
	Contract *Address // Generic contract binding to access the raw methods on
}

// AddressCallerRaw is an auto generated low-level read-only Go binding around a Klaytn contract.
type AddressCallerRaw struct {
	Contract *AddressCaller // Generic read-only contract binding to access the raw methods on
}

// AddressTransactorRaw is an auto generated low-level write-only Go binding around a Klaytn contract.
type AddressTransactorRaw struct {
	Contract *AddressTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAddress creates a new instance of Address, bound to a specific deployed contract.
func NewAddress(address common.Address, backend bind.ContractBackend) (*Address, error) {
	contract, err := bindAddress(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Address{AddressCaller: AddressCaller{contract: contract}, AddressTransactor: AddressTransactor{contract: contract}, AddressFilterer: AddressFilterer{contract: contract}}, nil
}

// NewAddressCaller creates a new read-only instance of Address, bound to a specific deployed contract.
func NewAddressCaller(address common.Address, caller bind.ContractCaller) (*AddressCaller, error) {
	contract, err := bindAddress(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AddressCaller{contract: contract}, nil
}

// NewAddressTransactor creates a new write-only instance of Address, bound to a specific deployed contract.
func NewAddressTransactor(address common.Address, transactor bind.ContractTransactor) (*AddressTransactor, error) {
	contract, err := bindAddress(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AddressTransactor{contract: contract}, nil
}

// NewAddressFilterer creates a new log filterer instance of Address, bound to a specific deployed contract.
func NewAddressFilterer(address common.Address, filterer bind.ContractFilterer) (*AddressFilterer, error) {
	contract, err := bindAddress(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AddressFilterer{contract: contract}, nil
}

// bindAddress binds a generic wrapper to an already deployed contract.
func bindAddress(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := AddressMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Address *AddressRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Address.Contract.AddressCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Address *AddressRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Address.Contract.AddressTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Address *AddressRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Address.Contract.AddressTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Address *AddressCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Address.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Address *AddressTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Address.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Address *AddressTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Address.Contract.contract.Transact(opts, method, params...)
}

// AddressUpgradeableMetaData contains all meta data concerning the AddressUpgradeable contract.
var AddressUpgradeableMetaData = &bind.MetaData{
	ABI: "[]",
	Bin: "0x60566037600b82828239805160001a607314602a57634e487b7160e01b600052600060045260246000fd5b30600052607381538281f3fe73000000000000000000000000000000000000000030146080604052600080fdfea26469706673582212209bf91cc68b5d2a0488d512025c875a923a298928d8f1d2cef3996b4dc164352c64736f6c63430008130033",
}

// AddressUpgradeableABI is the input ABI used to generate the binding from.
// Deprecated: Use AddressUpgradeableMetaData.ABI instead.
var AddressUpgradeableABI = AddressUpgradeableMetaData.ABI

// AddressUpgradeableBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const AddressUpgradeableBinRuntime = `73000000000000000000000000000000000000000030146080604052600080fdfea26469706673582212209bf91cc68b5d2a0488d512025c875a923a298928d8f1d2cef3996b4dc164352c64736f6c63430008130033`

// AddressUpgradeableBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use AddressUpgradeableMetaData.Bin instead.
var AddressUpgradeableBin = AddressUpgradeableMetaData.Bin

// DeployAddressUpgradeable deploys a new Klaytn contract, binding an instance of AddressUpgradeable to it.
func DeployAddressUpgradeable(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *AddressUpgradeable, error) {
	parsed, err := AddressUpgradeableMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(AddressUpgradeableBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &AddressUpgradeable{AddressUpgradeableCaller: AddressUpgradeableCaller{contract: contract}, AddressUpgradeableTransactor: AddressUpgradeableTransactor{contract: contract}, AddressUpgradeableFilterer: AddressUpgradeableFilterer{contract: contract}}, nil
}

// AddressUpgradeable is an auto generated Go binding around a Klaytn contract.
type AddressUpgradeable struct {
	AddressUpgradeableCaller     // Read-only binding to the contract
	AddressUpgradeableTransactor // Write-only binding to the contract
	AddressUpgradeableFilterer   // Log filterer for contract events
}

// AddressUpgradeableCaller is an auto generated read-only Go binding around a Klaytn contract.
type AddressUpgradeableCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AddressUpgradeableTransactor is an auto generated write-only Go binding around a Klaytn contract.
type AddressUpgradeableTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AddressUpgradeableFilterer is an auto generated log filtering Go binding around a Klaytn contract events.
type AddressUpgradeableFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AddressUpgradeableSession is an auto generated Go binding around a Klaytn contract,
// with pre-set call and transact options.
type AddressUpgradeableSession struct {
	Contract     *AddressUpgradeable // Generic contract binding to set the session for
	CallOpts     bind.CallOpts       // Call options to use throughout this session
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// AddressUpgradeableCallerSession is an auto generated read-only Go binding around a Klaytn contract,
// with pre-set call options.
type AddressUpgradeableCallerSession struct {
	Contract *AddressUpgradeableCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts             // Call options to use throughout this session
}

// AddressUpgradeableTransactorSession is an auto generated write-only Go binding around a Klaytn contract,
// with pre-set transact options.
type AddressUpgradeableTransactorSession struct {
	Contract     *AddressUpgradeableTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts             // Transaction auth options to use throughout this session
}

// AddressUpgradeableRaw is an auto generated low-level Go binding around a Klaytn contract.
type AddressUpgradeableRaw struct {
	Contract *AddressUpgradeable // Generic contract binding to access the raw methods on
}

// AddressUpgradeableCallerRaw is an auto generated low-level read-only Go binding around a Klaytn contract.
type AddressUpgradeableCallerRaw struct {
	Contract *AddressUpgradeableCaller // Generic read-only contract binding to access the raw methods on
}

// AddressUpgradeableTransactorRaw is an auto generated low-level write-only Go binding around a Klaytn contract.
type AddressUpgradeableTransactorRaw struct {
	Contract *AddressUpgradeableTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAddressUpgradeable creates a new instance of AddressUpgradeable, bound to a specific deployed contract.
func NewAddressUpgradeable(address common.Address, backend bind.ContractBackend) (*AddressUpgradeable, error) {
	contract, err := bindAddressUpgradeable(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AddressUpgradeable{AddressUpgradeableCaller: AddressUpgradeableCaller{contract: contract}, AddressUpgradeableTransactor: AddressUpgradeableTransactor{contract: contract}, AddressUpgradeableFilterer: AddressUpgradeableFilterer{contract: contract}}, nil
}

// NewAddressUpgradeableCaller creates a new read-only instance of AddressUpgradeable, bound to a specific deployed contract.
func NewAddressUpgradeableCaller(address common.Address, caller bind.ContractCaller) (*AddressUpgradeableCaller, error) {
	contract, err := bindAddressUpgradeable(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AddressUpgradeableCaller{contract: contract}, nil
}

// NewAddressUpgradeableTransactor creates a new write-only instance of AddressUpgradeable, bound to a specific deployed contract.
func NewAddressUpgradeableTransactor(address common.Address, transactor bind.ContractTransactor) (*AddressUpgradeableTransactor, error) {
	contract, err := bindAddressUpgradeable(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AddressUpgradeableTransactor{contract: contract}, nil
}

// NewAddressUpgradeableFilterer creates a new log filterer instance of AddressUpgradeable, bound to a specific deployed contract.
func NewAddressUpgradeableFilterer(address common.Address, filterer bind.ContractFilterer) (*AddressUpgradeableFilterer, error) {
	contract, err := bindAddressUpgradeable(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AddressUpgradeableFilterer{contract: contract}, nil
}

// bindAddressUpgradeable binds a generic wrapper to an already deployed contract.
func bindAddressUpgradeable(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := AddressUpgradeableMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AddressUpgradeable *AddressUpgradeableRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AddressUpgradeable.Contract.AddressUpgradeableCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AddressUpgradeable *AddressUpgradeableRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AddressUpgradeable.Contract.AddressUpgradeableTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AddressUpgradeable *AddressUpgradeableRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AddressUpgradeable.Contract.AddressUpgradeableTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AddressUpgradeable *AddressUpgradeableCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AddressUpgradeable.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AddressUpgradeable *AddressUpgradeableTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AddressUpgradeable.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AddressUpgradeable *AddressUpgradeableTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AddressUpgradeable.Contract.contract.Transact(opts, method, params...)
}

// ContextUpgradeableMetaData contains all meta data concerning the ContextUpgradeable contract.
var ContextUpgradeableMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"}]",
}

// ContextUpgradeableABI is the input ABI used to generate the binding from.
// Deprecated: Use ContextUpgradeableMetaData.ABI instead.
var ContextUpgradeableABI = ContextUpgradeableMetaData.ABI

// ContextUpgradeableBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const ContextUpgradeableBinRuntime = ``

// ContextUpgradeable is an auto generated Go binding around a Klaytn contract.
type ContextUpgradeable struct {
	ContextUpgradeableCaller     // Read-only binding to the contract
	ContextUpgradeableTransactor // Write-only binding to the contract
	ContextUpgradeableFilterer   // Log filterer for contract events
}

// ContextUpgradeableCaller is an auto generated read-only Go binding around a Klaytn contract.
type ContextUpgradeableCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContextUpgradeableTransactor is an auto generated write-only Go binding around a Klaytn contract.
type ContextUpgradeableTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContextUpgradeableFilterer is an auto generated log filtering Go binding around a Klaytn contract events.
type ContextUpgradeableFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContextUpgradeableSession is an auto generated Go binding around a Klaytn contract,
// with pre-set call and transact options.
type ContextUpgradeableSession struct {
	Contract     *ContextUpgradeable // Generic contract binding to set the session for
	CallOpts     bind.CallOpts       // Call options to use throughout this session
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// ContextUpgradeableCallerSession is an auto generated read-only Go binding around a Klaytn contract,
// with pre-set call options.
type ContextUpgradeableCallerSession struct {
	Contract *ContextUpgradeableCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts             // Call options to use throughout this session
}

// ContextUpgradeableTransactorSession is an auto generated write-only Go binding around a Klaytn contract,
// with pre-set transact options.
type ContextUpgradeableTransactorSession struct {
	Contract     *ContextUpgradeableTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts             // Transaction auth options to use throughout this session
}

// ContextUpgradeableRaw is an auto generated low-level Go binding around a Klaytn contract.
type ContextUpgradeableRaw struct {
	Contract *ContextUpgradeable // Generic contract binding to access the raw methods on
}

// ContextUpgradeableCallerRaw is an auto generated low-level read-only Go binding around a Klaytn contract.
type ContextUpgradeableCallerRaw struct {
	Contract *ContextUpgradeableCaller // Generic read-only contract binding to access the raw methods on
}

// ContextUpgradeableTransactorRaw is an auto generated low-level write-only Go binding around a Klaytn contract.
type ContextUpgradeableTransactorRaw struct {
	Contract *ContextUpgradeableTransactor // Generic write-only contract binding to access the raw methods on
}

// NewContextUpgradeable creates a new instance of ContextUpgradeable, bound to a specific deployed contract.
func NewContextUpgradeable(address common.Address, backend bind.ContractBackend) (*ContextUpgradeable, error) {
	contract, err := bindContextUpgradeable(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ContextUpgradeable{ContextUpgradeableCaller: ContextUpgradeableCaller{contract: contract}, ContextUpgradeableTransactor: ContextUpgradeableTransactor{contract: contract}, ContextUpgradeableFilterer: ContextUpgradeableFilterer{contract: contract}}, nil
}

// NewContextUpgradeableCaller creates a new read-only instance of ContextUpgradeable, bound to a specific deployed contract.
func NewContextUpgradeableCaller(address common.Address, caller bind.ContractCaller) (*ContextUpgradeableCaller, error) {
	contract, err := bindContextUpgradeable(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContextUpgradeableCaller{contract: contract}, nil
}

// NewContextUpgradeableTransactor creates a new write-only instance of ContextUpgradeable, bound to a specific deployed contract.
func NewContextUpgradeableTransactor(address common.Address, transactor bind.ContractTransactor) (*ContextUpgradeableTransactor, error) {
	contract, err := bindContextUpgradeable(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContextUpgradeableTransactor{contract: contract}, nil
}

// NewContextUpgradeableFilterer creates a new log filterer instance of ContextUpgradeable, bound to a specific deployed contract.
func NewContextUpgradeableFilterer(address common.Address, filterer bind.ContractFilterer) (*ContextUpgradeableFilterer, error) {
	contract, err := bindContextUpgradeable(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContextUpgradeableFilterer{contract: contract}, nil
}

// bindContextUpgradeable binds a generic wrapper to an already deployed contract.
func bindContextUpgradeable(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ContextUpgradeableMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContextUpgradeable *ContextUpgradeableRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContextUpgradeable.Contract.ContextUpgradeableCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContextUpgradeable *ContextUpgradeableRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContextUpgradeable.Contract.ContextUpgradeableTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContextUpgradeable *ContextUpgradeableRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContextUpgradeable.Contract.ContextUpgradeableTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContextUpgradeable *ContextUpgradeableCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContextUpgradeable.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContextUpgradeable *ContextUpgradeableTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContextUpgradeable.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContextUpgradeable *ContextUpgradeableTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContextUpgradeable.Contract.contract.Transact(opts, method, params...)
}

// ContextUpgradeableInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the ContextUpgradeable contract.
type ContextUpgradeableInitializedIterator struct {
	Event *ContextUpgradeableInitialized // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log      // Log channel receiving the found contract events
	sub  klaytn.Subscription // Subscription for errors, completion and termination
	done bool                // Whether the subscription completed delivering logs
	fail error               // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContextUpgradeableInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContextUpgradeableInitialized)
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
		it.Event = new(ContextUpgradeableInitialized)
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
func (it *ContextUpgradeableInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContextUpgradeableInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContextUpgradeableInitialized represents a Initialized event raised by the ContextUpgradeable contract.
type ContextUpgradeableInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_ContextUpgradeable *ContextUpgradeableFilterer) FilterInitialized(opts *bind.FilterOpts) (*ContextUpgradeableInitializedIterator, error) {
	logs, sub, err := _ContextUpgradeable.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &ContextUpgradeableInitializedIterator{contract: _ContextUpgradeable.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_ContextUpgradeable *ContextUpgradeableFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *ContextUpgradeableInitialized) (event.Subscription, error) {
	logs, sub, err := _ContextUpgradeable.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContextUpgradeableInitialized)
				if err := _ContextUpgradeable.contract.UnpackLog(event, "Initialized", log); err != nil {
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

// ParseInitialized is a log parse operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_ContextUpgradeable *ContextUpgradeableFilterer) ParseInitialized(log types.Log) (*ContextUpgradeableInitialized, error) {
	event := new(ContextUpgradeableInitialized)
	if err := _ContextUpgradeable.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ERC1967ProxyMetaData contains all meta data concerning the ERC1967Proxy contract.
var ERC1967ProxyMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_logic\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"}],\"stateMutability\":\"payable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"previousAdmin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"AdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"beacon\",\"type\":\"address\"}],\"name\":\"BeaconUpgraded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"},{\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Bin: "0x60806040526040516104e13803806104e1833981016040819052610022916102de565b61002e82826000610035565b50506103fb565b61003e83610061565b60008251118061004b5750805b1561005c5761005a83836100a1565b505b505050565b61006a816100cd565b6040516001600160a01b038216907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b90600090a250565b60606100c683836040518060600160405280602781526020016104ba60279139610180565b9392505050565b6001600160a01b0381163b61013f5760405162461bcd60e51b815260206004820152602d60248201527f455243313936373a206e657720696d706c656d656e746174696f6e206973206e60448201526c1bdd08184818dbdb9d1c9858dd609a1b60648201526084015b60405180910390fd5b7f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc80546001600160a01b0319166001600160a01b0392909216919091179055565b6060600080856001600160a01b03168560405161019d91906103ac565b600060405180830381855af49150503d80600081146101d8576040519150601f19603f3d011682016040523d82523d6000602084013e6101dd565b606091505b5090925090506101ef868383876101f9565b9695505050505050565b60608315610268578251600003610261576001600160a01b0385163b6102615760405162461bcd60e51b815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e74726163740000006044820152606401610136565b5081610272565b610272838361027a565b949350505050565b81511561028a5781518083602001fd5b8060405162461bcd60e51b815260040161013691906103c8565b634e487b7160e01b600052604160045260246000fd5b60005b838110156102d55781810151838201526020016102bd565b50506000910152565b600080604083850312156102f157600080fd5b82516001600160a01b038116811461030857600080fd5b60208401519092506001600160401b038082111561032557600080fd5b818501915085601f83011261033957600080fd5b81518181111561034b5761034b6102a4565b604051601f8201601f19908116603f01168101908382118183101715610373576103736102a4565b8160405282815288602084870101111561038c57600080fd5b61039d8360208301602088016102ba565b80955050505050509250929050565b600082516103be8184602087016102ba565b9190910192915050565b60208152600082518060208401526103e78160408501602087016102ba565b601f01601f19169190910160400192915050565b60b1806104096000396000f3fe608060405236601057600e6013565b005b600e5b601f601b6021565b6058565b565b600060537f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc546001600160a01b031690565b905090565b3660008037600080366000845af43d6000803e8080156076573d6000f35b3d6000fdfea264697066735822122032c702aadb8f2fab233834f027cfcd56871fb029b5b175f59fdf7b82bb7dd0dc64736f6c63430008130033416464726573733a206c6f772d6c6576656c2064656c65676174652063616c6c206661696c6564",
}

// ERC1967ProxyABI is the input ABI used to generate the binding from.
// Deprecated: Use ERC1967ProxyMetaData.ABI instead.
var ERC1967ProxyABI = ERC1967ProxyMetaData.ABI

// ERC1967ProxyBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const ERC1967ProxyBinRuntime = `608060405236601057600e6013565b005b600e5b601f601b6021565b6058565b565b600060537f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc546001600160a01b031690565b905090565b3660008037600080366000845af43d6000803e8080156076573d6000f35b3d6000fdfea264697066735822122032c702aadb8f2fab233834f027cfcd56871fb029b5b175f59fdf7b82bb7dd0dc64736f6c63430008130033`

// ERC1967ProxyBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ERC1967ProxyMetaData.Bin instead.
var ERC1967ProxyBin = ERC1967ProxyMetaData.Bin

// DeployERC1967Proxy deploys a new Klaytn contract, binding an instance of ERC1967Proxy to it.
func DeployERC1967Proxy(auth *bind.TransactOpts, backend bind.ContractBackend, _logic common.Address, _data []byte) (common.Address, *types.Transaction, *ERC1967Proxy, error) {
	parsed, err := ERC1967ProxyMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ERC1967ProxyBin), backend, _logic, _data)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ERC1967Proxy{ERC1967ProxyCaller: ERC1967ProxyCaller{contract: contract}, ERC1967ProxyTransactor: ERC1967ProxyTransactor{contract: contract}, ERC1967ProxyFilterer: ERC1967ProxyFilterer{contract: contract}}, nil
}

// ERC1967Proxy is an auto generated Go binding around a Klaytn contract.
type ERC1967Proxy struct {
	ERC1967ProxyCaller     // Read-only binding to the contract
	ERC1967ProxyTransactor // Write-only binding to the contract
	ERC1967ProxyFilterer   // Log filterer for contract events
}

// ERC1967ProxyCaller is an auto generated read-only Go binding around a Klaytn contract.
type ERC1967ProxyCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC1967ProxyTransactor is an auto generated write-only Go binding around a Klaytn contract.
type ERC1967ProxyTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC1967ProxyFilterer is an auto generated log filtering Go binding around a Klaytn contract events.
type ERC1967ProxyFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC1967ProxySession is an auto generated Go binding around a Klaytn contract,
// with pre-set call and transact options.
type ERC1967ProxySession struct {
	Contract     *ERC1967Proxy     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ERC1967ProxyCallerSession is an auto generated read-only Go binding around a Klaytn contract,
// with pre-set call options.
type ERC1967ProxyCallerSession struct {
	Contract *ERC1967ProxyCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// ERC1967ProxyTransactorSession is an auto generated write-only Go binding around a Klaytn contract,
// with pre-set transact options.
type ERC1967ProxyTransactorSession struct {
	Contract     *ERC1967ProxyTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// ERC1967ProxyRaw is an auto generated low-level Go binding around a Klaytn contract.
type ERC1967ProxyRaw struct {
	Contract *ERC1967Proxy // Generic contract binding to access the raw methods on
}

// ERC1967ProxyCallerRaw is an auto generated low-level read-only Go binding around a Klaytn contract.
type ERC1967ProxyCallerRaw struct {
	Contract *ERC1967ProxyCaller // Generic read-only contract binding to access the raw methods on
}

// ERC1967ProxyTransactorRaw is an auto generated low-level write-only Go binding around a Klaytn contract.
type ERC1967ProxyTransactorRaw struct {
	Contract *ERC1967ProxyTransactor // Generic write-only contract binding to access the raw methods on
}

// NewERC1967Proxy creates a new instance of ERC1967Proxy, bound to a specific deployed contract.
func NewERC1967Proxy(address common.Address, backend bind.ContractBackend) (*ERC1967Proxy, error) {
	contract, err := bindERC1967Proxy(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ERC1967Proxy{ERC1967ProxyCaller: ERC1967ProxyCaller{contract: contract}, ERC1967ProxyTransactor: ERC1967ProxyTransactor{contract: contract}, ERC1967ProxyFilterer: ERC1967ProxyFilterer{contract: contract}}, nil
}

// NewERC1967ProxyCaller creates a new read-only instance of ERC1967Proxy, bound to a specific deployed contract.
func NewERC1967ProxyCaller(address common.Address, caller bind.ContractCaller) (*ERC1967ProxyCaller, error) {
	contract, err := bindERC1967Proxy(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ERC1967ProxyCaller{contract: contract}, nil
}

// NewERC1967ProxyTransactor creates a new write-only instance of ERC1967Proxy, bound to a specific deployed contract.
func NewERC1967ProxyTransactor(address common.Address, transactor bind.ContractTransactor) (*ERC1967ProxyTransactor, error) {
	contract, err := bindERC1967Proxy(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ERC1967ProxyTransactor{contract: contract}, nil
}

// NewERC1967ProxyFilterer creates a new log filterer instance of ERC1967Proxy, bound to a specific deployed contract.
func NewERC1967ProxyFilterer(address common.Address, filterer bind.ContractFilterer) (*ERC1967ProxyFilterer, error) {
	contract, err := bindERC1967Proxy(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ERC1967ProxyFilterer{contract: contract}, nil
}

// bindERC1967Proxy binds a generic wrapper to an already deployed contract.
func bindERC1967Proxy(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ERC1967ProxyMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ERC1967Proxy *ERC1967ProxyRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ERC1967Proxy.Contract.ERC1967ProxyCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ERC1967Proxy *ERC1967ProxyRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC1967Proxy.Contract.ERC1967ProxyTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ERC1967Proxy *ERC1967ProxyRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ERC1967Proxy.Contract.ERC1967ProxyTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ERC1967Proxy *ERC1967ProxyCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ERC1967Proxy.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ERC1967Proxy *ERC1967ProxyTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC1967Proxy.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ERC1967Proxy *ERC1967ProxyTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ERC1967Proxy.Contract.contract.Transact(opts, method, params...)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_ERC1967Proxy *ERC1967ProxyTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _ERC1967Proxy.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_ERC1967Proxy *ERC1967ProxySession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _ERC1967Proxy.Contract.Fallback(&_ERC1967Proxy.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_ERC1967Proxy *ERC1967ProxyTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _ERC1967Proxy.Contract.Fallback(&_ERC1967Proxy.TransactOpts, calldata)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_ERC1967Proxy *ERC1967ProxyTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC1967Proxy.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_ERC1967Proxy *ERC1967ProxySession) Receive() (*types.Transaction, error) {
	return _ERC1967Proxy.Contract.Receive(&_ERC1967Proxy.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_ERC1967Proxy *ERC1967ProxyTransactorSession) Receive() (*types.Transaction, error) {
	return _ERC1967Proxy.Contract.Receive(&_ERC1967Proxy.TransactOpts)
}

// ERC1967ProxyAdminChangedIterator is returned from FilterAdminChanged and is used to iterate over the raw logs and unpacked data for AdminChanged events raised by the ERC1967Proxy contract.
type ERC1967ProxyAdminChangedIterator struct {
	Event *ERC1967ProxyAdminChanged // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log      // Log channel receiving the found contract events
	sub  klaytn.Subscription // Subscription for errors, completion and termination
	done bool                // Whether the subscription completed delivering logs
	fail error               // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ERC1967ProxyAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC1967ProxyAdminChanged)
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
		it.Event = new(ERC1967ProxyAdminChanged)
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
func (it *ERC1967ProxyAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC1967ProxyAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC1967ProxyAdminChanged represents a AdminChanged event raised by the ERC1967Proxy contract.
type ERC1967ProxyAdminChanged struct {
	PreviousAdmin common.Address
	NewAdmin      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterAdminChanged is a free log retrieval operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_ERC1967Proxy *ERC1967ProxyFilterer) FilterAdminChanged(opts *bind.FilterOpts) (*ERC1967ProxyAdminChangedIterator, error) {
	logs, sub, err := _ERC1967Proxy.contract.FilterLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return &ERC1967ProxyAdminChangedIterator{contract: _ERC1967Proxy.contract, event: "AdminChanged", logs: logs, sub: sub}, nil
}

// WatchAdminChanged is a free log subscription operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_ERC1967Proxy *ERC1967ProxyFilterer) WatchAdminChanged(opts *bind.WatchOpts, sink chan<- *ERC1967ProxyAdminChanged) (event.Subscription, error) {
	logs, sub, err := _ERC1967Proxy.contract.WatchLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC1967ProxyAdminChanged)
				if err := _ERC1967Proxy.contract.UnpackLog(event, "AdminChanged", log); err != nil {
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

// ParseAdminChanged is a log parse operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_ERC1967Proxy *ERC1967ProxyFilterer) ParseAdminChanged(log types.Log) (*ERC1967ProxyAdminChanged, error) {
	event := new(ERC1967ProxyAdminChanged)
	if err := _ERC1967Proxy.contract.UnpackLog(event, "AdminChanged", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ERC1967ProxyBeaconUpgradedIterator is returned from FilterBeaconUpgraded and is used to iterate over the raw logs and unpacked data for BeaconUpgraded events raised by the ERC1967Proxy contract.
type ERC1967ProxyBeaconUpgradedIterator struct {
	Event *ERC1967ProxyBeaconUpgraded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log      // Log channel receiving the found contract events
	sub  klaytn.Subscription // Subscription for errors, completion and termination
	done bool                // Whether the subscription completed delivering logs
	fail error               // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ERC1967ProxyBeaconUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC1967ProxyBeaconUpgraded)
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
		it.Event = new(ERC1967ProxyBeaconUpgraded)
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
func (it *ERC1967ProxyBeaconUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC1967ProxyBeaconUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC1967ProxyBeaconUpgraded represents a BeaconUpgraded event raised by the ERC1967Proxy contract.
type ERC1967ProxyBeaconUpgraded struct {
	Beacon common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterBeaconUpgraded is a free log retrieval operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_ERC1967Proxy *ERC1967ProxyFilterer) FilterBeaconUpgraded(opts *bind.FilterOpts, beacon []common.Address) (*ERC1967ProxyBeaconUpgradedIterator, error) {
	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _ERC1967Proxy.contract.FilterLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return &ERC1967ProxyBeaconUpgradedIterator{contract: _ERC1967Proxy.contract, event: "BeaconUpgraded", logs: logs, sub: sub}, nil
}

// WatchBeaconUpgraded is a free log subscription operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_ERC1967Proxy *ERC1967ProxyFilterer) WatchBeaconUpgraded(opts *bind.WatchOpts, sink chan<- *ERC1967ProxyBeaconUpgraded, beacon []common.Address) (event.Subscription, error) {
	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _ERC1967Proxy.contract.WatchLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC1967ProxyBeaconUpgraded)
				if err := _ERC1967Proxy.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
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

// ParseBeaconUpgraded is a log parse operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_ERC1967Proxy *ERC1967ProxyFilterer) ParseBeaconUpgraded(log types.Log) (*ERC1967ProxyBeaconUpgraded, error) {
	event := new(ERC1967ProxyBeaconUpgraded)
	if err := _ERC1967Proxy.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ERC1967ProxyUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the ERC1967Proxy contract.
type ERC1967ProxyUpgradedIterator struct {
	Event *ERC1967ProxyUpgraded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log      // Log channel receiving the found contract events
	sub  klaytn.Subscription // Subscription for errors, completion and termination
	done bool                // Whether the subscription completed delivering logs
	fail error               // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ERC1967ProxyUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC1967ProxyUpgraded)
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
		it.Event = new(ERC1967ProxyUpgraded)
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
func (it *ERC1967ProxyUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC1967ProxyUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC1967ProxyUpgraded represents a Upgraded event raised by the ERC1967Proxy contract.
type ERC1967ProxyUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_ERC1967Proxy *ERC1967ProxyFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*ERC1967ProxyUpgradedIterator, error) {
	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _ERC1967Proxy.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &ERC1967ProxyUpgradedIterator{contract: _ERC1967Proxy.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_ERC1967Proxy *ERC1967ProxyFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *ERC1967ProxyUpgraded, implementation []common.Address) (event.Subscription, error) {
	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _ERC1967Proxy.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC1967ProxyUpgraded)
				if err := _ERC1967Proxy.contract.UnpackLog(event, "Upgraded", log); err != nil {
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

// ParseUpgraded is a log parse operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_ERC1967Proxy *ERC1967ProxyFilterer) ParseUpgraded(log types.Log) (*ERC1967ProxyUpgraded, error) {
	event := new(ERC1967ProxyUpgraded)
	if err := _ERC1967Proxy.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ERC1967UpgradeMetaData contains all meta data concerning the ERC1967Upgrade contract.
var ERC1967UpgradeMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"previousAdmin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"AdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"beacon\",\"type\":\"address\"}],\"name\":\"BeaconUpgraded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"}]",
}

// ERC1967UpgradeABI is the input ABI used to generate the binding from.
// Deprecated: Use ERC1967UpgradeMetaData.ABI instead.
var ERC1967UpgradeABI = ERC1967UpgradeMetaData.ABI

// ERC1967UpgradeBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const ERC1967UpgradeBinRuntime = ``

// ERC1967Upgrade is an auto generated Go binding around a Klaytn contract.
type ERC1967Upgrade struct {
	ERC1967UpgradeCaller     // Read-only binding to the contract
	ERC1967UpgradeTransactor // Write-only binding to the contract
	ERC1967UpgradeFilterer   // Log filterer for contract events
}

// ERC1967UpgradeCaller is an auto generated read-only Go binding around a Klaytn contract.
type ERC1967UpgradeCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC1967UpgradeTransactor is an auto generated write-only Go binding around a Klaytn contract.
type ERC1967UpgradeTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC1967UpgradeFilterer is an auto generated log filtering Go binding around a Klaytn contract events.
type ERC1967UpgradeFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC1967UpgradeSession is an auto generated Go binding around a Klaytn contract,
// with pre-set call and transact options.
type ERC1967UpgradeSession struct {
	Contract     *ERC1967Upgrade   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ERC1967UpgradeCallerSession is an auto generated read-only Go binding around a Klaytn contract,
// with pre-set call options.
type ERC1967UpgradeCallerSession struct {
	Contract *ERC1967UpgradeCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// ERC1967UpgradeTransactorSession is an auto generated write-only Go binding around a Klaytn contract,
// with pre-set transact options.
type ERC1967UpgradeTransactorSession struct {
	Contract     *ERC1967UpgradeTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// ERC1967UpgradeRaw is an auto generated low-level Go binding around a Klaytn contract.
type ERC1967UpgradeRaw struct {
	Contract *ERC1967Upgrade // Generic contract binding to access the raw methods on
}

// ERC1967UpgradeCallerRaw is an auto generated low-level read-only Go binding around a Klaytn contract.
type ERC1967UpgradeCallerRaw struct {
	Contract *ERC1967UpgradeCaller // Generic read-only contract binding to access the raw methods on
}

// ERC1967UpgradeTransactorRaw is an auto generated low-level write-only Go binding around a Klaytn contract.
type ERC1967UpgradeTransactorRaw struct {
	Contract *ERC1967UpgradeTransactor // Generic write-only contract binding to access the raw methods on
}

// NewERC1967Upgrade creates a new instance of ERC1967Upgrade, bound to a specific deployed contract.
func NewERC1967Upgrade(address common.Address, backend bind.ContractBackend) (*ERC1967Upgrade, error) {
	contract, err := bindERC1967Upgrade(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ERC1967Upgrade{ERC1967UpgradeCaller: ERC1967UpgradeCaller{contract: contract}, ERC1967UpgradeTransactor: ERC1967UpgradeTransactor{contract: contract}, ERC1967UpgradeFilterer: ERC1967UpgradeFilterer{contract: contract}}, nil
}

// NewERC1967UpgradeCaller creates a new read-only instance of ERC1967Upgrade, bound to a specific deployed contract.
func NewERC1967UpgradeCaller(address common.Address, caller bind.ContractCaller) (*ERC1967UpgradeCaller, error) {
	contract, err := bindERC1967Upgrade(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ERC1967UpgradeCaller{contract: contract}, nil
}

// NewERC1967UpgradeTransactor creates a new write-only instance of ERC1967Upgrade, bound to a specific deployed contract.
func NewERC1967UpgradeTransactor(address common.Address, transactor bind.ContractTransactor) (*ERC1967UpgradeTransactor, error) {
	contract, err := bindERC1967Upgrade(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ERC1967UpgradeTransactor{contract: contract}, nil
}

// NewERC1967UpgradeFilterer creates a new log filterer instance of ERC1967Upgrade, bound to a specific deployed contract.
func NewERC1967UpgradeFilterer(address common.Address, filterer bind.ContractFilterer) (*ERC1967UpgradeFilterer, error) {
	contract, err := bindERC1967Upgrade(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ERC1967UpgradeFilterer{contract: contract}, nil
}

// bindERC1967Upgrade binds a generic wrapper to an already deployed contract.
func bindERC1967Upgrade(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ERC1967UpgradeMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ERC1967Upgrade *ERC1967UpgradeRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ERC1967Upgrade.Contract.ERC1967UpgradeCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ERC1967Upgrade *ERC1967UpgradeRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC1967Upgrade.Contract.ERC1967UpgradeTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ERC1967Upgrade *ERC1967UpgradeRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ERC1967Upgrade.Contract.ERC1967UpgradeTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ERC1967Upgrade *ERC1967UpgradeCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ERC1967Upgrade.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ERC1967Upgrade *ERC1967UpgradeTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC1967Upgrade.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ERC1967Upgrade *ERC1967UpgradeTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ERC1967Upgrade.Contract.contract.Transact(opts, method, params...)
}

// ERC1967UpgradeAdminChangedIterator is returned from FilterAdminChanged and is used to iterate over the raw logs and unpacked data for AdminChanged events raised by the ERC1967Upgrade contract.
type ERC1967UpgradeAdminChangedIterator struct {
	Event *ERC1967UpgradeAdminChanged // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log      // Log channel receiving the found contract events
	sub  klaytn.Subscription // Subscription for errors, completion and termination
	done bool                // Whether the subscription completed delivering logs
	fail error               // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ERC1967UpgradeAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC1967UpgradeAdminChanged)
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
		it.Event = new(ERC1967UpgradeAdminChanged)
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
func (it *ERC1967UpgradeAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC1967UpgradeAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC1967UpgradeAdminChanged represents a AdminChanged event raised by the ERC1967Upgrade contract.
type ERC1967UpgradeAdminChanged struct {
	PreviousAdmin common.Address
	NewAdmin      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterAdminChanged is a free log retrieval operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_ERC1967Upgrade *ERC1967UpgradeFilterer) FilterAdminChanged(opts *bind.FilterOpts) (*ERC1967UpgradeAdminChangedIterator, error) {
	logs, sub, err := _ERC1967Upgrade.contract.FilterLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return &ERC1967UpgradeAdminChangedIterator{contract: _ERC1967Upgrade.contract, event: "AdminChanged", logs: logs, sub: sub}, nil
}

// WatchAdminChanged is a free log subscription operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_ERC1967Upgrade *ERC1967UpgradeFilterer) WatchAdminChanged(opts *bind.WatchOpts, sink chan<- *ERC1967UpgradeAdminChanged) (event.Subscription, error) {
	logs, sub, err := _ERC1967Upgrade.contract.WatchLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC1967UpgradeAdminChanged)
				if err := _ERC1967Upgrade.contract.UnpackLog(event, "AdminChanged", log); err != nil {
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

// ParseAdminChanged is a log parse operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_ERC1967Upgrade *ERC1967UpgradeFilterer) ParseAdminChanged(log types.Log) (*ERC1967UpgradeAdminChanged, error) {
	event := new(ERC1967UpgradeAdminChanged)
	if err := _ERC1967Upgrade.contract.UnpackLog(event, "AdminChanged", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ERC1967UpgradeBeaconUpgradedIterator is returned from FilterBeaconUpgraded and is used to iterate over the raw logs and unpacked data for BeaconUpgraded events raised by the ERC1967Upgrade contract.
type ERC1967UpgradeBeaconUpgradedIterator struct {
	Event *ERC1967UpgradeBeaconUpgraded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log      // Log channel receiving the found contract events
	sub  klaytn.Subscription // Subscription for errors, completion and termination
	done bool                // Whether the subscription completed delivering logs
	fail error               // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ERC1967UpgradeBeaconUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC1967UpgradeBeaconUpgraded)
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
		it.Event = new(ERC1967UpgradeBeaconUpgraded)
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
func (it *ERC1967UpgradeBeaconUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC1967UpgradeBeaconUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC1967UpgradeBeaconUpgraded represents a BeaconUpgraded event raised by the ERC1967Upgrade contract.
type ERC1967UpgradeBeaconUpgraded struct {
	Beacon common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterBeaconUpgraded is a free log retrieval operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_ERC1967Upgrade *ERC1967UpgradeFilterer) FilterBeaconUpgraded(opts *bind.FilterOpts, beacon []common.Address) (*ERC1967UpgradeBeaconUpgradedIterator, error) {
	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _ERC1967Upgrade.contract.FilterLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return &ERC1967UpgradeBeaconUpgradedIterator{contract: _ERC1967Upgrade.contract, event: "BeaconUpgraded", logs: logs, sub: sub}, nil
}

// WatchBeaconUpgraded is a free log subscription operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_ERC1967Upgrade *ERC1967UpgradeFilterer) WatchBeaconUpgraded(opts *bind.WatchOpts, sink chan<- *ERC1967UpgradeBeaconUpgraded, beacon []common.Address) (event.Subscription, error) {
	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _ERC1967Upgrade.contract.WatchLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC1967UpgradeBeaconUpgraded)
				if err := _ERC1967Upgrade.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
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

// ParseBeaconUpgraded is a log parse operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_ERC1967Upgrade *ERC1967UpgradeFilterer) ParseBeaconUpgraded(log types.Log) (*ERC1967UpgradeBeaconUpgraded, error) {
	event := new(ERC1967UpgradeBeaconUpgraded)
	if err := _ERC1967Upgrade.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ERC1967UpgradeUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the ERC1967Upgrade contract.
type ERC1967UpgradeUpgradedIterator struct {
	Event *ERC1967UpgradeUpgraded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log      // Log channel receiving the found contract events
	sub  klaytn.Subscription // Subscription for errors, completion and termination
	done bool                // Whether the subscription completed delivering logs
	fail error               // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ERC1967UpgradeUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC1967UpgradeUpgraded)
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
		it.Event = new(ERC1967UpgradeUpgraded)
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
func (it *ERC1967UpgradeUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC1967UpgradeUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC1967UpgradeUpgraded represents a Upgraded event raised by the ERC1967Upgrade contract.
type ERC1967UpgradeUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_ERC1967Upgrade *ERC1967UpgradeFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*ERC1967UpgradeUpgradedIterator, error) {
	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _ERC1967Upgrade.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &ERC1967UpgradeUpgradedIterator{contract: _ERC1967Upgrade.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_ERC1967Upgrade *ERC1967UpgradeFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *ERC1967UpgradeUpgraded, implementation []common.Address) (event.Subscription, error) {
	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _ERC1967Upgrade.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC1967UpgradeUpgraded)
				if err := _ERC1967Upgrade.contract.UnpackLog(event, "Upgraded", log); err != nil {
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

// ParseUpgraded is a log parse operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_ERC1967Upgrade *ERC1967UpgradeFilterer) ParseUpgraded(log types.Log) (*ERC1967UpgradeUpgraded, error) {
	event := new(ERC1967UpgradeUpgraded)
	if err := _ERC1967Upgrade.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ERC1967UpgradeUpgradeableMetaData contains all meta data concerning the ERC1967UpgradeUpgradeable contract.
var ERC1967UpgradeUpgradeableMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"previousAdmin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"AdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"beacon\",\"type\":\"address\"}],\"name\":\"BeaconUpgraded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"}]",
}

// ERC1967UpgradeUpgradeableABI is the input ABI used to generate the binding from.
// Deprecated: Use ERC1967UpgradeUpgradeableMetaData.ABI instead.
var ERC1967UpgradeUpgradeableABI = ERC1967UpgradeUpgradeableMetaData.ABI

// ERC1967UpgradeUpgradeableBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const ERC1967UpgradeUpgradeableBinRuntime = ``

// ERC1967UpgradeUpgradeable is an auto generated Go binding around a Klaytn contract.
type ERC1967UpgradeUpgradeable struct {
	ERC1967UpgradeUpgradeableCaller     // Read-only binding to the contract
	ERC1967UpgradeUpgradeableTransactor // Write-only binding to the contract
	ERC1967UpgradeUpgradeableFilterer   // Log filterer for contract events
}

// ERC1967UpgradeUpgradeableCaller is an auto generated read-only Go binding around a Klaytn contract.
type ERC1967UpgradeUpgradeableCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC1967UpgradeUpgradeableTransactor is an auto generated write-only Go binding around a Klaytn contract.
type ERC1967UpgradeUpgradeableTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC1967UpgradeUpgradeableFilterer is an auto generated log filtering Go binding around a Klaytn contract events.
type ERC1967UpgradeUpgradeableFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC1967UpgradeUpgradeableSession is an auto generated Go binding around a Klaytn contract,
// with pre-set call and transact options.
type ERC1967UpgradeUpgradeableSession struct {
	Contract     *ERC1967UpgradeUpgradeable // Generic contract binding to set the session for
	CallOpts     bind.CallOpts              // Call options to use throughout this session
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// ERC1967UpgradeUpgradeableCallerSession is an auto generated read-only Go binding around a Klaytn contract,
// with pre-set call options.
type ERC1967UpgradeUpgradeableCallerSession struct {
	Contract *ERC1967UpgradeUpgradeableCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                    // Call options to use throughout this session
}

// ERC1967UpgradeUpgradeableTransactorSession is an auto generated write-only Go binding around a Klaytn contract,
// with pre-set transact options.
type ERC1967UpgradeUpgradeableTransactorSession struct {
	Contract     *ERC1967UpgradeUpgradeableTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                    // Transaction auth options to use throughout this session
}

// ERC1967UpgradeUpgradeableRaw is an auto generated low-level Go binding around a Klaytn contract.
type ERC1967UpgradeUpgradeableRaw struct {
	Contract *ERC1967UpgradeUpgradeable // Generic contract binding to access the raw methods on
}

// ERC1967UpgradeUpgradeableCallerRaw is an auto generated low-level read-only Go binding around a Klaytn contract.
type ERC1967UpgradeUpgradeableCallerRaw struct {
	Contract *ERC1967UpgradeUpgradeableCaller // Generic read-only contract binding to access the raw methods on
}

// ERC1967UpgradeUpgradeableTransactorRaw is an auto generated low-level write-only Go binding around a Klaytn contract.
type ERC1967UpgradeUpgradeableTransactorRaw struct {
	Contract *ERC1967UpgradeUpgradeableTransactor // Generic write-only contract binding to access the raw methods on
}

// NewERC1967UpgradeUpgradeable creates a new instance of ERC1967UpgradeUpgradeable, bound to a specific deployed contract.
func NewERC1967UpgradeUpgradeable(address common.Address, backend bind.ContractBackend) (*ERC1967UpgradeUpgradeable, error) {
	contract, err := bindERC1967UpgradeUpgradeable(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ERC1967UpgradeUpgradeable{ERC1967UpgradeUpgradeableCaller: ERC1967UpgradeUpgradeableCaller{contract: contract}, ERC1967UpgradeUpgradeableTransactor: ERC1967UpgradeUpgradeableTransactor{contract: contract}, ERC1967UpgradeUpgradeableFilterer: ERC1967UpgradeUpgradeableFilterer{contract: contract}}, nil
}

// NewERC1967UpgradeUpgradeableCaller creates a new read-only instance of ERC1967UpgradeUpgradeable, bound to a specific deployed contract.
func NewERC1967UpgradeUpgradeableCaller(address common.Address, caller bind.ContractCaller) (*ERC1967UpgradeUpgradeableCaller, error) {
	contract, err := bindERC1967UpgradeUpgradeable(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ERC1967UpgradeUpgradeableCaller{contract: contract}, nil
}

// NewERC1967UpgradeUpgradeableTransactor creates a new write-only instance of ERC1967UpgradeUpgradeable, bound to a specific deployed contract.
func NewERC1967UpgradeUpgradeableTransactor(address common.Address, transactor bind.ContractTransactor) (*ERC1967UpgradeUpgradeableTransactor, error) {
	contract, err := bindERC1967UpgradeUpgradeable(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ERC1967UpgradeUpgradeableTransactor{contract: contract}, nil
}

// NewERC1967UpgradeUpgradeableFilterer creates a new log filterer instance of ERC1967UpgradeUpgradeable, bound to a specific deployed contract.
func NewERC1967UpgradeUpgradeableFilterer(address common.Address, filterer bind.ContractFilterer) (*ERC1967UpgradeUpgradeableFilterer, error) {
	contract, err := bindERC1967UpgradeUpgradeable(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ERC1967UpgradeUpgradeableFilterer{contract: contract}, nil
}

// bindERC1967UpgradeUpgradeable binds a generic wrapper to an already deployed contract.
func bindERC1967UpgradeUpgradeable(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ERC1967UpgradeUpgradeableMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ERC1967UpgradeUpgradeable *ERC1967UpgradeUpgradeableRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ERC1967UpgradeUpgradeable.Contract.ERC1967UpgradeUpgradeableCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ERC1967UpgradeUpgradeable *ERC1967UpgradeUpgradeableRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC1967UpgradeUpgradeable.Contract.ERC1967UpgradeUpgradeableTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ERC1967UpgradeUpgradeable *ERC1967UpgradeUpgradeableRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ERC1967UpgradeUpgradeable.Contract.ERC1967UpgradeUpgradeableTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ERC1967UpgradeUpgradeable *ERC1967UpgradeUpgradeableCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ERC1967UpgradeUpgradeable.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ERC1967UpgradeUpgradeable *ERC1967UpgradeUpgradeableTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC1967UpgradeUpgradeable.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ERC1967UpgradeUpgradeable *ERC1967UpgradeUpgradeableTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ERC1967UpgradeUpgradeable.Contract.contract.Transact(opts, method, params...)
}

// ERC1967UpgradeUpgradeableAdminChangedIterator is returned from FilterAdminChanged and is used to iterate over the raw logs and unpacked data for AdminChanged events raised by the ERC1967UpgradeUpgradeable contract.
type ERC1967UpgradeUpgradeableAdminChangedIterator struct {
	Event *ERC1967UpgradeUpgradeableAdminChanged // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log      // Log channel receiving the found contract events
	sub  klaytn.Subscription // Subscription for errors, completion and termination
	done bool                // Whether the subscription completed delivering logs
	fail error               // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ERC1967UpgradeUpgradeableAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC1967UpgradeUpgradeableAdminChanged)
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
		it.Event = new(ERC1967UpgradeUpgradeableAdminChanged)
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
func (it *ERC1967UpgradeUpgradeableAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC1967UpgradeUpgradeableAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC1967UpgradeUpgradeableAdminChanged represents a AdminChanged event raised by the ERC1967UpgradeUpgradeable contract.
type ERC1967UpgradeUpgradeableAdminChanged struct {
	PreviousAdmin common.Address
	NewAdmin      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterAdminChanged is a free log retrieval operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_ERC1967UpgradeUpgradeable *ERC1967UpgradeUpgradeableFilterer) FilterAdminChanged(opts *bind.FilterOpts) (*ERC1967UpgradeUpgradeableAdminChangedIterator, error) {
	logs, sub, err := _ERC1967UpgradeUpgradeable.contract.FilterLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return &ERC1967UpgradeUpgradeableAdminChangedIterator{contract: _ERC1967UpgradeUpgradeable.contract, event: "AdminChanged", logs: logs, sub: sub}, nil
}

// WatchAdminChanged is a free log subscription operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_ERC1967UpgradeUpgradeable *ERC1967UpgradeUpgradeableFilterer) WatchAdminChanged(opts *bind.WatchOpts, sink chan<- *ERC1967UpgradeUpgradeableAdminChanged) (event.Subscription, error) {
	logs, sub, err := _ERC1967UpgradeUpgradeable.contract.WatchLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC1967UpgradeUpgradeableAdminChanged)
				if err := _ERC1967UpgradeUpgradeable.contract.UnpackLog(event, "AdminChanged", log); err != nil {
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

// ParseAdminChanged is a log parse operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_ERC1967UpgradeUpgradeable *ERC1967UpgradeUpgradeableFilterer) ParseAdminChanged(log types.Log) (*ERC1967UpgradeUpgradeableAdminChanged, error) {
	event := new(ERC1967UpgradeUpgradeableAdminChanged)
	if err := _ERC1967UpgradeUpgradeable.contract.UnpackLog(event, "AdminChanged", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ERC1967UpgradeUpgradeableBeaconUpgradedIterator is returned from FilterBeaconUpgraded and is used to iterate over the raw logs and unpacked data for BeaconUpgraded events raised by the ERC1967UpgradeUpgradeable contract.
type ERC1967UpgradeUpgradeableBeaconUpgradedIterator struct {
	Event *ERC1967UpgradeUpgradeableBeaconUpgraded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log      // Log channel receiving the found contract events
	sub  klaytn.Subscription // Subscription for errors, completion and termination
	done bool                // Whether the subscription completed delivering logs
	fail error               // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ERC1967UpgradeUpgradeableBeaconUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC1967UpgradeUpgradeableBeaconUpgraded)
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
		it.Event = new(ERC1967UpgradeUpgradeableBeaconUpgraded)
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
func (it *ERC1967UpgradeUpgradeableBeaconUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC1967UpgradeUpgradeableBeaconUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC1967UpgradeUpgradeableBeaconUpgraded represents a BeaconUpgraded event raised by the ERC1967UpgradeUpgradeable contract.
type ERC1967UpgradeUpgradeableBeaconUpgraded struct {
	Beacon common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterBeaconUpgraded is a free log retrieval operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_ERC1967UpgradeUpgradeable *ERC1967UpgradeUpgradeableFilterer) FilterBeaconUpgraded(opts *bind.FilterOpts, beacon []common.Address) (*ERC1967UpgradeUpgradeableBeaconUpgradedIterator, error) {
	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _ERC1967UpgradeUpgradeable.contract.FilterLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return &ERC1967UpgradeUpgradeableBeaconUpgradedIterator{contract: _ERC1967UpgradeUpgradeable.contract, event: "BeaconUpgraded", logs: logs, sub: sub}, nil
}

// WatchBeaconUpgraded is a free log subscription operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_ERC1967UpgradeUpgradeable *ERC1967UpgradeUpgradeableFilterer) WatchBeaconUpgraded(opts *bind.WatchOpts, sink chan<- *ERC1967UpgradeUpgradeableBeaconUpgraded, beacon []common.Address) (event.Subscription, error) {
	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _ERC1967UpgradeUpgradeable.contract.WatchLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC1967UpgradeUpgradeableBeaconUpgraded)
				if err := _ERC1967UpgradeUpgradeable.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
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

// ParseBeaconUpgraded is a log parse operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_ERC1967UpgradeUpgradeable *ERC1967UpgradeUpgradeableFilterer) ParseBeaconUpgraded(log types.Log) (*ERC1967UpgradeUpgradeableBeaconUpgraded, error) {
	event := new(ERC1967UpgradeUpgradeableBeaconUpgraded)
	if err := _ERC1967UpgradeUpgradeable.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ERC1967UpgradeUpgradeableInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the ERC1967UpgradeUpgradeable contract.
type ERC1967UpgradeUpgradeableInitializedIterator struct {
	Event *ERC1967UpgradeUpgradeableInitialized // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log      // Log channel receiving the found contract events
	sub  klaytn.Subscription // Subscription for errors, completion and termination
	done bool                // Whether the subscription completed delivering logs
	fail error               // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ERC1967UpgradeUpgradeableInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC1967UpgradeUpgradeableInitialized)
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
		it.Event = new(ERC1967UpgradeUpgradeableInitialized)
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
func (it *ERC1967UpgradeUpgradeableInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC1967UpgradeUpgradeableInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC1967UpgradeUpgradeableInitialized represents a Initialized event raised by the ERC1967UpgradeUpgradeable contract.
type ERC1967UpgradeUpgradeableInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_ERC1967UpgradeUpgradeable *ERC1967UpgradeUpgradeableFilterer) FilterInitialized(opts *bind.FilterOpts) (*ERC1967UpgradeUpgradeableInitializedIterator, error) {
	logs, sub, err := _ERC1967UpgradeUpgradeable.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &ERC1967UpgradeUpgradeableInitializedIterator{contract: _ERC1967UpgradeUpgradeable.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_ERC1967UpgradeUpgradeable *ERC1967UpgradeUpgradeableFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *ERC1967UpgradeUpgradeableInitialized) (event.Subscription, error) {
	logs, sub, err := _ERC1967UpgradeUpgradeable.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC1967UpgradeUpgradeableInitialized)
				if err := _ERC1967UpgradeUpgradeable.contract.UnpackLog(event, "Initialized", log); err != nil {
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

// ParseInitialized is a log parse operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_ERC1967UpgradeUpgradeable *ERC1967UpgradeUpgradeableFilterer) ParseInitialized(log types.Log) (*ERC1967UpgradeUpgradeableInitialized, error) {
	event := new(ERC1967UpgradeUpgradeableInitialized)
	if err := _ERC1967UpgradeUpgradeable.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ERC1967UpgradeUpgradeableUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the ERC1967UpgradeUpgradeable contract.
type ERC1967UpgradeUpgradeableUpgradedIterator struct {
	Event *ERC1967UpgradeUpgradeableUpgraded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log      // Log channel receiving the found contract events
	sub  klaytn.Subscription // Subscription for errors, completion and termination
	done bool                // Whether the subscription completed delivering logs
	fail error               // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ERC1967UpgradeUpgradeableUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC1967UpgradeUpgradeableUpgraded)
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
		it.Event = new(ERC1967UpgradeUpgradeableUpgraded)
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
func (it *ERC1967UpgradeUpgradeableUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC1967UpgradeUpgradeableUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC1967UpgradeUpgradeableUpgraded represents a Upgraded event raised by the ERC1967UpgradeUpgradeable contract.
type ERC1967UpgradeUpgradeableUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_ERC1967UpgradeUpgradeable *ERC1967UpgradeUpgradeableFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*ERC1967UpgradeUpgradeableUpgradedIterator, error) {
	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _ERC1967UpgradeUpgradeable.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &ERC1967UpgradeUpgradeableUpgradedIterator{contract: _ERC1967UpgradeUpgradeable.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_ERC1967UpgradeUpgradeable *ERC1967UpgradeUpgradeableFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *ERC1967UpgradeUpgradeableUpgraded, implementation []common.Address) (event.Subscription, error) {
	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _ERC1967UpgradeUpgradeable.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC1967UpgradeUpgradeableUpgraded)
				if err := _ERC1967UpgradeUpgradeable.contract.UnpackLog(event, "Upgraded", log); err != nil {
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

// ParseUpgraded is a log parse operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_ERC1967UpgradeUpgradeable *ERC1967UpgradeUpgradeableFilterer) ParseUpgraded(log types.Log) (*ERC1967UpgradeUpgradeableUpgraded, error) {
	event := new(ERC1967UpgradeUpgradeableUpgraded)
	if err := _ERC1967UpgradeUpgradeable.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	return event, nil
}

// IAddressBookMetaData contains all meta data concerning the IAddressBook contract.
var IAddressBookMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_adminList\",\"type\":\"address[]\"},{\"internalType\":\"uint256\",\"name\":\"_requirement\",\"type\":\"uint256\"}],\"name\":\"constructContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllAddress\",\"outputs\":[{\"internalType\":\"uint8[]\",\"name\":\"typeList\",\"type\":\"uint8[]\"},{\"internalType\":\"address[]\",\"name\":\"addressList\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllAddressInfo\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"cnNodeIdList\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"cnStakingContractList\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"cnRewardAddressList\",\"type\":\"address[]\"},{\"internalType\":\"address\",\"name\":\"pocContractAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"kirContractAddress\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_cnNodeId\",\"type\":\"address\"}],\"name\":\"getCnInfo\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"cnNodeId\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"cnStakingcontract\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"cnRewardAddress\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPendingRequestList\",\"outputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"pendingRequestList\",\"type\":\"bytes32[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_id\",\"type\":\"bytes32\"}],\"name\":\"getRequestInfo\",\"outputs\":[{\"internalType\":\"enumIAddressBook.Functions\",\"name\":\"functionId\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"firstArg\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"secondArg\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"thirdArg\",\"type\":\"bytes32\"},{\"internalType\":\"address[]\",\"name\":\"confirmers\",\"type\":\"address[]\"},{\"internalType\":\"uint256\",\"name\":\"initialProposedTime\",\"type\":\"uint256\"},{\"internalType\":\"enumIAddressBook.RequestState\",\"name\":\"state\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"enumIAddressBook.Functions\",\"name\":\"_functionId\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"_firstArg\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"_secondArg\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"_thirdArg\",\"type\":\"bytes32\"}],\"name\":\"getRequestInfoByArgs\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"},{\"internalType\":\"address[]\",\"name\":\"confirmers\",\"type\":\"address[]\"},{\"internalType\":\"uint256\",\"name\":\"initialProposedTime\",\"type\":\"uint256\"},{\"internalType\":\"enumIAddressBook.RequestState\",\"name\":\"state\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getState\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"adminList\",\"type\":\"address[]\"},{\"internalType\":\"uint256\",\"name\":\"requirement\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"isActivated\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"isConstructed\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"kirContractAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pocContractAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_rewardAddress\",\"type\":\"address\"}],\"name\":\"reviseRewardAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"enumIAddressBook.Functions\",\"name\":\"_functionId\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"_firstArg\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"_secondArg\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"_thirdArg\",\"type\":\"bytes32\"}],\"name\":\"revokeRequest\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"spareContractAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"submitActivateAddressBook\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_admin\",\"type\":\"address\"}],\"name\":\"submitAddAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"submitClearRequest\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_admin\",\"type\":\"address\"}],\"name\":\"submitDeleteAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_cnNodeId\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_cnStakingContractAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_cnRewardAddress\",\"type\":\"address\"}],\"name\":\"submitRegisterCnStakingContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_cnNodeId\",\"type\":\"address\"}],\"name\":\"submitUnregisterCnStakingContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_kirContractAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_version\",\"type\":\"uint256\"}],\"name\":\"submitUpdateKirContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_pocContractAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_version\",\"type\":\"uint256\"}],\"name\":\"submitUpdatePocContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_requirement\",\"type\":\"uint256\"}],\"name\":\"submitUpdateRequirement\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_spareContractAddress\",\"type\":\"address\"}],\"name\":\"submitUpdateSpareContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"7894c366": "constructContract(address[],uint256)",
		"715b208b": "getAllAddress()",
		"160370b8": "getAllAddressInfo()",
		"15575d5a": "getCnInfo(address)",
		"da34a0bd": "getPendingRequestList()",
		"82d67e5a": "getRequestInfo(bytes32)",
		"407091eb": "getRequestInfoByArgs(uint8,bytes32,bytes32,bytes32)",
		"1865c57d": "getState()",
		"4a8c1fb4": "isActivated()",
		"50a5bb69": "isConstructed()",
		"b858dd95": "kirContractAddress()",
		"d267eda5": "pocContractAddress()",
		"832a2aad": "reviseRewardAddress(address)",
		"3f0628b1": "revokeRequest(uint8,bytes32,bytes32,bytes32)",
		"6abd623d": "spareContractAddress()",
		"feb15ca1": "submitActivateAddressBook()",
		"863f5c0a": "submitAddAdmin(address)",
		"87cd9feb": "submitClearRequest()",
		"791b5123": "submitDeleteAdmin(address)",
		"cc11efc0": "submitRegisterCnStakingContract(address,address,address)",
		"b5067706": "submitUnregisterCnStakingContract(address)",
		"9258d768": "submitUpdateKirContract(address,uint256)",
		"21ac4ad4": "submitUpdatePocContract(address,uint256)",
		"e748357b": "submitUpdateRequirement(uint256)",
		"394a144a": "submitUpdateSpareContract(address)",
	},
}

// IAddressBookABI is the input ABI used to generate the binding from.
// Deprecated: Use IAddressBookMetaData.ABI instead.
var IAddressBookABI = IAddressBookMetaData.ABI

// IAddressBookBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const IAddressBookBinRuntime = ``

// IAddressBookFuncSigs maps the 4-byte function signature to its string representation.
// Deprecated: Use IAddressBookMetaData.Sigs instead.
var IAddressBookFuncSigs = IAddressBookMetaData.Sigs

// IAddressBook is an auto generated Go binding around a Klaytn contract.
type IAddressBook struct {
	IAddressBookCaller     // Read-only binding to the contract
	IAddressBookTransactor // Write-only binding to the contract
	IAddressBookFilterer   // Log filterer for contract events
}

// IAddressBookCaller is an auto generated read-only Go binding around a Klaytn contract.
type IAddressBookCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IAddressBookTransactor is an auto generated write-only Go binding around a Klaytn contract.
type IAddressBookTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IAddressBookFilterer is an auto generated log filtering Go binding around a Klaytn contract events.
type IAddressBookFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IAddressBookSession is an auto generated Go binding around a Klaytn contract,
// with pre-set call and transact options.
type IAddressBookSession struct {
	Contract     *IAddressBook     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IAddressBookCallerSession is an auto generated read-only Go binding around a Klaytn contract,
// with pre-set call options.
type IAddressBookCallerSession struct {
	Contract *IAddressBookCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// IAddressBookTransactorSession is an auto generated write-only Go binding around a Klaytn contract,
// with pre-set transact options.
type IAddressBookTransactorSession struct {
	Contract     *IAddressBookTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// IAddressBookRaw is an auto generated low-level Go binding around a Klaytn contract.
type IAddressBookRaw struct {
	Contract *IAddressBook // Generic contract binding to access the raw methods on
}

// IAddressBookCallerRaw is an auto generated low-level read-only Go binding around a Klaytn contract.
type IAddressBookCallerRaw struct {
	Contract *IAddressBookCaller // Generic read-only contract binding to access the raw methods on
}

// IAddressBookTransactorRaw is an auto generated low-level write-only Go binding around a Klaytn contract.
type IAddressBookTransactorRaw struct {
	Contract *IAddressBookTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIAddressBook creates a new instance of IAddressBook, bound to a specific deployed contract.
func NewIAddressBook(address common.Address, backend bind.ContractBackend) (*IAddressBook, error) {
	contract, err := bindIAddressBook(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IAddressBook{IAddressBookCaller: IAddressBookCaller{contract: contract}, IAddressBookTransactor: IAddressBookTransactor{contract: contract}, IAddressBookFilterer: IAddressBookFilterer{contract: contract}}, nil
}

// NewIAddressBookCaller creates a new read-only instance of IAddressBook, bound to a specific deployed contract.
func NewIAddressBookCaller(address common.Address, caller bind.ContractCaller) (*IAddressBookCaller, error) {
	contract, err := bindIAddressBook(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IAddressBookCaller{contract: contract}, nil
}

// NewIAddressBookTransactor creates a new write-only instance of IAddressBook, bound to a specific deployed contract.
func NewIAddressBookTransactor(address common.Address, transactor bind.ContractTransactor) (*IAddressBookTransactor, error) {
	contract, err := bindIAddressBook(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IAddressBookTransactor{contract: contract}, nil
}

// NewIAddressBookFilterer creates a new log filterer instance of IAddressBook, bound to a specific deployed contract.
func NewIAddressBookFilterer(address common.Address, filterer bind.ContractFilterer) (*IAddressBookFilterer, error) {
	contract, err := bindIAddressBook(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IAddressBookFilterer{contract: contract}, nil
}

// bindIAddressBook binds a generic wrapper to an already deployed contract.
func bindIAddressBook(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := IAddressBookMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IAddressBook *IAddressBookRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IAddressBook.Contract.IAddressBookCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IAddressBook *IAddressBookRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IAddressBook.Contract.IAddressBookTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IAddressBook *IAddressBookRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IAddressBook.Contract.IAddressBookTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IAddressBook *IAddressBookCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IAddressBook.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IAddressBook *IAddressBookTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IAddressBook.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IAddressBook *IAddressBookTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IAddressBook.Contract.contract.Transact(opts, method, params...)
}

// GetAllAddress is a free data retrieval call binding the contract method 0x715b208b.
//
// Solidity: function getAllAddress() view returns(uint8[] typeList, address[] addressList)
func (_IAddressBook *IAddressBookCaller) GetAllAddress(opts *bind.CallOpts) (struct {
	TypeList    []uint8
	AddressList []common.Address
}, error,
) {
	var out []interface{}
	err := _IAddressBook.contract.Call(opts, &out, "getAllAddress")

	outstruct := new(struct {
		TypeList    []uint8
		AddressList []common.Address
	})

	outstruct.TypeList = *abi.ConvertType(out[0], new([]uint8)).(*[]uint8)
	outstruct.AddressList = *abi.ConvertType(out[1], new([]common.Address)).(*[]common.Address)
	return *outstruct, err
}

// GetAllAddress is a free data retrieval call binding the contract method 0x715b208b.
//
// Solidity: function getAllAddress() view returns(uint8[] typeList, address[] addressList)
func (_IAddressBook *IAddressBookSession) GetAllAddress() (struct {
	TypeList    []uint8
	AddressList []common.Address
}, error,
) {
	return _IAddressBook.Contract.GetAllAddress(&_IAddressBook.CallOpts)
}

// GetAllAddress is a free data retrieval call binding the contract method 0x715b208b.
//
// Solidity: function getAllAddress() view returns(uint8[] typeList, address[] addressList)
func (_IAddressBook *IAddressBookCallerSession) GetAllAddress() (struct {
	TypeList    []uint8
	AddressList []common.Address
}, error,
) {
	return _IAddressBook.Contract.GetAllAddress(&_IAddressBook.CallOpts)
}

// GetAllAddressInfo is a free data retrieval call binding the contract method 0x160370b8.
//
// Solidity: function getAllAddressInfo() view returns(address[] cnNodeIdList, address[] cnStakingContractList, address[] cnRewardAddressList, address pocContractAddress, address kirContractAddress)
func (_IAddressBook *IAddressBookCaller) GetAllAddressInfo(opts *bind.CallOpts) (struct {
	CnNodeIdList          []common.Address
	CnStakingContractList []common.Address
	CnRewardAddressList   []common.Address
	PocContractAddress    common.Address
	KirContractAddress    common.Address
}, error,
) {
	var out []interface{}
	err := _IAddressBook.contract.Call(opts, &out, "getAllAddressInfo")

	outstruct := new(struct {
		CnNodeIdList          []common.Address
		CnStakingContractList []common.Address
		CnRewardAddressList   []common.Address
		PocContractAddress    common.Address
		KirContractAddress    common.Address
	})

	outstruct.CnNodeIdList = *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)
	outstruct.CnStakingContractList = *abi.ConvertType(out[1], new([]common.Address)).(*[]common.Address)
	outstruct.CnRewardAddressList = *abi.ConvertType(out[2], new([]common.Address)).(*[]common.Address)
	outstruct.PocContractAddress = *abi.ConvertType(out[3], new(common.Address)).(*common.Address)
	outstruct.KirContractAddress = *abi.ConvertType(out[4], new(common.Address)).(*common.Address)
	return *outstruct, err
}

// GetAllAddressInfo is a free data retrieval call binding the contract method 0x160370b8.
//
// Solidity: function getAllAddressInfo() view returns(address[] cnNodeIdList, address[] cnStakingContractList, address[] cnRewardAddressList, address pocContractAddress, address kirContractAddress)
func (_IAddressBook *IAddressBookSession) GetAllAddressInfo() (struct {
	CnNodeIdList          []common.Address
	CnStakingContractList []common.Address
	CnRewardAddressList   []common.Address
	PocContractAddress    common.Address
	KirContractAddress    common.Address
}, error,
) {
	return _IAddressBook.Contract.GetAllAddressInfo(&_IAddressBook.CallOpts)
}

// GetAllAddressInfo is a free data retrieval call binding the contract method 0x160370b8.
//
// Solidity: function getAllAddressInfo() view returns(address[] cnNodeIdList, address[] cnStakingContractList, address[] cnRewardAddressList, address pocContractAddress, address kirContractAddress)
func (_IAddressBook *IAddressBookCallerSession) GetAllAddressInfo() (struct {
	CnNodeIdList          []common.Address
	CnStakingContractList []common.Address
	CnRewardAddressList   []common.Address
	PocContractAddress    common.Address
	KirContractAddress    common.Address
}, error,
) {
	return _IAddressBook.Contract.GetAllAddressInfo(&_IAddressBook.CallOpts)
}

// GetCnInfo is a free data retrieval call binding the contract method 0x15575d5a.
//
// Solidity: function getCnInfo(address _cnNodeId) view returns(address cnNodeId, address cnStakingcontract, address cnRewardAddress)
func (_IAddressBook *IAddressBookCaller) GetCnInfo(opts *bind.CallOpts, _cnNodeId common.Address) (struct {
	CnNodeId          common.Address
	CnStakingcontract common.Address
	CnRewardAddress   common.Address
}, error,
) {
	var out []interface{}
	err := _IAddressBook.contract.Call(opts, &out, "getCnInfo", _cnNodeId)

	outstruct := new(struct {
		CnNodeId          common.Address
		CnStakingcontract common.Address
		CnRewardAddress   common.Address
	})

	outstruct.CnNodeId = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.CnStakingcontract = *abi.ConvertType(out[1], new(common.Address)).(*common.Address)
	outstruct.CnRewardAddress = *abi.ConvertType(out[2], new(common.Address)).(*common.Address)
	return *outstruct, err
}

// GetCnInfo is a free data retrieval call binding the contract method 0x15575d5a.
//
// Solidity: function getCnInfo(address _cnNodeId) view returns(address cnNodeId, address cnStakingcontract, address cnRewardAddress)
func (_IAddressBook *IAddressBookSession) GetCnInfo(_cnNodeId common.Address) (struct {
	CnNodeId          common.Address
	CnStakingcontract common.Address
	CnRewardAddress   common.Address
}, error,
) {
	return _IAddressBook.Contract.GetCnInfo(&_IAddressBook.CallOpts, _cnNodeId)
}

// GetCnInfo is a free data retrieval call binding the contract method 0x15575d5a.
//
// Solidity: function getCnInfo(address _cnNodeId) view returns(address cnNodeId, address cnStakingcontract, address cnRewardAddress)
func (_IAddressBook *IAddressBookCallerSession) GetCnInfo(_cnNodeId common.Address) (struct {
	CnNodeId          common.Address
	CnStakingcontract common.Address
	CnRewardAddress   common.Address
}, error,
) {
	return _IAddressBook.Contract.GetCnInfo(&_IAddressBook.CallOpts, _cnNodeId)
}

// GetPendingRequestList is a free data retrieval call binding the contract method 0xda34a0bd.
//
// Solidity: function getPendingRequestList() view returns(bytes32[] pendingRequestList)
func (_IAddressBook *IAddressBookCaller) GetPendingRequestList(opts *bind.CallOpts) ([][32]byte, error) {
	var out []interface{}
	err := _IAddressBook.contract.Call(opts, &out, "getPendingRequestList")
	if err != nil {
		return *new([][32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([][32]byte)).(*[][32]byte)

	return out0, err
}

// GetPendingRequestList is a free data retrieval call binding the contract method 0xda34a0bd.
//
// Solidity: function getPendingRequestList() view returns(bytes32[] pendingRequestList)
func (_IAddressBook *IAddressBookSession) GetPendingRequestList() ([][32]byte, error) {
	return _IAddressBook.Contract.GetPendingRequestList(&_IAddressBook.CallOpts)
}

// GetPendingRequestList is a free data retrieval call binding the contract method 0xda34a0bd.
//
// Solidity: function getPendingRequestList() view returns(bytes32[] pendingRequestList)
func (_IAddressBook *IAddressBookCallerSession) GetPendingRequestList() ([][32]byte, error) {
	return _IAddressBook.Contract.GetPendingRequestList(&_IAddressBook.CallOpts)
}

// GetRequestInfo is a free data retrieval call binding the contract method 0x82d67e5a.
//
// Solidity: function getRequestInfo(bytes32 _id) view returns(uint8 functionId, bytes32 firstArg, bytes32 secondArg, bytes32 thirdArg, address[] confirmers, uint256 initialProposedTime, uint8 state)
func (_IAddressBook *IAddressBookCaller) GetRequestInfo(opts *bind.CallOpts, _id [32]byte) (struct {
	FunctionId          uint8
	FirstArg            [32]byte
	SecondArg           [32]byte
	ThirdArg            [32]byte
	Confirmers          []common.Address
	InitialProposedTime *big.Int
	State               uint8
}, error,
) {
	var out []interface{}
	err := _IAddressBook.contract.Call(opts, &out, "getRequestInfo", _id)

	outstruct := new(struct {
		FunctionId          uint8
		FirstArg            [32]byte
		SecondArg           [32]byte
		ThirdArg            [32]byte
		Confirmers          []common.Address
		InitialProposedTime *big.Int
		State               uint8
	})

	outstruct.FunctionId = *abi.ConvertType(out[0], new(uint8)).(*uint8)
	outstruct.FirstArg = *abi.ConvertType(out[1], new([32]byte)).(*[32]byte)
	outstruct.SecondArg = *abi.ConvertType(out[2], new([32]byte)).(*[32]byte)
	outstruct.ThirdArg = *abi.ConvertType(out[3], new([32]byte)).(*[32]byte)
	outstruct.Confirmers = *abi.ConvertType(out[4], new([]common.Address)).(*[]common.Address)
	outstruct.InitialProposedTime = *abi.ConvertType(out[5], new(*big.Int)).(**big.Int)
	outstruct.State = *abi.ConvertType(out[6], new(uint8)).(*uint8)
	return *outstruct, err
}

// GetRequestInfo is a free data retrieval call binding the contract method 0x82d67e5a.
//
// Solidity: function getRequestInfo(bytes32 _id) view returns(uint8 functionId, bytes32 firstArg, bytes32 secondArg, bytes32 thirdArg, address[] confirmers, uint256 initialProposedTime, uint8 state)
func (_IAddressBook *IAddressBookSession) GetRequestInfo(_id [32]byte) (struct {
	FunctionId          uint8
	FirstArg            [32]byte
	SecondArg           [32]byte
	ThirdArg            [32]byte
	Confirmers          []common.Address
	InitialProposedTime *big.Int
	State               uint8
}, error,
) {
	return _IAddressBook.Contract.GetRequestInfo(&_IAddressBook.CallOpts, _id)
}

// GetRequestInfo is a free data retrieval call binding the contract method 0x82d67e5a.
//
// Solidity: function getRequestInfo(bytes32 _id) view returns(uint8 functionId, bytes32 firstArg, bytes32 secondArg, bytes32 thirdArg, address[] confirmers, uint256 initialProposedTime, uint8 state)
func (_IAddressBook *IAddressBookCallerSession) GetRequestInfo(_id [32]byte) (struct {
	FunctionId          uint8
	FirstArg            [32]byte
	SecondArg           [32]byte
	ThirdArg            [32]byte
	Confirmers          []common.Address
	InitialProposedTime *big.Int
	State               uint8
}, error,
) {
	return _IAddressBook.Contract.GetRequestInfo(&_IAddressBook.CallOpts, _id)
}

// GetRequestInfoByArgs is a free data retrieval call binding the contract method 0x407091eb.
//
// Solidity: function getRequestInfoByArgs(uint8 _functionId, bytes32 _firstArg, bytes32 _secondArg, bytes32 _thirdArg) view returns(bytes32 id, address[] confirmers, uint256 initialProposedTime, uint8 state)
func (_IAddressBook *IAddressBookCaller) GetRequestInfoByArgs(opts *bind.CallOpts, _functionId uint8, _firstArg [32]byte, _secondArg [32]byte, _thirdArg [32]byte) (struct {
	Id                  [32]byte
	Confirmers          []common.Address
	InitialProposedTime *big.Int
	State               uint8
}, error,
) {
	var out []interface{}
	err := _IAddressBook.contract.Call(opts, &out, "getRequestInfoByArgs", _functionId, _firstArg, _secondArg, _thirdArg)

	outstruct := new(struct {
		Id                  [32]byte
		Confirmers          []common.Address
		InitialProposedTime *big.Int
		State               uint8
	})

	outstruct.Id = *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	outstruct.Confirmers = *abi.ConvertType(out[1], new([]common.Address)).(*[]common.Address)
	outstruct.InitialProposedTime = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.State = *abi.ConvertType(out[3], new(uint8)).(*uint8)
	return *outstruct, err
}

// GetRequestInfoByArgs is a free data retrieval call binding the contract method 0x407091eb.
//
// Solidity: function getRequestInfoByArgs(uint8 _functionId, bytes32 _firstArg, bytes32 _secondArg, bytes32 _thirdArg) view returns(bytes32 id, address[] confirmers, uint256 initialProposedTime, uint8 state)
func (_IAddressBook *IAddressBookSession) GetRequestInfoByArgs(_functionId uint8, _firstArg [32]byte, _secondArg [32]byte, _thirdArg [32]byte) (struct {
	Id                  [32]byte
	Confirmers          []common.Address
	InitialProposedTime *big.Int
	State               uint8
}, error,
) {
	return _IAddressBook.Contract.GetRequestInfoByArgs(&_IAddressBook.CallOpts, _functionId, _firstArg, _secondArg, _thirdArg)
}

// GetRequestInfoByArgs is a free data retrieval call binding the contract method 0x407091eb.
//
// Solidity: function getRequestInfoByArgs(uint8 _functionId, bytes32 _firstArg, bytes32 _secondArg, bytes32 _thirdArg) view returns(bytes32 id, address[] confirmers, uint256 initialProposedTime, uint8 state)
func (_IAddressBook *IAddressBookCallerSession) GetRequestInfoByArgs(_functionId uint8, _firstArg [32]byte, _secondArg [32]byte, _thirdArg [32]byte) (struct {
	Id                  [32]byte
	Confirmers          []common.Address
	InitialProposedTime *big.Int
	State               uint8
}, error,
) {
	return _IAddressBook.Contract.GetRequestInfoByArgs(&_IAddressBook.CallOpts, _functionId, _firstArg, _secondArg, _thirdArg)
}

// GetState is a free data retrieval call binding the contract method 0x1865c57d.
//
// Solidity: function getState() view returns(address[] adminList, uint256 requirement)
func (_IAddressBook *IAddressBookCaller) GetState(opts *bind.CallOpts) (struct {
	AdminList   []common.Address
	Requirement *big.Int
}, error,
) {
	var out []interface{}
	err := _IAddressBook.contract.Call(opts, &out, "getState")

	outstruct := new(struct {
		AdminList   []common.Address
		Requirement *big.Int
	})

	outstruct.AdminList = *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)
	outstruct.Requirement = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	return *outstruct, err
}

// GetState is a free data retrieval call binding the contract method 0x1865c57d.
//
// Solidity: function getState() view returns(address[] adminList, uint256 requirement)
func (_IAddressBook *IAddressBookSession) GetState() (struct {
	AdminList   []common.Address
	Requirement *big.Int
}, error,
) {
	return _IAddressBook.Contract.GetState(&_IAddressBook.CallOpts)
}

// GetState is a free data retrieval call binding the contract method 0x1865c57d.
//
// Solidity: function getState() view returns(address[] adminList, uint256 requirement)
func (_IAddressBook *IAddressBookCallerSession) GetState() (struct {
	AdminList   []common.Address
	Requirement *big.Int
}, error,
) {
	return _IAddressBook.Contract.GetState(&_IAddressBook.CallOpts)
}

// IsActivated is a free data retrieval call binding the contract method 0x4a8c1fb4.
//
// Solidity: function isActivated() view returns(bool)
func (_IAddressBook *IAddressBookCaller) IsActivated(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _IAddressBook.contract.Call(opts, &out, "isActivated")
	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err
}

// IsActivated is a free data retrieval call binding the contract method 0x4a8c1fb4.
//
// Solidity: function isActivated() view returns(bool)
func (_IAddressBook *IAddressBookSession) IsActivated() (bool, error) {
	return _IAddressBook.Contract.IsActivated(&_IAddressBook.CallOpts)
}

// IsActivated is a free data retrieval call binding the contract method 0x4a8c1fb4.
//
// Solidity: function isActivated() view returns(bool)
func (_IAddressBook *IAddressBookCallerSession) IsActivated() (bool, error) {
	return _IAddressBook.Contract.IsActivated(&_IAddressBook.CallOpts)
}

// IsConstructed is a free data retrieval call binding the contract method 0x50a5bb69.
//
// Solidity: function isConstructed() view returns(bool)
func (_IAddressBook *IAddressBookCaller) IsConstructed(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _IAddressBook.contract.Call(opts, &out, "isConstructed")
	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err
}

// IsConstructed is a free data retrieval call binding the contract method 0x50a5bb69.
//
// Solidity: function isConstructed() view returns(bool)
func (_IAddressBook *IAddressBookSession) IsConstructed() (bool, error) {
	return _IAddressBook.Contract.IsConstructed(&_IAddressBook.CallOpts)
}

// IsConstructed is a free data retrieval call binding the contract method 0x50a5bb69.
//
// Solidity: function isConstructed() view returns(bool)
func (_IAddressBook *IAddressBookCallerSession) IsConstructed() (bool, error) {
	return _IAddressBook.Contract.IsConstructed(&_IAddressBook.CallOpts)
}

// KirContractAddress is a free data retrieval call binding the contract method 0xb858dd95.
//
// Solidity: function kirContractAddress() view returns(address)
func (_IAddressBook *IAddressBookCaller) KirContractAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IAddressBook.contract.Call(opts, &out, "kirContractAddress")
	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err
}

// KirContractAddress is a free data retrieval call binding the contract method 0xb858dd95.
//
// Solidity: function kirContractAddress() view returns(address)
func (_IAddressBook *IAddressBookSession) KirContractAddress() (common.Address, error) {
	return _IAddressBook.Contract.KirContractAddress(&_IAddressBook.CallOpts)
}

// KirContractAddress is a free data retrieval call binding the contract method 0xb858dd95.
//
// Solidity: function kirContractAddress() view returns(address)
func (_IAddressBook *IAddressBookCallerSession) KirContractAddress() (common.Address, error) {
	return _IAddressBook.Contract.KirContractAddress(&_IAddressBook.CallOpts)
}

// PocContractAddress is a free data retrieval call binding the contract method 0xd267eda5.
//
// Solidity: function pocContractAddress() view returns(address)
func (_IAddressBook *IAddressBookCaller) PocContractAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IAddressBook.contract.Call(opts, &out, "pocContractAddress")
	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err
}

// PocContractAddress is a free data retrieval call binding the contract method 0xd267eda5.
//
// Solidity: function pocContractAddress() view returns(address)
func (_IAddressBook *IAddressBookSession) PocContractAddress() (common.Address, error) {
	return _IAddressBook.Contract.PocContractAddress(&_IAddressBook.CallOpts)
}

// PocContractAddress is a free data retrieval call binding the contract method 0xd267eda5.
//
// Solidity: function pocContractAddress() view returns(address)
func (_IAddressBook *IAddressBookCallerSession) PocContractAddress() (common.Address, error) {
	return _IAddressBook.Contract.PocContractAddress(&_IAddressBook.CallOpts)
}

// SpareContractAddress is a free data retrieval call binding the contract method 0x6abd623d.
//
// Solidity: function spareContractAddress() view returns(address)
func (_IAddressBook *IAddressBookCaller) SpareContractAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IAddressBook.contract.Call(opts, &out, "spareContractAddress")
	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err
}

// SpareContractAddress is a free data retrieval call binding the contract method 0x6abd623d.
//
// Solidity: function spareContractAddress() view returns(address)
func (_IAddressBook *IAddressBookSession) SpareContractAddress() (common.Address, error) {
	return _IAddressBook.Contract.SpareContractAddress(&_IAddressBook.CallOpts)
}

// SpareContractAddress is a free data retrieval call binding the contract method 0x6abd623d.
//
// Solidity: function spareContractAddress() view returns(address)
func (_IAddressBook *IAddressBookCallerSession) SpareContractAddress() (common.Address, error) {
	return _IAddressBook.Contract.SpareContractAddress(&_IAddressBook.CallOpts)
}

// ConstructContract is a paid mutator transaction binding the contract method 0x7894c366.
//
// Solidity: function constructContract(address[] _adminList, uint256 _requirement) returns()
func (_IAddressBook *IAddressBookTransactor) ConstructContract(opts *bind.TransactOpts, _adminList []common.Address, _requirement *big.Int) (*types.Transaction, error) {
	return _IAddressBook.contract.Transact(opts, "constructContract", _adminList, _requirement)
}

// ConstructContract is a paid mutator transaction binding the contract method 0x7894c366.
//
// Solidity: function constructContract(address[] _adminList, uint256 _requirement) returns()
func (_IAddressBook *IAddressBookSession) ConstructContract(_adminList []common.Address, _requirement *big.Int) (*types.Transaction, error) {
	return _IAddressBook.Contract.ConstructContract(&_IAddressBook.TransactOpts, _adminList, _requirement)
}

// ConstructContract is a paid mutator transaction binding the contract method 0x7894c366.
//
// Solidity: function constructContract(address[] _adminList, uint256 _requirement) returns()
func (_IAddressBook *IAddressBookTransactorSession) ConstructContract(_adminList []common.Address, _requirement *big.Int) (*types.Transaction, error) {
	return _IAddressBook.Contract.ConstructContract(&_IAddressBook.TransactOpts, _adminList, _requirement)
}

// ReviseRewardAddress is a paid mutator transaction binding the contract method 0x832a2aad.
//
// Solidity: function reviseRewardAddress(address _rewardAddress) returns()
func (_IAddressBook *IAddressBookTransactor) ReviseRewardAddress(opts *bind.TransactOpts, _rewardAddress common.Address) (*types.Transaction, error) {
	return _IAddressBook.contract.Transact(opts, "reviseRewardAddress", _rewardAddress)
}

// ReviseRewardAddress is a paid mutator transaction binding the contract method 0x832a2aad.
//
// Solidity: function reviseRewardAddress(address _rewardAddress) returns()
func (_IAddressBook *IAddressBookSession) ReviseRewardAddress(_rewardAddress common.Address) (*types.Transaction, error) {
	return _IAddressBook.Contract.ReviseRewardAddress(&_IAddressBook.TransactOpts, _rewardAddress)
}

// ReviseRewardAddress is a paid mutator transaction binding the contract method 0x832a2aad.
//
// Solidity: function reviseRewardAddress(address _rewardAddress) returns()
func (_IAddressBook *IAddressBookTransactorSession) ReviseRewardAddress(_rewardAddress common.Address) (*types.Transaction, error) {
	return _IAddressBook.Contract.ReviseRewardAddress(&_IAddressBook.TransactOpts, _rewardAddress)
}

// RevokeRequest is a paid mutator transaction binding the contract method 0x3f0628b1.
//
// Solidity: function revokeRequest(uint8 _functionId, bytes32 _firstArg, bytes32 _secondArg, bytes32 _thirdArg) returns()
func (_IAddressBook *IAddressBookTransactor) RevokeRequest(opts *bind.TransactOpts, _functionId uint8, _firstArg [32]byte, _secondArg [32]byte, _thirdArg [32]byte) (*types.Transaction, error) {
	return _IAddressBook.contract.Transact(opts, "revokeRequest", _functionId, _firstArg, _secondArg, _thirdArg)
}

// RevokeRequest is a paid mutator transaction binding the contract method 0x3f0628b1.
//
// Solidity: function revokeRequest(uint8 _functionId, bytes32 _firstArg, bytes32 _secondArg, bytes32 _thirdArg) returns()
func (_IAddressBook *IAddressBookSession) RevokeRequest(_functionId uint8, _firstArg [32]byte, _secondArg [32]byte, _thirdArg [32]byte) (*types.Transaction, error) {
	return _IAddressBook.Contract.RevokeRequest(&_IAddressBook.TransactOpts, _functionId, _firstArg, _secondArg, _thirdArg)
}

// RevokeRequest is a paid mutator transaction binding the contract method 0x3f0628b1.
//
// Solidity: function revokeRequest(uint8 _functionId, bytes32 _firstArg, bytes32 _secondArg, bytes32 _thirdArg) returns()
func (_IAddressBook *IAddressBookTransactorSession) RevokeRequest(_functionId uint8, _firstArg [32]byte, _secondArg [32]byte, _thirdArg [32]byte) (*types.Transaction, error) {
	return _IAddressBook.Contract.RevokeRequest(&_IAddressBook.TransactOpts, _functionId, _firstArg, _secondArg, _thirdArg)
}

// SubmitActivateAddressBook is a paid mutator transaction binding the contract method 0xfeb15ca1.
//
// Solidity: function submitActivateAddressBook() returns()
func (_IAddressBook *IAddressBookTransactor) SubmitActivateAddressBook(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IAddressBook.contract.Transact(opts, "submitActivateAddressBook")
}

// SubmitActivateAddressBook is a paid mutator transaction binding the contract method 0xfeb15ca1.
//
// Solidity: function submitActivateAddressBook() returns()
func (_IAddressBook *IAddressBookSession) SubmitActivateAddressBook() (*types.Transaction, error) {
	return _IAddressBook.Contract.SubmitActivateAddressBook(&_IAddressBook.TransactOpts)
}

// SubmitActivateAddressBook is a paid mutator transaction binding the contract method 0xfeb15ca1.
//
// Solidity: function submitActivateAddressBook() returns()
func (_IAddressBook *IAddressBookTransactorSession) SubmitActivateAddressBook() (*types.Transaction, error) {
	return _IAddressBook.Contract.SubmitActivateAddressBook(&_IAddressBook.TransactOpts)
}

// SubmitAddAdmin is a paid mutator transaction binding the contract method 0x863f5c0a.
//
// Solidity: function submitAddAdmin(address _admin) returns()
func (_IAddressBook *IAddressBookTransactor) SubmitAddAdmin(opts *bind.TransactOpts, _admin common.Address) (*types.Transaction, error) {
	return _IAddressBook.contract.Transact(opts, "submitAddAdmin", _admin)
}

// SubmitAddAdmin is a paid mutator transaction binding the contract method 0x863f5c0a.
//
// Solidity: function submitAddAdmin(address _admin) returns()
func (_IAddressBook *IAddressBookSession) SubmitAddAdmin(_admin common.Address) (*types.Transaction, error) {
	return _IAddressBook.Contract.SubmitAddAdmin(&_IAddressBook.TransactOpts, _admin)
}

// SubmitAddAdmin is a paid mutator transaction binding the contract method 0x863f5c0a.
//
// Solidity: function submitAddAdmin(address _admin) returns()
func (_IAddressBook *IAddressBookTransactorSession) SubmitAddAdmin(_admin common.Address) (*types.Transaction, error) {
	return _IAddressBook.Contract.SubmitAddAdmin(&_IAddressBook.TransactOpts, _admin)
}

// SubmitClearRequest is a paid mutator transaction binding the contract method 0x87cd9feb.
//
// Solidity: function submitClearRequest() returns()
func (_IAddressBook *IAddressBookTransactor) SubmitClearRequest(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IAddressBook.contract.Transact(opts, "submitClearRequest")
}

// SubmitClearRequest is a paid mutator transaction binding the contract method 0x87cd9feb.
//
// Solidity: function submitClearRequest() returns()
func (_IAddressBook *IAddressBookSession) SubmitClearRequest() (*types.Transaction, error) {
	return _IAddressBook.Contract.SubmitClearRequest(&_IAddressBook.TransactOpts)
}

// SubmitClearRequest is a paid mutator transaction binding the contract method 0x87cd9feb.
//
// Solidity: function submitClearRequest() returns()
func (_IAddressBook *IAddressBookTransactorSession) SubmitClearRequest() (*types.Transaction, error) {
	return _IAddressBook.Contract.SubmitClearRequest(&_IAddressBook.TransactOpts)
}

// SubmitDeleteAdmin is a paid mutator transaction binding the contract method 0x791b5123.
//
// Solidity: function submitDeleteAdmin(address _admin) returns()
func (_IAddressBook *IAddressBookTransactor) SubmitDeleteAdmin(opts *bind.TransactOpts, _admin common.Address) (*types.Transaction, error) {
	return _IAddressBook.contract.Transact(opts, "submitDeleteAdmin", _admin)
}

// SubmitDeleteAdmin is a paid mutator transaction binding the contract method 0x791b5123.
//
// Solidity: function submitDeleteAdmin(address _admin) returns()
func (_IAddressBook *IAddressBookSession) SubmitDeleteAdmin(_admin common.Address) (*types.Transaction, error) {
	return _IAddressBook.Contract.SubmitDeleteAdmin(&_IAddressBook.TransactOpts, _admin)
}

// SubmitDeleteAdmin is a paid mutator transaction binding the contract method 0x791b5123.
//
// Solidity: function submitDeleteAdmin(address _admin) returns()
func (_IAddressBook *IAddressBookTransactorSession) SubmitDeleteAdmin(_admin common.Address) (*types.Transaction, error) {
	return _IAddressBook.Contract.SubmitDeleteAdmin(&_IAddressBook.TransactOpts, _admin)
}

// SubmitRegisterCnStakingContract is a paid mutator transaction binding the contract method 0xcc11efc0.
//
// Solidity: function submitRegisterCnStakingContract(address _cnNodeId, address _cnStakingContractAddress, address _cnRewardAddress) returns()
func (_IAddressBook *IAddressBookTransactor) SubmitRegisterCnStakingContract(opts *bind.TransactOpts, _cnNodeId common.Address, _cnStakingContractAddress common.Address, _cnRewardAddress common.Address) (*types.Transaction, error) {
	return _IAddressBook.contract.Transact(opts, "submitRegisterCnStakingContract", _cnNodeId, _cnStakingContractAddress, _cnRewardAddress)
}

// SubmitRegisterCnStakingContract is a paid mutator transaction binding the contract method 0xcc11efc0.
//
// Solidity: function submitRegisterCnStakingContract(address _cnNodeId, address _cnStakingContractAddress, address _cnRewardAddress) returns()
func (_IAddressBook *IAddressBookSession) SubmitRegisterCnStakingContract(_cnNodeId common.Address, _cnStakingContractAddress common.Address, _cnRewardAddress common.Address) (*types.Transaction, error) {
	return _IAddressBook.Contract.SubmitRegisterCnStakingContract(&_IAddressBook.TransactOpts, _cnNodeId, _cnStakingContractAddress, _cnRewardAddress)
}

// SubmitRegisterCnStakingContract is a paid mutator transaction binding the contract method 0xcc11efc0.
//
// Solidity: function submitRegisterCnStakingContract(address _cnNodeId, address _cnStakingContractAddress, address _cnRewardAddress) returns()
func (_IAddressBook *IAddressBookTransactorSession) SubmitRegisterCnStakingContract(_cnNodeId common.Address, _cnStakingContractAddress common.Address, _cnRewardAddress common.Address) (*types.Transaction, error) {
	return _IAddressBook.Contract.SubmitRegisterCnStakingContract(&_IAddressBook.TransactOpts, _cnNodeId, _cnStakingContractAddress, _cnRewardAddress)
}

// SubmitUnregisterCnStakingContract is a paid mutator transaction binding the contract method 0xb5067706.
//
// Solidity: function submitUnregisterCnStakingContract(address _cnNodeId) returns()
func (_IAddressBook *IAddressBookTransactor) SubmitUnregisterCnStakingContract(opts *bind.TransactOpts, _cnNodeId common.Address) (*types.Transaction, error) {
	return _IAddressBook.contract.Transact(opts, "submitUnregisterCnStakingContract", _cnNodeId)
}

// SubmitUnregisterCnStakingContract is a paid mutator transaction binding the contract method 0xb5067706.
//
// Solidity: function submitUnregisterCnStakingContract(address _cnNodeId) returns()
func (_IAddressBook *IAddressBookSession) SubmitUnregisterCnStakingContract(_cnNodeId common.Address) (*types.Transaction, error) {
	return _IAddressBook.Contract.SubmitUnregisterCnStakingContract(&_IAddressBook.TransactOpts, _cnNodeId)
}

// SubmitUnregisterCnStakingContract is a paid mutator transaction binding the contract method 0xb5067706.
//
// Solidity: function submitUnregisterCnStakingContract(address _cnNodeId) returns()
func (_IAddressBook *IAddressBookTransactorSession) SubmitUnregisterCnStakingContract(_cnNodeId common.Address) (*types.Transaction, error) {
	return _IAddressBook.Contract.SubmitUnregisterCnStakingContract(&_IAddressBook.TransactOpts, _cnNodeId)
}

// SubmitUpdateKirContract is a paid mutator transaction binding the contract method 0x9258d768.
//
// Solidity: function submitUpdateKirContract(address _kirContractAddress, uint256 _version) returns()
func (_IAddressBook *IAddressBookTransactor) SubmitUpdateKirContract(opts *bind.TransactOpts, _kirContractAddress common.Address, _version *big.Int) (*types.Transaction, error) {
	return _IAddressBook.contract.Transact(opts, "submitUpdateKirContract", _kirContractAddress, _version)
}

// SubmitUpdateKirContract is a paid mutator transaction binding the contract method 0x9258d768.
//
// Solidity: function submitUpdateKirContract(address _kirContractAddress, uint256 _version) returns()
func (_IAddressBook *IAddressBookSession) SubmitUpdateKirContract(_kirContractAddress common.Address, _version *big.Int) (*types.Transaction, error) {
	return _IAddressBook.Contract.SubmitUpdateKirContract(&_IAddressBook.TransactOpts, _kirContractAddress, _version)
}

// SubmitUpdateKirContract is a paid mutator transaction binding the contract method 0x9258d768.
//
// Solidity: function submitUpdateKirContract(address _kirContractAddress, uint256 _version) returns()
func (_IAddressBook *IAddressBookTransactorSession) SubmitUpdateKirContract(_kirContractAddress common.Address, _version *big.Int) (*types.Transaction, error) {
	return _IAddressBook.Contract.SubmitUpdateKirContract(&_IAddressBook.TransactOpts, _kirContractAddress, _version)
}

// SubmitUpdatePocContract is a paid mutator transaction binding the contract method 0x21ac4ad4.
//
// Solidity: function submitUpdatePocContract(address _pocContractAddress, uint256 _version) returns()
func (_IAddressBook *IAddressBookTransactor) SubmitUpdatePocContract(opts *bind.TransactOpts, _pocContractAddress common.Address, _version *big.Int) (*types.Transaction, error) {
	return _IAddressBook.contract.Transact(opts, "submitUpdatePocContract", _pocContractAddress, _version)
}

// SubmitUpdatePocContract is a paid mutator transaction binding the contract method 0x21ac4ad4.
//
// Solidity: function submitUpdatePocContract(address _pocContractAddress, uint256 _version) returns()
func (_IAddressBook *IAddressBookSession) SubmitUpdatePocContract(_pocContractAddress common.Address, _version *big.Int) (*types.Transaction, error) {
	return _IAddressBook.Contract.SubmitUpdatePocContract(&_IAddressBook.TransactOpts, _pocContractAddress, _version)
}

// SubmitUpdatePocContract is a paid mutator transaction binding the contract method 0x21ac4ad4.
//
// Solidity: function submitUpdatePocContract(address _pocContractAddress, uint256 _version) returns()
func (_IAddressBook *IAddressBookTransactorSession) SubmitUpdatePocContract(_pocContractAddress common.Address, _version *big.Int) (*types.Transaction, error) {
	return _IAddressBook.Contract.SubmitUpdatePocContract(&_IAddressBook.TransactOpts, _pocContractAddress, _version)
}

// SubmitUpdateRequirement is a paid mutator transaction binding the contract method 0xe748357b.
//
// Solidity: function submitUpdateRequirement(uint256 _requirement) returns()
func (_IAddressBook *IAddressBookTransactor) SubmitUpdateRequirement(opts *bind.TransactOpts, _requirement *big.Int) (*types.Transaction, error) {
	return _IAddressBook.contract.Transact(opts, "submitUpdateRequirement", _requirement)
}

// SubmitUpdateRequirement is a paid mutator transaction binding the contract method 0xe748357b.
//
// Solidity: function submitUpdateRequirement(uint256 _requirement) returns()
func (_IAddressBook *IAddressBookSession) SubmitUpdateRequirement(_requirement *big.Int) (*types.Transaction, error) {
	return _IAddressBook.Contract.SubmitUpdateRequirement(&_IAddressBook.TransactOpts, _requirement)
}

// SubmitUpdateRequirement is a paid mutator transaction binding the contract method 0xe748357b.
//
// Solidity: function submitUpdateRequirement(uint256 _requirement) returns()
func (_IAddressBook *IAddressBookTransactorSession) SubmitUpdateRequirement(_requirement *big.Int) (*types.Transaction, error) {
	return _IAddressBook.Contract.SubmitUpdateRequirement(&_IAddressBook.TransactOpts, _requirement)
}

// SubmitUpdateSpareContract is a paid mutator transaction binding the contract method 0x394a144a.
//
// Solidity: function submitUpdateSpareContract(address _spareContractAddress) returns()
func (_IAddressBook *IAddressBookTransactor) SubmitUpdateSpareContract(opts *bind.TransactOpts, _spareContractAddress common.Address) (*types.Transaction, error) {
	return _IAddressBook.contract.Transact(opts, "submitUpdateSpareContract", _spareContractAddress)
}

// SubmitUpdateSpareContract is a paid mutator transaction binding the contract method 0x394a144a.
//
// Solidity: function submitUpdateSpareContract(address _spareContractAddress) returns()
func (_IAddressBook *IAddressBookSession) SubmitUpdateSpareContract(_spareContractAddress common.Address) (*types.Transaction, error) {
	return _IAddressBook.Contract.SubmitUpdateSpareContract(&_IAddressBook.TransactOpts, _spareContractAddress)
}

// SubmitUpdateSpareContract is a paid mutator transaction binding the contract method 0x394a144a.
//
// Solidity: function submitUpdateSpareContract(address _spareContractAddress) returns()
func (_IAddressBook *IAddressBookTransactorSession) SubmitUpdateSpareContract(_spareContractAddress common.Address) (*types.Transaction, error) {
	return _IAddressBook.Contract.SubmitUpdateSpareContract(&_IAddressBook.TransactOpts, _spareContractAddress)
}

// IBeaconMetaData contains all meta data concerning the IBeacon contract.
var IBeaconMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"implementation\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"5c60da1b": "implementation()",
	},
}

// IBeaconABI is the input ABI used to generate the binding from.
// Deprecated: Use IBeaconMetaData.ABI instead.
var IBeaconABI = IBeaconMetaData.ABI

// IBeaconBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const IBeaconBinRuntime = ``

// IBeaconFuncSigs maps the 4-byte function signature to its string representation.
// Deprecated: Use IBeaconMetaData.Sigs instead.
var IBeaconFuncSigs = IBeaconMetaData.Sigs

// IBeacon is an auto generated Go binding around a Klaytn contract.
type IBeacon struct {
	IBeaconCaller     // Read-only binding to the contract
	IBeaconTransactor // Write-only binding to the contract
	IBeaconFilterer   // Log filterer for contract events
}

// IBeaconCaller is an auto generated read-only Go binding around a Klaytn contract.
type IBeaconCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IBeaconTransactor is an auto generated write-only Go binding around a Klaytn contract.
type IBeaconTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IBeaconFilterer is an auto generated log filtering Go binding around a Klaytn contract events.
type IBeaconFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IBeaconSession is an auto generated Go binding around a Klaytn contract,
// with pre-set call and transact options.
type IBeaconSession struct {
	Contract     *IBeacon          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IBeaconCallerSession is an auto generated read-only Go binding around a Klaytn contract,
// with pre-set call options.
type IBeaconCallerSession struct {
	Contract *IBeaconCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// IBeaconTransactorSession is an auto generated write-only Go binding around a Klaytn contract,
// with pre-set transact options.
type IBeaconTransactorSession struct {
	Contract     *IBeaconTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// IBeaconRaw is an auto generated low-level Go binding around a Klaytn contract.
type IBeaconRaw struct {
	Contract *IBeacon // Generic contract binding to access the raw methods on
}

// IBeaconCallerRaw is an auto generated low-level read-only Go binding around a Klaytn contract.
type IBeaconCallerRaw struct {
	Contract *IBeaconCaller // Generic read-only contract binding to access the raw methods on
}

// IBeaconTransactorRaw is an auto generated low-level write-only Go binding around a Klaytn contract.
type IBeaconTransactorRaw struct {
	Contract *IBeaconTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIBeacon creates a new instance of IBeacon, bound to a specific deployed contract.
func NewIBeacon(address common.Address, backend bind.ContractBackend) (*IBeacon, error) {
	contract, err := bindIBeacon(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IBeacon{IBeaconCaller: IBeaconCaller{contract: contract}, IBeaconTransactor: IBeaconTransactor{contract: contract}, IBeaconFilterer: IBeaconFilterer{contract: contract}}, nil
}

// NewIBeaconCaller creates a new read-only instance of IBeacon, bound to a specific deployed contract.
func NewIBeaconCaller(address common.Address, caller bind.ContractCaller) (*IBeaconCaller, error) {
	contract, err := bindIBeacon(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IBeaconCaller{contract: contract}, nil
}

// NewIBeaconTransactor creates a new write-only instance of IBeacon, bound to a specific deployed contract.
func NewIBeaconTransactor(address common.Address, transactor bind.ContractTransactor) (*IBeaconTransactor, error) {
	contract, err := bindIBeacon(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IBeaconTransactor{contract: contract}, nil
}

// NewIBeaconFilterer creates a new log filterer instance of IBeacon, bound to a specific deployed contract.
func NewIBeaconFilterer(address common.Address, filterer bind.ContractFilterer) (*IBeaconFilterer, error) {
	contract, err := bindIBeacon(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IBeaconFilterer{contract: contract}, nil
}

// bindIBeacon binds a generic wrapper to an already deployed contract.
func bindIBeacon(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := IBeaconMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IBeacon *IBeaconRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IBeacon.Contract.IBeaconCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IBeacon *IBeaconRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IBeacon.Contract.IBeaconTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IBeacon *IBeaconRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IBeacon.Contract.IBeaconTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IBeacon *IBeaconCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IBeacon.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IBeacon *IBeaconTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IBeacon.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IBeacon *IBeaconTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IBeacon.Contract.contract.Transact(opts, method, params...)
}

// Implementation is a free data retrieval call binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() view returns(address)
func (_IBeacon *IBeaconCaller) Implementation(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IBeacon.contract.Call(opts, &out, "implementation")
	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err
}

// Implementation is a free data retrieval call binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() view returns(address)
func (_IBeacon *IBeaconSession) Implementation() (common.Address, error) {
	return _IBeacon.Contract.Implementation(&_IBeacon.CallOpts)
}

// Implementation is a free data retrieval call binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() view returns(address)
func (_IBeacon *IBeaconCallerSession) Implementation() (common.Address, error) {
	return _IBeacon.Contract.Implementation(&_IBeacon.CallOpts)
}

// IBeaconUpgradeableMetaData contains all meta data concerning the IBeaconUpgradeable contract.
var IBeaconUpgradeableMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"implementation\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"5c60da1b": "implementation()",
	},
}

// IBeaconUpgradeableABI is the input ABI used to generate the binding from.
// Deprecated: Use IBeaconUpgradeableMetaData.ABI instead.
var IBeaconUpgradeableABI = IBeaconUpgradeableMetaData.ABI

// IBeaconUpgradeableBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const IBeaconUpgradeableBinRuntime = ``

// IBeaconUpgradeableFuncSigs maps the 4-byte function signature to its string representation.
// Deprecated: Use IBeaconUpgradeableMetaData.Sigs instead.
var IBeaconUpgradeableFuncSigs = IBeaconUpgradeableMetaData.Sigs

// IBeaconUpgradeable is an auto generated Go binding around a Klaytn contract.
type IBeaconUpgradeable struct {
	IBeaconUpgradeableCaller     // Read-only binding to the contract
	IBeaconUpgradeableTransactor // Write-only binding to the contract
	IBeaconUpgradeableFilterer   // Log filterer for contract events
}

// IBeaconUpgradeableCaller is an auto generated read-only Go binding around a Klaytn contract.
type IBeaconUpgradeableCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IBeaconUpgradeableTransactor is an auto generated write-only Go binding around a Klaytn contract.
type IBeaconUpgradeableTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IBeaconUpgradeableFilterer is an auto generated log filtering Go binding around a Klaytn contract events.
type IBeaconUpgradeableFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IBeaconUpgradeableSession is an auto generated Go binding around a Klaytn contract,
// with pre-set call and transact options.
type IBeaconUpgradeableSession struct {
	Contract     *IBeaconUpgradeable // Generic contract binding to set the session for
	CallOpts     bind.CallOpts       // Call options to use throughout this session
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// IBeaconUpgradeableCallerSession is an auto generated read-only Go binding around a Klaytn contract,
// with pre-set call options.
type IBeaconUpgradeableCallerSession struct {
	Contract *IBeaconUpgradeableCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts             // Call options to use throughout this session
}

// IBeaconUpgradeableTransactorSession is an auto generated write-only Go binding around a Klaytn contract,
// with pre-set transact options.
type IBeaconUpgradeableTransactorSession struct {
	Contract     *IBeaconUpgradeableTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts             // Transaction auth options to use throughout this session
}

// IBeaconUpgradeableRaw is an auto generated low-level Go binding around a Klaytn contract.
type IBeaconUpgradeableRaw struct {
	Contract *IBeaconUpgradeable // Generic contract binding to access the raw methods on
}

// IBeaconUpgradeableCallerRaw is an auto generated low-level read-only Go binding around a Klaytn contract.
type IBeaconUpgradeableCallerRaw struct {
	Contract *IBeaconUpgradeableCaller // Generic read-only contract binding to access the raw methods on
}

// IBeaconUpgradeableTransactorRaw is an auto generated low-level write-only Go binding around a Klaytn contract.
type IBeaconUpgradeableTransactorRaw struct {
	Contract *IBeaconUpgradeableTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIBeaconUpgradeable creates a new instance of IBeaconUpgradeable, bound to a specific deployed contract.
func NewIBeaconUpgradeable(address common.Address, backend bind.ContractBackend) (*IBeaconUpgradeable, error) {
	contract, err := bindIBeaconUpgradeable(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IBeaconUpgradeable{IBeaconUpgradeableCaller: IBeaconUpgradeableCaller{contract: contract}, IBeaconUpgradeableTransactor: IBeaconUpgradeableTransactor{contract: contract}, IBeaconUpgradeableFilterer: IBeaconUpgradeableFilterer{contract: contract}}, nil
}

// NewIBeaconUpgradeableCaller creates a new read-only instance of IBeaconUpgradeable, bound to a specific deployed contract.
func NewIBeaconUpgradeableCaller(address common.Address, caller bind.ContractCaller) (*IBeaconUpgradeableCaller, error) {
	contract, err := bindIBeaconUpgradeable(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IBeaconUpgradeableCaller{contract: contract}, nil
}

// NewIBeaconUpgradeableTransactor creates a new write-only instance of IBeaconUpgradeable, bound to a specific deployed contract.
func NewIBeaconUpgradeableTransactor(address common.Address, transactor bind.ContractTransactor) (*IBeaconUpgradeableTransactor, error) {
	contract, err := bindIBeaconUpgradeable(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IBeaconUpgradeableTransactor{contract: contract}, nil
}

// NewIBeaconUpgradeableFilterer creates a new log filterer instance of IBeaconUpgradeable, bound to a specific deployed contract.
func NewIBeaconUpgradeableFilterer(address common.Address, filterer bind.ContractFilterer) (*IBeaconUpgradeableFilterer, error) {
	contract, err := bindIBeaconUpgradeable(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IBeaconUpgradeableFilterer{contract: contract}, nil
}

// bindIBeaconUpgradeable binds a generic wrapper to an already deployed contract.
func bindIBeaconUpgradeable(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := IBeaconUpgradeableMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IBeaconUpgradeable *IBeaconUpgradeableRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IBeaconUpgradeable.Contract.IBeaconUpgradeableCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IBeaconUpgradeable *IBeaconUpgradeableRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IBeaconUpgradeable.Contract.IBeaconUpgradeableTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IBeaconUpgradeable *IBeaconUpgradeableRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IBeaconUpgradeable.Contract.IBeaconUpgradeableTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IBeaconUpgradeable *IBeaconUpgradeableCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IBeaconUpgradeable.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IBeaconUpgradeable *IBeaconUpgradeableTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IBeaconUpgradeable.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IBeaconUpgradeable *IBeaconUpgradeableTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IBeaconUpgradeable.Contract.contract.Transact(opts, method, params...)
}

// Implementation is a free data retrieval call binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() view returns(address)
func (_IBeaconUpgradeable *IBeaconUpgradeableCaller) Implementation(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IBeaconUpgradeable.contract.Call(opts, &out, "implementation")
	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err
}

// Implementation is a free data retrieval call binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() view returns(address)
func (_IBeaconUpgradeable *IBeaconUpgradeableSession) Implementation() (common.Address, error) {
	return _IBeaconUpgradeable.Contract.Implementation(&_IBeaconUpgradeable.CallOpts)
}

// Implementation is a free data retrieval call binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() view returns(address)
func (_IBeaconUpgradeable *IBeaconUpgradeableCallerSession) Implementation() (common.Address, error) {
	return _IBeaconUpgradeable.Contract.Implementation(&_IBeaconUpgradeable.CallOpts)
}

// IERC1822ProxiableMetaData contains all meta data concerning the IERC1822Proxiable contract.
var IERC1822ProxiableMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"proxiableUUID\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"52d1902d": "proxiableUUID()",
	},
}

// IERC1822ProxiableABI is the input ABI used to generate the binding from.
// Deprecated: Use IERC1822ProxiableMetaData.ABI instead.
var IERC1822ProxiableABI = IERC1822ProxiableMetaData.ABI

// IERC1822ProxiableBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const IERC1822ProxiableBinRuntime = ``

// IERC1822ProxiableFuncSigs maps the 4-byte function signature to its string representation.
// Deprecated: Use IERC1822ProxiableMetaData.Sigs instead.
var IERC1822ProxiableFuncSigs = IERC1822ProxiableMetaData.Sigs

// IERC1822Proxiable is an auto generated Go binding around a Klaytn contract.
type IERC1822Proxiable struct {
	IERC1822ProxiableCaller     // Read-only binding to the contract
	IERC1822ProxiableTransactor // Write-only binding to the contract
	IERC1822ProxiableFilterer   // Log filterer for contract events
}

// IERC1822ProxiableCaller is an auto generated read-only Go binding around a Klaytn contract.
type IERC1822ProxiableCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IERC1822ProxiableTransactor is an auto generated write-only Go binding around a Klaytn contract.
type IERC1822ProxiableTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IERC1822ProxiableFilterer is an auto generated log filtering Go binding around a Klaytn contract events.
type IERC1822ProxiableFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IERC1822ProxiableSession is an auto generated Go binding around a Klaytn contract,
// with pre-set call and transact options.
type IERC1822ProxiableSession struct {
	Contract     *IERC1822Proxiable // Generic contract binding to set the session for
	CallOpts     bind.CallOpts      // Call options to use throughout this session
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// IERC1822ProxiableCallerSession is an auto generated read-only Go binding around a Klaytn contract,
// with pre-set call options.
type IERC1822ProxiableCallerSession struct {
	Contract *IERC1822ProxiableCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts            // Call options to use throughout this session
}

// IERC1822ProxiableTransactorSession is an auto generated write-only Go binding around a Klaytn contract,
// with pre-set transact options.
type IERC1822ProxiableTransactorSession struct {
	Contract     *IERC1822ProxiableTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// IERC1822ProxiableRaw is an auto generated low-level Go binding around a Klaytn contract.
type IERC1822ProxiableRaw struct {
	Contract *IERC1822Proxiable // Generic contract binding to access the raw methods on
}

// IERC1822ProxiableCallerRaw is an auto generated low-level read-only Go binding around a Klaytn contract.
type IERC1822ProxiableCallerRaw struct {
	Contract *IERC1822ProxiableCaller // Generic read-only contract binding to access the raw methods on
}

// IERC1822ProxiableTransactorRaw is an auto generated low-level write-only Go binding around a Klaytn contract.
type IERC1822ProxiableTransactorRaw struct {
	Contract *IERC1822ProxiableTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIERC1822Proxiable creates a new instance of IERC1822Proxiable, bound to a specific deployed contract.
func NewIERC1822Proxiable(address common.Address, backend bind.ContractBackend) (*IERC1822Proxiable, error) {
	contract, err := bindIERC1822Proxiable(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IERC1822Proxiable{IERC1822ProxiableCaller: IERC1822ProxiableCaller{contract: contract}, IERC1822ProxiableTransactor: IERC1822ProxiableTransactor{contract: contract}, IERC1822ProxiableFilterer: IERC1822ProxiableFilterer{contract: contract}}, nil
}

// NewIERC1822ProxiableCaller creates a new read-only instance of IERC1822Proxiable, bound to a specific deployed contract.
func NewIERC1822ProxiableCaller(address common.Address, caller bind.ContractCaller) (*IERC1822ProxiableCaller, error) {
	contract, err := bindIERC1822Proxiable(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IERC1822ProxiableCaller{contract: contract}, nil
}

// NewIERC1822ProxiableTransactor creates a new write-only instance of IERC1822Proxiable, bound to a specific deployed contract.
func NewIERC1822ProxiableTransactor(address common.Address, transactor bind.ContractTransactor) (*IERC1822ProxiableTransactor, error) {
	contract, err := bindIERC1822Proxiable(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IERC1822ProxiableTransactor{contract: contract}, nil
}

// NewIERC1822ProxiableFilterer creates a new log filterer instance of IERC1822Proxiable, bound to a specific deployed contract.
func NewIERC1822ProxiableFilterer(address common.Address, filterer bind.ContractFilterer) (*IERC1822ProxiableFilterer, error) {
	contract, err := bindIERC1822Proxiable(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IERC1822ProxiableFilterer{contract: contract}, nil
}

// bindIERC1822Proxiable binds a generic wrapper to an already deployed contract.
func bindIERC1822Proxiable(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := IERC1822ProxiableMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IERC1822Proxiable *IERC1822ProxiableRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IERC1822Proxiable.Contract.IERC1822ProxiableCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IERC1822Proxiable *IERC1822ProxiableRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IERC1822Proxiable.Contract.IERC1822ProxiableTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IERC1822Proxiable *IERC1822ProxiableRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IERC1822Proxiable.Contract.IERC1822ProxiableTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IERC1822Proxiable *IERC1822ProxiableCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IERC1822Proxiable.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IERC1822Proxiable *IERC1822ProxiableTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IERC1822Proxiable.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IERC1822Proxiable *IERC1822ProxiableTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IERC1822Proxiable.Contract.contract.Transact(opts, method, params...)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_IERC1822Proxiable *IERC1822ProxiableCaller) ProxiableUUID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _IERC1822Proxiable.contract.Call(opts, &out, "proxiableUUID")
	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_IERC1822Proxiable *IERC1822ProxiableSession) ProxiableUUID() ([32]byte, error) {
	return _IERC1822Proxiable.Contract.ProxiableUUID(&_IERC1822Proxiable.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_IERC1822Proxiable *IERC1822ProxiableCallerSession) ProxiableUUID() ([32]byte, error) {
	return _IERC1822Proxiable.Contract.ProxiableUUID(&_IERC1822Proxiable.CallOpts)
}

// IERC1822ProxiableUpgradeableMetaData contains all meta data concerning the IERC1822ProxiableUpgradeable contract.
var IERC1822ProxiableUpgradeableMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"proxiableUUID\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"52d1902d": "proxiableUUID()",
	},
}

// IERC1822ProxiableUpgradeableABI is the input ABI used to generate the binding from.
// Deprecated: Use IERC1822ProxiableUpgradeableMetaData.ABI instead.
var IERC1822ProxiableUpgradeableABI = IERC1822ProxiableUpgradeableMetaData.ABI

// IERC1822ProxiableUpgradeableBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const IERC1822ProxiableUpgradeableBinRuntime = ``

// IERC1822ProxiableUpgradeableFuncSigs maps the 4-byte function signature to its string representation.
// Deprecated: Use IERC1822ProxiableUpgradeableMetaData.Sigs instead.
var IERC1822ProxiableUpgradeableFuncSigs = IERC1822ProxiableUpgradeableMetaData.Sigs

// IERC1822ProxiableUpgradeable is an auto generated Go binding around a Klaytn contract.
type IERC1822ProxiableUpgradeable struct {
	IERC1822ProxiableUpgradeableCaller     // Read-only binding to the contract
	IERC1822ProxiableUpgradeableTransactor // Write-only binding to the contract
	IERC1822ProxiableUpgradeableFilterer   // Log filterer for contract events
}

// IERC1822ProxiableUpgradeableCaller is an auto generated read-only Go binding around a Klaytn contract.
type IERC1822ProxiableUpgradeableCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IERC1822ProxiableUpgradeableTransactor is an auto generated write-only Go binding around a Klaytn contract.
type IERC1822ProxiableUpgradeableTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IERC1822ProxiableUpgradeableFilterer is an auto generated log filtering Go binding around a Klaytn contract events.
type IERC1822ProxiableUpgradeableFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IERC1822ProxiableUpgradeableSession is an auto generated Go binding around a Klaytn contract,
// with pre-set call and transact options.
type IERC1822ProxiableUpgradeableSession struct {
	Contract     *IERC1822ProxiableUpgradeable // Generic contract binding to set the session for
	CallOpts     bind.CallOpts                 // Call options to use throughout this session
	TransactOpts bind.TransactOpts             // Transaction auth options to use throughout this session
}

// IERC1822ProxiableUpgradeableCallerSession is an auto generated read-only Go binding around a Klaytn contract,
// with pre-set call options.
type IERC1822ProxiableUpgradeableCallerSession struct {
	Contract *IERC1822ProxiableUpgradeableCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                       // Call options to use throughout this session
}

// IERC1822ProxiableUpgradeableTransactorSession is an auto generated write-only Go binding around a Klaytn contract,
// with pre-set transact options.
type IERC1822ProxiableUpgradeableTransactorSession struct {
	Contract     *IERC1822ProxiableUpgradeableTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                       // Transaction auth options to use throughout this session
}

// IERC1822ProxiableUpgradeableRaw is an auto generated low-level Go binding around a Klaytn contract.
type IERC1822ProxiableUpgradeableRaw struct {
	Contract *IERC1822ProxiableUpgradeable // Generic contract binding to access the raw methods on
}

// IERC1822ProxiableUpgradeableCallerRaw is an auto generated low-level read-only Go binding around a Klaytn contract.
type IERC1822ProxiableUpgradeableCallerRaw struct {
	Contract *IERC1822ProxiableUpgradeableCaller // Generic read-only contract binding to access the raw methods on
}

// IERC1822ProxiableUpgradeableTransactorRaw is an auto generated low-level write-only Go binding around a Klaytn contract.
type IERC1822ProxiableUpgradeableTransactorRaw struct {
	Contract *IERC1822ProxiableUpgradeableTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIERC1822ProxiableUpgradeable creates a new instance of IERC1822ProxiableUpgradeable, bound to a specific deployed contract.
func NewIERC1822ProxiableUpgradeable(address common.Address, backend bind.ContractBackend) (*IERC1822ProxiableUpgradeable, error) {
	contract, err := bindIERC1822ProxiableUpgradeable(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IERC1822ProxiableUpgradeable{IERC1822ProxiableUpgradeableCaller: IERC1822ProxiableUpgradeableCaller{contract: contract}, IERC1822ProxiableUpgradeableTransactor: IERC1822ProxiableUpgradeableTransactor{contract: contract}, IERC1822ProxiableUpgradeableFilterer: IERC1822ProxiableUpgradeableFilterer{contract: contract}}, nil
}

// NewIERC1822ProxiableUpgradeableCaller creates a new read-only instance of IERC1822ProxiableUpgradeable, bound to a specific deployed contract.
func NewIERC1822ProxiableUpgradeableCaller(address common.Address, caller bind.ContractCaller) (*IERC1822ProxiableUpgradeableCaller, error) {
	contract, err := bindIERC1822ProxiableUpgradeable(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IERC1822ProxiableUpgradeableCaller{contract: contract}, nil
}

// NewIERC1822ProxiableUpgradeableTransactor creates a new write-only instance of IERC1822ProxiableUpgradeable, bound to a specific deployed contract.
func NewIERC1822ProxiableUpgradeableTransactor(address common.Address, transactor bind.ContractTransactor) (*IERC1822ProxiableUpgradeableTransactor, error) {
	contract, err := bindIERC1822ProxiableUpgradeable(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IERC1822ProxiableUpgradeableTransactor{contract: contract}, nil
}

// NewIERC1822ProxiableUpgradeableFilterer creates a new log filterer instance of IERC1822ProxiableUpgradeable, bound to a specific deployed contract.
func NewIERC1822ProxiableUpgradeableFilterer(address common.Address, filterer bind.ContractFilterer) (*IERC1822ProxiableUpgradeableFilterer, error) {
	contract, err := bindIERC1822ProxiableUpgradeable(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IERC1822ProxiableUpgradeableFilterer{contract: contract}, nil
}

// bindIERC1822ProxiableUpgradeable binds a generic wrapper to an already deployed contract.
func bindIERC1822ProxiableUpgradeable(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := IERC1822ProxiableUpgradeableMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IERC1822ProxiableUpgradeable *IERC1822ProxiableUpgradeableRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IERC1822ProxiableUpgradeable.Contract.IERC1822ProxiableUpgradeableCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IERC1822ProxiableUpgradeable *IERC1822ProxiableUpgradeableRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IERC1822ProxiableUpgradeable.Contract.IERC1822ProxiableUpgradeableTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IERC1822ProxiableUpgradeable *IERC1822ProxiableUpgradeableRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IERC1822ProxiableUpgradeable.Contract.IERC1822ProxiableUpgradeableTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IERC1822ProxiableUpgradeable *IERC1822ProxiableUpgradeableCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IERC1822ProxiableUpgradeable.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IERC1822ProxiableUpgradeable *IERC1822ProxiableUpgradeableTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IERC1822ProxiableUpgradeable.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IERC1822ProxiableUpgradeable *IERC1822ProxiableUpgradeableTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IERC1822ProxiableUpgradeable.Contract.contract.Transact(opts, method, params...)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_IERC1822ProxiableUpgradeable *IERC1822ProxiableUpgradeableCaller) ProxiableUUID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _IERC1822ProxiableUpgradeable.contract.Call(opts, &out, "proxiableUUID")
	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_IERC1822ProxiableUpgradeable *IERC1822ProxiableUpgradeableSession) ProxiableUUID() ([32]byte, error) {
	return _IERC1822ProxiableUpgradeable.Contract.ProxiableUUID(&_IERC1822ProxiableUpgradeable.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_IERC1822ProxiableUpgradeable *IERC1822ProxiableUpgradeableCallerSession) ProxiableUUID() ([32]byte, error) {
	return _IERC1822ProxiableUpgradeable.Contract.ProxiableUUID(&_IERC1822ProxiableUpgradeable.CallOpts)
}

// IERC1967MetaData contains all meta data concerning the IERC1967 contract.
var IERC1967MetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"previousAdmin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"AdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"beacon\",\"type\":\"address\"}],\"name\":\"BeaconUpgraded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"}]",
}

// IERC1967ABI is the input ABI used to generate the binding from.
// Deprecated: Use IERC1967MetaData.ABI instead.
var IERC1967ABI = IERC1967MetaData.ABI

// IERC1967BinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const IERC1967BinRuntime = ``

// IERC1967 is an auto generated Go binding around a Klaytn contract.
type IERC1967 struct {
	IERC1967Caller     // Read-only binding to the contract
	IERC1967Transactor // Write-only binding to the contract
	IERC1967Filterer   // Log filterer for contract events
}

// IERC1967Caller is an auto generated read-only Go binding around a Klaytn contract.
type IERC1967Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IERC1967Transactor is an auto generated write-only Go binding around a Klaytn contract.
type IERC1967Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IERC1967Filterer is an auto generated log filtering Go binding around a Klaytn contract events.
type IERC1967Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IERC1967Session is an auto generated Go binding around a Klaytn contract,
// with pre-set call and transact options.
type IERC1967Session struct {
	Contract     *IERC1967         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IERC1967CallerSession is an auto generated read-only Go binding around a Klaytn contract,
// with pre-set call options.
type IERC1967CallerSession struct {
	Contract *IERC1967Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// IERC1967TransactorSession is an auto generated write-only Go binding around a Klaytn contract,
// with pre-set transact options.
type IERC1967TransactorSession struct {
	Contract     *IERC1967Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// IERC1967Raw is an auto generated low-level Go binding around a Klaytn contract.
type IERC1967Raw struct {
	Contract *IERC1967 // Generic contract binding to access the raw methods on
}

// IERC1967CallerRaw is an auto generated low-level read-only Go binding around a Klaytn contract.
type IERC1967CallerRaw struct {
	Contract *IERC1967Caller // Generic read-only contract binding to access the raw methods on
}

// IERC1967TransactorRaw is an auto generated low-level write-only Go binding around a Klaytn contract.
type IERC1967TransactorRaw struct {
	Contract *IERC1967Transactor // Generic write-only contract binding to access the raw methods on
}

// NewIERC1967 creates a new instance of IERC1967, bound to a specific deployed contract.
func NewIERC1967(address common.Address, backend bind.ContractBackend) (*IERC1967, error) {
	contract, err := bindIERC1967(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IERC1967{IERC1967Caller: IERC1967Caller{contract: contract}, IERC1967Transactor: IERC1967Transactor{contract: contract}, IERC1967Filterer: IERC1967Filterer{contract: contract}}, nil
}

// NewIERC1967Caller creates a new read-only instance of IERC1967, bound to a specific deployed contract.
func NewIERC1967Caller(address common.Address, caller bind.ContractCaller) (*IERC1967Caller, error) {
	contract, err := bindIERC1967(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IERC1967Caller{contract: contract}, nil
}

// NewIERC1967Transactor creates a new write-only instance of IERC1967, bound to a specific deployed contract.
func NewIERC1967Transactor(address common.Address, transactor bind.ContractTransactor) (*IERC1967Transactor, error) {
	contract, err := bindIERC1967(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IERC1967Transactor{contract: contract}, nil
}

// NewIERC1967Filterer creates a new log filterer instance of IERC1967, bound to a specific deployed contract.
func NewIERC1967Filterer(address common.Address, filterer bind.ContractFilterer) (*IERC1967Filterer, error) {
	contract, err := bindIERC1967(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IERC1967Filterer{contract: contract}, nil
}

// bindIERC1967 binds a generic wrapper to an already deployed contract.
func bindIERC1967(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := IERC1967MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IERC1967 *IERC1967Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IERC1967.Contract.IERC1967Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IERC1967 *IERC1967Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IERC1967.Contract.IERC1967Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IERC1967 *IERC1967Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IERC1967.Contract.IERC1967Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IERC1967 *IERC1967CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IERC1967.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IERC1967 *IERC1967TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IERC1967.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IERC1967 *IERC1967TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IERC1967.Contract.contract.Transact(opts, method, params...)
}

// IERC1967AdminChangedIterator is returned from FilterAdminChanged and is used to iterate over the raw logs and unpacked data for AdminChanged events raised by the IERC1967 contract.
type IERC1967AdminChangedIterator struct {
	Event *IERC1967AdminChanged // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log      // Log channel receiving the found contract events
	sub  klaytn.Subscription // Subscription for errors, completion and termination
	done bool                // Whether the subscription completed delivering logs
	fail error               // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *IERC1967AdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IERC1967AdminChanged)
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
		it.Event = new(IERC1967AdminChanged)
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
func (it *IERC1967AdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IERC1967AdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IERC1967AdminChanged represents a AdminChanged event raised by the IERC1967 contract.
type IERC1967AdminChanged struct {
	PreviousAdmin common.Address
	NewAdmin      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterAdminChanged is a free log retrieval operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_IERC1967 *IERC1967Filterer) FilterAdminChanged(opts *bind.FilterOpts) (*IERC1967AdminChangedIterator, error) {
	logs, sub, err := _IERC1967.contract.FilterLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return &IERC1967AdminChangedIterator{contract: _IERC1967.contract, event: "AdminChanged", logs: logs, sub: sub}, nil
}

// WatchAdminChanged is a free log subscription operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_IERC1967 *IERC1967Filterer) WatchAdminChanged(opts *bind.WatchOpts, sink chan<- *IERC1967AdminChanged) (event.Subscription, error) {
	logs, sub, err := _IERC1967.contract.WatchLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IERC1967AdminChanged)
				if err := _IERC1967.contract.UnpackLog(event, "AdminChanged", log); err != nil {
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

// ParseAdminChanged is a log parse operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_IERC1967 *IERC1967Filterer) ParseAdminChanged(log types.Log) (*IERC1967AdminChanged, error) {
	event := new(IERC1967AdminChanged)
	if err := _IERC1967.contract.UnpackLog(event, "AdminChanged", log); err != nil {
		return nil, err
	}
	return event, nil
}

// IERC1967BeaconUpgradedIterator is returned from FilterBeaconUpgraded and is used to iterate over the raw logs and unpacked data for BeaconUpgraded events raised by the IERC1967 contract.
type IERC1967BeaconUpgradedIterator struct {
	Event *IERC1967BeaconUpgraded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log      // Log channel receiving the found contract events
	sub  klaytn.Subscription // Subscription for errors, completion and termination
	done bool                // Whether the subscription completed delivering logs
	fail error               // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *IERC1967BeaconUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IERC1967BeaconUpgraded)
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
		it.Event = new(IERC1967BeaconUpgraded)
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
func (it *IERC1967BeaconUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IERC1967BeaconUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IERC1967BeaconUpgraded represents a BeaconUpgraded event raised by the IERC1967 contract.
type IERC1967BeaconUpgraded struct {
	Beacon common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterBeaconUpgraded is a free log retrieval operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_IERC1967 *IERC1967Filterer) FilterBeaconUpgraded(opts *bind.FilterOpts, beacon []common.Address) (*IERC1967BeaconUpgradedIterator, error) {
	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _IERC1967.contract.FilterLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return &IERC1967BeaconUpgradedIterator{contract: _IERC1967.contract, event: "BeaconUpgraded", logs: logs, sub: sub}, nil
}

// WatchBeaconUpgraded is a free log subscription operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_IERC1967 *IERC1967Filterer) WatchBeaconUpgraded(opts *bind.WatchOpts, sink chan<- *IERC1967BeaconUpgraded, beacon []common.Address) (event.Subscription, error) {
	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _IERC1967.contract.WatchLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IERC1967BeaconUpgraded)
				if err := _IERC1967.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
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

// ParseBeaconUpgraded is a log parse operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_IERC1967 *IERC1967Filterer) ParseBeaconUpgraded(log types.Log) (*IERC1967BeaconUpgraded, error) {
	event := new(IERC1967BeaconUpgraded)
	if err := _IERC1967.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
		return nil, err
	}
	return event, nil
}

// IERC1967UpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the IERC1967 contract.
type IERC1967UpgradedIterator struct {
	Event *IERC1967Upgraded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log      // Log channel receiving the found contract events
	sub  klaytn.Subscription // Subscription for errors, completion and termination
	done bool                // Whether the subscription completed delivering logs
	fail error               // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *IERC1967UpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IERC1967Upgraded)
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
		it.Event = new(IERC1967Upgraded)
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
func (it *IERC1967UpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IERC1967UpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IERC1967Upgraded represents a Upgraded event raised by the IERC1967 contract.
type IERC1967Upgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_IERC1967 *IERC1967Filterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*IERC1967UpgradedIterator, error) {
	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _IERC1967.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &IERC1967UpgradedIterator{contract: _IERC1967.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_IERC1967 *IERC1967Filterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *IERC1967Upgraded, implementation []common.Address) (event.Subscription, error) {
	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _IERC1967.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IERC1967Upgraded)
				if err := _IERC1967.contract.UnpackLog(event, "Upgraded", log); err != nil {
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

// ParseUpgraded is a log parse operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_IERC1967 *IERC1967Filterer) ParseUpgraded(log types.Log) (*IERC1967Upgraded, error) {
	event := new(IERC1967Upgraded)
	if err := _IERC1967.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	return event, nil
}

// IERC1967UpgradeableMetaData contains all meta data concerning the IERC1967Upgradeable contract.
var IERC1967UpgradeableMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"previousAdmin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"AdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"beacon\",\"type\":\"address\"}],\"name\":\"BeaconUpgraded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"}]",
}

// IERC1967UpgradeableABI is the input ABI used to generate the binding from.
// Deprecated: Use IERC1967UpgradeableMetaData.ABI instead.
var IERC1967UpgradeableABI = IERC1967UpgradeableMetaData.ABI

// IERC1967UpgradeableBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const IERC1967UpgradeableBinRuntime = ``

// IERC1967Upgradeable is an auto generated Go binding around a Klaytn contract.
type IERC1967Upgradeable struct {
	IERC1967UpgradeableCaller     // Read-only binding to the contract
	IERC1967UpgradeableTransactor // Write-only binding to the contract
	IERC1967UpgradeableFilterer   // Log filterer for contract events
}

// IERC1967UpgradeableCaller is an auto generated read-only Go binding around a Klaytn contract.
type IERC1967UpgradeableCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IERC1967UpgradeableTransactor is an auto generated write-only Go binding around a Klaytn contract.
type IERC1967UpgradeableTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IERC1967UpgradeableFilterer is an auto generated log filtering Go binding around a Klaytn contract events.
type IERC1967UpgradeableFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IERC1967UpgradeableSession is an auto generated Go binding around a Klaytn contract,
// with pre-set call and transact options.
type IERC1967UpgradeableSession struct {
	Contract     *IERC1967Upgradeable // Generic contract binding to set the session for
	CallOpts     bind.CallOpts        // Call options to use throughout this session
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// IERC1967UpgradeableCallerSession is an auto generated read-only Go binding around a Klaytn contract,
// with pre-set call options.
type IERC1967UpgradeableCallerSession struct {
	Contract *IERC1967UpgradeableCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts              // Call options to use throughout this session
}

// IERC1967UpgradeableTransactorSession is an auto generated write-only Go binding around a Klaytn contract,
// with pre-set transact options.
type IERC1967UpgradeableTransactorSession struct {
	Contract     *IERC1967UpgradeableTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts              // Transaction auth options to use throughout this session
}

// IERC1967UpgradeableRaw is an auto generated low-level Go binding around a Klaytn contract.
type IERC1967UpgradeableRaw struct {
	Contract *IERC1967Upgradeable // Generic contract binding to access the raw methods on
}

// IERC1967UpgradeableCallerRaw is an auto generated low-level read-only Go binding around a Klaytn contract.
type IERC1967UpgradeableCallerRaw struct {
	Contract *IERC1967UpgradeableCaller // Generic read-only contract binding to access the raw methods on
}

// IERC1967UpgradeableTransactorRaw is an auto generated low-level write-only Go binding around a Klaytn contract.
type IERC1967UpgradeableTransactorRaw struct {
	Contract *IERC1967UpgradeableTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIERC1967Upgradeable creates a new instance of IERC1967Upgradeable, bound to a specific deployed contract.
func NewIERC1967Upgradeable(address common.Address, backend bind.ContractBackend) (*IERC1967Upgradeable, error) {
	contract, err := bindIERC1967Upgradeable(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IERC1967Upgradeable{IERC1967UpgradeableCaller: IERC1967UpgradeableCaller{contract: contract}, IERC1967UpgradeableTransactor: IERC1967UpgradeableTransactor{contract: contract}, IERC1967UpgradeableFilterer: IERC1967UpgradeableFilterer{contract: contract}}, nil
}

// NewIERC1967UpgradeableCaller creates a new read-only instance of IERC1967Upgradeable, bound to a specific deployed contract.
func NewIERC1967UpgradeableCaller(address common.Address, caller bind.ContractCaller) (*IERC1967UpgradeableCaller, error) {
	contract, err := bindIERC1967Upgradeable(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IERC1967UpgradeableCaller{contract: contract}, nil
}

// NewIERC1967UpgradeableTransactor creates a new write-only instance of IERC1967Upgradeable, bound to a specific deployed contract.
func NewIERC1967UpgradeableTransactor(address common.Address, transactor bind.ContractTransactor) (*IERC1967UpgradeableTransactor, error) {
	contract, err := bindIERC1967Upgradeable(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IERC1967UpgradeableTransactor{contract: contract}, nil
}

// NewIERC1967UpgradeableFilterer creates a new log filterer instance of IERC1967Upgradeable, bound to a specific deployed contract.
func NewIERC1967UpgradeableFilterer(address common.Address, filterer bind.ContractFilterer) (*IERC1967UpgradeableFilterer, error) {
	contract, err := bindIERC1967Upgradeable(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IERC1967UpgradeableFilterer{contract: contract}, nil
}

// bindIERC1967Upgradeable binds a generic wrapper to an already deployed contract.
func bindIERC1967Upgradeable(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := IERC1967UpgradeableMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IERC1967Upgradeable *IERC1967UpgradeableRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IERC1967Upgradeable.Contract.IERC1967UpgradeableCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IERC1967Upgradeable *IERC1967UpgradeableRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IERC1967Upgradeable.Contract.IERC1967UpgradeableTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IERC1967Upgradeable *IERC1967UpgradeableRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IERC1967Upgradeable.Contract.IERC1967UpgradeableTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IERC1967Upgradeable *IERC1967UpgradeableCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IERC1967Upgradeable.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IERC1967Upgradeable *IERC1967UpgradeableTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IERC1967Upgradeable.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IERC1967Upgradeable *IERC1967UpgradeableTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IERC1967Upgradeable.Contract.contract.Transact(opts, method, params...)
}

// IERC1967UpgradeableAdminChangedIterator is returned from FilterAdminChanged and is used to iterate over the raw logs and unpacked data for AdminChanged events raised by the IERC1967Upgradeable contract.
type IERC1967UpgradeableAdminChangedIterator struct {
	Event *IERC1967UpgradeableAdminChanged // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log      // Log channel receiving the found contract events
	sub  klaytn.Subscription // Subscription for errors, completion and termination
	done bool                // Whether the subscription completed delivering logs
	fail error               // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *IERC1967UpgradeableAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IERC1967UpgradeableAdminChanged)
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
		it.Event = new(IERC1967UpgradeableAdminChanged)
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
func (it *IERC1967UpgradeableAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IERC1967UpgradeableAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IERC1967UpgradeableAdminChanged represents a AdminChanged event raised by the IERC1967Upgradeable contract.
type IERC1967UpgradeableAdminChanged struct {
	PreviousAdmin common.Address
	NewAdmin      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterAdminChanged is a free log retrieval operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_IERC1967Upgradeable *IERC1967UpgradeableFilterer) FilterAdminChanged(opts *bind.FilterOpts) (*IERC1967UpgradeableAdminChangedIterator, error) {
	logs, sub, err := _IERC1967Upgradeable.contract.FilterLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return &IERC1967UpgradeableAdminChangedIterator{contract: _IERC1967Upgradeable.contract, event: "AdminChanged", logs: logs, sub: sub}, nil
}

// WatchAdminChanged is a free log subscription operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_IERC1967Upgradeable *IERC1967UpgradeableFilterer) WatchAdminChanged(opts *bind.WatchOpts, sink chan<- *IERC1967UpgradeableAdminChanged) (event.Subscription, error) {
	logs, sub, err := _IERC1967Upgradeable.contract.WatchLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IERC1967UpgradeableAdminChanged)
				if err := _IERC1967Upgradeable.contract.UnpackLog(event, "AdminChanged", log); err != nil {
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

// ParseAdminChanged is a log parse operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_IERC1967Upgradeable *IERC1967UpgradeableFilterer) ParseAdminChanged(log types.Log) (*IERC1967UpgradeableAdminChanged, error) {
	event := new(IERC1967UpgradeableAdminChanged)
	if err := _IERC1967Upgradeable.contract.UnpackLog(event, "AdminChanged", log); err != nil {
		return nil, err
	}
	return event, nil
}

// IERC1967UpgradeableBeaconUpgradedIterator is returned from FilterBeaconUpgraded and is used to iterate over the raw logs and unpacked data for BeaconUpgraded events raised by the IERC1967Upgradeable contract.
type IERC1967UpgradeableBeaconUpgradedIterator struct {
	Event *IERC1967UpgradeableBeaconUpgraded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log      // Log channel receiving the found contract events
	sub  klaytn.Subscription // Subscription for errors, completion and termination
	done bool                // Whether the subscription completed delivering logs
	fail error               // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *IERC1967UpgradeableBeaconUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IERC1967UpgradeableBeaconUpgraded)
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
		it.Event = new(IERC1967UpgradeableBeaconUpgraded)
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
func (it *IERC1967UpgradeableBeaconUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IERC1967UpgradeableBeaconUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IERC1967UpgradeableBeaconUpgraded represents a BeaconUpgraded event raised by the IERC1967Upgradeable contract.
type IERC1967UpgradeableBeaconUpgraded struct {
	Beacon common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterBeaconUpgraded is a free log retrieval operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_IERC1967Upgradeable *IERC1967UpgradeableFilterer) FilterBeaconUpgraded(opts *bind.FilterOpts, beacon []common.Address) (*IERC1967UpgradeableBeaconUpgradedIterator, error) {
	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _IERC1967Upgradeable.contract.FilterLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return &IERC1967UpgradeableBeaconUpgradedIterator{contract: _IERC1967Upgradeable.contract, event: "BeaconUpgraded", logs: logs, sub: sub}, nil
}

// WatchBeaconUpgraded is a free log subscription operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_IERC1967Upgradeable *IERC1967UpgradeableFilterer) WatchBeaconUpgraded(opts *bind.WatchOpts, sink chan<- *IERC1967UpgradeableBeaconUpgraded, beacon []common.Address) (event.Subscription, error) {
	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _IERC1967Upgradeable.contract.WatchLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IERC1967UpgradeableBeaconUpgraded)
				if err := _IERC1967Upgradeable.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
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

// ParseBeaconUpgraded is a log parse operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_IERC1967Upgradeable *IERC1967UpgradeableFilterer) ParseBeaconUpgraded(log types.Log) (*IERC1967UpgradeableBeaconUpgraded, error) {
	event := new(IERC1967UpgradeableBeaconUpgraded)
	if err := _IERC1967Upgradeable.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
		return nil, err
	}
	return event, nil
}

// IERC1967UpgradeableUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the IERC1967Upgradeable contract.
type IERC1967UpgradeableUpgradedIterator struct {
	Event *IERC1967UpgradeableUpgraded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log      // Log channel receiving the found contract events
	sub  klaytn.Subscription // Subscription for errors, completion and termination
	done bool                // Whether the subscription completed delivering logs
	fail error               // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *IERC1967UpgradeableUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IERC1967UpgradeableUpgraded)
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
		it.Event = new(IERC1967UpgradeableUpgraded)
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
func (it *IERC1967UpgradeableUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IERC1967UpgradeableUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IERC1967UpgradeableUpgraded represents a Upgraded event raised by the IERC1967Upgradeable contract.
type IERC1967UpgradeableUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_IERC1967Upgradeable *IERC1967UpgradeableFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*IERC1967UpgradeableUpgradedIterator, error) {
	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _IERC1967Upgradeable.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &IERC1967UpgradeableUpgradedIterator{contract: _IERC1967Upgradeable.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_IERC1967Upgradeable *IERC1967UpgradeableFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *IERC1967UpgradeableUpgraded, implementation []common.Address) (event.Subscription, error) {
	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _IERC1967Upgradeable.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IERC1967UpgradeableUpgraded)
				if err := _IERC1967Upgradeable.contract.UnpackLog(event, "Upgraded", log); err != nil {
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

// ParseUpgraded is a log parse operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_IERC1967Upgradeable *IERC1967UpgradeableFilterer) ParseUpgraded(log types.Log) (*IERC1967UpgradeableUpgraded, error) {
	event := new(IERC1967UpgradeableUpgraded)
	if err := _IERC1967Upgradeable.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	return event, nil
}

// IKIP113MetaData contains all meta data concerning the IKIP113 contract.
var IKIP113MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"getAllBlsInfo\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"nodeIdList\",\"type\":\"address[]\"},{\"components\":[{\"internalType\":\"bytes\",\"name\":\"publicKey\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"pop\",\"type\":\"bytes\"}],\"internalType\":\"structIKIP113.BlsPublicKeyInfo[]\",\"name\":\"pubkeyList\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"6968b53f": "getAllBlsInfo()",
	},
}

// IKIP113ABI is the input ABI used to generate the binding from.
// Deprecated: Use IKIP113MetaData.ABI instead.
var IKIP113ABI = IKIP113MetaData.ABI

// IKIP113BinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const IKIP113BinRuntime = ``

// IKIP113FuncSigs maps the 4-byte function signature to its string representation.
// Deprecated: Use IKIP113MetaData.Sigs instead.
var IKIP113FuncSigs = IKIP113MetaData.Sigs

// IKIP113 is an auto generated Go binding around a Klaytn contract.
type IKIP113 struct {
	IKIP113Caller     // Read-only binding to the contract
	IKIP113Transactor // Write-only binding to the contract
	IKIP113Filterer   // Log filterer for contract events
}

// IKIP113Caller is an auto generated read-only Go binding around a Klaytn contract.
type IKIP113Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IKIP113Transactor is an auto generated write-only Go binding around a Klaytn contract.
type IKIP113Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IKIP113Filterer is an auto generated log filtering Go binding around a Klaytn contract events.
type IKIP113Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IKIP113Session is an auto generated Go binding around a Klaytn contract,
// with pre-set call and transact options.
type IKIP113Session struct {
	Contract     *IKIP113          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IKIP113CallerSession is an auto generated read-only Go binding around a Klaytn contract,
// with pre-set call options.
type IKIP113CallerSession struct {
	Contract *IKIP113Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// IKIP113TransactorSession is an auto generated write-only Go binding around a Klaytn contract,
// with pre-set transact options.
type IKIP113TransactorSession struct {
	Contract     *IKIP113Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// IKIP113Raw is an auto generated low-level Go binding around a Klaytn contract.
type IKIP113Raw struct {
	Contract *IKIP113 // Generic contract binding to access the raw methods on
}

// IKIP113CallerRaw is an auto generated low-level read-only Go binding around a Klaytn contract.
type IKIP113CallerRaw struct {
	Contract *IKIP113Caller // Generic read-only contract binding to access the raw methods on
}

// IKIP113TransactorRaw is an auto generated low-level write-only Go binding around a Klaytn contract.
type IKIP113TransactorRaw struct {
	Contract *IKIP113Transactor // Generic write-only contract binding to access the raw methods on
}

// NewIKIP113 creates a new instance of IKIP113, bound to a specific deployed contract.
func NewIKIP113(address common.Address, backend bind.ContractBackend) (*IKIP113, error) {
	contract, err := bindIKIP113(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IKIP113{IKIP113Caller: IKIP113Caller{contract: contract}, IKIP113Transactor: IKIP113Transactor{contract: contract}, IKIP113Filterer: IKIP113Filterer{contract: contract}}, nil
}

// NewIKIP113Caller creates a new read-only instance of IKIP113, bound to a specific deployed contract.
func NewIKIP113Caller(address common.Address, caller bind.ContractCaller) (*IKIP113Caller, error) {
	contract, err := bindIKIP113(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IKIP113Caller{contract: contract}, nil
}

// NewIKIP113Transactor creates a new write-only instance of IKIP113, bound to a specific deployed contract.
func NewIKIP113Transactor(address common.Address, transactor bind.ContractTransactor) (*IKIP113Transactor, error) {
	contract, err := bindIKIP113(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IKIP113Transactor{contract: contract}, nil
}

// NewIKIP113Filterer creates a new log filterer instance of IKIP113, bound to a specific deployed contract.
func NewIKIP113Filterer(address common.Address, filterer bind.ContractFilterer) (*IKIP113Filterer, error) {
	contract, err := bindIKIP113(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IKIP113Filterer{contract: contract}, nil
}

// bindIKIP113 binds a generic wrapper to an already deployed contract.
func bindIKIP113(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := IKIP113MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IKIP113 *IKIP113Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IKIP113.Contract.IKIP113Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IKIP113 *IKIP113Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IKIP113.Contract.IKIP113Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IKIP113 *IKIP113Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IKIP113.Contract.IKIP113Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IKIP113 *IKIP113CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IKIP113.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IKIP113 *IKIP113TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IKIP113.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IKIP113 *IKIP113TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IKIP113.Contract.contract.Transact(opts, method, params...)
}

// GetAllBlsInfo is a free data retrieval call binding the contract method 0x6968b53f.
//
// Solidity: function getAllBlsInfo() view returns(address[] nodeIdList, (bytes,bytes)[] pubkeyList)
func (_IKIP113 *IKIP113Caller) GetAllBlsInfo(opts *bind.CallOpts) (struct {
	NodeIdList []common.Address
	PubkeyList []IKIP113BlsPublicKeyInfo
}, error,
) {
	var out []interface{}
	err := _IKIP113.contract.Call(opts, &out, "getAllBlsInfo")

	outstruct := new(struct {
		NodeIdList []common.Address
		PubkeyList []IKIP113BlsPublicKeyInfo
	})

	outstruct.NodeIdList = *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)
	outstruct.PubkeyList = *abi.ConvertType(out[1], new([]IKIP113BlsPublicKeyInfo)).(*[]IKIP113BlsPublicKeyInfo)
	return *outstruct, err
}

// GetAllBlsInfo is a free data retrieval call binding the contract method 0x6968b53f.
//
// Solidity: function getAllBlsInfo() view returns(address[] nodeIdList, (bytes,bytes)[] pubkeyList)
func (_IKIP113 *IKIP113Session) GetAllBlsInfo() (struct {
	NodeIdList []common.Address
	PubkeyList []IKIP113BlsPublicKeyInfo
}, error,
) {
	return _IKIP113.Contract.GetAllBlsInfo(&_IKIP113.CallOpts)
}

// GetAllBlsInfo is a free data retrieval call binding the contract method 0x6968b53f.
//
// Solidity: function getAllBlsInfo() view returns(address[] nodeIdList, (bytes,bytes)[] pubkeyList)
func (_IKIP113 *IKIP113CallerSession) GetAllBlsInfo() (struct {
	NodeIdList []common.Address
	PubkeyList []IKIP113BlsPublicKeyInfo
}, error,
) {
	return _IKIP113.Contract.GetAllBlsInfo(&_IKIP113.CallOpts)
}

// IRegistryMetaData contains all meta data concerning the IRegistry contract.
var IRegistryMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"activation\",\"type\":\"uint256\"}],\"name\":\"Registered\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"name\":\"getActiveAddr\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllNames\",\"outputs\":[{\"internalType\":\"string[]\",\"name\":\"\",\"type\":\"string[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"name\":\"getAllRecords\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"activation\",\"type\":\"uint256\"}],\"internalType\":\"structIRegistry.Record[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"names\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"records\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"activation\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"activation\",\"type\":\"uint256\"}],\"name\":\"register\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"e2693e3f": "getActiveAddr(string)",
		"fb825e5f": "getAllNames()",
		"78d573a2": "getAllRecords(string)",
		"4622ab03": "names(uint256)",
		"8da5cb5b": "owner()",
		"3b51650d": "records(string,uint256)",
		"d393c871": "register(string,address,uint256)",
		"f2fde38b": "transferOwnership(address)",
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
func (_IRegistry *IRegistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
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
func (_IRegistry *IRegistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
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

// GetAllNames is a free data retrieval call binding the contract method 0xfb825e5f.
//
// Solidity: function getAllNames() view returns(string[])
func (_IRegistry *IRegistryCaller) GetAllNames(opts *bind.CallOpts) ([]string, error) {
	var out []interface{}
	err := _IRegistry.contract.Call(opts, &out, "getAllNames")
	if err != nil {
		return *new([]string), err
	}

	out0 := *abi.ConvertType(out[0], new([]string)).(*[]string)

	return out0, err
}

// GetAllNames is a free data retrieval call binding the contract method 0xfb825e5f.
//
// Solidity: function getAllNames() view returns(string[])
func (_IRegistry *IRegistrySession) GetAllNames() ([]string, error) {
	return _IRegistry.Contract.GetAllNames(&_IRegistry.CallOpts)
}

// GetAllNames is a free data retrieval call binding the contract method 0xfb825e5f.
//
// Solidity: function getAllNames() view returns(string[])
func (_IRegistry *IRegistryCallerSession) GetAllNames() ([]string, error) {
	return _IRegistry.Contract.GetAllNames(&_IRegistry.CallOpts)
}

// GetAllRecords is a free data retrieval call binding the contract method 0x78d573a2.
//
// Solidity: function getAllRecords(string name) view returns((address,uint256)[])
func (_IRegistry *IRegistryCaller) GetAllRecords(opts *bind.CallOpts, name string) ([]IRegistryRecord, error) {
	var out []interface{}
	err := _IRegistry.contract.Call(opts, &out, "getAllRecords", name)
	if err != nil {
		return *new([]IRegistryRecord), err
	}

	out0 := *abi.ConvertType(out[0], new([]IRegistryRecord)).(*[]IRegistryRecord)

	return out0, err
}

// GetAllRecords is a free data retrieval call binding the contract method 0x78d573a2.
//
// Solidity: function getAllRecords(string name) view returns((address,uint256)[])
func (_IRegistry *IRegistrySession) GetAllRecords(name string) ([]IRegistryRecord, error) {
	return _IRegistry.Contract.GetAllRecords(&_IRegistry.CallOpts, name)
}

// GetAllRecords is a free data retrieval call binding the contract method 0x78d573a2.
//
// Solidity: function getAllRecords(string name) view returns((address,uint256)[])
func (_IRegistry *IRegistryCallerSession) GetAllRecords(name string) ([]IRegistryRecord, error) {
	return _IRegistry.Contract.GetAllRecords(&_IRegistry.CallOpts, name)
}

// Names is a free data retrieval call binding the contract method 0x4622ab03.
//
// Solidity: function names(uint256 ) view returns(string)
func (_IRegistry *IRegistryCaller) Names(opts *bind.CallOpts, arg0 *big.Int) (string, error) {
	var out []interface{}
	err := _IRegistry.contract.Call(opts, &out, "names", arg0)
	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err
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

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_IRegistry *IRegistryCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IRegistry.contract.Call(opts, &out, "owner")
	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_IRegistry *IRegistrySession) Owner() (common.Address, error) {
	return _IRegistry.Contract.Owner(&_IRegistry.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_IRegistry *IRegistryCallerSession) Owner() (common.Address, error) {
	return _IRegistry.Contract.Owner(&_IRegistry.CallOpts)
}

// Records is a free data retrieval call binding the contract method 0x3b51650d.
//
// Solidity: function records(string , uint256 ) view returns(address addr, uint256 activation)
func (_IRegistry *IRegistryCaller) Records(opts *bind.CallOpts, arg0 string, arg1 *big.Int) (struct {
	Addr       common.Address
	Activation *big.Int
}, error,
) {
	var out []interface{}
	err := _IRegistry.contract.Call(opts, &out, "records", arg0, arg1)

	outstruct := new(struct {
		Addr       common.Address
		Activation *big.Int
	})

	outstruct.Addr = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.Activation = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	return *outstruct, err
}

// Records is a free data retrieval call binding the contract method 0x3b51650d.
//
// Solidity: function records(string , uint256 ) view returns(address addr, uint256 activation)
func (_IRegistry *IRegistrySession) Records(arg0 string, arg1 *big.Int) (struct {
	Addr       common.Address
	Activation *big.Int
}, error,
) {
	return _IRegistry.Contract.Records(&_IRegistry.CallOpts, arg0, arg1)
}

// Records is a free data retrieval call binding the contract method 0x3b51650d.
//
// Solidity: function records(string , uint256 ) view returns(address addr, uint256 activation)
func (_IRegistry *IRegistryCallerSession) Records(arg0 string, arg1 *big.Int) (struct {
	Addr       common.Address
	Activation *big.Int
}, error,
) {
	return _IRegistry.Contract.Records(&_IRegistry.CallOpts, arg0, arg1)
}

// GetActiveAddr is a paid mutator transaction binding the contract method 0xe2693e3f.
//
// Solidity: function getActiveAddr(string name) returns(address)
func (_IRegistry *IRegistryTransactor) GetActiveAddr(opts *bind.TransactOpts, name string) (*types.Transaction, error) {
	return _IRegistry.contract.Transact(opts, "getActiveAddr", name)
}

// GetActiveAddr is a paid mutator transaction binding the contract method 0xe2693e3f.
//
// Solidity: function getActiveAddr(string name) returns(address)
func (_IRegistry *IRegistrySession) GetActiveAddr(name string) (*types.Transaction, error) {
	return _IRegistry.Contract.GetActiveAddr(&_IRegistry.TransactOpts, name)
}

// GetActiveAddr is a paid mutator transaction binding the contract method 0xe2693e3f.
//
// Solidity: function getActiveAddr(string name) returns(address)
func (_IRegistry *IRegistryTransactorSession) GetActiveAddr(name string) (*types.Transaction, error) {
	return _IRegistry.Contract.GetActiveAddr(&_IRegistry.TransactOpts, name)
}

// Register is a paid mutator transaction binding the contract method 0xd393c871.
//
// Solidity: function register(string name, address addr, uint256 activation) returns()
func (_IRegistry *IRegistryTransactor) Register(opts *bind.TransactOpts, name string, addr common.Address, activation *big.Int) (*types.Transaction, error) {
	return _IRegistry.contract.Transact(opts, "register", name, addr, activation)
}

// Register is a paid mutator transaction binding the contract method 0xd393c871.
//
// Solidity: function register(string name, address addr, uint256 activation) returns()
func (_IRegistry *IRegistrySession) Register(name string, addr common.Address, activation *big.Int) (*types.Transaction, error) {
	return _IRegistry.Contract.Register(&_IRegistry.TransactOpts, name, addr, activation)
}

// Register is a paid mutator transaction binding the contract method 0xd393c871.
//
// Solidity: function register(string name, address addr, uint256 activation) returns()
func (_IRegistry *IRegistryTransactorSession) Register(name string, addr common.Address, activation *big.Int) (*types.Transaction, error) {
	return _IRegistry.Contract.Register(&_IRegistry.TransactOpts, name, addr, activation)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_IRegistry *IRegistryTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _IRegistry.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_IRegistry *IRegistrySession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _IRegistry.Contract.TransferOwnership(&_IRegistry.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_IRegistry *IRegistryTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _IRegistry.Contract.TransferOwnership(&_IRegistry.TransactOpts, newOwner)
}

// IRegistryOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the IRegistry contract.
type IRegistryOwnershipTransferredIterator struct {
	Event *IRegistryOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log      // Log channel receiving the found contract events
	sub  klaytn.Subscription // Subscription for errors, completion and termination
	done bool                // Whether the subscription completed delivering logs
	fail error               // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *IRegistryOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IRegistryOwnershipTransferred)
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
		it.Event = new(IRegistryOwnershipTransferred)
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
func (it *IRegistryOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IRegistryOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IRegistryOwnershipTransferred represents a OwnershipTransferred event raised by the IRegistry contract.
type IRegistryOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_IRegistry *IRegistryFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*IRegistryOwnershipTransferredIterator, error) {
	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _IRegistry.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &IRegistryOwnershipTransferredIterator{contract: _IRegistry.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_IRegistry *IRegistryFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *IRegistryOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {
	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _IRegistry.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IRegistryOwnershipTransferred)
				if err := _IRegistry.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_IRegistry *IRegistryFilterer) ParseOwnershipTransferred(log types.Log) (*IRegistryOwnershipTransferred, error) {
	event := new(IRegistryOwnershipTransferred)
	if err := _IRegistry.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	return event, nil
}

// IRegistryRegisteredIterator is returned from FilterRegistered and is used to iterate over the raw logs and unpacked data for Registered events raised by the IRegistry contract.
type IRegistryRegisteredIterator struct {
	Event *IRegistryRegistered // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log      // Log channel receiving the found contract events
	sub  klaytn.Subscription // Subscription for errors, completion and termination
	done bool                // Whether the subscription completed delivering logs
	fail error               // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *IRegistryRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IRegistryRegistered)
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
		it.Event = new(IRegistryRegistered)
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
func (it *IRegistryRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IRegistryRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IRegistryRegistered represents a Registered event raised by the IRegistry contract.
type IRegistryRegistered struct {
	Name       string
	Addr       common.Address
	Activation *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterRegistered is a free log retrieval operation binding the contract event 0x142e1fdac7ecccbc62af925f0b4039db26847b625602e56b1421dfbc8a0e4f30.
//
// Solidity: event Registered(string name, address indexed addr, uint256 indexed activation)
func (_IRegistry *IRegistryFilterer) FilterRegistered(opts *bind.FilterOpts, addr []common.Address, activation []*big.Int) (*IRegistryRegisteredIterator, error) {
	var addrRule []interface{}
	for _, addrItem := range addr {
		addrRule = append(addrRule, addrItem)
	}
	var activationRule []interface{}
	for _, activationItem := range activation {
		activationRule = append(activationRule, activationItem)
	}

	logs, sub, err := _IRegistry.contract.FilterLogs(opts, "Registered", addrRule, activationRule)
	if err != nil {
		return nil, err
	}
	return &IRegistryRegisteredIterator{contract: _IRegistry.contract, event: "Registered", logs: logs, sub: sub}, nil
}

// WatchRegistered is a free log subscription operation binding the contract event 0x142e1fdac7ecccbc62af925f0b4039db26847b625602e56b1421dfbc8a0e4f30.
//
// Solidity: event Registered(string name, address indexed addr, uint256 indexed activation)
func (_IRegistry *IRegistryFilterer) WatchRegistered(opts *bind.WatchOpts, sink chan<- *IRegistryRegistered, addr []common.Address, activation []*big.Int) (event.Subscription, error) {
	var addrRule []interface{}
	for _, addrItem := range addr {
		addrRule = append(addrRule, addrItem)
	}
	var activationRule []interface{}
	for _, activationItem := range activation {
		activationRule = append(activationRule, activationItem)
	}

	logs, sub, err := _IRegistry.contract.WatchLogs(opts, "Registered", addrRule, activationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IRegistryRegistered)
				if err := _IRegistry.contract.UnpackLog(event, "Registered", log); err != nil {
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

// ParseRegistered is a log parse operation binding the contract event 0x142e1fdac7ecccbc62af925f0b4039db26847b625602e56b1421dfbc8a0e4f30.
//
// Solidity: event Registered(string name, address indexed addr, uint256 indexed activation)
func (_IRegistry *IRegistryFilterer) ParseRegistered(log types.Log) (*IRegistryRegistered, error) {
	event := new(IRegistryRegistered)
	if err := _IRegistry.contract.UnpackLog(event, "Registered", log); err != nil {
		return nil, err
	}
	return event, nil
}

// InitializableMetaData contains all meta data concerning the Initializable contract.
var InitializableMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"}]",
}

// InitializableABI is the input ABI used to generate the binding from.
// Deprecated: Use InitializableMetaData.ABI instead.
var InitializableABI = InitializableMetaData.ABI

// InitializableBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const InitializableBinRuntime = ``

// Initializable is an auto generated Go binding around a Klaytn contract.
type Initializable struct {
	InitializableCaller     // Read-only binding to the contract
	InitializableTransactor // Write-only binding to the contract
	InitializableFilterer   // Log filterer for contract events
}

// InitializableCaller is an auto generated read-only Go binding around a Klaytn contract.
type InitializableCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// InitializableTransactor is an auto generated write-only Go binding around a Klaytn contract.
type InitializableTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// InitializableFilterer is an auto generated log filtering Go binding around a Klaytn contract events.
type InitializableFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// InitializableSession is an auto generated Go binding around a Klaytn contract,
// with pre-set call and transact options.
type InitializableSession struct {
	Contract     *Initializable    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// InitializableCallerSession is an auto generated read-only Go binding around a Klaytn contract,
// with pre-set call options.
type InitializableCallerSession struct {
	Contract *InitializableCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// InitializableTransactorSession is an auto generated write-only Go binding around a Klaytn contract,
// with pre-set transact options.
type InitializableTransactorSession struct {
	Contract     *InitializableTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// InitializableRaw is an auto generated low-level Go binding around a Klaytn contract.
type InitializableRaw struct {
	Contract *Initializable // Generic contract binding to access the raw methods on
}

// InitializableCallerRaw is an auto generated low-level read-only Go binding around a Klaytn contract.
type InitializableCallerRaw struct {
	Contract *InitializableCaller // Generic read-only contract binding to access the raw methods on
}

// InitializableTransactorRaw is an auto generated low-level write-only Go binding around a Klaytn contract.
type InitializableTransactorRaw struct {
	Contract *InitializableTransactor // Generic write-only contract binding to access the raw methods on
}

// NewInitializable creates a new instance of Initializable, bound to a specific deployed contract.
func NewInitializable(address common.Address, backend bind.ContractBackend) (*Initializable, error) {
	contract, err := bindInitializable(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Initializable{InitializableCaller: InitializableCaller{contract: contract}, InitializableTransactor: InitializableTransactor{contract: contract}, InitializableFilterer: InitializableFilterer{contract: contract}}, nil
}

// NewInitializableCaller creates a new read-only instance of Initializable, bound to a specific deployed contract.
func NewInitializableCaller(address common.Address, caller bind.ContractCaller) (*InitializableCaller, error) {
	contract, err := bindInitializable(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &InitializableCaller{contract: contract}, nil
}

// NewInitializableTransactor creates a new write-only instance of Initializable, bound to a specific deployed contract.
func NewInitializableTransactor(address common.Address, transactor bind.ContractTransactor) (*InitializableTransactor, error) {
	contract, err := bindInitializable(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &InitializableTransactor{contract: contract}, nil
}

// NewInitializableFilterer creates a new log filterer instance of Initializable, bound to a specific deployed contract.
func NewInitializableFilterer(address common.Address, filterer bind.ContractFilterer) (*InitializableFilterer, error) {
	contract, err := bindInitializable(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &InitializableFilterer{contract: contract}, nil
}

// bindInitializable binds a generic wrapper to an already deployed contract.
func bindInitializable(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := InitializableMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Initializable *InitializableRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Initializable.Contract.InitializableCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Initializable *InitializableRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Initializable.Contract.InitializableTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Initializable *InitializableRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Initializable.Contract.InitializableTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Initializable *InitializableCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Initializable.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Initializable *InitializableTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Initializable.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Initializable *InitializableTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Initializable.Contract.contract.Transact(opts, method, params...)
}

// InitializableInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the Initializable contract.
type InitializableInitializedIterator struct {
	Event *InitializableInitialized // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log      // Log channel receiving the found contract events
	sub  klaytn.Subscription // Subscription for errors, completion and termination
	done bool                // Whether the subscription completed delivering logs
	fail error               // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *InitializableInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(InitializableInitialized)
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
		it.Event = new(InitializableInitialized)
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
func (it *InitializableInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *InitializableInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// InitializableInitialized represents a Initialized event raised by the Initializable contract.
type InitializableInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_Initializable *InitializableFilterer) FilterInitialized(opts *bind.FilterOpts) (*InitializableInitializedIterator, error) {
	logs, sub, err := _Initializable.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &InitializableInitializedIterator{contract: _Initializable.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_Initializable *InitializableFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *InitializableInitialized) (event.Subscription, error) {
	logs, sub, err := _Initializable.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(InitializableInitialized)
				if err := _Initializable.contract.UnpackLog(event, "Initialized", log); err != nil {
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

// ParseInitialized is a log parse operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_Initializable *InitializableFilterer) ParseInitialized(log types.Log) (*InitializableInitialized, error) {
	event := new(InitializableInitialized)
	if err := _Initializable.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	return event, nil
}

// KIP113MetaData contains all meta data concerning the KIP113 contract.
var KIP113MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"previousAdmin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"AdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"beacon\",\"type\":\"address\"}],\"name\":\"BeaconUpgraded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"cnNodeId\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"publicKey\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"pop\",\"type\":\"bytes\"}],\"name\":\"Registered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"cnNodeId\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"publicKey\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"pop\",\"type\":\"bytes\"}],\"name\":\"Unregistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"ZERO48HASH\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"ZERO96HASH\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"abook\",\"outputs\":[{\"internalType\":\"contractIAddressBook\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"allNodeIds\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllBlsInfo\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"nodeIdList\",\"type\":\"address[]\"},{\"components\":[{\"internalType\":\"bytes\",\"name\":\"publicKey\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"pop\",\"type\":\"bytes\"}],\"internalType\":\"structIKIP113.BlsPublicKeyInfo[]\",\"name\":\"pubkeyList\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"proxiableUUID\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"record\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"publicKey\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"pop\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"cnNodeId\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"publicKey\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"pop\",\"type\":\"bytes\"}],\"name\":\"register\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"cnNodeId\",\"type\":\"address\"}],\"name\":\"unregister\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"}],\"name\":\"upgradeTo\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"upgradeToAndCall\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"6fc522c6": "ZERO48HASH()",
		"20abd458": "ZERO96HASH()",
		"829d639d": "abook()",
		"a5834971": "allNodeIds(uint256)",
		"6968b53f": "getAllBlsInfo()",
		"8129fc1c": "initialize()",
		"8da5cb5b": "owner()",
		"52d1902d": "proxiableUUID()",
		"3465d6d5": "record(address)",
		"786cd4d7": "register(address,bytes,bytes)",
		"715018a6": "renounceOwnership()",
		"f2fde38b": "transferOwnership(address)",
		"2ec2c246": "unregister(address)",
		"3659cfe6": "upgradeTo(address)",
		"4f1ef286": "upgradeToAndCall(address,bytes)",
	},
	Bin: "0x60a06040523060805234801561001457600080fd5b5061001d610022565b6100e1565b600054610100900460ff161561008e5760405162461bcd60e51b815260206004820152602760248201527f496e697469616c697a61626c653a20636f6e747261637420697320696e697469604482015266616c697a696e6760c81b606482015260840160405180910390fd5b60005460ff908116146100df576000805460ff191660ff9081179091556040519081527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15b565b608051611ecc61011860003960008181610593015281816105d301528181610672015281816106b201526107450152611ecc6000f3fe6080604052600436106100e85760003560e01c80636fc522c61161008a578063829d639d11610059578063829d639d1461026d5780638da5cb5b1461029b578063a5834971146102b9578063f2fde38b146102d957600080fd5b80636fc522c6146101ef578063715018a614610223578063786cd4d7146102385780638129fc1c1461025857600080fd5b80633659cfe6116100c65780633659cfe6146101845780634f1ef286146101a457806352d1902d146101b75780636968b53f146101cc57600080fd5b806320abd458146100ed5780632ec2c246146101345780633465d6d514610156575b600080fd5b3480156100f957600080fd5b506101217f46700b4d40ac5c35af2c22dda2787a91eb567b06c924a8fb8ae9a05b20c08c2181565b6040519081526020015b60405180910390f35b34801561014057600080fd5b5061015461014f3660046116cb565b6102f9565b005b34801561016257600080fd5b506101766101713660046116cb565b61045d565b60405161012b92919061173f565b34801561019057600080fd5b5061015461019f3660046116cb565b610589565b6101546101b2366004611783565b610668565b3480156101c357600080fd5b50610121610738565b3480156101d857600080fd5b506101e16107eb565b60405161012b929190611847565b3480156101fb57600080fd5b506101217fc980e59163ce244bb4bb6211f48c7b46f88a4f40943e84eb99bdc41e129bd29381565b34801561022f57600080fd5b50610154610aa8565b34801561024457600080fd5b50610154610253366004611955565b610abc565b34801561026457600080fd5b50610154610e30565b34801561027957600080fd5b5061028361040081565b6040516001600160a01b03909116815260200161012b565b3480156102a757600080fd5b506097546001600160a01b0316610283565b3480156102c557600080fd5b506102836102d43660046119d8565b610f48565b3480156102e557600080fd5b506101546102f43660046116cb565b610f72565b610301610fe8565b61030a81611042565b1561035c5760405162461bcd60e51b815260206004820152601a60248201527f434e206973207374696c6c20696e2041646472657373426f6f6b00000000000060448201526064015b60405180910390fd5b6001600160a01b038116600090815260ca60205260409020805461037f906119f1565b90506000036103c75760405162461bcd60e51b815260206004820152601460248201527310d3881a5cc81b9bdd081c9959da5cdd195c995960621b6044820152606401610353565b6103d0816110be565b6001600160a01b038116600090815260ca60205260409081902090517fb98b07c4d52e8d65fa5416810f2746a810eb074b1ac7784e1b07e315c0dfd2d99161041f918491906001820190611aa8565b60405180910390a16001600160a01b038116600090815260ca602052604081209061044a8282611668565b610458600183016000611668565b505050565b60ca60205260009081526040902080548190610478906119f1565b80601f01602080910402602001604051908101604052809291908181526020018280546104a4906119f1565b80156104f15780601f106104c6576101008083540402835291602001916104f1565b820191906000526020600020905b8154815290600101906020018083116104d457829003601f168201915b505050505090806001018054610506906119f1565b80601f0160208091040260200160405190810160405280929190818152602001828054610532906119f1565b801561057f5780601f106105545761010080835404028352916020019161057f565b820191906000526020600020905b81548152906001019060200180831161056257829003601f168201915b5050505050905082565b6001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001630036105d15760405162461bcd60e51b815260040161035390611ade565b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031661061a600080516020611e50833981519152546001600160a01b031690565b6001600160a01b0316146106405760405162461bcd60e51b815260040161035390611b2a565b610649816111c5565b60408051600080825260208201909252610665918391906111cd565b50565b6001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001630036106b05760405162461bcd60e51b815260040161035390611ade565b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03166106f9600080516020611e50833981519152546001600160a01b031690565b6001600160a01b03161461071f5760405162461bcd60e51b815260040161035390611b2a565b610728826111c5565b610734828260016111cd565b5050565b6000306001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016146107d85760405162461bcd60e51b815260206004820152603860248201527f555550535570677261646561626c653a206d757374206e6f742062652063616c60448201527f6c6564207468726f7567682064656c656761746563616c6c00000000000000006064820152608401610353565b50600080516020611e5083398151915290565b60c954606090819067ffffffffffffffff81111561080b5761080b61176d565b604051908082528060200260200182016040528015610834578160200160208202803683370190505b5060c95490925067ffffffffffffffff8111156108535761085361176d565b60405190808252806020026020018201604052801561089857816020015b60408051808201909152606080825260208201528152602001906001900390816108715790505b50905060005b8251811015610aa35760c981815481106108ba576108ba611b76565b9060005260206000200160009054906101000a90046001600160a01b03168382815181106108ea576108ea611b76565b60200260200101906001600160a01b031690816001600160a01b03168152505060ca600060c9838154811061092157610921611b76565b60009182526020808320909101546001600160a01b031683528201929092526040908101909120815180830190925280548290829061095f906119f1565b80601f016020809104026020016040519081016040528092919081815260200182805461098b906119f1565b80156109d85780601f106109ad576101008083540402835291602001916109d8565b820191906000526020600020905b8154815290600101906020018083116109bb57829003601f168201915b505050505081526020016001820180546109f1906119f1565b80601f0160208091040260200160405190810160405280929190818152602001828054610a1d906119f1565b8015610a6a5780601f10610a3f57610100808354040283529160200191610a6a565b820191906000526020600020905b815481529060010190602001808311610a4d57829003601f168201915b505050505081525050828281518110610a8557610a85611b76565b60200260200101819052508080610a9b90611ba2565b91505061089e565b509091565b610ab0610fe8565b610aba6000611338565b565b610ac4610fe8565b838360308114610b165760405162461bcd60e51b815260206004820152601b60248201527f5075626c6963206b6579206d75737420626520343820627974657300000000006044820152606401610353565b6040517fc980e59163ce244bb4bb6211f48c7b46f88a4f40943e84eb99bdc41e129bd29390610b489084908490611bbb565b604051809103902003610b9d5760405162461bcd60e51b815260206004820152601960248201527f5075626c6963206b65792063616e6e6f74206265207a65726f000000000000006044820152606401610353565b838360608114610be65760405162461bcd60e51b8152602060048201526014602482015273506f70206d75737420626520393620627974657360601b6044820152606401610353565b6040517f46700b4d40ac5c35af2c22dda2787a91eb567b06c924a8fb8ae9a05b20c08c2190610c189084908490611bbb565b604051809103902003610c625760405162461bcd60e51b8152602060048201526012602482015271506f702063616e6e6f74206265207a65726f60701b6044820152606401610353565b610c6b89611042565b610cb75760405162461bcd60e51b815260206004820152601e60248201527f636e4e6f64654964206973206e6f7420696e2041646472657373426f6f6b00006044820152606401610353565b6001600160a01b038916600090815260ca602052604090208054610cda906119f1565b9050600003610d2f5760c980546001810182556000919091527f66be4f155c5ef2ebd3772b228f2f00681e4ed5826cdb3b1943cc11ad15ad1d280180546001600160a01b0319166001600160a01b038b161790555b6040805160606020601f8b018190040282018101835291810189815290918291908b908b9081908501838280828437600092019190915250505090825250604080516020601f8a018190048102820181019092528881529181019190899089908190840183828082843760009201829052509390945250506001600160a01b038c16815260ca6020526040902082519091508190610dcd9082611c19565b5060208201516001820190610de29082611c19565b509050507f79c75399e89a1f580d9a6252cb8bdcf4cd80f73b3597c94d845eb52174ad930f8989898989604051610e1d959493929190611d02565b60405180910390a1505050505050505050565b600054610100900460ff1615808015610e505750600054600160ff909116105b80610e6a5750303b158015610e6a575060005460ff166001145b610ecd5760405162461bcd60e51b815260206004820152602e60248201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160448201526d191e481a5b9a5d1a585b1a5e995960921b6064820152608401610353565b6000805460ff191660011790558015610ef0576000805461ff0019166101001790555b610ef861138a565b610f006113b9565b8015610665576000805461ff0019169055604051600181527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a150565b60c98181548110610f5857600080fd5b6000918252602090912001546001600160a01b0316905081565b610f7a610fe8565b6001600160a01b038116610fdf5760405162461bcd60e51b815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201526564647265737360d01b6064820152608401610353565b61066581611338565b6097546001600160a01b03163314610aba5760405162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e65726044820152606401610353565b604051630aabaead60e11b81526001600160a01b0382166004820152600090610400906315575d5a90602401606060405180830381865afa9250505080156110a7575060408051601f3d908101601f191682019092526110a491810190611d46565b60015b6110b357506000919050565b506001949350505050565b60005b60c95481101561073457816001600160a01b031660c982815481106110e8576110e8611b76565b6000918252602090912001546001600160a01b0316036111b35760c9805461111290600190611d93565b8154811061112257611122611b76565b60009182526020909120015460c980546001600160a01b03909216918390811061114e5761114e611b76565b9060005260206000200160006101000a8154816001600160a01b0302191690836001600160a01b0316021790555060c980548061118d5761118d611da6565b600082815260209020810160001990810180546001600160a01b03191690550190555050565b806111bd81611ba2565b9150506110c1565b610665610fe8565b7f4910fdfa16fed3260ed0e7147f7cc6da11a60208b5b9406d12a635614ffd91435460ff161561120057610458836113e0565b826001600160a01b03166352d1902d6040518163ffffffff1660e01b8152600401602060405180830381865afa92505050801561125a575060408051601f3d908101601f1916820190925261125791810190611dbc565b60015b6112bd5760405162461bcd60e51b815260206004820152602e60248201527f45524331393637557067726164653a206e657720696d706c656d656e7461746960448201526d6f6e206973206e6f74205555505360901b6064820152608401610353565b600080516020611e50833981519152811461132c5760405162461bcd60e51b815260206004820152602960248201527f45524331393637557067726164653a20756e737570706f727465642070726f786044820152681a58589b195555525160ba1b6064820152608401610353565b5061045883838361147c565b609780546001600160a01b038381166001600160a01b0319831681179093556040519116919082907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a35050565b600054610100900460ff166113b15760405162461bcd60e51b815260040161035390611dd5565b610aba6114a7565b600054610100900460ff16610aba5760405162461bcd60e51b815260040161035390611dd5565b6001600160a01b0381163b61144d5760405162461bcd60e51b815260206004820152602d60248201527f455243313936373a206e657720696d706c656d656e746174696f6e206973206e60448201526c1bdd08184818dbdb9d1c9858dd609a1b6064820152608401610353565b600080516020611e5083398151915280546001600160a01b0319166001600160a01b0392909216919091179055565b611485836114d7565b6000825111806114925750805b15610458576114a18383611517565b50505050565b600054610100900460ff166114ce5760405162461bcd60e51b815260040161035390611dd5565b610aba33611338565b6114e0816113e0565b6040516001600160a01b038216907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b90600090a250565b606061153c8383604051806060016040528060278152602001611e7060279139611545565b90505b92915050565b6060600080856001600160a01b0316856040516115629190611e20565b600060405180830381855af49150503d806000811461159d576040519150601f19603f3d011682016040523d82523d6000602084013e6115a2565b606091505b50915091506115b3868383876115bd565b9695505050505050565b6060831561162c578251600003611625576001600160a01b0385163b6116255760405162461bcd60e51b815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e74726163740000006044820152606401610353565b5081611636565b611636838361163e565b949350505050565b81511561164e5781518083602001fd5b8060405162461bcd60e51b81526004016103539190611e3c565b508054611674906119f1565b6000825580601f10611684575050565b601f01602090049060005260206000209081019061066591905b808211156116b2576000815560010161169e565b5090565b6001600160a01b038116811461066557600080fd5b6000602082840312156116dd57600080fd5b81356116e8816116b6565b9392505050565b60005b8381101561170a5781810151838201526020016116f2565b50506000910152565b6000815180845261172b8160208601602086016116ef565b601f01601f19169290920160200192915050565b6040815260006117526040830185611713565b82810360208401526117648185611713565b95945050505050565b634e487b7160e01b600052604160045260246000fd5b6000806040838503121561179657600080fd5b82356117a1816116b6565b9150602083013567ffffffffffffffff808211156117be57600080fd5b818501915085601f8301126117d257600080fd5b8135818111156117e4576117e461176d565b604051601f8201601f19908116603f0116810190838211818310171561180c5761180c61176d565b8160405282815288602084870101111561182557600080fd5b8260208601602083013760006020848301015280955050505050509250929050565b60408082528351828201819052600091906020906060850190828801855b8281101561188a5781516001600160a01b031684529284019290840190600101611865565b50505084810382860152855180825282820190600581901b8301840188850160005b838110156118fc57858303601f19018552815180518985526118d08a860182611713565b91890151858303868b01529190506118e88183611713565b9689019694505050908601906001016118ac565b50909a9950505050505050505050565b60008083601f84011261191e57600080fd5b50813567ffffffffffffffff81111561193657600080fd5b60208301915083602082850101111561194e57600080fd5b9250929050565b60008060008060006060868803121561196d57600080fd5b8535611978816116b6565b9450602086013567ffffffffffffffff8082111561199557600080fd5b6119a189838a0161190c565b909650945060408801359150808211156119ba57600080fd5b506119c78882890161190c565b969995985093965092949392505050565b6000602082840312156119ea57600080fd5b5035919050565b600181811c90821680611a0557607f821691505b602082108103611a2557634e487b7160e01b600052602260045260246000fd5b50919050565b60008154611a38816119f1565b808552602060018381168015611a555760018114611a6f57611a9d565b60ff1985168884015283151560051b880183019550611a9d565b866000528260002060005b85811015611a955781548a8201860152908301908401611a7a565b890184019650505b505050505092915050565b6001600160a01b0384168152606060208201819052600090611acc90830185611a2b565b82810360408401526115b38185611a2b565b6020808252602c908201527f46756e6374696f6e206d7573742062652063616c6c6564207468726f7567682060408201526b19195b1959d85d1958d85b1b60a21b606082015260800190565b6020808252602c908201527f46756e6374696f6e206d7573742062652063616c6c6564207468726f7567682060408201526b6163746976652070726f787960a01b606082015260800190565b634e487b7160e01b600052603260045260246000fd5b634e487b7160e01b600052601160045260246000fd5b600060018201611bb457611bb4611b8c565b5060010190565b8183823760009101908152919050565b601f82111561045857600081815260208120601f850160051c81016020861015611bf25750805b601f850160051c820191505b81811015611c1157828155600101611bfe565b505050505050565b815167ffffffffffffffff811115611c3357611c3361176d565b611c4781611c4184546119f1565b84611bcb565b602080601f831160018114611c7c5760008415611c645750858301515b600019600386901b1c1916600185901b178555611c11565b600085815260208120601f198616915b82811015611cab57888601518255948401946001909101908401611c8c565b5085821015611cc95787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b81835281816020850137506000828201602090810191909152601f909101601f19169091010190565b6001600160a01b0386168152606060208201819052600090611d279083018688611cd9565b8281036040840152611d3a818587611cd9565b98975050505050505050565b600080600060608486031215611d5b57600080fd5b8351611d66816116b6565b6020850151909350611d77816116b6565b6040850151909250611d88816116b6565b809150509250925092565b8181038181111561153f5761153f611b8c565b634e487b7160e01b600052603160045260246000fd5b600060208284031215611dce57600080fd5b5051919050565b6020808252602b908201527f496e697469616c697a61626c653a20636f6e7472616374206973206e6f74206960408201526a6e697469616c697a696e6760a81b606082015260800190565b60008251611e328184602087016116ef565b9190910192915050565b60208152600061153c602083018461171356fe360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc416464726573733a206c6f772d6c6576656c2064656c65676174652063616c6c206661696c6564a2646970667358221220d02b98764eebd34b392597d53a0f207a671d0bcaba08bb568c5dabdf106de2cd64736f6c63430008130033",
}

// KIP113ABI is the input ABI used to generate the binding from.
// Deprecated: Use KIP113MetaData.ABI instead.
var KIP113ABI = KIP113MetaData.ABI

// KIP113BinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const KIP113BinRuntime = `6080604052600436106100e85760003560e01c80636fc522c61161008a578063829d639d11610059578063829d639d1461026d5780638da5cb5b1461029b578063a5834971146102b9578063f2fde38b146102d957600080fd5b80636fc522c6146101ef578063715018a614610223578063786cd4d7146102385780638129fc1c1461025857600080fd5b80633659cfe6116100c65780633659cfe6146101845780634f1ef286146101a457806352d1902d146101b75780636968b53f146101cc57600080fd5b806320abd458146100ed5780632ec2c246146101345780633465d6d514610156575b600080fd5b3480156100f957600080fd5b506101217f46700b4d40ac5c35af2c22dda2787a91eb567b06c924a8fb8ae9a05b20c08c2181565b6040519081526020015b60405180910390f35b34801561014057600080fd5b5061015461014f3660046116cb565b6102f9565b005b34801561016257600080fd5b506101766101713660046116cb565b61045d565b60405161012b92919061173f565b34801561019057600080fd5b5061015461019f3660046116cb565b610589565b6101546101b2366004611783565b610668565b3480156101c357600080fd5b50610121610738565b3480156101d857600080fd5b506101e16107eb565b60405161012b929190611847565b3480156101fb57600080fd5b506101217fc980e59163ce244bb4bb6211f48c7b46f88a4f40943e84eb99bdc41e129bd29381565b34801561022f57600080fd5b50610154610aa8565b34801561024457600080fd5b50610154610253366004611955565b610abc565b34801561026457600080fd5b50610154610e30565b34801561027957600080fd5b5061028361040081565b6040516001600160a01b03909116815260200161012b565b3480156102a757600080fd5b506097546001600160a01b0316610283565b3480156102c557600080fd5b506102836102d43660046119d8565b610f48565b3480156102e557600080fd5b506101546102f43660046116cb565b610f72565b610301610fe8565b61030a81611042565b1561035c5760405162461bcd60e51b815260206004820152601a60248201527f434e206973207374696c6c20696e2041646472657373426f6f6b00000000000060448201526064015b60405180910390fd5b6001600160a01b038116600090815260ca60205260409020805461037f906119f1565b90506000036103c75760405162461bcd60e51b815260206004820152601460248201527310d3881a5cc81b9bdd081c9959da5cdd195c995960621b6044820152606401610353565b6103d0816110be565b6001600160a01b038116600090815260ca60205260409081902090517fb98b07c4d52e8d65fa5416810f2746a810eb074b1ac7784e1b07e315c0dfd2d99161041f918491906001820190611aa8565b60405180910390a16001600160a01b038116600090815260ca602052604081209061044a8282611668565b610458600183016000611668565b505050565b60ca60205260009081526040902080548190610478906119f1565b80601f01602080910402602001604051908101604052809291908181526020018280546104a4906119f1565b80156104f15780601f106104c6576101008083540402835291602001916104f1565b820191906000526020600020905b8154815290600101906020018083116104d457829003601f168201915b505050505090806001018054610506906119f1565b80601f0160208091040260200160405190810160405280929190818152602001828054610532906119f1565b801561057f5780601f106105545761010080835404028352916020019161057f565b820191906000526020600020905b81548152906001019060200180831161056257829003601f168201915b5050505050905082565b6001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001630036105d15760405162461bcd60e51b815260040161035390611ade565b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031661061a600080516020611e50833981519152546001600160a01b031690565b6001600160a01b0316146106405760405162461bcd60e51b815260040161035390611b2a565b610649816111c5565b60408051600080825260208201909252610665918391906111cd565b50565b6001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001630036106b05760405162461bcd60e51b815260040161035390611ade565b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03166106f9600080516020611e50833981519152546001600160a01b031690565b6001600160a01b03161461071f5760405162461bcd60e51b815260040161035390611b2a565b610728826111c5565b610734828260016111cd565b5050565b6000306001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016146107d85760405162461bcd60e51b815260206004820152603860248201527f555550535570677261646561626c653a206d757374206e6f742062652063616c60448201527f6c6564207468726f7567682064656c656761746563616c6c00000000000000006064820152608401610353565b50600080516020611e5083398151915290565b60c954606090819067ffffffffffffffff81111561080b5761080b61176d565b604051908082528060200260200182016040528015610834578160200160208202803683370190505b5060c95490925067ffffffffffffffff8111156108535761085361176d565b60405190808252806020026020018201604052801561089857816020015b60408051808201909152606080825260208201528152602001906001900390816108715790505b50905060005b8251811015610aa35760c981815481106108ba576108ba611b76565b9060005260206000200160009054906101000a90046001600160a01b03168382815181106108ea576108ea611b76565b60200260200101906001600160a01b031690816001600160a01b03168152505060ca600060c9838154811061092157610921611b76565b60009182526020808320909101546001600160a01b031683528201929092526040908101909120815180830190925280548290829061095f906119f1565b80601f016020809104026020016040519081016040528092919081815260200182805461098b906119f1565b80156109d85780601f106109ad576101008083540402835291602001916109d8565b820191906000526020600020905b8154815290600101906020018083116109bb57829003601f168201915b505050505081526020016001820180546109f1906119f1565b80601f0160208091040260200160405190810160405280929190818152602001828054610a1d906119f1565b8015610a6a5780601f10610a3f57610100808354040283529160200191610a6a565b820191906000526020600020905b815481529060010190602001808311610a4d57829003601f168201915b505050505081525050828281518110610a8557610a85611b76565b60200260200101819052508080610a9b90611ba2565b91505061089e565b509091565b610ab0610fe8565b610aba6000611338565b565b610ac4610fe8565b838360308114610b165760405162461bcd60e51b815260206004820152601b60248201527f5075626c6963206b6579206d75737420626520343820627974657300000000006044820152606401610353565b6040517fc980e59163ce244bb4bb6211f48c7b46f88a4f40943e84eb99bdc41e129bd29390610b489084908490611bbb565b604051809103902003610b9d5760405162461bcd60e51b815260206004820152601960248201527f5075626c6963206b65792063616e6e6f74206265207a65726f000000000000006044820152606401610353565b838360608114610be65760405162461bcd60e51b8152602060048201526014602482015273506f70206d75737420626520393620627974657360601b6044820152606401610353565b6040517f46700b4d40ac5c35af2c22dda2787a91eb567b06c924a8fb8ae9a05b20c08c2190610c189084908490611bbb565b604051809103902003610c625760405162461bcd60e51b8152602060048201526012602482015271506f702063616e6e6f74206265207a65726f60701b6044820152606401610353565b610c6b89611042565b610cb75760405162461bcd60e51b815260206004820152601e60248201527f636e4e6f64654964206973206e6f7420696e2041646472657373426f6f6b00006044820152606401610353565b6001600160a01b038916600090815260ca602052604090208054610cda906119f1565b9050600003610d2f5760c980546001810182556000919091527f66be4f155c5ef2ebd3772b228f2f00681e4ed5826cdb3b1943cc11ad15ad1d280180546001600160a01b0319166001600160a01b038b161790555b6040805160606020601f8b018190040282018101835291810189815290918291908b908b9081908501838280828437600092019190915250505090825250604080516020601f8a018190048102820181019092528881529181019190899089908190840183828082843760009201829052509390945250506001600160a01b038c16815260ca6020526040902082519091508190610dcd9082611c19565b5060208201516001820190610de29082611c19565b509050507f79c75399e89a1f580d9a6252cb8bdcf4cd80f73b3597c94d845eb52174ad930f8989898989604051610e1d959493929190611d02565b60405180910390a1505050505050505050565b600054610100900460ff1615808015610e505750600054600160ff909116105b80610e6a5750303b158015610e6a575060005460ff166001145b610ecd5760405162461bcd60e51b815260206004820152602e60248201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160448201526d191e481a5b9a5d1a585b1a5e995960921b6064820152608401610353565b6000805460ff191660011790558015610ef0576000805461ff0019166101001790555b610ef861138a565b610f006113b9565b8015610665576000805461ff0019169055604051600181527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a150565b60c98181548110610f5857600080fd5b6000918252602090912001546001600160a01b0316905081565b610f7a610fe8565b6001600160a01b038116610fdf5760405162461bcd60e51b815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201526564647265737360d01b6064820152608401610353565b61066581611338565b6097546001600160a01b03163314610aba5760405162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e65726044820152606401610353565b604051630aabaead60e11b81526001600160a01b0382166004820152600090610400906315575d5a90602401606060405180830381865afa9250505080156110a7575060408051601f3d908101601f191682019092526110a491810190611d46565b60015b6110b357506000919050565b506001949350505050565b60005b60c95481101561073457816001600160a01b031660c982815481106110e8576110e8611b76565b6000918252602090912001546001600160a01b0316036111b35760c9805461111290600190611d93565b8154811061112257611122611b76565b60009182526020909120015460c980546001600160a01b03909216918390811061114e5761114e611b76565b9060005260206000200160006101000a8154816001600160a01b0302191690836001600160a01b0316021790555060c980548061118d5761118d611da6565b600082815260209020810160001990810180546001600160a01b03191690550190555050565b806111bd81611ba2565b9150506110c1565b610665610fe8565b7f4910fdfa16fed3260ed0e7147f7cc6da11a60208b5b9406d12a635614ffd91435460ff161561120057610458836113e0565b826001600160a01b03166352d1902d6040518163ffffffff1660e01b8152600401602060405180830381865afa92505050801561125a575060408051601f3d908101601f1916820190925261125791810190611dbc565b60015b6112bd5760405162461bcd60e51b815260206004820152602e60248201527f45524331393637557067726164653a206e657720696d706c656d656e7461746960448201526d6f6e206973206e6f74205555505360901b6064820152608401610353565b600080516020611e50833981519152811461132c5760405162461bcd60e51b815260206004820152602960248201527f45524331393637557067726164653a20756e737570706f727465642070726f786044820152681a58589b195555525160ba1b6064820152608401610353565b5061045883838361147c565b609780546001600160a01b038381166001600160a01b0319831681179093556040519116919082907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a35050565b600054610100900460ff166113b15760405162461bcd60e51b815260040161035390611dd5565b610aba6114a7565b600054610100900460ff16610aba5760405162461bcd60e51b815260040161035390611dd5565b6001600160a01b0381163b61144d5760405162461bcd60e51b815260206004820152602d60248201527f455243313936373a206e657720696d706c656d656e746174696f6e206973206e60448201526c1bdd08184818dbdb9d1c9858dd609a1b6064820152608401610353565b600080516020611e5083398151915280546001600160a01b0319166001600160a01b0392909216919091179055565b611485836114d7565b6000825111806114925750805b15610458576114a18383611517565b50505050565b600054610100900460ff166114ce5760405162461bcd60e51b815260040161035390611dd5565b610aba33611338565b6114e0816113e0565b6040516001600160a01b038216907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b90600090a250565b606061153c8383604051806060016040528060278152602001611e7060279139611545565b90505b92915050565b6060600080856001600160a01b0316856040516115629190611e20565b600060405180830381855af49150503d806000811461159d576040519150601f19603f3d011682016040523d82523d6000602084013e6115a2565b606091505b50915091506115b3868383876115bd565b9695505050505050565b6060831561162c578251600003611625576001600160a01b0385163b6116255760405162461bcd60e51b815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e74726163740000006044820152606401610353565b5081611636565b611636838361163e565b949350505050565b81511561164e5781518083602001fd5b8060405162461bcd60e51b81526004016103539190611e3c565b508054611674906119f1565b6000825580601f10611684575050565b601f01602090049060005260206000209081019061066591905b808211156116b2576000815560010161169e565b5090565b6001600160a01b038116811461066557600080fd5b6000602082840312156116dd57600080fd5b81356116e8816116b6565b9392505050565b60005b8381101561170a5781810151838201526020016116f2565b50506000910152565b6000815180845261172b8160208601602086016116ef565b601f01601f19169290920160200192915050565b6040815260006117526040830185611713565b82810360208401526117648185611713565b95945050505050565b634e487b7160e01b600052604160045260246000fd5b6000806040838503121561179657600080fd5b82356117a1816116b6565b9150602083013567ffffffffffffffff808211156117be57600080fd5b818501915085601f8301126117d257600080fd5b8135818111156117e4576117e461176d565b604051601f8201601f19908116603f0116810190838211818310171561180c5761180c61176d565b8160405282815288602084870101111561182557600080fd5b8260208601602083013760006020848301015280955050505050509250929050565b60408082528351828201819052600091906020906060850190828801855b8281101561188a5781516001600160a01b031684529284019290840190600101611865565b50505084810382860152855180825282820190600581901b8301840188850160005b838110156118fc57858303601f19018552815180518985526118d08a860182611713565b91890151858303868b01529190506118e88183611713565b9689019694505050908601906001016118ac565b50909a9950505050505050505050565b60008083601f84011261191e57600080fd5b50813567ffffffffffffffff81111561193657600080fd5b60208301915083602082850101111561194e57600080fd5b9250929050565b60008060008060006060868803121561196d57600080fd5b8535611978816116b6565b9450602086013567ffffffffffffffff8082111561199557600080fd5b6119a189838a0161190c565b909650945060408801359150808211156119ba57600080fd5b506119c78882890161190c565b969995985093965092949392505050565b6000602082840312156119ea57600080fd5b5035919050565b600181811c90821680611a0557607f821691505b602082108103611a2557634e487b7160e01b600052602260045260246000fd5b50919050565b60008154611a38816119f1565b808552602060018381168015611a555760018114611a6f57611a9d565b60ff1985168884015283151560051b880183019550611a9d565b866000528260002060005b85811015611a955781548a8201860152908301908401611a7a565b890184019650505b505050505092915050565b6001600160a01b0384168152606060208201819052600090611acc90830185611a2b565b82810360408401526115b38185611a2b565b6020808252602c908201527f46756e6374696f6e206d7573742062652063616c6c6564207468726f7567682060408201526b19195b1959d85d1958d85b1b60a21b606082015260800190565b6020808252602c908201527f46756e6374696f6e206d7573742062652063616c6c6564207468726f7567682060408201526b6163746976652070726f787960a01b606082015260800190565b634e487b7160e01b600052603260045260246000fd5b634e487b7160e01b600052601160045260246000fd5b600060018201611bb457611bb4611b8c565b5060010190565b8183823760009101908152919050565b601f82111561045857600081815260208120601f850160051c81016020861015611bf25750805b601f850160051c820191505b81811015611c1157828155600101611bfe565b505050505050565b815167ffffffffffffffff811115611c3357611c3361176d565b611c4781611c4184546119f1565b84611bcb565b602080601f831160018114611c7c5760008415611c645750858301515b600019600386901b1c1916600185901b178555611c11565b600085815260208120601f198616915b82811015611cab57888601518255948401946001909101908401611c8c565b5085821015611cc95787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b81835281816020850137506000828201602090810191909152601f909101601f19169091010190565b6001600160a01b0386168152606060208201819052600090611d279083018688611cd9565b8281036040840152611d3a818587611cd9565b98975050505050505050565b600080600060608486031215611d5b57600080fd5b8351611d66816116b6565b6020850151909350611d77816116b6565b6040850151909250611d88816116b6565b809150509250925092565b8181038181111561153f5761153f611b8c565b634e487b7160e01b600052603160045260246000fd5b600060208284031215611dce57600080fd5b5051919050565b6020808252602b908201527f496e697469616c697a61626c653a20636f6e7472616374206973206e6f74206960408201526a6e697469616c697a696e6760a81b606082015260800190565b60008251611e328184602087016116ef565b9190910192915050565b60208152600061153c602083018461171356fe360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc416464726573733a206c6f772d6c6576656c2064656c65676174652063616c6c206661696c6564a2646970667358221220d02b98764eebd34b392597d53a0f207a671d0bcaba08bb568c5dabdf106de2cd64736f6c63430008130033`

// KIP113FuncSigs maps the 4-byte function signature to its string representation.
// Deprecated: Use KIP113MetaData.Sigs instead.
var KIP113FuncSigs = KIP113MetaData.Sigs

// KIP113Bin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use KIP113MetaData.Bin instead.
var KIP113Bin = KIP113MetaData.Bin

// DeployKIP113 deploys a new Klaytn contract, binding an instance of KIP113 to it.
func DeployKIP113(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *KIP113, error) {
	parsed, err := KIP113MetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(KIP113Bin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &KIP113{KIP113Caller: KIP113Caller{contract: contract}, KIP113Transactor: KIP113Transactor{contract: contract}, KIP113Filterer: KIP113Filterer{contract: contract}}, nil
}

// KIP113 is an auto generated Go binding around a Klaytn contract.
type KIP113 struct {
	KIP113Caller     // Read-only binding to the contract
	KIP113Transactor // Write-only binding to the contract
	KIP113Filterer   // Log filterer for contract events
}

// KIP113Caller is an auto generated read-only Go binding around a Klaytn contract.
type KIP113Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// KIP113Transactor is an auto generated write-only Go binding around a Klaytn contract.
type KIP113Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// KIP113Filterer is an auto generated log filtering Go binding around a Klaytn contract events.
type KIP113Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// KIP113Session is an auto generated Go binding around a Klaytn contract,
// with pre-set call and transact options.
type KIP113Session struct {
	Contract     *KIP113           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// KIP113CallerSession is an auto generated read-only Go binding around a Klaytn contract,
// with pre-set call options.
type KIP113CallerSession struct {
	Contract *KIP113Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// KIP113TransactorSession is an auto generated write-only Go binding around a Klaytn contract,
// with pre-set transact options.
type KIP113TransactorSession struct {
	Contract     *KIP113Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// KIP113Raw is an auto generated low-level Go binding around a Klaytn contract.
type KIP113Raw struct {
	Contract *KIP113 // Generic contract binding to access the raw methods on
}

// KIP113CallerRaw is an auto generated low-level read-only Go binding around a Klaytn contract.
type KIP113CallerRaw struct {
	Contract *KIP113Caller // Generic read-only contract binding to access the raw methods on
}

// KIP113TransactorRaw is an auto generated low-level write-only Go binding around a Klaytn contract.
type KIP113TransactorRaw struct {
	Contract *KIP113Transactor // Generic write-only contract binding to access the raw methods on
}

// NewKIP113 creates a new instance of KIP113, bound to a specific deployed contract.
func NewKIP113(address common.Address, backend bind.ContractBackend) (*KIP113, error) {
	contract, err := bindKIP113(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &KIP113{KIP113Caller: KIP113Caller{contract: contract}, KIP113Transactor: KIP113Transactor{contract: contract}, KIP113Filterer: KIP113Filterer{contract: contract}}, nil
}

// NewKIP113Caller creates a new read-only instance of KIP113, bound to a specific deployed contract.
func NewKIP113Caller(address common.Address, caller bind.ContractCaller) (*KIP113Caller, error) {
	contract, err := bindKIP113(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &KIP113Caller{contract: contract}, nil
}

// NewKIP113Transactor creates a new write-only instance of KIP113, bound to a specific deployed contract.
func NewKIP113Transactor(address common.Address, transactor bind.ContractTransactor) (*KIP113Transactor, error) {
	contract, err := bindKIP113(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &KIP113Transactor{contract: contract}, nil
}

// NewKIP113Filterer creates a new log filterer instance of KIP113, bound to a specific deployed contract.
func NewKIP113Filterer(address common.Address, filterer bind.ContractFilterer) (*KIP113Filterer, error) {
	contract, err := bindKIP113(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &KIP113Filterer{contract: contract}, nil
}

// bindKIP113 binds a generic wrapper to an already deployed contract.
func bindKIP113(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := KIP113MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_KIP113 *KIP113Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _KIP113.Contract.KIP113Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_KIP113 *KIP113Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KIP113.Contract.KIP113Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_KIP113 *KIP113Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KIP113.Contract.KIP113Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_KIP113 *KIP113CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _KIP113.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_KIP113 *KIP113TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KIP113.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_KIP113 *KIP113TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KIP113.Contract.contract.Transact(opts, method, params...)
}

// ZERO48HASH is a free data retrieval call binding the contract method 0x6fc522c6.
//
// Solidity: function ZERO48HASH() view returns(bytes32)
func (_KIP113 *KIP113Caller) ZERO48HASH(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _KIP113.contract.Call(opts, &out, "ZERO48HASH")
	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err
}

// ZERO48HASH is a free data retrieval call binding the contract method 0x6fc522c6.
//
// Solidity: function ZERO48HASH() view returns(bytes32)
func (_KIP113 *KIP113Session) ZERO48HASH() ([32]byte, error) {
	return _KIP113.Contract.ZERO48HASH(&_KIP113.CallOpts)
}

// ZERO48HASH is a free data retrieval call binding the contract method 0x6fc522c6.
//
// Solidity: function ZERO48HASH() view returns(bytes32)
func (_KIP113 *KIP113CallerSession) ZERO48HASH() ([32]byte, error) {
	return _KIP113.Contract.ZERO48HASH(&_KIP113.CallOpts)
}

// ZERO96HASH is a free data retrieval call binding the contract method 0x20abd458.
//
// Solidity: function ZERO96HASH() view returns(bytes32)
func (_KIP113 *KIP113Caller) ZERO96HASH(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _KIP113.contract.Call(opts, &out, "ZERO96HASH")
	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err
}

// ZERO96HASH is a free data retrieval call binding the contract method 0x20abd458.
//
// Solidity: function ZERO96HASH() view returns(bytes32)
func (_KIP113 *KIP113Session) ZERO96HASH() ([32]byte, error) {
	return _KIP113.Contract.ZERO96HASH(&_KIP113.CallOpts)
}

// ZERO96HASH is a free data retrieval call binding the contract method 0x20abd458.
//
// Solidity: function ZERO96HASH() view returns(bytes32)
func (_KIP113 *KIP113CallerSession) ZERO96HASH() ([32]byte, error) {
	return _KIP113.Contract.ZERO96HASH(&_KIP113.CallOpts)
}

// Abook is a free data retrieval call binding the contract method 0x829d639d.
//
// Solidity: function abook() view returns(address)
func (_KIP113 *KIP113Caller) Abook(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KIP113.contract.Call(opts, &out, "abook")
	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err
}

// Abook is a free data retrieval call binding the contract method 0x829d639d.
//
// Solidity: function abook() view returns(address)
func (_KIP113 *KIP113Session) Abook() (common.Address, error) {
	return _KIP113.Contract.Abook(&_KIP113.CallOpts)
}

// Abook is a free data retrieval call binding the contract method 0x829d639d.
//
// Solidity: function abook() view returns(address)
func (_KIP113 *KIP113CallerSession) Abook() (common.Address, error) {
	return _KIP113.Contract.Abook(&_KIP113.CallOpts)
}

// AllNodeIds is a free data retrieval call binding the contract method 0xa5834971.
//
// Solidity: function allNodeIds(uint256 ) view returns(address)
func (_KIP113 *KIP113Caller) AllNodeIds(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var out []interface{}
	err := _KIP113.contract.Call(opts, &out, "allNodeIds", arg0)
	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err
}

// AllNodeIds is a free data retrieval call binding the contract method 0xa5834971.
//
// Solidity: function allNodeIds(uint256 ) view returns(address)
func (_KIP113 *KIP113Session) AllNodeIds(arg0 *big.Int) (common.Address, error) {
	return _KIP113.Contract.AllNodeIds(&_KIP113.CallOpts, arg0)
}

// AllNodeIds is a free data retrieval call binding the contract method 0xa5834971.
//
// Solidity: function allNodeIds(uint256 ) view returns(address)
func (_KIP113 *KIP113CallerSession) AllNodeIds(arg0 *big.Int) (common.Address, error) {
	return _KIP113.Contract.AllNodeIds(&_KIP113.CallOpts, arg0)
}

// GetAllBlsInfo is a free data retrieval call binding the contract method 0x6968b53f.
//
// Solidity: function getAllBlsInfo() view returns(address[] nodeIdList, (bytes,bytes)[] pubkeyList)
func (_KIP113 *KIP113Caller) GetAllBlsInfo(opts *bind.CallOpts) (struct {
	NodeIdList []common.Address
	PubkeyList []IKIP113BlsPublicKeyInfo
}, error,
) {
	var out []interface{}
	err := _KIP113.contract.Call(opts, &out, "getAllBlsInfo")

	outstruct := new(struct {
		NodeIdList []common.Address
		PubkeyList []IKIP113BlsPublicKeyInfo
	})

	outstruct.NodeIdList = *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)
	outstruct.PubkeyList = *abi.ConvertType(out[1], new([]IKIP113BlsPublicKeyInfo)).(*[]IKIP113BlsPublicKeyInfo)
	return *outstruct, err
}

// GetAllBlsInfo is a free data retrieval call binding the contract method 0x6968b53f.
//
// Solidity: function getAllBlsInfo() view returns(address[] nodeIdList, (bytes,bytes)[] pubkeyList)
func (_KIP113 *KIP113Session) GetAllBlsInfo() (struct {
	NodeIdList []common.Address
	PubkeyList []IKIP113BlsPublicKeyInfo
}, error,
) {
	return _KIP113.Contract.GetAllBlsInfo(&_KIP113.CallOpts)
}

// GetAllBlsInfo is a free data retrieval call binding the contract method 0x6968b53f.
//
// Solidity: function getAllBlsInfo() view returns(address[] nodeIdList, (bytes,bytes)[] pubkeyList)
func (_KIP113 *KIP113CallerSession) GetAllBlsInfo() (struct {
	NodeIdList []common.Address
	PubkeyList []IKIP113BlsPublicKeyInfo
}, error,
) {
	return _KIP113.Contract.GetAllBlsInfo(&_KIP113.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_KIP113 *KIP113Caller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KIP113.contract.Call(opts, &out, "owner")
	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_KIP113 *KIP113Session) Owner() (common.Address, error) {
	return _KIP113.Contract.Owner(&_KIP113.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_KIP113 *KIP113CallerSession) Owner() (common.Address, error) {
	return _KIP113.Contract.Owner(&_KIP113.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_KIP113 *KIP113Caller) ProxiableUUID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _KIP113.contract.Call(opts, &out, "proxiableUUID")
	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_KIP113 *KIP113Session) ProxiableUUID() ([32]byte, error) {
	return _KIP113.Contract.ProxiableUUID(&_KIP113.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_KIP113 *KIP113CallerSession) ProxiableUUID() ([32]byte, error) {
	return _KIP113.Contract.ProxiableUUID(&_KIP113.CallOpts)
}

// Record is a free data retrieval call binding the contract method 0x3465d6d5.
//
// Solidity: function record(address ) view returns(bytes publicKey, bytes pop)
func (_KIP113 *KIP113Caller) Record(opts *bind.CallOpts, arg0 common.Address) (struct {
	PublicKey []byte
	Pop       []byte
}, error,
) {
	var out []interface{}
	err := _KIP113.contract.Call(opts, &out, "record", arg0)

	outstruct := new(struct {
		PublicKey []byte
		Pop       []byte
	})

	outstruct.PublicKey = *abi.ConvertType(out[0], new([]byte)).(*[]byte)
	outstruct.Pop = *abi.ConvertType(out[1], new([]byte)).(*[]byte)
	return *outstruct, err
}

// Record is a free data retrieval call binding the contract method 0x3465d6d5.
//
// Solidity: function record(address ) view returns(bytes publicKey, bytes pop)
func (_KIP113 *KIP113Session) Record(arg0 common.Address) (struct {
	PublicKey []byte
	Pop       []byte
}, error,
) {
	return _KIP113.Contract.Record(&_KIP113.CallOpts, arg0)
}

// Record is a free data retrieval call binding the contract method 0x3465d6d5.
//
// Solidity: function record(address ) view returns(bytes publicKey, bytes pop)
func (_KIP113 *KIP113CallerSession) Record(arg0 common.Address) (struct {
	PublicKey []byte
	Pop       []byte
}, error,
) {
	return _KIP113.Contract.Record(&_KIP113.CallOpts, arg0)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_KIP113 *KIP113Transactor) Initialize(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KIP113.contract.Transact(opts, "initialize")
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_KIP113 *KIP113Session) Initialize() (*types.Transaction, error) {
	return _KIP113.Contract.Initialize(&_KIP113.TransactOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_KIP113 *KIP113TransactorSession) Initialize() (*types.Transaction, error) {
	return _KIP113.Contract.Initialize(&_KIP113.TransactOpts)
}

// Register is a paid mutator transaction binding the contract method 0x786cd4d7.
//
// Solidity: function register(address cnNodeId, bytes publicKey, bytes pop) returns()
func (_KIP113 *KIP113Transactor) Register(opts *bind.TransactOpts, cnNodeId common.Address, publicKey []byte, pop []byte) (*types.Transaction, error) {
	return _KIP113.contract.Transact(opts, "register", cnNodeId, publicKey, pop)
}

// Register is a paid mutator transaction binding the contract method 0x786cd4d7.
//
// Solidity: function register(address cnNodeId, bytes publicKey, bytes pop) returns()
func (_KIP113 *KIP113Session) Register(cnNodeId common.Address, publicKey []byte, pop []byte) (*types.Transaction, error) {
	return _KIP113.Contract.Register(&_KIP113.TransactOpts, cnNodeId, publicKey, pop)
}

// Register is a paid mutator transaction binding the contract method 0x786cd4d7.
//
// Solidity: function register(address cnNodeId, bytes publicKey, bytes pop) returns()
func (_KIP113 *KIP113TransactorSession) Register(cnNodeId common.Address, publicKey []byte, pop []byte) (*types.Transaction, error) {
	return _KIP113.Contract.Register(&_KIP113.TransactOpts, cnNodeId, publicKey, pop)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_KIP113 *KIP113Transactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KIP113.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_KIP113 *KIP113Session) RenounceOwnership() (*types.Transaction, error) {
	return _KIP113.Contract.RenounceOwnership(&_KIP113.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_KIP113 *KIP113TransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _KIP113.Contract.RenounceOwnership(&_KIP113.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_KIP113 *KIP113Transactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _KIP113.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_KIP113 *KIP113Session) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _KIP113.Contract.TransferOwnership(&_KIP113.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_KIP113 *KIP113TransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _KIP113.Contract.TransferOwnership(&_KIP113.TransactOpts, newOwner)
}

// Unregister is a paid mutator transaction binding the contract method 0x2ec2c246.
//
// Solidity: function unregister(address cnNodeId) returns()
func (_KIP113 *KIP113Transactor) Unregister(opts *bind.TransactOpts, cnNodeId common.Address) (*types.Transaction, error) {
	return _KIP113.contract.Transact(opts, "unregister", cnNodeId)
}

// Unregister is a paid mutator transaction binding the contract method 0x2ec2c246.
//
// Solidity: function unregister(address cnNodeId) returns()
func (_KIP113 *KIP113Session) Unregister(cnNodeId common.Address) (*types.Transaction, error) {
	return _KIP113.Contract.Unregister(&_KIP113.TransactOpts, cnNodeId)
}

// Unregister is a paid mutator transaction binding the contract method 0x2ec2c246.
//
// Solidity: function unregister(address cnNodeId) returns()
func (_KIP113 *KIP113TransactorSession) Unregister(cnNodeId common.Address) (*types.Transaction, error) {
	return _KIP113.Contract.Unregister(&_KIP113.TransactOpts, cnNodeId)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_KIP113 *KIP113Transactor) UpgradeTo(opts *bind.TransactOpts, newImplementation common.Address) (*types.Transaction, error) {
	return _KIP113.contract.Transact(opts, "upgradeTo", newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_KIP113 *KIP113Session) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _KIP113.Contract.UpgradeTo(&_KIP113.TransactOpts, newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_KIP113 *KIP113TransactorSession) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _KIP113.Contract.UpgradeTo(&_KIP113.TransactOpts, newImplementation)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_KIP113 *KIP113Transactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _KIP113.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_KIP113 *KIP113Session) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _KIP113.Contract.UpgradeToAndCall(&_KIP113.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_KIP113 *KIP113TransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _KIP113.Contract.UpgradeToAndCall(&_KIP113.TransactOpts, newImplementation, data)
}

// KIP113AdminChangedIterator is returned from FilterAdminChanged and is used to iterate over the raw logs and unpacked data for AdminChanged events raised by the KIP113 contract.
type KIP113AdminChangedIterator struct {
	Event *KIP113AdminChanged // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log      // Log channel receiving the found contract events
	sub  klaytn.Subscription // Subscription for errors, completion and termination
	done bool                // Whether the subscription completed delivering logs
	fail error               // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *KIP113AdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KIP113AdminChanged)
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
		it.Event = new(KIP113AdminChanged)
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
func (it *KIP113AdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KIP113AdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KIP113AdminChanged represents a AdminChanged event raised by the KIP113 contract.
type KIP113AdminChanged struct {
	PreviousAdmin common.Address
	NewAdmin      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterAdminChanged is a free log retrieval operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_KIP113 *KIP113Filterer) FilterAdminChanged(opts *bind.FilterOpts) (*KIP113AdminChangedIterator, error) {
	logs, sub, err := _KIP113.contract.FilterLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return &KIP113AdminChangedIterator{contract: _KIP113.contract, event: "AdminChanged", logs: logs, sub: sub}, nil
}

// WatchAdminChanged is a free log subscription operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_KIP113 *KIP113Filterer) WatchAdminChanged(opts *bind.WatchOpts, sink chan<- *KIP113AdminChanged) (event.Subscription, error) {
	logs, sub, err := _KIP113.contract.WatchLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KIP113AdminChanged)
				if err := _KIP113.contract.UnpackLog(event, "AdminChanged", log); err != nil {
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

// ParseAdminChanged is a log parse operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_KIP113 *KIP113Filterer) ParseAdminChanged(log types.Log) (*KIP113AdminChanged, error) {
	event := new(KIP113AdminChanged)
	if err := _KIP113.contract.UnpackLog(event, "AdminChanged", log); err != nil {
		return nil, err
	}
	return event, nil
}

// KIP113BeaconUpgradedIterator is returned from FilterBeaconUpgraded and is used to iterate over the raw logs and unpacked data for BeaconUpgraded events raised by the KIP113 contract.
type KIP113BeaconUpgradedIterator struct {
	Event *KIP113BeaconUpgraded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log      // Log channel receiving the found contract events
	sub  klaytn.Subscription // Subscription for errors, completion and termination
	done bool                // Whether the subscription completed delivering logs
	fail error               // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *KIP113BeaconUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KIP113BeaconUpgraded)
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
		it.Event = new(KIP113BeaconUpgraded)
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
func (it *KIP113BeaconUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KIP113BeaconUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KIP113BeaconUpgraded represents a BeaconUpgraded event raised by the KIP113 contract.
type KIP113BeaconUpgraded struct {
	Beacon common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterBeaconUpgraded is a free log retrieval operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_KIP113 *KIP113Filterer) FilterBeaconUpgraded(opts *bind.FilterOpts, beacon []common.Address) (*KIP113BeaconUpgradedIterator, error) {
	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _KIP113.contract.FilterLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return &KIP113BeaconUpgradedIterator{contract: _KIP113.contract, event: "BeaconUpgraded", logs: logs, sub: sub}, nil
}

// WatchBeaconUpgraded is a free log subscription operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_KIP113 *KIP113Filterer) WatchBeaconUpgraded(opts *bind.WatchOpts, sink chan<- *KIP113BeaconUpgraded, beacon []common.Address) (event.Subscription, error) {
	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _KIP113.contract.WatchLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KIP113BeaconUpgraded)
				if err := _KIP113.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
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

// ParseBeaconUpgraded is a log parse operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_KIP113 *KIP113Filterer) ParseBeaconUpgraded(log types.Log) (*KIP113BeaconUpgraded, error) {
	event := new(KIP113BeaconUpgraded)
	if err := _KIP113.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
		return nil, err
	}
	return event, nil
}

// KIP113InitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the KIP113 contract.
type KIP113InitializedIterator struct {
	Event *KIP113Initialized // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log      // Log channel receiving the found contract events
	sub  klaytn.Subscription // Subscription for errors, completion and termination
	done bool                // Whether the subscription completed delivering logs
	fail error               // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *KIP113InitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KIP113Initialized)
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
		it.Event = new(KIP113Initialized)
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
func (it *KIP113InitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KIP113InitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KIP113Initialized represents a Initialized event raised by the KIP113 contract.
type KIP113Initialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_KIP113 *KIP113Filterer) FilterInitialized(opts *bind.FilterOpts) (*KIP113InitializedIterator, error) {
	logs, sub, err := _KIP113.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &KIP113InitializedIterator{contract: _KIP113.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_KIP113 *KIP113Filterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *KIP113Initialized) (event.Subscription, error) {
	logs, sub, err := _KIP113.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KIP113Initialized)
				if err := _KIP113.contract.UnpackLog(event, "Initialized", log); err != nil {
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

// ParseInitialized is a log parse operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_KIP113 *KIP113Filterer) ParseInitialized(log types.Log) (*KIP113Initialized, error) {
	event := new(KIP113Initialized)
	if err := _KIP113.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	return event, nil
}

// KIP113OwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the KIP113 contract.
type KIP113OwnershipTransferredIterator struct {
	Event *KIP113OwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log      // Log channel receiving the found contract events
	sub  klaytn.Subscription // Subscription for errors, completion and termination
	done bool                // Whether the subscription completed delivering logs
	fail error               // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *KIP113OwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KIP113OwnershipTransferred)
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
		it.Event = new(KIP113OwnershipTransferred)
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
func (it *KIP113OwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KIP113OwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KIP113OwnershipTransferred represents a OwnershipTransferred event raised by the KIP113 contract.
type KIP113OwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_KIP113 *KIP113Filterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*KIP113OwnershipTransferredIterator, error) {
	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _KIP113.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &KIP113OwnershipTransferredIterator{contract: _KIP113.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_KIP113 *KIP113Filterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *KIP113OwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {
	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _KIP113.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KIP113OwnershipTransferred)
				if err := _KIP113.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_KIP113 *KIP113Filterer) ParseOwnershipTransferred(log types.Log) (*KIP113OwnershipTransferred, error) {
	event := new(KIP113OwnershipTransferred)
	if err := _KIP113.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	return event, nil
}

// KIP113RegisteredIterator is returned from FilterRegistered and is used to iterate over the raw logs and unpacked data for Registered events raised by the KIP113 contract.
type KIP113RegisteredIterator struct {
	Event *KIP113Registered // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log      // Log channel receiving the found contract events
	sub  klaytn.Subscription // Subscription for errors, completion and termination
	done bool                // Whether the subscription completed delivering logs
	fail error               // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *KIP113RegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KIP113Registered)
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
		it.Event = new(KIP113Registered)
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
func (it *KIP113RegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KIP113RegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KIP113Registered represents a Registered event raised by the KIP113 contract.
type KIP113Registered struct {
	CnNodeId  common.Address
	PublicKey []byte
	Pop       []byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterRegistered is a free log retrieval operation binding the contract event 0x79c75399e89a1f580d9a6252cb8bdcf4cd80f73b3597c94d845eb52174ad930f.
//
// Solidity: event Registered(address cnNodeId, bytes publicKey, bytes pop)
func (_KIP113 *KIP113Filterer) FilterRegistered(opts *bind.FilterOpts) (*KIP113RegisteredIterator, error) {
	logs, sub, err := _KIP113.contract.FilterLogs(opts, "Registered")
	if err != nil {
		return nil, err
	}
	return &KIP113RegisteredIterator{contract: _KIP113.contract, event: "Registered", logs: logs, sub: sub}, nil
}

// WatchRegistered is a free log subscription operation binding the contract event 0x79c75399e89a1f580d9a6252cb8bdcf4cd80f73b3597c94d845eb52174ad930f.
//
// Solidity: event Registered(address cnNodeId, bytes publicKey, bytes pop)
func (_KIP113 *KIP113Filterer) WatchRegistered(opts *bind.WatchOpts, sink chan<- *KIP113Registered) (event.Subscription, error) {
	logs, sub, err := _KIP113.contract.WatchLogs(opts, "Registered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KIP113Registered)
				if err := _KIP113.contract.UnpackLog(event, "Registered", log); err != nil {
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

// ParseRegistered is a log parse operation binding the contract event 0x79c75399e89a1f580d9a6252cb8bdcf4cd80f73b3597c94d845eb52174ad930f.
//
// Solidity: event Registered(address cnNodeId, bytes publicKey, bytes pop)
func (_KIP113 *KIP113Filterer) ParseRegistered(log types.Log) (*KIP113Registered, error) {
	event := new(KIP113Registered)
	if err := _KIP113.contract.UnpackLog(event, "Registered", log); err != nil {
		return nil, err
	}
	return event, nil
}

// KIP113UnregisteredIterator is returned from FilterUnregistered and is used to iterate over the raw logs and unpacked data for Unregistered events raised by the KIP113 contract.
type KIP113UnregisteredIterator struct {
	Event *KIP113Unregistered // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log      // Log channel receiving the found contract events
	sub  klaytn.Subscription // Subscription for errors, completion and termination
	done bool                // Whether the subscription completed delivering logs
	fail error               // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *KIP113UnregisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KIP113Unregistered)
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
		it.Event = new(KIP113Unregistered)
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
func (it *KIP113UnregisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KIP113UnregisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KIP113Unregistered represents a Unregistered event raised by the KIP113 contract.
type KIP113Unregistered struct {
	CnNodeId  common.Address
	PublicKey []byte
	Pop       []byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterUnregistered is a free log retrieval operation binding the contract event 0xb98b07c4d52e8d65fa5416810f2746a810eb074b1ac7784e1b07e315c0dfd2d9.
//
// Solidity: event Unregistered(address cnNodeId, bytes publicKey, bytes pop)
func (_KIP113 *KIP113Filterer) FilterUnregistered(opts *bind.FilterOpts) (*KIP113UnregisteredIterator, error) {
	logs, sub, err := _KIP113.contract.FilterLogs(opts, "Unregistered")
	if err != nil {
		return nil, err
	}
	return &KIP113UnregisteredIterator{contract: _KIP113.contract, event: "Unregistered", logs: logs, sub: sub}, nil
}

// WatchUnregistered is a free log subscription operation binding the contract event 0xb98b07c4d52e8d65fa5416810f2746a810eb074b1ac7784e1b07e315c0dfd2d9.
//
// Solidity: event Unregistered(address cnNodeId, bytes publicKey, bytes pop)
func (_KIP113 *KIP113Filterer) WatchUnregistered(opts *bind.WatchOpts, sink chan<- *KIP113Unregistered) (event.Subscription, error) {
	logs, sub, err := _KIP113.contract.WatchLogs(opts, "Unregistered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KIP113Unregistered)
				if err := _KIP113.contract.UnpackLog(event, "Unregistered", log); err != nil {
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

// ParseUnregistered is a log parse operation binding the contract event 0xb98b07c4d52e8d65fa5416810f2746a810eb074b1ac7784e1b07e315c0dfd2d9.
//
// Solidity: event Unregistered(address cnNodeId, bytes publicKey, bytes pop)
func (_KIP113 *KIP113Filterer) ParseUnregistered(log types.Log) (*KIP113Unregistered, error) {
	event := new(KIP113Unregistered)
	if err := _KIP113.contract.UnpackLog(event, "Unregistered", log); err != nil {
		return nil, err
	}
	return event, nil
}

// KIP113UpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the KIP113 contract.
type KIP113UpgradedIterator struct {
	Event *KIP113Upgraded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log      // Log channel receiving the found contract events
	sub  klaytn.Subscription // Subscription for errors, completion and termination
	done bool                // Whether the subscription completed delivering logs
	fail error               // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *KIP113UpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KIP113Upgraded)
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
		it.Event = new(KIP113Upgraded)
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
func (it *KIP113UpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KIP113UpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KIP113Upgraded represents a Upgraded event raised by the KIP113 contract.
type KIP113Upgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_KIP113 *KIP113Filterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*KIP113UpgradedIterator, error) {
	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _KIP113.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &KIP113UpgradedIterator{contract: _KIP113.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_KIP113 *KIP113Filterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *KIP113Upgraded, implementation []common.Address) (event.Subscription, error) {
	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _KIP113.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KIP113Upgraded)
				if err := _KIP113.contract.UnpackLog(event, "Upgraded", log); err != nil {
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

// ParseUpgraded is a log parse operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_KIP113 *KIP113Filterer) ParseUpgraded(log types.Log) (*KIP113Upgraded, error) {
	event := new(KIP113Upgraded)
	if err := _KIP113.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	return event, nil
}

// KIP113MockMetaData contains all meta data concerning the KIP113Mock contract.
var KIP113MockMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"previousAdmin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"AdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"beacon\",\"type\":\"address\"}],\"name\":\"BeaconUpgraded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"cnNodeId\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"publicKey\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"pop\",\"type\":\"bytes\"}],\"name\":\"Registered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"cnNodeId\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"publicKey\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"pop\",\"type\":\"bytes\"}],\"name\":\"Unregistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"ZERO48HASH\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"ZERO96HASH\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"abook\",\"outputs\":[{\"internalType\":\"contractIAddressBook\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"allNodeIds\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllBlsInfo\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"nodeIdList\",\"type\":\"address[]\"},{\"components\":[{\"internalType\":\"bytes\",\"name\":\"publicKey\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"pop\",\"type\":\"bytes\"}],\"internalType\":\"structIKIP113.BlsPublicKeyInfo[]\",\"name\":\"pubkeyList\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"proxiableUUID\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"record\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"publicKey\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"pop\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"publicKey\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"pop\",\"type\":\"bytes\"}],\"name\":\"register\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"cnNodeId\",\"type\":\"address\"}],\"name\":\"unregister\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"}],\"name\":\"upgradeTo\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"upgradeToAndCall\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"6fc522c6": "ZERO48HASH()",
		"20abd458": "ZERO96HASH()",
		"829d639d": "abook()",
		"a5834971": "allNodeIds(uint256)",
		"6968b53f": "getAllBlsInfo()",
		"8129fc1c": "initialize()",
		"8da5cb5b": "owner()",
		"52d1902d": "proxiableUUID()",
		"3465d6d5": "record(address)",
		"786cd4d7": "register(address,bytes,bytes)",
		"715018a6": "renounceOwnership()",
		"f2fde38b": "transferOwnership(address)",
		"2ec2c246": "unregister(address)",
		"3659cfe6": "upgradeTo(address)",
		"4f1ef286": "upgradeToAndCall(address,bytes)",
	},
	Bin: "0x60a06040523060805234801561001457600080fd5b5061001d610022565b6100e1565b600054610100900460ff161561008e5760405162461bcd60e51b815260206004820152602760248201527f496e697469616c697a61626c653a20636f6e747261637420697320696e697469604482015266616c697a696e6760c81b606482015260840160405180910390fd5b60005460ff908116146100df576000805460ff191660ff9081179091556040519081527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15b565b608051611c0e61011860003960008181610593015281816105d301528181610672015281816106b201526107450152611c0e6000f3fe6080604052600436106100e85760003560e01c80636fc522c61161008a578063829d639d11610059578063829d639d1461026d5780638da5cb5b1461029b578063a5834971146102b9578063f2fde38b146102d957600080fd5b80636fc522c6146101ef578063715018a614610223578063786cd4d7146102385780638129fc1c1461025857600080fd5b80633659cfe6116100c65780633659cfe6146101845780634f1ef286146101a457806352d1902d146101b75780636968b53f146101cc57600080fd5b806320abd458146100ed5780632ec2c246146101345780633465d6d514610156575b600080fd5b3480156100f957600080fd5b506101217f46700b4d40ac5c35af2c22dda2787a91eb567b06c924a8fb8ae9a05b20c08c2181565b6040519081526020015b60405180910390f35b34801561014057600080fd5b5061015461014f36600461148a565b6102f9565b005b34801561016257600080fd5b5061017661017136600461148a565b61045d565b60405161012b9291906114fe565b34801561019057600080fd5b5061015461019f36600461148a565b610589565b6101546101b2366004611542565b610668565b3480156101c357600080fd5b50610121610738565b3480156101d857600080fd5b506101e16107eb565b60405161012b929190611606565b3480156101fb57600080fd5b506101217fc980e59163ce244bb4bb6211f48c7b46f88a4f40943e84eb99bdc41e129bd29381565b34801561022f57600080fd5b50610154610aa6565b34801561024457600080fd5b50610154610253366004611714565b610aba565b34801561026457600080fd5b50610154610bef565b34801561027957600080fd5b5061028361040081565b6040516001600160a01b03909116815260200161012b565b3480156102a757600080fd5b506097546001600160a01b0316610283565b3480156102c557600080fd5b506102836102d4366004611797565b610d07565b3480156102e557600080fd5b506101546102f436600461148a565b610d31565b610301610da7565b61030a81610e01565b1561035c5760405162461bcd60e51b815260206004820152601a60248201527f434e206973207374696c6c20696e2041646472657373426f6f6b00000000000060448201526064015b60405180910390fd5b6001600160a01b038116600090815260ca60205260409020805461037f906117b0565b90506000036103c75760405162461bcd60e51b815260206004820152601460248201527310d3881a5cc81b9bdd081c9959da5cdd195c995960621b6044820152606401610353565b6103d081610e7d565b6001600160a01b038116600090815260ca60205260409081902090517fb98b07c4d52e8d65fa5416810f2746a810eb074b1ac7784e1b07e315c0dfd2d99161041f918491906001820190611867565b60405180910390a16001600160a01b038116600090815260ca602052604081209061044a8282611427565b610458600183016000611427565b505050565b60ca60205260009081526040902080548190610478906117b0565b80601f01602080910402602001604051908101604052809291908181526020018280546104a4906117b0565b80156104f15780601f106104c6576101008083540402835291602001916104f1565b820191906000526020600020905b8154815290600101906020018083116104d457829003601f168201915b505050505090806001018054610506906117b0565b80601f0160208091040260200160405190810160405280929190818152602001828054610532906117b0565b801561057f5780601f106105545761010080835404028352916020019161057f565b820191906000526020600020905b81548152906001019060200180831161056257829003601f168201915b5050505050905082565b6001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001630036105d15760405162461bcd60e51b81526004016103539061189d565b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031661061a600080516020611b92833981519152546001600160a01b031690565b6001600160a01b0316146106405760405162461bcd60e51b8152600401610353906118e9565b61064981610f84565b6040805160008082526020820190925261066591839190610f8c565b50565b6001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001630036106b05760405162461bcd60e51b81526004016103539061189d565b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03166106f9600080516020611b92833981519152546001600160a01b031690565b6001600160a01b03161461071f5760405162461bcd60e51b8152600401610353906118e9565b61072882610f84565b61073482826001610f8c565b5050565b6000306001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016146107d85760405162461bcd60e51b815260206004820152603860248201527f555550535570677261646561626c653a206d757374206e6f742062652063616c60448201527f6c6564207468726f7567682064656c656761746563616c6c00000000000000006064820152608401610353565b50600080516020611b9283398151915290565b60c95460609081908067ffffffffffffffff81111561080c5761080c61152c565b604051908082528060200260200182016040528015610835578160200160208202803683370190505b5092508067ffffffffffffffff8111156108515761085161152c565b60405190808252806020026020018201604052801561089657816020015b604080518082019091526060808252602082015281526020019060019003908161086f5790505b50915060005b81811015610aa05760c981815481106108b7576108b7611935565b9060005260206000200160009054906101000a90046001600160a01b03168482815181106108e7576108e7611935565b60200260200101906001600160a01b031690816001600160a01b03168152505060ca600060c9838154811061091e5761091e611935565b60009182526020808320909101546001600160a01b031683528201929092526040908101909120815180830190925280548290829061095c906117b0565b80601f0160208091040260200160405190810160405280929190818152602001828054610988906117b0565b80156109d55780601f106109aa576101008083540402835291602001916109d5565b820191906000526020600020905b8154815290600101906020018083116109b857829003601f168201915b505050505081526020016001820180546109ee906117b0565b80601f0160208091040260200160405190810160405280929190818152602001828054610a1a906117b0565b8015610a675780601f10610a3c57610100808354040283529160200191610a67565b820191906000526020600020905b815481529060010190602001808311610a4a57829003601f168201915b505050505081525050838281518110610a8257610a82611935565b60200260200101819052508080610a9890611961565b91505061089c565b50509091565b610aae610da7565b610ab860006110f7565b565b6001600160a01b038516600090815260ca602052604090208054610add906117b0565b9050600003610b325760c980546001810182556000919091527f66be4f155c5ef2ebd3772b228f2f00681e4ed5826cdb3b1943cc11ad15ad1d280180546001600160a01b0319166001600160a01b0387161790555b6040805160606020601f87018190040282018101835291810185815290918291908790879081908501838280828437600092019190915250505090825250604080516020601f86018190048102820181019092528481529181019190859085908190840183828082843760009201829052509390945250506001600160a01b038816815260ca6020526040902082519091508190610bd090826119c8565b5060208201516001820190610be590826119c8565b5050505050505050565b600054610100900460ff1615808015610c0f5750600054600160ff909116105b80610c295750303b158015610c29575060005460ff166001145b610c8c5760405162461bcd60e51b815260206004820152602e60248201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160448201526d191e481a5b9a5d1a585b1a5e995960921b6064820152608401610353565b6000805460ff191660011790558015610caf576000805461ff0019166101001790555b610cb7611149565b610cbf611178565b8015610665576000805461ff0019169055604051600181527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a150565b60c98181548110610d1757600080fd5b6000918252602090912001546001600160a01b0316905081565b610d39610da7565b6001600160a01b038116610d9e5760405162461bcd60e51b815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201526564647265737360d01b6064820152608401610353565b610665816110f7565b6097546001600160a01b03163314610ab85760405162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e65726044820152606401610353565b604051630aabaead60e11b81526001600160a01b0382166004820152600090610400906315575d5a90602401606060405180830381865afa925050508015610e66575060408051601f3d908101601f19168201909252610e6391810190611a88565b60015b610e7257506000919050565b506001949350505050565b60005b60c95481101561073457816001600160a01b031660c98281548110610ea757610ea7611935565b6000918252602090912001546001600160a01b031603610f725760c98054610ed190600190611ad5565b81548110610ee157610ee1611935565b60009182526020909120015460c980546001600160a01b039092169183908110610f0d57610f0d611935565b9060005260206000200160006101000a8154816001600160a01b0302191690836001600160a01b0316021790555060c9805480610f4c57610f4c611ae8565b600082815260209020810160001990810180546001600160a01b03191690550190555050565b80610f7c81611961565b915050610e80565b610665610da7565b7f4910fdfa16fed3260ed0e7147f7cc6da11a60208b5b9406d12a635614ffd91435460ff1615610fbf576104588361119f565b826001600160a01b03166352d1902d6040518163ffffffff1660e01b8152600401602060405180830381865afa925050508015611019575060408051601f3d908101601f1916820190925261101691810190611afe565b60015b61107c5760405162461bcd60e51b815260206004820152602e60248201527f45524331393637557067726164653a206e657720696d706c656d656e7461746960448201526d6f6e206973206e6f74205555505360901b6064820152608401610353565b600080516020611b9283398151915281146110eb5760405162461bcd60e51b815260206004820152602960248201527f45524331393637557067726164653a20756e737570706f727465642070726f786044820152681a58589b195555525160ba1b6064820152608401610353565b5061045883838361123b565b609780546001600160a01b038381166001600160a01b0319831681179093556040519116919082907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a35050565b600054610100900460ff166111705760405162461bcd60e51b815260040161035390611b17565b610ab8611266565b600054610100900460ff16610ab85760405162461bcd60e51b815260040161035390611b17565b6001600160a01b0381163b61120c5760405162461bcd60e51b815260206004820152602d60248201527f455243313936373a206e657720696d706c656d656e746174696f6e206973206e60448201526c1bdd08184818dbdb9d1c9858dd609a1b6064820152608401610353565b600080516020611b9283398151915280546001600160a01b0319166001600160a01b0392909216919091179055565b61124483611296565b6000825111806112515750805b156104585761126083836112d6565b50505050565b600054610100900460ff1661128d5760405162461bcd60e51b815260040161035390611b17565b610ab8336110f7565b61129f8161119f565b6040516001600160a01b038216907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b90600090a250565b60606112fb8383604051806060016040528060278152602001611bb260279139611304565b90505b92915050565b6060600080856001600160a01b0316856040516113219190611b62565b600060405180830381855af49150503d806000811461135c576040519150601f19603f3d011682016040523d82523d6000602084013e611361565b606091505b50915091506113728683838761137c565b9695505050505050565b606083156113eb5782516000036113e4576001600160a01b0385163b6113e45760405162461bcd60e51b815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e74726163740000006044820152606401610353565b50816113f5565b6113f583836113fd565b949350505050565b81511561140d5781518083602001fd5b8060405162461bcd60e51b81526004016103539190611b7e565b508054611433906117b0565b6000825580601f10611443575050565b601f01602090049060005260206000209081019061066591905b80821115611471576000815560010161145d565b5090565b6001600160a01b038116811461066557600080fd5b60006020828403121561149c57600080fd5b81356114a781611475565b9392505050565b60005b838110156114c95781810151838201526020016114b1565b50506000910152565b600081518084526114ea8160208601602086016114ae565b601f01601f19169290920160200192915050565b60408152600061151160408301856114d2565b828103602084015261152381856114d2565b95945050505050565b634e487b7160e01b600052604160045260246000fd5b6000806040838503121561155557600080fd5b823561156081611475565b9150602083013567ffffffffffffffff8082111561157d57600080fd5b818501915085601f83011261159157600080fd5b8135818111156115a3576115a361152c565b604051601f8201601f19908116603f011681019083821181831017156115cb576115cb61152c565b816040528281528860208487010111156115e457600080fd5b8260208601602083013760006020848301015280955050505050509250929050565b60408082528351828201819052600091906020906060850190828801855b828110156116495781516001600160a01b031684529284019290840190600101611624565b50505084810382860152855180825282820190600581901b8301840188850160005b838110156116bb57858303601f190185528151805189855261168f8a8601826114d2565b91890151858303868b01529190506116a781836114d2565b96890196945050509086019060010161166b565b50909a9950505050505050505050565b60008083601f8401126116dd57600080fd5b50813567ffffffffffffffff8111156116f557600080fd5b60208301915083602082850101111561170d57600080fd5b9250929050565b60008060008060006060868803121561172c57600080fd5b853561173781611475565b9450602086013567ffffffffffffffff8082111561175457600080fd5b61176089838a016116cb565b9096509450604088013591508082111561177957600080fd5b50611786888289016116cb565b969995985093965092949392505050565b6000602082840312156117a957600080fd5b5035919050565b600181811c908216806117c457607f821691505b6020821081036117e457634e487b7160e01b600052602260045260246000fd5b50919050565b600081546117f7816117b0565b808552602060018381168015611814576001811461182e5761185c565b60ff1985168884015283151560051b88018301955061185c565b866000528260002060005b858110156118545781548a8201860152908301908401611839565b890184019650505b505050505092915050565b6001600160a01b038416815260606020820181905260009061188b908301856117ea565b828103604084015261137281856117ea565b6020808252602c908201527f46756e6374696f6e206d7573742062652063616c6c6564207468726f7567682060408201526b19195b1959d85d1958d85b1b60a21b606082015260800190565b6020808252602c908201527f46756e6374696f6e206d7573742062652063616c6c6564207468726f7567682060408201526b6163746976652070726f787960a01b606082015260800190565b634e487b7160e01b600052603260045260246000fd5b634e487b7160e01b600052601160045260246000fd5b6000600182016119735761197361194b565b5060010190565b601f82111561045857600081815260208120601f850160051c810160208610156119a15750805b601f850160051c820191505b818110156119c0578281556001016119ad565b505050505050565b815167ffffffffffffffff8111156119e2576119e261152c565b6119f6816119f084546117b0565b8461197a565b602080601f831160018114611a2b5760008415611a135750858301515b600019600386901b1c1916600185901b1785556119c0565b600085815260208120601f198616915b82811015611a5a57888601518255948401946001909101908401611a3b565b5085821015611a785787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b600080600060608486031215611a9d57600080fd5b8351611aa881611475565b6020850151909350611ab981611475565b6040850151909250611aca81611475565b809150509250925092565b818103818111156112fe576112fe61194b565b634e487b7160e01b600052603160045260246000fd5b600060208284031215611b1057600080fd5b5051919050565b6020808252602b908201527f496e697469616c697a61626c653a20636f6e7472616374206973206e6f74206960408201526a6e697469616c697a696e6760a81b606082015260800190565b60008251611b748184602087016114ae565b9190910192915050565b6020815260006112fb60208301846114d256fe360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc416464726573733a206c6f772d6c6576656c2064656c65676174652063616c6c206661696c6564a264697066735822122073ea4af54d3ffdbc93826462adf4b55a449fceb20536bcce5ec19edfa2b5128c64736f6c63430008130033",
}

// KIP113MockABI is the input ABI used to generate the binding from.
// Deprecated: Use KIP113MockMetaData.ABI instead.
var KIP113MockABI = KIP113MockMetaData.ABI

// KIP113MockBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const KIP113MockBinRuntime = `6080604052600436106100e85760003560e01c80636fc522c61161008a578063829d639d11610059578063829d639d1461026d5780638da5cb5b1461029b578063a5834971146102b9578063f2fde38b146102d957600080fd5b80636fc522c6146101ef578063715018a614610223578063786cd4d7146102385780638129fc1c1461025857600080fd5b80633659cfe6116100c65780633659cfe6146101845780634f1ef286146101a457806352d1902d146101b75780636968b53f146101cc57600080fd5b806320abd458146100ed5780632ec2c246146101345780633465d6d514610156575b600080fd5b3480156100f957600080fd5b506101217f46700b4d40ac5c35af2c22dda2787a91eb567b06c924a8fb8ae9a05b20c08c2181565b6040519081526020015b60405180910390f35b34801561014057600080fd5b5061015461014f36600461148a565b6102f9565b005b34801561016257600080fd5b5061017661017136600461148a565b61045d565b60405161012b9291906114fe565b34801561019057600080fd5b5061015461019f36600461148a565b610589565b6101546101b2366004611542565b610668565b3480156101c357600080fd5b50610121610738565b3480156101d857600080fd5b506101e16107eb565b60405161012b929190611606565b3480156101fb57600080fd5b506101217fc980e59163ce244bb4bb6211f48c7b46f88a4f40943e84eb99bdc41e129bd29381565b34801561022f57600080fd5b50610154610aa6565b34801561024457600080fd5b50610154610253366004611714565b610aba565b34801561026457600080fd5b50610154610bef565b34801561027957600080fd5b5061028361040081565b6040516001600160a01b03909116815260200161012b565b3480156102a757600080fd5b506097546001600160a01b0316610283565b3480156102c557600080fd5b506102836102d4366004611797565b610d07565b3480156102e557600080fd5b506101546102f436600461148a565b610d31565b610301610da7565b61030a81610e01565b1561035c5760405162461bcd60e51b815260206004820152601a60248201527f434e206973207374696c6c20696e2041646472657373426f6f6b00000000000060448201526064015b60405180910390fd5b6001600160a01b038116600090815260ca60205260409020805461037f906117b0565b90506000036103c75760405162461bcd60e51b815260206004820152601460248201527310d3881a5cc81b9bdd081c9959da5cdd195c995960621b6044820152606401610353565b6103d081610e7d565b6001600160a01b038116600090815260ca60205260409081902090517fb98b07c4d52e8d65fa5416810f2746a810eb074b1ac7784e1b07e315c0dfd2d99161041f918491906001820190611867565b60405180910390a16001600160a01b038116600090815260ca602052604081209061044a8282611427565b610458600183016000611427565b505050565b60ca60205260009081526040902080548190610478906117b0565b80601f01602080910402602001604051908101604052809291908181526020018280546104a4906117b0565b80156104f15780601f106104c6576101008083540402835291602001916104f1565b820191906000526020600020905b8154815290600101906020018083116104d457829003601f168201915b505050505090806001018054610506906117b0565b80601f0160208091040260200160405190810160405280929190818152602001828054610532906117b0565b801561057f5780601f106105545761010080835404028352916020019161057f565b820191906000526020600020905b81548152906001019060200180831161056257829003601f168201915b5050505050905082565b6001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001630036105d15760405162461bcd60e51b81526004016103539061189d565b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031661061a600080516020611b92833981519152546001600160a01b031690565b6001600160a01b0316146106405760405162461bcd60e51b8152600401610353906118e9565b61064981610f84565b6040805160008082526020820190925261066591839190610f8c565b50565b6001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001630036106b05760405162461bcd60e51b81526004016103539061189d565b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03166106f9600080516020611b92833981519152546001600160a01b031690565b6001600160a01b03161461071f5760405162461bcd60e51b8152600401610353906118e9565b61072882610f84565b61073482826001610f8c565b5050565b6000306001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016146107d85760405162461bcd60e51b815260206004820152603860248201527f555550535570677261646561626c653a206d757374206e6f742062652063616c60448201527f6c6564207468726f7567682064656c656761746563616c6c00000000000000006064820152608401610353565b50600080516020611b9283398151915290565b60c95460609081908067ffffffffffffffff81111561080c5761080c61152c565b604051908082528060200260200182016040528015610835578160200160208202803683370190505b5092508067ffffffffffffffff8111156108515761085161152c565b60405190808252806020026020018201604052801561089657816020015b604080518082019091526060808252602082015281526020019060019003908161086f5790505b50915060005b81811015610aa05760c981815481106108b7576108b7611935565b9060005260206000200160009054906101000a90046001600160a01b03168482815181106108e7576108e7611935565b60200260200101906001600160a01b031690816001600160a01b03168152505060ca600060c9838154811061091e5761091e611935565b60009182526020808320909101546001600160a01b031683528201929092526040908101909120815180830190925280548290829061095c906117b0565b80601f0160208091040260200160405190810160405280929190818152602001828054610988906117b0565b80156109d55780601f106109aa576101008083540402835291602001916109d5565b820191906000526020600020905b8154815290600101906020018083116109b857829003601f168201915b505050505081526020016001820180546109ee906117b0565b80601f0160208091040260200160405190810160405280929190818152602001828054610a1a906117b0565b8015610a675780601f10610a3c57610100808354040283529160200191610a67565b820191906000526020600020905b815481529060010190602001808311610a4a57829003601f168201915b505050505081525050838281518110610a8257610a82611935565b60200260200101819052508080610a9890611961565b91505061089c565b50509091565b610aae610da7565b610ab860006110f7565b565b6001600160a01b038516600090815260ca602052604090208054610add906117b0565b9050600003610b325760c980546001810182556000919091527f66be4f155c5ef2ebd3772b228f2f00681e4ed5826cdb3b1943cc11ad15ad1d280180546001600160a01b0319166001600160a01b0387161790555b6040805160606020601f87018190040282018101835291810185815290918291908790879081908501838280828437600092019190915250505090825250604080516020601f86018190048102820181019092528481529181019190859085908190840183828082843760009201829052509390945250506001600160a01b038816815260ca6020526040902082519091508190610bd090826119c8565b5060208201516001820190610be590826119c8565b5050505050505050565b600054610100900460ff1615808015610c0f5750600054600160ff909116105b80610c295750303b158015610c29575060005460ff166001145b610c8c5760405162461bcd60e51b815260206004820152602e60248201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160448201526d191e481a5b9a5d1a585b1a5e995960921b6064820152608401610353565b6000805460ff191660011790558015610caf576000805461ff0019166101001790555b610cb7611149565b610cbf611178565b8015610665576000805461ff0019169055604051600181527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a150565b60c98181548110610d1757600080fd5b6000918252602090912001546001600160a01b0316905081565b610d39610da7565b6001600160a01b038116610d9e5760405162461bcd60e51b815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201526564647265737360d01b6064820152608401610353565b610665816110f7565b6097546001600160a01b03163314610ab85760405162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e65726044820152606401610353565b604051630aabaead60e11b81526001600160a01b0382166004820152600090610400906315575d5a90602401606060405180830381865afa925050508015610e66575060408051601f3d908101601f19168201909252610e6391810190611a88565b60015b610e7257506000919050565b506001949350505050565b60005b60c95481101561073457816001600160a01b031660c98281548110610ea757610ea7611935565b6000918252602090912001546001600160a01b031603610f725760c98054610ed190600190611ad5565b81548110610ee157610ee1611935565b60009182526020909120015460c980546001600160a01b039092169183908110610f0d57610f0d611935565b9060005260206000200160006101000a8154816001600160a01b0302191690836001600160a01b0316021790555060c9805480610f4c57610f4c611ae8565b600082815260209020810160001990810180546001600160a01b03191690550190555050565b80610f7c81611961565b915050610e80565b610665610da7565b7f4910fdfa16fed3260ed0e7147f7cc6da11a60208b5b9406d12a635614ffd91435460ff1615610fbf576104588361119f565b826001600160a01b03166352d1902d6040518163ffffffff1660e01b8152600401602060405180830381865afa925050508015611019575060408051601f3d908101601f1916820190925261101691810190611afe565b60015b61107c5760405162461bcd60e51b815260206004820152602e60248201527f45524331393637557067726164653a206e657720696d706c656d656e7461746960448201526d6f6e206973206e6f74205555505360901b6064820152608401610353565b600080516020611b9283398151915281146110eb5760405162461bcd60e51b815260206004820152602960248201527f45524331393637557067726164653a20756e737570706f727465642070726f786044820152681a58589b195555525160ba1b6064820152608401610353565b5061045883838361123b565b609780546001600160a01b038381166001600160a01b0319831681179093556040519116919082907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a35050565b600054610100900460ff166111705760405162461bcd60e51b815260040161035390611b17565b610ab8611266565b600054610100900460ff16610ab85760405162461bcd60e51b815260040161035390611b17565b6001600160a01b0381163b61120c5760405162461bcd60e51b815260206004820152602d60248201527f455243313936373a206e657720696d706c656d656e746174696f6e206973206e60448201526c1bdd08184818dbdb9d1c9858dd609a1b6064820152608401610353565b600080516020611b9283398151915280546001600160a01b0319166001600160a01b0392909216919091179055565b61124483611296565b6000825111806112515750805b156104585761126083836112d6565b50505050565b600054610100900460ff1661128d5760405162461bcd60e51b815260040161035390611b17565b610ab8336110f7565b61129f8161119f565b6040516001600160a01b038216907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b90600090a250565b60606112fb8383604051806060016040528060278152602001611bb260279139611304565b90505b92915050565b6060600080856001600160a01b0316856040516113219190611b62565b600060405180830381855af49150503d806000811461135c576040519150601f19603f3d011682016040523d82523d6000602084013e611361565b606091505b50915091506113728683838761137c565b9695505050505050565b606083156113eb5782516000036113e4576001600160a01b0385163b6113e45760405162461bcd60e51b815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e74726163740000006044820152606401610353565b50816113f5565b6113f583836113fd565b949350505050565b81511561140d5781518083602001fd5b8060405162461bcd60e51b81526004016103539190611b7e565b508054611433906117b0565b6000825580601f10611443575050565b601f01602090049060005260206000209081019061066591905b80821115611471576000815560010161145d565b5090565b6001600160a01b038116811461066557600080fd5b60006020828403121561149c57600080fd5b81356114a781611475565b9392505050565b60005b838110156114c95781810151838201526020016114b1565b50506000910152565b600081518084526114ea8160208601602086016114ae565b601f01601f19169290920160200192915050565b60408152600061151160408301856114d2565b828103602084015261152381856114d2565b95945050505050565b634e487b7160e01b600052604160045260246000fd5b6000806040838503121561155557600080fd5b823561156081611475565b9150602083013567ffffffffffffffff8082111561157d57600080fd5b818501915085601f83011261159157600080fd5b8135818111156115a3576115a361152c565b604051601f8201601f19908116603f011681019083821181831017156115cb576115cb61152c565b816040528281528860208487010111156115e457600080fd5b8260208601602083013760006020848301015280955050505050509250929050565b60408082528351828201819052600091906020906060850190828801855b828110156116495781516001600160a01b031684529284019290840190600101611624565b50505084810382860152855180825282820190600581901b8301840188850160005b838110156116bb57858303601f190185528151805189855261168f8a8601826114d2565b91890151858303868b01529190506116a781836114d2565b96890196945050509086019060010161166b565b50909a9950505050505050505050565b60008083601f8401126116dd57600080fd5b50813567ffffffffffffffff8111156116f557600080fd5b60208301915083602082850101111561170d57600080fd5b9250929050565b60008060008060006060868803121561172c57600080fd5b853561173781611475565b9450602086013567ffffffffffffffff8082111561175457600080fd5b61176089838a016116cb565b9096509450604088013591508082111561177957600080fd5b50611786888289016116cb565b969995985093965092949392505050565b6000602082840312156117a957600080fd5b5035919050565b600181811c908216806117c457607f821691505b6020821081036117e457634e487b7160e01b600052602260045260246000fd5b50919050565b600081546117f7816117b0565b808552602060018381168015611814576001811461182e5761185c565b60ff1985168884015283151560051b88018301955061185c565b866000528260002060005b858110156118545781548a8201860152908301908401611839565b890184019650505b505050505092915050565b6001600160a01b038416815260606020820181905260009061188b908301856117ea565b828103604084015261137281856117ea565b6020808252602c908201527f46756e6374696f6e206d7573742062652063616c6c6564207468726f7567682060408201526b19195b1959d85d1958d85b1b60a21b606082015260800190565b6020808252602c908201527f46756e6374696f6e206d7573742062652063616c6c6564207468726f7567682060408201526b6163746976652070726f787960a01b606082015260800190565b634e487b7160e01b600052603260045260246000fd5b634e487b7160e01b600052601160045260246000fd5b6000600182016119735761197361194b565b5060010190565b601f82111561045857600081815260208120601f850160051c810160208610156119a15750805b601f850160051c820191505b818110156119c0578281556001016119ad565b505050505050565b815167ffffffffffffffff8111156119e2576119e261152c565b6119f6816119f084546117b0565b8461197a565b602080601f831160018114611a2b5760008415611a135750858301515b600019600386901b1c1916600185901b1785556119c0565b600085815260208120601f198616915b82811015611a5a57888601518255948401946001909101908401611a3b565b5085821015611a785787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b600080600060608486031215611a9d57600080fd5b8351611aa881611475565b6020850151909350611ab981611475565b6040850151909250611aca81611475565b809150509250925092565b818103818111156112fe576112fe61194b565b634e487b7160e01b600052603160045260246000fd5b600060208284031215611b1057600080fd5b5051919050565b6020808252602b908201527f496e697469616c697a61626c653a20636f6e7472616374206973206e6f74206960408201526a6e697469616c697a696e6760a81b606082015260800190565b60008251611b748184602087016114ae565b9190910192915050565b6020815260006112fb60208301846114d256fe360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc416464726573733a206c6f772d6c6576656c2064656c65676174652063616c6c206661696c6564a264697066735822122073ea4af54d3ffdbc93826462adf4b55a449fceb20536bcce5ec19edfa2b5128c64736f6c63430008130033`

// KIP113MockFuncSigs maps the 4-byte function signature to its string representation.
// Deprecated: Use KIP113MockMetaData.Sigs instead.
var KIP113MockFuncSigs = KIP113MockMetaData.Sigs

// KIP113MockBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use KIP113MockMetaData.Bin instead.
var KIP113MockBin = KIP113MockMetaData.Bin

// DeployKIP113Mock deploys a new Klaytn contract, binding an instance of KIP113Mock to it.
func DeployKIP113Mock(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *KIP113Mock, error) {
	parsed, err := KIP113MockMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(KIP113MockBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &KIP113Mock{KIP113MockCaller: KIP113MockCaller{contract: contract}, KIP113MockTransactor: KIP113MockTransactor{contract: contract}, KIP113MockFilterer: KIP113MockFilterer{contract: contract}}, nil
}

// KIP113Mock is an auto generated Go binding around a Klaytn contract.
type KIP113Mock struct {
	KIP113MockCaller     // Read-only binding to the contract
	KIP113MockTransactor // Write-only binding to the contract
	KIP113MockFilterer   // Log filterer for contract events
}

// KIP113MockCaller is an auto generated read-only Go binding around a Klaytn contract.
type KIP113MockCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// KIP113MockTransactor is an auto generated write-only Go binding around a Klaytn contract.
type KIP113MockTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// KIP113MockFilterer is an auto generated log filtering Go binding around a Klaytn contract events.
type KIP113MockFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// KIP113MockSession is an auto generated Go binding around a Klaytn contract,
// with pre-set call and transact options.
type KIP113MockSession struct {
	Contract     *KIP113Mock       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// KIP113MockCallerSession is an auto generated read-only Go binding around a Klaytn contract,
// with pre-set call options.
type KIP113MockCallerSession struct {
	Contract *KIP113MockCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// KIP113MockTransactorSession is an auto generated write-only Go binding around a Klaytn contract,
// with pre-set transact options.
type KIP113MockTransactorSession struct {
	Contract     *KIP113MockTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// KIP113MockRaw is an auto generated low-level Go binding around a Klaytn contract.
type KIP113MockRaw struct {
	Contract *KIP113Mock // Generic contract binding to access the raw methods on
}

// KIP113MockCallerRaw is an auto generated low-level read-only Go binding around a Klaytn contract.
type KIP113MockCallerRaw struct {
	Contract *KIP113MockCaller // Generic read-only contract binding to access the raw methods on
}

// KIP113MockTransactorRaw is an auto generated low-level write-only Go binding around a Klaytn contract.
type KIP113MockTransactorRaw struct {
	Contract *KIP113MockTransactor // Generic write-only contract binding to access the raw methods on
}

// NewKIP113Mock creates a new instance of KIP113Mock, bound to a specific deployed contract.
func NewKIP113Mock(address common.Address, backend bind.ContractBackend) (*KIP113Mock, error) {
	contract, err := bindKIP113Mock(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &KIP113Mock{KIP113MockCaller: KIP113MockCaller{contract: contract}, KIP113MockTransactor: KIP113MockTransactor{contract: contract}, KIP113MockFilterer: KIP113MockFilterer{contract: contract}}, nil
}

// NewKIP113MockCaller creates a new read-only instance of KIP113Mock, bound to a specific deployed contract.
func NewKIP113MockCaller(address common.Address, caller bind.ContractCaller) (*KIP113MockCaller, error) {
	contract, err := bindKIP113Mock(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &KIP113MockCaller{contract: contract}, nil
}

// NewKIP113MockTransactor creates a new write-only instance of KIP113Mock, bound to a specific deployed contract.
func NewKIP113MockTransactor(address common.Address, transactor bind.ContractTransactor) (*KIP113MockTransactor, error) {
	contract, err := bindKIP113Mock(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &KIP113MockTransactor{contract: contract}, nil
}

// NewKIP113MockFilterer creates a new log filterer instance of KIP113Mock, bound to a specific deployed contract.
func NewKIP113MockFilterer(address common.Address, filterer bind.ContractFilterer) (*KIP113MockFilterer, error) {
	contract, err := bindKIP113Mock(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &KIP113MockFilterer{contract: contract}, nil
}

// bindKIP113Mock binds a generic wrapper to an already deployed contract.
func bindKIP113Mock(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := KIP113MockMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_KIP113Mock *KIP113MockRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _KIP113Mock.Contract.KIP113MockCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_KIP113Mock *KIP113MockRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KIP113Mock.Contract.KIP113MockTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_KIP113Mock *KIP113MockRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KIP113Mock.Contract.KIP113MockTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_KIP113Mock *KIP113MockCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _KIP113Mock.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_KIP113Mock *KIP113MockTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KIP113Mock.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_KIP113Mock *KIP113MockTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KIP113Mock.Contract.contract.Transact(opts, method, params...)
}

// ZERO48HASH is a free data retrieval call binding the contract method 0x6fc522c6.
//
// Solidity: function ZERO48HASH() view returns(bytes32)
func (_KIP113Mock *KIP113MockCaller) ZERO48HASH(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _KIP113Mock.contract.Call(opts, &out, "ZERO48HASH")
	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err
}

// ZERO48HASH is a free data retrieval call binding the contract method 0x6fc522c6.
//
// Solidity: function ZERO48HASH() view returns(bytes32)
func (_KIP113Mock *KIP113MockSession) ZERO48HASH() ([32]byte, error) {
	return _KIP113Mock.Contract.ZERO48HASH(&_KIP113Mock.CallOpts)
}

// ZERO48HASH is a free data retrieval call binding the contract method 0x6fc522c6.
//
// Solidity: function ZERO48HASH() view returns(bytes32)
func (_KIP113Mock *KIP113MockCallerSession) ZERO48HASH() ([32]byte, error) {
	return _KIP113Mock.Contract.ZERO48HASH(&_KIP113Mock.CallOpts)
}

// ZERO96HASH is a free data retrieval call binding the contract method 0x20abd458.
//
// Solidity: function ZERO96HASH() view returns(bytes32)
func (_KIP113Mock *KIP113MockCaller) ZERO96HASH(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _KIP113Mock.contract.Call(opts, &out, "ZERO96HASH")
	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err
}

// ZERO96HASH is a free data retrieval call binding the contract method 0x20abd458.
//
// Solidity: function ZERO96HASH() view returns(bytes32)
func (_KIP113Mock *KIP113MockSession) ZERO96HASH() ([32]byte, error) {
	return _KIP113Mock.Contract.ZERO96HASH(&_KIP113Mock.CallOpts)
}

// ZERO96HASH is a free data retrieval call binding the contract method 0x20abd458.
//
// Solidity: function ZERO96HASH() view returns(bytes32)
func (_KIP113Mock *KIP113MockCallerSession) ZERO96HASH() ([32]byte, error) {
	return _KIP113Mock.Contract.ZERO96HASH(&_KIP113Mock.CallOpts)
}

// Abook is a free data retrieval call binding the contract method 0x829d639d.
//
// Solidity: function abook() view returns(address)
func (_KIP113Mock *KIP113MockCaller) Abook(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KIP113Mock.contract.Call(opts, &out, "abook")
	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err
}

// Abook is a free data retrieval call binding the contract method 0x829d639d.
//
// Solidity: function abook() view returns(address)
func (_KIP113Mock *KIP113MockSession) Abook() (common.Address, error) {
	return _KIP113Mock.Contract.Abook(&_KIP113Mock.CallOpts)
}

// Abook is a free data retrieval call binding the contract method 0x829d639d.
//
// Solidity: function abook() view returns(address)
func (_KIP113Mock *KIP113MockCallerSession) Abook() (common.Address, error) {
	return _KIP113Mock.Contract.Abook(&_KIP113Mock.CallOpts)
}

// AllNodeIds is a free data retrieval call binding the contract method 0xa5834971.
//
// Solidity: function allNodeIds(uint256 ) view returns(address)
func (_KIP113Mock *KIP113MockCaller) AllNodeIds(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var out []interface{}
	err := _KIP113Mock.contract.Call(opts, &out, "allNodeIds", arg0)
	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err
}

// AllNodeIds is a free data retrieval call binding the contract method 0xa5834971.
//
// Solidity: function allNodeIds(uint256 ) view returns(address)
func (_KIP113Mock *KIP113MockSession) AllNodeIds(arg0 *big.Int) (common.Address, error) {
	return _KIP113Mock.Contract.AllNodeIds(&_KIP113Mock.CallOpts, arg0)
}

// AllNodeIds is a free data retrieval call binding the contract method 0xa5834971.
//
// Solidity: function allNodeIds(uint256 ) view returns(address)
func (_KIP113Mock *KIP113MockCallerSession) AllNodeIds(arg0 *big.Int) (common.Address, error) {
	return _KIP113Mock.Contract.AllNodeIds(&_KIP113Mock.CallOpts, arg0)
}

// GetAllBlsInfo is a free data retrieval call binding the contract method 0x6968b53f.
//
// Solidity: function getAllBlsInfo() view returns(address[] nodeIdList, (bytes,bytes)[] pubkeyList)
func (_KIP113Mock *KIP113MockCaller) GetAllBlsInfo(opts *bind.CallOpts) (struct {
	NodeIdList []common.Address
	PubkeyList []IKIP113BlsPublicKeyInfo
}, error,
) {
	var out []interface{}
	err := _KIP113Mock.contract.Call(opts, &out, "getAllBlsInfo")

	outstruct := new(struct {
		NodeIdList []common.Address
		PubkeyList []IKIP113BlsPublicKeyInfo
	})

	outstruct.NodeIdList = *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)
	outstruct.PubkeyList = *abi.ConvertType(out[1], new([]IKIP113BlsPublicKeyInfo)).(*[]IKIP113BlsPublicKeyInfo)
	return *outstruct, err
}

// GetAllBlsInfo is a free data retrieval call binding the contract method 0x6968b53f.
//
// Solidity: function getAllBlsInfo() view returns(address[] nodeIdList, (bytes,bytes)[] pubkeyList)
func (_KIP113Mock *KIP113MockSession) GetAllBlsInfo() (struct {
	NodeIdList []common.Address
	PubkeyList []IKIP113BlsPublicKeyInfo
}, error,
) {
	return _KIP113Mock.Contract.GetAllBlsInfo(&_KIP113Mock.CallOpts)
}

// GetAllBlsInfo is a free data retrieval call binding the contract method 0x6968b53f.
//
// Solidity: function getAllBlsInfo() view returns(address[] nodeIdList, (bytes,bytes)[] pubkeyList)
func (_KIP113Mock *KIP113MockCallerSession) GetAllBlsInfo() (struct {
	NodeIdList []common.Address
	PubkeyList []IKIP113BlsPublicKeyInfo
}, error,
) {
	return _KIP113Mock.Contract.GetAllBlsInfo(&_KIP113Mock.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_KIP113Mock *KIP113MockCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KIP113Mock.contract.Call(opts, &out, "owner")
	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_KIP113Mock *KIP113MockSession) Owner() (common.Address, error) {
	return _KIP113Mock.Contract.Owner(&_KIP113Mock.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_KIP113Mock *KIP113MockCallerSession) Owner() (common.Address, error) {
	return _KIP113Mock.Contract.Owner(&_KIP113Mock.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_KIP113Mock *KIP113MockCaller) ProxiableUUID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _KIP113Mock.contract.Call(opts, &out, "proxiableUUID")
	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_KIP113Mock *KIP113MockSession) ProxiableUUID() ([32]byte, error) {
	return _KIP113Mock.Contract.ProxiableUUID(&_KIP113Mock.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_KIP113Mock *KIP113MockCallerSession) ProxiableUUID() ([32]byte, error) {
	return _KIP113Mock.Contract.ProxiableUUID(&_KIP113Mock.CallOpts)
}

// Record is a free data retrieval call binding the contract method 0x3465d6d5.
//
// Solidity: function record(address ) view returns(bytes publicKey, bytes pop)
func (_KIP113Mock *KIP113MockCaller) Record(opts *bind.CallOpts, arg0 common.Address) (struct {
	PublicKey []byte
	Pop       []byte
}, error,
) {
	var out []interface{}
	err := _KIP113Mock.contract.Call(opts, &out, "record", arg0)

	outstruct := new(struct {
		PublicKey []byte
		Pop       []byte
	})

	outstruct.PublicKey = *abi.ConvertType(out[0], new([]byte)).(*[]byte)
	outstruct.Pop = *abi.ConvertType(out[1], new([]byte)).(*[]byte)
	return *outstruct, err
}

// Record is a free data retrieval call binding the contract method 0x3465d6d5.
//
// Solidity: function record(address ) view returns(bytes publicKey, bytes pop)
func (_KIP113Mock *KIP113MockSession) Record(arg0 common.Address) (struct {
	PublicKey []byte
	Pop       []byte
}, error,
) {
	return _KIP113Mock.Contract.Record(&_KIP113Mock.CallOpts, arg0)
}

// Record is a free data retrieval call binding the contract method 0x3465d6d5.
//
// Solidity: function record(address ) view returns(bytes publicKey, bytes pop)
func (_KIP113Mock *KIP113MockCallerSession) Record(arg0 common.Address) (struct {
	PublicKey []byte
	Pop       []byte
}, error,
) {
	return _KIP113Mock.Contract.Record(&_KIP113Mock.CallOpts, arg0)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_KIP113Mock *KIP113MockTransactor) Initialize(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KIP113Mock.contract.Transact(opts, "initialize")
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_KIP113Mock *KIP113MockSession) Initialize() (*types.Transaction, error) {
	return _KIP113Mock.Contract.Initialize(&_KIP113Mock.TransactOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_KIP113Mock *KIP113MockTransactorSession) Initialize() (*types.Transaction, error) {
	return _KIP113Mock.Contract.Initialize(&_KIP113Mock.TransactOpts)
}

// Register is a paid mutator transaction binding the contract method 0x786cd4d7.
//
// Solidity: function register(address addr, bytes publicKey, bytes pop) returns()
func (_KIP113Mock *KIP113MockTransactor) Register(opts *bind.TransactOpts, addr common.Address, publicKey []byte, pop []byte) (*types.Transaction, error) {
	return _KIP113Mock.contract.Transact(opts, "register", addr, publicKey, pop)
}

// Register is a paid mutator transaction binding the contract method 0x786cd4d7.
//
// Solidity: function register(address addr, bytes publicKey, bytes pop) returns()
func (_KIP113Mock *KIP113MockSession) Register(addr common.Address, publicKey []byte, pop []byte) (*types.Transaction, error) {
	return _KIP113Mock.Contract.Register(&_KIP113Mock.TransactOpts, addr, publicKey, pop)
}

// Register is a paid mutator transaction binding the contract method 0x786cd4d7.
//
// Solidity: function register(address addr, bytes publicKey, bytes pop) returns()
func (_KIP113Mock *KIP113MockTransactorSession) Register(addr common.Address, publicKey []byte, pop []byte) (*types.Transaction, error) {
	return _KIP113Mock.Contract.Register(&_KIP113Mock.TransactOpts, addr, publicKey, pop)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_KIP113Mock *KIP113MockTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KIP113Mock.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_KIP113Mock *KIP113MockSession) RenounceOwnership() (*types.Transaction, error) {
	return _KIP113Mock.Contract.RenounceOwnership(&_KIP113Mock.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_KIP113Mock *KIP113MockTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _KIP113Mock.Contract.RenounceOwnership(&_KIP113Mock.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_KIP113Mock *KIP113MockTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _KIP113Mock.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_KIP113Mock *KIP113MockSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _KIP113Mock.Contract.TransferOwnership(&_KIP113Mock.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_KIP113Mock *KIP113MockTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _KIP113Mock.Contract.TransferOwnership(&_KIP113Mock.TransactOpts, newOwner)
}

// Unregister is a paid mutator transaction binding the contract method 0x2ec2c246.
//
// Solidity: function unregister(address cnNodeId) returns()
func (_KIP113Mock *KIP113MockTransactor) Unregister(opts *bind.TransactOpts, cnNodeId common.Address) (*types.Transaction, error) {
	return _KIP113Mock.contract.Transact(opts, "unregister", cnNodeId)
}

// Unregister is a paid mutator transaction binding the contract method 0x2ec2c246.
//
// Solidity: function unregister(address cnNodeId) returns()
func (_KIP113Mock *KIP113MockSession) Unregister(cnNodeId common.Address) (*types.Transaction, error) {
	return _KIP113Mock.Contract.Unregister(&_KIP113Mock.TransactOpts, cnNodeId)
}

// Unregister is a paid mutator transaction binding the contract method 0x2ec2c246.
//
// Solidity: function unregister(address cnNodeId) returns()
func (_KIP113Mock *KIP113MockTransactorSession) Unregister(cnNodeId common.Address) (*types.Transaction, error) {
	return _KIP113Mock.Contract.Unregister(&_KIP113Mock.TransactOpts, cnNodeId)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_KIP113Mock *KIP113MockTransactor) UpgradeTo(opts *bind.TransactOpts, newImplementation common.Address) (*types.Transaction, error) {
	return _KIP113Mock.contract.Transact(opts, "upgradeTo", newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_KIP113Mock *KIP113MockSession) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _KIP113Mock.Contract.UpgradeTo(&_KIP113Mock.TransactOpts, newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_KIP113Mock *KIP113MockTransactorSession) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _KIP113Mock.Contract.UpgradeTo(&_KIP113Mock.TransactOpts, newImplementation)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_KIP113Mock *KIP113MockTransactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _KIP113Mock.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_KIP113Mock *KIP113MockSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _KIP113Mock.Contract.UpgradeToAndCall(&_KIP113Mock.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_KIP113Mock *KIP113MockTransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _KIP113Mock.Contract.UpgradeToAndCall(&_KIP113Mock.TransactOpts, newImplementation, data)
}

// KIP113MockAdminChangedIterator is returned from FilterAdminChanged and is used to iterate over the raw logs and unpacked data for AdminChanged events raised by the KIP113Mock contract.
type KIP113MockAdminChangedIterator struct {
	Event *KIP113MockAdminChanged // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log      // Log channel receiving the found contract events
	sub  klaytn.Subscription // Subscription for errors, completion and termination
	done bool                // Whether the subscription completed delivering logs
	fail error               // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *KIP113MockAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KIP113MockAdminChanged)
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
		it.Event = new(KIP113MockAdminChanged)
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
func (it *KIP113MockAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KIP113MockAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KIP113MockAdminChanged represents a AdminChanged event raised by the KIP113Mock contract.
type KIP113MockAdminChanged struct {
	PreviousAdmin common.Address
	NewAdmin      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterAdminChanged is a free log retrieval operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_KIP113Mock *KIP113MockFilterer) FilterAdminChanged(opts *bind.FilterOpts) (*KIP113MockAdminChangedIterator, error) {
	logs, sub, err := _KIP113Mock.contract.FilterLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return &KIP113MockAdminChangedIterator{contract: _KIP113Mock.contract, event: "AdminChanged", logs: logs, sub: sub}, nil
}

// WatchAdminChanged is a free log subscription operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_KIP113Mock *KIP113MockFilterer) WatchAdminChanged(opts *bind.WatchOpts, sink chan<- *KIP113MockAdminChanged) (event.Subscription, error) {
	logs, sub, err := _KIP113Mock.contract.WatchLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KIP113MockAdminChanged)
				if err := _KIP113Mock.contract.UnpackLog(event, "AdminChanged", log); err != nil {
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

// ParseAdminChanged is a log parse operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_KIP113Mock *KIP113MockFilterer) ParseAdminChanged(log types.Log) (*KIP113MockAdminChanged, error) {
	event := new(KIP113MockAdminChanged)
	if err := _KIP113Mock.contract.UnpackLog(event, "AdminChanged", log); err != nil {
		return nil, err
	}
	return event, nil
}

// KIP113MockBeaconUpgradedIterator is returned from FilterBeaconUpgraded and is used to iterate over the raw logs and unpacked data for BeaconUpgraded events raised by the KIP113Mock contract.
type KIP113MockBeaconUpgradedIterator struct {
	Event *KIP113MockBeaconUpgraded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log      // Log channel receiving the found contract events
	sub  klaytn.Subscription // Subscription for errors, completion and termination
	done bool                // Whether the subscription completed delivering logs
	fail error               // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *KIP113MockBeaconUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KIP113MockBeaconUpgraded)
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
		it.Event = new(KIP113MockBeaconUpgraded)
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
func (it *KIP113MockBeaconUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KIP113MockBeaconUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KIP113MockBeaconUpgraded represents a BeaconUpgraded event raised by the KIP113Mock contract.
type KIP113MockBeaconUpgraded struct {
	Beacon common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterBeaconUpgraded is a free log retrieval operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_KIP113Mock *KIP113MockFilterer) FilterBeaconUpgraded(opts *bind.FilterOpts, beacon []common.Address) (*KIP113MockBeaconUpgradedIterator, error) {
	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _KIP113Mock.contract.FilterLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return &KIP113MockBeaconUpgradedIterator{contract: _KIP113Mock.contract, event: "BeaconUpgraded", logs: logs, sub: sub}, nil
}

// WatchBeaconUpgraded is a free log subscription operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_KIP113Mock *KIP113MockFilterer) WatchBeaconUpgraded(opts *bind.WatchOpts, sink chan<- *KIP113MockBeaconUpgraded, beacon []common.Address) (event.Subscription, error) {
	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _KIP113Mock.contract.WatchLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KIP113MockBeaconUpgraded)
				if err := _KIP113Mock.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
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

// ParseBeaconUpgraded is a log parse operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_KIP113Mock *KIP113MockFilterer) ParseBeaconUpgraded(log types.Log) (*KIP113MockBeaconUpgraded, error) {
	event := new(KIP113MockBeaconUpgraded)
	if err := _KIP113Mock.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
		return nil, err
	}
	return event, nil
}

// KIP113MockInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the KIP113Mock contract.
type KIP113MockInitializedIterator struct {
	Event *KIP113MockInitialized // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log      // Log channel receiving the found contract events
	sub  klaytn.Subscription // Subscription for errors, completion and termination
	done bool                // Whether the subscription completed delivering logs
	fail error               // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *KIP113MockInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KIP113MockInitialized)
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
		it.Event = new(KIP113MockInitialized)
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
func (it *KIP113MockInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KIP113MockInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KIP113MockInitialized represents a Initialized event raised by the KIP113Mock contract.
type KIP113MockInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_KIP113Mock *KIP113MockFilterer) FilterInitialized(opts *bind.FilterOpts) (*KIP113MockInitializedIterator, error) {
	logs, sub, err := _KIP113Mock.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &KIP113MockInitializedIterator{contract: _KIP113Mock.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_KIP113Mock *KIP113MockFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *KIP113MockInitialized) (event.Subscription, error) {
	logs, sub, err := _KIP113Mock.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KIP113MockInitialized)
				if err := _KIP113Mock.contract.UnpackLog(event, "Initialized", log); err != nil {
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

// ParseInitialized is a log parse operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_KIP113Mock *KIP113MockFilterer) ParseInitialized(log types.Log) (*KIP113MockInitialized, error) {
	event := new(KIP113MockInitialized)
	if err := _KIP113Mock.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	return event, nil
}

// KIP113MockOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the KIP113Mock contract.
type KIP113MockOwnershipTransferredIterator struct {
	Event *KIP113MockOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log      // Log channel receiving the found contract events
	sub  klaytn.Subscription // Subscription for errors, completion and termination
	done bool                // Whether the subscription completed delivering logs
	fail error               // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *KIP113MockOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KIP113MockOwnershipTransferred)
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
		it.Event = new(KIP113MockOwnershipTransferred)
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
func (it *KIP113MockOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KIP113MockOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KIP113MockOwnershipTransferred represents a OwnershipTransferred event raised by the KIP113Mock contract.
type KIP113MockOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_KIP113Mock *KIP113MockFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*KIP113MockOwnershipTransferredIterator, error) {
	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _KIP113Mock.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &KIP113MockOwnershipTransferredIterator{contract: _KIP113Mock.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_KIP113Mock *KIP113MockFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *KIP113MockOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {
	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _KIP113Mock.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KIP113MockOwnershipTransferred)
				if err := _KIP113Mock.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_KIP113Mock *KIP113MockFilterer) ParseOwnershipTransferred(log types.Log) (*KIP113MockOwnershipTransferred, error) {
	event := new(KIP113MockOwnershipTransferred)
	if err := _KIP113Mock.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	return event, nil
}

// KIP113MockRegisteredIterator is returned from FilterRegistered and is used to iterate over the raw logs and unpacked data for Registered events raised by the KIP113Mock contract.
type KIP113MockRegisteredIterator struct {
	Event *KIP113MockRegistered // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log      // Log channel receiving the found contract events
	sub  klaytn.Subscription // Subscription for errors, completion and termination
	done bool                // Whether the subscription completed delivering logs
	fail error               // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *KIP113MockRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KIP113MockRegistered)
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
		it.Event = new(KIP113MockRegistered)
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
func (it *KIP113MockRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KIP113MockRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KIP113MockRegistered represents a Registered event raised by the KIP113Mock contract.
type KIP113MockRegistered struct {
	CnNodeId  common.Address
	PublicKey []byte
	Pop       []byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterRegistered is a free log retrieval operation binding the contract event 0x79c75399e89a1f580d9a6252cb8bdcf4cd80f73b3597c94d845eb52174ad930f.
//
// Solidity: event Registered(address cnNodeId, bytes publicKey, bytes pop)
func (_KIP113Mock *KIP113MockFilterer) FilterRegistered(opts *bind.FilterOpts) (*KIP113MockRegisteredIterator, error) {
	logs, sub, err := _KIP113Mock.contract.FilterLogs(opts, "Registered")
	if err != nil {
		return nil, err
	}
	return &KIP113MockRegisteredIterator{contract: _KIP113Mock.contract, event: "Registered", logs: logs, sub: sub}, nil
}

// WatchRegistered is a free log subscription operation binding the contract event 0x79c75399e89a1f580d9a6252cb8bdcf4cd80f73b3597c94d845eb52174ad930f.
//
// Solidity: event Registered(address cnNodeId, bytes publicKey, bytes pop)
func (_KIP113Mock *KIP113MockFilterer) WatchRegistered(opts *bind.WatchOpts, sink chan<- *KIP113MockRegistered) (event.Subscription, error) {
	logs, sub, err := _KIP113Mock.contract.WatchLogs(opts, "Registered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KIP113MockRegistered)
				if err := _KIP113Mock.contract.UnpackLog(event, "Registered", log); err != nil {
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

// ParseRegistered is a log parse operation binding the contract event 0x79c75399e89a1f580d9a6252cb8bdcf4cd80f73b3597c94d845eb52174ad930f.
//
// Solidity: event Registered(address cnNodeId, bytes publicKey, bytes pop)
func (_KIP113Mock *KIP113MockFilterer) ParseRegistered(log types.Log) (*KIP113MockRegistered, error) {
	event := new(KIP113MockRegistered)
	if err := _KIP113Mock.contract.UnpackLog(event, "Registered", log); err != nil {
		return nil, err
	}
	return event, nil
}

// KIP113MockUnregisteredIterator is returned from FilterUnregistered and is used to iterate over the raw logs and unpacked data for Unregistered events raised by the KIP113Mock contract.
type KIP113MockUnregisteredIterator struct {
	Event *KIP113MockUnregistered // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log      // Log channel receiving the found contract events
	sub  klaytn.Subscription // Subscription for errors, completion and termination
	done bool                // Whether the subscription completed delivering logs
	fail error               // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *KIP113MockUnregisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KIP113MockUnregistered)
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
		it.Event = new(KIP113MockUnregistered)
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
func (it *KIP113MockUnregisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KIP113MockUnregisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KIP113MockUnregistered represents a Unregistered event raised by the KIP113Mock contract.
type KIP113MockUnregistered struct {
	CnNodeId  common.Address
	PublicKey []byte
	Pop       []byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterUnregistered is a free log retrieval operation binding the contract event 0xb98b07c4d52e8d65fa5416810f2746a810eb074b1ac7784e1b07e315c0dfd2d9.
//
// Solidity: event Unregistered(address cnNodeId, bytes publicKey, bytes pop)
func (_KIP113Mock *KIP113MockFilterer) FilterUnregistered(opts *bind.FilterOpts) (*KIP113MockUnregisteredIterator, error) {
	logs, sub, err := _KIP113Mock.contract.FilterLogs(opts, "Unregistered")
	if err != nil {
		return nil, err
	}
	return &KIP113MockUnregisteredIterator{contract: _KIP113Mock.contract, event: "Unregistered", logs: logs, sub: sub}, nil
}

// WatchUnregistered is a free log subscription operation binding the contract event 0xb98b07c4d52e8d65fa5416810f2746a810eb074b1ac7784e1b07e315c0dfd2d9.
//
// Solidity: event Unregistered(address cnNodeId, bytes publicKey, bytes pop)
func (_KIP113Mock *KIP113MockFilterer) WatchUnregistered(opts *bind.WatchOpts, sink chan<- *KIP113MockUnregistered) (event.Subscription, error) {
	logs, sub, err := _KIP113Mock.contract.WatchLogs(opts, "Unregistered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KIP113MockUnregistered)
				if err := _KIP113Mock.contract.UnpackLog(event, "Unregistered", log); err != nil {
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

// ParseUnregistered is a log parse operation binding the contract event 0xb98b07c4d52e8d65fa5416810f2746a810eb074b1ac7784e1b07e315c0dfd2d9.
//
// Solidity: event Unregistered(address cnNodeId, bytes publicKey, bytes pop)
func (_KIP113Mock *KIP113MockFilterer) ParseUnregistered(log types.Log) (*KIP113MockUnregistered, error) {
	event := new(KIP113MockUnregistered)
	if err := _KIP113Mock.contract.UnpackLog(event, "Unregistered", log); err != nil {
		return nil, err
	}
	return event, nil
}

// KIP113MockUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the KIP113Mock contract.
type KIP113MockUpgradedIterator struct {
	Event *KIP113MockUpgraded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log      // Log channel receiving the found contract events
	sub  klaytn.Subscription // Subscription for errors, completion and termination
	done bool                // Whether the subscription completed delivering logs
	fail error               // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *KIP113MockUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KIP113MockUpgraded)
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
		it.Event = new(KIP113MockUpgraded)
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
func (it *KIP113MockUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KIP113MockUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KIP113MockUpgraded represents a Upgraded event raised by the KIP113Mock contract.
type KIP113MockUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_KIP113Mock *KIP113MockFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*KIP113MockUpgradedIterator, error) {
	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _KIP113Mock.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &KIP113MockUpgradedIterator{contract: _KIP113Mock.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_KIP113Mock *KIP113MockFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *KIP113MockUpgraded, implementation []common.Address) (event.Subscription, error) {
	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _KIP113Mock.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KIP113MockUpgraded)
				if err := _KIP113Mock.contract.UnpackLog(event, "Upgraded", log); err != nil {
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

// ParseUpgraded is a log parse operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_KIP113Mock *KIP113MockFilterer) ParseUpgraded(log types.Log) (*KIP113MockUpgraded, error) {
	event := new(KIP113MockUpgraded)
	if err := _KIP113Mock.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	return event, nil
}

// OwnableUpgradeableMetaData contains all meta data concerning the OwnableUpgradeable contract.
var OwnableUpgradeableMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"8da5cb5b": "owner()",
		"715018a6": "renounceOwnership()",
		"f2fde38b": "transferOwnership(address)",
	},
}

// OwnableUpgradeableABI is the input ABI used to generate the binding from.
// Deprecated: Use OwnableUpgradeableMetaData.ABI instead.
var OwnableUpgradeableABI = OwnableUpgradeableMetaData.ABI

// OwnableUpgradeableBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const OwnableUpgradeableBinRuntime = ``

// OwnableUpgradeableFuncSigs maps the 4-byte function signature to its string representation.
// Deprecated: Use OwnableUpgradeableMetaData.Sigs instead.
var OwnableUpgradeableFuncSigs = OwnableUpgradeableMetaData.Sigs

// OwnableUpgradeable is an auto generated Go binding around a Klaytn contract.
type OwnableUpgradeable struct {
	OwnableUpgradeableCaller     // Read-only binding to the contract
	OwnableUpgradeableTransactor // Write-only binding to the contract
	OwnableUpgradeableFilterer   // Log filterer for contract events
}

// OwnableUpgradeableCaller is an auto generated read-only Go binding around a Klaytn contract.
type OwnableUpgradeableCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OwnableUpgradeableTransactor is an auto generated write-only Go binding around a Klaytn contract.
type OwnableUpgradeableTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OwnableUpgradeableFilterer is an auto generated log filtering Go binding around a Klaytn contract events.
type OwnableUpgradeableFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OwnableUpgradeableSession is an auto generated Go binding around a Klaytn contract,
// with pre-set call and transact options.
type OwnableUpgradeableSession struct {
	Contract     *OwnableUpgradeable // Generic contract binding to set the session for
	CallOpts     bind.CallOpts       // Call options to use throughout this session
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// OwnableUpgradeableCallerSession is an auto generated read-only Go binding around a Klaytn contract,
// with pre-set call options.
type OwnableUpgradeableCallerSession struct {
	Contract *OwnableUpgradeableCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts             // Call options to use throughout this session
}

// OwnableUpgradeableTransactorSession is an auto generated write-only Go binding around a Klaytn contract,
// with pre-set transact options.
type OwnableUpgradeableTransactorSession struct {
	Contract     *OwnableUpgradeableTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts             // Transaction auth options to use throughout this session
}

// OwnableUpgradeableRaw is an auto generated low-level Go binding around a Klaytn contract.
type OwnableUpgradeableRaw struct {
	Contract *OwnableUpgradeable // Generic contract binding to access the raw methods on
}

// OwnableUpgradeableCallerRaw is an auto generated low-level read-only Go binding around a Klaytn contract.
type OwnableUpgradeableCallerRaw struct {
	Contract *OwnableUpgradeableCaller // Generic read-only contract binding to access the raw methods on
}

// OwnableUpgradeableTransactorRaw is an auto generated low-level write-only Go binding around a Klaytn contract.
type OwnableUpgradeableTransactorRaw struct {
	Contract *OwnableUpgradeableTransactor // Generic write-only contract binding to access the raw methods on
}

// NewOwnableUpgradeable creates a new instance of OwnableUpgradeable, bound to a specific deployed contract.
func NewOwnableUpgradeable(address common.Address, backend bind.ContractBackend) (*OwnableUpgradeable, error) {
	contract, err := bindOwnableUpgradeable(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &OwnableUpgradeable{OwnableUpgradeableCaller: OwnableUpgradeableCaller{contract: contract}, OwnableUpgradeableTransactor: OwnableUpgradeableTransactor{contract: contract}, OwnableUpgradeableFilterer: OwnableUpgradeableFilterer{contract: contract}}, nil
}

// NewOwnableUpgradeableCaller creates a new read-only instance of OwnableUpgradeable, bound to a specific deployed contract.
func NewOwnableUpgradeableCaller(address common.Address, caller bind.ContractCaller) (*OwnableUpgradeableCaller, error) {
	contract, err := bindOwnableUpgradeable(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OwnableUpgradeableCaller{contract: contract}, nil
}

// NewOwnableUpgradeableTransactor creates a new write-only instance of OwnableUpgradeable, bound to a specific deployed contract.
func NewOwnableUpgradeableTransactor(address common.Address, transactor bind.ContractTransactor) (*OwnableUpgradeableTransactor, error) {
	contract, err := bindOwnableUpgradeable(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OwnableUpgradeableTransactor{contract: contract}, nil
}

// NewOwnableUpgradeableFilterer creates a new log filterer instance of OwnableUpgradeable, bound to a specific deployed contract.
func NewOwnableUpgradeableFilterer(address common.Address, filterer bind.ContractFilterer) (*OwnableUpgradeableFilterer, error) {
	contract, err := bindOwnableUpgradeable(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OwnableUpgradeableFilterer{contract: contract}, nil
}

// bindOwnableUpgradeable binds a generic wrapper to an already deployed contract.
func bindOwnableUpgradeable(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := OwnableUpgradeableMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OwnableUpgradeable *OwnableUpgradeableRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OwnableUpgradeable.Contract.OwnableUpgradeableCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OwnableUpgradeable *OwnableUpgradeableRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OwnableUpgradeable.Contract.OwnableUpgradeableTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OwnableUpgradeable *OwnableUpgradeableRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OwnableUpgradeable.Contract.OwnableUpgradeableTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OwnableUpgradeable *OwnableUpgradeableCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OwnableUpgradeable.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OwnableUpgradeable *OwnableUpgradeableTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OwnableUpgradeable.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OwnableUpgradeable *OwnableUpgradeableTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OwnableUpgradeable.Contract.contract.Transact(opts, method, params...)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_OwnableUpgradeable *OwnableUpgradeableCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _OwnableUpgradeable.contract.Call(opts, &out, "owner")
	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_OwnableUpgradeable *OwnableUpgradeableSession) Owner() (common.Address, error) {
	return _OwnableUpgradeable.Contract.Owner(&_OwnableUpgradeable.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_OwnableUpgradeable *OwnableUpgradeableCallerSession) Owner() (common.Address, error) {
	return _OwnableUpgradeable.Contract.Owner(&_OwnableUpgradeable.CallOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_OwnableUpgradeable *OwnableUpgradeableTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OwnableUpgradeable.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_OwnableUpgradeable *OwnableUpgradeableSession) RenounceOwnership() (*types.Transaction, error) {
	return _OwnableUpgradeable.Contract.RenounceOwnership(&_OwnableUpgradeable.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_OwnableUpgradeable *OwnableUpgradeableTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _OwnableUpgradeable.Contract.RenounceOwnership(&_OwnableUpgradeable.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_OwnableUpgradeable *OwnableUpgradeableTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _OwnableUpgradeable.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_OwnableUpgradeable *OwnableUpgradeableSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _OwnableUpgradeable.Contract.TransferOwnership(&_OwnableUpgradeable.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_OwnableUpgradeable *OwnableUpgradeableTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _OwnableUpgradeable.Contract.TransferOwnership(&_OwnableUpgradeable.TransactOpts, newOwner)
}

// OwnableUpgradeableInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the OwnableUpgradeable contract.
type OwnableUpgradeableInitializedIterator struct {
	Event *OwnableUpgradeableInitialized // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log      // Log channel receiving the found contract events
	sub  klaytn.Subscription // Subscription for errors, completion and termination
	done bool                // Whether the subscription completed delivering logs
	fail error               // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OwnableUpgradeableInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OwnableUpgradeableInitialized)
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
		it.Event = new(OwnableUpgradeableInitialized)
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
func (it *OwnableUpgradeableInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OwnableUpgradeableInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OwnableUpgradeableInitialized represents a Initialized event raised by the OwnableUpgradeable contract.
type OwnableUpgradeableInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_OwnableUpgradeable *OwnableUpgradeableFilterer) FilterInitialized(opts *bind.FilterOpts) (*OwnableUpgradeableInitializedIterator, error) {
	logs, sub, err := _OwnableUpgradeable.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &OwnableUpgradeableInitializedIterator{contract: _OwnableUpgradeable.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_OwnableUpgradeable *OwnableUpgradeableFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *OwnableUpgradeableInitialized) (event.Subscription, error) {
	logs, sub, err := _OwnableUpgradeable.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OwnableUpgradeableInitialized)
				if err := _OwnableUpgradeable.contract.UnpackLog(event, "Initialized", log); err != nil {
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

// ParseInitialized is a log parse operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_OwnableUpgradeable *OwnableUpgradeableFilterer) ParseInitialized(log types.Log) (*OwnableUpgradeableInitialized, error) {
	event := new(OwnableUpgradeableInitialized)
	if err := _OwnableUpgradeable.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	return event, nil
}

// OwnableUpgradeableOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the OwnableUpgradeable contract.
type OwnableUpgradeableOwnershipTransferredIterator struct {
	Event *OwnableUpgradeableOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log      // Log channel receiving the found contract events
	sub  klaytn.Subscription // Subscription for errors, completion and termination
	done bool                // Whether the subscription completed delivering logs
	fail error               // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OwnableUpgradeableOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OwnableUpgradeableOwnershipTransferred)
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
		it.Event = new(OwnableUpgradeableOwnershipTransferred)
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
func (it *OwnableUpgradeableOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OwnableUpgradeableOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OwnableUpgradeableOwnershipTransferred represents a OwnershipTransferred event raised by the OwnableUpgradeable contract.
type OwnableUpgradeableOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_OwnableUpgradeable *OwnableUpgradeableFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*OwnableUpgradeableOwnershipTransferredIterator, error) {
	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _OwnableUpgradeable.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &OwnableUpgradeableOwnershipTransferredIterator{contract: _OwnableUpgradeable.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_OwnableUpgradeable *OwnableUpgradeableFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *OwnableUpgradeableOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {
	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _OwnableUpgradeable.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OwnableUpgradeableOwnershipTransferred)
				if err := _OwnableUpgradeable.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_OwnableUpgradeable *OwnableUpgradeableFilterer) ParseOwnershipTransferred(log types.Log) (*OwnableUpgradeableOwnershipTransferred, error) {
	event := new(OwnableUpgradeableOwnershipTransferred)
	if err := _OwnableUpgradeable.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ProxyMetaData contains all meta data concerning the Proxy contract.
var ProxyMetaData = &bind.MetaData{
	ABI: "[{\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
}

// ProxyABI is the input ABI used to generate the binding from.
// Deprecated: Use ProxyMetaData.ABI instead.
var ProxyABI = ProxyMetaData.ABI

// ProxyBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const ProxyBinRuntime = ``

// Proxy is an auto generated Go binding around a Klaytn contract.
type Proxy struct {
	ProxyCaller     // Read-only binding to the contract
	ProxyTransactor // Write-only binding to the contract
	ProxyFilterer   // Log filterer for contract events
}

// ProxyCaller is an auto generated read-only Go binding around a Klaytn contract.
type ProxyCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ProxyTransactor is an auto generated write-only Go binding around a Klaytn contract.
type ProxyTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ProxyFilterer is an auto generated log filtering Go binding around a Klaytn contract events.
type ProxyFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ProxySession is an auto generated Go binding around a Klaytn contract,
// with pre-set call and transact options.
type ProxySession struct {
	Contract     *Proxy            // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ProxyCallerSession is an auto generated read-only Go binding around a Klaytn contract,
// with pre-set call options.
type ProxyCallerSession struct {
	Contract *ProxyCaller  // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// ProxyTransactorSession is an auto generated write-only Go binding around a Klaytn contract,
// with pre-set transact options.
type ProxyTransactorSession struct {
	Contract     *ProxyTransactor  // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ProxyRaw is an auto generated low-level Go binding around a Klaytn contract.
type ProxyRaw struct {
	Contract *Proxy // Generic contract binding to access the raw methods on
}

// ProxyCallerRaw is an auto generated low-level read-only Go binding around a Klaytn contract.
type ProxyCallerRaw struct {
	Contract *ProxyCaller // Generic read-only contract binding to access the raw methods on
}

// ProxyTransactorRaw is an auto generated low-level write-only Go binding around a Klaytn contract.
type ProxyTransactorRaw struct {
	Contract *ProxyTransactor // Generic write-only contract binding to access the raw methods on
}

// NewProxy creates a new instance of Proxy, bound to a specific deployed contract.
func NewProxy(address common.Address, backend bind.ContractBackend) (*Proxy, error) {
	contract, err := bindProxy(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Proxy{ProxyCaller: ProxyCaller{contract: contract}, ProxyTransactor: ProxyTransactor{contract: contract}, ProxyFilterer: ProxyFilterer{contract: contract}}, nil
}

// NewProxyCaller creates a new read-only instance of Proxy, bound to a specific deployed contract.
func NewProxyCaller(address common.Address, caller bind.ContractCaller) (*ProxyCaller, error) {
	contract, err := bindProxy(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ProxyCaller{contract: contract}, nil
}

// NewProxyTransactor creates a new write-only instance of Proxy, bound to a specific deployed contract.
func NewProxyTransactor(address common.Address, transactor bind.ContractTransactor) (*ProxyTransactor, error) {
	contract, err := bindProxy(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ProxyTransactor{contract: contract}, nil
}

// NewProxyFilterer creates a new log filterer instance of Proxy, bound to a specific deployed contract.
func NewProxyFilterer(address common.Address, filterer bind.ContractFilterer) (*ProxyFilterer, error) {
	contract, err := bindProxy(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ProxyFilterer{contract: contract}, nil
}

// bindProxy binds a generic wrapper to an already deployed contract.
func bindProxy(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ProxyMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Proxy *ProxyRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Proxy.Contract.ProxyCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Proxy *ProxyRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Proxy.Contract.ProxyTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Proxy *ProxyRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Proxy.Contract.ProxyTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Proxy *ProxyCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Proxy.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Proxy *ProxyTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Proxy.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Proxy *ProxyTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Proxy.Contract.contract.Transact(opts, method, params...)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Proxy *ProxyTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _Proxy.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Proxy *ProxySession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _Proxy.Contract.Fallback(&_Proxy.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Proxy *ProxyTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _Proxy.Contract.Fallback(&_Proxy.TransactOpts, calldata)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Proxy *ProxyTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Proxy.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Proxy *ProxySession) Receive() (*types.Transaction, error) {
	return _Proxy.Contract.Receive(&_Proxy.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Proxy *ProxyTransactorSession) Receive() (*types.Transaction, error) {
	return _Proxy.Contract.Receive(&_Proxy.TransactOpts)
}

// RegistryMetaData contains all meta data concerning the Registry contract.
var RegistryMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"activation\",\"type\":\"uint256\"}],\"name\":\"Registered\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"name\":\"getActiveAddr\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllNames\",\"outputs\":[{\"internalType\":\"string[]\",\"name\":\"\",\"type\":\"string[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"name\":\"getAllRecords\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"activation\",\"type\":\"uint256\"}],\"internalType\":\"structIRegistry.Record[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"names\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"records\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"activation\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"activation\",\"type\":\"uint256\"}],\"name\":\"register\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"e2693e3f": "getActiveAddr(string)",
		"fb825e5f": "getAllNames()",
		"78d573a2": "getAllRecords(string)",
		"4622ab03": "names(uint256)",
		"8da5cb5b": "owner()",
		"3b51650d": "records(string,uint256)",
		"d393c871": "register(string,address,uint256)",
		"f2fde38b": "transferOwnership(address)",
	},
	Bin: "0x608060405234801561001057600080fd5b50610db9806100206000396000f3fe608060405234801561001057600080fd5b50600436106100885760003560e01c8063d393c8711161005b578063d393c87114610129578063e2693e3f1461013e578063f2fde38b14610151578063fb825e5f1461016457600080fd5b80633b51650d1461008d5780634622ab03146100c457806378d573a2146100e45780638da5cb5b14610104575b600080fd5b6100a061009b366004610975565b610179565b604080516001600160a01b0390931683526020830191909152015b60405180910390f35b6100d76100d23660046109ba565b6101ce565b6040516100bb9190610a23565b6100f76100f2366004610a3d565b61027a565b6040516100bb9190610a7a565b6002546001600160a01b03165b6040516001600160a01b0390911681526020016100bb565b61013c610137366004610aee565b61030d565b005b61011161014c366004610a3d565b61062b565b61013c61015f366004610b45565b610722565b61016c6107f9565b6040516100bb9190610b60565b815160208184018101805160008252928201918501919091209190528054829081106101a457600080fd5b6000918252602090912060029091020180546001909101546001600160a01b039091169250905082565b600181815481106101de57600080fd5b9060005260206000200160009150905080546101f990610bc2565b80601f016020809104026020016040519081016040528092919081815260200182805461022590610bc2565b80156102725780601f1061024757610100808354040283529160200191610272565b820191906000526020600020905b81548152906001019060200180831161025557829003601f168201915b505050505081565b606060008260405161028c9190610bfc565b9081526020016040518091039020805480602002602001604051908101604052809291908181526020016000905b82821015610302576000848152602090819020604080518082019091526002850290910180546001600160a01b031682526001908101548284015290835290920191016102ba565b505050509050919050565b6002546001600160a01b031633146103585760405162461bcd60e51b81526020600482015260096024820152682737ba1037bbb732b960b91b60448201526064015b60405180910390fd5b8260008160405160200161036c9190610bfc565b604051602081830303815290604052905080516000036103bd5760405162461bcd60e51b815260206004820152600c60248201526b456d70747920737472696e6760a01b604482015260640161034f565b4383116104165760405162461bcd60e51b815260206004820152602160248201527f43616e277420726567697374657220636f6e74726163742066726f6d207061736044820152601d60fa1b606482015260840161034f565b600080866040516104279190610bfc565b90815260405190819003602001902054905060008190036104f3576001805480820182556000919091527fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf60161047d8782610c67565b5060008660405161048e9190610bfc565b90815260408051602092819003830181208183019092526001600160a01b0388811682528382018881528354600180820186556000958652959094209251600290940290920180546001600160a01b03191693909116929092178255519101556105e1565b600080876040516105049190610bfc565b90815260405190819003602001902061051e600184610d3d565b8154811061052e5761052e610d56565b90600052602060002090600202019050438160010154116105be576000876040516105599190610bfc565b90815260408051602092819003830181208183019092526001600160a01b0389811682528382018981528354600180820186556000958652959094209251600290940290920180546001600160a01b03191693909116929092178255519101556105df565b80546001600160a01b0319166001600160a01b038716178155600181018590555b505b83856001600160a01b03167f142e1fdac7ecccbc62af925f0b4039db26847b625602e56b1421dfbc8a0e4f308860405161061b9190610a23565b60405180910390a3505050505050565b60008060008360405161063e9190610bfc565b908152604051908190036020019020549050805b801561071857436000856040516106699190610bfc565b908152604051908190036020019020610683600184610d3d565b8154811061069357610693610d56565b90600052602060002090600202016001015411610706576000846040516106ba9190610bfc565b9081526040519081900360200190206106d4600183610d3d565b815481106106e4576106e4610d56565b60009182526020909120600290910201546001600160a01b0316949350505050565b8061071081610d6c565b915050610652565b5060009392505050565b6002546001600160a01b031633146107685760405162461bcd60e51b81526020600482015260096024820152682737ba1037bbb732b960b91b604482015260640161034f565b6001600160a01b0381166107ad5760405162461bcd60e51b815260206004820152600c60248201526b5a65726f206164647265737360a01b604482015260640161034f565b600280546001600160a01b0319166001600160a01b03831690811790915560405133907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a350565b60606001805480602002602001604051908101604052809291908181526020016000905b828210156108c957838290600052602060002001805461083c90610bc2565b80601f016020809104026020016040519081016040528092919081815260200182805461086890610bc2565b80156108b55780601f1061088a576101008083540402835291602001916108b5565b820191906000526020600020905b81548152906001019060200180831161089857829003601f168201915b50505050508152602001906001019061081d565b50505050905090565b634e487b7160e01b600052604160045260246000fd5b600082601f8301126108f957600080fd5b813567ffffffffffffffff80821115610914576109146108d2565b604051601f8301601f19908116603f0116810190828211818310171561093c5761093c6108d2565b8160405283815286602085880101111561095557600080fd5b836020870160208301376000602085830101528094505050505092915050565b6000806040838503121561098857600080fd5b823567ffffffffffffffff81111561099f57600080fd5b6109ab858286016108e8565b95602094909401359450505050565b6000602082840312156109cc57600080fd5b5035919050565b60005b838110156109ee5781810151838201526020016109d6565b50506000910152565b60008151808452610a0f8160208601602086016109d3565b601f01601f19169290920160200192915050565b602081526000610a3660208301846109f7565b9392505050565b600060208284031215610a4f57600080fd5b813567ffffffffffffffff811115610a6657600080fd5b610a72848285016108e8565b949350505050565b602080825282518282018190526000919060409081850190868401855b82811015610ac557815180516001600160a01b03168552860151868501529284019290850190600101610a97565b5091979650505050505050565b80356001600160a01b0381168114610ae957600080fd5b919050565b600080600060608486031215610b0357600080fd5b833567ffffffffffffffff811115610b1a57600080fd5b610b26868287016108e8565b935050610b3560208501610ad2565b9150604084013590509250925092565b600060208284031215610b5757600080fd5b610a3682610ad2565b6000602080830181845280855180835260408601915060408160051b870101925083870160005b82811015610bb557603f19888603018452610ba38583516109f7565b94509285019290850190600101610b87565b5092979650505050505050565b600181811c90821680610bd657607f821691505b602082108103610bf657634e487b7160e01b600052602260045260246000fd5b50919050565b60008251610c0e8184602087016109d3565b9190910192915050565b601f821115610c6257600081815260208120601f850160051c81016020861015610c3f5750805b601f850160051c820191505b81811015610c5e57828155600101610c4b565b5050505b505050565b815167ffffffffffffffff811115610c8157610c816108d2565b610c9581610c8f8454610bc2565b84610c18565b602080601f831160018114610cca5760008415610cb25750858301515b600019600386901b1c1916600185901b178555610c5e565b600085815260208120601f198616915b82811015610cf957888601518255948401946001909101908401610cda565b5085821015610d175787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b634e487b7160e01b600052601160045260246000fd5b81810381811115610d5057610d50610d27565b92915050565b634e487b7160e01b600052603260045260246000fd5b600081610d7b57610d7b610d27565b50600019019056fea2646970667358221220a3f6e37a5b67f7bb6210f0cfef969cb855bda4c8699bb46ed4a49ee814b7765864736f6c63430008130033",
}

// RegistryABI is the input ABI used to generate the binding from.
// Deprecated: Use RegistryMetaData.ABI instead.
var RegistryABI = RegistryMetaData.ABI

// RegistryBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const RegistryBinRuntime = `608060405234801561001057600080fd5b50600436106100885760003560e01c8063d393c8711161005b578063d393c87114610129578063e2693e3f1461013e578063f2fde38b14610151578063fb825e5f1461016457600080fd5b80633b51650d1461008d5780634622ab03146100c457806378d573a2146100e45780638da5cb5b14610104575b600080fd5b6100a061009b366004610975565b610179565b604080516001600160a01b0390931683526020830191909152015b60405180910390f35b6100d76100d23660046109ba565b6101ce565b6040516100bb9190610a23565b6100f76100f2366004610a3d565b61027a565b6040516100bb9190610a7a565b6002546001600160a01b03165b6040516001600160a01b0390911681526020016100bb565b61013c610137366004610aee565b61030d565b005b61011161014c366004610a3d565b61062b565b61013c61015f366004610b45565b610722565b61016c6107f9565b6040516100bb9190610b60565b815160208184018101805160008252928201918501919091209190528054829081106101a457600080fd5b6000918252602090912060029091020180546001909101546001600160a01b039091169250905082565b600181815481106101de57600080fd5b9060005260206000200160009150905080546101f990610bc2565b80601f016020809104026020016040519081016040528092919081815260200182805461022590610bc2565b80156102725780601f1061024757610100808354040283529160200191610272565b820191906000526020600020905b81548152906001019060200180831161025557829003601f168201915b505050505081565b606060008260405161028c9190610bfc565b9081526020016040518091039020805480602002602001604051908101604052809291908181526020016000905b82821015610302576000848152602090819020604080518082019091526002850290910180546001600160a01b031682526001908101548284015290835290920191016102ba565b505050509050919050565b6002546001600160a01b031633146103585760405162461bcd60e51b81526020600482015260096024820152682737ba1037bbb732b960b91b60448201526064015b60405180910390fd5b8260008160405160200161036c9190610bfc565b604051602081830303815290604052905080516000036103bd5760405162461bcd60e51b815260206004820152600c60248201526b456d70747920737472696e6760a01b604482015260640161034f565b4383116104165760405162461bcd60e51b815260206004820152602160248201527f43616e277420726567697374657220636f6e74726163742066726f6d207061736044820152601d60fa1b606482015260840161034f565b600080866040516104279190610bfc565b90815260405190819003602001902054905060008190036104f3576001805480820182556000919091527fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf60161047d8782610c67565b5060008660405161048e9190610bfc565b90815260408051602092819003830181208183019092526001600160a01b0388811682528382018881528354600180820186556000958652959094209251600290940290920180546001600160a01b03191693909116929092178255519101556105e1565b600080876040516105049190610bfc565b90815260405190819003602001902061051e600184610d3d565b8154811061052e5761052e610d56565b90600052602060002090600202019050438160010154116105be576000876040516105599190610bfc565b90815260408051602092819003830181208183019092526001600160a01b0389811682528382018981528354600180820186556000958652959094209251600290940290920180546001600160a01b03191693909116929092178255519101556105df565b80546001600160a01b0319166001600160a01b038716178155600181018590555b505b83856001600160a01b03167f142e1fdac7ecccbc62af925f0b4039db26847b625602e56b1421dfbc8a0e4f308860405161061b9190610a23565b60405180910390a3505050505050565b60008060008360405161063e9190610bfc565b908152604051908190036020019020549050805b801561071857436000856040516106699190610bfc565b908152604051908190036020019020610683600184610d3d565b8154811061069357610693610d56565b90600052602060002090600202016001015411610706576000846040516106ba9190610bfc565b9081526040519081900360200190206106d4600183610d3d565b815481106106e4576106e4610d56565b60009182526020909120600290910201546001600160a01b0316949350505050565b8061071081610d6c565b915050610652565b5060009392505050565b6002546001600160a01b031633146107685760405162461bcd60e51b81526020600482015260096024820152682737ba1037bbb732b960b91b604482015260640161034f565b6001600160a01b0381166107ad5760405162461bcd60e51b815260206004820152600c60248201526b5a65726f206164647265737360a01b604482015260640161034f565b600280546001600160a01b0319166001600160a01b03831690811790915560405133907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a350565b60606001805480602002602001604051908101604052809291908181526020016000905b828210156108c957838290600052602060002001805461083c90610bc2565b80601f016020809104026020016040519081016040528092919081815260200182805461086890610bc2565b80156108b55780601f1061088a576101008083540402835291602001916108b5565b820191906000526020600020905b81548152906001019060200180831161089857829003601f168201915b50505050508152602001906001019061081d565b50505050905090565b634e487b7160e01b600052604160045260246000fd5b600082601f8301126108f957600080fd5b813567ffffffffffffffff80821115610914576109146108d2565b604051601f8301601f19908116603f0116810190828211818310171561093c5761093c6108d2565b8160405283815286602085880101111561095557600080fd5b836020870160208301376000602085830101528094505050505092915050565b6000806040838503121561098857600080fd5b823567ffffffffffffffff81111561099f57600080fd5b6109ab858286016108e8565b95602094909401359450505050565b6000602082840312156109cc57600080fd5b5035919050565b60005b838110156109ee5781810151838201526020016109d6565b50506000910152565b60008151808452610a0f8160208601602086016109d3565b601f01601f19169290920160200192915050565b602081526000610a3660208301846109f7565b9392505050565b600060208284031215610a4f57600080fd5b813567ffffffffffffffff811115610a6657600080fd5b610a72848285016108e8565b949350505050565b602080825282518282018190526000919060409081850190868401855b82811015610ac557815180516001600160a01b03168552860151868501529284019290850190600101610a97565b5091979650505050505050565b80356001600160a01b0381168114610ae957600080fd5b919050565b600080600060608486031215610b0357600080fd5b833567ffffffffffffffff811115610b1a57600080fd5b610b26868287016108e8565b935050610b3560208501610ad2565b9150604084013590509250925092565b600060208284031215610b5757600080fd5b610a3682610ad2565b6000602080830181845280855180835260408601915060408160051b870101925083870160005b82811015610bb557603f19888603018452610ba38583516109f7565b94509285019290850190600101610b87565b5092979650505050505050565b600181811c90821680610bd657607f821691505b602082108103610bf657634e487b7160e01b600052602260045260246000fd5b50919050565b60008251610c0e8184602087016109d3565b9190910192915050565b601f821115610c6257600081815260208120601f850160051c81016020861015610c3f5750805b601f850160051c820191505b81811015610c5e57828155600101610c4b565b5050505b505050565b815167ffffffffffffffff811115610c8157610c816108d2565b610c9581610c8f8454610bc2565b84610c18565b602080601f831160018114610cca5760008415610cb25750858301515b600019600386901b1c1916600185901b178555610c5e565b600085815260208120601f198616915b82811015610cf957888601518255948401946001909101908401610cda565b5085821015610d175787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b634e487b7160e01b600052601160045260246000fd5b81810381811115610d5057610d50610d27565b92915050565b634e487b7160e01b600052603260045260246000fd5b600081610d7b57610d7b610d27565b50600019019056fea2646970667358221220a3f6e37a5b67f7bb6210f0cfef969cb855bda4c8699bb46ed4a49ee814b7765864736f6c63430008130033`

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
func (_Registry *RegistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
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
func (_Registry *RegistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
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
	var out []interface{}
	err := _Registry.contract.Call(opts, &out, "getActiveAddr", name)
	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err
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

// GetAllNames is a free data retrieval call binding the contract method 0xfb825e5f.
//
// Solidity: function getAllNames() view returns(string[])
func (_Registry *RegistryCaller) GetAllNames(opts *bind.CallOpts) ([]string, error) {
	var out []interface{}
	err := _Registry.contract.Call(opts, &out, "getAllNames")
	if err != nil {
		return *new([]string), err
	}

	out0 := *abi.ConvertType(out[0], new([]string)).(*[]string)

	return out0, err
}

// GetAllNames is a free data retrieval call binding the contract method 0xfb825e5f.
//
// Solidity: function getAllNames() view returns(string[])
func (_Registry *RegistrySession) GetAllNames() ([]string, error) {
	return _Registry.Contract.GetAllNames(&_Registry.CallOpts)
}

// GetAllNames is a free data retrieval call binding the contract method 0xfb825e5f.
//
// Solidity: function getAllNames() view returns(string[])
func (_Registry *RegistryCallerSession) GetAllNames() ([]string, error) {
	return _Registry.Contract.GetAllNames(&_Registry.CallOpts)
}

// GetAllRecords is a free data retrieval call binding the contract method 0x78d573a2.
//
// Solidity: function getAllRecords(string name) view returns((address,uint256)[])
func (_Registry *RegistryCaller) GetAllRecords(opts *bind.CallOpts, name string) ([]IRegistryRecord, error) {
	var out []interface{}
	err := _Registry.contract.Call(opts, &out, "getAllRecords", name)
	if err != nil {
		return *new([]IRegistryRecord), err
	}

	out0 := *abi.ConvertType(out[0], new([]IRegistryRecord)).(*[]IRegistryRecord)

	return out0, err
}

// GetAllRecords is a free data retrieval call binding the contract method 0x78d573a2.
//
// Solidity: function getAllRecords(string name) view returns((address,uint256)[])
func (_Registry *RegistrySession) GetAllRecords(name string) ([]IRegistryRecord, error) {
	return _Registry.Contract.GetAllRecords(&_Registry.CallOpts, name)
}

// GetAllRecords is a free data retrieval call binding the contract method 0x78d573a2.
//
// Solidity: function getAllRecords(string name) view returns((address,uint256)[])
func (_Registry *RegistryCallerSession) GetAllRecords(name string) ([]IRegistryRecord, error) {
	return _Registry.Contract.GetAllRecords(&_Registry.CallOpts, name)
}

// Names is a free data retrieval call binding the contract method 0x4622ab03.
//
// Solidity: function names(uint256 ) view returns(string)
func (_Registry *RegistryCaller) Names(opts *bind.CallOpts, arg0 *big.Int) (string, error) {
	var out []interface{}
	err := _Registry.contract.Call(opts, &out, "names", arg0)
	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err
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

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Registry *RegistryCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Registry.contract.Call(opts, &out, "owner")
	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Registry *RegistrySession) Owner() (common.Address, error) {
	return _Registry.Contract.Owner(&_Registry.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Registry *RegistryCallerSession) Owner() (common.Address, error) {
	return _Registry.Contract.Owner(&_Registry.CallOpts)
}

// Records is a free data retrieval call binding the contract method 0x3b51650d.
//
// Solidity: function records(string , uint256 ) view returns(address addr, uint256 activation)
func (_Registry *RegistryCaller) Records(opts *bind.CallOpts, arg0 string, arg1 *big.Int) (struct {
	Addr       common.Address
	Activation *big.Int
}, error,
) {
	var out []interface{}
	err := _Registry.contract.Call(opts, &out, "records", arg0, arg1)

	outstruct := new(struct {
		Addr       common.Address
		Activation *big.Int
	})

	outstruct.Addr = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.Activation = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	return *outstruct, err
}

// Records is a free data retrieval call binding the contract method 0x3b51650d.
//
// Solidity: function records(string , uint256 ) view returns(address addr, uint256 activation)
func (_Registry *RegistrySession) Records(arg0 string, arg1 *big.Int) (struct {
	Addr       common.Address
	Activation *big.Int
}, error,
) {
	return _Registry.Contract.Records(&_Registry.CallOpts, arg0, arg1)
}

// Records is a free data retrieval call binding the contract method 0x3b51650d.
//
// Solidity: function records(string , uint256 ) view returns(address addr, uint256 activation)
func (_Registry *RegistryCallerSession) Records(arg0 string, arg1 *big.Int) (struct {
	Addr       common.Address
	Activation *big.Int
}, error,
) {
	return _Registry.Contract.Records(&_Registry.CallOpts, arg0, arg1)
}

// Register is a paid mutator transaction binding the contract method 0xd393c871.
//
// Solidity: function register(string name, address addr, uint256 activation) returns()
func (_Registry *RegistryTransactor) Register(opts *bind.TransactOpts, name string, addr common.Address, activation *big.Int) (*types.Transaction, error) {
	return _Registry.contract.Transact(opts, "register", name, addr, activation)
}

// Register is a paid mutator transaction binding the contract method 0xd393c871.
//
// Solidity: function register(string name, address addr, uint256 activation) returns()
func (_Registry *RegistrySession) Register(name string, addr common.Address, activation *big.Int) (*types.Transaction, error) {
	return _Registry.Contract.Register(&_Registry.TransactOpts, name, addr, activation)
}

// Register is a paid mutator transaction binding the contract method 0xd393c871.
//
// Solidity: function register(string name, address addr, uint256 activation) returns()
func (_Registry *RegistryTransactorSession) Register(name string, addr common.Address, activation *big.Int) (*types.Transaction, error) {
	return _Registry.Contract.Register(&_Registry.TransactOpts, name, addr, activation)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Registry *RegistryTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Registry.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Registry *RegistrySession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Registry.Contract.TransferOwnership(&_Registry.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Registry *RegistryTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Registry.Contract.TransferOwnership(&_Registry.TransactOpts, newOwner)
}

// RegistryOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Registry contract.
type RegistryOwnershipTransferredIterator struct {
	Event *RegistryOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log      // Log channel receiving the found contract events
	sub  klaytn.Subscription // Subscription for errors, completion and termination
	done bool                // Whether the subscription completed delivering logs
	fail error               // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *RegistryOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RegistryOwnershipTransferred)
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
		it.Event = new(RegistryOwnershipTransferred)
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
func (it *RegistryOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RegistryOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RegistryOwnershipTransferred represents a OwnershipTransferred event raised by the Registry contract.
type RegistryOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Registry *RegistryFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*RegistryOwnershipTransferredIterator, error) {
	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Registry.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &RegistryOwnershipTransferredIterator{contract: _Registry.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Registry *RegistryFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *RegistryOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {
	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Registry.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RegistryOwnershipTransferred)
				if err := _Registry.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Registry *RegistryFilterer) ParseOwnershipTransferred(log types.Log) (*RegistryOwnershipTransferred, error) {
	event := new(RegistryOwnershipTransferred)
	if err := _Registry.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	return event, nil
}

// RegistryRegisteredIterator is returned from FilterRegistered and is used to iterate over the raw logs and unpacked data for Registered events raised by the Registry contract.
type RegistryRegisteredIterator struct {
	Event *RegistryRegistered // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log      // Log channel receiving the found contract events
	sub  klaytn.Subscription // Subscription for errors, completion and termination
	done bool                // Whether the subscription completed delivering logs
	fail error               // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *RegistryRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RegistryRegistered)
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
		it.Event = new(RegistryRegistered)
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
func (it *RegistryRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RegistryRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RegistryRegistered represents a Registered event raised by the Registry contract.
type RegistryRegistered struct {
	Name       string
	Addr       common.Address
	Activation *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterRegistered is a free log retrieval operation binding the contract event 0x142e1fdac7ecccbc62af925f0b4039db26847b625602e56b1421dfbc8a0e4f30.
//
// Solidity: event Registered(string name, address indexed addr, uint256 indexed activation)
func (_Registry *RegistryFilterer) FilterRegistered(opts *bind.FilterOpts, addr []common.Address, activation []*big.Int) (*RegistryRegisteredIterator, error) {
	var addrRule []interface{}
	for _, addrItem := range addr {
		addrRule = append(addrRule, addrItem)
	}
	var activationRule []interface{}
	for _, activationItem := range activation {
		activationRule = append(activationRule, activationItem)
	}

	logs, sub, err := _Registry.contract.FilterLogs(opts, "Registered", addrRule, activationRule)
	if err != nil {
		return nil, err
	}
	return &RegistryRegisteredIterator{contract: _Registry.contract, event: "Registered", logs: logs, sub: sub}, nil
}

// WatchRegistered is a free log subscription operation binding the contract event 0x142e1fdac7ecccbc62af925f0b4039db26847b625602e56b1421dfbc8a0e4f30.
//
// Solidity: event Registered(string name, address indexed addr, uint256 indexed activation)
func (_Registry *RegistryFilterer) WatchRegistered(opts *bind.WatchOpts, sink chan<- *RegistryRegistered, addr []common.Address, activation []*big.Int) (event.Subscription, error) {
	var addrRule []interface{}
	for _, addrItem := range addr {
		addrRule = append(addrRule, addrItem)
	}
	var activationRule []interface{}
	for _, activationItem := range activation {
		activationRule = append(activationRule, activationItem)
	}

	logs, sub, err := _Registry.contract.WatchLogs(opts, "Registered", addrRule, activationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RegistryRegistered)
				if err := _Registry.contract.UnpackLog(event, "Registered", log); err != nil {
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

// ParseRegistered is a log parse operation binding the contract event 0x142e1fdac7ecccbc62af925f0b4039db26847b625602e56b1421dfbc8a0e4f30.
//
// Solidity: event Registered(string name, address indexed addr, uint256 indexed activation)
func (_Registry *RegistryFilterer) ParseRegistered(log types.Log) (*RegistryRegistered, error) {
	event := new(RegistryRegistered)
	if err := _Registry.contract.UnpackLog(event, "Registered", log); err != nil {
		return nil, err
	}
	return event, nil
}

// RegistryMockMetaData contains all meta data concerning the RegistryMock contract.
var RegistryMockMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"activation\",\"type\":\"uint256\"}],\"name\":\"Registered\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"name\":\"getActiveAddr\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllNames\",\"outputs\":[{\"internalType\":\"string[]\",\"name\":\"\",\"type\":\"string[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"name\":\"getAllRecords\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"activation\",\"type\":\"uint256\"}],\"internalType\":\"structIRegistry.Record[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"names\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"records\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"activation\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"activation\",\"type\":\"uint256\"}],\"name\":\"register\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"e2693e3f": "getActiveAddr(string)",
		"fb825e5f": "getAllNames()",
		"78d573a2": "getAllRecords(string)",
		"4622ab03": "names(uint256)",
		"8da5cb5b": "owner()",
		"3b51650d": "records(string,uint256)",
		"d393c871": "register(string,address,uint256)",
		"f2fde38b": "transferOwnership(address)",
	},
	Bin: "0x608060405234801561001057600080fd5b50610a30806100206000396000f3fe608060405234801561001057600080fd5b50600436106100885760003560e01c8063d393c8711161005b578063d393c87114610129578063e2693e3f1461013e578063f2fde38b14610151578063fb825e5f1461018157600080fd5b80633b51650d1461008d5780634622ab03146100c457806378d573a2146100e45780638da5cb5b14610104575b600080fd5b6100a061009b366004610611565b610196565b604080516001600160a01b0390931683526020830191909152015b60405180910390f35b6100d76100d2366004610656565b6101eb565b6040516100bb91906106bf565b6100f76100f23660046106d9565b610297565b6040516100bb9190610716565b6002546001600160a01b03165b6040516001600160a01b0390911681526020016100bb565b61013c61013736600461078a565b61032a565b005b61011161014c3660046106d9565b6103fd565b61013c61015f3660046107e1565b600280546001600160a01b0319166001600160a01b0392909216919091179055565b610189610495565b6040516100bb91906107fc565b815160208184018101805160008252928201918501919091209190528054829081106101c157600080fd5b6000918252602090912060029091020180546001909101546001600160a01b039091169250905082565b600181815481106101fb57600080fd5b9060005260206000200160009150905080546102169061085e565b80601f01602080910402602001604051908101604052809291908181526020018280546102429061085e565b801561028f5780601f106102645761010080835404028352916020019161028f565b820191906000526020600020905b81548152906001019060200180831161027257829003601f168201915b505050505081565b60606000826040516102a99190610892565b9081526020016040518091039020805480602002602001604051908101604052809291908181526020016000905b8282101561031f576000848152602090819020604080518082019091526002850290910180546001600160a01b031682526001908101548284015290835290920191016102d7565b505050509050919050565b60008360405161033a9190610892565b9081526040519081900360200190205460000361038e576001805480820182556000919091527fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf60161038c84826108fd565b505b60008360405161039e9190610892565b90815260408051602092819003830181208183019092526001600160a01b039485168152828101938452815460018082018455600093845293909220905160029092020180546001600160a01b03191691909416178355905191015550565b6000806000836040516104109190610892565b90815260405190819003602001902054905060008190036104345750600092915050565b6000836040516104449190610892565b90815260405190819003602001902061045e6001836109bd565b8154811061046e5761046e6109e4565b60009182526020909120600290910201546001600160a01b03169392505050565b50919050565b60606001805480602002602001604051908101604052809291908181526020016000905b828210156105655783829060005260206000200180546104d89061085e565b80601f01602080910402602001604051908101604052809291908181526020018280546105049061085e565b80156105515780601f1061052657610100808354040283529160200191610551565b820191906000526020600020905b81548152906001019060200180831161053457829003601f168201915b5050505050815260200190600101906104b9565b50505050905090565b634e487b7160e01b600052604160045260246000fd5b600082601f83011261059557600080fd5b813567ffffffffffffffff808211156105b0576105b061056e565b604051601f8301601f19908116603f011681019082821181831017156105d8576105d861056e565b816040528381528660208588010111156105f157600080fd5b836020870160208301376000602085830101528094505050505092915050565b6000806040838503121561062457600080fd5b823567ffffffffffffffff81111561063b57600080fd5b61064785828601610584565b95602094909401359450505050565b60006020828403121561066857600080fd5b5035919050565b60005b8381101561068a578181015183820152602001610672565b50506000910152565b600081518084526106ab81602086016020860161066f565b601f01601f19169290920160200192915050565b6020815260006106d26020830184610693565b9392505050565b6000602082840312156106eb57600080fd5b813567ffffffffffffffff81111561070257600080fd5b61070e84828501610584565b949350505050565b602080825282518282018190526000919060409081850190868401855b8281101561076157815180516001600160a01b03168552860151868501529284019290850190600101610733565b5091979650505050505050565b80356001600160a01b038116811461078557600080fd5b919050565b60008060006060848603121561079f57600080fd5b833567ffffffffffffffff8111156107b657600080fd5b6107c286828701610584565b9350506107d16020850161076e565b9150604084013590509250925092565b6000602082840312156107f357600080fd5b6106d28261076e565b6000602080830181845280855180835260408601915060408160051b870101925083870160005b8281101561085157603f1988860301845261083f858351610693565b94509285019290850190600101610823565b5092979650505050505050565b600181811c9082168061087257607f821691505b60208210810361048f57634e487b7160e01b600052602260045260246000fd5b600082516108a481846020870161066f565b9190910192915050565b601f8211156108f857600081815260208120601f850160051c810160208610156108d55750805b601f850160051c820191505b818110156108f4578281556001016108e1565b5050505b505050565b815167ffffffffffffffff8111156109175761091761056e565b61092b81610925845461085e565b846108ae565b602080601f83116001811461096057600084156109485750858301515b600019600386901b1c1916600185901b1785556108f4565b600085815260208120601f198616915b8281101561098f57888601518255948401946001909101908401610970565b50858210156109ad5787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b818103818111156109de57634e487b7160e01b600052601160045260246000fd5b92915050565b634e487b7160e01b600052603260045260246000fdfea2646970667358221220055dcb65f93dbbcf001a6a6f0c053ff5f256e00e0debe154e666efbf0caa626364736f6c63430008130033",
}

// RegistryMockABI is the input ABI used to generate the binding from.
// Deprecated: Use RegistryMockMetaData.ABI instead.
var RegistryMockABI = RegistryMockMetaData.ABI

// RegistryMockBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const RegistryMockBinRuntime = `608060405234801561001057600080fd5b50600436106100885760003560e01c8063d393c8711161005b578063d393c87114610129578063e2693e3f1461013e578063f2fde38b14610151578063fb825e5f1461018157600080fd5b80633b51650d1461008d5780634622ab03146100c457806378d573a2146100e45780638da5cb5b14610104575b600080fd5b6100a061009b366004610611565b610196565b604080516001600160a01b0390931683526020830191909152015b60405180910390f35b6100d76100d2366004610656565b6101eb565b6040516100bb91906106bf565b6100f76100f23660046106d9565b610297565b6040516100bb9190610716565b6002546001600160a01b03165b6040516001600160a01b0390911681526020016100bb565b61013c61013736600461078a565b61032a565b005b61011161014c3660046106d9565b6103fd565b61013c61015f3660046107e1565b600280546001600160a01b0319166001600160a01b0392909216919091179055565b610189610495565b6040516100bb91906107fc565b815160208184018101805160008252928201918501919091209190528054829081106101c157600080fd5b6000918252602090912060029091020180546001909101546001600160a01b039091169250905082565b600181815481106101fb57600080fd5b9060005260206000200160009150905080546102169061085e565b80601f01602080910402602001604051908101604052809291908181526020018280546102429061085e565b801561028f5780601f106102645761010080835404028352916020019161028f565b820191906000526020600020905b81548152906001019060200180831161027257829003601f168201915b505050505081565b60606000826040516102a99190610892565b9081526020016040518091039020805480602002602001604051908101604052809291908181526020016000905b8282101561031f576000848152602090819020604080518082019091526002850290910180546001600160a01b031682526001908101548284015290835290920191016102d7565b505050509050919050565b60008360405161033a9190610892565b9081526040519081900360200190205460000361038e576001805480820182556000919091527fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf60161038c84826108fd565b505b60008360405161039e9190610892565b90815260408051602092819003830181208183019092526001600160a01b039485168152828101938452815460018082018455600093845293909220905160029092020180546001600160a01b03191691909416178355905191015550565b6000806000836040516104109190610892565b90815260405190819003602001902054905060008190036104345750600092915050565b6000836040516104449190610892565b90815260405190819003602001902061045e6001836109bd565b8154811061046e5761046e6109e4565b60009182526020909120600290910201546001600160a01b03169392505050565b50919050565b60606001805480602002602001604051908101604052809291908181526020016000905b828210156105655783829060005260206000200180546104d89061085e565b80601f01602080910402602001604051908101604052809291908181526020018280546105049061085e565b80156105515780601f1061052657610100808354040283529160200191610551565b820191906000526020600020905b81548152906001019060200180831161053457829003601f168201915b5050505050815260200190600101906104b9565b50505050905090565b634e487b7160e01b600052604160045260246000fd5b600082601f83011261059557600080fd5b813567ffffffffffffffff808211156105b0576105b061056e565b604051601f8301601f19908116603f011681019082821181831017156105d8576105d861056e565b816040528381528660208588010111156105f157600080fd5b836020870160208301376000602085830101528094505050505092915050565b6000806040838503121561062457600080fd5b823567ffffffffffffffff81111561063b57600080fd5b61064785828601610584565b95602094909401359450505050565b60006020828403121561066857600080fd5b5035919050565b60005b8381101561068a578181015183820152602001610672565b50506000910152565b600081518084526106ab81602086016020860161066f565b601f01601f19169290920160200192915050565b6020815260006106d26020830184610693565b9392505050565b6000602082840312156106eb57600080fd5b813567ffffffffffffffff81111561070257600080fd5b61070e84828501610584565b949350505050565b602080825282518282018190526000919060409081850190868401855b8281101561076157815180516001600160a01b03168552860151868501529284019290850190600101610733565b5091979650505050505050565b80356001600160a01b038116811461078557600080fd5b919050565b60008060006060848603121561079f57600080fd5b833567ffffffffffffffff8111156107b657600080fd5b6107c286828701610584565b9350506107d16020850161076e565b9150604084013590509250925092565b6000602082840312156107f357600080fd5b6106d28261076e565b6000602080830181845280855180835260408601915060408160051b870101925083870160005b8281101561085157603f1988860301845261083f858351610693565b94509285019290850190600101610823565b5092979650505050505050565b600181811c9082168061087257607f821691505b60208210810361048f57634e487b7160e01b600052602260045260246000fd5b600082516108a481846020870161066f565b9190910192915050565b601f8211156108f857600081815260208120601f850160051c810160208610156108d55750805b601f850160051c820191505b818110156108f4578281556001016108e1565b5050505b505050565b815167ffffffffffffffff8111156109175761091761056e565b61092b81610925845461085e565b846108ae565b602080601f83116001811461096057600084156109485750858301515b600019600386901b1c1916600185901b1785556108f4565b600085815260208120601f198616915b8281101561098f57888601518255948401946001909101908401610970565b50858210156109ad5787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b818103818111156109de57634e487b7160e01b600052601160045260246000fd5b92915050565b634e487b7160e01b600052603260045260246000fdfea2646970667358221220055dcb65f93dbbcf001a6a6f0c053ff5f256e00e0debe154e666efbf0caa626364736f6c63430008130033`

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
func (_RegistryMock *RegistryMockRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
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
func (_RegistryMock *RegistryMockCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
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
	var out []interface{}
	err := _RegistryMock.contract.Call(opts, &out, "getActiveAddr", name)
	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err
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

// GetAllNames is a free data retrieval call binding the contract method 0xfb825e5f.
//
// Solidity: function getAllNames() view returns(string[])
func (_RegistryMock *RegistryMockCaller) GetAllNames(opts *bind.CallOpts) ([]string, error) {
	var out []interface{}
	err := _RegistryMock.contract.Call(opts, &out, "getAllNames")
	if err != nil {
		return *new([]string), err
	}

	out0 := *abi.ConvertType(out[0], new([]string)).(*[]string)

	return out0, err
}

// GetAllNames is a free data retrieval call binding the contract method 0xfb825e5f.
//
// Solidity: function getAllNames() view returns(string[])
func (_RegistryMock *RegistryMockSession) GetAllNames() ([]string, error) {
	return _RegistryMock.Contract.GetAllNames(&_RegistryMock.CallOpts)
}

// GetAllNames is a free data retrieval call binding the contract method 0xfb825e5f.
//
// Solidity: function getAllNames() view returns(string[])
func (_RegistryMock *RegistryMockCallerSession) GetAllNames() ([]string, error) {
	return _RegistryMock.Contract.GetAllNames(&_RegistryMock.CallOpts)
}

// GetAllRecords is a free data retrieval call binding the contract method 0x78d573a2.
//
// Solidity: function getAllRecords(string name) view returns((address,uint256)[])
func (_RegistryMock *RegistryMockCaller) GetAllRecords(opts *bind.CallOpts, name string) ([]IRegistryRecord, error) {
	var out []interface{}
	err := _RegistryMock.contract.Call(opts, &out, "getAllRecords", name)
	if err != nil {
		return *new([]IRegistryRecord), err
	}

	out0 := *abi.ConvertType(out[0], new([]IRegistryRecord)).(*[]IRegistryRecord)

	return out0, err
}

// GetAllRecords is a free data retrieval call binding the contract method 0x78d573a2.
//
// Solidity: function getAllRecords(string name) view returns((address,uint256)[])
func (_RegistryMock *RegistryMockSession) GetAllRecords(name string) ([]IRegistryRecord, error) {
	return _RegistryMock.Contract.GetAllRecords(&_RegistryMock.CallOpts, name)
}

// GetAllRecords is a free data retrieval call binding the contract method 0x78d573a2.
//
// Solidity: function getAllRecords(string name) view returns((address,uint256)[])
func (_RegistryMock *RegistryMockCallerSession) GetAllRecords(name string) ([]IRegistryRecord, error) {
	return _RegistryMock.Contract.GetAllRecords(&_RegistryMock.CallOpts, name)
}

// Names is a free data retrieval call binding the contract method 0x4622ab03.
//
// Solidity: function names(uint256 ) view returns(string)
func (_RegistryMock *RegistryMockCaller) Names(opts *bind.CallOpts, arg0 *big.Int) (string, error) {
	var out []interface{}
	err := _RegistryMock.contract.Call(opts, &out, "names", arg0)
	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err
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

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_RegistryMock *RegistryMockCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _RegistryMock.contract.Call(opts, &out, "owner")
	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_RegistryMock *RegistryMockSession) Owner() (common.Address, error) {
	return _RegistryMock.Contract.Owner(&_RegistryMock.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_RegistryMock *RegistryMockCallerSession) Owner() (common.Address, error) {
	return _RegistryMock.Contract.Owner(&_RegistryMock.CallOpts)
}

// Records is a free data retrieval call binding the contract method 0x3b51650d.
//
// Solidity: function records(string , uint256 ) view returns(address addr, uint256 activation)
func (_RegistryMock *RegistryMockCaller) Records(opts *bind.CallOpts, arg0 string, arg1 *big.Int) (struct {
	Addr       common.Address
	Activation *big.Int
}, error,
) {
	var out []interface{}
	err := _RegistryMock.contract.Call(opts, &out, "records", arg0, arg1)

	outstruct := new(struct {
		Addr       common.Address
		Activation *big.Int
	})

	outstruct.Addr = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.Activation = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	return *outstruct, err
}

// Records is a free data retrieval call binding the contract method 0x3b51650d.
//
// Solidity: function records(string , uint256 ) view returns(address addr, uint256 activation)
func (_RegistryMock *RegistryMockSession) Records(arg0 string, arg1 *big.Int) (struct {
	Addr       common.Address
	Activation *big.Int
}, error,
) {
	return _RegistryMock.Contract.Records(&_RegistryMock.CallOpts, arg0, arg1)
}

// Records is a free data retrieval call binding the contract method 0x3b51650d.
//
// Solidity: function records(string , uint256 ) view returns(address addr, uint256 activation)
func (_RegistryMock *RegistryMockCallerSession) Records(arg0 string, arg1 *big.Int) (struct {
	Addr       common.Address
	Activation *big.Int
}, error,
) {
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

// RegistryMockOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the RegistryMock contract.
type RegistryMockOwnershipTransferredIterator struct {
	Event *RegistryMockOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log      // Log channel receiving the found contract events
	sub  klaytn.Subscription // Subscription for errors, completion and termination
	done bool                // Whether the subscription completed delivering logs
	fail error               // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *RegistryMockOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RegistryMockOwnershipTransferred)
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
		it.Event = new(RegistryMockOwnershipTransferred)
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
func (it *RegistryMockOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RegistryMockOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RegistryMockOwnershipTransferred represents a OwnershipTransferred event raised by the RegistryMock contract.
type RegistryMockOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_RegistryMock *RegistryMockFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*RegistryMockOwnershipTransferredIterator, error) {
	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _RegistryMock.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &RegistryMockOwnershipTransferredIterator{contract: _RegistryMock.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_RegistryMock *RegistryMockFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *RegistryMockOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {
	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _RegistryMock.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RegistryMockOwnershipTransferred)
				if err := _RegistryMock.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_RegistryMock *RegistryMockFilterer) ParseOwnershipTransferred(log types.Log) (*RegistryMockOwnershipTransferred, error) {
	event := new(RegistryMockOwnershipTransferred)
	if err := _RegistryMock.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	return event, nil
}

// RegistryMockRegisteredIterator is returned from FilterRegistered and is used to iterate over the raw logs and unpacked data for Registered events raised by the RegistryMock contract.
type RegistryMockRegisteredIterator struct {
	Event *RegistryMockRegistered // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log      // Log channel receiving the found contract events
	sub  klaytn.Subscription // Subscription for errors, completion and termination
	done bool                // Whether the subscription completed delivering logs
	fail error               // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *RegistryMockRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RegistryMockRegistered)
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
		it.Event = new(RegistryMockRegistered)
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
func (it *RegistryMockRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RegistryMockRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RegistryMockRegistered represents a Registered event raised by the RegistryMock contract.
type RegistryMockRegistered struct {
	Name       string
	Addr       common.Address
	Activation *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterRegistered is a free log retrieval operation binding the contract event 0x142e1fdac7ecccbc62af925f0b4039db26847b625602e56b1421dfbc8a0e4f30.
//
// Solidity: event Registered(string name, address indexed addr, uint256 indexed activation)
func (_RegistryMock *RegistryMockFilterer) FilterRegistered(opts *bind.FilterOpts, addr []common.Address, activation []*big.Int) (*RegistryMockRegisteredIterator, error) {
	var addrRule []interface{}
	for _, addrItem := range addr {
		addrRule = append(addrRule, addrItem)
	}
	var activationRule []interface{}
	for _, activationItem := range activation {
		activationRule = append(activationRule, activationItem)
	}

	logs, sub, err := _RegistryMock.contract.FilterLogs(opts, "Registered", addrRule, activationRule)
	if err != nil {
		return nil, err
	}
	return &RegistryMockRegisteredIterator{contract: _RegistryMock.contract, event: "Registered", logs: logs, sub: sub}, nil
}

// WatchRegistered is a free log subscription operation binding the contract event 0x142e1fdac7ecccbc62af925f0b4039db26847b625602e56b1421dfbc8a0e4f30.
//
// Solidity: event Registered(string name, address indexed addr, uint256 indexed activation)
func (_RegistryMock *RegistryMockFilterer) WatchRegistered(opts *bind.WatchOpts, sink chan<- *RegistryMockRegistered, addr []common.Address, activation []*big.Int) (event.Subscription, error) {
	var addrRule []interface{}
	for _, addrItem := range addr {
		addrRule = append(addrRule, addrItem)
	}
	var activationRule []interface{}
	for _, activationItem := range activation {
		activationRule = append(activationRule, activationItem)
	}

	logs, sub, err := _RegistryMock.contract.WatchLogs(opts, "Registered", addrRule, activationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RegistryMockRegistered)
				if err := _RegistryMock.contract.UnpackLog(event, "Registered", log); err != nil {
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

// ParseRegistered is a log parse operation binding the contract event 0x142e1fdac7ecccbc62af925f0b4039db26847b625602e56b1421dfbc8a0e4f30.
//
// Solidity: event Registered(string name, address indexed addr, uint256 indexed activation)
func (_RegistryMock *RegistryMockFilterer) ParseRegistered(log types.Log) (*RegistryMockRegistered, error) {
	event := new(RegistryMockRegistered)
	if err := _RegistryMock.contract.UnpackLog(event, "Registered", log); err != nil {
		return nil, err
	}
	return event, nil
}

// StorageSlotMetaData contains all meta data concerning the StorageSlot contract.
var StorageSlotMetaData = &bind.MetaData{
	ABI: "[]",
	Bin: "0x60566037600b82828239805160001a607314602a57634e487b7160e01b600052600060045260246000fd5b30600052607381538281f3fe73000000000000000000000000000000000000000030146080604052600080fdfea264697066735822122061ca98bd9a9672bf10c068dcf42a45d13fc17620bb619c26bdf064388a883db964736f6c63430008130033",
}

// StorageSlotABI is the input ABI used to generate the binding from.
// Deprecated: Use StorageSlotMetaData.ABI instead.
var StorageSlotABI = StorageSlotMetaData.ABI

// StorageSlotBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const StorageSlotBinRuntime = `73000000000000000000000000000000000000000030146080604052600080fdfea264697066735822122061ca98bd9a9672bf10c068dcf42a45d13fc17620bb619c26bdf064388a883db964736f6c63430008130033`

// StorageSlotBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use StorageSlotMetaData.Bin instead.
var StorageSlotBin = StorageSlotMetaData.Bin

// DeployStorageSlot deploys a new Klaytn contract, binding an instance of StorageSlot to it.
func DeployStorageSlot(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *StorageSlot, error) {
	parsed, err := StorageSlotMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(StorageSlotBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &StorageSlot{StorageSlotCaller: StorageSlotCaller{contract: contract}, StorageSlotTransactor: StorageSlotTransactor{contract: contract}, StorageSlotFilterer: StorageSlotFilterer{contract: contract}}, nil
}

// StorageSlot is an auto generated Go binding around a Klaytn contract.
type StorageSlot struct {
	StorageSlotCaller     // Read-only binding to the contract
	StorageSlotTransactor // Write-only binding to the contract
	StorageSlotFilterer   // Log filterer for contract events
}

// StorageSlotCaller is an auto generated read-only Go binding around a Klaytn contract.
type StorageSlotCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StorageSlotTransactor is an auto generated write-only Go binding around a Klaytn contract.
type StorageSlotTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StorageSlotFilterer is an auto generated log filtering Go binding around a Klaytn contract events.
type StorageSlotFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StorageSlotSession is an auto generated Go binding around a Klaytn contract,
// with pre-set call and transact options.
type StorageSlotSession struct {
	Contract     *StorageSlot      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// StorageSlotCallerSession is an auto generated read-only Go binding around a Klaytn contract,
// with pre-set call options.
type StorageSlotCallerSession struct {
	Contract *StorageSlotCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// StorageSlotTransactorSession is an auto generated write-only Go binding around a Klaytn contract,
// with pre-set transact options.
type StorageSlotTransactorSession struct {
	Contract     *StorageSlotTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// StorageSlotRaw is an auto generated low-level Go binding around a Klaytn contract.
type StorageSlotRaw struct {
	Contract *StorageSlot // Generic contract binding to access the raw methods on
}

// StorageSlotCallerRaw is an auto generated low-level read-only Go binding around a Klaytn contract.
type StorageSlotCallerRaw struct {
	Contract *StorageSlotCaller // Generic read-only contract binding to access the raw methods on
}

// StorageSlotTransactorRaw is an auto generated low-level write-only Go binding around a Klaytn contract.
type StorageSlotTransactorRaw struct {
	Contract *StorageSlotTransactor // Generic write-only contract binding to access the raw methods on
}

// NewStorageSlot creates a new instance of StorageSlot, bound to a specific deployed contract.
func NewStorageSlot(address common.Address, backend bind.ContractBackend) (*StorageSlot, error) {
	contract, err := bindStorageSlot(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &StorageSlot{StorageSlotCaller: StorageSlotCaller{contract: contract}, StorageSlotTransactor: StorageSlotTransactor{contract: contract}, StorageSlotFilterer: StorageSlotFilterer{contract: contract}}, nil
}

// NewStorageSlotCaller creates a new read-only instance of StorageSlot, bound to a specific deployed contract.
func NewStorageSlotCaller(address common.Address, caller bind.ContractCaller) (*StorageSlotCaller, error) {
	contract, err := bindStorageSlot(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &StorageSlotCaller{contract: contract}, nil
}

// NewStorageSlotTransactor creates a new write-only instance of StorageSlot, bound to a specific deployed contract.
func NewStorageSlotTransactor(address common.Address, transactor bind.ContractTransactor) (*StorageSlotTransactor, error) {
	contract, err := bindStorageSlot(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &StorageSlotTransactor{contract: contract}, nil
}

// NewStorageSlotFilterer creates a new log filterer instance of StorageSlot, bound to a specific deployed contract.
func NewStorageSlotFilterer(address common.Address, filterer bind.ContractFilterer) (*StorageSlotFilterer, error) {
	contract, err := bindStorageSlot(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &StorageSlotFilterer{contract: contract}, nil
}

// bindStorageSlot binds a generic wrapper to an already deployed contract.
func bindStorageSlot(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := StorageSlotMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_StorageSlot *StorageSlotRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _StorageSlot.Contract.StorageSlotCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_StorageSlot *StorageSlotRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StorageSlot.Contract.StorageSlotTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_StorageSlot *StorageSlotRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _StorageSlot.Contract.StorageSlotTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_StorageSlot *StorageSlotCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _StorageSlot.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_StorageSlot *StorageSlotTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StorageSlot.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_StorageSlot *StorageSlotTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _StorageSlot.Contract.contract.Transact(opts, method, params...)
}

// StorageSlotUpgradeableMetaData contains all meta data concerning the StorageSlotUpgradeable contract.
var StorageSlotUpgradeableMetaData = &bind.MetaData{
	ABI: "[]",
	Bin: "0x60566037600b82828239805160001a607314602a57634e487b7160e01b600052600060045260246000fd5b30600052607381538281f3fe73000000000000000000000000000000000000000030146080604052600080fdfea2646970667358221220591c77dfec17118a2ca0803e321cdaa6629b0a43ce207b858ac44c985568473064736f6c63430008130033",
}

// StorageSlotUpgradeableABI is the input ABI used to generate the binding from.
// Deprecated: Use StorageSlotUpgradeableMetaData.ABI instead.
var StorageSlotUpgradeableABI = StorageSlotUpgradeableMetaData.ABI

// StorageSlotUpgradeableBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const StorageSlotUpgradeableBinRuntime = `73000000000000000000000000000000000000000030146080604052600080fdfea2646970667358221220591c77dfec17118a2ca0803e321cdaa6629b0a43ce207b858ac44c985568473064736f6c63430008130033`

// StorageSlotUpgradeableBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use StorageSlotUpgradeableMetaData.Bin instead.
var StorageSlotUpgradeableBin = StorageSlotUpgradeableMetaData.Bin

// DeployStorageSlotUpgradeable deploys a new Klaytn contract, binding an instance of StorageSlotUpgradeable to it.
func DeployStorageSlotUpgradeable(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *StorageSlotUpgradeable, error) {
	parsed, err := StorageSlotUpgradeableMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(StorageSlotUpgradeableBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &StorageSlotUpgradeable{StorageSlotUpgradeableCaller: StorageSlotUpgradeableCaller{contract: contract}, StorageSlotUpgradeableTransactor: StorageSlotUpgradeableTransactor{contract: contract}, StorageSlotUpgradeableFilterer: StorageSlotUpgradeableFilterer{contract: contract}}, nil
}

// StorageSlotUpgradeable is an auto generated Go binding around a Klaytn contract.
type StorageSlotUpgradeable struct {
	StorageSlotUpgradeableCaller     // Read-only binding to the contract
	StorageSlotUpgradeableTransactor // Write-only binding to the contract
	StorageSlotUpgradeableFilterer   // Log filterer for contract events
}

// StorageSlotUpgradeableCaller is an auto generated read-only Go binding around a Klaytn contract.
type StorageSlotUpgradeableCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StorageSlotUpgradeableTransactor is an auto generated write-only Go binding around a Klaytn contract.
type StorageSlotUpgradeableTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StorageSlotUpgradeableFilterer is an auto generated log filtering Go binding around a Klaytn contract events.
type StorageSlotUpgradeableFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StorageSlotUpgradeableSession is an auto generated Go binding around a Klaytn contract,
// with pre-set call and transact options.
type StorageSlotUpgradeableSession struct {
	Contract     *StorageSlotUpgradeable // Generic contract binding to set the session for
	CallOpts     bind.CallOpts           // Call options to use throughout this session
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// StorageSlotUpgradeableCallerSession is an auto generated read-only Go binding around a Klaytn contract,
// with pre-set call options.
type StorageSlotUpgradeableCallerSession struct {
	Contract *StorageSlotUpgradeableCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                 // Call options to use throughout this session
}

// StorageSlotUpgradeableTransactorSession is an auto generated write-only Go binding around a Klaytn contract,
// with pre-set transact options.
type StorageSlotUpgradeableTransactorSession struct {
	Contract     *StorageSlotUpgradeableTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                 // Transaction auth options to use throughout this session
}

// StorageSlotUpgradeableRaw is an auto generated low-level Go binding around a Klaytn contract.
type StorageSlotUpgradeableRaw struct {
	Contract *StorageSlotUpgradeable // Generic contract binding to access the raw methods on
}

// StorageSlotUpgradeableCallerRaw is an auto generated low-level read-only Go binding around a Klaytn contract.
type StorageSlotUpgradeableCallerRaw struct {
	Contract *StorageSlotUpgradeableCaller // Generic read-only contract binding to access the raw methods on
}

// StorageSlotUpgradeableTransactorRaw is an auto generated low-level write-only Go binding around a Klaytn contract.
type StorageSlotUpgradeableTransactorRaw struct {
	Contract *StorageSlotUpgradeableTransactor // Generic write-only contract binding to access the raw methods on
}

// NewStorageSlotUpgradeable creates a new instance of StorageSlotUpgradeable, bound to a specific deployed contract.
func NewStorageSlotUpgradeable(address common.Address, backend bind.ContractBackend) (*StorageSlotUpgradeable, error) {
	contract, err := bindStorageSlotUpgradeable(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &StorageSlotUpgradeable{StorageSlotUpgradeableCaller: StorageSlotUpgradeableCaller{contract: contract}, StorageSlotUpgradeableTransactor: StorageSlotUpgradeableTransactor{contract: contract}, StorageSlotUpgradeableFilterer: StorageSlotUpgradeableFilterer{contract: contract}}, nil
}

// NewStorageSlotUpgradeableCaller creates a new read-only instance of StorageSlotUpgradeable, bound to a specific deployed contract.
func NewStorageSlotUpgradeableCaller(address common.Address, caller bind.ContractCaller) (*StorageSlotUpgradeableCaller, error) {
	contract, err := bindStorageSlotUpgradeable(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &StorageSlotUpgradeableCaller{contract: contract}, nil
}

// NewStorageSlotUpgradeableTransactor creates a new write-only instance of StorageSlotUpgradeable, bound to a specific deployed contract.
func NewStorageSlotUpgradeableTransactor(address common.Address, transactor bind.ContractTransactor) (*StorageSlotUpgradeableTransactor, error) {
	contract, err := bindStorageSlotUpgradeable(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &StorageSlotUpgradeableTransactor{contract: contract}, nil
}

// NewStorageSlotUpgradeableFilterer creates a new log filterer instance of StorageSlotUpgradeable, bound to a specific deployed contract.
func NewStorageSlotUpgradeableFilterer(address common.Address, filterer bind.ContractFilterer) (*StorageSlotUpgradeableFilterer, error) {
	contract, err := bindStorageSlotUpgradeable(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &StorageSlotUpgradeableFilterer{contract: contract}, nil
}

// bindStorageSlotUpgradeable binds a generic wrapper to an already deployed contract.
func bindStorageSlotUpgradeable(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := StorageSlotUpgradeableMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_StorageSlotUpgradeable *StorageSlotUpgradeableRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _StorageSlotUpgradeable.Contract.StorageSlotUpgradeableCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_StorageSlotUpgradeable *StorageSlotUpgradeableRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StorageSlotUpgradeable.Contract.StorageSlotUpgradeableTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_StorageSlotUpgradeable *StorageSlotUpgradeableRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _StorageSlotUpgradeable.Contract.StorageSlotUpgradeableTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_StorageSlotUpgradeable *StorageSlotUpgradeableCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _StorageSlotUpgradeable.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_StorageSlotUpgradeable *StorageSlotUpgradeableTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StorageSlotUpgradeable.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_StorageSlotUpgradeable *StorageSlotUpgradeableTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _StorageSlotUpgradeable.Contract.contract.Transact(opts, method, params...)
}

// UUPSUpgradeableMetaData contains all meta data concerning the UUPSUpgradeable contract.
var UUPSUpgradeableMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"previousAdmin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"AdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"beacon\",\"type\":\"address\"}],\"name\":\"BeaconUpgraded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"proxiableUUID\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"}],\"name\":\"upgradeTo\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"upgradeToAndCall\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"52d1902d": "proxiableUUID()",
		"3659cfe6": "upgradeTo(address)",
		"4f1ef286": "upgradeToAndCall(address,bytes)",
	},
}

// UUPSUpgradeableABI is the input ABI used to generate the binding from.
// Deprecated: Use UUPSUpgradeableMetaData.ABI instead.
var UUPSUpgradeableABI = UUPSUpgradeableMetaData.ABI

// UUPSUpgradeableBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const UUPSUpgradeableBinRuntime = ``

// UUPSUpgradeableFuncSigs maps the 4-byte function signature to its string representation.
// Deprecated: Use UUPSUpgradeableMetaData.Sigs instead.
var UUPSUpgradeableFuncSigs = UUPSUpgradeableMetaData.Sigs

// UUPSUpgradeable is an auto generated Go binding around a Klaytn contract.
type UUPSUpgradeable struct {
	UUPSUpgradeableCaller     // Read-only binding to the contract
	UUPSUpgradeableTransactor // Write-only binding to the contract
	UUPSUpgradeableFilterer   // Log filterer for contract events
}

// UUPSUpgradeableCaller is an auto generated read-only Go binding around a Klaytn contract.
type UUPSUpgradeableCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UUPSUpgradeableTransactor is an auto generated write-only Go binding around a Klaytn contract.
type UUPSUpgradeableTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UUPSUpgradeableFilterer is an auto generated log filtering Go binding around a Klaytn contract events.
type UUPSUpgradeableFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UUPSUpgradeableSession is an auto generated Go binding around a Klaytn contract,
// with pre-set call and transact options.
type UUPSUpgradeableSession struct {
	Contract     *UUPSUpgradeable  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// UUPSUpgradeableCallerSession is an auto generated read-only Go binding around a Klaytn contract,
// with pre-set call options.
type UUPSUpgradeableCallerSession struct {
	Contract *UUPSUpgradeableCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// UUPSUpgradeableTransactorSession is an auto generated write-only Go binding around a Klaytn contract,
// with pre-set transact options.
type UUPSUpgradeableTransactorSession struct {
	Contract     *UUPSUpgradeableTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// UUPSUpgradeableRaw is an auto generated low-level Go binding around a Klaytn contract.
type UUPSUpgradeableRaw struct {
	Contract *UUPSUpgradeable // Generic contract binding to access the raw methods on
}

// UUPSUpgradeableCallerRaw is an auto generated low-level read-only Go binding around a Klaytn contract.
type UUPSUpgradeableCallerRaw struct {
	Contract *UUPSUpgradeableCaller // Generic read-only contract binding to access the raw methods on
}

// UUPSUpgradeableTransactorRaw is an auto generated low-level write-only Go binding around a Klaytn contract.
type UUPSUpgradeableTransactorRaw struct {
	Contract *UUPSUpgradeableTransactor // Generic write-only contract binding to access the raw methods on
}

// NewUUPSUpgradeable creates a new instance of UUPSUpgradeable, bound to a specific deployed contract.
func NewUUPSUpgradeable(address common.Address, backend bind.ContractBackend) (*UUPSUpgradeable, error) {
	contract, err := bindUUPSUpgradeable(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &UUPSUpgradeable{UUPSUpgradeableCaller: UUPSUpgradeableCaller{contract: contract}, UUPSUpgradeableTransactor: UUPSUpgradeableTransactor{contract: contract}, UUPSUpgradeableFilterer: UUPSUpgradeableFilterer{contract: contract}}, nil
}

// NewUUPSUpgradeableCaller creates a new read-only instance of UUPSUpgradeable, bound to a specific deployed contract.
func NewUUPSUpgradeableCaller(address common.Address, caller bind.ContractCaller) (*UUPSUpgradeableCaller, error) {
	contract, err := bindUUPSUpgradeable(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &UUPSUpgradeableCaller{contract: contract}, nil
}

// NewUUPSUpgradeableTransactor creates a new write-only instance of UUPSUpgradeable, bound to a specific deployed contract.
func NewUUPSUpgradeableTransactor(address common.Address, transactor bind.ContractTransactor) (*UUPSUpgradeableTransactor, error) {
	contract, err := bindUUPSUpgradeable(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &UUPSUpgradeableTransactor{contract: contract}, nil
}

// NewUUPSUpgradeableFilterer creates a new log filterer instance of UUPSUpgradeable, bound to a specific deployed contract.
func NewUUPSUpgradeableFilterer(address common.Address, filterer bind.ContractFilterer) (*UUPSUpgradeableFilterer, error) {
	contract, err := bindUUPSUpgradeable(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &UUPSUpgradeableFilterer{contract: contract}, nil
}

// bindUUPSUpgradeable binds a generic wrapper to an already deployed contract.
func bindUUPSUpgradeable(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := UUPSUpgradeableMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_UUPSUpgradeable *UUPSUpgradeableRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _UUPSUpgradeable.Contract.UUPSUpgradeableCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_UUPSUpgradeable *UUPSUpgradeableRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UUPSUpgradeable.Contract.UUPSUpgradeableTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_UUPSUpgradeable *UUPSUpgradeableRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _UUPSUpgradeable.Contract.UUPSUpgradeableTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_UUPSUpgradeable *UUPSUpgradeableCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _UUPSUpgradeable.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_UUPSUpgradeable *UUPSUpgradeableTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UUPSUpgradeable.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_UUPSUpgradeable *UUPSUpgradeableTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _UUPSUpgradeable.Contract.contract.Transact(opts, method, params...)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_UUPSUpgradeable *UUPSUpgradeableCaller) ProxiableUUID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _UUPSUpgradeable.contract.Call(opts, &out, "proxiableUUID")
	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_UUPSUpgradeable *UUPSUpgradeableSession) ProxiableUUID() ([32]byte, error) {
	return _UUPSUpgradeable.Contract.ProxiableUUID(&_UUPSUpgradeable.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_UUPSUpgradeable *UUPSUpgradeableCallerSession) ProxiableUUID() ([32]byte, error) {
	return _UUPSUpgradeable.Contract.ProxiableUUID(&_UUPSUpgradeable.CallOpts)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_UUPSUpgradeable *UUPSUpgradeableTransactor) UpgradeTo(opts *bind.TransactOpts, newImplementation common.Address) (*types.Transaction, error) {
	return _UUPSUpgradeable.contract.Transact(opts, "upgradeTo", newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_UUPSUpgradeable *UUPSUpgradeableSession) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _UUPSUpgradeable.Contract.UpgradeTo(&_UUPSUpgradeable.TransactOpts, newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_UUPSUpgradeable *UUPSUpgradeableTransactorSession) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _UUPSUpgradeable.Contract.UpgradeTo(&_UUPSUpgradeable.TransactOpts, newImplementation)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_UUPSUpgradeable *UUPSUpgradeableTransactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _UUPSUpgradeable.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_UUPSUpgradeable *UUPSUpgradeableSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _UUPSUpgradeable.Contract.UpgradeToAndCall(&_UUPSUpgradeable.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_UUPSUpgradeable *UUPSUpgradeableTransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _UUPSUpgradeable.Contract.UpgradeToAndCall(&_UUPSUpgradeable.TransactOpts, newImplementation, data)
}

// UUPSUpgradeableAdminChangedIterator is returned from FilterAdminChanged and is used to iterate over the raw logs and unpacked data for AdminChanged events raised by the UUPSUpgradeable contract.
type UUPSUpgradeableAdminChangedIterator struct {
	Event *UUPSUpgradeableAdminChanged // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log      // Log channel receiving the found contract events
	sub  klaytn.Subscription // Subscription for errors, completion and termination
	done bool                // Whether the subscription completed delivering logs
	fail error               // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *UUPSUpgradeableAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UUPSUpgradeableAdminChanged)
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
		it.Event = new(UUPSUpgradeableAdminChanged)
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
func (it *UUPSUpgradeableAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UUPSUpgradeableAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UUPSUpgradeableAdminChanged represents a AdminChanged event raised by the UUPSUpgradeable contract.
type UUPSUpgradeableAdminChanged struct {
	PreviousAdmin common.Address
	NewAdmin      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterAdminChanged is a free log retrieval operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_UUPSUpgradeable *UUPSUpgradeableFilterer) FilterAdminChanged(opts *bind.FilterOpts) (*UUPSUpgradeableAdminChangedIterator, error) {
	logs, sub, err := _UUPSUpgradeable.contract.FilterLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return &UUPSUpgradeableAdminChangedIterator{contract: _UUPSUpgradeable.contract, event: "AdminChanged", logs: logs, sub: sub}, nil
}

// WatchAdminChanged is a free log subscription operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_UUPSUpgradeable *UUPSUpgradeableFilterer) WatchAdminChanged(opts *bind.WatchOpts, sink chan<- *UUPSUpgradeableAdminChanged) (event.Subscription, error) {
	logs, sub, err := _UUPSUpgradeable.contract.WatchLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UUPSUpgradeableAdminChanged)
				if err := _UUPSUpgradeable.contract.UnpackLog(event, "AdminChanged", log); err != nil {
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

// ParseAdminChanged is a log parse operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_UUPSUpgradeable *UUPSUpgradeableFilterer) ParseAdminChanged(log types.Log) (*UUPSUpgradeableAdminChanged, error) {
	event := new(UUPSUpgradeableAdminChanged)
	if err := _UUPSUpgradeable.contract.UnpackLog(event, "AdminChanged", log); err != nil {
		return nil, err
	}
	return event, nil
}

// UUPSUpgradeableBeaconUpgradedIterator is returned from FilterBeaconUpgraded and is used to iterate over the raw logs and unpacked data for BeaconUpgraded events raised by the UUPSUpgradeable contract.
type UUPSUpgradeableBeaconUpgradedIterator struct {
	Event *UUPSUpgradeableBeaconUpgraded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log      // Log channel receiving the found contract events
	sub  klaytn.Subscription // Subscription for errors, completion and termination
	done bool                // Whether the subscription completed delivering logs
	fail error               // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *UUPSUpgradeableBeaconUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UUPSUpgradeableBeaconUpgraded)
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
		it.Event = new(UUPSUpgradeableBeaconUpgraded)
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
func (it *UUPSUpgradeableBeaconUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UUPSUpgradeableBeaconUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UUPSUpgradeableBeaconUpgraded represents a BeaconUpgraded event raised by the UUPSUpgradeable contract.
type UUPSUpgradeableBeaconUpgraded struct {
	Beacon common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterBeaconUpgraded is a free log retrieval operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_UUPSUpgradeable *UUPSUpgradeableFilterer) FilterBeaconUpgraded(opts *bind.FilterOpts, beacon []common.Address) (*UUPSUpgradeableBeaconUpgradedIterator, error) {
	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _UUPSUpgradeable.contract.FilterLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return &UUPSUpgradeableBeaconUpgradedIterator{contract: _UUPSUpgradeable.contract, event: "BeaconUpgraded", logs: logs, sub: sub}, nil
}

// WatchBeaconUpgraded is a free log subscription operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_UUPSUpgradeable *UUPSUpgradeableFilterer) WatchBeaconUpgraded(opts *bind.WatchOpts, sink chan<- *UUPSUpgradeableBeaconUpgraded, beacon []common.Address) (event.Subscription, error) {
	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _UUPSUpgradeable.contract.WatchLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UUPSUpgradeableBeaconUpgraded)
				if err := _UUPSUpgradeable.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
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

// ParseBeaconUpgraded is a log parse operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_UUPSUpgradeable *UUPSUpgradeableFilterer) ParseBeaconUpgraded(log types.Log) (*UUPSUpgradeableBeaconUpgraded, error) {
	event := new(UUPSUpgradeableBeaconUpgraded)
	if err := _UUPSUpgradeable.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
		return nil, err
	}
	return event, nil
}

// UUPSUpgradeableInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the UUPSUpgradeable contract.
type UUPSUpgradeableInitializedIterator struct {
	Event *UUPSUpgradeableInitialized // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log      // Log channel receiving the found contract events
	sub  klaytn.Subscription // Subscription for errors, completion and termination
	done bool                // Whether the subscription completed delivering logs
	fail error               // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *UUPSUpgradeableInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UUPSUpgradeableInitialized)
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
		it.Event = new(UUPSUpgradeableInitialized)
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
func (it *UUPSUpgradeableInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UUPSUpgradeableInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UUPSUpgradeableInitialized represents a Initialized event raised by the UUPSUpgradeable contract.
type UUPSUpgradeableInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_UUPSUpgradeable *UUPSUpgradeableFilterer) FilterInitialized(opts *bind.FilterOpts) (*UUPSUpgradeableInitializedIterator, error) {
	logs, sub, err := _UUPSUpgradeable.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &UUPSUpgradeableInitializedIterator{contract: _UUPSUpgradeable.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_UUPSUpgradeable *UUPSUpgradeableFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *UUPSUpgradeableInitialized) (event.Subscription, error) {
	logs, sub, err := _UUPSUpgradeable.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UUPSUpgradeableInitialized)
				if err := _UUPSUpgradeable.contract.UnpackLog(event, "Initialized", log); err != nil {
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

// ParseInitialized is a log parse operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_UUPSUpgradeable *UUPSUpgradeableFilterer) ParseInitialized(log types.Log) (*UUPSUpgradeableInitialized, error) {
	event := new(UUPSUpgradeableInitialized)
	if err := _UUPSUpgradeable.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	return event, nil
}

// UUPSUpgradeableUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the UUPSUpgradeable contract.
type UUPSUpgradeableUpgradedIterator struct {
	Event *UUPSUpgradeableUpgraded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log      // Log channel receiving the found contract events
	sub  klaytn.Subscription // Subscription for errors, completion and termination
	done bool                // Whether the subscription completed delivering logs
	fail error               // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *UUPSUpgradeableUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UUPSUpgradeableUpgraded)
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
		it.Event = new(UUPSUpgradeableUpgraded)
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
func (it *UUPSUpgradeableUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UUPSUpgradeableUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UUPSUpgradeableUpgraded represents a Upgraded event raised by the UUPSUpgradeable contract.
type UUPSUpgradeableUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_UUPSUpgradeable *UUPSUpgradeableFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*UUPSUpgradeableUpgradedIterator, error) {
	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _UUPSUpgradeable.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &UUPSUpgradeableUpgradedIterator{contract: _UUPSUpgradeable.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_UUPSUpgradeable *UUPSUpgradeableFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *UUPSUpgradeableUpgraded, implementation []common.Address) (event.Subscription, error) {
	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _UUPSUpgradeable.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UUPSUpgradeableUpgraded)
				if err := _UUPSUpgradeable.contract.UnpackLog(event, "Upgraded", log); err != nil {
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

// ParseUpgraded is a log parse operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_UUPSUpgradeable *UUPSUpgradeableFilterer) ParseUpgraded(log types.Log) (*UUPSUpgradeableUpgraded, error) {
	event := new(UUPSUpgradeableUpgraded)
	if err := _UUPSUpgradeable.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	return event, nil
}
