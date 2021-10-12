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

// MintBurnWithPermissionABI is the input ABI used to generate the binding from.
const MintBurnWithPermissionABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"burner\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"burnee\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"burnAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"accBurnAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"totalSupply\",\"type\":\"uint256\"}],\"name\":\"Burn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"minter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"mintee\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"mintAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"accMintAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"totalSupply\",\"type\":\"uint256\"}],\"name\":\"Mint\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"burnt\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBurnee\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBurner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getMinter\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"mint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"minted\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"resetBurnAmount\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"resetMintAmount\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"setBurnee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"setBurner\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"setMinter\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"

// MintBurnWithPermissionBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const MintBurnWithPermissionBinRuntime = `608060405234801561001057600080fd5b50600436106100cf5760003560e01c8063b192da2d1161008c578063d773cca911610066578063d773cca91461015d578063e236b8e214610182578063f366751714610193578063fca3b5aa146101a457600080fd5b8063b192da2d1461013a578063bf6322c314610142578063c46d4c0f1461015557600080fd5b806318160ddd146100d457806340c10f19146100ef57806342966c68146101045780634f02c42014610117578063a996d6ce1461011f578063ace8923414610132575b600080fd5b6100dc6101b7565b6040519081526020015b60405180910390f35b6101026100fd3660046107e9565b6101ce565b005b610102610112366004610813565b61041d565b6003546100dc565b61010261012d3660046107c7565b610665565b6101026106b1565b6004546100dc565b6101026101503660046107c7565b6106e2565b61010261072e565b6002546001600160a01b03165b6040516001600160a01b0390911681526020016100e6565b6001546001600160a01b031661016a565b6000546001600160a01b031661016a565b6101026101b23660046107c7565b61075f565b60006004546003546101c991906108dc565b905090565b6000546001600160a01b031633146102015760405162461bcd60e51b81526004016101f890610867565b60405180910390fd5b6040516bffffffffffffffffffffffff19606084901b166020820152603481018290526000906054016040516020818303038152906040529050600060606103fb6001600160a01b031683604051610259919061082c565b6000604051808303816000865af19150503d8060008114610296576040519150601f19603f3d011682016040523d82523d6000602084013e61029b565b606091505b5060405191935091506009906102b290839061082c565b6000604051808303816000865af19150503d80600081146102ef576040519150601f19603f3d011682016040523d82523d6000602084013e6102f4565b606091505b505050836003600082825461030991906108c4565b9091555082905061035c5760405162461bcd60e51b815260206004820152601f60248201527f546865726520697320616e206572726f72207768696c65206d696e74696e670060448201526064016101f8565b8051156103ab5760405162461bcd60e51b815260206004820152601f60248201527f546865726520697320616e206572726f72207768696c65206d696e74696e670060448201526064016101f8565b7f458f5fa412d0f69b08dd84872b0215675cc67bc1d5b6fd93300a1c3878b861963386866103d860035490565b6103e06101b7565b604080516001600160a01b039687168152959094166020860152928401919091526060830152608082015260a00160405180910390a15050505050565b6001546001600160a01b031633146104475760405162461bcd60e51b81526004016101f890610896565b600254604080516bffffffffffffffffffffffff19606093841b166020820152603480820185905282518083039091018152605490910191829052916000916103fc9061049590859061082c565b6000604051808303816000865af19150503d80600081146104d2576040519150601f19603f3d011682016040523d82523d6000602084013e6104d7565b606091505b5060405191935091506009906104ee90839061082c565b6000604051808303816000865af19150503d806000811461052b576040519150601f19603f3d011682016040523d82523d6000602084013e610530565b606091505b505050836004600082825461054591906108c4565b909155508290506105985760405162461bcd60e51b815260206004820152601f60248201527f546865726520697320616e206572726f72207768696c65206275726e696e670060448201526064016101f8565b8051156105e75760405162461bcd60e51b815260206004820152601f60248201527f546865726520697320616e206572726f72207768696c65206275726e696e670060448201526064016101f8565b6002547f4cf25bc1d991c17529c25213d3cc0cda295eeaad5f13f361969b12ea48015f909033906001600160a01b03168661062160045490565b6106296101b7565b604080516001600160a01b039687168152959094166020860152928401919091526060830152608082015260a00160405180910390a150505050565b6001546001600160a01b0316331461068f5760405162461bcd60e51b81526004016101f890610896565b600180546001600160a01b0319166001600160a01b0392909216919091179055565b6001546001600160a01b031633146106db5760405162461bcd60e51b81526004016101f890610896565b6000600455565b6001546001600160a01b0316331461070c5760405162461bcd60e51b81526004016101f890610896565b600280546001600160a01b0319166001600160a01b0392909216919091179055565b6000546001600160a01b031633146107585760405162461bcd60e51b81526004016101f890610867565b6000600355565b6000546001600160a01b031633146107895760405162461bcd60e51b81526004016101f890610867565b600080546001600160a01b0319166001600160a01b0392909216919091179055565b80356001600160a01b03811681146107c257600080fd5b919050565b6000602082840312156107d957600080fd5b6107e2826107ab565b9392505050565b600080604083850312156107fc57600080fd5b610805836107ab565b946020939093013593505050565b60006020828403121561082557600080fd5b5035919050565b6000825160005b8181101561084d5760208186018101518583015201610833565b8181111561085c576000828501525b509190910192915050565b60208082526015908201527427b7363c9036b4b73a32b91031b0b71031b0b6361760591b604082015260600190565b60208082526014908201527313db9b1e48189d5c9b995c8818d85b8818d85b1b60621b604082015260600190565b600082198211156108d7576108d76108f3565b500190565b6000828210156108ee576108ee6108f3565b500390565b634e487b7160e01b600052601160045260246000fdfea26469706673582212209b7141076237bd55d0807b09f2114bd9fe5d40bf48c78c40e767d1367cb7565364736f6c63430008070033`

// MintBurnWithPermissionFuncSigs maps the 4-byte function signature to its string representation.
var MintBurnWithPermissionFuncSigs = map[string]string{
	"42966c68": "burn(uint256)",
	"b192da2d": "burnt()",
	"d773cca9": "getBurnee()",
	"e236b8e2": "getBurner()",
	"f3667517": "getMinter()",
	"40c10f19": "mint(address,uint256)",
	"4f02c420": "minted()",
	"ace89234": "resetBurnAmount()",
	"c46d4c0f": "resetMintAmount()",
	"bf6322c3": "setBurnee(address)",
	"a996d6ce": "setBurner(address)",
	"fca3b5aa": "setMinter(address)",
	"18160ddd": "totalSupply()",
}

// MintBurnWithPermissionBin is the compiled bytecode used for deploying new contracts.
var MintBurnWithPermissionBin = "0x608060405234801561001057600080fd5b5061093f806100206000396000f3fe608060405234801561001057600080fd5b50600436106100cf5760003560e01c8063b192da2d1161008c578063d773cca911610066578063d773cca91461015d578063e236b8e214610182578063f366751714610193578063fca3b5aa146101a457600080fd5b8063b192da2d1461013a578063bf6322c314610142578063c46d4c0f1461015557600080fd5b806318160ddd146100d457806340c10f19146100ef57806342966c68146101045780634f02c42014610117578063a996d6ce1461011f578063ace8923414610132575b600080fd5b6100dc6101b7565b6040519081526020015b60405180910390f35b6101026100fd3660046107e9565b6101ce565b005b610102610112366004610813565b61041d565b6003546100dc565b61010261012d3660046107c7565b610665565b6101026106b1565b6004546100dc565b6101026101503660046107c7565b6106e2565b61010261072e565b6002546001600160a01b03165b6040516001600160a01b0390911681526020016100e6565b6001546001600160a01b031661016a565b6000546001600160a01b031661016a565b6101026101b23660046107c7565b61075f565b60006004546003546101c991906108dc565b905090565b6000546001600160a01b031633146102015760405162461bcd60e51b81526004016101f890610867565b60405180910390fd5b6040516bffffffffffffffffffffffff19606084901b166020820152603481018290526000906054016040516020818303038152906040529050600060606103fb6001600160a01b031683604051610259919061082c565b6000604051808303816000865af19150503d8060008114610296576040519150601f19603f3d011682016040523d82523d6000602084013e61029b565b606091505b5060405191935091506009906102b290839061082c565b6000604051808303816000865af19150503d80600081146102ef576040519150601f19603f3d011682016040523d82523d6000602084013e6102f4565b606091505b505050836003600082825461030991906108c4565b9091555082905061035c5760405162461bcd60e51b815260206004820152601f60248201527f546865726520697320616e206572726f72207768696c65206d696e74696e670060448201526064016101f8565b8051156103ab5760405162461bcd60e51b815260206004820152601f60248201527f546865726520697320616e206572726f72207768696c65206d696e74696e670060448201526064016101f8565b7f458f5fa412d0f69b08dd84872b0215675cc67bc1d5b6fd93300a1c3878b861963386866103d860035490565b6103e06101b7565b604080516001600160a01b039687168152959094166020860152928401919091526060830152608082015260a00160405180910390a15050505050565b6001546001600160a01b031633146104475760405162461bcd60e51b81526004016101f890610896565b600254604080516bffffffffffffffffffffffff19606093841b166020820152603480820185905282518083039091018152605490910191829052916000916103fc9061049590859061082c565b6000604051808303816000865af19150503d80600081146104d2576040519150601f19603f3d011682016040523d82523d6000602084013e6104d7565b606091505b5060405191935091506009906104ee90839061082c565b6000604051808303816000865af19150503d806000811461052b576040519150601f19603f3d011682016040523d82523d6000602084013e610530565b606091505b505050836004600082825461054591906108c4565b909155508290506105985760405162461bcd60e51b815260206004820152601f60248201527f546865726520697320616e206572726f72207768696c65206275726e696e670060448201526064016101f8565b8051156105e75760405162461bcd60e51b815260206004820152601f60248201527f546865726520697320616e206572726f72207768696c65206275726e696e670060448201526064016101f8565b6002547f4cf25bc1d991c17529c25213d3cc0cda295eeaad5f13f361969b12ea48015f909033906001600160a01b03168661062160045490565b6106296101b7565b604080516001600160a01b039687168152959094166020860152928401919091526060830152608082015260a00160405180910390a150505050565b6001546001600160a01b0316331461068f5760405162461bcd60e51b81526004016101f890610896565b600180546001600160a01b0319166001600160a01b0392909216919091179055565b6001546001600160a01b031633146106db5760405162461bcd60e51b81526004016101f890610896565b6000600455565b6001546001600160a01b0316331461070c5760405162461bcd60e51b81526004016101f890610896565b600280546001600160a01b0319166001600160a01b0392909216919091179055565b6000546001600160a01b031633146107585760405162461bcd60e51b81526004016101f890610867565b6000600355565b6000546001600160a01b031633146107895760405162461bcd60e51b81526004016101f890610867565b600080546001600160a01b0319166001600160a01b0392909216919091179055565b80356001600160a01b03811681146107c257600080fd5b919050565b6000602082840312156107d957600080fd5b6107e2826107ab565b9392505050565b600080604083850312156107fc57600080fd5b610805836107ab565b946020939093013593505050565b60006020828403121561082557600080fd5b5035919050565b6000825160005b8181101561084d5760208186018101518583015201610833565b8181111561085c576000828501525b509190910192915050565b60208082526015908201527427b7363c9036b4b73a32b91031b0b71031b0b6361760591b604082015260600190565b60208082526014908201527313db9b1e48189d5c9b995c8818d85b8818d85b1b60621b604082015260600190565b600082198211156108d7576108d76108f3565b500190565b6000828210156108ee576108ee6108f3565b500390565b634e487b7160e01b600052601160045260246000fdfea26469706673582212209b7141076237bd55d0807b09f2114bd9fe5d40bf48c78c40e767d1367cb7565364736f6c63430008070033"

// DeployMintBurnWithPermission deploys a new Klaytn contract, binding an instance of MintBurnWithPermission to it.
func DeployMintBurnWithPermission(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *MintBurnWithPermission, error) {
	parsed, err := abi.JSON(strings.NewReader(MintBurnWithPermissionABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(MintBurnWithPermissionBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &MintBurnWithPermission{MintBurnWithPermissionCaller: MintBurnWithPermissionCaller{contract: contract}, MintBurnWithPermissionTransactor: MintBurnWithPermissionTransactor{contract: contract}, MintBurnWithPermissionFilterer: MintBurnWithPermissionFilterer{contract: contract}}, nil
}

// MintBurnWithPermission is an auto generated Go binding around a Klaytn contract.
type MintBurnWithPermission struct {
	MintBurnWithPermissionCaller     // Read-only binding to the contract
	MintBurnWithPermissionTransactor // Write-only binding to the contract
	MintBurnWithPermissionFilterer   // Log filterer for contract events
}

// MintBurnWithPermissionCaller is an auto generated read-only Go binding around a Klaytn contract.
type MintBurnWithPermissionCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MintBurnWithPermissionTransactor is an auto generated write-only Go binding around a Klaytn contract.
type MintBurnWithPermissionTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MintBurnWithPermissionFilterer is an auto generated log filtering Go binding around a Klaytn contract events.
type MintBurnWithPermissionFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MintBurnWithPermissionSession is an auto generated Go binding around a Klaytn contract,
// with pre-set call and transact options.
type MintBurnWithPermissionSession struct {
	Contract     *MintBurnWithPermission // Generic contract binding to set the session for
	CallOpts     bind.CallOpts           // Call options to use throughout this session
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// MintBurnWithPermissionCallerSession is an auto generated read-only Go binding around a Klaytn contract,
// with pre-set call options.
type MintBurnWithPermissionCallerSession struct {
	Contract *MintBurnWithPermissionCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                 // Call options to use throughout this session
}

// MintBurnWithPermissionTransactorSession is an auto generated write-only Go binding around a Klaytn contract,
// with pre-set transact options.
type MintBurnWithPermissionTransactorSession struct {
	Contract     *MintBurnWithPermissionTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                 // Transaction auth options to use throughout this session
}

// MintBurnWithPermissionRaw is an auto generated low-level Go binding around a Klaytn contract.
type MintBurnWithPermissionRaw struct {
	Contract *MintBurnWithPermission // Generic contract binding to access the raw methods on
}

// MintBurnWithPermissionCallerRaw is an auto generated low-level read-only Go binding around a Klaytn contract.
type MintBurnWithPermissionCallerRaw struct {
	Contract *MintBurnWithPermissionCaller // Generic read-only contract binding to access the raw methods on
}

// MintBurnWithPermissionTransactorRaw is an auto generated low-level write-only Go binding around a Klaytn contract.
type MintBurnWithPermissionTransactorRaw struct {
	Contract *MintBurnWithPermissionTransactor // Generic write-only contract binding to access the raw methods on
}

// NewMintBurnWithPermission creates a new instance of MintBurnWithPermission, bound to a specific deployed contract.
func NewMintBurnWithPermission(address common.Address, backend bind.ContractBackend) (*MintBurnWithPermission, error) {
	contract, err := bindMintBurnWithPermission(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &MintBurnWithPermission{MintBurnWithPermissionCaller: MintBurnWithPermissionCaller{contract: contract}, MintBurnWithPermissionTransactor: MintBurnWithPermissionTransactor{contract: contract}, MintBurnWithPermissionFilterer: MintBurnWithPermissionFilterer{contract: contract}}, nil
}

// NewMintBurnWithPermissionCaller creates a new read-only instance of MintBurnWithPermission, bound to a specific deployed contract.
func NewMintBurnWithPermissionCaller(address common.Address, caller bind.ContractCaller) (*MintBurnWithPermissionCaller, error) {
	contract, err := bindMintBurnWithPermission(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MintBurnWithPermissionCaller{contract: contract}, nil
}

// NewMintBurnWithPermissionTransactor creates a new write-only instance of MintBurnWithPermission, bound to a specific deployed contract.
func NewMintBurnWithPermissionTransactor(address common.Address, transactor bind.ContractTransactor) (*MintBurnWithPermissionTransactor, error) {
	contract, err := bindMintBurnWithPermission(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MintBurnWithPermissionTransactor{contract: contract}, nil
}

// NewMintBurnWithPermissionFilterer creates a new log filterer instance of MintBurnWithPermission, bound to a specific deployed contract.
func NewMintBurnWithPermissionFilterer(address common.Address, filterer bind.ContractFilterer) (*MintBurnWithPermissionFilterer, error) {
	contract, err := bindMintBurnWithPermission(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MintBurnWithPermissionFilterer{contract: contract}, nil
}

// bindMintBurnWithPermission binds a generic wrapper to an already deployed contract.
func bindMintBurnWithPermission(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(MintBurnWithPermissionABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MintBurnWithPermission *MintBurnWithPermissionRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _MintBurnWithPermission.Contract.MintBurnWithPermissionCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MintBurnWithPermission *MintBurnWithPermissionRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MintBurnWithPermission.Contract.MintBurnWithPermissionTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MintBurnWithPermission *MintBurnWithPermissionRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MintBurnWithPermission.Contract.MintBurnWithPermissionTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MintBurnWithPermission *MintBurnWithPermissionCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _MintBurnWithPermission.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MintBurnWithPermission *MintBurnWithPermissionTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MintBurnWithPermission.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MintBurnWithPermission *MintBurnWithPermissionTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MintBurnWithPermission.Contract.contract.Transact(opts, method, params...)
}

// Burnt is a free data retrieval call binding the contract method 0xb192da2d.
//
// Solidity: function burnt() view returns(uint256)
func (_MintBurnWithPermission *MintBurnWithPermissionCaller) Burnt(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _MintBurnWithPermission.contract.Call(opts, out, "burnt")
	return *ret0, err
}

// Burnt is a free data retrieval call binding the contract method 0xb192da2d.
//
// Solidity: function burnt() view returns(uint256)
func (_MintBurnWithPermission *MintBurnWithPermissionSession) Burnt() (*big.Int, error) {
	return _MintBurnWithPermission.Contract.Burnt(&_MintBurnWithPermission.CallOpts)
}

// Burnt is a free data retrieval call binding the contract method 0xb192da2d.
//
// Solidity: function burnt() view returns(uint256)
func (_MintBurnWithPermission *MintBurnWithPermissionCallerSession) Burnt() (*big.Int, error) {
	return _MintBurnWithPermission.Contract.Burnt(&_MintBurnWithPermission.CallOpts)
}

// GetBurnee is a free data retrieval call binding the contract method 0xd773cca9.
//
// Solidity: function getBurnee() view returns(address)
func (_MintBurnWithPermission *MintBurnWithPermissionCaller) GetBurnee(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _MintBurnWithPermission.contract.Call(opts, out, "getBurnee")
	return *ret0, err
}

// GetBurnee is a free data retrieval call binding the contract method 0xd773cca9.
//
// Solidity: function getBurnee() view returns(address)
func (_MintBurnWithPermission *MintBurnWithPermissionSession) GetBurnee() (common.Address, error) {
	return _MintBurnWithPermission.Contract.GetBurnee(&_MintBurnWithPermission.CallOpts)
}

// GetBurnee is a free data retrieval call binding the contract method 0xd773cca9.
//
// Solidity: function getBurnee() view returns(address)
func (_MintBurnWithPermission *MintBurnWithPermissionCallerSession) GetBurnee() (common.Address, error) {
	return _MintBurnWithPermission.Contract.GetBurnee(&_MintBurnWithPermission.CallOpts)
}

// GetBurner is a free data retrieval call binding the contract method 0xe236b8e2.
//
// Solidity: function getBurner() view returns(address)
func (_MintBurnWithPermission *MintBurnWithPermissionCaller) GetBurner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _MintBurnWithPermission.contract.Call(opts, out, "getBurner")
	return *ret0, err
}

// GetBurner is a free data retrieval call binding the contract method 0xe236b8e2.
//
// Solidity: function getBurner() view returns(address)
func (_MintBurnWithPermission *MintBurnWithPermissionSession) GetBurner() (common.Address, error) {
	return _MintBurnWithPermission.Contract.GetBurner(&_MintBurnWithPermission.CallOpts)
}

// GetBurner is a free data retrieval call binding the contract method 0xe236b8e2.
//
// Solidity: function getBurner() view returns(address)
func (_MintBurnWithPermission *MintBurnWithPermissionCallerSession) GetBurner() (common.Address, error) {
	return _MintBurnWithPermission.Contract.GetBurner(&_MintBurnWithPermission.CallOpts)
}

// GetMinter is a free data retrieval call binding the contract method 0xf3667517.
//
// Solidity: function getMinter() view returns(address)
func (_MintBurnWithPermission *MintBurnWithPermissionCaller) GetMinter(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _MintBurnWithPermission.contract.Call(opts, out, "getMinter")
	return *ret0, err
}

// GetMinter is a free data retrieval call binding the contract method 0xf3667517.
//
// Solidity: function getMinter() view returns(address)
func (_MintBurnWithPermission *MintBurnWithPermissionSession) GetMinter() (common.Address, error) {
	return _MintBurnWithPermission.Contract.GetMinter(&_MintBurnWithPermission.CallOpts)
}

// GetMinter is a free data retrieval call binding the contract method 0xf3667517.
//
// Solidity: function getMinter() view returns(address)
func (_MintBurnWithPermission *MintBurnWithPermissionCallerSession) GetMinter() (common.Address, error) {
	return _MintBurnWithPermission.Contract.GetMinter(&_MintBurnWithPermission.CallOpts)
}

// Minted is a free data retrieval call binding the contract method 0x4f02c420.
//
// Solidity: function minted() view returns(uint256)
func (_MintBurnWithPermission *MintBurnWithPermissionCaller) Minted(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _MintBurnWithPermission.contract.Call(opts, out, "minted")
	return *ret0, err
}

// Minted is a free data retrieval call binding the contract method 0x4f02c420.
//
// Solidity: function minted() view returns(uint256)
func (_MintBurnWithPermission *MintBurnWithPermissionSession) Minted() (*big.Int, error) {
	return _MintBurnWithPermission.Contract.Minted(&_MintBurnWithPermission.CallOpts)
}

// Minted is a free data retrieval call binding the contract method 0x4f02c420.
//
// Solidity: function minted() view returns(uint256)
func (_MintBurnWithPermission *MintBurnWithPermissionCallerSession) Minted() (*big.Int, error) {
	return _MintBurnWithPermission.Contract.Minted(&_MintBurnWithPermission.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_MintBurnWithPermission *MintBurnWithPermissionCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _MintBurnWithPermission.contract.Call(opts, out, "totalSupply")
	return *ret0, err
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_MintBurnWithPermission *MintBurnWithPermissionSession) TotalSupply() (*big.Int, error) {
	return _MintBurnWithPermission.Contract.TotalSupply(&_MintBurnWithPermission.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_MintBurnWithPermission *MintBurnWithPermissionCallerSession) TotalSupply() (*big.Int, error) {
	return _MintBurnWithPermission.Contract.TotalSupply(&_MintBurnWithPermission.CallOpts)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 amount) returns()
func (_MintBurnWithPermission *MintBurnWithPermissionTransactor) Burn(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _MintBurnWithPermission.contract.Transact(opts, "burn", amount)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 amount) returns()
func (_MintBurnWithPermission *MintBurnWithPermissionSession) Burn(amount *big.Int) (*types.Transaction, error) {
	return _MintBurnWithPermission.Contract.Burn(&_MintBurnWithPermission.TransactOpts, amount)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 amount) returns()
func (_MintBurnWithPermission *MintBurnWithPermissionTransactorSession) Burn(amount *big.Int) (*types.Transaction, error) {
	return _MintBurnWithPermission.Contract.Burn(&_MintBurnWithPermission.TransactOpts, amount)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address account, uint256 amount) returns()
func (_MintBurnWithPermission *MintBurnWithPermissionTransactor) Mint(opts *bind.TransactOpts, account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _MintBurnWithPermission.contract.Transact(opts, "mint", account, amount)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address account, uint256 amount) returns()
func (_MintBurnWithPermission *MintBurnWithPermissionSession) Mint(account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _MintBurnWithPermission.Contract.Mint(&_MintBurnWithPermission.TransactOpts, account, amount)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address account, uint256 amount) returns()
func (_MintBurnWithPermission *MintBurnWithPermissionTransactorSession) Mint(account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _MintBurnWithPermission.Contract.Mint(&_MintBurnWithPermission.TransactOpts, account, amount)
}

// ResetBurnAmount is a paid mutator transaction binding the contract method 0xace89234.
//
// Solidity: function resetBurnAmount() returns()
func (_MintBurnWithPermission *MintBurnWithPermissionTransactor) ResetBurnAmount(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MintBurnWithPermission.contract.Transact(opts, "resetBurnAmount")
}

// ResetBurnAmount is a paid mutator transaction binding the contract method 0xace89234.
//
// Solidity: function resetBurnAmount() returns()
func (_MintBurnWithPermission *MintBurnWithPermissionSession) ResetBurnAmount() (*types.Transaction, error) {
	return _MintBurnWithPermission.Contract.ResetBurnAmount(&_MintBurnWithPermission.TransactOpts)
}

// ResetBurnAmount is a paid mutator transaction binding the contract method 0xace89234.
//
// Solidity: function resetBurnAmount() returns()
func (_MintBurnWithPermission *MintBurnWithPermissionTransactorSession) ResetBurnAmount() (*types.Transaction, error) {
	return _MintBurnWithPermission.Contract.ResetBurnAmount(&_MintBurnWithPermission.TransactOpts)
}

// ResetMintAmount is a paid mutator transaction binding the contract method 0xc46d4c0f.
//
// Solidity: function resetMintAmount() returns()
func (_MintBurnWithPermission *MintBurnWithPermissionTransactor) ResetMintAmount(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MintBurnWithPermission.contract.Transact(opts, "resetMintAmount")
}

// ResetMintAmount is a paid mutator transaction binding the contract method 0xc46d4c0f.
//
// Solidity: function resetMintAmount() returns()
func (_MintBurnWithPermission *MintBurnWithPermissionSession) ResetMintAmount() (*types.Transaction, error) {
	return _MintBurnWithPermission.Contract.ResetMintAmount(&_MintBurnWithPermission.TransactOpts)
}

// ResetMintAmount is a paid mutator transaction binding the contract method 0xc46d4c0f.
//
// Solidity: function resetMintAmount() returns()
func (_MintBurnWithPermission *MintBurnWithPermissionTransactorSession) ResetMintAmount() (*types.Transaction, error) {
	return _MintBurnWithPermission.Contract.ResetMintAmount(&_MintBurnWithPermission.TransactOpts)
}

// SetBurnee is a paid mutator transaction binding the contract method 0xbf6322c3.
//
// Solidity: function setBurnee(address account) returns()
func (_MintBurnWithPermission *MintBurnWithPermissionTransactor) SetBurnee(opts *bind.TransactOpts, account common.Address) (*types.Transaction, error) {
	return _MintBurnWithPermission.contract.Transact(opts, "setBurnee", account)
}

// SetBurnee is a paid mutator transaction binding the contract method 0xbf6322c3.
//
// Solidity: function setBurnee(address account) returns()
func (_MintBurnWithPermission *MintBurnWithPermissionSession) SetBurnee(account common.Address) (*types.Transaction, error) {
	return _MintBurnWithPermission.Contract.SetBurnee(&_MintBurnWithPermission.TransactOpts, account)
}

// SetBurnee is a paid mutator transaction binding the contract method 0xbf6322c3.
//
// Solidity: function setBurnee(address account) returns()
func (_MintBurnWithPermission *MintBurnWithPermissionTransactorSession) SetBurnee(account common.Address) (*types.Transaction, error) {
	return _MintBurnWithPermission.Contract.SetBurnee(&_MintBurnWithPermission.TransactOpts, account)
}

// SetBurner is a paid mutator transaction binding the contract method 0xa996d6ce.
//
// Solidity: function setBurner(address account) returns()
func (_MintBurnWithPermission *MintBurnWithPermissionTransactor) SetBurner(opts *bind.TransactOpts, account common.Address) (*types.Transaction, error) {
	return _MintBurnWithPermission.contract.Transact(opts, "setBurner", account)
}

// SetBurner is a paid mutator transaction binding the contract method 0xa996d6ce.
//
// Solidity: function setBurner(address account) returns()
func (_MintBurnWithPermission *MintBurnWithPermissionSession) SetBurner(account common.Address) (*types.Transaction, error) {
	return _MintBurnWithPermission.Contract.SetBurner(&_MintBurnWithPermission.TransactOpts, account)
}

// SetBurner is a paid mutator transaction binding the contract method 0xa996d6ce.
//
// Solidity: function setBurner(address account) returns()
func (_MintBurnWithPermission *MintBurnWithPermissionTransactorSession) SetBurner(account common.Address) (*types.Transaction, error) {
	return _MintBurnWithPermission.Contract.SetBurner(&_MintBurnWithPermission.TransactOpts, account)
}

// SetMinter is a paid mutator transaction binding the contract method 0xfca3b5aa.
//
// Solidity: function setMinter(address account) returns()
func (_MintBurnWithPermission *MintBurnWithPermissionTransactor) SetMinter(opts *bind.TransactOpts, account common.Address) (*types.Transaction, error) {
	return _MintBurnWithPermission.contract.Transact(opts, "setMinter", account)
}

// SetMinter is a paid mutator transaction binding the contract method 0xfca3b5aa.
//
// Solidity: function setMinter(address account) returns()
func (_MintBurnWithPermission *MintBurnWithPermissionSession) SetMinter(account common.Address) (*types.Transaction, error) {
	return _MintBurnWithPermission.Contract.SetMinter(&_MintBurnWithPermission.TransactOpts, account)
}

// SetMinter is a paid mutator transaction binding the contract method 0xfca3b5aa.
//
// Solidity: function setMinter(address account) returns()
func (_MintBurnWithPermission *MintBurnWithPermissionTransactorSession) SetMinter(account common.Address) (*types.Transaction, error) {
	return _MintBurnWithPermission.Contract.SetMinter(&_MintBurnWithPermission.TransactOpts, account)
}

// MintBurnWithPermissionBurnIterator is returned from FilterBurn and is used to iterate over the raw logs and unpacked data for Burn events raised by the MintBurnWithPermission contract.
type MintBurnWithPermissionBurnIterator struct {
	Event *MintBurnWithPermissionBurn // Event containing the contract specifics and raw log

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
func (it *MintBurnWithPermissionBurnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MintBurnWithPermissionBurn)
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
		it.Event = new(MintBurnWithPermissionBurn)
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
func (it *MintBurnWithPermissionBurnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MintBurnWithPermissionBurnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MintBurnWithPermissionBurn represents a Burn event raised by the MintBurnWithPermission contract.
type MintBurnWithPermissionBurn struct {
	Burner        common.Address
	Burnee        common.Address
	BurnAmount    *big.Int
	AccBurnAmount *big.Int
	TotalSupply   *big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterBurn is a free log retrieval operation binding the contract event 0x4cf25bc1d991c17529c25213d3cc0cda295eeaad5f13f361969b12ea48015f90.
//
// Solidity: event Burn(address burner, address burnee, uint256 burnAmount, uint256 accBurnAmount, uint256 totalSupply)
func (_MintBurnWithPermission *MintBurnWithPermissionFilterer) FilterBurn(opts *bind.FilterOpts) (*MintBurnWithPermissionBurnIterator, error) {

	logs, sub, err := _MintBurnWithPermission.contract.FilterLogs(opts, "Burn")
	if err != nil {
		return nil, err
	}
	return &MintBurnWithPermissionBurnIterator{contract: _MintBurnWithPermission.contract, event: "Burn", logs: logs, sub: sub}, nil
}

// WatchBurn is a free log subscription operation binding the contract event 0x4cf25bc1d991c17529c25213d3cc0cda295eeaad5f13f361969b12ea48015f90.
//
// Solidity: event Burn(address burner, address burnee, uint256 burnAmount, uint256 accBurnAmount, uint256 totalSupply)
func (_MintBurnWithPermission *MintBurnWithPermissionFilterer) WatchBurn(opts *bind.WatchOpts, sink chan<- *MintBurnWithPermissionBurn) (event.Subscription, error) {

	logs, sub, err := _MintBurnWithPermission.contract.WatchLogs(opts, "Burn")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MintBurnWithPermissionBurn)
				if err := _MintBurnWithPermission.contract.UnpackLog(event, "Burn", log); err != nil {
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

// ParseBurn is a log parse operation binding the contract event 0x4cf25bc1d991c17529c25213d3cc0cda295eeaad5f13f361969b12ea48015f90.
//
// Solidity: event Burn(address burner, address burnee, uint256 burnAmount, uint256 accBurnAmount, uint256 totalSupply)
func (_MintBurnWithPermission *MintBurnWithPermissionFilterer) ParseBurn(log types.Log) (*MintBurnWithPermissionBurn, error) {
	event := new(MintBurnWithPermissionBurn)
	if err := _MintBurnWithPermission.contract.UnpackLog(event, "Burn", log); err != nil {
		return nil, err
	}
	return event, nil
}

// MintBurnWithPermissionMintIterator is returned from FilterMint and is used to iterate over the raw logs and unpacked data for Mint events raised by the MintBurnWithPermission contract.
type MintBurnWithPermissionMintIterator struct {
	Event *MintBurnWithPermissionMint // Event containing the contract specifics and raw log

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
func (it *MintBurnWithPermissionMintIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MintBurnWithPermissionMint)
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
		it.Event = new(MintBurnWithPermissionMint)
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
func (it *MintBurnWithPermissionMintIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MintBurnWithPermissionMintIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MintBurnWithPermissionMint represents a Mint event raised by the MintBurnWithPermission contract.
type MintBurnWithPermissionMint struct {
	Minter        common.Address
	Mintee        common.Address
	MintAmount    *big.Int
	AccMintAmount *big.Int
	TotalSupply   *big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterMint is a free log retrieval operation binding the contract event 0x458f5fa412d0f69b08dd84872b0215675cc67bc1d5b6fd93300a1c3878b86196.
//
// Solidity: event Mint(address minter, address mintee, uint256 mintAmount, uint256 accMintAmount, uint256 totalSupply)
func (_MintBurnWithPermission *MintBurnWithPermissionFilterer) FilterMint(opts *bind.FilterOpts) (*MintBurnWithPermissionMintIterator, error) {

	logs, sub, err := _MintBurnWithPermission.contract.FilterLogs(opts, "Mint")
	if err != nil {
		return nil, err
	}
	return &MintBurnWithPermissionMintIterator{contract: _MintBurnWithPermission.contract, event: "Mint", logs: logs, sub: sub}, nil
}

// WatchMint is a free log subscription operation binding the contract event 0x458f5fa412d0f69b08dd84872b0215675cc67bc1d5b6fd93300a1c3878b86196.
//
// Solidity: event Mint(address minter, address mintee, uint256 mintAmount, uint256 accMintAmount, uint256 totalSupply)
func (_MintBurnWithPermission *MintBurnWithPermissionFilterer) WatchMint(opts *bind.WatchOpts, sink chan<- *MintBurnWithPermissionMint) (event.Subscription, error) {

	logs, sub, err := _MintBurnWithPermission.contract.WatchLogs(opts, "Mint")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MintBurnWithPermissionMint)
				if err := _MintBurnWithPermission.contract.UnpackLog(event, "Mint", log); err != nil {
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

// ParseMint is a log parse operation binding the contract event 0x458f5fa412d0f69b08dd84872b0215675cc67bc1d5b6fd93300a1c3878b86196.
//
// Solidity: event Mint(address minter, address mintee, uint256 mintAmount, uint256 accMintAmount, uint256 totalSupply)
func (_MintBurnWithPermission *MintBurnWithPermissionFilterer) ParseMint(log types.Log) (*MintBurnWithPermissionMint, error) {
	event := new(MintBurnWithPermissionMint)
	if err := _MintBurnWithPermission.contract.UnpackLog(event, "Mint", log); err != nil {
		return nil, err
	}
	return event, nil
}
