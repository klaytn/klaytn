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

package tests

import (
	"crypto/ecdsa"
	"strings"
	"testing"

	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/stretchr/testify/assert"
)

func BenchmarkTxHash(b *testing.B) {
	testfns := []genTx{
		genLegacyValueTransfer,
		genNewValueTransfer,
		genLegacySmartContractDeploy,
		genNewSmartContractDeploy,
		genLegacySmartContractExecution,
		genNewSmartContractExecution,
		genNewAccountUpdateAccountKeyPublic,
		genNewFeeDelegatedValueTransfer,
		genNewFeeDelegatedValueTransferWithRatio,
		genNewCancel,
	}

	for _, fn := range testfns {
		fnname := getFunctionName(fn)
		fnname = fnname[strings.LastIndex(fnname, ".")+1:]
		if strings.Contains(fnname, "New") {
			benchName = "New/" + strings.Split(fnname, "New")[1]
		} else {
			benchName = "Legacy/" + strings.Split(fnname, "Legacy")[1]
		}
		b.Run(benchName, func(b *testing.B) {
			benchmarkTxHash(b, fn)
		})
	}
}

func benchmarkTxHash(b *testing.B, genTx genTx) {
	signer := types.LatestSignerForChainID(common.Big1)

	// Initialize blockchain
	bcdata, err := NewBCData(6, 4)
	if err != nil {
		b.Fatal(err)
	}
	defer bcdata.Shutdown()

	// reservoir account
	reservoir := &TestAccountType{
		Addr:  *bcdata.addrs[0],
		Keys:  []*ecdsa.PrivateKey{bcdata.privKeys[0]},
		Nonce: uint64(0),
	}

	anon, err := createAnonymousAccount("ed580f5bd71a2ee4dae5cb43e331b7d0318596e561e6add7844271ed94156b20")
	assert.Equal(b, nil, err)

	tx := genTx(signer, reservoir, anon)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		signer.Hash(tx)
	}
}
