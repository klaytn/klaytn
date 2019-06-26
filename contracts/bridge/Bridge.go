// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package bridge

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

// BridgeABI is the input ABI used to generate the binding from.
const BridgeABI = "[{\"constant\":false,\"inputs\":[],\"name\":\"stop\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_amount\",\"type\":\"uint256\"},{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_contractAddress\",\"type\":\"address\"},{\"name\":\"_requestNonce\",\"type\":\"uint64\"},{\"name\":\"_requestBlockNumber\",\"type\":\"uint64\"}],\"name\":\"handleTokenTransfer\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"isRunning\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"counterpartBridge\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_token\",\"type\":\"address\"},{\"name\":\"_cToken\",\"type\":\"address\"}],\"name\":\"registerToken\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"from\",\"type\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\"},{\"name\":\"to\",\"type\":\"address\"}],\"name\":\"onNFTReceived\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes4\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"requestNonce\",\"outputs\":[{\"name\":\"\",\"type\":\"uint64\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_to\",\"type\":\"address\"}],\"name\":\"requestKLAYTransfer\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_bridge\",\"type\":\"address\"}],\"name\":\"setCounterPartBridge\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_uid\",\"type\":\"uint256\"},{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_contractAddress\",\"type\":\"address\"},{\"name\":\"_requestNonce\",\"type\":\"uint64\"},{\"name\":\"_requestBlockNumber\",\"type\":\"uint64\"}],\"name\":\"handleNFTTransfer\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"isOwner\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_token\",\"type\":\"address\"}],\"name\":\"deregisterToken\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"start\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"chargeWithoutEvent\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"allowedTokens\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_amount\",\"type\":\"uint256\"},{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_requestNonce\",\"type\":\"uint64\"},{\"name\":\"_requestBlockNumber\",\"type\":\"uint64\"}],\"name\":\"handleKLAYTransfer\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_from\",\"type\":\"address\"},{\"name\":\"_amount\",\"type\":\"uint256\"},{\"name\":\"_to\",\"type\":\"address\"}],\"name\":\"onTokenReceived\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes4\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"lastHandledRequestBlockNumber\",\"outputs\":[{\"name\":\"\",\"type\":\"uint64\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"handleNonce\",\"outputs\":[{\"name\":\"\",\"type\":\"uint64\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"VERSION\",\"outputs\":[{\"name\":\"\",\"type\":\"uint64\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"constructor\"},{\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"kind\",\"type\":\"uint8\"},{\"indexed\":false,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"contractAddress\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"requestNonce\",\"type\":\"uint64\"}],\"name\":\"RequestValueTransfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"kind\",\"type\":\"uint8\"},{\"indexed\":false,\"name\":\"contractAddress\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"handleNonce\",\"type\":\"uint64\"}],\"name\":\"HandleValueTransfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"}]"

// BridgeBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const BridgeBinRuntime = `0x6080604052600436106101325763ffffffff7c010000000000000000000000000000000000000000000000000000000060003504166307da68f581146102515780631d96fe1d146102685780632014e5d1146102a55780633a348533146102ce5780634739f7e5146102ff57806348f32f8814610326578063715018a6146103865780637c1a03021461039b578063879ae9a0146103cd57806387b04c55146103e15780638ce89ba8146104025780638da5cb5b1461043f5780638f32d59b14610454578063bab2af1d14610469578063be9a65551461048a578063dd9222d61461049f578063e744092e146104a7578063ed9c6da1146104c8578063f099d9bd146104ff578063f2fde38b1461052a578063f42b9aa11461054b578063fd07516914610560578063ffa1ad7414610575575b60015460a060020a900460ff161515610183576040805160e560020a62461bcd02815260206004820152600e60248201526000805160206112e1833981519152604482015290519081900360640190fd5b600034116101db576040805160e560020a62461bcd02815260206004820152600e60248201527f7a65726f206d73672e76616c7565000000000000000000000000000000000000604482015290519081900360640190fd5b600354604080516000808252336020830181905234838501526060830191909152608082015267ffffffffffffffff90921660a0830152516000805160206113018339815191529181900360c00190a16003805467ffffffffffffffff8082166001011667ffffffffffffffff19909116179055005b34801561025d57600080fd5b5061026661058a565b005b34801561027457600080fd5b50610266600435600160a060020a036024358116906044351667ffffffffffffffff606435811690608435166105bd565b3480156102b157600080fd5b506102ba6107e8565b604080519115158252519081900360200190f35b3480156102da57600080fd5b506102e36107f8565b60408051600160a060020a039092168252519081900360200190f35b34801561030b57600080fd5b50610266600160a060020a0360043581169060243516610807565b34801561033257600080fd5b50610351600160a060020a036004358116906024359060443516610855565b604080517fffffffff000000000000000000000000000000000000000000000000000000009092168252519081900360200190f35b34801561039257600080fd5b506102666109bd565b3480156103a757600080fd5b506103b0610a27565b6040805167ffffffffffffffff9092168252519081900360200190f35b610266600160a060020a0360043516610a37565b3480156103ed57600080fd5b50610266600160a060020a0360043516610b5c565b34801561040e57600080fd5b50610266600435600160a060020a036024358116906044351667ffffffffffffffff60643581169060843516610b9e565b34801561044b57600080fd5b506102e3610dba565b34801561046057600080fd5b506102ba610dc9565b34801561047557600080fd5b50610266600160a060020a0360043516610dda565b34801561049657600080fd5b50610266610e21565b610266610e5a565b3480156104b357600080fd5b506102e3600160a060020a0360043516610e5c565b3480156104d457600080fd5b50610266600435600160a060020a036024351667ffffffffffffffff60443581169060643516610e77565b34801561050b57600080fd5b50610351600160a060020a03600435811690602435906044351661103f565b34801561053657600080fd5b50610266600160a060020a03600435166111ff565b34801561055757600080fd5b506103b061121e565b34801561056c57600080fd5b506103b0611242565b34801561058157600080fd5b506103b061125e565b610592610dc9565b151561059d57600080fd5b6001805474ff000000000000000000000000000000000000000019169055565b6105c5610dc9565b15156105d057600080fd5b60035467ffffffffffffffff838116680100000000000000009092041614610668576040805160e560020a62461bcd02815260206004820152602160248201527f6d69736d6174636865642068616e646c65202f2072657175657374206e6f6e6360448201527f6500000000000000000000000000000000000000000000000000000000000000606482015290519081900360840190fd5b60035460408051600160a060020a03878116825260016020830152861681830152606081018890526801000000000000000090920467ffffffffffffffff166080830152517ff8bb02a7b538d2f29b0352ae5e0110751cf68d913d74b0a2b17fa06c73e4c8449181900360a00190a1600380546801000000000000000067ffffffffffffffff8085167001000000000000000000000000000000000277ffffffffffffffff0000000000000000000000000000000019909316929092178181048316600101909216026fffffffffffffffff000000000000000019909116179055604080517fa9059cbb000000000000000000000000000000000000000000000000000000008152600160a060020a0386811660048301526024820188905291519185169163a9059cbb916044808201926020929091908290030181600087803b1580156107b557600080fd5b505af11580156107c9573d6000803e3d6000fd5b505050506040513d60208110156107df57600080fd5b50505050505050565b60015460a060020a900460ff1681565b600154600160a060020a031681565b61080f610dc9565b151561081a57600080fd5b600160a060020a039182166000908152600260205260409020805473ffffffffffffffffffffffffffffffffffffffff191691909216179055565b60015460009060a060020a900460ff1615156108a9576040805160e560020a62461bcd02815260206004820152600e60248201526000805160206112e1833981519152604482015290519081900360640190fd5b33600090815260026020526040902054600160a060020a03161515610918576040805160e560020a62461bcd02815260206004820152601160248201527f4e6f7420612076616c696420746f6b656e000000000000000000000000000000604482015290519081900360640190fd5b6003546040805160028152600160a060020a0387811660208301528183018790523360608301528516608082015267ffffffffffffffff90921660a0830152516000805160206113018339815191529181900360c00190a1506003805467ffffffffffffffff8082166001011667ffffffffffffffff199091161790557f150b7a02000000000000000000000000000000000000000000000000000000009392505050565b6109c5610dc9565b15156109d057600080fd5b60008054604051600160a060020a03909116907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0908390a36000805473ffffffffffffffffffffffffffffffffffffffff19169055565b60035467ffffffffffffffff1681565b60015460a060020a900460ff161515610a88576040805160e560020a62461bcd02815260206004820152600e60248201526000805160206112e1833981519152604482015290519081900360640190fd5b60003411610ae0576040805160e560020a62461bcd02815260206004820152600e60248201527f7a65726f206d73672e76616c7565000000000000000000000000000000000000604482015290519081900360640190fd5b60035460408051600080825233602083015234828401526060820152600160a060020a038416608082015267ffffffffffffffff90921660a0830152516000805160206113018339815191529181900360c00190a1506003805467ffffffffffffffff8082166001011667ffffffffffffffff19909116179055565b610b64610dc9565b1515610b6f57600080fd5b6001805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a0392909216919091179055565b610ba6610dc9565b1515610bb157600080fd5b60035467ffffffffffffffff838116680100000000000000009092041614610c49576040805160e560020a62461bcd02815260206004820152602160248201527f6d69736d6174636865642068616e646c65202f2072657175657374206e6f6e6360448201527f6500000000000000000000000000000000000000000000000000000000000000606482015290519081900360840190fd5b60035460408051600160a060020a03878116825260026020830152861681830152606081018890526801000000000000000090920467ffffffffffffffff166080830152517ff8bb02a7b538d2f29b0352ae5e0110751cf68d913d74b0a2b17fa06c73e4c8449181900360a00190a1600380546801000000000000000067ffffffffffffffff8085167001000000000000000000000000000000000277ffffffffffffffff0000000000000000000000000000000019909316929092178181048316600101909216026fffffffffffffffff000000000000000019909116179055604080517f42842e0e000000000000000000000000000000000000000000000000000000008152306004820152600160a060020a038681166024830152604482018890529151918516916342842e0e9160648082019260009290919082900301818387803b158015610d9b57600080fd5b505af1158015610daf573d6000803e3d6000fd5b505050505050505050565b600054600160a060020a031690565b600054600160a060020a0316331490565b610de2610dc9565b1515610ded57600080fd5b600160a060020a03166000908152600260205260409020805473ffffffffffffffffffffffffffffffffffffffff19169055565b610e29610dc9565b1515610e3457600080fd5b6001805474ff0000000000000000000000000000000000000000191660a060020a179055565b565b600260205260009081526040902054600160a060020a031681565b610e7f610dc9565b1515610e8a57600080fd5b60035467ffffffffffffffff838116680100000000000000009092041614610f22576040805160e560020a62461bcd02815260206004820152602160248201527f6d69736d6174636865642068616e646c65202f2072657175657374206e6f6e6360448201527f6500000000000000000000000000000000000000000000000000000000000000606482015290519081900360840190fd5b60035460408051600160a060020a038616815260006020820181905281830152606081018790526801000000000000000090920467ffffffffffffffff166080830152517ff8bb02a7b538d2f29b0352ae5e0110751cf68d913d74b0a2b17fa06c73e4c8449181900360a00190a1600380546801000000000000000067ffffffffffffffff8085167001000000000000000000000000000000000277ffffffffffffffff0000000000000000000000000000000019909316929092178181048316600101909216026fffffffffffffffff000000000000000019909116179055604051600160a060020a0384169085156108fc029086906000818181858888f19350505050158015611038573d6000803e3d6000fd5b5050505050565b60015460009060a060020a900460ff161515611093576040805160e560020a62461bcd02815260206004820152600e60248201526000805160206112e1833981519152604482015290519081900360640190fd5b33600090815260026020526040902054600160a060020a03161515611102576040805160e560020a62461bcd02815260206004820152601160248201527f4e6f7420612076616c696420746f6b656e000000000000000000000000000000604482015290519081900360640190fd5b6000831161115a576040805160e560020a62461bcd02815260206004820152600b60248201527f7a65726f20616d6f756e74000000000000000000000000000000000000000000604482015290519081900360640190fd5b6003546040805160018152600160a060020a0387811660208301528183018790523360608301528516608082015267ffffffffffffffff90921660a0830152516000805160206113018339815191529181900360c00190a1506003805467ffffffffffffffff8082166001011667ffffffffffffffff199091161790557fbc04f0af000000000000000000000000000000000000000000000000000000009392505050565b611207610dc9565b151561121257600080fd5b61121b81611263565b50565b600354700100000000000000000000000000000000900467ffffffffffffffff1681565b60035468010000000000000000900467ffffffffffffffff1681565b600181565b600160a060020a038116151561127857600080fd5b60008054604051600160a060020a03808516939216917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a36000805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a0392909216919091179055560073746f70706564206272696467650000000000000000000000000000000000009cfefd96ef6f3968dfcb858820a8ee3b972faab26ddc7858691ffcd559beb005a165627a7a72305820f8bf433cda8b35df0c0400458b32e9ab200fec899bf800d5f1d90da425908cd60029`

// BridgeBin is the compiled bytecode used for deploying new contracts.
const BridgeBin = `0x6080604081905260008054600160a060020a0319163317808255600160a060020a0316917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0908290a36001805460a060020a60ff0219167401000000000000000000000000000000000000000017905561134c8061007e6000396000f3006080604052600436106101325763ffffffff7c010000000000000000000000000000000000000000000000000000000060003504166307da68f581146102515780631d96fe1d146102685780632014e5d1146102a55780633a348533146102ce5780634739f7e5146102ff57806348f32f8814610326578063715018a6146103865780637c1a03021461039b578063879ae9a0146103cd57806387b04c55146103e15780638ce89ba8146104025780638da5cb5b1461043f5780638f32d59b14610454578063bab2af1d14610469578063be9a65551461048a578063dd9222d61461049f578063e744092e146104a7578063ed9c6da1146104c8578063f099d9bd146104ff578063f2fde38b1461052a578063f42b9aa11461054b578063fd07516914610560578063ffa1ad7414610575575b60015460a060020a900460ff161515610183576040805160e560020a62461bcd02815260206004820152600e60248201526000805160206112e1833981519152604482015290519081900360640190fd5b600034116101db576040805160e560020a62461bcd02815260206004820152600e60248201527f7a65726f206d73672e76616c7565000000000000000000000000000000000000604482015290519081900360640190fd5b600354604080516000808252336020830181905234838501526060830191909152608082015267ffffffffffffffff90921660a0830152516000805160206113018339815191529181900360c00190a16003805467ffffffffffffffff8082166001011667ffffffffffffffff19909116179055005b34801561025d57600080fd5b5061026661058a565b005b34801561027457600080fd5b50610266600435600160a060020a036024358116906044351667ffffffffffffffff606435811690608435166105bd565b3480156102b157600080fd5b506102ba6107e8565b604080519115158252519081900360200190f35b3480156102da57600080fd5b506102e36107f8565b60408051600160a060020a039092168252519081900360200190f35b34801561030b57600080fd5b50610266600160a060020a0360043581169060243516610807565b34801561033257600080fd5b50610351600160a060020a036004358116906024359060443516610855565b604080517fffffffff000000000000000000000000000000000000000000000000000000009092168252519081900360200190f35b34801561039257600080fd5b506102666109bd565b3480156103a757600080fd5b506103b0610a27565b6040805167ffffffffffffffff9092168252519081900360200190f35b610266600160a060020a0360043516610a37565b3480156103ed57600080fd5b50610266600160a060020a0360043516610b5c565b34801561040e57600080fd5b50610266600435600160a060020a036024358116906044351667ffffffffffffffff60643581169060843516610b9e565b34801561044b57600080fd5b506102e3610dba565b34801561046057600080fd5b506102ba610dc9565b34801561047557600080fd5b50610266600160a060020a0360043516610dda565b34801561049657600080fd5b50610266610e21565b610266610e5a565b3480156104b357600080fd5b506102e3600160a060020a0360043516610e5c565b3480156104d457600080fd5b50610266600435600160a060020a036024351667ffffffffffffffff60443581169060643516610e77565b34801561050b57600080fd5b50610351600160a060020a03600435811690602435906044351661103f565b34801561053657600080fd5b50610266600160a060020a03600435166111ff565b34801561055757600080fd5b506103b061121e565b34801561056c57600080fd5b506103b0611242565b34801561058157600080fd5b506103b061125e565b610592610dc9565b151561059d57600080fd5b6001805474ff000000000000000000000000000000000000000019169055565b6105c5610dc9565b15156105d057600080fd5b60035467ffffffffffffffff838116680100000000000000009092041614610668576040805160e560020a62461bcd02815260206004820152602160248201527f6d69736d6174636865642068616e646c65202f2072657175657374206e6f6e6360448201527f6500000000000000000000000000000000000000000000000000000000000000606482015290519081900360840190fd5b60035460408051600160a060020a03878116825260016020830152861681830152606081018890526801000000000000000090920467ffffffffffffffff166080830152517ff8bb02a7b538d2f29b0352ae5e0110751cf68d913d74b0a2b17fa06c73e4c8449181900360a00190a1600380546801000000000000000067ffffffffffffffff8085167001000000000000000000000000000000000277ffffffffffffffff0000000000000000000000000000000019909316929092178181048316600101909216026fffffffffffffffff000000000000000019909116179055604080517fa9059cbb000000000000000000000000000000000000000000000000000000008152600160a060020a0386811660048301526024820188905291519185169163a9059cbb916044808201926020929091908290030181600087803b1580156107b557600080fd5b505af11580156107c9573d6000803e3d6000fd5b505050506040513d60208110156107df57600080fd5b50505050505050565b60015460a060020a900460ff1681565b600154600160a060020a031681565b61080f610dc9565b151561081a57600080fd5b600160a060020a039182166000908152600260205260409020805473ffffffffffffffffffffffffffffffffffffffff191691909216179055565b60015460009060a060020a900460ff1615156108a9576040805160e560020a62461bcd02815260206004820152600e60248201526000805160206112e1833981519152604482015290519081900360640190fd5b33600090815260026020526040902054600160a060020a03161515610918576040805160e560020a62461bcd02815260206004820152601160248201527f4e6f7420612076616c696420746f6b656e000000000000000000000000000000604482015290519081900360640190fd5b6003546040805160028152600160a060020a0387811660208301528183018790523360608301528516608082015267ffffffffffffffff90921660a0830152516000805160206113018339815191529181900360c00190a1506003805467ffffffffffffffff8082166001011667ffffffffffffffff199091161790557f150b7a02000000000000000000000000000000000000000000000000000000009392505050565b6109c5610dc9565b15156109d057600080fd5b60008054604051600160a060020a03909116907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0908390a36000805473ffffffffffffffffffffffffffffffffffffffff19169055565b60035467ffffffffffffffff1681565b60015460a060020a900460ff161515610a88576040805160e560020a62461bcd02815260206004820152600e60248201526000805160206112e1833981519152604482015290519081900360640190fd5b60003411610ae0576040805160e560020a62461bcd02815260206004820152600e60248201527f7a65726f206d73672e76616c7565000000000000000000000000000000000000604482015290519081900360640190fd5b60035460408051600080825233602083015234828401526060820152600160a060020a038416608082015267ffffffffffffffff90921660a0830152516000805160206113018339815191529181900360c00190a1506003805467ffffffffffffffff8082166001011667ffffffffffffffff19909116179055565b610b64610dc9565b1515610b6f57600080fd5b6001805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a0392909216919091179055565b610ba6610dc9565b1515610bb157600080fd5b60035467ffffffffffffffff838116680100000000000000009092041614610c49576040805160e560020a62461bcd02815260206004820152602160248201527f6d69736d6174636865642068616e646c65202f2072657175657374206e6f6e6360448201527f6500000000000000000000000000000000000000000000000000000000000000606482015290519081900360840190fd5b60035460408051600160a060020a03878116825260026020830152861681830152606081018890526801000000000000000090920467ffffffffffffffff166080830152517ff8bb02a7b538d2f29b0352ae5e0110751cf68d913d74b0a2b17fa06c73e4c8449181900360a00190a1600380546801000000000000000067ffffffffffffffff8085167001000000000000000000000000000000000277ffffffffffffffff0000000000000000000000000000000019909316929092178181048316600101909216026fffffffffffffffff000000000000000019909116179055604080517f42842e0e000000000000000000000000000000000000000000000000000000008152306004820152600160a060020a038681166024830152604482018890529151918516916342842e0e9160648082019260009290919082900301818387803b158015610d9b57600080fd5b505af1158015610daf573d6000803e3d6000fd5b505050505050505050565b600054600160a060020a031690565b600054600160a060020a0316331490565b610de2610dc9565b1515610ded57600080fd5b600160a060020a03166000908152600260205260409020805473ffffffffffffffffffffffffffffffffffffffff19169055565b610e29610dc9565b1515610e3457600080fd5b6001805474ff0000000000000000000000000000000000000000191660a060020a179055565b565b600260205260009081526040902054600160a060020a031681565b610e7f610dc9565b1515610e8a57600080fd5b60035467ffffffffffffffff838116680100000000000000009092041614610f22576040805160e560020a62461bcd02815260206004820152602160248201527f6d69736d6174636865642068616e646c65202f2072657175657374206e6f6e6360448201527f6500000000000000000000000000000000000000000000000000000000000000606482015290519081900360840190fd5b60035460408051600160a060020a038616815260006020820181905281830152606081018790526801000000000000000090920467ffffffffffffffff166080830152517ff8bb02a7b538d2f29b0352ae5e0110751cf68d913d74b0a2b17fa06c73e4c8449181900360a00190a1600380546801000000000000000067ffffffffffffffff8085167001000000000000000000000000000000000277ffffffffffffffff0000000000000000000000000000000019909316929092178181048316600101909216026fffffffffffffffff000000000000000019909116179055604051600160a060020a0384169085156108fc029086906000818181858888f19350505050158015611038573d6000803e3d6000fd5b5050505050565b60015460009060a060020a900460ff161515611093576040805160e560020a62461bcd02815260206004820152600e60248201526000805160206112e1833981519152604482015290519081900360640190fd5b33600090815260026020526040902054600160a060020a03161515611102576040805160e560020a62461bcd02815260206004820152601160248201527f4e6f7420612076616c696420746f6b656e000000000000000000000000000000604482015290519081900360640190fd5b6000831161115a576040805160e560020a62461bcd02815260206004820152600b60248201527f7a65726f20616d6f756e74000000000000000000000000000000000000000000604482015290519081900360640190fd5b6003546040805160018152600160a060020a0387811660208301528183018790523360608301528516608082015267ffffffffffffffff90921660a0830152516000805160206113018339815191529181900360c00190a1506003805467ffffffffffffffff8082166001011667ffffffffffffffff199091161790557fbc04f0af000000000000000000000000000000000000000000000000000000009392505050565b611207610dc9565b151561121257600080fd5b61121b81611263565b50565b600354700100000000000000000000000000000000900467ffffffffffffffff1681565b60035468010000000000000000900467ffffffffffffffff1681565b600181565b600160a060020a038116151561127857600080fd5b60008054604051600160a060020a03808516939216917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a36000805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a0392909216919091179055560073746f70706564206272696467650000000000000000000000000000000000009cfefd96ef6f3968dfcb858820a8ee3b972faab26ddc7858691ffcd559beb005a165627a7a72305820f8bf433cda8b35df0c0400458b32e9ab200fec899bf800d5f1d90da425908cd60029`

// DeployBridge deploys a new Klaytn contract, binding an instance of Bridge to it.
func DeployBridge(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Bridge, error) {
	parsed, err := abi.JSON(strings.NewReader(BridgeABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(BridgeBin), backend)
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

// HandleNonce is a free data retrieval call binding the contract method 0xfd075169.
//
// Solidity: function handleNonce() constant returns(uint64)
func (_Bridge *BridgeCaller) HandleNonce(opts *bind.CallOpts) (uint64, error) {
	var (
		ret0 = new(uint64)
	)
	out := ret0
	err := _Bridge.contract.Call(opts, out, "handleNonce")
	return *ret0, err
}

// HandleNonce is a free data retrieval call binding the contract method 0xfd075169.
//
// Solidity: function handleNonce() constant returns(uint64)
func (_Bridge *BridgeSession) HandleNonce() (uint64, error) {
	return _Bridge.Contract.HandleNonce(&_Bridge.CallOpts)
}

// HandleNonce is a free data retrieval call binding the contract method 0xfd075169.
//
// Solidity: function handleNonce() constant returns(uint64)
func (_Bridge *BridgeCallerSession) HandleNonce() (uint64, error) {
	return _Bridge.Contract.HandleNonce(&_Bridge.CallOpts)
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

// LastHandledRequestBlockNumber is a free data retrieval call binding the contract method 0xf42b9aa1.
//
// Solidity: function lastHandledRequestBlockNumber() constant returns(uint64)
func (_Bridge *BridgeCaller) LastHandledRequestBlockNumber(opts *bind.CallOpts) (uint64, error) {
	var (
		ret0 = new(uint64)
	)
	out := ret0
	err := _Bridge.contract.Call(opts, out, "lastHandledRequestBlockNumber")
	return *ret0, err
}

// LastHandledRequestBlockNumber is a free data retrieval call binding the contract method 0xf42b9aa1.
//
// Solidity: function lastHandledRequestBlockNumber() constant returns(uint64)
func (_Bridge *BridgeSession) LastHandledRequestBlockNumber() (uint64, error) {
	return _Bridge.Contract.LastHandledRequestBlockNumber(&_Bridge.CallOpts)
}

// LastHandledRequestBlockNumber is a free data retrieval call binding the contract method 0xf42b9aa1.
//
// Solidity: function lastHandledRequestBlockNumber() constant returns(uint64)
func (_Bridge *BridgeCallerSession) LastHandledRequestBlockNumber() (uint64, error) {
	return _Bridge.Contract.LastHandledRequestBlockNumber(&_Bridge.CallOpts)
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

// HandleKLAYTransfer is a paid mutator transaction binding the contract method 0xed9c6da1.
//
// Solidity: function handleKLAYTransfer(_amount uint256, _to address, _requestNonce uint64, _requestBlockNumber uint64) returns()
func (_Bridge *BridgeTransactor) HandleKLAYTransfer(opts *bind.TransactOpts, _amount *big.Int, _to common.Address, _requestNonce uint64, _requestBlockNumber uint64) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "handleKLAYTransfer", _amount, _to, _requestNonce, _requestBlockNumber)
}

// HandleKLAYTransfer is a paid mutator transaction binding the contract method 0xed9c6da1.
//
// Solidity: function handleKLAYTransfer(_amount uint256, _to address, _requestNonce uint64, _requestBlockNumber uint64) returns()
func (_Bridge *BridgeSession) HandleKLAYTransfer(_amount *big.Int, _to common.Address, _requestNonce uint64, _requestBlockNumber uint64) (*types.Transaction, error) {
	return _Bridge.Contract.HandleKLAYTransfer(&_Bridge.TransactOpts, _amount, _to, _requestNonce, _requestBlockNumber)
}

// HandleKLAYTransfer is a paid mutator transaction binding the contract method 0xed9c6da1.
//
// Solidity: function handleKLAYTransfer(_amount uint256, _to address, _requestNonce uint64, _requestBlockNumber uint64) returns()
func (_Bridge *BridgeTransactorSession) HandleKLAYTransfer(_amount *big.Int, _to common.Address, _requestNonce uint64, _requestBlockNumber uint64) (*types.Transaction, error) {
	return _Bridge.Contract.HandleKLAYTransfer(&_Bridge.TransactOpts, _amount, _to, _requestNonce, _requestBlockNumber)
}

// HandleNFTTransfer is a paid mutator transaction binding the contract method 0x8ce89ba8.
//
// Solidity: function handleNFTTransfer(_uid uint256, _to address, _contractAddress address, _requestNonce uint64, _requestBlockNumber uint64) returns()
func (_Bridge *BridgeTransactor) HandleNFTTransfer(opts *bind.TransactOpts, _uid *big.Int, _to common.Address, _contractAddress common.Address, _requestNonce uint64, _requestBlockNumber uint64) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "handleNFTTransfer", _uid, _to, _contractAddress, _requestNonce, _requestBlockNumber)
}

// HandleNFTTransfer is a paid mutator transaction binding the contract method 0x8ce89ba8.
//
// Solidity: function handleNFTTransfer(_uid uint256, _to address, _contractAddress address, _requestNonce uint64, _requestBlockNumber uint64) returns()
func (_Bridge *BridgeSession) HandleNFTTransfer(_uid *big.Int, _to common.Address, _contractAddress common.Address, _requestNonce uint64, _requestBlockNumber uint64) (*types.Transaction, error) {
	return _Bridge.Contract.HandleNFTTransfer(&_Bridge.TransactOpts, _uid, _to, _contractAddress, _requestNonce, _requestBlockNumber)
}

// HandleNFTTransfer is a paid mutator transaction binding the contract method 0x8ce89ba8.
//
// Solidity: function handleNFTTransfer(_uid uint256, _to address, _contractAddress address, _requestNonce uint64, _requestBlockNumber uint64) returns()
func (_Bridge *BridgeTransactorSession) HandleNFTTransfer(_uid *big.Int, _to common.Address, _contractAddress common.Address, _requestNonce uint64, _requestBlockNumber uint64) (*types.Transaction, error) {
	return _Bridge.Contract.HandleNFTTransfer(&_Bridge.TransactOpts, _uid, _to, _contractAddress, _requestNonce, _requestBlockNumber)
}

// HandleTokenTransfer is a paid mutator transaction binding the contract method 0x1d96fe1d.
//
// Solidity: function handleTokenTransfer(_amount uint256, _to address, _contractAddress address, _requestNonce uint64, _requestBlockNumber uint64) returns()
func (_Bridge *BridgeTransactor) HandleTokenTransfer(opts *bind.TransactOpts, _amount *big.Int, _to common.Address, _contractAddress common.Address, _requestNonce uint64, _requestBlockNumber uint64) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "handleTokenTransfer", _amount, _to, _contractAddress, _requestNonce, _requestBlockNumber)
}

// HandleTokenTransfer is a paid mutator transaction binding the contract method 0x1d96fe1d.
//
// Solidity: function handleTokenTransfer(_amount uint256, _to address, _contractAddress address, _requestNonce uint64, _requestBlockNumber uint64) returns()
func (_Bridge *BridgeSession) HandleTokenTransfer(_amount *big.Int, _to common.Address, _contractAddress common.Address, _requestNonce uint64, _requestBlockNumber uint64) (*types.Transaction, error) {
	return _Bridge.Contract.HandleTokenTransfer(&_Bridge.TransactOpts, _amount, _to, _contractAddress, _requestNonce, _requestBlockNumber)
}

// HandleTokenTransfer is a paid mutator transaction binding the contract method 0x1d96fe1d.
//
// Solidity: function handleTokenTransfer(_amount uint256, _to address, _contractAddress address, _requestNonce uint64, _requestBlockNumber uint64) returns()
func (_Bridge *BridgeTransactorSession) HandleTokenTransfer(_amount *big.Int, _to common.Address, _contractAddress common.Address, _requestNonce uint64, _requestBlockNumber uint64) (*types.Transaction, error) {
	return _Bridge.Contract.HandleTokenTransfer(&_Bridge.TransactOpts, _amount, _to, _contractAddress, _requestNonce, _requestBlockNumber)
}

// OnNFTReceived is a paid mutator transaction binding the contract method 0x48f32f88.
//
// Solidity: function onNFTReceived(from address, tokenId uint256, to address) returns(bytes4)
func (_Bridge *BridgeTransactor) OnNFTReceived(opts *bind.TransactOpts, from common.Address, tokenId *big.Int, to common.Address) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "onNFTReceived", from, tokenId, to)
}

// OnNFTReceived is a paid mutator transaction binding the contract method 0x48f32f88.
//
// Solidity: function onNFTReceived(from address, tokenId uint256, to address) returns(bytes4)
func (_Bridge *BridgeSession) OnNFTReceived(from common.Address, tokenId *big.Int, to common.Address) (*types.Transaction, error) {
	return _Bridge.Contract.OnNFTReceived(&_Bridge.TransactOpts, from, tokenId, to)
}

// OnNFTReceived is a paid mutator transaction binding the contract method 0x48f32f88.
//
// Solidity: function onNFTReceived(from address, tokenId uint256, to address) returns(bytes4)
func (_Bridge *BridgeTransactorSession) OnNFTReceived(from common.Address, tokenId *big.Int, to common.Address) (*types.Transaction, error) {
	return _Bridge.Contract.OnNFTReceived(&_Bridge.TransactOpts, from, tokenId, to)
}

// OnTokenReceived is a paid mutator transaction binding the contract method 0xf099d9bd.
//
// Solidity: function onTokenReceived(_from address, _amount uint256, _to address) returns(bytes4)
func (_Bridge *BridgeTransactor) OnTokenReceived(opts *bind.TransactOpts, _from common.Address, _amount *big.Int, _to common.Address) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "onTokenReceived", _from, _amount, _to)
}

// OnTokenReceived is a paid mutator transaction binding the contract method 0xf099d9bd.
//
// Solidity: function onTokenReceived(_from address, _amount uint256, _to address) returns(bytes4)
func (_Bridge *BridgeSession) OnTokenReceived(_from common.Address, _amount *big.Int, _to common.Address) (*types.Transaction, error) {
	return _Bridge.Contract.OnTokenReceived(&_Bridge.TransactOpts, _from, _amount, _to)
}

// OnTokenReceived is a paid mutator transaction binding the contract method 0xf099d9bd.
//
// Solidity: function onTokenReceived(_from address, _amount uint256, _to address) returns(bytes4)
func (_Bridge *BridgeTransactorSession) OnTokenReceived(_from common.Address, _amount *big.Int, _to common.Address) (*types.Transaction, error) {
	return _Bridge.Contract.OnTokenReceived(&_Bridge.TransactOpts, _from, _amount, _to)
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

// RequestKLAYTransfer is a paid mutator transaction binding the contract method 0x879ae9a0.
//
// Solidity: function requestKLAYTransfer(_to address) returns()
func (_Bridge *BridgeTransactor) RequestKLAYTransfer(opts *bind.TransactOpts, _to common.Address) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "requestKLAYTransfer", _to)
}

// RequestKLAYTransfer is a paid mutator transaction binding the contract method 0x879ae9a0.
//
// Solidity: function requestKLAYTransfer(_to address) returns()
func (_Bridge *BridgeSession) RequestKLAYTransfer(_to common.Address) (*types.Transaction, error) {
	return _Bridge.Contract.RequestKLAYTransfer(&_Bridge.TransactOpts, _to)
}

// RequestKLAYTransfer is a paid mutator transaction binding the contract method 0x879ae9a0.
//
// Solidity: function requestKLAYTransfer(_to address) returns()
func (_Bridge *BridgeTransactorSession) RequestKLAYTransfer(_to common.Address) (*types.Transaction, error) {
	return _Bridge.Contract.RequestKLAYTransfer(&_Bridge.TransactOpts, _to)
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

// Start is a paid mutator transaction binding the contract method 0xbe9a6555.
//
// Solidity: function start() returns()
func (_Bridge *BridgeTransactor) Start(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "start")
}

// Start is a paid mutator transaction binding the contract method 0xbe9a6555.
//
// Solidity: function start() returns()
func (_Bridge *BridgeSession) Start() (*types.Transaction, error) {
	return _Bridge.Contract.Start(&_Bridge.TransactOpts)
}

// Start is a paid mutator transaction binding the contract method 0xbe9a6555.
//
// Solidity: function start() returns()
func (_Bridge *BridgeTransactorSession) Start() (*types.Transaction, error) {
	return _Bridge.Contract.Start(&_Bridge.TransactOpts)
}

// Stop is a paid mutator transaction binding the contract method 0x07da68f5.
//
// Solidity: function stop() returns()
func (_Bridge *BridgeTransactor) Stop(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "stop")
}

// Stop is a paid mutator transaction binding the contract method 0x07da68f5.
//
// Solidity: function stop() returns()
func (_Bridge *BridgeSession) Stop() (*types.Transaction, error) {
	return _Bridge.Contract.Stop(&_Bridge.TransactOpts)
}

// Stop is a paid mutator transaction binding the contract method 0x07da68f5.
//
// Solidity: function stop() returns()
func (_Bridge *BridgeTransactorSession) Stop() (*types.Transaction, error) {
	return _Bridge.Contract.Stop(&_Bridge.TransactOpts)
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
	Owner           common.Address
	Kind            uint8
	ContractAddress common.Address
	Value           *big.Int
	HandleNonce     uint64
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterHandleValueTransfer is a free log retrieval operation binding the contract event 0xf8bb02a7b538d2f29b0352ae5e0110751cf68d913d74b0a2b17fa06c73e4c844.
//
// Solidity: e HandleValueTransfer(owner address, kind uint8, contractAddress address, value uint256, handleNonce uint64)
func (_Bridge *BridgeFilterer) FilterHandleValueTransfer(opts *bind.FilterOpts) (*BridgeHandleValueTransferIterator, error) {

	logs, sub, err := _Bridge.contract.FilterLogs(opts, "HandleValueTransfer")
	if err != nil {
		return nil, err
	}
	return &BridgeHandleValueTransferIterator{contract: _Bridge.contract, event: "HandleValueTransfer", logs: logs, sub: sub}, nil
}

// WatchHandleValueTransfer is a free log subscription operation binding the contract event 0xf8bb02a7b538d2f29b0352ae5e0110751cf68d913d74b0a2b17fa06c73e4c844.
//
// Solidity: e HandleValueTransfer(owner address, kind uint8, contractAddress address, value uint256, handleNonce uint64)
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
	Kind            uint8
	From            common.Address
	Amount          *big.Int
	ContractAddress common.Address
	To              common.Address
	RequestNonce    uint64
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterRequestValueTransfer is a free log retrieval operation binding the contract event 0x9cfefd96ef6f3968dfcb858820a8ee3b972faab26ddc7858691ffcd559beb005.
//
// Solidity: e RequestValueTransfer(kind uint8, from address, amount uint256, contractAddress address, to address, requestNonce uint64)
func (_Bridge *BridgeFilterer) FilterRequestValueTransfer(opts *bind.FilterOpts) (*BridgeRequestValueTransferIterator, error) {

	logs, sub, err := _Bridge.contract.FilterLogs(opts, "RequestValueTransfer")
	if err != nil {
		return nil, err
	}
	return &BridgeRequestValueTransferIterator{contract: _Bridge.contract, event: "RequestValueTransfer", logs: logs, sub: sub}, nil
}

// WatchRequestValueTransfer is a free log subscription operation binding the contract event 0x9cfefd96ef6f3968dfcb858820a8ee3b972faab26ddc7858691ffcd559beb005.
//
// Solidity: e RequestValueTransfer(kind uint8, from address, amount uint256, contractAddress address, to address, requestNonce uint64)
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

// INFTReceiverABI is the input ABI used to generate the binding from.
const INFTReceiverABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"from\",\"type\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\"},{\"name\":\"to\",\"type\":\"address\"}],\"name\":\"onNFTReceived\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes4\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// INFTReceiverBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const INFTReceiverBinRuntime = `0x`

// INFTReceiverBin is the compiled bytecode used for deploying new contracts.
const INFTReceiverBin = `0x`

// DeployINFTReceiver deploys a new Klaytn contract, binding an instance of INFTReceiver to it.
func DeployINFTReceiver(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *INFTReceiver, error) {
	parsed, err := abi.JSON(strings.NewReader(INFTReceiverABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(INFTReceiverBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &INFTReceiver{INFTReceiverCaller: INFTReceiverCaller{contract: contract}, INFTReceiverTransactor: INFTReceiverTransactor{contract: contract}, INFTReceiverFilterer: INFTReceiverFilterer{contract: contract}}, nil
}

// INFTReceiver is an auto generated Go binding around a Klaytn contract.
type INFTReceiver struct {
	INFTReceiverCaller     // Read-only binding to the contract
	INFTReceiverTransactor // Write-only binding to the contract
	INFTReceiverFilterer   // Log filterer for contract events
}

// INFTReceiverCaller is an auto generated read-only Go binding around a Klaytn contract.
type INFTReceiverCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// INFTReceiverTransactor is an auto generated write-only Go binding around a Klaytn contract.
type INFTReceiverTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// INFTReceiverFilterer is an auto generated log filtering Go binding around a Klaytn contract events.
type INFTReceiverFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// INFTReceiverSession is an auto generated Go binding around a Klaytn contract,
// with pre-set call and transact options.
type INFTReceiverSession struct {
	Contract     *INFTReceiver     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// INFTReceiverCallerSession is an auto generated read-only Go binding around a Klaytn contract,
// with pre-set call options.
type INFTReceiverCallerSession struct {
	Contract *INFTReceiverCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// INFTReceiverTransactorSession is an auto generated write-only Go binding around a Klaytn contract,
// with pre-set transact options.
type INFTReceiverTransactorSession struct {
	Contract     *INFTReceiverTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// INFTReceiverRaw is an auto generated low-level Go binding around a Klaytn contract.
type INFTReceiverRaw struct {
	Contract *INFTReceiver // Generic contract binding to access the raw methods on
}

// INFTReceiverCallerRaw is an auto generated low-level read-only Go binding around a Klaytn contract.
type INFTReceiverCallerRaw struct {
	Contract *INFTReceiverCaller // Generic read-only contract binding to access the raw methods on
}

// INFTReceiverTransactorRaw is an auto generated low-level write-only Go binding around a Klaytn contract.
type INFTReceiverTransactorRaw struct {
	Contract *INFTReceiverTransactor // Generic write-only contract binding to access the raw methods on
}

// NewINFTReceiver creates a new instance of INFTReceiver, bound to a specific deployed contract.
func NewINFTReceiver(address common.Address, backend bind.ContractBackend) (*INFTReceiver, error) {
	contract, err := bindINFTReceiver(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &INFTReceiver{INFTReceiverCaller: INFTReceiverCaller{contract: contract}, INFTReceiverTransactor: INFTReceiverTransactor{contract: contract}, INFTReceiverFilterer: INFTReceiverFilterer{contract: contract}}, nil
}

// NewINFTReceiverCaller creates a new read-only instance of INFTReceiver, bound to a specific deployed contract.
func NewINFTReceiverCaller(address common.Address, caller bind.ContractCaller) (*INFTReceiverCaller, error) {
	contract, err := bindINFTReceiver(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &INFTReceiverCaller{contract: contract}, nil
}

// NewINFTReceiverTransactor creates a new write-only instance of INFTReceiver, bound to a specific deployed contract.
func NewINFTReceiverTransactor(address common.Address, transactor bind.ContractTransactor) (*INFTReceiverTransactor, error) {
	contract, err := bindINFTReceiver(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &INFTReceiverTransactor{contract: contract}, nil
}

// NewINFTReceiverFilterer creates a new log filterer instance of INFTReceiver, bound to a specific deployed contract.
func NewINFTReceiverFilterer(address common.Address, filterer bind.ContractFilterer) (*INFTReceiverFilterer, error) {
	contract, err := bindINFTReceiver(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &INFTReceiverFilterer{contract: contract}, nil
}

// bindINFTReceiver binds a generic wrapper to an already deployed contract.
func bindINFTReceiver(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(INFTReceiverABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_INFTReceiver *INFTReceiverRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _INFTReceiver.Contract.INFTReceiverCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_INFTReceiver *INFTReceiverRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _INFTReceiver.Contract.INFTReceiverTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_INFTReceiver *INFTReceiverRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _INFTReceiver.Contract.INFTReceiverTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_INFTReceiver *INFTReceiverCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _INFTReceiver.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_INFTReceiver *INFTReceiverTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _INFTReceiver.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_INFTReceiver *INFTReceiverTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _INFTReceiver.Contract.contract.Transact(opts, method, params...)
}

// OnNFTReceived is a paid mutator transaction binding the contract method 0x48f32f88.
//
// Solidity: function onNFTReceived(from address, tokenId uint256, to address) returns(bytes4)
func (_INFTReceiver *INFTReceiverTransactor) OnNFTReceived(opts *bind.TransactOpts, from common.Address, tokenId *big.Int, to common.Address) (*types.Transaction, error) {
	return _INFTReceiver.contract.Transact(opts, "onNFTReceived", from, tokenId, to)
}

// OnNFTReceived is a paid mutator transaction binding the contract method 0x48f32f88.
//
// Solidity: function onNFTReceived(from address, tokenId uint256, to address) returns(bytes4)
func (_INFTReceiver *INFTReceiverSession) OnNFTReceived(from common.Address, tokenId *big.Int, to common.Address) (*types.Transaction, error) {
	return _INFTReceiver.Contract.OnNFTReceived(&_INFTReceiver.TransactOpts, from, tokenId, to)
}

// OnNFTReceived is a paid mutator transaction binding the contract method 0x48f32f88.
//
// Solidity: function onNFTReceived(from address, tokenId uint256, to address) returns(bytes4)
func (_INFTReceiver *INFTReceiverTransactorSession) OnNFTReceived(from common.Address, tokenId *big.Int, to common.Address) (*types.Transaction, error) {
	return _INFTReceiver.Contract.OnNFTReceived(&_INFTReceiver.TransactOpts, from, tokenId, to)
}

// ITokenReceiverABI is the input ABI used to generate the binding from.
const ITokenReceiverABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"_from\",\"type\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\"},{\"name\":\"_to\",\"type\":\"address\"}],\"name\":\"onTokenReceived\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes4\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// ITokenReceiverBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const ITokenReceiverBinRuntime = `0x`

// ITokenReceiverBin is the compiled bytecode used for deploying new contracts.
const ITokenReceiverBin = `0x`

// DeployITokenReceiver deploys a new Klaytn contract, binding an instance of ITokenReceiver to it.
func DeployITokenReceiver(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ITokenReceiver, error) {
	parsed, err := abi.JSON(strings.NewReader(ITokenReceiverABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(ITokenReceiverBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ITokenReceiver{ITokenReceiverCaller: ITokenReceiverCaller{contract: contract}, ITokenReceiverTransactor: ITokenReceiverTransactor{contract: contract}, ITokenReceiverFilterer: ITokenReceiverFilterer{contract: contract}}, nil
}

// ITokenReceiver is an auto generated Go binding around a Klaytn contract.
type ITokenReceiver struct {
	ITokenReceiverCaller     // Read-only binding to the contract
	ITokenReceiverTransactor // Write-only binding to the contract
	ITokenReceiverFilterer   // Log filterer for contract events
}

// ITokenReceiverCaller is an auto generated read-only Go binding around a Klaytn contract.
type ITokenReceiverCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ITokenReceiverTransactor is an auto generated write-only Go binding around a Klaytn contract.
type ITokenReceiverTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ITokenReceiverFilterer is an auto generated log filtering Go binding around a Klaytn contract events.
type ITokenReceiverFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ITokenReceiverSession is an auto generated Go binding around a Klaytn contract,
// with pre-set call and transact options.
type ITokenReceiverSession struct {
	Contract     *ITokenReceiver   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ITokenReceiverCallerSession is an auto generated read-only Go binding around a Klaytn contract,
// with pre-set call options.
type ITokenReceiverCallerSession struct {
	Contract *ITokenReceiverCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// ITokenReceiverTransactorSession is an auto generated write-only Go binding around a Klaytn contract,
// with pre-set transact options.
type ITokenReceiverTransactorSession struct {
	Contract     *ITokenReceiverTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// ITokenReceiverRaw is an auto generated low-level Go binding around a Klaytn contract.
type ITokenReceiverRaw struct {
	Contract *ITokenReceiver // Generic contract binding to access the raw methods on
}

// ITokenReceiverCallerRaw is an auto generated low-level read-only Go binding around a Klaytn contract.
type ITokenReceiverCallerRaw struct {
	Contract *ITokenReceiverCaller // Generic read-only contract binding to access the raw methods on
}

// ITokenReceiverTransactorRaw is an auto generated low-level write-only Go binding around a Klaytn contract.
type ITokenReceiverTransactorRaw struct {
	Contract *ITokenReceiverTransactor // Generic write-only contract binding to access the raw methods on
}

// NewITokenReceiver creates a new instance of ITokenReceiver, bound to a specific deployed contract.
func NewITokenReceiver(address common.Address, backend bind.ContractBackend) (*ITokenReceiver, error) {
	contract, err := bindITokenReceiver(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ITokenReceiver{ITokenReceiverCaller: ITokenReceiverCaller{contract: contract}, ITokenReceiverTransactor: ITokenReceiverTransactor{contract: contract}, ITokenReceiverFilterer: ITokenReceiverFilterer{contract: contract}}, nil
}

// NewITokenReceiverCaller creates a new read-only instance of ITokenReceiver, bound to a specific deployed contract.
func NewITokenReceiverCaller(address common.Address, caller bind.ContractCaller) (*ITokenReceiverCaller, error) {
	contract, err := bindITokenReceiver(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ITokenReceiverCaller{contract: contract}, nil
}

// NewITokenReceiverTransactor creates a new write-only instance of ITokenReceiver, bound to a specific deployed contract.
func NewITokenReceiverTransactor(address common.Address, transactor bind.ContractTransactor) (*ITokenReceiverTransactor, error) {
	contract, err := bindITokenReceiver(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ITokenReceiverTransactor{contract: contract}, nil
}

// NewITokenReceiverFilterer creates a new log filterer instance of ITokenReceiver, bound to a specific deployed contract.
func NewITokenReceiverFilterer(address common.Address, filterer bind.ContractFilterer) (*ITokenReceiverFilterer, error) {
	contract, err := bindITokenReceiver(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ITokenReceiverFilterer{contract: contract}, nil
}

// bindITokenReceiver binds a generic wrapper to an already deployed contract.
func bindITokenReceiver(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ITokenReceiverABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ITokenReceiver *ITokenReceiverRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _ITokenReceiver.Contract.ITokenReceiverCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ITokenReceiver *ITokenReceiverRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ITokenReceiver.Contract.ITokenReceiverTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ITokenReceiver *ITokenReceiverRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ITokenReceiver.Contract.ITokenReceiverTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ITokenReceiver *ITokenReceiverCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _ITokenReceiver.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ITokenReceiver *ITokenReceiverTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ITokenReceiver.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ITokenReceiver *ITokenReceiverTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ITokenReceiver.Contract.contract.Transact(opts, method, params...)
}

// OnTokenReceived is a paid mutator transaction binding the contract method 0xf099d9bd.
//
// Solidity: function onTokenReceived(_from address, amount uint256, _to address) returns(bytes4)
func (_ITokenReceiver *ITokenReceiverTransactor) OnTokenReceived(opts *bind.TransactOpts, _from common.Address, amount *big.Int, _to common.Address) (*types.Transaction, error) {
	return _ITokenReceiver.contract.Transact(opts, "onTokenReceived", _from, amount, _to)
}

// OnTokenReceived is a paid mutator transaction binding the contract method 0xf099d9bd.
//
// Solidity: function onTokenReceived(_from address, amount uint256, _to address) returns(bytes4)
func (_ITokenReceiver *ITokenReceiverSession) OnTokenReceived(_from common.Address, amount *big.Int, _to common.Address) (*types.Transaction, error) {
	return _ITokenReceiver.Contract.OnTokenReceived(&_ITokenReceiver.TransactOpts, _from, amount, _to)
}

// OnTokenReceived is a paid mutator transaction binding the contract method 0xf099d9bd.
//
// Solidity: function onTokenReceived(_from address, amount uint256, _to address) returns(bytes4)
func (_ITokenReceiver *ITokenReceiverTransactorSession) OnTokenReceived(_from common.Address, amount *big.Int, _to common.Address) (*types.Transaction, error) {
	return _ITokenReceiver.Contract.OnTokenReceived(&_ITokenReceiver.TransactOpts, _from, amount, _to)
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

// SafeMathABI is the input ABI used to generate the binding from.
const SafeMathABI = "[]"

// SafeMathBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const SafeMathBinRuntime = `0x73000000000000000000000000000000000000000030146080604052600080fd00a165627a7a72305820f757a578764ac5a3ad865cab91db235b14c4e61648a4e8f107bbd2faabbcaa490029`

// SafeMathBin is the compiled bytecode used for deploying new contracts.
const SafeMathBin = `0x604c602c600b82828239805160001a60731460008114601c57601e565bfe5b5030600052607381538281f30073000000000000000000000000000000000000000030146080604052600080fd00a165627a7a72305820f757a578764ac5a3ad865cab91db235b14c4e61648a4e8f107bbd2faabbcaa490029`

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
