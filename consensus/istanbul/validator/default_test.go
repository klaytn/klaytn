// Modifications Copyright 2018 The klaytn Authors
// Copyright 2017 The go-ethereum Authors
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
// This file is derived from quorum/consensus/istanbul/validator/default_test.go (2018/06/04).
// Modified and improved for the klaytn development.

package validator

import (
	"fmt"
	"math/big"
	"reflect"
	"strings"
	"testing"

	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/consensus/istanbul"
	"github.com/klaytn/klaytn/crypto"
	"github.com/stretchr/testify/assert"
)

var (
	testAddress  = "136807B12327a8AfF9831F09617dA1B9D398cda2"
	testAddress2 = "4dd324F9821485caE941640B32c3Bcf1fA6E93E6"
	testAddress3 = "62E47d858bf8513fc401886B94E33e7DCec2Bfb7"
	testAddress4 = "8aD8F547fa00f58A8c4fb3B671Ee5f1A75bA028a"
	testAddress5 = "dd197E88fd97aF3877023cf20d69543fc72e6298"
)

func TestNewValidatorSet(t *testing.T) {
	var validators []istanbul.Validator
	const ValCnt = 100

	// Create 100 validators with random addresses
	b := []byte{}
	for i := 0; i < ValCnt; i++ {
		key, _ := crypto.GenerateKey()
		addr := crypto.PubkeyToAddress(key.PublicKey)
		val := New(addr)
		validators = append(validators, val)
		b = append(b, val.Address().Bytes()...)
	}

	// Create ValidatorSet
	valSet := NewSet(ExtractValidators(b), istanbul.RoundRobin)
	if valSet == nil {
		t.Errorf("the validator byte array cannot be parsed")
		t.FailNow()
	}

	// Check validators sorting: should be in ascending order
	for i := 0; i < ValCnt-1; i++ {
		val := valSet.GetByIndex(uint64(i))
		nextVal := valSet.GetByIndex(uint64(i + 1))
		if strings.Compare(val.String(), nextVal.String()) >= 0 {
			t.Errorf("validator set is not sorted in descending order")
		}
	}
}

func TestNormalValSet(t *testing.T) {
	b1 := common.Hex2Bytes(testAddress)
	b2 := common.Hex2Bytes(testAddress2)
	addr1 := common.BytesToAddress(b1)
	addr2 := common.BytesToAddress(b2)
	val1 := New(addr1)
	val2 := New(addr2)

	valSet := newDefaultSet([]common.Address{addr1, addr2}, istanbul.RoundRobin)
	if valSet == nil {
		t.Errorf("the format of validator set is invalid")
		t.FailNow()
	}

	// check size
	if size := valSet.Size(); size != 2 {
		t.Errorf("the size of validator set is wrong: have %v, want 2", size)
	}
	// test get by index
	if val := valSet.GetByIndex(uint64(0)); !reflect.DeepEqual(val, val1) {
		t.Errorf("validator mismatch: have %v, want %v", val, val1)
	}
	// test get by invalid index
	if val := valSet.GetByIndex(uint64(2)); val != nil {
		t.Errorf("validator mismatch: have %v, want nil", val)
	}
	// test get by address
	if _, val := valSet.GetByAddress(addr2); !reflect.DeepEqual(val, val2) {
		t.Errorf("validator mismatch: have %v, want %v", val, val2)
	}
	// test get by invalid address
	invalidAddr := common.HexToAddress("0x9535b2e7faaba5288511d89341d94a38063a349b")
	if _, val := valSet.GetByAddress(invalidAddr); val != nil {
		t.Errorf("validator mismatch: have %v, want nil", val)
	}
	// test get proposer
	if val := valSet.GetProposer(); !reflect.DeepEqual(val, val1) {
		t.Errorf("proposer mismatch: have %v, want %v", val, val1)
	}
	// test calculate proposer
	lastProposer := addr1
	valSet.CalcProposer(lastProposer, uint64(0))
	if val := valSet.GetProposer(); !reflect.DeepEqual(val, val2) {
		t.Errorf("proposer mismatch: have %v, want %v", val, val2)
	}
	valSet.CalcProposer(lastProposer, uint64(3))
	if val := valSet.GetProposer(); !reflect.DeepEqual(val, val1) {
		t.Errorf("proposer mismatch: have %v, want %v", val, val1)
	}
	// test empty last proposer
	lastProposer = common.Address{}
	valSet.CalcProposer(lastProposer, uint64(3))
	if val := valSet.GetProposer(); !reflect.DeepEqual(val, val2) {
		t.Errorf("proposer mismatch: have %v, want %v", val, val2)
	}
}

func TestEmptyValSet(t *testing.T) {
	valSet := NewSet(ExtractValidators([]byte{}), istanbul.RoundRobin)
	if valSet == nil {
		t.Errorf("validator set should not be nil")
	}
}

func TestAddAndRemoveValidator(t *testing.T) {
	valSet := NewSet(ExtractValidators([]byte{}), istanbul.RoundRobin)
	if !valSet.AddValidator(common.StringToAddress(string(rune(2)))) {
		t.Error("the validator should be added")
	}
	if valSet.AddValidator(common.StringToAddress(string(rune(2)))) {
		t.Error("the existing validator should not be added")
	}
	valSet.AddValidator(common.StringToAddress(string(rune(1))))
	valSet.AddValidator(common.StringToAddress(string(rune(0))))
	if len(valSet.List()) != 3 {
		t.Error("the size of validator set should be 3")
	}

	for i, v := range valSet.List() {
		expected := common.StringToAddress(string(rune(i)))
		if v.Address() != expected {
			t.Errorf("the order of validators is wrong: have %v, want %v", v.Address().Hex(), expected.Hex())
		}
	}

	if !valSet.RemoveValidator(common.StringToAddress(string(rune(2)))) {
		t.Error("the validator should be removed")
	}
	if valSet.RemoveValidator(common.StringToAddress(string(rune(2)))) {
		t.Error("the non-existing validator should not be removed")
	}
	if len(valSet.List()) != 2 {
		t.Error("the size of validator set should be 2")
	}
	valSet.RemoveValidator(common.StringToAddress(string(rune(1))))
	if len(valSet.List()) != 1 {
		t.Error("the size of validator set should be 1")
	}
	valSet.RemoveValidator(common.StringToAddress(string(rune(0))))
	if len(valSet.List()) != 0 {
		t.Error("the size of validator set should be 0")
	}
}

func TestStickyProposer(t *testing.T) {
	b1 := common.Hex2Bytes(testAddress)
	b2 := common.Hex2Bytes(testAddress2)
	addr1 := common.BytesToAddress(b1)
	addr2 := common.BytesToAddress(b2)
	val1 := New(addr1)
	val2 := New(addr2)

	valSet := newDefaultSet([]common.Address{addr1, addr2}, istanbul.Sticky)

	// test get proposer
	if val := valSet.GetProposer(); !reflect.DeepEqual(val, val1) {
		t.Errorf("proposer mismatch: have %v, want %v", val, val1)
	}
	// test calculate proposer
	lastProposer := addr1
	valSet.CalcProposer(lastProposer, uint64(0))
	if val := valSet.GetProposer(); !reflect.DeepEqual(val, val1) {
		t.Errorf("proposer mismatch: have %v, want %v", val, val1)
	}

	valSet.CalcProposer(lastProposer, uint64(1))
	if val := valSet.GetProposer(); !reflect.DeepEqual(val, val2) {
		t.Errorf("proposer mismatch: have %v, want %v", val, val2)
	}
	// test empty last proposer
	lastProposer = common.Address{}
	valSet.CalcProposer(lastProposer, uint64(3))
	if val := valSet.GetProposer(); !reflect.DeepEqual(val, val2) {
		t.Errorf("proposer mismatch: have %v, want %v", val, val2)
	}
}

func TestDefaultSet_SubList(t *testing.T) {
	b1 := common.Hex2Bytes(testAddress)
	b2 := common.Hex2Bytes(testAddress2)
	b3 := common.Hex2Bytes(testAddress3)
	b4 := common.Hex2Bytes(testAddress4)
	b5 := common.Hex2Bytes(testAddress5)
	addr1 := common.BytesToAddress(b1)
	addr2 := common.BytesToAddress(b2)
	addr3 := common.BytesToAddress(b3)
	addr4 := common.BytesToAddress(b4)
	addr5 := common.BytesToAddress(b5)
	testAddresses := []common.Address{addr1, addr2, addr3, addr4, addr5}

	valSet := NewSet(testAddresses, istanbul.RoundRobin)
	if valSet == nil {
		t.Errorf("the format of validator set is invalid")
		t.FailNow()
	}
	valSet.SetSubGroupSize(3)

	hash := istanbul.RLPHash("This is a test hash")
	view := &istanbul.View{
		Sequence: new(big.Int).SetInt64(1),
		Round:    new(big.Int).SetInt64(0),
	}

	lenAddress := len(testAddresses)
	for i := 0; i < lenAddress*2; i++ {
		currentProposer := valSet.GetProposer()
		assert.Equal(t, testAddresses[i%lenAddress], currentProposer.Address())

		committee := valSet.SubList(hash, view, true)

		assert.Equal(t, testAddresses[i%lenAddress].String(), committee[0].String())
		assert.Equal(t, testAddresses[(i+1)%lenAddress].String(), committee[1].String())

		valSet.CalcProposer(currentProposer.Address(), view.Round.Uint64())
	}
}

func TestDefaultSet_Copy(t *testing.T) {
	b1 := common.Hex2Bytes(testAddress)
	b2 := common.Hex2Bytes(testAddress2)
	b3 := common.Hex2Bytes(testAddress3)
	b4 := common.Hex2Bytes(testAddress4)
	b5 := common.Hex2Bytes(testAddress5)
	addr1 := common.BytesToAddress(b1)
	addr2 := common.BytesToAddress(b2)
	addr3 := common.BytesToAddress(b3)
	addr4 := common.BytesToAddress(b4)
	addr5 := common.BytesToAddress(b5)
	testAddresses := []common.Address{addr1, addr2, addr3, addr4, addr5}

	valSet := NewSet(testAddresses, istanbul.RoundRobin)
	copiedValSet := valSet.Copy()

	assert.NotEqual(t, fmt.Sprintf("%p", &valSet), fmt.Sprintf("%p", &copiedValSet))

	assert.Equal(t, valSet.List(), copiedValSet.List())
	assert.NotEqual(t, fmt.Sprintf("%p", valSet.List()), fmt.Sprintf("%p", copiedValSet.List()))

	for i := uint64(0); i < valSet.Size(); i++ {
		assert.Equal(t, valSet.List()[i], copiedValSet.List()[i])
		assert.NotEqual(t, fmt.Sprintf("%p", valSet.List()[i]), fmt.Sprintf("%p", copiedValSet.List())[i])
	}

	assert.Equal(t, valSet.GetProposer(), copiedValSet.GetProposer())
	assert.NotEqual(t, fmt.Sprintf("%p", valSet.GetProposer()), fmt.Sprintf("%p", copiedValSet.GetProposer()))

	assert.Equal(t, valSet.Policy(), copiedValSet.Policy())
	assert.Equal(t, valSet.SubGroupSize(), copiedValSet.SubGroupSize())
	assert.Equal(t, valSet.TotalVotingPower(), copiedValSet.TotalVotingPower())
}
