// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package nodewhitelist

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

// NodeWhitelistABI is the input ABI used to generate the binding from.
const NodeWhitelistABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"addedNode\",\"type\":\"string\"}],\"name\":\"AddNode\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"deletedNode\",\"type\":\"string\"}],\"name\":\"DelNode\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"node\",\"type\":\"string\"}],\"name\":\"addNode\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"node\",\"type\":\"string\"}],\"name\":\"delNode\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAdmin\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getWhitelist\",\"outputs\":[{\"internalType\":\"string[]\",\"name\":\"\",\"type\":\"string[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"setAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// NodeWhitelistBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const NodeWhitelistBinRuntime = `608060405234801561001057600080fd5b50600436106100575760003560e01c806311c903051461005c5780636e9960c314610071578063704b6c02146100915780638994dd8e146100a4578063d01f63f5146100b7575b600080fd5b61006f61006a366004610669565b6100cc565b005b6000546040516001600160a01b0390911681526020015b60405180910390f35b61006f61009f36600461071a565b6101c0565b61006f6100b2366004610669565b61020c565b6100bf61031b565b6040516100889190610797565b6000546001600160a01b031633146100ff5760405162461bcd60e51b81526004016100f6906107f9565b60405180910390fd5b610108816103f4565b60001914156101655760405162461bcd60e51b815260206004820152602360248201527f676976656e206e6f6465206973206e6f74206f6e207468652077686974656c6960448201526273742160e81b60648201526084016100f6565b610176610171826103f4565b610464565b6000546040517ffd752733394b8065fe1f8210061450782b2b842513735a01823786ecbf8889ae916101b5916001600160a01b03909116908490610827565b60405180910390a150565b6000546001600160a01b031633146101ea5760405162461bcd60e51b81526004016100f6906107f9565b600080546001600160a01b0319166001600160a01b0392909216919091179055565b6000546001600160a01b031633146102365760405162461bcd60e51b81526004016100f6906107f9565b61023f816103f4565b600019146102995760405162461bcd60e51b815260206004820152602160248201527f676976656e206e6f646520697320616c726561647920726567697374657265646044820152602160f81b60648201526084016100f6565b60018054808201825560009190915281516102db917fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf601906020840190610502565b506000546040517f1467d566c65746255852dd30f89b5cd7a10ab646a0dccaa3e6348725a5810099916101b5916001600160a01b03909116908490610827565b60606001805480602002602001604051908101604052809291908181526020016000905b828210156103eb57838290600052602060002001805461035e90610853565b80601f016020809104026020016040519081016040528092919081815260200182805461038a90610853565b80156103d75780601f106103ac576101008083540402835291602001916103d7565b820191906000526020600020905b8154815290600101906020018083116103ba57829003601f168201915b50505050508152602001906001019061033f565b50505050905090565b6000805b60015481101561045a5782805190602001206001828154811061041d5761041d61088e565b9060005260206000200160405161043491906108a4565b604051809103902014156104485792915050565b8061045281610956565b9150506103f8565b5060001992915050565b600154811061047257600080fd5b60018054610481908290610971565b815481106104915761049161088e565b90600052602060002001600182815481106104ae576104ae61088e565b906000526020600020019080546104c490610853565b6104cf929190610586565b5060018054806104e1576104e1610988565b6001900381819060005260206000200160006104fd9190610601565b905550565b82805461050e90610853565b90600052602060002090601f0160209004810192826105305760008555610576565b82601f1061054957805160ff1916838001178555610576565b82800160010185558215610576579182015b8281111561057657825182559160200191906001019061055b565b5061058292915061063e565b5090565b82805461059290610853565b90600052602060002090601f0160209004810192826105b45760008555610576565b82601f106105c55780548555610576565b8280016001018555821561057657600052602060002091601f016020900482015b828111156105765782548255916001019190600101906105e6565b50805461060d90610853565b6000825580601f1061061d575050565b601f01602090049060005260206000209081019061063b919061063e565b50565b5b80821115610582576000815560010161063f565b634e487b7160e01b600052604160045260246000fd5b60006020828403121561067b57600080fd5b813567ffffffffffffffff8082111561069357600080fd5b818401915084601f8301126106a757600080fd5b8135818111156106b9576106b9610653565b604051601f8201601f19908116603f011681019083821181831017156106e1576106e1610653565b816040528281528760208487010111156106fa57600080fd5b826020860160208301376000928101602001929092525095945050505050565b60006020828403121561072c57600080fd5b81356001600160a01b038116811461074357600080fd5b9392505050565b6000815180845260005b8181101561077057602081850181015186830182015201610754565b81811115610782576000602083870101525b50601f01601f19169290920160200192915050565b6000602080830181845280855180835260408601915060408160051b870101925083870160005b828110156107ec57603f198886030184526107da85835161074a565b945092850192908501906001016107be565b5092979650505050505050565b60208082526014908201527327b7363c9030b236b4b71031b0b71031b0b6361760611b604082015260600190565b6001600160a01b038316815260406020820181905260009061084b9083018461074a565b949350505050565b600181811c9082168061086757607f821691505b6020821081141561088857634e487b7160e01b600052602260045260246000fd5b50919050565b634e487b7160e01b600052603260045260246000fd5b600080835481600182811c9150808316806108c057607f831692505b60208084108214156108e057634e487b7160e01b86526022600452602486fd5b8180156108f4576001811461090557610932565b60ff19861689528489019650610932565b60008a81526020902060005b8681101561092a5781548b820152908501908301610911565b505084890196505b509498975050505050505050565b634e487b7160e01b600052601160045260246000fd5b600060001982141561096a5761096a610940565b5060010190565b60008282101561098357610983610940565b500390565b634e487b7160e01b600052603160045260246000fdfea2646970667358221220b5b3a4d3c7196b6cd72144f9e9da9d14c1dafefd40e09111333247b44bcc590564736f6c63430008090033`

// NodeWhitelistFuncSigs maps the 4-byte function signature to its string representation.
var NodeWhitelistFuncSigs = map[string]string{
	"8994dd8e": "addNode(string)",
	"11c90305": "delNode(string)",
	"6e9960c3": "getAdmin()",
	"d01f63f5": "getWhitelist()",
	"704b6c02": "setAdmin(address)",
}

// NodeWhitelistBin is the compiled bytecode used for deploying new contracts.
var NodeWhitelistBin = "0x608060405234801561001057600080fd5b506109d4806100206000396000f3fe608060405234801561001057600080fd5b50600436106100575760003560e01c806311c903051461005c5780636e9960c314610071578063704b6c02146100915780638994dd8e146100a4578063d01f63f5146100b7575b600080fd5b61006f61006a366004610669565b6100cc565b005b6000546040516001600160a01b0390911681526020015b60405180910390f35b61006f61009f36600461071a565b6101c0565b61006f6100b2366004610669565b61020c565b6100bf61031b565b6040516100889190610797565b6000546001600160a01b031633146100ff5760405162461bcd60e51b81526004016100f6906107f9565b60405180910390fd5b610108816103f4565b60001914156101655760405162461bcd60e51b815260206004820152602360248201527f676976656e206e6f6465206973206e6f74206f6e207468652077686974656c6960448201526273742160e81b60648201526084016100f6565b610176610171826103f4565b610464565b6000546040517ffd752733394b8065fe1f8210061450782b2b842513735a01823786ecbf8889ae916101b5916001600160a01b03909116908490610827565b60405180910390a150565b6000546001600160a01b031633146101ea5760405162461bcd60e51b81526004016100f6906107f9565b600080546001600160a01b0319166001600160a01b0392909216919091179055565b6000546001600160a01b031633146102365760405162461bcd60e51b81526004016100f6906107f9565b61023f816103f4565b600019146102995760405162461bcd60e51b815260206004820152602160248201527f676976656e206e6f646520697320616c726561647920726567697374657265646044820152602160f81b60648201526084016100f6565b60018054808201825560009190915281516102db917fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf601906020840190610502565b506000546040517f1467d566c65746255852dd30f89b5cd7a10ab646a0dccaa3e6348725a5810099916101b5916001600160a01b03909116908490610827565b60606001805480602002602001604051908101604052809291908181526020016000905b828210156103eb57838290600052602060002001805461035e90610853565b80601f016020809104026020016040519081016040528092919081815260200182805461038a90610853565b80156103d75780601f106103ac576101008083540402835291602001916103d7565b820191906000526020600020905b8154815290600101906020018083116103ba57829003601f168201915b50505050508152602001906001019061033f565b50505050905090565b6000805b60015481101561045a5782805190602001206001828154811061041d5761041d61088e565b9060005260206000200160405161043491906108a4565b604051809103902014156104485792915050565b8061045281610956565b9150506103f8565b5060001992915050565b600154811061047257600080fd5b60018054610481908290610971565b815481106104915761049161088e565b90600052602060002001600182815481106104ae576104ae61088e565b906000526020600020019080546104c490610853565b6104cf929190610586565b5060018054806104e1576104e1610988565b6001900381819060005260206000200160006104fd9190610601565b905550565b82805461050e90610853565b90600052602060002090601f0160209004810192826105305760008555610576565b82601f1061054957805160ff1916838001178555610576565b82800160010185558215610576579182015b8281111561057657825182559160200191906001019061055b565b5061058292915061063e565b5090565b82805461059290610853565b90600052602060002090601f0160209004810192826105b45760008555610576565b82601f106105c55780548555610576565b8280016001018555821561057657600052602060002091601f016020900482015b828111156105765782548255916001019190600101906105e6565b50805461060d90610853565b6000825580601f1061061d575050565b601f01602090049060005260206000209081019061063b919061063e565b50565b5b80821115610582576000815560010161063f565b634e487b7160e01b600052604160045260246000fd5b60006020828403121561067b57600080fd5b813567ffffffffffffffff8082111561069357600080fd5b818401915084601f8301126106a757600080fd5b8135818111156106b9576106b9610653565b604051601f8201601f19908116603f011681019083821181831017156106e1576106e1610653565b816040528281528760208487010111156106fa57600080fd5b826020860160208301376000928101602001929092525095945050505050565b60006020828403121561072c57600080fd5b81356001600160a01b038116811461074357600080fd5b9392505050565b6000815180845260005b8181101561077057602081850181015186830182015201610754565b81811115610782576000602083870101525b50601f01601f19169290920160200192915050565b6000602080830181845280855180835260408601915060408160051b870101925083870160005b828110156107ec57603f198886030184526107da85835161074a565b945092850192908501906001016107be565b5092979650505050505050565b60208082526014908201527327b7363c9030b236b4b71031b0b71031b0b6361760611b604082015260600190565b6001600160a01b038316815260406020820181905260009061084b9083018461074a565b949350505050565b600181811c9082168061086757607f821691505b6020821081141561088857634e487b7160e01b600052602260045260246000fd5b50919050565b634e487b7160e01b600052603260045260246000fd5b600080835481600182811c9150808316806108c057607f831692505b60208084108214156108e057634e487b7160e01b86526022600452602486fd5b8180156108f4576001811461090557610932565b60ff19861689528489019650610932565b60008a81526020902060005b8681101561092a5781548b820152908501908301610911565b505084890196505b509498975050505050505050565b634e487b7160e01b600052601160045260246000fd5b600060001982141561096a5761096a610940565b5060010190565b60008282101561098357610983610940565b500390565b634e487b7160e01b600052603160045260246000fdfea2646970667358221220b5b3a4d3c7196b6cd72144f9e9da9d14c1dafefd40e09111333247b44bcc590564736f6c63430008090033"

// DeployNodeWhitelist deploys a new Klaytn contract, binding an instance of NodeWhitelist to it.
func DeployNodeWhitelist(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *NodeWhitelist, error) {
	parsed, err := abi.JSON(strings.NewReader(NodeWhitelistABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(NodeWhitelistBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &NodeWhitelist{NodeWhitelistCaller: NodeWhitelistCaller{contract: contract}, NodeWhitelistTransactor: NodeWhitelistTransactor{contract: contract}, NodeWhitelistFilterer: NodeWhitelistFilterer{contract: contract}}, nil
}

// NodeWhitelist is an auto generated Go binding around a Klaytn contract.
type NodeWhitelist struct {
	NodeWhitelistCaller     // Read-only binding to the contract
	NodeWhitelistTransactor // Write-only binding to the contract
	NodeWhitelistFilterer   // Log filterer for contract events
}

// NodeWhitelistCaller is an auto generated read-only Go binding around a Klaytn contract.
type NodeWhitelistCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NodeWhitelistTransactor is an auto generated write-only Go binding around a Klaytn contract.
type NodeWhitelistTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NodeWhitelistFilterer is an auto generated log filtering Go binding around a Klaytn contract events.
type NodeWhitelistFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NodeWhitelistSession is an auto generated Go binding around a Klaytn contract,
// with pre-set call and transact options.
type NodeWhitelistSession struct {
	Contract     *NodeWhitelist    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// NodeWhitelistCallerSession is an auto generated read-only Go binding around a Klaytn contract,
// with pre-set call options.
type NodeWhitelistCallerSession struct {
	Contract *NodeWhitelistCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// NodeWhitelistTransactorSession is an auto generated write-only Go binding around a Klaytn contract,
// with pre-set transact options.
type NodeWhitelistTransactorSession struct {
	Contract     *NodeWhitelistTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// NodeWhitelistRaw is an auto generated low-level Go binding around a Klaytn contract.
type NodeWhitelistRaw struct {
	Contract *NodeWhitelist // Generic contract binding to access the raw methods on
}

// NodeWhitelistCallerRaw is an auto generated low-level read-only Go binding around a Klaytn contract.
type NodeWhitelistCallerRaw struct {
	Contract *NodeWhitelistCaller // Generic read-only contract binding to access the raw methods on
}

// NodeWhitelistTransactorRaw is an auto generated low-level write-only Go binding around a Klaytn contract.
type NodeWhitelistTransactorRaw struct {
	Contract *NodeWhitelistTransactor // Generic write-only contract binding to access the raw methods on
}

// NewNodeWhitelist creates a new instance of NodeWhitelist, bound to a specific deployed contract.
func NewNodeWhitelist(address common.Address, backend bind.ContractBackend) (*NodeWhitelist, error) {
	contract, err := bindNodeWhitelist(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &NodeWhitelist{NodeWhitelistCaller: NodeWhitelistCaller{contract: contract}, NodeWhitelistTransactor: NodeWhitelistTransactor{contract: contract}, NodeWhitelistFilterer: NodeWhitelistFilterer{contract: contract}}, nil
}

// NewNodeWhitelistCaller creates a new read-only instance of NodeWhitelist, bound to a specific deployed contract.
func NewNodeWhitelistCaller(address common.Address, caller bind.ContractCaller) (*NodeWhitelistCaller, error) {
	contract, err := bindNodeWhitelist(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &NodeWhitelistCaller{contract: contract}, nil
}

// NewNodeWhitelistTransactor creates a new write-only instance of NodeWhitelist, bound to a specific deployed contract.
func NewNodeWhitelistTransactor(address common.Address, transactor bind.ContractTransactor) (*NodeWhitelistTransactor, error) {
	contract, err := bindNodeWhitelist(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &NodeWhitelistTransactor{contract: contract}, nil
}

// NewNodeWhitelistFilterer creates a new log filterer instance of NodeWhitelist, bound to a specific deployed contract.
func NewNodeWhitelistFilterer(address common.Address, filterer bind.ContractFilterer) (*NodeWhitelistFilterer, error) {
	contract, err := bindNodeWhitelist(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &NodeWhitelistFilterer{contract: contract}, nil
}

// bindNodeWhitelist binds a generic wrapper to an already deployed contract.
func bindNodeWhitelist(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(NodeWhitelistABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_NodeWhitelist *NodeWhitelistRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _NodeWhitelist.Contract.NodeWhitelistCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_NodeWhitelist *NodeWhitelistRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NodeWhitelist.Contract.NodeWhitelistTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_NodeWhitelist *NodeWhitelistRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _NodeWhitelist.Contract.NodeWhitelistTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_NodeWhitelist *NodeWhitelistCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _NodeWhitelist.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_NodeWhitelist *NodeWhitelistTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NodeWhitelist.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_NodeWhitelist *NodeWhitelistTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _NodeWhitelist.Contract.contract.Transact(opts, method, params...)
}

// GetAdmin is a free data retrieval call binding the contract method 0x6e9960c3.
//
// Solidity: function getAdmin() view returns(address)
func (_NodeWhitelist *NodeWhitelistCaller) GetAdmin(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _NodeWhitelist.contract.Call(opts, out, "getAdmin")
	return *ret0, err
}

// GetAdmin is a free data retrieval call binding the contract method 0x6e9960c3.
//
// Solidity: function getAdmin() view returns(address)
func (_NodeWhitelist *NodeWhitelistSession) GetAdmin() (common.Address, error) {
	return _NodeWhitelist.Contract.GetAdmin(&_NodeWhitelist.CallOpts)
}

// GetAdmin is a free data retrieval call binding the contract method 0x6e9960c3.
//
// Solidity: function getAdmin() view returns(address)
func (_NodeWhitelist *NodeWhitelistCallerSession) GetAdmin() (common.Address, error) {
	return _NodeWhitelist.Contract.GetAdmin(&_NodeWhitelist.CallOpts)
}

// GetWhitelist is a free data retrieval call binding the contract method 0xd01f63f5.
//
// Solidity: function getWhitelist() view returns(string[])
func (_NodeWhitelist *NodeWhitelistCaller) GetWhitelist(opts *bind.CallOpts) ([]string, error) {
	var (
		ret0 = new([]string)
	)
	out := ret0
	err := _NodeWhitelist.contract.Call(opts, out, "getWhitelist")
	return *ret0, err
}

// GetWhitelist is a free data retrieval call binding the contract method 0xd01f63f5.
//
// Solidity: function getWhitelist() view returns(string[])
func (_NodeWhitelist *NodeWhitelistSession) GetNodeWhitelist() ([]string, error) {
	return _NodeWhitelist.Contract.GetWhitelist(&_NodeWhitelist.CallOpts)
}

// GetWhitelist is a free data retrieval call binding the contract method 0xd01f63f5.
//
// Solidity: function getWhitelist() view returns(string[])
func (_NodeWhitelist *NodeWhitelistCallerSession) GetNodeWhitelist() ([]string, error) {
	return _NodeWhitelist.Contract.GetWhitelist(&_NodeWhitelist.CallOpts)
}

// AddNode is a paid mutator transaction binding the contract method 0x8994dd8e.
//
// Solidity: function addNode(string node) returns()
func (_NodeWhitelist *NodeWhitelistTransactor) AddNode(opts *bind.TransactOpts, node string) (*types.Transaction, error) {
	return _NodeWhitelist.contract.Transact(opts, "addNode", node)
}

// AddNode is a paid mutator transaction binding the contract method 0x8994dd8e.
//
// Solidity: function addNode(string node) returns()
func (_NodeWhitelist *NodeWhitelistSession) AddNode(node string) (*types.Transaction, error) {
	return _NodeWhitelist.Contract.AddNode(&_NodeWhitelist.TransactOpts, node)
}

// AddNode is a paid mutator transaction binding the contract method 0x8994dd8e.
//
// Solidity: function addNode(string node) returns()
func (_NodeWhitelist *NodeWhitelistTransactorSession) AddNode(node string) (*types.Transaction, error) {
	return _NodeWhitelist.Contract.AddNode(&_NodeWhitelist.TransactOpts, node)
}

// DelNode is a paid mutator transaction binding the contract method 0x11c90305.
//
// Solidity: function delNode(string node) returns()
func (_NodeWhitelist *NodeWhitelistTransactor) DelNode(opts *bind.TransactOpts, node string) (*types.Transaction, error) {
	return _NodeWhitelist.contract.Transact(opts, "delNode", node)
}

// DelNode is a paid mutator transaction binding the contract method 0x11c90305.
//
// Solidity: function delNode(string node) returns()
func (_NodeWhitelist *NodeWhitelistSession) DelNode(node string) (*types.Transaction, error) {
	return _NodeWhitelist.Contract.DelNode(&_NodeWhitelist.TransactOpts, node)
}

// DelNode is a paid mutator transaction binding the contract method 0x11c90305.
//
// Solidity: function delNode(string node) returns()
func (_NodeWhitelist *NodeWhitelistTransactorSession) DelNode(node string) (*types.Transaction, error) {
	return _NodeWhitelist.Contract.DelNode(&_NodeWhitelist.TransactOpts, node)
}

// SetAdmin is a paid mutator transaction binding the contract method 0x704b6c02.
//
// Solidity: function setAdmin(address newAdmin) returns()
func (_NodeWhitelist *NodeWhitelistTransactor) SetAdmin(opts *bind.TransactOpts, newAdmin common.Address) (*types.Transaction, error) {
	return _NodeWhitelist.contract.Transact(opts, "setAdmin", newAdmin)
}

// SetAdmin is a paid mutator transaction binding the contract method 0x704b6c02.
//
// Solidity: function setAdmin(address newAdmin) returns()
func (_NodeWhitelist *NodeWhitelistSession) SetAdmin(newAdmin common.Address) (*types.Transaction, error) {
	return _NodeWhitelist.Contract.SetAdmin(&_NodeWhitelist.TransactOpts, newAdmin)
}

// SetAdmin is a paid mutator transaction binding the contract method 0x704b6c02.
//
// Solidity: function setAdmin(address newAdmin) returns()
func (_NodeWhitelist *NodeWhitelistTransactorSession) SetAdmin(newAdmin common.Address) (*types.Transaction, error) {
	return _NodeWhitelist.Contract.SetAdmin(&_NodeWhitelist.TransactOpts, newAdmin)
}

// NodeWhitelistAddNodeIterator is returned from FilterAddNode and is used to iterate over the raw logs and unpacked data for AddNode events raised by the NodeWhitelist contract.
type NodeWhitelistAddNodeIterator struct {
	Event *NodeWhitelistAddNode // Event containing the contract specifics and raw log

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
func (it *NodeWhitelistAddNodeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodeWhitelistAddNode)
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
		it.Event = new(NodeWhitelistAddNode)
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
func (it *NodeWhitelistAddNodeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodeWhitelistAddNodeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodeWhitelistAddNode represents a AddNode event raised by the NodeWhitelist contract.
type NodeWhitelistAddNode struct {
	Admin     common.Address
	AddedNode string
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterAddNode is a free log retrieval operation binding the contract event 0x1467d566c65746255852dd30f89b5cd7a10ab646a0dccaa3e6348725a5810099.
//
// Solidity: event AddNode(address admin, string addedNode)
func (_NodeWhitelist *NodeWhitelistFilterer) FilterAddNode(opts *bind.FilterOpts) (*NodeWhitelistAddNodeIterator, error) {

	logs, sub, err := _NodeWhitelist.contract.FilterLogs(opts, "AddNode")
	if err != nil {
		return nil, err
	}
	return &NodeWhitelistAddNodeIterator{contract: _NodeWhitelist.contract, event: "AddNode", logs: logs, sub: sub}, nil
}

// WatchAddNode is a free log subscription operation binding the contract event 0x1467d566c65746255852dd30f89b5cd7a10ab646a0dccaa3e6348725a5810099.
//
// Solidity: event AddNode(address admin, string addedNode)
func (_NodeWhitelist *NodeWhitelistFilterer) WatchAddNode(opts *bind.WatchOpts, sink chan<- *NodeWhitelistAddNode) (event.Subscription, error) {

	logs, sub, err := _NodeWhitelist.contract.WatchLogs(opts, "AddNode")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodeWhitelistAddNode)
				if err := _NodeWhitelist.contract.UnpackLog(event, "AddNode", log); err != nil {
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

// ParseAddNode is a log parse operation binding the contract event 0x1467d566c65746255852dd30f89b5cd7a10ab646a0dccaa3e6348725a5810099.
//
// Solidity: event AddNode(address admin, string addedNode)
func (_NodeWhitelist *NodeWhitelistFilterer) ParseAddNode(log types.Log) (*NodeWhitelistAddNode, error) {
	event := new(NodeWhitelistAddNode)
	if err := _NodeWhitelist.contract.UnpackLog(event, "AddNode", log); err != nil {
		return nil, err
	}
	return event, nil
}

// NodeWhitelistDelNodeIterator is returned from FilterDelNode and is used to iterate over the raw logs and unpacked data for DelNode events raised by the NodeWhitelist contract.
type NodeWhitelistDelNodeIterator struct {
	Event *NodeWhitelistDelNode // Event containing the contract specifics and raw log

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
func (it *NodeWhitelistDelNodeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodeWhitelistDelNode)
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
		it.Event = new(NodeWhitelistDelNode)
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
func (it *NodeWhitelistDelNodeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodeWhitelistDelNodeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodeWhitelistDelNode represents a DelNode event raised by the NodeWhitelist contract.
type NodeWhitelistDelNode struct {
	Admin       common.Address
	DeletedNode string
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterDelNode is a free log retrieval operation binding the contract event 0xfd752733394b8065fe1f8210061450782b2b842513735a01823786ecbf8889ae.
//
// Solidity: event DelNode(address admin, string deletedNode)
func (_NodeWhitelist *NodeWhitelistFilterer) FilterDelNode(opts *bind.FilterOpts) (*NodeWhitelistDelNodeIterator, error) {

	logs, sub, err := _NodeWhitelist.contract.FilterLogs(opts, "DelNode")
	if err != nil {
		return nil, err
	}
	return &NodeWhitelistDelNodeIterator{contract: _NodeWhitelist.contract, event: "DelNode", logs: logs, sub: sub}, nil
}

// WatchDelNode is a free log subscription operation binding the contract event 0xfd752733394b8065fe1f8210061450782b2b842513735a01823786ecbf8889ae.
//
// Solidity: event DelNode(address admin, string deletedNode)
func (_NodeWhitelist *NodeWhitelistFilterer) WatchDelNode(opts *bind.WatchOpts, sink chan<- *NodeWhitelistDelNode) (event.Subscription, error) {

	logs, sub, err := _NodeWhitelist.contract.WatchLogs(opts, "DelNode")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodeWhitelistDelNode)
				if err := _NodeWhitelist.contract.UnpackLog(event, "DelNode", log); err != nil {
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

// ParseDelNode is a log parse operation binding the contract event 0xfd752733394b8065fe1f8210061450782b2b842513735a01823786ecbf8889ae.
//
// Solidity: event DelNode(address admin, string deletedNode)
func (_NodeWhitelist *NodeWhitelistFilterer) ParseDelNode(log types.Log) (*NodeWhitelistDelNode, error) {
	event := new(NodeWhitelistDelNode)
	if err := _NodeWhitelist.contract.UnpackLog(event, "DelNode", log); err != nil {
		return nil, err
	}
	return event, nil
}
