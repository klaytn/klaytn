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

import "testing"

// TestFeeRatioCheck checks the txs with a fee-ratio implement GetFeeRatio() or not.
// This prohibits the case that GetFeeRatio() is not implemented for TxTypeFeeDelegatedWithRatio types.
func TestFeeRatioCheck(t *testing.T) {
	for i := TxTypeLegacyTransaction; i < TxTypeEthereumLast; i++ {
		if i == TxTypeKlaytnLast {
			i = TxTypeEthereumAccessList
		}

		tx, err := NewTxInternalData(i)
		if err == nil && tx.Type().IsFeeDelegatedWithRatioTransaction() {
			if _, ok := tx.(TxInternalDataFeeRatio); !ok {
				t.Fatalf("GetFeeRatio() is not implemented. tx=%s", tx.String())
			}
		}
	}
}

// TestFeeDelegatedCheck checks that fee-delegated tx types implement GetFeePayer() or not.
// This prohibits the case that GetFeePayer() is not implemented for TxTypeFeeDelegatedXXX types.
func TestFeeDelegatedCheck(t *testing.T) {
	for i := TxTypeLegacyTransaction; i < TxTypeEthereumLast; i++ {
		if i == TxTypeKlaytnLast {
			i = TxTypeEthereumAccessList
		}

		tx, err := NewTxInternalData(i)
		if err == nil && tx.Type().IsFeeDelegatedTransaction() {
			if _, ok := tx.(TxInternalDataFeePayer); !ok {
				t.Fatalf("GetFeePayer() is not implemented. tx=%s", tx.String())
			}
		}
	}

}
