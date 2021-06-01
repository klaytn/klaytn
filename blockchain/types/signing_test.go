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

package types

import (
	"math/big"
	"testing"

	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/crypto"
)

func BenchmarkSingleRecoverEIP155Signer(b *testing.B) {
	signer := NewEIP155Signer(big.NewInt(2018))

	from, _ := crypto.GenerateKey()
	to, _ := crypto.GenerateKey()

	toAddr := crypto.PubkeyToAddress(to.PublicKey)

	tx := NewTransaction(
		0,
		toAddr,
		big.NewInt(1),
		100000,
		big.NewInt(0),
		nil)

	signTx, _ := SignTx(tx, NewEIP155Signer(big.NewInt(2018)), from)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = signer.Sender(signTx)
	}
}

func BenchmarkDoubleRecoverEIP155Signer(b *testing.B) {
	signer := NewEIP155Signer(big.NewInt(2018))

	from, _ := crypto.GenerateKey()
	to, _ := crypto.GenerateKey()

	toAddr := crypto.PubkeyToAddress(to.PublicKey)

	tx := NewTransaction(
		0,
		toAddr,
		big.NewInt(1),
		100000,
		big.NewInt(0),
		nil)

	signTx, _ := SignTx(tx, NewEIP155Signer(big.NewInt(2018)), from)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = signer.Sender(signTx)
		_, _ = signer.Sender(signTx)
	}
}

func BenchmarkDoubleRecoverEIP155SignerDoubleGoroutine(b *testing.B) {
	signer := NewEIP155Signer(big.NewInt(2018))

	from, _ := crypto.GenerateKey()
	to, _ := crypto.GenerateKey()

	toAddr := crypto.PubkeyToAddress(to.PublicKey)

	tx := NewTransaction(
		0,
		toAddr,
		big.NewInt(1),
		100000,
		big.NewInt(0),
		nil)

	signTx, _ := SignTx(tx, NewEIP155Signer(big.NewInt(2018)), from)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		channel1 := make(chan common.Address)
		channel2 := make(chan common.Address)
		go func() {
			addr, _ := signer.Sender(signTx)
			channel1 <- addr
		}()
		go func() {
			addr, _ := signer.Sender(signTx)
			channel2 <- addr
		}()
		_ = <-channel1
		_ = <-channel2
	}
}

func BenchmarkDoubleRecoverEIP155SignerSingleGoroutine(b *testing.B) {
	signer := NewEIP155Signer(big.NewInt(2018))

	from, _ := crypto.GenerateKey()
	to, _ := crypto.GenerateKey()

	toAddr := crypto.PubkeyToAddress(to.PublicKey)

	tx := NewTransaction(
		0,
		toAddr,
		big.NewInt(1),
		100000,
		big.NewInt(0),
		nil)

	signTx, _ := SignTx(tx, NewEIP155Signer(big.NewInt(2018)), from)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		channel1 := make(chan common.Address)
		go func() {
			addr, _ := signer.Sender(signTx)
			channel1 <- addr
		}()
		_, _ = signer.Sender(signTx)
		_ = <-channel1
	}
}

type param struct {
	signTx   *Transaction
	addrChan chan common.Address
}

func launchSenderGoroutines(signer *EIP155Signer, num int, paramCh chan param, quitCh chan struct{}) {
	for i := 0; i < num; i++ {
		go func() {
			for {
				select {
				case p := <-paramCh:
					addr, _ := signer.Sender(p.signTx)
					p.addrChan <- addr
				case <-quitCh:
					return
				}
			}
		}()
	}
}

func BenchmarkDoubleRecoverEIP155SignerReusedGoroutines(b *testing.B) {
	signer := NewEIP155Signer(big.NewInt(2018))

	from, _ := crypto.GenerateKey()
	to, _ := crypto.GenerateKey()

	toAddr := crypto.PubkeyToAddress(to.PublicKey)

	tx := NewTransaction(
		0,
		toAddr,
		big.NewInt(1),
		100000,
		big.NewInt(0),
		nil)

	signTx, _ := SignTx(tx, NewEIP155Signer(big.NewInt(2018)), from)

	numGoroutine := 2
	paramCh := make(chan param, numGoroutine)

	quitCh := make(chan struct{}, numGoroutine)
	launchSenderGoroutines(&signer, numGoroutine, paramCh, quitCh)

	ch0 := make(chan common.Address)
	ch1 := make(chan common.Address)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		paramCh <- param{signTx, ch0}
		paramCh <- param{signTx, ch1}

		_ = <-ch0
		_ = <-ch1
	}
	b.StopTimer()

	for i := 0; i < numGoroutine; i++ {
		quitCh <- struct{}{}
	}
}
