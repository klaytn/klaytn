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
	"encoding/json"
	"io"

	"github.com/klaytn/klaytn/rlp"
)

// TxInternalDataSerializer serializes an object that implements `TxInternalData`.
type TxInternalDataSerializer struct {
	txType TxType
	tx     TxInternalData
}

// txInternalDataJSON is an internal object for JSON serialization.
type txInternalDataJSON struct {
	TxType TxType `json:"typeInt"`
}

// newTxInternalDataSerializerWithValues creates a new TxInternalDataSerializer object with the given TxInternalData object.
func newTxInternalDataSerializerWithValues(t TxInternalData) *TxInternalDataSerializer {
	return &TxInternalDataSerializer{t.Type(), t}
}

// newTxInternalDataSerializer creates an empty TxInternalDataSerializer object for decoding.
func newTxInternalDataSerializer() *TxInternalDataSerializer {
	return &TxInternalDataSerializer{}
}

func (serializer *TxInternalDataSerializer) EncodeRLP(w io.Writer) error {
	// if it is the original transaction, do not encode type.
	if serializer.txType == TxTypeLegacyTransaction {
		return rlp.Encode(w, serializer.tx)
	}

	if serializer.txType.IsEthTypedTransaction() {
		ethType := byte(serializer.txType)
		if _, err := w.Write([]byte{byte(EthereumTxTypeEnvelope), ethType}); err != nil {
			return err
		}
	} else {
		if err := rlp.Encode(w, serializer.txType); err != nil {
			return err
		}
	}

	return rlp.Encode(w, serializer.tx)
}

func (serializer *TxInternalDataSerializer) DecodeRLP(s *rlp.Stream) error {
	if err := s.Decode(&serializer.txType); err != nil {
		// fallback to the original transaction decoding.
		txd := newEmptyTxInternalDataLegacy()
		if err := s.Decode(txd); err != nil {
			return err
		}
		serializer.txType = TxTypeLegacyTransaction
		serializer.tx = txd
		return nil
	}

	if serializer.txType == EthereumTxTypeEnvelope {
		var ethType TxType
		if err := s.Decode(&ethType); err != nil {
			return err
		}
		serializer.txType = serializer.txType<<8 | ethType
	}

	var err error
	serializer.tx, err = NewTxInternalData(serializer.txType)
	if err != nil {
		return err
	}

	return s.Decode(serializer.tx)
}

func (serializer *TxInternalDataSerializer) MarshalJSON() ([]byte, error) {
	// if it is the original transaction, do not marshal type.
	if serializer.txType == TxTypeLegacyTransaction {
		return json.Marshal(serializer.tx)
	}

	return json.Marshal(serializer.tx)
}

func (serializer *TxInternalDataSerializer) UnmarshalJSON(b []byte) error {
	var dec txInternalDataJSON

	if err := json.Unmarshal(b, &dec); err != nil {
		return err
	}

	if dec.TxType == TxTypeLegacyTransaction {
		// fallback to unmarshal the legacy transaction.
		txd := newTxInternalDataLegacy()
		if err := json.Unmarshal(b, txd); err != nil {
			return err
		}
		serializer.txType = TxTypeLegacyTransaction
		serializer.tx = txd

		return nil
	}

	serializer.txType = dec.TxType

	var err error
	serializer.tx, err = NewTxInternalData(serializer.txType)
	if err != nil {
		return err
	}

	return json.Unmarshal(b, serializer.tx)
}
