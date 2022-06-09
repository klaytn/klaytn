// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package gov

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

// GovParamParam is an auto generated low-level Go binding around an user-defined struct.
type GovParamParam struct {
	Exists  bool
	Votable bool
	Value   []byte
}

// AddressUpgradeableABI is the input ABI used to generate the binding from.
const AddressUpgradeableABI = "[]"

// AddressUpgradeableBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const AddressUpgradeableBinRuntime = `73000000000000000000000000000000000000000030146080604052600080fdfea26469706673582212208a4aba28524878c6400bb81f8dc1e1a301c0cbbe6a6d969f445f02128571c88964736f6c634300080e0033`

// AddressUpgradeableBin is the compiled bytecode used for deploying new contracts.
var AddressUpgradeableBin = "0x60566037600b82828239805160001a607314602a57634e487b7160e01b600052600060045260246000fd5b30600052607381538281f3fe73000000000000000000000000000000000000000030146080604052600080fdfea26469706673582212208a4aba28524878c6400bb81f8dc1e1a301c0cbbe6a6d969f445f02128571c88964736f6c634300080e0033"

// DeployAddressUpgradeable deploys a new Klaytn contract, binding an instance of AddressUpgradeable to it.
func DeployAddressUpgradeable(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *AddressUpgradeable, error) {
	parsed, err := abi.JSON(strings.NewReader(AddressUpgradeableABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(AddressUpgradeableBin), backend)
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
	parsed, err := abi.JSON(strings.NewReader(AddressUpgradeableABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AddressUpgradeable *AddressUpgradeableRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
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
func (_AddressUpgradeable *AddressUpgradeableCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
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

// ContextUpgradeableABI is the input ABI used to generate the binding from.
const ContextUpgradeableABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"}]"

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
	parsed, err := abi.JSON(strings.NewReader(ContextUpgradeableABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContextUpgradeable *ContextUpgradeableRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
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
func (_ContextUpgradeable *ContextUpgradeableCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
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

// GovParamABI is the input ABI used to generate the binding from.
const GovParamABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"SetParam\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"name\":\"SetParamVotable\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"getAllParams\",\"outputs\":[{\"internalType\":\"string[]\",\"name\":\"\",\"type\":\"string[]\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"exists\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"votable\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"value\",\"type\":\"bytes\"}],\"internalType\":\"structGovParam.Param[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getGovernaceAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"name\":\"getParam\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"governanceAddr\",\"type\":\"address\"}],\"name\":\"setGovernanceAddr\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"value\",\"type\":\"bytes\"}],\"name\":\"setParam\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"value\",\"type\":\"bytes\"}],\"name\":\"setParamByOwner\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"name\":\"setParamVotable\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// GovParamBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const GovParamBinRuntime = `608060405234801561001057600080fd5b50600436106100a95760003560e01c80638da5cb5b116100715780638da5cb5b1461011a5780638e87fc1a1461013f578063a170052e14610152578063c4d66de814610168578063ee9c780a1461017b578063f2fde38b1461018c57600080fd5b806305d2b1bb146100ae5780632b813437146100c3578063400e3f6f146100d65780635d4f71d4146100e9578063715018a614610112575b600080fd5b6100c16100bc366004610ea0565b61019f565b005b6100c16100d1366004610ea0565b610270565b6100c16100e4366004610f22565b610379565b6100fc6100f7366004610f22565b61050b565b604051610109919061102b565b60405180910390f35b6100c16105f9565b6033546001600160a01b03165b6040516001600160a01b039091168152602001610109565b6100c161014d366004611045565b61062f565b61015a6106ee565b60405161010992919061106e565b6100c1610176366004611045565b610a02565b6067546001600160a01b0316610127565b6100c161019a366004611045565b610a90565b6033546001600160a01b031633146101d25760405162461bcd60e51b81526004016101c990611145565b60405180910390fd5b606684846040516101e492919061117a565b9081526040519081900360200190205460ff610100909104161561025e5760405162461bcd60e51b815260206004820152602b60248201527f476f76506172616d3a20706172616d566f7461626c65206d757374206265207360448201526a657420746f2066616c736560a81b60648201526084016101c9565b61026a84848484610b2b565b50505050565b6067546001600160a01b0316336001600160a01b0316146102e95760405162461bcd60e51b815260206004820152602d60248201527f476f76506172616d3a2063616c6c6572206973206e6f74206120476f7665726e60448201526c185b98d94810dbdb9d1c9858dd609a1b60648201526084016101c9565b606684846040516102fb92919061117a565b9081526040519081900360200190205460ff61010090910416151560011461025e5760405162461bcd60e51b815260206004820152602b60248201527f476f76506172616d3a20706172616d566f7461626c65206d757374206265202060448201526a73657420746f207472756560a81b60648201526084016101c9565b6033546001600160a01b031633146103a35760405162461bcd60e51b81526004016101c990611145565b60008151116103f45760405162461bcd60e51b815260206004820152601e60248201527f476f76506172616d3a206e616d652063616e6e6f7420626520656d707479000060448201526064016101c9565b606681604051610404919061118a565b9081526040519081900360200190205460ff1661046f5760405162461bcd60e51b815260206004820152602360248201527f476f76506172616d3a20706172616d6574657220646f6573206e6f742065786960448201526273747360e81b60648201526084016101c9565b6001606682604051610481919061118a565b90815260405190819003602001812080549215156101000261ff0019909316929092179091557f064a476da825aa99a31c682cc490b34f462457006480bda5cb27151a201213f59082906066906104d990839061118a565b90815260405190819003602001812054610500929160ff61010090920491909116906111a6565b60405180910390a150565b606060668260405161051d919061118a565b9081526040519081900360200190205460ff1661054857505060408051602081019091526000815290565b606682604051610558919061118a565b90815260200160405180910390206001018054610574906111ca565b80601f01602080910402602001604051908101604052809291908181526020018280546105a0906111ca565b80156105ed5780601f106105c2576101008083540402835291602001916105ed565b820191906000526020600020905b8154815290600101906020018083116105d057829003601f168201915b50505050509050919050565b6033546001600160a01b031633146106235760405162461bcd60e51b81526004016101c990611145565b61062d6000610c8a565b565b6033546001600160a01b031633146106595760405162461bcd60e51b81526004016101c990611145565b6001600160a01b0381166106cc5760405162461bcd60e51b815260206004820152603460248201527f476f76506172616d3a20676f7665726e616e636520636f6e74726163742063616044820152736e6e6f74206265207a65726f206164647265737360601b60648201526084016101c9565b606780546001600160a01b0319166001600160a01b0392909216919091179055565b606080600060658054905067ffffffffffffffff81111561071157610711610f0c565b60405190808252806020026020018201604052801561075e57816020015b6040805160608082018352600080835260208301529181019190915281526020019060019003908161072f5790505b50905060005b6065548110156109225760006065828154811061078357610783611204565b906000526020600020018054610798906111ca565b80601f01602080910402602001604051908101604052809291908181526020018280546107c4906111ca565b80156108115780601f106107e657610100808354040283529160200191610811565b820191906000526020600020905b8154815290600101906020018083116107f457829003601f168201915b50505050509050606681604051610828919061118a565b908152604080519182900360209081018320606084018352805460ff8082161515865261010090910416151591840191909152600181018054919284019161086f906111ca565b80601f016020809104026020016040519081016040528092919081815260200182805461089b906111ca565b80156108e85780601f106108bd576101008083540402835291602001916108e8565b820191906000526020600020905b8154815290600101906020018083116108cb57829003601f168201915b50505050508152505083838151811061090357610903611204565b602002602001018190525050808061091a9061121a565b915050610764565b5060658181805480602002602001604051908101604052809291908181526020016000905b828210156109f3578382906000526020600020018054610966906111ca565b80601f0160208091040260200160405190810160405280929190818152602001828054610992906111ca565b80156109df5780601f106109b4576101008083540402835291602001916109df565b820191906000526020600020905b8154815290600101906020018083116109c257829003601f168201915b505050505081526020019060010190610947565b50505050915092509250509091565b6000610a0e6001610cdc565b90508015610a26576000805461ff0019166101001790555b610a2e610d64565b6001600160a01b03821615610a4657610a4682610c8a565b8015610a8c576000805461ff0019169055604051600181527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15b5050565b6033546001600160a01b03163314610aba5760405162461bcd60e51b81526004016101c990611145565b6001600160a01b038116610b1f5760405162461bcd60e51b815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201526564647265737360d01b60648201526084016101c9565b610b2881610c8a565b50565b82610b785760405162461bcd60e51b815260206004820152601e60248201527f476f76506172616d3a206e616d652063616e6e6f7420626520656d707479000060448201526064016101c9565b60668484604051610b8a92919061117a565b9081526040519081900360200190205460ff16610c1457600160668585604051610bb592919061117a565b908152604051908190036020019020805491151560ff1990921691909117905560658054600181018255600091909152610c12907f8ff97419363ffd7000167f130ef7168fbea05faf9251824ca5043f113cc6a7c7018585610dbe565b505b818160668686604051610c2892919061117a565b90815260200160405180910390206001019190610c46929190610dbe565b507f9067e93c0b7173962f200e7f543c8b0496b9d00bcf35715553da6a4a9f66d0a784848484604051610c7c949392919061126a565b60405180910390a150505050565b603380546001600160a01b038381166001600160a01b0319831681179093556040519116919082907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a35050565b60008054610100900460ff1615610d23578160ff166001148015610cff5750303b155b610d1b5760405162461bcd60e51b81526004016101c99061129c565b506000919050565b60005460ff808416911610610d4a5760405162461bcd60e51b81526004016101c99061129c565b506000805460ff191660ff92909216919091179055600190565b600054610100900460ff16610d8b5760405162461bcd60e51b81526004016101c9906112ea565b61062d600054610100900460ff16610db55760405162461bcd60e51b81526004016101c9906112ea565b61062d33610c8a565b828054610dca906111ca565b90600052602060002090601f016020900481019282610dec5760008555610e32565b82601f10610e055782800160ff19823516178555610e32565b82800160010185558215610e32579182015b82811115610e32578235825591602001919060010190610e17565b50610e3e929150610e42565b5090565b5b80821115610e3e5760008155600101610e43565b60008083601f840112610e6957600080fd5b50813567ffffffffffffffff811115610e8157600080fd5b602083019150836020828501011115610e9957600080fd5b9250929050565b60008060008060408587031215610eb657600080fd5b843567ffffffffffffffff80821115610ece57600080fd5b610eda88838901610e57565b90965094506020870135915080821115610ef357600080fd5b50610f0087828801610e57565b95989497509550505050565b634e487b7160e01b600052604160045260246000fd5b600060208284031215610f3457600080fd5b813567ffffffffffffffff80821115610f4c57600080fd5b818401915084601f830112610f6057600080fd5b813581811115610f7257610f72610f0c565b604051601f8201601f19908116603f01168101908382118183101715610f9a57610f9a610f0c565b81604052828152876020848701011115610fb357600080fd5b826020860160208301376000928101602001929092525095945050505050565b60005b83811015610fee578181015183820152602001610fd6565b8381111561026a5750506000910152565b60008151808452611017816020860160208601610fd3565b601f01601f19169290920160200192915050565b60208152600061103e6020830184610fff565b9392505050565b60006020828403121561105757600080fd5b81356001600160a01b038116811461103e57600080fd5b60006040808301818452808651808352606092508286019150828160051b8701016020808a0160005b848110156110c557605f198a85030186526110b3848351610fff565b95830195935090820190600101611097565b505087820381890152885180835281830194509250600583901b8201810189820160005b8581101561113457848303601f1901875281518051151584528481015115158585015289015189840189905261112189850182610fff565b97850197935050908301906001016110e9565b50909b9a5050505050505050505050565b6020808252818101527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572604082015260600190565b8183823760009101908152919050565b6000825161119c818460208701610fd3565b9190910192915050565b6040815260006111b96040830185610fff565b905082151560208301529392505050565b600181811c908216806111de57607f821691505b6020821081036111fe57634e487b7160e01b600052602260045260246000fd5b50919050565b634e487b7160e01b600052603260045260246000fd5b60006001820161123a57634e487b7160e01b600052601160045260246000fd5b5060010190565b81835281816020850137506000828201602090810191909152601f909101601f19169091010190565b60408152600061127e604083018688611241565b8281036020840152611291818587611241565b979650505050505050565b6020808252602e908201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160408201526d191e481a5b9a5d1a585b1a5e995960921b606082015260800190565b6020808252602b908201527f496e697469616c697a61626c653a20636f6e7472616374206973206e6f74206960408201526a6e697469616c697a696e6760a81b60608201526080019056fea2646970667358221220c462d6dbdae6628f5995bb10b97bc2f5d900ffe41f89c1a37108105f6d0091a464736f6c634300080e0033`

// GovParamFuncSigs maps the 4-byte function signature to its string representation.
var GovParamFuncSigs = map[string]string{
	"a170052e": "getAllParams()",
	"ee9c780a": "getGovernaceAddress()",
	"5d4f71d4": "getParam(string)",
	"c4d66de8": "initialize(address)",
	"8da5cb5b": "owner()",
	"715018a6": "renounceOwnership()",
	"8e87fc1a": "setGovernanceAddr(address)",
	"2b813437": "setParam(string,bytes)",
	"05d2b1bb": "setParamByOwner(string,bytes)",
	"400e3f6f": "setParamVotable(string)",
	"f2fde38b": "transferOwnership(address)",
}

// GovParamBin is the compiled bytecode used for deploying new contracts.
var GovParamBin = "0x608060405234801561001057600080fd5b5061136b806100206000396000f3fe608060405234801561001057600080fd5b50600436106100a95760003560e01c80638da5cb5b116100715780638da5cb5b1461011a5780638e87fc1a1461013f578063a170052e14610152578063c4d66de814610168578063ee9c780a1461017b578063f2fde38b1461018c57600080fd5b806305d2b1bb146100ae5780632b813437146100c3578063400e3f6f146100d65780635d4f71d4146100e9578063715018a614610112575b600080fd5b6100c16100bc366004610ea0565b61019f565b005b6100c16100d1366004610ea0565b610270565b6100c16100e4366004610f22565b610379565b6100fc6100f7366004610f22565b61050b565b604051610109919061102b565b60405180910390f35b6100c16105f9565b6033546001600160a01b03165b6040516001600160a01b039091168152602001610109565b6100c161014d366004611045565b61062f565b61015a6106ee565b60405161010992919061106e565b6100c1610176366004611045565b610a02565b6067546001600160a01b0316610127565b6100c161019a366004611045565b610a90565b6033546001600160a01b031633146101d25760405162461bcd60e51b81526004016101c990611145565b60405180910390fd5b606684846040516101e492919061117a565b9081526040519081900360200190205460ff610100909104161561025e5760405162461bcd60e51b815260206004820152602b60248201527f476f76506172616d3a20706172616d566f7461626c65206d757374206265207360448201526a657420746f2066616c736560a81b60648201526084016101c9565b61026a84848484610b2b565b50505050565b6067546001600160a01b0316336001600160a01b0316146102e95760405162461bcd60e51b815260206004820152602d60248201527f476f76506172616d3a2063616c6c6572206973206e6f74206120476f7665726e60448201526c185b98d94810dbdb9d1c9858dd609a1b60648201526084016101c9565b606684846040516102fb92919061117a565b9081526040519081900360200190205460ff61010090910416151560011461025e5760405162461bcd60e51b815260206004820152602b60248201527f476f76506172616d3a20706172616d566f7461626c65206d757374206265202060448201526a73657420746f207472756560a81b60648201526084016101c9565b6033546001600160a01b031633146103a35760405162461bcd60e51b81526004016101c990611145565b60008151116103f45760405162461bcd60e51b815260206004820152601e60248201527f476f76506172616d3a206e616d652063616e6e6f7420626520656d707479000060448201526064016101c9565b606681604051610404919061118a565b9081526040519081900360200190205460ff1661046f5760405162461bcd60e51b815260206004820152602360248201527f476f76506172616d3a20706172616d6574657220646f6573206e6f742065786960448201526273747360e81b60648201526084016101c9565b6001606682604051610481919061118a565b90815260405190819003602001812080549215156101000261ff0019909316929092179091557f064a476da825aa99a31c682cc490b34f462457006480bda5cb27151a201213f59082906066906104d990839061118a565b90815260405190819003602001812054610500929160ff61010090920491909116906111a6565b60405180910390a150565b606060668260405161051d919061118a565b9081526040519081900360200190205460ff1661054857505060408051602081019091526000815290565b606682604051610558919061118a565b90815260200160405180910390206001018054610574906111ca565b80601f01602080910402602001604051908101604052809291908181526020018280546105a0906111ca565b80156105ed5780601f106105c2576101008083540402835291602001916105ed565b820191906000526020600020905b8154815290600101906020018083116105d057829003601f168201915b50505050509050919050565b6033546001600160a01b031633146106235760405162461bcd60e51b81526004016101c990611145565b61062d6000610c8a565b565b6033546001600160a01b031633146106595760405162461bcd60e51b81526004016101c990611145565b6001600160a01b0381166106cc5760405162461bcd60e51b815260206004820152603460248201527f476f76506172616d3a20676f7665726e616e636520636f6e74726163742063616044820152736e6e6f74206265207a65726f206164647265737360601b60648201526084016101c9565b606780546001600160a01b0319166001600160a01b0392909216919091179055565b606080600060658054905067ffffffffffffffff81111561071157610711610f0c565b60405190808252806020026020018201604052801561075e57816020015b6040805160608082018352600080835260208301529181019190915281526020019060019003908161072f5790505b50905060005b6065548110156109225760006065828154811061078357610783611204565b906000526020600020018054610798906111ca565b80601f01602080910402602001604051908101604052809291908181526020018280546107c4906111ca565b80156108115780601f106107e657610100808354040283529160200191610811565b820191906000526020600020905b8154815290600101906020018083116107f457829003601f168201915b50505050509050606681604051610828919061118a565b908152604080519182900360209081018320606084018352805460ff8082161515865261010090910416151591840191909152600181018054919284019161086f906111ca565b80601f016020809104026020016040519081016040528092919081815260200182805461089b906111ca565b80156108e85780601f106108bd576101008083540402835291602001916108e8565b820191906000526020600020905b8154815290600101906020018083116108cb57829003601f168201915b50505050508152505083838151811061090357610903611204565b602002602001018190525050808061091a9061121a565b915050610764565b5060658181805480602002602001604051908101604052809291908181526020016000905b828210156109f3578382906000526020600020018054610966906111ca565b80601f0160208091040260200160405190810160405280929190818152602001828054610992906111ca565b80156109df5780601f106109b4576101008083540402835291602001916109df565b820191906000526020600020905b8154815290600101906020018083116109c257829003601f168201915b505050505081526020019060010190610947565b50505050915092509250509091565b6000610a0e6001610cdc565b90508015610a26576000805461ff0019166101001790555b610a2e610d64565b6001600160a01b03821615610a4657610a4682610c8a565b8015610a8c576000805461ff0019169055604051600181527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15b5050565b6033546001600160a01b03163314610aba5760405162461bcd60e51b81526004016101c990611145565b6001600160a01b038116610b1f5760405162461bcd60e51b815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201526564647265737360d01b60648201526084016101c9565b610b2881610c8a565b50565b82610b785760405162461bcd60e51b815260206004820152601e60248201527f476f76506172616d3a206e616d652063616e6e6f7420626520656d707479000060448201526064016101c9565b60668484604051610b8a92919061117a565b9081526040519081900360200190205460ff16610c1457600160668585604051610bb592919061117a565b908152604051908190036020019020805491151560ff1990921691909117905560658054600181018255600091909152610c12907f8ff97419363ffd7000167f130ef7168fbea05faf9251824ca5043f113cc6a7c7018585610dbe565b505b818160668686604051610c2892919061117a565b90815260200160405180910390206001019190610c46929190610dbe565b507f9067e93c0b7173962f200e7f543c8b0496b9d00bcf35715553da6a4a9f66d0a784848484604051610c7c949392919061126a565b60405180910390a150505050565b603380546001600160a01b038381166001600160a01b0319831681179093556040519116919082907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a35050565b60008054610100900460ff1615610d23578160ff166001148015610cff5750303b155b610d1b5760405162461bcd60e51b81526004016101c99061129c565b506000919050565b60005460ff808416911610610d4a5760405162461bcd60e51b81526004016101c99061129c565b506000805460ff191660ff92909216919091179055600190565b600054610100900460ff16610d8b5760405162461bcd60e51b81526004016101c9906112ea565b61062d600054610100900460ff16610db55760405162461bcd60e51b81526004016101c9906112ea565b61062d33610c8a565b828054610dca906111ca565b90600052602060002090601f016020900481019282610dec5760008555610e32565b82601f10610e055782800160ff19823516178555610e32565b82800160010185558215610e32579182015b82811115610e32578235825591602001919060010190610e17565b50610e3e929150610e42565b5090565b5b80821115610e3e5760008155600101610e43565b60008083601f840112610e6957600080fd5b50813567ffffffffffffffff811115610e8157600080fd5b602083019150836020828501011115610e9957600080fd5b9250929050565b60008060008060408587031215610eb657600080fd5b843567ffffffffffffffff80821115610ece57600080fd5b610eda88838901610e57565b90965094506020870135915080821115610ef357600080fd5b50610f0087828801610e57565b95989497509550505050565b634e487b7160e01b600052604160045260246000fd5b600060208284031215610f3457600080fd5b813567ffffffffffffffff80821115610f4c57600080fd5b818401915084601f830112610f6057600080fd5b813581811115610f7257610f72610f0c565b604051601f8201601f19908116603f01168101908382118183101715610f9a57610f9a610f0c565b81604052828152876020848701011115610fb357600080fd5b826020860160208301376000928101602001929092525095945050505050565b60005b83811015610fee578181015183820152602001610fd6565b8381111561026a5750506000910152565b60008151808452611017816020860160208601610fd3565b601f01601f19169290920160200192915050565b60208152600061103e6020830184610fff565b9392505050565b60006020828403121561105757600080fd5b81356001600160a01b038116811461103e57600080fd5b60006040808301818452808651808352606092508286019150828160051b8701016020808a0160005b848110156110c557605f198a85030186526110b3848351610fff565b95830195935090820190600101611097565b505087820381890152885180835281830194509250600583901b8201810189820160005b8581101561113457848303601f1901875281518051151584528481015115158585015289015189840189905261112189850182610fff565b97850197935050908301906001016110e9565b50909b9a5050505050505050505050565b6020808252818101527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572604082015260600190565b8183823760009101908152919050565b6000825161119c818460208701610fd3565b9190910192915050565b6040815260006111b96040830185610fff565b905082151560208301529392505050565b600181811c908216806111de57607f821691505b6020821081036111fe57634e487b7160e01b600052602260045260246000fd5b50919050565b634e487b7160e01b600052603260045260246000fd5b60006001820161123a57634e487b7160e01b600052601160045260246000fd5b5060010190565b81835281816020850137506000828201602090810191909152601f909101601f19169091010190565b60408152600061127e604083018688611241565b8281036020840152611291818587611241565b979650505050505050565b6020808252602e908201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160408201526d191e481a5b9a5d1a585b1a5e995960921b606082015260800190565b6020808252602b908201527f496e697469616c697a61626c653a20636f6e7472616374206973206e6f74206960408201526a6e697469616c697a696e6760a81b60608201526080019056fea2646970667358221220c462d6dbdae6628f5995bb10b97bc2f5d900ffe41f89c1a37108105f6d0091a464736f6c634300080e0033"

// DeployGovParam deploys a new Klaytn contract, binding an instance of GovParam to it.
func DeployGovParam(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *GovParam, error) {
	parsed, err := abi.JSON(strings.NewReader(GovParamABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(GovParamBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &GovParam{GovParamCaller: GovParamCaller{contract: contract}, GovParamTransactor: GovParamTransactor{contract: contract}, GovParamFilterer: GovParamFilterer{contract: contract}}, nil
}

// GovParam is an auto generated Go binding around a Klaytn contract.
type GovParam struct {
	GovParamCaller     // Read-only binding to the contract
	GovParamTransactor // Write-only binding to the contract
	GovParamFilterer   // Log filterer for contract events
}

// GovParamCaller is an auto generated read-only Go binding around a Klaytn contract.
type GovParamCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GovParamTransactor is an auto generated write-only Go binding around a Klaytn contract.
type GovParamTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GovParamFilterer is an auto generated log filtering Go binding around a Klaytn contract events.
type GovParamFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GovParamSession is an auto generated Go binding around a Klaytn contract,
// with pre-set call and transact options.
type GovParamSession struct {
	Contract     *GovParam         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// GovParamCallerSession is an auto generated read-only Go binding around a Klaytn contract,
// with pre-set call options.
type GovParamCallerSession struct {
	Contract *GovParamCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// GovParamTransactorSession is an auto generated write-only Go binding around a Klaytn contract,
// with pre-set transact options.
type GovParamTransactorSession struct {
	Contract     *GovParamTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// GovParamRaw is an auto generated low-level Go binding around a Klaytn contract.
type GovParamRaw struct {
	Contract *GovParam // Generic contract binding to access the raw methods on
}

// GovParamCallerRaw is an auto generated low-level read-only Go binding around a Klaytn contract.
type GovParamCallerRaw struct {
	Contract *GovParamCaller // Generic read-only contract binding to access the raw methods on
}

// GovParamTransactorRaw is an auto generated low-level write-only Go binding around a Klaytn contract.
type GovParamTransactorRaw struct {
	Contract *GovParamTransactor // Generic write-only contract binding to access the raw methods on
}

// NewGovParam creates a new instance of GovParam, bound to a specific deployed contract.
func NewGovParam(address common.Address, backend bind.ContractBackend) (*GovParam, error) {
	contract, err := bindGovParam(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &GovParam{GovParamCaller: GovParamCaller{contract: contract}, GovParamTransactor: GovParamTransactor{contract: contract}, GovParamFilterer: GovParamFilterer{contract: contract}}, nil
}

// NewGovParamCaller creates a new read-only instance of GovParam, bound to a specific deployed contract.
func NewGovParamCaller(address common.Address, caller bind.ContractCaller) (*GovParamCaller, error) {
	contract, err := bindGovParam(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &GovParamCaller{contract: contract}, nil
}

// NewGovParamTransactor creates a new write-only instance of GovParam, bound to a specific deployed contract.
func NewGovParamTransactor(address common.Address, transactor bind.ContractTransactor) (*GovParamTransactor, error) {
	contract, err := bindGovParam(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &GovParamTransactor{contract: contract}, nil
}

// NewGovParamFilterer creates a new log filterer instance of GovParam, bound to a specific deployed contract.
func NewGovParamFilterer(address common.Address, filterer bind.ContractFilterer) (*GovParamFilterer, error) {
	contract, err := bindGovParam(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &GovParamFilterer{contract: contract}, nil
}

// bindGovParam binds a generic wrapper to an already deployed contract.
func bindGovParam(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(GovParamABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_GovParam *GovParamRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _GovParam.Contract.GovParamCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_GovParam *GovParamRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GovParam.Contract.GovParamTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_GovParam *GovParamRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _GovParam.Contract.GovParamTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_GovParam *GovParamCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _GovParam.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_GovParam *GovParamTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GovParam.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_GovParam *GovParamTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _GovParam.Contract.contract.Transact(opts, method, params...)
}

// GetAllParams is a free data retrieval call binding the contract method 0xa170052e.
//
// Solidity: function getAllParams() view returns(string[], (bool,bool,bytes)[])
func (_GovParam *GovParamCaller) GetAllParams(opts *bind.CallOpts) ([]string, []GovParamParam, error) {
	var (
		ret0 = new([]string)
		ret1 = new([]GovParamParam)
	)
	out := &[]interface{}{
		ret0,
		ret1,
	}
	err := _GovParam.contract.Call(opts, out, "getAllParams")
	return *ret0, *ret1, err
}

// GetAllParams is a free data retrieval call binding the contract method 0xa170052e.
//
// Solidity: function getAllParams() view returns(string[], (bool,bool,bytes)[])
func (_GovParam *GovParamSession) GetAllParams() ([]string, []GovParamParam, error) {
	return _GovParam.Contract.GetAllParams(&_GovParam.CallOpts)
}

// GetAllParams is a free data retrieval call binding the contract method 0xa170052e.
//
// Solidity: function getAllParams() view returns(string[], (bool,bool,bytes)[])
func (_GovParam *GovParamCallerSession) GetAllParams() ([]string, []GovParamParam, error) {
	return _GovParam.Contract.GetAllParams(&_GovParam.CallOpts)
}

// GetGovernaceAddress is a free data retrieval call binding the contract method 0xee9c780a.
//
// Solidity: function getGovernaceAddress() view returns(address)
func (_GovParam *GovParamCaller) GetGovernaceAddress(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _GovParam.contract.Call(opts, out, "getGovernaceAddress")
	return *ret0, err
}

// GetGovernaceAddress is a free data retrieval call binding the contract method 0xee9c780a.
//
// Solidity: function getGovernaceAddress() view returns(address)
func (_GovParam *GovParamSession) GetGovernaceAddress() (common.Address, error) {
	return _GovParam.Contract.GetGovernaceAddress(&_GovParam.CallOpts)
}

// GetGovernaceAddress is a free data retrieval call binding the contract method 0xee9c780a.
//
// Solidity: function getGovernaceAddress() view returns(address)
func (_GovParam *GovParamCallerSession) GetGovernaceAddress() (common.Address, error) {
	return _GovParam.Contract.GetGovernaceAddress(&_GovParam.CallOpts)
}

// GetParam is a free data retrieval call binding the contract method 0x5d4f71d4.
//
// Solidity: function getParam(string name) view returns(bytes)
func (_GovParam *GovParamCaller) GetParam(opts *bind.CallOpts, name string) ([]byte, error) {
	var (
		ret0 = new([]byte)
	)
	out := ret0
	err := _GovParam.contract.Call(opts, out, "getParam", name)
	return *ret0, err
}

// GetParam is a free data retrieval call binding the contract method 0x5d4f71d4.
//
// Solidity: function getParam(string name) view returns(bytes)
func (_GovParam *GovParamSession) GetParam(name string) ([]byte, error) {
	return _GovParam.Contract.GetParam(&_GovParam.CallOpts, name)
}

// GetParam is a free data retrieval call binding the contract method 0x5d4f71d4.
//
// Solidity: function getParam(string name) view returns(bytes)
func (_GovParam *GovParamCallerSession) GetParam(name string) ([]byte, error) {
	return _GovParam.Contract.GetParam(&_GovParam.CallOpts, name)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_GovParam *GovParamCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _GovParam.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_GovParam *GovParamSession) Owner() (common.Address, error) {
	return _GovParam.Contract.Owner(&_GovParam.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_GovParam *GovParamCallerSession) Owner() (common.Address, error) {
	return _GovParam.Contract.Owner(&_GovParam.CallOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address owner) returns()
func (_GovParam *GovParamTransactor) Initialize(opts *bind.TransactOpts, owner common.Address) (*types.Transaction, error) {
	return _GovParam.contract.Transact(opts, "initialize", owner)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address owner) returns()
func (_GovParam *GovParamSession) Initialize(owner common.Address) (*types.Transaction, error) {
	return _GovParam.Contract.Initialize(&_GovParam.TransactOpts, owner)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address owner) returns()
func (_GovParam *GovParamTransactorSession) Initialize(owner common.Address) (*types.Transaction, error) {
	return _GovParam.Contract.Initialize(&_GovParam.TransactOpts, owner)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_GovParam *GovParamTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GovParam.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_GovParam *GovParamSession) RenounceOwnership() (*types.Transaction, error) {
	return _GovParam.Contract.RenounceOwnership(&_GovParam.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_GovParam *GovParamTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _GovParam.Contract.RenounceOwnership(&_GovParam.TransactOpts)
}

// SetGovernanceAddr is a paid mutator transaction binding the contract method 0x8e87fc1a.
//
// Solidity: function setGovernanceAddr(address governanceAddr) returns()
func (_GovParam *GovParamTransactor) SetGovernanceAddr(opts *bind.TransactOpts, governanceAddr common.Address) (*types.Transaction, error) {
	return _GovParam.contract.Transact(opts, "setGovernanceAddr", governanceAddr)
}

// SetGovernanceAddr is a paid mutator transaction binding the contract method 0x8e87fc1a.
//
// Solidity: function setGovernanceAddr(address governanceAddr) returns()
func (_GovParam *GovParamSession) SetGovernanceAddr(governanceAddr common.Address) (*types.Transaction, error) {
	return _GovParam.Contract.SetGovernanceAddr(&_GovParam.TransactOpts, governanceAddr)
}

// SetGovernanceAddr is a paid mutator transaction binding the contract method 0x8e87fc1a.
//
// Solidity: function setGovernanceAddr(address governanceAddr) returns()
func (_GovParam *GovParamTransactorSession) SetGovernanceAddr(governanceAddr common.Address) (*types.Transaction, error) {
	return _GovParam.Contract.SetGovernanceAddr(&_GovParam.TransactOpts, governanceAddr)
}

// SetParam is a paid mutator transaction binding the contract method 0x2b813437.
//
// Solidity: function setParam(string name, bytes value) returns()
func (_GovParam *GovParamTransactor) SetParam(opts *bind.TransactOpts, name string, value []byte) (*types.Transaction, error) {
	return _GovParam.contract.Transact(opts, "setParam", name, value)
}

// SetParam is a paid mutator transaction binding the contract method 0x2b813437.
//
// Solidity: function setParam(string name, bytes value) returns()
func (_GovParam *GovParamSession) SetParam(name string, value []byte) (*types.Transaction, error) {
	return _GovParam.Contract.SetParam(&_GovParam.TransactOpts, name, value)
}

// SetParam is a paid mutator transaction binding the contract method 0x2b813437.
//
// Solidity: function setParam(string name, bytes value) returns()
func (_GovParam *GovParamTransactorSession) SetParam(name string, value []byte) (*types.Transaction, error) {
	return _GovParam.Contract.SetParam(&_GovParam.TransactOpts, name, value)
}

// SetParamByOwner is a paid mutator transaction binding the contract method 0x05d2b1bb.
//
// Solidity: function setParamByOwner(string name, bytes value) returns()
func (_GovParam *GovParamTransactor) SetParamByOwner(opts *bind.TransactOpts, name string, value []byte) (*types.Transaction, error) {
	return _GovParam.contract.Transact(opts, "setParamByOwner", name, value)
}

// SetParamByOwner is a paid mutator transaction binding the contract method 0x05d2b1bb.
//
// Solidity: function setParamByOwner(string name, bytes value) returns()
func (_GovParam *GovParamSession) SetParamByOwner(name string, value []byte) (*types.Transaction, error) {
	return _GovParam.Contract.SetParamByOwner(&_GovParam.TransactOpts, name, value)
}

// SetParamByOwner is a paid mutator transaction binding the contract method 0x05d2b1bb.
//
// Solidity: function setParamByOwner(string name, bytes value) returns()
func (_GovParam *GovParamTransactorSession) SetParamByOwner(name string, value []byte) (*types.Transaction, error) {
	return _GovParam.Contract.SetParamByOwner(&_GovParam.TransactOpts, name, value)
}

// SetParamVotable is a paid mutator transaction binding the contract method 0x400e3f6f.
//
// Solidity: function setParamVotable(string name) returns()
func (_GovParam *GovParamTransactor) SetParamVotable(opts *bind.TransactOpts, name string) (*types.Transaction, error) {
	return _GovParam.contract.Transact(opts, "setParamVotable", name)
}

// SetParamVotable is a paid mutator transaction binding the contract method 0x400e3f6f.
//
// Solidity: function setParamVotable(string name) returns()
func (_GovParam *GovParamSession) SetParamVotable(name string) (*types.Transaction, error) {
	return _GovParam.Contract.SetParamVotable(&_GovParam.TransactOpts, name)
}

// SetParamVotable is a paid mutator transaction binding the contract method 0x400e3f6f.
//
// Solidity: function setParamVotable(string name) returns()
func (_GovParam *GovParamTransactorSession) SetParamVotable(name string) (*types.Transaction, error) {
	return _GovParam.Contract.SetParamVotable(&_GovParam.TransactOpts, name)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_GovParam *GovParamTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _GovParam.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_GovParam *GovParamSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _GovParam.Contract.TransferOwnership(&_GovParam.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_GovParam *GovParamTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _GovParam.Contract.TransferOwnership(&_GovParam.TransactOpts, newOwner)
}

// GovParamInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the GovParam contract.
type GovParamInitializedIterator struct {
	Event *GovParamInitialized // Event containing the contract specifics and raw log

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
func (it *GovParamInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GovParamInitialized)
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
		it.Event = new(GovParamInitialized)
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
func (it *GovParamInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GovParamInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GovParamInitialized represents a Initialized event raised by the GovParam contract.
type GovParamInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_GovParam *GovParamFilterer) FilterInitialized(opts *bind.FilterOpts) (*GovParamInitializedIterator, error) {

	logs, sub, err := _GovParam.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &GovParamInitializedIterator{contract: _GovParam.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_GovParam *GovParamFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *GovParamInitialized) (event.Subscription, error) {

	logs, sub, err := _GovParam.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GovParamInitialized)
				if err := _GovParam.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_GovParam *GovParamFilterer) ParseInitialized(log types.Log) (*GovParamInitialized, error) {
	event := new(GovParamInitialized)
	if err := _GovParam.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	return event, nil
}

// GovParamOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the GovParam contract.
type GovParamOwnershipTransferredIterator struct {
	Event *GovParamOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *GovParamOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GovParamOwnershipTransferred)
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
		it.Event = new(GovParamOwnershipTransferred)
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
func (it *GovParamOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GovParamOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GovParamOwnershipTransferred represents a OwnershipTransferred event raised by the GovParam contract.
type GovParamOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_GovParam *GovParamFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*GovParamOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _GovParam.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &GovParamOwnershipTransferredIterator{contract: _GovParam.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_GovParam *GovParamFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *GovParamOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _GovParam.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GovParamOwnershipTransferred)
				if err := _GovParam.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_GovParam *GovParamFilterer) ParseOwnershipTransferred(log types.Log) (*GovParamOwnershipTransferred, error) {
	event := new(GovParamOwnershipTransferred)
	if err := _GovParam.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	return event, nil
}

// GovParamSetParamIterator is returned from FilterSetParam and is used to iterate over the raw logs and unpacked data for SetParam events raised by the GovParam contract.
type GovParamSetParamIterator struct {
	Event *GovParamSetParam // Event containing the contract specifics and raw log

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
func (it *GovParamSetParamIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GovParamSetParam)
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
		it.Event = new(GovParamSetParam)
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
func (it *GovParamSetParamIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GovParamSetParamIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GovParamSetParam represents a SetParam event raised by the GovParam contract.
type GovParamSetParam struct {
	Arg0 string
	Arg1 []byte
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterSetParam is a free log retrieval operation binding the contract event 0x9067e93c0b7173962f200e7f543c8b0496b9d00bcf35715553da6a4a9f66d0a7.
//
// Solidity: event SetParam(string arg0, bytes arg1)
func (_GovParam *GovParamFilterer) FilterSetParam(opts *bind.FilterOpts) (*GovParamSetParamIterator, error) {

	logs, sub, err := _GovParam.contract.FilterLogs(opts, "SetParam")
	if err != nil {
		return nil, err
	}
	return &GovParamSetParamIterator{contract: _GovParam.contract, event: "SetParam", logs: logs, sub: sub}, nil
}

// WatchSetParam is a free log subscription operation binding the contract event 0x9067e93c0b7173962f200e7f543c8b0496b9d00bcf35715553da6a4a9f66d0a7.
//
// Solidity: event SetParam(string arg0, bytes arg1)
func (_GovParam *GovParamFilterer) WatchSetParam(opts *bind.WatchOpts, sink chan<- *GovParamSetParam) (event.Subscription, error) {

	logs, sub, err := _GovParam.contract.WatchLogs(opts, "SetParam")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GovParamSetParam)
				if err := _GovParam.contract.UnpackLog(event, "SetParam", log); err != nil {
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

// ParseSetParam is a log parse operation binding the contract event 0x9067e93c0b7173962f200e7f543c8b0496b9d00bcf35715553da6a4a9f66d0a7.
//
// Solidity: event SetParam(string arg0, bytes arg1)
func (_GovParam *GovParamFilterer) ParseSetParam(log types.Log) (*GovParamSetParam, error) {
	event := new(GovParamSetParam)
	if err := _GovParam.contract.UnpackLog(event, "SetParam", log); err != nil {
		return nil, err
	}
	return event, nil
}

// GovParamSetParamVotableIterator is returned from FilterSetParamVotable and is used to iterate over the raw logs and unpacked data for SetParamVotable events raised by the GovParam contract.
type GovParamSetParamVotableIterator struct {
	Event *GovParamSetParamVotable // Event containing the contract specifics and raw log

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
func (it *GovParamSetParamVotableIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GovParamSetParamVotable)
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
		it.Event = new(GovParamSetParamVotable)
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
func (it *GovParamSetParamVotableIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GovParamSetParamVotableIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GovParamSetParamVotable represents a SetParamVotable event raised by the GovParam contract.
type GovParamSetParamVotable struct {
	Arg0 string
	Arg1 bool
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterSetParamVotable is a free log retrieval operation binding the contract event 0x064a476da825aa99a31c682cc490b34f462457006480bda5cb27151a201213f5.
//
// Solidity: event SetParamVotable(string arg0, bool arg1)
func (_GovParam *GovParamFilterer) FilterSetParamVotable(opts *bind.FilterOpts) (*GovParamSetParamVotableIterator, error) {

	logs, sub, err := _GovParam.contract.FilterLogs(opts, "SetParamVotable")
	if err != nil {
		return nil, err
	}
	return &GovParamSetParamVotableIterator{contract: _GovParam.contract, event: "SetParamVotable", logs: logs, sub: sub}, nil
}

// WatchSetParamVotable is a free log subscription operation binding the contract event 0x064a476da825aa99a31c682cc490b34f462457006480bda5cb27151a201213f5.
//
// Solidity: event SetParamVotable(string arg0, bool arg1)
func (_GovParam *GovParamFilterer) WatchSetParamVotable(opts *bind.WatchOpts, sink chan<- *GovParamSetParamVotable) (event.Subscription, error) {

	logs, sub, err := _GovParam.contract.WatchLogs(opts, "SetParamVotable")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GovParamSetParamVotable)
				if err := _GovParam.contract.UnpackLog(event, "SetParamVotable", log); err != nil {
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

// ParseSetParamVotable is a log parse operation binding the contract event 0x064a476da825aa99a31c682cc490b34f462457006480bda5cb27151a201213f5.
//
// Solidity: event SetParamVotable(string arg0, bool arg1)
func (_GovParam *GovParamFilterer) ParseSetParamVotable(log types.Log) (*GovParamSetParamVotable, error) {
	event := new(GovParamSetParamVotable)
	if err := _GovParam.contract.UnpackLog(event, "SetParamVotable", log); err != nil {
		return nil, err
	}
	return event, nil
}

// InitializableABI is the input ABI used to generate the binding from.
const InitializableABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"}]"

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
	parsed, err := abi.JSON(strings.NewReader(InitializableABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Initializable *InitializableRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
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
func (_Initializable *InitializableCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
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

// OwnableUpgradeableABI is the input ABI used to generate the binding from.
const OwnableUpgradeableABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// OwnableUpgradeableBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const OwnableUpgradeableBinRuntime = ``

// OwnableUpgradeableFuncSigs maps the 4-byte function signature to its string representation.
var OwnableUpgradeableFuncSigs = map[string]string{
	"8da5cb5b": "owner()",
	"715018a6": "renounceOwnership()",
	"f2fde38b": "transferOwnership(address)",
}

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
	parsed, err := abi.JSON(strings.NewReader(OwnableUpgradeableABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OwnableUpgradeable *OwnableUpgradeableRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
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
func (_OwnableUpgradeable *OwnableUpgradeableCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
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
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _OwnableUpgradeable.contract.Call(opts, out, "owner")
	return *ret0, err
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
