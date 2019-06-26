// Modifications Copyright 2018 The klaytn Authors
// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.
//
// This file is derived from params/protocol_params.go (2018/06/04).
// Modified and improved for the klaytn development.

package params

import (
	"math/big"
	"time"
)

var (
	TargetGasLimit = GenesisGasLimit // The artificial target
)

const (
	// Fee schedule parameters

	CallValueTransferGas  uint64 = 9000  // Paid for CALL when the value transfer is non-zero.                  // G_callvalue
	CallNewAccountGas     uint64 = 25000 // Paid for CALL when the destination address didn't exist prior.      // G_newaccount
	TxGas                 uint64 = 21000 // Per transaction not creating a contract. NOTE: Not payable on data of calls between transactions. // G_transaction
	TxGasContractCreation uint64 = 53000 // Per transaction that creates a contract. NOTE: Not payable on data of calls between transactions. // G_transaction + G_create
	TxDataZeroGas         uint64 = 4     // Per byte of data attached to a transaction that equals zero. NOTE: Not payable on data of calls between transactions. // G_txdatazero
	QuadCoeffDiv          uint64 = 512   // Divisor for the quadratic particle of the memory cost equation.
	SstoreSetGas          uint64 = 20000 // Once per SLOAD operation.                                           // G_sset
	LogDataGas            uint64 = 8     // Per byte in a LOG* operation's data.                                // G_logdata
	CallStipend           uint64 = 2300  // Free gas given at beginning of call.                                // G_callstipend
	Sha3Gas               uint64 = 30    // Once per SHA3 operation.                                                 // G_sha3
	Sha3WordGas           uint64 = 6     // Once per word of the SHA3 operation's data.                              // G_sha3word
	SstoreResetGas        uint64 = 5000  // Once per SSTORE operation if the zeroness changes from zero.             // G_sreset
	SstoreClearGas        uint64 = 5000  // Once per SSTORE operation if the zeroness doesn't change.                // G_sreset
	SstoreRefundGas       uint64 = 15000 // Once per SSTORE operation if the zeroness changes to zero.               // R_sclear
	JumpdestGas           uint64 = 1     // Once per JUMPDEST operation.
	CreateDataGas         uint64 = 200   // Paid per byte for a CREATE operation to succeed in placing code into state. // G_codedeposit
	ExpGas                uint64 = 10    // Once per EXP instruction
	LogGas                uint64 = 375   // Per LOG* operation.                                                          // G_log
	CopyGas               uint64 = 3     // Partial payment for COPY operations, multiplied by words copied, rounded up. // G_copy
	CreateGas             uint64 = 32000 // Once per CREATE operation & contract-creation transaction.               // G_create
	Create2Gas            uint64 = 32000 // Once per CREATE2 operation
	SuicideRefundGas      uint64 = 24000 // Refunded following a suicide operation.                                  // R_selfdestruct
	MemoryGas             uint64 = 3     // Times the address of the (highest referenced byte in memory + 1). NOTE: referencing happens on read, write and in instructions such as RETURN and CALL. // G_memory
	LogTopicGas           uint64 = 375   // Multiplied by the * of the LOG*, per LOG transaction. e.g. LOG0 incurs 0 * c_txLogTopicGas, LOG4 incurs 4 * c_txLogTopicGas.   // G_logtopic
	TxDataNonZeroGas      uint64 = 68    // Per byte of data attached to a transaction that is not equal to zero. NOTE: Not payable on data of calls between transactions. // G_txdatanonzero

	// Fee for Service Chain
	// TODO-Klaytn-ServiceChain The following parameters should be fixed.
	// TODO-Klaytn-Governance The following parameters should be able to be modified by governance.
	TxChainDataAnchoringGas uint64 = 21000 // Per transaction anchoring chain data. NOTE: Not payable on data of calls between transactions. // G_transactionchaindataanchoring
	ChainDataAnchoringGas   uint64 = 100   // Per byte of anchoring chain data NOTE: Not payable on data of calls between transactions. // G_chaindataanchoring

	// Precompiled contract gas prices

	EcrecoverGas            uint64 = 3000   // Elliptic curve sender recovery gas price
	Sha256BaseGas           uint64 = 60     // Base price for a SHA256 operation
	Sha256PerWordGas        uint64 = 12     // Per-word price for a SHA256 operation
	Ripemd160BaseGas        uint64 = 600    // Base price for a RIPEMD160 operation
	Ripemd160PerWordGas     uint64 = 120    // Per-word price for a RIPEMD160 operation
	IdentityBaseGas         uint64 = 15     // Base price for a data copy operation
	IdentityPerWordGas      uint64 = 3      // Per-work price for a data copy operation
	ModExpQuadCoeffDiv      uint64 = 20     // Divisor for the quadratic particle of the big int modular exponentiation
	Bn256AddGas             uint64 = 500    // Gas needed for an elliptic curve addition
	Bn256ScalarMulGas       uint64 = 40000  // Gas needed for an elliptic curve scalar multiplication
	Bn256PairingBaseGas     uint64 = 100000 // Base price for an elliptic curve pairing check
	Bn256PairingPerPointGas uint64 = 80000  // Per-point price for an elliptic curve pairing check
	VMLogBaseGas            uint64 = 100    // Base price for a VMLOG operation
	VMLogPerByteGas         uint64 = 20     // Per-byte price for a VMLOG operation
	FeePayerGas             uint64 = 300    // Gas needed for calculating the fee payer of the transaction in a smart contract.
	ValidateSenderGas       uint64 = 5000   // Gas needed for validating the signature of a message.

	GasLimitBoundDivisor uint64 = 1024    // The bound divisor of the gas limit, used in update calculations.
	MinGasLimit          uint64 = 5000    // Minimum the gas limit may ever be.
	GenesisGasLimit      uint64 = 4712388 // Gas limit of the Genesis block.

	MaximumExtraDataSize uint64 = 32 // Maximum size extra data may be after Genesis.

	EpochDuration   uint64 = 30000 // Duration between proof-of-work epochs.
	CallCreateDepth uint64 = 1024  // Maximum depth of call/create stack.
	StackLimit      uint64 = 1024  // Maximum size of VM stack allowed.

	MaxCodeSize = 24576 // Maximum bytecode to permit for a contract

	// istanbul BFT
	BFTMaximumExtraDataSize uint64 = 65 // Maximum size extra data may be after Genesis.

	// AccountKey
	// TODO-Klaytn: Need to fix below values.
	TxAccountCreationGasDefault uint64 = 0
	TxValidationGasDefault      uint64 = 0
	TxAccountCreationGasPerKey  uint64 = 20000 // WARNING: With integer overflow in mind before changing this value.
	TxValidationGasPerKey       uint64 = 15000 // WARNING: With integer overflow in mind before changing this value.

	// Fee for new tx types
	// TODO-Klaytn: Need to fix values
	TxGasAccountCreation       uint64 = 21000
	TxGasAccountUpdate         uint64 = 21000
	TxGasFeeDelegated          uint64 = 10000
	TxGasFeeDelegatedWithRatio uint64 = 15000
	TxGasCancel                uint64 = 21000

	// Network Id
	UnusedNetworkId  uint64 = 0
	AspenNetworkId   uint64 = 1000
	BaobabNetworkId  uint64 = 1001
	CypressNetworkId uint64 = 8217

	TxGasValueTransfer     uint64 = 21000
	TxGasContractExecution uint64 = 21000

	TxDataGas uint64 = 100
)

var (
	TxGasHumanReadable     uint64 = 4000000000
	BlockScoreBoundDivisor        = big.NewInt(2048)   // The bound divisor of the blockscore, used in the update calculations.
	GenesisBlockScore             = big.NewInt(131072) // BlockScore of the Genesis block.
	MinimumBlockScore             = big.NewInt(131072) // The minimum that the blockscore may ever be.
	DurationLimit                 = big.NewInt(13)     // The decision boundary on the blocktime duration used to determine whether blockscore should go up or not.
)

// Parameters for execution time limit
var (
	// TODO-Klaytn Determine more practical values through actual running experience
	TotalTimeLimit             = 250 * time.Millisecond // Execution time limit for all txs in a block
	OpcodeComputationCostLimit = uint64(100000000)      // Computation cost limit for a tx. For now, it is approximately 100 ms.
)

// istanbul BFT
func GetMaximumExtraDataSize() uint64 {
	return BFTMaximumExtraDataSize
}

type CodeFormat uint8

const (
	CodeFormatEVM CodeFormat = iota
	CodeFormatLast
)

func (t CodeFormat) String() string {
	switch t {
	case CodeFormatEVM:
		return "CodeFormatEVM"
	}

	return "UndefinedCodeFormat"
}

func (t CodeFormat) Validate() bool {
	if t < CodeFormatLast {
		return true
	}
	return false
}
