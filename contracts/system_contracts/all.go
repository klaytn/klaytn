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

// IKIP113BlsPublicKeyInfo is an auto generated low-level Go binding around an user-defined struct.
type IKIP113BlsPublicKeyInfo struct {
	PublicKey []byte
	Pop       []byte
}

// IKIP113ABI is the input ABI used to generate the binding from.
const IKIP113ABI = "[{\"inputs\":[],\"name\":\"getAllBlsInfo\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"nodeIdList\",\"type\":\"address[]\"},{\"components\":[{\"internalType\":\"bytes\",\"name\":\"publicKey\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"pop\",\"type\":\"bytes\"}],\"internalType\":\"structIKIP113.BlsPublicKeyInfo[]\",\"name\":\"pubkeyList\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"

// IKIP113BinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const IKIP113BinRuntime = ``

// IKIP113FuncSigs maps the 4-byte function signature to its string representation.
var IKIP113FuncSigs = map[string]string{
	"6968b53f": "getAllBlsInfo()",
}

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
	parsed, err := abi.JSON(strings.NewReader(IKIP113ABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IKIP113 *IKIP113Raw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
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
func (_IKIP113 *IKIP113CallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
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
}, error) {
	ret := new(struct {
		NodeIdList []common.Address
		PubkeyList []IKIP113BlsPublicKeyInfo
	})
	out := ret
	err := _IKIP113.contract.Call(opts, out, "getAllBlsInfo")
	return *ret, err
}

// GetAllBlsInfo is a free data retrieval call binding the contract method 0x6968b53f.
//
// Solidity: function getAllBlsInfo() view returns(address[] nodeIdList, (bytes,bytes)[] pubkeyList)
func (_IKIP113 *IKIP113Session) GetAllBlsInfo() (struct {
	NodeIdList []common.Address
	PubkeyList []IKIP113BlsPublicKeyInfo
}, error) {
	return _IKIP113.Contract.GetAllBlsInfo(&_IKIP113.CallOpts)
}

// GetAllBlsInfo is a free data retrieval call binding the contract method 0x6968b53f.
//
// Solidity: function getAllBlsInfo() view returns(address[] nodeIdList, (bytes,bytes)[] pubkeyList)
func (_IKIP113 *IKIP113CallerSession) GetAllBlsInfo() (struct {
	NodeIdList []common.Address
	PubkeyList []IKIP113BlsPublicKeyInfo
}, error) {
	return _IKIP113.Contract.GetAllBlsInfo(&_IKIP113.CallOpts)
}

// IRegistryABI is the input ABI used to generate the binding from.
const IRegistryABI = "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"name\":\"getActiveAddr\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"names\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"records\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"activation\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"

// IRegistryBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const IRegistryBinRuntime = ``

// IRegistryFuncSigs maps the 4-byte function signature to its string representation.
var IRegistryFuncSigs = map[string]string{
	"e2693e3f": "getActiveAddr(string)",
	"4622ab03": "names(uint256)",
	"3b51650d": "records(string,uint256)",
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

// KIP113MockABI is the input ABI used to generate the binding from.
const KIP113MockABI = "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"allNodeIds\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllBlsInfo\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"nodeIdList\",\"type\":\"address[]\"},{\"components\":[{\"internalType\":\"bytes\",\"name\":\"publicKey\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"pop\",\"type\":\"bytes\"}],\"internalType\":\"structIKIP113.BlsPublicKeyInfo[]\",\"name\":\"pubkeyList\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"record\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"publicKey\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"pop\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"publicKey\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"pop\",\"type\":\"bytes\"}],\"name\":\"register\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// KIP113MockBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const KIP113MockBinRuntime = `608060405234801561001057600080fd5b506004361061004c5760003560e01c80633465d6d5146100515780636968b53f1461007b578063786cd4d714610091578063a5834971146100a6575b600080fd5b61006461005f366004610631565b6100d1565b604051610072929190610699565b60405180910390f35b6100836101fd565b6040516100729291906106c7565b6100a461009f3660046107d5565b6104b7565b005b6100b96100b4366004610856565b6105eb565b6040516001600160a01b039091168152602001610072565b6000602081905290815260409020805481906100ec9061086f565b80601f01602080910402602001604051908101604052809291908181526020018280546101189061086f565b80156101655780601f1061013a57610100808354040283529160200191610165565b820191906000526020600020905b81548152906001019060200180831161014857829003601f168201915b50505050509080600101805461017a9061086f565b80601f01602080910402602001604051908101604052809291908181526020018280546101a69061086f565b80156101f35780601f106101c8576101008083540402835291602001916101f3565b820191906000526020600020905b8154815290600101906020018083116101d657829003601f168201915b5050505050905082565b60015460609081908067ffffffffffffffff81111561021e5761021e6108a9565b604051908082528060200260200182016040528015610247578160200160208202803683370190505b5092508067ffffffffffffffff811115610263576102636108a9565b6040519080825280602002602001820160405280156102a857816020015b60408051808201909152606080825260208201528152602001906001900390816102815790505b50915060005b818110156104b157600181815481106102c9576102c96108bf565b9060005260206000200160009054906101000a90046001600160a01b03168482815181106102f9576102f96108bf565b60200260200101906001600160a01b031690816001600160a01b0316815250506000806001838154811061032f5761032f6108bf565b60009182526020808320909101546001600160a01b031683528201929092526040908101909120815180830190925280548290829061036d9061086f565b80601f01602080910402602001604051908101604052809291908181526020018280546103999061086f565b80156103e65780601f106103bb576101008083540402835291602001916103e6565b820191906000526020600020905b8154815290600101906020018083116103c957829003601f168201915b505050505081526020016001820180546103ff9061086f565b80601f016020809104026020016040519081016040528092919081815260200182805461042b9061086f565b80156104785780601f1061044d57610100808354040283529160200191610478565b820191906000526020600020905b81548152906001019060200180831161045b57829003601f168201915b505050505081525050838281518110610493576104936108bf565b602002602001018190525080806104a9906108d5565b9150506102ae565b50509091565b6001600160a01b038516600090815260208190526040902080546104da9061086f565b905060000361052e576001805480820182556000919091527fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf60180546001600160a01b0319166001600160a01b0387161790555b6040805160606020601f87018190040282018101835291810185815290918291908790879081908501838280828437600092019190915250505090825250604080516020601f86018190048102820181019092528481529181019190859085908190840183828082843760009201829052509390945250506001600160a01b0388168152602081905260409020825190915081906105cc908261094b565b50602082015160018201906105e1908261094b565b5050505050505050565b600181815481106105fb57600080fd5b6000918252602090912001546001600160a01b0316905081565b80356001600160a01b038116811461062c57600080fd5b919050565b60006020828403121561064357600080fd5b61064c82610615565b9392505050565b6000815180845260005b818110156106795760208185018101518683018201520161065d565b506000602082860101526020601f19601f83011685010191505092915050565b6040815260006106ac6040830185610653565b82810360208401526106be8185610653565b95945050505050565b60408082528351828201819052600091906020906060850190828801855b8281101561070a5781516001600160a01b0316845292840192908401906001016106e5565b50505084810382860152855180825282820190600581901b8301840188850160005b8381101561077c57858303601f19018552815180518985526107508a860182610653565b91890151858303868b01529190506107688183610653565b96890196945050509086019060010161072c565b50909a9950505050505050505050565b60008083601f84011261079e57600080fd5b50813567ffffffffffffffff8111156107b657600080fd5b6020830191508360208285010111156107ce57600080fd5b9250929050565b6000806000806000606086880312156107ed57600080fd5b6107f686610615565b9450602086013567ffffffffffffffff8082111561081357600080fd5b61081f89838a0161078c565b9096509450604088013591508082111561083857600080fd5b506108458882890161078c565b969995985093965092949392505050565b60006020828403121561086857600080fd5b5035919050565b600181811c9082168061088357607f821691505b6020821081036108a357634e487b7160e01b600052602260045260246000fd5b50919050565b634e487b7160e01b600052604160045260246000fd5b634e487b7160e01b600052603260045260246000fd5b6000600182016108f557634e487b7160e01b600052601160045260246000fd5b5060010190565b601f82111561094657600081815260208120601f850160051c810160208610156109235750805b601f850160051c820191505b818110156109425782815560010161092f565b5050505b505050565b815167ffffffffffffffff811115610965576109656108a9565b61097981610973845461086f565b846108fc565b602080601f8311600181146109ae57600084156109965750858301515b600019600386901b1c1916600185901b178555610942565b600085815260208120601f198616915b828110156109dd578886015182559484019460019091019084016109be565b50858210156109fb5787850151600019600388901b60f8161c191681555b5050505050600190811b0190555056fea26469706673582212204113607f3e91773b3f06c6cb723f751b7dfbc7877c631e26f3df94555695da0d64736f6c63430008110033`

// KIP113MockFuncSigs maps the 4-byte function signature to its string representation.
var KIP113MockFuncSigs = map[string]string{
	"a5834971": "allNodeIds(uint256)",
	"6968b53f": "getAllBlsInfo()",
	"3465d6d5": "record(address)",
	"786cd4d7": "register(address,bytes,bytes)",
}

// KIP113MockBin is the compiled bytecode used for deploying new contracts.
var KIP113MockBin = "0x608060405234801561001057600080fd5b50610a41806100206000396000f3fe608060405234801561001057600080fd5b506004361061004c5760003560e01c80633465d6d5146100515780636968b53f1461007b578063786cd4d714610091578063a5834971146100a6575b600080fd5b61006461005f366004610631565b6100d1565b604051610072929190610699565b60405180910390f35b6100836101fd565b6040516100729291906106c7565b6100a461009f3660046107d5565b6104b7565b005b6100b96100b4366004610856565b6105eb565b6040516001600160a01b039091168152602001610072565b6000602081905290815260409020805481906100ec9061086f565b80601f01602080910402602001604051908101604052809291908181526020018280546101189061086f565b80156101655780601f1061013a57610100808354040283529160200191610165565b820191906000526020600020905b81548152906001019060200180831161014857829003601f168201915b50505050509080600101805461017a9061086f565b80601f01602080910402602001604051908101604052809291908181526020018280546101a69061086f565b80156101f35780601f106101c8576101008083540402835291602001916101f3565b820191906000526020600020905b8154815290600101906020018083116101d657829003601f168201915b5050505050905082565b60015460609081908067ffffffffffffffff81111561021e5761021e6108a9565b604051908082528060200260200182016040528015610247578160200160208202803683370190505b5092508067ffffffffffffffff811115610263576102636108a9565b6040519080825280602002602001820160405280156102a857816020015b60408051808201909152606080825260208201528152602001906001900390816102815790505b50915060005b818110156104b157600181815481106102c9576102c96108bf565b9060005260206000200160009054906101000a90046001600160a01b03168482815181106102f9576102f96108bf565b60200260200101906001600160a01b031690816001600160a01b0316815250506000806001838154811061032f5761032f6108bf565b60009182526020808320909101546001600160a01b031683528201929092526040908101909120815180830190925280548290829061036d9061086f565b80601f01602080910402602001604051908101604052809291908181526020018280546103999061086f565b80156103e65780601f106103bb576101008083540402835291602001916103e6565b820191906000526020600020905b8154815290600101906020018083116103c957829003601f168201915b505050505081526020016001820180546103ff9061086f565b80601f016020809104026020016040519081016040528092919081815260200182805461042b9061086f565b80156104785780601f1061044d57610100808354040283529160200191610478565b820191906000526020600020905b81548152906001019060200180831161045b57829003601f168201915b505050505081525050838281518110610493576104936108bf565b602002602001018190525080806104a9906108d5565b9150506102ae565b50509091565b6001600160a01b038516600090815260208190526040902080546104da9061086f565b905060000361052e576001805480820182556000919091527fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf60180546001600160a01b0319166001600160a01b0387161790555b6040805160606020601f87018190040282018101835291810185815290918291908790879081908501838280828437600092019190915250505090825250604080516020601f86018190048102820181019092528481529181019190859085908190840183828082843760009201829052509390945250506001600160a01b0388168152602081905260409020825190915081906105cc908261094b565b50602082015160018201906105e1908261094b565b5050505050505050565b600181815481106105fb57600080fd5b6000918252602090912001546001600160a01b0316905081565b80356001600160a01b038116811461062c57600080fd5b919050565b60006020828403121561064357600080fd5b61064c82610615565b9392505050565b6000815180845260005b818110156106795760208185018101518683018201520161065d565b506000602082860101526020601f19601f83011685010191505092915050565b6040815260006106ac6040830185610653565b82810360208401526106be8185610653565b95945050505050565b60408082528351828201819052600091906020906060850190828801855b8281101561070a5781516001600160a01b0316845292840192908401906001016106e5565b50505084810382860152855180825282820190600581901b8301840188850160005b8381101561077c57858303601f19018552815180518985526107508a860182610653565b91890151858303868b01529190506107688183610653565b96890196945050509086019060010161072c565b50909a9950505050505050505050565b60008083601f84011261079e57600080fd5b50813567ffffffffffffffff8111156107b657600080fd5b6020830191508360208285010111156107ce57600080fd5b9250929050565b6000806000806000606086880312156107ed57600080fd5b6107f686610615565b9450602086013567ffffffffffffffff8082111561081357600080fd5b61081f89838a0161078c565b9096509450604088013591508082111561083857600080fd5b506108458882890161078c565b969995985093965092949392505050565b60006020828403121561086857600080fd5b5035919050565b600181811c9082168061088357607f821691505b6020821081036108a357634e487b7160e01b600052602260045260246000fd5b50919050565b634e487b7160e01b600052604160045260246000fd5b634e487b7160e01b600052603260045260246000fd5b6000600182016108f557634e487b7160e01b600052601160045260246000fd5b5060010190565b601f82111561094657600081815260208120601f850160051c810160208610156109235750805b601f850160051c820191505b818110156109425782815560010161092f565b5050505b505050565b815167ffffffffffffffff811115610965576109656108a9565b61097981610973845461086f565b846108fc565b602080601f8311600181146109ae57600084156109965750858301515b600019600386901b1c1916600185901b178555610942565b600085815260208120601f198616915b828110156109dd578886015182559484019460019091019084016109be565b50858210156109fb5787850151600019600388901b60f8161c191681555b5050505050600190811b0190555056fea26469706673582212204113607f3e91773b3f06c6cb723f751b7dfbc7877c631e26f3df94555695da0d64736f6c63430008110033"

// DeployKIP113Mock deploys a new Klaytn contract, binding an instance of KIP113Mock to it.
func DeployKIP113Mock(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *KIP113Mock, error) {
	parsed, err := abi.JSON(strings.NewReader(KIP113MockABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(KIP113MockBin), backend)
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
	parsed, err := abi.JSON(strings.NewReader(KIP113MockABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_KIP113Mock *KIP113MockRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
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
func (_KIP113Mock *KIP113MockCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
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

// AllNodeIds is a free data retrieval call binding the contract method 0xa5834971.
//
// Solidity: function allNodeIds(uint256 ) view returns(address)
func (_KIP113Mock *KIP113MockCaller) AllNodeIds(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _KIP113Mock.contract.Call(opts, out, "allNodeIds", arg0)
	return *ret0, err
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
}, error) {
	ret := new(struct {
		NodeIdList []common.Address
		PubkeyList []IKIP113BlsPublicKeyInfo
	})
	out := ret
	err := _KIP113Mock.contract.Call(opts, out, "getAllBlsInfo")
	return *ret, err
}

// GetAllBlsInfo is a free data retrieval call binding the contract method 0x6968b53f.
//
// Solidity: function getAllBlsInfo() view returns(address[] nodeIdList, (bytes,bytes)[] pubkeyList)
func (_KIP113Mock *KIP113MockSession) GetAllBlsInfo() (struct {
	NodeIdList []common.Address
	PubkeyList []IKIP113BlsPublicKeyInfo
}, error) {
	return _KIP113Mock.Contract.GetAllBlsInfo(&_KIP113Mock.CallOpts)
}

// GetAllBlsInfo is a free data retrieval call binding the contract method 0x6968b53f.
//
// Solidity: function getAllBlsInfo() view returns(address[] nodeIdList, (bytes,bytes)[] pubkeyList)
func (_KIP113Mock *KIP113MockCallerSession) GetAllBlsInfo() (struct {
	NodeIdList []common.Address
	PubkeyList []IKIP113BlsPublicKeyInfo
}, error) {
	return _KIP113Mock.Contract.GetAllBlsInfo(&_KIP113Mock.CallOpts)
}

// Record is a free data retrieval call binding the contract method 0x3465d6d5.
//
// Solidity: function record(address ) view returns(bytes publicKey, bytes pop)
func (_KIP113Mock *KIP113MockCaller) Record(opts *bind.CallOpts, arg0 common.Address) (struct {
	PublicKey []byte
	Pop       []byte
}, error) {
	ret := new(struct {
		PublicKey []byte
		Pop       []byte
	})
	out := ret
	err := _KIP113Mock.contract.Call(opts, out, "record", arg0)
	return *ret, err
}

// Record is a free data retrieval call binding the contract method 0x3465d6d5.
//
// Solidity: function record(address ) view returns(bytes publicKey, bytes pop)
func (_KIP113Mock *KIP113MockSession) Record(arg0 common.Address) (struct {
	PublicKey []byte
	Pop       []byte
}, error) {
	return _KIP113Mock.Contract.Record(&_KIP113Mock.CallOpts, arg0)
}

// Record is a free data retrieval call binding the contract method 0x3465d6d5.
//
// Solidity: function record(address ) view returns(bytes publicKey, bytes pop)
func (_KIP113Mock *KIP113MockCallerSession) Record(arg0 common.Address) (struct {
	PublicKey []byte
	Pop       []byte
}, error) {
	return _KIP113Mock.Contract.Record(&_KIP113Mock.CallOpts, arg0)
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

// RegistryABI is the input ABI used to generate the binding from.
const RegistryABI = "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"name\":\"getActiveAddr\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"names\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"records\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"activation\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"

// RegistryBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const RegistryBinRuntime = `608060405234801561001057600080fd5b50600436106100415760003560e01c80633b51650d146100465780634622ab031461007d578063e2693e3f1461009d575b600080fd5b6100596100543660046102ad565b6100c8565b604080516001600160a01b0390931683526020830191909152015b60405180910390f35b61009061008b3660046102f2565b61011d565b604051610074919061030b565b6100b06100ab366004610359565b6101c9565b6040516001600160a01b039091168152602001610074565b815160208184018101805160008252928201918501919091209190528054829081106100f357600080fd5b6000918252602090912060029091020180546001909101546001600160a01b039091169250905082565b6001818154811061012d57600080fd5b90600052602060002001600091509050805461014890610396565b80601f016020809104026020016040519081016040528092919081815260200182805461017490610396565b80156101c15780601f10610196576101008083540402835291602001916101c1565b820191906000526020600020905b8154815290600101906020018083116101a457829003601f168201915b505050505081565b60405162461bcd60e51b815260206004820152600f60248201526e1b9bdd081a5b5c1b195b595b9d1959608a1b604482015260009060640160405180910390fd5b634e487b7160e01b600052604160045260246000fd5b600082601f83011261023157600080fd5b813567ffffffffffffffff8082111561024c5761024c61020a565b604051601f8301601f19908116603f011681019082821181831017156102745761027461020a565b8160405283815286602085880101111561028d57600080fd5b836020870160208301376000602085830101528094505050505092915050565b600080604083850312156102c057600080fd5b823567ffffffffffffffff8111156102d757600080fd5b6102e385828601610220565b95602094909401359450505050565b60006020828403121561030457600080fd5b5035919050565b600060208083528351808285015260005b818110156103385785810183015185820160400152820161031c565b506000604082860101526040601f19601f8301168501019250505092915050565b60006020828403121561036b57600080fd5b813567ffffffffffffffff81111561038257600080fd5b61038e84828501610220565b949350505050565b600181811c908216806103aa57607f821691505b6020821081036103ca57634e487b7160e01b600052602260045260246000fd5b5091905056fea264697066735822122081d8382cad173f430e0005a96faaba1ff204fe72b466874499e03026c424a33964736f6c63430008110033`

// RegistryFuncSigs maps the 4-byte function signature to its string representation.
var RegistryFuncSigs = map[string]string{
	"e2693e3f": "getActiveAddr(string)",
	"4622ab03": "names(uint256)",
	"3b51650d": "records(string,uint256)",
}

// RegistryBin is the compiled bytecode used for deploying new contracts.
var RegistryBin = "0x608060405234801561001057600080fd5b50610406806100206000396000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c80633b51650d146100465780634622ab031461007d578063e2693e3f1461009d575b600080fd5b6100596100543660046102ad565b6100c8565b604080516001600160a01b0390931683526020830191909152015b60405180910390f35b61009061008b3660046102f2565b61011d565b604051610074919061030b565b6100b06100ab366004610359565b6101c9565b6040516001600160a01b039091168152602001610074565b815160208184018101805160008252928201918501919091209190528054829081106100f357600080fd5b6000918252602090912060029091020180546001909101546001600160a01b039091169250905082565b6001818154811061012d57600080fd5b90600052602060002001600091509050805461014890610396565b80601f016020809104026020016040519081016040528092919081815260200182805461017490610396565b80156101c15780601f10610196576101008083540402835291602001916101c1565b820191906000526020600020905b8154815290600101906020018083116101a457829003601f168201915b505050505081565b60405162461bcd60e51b815260206004820152600f60248201526e1b9bdd081a5b5c1b195b595b9d1959608a1b604482015260009060640160405180910390fd5b634e487b7160e01b600052604160045260246000fd5b600082601f83011261023157600080fd5b813567ffffffffffffffff8082111561024c5761024c61020a565b604051601f8301601f19908116603f011681019082821181831017156102745761027461020a565b8160405283815286602085880101111561028d57600080fd5b836020870160208301376000602085830101528094505050505092915050565b600080604083850312156102c057600080fd5b823567ffffffffffffffff8111156102d757600080fd5b6102e385828601610220565b95602094909401359450505050565b60006020828403121561030457600080fd5b5035919050565b600060208083528351808285015260005b818110156103385785810183015185820160400152820161031c565b506000604082860101526040601f19601f8301168501019250505092915050565b60006020828403121561036b57600080fd5b813567ffffffffffffffff81111561038257600080fd5b61038e84828501610220565b949350505050565b600181811c908216806103aa57607f821691505b6020821081036103ca57634e487b7160e01b600052602260045260246000fd5b5091905056fea264697066735822122081d8382cad173f430e0005a96faaba1ff204fe72b466874499e03026c424a33964736f6c63430008110033"

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

// RegistryMockABI is the input ABI used to generate the binding from.
const RegistryMockABI = "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"name\":\"getActiveAddr\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"names\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"records\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"activation\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"activation\",\"type\":\"uint256\"}],\"name\":\"register\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// RegistryMockBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const RegistryMockBinRuntime = `608060405234801561001057600080fd5b50600436106100575760003560e01c80633b51650d1461005c5780634622ab0314610093578063d393c871146100b3578063e2693e3f146100c8578063f2fde38b146100f3575b600080fd5b61006f61006a366004610432565b610123565b604080516001600160a01b0390931683526020830191909152015b60405180910390f35b6100a66100a1366004610477565b610178565b60405161008a91906104b4565b6100c66100c1366004610503565b610224565b005b6100db6100d636600461055a565b6102f7565b6040516001600160a01b03909116815260200161008a565b6100c6610101366004610597565b600280546001600160a01b0319166001600160a01b0392909216919091179055565b8151602081840181018051600082529282019185019190912091905280548290811061014e57600080fd5b6000918252602090912060029091020180546001909101546001600160a01b039091169250905082565b6001818154811061018857600080fd5b9060005260206000200160009150905080546101a3906105b9565b80601f01602080910402602001604051908101604052809291908181526020018280546101cf906105b9565b801561021c5780601f106101f15761010080835404028352916020019161021c565b820191906000526020600020905b8154815290600101906020018083116101ff57829003601f168201915b505050505081565b60008360405161023491906105ed565b90815260405190819003602001902054600003610288576001805480820182556000919091527fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf6016102868482610658565b505b60008360405161029891906105ed565b90815260408051602092819003830181208183019092526001600160a01b039485168152828101938452815460018082018455600093845293909220905160029092020180546001600160a01b03191691909416178355905191015550565b60008060008360405161030a91906105ed565b908152604051908190036020019020549050600081900361032e5750600092915050565b60008360405161033e91906105ed565b908152604051908190036020019020610358600183610718565b815481106103685761036861073f565b60009182526020909120600290910201546001600160a01b03169392505050565b50919050565b634e487b7160e01b600052604160045260246000fd5b600082601f8301126103b657600080fd5b813567ffffffffffffffff808211156103d1576103d161038f565b604051601f8301601f19908116603f011681019082821181831017156103f9576103f961038f565b8160405283815286602085880101111561041257600080fd5b836020870160208301376000602085830101528094505050505092915050565b6000806040838503121561044557600080fd5b823567ffffffffffffffff81111561045c57600080fd5b610468858286016103a5565b95602094909401359450505050565b60006020828403121561048957600080fd5b5035919050565b60005b838110156104ab578181015183820152602001610493565b50506000910152565b60208152600082518060208401526104d3816040850160208701610490565b601f01601f19169190910160400192915050565b80356001600160a01b03811681146104fe57600080fd5b919050565b60008060006060848603121561051857600080fd5b833567ffffffffffffffff81111561052f57600080fd5b61053b868287016103a5565b93505061054a602085016104e7565b9150604084013590509250925092565b60006020828403121561056c57600080fd5b813567ffffffffffffffff81111561058357600080fd5b61058f848285016103a5565b949350505050565b6000602082840312156105a957600080fd5b6105b2826104e7565b9392505050565b600181811c908216806105cd57607f821691505b60208210810361038957634e487b7160e01b600052602260045260246000fd5b600082516105ff818460208701610490565b9190910192915050565b601f82111561065357600081815260208120601f850160051c810160208610156106305750805b601f850160051c820191505b8181101561064f5782815560010161063c565b5050505b505050565b815167ffffffffffffffff8111156106725761067261038f565b6106868161068084546105b9565b84610609565b602080601f8311600181146106bb57600084156106a35750858301515b600019600386901b1c1916600185901b17855561064f565b600085815260208120601f198616915b828110156106ea578886015182559484019460019091019084016106cb565b50858210156107085787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b8181038181111561073957634e487b7160e01b600052601160045260246000fd5b92915050565b634e487b7160e01b600052603260045260246000fdfea2646970667358221220614d57a0d78769a7f0862854ec412095f7372467ccb4c829a29ffd97a72e5b1464736f6c63430008110033`

// RegistryMockFuncSigs maps the 4-byte function signature to its string representation.
var RegistryMockFuncSigs = map[string]string{
	"e2693e3f": "getActiveAddr(string)",
	"4622ab03": "names(uint256)",
	"3b51650d": "records(string,uint256)",
	"d393c871": "register(string,address,uint256)",
	"f2fde38b": "transferOwnership(address)",
}

// RegistryMockBin is the compiled bytecode used for deploying new contracts.
var RegistryMockBin = "0x608060405234801561001057600080fd5b5061078b806100206000396000f3fe608060405234801561001057600080fd5b50600436106100575760003560e01c80633b51650d1461005c5780634622ab0314610093578063d393c871146100b3578063e2693e3f146100c8578063f2fde38b146100f3575b600080fd5b61006f61006a366004610432565b610123565b604080516001600160a01b0390931683526020830191909152015b60405180910390f35b6100a66100a1366004610477565b610178565b60405161008a91906104b4565b6100c66100c1366004610503565b610224565b005b6100db6100d636600461055a565b6102f7565b6040516001600160a01b03909116815260200161008a565b6100c6610101366004610597565b600280546001600160a01b0319166001600160a01b0392909216919091179055565b8151602081840181018051600082529282019185019190912091905280548290811061014e57600080fd5b6000918252602090912060029091020180546001909101546001600160a01b039091169250905082565b6001818154811061018857600080fd5b9060005260206000200160009150905080546101a3906105b9565b80601f01602080910402602001604051908101604052809291908181526020018280546101cf906105b9565b801561021c5780601f106101f15761010080835404028352916020019161021c565b820191906000526020600020905b8154815290600101906020018083116101ff57829003601f168201915b505050505081565b60008360405161023491906105ed565b90815260405190819003602001902054600003610288576001805480820182556000919091527fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf6016102868482610658565b505b60008360405161029891906105ed565b90815260408051602092819003830181208183019092526001600160a01b039485168152828101938452815460018082018455600093845293909220905160029092020180546001600160a01b03191691909416178355905191015550565b60008060008360405161030a91906105ed565b908152604051908190036020019020549050600081900361032e5750600092915050565b60008360405161033e91906105ed565b908152604051908190036020019020610358600183610718565b815481106103685761036861073f565b60009182526020909120600290910201546001600160a01b03169392505050565b50919050565b634e487b7160e01b600052604160045260246000fd5b600082601f8301126103b657600080fd5b813567ffffffffffffffff808211156103d1576103d161038f565b604051601f8301601f19908116603f011681019082821181831017156103f9576103f961038f565b8160405283815286602085880101111561041257600080fd5b836020870160208301376000602085830101528094505050505092915050565b6000806040838503121561044557600080fd5b823567ffffffffffffffff81111561045c57600080fd5b610468858286016103a5565b95602094909401359450505050565b60006020828403121561048957600080fd5b5035919050565b60005b838110156104ab578181015183820152602001610493565b50506000910152565b60208152600082518060208401526104d3816040850160208701610490565b601f01601f19169190910160400192915050565b80356001600160a01b03811681146104fe57600080fd5b919050565b60008060006060848603121561051857600080fd5b833567ffffffffffffffff81111561052f57600080fd5b61053b868287016103a5565b93505061054a602085016104e7565b9150604084013590509250925092565b60006020828403121561056c57600080fd5b813567ffffffffffffffff81111561058357600080fd5b61058f848285016103a5565b949350505050565b6000602082840312156105a957600080fd5b6105b2826104e7565b9392505050565b600181811c908216806105cd57607f821691505b60208210810361038957634e487b7160e01b600052602260045260246000fd5b600082516105ff818460208701610490565b9190910192915050565b601f82111561065357600081815260208120601f850160051c810160208610156106305750805b601f850160051c820191505b8181101561064f5782815560010161063c565b5050505b505050565b815167ffffffffffffffff8111156106725761067261038f565b6106868161068084546105b9565b84610609565b602080601f8311600181146106bb57600084156106a35750858301515b600019600386901b1c1916600185901b17855561064f565b600085815260208120601f198616915b828110156106ea578886015182559484019460019091019084016106cb565b50858210156107085787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b8181038181111561073957634e487b7160e01b600052601160045260246000fd5b92915050565b634e487b7160e01b600052603260045260246000fdfea2646970667358221220614d57a0d78769a7f0862854ec412095f7372467ccb4c829a29ffd97a72e5b1464736f6c63430008110033"

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
