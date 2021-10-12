// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package MintBurnWithPermission

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

// MintBurnTestABI is the input ABI used to generate the binding from.
const MintBurnTestABI = "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"mint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// MintBurnTestBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const MintBurnTestBinRuntime = `608060405234801561001057600080fd5b50600436106100365760003560e01c806340c10f191461003b5780639dc29fac14610050575b600080fd5b61004e610049366004610393565b610063565b005b61004e61005e366004610393565b610201565b6040516bffffffffffffffffffffffff19606084901b166020820152603481018290526000906054016040516020818303038152906040529050600060606103fb6001600160a01b0316836040516100bb91906103cb565b6000604051808303816000865af19150503d80600081146100f8576040519150601f19603f3d011682016040523d82523d6000602084013e6100fd565b606091505b5060405191935091506009906101149083906103cb565b6000604051808303816000865af19150503d8060008114610151576040519150601f19603f3d011682016040523d82523d6000602084013e610156565b606091505b505050816101ab5760405162461bcd60e51b815260206004820152601f60248201527f546865726520697320616e206572726f72207768696c65206d696e74696e670060448201526064015b60405180910390fd5b8051156101fa5760405162461bcd60e51b815260206004820152601f60248201527f546865726520697320616e206572726f72207768696c65206d696e74696e670060448201526064016101a2565b5050505050565b6040516bffffffffffffffffffffffff19606084901b166020820152603481018290526000906054016040516020818303038152906040529050600060606103fc6001600160a01b03168360405161025991906103cb565b6000604051808303816000865af19150503d8060008114610296576040519150601f19603f3d011682016040523d82523d6000602084013e61029b565b606091505b5060405191935091506009906102b29083906103cb565b6000604051808303816000865af19150503d80600081146102ef576040519150601f19603f3d011682016040523d82523d6000602084013e6102f4565b606091505b505050816103445760405162461bcd60e51b815260206004820152601f60248201527f546865726520697320616e206572726f72207768696c65206275726e696e670060448201526064016101a2565b8051156101fa5760405162461bcd60e51b815260206004820152601f60248201527f546865726520697320616e206572726f72207768696c65206275726e696e670060448201526064016101a2565b600080604083850312156103a657600080fd5b82356001600160a01b03811681146103bd57600080fd5b946020939093013593505050565b6000825160005b818110156103ec57602081860181015185830152016103d2565b818111156103fb576000828501525b50919091019291505056fea2646970667358221220aa3230aa2646271a18711d1d0d3713c5169de2d2ffef0b6f06904411083afb1464736f6c63430008070033`

// MintBurnTestFuncSigs maps the 4-byte function signature to its string representation.
var MintBurnTestFuncSigs = map[string]string{
	"9dc29fac": "burn(address,uint256)",
	"40c10f19": "mint(address,uint256)",
}

// MintBurnTestBin is the compiled bytecode used for deploying new contracts.
var MintBurnTestBin = "0x608060405234801561001057600080fd5b5061043c806100206000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c806340c10f191461003b5780639dc29fac14610050575b600080fd5b61004e610049366004610393565b610063565b005b61004e61005e366004610393565b610201565b6040516bffffffffffffffffffffffff19606084901b166020820152603481018290526000906054016040516020818303038152906040529050600060606103fb6001600160a01b0316836040516100bb91906103cb565b6000604051808303816000865af19150503d80600081146100f8576040519150601f19603f3d011682016040523d82523d6000602084013e6100fd565b606091505b5060405191935091506009906101149083906103cb565b6000604051808303816000865af19150503d8060008114610151576040519150601f19603f3d011682016040523d82523d6000602084013e610156565b606091505b505050816101ab5760405162461bcd60e51b815260206004820152601f60248201527f546865726520697320616e206572726f72207768696c65206d696e74696e670060448201526064015b60405180910390fd5b8051156101fa5760405162461bcd60e51b815260206004820152601f60248201527f546865726520697320616e206572726f72207768696c65206d696e74696e670060448201526064016101a2565b5050505050565b6040516bffffffffffffffffffffffff19606084901b166020820152603481018290526000906054016040516020818303038152906040529050600060606103fc6001600160a01b03168360405161025991906103cb565b6000604051808303816000865af19150503d8060008114610296576040519150601f19603f3d011682016040523d82523d6000602084013e61029b565b606091505b5060405191935091506009906102b29083906103cb565b6000604051808303816000865af19150503d80600081146102ef576040519150601f19603f3d011682016040523d82523d6000602084013e6102f4565b606091505b505050816103445760405162461bcd60e51b815260206004820152601f60248201527f546865726520697320616e206572726f72207768696c65206275726e696e670060448201526064016101a2565b8051156101fa5760405162461bcd60e51b815260206004820152601f60248201527f546865726520697320616e206572726f72207768696c65206275726e696e670060448201526064016101a2565b600080604083850312156103a657600080fd5b82356001600160a01b03811681146103bd57600080fd5b946020939093013593505050565b6000825160005b818110156103ec57602081860181015185830152016103d2565b818111156103fb576000828501525b50919091019291505056fea2646970667358221220aa3230aa2646271a18711d1d0d3713c5169de2d2ffef0b6f06904411083afb1464736f6c63430008070033"

// DeployMintBurnTest deploys a new Klaytn contract, binding an instance of MintBurnTest to it.
func DeployMintBurnTest(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *MintBurnTest, error) {
	parsed, err := abi.JSON(strings.NewReader(MintBurnTestABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(MintBurnTestBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &MintBurnTest{MintBurnTestCaller: MintBurnTestCaller{contract: contract}, MintBurnTestTransactor: MintBurnTestTransactor{contract: contract}, MintBurnTestFilterer: MintBurnTestFilterer{contract: contract}}, nil
}

// MintBurnTest is an auto generated Go binding around a Klaytn contract.
type MintBurnTest struct {
	MintBurnTestCaller     // Read-only binding to the contract
	MintBurnTestTransactor // Write-only binding to the contract
	MintBurnTestFilterer   // Log filterer for contract events
}

// MintBurnTestCaller is an auto generated read-only Go binding around a Klaytn contract.
type MintBurnTestCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MintBurnTestTransactor is an auto generated write-only Go binding around a Klaytn contract.
type MintBurnTestTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MintBurnTestFilterer is an auto generated log filtering Go binding around a Klaytn contract events.
type MintBurnTestFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MintBurnTestSession is an auto generated Go binding around a Klaytn contract,
// with pre-set call and transact options.
type MintBurnTestSession struct {
	Contract     *MintBurnTest     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// MintBurnTestCallerSession is an auto generated read-only Go binding around a Klaytn contract,
// with pre-set call options.
type MintBurnTestCallerSession struct {
	Contract *MintBurnTestCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// MintBurnTestTransactorSession is an auto generated write-only Go binding around a Klaytn contract,
// with pre-set transact options.
type MintBurnTestTransactorSession struct {
	Contract     *MintBurnTestTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// MintBurnTestRaw is an auto generated low-level Go binding around a Klaytn contract.
type MintBurnTestRaw struct {
	Contract *MintBurnTest // Generic contract binding to access the raw methods on
}

// MintBurnTestCallerRaw is an auto generated low-level read-only Go binding around a Klaytn contract.
type MintBurnTestCallerRaw struct {
	Contract *MintBurnTestCaller // Generic read-only contract binding to access the raw methods on
}

// MintBurnTestTransactorRaw is an auto generated low-level write-only Go binding around a Klaytn contract.
type MintBurnTestTransactorRaw struct {
	Contract *MintBurnTestTransactor // Generic write-only contract binding to access the raw methods on
}

// NewMintBurnTest creates a new instance of MintBurnTest, bound to a specific deployed contract.
func NewMintBurnTest(address common.Address, backend bind.ContractBackend) (*MintBurnTest, error) {
	contract, err := bindMintBurnTest(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &MintBurnTest{MintBurnTestCaller: MintBurnTestCaller{contract: contract}, MintBurnTestTransactor: MintBurnTestTransactor{contract: contract}, MintBurnTestFilterer: MintBurnTestFilterer{contract: contract}}, nil
}

// NewMintBurnTestCaller creates a new read-only instance of MintBurnTest, bound to a specific deployed contract.
func NewMintBurnTestCaller(address common.Address, caller bind.ContractCaller) (*MintBurnTestCaller, error) {
	contract, err := bindMintBurnTest(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MintBurnTestCaller{contract: contract}, nil
}

// NewMintBurnTestTransactor creates a new write-only instance of MintBurnTest, bound to a specific deployed contract.
func NewMintBurnTestTransactor(address common.Address, transactor bind.ContractTransactor) (*MintBurnTestTransactor, error) {
	contract, err := bindMintBurnTest(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MintBurnTestTransactor{contract: contract}, nil
}

// NewMintBurnTestFilterer creates a new log filterer instance of MintBurnTest, bound to a specific deployed contract.
func NewMintBurnTestFilterer(address common.Address, filterer bind.ContractFilterer) (*MintBurnTestFilterer, error) {
	contract, err := bindMintBurnTest(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MintBurnTestFilterer{contract: contract}, nil
}

// bindMintBurnTest binds a generic wrapper to an already deployed contract.
func bindMintBurnTest(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(MintBurnTestABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MintBurnTest *MintBurnTestRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _MintBurnTest.Contract.MintBurnTestCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MintBurnTest *MintBurnTestRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MintBurnTest.Contract.MintBurnTestTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MintBurnTest *MintBurnTestRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MintBurnTest.Contract.MintBurnTestTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MintBurnTest *MintBurnTestCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _MintBurnTest.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MintBurnTest *MintBurnTestTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MintBurnTest.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MintBurnTest *MintBurnTestTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MintBurnTest.Contract.contract.Transact(opts, method, params...)
}

// Burn is a paid mutator transaction binding the contract method 0x9dc29fac.
//
// Solidity: function burn(address account, uint256 amount) returns()
func (_MintBurnTest *MintBurnTestTransactor) Burn(opts *bind.TransactOpts, account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _MintBurnTest.contract.Transact(opts, "burn", account, amount)
}

// Burn is a paid mutator transaction binding the contract method 0x9dc29fac.
//
// Solidity: function burn(address account, uint256 amount) returns()
func (_MintBurnTest *MintBurnTestSession) Burn(account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _MintBurnTest.Contract.Burn(&_MintBurnTest.TransactOpts, account, amount)
}

// Burn is a paid mutator transaction binding the contract method 0x9dc29fac.
//
// Solidity: function burn(address account, uint256 amount) returns()
func (_MintBurnTest *MintBurnTestTransactorSession) Burn(account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _MintBurnTest.Contract.Burn(&_MintBurnTest.TransactOpts, account, amount)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address account, uint256 amount) returns()
func (_MintBurnTest *MintBurnTestTransactor) Mint(opts *bind.TransactOpts, account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _MintBurnTest.contract.Transact(opts, "mint", account, amount)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address account, uint256 amount) returns()
func (_MintBurnTest *MintBurnTestSession) Mint(account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _MintBurnTest.Contract.Mint(&_MintBurnTest.TransactOpts, account, amount)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address account, uint256 amount) returns()
func (_MintBurnTest *MintBurnTestTransactorSession) Mint(account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _MintBurnTest.Contract.Mint(&_MintBurnTest.TransactOpts, account, amount)
}
