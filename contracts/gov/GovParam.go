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

// GovParamParamView is an auto generated low-level Go binding around an user-defined struct.
type GovParamParamView struct {
	Name  string
	Value []byte
}

// ContextABI is the input ABI used to generate the binding from.
const ContextABI = "[]"

// ContextBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const ContextBinRuntime = ``

// Context is an auto generated Go binding around a Klaytn contract.
type Context struct {
	ContextCaller     // Read-only binding to the contract
	ContextTransactor // Write-only binding to the contract
	ContextFilterer   // Log filterer for contract events
}

// ContextCaller is an auto generated read-only Go binding around a Klaytn contract.
type ContextCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContextTransactor is an auto generated write-only Go binding around a Klaytn contract.
type ContextTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContextFilterer is an auto generated log filtering Go binding around a Klaytn contract events.
type ContextFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContextSession is an auto generated Go binding around a Klaytn contract,
// with pre-set call and transact options.
type ContextSession struct {
	Contract     *Context          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ContextCallerSession is an auto generated read-only Go binding around a Klaytn contract,
// with pre-set call options.
type ContextCallerSession struct {
	Contract *ContextCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// ContextTransactorSession is an auto generated write-only Go binding around a Klaytn contract,
// with pre-set transact options.
type ContextTransactorSession struct {
	Contract     *ContextTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// ContextRaw is an auto generated low-level Go binding around a Klaytn contract.
type ContextRaw struct {
	Contract *Context // Generic contract binding to access the raw methods on
}

// ContextCallerRaw is an auto generated low-level read-only Go binding around a Klaytn contract.
type ContextCallerRaw struct {
	Contract *ContextCaller // Generic read-only contract binding to access the raw methods on
}

// ContextTransactorRaw is an auto generated low-level write-only Go binding around a Klaytn contract.
type ContextTransactorRaw struct {
	Contract *ContextTransactor // Generic write-only contract binding to access the raw methods on
}

// NewContext creates a new instance of Context, bound to a specific deployed contract.
func NewContext(address common.Address, backend bind.ContractBackend) (*Context, error) {
	contract, err := bindContext(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Context{ContextCaller: ContextCaller{contract: contract}, ContextTransactor: ContextTransactor{contract: contract}, ContextFilterer: ContextFilterer{contract: contract}}, nil
}

// NewContextCaller creates a new read-only instance of Context, bound to a specific deployed contract.
func NewContextCaller(address common.Address, caller bind.ContractCaller) (*ContextCaller, error) {
	contract, err := bindContext(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContextCaller{contract: contract}, nil
}

// NewContextTransactor creates a new write-only instance of Context, bound to a specific deployed contract.
func NewContextTransactor(address common.Address, transactor bind.ContractTransactor) (*ContextTransactor, error) {
	contract, err := bindContext(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContextTransactor{contract: contract}, nil
}

// NewContextFilterer creates a new log filterer instance of Context, bound to a specific deployed contract.
func NewContextFilterer(address common.Address, filterer bind.ContractFilterer) (*ContextFilterer, error) {
	contract, err := bindContext(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContextFilterer{contract: contract}, nil
}

// bindContext binds a generic wrapper to an already deployed contract.
func bindContext(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ContextABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Context *ContextRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Context.Contract.ContextCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Context *ContextRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Context.Contract.ContextTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Context *ContextRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Context.Contract.ContextTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Context *ContextCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Context.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Context *ContextTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Context.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Context *ContextTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Context.Contract.contract.Transact(opts, method, params...)
}

// GovParamABI is the input ABI used to generate the binding from.
const GovParamABI = "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"name\":\"SetParam\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"getAllParams\",\"outputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"value\",\"type\":\"bytes\"}],\"internalType\":\"structGovParam.ParamView[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllParamsAtNextBlock\",\"outputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"value\",\"type\":\"bytes\"}],\"internalType\":\"structGovParam.ParamView[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"name\":\"getParam\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"name\":\"getParamAtNextBlock\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"value\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"activation\",\"type\":\"uint64\"}],\"name\":\"setParam\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// GovParamBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const GovParamBinRuntime = `608060405234801561001057600080fd5b50600436106100885760003560e01c80639c9c03bb1161005b5780639c9c03bb146100ee578063a170052e14610101578063def2466814610116578063f2fde38b1461011e57600080fd5b80635d4f71d41461008d578063715018a6146100b65780638841067a146100c05780638da5cb5b146100d3575b600080fd5b6100a061009b366004610af7565b610131565b6040516100ad9190610c04565b60405180910390f35b6100be610143565b005b6100be6100ce366004610c67565b610182565b6000546040516001600160a01b0390911681526020016100ad565b6100a06100fc366004610af7565b610483565b610109610499565b6040516100ad9190610cf1565b6101096104a9565b6100be61012c366004610d76565b6104be565b606061013d8243610559565b92915050565b6000546001600160a01b031633146101765760405162461bcd60e51b815260040161016d90610d9f565b60405180910390fd5b6101806000610756565b565b6000546001600160a01b031633146101ac5760405162461bcd60e51b815260040161016d90610d9f565b836101f05760405162461bcd60e51b81526020600482015260146024820152736e616d652063616e6e6f7420626520656d70747960601b604482015260640161016d565b4360028686604051610203929190610dd4565b9081526040519081900360200190205467ffffffffffffffff61010090910416106102705760405162461bcd60e51b815260206004820152601d60248201527f616c7265616479206861766520612070656e64696e67206368616e6765000000604482015260640161016d565b438167ffffffffffffffff16116102c95760405162461bcd60e51b815260206004820152601e60248201527f61637469766174696f6e206d75737420626520696e2061206675747572650000604482015260640161016d565b600285856040516102db929190610dd4565b9081526040519081900360200190205460ff1661036457600160028686604051610306929190610dd4565b908152604051908190036020019020805491151560ff19909216919091179055600180548082018255600091909152610362907fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf60186866109cd565b505b8060028686604051610377929190610dd4565b908152604051908190036020018120805467ffffffffffffffff939093166101000268ffffffffffffffff0019909316929092179091556002906103be9087908790610dd4565b9081526020016040518091039020600201600286866040516103e1929190610dd4565b90815260200160405180910390206001019080546103fe90610de4565b610409929190610a51565b5082826002878760405161041e929190610dd4565b9081526020016040518091039020600201919061043c9291906109cd565b507ff8d4c4710beeb9d48e0bc6888dfcaf21467794f5da39c3bfb2874a2ac3a73c568585858585604051610474959493929190610e48565b60405180910390a15050505050565b606061013d82610494436001610ea2565b610559565b60606104a4436107a6565b905090565b60606104a46104b9436001610ea2565b6107a6565b6000546001600160a01b031633146104e85760405162461bcd60e51b815260040161016d90610d9f565b6001600160a01b03811661054d5760405162461bcd60e51b815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201526564647265737360d01b606482015260840161016d565b61055681610756565b50565b6060438210156105a35760405162461bcd60e51b815260206004820152601560248201527418d85b9b9bdd081c5d595c9e481d1a19481c185cdd605a1b604482015260640161016d565b6105ae436001610ea2565b8211156106075760405162461bcd60e51b815260206004820152602160248201527f63616e6e6f74207175657279207468652074656e7461746976652066757475726044820152606560f81b606482015260840161016d565b6002836040516106179190610eba565b9081526040519081900360200190205460ff16610643575060408051602081019091526000815261013d565b6002836040516106539190610eba565b9081526040519081900360200190205467ffffffffffffffff61010090910416821061072a576002836040516106899190610eba565b908152602001604051809103902060020180546106a590610de4565b80601f01602080910402602001604051908101604052809291908181526020018280546106d190610de4565b801561071e5780601f106106f35761010080835404028352916020019161071e565b820191906000526020600020905b81548152906001019060200180831161070157829003601f168201915b5050505050905061013d565b60028360405161073a9190610eba565b908152602001604051809103902060010180546106a590610de4565b600080546001600160a01b038381166001600160a01b0319831681178455604051919092169283917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e09190a35050565b60015460609060009067ffffffffffffffff8111156107c7576107c7610ae1565b60405190808252806020026020018201604052801561080c57816020015b60408051808201909152606080825260208201528152602001906001900390816107e55790505b50905060005b6001548110156109c6576001818154811061082f5761082f610ed6565b90600052602060002001805461084490610de4565b80601f016020809104026020016040519081016040528092919081815260200182805461087090610de4565b80156108bd5780601f10610892576101008083540402835291602001916108bd565b820191906000526020600020905b8154815290600101906020018083116108a057829003601f168201915b50505050508282815181106108d4576108d4610ed6565b602002602001015160000181905250610992600182815481106108f9576108f9610ed6565b90600052602060002001805461090e90610de4565b80601f016020809104026020016040519081016040528092919081815260200182805461093a90610de4565b80156109875780601f1061095c57610100808354040283529160200191610987565b820191906000526020600020905b81548152906001019060200180831161096a57829003601f168201915b505050505085610559565b8282815181106109a4576109a4610ed6565b60200260200101516020018190525080806109be90610eec565b915050610812565b5092915050565b8280546109d990610de4565b90600052602060002090601f0160209004810192826109fb5760008555610a41565b82601f10610a145782800160ff19823516178555610a41565b82800160010185558215610a41579182015b82811115610a41578235825591602001919060010190610a26565b50610a4d929150610acc565b5090565b828054610a5d90610de4565b90600052602060002090601f016020900481019282610a7f5760008555610a41565b82601f10610a905780548555610a41565b82800160010185558215610a4157600052602060002091601f016020900482015b82811115610a41578254825591600101919060010190610ab1565b5b80821115610a4d5760008155600101610acd565b634e487b7160e01b600052604160045260246000fd5b600060208284031215610b0957600080fd5b813567ffffffffffffffff80821115610b2157600080fd5b818401915084601f830112610b3557600080fd5b813581811115610b4757610b47610ae1565b604051601f8201601f19908116603f01168101908382118183101715610b6f57610b6f610ae1565b81604052828152876020848701011115610b8857600080fd5b826020860160208301376000928101602001929092525095945050505050565b60005b83811015610bc3578181015183820152602001610bab565b83811115610bd2576000848401525b50505050565b60008151808452610bf0816020860160208601610ba8565b601f01601f19169290920160200192915050565b602081526000610c176020830184610bd8565b9392505050565b60008083601f840112610c3057600080fd5b50813567ffffffffffffffff811115610c4857600080fd5b602083019150836020828501011115610c6057600080fd5b9250929050565b600080600080600060608688031215610c7f57600080fd5b853567ffffffffffffffff80821115610c9757600080fd5b610ca389838a01610c1e565b90975095506020880135915080821115610cbc57600080fd5b610cc889838a01610c1e565b9095509350604088013591508082168214610ce257600080fd5b50809150509295509295909350565b60006020808301818452808551808352604092508286019150828160051b87010184880160005b83811015610d6857888303603f1901855281518051878552610d3c88860182610bd8565b91890151858303868b0152919050610d548183610bd8565b968901969450505090860190600101610d18565b509098975050505050505050565b600060208284031215610d8857600080fd5b81356001600160a01b0381168114610c1757600080fd5b6020808252818101527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572604082015260600190565b8183823760009101908152919050565b600181811c90821680610df857607f821691505b60208210811415610e1957634e487b7160e01b600052602260045260246000fd5b50919050565b81835281816020850137506000828201602090810191909152601f909101601f19169091010190565b606081526000610e5c606083018789610e1f565b8281036020840152610e6f818688610e1f565b91505067ffffffffffffffff831660408301529695505050505050565b634e487b7160e01b600052601160045260246000fd5b60008219821115610eb557610eb5610e8c565b500190565b60008251610ecc818460208701610ba8565b9190910192915050565b634e487b7160e01b600052603260045260246000fd5b6000600019821415610f0057610f00610e8c565b506001019056fea2646970667358221220df4c0a1c90ea88b7949359549af412d98933a10a15a2bdfd01ad5eb05457002f64736f6c634300080b0033`

// GovParamFuncSigs maps the 4-byte function signature to its string representation.
var GovParamFuncSigs = map[string]string{
	"a170052e": "getAllParams()",
	"def24668": "getAllParamsAtNextBlock()",
	"5d4f71d4": "getParam(string)",
	"9c9c03bb": "getParamAtNextBlock(string)",
	"8da5cb5b": "owner()",
	"715018a6": "renounceOwnership()",
	"8841067a": "setParam(string,bytes,uint64)",
	"f2fde38b": "transferOwnership(address)",
}

// GovParamBin is the compiled bytecode used for deploying new contracts.
var GovParamBin = "0x608060405234801561001057600080fd5b5060405161102238038061102283398101604081905261002f916100a6565b61003833610056565b6001600160a01b038116156100505761005081610056565b506100d6565b600080546001600160a01b038381166001600160a01b0319831681178455604051919092169283917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e09190a35050565b6000602082840312156100b857600080fd5b81516001600160a01b03811681146100cf57600080fd5b9392505050565b610f3d806100e56000396000f3fe608060405234801561001057600080fd5b50600436106100885760003560e01c80639c9c03bb1161005b5780639c9c03bb146100ee578063a170052e14610101578063def2466814610116578063f2fde38b1461011e57600080fd5b80635d4f71d41461008d578063715018a6146100b65780638841067a146100c05780638da5cb5b146100d3575b600080fd5b6100a061009b366004610af7565b610131565b6040516100ad9190610c04565b60405180910390f35b6100be610143565b005b6100be6100ce366004610c67565b610182565b6000546040516001600160a01b0390911681526020016100ad565b6100a06100fc366004610af7565b610483565b610109610499565b6040516100ad9190610cf1565b6101096104a9565b6100be61012c366004610d76565b6104be565b606061013d8243610559565b92915050565b6000546001600160a01b031633146101765760405162461bcd60e51b815260040161016d90610d9f565b60405180910390fd5b6101806000610756565b565b6000546001600160a01b031633146101ac5760405162461bcd60e51b815260040161016d90610d9f565b836101f05760405162461bcd60e51b81526020600482015260146024820152736e616d652063616e6e6f7420626520656d70747960601b604482015260640161016d565b4360028686604051610203929190610dd4565b9081526040519081900360200190205467ffffffffffffffff61010090910416106102705760405162461bcd60e51b815260206004820152601d60248201527f616c7265616479206861766520612070656e64696e67206368616e6765000000604482015260640161016d565b438167ffffffffffffffff16116102c95760405162461bcd60e51b815260206004820152601e60248201527f61637469766174696f6e206d75737420626520696e2061206675747572650000604482015260640161016d565b600285856040516102db929190610dd4565b9081526040519081900360200190205460ff1661036457600160028686604051610306929190610dd4565b908152604051908190036020019020805491151560ff19909216919091179055600180548082018255600091909152610362907fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf60186866109cd565b505b8060028686604051610377929190610dd4565b908152604051908190036020018120805467ffffffffffffffff939093166101000268ffffffffffffffff0019909316929092179091556002906103be9087908790610dd4565b9081526020016040518091039020600201600286866040516103e1929190610dd4565b90815260200160405180910390206001019080546103fe90610de4565b610409929190610a51565b5082826002878760405161041e929190610dd4565b9081526020016040518091039020600201919061043c9291906109cd565b507ff8d4c4710beeb9d48e0bc6888dfcaf21467794f5da39c3bfb2874a2ac3a73c568585858585604051610474959493929190610e48565b60405180910390a15050505050565b606061013d82610494436001610ea2565b610559565b60606104a4436107a6565b905090565b60606104a46104b9436001610ea2565b6107a6565b6000546001600160a01b031633146104e85760405162461bcd60e51b815260040161016d90610d9f565b6001600160a01b03811661054d5760405162461bcd60e51b815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201526564647265737360d01b606482015260840161016d565b61055681610756565b50565b6060438210156105a35760405162461bcd60e51b815260206004820152601560248201527418d85b9b9bdd081c5d595c9e481d1a19481c185cdd605a1b604482015260640161016d565b6105ae436001610ea2565b8211156106075760405162461bcd60e51b815260206004820152602160248201527f63616e6e6f74207175657279207468652074656e7461746976652066757475726044820152606560f81b606482015260840161016d565b6002836040516106179190610eba565b9081526040519081900360200190205460ff16610643575060408051602081019091526000815261013d565b6002836040516106539190610eba565b9081526040519081900360200190205467ffffffffffffffff61010090910416821061072a576002836040516106899190610eba565b908152602001604051809103902060020180546106a590610de4565b80601f01602080910402602001604051908101604052809291908181526020018280546106d190610de4565b801561071e5780601f106106f35761010080835404028352916020019161071e565b820191906000526020600020905b81548152906001019060200180831161070157829003601f168201915b5050505050905061013d565b60028360405161073a9190610eba565b908152602001604051809103902060010180546106a590610de4565b600080546001600160a01b038381166001600160a01b0319831681178455604051919092169283917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e09190a35050565b60015460609060009067ffffffffffffffff8111156107c7576107c7610ae1565b60405190808252806020026020018201604052801561080c57816020015b60408051808201909152606080825260208201528152602001906001900390816107e55790505b50905060005b6001548110156109c6576001818154811061082f5761082f610ed6565b90600052602060002001805461084490610de4565b80601f016020809104026020016040519081016040528092919081815260200182805461087090610de4565b80156108bd5780601f10610892576101008083540402835291602001916108bd565b820191906000526020600020905b8154815290600101906020018083116108a057829003601f168201915b50505050508282815181106108d4576108d4610ed6565b602002602001015160000181905250610992600182815481106108f9576108f9610ed6565b90600052602060002001805461090e90610de4565b80601f016020809104026020016040519081016040528092919081815260200182805461093a90610de4565b80156109875780601f1061095c57610100808354040283529160200191610987565b820191906000526020600020905b81548152906001019060200180831161096a57829003601f168201915b505050505085610559565b8282815181106109a4576109a4610ed6565b60200260200101516020018190525080806109be90610eec565b915050610812565b5092915050565b8280546109d990610de4565b90600052602060002090601f0160209004810192826109fb5760008555610a41565b82601f10610a145782800160ff19823516178555610a41565b82800160010185558215610a41579182015b82811115610a41578235825591602001919060010190610a26565b50610a4d929150610acc565b5090565b828054610a5d90610de4565b90600052602060002090601f016020900481019282610a7f5760008555610a41565b82601f10610a905780548555610a41565b82800160010185558215610a4157600052602060002091601f016020900482015b82811115610a41578254825591600101919060010190610ab1565b5b80821115610a4d5760008155600101610acd565b634e487b7160e01b600052604160045260246000fd5b600060208284031215610b0957600080fd5b813567ffffffffffffffff80821115610b2157600080fd5b818401915084601f830112610b3557600080fd5b813581811115610b4757610b47610ae1565b604051601f8201601f19908116603f01168101908382118183101715610b6f57610b6f610ae1565b81604052828152876020848701011115610b8857600080fd5b826020860160208301376000928101602001929092525095945050505050565b60005b83811015610bc3578181015183820152602001610bab565b83811115610bd2576000848401525b50505050565b60008151808452610bf0816020860160208601610ba8565b601f01601f19169290920160200192915050565b602081526000610c176020830184610bd8565b9392505050565b60008083601f840112610c3057600080fd5b50813567ffffffffffffffff811115610c4857600080fd5b602083019150836020828501011115610c6057600080fd5b9250929050565b600080600080600060608688031215610c7f57600080fd5b853567ffffffffffffffff80821115610c9757600080fd5b610ca389838a01610c1e565b90975095506020880135915080821115610cbc57600080fd5b610cc889838a01610c1e565b9095509350604088013591508082168214610ce257600080fd5b50809150509295509295909350565b60006020808301818452808551808352604092508286019150828160051b87010184880160005b83811015610d6857888303603f1901855281518051878552610d3c88860182610bd8565b91890151858303868b0152919050610d548183610bd8565b968901969450505090860190600101610d18565b509098975050505050505050565b600060208284031215610d8857600080fd5b81356001600160a01b0381168114610c1757600080fd5b6020808252818101527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572604082015260600190565b8183823760009101908152919050565b600181811c90821680610df857607f821691505b60208210811415610e1957634e487b7160e01b600052602260045260246000fd5b50919050565b81835281816020850137506000828201602090810191909152601f909101601f19169091010190565b606081526000610e5c606083018789610e1f565b8281036020840152610e6f818688610e1f565b91505067ffffffffffffffff831660408301529695505050505050565b634e487b7160e01b600052601160045260246000fd5b60008219821115610eb557610eb5610e8c565b500190565b60008251610ecc818460208701610ba8565b9190910192915050565b634e487b7160e01b600052603260045260246000fd5b6000600019821415610f0057610f00610e8c565b506001019056fea2646970667358221220df4c0a1c90ea88b7949359549af412d98933a10a15a2bdfd01ad5eb05457002f64736f6c634300080b0033"

// DeployGovParam deploys a new Klaytn contract, binding an instance of GovParam to it.
func DeployGovParam(auth *bind.TransactOpts, backend bind.ContractBackend, _owner common.Address) (common.Address, *types.Transaction, *GovParam, error) {
	parsed, err := abi.JSON(strings.NewReader(GovParamABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(GovParamBin), backend, _owner)
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
// Solidity: function getAllParams() view returns((string,bytes)[])
func (_GovParam *GovParamCaller) GetAllParams(opts *bind.CallOpts) ([]GovParamParamView, error) {
	var (
		ret0 = new([]GovParamParamView)
	)
	out := ret0
	err := _GovParam.contract.Call(opts, out, "getAllParams")
	return *ret0, err
}

// GetAllParams is a free data retrieval call binding the contract method 0xa170052e.
//
// Solidity: function getAllParams() view returns((string,bytes)[])
func (_GovParam *GovParamSession) GetAllParams() ([]GovParamParamView, error) {
	return _GovParam.Contract.GetAllParams(&_GovParam.CallOpts)
}

// GetAllParams is a free data retrieval call binding the contract method 0xa170052e.
//
// Solidity: function getAllParams() view returns((string,bytes)[])
func (_GovParam *GovParamCallerSession) GetAllParams() ([]GovParamParamView, error) {
	return _GovParam.Contract.GetAllParams(&_GovParam.CallOpts)
}

// GetAllParamsAtNextBlock is a free data retrieval call binding the contract method 0xdef24668.
//
// Solidity: function getAllParamsAtNextBlock() view returns((string,bytes)[])
func (_GovParam *GovParamCaller) GetAllParamsAtNextBlock(opts *bind.CallOpts) ([]GovParamParamView, error) {
	var (
		ret0 = new([]GovParamParamView)
	)
	out := ret0
	err := _GovParam.contract.Call(opts, out, "getAllParamsAtNextBlock")
	return *ret0, err
}

// GetAllParamsAtNextBlock is a free data retrieval call binding the contract method 0xdef24668.
//
// Solidity: function getAllParamsAtNextBlock() view returns((string,bytes)[])
func (_GovParam *GovParamSession) GetAllParamsAtNextBlock() ([]GovParamParamView, error) {
	return _GovParam.Contract.GetAllParamsAtNextBlock(&_GovParam.CallOpts)
}

// GetAllParamsAtNextBlock is a free data retrieval call binding the contract method 0xdef24668.
//
// Solidity: function getAllParamsAtNextBlock() view returns((string,bytes)[])
func (_GovParam *GovParamCallerSession) GetAllParamsAtNextBlock() ([]GovParamParamView, error) {
	return _GovParam.Contract.GetAllParamsAtNextBlock(&_GovParam.CallOpts)
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

// GetParamAtNextBlock is a free data retrieval call binding the contract method 0x9c9c03bb.
//
// Solidity: function getParamAtNextBlock(string name) view returns(bytes)
func (_GovParam *GovParamCaller) GetParamAtNextBlock(opts *bind.CallOpts, name string) ([]byte, error) {
	var (
		ret0 = new([]byte)
	)
	out := ret0
	err := _GovParam.contract.Call(opts, out, "getParamAtNextBlock", name)
	return *ret0, err
}

// GetParamAtNextBlock is a free data retrieval call binding the contract method 0x9c9c03bb.
//
// Solidity: function getParamAtNextBlock(string name) view returns(bytes)
func (_GovParam *GovParamSession) GetParamAtNextBlock(name string) ([]byte, error) {
	return _GovParam.Contract.GetParamAtNextBlock(&_GovParam.CallOpts, name)
}

// GetParamAtNextBlock is a free data retrieval call binding the contract method 0x9c9c03bb.
//
// Solidity: function getParamAtNextBlock(string name) view returns(bytes)
func (_GovParam *GovParamCallerSession) GetParamAtNextBlock(name string) ([]byte, error) {
	return _GovParam.Contract.GetParamAtNextBlock(&_GovParam.CallOpts, name)
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

// SetParam is a paid mutator transaction binding the contract method 0x8841067a.
//
// Solidity: function setParam(string name, bytes value, uint64 activation) returns()
func (_GovParam *GovParamTransactor) SetParam(opts *bind.TransactOpts, name string, value []byte, activation uint64) (*types.Transaction, error) {
	return _GovParam.contract.Transact(opts, "setParam", name, value, activation)
}

// SetParam is a paid mutator transaction binding the contract method 0x8841067a.
//
// Solidity: function setParam(string name, bytes value, uint64 activation) returns()
func (_GovParam *GovParamSession) SetParam(name string, value []byte, activation uint64) (*types.Transaction, error) {
	return _GovParam.Contract.SetParam(&_GovParam.TransactOpts, name, value, activation)
}

// SetParam is a paid mutator transaction binding the contract method 0x8841067a.
//
// Solidity: function setParam(string name, bytes value, uint64 activation) returns()
func (_GovParam *GovParamTransactorSession) SetParam(name string, value []byte, activation uint64) (*types.Transaction, error) {
	return _GovParam.Contract.SetParam(&_GovParam.TransactOpts, name, value, activation)
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
	Arg2 uint64
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterSetParam is a free log retrieval operation binding the contract event 0xf8d4c4710beeb9d48e0bc6888dfcaf21467794f5da39c3bfb2874a2ac3a73c56.
//
// Solidity: event SetParam(string arg0, bytes arg1, uint64 arg2)
func (_GovParam *GovParamFilterer) FilterSetParam(opts *bind.FilterOpts) (*GovParamSetParamIterator, error) {

	logs, sub, err := _GovParam.contract.FilterLogs(opts, "SetParam")
	if err != nil {
		return nil, err
	}
	return &GovParamSetParamIterator{contract: _GovParam.contract, event: "SetParam", logs: logs, sub: sub}, nil
}

// WatchSetParam is a free log subscription operation binding the contract event 0xf8d4c4710beeb9d48e0bc6888dfcaf21467794f5da39c3bfb2874a2ac3a73c56.
//
// Solidity: event SetParam(string arg0, bytes arg1, uint64 arg2)
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

// ParseSetParam is a log parse operation binding the contract event 0xf8d4c4710beeb9d48e0bc6888dfcaf21467794f5da39c3bfb2874a2ac3a73c56.
//
// Solidity: event SetParam(string arg0, bytes arg1, uint64 arg2)
func (_GovParam *GovParamFilterer) ParseSetParam(log types.Log) (*GovParamSetParam, error) {
	event := new(GovParamSetParam)
	if err := _GovParam.contract.UnpackLog(event, "SetParam", log); err != nil {
		return nil, err
	}
	return event, nil
}

// OwnableABI is the input ABI used to generate the binding from.
const OwnableABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// OwnableBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const OwnableBinRuntime = ``

// OwnableFuncSigs maps the 4-byte function signature to its string representation.
var OwnableFuncSigs = map[string]string{
	"8da5cb5b": "owner()",
	"715018a6": "renounceOwnership()",
	"f2fde38b": "transferOwnership(address)",
}

// Ownable is an auto generated Go binding around a Klaytn contract.
type Ownable struct {
	OwnableCaller     // Read-only binding to the contract
	OwnableTransactor // Write-only binding to the contract
	OwnableFilterer   // Log filterer for contract events
}

// OwnableCaller is an auto generated read-only Go binding around a Klaytn contract.
type OwnableCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OwnableTransactor is an auto generated write-only Go binding around a Klaytn contract.
type OwnableTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OwnableFilterer is an auto generated log filtering Go binding around a Klaytn contract events.
type OwnableFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OwnableSession is an auto generated Go binding around a Klaytn contract,
// with pre-set call and transact options.
type OwnableSession struct {
	Contract     *Ownable          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// OwnableCallerSession is an auto generated read-only Go binding around a Klaytn contract,
// with pre-set call options.
type OwnableCallerSession struct {
	Contract *OwnableCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// OwnableTransactorSession is an auto generated write-only Go binding around a Klaytn contract,
// with pre-set transact options.
type OwnableTransactorSession struct {
	Contract     *OwnableTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// OwnableRaw is an auto generated low-level Go binding around a Klaytn contract.
type OwnableRaw struct {
	Contract *Ownable // Generic contract binding to access the raw methods on
}

// OwnableCallerRaw is an auto generated low-level read-only Go binding around a Klaytn contract.
type OwnableCallerRaw struct {
	Contract *OwnableCaller // Generic read-only contract binding to access the raw methods on
}

// OwnableTransactorRaw is an auto generated low-level write-only Go binding around a Klaytn contract.
type OwnableTransactorRaw struct {
	Contract *OwnableTransactor // Generic write-only contract binding to access the raw methods on
}

// NewOwnable creates a new instance of Ownable, bound to a specific deployed contract.
func NewOwnable(address common.Address, backend bind.ContractBackend) (*Ownable, error) {
	contract, err := bindOwnable(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Ownable{OwnableCaller: OwnableCaller{contract: contract}, OwnableTransactor: OwnableTransactor{contract: contract}, OwnableFilterer: OwnableFilterer{contract: contract}}, nil
}

// NewOwnableCaller creates a new read-only instance of Ownable, bound to a specific deployed contract.
func NewOwnableCaller(address common.Address, caller bind.ContractCaller) (*OwnableCaller, error) {
	contract, err := bindOwnable(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OwnableCaller{contract: contract}, nil
}

// NewOwnableTransactor creates a new write-only instance of Ownable, bound to a specific deployed contract.
func NewOwnableTransactor(address common.Address, transactor bind.ContractTransactor) (*OwnableTransactor, error) {
	contract, err := bindOwnable(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OwnableTransactor{contract: contract}, nil
}

// NewOwnableFilterer creates a new log filterer instance of Ownable, bound to a specific deployed contract.
func NewOwnableFilterer(address common.Address, filterer bind.ContractFilterer) (*OwnableFilterer, error) {
	contract, err := bindOwnable(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OwnableFilterer{contract: contract}, nil
}

// bindOwnable binds a generic wrapper to an already deployed contract.
func bindOwnable(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(OwnableABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Ownable *OwnableRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Ownable.Contract.OwnableCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Ownable *OwnableRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Ownable.Contract.OwnableTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Ownable *OwnableRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Ownable.Contract.OwnableTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Ownable *OwnableCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Ownable.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Ownable *OwnableTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Ownable.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Ownable *OwnableTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Ownable.Contract.contract.Transact(opts, method, params...)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Ownable *OwnableCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Ownable.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Ownable *OwnableSession) Owner() (common.Address, error) {
	return _Ownable.Contract.Owner(&_Ownable.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Ownable *OwnableCallerSession) Owner() (common.Address, error) {
	return _Ownable.Contract.Owner(&_Ownable.CallOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Ownable *OwnableTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Ownable.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Ownable *OwnableSession) RenounceOwnership() (*types.Transaction, error) {
	return _Ownable.Contract.RenounceOwnership(&_Ownable.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Ownable *OwnableTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Ownable.Contract.RenounceOwnership(&_Ownable.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Ownable *OwnableTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Ownable.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Ownable *OwnableSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Ownable.Contract.TransferOwnership(&_Ownable.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Ownable *OwnableTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Ownable.Contract.TransferOwnership(&_Ownable.TransactOpts, newOwner)
}

// OwnableOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Ownable contract.
type OwnableOwnershipTransferredIterator struct {
	Event *OwnableOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *OwnableOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OwnableOwnershipTransferred)
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
		it.Event = new(OwnableOwnershipTransferred)
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
func (it *OwnableOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OwnableOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OwnableOwnershipTransferred represents a OwnershipTransferred event raised by the Ownable contract.
type OwnableOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Ownable *OwnableFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*OwnableOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Ownable.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &OwnableOwnershipTransferredIterator{contract: _Ownable.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Ownable *OwnableFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *OwnableOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Ownable.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OwnableOwnershipTransferred)
				if err := _Ownable.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_Ownable *OwnableFilterer) ParseOwnershipTransferred(log types.Log) (*OwnableOwnershipTransferred, error) {
	event := new(OwnableOwnershipTransferred)
	if err := _Ownable.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	return event, nil
}
