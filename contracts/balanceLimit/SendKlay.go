// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package balanceLimit

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

// SendKlayABI is the input ABI used to generate the binding from.
const SendKlayABI = "[{\"inputs\":[],\"stateMutability\":\"payable\",\"type\":\"constructor\"},{\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"contract_call\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"contract_payable\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"contract_send\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"contract_transfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]"

// SendKlayBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const SendKlayBinRuntime = `6080604052600436106100405760003560e01c80631d40a1701461004957806364c561e414610069578063bef22f0e14610089578063fdbdb1cd1461004757005b3661004757005b005b34801561005557600080fd5b506100476100643660046101f8565b6100a9565b34801561007557600080fd5b506100476100843660046101f8565b610150565b34801561009557600080fd5b506100476100a43660046101f8565b6101c2565b600080836001600160a01b03168360405160006040518083038185875af1925050503d80600081146100f7576040519150601f19603f3d011682016040523d82523d6000602084013e6100fc565b606091505b50915091508161014a5760405162461bcd60e51b81526020600482015260146024820152732330b4b632b2103a379039b2b7321022ba3432b960611b60448201526064015b60405180910390fd5b50505050565b6040516000906001600160a01b0384169083156108fc0290849084818181858888f193505050509050806101bd5760405162461bcd60e51b81526020600482015260146024820152732330b4b632b2103a379039b2b7321022ba3432b960611b6044820152606401610141565b505050565b6040516001600160a01b0383169082156108fc029083906000818181858888f193505050501580156101bd573d6000803e3d6000fd5b6000806040838503121561020b57600080fd5b82356001600160a01b038116811461022257600080fd5b94602093909301359350505056fea264697066735822122097a2402f932600e3d2ba6d07a9af02b259e448a980a85ddd1ed630b6efd6cc0164736f6c63430008070033`

// SendKlayFuncSigs maps the 4-byte function signature to its string representation.
var SendKlayFuncSigs = map[string]string{
	"1d40a170": "contract_call(address,uint256)",
	"fdbdb1cd": "contract_payable()",
	"64c561e4": "contract_send(address,uint256)",
	"bef22f0e": "contract_transfer(address,uint256)",
}

// SendKlayBin is the compiled bytecode used for deploying new contracts.
var SendKlayBin = "0x6080604052610266806100136000396000f3fe6080604052600436106100405760003560e01c80631d40a1701461004957806364c561e414610069578063bef22f0e14610089578063fdbdb1cd1461004757005b3661004757005b005b34801561005557600080fd5b506100476100643660046101f8565b6100a9565b34801561007557600080fd5b506100476100843660046101f8565b610150565b34801561009557600080fd5b506100476100a43660046101f8565b6101c2565b600080836001600160a01b03168360405160006040518083038185875af1925050503d80600081146100f7576040519150601f19603f3d011682016040523d82523d6000602084013e6100fc565b606091505b50915091508161014a5760405162461bcd60e51b81526020600482015260146024820152732330b4b632b2103a379039b2b7321022ba3432b960611b60448201526064015b60405180910390fd5b50505050565b6040516000906001600160a01b0384169083156108fc0290849084818181858888f193505050509050806101bd5760405162461bcd60e51b81526020600482015260146024820152732330b4b632b2103a379039b2b7321022ba3432b960611b6044820152606401610141565b505050565b6040516001600160a01b0383169082156108fc029083906000818181858888f193505050501580156101bd573d6000803e3d6000fd5b6000806040838503121561020b57600080fd5b82356001600160a01b038116811461022257600080fd5b94602093909301359350505056fea264697066735822122097a2402f932600e3d2ba6d07a9af02b259e448a980a85ddd1ed630b6efd6cc0164736f6c63430008070033"

// DeploySendKlay deploys a new Klaytn contract, binding an instance of SendKlay to it.
func DeploySendKlay(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *SendKlay, error) {
	parsed, err := abi.JSON(strings.NewReader(SendKlayABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(SendKlayBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &SendKlay{SendKlayCaller: SendKlayCaller{contract: contract}, SendKlayTransactor: SendKlayTransactor{contract: contract}, SendKlayFilterer: SendKlayFilterer{contract: contract}}, nil
}

// SendKlay is an auto generated Go binding around a Klaytn contract.
type SendKlay struct {
	SendKlayCaller     // Read-only binding to the contract
	SendKlayTransactor // Write-only binding to the contract
	SendKlayFilterer   // Log filterer for contract events
}

// SendKlayCaller is an auto generated read-only Go binding around a Klaytn contract.
type SendKlayCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SendKlayTransactor is an auto generated write-only Go binding around a Klaytn contract.
type SendKlayTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SendKlayFilterer is an auto generated log filtering Go binding around a Klaytn contract events.
type SendKlayFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SendKlaySession is an auto generated Go binding around a Klaytn contract,
// with pre-set call and transact options.
type SendKlaySession struct {
	Contract     *SendKlay         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SendKlayCallerSession is an auto generated read-only Go binding around a Klaytn contract,
// with pre-set call options.
type SendKlayCallerSession struct {
	Contract *SendKlayCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// SendKlayTransactorSession is an auto generated write-only Go binding around a Klaytn contract,
// with pre-set transact options.
type SendKlayTransactorSession struct {
	Contract     *SendKlayTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// SendKlayRaw is an auto generated low-level Go binding around a Klaytn contract.
type SendKlayRaw struct {
	Contract *SendKlay // Generic contract binding to access the raw methods on
}

// SendKlayCallerRaw is an auto generated low-level read-only Go binding around a Klaytn contract.
type SendKlayCallerRaw struct {
	Contract *SendKlayCaller // Generic read-only contract binding to access the raw methods on
}

// SendKlayTransactorRaw is an auto generated low-level write-only Go binding around a Klaytn contract.
type SendKlayTransactorRaw struct {
	Contract *SendKlayTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSendKlay creates a new instance of SendKlay, bound to a specific deployed contract.
func NewSendKlay(address common.Address, backend bind.ContractBackend) (*SendKlay, error) {
	contract, err := bindSendKlay(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SendKlay{SendKlayCaller: SendKlayCaller{contract: contract}, SendKlayTransactor: SendKlayTransactor{contract: contract}, SendKlayFilterer: SendKlayFilterer{contract: contract}}, nil
}

// NewSendKlayCaller creates a new read-only instance of SendKlay, bound to a specific deployed contract.
func NewSendKlayCaller(address common.Address, caller bind.ContractCaller) (*SendKlayCaller, error) {
	contract, err := bindSendKlay(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SendKlayCaller{contract: contract}, nil
}

// NewSendKlayTransactor creates a new write-only instance of SendKlay, bound to a specific deployed contract.
func NewSendKlayTransactor(address common.Address, transactor bind.ContractTransactor) (*SendKlayTransactor, error) {
	contract, err := bindSendKlay(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SendKlayTransactor{contract: contract}, nil
}

// NewSendKlayFilterer creates a new log filterer instance of SendKlay, bound to a specific deployed contract.
func NewSendKlayFilterer(address common.Address, filterer bind.ContractFilterer) (*SendKlayFilterer, error) {
	contract, err := bindSendKlay(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SendKlayFilterer{contract: contract}, nil
}

// bindSendKlay binds a generic wrapper to an already deployed contract.
func bindSendKlay(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(SendKlayABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SendKlay *SendKlayRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _SendKlay.Contract.SendKlayCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SendKlay *SendKlayRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SendKlay.Contract.SendKlayTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SendKlay *SendKlayRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SendKlay.Contract.SendKlayTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SendKlay *SendKlayCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _SendKlay.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SendKlay *SendKlayTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SendKlay.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SendKlay *SendKlayTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SendKlay.Contract.contract.Transact(opts, method, params...)
}

// ContractCall is a paid mutator transaction binding the contract method 0x1d40a170.
//
// Solidity: function contract_call(address _to, uint256 amount) returns()
func (_SendKlay *SendKlayTransactor) ContractCall(opts *bind.TransactOpts, _to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _SendKlay.contract.Transact(opts, "contract_call", _to, amount)
}

// ContractCall is a paid mutator transaction binding the contract method 0x1d40a170.
//
// Solidity: function contract_call(address _to, uint256 amount) returns()
func (_SendKlay *SendKlaySession) ContractCall(_to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _SendKlay.Contract.ContractCall(&_SendKlay.TransactOpts, _to, amount)
}

// ContractCall is a paid mutator transaction binding the contract method 0x1d40a170.
//
// Solidity: function contract_call(address _to, uint256 amount) returns()
func (_SendKlay *SendKlayTransactorSession) ContractCall(_to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _SendKlay.Contract.ContractCall(&_SendKlay.TransactOpts, _to, amount)
}

// ContractPayable is a paid mutator transaction binding the contract method 0xfdbdb1cd.
//
// Solidity: function contract_payable() payable returns()
func (_SendKlay *SendKlayTransactor) ContractPayable(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SendKlay.contract.Transact(opts, "contract_payable")
}

// ContractPayable is a paid mutator transaction binding the contract method 0xfdbdb1cd.
//
// Solidity: function contract_payable() payable returns()
func (_SendKlay *SendKlaySession) ContractPayable() (*types.Transaction, error) {
	return _SendKlay.Contract.ContractPayable(&_SendKlay.TransactOpts)
}

// ContractPayable is a paid mutator transaction binding the contract method 0xfdbdb1cd.
//
// Solidity: function contract_payable() payable returns()
func (_SendKlay *SendKlayTransactorSession) ContractPayable() (*types.Transaction, error) {
	return _SendKlay.Contract.ContractPayable(&_SendKlay.TransactOpts)
}

// ContractSend is a paid mutator transaction binding the contract method 0x64c561e4.
//
// Solidity: function contract_send(address _to, uint256 amount) returns()
func (_SendKlay *SendKlayTransactor) ContractSend(opts *bind.TransactOpts, _to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _SendKlay.contract.Transact(opts, "contract_send", _to, amount)
}

// ContractSend is a paid mutator transaction binding the contract method 0x64c561e4.
//
// Solidity: function contract_send(address _to, uint256 amount) returns()
func (_SendKlay *SendKlaySession) ContractSend(_to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _SendKlay.Contract.ContractSend(&_SendKlay.TransactOpts, _to, amount)
}

// ContractSend is a paid mutator transaction binding the contract method 0x64c561e4.
//
// Solidity: function contract_send(address _to, uint256 amount) returns()
func (_SendKlay *SendKlayTransactorSession) ContractSend(_to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _SendKlay.Contract.ContractSend(&_SendKlay.TransactOpts, _to, amount)
}

// ContractTransfer is a paid mutator transaction binding the contract method 0xbef22f0e.
//
// Solidity: function contract_transfer(address _to, uint256 amount) returns()
func (_SendKlay *SendKlayTransactor) ContractTransfer(opts *bind.TransactOpts, _to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _SendKlay.contract.Transact(opts, "contract_transfer", _to, amount)
}

// ContractTransfer is a paid mutator transaction binding the contract method 0xbef22f0e.
//
// Solidity: function contract_transfer(address _to, uint256 amount) returns()
func (_SendKlay *SendKlaySession) ContractTransfer(_to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _SendKlay.Contract.ContractTransfer(&_SendKlay.TransactOpts, _to, amount)
}

// ContractTransfer is a paid mutator transaction binding the contract method 0xbef22f0e.
//
// Solidity: function contract_transfer(address _to, uint256 amount) returns()
func (_SendKlay *SendKlayTransactorSession) ContractTransfer(_to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _SendKlay.Contract.ContractTransfer(&_SendKlay.TransactOpts, _to, amount)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_SendKlay *SendKlayTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _SendKlay.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_SendKlay *SendKlaySession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _SendKlay.Contract.Fallback(&_SendKlay.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_SendKlay *SendKlayTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _SendKlay.Contract.Fallback(&_SendKlay.TransactOpts, calldata)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_SendKlay *SendKlayTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SendKlay.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_SendKlay *SendKlaySession) Receive() (*types.Transaction, error) {
	return _SendKlay.Contract.Receive(&_SendKlay.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_SendKlay *SendKlayTransactorSession) Receive() (*types.Transaction, error) {
	return _SendKlay.Contract.Receive(&_SendKlay.TransactOpts)
}
