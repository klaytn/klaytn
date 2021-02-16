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
// This file is derived from core/types/receipt.go (2018/06/04).
// Modified and improved for the klaytn development.

package types

import (
	"fmt"
	"io"
	"unsafe"

	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/hexutil"
	"github.com/klaytn/klaytn/log"
	"github.com/klaytn/klaytn/rlp"
)

//go:generate gencodec -type Receipt -field-override receiptMarshaling -out gen_receipt_json.go

var (
	receiptStatusFailedRLP     = []byte{}
	receiptStatusSuccessfulRLP = []byte{0x01}
	logger                     = log.NewModuleLogger(log.BlockchainTypes)
)

const (
	// ReceiptStatusFailed is the status code of a transaction if execution failed.
	ReceiptStatusFailed = uint(0)

	// ReceiptStatusSuccessful is the status code of a transaction if execution succeeded.
	ReceiptStatusSuccessful = uint(1)

	// TODO-Klaytn Enable more error below.
	// Klaytn specific
	// NOTE-Klaytn Value should be consecutive from ReceiptStatusFailed to the last ReceiptStatusLast
	//         Add a new ReceiptStatusErrXXX before ReceiptStatusLast
	ReceiptStatusErrDefault                              = uint(0x02) // Default
	ReceiptStatusErrDepth                                = uint(0x03)
	ReceiptStatusErrContractAddressCollision             = uint(0x04)
	ReceiptStatusErrCodeStoreOutOfGas                    = uint(0x05)
	ReceiptStatuserrMaxCodeSizeExceed                    = uint(0x06)
	ReceiptStatusErrOutOfGas                             = uint(0x07)
	ReceiptStatusErrWriteProtection                      = uint(0x08)
	ReceiptStatusErrExecutionReverted                    = uint(0x09)
	ReceiptStatusErrOpcodeComputationCostLimitReached    = uint(0x0a)
	ReceiptStatusErrAddressAlreadyExists                 = uint(0x0b)
	ReceiptStatusErrNotAProgramAccount                   = uint(0x0c)
	ReceiptStatusErrNotHumanReadableAddress              = uint(0x0d)
	ReceiptStatusErrFeeRatioOutOfRange                   = uint(0x0e)
	ReceiptStatusErrAccountKeyFailNotUpdatable           = uint(0x0f)
	ReceiptStatusErrDifferentAccountKeyType              = uint(0x10)
	ReceiptStatusErrAccountKeyNilUninitializable         = uint(0x11)
	ReceiptStatusErrNotOnCurve                           = uint(0x12)
	ReceiptStatusErrZeroKeyWeight                        = uint(0x13)
	ReceiptStatusErrUnserializableKey                    = uint(0x14)
	ReceiptStatusErrDuplicatedKey                        = uint(0x15)
	ReceiptStatusErrWeightedSumOverflow                  = uint(0x16)
	ReceiptStatusErrUnsatisfiableThreshold               = uint(0x17)
	ReceiptStatusErrZeroLength                           = uint(0x18)
	ReceiptStatusErrLengthTooLong                        = uint(0x19)
	ReceiptStatusErrNestedRoleBasedKey                   = uint(0x1a)
	ReceiptStatusErrLegacyTransactionMustBeWithLegacyKey = uint(0x1b)
	ReceiptStatusErrDeprecated                           = uint(0x1c)
	ReceiptStatusErrNotSupported                         = uint(0x1d)
	ReceiptStatusErrInvalidCodeFormat                    = uint(0x1e)
	ReceiptStatusLast                                    = uint(0x1f) // Last value which is not an actual ReceiptStatus
//	ReceiptStatusErrInvalidJumpDestination   // TODO-Klaytn-Issue615
//	ReceiptStatusErrInvalidOpcode            // Default case, because no static message available
//	ReceiptStatusErrStackUnderflow           // Default case, because no static message available
//	ReceiptStatusErrStackOverflow            // Default case, because no static message available
//	ReceiptStatusErrInsufficientBalance      // No receipt available for this error
//	ReceiptStatusErrTotalTimeLimitReached    // No receipt available for this error
//	ReceiptStatusErrGasUintOverflow          // TODO-Klaytn-Issue615

)

// Receipt represents the results of a transaction.
type Receipt struct {
	// Consensus fields
	Status uint   `json:"status"`
	Bloom  Bloom  `json:"logsBloom"         gencodec:"required"`
	Logs   []*Log `json:"logs"              gencodec:"required"`

	// Implementation fields (don't reorder!)
	TxHash          common.Hash    `json:"transactionHash" gencodec:"required"`
	ContractAddress common.Address `json:"contractAddress"`
	GasUsed         uint64         `json:"gasUsed" gencodec:"required"`
}

type receiptMarshaling struct {
	Status  hexutil.Uint
	GasUsed hexutil.Uint64
}

// receiptRLP is the consensus encoding of a receipt.
type receiptRLP struct {
	Status  uint
	GasUsed uint64
	Bloom   Bloom
	Logs    []*Log
}

type receiptStorageRLP struct {
	Status          uint
	Bloom           Bloom
	TxHash          common.Hash
	ContractAddress common.Address
	Logs            []*LogForStorage
	GasUsed         uint64
}

// NewReceipt creates a barebone transaction receipt, copying the init fields.
func NewReceipt(status uint, txHash common.Hash, gasUsed uint64) *Receipt {
	return &Receipt{
		Status:  status,
		TxHash:  txHash,
		GasUsed: gasUsed,
	}
}

// EncodeRLP implements rlp.Encoder, and flattens the consensus fields of a receipt
// into an RLP stream. If no post state is present, byzantium fork is assumed.
func (r *Receipt) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, &receiptRLP{r.Status, r.GasUsed, r.Bloom, r.Logs})
}

// DecodeRLP implements rlp.Decoder, and loads the consensus fields of a receipt
// from an RLP stream.
func (r *Receipt) DecodeRLP(s *rlp.Stream) error {
	var dec receiptRLP
	if err := s.Decode(&dec); err != nil {
		return err
	}
	r.Status = dec.Status
	r.GasUsed, r.Bloom, r.Logs = dec.GasUsed, dec.Bloom, dec.Logs
	return nil
}

// Size returns the approximate memory used by all internal contents. It is used
// to approximate and limit the memory consumption of various caches.
func (r *Receipt) Size() common.StorageSize {
	size := common.StorageSize(unsafe.Sizeof(*r))

	size += common.StorageSize(len(r.Logs)) * common.StorageSize(unsafe.Sizeof(Log{}))
	for _, log := range r.Logs {
		size += common.StorageSize(len(log.Topics)*common.HashLength + len(log.Data))
	}
	return size
}

// String implements the Stringer interface.
func (r *Receipt) String() string {
	return fmt.Sprintf("receipt{status=%d gas=%v bloom=%x logs=%v}", r.Status, r.GasUsed, r.Bloom, r.Logs)
}

// ReceiptForStorage is a wrapper around a Receipt that flattens and parses the
// entire content of a receipt, as opposed to only the consensus fields originally.
type ReceiptForStorage Receipt

// EncodeRLP implements rlp.Encoder, and flattens all content fields of a receipt
// into an RLP stream.
func (r *ReceiptForStorage) EncodeRLP(w io.Writer) error {
	enc := &receiptStorageRLP{
		Status:          r.Status,
		Bloom:           r.Bloom,
		TxHash:          r.TxHash,
		ContractAddress: r.ContractAddress,
		Logs:            make([]*LogForStorage, len(r.Logs)),
		GasUsed:         r.GasUsed,
	}
	for i, log := range r.Logs {
		enc.Logs[i] = (*LogForStorage)(log)
	}
	return rlp.Encode(w, enc)
}

// DecodeRLP implements rlp.Decoder, and loads both consensus and implementation
// fields of a receipt from an RLP stream.
func (r *ReceiptForStorage) DecodeRLP(s *rlp.Stream) error {
	var dec receiptStorageRLP
	if err := s.Decode(&dec); err != nil {
		return err
	}
	r.Status = dec.Status

	// Assign the consensus fields
	r.Bloom = dec.Bloom
	r.Logs = make([]*Log, len(dec.Logs))
	for i, log := range dec.Logs {
		r.Logs[i] = (*Log)(log)
	}
	// Assign the implementation fields
	r.TxHash, r.ContractAddress, r.GasUsed = dec.TxHash, dec.ContractAddress, dec.GasUsed
	return nil
}

// Receipts is a wrapper around a Receipt array to implement DerivableList.
type Receipts []*Receipt

// Len returns the number of receipts in this list.
func (r Receipts) Len() int { return len(r) }

// GetRlp returns the RLP encoding of one receipt from the list.
func (r Receipts) GetRlp(i int) []byte {
	bytes, err := rlp.EncodeToBytes(r[i])
	if err != nil {
		panic(err)
	}
	return bytes
}
