// Copyright 2023 The klaytn Authors
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

package backends

import (
	"context"
	"errors"
	"math/big"
	"strings"
	"testing"

	"github.com/klaytn/klaytn"
	"github.com/klaytn/klaytn/accounts/abi"
	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/blockchain/vm"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/consensus/gxhash"
	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/params"
	"github.com/klaytn/klaytn/storage/database"
	"github.com/stretchr/testify/assert"
)

var (
	testAddr  = crypto.PubkeyToAddress(testKey.PublicKey)
	code1Addr = common.HexToAddress("0x1111111111111111111111111111111111111111")
	code2Addr = common.HexToAddress("0x2222222222222222222222222222222222222222")

	parsedAbi1, _ = abi.JSON(strings.NewReader(abiJSON))
	parsedAbi2, _ = abi.JSON(strings.NewReader(reverterABI))
	code1Bytes    = common.FromHex(deployedCode)
	code2Bytes    = common.FromHex(reverterDeployedBin)
)

func newTestBlockchain() *blockchain.BlockChain {
	config := params.TestChainConfig.Copy()
	alloc := blockchain.GenesisAlloc{
		testAddr:  {Balance: big.NewInt(10000000000)},
		code1Addr: {Balance: big.NewInt(0), Code: code1Bytes},
		code2Addr: {Balance: big.NewInt(0), Code: code2Bytes},
	}

	db := database.NewMemoryDBManager()
	genesis := blockchain.Genesis{Config: config, Alloc: alloc}
	genesis.MustCommit(db)

	bc, _ := blockchain.NewBlockChain(db, nil, genesis.Config, gxhash.NewFaker(), vm.Config{})

	// Append 10 blocks to test with block numbers other than 0
	block := bc.CurrentBlock()
	blocks, _ := blockchain.GenerateChain(config, block, gxhash.NewFaker(), db, 10, func(i int, b *blockchain.BlockGen) {})
	bc.InsertChain(blocks)

	return bc
}

func TestBlockchainCodeAt(t *testing.T) {
	bc := newTestBlockchain()
	c := NewBlockchainContractCaller(bc)

	// Normal cases
	code, err := c.CodeAt(context.Background(), code1Addr, nil)
	assert.Nil(t, err)
	assert.Equal(t, code1Bytes, code)

	code, err = c.CodeAt(context.Background(), code2Addr, nil)
	assert.Nil(t, err)
	assert.Equal(t, code2Bytes, code)

	code, err = c.CodeAt(context.Background(), code1Addr, common.Big0)
	assert.Nil(t, err)
	assert.Equal(t, code1Bytes, code)

	code, err = c.CodeAt(context.Background(), code1Addr, common.Big1)
	assert.Nil(t, err)
	assert.Equal(t, code1Bytes, code)

	code, err = c.CodeAt(context.Background(), code1Addr, big.NewInt(10))
	assert.Nil(t, err)
	assert.Equal(t, code1Bytes, code)

	// Non-code address
	code, err = c.CodeAt(context.Background(), testAddr, nil)
	assert.True(t, code == nil || err != nil)

	// Invalid block number
	code, err = c.CodeAt(context.Background(), code1Addr, big.NewInt(11))
	assert.True(t, code == nil || err != nil)
}

func TestBlockchainCallContract(t *testing.T) {
	bc := newTestBlockchain()
	c := NewBlockchainContractCaller(bc)

	data_receive, _ := parsedAbi1.Pack("receive", []byte("X"))
	data_revertString, _ := parsedAbi2.Pack("revertString")
	data_revertNoString, _ := parsedAbi2.Pack("revertNoString")

	// Normal case
	ret, err := c.CallContract(context.Background(), klaytn.CallMsg{
		From: testAddr,
		To:   &code1Addr,
		Gas:  1000000,
		Data: data_receive,
	}, nil)
	assert.Nil(t, err)
	assert.Equal(t, expectedReturn, ret)

	// Error outside VM - Intrinsic Gas
	ret, err = c.CallContract(context.Background(), klaytn.CallMsg{
		From: testAddr,
		To:   &code1Addr,
		Gas:  20000,
		Data: data_receive,
	}, nil)
	assert.True(t, errors.Is(err, blockchain.ErrIntrinsicGas))

	// VM revert error - empty reason
	ret, err = c.CallContract(context.Background(), klaytn.CallMsg{
		From: testAddr,
		To:   &code2Addr,
		Gas:  100000,
		Data: data_revertNoString,
	}, nil)
	assert.Equal(t, "execution reverted: ", err.Error())

	// VM revert error - string reason
	ret, err = c.CallContract(context.Background(), klaytn.CallMsg{
		From: testAddr,
		To:   &code2Addr,
		Gas:  100000,
		Data: data_revertString,
	}, nil)
	assert.Equal(t, "execution reverted: some error", err.Error())
}
