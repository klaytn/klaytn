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
	"math/big"
	"testing"

	"github.com/klaytn/klaytn/blockchain/types/accountkey"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/crypto/sha3"
	"github.com/klaytn/klaytn/params"
	"github.com/klaytn/klaytn/rlp"
	"github.com/stretchr/testify/assert"
)

// TestTransactionSenderTxHash tests SenderTxHash() of all transaction types.
func TestTransactionSenderTxHash(t *testing.T) {
	var txs = []struct {
		Name string
		tx   TxInternalData
	}{
		{"OriginalTx", genLegacyTransaction()},
		{"SmartContractDeploy", genSmartContractDeployTransaction()},
		{"FeeDelegatedSmartContractDeploy", genFeeDelegatedSmartContractDeployTransaction()},
		{"FeeDelegatedSmartContractDeployWithRatio", genFeeDelegatedSmartContractDeployWithRatioTransaction()},
		{"ValueTransfer", genValueTransferTransaction()},
		{"ValueTransferMemo", genValueTransferMemoTransaction()},
		{"FeeDelegatedValueTransferMemo", genFeeDelegatedValueTransferMemoTransaction()},
		{"FeeDelegatedValueTransferMemoWithRatio", genFeeDelegatedValueTransferMemoWithRatioTransaction()},
		{"ChainDataAnchoring", genChainDataTransaction()},
		{"FeeDelegatedChainDataAnchoring", genFeeDelegatedChainDataTransaction()},
		{"FeeDelegatedChainDataAnchoringWithRatio", genFeeDelegatedChainDataWithRatioTransaction()},
		{"AccountUpdate", genAccountUpdateTransaction()},
		{"FeeDelegatedAccountUpdate", genFeeDelegatedAccountUpdateTransaction()},
		{"FeeDelegatedAccountUpdateWithRatio", genFeeDelegatedAccountUpdateWithRatioTransaction()},
		{"FeeDelegatedValueTransfer", genFeeDelegatedValueTransferTransaction()},
		{"SmartContractExecution", genSmartContractExecutionTransaction()},
		{"FeeDelegatedSmartContractExecution", genFeeDelegatedSmartContractExecutionTransaction()},
		{"FeeDelegatedSmartContractExecutionWithRatio", genFeeDelegatedSmartContractExecutionWithRatioTransaction()},
		{"FeeDelegatedValueTransferWithRatio", genFeeDelegatedValueTransferWithRatioTransaction()},
		{"Cancel", genCancelTransaction()},
		{"FeeDelegatedCancel", genFeeDelegatedCancelTransaction()},
		{"FeeDelegatedCancelWithRatio", genFeeDelegatedCancelWithRatioTransaction()},
		{"AccessList", genAccessListTransaction()},
		{"DynamicFee", genDynamicFeeTransaction()},
	}

	var testcases = []struct {
		Name string
		fn   func(t *testing.T, tx TxInternalData)
	}{
		{"SenderTxHash", testTransactionSenderTxHash},
	}

	txMap := make(map[TxType]TxInternalData)
	for _, test := range testcases {
		for _, tx := range txs {
			txMap[tx.tx.Type()] = tx.tx
			Name := test.Name + "/" + tx.Name
			t.Run(Name, func(t *testing.T) {
				test.fn(t, tx.tx)
			})
		}
	}

	// Below code checks whether serialization for all tx implementations is done or not.
	// If no serialization, make test failed.
	for i := TxTypeLegacyTransaction; i < TxTypeEthereumLast; i++ {
		if i == TxTypeKlaytnLast {
			i = TxTypeEthereumAccessList
		}

		// TxTypeAccountCreation is not supported now
		if i == TxTypeAccountCreation {
			continue
		}
		tx, err := NewTxInternalData(i)
		if err == nil {
			if _, ok := txMap[tx.Type()]; !ok {
				t.Errorf("No serialization test for tx %s", tx.Type().String())
			}
		}
	}
}

func testTransactionSenderTxHash(t *testing.T, tx TxInternalData) {
	signer := MakeSigner(params.BFTTestChainConfig, big.NewInt(2))
	rawTx := &Transaction{data: tx}
	rawTx.Sign(signer, key)

	if _, ok := tx.(TxInternalDataFeePayer); ok {
		rawTx.SignFeePayer(signer, key)
	}

	switch v := rawTx.data.(type) {
	case *TxInternalDataLegacy:
		senderTxHash := rawTx.GetTxInternalData().SenderTxHash()
		assert.Equal(t, rawTx.Hash(), senderTxHash)

	case *TxInternalDataValueTransfer:
		senderTxHash := rawTx.GetTxInternalData().SenderTxHash()
		assert.Equal(t, rawTx.Hash(), senderTxHash)

	case *TxInternalDataFeeDelegatedValueTransfer:
		hw := sha3.NewKeccak256()
		rlp.Encode(hw, rawTx.Type())
		rlp.Encode(hw, []interface{}{
			v.AccountNonce,
			v.Price,
			v.GasLimit,
			v.Recipient,
			v.Amount,
			v.From,
			v.TxSignatures,
		})

		h := common.Hash{}

		hw.Sum(h[:0])
		senderTxHash := rawTx.GetTxInternalData().SenderTxHash()
		assert.Equal(t, h, senderTxHash)

	case *TxInternalDataFeeDelegatedValueTransferWithRatio:
		hw := sha3.NewKeccak256()
		rlp.Encode(hw, rawTx.Type())
		rlp.Encode(hw, []interface{}{
			v.AccountNonce,
			v.Price,
			v.GasLimit,
			v.Recipient,
			v.Amount,
			v.From,
			v.FeeRatio,
			v.TxSignatures,
		})

		h := common.Hash{}

		hw.Sum(h[:0])
		senderTxHash := rawTx.GetTxInternalData().SenderTxHash()
		assert.Equal(t, h, senderTxHash)

	case *TxInternalDataValueTransferMemo:
		senderTxHash := rawTx.GetTxInternalData().SenderTxHash()
		assert.Equal(t, rawTx.Hash(), senderTxHash)

	case *TxInternalDataFeeDelegatedValueTransferMemo:
		hw := sha3.NewKeccak256()
		rlp.Encode(hw, rawTx.Type())
		rlp.Encode(hw, []interface{}{
			v.AccountNonce,
			v.Price,
			v.GasLimit,
			v.Recipient,
			v.Amount,
			v.From,
			v.Payload,
			v.TxSignatures,
		})

		h := common.Hash{}

		hw.Sum(h[:0])
		senderTxHash := rawTx.GetTxInternalData().SenderTxHash()
		assert.Equal(t, h, senderTxHash)

	case *TxInternalDataFeeDelegatedValueTransferMemoWithRatio:
		hw := sha3.NewKeccak256()
		rlp.Encode(hw, rawTx.Type())
		rlp.Encode(hw, []interface{}{
			v.AccountNonce,
			v.Price,
			v.GasLimit,
			v.Recipient,
			v.Amount,
			v.From,
			v.Payload,
			v.FeeRatio,
			v.TxSignatures,
		})

		h := common.Hash{}

		hw.Sum(h[:0])
		senderTxHash := rawTx.GetTxInternalData().SenderTxHash()
		assert.Equal(t, h, senderTxHash)

	//case *TxInternalDataAccountCreation:
	//	senderTxHash := rawTx.GetTxInternalData().SenderTxHash()
	//	assert.Equal(t, rawTx.Hash(), senderTxHash)

	case *TxInternalDataAccountUpdate:
		senderTxHash := rawTx.GetTxInternalData().SenderTxHash()
		assert.Equal(t, rawTx.Hash(), senderTxHash)

	case *TxInternalDataFeeDelegatedAccountUpdate:
		serializer := accountkey.NewAccountKeySerializerWithAccountKey(v.Key)
		keyEnc, _ := rlp.EncodeToBytes(serializer)
		hw := sha3.NewKeccak256()
		rlp.Encode(hw, rawTx.Type())
		rlp.Encode(hw, []interface{}{
			v.AccountNonce,
			v.Price,
			v.GasLimit,
			v.From,
			keyEnc,
			v.TxSignatures,
		})

		h := common.Hash{}

		hw.Sum(h[:0])
		senderTxHash := rawTx.GetTxInternalData().SenderTxHash()
		assert.Equal(t, h, senderTxHash)

	case *TxInternalDataFeeDelegatedAccountUpdateWithRatio:
		serializer := accountkey.NewAccountKeySerializerWithAccountKey(v.Key)
		keyEnc, _ := rlp.EncodeToBytes(serializer)
		hw := sha3.NewKeccak256()
		rlp.Encode(hw, rawTx.Type())
		rlp.Encode(hw, []interface{}{
			v.AccountNonce,
			v.Price,
			v.GasLimit,
			v.From,
			keyEnc,
			v.FeeRatio,
			v.TxSignatures,
		})

		h := common.Hash{}

		hw.Sum(h[:0])
		senderTxHash := rawTx.GetTxInternalData().SenderTxHash()
		assert.Equal(t, h, senderTxHash)

	case *TxInternalDataSmartContractDeploy:
		senderTxHash := rawTx.GetTxInternalData().SenderTxHash()
		assert.Equal(t, rawTx.Hash(), senderTxHash)

	case *TxInternalDataFeeDelegatedSmartContractDeploy:
		hw := sha3.NewKeccak256()
		rlp.Encode(hw, rawTx.Type())
		rlp.Encode(hw, []interface{}{
			v.AccountNonce,
			v.Price,
			v.GasLimit,
			v.Recipient,
			v.Amount,
			v.From,
			v.Payload,
			v.HumanReadable,
			v.CodeFormat,
			v.TxSignatures,
		})

		h := common.Hash{}

		hw.Sum(h[:0])
		senderTxHash := rawTx.GetTxInternalData().SenderTxHash()
		assert.Equal(t, h, senderTxHash)

	case *TxInternalDataFeeDelegatedSmartContractDeployWithRatio:
		hw := sha3.NewKeccak256()
		rlp.Encode(hw, rawTx.Type())
		rlp.Encode(hw, []interface{}{
			v.AccountNonce,
			v.Price,
			v.GasLimit,
			v.Recipient,
			v.Amount,
			v.From,
			v.Payload,
			v.HumanReadable,
			v.FeeRatio,
			v.CodeFormat,
			v.TxSignatures,
		})

		h := common.Hash{}

		hw.Sum(h[:0])
		senderTxHash := rawTx.GetTxInternalData().SenderTxHash()
		assert.Equal(t, h, senderTxHash)

	case *TxInternalDataSmartContractExecution:
		senderTxHash := rawTx.GetTxInternalData().SenderTxHash()
		assert.Equal(t, rawTx.Hash(), senderTxHash)

	case *TxInternalDataFeeDelegatedSmartContractExecution:
		hw := sha3.NewKeccak256()
		rlp.Encode(hw, rawTx.Type())
		rlp.Encode(hw, []interface{}{
			v.AccountNonce,
			v.Price,
			v.GasLimit,
			v.Recipient,
			v.Amount,
			v.From,
			v.Payload,
			v.TxSignatures,
		})

		h := common.Hash{}

		hw.Sum(h[:0])
		senderTxHash := rawTx.GetTxInternalData().SenderTxHash()
		assert.Equal(t, h, senderTxHash)

	case *TxInternalDataFeeDelegatedSmartContractExecutionWithRatio:
		hw := sha3.NewKeccak256()
		rlp.Encode(hw, rawTx.Type())
		rlp.Encode(hw, []interface{}{
			v.AccountNonce,
			v.Price,
			v.GasLimit,
			v.Recipient,
			v.Amount,
			v.From,
			v.Payload,
			v.FeeRatio,
			v.TxSignatures,
		})

		h := common.Hash{}

		hw.Sum(h[:0])
		senderTxHash := rawTx.GetTxInternalData().SenderTxHash()
		assert.Equal(t, h, senderTxHash)

	case *TxInternalDataCancel:
		senderTxHash := rawTx.GetTxInternalData().SenderTxHash()
		assert.Equal(t, rawTx.Hash(), senderTxHash)

	case *TxInternalDataFeeDelegatedCancel:
		hw := sha3.NewKeccak256()
		rlp.Encode(hw, rawTx.Type())
		rlp.Encode(hw, []interface{}{
			v.AccountNonce,
			v.Price,
			v.GasLimit,
			v.From,
			v.TxSignatures,
		})

		h := common.Hash{}

		hw.Sum(h[:0])
		senderTxHash := rawTx.GetTxInternalData().SenderTxHash()
		assert.Equal(t, h, senderTxHash)

	case *TxInternalDataFeeDelegatedCancelWithRatio:
		hw := sha3.NewKeccak256()
		rlp.Encode(hw, rawTx.Type())
		rlp.Encode(hw, []interface{}{
			v.AccountNonce,
			v.Price,
			v.GasLimit,
			v.From,
			v.FeeRatio,
			v.TxSignatures,
		})

		h := common.Hash{}

		hw.Sum(h[:0])
		senderTxHash := rawTx.GetTxInternalData().SenderTxHash()
		assert.Equal(t, h, senderTxHash)

	case *TxInternalDataChainDataAnchoring:
		senderTxHash := rawTx.GetTxInternalData().SenderTxHash()
		assert.Equal(t, rawTx.Hash(), senderTxHash)

	case *TxInternalDataFeeDelegatedChainDataAnchoring:
		hw := sha3.NewKeccak256()
		rlp.Encode(hw, rawTx.Type())
		rlp.Encode(hw, []interface{}{
			v.AccountNonce,
			v.Price,
			v.GasLimit,
			v.From,
			v.Payload,
			v.TxSignatures,
		})

		h := common.Hash{}

		hw.Sum(h[:0])
		senderTxHash := rawTx.GetTxInternalData().SenderTxHash()
		assert.Equal(t, h, senderTxHash)

	case *TxInternalDataFeeDelegatedChainDataAnchoringWithRatio:
		hw := sha3.NewKeccak256()
		rlp.Encode(hw, rawTx.Type())
		rlp.Encode(hw, []interface{}{
			v.AccountNonce,
			v.Price,
			v.GasLimit,
			v.From,
			v.Payload,
			v.FeeRatio,
			v.TxSignatures,
		})

		h := common.Hash{}

		hw.Sum(h[:0])
		senderTxHash := rawTx.GetTxInternalData().SenderTxHash()
		assert.Equal(t, h, senderTxHash)
	case *TxInternalDataEthereumAccessList:
		hw := sha3.NewKeccak256()
		rlp.Encode(hw, byte(rawTx.Type()))
		rlp.Encode(hw, []interface{}{
			v.ChainID,
			v.AccountNonce,
			v.Price,
			v.GasLimit,
			v.Recipient,
			v.Amount,
			v.Payload,
			v.AccessList,
			v.V,
			v.R,
			v.S,
		})

		h := common.Hash{}

		hw.Sum(h[:0])
		senderTxHash := rawTx.GetTxInternalData().SenderTxHash()
		assert.Equal(t, h, senderTxHash)
	case *TxInternalDataEthereumDynamicFee:
		hw := sha3.NewKeccak256()
		rlp.Encode(hw, byte(rawTx.Type()))
		rlp.Encode(hw, []interface{}{
			v.ChainID,
			v.AccountNonce,
			v.GasTipCap,
			v.GasFeeCap,
			v.GasLimit,
			v.Recipient,
			v.Amount,
			v.Payload,
			v.AccessList,
			v.V,
			v.R,
			v.S,
		})

		h := common.Hash{}

		hw.Sum(h[:0])
		senderTxHash := rawTx.GetTxInternalData().SenderTxHash()
		assert.Equal(t, h, senderTxHash)
	default:
		t.Fatal("Undefined tx type.")
	}
}
