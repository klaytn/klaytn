// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contract

import (
	"math/big"
	"strings"

	"github.com/klaytn/klaytn/accounts/abi"
	"github.com/klaytn/klaytn/accounts/abi/bind"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
)

// KlaytnRewardABI is the input ABI used to generate the binding from.
const KlaytnRewardABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"totalAmount\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"reward\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"safeWithdrawal\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"fallback\"}]"

// KlaytnRewardBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const KlaytnRewardBinRuntime = `0x6080604052600436106100615763ffffffff7c01000000000000000000000000000000000000000000000000000000006000350416631a39d8ef81146100805780636353586b146100a757806370a08231146100ca578063fd6b7ef8146100f8575b3360009081526001602052604081208054349081019091558154019055005b34801561008c57600080fd5b5061009561010d565b60408051918252519081900360200190f35b6100c873ffffffffffffffffffffffffffffffffffffffff60043516610113565b005b3480156100d657600080fd5b5061009573ffffffffffffffffffffffffffffffffffffffff60043516610147565b34801561010457600080fd5b506100c8610159565b60005481565b73ffffffffffffffffffffffffffffffffffffffff1660009081526001602052604081208054349081019091558154019055565b60016020526000908152604090205481565b336000908152600160205260408120805490829055908111156101af57604051339082156108fc029083906000818181858888f193505050501561019c576101af565b3360009081526001602052604090208190555b505600a165627a7a723058200328c97654866fd7cb18d8f3e902b9bf0331281c7d5d6c1e47a8b1c50634b1100029`

// KlaytnRewardBin is the compiled bytecode used for deploying new contracts.
const KlaytnRewardBin = `0x608060405234801561001057600080fd5b506101de806100206000396000f3006080604052600436106100615763ffffffff7c01000000000000000000000000000000000000000000000000000000006000350416631a39d8ef81146100805780636353586b146100a757806370a08231146100ca578063fd6b7ef8146100f8575b3360009081526001602052604081208054349081019091558154019055005b34801561008c57600080fd5b5061009561010d565b60408051918252519081900360200190f35b6100c873ffffffffffffffffffffffffffffffffffffffff60043516610113565b005b3480156100d657600080fd5b5061009573ffffffffffffffffffffffffffffffffffffffff60043516610147565b34801561010457600080fd5b506100c8610159565b60005481565b73ffffffffffffffffffffffffffffffffffffffff1660009081526001602052604081208054349081019091558154019055565b60016020526000908152604090205481565b336000908152600160205260408120805490829055908111156101af57604051339082156108fc029083906000818181858888f193505050501561019c576101af565b3360009081526001602052604090208190555b505600a165627a7a723058200328c97654866fd7cb18d8f3e902b9bf0331281c7d5d6c1e47a8b1c50634b1100029`

// DeployKlaytnReward deploys a new Klaytn contract, binding an instance of KlaytnReward to it.
func DeployKlaytnReward(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *KlaytnReward, error) {
	parsed, err := abi.JSON(strings.NewReader(KlaytnRewardABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(KlaytnRewardBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &KlaytnReward{KlaytnRewardCaller: KlaytnRewardCaller{contract: contract}, KlaytnRewardTransactor: KlaytnRewardTransactor{contract: contract}, KlaytnRewardFilterer: KlaytnRewardFilterer{contract: contract}}, nil
}

// KlaytnReward is an auto generated Go binding around a Klaytn contract.
type KlaytnReward struct {
	KlaytnRewardCaller     // Read-only binding to the contract
	KlaytnRewardTransactor // Write-only binding to the contract
	KlaytnRewardFilterer   // Log filterer for contract events
}

// KlaytnRewardCaller is an auto generated read-only Go binding around a Klaytn contract.
type KlaytnRewardCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// KlaytnRewardTransactor is an auto generated write-only Go binding around a Klaytn contract.
type KlaytnRewardTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// KlaytnRewardFilterer is an auto generated log filtering Go binding around a Klaytn contract events.
type KlaytnRewardFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// KlaytnRewardSession is an auto generated Go binding around a Klaytn contract,
// with pre-set call and transact options.
type KlaytnRewardSession struct {
	Contract     *KlaytnReward     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// KlaytnRewardCallerSession is an auto generated read-only Go binding around a Klaytn contract,
// with pre-set call options.
type KlaytnRewardCallerSession struct {
	Contract *KlaytnRewardCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// KlaytnRewardTransactorSession is an auto generated write-only Go binding around a Klaytn contract,
// with pre-set transact options.
type KlaytnRewardTransactorSession struct {
	Contract     *KlaytnRewardTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// KlaytnRewardRaw is an auto generated low-level Go binding around a Klaytn contract.
type KlaytnRewardRaw struct {
	Contract *KlaytnReward // Generic contract binding to access the raw methods on
}

// KlaytnRewardCallerRaw is an auto generated low-level read-only Go binding around a Klaytn contract.
type KlaytnRewardCallerRaw struct {
	Contract *KlaytnRewardCaller // Generic read-only contract binding to access the raw methods on
}

// KlaytnRewardTransactorRaw is an auto generated low-level write-only Go binding around a Klaytn contract.
type KlaytnRewardTransactorRaw struct {
	Contract *KlaytnRewardTransactor // Generic write-only contract binding to access the raw methods on
}

// NewKlaytnReward creates a new instance of KlaytnReward, bound to a specific deployed contract.
func NewKlaytnReward(address common.Address, backend bind.ContractBackend) (*KlaytnReward, error) {
	contract, err := bindKlaytnReward(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &KlaytnReward{KlaytnRewardCaller: KlaytnRewardCaller{contract: contract}, KlaytnRewardTransactor: KlaytnRewardTransactor{contract: contract}, KlaytnRewardFilterer: KlaytnRewardFilterer{contract: contract}}, nil
}

// NewKlaytnRewardCaller creates a new read-only instance of KlaytnReward, bound to a specific deployed contract.
func NewKlaytnRewardCaller(address common.Address, caller bind.ContractCaller) (*KlaytnRewardCaller, error) {
	contract, err := bindKlaytnReward(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &KlaytnRewardCaller{contract: contract}, nil
}

// NewKlaytnRewardTransactor creates a new write-only instance of KlaytnReward, bound to a specific deployed contract.
func NewKlaytnRewardTransactor(address common.Address, transactor bind.ContractTransactor) (*KlaytnRewardTransactor, error) {
	contract, err := bindKlaytnReward(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &KlaytnRewardTransactor{contract: contract}, nil
}

// NewKlaytnRewardFilterer creates a new log filterer instance of KlaytnReward, bound to a specific deployed contract.
func NewKlaytnRewardFilterer(address common.Address, filterer bind.ContractFilterer) (*KlaytnRewardFilterer, error) {
	contract, err := bindKlaytnReward(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &KlaytnRewardFilterer{contract: contract}, nil
}

// bindKlaytnReward binds a generic wrapper to an already deployed contract.
func bindKlaytnReward(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(KlaytnRewardABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_KlaytnReward *KlaytnRewardRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _KlaytnReward.Contract.KlaytnRewardCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_KlaytnReward *KlaytnRewardRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KlaytnReward.Contract.KlaytnRewardTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_KlaytnReward *KlaytnRewardRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KlaytnReward.Contract.KlaytnRewardTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_KlaytnReward *KlaytnRewardCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _KlaytnReward.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_KlaytnReward *KlaytnRewardTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KlaytnReward.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_KlaytnReward *KlaytnRewardTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KlaytnReward.Contract.contract.Transact(opts, method, params...)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf( address) constant returns(uint256)
func (_KlaytnReward *KlaytnRewardCaller) BalanceOf(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _KlaytnReward.contract.Call(opts, out, "balanceOf", arg0)
	return *ret0, err
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf( address) constant returns(uint256)
func (_KlaytnReward *KlaytnRewardSession) BalanceOf(arg0 common.Address) (*big.Int, error) {
	return _KlaytnReward.Contract.BalanceOf(&_KlaytnReward.CallOpts, arg0)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf( address) constant returns(uint256)
func (_KlaytnReward *KlaytnRewardCallerSession) BalanceOf(arg0 common.Address) (*big.Int, error) {
	return _KlaytnReward.Contract.BalanceOf(&_KlaytnReward.CallOpts, arg0)
}

// TotalAmount is a free data retrieval call binding the contract method 0x1a39d8ef.
//
// Solidity: function totalAmount() constant returns(uint256)
func (_KlaytnReward *KlaytnRewardCaller) TotalAmount(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _KlaytnReward.contract.Call(opts, out, "totalAmount")
	return *ret0, err
}

// TotalAmount is a free data retrieval call binding the contract method 0x1a39d8ef.
//
// Solidity: function totalAmount() constant returns(uint256)
func (_KlaytnReward *KlaytnRewardSession) TotalAmount() (*big.Int, error) {
	return _KlaytnReward.Contract.TotalAmount(&_KlaytnReward.CallOpts)
}

// TotalAmount is a free data retrieval call binding the contract method 0x1a39d8ef.
//
// Solidity: function totalAmount() constant returns(uint256)
func (_KlaytnReward *KlaytnRewardCallerSession) TotalAmount() (*big.Int, error) {
	return _KlaytnReward.Contract.TotalAmount(&_KlaytnReward.CallOpts)
}

// Reward is a paid mutator transaction binding the contract method 0x6353586b.
//
// Solidity: function reward(receiver address) returns()
func (_KlaytnReward *KlaytnRewardTransactor) Reward(opts *bind.TransactOpts, receiver common.Address) (*types.Transaction, error) {
	return _KlaytnReward.contract.Transact(opts, "reward", receiver)
}

// Reward is a paid mutator transaction binding the contract method 0x6353586b.
//
// Solidity: function reward(receiver address) returns()
func (_KlaytnReward *KlaytnRewardSession) Reward(receiver common.Address) (*types.Transaction, error) {
	return _KlaytnReward.Contract.Reward(&_KlaytnReward.TransactOpts, receiver)
}

// Reward is a paid mutator transaction binding the contract method 0x6353586b.
//
// Solidity: function reward(receiver address) returns()
func (_KlaytnReward *KlaytnRewardTransactorSession) Reward(receiver common.Address) (*types.Transaction, error) {
	return _KlaytnReward.Contract.Reward(&_KlaytnReward.TransactOpts, receiver)
}

// SafeWithdrawal is a paid mutator transaction binding the contract method 0xfd6b7ef8.
//
// Solidity: function safeWithdrawal() returns()
func (_KlaytnReward *KlaytnRewardTransactor) SafeWithdrawal(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KlaytnReward.contract.Transact(opts, "safeWithdrawal")
}

// SafeWithdrawal is a paid mutator transaction binding the contract method 0xfd6b7ef8.
//
// Solidity: function safeWithdrawal() returns()
func (_KlaytnReward *KlaytnRewardSession) SafeWithdrawal() (*types.Transaction, error) {
	return _KlaytnReward.Contract.SafeWithdrawal(&_KlaytnReward.TransactOpts)
}

// SafeWithdrawal is a paid mutator transaction binding the contract method 0xfd6b7ef8.
//
// Solidity: function safeWithdrawal() returns()
func (_KlaytnReward *KlaytnRewardTransactorSession) SafeWithdrawal() (*types.Transaction, error) {
	return _KlaytnReward.Contract.SafeWithdrawal(&_KlaytnReward.TransactOpts)
}
