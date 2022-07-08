// Copyright 2019 The klaytn Authors
// This file is part of the klaytn library.
//
// The klaytn library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The klaytn library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the klaytn library. If not, see <http://www.gnu.org/licenses/>.

package types

import (
	"bytes"
	"strings"
	"testing"

	"github.com/klaytn/klaytn/blockchain/types/accountkey"
	"github.com/klaytn/klaytn/rlp"
	"github.com/stretchr/testify/assert"
)

type testingF func(t *testing.T)

// TestTxRLPDecode tests decoding transactions with encoding their fields converted into []interface{}.
func TestTxRLPDecode(t *testing.T) {
	funcs := []testingF{
		testTxRLPDecodeLegacy,
		testTxRLPDecodeAccessList,
		testTxRLPDecodeDynamicFee,

		testTxRLPDecodeValueTransfer,
		testTxRLPDecodeValueTransferMemo,
		testTxRLPDecodeAccountUpdate,
		testTxRLPDecodeSmartContractDeploy,
		testTxRLPDecodeSmartContractExecution,
		testTxRLPDecodeCancel,
		testTxRLPDecodeChainDataAnchoring,

		testTxRLPDecodeFeeDelegatedValueTransfer,
		testTxRLPDecodeFeeDelegatedValueTransferMemo,
		testTxRLPDecodeFeeDelegatedAccountUpdate,
		testTxRLPDecodeFeeDelegatedSmartContractDeploy,
		testTxRLPDecodeFeeDelegatedSmartContractExecution,
		testTxRLPDecodeFeeDelegatedCancel,
		testTxRLPDecodeFeeDelegatedChainDataAnchoring,

		testTxRLPDecodeFeeDelegatedValueTransferWithRatio,
		testTxRLPDecodeFeeDelegatedValueTransferMemoWithRatio,
		testTxRLPDecodeFeeDelegatedAccountUpdateWithRatio,
		testTxRLPDecodeFeeDelegatedSmartContractDeployWithRatio,
		testTxRLPDecodeFeeDelegatedSmartContractExecutionWithRatio,
		testTxRLPDecodeFeeDelegatedCancelWithRatio,
		testTxRLPDecodeFeeDelegatedChainDataAnchoringWithRatio,
	}

	for _, f := range funcs {
		fnname := getFunctionName(f)
		fnname = fnname[strings.LastIndex(fnname, ".")+1:]
		t.Run(fnname, func(t *testing.T) {
			f(t)
		})
	}
}

func testTxRLPDecodeLegacy(t *testing.T) {
	tx := genLegacyTransaction().(*TxInternalDataLegacy)

	buffer := new(bytes.Buffer)
	err := rlp.Encode(buffer, tx.Type())
	assert.Equal(t, nil, err)

	err = rlp.Encode(buffer, []interface{}{
		tx.AccountNonce,
		tx.Price, // in case big.Int, it does not allow nil.
		tx.GasLimit,
		tx.Recipient,
		tx.Amount,
		tx.Payload, // in case of []bytes, it does not allow nil.
		tx.V,
		tx.R,
		tx.S,
	})
	assert.Equal(t, nil, err)

	dec := newTxInternalDataSerializer()

	if err := rlp.DecodeBytes(buffer.Bytes(), &dec); err != nil {
		panic(err)
	}

	if !tx.Equal(dec.tx) {
		t.Fatalf("tx != dec.tx\ntx=%v\ndec.tx=%v", tx, dec.tx)
	}
}

func testTxRLPDecodeValueTransfer(t *testing.T) {
	tx := genValueTransferTransaction().(*TxInternalDataValueTransfer)

	buffer := new(bytes.Buffer)
	err := rlp.Encode(buffer, tx.Type())
	assert.Equal(t, nil, err)

	err = rlp.Encode(buffer, []interface{}{
		tx.AccountNonce,
		tx.Price,
		tx.GasLimit,
		tx.Recipient,
		tx.Amount,
		tx.From,
		tx.TxSignatures,
	})
	assert.Equal(t, nil, err)

	dec := newTxInternalDataSerializer()

	if err := rlp.DecodeBytes(buffer.Bytes(), &dec); err != nil {
		panic(err)
	}

	if !tx.Equal(dec.tx) {
		t.Fatalf("tx != dec.tx\ntx=%v\ndec.tx=%v", tx, dec.tx)
	}
}

func testTxRLPDecodeValueTransferMemo(t *testing.T) {
	tx := genValueTransferMemoTransaction().(*TxInternalDataValueTransferMemo)

	buffer := new(bytes.Buffer)
	err := rlp.Encode(buffer, tx.Type())
	assert.Equal(t, nil, err)

	err = rlp.Encode(buffer, []interface{}{
		tx.AccountNonce,
		tx.Price,
		tx.GasLimit,
		tx.Recipient,
		tx.Amount,
		tx.From,
		tx.Payload,
		tx.TxSignatures,
	})
	assert.Equal(t, nil, err)

	dec := newTxInternalDataSerializer()

	if err := rlp.DecodeBytes(buffer.Bytes(), &dec); err != nil {
		panic(err)
	}

	if !tx.Equal(dec.tx) {
		t.Fatalf("tx != dec.tx\ntx=%v\ndec.tx=%v", tx, dec.tx)
	}
}

//func testTxRLPDecodeAccountCreation(t *testing.T) {
//	tx := genAccountCreationTransaction().(*TxInternalDataAccountCreation)
//
//	buffer := new(bytes.Buffer)
//	err := rlp.Encode(buffer, tx.Type())
//	assert.Equal(t, nil, err)
//
//	encodedKey, err := rlp.EncodeToBytes(accountkey.NewAccountKeySerializerWithAccountKey(tx.Key))
//	assert.Equal(t, nil, err)
//
//	err = rlp.Encode(buffer, []interface{}{
//		tx.AccountNonce,
//		tx.Price,
//		tx.GasLimit,
//		tx.Recipient,
//		tx.Amount,
//		tx.From,
//		tx.HumanReadable, // bool only allows 0 or 1.
//		encodedKey,
//		tx.TxSignatures,
//	})
//	assert.Equal(t, nil, err)
//
//	dec := newTxInternalDataSerializer()
//
//	if err := rlp.DecodeBytes(buffer.Bytes(), &dec); err != nil {
//		panic(err)
//	}
//
//	if !tx.Equal(dec.tx) {
//		t.Fatalf("tx != dec.tx\ntx=%v\ndec.tx=%v", tx, dec.tx)
//	}
//}

func testTxRLPDecodeAccountUpdate(t *testing.T) {
	tx := genAccountUpdateTransaction().(*TxInternalDataAccountUpdate)

	buffer := new(bytes.Buffer)
	err := rlp.Encode(buffer, tx.Type())
	assert.Equal(t, nil, err)

	encodedKey, err := rlp.EncodeToBytes(accountkey.NewAccountKeySerializerWithAccountKey(tx.Key))
	assert.Equal(t, nil, err)

	err = rlp.Encode(buffer, []interface{}{
		tx.AccountNonce,
		tx.Price,
		tx.GasLimit,
		tx.From,
		encodedKey,
		tx.TxSignatures,
	})
	assert.Equal(t, nil, err)

	dec := newTxInternalDataSerializer()

	if err := rlp.DecodeBytes(buffer.Bytes(), &dec); err != nil {
		panic(err)
	}

	if !tx.Equal(dec.tx) {
		t.Fatalf("tx != dec.tx\ntx=%v\ndec.tx=%v", tx, dec.tx)
	}
}

func testTxRLPDecodeSmartContractDeploy(t *testing.T) {
	tx := genSmartContractDeployTransaction().(*TxInternalDataSmartContractDeploy)

	buffer := new(bytes.Buffer)
	err := rlp.Encode(buffer, tx.Type())
	assert.Equal(t, nil, err)

	err = rlp.Encode(buffer, []interface{}{
		tx.AccountNonce,
		tx.Price,
		tx.GasLimit,
		tx.Recipient,
		tx.Amount,
		tx.From,
		tx.Payload,
		tx.HumanReadable, // bool only allows 0 or 1.
		tx.CodeFormat,
		tx.TxSignatures,
	})
	assert.Equal(t, nil, err)

	dec := newTxInternalDataSerializer()

	if err := rlp.DecodeBytes(buffer.Bytes(), &dec); err != nil {
		panic(err)
	}

	if !tx.Equal(dec.tx) {
		t.Fatalf("tx != dec.tx\ntx=%v\ndec.tx=%v", tx, dec.tx)
	}
}

func testTxRLPDecodeSmartContractExecution(t *testing.T) {
	tx := genSmartContractExecutionTransaction().(*TxInternalDataSmartContractExecution)

	buffer := new(bytes.Buffer)
	err := rlp.Encode(buffer, tx.Type())
	assert.Equal(t, nil, err)

	err = rlp.Encode(buffer, []interface{}{
		tx.AccountNonce,
		tx.Price,
		tx.GasLimit,
		tx.Recipient,
		tx.Amount,
		tx.From,
		tx.Payload,
		tx.TxSignatures,
	})
	assert.Equal(t, nil, err)

	dec := newTxInternalDataSerializer()

	if err := rlp.DecodeBytes(buffer.Bytes(), &dec); err != nil {
		panic(err)
	}

	if !tx.Equal(dec.tx) {
		t.Fatalf("tx != dec.tx\ntx=%v\ndec.tx=%v", tx, dec.tx)
	}
}

func testTxRLPDecodeCancel(t *testing.T) {
	tx := genCancelTransaction().(*TxInternalDataCancel)

	buffer := new(bytes.Buffer)
	err := rlp.Encode(buffer, tx.Type())
	assert.Equal(t, nil, err)

	err = rlp.Encode(buffer, []interface{}{
		tx.AccountNonce,
		tx.Price,
		tx.GasLimit,
		tx.From,
		tx.TxSignatures,
	})
	assert.Equal(t, nil, err)

	dec := newTxInternalDataSerializer()

	if err := rlp.DecodeBytes(buffer.Bytes(), &dec); err != nil {
		panic(err)
	}

	if !tx.Equal(dec.tx) {
		t.Fatalf("tx != dec.tx\ntx=%v\ndec.tx=%v", tx, dec.tx)
	}
}

func testTxRLPDecodeChainDataAnchoring(t *testing.T) {
	tx := genChainDataTransaction().(*TxInternalDataChainDataAnchoring)

	buffer := new(bytes.Buffer)
	err := rlp.Encode(buffer, tx.Type())
	assert.Equal(t, nil, err)

	err = rlp.Encode(buffer, []interface{}{
		tx.AccountNonce,
		tx.Price,
		tx.GasLimit,
		tx.From,
		tx.Payload,
		tx.TxSignatures,
	})
	assert.Equal(t, nil, err)

	dec := newTxInternalDataSerializer()

	if err := rlp.DecodeBytes(buffer.Bytes(), &dec); err != nil {
		panic(err)
	}

	if !tx.Equal(dec.tx) {
		t.Fatalf("tx != dec.tx\ntx=%v\ndec.tx=%v", tx, dec.tx)
	}
}

func testTxRLPDecodeFeeDelegatedValueTransfer(t *testing.T) {
	tx := genFeeDelegatedValueTransferTransaction().(*TxInternalDataFeeDelegatedValueTransfer)

	buffer := new(bytes.Buffer)
	err := rlp.Encode(buffer, tx.Type())
	assert.Equal(t, nil, err)

	err = rlp.Encode(buffer, []interface{}{
		tx.AccountNonce,
		tx.Price,
		tx.GasLimit,
		tx.Recipient,
		tx.Amount,
		tx.From,
		tx.TxSignatures,
		tx.FeePayer,
		tx.FeePayerSignatures,
	})
	assert.Equal(t, nil, err)

	dec := newTxInternalDataSerializer()

	if err := rlp.DecodeBytes(buffer.Bytes(), &dec); err != nil {
		panic(err)
	}

	if !tx.Equal(dec.tx) {
		t.Fatalf("tx != dec.tx\ntx=%v\ndec.tx=%v", tx, dec.tx)
	}
}

func testTxRLPDecodeFeeDelegatedValueTransferMemo(t *testing.T) {
	tx := genFeeDelegatedValueTransferMemoTransaction().(*TxInternalDataFeeDelegatedValueTransferMemo)

	buffer := new(bytes.Buffer)
	err := rlp.Encode(buffer, tx.Type())
	assert.Equal(t, nil, err)

	err = rlp.Encode(buffer, []interface{}{
		tx.AccountNonce,
		tx.Price,
		tx.GasLimit,
		tx.Recipient,
		tx.Amount,
		tx.From,
		tx.Payload,
		tx.TxSignatures,
		tx.FeePayer,
		tx.FeePayerSignatures,
	})
	assert.Equal(t, nil, err)

	dec := newTxInternalDataSerializer()

	if err := rlp.DecodeBytes(buffer.Bytes(), &dec); err != nil {
		panic(err)
	}

	if !tx.Equal(dec.tx) {
		t.Fatalf("tx != dec.tx\ntx=%v\ndec.tx=%v", tx, dec.tx)
	}
}

func testTxRLPDecodeFeeDelegatedAccountUpdate(t *testing.T) {
	tx := genFeeDelegatedAccountUpdateTransaction().(*TxInternalDataFeeDelegatedAccountUpdate)

	buffer := new(bytes.Buffer)
	err := rlp.Encode(buffer, tx.Type())
	assert.Equal(t, nil, err)

	encodedKey, err := rlp.EncodeToBytes(accountkey.NewAccountKeySerializerWithAccountKey(tx.Key))
	assert.Equal(t, nil, err)

	err = rlp.Encode(buffer, []interface{}{
		tx.AccountNonce,
		tx.Price,
		tx.GasLimit,
		tx.From,
		encodedKey,
		tx.TxSignatures,
		tx.FeePayer,
		tx.FeePayerSignatures,
	})
	assert.Equal(t, nil, err)

	dec := newTxInternalDataSerializer()

	if err := rlp.DecodeBytes(buffer.Bytes(), &dec); err != nil {
		panic(err)
	}

	if !tx.Equal(dec.tx) {
		t.Fatalf("tx != dec.tx\ntx=%v\ndec.tx=%v", tx, dec.tx)
	}
}

func testTxRLPDecodeFeeDelegatedSmartContractDeploy(t *testing.T) {
	tx := genFeeDelegatedSmartContractDeployTransaction().(*TxInternalDataFeeDelegatedSmartContractDeploy)

	buffer := new(bytes.Buffer)
	err := rlp.Encode(buffer, tx.Type())
	assert.Equal(t, nil, err)

	err = rlp.Encode(buffer, []interface{}{
		tx.AccountNonce,
		tx.Price,
		tx.GasLimit,
		tx.Recipient,
		tx.Amount,
		tx.From,
		tx.Payload,
		tx.HumanReadable,
		tx.CodeFormat,
		tx.TxSignatures,
		tx.FeePayer,
		tx.FeePayerSignatures,
	})
	assert.Equal(t, nil, err)

	dec := newTxInternalDataSerializer()

	if err := rlp.DecodeBytes(buffer.Bytes(), &dec); err != nil {
		panic(err)
	}

	if !tx.Equal(dec.tx) {
		t.Fatalf("tx != dec.tx\ntx=%v\ndec.tx=%v", tx, dec.tx)
	}
}

func testTxRLPDecodeFeeDelegatedSmartContractExecution(t *testing.T) {
	tx := genFeeDelegatedSmartContractExecutionTransaction().(*TxInternalDataFeeDelegatedSmartContractExecution)

	buffer := new(bytes.Buffer)
	err := rlp.Encode(buffer, tx.Type())
	assert.Equal(t, nil, err)

	err = rlp.Encode(buffer, []interface{}{
		tx.AccountNonce,
		tx.Price,
		tx.GasLimit,
		tx.Recipient,
		tx.Amount,
		tx.From,
		tx.Payload,
		tx.TxSignatures,
		tx.FeePayer,
		tx.FeePayerSignatures,
	})
	assert.Equal(t, nil, err)

	dec := newTxInternalDataSerializer()

	if err := rlp.DecodeBytes(buffer.Bytes(), &dec); err != nil {
		panic(err)
	}

	if !tx.Equal(dec.tx) {
		t.Fatalf("tx != dec.tx\ntx=%v\ndec.tx=%v", tx, dec.tx)
	}
}

func testTxRLPDecodeFeeDelegatedCancel(t *testing.T) {
	tx := genFeeDelegatedCancelTransaction().(*TxInternalDataFeeDelegatedCancel)

	buffer := new(bytes.Buffer)
	err := rlp.Encode(buffer, tx.Type())
	assert.Equal(t, nil, err)

	err = rlp.Encode(buffer, []interface{}{
		tx.AccountNonce,
		tx.Price,
		tx.GasLimit,
		tx.From,
		tx.TxSignatures,
		tx.FeePayer,
		tx.FeePayerSignatures,
	})
	assert.Equal(t, nil, err)

	dec := newTxInternalDataSerializer()

	if err := rlp.DecodeBytes(buffer.Bytes(), &dec); err != nil {
		panic(err)
	}

	if !tx.Equal(dec.tx) {
		t.Fatalf("tx != dec.tx\ntx=%v\ndec.tx=%v", tx, dec.tx)
	}
}

func testTxRLPDecodeFeeDelegatedChainDataAnchoring(t *testing.T) {
	tx := genFeeDelegatedChainDataTransaction().(*TxInternalDataFeeDelegatedChainDataAnchoring)

	buffer := new(bytes.Buffer)
	err := rlp.Encode(buffer, tx.Type())
	assert.Equal(t, nil, err)

	err = rlp.Encode(buffer, []interface{}{
		tx.AccountNonce,
		tx.Price,
		tx.GasLimit,
		tx.From,
		tx.Payload,
		tx.TxSignatures,
		tx.FeePayer,
		tx.FeePayerSignatures,
	})
	assert.Equal(t, nil, err)

	dec := newTxInternalDataSerializer()

	if err := rlp.DecodeBytes(buffer.Bytes(), &dec); err != nil {
		panic(err)
	}

	if !tx.Equal(dec.tx) {
		t.Fatalf("tx != dec.tx\ntx=%v\ndec.tx=%v", tx, dec.tx)
	}
}

func testTxRLPDecodeFeeDelegatedValueTransferWithRatio(t *testing.T) {
	tx := genFeeDelegatedValueTransferWithRatioTransaction().(*TxInternalDataFeeDelegatedValueTransferWithRatio)

	buffer := new(bytes.Buffer)
	err := rlp.Encode(buffer, tx.Type())
	assert.Equal(t, nil, err)

	err = rlp.Encode(buffer, []interface{}{
		tx.AccountNonce,
		tx.Price,
		tx.GasLimit,
		tx.Recipient,
		tx.Amount,
		tx.From,
		tx.FeeRatio,
		tx.TxSignatures,
		tx.FeePayer,
		tx.FeePayerSignatures,
	})
	assert.Equal(t, nil, err)

	dec := newTxInternalDataSerializer()

	if err := rlp.DecodeBytes(buffer.Bytes(), &dec); err != nil {
		panic(err)
	}

	if !tx.Equal(dec.tx) {
		t.Fatalf("tx != dec.tx\ntx=%v\ndec.tx=%v", tx, dec.tx)
	}
}

func testTxRLPDecodeFeeDelegatedValueTransferMemoWithRatio(t *testing.T) {
	tx := genFeeDelegatedValueTransferMemoWithRatioTransaction().(*TxInternalDataFeeDelegatedValueTransferMemoWithRatio)

	buffer := new(bytes.Buffer)
	err := rlp.Encode(buffer, tx.Type())
	assert.Equal(t, nil, err)

	err = rlp.Encode(buffer, []interface{}{
		tx.AccountNonce,
		tx.Price,
		tx.GasLimit,
		tx.Recipient,
		tx.Amount,
		tx.From,
		tx.Payload,
		tx.FeeRatio,
		tx.TxSignatures,
		tx.FeePayer,
		tx.FeePayerSignatures,
	})
	assert.Equal(t, nil, err)

	dec := newTxInternalDataSerializer()

	if err := rlp.DecodeBytes(buffer.Bytes(), &dec); err != nil {
		panic(err)
	}

	if !tx.Equal(dec.tx) {
		t.Fatalf("tx != dec.tx\ntx=%v\ndec.tx=%v", tx, dec.tx)
	}
}

func testTxRLPDecodeFeeDelegatedAccountUpdateWithRatio(t *testing.T) {
	tx := genFeeDelegatedAccountUpdateWithRatioTransaction().(*TxInternalDataFeeDelegatedAccountUpdateWithRatio)

	buffer := new(bytes.Buffer)
	err := rlp.Encode(buffer, tx.Type())
	assert.Equal(t, nil, err)

	encodedKey, err := rlp.EncodeToBytes(accountkey.NewAccountKeySerializerWithAccountKey(tx.Key))
	assert.Equal(t, nil, err)

	err = rlp.Encode(buffer, []interface{}{
		tx.AccountNonce,
		tx.Price,
		tx.GasLimit,
		tx.From,
		encodedKey,
		tx.FeeRatio,
		tx.TxSignatures,
		tx.FeePayer,
		tx.FeePayerSignatures,
	})
	assert.Equal(t, nil, err)

	dec := newTxInternalDataSerializer()

	if err := rlp.DecodeBytes(buffer.Bytes(), &dec); err != nil {
		panic(err)
	}

	if !tx.Equal(dec.tx) {
		t.Fatalf("tx != dec.tx\ntx=%v\ndec.tx=%v", tx, dec.tx)
	}
}

func testTxRLPDecodeFeeDelegatedSmartContractDeployWithRatio(t *testing.T) {
	tx := genFeeDelegatedSmartContractDeployWithRatioTransaction().(*TxInternalDataFeeDelegatedSmartContractDeployWithRatio)

	buffer := new(bytes.Buffer)
	err := rlp.Encode(buffer, tx.Type())
	assert.Equal(t, nil, err)

	err = rlp.Encode(buffer, []interface{}{
		tx.AccountNonce,
		tx.Price,
		tx.GasLimit,
		tx.Recipient,
		tx.Amount,
		tx.From,
		tx.Payload,
		tx.HumanReadable,
		tx.FeeRatio,
		tx.CodeFormat,
		tx.TxSignatures,
		tx.FeePayer,
		tx.FeePayerSignatures,
	})
	assert.Equal(t, nil, err)

	dec := newTxInternalDataSerializer()

	if err := rlp.DecodeBytes(buffer.Bytes(), &dec); err != nil {
		panic(err)
	}

	if !tx.Equal(dec.tx) {
		t.Fatalf("tx != dec.tx\ntx=%v\ndec.tx=%v", tx, dec.tx)
	}
}

func testTxRLPDecodeFeeDelegatedSmartContractExecutionWithRatio(t *testing.T) {
	tx := genFeeDelegatedSmartContractExecutionWithRatioTransaction().(*TxInternalDataFeeDelegatedSmartContractExecutionWithRatio)

	buffer := new(bytes.Buffer)
	err := rlp.Encode(buffer, tx.Type())
	assert.Equal(t, nil, err)

	err = rlp.Encode(buffer, []interface{}{
		tx.AccountNonce,
		tx.Price,
		tx.GasLimit,
		tx.Recipient,
		tx.Amount,
		tx.From,
		tx.Payload,
		tx.FeeRatio,
		tx.TxSignatures,
		tx.FeePayer,
		tx.FeePayerSignatures,
	})
	assert.Equal(t, nil, err)

	dec := newTxInternalDataSerializer()

	if err := rlp.DecodeBytes(buffer.Bytes(), &dec); err != nil {
		panic(err)
	}

	if !tx.Equal(dec.tx) {
		t.Fatalf("tx != dec.tx\ntx=%v\ndec.tx=%v", tx, dec.tx)
	}
}

func testTxRLPDecodeFeeDelegatedCancelWithRatio(t *testing.T) {
	tx := genFeeDelegatedCancelWithRatioTransaction().(*TxInternalDataFeeDelegatedCancelWithRatio)

	buffer := new(bytes.Buffer)
	err := rlp.Encode(buffer, tx.Type())
	assert.Equal(t, nil, err)

	err = rlp.Encode(buffer, []interface{}{
		tx.AccountNonce,
		tx.Price,
		tx.GasLimit,
		tx.From,
		tx.FeeRatio,
		tx.TxSignatures,
		tx.FeePayer,
		tx.FeePayerSignatures,
	})
	assert.Equal(t, nil, err)

	dec := newTxInternalDataSerializer()

	if err := rlp.DecodeBytes(buffer.Bytes(), &dec); err != nil {
		panic(err)
	}

	if !tx.Equal(dec.tx) {
		t.Fatalf("tx != dec.tx\ntx=%v\ndec.tx=%v", tx, dec.tx)
	}
}

func testTxRLPDecodeFeeDelegatedChainDataAnchoringWithRatio(t *testing.T) {
	tx := genFeeDelegatedChainDataWithRatioTransaction().(*TxInternalDataFeeDelegatedChainDataAnchoringWithRatio)

	buffer := new(bytes.Buffer)
	err := rlp.Encode(buffer, tx.Type())
	assert.Equal(t, nil, err)

	err = rlp.Encode(buffer, []interface{}{
		tx.AccountNonce,
		tx.Price,
		tx.GasLimit,
		tx.From,
		tx.Payload,
		tx.FeeRatio,
		tx.TxSignatures,
		tx.FeePayer,
		tx.FeePayerSignatures,
	})
	assert.Equal(t, nil, err)

	dec := newTxInternalDataSerializer()

	if err := rlp.DecodeBytes(buffer.Bytes(), &dec); err != nil {
		panic(err)
	}

	if !tx.Equal(dec.tx) {
		t.Fatalf("tx != dec.tx\ntx=%v\ndec.tx=%v", tx, dec.tx)
	}
}

func testTxRLPDecodeAccessList(t *testing.T) {
	tx := genAccessListTransaction().(*TxInternalDataEthereumAccessList)

	buffer := new(bytes.Buffer)
	err := rlp.Encode(buffer, tx.Type())
	assert.Equal(t, nil, err)

	err = rlp.Encode(buffer, []interface{}{
		tx.ChainID,
		tx.AccountNonce,
		tx.Price,
		tx.GasLimit,
		tx.Recipient,
		tx.Amount,
		tx.Payload,
		tx.AccessList,
		tx.V,
		tx.R,
		tx.S,
	})

	assert.Equal(t, nil, err)

	dec := newTxInternalDataSerializer()

	if err := rlp.DecodeBytes(buffer.Bytes(), &dec); err != nil {
		panic(err)
	}

	if !tx.Equal(dec.tx) {
		t.Fatalf("tx != dec.tx\ntx=%v\ndec.tx=%v", tx, dec.tx)
	}
}

func testTxRLPDecodeDynamicFee(t *testing.T) {
	tx := genDynamicFeeTransaction().(*TxInternalDataEthereumDynamicFee)

	buffer := new(bytes.Buffer)
	err := rlp.Encode(buffer, tx.Type())
	assert.Equal(t, nil, err)

	err = rlp.Encode(buffer, []interface{}{
		tx.ChainID,
		tx.AccountNonce,
		tx.GasTipCap,
		tx.GasFeeCap,
		tx.GasLimit,
		tx.Recipient,
		tx.Amount,
		tx.Payload,
		tx.AccessList,
		tx.V,
		tx.R,
		tx.S,
	})

	assert.Equal(t, nil, err)

	dec := newTxInternalDataSerializer()

	if err := rlp.DecodeBytes(buffer.Bytes(), &dec); err != nil {
		panic(err)
	}

	if !tx.Equal(dec.tx) {
		t.Fatalf("tx != dec.tx\ntx=%v\ndec.tx=%v", tx, dec.tx)
	}
}
