// Copyright 2018 The klaytn Authors
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
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/blockchain/types/derivesha"
	"github.com/klaytn/klaytn/common/profile"
	"github.com/klaytn/klaytn/crypto"
)

func BenchmarkDeriveSha(b *testing.B) {
	funcs := map[string]derivesha.IDeriveSha{
		"Orig":   derivesha.DeriveShaOrig{},
		"Simple": derivesha.DeriveShaSimple{},
		"Concat": derivesha.DeriveShaConcat{},
	}

	NTS := []int{1000}

	for k, f := range funcs {
		for _, nt := range NTS {
			testName := fmt.Sprintf("%s,%d", k, nt)
			b.Run(testName, func(b *testing.B) {
				benchDeriveSha(b, nt, 4, f)
			})
		}
	}
}

func BenchmarkDeriveShaSingleAccount(b *testing.B) {
	txs, err := genTxs(4000)
	if err != nil {
		b.Fatal(err)
	}

	funcs := map[string]derivesha.IDeriveSha{
		"Orig":   derivesha.DeriveShaOrig{},
		"Simple": derivesha.DeriveShaSimple{},
		"Concat": derivesha.DeriveShaConcat{},
	}

	for k, f := range funcs {
		b.Run(k, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				f.DeriveSha(txs)
			}
		})
	}
}

func BenchmarkDeriveShaOrig(b *testing.B) {
	b.Run("1000test-stackTrie", func(b *testing.B) {
		benchDeriveSha(b, 1000, 4, derivesha.DeriveShaOrig{})
	})
}

func benchDeriveSha(b *testing.B, numTransactions, numValidators int, sha derivesha.IDeriveSha) {
	// Initialize blockchain
	start := time.Now()
	maxAccounts := numTransactions * 2
	bcdata, err := NewBCData(maxAccounts, numValidators)
	if err != nil {
		b.Fatal(err)
	}
	profile.Prof.Profile("main_init_blockchain", time.Now().Sub(start))
	defer bcdata.Shutdown()

	// Initialize address-balance map for verification
	start = time.Now()
	accountMap := NewAccountMap()
	if err := accountMap.Initialize(bcdata); err != nil {
		b.Fatal(err)
	}
	profile.Prof.Profile("main_init_accountMap", time.Now().Sub(start))

	signer := types.MakeSigner(bcdata.bc.Config(), bcdata.bc.CurrentBlock().Number())

	txs, err := makeTransactions(accountMap,
		bcdata.addrs[0:numTransactions], bcdata.privKeys[0:numTransactions],
		signer, bcdata.addrs[numTransactions:numTransactions*2], nil, 0, false)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sha.DeriveSha(txs)
	}
}

func genTxs(num uint64) (types.Transactions, error) {
	key, err := crypto.HexToECDSA("deadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef")
	if err != nil {
		return nil, err
	}
	addr := crypto.PubkeyToAddress(key.PublicKey)
	newTx := func(i uint64) (*types.Transaction, error) {
		signer := types.NewEIP155Signer(big.NewInt(18))
		utx := types.NewTransaction(i, addr, new(big.Int), 0, new(big.Int).SetUint64(10000000), nil)
		tx, err := types.SignTx(utx, signer, key)
		return tx, err
	}
	var txs types.Transactions
	for i := uint64(0); i < num; i++ {
		tx, err := newTx(i)
		if err != nil {
			return nil, err
		}
		txs = append(txs, tx)
	}
	return txs, nil
}
