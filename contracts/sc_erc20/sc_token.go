// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package sctoken

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

// AddressABI is the input ABI used to generate the binding from.
const AddressABI = "[]"

// AddressBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const AddressBinRuntime = `73000000000000000000000000000000000000000030146080604052600080fdfea165627a7a72305820b7534b8b203852d9ead56352219b5abb5424533d9bb1b9cf84c0b1d15d11fee50029`

// AddressBin is the compiled bytecode used for deploying new contracts.
var AddressBin = "0x604c6023600b82828239805160001a607314601657fe5b30600052607381538281f3fe73000000000000000000000000000000000000000030146080604052600080fdfea165627a7a72305820b7534b8b203852d9ead56352219b5abb5424533d9bb1b9cf84c0b1d15d11fee50029"

// DeployAddress deploys a new Klaytn contract, binding an instance of Address to it.
func DeployAddress(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Address, error) {
	parsed, err := abi.JSON(strings.NewReader(AddressABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(AddressBin), backend)
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
	parsed, err := abi.JSON(strings.NewReader(AddressABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Address *AddressRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
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
func (_Address *AddressCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
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

// ERC20ABI is the input ABI used to generate the binding from.
const ERC20ABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"spender\",\"type\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"sender\",\"type\":\"address\"},{\"name\":\"recipient\",\"type\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"spender\",\"type\":\"address\"},{\"name\":\"addedValue\",\"type\":\"uint256\"}],\"name\":\"increaseAllowance\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"spender\",\"type\":\"address\"},{\"name\":\"subtractedValue\",\"type\":\"uint256\"}],\"name\":\"decreaseAllowance\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"recipient\",\"type\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"owner\",\"type\":\"address\"},{\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"}]"

// ERC20BinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const ERC20BinRuntime = `608060405234801561001057600080fd5b50600436106100885760003560e01c806370a082311161005b57806370a0823114610149578063a457c2d71461016f578063a9059cbb1461019b578063dd62ed3e146101c757610088565b8063095ea7b31461008d57806318160ddd146100cd57806323b872dd146100e7578063395093511461011d575b600080fd5b6100b9600480360360408110156100a357600080fd5b506001600160a01b0381351690602001356101f5565b604080519115158252519081900360200190f35b6100d561020b565b60408051918252519081900360200190f35b6100b9600480360360608110156100fd57600080fd5b506001600160a01b03813581169160208101359091169060400135610211565b6100b96004803603604081101561013357600080fd5b506001600160a01b038135169060200135610268565b6100d56004803603602081101561015f57600080fd5b50356001600160a01b03166102a4565b6100b96004803603604081101561018557600080fd5b506001600160a01b0381351690602001356102bf565b6100b9600480360360408110156101b157600080fd5b506001600160a01b0381351690602001356102fb565b6100d5600480360360408110156101dd57600080fd5b506001600160a01b0381358116916020013516610308565b6000610202338484610333565b50600192915050565b60025490565b600061021e848484610425565b6001600160a01b03841660009081526001602090815260408083203380855292529091205461025e918691610259908663ffffffff61056d16565b610333565b5060019392505050565b3360008181526001602090815260408083206001600160a01b03871684529091528120549091610202918590610259908663ffffffff6105cd16565b6001600160a01b031660009081526020819052604090205490565b3360008181526001602090815260408083206001600160a01b03871684529091528120549091610202918590610259908663ffffffff61056d16565b6000610202338484610425565b6001600160a01b03918216600090815260016020908152604080832093909416825291909152205490565b6001600160a01b03831661037b57604051600160e51b62461bcd02815260040180806020018281038252602481526020018061069c6024913960400191505060405180910390fd5b6001600160a01b0382166103c357604051600160e51b62461bcd0281526004018080602001828103825260228152602001806106556022913960400191505060405180910390fd5b6001600160a01b03808416600081815260016020908152604080832094871680845294825291829020859055815185815291517f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b9259281900390910190a3505050565b6001600160a01b03831661046d57604051600160e51b62461bcd0281526004018080602001828103825260258152602001806106776025913960400191505060405180910390fd5b6001600160a01b0382166104b557604051600160e51b62461bcd0281526004018080602001828103825260238152602001806106326023913960400191505060405180910390fd5b6001600160a01b0383166000908152602081905260409020546104de908263ffffffff61056d16565b6001600160a01b038085166000908152602081905260408082209390935590841681522054610513908263ffffffff6105cd16565b6001600160a01b038084166000818152602081815260409182902094909455805185815290519193928716927fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef92918290030190a3505050565b6000828211156105c75760408051600160e51b62461bcd02815260206004820152601e60248201527f536166654d6174683a207375627472616374696f6e206f766572666c6f770000604482015290519081900360640190fd5b50900390565b60008282018381101561062a5760408051600160e51b62461bcd02815260206004820152601b60248201527f536166654d6174683a206164646974696f6e206f766572666c6f770000000000604482015290519081900360640190fd5b939250505056fe45524332303a207472616e7366657220746f20746865207a65726f206164647265737345524332303a20617070726f766520746f20746865207a65726f206164647265737345524332303a207472616e736665722066726f6d20746865207a65726f206164647265737345524332303a20617070726f76652066726f6d20746865207a65726f2061646472657373a165627a7a72305820802770bd4b0917c80d6b47964550f6c3e437f4a0b7286b55aa9963a98c5aaf400029`

// ERC20FuncSigs maps the 4-byte function signature to its string representation.
var ERC20FuncSigs = map[string]string{
	"dd62ed3e": "allowance(address,address)",
	"095ea7b3": "approve(address,uint256)",
	"70a08231": "balanceOf(address)",
	"a457c2d7": "decreaseAllowance(address,uint256)",
	"39509351": "increaseAllowance(address,uint256)",
	"18160ddd": "totalSupply()",
	"a9059cbb": "transfer(address,uint256)",
	"23b872dd": "transferFrom(address,address,uint256)",
}

// ERC20Bin is the compiled bytecode used for deploying new contracts.
var ERC20Bin = "0x608060405234801561001057600080fd5b506106eb806100206000396000f3fe608060405234801561001057600080fd5b50600436106100885760003560e01c806370a082311161005b57806370a0823114610149578063a457c2d71461016f578063a9059cbb1461019b578063dd62ed3e146101c757610088565b8063095ea7b31461008d57806318160ddd146100cd57806323b872dd146100e7578063395093511461011d575b600080fd5b6100b9600480360360408110156100a357600080fd5b506001600160a01b0381351690602001356101f5565b604080519115158252519081900360200190f35b6100d561020b565b60408051918252519081900360200190f35b6100b9600480360360608110156100fd57600080fd5b506001600160a01b03813581169160208101359091169060400135610211565b6100b96004803603604081101561013357600080fd5b506001600160a01b038135169060200135610268565b6100d56004803603602081101561015f57600080fd5b50356001600160a01b03166102a4565b6100b96004803603604081101561018557600080fd5b506001600160a01b0381351690602001356102bf565b6100b9600480360360408110156101b157600080fd5b506001600160a01b0381351690602001356102fb565b6100d5600480360360408110156101dd57600080fd5b506001600160a01b0381358116916020013516610308565b6000610202338484610333565b50600192915050565b60025490565b600061021e848484610425565b6001600160a01b03841660009081526001602090815260408083203380855292529091205461025e918691610259908663ffffffff61056d16565b610333565b5060019392505050565b3360008181526001602090815260408083206001600160a01b03871684529091528120549091610202918590610259908663ffffffff6105cd16565b6001600160a01b031660009081526020819052604090205490565b3360008181526001602090815260408083206001600160a01b03871684529091528120549091610202918590610259908663ffffffff61056d16565b6000610202338484610425565b6001600160a01b03918216600090815260016020908152604080832093909416825291909152205490565b6001600160a01b03831661037b57604051600160e51b62461bcd02815260040180806020018281038252602481526020018061069c6024913960400191505060405180910390fd5b6001600160a01b0382166103c357604051600160e51b62461bcd0281526004018080602001828103825260228152602001806106556022913960400191505060405180910390fd5b6001600160a01b03808416600081815260016020908152604080832094871680845294825291829020859055815185815291517f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b9259281900390910190a3505050565b6001600160a01b03831661046d57604051600160e51b62461bcd0281526004018080602001828103825260258152602001806106776025913960400191505060405180910390fd5b6001600160a01b0382166104b557604051600160e51b62461bcd0281526004018080602001828103825260238152602001806106326023913960400191505060405180910390fd5b6001600160a01b0383166000908152602081905260409020546104de908263ffffffff61056d16565b6001600160a01b038085166000908152602081905260408082209390935590841681522054610513908263ffffffff6105cd16565b6001600160a01b038084166000818152602081815260409182902094909455805185815290519193928716927fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef92918290030190a3505050565b6000828211156105c75760408051600160e51b62461bcd02815260206004820152601e60248201527f536166654d6174683a207375627472616374696f6e206f766572666c6f770000604482015290519081900360640190fd5b50900390565b60008282018381101561062a5760408051600160e51b62461bcd02815260206004820152601b60248201527f536166654d6174683a206164646974696f6e206f766572666c6f770000000000604482015290519081900360640190fd5b939250505056fe45524332303a207472616e7366657220746f20746865207a65726f206164647265737345524332303a20617070726f766520746f20746865207a65726f206164647265737345524332303a207472616e736665722066726f6d20746865207a65726f206164647265737345524332303a20617070726f76652066726f6d20746865207a65726f2061646472657373a165627a7a72305820802770bd4b0917c80d6b47964550f6c3e437f4a0b7286b55aa9963a98c5aaf400029"

// DeployERC20 deploys a new Klaytn contract, binding an instance of ERC20 to it.
func DeployERC20(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ERC20, error) {
	parsed, err := abi.JSON(strings.NewReader(ERC20ABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(ERC20Bin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ERC20{ERC20Caller: ERC20Caller{contract: contract}, ERC20Transactor: ERC20Transactor{contract: contract}, ERC20Filterer: ERC20Filterer{contract: contract}}, nil
}

// ERC20 is an auto generated Go binding around a Klaytn contract.
type ERC20 struct {
	ERC20Caller     // Read-only binding to the contract
	ERC20Transactor // Write-only binding to the contract
	ERC20Filterer   // Log filterer for contract events
}

// ERC20Caller is an auto generated read-only Go binding around a Klaytn contract.
type ERC20Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC20Transactor is an auto generated write-only Go binding around a Klaytn contract.
type ERC20Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC20Filterer is an auto generated log filtering Go binding around a Klaytn contract events.
type ERC20Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC20Session is an auto generated Go binding around a Klaytn contract,
// with pre-set call and transact options.
type ERC20Session struct {
	Contract     *ERC20            // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ERC20CallerSession is an auto generated read-only Go binding around a Klaytn contract,
// with pre-set call options.
type ERC20CallerSession struct {
	Contract *ERC20Caller  // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// ERC20TransactorSession is an auto generated write-only Go binding around a Klaytn contract,
// with pre-set transact options.
type ERC20TransactorSession struct {
	Contract     *ERC20Transactor  // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ERC20Raw is an auto generated low-level Go binding around a Klaytn contract.
type ERC20Raw struct {
	Contract *ERC20 // Generic contract binding to access the raw methods on
}

// ERC20CallerRaw is an auto generated low-level read-only Go binding around a Klaytn contract.
type ERC20CallerRaw struct {
	Contract *ERC20Caller // Generic read-only contract binding to access the raw methods on
}

// ERC20TransactorRaw is an auto generated low-level write-only Go binding around a Klaytn contract.
type ERC20TransactorRaw struct {
	Contract *ERC20Transactor // Generic write-only contract binding to access the raw methods on
}

// NewERC20 creates a new instance of ERC20, bound to a specific deployed contract.
func NewERC20(address common.Address, backend bind.ContractBackend) (*ERC20, error) {
	contract, err := bindERC20(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ERC20{ERC20Caller: ERC20Caller{contract: contract}, ERC20Transactor: ERC20Transactor{contract: contract}, ERC20Filterer: ERC20Filterer{contract: contract}}, nil
}

// NewERC20Caller creates a new read-only instance of ERC20, bound to a specific deployed contract.
func NewERC20Caller(address common.Address, caller bind.ContractCaller) (*ERC20Caller, error) {
	contract, err := bindERC20(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ERC20Caller{contract: contract}, nil
}

// NewERC20Transactor creates a new write-only instance of ERC20, bound to a specific deployed contract.
func NewERC20Transactor(address common.Address, transactor bind.ContractTransactor) (*ERC20Transactor, error) {
	contract, err := bindERC20(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ERC20Transactor{contract: contract}, nil
}

// NewERC20Filterer creates a new log filterer instance of ERC20, bound to a specific deployed contract.
func NewERC20Filterer(address common.Address, filterer bind.ContractFilterer) (*ERC20Filterer, error) {
	contract, err := bindERC20(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ERC20Filterer{contract: contract}, nil
}

// bindERC20 binds a generic wrapper to an already deployed contract.
func bindERC20(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ERC20ABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ERC20 *ERC20Raw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _ERC20.Contract.ERC20Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ERC20 *ERC20Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC20.Contract.ERC20Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ERC20 *ERC20Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ERC20.Contract.ERC20Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ERC20 *ERC20CallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _ERC20.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ERC20 *ERC20TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC20.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ERC20 *ERC20TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ERC20.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_ERC20 *ERC20Caller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ERC20.contract.Call(opts, out, "allowance", owner, spender)
	return *ret0, err
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_ERC20 *ERC20Session) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _ERC20.Contract.Allowance(&_ERC20.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_ERC20 *ERC20CallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _ERC20.Contract.Allowance(&_ERC20.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_ERC20 *ERC20Caller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ERC20.contract.Call(opts, out, "balanceOf", account)
	return *ret0, err
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_ERC20 *ERC20Session) BalanceOf(account common.Address) (*big.Int, error) {
	return _ERC20.Contract.BalanceOf(&_ERC20.CallOpts, account)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_ERC20 *ERC20CallerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _ERC20.Contract.BalanceOf(&_ERC20.CallOpts, account)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_ERC20 *ERC20Caller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ERC20.contract.Call(opts, out, "totalSupply")
	return *ret0, err
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_ERC20 *ERC20Session) TotalSupply() (*big.Int, error) {
	return _ERC20.Contract.TotalSupply(&_ERC20.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_ERC20 *ERC20CallerSession) TotalSupply() (*big.Int, error) {
	return _ERC20.Contract.TotalSupply(&_ERC20.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_ERC20 *ERC20Transactor) Approve(opts *bind.TransactOpts, spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _ERC20.contract.Transact(opts, "approve", spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_ERC20 *ERC20Session) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _ERC20.Contract.Approve(&_ERC20.TransactOpts, spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_ERC20 *ERC20TransactorSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _ERC20.Contract.Approve(&_ERC20.TransactOpts, spender, value)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 subtractedValue) returns(bool)
func (_ERC20 *ERC20Transactor) DecreaseAllowance(opts *bind.TransactOpts, spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _ERC20.contract.Transact(opts, "decreaseAllowance", spender, subtractedValue)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 subtractedValue) returns(bool)
func (_ERC20 *ERC20Session) DecreaseAllowance(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _ERC20.Contract.DecreaseAllowance(&_ERC20.TransactOpts, spender, subtractedValue)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 subtractedValue) returns(bool)
func (_ERC20 *ERC20TransactorSession) DecreaseAllowance(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _ERC20.Contract.DecreaseAllowance(&_ERC20.TransactOpts, spender, subtractedValue)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 addedValue) returns(bool)
func (_ERC20 *ERC20Transactor) IncreaseAllowance(opts *bind.TransactOpts, spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _ERC20.contract.Transact(opts, "increaseAllowance", spender, addedValue)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 addedValue) returns(bool)
func (_ERC20 *ERC20Session) IncreaseAllowance(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _ERC20.Contract.IncreaseAllowance(&_ERC20.TransactOpts, spender, addedValue)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 addedValue) returns(bool)
func (_ERC20 *ERC20TransactorSession) IncreaseAllowance(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _ERC20.Contract.IncreaseAllowance(&_ERC20.TransactOpts, spender, addedValue)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_ERC20 *ERC20Transactor) Transfer(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ERC20.contract.Transact(opts, "transfer", recipient, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_ERC20 *ERC20Session) Transfer(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ERC20.Contract.Transfer(&_ERC20.TransactOpts, recipient, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_ERC20 *ERC20TransactorSession) Transfer(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ERC20.Contract.Transfer(&_ERC20.TransactOpts, recipient, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_ERC20 *ERC20Transactor) TransferFrom(opts *bind.TransactOpts, sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ERC20.contract.Transact(opts, "transferFrom", sender, recipient, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_ERC20 *ERC20Session) TransferFrom(sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ERC20.Contract.TransferFrom(&_ERC20.TransactOpts, sender, recipient, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_ERC20 *ERC20TransactorSession) TransferFrom(sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ERC20.Contract.TransferFrom(&_ERC20.TransactOpts, sender, recipient, amount)
}

// ERC20ApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the ERC20 contract.
type ERC20ApprovalIterator struct {
	Event *ERC20Approval // Event containing the contract specifics and raw log

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
func (it *ERC20ApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC20Approval)
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
		it.Event = new(ERC20Approval)
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
func (it *ERC20ApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC20ApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC20Approval represents a Approval event raised by the ERC20 contract.
type ERC20Approval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_ERC20 *ERC20Filterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*ERC20ApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _ERC20.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &ERC20ApprovalIterator{contract: _ERC20.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_ERC20 *ERC20Filterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *ERC20Approval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _ERC20.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC20Approval)
				if err := _ERC20.contract.UnpackLog(event, "Approval", log); err != nil {
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

// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_ERC20 *ERC20Filterer) ParseApproval(log types.Log) (*ERC20Approval, error) {
	event := new(ERC20Approval)
	if err := _ERC20.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ERC20TransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the ERC20 contract.
type ERC20TransferIterator struct {
	Event *ERC20Transfer // Event containing the contract specifics and raw log

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
func (it *ERC20TransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC20Transfer)
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
		it.Event = new(ERC20Transfer)
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
func (it *ERC20TransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC20TransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC20Transfer represents a Transfer event raised by the ERC20 contract.
type ERC20Transfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_ERC20 *ERC20Filterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*ERC20TransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ERC20.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &ERC20TransferIterator{contract: _ERC20.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_ERC20 *ERC20Filterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *ERC20Transfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ERC20.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC20Transfer)
				if err := _ERC20.contract.UnpackLog(event, "Transfer", log); err != nil {
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

// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_ERC20 *ERC20Filterer) ParseTransfer(log types.Log) (*ERC20Transfer, error) {
	event := new(ERC20Transfer)
	if err := _ERC20.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ERC20BurnableABI is the input ABI used to generate the binding from.
const ERC20BurnableABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"spender\",\"type\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"sender\",\"type\":\"address\"},{\"name\":\"recipient\",\"type\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"spender\",\"type\":\"address\"},{\"name\":\"addedValue\",\"type\":\"uint256\"}],\"name\":\"increaseAllowance\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"account\",\"type\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"burnFrom\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"spender\",\"type\":\"address\"},{\"name\":\"subtractedValue\",\"type\":\"uint256\"}],\"name\":\"decreaseAllowance\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"recipient\",\"type\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"owner\",\"type\":\"address\"},{\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"}]"

// ERC20BurnableBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const ERC20BurnableBinRuntime = `608060405234801561001057600080fd5b506004361061009e5760003560e01c806370a082311161006657806370a082311461017e57806379cc6790146101a4578063a457c2d7146101d0578063a9059cbb146101fc578063dd62ed3e146102285761009e565b8063095ea7b3146100a357806318160ddd146100e357806323b872dd146100fd578063395093511461013357806342966c681461015f575b600080fd5b6100cf600480360360408110156100b957600080fd5b506001600160a01b038135169060200135610256565b604080519115158252519081900360200190f35b6100eb61026c565b60408051918252519081900360200190f35b6100cf6004803603606081101561011357600080fd5b506001600160a01b03813581169160208101359091169060400135610272565b6100cf6004803603604081101561014957600080fd5b506001600160a01b0381351690602001356102c9565b61017c6004803603602081101561017557600080fd5b5035610305565b005b6100eb6004803603602081101561019457600080fd5b50356001600160a01b0316610312565b61017c600480360360408110156101ba57600080fd5b506001600160a01b03813516906020013561032d565b6100cf600480360360408110156101e657600080fd5b506001600160a01b03813516906020013561033b565b6100cf6004803603604081101561021257600080fd5b506001600160a01b038135169060200135610377565b6100eb6004803603604081101561023e57600080fd5b506001600160a01b0381358116916020013516610384565b60006102633384846103af565b50600192915050565b60025490565b600061027f8484846104a1565b6001600160a01b0384166000908152600160209081526040808320338085529252909120546102bf9186916102ba908663ffffffff6105e916565b6103af565b5060019392505050565b3360008181526001602090815260408083206001600160a01b038716845290915281205490916102639185906102ba908663ffffffff61064916565b61030f33826106ad565b50565b6001600160a01b031660009081526020819052604090205490565b6103378282610789565b5050565b3360008181526001602090815260408083206001600160a01b038716845290915281205490916102639185906102ba908663ffffffff6105e916565b60006102633384846104a1565b6001600160a01b03918216600090815260016020908152604080832093909416825291909152205490565b6001600160a01b0383166103f757604051600160e51b62461bcd02815260040180806020018281038252602481526020018061085a6024913960400191505060405180910390fd5b6001600160a01b03821661043f57604051600160e51b62461bcd0281526004018080602001828103825260228152602001806107f26022913960400191505060405180910390fd5b6001600160a01b03808416600081815260016020908152604080832094871680845294825291829020859055815185815291517f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b9259281900390910190a3505050565b6001600160a01b0383166104e957604051600160e51b62461bcd0281526004018080602001828103825260258152602001806108356025913960400191505060405180910390fd5b6001600160a01b03821661053157604051600160e51b62461bcd0281526004018080602001828103825260238152602001806107cf6023913960400191505060405180910390fd5b6001600160a01b03831660009081526020819052604090205461055a908263ffffffff6105e916565b6001600160a01b03808516600090815260208190526040808220939093559084168152205461058f908263ffffffff61064916565b6001600160a01b038084166000818152602081815260409182902094909455805185815290519193928716927fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef92918290030190a3505050565b6000828211156106435760408051600160e51b62461bcd02815260206004820152601e60248201527f536166654d6174683a207375627472616374696f6e206f766572666c6f770000604482015290519081900360640190fd5b50900390565b6000828201838110156106a65760408051600160e51b62461bcd02815260206004820152601b60248201527f536166654d6174683a206164646974696f6e206f766572666c6f770000000000604482015290519081900360640190fd5b9392505050565b6001600160a01b0382166106f557604051600160e51b62461bcd0281526004018080602001828103825260218152602001806108146021913960400191505060405180910390fd5b600254610708908263ffffffff6105e916565b6002556001600160a01b038216600090815260208190526040902054610734908263ffffffff6105e916565b6001600160a01b038316600081815260208181526040808320949094558351858152935191937fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef929081900390910190a35050565b61079382826106ad565b6001600160a01b0382166000908152600160209081526040808320338085529252909120546103379184916102ba908563ffffffff6105e91656fe45524332303a207472616e7366657220746f20746865207a65726f206164647265737345524332303a20617070726f766520746f20746865207a65726f206164647265737345524332303a206275726e2066726f6d20746865207a65726f206164647265737345524332303a207472616e736665722066726f6d20746865207a65726f206164647265737345524332303a20617070726f76652066726f6d20746865207a65726f2061646472657373a165627a7a7230582013b4ee4a596e91b545aeb6447c3788b121a99567674027b21b642c7da8dc493c0029`

// ERC20BurnableFuncSigs maps the 4-byte function signature to its string representation.
var ERC20BurnableFuncSigs = map[string]string{
	"dd62ed3e": "allowance(address,address)",
	"095ea7b3": "approve(address,uint256)",
	"70a08231": "balanceOf(address)",
	"42966c68": "burn(uint256)",
	"79cc6790": "burnFrom(address,uint256)",
	"a457c2d7": "decreaseAllowance(address,uint256)",
	"39509351": "increaseAllowance(address,uint256)",
	"18160ddd": "totalSupply()",
	"a9059cbb": "transfer(address,uint256)",
	"23b872dd": "transferFrom(address,address,uint256)",
}

// ERC20BurnableBin is the compiled bytecode used for deploying new contracts.
var ERC20BurnableBin = "0x608060405234801561001057600080fd5b506108a9806100206000396000f3fe608060405234801561001057600080fd5b506004361061009e5760003560e01c806370a082311161006657806370a082311461017e57806379cc6790146101a4578063a457c2d7146101d0578063a9059cbb146101fc578063dd62ed3e146102285761009e565b8063095ea7b3146100a357806318160ddd146100e357806323b872dd146100fd578063395093511461013357806342966c681461015f575b600080fd5b6100cf600480360360408110156100b957600080fd5b506001600160a01b038135169060200135610256565b604080519115158252519081900360200190f35b6100eb61026c565b60408051918252519081900360200190f35b6100cf6004803603606081101561011357600080fd5b506001600160a01b03813581169160208101359091169060400135610272565b6100cf6004803603604081101561014957600080fd5b506001600160a01b0381351690602001356102c9565b61017c6004803603602081101561017557600080fd5b5035610305565b005b6100eb6004803603602081101561019457600080fd5b50356001600160a01b0316610312565b61017c600480360360408110156101ba57600080fd5b506001600160a01b03813516906020013561032d565b6100cf600480360360408110156101e657600080fd5b506001600160a01b03813516906020013561033b565b6100cf6004803603604081101561021257600080fd5b506001600160a01b038135169060200135610377565b6100eb6004803603604081101561023e57600080fd5b506001600160a01b0381358116916020013516610384565b60006102633384846103af565b50600192915050565b60025490565b600061027f8484846104a1565b6001600160a01b0384166000908152600160209081526040808320338085529252909120546102bf9186916102ba908663ffffffff6105e916565b6103af565b5060019392505050565b3360008181526001602090815260408083206001600160a01b038716845290915281205490916102639185906102ba908663ffffffff61064916565b61030f33826106ad565b50565b6001600160a01b031660009081526020819052604090205490565b6103378282610789565b5050565b3360008181526001602090815260408083206001600160a01b038716845290915281205490916102639185906102ba908663ffffffff6105e916565b60006102633384846104a1565b6001600160a01b03918216600090815260016020908152604080832093909416825291909152205490565b6001600160a01b0383166103f757604051600160e51b62461bcd02815260040180806020018281038252602481526020018061085a6024913960400191505060405180910390fd5b6001600160a01b03821661043f57604051600160e51b62461bcd0281526004018080602001828103825260228152602001806107f26022913960400191505060405180910390fd5b6001600160a01b03808416600081815260016020908152604080832094871680845294825291829020859055815185815291517f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b9259281900390910190a3505050565b6001600160a01b0383166104e957604051600160e51b62461bcd0281526004018080602001828103825260258152602001806108356025913960400191505060405180910390fd5b6001600160a01b03821661053157604051600160e51b62461bcd0281526004018080602001828103825260238152602001806107cf6023913960400191505060405180910390fd5b6001600160a01b03831660009081526020819052604090205461055a908263ffffffff6105e916565b6001600160a01b03808516600090815260208190526040808220939093559084168152205461058f908263ffffffff61064916565b6001600160a01b038084166000818152602081815260409182902094909455805185815290519193928716927fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef92918290030190a3505050565b6000828211156106435760408051600160e51b62461bcd02815260206004820152601e60248201527f536166654d6174683a207375627472616374696f6e206f766572666c6f770000604482015290519081900360640190fd5b50900390565b6000828201838110156106a65760408051600160e51b62461bcd02815260206004820152601b60248201527f536166654d6174683a206164646974696f6e206f766572666c6f770000000000604482015290519081900360640190fd5b9392505050565b6001600160a01b0382166106f557604051600160e51b62461bcd0281526004018080602001828103825260218152602001806108146021913960400191505060405180910390fd5b600254610708908263ffffffff6105e916565b6002556001600160a01b038216600090815260208190526040902054610734908263ffffffff6105e916565b6001600160a01b038316600081815260208181526040808320949094558351858152935191937fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef929081900390910190a35050565b61079382826106ad565b6001600160a01b0382166000908152600160209081526040808320338085529252909120546103379184916102ba908563ffffffff6105e91656fe45524332303a207472616e7366657220746f20746865207a65726f206164647265737345524332303a20617070726f766520746f20746865207a65726f206164647265737345524332303a206275726e2066726f6d20746865207a65726f206164647265737345524332303a207472616e736665722066726f6d20746865207a65726f206164647265737345524332303a20617070726f76652066726f6d20746865207a65726f2061646472657373a165627a7a7230582013b4ee4a596e91b545aeb6447c3788b121a99567674027b21b642c7da8dc493c0029"

// DeployERC20Burnable deploys a new Klaytn contract, binding an instance of ERC20Burnable to it.
func DeployERC20Burnable(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ERC20Burnable, error) {
	parsed, err := abi.JSON(strings.NewReader(ERC20BurnableABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(ERC20BurnableBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ERC20Burnable{ERC20BurnableCaller: ERC20BurnableCaller{contract: contract}, ERC20BurnableTransactor: ERC20BurnableTransactor{contract: contract}, ERC20BurnableFilterer: ERC20BurnableFilterer{contract: contract}}, nil
}

// ERC20Burnable is an auto generated Go binding around a Klaytn contract.
type ERC20Burnable struct {
	ERC20BurnableCaller     // Read-only binding to the contract
	ERC20BurnableTransactor // Write-only binding to the contract
	ERC20BurnableFilterer   // Log filterer for contract events
}

// ERC20BurnableCaller is an auto generated read-only Go binding around a Klaytn contract.
type ERC20BurnableCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC20BurnableTransactor is an auto generated write-only Go binding around a Klaytn contract.
type ERC20BurnableTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC20BurnableFilterer is an auto generated log filtering Go binding around a Klaytn contract events.
type ERC20BurnableFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC20BurnableSession is an auto generated Go binding around a Klaytn contract,
// with pre-set call and transact options.
type ERC20BurnableSession struct {
	Contract     *ERC20Burnable    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ERC20BurnableCallerSession is an auto generated read-only Go binding around a Klaytn contract,
// with pre-set call options.
type ERC20BurnableCallerSession struct {
	Contract *ERC20BurnableCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// ERC20BurnableTransactorSession is an auto generated write-only Go binding around a Klaytn contract,
// with pre-set transact options.
type ERC20BurnableTransactorSession struct {
	Contract     *ERC20BurnableTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// ERC20BurnableRaw is an auto generated low-level Go binding around a Klaytn contract.
type ERC20BurnableRaw struct {
	Contract *ERC20Burnable // Generic contract binding to access the raw methods on
}

// ERC20BurnableCallerRaw is an auto generated low-level read-only Go binding around a Klaytn contract.
type ERC20BurnableCallerRaw struct {
	Contract *ERC20BurnableCaller // Generic read-only contract binding to access the raw methods on
}

// ERC20BurnableTransactorRaw is an auto generated low-level write-only Go binding around a Klaytn contract.
type ERC20BurnableTransactorRaw struct {
	Contract *ERC20BurnableTransactor // Generic write-only contract binding to access the raw methods on
}

// NewERC20Burnable creates a new instance of ERC20Burnable, bound to a specific deployed contract.
func NewERC20Burnable(address common.Address, backend bind.ContractBackend) (*ERC20Burnable, error) {
	contract, err := bindERC20Burnable(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ERC20Burnable{ERC20BurnableCaller: ERC20BurnableCaller{contract: contract}, ERC20BurnableTransactor: ERC20BurnableTransactor{contract: contract}, ERC20BurnableFilterer: ERC20BurnableFilterer{contract: contract}}, nil
}

// NewERC20BurnableCaller creates a new read-only instance of ERC20Burnable, bound to a specific deployed contract.
func NewERC20BurnableCaller(address common.Address, caller bind.ContractCaller) (*ERC20BurnableCaller, error) {
	contract, err := bindERC20Burnable(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ERC20BurnableCaller{contract: contract}, nil
}

// NewERC20BurnableTransactor creates a new write-only instance of ERC20Burnable, bound to a specific deployed contract.
func NewERC20BurnableTransactor(address common.Address, transactor bind.ContractTransactor) (*ERC20BurnableTransactor, error) {
	contract, err := bindERC20Burnable(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ERC20BurnableTransactor{contract: contract}, nil
}

// NewERC20BurnableFilterer creates a new log filterer instance of ERC20Burnable, bound to a specific deployed contract.
func NewERC20BurnableFilterer(address common.Address, filterer bind.ContractFilterer) (*ERC20BurnableFilterer, error) {
	contract, err := bindERC20Burnable(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ERC20BurnableFilterer{contract: contract}, nil
}

// bindERC20Burnable binds a generic wrapper to an already deployed contract.
func bindERC20Burnable(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ERC20BurnableABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ERC20Burnable *ERC20BurnableRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _ERC20Burnable.Contract.ERC20BurnableCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ERC20Burnable *ERC20BurnableRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC20Burnable.Contract.ERC20BurnableTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ERC20Burnable *ERC20BurnableRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ERC20Burnable.Contract.ERC20BurnableTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ERC20Burnable *ERC20BurnableCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _ERC20Burnable.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ERC20Burnable *ERC20BurnableTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC20Burnable.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ERC20Burnable *ERC20BurnableTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ERC20Burnable.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_ERC20Burnable *ERC20BurnableCaller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ERC20Burnable.contract.Call(opts, out, "allowance", owner, spender)
	return *ret0, err
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_ERC20Burnable *ERC20BurnableSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _ERC20Burnable.Contract.Allowance(&_ERC20Burnable.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_ERC20Burnable *ERC20BurnableCallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _ERC20Burnable.Contract.Allowance(&_ERC20Burnable.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_ERC20Burnable *ERC20BurnableCaller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ERC20Burnable.contract.Call(opts, out, "balanceOf", account)
	return *ret0, err
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_ERC20Burnable *ERC20BurnableSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _ERC20Burnable.Contract.BalanceOf(&_ERC20Burnable.CallOpts, account)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_ERC20Burnable *ERC20BurnableCallerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _ERC20Burnable.Contract.BalanceOf(&_ERC20Burnable.CallOpts, account)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_ERC20Burnable *ERC20BurnableCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ERC20Burnable.contract.Call(opts, out, "totalSupply")
	return *ret0, err
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_ERC20Burnable *ERC20BurnableSession) TotalSupply() (*big.Int, error) {
	return _ERC20Burnable.Contract.TotalSupply(&_ERC20Burnable.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_ERC20Burnable *ERC20BurnableCallerSession) TotalSupply() (*big.Int, error) {
	return _ERC20Burnable.Contract.TotalSupply(&_ERC20Burnable.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_ERC20Burnable *ERC20BurnableTransactor) Approve(opts *bind.TransactOpts, spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _ERC20Burnable.contract.Transact(opts, "approve", spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_ERC20Burnable *ERC20BurnableSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _ERC20Burnable.Contract.Approve(&_ERC20Burnable.TransactOpts, spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_ERC20Burnable *ERC20BurnableTransactorSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _ERC20Burnable.Contract.Approve(&_ERC20Burnable.TransactOpts, spender, value)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 amount) returns()
func (_ERC20Burnable *ERC20BurnableTransactor) Burn(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _ERC20Burnable.contract.Transact(opts, "burn", amount)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 amount) returns()
func (_ERC20Burnable *ERC20BurnableSession) Burn(amount *big.Int) (*types.Transaction, error) {
	return _ERC20Burnable.Contract.Burn(&_ERC20Burnable.TransactOpts, amount)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 amount) returns()
func (_ERC20Burnable *ERC20BurnableTransactorSession) Burn(amount *big.Int) (*types.Transaction, error) {
	return _ERC20Burnable.Contract.Burn(&_ERC20Burnable.TransactOpts, amount)
}

// BurnFrom is a paid mutator transaction binding the contract method 0x79cc6790.
//
// Solidity: function burnFrom(address account, uint256 amount) returns()
func (_ERC20Burnable *ERC20BurnableTransactor) BurnFrom(opts *bind.TransactOpts, account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ERC20Burnable.contract.Transact(opts, "burnFrom", account, amount)
}

// BurnFrom is a paid mutator transaction binding the contract method 0x79cc6790.
//
// Solidity: function burnFrom(address account, uint256 amount) returns()
func (_ERC20Burnable *ERC20BurnableSession) BurnFrom(account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ERC20Burnable.Contract.BurnFrom(&_ERC20Burnable.TransactOpts, account, amount)
}

// BurnFrom is a paid mutator transaction binding the contract method 0x79cc6790.
//
// Solidity: function burnFrom(address account, uint256 amount) returns()
func (_ERC20Burnable *ERC20BurnableTransactorSession) BurnFrom(account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ERC20Burnable.Contract.BurnFrom(&_ERC20Burnable.TransactOpts, account, amount)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 subtractedValue) returns(bool)
func (_ERC20Burnable *ERC20BurnableTransactor) DecreaseAllowance(opts *bind.TransactOpts, spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _ERC20Burnable.contract.Transact(opts, "decreaseAllowance", spender, subtractedValue)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 subtractedValue) returns(bool)
func (_ERC20Burnable *ERC20BurnableSession) DecreaseAllowance(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _ERC20Burnable.Contract.DecreaseAllowance(&_ERC20Burnable.TransactOpts, spender, subtractedValue)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 subtractedValue) returns(bool)
func (_ERC20Burnable *ERC20BurnableTransactorSession) DecreaseAllowance(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _ERC20Burnable.Contract.DecreaseAllowance(&_ERC20Burnable.TransactOpts, spender, subtractedValue)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 addedValue) returns(bool)
func (_ERC20Burnable *ERC20BurnableTransactor) IncreaseAllowance(opts *bind.TransactOpts, spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _ERC20Burnable.contract.Transact(opts, "increaseAllowance", spender, addedValue)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 addedValue) returns(bool)
func (_ERC20Burnable *ERC20BurnableSession) IncreaseAllowance(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _ERC20Burnable.Contract.IncreaseAllowance(&_ERC20Burnable.TransactOpts, spender, addedValue)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 addedValue) returns(bool)
func (_ERC20Burnable *ERC20BurnableTransactorSession) IncreaseAllowance(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _ERC20Burnable.Contract.IncreaseAllowance(&_ERC20Burnable.TransactOpts, spender, addedValue)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_ERC20Burnable *ERC20BurnableTransactor) Transfer(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ERC20Burnable.contract.Transact(opts, "transfer", recipient, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_ERC20Burnable *ERC20BurnableSession) Transfer(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ERC20Burnable.Contract.Transfer(&_ERC20Burnable.TransactOpts, recipient, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_ERC20Burnable *ERC20BurnableTransactorSession) Transfer(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ERC20Burnable.Contract.Transfer(&_ERC20Burnable.TransactOpts, recipient, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_ERC20Burnable *ERC20BurnableTransactor) TransferFrom(opts *bind.TransactOpts, sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ERC20Burnable.contract.Transact(opts, "transferFrom", sender, recipient, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_ERC20Burnable *ERC20BurnableSession) TransferFrom(sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ERC20Burnable.Contract.TransferFrom(&_ERC20Burnable.TransactOpts, sender, recipient, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_ERC20Burnable *ERC20BurnableTransactorSession) TransferFrom(sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ERC20Burnable.Contract.TransferFrom(&_ERC20Burnable.TransactOpts, sender, recipient, amount)
}

// ERC20BurnableApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the ERC20Burnable contract.
type ERC20BurnableApprovalIterator struct {
	Event *ERC20BurnableApproval // Event containing the contract specifics and raw log

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
func (it *ERC20BurnableApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC20BurnableApproval)
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
		it.Event = new(ERC20BurnableApproval)
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
func (it *ERC20BurnableApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC20BurnableApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC20BurnableApproval represents a Approval event raised by the ERC20Burnable contract.
type ERC20BurnableApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_ERC20Burnable *ERC20BurnableFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*ERC20BurnableApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _ERC20Burnable.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &ERC20BurnableApprovalIterator{contract: _ERC20Burnable.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_ERC20Burnable *ERC20BurnableFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *ERC20BurnableApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _ERC20Burnable.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC20BurnableApproval)
				if err := _ERC20Burnable.contract.UnpackLog(event, "Approval", log); err != nil {
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

// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_ERC20Burnable *ERC20BurnableFilterer) ParseApproval(log types.Log) (*ERC20BurnableApproval, error) {
	event := new(ERC20BurnableApproval)
	if err := _ERC20Burnable.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ERC20BurnableTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the ERC20Burnable contract.
type ERC20BurnableTransferIterator struct {
	Event *ERC20BurnableTransfer // Event containing the contract specifics and raw log

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
func (it *ERC20BurnableTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC20BurnableTransfer)
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
		it.Event = new(ERC20BurnableTransfer)
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
func (it *ERC20BurnableTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC20BurnableTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC20BurnableTransfer represents a Transfer event raised by the ERC20Burnable contract.
type ERC20BurnableTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_ERC20Burnable *ERC20BurnableFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*ERC20BurnableTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ERC20Burnable.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &ERC20BurnableTransferIterator{contract: _ERC20Burnable.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_ERC20Burnable *ERC20BurnableFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *ERC20BurnableTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ERC20Burnable.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC20BurnableTransfer)
				if err := _ERC20Burnable.contract.UnpackLog(event, "Transfer", log); err != nil {
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

// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_ERC20Burnable *ERC20BurnableFilterer) ParseTransfer(log types.Log) (*ERC20BurnableTransfer, error) {
	event := new(ERC20BurnableTransfer)
	if err := _ERC20Burnable.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ERC20MintableABI is the input ABI used to generate the binding from.
const ERC20MintableABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"spender\",\"type\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"sender\",\"type\":\"address\"},{\"name\":\"recipient\",\"type\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"spender\",\"type\":\"address\"},{\"name\":\"addedValue\",\"type\":\"uint256\"}],\"name\":\"increaseAllowance\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"account\",\"type\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"mint\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"account\",\"type\":\"address\"}],\"name\":\"addMinter\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"renounceMinter\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"spender\",\"type\":\"address\"},{\"name\":\"subtractedValue\",\"type\":\"uint256\"}],\"name\":\"decreaseAllowance\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"recipient\",\"type\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"account\",\"type\":\"address\"}],\"name\":\"isMinter\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"owner\",\"type\":\"address\"},{\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"account\",\"type\":\"address\"}],\"name\":\"MinterAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"account\",\"type\":\"address\"}],\"name\":\"MinterRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"}]"

// ERC20MintableBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const ERC20MintableBinRuntime = `608060405234801561001057600080fd5b50600436106100b45760003560e01c8063983b2d5611610071578063983b2d56146101c757806398650275146101ef578063a457c2d7146101f7578063a9059cbb14610223578063aa271e1a1461024f578063dd62ed3e14610275576100b4565b8063095ea7b3146100b957806318160ddd146100f957806323b872dd14610113578063395093511461014957806340c10f191461017557806370a08231146101a1575b600080fd5b6100e5600480360360408110156100cf57600080fd5b506001600160a01b0381351690602001356102a3565b604080519115158252519081900360200190f35b6101016102b9565b60408051918252519081900360200190f35b6100e56004803603606081101561012957600080fd5b506001600160a01b038135811691602081013590911690604001356102bf565b6100e56004803603604081101561015f57600080fd5b506001600160a01b038135169060200135610316565b6100e56004803603604081101561018b57600080fd5b506001600160a01b038135169060200135610352565b610101600480360360208110156101b757600080fd5b50356001600160a01b03166103a5565b6101ed600480360360208110156101dd57600080fd5b50356001600160a01b03166103c0565b005b6101ed610413565b6100e56004803603604081101561020d57600080fd5b506001600160a01b03813516906020013561041e565b6100e56004803603604081101561023957600080fd5b506001600160a01b03813516906020013561045a565b6100e56004803603602081101561026557600080fd5b50356001600160a01b0316610467565b6101016004803603604081101561028b57600080fd5b506001600160a01b0381358116916020013516610480565b60006102b03384846104ab565b50600192915050565b60025490565b60006102cc84848461059d565b6001600160a01b03841660009081526001602090815260408083203380855292529091205461030c918691610307908663ffffffff6106e516565b6104ab565b5060019392505050565b3360008181526001602090815260408083206001600160a01b038716845290915281205490916102b0918590610307908663ffffffff61074516565b600061035d33610467565b61039b57604051600160e51b62461bcd028152600401808060200182810382526030815260200180610aca6030913960400191505060405180910390fd5b6102b083836107a9565b6001600160a01b031660009081526020819052604090205490565b6103c933610467565b61040757604051600160e51b62461bcd028152600401808060200182810382526030815260200180610aca6030913960400191505060405180910390fd5b6104108161089c565b50565b61041c336108e4565b565b3360008181526001602090815260408083206001600160a01b038716845290915281205490916102b0918590610307908663ffffffff6106e516565b60006102b033848461059d565b600061047a60038363ffffffff61092c16565b92915050565b6001600160a01b03918216600090815260016020908152604080832093909416825291909152205490565b6001600160a01b0383166104f357604051600160e51b62461bcd028152600401808060200182810382526024815260200180610b626024913960400191505060405180910390fd5b6001600160a01b03821661053b57604051600160e51b62461bcd028152600401808060200182810382526022815260200180610aa86022913960400191505060405180910390fd5b6001600160a01b03808416600081815260016020908152604080832094871680845294825291829020859055815185815291517f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b9259281900390910190a3505050565b6001600160a01b0383166105e557604051600160e51b62461bcd028152600401808060200182810382526025815260200180610b3d6025913960400191505060405180910390fd5b6001600160a01b03821661062d57604051600160e51b62461bcd028152600401808060200182810382526023815260200180610a856023913960400191505060405180910390fd5b6001600160a01b038316600090815260208190526040902054610656908263ffffffff6106e516565b6001600160a01b03808516600090815260208190526040808220939093559084168152205461068b908263ffffffff61074516565b6001600160a01b038084166000818152602081815260409182902094909455805185815290519193928716927fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef92918290030190a3505050565b60008282111561073f5760408051600160e51b62461bcd02815260206004820152601e60248201527f536166654d6174683a207375627472616374696f6e206f766572666c6f770000604482015290519081900360640190fd5b50900390565b6000828201838110156107a25760408051600160e51b62461bcd02815260206004820152601b60248201527f536166654d6174683a206164646974696f6e206f766572666c6f770000000000604482015290519081900360640190fd5b9392505050565b6001600160a01b0382166108075760408051600160e51b62461bcd02815260206004820152601f60248201527f45524332303a206d696e7420746f20746865207a65726f206164647265737300604482015290519081900360640190fd5b60025461081a908263ffffffff61074516565b6002556001600160a01b038216600090815260208190526040902054610846908263ffffffff61074516565b6001600160a01b0383166000818152602081815260408083209490945583518581529351929391927fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef9281900390910190a35050565b6108ad60038263ffffffff61099616565b6040516001600160a01b038216907f6ae172837ea30b801fbfcdd4108aa1d5bf8ff775444fd70256b44e6bf3dfc3f690600090a250565b6108f560038263ffffffff610a1a16565b6040516001600160a01b038216907fe94479a9f7e1952cc78f2d6baab678adc1b772d936c6583def489e524cb6669290600090a250565b60006001600160a01b03821661097657604051600160e51b62461bcd028152600401808060200182810382526022815260200180610b1b6022913960400191505060405180910390fd5b506001600160a01b03166000908152602091909152604090205460ff1690565b6109a0828261092c565b156109f55760408051600160e51b62461bcd02815260206004820152601f60248201527f526f6c65733a206163636f756e7420616c72656164792068617320726f6c6500604482015290519081900360640190fd5b6001600160a01b0316600090815260209190915260409020805460ff19166001179055565b610a24828261092c565b610a6257604051600160e51b62461bcd028152600401808060200182810382526021815260200180610afa6021913960400191505060405180910390fd5b6001600160a01b0316600090815260209190915260409020805460ff1916905556fe45524332303a207472616e7366657220746f20746865207a65726f206164647265737345524332303a20617070726f766520746f20746865207a65726f20616464726573734d696e746572526f6c653a2063616c6c657220646f6573206e6f74206861766520746865204d696e74657220726f6c65526f6c65733a206163636f756e7420646f6573206e6f74206861766520726f6c65526f6c65733a206163636f756e7420697320746865207a65726f206164647265737345524332303a207472616e736665722066726f6d20746865207a65726f206164647265737345524332303a20617070726f76652066726f6d20746865207a65726f2061646472657373a165627a7a723058208533d6897cadaf6784525a9d32e613a3c81da44cfdedf624f1bdb497f8b68f4a0029`

// ERC20MintableFuncSigs maps the 4-byte function signature to its string representation.
var ERC20MintableFuncSigs = map[string]string{
	"983b2d56": "addMinter(address)",
	"dd62ed3e": "allowance(address,address)",
	"095ea7b3": "approve(address,uint256)",
	"70a08231": "balanceOf(address)",
	"a457c2d7": "decreaseAllowance(address,uint256)",
	"39509351": "increaseAllowance(address,uint256)",
	"aa271e1a": "isMinter(address)",
	"40c10f19": "mint(address,uint256)",
	"98650275": "renounceMinter()",
	"18160ddd": "totalSupply()",
	"a9059cbb": "transfer(address,uint256)",
	"23b872dd": "transferFrom(address,address,uint256)",
}

// ERC20MintableBin is the compiled bytecode used for deploying new contracts.
var ERC20MintableBin = "0x60806040526100133361001860201b60201c565b610189565b61003081600361006760201b6109961790919060201c565b6040516001600160a01b038216907f6ae172837ea30b801fbfcdd4108aa1d5bf8ff775444fd70256b44e6bf3dfc3f690600090a250565b610077828261010860201b60201c565b156100e357604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601f60248201527f526f6c65733a206163636f756e7420616c72656164792068617320726f6c6500604482015290519081900360640190fd5b6001600160a01b0316600090815260209190915260409020805460ff19166001179055565b60006001600160a01b038216610169576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401808060200182810382526022815260200180610d496022913960400191505060405180910390fd5b506001600160a01b03166000908152602091909152604090205460ff1690565b610bb1806101986000396000f3fe608060405234801561001057600080fd5b50600436106100b45760003560e01c8063983b2d5611610071578063983b2d56146101c757806398650275146101ef578063a457c2d7146101f7578063a9059cbb14610223578063aa271e1a1461024f578063dd62ed3e14610275576100b4565b8063095ea7b3146100b957806318160ddd146100f957806323b872dd14610113578063395093511461014957806340c10f191461017557806370a08231146101a1575b600080fd5b6100e5600480360360408110156100cf57600080fd5b506001600160a01b0381351690602001356102a3565b604080519115158252519081900360200190f35b6101016102b9565b60408051918252519081900360200190f35b6100e56004803603606081101561012957600080fd5b506001600160a01b038135811691602081013590911690604001356102bf565b6100e56004803603604081101561015f57600080fd5b506001600160a01b038135169060200135610316565b6100e56004803603604081101561018b57600080fd5b506001600160a01b038135169060200135610352565b610101600480360360208110156101b757600080fd5b50356001600160a01b03166103a5565b6101ed600480360360208110156101dd57600080fd5b50356001600160a01b03166103c0565b005b6101ed610413565b6100e56004803603604081101561020d57600080fd5b506001600160a01b03813516906020013561041e565b6100e56004803603604081101561023957600080fd5b506001600160a01b03813516906020013561045a565b6100e56004803603602081101561026557600080fd5b50356001600160a01b0316610467565b6101016004803603604081101561028b57600080fd5b506001600160a01b0381358116916020013516610480565b60006102b03384846104ab565b50600192915050565b60025490565b60006102cc84848461059d565b6001600160a01b03841660009081526001602090815260408083203380855292529091205461030c918691610307908663ffffffff6106e516565b6104ab565b5060019392505050565b3360008181526001602090815260408083206001600160a01b038716845290915281205490916102b0918590610307908663ffffffff61074516565b600061035d33610467565b61039b57604051600160e51b62461bcd028152600401808060200182810382526030815260200180610aca6030913960400191505060405180910390fd5b6102b083836107a9565b6001600160a01b031660009081526020819052604090205490565b6103c933610467565b61040757604051600160e51b62461bcd028152600401808060200182810382526030815260200180610aca6030913960400191505060405180910390fd5b6104108161089c565b50565b61041c336108e4565b565b3360008181526001602090815260408083206001600160a01b038716845290915281205490916102b0918590610307908663ffffffff6106e516565b60006102b033848461059d565b600061047a60038363ffffffff61092c16565b92915050565b6001600160a01b03918216600090815260016020908152604080832093909416825291909152205490565b6001600160a01b0383166104f357604051600160e51b62461bcd028152600401808060200182810382526024815260200180610b626024913960400191505060405180910390fd5b6001600160a01b03821661053b57604051600160e51b62461bcd028152600401808060200182810382526022815260200180610aa86022913960400191505060405180910390fd5b6001600160a01b03808416600081815260016020908152604080832094871680845294825291829020859055815185815291517f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b9259281900390910190a3505050565b6001600160a01b0383166105e557604051600160e51b62461bcd028152600401808060200182810382526025815260200180610b3d6025913960400191505060405180910390fd5b6001600160a01b03821661062d57604051600160e51b62461bcd028152600401808060200182810382526023815260200180610a856023913960400191505060405180910390fd5b6001600160a01b038316600090815260208190526040902054610656908263ffffffff6106e516565b6001600160a01b03808516600090815260208190526040808220939093559084168152205461068b908263ffffffff61074516565b6001600160a01b038084166000818152602081815260409182902094909455805185815290519193928716927fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef92918290030190a3505050565b60008282111561073f5760408051600160e51b62461bcd02815260206004820152601e60248201527f536166654d6174683a207375627472616374696f6e206f766572666c6f770000604482015290519081900360640190fd5b50900390565b6000828201838110156107a25760408051600160e51b62461bcd02815260206004820152601b60248201527f536166654d6174683a206164646974696f6e206f766572666c6f770000000000604482015290519081900360640190fd5b9392505050565b6001600160a01b0382166108075760408051600160e51b62461bcd02815260206004820152601f60248201527f45524332303a206d696e7420746f20746865207a65726f206164647265737300604482015290519081900360640190fd5b60025461081a908263ffffffff61074516565b6002556001600160a01b038216600090815260208190526040902054610846908263ffffffff61074516565b6001600160a01b0383166000818152602081815260408083209490945583518581529351929391927fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef9281900390910190a35050565b6108ad60038263ffffffff61099616565b6040516001600160a01b038216907f6ae172837ea30b801fbfcdd4108aa1d5bf8ff775444fd70256b44e6bf3dfc3f690600090a250565b6108f560038263ffffffff610a1a16565b6040516001600160a01b038216907fe94479a9f7e1952cc78f2d6baab678adc1b772d936c6583def489e524cb6669290600090a250565b60006001600160a01b03821661097657604051600160e51b62461bcd028152600401808060200182810382526022815260200180610b1b6022913960400191505060405180910390fd5b506001600160a01b03166000908152602091909152604090205460ff1690565b6109a0828261092c565b156109f55760408051600160e51b62461bcd02815260206004820152601f60248201527f526f6c65733a206163636f756e7420616c72656164792068617320726f6c6500604482015290519081900360640190fd5b6001600160a01b0316600090815260209190915260409020805460ff19166001179055565b610a24828261092c565b610a6257604051600160e51b62461bcd028152600401808060200182810382526021815260200180610afa6021913960400191505060405180910390fd5b6001600160a01b0316600090815260209190915260409020805460ff1916905556fe45524332303a207472616e7366657220746f20746865207a65726f206164647265737345524332303a20617070726f766520746f20746865207a65726f20616464726573734d696e746572526f6c653a2063616c6c657220646f6573206e6f74206861766520746865204d696e74657220726f6c65526f6c65733a206163636f756e7420646f6573206e6f74206861766520726f6c65526f6c65733a206163636f756e7420697320746865207a65726f206164647265737345524332303a207472616e736665722066726f6d20746865207a65726f206164647265737345524332303a20617070726f76652066726f6d20746865207a65726f2061646472657373a165627a7a723058208533d6897cadaf6784525a9d32e613a3c81da44cfdedf624f1bdb497f8b68f4a0029526f6c65733a206163636f756e7420697320746865207a65726f2061646472657373"

// DeployERC20Mintable deploys a new Klaytn contract, binding an instance of ERC20Mintable to it.
func DeployERC20Mintable(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ERC20Mintable, error) {
	parsed, err := abi.JSON(strings.NewReader(ERC20MintableABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(ERC20MintableBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ERC20Mintable{ERC20MintableCaller: ERC20MintableCaller{contract: contract}, ERC20MintableTransactor: ERC20MintableTransactor{contract: contract}, ERC20MintableFilterer: ERC20MintableFilterer{contract: contract}}, nil
}

// ERC20Mintable is an auto generated Go binding around a Klaytn contract.
type ERC20Mintable struct {
	ERC20MintableCaller     // Read-only binding to the contract
	ERC20MintableTransactor // Write-only binding to the contract
	ERC20MintableFilterer   // Log filterer for contract events
}

// ERC20MintableCaller is an auto generated read-only Go binding around a Klaytn contract.
type ERC20MintableCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC20MintableTransactor is an auto generated write-only Go binding around a Klaytn contract.
type ERC20MintableTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC20MintableFilterer is an auto generated log filtering Go binding around a Klaytn contract events.
type ERC20MintableFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC20MintableSession is an auto generated Go binding around a Klaytn contract,
// with pre-set call and transact options.
type ERC20MintableSession struct {
	Contract     *ERC20Mintable    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ERC20MintableCallerSession is an auto generated read-only Go binding around a Klaytn contract,
// with pre-set call options.
type ERC20MintableCallerSession struct {
	Contract *ERC20MintableCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// ERC20MintableTransactorSession is an auto generated write-only Go binding around a Klaytn contract,
// with pre-set transact options.
type ERC20MintableTransactorSession struct {
	Contract     *ERC20MintableTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// ERC20MintableRaw is an auto generated low-level Go binding around a Klaytn contract.
type ERC20MintableRaw struct {
	Contract *ERC20Mintable // Generic contract binding to access the raw methods on
}

// ERC20MintableCallerRaw is an auto generated low-level read-only Go binding around a Klaytn contract.
type ERC20MintableCallerRaw struct {
	Contract *ERC20MintableCaller // Generic read-only contract binding to access the raw methods on
}

// ERC20MintableTransactorRaw is an auto generated low-level write-only Go binding around a Klaytn contract.
type ERC20MintableTransactorRaw struct {
	Contract *ERC20MintableTransactor // Generic write-only contract binding to access the raw methods on
}

// NewERC20Mintable creates a new instance of ERC20Mintable, bound to a specific deployed contract.
func NewERC20Mintable(address common.Address, backend bind.ContractBackend) (*ERC20Mintable, error) {
	contract, err := bindERC20Mintable(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ERC20Mintable{ERC20MintableCaller: ERC20MintableCaller{contract: contract}, ERC20MintableTransactor: ERC20MintableTransactor{contract: contract}, ERC20MintableFilterer: ERC20MintableFilterer{contract: contract}}, nil
}

// NewERC20MintableCaller creates a new read-only instance of ERC20Mintable, bound to a specific deployed contract.
func NewERC20MintableCaller(address common.Address, caller bind.ContractCaller) (*ERC20MintableCaller, error) {
	contract, err := bindERC20Mintable(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ERC20MintableCaller{contract: contract}, nil
}

// NewERC20MintableTransactor creates a new write-only instance of ERC20Mintable, bound to a specific deployed contract.
func NewERC20MintableTransactor(address common.Address, transactor bind.ContractTransactor) (*ERC20MintableTransactor, error) {
	contract, err := bindERC20Mintable(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ERC20MintableTransactor{contract: contract}, nil
}

// NewERC20MintableFilterer creates a new log filterer instance of ERC20Mintable, bound to a specific deployed contract.
func NewERC20MintableFilterer(address common.Address, filterer bind.ContractFilterer) (*ERC20MintableFilterer, error) {
	contract, err := bindERC20Mintable(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ERC20MintableFilterer{contract: contract}, nil
}

// bindERC20Mintable binds a generic wrapper to an already deployed contract.
func bindERC20Mintable(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ERC20MintableABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ERC20Mintable *ERC20MintableRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _ERC20Mintable.Contract.ERC20MintableCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ERC20Mintable *ERC20MintableRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC20Mintable.Contract.ERC20MintableTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ERC20Mintable *ERC20MintableRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ERC20Mintable.Contract.ERC20MintableTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ERC20Mintable *ERC20MintableCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _ERC20Mintable.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ERC20Mintable *ERC20MintableTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC20Mintable.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ERC20Mintable *ERC20MintableTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ERC20Mintable.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_ERC20Mintable *ERC20MintableCaller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ERC20Mintable.contract.Call(opts, out, "allowance", owner, spender)
	return *ret0, err
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_ERC20Mintable *ERC20MintableSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _ERC20Mintable.Contract.Allowance(&_ERC20Mintable.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_ERC20Mintable *ERC20MintableCallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _ERC20Mintable.Contract.Allowance(&_ERC20Mintable.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_ERC20Mintable *ERC20MintableCaller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ERC20Mintable.contract.Call(opts, out, "balanceOf", account)
	return *ret0, err
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_ERC20Mintable *ERC20MintableSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _ERC20Mintable.Contract.BalanceOf(&_ERC20Mintable.CallOpts, account)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_ERC20Mintable *ERC20MintableCallerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _ERC20Mintable.Contract.BalanceOf(&_ERC20Mintable.CallOpts, account)
}

// IsMinter is a free data retrieval call binding the contract method 0xaa271e1a.
//
// Solidity: function isMinter(address account) view returns(bool)
func (_ERC20Mintable *ERC20MintableCaller) IsMinter(opts *bind.CallOpts, account common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _ERC20Mintable.contract.Call(opts, out, "isMinter", account)
	return *ret0, err
}

// IsMinter is a free data retrieval call binding the contract method 0xaa271e1a.
//
// Solidity: function isMinter(address account) view returns(bool)
func (_ERC20Mintable *ERC20MintableSession) IsMinter(account common.Address) (bool, error) {
	return _ERC20Mintable.Contract.IsMinter(&_ERC20Mintable.CallOpts, account)
}

// IsMinter is a free data retrieval call binding the contract method 0xaa271e1a.
//
// Solidity: function isMinter(address account) view returns(bool)
func (_ERC20Mintable *ERC20MintableCallerSession) IsMinter(account common.Address) (bool, error) {
	return _ERC20Mintable.Contract.IsMinter(&_ERC20Mintable.CallOpts, account)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_ERC20Mintable *ERC20MintableCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ERC20Mintable.contract.Call(opts, out, "totalSupply")
	return *ret0, err
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_ERC20Mintable *ERC20MintableSession) TotalSupply() (*big.Int, error) {
	return _ERC20Mintable.Contract.TotalSupply(&_ERC20Mintable.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_ERC20Mintable *ERC20MintableCallerSession) TotalSupply() (*big.Int, error) {
	return _ERC20Mintable.Contract.TotalSupply(&_ERC20Mintable.CallOpts)
}

// AddMinter is a paid mutator transaction binding the contract method 0x983b2d56.
//
// Solidity: function addMinter(address account) returns()
func (_ERC20Mintable *ERC20MintableTransactor) AddMinter(opts *bind.TransactOpts, account common.Address) (*types.Transaction, error) {
	return _ERC20Mintable.contract.Transact(opts, "addMinter", account)
}

// AddMinter is a paid mutator transaction binding the contract method 0x983b2d56.
//
// Solidity: function addMinter(address account) returns()
func (_ERC20Mintable *ERC20MintableSession) AddMinter(account common.Address) (*types.Transaction, error) {
	return _ERC20Mintable.Contract.AddMinter(&_ERC20Mintable.TransactOpts, account)
}

// AddMinter is a paid mutator transaction binding the contract method 0x983b2d56.
//
// Solidity: function addMinter(address account) returns()
func (_ERC20Mintable *ERC20MintableTransactorSession) AddMinter(account common.Address) (*types.Transaction, error) {
	return _ERC20Mintable.Contract.AddMinter(&_ERC20Mintable.TransactOpts, account)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_ERC20Mintable *ERC20MintableTransactor) Approve(opts *bind.TransactOpts, spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _ERC20Mintable.contract.Transact(opts, "approve", spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_ERC20Mintable *ERC20MintableSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _ERC20Mintable.Contract.Approve(&_ERC20Mintable.TransactOpts, spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_ERC20Mintable *ERC20MintableTransactorSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _ERC20Mintable.Contract.Approve(&_ERC20Mintable.TransactOpts, spender, value)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 subtractedValue) returns(bool)
func (_ERC20Mintable *ERC20MintableTransactor) DecreaseAllowance(opts *bind.TransactOpts, spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _ERC20Mintable.contract.Transact(opts, "decreaseAllowance", spender, subtractedValue)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 subtractedValue) returns(bool)
func (_ERC20Mintable *ERC20MintableSession) DecreaseAllowance(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _ERC20Mintable.Contract.DecreaseAllowance(&_ERC20Mintable.TransactOpts, spender, subtractedValue)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 subtractedValue) returns(bool)
func (_ERC20Mintable *ERC20MintableTransactorSession) DecreaseAllowance(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _ERC20Mintable.Contract.DecreaseAllowance(&_ERC20Mintable.TransactOpts, spender, subtractedValue)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 addedValue) returns(bool)
func (_ERC20Mintable *ERC20MintableTransactor) IncreaseAllowance(opts *bind.TransactOpts, spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _ERC20Mintable.contract.Transact(opts, "increaseAllowance", spender, addedValue)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 addedValue) returns(bool)
func (_ERC20Mintable *ERC20MintableSession) IncreaseAllowance(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _ERC20Mintable.Contract.IncreaseAllowance(&_ERC20Mintable.TransactOpts, spender, addedValue)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 addedValue) returns(bool)
func (_ERC20Mintable *ERC20MintableTransactorSession) IncreaseAllowance(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _ERC20Mintable.Contract.IncreaseAllowance(&_ERC20Mintable.TransactOpts, spender, addedValue)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address account, uint256 amount) returns(bool)
func (_ERC20Mintable *ERC20MintableTransactor) Mint(opts *bind.TransactOpts, account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ERC20Mintable.contract.Transact(opts, "mint", account, amount)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address account, uint256 amount) returns(bool)
func (_ERC20Mintable *ERC20MintableSession) Mint(account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ERC20Mintable.Contract.Mint(&_ERC20Mintable.TransactOpts, account, amount)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address account, uint256 amount) returns(bool)
func (_ERC20Mintable *ERC20MintableTransactorSession) Mint(account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ERC20Mintable.Contract.Mint(&_ERC20Mintable.TransactOpts, account, amount)
}

// RenounceMinter is a paid mutator transaction binding the contract method 0x98650275.
//
// Solidity: function renounceMinter() returns()
func (_ERC20Mintable *ERC20MintableTransactor) RenounceMinter(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC20Mintable.contract.Transact(opts, "renounceMinter")
}

// RenounceMinter is a paid mutator transaction binding the contract method 0x98650275.
//
// Solidity: function renounceMinter() returns()
func (_ERC20Mintable *ERC20MintableSession) RenounceMinter() (*types.Transaction, error) {
	return _ERC20Mintable.Contract.RenounceMinter(&_ERC20Mintable.TransactOpts)
}

// RenounceMinter is a paid mutator transaction binding the contract method 0x98650275.
//
// Solidity: function renounceMinter() returns()
func (_ERC20Mintable *ERC20MintableTransactorSession) RenounceMinter() (*types.Transaction, error) {
	return _ERC20Mintable.Contract.RenounceMinter(&_ERC20Mintable.TransactOpts)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_ERC20Mintable *ERC20MintableTransactor) Transfer(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ERC20Mintable.contract.Transact(opts, "transfer", recipient, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_ERC20Mintable *ERC20MintableSession) Transfer(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ERC20Mintable.Contract.Transfer(&_ERC20Mintable.TransactOpts, recipient, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_ERC20Mintable *ERC20MintableTransactorSession) Transfer(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ERC20Mintable.Contract.Transfer(&_ERC20Mintable.TransactOpts, recipient, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_ERC20Mintable *ERC20MintableTransactor) TransferFrom(opts *bind.TransactOpts, sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ERC20Mintable.contract.Transact(opts, "transferFrom", sender, recipient, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_ERC20Mintable *ERC20MintableSession) TransferFrom(sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ERC20Mintable.Contract.TransferFrom(&_ERC20Mintable.TransactOpts, sender, recipient, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_ERC20Mintable *ERC20MintableTransactorSession) TransferFrom(sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ERC20Mintable.Contract.TransferFrom(&_ERC20Mintable.TransactOpts, sender, recipient, amount)
}

// ERC20MintableApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the ERC20Mintable contract.
type ERC20MintableApprovalIterator struct {
	Event *ERC20MintableApproval // Event containing the contract specifics and raw log

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
func (it *ERC20MintableApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC20MintableApproval)
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
		it.Event = new(ERC20MintableApproval)
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
func (it *ERC20MintableApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC20MintableApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC20MintableApproval represents a Approval event raised by the ERC20Mintable contract.
type ERC20MintableApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_ERC20Mintable *ERC20MintableFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*ERC20MintableApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _ERC20Mintable.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &ERC20MintableApprovalIterator{contract: _ERC20Mintable.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_ERC20Mintable *ERC20MintableFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *ERC20MintableApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _ERC20Mintable.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC20MintableApproval)
				if err := _ERC20Mintable.contract.UnpackLog(event, "Approval", log); err != nil {
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

// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_ERC20Mintable *ERC20MintableFilterer) ParseApproval(log types.Log) (*ERC20MintableApproval, error) {
	event := new(ERC20MintableApproval)
	if err := _ERC20Mintable.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ERC20MintableMinterAddedIterator is returned from FilterMinterAdded and is used to iterate over the raw logs and unpacked data for MinterAdded events raised by the ERC20Mintable contract.
type ERC20MintableMinterAddedIterator struct {
	Event *ERC20MintableMinterAdded // Event containing the contract specifics and raw log

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
func (it *ERC20MintableMinterAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC20MintableMinterAdded)
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
		it.Event = new(ERC20MintableMinterAdded)
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
func (it *ERC20MintableMinterAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC20MintableMinterAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC20MintableMinterAdded represents a MinterAdded event raised by the ERC20Mintable contract.
type ERC20MintableMinterAdded struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterMinterAdded is a free log retrieval operation binding the contract event 0x6ae172837ea30b801fbfcdd4108aa1d5bf8ff775444fd70256b44e6bf3dfc3f6.
//
// Solidity: event MinterAdded(address indexed account)
func (_ERC20Mintable *ERC20MintableFilterer) FilterMinterAdded(opts *bind.FilterOpts, account []common.Address) (*ERC20MintableMinterAddedIterator, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _ERC20Mintable.contract.FilterLogs(opts, "MinterAdded", accountRule)
	if err != nil {
		return nil, err
	}
	return &ERC20MintableMinterAddedIterator{contract: _ERC20Mintable.contract, event: "MinterAdded", logs: logs, sub: sub}, nil
}

// WatchMinterAdded is a free log subscription operation binding the contract event 0x6ae172837ea30b801fbfcdd4108aa1d5bf8ff775444fd70256b44e6bf3dfc3f6.
//
// Solidity: event MinterAdded(address indexed account)
func (_ERC20Mintable *ERC20MintableFilterer) WatchMinterAdded(opts *bind.WatchOpts, sink chan<- *ERC20MintableMinterAdded, account []common.Address) (event.Subscription, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _ERC20Mintable.contract.WatchLogs(opts, "MinterAdded", accountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC20MintableMinterAdded)
				if err := _ERC20Mintable.contract.UnpackLog(event, "MinterAdded", log); err != nil {
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

// ParseMinterAdded is a log parse operation binding the contract event 0x6ae172837ea30b801fbfcdd4108aa1d5bf8ff775444fd70256b44e6bf3dfc3f6.
//
// Solidity: event MinterAdded(address indexed account)
func (_ERC20Mintable *ERC20MintableFilterer) ParseMinterAdded(log types.Log) (*ERC20MintableMinterAdded, error) {
	event := new(ERC20MintableMinterAdded)
	if err := _ERC20Mintable.contract.UnpackLog(event, "MinterAdded", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ERC20MintableMinterRemovedIterator is returned from FilterMinterRemoved and is used to iterate over the raw logs and unpacked data for MinterRemoved events raised by the ERC20Mintable contract.
type ERC20MintableMinterRemovedIterator struct {
	Event *ERC20MintableMinterRemoved // Event containing the contract specifics and raw log

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
func (it *ERC20MintableMinterRemovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC20MintableMinterRemoved)
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
		it.Event = new(ERC20MintableMinterRemoved)
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
func (it *ERC20MintableMinterRemovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC20MintableMinterRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC20MintableMinterRemoved represents a MinterRemoved event raised by the ERC20Mintable contract.
type ERC20MintableMinterRemoved struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterMinterRemoved is a free log retrieval operation binding the contract event 0xe94479a9f7e1952cc78f2d6baab678adc1b772d936c6583def489e524cb66692.
//
// Solidity: event MinterRemoved(address indexed account)
func (_ERC20Mintable *ERC20MintableFilterer) FilterMinterRemoved(opts *bind.FilterOpts, account []common.Address) (*ERC20MintableMinterRemovedIterator, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _ERC20Mintable.contract.FilterLogs(opts, "MinterRemoved", accountRule)
	if err != nil {
		return nil, err
	}
	return &ERC20MintableMinterRemovedIterator{contract: _ERC20Mintable.contract, event: "MinterRemoved", logs: logs, sub: sub}, nil
}

// WatchMinterRemoved is a free log subscription operation binding the contract event 0xe94479a9f7e1952cc78f2d6baab678adc1b772d936c6583def489e524cb66692.
//
// Solidity: event MinterRemoved(address indexed account)
func (_ERC20Mintable *ERC20MintableFilterer) WatchMinterRemoved(opts *bind.WatchOpts, sink chan<- *ERC20MintableMinterRemoved, account []common.Address) (event.Subscription, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _ERC20Mintable.contract.WatchLogs(opts, "MinterRemoved", accountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC20MintableMinterRemoved)
				if err := _ERC20Mintable.contract.UnpackLog(event, "MinterRemoved", log); err != nil {
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

// ParseMinterRemoved is a log parse operation binding the contract event 0xe94479a9f7e1952cc78f2d6baab678adc1b772d936c6583def489e524cb66692.
//
// Solidity: event MinterRemoved(address indexed account)
func (_ERC20Mintable *ERC20MintableFilterer) ParseMinterRemoved(log types.Log) (*ERC20MintableMinterRemoved, error) {
	event := new(ERC20MintableMinterRemoved)
	if err := _ERC20Mintable.contract.UnpackLog(event, "MinterRemoved", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ERC20MintableTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the ERC20Mintable contract.
type ERC20MintableTransferIterator struct {
	Event *ERC20MintableTransfer // Event containing the contract specifics and raw log

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
func (it *ERC20MintableTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC20MintableTransfer)
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
		it.Event = new(ERC20MintableTransfer)
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
func (it *ERC20MintableTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC20MintableTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC20MintableTransfer represents a Transfer event raised by the ERC20Mintable contract.
type ERC20MintableTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_ERC20Mintable *ERC20MintableFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*ERC20MintableTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ERC20Mintable.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &ERC20MintableTransferIterator{contract: _ERC20Mintable.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_ERC20Mintable *ERC20MintableFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *ERC20MintableTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ERC20Mintable.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC20MintableTransfer)
				if err := _ERC20Mintable.contract.UnpackLog(event, "Transfer", log); err != nil {
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

// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_ERC20Mintable *ERC20MintableFilterer) ParseTransfer(log types.Log) (*ERC20MintableTransfer, error) {
	event := new(ERC20MintableTransfer)
	if err := _ERC20Mintable.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ERC20ServiceChainABI is the input ABI used to generate the binding from.
const ERC20ServiceChainABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"spender\",\"type\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"sender\",\"type\":\"address\"},{\"name\":\"recipient\",\"type\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"spender\",\"type\":\"address\"},{\"name\":\"addedValue\",\"type\":\"uint256\"}],\"name\":\"increaseAllowance\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_amount\",\"type\":\"uint256\"},{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_feeLimit\",\"type\":\"uint256\"},{\"name\":\"_extraData\",\"type\":\"bytes\"}],\"name\":\"requestValueTransfer\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_bridge\",\"type\":\"address\"}],\"name\":\"setBridge\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"isOwner\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"spender\",\"type\":\"address\"},{\"name\":\"subtractedValue\",\"type\":\"uint256\"}],\"name\":\"decreaseAllowance\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"recipient\",\"type\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"owner\",\"type\":\"address\"},{\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"bridge\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"_bridge\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"}]"

// ERC20ServiceChainBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const ERC20ServiceChainBinRuntime = ``

// ERC20ServiceChainFuncSigs maps the 4-byte function signature to its string representation.
var ERC20ServiceChainFuncSigs = map[string]string{
	"dd62ed3e": "allowance(address,address)",
	"095ea7b3": "approve(address,uint256)",
	"70a08231": "balanceOf(address)",
	"e78cea92": "bridge()",
	"a457c2d7": "decreaseAllowance(address,uint256)",
	"39509351": "increaseAllowance(address,uint256)",
	"8f32d59b": "isOwner()",
	"8da5cb5b": "owner()",
	"715018a6": "renounceOwnership()",
	"3dc3c9e1": "requestValueTransfer(uint256,address,uint256,bytes)",
	"8dd14802": "setBridge(address)",
	"18160ddd": "totalSupply()",
	"a9059cbb": "transfer(address,uint256)",
	"23b872dd": "transferFrom(address,address,uint256)",
	"f2fde38b": "transferOwnership(address)",
}

// ERC20ServiceChain is an auto generated Go binding around a Klaytn contract.
type ERC20ServiceChain struct {
	ERC20ServiceChainCaller     // Read-only binding to the contract
	ERC20ServiceChainTransactor // Write-only binding to the contract
	ERC20ServiceChainFilterer   // Log filterer for contract events
}

// ERC20ServiceChainCaller is an auto generated read-only Go binding around a Klaytn contract.
type ERC20ServiceChainCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC20ServiceChainTransactor is an auto generated write-only Go binding around a Klaytn contract.
type ERC20ServiceChainTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC20ServiceChainFilterer is an auto generated log filtering Go binding around a Klaytn contract events.
type ERC20ServiceChainFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC20ServiceChainSession is an auto generated Go binding around a Klaytn contract,
// with pre-set call and transact options.
type ERC20ServiceChainSession struct {
	Contract     *ERC20ServiceChain // Generic contract binding to set the session for
	CallOpts     bind.CallOpts      // Call options to use throughout this session
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// ERC20ServiceChainCallerSession is an auto generated read-only Go binding around a Klaytn contract,
// with pre-set call options.
type ERC20ServiceChainCallerSession struct {
	Contract *ERC20ServiceChainCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts            // Call options to use throughout this session
}

// ERC20ServiceChainTransactorSession is an auto generated write-only Go binding around a Klaytn contract,
// with pre-set transact options.
type ERC20ServiceChainTransactorSession struct {
	Contract     *ERC20ServiceChainTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// ERC20ServiceChainRaw is an auto generated low-level Go binding around a Klaytn contract.
type ERC20ServiceChainRaw struct {
	Contract *ERC20ServiceChain // Generic contract binding to access the raw methods on
}

// ERC20ServiceChainCallerRaw is an auto generated low-level read-only Go binding around a Klaytn contract.
type ERC20ServiceChainCallerRaw struct {
	Contract *ERC20ServiceChainCaller // Generic read-only contract binding to access the raw methods on
}

// ERC20ServiceChainTransactorRaw is an auto generated low-level write-only Go binding around a Klaytn contract.
type ERC20ServiceChainTransactorRaw struct {
	Contract *ERC20ServiceChainTransactor // Generic write-only contract binding to access the raw methods on
}

// NewERC20ServiceChain creates a new instance of ERC20ServiceChain, bound to a specific deployed contract.
func NewERC20ServiceChain(address common.Address, backend bind.ContractBackend) (*ERC20ServiceChain, error) {
	contract, err := bindERC20ServiceChain(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ERC20ServiceChain{ERC20ServiceChainCaller: ERC20ServiceChainCaller{contract: contract}, ERC20ServiceChainTransactor: ERC20ServiceChainTransactor{contract: contract}, ERC20ServiceChainFilterer: ERC20ServiceChainFilterer{contract: contract}}, nil
}

// NewERC20ServiceChainCaller creates a new read-only instance of ERC20ServiceChain, bound to a specific deployed contract.
func NewERC20ServiceChainCaller(address common.Address, caller bind.ContractCaller) (*ERC20ServiceChainCaller, error) {
	contract, err := bindERC20ServiceChain(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ERC20ServiceChainCaller{contract: contract}, nil
}

// NewERC20ServiceChainTransactor creates a new write-only instance of ERC20ServiceChain, bound to a specific deployed contract.
func NewERC20ServiceChainTransactor(address common.Address, transactor bind.ContractTransactor) (*ERC20ServiceChainTransactor, error) {
	contract, err := bindERC20ServiceChain(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ERC20ServiceChainTransactor{contract: contract}, nil
}

// NewERC20ServiceChainFilterer creates a new log filterer instance of ERC20ServiceChain, bound to a specific deployed contract.
func NewERC20ServiceChainFilterer(address common.Address, filterer bind.ContractFilterer) (*ERC20ServiceChainFilterer, error) {
	contract, err := bindERC20ServiceChain(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ERC20ServiceChainFilterer{contract: contract}, nil
}

// bindERC20ServiceChain binds a generic wrapper to an already deployed contract.
func bindERC20ServiceChain(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ERC20ServiceChainABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ERC20ServiceChain *ERC20ServiceChainRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _ERC20ServiceChain.Contract.ERC20ServiceChainCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ERC20ServiceChain *ERC20ServiceChainRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC20ServiceChain.Contract.ERC20ServiceChainTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ERC20ServiceChain *ERC20ServiceChainRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ERC20ServiceChain.Contract.ERC20ServiceChainTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ERC20ServiceChain *ERC20ServiceChainCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _ERC20ServiceChain.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ERC20ServiceChain *ERC20ServiceChainTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC20ServiceChain.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ERC20ServiceChain *ERC20ServiceChainTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ERC20ServiceChain.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_ERC20ServiceChain *ERC20ServiceChainCaller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ERC20ServiceChain.contract.Call(opts, out, "allowance", owner, spender)
	return *ret0, err
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_ERC20ServiceChain *ERC20ServiceChainSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _ERC20ServiceChain.Contract.Allowance(&_ERC20ServiceChain.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_ERC20ServiceChain *ERC20ServiceChainCallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _ERC20ServiceChain.Contract.Allowance(&_ERC20ServiceChain.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_ERC20ServiceChain *ERC20ServiceChainCaller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ERC20ServiceChain.contract.Call(opts, out, "balanceOf", account)
	return *ret0, err
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_ERC20ServiceChain *ERC20ServiceChainSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _ERC20ServiceChain.Contract.BalanceOf(&_ERC20ServiceChain.CallOpts, account)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_ERC20ServiceChain *ERC20ServiceChainCallerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _ERC20ServiceChain.Contract.BalanceOf(&_ERC20ServiceChain.CallOpts, account)
}

// Bridge is a free data retrieval call binding the contract method 0xe78cea92.
//
// Solidity: function bridge() view returns(address)
func (_ERC20ServiceChain *ERC20ServiceChainCaller) Bridge(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _ERC20ServiceChain.contract.Call(opts, out, "bridge")
	return *ret0, err
}

// Bridge is a free data retrieval call binding the contract method 0xe78cea92.
//
// Solidity: function bridge() view returns(address)
func (_ERC20ServiceChain *ERC20ServiceChainSession) Bridge() (common.Address, error) {
	return _ERC20ServiceChain.Contract.Bridge(&_ERC20ServiceChain.CallOpts)
}

// Bridge is a free data retrieval call binding the contract method 0xe78cea92.
//
// Solidity: function bridge() view returns(address)
func (_ERC20ServiceChain *ERC20ServiceChainCallerSession) Bridge() (common.Address, error) {
	return _ERC20ServiceChain.Contract.Bridge(&_ERC20ServiceChain.CallOpts)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() view returns(bool)
func (_ERC20ServiceChain *ERC20ServiceChainCaller) IsOwner(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _ERC20ServiceChain.contract.Call(opts, out, "isOwner")
	return *ret0, err
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() view returns(bool)
func (_ERC20ServiceChain *ERC20ServiceChainSession) IsOwner() (bool, error) {
	return _ERC20ServiceChain.Contract.IsOwner(&_ERC20ServiceChain.CallOpts)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() view returns(bool)
func (_ERC20ServiceChain *ERC20ServiceChainCallerSession) IsOwner() (bool, error) {
	return _ERC20ServiceChain.Contract.IsOwner(&_ERC20ServiceChain.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ERC20ServiceChain *ERC20ServiceChainCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _ERC20ServiceChain.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ERC20ServiceChain *ERC20ServiceChainSession) Owner() (common.Address, error) {
	return _ERC20ServiceChain.Contract.Owner(&_ERC20ServiceChain.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ERC20ServiceChain *ERC20ServiceChainCallerSession) Owner() (common.Address, error) {
	return _ERC20ServiceChain.Contract.Owner(&_ERC20ServiceChain.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_ERC20ServiceChain *ERC20ServiceChainCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ERC20ServiceChain.contract.Call(opts, out, "totalSupply")
	return *ret0, err
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_ERC20ServiceChain *ERC20ServiceChainSession) TotalSupply() (*big.Int, error) {
	return _ERC20ServiceChain.Contract.TotalSupply(&_ERC20ServiceChain.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_ERC20ServiceChain *ERC20ServiceChainCallerSession) TotalSupply() (*big.Int, error) {
	return _ERC20ServiceChain.Contract.TotalSupply(&_ERC20ServiceChain.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_ERC20ServiceChain *ERC20ServiceChainTransactor) Approve(opts *bind.TransactOpts, spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _ERC20ServiceChain.contract.Transact(opts, "approve", spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_ERC20ServiceChain *ERC20ServiceChainSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _ERC20ServiceChain.Contract.Approve(&_ERC20ServiceChain.TransactOpts, spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_ERC20ServiceChain *ERC20ServiceChainTransactorSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _ERC20ServiceChain.Contract.Approve(&_ERC20ServiceChain.TransactOpts, spender, value)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 subtractedValue) returns(bool)
func (_ERC20ServiceChain *ERC20ServiceChainTransactor) DecreaseAllowance(opts *bind.TransactOpts, spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _ERC20ServiceChain.contract.Transact(opts, "decreaseAllowance", spender, subtractedValue)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 subtractedValue) returns(bool)
func (_ERC20ServiceChain *ERC20ServiceChainSession) DecreaseAllowance(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _ERC20ServiceChain.Contract.DecreaseAllowance(&_ERC20ServiceChain.TransactOpts, spender, subtractedValue)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 subtractedValue) returns(bool)
func (_ERC20ServiceChain *ERC20ServiceChainTransactorSession) DecreaseAllowance(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _ERC20ServiceChain.Contract.DecreaseAllowance(&_ERC20ServiceChain.TransactOpts, spender, subtractedValue)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 addedValue) returns(bool)
func (_ERC20ServiceChain *ERC20ServiceChainTransactor) IncreaseAllowance(opts *bind.TransactOpts, spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _ERC20ServiceChain.contract.Transact(opts, "increaseAllowance", spender, addedValue)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 addedValue) returns(bool)
func (_ERC20ServiceChain *ERC20ServiceChainSession) IncreaseAllowance(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _ERC20ServiceChain.Contract.IncreaseAllowance(&_ERC20ServiceChain.TransactOpts, spender, addedValue)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 addedValue) returns(bool)
func (_ERC20ServiceChain *ERC20ServiceChainTransactorSession) IncreaseAllowance(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _ERC20ServiceChain.Contract.IncreaseAllowance(&_ERC20ServiceChain.TransactOpts, spender, addedValue)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ERC20ServiceChain *ERC20ServiceChainTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC20ServiceChain.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ERC20ServiceChain *ERC20ServiceChainSession) RenounceOwnership() (*types.Transaction, error) {
	return _ERC20ServiceChain.Contract.RenounceOwnership(&_ERC20ServiceChain.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ERC20ServiceChain *ERC20ServiceChainTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _ERC20ServiceChain.Contract.RenounceOwnership(&_ERC20ServiceChain.TransactOpts)
}

// RequestValueTransfer is a paid mutator transaction binding the contract method 0x3dc3c9e1.
//
// Solidity: function requestValueTransfer(uint256 _amount, address _to, uint256 _feeLimit, bytes _extraData) returns()
func (_ERC20ServiceChain *ERC20ServiceChainTransactor) RequestValueTransfer(opts *bind.TransactOpts, _amount *big.Int, _to common.Address, _feeLimit *big.Int, _extraData []byte) (*types.Transaction, error) {
	return _ERC20ServiceChain.contract.Transact(opts, "requestValueTransfer", _amount, _to, _feeLimit, _extraData)
}

// RequestValueTransfer is a paid mutator transaction binding the contract method 0x3dc3c9e1.
//
// Solidity: function requestValueTransfer(uint256 _amount, address _to, uint256 _feeLimit, bytes _extraData) returns()
func (_ERC20ServiceChain *ERC20ServiceChainSession) RequestValueTransfer(_amount *big.Int, _to common.Address, _feeLimit *big.Int, _extraData []byte) (*types.Transaction, error) {
	return _ERC20ServiceChain.Contract.RequestValueTransfer(&_ERC20ServiceChain.TransactOpts, _amount, _to, _feeLimit, _extraData)
}

// RequestValueTransfer is a paid mutator transaction binding the contract method 0x3dc3c9e1.
//
// Solidity: function requestValueTransfer(uint256 _amount, address _to, uint256 _feeLimit, bytes _extraData) returns()
func (_ERC20ServiceChain *ERC20ServiceChainTransactorSession) RequestValueTransfer(_amount *big.Int, _to common.Address, _feeLimit *big.Int, _extraData []byte) (*types.Transaction, error) {
	return _ERC20ServiceChain.Contract.RequestValueTransfer(&_ERC20ServiceChain.TransactOpts, _amount, _to, _feeLimit, _extraData)
}

// SetBridge is a paid mutator transaction binding the contract method 0x8dd14802.
//
// Solidity: function setBridge(address _bridge) returns()
func (_ERC20ServiceChain *ERC20ServiceChainTransactor) SetBridge(opts *bind.TransactOpts, _bridge common.Address) (*types.Transaction, error) {
	return _ERC20ServiceChain.contract.Transact(opts, "setBridge", _bridge)
}

// SetBridge is a paid mutator transaction binding the contract method 0x8dd14802.
//
// Solidity: function setBridge(address _bridge) returns()
func (_ERC20ServiceChain *ERC20ServiceChainSession) SetBridge(_bridge common.Address) (*types.Transaction, error) {
	return _ERC20ServiceChain.Contract.SetBridge(&_ERC20ServiceChain.TransactOpts, _bridge)
}

// SetBridge is a paid mutator transaction binding the contract method 0x8dd14802.
//
// Solidity: function setBridge(address _bridge) returns()
func (_ERC20ServiceChain *ERC20ServiceChainTransactorSession) SetBridge(_bridge common.Address) (*types.Transaction, error) {
	return _ERC20ServiceChain.Contract.SetBridge(&_ERC20ServiceChain.TransactOpts, _bridge)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_ERC20ServiceChain *ERC20ServiceChainTransactor) Transfer(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ERC20ServiceChain.contract.Transact(opts, "transfer", recipient, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_ERC20ServiceChain *ERC20ServiceChainSession) Transfer(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ERC20ServiceChain.Contract.Transfer(&_ERC20ServiceChain.TransactOpts, recipient, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_ERC20ServiceChain *ERC20ServiceChainTransactorSession) Transfer(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ERC20ServiceChain.Contract.Transfer(&_ERC20ServiceChain.TransactOpts, recipient, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_ERC20ServiceChain *ERC20ServiceChainTransactor) TransferFrom(opts *bind.TransactOpts, sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ERC20ServiceChain.contract.Transact(opts, "transferFrom", sender, recipient, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_ERC20ServiceChain *ERC20ServiceChainSession) TransferFrom(sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ERC20ServiceChain.Contract.TransferFrom(&_ERC20ServiceChain.TransactOpts, sender, recipient, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_ERC20ServiceChain *ERC20ServiceChainTransactorSession) TransferFrom(sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ERC20ServiceChain.Contract.TransferFrom(&_ERC20ServiceChain.TransactOpts, sender, recipient, amount)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ERC20ServiceChain *ERC20ServiceChainTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _ERC20ServiceChain.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ERC20ServiceChain *ERC20ServiceChainSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ERC20ServiceChain.Contract.TransferOwnership(&_ERC20ServiceChain.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ERC20ServiceChain *ERC20ServiceChainTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ERC20ServiceChain.Contract.TransferOwnership(&_ERC20ServiceChain.TransactOpts, newOwner)
}

// ERC20ServiceChainApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the ERC20ServiceChain contract.
type ERC20ServiceChainApprovalIterator struct {
	Event *ERC20ServiceChainApproval // Event containing the contract specifics and raw log

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
func (it *ERC20ServiceChainApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC20ServiceChainApproval)
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
		it.Event = new(ERC20ServiceChainApproval)
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
func (it *ERC20ServiceChainApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC20ServiceChainApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC20ServiceChainApproval represents a Approval event raised by the ERC20ServiceChain contract.
type ERC20ServiceChainApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_ERC20ServiceChain *ERC20ServiceChainFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*ERC20ServiceChainApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _ERC20ServiceChain.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &ERC20ServiceChainApprovalIterator{contract: _ERC20ServiceChain.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_ERC20ServiceChain *ERC20ServiceChainFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *ERC20ServiceChainApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _ERC20ServiceChain.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC20ServiceChainApproval)
				if err := _ERC20ServiceChain.contract.UnpackLog(event, "Approval", log); err != nil {
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

// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_ERC20ServiceChain *ERC20ServiceChainFilterer) ParseApproval(log types.Log) (*ERC20ServiceChainApproval, error) {
	event := new(ERC20ServiceChainApproval)
	if err := _ERC20ServiceChain.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ERC20ServiceChainOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the ERC20ServiceChain contract.
type ERC20ServiceChainOwnershipTransferredIterator struct {
	Event *ERC20ServiceChainOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *ERC20ServiceChainOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC20ServiceChainOwnershipTransferred)
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
		it.Event = new(ERC20ServiceChainOwnershipTransferred)
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
func (it *ERC20ServiceChainOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC20ServiceChainOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC20ServiceChainOwnershipTransferred represents a OwnershipTransferred event raised by the ERC20ServiceChain contract.
type ERC20ServiceChainOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ERC20ServiceChain *ERC20ServiceChainFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*ERC20ServiceChainOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ERC20ServiceChain.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &ERC20ServiceChainOwnershipTransferredIterator{contract: _ERC20ServiceChain.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ERC20ServiceChain *ERC20ServiceChainFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ERC20ServiceChainOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ERC20ServiceChain.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC20ServiceChainOwnershipTransferred)
				if err := _ERC20ServiceChain.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_ERC20ServiceChain *ERC20ServiceChainFilterer) ParseOwnershipTransferred(log types.Log) (*ERC20ServiceChainOwnershipTransferred, error) {
	event := new(ERC20ServiceChainOwnershipTransferred)
	if err := _ERC20ServiceChain.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ERC20ServiceChainTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the ERC20ServiceChain contract.
type ERC20ServiceChainTransferIterator struct {
	Event *ERC20ServiceChainTransfer // Event containing the contract specifics and raw log

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
func (it *ERC20ServiceChainTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC20ServiceChainTransfer)
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
		it.Event = new(ERC20ServiceChainTransfer)
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
func (it *ERC20ServiceChainTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC20ServiceChainTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC20ServiceChainTransfer represents a Transfer event raised by the ERC20ServiceChain contract.
type ERC20ServiceChainTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_ERC20ServiceChain *ERC20ServiceChainFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*ERC20ServiceChainTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ERC20ServiceChain.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &ERC20ServiceChainTransferIterator{contract: _ERC20ServiceChain.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_ERC20ServiceChain *ERC20ServiceChainFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *ERC20ServiceChainTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ERC20ServiceChain.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC20ServiceChainTransfer)
				if err := _ERC20ServiceChain.contract.UnpackLog(event, "Transfer", log); err != nil {
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

// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_ERC20ServiceChain *ERC20ServiceChainFilterer) ParseTransfer(log types.Log) (*ERC20ServiceChainTransfer, error) {
	event := new(ERC20ServiceChainTransfer)
	if err := _ERC20ServiceChain.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	return event, nil
}

// IERC20ABI is the input ABI used to generate the binding from.
const IERC20ABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"spender\",\"type\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"sender\",\"type\":\"address\"},{\"name\":\"recipient\",\"type\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"recipient\",\"type\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"owner\",\"type\":\"address\"},{\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"}]"

// IERC20BinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const IERC20BinRuntime = ``

// IERC20FuncSigs maps the 4-byte function signature to its string representation.
var IERC20FuncSigs = map[string]string{
	"dd62ed3e": "allowance(address,address)",
	"095ea7b3": "approve(address,uint256)",
	"70a08231": "balanceOf(address)",
	"18160ddd": "totalSupply()",
	"a9059cbb": "transfer(address,uint256)",
	"23b872dd": "transferFrom(address,address,uint256)",
}

// IERC20 is an auto generated Go binding around a Klaytn contract.
type IERC20 struct {
	IERC20Caller     // Read-only binding to the contract
	IERC20Transactor // Write-only binding to the contract
	IERC20Filterer   // Log filterer for contract events
}

// IERC20Caller is an auto generated read-only Go binding around a Klaytn contract.
type IERC20Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IERC20Transactor is an auto generated write-only Go binding around a Klaytn contract.
type IERC20Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IERC20Filterer is an auto generated log filtering Go binding around a Klaytn contract events.
type IERC20Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IERC20Session is an auto generated Go binding around a Klaytn contract,
// with pre-set call and transact options.
type IERC20Session struct {
	Contract     *IERC20           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IERC20CallerSession is an auto generated read-only Go binding around a Klaytn contract,
// with pre-set call options.
type IERC20CallerSession struct {
	Contract *IERC20Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// IERC20TransactorSession is an auto generated write-only Go binding around a Klaytn contract,
// with pre-set transact options.
type IERC20TransactorSession struct {
	Contract     *IERC20Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IERC20Raw is an auto generated low-level Go binding around a Klaytn contract.
type IERC20Raw struct {
	Contract *IERC20 // Generic contract binding to access the raw methods on
}

// IERC20CallerRaw is an auto generated low-level read-only Go binding around a Klaytn contract.
type IERC20CallerRaw struct {
	Contract *IERC20Caller // Generic read-only contract binding to access the raw methods on
}

// IERC20TransactorRaw is an auto generated low-level write-only Go binding around a Klaytn contract.
type IERC20TransactorRaw struct {
	Contract *IERC20Transactor // Generic write-only contract binding to access the raw methods on
}

// NewIERC20 creates a new instance of IERC20, bound to a specific deployed contract.
func NewIERC20(address common.Address, backend bind.ContractBackend) (*IERC20, error) {
	contract, err := bindIERC20(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IERC20{IERC20Caller: IERC20Caller{contract: contract}, IERC20Transactor: IERC20Transactor{contract: contract}, IERC20Filterer: IERC20Filterer{contract: contract}}, nil
}

// NewIERC20Caller creates a new read-only instance of IERC20, bound to a specific deployed contract.
func NewIERC20Caller(address common.Address, caller bind.ContractCaller) (*IERC20Caller, error) {
	contract, err := bindIERC20(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IERC20Caller{contract: contract}, nil
}

// NewIERC20Transactor creates a new write-only instance of IERC20, bound to a specific deployed contract.
func NewIERC20Transactor(address common.Address, transactor bind.ContractTransactor) (*IERC20Transactor, error) {
	contract, err := bindIERC20(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IERC20Transactor{contract: contract}, nil
}

// NewIERC20Filterer creates a new log filterer instance of IERC20, bound to a specific deployed contract.
func NewIERC20Filterer(address common.Address, filterer bind.ContractFilterer) (*IERC20Filterer, error) {
	contract, err := bindIERC20(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IERC20Filterer{contract: contract}, nil
}

// bindIERC20 binds a generic wrapper to an already deployed contract.
func bindIERC20(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(IERC20ABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IERC20 *IERC20Raw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _IERC20.Contract.IERC20Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IERC20 *IERC20Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IERC20.Contract.IERC20Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IERC20 *IERC20Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IERC20.Contract.IERC20Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IERC20 *IERC20CallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _IERC20.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IERC20 *IERC20TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IERC20.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IERC20 *IERC20TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IERC20.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_IERC20 *IERC20Caller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _IERC20.contract.Call(opts, out, "allowance", owner, spender)
	return *ret0, err
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_IERC20 *IERC20Session) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _IERC20.Contract.Allowance(&_IERC20.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_IERC20 *IERC20CallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _IERC20.Contract.Allowance(&_IERC20.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_IERC20 *IERC20Caller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _IERC20.contract.Call(opts, out, "balanceOf", account)
	return *ret0, err
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_IERC20 *IERC20Session) BalanceOf(account common.Address) (*big.Int, error) {
	return _IERC20.Contract.BalanceOf(&_IERC20.CallOpts, account)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_IERC20 *IERC20CallerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _IERC20.Contract.BalanceOf(&_IERC20.CallOpts, account)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_IERC20 *IERC20Caller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _IERC20.contract.Call(opts, out, "totalSupply")
	return *ret0, err
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_IERC20 *IERC20Session) TotalSupply() (*big.Int, error) {
	return _IERC20.Contract.TotalSupply(&_IERC20.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_IERC20 *IERC20CallerSession) TotalSupply() (*big.Int, error) {
	return _IERC20.Contract.TotalSupply(&_IERC20.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_IERC20 *IERC20Transactor) Approve(opts *bind.TransactOpts, spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _IERC20.contract.Transact(opts, "approve", spender, amount)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_IERC20 *IERC20Session) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _IERC20.Contract.Approve(&_IERC20.TransactOpts, spender, amount)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_IERC20 *IERC20TransactorSession) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _IERC20.Contract.Approve(&_IERC20.TransactOpts, spender, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_IERC20 *IERC20Transactor) Transfer(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _IERC20.contract.Transact(opts, "transfer", recipient, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_IERC20 *IERC20Session) Transfer(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _IERC20.Contract.Transfer(&_IERC20.TransactOpts, recipient, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_IERC20 *IERC20TransactorSession) Transfer(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _IERC20.Contract.Transfer(&_IERC20.TransactOpts, recipient, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_IERC20 *IERC20Transactor) TransferFrom(opts *bind.TransactOpts, sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _IERC20.contract.Transact(opts, "transferFrom", sender, recipient, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_IERC20 *IERC20Session) TransferFrom(sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _IERC20.Contract.TransferFrom(&_IERC20.TransactOpts, sender, recipient, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_IERC20 *IERC20TransactorSession) TransferFrom(sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _IERC20.Contract.TransferFrom(&_IERC20.TransactOpts, sender, recipient, amount)
}

// IERC20ApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the IERC20 contract.
type IERC20ApprovalIterator struct {
	Event *IERC20Approval // Event containing the contract specifics and raw log

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
func (it *IERC20ApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IERC20Approval)
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
		it.Event = new(IERC20Approval)
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
func (it *IERC20ApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IERC20ApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IERC20Approval represents a Approval event raised by the IERC20 contract.
type IERC20Approval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_IERC20 *IERC20Filterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*IERC20ApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _IERC20.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &IERC20ApprovalIterator{contract: _IERC20.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_IERC20 *IERC20Filterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *IERC20Approval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _IERC20.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IERC20Approval)
				if err := _IERC20.contract.UnpackLog(event, "Approval", log); err != nil {
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

// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_IERC20 *IERC20Filterer) ParseApproval(log types.Log) (*IERC20Approval, error) {
	event := new(IERC20Approval)
	if err := _IERC20.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	return event, nil
}

// IERC20TransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the IERC20 contract.
type IERC20TransferIterator struct {
	Event *IERC20Transfer // Event containing the contract specifics and raw log

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
func (it *IERC20TransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IERC20Transfer)
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
		it.Event = new(IERC20Transfer)
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
func (it *IERC20TransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IERC20TransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IERC20Transfer represents a Transfer event raised by the IERC20 contract.
type IERC20Transfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_IERC20 *IERC20Filterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*IERC20TransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _IERC20.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &IERC20TransferIterator{contract: _IERC20.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_IERC20 *IERC20Filterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *IERC20Transfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _IERC20.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IERC20Transfer)
				if err := _IERC20.contract.UnpackLog(event, "Transfer", log); err != nil {
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

// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_IERC20 *IERC20Filterer) ParseTransfer(log types.Log) (*IERC20Transfer, error) {
	event := new(IERC20Transfer)
	if err := _IERC20.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	return event, nil
}

// IERC20BridgeReceiverABI is the input ABI used to generate the binding from.
const IERC20BridgeReceiverABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"_from\",\"type\":\"address\"},{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_amount\",\"type\":\"uint256\"},{\"name\":\"_feeLimit\",\"type\":\"uint256\"},{\"name\":\"_extraData\",\"type\":\"bytes\"}],\"name\":\"onERC20Received\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// IERC20BridgeReceiverBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const IERC20BridgeReceiverBinRuntime = ``

// IERC20BridgeReceiverFuncSigs maps the 4-byte function signature to its string representation.
var IERC20BridgeReceiverFuncSigs = map[string]string{
	"f1656e53": "onERC20Received(address,address,uint256,uint256,bytes)",
}

// IERC20BridgeReceiver is an auto generated Go binding around a Klaytn contract.
type IERC20BridgeReceiver struct {
	IERC20BridgeReceiverCaller     // Read-only binding to the contract
	IERC20BridgeReceiverTransactor // Write-only binding to the contract
	IERC20BridgeReceiverFilterer   // Log filterer for contract events
}

// IERC20BridgeReceiverCaller is an auto generated read-only Go binding around a Klaytn contract.
type IERC20BridgeReceiverCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IERC20BridgeReceiverTransactor is an auto generated write-only Go binding around a Klaytn contract.
type IERC20BridgeReceiverTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IERC20BridgeReceiverFilterer is an auto generated log filtering Go binding around a Klaytn contract events.
type IERC20BridgeReceiverFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IERC20BridgeReceiverSession is an auto generated Go binding around a Klaytn contract,
// with pre-set call and transact options.
type IERC20BridgeReceiverSession struct {
	Contract     *IERC20BridgeReceiver // Generic contract binding to set the session for
	CallOpts     bind.CallOpts         // Call options to use throughout this session
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// IERC20BridgeReceiverCallerSession is an auto generated read-only Go binding around a Klaytn contract,
// with pre-set call options.
type IERC20BridgeReceiverCallerSession struct {
	Contract *IERC20BridgeReceiverCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts               // Call options to use throughout this session
}

// IERC20BridgeReceiverTransactorSession is an auto generated write-only Go binding around a Klaytn contract,
// with pre-set transact options.
type IERC20BridgeReceiverTransactorSession struct {
	Contract     *IERC20BridgeReceiverTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts               // Transaction auth options to use throughout this session
}

// IERC20BridgeReceiverRaw is an auto generated low-level Go binding around a Klaytn contract.
type IERC20BridgeReceiverRaw struct {
	Contract *IERC20BridgeReceiver // Generic contract binding to access the raw methods on
}

// IERC20BridgeReceiverCallerRaw is an auto generated low-level read-only Go binding around a Klaytn contract.
type IERC20BridgeReceiverCallerRaw struct {
	Contract *IERC20BridgeReceiverCaller // Generic read-only contract binding to access the raw methods on
}

// IERC20BridgeReceiverTransactorRaw is an auto generated low-level write-only Go binding around a Klaytn contract.
type IERC20BridgeReceiverTransactorRaw struct {
	Contract *IERC20BridgeReceiverTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIERC20BridgeReceiver creates a new instance of IERC20BridgeReceiver, bound to a specific deployed contract.
func NewIERC20BridgeReceiver(address common.Address, backend bind.ContractBackend) (*IERC20BridgeReceiver, error) {
	contract, err := bindIERC20BridgeReceiver(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IERC20BridgeReceiver{IERC20BridgeReceiverCaller: IERC20BridgeReceiverCaller{contract: contract}, IERC20BridgeReceiverTransactor: IERC20BridgeReceiverTransactor{contract: contract}, IERC20BridgeReceiverFilterer: IERC20BridgeReceiverFilterer{contract: contract}}, nil
}

// NewIERC20BridgeReceiverCaller creates a new read-only instance of IERC20BridgeReceiver, bound to a specific deployed contract.
func NewIERC20BridgeReceiverCaller(address common.Address, caller bind.ContractCaller) (*IERC20BridgeReceiverCaller, error) {
	contract, err := bindIERC20BridgeReceiver(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IERC20BridgeReceiverCaller{contract: contract}, nil
}

// NewIERC20BridgeReceiverTransactor creates a new write-only instance of IERC20BridgeReceiver, bound to a specific deployed contract.
func NewIERC20BridgeReceiverTransactor(address common.Address, transactor bind.ContractTransactor) (*IERC20BridgeReceiverTransactor, error) {
	contract, err := bindIERC20BridgeReceiver(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IERC20BridgeReceiverTransactor{contract: contract}, nil
}

// NewIERC20BridgeReceiverFilterer creates a new log filterer instance of IERC20BridgeReceiver, bound to a specific deployed contract.
func NewIERC20BridgeReceiverFilterer(address common.Address, filterer bind.ContractFilterer) (*IERC20BridgeReceiverFilterer, error) {
	contract, err := bindIERC20BridgeReceiver(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IERC20BridgeReceiverFilterer{contract: contract}, nil
}

// bindIERC20BridgeReceiver binds a generic wrapper to an already deployed contract.
func bindIERC20BridgeReceiver(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(IERC20BridgeReceiverABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IERC20BridgeReceiver *IERC20BridgeReceiverRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _IERC20BridgeReceiver.Contract.IERC20BridgeReceiverCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IERC20BridgeReceiver *IERC20BridgeReceiverRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IERC20BridgeReceiver.Contract.IERC20BridgeReceiverTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IERC20BridgeReceiver *IERC20BridgeReceiverRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IERC20BridgeReceiver.Contract.IERC20BridgeReceiverTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IERC20BridgeReceiver *IERC20BridgeReceiverCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _IERC20BridgeReceiver.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IERC20BridgeReceiver *IERC20BridgeReceiverTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IERC20BridgeReceiver.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IERC20BridgeReceiver *IERC20BridgeReceiverTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IERC20BridgeReceiver.Contract.contract.Transact(opts, method, params...)
}

// OnERC20Received is a paid mutator transaction binding the contract method 0xf1656e53.
//
// Solidity: function onERC20Received(address _from, address _to, uint256 _amount, uint256 _feeLimit, bytes _extraData) returns()
func (_IERC20BridgeReceiver *IERC20BridgeReceiverTransactor) OnERC20Received(opts *bind.TransactOpts, _from common.Address, _to common.Address, _amount *big.Int, _feeLimit *big.Int, _extraData []byte) (*types.Transaction, error) {
	return _IERC20BridgeReceiver.contract.Transact(opts, "onERC20Received", _from, _to, _amount, _feeLimit, _extraData)
}

// OnERC20Received is a paid mutator transaction binding the contract method 0xf1656e53.
//
// Solidity: function onERC20Received(address _from, address _to, uint256 _amount, uint256 _feeLimit, bytes _extraData) returns()
func (_IERC20BridgeReceiver *IERC20BridgeReceiverSession) OnERC20Received(_from common.Address, _to common.Address, _amount *big.Int, _feeLimit *big.Int, _extraData []byte) (*types.Transaction, error) {
	return _IERC20BridgeReceiver.Contract.OnERC20Received(&_IERC20BridgeReceiver.TransactOpts, _from, _to, _amount, _feeLimit, _extraData)
}

// OnERC20Received is a paid mutator transaction binding the contract method 0xf1656e53.
//
// Solidity: function onERC20Received(address _from, address _to, uint256 _amount, uint256 _feeLimit, bytes _extraData) returns()
func (_IERC20BridgeReceiver *IERC20BridgeReceiverTransactorSession) OnERC20Received(_from common.Address, _to common.Address, _amount *big.Int, _feeLimit *big.Int, _extraData []byte) (*types.Transaction, error) {
	return _IERC20BridgeReceiver.Contract.OnERC20Received(&_IERC20BridgeReceiver.TransactOpts, _from, _to, _amount, _feeLimit, _extraData)
}

// MinterRoleABI is the input ABI used to generate the binding from.
const MinterRoleABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"account\",\"type\":\"address\"}],\"name\":\"addMinter\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"renounceMinter\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"account\",\"type\":\"address\"}],\"name\":\"isMinter\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"account\",\"type\":\"address\"}],\"name\":\"MinterAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"account\",\"type\":\"address\"}],\"name\":\"MinterRemoved\",\"type\":\"event\"}]"

// MinterRoleBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const MinterRoleBinRuntime = ``

// MinterRoleFuncSigs maps the 4-byte function signature to its string representation.
var MinterRoleFuncSigs = map[string]string{
	"983b2d56": "addMinter(address)",
	"aa271e1a": "isMinter(address)",
	"98650275": "renounceMinter()",
}

// MinterRole is an auto generated Go binding around a Klaytn contract.
type MinterRole struct {
	MinterRoleCaller     // Read-only binding to the contract
	MinterRoleTransactor // Write-only binding to the contract
	MinterRoleFilterer   // Log filterer for contract events
}

// MinterRoleCaller is an auto generated read-only Go binding around a Klaytn contract.
type MinterRoleCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MinterRoleTransactor is an auto generated write-only Go binding around a Klaytn contract.
type MinterRoleTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MinterRoleFilterer is an auto generated log filtering Go binding around a Klaytn contract events.
type MinterRoleFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MinterRoleSession is an auto generated Go binding around a Klaytn contract,
// with pre-set call and transact options.
type MinterRoleSession struct {
	Contract     *MinterRole       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// MinterRoleCallerSession is an auto generated read-only Go binding around a Klaytn contract,
// with pre-set call options.
type MinterRoleCallerSession struct {
	Contract *MinterRoleCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// MinterRoleTransactorSession is an auto generated write-only Go binding around a Klaytn contract,
// with pre-set transact options.
type MinterRoleTransactorSession struct {
	Contract     *MinterRoleTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// MinterRoleRaw is an auto generated low-level Go binding around a Klaytn contract.
type MinterRoleRaw struct {
	Contract *MinterRole // Generic contract binding to access the raw methods on
}

// MinterRoleCallerRaw is an auto generated low-level read-only Go binding around a Klaytn contract.
type MinterRoleCallerRaw struct {
	Contract *MinterRoleCaller // Generic read-only contract binding to access the raw methods on
}

// MinterRoleTransactorRaw is an auto generated low-level write-only Go binding around a Klaytn contract.
type MinterRoleTransactorRaw struct {
	Contract *MinterRoleTransactor // Generic write-only contract binding to access the raw methods on
}

// NewMinterRole creates a new instance of MinterRole, bound to a specific deployed contract.
func NewMinterRole(address common.Address, backend bind.ContractBackend) (*MinterRole, error) {
	contract, err := bindMinterRole(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &MinterRole{MinterRoleCaller: MinterRoleCaller{contract: contract}, MinterRoleTransactor: MinterRoleTransactor{contract: contract}, MinterRoleFilterer: MinterRoleFilterer{contract: contract}}, nil
}

// NewMinterRoleCaller creates a new read-only instance of MinterRole, bound to a specific deployed contract.
func NewMinterRoleCaller(address common.Address, caller bind.ContractCaller) (*MinterRoleCaller, error) {
	contract, err := bindMinterRole(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MinterRoleCaller{contract: contract}, nil
}

// NewMinterRoleTransactor creates a new write-only instance of MinterRole, bound to a specific deployed contract.
func NewMinterRoleTransactor(address common.Address, transactor bind.ContractTransactor) (*MinterRoleTransactor, error) {
	contract, err := bindMinterRole(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MinterRoleTransactor{contract: contract}, nil
}

// NewMinterRoleFilterer creates a new log filterer instance of MinterRole, bound to a specific deployed contract.
func NewMinterRoleFilterer(address common.Address, filterer bind.ContractFilterer) (*MinterRoleFilterer, error) {
	contract, err := bindMinterRole(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MinterRoleFilterer{contract: contract}, nil
}

// bindMinterRole binds a generic wrapper to an already deployed contract.
func bindMinterRole(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(MinterRoleABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MinterRole *MinterRoleRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _MinterRole.Contract.MinterRoleCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MinterRole *MinterRoleRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MinterRole.Contract.MinterRoleTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MinterRole *MinterRoleRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MinterRole.Contract.MinterRoleTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MinterRole *MinterRoleCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _MinterRole.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MinterRole *MinterRoleTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MinterRole.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MinterRole *MinterRoleTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MinterRole.Contract.contract.Transact(opts, method, params...)
}

// IsMinter is a free data retrieval call binding the contract method 0xaa271e1a.
//
// Solidity: function isMinter(address account) view returns(bool)
func (_MinterRole *MinterRoleCaller) IsMinter(opts *bind.CallOpts, account common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _MinterRole.contract.Call(opts, out, "isMinter", account)
	return *ret0, err
}

// IsMinter is a free data retrieval call binding the contract method 0xaa271e1a.
//
// Solidity: function isMinter(address account) view returns(bool)
func (_MinterRole *MinterRoleSession) IsMinter(account common.Address) (bool, error) {
	return _MinterRole.Contract.IsMinter(&_MinterRole.CallOpts, account)
}

// IsMinter is a free data retrieval call binding the contract method 0xaa271e1a.
//
// Solidity: function isMinter(address account) view returns(bool)
func (_MinterRole *MinterRoleCallerSession) IsMinter(account common.Address) (bool, error) {
	return _MinterRole.Contract.IsMinter(&_MinterRole.CallOpts, account)
}

// AddMinter is a paid mutator transaction binding the contract method 0x983b2d56.
//
// Solidity: function addMinter(address account) returns()
func (_MinterRole *MinterRoleTransactor) AddMinter(opts *bind.TransactOpts, account common.Address) (*types.Transaction, error) {
	return _MinterRole.contract.Transact(opts, "addMinter", account)
}

// AddMinter is a paid mutator transaction binding the contract method 0x983b2d56.
//
// Solidity: function addMinter(address account) returns()
func (_MinterRole *MinterRoleSession) AddMinter(account common.Address) (*types.Transaction, error) {
	return _MinterRole.Contract.AddMinter(&_MinterRole.TransactOpts, account)
}

// AddMinter is a paid mutator transaction binding the contract method 0x983b2d56.
//
// Solidity: function addMinter(address account) returns()
func (_MinterRole *MinterRoleTransactorSession) AddMinter(account common.Address) (*types.Transaction, error) {
	return _MinterRole.Contract.AddMinter(&_MinterRole.TransactOpts, account)
}

// RenounceMinter is a paid mutator transaction binding the contract method 0x98650275.
//
// Solidity: function renounceMinter() returns()
func (_MinterRole *MinterRoleTransactor) RenounceMinter(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MinterRole.contract.Transact(opts, "renounceMinter")
}

// RenounceMinter is a paid mutator transaction binding the contract method 0x98650275.
//
// Solidity: function renounceMinter() returns()
func (_MinterRole *MinterRoleSession) RenounceMinter() (*types.Transaction, error) {
	return _MinterRole.Contract.RenounceMinter(&_MinterRole.TransactOpts)
}

// RenounceMinter is a paid mutator transaction binding the contract method 0x98650275.
//
// Solidity: function renounceMinter() returns()
func (_MinterRole *MinterRoleTransactorSession) RenounceMinter() (*types.Transaction, error) {
	return _MinterRole.Contract.RenounceMinter(&_MinterRole.TransactOpts)
}

// MinterRoleMinterAddedIterator is returned from FilterMinterAdded and is used to iterate over the raw logs and unpacked data for MinterAdded events raised by the MinterRole contract.
type MinterRoleMinterAddedIterator struct {
	Event *MinterRoleMinterAdded // Event containing the contract specifics and raw log

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
func (it *MinterRoleMinterAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MinterRoleMinterAdded)
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
		it.Event = new(MinterRoleMinterAdded)
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
func (it *MinterRoleMinterAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MinterRoleMinterAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MinterRoleMinterAdded represents a MinterAdded event raised by the MinterRole contract.
type MinterRoleMinterAdded struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterMinterAdded is a free log retrieval operation binding the contract event 0x6ae172837ea30b801fbfcdd4108aa1d5bf8ff775444fd70256b44e6bf3dfc3f6.
//
// Solidity: event MinterAdded(address indexed account)
func (_MinterRole *MinterRoleFilterer) FilterMinterAdded(opts *bind.FilterOpts, account []common.Address) (*MinterRoleMinterAddedIterator, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _MinterRole.contract.FilterLogs(opts, "MinterAdded", accountRule)
	if err != nil {
		return nil, err
	}
	return &MinterRoleMinterAddedIterator{contract: _MinterRole.contract, event: "MinterAdded", logs: logs, sub: sub}, nil
}

// WatchMinterAdded is a free log subscription operation binding the contract event 0x6ae172837ea30b801fbfcdd4108aa1d5bf8ff775444fd70256b44e6bf3dfc3f6.
//
// Solidity: event MinterAdded(address indexed account)
func (_MinterRole *MinterRoleFilterer) WatchMinterAdded(opts *bind.WatchOpts, sink chan<- *MinterRoleMinterAdded, account []common.Address) (event.Subscription, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _MinterRole.contract.WatchLogs(opts, "MinterAdded", accountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MinterRoleMinterAdded)
				if err := _MinterRole.contract.UnpackLog(event, "MinterAdded", log); err != nil {
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

// ParseMinterAdded is a log parse operation binding the contract event 0x6ae172837ea30b801fbfcdd4108aa1d5bf8ff775444fd70256b44e6bf3dfc3f6.
//
// Solidity: event MinterAdded(address indexed account)
func (_MinterRole *MinterRoleFilterer) ParseMinterAdded(log types.Log) (*MinterRoleMinterAdded, error) {
	event := new(MinterRoleMinterAdded)
	if err := _MinterRole.contract.UnpackLog(event, "MinterAdded", log); err != nil {
		return nil, err
	}
	return event, nil
}

// MinterRoleMinterRemovedIterator is returned from FilterMinterRemoved and is used to iterate over the raw logs and unpacked data for MinterRemoved events raised by the MinterRole contract.
type MinterRoleMinterRemovedIterator struct {
	Event *MinterRoleMinterRemoved // Event containing the contract specifics and raw log

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
func (it *MinterRoleMinterRemovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MinterRoleMinterRemoved)
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
		it.Event = new(MinterRoleMinterRemoved)
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
func (it *MinterRoleMinterRemovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MinterRoleMinterRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MinterRoleMinterRemoved represents a MinterRemoved event raised by the MinterRole contract.
type MinterRoleMinterRemoved struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterMinterRemoved is a free log retrieval operation binding the contract event 0xe94479a9f7e1952cc78f2d6baab678adc1b772d936c6583def489e524cb66692.
//
// Solidity: event MinterRemoved(address indexed account)
func (_MinterRole *MinterRoleFilterer) FilterMinterRemoved(opts *bind.FilterOpts, account []common.Address) (*MinterRoleMinterRemovedIterator, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _MinterRole.contract.FilterLogs(opts, "MinterRemoved", accountRule)
	if err != nil {
		return nil, err
	}
	return &MinterRoleMinterRemovedIterator{contract: _MinterRole.contract, event: "MinterRemoved", logs: logs, sub: sub}, nil
}

// WatchMinterRemoved is a free log subscription operation binding the contract event 0xe94479a9f7e1952cc78f2d6baab678adc1b772d936c6583def489e524cb66692.
//
// Solidity: event MinterRemoved(address indexed account)
func (_MinterRole *MinterRoleFilterer) WatchMinterRemoved(opts *bind.WatchOpts, sink chan<- *MinterRoleMinterRemoved, account []common.Address) (event.Subscription, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _MinterRole.contract.WatchLogs(opts, "MinterRemoved", accountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MinterRoleMinterRemoved)
				if err := _MinterRole.contract.UnpackLog(event, "MinterRemoved", log); err != nil {
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

// ParseMinterRemoved is a log parse operation binding the contract event 0xe94479a9f7e1952cc78f2d6baab678adc1b772d936c6583def489e524cb66692.
//
// Solidity: event MinterRemoved(address indexed account)
func (_MinterRole *MinterRoleFilterer) ParseMinterRemoved(log types.Log) (*MinterRoleMinterRemoved, error) {
	event := new(MinterRoleMinterRemoved)
	if err := _MinterRole.contract.UnpackLog(event, "MinterRemoved", log); err != nil {
		return nil, err
	}
	return event, nil
}

// OwnableABI is the input ABI used to generate the binding from.
const OwnableABI = "[{\"constant\":false,\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"isOwner\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"}]"

// OwnableBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const OwnableBinRuntime = ``

// OwnableFuncSigs maps the 4-byte function signature to its string representation.
var OwnableFuncSigs = map[string]string{
	"8f32d59b": "isOwner()",
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

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() view returns(bool)
func (_Ownable *OwnableCaller) IsOwner(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Ownable.contract.Call(opts, out, "isOwner")
	return *ret0, err
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() view returns(bool)
func (_Ownable *OwnableSession) IsOwner() (bool, error) {
	return _Ownable.Contract.IsOwner(&_Ownable.CallOpts)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() view returns(bool)
func (_Ownable *OwnableCallerSession) IsOwner() (bool, error) {
	return _Ownable.Contract.IsOwner(&_Ownable.CallOpts)
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

// RolesABI is the input ABI used to generate the binding from.
const RolesABI = "[]"

// RolesBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const RolesBinRuntime = `73000000000000000000000000000000000000000030146080604052600080fdfea165627a7a7230582026429ed555f83236ac076363d644804ff5b0baa4e098f1d4ef783a49049f7df70029`

// RolesBin is the compiled bytecode used for deploying new contracts.
var RolesBin = "0x604c6023600b82828239805160001a607314601657fe5b30600052607381538281f3fe73000000000000000000000000000000000000000030146080604052600080fdfea165627a7a7230582026429ed555f83236ac076363d644804ff5b0baa4e098f1d4ef783a49049f7df70029"

// DeployRoles deploys a new Klaytn contract, binding an instance of Roles to it.
func DeployRoles(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Roles, error) {
	parsed, err := abi.JSON(strings.NewReader(RolesABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(RolesBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Roles{RolesCaller: RolesCaller{contract: contract}, RolesTransactor: RolesTransactor{contract: contract}, RolesFilterer: RolesFilterer{contract: contract}}, nil
}

// Roles is an auto generated Go binding around a Klaytn contract.
type Roles struct {
	RolesCaller     // Read-only binding to the contract
	RolesTransactor // Write-only binding to the contract
	RolesFilterer   // Log filterer for contract events
}

// RolesCaller is an auto generated read-only Go binding around a Klaytn contract.
type RolesCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RolesTransactor is an auto generated write-only Go binding around a Klaytn contract.
type RolesTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RolesFilterer is an auto generated log filtering Go binding around a Klaytn contract events.
type RolesFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RolesSession is an auto generated Go binding around a Klaytn contract,
// with pre-set call and transact options.
type RolesSession struct {
	Contract     *Roles            // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// RolesCallerSession is an auto generated read-only Go binding around a Klaytn contract,
// with pre-set call options.
type RolesCallerSession struct {
	Contract *RolesCaller  // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// RolesTransactorSession is an auto generated write-only Go binding around a Klaytn contract,
// with pre-set transact options.
type RolesTransactorSession struct {
	Contract     *RolesTransactor  // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// RolesRaw is an auto generated low-level Go binding around a Klaytn contract.
type RolesRaw struct {
	Contract *Roles // Generic contract binding to access the raw methods on
}

// RolesCallerRaw is an auto generated low-level read-only Go binding around a Klaytn contract.
type RolesCallerRaw struct {
	Contract *RolesCaller // Generic read-only contract binding to access the raw methods on
}

// RolesTransactorRaw is an auto generated low-level write-only Go binding around a Klaytn contract.
type RolesTransactorRaw struct {
	Contract *RolesTransactor // Generic write-only contract binding to access the raw methods on
}

// NewRoles creates a new instance of Roles, bound to a specific deployed contract.
func NewRoles(address common.Address, backend bind.ContractBackend) (*Roles, error) {
	contract, err := bindRoles(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Roles{RolesCaller: RolesCaller{contract: contract}, RolesTransactor: RolesTransactor{contract: contract}, RolesFilterer: RolesFilterer{contract: contract}}, nil
}

// NewRolesCaller creates a new read-only instance of Roles, bound to a specific deployed contract.
func NewRolesCaller(address common.Address, caller bind.ContractCaller) (*RolesCaller, error) {
	contract, err := bindRoles(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &RolesCaller{contract: contract}, nil
}

// NewRolesTransactor creates a new write-only instance of Roles, bound to a specific deployed contract.
func NewRolesTransactor(address common.Address, transactor bind.ContractTransactor) (*RolesTransactor, error) {
	contract, err := bindRoles(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &RolesTransactor{contract: contract}, nil
}

// NewRolesFilterer creates a new log filterer instance of Roles, bound to a specific deployed contract.
func NewRolesFilterer(address common.Address, filterer bind.ContractFilterer) (*RolesFilterer, error) {
	contract, err := bindRoles(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &RolesFilterer{contract: contract}, nil
}

// bindRoles binds a generic wrapper to an already deployed contract.
func bindRoles(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(RolesABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Roles *RolesRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Roles.Contract.RolesCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Roles *RolesRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Roles.Contract.RolesTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Roles *RolesRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Roles.Contract.RolesTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Roles *RolesCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Roles.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Roles *RolesTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Roles.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Roles *RolesTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Roles.Contract.contract.Transact(opts, method, params...)
}

// SafeMathABI is the input ABI used to generate the binding from.
const SafeMathABI = "[]"

// SafeMathBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const SafeMathBinRuntime = `73000000000000000000000000000000000000000030146080604052600080fdfea165627a7a723058201aa360e6241ae1bc829140fb38ae5ca1672181b2a5c529c4ca3a045ed2896feb0029`

// SafeMathBin is the compiled bytecode used for deploying new contracts.
var SafeMathBin = "0x604c6023600b82828239805160001a607314601657fe5b30600052607381538281f3fe73000000000000000000000000000000000000000030146080604052600080fdfea165627a7a723058201aa360e6241ae1bc829140fb38ae5ca1672181b2a5c529c4ca3a045ed2896feb0029"

// DeploySafeMath deploys a new Klaytn contract, binding an instance of SafeMath to it.
func DeploySafeMath(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *SafeMath, error) {
	parsed, err := abi.JSON(strings.NewReader(SafeMathABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(SafeMathBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &SafeMath{SafeMathCaller: SafeMathCaller{contract: contract}, SafeMathTransactor: SafeMathTransactor{contract: contract}, SafeMathFilterer: SafeMathFilterer{contract: contract}}, nil
}

// SafeMath is an auto generated Go binding around a Klaytn contract.
type SafeMath struct {
	SafeMathCaller     // Read-only binding to the contract
	SafeMathTransactor // Write-only binding to the contract
	SafeMathFilterer   // Log filterer for contract events
}

// SafeMathCaller is an auto generated read-only Go binding around a Klaytn contract.
type SafeMathCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SafeMathTransactor is an auto generated write-only Go binding around a Klaytn contract.
type SafeMathTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SafeMathFilterer is an auto generated log filtering Go binding around a Klaytn contract events.
type SafeMathFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SafeMathSession is an auto generated Go binding around a Klaytn contract,
// with pre-set call and transact options.
type SafeMathSession struct {
	Contract     *SafeMath         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SafeMathCallerSession is an auto generated read-only Go binding around a Klaytn contract,
// with pre-set call options.
type SafeMathCallerSession struct {
	Contract *SafeMathCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// SafeMathTransactorSession is an auto generated write-only Go binding around a Klaytn contract,
// with pre-set transact options.
type SafeMathTransactorSession struct {
	Contract     *SafeMathTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// SafeMathRaw is an auto generated low-level Go binding around a Klaytn contract.
type SafeMathRaw struct {
	Contract *SafeMath // Generic contract binding to access the raw methods on
}

// SafeMathCallerRaw is an auto generated low-level read-only Go binding around a Klaytn contract.
type SafeMathCallerRaw struct {
	Contract *SafeMathCaller // Generic read-only contract binding to access the raw methods on
}

// SafeMathTransactorRaw is an auto generated low-level write-only Go binding around a Klaytn contract.
type SafeMathTransactorRaw struct {
	Contract *SafeMathTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSafeMath creates a new instance of SafeMath, bound to a specific deployed contract.
func NewSafeMath(address common.Address, backend bind.ContractBackend) (*SafeMath, error) {
	contract, err := bindSafeMath(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SafeMath{SafeMathCaller: SafeMathCaller{contract: contract}, SafeMathTransactor: SafeMathTransactor{contract: contract}, SafeMathFilterer: SafeMathFilterer{contract: contract}}, nil
}

// NewSafeMathCaller creates a new read-only instance of SafeMath, bound to a specific deployed contract.
func NewSafeMathCaller(address common.Address, caller bind.ContractCaller) (*SafeMathCaller, error) {
	contract, err := bindSafeMath(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SafeMathCaller{contract: contract}, nil
}

// NewSafeMathTransactor creates a new write-only instance of SafeMath, bound to a specific deployed contract.
func NewSafeMathTransactor(address common.Address, transactor bind.ContractTransactor) (*SafeMathTransactor, error) {
	contract, err := bindSafeMath(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SafeMathTransactor{contract: contract}, nil
}

// NewSafeMathFilterer creates a new log filterer instance of SafeMath, bound to a specific deployed contract.
func NewSafeMathFilterer(address common.Address, filterer bind.ContractFilterer) (*SafeMathFilterer, error) {
	contract, err := bindSafeMath(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SafeMathFilterer{contract: contract}, nil
}

// bindSafeMath binds a generic wrapper to an already deployed contract.
func bindSafeMath(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(SafeMathABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SafeMath *SafeMathRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _SafeMath.Contract.SafeMathCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SafeMath *SafeMathRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SafeMath.Contract.SafeMathTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SafeMath *SafeMathRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SafeMath.Contract.SafeMathTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SafeMath *SafeMathCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _SafeMath.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SafeMath *SafeMathTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SafeMath.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SafeMath *SafeMathTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SafeMath.Contract.contract.Transact(opts, method, params...)
}

// ServiceChainTokenABI is the input ABI used to generate the binding from.
const ServiceChainTokenABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"spender\",\"type\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"sender\",\"type\":\"address\"},{\"name\":\"recipient\",\"type\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"DECIMALS\",\"outputs\":[{\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"INITIAL_SUPPLY\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"spender\",\"type\":\"address\"},{\"name\":\"addedValue\",\"type\":\"uint256\"}],\"name\":\"increaseAllowance\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_amount\",\"type\":\"uint256\"},{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_feeLimit\",\"type\":\"uint256\"},{\"name\":\"_extraData\",\"type\":\"bytes\"}],\"name\":\"requestValueTransfer\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"account\",\"type\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"mint\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"account\",\"type\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"burnFrom\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_bridge\",\"type\":\"address\"}],\"name\":\"setBridge\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"isOwner\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"account\",\"type\":\"address\"}],\"name\":\"addMinter\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"renounceMinter\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"NAME\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"spender\",\"type\":\"address\"},{\"name\":\"subtractedValue\",\"type\":\"uint256\"}],\"name\":\"decreaseAllowance\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"recipient\",\"type\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"account\",\"type\":\"address\"}],\"name\":\"isMinter\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"owner\",\"type\":\"address\"},{\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"bridge\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"SYMBOL\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"_bridge\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"account\",\"type\":\"address\"}],\"name\":\"MinterAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"account\",\"type\":\"address\"}],\"name\":\"MinterRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"}]"

// ServiceChainTokenBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const ServiceChainTokenBinRuntime = `608060405234801561001057600080fd5b50600436106101735760003560e01c80638da5cb5b116100de578063a457c2d711610097578063dd62ed3e11610071578063dd62ed3e14610506578063e78cea9214610534578063f2fde38b1461053c578063f76f8d781461056257610173565b8063a457c2d714610488578063a9059cbb146104b4578063aa271e1a146104e057610173565b80638da5cb5b1461038b5780638dd14802146103af5780638f32d59b146103d5578063983b2d56146103dd5780639865027514610403578063a3f4df7e1461040b57610173565b80633dc3c9e1116101305780633dc3c9e11461025a57806340c10f19146102e857806342966c681461031457806370a0823114610331578063715018a61461035757806379cc67901461035f57610173565b8063095ea7b31461017857806318160ddd146101b857806323b872dd146101d25780632e0f2625146102085780632ff2e9dc14610226578063395093511461022e575b600080fd5b6101a46004803603604081101561018e57600080fd5b506001600160a01b03813516906020013561056a565b604080519115158252519081900360200190f35b6101c0610580565b60408051918252519081900360200190f35b6101a4600480360360608110156101e857600080fd5b506001600160a01b03813581169160208101359091169060400135610586565b6102106105dd565b6040805160ff9092168252519081900360200190f35b6101c06105e2565b6101a46004803603604081101561024457600080fd5b506001600160a01b0381351690602001356105f2565b6102e66004803603608081101561027057600080fd5b8135916001600160a01b0360208201351691604082013591908101906080810160608201356401000000008111156102a757600080fd5b8201836020820111156102b957600080fd5b803590602001918460018302840111640100000000831117156102db57600080fd5b50909250905061062e565b005b6101a4600480360360408110156102fe57600080fd5b506001600160a01b038135169060200135610751565b6102e66004803603602081101561032a57600080fd5b50356107a4565b6101c06004803603602081101561034757600080fd5b50356001600160a01b03166107b1565b6102e66107cc565b6102e66004803603604081101561037557600080fd5b506001600160a01b038135169060200135610872565b610393610880565b604080516001600160a01b039092168252519081900360200190f35b6102e6600480360360208110156103c557600080fd5b50356001600160a01b031661088f565b6101a4610973565b6102e6600480360360208110156103f357600080fd5b50356001600160a01b0316610984565b6102e66109d4565b6104136109df565b6040805160208082528351818301528351919283929083019185019080838360005b8381101561044d578181015183820152602001610435565b50505050905090810190601f16801561047a5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6101a46004803603604081101561049e57600080fd5b506001600160a01b038135169060200135610a18565b6101a4600480360360408110156104ca57600080fd5b506001600160a01b038135169060200135610a54565b6101a4600480360360208110156104f657600080fd5b50356001600160a01b0316610a61565b6101c06004803603604081101561051c57600080fd5b506001600160a01b0381358116916020013516610a7a565b610393610aa5565b6102e66004803603602081101561055257600080fd5b50356001600160a01b0316610ab4565b610413610b19565b6000610577338484610b3b565b50600192915050565b60025490565b6000610593848484610c2d565b6001600160a01b0384166000908152600160209081526040808320338085529252909120546105d39186916105ce908663ffffffff610d7516565b610b3b565b5060019392505050565b601281565b6b033b2e3c9fd0803ce800000081565b3360008181526001602090815260408083206001600160a01b038716845290915281205490916105779185906105ce908663ffffffff610dd516565b600554610653906001600160a01b031661064e878663ffffffff610dd516565b610a54565b61069157604051600160e51b62461bcd0281526004018080602001828103825260258152602001806114286025913960400191505060405180910390fd5b600554604051600160e01b63f1656e5302815233600482018181526001600160a01b038881166024850152604484018a90526064840188905260a06084850190815260a4850187905294169363f1656e539389928b928a928a928a92909160c401848480828437600081840152601f19601f820116905080830192505050975050505050505050600060405180830381600087803b15801561073257600080fd5b505af1158015610746573d6000803e3d6000fd5b505050505050505050565b600061075c33610a61565b61079a57604051600160e51b62461bcd02815260040180806020018281038252603081526020018061134b6030913960400191505060405180910390fd5b6105778383610e39565b6107ae3382610f2c565b50565b6001600160a01b031660009081526020819052604090205490565b6107d4610973565b6108285760408051600160e51b62461bcd02815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572604482015290519081900360640190fd5b6004546040516000916001600160a01b0316907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0908390a3600480546001600160a01b0319169055565b61087c8282611008565b5050565b6004546001600160a01b031690565b610897610973565b6108eb5760408051600160e51b62461bcd02815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572604482015290519081900360640190fd5b6108fd816001600160a01b031661104d565b6109515760408051600160e51b62461bcd02815260206004820152601860248201527f627269646765206973206e6f74206120636f6e74726163740000000000000000604482015290519081900360640190fd5b600580546001600160a01b0319166001600160a01b0392909216919091179055565b6004546001600160a01b0316331490565b61098d33610a61565b6109cb57604051600160e51b62461bcd02815260040180806020018281038252603081526020018061134b6030913960400191505060405180910390fd5b6107ae81611053565b6109dd3361109b565b565b6040518060400160405280601181526020017f53657276696365436861696e546f6b656e00000000000000000000000000000081525081565b3360008181526001602090815260408083206001600160a01b038716845290915281205490916105779185906105ce908663ffffffff610d7516565b6000610577338484610c2d565b6000610a7460038363ffffffff6110e316565b92915050565b6001600160a01b03918216600090815260016020908152604080832093909416825291909152205490565b6005546001600160a01b031681565b610abc610973565b610b105760408051600160e51b62461bcd02815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572604482015290519081900360640190fd5b6107ae8161114d565b604051806040016040528060038152602001600160ea1b6214d0d50281525081565b6001600160a01b038316610b8357604051600160e51b62461bcd0281526004018080602001828103825260248152602001806114046024913960400191505060405180910390fd5b6001600160a01b038216610bcb57604051600160e51b62461bcd0281526004018080602001828103825260228152602001806113296022913960400191505060405180910390fd5b6001600160a01b03808416600081815260016020908152604080832094871680845294825291829020859055815185815291517f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b9259281900390910190a3505050565b6001600160a01b038316610c7557604051600160e51b62461bcd0281526004018080602001828103825260258152602001806113df6025913960400191505060405180910390fd5b6001600160a01b038216610cbd57604051600160e51b62461bcd0281526004018080602001828103825260238152602001806112e06023913960400191505060405180910390fd5b6001600160a01b038316600090815260208190526040902054610ce6908263ffffffff610d7516565b6001600160a01b038085166000908152602081905260408082209390935590841681522054610d1b908263ffffffff610dd516565b6001600160a01b038084166000818152602081815260409182902094909455805185815290519193928716927fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef92918290030190a3505050565b600082821115610dcf5760408051600160e51b62461bcd02815260206004820152601e60248201527f536166654d6174683a207375627472616374696f6e206f766572666c6f770000604482015290519081900360640190fd5b50900390565b600082820183811015610e325760408051600160e51b62461bcd02815260206004820152601b60248201527f536166654d6174683a206164646974696f6e206f766572666c6f770000000000604482015290519081900360640190fd5b9392505050565b6001600160a01b038216610e975760408051600160e51b62461bcd02815260206004820152601f60248201527f45524332303a206d696e7420746f20746865207a65726f206164647265737300604482015290519081900360640190fd5b600254610eaa908263ffffffff610dd516565b6002556001600160a01b038216600090815260208190526040902054610ed6908263ffffffff610dd516565b6001600160a01b0383166000818152602081815260408083209490945583518581529351929391927fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef9281900390910190a35050565b6001600160a01b038216610f7457604051600160e51b62461bcd0281526004018080602001828103825260218152602001806113be6021913960400191505060405180910390fd5b600254610f87908263ffffffff610d7516565b6002556001600160a01b038216600090815260208190526040902054610fb3908263ffffffff610d7516565b6001600160a01b038316600081815260208181526040808320949094558351858152935191937fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef929081900390910190a35050565b6110128282610f2c565b6001600160a01b03821660009081526001602090815260408083203380855292529091205461087c9184916105ce908563ffffffff610d7516565b3b151590565b61106460038263ffffffff6111f116565b6040516001600160a01b038216907f6ae172837ea30b801fbfcdd4108aa1d5bf8ff775444fd70256b44e6bf3dfc3f690600090a250565b6110ac60038263ffffffff61127516565b6040516001600160a01b038216907fe94479a9f7e1952cc78f2d6baab678adc1b772d936c6583def489e524cb6669290600090a250565b60006001600160a01b03821661112d57604051600160e51b62461bcd02815260040180806020018281038252602281526020018061139c6022913960400191505060405180910390fd5b506001600160a01b03166000908152602091909152604090205460ff1690565b6001600160a01b03811661119557604051600160e51b62461bcd0281526004018080602001828103825260268152602001806113036026913960400191505060405180910390fd5b6004546040516001600160a01b038084169216907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a3600480546001600160a01b0319166001600160a01b0392909216919091179055565b6111fb82826110e3565b156112505760408051600160e51b62461bcd02815260206004820152601f60248201527f526f6c65733a206163636f756e7420616c72656164792068617320726f6c6500604482015290519081900360640190fd5b6001600160a01b0316600090815260209190915260409020805460ff19166001179055565b61127f82826110e3565b6112bd57604051600160e51b62461bcd02815260040180806020018281038252602181526020018061137b6021913960400191505060405180910390fd5b6001600160a01b0316600090815260209190915260409020805460ff1916905556fe45524332303a207472616e7366657220746f20746865207a65726f20616464726573734f776e61626c653a206e6577206f776e657220697320746865207a65726f206164647265737345524332303a20617070726f766520746f20746865207a65726f20616464726573734d696e746572526f6c653a2063616c6c657220646f6573206e6f74206861766520746865204d696e74657220726f6c65526f6c65733a206163636f756e7420646f6573206e6f74206861766520726f6c65526f6c65733a206163636f756e7420697320746865207a65726f206164647265737345524332303a206275726e2066726f6d20746865207a65726f206164647265737345524332303a207472616e736665722066726f6d20746865207a65726f206164647265737345524332303a20617070726f76652066726f6d20746865207a65726f20616464726573737265717565737456616c75655472616e736665723a207472616e73666572206661696c6564a165627a7a72305820b3517cc790738552a1d3bd452fa586f1db04bfff07e57508df3fbd876f93717e0029`

// ServiceChainTokenFuncSigs maps the 4-byte function signature to its string representation.
var ServiceChainTokenFuncSigs = map[string]string{
	"2e0f2625": "DECIMALS()",
	"2ff2e9dc": "INITIAL_SUPPLY()",
	"a3f4df7e": "NAME()",
	"f76f8d78": "SYMBOL()",
	"983b2d56": "addMinter(address)",
	"dd62ed3e": "allowance(address,address)",
	"095ea7b3": "approve(address,uint256)",
	"70a08231": "balanceOf(address)",
	"e78cea92": "bridge()",
	"42966c68": "burn(uint256)",
	"79cc6790": "burnFrom(address,uint256)",
	"a457c2d7": "decreaseAllowance(address,uint256)",
	"39509351": "increaseAllowance(address,uint256)",
	"aa271e1a": "isMinter(address)",
	"8f32d59b": "isOwner()",
	"40c10f19": "mint(address,uint256)",
	"8da5cb5b": "owner()",
	"98650275": "renounceMinter()",
	"715018a6": "renounceOwnership()",
	"3dc3c9e1": "requestValueTransfer(uint256,address,uint256,bytes)",
	"8dd14802": "setBridge(address)",
	"18160ddd": "totalSupply()",
	"a9059cbb": "transfer(address,uint256)",
	"23b872dd": "transferFrom(address,address,uint256)",
	"f2fde38b": "transferOwnership(address)",
}

// ServiceChainTokenBin is the compiled bytecode used for deploying new contracts.
var ServiceChainTokenBin = "0x60806040523480156200001157600080fd5b50604051602080620019c6833981018060405260208110156200003357600080fd5b5051806200004833620000cc602090811b901c565b600480546001600160a01b0319163317908190556040516001600160a01b0391909116906000907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0908290a3620000a5816200011e60201b60201c565b50620000c5336b033b2e3c9fd0803ce800000062000247602090811b901c565b506200051c565b620000e78160036200036260201b620011f11790919060201c565b6040516001600160a01b038216907f6ae172837ea30b801fbfcdd4108aa1d5bf8ff775444fd70256b44e6bf3dfc3f690600090a250565b6200012e6200040660201b60201c565b6200019a57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572604482015290519081900360640190fd5b620001b9816001600160a01b03166200041760201b6200104d1760201c565b6200022557604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601860248201527f627269646765206973206e6f74206120636f6e74726163740000000000000000604482015290519081900360640190fd5b600580546001600160a01b0319166001600160a01b0392909216919091179055565b6001600160a01b038216620002bd57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601f60248201527f45524332303a206d696e7420746f20746865207a65726f206164647265737300604482015290519081900360640190fd5b620002d9816002546200041d60201b62000dd51790919060201c565b6002556001600160a01b038216600090815260208181526040909120546200030c91839062000dd56200041d821b17901c565b6001600160a01b0383166000818152602081815260408083209490945583518581529351929391927fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef9281900390910190a35050565b6200037482826200049960201b60201c565b15620003e157604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601f60248201527f526f6c65733a206163636f756e7420616c72656164792068617320726f6c6500604482015290519081900360640190fd5b6001600160a01b0316600090815260209190915260409020805460ff19166001179055565b6004546001600160a01b0316331490565b3b151590565b6000828201838110156200049257604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601b60248201527f536166654d6174683a206164646974696f6e206f766572666c6f770000000000604482015290519081900360640190fd5b9392505050565b60006001600160a01b038216620004fc576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401808060200182810382526022815260200180620019a46022913960400191505060405180910390fd5b506001600160a01b03166000908152602091909152604090205460ff1690565b611478806200052c6000396000f3fe608060405234801561001057600080fd5b50600436106101735760003560e01c80638da5cb5b116100de578063a457c2d711610097578063dd62ed3e11610071578063dd62ed3e14610506578063e78cea9214610534578063f2fde38b1461053c578063f76f8d781461056257610173565b8063a457c2d714610488578063a9059cbb146104b4578063aa271e1a146104e057610173565b80638da5cb5b1461038b5780638dd14802146103af5780638f32d59b146103d5578063983b2d56146103dd5780639865027514610403578063a3f4df7e1461040b57610173565b80633dc3c9e1116101305780633dc3c9e11461025a57806340c10f19146102e857806342966c681461031457806370a0823114610331578063715018a61461035757806379cc67901461035f57610173565b8063095ea7b31461017857806318160ddd146101b857806323b872dd146101d25780632e0f2625146102085780632ff2e9dc14610226578063395093511461022e575b600080fd5b6101a46004803603604081101561018e57600080fd5b506001600160a01b03813516906020013561056a565b604080519115158252519081900360200190f35b6101c0610580565b60408051918252519081900360200190f35b6101a4600480360360608110156101e857600080fd5b506001600160a01b03813581169160208101359091169060400135610586565b6102106105dd565b6040805160ff9092168252519081900360200190f35b6101c06105e2565b6101a46004803603604081101561024457600080fd5b506001600160a01b0381351690602001356105f2565b6102e66004803603608081101561027057600080fd5b8135916001600160a01b0360208201351691604082013591908101906080810160608201356401000000008111156102a757600080fd5b8201836020820111156102b957600080fd5b803590602001918460018302840111640100000000831117156102db57600080fd5b50909250905061062e565b005b6101a4600480360360408110156102fe57600080fd5b506001600160a01b038135169060200135610751565b6102e66004803603602081101561032a57600080fd5b50356107a4565b6101c06004803603602081101561034757600080fd5b50356001600160a01b03166107b1565b6102e66107cc565b6102e66004803603604081101561037557600080fd5b506001600160a01b038135169060200135610872565b610393610880565b604080516001600160a01b039092168252519081900360200190f35b6102e6600480360360208110156103c557600080fd5b50356001600160a01b031661088f565b6101a4610973565b6102e6600480360360208110156103f357600080fd5b50356001600160a01b0316610984565b6102e66109d4565b6104136109df565b6040805160208082528351818301528351919283929083019185019080838360005b8381101561044d578181015183820152602001610435565b50505050905090810190601f16801561047a5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6101a46004803603604081101561049e57600080fd5b506001600160a01b038135169060200135610a18565b6101a4600480360360408110156104ca57600080fd5b506001600160a01b038135169060200135610a54565b6101a4600480360360208110156104f657600080fd5b50356001600160a01b0316610a61565b6101c06004803603604081101561051c57600080fd5b506001600160a01b0381358116916020013516610a7a565b610393610aa5565b6102e66004803603602081101561055257600080fd5b50356001600160a01b0316610ab4565b610413610b19565b6000610577338484610b3b565b50600192915050565b60025490565b6000610593848484610c2d565b6001600160a01b0384166000908152600160209081526040808320338085529252909120546105d39186916105ce908663ffffffff610d7516565b610b3b565b5060019392505050565b601281565b6b033b2e3c9fd0803ce800000081565b3360008181526001602090815260408083206001600160a01b038716845290915281205490916105779185906105ce908663ffffffff610dd516565b600554610653906001600160a01b031661064e878663ffffffff610dd516565b610a54565b61069157604051600160e51b62461bcd0281526004018080602001828103825260258152602001806114286025913960400191505060405180910390fd5b600554604051600160e01b63f1656e5302815233600482018181526001600160a01b038881166024850152604484018a90526064840188905260a06084850190815260a4850187905294169363f1656e539389928b928a928a928a92909160c401848480828437600081840152601f19601f820116905080830192505050975050505050505050600060405180830381600087803b15801561073257600080fd5b505af1158015610746573d6000803e3d6000fd5b505050505050505050565b600061075c33610a61565b61079a57604051600160e51b62461bcd02815260040180806020018281038252603081526020018061134b6030913960400191505060405180910390fd5b6105778383610e39565b6107ae3382610f2c565b50565b6001600160a01b031660009081526020819052604090205490565b6107d4610973565b6108285760408051600160e51b62461bcd02815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572604482015290519081900360640190fd5b6004546040516000916001600160a01b0316907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0908390a3600480546001600160a01b0319169055565b61087c8282611008565b5050565b6004546001600160a01b031690565b610897610973565b6108eb5760408051600160e51b62461bcd02815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572604482015290519081900360640190fd5b6108fd816001600160a01b031661104d565b6109515760408051600160e51b62461bcd02815260206004820152601860248201527f627269646765206973206e6f74206120636f6e74726163740000000000000000604482015290519081900360640190fd5b600580546001600160a01b0319166001600160a01b0392909216919091179055565b6004546001600160a01b0316331490565b61098d33610a61565b6109cb57604051600160e51b62461bcd02815260040180806020018281038252603081526020018061134b6030913960400191505060405180910390fd5b6107ae81611053565b6109dd3361109b565b565b6040518060400160405280601181526020017f53657276696365436861696e546f6b656e00000000000000000000000000000081525081565b3360008181526001602090815260408083206001600160a01b038716845290915281205490916105779185906105ce908663ffffffff610d7516565b6000610577338484610c2d565b6000610a7460038363ffffffff6110e316565b92915050565b6001600160a01b03918216600090815260016020908152604080832093909416825291909152205490565b6005546001600160a01b031681565b610abc610973565b610b105760408051600160e51b62461bcd02815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572604482015290519081900360640190fd5b6107ae8161114d565b604051806040016040528060038152602001600160ea1b6214d0d50281525081565b6001600160a01b038316610b8357604051600160e51b62461bcd0281526004018080602001828103825260248152602001806114046024913960400191505060405180910390fd5b6001600160a01b038216610bcb57604051600160e51b62461bcd0281526004018080602001828103825260228152602001806113296022913960400191505060405180910390fd5b6001600160a01b03808416600081815260016020908152604080832094871680845294825291829020859055815185815291517f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b9259281900390910190a3505050565b6001600160a01b038316610c7557604051600160e51b62461bcd0281526004018080602001828103825260258152602001806113df6025913960400191505060405180910390fd5b6001600160a01b038216610cbd57604051600160e51b62461bcd0281526004018080602001828103825260238152602001806112e06023913960400191505060405180910390fd5b6001600160a01b038316600090815260208190526040902054610ce6908263ffffffff610d7516565b6001600160a01b038085166000908152602081905260408082209390935590841681522054610d1b908263ffffffff610dd516565b6001600160a01b038084166000818152602081815260409182902094909455805185815290519193928716927fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef92918290030190a3505050565b600082821115610dcf5760408051600160e51b62461bcd02815260206004820152601e60248201527f536166654d6174683a207375627472616374696f6e206f766572666c6f770000604482015290519081900360640190fd5b50900390565b600082820183811015610e325760408051600160e51b62461bcd02815260206004820152601b60248201527f536166654d6174683a206164646974696f6e206f766572666c6f770000000000604482015290519081900360640190fd5b9392505050565b6001600160a01b038216610e975760408051600160e51b62461bcd02815260206004820152601f60248201527f45524332303a206d696e7420746f20746865207a65726f206164647265737300604482015290519081900360640190fd5b600254610eaa908263ffffffff610dd516565b6002556001600160a01b038216600090815260208190526040902054610ed6908263ffffffff610dd516565b6001600160a01b0383166000818152602081815260408083209490945583518581529351929391927fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef9281900390910190a35050565b6001600160a01b038216610f7457604051600160e51b62461bcd0281526004018080602001828103825260218152602001806113be6021913960400191505060405180910390fd5b600254610f87908263ffffffff610d7516565b6002556001600160a01b038216600090815260208190526040902054610fb3908263ffffffff610d7516565b6001600160a01b038316600081815260208181526040808320949094558351858152935191937fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef929081900390910190a35050565b6110128282610f2c565b6001600160a01b03821660009081526001602090815260408083203380855292529091205461087c9184916105ce908563ffffffff610d7516565b3b151590565b61106460038263ffffffff6111f116565b6040516001600160a01b038216907f6ae172837ea30b801fbfcdd4108aa1d5bf8ff775444fd70256b44e6bf3dfc3f690600090a250565b6110ac60038263ffffffff61127516565b6040516001600160a01b038216907fe94479a9f7e1952cc78f2d6baab678adc1b772d936c6583def489e524cb6669290600090a250565b60006001600160a01b03821661112d57604051600160e51b62461bcd02815260040180806020018281038252602281526020018061139c6022913960400191505060405180910390fd5b506001600160a01b03166000908152602091909152604090205460ff1690565b6001600160a01b03811661119557604051600160e51b62461bcd0281526004018080602001828103825260268152602001806113036026913960400191505060405180910390fd5b6004546040516001600160a01b038084169216907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a3600480546001600160a01b0319166001600160a01b0392909216919091179055565b6111fb82826110e3565b156112505760408051600160e51b62461bcd02815260206004820152601f60248201527f526f6c65733a206163636f756e7420616c72656164792068617320726f6c6500604482015290519081900360640190fd5b6001600160a01b0316600090815260209190915260409020805460ff19166001179055565b61127f82826110e3565b6112bd57604051600160e51b62461bcd02815260040180806020018281038252602181526020018061137b6021913960400191505060405180910390fd5b6001600160a01b0316600090815260209190915260409020805460ff1916905556fe45524332303a207472616e7366657220746f20746865207a65726f20616464726573734f776e61626c653a206e6577206f776e657220697320746865207a65726f206164647265737345524332303a20617070726f766520746f20746865207a65726f20616464726573734d696e746572526f6c653a2063616c6c657220646f6573206e6f74206861766520746865204d696e74657220726f6c65526f6c65733a206163636f756e7420646f6573206e6f74206861766520726f6c65526f6c65733a206163636f756e7420697320746865207a65726f206164647265737345524332303a206275726e2066726f6d20746865207a65726f206164647265737345524332303a207472616e736665722066726f6d20746865207a65726f206164647265737345524332303a20617070726f76652066726f6d20746865207a65726f20616464726573737265717565737456616c75655472616e736665723a207472616e73666572206661696c6564a165627a7a72305820b3517cc790738552a1d3bd452fa586f1db04bfff07e57508df3fbd876f93717e0029526f6c65733a206163636f756e7420697320746865207a65726f2061646472657373"

// DeployServiceChainToken deploys a new Klaytn contract, binding an instance of ServiceChainToken to it.
func DeployServiceChainToken(auth *bind.TransactOpts, backend bind.ContractBackend, _bridge common.Address) (common.Address, *types.Transaction, *ServiceChainToken, error) {
	parsed, err := abi.JSON(strings.NewReader(ServiceChainTokenABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(ServiceChainTokenBin), backend, _bridge)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ServiceChainToken{ServiceChainTokenCaller: ServiceChainTokenCaller{contract: contract}, ServiceChainTokenTransactor: ServiceChainTokenTransactor{contract: contract}, ServiceChainTokenFilterer: ServiceChainTokenFilterer{contract: contract}}, nil
}

// ServiceChainToken is an auto generated Go binding around a Klaytn contract.
type ServiceChainToken struct {
	ServiceChainTokenCaller     // Read-only binding to the contract
	ServiceChainTokenTransactor // Write-only binding to the contract
	ServiceChainTokenFilterer   // Log filterer for contract events
}

// ServiceChainTokenCaller is an auto generated read-only Go binding around a Klaytn contract.
type ServiceChainTokenCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ServiceChainTokenTransactor is an auto generated write-only Go binding around a Klaytn contract.
type ServiceChainTokenTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ServiceChainTokenFilterer is an auto generated log filtering Go binding around a Klaytn contract events.
type ServiceChainTokenFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ServiceChainTokenSession is an auto generated Go binding around a Klaytn contract,
// with pre-set call and transact options.
type ServiceChainTokenSession struct {
	Contract     *ServiceChainToken // Generic contract binding to set the session for
	CallOpts     bind.CallOpts      // Call options to use throughout this session
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// ServiceChainTokenCallerSession is an auto generated read-only Go binding around a Klaytn contract,
// with pre-set call options.
type ServiceChainTokenCallerSession struct {
	Contract *ServiceChainTokenCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts            // Call options to use throughout this session
}

// ServiceChainTokenTransactorSession is an auto generated write-only Go binding around a Klaytn contract,
// with pre-set transact options.
type ServiceChainTokenTransactorSession struct {
	Contract     *ServiceChainTokenTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// ServiceChainTokenRaw is an auto generated low-level Go binding around a Klaytn contract.
type ServiceChainTokenRaw struct {
	Contract *ServiceChainToken // Generic contract binding to access the raw methods on
}

// ServiceChainTokenCallerRaw is an auto generated low-level read-only Go binding around a Klaytn contract.
type ServiceChainTokenCallerRaw struct {
	Contract *ServiceChainTokenCaller // Generic read-only contract binding to access the raw methods on
}

// ServiceChainTokenTransactorRaw is an auto generated low-level write-only Go binding around a Klaytn contract.
type ServiceChainTokenTransactorRaw struct {
	Contract *ServiceChainTokenTransactor // Generic write-only contract binding to access the raw methods on
}

// NewServiceChainToken creates a new instance of ServiceChainToken, bound to a specific deployed contract.
func NewServiceChainToken(address common.Address, backend bind.ContractBackend) (*ServiceChainToken, error) {
	contract, err := bindServiceChainToken(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ServiceChainToken{ServiceChainTokenCaller: ServiceChainTokenCaller{contract: contract}, ServiceChainTokenTransactor: ServiceChainTokenTransactor{contract: contract}, ServiceChainTokenFilterer: ServiceChainTokenFilterer{contract: contract}}, nil
}

// NewServiceChainTokenCaller creates a new read-only instance of ServiceChainToken, bound to a specific deployed contract.
func NewServiceChainTokenCaller(address common.Address, caller bind.ContractCaller) (*ServiceChainTokenCaller, error) {
	contract, err := bindServiceChainToken(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ServiceChainTokenCaller{contract: contract}, nil
}

// NewServiceChainTokenTransactor creates a new write-only instance of ServiceChainToken, bound to a specific deployed contract.
func NewServiceChainTokenTransactor(address common.Address, transactor bind.ContractTransactor) (*ServiceChainTokenTransactor, error) {
	contract, err := bindServiceChainToken(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ServiceChainTokenTransactor{contract: contract}, nil
}

// NewServiceChainTokenFilterer creates a new log filterer instance of ServiceChainToken, bound to a specific deployed contract.
func NewServiceChainTokenFilterer(address common.Address, filterer bind.ContractFilterer) (*ServiceChainTokenFilterer, error) {
	contract, err := bindServiceChainToken(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ServiceChainTokenFilterer{contract: contract}, nil
}

// bindServiceChainToken binds a generic wrapper to an already deployed contract.
func bindServiceChainToken(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ServiceChainTokenABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ServiceChainToken *ServiceChainTokenRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _ServiceChainToken.Contract.ServiceChainTokenCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ServiceChainToken *ServiceChainTokenRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ServiceChainToken.Contract.ServiceChainTokenTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ServiceChainToken *ServiceChainTokenRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ServiceChainToken.Contract.ServiceChainTokenTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ServiceChainToken *ServiceChainTokenCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _ServiceChainToken.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ServiceChainToken *ServiceChainTokenTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ServiceChainToken.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ServiceChainToken *ServiceChainTokenTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ServiceChainToken.Contract.contract.Transact(opts, method, params...)
}

// DECIMALS is a free data retrieval call binding the contract method 0x2e0f2625.
//
// Solidity: function DECIMALS() view returns(uint8)
func (_ServiceChainToken *ServiceChainTokenCaller) DECIMALS(opts *bind.CallOpts) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _ServiceChainToken.contract.Call(opts, out, "DECIMALS")
	return *ret0, err
}

// DECIMALS is a free data retrieval call binding the contract method 0x2e0f2625.
//
// Solidity: function DECIMALS() view returns(uint8)
func (_ServiceChainToken *ServiceChainTokenSession) DECIMALS() (uint8, error) {
	return _ServiceChainToken.Contract.DECIMALS(&_ServiceChainToken.CallOpts)
}

// DECIMALS is a free data retrieval call binding the contract method 0x2e0f2625.
//
// Solidity: function DECIMALS() view returns(uint8)
func (_ServiceChainToken *ServiceChainTokenCallerSession) DECIMALS() (uint8, error) {
	return _ServiceChainToken.Contract.DECIMALS(&_ServiceChainToken.CallOpts)
}

// INITIALSUPPLY is a free data retrieval call binding the contract method 0x2ff2e9dc.
//
// Solidity: function INITIAL_SUPPLY() view returns(uint256)
func (_ServiceChainToken *ServiceChainTokenCaller) INITIALSUPPLY(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ServiceChainToken.contract.Call(opts, out, "INITIAL_SUPPLY")
	return *ret0, err
}

// INITIALSUPPLY is a free data retrieval call binding the contract method 0x2ff2e9dc.
//
// Solidity: function INITIAL_SUPPLY() view returns(uint256)
func (_ServiceChainToken *ServiceChainTokenSession) INITIALSUPPLY() (*big.Int, error) {
	return _ServiceChainToken.Contract.INITIALSUPPLY(&_ServiceChainToken.CallOpts)
}

// INITIALSUPPLY is a free data retrieval call binding the contract method 0x2ff2e9dc.
//
// Solidity: function INITIAL_SUPPLY() view returns(uint256)
func (_ServiceChainToken *ServiceChainTokenCallerSession) INITIALSUPPLY() (*big.Int, error) {
	return _ServiceChainToken.Contract.INITIALSUPPLY(&_ServiceChainToken.CallOpts)
}

// NAME is a free data retrieval call binding the contract method 0xa3f4df7e.
//
// Solidity: function NAME() view returns(string)
func (_ServiceChainToken *ServiceChainTokenCaller) NAME(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _ServiceChainToken.contract.Call(opts, out, "NAME")
	return *ret0, err
}

// NAME is a free data retrieval call binding the contract method 0xa3f4df7e.
//
// Solidity: function NAME() view returns(string)
func (_ServiceChainToken *ServiceChainTokenSession) NAME() (string, error) {
	return _ServiceChainToken.Contract.NAME(&_ServiceChainToken.CallOpts)
}

// NAME is a free data retrieval call binding the contract method 0xa3f4df7e.
//
// Solidity: function NAME() view returns(string)
func (_ServiceChainToken *ServiceChainTokenCallerSession) NAME() (string, error) {
	return _ServiceChainToken.Contract.NAME(&_ServiceChainToken.CallOpts)
}

// SYMBOL is a free data retrieval call binding the contract method 0xf76f8d78.
//
// Solidity: function SYMBOL() view returns(string)
func (_ServiceChainToken *ServiceChainTokenCaller) SYMBOL(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _ServiceChainToken.contract.Call(opts, out, "SYMBOL")
	return *ret0, err
}

// SYMBOL is a free data retrieval call binding the contract method 0xf76f8d78.
//
// Solidity: function SYMBOL() view returns(string)
func (_ServiceChainToken *ServiceChainTokenSession) SYMBOL() (string, error) {
	return _ServiceChainToken.Contract.SYMBOL(&_ServiceChainToken.CallOpts)
}

// SYMBOL is a free data retrieval call binding the contract method 0xf76f8d78.
//
// Solidity: function SYMBOL() view returns(string)
func (_ServiceChainToken *ServiceChainTokenCallerSession) SYMBOL() (string, error) {
	return _ServiceChainToken.Contract.SYMBOL(&_ServiceChainToken.CallOpts)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_ServiceChainToken *ServiceChainTokenCaller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ServiceChainToken.contract.Call(opts, out, "allowance", owner, spender)
	return *ret0, err
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_ServiceChainToken *ServiceChainTokenSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _ServiceChainToken.Contract.Allowance(&_ServiceChainToken.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_ServiceChainToken *ServiceChainTokenCallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _ServiceChainToken.Contract.Allowance(&_ServiceChainToken.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_ServiceChainToken *ServiceChainTokenCaller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ServiceChainToken.contract.Call(opts, out, "balanceOf", account)
	return *ret0, err
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_ServiceChainToken *ServiceChainTokenSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _ServiceChainToken.Contract.BalanceOf(&_ServiceChainToken.CallOpts, account)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_ServiceChainToken *ServiceChainTokenCallerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _ServiceChainToken.Contract.BalanceOf(&_ServiceChainToken.CallOpts, account)
}

// Bridge is a free data retrieval call binding the contract method 0xe78cea92.
//
// Solidity: function bridge() view returns(address)
func (_ServiceChainToken *ServiceChainTokenCaller) Bridge(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _ServiceChainToken.contract.Call(opts, out, "bridge")
	return *ret0, err
}

// Bridge is a free data retrieval call binding the contract method 0xe78cea92.
//
// Solidity: function bridge() view returns(address)
func (_ServiceChainToken *ServiceChainTokenSession) Bridge() (common.Address, error) {
	return _ServiceChainToken.Contract.Bridge(&_ServiceChainToken.CallOpts)
}

// Bridge is a free data retrieval call binding the contract method 0xe78cea92.
//
// Solidity: function bridge() view returns(address)
func (_ServiceChainToken *ServiceChainTokenCallerSession) Bridge() (common.Address, error) {
	return _ServiceChainToken.Contract.Bridge(&_ServiceChainToken.CallOpts)
}

// IsMinter is a free data retrieval call binding the contract method 0xaa271e1a.
//
// Solidity: function isMinter(address account) view returns(bool)
func (_ServiceChainToken *ServiceChainTokenCaller) IsMinter(opts *bind.CallOpts, account common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _ServiceChainToken.contract.Call(opts, out, "isMinter", account)
	return *ret0, err
}

// IsMinter is a free data retrieval call binding the contract method 0xaa271e1a.
//
// Solidity: function isMinter(address account) view returns(bool)
func (_ServiceChainToken *ServiceChainTokenSession) IsMinter(account common.Address) (bool, error) {
	return _ServiceChainToken.Contract.IsMinter(&_ServiceChainToken.CallOpts, account)
}

// IsMinter is a free data retrieval call binding the contract method 0xaa271e1a.
//
// Solidity: function isMinter(address account) view returns(bool)
func (_ServiceChainToken *ServiceChainTokenCallerSession) IsMinter(account common.Address) (bool, error) {
	return _ServiceChainToken.Contract.IsMinter(&_ServiceChainToken.CallOpts, account)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() view returns(bool)
func (_ServiceChainToken *ServiceChainTokenCaller) IsOwner(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _ServiceChainToken.contract.Call(opts, out, "isOwner")
	return *ret0, err
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() view returns(bool)
func (_ServiceChainToken *ServiceChainTokenSession) IsOwner() (bool, error) {
	return _ServiceChainToken.Contract.IsOwner(&_ServiceChainToken.CallOpts)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() view returns(bool)
func (_ServiceChainToken *ServiceChainTokenCallerSession) IsOwner() (bool, error) {
	return _ServiceChainToken.Contract.IsOwner(&_ServiceChainToken.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ServiceChainToken *ServiceChainTokenCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _ServiceChainToken.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ServiceChainToken *ServiceChainTokenSession) Owner() (common.Address, error) {
	return _ServiceChainToken.Contract.Owner(&_ServiceChainToken.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ServiceChainToken *ServiceChainTokenCallerSession) Owner() (common.Address, error) {
	return _ServiceChainToken.Contract.Owner(&_ServiceChainToken.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_ServiceChainToken *ServiceChainTokenCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ServiceChainToken.contract.Call(opts, out, "totalSupply")
	return *ret0, err
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_ServiceChainToken *ServiceChainTokenSession) TotalSupply() (*big.Int, error) {
	return _ServiceChainToken.Contract.TotalSupply(&_ServiceChainToken.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_ServiceChainToken *ServiceChainTokenCallerSession) TotalSupply() (*big.Int, error) {
	return _ServiceChainToken.Contract.TotalSupply(&_ServiceChainToken.CallOpts)
}

// AddMinter is a paid mutator transaction binding the contract method 0x983b2d56.
//
// Solidity: function addMinter(address account) returns()
func (_ServiceChainToken *ServiceChainTokenTransactor) AddMinter(opts *bind.TransactOpts, account common.Address) (*types.Transaction, error) {
	return _ServiceChainToken.contract.Transact(opts, "addMinter", account)
}

// AddMinter is a paid mutator transaction binding the contract method 0x983b2d56.
//
// Solidity: function addMinter(address account) returns()
func (_ServiceChainToken *ServiceChainTokenSession) AddMinter(account common.Address) (*types.Transaction, error) {
	return _ServiceChainToken.Contract.AddMinter(&_ServiceChainToken.TransactOpts, account)
}

// AddMinter is a paid mutator transaction binding the contract method 0x983b2d56.
//
// Solidity: function addMinter(address account) returns()
func (_ServiceChainToken *ServiceChainTokenTransactorSession) AddMinter(account common.Address) (*types.Transaction, error) {
	return _ServiceChainToken.Contract.AddMinter(&_ServiceChainToken.TransactOpts, account)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_ServiceChainToken *ServiceChainTokenTransactor) Approve(opts *bind.TransactOpts, spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _ServiceChainToken.contract.Transact(opts, "approve", spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_ServiceChainToken *ServiceChainTokenSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _ServiceChainToken.Contract.Approve(&_ServiceChainToken.TransactOpts, spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_ServiceChainToken *ServiceChainTokenTransactorSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _ServiceChainToken.Contract.Approve(&_ServiceChainToken.TransactOpts, spender, value)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 amount) returns()
func (_ServiceChainToken *ServiceChainTokenTransactor) Burn(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _ServiceChainToken.contract.Transact(opts, "burn", amount)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 amount) returns()
func (_ServiceChainToken *ServiceChainTokenSession) Burn(amount *big.Int) (*types.Transaction, error) {
	return _ServiceChainToken.Contract.Burn(&_ServiceChainToken.TransactOpts, amount)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 amount) returns()
func (_ServiceChainToken *ServiceChainTokenTransactorSession) Burn(amount *big.Int) (*types.Transaction, error) {
	return _ServiceChainToken.Contract.Burn(&_ServiceChainToken.TransactOpts, amount)
}

// BurnFrom is a paid mutator transaction binding the contract method 0x79cc6790.
//
// Solidity: function burnFrom(address account, uint256 amount) returns()
func (_ServiceChainToken *ServiceChainTokenTransactor) BurnFrom(opts *bind.TransactOpts, account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ServiceChainToken.contract.Transact(opts, "burnFrom", account, amount)
}

// BurnFrom is a paid mutator transaction binding the contract method 0x79cc6790.
//
// Solidity: function burnFrom(address account, uint256 amount) returns()
func (_ServiceChainToken *ServiceChainTokenSession) BurnFrom(account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ServiceChainToken.Contract.BurnFrom(&_ServiceChainToken.TransactOpts, account, amount)
}

// BurnFrom is a paid mutator transaction binding the contract method 0x79cc6790.
//
// Solidity: function burnFrom(address account, uint256 amount) returns()
func (_ServiceChainToken *ServiceChainTokenTransactorSession) BurnFrom(account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ServiceChainToken.Contract.BurnFrom(&_ServiceChainToken.TransactOpts, account, amount)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 subtractedValue) returns(bool)
func (_ServiceChainToken *ServiceChainTokenTransactor) DecreaseAllowance(opts *bind.TransactOpts, spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _ServiceChainToken.contract.Transact(opts, "decreaseAllowance", spender, subtractedValue)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 subtractedValue) returns(bool)
func (_ServiceChainToken *ServiceChainTokenSession) DecreaseAllowance(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _ServiceChainToken.Contract.DecreaseAllowance(&_ServiceChainToken.TransactOpts, spender, subtractedValue)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 subtractedValue) returns(bool)
func (_ServiceChainToken *ServiceChainTokenTransactorSession) DecreaseAllowance(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _ServiceChainToken.Contract.DecreaseAllowance(&_ServiceChainToken.TransactOpts, spender, subtractedValue)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 addedValue) returns(bool)
func (_ServiceChainToken *ServiceChainTokenTransactor) IncreaseAllowance(opts *bind.TransactOpts, spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _ServiceChainToken.contract.Transact(opts, "increaseAllowance", spender, addedValue)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 addedValue) returns(bool)
func (_ServiceChainToken *ServiceChainTokenSession) IncreaseAllowance(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _ServiceChainToken.Contract.IncreaseAllowance(&_ServiceChainToken.TransactOpts, spender, addedValue)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 addedValue) returns(bool)
func (_ServiceChainToken *ServiceChainTokenTransactorSession) IncreaseAllowance(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _ServiceChainToken.Contract.IncreaseAllowance(&_ServiceChainToken.TransactOpts, spender, addedValue)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address account, uint256 amount) returns(bool)
func (_ServiceChainToken *ServiceChainTokenTransactor) Mint(opts *bind.TransactOpts, account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ServiceChainToken.contract.Transact(opts, "mint", account, amount)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address account, uint256 amount) returns(bool)
func (_ServiceChainToken *ServiceChainTokenSession) Mint(account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ServiceChainToken.Contract.Mint(&_ServiceChainToken.TransactOpts, account, amount)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address account, uint256 amount) returns(bool)
func (_ServiceChainToken *ServiceChainTokenTransactorSession) Mint(account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ServiceChainToken.Contract.Mint(&_ServiceChainToken.TransactOpts, account, amount)
}

// RenounceMinter is a paid mutator transaction binding the contract method 0x98650275.
//
// Solidity: function renounceMinter() returns()
func (_ServiceChainToken *ServiceChainTokenTransactor) RenounceMinter(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ServiceChainToken.contract.Transact(opts, "renounceMinter")
}

// RenounceMinter is a paid mutator transaction binding the contract method 0x98650275.
//
// Solidity: function renounceMinter() returns()
func (_ServiceChainToken *ServiceChainTokenSession) RenounceMinter() (*types.Transaction, error) {
	return _ServiceChainToken.Contract.RenounceMinter(&_ServiceChainToken.TransactOpts)
}

// RenounceMinter is a paid mutator transaction binding the contract method 0x98650275.
//
// Solidity: function renounceMinter() returns()
func (_ServiceChainToken *ServiceChainTokenTransactorSession) RenounceMinter() (*types.Transaction, error) {
	return _ServiceChainToken.Contract.RenounceMinter(&_ServiceChainToken.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ServiceChainToken *ServiceChainTokenTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ServiceChainToken.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ServiceChainToken *ServiceChainTokenSession) RenounceOwnership() (*types.Transaction, error) {
	return _ServiceChainToken.Contract.RenounceOwnership(&_ServiceChainToken.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ServiceChainToken *ServiceChainTokenTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _ServiceChainToken.Contract.RenounceOwnership(&_ServiceChainToken.TransactOpts)
}

// RequestValueTransfer is a paid mutator transaction binding the contract method 0x3dc3c9e1.
//
// Solidity: function requestValueTransfer(uint256 _amount, address _to, uint256 _feeLimit, bytes _extraData) returns()
func (_ServiceChainToken *ServiceChainTokenTransactor) RequestValueTransfer(opts *bind.TransactOpts, _amount *big.Int, _to common.Address, _feeLimit *big.Int, _extraData []byte) (*types.Transaction, error) {
	return _ServiceChainToken.contract.Transact(opts, "requestValueTransfer", _amount, _to, _feeLimit, _extraData)
}

// RequestValueTransfer is a paid mutator transaction binding the contract method 0x3dc3c9e1.
//
// Solidity: function requestValueTransfer(uint256 _amount, address _to, uint256 _feeLimit, bytes _extraData) returns()
func (_ServiceChainToken *ServiceChainTokenSession) RequestValueTransfer(_amount *big.Int, _to common.Address, _feeLimit *big.Int, _extraData []byte) (*types.Transaction, error) {
	return _ServiceChainToken.Contract.RequestValueTransfer(&_ServiceChainToken.TransactOpts, _amount, _to, _feeLimit, _extraData)
}

// RequestValueTransfer is a paid mutator transaction binding the contract method 0x3dc3c9e1.
//
// Solidity: function requestValueTransfer(uint256 _amount, address _to, uint256 _feeLimit, bytes _extraData) returns()
func (_ServiceChainToken *ServiceChainTokenTransactorSession) RequestValueTransfer(_amount *big.Int, _to common.Address, _feeLimit *big.Int, _extraData []byte) (*types.Transaction, error) {
	return _ServiceChainToken.Contract.RequestValueTransfer(&_ServiceChainToken.TransactOpts, _amount, _to, _feeLimit, _extraData)
}

// SetBridge is a paid mutator transaction binding the contract method 0x8dd14802.
//
// Solidity: function setBridge(address _bridge) returns()
func (_ServiceChainToken *ServiceChainTokenTransactor) SetBridge(opts *bind.TransactOpts, _bridge common.Address) (*types.Transaction, error) {
	return _ServiceChainToken.contract.Transact(opts, "setBridge", _bridge)
}

// SetBridge is a paid mutator transaction binding the contract method 0x8dd14802.
//
// Solidity: function setBridge(address _bridge) returns()
func (_ServiceChainToken *ServiceChainTokenSession) SetBridge(_bridge common.Address) (*types.Transaction, error) {
	return _ServiceChainToken.Contract.SetBridge(&_ServiceChainToken.TransactOpts, _bridge)
}

// SetBridge is a paid mutator transaction binding the contract method 0x8dd14802.
//
// Solidity: function setBridge(address _bridge) returns()
func (_ServiceChainToken *ServiceChainTokenTransactorSession) SetBridge(_bridge common.Address) (*types.Transaction, error) {
	return _ServiceChainToken.Contract.SetBridge(&_ServiceChainToken.TransactOpts, _bridge)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_ServiceChainToken *ServiceChainTokenTransactor) Transfer(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ServiceChainToken.contract.Transact(opts, "transfer", recipient, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_ServiceChainToken *ServiceChainTokenSession) Transfer(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ServiceChainToken.Contract.Transfer(&_ServiceChainToken.TransactOpts, recipient, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_ServiceChainToken *ServiceChainTokenTransactorSession) Transfer(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ServiceChainToken.Contract.Transfer(&_ServiceChainToken.TransactOpts, recipient, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_ServiceChainToken *ServiceChainTokenTransactor) TransferFrom(opts *bind.TransactOpts, sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ServiceChainToken.contract.Transact(opts, "transferFrom", sender, recipient, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_ServiceChainToken *ServiceChainTokenSession) TransferFrom(sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ServiceChainToken.Contract.TransferFrom(&_ServiceChainToken.TransactOpts, sender, recipient, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_ServiceChainToken *ServiceChainTokenTransactorSession) TransferFrom(sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ServiceChainToken.Contract.TransferFrom(&_ServiceChainToken.TransactOpts, sender, recipient, amount)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ServiceChainToken *ServiceChainTokenTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _ServiceChainToken.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ServiceChainToken *ServiceChainTokenSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ServiceChainToken.Contract.TransferOwnership(&_ServiceChainToken.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ServiceChainToken *ServiceChainTokenTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ServiceChainToken.Contract.TransferOwnership(&_ServiceChainToken.TransactOpts, newOwner)
}

// ServiceChainTokenApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the ServiceChainToken contract.
type ServiceChainTokenApprovalIterator struct {
	Event *ServiceChainTokenApproval // Event containing the contract specifics and raw log

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
func (it *ServiceChainTokenApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ServiceChainTokenApproval)
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
		it.Event = new(ServiceChainTokenApproval)
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
func (it *ServiceChainTokenApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ServiceChainTokenApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ServiceChainTokenApproval represents a Approval event raised by the ServiceChainToken contract.
type ServiceChainTokenApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_ServiceChainToken *ServiceChainTokenFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*ServiceChainTokenApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _ServiceChainToken.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &ServiceChainTokenApprovalIterator{contract: _ServiceChainToken.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_ServiceChainToken *ServiceChainTokenFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *ServiceChainTokenApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _ServiceChainToken.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ServiceChainTokenApproval)
				if err := _ServiceChainToken.contract.UnpackLog(event, "Approval", log); err != nil {
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

// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_ServiceChainToken *ServiceChainTokenFilterer) ParseApproval(log types.Log) (*ServiceChainTokenApproval, error) {
	event := new(ServiceChainTokenApproval)
	if err := _ServiceChainToken.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ServiceChainTokenMinterAddedIterator is returned from FilterMinterAdded and is used to iterate over the raw logs and unpacked data for MinterAdded events raised by the ServiceChainToken contract.
type ServiceChainTokenMinterAddedIterator struct {
	Event *ServiceChainTokenMinterAdded // Event containing the contract specifics and raw log

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
func (it *ServiceChainTokenMinterAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ServiceChainTokenMinterAdded)
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
		it.Event = new(ServiceChainTokenMinterAdded)
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
func (it *ServiceChainTokenMinterAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ServiceChainTokenMinterAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ServiceChainTokenMinterAdded represents a MinterAdded event raised by the ServiceChainToken contract.
type ServiceChainTokenMinterAdded struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterMinterAdded is a free log retrieval operation binding the contract event 0x6ae172837ea30b801fbfcdd4108aa1d5bf8ff775444fd70256b44e6bf3dfc3f6.
//
// Solidity: event MinterAdded(address indexed account)
func (_ServiceChainToken *ServiceChainTokenFilterer) FilterMinterAdded(opts *bind.FilterOpts, account []common.Address) (*ServiceChainTokenMinterAddedIterator, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _ServiceChainToken.contract.FilterLogs(opts, "MinterAdded", accountRule)
	if err != nil {
		return nil, err
	}
	return &ServiceChainTokenMinterAddedIterator{contract: _ServiceChainToken.contract, event: "MinterAdded", logs: logs, sub: sub}, nil
}

// WatchMinterAdded is a free log subscription operation binding the contract event 0x6ae172837ea30b801fbfcdd4108aa1d5bf8ff775444fd70256b44e6bf3dfc3f6.
//
// Solidity: event MinterAdded(address indexed account)
func (_ServiceChainToken *ServiceChainTokenFilterer) WatchMinterAdded(opts *bind.WatchOpts, sink chan<- *ServiceChainTokenMinterAdded, account []common.Address) (event.Subscription, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _ServiceChainToken.contract.WatchLogs(opts, "MinterAdded", accountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ServiceChainTokenMinterAdded)
				if err := _ServiceChainToken.contract.UnpackLog(event, "MinterAdded", log); err != nil {
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

// ParseMinterAdded is a log parse operation binding the contract event 0x6ae172837ea30b801fbfcdd4108aa1d5bf8ff775444fd70256b44e6bf3dfc3f6.
//
// Solidity: event MinterAdded(address indexed account)
func (_ServiceChainToken *ServiceChainTokenFilterer) ParseMinterAdded(log types.Log) (*ServiceChainTokenMinterAdded, error) {
	event := new(ServiceChainTokenMinterAdded)
	if err := _ServiceChainToken.contract.UnpackLog(event, "MinterAdded", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ServiceChainTokenMinterRemovedIterator is returned from FilterMinterRemoved and is used to iterate over the raw logs and unpacked data for MinterRemoved events raised by the ServiceChainToken contract.
type ServiceChainTokenMinterRemovedIterator struct {
	Event *ServiceChainTokenMinterRemoved // Event containing the contract specifics and raw log

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
func (it *ServiceChainTokenMinterRemovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ServiceChainTokenMinterRemoved)
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
		it.Event = new(ServiceChainTokenMinterRemoved)
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
func (it *ServiceChainTokenMinterRemovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ServiceChainTokenMinterRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ServiceChainTokenMinterRemoved represents a MinterRemoved event raised by the ServiceChainToken contract.
type ServiceChainTokenMinterRemoved struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterMinterRemoved is a free log retrieval operation binding the contract event 0xe94479a9f7e1952cc78f2d6baab678adc1b772d936c6583def489e524cb66692.
//
// Solidity: event MinterRemoved(address indexed account)
func (_ServiceChainToken *ServiceChainTokenFilterer) FilterMinterRemoved(opts *bind.FilterOpts, account []common.Address) (*ServiceChainTokenMinterRemovedIterator, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _ServiceChainToken.contract.FilterLogs(opts, "MinterRemoved", accountRule)
	if err != nil {
		return nil, err
	}
	return &ServiceChainTokenMinterRemovedIterator{contract: _ServiceChainToken.contract, event: "MinterRemoved", logs: logs, sub: sub}, nil
}

// WatchMinterRemoved is a free log subscription operation binding the contract event 0xe94479a9f7e1952cc78f2d6baab678adc1b772d936c6583def489e524cb66692.
//
// Solidity: event MinterRemoved(address indexed account)
func (_ServiceChainToken *ServiceChainTokenFilterer) WatchMinterRemoved(opts *bind.WatchOpts, sink chan<- *ServiceChainTokenMinterRemoved, account []common.Address) (event.Subscription, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _ServiceChainToken.contract.WatchLogs(opts, "MinterRemoved", accountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ServiceChainTokenMinterRemoved)
				if err := _ServiceChainToken.contract.UnpackLog(event, "MinterRemoved", log); err != nil {
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

// ParseMinterRemoved is a log parse operation binding the contract event 0xe94479a9f7e1952cc78f2d6baab678adc1b772d936c6583def489e524cb66692.
//
// Solidity: event MinterRemoved(address indexed account)
func (_ServiceChainToken *ServiceChainTokenFilterer) ParseMinterRemoved(log types.Log) (*ServiceChainTokenMinterRemoved, error) {
	event := new(ServiceChainTokenMinterRemoved)
	if err := _ServiceChainToken.contract.UnpackLog(event, "MinterRemoved", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ServiceChainTokenOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the ServiceChainToken contract.
type ServiceChainTokenOwnershipTransferredIterator struct {
	Event *ServiceChainTokenOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *ServiceChainTokenOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ServiceChainTokenOwnershipTransferred)
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
		it.Event = new(ServiceChainTokenOwnershipTransferred)
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
func (it *ServiceChainTokenOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ServiceChainTokenOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ServiceChainTokenOwnershipTransferred represents a OwnershipTransferred event raised by the ServiceChainToken contract.
type ServiceChainTokenOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ServiceChainToken *ServiceChainTokenFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*ServiceChainTokenOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ServiceChainToken.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &ServiceChainTokenOwnershipTransferredIterator{contract: _ServiceChainToken.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ServiceChainToken *ServiceChainTokenFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ServiceChainTokenOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ServiceChainToken.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ServiceChainTokenOwnershipTransferred)
				if err := _ServiceChainToken.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_ServiceChainToken *ServiceChainTokenFilterer) ParseOwnershipTransferred(log types.Log) (*ServiceChainTokenOwnershipTransferred, error) {
	event := new(ServiceChainTokenOwnershipTransferred)
	if err := _ServiceChainToken.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ServiceChainTokenTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the ServiceChainToken contract.
type ServiceChainTokenTransferIterator struct {
	Event *ServiceChainTokenTransfer // Event containing the contract specifics and raw log

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
func (it *ServiceChainTokenTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ServiceChainTokenTransfer)
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
		it.Event = new(ServiceChainTokenTransfer)
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
func (it *ServiceChainTokenTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ServiceChainTokenTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ServiceChainTokenTransfer represents a Transfer event raised by the ServiceChainToken contract.
type ServiceChainTokenTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_ServiceChainToken *ServiceChainTokenFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*ServiceChainTokenTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ServiceChainToken.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &ServiceChainTokenTransferIterator{contract: _ServiceChainToken.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_ServiceChainToken *ServiceChainTokenFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *ServiceChainTokenTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ServiceChainToken.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ServiceChainTokenTransfer)
				if err := _ServiceChainToken.contract.UnpackLog(event, "Transfer", log); err != nil {
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

// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_ServiceChainToken *ServiceChainTokenFilterer) ParseTransfer(log types.Log) (*ServiceChainTokenTransfer, error) {
	event := new(ServiceChainTokenTransfer)
	if err := _ServiceChainToken.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	return event, nil
}
