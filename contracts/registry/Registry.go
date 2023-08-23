// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package registry

import (
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
	_ = big.NewInt
	_ = strings.NewReader
	_ = klaytn.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)
// RegistryMetaData contains all meta data concerning the Registry contract.
var RegistryMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"activationBlockNumber\",\"type\":\"uint256\"}],\"name\":\"Activated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"governance\",\"type\":\"address\"}],\"name\":\"ConstructContract\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"deprecateBlockNumber\",\"type\":\"uint256\"}],\"name\":\"Deprecated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"Registered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"prevName\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"newName\",\"type\":\"string\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"replaceBlockNumber\",\"type\":\"uint256\"}],\"name\":\"Replaced\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newGovernance\",\"type\":\"address\"}],\"name\":\"UpdateGovernance\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"activationBlockNumber\",\"type\":\"uint256\"}],\"name\":\"activate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_governance\",\"type\":\"address\"}],\"name\":\"constructContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"contractNames\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"deprecateBlockNumber\",\"type\":\"uint256\"}],\"name\":\"deprecate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllContractNames\",\"outputs\":[{\"internalType\":\"string[]\",\"name\":\"\",\"type\":\"string[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"name\":\"getContractIfActive\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"governance\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"internalType\":\"enumIRegistry.State\",\"name\":\"state\",\"type\":\"uint8\"}],\"name\":\"readAllContractsAtGivenState\",\"outputs\":[{\"internalType\":\"string[]\",\"name\":\"\",\"type\":\"string[]\"},{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"activation\",\"type\":\"bool\"}],\"name\":\"register\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"name\":\"registry\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"activationBlockNumber\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deprecateBlockNumber\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"prevName\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"newName\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"replaceBlockNumber\",\"type\":\"uint256\"}],\"name\":\"replace\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"}],\"name\":\"stateAt\",\"outputs\":[{\"internalType\":\"enumIRegistry.State\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_newGovernance\",\"type\":\"address\"}],\"name\":\"updateGovernance\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

const RegistryBinRuntime = "0x608060405234801561001057600080fd5b50600436106100cf5760003560e01c806360d7a2781161008c578063d7a9b81611610066578063d7a9b81614610212578063e14c672d14610227578063e47f11601461023a578063fa96b37f1461024d57600080fd5b806360d7a2781461018357806392a296c914610196578063b2561263146101ff57600080fd5b806305e4f16e146100d457806335b5cfbb146100fd5780633ca6bb92146101285780634d0e66061461014857806351afb6a31461015d5780635aa6e67514610170575b600080fd5b6100e76100e236600461106e565b61026e565b6040516100f491906110c9565b60405180910390f35b61011061010b3660046110f1565b6103e5565b6040516001600160a01b0390911681526020016100f4565b61013b61013636600461112e565b610476565b6040516100f49190611197565b61015b61015636600461106e565b610522565b005b61015b61016b36600461106e565b610645565b600254610110906001600160a01b031681565b61015b6101913660046111cd565b610758565b6101da6101a43660046110f1565b80516020818301810180516000825292820191909301209152805460018201546002909201546001600160a01b03909116919083565b604080516001600160a01b0390941684526020840192909252908201526060016100f4565b61015b61020d366004611234565b610983565b61021a610a42565b6040516100f491906112a4565b61015b610235366004611234565b610b1b565b61015b6102483660046112b7565b610b8d565b61026061025b366004611324565b610be0565b6040516100f4929190611358565b60008260008160405160200161028491906113b8565b604051602081830303815290604052905080516000036102da5760405162461bcd60e51b815260206004820152600c60248201526b456d70747920737472696e6760a01b60448201526064015b60405180910390fd5b6000856040516102ea91906113b8565b908152604051908190036020019020546001600160a01b0316806103205760405162461bcd60e51b81526004016102d1906113d4565b6000808760405161033191906113b8565b90815260200160405180910390206001015490506000808860405161035691906113b8565b90815260200160405180910390206002015490508160000361037d576000955050506103dc565b806000036103a35781871015610398576000955050506103dc565b6001955050506103dc565b818710156103b6576000955050506103dc565b8187101580156103c557508087105b156103d5576001955050506103dc565b6002955050505b50505092915050565b600060016103f3834361026e565b6002811115610404576104046110b3565b146104475760405162461bcd60e51b8152602060048201526013602482015272139bdd081858dd1a5d994818dbdb9d1c9858dd606a1b60448201526064016102d1565b60008260405161045791906113b8565b908152604051908190036020019020546001600160a01b031692915050565b6001818154811061048657600080fd5b9060005260206000200160009150905080546104a1906113fa565b80601f01602080910402602001604051908101604052809291908181526020018280546104cd906113fa565b801561051a5780601f106104ef5761010080835404028352916020019161051a565b820191906000526020600020905b8154815290600101906020018083116104fd57829003601f168201915b505050505081565b81600180610530834361026e565b6002811115610541576105416110b3565b146105835760405162461bcd60e51b81526020600482015260126024820152714e6f7420617420676976656e20737461746560701b60448201526064016102d1565b438310156105de5760405162461bcd60e51b815260206004820152602260248201527f43616e27742064657072656361746520636f6e74726163742066726f6d2070616044820152611cdd60f21b60648201526084016102d1565b600080856040516105ef91906113b8565b90815260405190819003602001812060028101869055915084907f058d0e8a83e7a5228ca32808c88378cac5ff2eb4b3ef45a109be564227420c0590610636908890611197565b60405180910390a25050505050565b81600080610653834361026e565b6002811115610664576106646110b3565b146106a65760405162461bcd60e51b81526020600482015260126024820152714e6f7420617420676976656e20737461746560701b60448201526064016102d1565b438310156107005760405162461bcd60e51b815260206004820152602160248201527f43616e277420616374697661746520636f6e74726163742066726f6d207061736044820152601d60fa1b60648201526084016102d1565b6000808560405161071191906113b8565b90815260405190819003602001812060018101869055915084907f46ca43c7152dee1b324184afa3e3edec374f29a3880b6189ec2106c94c3f4b5d90610636908890611197565b6002546001600160a01b031633146107a55760405162461bcd60e51b815260206004820152601060248201526f4e6f74206120676f7665726e616e636560801b60448201526064016102d1565b8260006001600160a01b03166000826040516107c191906113b8565b908152604051908190036020019020546001600160a01b03161461081c5760405162461bcd60e51b8152602060048201526012602482015271105b1c9958591e481c9959da5cdd195c995960721b60448201526064016102d1565b8360008160405160200161083091906113b8565b604051602081830303815290604052905080516000036108815760405162461bcd60e51b815260206004820152600c60248201526b456d70747920737472696e6760a01b60448201526064016102d1565b846001600160a01b0381166108a85760405162461bcd60e51b81526004016102d1906113d4565b600080886040516108b991906113b8565b90815260405190819003602001902080546001600160a01b0389166001600160a01b0319909116178155905085156108fc576108f643600161144a565b60018201555b6001805480820182556000919091527fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf60161093789826114b2565b50866001600160a01b03167f50f74ca45caac8020b8d891bd13ea5a2d79564986ee6a839f0d914896388322d896040516109719190611197565b60405180910390a25050505050505050565b6002546001600160a01b031633146109d05760405162461bcd60e51b815260206004820152601060248201526f4e6f74206120676f7665726e616e636560801b60448201526064016102d1565b806001600160a01b0381166109f75760405162461bcd60e51b81526004016102d1906113d4565b600280546001600160a01b0319166001600160a01b0384169081179091556040517f8d55d160c0009eb3d739442df0a3ca089ed64378bfac017e7ddad463f9815b8790600090a25050565b60606001805480602002602001604051908101604052809291908181526020016000905b82821015610b12578382906000526020600020018054610a85906113fa565b80601f0160208091040260200160405190810160405280929190818152602001828054610ab1906113fa565b8015610afe5780601f10610ad357610100808354040283529160200191610afe565b820191906000526020600020905b815481529060010190602001808311610ae157829003601f168201915b505050505081526020019060010190610a66565b50505050905090565b806001600160a01b038116610b425760405162461bcd60e51b81526004016102d1906113d4565b600280546001600160a01b0319166001600160a01b0384169081179091556040517f15755ba3b1deec9fcf75c3db96f2d79482964aea97c6da27ceca7411850e405490600090a25050565b610b978382610522565b610ba18282610645565b807fe56cea40c3f094f642625199671fa50eeee44fd97d16ce800adddce75019e9568484604051610bd3929190611572565b60405180910390a2505050565b6060806000805b600154811015610ce857846002811115610c0357610c036110b3565b610cb260018381548110610c1957610c196115a0565b906000526020600020018054610c2e906113fa565b80601f0160208091040260200160405190810160405280929190818152602001828054610c5a906113fa565b8015610ca75780601f10610c7c57610100808354040283529160200191610ca7565b820191906000526020600020905b815481529060010190602001808311610c8a57829003601f168201915b50505050508861026e565b6002811115610cc357610cc36110b3565b03610cd65781610cd2816115b6565b9250505b80610ce0816115b6565b915050610be7565b5060008167ffffffffffffffff811115610d0457610d04610fcb565b604051908082528060200260200182016040528015610d3757816020015b6060815260200190600190039081610d225790505b50905060008267ffffffffffffffff811115610d5557610d55610fcb565b604051908082528060200260200182016040528015610d7e578160200160208202803683370190505b5090506000805b600154811015610fbd57876002811115610da157610da16110b3565b610e5060018381548110610db757610db76115a0565b906000526020600020018054610dcc906113fa565b80601f0160208091040260200160405190810160405280929190818152602001828054610df8906113fa565b8015610e455780601f10610e1a57610100808354040283529160200191610e45565b820191906000526020600020905b815481529060010190602001808311610e2857829003601f168201915b50505050508b61026e565b6002811115610e6157610e616110b3565b03610fab5760018181548110610e7957610e796115a0565b906000526020600020018054610e8e906113fa565b80601f0160208091040260200160405190810160405280929190818152602001828054610eba906113fa565b8015610f075780601f10610edc57610100808354040283529160200191610f07565b820191906000526020600020905b815481529060010190602001808311610eea57829003601f168201915b5050505050848381518110610f1e57610f1e6115a0565b6020026020010181905250600060018281548110610f3e57610f3e6115a0565b90600052602060002001604051610f5591906115cf565b9081526040519081900360200190205483516001600160a01b0390911690849084908110610f8557610f856115a0565b6001600160a01b039092166020928302919091019091015281610fa7816115b6565b9250505b80610fb5816115b6565b915050610d85565b509197909650945050505050565b634e487b7160e01b600052604160045260246000fd5b600082601f830112610ff257600080fd5b813567ffffffffffffffff8082111561100d5761100d610fcb565b604051601f8301601f19908116603f0116810190828211818310171561103557611035610fcb565b8160405283815286602085880101111561104e57600080fd5b836020870160208301376000602085830101528094505050505092915050565b6000806040838503121561108157600080fd5b823567ffffffffffffffff81111561109857600080fd5b6110a485828601610fe1565b95602094909401359450505050565b634e487b7160e01b600052602160045260246000fd5b60208101600383106110eb57634e487b7160e01b600052602160045260246000fd5b91905290565b60006020828403121561110357600080fd5b813567ffffffffffffffff81111561111a57600080fd5b61112684828501610fe1565b949350505050565b60006020828403121561114057600080fd5b5035919050565b60005b8381101561116257818101518382015260200161114a565b50506000910152565b60008151808452611183816020860160208601611147565b601f01601f19169290920160200192915050565b6020815260006111aa602083018461116b565b9392505050565b80356001600160a01b03811681146111c857600080fd5b919050565b6000806000606084860312156111e257600080fd5b833567ffffffffffffffff8111156111f957600080fd5b61120586828701610fe1565b935050611214602085016111b1565b91506040840135801515811461122957600080fd5b809150509250925092565b60006020828403121561124657600080fd5b6111aa826111b1565b600081518084526020808501808196508360051b8101915082860160005b8581101561129757828403895261128584835161116b565b9885019893509084019060010161126d565b5091979650505050505050565b6020815260006111aa602083018461124f565b6000806000606084860312156112cc57600080fd5b833567ffffffffffffffff808211156112e457600080fd5b6112f087838801610fe1565b9450602086013591508082111561130657600080fd5b5061131386828701610fe1565b925050604084013590509250925092565b6000806040838503121561133757600080fd5b8235915060208301356003811061134d57600080fd5b809150509250929050565b60408152600061136b604083018561124f565b82810360208481019190915284518083528582019282019060005b818110156113ab5784516001600160a01b031683529383019391830191600101611386565b5090979650505050505050565b600082516113ca818460208701611147565b9190910192915050565b6020808252600c908201526b5a65726f206164647265737360a01b604082015260600190565b600181811c9082168061140e57607f821691505b60208210810361142e57634e487b7160e01b600052602260045260246000fd5b50919050565b634e487b7160e01b600052601160045260246000fd5b8082018082111561145d5761145d611434565b92915050565b601f8211156114ad57600081815260208120601f850160051c8101602086101561148a5750805b601f850160051c820191505b818110156114a957828155600101611496565b5050505b505050565b815167ffffffffffffffff8111156114cc576114cc610fcb565b6114e0816114da84546113fa565b84611463565b602080601f83116001811461151557600084156114fd5750858301515b600019600386901b1c1916600185901b1785556114a9565b600085815260208120601f198616915b8281101561154457888601518255948401946001909101908401611525565b50858210156115625787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b604081526000611585604083018561116b565b8281036020840152611597818561116b565b95945050505050565b634e487b7160e01b600052603260045260246000fd5b6000600182016115c8576115c8611434565b5060010190565b60008083546115dd816113fa565b600182811680156115f5576001811461160a57611639565b60ff1984168752821515830287019450611639565b8760005260208060002060005b858110156116305781548a820152908401908201611617565b50505082870194505b5092969550505050505056fea264697066735822122044479653d70716308dc1faa779fd3bb2a9aeb3562ed35a360cd10fac3eb0d66e64736f6c63430008120033"

const RegistryContractAddress = "0x0000000000000000000000000000000000000401"

// RegistryABI is the input ABI used to generate the binding from.
// Deprecated: Use RegistryMetaData.ABI instead.
var RegistryABI = RegistryMetaData.ABI

// Registry is an auto generated Go binding around an Ethereum contract.
type Registry struct {
	RegistryCaller     // Read-only binding to the contract
	RegistryTransactor // Write-only binding to the contract
	RegistryFilterer   // Log filterer for contract events
}

// RegistryCaller is an auto generated read-only Go binding around an Ethereum contract.
type RegistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RegistryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type RegistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RegistryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type RegistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RegistrySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type RegistrySession struct {
	Contract     *Registry         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// RegistryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type RegistryCallerSession struct {
	Contract *RegistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// RegistryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type RegistryTransactorSession struct {
	Contract     *RegistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// RegistryRaw is an auto generated low-level Go binding around an Ethereum contract.
type RegistryRaw struct {
	Contract *Registry // Generic contract binding to access the raw methods on
}

// RegistryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type RegistryCallerRaw struct {
	Contract *RegistryCaller // Generic read-only contract binding to access the raw methods on
}

// RegistryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
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

// ContractNames is a free data retrieval call binding the contract method 0x3ca6bb92.
//
// Solidity: function contractNames(uint256 ) view returns(string)
func (_Registry *RegistryCaller) ContractNames(opts *bind.CallOpts, arg0 *big.Int) (string, error) {
	ret0 := new(string)
	out := ret0
	err := _Registry.contract.Call(opts, out, "contractNames", arg0)
	return *ret0, err
}

// ContractNames is a free data retrieval call binding the contract method 0x3ca6bb92.
//
// Solidity: function contractNames(uint256 ) view returns(string)
func (_Registry *RegistrySession) ContractNames(arg0 *big.Int) (string, error) {
	return _Registry.Contract.ContractNames(&_Registry.CallOpts, arg0)
}

// ContractNames is a free data retrieval call binding the contract method 0x3ca6bb92.
//
// Solidity: function contractNames(uint256 ) view returns(string)
func (_Registry *RegistryCallerSession) ContractNames(arg0 *big.Int) (string, error) {
	return _Registry.Contract.ContractNames(&_Registry.CallOpts, arg0)
}

// GetAllContractNames is a free data retrieval call binding the contract method 0xd7a9b816.
//
// Solidity: function getAllContractNames() view returns(string[])
func (_Registry *RegistryCaller) GetAllContractNames(opts *bind.CallOpts) ([]string, error) {
	ret0 := new([]string)
	out := ret0
	err := _Registry.contract.Call(opts, out, "getAllContractNames")
	return *ret0, err
}

// GetAllContractNames is a free data retrieval call binding the contract method 0xd7a9b816.
//
// Solidity: function getAllContractNames() view returns(string[])
func (_Registry *RegistrySession) GetAllContractNames() ([]string, error) {
	return _Registry.Contract.GetAllContractNames(&_Registry.CallOpts)
}

// GetAllContractNames is a free data retrieval call binding the contract method 0xd7a9b816.
//
// Solidity: function getAllContractNames() view returns(string[])
func (_Registry *RegistryCallerSession) GetAllContractNames() ([]string, error) {
	return _Registry.Contract.GetAllContractNames(&_Registry.CallOpts)
}

// GetContractIfActive is a free data retrieval call binding the contract method 0x35b5cfbb.
//
// Solidity: function getContractIfActive(string name) view returns(address)
func (_Registry *RegistryCaller) GetContractIfActive(opts *bind.CallOpts, name string) (common.Address, error) {
	ret0 := new(common.Address)
	out := ret0
	err := _Registry.contract.Call(opts, out, "getContractIfActive", name)
	return *ret0, err
}

// GetContractIfActive is a free data retrieval call binding the contract method 0x35b5cfbb.
//
// Solidity: function getContractIfActive(string name) view returns(address)
func (_Registry *RegistrySession) GetContractIfActive(name string) (common.Address, error) {
	return _Registry.Contract.GetContractIfActive(&_Registry.CallOpts, name)
}

// GetContractIfActive is a free data retrieval call binding the contract method 0x35b5cfbb.
//
// Solidity: function getContractIfActive(string name) view returns(address)
func (_Registry *RegistryCallerSession) GetContractIfActive(name string) (common.Address, error) {
	return _Registry.Contract.GetContractIfActive(&_Registry.CallOpts, name)
}

// Governance is a free data retrieval call binding the contract method 0x5aa6e675.
//
// Solidity: function governance() view returns(address)
func (_Registry *RegistryCaller) Governance(opts *bind.CallOpts) (common.Address, error) {
	ret0 := new(common.Address)
	out := ret0
	err := _Registry.contract.Call(opts, out, "governance")
	return *ret0, err
}

// Governance is a free data retrieval call binding the contract method 0x5aa6e675.
//
// Solidity: function governance() view returns(address)
func (_Registry *RegistrySession) Governance() (common.Address, error) {
	return _Registry.Contract.Governance(&_Registry.CallOpts)
}

// Governance is a free data retrieval call binding the contract method 0x5aa6e675.
//
// Solidity: function governance() view returns(address)
func (_Registry *RegistryCallerSession) Governance() (common.Address, error) {
	return _Registry.Contract.Governance(&_Registry.CallOpts)
}

// ReadAllContractsAtGivenState is a free data retrieval call binding the contract method 0xfa96b37f.
//
// Solidity: function readAllContractsAtGivenState(uint256 blockNumber, uint8 state) view returns(string[], address[])
func (_Registry *RegistryCaller) ReadAllContractsAtGivenState(opts *bind.CallOpts, blockNumber *big.Int, state uint8) ([]string, []common.Address, error) {
	var (
		ret0 = new([]string)
		ret1 = new([]common.Address)
	)
	out := &[]interface{}{
		ret0,
		ret1,
	}
	err := _Registry.contract.Call(opts, out, "readAllContractsAtGivenState", blockNumber, state)

	return *ret0, *ret1, err
}

// ReadAllContractsAtGivenState is a free data retrieval call binding the contract method 0xfa96b37f.
//
// Solidity: function readAllContractsAtGivenState(uint256 blockNumber, uint8 state) view returns(string[], address[])
func (_Registry *RegistrySession) ReadAllContractsAtGivenState(blockNumber *big.Int, state uint8) ([]string, []common.Address, error) {
	return _Registry.Contract.ReadAllContractsAtGivenState(&_Registry.CallOpts, blockNumber, state)
}

// ReadAllContractsAtGivenState is a free data retrieval call binding the contract method 0xfa96b37f.
//
// Solidity: function readAllContractsAtGivenState(uint256 blockNumber, uint8 state) view returns(string[], address[])
func (_Registry *RegistryCallerSession) ReadAllContractsAtGivenState(blockNumber *big.Int, state uint8) ([]string, []common.Address, error) {
	return _Registry.Contract.ReadAllContractsAtGivenState(&_Registry.CallOpts, blockNumber, state)
}

// Registry is a free data retrieval call binding the contract method 0x92a296c9.
//
// Solidity: function registry(string ) view returns(address addr, uint256 activationBlockNumber, uint256 deprecateBlockNumber)
func (_Registry *RegistryCaller) Registry(opts *bind.CallOpts, arg0 string) (struct {
	Addr                  common.Address
	ActivationBlockNumber *big.Int
	DeprecateBlockNumber  *big.Int
}, error,
) {
	ret := new(struct {
		Addr                  common.Address
		ActivationBlockNumber *big.Int
		DeprecateBlockNumber  *big.Int
	})
	out := ret
	err := _Registry.contract.Call(opts, out, "registry", arg0)
	return *ret, err
}

// Registry is a free data retrieval call binding the contract method 0x92a296c9.
//
// Solidity: function registry(string ) view returns(address addr, uint256 activationBlockNumber, uint256 deprecateBlockNumber)
func (_Registry *RegistrySession) Registry(arg0 string) (struct {
	Addr                  common.Address
	ActivationBlockNumber *big.Int
	DeprecateBlockNumber  *big.Int
}, error) {
	return _Registry.Contract.Registry(&_Registry.CallOpts, arg0)
}

// Registry is a free data retrieval call binding the contract method 0x92a296c9.
//
// Solidity: function registry(string ) view returns(address addr, uint256 activationBlockNumber, uint256 deprecateBlockNumber)
func (_Registry *RegistryCallerSession) Registry(arg0 string) (struct {
	Addr                  common.Address
	ActivationBlockNumber *big.Int
	DeprecateBlockNumber  *big.Int
}, error) {
	return _Registry.Contract.Registry(&_Registry.CallOpts, arg0)
}

// StateAt is a free data retrieval call binding the contract method 0x05e4f16e.
//
// Solidity: function stateAt(string name, uint256 blockNumber) view returns(uint8)
func (_Registry *RegistryCaller) StateAt(opts *bind.CallOpts, name string, blockNumber *big.Int) (uint8, error) {
	ret0 := new(uint8)
	out := ret0
	err := _Registry.contract.Call(opts, out, "stateAt", name, blockNumber)
	return *ret0, err
}

// StateAt is a free data retrieval call binding the contract method 0x05e4f16e.
//
// Solidity: function stateAt(string name, uint256 blockNumber) view returns(uint8)
func (_Registry *RegistrySession) StateAt(name string, blockNumber *big.Int) (uint8, error) {
	return _Registry.Contract.StateAt(&_Registry.CallOpts, name, blockNumber)
}

// StateAt is a free data retrieval call binding the contract method 0x05e4f16e.
//
// Solidity: function stateAt(string name, uint256 blockNumber) view returns(uint8)
func (_Registry *RegistryCallerSession) StateAt(name string, blockNumber *big.Int) (uint8, error) {
	return _Registry.Contract.StateAt(&_Registry.CallOpts, name, blockNumber)
}

// Activate is a paid mutator transaction binding the contract method 0x51afb6a3.
//
// Solidity: function activate(string name, uint256 activationBlockNumber) returns()
func (_Registry *RegistryTransactor) Activate(opts *bind.TransactOpts, name string, activationBlockNumber *big.Int) (*types.Transaction, error) {
	return _Registry.contract.Transact(opts, "activate", name, activationBlockNumber)
}

// Activate is a paid mutator transaction binding the contract method 0x51afb6a3.
//
// Solidity: function activate(string name, uint256 activationBlockNumber) returns()
func (_Registry *RegistrySession) Activate(name string, activationBlockNumber *big.Int) (*types.Transaction, error) {
	return _Registry.Contract.Activate(&_Registry.TransactOpts, name, activationBlockNumber)
}

// Activate is a paid mutator transaction binding the contract method 0x51afb6a3.
//
// Solidity: function activate(string name, uint256 activationBlockNumber) returns()
func (_Registry *RegistryTransactorSession) Activate(name string, activationBlockNumber *big.Int) (*types.Transaction, error) {
	return _Registry.Contract.Activate(&_Registry.TransactOpts, name, activationBlockNumber)
}

// ConstructContract is a paid mutator transaction binding the contract method 0xe14c672d.
//
// Solidity: function constructContract(address _governance) returns()
func (_Registry *RegistryTransactor) ConstructContract(opts *bind.TransactOpts, _governance common.Address) (*types.Transaction, error) {
	return _Registry.contract.Transact(opts, "constructContract", _governance)
}

// ConstructContract is a paid mutator transaction binding the contract method 0xe14c672d.
//
// Solidity: function constructContract(address _governance) returns()
func (_Registry *RegistrySession) ConstructContract(_governance common.Address) (*types.Transaction, error) {
	return _Registry.Contract.ConstructContract(&_Registry.TransactOpts, _governance)
}

// ConstructContract is a paid mutator transaction binding the contract method 0xe14c672d.
//
// Solidity: function constructContract(address _governance) returns()
func (_Registry *RegistryTransactorSession) ConstructContract(_governance common.Address) (*types.Transaction, error) {
	return _Registry.Contract.ConstructContract(&_Registry.TransactOpts, _governance)
}

// Deprecate is a paid mutator transaction binding the contract method 0x4d0e6606.
//
// Solidity: function deprecate(string name, uint256 deprecateBlockNumber) returns()
func (_Registry *RegistryTransactor) Deprecate(opts *bind.TransactOpts, name string, deprecateBlockNumber *big.Int) (*types.Transaction, error) {
	return _Registry.contract.Transact(opts, "deprecate", name, deprecateBlockNumber)
}

// Deprecate is a paid mutator transaction binding the contract method 0x4d0e6606.
//
// Solidity: function deprecate(string name, uint256 deprecateBlockNumber) returns()
func (_Registry *RegistrySession) Deprecate(name string, deprecateBlockNumber *big.Int) (*types.Transaction, error) {
	return _Registry.Contract.Deprecate(&_Registry.TransactOpts, name, deprecateBlockNumber)
}

// Deprecate is a paid mutator transaction binding the contract method 0x4d0e6606.
//
// Solidity: function deprecate(string name, uint256 deprecateBlockNumber) returns()
func (_Registry *RegistryTransactorSession) Deprecate(name string, deprecateBlockNumber *big.Int) (*types.Transaction, error) {
	return _Registry.Contract.Deprecate(&_Registry.TransactOpts, name, deprecateBlockNumber)
}

// Register is a paid mutator transaction binding the contract method 0x60d7a278.
//
// Solidity: function register(string name, address addr, bool activation) returns()
func (_Registry *RegistryTransactor) Register(opts *bind.TransactOpts, name string, addr common.Address, activation bool) (*types.Transaction, error) {
	return _Registry.contract.Transact(opts, "register", name, addr, activation)
}

// Register is a paid mutator transaction binding the contract method 0x60d7a278.
//
// Solidity: function register(string name, address addr, bool activation) returns()
func (_Registry *RegistrySession) Register(name string, addr common.Address, activation bool) (*types.Transaction, error) {
	return _Registry.Contract.Register(&_Registry.TransactOpts, name, addr, activation)
}

// Register is a paid mutator transaction binding the contract method 0x60d7a278.
//
// Solidity: function register(string name, address addr, bool activation) returns()
func (_Registry *RegistryTransactorSession) Register(name string, addr common.Address, activation bool) (*types.Transaction, error) {
	return _Registry.Contract.Register(&_Registry.TransactOpts, name, addr, activation)
}

// Replace is a paid mutator transaction binding the contract method 0xe47f1160.
//
// Solidity: function replace(string prevName, string newName, uint256 replaceBlockNumber) returns()
func (_Registry *RegistryTransactor) Replace(opts *bind.TransactOpts, prevName string, newName string, replaceBlockNumber *big.Int) (*types.Transaction, error) {
	return _Registry.contract.Transact(opts, "replace", prevName, newName, replaceBlockNumber)
}

// Replace is a paid mutator transaction binding the contract method 0xe47f1160.
//
// Solidity: function replace(string prevName, string newName, uint256 replaceBlockNumber) returns()
func (_Registry *RegistrySession) Replace(prevName string, newName string, replaceBlockNumber *big.Int) (*types.Transaction, error) {
	return _Registry.Contract.Replace(&_Registry.TransactOpts, prevName, newName, replaceBlockNumber)
}

// Replace is a paid mutator transaction binding the contract method 0xe47f1160.
//
// Solidity: function replace(string prevName, string newName, uint256 replaceBlockNumber) returns()
func (_Registry *RegistryTransactorSession) Replace(prevName string, newName string, replaceBlockNumber *big.Int) (*types.Transaction, error) {
	return _Registry.Contract.Replace(&_Registry.TransactOpts, prevName, newName, replaceBlockNumber)
}

// UpdateGovernance is a paid mutator transaction binding the contract method 0xb2561263.
//
// Solidity: function updateGovernance(address _newGovernance) returns()
func (_Registry *RegistryTransactor) UpdateGovernance(opts *bind.TransactOpts, _newGovernance common.Address) (*types.Transaction, error) {
	return _Registry.contract.Transact(opts, "updateGovernance", _newGovernance)
}

// UpdateGovernance is a paid mutator transaction binding the contract method 0xb2561263.
//
// Solidity: function updateGovernance(address _newGovernance) returns()
func (_Registry *RegistrySession) UpdateGovernance(_newGovernance common.Address) (*types.Transaction, error) {
	return _Registry.Contract.UpdateGovernance(&_Registry.TransactOpts, _newGovernance)
}

// UpdateGovernance is a paid mutator transaction binding the contract method 0xb2561263.
//
// Solidity: function updateGovernance(address _newGovernance) returns()
func (_Registry *RegistryTransactorSession) UpdateGovernance(_newGovernance common.Address) (*types.Transaction, error) {
	return _Registry.Contract.UpdateGovernance(&_Registry.TransactOpts, _newGovernance)
}

// RegistryActivatedIterator is returned from FilterActivated and is used to iterate over the raw logs and unpacked data for Activated events raised by the Registry contract.
type RegistryActivatedIterator struct {
	Event *RegistryActivated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  klaytn.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *RegistryActivatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RegistryActivated)
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
		it.Event = new(RegistryActivated)
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
func (it *RegistryActivatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RegistryActivatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RegistryActivated represents a Activated event raised by the Registry contract.
type RegistryActivated struct {
	Name                  string
	ActivationBlockNumber *big.Int
	Raw                   types.Log // Blockchain specific contextual infos
}

// FilterActivated is a free log retrieval operation binding the contract event 0x46ca43c7152dee1b324184afa3e3edec374f29a3880b6189ec2106c94c3f4b5d.
//
// Solidity: event Activated(string name, uint256 indexed activationBlockNumber)
func (_Registry *RegistryFilterer) FilterActivated(opts *bind.FilterOpts, activationBlockNumber []*big.Int) (*RegistryActivatedIterator, error) {

	var activationBlockNumberRule []interface{}
	for _, activationBlockNumberItem := range activationBlockNumber {
		activationBlockNumberRule = append(activationBlockNumberRule, activationBlockNumberItem)
	}

	logs, sub, err := _Registry.contract.FilterLogs(opts, "Activated", activationBlockNumberRule)
	if err != nil {
		return nil, err
	}
	return &RegistryActivatedIterator{contract: _Registry.contract, event: "Activated", logs: logs, sub: sub}, nil
}

// WatchActivated is a free log subscription operation binding the contract event 0x46ca43c7152dee1b324184afa3e3edec374f29a3880b6189ec2106c94c3f4b5d.
//
// Solidity: event Activated(string name, uint256 indexed activationBlockNumber)
func (_Registry *RegistryFilterer) WatchActivated(opts *bind.WatchOpts, sink chan<- *RegistryActivated, activationBlockNumber []*big.Int) (event.Subscription, error) {

	var activationBlockNumberRule []interface{}
	for _, activationBlockNumberItem := range activationBlockNumber {
		activationBlockNumberRule = append(activationBlockNumberRule, activationBlockNumberItem)
	}

	logs, sub, err := _Registry.contract.WatchLogs(opts, "Activated", activationBlockNumberRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RegistryActivated)
				if err := _Registry.contract.UnpackLog(event, "Activated", log); err != nil {
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

// ParseActivated is a log parse operation binding the contract event 0x46ca43c7152dee1b324184afa3e3edec374f29a3880b6189ec2106c94c3f4b5d.
//
// Solidity: event Activated(string name, uint256 indexed activationBlockNumber)
func (_Registry *RegistryFilterer) ParseActivated(log types.Log) (*RegistryActivated, error) {
	event := new(RegistryActivated)
	if err := _Registry.contract.UnpackLog(event, "Activated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RegistryConstructContractIterator is returned from FilterConstructContract and is used to iterate over the raw logs and unpacked data for ConstructContract events raised by the Registry contract.
type RegistryConstructContractIterator struct {
	Event *RegistryConstructContract // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  klaytn.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *RegistryConstructContractIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RegistryConstructContract)
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
		it.Event = new(RegistryConstructContract)
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
func (it *RegistryConstructContractIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RegistryConstructContractIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RegistryConstructContract represents a ConstructContract event raised by the Registry contract.
type RegistryConstructContract struct {
	Governance common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterConstructContract is a free log retrieval operation binding the contract event 0x15755ba3b1deec9fcf75c3db96f2d79482964aea97c6da27ceca7411850e4054.
//
// Solidity: event ConstructContract(address indexed governance)
func (_Registry *RegistryFilterer) FilterConstructContract(opts *bind.FilterOpts, governance []common.Address) (*RegistryConstructContractIterator, error) {

	var governanceRule []interface{}
	for _, governanceItem := range governance {
		governanceRule = append(governanceRule, governanceItem)
	}

	logs, sub, err := _Registry.contract.FilterLogs(opts, "ConstructContract", governanceRule)
	if err != nil {
		return nil, err
	}
	return &RegistryConstructContractIterator{contract: _Registry.contract, event: "ConstructContract", logs: logs, sub: sub}, nil
}

// WatchConstructContract is a free log subscription operation binding the contract event 0x15755ba3b1deec9fcf75c3db96f2d79482964aea97c6da27ceca7411850e4054.
//
// Solidity: event ConstructContract(address indexed governance)
func (_Registry *RegistryFilterer) WatchConstructContract(opts *bind.WatchOpts, sink chan<- *RegistryConstructContract, governance []common.Address) (event.Subscription, error) {

	var governanceRule []interface{}
	for _, governanceItem := range governance {
		governanceRule = append(governanceRule, governanceItem)
	}

	logs, sub, err := _Registry.contract.WatchLogs(opts, "ConstructContract", governanceRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RegistryConstructContract)
				if err := _Registry.contract.UnpackLog(event, "ConstructContract", log); err != nil {
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

// ParseConstructContract is a log parse operation binding the contract event 0x15755ba3b1deec9fcf75c3db96f2d79482964aea97c6da27ceca7411850e4054.
//
// Solidity: event ConstructContract(address indexed governance)
func (_Registry *RegistryFilterer) ParseConstructContract(log types.Log) (*RegistryConstructContract, error) {
	event := new(RegistryConstructContract)
	if err := _Registry.contract.UnpackLog(event, "ConstructContract", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RegistryDeprecatedIterator is returned from FilterDeprecated and is used to iterate over the raw logs and unpacked data for Deprecated events raised by the Registry contract.
type RegistryDeprecatedIterator struct {
	Event *RegistryDeprecated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  klaytn.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *RegistryDeprecatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RegistryDeprecated)
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
		it.Event = new(RegistryDeprecated)
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
func (it *RegistryDeprecatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RegistryDeprecatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RegistryDeprecated represents a Deprecated event raised by the Registry contract.
type RegistryDeprecated struct {
	Name                 string
	DeprecateBlockNumber *big.Int
	Raw                  types.Log // Blockchain specific contextual infos
}

// FilterDeprecated is a free log retrieval operation binding the contract event 0x058d0e8a83e7a5228ca32808c88378cac5ff2eb4b3ef45a109be564227420c05.
//
// Solidity: event Deprecated(string name, uint256 indexed deprecateBlockNumber)
func (_Registry *RegistryFilterer) FilterDeprecated(opts *bind.FilterOpts, deprecateBlockNumber []*big.Int) (*RegistryDeprecatedIterator, error) {

	var deprecateBlockNumberRule []interface{}
	for _, deprecateBlockNumberItem := range deprecateBlockNumber {
		deprecateBlockNumberRule = append(deprecateBlockNumberRule, deprecateBlockNumberItem)
	}

	logs, sub, err := _Registry.contract.FilterLogs(opts, "Deprecated", deprecateBlockNumberRule)
	if err != nil {
		return nil, err
	}
	return &RegistryDeprecatedIterator{contract: _Registry.contract, event: "Deprecated", logs: logs, sub: sub}, nil
}

// WatchDeprecated is a free log subscription operation binding the contract event 0x058d0e8a83e7a5228ca32808c88378cac5ff2eb4b3ef45a109be564227420c05.
//
// Solidity: event Deprecated(string name, uint256 indexed deprecateBlockNumber)
func (_Registry *RegistryFilterer) WatchDeprecated(opts *bind.WatchOpts, sink chan<- *RegistryDeprecated, deprecateBlockNumber []*big.Int) (event.Subscription, error) {

	var deprecateBlockNumberRule []interface{}
	for _, deprecateBlockNumberItem := range deprecateBlockNumber {
		deprecateBlockNumberRule = append(deprecateBlockNumberRule, deprecateBlockNumberItem)
	}

	logs, sub, err := _Registry.contract.WatchLogs(opts, "Deprecated", deprecateBlockNumberRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RegistryDeprecated)
				if err := _Registry.contract.UnpackLog(event, "Deprecated", log); err != nil {
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

// ParseDeprecated is a log parse operation binding the contract event 0x058d0e8a83e7a5228ca32808c88378cac5ff2eb4b3ef45a109be564227420c05.
//
// Solidity: event Deprecated(string name, uint256 indexed deprecateBlockNumber)
func (_Registry *RegistryFilterer) ParseDeprecated(log types.Log) (*RegistryDeprecated, error) {
	event := new(RegistryDeprecated)
	if err := _Registry.contract.UnpackLog(event, "Deprecated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RegistryRegisteredIterator is returned from FilterRegistered and is used to iterate over the raw logs and unpacked data for Registered events raised by the Registry contract.
type RegistryRegisteredIterator struct {
	Event *RegistryRegistered // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  klaytn.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
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
	Name string
	Addr common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterRegistered is a free log retrieval operation binding the contract event 0x50f74ca45caac8020b8d891bd13ea5a2d79564986ee6a839f0d914896388322d.
//
// Solidity: event Registered(string name, address indexed addr)
func (_Registry *RegistryFilterer) FilterRegistered(opts *bind.FilterOpts, addr []common.Address) (*RegistryRegisteredIterator, error) {

	var addrRule []interface{}
	for _, addrItem := range addr {
		addrRule = append(addrRule, addrItem)
	}

	logs, sub, err := _Registry.contract.FilterLogs(opts, "Registered", addrRule)
	if err != nil {
		return nil, err
	}
	return &RegistryRegisteredIterator{contract: _Registry.contract, event: "Registered", logs: logs, sub: sub}, nil
}

// WatchRegistered is a free log subscription operation binding the contract event 0x50f74ca45caac8020b8d891bd13ea5a2d79564986ee6a839f0d914896388322d.
//
// Solidity: event Registered(string name, address indexed addr)
func (_Registry *RegistryFilterer) WatchRegistered(opts *bind.WatchOpts, sink chan<- *RegistryRegistered, addr []common.Address) (event.Subscription, error) {

	var addrRule []interface{}
	for _, addrItem := range addr {
		addrRule = append(addrRule, addrItem)
	}

	logs, sub, err := _Registry.contract.WatchLogs(opts, "Registered", addrRule)
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

// ParseRegistered is a log parse operation binding the contract event 0x50f74ca45caac8020b8d891bd13ea5a2d79564986ee6a839f0d914896388322d.
//
// Solidity: event Registered(string name, address indexed addr)
func (_Registry *RegistryFilterer) ParseRegistered(log types.Log) (*RegistryRegistered, error) {
	event := new(RegistryRegistered)
	if err := _Registry.contract.UnpackLog(event, "Registered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RegistryReplacedIterator is returned from FilterReplaced and is used to iterate over the raw logs and unpacked data for Replaced events raised by the Registry contract.
type RegistryReplacedIterator struct {
	Event *RegistryReplaced // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  klaytn.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *RegistryReplacedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RegistryReplaced)
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
		it.Event = new(RegistryReplaced)
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
func (it *RegistryReplacedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RegistryReplacedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RegistryReplaced represents a Replaced event raised by the Registry contract.
type RegistryReplaced struct {
	PrevName           string
	NewName            string
	ReplaceBlockNumber *big.Int
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterReplaced is a free log retrieval operation binding the contract event 0xe56cea40c3f094f642625199671fa50eeee44fd97d16ce800adddce75019e956.
//
// Solidity: event Replaced(string prevName, string newName, uint256 indexed replaceBlockNumber)
func (_Registry *RegistryFilterer) FilterReplaced(opts *bind.FilterOpts, replaceBlockNumber []*big.Int) (*RegistryReplacedIterator, error) {

	var replaceBlockNumberRule []interface{}
	for _, replaceBlockNumberItem := range replaceBlockNumber {
		replaceBlockNumberRule = append(replaceBlockNumberRule, replaceBlockNumberItem)
	}

	logs, sub, err := _Registry.contract.FilterLogs(opts, "Replaced", replaceBlockNumberRule)
	if err != nil {
		return nil, err
	}
	return &RegistryReplacedIterator{contract: _Registry.contract, event: "Replaced", logs: logs, sub: sub}, nil
}

// WatchReplaced is a free log subscription operation binding the contract event 0xe56cea40c3f094f642625199671fa50eeee44fd97d16ce800adddce75019e956.
//
// Solidity: event Replaced(string prevName, string newName, uint256 indexed replaceBlockNumber)
func (_Registry *RegistryFilterer) WatchReplaced(opts *bind.WatchOpts, sink chan<- *RegistryReplaced, replaceBlockNumber []*big.Int) (event.Subscription, error) {

	var replaceBlockNumberRule []interface{}
	for _, replaceBlockNumberItem := range replaceBlockNumber {
		replaceBlockNumberRule = append(replaceBlockNumberRule, replaceBlockNumberItem)
	}

	logs, sub, err := _Registry.contract.WatchLogs(opts, "Replaced", replaceBlockNumberRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RegistryReplaced)
				if err := _Registry.contract.UnpackLog(event, "Replaced", log); err != nil {
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

// ParseReplaced is a log parse operation binding the contract event 0xe56cea40c3f094f642625199671fa50eeee44fd97d16ce800adddce75019e956.
//
// Solidity: event Replaced(string prevName, string newName, uint256 indexed replaceBlockNumber)
func (_Registry *RegistryFilterer) ParseReplaced(log types.Log) (*RegistryReplaced, error) {
	event := new(RegistryReplaced)
	if err := _Registry.contract.UnpackLog(event, "Replaced", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RegistryUpdateGovernanceIterator is returned from FilterUpdateGovernance and is used to iterate over the raw logs and unpacked data for UpdateGovernance events raised by the Registry contract.
type RegistryUpdateGovernanceIterator struct {
	Event *RegistryUpdateGovernance // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  klaytn.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *RegistryUpdateGovernanceIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RegistryUpdateGovernance)
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
		it.Event = new(RegistryUpdateGovernance)
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
func (it *RegistryUpdateGovernanceIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RegistryUpdateGovernanceIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RegistryUpdateGovernance represents a UpdateGovernance event raised by the Registry contract.
type RegistryUpdateGovernance struct {
	NewGovernance common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterUpdateGovernance is a free log retrieval operation binding the contract event 0x8d55d160c0009eb3d739442df0a3ca089ed64378bfac017e7ddad463f9815b87.
//
// Solidity: event UpdateGovernance(address indexed newGovernance)
func (_Registry *RegistryFilterer) FilterUpdateGovernance(opts *bind.FilterOpts, newGovernance []common.Address) (*RegistryUpdateGovernanceIterator, error) {

	var newGovernanceRule []interface{}
	for _, newGovernanceItem := range newGovernance {
		newGovernanceRule = append(newGovernanceRule, newGovernanceItem)
	}

	logs, sub, err := _Registry.contract.FilterLogs(opts, "UpdateGovernance", newGovernanceRule)
	if err != nil {
		return nil, err
	}
	return &RegistryUpdateGovernanceIterator{contract: _Registry.contract, event: "UpdateGovernance", logs: logs, sub: sub}, nil
}

// WatchUpdateGovernance is a free log subscription operation binding the contract event 0x8d55d160c0009eb3d739442df0a3ca089ed64378bfac017e7ddad463f9815b87.
//
// Solidity: event UpdateGovernance(address indexed newGovernance)
func (_Registry *RegistryFilterer) WatchUpdateGovernance(opts *bind.WatchOpts, sink chan<- *RegistryUpdateGovernance, newGovernance []common.Address) (event.Subscription, error) {

	var newGovernanceRule []interface{}
	for _, newGovernanceItem := range newGovernance {
		newGovernanceRule = append(newGovernanceRule, newGovernanceItem)
	}

	logs, sub, err := _Registry.contract.WatchLogs(opts, "UpdateGovernance", newGovernanceRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RegistryUpdateGovernance)
				if err := _Registry.contract.UnpackLog(event, "UpdateGovernance", log); err != nil {
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

// ParseUpdateGovernance is a log parse operation binding the contract event 0x8d55d160c0009eb3d739442df0a3ca089ed64378bfac017e7ddad463f9815b87.
//
// Solidity: event UpdateGovernance(address indexed newGovernance)
func (_Registry *RegistryFilterer) ParseUpdateGovernance(log types.Log) (*RegistryUpdateGovernance, error) {
	event := new(RegistryUpdateGovernance)
	if err := _Registry.contract.UnpackLog(event, "UpdateGovernance", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
