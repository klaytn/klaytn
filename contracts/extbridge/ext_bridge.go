// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package extbridge

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

// AddressABI is the input ABI used to generate the binding from.
const AddressABI = "[]"

// AddressBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const AddressBinRuntime = `0x73000000000000000000000000000000000000000030146080604052600080fd00a165627a7a72305820145ddad4b6b54b6baf7ccf337abfe1e51b0c28e420d1e2d4eeb65780f52c2b770029`

// AddressBin is the compiled bytecode used for deploying new contracts.
const AddressBin = `0x604c602c600b82828239805160001a60731460008114601c57601e565bfe5b5030600052607381538281f30073000000000000000000000000000000000000000030146080604052600080fd00a165627a7a72305820145ddad4b6b54b6baf7ccf337abfe1e51b0c28e420d1e2d4eeb65780f52c2b770029`

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

// BridgeABI is the input ABI used to generate the binding from.
const BridgeABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"_tokenAddress\",\"type\":\"address\"},{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_tokenId\",\"type\":\"uint256\"},{\"name\":\"_extraData\",\"type\":\"uint256[]\"}],\"name\":\"requestERC721Transfer\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"operators\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_fee\",\"type\":\"uint256\"},{\"name\":\"_requestNonce\",\"type\":\"uint64\"}],\"name\":\"setKLAYFee\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"isRunning\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_requestTxHash\",\"type\":\"bytes32\"},{\"name\":\"_from\",\"type\":\"address\"},{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_tokenAddress\",\"type\":\"address\"},{\"name\":\"_tokenId\",\"type\":\"uint256\"},{\"name\":\"_requestedNonce\",\"type\":\"uint64\"},{\"name\":\"_requestedBlockNumber\",\"type\":\"uint64\"},{\"name\":\"_tokenURI\",\"type\":\"string\"},{\"name\":\"_extraData\",\"type\":\"uint256[]\"}],\"name\":\"handleERC721Transfer\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_token\",\"type\":\"address\"},{\"name\":\"_fee\",\"type\":\"uint256\"},{\"name\":\"_requestNonce\",\"type\":\"uint64\"}],\"name\":\"setERC20Fee\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"sequentialHandledRequestBlockNumber\",\"outputs\":[{\"name\":\"\",\"type\":\"uint64\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_from\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"},{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_feeLimit\",\"type\":\"uint256\"},{\"name\":\"_extraData\",\"type\":\"uint256[]\"}],\"name\":\"onERC20Received\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_operator\",\"type\":\"address\"}],\"name\":\"registerOperator\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"counterpartBridge\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_requestTxHash\",\"type\":\"bytes32\"},{\"name\":\"_from\",\"type\":\"address\"},{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_tokenAddress\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"},{\"name\":\"_requestedNonce\",\"type\":\"uint64\"},{\"name\":\"_requestedBlockNumber\",\"type\":\"uint64\"},{\"name\":\"_extraData\",\"type\":\"uint256[]\"}],\"name\":\"handleERC20Transfer\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_token\",\"type\":\"address\"},{\"name\":\"_cToken\",\"type\":\"address\"}],\"name\":\"registerToken\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"feeOfERC20\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint8\"}],\"name\":\"operatorThresholds\",\"outputs\":[{\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"maxHandledRequestedNonce\",\"outputs\":[{\"name\":\"\",\"type\":\"uint64\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint64\"}],\"name\":\"handledHistory\",\"outputs\":[{\"name\":\"\",\"type\":\"uint64\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"modeMintBurn\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_tokenAddress\",\"type\":\"address\"},{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"},{\"name\":\"_feeLimit\",\"type\":\"uint256\"},{\"name\":\"_extraData\",\"type\":\"uint256[]\"}],\"name\":\"requestERC20Transfer\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"requestNonce\",\"outputs\":[{\"name\":\"\",\"type\":\"uint64\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_bridge\",\"type\":\"address\"}],\"name\":\"setCounterPartBridge\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"handledRequestTx\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"isOwner\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\"},{\"name\":\"\",\"type\":\"address\"}],\"name\":\"votes\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint64\"}],\"name\":\"closedValueTransferVotes\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_from\",\"type\":\"address\"},{\"name\":\"_tokenId\",\"type\":\"uint256\"},{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_extraData\",\"type\":\"uint256[]\"}],\"name\":\"onERC721Received\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"votesCounts\",\"outputs\":[{\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"configurationNonce\",\"outputs\":[{\"name\":\"\",\"type\":\"uint64\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"feeReceiver\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_token\",\"type\":\"address\"}],\"name\":\"deregisterToken\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"feeOfKLAY\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_status\",\"type\":\"bool\"}],\"name\":\"start\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_operator\",\"type\":\"address\"}],\"name\":\"deregisterOperator\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"chargeWithoutEvent\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"allowedTokens\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"},{\"name\":\"_extraData\",\"type\":\"uint256[]\"}],\"name\":\"requestKLAYTransfer\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_requestTxHash\",\"type\":\"bytes32\"},{\"name\":\"_from\",\"type\":\"address\"},{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"},{\"name\":\"_requestedNonce\",\"type\":\"uint64\"},{\"name\":\"_requestedBlockNumber\",\"type\":\"uint64\"},{\"name\":\"_extraData\",\"type\":\"uint256[]\"}],\"name\":\"handleKLAYTransfer\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_voteType\",\"type\":\"uint8\"},{\"name\":\"_threshold\",\"type\":\"uint8\"}],\"name\":\"setOperatorThreshold\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_feeReceiver\",\"type\":\"address\"}],\"name\":\"setFeeReceiver\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"sequentialHandleNonce\",\"outputs\":[{\"name\":\"\",\"type\":\"uint64\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"VERSION\",\"outputs\":[{\"name\":\"\",\"type\":\"uint64\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"_modeMintBurn\",\"type\":\"bool\"}],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"constructor\"},{\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"tokenType\",\"type\":\"uint8\"},{\"indexed\":false,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"tokenAddress\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"valueOrTokenId\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"requestNonce\",\"type\":\"uint64\"},{\"indexed\":false,\"name\":\"uri\",\"type\":\"string\"},{\"indexed\":false,\"name\":\"fee\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"extraData\",\"type\":\"uint256[]\"}],\"name\":\"RequestValueTransfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"requestTxHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"tokenType\",\"type\":\"uint8\"},{\"indexed\":false,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"tokenAddress\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"valueOrTokenId\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"handleNonce\",\"type\":\"uint64\"},{\"indexed\":false,\"name\":\"extraData\",\"type\":\"uint256[]\"}],\"name\":\"HandleValueTransfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"fee\",\"type\":\"uint256\"}],\"name\":\"KLAYFeeChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"token\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"fee\",\"type\":\"uint256\"}],\"name\":\"ERC20FeeChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"feeReceiver\",\"type\":\"address\"}],\"name\":\"FeeReceiverChanged\",\"type\":\"event\"}]"

// BridgeBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const BridgeBinRuntime = `0x6080604052600436106101f55763ffffffff60e060020a600035041663098540b8811461022f57806313e7c9d8146102685780631a2ae53e1461029d5780632014e5d1146102c25780632e4bb197146102d75780632f88396c146103a757806330f436d0146103d8578063349449b21461040a5780633682a4501461047f5780633a348533146104a0578063452fbc62146104d15780634739f7e514610561578063488af871146105885780635526f76b146105bb57806360f17a9b146105ec57806367bda76b146106015780636e176ec214610623578063715018a61461063857806374f41f151461064d5780637c1a03021461068a57806387b04c551461069f5780638a75eee2146106c05780638da5cb5b146106d85780638f32d59b146106ed5780639386775a146107025780639832c1d714610726578063a319de9e14610748578063a4155e5d146107bb578063ac6fff0b146107d3578063b3f00674146107e8578063bab2af1d146107fd578063c263b5d61461081e578063c877cf3714610833578063d8cf98ca1461084d578063dd9222d61461086e578063e744092e14610876578063eb7f0cbe14610897578063ee273ad7146108bb578063ee2aec6514610943578063efdcd97414610964578063f2fde38b14610985578063febb0cc2146109a6578063ffa1ad74146109bb575b61022d336001546000604051908082528060200260200182016040528015610227578160200160208202803883390190505b506109d0565b005b34801561023b57600080fd5b5061022d600160a060020a0360048035821691602480359091169160443591606435908101910135610c23565b34801561027457600080fd5b50610289600160a060020a0360043516610ce9565b604080519115158252519081900360200190f35b3480156102a957600080fd5b5061022d60043567ffffffffffffffff60243516610cfe565b3480156102ce57600080fd5b50610289610d3d565b3480156102e357600080fd5b50604080516020601f60e43560048181013592830184900484028501840190955281845261022d94803594600160a060020a0360248035821696604435831696606435909316956084359567ffffffffffffffff60a43581169660c435909116953695929461010494929392909201918190840183828082843750506040805187358901803560208181028481018201909552818452989b9a998901989297509082019550935083925085019084908082843750949750610d4d9650505050505050565b3480156103b357600080fd5b5061022d600160a060020a036004351660243567ffffffffffffffff60443516611007565b3480156103e457600080fd5b506103ed611048565b6040805167ffffffffffffffff9092168252519081900360200190f35b34801561041657600080fd5b50604080516084356004818101356020818102858101820190965281855261022d95600160a060020a038435811696602480359760443590931696606435963696919560a495909301929091829190850190849080828437509497506110589650505050505050565b34801561048b57600080fd5b5061022d600160a060020a0360043516611066565b3480156104ac57600080fd5b506104b561109d565b60408051600160a060020a039092168252519081900360200190f35b3480156104dd57600080fd5b50604080516020600460e43581810135838102808601850190965280855261022d9583359560248035600160a060020a0390811697604435821697606435909216966084359660a43567ffffffffffffffff9081169760c4359091169636966101049591940192918291850190849080828437509497506110b19650505050505050565b34801561056d57600080fd5b5061022d600160a060020a0360043581169060243516611306565b34801561059457600080fd5b506105a9600160a060020a0360043516611354565b60408051918252519081900360200190f35b3480156105c757600080fd5b506105d660ff60043516611366565b6040805160ff9092168252519081900360200190f35b3480156105f857600080fd5b506103ed61137b565b34801561060d57600080fd5b506103ed67ffffffffffffffff6004351661139f565b34801561062f57600080fd5b506102896113bb565b34801561064457600080fd5b5061022d6113c4565b34801561065957600080fd5b5061022d600160a060020a03600480358216916024803590911691604435916064359160843591820191013561142e565b34801561069657600080fd5b506103ed611516565b3480156106ab57600080fd5b5061022d600160a060020a036004351661152d565b3480156106cc57600080fd5b50610289600435611575565b3480156106e457600080fd5b506104b561158a565b3480156106f957600080fd5b50610289611599565b34801561070e57600080fd5b50610289600435600160a060020a03602435166115aa565b34801561073257600080fd5b5061028967ffffffffffffffff600435166115ca565b34801561075457600080fd5b50604080516064356004818101356020818102858101820190965281855261022d95600160a060020a038435811696602480359760443590931696369690956084959394909201929091829190850190849080828437509497506115df9650505050505050565b3480156107c757600080fd5b506105d66004356115f2565b3480156107df57600080fd5b506103ed611607565b3480156107f457600080fd5b506104b5611617565b34801561080957600080fd5b5061022d600160a060020a0360043516611626565b34801561082a57600080fd5b506105a961166d565b34801561083f57600080fd5b5061022d6004351515611673565b34801561085957600080fd5b5061022d600160a060020a03600435166116b6565b61022d6116ea565b34801561088257600080fd5b506104b5600160a060020a03600435166116ec565b61022d60048035600160a060020a0316906024803591604435918201910135611707565b3480156108c757600080fd5b506040805160c4356004818101356020818102858101820190965281855261022d95833595600160a060020a0360248035821697604435909216966064359667ffffffffffffffff60843581169760a435909116963696939560e495019291829190850190849080828437509497506117539650505050505050565b34801561094f57600080fd5b5061022d60ff600435811690602435166118ab565b34801561097057600080fd5b5061022d600160a060020a03600435166118f8565b34801561099157600080fd5b5061022d600160a060020a0360043516611917565b3480156109b257600080fd5b506103ed611933565b3480156109c757600080fd5b506103ed61194f565b600b5460009060a860020a900460ff161515610a36576040805160e560020a62461bcd02815260206004820152600e60248201527f73746f7070656420627269646765000000000000000000000000000000000000604482015290519081900360640190fd5b348310610a8d576040805160e560020a62461bcd02815260206004820152601360248201527f696e73756666696369656e7420616d6f756e7400000000000000000000000000604482015290519081900360640190fd5b610a9683611954565b90507f09a722fc78123f7b2ba788ae83b01221477757cfb4b83302421e7c5c1fd01da96000338682610ace348963ffffffff611a8a16565b600b60169054906101000a900467ffffffffffffffff16878960405180896002811115610af757fe5b60ff16815260200188600160a060020a0316600160a060020a0316815260200187600160a060020a0316600160a060020a0316815260200186600160a060020a0316600160a060020a031681526020018581526020018467ffffffffffffffff1667ffffffffffffffff1681526020018060200184815260200180602001838103835260008152602001602001838103825284818151815260200191508051906020019060200280838360005b83811015610bbc578181015183820152602001610ba4565b505050509050019a505050505050505050505060405180910390a15050600b8054600167ffffffffffffffff60b060020a80840482169290920116027fffff0000000000000000ffffffffffffffffffffffffffffffffffffffffffff9091161790555050565b604080517f23b872dd000000000000000000000000000000000000000000000000000000008152336004820152306024820152604481018590529051600160a060020a038716916323b872dd91606480830192600092919082900301818387803b158015610c9057600080fd5b505af1158015610ca4573d6000803e3d6000fd5b50505050610ce28533868686868080602002602001604051908101604052809392919081815260200183836020028082843750611aa1945050505050565b5050505050565b60046020526000908152604090205460ff1681565b3360009081526004602052604090205460ff161515610d1c57600080fd5b610d2581611e8b565b1515610d3057610d39565b610d3982611f57565b5050565b600b5460a860020a900460ff1681565b3360009081526004602052604090205460ff161515610d6b57600080fd5b610d7484611f8a565b1515610d7f57610ffc565b604080518a815260026020808301829052600160a060020a03808d1694840194909452838b166060840152928916608083015260a0820188905267ffffffffffffffff871660c083015261010060e0830181815285519184019190915284517f91f1d4b9864c9c091ade654e58899919d2a56e194153899c7e292cd3ae2ac47b948e948e938e938e938e938e938c9391929091610120840191808601910280838360005b83811015610e3b578181015183820152602001610e23565b50505050905001995050505050505050505060405180910390a1610e5e89612055565b610e688484612070565b600b5460ff1615610f725785600160a060020a03166350bb4e7f8887856040518463ffffffff1660e060020a0281526004018084600160a060020a0316600160a060020a0316815260200183815260200180602001828103825283818151815260200191508051906020019080838360005b83811015610ef2578181015183820152602001610eda565b50505050905090810190601f168015610f1f5780820380516001836020036101000a031916815260200191505b50945050505050602060405180830381600087803b158015610f4057600080fd5b505af1158015610f54573d6000803e3d6000fd5b505050506040513d6020811015610f6a57600080fd5b50610ffc9050565b604080517f42842e0e000000000000000000000000000000000000000000000000000000008152306004820152600160a060020a038981166024830152604482018890529151918816916342842e0e9160648082019260009290919082900301818387803b158015610fe357600080fd5b505af1158015610ff7573d6000803e3d6000fd5b505050505b505050505050505050565b3360009081526004602052604090205460ff16151561102557600080fd5b61102e81611e8b565b151561103957611043565b61104383836121dd565b505050565b600c5467ffffffffffffffff1681565b610ce2338685878686612232565b61106e611599565b151561107957600080fd5b600160a060020a03166000908152600460205260409020805460ff19166001179055565b600b546101009004600160a060020a031681565b3360009081526004602052604090205460ff1615156110cf57600080fd5b6110d883611f8a565b15156110e3576112fc565b6040805189815260016020808301829052600160a060020a03808c1694840194909452838a166060840152928816608083015260a0820187905267ffffffffffffffff861660c083015261010060e0830181815285519184019190915284517f91f1d4b9864c9c091ade654e58899919d2a56e194153899c7e292cd3ae2ac47b948d948d938d938d938d938d938c9391929091610120840191808601910280838360005b8381101561119f578181015183820152602001611187565b50505050905001995050505050505050505060405180910390a16111c288612055565b6111cc8383612070565b600b5460ff161561126c5784600160a060020a03166340c10f1987866040518363ffffffff1660e060020a0281526004018083600160a060020a0316600160a060020a0316815260200182815260200192505050602060405180830381600087803b15801561123a57600080fd5b505af115801561124e573d6000803e3d6000fd5b505050506040513d602081101561126457600080fd5b506112fc9050565b84600160a060020a031663a9059cbb87866040518363ffffffff1660e060020a0281526004018083600160a060020a0316600160a060020a0316815260200182815260200192505050602060405180830381600087803b1580156112cf57600080fd5b505af11580156112e3573d6000803e3d6000fd5b505050506040513d60208110156112f957600080fd5b50505b5050505050505050565b61130e611599565b151561131957600080fd5b600160a060020a039182166000908152600e60205260409020805473ffffffffffffffffffffffffffffffffffffffff191691909216179055565b60026020526000908152604090205481565b60086020526000908152604090205460ff1681565b600c54700100000000000000000000000000000000900467ffffffffffffffff1681565b600d6020526000908152604090205467ffffffffffffffff1681565b600b5460ff1681565b6113cc611599565b15156113d757600080fd5b600354604051600091600160a060020a0316907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0908390a36003805473ffffffffffffffffffffffffffffffffffffffff19169055565b600160a060020a0386166323b872dd333061144f888863ffffffff61255b16565b6040805160e060020a63ffffffff8716028152600160a060020a0394851660048201529290931660248301526044820152905160648083019260209291908290030181600087803b1580156114a357600080fd5b505af11580156114b7573d6000803e3d6000fd5b505050506040513d60208110156114cd57600080fd5b5050604080516020838102828101820190935283825261150e92899233928a928a928a9290918a918a91829185019084908082843750612232945050505050565b505050505050565b600b5460b060020a900467ffffffffffffffff1681565b611535611599565b151561154057600080fd5b600b8054600160a060020a039092166101000274ffffffffffffffffffffffffffffffffffffffff0019909216919091179055565b600a6020526000908152604090205460ff1681565b600354600160a060020a031690565b600354600160a060020a0316331490565b600560209081526000928352604080842090915290825290205460ff1681565b60076020526000908152604090205460ff1681565b6115ec3385848685611aa1565b50505050565b60066020526000908152604090205460ff1681565b60095467ffffffffffffffff1681565b600054600160a060020a031681565b61162e611599565b151561163957600080fd5b600160a060020a03166000908152600e60205260409020805473ffffffffffffffffffffffffffffffffffffffff19169055565b60015481565b61167b611599565b151561168657600080fd5b600b805491151560a860020a0275ff00000000000000000000000000000000000000000019909216919091179055565b6116be611599565b15156116c957600080fd5b600160a060020a03166000908152600460205260409020805460ff19169055565b565b600e60205260009081526040902054600160a060020a031681565b6000611719348563ffffffff611a8a16565b9050610ce28582858580806020026020016040519081016040528093929190818152602001838360200280828437506109d0945050505050565b3360009081526004602052604090205460ff16151561177157600080fd5b61177a83611f8a565b1515611785576118a2565b6040805188815260006020808301829052600160a060020a03808b169484019490945292881660608301526080820181905260a0820187905267ffffffffffffffff861660c083015261010060e0830181815285519184019190915284517f91f1d4b9864c9c091ade654e58899919d2a56e194153899c7e292cd3ae2ac47b948c948c938c9386938d938d938c939192909161012084019180860191028083838b5b8381101561183f578181015183820152602001611827565b50505050905001995050505050505050505060405180910390a161186287612055565b61186c8383612070565b604051600160a060020a0386169085156108fc029086906000818181858888f193505050501580156112fc573d6000803e3d6000fd5b50505050505050565b6118b3611599565b15156118be57600080fd5b80600860008460028111156118cf57fe5b60ff90811682526020820192909252604001600020805460ff1916929091169190911790555050565b611900611599565b151561190b57600080fd5b61191481612574565b50565b61191f611599565b151561192a57600080fd5b611914816125c9565b600c5468010000000000000000900467ffffffffffffffff1681565b600181565b60015460008054909190600160a060020a0316158015906119755750600081115b15611a5157808310156119d2576040805160e560020a62461bcd02815260206004820152601560248201527f696e73756666696369656e74206665654c696d69740000000000000000000000604482015290519081900360640190fd5b60008054604051600160a060020a039091169183156108fc02918491818181858888f19350505050158015611a0b573d6000803e3d6000fd5b50336108fc611a20858463ffffffff611a8a16565b6040518115909202916000818181858888f19350505050158015611a48573d6000803e3d6000fd5b50809150611a84565b604051339084156108fc029085906000818181858888f19350505050158015611a7e573d6000803e3d6000fd5b50600091505b50919050565b60008083831115611a9a57600080fd5b5050900390565b600b5460609060a860020a900460ff161515611b07576040805160e560020a62461bcd02815260206004820152600e60248201527f73746f7070656420627269646765000000000000000000000000000000000000604482015290519081900360640190fd5b600160a060020a038681166000908152600e6020526040902054161515611b78576040805160e560020a62461bcd02815260206004820152600d60248201527f696e76616c696420746f6b656e00000000000000000000000000000000000000604482015290519081900360640190fd5b85600160a060020a031663c87b56dd846040518263ffffffff1660e060020a02815260040180828152602001915050600060405180830381600087803b158015611bc157600080fd5b505af1158015611bd5573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f191682016040526020811015611bfe57600080fd5b810190808051640100000000811115611c1657600080fd5b82016020810184811115611c2957600080fd5b8151640100000000811182820187101715611c4357600080fd5b5050600b5490945060ff16159250611cba9150505785600160a060020a03166342966c68846040518263ffffffff1660e060020a02815260040180828152602001915050600060405180830381600087803b158015611ca157600080fd5b505af1158015611cb5573d6000803e3d6000fd5b505050505b7f09a722fc78123f7b2ba788ae83b01221477757cfb4b83302421e7c5c1fd01da9600286868987600b60169054906101000a900467ffffffffffffffff168760008a604051808a6002811115611d0c57fe5b60ff16815260200189600160a060020a0316600160a060020a0316815260200188600160a060020a0316600160a060020a0316815260200187600160a060020a0316600160a060020a031681526020018681526020018567ffffffffffffffff1667ffffffffffffffff1681526020018060200184815260200180602001838103835286818151815260200191508051906020019080838360005b83811015611dbf578181015183820152602001611da7565b50505050905090810190601f168015611dec5780820380516001836020036101000a031916815260200191505b508381038252845181528451602091820191808701910280838360005b83811015611e21578181015183820152602001611e09565b505050509050019b50505050505050505050505060405180910390a15050600b8054600167ffffffffffffffff60b060020a80840482169290920116027fffff0000000000000000ffffffffffffffffffffffffffffffffffffffffffff90911617905550505050565b60095460009067ffffffffffffffff838116911614611ef4576040805160e560020a62461bcd02815260206004820152600e60248201527f6e6f6e6365206d69736d61746368000000000000000000000000000000000000604482015290519081900360640190fd5b611f1c6001600036604051808383808284378201915050925050506040518091039020612647565b15611f4e57506009805467ffffffffffffffff198116600167ffffffffffffffff928316810190921617909155611f52565b5060005b919050565b600181905560405181907fa7a33d0996347e1aa55ca2206015b61b9534bdd881d59d59aa680e25eefac36590600090a250565b67ffffffffffffffff811660009081526007602052604081205460ff1615611ffc576040805160e560020a62461bcd02815260206004820152600b60248201527f636c6f73656420766f7465000000000000000000000000000000000000000000604482015290519081900360640190fd5b61202360008036604051808383808284378201915050925050506040518091039020612647565b15611f4e575067ffffffffffffffff81166000908152600760205260409020805460ff19166001908117909155611f52565b6000908152600a60205260409020805460ff19166001179055565b67ffffffffffffffff8281166000818152600d60205260408120805467ffffffffffffffff1916858516179055600c5490927001000000000000000000000000000000009091041610156120ff57600c805477ffffffffffffffff00000000000000000000000000000000191670010000000000000000000000000000000067ffffffffffffffff8616021790555b50600c5468010000000000000000900467ffffffffffffffff165b600c5467ffffffffffffffff700100000000000000000000000000000000909104811690821611801590612169575067ffffffffffffffff8082166000908152600d6020526040812054909116115b156121765760010161211a565b600c805467ffffffffffffffff80841668010000000000000000026fffffffffffffffff0000000000000000199092169190911780835560001990930181166000908152600d602052604090205467ffffffffffffffff1990931692169190911790555050565b600160a060020a0382166000818152600260209081526040918290208490558151928352905183927fdb5ad2e76ae20cfa4e7adbc7305d7538442164d85ead9937c98620a1aa4c255b92908290030190a25050565b600b5460009060a860020a900460ff161515612298576040805160e560020a62461bcd02815260206004820152600e60248201527f73746f7070656420627269646765000000000000000000000000000000000000604482015290519081900360640190fd5b600084116122f0576040805160e560020a62461bcd02815260206004820152600e60248201527f7a65726f206d73672e76616c7565000000000000000000000000000000000000604482015290519081900360640190fd5b600160a060020a038781166000908152600e6020526040902054161515612361576040805160e560020a62461bcd02815260206004820152600d60248201527f696e76616c696420746f6b656e00000000000000000000000000000000000000604482015290519081900360640190fd5b61236c868885612727565b600b5490915060ff16156123dc5786600160a060020a03166342966c68856040518263ffffffff1660e060020a02815260040180828152602001915050600060405180830381600087803b1580156123c357600080fd5b505af11580156123d7573d6000803e3d6000fd5b505050505b7f09a722fc78123f7b2ba788ae83b01221477757cfb4b83302421e7c5c1fd01da9600187878a88600b60169054906101000a900467ffffffffffffffff1687896040518089600281111561242c57fe5b60ff16815260200188600160a060020a0316600160a060020a0316815260200187600160a060020a0316600160a060020a0316815260200186600160a060020a0316600160a060020a031681526020018581526020018467ffffffffffffffff1667ffffffffffffffff1681526020018060200184815260200180602001838103835260008152602001602001838103825284818151815260200191508051906020019060200280838360005b838110156124f15781810151838201526020016124d9565b505050509050019a505050505050505050505060405180910390a15050600b8054600167ffffffffffffffff60b060020a80840482169290920116027fffff0000000000000000ffffffffffffffffffffffffffffffffffffffffffff9091161790555050505050565b60008282018381101561256d57600080fd5b9392505050565b6000805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a038316908117825560405190917f647672599d3468abcfa241a13c9e3d34383caadb5cc80fb67c3cdfcd5f78605991a250565b600160a060020a03811615156125de57600080fd5b600354604051600160a060020a038084169216907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a36003805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a0392909216919091179055565b600081815260056020908152604080832033845290915281205460ff1615156126d35760008281526005602090815260408083203384528252808320805460ff1916600190811790915585845260069092529091205460ff90811691820116116126b057600080fd5b6000828152600660205260409020805460ff8082166001011660ff199091161790555b600860008460028111156126e357fe5b60ff908116825260208083019390935260409182016000908120548682526006909452919091205491811691161061271d57506001612721565b5060005b92915050565b600160a060020a0380831660009081526002602052604081205481549192909116158015906127565750600081115b156128f757808310156127b3576040805160e560020a62461bcd02815260206004820152601560248201527f696e73756666696369656e74206665654c696d69740000000000000000000000604482015290519081900360640190fd5b60008054604080517fa9059cbb000000000000000000000000000000000000000000000000000000008152600160a060020a0392831660048201526024810185905290519187169263a9059cbb926044808401936020939083900390910190829087803b15801561282357600080fd5b505af1158015612837573d6000803e3d6000fd5b505050506040513d602081101561284d57600080fd5b5050600160a060020a03841663a9059cbb8661286f868563ffffffff611a8a16565b6040518363ffffffff1660e060020a0281526004018083600160a060020a0316600160a060020a0316815260200182815260200192505050602060405180830381600087803b1580156128c157600080fd5b505af11580156128d5573d6000803e3d6000fd5b505050506040513d60208110156128eb57600080fd5b5090915081905061298b565b83600160a060020a031663a9059cbb86856040518363ffffffff1660e060020a0281526004018083600160a060020a0316600160a060020a0316815260200182815260200192505050602060405180830381600087803b15801561295a57600080fd5b505af115801561296e573d6000803e3d6000fd5b505050506040513d602081101561298457600080fd5b5060009250505b5093925050505600a165627a7a72305820b26ab42c0688c114123227adc844724b08ff5886ef6f44271ff99c60751076850029`

// BridgeBin is the compiled bytecode used for deploying new contracts.
const BridgeBin = `0x6080604081905260008054600160a060020a0319168155600155600b805460ff19169055602080612ade833981016040819052905160008054600160a060020a031990811682556003805490911633179081905591929091600160a060020a03169082907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0908290a35060005b600260ff821610156100be5760ff81166000908152600860205260409020805460ff191660019081179091550161008c565b50336000908152600460205260409020805460ff19908116600117909155600b805492151560a860020a60ff02199093167501000000000000000000000000000000000000000000179091169190911790556129bf8061011f6000396000f3006080604052600436106101f55763ffffffff60e060020a600035041663098540b8811461022f57806313e7c9d8146102685780631a2ae53e1461029d5780632014e5d1146102c25780632e4bb197146102d75780632f88396c146103a757806330f436d0146103d8578063349449b21461040a5780633682a4501461047f5780633a348533146104a0578063452fbc62146104d15780634739f7e514610561578063488af871146105885780635526f76b146105bb57806360f17a9b146105ec57806367bda76b146106015780636e176ec214610623578063715018a61461063857806374f41f151461064d5780637c1a03021461068a57806387b04c551461069f5780638a75eee2146106c05780638da5cb5b146106d85780638f32d59b146106ed5780639386775a146107025780639832c1d714610726578063a319de9e14610748578063a4155e5d146107bb578063ac6fff0b146107d3578063b3f00674146107e8578063bab2af1d146107fd578063c263b5d61461081e578063c877cf3714610833578063d8cf98ca1461084d578063dd9222d61461086e578063e744092e14610876578063eb7f0cbe14610897578063ee273ad7146108bb578063ee2aec6514610943578063efdcd97414610964578063f2fde38b14610985578063febb0cc2146109a6578063ffa1ad74146109bb575b61022d336001546000604051908082528060200260200182016040528015610227578160200160208202803883390190505b506109d0565b005b34801561023b57600080fd5b5061022d600160a060020a0360048035821691602480359091169160443591606435908101910135610c23565b34801561027457600080fd5b50610289600160a060020a0360043516610ce9565b604080519115158252519081900360200190f35b3480156102a957600080fd5b5061022d60043567ffffffffffffffff60243516610cfe565b3480156102ce57600080fd5b50610289610d3d565b3480156102e357600080fd5b50604080516020601f60e43560048181013592830184900484028501840190955281845261022d94803594600160a060020a0360248035821696604435831696606435909316956084359567ffffffffffffffff60a43581169660c435909116953695929461010494929392909201918190840183828082843750506040805187358901803560208181028481018201909552818452989b9a998901989297509082019550935083925085019084908082843750949750610d4d9650505050505050565b3480156103b357600080fd5b5061022d600160a060020a036004351660243567ffffffffffffffff60443516611007565b3480156103e457600080fd5b506103ed611048565b6040805167ffffffffffffffff9092168252519081900360200190f35b34801561041657600080fd5b50604080516084356004818101356020818102858101820190965281855261022d95600160a060020a038435811696602480359760443590931696606435963696919560a495909301929091829190850190849080828437509497506110589650505050505050565b34801561048b57600080fd5b5061022d600160a060020a0360043516611066565b3480156104ac57600080fd5b506104b561109d565b60408051600160a060020a039092168252519081900360200190f35b3480156104dd57600080fd5b50604080516020600460e43581810135838102808601850190965280855261022d9583359560248035600160a060020a0390811697604435821697606435909216966084359660a43567ffffffffffffffff9081169760c4359091169636966101049591940192918291850190849080828437509497506110b19650505050505050565b34801561056d57600080fd5b5061022d600160a060020a0360043581169060243516611306565b34801561059457600080fd5b506105a9600160a060020a0360043516611354565b60408051918252519081900360200190f35b3480156105c757600080fd5b506105d660ff60043516611366565b6040805160ff9092168252519081900360200190f35b3480156105f857600080fd5b506103ed61137b565b34801561060d57600080fd5b506103ed67ffffffffffffffff6004351661139f565b34801561062f57600080fd5b506102896113bb565b34801561064457600080fd5b5061022d6113c4565b34801561065957600080fd5b5061022d600160a060020a03600480358216916024803590911691604435916064359160843591820191013561142e565b34801561069657600080fd5b506103ed611516565b3480156106ab57600080fd5b5061022d600160a060020a036004351661152d565b3480156106cc57600080fd5b50610289600435611575565b3480156106e457600080fd5b506104b561158a565b3480156106f957600080fd5b50610289611599565b34801561070e57600080fd5b50610289600435600160a060020a03602435166115aa565b34801561073257600080fd5b5061028967ffffffffffffffff600435166115ca565b34801561075457600080fd5b50604080516064356004818101356020818102858101820190965281855261022d95600160a060020a038435811696602480359760443590931696369690956084959394909201929091829190850190849080828437509497506115df9650505050505050565b3480156107c757600080fd5b506105d66004356115f2565b3480156107df57600080fd5b506103ed611607565b3480156107f457600080fd5b506104b5611617565b34801561080957600080fd5b5061022d600160a060020a0360043516611626565b34801561082a57600080fd5b506105a961166d565b34801561083f57600080fd5b5061022d6004351515611673565b34801561085957600080fd5b5061022d600160a060020a03600435166116b6565b61022d6116ea565b34801561088257600080fd5b506104b5600160a060020a03600435166116ec565b61022d60048035600160a060020a0316906024803591604435918201910135611707565b3480156108c757600080fd5b506040805160c4356004818101356020818102858101820190965281855261022d95833595600160a060020a0360248035821697604435909216966064359667ffffffffffffffff60843581169760a435909116963696939560e495019291829190850190849080828437509497506117539650505050505050565b34801561094f57600080fd5b5061022d60ff600435811690602435166118ab565b34801561097057600080fd5b5061022d600160a060020a03600435166118f8565b34801561099157600080fd5b5061022d600160a060020a0360043516611917565b3480156109b257600080fd5b506103ed611933565b3480156109c757600080fd5b506103ed61194f565b600b5460009060a860020a900460ff161515610a36576040805160e560020a62461bcd02815260206004820152600e60248201527f73746f7070656420627269646765000000000000000000000000000000000000604482015290519081900360640190fd5b348310610a8d576040805160e560020a62461bcd02815260206004820152601360248201527f696e73756666696369656e7420616d6f756e7400000000000000000000000000604482015290519081900360640190fd5b610a9683611954565b90507f09a722fc78123f7b2ba788ae83b01221477757cfb4b83302421e7c5c1fd01da96000338682610ace348963ffffffff611a8a16565b600b60169054906101000a900467ffffffffffffffff16878960405180896002811115610af757fe5b60ff16815260200188600160a060020a0316600160a060020a0316815260200187600160a060020a0316600160a060020a0316815260200186600160a060020a0316600160a060020a031681526020018581526020018467ffffffffffffffff1667ffffffffffffffff1681526020018060200184815260200180602001838103835260008152602001602001838103825284818151815260200191508051906020019060200280838360005b83811015610bbc578181015183820152602001610ba4565b505050509050019a505050505050505050505060405180910390a15050600b8054600167ffffffffffffffff60b060020a80840482169290920116027fffff0000000000000000ffffffffffffffffffffffffffffffffffffffffffff9091161790555050565b604080517f23b872dd000000000000000000000000000000000000000000000000000000008152336004820152306024820152604481018590529051600160a060020a038716916323b872dd91606480830192600092919082900301818387803b158015610c9057600080fd5b505af1158015610ca4573d6000803e3d6000fd5b50505050610ce28533868686868080602002602001604051908101604052809392919081815260200183836020028082843750611aa1945050505050565b5050505050565b60046020526000908152604090205460ff1681565b3360009081526004602052604090205460ff161515610d1c57600080fd5b610d2581611e8b565b1515610d3057610d39565b610d3982611f57565b5050565b600b5460a860020a900460ff1681565b3360009081526004602052604090205460ff161515610d6b57600080fd5b610d7484611f8a565b1515610d7f57610ffc565b604080518a815260026020808301829052600160a060020a03808d1694840194909452838b166060840152928916608083015260a0820188905267ffffffffffffffff871660c083015261010060e0830181815285519184019190915284517f91f1d4b9864c9c091ade654e58899919d2a56e194153899c7e292cd3ae2ac47b948e948e938e938e938e938e938c9391929091610120840191808601910280838360005b83811015610e3b578181015183820152602001610e23565b50505050905001995050505050505050505060405180910390a1610e5e89612055565b610e688484612070565b600b5460ff1615610f725785600160a060020a03166350bb4e7f8887856040518463ffffffff1660e060020a0281526004018084600160a060020a0316600160a060020a0316815260200183815260200180602001828103825283818151815260200191508051906020019080838360005b83811015610ef2578181015183820152602001610eda565b50505050905090810190601f168015610f1f5780820380516001836020036101000a031916815260200191505b50945050505050602060405180830381600087803b158015610f4057600080fd5b505af1158015610f54573d6000803e3d6000fd5b505050506040513d6020811015610f6a57600080fd5b50610ffc9050565b604080517f42842e0e000000000000000000000000000000000000000000000000000000008152306004820152600160a060020a038981166024830152604482018890529151918816916342842e0e9160648082019260009290919082900301818387803b158015610fe357600080fd5b505af1158015610ff7573d6000803e3d6000fd5b505050505b505050505050505050565b3360009081526004602052604090205460ff16151561102557600080fd5b61102e81611e8b565b151561103957611043565b61104383836121dd565b505050565b600c5467ffffffffffffffff1681565b610ce2338685878686612232565b61106e611599565b151561107957600080fd5b600160a060020a03166000908152600460205260409020805460ff19166001179055565b600b546101009004600160a060020a031681565b3360009081526004602052604090205460ff1615156110cf57600080fd5b6110d883611f8a565b15156110e3576112fc565b6040805189815260016020808301829052600160a060020a03808c1694840194909452838a166060840152928816608083015260a0820187905267ffffffffffffffff861660c083015261010060e0830181815285519184019190915284517f91f1d4b9864c9c091ade654e58899919d2a56e194153899c7e292cd3ae2ac47b948d948d938d938d938d938d938c9391929091610120840191808601910280838360005b8381101561119f578181015183820152602001611187565b50505050905001995050505050505050505060405180910390a16111c288612055565b6111cc8383612070565b600b5460ff161561126c5784600160a060020a03166340c10f1987866040518363ffffffff1660e060020a0281526004018083600160a060020a0316600160a060020a0316815260200182815260200192505050602060405180830381600087803b15801561123a57600080fd5b505af115801561124e573d6000803e3d6000fd5b505050506040513d602081101561126457600080fd5b506112fc9050565b84600160a060020a031663a9059cbb87866040518363ffffffff1660e060020a0281526004018083600160a060020a0316600160a060020a0316815260200182815260200192505050602060405180830381600087803b1580156112cf57600080fd5b505af11580156112e3573d6000803e3d6000fd5b505050506040513d60208110156112f957600080fd5b50505b5050505050505050565b61130e611599565b151561131957600080fd5b600160a060020a039182166000908152600e60205260409020805473ffffffffffffffffffffffffffffffffffffffff191691909216179055565b60026020526000908152604090205481565b60086020526000908152604090205460ff1681565b600c54700100000000000000000000000000000000900467ffffffffffffffff1681565b600d6020526000908152604090205467ffffffffffffffff1681565b600b5460ff1681565b6113cc611599565b15156113d757600080fd5b600354604051600091600160a060020a0316907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0908390a36003805473ffffffffffffffffffffffffffffffffffffffff19169055565b600160a060020a0386166323b872dd333061144f888863ffffffff61255b16565b6040805160e060020a63ffffffff8716028152600160a060020a0394851660048201529290931660248301526044820152905160648083019260209291908290030181600087803b1580156114a357600080fd5b505af11580156114b7573d6000803e3d6000fd5b505050506040513d60208110156114cd57600080fd5b5050604080516020838102828101820190935283825261150e92899233928a928a928a9290918a918a91829185019084908082843750612232945050505050565b505050505050565b600b5460b060020a900467ffffffffffffffff1681565b611535611599565b151561154057600080fd5b600b8054600160a060020a039092166101000274ffffffffffffffffffffffffffffffffffffffff0019909216919091179055565b600a6020526000908152604090205460ff1681565b600354600160a060020a031690565b600354600160a060020a0316331490565b600560209081526000928352604080842090915290825290205460ff1681565b60076020526000908152604090205460ff1681565b6115ec3385848685611aa1565b50505050565b60066020526000908152604090205460ff1681565b60095467ffffffffffffffff1681565b600054600160a060020a031681565b61162e611599565b151561163957600080fd5b600160a060020a03166000908152600e60205260409020805473ffffffffffffffffffffffffffffffffffffffff19169055565b60015481565b61167b611599565b151561168657600080fd5b600b805491151560a860020a0275ff00000000000000000000000000000000000000000019909216919091179055565b6116be611599565b15156116c957600080fd5b600160a060020a03166000908152600460205260409020805460ff19169055565b565b600e60205260009081526040902054600160a060020a031681565b6000611719348563ffffffff611a8a16565b9050610ce28582858580806020026020016040519081016040528093929190818152602001838360200280828437506109d0945050505050565b3360009081526004602052604090205460ff16151561177157600080fd5b61177a83611f8a565b1515611785576118a2565b6040805188815260006020808301829052600160a060020a03808b169484019490945292881660608301526080820181905260a0820187905267ffffffffffffffff861660c083015261010060e0830181815285519184019190915284517f91f1d4b9864c9c091ade654e58899919d2a56e194153899c7e292cd3ae2ac47b948c948c938c9386938d938d938c939192909161012084019180860191028083838b5b8381101561183f578181015183820152602001611827565b50505050905001995050505050505050505060405180910390a161186287612055565b61186c8383612070565b604051600160a060020a0386169085156108fc029086906000818181858888f193505050501580156112fc573d6000803e3d6000fd5b50505050505050565b6118b3611599565b15156118be57600080fd5b80600860008460028111156118cf57fe5b60ff90811682526020820192909252604001600020805460ff1916929091169190911790555050565b611900611599565b151561190b57600080fd5b61191481612574565b50565b61191f611599565b151561192a57600080fd5b611914816125c9565b600c5468010000000000000000900467ffffffffffffffff1681565b600181565b60015460008054909190600160a060020a0316158015906119755750600081115b15611a5157808310156119d2576040805160e560020a62461bcd02815260206004820152601560248201527f696e73756666696369656e74206665654c696d69740000000000000000000000604482015290519081900360640190fd5b60008054604051600160a060020a039091169183156108fc02918491818181858888f19350505050158015611a0b573d6000803e3d6000fd5b50336108fc611a20858463ffffffff611a8a16565b6040518115909202916000818181858888f19350505050158015611a48573d6000803e3d6000fd5b50809150611a84565b604051339084156108fc029085906000818181858888f19350505050158015611a7e573d6000803e3d6000fd5b50600091505b50919050565b60008083831115611a9a57600080fd5b5050900390565b600b5460609060a860020a900460ff161515611b07576040805160e560020a62461bcd02815260206004820152600e60248201527f73746f7070656420627269646765000000000000000000000000000000000000604482015290519081900360640190fd5b600160a060020a038681166000908152600e6020526040902054161515611b78576040805160e560020a62461bcd02815260206004820152600d60248201527f696e76616c696420746f6b656e00000000000000000000000000000000000000604482015290519081900360640190fd5b85600160a060020a031663c87b56dd846040518263ffffffff1660e060020a02815260040180828152602001915050600060405180830381600087803b158015611bc157600080fd5b505af1158015611bd5573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f191682016040526020811015611bfe57600080fd5b810190808051640100000000811115611c1657600080fd5b82016020810184811115611c2957600080fd5b8151640100000000811182820187101715611c4357600080fd5b5050600b5490945060ff16159250611cba9150505785600160a060020a03166342966c68846040518263ffffffff1660e060020a02815260040180828152602001915050600060405180830381600087803b158015611ca157600080fd5b505af1158015611cb5573d6000803e3d6000fd5b505050505b7f09a722fc78123f7b2ba788ae83b01221477757cfb4b83302421e7c5c1fd01da9600286868987600b60169054906101000a900467ffffffffffffffff168760008a604051808a6002811115611d0c57fe5b60ff16815260200189600160a060020a0316600160a060020a0316815260200188600160a060020a0316600160a060020a0316815260200187600160a060020a0316600160a060020a031681526020018681526020018567ffffffffffffffff1667ffffffffffffffff1681526020018060200184815260200180602001838103835286818151815260200191508051906020019080838360005b83811015611dbf578181015183820152602001611da7565b50505050905090810190601f168015611dec5780820380516001836020036101000a031916815260200191505b508381038252845181528451602091820191808701910280838360005b83811015611e21578181015183820152602001611e09565b505050509050019b50505050505050505050505060405180910390a15050600b8054600167ffffffffffffffff60b060020a80840482169290920116027fffff0000000000000000ffffffffffffffffffffffffffffffffffffffffffff90911617905550505050565b60095460009067ffffffffffffffff838116911614611ef4576040805160e560020a62461bcd02815260206004820152600e60248201527f6e6f6e6365206d69736d61746368000000000000000000000000000000000000604482015290519081900360640190fd5b611f1c6001600036604051808383808284378201915050925050506040518091039020612647565b15611f4e57506009805467ffffffffffffffff198116600167ffffffffffffffff928316810190921617909155611f52565b5060005b919050565b600181905560405181907fa7a33d0996347e1aa55ca2206015b61b9534bdd881d59d59aa680e25eefac36590600090a250565b67ffffffffffffffff811660009081526007602052604081205460ff1615611ffc576040805160e560020a62461bcd02815260206004820152600b60248201527f636c6f73656420766f7465000000000000000000000000000000000000000000604482015290519081900360640190fd5b61202360008036604051808383808284378201915050925050506040518091039020612647565b15611f4e575067ffffffffffffffff81166000908152600760205260409020805460ff19166001908117909155611f52565b6000908152600a60205260409020805460ff19166001179055565b67ffffffffffffffff8281166000818152600d60205260408120805467ffffffffffffffff1916858516179055600c5490927001000000000000000000000000000000009091041610156120ff57600c805477ffffffffffffffff00000000000000000000000000000000191670010000000000000000000000000000000067ffffffffffffffff8616021790555b50600c5468010000000000000000900467ffffffffffffffff165b600c5467ffffffffffffffff700100000000000000000000000000000000909104811690821611801590612169575067ffffffffffffffff8082166000908152600d6020526040812054909116115b156121765760010161211a565b600c805467ffffffffffffffff80841668010000000000000000026fffffffffffffffff0000000000000000199092169190911780835560001990930181166000908152600d602052604090205467ffffffffffffffff1990931692169190911790555050565b600160a060020a0382166000818152600260209081526040918290208490558151928352905183927fdb5ad2e76ae20cfa4e7adbc7305d7538442164d85ead9937c98620a1aa4c255b92908290030190a25050565b600b5460009060a860020a900460ff161515612298576040805160e560020a62461bcd02815260206004820152600e60248201527f73746f7070656420627269646765000000000000000000000000000000000000604482015290519081900360640190fd5b600084116122f0576040805160e560020a62461bcd02815260206004820152600e60248201527f7a65726f206d73672e76616c7565000000000000000000000000000000000000604482015290519081900360640190fd5b600160a060020a038781166000908152600e6020526040902054161515612361576040805160e560020a62461bcd02815260206004820152600d60248201527f696e76616c696420746f6b656e00000000000000000000000000000000000000604482015290519081900360640190fd5b61236c868885612727565b600b5490915060ff16156123dc5786600160a060020a03166342966c68856040518263ffffffff1660e060020a02815260040180828152602001915050600060405180830381600087803b1580156123c357600080fd5b505af11580156123d7573d6000803e3d6000fd5b505050505b7f09a722fc78123f7b2ba788ae83b01221477757cfb4b83302421e7c5c1fd01da9600187878a88600b60169054906101000a900467ffffffffffffffff1687896040518089600281111561242c57fe5b60ff16815260200188600160a060020a0316600160a060020a0316815260200187600160a060020a0316600160a060020a0316815260200186600160a060020a0316600160a060020a031681526020018581526020018467ffffffffffffffff1667ffffffffffffffff1681526020018060200184815260200180602001838103835260008152602001602001838103825284818151815260200191508051906020019060200280838360005b838110156124f15781810151838201526020016124d9565b505050509050019a505050505050505050505060405180910390a15050600b8054600167ffffffffffffffff60b060020a80840482169290920116027fffff0000000000000000ffffffffffffffffffffffffffffffffffffffffffff9091161790555050505050565b60008282018381101561256d57600080fd5b9392505050565b6000805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a038316908117825560405190917f647672599d3468abcfa241a13c9e3d34383caadb5cc80fb67c3cdfcd5f78605991a250565b600160a060020a03811615156125de57600080fd5b600354604051600160a060020a038084169216907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a36003805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a0392909216919091179055565b600081815260056020908152604080832033845290915281205460ff1615156126d35760008281526005602090815260408083203384528252808320805460ff1916600190811790915585845260069092529091205460ff90811691820116116126b057600080fd5b6000828152600660205260409020805460ff8082166001011660ff199091161790555b600860008460028111156126e357fe5b60ff908116825260208083019390935260409182016000908120548682526006909452919091205491811691161061271d57506001612721565b5060005b92915050565b600160a060020a0380831660009081526002602052604081205481549192909116158015906127565750600081115b156128f757808310156127b3576040805160e560020a62461bcd02815260206004820152601560248201527f696e73756666696369656e74206665654c696d69740000000000000000000000604482015290519081900360640190fd5b60008054604080517fa9059cbb000000000000000000000000000000000000000000000000000000008152600160a060020a0392831660048201526024810185905290519187169263a9059cbb926044808401936020939083900390910190829087803b15801561282357600080fd5b505af1158015612837573d6000803e3d6000fd5b505050506040513d602081101561284d57600080fd5b5050600160a060020a03841663a9059cbb8661286f868563ffffffff611a8a16565b6040518363ffffffff1660e060020a0281526004018083600160a060020a0316600160a060020a0316815260200182815260200192505050602060405180830381600087803b1580156128c157600080fd5b505af11580156128d5573d6000803e3d6000fd5b505050506040513d60208110156128eb57600080fd5b5090915081905061298b565b83600160a060020a031663a9059cbb86856040518363ffffffff1660e060020a0281526004018083600160a060020a0316600160a060020a0316815260200182815260200192505050602060405180830381600087803b15801561295a57600080fd5b505af115801561296e573d6000803e3d6000fd5b505050506040513d602081101561298457600080fd5b5060009250505b5093925050505600a165627a7a72305820b26ab42c0688c114123227adc844724b08ff5886ef6f44271ff99c60751076850029`

// DeployBridge deploys a new Klaytn contract, binding an instance of Bridge to it.
func DeployBridge(auth *bind.TransactOpts, backend bind.ContractBackend, _modeMintBurn bool) (common.Address, *types.Transaction, *Bridge, error) {
	parsed, err := abi.JSON(strings.NewReader(BridgeABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(BridgeBin), backend, _modeMintBurn)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Bridge{BridgeCaller: BridgeCaller{contract: contract}, BridgeTransactor: BridgeTransactor{contract: contract}, BridgeFilterer: BridgeFilterer{contract: contract}}, nil
}

// Bridge is an auto generated Go binding around a Klaytn contract.
type Bridge struct {
	BridgeCaller     // Read-only binding to the contract
	BridgeTransactor // Write-only binding to the contract
	BridgeFilterer   // Log filterer for contract events
}

// BridgeCaller is an auto generated read-only Go binding around a Klaytn contract.
type BridgeCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BridgeTransactor is an auto generated write-only Go binding around a Klaytn contract.
type BridgeTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BridgeFilterer is an auto generated log filtering Go binding around a Klaytn contract events.
type BridgeFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BridgeSession is an auto generated Go binding around a Klaytn contract,
// with pre-set call and transact options.
type BridgeSession struct {
	Contract     *Bridge           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// BridgeCallerSession is an auto generated read-only Go binding around a Klaytn contract,
// with pre-set call options.
type BridgeCallerSession struct {
	Contract *BridgeCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// BridgeTransactorSession is an auto generated write-only Go binding around a Klaytn contract,
// with pre-set transact options.
type BridgeTransactorSession struct {
	Contract     *BridgeTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// BridgeRaw is an auto generated low-level Go binding around a Klaytn contract.
type BridgeRaw struct {
	Contract *Bridge // Generic contract binding to access the raw methods on
}

// BridgeCallerRaw is an auto generated low-level read-only Go binding around a Klaytn contract.
type BridgeCallerRaw struct {
	Contract *BridgeCaller // Generic read-only contract binding to access the raw methods on
}

// BridgeTransactorRaw is an auto generated low-level write-only Go binding around a Klaytn contract.
type BridgeTransactorRaw struct {
	Contract *BridgeTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBridge creates a new instance of Bridge, bound to a specific deployed contract.
func NewBridge(address common.Address, backend bind.ContractBackend) (*Bridge, error) {
	contract, err := bindBridge(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Bridge{BridgeCaller: BridgeCaller{contract: contract}, BridgeTransactor: BridgeTransactor{contract: contract}, BridgeFilterer: BridgeFilterer{contract: contract}}, nil
}

// NewBridgeCaller creates a new read-only instance of Bridge, bound to a specific deployed contract.
func NewBridgeCaller(address common.Address, caller bind.ContractCaller) (*BridgeCaller, error) {
	contract, err := bindBridge(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BridgeCaller{contract: contract}, nil
}

// NewBridgeTransactor creates a new write-only instance of Bridge, bound to a specific deployed contract.
func NewBridgeTransactor(address common.Address, transactor bind.ContractTransactor) (*BridgeTransactor, error) {
	contract, err := bindBridge(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BridgeTransactor{contract: contract}, nil
}

// NewBridgeFilterer creates a new log filterer instance of Bridge, bound to a specific deployed contract.
func NewBridgeFilterer(address common.Address, filterer bind.ContractFilterer) (*BridgeFilterer, error) {
	contract, err := bindBridge(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BridgeFilterer{contract: contract}, nil
}

// bindBridge binds a generic wrapper to an already deployed contract.
func bindBridge(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(BridgeABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Bridge *BridgeRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Bridge.Contract.BridgeCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Bridge *BridgeRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Bridge.Contract.BridgeTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Bridge *BridgeRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Bridge.Contract.BridgeTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Bridge *BridgeCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Bridge.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Bridge *BridgeTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Bridge.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Bridge *BridgeTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Bridge.Contract.contract.Transact(opts, method, params...)
}

// VERSION is a free data retrieval call binding the contract method 0xffa1ad74.
//
// Solidity: function VERSION() constant returns(uint64)
func (_Bridge *BridgeCaller) VERSION(opts *bind.CallOpts) (uint64, error) {
	var (
		ret0 = new(uint64)
	)
	out := ret0
	err := _Bridge.contract.Call(opts, out, "VERSION")
	return *ret0, err
}

// VERSION is a free data retrieval call binding the contract method 0xffa1ad74.
//
// Solidity: function VERSION() constant returns(uint64)
func (_Bridge *BridgeSession) VERSION() (uint64, error) {
	return _Bridge.Contract.VERSION(&_Bridge.CallOpts)
}

// VERSION is a free data retrieval call binding the contract method 0xffa1ad74.
//
// Solidity: function VERSION() constant returns(uint64)
func (_Bridge *BridgeCallerSession) VERSION() (uint64, error) {
	return _Bridge.Contract.VERSION(&_Bridge.CallOpts)
}

// AllowedTokens is a free data retrieval call binding the contract method 0xe744092e.
//
// Solidity: function allowedTokens( address) constant returns(address)
func (_Bridge *BridgeCaller) AllowedTokens(opts *bind.CallOpts, arg0 common.Address) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Bridge.contract.Call(opts, out, "allowedTokens", arg0)
	return *ret0, err
}

// AllowedTokens is a free data retrieval call binding the contract method 0xe744092e.
//
// Solidity: function allowedTokens( address) constant returns(address)
func (_Bridge *BridgeSession) AllowedTokens(arg0 common.Address) (common.Address, error) {
	return _Bridge.Contract.AllowedTokens(&_Bridge.CallOpts, arg0)
}

// AllowedTokens is a free data retrieval call binding the contract method 0xe744092e.
//
// Solidity: function allowedTokens( address) constant returns(address)
func (_Bridge *BridgeCallerSession) AllowedTokens(arg0 common.Address) (common.Address, error) {
	return _Bridge.Contract.AllowedTokens(&_Bridge.CallOpts, arg0)
}

// ClosedValueTransferVotes is a free data retrieval call binding the contract method 0x9832c1d7.
//
// Solidity: function closedValueTransferVotes( uint64) constant returns(bool)
func (_Bridge *BridgeCaller) ClosedValueTransferVotes(opts *bind.CallOpts, arg0 uint64) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Bridge.contract.Call(opts, out, "closedValueTransferVotes", arg0)
	return *ret0, err
}

// ClosedValueTransferVotes is a free data retrieval call binding the contract method 0x9832c1d7.
//
// Solidity: function closedValueTransferVotes( uint64) constant returns(bool)
func (_Bridge *BridgeSession) ClosedValueTransferVotes(arg0 uint64) (bool, error) {
	return _Bridge.Contract.ClosedValueTransferVotes(&_Bridge.CallOpts, arg0)
}

// ClosedValueTransferVotes is a free data retrieval call binding the contract method 0x9832c1d7.
//
// Solidity: function closedValueTransferVotes( uint64) constant returns(bool)
func (_Bridge *BridgeCallerSession) ClosedValueTransferVotes(arg0 uint64) (bool, error) {
	return _Bridge.Contract.ClosedValueTransferVotes(&_Bridge.CallOpts, arg0)
}

// ConfigurationNonce is a free data retrieval call binding the contract method 0xac6fff0b.
//
// Solidity: function configurationNonce() constant returns(uint64)
func (_Bridge *BridgeCaller) ConfigurationNonce(opts *bind.CallOpts) (uint64, error) {
	var (
		ret0 = new(uint64)
	)
	out := ret0
	err := _Bridge.contract.Call(opts, out, "configurationNonce")
	return *ret0, err
}

// ConfigurationNonce is a free data retrieval call binding the contract method 0xac6fff0b.
//
// Solidity: function configurationNonce() constant returns(uint64)
func (_Bridge *BridgeSession) ConfigurationNonce() (uint64, error) {
	return _Bridge.Contract.ConfigurationNonce(&_Bridge.CallOpts)
}

// ConfigurationNonce is a free data retrieval call binding the contract method 0xac6fff0b.
//
// Solidity: function configurationNonce() constant returns(uint64)
func (_Bridge *BridgeCallerSession) ConfigurationNonce() (uint64, error) {
	return _Bridge.Contract.ConfigurationNonce(&_Bridge.CallOpts)
}

// CounterpartBridge is a free data retrieval call binding the contract method 0x3a348533.
//
// Solidity: function counterpartBridge() constant returns(address)
func (_Bridge *BridgeCaller) CounterpartBridge(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Bridge.contract.Call(opts, out, "counterpartBridge")
	return *ret0, err
}

// CounterpartBridge is a free data retrieval call binding the contract method 0x3a348533.
//
// Solidity: function counterpartBridge() constant returns(address)
func (_Bridge *BridgeSession) CounterpartBridge() (common.Address, error) {
	return _Bridge.Contract.CounterpartBridge(&_Bridge.CallOpts)
}

// CounterpartBridge is a free data retrieval call binding the contract method 0x3a348533.
//
// Solidity: function counterpartBridge() constant returns(address)
func (_Bridge *BridgeCallerSession) CounterpartBridge() (common.Address, error) {
	return _Bridge.Contract.CounterpartBridge(&_Bridge.CallOpts)
}

// FeeOfERC20 is a free data retrieval call binding the contract method 0x488af871.
//
// Solidity: function feeOfERC20( address) constant returns(uint256)
func (_Bridge *BridgeCaller) FeeOfERC20(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Bridge.contract.Call(opts, out, "feeOfERC20", arg0)
	return *ret0, err
}

// FeeOfERC20 is a free data retrieval call binding the contract method 0x488af871.
//
// Solidity: function feeOfERC20( address) constant returns(uint256)
func (_Bridge *BridgeSession) FeeOfERC20(arg0 common.Address) (*big.Int, error) {
	return _Bridge.Contract.FeeOfERC20(&_Bridge.CallOpts, arg0)
}

// FeeOfERC20 is a free data retrieval call binding the contract method 0x488af871.
//
// Solidity: function feeOfERC20( address) constant returns(uint256)
func (_Bridge *BridgeCallerSession) FeeOfERC20(arg0 common.Address) (*big.Int, error) {
	return _Bridge.Contract.FeeOfERC20(&_Bridge.CallOpts, arg0)
}

// FeeOfKLAY is a free data retrieval call binding the contract method 0xc263b5d6.
//
// Solidity: function feeOfKLAY() constant returns(uint256)
func (_Bridge *BridgeCaller) FeeOfKLAY(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Bridge.contract.Call(opts, out, "feeOfKLAY")
	return *ret0, err
}

// FeeOfKLAY is a free data retrieval call binding the contract method 0xc263b5d6.
//
// Solidity: function feeOfKLAY() constant returns(uint256)
func (_Bridge *BridgeSession) FeeOfKLAY() (*big.Int, error) {
	return _Bridge.Contract.FeeOfKLAY(&_Bridge.CallOpts)
}

// FeeOfKLAY is a free data retrieval call binding the contract method 0xc263b5d6.
//
// Solidity: function feeOfKLAY() constant returns(uint256)
func (_Bridge *BridgeCallerSession) FeeOfKLAY() (*big.Int, error) {
	return _Bridge.Contract.FeeOfKLAY(&_Bridge.CallOpts)
}

// FeeReceiver is a free data retrieval call binding the contract method 0xb3f00674.
//
// Solidity: function feeReceiver() constant returns(address)
func (_Bridge *BridgeCaller) FeeReceiver(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Bridge.contract.Call(opts, out, "feeReceiver")
	return *ret0, err
}

// FeeReceiver is a free data retrieval call binding the contract method 0xb3f00674.
//
// Solidity: function feeReceiver() constant returns(address)
func (_Bridge *BridgeSession) FeeReceiver() (common.Address, error) {
	return _Bridge.Contract.FeeReceiver(&_Bridge.CallOpts)
}

// FeeReceiver is a free data retrieval call binding the contract method 0xb3f00674.
//
// Solidity: function feeReceiver() constant returns(address)
func (_Bridge *BridgeCallerSession) FeeReceiver() (common.Address, error) {
	return _Bridge.Contract.FeeReceiver(&_Bridge.CallOpts)
}

// HandledHistory is a free data retrieval call binding the contract method 0x67bda76b.
//
// Solidity: function handledHistory( uint64) constant returns(uint64)
func (_Bridge *BridgeCaller) HandledHistory(opts *bind.CallOpts, arg0 uint64) (uint64, error) {
	var (
		ret0 = new(uint64)
	)
	out := ret0
	err := _Bridge.contract.Call(opts, out, "handledHistory", arg0)
	return *ret0, err
}

// HandledHistory is a free data retrieval call binding the contract method 0x67bda76b.
//
// Solidity: function handledHistory( uint64) constant returns(uint64)
func (_Bridge *BridgeSession) HandledHistory(arg0 uint64) (uint64, error) {
	return _Bridge.Contract.HandledHistory(&_Bridge.CallOpts, arg0)
}

// HandledHistory is a free data retrieval call binding the contract method 0x67bda76b.
//
// Solidity: function handledHistory( uint64) constant returns(uint64)
func (_Bridge *BridgeCallerSession) HandledHistory(arg0 uint64) (uint64, error) {
	return _Bridge.Contract.HandledHistory(&_Bridge.CallOpts, arg0)
}

// HandledRequestTx is a free data retrieval call binding the contract method 0x8a75eee2.
//
// Solidity: function handledRequestTx( bytes32) constant returns(bool)
func (_Bridge *BridgeCaller) HandledRequestTx(opts *bind.CallOpts, arg0 [32]byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Bridge.contract.Call(opts, out, "handledRequestTx", arg0)
	return *ret0, err
}

// HandledRequestTx is a free data retrieval call binding the contract method 0x8a75eee2.
//
// Solidity: function handledRequestTx( bytes32) constant returns(bool)
func (_Bridge *BridgeSession) HandledRequestTx(arg0 [32]byte) (bool, error) {
	return _Bridge.Contract.HandledRequestTx(&_Bridge.CallOpts, arg0)
}

// HandledRequestTx is a free data retrieval call binding the contract method 0x8a75eee2.
//
// Solidity: function handledRequestTx( bytes32) constant returns(bool)
func (_Bridge *BridgeCallerSession) HandledRequestTx(arg0 [32]byte) (bool, error) {
	return _Bridge.Contract.HandledRequestTx(&_Bridge.CallOpts, arg0)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() constant returns(bool)
func (_Bridge *BridgeCaller) IsOwner(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Bridge.contract.Call(opts, out, "isOwner")
	return *ret0, err
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() constant returns(bool)
func (_Bridge *BridgeSession) IsOwner() (bool, error) {
	return _Bridge.Contract.IsOwner(&_Bridge.CallOpts)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() constant returns(bool)
func (_Bridge *BridgeCallerSession) IsOwner() (bool, error) {
	return _Bridge.Contract.IsOwner(&_Bridge.CallOpts)
}

// IsRunning is a free data retrieval call binding the contract method 0x2014e5d1.
//
// Solidity: function isRunning() constant returns(bool)
func (_Bridge *BridgeCaller) IsRunning(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Bridge.contract.Call(opts, out, "isRunning")
	return *ret0, err
}

// IsRunning is a free data retrieval call binding the contract method 0x2014e5d1.
//
// Solidity: function isRunning() constant returns(bool)
func (_Bridge *BridgeSession) IsRunning() (bool, error) {
	return _Bridge.Contract.IsRunning(&_Bridge.CallOpts)
}

// IsRunning is a free data retrieval call binding the contract method 0x2014e5d1.
//
// Solidity: function isRunning() constant returns(bool)
func (_Bridge *BridgeCallerSession) IsRunning() (bool, error) {
	return _Bridge.Contract.IsRunning(&_Bridge.CallOpts)
}

// MaxHandledRequestedNonce is a free data retrieval call binding the contract method 0x60f17a9b.
//
// Solidity: function maxHandledRequestedNonce() constant returns(uint64)
func (_Bridge *BridgeCaller) MaxHandledRequestedNonce(opts *bind.CallOpts) (uint64, error) {
	var (
		ret0 = new(uint64)
	)
	out := ret0
	err := _Bridge.contract.Call(opts, out, "maxHandledRequestedNonce")
	return *ret0, err
}

// MaxHandledRequestedNonce is a free data retrieval call binding the contract method 0x60f17a9b.
//
// Solidity: function maxHandledRequestedNonce() constant returns(uint64)
func (_Bridge *BridgeSession) MaxHandledRequestedNonce() (uint64, error) {
	return _Bridge.Contract.MaxHandledRequestedNonce(&_Bridge.CallOpts)
}

// MaxHandledRequestedNonce is a free data retrieval call binding the contract method 0x60f17a9b.
//
// Solidity: function maxHandledRequestedNonce() constant returns(uint64)
func (_Bridge *BridgeCallerSession) MaxHandledRequestedNonce() (uint64, error) {
	return _Bridge.Contract.MaxHandledRequestedNonce(&_Bridge.CallOpts)
}

// ModeMintBurn is a free data retrieval call binding the contract method 0x6e176ec2.
//
// Solidity: function modeMintBurn() constant returns(bool)
func (_Bridge *BridgeCaller) ModeMintBurn(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Bridge.contract.Call(opts, out, "modeMintBurn")
	return *ret0, err
}

// ModeMintBurn is a free data retrieval call binding the contract method 0x6e176ec2.
//
// Solidity: function modeMintBurn() constant returns(bool)
func (_Bridge *BridgeSession) ModeMintBurn() (bool, error) {
	return _Bridge.Contract.ModeMintBurn(&_Bridge.CallOpts)
}

// ModeMintBurn is a free data retrieval call binding the contract method 0x6e176ec2.
//
// Solidity: function modeMintBurn() constant returns(bool)
func (_Bridge *BridgeCallerSession) ModeMintBurn() (bool, error) {
	return _Bridge.Contract.ModeMintBurn(&_Bridge.CallOpts)
}

// OperatorThresholds is a free data retrieval call binding the contract method 0x5526f76b.
//
// Solidity: function operatorThresholds( uint8) constant returns(uint8)
func (_Bridge *BridgeCaller) OperatorThresholds(opts *bind.CallOpts, arg0 uint8) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _Bridge.contract.Call(opts, out, "operatorThresholds", arg0)
	return *ret0, err
}

// OperatorThresholds is a free data retrieval call binding the contract method 0x5526f76b.
//
// Solidity: function operatorThresholds( uint8) constant returns(uint8)
func (_Bridge *BridgeSession) OperatorThresholds(arg0 uint8) (uint8, error) {
	return _Bridge.Contract.OperatorThresholds(&_Bridge.CallOpts, arg0)
}

// OperatorThresholds is a free data retrieval call binding the contract method 0x5526f76b.
//
// Solidity: function operatorThresholds( uint8) constant returns(uint8)
func (_Bridge *BridgeCallerSession) OperatorThresholds(arg0 uint8) (uint8, error) {
	return _Bridge.Contract.OperatorThresholds(&_Bridge.CallOpts, arg0)
}

// Operators is a free data retrieval call binding the contract method 0x13e7c9d8.
//
// Solidity: function operators( address) constant returns(bool)
func (_Bridge *BridgeCaller) Operators(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Bridge.contract.Call(opts, out, "operators", arg0)
	return *ret0, err
}

// Operators is a free data retrieval call binding the contract method 0x13e7c9d8.
//
// Solidity: function operators( address) constant returns(bool)
func (_Bridge *BridgeSession) Operators(arg0 common.Address) (bool, error) {
	return _Bridge.Contract.Operators(&_Bridge.CallOpts, arg0)
}

// Operators is a free data retrieval call binding the contract method 0x13e7c9d8.
//
// Solidity: function operators( address) constant returns(bool)
func (_Bridge *BridgeCallerSession) Operators(arg0 common.Address) (bool, error) {
	return _Bridge.Contract.Operators(&_Bridge.CallOpts, arg0)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Bridge *BridgeCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Bridge.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Bridge *BridgeSession) Owner() (common.Address, error) {
	return _Bridge.Contract.Owner(&_Bridge.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Bridge *BridgeCallerSession) Owner() (common.Address, error) {
	return _Bridge.Contract.Owner(&_Bridge.CallOpts)
}

// RequestNonce is a free data retrieval call binding the contract method 0x7c1a0302.
//
// Solidity: function requestNonce() constant returns(uint64)
func (_Bridge *BridgeCaller) RequestNonce(opts *bind.CallOpts) (uint64, error) {
	var (
		ret0 = new(uint64)
	)
	out := ret0
	err := _Bridge.contract.Call(opts, out, "requestNonce")
	return *ret0, err
}

// RequestNonce is a free data retrieval call binding the contract method 0x7c1a0302.
//
// Solidity: function requestNonce() constant returns(uint64)
func (_Bridge *BridgeSession) RequestNonce() (uint64, error) {
	return _Bridge.Contract.RequestNonce(&_Bridge.CallOpts)
}

// RequestNonce is a free data retrieval call binding the contract method 0x7c1a0302.
//
// Solidity: function requestNonce() constant returns(uint64)
func (_Bridge *BridgeCallerSession) RequestNonce() (uint64, error) {
	return _Bridge.Contract.RequestNonce(&_Bridge.CallOpts)
}

// SequentialHandleNonce is a free data retrieval call binding the contract method 0xfebb0cc2.
//
// Solidity: function sequentialHandleNonce() constant returns(uint64)
func (_Bridge *BridgeCaller) SequentialHandleNonce(opts *bind.CallOpts) (uint64, error) {
	var (
		ret0 = new(uint64)
	)
	out := ret0
	err := _Bridge.contract.Call(opts, out, "sequentialHandleNonce")
	return *ret0, err
}

// SequentialHandleNonce is a free data retrieval call binding the contract method 0xfebb0cc2.
//
// Solidity: function sequentialHandleNonce() constant returns(uint64)
func (_Bridge *BridgeSession) SequentialHandleNonce() (uint64, error) {
	return _Bridge.Contract.SequentialHandleNonce(&_Bridge.CallOpts)
}

// SequentialHandleNonce is a free data retrieval call binding the contract method 0xfebb0cc2.
//
// Solidity: function sequentialHandleNonce() constant returns(uint64)
func (_Bridge *BridgeCallerSession) SequentialHandleNonce() (uint64, error) {
	return _Bridge.Contract.SequentialHandleNonce(&_Bridge.CallOpts)
}

// SequentialHandledRequestBlockNumber is a free data retrieval call binding the contract method 0x30f436d0.
//
// Solidity: function sequentialHandledRequestBlockNumber() constant returns(uint64)
func (_Bridge *BridgeCaller) SequentialHandledRequestBlockNumber(opts *bind.CallOpts) (uint64, error) {
	var (
		ret0 = new(uint64)
	)
	out := ret0
	err := _Bridge.contract.Call(opts, out, "sequentialHandledRequestBlockNumber")
	return *ret0, err
}

// SequentialHandledRequestBlockNumber is a free data retrieval call binding the contract method 0x30f436d0.
//
// Solidity: function sequentialHandledRequestBlockNumber() constant returns(uint64)
func (_Bridge *BridgeSession) SequentialHandledRequestBlockNumber() (uint64, error) {
	return _Bridge.Contract.SequentialHandledRequestBlockNumber(&_Bridge.CallOpts)
}

// SequentialHandledRequestBlockNumber is a free data retrieval call binding the contract method 0x30f436d0.
//
// Solidity: function sequentialHandledRequestBlockNumber() constant returns(uint64)
func (_Bridge *BridgeCallerSession) SequentialHandledRequestBlockNumber() (uint64, error) {
	return _Bridge.Contract.SequentialHandledRequestBlockNumber(&_Bridge.CallOpts)
}

// Votes is a free data retrieval call binding the contract method 0x9386775a.
//
// Solidity: function votes( bytes32,  address) constant returns(bool)
func (_Bridge *BridgeCaller) Votes(opts *bind.CallOpts, arg0 [32]byte, arg1 common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Bridge.contract.Call(opts, out, "votes", arg0, arg1)
	return *ret0, err
}

// Votes is a free data retrieval call binding the contract method 0x9386775a.
//
// Solidity: function votes( bytes32,  address) constant returns(bool)
func (_Bridge *BridgeSession) Votes(arg0 [32]byte, arg1 common.Address) (bool, error) {
	return _Bridge.Contract.Votes(&_Bridge.CallOpts, arg0, arg1)
}

// Votes is a free data retrieval call binding the contract method 0x9386775a.
//
// Solidity: function votes( bytes32,  address) constant returns(bool)
func (_Bridge *BridgeCallerSession) Votes(arg0 [32]byte, arg1 common.Address) (bool, error) {
	return _Bridge.Contract.Votes(&_Bridge.CallOpts, arg0, arg1)
}

// VotesCounts is a free data retrieval call binding the contract method 0xa4155e5d.
//
// Solidity: function votesCounts( bytes32) constant returns(uint8)
func (_Bridge *BridgeCaller) VotesCounts(opts *bind.CallOpts, arg0 [32]byte) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _Bridge.contract.Call(opts, out, "votesCounts", arg0)
	return *ret0, err
}

// VotesCounts is a free data retrieval call binding the contract method 0xa4155e5d.
//
// Solidity: function votesCounts( bytes32) constant returns(uint8)
func (_Bridge *BridgeSession) VotesCounts(arg0 [32]byte) (uint8, error) {
	return _Bridge.Contract.VotesCounts(&_Bridge.CallOpts, arg0)
}

// VotesCounts is a free data retrieval call binding the contract method 0xa4155e5d.
//
// Solidity: function votesCounts( bytes32) constant returns(uint8)
func (_Bridge *BridgeCallerSession) VotesCounts(arg0 [32]byte) (uint8, error) {
	return _Bridge.Contract.VotesCounts(&_Bridge.CallOpts, arg0)
}

// ChargeWithoutEvent is a paid mutator transaction binding the contract method 0xdd9222d6.
//
// Solidity: function chargeWithoutEvent() returns()
func (_Bridge *BridgeTransactor) ChargeWithoutEvent(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "chargeWithoutEvent")
}

// ChargeWithoutEvent is a paid mutator transaction binding the contract method 0xdd9222d6.
//
// Solidity: function chargeWithoutEvent() returns()
func (_Bridge *BridgeSession) ChargeWithoutEvent() (*types.Transaction, error) {
	return _Bridge.Contract.ChargeWithoutEvent(&_Bridge.TransactOpts)
}

// ChargeWithoutEvent is a paid mutator transaction binding the contract method 0xdd9222d6.
//
// Solidity: function chargeWithoutEvent() returns()
func (_Bridge *BridgeTransactorSession) ChargeWithoutEvent() (*types.Transaction, error) {
	return _Bridge.Contract.ChargeWithoutEvent(&_Bridge.TransactOpts)
}

// DeregisterOperator is a paid mutator transaction binding the contract method 0xd8cf98ca.
//
// Solidity: function deregisterOperator(_operator address) returns()
func (_Bridge *BridgeTransactor) DeregisterOperator(opts *bind.TransactOpts, _operator common.Address) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "deregisterOperator", _operator)
}

// DeregisterOperator is a paid mutator transaction binding the contract method 0xd8cf98ca.
//
// Solidity: function deregisterOperator(_operator address) returns()
func (_Bridge *BridgeSession) DeregisterOperator(_operator common.Address) (*types.Transaction, error) {
	return _Bridge.Contract.DeregisterOperator(&_Bridge.TransactOpts, _operator)
}

// DeregisterOperator is a paid mutator transaction binding the contract method 0xd8cf98ca.
//
// Solidity: function deregisterOperator(_operator address) returns()
func (_Bridge *BridgeTransactorSession) DeregisterOperator(_operator common.Address) (*types.Transaction, error) {
	return _Bridge.Contract.DeregisterOperator(&_Bridge.TransactOpts, _operator)
}

// DeregisterToken is a paid mutator transaction binding the contract method 0xbab2af1d.
//
// Solidity: function deregisterToken(_token address) returns()
func (_Bridge *BridgeTransactor) DeregisterToken(opts *bind.TransactOpts, _token common.Address) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "deregisterToken", _token)
}

// DeregisterToken is a paid mutator transaction binding the contract method 0xbab2af1d.
//
// Solidity: function deregisterToken(_token address) returns()
func (_Bridge *BridgeSession) DeregisterToken(_token common.Address) (*types.Transaction, error) {
	return _Bridge.Contract.DeregisterToken(&_Bridge.TransactOpts, _token)
}

// DeregisterToken is a paid mutator transaction binding the contract method 0xbab2af1d.
//
// Solidity: function deregisterToken(_token address) returns()
func (_Bridge *BridgeTransactorSession) DeregisterToken(_token common.Address) (*types.Transaction, error) {
	return _Bridge.Contract.DeregisterToken(&_Bridge.TransactOpts, _token)
}

// HandleERC20Transfer is a paid mutator transaction binding the contract method 0x452fbc62.
//
// Solidity: function handleERC20Transfer(_requestTxHash bytes32, _from address, _to address, _tokenAddress address, _value uint256, _requestedNonce uint64, _requestedBlockNumber uint64, _extraData uint256[]) returns()
func (_Bridge *BridgeTransactor) HandleERC20Transfer(opts *bind.TransactOpts, _requestTxHash [32]byte, _from common.Address, _to common.Address, _tokenAddress common.Address, _value *big.Int, _requestedNonce uint64, _requestedBlockNumber uint64, _extraData []*big.Int) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "handleERC20Transfer", _requestTxHash, _from, _to, _tokenAddress, _value, _requestedNonce, _requestedBlockNumber, _extraData)
}

// HandleERC20Transfer is a paid mutator transaction binding the contract method 0x452fbc62.
//
// Solidity: function handleERC20Transfer(_requestTxHash bytes32, _from address, _to address, _tokenAddress address, _value uint256, _requestedNonce uint64, _requestedBlockNumber uint64, _extraData uint256[]) returns()
func (_Bridge *BridgeSession) HandleERC20Transfer(_requestTxHash [32]byte, _from common.Address, _to common.Address, _tokenAddress common.Address, _value *big.Int, _requestedNonce uint64, _requestedBlockNumber uint64, _extraData []*big.Int) (*types.Transaction, error) {
	return _Bridge.Contract.HandleERC20Transfer(&_Bridge.TransactOpts, _requestTxHash, _from, _to, _tokenAddress, _value, _requestedNonce, _requestedBlockNumber, _extraData)
}

// HandleERC20Transfer is a paid mutator transaction binding the contract method 0x452fbc62.
//
// Solidity: function handleERC20Transfer(_requestTxHash bytes32, _from address, _to address, _tokenAddress address, _value uint256, _requestedNonce uint64, _requestedBlockNumber uint64, _extraData uint256[]) returns()
func (_Bridge *BridgeTransactorSession) HandleERC20Transfer(_requestTxHash [32]byte, _from common.Address, _to common.Address, _tokenAddress common.Address, _value *big.Int, _requestedNonce uint64, _requestedBlockNumber uint64, _extraData []*big.Int) (*types.Transaction, error) {
	return _Bridge.Contract.HandleERC20Transfer(&_Bridge.TransactOpts, _requestTxHash, _from, _to, _tokenAddress, _value, _requestedNonce, _requestedBlockNumber, _extraData)
}

// HandleERC721Transfer is a paid mutator transaction binding the contract method 0x2e4bb197.
//
// Solidity: function handleERC721Transfer(_requestTxHash bytes32, _from address, _to address, _tokenAddress address, _tokenId uint256, _requestedNonce uint64, _requestedBlockNumber uint64, _tokenURI string, _extraData uint256[]) returns()
func (_Bridge *BridgeTransactor) HandleERC721Transfer(opts *bind.TransactOpts, _requestTxHash [32]byte, _from common.Address, _to common.Address, _tokenAddress common.Address, _tokenId *big.Int, _requestedNonce uint64, _requestedBlockNumber uint64, _tokenURI string, _extraData []*big.Int) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "handleERC721Transfer", _requestTxHash, _from, _to, _tokenAddress, _tokenId, _requestedNonce, _requestedBlockNumber, _tokenURI, _extraData)
}

// HandleERC721Transfer is a paid mutator transaction binding the contract method 0x2e4bb197.
//
// Solidity: function handleERC721Transfer(_requestTxHash bytes32, _from address, _to address, _tokenAddress address, _tokenId uint256, _requestedNonce uint64, _requestedBlockNumber uint64, _tokenURI string, _extraData uint256[]) returns()
func (_Bridge *BridgeSession) HandleERC721Transfer(_requestTxHash [32]byte, _from common.Address, _to common.Address, _tokenAddress common.Address, _tokenId *big.Int, _requestedNonce uint64, _requestedBlockNumber uint64, _tokenURI string, _extraData []*big.Int) (*types.Transaction, error) {
	return _Bridge.Contract.HandleERC721Transfer(&_Bridge.TransactOpts, _requestTxHash, _from, _to, _tokenAddress, _tokenId, _requestedNonce, _requestedBlockNumber, _tokenURI, _extraData)
}

// HandleERC721Transfer is a paid mutator transaction binding the contract method 0x2e4bb197.
//
// Solidity: function handleERC721Transfer(_requestTxHash bytes32, _from address, _to address, _tokenAddress address, _tokenId uint256, _requestedNonce uint64, _requestedBlockNumber uint64, _tokenURI string, _extraData uint256[]) returns()
func (_Bridge *BridgeTransactorSession) HandleERC721Transfer(_requestTxHash [32]byte, _from common.Address, _to common.Address, _tokenAddress common.Address, _tokenId *big.Int, _requestedNonce uint64, _requestedBlockNumber uint64, _tokenURI string, _extraData []*big.Int) (*types.Transaction, error) {
	return _Bridge.Contract.HandleERC721Transfer(&_Bridge.TransactOpts, _requestTxHash, _from, _to, _tokenAddress, _tokenId, _requestedNonce, _requestedBlockNumber, _tokenURI, _extraData)
}

// HandleKLAYTransfer is a paid mutator transaction binding the contract method 0xee273ad7.
//
// Solidity: function handleKLAYTransfer(_requestTxHash bytes32, _from address, _to address, _value uint256, _requestedNonce uint64, _requestedBlockNumber uint64, _extraData uint256[]) returns()
func (_Bridge *BridgeTransactor) HandleKLAYTransfer(opts *bind.TransactOpts, _requestTxHash [32]byte, _from common.Address, _to common.Address, _value *big.Int, _requestedNonce uint64, _requestedBlockNumber uint64, _extraData []*big.Int) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "handleKLAYTransfer", _requestTxHash, _from, _to, _value, _requestedNonce, _requestedBlockNumber, _extraData)
}

// HandleKLAYTransfer is a paid mutator transaction binding the contract method 0xee273ad7.
//
// Solidity: function handleKLAYTransfer(_requestTxHash bytes32, _from address, _to address, _value uint256, _requestedNonce uint64, _requestedBlockNumber uint64, _extraData uint256[]) returns()
func (_Bridge *BridgeSession) HandleKLAYTransfer(_requestTxHash [32]byte, _from common.Address, _to common.Address, _value *big.Int, _requestedNonce uint64, _requestedBlockNumber uint64, _extraData []*big.Int) (*types.Transaction, error) {
	return _Bridge.Contract.HandleKLAYTransfer(&_Bridge.TransactOpts, _requestTxHash, _from, _to, _value, _requestedNonce, _requestedBlockNumber, _extraData)
}

// HandleKLAYTransfer is a paid mutator transaction binding the contract method 0xee273ad7.
//
// Solidity: function handleKLAYTransfer(_requestTxHash bytes32, _from address, _to address, _value uint256, _requestedNonce uint64, _requestedBlockNumber uint64, _extraData uint256[]) returns()
func (_Bridge *BridgeTransactorSession) HandleKLAYTransfer(_requestTxHash [32]byte, _from common.Address, _to common.Address, _value *big.Int, _requestedNonce uint64, _requestedBlockNumber uint64, _extraData []*big.Int) (*types.Transaction, error) {
	return _Bridge.Contract.HandleKLAYTransfer(&_Bridge.TransactOpts, _requestTxHash, _from, _to, _value, _requestedNonce, _requestedBlockNumber, _extraData)
}

// OnERC20Received is a paid mutator transaction binding the contract method 0x349449b2.
//
// Solidity: function onERC20Received(_from address, _value uint256, _to address, _feeLimit uint256, _extraData uint256[]) returns()
func (_Bridge *BridgeTransactor) OnERC20Received(opts *bind.TransactOpts, _from common.Address, _value *big.Int, _to common.Address, _feeLimit *big.Int, _extraData []*big.Int) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "onERC20Received", _from, _value, _to, _feeLimit, _extraData)
}

// OnERC20Received is a paid mutator transaction binding the contract method 0x349449b2.
//
// Solidity: function onERC20Received(_from address, _value uint256, _to address, _feeLimit uint256, _extraData uint256[]) returns()
func (_Bridge *BridgeSession) OnERC20Received(_from common.Address, _value *big.Int, _to common.Address, _feeLimit *big.Int, _extraData []*big.Int) (*types.Transaction, error) {
	return _Bridge.Contract.OnERC20Received(&_Bridge.TransactOpts, _from, _value, _to, _feeLimit, _extraData)
}

// OnERC20Received is a paid mutator transaction binding the contract method 0x349449b2.
//
// Solidity: function onERC20Received(_from address, _value uint256, _to address, _feeLimit uint256, _extraData uint256[]) returns()
func (_Bridge *BridgeTransactorSession) OnERC20Received(_from common.Address, _value *big.Int, _to common.Address, _feeLimit *big.Int, _extraData []*big.Int) (*types.Transaction, error) {
	return _Bridge.Contract.OnERC20Received(&_Bridge.TransactOpts, _from, _value, _to, _feeLimit, _extraData)
}

// OnERC721Received is a paid mutator transaction binding the contract method 0xa319de9e.
//
// Solidity: function onERC721Received(_from address, _tokenId uint256, _to address, _extraData uint256[]) returns()
func (_Bridge *BridgeTransactor) OnERC721Received(opts *bind.TransactOpts, _from common.Address, _tokenId *big.Int, _to common.Address, _extraData []*big.Int) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "onERC721Received", _from, _tokenId, _to, _extraData)
}

// OnERC721Received is a paid mutator transaction binding the contract method 0xa319de9e.
//
// Solidity: function onERC721Received(_from address, _tokenId uint256, _to address, _extraData uint256[]) returns()
func (_Bridge *BridgeSession) OnERC721Received(_from common.Address, _tokenId *big.Int, _to common.Address, _extraData []*big.Int) (*types.Transaction, error) {
	return _Bridge.Contract.OnERC721Received(&_Bridge.TransactOpts, _from, _tokenId, _to, _extraData)
}

// OnERC721Received is a paid mutator transaction binding the contract method 0xa319de9e.
//
// Solidity: function onERC721Received(_from address, _tokenId uint256, _to address, _extraData uint256[]) returns()
func (_Bridge *BridgeTransactorSession) OnERC721Received(_from common.Address, _tokenId *big.Int, _to common.Address, _extraData []*big.Int) (*types.Transaction, error) {
	return _Bridge.Contract.OnERC721Received(&_Bridge.TransactOpts, _from, _tokenId, _to, _extraData)
}

// RegisterOperator is a paid mutator transaction binding the contract method 0x3682a450.
//
// Solidity: function registerOperator(_operator address) returns()
func (_Bridge *BridgeTransactor) RegisterOperator(opts *bind.TransactOpts, _operator common.Address) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "registerOperator", _operator)
}

// RegisterOperator is a paid mutator transaction binding the contract method 0x3682a450.
//
// Solidity: function registerOperator(_operator address) returns()
func (_Bridge *BridgeSession) RegisterOperator(_operator common.Address) (*types.Transaction, error) {
	return _Bridge.Contract.RegisterOperator(&_Bridge.TransactOpts, _operator)
}

// RegisterOperator is a paid mutator transaction binding the contract method 0x3682a450.
//
// Solidity: function registerOperator(_operator address) returns()
func (_Bridge *BridgeTransactorSession) RegisterOperator(_operator common.Address) (*types.Transaction, error) {
	return _Bridge.Contract.RegisterOperator(&_Bridge.TransactOpts, _operator)
}

// RegisterToken is a paid mutator transaction binding the contract method 0x4739f7e5.
//
// Solidity: function registerToken(_token address, _cToken address) returns()
func (_Bridge *BridgeTransactor) RegisterToken(opts *bind.TransactOpts, _token common.Address, _cToken common.Address) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "registerToken", _token, _cToken)
}

// RegisterToken is a paid mutator transaction binding the contract method 0x4739f7e5.
//
// Solidity: function registerToken(_token address, _cToken address) returns()
func (_Bridge *BridgeSession) RegisterToken(_token common.Address, _cToken common.Address) (*types.Transaction, error) {
	return _Bridge.Contract.RegisterToken(&_Bridge.TransactOpts, _token, _cToken)
}

// RegisterToken is a paid mutator transaction binding the contract method 0x4739f7e5.
//
// Solidity: function registerToken(_token address, _cToken address) returns()
func (_Bridge *BridgeTransactorSession) RegisterToken(_token common.Address, _cToken common.Address) (*types.Transaction, error) {
	return _Bridge.Contract.RegisterToken(&_Bridge.TransactOpts, _token, _cToken)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Bridge *BridgeTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Bridge *BridgeSession) RenounceOwnership() (*types.Transaction, error) {
	return _Bridge.Contract.RenounceOwnership(&_Bridge.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Bridge *BridgeTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Bridge.Contract.RenounceOwnership(&_Bridge.TransactOpts)
}

// RequestERC20Transfer is a paid mutator transaction binding the contract method 0x74f41f15.
//
// Solidity: function requestERC20Transfer(_tokenAddress address, _to address, _value uint256, _feeLimit uint256, _extraData uint256[]) returns()
func (_Bridge *BridgeTransactor) RequestERC20Transfer(opts *bind.TransactOpts, _tokenAddress common.Address, _to common.Address, _value *big.Int, _feeLimit *big.Int, _extraData []*big.Int) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "requestERC20Transfer", _tokenAddress, _to, _value, _feeLimit, _extraData)
}

// RequestERC20Transfer is a paid mutator transaction binding the contract method 0x74f41f15.
//
// Solidity: function requestERC20Transfer(_tokenAddress address, _to address, _value uint256, _feeLimit uint256, _extraData uint256[]) returns()
func (_Bridge *BridgeSession) RequestERC20Transfer(_tokenAddress common.Address, _to common.Address, _value *big.Int, _feeLimit *big.Int, _extraData []*big.Int) (*types.Transaction, error) {
	return _Bridge.Contract.RequestERC20Transfer(&_Bridge.TransactOpts, _tokenAddress, _to, _value, _feeLimit, _extraData)
}

// RequestERC20Transfer is a paid mutator transaction binding the contract method 0x74f41f15.
//
// Solidity: function requestERC20Transfer(_tokenAddress address, _to address, _value uint256, _feeLimit uint256, _extraData uint256[]) returns()
func (_Bridge *BridgeTransactorSession) RequestERC20Transfer(_tokenAddress common.Address, _to common.Address, _value *big.Int, _feeLimit *big.Int, _extraData []*big.Int) (*types.Transaction, error) {
	return _Bridge.Contract.RequestERC20Transfer(&_Bridge.TransactOpts, _tokenAddress, _to, _value, _feeLimit, _extraData)
}

// RequestERC721Transfer is a paid mutator transaction binding the contract method 0x098540b8.
//
// Solidity: function requestERC721Transfer(_tokenAddress address, _to address, _tokenId uint256, _extraData uint256[]) returns()
func (_Bridge *BridgeTransactor) RequestERC721Transfer(opts *bind.TransactOpts, _tokenAddress common.Address, _to common.Address, _tokenId *big.Int, _extraData []*big.Int) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "requestERC721Transfer", _tokenAddress, _to, _tokenId, _extraData)
}

// RequestERC721Transfer is a paid mutator transaction binding the contract method 0x098540b8.
//
// Solidity: function requestERC721Transfer(_tokenAddress address, _to address, _tokenId uint256, _extraData uint256[]) returns()
func (_Bridge *BridgeSession) RequestERC721Transfer(_tokenAddress common.Address, _to common.Address, _tokenId *big.Int, _extraData []*big.Int) (*types.Transaction, error) {
	return _Bridge.Contract.RequestERC721Transfer(&_Bridge.TransactOpts, _tokenAddress, _to, _tokenId, _extraData)
}

// RequestERC721Transfer is a paid mutator transaction binding the contract method 0x098540b8.
//
// Solidity: function requestERC721Transfer(_tokenAddress address, _to address, _tokenId uint256, _extraData uint256[]) returns()
func (_Bridge *BridgeTransactorSession) RequestERC721Transfer(_tokenAddress common.Address, _to common.Address, _tokenId *big.Int, _extraData []*big.Int) (*types.Transaction, error) {
	return _Bridge.Contract.RequestERC721Transfer(&_Bridge.TransactOpts, _tokenAddress, _to, _tokenId, _extraData)
}

// RequestKLAYTransfer is a paid mutator transaction binding the contract method 0xeb7f0cbe.
//
// Solidity: function requestKLAYTransfer(_to address, _value uint256, _extraData uint256[]) returns()
func (_Bridge *BridgeTransactor) RequestKLAYTransfer(opts *bind.TransactOpts, _to common.Address, _value *big.Int, _extraData []*big.Int) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "requestKLAYTransfer", _to, _value, _extraData)
}

// RequestKLAYTransfer is a paid mutator transaction binding the contract method 0xeb7f0cbe.
//
// Solidity: function requestKLAYTransfer(_to address, _value uint256, _extraData uint256[]) returns()
func (_Bridge *BridgeSession) RequestKLAYTransfer(_to common.Address, _value *big.Int, _extraData []*big.Int) (*types.Transaction, error) {
	return _Bridge.Contract.RequestKLAYTransfer(&_Bridge.TransactOpts, _to, _value, _extraData)
}

// RequestKLAYTransfer is a paid mutator transaction binding the contract method 0xeb7f0cbe.
//
// Solidity: function requestKLAYTransfer(_to address, _value uint256, _extraData uint256[]) returns()
func (_Bridge *BridgeTransactorSession) RequestKLAYTransfer(_to common.Address, _value *big.Int, _extraData []*big.Int) (*types.Transaction, error) {
	return _Bridge.Contract.RequestKLAYTransfer(&_Bridge.TransactOpts, _to, _value, _extraData)
}

// SetCounterPartBridge is a paid mutator transaction binding the contract method 0x87b04c55.
//
// Solidity: function setCounterPartBridge(_bridge address) returns()
func (_Bridge *BridgeTransactor) SetCounterPartBridge(opts *bind.TransactOpts, _bridge common.Address) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "setCounterPartBridge", _bridge)
}

// SetCounterPartBridge is a paid mutator transaction binding the contract method 0x87b04c55.
//
// Solidity: function setCounterPartBridge(_bridge address) returns()
func (_Bridge *BridgeSession) SetCounterPartBridge(_bridge common.Address) (*types.Transaction, error) {
	return _Bridge.Contract.SetCounterPartBridge(&_Bridge.TransactOpts, _bridge)
}

// SetCounterPartBridge is a paid mutator transaction binding the contract method 0x87b04c55.
//
// Solidity: function setCounterPartBridge(_bridge address) returns()
func (_Bridge *BridgeTransactorSession) SetCounterPartBridge(_bridge common.Address) (*types.Transaction, error) {
	return _Bridge.Contract.SetCounterPartBridge(&_Bridge.TransactOpts, _bridge)
}

// SetERC20Fee is a paid mutator transaction binding the contract method 0x2f88396c.
//
// Solidity: function setERC20Fee(_token address, _fee uint256, _requestNonce uint64) returns()
func (_Bridge *BridgeTransactor) SetERC20Fee(opts *bind.TransactOpts, _token common.Address, _fee *big.Int, _requestNonce uint64) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "setERC20Fee", _token, _fee, _requestNonce)
}

// SetERC20Fee is a paid mutator transaction binding the contract method 0x2f88396c.
//
// Solidity: function setERC20Fee(_token address, _fee uint256, _requestNonce uint64) returns()
func (_Bridge *BridgeSession) SetERC20Fee(_token common.Address, _fee *big.Int, _requestNonce uint64) (*types.Transaction, error) {
	return _Bridge.Contract.SetERC20Fee(&_Bridge.TransactOpts, _token, _fee, _requestNonce)
}

// SetERC20Fee is a paid mutator transaction binding the contract method 0x2f88396c.
//
// Solidity: function setERC20Fee(_token address, _fee uint256, _requestNonce uint64) returns()
func (_Bridge *BridgeTransactorSession) SetERC20Fee(_token common.Address, _fee *big.Int, _requestNonce uint64) (*types.Transaction, error) {
	return _Bridge.Contract.SetERC20Fee(&_Bridge.TransactOpts, _token, _fee, _requestNonce)
}

// SetFeeReceiver is a paid mutator transaction binding the contract method 0xefdcd974.
//
// Solidity: function setFeeReceiver(_feeReceiver address) returns()
func (_Bridge *BridgeTransactor) SetFeeReceiver(opts *bind.TransactOpts, _feeReceiver common.Address) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "setFeeReceiver", _feeReceiver)
}

// SetFeeReceiver is a paid mutator transaction binding the contract method 0xefdcd974.
//
// Solidity: function setFeeReceiver(_feeReceiver address) returns()
func (_Bridge *BridgeSession) SetFeeReceiver(_feeReceiver common.Address) (*types.Transaction, error) {
	return _Bridge.Contract.SetFeeReceiver(&_Bridge.TransactOpts, _feeReceiver)
}

// SetFeeReceiver is a paid mutator transaction binding the contract method 0xefdcd974.
//
// Solidity: function setFeeReceiver(_feeReceiver address) returns()
func (_Bridge *BridgeTransactorSession) SetFeeReceiver(_feeReceiver common.Address) (*types.Transaction, error) {
	return _Bridge.Contract.SetFeeReceiver(&_Bridge.TransactOpts, _feeReceiver)
}

// SetKLAYFee is a paid mutator transaction binding the contract method 0x1a2ae53e.
//
// Solidity: function setKLAYFee(_fee uint256, _requestNonce uint64) returns()
func (_Bridge *BridgeTransactor) SetKLAYFee(opts *bind.TransactOpts, _fee *big.Int, _requestNonce uint64) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "setKLAYFee", _fee, _requestNonce)
}

// SetKLAYFee is a paid mutator transaction binding the contract method 0x1a2ae53e.
//
// Solidity: function setKLAYFee(_fee uint256, _requestNonce uint64) returns()
func (_Bridge *BridgeSession) SetKLAYFee(_fee *big.Int, _requestNonce uint64) (*types.Transaction, error) {
	return _Bridge.Contract.SetKLAYFee(&_Bridge.TransactOpts, _fee, _requestNonce)
}

// SetKLAYFee is a paid mutator transaction binding the contract method 0x1a2ae53e.
//
// Solidity: function setKLAYFee(_fee uint256, _requestNonce uint64) returns()
func (_Bridge *BridgeTransactorSession) SetKLAYFee(_fee *big.Int, _requestNonce uint64) (*types.Transaction, error) {
	return _Bridge.Contract.SetKLAYFee(&_Bridge.TransactOpts, _fee, _requestNonce)
}

// SetOperatorThreshold is a paid mutator transaction binding the contract method 0xee2aec65.
//
// Solidity: function setOperatorThreshold(_voteType uint8, _threshold uint8) returns()
func (_Bridge *BridgeTransactor) SetOperatorThreshold(opts *bind.TransactOpts, _voteType uint8, _threshold uint8) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "setOperatorThreshold", _voteType, _threshold)
}

// SetOperatorThreshold is a paid mutator transaction binding the contract method 0xee2aec65.
//
// Solidity: function setOperatorThreshold(_voteType uint8, _threshold uint8) returns()
func (_Bridge *BridgeSession) SetOperatorThreshold(_voteType uint8, _threshold uint8) (*types.Transaction, error) {
	return _Bridge.Contract.SetOperatorThreshold(&_Bridge.TransactOpts, _voteType, _threshold)
}

// SetOperatorThreshold is a paid mutator transaction binding the contract method 0xee2aec65.
//
// Solidity: function setOperatorThreshold(_voteType uint8, _threshold uint8) returns()
func (_Bridge *BridgeTransactorSession) SetOperatorThreshold(_voteType uint8, _threshold uint8) (*types.Transaction, error) {
	return _Bridge.Contract.SetOperatorThreshold(&_Bridge.TransactOpts, _voteType, _threshold)
}

// Start is a paid mutator transaction binding the contract method 0xc877cf37.
//
// Solidity: function start(_status bool) returns()
func (_Bridge *BridgeTransactor) Start(opts *bind.TransactOpts, _status bool) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "start", _status)
}

// Start is a paid mutator transaction binding the contract method 0xc877cf37.
//
// Solidity: function start(_status bool) returns()
func (_Bridge *BridgeSession) Start(_status bool) (*types.Transaction, error) {
	return _Bridge.Contract.Start(&_Bridge.TransactOpts, _status)
}

// Start is a paid mutator transaction binding the contract method 0xc877cf37.
//
// Solidity: function start(_status bool) returns()
func (_Bridge *BridgeTransactorSession) Start(_status bool) (*types.Transaction, error) {
	return _Bridge.Contract.Start(&_Bridge.TransactOpts, _status)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(newOwner address) returns()
func (_Bridge *BridgeTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(newOwner address) returns()
func (_Bridge *BridgeSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Bridge.Contract.TransferOwnership(&_Bridge.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(newOwner address) returns()
func (_Bridge *BridgeTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Bridge.Contract.TransferOwnership(&_Bridge.TransactOpts, newOwner)
}

// BridgeERC20FeeChangedIterator is returned from FilterERC20FeeChanged and is used to iterate over the raw logs and unpacked data for ERC20FeeChanged events raised by the Bridge contract.
type BridgeERC20FeeChangedIterator struct {
	Event *BridgeERC20FeeChanged // Event containing the contract specifics and raw log

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
func (it *BridgeERC20FeeChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BridgeERC20FeeChanged)
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
		it.Event = new(BridgeERC20FeeChanged)
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
func (it *BridgeERC20FeeChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BridgeERC20FeeChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BridgeERC20FeeChanged represents a ERC20FeeChanged event raised by the Bridge contract.
type BridgeERC20FeeChanged struct {
	Token common.Address
	Fee   *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterERC20FeeChanged is a free log retrieval operation binding the contract event 0xdb5ad2e76ae20cfa4e7adbc7305d7538442164d85ead9937c98620a1aa4c255b.
//
// Solidity: e ERC20FeeChanged(token address, fee indexed uint256)
func (_Bridge *BridgeFilterer) FilterERC20FeeChanged(opts *bind.FilterOpts, fee []*big.Int) (*BridgeERC20FeeChangedIterator, error) {

	var feeRule []interface{}
	for _, feeItem := range fee {
		feeRule = append(feeRule, feeItem)
	}

	logs, sub, err := _Bridge.contract.FilterLogs(opts, "ERC20FeeChanged", feeRule)
	if err != nil {
		return nil, err
	}
	return &BridgeERC20FeeChangedIterator{contract: _Bridge.contract, event: "ERC20FeeChanged", logs: logs, sub: sub}, nil
}

// WatchERC20FeeChanged is a free log subscription operation binding the contract event 0xdb5ad2e76ae20cfa4e7adbc7305d7538442164d85ead9937c98620a1aa4c255b.
//
// Solidity: e ERC20FeeChanged(token address, fee indexed uint256)
func (_Bridge *BridgeFilterer) WatchERC20FeeChanged(opts *bind.WatchOpts, sink chan<- *BridgeERC20FeeChanged, fee []*big.Int) (event.Subscription, error) {

	var feeRule []interface{}
	for _, feeItem := range fee {
		feeRule = append(feeRule, feeItem)
	}

	logs, sub, err := _Bridge.contract.WatchLogs(opts, "ERC20FeeChanged", feeRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BridgeERC20FeeChanged)
				if err := _Bridge.contract.UnpackLog(event, "ERC20FeeChanged", log); err != nil {
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

// BridgeFeeReceiverChangedIterator is returned from FilterFeeReceiverChanged and is used to iterate over the raw logs and unpacked data for FeeReceiverChanged events raised by the Bridge contract.
type BridgeFeeReceiverChangedIterator struct {
	Event *BridgeFeeReceiverChanged // Event containing the contract specifics and raw log

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
func (it *BridgeFeeReceiverChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BridgeFeeReceiverChanged)
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
		it.Event = new(BridgeFeeReceiverChanged)
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
func (it *BridgeFeeReceiverChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BridgeFeeReceiverChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BridgeFeeReceiverChanged represents a FeeReceiverChanged event raised by the Bridge contract.
type BridgeFeeReceiverChanged struct {
	FeeReceiver common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterFeeReceiverChanged is a free log retrieval operation binding the contract event 0x647672599d3468abcfa241a13c9e3d34383caadb5cc80fb67c3cdfcd5f786059.
//
// Solidity: e FeeReceiverChanged(feeReceiver indexed address)
func (_Bridge *BridgeFilterer) FilterFeeReceiverChanged(opts *bind.FilterOpts, feeReceiver []common.Address) (*BridgeFeeReceiverChangedIterator, error) {

	var feeReceiverRule []interface{}
	for _, feeReceiverItem := range feeReceiver {
		feeReceiverRule = append(feeReceiverRule, feeReceiverItem)
	}

	logs, sub, err := _Bridge.contract.FilterLogs(opts, "FeeReceiverChanged", feeReceiverRule)
	if err != nil {
		return nil, err
	}
	return &BridgeFeeReceiverChangedIterator{contract: _Bridge.contract, event: "FeeReceiverChanged", logs: logs, sub: sub}, nil
}

// WatchFeeReceiverChanged is a free log subscription operation binding the contract event 0x647672599d3468abcfa241a13c9e3d34383caadb5cc80fb67c3cdfcd5f786059.
//
// Solidity: e FeeReceiverChanged(feeReceiver indexed address)
func (_Bridge *BridgeFilterer) WatchFeeReceiverChanged(opts *bind.WatchOpts, sink chan<- *BridgeFeeReceiverChanged, feeReceiver []common.Address) (event.Subscription, error) {

	var feeReceiverRule []interface{}
	for _, feeReceiverItem := range feeReceiver {
		feeReceiverRule = append(feeReceiverRule, feeReceiverItem)
	}

	logs, sub, err := _Bridge.contract.WatchLogs(opts, "FeeReceiverChanged", feeReceiverRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BridgeFeeReceiverChanged)
				if err := _Bridge.contract.UnpackLog(event, "FeeReceiverChanged", log); err != nil {
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

// BridgeHandleValueTransferIterator is returned from FilterHandleValueTransfer and is used to iterate over the raw logs and unpacked data for HandleValueTransfer events raised by the Bridge contract.
type BridgeHandleValueTransferIterator struct {
	Event *BridgeHandleValueTransfer // Event containing the contract specifics and raw log

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
func (it *BridgeHandleValueTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BridgeHandleValueTransfer)
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
		it.Event = new(BridgeHandleValueTransfer)
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
func (it *BridgeHandleValueTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BridgeHandleValueTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BridgeHandleValueTransfer represents a HandleValueTransfer event raised by the Bridge contract.
type BridgeHandleValueTransfer struct {
	RequestTxHash  [32]byte
	TokenType      uint8
	From           common.Address
	To             common.Address
	TokenAddress   common.Address
	ValueOrTokenId *big.Int
	HandleNonce    uint64
	ExtraData      []*big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterHandleValueTransfer is a free log retrieval operation binding the contract event 0x91f1d4b9864c9c091ade654e58899919d2a56e194153899c7e292cd3ae2ac47b.
//
// Solidity: e HandleValueTransfer(requestTxHash bytes32, tokenType uint8, from address, to address, tokenAddress address, valueOrTokenId uint256, handleNonce uint64, extraData uint256[])
func (_Bridge *BridgeFilterer) FilterHandleValueTransfer(opts *bind.FilterOpts) (*BridgeHandleValueTransferIterator, error) {

	logs, sub, err := _Bridge.contract.FilterLogs(opts, "HandleValueTransfer")
	if err != nil {
		return nil, err
	}
	return &BridgeHandleValueTransferIterator{contract: _Bridge.contract, event: "HandleValueTransfer", logs: logs, sub: sub}, nil
}

// WatchHandleValueTransfer is a free log subscription operation binding the contract event 0x91f1d4b9864c9c091ade654e58899919d2a56e194153899c7e292cd3ae2ac47b.
//
// Solidity: e HandleValueTransfer(requestTxHash bytes32, tokenType uint8, from address, to address, tokenAddress address, valueOrTokenId uint256, handleNonce uint64, extraData uint256[])
func (_Bridge *BridgeFilterer) WatchHandleValueTransfer(opts *bind.WatchOpts, sink chan<- *BridgeHandleValueTransfer) (event.Subscription, error) {

	logs, sub, err := _Bridge.contract.WatchLogs(opts, "HandleValueTransfer")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BridgeHandleValueTransfer)
				if err := _Bridge.contract.UnpackLog(event, "HandleValueTransfer", log); err != nil {
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

// BridgeKLAYFeeChangedIterator is returned from FilterKLAYFeeChanged and is used to iterate over the raw logs and unpacked data for KLAYFeeChanged events raised by the Bridge contract.
type BridgeKLAYFeeChangedIterator struct {
	Event *BridgeKLAYFeeChanged // Event containing the contract specifics and raw log

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
func (it *BridgeKLAYFeeChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BridgeKLAYFeeChanged)
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
		it.Event = new(BridgeKLAYFeeChanged)
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
func (it *BridgeKLAYFeeChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BridgeKLAYFeeChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BridgeKLAYFeeChanged represents a KLAYFeeChanged event raised by the Bridge contract.
type BridgeKLAYFeeChanged struct {
	Fee *big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// FilterKLAYFeeChanged is a free log retrieval operation binding the contract event 0xa7a33d0996347e1aa55ca2206015b61b9534bdd881d59d59aa680e25eefac365.
//
// Solidity: e KLAYFeeChanged(fee indexed uint256)
func (_Bridge *BridgeFilterer) FilterKLAYFeeChanged(opts *bind.FilterOpts, fee []*big.Int) (*BridgeKLAYFeeChangedIterator, error) {

	var feeRule []interface{}
	for _, feeItem := range fee {
		feeRule = append(feeRule, feeItem)
	}

	logs, sub, err := _Bridge.contract.FilterLogs(opts, "KLAYFeeChanged", feeRule)
	if err != nil {
		return nil, err
	}
	return &BridgeKLAYFeeChangedIterator{contract: _Bridge.contract, event: "KLAYFeeChanged", logs: logs, sub: sub}, nil
}

// WatchKLAYFeeChanged is a free log subscription operation binding the contract event 0xa7a33d0996347e1aa55ca2206015b61b9534bdd881d59d59aa680e25eefac365.
//
// Solidity: e KLAYFeeChanged(fee indexed uint256)
func (_Bridge *BridgeFilterer) WatchKLAYFeeChanged(opts *bind.WatchOpts, sink chan<- *BridgeKLAYFeeChanged, fee []*big.Int) (event.Subscription, error) {

	var feeRule []interface{}
	for _, feeItem := range fee {
		feeRule = append(feeRule, feeItem)
	}

	logs, sub, err := _Bridge.contract.WatchLogs(opts, "KLAYFeeChanged", feeRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BridgeKLAYFeeChanged)
				if err := _Bridge.contract.UnpackLog(event, "KLAYFeeChanged", log); err != nil {
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

// BridgeOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Bridge contract.
type BridgeOwnershipTransferredIterator struct {
	Event *BridgeOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *BridgeOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BridgeOwnershipTransferred)
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
		it.Event = new(BridgeOwnershipTransferred)
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
func (it *BridgeOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BridgeOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BridgeOwnershipTransferred represents a OwnershipTransferred event raised by the Bridge contract.
type BridgeOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: e OwnershipTransferred(previousOwner indexed address, newOwner indexed address)
func (_Bridge *BridgeFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*BridgeOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Bridge.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &BridgeOwnershipTransferredIterator{contract: _Bridge.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: e OwnershipTransferred(previousOwner indexed address, newOwner indexed address)
func (_Bridge *BridgeFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *BridgeOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Bridge.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BridgeOwnershipTransferred)
				if err := _Bridge.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// BridgeRequestValueTransferIterator is returned from FilterRequestValueTransfer and is used to iterate over the raw logs and unpacked data for RequestValueTransfer events raised by the Bridge contract.
type BridgeRequestValueTransferIterator struct {
	Event *BridgeRequestValueTransfer // Event containing the contract specifics and raw log

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
func (it *BridgeRequestValueTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BridgeRequestValueTransfer)
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
		it.Event = new(BridgeRequestValueTransfer)
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
func (it *BridgeRequestValueTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BridgeRequestValueTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BridgeRequestValueTransfer represents a RequestValueTransfer event raised by the Bridge contract.
type BridgeRequestValueTransfer struct {
	TokenType      uint8
	From           common.Address
	To             common.Address
	TokenAddress   common.Address
	ValueOrTokenId *big.Int
	RequestNonce   uint64
	Uri            string
	Fee            *big.Int
	ExtraData      []*big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterRequestValueTransfer is a free log retrieval operation binding the contract event 0x09a722fc78123f7b2ba788ae83b01221477757cfb4b83302421e7c5c1fd01da9.
//
// Solidity: e RequestValueTransfer(tokenType uint8, from address, to address, tokenAddress address, valueOrTokenId uint256, requestNonce uint64, uri string, fee uint256, extraData uint256[])
func (_Bridge *BridgeFilterer) FilterRequestValueTransfer(opts *bind.FilterOpts) (*BridgeRequestValueTransferIterator, error) {

	logs, sub, err := _Bridge.contract.FilterLogs(opts, "RequestValueTransfer")
	if err != nil {
		return nil, err
	}
	return &BridgeRequestValueTransferIterator{contract: _Bridge.contract, event: "RequestValueTransfer", logs: logs, sub: sub}, nil
}

// WatchRequestValueTransfer is a free log subscription operation binding the contract event 0x09a722fc78123f7b2ba788ae83b01221477757cfb4b83302421e7c5c1fd01da9.
//
// Solidity: e RequestValueTransfer(tokenType uint8, from address, to address, tokenAddress address, valueOrTokenId uint256, requestNonce uint64, uri string, fee uint256, extraData uint256[])
func (_Bridge *BridgeFilterer) WatchRequestValueTransfer(opts *bind.WatchOpts, sink chan<- *BridgeRequestValueTransfer) (event.Subscription, error) {

	logs, sub, err := _Bridge.contract.WatchLogs(opts, "RequestValueTransfer")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BridgeRequestValueTransfer)
				if err := _Bridge.contract.UnpackLog(event, "RequestValueTransfer", log); err != nil {
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

// BridgeFeeABI is the input ABI used to generate the binding from.
const BridgeFeeABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"feeOfERC20\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"feeReceiver\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"feeOfKLAY\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"_feeReceiver\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"fee\",\"type\":\"uint256\"}],\"name\":\"KLAYFeeChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"token\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"fee\",\"type\":\"uint256\"}],\"name\":\"ERC20FeeChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"feeReceiver\",\"type\":\"address\"}],\"name\":\"FeeReceiverChanged\",\"type\":\"event\"}]"

// BridgeFeeBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const BridgeFeeBinRuntime = `0x`

// BridgeFeeBin is the compiled bytecode used for deploying new contracts.
const BridgeFeeBin = `0x`

// DeployBridgeFee deploys a new Klaytn contract, binding an instance of BridgeFee to it.
func DeployBridgeFee(auth *bind.TransactOpts, backend bind.ContractBackend, _feeReceiver common.Address) (common.Address, *types.Transaction, *BridgeFee, error) {
	parsed, err := abi.JSON(strings.NewReader(BridgeFeeABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(BridgeFeeBin), backend, _feeReceiver)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &BridgeFee{BridgeFeeCaller: BridgeFeeCaller{contract: contract}, BridgeFeeTransactor: BridgeFeeTransactor{contract: contract}, BridgeFeeFilterer: BridgeFeeFilterer{contract: contract}}, nil
}

// BridgeFee is an auto generated Go binding around a Klaytn contract.
type BridgeFee struct {
	BridgeFeeCaller     // Read-only binding to the contract
	BridgeFeeTransactor // Write-only binding to the contract
	BridgeFeeFilterer   // Log filterer for contract events
}

// BridgeFeeCaller is an auto generated read-only Go binding around a Klaytn contract.
type BridgeFeeCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BridgeFeeTransactor is an auto generated write-only Go binding around a Klaytn contract.
type BridgeFeeTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BridgeFeeFilterer is an auto generated log filtering Go binding around a Klaytn contract events.
type BridgeFeeFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BridgeFeeSession is an auto generated Go binding around a Klaytn contract,
// with pre-set call and transact options.
type BridgeFeeSession struct {
	Contract     *BridgeFee        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// BridgeFeeCallerSession is an auto generated read-only Go binding around a Klaytn contract,
// with pre-set call options.
type BridgeFeeCallerSession struct {
	Contract *BridgeFeeCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// BridgeFeeTransactorSession is an auto generated write-only Go binding around a Klaytn contract,
// with pre-set transact options.
type BridgeFeeTransactorSession struct {
	Contract     *BridgeFeeTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// BridgeFeeRaw is an auto generated low-level Go binding around a Klaytn contract.
type BridgeFeeRaw struct {
	Contract *BridgeFee // Generic contract binding to access the raw methods on
}

// BridgeFeeCallerRaw is an auto generated low-level read-only Go binding around a Klaytn contract.
type BridgeFeeCallerRaw struct {
	Contract *BridgeFeeCaller // Generic read-only contract binding to access the raw methods on
}

// BridgeFeeTransactorRaw is an auto generated low-level write-only Go binding around a Klaytn contract.
type BridgeFeeTransactorRaw struct {
	Contract *BridgeFeeTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBridgeFee creates a new instance of BridgeFee, bound to a specific deployed contract.
func NewBridgeFee(address common.Address, backend bind.ContractBackend) (*BridgeFee, error) {
	contract, err := bindBridgeFee(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &BridgeFee{BridgeFeeCaller: BridgeFeeCaller{contract: contract}, BridgeFeeTransactor: BridgeFeeTransactor{contract: contract}, BridgeFeeFilterer: BridgeFeeFilterer{contract: contract}}, nil
}

// NewBridgeFeeCaller creates a new read-only instance of BridgeFee, bound to a specific deployed contract.
func NewBridgeFeeCaller(address common.Address, caller bind.ContractCaller) (*BridgeFeeCaller, error) {
	contract, err := bindBridgeFee(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BridgeFeeCaller{contract: contract}, nil
}

// NewBridgeFeeTransactor creates a new write-only instance of BridgeFee, bound to a specific deployed contract.
func NewBridgeFeeTransactor(address common.Address, transactor bind.ContractTransactor) (*BridgeFeeTransactor, error) {
	contract, err := bindBridgeFee(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BridgeFeeTransactor{contract: contract}, nil
}

// NewBridgeFeeFilterer creates a new log filterer instance of BridgeFee, bound to a specific deployed contract.
func NewBridgeFeeFilterer(address common.Address, filterer bind.ContractFilterer) (*BridgeFeeFilterer, error) {
	contract, err := bindBridgeFee(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BridgeFeeFilterer{contract: contract}, nil
}

// bindBridgeFee binds a generic wrapper to an already deployed contract.
func bindBridgeFee(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(BridgeFeeABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BridgeFee *BridgeFeeRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _BridgeFee.Contract.BridgeFeeCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BridgeFee *BridgeFeeRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BridgeFee.Contract.BridgeFeeTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BridgeFee *BridgeFeeRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BridgeFee.Contract.BridgeFeeTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BridgeFee *BridgeFeeCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _BridgeFee.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BridgeFee *BridgeFeeTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BridgeFee.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BridgeFee *BridgeFeeTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BridgeFee.Contract.contract.Transact(opts, method, params...)
}

// FeeOfERC20 is a free data retrieval call binding the contract method 0x488af871.
//
// Solidity: function feeOfERC20( address) constant returns(uint256)
func (_BridgeFee *BridgeFeeCaller) FeeOfERC20(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _BridgeFee.contract.Call(opts, out, "feeOfERC20", arg0)
	return *ret0, err
}

// FeeOfERC20 is a free data retrieval call binding the contract method 0x488af871.
//
// Solidity: function feeOfERC20( address) constant returns(uint256)
func (_BridgeFee *BridgeFeeSession) FeeOfERC20(arg0 common.Address) (*big.Int, error) {
	return _BridgeFee.Contract.FeeOfERC20(&_BridgeFee.CallOpts, arg0)
}

// FeeOfERC20 is a free data retrieval call binding the contract method 0x488af871.
//
// Solidity: function feeOfERC20( address) constant returns(uint256)
func (_BridgeFee *BridgeFeeCallerSession) FeeOfERC20(arg0 common.Address) (*big.Int, error) {
	return _BridgeFee.Contract.FeeOfERC20(&_BridgeFee.CallOpts, arg0)
}

// FeeOfKLAY is a free data retrieval call binding the contract method 0xc263b5d6.
//
// Solidity: function feeOfKLAY() constant returns(uint256)
func (_BridgeFee *BridgeFeeCaller) FeeOfKLAY(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _BridgeFee.contract.Call(opts, out, "feeOfKLAY")
	return *ret0, err
}

// FeeOfKLAY is a free data retrieval call binding the contract method 0xc263b5d6.
//
// Solidity: function feeOfKLAY() constant returns(uint256)
func (_BridgeFee *BridgeFeeSession) FeeOfKLAY() (*big.Int, error) {
	return _BridgeFee.Contract.FeeOfKLAY(&_BridgeFee.CallOpts)
}

// FeeOfKLAY is a free data retrieval call binding the contract method 0xc263b5d6.
//
// Solidity: function feeOfKLAY() constant returns(uint256)
func (_BridgeFee *BridgeFeeCallerSession) FeeOfKLAY() (*big.Int, error) {
	return _BridgeFee.Contract.FeeOfKLAY(&_BridgeFee.CallOpts)
}

// FeeReceiver is a free data retrieval call binding the contract method 0xb3f00674.
//
// Solidity: function feeReceiver() constant returns(address)
func (_BridgeFee *BridgeFeeCaller) FeeReceiver(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _BridgeFee.contract.Call(opts, out, "feeReceiver")
	return *ret0, err
}

// FeeReceiver is a free data retrieval call binding the contract method 0xb3f00674.
//
// Solidity: function feeReceiver() constant returns(address)
func (_BridgeFee *BridgeFeeSession) FeeReceiver() (common.Address, error) {
	return _BridgeFee.Contract.FeeReceiver(&_BridgeFee.CallOpts)
}

// FeeReceiver is a free data retrieval call binding the contract method 0xb3f00674.
//
// Solidity: function feeReceiver() constant returns(address)
func (_BridgeFee *BridgeFeeCallerSession) FeeReceiver() (common.Address, error) {
	return _BridgeFee.Contract.FeeReceiver(&_BridgeFee.CallOpts)
}

// BridgeFeeERC20FeeChangedIterator is returned from FilterERC20FeeChanged and is used to iterate over the raw logs and unpacked data for ERC20FeeChanged events raised by the BridgeFee contract.
type BridgeFeeERC20FeeChangedIterator struct {
	Event *BridgeFeeERC20FeeChanged // Event containing the contract specifics and raw log

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
func (it *BridgeFeeERC20FeeChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BridgeFeeERC20FeeChanged)
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
		it.Event = new(BridgeFeeERC20FeeChanged)
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
func (it *BridgeFeeERC20FeeChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BridgeFeeERC20FeeChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BridgeFeeERC20FeeChanged represents a ERC20FeeChanged event raised by the BridgeFee contract.
type BridgeFeeERC20FeeChanged struct {
	Token common.Address
	Fee   *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterERC20FeeChanged is a free log retrieval operation binding the contract event 0xdb5ad2e76ae20cfa4e7adbc7305d7538442164d85ead9937c98620a1aa4c255b.
//
// Solidity: e ERC20FeeChanged(token address, fee indexed uint256)
func (_BridgeFee *BridgeFeeFilterer) FilterERC20FeeChanged(opts *bind.FilterOpts, fee []*big.Int) (*BridgeFeeERC20FeeChangedIterator, error) {

	var feeRule []interface{}
	for _, feeItem := range fee {
		feeRule = append(feeRule, feeItem)
	}

	logs, sub, err := _BridgeFee.contract.FilterLogs(opts, "ERC20FeeChanged", feeRule)
	if err != nil {
		return nil, err
	}
	return &BridgeFeeERC20FeeChangedIterator{contract: _BridgeFee.contract, event: "ERC20FeeChanged", logs: logs, sub: sub}, nil
}

// WatchERC20FeeChanged is a free log subscription operation binding the contract event 0xdb5ad2e76ae20cfa4e7adbc7305d7538442164d85ead9937c98620a1aa4c255b.
//
// Solidity: e ERC20FeeChanged(token address, fee indexed uint256)
func (_BridgeFee *BridgeFeeFilterer) WatchERC20FeeChanged(opts *bind.WatchOpts, sink chan<- *BridgeFeeERC20FeeChanged, fee []*big.Int) (event.Subscription, error) {

	var feeRule []interface{}
	for _, feeItem := range fee {
		feeRule = append(feeRule, feeItem)
	}

	logs, sub, err := _BridgeFee.contract.WatchLogs(opts, "ERC20FeeChanged", feeRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BridgeFeeERC20FeeChanged)
				if err := _BridgeFee.contract.UnpackLog(event, "ERC20FeeChanged", log); err != nil {
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

// BridgeFeeFeeReceiverChangedIterator is returned from FilterFeeReceiverChanged and is used to iterate over the raw logs and unpacked data for FeeReceiverChanged events raised by the BridgeFee contract.
type BridgeFeeFeeReceiverChangedIterator struct {
	Event *BridgeFeeFeeReceiverChanged // Event containing the contract specifics and raw log

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
func (it *BridgeFeeFeeReceiverChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BridgeFeeFeeReceiverChanged)
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
		it.Event = new(BridgeFeeFeeReceiverChanged)
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
func (it *BridgeFeeFeeReceiverChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BridgeFeeFeeReceiverChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BridgeFeeFeeReceiverChanged represents a FeeReceiverChanged event raised by the BridgeFee contract.
type BridgeFeeFeeReceiverChanged struct {
	FeeReceiver common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterFeeReceiverChanged is a free log retrieval operation binding the contract event 0x647672599d3468abcfa241a13c9e3d34383caadb5cc80fb67c3cdfcd5f786059.
//
// Solidity: e FeeReceiverChanged(feeReceiver indexed address)
func (_BridgeFee *BridgeFeeFilterer) FilterFeeReceiverChanged(opts *bind.FilterOpts, feeReceiver []common.Address) (*BridgeFeeFeeReceiverChangedIterator, error) {

	var feeReceiverRule []interface{}
	for _, feeReceiverItem := range feeReceiver {
		feeReceiverRule = append(feeReceiverRule, feeReceiverItem)
	}

	logs, sub, err := _BridgeFee.contract.FilterLogs(opts, "FeeReceiverChanged", feeReceiverRule)
	if err != nil {
		return nil, err
	}
	return &BridgeFeeFeeReceiverChangedIterator{contract: _BridgeFee.contract, event: "FeeReceiverChanged", logs: logs, sub: sub}, nil
}

// WatchFeeReceiverChanged is a free log subscription operation binding the contract event 0x647672599d3468abcfa241a13c9e3d34383caadb5cc80fb67c3cdfcd5f786059.
//
// Solidity: e FeeReceiverChanged(feeReceiver indexed address)
func (_BridgeFee *BridgeFeeFilterer) WatchFeeReceiverChanged(opts *bind.WatchOpts, sink chan<- *BridgeFeeFeeReceiverChanged, feeReceiver []common.Address) (event.Subscription, error) {

	var feeReceiverRule []interface{}
	for _, feeReceiverItem := range feeReceiver {
		feeReceiverRule = append(feeReceiverRule, feeReceiverItem)
	}

	logs, sub, err := _BridgeFee.contract.WatchLogs(opts, "FeeReceiverChanged", feeReceiverRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BridgeFeeFeeReceiverChanged)
				if err := _BridgeFee.contract.UnpackLog(event, "FeeReceiverChanged", log); err != nil {
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

// BridgeFeeKLAYFeeChangedIterator is returned from FilterKLAYFeeChanged and is used to iterate over the raw logs and unpacked data for KLAYFeeChanged events raised by the BridgeFee contract.
type BridgeFeeKLAYFeeChangedIterator struct {
	Event *BridgeFeeKLAYFeeChanged // Event containing the contract specifics and raw log

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
func (it *BridgeFeeKLAYFeeChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BridgeFeeKLAYFeeChanged)
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
		it.Event = new(BridgeFeeKLAYFeeChanged)
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
func (it *BridgeFeeKLAYFeeChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BridgeFeeKLAYFeeChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BridgeFeeKLAYFeeChanged represents a KLAYFeeChanged event raised by the BridgeFee contract.
type BridgeFeeKLAYFeeChanged struct {
	Fee *big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// FilterKLAYFeeChanged is a free log retrieval operation binding the contract event 0xa7a33d0996347e1aa55ca2206015b61b9534bdd881d59d59aa680e25eefac365.
//
// Solidity: e KLAYFeeChanged(fee indexed uint256)
func (_BridgeFee *BridgeFeeFilterer) FilterKLAYFeeChanged(opts *bind.FilterOpts, fee []*big.Int) (*BridgeFeeKLAYFeeChangedIterator, error) {

	var feeRule []interface{}
	for _, feeItem := range fee {
		feeRule = append(feeRule, feeItem)
	}

	logs, sub, err := _BridgeFee.contract.FilterLogs(opts, "KLAYFeeChanged", feeRule)
	if err != nil {
		return nil, err
	}
	return &BridgeFeeKLAYFeeChangedIterator{contract: _BridgeFee.contract, event: "KLAYFeeChanged", logs: logs, sub: sub}, nil
}

// WatchKLAYFeeChanged is a free log subscription operation binding the contract event 0xa7a33d0996347e1aa55ca2206015b61b9534bdd881d59d59aa680e25eefac365.
//
// Solidity: e KLAYFeeChanged(fee indexed uint256)
func (_BridgeFee *BridgeFeeFilterer) WatchKLAYFeeChanged(opts *bind.WatchOpts, sink chan<- *BridgeFeeKLAYFeeChanged, fee []*big.Int) (event.Subscription, error) {

	var feeRule []interface{}
	for _, feeItem := range fee {
		feeRule = append(feeRule, feeItem)
	}

	logs, sub, err := _BridgeFee.contract.WatchLogs(opts, "KLAYFeeChanged", feeRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BridgeFeeKLAYFeeChanged)
				if err := _BridgeFee.contract.UnpackLog(event, "KLAYFeeChanged", log); err != nil {
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

// BridgeOperatorABI is the input ABI used to generate the binding from.
const BridgeOperatorABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"operators\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_operator\",\"type\":\"address\"}],\"name\":\"registerOperator\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint8\"}],\"name\":\"operatorThresholds\",\"outputs\":[{\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"isOwner\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\"},{\"name\":\"\",\"type\":\"address\"}],\"name\":\"votes\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint64\"}],\"name\":\"closedValueTransferVotes\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"votesCounts\",\"outputs\":[{\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"configurationNonce\",\"outputs\":[{\"name\":\"\",\"type\":\"uint64\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_operator\",\"type\":\"address\"}],\"name\":\"deregisterOperator\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_voteType\",\"type\":\"uint8\"},{\"name\":\"_threshold\",\"type\":\"uint8\"}],\"name\":\"setOperatorThreshold\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"}]"

// BridgeOperatorBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const BridgeOperatorBinRuntime = `0x`

// BridgeOperatorBin is the compiled bytecode used for deploying new contracts.
const BridgeOperatorBin = `0x`

// DeployBridgeOperator deploys a new Klaytn contract, binding an instance of BridgeOperator to it.
func DeployBridgeOperator(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *BridgeOperator, error) {
	parsed, err := abi.JSON(strings.NewReader(BridgeOperatorABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(BridgeOperatorBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &BridgeOperator{BridgeOperatorCaller: BridgeOperatorCaller{contract: contract}, BridgeOperatorTransactor: BridgeOperatorTransactor{contract: contract}, BridgeOperatorFilterer: BridgeOperatorFilterer{contract: contract}}, nil
}

// BridgeOperator is an auto generated Go binding around a Klaytn contract.
type BridgeOperator struct {
	BridgeOperatorCaller     // Read-only binding to the contract
	BridgeOperatorTransactor // Write-only binding to the contract
	BridgeOperatorFilterer   // Log filterer for contract events
}

// BridgeOperatorCaller is an auto generated read-only Go binding around a Klaytn contract.
type BridgeOperatorCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BridgeOperatorTransactor is an auto generated write-only Go binding around a Klaytn contract.
type BridgeOperatorTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BridgeOperatorFilterer is an auto generated log filtering Go binding around a Klaytn contract events.
type BridgeOperatorFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BridgeOperatorSession is an auto generated Go binding around a Klaytn contract,
// with pre-set call and transact options.
type BridgeOperatorSession struct {
	Contract     *BridgeOperator   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// BridgeOperatorCallerSession is an auto generated read-only Go binding around a Klaytn contract,
// with pre-set call options.
type BridgeOperatorCallerSession struct {
	Contract *BridgeOperatorCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// BridgeOperatorTransactorSession is an auto generated write-only Go binding around a Klaytn contract,
// with pre-set transact options.
type BridgeOperatorTransactorSession struct {
	Contract     *BridgeOperatorTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// BridgeOperatorRaw is an auto generated low-level Go binding around a Klaytn contract.
type BridgeOperatorRaw struct {
	Contract *BridgeOperator // Generic contract binding to access the raw methods on
}

// BridgeOperatorCallerRaw is an auto generated low-level read-only Go binding around a Klaytn contract.
type BridgeOperatorCallerRaw struct {
	Contract *BridgeOperatorCaller // Generic read-only contract binding to access the raw methods on
}

// BridgeOperatorTransactorRaw is an auto generated low-level write-only Go binding around a Klaytn contract.
type BridgeOperatorTransactorRaw struct {
	Contract *BridgeOperatorTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBridgeOperator creates a new instance of BridgeOperator, bound to a specific deployed contract.
func NewBridgeOperator(address common.Address, backend bind.ContractBackend) (*BridgeOperator, error) {
	contract, err := bindBridgeOperator(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &BridgeOperator{BridgeOperatorCaller: BridgeOperatorCaller{contract: contract}, BridgeOperatorTransactor: BridgeOperatorTransactor{contract: contract}, BridgeOperatorFilterer: BridgeOperatorFilterer{contract: contract}}, nil
}

// NewBridgeOperatorCaller creates a new read-only instance of BridgeOperator, bound to a specific deployed contract.
func NewBridgeOperatorCaller(address common.Address, caller bind.ContractCaller) (*BridgeOperatorCaller, error) {
	contract, err := bindBridgeOperator(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BridgeOperatorCaller{contract: contract}, nil
}

// NewBridgeOperatorTransactor creates a new write-only instance of BridgeOperator, bound to a specific deployed contract.
func NewBridgeOperatorTransactor(address common.Address, transactor bind.ContractTransactor) (*BridgeOperatorTransactor, error) {
	contract, err := bindBridgeOperator(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BridgeOperatorTransactor{contract: contract}, nil
}

// NewBridgeOperatorFilterer creates a new log filterer instance of BridgeOperator, bound to a specific deployed contract.
func NewBridgeOperatorFilterer(address common.Address, filterer bind.ContractFilterer) (*BridgeOperatorFilterer, error) {
	contract, err := bindBridgeOperator(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BridgeOperatorFilterer{contract: contract}, nil
}

// bindBridgeOperator binds a generic wrapper to an already deployed contract.
func bindBridgeOperator(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(BridgeOperatorABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BridgeOperator *BridgeOperatorRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _BridgeOperator.Contract.BridgeOperatorCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BridgeOperator *BridgeOperatorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BridgeOperator.Contract.BridgeOperatorTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BridgeOperator *BridgeOperatorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BridgeOperator.Contract.BridgeOperatorTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BridgeOperator *BridgeOperatorCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _BridgeOperator.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BridgeOperator *BridgeOperatorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BridgeOperator.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BridgeOperator *BridgeOperatorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BridgeOperator.Contract.contract.Transact(opts, method, params...)
}

// ClosedValueTransferVotes is a free data retrieval call binding the contract method 0x9832c1d7.
//
// Solidity: function closedValueTransferVotes( uint64) constant returns(bool)
func (_BridgeOperator *BridgeOperatorCaller) ClosedValueTransferVotes(opts *bind.CallOpts, arg0 uint64) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _BridgeOperator.contract.Call(opts, out, "closedValueTransferVotes", arg0)
	return *ret0, err
}

// ClosedValueTransferVotes is a free data retrieval call binding the contract method 0x9832c1d7.
//
// Solidity: function closedValueTransferVotes( uint64) constant returns(bool)
func (_BridgeOperator *BridgeOperatorSession) ClosedValueTransferVotes(arg0 uint64) (bool, error) {
	return _BridgeOperator.Contract.ClosedValueTransferVotes(&_BridgeOperator.CallOpts, arg0)
}

// ClosedValueTransferVotes is a free data retrieval call binding the contract method 0x9832c1d7.
//
// Solidity: function closedValueTransferVotes( uint64) constant returns(bool)
func (_BridgeOperator *BridgeOperatorCallerSession) ClosedValueTransferVotes(arg0 uint64) (bool, error) {
	return _BridgeOperator.Contract.ClosedValueTransferVotes(&_BridgeOperator.CallOpts, arg0)
}

// ConfigurationNonce is a free data retrieval call binding the contract method 0xac6fff0b.
//
// Solidity: function configurationNonce() constant returns(uint64)
func (_BridgeOperator *BridgeOperatorCaller) ConfigurationNonce(opts *bind.CallOpts) (uint64, error) {
	var (
		ret0 = new(uint64)
	)
	out := ret0
	err := _BridgeOperator.contract.Call(opts, out, "configurationNonce")
	return *ret0, err
}

// ConfigurationNonce is a free data retrieval call binding the contract method 0xac6fff0b.
//
// Solidity: function configurationNonce() constant returns(uint64)
func (_BridgeOperator *BridgeOperatorSession) ConfigurationNonce() (uint64, error) {
	return _BridgeOperator.Contract.ConfigurationNonce(&_BridgeOperator.CallOpts)
}

// ConfigurationNonce is a free data retrieval call binding the contract method 0xac6fff0b.
//
// Solidity: function configurationNonce() constant returns(uint64)
func (_BridgeOperator *BridgeOperatorCallerSession) ConfigurationNonce() (uint64, error) {
	return _BridgeOperator.Contract.ConfigurationNonce(&_BridgeOperator.CallOpts)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() constant returns(bool)
func (_BridgeOperator *BridgeOperatorCaller) IsOwner(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _BridgeOperator.contract.Call(opts, out, "isOwner")
	return *ret0, err
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() constant returns(bool)
func (_BridgeOperator *BridgeOperatorSession) IsOwner() (bool, error) {
	return _BridgeOperator.Contract.IsOwner(&_BridgeOperator.CallOpts)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() constant returns(bool)
func (_BridgeOperator *BridgeOperatorCallerSession) IsOwner() (bool, error) {
	return _BridgeOperator.Contract.IsOwner(&_BridgeOperator.CallOpts)
}

// OperatorThresholds is a free data retrieval call binding the contract method 0x5526f76b.
//
// Solidity: function operatorThresholds( uint8) constant returns(uint8)
func (_BridgeOperator *BridgeOperatorCaller) OperatorThresholds(opts *bind.CallOpts, arg0 uint8) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _BridgeOperator.contract.Call(opts, out, "operatorThresholds", arg0)
	return *ret0, err
}

// OperatorThresholds is a free data retrieval call binding the contract method 0x5526f76b.
//
// Solidity: function operatorThresholds( uint8) constant returns(uint8)
func (_BridgeOperator *BridgeOperatorSession) OperatorThresholds(arg0 uint8) (uint8, error) {
	return _BridgeOperator.Contract.OperatorThresholds(&_BridgeOperator.CallOpts, arg0)
}

// OperatorThresholds is a free data retrieval call binding the contract method 0x5526f76b.
//
// Solidity: function operatorThresholds( uint8) constant returns(uint8)
func (_BridgeOperator *BridgeOperatorCallerSession) OperatorThresholds(arg0 uint8) (uint8, error) {
	return _BridgeOperator.Contract.OperatorThresholds(&_BridgeOperator.CallOpts, arg0)
}

// Operators is a free data retrieval call binding the contract method 0x13e7c9d8.
//
// Solidity: function operators( address) constant returns(bool)
func (_BridgeOperator *BridgeOperatorCaller) Operators(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _BridgeOperator.contract.Call(opts, out, "operators", arg0)
	return *ret0, err
}

// Operators is a free data retrieval call binding the contract method 0x13e7c9d8.
//
// Solidity: function operators( address) constant returns(bool)
func (_BridgeOperator *BridgeOperatorSession) Operators(arg0 common.Address) (bool, error) {
	return _BridgeOperator.Contract.Operators(&_BridgeOperator.CallOpts, arg0)
}

// Operators is a free data retrieval call binding the contract method 0x13e7c9d8.
//
// Solidity: function operators( address) constant returns(bool)
func (_BridgeOperator *BridgeOperatorCallerSession) Operators(arg0 common.Address) (bool, error) {
	return _BridgeOperator.Contract.Operators(&_BridgeOperator.CallOpts, arg0)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_BridgeOperator *BridgeOperatorCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _BridgeOperator.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_BridgeOperator *BridgeOperatorSession) Owner() (common.Address, error) {
	return _BridgeOperator.Contract.Owner(&_BridgeOperator.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_BridgeOperator *BridgeOperatorCallerSession) Owner() (common.Address, error) {
	return _BridgeOperator.Contract.Owner(&_BridgeOperator.CallOpts)
}

// Votes is a free data retrieval call binding the contract method 0x9386775a.
//
// Solidity: function votes( bytes32,  address) constant returns(bool)
func (_BridgeOperator *BridgeOperatorCaller) Votes(opts *bind.CallOpts, arg0 [32]byte, arg1 common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _BridgeOperator.contract.Call(opts, out, "votes", arg0, arg1)
	return *ret0, err
}

// Votes is a free data retrieval call binding the contract method 0x9386775a.
//
// Solidity: function votes( bytes32,  address) constant returns(bool)
func (_BridgeOperator *BridgeOperatorSession) Votes(arg0 [32]byte, arg1 common.Address) (bool, error) {
	return _BridgeOperator.Contract.Votes(&_BridgeOperator.CallOpts, arg0, arg1)
}

// Votes is a free data retrieval call binding the contract method 0x9386775a.
//
// Solidity: function votes( bytes32,  address) constant returns(bool)
func (_BridgeOperator *BridgeOperatorCallerSession) Votes(arg0 [32]byte, arg1 common.Address) (bool, error) {
	return _BridgeOperator.Contract.Votes(&_BridgeOperator.CallOpts, arg0, arg1)
}

// VotesCounts is a free data retrieval call binding the contract method 0xa4155e5d.
//
// Solidity: function votesCounts( bytes32) constant returns(uint8)
func (_BridgeOperator *BridgeOperatorCaller) VotesCounts(opts *bind.CallOpts, arg0 [32]byte) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _BridgeOperator.contract.Call(opts, out, "votesCounts", arg0)
	return *ret0, err
}

// VotesCounts is a free data retrieval call binding the contract method 0xa4155e5d.
//
// Solidity: function votesCounts( bytes32) constant returns(uint8)
func (_BridgeOperator *BridgeOperatorSession) VotesCounts(arg0 [32]byte) (uint8, error) {
	return _BridgeOperator.Contract.VotesCounts(&_BridgeOperator.CallOpts, arg0)
}

// VotesCounts is a free data retrieval call binding the contract method 0xa4155e5d.
//
// Solidity: function votesCounts( bytes32) constant returns(uint8)
func (_BridgeOperator *BridgeOperatorCallerSession) VotesCounts(arg0 [32]byte) (uint8, error) {
	return _BridgeOperator.Contract.VotesCounts(&_BridgeOperator.CallOpts, arg0)
}

// DeregisterOperator is a paid mutator transaction binding the contract method 0xd8cf98ca.
//
// Solidity: function deregisterOperator(_operator address) returns()
func (_BridgeOperator *BridgeOperatorTransactor) DeregisterOperator(opts *bind.TransactOpts, _operator common.Address) (*types.Transaction, error) {
	return _BridgeOperator.contract.Transact(opts, "deregisterOperator", _operator)
}

// DeregisterOperator is a paid mutator transaction binding the contract method 0xd8cf98ca.
//
// Solidity: function deregisterOperator(_operator address) returns()
func (_BridgeOperator *BridgeOperatorSession) DeregisterOperator(_operator common.Address) (*types.Transaction, error) {
	return _BridgeOperator.Contract.DeregisterOperator(&_BridgeOperator.TransactOpts, _operator)
}

// DeregisterOperator is a paid mutator transaction binding the contract method 0xd8cf98ca.
//
// Solidity: function deregisterOperator(_operator address) returns()
func (_BridgeOperator *BridgeOperatorTransactorSession) DeregisterOperator(_operator common.Address) (*types.Transaction, error) {
	return _BridgeOperator.Contract.DeregisterOperator(&_BridgeOperator.TransactOpts, _operator)
}

// RegisterOperator is a paid mutator transaction binding the contract method 0x3682a450.
//
// Solidity: function registerOperator(_operator address) returns()
func (_BridgeOperator *BridgeOperatorTransactor) RegisterOperator(opts *bind.TransactOpts, _operator common.Address) (*types.Transaction, error) {
	return _BridgeOperator.contract.Transact(opts, "registerOperator", _operator)
}

// RegisterOperator is a paid mutator transaction binding the contract method 0x3682a450.
//
// Solidity: function registerOperator(_operator address) returns()
func (_BridgeOperator *BridgeOperatorSession) RegisterOperator(_operator common.Address) (*types.Transaction, error) {
	return _BridgeOperator.Contract.RegisterOperator(&_BridgeOperator.TransactOpts, _operator)
}

// RegisterOperator is a paid mutator transaction binding the contract method 0x3682a450.
//
// Solidity: function registerOperator(_operator address) returns()
func (_BridgeOperator *BridgeOperatorTransactorSession) RegisterOperator(_operator common.Address) (*types.Transaction, error) {
	return _BridgeOperator.Contract.RegisterOperator(&_BridgeOperator.TransactOpts, _operator)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_BridgeOperator *BridgeOperatorTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BridgeOperator.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_BridgeOperator *BridgeOperatorSession) RenounceOwnership() (*types.Transaction, error) {
	return _BridgeOperator.Contract.RenounceOwnership(&_BridgeOperator.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_BridgeOperator *BridgeOperatorTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _BridgeOperator.Contract.RenounceOwnership(&_BridgeOperator.TransactOpts)
}

// SetOperatorThreshold is a paid mutator transaction binding the contract method 0xee2aec65.
//
// Solidity: function setOperatorThreshold(_voteType uint8, _threshold uint8) returns()
func (_BridgeOperator *BridgeOperatorTransactor) SetOperatorThreshold(opts *bind.TransactOpts, _voteType uint8, _threshold uint8) (*types.Transaction, error) {
	return _BridgeOperator.contract.Transact(opts, "setOperatorThreshold", _voteType, _threshold)
}

// SetOperatorThreshold is a paid mutator transaction binding the contract method 0xee2aec65.
//
// Solidity: function setOperatorThreshold(_voteType uint8, _threshold uint8) returns()
func (_BridgeOperator *BridgeOperatorSession) SetOperatorThreshold(_voteType uint8, _threshold uint8) (*types.Transaction, error) {
	return _BridgeOperator.Contract.SetOperatorThreshold(&_BridgeOperator.TransactOpts, _voteType, _threshold)
}

// SetOperatorThreshold is a paid mutator transaction binding the contract method 0xee2aec65.
//
// Solidity: function setOperatorThreshold(_voteType uint8, _threshold uint8) returns()
func (_BridgeOperator *BridgeOperatorTransactorSession) SetOperatorThreshold(_voteType uint8, _threshold uint8) (*types.Transaction, error) {
	return _BridgeOperator.Contract.SetOperatorThreshold(&_BridgeOperator.TransactOpts, _voteType, _threshold)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(newOwner address) returns()
func (_BridgeOperator *BridgeOperatorTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _BridgeOperator.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(newOwner address) returns()
func (_BridgeOperator *BridgeOperatorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _BridgeOperator.Contract.TransferOwnership(&_BridgeOperator.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(newOwner address) returns()
func (_BridgeOperator *BridgeOperatorTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _BridgeOperator.Contract.TransferOwnership(&_BridgeOperator.TransactOpts, newOwner)
}

// BridgeOperatorOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the BridgeOperator contract.
type BridgeOperatorOwnershipTransferredIterator struct {
	Event *BridgeOperatorOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *BridgeOperatorOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BridgeOperatorOwnershipTransferred)
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
		it.Event = new(BridgeOperatorOwnershipTransferred)
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
func (it *BridgeOperatorOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BridgeOperatorOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BridgeOperatorOwnershipTransferred represents a OwnershipTransferred event raised by the BridgeOperator contract.
type BridgeOperatorOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: e OwnershipTransferred(previousOwner indexed address, newOwner indexed address)
func (_BridgeOperator *BridgeOperatorFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*BridgeOperatorOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _BridgeOperator.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &BridgeOperatorOwnershipTransferredIterator{contract: _BridgeOperator.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: e OwnershipTransferred(previousOwner indexed address, newOwner indexed address)
func (_BridgeOperator *BridgeOperatorFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *BridgeOperatorOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _BridgeOperator.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BridgeOperatorOwnershipTransferred)
				if err := _BridgeOperator.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// CallbackABI is the input ABI used to generate the binding from.
const CallbackABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\"},{\"name\":\"_valueOrID\",\"type\":\"uint256\"},{\"name\":\"_tokenAddress\",\"type\":\"address\"},{\"name\":\"_price\",\"type\":\"uint256\"}],\"name\":\"RegisterOffer\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"valueOrID\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"tokenAddress\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"price\",\"type\":\"uint256\"}],\"name\":\"RegisteredOffer\",\"type\":\"event\"}]"

// CallbackBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const CallbackBinRuntime = `0x608060405260043610603e5763ffffffff7c010000000000000000000000000000000000000000000000000000000060003504166344b6d69d81146043575b600080fd5b348015604e57600080fd5b50607b73ffffffffffffffffffffffffffffffffffffffff6004358116906024359060443516606435607d565b005b6040805173ffffffffffffffffffffffffffffffffffffffff8087168252602082018690528416818301526060810183905290517f6e0b5117e49b57aaf37c635363b1b78a14ad521875ec99079d95bee2838cfeb89181900360800190a1505050505600a165627a7a7230582070034f78676f8df8fb23f76428a64ec5a8394ca95f631824ef66e7b4d3064ea50029`

// CallbackBin is the compiled bytecode used for deploying new contracts.
const CallbackBin = `0x608060405234801561001057600080fd5b5061010d806100206000396000f300608060405260043610603e5763ffffffff7c010000000000000000000000000000000000000000000000000000000060003504166344b6d69d81146043575b600080fd5b348015604e57600080fd5b50607b73ffffffffffffffffffffffffffffffffffffffff6004358116906024359060443516606435607d565b005b6040805173ffffffffffffffffffffffffffffffffffffffff8087168252602082018690528416818301526060810183905290517f6e0b5117e49b57aaf37c635363b1b78a14ad521875ec99079d95bee2838cfeb89181900360800190a1505050505600a165627a7a7230582070034f78676f8df8fb23f76428a64ec5a8394ca95f631824ef66e7b4d3064ea50029`

// DeployCallback deploys a new Klaytn contract, binding an instance of Callback to it.
func DeployCallback(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Callback, error) {
	parsed, err := abi.JSON(strings.NewReader(CallbackABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(CallbackBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Callback{CallbackCaller: CallbackCaller{contract: contract}, CallbackTransactor: CallbackTransactor{contract: contract}, CallbackFilterer: CallbackFilterer{contract: contract}}, nil
}

// Callback is an auto generated Go binding around a Klaytn contract.
type Callback struct {
	CallbackCaller     // Read-only binding to the contract
	CallbackTransactor // Write-only binding to the contract
	CallbackFilterer   // Log filterer for contract events
}

// CallbackCaller is an auto generated read-only Go binding around a Klaytn contract.
type CallbackCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CallbackTransactor is an auto generated write-only Go binding around a Klaytn contract.
type CallbackTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CallbackFilterer is an auto generated log filtering Go binding around a Klaytn contract events.
type CallbackFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CallbackSession is an auto generated Go binding around a Klaytn contract,
// with pre-set call and transact options.
type CallbackSession struct {
	Contract     *Callback         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// CallbackCallerSession is an auto generated read-only Go binding around a Klaytn contract,
// with pre-set call options.
type CallbackCallerSession struct {
	Contract *CallbackCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// CallbackTransactorSession is an auto generated write-only Go binding around a Klaytn contract,
// with pre-set transact options.
type CallbackTransactorSession struct {
	Contract     *CallbackTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// CallbackRaw is an auto generated low-level Go binding around a Klaytn contract.
type CallbackRaw struct {
	Contract *Callback // Generic contract binding to access the raw methods on
}

// CallbackCallerRaw is an auto generated low-level read-only Go binding around a Klaytn contract.
type CallbackCallerRaw struct {
	Contract *CallbackCaller // Generic read-only contract binding to access the raw methods on
}

// CallbackTransactorRaw is an auto generated low-level write-only Go binding around a Klaytn contract.
type CallbackTransactorRaw struct {
	Contract *CallbackTransactor // Generic write-only contract binding to access the raw methods on
}

// NewCallback creates a new instance of Callback, bound to a specific deployed contract.
func NewCallback(address common.Address, backend bind.ContractBackend) (*Callback, error) {
	contract, err := bindCallback(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Callback{CallbackCaller: CallbackCaller{contract: contract}, CallbackTransactor: CallbackTransactor{contract: contract}, CallbackFilterer: CallbackFilterer{contract: contract}}, nil
}

// NewCallbackCaller creates a new read-only instance of Callback, bound to a specific deployed contract.
func NewCallbackCaller(address common.Address, caller bind.ContractCaller) (*CallbackCaller, error) {
	contract, err := bindCallback(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &CallbackCaller{contract: contract}, nil
}

// NewCallbackTransactor creates a new write-only instance of Callback, bound to a specific deployed contract.
func NewCallbackTransactor(address common.Address, transactor bind.ContractTransactor) (*CallbackTransactor, error) {
	contract, err := bindCallback(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &CallbackTransactor{contract: contract}, nil
}

// NewCallbackFilterer creates a new log filterer instance of Callback, bound to a specific deployed contract.
func NewCallbackFilterer(address common.Address, filterer bind.ContractFilterer) (*CallbackFilterer, error) {
	contract, err := bindCallback(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &CallbackFilterer{contract: contract}, nil
}

// bindCallback binds a generic wrapper to an already deployed contract.
func bindCallback(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(CallbackABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Callback *CallbackRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Callback.Contract.CallbackCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Callback *CallbackRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Callback.Contract.CallbackTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Callback *CallbackRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Callback.Contract.CallbackTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Callback *CallbackCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Callback.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Callback *CallbackTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Callback.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Callback *CallbackTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Callback.Contract.contract.Transact(opts, method, params...)
}

// RegisterOffer is a paid mutator transaction binding the contract method 0x44b6d69d.
//
// Solidity: function RegisterOffer(_owner address, _valueOrID uint256, _tokenAddress address, _price uint256) returns()
func (_Callback *CallbackTransactor) RegisterOffer(opts *bind.TransactOpts, _owner common.Address, _valueOrID *big.Int, _tokenAddress common.Address, _price *big.Int) (*types.Transaction, error) {
	return _Callback.contract.Transact(opts, "RegisterOffer", _owner, _valueOrID, _tokenAddress, _price)
}

// RegisterOffer is a paid mutator transaction binding the contract method 0x44b6d69d.
//
// Solidity: function RegisterOffer(_owner address, _valueOrID uint256, _tokenAddress address, _price uint256) returns()
func (_Callback *CallbackSession) RegisterOffer(_owner common.Address, _valueOrID *big.Int, _tokenAddress common.Address, _price *big.Int) (*types.Transaction, error) {
	return _Callback.Contract.RegisterOffer(&_Callback.TransactOpts, _owner, _valueOrID, _tokenAddress, _price)
}

// RegisterOffer is a paid mutator transaction binding the contract method 0x44b6d69d.
//
// Solidity: function RegisterOffer(_owner address, _valueOrID uint256, _tokenAddress address, _price uint256) returns()
func (_Callback *CallbackTransactorSession) RegisterOffer(_owner common.Address, _valueOrID *big.Int, _tokenAddress common.Address, _price *big.Int) (*types.Transaction, error) {
	return _Callback.Contract.RegisterOffer(&_Callback.TransactOpts, _owner, _valueOrID, _tokenAddress, _price)
}

// CallbackRegisteredOfferIterator is returned from FilterRegisteredOffer and is used to iterate over the raw logs and unpacked data for RegisteredOffer events raised by the Callback contract.
type CallbackRegisteredOfferIterator struct {
	Event *CallbackRegisteredOffer // Event containing the contract specifics and raw log

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
func (it *CallbackRegisteredOfferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CallbackRegisteredOffer)
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
		it.Event = new(CallbackRegisteredOffer)
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
func (it *CallbackRegisteredOfferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CallbackRegisteredOfferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CallbackRegisteredOffer represents a RegisteredOffer event raised by the Callback contract.
type CallbackRegisteredOffer struct {
	Owner        common.Address
	ValueOrID    *big.Int
	TokenAddress common.Address
	Price        *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterRegisteredOffer is a free log retrieval operation binding the contract event 0x6e0b5117e49b57aaf37c635363b1b78a14ad521875ec99079d95bee2838cfeb8.
//
// Solidity: e RegisteredOffer(owner address, valueOrID uint256, tokenAddress address, price uint256)
func (_Callback *CallbackFilterer) FilterRegisteredOffer(opts *bind.FilterOpts) (*CallbackRegisteredOfferIterator, error) {

	logs, sub, err := _Callback.contract.FilterLogs(opts, "RegisteredOffer")
	if err != nil {
		return nil, err
	}
	return &CallbackRegisteredOfferIterator{contract: _Callback.contract, event: "RegisteredOffer", logs: logs, sub: sub}, nil
}

// WatchRegisteredOffer is a free log subscription operation binding the contract event 0x6e0b5117e49b57aaf37c635363b1b78a14ad521875ec99079d95bee2838cfeb8.
//
// Solidity: e RegisteredOffer(owner address, valueOrID uint256, tokenAddress address, price uint256)
func (_Callback *CallbackFilterer) WatchRegisteredOffer(opts *bind.WatchOpts, sink chan<- *CallbackRegisteredOffer) (event.Subscription, error) {

	logs, sub, err := _Callback.contract.WatchLogs(opts, "RegisteredOffer")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CallbackRegisteredOffer)
				if err := _Callback.contract.UnpackLog(event, "RegisteredOffer", log); err != nil {
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

// ERC165ABI is the input ABI used to generate the binding from.
const ERC165ABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"}]"

// ERC165BinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const ERC165BinRuntime = `0x`

// ERC165Bin is the compiled bytecode used for deploying new contracts.
const ERC165Bin = `0x`

// DeployERC165 deploys a new Klaytn contract, binding an instance of ERC165 to it.
func DeployERC165(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ERC165, error) {
	parsed, err := abi.JSON(strings.NewReader(ERC165ABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(ERC165Bin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ERC165{ERC165Caller: ERC165Caller{contract: contract}, ERC165Transactor: ERC165Transactor{contract: contract}, ERC165Filterer: ERC165Filterer{contract: contract}}, nil
}

// ERC165 is an auto generated Go binding around a Klaytn contract.
type ERC165 struct {
	ERC165Caller     // Read-only binding to the contract
	ERC165Transactor // Write-only binding to the contract
	ERC165Filterer   // Log filterer for contract events
}

// ERC165Caller is an auto generated read-only Go binding around a Klaytn contract.
type ERC165Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC165Transactor is an auto generated write-only Go binding around a Klaytn contract.
type ERC165Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC165Filterer is an auto generated log filtering Go binding around a Klaytn contract events.
type ERC165Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC165Session is an auto generated Go binding around a Klaytn contract,
// with pre-set call and transact options.
type ERC165Session struct {
	Contract     *ERC165           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ERC165CallerSession is an auto generated read-only Go binding around a Klaytn contract,
// with pre-set call options.
type ERC165CallerSession struct {
	Contract *ERC165Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// ERC165TransactorSession is an auto generated write-only Go binding around a Klaytn contract,
// with pre-set transact options.
type ERC165TransactorSession struct {
	Contract     *ERC165Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ERC165Raw is an auto generated low-level Go binding around a Klaytn contract.
type ERC165Raw struct {
	Contract *ERC165 // Generic contract binding to access the raw methods on
}

// ERC165CallerRaw is an auto generated low-level read-only Go binding around a Klaytn contract.
type ERC165CallerRaw struct {
	Contract *ERC165Caller // Generic read-only contract binding to access the raw methods on
}

// ERC165TransactorRaw is an auto generated low-level write-only Go binding around a Klaytn contract.
type ERC165TransactorRaw struct {
	Contract *ERC165Transactor // Generic write-only contract binding to access the raw methods on
}

// NewERC165 creates a new instance of ERC165, bound to a specific deployed contract.
func NewERC165(address common.Address, backend bind.ContractBackend) (*ERC165, error) {
	contract, err := bindERC165(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ERC165{ERC165Caller: ERC165Caller{contract: contract}, ERC165Transactor: ERC165Transactor{contract: contract}, ERC165Filterer: ERC165Filterer{contract: contract}}, nil
}

// NewERC165Caller creates a new read-only instance of ERC165, bound to a specific deployed contract.
func NewERC165Caller(address common.Address, caller bind.ContractCaller) (*ERC165Caller, error) {
	contract, err := bindERC165(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ERC165Caller{contract: contract}, nil
}

// NewERC165Transactor creates a new write-only instance of ERC165, bound to a specific deployed contract.
func NewERC165Transactor(address common.Address, transactor bind.ContractTransactor) (*ERC165Transactor, error) {
	contract, err := bindERC165(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ERC165Transactor{contract: contract}, nil
}

// NewERC165Filterer creates a new log filterer instance of ERC165, bound to a specific deployed contract.
func NewERC165Filterer(address common.Address, filterer bind.ContractFilterer) (*ERC165Filterer, error) {
	contract, err := bindERC165(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ERC165Filterer{contract: contract}, nil
}

// bindERC165 binds a generic wrapper to an already deployed contract.
func bindERC165(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ERC165ABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ERC165 *ERC165Raw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _ERC165.Contract.ERC165Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ERC165 *ERC165Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC165.Contract.ERC165Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ERC165 *ERC165Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ERC165.Contract.ERC165Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ERC165 *ERC165CallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _ERC165.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ERC165 *ERC165TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC165.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ERC165 *ERC165TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ERC165.Contract.contract.Transact(opts, method, params...)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(interfaceId bytes4) constant returns(bool)
func (_ERC165 *ERC165Caller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _ERC165.contract.Call(opts, out, "supportsInterface", interfaceId)
	return *ret0, err
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(interfaceId bytes4) constant returns(bool)
func (_ERC165 *ERC165Session) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _ERC165.Contract.SupportsInterface(&_ERC165.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(interfaceId bytes4) constant returns(bool)
func (_ERC165 *ERC165CallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _ERC165.Contract.SupportsInterface(&_ERC165.CallOpts, interfaceId)
}

// ERC20ABI is the input ABI used to generate the binding from.
const ERC20ABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"spender\",\"type\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"from\",\"type\":\"address\"},{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"spender\",\"type\":\"address\"},{\"name\":\"addedValue\",\"type\":\"uint256\"}],\"name\":\"increaseAllowance\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"spender\",\"type\":\"address\"},{\"name\":\"subtractedValue\",\"type\":\"uint256\"}],\"name\":\"decreaseAllowance\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"owner\",\"type\":\"address\"},{\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"}]"

// ERC20BinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const ERC20BinRuntime = `0x60806040526004361061008d5763ffffffff7c0100000000000000000000000000000000000000000000000000000000600035041663095ea7b3811461009257806318160ddd146100ca57806323b872dd146100f1578063395093511461011b57806370a082311461013f578063a457c2d714610160578063a9059cbb14610184578063dd62ed3e146101a8575b600080fd5b34801561009e57600080fd5b506100b6600160a060020a03600435166024356101cf565b604080519115158252519081900360200190f35b3480156100d657600080fd5b506100df61024d565b60408051918252519081900360200190f35b3480156100fd57600080fd5b506100b6600160a060020a0360043581169060243516604435610253565b34801561012757600080fd5b506100b6600160a060020a03600435166024356102f0565b34801561014b57600080fd5b506100df600160a060020a03600435166103a0565b34801561016c57600080fd5b506100b6600160a060020a03600435166024356103bb565b34801561019057600080fd5b506100b6600160a060020a0360043516602435610406565b3480156101b457600080fd5b506100df600160a060020a036004358116906024351661041c565b6000600160a060020a03831615156101e657600080fd5b336000818152600160209081526040808320600160a060020a03881680855290835292819020869055805186815290519293927f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925929181900390910190a350600192915050565b60025490565b600160a060020a038316600090815260016020908152604080832033845290915281205482111561028357600080fd5b600160a060020a03841660009081526001602090815260408083203384529091529020546102b7908363ffffffff61044716565b600160a060020a03851660009081526001602090815260408083203384529091529020556102e684848461045e565b5060019392505050565b6000600160a060020a038316151561030757600080fd5b336000908152600160209081526040808320600160a060020a038716845290915290205461033b908363ffffffff61055016565b336000818152600160209081526040808320600160a060020a0389168085529083529281902085905580519485525191937f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925929081900390910190a350600192915050565b600160a060020a031660009081526020819052604090205490565b6000600160a060020a03831615156103d257600080fd5b336000908152600160209081526040808320600160a060020a038716845290915290205461033b908363ffffffff61044716565b600061041333848461045e565b50600192915050565b600160a060020a03918216600090815260016020908152604080832093909416825291909152205490565b6000808383111561045757600080fd5b5050900390565b600160a060020a03831660009081526020819052604090205481111561048357600080fd5b600160a060020a038216151561049857600080fd5b600160a060020a0383166000908152602081905260409020546104c1908263ffffffff61044716565b600160a060020a0380851660009081526020819052604080822093909355908416815220546104f6908263ffffffff61055016565b600160a060020a038084166000818152602081815260409182902094909455805185815290519193928716927fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef92918290030190a3505050565b60008282018381101561056257600080fd5b93925050505600a165627a7a723058205d5177bab5d25649c5887b280a988a805eb604c9ea185163fdb6dc7f8b57f62b0029`

// ERC20Bin is the compiled bytecode used for deploying new contracts.
const ERC20Bin = `0x608060405234801561001057600080fd5b50610595806100206000396000f30060806040526004361061008d5763ffffffff7c0100000000000000000000000000000000000000000000000000000000600035041663095ea7b3811461009257806318160ddd146100ca57806323b872dd146100f1578063395093511461011b57806370a082311461013f578063a457c2d714610160578063a9059cbb14610184578063dd62ed3e146101a8575b600080fd5b34801561009e57600080fd5b506100b6600160a060020a03600435166024356101cf565b604080519115158252519081900360200190f35b3480156100d657600080fd5b506100df61024d565b60408051918252519081900360200190f35b3480156100fd57600080fd5b506100b6600160a060020a0360043581169060243516604435610253565b34801561012757600080fd5b506100b6600160a060020a03600435166024356102f0565b34801561014b57600080fd5b506100df600160a060020a03600435166103a0565b34801561016c57600080fd5b506100b6600160a060020a03600435166024356103bb565b34801561019057600080fd5b506100b6600160a060020a0360043516602435610406565b3480156101b457600080fd5b506100df600160a060020a036004358116906024351661041c565b6000600160a060020a03831615156101e657600080fd5b336000818152600160209081526040808320600160a060020a03881680855290835292819020869055805186815290519293927f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925929181900390910190a350600192915050565b60025490565b600160a060020a038316600090815260016020908152604080832033845290915281205482111561028357600080fd5b600160a060020a03841660009081526001602090815260408083203384529091529020546102b7908363ffffffff61044716565b600160a060020a03851660009081526001602090815260408083203384529091529020556102e684848461045e565b5060019392505050565b6000600160a060020a038316151561030757600080fd5b336000908152600160209081526040808320600160a060020a038716845290915290205461033b908363ffffffff61055016565b336000818152600160209081526040808320600160a060020a0389168085529083529281902085905580519485525191937f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925929081900390910190a350600192915050565b600160a060020a031660009081526020819052604090205490565b6000600160a060020a03831615156103d257600080fd5b336000908152600160209081526040808320600160a060020a038716845290915290205461033b908363ffffffff61044716565b600061041333848461045e565b50600192915050565b600160a060020a03918216600090815260016020908152604080832093909416825291909152205490565b6000808383111561045757600080fd5b5050900390565b600160a060020a03831660009081526020819052604090205481111561048357600080fd5b600160a060020a038216151561049857600080fd5b600160a060020a0383166000908152602081905260409020546104c1908263ffffffff61044716565b600160a060020a0380851660009081526020819052604080822093909355908416815220546104f6908263ffffffff61055016565b600160a060020a038084166000818152602081815260409182902094909455805185815290519193928716927fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef92918290030190a3505050565b60008282018381101561056257600080fd5b93925050505600a165627a7a723058205d5177bab5d25649c5887b280a988a805eb604c9ea185163fdb6dc7f8b57f62b0029`

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
// Solidity: function allowance(owner address, spender address) constant returns(uint256)
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
// Solidity: function allowance(owner address, spender address) constant returns(uint256)
func (_ERC20 *ERC20Session) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _ERC20.Contract.Allowance(&_ERC20.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(owner address, spender address) constant returns(uint256)
func (_ERC20 *ERC20CallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _ERC20.Contract.Allowance(&_ERC20.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(owner address) constant returns(uint256)
func (_ERC20 *ERC20Caller) BalanceOf(opts *bind.CallOpts, owner common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ERC20.contract.Call(opts, out, "balanceOf", owner)
	return *ret0, err
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(owner address) constant returns(uint256)
func (_ERC20 *ERC20Session) BalanceOf(owner common.Address) (*big.Int, error) {
	return _ERC20.Contract.BalanceOf(&_ERC20.CallOpts, owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(owner address) constant returns(uint256)
func (_ERC20 *ERC20CallerSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _ERC20.Contract.BalanceOf(&_ERC20.CallOpts, owner)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
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
// Solidity: function totalSupply() constant returns(uint256)
func (_ERC20 *ERC20Session) TotalSupply() (*big.Int, error) {
	return _ERC20.Contract.TotalSupply(&_ERC20.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_ERC20 *ERC20CallerSession) TotalSupply() (*big.Int, error) {
	return _ERC20.Contract.TotalSupply(&_ERC20.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(spender address, value uint256) returns(bool)
func (_ERC20 *ERC20Transactor) Approve(opts *bind.TransactOpts, spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _ERC20.contract.Transact(opts, "approve", spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(spender address, value uint256) returns(bool)
func (_ERC20 *ERC20Session) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _ERC20.Contract.Approve(&_ERC20.TransactOpts, spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(spender address, value uint256) returns(bool)
func (_ERC20 *ERC20TransactorSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _ERC20.Contract.Approve(&_ERC20.TransactOpts, spender, value)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(spender address, subtractedValue uint256) returns(bool)
func (_ERC20 *ERC20Transactor) DecreaseAllowance(opts *bind.TransactOpts, spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _ERC20.contract.Transact(opts, "decreaseAllowance", spender, subtractedValue)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(spender address, subtractedValue uint256) returns(bool)
func (_ERC20 *ERC20Session) DecreaseAllowance(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _ERC20.Contract.DecreaseAllowance(&_ERC20.TransactOpts, spender, subtractedValue)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(spender address, subtractedValue uint256) returns(bool)
func (_ERC20 *ERC20TransactorSession) DecreaseAllowance(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _ERC20.Contract.DecreaseAllowance(&_ERC20.TransactOpts, spender, subtractedValue)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(spender address, addedValue uint256) returns(bool)
func (_ERC20 *ERC20Transactor) IncreaseAllowance(opts *bind.TransactOpts, spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _ERC20.contract.Transact(opts, "increaseAllowance", spender, addedValue)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(spender address, addedValue uint256) returns(bool)
func (_ERC20 *ERC20Session) IncreaseAllowance(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _ERC20.Contract.IncreaseAllowance(&_ERC20.TransactOpts, spender, addedValue)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(spender address, addedValue uint256) returns(bool)
func (_ERC20 *ERC20TransactorSession) IncreaseAllowance(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _ERC20.Contract.IncreaseAllowance(&_ERC20.TransactOpts, spender, addedValue)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(to address, value uint256) returns(bool)
func (_ERC20 *ERC20Transactor) Transfer(opts *bind.TransactOpts, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ERC20.contract.Transact(opts, "transfer", to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(to address, value uint256) returns(bool)
func (_ERC20 *ERC20Session) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ERC20.Contract.Transfer(&_ERC20.TransactOpts, to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(to address, value uint256) returns(bool)
func (_ERC20 *ERC20TransactorSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ERC20.Contract.Transfer(&_ERC20.TransactOpts, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(from address, to address, value uint256) returns(bool)
func (_ERC20 *ERC20Transactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ERC20.contract.Transact(opts, "transferFrom", from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(from address, to address, value uint256) returns(bool)
func (_ERC20 *ERC20Session) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ERC20.Contract.TransferFrom(&_ERC20.TransactOpts, from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(from address, to address, value uint256) returns(bool)
func (_ERC20 *ERC20TransactorSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ERC20.Contract.TransferFrom(&_ERC20.TransactOpts, from, to, value)
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
// Solidity: e Approval(owner indexed address, spender indexed address, value uint256)
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
// Solidity: e Approval(owner indexed address, spender indexed address, value uint256)
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
// Solidity: e Transfer(from indexed address, to indexed address, value uint256)
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
// Solidity: e Transfer(from indexed address, to indexed address, value uint256)
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

// ERC20BurnableABI is the input ABI used to generate the binding from.
const ERC20BurnableABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"spender\",\"type\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"from\",\"type\":\"address\"},{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"spender\",\"type\":\"address\"},{\"name\":\"addedValue\",\"type\":\"uint256\"}],\"name\":\"increaseAllowance\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"from\",\"type\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"burnFrom\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"spender\",\"type\":\"address\"},{\"name\":\"subtractedValue\",\"type\":\"uint256\"}],\"name\":\"decreaseAllowance\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"owner\",\"type\":\"address\"},{\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"}]"

// ERC20BurnableBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const ERC20BurnableBinRuntime = `0x6080604052600436106100a35763ffffffff7c0100000000000000000000000000000000000000000000000000000000600035041663095ea7b381146100a857806318160ddd146100e057806323b872dd14610107578063395093511461013157806342966c681461015557806370a082311461016f57806379cc679014610190578063a457c2d7146101b4578063a9059cbb146101d8578063dd62ed3e146101fc575b600080fd5b3480156100b457600080fd5b506100cc600160a060020a0360043516602435610223565b604080519115158252519081900360200190f35b3480156100ec57600080fd5b506100f56102a1565b60408051918252519081900360200190f35b34801561011357600080fd5b506100cc600160a060020a03600435811690602435166044356102a7565b34801561013d57600080fd5b506100cc600160a060020a0360043516602435610344565b34801561016157600080fd5b5061016d6004356103f4565b005b34801561017b57600080fd5b506100f5600160a060020a0360043516610401565b34801561019c57600080fd5b5061016d600160a060020a036004351660243561041c565b3480156101c057600080fd5b506100cc600160a060020a036004351660243561042a565b3480156101e457600080fd5b506100cc600160a060020a0360043516602435610475565b34801561020857600080fd5b506100f5600160a060020a036004358116906024351661048b565b6000600160a060020a038316151561023a57600080fd5b336000818152600160209081526040808320600160a060020a03881680855290835292819020869055805186815290519293927f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925929181900390910190a350600192915050565b60025490565b600160a060020a03831660009081526001602090815260408083203384529091528120548211156102d757600080fd5b600160a060020a038416600090815260016020908152604080832033845290915290205461030b908363ffffffff6104b616565b600160a060020a038516600090815260016020908152604080832033845290915290205561033a8484846104cd565b5060019392505050565b6000600160a060020a038316151561035b57600080fd5b336000908152600160209081526040808320600160a060020a038716845290915290205461038f908363ffffffff6105bf16565b336000818152600160209081526040808320600160a060020a0389168085529083529281902085905580519485525191937f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925929081900390910190a350600192915050565b6103fe33826105d8565b50565b600160a060020a031660009081526020819052604090205490565b61042682826106a6565b5050565b6000600160a060020a038316151561044157600080fd5b336000908152600160209081526040808320600160a060020a038716845290915290205461038f908363ffffffff6104b616565b60006104823384846104cd565b50600192915050565b600160a060020a03918216600090815260016020908152604080832093909416825291909152205490565b600080838311156104c657600080fd5b5050900390565b600160a060020a0383166000908152602081905260409020548111156104f257600080fd5b600160a060020a038216151561050757600080fd5b600160a060020a038316600090815260208190526040902054610530908263ffffffff6104b616565b600160a060020a038085166000908152602081905260408082209390935590841681522054610565908263ffffffff6105bf16565b600160a060020a038084166000818152602081815260409182902094909455805185815290519193928716927fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef92918290030190a3505050565b6000828201838110156105d157600080fd5b9392505050565b600160a060020a03821615156105ed57600080fd5b600160a060020a03821660009081526020819052604090205481111561061257600080fd5b600254610625908263ffffffff6104b616565b600255600160a060020a038216600090815260208190526040902054610651908263ffffffff6104b616565b600160a060020a038316600081815260208181526040808320949094558351858152935191937fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef929081900390910190a35050565b600160a060020a03821660009081526001602090815260408083203384529091529020548111156106d657600080fd5b600160a060020a038216600090815260016020908152604080832033845290915290205461070a908263ffffffff6104b616565b600160a060020a038316600090815260016020908152604080832033845290915290205561042682826105d85600a165627a7a7230582075e2221d8c9d72178ed81a57f32fb409a84557421f079e9d482f661c034cf95c0029`

// ERC20BurnableBin is the compiled bytecode used for deploying new contracts.
const ERC20BurnableBin = `0x608060405234801561001057600080fd5b50610764806100206000396000f3006080604052600436106100a35763ffffffff7c0100000000000000000000000000000000000000000000000000000000600035041663095ea7b381146100a857806318160ddd146100e057806323b872dd14610107578063395093511461013157806342966c681461015557806370a082311461016f57806379cc679014610190578063a457c2d7146101b4578063a9059cbb146101d8578063dd62ed3e146101fc575b600080fd5b3480156100b457600080fd5b506100cc600160a060020a0360043516602435610223565b604080519115158252519081900360200190f35b3480156100ec57600080fd5b506100f56102a1565b60408051918252519081900360200190f35b34801561011357600080fd5b506100cc600160a060020a03600435811690602435166044356102a7565b34801561013d57600080fd5b506100cc600160a060020a0360043516602435610344565b34801561016157600080fd5b5061016d6004356103f4565b005b34801561017b57600080fd5b506100f5600160a060020a0360043516610401565b34801561019c57600080fd5b5061016d600160a060020a036004351660243561041c565b3480156101c057600080fd5b506100cc600160a060020a036004351660243561042a565b3480156101e457600080fd5b506100cc600160a060020a0360043516602435610475565b34801561020857600080fd5b506100f5600160a060020a036004358116906024351661048b565b6000600160a060020a038316151561023a57600080fd5b336000818152600160209081526040808320600160a060020a03881680855290835292819020869055805186815290519293927f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925929181900390910190a350600192915050565b60025490565b600160a060020a03831660009081526001602090815260408083203384529091528120548211156102d757600080fd5b600160a060020a038416600090815260016020908152604080832033845290915290205461030b908363ffffffff6104b616565b600160a060020a038516600090815260016020908152604080832033845290915290205561033a8484846104cd565b5060019392505050565b6000600160a060020a038316151561035b57600080fd5b336000908152600160209081526040808320600160a060020a038716845290915290205461038f908363ffffffff6105bf16565b336000818152600160209081526040808320600160a060020a0389168085529083529281902085905580519485525191937f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925929081900390910190a350600192915050565b6103fe33826105d8565b50565b600160a060020a031660009081526020819052604090205490565b61042682826106a6565b5050565b6000600160a060020a038316151561044157600080fd5b336000908152600160209081526040808320600160a060020a038716845290915290205461038f908363ffffffff6104b616565b60006104823384846104cd565b50600192915050565b600160a060020a03918216600090815260016020908152604080832093909416825291909152205490565b600080838311156104c657600080fd5b5050900390565b600160a060020a0383166000908152602081905260409020548111156104f257600080fd5b600160a060020a038216151561050757600080fd5b600160a060020a038316600090815260208190526040902054610530908263ffffffff6104b616565b600160a060020a038085166000908152602081905260408082209390935590841681522054610565908263ffffffff6105bf16565b600160a060020a038084166000818152602081815260409182902094909455805185815290519193928716927fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef92918290030190a3505050565b6000828201838110156105d157600080fd5b9392505050565b600160a060020a03821615156105ed57600080fd5b600160a060020a03821660009081526020819052604090205481111561061257600080fd5b600254610625908263ffffffff6104b616565b600255600160a060020a038216600090815260208190526040902054610651908263ffffffff6104b616565b600160a060020a038316600081815260208181526040808320949094558351858152935191937fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef929081900390910190a35050565b600160a060020a03821660009081526001602090815260408083203384529091529020548111156106d657600080fd5b600160a060020a038216600090815260016020908152604080832033845290915290205461070a908263ffffffff6104b616565b600160a060020a038316600090815260016020908152604080832033845290915290205561042682826105d85600a165627a7a7230582075e2221d8c9d72178ed81a57f32fb409a84557421f079e9d482f661c034cf95c0029`

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
// Solidity: function allowance(owner address, spender address) constant returns(uint256)
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
// Solidity: function allowance(owner address, spender address) constant returns(uint256)
func (_ERC20Burnable *ERC20BurnableSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _ERC20Burnable.Contract.Allowance(&_ERC20Burnable.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(owner address, spender address) constant returns(uint256)
func (_ERC20Burnable *ERC20BurnableCallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _ERC20Burnable.Contract.Allowance(&_ERC20Burnable.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(owner address) constant returns(uint256)
func (_ERC20Burnable *ERC20BurnableCaller) BalanceOf(opts *bind.CallOpts, owner common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ERC20Burnable.contract.Call(opts, out, "balanceOf", owner)
	return *ret0, err
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(owner address) constant returns(uint256)
func (_ERC20Burnable *ERC20BurnableSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _ERC20Burnable.Contract.BalanceOf(&_ERC20Burnable.CallOpts, owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(owner address) constant returns(uint256)
func (_ERC20Burnable *ERC20BurnableCallerSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _ERC20Burnable.Contract.BalanceOf(&_ERC20Burnable.CallOpts, owner)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
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
// Solidity: function totalSupply() constant returns(uint256)
func (_ERC20Burnable *ERC20BurnableSession) TotalSupply() (*big.Int, error) {
	return _ERC20Burnable.Contract.TotalSupply(&_ERC20Burnable.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_ERC20Burnable *ERC20BurnableCallerSession) TotalSupply() (*big.Int, error) {
	return _ERC20Burnable.Contract.TotalSupply(&_ERC20Burnable.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(spender address, value uint256) returns(bool)
func (_ERC20Burnable *ERC20BurnableTransactor) Approve(opts *bind.TransactOpts, spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _ERC20Burnable.contract.Transact(opts, "approve", spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(spender address, value uint256) returns(bool)
func (_ERC20Burnable *ERC20BurnableSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _ERC20Burnable.Contract.Approve(&_ERC20Burnable.TransactOpts, spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(spender address, value uint256) returns(bool)
func (_ERC20Burnable *ERC20BurnableTransactorSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _ERC20Burnable.Contract.Approve(&_ERC20Burnable.TransactOpts, spender, value)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(value uint256) returns()
func (_ERC20Burnable *ERC20BurnableTransactor) Burn(opts *bind.TransactOpts, value *big.Int) (*types.Transaction, error) {
	return _ERC20Burnable.contract.Transact(opts, "burn", value)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(value uint256) returns()
func (_ERC20Burnable *ERC20BurnableSession) Burn(value *big.Int) (*types.Transaction, error) {
	return _ERC20Burnable.Contract.Burn(&_ERC20Burnable.TransactOpts, value)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(value uint256) returns()
func (_ERC20Burnable *ERC20BurnableTransactorSession) Burn(value *big.Int) (*types.Transaction, error) {
	return _ERC20Burnable.Contract.Burn(&_ERC20Burnable.TransactOpts, value)
}

// BurnFrom is a paid mutator transaction binding the contract method 0x79cc6790.
//
// Solidity: function burnFrom(from address, value uint256) returns()
func (_ERC20Burnable *ERC20BurnableTransactor) BurnFrom(opts *bind.TransactOpts, from common.Address, value *big.Int) (*types.Transaction, error) {
	return _ERC20Burnable.contract.Transact(opts, "burnFrom", from, value)
}

// BurnFrom is a paid mutator transaction binding the contract method 0x79cc6790.
//
// Solidity: function burnFrom(from address, value uint256) returns()
func (_ERC20Burnable *ERC20BurnableSession) BurnFrom(from common.Address, value *big.Int) (*types.Transaction, error) {
	return _ERC20Burnable.Contract.BurnFrom(&_ERC20Burnable.TransactOpts, from, value)
}

// BurnFrom is a paid mutator transaction binding the contract method 0x79cc6790.
//
// Solidity: function burnFrom(from address, value uint256) returns()
func (_ERC20Burnable *ERC20BurnableTransactorSession) BurnFrom(from common.Address, value *big.Int) (*types.Transaction, error) {
	return _ERC20Burnable.Contract.BurnFrom(&_ERC20Burnable.TransactOpts, from, value)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(spender address, subtractedValue uint256) returns(bool)
func (_ERC20Burnable *ERC20BurnableTransactor) DecreaseAllowance(opts *bind.TransactOpts, spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _ERC20Burnable.contract.Transact(opts, "decreaseAllowance", spender, subtractedValue)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(spender address, subtractedValue uint256) returns(bool)
func (_ERC20Burnable *ERC20BurnableSession) DecreaseAllowance(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _ERC20Burnable.Contract.DecreaseAllowance(&_ERC20Burnable.TransactOpts, spender, subtractedValue)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(spender address, subtractedValue uint256) returns(bool)
func (_ERC20Burnable *ERC20BurnableTransactorSession) DecreaseAllowance(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _ERC20Burnable.Contract.DecreaseAllowance(&_ERC20Burnable.TransactOpts, spender, subtractedValue)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(spender address, addedValue uint256) returns(bool)
func (_ERC20Burnable *ERC20BurnableTransactor) IncreaseAllowance(opts *bind.TransactOpts, spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _ERC20Burnable.contract.Transact(opts, "increaseAllowance", spender, addedValue)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(spender address, addedValue uint256) returns(bool)
func (_ERC20Burnable *ERC20BurnableSession) IncreaseAllowance(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _ERC20Burnable.Contract.IncreaseAllowance(&_ERC20Burnable.TransactOpts, spender, addedValue)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(spender address, addedValue uint256) returns(bool)
func (_ERC20Burnable *ERC20BurnableTransactorSession) IncreaseAllowance(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _ERC20Burnable.Contract.IncreaseAllowance(&_ERC20Burnable.TransactOpts, spender, addedValue)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(to address, value uint256) returns(bool)
func (_ERC20Burnable *ERC20BurnableTransactor) Transfer(opts *bind.TransactOpts, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ERC20Burnable.contract.Transact(opts, "transfer", to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(to address, value uint256) returns(bool)
func (_ERC20Burnable *ERC20BurnableSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ERC20Burnable.Contract.Transfer(&_ERC20Burnable.TransactOpts, to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(to address, value uint256) returns(bool)
func (_ERC20Burnable *ERC20BurnableTransactorSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ERC20Burnable.Contract.Transfer(&_ERC20Burnable.TransactOpts, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(from address, to address, value uint256) returns(bool)
func (_ERC20Burnable *ERC20BurnableTransactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ERC20Burnable.contract.Transact(opts, "transferFrom", from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(from address, to address, value uint256) returns(bool)
func (_ERC20Burnable *ERC20BurnableSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ERC20Burnable.Contract.TransferFrom(&_ERC20Burnable.TransactOpts, from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(from address, to address, value uint256) returns(bool)
func (_ERC20Burnable *ERC20BurnableTransactorSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ERC20Burnable.Contract.TransferFrom(&_ERC20Burnable.TransactOpts, from, to, value)
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
// Solidity: e Approval(owner indexed address, spender indexed address, value uint256)
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
// Solidity: e Approval(owner indexed address, spender indexed address, value uint256)
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
// Solidity: e Transfer(from indexed address, to indexed address, value uint256)
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
// Solidity: e Transfer(from indexed address, to indexed address, value uint256)
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

// ERC20MintableABI is the input ABI used to generate the binding from.
const ERC20MintableABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"spender\",\"type\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"from\",\"type\":\"address\"},{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"spender\",\"type\":\"address\"},{\"name\":\"addedValue\",\"type\":\"uint256\"}],\"name\":\"increaseAllowance\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"mint\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"account\",\"type\":\"address\"}],\"name\":\"addMinter\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"renounceMinter\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"spender\",\"type\":\"address\"},{\"name\":\"subtractedValue\",\"type\":\"uint256\"}],\"name\":\"decreaseAllowance\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"account\",\"type\":\"address\"}],\"name\":\"isMinter\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"owner\",\"type\":\"address\"},{\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"account\",\"type\":\"address\"}],\"name\":\"MinterAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"account\",\"type\":\"address\"}],\"name\":\"MinterRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"}]"

// ERC20MintableBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const ERC20MintableBinRuntime = `0x6080604052600436106100b95763ffffffff7c0100000000000000000000000000000000000000000000000000000000600035041663095ea7b381146100be57806318160ddd146100f657806323b872dd1461011d578063395093511461014757806340c10f191461016b57806370a082311461018f578063983b2d56146101b057806398650275146101d3578063a457c2d7146101e8578063a9059cbb1461020c578063aa271e1a14610230578063dd62ed3e14610251575b600080fd5b3480156100ca57600080fd5b506100e2600160a060020a0360043516602435610278565b604080519115158252519081900360200190f35b34801561010257600080fd5b5061010b6102f6565b60408051918252519081900360200190f35b34801561012957600080fd5b506100e2600160a060020a03600435811690602435166044356102fc565b34801561015357600080fd5b506100e2600160a060020a0360043516602435610399565b34801561017757600080fd5b506100e2600160a060020a0360043516602435610449565b34801561019b57600080fd5b5061010b600160a060020a0360043516610472565b3480156101bc57600080fd5b506101d1600160a060020a036004351661048d565b005b3480156101df57600080fd5b506101d16104ad565b3480156101f457600080fd5b506100e2600160a060020a03600435166024356104b8565b34801561021857600080fd5b506100e2600160a060020a0360043516602435610503565b34801561023c57600080fd5b506100e2600160a060020a0360043516610510565b34801561025d57600080fd5b5061010b600160a060020a0360043581169060243516610529565b6000600160a060020a038316151561028f57600080fd5b336000818152600160209081526040808320600160a060020a03881680855290835292819020869055805186815290519293927f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925929181900390910190a350600192915050565b60025490565b600160a060020a038316600090815260016020908152604080832033845290915281205482111561032c57600080fd5b600160a060020a0384166000908152600160209081526040808320338452909152902054610360908363ffffffff61055416565b600160a060020a038516600090815260016020908152604080832033845290915290205561038f84848461056b565b5060019392505050565b6000600160a060020a03831615156103b057600080fd5b336000908152600160209081526040808320600160a060020a03871684529091529020546103e4908363ffffffff61065d16565b336000818152600160209081526040808320600160a060020a0389168085529083529281902085905580519485525191937f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925929081900390910190a350600192915050565b600061045433610510565b151561045f57600080fd5b6104698383610676565b50600192915050565b600160a060020a031660009081526020819052604090205490565b61049633610510565b15156104a157600080fd5b6104aa81610720565b50565b6104b633610768565b565b6000600160a060020a03831615156104cf57600080fd5b336000908152600160209081526040808320600160a060020a03871684529091529020546103e4908363ffffffff61055416565b600061046933848461056b565b600061052360038363ffffffff6107b016565b92915050565b600160a060020a03918216600090815260016020908152604080832093909416825291909152205490565b6000808383111561056457600080fd5b5050900390565b600160a060020a03831660009081526020819052604090205481111561059057600080fd5b600160a060020a03821615156105a557600080fd5b600160a060020a0383166000908152602081905260409020546105ce908263ffffffff61055416565b600160a060020a038085166000908152602081905260408082209390935590841681522054610603908263ffffffff61065d16565b600160a060020a038084166000818152602081815260409182902094909455805185815290519193928716927fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef92918290030190a3505050565b60008282018381101561066f57600080fd5b9392505050565b600160a060020a038216151561068b57600080fd5b60025461069e908263ffffffff61065d16565b600255600160a060020a0382166000908152602081905260409020546106ca908263ffffffff61065d16565b600160a060020a0383166000818152602081815260408083209490945583518581529351929391927fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef9281900390910190a35050565b61073160038263ffffffff6107e716565b604051600160a060020a038216907f6ae172837ea30b801fbfcdd4108aa1d5bf8ff775444fd70256b44e6bf3dfc3f690600090a250565b61077960038263ffffffff61083516565b604051600160a060020a038216907fe94479a9f7e1952cc78f2d6baab678adc1b772d936c6583def489e524cb6669290600090a250565b6000600160a060020a03821615156107c757600080fd5b50600160a060020a03166000908152602091909152604090205460ff1690565b600160a060020a03811615156107fc57600080fd5b61080682826107b0565b1561081057600080fd5b600160a060020a0316600090815260209190915260409020805460ff19166001179055565b600160a060020a038116151561084a57600080fd5b61085482826107b0565b151561085f57600080fd5b600160a060020a0316600090815260209190915260409020805460ff191690555600a165627a7a72305820541b78840c55fb5ccc36d17aa05428c6dd79ac99025ddae432990c50c1d8900f0029`

// ERC20MintableBin is the compiled bytecode used for deploying new contracts.
const ERC20MintableBin = `0x60806040526100163364010000000061001b810204565b6100f8565b6100336003826401000000006107e761006a82021704565b604051600160a060020a038216907f6ae172837ea30b801fbfcdd4108aa1d5bf8ff775444fd70256b44e6bf3dfc3f690600090a250565b600160a060020a038116151561007f57600080fd5b61009282826401000000006100c1810204565b1561009c57600080fd5b600160a060020a0316600090815260209190915260409020805460ff19166001179055565b6000600160a060020a03821615156100d857600080fd5b50600160a060020a03166000908152602091909152604090205460ff1690565b6108ad806101076000396000f3006080604052600436106100b95763ffffffff7c0100000000000000000000000000000000000000000000000000000000600035041663095ea7b381146100be57806318160ddd146100f657806323b872dd1461011d578063395093511461014757806340c10f191461016b57806370a082311461018f578063983b2d56146101b057806398650275146101d3578063a457c2d7146101e8578063a9059cbb1461020c578063aa271e1a14610230578063dd62ed3e14610251575b600080fd5b3480156100ca57600080fd5b506100e2600160a060020a0360043516602435610278565b604080519115158252519081900360200190f35b34801561010257600080fd5b5061010b6102f6565b60408051918252519081900360200190f35b34801561012957600080fd5b506100e2600160a060020a03600435811690602435166044356102fc565b34801561015357600080fd5b506100e2600160a060020a0360043516602435610399565b34801561017757600080fd5b506100e2600160a060020a0360043516602435610449565b34801561019b57600080fd5b5061010b600160a060020a0360043516610472565b3480156101bc57600080fd5b506101d1600160a060020a036004351661048d565b005b3480156101df57600080fd5b506101d16104ad565b3480156101f457600080fd5b506100e2600160a060020a03600435166024356104b8565b34801561021857600080fd5b506100e2600160a060020a0360043516602435610503565b34801561023c57600080fd5b506100e2600160a060020a0360043516610510565b34801561025d57600080fd5b5061010b600160a060020a0360043581169060243516610529565b6000600160a060020a038316151561028f57600080fd5b336000818152600160209081526040808320600160a060020a03881680855290835292819020869055805186815290519293927f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925929181900390910190a350600192915050565b60025490565b600160a060020a038316600090815260016020908152604080832033845290915281205482111561032c57600080fd5b600160a060020a0384166000908152600160209081526040808320338452909152902054610360908363ffffffff61055416565b600160a060020a038516600090815260016020908152604080832033845290915290205561038f84848461056b565b5060019392505050565b6000600160a060020a03831615156103b057600080fd5b336000908152600160209081526040808320600160a060020a03871684529091529020546103e4908363ffffffff61065d16565b336000818152600160209081526040808320600160a060020a0389168085529083529281902085905580519485525191937f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925929081900390910190a350600192915050565b600061045433610510565b151561045f57600080fd5b6104698383610676565b50600192915050565b600160a060020a031660009081526020819052604090205490565b61049633610510565b15156104a157600080fd5b6104aa81610720565b50565b6104b633610768565b565b6000600160a060020a03831615156104cf57600080fd5b336000908152600160209081526040808320600160a060020a03871684529091529020546103e4908363ffffffff61055416565b600061046933848461056b565b600061052360038363ffffffff6107b016565b92915050565b600160a060020a03918216600090815260016020908152604080832093909416825291909152205490565b6000808383111561056457600080fd5b5050900390565b600160a060020a03831660009081526020819052604090205481111561059057600080fd5b600160a060020a03821615156105a557600080fd5b600160a060020a0383166000908152602081905260409020546105ce908263ffffffff61055416565b600160a060020a038085166000908152602081905260408082209390935590841681522054610603908263ffffffff61065d16565b600160a060020a038084166000818152602081815260409182902094909455805185815290519193928716927fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef92918290030190a3505050565b60008282018381101561066f57600080fd5b9392505050565b600160a060020a038216151561068b57600080fd5b60025461069e908263ffffffff61065d16565b600255600160a060020a0382166000908152602081905260409020546106ca908263ffffffff61065d16565b600160a060020a0383166000818152602081815260408083209490945583518581529351929391927fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef9281900390910190a35050565b61073160038263ffffffff6107e716565b604051600160a060020a038216907f6ae172837ea30b801fbfcdd4108aa1d5bf8ff775444fd70256b44e6bf3dfc3f690600090a250565b61077960038263ffffffff61083516565b604051600160a060020a038216907fe94479a9f7e1952cc78f2d6baab678adc1b772d936c6583def489e524cb6669290600090a250565b6000600160a060020a03821615156107c757600080fd5b50600160a060020a03166000908152602091909152604090205460ff1690565b600160a060020a03811615156107fc57600080fd5b61080682826107b0565b1561081057600080fd5b600160a060020a0316600090815260209190915260409020805460ff19166001179055565b600160a060020a038116151561084a57600080fd5b61085482826107b0565b151561085f57600080fd5b600160a060020a0316600090815260209190915260409020805460ff191690555600a165627a7a72305820541b78840c55fb5ccc36d17aa05428c6dd79ac99025ddae432990c50c1d8900f0029`

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
// Solidity: function allowance(owner address, spender address) constant returns(uint256)
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
// Solidity: function allowance(owner address, spender address) constant returns(uint256)
func (_ERC20Mintable *ERC20MintableSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _ERC20Mintable.Contract.Allowance(&_ERC20Mintable.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(owner address, spender address) constant returns(uint256)
func (_ERC20Mintable *ERC20MintableCallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _ERC20Mintable.Contract.Allowance(&_ERC20Mintable.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(owner address) constant returns(uint256)
func (_ERC20Mintable *ERC20MintableCaller) BalanceOf(opts *bind.CallOpts, owner common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ERC20Mintable.contract.Call(opts, out, "balanceOf", owner)
	return *ret0, err
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(owner address) constant returns(uint256)
func (_ERC20Mintable *ERC20MintableSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _ERC20Mintable.Contract.BalanceOf(&_ERC20Mintable.CallOpts, owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(owner address) constant returns(uint256)
func (_ERC20Mintable *ERC20MintableCallerSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _ERC20Mintable.Contract.BalanceOf(&_ERC20Mintable.CallOpts, owner)
}

// IsMinter is a free data retrieval call binding the contract method 0xaa271e1a.
//
// Solidity: function isMinter(account address) constant returns(bool)
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
// Solidity: function isMinter(account address) constant returns(bool)
func (_ERC20Mintable *ERC20MintableSession) IsMinter(account common.Address) (bool, error) {
	return _ERC20Mintable.Contract.IsMinter(&_ERC20Mintable.CallOpts, account)
}

// IsMinter is a free data retrieval call binding the contract method 0xaa271e1a.
//
// Solidity: function isMinter(account address) constant returns(bool)
func (_ERC20Mintable *ERC20MintableCallerSession) IsMinter(account common.Address) (bool, error) {
	return _ERC20Mintable.Contract.IsMinter(&_ERC20Mintable.CallOpts, account)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
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
// Solidity: function totalSupply() constant returns(uint256)
func (_ERC20Mintable *ERC20MintableSession) TotalSupply() (*big.Int, error) {
	return _ERC20Mintable.Contract.TotalSupply(&_ERC20Mintable.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_ERC20Mintable *ERC20MintableCallerSession) TotalSupply() (*big.Int, error) {
	return _ERC20Mintable.Contract.TotalSupply(&_ERC20Mintable.CallOpts)
}

// AddMinter is a paid mutator transaction binding the contract method 0x983b2d56.
//
// Solidity: function addMinter(account address) returns()
func (_ERC20Mintable *ERC20MintableTransactor) AddMinter(opts *bind.TransactOpts, account common.Address) (*types.Transaction, error) {
	return _ERC20Mintable.contract.Transact(opts, "addMinter", account)
}

// AddMinter is a paid mutator transaction binding the contract method 0x983b2d56.
//
// Solidity: function addMinter(account address) returns()
func (_ERC20Mintable *ERC20MintableSession) AddMinter(account common.Address) (*types.Transaction, error) {
	return _ERC20Mintable.Contract.AddMinter(&_ERC20Mintable.TransactOpts, account)
}

// AddMinter is a paid mutator transaction binding the contract method 0x983b2d56.
//
// Solidity: function addMinter(account address) returns()
func (_ERC20Mintable *ERC20MintableTransactorSession) AddMinter(account common.Address) (*types.Transaction, error) {
	return _ERC20Mintable.Contract.AddMinter(&_ERC20Mintable.TransactOpts, account)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(spender address, value uint256) returns(bool)
func (_ERC20Mintable *ERC20MintableTransactor) Approve(opts *bind.TransactOpts, spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _ERC20Mintable.contract.Transact(opts, "approve", spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(spender address, value uint256) returns(bool)
func (_ERC20Mintable *ERC20MintableSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _ERC20Mintable.Contract.Approve(&_ERC20Mintable.TransactOpts, spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(spender address, value uint256) returns(bool)
func (_ERC20Mintable *ERC20MintableTransactorSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _ERC20Mintable.Contract.Approve(&_ERC20Mintable.TransactOpts, spender, value)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(spender address, subtractedValue uint256) returns(bool)
func (_ERC20Mintable *ERC20MintableTransactor) DecreaseAllowance(opts *bind.TransactOpts, spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _ERC20Mintable.contract.Transact(opts, "decreaseAllowance", spender, subtractedValue)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(spender address, subtractedValue uint256) returns(bool)
func (_ERC20Mintable *ERC20MintableSession) DecreaseAllowance(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _ERC20Mintable.Contract.DecreaseAllowance(&_ERC20Mintable.TransactOpts, spender, subtractedValue)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(spender address, subtractedValue uint256) returns(bool)
func (_ERC20Mintable *ERC20MintableTransactorSession) DecreaseAllowance(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _ERC20Mintable.Contract.DecreaseAllowance(&_ERC20Mintable.TransactOpts, spender, subtractedValue)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(spender address, addedValue uint256) returns(bool)
func (_ERC20Mintable *ERC20MintableTransactor) IncreaseAllowance(opts *bind.TransactOpts, spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _ERC20Mintable.contract.Transact(opts, "increaseAllowance", spender, addedValue)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(spender address, addedValue uint256) returns(bool)
func (_ERC20Mintable *ERC20MintableSession) IncreaseAllowance(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _ERC20Mintable.Contract.IncreaseAllowance(&_ERC20Mintable.TransactOpts, spender, addedValue)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(spender address, addedValue uint256) returns(bool)
func (_ERC20Mintable *ERC20MintableTransactorSession) IncreaseAllowance(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _ERC20Mintable.Contract.IncreaseAllowance(&_ERC20Mintable.TransactOpts, spender, addedValue)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(to address, value uint256) returns(bool)
func (_ERC20Mintable *ERC20MintableTransactor) Mint(opts *bind.TransactOpts, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ERC20Mintable.contract.Transact(opts, "mint", to, value)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(to address, value uint256) returns(bool)
func (_ERC20Mintable *ERC20MintableSession) Mint(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ERC20Mintable.Contract.Mint(&_ERC20Mintable.TransactOpts, to, value)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(to address, value uint256) returns(bool)
func (_ERC20Mintable *ERC20MintableTransactorSession) Mint(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ERC20Mintable.Contract.Mint(&_ERC20Mintable.TransactOpts, to, value)
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
// Solidity: function transfer(to address, value uint256) returns(bool)
func (_ERC20Mintable *ERC20MintableTransactor) Transfer(opts *bind.TransactOpts, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ERC20Mintable.contract.Transact(opts, "transfer", to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(to address, value uint256) returns(bool)
func (_ERC20Mintable *ERC20MintableSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ERC20Mintable.Contract.Transfer(&_ERC20Mintable.TransactOpts, to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(to address, value uint256) returns(bool)
func (_ERC20Mintable *ERC20MintableTransactorSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ERC20Mintable.Contract.Transfer(&_ERC20Mintable.TransactOpts, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(from address, to address, value uint256) returns(bool)
func (_ERC20Mintable *ERC20MintableTransactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ERC20Mintable.contract.Transact(opts, "transferFrom", from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(from address, to address, value uint256) returns(bool)
func (_ERC20Mintable *ERC20MintableSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ERC20Mintable.Contract.TransferFrom(&_ERC20Mintable.TransactOpts, from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(from address, to address, value uint256) returns(bool)
func (_ERC20Mintable *ERC20MintableTransactorSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ERC20Mintable.Contract.TransferFrom(&_ERC20Mintable.TransactOpts, from, to, value)
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
// Solidity: e Approval(owner indexed address, spender indexed address, value uint256)
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
// Solidity: e Approval(owner indexed address, spender indexed address, value uint256)
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
// Solidity: e MinterAdded(account indexed address)
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
// Solidity: e MinterAdded(account indexed address)
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
// Solidity: e MinterRemoved(account indexed address)
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
// Solidity: e MinterRemoved(account indexed address)
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
// Solidity: e Transfer(from indexed address, to indexed address, value uint256)
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
// Solidity: e Transfer(from indexed address, to indexed address, value uint256)
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

// ERC721ABI is the input ABI used to generate the binding from.
const ERC721ABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"getApproved\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"from\",\"type\":\"address\"},{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"from\",\"type\":\"address\"},{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"safeTransferFrom\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"ownerOf\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"setApprovalForAll\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"from\",\"type\":\"address\"},{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\"},{\"name\":\"_data\",\"type\":\"bytes\"}],\"name\":\"safeTransferFrom\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"owner\",\"type\":\"address\"},{\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"isApprovedForAll\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"approved\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"ApprovalForAll\",\"type\":\"event\"}]"

// ERC721BinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const ERC721BinRuntime = `0x6080604052600436106100a35763ffffffff7c010000000000000000000000000000000000000000000000000000000060003504166301ffc9a781146100a8578063081812fc146100f3578063095ea7b31461012757806323b872dd1461014d57806342842e0e146101775780636352211e146101a157806370a08231146101b9578063a22cb465146101ec578063b88d4fde14610212578063e985e9c514610281575b600080fd5b3480156100b457600080fd5b506100df7bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19600435166102a8565b604080519115158252519081900360200190f35b3480156100ff57600080fd5b5061010b6004356102dc565b60408051600160a060020a039092168252519081900360200190f35b34801561013357600080fd5b5061014b600160a060020a036004351660243561030e565b005b34801561015957600080fd5b5061014b600160a060020a03600435811690602435166044356103c4565b34801561018357600080fd5b5061014b600160a060020a0360043581169060243516604435610452565b3480156101ad57600080fd5b5061010b600435610473565b3480156101c557600080fd5b506101da600160a060020a036004351661049d565b60408051918252519081900360200190f35b3480156101f857600080fd5b5061014b600160a060020a036004351660243515156104d0565b34801561021e57600080fd5b50604080516020601f60643560048181013592830184900484028501840190955281845261014b94600160a060020a0381358116956024803590921695604435953695608494019181908401838280828437509497506105549650505050505050565b34801561028d57600080fd5b506100df600160a060020a036004358116906024351661057c565b7bffffffffffffffffffffffffffffffffffffffffffffffffffffffff191660009081526020819052604090205460ff1690565b60006102e7826105aa565b15156102f257600080fd5b50600090815260026020526040902054600160a060020a031690565b600061031982610473565b9050600160a060020a03838116908216141561033457600080fd5b33600160a060020a03821614806103505750610350813361057c565b151561035b57600080fd5b600082815260026020526040808220805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a0387811691821790925591518593918516917f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b92591a4505050565b6103ce33826105c7565b15156103d957600080fd5b600160a060020a03821615156103ee57600080fd5b6103f88382610626565b6104028382610697565b61040c828261072d565b8082600160a060020a031684600160a060020a03167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef60405160405180910390a4505050565b61046e8383836020604051908101604052806000815250610554565b505050565b600081815260016020526040812054600160a060020a031680151561049757600080fd5b92915050565b6000600160a060020a03821615156104b457600080fd5b50600160a060020a031660009081526003602052604090205490565b600160a060020a0382163314156104e657600080fd5b336000818152600460209081526040808320600160a060020a03871680855290835292819020805460ff1916861515908117909155815190815290519293927f17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31929181900390910190a35050565b61055f8484846103c4565b61056b848484846107bd565b151561057657600080fd5b50505050565b600160a060020a03918216600090815260046020908152604080832093909416825291909152205460ff1690565b600090815260016020526040902054600160a060020a0316151590565b6000806105d383610473565b905080600160a060020a031684600160a060020a0316148061060e575083600160a060020a0316610603846102dc565b600160a060020a0316145b8061061e575061061e818561057c565b949350505050565b81600160a060020a031661063982610473565b600160a060020a03161461064c57600080fd5b600081815260026020526040902054600160a060020a031615610693576000818152600260205260409020805473ffffffffffffffffffffffffffffffffffffffff191690555b5050565b81600160a060020a03166106aa82610473565b600160a060020a0316146106bd57600080fd5b600160a060020a0382166000908152600360205260409020546106e790600163ffffffff61093f16565b600160a060020a03909216600090815260036020908152604080832094909455918152600190915220805473ffffffffffffffffffffffffffffffffffffffff19169055565b600081815260016020526040902054600160a060020a03161561074f57600080fd5b6000818152600160208181526040808420805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a038816908117909155845260039091529091205461079d91610956565b600160a060020a0390921660009081526003602052604090209190915550565b6000806107d285600160a060020a031661096f565b15156107e15760019150610936565b6040517f150b7a020000000000000000000000000000000000000000000000000000000081523360048201818152600160a060020a03898116602485015260448401889052608060648501908152875160848601528751918a169463150b7a0294938c938b938b93909160a490910190602085019080838360005b8381101561087457818101518382015260200161085c565b50505050905090810190601f1680156108a15780820380516001836020036101000a031916815260200191505b5095505050505050602060405180830381600087803b1580156108c357600080fd5b505af11580156108d7573d6000803e3d6000fd5b505050506040513d60208110156108ed57600080fd5b50517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1981167f150b7a020000000000000000000000000000000000000000000000000000000014925090505b50949350505050565b6000808383111561094f57600080fd5b5050900390565b60008282018381101561096857600080fd5b9392505050565b6000903b11905600a165627a7a723058204cff0e10440e2cdf8989cae2108c0521158194186e30706c1350d26f33bc67810029`

// ERC721Bin is the compiled bytecode used for deploying new contracts.
const ERC721Bin = `0x608060405234801561001057600080fd5b506100437f01ffc9a70000000000000000000000000000000000000000000000000000000064010000000061007a810204565b6100757f80ac58cd0000000000000000000000000000000000000000000000000000000064010000000061007a810204565b6100e6565b7fffffffff0000000000000000000000000000000000000000000000000000000080821614156100a957600080fd5b7fffffffff00000000000000000000000000000000000000000000000000000000166000908152602081905260409020805460ff19166001179055565b6109a3806100f56000396000f3006080604052600436106100a35763ffffffff7c010000000000000000000000000000000000000000000000000000000060003504166301ffc9a781146100a8578063081812fc146100f3578063095ea7b31461012757806323b872dd1461014d57806342842e0e146101775780636352211e146101a157806370a08231146101b9578063a22cb465146101ec578063b88d4fde14610212578063e985e9c514610281575b600080fd5b3480156100b457600080fd5b506100df7bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19600435166102a8565b604080519115158252519081900360200190f35b3480156100ff57600080fd5b5061010b6004356102dc565b60408051600160a060020a039092168252519081900360200190f35b34801561013357600080fd5b5061014b600160a060020a036004351660243561030e565b005b34801561015957600080fd5b5061014b600160a060020a03600435811690602435166044356103c4565b34801561018357600080fd5b5061014b600160a060020a0360043581169060243516604435610452565b3480156101ad57600080fd5b5061010b600435610473565b3480156101c557600080fd5b506101da600160a060020a036004351661049d565b60408051918252519081900360200190f35b3480156101f857600080fd5b5061014b600160a060020a036004351660243515156104d0565b34801561021e57600080fd5b50604080516020601f60643560048181013592830184900484028501840190955281845261014b94600160a060020a0381358116956024803590921695604435953695608494019181908401838280828437509497506105549650505050505050565b34801561028d57600080fd5b506100df600160a060020a036004358116906024351661057c565b7bffffffffffffffffffffffffffffffffffffffffffffffffffffffff191660009081526020819052604090205460ff1690565b60006102e7826105aa565b15156102f257600080fd5b50600090815260026020526040902054600160a060020a031690565b600061031982610473565b9050600160a060020a03838116908216141561033457600080fd5b33600160a060020a03821614806103505750610350813361057c565b151561035b57600080fd5b600082815260026020526040808220805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a0387811691821790925591518593918516917f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b92591a4505050565b6103ce33826105c7565b15156103d957600080fd5b600160a060020a03821615156103ee57600080fd5b6103f88382610626565b6104028382610697565b61040c828261072d565b8082600160a060020a031684600160a060020a03167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef60405160405180910390a4505050565b61046e8383836020604051908101604052806000815250610554565b505050565b600081815260016020526040812054600160a060020a031680151561049757600080fd5b92915050565b6000600160a060020a03821615156104b457600080fd5b50600160a060020a031660009081526003602052604090205490565b600160a060020a0382163314156104e657600080fd5b336000818152600460209081526040808320600160a060020a03871680855290835292819020805460ff1916861515908117909155815190815290519293927f17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31929181900390910190a35050565b61055f8484846103c4565b61056b848484846107bd565b151561057657600080fd5b50505050565b600160a060020a03918216600090815260046020908152604080832093909416825291909152205460ff1690565b600090815260016020526040902054600160a060020a0316151590565b6000806105d383610473565b905080600160a060020a031684600160a060020a0316148061060e575083600160a060020a0316610603846102dc565b600160a060020a0316145b8061061e575061061e818561057c565b949350505050565b81600160a060020a031661063982610473565b600160a060020a03161461064c57600080fd5b600081815260026020526040902054600160a060020a031615610693576000818152600260205260409020805473ffffffffffffffffffffffffffffffffffffffff191690555b5050565b81600160a060020a03166106aa82610473565b600160a060020a0316146106bd57600080fd5b600160a060020a0382166000908152600360205260409020546106e790600163ffffffff61093f16565b600160a060020a03909216600090815260036020908152604080832094909455918152600190915220805473ffffffffffffffffffffffffffffffffffffffff19169055565b600081815260016020526040902054600160a060020a03161561074f57600080fd5b6000818152600160208181526040808420805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a038816908117909155845260039091529091205461079d91610956565b600160a060020a0390921660009081526003602052604090209190915550565b6000806107d285600160a060020a031661096f565b15156107e15760019150610936565b6040517f150b7a020000000000000000000000000000000000000000000000000000000081523360048201818152600160a060020a03898116602485015260448401889052608060648501908152875160848601528751918a169463150b7a0294938c938b938b93909160a490910190602085019080838360005b8381101561087457818101518382015260200161085c565b50505050905090810190601f1680156108a15780820380516001836020036101000a031916815260200191505b5095505050505050602060405180830381600087803b1580156108c357600080fd5b505af11580156108d7573d6000803e3d6000fd5b505050506040513d60208110156108ed57600080fd5b50517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1981167f150b7a020000000000000000000000000000000000000000000000000000000014925090505b50949350505050565b6000808383111561094f57600080fd5b5050900390565b60008282018381101561096857600080fd5b9392505050565b6000903b11905600a165627a7a723058204cff0e10440e2cdf8989cae2108c0521158194186e30706c1350d26f33bc67810029`

// DeployERC721 deploys a new Klaytn contract, binding an instance of ERC721 to it.
func DeployERC721(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ERC721, error) {
	parsed, err := abi.JSON(strings.NewReader(ERC721ABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(ERC721Bin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ERC721{ERC721Caller: ERC721Caller{contract: contract}, ERC721Transactor: ERC721Transactor{contract: contract}, ERC721Filterer: ERC721Filterer{contract: contract}}, nil
}

// ERC721 is an auto generated Go binding around a Klaytn contract.
type ERC721 struct {
	ERC721Caller     // Read-only binding to the contract
	ERC721Transactor // Write-only binding to the contract
	ERC721Filterer   // Log filterer for contract events
}

// ERC721Caller is an auto generated read-only Go binding around a Klaytn contract.
type ERC721Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC721Transactor is an auto generated write-only Go binding around a Klaytn contract.
type ERC721Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC721Filterer is an auto generated log filtering Go binding around a Klaytn contract events.
type ERC721Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC721Session is an auto generated Go binding around a Klaytn contract,
// with pre-set call and transact options.
type ERC721Session struct {
	Contract     *ERC721           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ERC721CallerSession is an auto generated read-only Go binding around a Klaytn contract,
// with pre-set call options.
type ERC721CallerSession struct {
	Contract *ERC721Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// ERC721TransactorSession is an auto generated write-only Go binding around a Klaytn contract,
// with pre-set transact options.
type ERC721TransactorSession struct {
	Contract     *ERC721Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ERC721Raw is an auto generated low-level Go binding around a Klaytn contract.
type ERC721Raw struct {
	Contract *ERC721 // Generic contract binding to access the raw methods on
}

// ERC721CallerRaw is an auto generated low-level read-only Go binding around a Klaytn contract.
type ERC721CallerRaw struct {
	Contract *ERC721Caller // Generic read-only contract binding to access the raw methods on
}

// ERC721TransactorRaw is an auto generated low-level write-only Go binding around a Klaytn contract.
type ERC721TransactorRaw struct {
	Contract *ERC721Transactor // Generic write-only contract binding to access the raw methods on
}

// NewERC721 creates a new instance of ERC721, bound to a specific deployed contract.
func NewERC721(address common.Address, backend bind.ContractBackend) (*ERC721, error) {
	contract, err := bindERC721(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ERC721{ERC721Caller: ERC721Caller{contract: contract}, ERC721Transactor: ERC721Transactor{contract: contract}, ERC721Filterer: ERC721Filterer{contract: contract}}, nil
}

// NewERC721Caller creates a new read-only instance of ERC721, bound to a specific deployed contract.
func NewERC721Caller(address common.Address, caller bind.ContractCaller) (*ERC721Caller, error) {
	contract, err := bindERC721(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ERC721Caller{contract: contract}, nil
}

// NewERC721Transactor creates a new write-only instance of ERC721, bound to a specific deployed contract.
func NewERC721Transactor(address common.Address, transactor bind.ContractTransactor) (*ERC721Transactor, error) {
	contract, err := bindERC721(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ERC721Transactor{contract: contract}, nil
}

// NewERC721Filterer creates a new log filterer instance of ERC721, bound to a specific deployed contract.
func NewERC721Filterer(address common.Address, filterer bind.ContractFilterer) (*ERC721Filterer, error) {
	contract, err := bindERC721(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ERC721Filterer{contract: contract}, nil
}

// bindERC721 binds a generic wrapper to an already deployed contract.
func bindERC721(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ERC721ABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ERC721 *ERC721Raw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _ERC721.Contract.ERC721Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ERC721 *ERC721Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC721.Contract.ERC721Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ERC721 *ERC721Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ERC721.Contract.ERC721Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ERC721 *ERC721CallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _ERC721.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ERC721 *ERC721TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC721.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ERC721 *ERC721TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ERC721.Contract.contract.Transact(opts, method, params...)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(owner address) constant returns(uint256)
func (_ERC721 *ERC721Caller) BalanceOf(opts *bind.CallOpts, owner common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ERC721.contract.Call(opts, out, "balanceOf", owner)
	return *ret0, err
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(owner address) constant returns(uint256)
func (_ERC721 *ERC721Session) BalanceOf(owner common.Address) (*big.Int, error) {
	return _ERC721.Contract.BalanceOf(&_ERC721.CallOpts, owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(owner address) constant returns(uint256)
func (_ERC721 *ERC721CallerSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _ERC721.Contract.BalanceOf(&_ERC721.CallOpts, owner)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(tokenId uint256) constant returns(address)
func (_ERC721 *ERC721Caller) GetApproved(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _ERC721.contract.Call(opts, out, "getApproved", tokenId)
	return *ret0, err
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(tokenId uint256) constant returns(address)
func (_ERC721 *ERC721Session) GetApproved(tokenId *big.Int) (common.Address, error) {
	return _ERC721.Contract.GetApproved(&_ERC721.CallOpts, tokenId)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(tokenId uint256) constant returns(address)
func (_ERC721 *ERC721CallerSession) GetApproved(tokenId *big.Int) (common.Address, error) {
	return _ERC721.Contract.GetApproved(&_ERC721.CallOpts, tokenId)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(owner address, operator address) constant returns(bool)
func (_ERC721 *ERC721Caller) IsApprovedForAll(opts *bind.CallOpts, owner common.Address, operator common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _ERC721.contract.Call(opts, out, "isApprovedForAll", owner, operator)
	return *ret0, err
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(owner address, operator address) constant returns(bool)
func (_ERC721 *ERC721Session) IsApprovedForAll(owner common.Address, operator common.Address) (bool, error) {
	return _ERC721.Contract.IsApprovedForAll(&_ERC721.CallOpts, owner, operator)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(owner address, operator address) constant returns(bool)
func (_ERC721 *ERC721CallerSession) IsApprovedForAll(owner common.Address, operator common.Address) (bool, error) {
	return _ERC721.Contract.IsApprovedForAll(&_ERC721.CallOpts, owner, operator)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(tokenId uint256) constant returns(address)
func (_ERC721 *ERC721Caller) OwnerOf(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _ERC721.contract.Call(opts, out, "ownerOf", tokenId)
	return *ret0, err
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(tokenId uint256) constant returns(address)
func (_ERC721 *ERC721Session) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _ERC721.Contract.OwnerOf(&_ERC721.CallOpts, tokenId)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(tokenId uint256) constant returns(address)
func (_ERC721 *ERC721CallerSession) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _ERC721.Contract.OwnerOf(&_ERC721.CallOpts, tokenId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(interfaceId bytes4) constant returns(bool)
func (_ERC721 *ERC721Caller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _ERC721.contract.Call(opts, out, "supportsInterface", interfaceId)
	return *ret0, err
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(interfaceId bytes4) constant returns(bool)
func (_ERC721 *ERC721Session) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _ERC721.Contract.SupportsInterface(&_ERC721.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(interfaceId bytes4) constant returns(bool)
func (_ERC721 *ERC721CallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _ERC721.Contract.SupportsInterface(&_ERC721.CallOpts, interfaceId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(to address, tokenId uint256) returns()
func (_ERC721 *ERC721Transactor) Approve(opts *bind.TransactOpts, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _ERC721.contract.Transact(opts, "approve", to, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(to address, tokenId uint256) returns()
func (_ERC721 *ERC721Session) Approve(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _ERC721.Contract.Approve(&_ERC721.TransactOpts, to, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(to address, tokenId uint256) returns()
func (_ERC721 *ERC721TransactorSession) Approve(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _ERC721.Contract.Approve(&_ERC721.TransactOpts, to, tokenId)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(from address, to address, tokenId uint256, _data bytes) returns()
func (_ERC721 *ERC721Transactor) SafeTransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int, _data []byte) (*types.Transaction, error) {
	return _ERC721.contract.Transact(opts, "safeTransferFrom", from, to, tokenId, _data)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(from address, to address, tokenId uint256, _data bytes) returns()
func (_ERC721 *ERC721Session) SafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int, _data []byte) (*types.Transaction, error) {
	return _ERC721.Contract.SafeTransferFrom(&_ERC721.TransactOpts, from, to, tokenId, _data)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(from address, to address, tokenId uint256, _data bytes) returns()
func (_ERC721 *ERC721TransactorSession) SafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int, _data []byte) (*types.Transaction, error) {
	return _ERC721.Contract.SafeTransferFrom(&_ERC721.TransactOpts, from, to, tokenId, _data)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(to address, approved bool) returns()
func (_ERC721 *ERC721Transactor) SetApprovalForAll(opts *bind.TransactOpts, to common.Address, approved bool) (*types.Transaction, error) {
	return _ERC721.contract.Transact(opts, "setApprovalForAll", to, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(to address, approved bool) returns()
func (_ERC721 *ERC721Session) SetApprovalForAll(to common.Address, approved bool) (*types.Transaction, error) {
	return _ERC721.Contract.SetApprovalForAll(&_ERC721.TransactOpts, to, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(to address, approved bool) returns()
func (_ERC721 *ERC721TransactorSession) SetApprovalForAll(to common.Address, approved bool) (*types.Transaction, error) {
	return _ERC721.Contract.SetApprovalForAll(&_ERC721.TransactOpts, to, approved)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(from address, to address, tokenId uint256) returns()
func (_ERC721 *ERC721Transactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _ERC721.contract.Transact(opts, "transferFrom", from, to, tokenId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(from address, to address, tokenId uint256) returns()
func (_ERC721 *ERC721Session) TransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _ERC721.Contract.TransferFrom(&_ERC721.TransactOpts, from, to, tokenId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(from address, to address, tokenId uint256) returns()
func (_ERC721 *ERC721TransactorSession) TransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _ERC721.Contract.TransferFrom(&_ERC721.TransactOpts, from, to, tokenId)
}

// ERC721ApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the ERC721 contract.
type ERC721ApprovalIterator struct {
	Event *ERC721Approval // Event containing the contract specifics and raw log

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
func (it *ERC721ApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC721Approval)
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
		it.Event = new(ERC721Approval)
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
func (it *ERC721ApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC721ApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC721Approval represents a Approval event raised by the ERC721 contract.
type ERC721Approval struct {
	Owner    common.Address
	Approved common.Address
	TokenId  *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: e Approval(owner indexed address, approved indexed address, tokenId indexed uint256)
func (_ERC721 *ERC721Filterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, approved []common.Address, tokenId []*big.Int) (*ERC721ApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var approvedRule []interface{}
	for _, approvedItem := range approved {
		approvedRule = append(approvedRule, approvedItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _ERC721.contract.FilterLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &ERC721ApprovalIterator{contract: _ERC721.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: e Approval(owner indexed address, approved indexed address, tokenId indexed uint256)
func (_ERC721 *ERC721Filterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *ERC721Approval, owner []common.Address, approved []common.Address, tokenId []*big.Int) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var approvedRule []interface{}
	for _, approvedItem := range approved {
		approvedRule = append(approvedRule, approvedItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _ERC721.contract.WatchLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC721Approval)
				if err := _ERC721.contract.UnpackLog(event, "Approval", log); err != nil {
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

// ERC721ApprovalForAllIterator is returned from FilterApprovalForAll and is used to iterate over the raw logs and unpacked data for ApprovalForAll events raised by the ERC721 contract.
type ERC721ApprovalForAllIterator struct {
	Event *ERC721ApprovalForAll // Event containing the contract specifics and raw log

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
func (it *ERC721ApprovalForAllIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC721ApprovalForAll)
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
		it.Event = new(ERC721ApprovalForAll)
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
func (it *ERC721ApprovalForAllIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC721ApprovalForAllIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC721ApprovalForAll represents a ApprovalForAll event raised by the ERC721 contract.
type ERC721ApprovalForAll struct {
	Owner    common.Address
	Operator common.Address
	Approved bool
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApprovalForAll is a free log retrieval operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: e ApprovalForAll(owner indexed address, operator indexed address, approved bool)
func (_ERC721 *ERC721Filterer) FilterApprovalForAll(opts *bind.FilterOpts, owner []common.Address, operator []common.Address) (*ERC721ApprovalForAllIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _ERC721.contract.FilterLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return &ERC721ApprovalForAllIterator{contract: _ERC721.contract, event: "ApprovalForAll", logs: logs, sub: sub}, nil
}

// WatchApprovalForAll is a free log subscription operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: e ApprovalForAll(owner indexed address, operator indexed address, approved bool)
func (_ERC721 *ERC721Filterer) WatchApprovalForAll(opts *bind.WatchOpts, sink chan<- *ERC721ApprovalForAll, owner []common.Address, operator []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _ERC721.contract.WatchLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC721ApprovalForAll)
				if err := _ERC721.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
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

// ERC721TransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the ERC721 contract.
type ERC721TransferIterator struct {
	Event *ERC721Transfer // Event containing the contract specifics and raw log

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
func (it *ERC721TransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC721Transfer)
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
		it.Event = new(ERC721Transfer)
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
func (it *ERC721TransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC721TransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC721Transfer represents a Transfer event raised by the ERC721 contract.
type ERC721Transfer struct {
	From    common.Address
	To      common.Address
	TokenId *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: e Transfer(from indexed address, to indexed address, tokenId indexed uint256)
func (_ERC721 *ERC721Filterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address, tokenId []*big.Int) (*ERC721TransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _ERC721.contract.FilterLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &ERC721TransferIterator{contract: _ERC721.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: e Transfer(from indexed address, to indexed address, tokenId indexed uint256)
func (_ERC721 *ERC721Filterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *ERC721Transfer, from []common.Address, to []common.Address, tokenId []*big.Int) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _ERC721.contract.WatchLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC721Transfer)
				if err := _ERC721.contract.UnpackLog(event, "Transfer", log); err != nil {
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

// ERC721BurnableABI is the input ABI used to generate the binding from.
const ERC721BurnableABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"getApproved\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"from\",\"type\":\"address\"},{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"from\",\"type\":\"address\"},{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"safeTransferFrom\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"ownerOf\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"setApprovalForAll\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"from\",\"type\":\"address\"},{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\"},{\"name\":\"_data\",\"type\":\"bytes\"}],\"name\":\"safeTransferFrom\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"owner\",\"type\":\"address\"},{\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"isApprovedForAll\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"approved\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"ApprovalForAll\",\"type\":\"event\"}]"

// ERC721BurnableBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const ERC721BurnableBinRuntime = `0x6080604052600436106100ae5763ffffffff7c010000000000000000000000000000000000000000000000000000000060003504166301ffc9a781146100b3578063081812fc146100fe578063095ea7b31461013257806323b872dd1461015857806342842e0e1461018257806342966c68146101ac5780636352211e146101c457806370a08231146101dc578063a22cb4651461020f578063b88d4fde14610235578063e985e9c5146102a4575b600080fd5b3480156100bf57600080fd5b506100ea7bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19600435166102cb565b604080519115158252519081900360200190f35b34801561010a57600080fd5b506101166004356102ff565b60408051600160a060020a039092168252519081900360200190f35b34801561013e57600080fd5b50610156600160a060020a0360043516602435610331565b005b34801561016457600080fd5b50610156600160a060020a03600435811690602435166044356103e7565b34801561018e57600080fd5b50610156600160a060020a0360043581169060243516604435610475565b3480156101b857600080fd5b50610156600435610496565b3480156101d057600080fd5b506101166004356104c0565b3480156101e857600080fd5b506101fd600160a060020a03600435166104ea565b60408051918252519081900360200190f35b34801561021b57600080fd5b50610156600160a060020a0360043516602435151561051d565b34801561024157600080fd5b50604080516020601f60643560048181013592830184900484028501840190955281845261015694600160a060020a0381358116956024803590921695604435953695608494019181908401838280828437509497506105a19650505050505050565b3480156102b057600080fd5b506100ea600160a060020a03600435811690602435166105c9565b7bffffffffffffffffffffffffffffffffffffffffffffffffffffffff191660009081526020819052604090205460ff1690565b600061030a826105f7565b151561031557600080fd5b50600090815260026020526040902054600160a060020a031690565b600061033c826104c0565b9050600160a060020a03838116908216141561035757600080fd5b33600160a060020a0382161480610373575061037381336105c9565b151561037e57600080fd5b600082815260026020526040808220805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a0387811691821790925591518593918516917f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b92591a4505050565b6103f13382610614565b15156103fc57600080fd5b600160a060020a038216151561041157600080fd5b61041b8382610673565b61042583826106e4565b61042f828261077a565b8082600160a060020a031684600160a060020a03167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef60405160405180910390a4505050565b61049183838360206040519081016040528060008152506105a1565b505050565b6104a03382610614565b15156104ab57600080fd5b6104bd6104b7826104c0565b8261080a565b50565b600081815260016020526040812054600160a060020a03168015156104e457600080fd5b92915050565b6000600160a060020a038216151561050157600080fd5b50600160a060020a031660009081526003602052604090205490565b600160a060020a03821633141561053357600080fd5b336000818152600460209081526040808320600160a060020a03871680855290835292819020805460ff1916861515908117909155815190815290519293927f17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31929181900390910190a35050565b6105ac8484846103e7565b6105b88484848461085a565b15156105c357600080fd5b50505050565b600160a060020a03918216600090815260046020908152604080832093909416825291909152205460ff1690565b600090815260016020526040902054600160a060020a0316151590565b600080610620836104c0565b905080600160a060020a031684600160a060020a0316148061065b575083600160a060020a0316610650846102ff565b600160a060020a0316145b8061066b575061066b81856105c9565b949350505050565b81600160a060020a0316610686826104c0565b600160a060020a03161461069957600080fd5b600081815260026020526040902054600160a060020a0316156106e0576000818152600260205260409020805473ffffffffffffffffffffffffffffffffffffffff191690555b5050565b81600160a060020a03166106f7826104c0565b600160a060020a03161461070a57600080fd5b600160a060020a03821660009081526003602052604090205461073490600163ffffffff6109dc16565b600160a060020a03909216600090815260036020908152604080832094909455918152600190915220805473ffffffffffffffffffffffffffffffffffffffff19169055565b600081815260016020526040902054600160a060020a03161561079c57600080fd5b6000818152600160208181526040808420805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a03881690811790915584526003909152909120546107ea916109f3565b600160a060020a0390921660009081526003602052604090209190915550565b6108148282610673565b61081e82826106e4565b6040518190600090600160a060020a038516907fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef908390a45050565b60008061086f85600160a060020a0316610a0c565b151561087e57600191506109d3565b6040517f150b7a020000000000000000000000000000000000000000000000000000000081523360048201818152600160a060020a03898116602485015260448401889052608060648501908152875160848601528751918a169463150b7a0294938c938b938b93909160a490910190602085019080838360005b838110156109115781810151838201526020016108f9565b50505050905090810190601f16801561093e5780820380516001836020036101000a031916815260200191505b5095505050505050602060405180830381600087803b15801561096057600080fd5b505af1158015610974573d6000803e3d6000fd5b505050506040513d602081101561098a57600080fd5b50517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1981167f150b7a020000000000000000000000000000000000000000000000000000000014925090505b50949350505050565b600080838311156109ec57600080fd5b5050900390565b600082820183811015610a0557600080fd5b9392505050565b6000903b11905600a165627a7a72305820a2496ef201c80238d6d8b6a84fb8e7d5fbaf02593e1d07723e531a7d4cf5f9f20029`

// ERC721BurnableBin is the compiled bytecode used for deploying new contracts.
const ERC721BurnableBin = `0x60806040526100367f01ffc9a70000000000000000000000000000000000000000000000000000000064010000000061006d810204565b6100687f80ac58cd0000000000000000000000000000000000000000000000000000000064010000000061006d810204565b6100d9565b7fffffffff00000000000000000000000000000000000000000000000000000000808216141561009c57600080fd5b7fffffffff00000000000000000000000000000000000000000000000000000000166000908152602081905260409020805460ff19166001179055565b610a40806100e86000396000f3006080604052600436106100ae5763ffffffff7c010000000000000000000000000000000000000000000000000000000060003504166301ffc9a781146100b3578063081812fc146100fe578063095ea7b31461013257806323b872dd1461015857806342842e0e1461018257806342966c68146101ac5780636352211e146101c457806370a08231146101dc578063a22cb4651461020f578063b88d4fde14610235578063e985e9c5146102a4575b600080fd5b3480156100bf57600080fd5b506100ea7bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19600435166102cb565b604080519115158252519081900360200190f35b34801561010a57600080fd5b506101166004356102ff565b60408051600160a060020a039092168252519081900360200190f35b34801561013e57600080fd5b50610156600160a060020a0360043516602435610331565b005b34801561016457600080fd5b50610156600160a060020a03600435811690602435166044356103e7565b34801561018e57600080fd5b50610156600160a060020a0360043581169060243516604435610475565b3480156101b857600080fd5b50610156600435610496565b3480156101d057600080fd5b506101166004356104c0565b3480156101e857600080fd5b506101fd600160a060020a03600435166104ea565b60408051918252519081900360200190f35b34801561021b57600080fd5b50610156600160a060020a0360043516602435151561051d565b34801561024157600080fd5b50604080516020601f60643560048181013592830184900484028501840190955281845261015694600160a060020a0381358116956024803590921695604435953695608494019181908401838280828437509497506105a19650505050505050565b3480156102b057600080fd5b506100ea600160a060020a03600435811690602435166105c9565b7bffffffffffffffffffffffffffffffffffffffffffffffffffffffff191660009081526020819052604090205460ff1690565b600061030a826105f7565b151561031557600080fd5b50600090815260026020526040902054600160a060020a031690565b600061033c826104c0565b9050600160a060020a03838116908216141561035757600080fd5b33600160a060020a0382161480610373575061037381336105c9565b151561037e57600080fd5b600082815260026020526040808220805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a0387811691821790925591518593918516917f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b92591a4505050565b6103f13382610614565b15156103fc57600080fd5b600160a060020a038216151561041157600080fd5b61041b8382610673565b61042583826106e4565b61042f828261077a565b8082600160a060020a031684600160a060020a03167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef60405160405180910390a4505050565b61049183838360206040519081016040528060008152506105a1565b505050565b6104a03382610614565b15156104ab57600080fd5b6104bd6104b7826104c0565b8261080a565b50565b600081815260016020526040812054600160a060020a03168015156104e457600080fd5b92915050565b6000600160a060020a038216151561050157600080fd5b50600160a060020a031660009081526003602052604090205490565b600160a060020a03821633141561053357600080fd5b336000818152600460209081526040808320600160a060020a03871680855290835292819020805460ff1916861515908117909155815190815290519293927f17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31929181900390910190a35050565b6105ac8484846103e7565b6105b88484848461085a565b15156105c357600080fd5b50505050565b600160a060020a03918216600090815260046020908152604080832093909416825291909152205460ff1690565b600090815260016020526040902054600160a060020a0316151590565b600080610620836104c0565b905080600160a060020a031684600160a060020a0316148061065b575083600160a060020a0316610650846102ff565b600160a060020a0316145b8061066b575061066b81856105c9565b949350505050565b81600160a060020a0316610686826104c0565b600160a060020a03161461069957600080fd5b600081815260026020526040902054600160a060020a0316156106e0576000818152600260205260409020805473ffffffffffffffffffffffffffffffffffffffff191690555b5050565b81600160a060020a03166106f7826104c0565b600160a060020a03161461070a57600080fd5b600160a060020a03821660009081526003602052604090205461073490600163ffffffff6109dc16565b600160a060020a03909216600090815260036020908152604080832094909455918152600190915220805473ffffffffffffffffffffffffffffffffffffffff19169055565b600081815260016020526040902054600160a060020a03161561079c57600080fd5b6000818152600160208181526040808420805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a03881690811790915584526003909152909120546107ea916109f3565b600160a060020a0390921660009081526003602052604090209190915550565b6108148282610673565b61081e82826106e4565b6040518190600090600160a060020a038516907fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef908390a45050565b60008061086f85600160a060020a0316610a0c565b151561087e57600191506109d3565b6040517f150b7a020000000000000000000000000000000000000000000000000000000081523360048201818152600160a060020a03898116602485015260448401889052608060648501908152875160848601528751918a169463150b7a0294938c938b938b93909160a490910190602085019080838360005b838110156109115781810151838201526020016108f9565b50505050905090810190601f16801561093e5780820380516001836020036101000a031916815260200191505b5095505050505050602060405180830381600087803b15801561096057600080fd5b505af1158015610974573d6000803e3d6000fd5b505050506040513d602081101561098a57600080fd5b50517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1981167f150b7a020000000000000000000000000000000000000000000000000000000014925090505b50949350505050565b600080838311156109ec57600080fd5b5050900390565b600082820183811015610a0557600080fd5b9392505050565b6000903b11905600a165627a7a72305820a2496ef201c80238d6d8b6a84fb8e7d5fbaf02593e1d07723e531a7d4cf5f9f20029`

// DeployERC721Burnable deploys a new Klaytn contract, binding an instance of ERC721Burnable to it.
func DeployERC721Burnable(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ERC721Burnable, error) {
	parsed, err := abi.JSON(strings.NewReader(ERC721BurnableABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(ERC721BurnableBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ERC721Burnable{ERC721BurnableCaller: ERC721BurnableCaller{contract: contract}, ERC721BurnableTransactor: ERC721BurnableTransactor{contract: contract}, ERC721BurnableFilterer: ERC721BurnableFilterer{contract: contract}}, nil
}

// ERC721Burnable is an auto generated Go binding around a Klaytn contract.
type ERC721Burnable struct {
	ERC721BurnableCaller     // Read-only binding to the contract
	ERC721BurnableTransactor // Write-only binding to the contract
	ERC721BurnableFilterer   // Log filterer for contract events
}

// ERC721BurnableCaller is an auto generated read-only Go binding around a Klaytn contract.
type ERC721BurnableCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC721BurnableTransactor is an auto generated write-only Go binding around a Klaytn contract.
type ERC721BurnableTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC721BurnableFilterer is an auto generated log filtering Go binding around a Klaytn contract events.
type ERC721BurnableFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC721BurnableSession is an auto generated Go binding around a Klaytn contract,
// with pre-set call and transact options.
type ERC721BurnableSession struct {
	Contract     *ERC721Burnable   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ERC721BurnableCallerSession is an auto generated read-only Go binding around a Klaytn contract,
// with pre-set call options.
type ERC721BurnableCallerSession struct {
	Contract *ERC721BurnableCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// ERC721BurnableTransactorSession is an auto generated write-only Go binding around a Klaytn contract,
// with pre-set transact options.
type ERC721BurnableTransactorSession struct {
	Contract     *ERC721BurnableTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// ERC721BurnableRaw is an auto generated low-level Go binding around a Klaytn contract.
type ERC721BurnableRaw struct {
	Contract *ERC721Burnable // Generic contract binding to access the raw methods on
}

// ERC721BurnableCallerRaw is an auto generated low-level read-only Go binding around a Klaytn contract.
type ERC721BurnableCallerRaw struct {
	Contract *ERC721BurnableCaller // Generic read-only contract binding to access the raw methods on
}

// ERC721BurnableTransactorRaw is an auto generated low-level write-only Go binding around a Klaytn contract.
type ERC721BurnableTransactorRaw struct {
	Contract *ERC721BurnableTransactor // Generic write-only contract binding to access the raw methods on
}

// NewERC721Burnable creates a new instance of ERC721Burnable, bound to a specific deployed contract.
func NewERC721Burnable(address common.Address, backend bind.ContractBackend) (*ERC721Burnable, error) {
	contract, err := bindERC721Burnable(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ERC721Burnable{ERC721BurnableCaller: ERC721BurnableCaller{contract: contract}, ERC721BurnableTransactor: ERC721BurnableTransactor{contract: contract}, ERC721BurnableFilterer: ERC721BurnableFilterer{contract: contract}}, nil
}

// NewERC721BurnableCaller creates a new read-only instance of ERC721Burnable, bound to a specific deployed contract.
func NewERC721BurnableCaller(address common.Address, caller bind.ContractCaller) (*ERC721BurnableCaller, error) {
	contract, err := bindERC721Burnable(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ERC721BurnableCaller{contract: contract}, nil
}

// NewERC721BurnableTransactor creates a new write-only instance of ERC721Burnable, bound to a specific deployed contract.
func NewERC721BurnableTransactor(address common.Address, transactor bind.ContractTransactor) (*ERC721BurnableTransactor, error) {
	contract, err := bindERC721Burnable(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ERC721BurnableTransactor{contract: contract}, nil
}

// NewERC721BurnableFilterer creates a new log filterer instance of ERC721Burnable, bound to a specific deployed contract.
func NewERC721BurnableFilterer(address common.Address, filterer bind.ContractFilterer) (*ERC721BurnableFilterer, error) {
	contract, err := bindERC721Burnable(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ERC721BurnableFilterer{contract: contract}, nil
}

// bindERC721Burnable binds a generic wrapper to an already deployed contract.
func bindERC721Burnable(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ERC721BurnableABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ERC721Burnable *ERC721BurnableRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _ERC721Burnable.Contract.ERC721BurnableCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ERC721Burnable *ERC721BurnableRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC721Burnable.Contract.ERC721BurnableTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ERC721Burnable *ERC721BurnableRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ERC721Burnable.Contract.ERC721BurnableTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ERC721Burnable *ERC721BurnableCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _ERC721Burnable.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ERC721Burnable *ERC721BurnableTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC721Burnable.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ERC721Burnable *ERC721BurnableTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ERC721Burnable.Contract.contract.Transact(opts, method, params...)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(owner address) constant returns(uint256)
func (_ERC721Burnable *ERC721BurnableCaller) BalanceOf(opts *bind.CallOpts, owner common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ERC721Burnable.contract.Call(opts, out, "balanceOf", owner)
	return *ret0, err
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(owner address) constant returns(uint256)
func (_ERC721Burnable *ERC721BurnableSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _ERC721Burnable.Contract.BalanceOf(&_ERC721Burnable.CallOpts, owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(owner address) constant returns(uint256)
func (_ERC721Burnable *ERC721BurnableCallerSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _ERC721Burnable.Contract.BalanceOf(&_ERC721Burnable.CallOpts, owner)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(tokenId uint256) constant returns(address)
func (_ERC721Burnable *ERC721BurnableCaller) GetApproved(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _ERC721Burnable.contract.Call(opts, out, "getApproved", tokenId)
	return *ret0, err
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(tokenId uint256) constant returns(address)
func (_ERC721Burnable *ERC721BurnableSession) GetApproved(tokenId *big.Int) (common.Address, error) {
	return _ERC721Burnable.Contract.GetApproved(&_ERC721Burnable.CallOpts, tokenId)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(tokenId uint256) constant returns(address)
func (_ERC721Burnable *ERC721BurnableCallerSession) GetApproved(tokenId *big.Int) (common.Address, error) {
	return _ERC721Burnable.Contract.GetApproved(&_ERC721Burnable.CallOpts, tokenId)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(owner address, operator address) constant returns(bool)
func (_ERC721Burnable *ERC721BurnableCaller) IsApprovedForAll(opts *bind.CallOpts, owner common.Address, operator common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _ERC721Burnable.contract.Call(opts, out, "isApprovedForAll", owner, operator)
	return *ret0, err
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(owner address, operator address) constant returns(bool)
func (_ERC721Burnable *ERC721BurnableSession) IsApprovedForAll(owner common.Address, operator common.Address) (bool, error) {
	return _ERC721Burnable.Contract.IsApprovedForAll(&_ERC721Burnable.CallOpts, owner, operator)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(owner address, operator address) constant returns(bool)
func (_ERC721Burnable *ERC721BurnableCallerSession) IsApprovedForAll(owner common.Address, operator common.Address) (bool, error) {
	return _ERC721Burnable.Contract.IsApprovedForAll(&_ERC721Burnable.CallOpts, owner, operator)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(tokenId uint256) constant returns(address)
func (_ERC721Burnable *ERC721BurnableCaller) OwnerOf(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _ERC721Burnable.contract.Call(opts, out, "ownerOf", tokenId)
	return *ret0, err
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(tokenId uint256) constant returns(address)
func (_ERC721Burnable *ERC721BurnableSession) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _ERC721Burnable.Contract.OwnerOf(&_ERC721Burnable.CallOpts, tokenId)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(tokenId uint256) constant returns(address)
func (_ERC721Burnable *ERC721BurnableCallerSession) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _ERC721Burnable.Contract.OwnerOf(&_ERC721Burnable.CallOpts, tokenId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(interfaceId bytes4) constant returns(bool)
func (_ERC721Burnable *ERC721BurnableCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _ERC721Burnable.contract.Call(opts, out, "supportsInterface", interfaceId)
	return *ret0, err
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(interfaceId bytes4) constant returns(bool)
func (_ERC721Burnable *ERC721BurnableSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _ERC721Burnable.Contract.SupportsInterface(&_ERC721Burnable.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(interfaceId bytes4) constant returns(bool)
func (_ERC721Burnable *ERC721BurnableCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _ERC721Burnable.Contract.SupportsInterface(&_ERC721Burnable.CallOpts, interfaceId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(to address, tokenId uint256) returns()
func (_ERC721Burnable *ERC721BurnableTransactor) Approve(opts *bind.TransactOpts, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _ERC721Burnable.contract.Transact(opts, "approve", to, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(to address, tokenId uint256) returns()
func (_ERC721Burnable *ERC721BurnableSession) Approve(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _ERC721Burnable.Contract.Approve(&_ERC721Burnable.TransactOpts, to, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(to address, tokenId uint256) returns()
func (_ERC721Burnable *ERC721BurnableTransactorSession) Approve(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _ERC721Burnable.Contract.Approve(&_ERC721Burnable.TransactOpts, to, tokenId)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(tokenId uint256) returns()
func (_ERC721Burnable *ERC721BurnableTransactor) Burn(opts *bind.TransactOpts, tokenId *big.Int) (*types.Transaction, error) {
	return _ERC721Burnable.contract.Transact(opts, "burn", tokenId)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(tokenId uint256) returns()
func (_ERC721Burnable *ERC721BurnableSession) Burn(tokenId *big.Int) (*types.Transaction, error) {
	return _ERC721Burnable.Contract.Burn(&_ERC721Burnable.TransactOpts, tokenId)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(tokenId uint256) returns()
func (_ERC721Burnable *ERC721BurnableTransactorSession) Burn(tokenId *big.Int) (*types.Transaction, error) {
	return _ERC721Burnable.Contract.Burn(&_ERC721Burnable.TransactOpts, tokenId)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(from address, to address, tokenId uint256, _data bytes) returns()
func (_ERC721Burnable *ERC721BurnableTransactor) SafeTransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int, _data []byte) (*types.Transaction, error) {
	return _ERC721Burnable.contract.Transact(opts, "safeTransferFrom", from, to, tokenId, _data)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(from address, to address, tokenId uint256, _data bytes) returns()
func (_ERC721Burnable *ERC721BurnableSession) SafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int, _data []byte) (*types.Transaction, error) {
	return _ERC721Burnable.Contract.SafeTransferFrom(&_ERC721Burnable.TransactOpts, from, to, tokenId, _data)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(from address, to address, tokenId uint256, _data bytes) returns()
func (_ERC721Burnable *ERC721BurnableTransactorSession) SafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int, _data []byte) (*types.Transaction, error) {
	return _ERC721Burnable.Contract.SafeTransferFrom(&_ERC721Burnable.TransactOpts, from, to, tokenId, _data)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(to address, approved bool) returns()
func (_ERC721Burnable *ERC721BurnableTransactor) SetApprovalForAll(opts *bind.TransactOpts, to common.Address, approved bool) (*types.Transaction, error) {
	return _ERC721Burnable.contract.Transact(opts, "setApprovalForAll", to, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(to address, approved bool) returns()
func (_ERC721Burnable *ERC721BurnableSession) SetApprovalForAll(to common.Address, approved bool) (*types.Transaction, error) {
	return _ERC721Burnable.Contract.SetApprovalForAll(&_ERC721Burnable.TransactOpts, to, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(to address, approved bool) returns()
func (_ERC721Burnable *ERC721BurnableTransactorSession) SetApprovalForAll(to common.Address, approved bool) (*types.Transaction, error) {
	return _ERC721Burnable.Contract.SetApprovalForAll(&_ERC721Burnable.TransactOpts, to, approved)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(from address, to address, tokenId uint256) returns()
func (_ERC721Burnable *ERC721BurnableTransactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _ERC721Burnable.contract.Transact(opts, "transferFrom", from, to, tokenId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(from address, to address, tokenId uint256) returns()
func (_ERC721Burnable *ERC721BurnableSession) TransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _ERC721Burnable.Contract.TransferFrom(&_ERC721Burnable.TransactOpts, from, to, tokenId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(from address, to address, tokenId uint256) returns()
func (_ERC721Burnable *ERC721BurnableTransactorSession) TransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _ERC721Burnable.Contract.TransferFrom(&_ERC721Burnable.TransactOpts, from, to, tokenId)
}

// ERC721BurnableApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the ERC721Burnable contract.
type ERC721BurnableApprovalIterator struct {
	Event *ERC721BurnableApproval // Event containing the contract specifics and raw log

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
func (it *ERC721BurnableApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC721BurnableApproval)
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
		it.Event = new(ERC721BurnableApproval)
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
func (it *ERC721BurnableApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC721BurnableApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC721BurnableApproval represents a Approval event raised by the ERC721Burnable contract.
type ERC721BurnableApproval struct {
	Owner    common.Address
	Approved common.Address
	TokenId  *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: e Approval(owner indexed address, approved indexed address, tokenId indexed uint256)
func (_ERC721Burnable *ERC721BurnableFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, approved []common.Address, tokenId []*big.Int) (*ERC721BurnableApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var approvedRule []interface{}
	for _, approvedItem := range approved {
		approvedRule = append(approvedRule, approvedItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _ERC721Burnable.contract.FilterLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &ERC721BurnableApprovalIterator{contract: _ERC721Burnable.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: e Approval(owner indexed address, approved indexed address, tokenId indexed uint256)
func (_ERC721Burnable *ERC721BurnableFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *ERC721BurnableApproval, owner []common.Address, approved []common.Address, tokenId []*big.Int) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var approvedRule []interface{}
	for _, approvedItem := range approved {
		approvedRule = append(approvedRule, approvedItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _ERC721Burnable.contract.WatchLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC721BurnableApproval)
				if err := _ERC721Burnable.contract.UnpackLog(event, "Approval", log); err != nil {
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

// ERC721BurnableApprovalForAllIterator is returned from FilterApprovalForAll and is used to iterate over the raw logs and unpacked data for ApprovalForAll events raised by the ERC721Burnable contract.
type ERC721BurnableApprovalForAllIterator struct {
	Event *ERC721BurnableApprovalForAll // Event containing the contract specifics and raw log

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
func (it *ERC721BurnableApprovalForAllIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC721BurnableApprovalForAll)
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
		it.Event = new(ERC721BurnableApprovalForAll)
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
func (it *ERC721BurnableApprovalForAllIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC721BurnableApprovalForAllIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC721BurnableApprovalForAll represents a ApprovalForAll event raised by the ERC721Burnable contract.
type ERC721BurnableApprovalForAll struct {
	Owner    common.Address
	Operator common.Address
	Approved bool
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApprovalForAll is a free log retrieval operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: e ApprovalForAll(owner indexed address, operator indexed address, approved bool)
func (_ERC721Burnable *ERC721BurnableFilterer) FilterApprovalForAll(opts *bind.FilterOpts, owner []common.Address, operator []common.Address) (*ERC721BurnableApprovalForAllIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _ERC721Burnable.contract.FilterLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return &ERC721BurnableApprovalForAllIterator{contract: _ERC721Burnable.contract, event: "ApprovalForAll", logs: logs, sub: sub}, nil
}

// WatchApprovalForAll is a free log subscription operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: e ApprovalForAll(owner indexed address, operator indexed address, approved bool)
func (_ERC721Burnable *ERC721BurnableFilterer) WatchApprovalForAll(opts *bind.WatchOpts, sink chan<- *ERC721BurnableApprovalForAll, owner []common.Address, operator []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _ERC721Burnable.contract.WatchLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC721BurnableApprovalForAll)
				if err := _ERC721Burnable.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
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

// ERC721BurnableTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the ERC721Burnable contract.
type ERC721BurnableTransferIterator struct {
	Event *ERC721BurnableTransfer // Event containing the contract specifics and raw log

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
func (it *ERC721BurnableTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC721BurnableTransfer)
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
		it.Event = new(ERC721BurnableTransfer)
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
func (it *ERC721BurnableTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC721BurnableTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC721BurnableTransfer represents a Transfer event raised by the ERC721Burnable contract.
type ERC721BurnableTransfer struct {
	From    common.Address
	To      common.Address
	TokenId *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: e Transfer(from indexed address, to indexed address, tokenId indexed uint256)
func (_ERC721Burnable *ERC721BurnableFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address, tokenId []*big.Int) (*ERC721BurnableTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _ERC721Burnable.contract.FilterLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &ERC721BurnableTransferIterator{contract: _ERC721Burnable.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: e Transfer(from indexed address, to indexed address, tokenId indexed uint256)
func (_ERC721Burnable *ERC721BurnableFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *ERC721BurnableTransfer, from []common.Address, to []common.Address, tokenId []*big.Int) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _ERC721Burnable.contract.WatchLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC721BurnableTransfer)
				if err := _ERC721Burnable.contract.UnpackLog(event, "Transfer", log); err != nil {
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

// ERC721MetadataABI is the input ABI used to generate the binding from.
const ERC721MetadataABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"getApproved\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"from\",\"type\":\"address\"},{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"from\",\"type\":\"address\"},{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"safeTransferFrom\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"ownerOf\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"setApprovalForAll\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"from\",\"type\":\"address\"},{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\"},{\"name\":\"_data\",\"type\":\"bytes\"}],\"name\":\"safeTransferFrom\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"tokenURI\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"owner\",\"type\":\"address\"},{\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"isApprovedForAll\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"name\",\"type\":\"string\"},{\"name\":\"symbol\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"approved\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"ApprovalForAll\",\"type\":\"event\"}]"

// ERC721MetadataBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const ERC721MetadataBinRuntime = `0x6080604052600436106100c45763ffffffff7c010000000000000000000000000000000000000000000000000000000060003504166301ffc9a781146100c957806306fdde0314610114578063081812fc1461019e578063095ea7b3146101d257806323b872dd146101f857806342842e0e146102225780636352211e1461024c57806370a082311461026457806395d89b4114610297578063a22cb465146102ac578063b88d4fde146102d2578063c87b56dd14610341578063e985e9c514610359575b600080fd5b3480156100d557600080fd5b506101007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1960043516610380565b604080519115158252519081900360200190f35b34801561012057600080fd5b506101296103b4565b6040805160208082528351818301528351919283929083019185019080838360005b8381101561016357818101518382015260200161014b565b50505050905090810190601f1680156101905780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b3480156101aa57600080fd5b506101b660043561044a565b60408051600160a060020a039092168252519081900360200190f35b3480156101de57600080fd5b506101f6600160a060020a036004351660243561047c565b005b34801561020457600080fd5b506101f6600160a060020a0360043581169060243516604435610532565b34801561022e57600080fd5b506101f6600160a060020a03600435811690602435166044356105c0565b34801561025857600080fd5b506101b66004356105e1565b34801561027057600080fd5b50610285600160a060020a036004351661060b565b60408051918252519081900360200190f35b3480156102a357600080fd5b5061012961063e565b3480156102b857600080fd5b506101f6600160a060020a0360043516602435151561069f565b3480156102de57600080fd5b50604080516020601f6064356004818101359283018490048402850184019095528184526101f694600160a060020a0381358116956024803590921695604435953695608494019181908401838280828437509497506107239650505050505050565b34801561034d57600080fd5b5061012960043561074b565b34801561036557600080fd5b50610100600160a060020a0360043581169060243516610800565b7bffffffffffffffffffffffffffffffffffffffffffffffffffffffff191660009081526020819052604090205460ff1690565b60058054604080516020601f60026000196101006001881615020190951694909404938401819004810282018101909252828152606093909290918301828280156104405780601f1061041557610100808354040283529160200191610440565b820191906000526020600020905b81548152906001019060200180831161042357829003601f168201915b5050505050905090565b60006104558261082e565b151561046057600080fd5b50600090815260026020526040902054600160a060020a031690565b6000610487826105e1565b9050600160a060020a0383811690821614156104a257600080fd5b33600160a060020a03821614806104be57506104be8133610800565b15156104c957600080fd5b600082815260026020526040808220805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a0387811691821790925591518593918516917f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b92591a4505050565b61053c338261084b565b151561054757600080fd5b600160a060020a038216151561055c57600080fd5b61056683826108aa565b610570838261091b565b61057a82826109b1565b8082600160a060020a031684600160a060020a03167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef60405160405180910390a4505050565b6105dc8383836020604051908101604052806000815250610723565b505050565b600081815260016020526040812054600160a060020a031680151561060557600080fd5b92915050565b6000600160a060020a038216151561062257600080fd5b50600160a060020a031660009081526003602052604090205490565b60068054604080516020601f60026000196101006001881615020190951694909404938401819004810282018101909252828152606093909290918301828280156104405780601f1061041557610100808354040283529160200191610440565b600160a060020a0382163314156106b557600080fd5b336000818152600460209081526040808320600160a060020a03871680855290835292819020805460ff1916861515908117909155815190815290519293927f17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31929181900390910190a35050565b61072e848484610532565b61073a84848484610a41565b151561074557600080fd5b50505050565b60606107568261082e565b151561076157600080fd5b60008281526007602090815260409182902080548351601f6002600019610100600186161502019093169290920491820184900484028101840190945280845290918301828280156107f45780601f106107c9576101008083540402835291602001916107f4565b820191906000526020600020905b8154815290600101906020018083116107d757829003601f168201915b50505050509050919050565b600160a060020a03918216600090815260046020908152604080832093909416825291909152205460ff1690565b600090815260016020526040902054600160a060020a0316151590565b600080610857836105e1565b905080600160a060020a031684600160a060020a03161480610892575083600160a060020a03166108878461044a565b600160a060020a0316145b806108a257506108a28185610800565b949350505050565b81600160a060020a03166108bd826105e1565b600160a060020a0316146108d057600080fd5b600081815260026020526040902054600160a060020a031615610917576000818152600260205260409020805473ffffffffffffffffffffffffffffffffffffffff191690555b5050565b81600160a060020a031661092e826105e1565b600160a060020a03161461094157600080fd5b600160a060020a03821660009081526003602052604090205461096b90600163ffffffff610bc316565b600160a060020a03909216600090815260036020908152604080832094909455918152600190915220805473ffffffffffffffffffffffffffffffffffffffff19169055565b600081815260016020526040902054600160a060020a0316156109d357600080fd5b6000818152600160208181526040808420805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a0388169081179091558452600390915290912054610a2191610bda565b600160a060020a0390921660009081526003602052604090209190915550565b600080610a5685600160a060020a0316610bf3565b1515610a655760019150610bba565b6040517f150b7a020000000000000000000000000000000000000000000000000000000081523360048201818152600160a060020a03898116602485015260448401889052608060648501908152875160848601528751918a169463150b7a0294938c938b938b93909160a490910190602085019080838360005b83811015610af8578181015183820152602001610ae0565b50505050905090810190601f168015610b255780820380516001836020036101000a031916815260200191505b5095505050505050602060405180830381600087803b158015610b4757600080fd5b505af1158015610b5b573d6000803e3d6000fd5b505050506040513d6020811015610b7157600080fd5b50517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1981167f150b7a020000000000000000000000000000000000000000000000000000000014925090505b50949350505050565b60008083831115610bd357600080fd5b5050900390565b600082820183811015610bec57600080fd5b9392505050565b6000903b11905600a165627a7a72305820a0c7eb87b832cc062c7aaee5f3bbc6a489e238b81ff279a0c7f4616838286de20029`

// ERC721MetadataBin is the compiled bytecode used for deploying new contracts.
const ERC721MetadataBin = `0x60806040523480156200001157600080fd5b5060405162000e4c38038062000e4c83398101604052805160208201519082019101620000677f01ffc9a70000000000000000000000000000000000000000000000000000000064010000000062000103810204565b6200009b7f80ac58cd0000000000000000000000000000000000000000000000000000000064010000000062000103810204565b8151620000b090600590602085019062000170565b508051620000c690600690602084019062000170565b50620000fb7f5b5e139f0000000000000000000000000000000000000000000000000000000064010000000062000103810204565b505062000215565b7fffffffff0000000000000000000000000000000000000000000000000000000080821614156200013357600080fd5b7fffffffff00000000000000000000000000000000000000000000000000000000166000908152602081905260409020805460ff19166001179055565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f10620001b357805160ff1916838001178555620001e3565b82800160010185558215620001e3579182015b82811115620001e3578251825591602001919060010190620001c6565b50620001f1929150620001f5565b5090565b6200021291905b80821115620001f15760008155600101620001fc565b90565b610c2780620002256000396000f3006080604052600436106100c45763ffffffff7c010000000000000000000000000000000000000000000000000000000060003504166301ffc9a781146100c957806306fdde0314610114578063081812fc1461019e578063095ea7b3146101d257806323b872dd146101f857806342842e0e146102225780636352211e1461024c57806370a082311461026457806395d89b4114610297578063a22cb465146102ac578063b88d4fde146102d2578063c87b56dd14610341578063e985e9c514610359575b600080fd5b3480156100d557600080fd5b506101007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1960043516610380565b604080519115158252519081900360200190f35b34801561012057600080fd5b506101296103b4565b6040805160208082528351818301528351919283929083019185019080838360005b8381101561016357818101518382015260200161014b565b50505050905090810190601f1680156101905780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b3480156101aa57600080fd5b506101b660043561044a565b60408051600160a060020a039092168252519081900360200190f35b3480156101de57600080fd5b506101f6600160a060020a036004351660243561047c565b005b34801561020457600080fd5b506101f6600160a060020a0360043581169060243516604435610532565b34801561022e57600080fd5b506101f6600160a060020a03600435811690602435166044356105c0565b34801561025857600080fd5b506101b66004356105e1565b34801561027057600080fd5b50610285600160a060020a036004351661060b565b60408051918252519081900360200190f35b3480156102a357600080fd5b5061012961063e565b3480156102b857600080fd5b506101f6600160a060020a0360043516602435151561069f565b3480156102de57600080fd5b50604080516020601f6064356004818101359283018490048402850184019095528184526101f694600160a060020a0381358116956024803590921695604435953695608494019181908401838280828437509497506107239650505050505050565b34801561034d57600080fd5b5061012960043561074b565b34801561036557600080fd5b50610100600160a060020a0360043581169060243516610800565b7bffffffffffffffffffffffffffffffffffffffffffffffffffffffff191660009081526020819052604090205460ff1690565b60058054604080516020601f60026000196101006001881615020190951694909404938401819004810282018101909252828152606093909290918301828280156104405780601f1061041557610100808354040283529160200191610440565b820191906000526020600020905b81548152906001019060200180831161042357829003601f168201915b5050505050905090565b60006104558261082e565b151561046057600080fd5b50600090815260026020526040902054600160a060020a031690565b6000610487826105e1565b9050600160a060020a0383811690821614156104a257600080fd5b33600160a060020a03821614806104be57506104be8133610800565b15156104c957600080fd5b600082815260026020526040808220805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a0387811691821790925591518593918516917f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b92591a4505050565b61053c338261084b565b151561054757600080fd5b600160a060020a038216151561055c57600080fd5b61056683826108aa565b610570838261091b565b61057a82826109b1565b8082600160a060020a031684600160a060020a03167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef60405160405180910390a4505050565b6105dc8383836020604051908101604052806000815250610723565b505050565b600081815260016020526040812054600160a060020a031680151561060557600080fd5b92915050565b6000600160a060020a038216151561062257600080fd5b50600160a060020a031660009081526003602052604090205490565b60068054604080516020601f60026000196101006001881615020190951694909404938401819004810282018101909252828152606093909290918301828280156104405780601f1061041557610100808354040283529160200191610440565b600160a060020a0382163314156106b557600080fd5b336000818152600460209081526040808320600160a060020a03871680855290835292819020805460ff1916861515908117909155815190815290519293927f17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31929181900390910190a35050565b61072e848484610532565b61073a84848484610a41565b151561074557600080fd5b50505050565b60606107568261082e565b151561076157600080fd5b60008281526007602090815260409182902080548351601f6002600019610100600186161502019093169290920491820184900484028101840190945280845290918301828280156107f45780601f106107c9576101008083540402835291602001916107f4565b820191906000526020600020905b8154815290600101906020018083116107d757829003601f168201915b50505050509050919050565b600160a060020a03918216600090815260046020908152604080832093909416825291909152205460ff1690565b600090815260016020526040902054600160a060020a0316151590565b600080610857836105e1565b905080600160a060020a031684600160a060020a03161480610892575083600160a060020a03166108878461044a565b600160a060020a0316145b806108a257506108a28185610800565b949350505050565b81600160a060020a03166108bd826105e1565b600160a060020a0316146108d057600080fd5b600081815260026020526040902054600160a060020a031615610917576000818152600260205260409020805473ffffffffffffffffffffffffffffffffffffffff191690555b5050565b81600160a060020a031661092e826105e1565b600160a060020a03161461094157600080fd5b600160a060020a03821660009081526003602052604090205461096b90600163ffffffff610bc316565b600160a060020a03909216600090815260036020908152604080832094909455918152600190915220805473ffffffffffffffffffffffffffffffffffffffff19169055565b600081815260016020526040902054600160a060020a0316156109d357600080fd5b6000818152600160208181526040808420805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a0388169081179091558452600390915290912054610a2191610bda565b600160a060020a0390921660009081526003602052604090209190915550565b600080610a5685600160a060020a0316610bf3565b1515610a655760019150610bba565b6040517f150b7a020000000000000000000000000000000000000000000000000000000081523360048201818152600160a060020a03898116602485015260448401889052608060648501908152875160848601528751918a169463150b7a0294938c938b938b93909160a490910190602085019080838360005b83811015610af8578181015183820152602001610ae0565b50505050905090810190601f168015610b255780820380516001836020036101000a031916815260200191505b5095505050505050602060405180830381600087803b158015610b4757600080fd5b505af1158015610b5b573d6000803e3d6000fd5b505050506040513d6020811015610b7157600080fd5b50517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1981167f150b7a020000000000000000000000000000000000000000000000000000000014925090505b50949350505050565b60008083831115610bd357600080fd5b5050900390565b600082820183811015610bec57600080fd5b9392505050565b6000903b11905600a165627a7a72305820a0c7eb87b832cc062c7aaee5f3bbc6a489e238b81ff279a0c7f4616838286de20029`

// DeployERC721Metadata deploys a new Klaytn contract, binding an instance of ERC721Metadata to it.
func DeployERC721Metadata(auth *bind.TransactOpts, backend bind.ContractBackend, name string, symbol string) (common.Address, *types.Transaction, *ERC721Metadata, error) {
	parsed, err := abi.JSON(strings.NewReader(ERC721MetadataABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(ERC721MetadataBin), backend, name, symbol)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ERC721Metadata{ERC721MetadataCaller: ERC721MetadataCaller{contract: contract}, ERC721MetadataTransactor: ERC721MetadataTransactor{contract: contract}, ERC721MetadataFilterer: ERC721MetadataFilterer{contract: contract}}, nil
}

// ERC721Metadata is an auto generated Go binding around a Klaytn contract.
type ERC721Metadata struct {
	ERC721MetadataCaller     // Read-only binding to the contract
	ERC721MetadataTransactor // Write-only binding to the contract
	ERC721MetadataFilterer   // Log filterer for contract events
}

// ERC721MetadataCaller is an auto generated read-only Go binding around a Klaytn contract.
type ERC721MetadataCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC721MetadataTransactor is an auto generated write-only Go binding around a Klaytn contract.
type ERC721MetadataTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC721MetadataFilterer is an auto generated log filtering Go binding around a Klaytn contract events.
type ERC721MetadataFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC721MetadataSession is an auto generated Go binding around a Klaytn contract,
// with pre-set call and transact options.
type ERC721MetadataSession struct {
	Contract     *ERC721Metadata   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ERC721MetadataCallerSession is an auto generated read-only Go binding around a Klaytn contract,
// with pre-set call options.
type ERC721MetadataCallerSession struct {
	Contract *ERC721MetadataCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// ERC721MetadataTransactorSession is an auto generated write-only Go binding around a Klaytn contract,
// with pre-set transact options.
type ERC721MetadataTransactorSession struct {
	Contract     *ERC721MetadataTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// ERC721MetadataRaw is an auto generated low-level Go binding around a Klaytn contract.
type ERC721MetadataRaw struct {
	Contract *ERC721Metadata // Generic contract binding to access the raw methods on
}

// ERC721MetadataCallerRaw is an auto generated low-level read-only Go binding around a Klaytn contract.
type ERC721MetadataCallerRaw struct {
	Contract *ERC721MetadataCaller // Generic read-only contract binding to access the raw methods on
}

// ERC721MetadataTransactorRaw is an auto generated low-level write-only Go binding around a Klaytn contract.
type ERC721MetadataTransactorRaw struct {
	Contract *ERC721MetadataTransactor // Generic write-only contract binding to access the raw methods on
}

// NewERC721Metadata creates a new instance of ERC721Metadata, bound to a specific deployed contract.
func NewERC721Metadata(address common.Address, backend bind.ContractBackend) (*ERC721Metadata, error) {
	contract, err := bindERC721Metadata(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ERC721Metadata{ERC721MetadataCaller: ERC721MetadataCaller{contract: contract}, ERC721MetadataTransactor: ERC721MetadataTransactor{contract: contract}, ERC721MetadataFilterer: ERC721MetadataFilterer{contract: contract}}, nil
}

// NewERC721MetadataCaller creates a new read-only instance of ERC721Metadata, bound to a specific deployed contract.
func NewERC721MetadataCaller(address common.Address, caller bind.ContractCaller) (*ERC721MetadataCaller, error) {
	contract, err := bindERC721Metadata(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ERC721MetadataCaller{contract: contract}, nil
}

// NewERC721MetadataTransactor creates a new write-only instance of ERC721Metadata, bound to a specific deployed contract.
func NewERC721MetadataTransactor(address common.Address, transactor bind.ContractTransactor) (*ERC721MetadataTransactor, error) {
	contract, err := bindERC721Metadata(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ERC721MetadataTransactor{contract: contract}, nil
}

// NewERC721MetadataFilterer creates a new log filterer instance of ERC721Metadata, bound to a specific deployed contract.
func NewERC721MetadataFilterer(address common.Address, filterer bind.ContractFilterer) (*ERC721MetadataFilterer, error) {
	contract, err := bindERC721Metadata(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ERC721MetadataFilterer{contract: contract}, nil
}

// bindERC721Metadata binds a generic wrapper to an already deployed contract.
func bindERC721Metadata(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ERC721MetadataABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ERC721Metadata *ERC721MetadataRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _ERC721Metadata.Contract.ERC721MetadataCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ERC721Metadata *ERC721MetadataRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC721Metadata.Contract.ERC721MetadataTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ERC721Metadata *ERC721MetadataRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ERC721Metadata.Contract.ERC721MetadataTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ERC721Metadata *ERC721MetadataCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _ERC721Metadata.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ERC721Metadata *ERC721MetadataTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC721Metadata.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ERC721Metadata *ERC721MetadataTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ERC721Metadata.Contract.contract.Transact(opts, method, params...)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(owner address) constant returns(uint256)
func (_ERC721Metadata *ERC721MetadataCaller) BalanceOf(opts *bind.CallOpts, owner common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ERC721Metadata.contract.Call(opts, out, "balanceOf", owner)
	return *ret0, err
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(owner address) constant returns(uint256)
func (_ERC721Metadata *ERC721MetadataSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _ERC721Metadata.Contract.BalanceOf(&_ERC721Metadata.CallOpts, owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(owner address) constant returns(uint256)
func (_ERC721Metadata *ERC721MetadataCallerSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _ERC721Metadata.Contract.BalanceOf(&_ERC721Metadata.CallOpts, owner)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(tokenId uint256) constant returns(address)
func (_ERC721Metadata *ERC721MetadataCaller) GetApproved(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _ERC721Metadata.contract.Call(opts, out, "getApproved", tokenId)
	return *ret0, err
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(tokenId uint256) constant returns(address)
func (_ERC721Metadata *ERC721MetadataSession) GetApproved(tokenId *big.Int) (common.Address, error) {
	return _ERC721Metadata.Contract.GetApproved(&_ERC721Metadata.CallOpts, tokenId)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(tokenId uint256) constant returns(address)
func (_ERC721Metadata *ERC721MetadataCallerSession) GetApproved(tokenId *big.Int) (common.Address, error) {
	return _ERC721Metadata.Contract.GetApproved(&_ERC721Metadata.CallOpts, tokenId)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(owner address, operator address) constant returns(bool)
func (_ERC721Metadata *ERC721MetadataCaller) IsApprovedForAll(opts *bind.CallOpts, owner common.Address, operator common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _ERC721Metadata.contract.Call(opts, out, "isApprovedForAll", owner, operator)
	return *ret0, err
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(owner address, operator address) constant returns(bool)
func (_ERC721Metadata *ERC721MetadataSession) IsApprovedForAll(owner common.Address, operator common.Address) (bool, error) {
	return _ERC721Metadata.Contract.IsApprovedForAll(&_ERC721Metadata.CallOpts, owner, operator)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(owner address, operator address) constant returns(bool)
func (_ERC721Metadata *ERC721MetadataCallerSession) IsApprovedForAll(owner common.Address, operator common.Address) (bool, error) {
	return _ERC721Metadata.Contract.IsApprovedForAll(&_ERC721Metadata.CallOpts, owner, operator)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() constant returns(string)
func (_ERC721Metadata *ERC721MetadataCaller) Name(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _ERC721Metadata.contract.Call(opts, out, "name")
	return *ret0, err
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() constant returns(string)
func (_ERC721Metadata *ERC721MetadataSession) Name() (string, error) {
	return _ERC721Metadata.Contract.Name(&_ERC721Metadata.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() constant returns(string)
func (_ERC721Metadata *ERC721MetadataCallerSession) Name() (string, error) {
	return _ERC721Metadata.Contract.Name(&_ERC721Metadata.CallOpts)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(tokenId uint256) constant returns(address)
func (_ERC721Metadata *ERC721MetadataCaller) OwnerOf(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _ERC721Metadata.contract.Call(opts, out, "ownerOf", tokenId)
	return *ret0, err
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(tokenId uint256) constant returns(address)
func (_ERC721Metadata *ERC721MetadataSession) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _ERC721Metadata.Contract.OwnerOf(&_ERC721Metadata.CallOpts, tokenId)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(tokenId uint256) constant returns(address)
func (_ERC721Metadata *ERC721MetadataCallerSession) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _ERC721Metadata.Contract.OwnerOf(&_ERC721Metadata.CallOpts, tokenId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(interfaceId bytes4) constant returns(bool)
func (_ERC721Metadata *ERC721MetadataCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _ERC721Metadata.contract.Call(opts, out, "supportsInterface", interfaceId)
	return *ret0, err
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(interfaceId bytes4) constant returns(bool)
func (_ERC721Metadata *ERC721MetadataSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _ERC721Metadata.Contract.SupportsInterface(&_ERC721Metadata.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(interfaceId bytes4) constant returns(bool)
func (_ERC721Metadata *ERC721MetadataCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _ERC721Metadata.Contract.SupportsInterface(&_ERC721Metadata.CallOpts, interfaceId)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() constant returns(string)
func (_ERC721Metadata *ERC721MetadataCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _ERC721Metadata.contract.Call(opts, out, "symbol")
	return *ret0, err
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() constant returns(string)
func (_ERC721Metadata *ERC721MetadataSession) Symbol() (string, error) {
	return _ERC721Metadata.Contract.Symbol(&_ERC721Metadata.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() constant returns(string)
func (_ERC721Metadata *ERC721MetadataCallerSession) Symbol() (string, error) {
	return _ERC721Metadata.Contract.Symbol(&_ERC721Metadata.CallOpts)
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(tokenId uint256) constant returns(string)
func (_ERC721Metadata *ERC721MetadataCaller) TokenURI(opts *bind.CallOpts, tokenId *big.Int) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _ERC721Metadata.contract.Call(opts, out, "tokenURI", tokenId)
	return *ret0, err
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(tokenId uint256) constant returns(string)
func (_ERC721Metadata *ERC721MetadataSession) TokenURI(tokenId *big.Int) (string, error) {
	return _ERC721Metadata.Contract.TokenURI(&_ERC721Metadata.CallOpts, tokenId)
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(tokenId uint256) constant returns(string)
func (_ERC721Metadata *ERC721MetadataCallerSession) TokenURI(tokenId *big.Int) (string, error) {
	return _ERC721Metadata.Contract.TokenURI(&_ERC721Metadata.CallOpts, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(to address, tokenId uint256) returns()
func (_ERC721Metadata *ERC721MetadataTransactor) Approve(opts *bind.TransactOpts, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _ERC721Metadata.contract.Transact(opts, "approve", to, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(to address, tokenId uint256) returns()
func (_ERC721Metadata *ERC721MetadataSession) Approve(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _ERC721Metadata.Contract.Approve(&_ERC721Metadata.TransactOpts, to, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(to address, tokenId uint256) returns()
func (_ERC721Metadata *ERC721MetadataTransactorSession) Approve(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _ERC721Metadata.Contract.Approve(&_ERC721Metadata.TransactOpts, to, tokenId)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(from address, to address, tokenId uint256, _data bytes) returns()
func (_ERC721Metadata *ERC721MetadataTransactor) SafeTransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int, _data []byte) (*types.Transaction, error) {
	return _ERC721Metadata.contract.Transact(opts, "safeTransferFrom", from, to, tokenId, _data)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(from address, to address, tokenId uint256, _data bytes) returns()
func (_ERC721Metadata *ERC721MetadataSession) SafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int, _data []byte) (*types.Transaction, error) {
	return _ERC721Metadata.Contract.SafeTransferFrom(&_ERC721Metadata.TransactOpts, from, to, tokenId, _data)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(from address, to address, tokenId uint256, _data bytes) returns()
func (_ERC721Metadata *ERC721MetadataTransactorSession) SafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int, _data []byte) (*types.Transaction, error) {
	return _ERC721Metadata.Contract.SafeTransferFrom(&_ERC721Metadata.TransactOpts, from, to, tokenId, _data)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(to address, approved bool) returns()
func (_ERC721Metadata *ERC721MetadataTransactor) SetApprovalForAll(opts *bind.TransactOpts, to common.Address, approved bool) (*types.Transaction, error) {
	return _ERC721Metadata.contract.Transact(opts, "setApprovalForAll", to, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(to address, approved bool) returns()
func (_ERC721Metadata *ERC721MetadataSession) SetApprovalForAll(to common.Address, approved bool) (*types.Transaction, error) {
	return _ERC721Metadata.Contract.SetApprovalForAll(&_ERC721Metadata.TransactOpts, to, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(to address, approved bool) returns()
func (_ERC721Metadata *ERC721MetadataTransactorSession) SetApprovalForAll(to common.Address, approved bool) (*types.Transaction, error) {
	return _ERC721Metadata.Contract.SetApprovalForAll(&_ERC721Metadata.TransactOpts, to, approved)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(from address, to address, tokenId uint256) returns()
func (_ERC721Metadata *ERC721MetadataTransactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _ERC721Metadata.contract.Transact(opts, "transferFrom", from, to, tokenId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(from address, to address, tokenId uint256) returns()
func (_ERC721Metadata *ERC721MetadataSession) TransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _ERC721Metadata.Contract.TransferFrom(&_ERC721Metadata.TransactOpts, from, to, tokenId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(from address, to address, tokenId uint256) returns()
func (_ERC721Metadata *ERC721MetadataTransactorSession) TransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _ERC721Metadata.Contract.TransferFrom(&_ERC721Metadata.TransactOpts, from, to, tokenId)
}

// ERC721MetadataApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the ERC721Metadata contract.
type ERC721MetadataApprovalIterator struct {
	Event *ERC721MetadataApproval // Event containing the contract specifics and raw log

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
func (it *ERC721MetadataApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC721MetadataApproval)
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
		it.Event = new(ERC721MetadataApproval)
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
func (it *ERC721MetadataApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC721MetadataApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC721MetadataApproval represents a Approval event raised by the ERC721Metadata contract.
type ERC721MetadataApproval struct {
	Owner    common.Address
	Approved common.Address
	TokenId  *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: e Approval(owner indexed address, approved indexed address, tokenId indexed uint256)
func (_ERC721Metadata *ERC721MetadataFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, approved []common.Address, tokenId []*big.Int) (*ERC721MetadataApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var approvedRule []interface{}
	for _, approvedItem := range approved {
		approvedRule = append(approvedRule, approvedItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _ERC721Metadata.contract.FilterLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &ERC721MetadataApprovalIterator{contract: _ERC721Metadata.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: e Approval(owner indexed address, approved indexed address, tokenId indexed uint256)
func (_ERC721Metadata *ERC721MetadataFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *ERC721MetadataApproval, owner []common.Address, approved []common.Address, tokenId []*big.Int) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var approvedRule []interface{}
	for _, approvedItem := range approved {
		approvedRule = append(approvedRule, approvedItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _ERC721Metadata.contract.WatchLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC721MetadataApproval)
				if err := _ERC721Metadata.contract.UnpackLog(event, "Approval", log); err != nil {
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

// ERC721MetadataApprovalForAllIterator is returned from FilterApprovalForAll and is used to iterate over the raw logs and unpacked data for ApprovalForAll events raised by the ERC721Metadata contract.
type ERC721MetadataApprovalForAllIterator struct {
	Event *ERC721MetadataApprovalForAll // Event containing the contract specifics and raw log

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
func (it *ERC721MetadataApprovalForAllIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC721MetadataApprovalForAll)
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
		it.Event = new(ERC721MetadataApprovalForAll)
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
func (it *ERC721MetadataApprovalForAllIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC721MetadataApprovalForAllIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC721MetadataApprovalForAll represents a ApprovalForAll event raised by the ERC721Metadata contract.
type ERC721MetadataApprovalForAll struct {
	Owner    common.Address
	Operator common.Address
	Approved bool
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApprovalForAll is a free log retrieval operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: e ApprovalForAll(owner indexed address, operator indexed address, approved bool)
func (_ERC721Metadata *ERC721MetadataFilterer) FilterApprovalForAll(opts *bind.FilterOpts, owner []common.Address, operator []common.Address) (*ERC721MetadataApprovalForAllIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _ERC721Metadata.contract.FilterLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return &ERC721MetadataApprovalForAllIterator{contract: _ERC721Metadata.contract, event: "ApprovalForAll", logs: logs, sub: sub}, nil
}

// WatchApprovalForAll is a free log subscription operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: e ApprovalForAll(owner indexed address, operator indexed address, approved bool)
func (_ERC721Metadata *ERC721MetadataFilterer) WatchApprovalForAll(opts *bind.WatchOpts, sink chan<- *ERC721MetadataApprovalForAll, owner []common.Address, operator []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _ERC721Metadata.contract.WatchLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC721MetadataApprovalForAll)
				if err := _ERC721Metadata.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
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

// ERC721MetadataTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the ERC721Metadata contract.
type ERC721MetadataTransferIterator struct {
	Event *ERC721MetadataTransfer // Event containing the contract specifics and raw log

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
func (it *ERC721MetadataTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC721MetadataTransfer)
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
		it.Event = new(ERC721MetadataTransfer)
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
func (it *ERC721MetadataTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC721MetadataTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC721MetadataTransfer represents a Transfer event raised by the ERC721Metadata contract.
type ERC721MetadataTransfer struct {
	From    common.Address
	To      common.Address
	TokenId *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: e Transfer(from indexed address, to indexed address, tokenId indexed uint256)
func (_ERC721Metadata *ERC721MetadataFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address, tokenId []*big.Int) (*ERC721MetadataTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _ERC721Metadata.contract.FilterLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &ERC721MetadataTransferIterator{contract: _ERC721Metadata.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: e Transfer(from indexed address, to indexed address, tokenId indexed uint256)
func (_ERC721Metadata *ERC721MetadataFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *ERC721MetadataTransfer, from []common.Address, to []common.Address, tokenId []*big.Int) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _ERC721Metadata.contract.WatchLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC721MetadataTransfer)
				if err := _ERC721Metadata.contract.UnpackLog(event, "Transfer", log); err != nil {
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

// ERC721MetadataMintableABI is the input ABI used to generate the binding from.
const ERC721MetadataMintableABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"getApproved\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"from\",\"type\":\"address\"},{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"from\",\"type\":\"address\"},{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"safeTransferFrom\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\"},{\"name\":\"tokenURI\",\"type\":\"string\"}],\"name\":\"mintWithTokenURI\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"ownerOf\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"account\",\"type\":\"address\"}],\"name\":\"addMinter\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"renounceMinter\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"setApprovalForAll\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"account\",\"type\":\"address\"}],\"name\":\"isMinter\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"from\",\"type\":\"address\"},{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\"},{\"name\":\"_data\",\"type\":\"bytes\"}],\"name\":\"safeTransferFrom\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"tokenURI\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"owner\",\"type\":\"address\"},{\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"isApprovedForAll\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"account\",\"type\":\"address\"}],\"name\":\"MinterAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"account\",\"type\":\"address\"}],\"name\":\"MinterRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"approved\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"ApprovalForAll\",\"type\":\"event\"}]"

// ERC721MetadataMintableBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const ERC721MetadataMintableBinRuntime = `0x`

// ERC721MetadataMintableBin is the compiled bytecode used for deploying new contracts.
const ERC721MetadataMintableBin = `0x`

// DeployERC721MetadataMintable deploys a new Klaytn contract, binding an instance of ERC721MetadataMintable to it.
func DeployERC721MetadataMintable(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ERC721MetadataMintable, error) {
	parsed, err := abi.JSON(strings.NewReader(ERC721MetadataMintableABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(ERC721MetadataMintableBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ERC721MetadataMintable{ERC721MetadataMintableCaller: ERC721MetadataMintableCaller{contract: contract}, ERC721MetadataMintableTransactor: ERC721MetadataMintableTransactor{contract: contract}, ERC721MetadataMintableFilterer: ERC721MetadataMintableFilterer{contract: contract}}, nil
}

// ERC721MetadataMintable is an auto generated Go binding around a Klaytn contract.
type ERC721MetadataMintable struct {
	ERC721MetadataMintableCaller     // Read-only binding to the contract
	ERC721MetadataMintableTransactor // Write-only binding to the contract
	ERC721MetadataMintableFilterer   // Log filterer for contract events
}

// ERC721MetadataMintableCaller is an auto generated read-only Go binding around a Klaytn contract.
type ERC721MetadataMintableCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC721MetadataMintableTransactor is an auto generated write-only Go binding around a Klaytn contract.
type ERC721MetadataMintableTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC721MetadataMintableFilterer is an auto generated log filtering Go binding around a Klaytn contract events.
type ERC721MetadataMintableFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC721MetadataMintableSession is an auto generated Go binding around a Klaytn contract,
// with pre-set call and transact options.
type ERC721MetadataMintableSession struct {
	Contract     *ERC721MetadataMintable // Generic contract binding to set the session for
	CallOpts     bind.CallOpts           // Call options to use throughout this session
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// ERC721MetadataMintableCallerSession is an auto generated read-only Go binding around a Klaytn contract,
// with pre-set call options.
type ERC721MetadataMintableCallerSession struct {
	Contract *ERC721MetadataMintableCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                 // Call options to use throughout this session
}

// ERC721MetadataMintableTransactorSession is an auto generated write-only Go binding around a Klaytn contract,
// with pre-set transact options.
type ERC721MetadataMintableTransactorSession struct {
	Contract     *ERC721MetadataMintableTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                 // Transaction auth options to use throughout this session
}

// ERC721MetadataMintableRaw is an auto generated low-level Go binding around a Klaytn contract.
type ERC721MetadataMintableRaw struct {
	Contract *ERC721MetadataMintable // Generic contract binding to access the raw methods on
}

// ERC721MetadataMintableCallerRaw is an auto generated low-level read-only Go binding around a Klaytn contract.
type ERC721MetadataMintableCallerRaw struct {
	Contract *ERC721MetadataMintableCaller // Generic read-only contract binding to access the raw methods on
}

// ERC721MetadataMintableTransactorRaw is an auto generated low-level write-only Go binding around a Klaytn contract.
type ERC721MetadataMintableTransactorRaw struct {
	Contract *ERC721MetadataMintableTransactor // Generic write-only contract binding to access the raw methods on
}

// NewERC721MetadataMintable creates a new instance of ERC721MetadataMintable, bound to a specific deployed contract.
func NewERC721MetadataMintable(address common.Address, backend bind.ContractBackend) (*ERC721MetadataMintable, error) {
	contract, err := bindERC721MetadataMintable(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ERC721MetadataMintable{ERC721MetadataMintableCaller: ERC721MetadataMintableCaller{contract: contract}, ERC721MetadataMintableTransactor: ERC721MetadataMintableTransactor{contract: contract}, ERC721MetadataMintableFilterer: ERC721MetadataMintableFilterer{contract: contract}}, nil
}

// NewERC721MetadataMintableCaller creates a new read-only instance of ERC721MetadataMintable, bound to a specific deployed contract.
func NewERC721MetadataMintableCaller(address common.Address, caller bind.ContractCaller) (*ERC721MetadataMintableCaller, error) {
	contract, err := bindERC721MetadataMintable(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ERC721MetadataMintableCaller{contract: contract}, nil
}

// NewERC721MetadataMintableTransactor creates a new write-only instance of ERC721MetadataMintable, bound to a specific deployed contract.
func NewERC721MetadataMintableTransactor(address common.Address, transactor bind.ContractTransactor) (*ERC721MetadataMintableTransactor, error) {
	contract, err := bindERC721MetadataMintable(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ERC721MetadataMintableTransactor{contract: contract}, nil
}

// NewERC721MetadataMintableFilterer creates a new log filterer instance of ERC721MetadataMintable, bound to a specific deployed contract.
func NewERC721MetadataMintableFilterer(address common.Address, filterer bind.ContractFilterer) (*ERC721MetadataMintableFilterer, error) {
	contract, err := bindERC721MetadataMintable(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ERC721MetadataMintableFilterer{contract: contract}, nil
}

// bindERC721MetadataMintable binds a generic wrapper to an already deployed contract.
func bindERC721MetadataMintable(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ERC721MetadataMintableABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ERC721MetadataMintable *ERC721MetadataMintableRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _ERC721MetadataMintable.Contract.ERC721MetadataMintableCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ERC721MetadataMintable *ERC721MetadataMintableRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC721MetadataMintable.Contract.ERC721MetadataMintableTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ERC721MetadataMintable *ERC721MetadataMintableRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ERC721MetadataMintable.Contract.ERC721MetadataMintableTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ERC721MetadataMintable *ERC721MetadataMintableCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _ERC721MetadataMintable.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ERC721MetadataMintable *ERC721MetadataMintableTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC721MetadataMintable.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ERC721MetadataMintable *ERC721MetadataMintableTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ERC721MetadataMintable.Contract.contract.Transact(opts, method, params...)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(owner address) constant returns(uint256)
func (_ERC721MetadataMintable *ERC721MetadataMintableCaller) BalanceOf(opts *bind.CallOpts, owner common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ERC721MetadataMintable.contract.Call(opts, out, "balanceOf", owner)
	return *ret0, err
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(owner address) constant returns(uint256)
func (_ERC721MetadataMintable *ERC721MetadataMintableSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _ERC721MetadataMintable.Contract.BalanceOf(&_ERC721MetadataMintable.CallOpts, owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(owner address) constant returns(uint256)
func (_ERC721MetadataMintable *ERC721MetadataMintableCallerSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _ERC721MetadataMintable.Contract.BalanceOf(&_ERC721MetadataMintable.CallOpts, owner)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(tokenId uint256) constant returns(address)
func (_ERC721MetadataMintable *ERC721MetadataMintableCaller) GetApproved(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _ERC721MetadataMintable.contract.Call(opts, out, "getApproved", tokenId)
	return *ret0, err
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(tokenId uint256) constant returns(address)
func (_ERC721MetadataMintable *ERC721MetadataMintableSession) GetApproved(tokenId *big.Int) (common.Address, error) {
	return _ERC721MetadataMintable.Contract.GetApproved(&_ERC721MetadataMintable.CallOpts, tokenId)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(tokenId uint256) constant returns(address)
func (_ERC721MetadataMintable *ERC721MetadataMintableCallerSession) GetApproved(tokenId *big.Int) (common.Address, error) {
	return _ERC721MetadataMintable.Contract.GetApproved(&_ERC721MetadataMintable.CallOpts, tokenId)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(owner address, operator address) constant returns(bool)
func (_ERC721MetadataMintable *ERC721MetadataMintableCaller) IsApprovedForAll(opts *bind.CallOpts, owner common.Address, operator common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _ERC721MetadataMintable.contract.Call(opts, out, "isApprovedForAll", owner, operator)
	return *ret0, err
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(owner address, operator address) constant returns(bool)
func (_ERC721MetadataMintable *ERC721MetadataMintableSession) IsApprovedForAll(owner common.Address, operator common.Address) (bool, error) {
	return _ERC721MetadataMintable.Contract.IsApprovedForAll(&_ERC721MetadataMintable.CallOpts, owner, operator)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(owner address, operator address) constant returns(bool)
func (_ERC721MetadataMintable *ERC721MetadataMintableCallerSession) IsApprovedForAll(owner common.Address, operator common.Address) (bool, error) {
	return _ERC721MetadataMintable.Contract.IsApprovedForAll(&_ERC721MetadataMintable.CallOpts, owner, operator)
}

// IsMinter is a free data retrieval call binding the contract method 0xaa271e1a.
//
// Solidity: function isMinter(account address) constant returns(bool)
func (_ERC721MetadataMintable *ERC721MetadataMintableCaller) IsMinter(opts *bind.CallOpts, account common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _ERC721MetadataMintable.contract.Call(opts, out, "isMinter", account)
	return *ret0, err
}

// IsMinter is a free data retrieval call binding the contract method 0xaa271e1a.
//
// Solidity: function isMinter(account address) constant returns(bool)
func (_ERC721MetadataMintable *ERC721MetadataMintableSession) IsMinter(account common.Address) (bool, error) {
	return _ERC721MetadataMintable.Contract.IsMinter(&_ERC721MetadataMintable.CallOpts, account)
}

// IsMinter is a free data retrieval call binding the contract method 0xaa271e1a.
//
// Solidity: function isMinter(account address) constant returns(bool)
func (_ERC721MetadataMintable *ERC721MetadataMintableCallerSession) IsMinter(account common.Address) (bool, error) {
	return _ERC721MetadataMintable.Contract.IsMinter(&_ERC721MetadataMintable.CallOpts, account)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() constant returns(string)
func (_ERC721MetadataMintable *ERC721MetadataMintableCaller) Name(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _ERC721MetadataMintable.contract.Call(opts, out, "name")
	return *ret0, err
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() constant returns(string)
func (_ERC721MetadataMintable *ERC721MetadataMintableSession) Name() (string, error) {
	return _ERC721MetadataMintable.Contract.Name(&_ERC721MetadataMintable.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() constant returns(string)
func (_ERC721MetadataMintable *ERC721MetadataMintableCallerSession) Name() (string, error) {
	return _ERC721MetadataMintable.Contract.Name(&_ERC721MetadataMintable.CallOpts)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(tokenId uint256) constant returns(address)
func (_ERC721MetadataMintable *ERC721MetadataMintableCaller) OwnerOf(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _ERC721MetadataMintable.contract.Call(opts, out, "ownerOf", tokenId)
	return *ret0, err
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(tokenId uint256) constant returns(address)
func (_ERC721MetadataMintable *ERC721MetadataMintableSession) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _ERC721MetadataMintable.Contract.OwnerOf(&_ERC721MetadataMintable.CallOpts, tokenId)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(tokenId uint256) constant returns(address)
func (_ERC721MetadataMintable *ERC721MetadataMintableCallerSession) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _ERC721MetadataMintable.Contract.OwnerOf(&_ERC721MetadataMintable.CallOpts, tokenId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(interfaceId bytes4) constant returns(bool)
func (_ERC721MetadataMintable *ERC721MetadataMintableCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _ERC721MetadataMintable.contract.Call(opts, out, "supportsInterface", interfaceId)
	return *ret0, err
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(interfaceId bytes4) constant returns(bool)
func (_ERC721MetadataMintable *ERC721MetadataMintableSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _ERC721MetadataMintable.Contract.SupportsInterface(&_ERC721MetadataMintable.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(interfaceId bytes4) constant returns(bool)
func (_ERC721MetadataMintable *ERC721MetadataMintableCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _ERC721MetadataMintable.Contract.SupportsInterface(&_ERC721MetadataMintable.CallOpts, interfaceId)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() constant returns(string)
func (_ERC721MetadataMintable *ERC721MetadataMintableCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _ERC721MetadataMintable.contract.Call(opts, out, "symbol")
	return *ret0, err
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() constant returns(string)
func (_ERC721MetadataMintable *ERC721MetadataMintableSession) Symbol() (string, error) {
	return _ERC721MetadataMintable.Contract.Symbol(&_ERC721MetadataMintable.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() constant returns(string)
func (_ERC721MetadataMintable *ERC721MetadataMintableCallerSession) Symbol() (string, error) {
	return _ERC721MetadataMintable.Contract.Symbol(&_ERC721MetadataMintable.CallOpts)
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(tokenId uint256) constant returns(string)
func (_ERC721MetadataMintable *ERC721MetadataMintableCaller) TokenURI(opts *bind.CallOpts, tokenId *big.Int) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _ERC721MetadataMintable.contract.Call(opts, out, "tokenURI", tokenId)
	return *ret0, err
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(tokenId uint256) constant returns(string)
func (_ERC721MetadataMintable *ERC721MetadataMintableSession) TokenURI(tokenId *big.Int) (string, error) {
	return _ERC721MetadataMintable.Contract.TokenURI(&_ERC721MetadataMintable.CallOpts, tokenId)
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(tokenId uint256) constant returns(string)
func (_ERC721MetadataMintable *ERC721MetadataMintableCallerSession) TokenURI(tokenId *big.Int) (string, error) {
	return _ERC721MetadataMintable.Contract.TokenURI(&_ERC721MetadataMintable.CallOpts, tokenId)
}

// AddMinter is a paid mutator transaction binding the contract method 0x983b2d56.
//
// Solidity: function addMinter(account address) returns()
func (_ERC721MetadataMintable *ERC721MetadataMintableTransactor) AddMinter(opts *bind.TransactOpts, account common.Address) (*types.Transaction, error) {
	return _ERC721MetadataMintable.contract.Transact(opts, "addMinter", account)
}

// AddMinter is a paid mutator transaction binding the contract method 0x983b2d56.
//
// Solidity: function addMinter(account address) returns()
func (_ERC721MetadataMintable *ERC721MetadataMintableSession) AddMinter(account common.Address) (*types.Transaction, error) {
	return _ERC721MetadataMintable.Contract.AddMinter(&_ERC721MetadataMintable.TransactOpts, account)
}

// AddMinter is a paid mutator transaction binding the contract method 0x983b2d56.
//
// Solidity: function addMinter(account address) returns()
func (_ERC721MetadataMintable *ERC721MetadataMintableTransactorSession) AddMinter(account common.Address) (*types.Transaction, error) {
	return _ERC721MetadataMintable.Contract.AddMinter(&_ERC721MetadataMintable.TransactOpts, account)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(to address, tokenId uint256) returns()
func (_ERC721MetadataMintable *ERC721MetadataMintableTransactor) Approve(opts *bind.TransactOpts, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _ERC721MetadataMintable.contract.Transact(opts, "approve", to, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(to address, tokenId uint256) returns()
func (_ERC721MetadataMintable *ERC721MetadataMintableSession) Approve(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _ERC721MetadataMintable.Contract.Approve(&_ERC721MetadataMintable.TransactOpts, to, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(to address, tokenId uint256) returns()
func (_ERC721MetadataMintable *ERC721MetadataMintableTransactorSession) Approve(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _ERC721MetadataMintable.Contract.Approve(&_ERC721MetadataMintable.TransactOpts, to, tokenId)
}

// MintWithTokenURI is a paid mutator transaction binding the contract method 0x50bb4e7f.
//
// Solidity: function mintWithTokenURI(to address, tokenId uint256, tokenURI string) returns(bool)
func (_ERC721MetadataMintable *ERC721MetadataMintableTransactor) MintWithTokenURI(opts *bind.TransactOpts, to common.Address, tokenId *big.Int, tokenURI string) (*types.Transaction, error) {
	return _ERC721MetadataMintable.contract.Transact(opts, "mintWithTokenURI", to, tokenId, tokenURI)
}

// MintWithTokenURI is a paid mutator transaction binding the contract method 0x50bb4e7f.
//
// Solidity: function mintWithTokenURI(to address, tokenId uint256, tokenURI string) returns(bool)
func (_ERC721MetadataMintable *ERC721MetadataMintableSession) MintWithTokenURI(to common.Address, tokenId *big.Int, tokenURI string) (*types.Transaction, error) {
	return _ERC721MetadataMintable.Contract.MintWithTokenURI(&_ERC721MetadataMintable.TransactOpts, to, tokenId, tokenURI)
}

// MintWithTokenURI is a paid mutator transaction binding the contract method 0x50bb4e7f.
//
// Solidity: function mintWithTokenURI(to address, tokenId uint256, tokenURI string) returns(bool)
func (_ERC721MetadataMintable *ERC721MetadataMintableTransactorSession) MintWithTokenURI(to common.Address, tokenId *big.Int, tokenURI string) (*types.Transaction, error) {
	return _ERC721MetadataMintable.Contract.MintWithTokenURI(&_ERC721MetadataMintable.TransactOpts, to, tokenId, tokenURI)
}

// RenounceMinter is a paid mutator transaction binding the contract method 0x98650275.
//
// Solidity: function renounceMinter() returns()
func (_ERC721MetadataMintable *ERC721MetadataMintableTransactor) RenounceMinter(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC721MetadataMintable.contract.Transact(opts, "renounceMinter")
}

// RenounceMinter is a paid mutator transaction binding the contract method 0x98650275.
//
// Solidity: function renounceMinter() returns()
func (_ERC721MetadataMintable *ERC721MetadataMintableSession) RenounceMinter() (*types.Transaction, error) {
	return _ERC721MetadataMintable.Contract.RenounceMinter(&_ERC721MetadataMintable.TransactOpts)
}

// RenounceMinter is a paid mutator transaction binding the contract method 0x98650275.
//
// Solidity: function renounceMinter() returns()
func (_ERC721MetadataMintable *ERC721MetadataMintableTransactorSession) RenounceMinter() (*types.Transaction, error) {
	return _ERC721MetadataMintable.Contract.RenounceMinter(&_ERC721MetadataMintable.TransactOpts)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(from address, to address, tokenId uint256, _data bytes) returns()
func (_ERC721MetadataMintable *ERC721MetadataMintableTransactor) SafeTransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int, _data []byte) (*types.Transaction, error) {
	return _ERC721MetadataMintable.contract.Transact(opts, "safeTransferFrom", from, to, tokenId, _data)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(from address, to address, tokenId uint256, _data bytes) returns()
func (_ERC721MetadataMintable *ERC721MetadataMintableSession) SafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int, _data []byte) (*types.Transaction, error) {
	return _ERC721MetadataMintable.Contract.SafeTransferFrom(&_ERC721MetadataMintable.TransactOpts, from, to, tokenId, _data)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(from address, to address, tokenId uint256, _data bytes) returns()
func (_ERC721MetadataMintable *ERC721MetadataMintableTransactorSession) SafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int, _data []byte) (*types.Transaction, error) {
	return _ERC721MetadataMintable.Contract.SafeTransferFrom(&_ERC721MetadataMintable.TransactOpts, from, to, tokenId, _data)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(to address, approved bool) returns()
func (_ERC721MetadataMintable *ERC721MetadataMintableTransactor) SetApprovalForAll(opts *bind.TransactOpts, to common.Address, approved bool) (*types.Transaction, error) {
	return _ERC721MetadataMintable.contract.Transact(opts, "setApprovalForAll", to, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(to address, approved bool) returns()
func (_ERC721MetadataMintable *ERC721MetadataMintableSession) SetApprovalForAll(to common.Address, approved bool) (*types.Transaction, error) {
	return _ERC721MetadataMintable.Contract.SetApprovalForAll(&_ERC721MetadataMintable.TransactOpts, to, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(to address, approved bool) returns()
func (_ERC721MetadataMintable *ERC721MetadataMintableTransactorSession) SetApprovalForAll(to common.Address, approved bool) (*types.Transaction, error) {
	return _ERC721MetadataMintable.Contract.SetApprovalForAll(&_ERC721MetadataMintable.TransactOpts, to, approved)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(from address, to address, tokenId uint256) returns()
func (_ERC721MetadataMintable *ERC721MetadataMintableTransactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _ERC721MetadataMintable.contract.Transact(opts, "transferFrom", from, to, tokenId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(from address, to address, tokenId uint256) returns()
func (_ERC721MetadataMintable *ERC721MetadataMintableSession) TransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _ERC721MetadataMintable.Contract.TransferFrom(&_ERC721MetadataMintable.TransactOpts, from, to, tokenId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(from address, to address, tokenId uint256) returns()
func (_ERC721MetadataMintable *ERC721MetadataMintableTransactorSession) TransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _ERC721MetadataMintable.Contract.TransferFrom(&_ERC721MetadataMintable.TransactOpts, from, to, tokenId)
}

// ERC721MetadataMintableApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the ERC721MetadataMintable contract.
type ERC721MetadataMintableApprovalIterator struct {
	Event *ERC721MetadataMintableApproval // Event containing the contract specifics and raw log

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
func (it *ERC721MetadataMintableApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC721MetadataMintableApproval)
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
		it.Event = new(ERC721MetadataMintableApproval)
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
func (it *ERC721MetadataMintableApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC721MetadataMintableApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC721MetadataMintableApproval represents a Approval event raised by the ERC721MetadataMintable contract.
type ERC721MetadataMintableApproval struct {
	Owner    common.Address
	Approved common.Address
	TokenId  *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: e Approval(owner indexed address, approved indexed address, tokenId indexed uint256)
func (_ERC721MetadataMintable *ERC721MetadataMintableFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, approved []common.Address, tokenId []*big.Int) (*ERC721MetadataMintableApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var approvedRule []interface{}
	for _, approvedItem := range approved {
		approvedRule = append(approvedRule, approvedItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _ERC721MetadataMintable.contract.FilterLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &ERC721MetadataMintableApprovalIterator{contract: _ERC721MetadataMintable.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: e Approval(owner indexed address, approved indexed address, tokenId indexed uint256)
func (_ERC721MetadataMintable *ERC721MetadataMintableFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *ERC721MetadataMintableApproval, owner []common.Address, approved []common.Address, tokenId []*big.Int) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var approvedRule []interface{}
	for _, approvedItem := range approved {
		approvedRule = append(approvedRule, approvedItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _ERC721MetadataMintable.contract.WatchLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC721MetadataMintableApproval)
				if err := _ERC721MetadataMintable.contract.UnpackLog(event, "Approval", log); err != nil {
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

// ERC721MetadataMintableApprovalForAllIterator is returned from FilterApprovalForAll and is used to iterate over the raw logs and unpacked data for ApprovalForAll events raised by the ERC721MetadataMintable contract.
type ERC721MetadataMintableApprovalForAllIterator struct {
	Event *ERC721MetadataMintableApprovalForAll // Event containing the contract specifics and raw log

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
func (it *ERC721MetadataMintableApprovalForAllIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC721MetadataMintableApprovalForAll)
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
		it.Event = new(ERC721MetadataMintableApprovalForAll)
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
func (it *ERC721MetadataMintableApprovalForAllIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC721MetadataMintableApprovalForAllIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC721MetadataMintableApprovalForAll represents a ApprovalForAll event raised by the ERC721MetadataMintable contract.
type ERC721MetadataMintableApprovalForAll struct {
	Owner    common.Address
	Operator common.Address
	Approved bool
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApprovalForAll is a free log retrieval operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: e ApprovalForAll(owner indexed address, operator indexed address, approved bool)
func (_ERC721MetadataMintable *ERC721MetadataMintableFilterer) FilterApprovalForAll(opts *bind.FilterOpts, owner []common.Address, operator []common.Address) (*ERC721MetadataMintableApprovalForAllIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _ERC721MetadataMintable.contract.FilterLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return &ERC721MetadataMintableApprovalForAllIterator{contract: _ERC721MetadataMintable.contract, event: "ApprovalForAll", logs: logs, sub: sub}, nil
}

// WatchApprovalForAll is a free log subscription operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: e ApprovalForAll(owner indexed address, operator indexed address, approved bool)
func (_ERC721MetadataMintable *ERC721MetadataMintableFilterer) WatchApprovalForAll(opts *bind.WatchOpts, sink chan<- *ERC721MetadataMintableApprovalForAll, owner []common.Address, operator []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _ERC721MetadataMintable.contract.WatchLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC721MetadataMintableApprovalForAll)
				if err := _ERC721MetadataMintable.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
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

// ERC721MetadataMintableMinterAddedIterator is returned from FilterMinterAdded and is used to iterate over the raw logs and unpacked data for MinterAdded events raised by the ERC721MetadataMintable contract.
type ERC721MetadataMintableMinterAddedIterator struct {
	Event *ERC721MetadataMintableMinterAdded // Event containing the contract specifics and raw log

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
func (it *ERC721MetadataMintableMinterAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC721MetadataMintableMinterAdded)
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
		it.Event = new(ERC721MetadataMintableMinterAdded)
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
func (it *ERC721MetadataMintableMinterAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC721MetadataMintableMinterAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC721MetadataMintableMinterAdded represents a MinterAdded event raised by the ERC721MetadataMintable contract.
type ERC721MetadataMintableMinterAdded struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterMinterAdded is a free log retrieval operation binding the contract event 0x6ae172837ea30b801fbfcdd4108aa1d5bf8ff775444fd70256b44e6bf3dfc3f6.
//
// Solidity: e MinterAdded(account indexed address)
func (_ERC721MetadataMintable *ERC721MetadataMintableFilterer) FilterMinterAdded(opts *bind.FilterOpts, account []common.Address) (*ERC721MetadataMintableMinterAddedIterator, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _ERC721MetadataMintable.contract.FilterLogs(opts, "MinterAdded", accountRule)
	if err != nil {
		return nil, err
	}
	return &ERC721MetadataMintableMinterAddedIterator{contract: _ERC721MetadataMintable.contract, event: "MinterAdded", logs: logs, sub: sub}, nil
}

// WatchMinterAdded is a free log subscription operation binding the contract event 0x6ae172837ea30b801fbfcdd4108aa1d5bf8ff775444fd70256b44e6bf3dfc3f6.
//
// Solidity: e MinterAdded(account indexed address)
func (_ERC721MetadataMintable *ERC721MetadataMintableFilterer) WatchMinterAdded(opts *bind.WatchOpts, sink chan<- *ERC721MetadataMintableMinterAdded, account []common.Address) (event.Subscription, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _ERC721MetadataMintable.contract.WatchLogs(opts, "MinterAdded", accountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC721MetadataMintableMinterAdded)
				if err := _ERC721MetadataMintable.contract.UnpackLog(event, "MinterAdded", log); err != nil {
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

// ERC721MetadataMintableMinterRemovedIterator is returned from FilterMinterRemoved and is used to iterate over the raw logs and unpacked data for MinterRemoved events raised by the ERC721MetadataMintable contract.
type ERC721MetadataMintableMinterRemovedIterator struct {
	Event *ERC721MetadataMintableMinterRemoved // Event containing the contract specifics and raw log

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
func (it *ERC721MetadataMintableMinterRemovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC721MetadataMintableMinterRemoved)
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
		it.Event = new(ERC721MetadataMintableMinterRemoved)
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
func (it *ERC721MetadataMintableMinterRemovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC721MetadataMintableMinterRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC721MetadataMintableMinterRemoved represents a MinterRemoved event raised by the ERC721MetadataMintable contract.
type ERC721MetadataMintableMinterRemoved struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterMinterRemoved is a free log retrieval operation binding the contract event 0xe94479a9f7e1952cc78f2d6baab678adc1b772d936c6583def489e524cb66692.
//
// Solidity: e MinterRemoved(account indexed address)
func (_ERC721MetadataMintable *ERC721MetadataMintableFilterer) FilterMinterRemoved(opts *bind.FilterOpts, account []common.Address) (*ERC721MetadataMintableMinterRemovedIterator, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _ERC721MetadataMintable.contract.FilterLogs(opts, "MinterRemoved", accountRule)
	if err != nil {
		return nil, err
	}
	return &ERC721MetadataMintableMinterRemovedIterator{contract: _ERC721MetadataMintable.contract, event: "MinterRemoved", logs: logs, sub: sub}, nil
}

// WatchMinterRemoved is a free log subscription operation binding the contract event 0xe94479a9f7e1952cc78f2d6baab678adc1b772d936c6583def489e524cb66692.
//
// Solidity: e MinterRemoved(account indexed address)
func (_ERC721MetadataMintable *ERC721MetadataMintableFilterer) WatchMinterRemoved(opts *bind.WatchOpts, sink chan<- *ERC721MetadataMintableMinterRemoved, account []common.Address) (event.Subscription, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _ERC721MetadataMintable.contract.WatchLogs(opts, "MinterRemoved", accountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC721MetadataMintableMinterRemoved)
				if err := _ERC721MetadataMintable.contract.UnpackLog(event, "MinterRemoved", log); err != nil {
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

// ERC721MetadataMintableTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the ERC721MetadataMintable contract.
type ERC721MetadataMintableTransferIterator struct {
	Event *ERC721MetadataMintableTransfer // Event containing the contract specifics and raw log

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
func (it *ERC721MetadataMintableTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC721MetadataMintableTransfer)
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
		it.Event = new(ERC721MetadataMintableTransfer)
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
func (it *ERC721MetadataMintableTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC721MetadataMintableTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC721MetadataMintableTransfer represents a Transfer event raised by the ERC721MetadataMintable contract.
type ERC721MetadataMintableTransfer struct {
	From    common.Address
	To      common.Address
	TokenId *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: e Transfer(from indexed address, to indexed address, tokenId indexed uint256)
func (_ERC721MetadataMintable *ERC721MetadataMintableFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address, tokenId []*big.Int) (*ERC721MetadataMintableTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _ERC721MetadataMintable.contract.FilterLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &ERC721MetadataMintableTransferIterator{contract: _ERC721MetadataMintable.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: e Transfer(from indexed address, to indexed address, tokenId indexed uint256)
func (_ERC721MetadataMintable *ERC721MetadataMintableFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *ERC721MetadataMintableTransfer, from []common.Address, to []common.Address, tokenId []*big.Int) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _ERC721MetadataMintable.contract.WatchLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC721MetadataMintableTransfer)
				if err := _ERC721MetadataMintable.contract.UnpackLog(event, "Transfer", log); err != nil {
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

// ExtBridgeABI is the input ABI used to generate the binding from.
const ExtBridgeABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"callback\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_tokenAddress\",\"type\":\"address\"},{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_tokenId\",\"type\":\"uint256\"},{\"name\":\"_extraData\",\"type\":\"uint256[]\"}],\"name\":\"requestERC721Transfer\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"operators\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_fee\",\"type\":\"uint256\"},{\"name\":\"_requestNonce\",\"type\":\"uint64\"}],\"name\":\"setKLAYFee\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"isRunning\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_requestTxHash\",\"type\":\"bytes32\"},{\"name\":\"_from\",\"type\":\"address\"},{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_tokenAddress\",\"type\":\"address\"},{\"name\":\"_tokenId\",\"type\":\"uint256\"},{\"name\":\"_requestNonce\",\"type\":\"uint64\"},{\"name\":\"_requestBlockNumber\",\"type\":\"uint64\"},{\"name\":\"_tokenURI\",\"type\":\"string\"},{\"name\":\"_extraData\",\"type\":\"uint256[]\"}],\"name\":\"handleERC721Transfer\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_token\",\"type\":\"address\"},{\"name\":\"_fee\",\"type\":\"uint256\"},{\"name\":\"_requestNonce\",\"type\":\"uint64\"}],\"name\":\"setERC20Fee\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"sequentialHandledRequestBlockNumber\",\"outputs\":[{\"name\":\"\",\"type\":\"uint64\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_from\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"},{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_feeLimit\",\"type\":\"uint256\"},{\"name\":\"_extraData\",\"type\":\"uint256[]\"}],\"name\":\"onERC20Received\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_operator\",\"type\":\"address\"}],\"name\":\"registerOperator\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"counterpartBridge\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_requestTxHash\",\"type\":\"bytes32\"},{\"name\":\"_from\",\"type\":\"address\"},{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_tokenAddress\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"},{\"name\":\"_requestNonce\",\"type\":\"uint64\"},{\"name\":\"_requestBlockNumber\",\"type\":\"uint64\"},{\"name\":\"_extraData\",\"type\":\"uint256[]\"}],\"name\":\"handleERC20Transfer\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_token\",\"type\":\"address\"},{\"name\":\"_cToken\",\"type\":\"address\"}],\"name\":\"registerToken\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"feeOfERC20\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint8\"}],\"name\":\"operatorThresholds\",\"outputs\":[{\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"maxHandledRequestedNonce\",\"outputs\":[{\"name\":\"\",\"type\":\"uint64\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint64\"}],\"name\":\"handledHistory\",\"outputs\":[{\"name\":\"\",\"type\":\"uint64\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"modeMintBurn\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_tokenAddress\",\"type\":\"address\"},{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"},{\"name\":\"_feeLimit\",\"type\":\"uint256\"},{\"name\":\"_extraData\",\"type\":\"uint256[]\"}],\"name\":\"requestERC20Transfer\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"requestNonce\",\"outputs\":[{\"name\":\"\",\"type\":\"uint64\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_bridge\",\"type\":\"address\"}],\"name\":\"setCounterPartBridge\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"handledRequestTx\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_addr\",\"type\":\"address\"}],\"name\":\"setCallback\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"isOwner\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\"},{\"name\":\"\",\"type\":\"address\"}],\"name\":\"votes\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint64\"}],\"name\":\"closedValueTransferVotes\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_from\",\"type\":\"address\"},{\"name\":\"_tokenId\",\"type\":\"uint256\"},{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_extraData\",\"type\":\"uint256[]\"}],\"name\":\"onERC721Received\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"votesCounts\",\"outputs\":[{\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"configurationNonce\",\"outputs\":[{\"name\":\"\",\"type\":\"uint64\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"feeReceiver\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_token\",\"type\":\"address\"}],\"name\":\"deregisterToken\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"feeOfKLAY\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_status\",\"type\":\"bool\"}],\"name\":\"start\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_operator\",\"type\":\"address\"}],\"name\":\"deregisterOperator\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"chargeWithoutEvent\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"allowedTokens\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"},{\"name\":\"_extraData\",\"type\":\"uint256[]\"}],\"name\":\"requestKLAYTransfer\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_requestTxHash\",\"type\":\"bytes32\"},{\"name\":\"_from\",\"type\":\"address\"},{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"},{\"name\":\"_requestedNonce\",\"type\":\"uint64\"},{\"name\":\"_requestedBlockNumber\",\"type\":\"uint64\"},{\"name\":\"_extraData\",\"type\":\"uint256[]\"}],\"name\":\"handleKLAYTransfer\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_voteType\",\"type\":\"uint8\"},{\"name\":\"_threshold\",\"type\":\"uint8\"}],\"name\":\"setOperatorThreshold\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_feeReceiver\",\"type\":\"address\"}],\"name\":\"setFeeReceiver\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"sequentialHandleNonce\",\"outputs\":[{\"name\":\"\",\"type\":\"uint64\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"VERSION\",\"outputs\":[{\"name\":\"\",\"type\":\"uint64\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"_modeMintBurn\",\"type\":\"bool\"}],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"constructor\"},{\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"tokenType\",\"type\":\"uint8\"},{\"indexed\":false,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"tokenAddress\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"valueOrTokenId\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"requestNonce\",\"type\":\"uint64\"},{\"indexed\":false,\"name\":\"uri\",\"type\":\"string\"},{\"indexed\":false,\"name\":\"fee\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"extraData\",\"type\":\"uint256[]\"}],\"name\":\"RequestValueTransfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"requestTxHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"tokenType\",\"type\":\"uint8\"},{\"indexed\":false,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"tokenAddress\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"valueOrTokenId\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"handleNonce\",\"type\":\"uint64\"},{\"indexed\":false,\"name\":\"extraData\",\"type\":\"uint256[]\"}],\"name\":\"HandleValueTransfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"fee\",\"type\":\"uint256\"}],\"name\":\"KLAYFeeChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"token\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"fee\",\"type\":\"uint256\"}],\"name\":\"ERC20FeeChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"feeReceiver\",\"type\":\"address\"}],\"name\":\"FeeReceiverChanged\",\"type\":\"event\"}]"

// ExtBridgeBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const ExtBridgeBinRuntime = `0x60806040526004361061020b5763ffffffff60e060020a600035041663083b27328114610245578063098540b81461027657806313e7c9d8146102af5780631a2ae53e146102e45780632014e5d1146103095780632e4bb1971461031e5780632f88396c146103ee57806330f436d01461041f578063349449b2146104515780633682a450146104c65780633a348533146104e7578063452fbc62146104fc5780634739f7e51461058c578063488af871146105b35780635526f76b146105e657806360f17a9b1461061757806367bda76b1461062c5780636e176ec21461064e578063715018a61461066357806374f41f15146106785780637c1a0302146106b557806387b04c55146106ca5780638a75eee2146106eb5780638da5cb5b146107035780638daa63ac146107185780638f32d59b146107395780639386775a1461074e5780639832c1d714610772578063a319de9e14610794578063a4155e5d14610807578063ac6fff0b1461081f578063b3f0067414610834578063bab2af1d14610849578063c263b5d61461086a578063c877cf371461087f578063d8cf98ca14610899578063dd9222d6146108ba578063e744092e146108c2578063eb7f0cbe146108e3578063ee273ad714610907578063ee2aec651461098f578063efdcd974146109b0578063f2fde38b146109d1578063febb0cc2146109f2578063ffa1ad7414610a07575b61024333600154600060405190808252806020026020018201604052801561023d578160200160208202803883390190505b50610a1c565b005b34801561025157600080fd5b5061025a610c6f565b60408051600160a060020a039092168252519081900360200190f35b34801561028257600080fd5b50610243600160a060020a0360048035821691602480359091169160443591606435908101910135610c7e565b3480156102bb57600080fd5b506102d0600160a060020a0360043516610d44565b604080519115158252519081900360200190f35b3480156102f057600080fd5b5061024360043567ffffffffffffffff60243516610d59565b34801561031557600080fd5b506102d0610d98565b34801561032a57600080fd5b50604080516020601f60e43560048181013592830184900484028501840190955281845261024394803594600160a060020a0360248035821696604435831696606435909316956084359567ffffffffffffffff60a43581169660c435909116953695929461010494929392909201918190840183828082843750506040805187358901803560208181028481018201909552818452989b9a998901989297509082019550935083925085019084908082843750949750610da89650505050505050565b3480156103fa57600080fd5b50610243600160a060020a036004351660243567ffffffffffffffff60443516610ec3565b34801561042b57600080fd5b50610434610f04565b6040805167ffffffffffffffff9092168252519081900360200190f35b34801561045d57600080fd5b50604080516084356004818101356020818102858101820190965281855261024395600160a060020a038435811696602480359760443590931696606435963696919560a49590930192909182919085019084908082843750949750610f149650505050505050565b3480156104d257600080fd5b50610243600160a060020a0360043516610f22565b3480156104f357600080fd5b5061025a610f59565b34801561050857600080fd5b50604080516020600460e4358181013583810280860185019096528085526102439583359560248035600160a060020a0390811697604435821697606435909216966084359660a43567ffffffffffffffff9081169760c435909116963696610104959194019291829185019084908082843750949750610f6d9650505050505050565b34801561059857600080fd5b50610243600160a060020a0360043581169060243516611085565b3480156105bf57600080fd5b506105d4600160a060020a03600435166110d3565b60408051918252519081900360200190f35b3480156105f257600080fd5b5061060160ff600435166110e5565b6040805160ff9092168252519081900360200190f35b34801561062357600080fd5b506104346110fa565b34801561063857600080fd5b5061043467ffffffffffffffff6004351661111e565b34801561065a57600080fd5b506102d061113a565b34801561066f57600080fd5b50610243611143565b34801561068457600080fd5b50610243600160a060020a0360048035821691602480359091169160443591606435916084359182019101356111ad565b3480156106c157600080fd5b50610434611295565b3480156106d657600080fd5b50610243600160a060020a03600435166112ac565b3480156106f757600080fd5b506102d06004356112f4565b34801561070f57600080fd5b5061025a611309565b34801561072457600080fd5b50610243600160a060020a0360043516611318565b34801561074557600080fd5b506102d061135a565b34801561075a57600080fd5b506102d0600435600160a060020a036024351661136b565b34801561077e57600080fd5b506102d067ffffffffffffffff6004351661138b565b3480156107a057600080fd5b50604080516064356004818101356020818102858101820190965281855261024395600160a060020a038435811696602480359760443590931696369690956084959394909201929091829190850190849080828437509497506113a09650505050505050565b34801561081357600080fd5b506106016004356113b3565b34801561082b57600080fd5b506104346113c8565b34801561084057600080fd5b5061025a6113d8565b34801561085557600080fd5b50610243600160a060020a03600435166113e7565b34801561087657600080fd5b506105d461142e565b34801561088b57600080fd5b506102436004351515611434565b3480156108a557600080fd5b50610243600160a060020a0360043516611477565b6102436114ab565b3480156108ce57600080fd5b5061025a600160a060020a03600435166114ad565b61024360048035600160a060020a03169060248035916044359182019101356114c8565b34801561091357600080fd5b506040805160c4356004818101356020818102858101820190965281855261024395833595600160a060020a0360248035821697604435909216966064359667ffffffffffffffff60843581169760a435909116963696939560e495019291829190850190849080828437509497506115149650505050505050565b34801561099b57600080fd5b5061024360ff6004358116906024351661166e565b3480156109bc57600080fd5b50610243600160a060020a03600435166116bb565b3480156109dd57600080fd5b50610243600160a060020a03600435166116da565b3480156109fe57600080fd5b506104346116f6565b348015610a1357600080fd5b50610434611712565b600b5460009060a860020a900460ff161515610a82576040805160e560020a62461bcd02815260206004820152600e60248201527f73746f7070656420627269646765000000000000000000000000000000000000604482015290519081900360640190fd5b348310610ad9576040805160e560020a62461bcd02815260206004820152601360248201527f696e73756666696369656e7420616d6f756e7400000000000000000000000000604482015290519081900360640190fd5b610ae283611717565b90507f09a722fc78123f7b2ba788ae83b01221477757cfb4b83302421e7c5c1fd01da96000338682610b1a348963ffffffff61184d16565b600b60169054906101000a900467ffffffffffffffff16878960405180896002811115610b4357fe5b60ff16815260200188600160a060020a0316600160a060020a0316815260200187600160a060020a0316600160a060020a0316815260200186600160a060020a0316600160a060020a031681526020018581526020018467ffffffffffffffff1667ffffffffffffffff1681526020018060200184815260200180602001838103835260008152602001602001838103825284818151815260200191508051906020019060200280838360005b83811015610c08578181015183820152602001610bf0565b505050509050019a505050505050505050505060405180910390a15050600b8054600167ffffffffffffffff60b060020a80840482169290920116027fffff0000000000000000ffffffffffffffffffffffffffffffffffffffffffff9091161790555050565b600f54600160a060020a031681565b604080517f23b872dd000000000000000000000000000000000000000000000000000000008152336004820152306024820152604481018590529051600160a060020a038716916323b872dd91606480830192600092919082900301818387803b158015610ceb57600080fd5b505af1158015610cff573d6000803e3d6000fd5b50505050610d3d8533868686868080602002602001604051908101604052809392919081815260200183836020028082843750611864945050505050565b5050505050565b60046020526000908152604090205460ff1681565b3360009081526004602052604090205460ff161515610d7757600080fd5b610d8081611c4e565b1515610d8b57610d94565b610d9482611d1a565b5050565b600b5460a860020a900460ff1681565b60008082511115610ea657816000815181101515610dc257fe5b906020019060200201519050600081118015610de85750600f54600160a060020a031615155b15610ea657600f54610e0c908b908b90600160a060020a03168a8a8a8a8a8a611d4d565b600f54604080517f44b6d69d000000000000000000000000000000000000000000000000000000008152600160a060020a038b81166004830152602482018a90528a8116604483015260648201859052915191909216916344b6d69d91608480830192600092919082900301818387803b158015610e8957600080fd5b505af1158015610e9d573d6000803e3d6000fd5b50505050610eb7565b610eb78a8a8a8a8a8a8a8a8a611d4d565b50505050505050505050565b3360009081526004602052604090205460ff161515610ee157600080fd5b610eea81611c4e565b1515610ef557610eff565b610eff8383612006565b505050565b600c5467ffffffffffffffff1681565b610d3d33868587868661205b565b610f2a61135a565b1515610f3557600080fd5b600160a060020a03166000908152600460205260409020805460ff19166001179055565b600b546101009004600160a060020a031681565b6000808251111561106a57816000815181101515610f8757fe5b906020019060200201519050600081118015610fad5750600f54600160a060020a031615155b1561106a57600f54610fd0908a908a90600160a060020a03168989898989612384565b600f54604080517f44b6d69d000000000000000000000000000000000000000000000000000000008152600160a060020a038a8116600483015260248201899052898116604483015260648201859052915191909216916344b6d69d91608480830192600092919082900301818387803b15801561104d57600080fd5b505af1158015611061573d6000803e3d6000fd5b5050505061107a565b61107a8989898989898989612384565b505050505050505050565b61108d61135a565b151561109857600080fd5b600160a060020a039182166000908152600e60205260409020805473ffffffffffffffffffffffffffffffffffffffff191691909216179055565b60026020526000908152604090205481565b60086020526000908152604090205460ff1681565b600c54700100000000000000000000000000000000900467ffffffffffffffff1681565b600d6020526000908152604090205467ffffffffffffffff1681565b600b5460ff1681565b61114b61135a565b151561115657600080fd5b600354604051600091600160a060020a0316907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0908390a36003805473ffffffffffffffffffffffffffffffffffffffff19169055565b600160a060020a0386166323b872dd33306111ce888863ffffffff6125cc16565b6040805160e060020a63ffffffff8716028152600160a060020a0394851660048201529290931660248301526044820152905160648083019260209291908290030181600087803b15801561122257600080fd5b505af1158015611236573d6000803e3d6000fd5b505050506040513d602081101561124c57600080fd5b5050604080516020838102828101820190935283825261128d92899233928a928a928a9290918a918a9182918501908490808284375061205b945050505050565b505050505050565b600b5460b060020a900467ffffffffffffffff1681565b6112b461135a565b15156112bf57600080fd5b600b8054600160a060020a039092166101000274ffffffffffffffffffffffffffffffffffffffff0019909216919091179055565b600a6020526000908152604090205460ff1681565b600354600160a060020a031690565b61132061135a565b151561132b57600080fd5b600f805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a0392909216919091179055565b600354600160a060020a0316331490565b600560209081526000928352604080842090915290825290205460ff1681565b60076020526000908152604090205460ff1681565b6113ad3385848685611864565b50505050565b60066020526000908152604090205460ff1681565b60095467ffffffffffffffff1681565b600054600160a060020a031681565b6113ef61135a565b15156113fa57600080fd5b600160a060020a03166000908152600e60205260409020805473ffffffffffffffffffffffffffffffffffffffff19169055565b60015481565b61143c61135a565b151561144757600080fd5b600b805491151560a860020a0275ff00000000000000000000000000000000000000000019909216919091179055565b61147f61135a565b151561148a57600080fd5b600160a060020a03166000908152600460205260409020805460ff19169055565b565b600e60205260009081526040902054600160a060020a031681565b60006114da348563ffffffff61184d16565b9050610d3d858285858080602002602001604051908101604052809392919081815260200183836020028082843750610a1c945050505050565b3360009081526004602052604090205460ff16151561153257600080fd5b61153b836125e5565b151561154657611665565b6040805188815260006020808301829052600160a060020a03808b169484019490945292881660608301526080820181905260a0820187905267ffffffffffffffff861660c083015261010060e0830181815285519184019190915284517f91f1d4b9864c9c091ade654e58899919d2a56e194153899c7e292cd3ae2ac47b948c948c938c9386938d938d938c939192909161012084019180860191028083838b5b838110156116005781810151838201526020016115e8565b50505050905001995050505050505050505060405180910390a1611623876126b0565b61162d83836126cb565b604051600160a060020a0386169085156108fc029086906000818181858888f19350505050158015611663573d6000803e3d6000fd5b505b50505050505050565b61167661135a565b151561168157600080fd5b806008600084600281111561169257fe5b60ff90811682526020820192909252604001600020805460ff1916929091169190911790555050565b6116c361135a565b15156116ce57600080fd5b6116d781612838565b50565b6116e261135a565b15156116ed57600080fd5b6116d78161288d565b600c5468010000000000000000900467ffffffffffffffff1681565b600181565b60015460008054909190600160a060020a0316158015906117385750600081115b156118145780831015611795576040805160e560020a62461bcd02815260206004820152601560248201527f696e73756666696369656e74206665654c696d69740000000000000000000000604482015290519081900360640190fd5b60008054604051600160a060020a039091169183156108fc02918491818181858888f193505050501580156117ce573d6000803e3d6000fd5b50336108fc6117e3858463ffffffff61184d16565b6040518115909202916000818181858888f1935050505015801561180b573d6000803e3d6000fd5b50809150611847565b604051339084156108fc029085906000818181858888f19350505050158015611841573d6000803e3d6000fd5b50600091505b50919050565b6000808383111561185d57600080fd5b5050900390565b600b5460609060a860020a900460ff1615156118ca576040805160e560020a62461bcd02815260206004820152600e60248201527f73746f7070656420627269646765000000000000000000000000000000000000604482015290519081900360640190fd5b600160a060020a038681166000908152600e602052604090205416151561193b576040805160e560020a62461bcd02815260206004820152600d60248201527f696e76616c696420746f6b656e00000000000000000000000000000000000000604482015290519081900360640190fd5b85600160a060020a031663c87b56dd846040518263ffffffff1660e060020a02815260040180828152602001915050600060405180830381600087803b15801561198457600080fd5b505af1158015611998573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f1916820160405260208110156119c157600080fd5b8101908080516401000000008111156119d957600080fd5b820160208101848111156119ec57600080fd5b8151640100000000811182820187101715611a0657600080fd5b5050600b5490945060ff16159250611a7d9150505785600160a060020a03166342966c68846040518263ffffffff1660e060020a02815260040180828152602001915050600060405180830381600087803b158015611a6457600080fd5b505af1158015611a78573d6000803e3d6000fd5b505050505b7f09a722fc78123f7b2ba788ae83b01221477757cfb4b83302421e7c5c1fd01da9600286868987600b60169054906101000a900467ffffffffffffffff168760008a604051808a6002811115611acf57fe5b60ff16815260200189600160a060020a0316600160a060020a0316815260200188600160a060020a0316600160a060020a0316815260200187600160a060020a0316600160a060020a031681526020018681526020018567ffffffffffffffff1667ffffffffffffffff1681526020018060200184815260200180602001838103835286818151815260200191508051906020019080838360005b83811015611b82578181015183820152602001611b6a565b50505050905090810190601f168015611baf5780820380516001836020036101000a031916815260200191505b508381038252845181528451602091820191808701910280838360005b83811015611be4578181015183820152602001611bcc565b505050509050019b50505050505050505050505060405180910390a15050600b8054600167ffffffffffffffff60b060020a80840482169290920116027fffff0000000000000000ffffffffffffffffffffffffffffffffffffffffffff90911617905550505050565b60095460009067ffffffffffffffff838116911614611cb7576040805160e560020a62461bcd02815260206004820152600e60248201527f6e6f6e6365206d69736d61746368000000000000000000000000000000000000604482015290519081900360640190fd5b611cdf600160003660405180838380828437820191505092505050604051809103902061290b565b15611d1157506009805467ffffffffffffffff198116600167ffffffffffffffff928316810190921617909155611d15565b5060005b919050565b600181905560405181907fa7a33d0996347e1aa55ca2206015b61b9534bdd881d59d59aa680e25eefac36590600090a250565b3360009081526004602052604090205460ff161515611d6b57600080fd5b611d74846125e5565b1515611d7f5761107a565b604080518a815260026020808301829052600160a060020a03808d1694840194909452838b166060840152928916608083015260a0820188905267ffffffffffffffff871660c083015261010060e0830181815285519184019190915284517f91f1d4b9864c9c091ade654e58899919d2a56e194153899c7e292cd3ae2ac47b948e948e938e938e938e938e938c9391929091610120840191808601910280838360005b83811015611e3b578181015183820152602001611e23565b50505050905001995050505050505050505060405180910390a1611e5e896126b0565b611e6884846126cb565b600b5460ff1615611f725785600160a060020a03166350bb4e7f8887856040518463ffffffff1660e060020a0281526004018084600160a060020a0316600160a060020a0316815260200183815260200180602001828103825283818151815260200191508051906020019080838360005b83811015611ef2578181015183820152602001611eda565b50505050905090810190601f168015611f1f5780820380516001836020036101000a031916815260200191505b50945050505050602060405180830381600087803b158015611f4057600080fd5b505af1158015611f54573d6000803e3d6000fd5b505050506040513d6020811015611f6a57600080fd5b5061107a9050565b604080517f42842e0e000000000000000000000000000000000000000000000000000000008152306004820152600160a060020a038981166024830152604482018890529151918816916342842e0e9160648082019260009290919082900301818387803b158015611fe357600080fd5b505af1158015611ff7573d6000803e3d6000fd5b50505050505050505050505050565b600160a060020a0382166000818152600260209081526040918290208490558151928352905183927fdb5ad2e76ae20cfa4e7adbc7305d7538442164d85ead9937c98620a1aa4c255b92908290030190a25050565b600b5460009060a860020a900460ff1615156120c1576040805160e560020a62461bcd02815260206004820152600e60248201527f73746f7070656420627269646765000000000000000000000000000000000000604482015290519081900360640190fd5b60008411612119576040805160e560020a62461bcd02815260206004820152600e60248201527f7a65726f206d73672e76616c7565000000000000000000000000000000000000604482015290519081900360640190fd5b600160a060020a038781166000908152600e602052604090205416151561218a576040805160e560020a62461bcd02815260206004820152600d60248201527f696e76616c696420746f6b656e00000000000000000000000000000000000000604482015290519081900360640190fd5b6121958688856129eb565b600b5490915060ff16156122055786600160a060020a03166342966c68856040518263ffffffff1660e060020a02815260040180828152602001915050600060405180830381600087803b1580156121ec57600080fd5b505af1158015612200573d6000803e3d6000fd5b505050505b7f09a722fc78123f7b2ba788ae83b01221477757cfb4b83302421e7c5c1fd01da9600187878a88600b60169054906101000a900467ffffffffffffffff1687896040518089600281111561225557fe5b60ff16815260200188600160a060020a0316600160a060020a0316815260200187600160a060020a0316600160a060020a0316815260200186600160a060020a0316600160a060020a031681526020018581526020018467ffffffffffffffff1667ffffffffffffffff1681526020018060200184815260200180602001838103835260008152602001602001838103825284818151815260200191508051906020019060200280838360005b8381101561231a578181015183820152602001612302565b505050509050019a505050505050505050505060405180910390a15050600b8054600167ffffffffffffffff60b060020a80840482169290920116027fffff0000000000000000ffffffffffffffffffffffffffffffffffffffffffff9091161790555050505050565b3360009081526004602052604090205460ff1615156123a257600080fd5b6123ab836125e5565b15156123b657611663565b6040805189815260016020808301829052600160a060020a03808c1694840194909452838a166060840152928816608083015260a0820187905267ffffffffffffffff861660c083015261010060e0830181815285519184019190915284517f91f1d4b9864c9c091ade654e58899919d2a56e194153899c7e292cd3ae2ac47b948d948d938d938d938d938d938c9391929091610120840191808601910280838360005b8381101561247257818101518382015260200161245a565b50505050905001995050505050505050505060405180910390a1612495886126b0565b61249f83836126cb565b600b5460ff161561253f5784600160a060020a03166340c10f1987866040518363ffffffff1660e060020a0281526004018083600160a060020a0316600160a060020a0316815260200182815260200192505050602060405180830381600087803b15801561250d57600080fd5b505af1158015612521573d6000803e3d6000fd5b505050506040513d602081101561253757600080fd5b506116639050565b84600160a060020a031663a9059cbb87866040518363ffffffff1660e060020a0281526004018083600160a060020a0316600160a060020a0316815260200182815260200192505050602060405180830381600087803b1580156125a257600080fd5b505af11580156125b6573d6000803e3d6000fd5b505050506040513d6020811015610eb757600080fd5b6000828201838110156125de57600080fd5b9392505050565b67ffffffffffffffff811660009081526007602052604081205460ff1615612657576040805160e560020a62461bcd02815260206004820152600b60248201527f636c6f73656420766f7465000000000000000000000000000000000000000000604482015290519081900360640190fd5b61267e6000803660405180838380828437820191505092505050604051809103902061290b565b15611d11575067ffffffffffffffff81166000908152600760205260409020805460ff19166001908117909155611d15565b6000908152600a60205260409020805460ff19166001179055565b67ffffffffffffffff8281166000818152600d60205260408120805467ffffffffffffffff1916858516179055600c54909270010000000000000000000000000000000090910416101561275a57600c805477ffffffffffffffff00000000000000000000000000000000191670010000000000000000000000000000000067ffffffffffffffff8616021790555b50600c5468010000000000000000900467ffffffffffffffff165b600c5467ffffffffffffffff7001000000000000000000000000000000009091048116908216118015906127c4575067ffffffffffffffff8082166000908152600d6020526040812054909116115b156127d157600101612775565b600c805467ffffffffffffffff80841668010000000000000000026fffffffffffffffff0000000000000000199092169190911780835560001990930181166000908152600d602052604090205467ffffffffffffffff1990931692169190911790555050565b6000805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a038316908117825560405190917f647672599d3468abcfa241a13c9e3d34383caadb5cc80fb67c3cdfcd5f78605991a250565b600160a060020a03811615156128a257600080fd5b600354604051600160a060020a038084169216907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a36003805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a0392909216919091179055565b600081815260056020908152604080832033845290915281205460ff1615156129975760008281526005602090815260408083203384528252808320805460ff1916600190811790915585845260069092529091205460ff908116918201161161297457600080fd5b6000828152600660205260409020805460ff8082166001011660ff199091161790555b600860008460028111156129a757fe5b60ff90811682526020808301939093526040918201600090812054868252600690945291909120549181169116106129e1575060016129e5565b5060005b92915050565b600160a060020a038083166000908152600260205260408120548154919290911615801590612a1a5750600081115b15612bbb5780831015612a77576040805160e560020a62461bcd02815260206004820152601560248201527f696e73756666696369656e74206665654c696d69740000000000000000000000604482015290519081900360640190fd5b60008054604080517fa9059cbb000000000000000000000000000000000000000000000000000000008152600160a060020a0392831660048201526024810185905290519187169263a9059cbb926044808401936020939083900390910190829087803b158015612ae757600080fd5b505af1158015612afb573d6000803e3d6000fd5b505050506040513d6020811015612b1157600080fd5b5050600160a060020a03841663a9059cbb86612b33868563ffffffff61184d16565b6040518363ffffffff1660e060020a0281526004018083600160a060020a0316600160a060020a0316815260200182815260200192505050602060405180830381600087803b158015612b8557600080fd5b505af1158015612b99573d6000803e3d6000fd5b505050506040513d6020811015612baf57600080fd5b50909150819050612c4f565b83600160a060020a031663a9059cbb86856040518363ffffffff1660e060020a0281526004018083600160a060020a0316600160a060020a0316815260200182815260200192505050602060405180830381600087803b158015612c1e57600080fd5b505af1158015612c32573d6000803e3d6000fd5b505050506040513d6020811015612c4857600080fd5b5060009250505b5093925050505600a165627a7a723058205d027aa195530081f9d8232f5e63b67cfbe2c185ccf31856068f315b878b43190029`

// ExtBridgeBin is the compiled bytecode used for deploying new contracts.
const ExtBridgeBin = `0x6080604081905260008054600160a060020a03199081168255600191909155600b805460ff19169055600f80549091169055602080612db1833981016040819052905160008054600160a060020a031990811682556003805490911633179081905591928392600160a060020a03169082907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0908290a35060005b600260ff821610156100cc5760ff81166000908152600860205260409020805460ff191660019081179091550161009a565b50336000908152600460205260409020805460ff19908116600117909155600b805492151560a860020a60ff021990931675010000000000000000000000000000000000000000001790911691909117905550612c838061012e6000396000f30060806040526004361061020b5763ffffffff60e060020a600035041663083b27328114610245578063098540b81461027657806313e7c9d8146102af5780631a2ae53e146102e45780632014e5d1146103095780632e4bb1971461031e5780632f88396c146103ee57806330f436d01461041f578063349449b2146104515780633682a450146104c65780633a348533146104e7578063452fbc62146104fc5780634739f7e51461058c578063488af871146105b35780635526f76b146105e657806360f17a9b1461061757806367bda76b1461062c5780636e176ec21461064e578063715018a61461066357806374f41f15146106785780637c1a0302146106b557806387b04c55146106ca5780638a75eee2146106eb5780638da5cb5b146107035780638daa63ac146107185780638f32d59b146107395780639386775a1461074e5780639832c1d714610772578063a319de9e14610794578063a4155e5d14610807578063ac6fff0b1461081f578063b3f0067414610834578063bab2af1d14610849578063c263b5d61461086a578063c877cf371461087f578063d8cf98ca14610899578063dd9222d6146108ba578063e744092e146108c2578063eb7f0cbe146108e3578063ee273ad714610907578063ee2aec651461098f578063efdcd974146109b0578063f2fde38b146109d1578063febb0cc2146109f2578063ffa1ad7414610a07575b61024333600154600060405190808252806020026020018201604052801561023d578160200160208202803883390190505b50610a1c565b005b34801561025157600080fd5b5061025a610c6f565b60408051600160a060020a039092168252519081900360200190f35b34801561028257600080fd5b50610243600160a060020a0360048035821691602480359091169160443591606435908101910135610c7e565b3480156102bb57600080fd5b506102d0600160a060020a0360043516610d44565b604080519115158252519081900360200190f35b3480156102f057600080fd5b5061024360043567ffffffffffffffff60243516610d59565b34801561031557600080fd5b506102d0610d98565b34801561032a57600080fd5b50604080516020601f60e43560048181013592830184900484028501840190955281845261024394803594600160a060020a0360248035821696604435831696606435909316956084359567ffffffffffffffff60a43581169660c435909116953695929461010494929392909201918190840183828082843750506040805187358901803560208181028481018201909552818452989b9a998901989297509082019550935083925085019084908082843750949750610da89650505050505050565b3480156103fa57600080fd5b50610243600160a060020a036004351660243567ffffffffffffffff60443516610ec3565b34801561042b57600080fd5b50610434610f04565b6040805167ffffffffffffffff9092168252519081900360200190f35b34801561045d57600080fd5b50604080516084356004818101356020818102858101820190965281855261024395600160a060020a038435811696602480359760443590931696606435963696919560a49590930192909182919085019084908082843750949750610f149650505050505050565b3480156104d257600080fd5b50610243600160a060020a0360043516610f22565b3480156104f357600080fd5b5061025a610f59565b34801561050857600080fd5b50604080516020600460e4358181013583810280860185019096528085526102439583359560248035600160a060020a0390811697604435821697606435909216966084359660a43567ffffffffffffffff9081169760c435909116963696610104959194019291829185019084908082843750949750610f6d9650505050505050565b34801561059857600080fd5b50610243600160a060020a0360043581169060243516611085565b3480156105bf57600080fd5b506105d4600160a060020a03600435166110d3565b60408051918252519081900360200190f35b3480156105f257600080fd5b5061060160ff600435166110e5565b6040805160ff9092168252519081900360200190f35b34801561062357600080fd5b506104346110fa565b34801561063857600080fd5b5061043467ffffffffffffffff6004351661111e565b34801561065a57600080fd5b506102d061113a565b34801561066f57600080fd5b50610243611143565b34801561068457600080fd5b50610243600160a060020a0360048035821691602480359091169160443591606435916084359182019101356111ad565b3480156106c157600080fd5b50610434611295565b3480156106d657600080fd5b50610243600160a060020a03600435166112ac565b3480156106f757600080fd5b506102d06004356112f4565b34801561070f57600080fd5b5061025a611309565b34801561072457600080fd5b50610243600160a060020a0360043516611318565b34801561074557600080fd5b506102d061135a565b34801561075a57600080fd5b506102d0600435600160a060020a036024351661136b565b34801561077e57600080fd5b506102d067ffffffffffffffff6004351661138b565b3480156107a057600080fd5b50604080516064356004818101356020818102858101820190965281855261024395600160a060020a038435811696602480359760443590931696369690956084959394909201929091829190850190849080828437509497506113a09650505050505050565b34801561081357600080fd5b506106016004356113b3565b34801561082b57600080fd5b506104346113c8565b34801561084057600080fd5b5061025a6113d8565b34801561085557600080fd5b50610243600160a060020a03600435166113e7565b34801561087657600080fd5b506105d461142e565b34801561088b57600080fd5b506102436004351515611434565b3480156108a557600080fd5b50610243600160a060020a0360043516611477565b6102436114ab565b3480156108ce57600080fd5b5061025a600160a060020a03600435166114ad565b61024360048035600160a060020a03169060248035916044359182019101356114c8565b34801561091357600080fd5b506040805160c4356004818101356020818102858101820190965281855261024395833595600160a060020a0360248035821697604435909216966064359667ffffffffffffffff60843581169760a435909116963696939560e495019291829190850190849080828437509497506115149650505050505050565b34801561099b57600080fd5b5061024360ff6004358116906024351661166e565b3480156109bc57600080fd5b50610243600160a060020a03600435166116bb565b3480156109dd57600080fd5b50610243600160a060020a03600435166116da565b3480156109fe57600080fd5b506104346116f6565b348015610a1357600080fd5b50610434611712565b600b5460009060a860020a900460ff161515610a82576040805160e560020a62461bcd02815260206004820152600e60248201527f73746f7070656420627269646765000000000000000000000000000000000000604482015290519081900360640190fd5b348310610ad9576040805160e560020a62461bcd02815260206004820152601360248201527f696e73756666696369656e7420616d6f756e7400000000000000000000000000604482015290519081900360640190fd5b610ae283611717565b90507f09a722fc78123f7b2ba788ae83b01221477757cfb4b83302421e7c5c1fd01da96000338682610b1a348963ffffffff61184d16565b600b60169054906101000a900467ffffffffffffffff16878960405180896002811115610b4357fe5b60ff16815260200188600160a060020a0316600160a060020a0316815260200187600160a060020a0316600160a060020a0316815260200186600160a060020a0316600160a060020a031681526020018581526020018467ffffffffffffffff1667ffffffffffffffff1681526020018060200184815260200180602001838103835260008152602001602001838103825284818151815260200191508051906020019060200280838360005b83811015610c08578181015183820152602001610bf0565b505050509050019a505050505050505050505060405180910390a15050600b8054600167ffffffffffffffff60b060020a80840482169290920116027fffff0000000000000000ffffffffffffffffffffffffffffffffffffffffffff9091161790555050565b600f54600160a060020a031681565b604080517f23b872dd000000000000000000000000000000000000000000000000000000008152336004820152306024820152604481018590529051600160a060020a038716916323b872dd91606480830192600092919082900301818387803b158015610ceb57600080fd5b505af1158015610cff573d6000803e3d6000fd5b50505050610d3d8533868686868080602002602001604051908101604052809392919081815260200183836020028082843750611864945050505050565b5050505050565b60046020526000908152604090205460ff1681565b3360009081526004602052604090205460ff161515610d7757600080fd5b610d8081611c4e565b1515610d8b57610d94565b610d9482611d1a565b5050565b600b5460a860020a900460ff1681565b60008082511115610ea657816000815181101515610dc257fe5b906020019060200201519050600081118015610de85750600f54600160a060020a031615155b15610ea657600f54610e0c908b908b90600160a060020a03168a8a8a8a8a8a611d4d565b600f54604080517f44b6d69d000000000000000000000000000000000000000000000000000000008152600160a060020a038b81166004830152602482018a90528a8116604483015260648201859052915191909216916344b6d69d91608480830192600092919082900301818387803b158015610e8957600080fd5b505af1158015610e9d573d6000803e3d6000fd5b50505050610eb7565b610eb78a8a8a8a8a8a8a8a8a611d4d565b50505050505050505050565b3360009081526004602052604090205460ff161515610ee157600080fd5b610eea81611c4e565b1515610ef557610eff565b610eff8383612006565b505050565b600c5467ffffffffffffffff1681565b610d3d33868587868661205b565b610f2a61135a565b1515610f3557600080fd5b600160a060020a03166000908152600460205260409020805460ff19166001179055565b600b546101009004600160a060020a031681565b6000808251111561106a57816000815181101515610f8757fe5b906020019060200201519050600081118015610fad5750600f54600160a060020a031615155b1561106a57600f54610fd0908a908a90600160a060020a03168989898989612384565b600f54604080517f44b6d69d000000000000000000000000000000000000000000000000000000008152600160a060020a038a8116600483015260248201899052898116604483015260648201859052915191909216916344b6d69d91608480830192600092919082900301818387803b15801561104d57600080fd5b505af1158015611061573d6000803e3d6000fd5b5050505061107a565b61107a8989898989898989612384565b505050505050505050565b61108d61135a565b151561109857600080fd5b600160a060020a039182166000908152600e60205260409020805473ffffffffffffffffffffffffffffffffffffffff191691909216179055565b60026020526000908152604090205481565b60086020526000908152604090205460ff1681565b600c54700100000000000000000000000000000000900467ffffffffffffffff1681565b600d6020526000908152604090205467ffffffffffffffff1681565b600b5460ff1681565b61114b61135a565b151561115657600080fd5b600354604051600091600160a060020a0316907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0908390a36003805473ffffffffffffffffffffffffffffffffffffffff19169055565b600160a060020a0386166323b872dd33306111ce888863ffffffff6125cc16565b6040805160e060020a63ffffffff8716028152600160a060020a0394851660048201529290931660248301526044820152905160648083019260209291908290030181600087803b15801561122257600080fd5b505af1158015611236573d6000803e3d6000fd5b505050506040513d602081101561124c57600080fd5b5050604080516020838102828101820190935283825261128d92899233928a928a928a9290918a918a9182918501908490808284375061205b945050505050565b505050505050565b600b5460b060020a900467ffffffffffffffff1681565b6112b461135a565b15156112bf57600080fd5b600b8054600160a060020a039092166101000274ffffffffffffffffffffffffffffffffffffffff0019909216919091179055565b600a6020526000908152604090205460ff1681565b600354600160a060020a031690565b61132061135a565b151561132b57600080fd5b600f805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a0392909216919091179055565b600354600160a060020a0316331490565b600560209081526000928352604080842090915290825290205460ff1681565b60076020526000908152604090205460ff1681565b6113ad3385848685611864565b50505050565b60066020526000908152604090205460ff1681565b60095467ffffffffffffffff1681565b600054600160a060020a031681565b6113ef61135a565b15156113fa57600080fd5b600160a060020a03166000908152600e60205260409020805473ffffffffffffffffffffffffffffffffffffffff19169055565b60015481565b61143c61135a565b151561144757600080fd5b600b805491151560a860020a0275ff00000000000000000000000000000000000000000019909216919091179055565b61147f61135a565b151561148a57600080fd5b600160a060020a03166000908152600460205260409020805460ff19169055565b565b600e60205260009081526040902054600160a060020a031681565b60006114da348563ffffffff61184d16565b9050610d3d858285858080602002602001604051908101604052809392919081815260200183836020028082843750610a1c945050505050565b3360009081526004602052604090205460ff16151561153257600080fd5b61153b836125e5565b151561154657611665565b6040805188815260006020808301829052600160a060020a03808b169484019490945292881660608301526080820181905260a0820187905267ffffffffffffffff861660c083015261010060e0830181815285519184019190915284517f91f1d4b9864c9c091ade654e58899919d2a56e194153899c7e292cd3ae2ac47b948c948c938c9386938d938d938c939192909161012084019180860191028083838b5b838110156116005781810151838201526020016115e8565b50505050905001995050505050505050505060405180910390a1611623876126b0565b61162d83836126cb565b604051600160a060020a0386169085156108fc029086906000818181858888f19350505050158015611663573d6000803e3d6000fd5b505b50505050505050565b61167661135a565b151561168157600080fd5b806008600084600281111561169257fe5b60ff90811682526020820192909252604001600020805460ff1916929091169190911790555050565b6116c361135a565b15156116ce57600080fd5b6116d781612838565b50565b6116e261135a565b15156116ed57600080fd5b6116d78161288d565b600c5468010000000000000000900467ffffffffffffffff1681565b600181565b60015460008054909190600160a060020a0316158015906117385750600081115b156118145780831015611795576040805160e560020a62461bcd02815260206004820152601560248201527f696e73756666696369656e74206665654c696d69740000000000000000000000604482015290519081900360640190fd5b60008054604051600160a060020a039091169183156108fc02918491818181858888f193505050501580156117ce573d6000803e3d6000fd5b50336108fc6117e3858463ffffffff61184d16565b6040518115909202916000818181858888f1935050505015801561180b573d6000803e3d6000fd5b50809150611847565b604051339084156108fc029085906000818181858888f19350505050158015611841573d6000803e3d6000fd5b50600091505b50919050565b6000808383111561185d57600080fd5b5050900390565b600b5460609060a860020a900460ff1615156118ca576040805160e560020a62461bcd02815260206004820152600e60248201527f73746f7070656420627269646765000000000000000000000000000000000000604482015290519081900360640190fd5b600160a060020a038681166000908152600e602052604090205416151561193b576040805160e560020a62461bcd02815260206004820152600d60248201527f696e76616c696420746f6b656e00000000000000000000000000000000000000604482015290519081900360640190fd5b85600160a060020a031663c87b56dd846040518263ffffffff1660e060020a02815260040180828152602001915050600060405180830381600087803b15801561198457600080fd5b505af1158015611998573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f1916820160405260208110156119c157600080fd5b8101908080516401000000008111156119d957600080fd5b820160208101848111156119ec57600080fd5b8151640100000000811182820187101715611a0657600080fd5b5050600b5490945060ff16159250611a7d9150505785600160a060020a03166342966c68846040518263ffffffff1660e060020a02815260040180828152602001915050600060405180830381600087803b158015611a6457600080fd5b505af1158015611a78573d6000803e3d6000fd5b505050505b7f09a722fc78123f7b2ba788ae83b01221477757cfb4b83302421e7c5c1fd01da9600286868987600b60169054906101000a900467ffffffffffffffff168760008a604051808a6002811115611acf57fe5b60ff16815260200189600160a060020a0316600160a060020a0316815260200188600160a060020a0316600160a060020a0316815260200187600160a060020a0316600160a060020a031681526020018681526020018567ffffffffffffffff1667ffffffffffffffff1681526020018060200184815260200180602001838103835286818151815260200191508051906020019080838360005b83811015611b82578181015183820152602001611b6a565b50505050905090810190601f168015611baf5780820380516001836020036101000a031916815260200191505b508381038252845181528451602091820191808701910280838360005b83811015611be4578181015183820152602001611bcc565b505050509050019b50505050505050505050505060405180910390a15050600b8054600167ffffffffffffffff60b060020a80840482169290920116027fffff0000000000000000ffffffffffffffffffffffffffffffffffffffffffff90911617905550505050565b60095460009067ffffffffffffffff838116911614611cb7576040805160e560020a62461bcd02815260206004820152600e60248201527f6e6f6e6365206d69736d61746368000000000000000000000000000000000000604482015290519081900360640190fd5b611cdf600160003660405180838380828437820191505092505050604051809103902061290b565b15611d1157506009805467ffffffffffffffff198116600167ffffffffffffffff928316810190921617909155611d15565b5060005b919050565b600181905560405181907fa7a33d0996347e1aa55ca2206015b61b9534bdd881d59d59aa680e25eefac36590600090a250565b3360009081526004602052604090205460ff161515611d6b57600080fd5b611d74846125e5565b1515611d7f5761107a565b604080518a815260026020808301829052600160a060020a03808d1694840194909452838b166060840152928916608083015260a0820188905267ffffffffffffffff871660c083015261010060e0830181815285519184019190915284517f91f1d4b9864c9c091ade654e58899919d2a56e194153899c7e292cd3ae2ac47b948e948e938e938e938e938e938c9391929091610120840191808601910280838360005b83811015611e3b578181015183820152602001611e23565b50505050905001995050505050505050505060405180910390a1611e5e896126b0565b611e6884846126cb565b600b5460ff1615611f725785600160a060020a03166350bb4e7f8887856040518463ffffffff1660e060020a0281526004018084600160a060020a0316600160a060020a0316815260200183815260200180602001828103825283818151815260200191508051906020019080838360005b83811015611ef2578181015183820152602001611eda565b50505050905090810190601f168015611f1f5780820380516001836020036101000a031916815260200191505b50945050505050602060405180830381600087803b158015611f4057600080fd5b505af1158015611f54573d6000803e3d6000fd5b505050506040513d6020811015611f6a57600080fd5b5061107a9050565b604080517f42842e0e000000000000000000000000000000000000000000000000000000008152306004820152600160a060020a038981166024830152604482018890529151918816916342842e0e9160648082019260009290919082900301818387803b158015611fe357600080fd5b505af1158015611ff7573d6000803e3d6000fd5b50505050505050505050505050565b600160a060020a0382166000818152600260209081526040918290208490558151928352905183927fdb5ad2e76ae20cfa4e7adbc7305d7538442164d85ead9937c98620a1aa4c255b92908290030190a25050565b600b5460009060a860020a900460ff1615156120c1576040805160e560020a62461bcd02815260206004820152600e60248201527f73746f7070656420627269646765000000000000000000000000000000000000604482015290519081900360640190fd5b60008411612119576040805160e560020a62461bcd02815260206004820152600e60248201527f7a65726f206d73672e76616c7565000000000000000000000000000000000000604482015290519081900360640190fd5b600160a060020a038781166000908152600e602052604090205416151561218a576040805160e560020a62461bcd02815260206004820152600d60248201527f696e76616c696420746f6b656e00000000000000000000000000000000000000604482015290519081900360640190fd5b6121958688856129eb565b600b5490915060ff16156122055786600160a060020a03166342966c68856040518263ffffffff1660e060020a02815260040180828152602001915050600060405180830381600087803b1580156121ec57600080fd5b505af1158015612200573d6000803e3d6000fd5b505050505b7f09a722fc78123f7b2ba788ae83b01221477757cfb4b83302421e7c5c1fd01da9600187878a88600b60169054906101000a900467ffffffffffffffff1687896040518089600281111561225557fe5b60ff16815260200188600160a060020a0316600160a060020a0316815260200187600160a060020a0316600160a060020a0316815260200186600160a060020a0316600160a060020a031681526020018581526020018467ffffffffffffffff1667ffffffffffffffff1681526020018060200184815260200180602001838103835260008152602001602001838103825284818151815260200191508051906020019060200280838360005b8381101561231a578181015183820152602001612302565b505050509050019a505050505050505050505060405180910390a15050600b8054600167ffffffffffffffff60b060020a80840482169290920116027fffff0000000000000000ffffffffffffffffffffffffffffffffffffffffffff9091161790555050505050565b3360009081526004602052604090205460ff1615156123a257600080fd5b6123ab836125e5565b15156123b657611663565b6040805189815260016020808301829052600160a060020a03808c1694840194909452838a166060840152928816608083015260a0820187905267ffffffffffffffff861660c083015261010060e0830181815285519184019190915284517f91f1d4b9864c9c091ade654e58899919d2a56e194153899c7e292cd3ae2ac47b948d948d938d938d938d938d938c9391929091610120840191808601910280838360005b8381101561247257818101518382015260200161245a565b50505050905001995050505050505050505060405180910390a1612495886126b0565b61249f83836126cb565b600b5460ff161561253f5784600160a060020a03166340c10f1987866040518363ffffffff1660e060020a0281526004018083600160a060020a0316600160a060020a0316815260200182815260200192505050602060405180830381600087803b15801561250d57600080fd5b505af1158015612521573d6000803e3d6000fd5b505050506040513d602081101561253757600080fd5b506116639050565b84600160a060020a031663a9059cbb87866040518363ffffffff1660e060020a0281526004018083600160a060020a0316600160a060020a0316815260200182815260200192505050602060405180830381600087803b1580156125a257600080fd5b505af11580156125b6573d6000803e3d6000fd5b505050506040513d6020811015610eb757600080fd5b6000828201838110156125de57600080fd5b9392505050565b67ffffffffffffffff811660009081526007602052604081205460ff1615612657576040805160e560020a62461bcd02815260206004820152600b60248201527f636c6f73656420766f7465000000000000000000000000000000000000000000604482015290519081900360640190fd5b61267e6000803660405180838380828437820191505092505050604051809103902061290b565b15611d11575067ffffffffffffffff81166000908152600760205260409020805460ff19166001908117909155611d15565b6000908152600a60205260409020805460ff19166001179055565b67ffffffffffffffff8281166000818152600d60205260408120805467ffffffffffffffff1916858516179055600c54909270010000000000000000000000000000000090910416101561275a57600c805477ffffffffffffffff00000000000000000000000000000000191670010000000000000000000000000000000067ffffffffffffffff8616021790555b50600c5468010000000000000000900467ffffffffffffffff165b600c5467ffffffffffffffff7001000000000000000000000000000000009091048116908216118015906127c4575067ffffffffffffffff8082166000908152600d6020526040812054909116115b156127d157600101612775565b600c805467ffffffffffffffff80841668010000000000000000026fffffffffffffffff0000000000000000199092169190911780835560001990930181166000908152600d602052604090205467ffffffffffffffff1990931692169190911790555050565b6000805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a038316908117825560405190917f647672599d3468abcfa241a13c9e3d34383caadb5cc80fb67c3cdfcd5f78605991a250565b600160a060020a03811615156128a257600080fd5b600354604051600160a060020a038084169216907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a36003805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a0392909216919091179055565b600081815260056020908152604080832033845290915281205460ff1615156129975760008281526005602090815260408083203384528252808320805460ff1916600190811790915585845260069092529091205460ff908116918201161161297457600080fd5b6000828152600660205260409020805460ff8082166001011660ff199091161790555b600860008460028111156129a757fe5b60ff90811682526020808301939093526040918201600090812054868252600690945291909120549181169116106129e1575060016129e5565b5060005b92915050565b600160a060020a038083166000908152600260205260408120548154919290911615801590612a1a5750600081115b15612bbb5780831015612a77576040805160e560020a62461bcd02815260206004820152601560248201527f696e73756666696369656e74206665654c696d69740000000000000000000000604482015290519081900360640190fd5b60008054604080517fa9059cbb000000000000000000000000000000000000000000000000000000008152600160a060020a0392831660048201526024810185905290519187169263a9059cbb926044808401936020939083900390910190829087803b158015612ae757600080fd5b505af1158015612afb573d6000803e3d6000fd5b505050506040513d6020811015612b1157600080fd5b5050600160a060020a03841663a9059cbb86612b33868563ffffffff61184d16565b6040518363ffffffff1660e060020a0281526004018083600160a060020a0316600160a060020a0316815260200182815260200192505050602060405180830381600087803b158015612b8557600080fd5b505af1158015612b99573d6000803e3d6000fd5b505050506040513d6020811015612baf57600080fd5b50909150819050612c4f565b83600160a060020a031663a9059cbb86856040518363ffffffff1660e060020a0281526004018083600160a060020a0316600160a060020a0316815260200182815260200192505050602060405180830381600087803b158015612c1e57600080fd5b505af1158015612c32573d6000803e3d6000fd5b505050506040513d6020811015612c4857600080fd5b5060009250505b5093925050505600a165627a7a723058205d027aa195530081f9d8232f5e63b67cfbe2c185ccf31856068f315b878b43190029`

// DeployExtBridge deploys a new Klaytn contract, binding an instance of ExtBridge to it.
func DeployExtBridge(auth *bind.TransactOpts, backend bind.ContractBackend, _modeMintBurn bool) (common.Address, *types.Transaction, *ExtBridge, error) {
	parsed, err := abi.JSON(strings.NewReader(ExtBridgeABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(ExtBridgeBin), backend, _modeMintBurn)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ExtBridge{ExtBridgeCaller: ExtBridgeCaller{contract: contract}, ExtBridgeTransactor: ExtBridgeTransactor{contract: contract}, ExtBridgeFilterer: ExtBridgeFilterer{contract: contract}}, nil
}

// ExtBridge is an auto generated Go binding around a Klaytn contract.
type ExtBridge struct {
	ExtBridgeCaller     // Read-only binding to the contract
	ExtBridgeTransactor // Write-only binding to the contract
	ExtBridgeFilterer   // Log filterer for contract events
}

// ExtBridgeCaller is an auto generated read-only Go binding around a Klaytn contract.
type ExtBridgeCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ExtBridgeTransactor is an auto generated write-only Go binding around a Klaytn contract.
type ExtBridgeTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ExtBridgeFilterer is an auto generated log filtering Go binding around a Klaytn contract events.
type ExtBridgeFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ExtBridgeSession is an auto generated Go binding around a Klaytn contract,
// with pre-set call and transact options.
type ExtBridgeSession struct {
	Contract     *ExtBridge        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ExtBridgeCallerSession is an auto generated read-only Go binding around a Klaytn contract,
// with pre-set call options.
type ExtBridgeCallerSession struct {
	Contract *ExtBridgeCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// ExtBridgeTransactorSession is an auto generated write-only Go binding around a Klaytn contract,
// with pre-set transact options.
type ExtBridgeTransactorSession struct {
	Contract     *ExtBridgeTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// ExtBridgeRaw is an auto generated low-level Go binding around a Klaytn contract.
type ExtBridgeRaw struct {
	Contract *ExtBridge // Generic contract binding to access the raw methods on
}

// ExtBridgeCallerRaw is an auto generated low-level read-only Go binding around a Klaytn contract.
type ExtBridgeCallerRaw struct {
	Contract *ExtBridgeCaller // Generic read-only contract binding to access the raw methods on
}

// ExtBridgeTransactorRaw is an auto generated low-level write-only Go binding around a Klaytn contract.
type ExtBridgeTransactorRaw struct {
	Contract *ExtBridgeTransactor // Generic write-only contract binding to access the raw methods on
}

// NewExtBridge creates a new instance of ExtBridge, bound to a specific deployed contract.
func NewExtBridge(address common.Address, backend bind.ContractBackend) (*ExtBridge, error) {
	contract, err := bindExtBridge(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ExtBridge{ExtBridgeCaller: ExtBridgeCaller{contract: contract}, ExtBridgeTransactor: ExtBridgeTransactor{contract: contract}, ExtBridgeFilterer: ExtBridgeFilterer{contract: contract}}, nil
}

// NewExtBridgeCaller creates a new read-only instance of ExtBridge, bound to a specific deployed contract.
func NewExtBridgeCaller(address common.Address, caller bind.ContractCaller) (*ExtBridgeCaller, error) {
	contract, err := bindExtBridge(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ExtBridgeCaller{contract: contract}, nil
}

// NewExtBridgeTransactor creates a new write-only instance of ExtBridge, bound to a specific deployed contract.
func NewExtBridgeTransactor(address common.Address, transactor bind.ContractTransactor) (*ExtBridgeTransactor, error) {
	contract, err := bindExtBridge(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ExtBridgeTransactor{contract: contract}, nil
}

// NewExtBridgeFilterer creates a new log filterer instance of ExtBridge, bound to a specific deployed contract.
func NewExtBridgeFilterer(address common.Address, filterer bind.ContractFilterer) (*ExtBridgeFilterer, error) {
	contract, err := bindExtBridge(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ExtBridgeFilterer{contract: contract}, nil
}

// bindExtBridge binds a generic wrapper to an already deployed contract.
func bindExtBridge(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ExtBridgeABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ExtBridge *ExtBridgeRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _ExtBridge.Contract.ExtBridgeCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ExtBridge *ExtBridgeRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ExtBridge.Contract.ExtBridgeTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ExtBridge *ExtBridgeRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ExtBridge.Contract.ExtBridgeTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ExtBridge *ExtBridgeCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _ExtBridge.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ExtBridge *ExtBridgeTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ExtBridge.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ExtBridge *ExtBridgeTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ExtBridge.Contract.contract.Transact(opts, method, params...)
}

// VERSION is a free data retrieval call binding the contract method 0xffa1ad74.
//
// Solidity: function VERSION() constant returns(uint64)
func (_ExtBridge *ExtBridgeCaller) VERSION(opts *bind.CallOpts) (uint64, error) {
	var (
		ret0 = new(uint64)
	)
	out := ret0
	err := _ExtBridge.contract.Call(opts, out, "VERSION")
	return *ret0, err
}

// VERSION is a free data retrieval call binding the contract method 0xffa1ad74.
//
// Solidity: function VERSION() constant returns(uint64)
func (_ExtBridge *ExtBridgeSession) VERSION() (uint64, error) {
	return _ExtBridge.Contract.VERSION(&_ExtBridge.CallOpts)
}

// VERSION is a free data retrieval call binding the contract method 0xffa1ad74.
//
// Solidity: function VERSION() constant returns(uint64)
func (_ExtBridge *ExtBridgeCallerSession) VERSION() (uint64, error) {
	return _ExtBridge.Contract.VERSION(&_ExtBridge.CallOpts)
}

// AllowedTokens is a free data retrieval call binding the contract method 0xe744092e.
//
// Solidity: function allowedTokens( address) constant returns(address)
func (_ExtBridge *ExtBridgeCaller) AllowedTokens(opts *bind.CallOpts, arg0 common.Address) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _ExtBridge.contract.Call(opts, out, "allowedTokens", arg0)
	return *ret0, err
}

// AllowedTokens is a free data retrieval call binding the contract method 0xe744092e.
//
// Solidity: function allowedTokens( address) constant returns(address)
func (_ExtBridge *ExtBridgeSession) AllowedTokens(arg0 common.Address) (common.Address, error) {
	return _ExtBridge.Contract.AllowedTokens(&_ExtBridge.CallOpts, arg0)
}

// AllowedTokens is a free data retrieval call binding the contract method 0xe744092e.
//
// Solidity: function allowedTokens( address) constant returns(address)
func (_ExtBridge *ExtBridgeCallerSession) AllowedTokens(arg0 common.Address) (common.Address, error) {
	return _ExtBridge.Contract.AllowedTokens(&_ExtBridge.CallOpts, arg0)
}

// Callback is a free data retrieval call binding the contract method 0x083b2732.
//
// Solidity: function callback() constant returns(address)
func (_ExtBridge *ExtBridgeCaller) Callback(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _ExtBridge.contract.Call(opts, out, "callback")
	return *ret0, err
}

// Callback is a free data retrieval call binding the contract method 0x083b2732.
//
// Solidity: function callback() constant returns(address)
func (_ExtBridge *ExtBridgeSession) Callback() (common.Address, error) {
	return _ExtBridge.Contract.Callback(&_ExtBridge.CallOpts)
}

// Callback is a free data retrieval call binding the contract method 0x083b2732.
//
// Solidity: function callback() constant returns(address)
func (_ExtBridge *ExtBridgeCallerSession) Callback() (common.Address, error) {
	return _ExtBridge.Contract.Callback(&_ExtBridge.CallOpts)
}

// ClosedValueTransferVotes is a free data retrieval call binding the contract method 0x9832c1d7.
//
// Solidity: function closedValueTransferVotes( uint64) constant returns(bool)
func (_ExtBridge *ExtBridgeCaller) ClosedValueTransferVotes(opts *bind.CallOpts, arg0 uint64) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _ExtBridge.contract.Call(opts, out, "closedValueTransferVotes", arg0)
	return *ret0, err
}

// ClosedValueTransferVotes is a free data retrieval call binding the contract method 0x9832c1d7.
//
// Solidity: function closedValueTransferVotes( uint64) constant returns(bool)
func (_ExtBridge *ExtBridgeSession) ClosedValueTransferVotes(arg0 uint64) (bool, error) {
	return _ExtBridge.Contract.ClosedValueTransferVotes(&_ExtBridge.CallOpts, arg0)
}

// ClosedValueTransferVotes is a free data retrieval call binding the contract method 0x9832c1d7.
//
// Solidity: function closedValueTransferVotes( uint64) constant returns(bool)
func (_ExtBridge *ExtBridgeCallerSession) ClosedValueTransferVotes(arg0 uint64) (bool, error) {
	return _ExtBridge.Contract.ClosedValueTransferVotes(&_ExtBridge.CallOpts, arg0)
}

// ConfigurationNonce is a free data retrieval call binding the contract method 0xac6fff0b.
//
// Solidity: function configurationNonce() constant returns(uint64)
func (_ExtBridge *ExtBridgeCaller) ConfigurationNonce(opts *bind.CallOpts) (uint64, error) {
	var (
		ret0 = new(uint64)
	)
	out := ret0
	err := _ExtBridge.contract.Call(opts, out, "configurationNonce")
	return *ret0, err
}

// ConfigurationNonce is a free data retrieval call binding the contract method 0xac6fff0b.
//
// Solidity: function configurationNonce() constant returns(uint64)
func (_ExtBridge *ExtBridgeSession) ConfigurationNonce() (uint64, error) {
	return _ExtBridge.Contract.ConfigurationNonce(&_ExtBridge.CallOpts)
}

// ConfigurationNonce is a free data retrieval call binding the contract method 0xac6fff0b.
//
// Solidity: function configurationNonce() constant returns(uint64)
func (_ExtBridge *ExtBridgeCallerSession) ConfigurationNonce() (uint64, error) {
	return _ExtBridge.Contract.ConfigurationNonce(&_ExtBridge.CallOpts)
}

// CounterpartBridge is a free data retrieval call binding the contract method 0x3a348533.
//
// Solidity: function counterpartBridge() constant returns(address)
func (_ExtBridge *ExtBridgeCaller) CounterpartBridge(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _ExtBridge.contract.Call(opts, out, "counterpartBridge")
	return *ret0, err
}

// CounterpartBridge is a free data retrieval call binding the contract method 0x3a348533.
//
// Solidity: function counterpartBridge() constant returns(address)
func (_ExtBridge *ExtBridgeSession) CounterpartBridge() (common.Address, error) {
	return _ExtBridge.Contract.CounterpartBridge(&_ExtBridge.CallOpts)
}

// CounterpartBridge is a free data retrieval call binding the contract method 0x3a348533.
//
// Solidity: function counterpartBridge() constant returns(address)
func (_ExtBridge *ExtBridgeCallerSession) CounterpartBridge() (common.Address, error) {
	return _ExtBridge.Contract.CounterpartBridge(&_ExtBridge.CallOpts)
}

// FeeOfERC20 is a free data retrieval call binding the contract method 0x488af871.
//
// Solidity: function feeOfERC20( address) constant returns(uint256)
func (_ExtBridge *ExtBridgeCaller) FeeOfERC20(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ExtBridge.contract.Call(opts, out, "feeOfERC20", arg0)
	return *ret0, err
}

// FeeOfERC20 is a free data retrieval call binding the contract method 0x488af871.
//
// Solidity: function feeOfERC20( address) constant returns(uint256)
func (_ExtBridge *ExtBridgeSession) FeeOfERC20(arg0 common.Address) (*big.Int, error) {
	return _ExtBridge.Contract.FeeOfERC20(&_ExtBridge.CallOpts, arg0)
}

// FeeOfERC20 is a free data retrieval call binding the contract method 0x488af871.
//
// Solidity: function feeOfERC20( address) constant returns(uint256)
func (_ExtBridge *ExtBridgeCallerSession) FeeOfERC20(arg0 common.Address) (*big.Int, error) {
	return _ExtBridge.Contract.FeeOfERC20(&_ExtBridge.CallOpts, arg0)
}

// FeeOfKLAY is a free data retrieval call binding the contract method 0xc263b5d6.
//
// Solidity: function feeOfKLAY() constant returns(uint256)
func (_ExtBridge *ExtBridgeCaller) FeeOfKLAY(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ExtBridge.contract.Call(opts, out, "feeOfKLAY")
	return *ret0, err
}

// FeeOfKLAY is a free data retrieval call binding the contract method 0xc263b5d6.
//
// Solidity: function feeOfKLAY() constant returns(uint256)
func (_ExtBridge *ExtBridgeSession) FeeOfKLAY() (*big.Int, error) {
	return _ExtBridge.Contract.FeeOfKLAY(&_ExtBridge.CallOpts)
}

// FeeOfKLAY is a free data retrieval call binding the contract method 0xc263b5d6.
//
// Solidity: function feeOfKLAY() constant returns(uint256)
func (_ExtBridge *ExtBridgeCallerSession) FeeOfKLAY() (*big.Int, error) {
	return _ExtBridge.Contract.FeeOfKLAY(&_ExtBridge.CallOpts)
}

// FeeReceiver is a free data retrieval call binding the contract method 0xb3f00674.
//
// Solidity: function feeReceiver() constant returns(address)
func (_ExtBridge *ExtBridgeCaller) FeeReceiver(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _ExtBridge.contract.Call(opts, out, "feeReceiver")
	return *ret0, err
}

// FeeReceiver is a free data retrieval call binding the contract method 0xb3f00674.
//
// Solidity: function feeReceiver() constant returns(address)
func (_ExtBridge *ExtBridgeSession) FeeReceiver() (common.Address, error) {
	return _ExtBridge.Contract.FeeReceiver(&_ExtBridge.CallOpts)
}

// FeeReceiver is a free data retrieval call binding the contract method 0xb3f00674.
//
// Solidity: function feeReceiver() constant returns(address)
func (_ExtBridge *ExtBridgeCallerSession) FeeReceiver() (common.Address, error) {
	return _ExtBridge.Contract.FeeReceiver(&_ExtBridge.CallOpts)
}

// HandledHistory is a free data retrieval call binding the contract method 0x67bda76b.
//
// Solidity: function handledHistory( uint64) constant returns(uint64)
func (_ExtBridge *ExtBridgeCaller) HandledHistory(opts *bind.CallOpts, arg0 uint64) (uint64, error) {
	var (
		ret0 = new(uint64)
	)
	out := ret0
	err := _ExtBridge.contract.Call(opts, out, "handledHistory", arg0)
	return *ret0, err
}

// HandledHistory is a free data retrieval call binding the contract method 0x67bda76b.
//
// Solidity: function handledHistory( uint64) constant returns(uint64)
func (_ExtBridge *ExtBridgeSession) HandledHistory(arg0 uint64) (uint64, error) {
	return _ExtBridge.Contract.HandledHistory(&_ExtBridge.CallOpts, arg0)
}

// HandledHistory is a free data retrieval call binding the contract method 0x67bda76b.
//
// Solidity: function handledHistory( uint64) constant returns(uint64)
func (_ExtBridge *ExtBridgeCallerSession) HandledHistory(arg0 uint64) (uint64, error) {
	return _ExtBridge.Contract.HandledHistory(&_ExtBridge.CallOpts, arg0)
}

// HandledRequestTx is a free data retrieval call binding the contract method 0x8a75eee2.
//
// Solidity: function handledRequestTx( bytes32) constant returns(bool)
func (_ExtBridge *ExtBridgeCaller) HandledRequestTx(opts *bind.CallOpts, arg0 [32]byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _ExtBridge.contract.Call(opts, out, "handledRequestTx", arg0)
	return *ret0, err
}

// HandledRequestTx is a free data retrieval call binding the contract method 0x8a75eee2.
//
// Solidity: function handledRequestTx( bytes32) constant returns(bool)
func (_ExtBridge *ExtBridgeSession) HandledRequestTx(arg0 [32]byte) (bool, error) {
	return _ExtBridge.Contract.HandledRequestTx(&_ExtBridge.CallOpts, arg0)
}

// HandledRequestTx is a free data retrieval call binding the contract method 0x8a75eee2.
//
// Solidity: function handledRequestTx( bytes32) constant returns(bool)
func (_ExtBridge *ExtBridgeCallerSession) HandledRequestTx(arg0 [32]byte) (bool, error) {
	return _ExtBridge.Contract.HandledRequestTx(&_ExtBridge.CallOpts, arg0)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() constant returns(bool)
func (_ExtBridge *ExtBridgeCaller) IsOwner(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _ExtBridge.contract.Call(opts, out, "isOwner")
	return *ret0, err
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() constant returns(bool)
func (_ExtBridge *ExtBridgeSession) IsOwner() (bool, error) {
	return _ExtBridge.Contract.IsOwner(&_ExtBridge.CallOpts)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() constant returns(bool)
func (_ExtBridge *ExtBridgeCallerSession) IsOwner() (bool, error) {
	return _ExtBridge.Contract.IsOwner(&_ExtBridge.CallOpts)
}

// IsRunning is a free data retrieval call binding the contract method 0x2014e5d1.
//
// Solidity: function isRunning() constant returns(bool)
func (_ExtBridge *ExtBridgeCaller) IsRunning(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _ExtBridge.contract.Call(opts, out, "isRunning")
	return *ret0, err
}

// IsRunning is a free data retrieval call binding the contract method 0x2014e5d1.
//
// Solidity: function isRunning() constant returns(bool)
func (_ExtBridge *ExtBridgeSession) IsRunning() (bool, error) {
	return _ExtBridge.Contract.IsRunning(&_ExtBridge.CallOpts)
}

// IsRunning is a free data retrieval call binding the contract method 0x2014e5d1.
//
// Solidity: function isRunning() constant returns(bool)
func (_ExtBridge *ExtBridgeCallerSession) IsRunning() (bool, error) {
	return _ExtBridge.Contract.IsRunning(&_ExtBridge.CallOpts)
}

// MaxHandledRequestedNonce is a free data retrieval call binding the contract method 0x60f17a9b.
//
// Solidity: function maxHandledRequestedNonce() constant returns(uint64)
func (_ExtBridge *ExtBridgeCaller) MaxHandledRequestedNonce(opts *bind.CallOpts) (uint64, error) {
	var (
		ret0 = new(uint64)
	)
	out := ret0
	err := _ExtBridge.contract.Call(opts, out, "maxHandledRequestedNonce")
	return *ret0, err
}

// MaxHandledRequestedNonce is a free data retrieval call binding the contract method 0x60f17a9b.
//
// Solidity: function maxHandledRequestedNonce() constant returns(uint64)
func (_ExtBridge *ExtBridgeSession) MaxHandledRequestedNonce() (uint64, error) {
	return _ExtBridge.Contract.MaxHandledRequestedNonce(&_ExtBridge.CallOpts)
}

// MaxHandledRequestedNonce is a free data retrieval call binding the contract method 0x60f17a9b.
//
// Solidity: function maxHandledRequestedNonce() constant returns(uint64)
func (_ExtBridge *ExtBridgeCallerSession) MaxHandledRequestedNonce() (uint64, error) {
	return _ExtBridge.Contract.MaxHandledRequestedNonce(&_ExtBridge.CallOpts)
}

// ModeMintBurn is a free data retrieval call binding the contract method 0x6e176ec2.
//
// Solidity: function modeMintBurn() constant returns(bool)
func (_ExtBridge *ExtBridgeCaller) ModeMintBurn(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _ExtBridge.contract.Call(opts, out, "modeMintBurn")
	return *ret0, err
}

// ModeMintBurn is a free data retrieval call binding the contract method 0x6e176ec2.
//
// Solidity: function modeMintBurn() constant returns(bool)
func (_ExtBridge *ExtBridgeSession) ModeMintBurn() (bool, error) {
	return _ExtBridge.Contract.ModeMintBurn(&_ExtBridge.CallOpts)
}

// ModeMintBurn is a free data retrieval call binding the contract method 0x6e176ec2.
//
// Solidity: function modeMintBurn() constant returns(bool)
func (_ExtBridge *ExtBridgeCallerSession) ModeMintBurn() (bool, error) {
	return _ExtBridge.Contract.ModeMintBurn(&_ExtBridge.CallOpts)
}

// OperatorThresholds is a free data retrieval call binding the contract method 0x5526f76b.
//
// Solidity: function operatorThresholds( uint8) constant returns(uint8)
func (_ExtBridge *ExtBridgeCaller) OperatorThresholds(opts *bind.CallOpts, arg0 uint8) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _ExtBridge.contract.Call(opts, out, "operatorThresholds", arg0)
	return *ret0, err
}

// OperatorThresholds is a free data retrieval call binding the contract method 0x5526f76b.
//
// Solidity: function operatorThresholds( uint8) constant returns(uint8)
func (_ExtBridge *ExtBridgeSession) OperatorThresholds(arg0 uint8) (uint8, error) {
	return _ExtBridge.Contract.OperatorThresholds(&_ExtBridge.CallOpts, arg0)
}

// OperatorThresholds is a free data retrieval call binding the contract method 0x5526f76b.
//
// Solidity: function operatorThresholds( uint8) constant returns(uint8)
func (_ExtBridge *ExtBridgeCallerSession) OperatorThresholds(arg0 uint8) (uint8, error) {
	return _ExtBridge.Contract.OperatorThresholds(&_ExtBridge.CallOpts, arg0)
}

// Operators is a free data retrieval call binding the contract method 0x13e7c9d8.
//
// Solidity: function operators( address) constant returns(bool)
func (_ExtBridge *ExtBridgeCaller) Operators(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _ExtBridge.contract.Call(opts, out, "operators", arg0)
	return *ret0, err
}

// Operators is a free data retrieval call binding the contract method 0x13e7c9d8.
//
// Solidity: function operators( address) constant returns(bool)
func (_ExtBridge *ExtBridgeSession) Operators(arg0 common.Address) (bool, error) {
	return _ExtBridge.Contract.Operators(&_ExtBridge.CallOpts, arg0)
}

// Operators is a free data retrieval call binding the contract method 0x13e7c9d8.
//
// Solidity: function operators( address) constant returns(bool)
func (_ExtBridge *ExtBridgeCallerSession) Operators(arg0 common.Address) (bool, error) {
	return _ExtBridge.Contract.Operators(&_ExtBridge.CallOpts, arg0)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_ExtBridge *ExtBridgeCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _ExtBridge.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_ExtBridge *ExtBridgeSession) Owner() (common.Address, error) {
	return _ExtBridge.Contract.Owner(&_ExtBridge.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_ExtBridge *ExtBridgeCallerSession) Owner() (common.Address, error) {
	return _ExtBridge.Contract.Owner(&_ExtBridge.CallOpts)
}

// RequestNonce is a free data retrieval call binding the contract method 0x7c1a0302.
//
// Solidity: function requestNonce() constant returns(uint64)
func (_ExtBridge *ExtBridgeCaller) RequestNonce(opts *bind.CallOpts) (uint64, error) {
	var (
		ret0 = new(uint64)
	)
	out := ret0
	err := _ExtBridge.contract.Call(opts, out, "requestNonce")
	return *ret0, err
}

// RequestNonce is a free data retrieval call binding the contract method 0x7c1a0302.
//
// Solidity: function requestNonce() constant returns(uint64)
func (_ExtBridge *ExtBridgeSession) RequestNonce() (uint64, error) {
	return _ExtBridge.Contract.RequestNonce(&_ExtBridge.CallOpts)
}

// RequestNonce is a free data retrieval call binding the contract method 0x7c1a0302.
//
// Solidity: function requestNonce() constant returns(uint64)
func (_ExtBridge *ExtBridgeCallerSession) RequestNonce() (uint64, error) {
	return _ExtBridge.Contract.RequestNonce(&_ExtBridge.CallOpts)
}

// SequentialHandleNonce is a free data retrieval call binding the contract method 0xfebb0cc2.
//
// Solidity: function sequentialHandleNonce() constant returns(uint64)
func (_ExtBridge *ExtBridgeCaller) SequentialHandleNonce(opts *bind.CallOpts) (uint64, error) {
	var (
		ret0 = new(uint64)
	)
	out := ret0
	err := _ExtBridge.contract.Call(opts, out, "sequentialHandleNonce")
	return *ret0, err
}

// SequentialHandleNonce is a free data retrieval call binding the contract method 0xfebb0cc2.
//
// Solidity: function sequentialHandleNonce() constant returns(uint64)
func (_ExtBridge *ExtBridgeSession) SequentialHandleNonce() (uint64, error) {
	return _ExtBridge.Contract.SequentialHandleNonce(&_ExtBridge.CallOpts)
}

// SequentialHandleNonce is a free data retrieval call binding the contract method 0xfebb0cc2.
//
// Solidity: function sequentialHandleNonce() constant returns(uint64)
func (_ExtBridge *ExtBridgeCallerSession) SequentialHandleNonce() (uint64, error) {
	return _ExtBridge.Contract.SequentialHandleNonce(&_ExtBridge.CallOpts)
}

// SequentialHandledRequestBlockNumber is a free data retrieval call binding the contract method 0x30f436d0.
//
// Solidity: function sequentialHandledRequestBlockNumber() constant returns(uint64)
func (_ExtBridge *ExtBridgeCaller) SequentialHandledRequestBlockNumber(opts *bind.CallOpts) (uint64, error) {
	var (
		ret0 = new(uint64)
	)
	out := ret0
	err := _ExtBridge.contract.Call(opts, out, "sequentialHandledRequestBlockNumber")
	return *ret0, err
}

// SequentialHandledRequestBlockNumber is a free data retrieval call binding the contract method 0x30f436d0.
//
// Solidity: function sequentialHandledRequestBlockNumber() constant returns(uint64)
func (_ExtBridge *ExtBridgeSession) SequentialHandledRequestBlockNumber() (uint64, error) {
	return _ExtBridge.Contract.SequentialHandledRequestBlockNumber(&_ExtBridge.CallOpts)
}

// SequentialHandledRequestBlockNumber is a free data retrieval call binding the contract method 0x30f436d0.
//
// Solidity: function sequentialHandledRequestBlockNumber() constant returns(uint64)
func (_ExtBridge *ExtBridgeCallerSession) SequentialHandledRequestBlockNumber() (uint64, error) {
	return _ExtBridge.Contract.SequentialHandledRequestBlockNumber(&_ExtBridge.CallOpts)
}

// Votes is a free data retrieval call binding the contract method 0x9386775a.
//
// Solidity: function votes( bytes32,  address) constant returns(bool)
func (_ExtBridge *ExtBridgeCaller) Votes(opts *bind.CallOpts, arg0 [32]byte, arg1 common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _ExtBridge.contract.Call(opts, out, "votes", arg0, arg1)
	return *ret0, err
}

// Votes is a free data retrieval call binding the contract method 0x9386775a.
//
// Solidity: function votes( bytes32,  address) constant returns(bool)
func (_ExtBridge *ExtBridgeSession) Votes(arg0 [32]byte, arg1 common.Address) (bool, error) {
	return _ExtBridge.Contract.Votes(&_ExtBridge.CallOpts, arg0, arg1)
}

// Votes is a free data retrieval call binding the contract method 0x9386775a.
//
// Solidity: function votes( bytes32,  address) constant returns(bool)
func (_ExtBridge *ExtBridgeCallerSession) Votes(arg0 [32]byte, arg1 common.Address) (bool, error) {
	return _ExtBridge.Contract.Votes(&_ExtBridge.CallOpts, arg0, arg1)
}

// VotesCounts is a free data retrieval call binding the contract method 0xa4155e5d.
//
// Solidity: function votesCounts( bytes32) constant returns(uint8)
func (_ExtBridge *ExtBridgeCaller) VotesCounts(opts *bind.CallOpts, arg0 [32]byte) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _ExtBridge.contract.Call(opts, out, "votesCounts", arg0)
	return *ret0, err
}

// VotesCounts is a free data retrieval call binding the contract method 0xa4155e5d.
//
// Solidity: function votesCounts( bytes32) constant returns(uint8)
func (_ExtBridge *ExtBridgeSession) VotesCounts(arg0 [32]byte) (uint8, error) {
	return _ExtBridge.Contract.VotesCounts(&_ExtBridge.CallOpts, arg0)
}

// VotesCounts is a free data retrieval call binding the contract method 0xa4155e5d.
//
// Solidity: function votesCounts( bytes32) constant returns(uint8)
func (_ExtBridge *ExtBridgeCallerSession) VotesCounts(arg0 [32]byte) (uint8, error) {
	return _ExtBridge.Contract.VotesCounts(&_ExtBridge.CallOpts, arg0)
}

// ChargeWithoutEvent is a paid mutator transaction binding the contract method 0xdd9222d6.
//
// Solidity: function chargeWithoutEvent() returns()
func (_ExtBridge *ExtBridgeTransactor) ChargeWithoutEvent(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ExtBridge.contract.Transact(opts, "chargeWithoutEvent")
}

// ChargeWithoutEvent is a paid mutator transaction binding the contract method 0xdd9222d6.
//
// Solidity: function chargeWithoutEvent() returns()
func (_ExtBridge *ExtBridgeSession) ChargeWithoutEvent() (*types.Transaction, error) {
	return _ExtBridge.Contract.ChargeWithoutEvent(&_ExtBridge.TransactOpts)
}

// ChargeWithoutEvent is a paid mutator transaction binding the contract method 0xdd9222d6.
//
// Solidity: function chargeWithoutEvent() returns()
func (_ExtBridge *ExtBridgeTransactorSession) ChargeWithoutEvent() (*types.Transaction, error) {
	return _ExtBridge.Contract.ChargeWithoutEvent(&_ExtBridge.TransactOpts)
}

// DeregisterOperator is a paid mutator transaction binding the contract method 0xd8cf98ca.
//
// Solidity: function deregisterOperator(_operator address) returns()
func (_ExtBridge *ExtBridgeTransactor) DeregisterOperator(opts *bind.TransactOpts, _operator common.Address) (*types.Transaction, error) {
	return _ExtBridge.contract.Transact(opts, "deregisterOperator", _operator)
}

// DeregisterOperator is a paid mutator transaction binding the contract method 0xd8cf98ca.
//
// Solidity: function deregisterOperator(_operator address) returns()
func (_ExtBridge *ExtBridgeSession) DeregisterOperator(_operator common.Address) (*types.Transaction, error) {
	return _ExtBridge.Contract.DeregisterOperator(&_ExtBridge.TransactOpts, _operator)
}

// DeregisterOperator is a paid mutator transaction binding the contract method 0xd8cf98ca.
//
// Solidity: function deregisterOperator(_operator address) returns()
func (_ExtBridge *ExtBridgeTransactorSession) DeregisterOperator(_operator common.Address) (*types.Transaction, error) {
	return _ExtBridge.Contract.DeregisterOperator(&_ExtBridge.TransactOpts, _operator)
}

// DeregisterToken is a paid mutator transaction binding the contract method 0xbab2af1d.
//
// Solidity: function deregisterToken(_token address) returns()
func (_ExtBridge *ExtBridgeTransactor) DeregisterToken(opts *bind.TransactOpts, _token common.Address) (*types.Transaction, error) {
	return _ExtBridge.contract.Transact(opts, "deregisterToken", _token)
}

// DeregisterToken is a paid mutator transaction binding the contract method 0xbab2af1d.
//
// Solidity: function deregisterToken(_token address) returns()
func (_ExtBridge *ExtBridgeSession) DeregisterToken(_token common.Address) (*types.Transaction, error) {
	return _ExtBridge.Contract.DeregisterToken(&_ExtBridge.TransactOpts, _token)
}

// DeregisterToken is a paid mutator transaction binding the contract method 0xbab2af1d.
//
// Solidity: function deregisterToken(_token address) returns()
func (_ExtBridge *ExtBridgeTransactorSession) DeregisterToken(_token common.Address) (*types.Transaction, error) {
	return _ExtBridge.Contract.DeregisterToken(&_ExtBridge.TransactOpts, _token)
}

// HandleERC20Transfer is a paid mutator transaction binding the contract method 0x452fbc62.
//
// Solidity: function handleERC20Transfer(_requestTxHash bytes32, _from address, _to address, _tokenAddress address, _value uint256, _requestNonce uint64, _requestBlockNumber uint64, _extraData uint256[]) returns()
func (_ExtBridge *ExtBridgeTransactor) HandleERC20Transfer(opts *bind.TransactOpts, _requestTxHash [32]byte, _from common.Address, _to common.Address, _tokenAddress common.Address, _value *big.Int, _requestNonce uint64, _requestBlockNumber uint64, _extraData []*big.Int) (*types.Transaction, error) {
	return _ExtBridge.contract.Transact(opts, "handleERC20Transfer", _requestTxHash, _from, _to, _tokenAddress, _value, _requestNonce, _requestBlockNumber, _extraData)
}

// HandleERC20Transfer is a paid mutator transaction binding the contract method 0x452fbc62.
//
// Solidity: function handleERC20Transfer(_requestTxHash bytes32, _from address, _to address, _tokenAddress address, _value uint256, _requestNonce uint64, _requestBlockNumber uint64, _extraData uint256[]) returns()
func (_ExtBridge *ExtBridgeSession) HandleERC20Transfer(_requestTxHash [32]byte, _from common.Address, _to common.Address, _tokenAddress common.Address, _value *big.Int, _requestNonce uint64, _requestBlockNumber uint64, _extraData []*big.Int) (*types.Transaction, error) {
	return _ExtBridge.Contract.HandleERC20Transfer(&_ExtBridge.TransactOpts, _requestTxHash, _from, _to, _tokenAddress, _value, _requestNonce, _requestBlockNumber, _extraData)
}

// HandleERC20Transfer is a paid mutator transaction binding the contract method 0x452fbc62.
//
// Solidity: function handleERC20Transfer(_requestTxHash bytes32, _from address, _to address, _tokenAddress address, _value uint256, _requestNonce uint64, _requestBlockNumber uint64, _extraData uint256[]) returns()
func (_ExtBridge *ExtBridgeTransactorSession) HandleERC20Transfer(_requestTxHash [32]byte, _from common.Address, _to common.Address, _tokenAddress common.Address, _value *big.Int, _requestNonce uint64, _requestBlockNumber uint64, _extraData []*big.Int) (*types.Transaction, error) {
	return _ExtBridge.Contract.HandleERC20Transfer(&_ExtBridge.TransactOpts, _requestTxHash, _from, _to, _tokenAddress, _value, _requestNonce, _requestBlockNumber, _extraData)
}

// HandleERC721Transfer is a paid mutator transaction binding the contract method 0x2e4bb197.
//
// Solidity: function handleERC721Transfer(_requestTxHash bytes32, _from address, _to address, _tokenAddress address, _tokenId uint256, _requestNonce uint64, _requestBlockNumber uint64, _tokenURI string, _extraData uint256[]) returns()
func (_ExtBridge *ExtBridgeTransactor) HandleERC721Transfer(opts *bind.TransactOpts, _requestTxHash [32]byte, _from common.Address, _to common.Address, _tokenAddress common.Address, _tokenId *big.Int, _requestNonce uint64, _requestBlockNumber uint64, _tokenURI string, _extraData []*big.Int) (*types.Transaction, error) {
	return _ExtBridge.contract.Transact(opts, "handleERC721Transfer", _requestTxHash, _from, _to, _tokenAddress, _tokenId, _requestNonce, _requestBlockNumber, _tokenURI, _extraData)
}

// HandleERC721Transfer is a paid mutator transaction binding the contract method 0x2e4bb197.
//
// Solidity: function handleERC721Transfer(_requestTxHash bytes32, _from address, _to address, _tokenAddress address, _tokenId uint256, _requestNonce uint64, _requestBlockNumber uint64, _tokenURI string, _extraData uint256[]) returns()
func (_ExtBridge *ExtBridgeSession) HandleERC721Transfer(_requestTxHash [32]byte, _from common.Address, _to common.Address, _tokenAddress common.Address, _tokenId *big.Int, _requestNonce uint64, _requestBlockNumber uint64, _tokenURI string, _extraData []*big.Int) (*types.Transaction, error) {
	return _ExtBridge.Contract.HandleERC721Transfer(&_ExtBridge.TransactOpts, _requestTxHash, _from, _to, _tokenAddress, _tokenId, _requestNonce, _requestBlockNumber, _tokenURI, _extraData)
}

// HandleERC721Transfer is a paid mutator transaction binding the contract method 0x2e4bb197.
//
// Solidity: function handleERC721Transfer(_requestTxHash bytes32, _from address, _to address, _tokenAddress address, _tokenId uint256, _requestNonce uint64, _requestBlockNumber uint64, _tokenURI string, _extraData uint256[]) returns()
func (_ExtBridge *ExtBridgeTransactorSession) HandleERC721Transfer(_requestTxHash [32]byte, _from common.Address, _to common.Address, _tokenAddress common.Address, _tokenId *big.Int, _requestNonce uint64, _requestBlockNumber uint64, _tokenURI string, _extraData []*big.Int) (*types.Transaction, error) {
	return _ExtBridge.Contract.HandleERC721Transfer(&_ExtBridge.TransactOpts, _requestTxHash, _from, _to, _tokenAddress, _tokenId, _requestNonce, _requestBlockNumber, _tokenURI, _extraData)
}

// HandleKLAYTransfer is a paid mutator transaction binding the contract method 0xee273ad7.
//
// Solidity: function handleKLAYTransfer(_requestTxHash bytes32, _from address, _to address, _value uint256, _requestedNonce uint64, _requestedBlockNumber uint64, _extraData uint256[]) returns()
func (_ExtBridge *ExtBridgeTransactor) HandleKLAYTransfer(opts *bind.TransactOpts, _requestTxHash [32]byte, _from common.Address, _to common.Address, _value *big.Int, _requestedNonce uint64, _requestedBlockNumber uint64, _extraData []*big.Int) (*types.Transaction, error) {
	return _ExtBridge.contract.Transact(opts, "handleKLAYTransfer", _requestTxHash, _from, _to, _value, _requestedNonce, _requestedBlockNumber, _extraData)
}

// HandleKLAYTransfer is a paid mutator transaction binding the contract method 0xee273ad7.
//
// Solidity: function handleKLAYTransfer(_requestTxHash bytes32, _from address, _to address, _value uint256, _requestedNonce uint64, _requestedBlockNumber uint64, _extraData uint256[]) returns()
func (_ExtBridge *ExtBridgeSession) HandleKLAYTransfer(_requestTxHash [32]byte, _from common.Address, _to common.Address, _value *big.Int, _requestedNonce uint64, _requestedBlockNumber uint64, _extraData []*big.Int) (*types.Transaction, error) {
	return _ExtBridge.Contract.HandleKLAYTransfer(&_ExtBridge.TransactOpts, _requestTxHash, _from, _to, _value, _requestedNonce, _requestedBlockNumber, _extraData)
}

// HandleKLAYTransfer is a paid mutator transaction binding the contract method 0xee273ad7.
//
// Solidity: function handleKLAYTransfer(_requestTxHash bytes32, _from address, _to address, _value uint256, _requestedNonce uint64, _requestedBlockNumber uint64, _extraData uint256[]) returns()
func (_ExtBridge *ExtBridgeTransactorSession) HandleKLAYTransfer(_requestTxHash [32]byte, _from common.Address, _to common.Address, _value *big.Int, _requestedNonce uint64, _requestedBlockNumber uint64, _extraData []*big.Int) (*types.Transaction, error) {
	return _ExtBridge.Contract.HandleKLAYTransfer(&_ExtBridge.TransactOpts, _requestTxHash, _from, _to, _value, _requestedNonce, _requestedBlockNumber, _extraData)
}

// OnERC20Received is a paid mutator transaction binding the contract method 0x349449b2.
//
// Solidity: function onERC20Received(_from address, _value uint256, _to address, _feeLimit uint256, _extraData uint256[]) returns()
func (_ExtBridge *ExtBridgeTransactor) OnERC20Received(opts *bind.TransactOpts, _from common.Address, _value *big.Int, _to common.Address, _feeLimit *big.Int, _extraData []*big.Int) (*types.Transaction, error) {
	return _ExtBridge.contract.Transact(opts, "onERC20Received", _from, _value, _to, _feeLimit, _extraData)
}

// OnERC20Received is a paid mutator transaction binding the contract method 0x349449b2.
//
// Solidity: function onERC20Received(_from address, _value uint256, _to address, _feeLimit uint256, _extraData uint256[]) returns()
func (_ExtBridge *ExtBridgeSession) OnERC20Received(_from common.Address, _value *big.Int, _to common.Address, _feeLimit *big.Int, _extraData []*big.Int) (*types.Transaction, error) {
	return _ExtBridge.Contract.OnERC20Received(&_ExtBridge.TransactOpts, _from, _value, _to, _feeLimit, _extraData)
}

// OnERC20Received is a paid mutator transaction binding the contract method 0x349449b2.
//
// Solidity: function onERC20Received(_from address, _value uint256, _to address, _feeLimit uint256, _extraData uint256[]) returns()
func (_ExtBridge *ExtBridgeTransactorSession) OnERC20Received(_from common.Address, _value *big.Int, _to common.Address, _feeLimit *big.Int, _extraData []*big.Int) (*types.Transaction, error) {
	return _ExtBridge.Contract.OnERC20Received(&_ExtBridge.TransactOpts, _from, _value, _to, _feeLimit, _extraData)
}

// OnERC721Received is a paid mutator transaction binding the contract method 0xa319de9e.
//
// Solidity: function onERC721Received(_from address, _tokenId uint256, _to address, _extraData uint256[]) returns()
func (_ExtBridge *ExtBridgeTransactor) OnERC721Received(opts *bind.TransactOpts, _from common.Address, _tokenId *big.Int, _to common.Address, _extraData []*big.Int) (*types.Transaction, error) {
	return _ExtBridge.contract.Transact(opts, "onERC721Received", _from, _tokenId, _to, _extraData)
}

// OnERC721Received is a paid mutator transaction binding the contract method 0xa319de9e.
//
// Solidity: function onERC721Received(_from address, _tokenId uint256, _to address, _extraData uint256[]) returns()
func (_ExtBridge *ExtBridgeSession) OnERC721Received(_from common.Address, _tokenId *big.Int, _to common.Address, _extraData []*big.Int) (*types.Transaction, error) {
	return _ExtBridge.Contract.OnERC721Received(&_ExtBridge.TransactOpts, _from, _tokenId, _to, _extraData)
}

// OnERC721Received is a paid mutator transaction binding the contract method 0xa319de9e.
//
// Solidity: function onERC721Received(_from address, _tokenId uint256, _to address, _extraData uint256[]) returns()
func (_ExtBridge *ExtBridgeTransactorSession) OnERC721Received(_from common.Address, _tokenId *big.Int, _to common.Address, _extraData []*big.Int) (*types.Transaction, error) {
	return _ExtBridge.Contract.OnERC721Received(&_ExtBridge.TransactOpts, _from, _tokenId, _to, _extraData)
}

// RegisterOperator is a paid mutator transaction binding the contract method 0x3682a450.
//
// Solidity: function registerOperator(_operator address) returns()
func (_ExtBridge *ExtBridgeTransactor) RegisterOperator(opts *bind.TransactOpts, _operator common.Address) (*types.Transaction, error) {
	return _ExtBridge.contract.Transact(opts, "registerOperator", _operator)
}

// RegisterOperator is a paid mutator transaction binding the contract method 0x3682a450.
//
// Solidity: function registerOperator(_operator address) returns()
func (_ExtBridge *ExtBridgeSession) RegisterOperator(_operator common.Address) (*types.Transaction, error) {
	return _ExtBridge.Contract.RegisterOperator(&_ExtBridge.TransactOpts, _operator)
}

// RegisterOperator is a paid mutator transaction binding the contract method 0x3682a450.
//
// Solidity: function registerOperator(_operator address) returns()
func (_ExtBridge *ExtBridgeTransactorSession) RegisterOperator(_operator common.Address) (*types.Transaction, error) {
	return _ExtBridge.Contract.RegisterOperator(&_ExtBridge.TransactOpts, _operator)
}

// RegisterToken is a paid mutator transaction binding the contract method 0x4739f7e5.
//
// Solidity: function registerToken(_token address, _cToken address) returns()
func (_ExtBridge *ExtBridgeTransactor) RegisterToken(opts *bind.TransactOpts, _token common.Address, _cToken common.Address) (*types.Transaction, error) {
	return _ExtBridge.contract.Transact(opts, "registerToken", _token, _cToken)
}

// RegisterToken is a paid mutator transaction binding the contract method 0x4739f7e5.
//
// Solidity: function registerToken(_token address, _cToken address) returns()
func (_ExtBridge *ExtBridgeSession) RegisterToken(_token common.Address, _cToken common.Address) (*types.Transaction, error) {
	return _ExtBridge.Contract.RegisterToken(&_ExtBridge.TransactOpts, _token, _cToken)
}

// RegisterToken is a paid mutator transaction binding the contract method 0x4739f7e5.
//
// Solidity: function registerToken(_token address, _cToken address) returns()
func (_ExtBridge *ExtBridgeTransactorSession) RegisterToken(_token common.Address, _cToken common.Address) (*types.Transaction, error) {
	return _ExtBridge.Contract.RegisterToken(&_ExtBridge.TransactOpts, _token, _cToken)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ExtBridge *ExtBridgeTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ExtBridge.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ExtBridge *ExtBridgeSession) RenounceOwnership() (*types.Transaction, error) {
	return _ExtBridge.Contract.RenounceOwnership(&_ExtBridge.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ExtBridge *ExtBridgeTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _ExtBridge.Contract.RenounceOwnership(&_ExtBridge.TransactOpts)
}

// RequestERC20Transfer is a paid mutator transaction binding the contract method 0x74f41f15.
//
// Solidity: function requestERC20Transfer(_tokenAddress address, _to address, _value uint256, _feeLimit uint256, _extraData uint256[]) returns()
func (_ExtBridge *ExtBridgeTransactor) RequestERC20Transfer(opts *bind.TransactOpts, _tokenAddress common.Address, _to common.Address, _value *big.Int, _feeLimit *big.Int, _extraData []*big.Int) (*types.Transaction, error) {
	return _ExtBridge.contract.Transact(opts, "requestERC20Transfer", _tokenAddress, _to, _value, _feeLimit, _extraData)
}

// RequestERC20Transfer is a paid mutator transaction binding the contract method 0x74f41f15.
//
// Solidity: function requestERC20Transfer(_tokenAddress address, _to address, _value uint256, _feeLimit uint256, _extraData uint256[]) returns()
func (_ExtBridge *ExtBridgeSession) RequestERC20Transfer(_tokenAddress common.Address, _to common.Address, _value *big.Int, _feeLimit *big.Int, _extraData []*big.Int) (*types.Transaction, error) {
	return _ExtBridge.Contract.RequestERC20Transfer(&_ExtBridge.TransactOpts, _tokenAddress, _to, _value, _feeLimit, _extraData)
}

// RequestERC20Transfer is a paid mutator transaction binding the contract method 0x74f41f15.
//
// Solidity: function requestERC20Transfer(_tokenAddress address, _to address, _value uint256, _feeLimit uint256, _extraData uint256[]) returns()
func (_ExtBridge *ExtBridgeTransactorSession) RequestERC20Transfer(_tokenAddress common.Address, _to common.Address, _value *big.Int, _feeLimit *big.Int, _extraData []*big.Int) (*types.Transaction, error) {
	return _ExtBridge.Contract.RequestERC20Transfer(&_ExtBridge.TransactOpts, _tokenAddress, _to, _value, _feeLimit, _extraData)
}

// RequestERC721Transfer is a paid mutator transaction binding the contract method 0x098540b8.
//
// Solidity: function requestERC721Transfer(_tokenAddress address, _to address, _tokenId uint256, _extraData uint256[]) returns()
func (_ExtBridge *ExtBridgeTransactor) RequestERC721Transfer(opts *bind.TransactOpts, _tokenAddress common.Address, _to common.Address, _tokenId *big.Int, _extraData []*big.Int) (*types.Transaction, error) {
	return _ExtBridge.contract.Transact(opts, "requestERC721Transfer", _tokenAddress, _to, _tokenId, _extraData)
}

// RequestERC721Transfer is a paid mutator transaction binding the contract method 0x098540b8.
//
// Solidity: function requestERC721Transfer(_tokenAddress address, _to address, _tokenId uint256, _extraData uint256[]) returns()
func (_ExtBridge *ExtBridgeSession) RequestERC721Transfer(_tokenAddress common.Address, _to common.Address, _tokenId *big.Int, _extraData []*big.Int) (*types.Transaction, error) {
	return _ExtBridge.Contract.RequestERC721Transfer(&_ExtBridge.TransactOpts, _tokenAddress, _to, _tokenId, _extraData)
}

// RequestERC721Transfer is a paid mutator transaction binding the contract method 0x098540b8.
//
// Solidity: function requestERC721Transfer(_tokenAddress address, _to address, _tokenId uint256, _extraData uint256[]) returns()
func (_ExtBridge *ExtBridgeTransactorSession) RequestERC721Transfer(_tokenAddress common.Address, _to common.Address, _tokenId *big.Int, _extraData []*big.Int) (*types.Transaction, error) {
	return _ExtBridge.Contract.RequestERC721Transfer(&_ExtBridge.TransactOpts, _tokenAddress, _to, _tokenId, _extraData)
}

// RequestKLAYTransfer is a paid mutator transaction binding the contract method 0xeb7f0cbe.
//
// Solidity: function requestKLAYTransfer(_to address, _value uint256, _extraData uint256[]) returns()
func (_ExtBridge *ExtBridgeTransactor) RequestKLAYTransfer(opts *bind.TransactOpts, _to common.Address, _value *big.Int, _extraData []*big.Int) (*types.Transaction, error) {
	return _ExtBridge.contract.Transact(opts, "requestKLAYTransfer", _to, _value, _extraData)
}

// RequestKLAYTransfer is a paid mutator transaction binding the contract method 0xeb7f0cbe.
//
// Solidity: function requestKLAYTransfer(_to address, _value uint256, _extraData uint256[]) returns()
func (_ExtBridge *ExtBridgeSession) RequestKLAYTransfer(_to common.Address, _value *big.Int, _extraData []*big.Int) (*types.Transaction, error) {
	return _ExtBridge.Contract.RequestKLAYTransfer(&_ExtBridge.TransactOpts, _to, _value, _extraData)
}

// RequestKLAYTransfer is a paid mutator transaction binding the contract method 0xeb7f0cbe.
//
// Solidity: function requestKLAYTransfer(_to address, _value uint256, _extraData uint256[]) returns()
func (_ExtBridge *ExtBridgeTransactorSession) RequestKLAYTransfer(_to common.Address, _value *big.Int, _extraData []*big.Int) (*types.Transaction, error) {
	return _ExtBridge.Contract.RequestKLAYTransfer(&_ExtBridge.TransactOpts, _to, _value, _extraData)
}

// SetCallback is a paid mutator transaction binding the contract method 0x8daa63ac.
//
// Solidity: function setCallback(_addr address) returns()
func (_ExtBridge *ExtBridgeTransactor) SetCallback(opts *bind.TransactOpts, _addr common.Address) (*types.Transaction, error) {
	return _ExtBridge.contract.Transact(opts, "setCallback", _addr)
}

// SetCallback is a paid mutator transaction binding the contract method 0x8daa63ac.
//
// Solidity: function setCallback(_addr address) returns()
func (_ExtBridge *ExtBridgeSession) SetCallback(_addr common.Address) (*types.Transaction, error) {
	return _ExtBridge.Contract.SetCallback(&_ExtBridge.TransactOpts, _addr)
}

// SetCallback is a paid mutator transaction binding the contract method 0x8daa63ac.
//
// Solidity: function setCallback(_addr address) returns()
func (_ExtBridge *ExtBridgeTransactorSession) SetCallback(_addr common.Address) (*types.Transaction, error) {
	return _ExtBridge.Contract.SetCallback(&_ExtBridge.TransactOpts, _addr)
}

// SetCounterPartBridge is a paid mutator transaction binding the contract method 0x87b04c55.
//
// Solidity: function setCounterPartBridge(_bridge address) returns()
func (_ExtBridge *ExtBridgeTransactor) SetCounterPartBridge(opts *bind.TransactOpts, _bridge common.Address) (*types.Transaction, error) {
	return _ExtBridge.contract.Transact(opts, "setCounterPartBridge", _bridge)
}

// SetCounterPartBridge is a paid mutator transaction binding the contract method 0x87b04c55.
//
// Solidity: function setCounterPartBridge(_bridge address) returns()
func (_ExtBridge *ExtBridgeSession) SetCounterPartBridge(_bridge common.Address) (*types.Transaction, error) {
	return _ExtBridge.Contract.SetCounterPartBridge(&_ExtBridge.TransactOpts, _bridge)
}

// SetCounterPartBridge is a paid mutator transaction binding the contract method 0x87b04c55.
//
// Solidity: function setCounterPartBridge(_bridge address) returns()
func (_ExtBridge *ExtBridgeTransactorSession) SetCounterPartBridge(_bridge common.Address) (*types.Transaction, error) {
	return _ExtBridge.Contract.SetCounterPartBridge(&_ExtBridge.TransactOpts, _bridge)
}

// SetERC20Fee is a paid mutator transaction binding the contract method 0x2f88396c.
//
// Solidity: function setERC20Fee(_token address, _fee uint256, _requestNonce uint64) returns()
func (_ExtBridge *ExtBridgeTransactor) SetERC20Fee(opts *bind.TransactOpts, _token common.Address, _fee *big.Int, _requestNonce uint64) (*types.Transaction, error) {
	return _ExtBridge.contract.Transact(opts, "setERC20Fee", _token, _fee, _requestNonce)
}

// SetERC20Fee is a paid mutator transaction binding the contract method 0x2f88396c.
//
// Solidity: function setERC20Fee(_token address, _fee uint256, _requestNonce uint64) returns()
func (_ExtBridge *ExtBridgeSession) SetERC20Fee(_token common.Address, _fee *big.Int, _requestNonce uint64) (*types.Transaction, error) {
	return _ExtBridge.Contract.SetERC20Fee(&_ExtBridge.TransactOpts, _token, _fee, _requestNonce)
}

// SetERC20Fee is a paid mutator transaction binding the contract method 0x2f88396c.
//
// Solidity: function setERC20Fee(_token address, _fee uint256, _requestNonce uint64) returns()
func (_ExtBridge *ExtBridgeTransactorSession) SetERC20Fee(_token common.Address, _fee *big.Int, _requestNonce uint64) (*types.Transaction, error) {
	return _ExtBridge.Contract.SetERC20Fee(&_ExtBridge.TransactOpts, _token, _fee, _requestNonce)
}

// SetFeeReceiver is a paid mutator transaction binding the contract method 0xefdcd974.
//
// Solidity: function setFeeReceiver(_feeReceiver address) returns()
func (_ExtBridge *ExtBridgeTransactor) SetFeeReceiver(opts *bind.TransactOpts, _feeReceiver common.Address) (*types.Transaction, error) {
	return _ExtBridge.contract.Transact(opts, "setFeeReceiver", _feeReceiver)
}

// SetFeeReceiver is a paid mutator transaction binding the contract method 0xefdcd974.
//
// Solidity: function setFeeReceiver(_feeReceiver address) returns()
func (_ExtBridge *ExtBridgeSession) SetFeeReceiver(_feeReceiver common.Address) (*types.Transaction, error) {
	return _ExtBridge.Contract.SetFeeReceiver(&_ExtBridge.TransactOpts, _feeReceiver)
}

// SetFeeReceiver is a paid mutator transaction binding the contract method 0xefdcd974.
//
// Solidity: function setFeeReceiver(_feeReceiver address) returns()
func (_ExtBridge *ExtBridgeTransactorSession) SetFeeReceiver(_feeReceiver common.Address) (*types.Transaction, error) {
	return _ExtBridge.Contract.SetFeeReceiver(&_ExtBridge.TransactOpts, _feeReceiver)
}

// SetKLAYFee is a paid mutator transaction binding the contract method 0x1a2ae53e.
//
// Solidity: function setKLAYFee(_fee uint256, _requestNonce uint64) returns()
func (_ExtBridge *ExtBridgeTransactor) SetKLAYFee(opts *bind.TransactOpts, _fee *big.Int, _requestNonce uint64) (*types.Transaction, error) {
	return _ExtBridge.contract.Transact(opts, "setKLAYFee", _fee, _requestNonce)
}

// SetKLAYFee is a paid mutator transaction binding the contract method 0x1a2ae53e.
//
// Solidity: function setKLAYFee(_fee uint256, _requestNonce uint64) returns()
func (_ExtBridge *ExtBridgeSession) SetKLAYFee(_fee *big.Int, _requestNonce uint64) (*types.Transaction, error) {
	return _ExtBridge.Contract.SetKLAYFee(&_ExtBridge.TransactOpts, _fee, _requestNonce)
}

// SetKLAYFee is a paid mutator transaction binding the contract method 0x1a2ae53e.
//
// Solidity: function setKLAYFee(_fee uint256, _requestNonce uint64) returns()
func (_ExtBridge *ExtBridgeTransactorSession) SetKLAYFee(_fee *big.Int, _requestNonce uint64) (*types.Transaction, error) {
	return _ExtBridge.Contract.SetKLAYFee(&_ExtBridge.TransactOpts, _fee, _requestNonce)
}

// SetOperatorThreshold is a paid mutator transaction binding the contract method 0xee2aec65.
//
// Solidity: function setOperatorThreshold(_voteType uint8, _threshold uint8) returns()
func (_ExtBridge *ExtBridgeTransactor) SetOperatorThreshold(opts *bind.TransactOpts, _voteType uint8, _threshold uint8) (*types.Transaction, error) {
	return _ExtBridge.contract.Transact(opts, "setOperatorThreshold", _voteType, _threshold)
}

// SetOperatorThreshold is a paid mutator transaction binding the contract method 0xee2aec65.
//
// Solidity: function setOperatorThreshold(_voteType uint8, _threshold uint8) returns()
func (_ExtBridge *ExtBridgeSession) SetOperatorThreshold(_voteType uint8, _threshold uint8) (*types.Transaction, error) {
	return _ExtBridge.Contract.SetOperatorThreshold(&_ExtBridge.TransactOpts, _voteType, _threshold)
}

// SetOperatorThreshold is a paid mutator transaction binding the contract method 0xee2aec65.
//
// Solidity: function setOperatorThreshold(_voteType uint8, _threshold uint8) returns()
func (_ExtBridge *ExtBridgeTransactorSession) SetOperatorThreshold(_voteType uint8, _threshold uint8) (*types.Transaction, error) {
	return _ExtBridge.Contract.SetOperatorThreshold(&_ExtBridge.TransactOpts, _voteType, _threshold)
}

// Start is a paid mutator transaction binding the contract method 0xc877cf37.
//
// Solidity: function start(_status bool) returns()
func (_ExtBridge *ExtBridgeTransactor) Start(opts *bind.TransactOpts, _status bool) (*types.Transaction, error) {
	return _ExtBridge.contract.Transact(opts, "start", _status)
}

// Start is a paid mutator transaction binding the contract method 0xc877cf37.
//
// Solidity: function start(_status bool) returns()
func (_ExtBridge *ExtBridgeSession) Start(_status bool) (*types.Transaction, error) {
	return _ExtBridge.Contract.Start(&_ExtBridge.TransactOpts, _status)
}

// Start is a paid mutator transaction binding the contract method 0xc877cf37.
//
// Solidity: function start(_status bool) returns()
func (_ExtBridge *ExtBridgeTransactorSession) Start(_status bool) (*types.Transaction, error) {
	return _ExtBridge.Contract.Start(&_ExtBridge.TransactOpts, _status)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(newOwner address) returns()
func (_ExtBridge *ExtBridgeTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _ExtBridge.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(newOwner address) returns()
func (_ExtBridge *ExtBridgeSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ExtBridge.Contract.TransferOwnership(&_ExtBridge.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(newOwner address) returns()
func (_ExtBridge *ExtBridgeTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ExtBridge.Contract.TransferOwnership(&_ExtBridge.TransactOpts, newOwner)
}

// ExtBridgeERC20FeeChangedIterator is returned from FilterERC20FeeChanged and is used to iterate over the raw logs and unpacked data for ERC20FeeChanged events raised by the ExtBridge contract.
type ExtBridgeERC20FeeChangedIterator struct {
	Event *ExtBridgeERC20FeeChanged // Event containing the contract specifics and raw log

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
func (it *ExtBridgeERC20FeeChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExtBridgeERC20FeeChanged)
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
		it.Event = new(ExtBridgeERC20FeeChanged)
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
func (it *ExtBridgeERC20FeeChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExtBridgeERC20FeeChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExtBridgeERC20FeeChanged represents a ERC20FeeChanged event raised by the ExtBridge contract.
type ExtBridgeERC20FeeChanged struct {
	Token common.Address
	Fee   *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterERC20FeeChanged is a free log retrieval operation binding the contract event 0xdb5ad2e76ae20cfa4e7adbc7305d7538442164d85ead9937c98620a1aa4c255b.
//
// Solidity: e ERC20FeeChanged(token address, fee indexed uint256)
func (_ExtBridge *ExtBridgeFilterer) FilterERC20FeeChanged(opts *bind.FilterOpts, fee []*big.Int) (*ExtBridgeERC20FeeChangedIterator, error) {

	var feeRule []interface{}
	for _, feeItem := range fee {
		feeRule = append(feeRule, feeItem)
	}

	logs, sub, err := _ExtBridge.contract.FilterLogs(opts, "ERC20FeeChanged", feeRule)
	if err != nil {
		return nil, err
	}
	return &ExtBridgeERC20FeeChangedIterator{contract: _ExtBridge.contract, event: "ERC20FeeChanged", logs: logs, sub: sub}, nil
}

// WatchERC20FeeChanged is a free log subscription operation binding the contract event 0xdb5ad2e76ae20cfa4e7adbc7305d7538442164d85ead9937c98620a1aa4c255b.
//
// Solidity: e ERC20FeeChanged(token address, fee indexed uint256)
func (_ExtBridge *ExtBridgeFilterer) WatchERC20FeeChanged(opts *bind.WatchOpts, sink chan<- *ExtBridgeERC20FeeChanged, fee []*big.Int) (event.Subscription, error) {

	var feeRule []interface{}
	for _, feeItem := range fee {
		feeRule = append(feeRule, feeItem)
	}

	logs, sub, err := _ExtBridge.contract.WatchLogs(opts, "ERC20FeeChanged", feeRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExtBridgeERC20FeeChanged)
				if err := _ExtBridge.contract.UnpackLog(event, "ERC20FeeChanged", log); err != nil {
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

// ExtBridgeFeeReceiverChangedIterator is returned from FilterFeeReceiverChanged and is used to iterate over the raw logs and unpacked data for FeeReceiverChanged events raised by the ExtBridge contract.
type ExtBridgeFeeReceiverChangedIterator struct {
	Event *ExtBridgeFeeReceiverChanged // Event containing the contract specifics and raw log

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
func (it *ExtBridgeFeeReceiverChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExtBridgeFeeReceiverChanged)
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
		it.Event = new(ExtBridgeFeeReceiverChanged)
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
func (it *ExtBridgeFeeReceiverChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExtBridgeFeeReceiverChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExtBridgeFeeReceiverChanged represents a FeeReceiverChanged event raised by the ExtBridge contract.
type ExtBridgeFeeReceiverChanged struct {
	FeeReceiver common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterFeeReceiverChanged is a free log retrieval operation binding the contract event 0x647672599d3468abcfa241a13c9e3d34383caadb5cc80fb67c3cdfcd5f786059.
//
// Solidity: e FeeReceiverChanged(feeReceiver indexed address)
func (_ExtBridge *ExtBridgeFilterer) FilterFeeReceiverChanged(opts *bind.FilterOpts, feeReceiver []common.Address) (*ExtBridgeFeeReceiverChangedIterator, error) {

	var feeReceiverRule []interface{}
	for _, feeReceiverItem := range feeReceiver {
		feeReceiverRule = append(feeReceiverRule, feeReceiverItem)
	}

	logs, sub, err := _ExtBridge.contract.FilterLogs(opts, "FeeReceiverChanged", feeReceiverRule)
	if err != nil {
		return nil, err
	}
	return &ExtBridgeFeeReceiverChangedIterator{contract: _ExtBridge.contract, event: "FeeReceiverChanged", logs: logs, sub: sub}, nil
}

// WatchFeeReceiverChanged is a free log subscription operation binding the contract event 0x647672599d3468abcfa241a13c9e3d34383caadb5cc80fb67c3cdfcd5f786059.
//
// Solidity: e FeeReceiverChanged(feeReceiver indexed address)
func (_ExtBridge *ExtBridgeFilterer) WatchFeeReceiverChanged(opts *bind.WatchOpts, sink chan<- *ExtBridgeFeeReceiverChanged, feeReceiver []common.Address) (event.Subscription, error) {

	var feeReceiverRule []interface{}
	for _, feeReceiverItem := range feeReceiver {
		feeReceiverRule = append(feeReceiverRule, feeReceiverItem)
	}

	logs, sub, err := _ExtBridge.contract.WatchLogs(opts, "FeeReceiverChanged", feeReceiverRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExtBridgeFeeReceiverChanged)
				if err := _ExtBridge.contract.UnpackLog(event, "FeeReceiverChanged", log); err != nil {
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

// ExtBridgeHandleValueTransferIterator is returned from FilterHandleValueTransfer and is used to iterate over the raw logs and unpacked data for HandleValueTransfer events raised by the ExtBridge contract.
type ExtBridgeHandleValueTransferIterator struct {
	Event *ExtBridgeHandleValueTransfer // Event containing the contract specifics and raw log

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
func (it *ExtBridgeHandleValueTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExtBridgeHandleValueTransfer)
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
		it.Event = new(ExtBridgeHandleValueTransfer)
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
func (it *ExtBridgeHandleValueTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExtBridgeHandleValueTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExtBridgeHandleValueTransfer represents a HandleValueTransfer event raised by the ExtBridge contract.
type ExtBridgeHandleValueTransfer struct {
	RequestTxHash  [32]byte
	TokenType      uint8
	From           common.Address
	To             common.Address
	TokenAddress   common.Address
	ValueOrTokenId *big.Int
	HandleNonce    uint64
	ExtraData      []*big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterHandleValueTransfer is a free log retrieval operation binding the contract event 0x91f1d4b9864c9c091ade654e58899919d2a56e194153899c7e292cd3ae2ac47b.
//
// Solidity: e HandleValueTransfer(requestTxHash bytes32, tokenType uint8, from address, to address, tokenAddress address, valueOrTokenId uint256, handleNonce uint64, extraData uint256[])
func (_ExtBridge *ExtBridgeFilterer) FilterHandleValueTransfer(opts *bind.FilterOpts) (*ExtBridgeHandleValueTransferIterator, error) {

	logs, sub, err := _ExtBridge.contract.FilterLogs(opts, "HandleValueTransfer")
	if err != nil {
		return nil, err
	}
	return &ExtBridgeHandleValueTransferIterator{contract: _ExtBridge.contract, event: "HandleValueTransfer", logs: logs, sub: sub}, nil
}

// WatchHandleValueTransfer is a free log subscription operation binding the contract event 0x91f1d4b9864c9c091ade654e58899919d2a56e194153899c7e292cd3ae2ac47b.
//
// Solidity: e HandleValueTransfer(requestTxHash bytes32, tokenType uint8, from address, to address, tokenAddress address, valueOrTokenId uint256, handleNonce uint64, extraData uint256[])
func (_ExtBridge *ExtBridgeFilterer) WatchHandleValueTransfer(opts *bind.WatchOpts, sink chan<- *ExtBridgeHandleValueTransfer) (event.Subscription, error) {

	logs, sub, err := _ExtBridge.contract.WatchLogs(opts, "HandleValueTransfer")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExtBridgeHandleValueTransfer)
				if err := _ExtBridge.contract.UnpackLog(event, "HandleValueTransfer", log); err != nil {
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

// ExtBridgeKLAYFeeChangedIterator is returned from FilterKLAYFeeChanged and is used to iterate over the raw logs and unpacked data for KLAYFeeChanged events raised by the ExtBridge contract.
type ExtBridgeKLAYFeeChangedIterator struct {
	Event *ExtBridgeKLAYFeeChanged // Event containing the contract specifics and raw log

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
func (it *ExtBridgeKLAYFeeChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExtBridgeKLAYFeeChanged)
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
		it.Event = new(ExtBridgeKLAYFeeChanged)
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
func (it *ExtBridgeKLAYFeeChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExtBridgeKLAYFeeChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExtBridgeKLAYFeeChanged represents a KLAYFeeChanged event raised by the ExtBridge contract.
type ExtBridgeKLAYFeeChanged struct {
	Fee *big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// FilterKLAYFeeChanged is a free log retrieval operation binding the contract event 0xa7a33d0996347e1aa55ca2206015b61b9534bdd881d59d59aa680e25eefac365.
//
// Solidity: e KLAYFeeChanged(fee indexed uint256)
func (_ExtBridge *ExtBridgeFilterer) FilterKLAYFeeChanged(opts *bind.FilterOpts, fee []*big.Int) (*ExtBridgeKLAYFeeChangedIterator, error) {

	var feeRule []interface{}
	for _, feeItem := range fee {
		feeRule = append(feeRule, feeItem)
	}

	logs, sub, err := _ExtBridge.contract.FilterLogs(opts, "KLAYFeeChanged", feeRule)
	if err != nil {
		return nil, err
	}
	return &ExtBridgeKLAYFeeChangedIterator{contract: _ExtBridge.contract, event: "KLAYFeeChanged", logs: logs, sub: sub}, nil
}

// WatchKLAYFeeChanged is a free log subscription operation binding the contract event 0xa7a33d0996347e1aa55ca2206015b61b9534bdd881d59d59aa680e25eefac365.
//
// Solidity: e KLAYFeeChanged(fee indexed uint256)
func (_ExtBridge *ExtBridgeFilterer) WatchKLAYFeeChanged(opts *bind.WatchOpts, sink chan<- *ExtBridgeKLAYFeeChanged, fee []*big.Int) (event.Subscription, error) {

	var feeRule []interface{}
	for _, feeItem := range fee {
		feeRule = append(feeRule, feeItem)
	}

	logs, sub, err := _ExtBridge.contract.WatchLogs(opts, "KLAYFeeChanged", feeRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExtBridgeKLAYFeeChanged)
				if err := _ExtBridge.contract.UnpackLog(event, "KLAYFeeChanged", log); err != nil {
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

// ExtBridgeOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the ExtBridge contract.
type ExtBridgeOwnershipTransferredIterator struct {
	Event *ExtBridgeOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *ExtBridgeOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExtBridgeOwnershipTransferred)
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
		it.Event = new(ExtBridgeOwnershipTransferred)
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
func (it *ExtBridgeOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExtBridgeOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExtBridgeOwnershipTransferred represents a OwnershipTransferred event raised by the ExtBridge contract.
type ExtBridgeOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: e OwnershipTransferred(previousOwner indexed address, newOwner indexed address)
func (_ExtBridge *ExtBridgeFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*ExtBridgeOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ExtBridge.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &ExtBridgeOwnershipTransferredIterator{contract: _ExtBridge.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: e OwnershipTransferred(previousOwner indexed address, newOwner indexed address)
func (_ExtBridge *ExtBridgeFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ExtBridgeOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ExtBridge.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExtBridgeOwnershipTransferred)
				if err := _ExtBridge.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// ExtBridgeRequestValueTransferIterator is returned from FilterRequestValueTransfer and is used to iterate over the raw logs and unpacked data for RequestValueTransfer events raised by the ExtBridge contract.
type ExtBridgeRequestValueTransferIterator struct {
	Event *ExtBridgeRequestValueTransfer // Event containing the contract specifics and raw log

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
func (it *ExtBridgeRequestValueTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExtBridgeRequestValueTransfer)
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
		it.Event = new(ExtBridgeRequestValueTransfer)
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
func (it *ExtBridgeRequestValueTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExtBridgeRequestValueTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExtBridgeRequestValueTransfer represents a RequestValueTransfer event raised by the ExtBridge contract.
type ExtBridgeRequestValueTransfer struct {
	TokenType      uint8
	From           common.Address
	To             common.Address
	TokenAddress   common.Address
	ValueOrTokenId *big.Int
	RequestNonce   uint64
	Uri            string
	Fee            *big.Int
	ExtraData      []*big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterRequestValueTransfer is a free log retrieval operation binding the contract event 0x09a722fc78123f7b2ba788ae83b01221477757cfb4b83302421e7c5c1fd01da9.
//
// Solidity: e RequestValueTransfer(tokenType uint8, from address, to address, tokenAddress address, valueOrTokenId uint256, requestNonce uint64, uri string, fee uint256, extraData uint256[])
func (_ExtBridge *ExtBridgeFilterer) FilterRequestValueTransfer(opts *bind.FilterOpts) (*ExtBridgeRequestValueTransferIterator, error) {

	logs, sub, err := _ExtBridge.contract.FilterLogs(opts, "RequestValueTransfer")
	if err != nil {
		return nil, err
	}
	return &ExtBridgeRequestValueTransferIterator{contract: _ExtBridge.contract, event: "RequestValueTransfer", logs: logs, sub: sub}, nil
}

// WatchRequestValueTransfer is a free log subscription operation binding the contract event 0x09a722fc78123f7b2ba788ae83b01221477757cfb4b83302421e7c5c1fd01da9.
//
// Solidity: e RequestValueTransfer(tokenType uint8, from address, to address, tokenAddress address, valueOrTokenId uint256, requestNonce uint64, uri string, fee uint256, extraData uint256[])
func (_ExtBridge *ExtBridgeFilterer) WatchRequestValueTransfer(opts *bind.WatchOpts, sink chan<- *ExtBridgeRequestValueTransfer) (event.Subscription, error) {

	logs, sub, err := _ExtBridge.contract.WatchLogs(opts, "RequestValueTransfer")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExtBridgeRequestValueTransfer)
				if err := _ExtBridge.contract.UnpackLog(event, "RequestValueTransfer", log); err != nil {
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

// HandledRequestsABI is the input ABI used to generate the binding from.
const HandledRequestsABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"handledRequestTx\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"

// HandledRequestsBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const HandledRequestsBinRuntime = `0x608060405260043610603e5763ffffffff7c01000000000000000000000000000000000000000000000000000000006000350416638a75eee281146043575b600080fd5b348015604e57600080fd5b506058600435606c565b604080519115158252519081900360200190f35b60006020819052908152604090205460ff16815600a165627a7a723058201bbae6a7d37dd53d64e287b454677e6fe47980c1d565a44b026a0030a12a827d0029`

// HandledRequestsBin is the compiled bytecode used for deploying new contracts.
const HandledRequestsBin = `0x6080604052348015600f57600080fd5b5060ad8061001e6000396000f300608060405260043610603e5763ffffffff7c01000000000000000000000000000000000000000000000000000000006000350416638a75eee281146043575b600080fd5b348015604e57600080fd5b506058600435606c565b604080519115158252519081900360200190f35b60006020819052908152604090205460ff16815600a165627a7a723058201bbae6a7d37dd53d64e287b454677e6fe47980c1d565a44b026a0030a12a827d0029`

// DeployHandledRequests deploys a new Klaytn contract, binding an instance of HandledRequests to it.
func DeployHandledRequests(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *HandledRequests, error) {
	parsed, err := abi.JSON(strings.NewReader(HandledRequestsABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(HandledRequestsBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &HandledRequests{HandledRequestsCaller: HandledRequestsCaller{contract: contract}, HandledRequestsTransactor: HandledRequestsTransactor{contract: contract}, HandledRequestsFilterer: HandledRequestsFilterer{contract: contract}}, nil
}

// HandledRequests is an auto generated Go binding around a Klaytn contract.
type HandledRequests struct {
	HandledRequestsCaller     // Read-only binding to the contract
	HandledRequestsTransactor // Write-only binding to the contract
	HandledRequestsFilterer   // Log filterer for contract events
}

// HandledRequestsCaller is an auto generated read-only Go binding around a Klaytn contract.
type HandledRequestsCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// HandledRequestsTransactor is an auto generated write-only Go binding around a Klaytn contract.
type HandledRequestsTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// HandledRequestsFilterer is an auto generated log filtering Go binding around a Klaytn contract events.
type HandledRequestsFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// HandledRequestsSession is an auto generated Go binding around a Klaytn contract,
// with pre-set call and transact options.
type HandledRequestsSession struct {
	Contract     *HandledRequests  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// HandledRequestsCallerSession is an auto generated read-only Go binding around a Klaytn contract,
// with pre-set call options.
type HandledRequestsCallerSession struct {
	Contract *HandledRequestsCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// HandledRequestsTransactorSession is an auto generated write-only Go binding around a Klaytn contract,
// with pre-set transact options.
type HandledRequestsTransactorSession struct {
	Contract     *HandledRequestsTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// HandledRequestsRaw is an auto generated low-level Go binding around a Klaytn contract.
type HandledRequestsRaw struct {
	Contract *HandledRequests // Generic contract binding to access the raw methods on
}

// HandledRequestsCallerRaw is an auto generated low-level read-only Go binding around a Klaytn contract.
type HandledRequestsCallerRaw struct {
	Contract *HandledRequestsCaller // Generic read-only contract binding to access the raw methods on
}

// HandledRequestsTransactorRaw is an auto generated low-level write-only Go binding around a Klaytn contract.
type HandledRequestsTransactorRaw struct {
	Contract *HandledRequestsTransactor // Generic write-only contract binding to access the raw methods on
}

// NewHandledRequests creates a new instance of HandledRequests, bound to a specific deployed contract.
func NewHandledRequests(address common.Address, backend bind.ContractBackend) (*HandledRequests, error) {
	contract, err := bindHandledRequests(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &HandledRequests{HandledRequestsCaller: HandledRequestsCaller{contract: contract}, HandledRequestsTransactor: HandledRequestsTransactor{contract: contract}, HandledRequestsFilterer: HandledRequestsFilterer{contract: contract}}, nil
}

// NewHandledRequestsCaller creates a new read-only instance of HandledRequests, bound to a specific deployed contract.
func NewHandledRequestsCaller(address common.Address, caller bind.ContractCaller) (*HandledRequestsCaller, error) {
	contract, err := bindHandledRequests(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &HandledRequestsCaller{contract: contract}, nil
}

// NewHandledRequestsTransactor creates a new write-only instance of HandledRequests, bound to a specific deployed contract.
func NewHandledRequestsTransactor(address common.Address, transactor bind.ContractTransactor) (*HandledRequestsTransactor, error) {
	contract, err := bindHandledRequests(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &HandledRequestsTransactor{contract: contract}, nil
}

// NewHandledRequestsFilterer creates a new log filterer instance of HandledRequests, bound to a specific deployed contract.
func NewHandledRequestsFilterer(address common.Address, filterer bind.ContractFilterer) (*HandledRequestsFilterer, error) {
	contract, err := bindHandledRequests(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &HandledRequestsFilterer{contract: contract}, nil
}

// bindHandledRequests binds a generic wrapper to an already deployed contract.
func bindHandledRequests(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(HandledRequestsABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_HandledRequests *HandledRequestsRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _HandledRequests.Contract.HandledRequestsCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_HandledRequests *HandledRequestsRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _HandledRequests.Contract.HandledRequestsTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_HandledRequests *HandledRequestsRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _HandledRequests.Contract.HandledRequestsTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_HandledRequests *HandledRequestsCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _HandledRequests.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_HandledRequests *HandledRequestsTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _HandledRequests.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_HandledRequests *HandledRequestsTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _HandledRequests.Contract.contract.Transact(opts, method, params...)
}

// HandledRequestTx is a free data retrieval call binding the contract method 0x8a75eee2.
//
// Solidity: function handledRequestTx( bytes32) constant returns(bool)
func (_HandledRequests *HandledRequestsCaller) HandledRequestTx(opts *bind.CallOpts, arg0 [32]byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _HandledRequests.contract.Call(opts, out, "handledRequestTx", arg0)
	return *ret0, err
}

// HandledRequestTx is a free data retrieval call binding the contract method 0x8a75eee2.
//
// Solidity: function handledRequestTx( bytes32) constant returns(bool)
func (_HandledRequests *HandledRequestsSession) HandledRequestTx(arg0 [32]byte) (bool, error) {
	return _HandledRequests.Contract.HandledRequestTx(&_HandledRequests.CallOpts, arg0)
}

// HandledRequestTx is a free data retrieval call binding the contract method 0x8a75eee2.
//
// Solidity: function handledRequestTx( bytes32) constant returns(bool)
func (_HandledRequests *HandledRequestsCallerSession) HandledRequestTx(arg0 [32]byte) (bool, error) {
	return _HandledRequests.Contract.HandledRequestTx(&_HandledRequests.CallOpts, arg0)
}

// IERC165ABI is the input ABI used to generate the binding from.
const IERC165ABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"

// IERC165BinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const IERC165BinRuntime = `0x`

// IERC165Bin is the compiled bytecode used for deploying new contracts.
const IERC165Bin = `0x`

// DeployIERC165 deploys a new Klaytn contract, binding an instance of IERC165 to it.
func DeployIERC165(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *IERC165, error) {
	parsed, err := abi.JSON(strings.NewReader(IERC165ABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(IERC165Bin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &IERC165{IERC165Caller: IERC165Caller{contract: contract}, IERC165Transactor: IERC165Transactor{contract: contract}, IERC165Filterer: IERC165Filterer{contract: contract}}, nil
}

// IERC165 is an auto generated Go binding around a Klaytn contract.
type IERC165 struct {
	IERC165Caller     // Read-only binding to the contract
	IERC165Transactor // Write-only binding to the contract
	IERC165Filterer   // Log filterer for contract events
}

// IERC165Caller is an auto generated read-only Go binding around a Klaytn contract.
type IERC165Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IERC165Transactor is an auto generated write-only Go binding around a Klaytn contract.
type IERC165Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IERC165Filterer is an auto generated log filtering Go binding around a Klaytn contract events.
type IERC165Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IERC165Session is an auto generated Go binding around a Klaytn contract,
// with pre-set call and transact options.
type IERC165Session struct {
	Contract     *IERC165          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IERC165CallerSession is an auto generated read-only Go binding around a Klaytn contract,
// with pre-set call options.
type IERC165CallerSession struct {
	Contract *IERC165Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// IERC165TransactorSession is an auto generated write-only Go binding around a Klaytn contract,
// with pre-set transact options.
type IERC165TransactorSession struct {
	Contract     *IERC165Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// IERC165Raw is an auto generated low-level Go binding around a Klaytn contract.
type IERC165Raw struct {
	Contract *IERC165 // Generic contract binding to access the raw methods on
}

// IERC165CallerRaw is an auto generated low-level read-only Go binding around a Klaytn contract.
type IERC165CallerRaw struct {
	Contract *IERC165Caller // Generic read-only contract binding to access the raw methods on
}

// IERC165TransactorRaw is an auto generated low-level write-only Go binding around a Klaytn contract.
type IERC165TransactorRaw struct {
	Contract *IERC165Transactor // Generic write-only contract binding to access the raw methods on
}

// NewIERC165 creates a new instance of IERC165, bound to a specific deployed contract.
func NewIERC165(address common.Address, backend bind.ContractBackend) (*IERC165, error) {
	contract, err := bindIERC165(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IERC165{IERC165Caller: IERC165Caller{contract: contract}, IERC165Transactor: IERC165Transactor{contract: contract}, IERC165Filterer: IERC165Filterer{contract: contract}}, nil
}

// NewIERC165Caller creates a new read-only instance of IERC165, bound to a specific deployed contract.
func NewIERC165Caller(address common.Address, caller bind.ContractCaller) (*IERC165Caller, error) {
	contract, err := bindIERC165(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IERC165Caller{contract: contract}, nil
}

// NewIERC165Transactor creates a new write-only instance of IERC165, bound to a specific deployed contract.
func NewIERC165Transactor(address common.Address, transactor bind.ContractTransactor) (*IERC165Transactor, error) {
	contract, err := bindIERC165(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IERC165Transactor{contract: contract}, nil
}

// NewIERC165Filterer creates a new log filterer instance of IERC165, bound to a specific deployed contract.
func NewIERC165Filterer(address common.Address, filterer bind.ContractFilterer) (*IERC165Filterer, error) {
	contract, err := bindIERC165(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IERC165Filterer{contract: contract}, nil
}

// bindIERC165 binds a generic wrapper to an already deployed contract.
func bindIERC165(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(IERC165ABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IERC165 *IERC165Raw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _IERC165.Contract.IERC165Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IERC165 *IERC165Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IERC165.Contract.IERC165Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IERC165 *IERC165Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IERC165.Contract.IERC165Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IERC165 *IERC165CallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _IERC165.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IERC165 *IERC165TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IERC165.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IERC165 *IERC165TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IERC165.Contract.contract.Transact(opts, method, params...)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(interfaceId bytes4) constant returns(bool)
func (_IERC165 *IERC165Caller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _IERC165.contract.Call(opts, out, "supportsInterface", interfaceId)
	return *ret0, err
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(interfaceId bytes4) constant returns(bool)
func (_IERC165 *IERC165Session) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _IERC165.Contract.SupportsInterface(&_IERC165.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(interfaceId bytes4) constant returns(bool)
func (_IERC165 *IERC165CallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _IERC165.Contract.SupportsInterface(&_IERC165.CallOpts, interfaceId)
}

// IERC20ABI is the input ABI used to generate the binding from.
const IERC20ABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"spender\",\"type\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"from\",\"type\":\"address\"},{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"who\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"owner\",\"type\":\"address\"},{\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"}]"

// IERC20BinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const IERC20BinRuntime = `0x`

// IERC20Bin is the compiled bytecode used for deploying new contracts.
const IERC20Bin = `0x`

// DeployIERC20 deploys a new Klaytn contract, binding an instance of IERC20 to it.
func DeployIERC20(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *IERC20, error) {
	parsed, err := abi.JSON(strings.NewReader(IERC20ABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(IERC20Bin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &IERC20{IERC20Caller: IERC20Caller{contract: contract}, IERC20Transactor: IERC20Transactor{contract: contract}, IERC20Filterer: IERC20Filterer{contract: contract}}, nil
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
// Solidity: function allowance(owner address, spender address) constant returns(uint256)
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
// Solidity: function allowance(owner address, spender address) constant returns(uint256)
func (_IERC20 *IERC20Session) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _IERC20.Contract.Allowance(&_IERC20.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(owner address, spender address) constant returns(uint256)
func (_IERC20 *IERC20CallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _IERC20.Contract.Allowance(&_IERC20.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(who address) constant returns(uint256)
func (_IERC20 *IERC20Caller) BalanceOf(opts *bind.CallOpts, who common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _IERC20.contract.Call(opts, out, "balanceOf", who)
	return *ret0, err
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(who address) constant returns(uint256)
func (_IERC20 *IERC20Session) BalanceOf(who common.Address) (*big.Int, error) {
	return _IERC20.Contract.BalanceOf(&_IERC20.CallOpts, who)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(who address) constant returns(uint256)
func (_IERC20 *IERC20CallerSession) BalanceOf(who common.Address) (*big.Int, error) {
	return _IERC20.Contract.BalanceOf(&_IERC20.CallOpts, who)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
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
// Solidity: function totalSupply() constant returns(uint256)
func (_IERC20 *IERC20Session) TotalSupply() (*big.Int, error) {
	return _IERC20.Contract.TotalSupply(&_IERC20.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_IERC20 *IERC20CallerSession) TotalSupply() (*big.Int, error) {
	return _IERC20.Contract.TotalSupply(&_IERC20.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(spender address, value uint256) returns(bool)
func (_IERC20 *IERC20Transactor) Approve(opts *bind.TransactOpts, spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _IERC20.contract.Transact(opts, "approve", spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(spender address, value uint256) returns(bool)
func (_IERC20 *IERC20Session) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _IERC20.Contract.Approve(&_IERC20.TransactOpts, spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(spender address, value uint256) returns(bool)
func (_IERC20 *IERC20TransactorSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _IERC20.Contract.Approve(&_IERC20.TransactOpts, spender, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(to address, value uint256) returns(bool)
func (_IERC20 *IERC20Transactor) Transfer(opts *bind.TransactOpts, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _IERC20.contract.Transact(opts, "transfer", to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(to address, value uint256) returns(bool)
func (_IERC20 *IERC20Session) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _IERC20.Contract.Transfer(&_IERC20.TransactOpts, to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(to address, value uint256) returns(bool)
func (_IERC20 *IERC20TransactorSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _IERC20.Contract.Transfer(&_IERC20.TransactOpts, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(from address, to address, value uint256) returns(bool)
func (_IERC20 *IERC20Transactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _IERC20.contract.Transact(opts, "transferFrom", from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(from address, to address, value uint256) returns(bool)
func (_IERC20 *IERC20Session) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _IERC20.Contract.TransferFrom(&_IERC20.TransactOpts, from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(from address, to address, value uint256) returns(bool)
func (_IERC20 *IERC20TransactorSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _IERC20.Contract.TransferFrom(&_IERC20.TransactOpts, from, to, value)
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
// Solidity: e Approval(owner indexed address, spender indexed address, value uint256)
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
// Solidity: e Approval(owner indexed address, spender indexed address, value uint256)
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
// Solidity: e Transfer(from indexed address, to indexed address, value uint256)
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
// Solidity: e Transfer(from indexed address, to indexed address, value uint256)
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

// IERC20BridgeReceiverABI is the input ABI used to generate the binding from.
const IERC20BridgeReceiverABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"_from\",\"type\":\"address\"},{\"name\":\"_amount\",\"type\":\"uint256\"},{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_feeLimit\",\"type\":\"uint256\"},{\"name\":\"_extraData\",\"type\":\"uint256[]\"}],\"name\":\"onERC20Received\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// IERC20BridgeReceiverBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const IERC20BridgeReceiverBinRuntime = `0x`

// IERC20BridgeReceiverBin is the compiled bytecode used for deploying new contracts.
const IERC20BridgeReceiverBin = `0x`

// DeployIERC20BridgeReceiver deploys a new Klaytn contract, binding an instance of IERC20BridgeReceiver to it.
func DeployIERC20BridgeReceiver(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *IERC20BridgeReceiver, error) {
	parsed, err := abi.JSON(strings.NewReader(IERC20BridgeReceiverABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(IERC20BridgeReceiverBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &IERC20BridgeReceiver{IERC20BridgeReceiverCaller: IERC20BridgeReceiverCaller{contract: contract}, IERC20BridgeReceiverTransactor: IERC20BridgeReceiverTransactor{contract: contract}, IERC20BridgeReceiverFilterer: IERC20BridgeReceiverFilterer{contract: contract}}, nil
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

// OnERC20Received is a paid mutator transaction binding the contract method 0x349449b2.
//
// Solidity: function onERC20Received(_from address, _amount uint256, _to address, _feeLimit uint256, _extraData uint256[]) returns()
func (_IERC20BridgeReceiver *IERC20BridgeReceiverTransactor) OnERC20Received(opts *bind.TransactOpts, _from common.Address, _amount *big.Int, _to common.Address, _feeLimit *big.Int, _extraData []*big.Int) (*types.Transaction, error) {
	return _IERC20BridgeReceiver.contract.Transact(opts, "onERC20Received", _from, _amount, _to, _feeLimit, _extraData)
}

// OnERC20Received is a paid mutator transaction binding the contract method 0x349449b2.
//
// Solidity: function onERC20Received(_from address, _amount uint256, _to address, _feeLimit uint256, _extraData uint256[]) returns()
func (_IERC20BridgeReceiver *IERC20BridgeReceiverSession) OnERC20Received(_from common.Address, _amount *big.Int, _to common.Address, _feeLimit *big.Int, _extraData []*big.Int) (*types.Transaction, error) {
	return _IERC20BridgeReceiver.Contract.OnERC20Received(&_IERC20BridgeReceiver.TransactOpts, _from, _amount, _to, _feeLimit, _extraData)
}

// OnERC20Received is a paid mutator transaction binding the contract method 0x349449b2.
//
// Solidity: function onERC20Received(_from address, _amount uint256, _to address, _feeLimit uint256, _extraData uint256[]) returns()
func (_IERC20BridgeReceiver *IERC20BridgeReceiverTransactorSession) OnERC20Received(_from common.Address, _amount *big.Int, _to common.Address, _feeLimit *big.Int, _extraData []*big.Int) (*types.Transaction, error) {
	return _IERC20BridgeReceiver.Contract.OnERC20Received(&_IERC20BridgeReceiver.TransactOpts, _from, _amount, _to, _feeLimit, _extraData)
}

// IERC721ABI is the input ABI used to generate the binding from.
const IERC721ABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"getApproved\",\"outputs\":[{\"name\":\"operator\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"from\",\"type\":\"address\"},{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"from\",\"type\":\"address\"},{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"safeTransferFrom\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"ownerOf\",\"outputs\":[{\"name\":\"owner\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"balance\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"operator\",\"type\":\"address\"},{\"name\":\"_approved\",\"type\":\"bool\"}],\"name\":\"setApprovalForAll\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"from\",\"type\":\"address\"},{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\"},{\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"safeTransferFrom\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"owner\",\"type\":\"address\"},{\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"isApprovedForAll\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"approved\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"ApprovalForAll\",\"type\":\"event\"}]"

// IERC721BinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const IERC721BinRuntime = `0x`

// IERC721Bin is the compiled bytecode used for deploying new contracts.
const IERC721Bin = `0x`

// DeployIERC721 deploys a new Klaytn contract, binding an instance of IERC721 to it.
func DeployIERC721(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *IERC721, error) {
	parsed, err := abi.JSON(strings.NewReader(IERC721ABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(IERC721Bin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &IERC721{IERC721Caller: IERC721Caller{contract: contract}, IERC721Transactor: IERC721Transactor{contract: contract}, IERC721Filterer: IERC721Filterer{contract: contract}}, nil
}

// IERC721 is an auto generated Go binding around a Klaytn contract.
type IERC721 struct {
	IERC721Caller     // Read-only binding to the contract
	IERC721Transactor // Write-only binding to the contract
	IERC721Filterer   // Log filterer for contract events
}

// IERC721Caller is an auto generated read-only Go binding around a Klaytn contract.
type IERC721Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IERC721Transactor is an auto generated write-only Go binding around a Klaytn contract.
type IERC721Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IERC721Filterer is an auto generated log filtering Go binding around a Klaytn contract events.
type IERC721Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IERC721Session is an auto generated Go binding around a Klaytn contract,
// with pre-set call and transact options.
type IERC721Session struct {
	Contract     *IERC721          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IERC721CallerSession is an auto generated read-only Go binding around a Klaytn contract,
// with pre-set call options.
type IERC721CallerSession struct {
	Contract *IERC721Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// IERC721TransactorSession is an auto generated write-only Go binding around a Klaytn contract,
// with pre-set transact options.
type IERC721TransactorSession struct {
	Contract     *IERC721Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// IERC721Raw is an auto generated low-level Go binding around a Klaytn contract.
type IERC721Raw struct {
	Contract *IERC721 // Generic contract binding to access the raw methods on
}

// IERC721CallerRaw is an auto generated low-level read-only Go binding around a Klaytn contract.
type IERC721CallerRaw struct {
	Contract *IERC721Caller // Generic read-only contract binding to access the raw methods on
}

// IERC721TransactorRaw is an auto generated low-level write-only Go binding around a Klaytn contract.
type IERC721TransactorRaw struct {
	Contract *IERC721Transactor // Generic write-only contract binding to access the raw methods on
}

// NewIERC721 creates a new instance of IERC721, bound to a specific deployed contract.
func NewIERC721(address common.Address, backend bind.ContractBackend) (*IERC721, error) {
	contract, err := bindIERC721(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IERC721{IERC721Caller: IERC721Caller{contract: contract}, IERC721Transactor: IERC721Transactor{contract: contract}, IERC721Filterer: IERC721Filterer{contract: contract}}, nil
}

// NewIERC721Caller creates a new read-only instance of IERC721, bound to a specific deployed contract.
func NewIERC721Caller(address common.Address, caller bind.ContractCaller) (*IERC721Caller, error) {
	contract, err := bindIERC721(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IERC721Caller{contract: contract}, nil
}

// NewIERC721Transactor creates a new write-only instance of IERC721, bound to a specific deployed contract.
func NewIERC721Transactor(address common.Address, transactor bind.ContractTransactor) (*IERC721Transactor, error) {
	contract, err := bindIERC721(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IERC721Transactor{contract: contract}, nil
}

// NewIERC721Filterer creates a new log filterer instance of IERC721, bound to a specific deployed contract.
func NewIERC721Filterer(address common.Address, filterer bind.ContractFilterer) (*IERC721Filterer, error) {
	contract, err := bindIERC721(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IERC721Filterer{contract: contract}, nil
}

// bindIERC721 binds a generic wrapper to an already deployed contract.
func bindIERC721(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(IERC721ABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IERC721 *IERC721Raw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _IERC721.Contract.IERC721Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IERC721 *IERC721Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IERC721.Contract.IERC721Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IERC721 *IERC721Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IERC721.Contract.IERC721Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IERC721 *IERC721CallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _IERC721.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IERC721 *IERC721TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IERC721.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IERC721 *IERC721TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IERC721.Contract.contract.Transact(opts, method, params...)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(owner address) constant returns(balance uint256)
func (_IERC721 *IERC721Caller) BalanceOf(opts *bind.CallOpts, owner common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _IERC721.contract.Call(opts, out, "balanceOf", owner)
	return *ret0, err
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(owner address) constant returns(balance uint256)
func (_IERC721 *IERC721Session) BalanceOf(owner common.Address) (*big.Int, error) {
	return _IERC721.Contract.BalanceOf(&_IERC721.CallOpts, owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(owner address) constant returns(balance uint256)
func (_IERC721 *IERC721CallerSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _IERC721.Contract.BalanceOf(&_IERC721.CallOpts, owner)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(tokenId uint256) constant returns(operator address)
func (_IERC721 *IERC721Caller) GetApproved(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _IERC721.contract.Call(opts, out, "getApproved", tokenId)
	return *ret0, err
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(tokenId uint256) constant returns(operator address)
func (_IERC721 *IERC721Session) GetApproved(tokenId *big.Int) (common.Address, error) {
	return _IERC721.Contract.GetApproved(&_IERC721.CallOpts, tokenId)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(tokenId uint256) constant returns(operator address)
func (_IERC721 *IERC721CallerSession) GetApproved(tokenId *big.Int) (common.Address, error) {
	return _IERC721.Contract.GetApproved(&_IERC721.CallOpts, tokenId)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(owner address, operator address) constant returns(bool)
func (_IERC721 *IERC721Caller) IsApprovedForAll(opts *bind.CallOpts, owner common.Address, operator common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _IERC721.contract.Call(opts, out, "isApprovedForAll", owner, operator)
	return *ret0, err
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(owner address, operator address) constant returns(bool)
func (_IERC721 *IERC721Session) IsApprovedForAll(owner common.Address, operator common.Address) (bool, error) {
	return _IERC721.Contract.IsApprovedForAll(&_IERC721.CallOpts, owner, operator)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(owner address, operator address) constant returns(bool)
func (_IERC721 *IERC721CallerSession) IsApprovedForAll(owner common.Address, operator common.Address) (bool, error) {
	return _IERC721.Contract.IsApprovedForAll(&_IERC721.CallOpts, owner, operator)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(tokenId uint256) constant returns(owner address)
func (_IERC721 *IERC721Caller) OwnerOf(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _IERC721.contract.Call(opts, out, "ownerOf", tokenId)
	return *ret0, err
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(tokenId uint256) constant returns(owner address)
func (_IERC721 *IERC721Session) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _IERC721.Contract.OwnerOf(&_IERC721.CallOpts, tokenId)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(tokenId uint256) constant returns(owner address)
func (_IERC721 *IERC721CallerSession) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _IERC721.Contract.OwnerOf(&_IERC721.CallOpts, tokenId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(interfaceId bytes4) constant returns(bool)
func (_IERC721 *IERC721Caller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _IERC721.contract.Call(opts, out, "supportsInterface", interfaceId)
	return *ret0, err
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(interfaceId bytes4) constant returns(bool)
func (_IERC721 *IERC721Session) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _IERC721.Contract.SupportsInterface(&_IERC721.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(interfaceId bytes4) constant returns(bool)
func (_IERC721 *IERC721CallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _IERC721.Contract.SupportsInterface(&_IERC721.CallOpts, interfaceId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(to address, tokenId uint256) returns()
func (_IERC721 *IERC721Transactor) Approve(opts *bind.TransactOpts, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _IERC721.contract.Transact(opts, "approve", to, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(to address, tokenId uint256) returns()
func (_IERC721 *IERC721Session) Approve(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _IERC721.Contract.Approve(&_IERC721.TransactOpts, to, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(to address, tokenId uint256) returns()
func (_IERC721 *IERC721TransactorSession) Approve(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _IERC721.Contract.Approve(&_IERC721.TransactOpts, to, tokenId)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(from address, to address, tokenId uint256, data bytes) returns()
func (_IERC721 *IERC721Transactor) SafeTransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _IERC721.contract.Transact(opts, "safeTransferFrom", from, to, tokenId, data)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(from address, to address, tokenId uint256, data bytes) returns()
func (_IERC721 *IERC721Session) SafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _IERC721.Contract.SafeTransferFrom(&_IERC721.TransactOpts, from, to, tokenId, data)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(from address, to address, tokenId uint256, data bytes) returns()
func (_IERC721 *IERC721TransactorSession) SafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _IERC721.Contract.SafeTransferFrom(&_IERC721.TransactOpts, from, to, tokenId, data)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(operator address, _approved bool) returns()
func (_IERC721 *IERC721Transactor) SetApprovalForAll(opts *bind.TransactOpts, operator common.Address, _approved bool) (*types.Transaction, error) {
	return _IERC721.contract.Transact(opts, "setApprovalForAll", operator, _approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(operator address, _approved bool) returns()
func (_IERC721 *IERC721Session) SetApprovalForAll(operator common.Address, _approved bool) (*types.Transaction, error) {
	return _IERC721.Contract.SetApprovalForAll(&_IERC721.TransactOpts, operator, _approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(operator address, _approved bool) returns()
func (_IERC721 *IERC721TransactorSession) SetApprovalForAll(operator common.Address, _approved bool) (*types.Transaction, error) {
	return _IERC721.Contract.SetApprovalForAll(&_IERC721.TransactOpts, operator, _approved)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(from address, to address, tokenId uint256) returns()
func (_IERC721 *IERC721Transactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _IERC721.contract.Transact(opts, "transferFrom", from, to, tokenId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(from address, to address, tokenId uint256) returns()
func (_IERC721 *IERC721Session) TransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _IERC721.Contract.TransferFrom(&_IERC721.TransactOpts, from, to, tokenId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(from address, to address, tokenId uint256) returns()
func (_IERC721 *IERC721TransactorSession) TransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _IERC721.Contract.TransferFrom(&_IERC721.TransactOpts, from, to, tokenId)
}

// IERC721ApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the IERC721 contract.
type IERC721ApprovalIterator struct {
	Event *IERC721Approval // Event containing the contract specifics and raw log

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
func (it *IERC721ApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IERC721Approval)
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
		it.Event = new(IERC721Approval)
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
func (it *IERC721ApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IERC721ApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IERC721Approval represents a Approval event raised by the IERC721 contract.
type IERC721Approval struct {
	Owner    common.Address
	Approved common.Address
	TokenId  *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: e Approval(owner indexed address, approved indexed address, tokenId indexed uint256)
func (_IERC721 *IERC721Filterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, approved []common.Address, tokenId []*big.Int) (*IERC721ApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var approvedRule []interface{}
	for _, approvedItem := range approved {
		approvedRule = append(approvedRule, approvedItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _IERC721.contract.FilterLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &IERC721ApprovalIterator{contract: _IERC721.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: e Approval(owner indexed address, approved indexed address, tokenId indexed uint256)
func (_IERC721 *IERC721Filterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *IERC721Approval, owner []common.Address, approved []common.Address, tokenId []*big.Int) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var approvedRule []interface{}
	for _, approvedItem := range approved {
		approvedRule = append(approvedRule, approvedItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _IERC721.contract.WatchLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IERC721Approval)
				if err := _IERC721.contract.UnpackLog(event, "Approval", log); err != nil {
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

// IERC721ApprovalForAllIterator is returned from FilterApprovalForAll and is used to iterate over the raw logs and unpacked data for ApprovalForAll events raised by the IERC721 contract.
type IERC721ApprovalForAllIterator struct {
	Event *IERC721ApprovalForAll // Event containing the contract specifics and raw log

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
func (it *IERC721ApprovalForAllIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IERC721ApprovalForAll)
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
		it.Event = new(IERC721ApprovalForAll)
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
func (it *IERC721ApprovalForAllIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IERC721ApprovalForAllIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IERC721ApprovalForAll represents a ApprovalForAll event raised by the IERC721 contract.
type IERC721ApprovalForAll struct {
	Owner    common.Address
	Operator common.Address
	Approved bool
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApprovalForAll is a free log retrieval operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: e ApprovalForAll(owner indexed address, operator indexed address, approved bool)
func (_IERC721 *IERC721Filterer) FilterApprovalForAll(opts *bind.FilterOpts, owner []common.Address, operator []common.Address) (*IERC721ApprovalForAllIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _IERC721.contract.FilterLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return &IERC721ApprovalForAllIterator{contract: _IERC721.contract, event: "ApprovalForAll", logs: logs, sub: sub}, nil
}

// WatchApprovalForAll is a free log subscription operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: e ApprovalForAll(owner indexed address, operator indexed address, approved bool)
func (_IERC721 *IERC721Filterer) WatchApprovalForAll(opts *bind.WatchOpts, sink chan<- *IERC721ApprovalForAll, owner []common.Address, operator []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _IERC721.contract.WatchLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IERC721ApprovalForAll)
				if err := _IERC721.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
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

// IERC721TransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the IERC721 contract.
type IERC721TransferIterator struct {
	Event *IERC721Transfer // Event containing the contract specifics and raw log

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
func (it *IERC721TransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IERC721Transfer)
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
		it.Event = new(IERC721Transfer)
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
func (it *IERC721TransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IERC721TransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IERC721Transfer represents a Transfer event raised by the IERC721 contract.
type IERC721Transfer struct {
	From    common.Address
	To      common.Address
	TokenId *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: e Transfer(from indexed address, to indexed address, tokenId indexed uint256)
func (_IERC721 *IERC721Filterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address, tokenId []*big.Int) (*IERC721TransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _IERC721.contract.FilterLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &IERC721TransferIterator{contract: _IERC721.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: e Transfer(from indexed address, to indexed address, tokenId indexed uint256)
func (_IERC721 *IERC721Filterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *IERC721Transfer, from []common.Address, to []common.Address, tokenId []*big.Int) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _IERC721.contract.WatchLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IERC721Transfer)
				if err := _IERC721.contract.UnpackLog(event, "Transfer", log); err != nil {
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

// IERC721BridgeReceiverABI is the input ABI used to generate the binding from.
const IERC721BridgeReceiverABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"_from\",\"type\":\"address\"},{\"name\":\"_tokenId\",\"type\":\"uint256\"},{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_extraData\",\"type\":\"uint256[]\"}],\"name\":\"onERC721Received\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// IERC721BridgeReceiverBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const IERC721BridgeReceiverBinRuntime = `0x`

// IERC721BridgeReceiverBin is the compiled bytecode used for deploying new contracts.
const IERC721BridgeReceiverBin = `0x`

// DeployIERC721BridgeReceiver deploys a new Klaytn contract, binding an instance of IERC721BridgeReceiver to it.
func DeployIERC721BridgeReceiver(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *IERC721BridgeReceiver, error) {
	parsed, err := abi.JSON(strings.NewReader(IERC721BridgeReceiverABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(IERC721BridgeReceiverBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &IERC721BridgeReceiver{IERC721BridgeReceiverCaller: IERC721BridgeReceiverCaller{contract: contract}, IERC721BridgeReceiverTransactor: IERC721BridgeReceiverTransactor{contract: contract}, IERC721BridgeReceiverFilterer: IERC721BridgeReceiverFilterer{contract: contract}}, nil
}

// IERC721BridgeReceiver is an auto generated Go binding around a Klaytn contract.
type IERC721BridgeReceiver struct {
	IERC721BridgeReceiverCaller     // Read-only binding to the contract
	IERC721BridgeReceiverTransactor // Write-only binding to the contract
	IERC721BridgeReceiverFilterer   // Log filterer for contract events
}

// IERC721BridgeReceiverCaller is an auto generated read-only Go binding around a Klaytn contract.
type IERC721BridgeReceiverCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IERC721BridgeReceiverTransactor is an auto generated write-only Go binding around a Klaytn contract.
type IERC721BridgeReceiverTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IERC721BridgeReceiverFilterer is an auto generated log filtering Go binding around a Klaytn contract events.
type IERC721BridgeReceiverFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IERC721BridgeReceiverSession is an auto generated Go binding around a Klaytn contract,
// with pre-set call and transact options.
type IERC721BridgeReceiverSession struct {
	Contract     *IERC721BridgeReceiver // Generic contract binding to set the session for
	CallOpts     bind.CallOpts          // Call options to use throughout this session
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// IERC721BridgeReceiverCallerSession is an auto generated read-only Go binding around a Klaytn contract,
// with pre-set call options.
type IERC721BridgeReceiverCallerSession struct {
	Contract *IERC721BridgeReceiverCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                // Call options to use throughout this session
}

// IERC721BridgeReceiverTransactorSession is an auto generated write-only Go binding around a Klaytn contract,
// with pre-set transact options.
type IERC721BridgeReceiverTransactorSession struct {
	Contract     *IERC721BridgeReceiverTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                // Transaction auth options to use throughout this session
}

// IERC721BridgeReceiverRaw is an auto generated low-level Go binding around a Klaytn contract.
type IERC721BridgeReceiverRaw struct {
	Contract *IERC721BridgeReceiver // Generic contract binding to access the raw methods on
}

// IERC721BridgeReceiverCallerRaw is an auto generated low-level read-only Go binding around a Klaytn contract.
type IERC721BridgeReceiverCallerRaw struct {
	Contract *IERC721BridgeReceiverCaller // Generic read-only contract binding to access the raw methods on
}

// IERC721BridgeReceiverTransactorRaw is an auto generated low-level write-only Go binding around a Klaytn contract.
type IERC721BridgeReceiverTransactorRaw struct {
	Contract *IERC721BridgeReceiverTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIERC721BridgeReceiver creates a new instance of IERC721BridgeReceiver, bound to a specific deployed contract.
func NewIERC721BridgeReceiver(address common.Address, backend bind.ContractBackend) (*IERC721BridgeReceiver, error) {
	contract, err := bindIERC721BridgeReceiver(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IERC721BridgeReceiver{IERC721BridgeReceiverCaller: IERC721BridgeReceiverCaller{contract: contract}, IERC721BridgeReceiverTransactor: IERC721BridgeReceiverTransactor{contract: contract}, IERC721BridgeReceiverFilterer: IERC721BridgeReceiverFilterer{contract: contract}}, nil
}

// NewIERC721BridgeReceiverCaller creates a new read-only instance of IERC721BridgeReceiver, bound to a specific deployed contract.
func NewIERC721BridgeReceiverCaller(address common.Address, caller bind.ContractCaller) (*IERC721BridgeReceiverCaller, error) {
	contract, err := bindIERC721BridgeReceiver(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IERC721BridgeReceiverCaller{contract: contract}, nil
}

// NewIERC721BridgeReceiverTransactor creates a new write-only instance of IERC721BridgeReceiver, bound to a specific deployed contract.
func NewIERC721BridgeReceiverTransactor(address common.Address, transactor bind.ContractTransactor) (*IERC721BridgeReceiverTransactor, error) {
	contract, err := bindIERC721BridgeReceiver(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IERC721BridgeReceiverTransactor{contract: contract}, nil
}

// NewIERC721BridgeReceiverFilterer creates a new log filterer instance of IERC721BridgeReceiver, bound to a specific deployed contract.
func NewIERC721BridgeReceiverFilterer(address common.Address, filterer bind.ContractFilterer) (*IERC721BridgeReceiverFilterer, error) {
	contract, err := bindIERC721BridgeReceiver(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IERC721BridgeReceiverFilterer{contract: contract}, nil
}

// bindIERC721BridgeReceiver binds a generic wrapper to an already deployed contract.
func bindIERC721BridgeReceiver(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(IERC721BridgeReceiverABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IERC721BridgeReceiver *IERC721BridgeReceiverRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _IERC721BridgeReceiver.Contract.IERC721BridgeReceiverCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IERC721BridgeReceiver *IERC721BridgeReceiverRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IERC721BridgeReceiver.Contract.IERC721BridgeReceiverTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IERC721BridgeReceiver *IERC721BridgeReceiverRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IERC721BridgeReceiver.Contract.IERC721BridgeReceiverTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IERC721BridgeReceiver *IERC721BridgeReceiverCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _IERC721BridgeReceiver.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IERC721BridgeReceiver *IERC721BridgeReceiverTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IERC721BridgeReceiver.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IERC721BridgeReceiver *IERC721BridgeReceiverTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IERC721BridgeReceiver.Contract.contract.Transact(opts, method, params...)
}

// OnERC721Received is a paid mutator transaction binding the contract method 0xa319de9e.
//
// Solidity: function onERC721Received(_from address, _tokenId uint256, _to address, _extraData uint256[]) returns()
func (_IERC721BridgeReceiver *IERC721BridgeReceiverTransactor) OnERC721Received(opts *bind.TransactOpts, _from common.Address, _tokenId *big.Int, _to common.Address, _extraData []*big.Int) (*types.Transaction, error) {
	return _IERC721BridgeReceiver.contract.Transact(opts, "onERC721Received", _from, _tokenId, _to, _extraData)
}

// OnERC721Received is a paid mutator transaction binding the contract method 0xa319de9e.
//
// Solidity: function onERC721Received(_from address, _tokenId uint256, _to address, _extraData uint256[]) returns()
func (_IERC721BridgeReceiver *IERC721BridgeReceiverSession) OnERC721Received(_from common.Address, _tokenId *big.Int, _to common.Address, _extraData []*big.Int) (*types.Transaction, error) {
	return _IERC721BridgeReceiver.Contract.OnERC721Received(&_IERC721BridgeReceiver.TransactOpts, _from, _tokenId, _to, _extraData)
}

// OnERC721Received is a paid mutator transaction binding the contract method 0xa319de9e.
//
// Solidity: function onERC721Received(_from address, _tokenId uint256, _to address, _extraData uint256[]) returns()
func (_IERC721BridgeReceiver *IERC721BridgeReceiverTransactorSession) OnERC721Received(_from common.Address, _tokenId *big.Int, _to common.Address, _extraData []*big.Int) (*types.Transaction, error) {
	return _IERC721BridgeReceiver.Contract.OnERC721Received(&_IERC721BridgeReceiver.TransactOpts, _from, _tokenId, _to, _extraData)
}

// IERC721MetadataABI is the input ABI used to generate the binding from.
const IERC721MetadataABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"getApproved\",\"outputs\":[{\"name\":\"operator\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"from\",\"type\":\"address\"},{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"from\",\"type\":\"address\"},{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"safeTransferFrom\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"ownerOf\",\"outputs\":[{\"name\":\"owner\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"balance\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"operator\",\"type\":\"address\"},{\"name\":\"_approved\",\"type\":\"bool\"}],\"name\":\"setApprovalForAll\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"from\",\"type\":\"address\"},{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\"},{\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"safeTransferFrom\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"tokenURI\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"owner\",\"type\":\"address\"},{\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"isApprovedForAll\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"approved\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"ApprovalForAll\",\"type\":\"event\"}]"

// IERC721MetadataBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const IERC721MetadataBinRuntime = `0x`

// IERC721MetadataBin is the compiled bytecode used for deploying new contracts.
const IERC721MetadataBin = `0x`

// DeployIERC721Metadata deploys a new Klaytn contract, binding an instance of IERC721Metadata to it.
func DeployIERC721Metadata(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *IERC721Metadata, error) {
	parsed, err := abi.JSON(strings.NewReader(IERC721MetadataABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(IERC721MetadataBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &IERC721Metadata{IERC721MetadataCaller: IERC721MetadataCaller{contract: contract}, IERC721MetadataTransactor: IERC721MetadataTransactor{contract: contract}, IERC721MetadataFilterer: IERC721MetadataFilterer{contract: contract}}, nil
}

// IERC721Metadata is an auto generated Go binding around a Klaytn contract.
type IERC721Metadata struct {
	IERC721MetadataCaller     // Read-only binding to the contract
	IERC721MetadataTransactor // Write-only binding to the contract
	IERC721MetadataFilterer   // Log filterer for contract events
}

// IERC721MetadataCaller is an auto generated read-only Go binding around a Klaytn contract.
type IERC721MetadataCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IERC721MetadataTransactor is an auto generated write-only Go binding around a Klaytn contract.
type IERC721MetadataTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IERC721MetadataFilterer is an auto generated log filtering Go binding around a Klaytn contract events.
type IERC721MetadataFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IERC721MetadataSession is an auto generated Go binding around a Klaytn contract,
// with pre-set call and transact options.
type IERC721MetadataSession struct {
	Contract     *IERC721Metadata  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IERC721MetadataCallerSession is an auto generated read-only Go binding around a Klaytn contract,
// with pre-set call options.
type IERC721MetadataCallerSession struct {
	Contract *IERC721MetadataCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// IERC721MetadataTransactorSession is an auto generated write-only Go binding around a Klaytn contract,
// with pre-set transact options.
type IERC721MetadataTransactorSession struct {
	Contract     *IERC721MetadataTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// IERC721MetadataRaw is an auto generated low-level Go binding around a Klaytn contract.
type IERC721MetadataRaw struct {
	Contract *IERC721Metadata // Generic contract binding to access the raw methods on
}

// IERC721MetadataCallerRaw is an auto generated low-level read-only Go binding around a Klaytn contract.
type IERC721MetadataCallerRaw struct {
	Contract *IERC721MetadataCaller // Generic read-only contract binding to access the raw methods on
}

// IERC721MetadataTransactorRaw is an auto generated low-level write-only Go binding around a Klaytn contract.
type IERC721MetadataTransactorRaw struct {
	Contract *IERC721MetadataTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIERC721Metadata creates a new instance of IERC721Metadata, bound to a specific deployed contract.
func NewIERC721Metadata(address common.Address, backend bind.ContractBackend) (*IERC721Metadata, error) {
	contract, err := bindIERC721Metadata(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IERC721Metadata{IERC721MetadataCaller: IERC721MetadataCaller{contract: contract}, IERC721MetadataTransactor: IERC721MetadataTransactor{contract: contract}, IERC721MetadataFilterer: IERC721MetadataFilterer{contract: contract}}, nil
}

// NewIERC721MetadataCaller creates a new read-only instance of IERC721Metadata, bound to a specific deployed contract.
func NewIERC721MetadataCaller(address common.Address, caller bind.ContractCaller) (*IERC721MetadataCaller, error) {
	contract, err := bindIERC721Metadata(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IERC721MetadataCaller{contract: contract}, nil
}

// NewIERC721MetadataTransactor creates a new write-only instance of IERC721Metadata, bound to a specific deployed contract.
func NewIERC721MetadataTransactor(address common.Address, transactor bind.ContractTransactor) (*IERC721MetadataTransactor, error) {
	contract, err := bindIERC721Metadata(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IERC721MetadataTransactor{contract: contract}, nil
}

// NewIERC721MetadataFilterer creates a new log filterer instance of IERC721Metadata, bound to a specific deployed contract.
func NewIERC721MetadataFilterer(address common.Address, filterer bind.ContractFilterer) (*IERC721MetadataFilterer, error) {
	contract, err := bindIERC721Metadata(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IERC721MetadataFilterer{contract: contract}, nil
}

// bindIERC721Metadata binds a generic wrapper to an already deployed contract.
func bindIERC721Metadata(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(IERC721MetadataABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IERC721Metadata *IERC721MetadataRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _IERC721Metadata.Contract.IERC721MetadataCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IERC721Metadata *IERC721MetadataRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IERC721Metadata.Contract.IERC721MetadataTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IERC721Metadata *IERC721MetadataRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IERC721Metadata.Contract.IERC721MetadataTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IERC721Metadata *IERC721MetadataCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _IERC721Metadata.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IERC721Metadata *IERC721MetadataTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IERC721Metadata.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IERC721Metadata *IERC721MetadataTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IERC721Metadata.Contract.contract.Transact(opts, method, params...)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(owner address) constant returns(balance uint256)
func (_IERC721Metadata *IERC721MetadataCaller) BalanceOf(opts *bind.CallOpts, owner common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _IERC721Metadata.contract.Call(opts, out, "balanceOf", owner)
	return *ret0, err
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(owner address) constant returns(balance uint256)
func (_IERC721Metadata *IERC721MetadataSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _IERC721Metadata.Contract.BalanceOf(&_IERC721Metadata.CallOpts, owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(owner address) constant returns(balance uint256)
func (_IERC721Metadata *IERC721MetadataCallerSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _IERC721Metadata.Contract.BalanceOf(&_IERC721Metadata.CallOpts, owner)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(tokenId uint256) constant returns(operator address)
func (_IERC721Metadata *IERC721MetadataCaller) GetApproved(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _IERC721Metadata.contract.Call(opts, out, "getApproved", tokenId)
	return *ret0, err
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(tokenId uint256) constant returns(operator address)
func (_IERC721Metadata *IERC721MetadataSession) GetApproved(tokenId *big.Int) (common.Address, error) {
	return _IERC721Metadata.Contract.GetApproved(&_IERC721Metadata.CallOpts, tokenId)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(tokenId uint256) constant returns(operator address)
func (_IERC721Metadata *IERC721MetadataCallerSession) GetApproved(tokenId *big.Int) (common.Address, error) {
	return _IERC721Metadata.Contract.GetApproved(&_IERC721Metadata.CallOpts, tokenId)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(owner address, operator address) constant returns(bool)
func (_IERC721Metadata *IERC721MetadataCaller) IsApprovedForAll(opts *bind.CallOpts, owner common.Address, operator common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _IERC721Metadata.contract.Call(opts, out, "isApprovedForAll", owner, operator)
	return *ret0, err
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(owner address, operator address) constant returns(bool)
func (_IERC721Metadata *IERC721MetadataSession) IsApprovedForAll(owner common.Address, operator common.Address) (bool, error) {
	return _IERC721Metadata.Contract.IsApprovedForAll(&_IERC721Metadata.CallOpts, owner, operator)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(owner address, operator address) constant returns(bool)
func (_IERC721Metadata *IERC721MetadataCallerSession) IsApprovedForAll(owner common.Address, operator common.Address) (bool, error) {
	return _IERC721Metadata.Contract.IsApprovedForAll(&_IERC721Metadata.CallOpts, owner, operator)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() constant returns(string)
func (_IERC721Metadata *IERC721MetadataCaller) Name(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _IERC721Metadata.contract.Call(opts, out, "name")
	return *ret0, err
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() constant returns(string)
func (_IERC721Metadata *IERC721MetadataSession) Name() (string, error) {
	return _IERC721Metadata.Contract.Name(&_IERC721Metadata.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() constant returns(string)
func (_IERC721Metadata *IERC721MetadataCallerSession) Name() (string, error) {
	return _IERC721Metadata.Contract.Name(&_IERC721Metadata.CallOpts)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(tokenId uint256) constant returns(owner address)
func (_IERC721Metadata *IERC721MetadataCaller) OwnerOf(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _IERC721Metadata.contract.Call(opts, out, "ownerOf", tokenId)
	return *ret0, err
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(tokenId uint256) constant returns(owner address)
func (_IERC721Metadata *IERC721MetadataSession) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _IERC721Metadata.Contract.OwnerOf(&_IERC721Metadata.CallOpts, tokenId)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(tokenId uint256) constant returns(owner address)
func (_IERC721Metadata *IERC721MetadataCallerSession) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _IERC721Metadata.Contract.OwnerOf(&_IERC721Metadata.CallOpts, tokenId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(interfaceId bytes4) constant returns(bool)
func (_IERC721Metadata *IERC721MetadataCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _IERC721Metadata.contract.Call(opts, out, "supportsInterface", interfaceId)
	return *ret0, err
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(interfaceId bytes4) constant returns(bool)
func (_IERC721Metadata *IERC721MetadataSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _IERC721Metadata.Contract.SupportsInterface(&_IERC721Metadata.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(interfaceId bytes4) constant returns(bool)
func (_IERC721Metadata *IERC721MetadataCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _IERC721Metadata.Contract.SupportsInterface(&_IERC721Metadata.CallOpts, interfaceId)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() constant returns(string)
func (_IERC721Metadata *IERC721MetadataCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _IERC721Metadata.contract.Call(opts, out, "symbol")
	return *ret0, err
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() constant returns(string)
func (_IERC721Metadata *IERC721MetadataSession) Symbol() (string, error) {
	return _IERC721Metadata.Contract.Symbol(&_IERC721Metadata.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() constant returns(string)
func (_IERC721Metadata *IERC721MetadataCallerSession) Symbol() (string, error) {
	return _IERC721Metadata.Contract.Symbol(&_IERC721Metadata.CallOpts)
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(tokenId uint256) constant returns(string)
func (_IERC721Metadata *IERC721MetadataCaller) TokenURI(opts *bind.CallOpts, tokenId *big.Int) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _IERC721Metadata.contract.Call(opts, out, "tokenURI", tokenId)
	return *ret0, err
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(tokenId uint256) constant returns(string)
func (_IERC721Metadata *IERC721MetadataSession) TokenURI(tokenId *big.Int) (string, error) {
	return _IERC721Metadata.Contract.TokenURI(&_IERC721Metadata.CallOpts, tokenId)
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(tokenId uint256) constant returns(string)
func (_IERC721Metadata *IERC721MetadataCallerSession) TokenURI(tokenId *big.Int) (string, error) {
	return _IERC721Metadata.Contract.TokenURI(&_IERC721Metadata.CallOpts, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(to address, tokenId uint256) returns()
func (_IERC721Metadata *IERC721MetadataTransactor) Approve(opts *bind.TransactOpts, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _IERC721Metadata.contract.Transact(opts, "approve", to, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(to address, tokenId uint256) returns()
func (_IERC721Metadata *IERC721MetadataSession) Approve(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _IERC721Metadata.Contract.Approve(&_IERC721Metadata.TransactOpts, to, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(to address, tokenId uint256) returns()
func (_IERC721Metadata *IERC721MetadataTransactorSession) Approve(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _IERC721Metadata.Contract.Approve(&_IERC721Metadata.TransactOpts, to, tokenId)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(from address, to address, tokenId uint256, data bytes) returns()
func (_IERC721Metadata *IERC721MetadataTransactor) SafeTransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _IERC721Metadata.contract.Transact(opts, "safeTransferFrom", from, to, tokenId, data)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(from address, to address, tokenId uint256, data bytes) returns()
func (_IERC721Metadata *IERC721MetadataSession) SafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _IERC721Metadata.Contract.SafeTransferFrom(&_IERC721Metadata.TransactOpts, from, to, tokenId, data)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(from address, to address, tokenId uint256, data bytes) returns()
func (_IERC721Metadata *IERC721MetadataTransactorSession) SafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _IERC721Metadata.Contract.SafeTransferFrom(&_IERC721Metadata.TransactOpts, from, to, tokenId, data)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(operator address, _approved bool) returns()
func (_IERC721Metadata *IERC721MetadataTransactor) SetApprovalForAll(opts *bind.TransactOpts, operator common.Address, _approved bool) (*types.Transaction, error) {
	return _IERC721Metadata.contract.Transact(opts, "setApprovalForAll", operator, _approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(operator address, _approved bool) returns()
func (_IERC721Metadata *IERC721MetadataSession) SetApprovalForAll(operator common.Address, _approved bool) (*types.Transaction, error) {
	return _IERC721Metadata.Contract.SetApprovalForAll(&_IERC721Metadata.TransactOpts, operator, _approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(operator address, _approved bool) returns()
func (_IERC721Metadata *IERC721MetadataTransactorSession) SetApprovalForAll(operator common.Address, _approved bool) (*types.Transaction, error) {
	return _IERC721Metadata.Contract.SetApprovalForAll(&_IERC721Metadata.TransactOpts, operator, _approved)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(from address, to address, tokenId uint256) returns()
func (_IERC721Metadata *IERC721MetadataTransactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _IERC721Metadata.contract.Transact(opts, "transferFrom", from, to, tokenId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(from address, to address, tokenId uint256) returns()
func (_IERC721Metadata *IERC721MetadataSession) TransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _IERC721Metadata.Contract.TransferFrom(&_IERC721Metadata.TransactOpts, from, to, tokenId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(from address, to address, tokenId uint256) returns()
func (_IERC721Metadata *IERC721MetadataTransactorSession) TransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _IERC721Metadata.Contract.TransferFrom(&_IERC721Metadata.TransactOpts, from, to, tokenId)
}

// IERC721MetadataApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the IERC721Metadata contract.
type IERC721MetadataApprovalIterator struct {
	Event *IERC721MetadataApproval // Event containing the contract specifics and raw log

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
func (it *IERC721MetadataApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IERC721MetadataApproval)
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
		it.Event = new(IERC721MetadataApproval)
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
func (it *IERC721MetadataApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IERC721MetadataApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IERC721MetadataApproval represents a Approval event raised by the IERC721Metadata contract.
type IERC721MetadataApproval struct {
	Owner    common.Address
	Approved common.Address
	TokenId  *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: e Approval(owner indexed address, approved indexed address, tokenId indexed uint256)
func (_IERC721Metadata *IERC721MetadataFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, approved []common.Address, tokenId []*big.Int) (*IERC721MetadataApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var approvedRule []interface{}
	for _, approvedItem := range approved {
		approvedRule = append(approvedRule, approvedItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _IERC721Metadata.contract.FilterLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &IERC721MetadataApprovalIterator{contract: _IERC721Metadata.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: e Approval(owner indexed address, approved indexed address, tokenId indexed uint256)
func (_IERC721Metadata *IERC721MetadataFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *IERC721MetadataApproval, owner []common.Address, approved []common.Address, tokenId []*big.Int) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var approvedRule []interface{}
	for _, approvedItem := range approved {
		approvedRule = append(approvedRule, approvedItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _IERC721Metadata.contract.WatchLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IERC721MetadataApproval)
				if err := _IERC721Metadata.contract.UnpackLog(event, "Approval", log); err != nil {
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

// IERC721MetadataApprovalForAllIterator is returned from FilterApprovalForAll and is used to iterate over the raw logs and unpacked data for ApprovalForAll events raised by the IERC721Metadata contract.
type IERC721MetadataApprovalForAllIterator struct {
	Event *IERC721MetadataApprovalForAll // Event containing the contract specifics and raw log

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
func (it *IERC721MetadataApprovalForAllIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IERC721MetadataApprovalForAll)
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
		it.Event = new(IERC721MetadataApprovalForAll)
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
func (it *IERC721MetadataApprovalForAllIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IERC721MetadataApprovalForAllIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IERC721MetadataApprovalForAll represents a ApprovalForAll event raised by the IERC721Metadata contract.
type IERC721MetadataApprovalForAll struct {
	Owner    common.Address
	Operator common.Address
	Approved bool
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApprovalForAll is a free log retrieval operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: e ApprovalForAll(owner indexed address, operator indexed address, approved bool)
func (_IERC721Metadata *IERC721MetadataFilterer) FilterApprovalForAll(opts *bind.FilterOpts, owner []common.Address, operator []common.Address) (*IERC721MetadataApprovalForAllIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _IERC721Metadata.contract.FilterLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return &IERC721MetadataApprovalForAllIterator{contract: _IERC721Metadata.contract, event: "ApprovalForAll", logs: logs, sub: sub}, nil
}

// WatchApprovalForAll is a free log subscription operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: e ApprovalForAll(owner indexed address, operator indexed address, approved bool)
func (_IERC721Metadata *IERC721MetadataFilterer) WatchApprovalForAll(opts *bind.WatchOpts, sink chan<- *IERC721MetadataApprovalForAll, owner []common.Address, operator []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _IERC721Metadata.contract.WatchLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IERC721MetadataApprovalForAll)
				if err := _IERC721Metadata.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
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

// IERC721MetadataTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the IERC721Metadata contract.
type IERC721MetadataTransferIterator struct {
	Event *IERC721MetadataTransfer // Event containing the contract specifics and raw log

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
func (it *IERC721MetadataTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IERC721MetadataTransfer)
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
		it.Event = new(IERC721MetadataTransfer)
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
func (it *IERC721MetadataTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IERC721MetadataTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IERC721MetadataTransfer represents a Transfer event raised by the IERC721Metadata contract.
type IERC721MetadataTransfer struct {
	From    common.Address
	To      common.Address
	TokenId *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: e Transfer(from indexed address, to indexed address, tokenId indexed uint256)
func (_IERC721Metadata *IERC721MetadataFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address, tokenId []*big.Int) (*IERC721MetadataTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _IERC721Metadata.contract.FilterLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &IERC721MetadataTransferIterator{contract: _IERC721Metadata.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: e Transfer(from indexed address, to indexed address, tokenId indexed uint256)
func (_IERC721Metadata *IERC721MetadataFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *IERC721MetadataTransfer, from []common.Address, to []common.Address, tokenId []*big.Int) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _IERC721Metadata.contract.WatchLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IERC721MetadataTransfer)
				if err := _IERC721Metadata.contract.UnpackLog(event, "Transfer", log); err != nil {
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

// IERC721ReceiverABI is the input ABI used to generate the binding from.
const IERC721ReceiverABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"operator\",\"type\":\"address\"},{\"name\":\"from\",\"type\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\"},{\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"onERC721Received\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes4\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// IERC721ReceiverBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const IERC721ReceiverBinRuntime = `0x`

// IERC721ReceiverBin is the compiled bytecode used for deploying new contracts.
const IERC721ReceiverBin = `0x`

// DeployIERC721Receiver deploys a new Klaytn contract, binding an instance of IERC721Receiver to it.
func DeployIERC721Receiver(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *IERC721Receiver, error) {
	parsed, err := abi.JSON(strings.NewReader(IERC721ReceiverABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(IERC721ReceiverBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &IERC721Receiver{IERC721ReceiverCaller: IERC721ReceiverCaller{contract: contract}, IERC721ReceiverTransactor: IERC721ReceiverTransactor{contract: contract}, IERC721ReceiverFilterer: IERC721ReceiverFilterer{contract: contract}}, nil
}

// IERC721Receiver is an auto generated Go binding around a Klaytn contract.
type IERC721Receiver struct {
	IERC721ReceiverCaller     // Read-only binding to the contract
	IERC721ReceiverTransactor // Write-only binding to the contract
	IERC721ReceiverFilterer   // Log filterer for contract events
}

// IERC721ReceiverCaller is an auto generated read-only Go binding around a Klaytn contract.
type IERC721ReceiverCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IERC721ReceiverTransactor is an auto generated write-only Go binding around a Klaytn contract.
type IERC721ReceiverTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IERC721ReceiverFilterer is an auto generated log filtering Go binding around a Klaytn contract events.
type IERC721ReceiverFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IERC721ReceiverSession is an auto generated Go binding around a Klaytn contract,
// with pre-set call and transact options.
type IERC721ReceiverSession struct {
	Contract     *IERC721Receiver  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IERC721ReceiverCallerSession is an auto generated read-only Go binding around a Klaytn contract,
// with pre-set call options.
type IERC721ReceiverCallerSession struct {
	Contract *IERC721ReceiverCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// IERC721ReceiverTransactorSession is an auto generated write-only Go binding around a Klaytn contract,
// with pre-set transact options.
type IERC721ReceiverTransactorSession struct {
	Contract     *IERC721ReceiverTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// IERC721ReceiverRaw is an auto generated low-level Go binding around a Klaytn contract.
type IERC721ReceiverRaw struct {
	Contract *IERC721Receiver // Generic contract binding to access the raw methods on
}

// IERC721ReceiverCallerRaw is an auto generated low-level read-only Go binding around a Klaytn contract.
type IERC721ReceiverCallerRaw struct {
	Contract *IERC721ReceiverCaller // Generic read-only contract binding to access the raw methods on
}

// IERC721ReceiverTransactorRaw is an auto generated low-level write-only Go binding around a Klaytn contract.
type IERC721ReceiverTransactorRaw struct {
	Contract *IERC721ReceiverTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIERC721Receiver creates a new instance of IERC721Receiver, bound to a specific deployed contract.
func NewIERC721Receiver(address common.Address, backend bind.ContractBackend) (*IERC721Receiver, error) {
	contract, err := bindIERC721Receiver(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IERC721Receiver{IERC721ReceiverCaller: IERC721ReceiverCaller{contract: contract}, IERC721ReceiverTransactor: IERC721ReceiverTransactor{contract: contract}, IERC721ReceiverFilterer: IERC721ReceiverFilterer{contract: contract}}, nil
}

// NewIERC721ReceiverCaller creates a new read-only instance of IERC721Receiver, bound to a specific deployed contract.
func NewIERC721ReceiverCaller(address common.Address, caller bind.ContractCaller) (*IERC721ReceiverCaller, error) {
	contract, err := bindIERC721Receiver(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IERC721ReceiverCaller{contract: contract}, nil
}

// NewIERC721ReceiverTransactor creates a new write-only instance of IERC721Receiver, bound to a specific deployed contract.
func NewIERC721ReceiverTransactor(address common.Address, transactor bind.ContractTransactor) (*IERC721ReceiverTransactor, error) {
	contract, err := bindIERC721Receiver(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IERC721ReceiverTransactor{contract: contract}, nil
}

// NewIERC721ReceiverFilterer creates a new log filterer instance of IERC721Receiver, bound to a specific deployed contract.
func NewIERC721ReceiverFilterer(address common.Address, filterer bind.ContractFilterer) (*IERC721ReceiverFilterer, error) {
	contract, err := bindIERC721Receiver(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IERC721ReceiverFilterer{contract: contract}, nil
}

// bindIERC721Receiver binds a generic wrapper to an already deployed contract.
func bindIERC721Receiver(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(IERC721ReceiverABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IERC721Receiver *IERC721ReceiverRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _IERC721Receiver.Contract.IERC721ReceiverCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IERC721Receiver *IERC721ReceiverRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IERC721Receiver.Contract.IERC721ReceiverTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IERC721Receiver *IERC721ReceiverRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IERC721Receiver.Contract.IERC721ReceiverTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IERC721Receiver *IERC721ReceiverCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _IERC721Receiver.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IERC721Receiver *IERC721ReceiverTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IERC721Receiver.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IERC721Receiver *IERC721ReceiverTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IERC721Receiver.Contract.contract.Transact(opts, method, params...)
}

// OnERC721Received is a paid mutator transaction binding the contract method 0x150b7a02.
//
// Solidity: function onERC721Received(operator address, from address, tokenId uint256, data bytes) returns(bytes4)
func (_IERC721Receiver *IERC721ReceiverTransactor) OnERC721Received(opts *bind.TransactOpts, operator common.Address, from common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _IERC721Receiver.contract.Transact(opts, "onERC721Received", operator, from, tokenId, data)
}

// OnERC721Received is a paid mutator transaction binding the contract method 0x150b7a02.
//
// Solidity: function onERC721Received(operator address, from address, tokenId uint256, data bytes) returns(bytes4)
func (_IERC721Receiver *IERC721ReceiverSession) OnERC721Received(operator common.Address, from common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _IERC721Receiver.Contract.OnERC721Received(&_IERC721Receiver.TransactOpts, operator, from, tokenId, data)
}

// OnERC721Received is a paid mutator transaction binding the contract method 0x150b7a02.
//
// Solidity: function onERC721Received(operator address, from address, tokenId uint256, data bytes) returns(bytes4)
func (_IERC721Receiver *IERC721ReceiverTransactorSession) OnERC721Received(operator common.Address, from common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _IERC721Receiver.Contract.OnERC721Received(&_IERC721Receiver.TransactOpts, operator, from, tokenId, data)
}

// MinterRoleABI is the input ABI used to generate the binding from.
const MinterRoleABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"account\",\"type\":\"address\"}],\"name\":\"addMinter\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"renounceMinter\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"account\",\"type\":\"address\"}],\"name\":\"isMinter\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"account\",\"type\":\"address\"}],\"name\":\"MinterAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"account\",\"type\":\"address\"}],\"name\":\"MinterRemoved\",\"type\":\"event\"}]"

// MinterRoleBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const MinterRoleBinRuntime = `0x`

// MinterRoleBin is the compiled bytecode used for deploying new contracts.
const MinterRoleBin = `0x`

// DeployMinterRole deploys a new Klaytn contract, binding an instance of MinterRole to it.
func DeployMinterRole(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *MinterRole, error) {
	parsed, err := abi.JSON(strings.NewReader(MinterRoleABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(MinterRoleBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &MinterRole{MinterRoleCaller: MinterRoleCaller{contract: contract}, MinterRoleTransactor: MinterRoleTransactor{contract: contract}, MinterRoleFilterer: MinterRoleFilterer{contract: contract}}, nil
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
// Solidity: function isMinter(account address) constant returns(bool)
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
// Solidity: function isMinter(account address) constant returns(bool)
func (_MinterRole *MinterRoleSession) IsMinter(account common.Address) (bool, error) {
	return _MinterRole.Contract.IsMinter(&_MinterRole.CallOpts, account)
}

// IsMinter is a free data retrieval call binding the contract method 0xaa271e1a.
//
// Solidity: function isMinter(account address) constant returns(bool)
func (_MinterRole *MinterRoleCallerSession) IsMinter(account common.Address) (bool, error) {
	return _MinterRole.Contract.IsMinter(&_MinterRole.CallOpts, account)
}

// AddMinter is a paid mutator transaction binding the contract method 0x983b2d56.
//
// Solidity: function addMinter(account address) returns()
func (_MinterRole *MinterRoleTransactor) AddMinter(opts *bind.TransactOpts, account common.Address) (*types.Transaction, error) {
	return _MinterRole.contract.Transact(opts, "addMinter", account)
}

// AddMinter is a paid mutator transaction binding the contract method 0x983b2d56.
//
// Solidity: function addMinter(account address) returns()
func (_MinterRole *MinterRoleSession) AddMinter(account common.Address) (*types.Transaction, error) {
	return _MinterRole.Contract.AddMinter(&_MinterRole.TransactOpts, account)
}

// AddMinter is a paid mutator transaction binding the contract method 0x983b2d56.
//
// Solidity: function addMinter(account address) returns()
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
// Solidity: e MinterAdded(account indexed address)
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
// Solidity: e MinterAdded(account indexed address)
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
// Solidity: e MinterRemoved(account indexed address)
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
// Solidity: e MinterRemoved(account indexed address)
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

// OwnableABI is the input ABI used to generate the binding from.
const OwnableABI = "[{\"constant\":false,\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"isOwner\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"}]"

// OwnableBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const OwnableBinRuntime = `0x`

// OwnableBin is the compiled bytecode used for deploying new contracts.
const OwnableBin = `0x`

// DeployOwnable deploys a new Klaytn contract, binding an instance of Ownable to it.
func DeployOwnable(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Ownable, error) {
	parsed, err := abi.JSON(strings.NewReader(OwnableABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(OwnableBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Ownable{OwnableCaller: OwnableCaller{contract: contract}, OwnableTransactor: OwnableTransactor{contract: contract}, OwnableFilterer: OwnableFilterer{contract: contract}}, nil
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
// Solidity: function isOwner() constant returns(bool)
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
// Solidity: function isOwner() constant returns(bool)
func (_Ownable *OwnableSession) IsOwner() (bool, error) {
	return _Ownable.Contract.IsOwner(&_Ownable.CallOpts)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() constant returns(bool)
func (_Ownable *OwnableCallerSession) IsOwner() (bool, error) {
	return _Ownable.Contract.IsOwner(&_Ownable.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
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
// Solidity: function owner() constant returns(address)
func (_Ownable *OwnableSession) Owner() (common.Address, error) {
	return _Ownable.Contract.Owner(&_Ownable.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
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
// Solidity: function transferOwnership(newOwner address) returns()
func (_Ownable *OwnableTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Ownable.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(newOwner address) returns()
func (_Ownable *OwnableSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Ownable.Contract.TransferOwnership(&_Ownable.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(newOwner address) returns()
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
// Solidity: e OwnershipTransferred(previousOwner indexed address, newOwner indexed address)
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
// Solidity: e OwnershipTransferred(previousOwner indexed address, newOwner indexed address)
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

// RolesABI is the input ABI used to generate the binding from.
const RolesABI = "[]"

// RolesBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const RolesBinRuntime = `0x73000000000000000000000000000000000000000030146080604052600080fd00a165627a7a72305820f16317d9d8c47a969121308df1f7c1d4e40ed1c79645f6236266f3e12766e1b80029`

// RolesBin is the compiled bytecode used for deploying new contracts.
const RolesBin = `0x604c602c600b82828239805160001a60731460008114601c57601e565bfe5b5030600052607381538281f30073000000000000000000000000000000000000000030146080604052600080fd00a165627a7a72305820f16317d9d8c47a969121308df1f7c1d4e40ed1c79645f6236266f3e12766e1b80029`

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
const SafeMathBinRuntime = `0x73000000000000000000000000000000000000000030146080604052600080fd00a165627a7a72305820f964156d21d4855472aa362aff743d0d42dd5e6fda9399e6ad50e2b4ec7bcad00029`

// SafeMathBin is the compiled bytecode used for deploying new contracts.
const SafeMathBin = `0x604c602c600b82828239805160001a60731460008114601c57601e565bfe5b5030600052607381538281f30073000000000000000000000000000000000000000030146080604052600080fd00a165627a7a72305820f964156d21d4855472aa362aff743d0d42dd5e6fda9399e6ad50e2b4ec7bcad00029`

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
