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

package sc

import (
	"github.com/klaytn/klaytn/accounts/abi/bind"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/contracts/sc_erc20"
	"github.com/klaytn/klaytn/contracts/sc_erc721"
	"github.com/stretchr/testify/assert"
	"math/big"
	"os"
	"path"
	"strconv"
	"testing"
)

// TestTokenPublicVariables checks the results of the public variables.
func TestTokenPublicVariables(t *testing.T) {
	defer func() {
		if err := os.Remove(path.Join(os.TempDir(), BridgeAddrJournal)); err != nil {
			t.Fatalf("fail to delete journal file %v", err)
		}
	}()

	info := prepare(t, func(info *testInfo) {
		for i := 0; i < testTxCount; i++ {
			ops[KLAY].request(info, info.localInfo)
		}
	})

	initSupply, err := info.tokenLocalBridge.INITIALSUPPLY(nil)
	assert.Equal(t, "1000000000000000000000000000", initSupply.String())

	allowance, err := info.tokenLocalBridge.Allowance(nil, info.chainAuth.From, info.chainAuth.From)
	assert.Equal(t, "0", allowance.String())

	balance, err := info.tokenLocalBridge.BalanceOf(nil, info.nodeAuth.From)
	assert.Equal(t, "1000000000000000000000000000", balance.String())

	decimal, err := info.tokenLocalBridge.Decimals(nil)
	assert.Equal(t, uint8(0x12), decimal)

	name, err := info.tokenLocalBridge.Name(nil)
	assert.Equal(t, "ServiceChainToken", name)

	symbol, err := info.tokenLocalBridge.Symbol(nil)
	assert.Equal(t, "SCT", symbol)

	_, _, _, err = sctoken.DeployServiceChainToken(info.nodeAuth, info.sim, common.Address{0})
	assert.NotEqual(t, nil, err)
}

// TestTokenPublicVariables checks the results of the public variables.
func TestNFTPublicVariables(t *testing.T) {
	defer func() {
		if err := os.Remove(path.Join(os.TempDir(), BridgeAddrJournal)); err != nil {
			t.Fatalf("fail to delete journal file %v", err)
		}
	}()

	info := prepare(t, func(info *testInfo) {
		for i := 0; i < testTxCount; i++ {
			ops[KLAY].request(info, info.localInfo)
		}
	})

	_, _, _, err := scnft.DeployServiceChainNFT(info.nodeAuth, info.sim, common.Address{0})
	assert.NotEqual(t, nil, err)

	balance, err := info.nftLocalBridge.BalanceOf(nil, info.nodeAuth.From)
	assert.Equal(t, strconv.FormatInt(testTxCount, 10), balance.String())

	bride, err := info.nftLocalBridge.Bridge(nil)
	assert.Equal(t, info.localInfo.address, bride)

	bof, err := info.nftLocalBridge.BalanceOf(nil, info.nodeAuth.From)
	assert.Equal(t, strconv.FormatInt(testTxCount, 10), bof.String())

	approved, err := info.nftLocalBridge.GetApproved(&bind.CallOpts{From: info.nodeAuth.From}, big.NewInt(int64(testNFT)))
	assert.Equal(t, common.Address{0}, approved)

	isApproved, err := info.nftLocalBridge.IsApprovedForAll(nil, info.nodeAuth.From, info.nodeAuth.From)
	assert.Equal(t, false, isApproved)

	isOwner, err := info.nftLocalBridge.IsOwner(&bind.CallOpts{From: info.nodeAuth.From})
	assert.Equal(t, true, isOwner)

	name, err := info.nftLocalBridge.Name(nil)
	assert.Equal(t, "ServiceChainNFT", name)

	owner, err := info.nftLocalBridge.Owner(nil)
	assert.Equal(t, info.nodeAuth.From, owner)

	ownerOf, err := info.nftLocalBridge.OwnerOf(nil, big.NewInt(int64(testNFT)))
	assert.Equal(t, info.nodeAuth.From, ownerOf)

	ifid := [4]byte{0}
	sif, err := info.nftLocalBridge.SupportsInterface(nil, ifid)
	assert.Equal(t, false, sif)

	symbol, err := info.nftLocalBridge.Symbol(nil)
	assert.Equal(t, "SCN", symbol)

	tindex, err := info.nftLocalBridge.TokenByIndex(nil, big.NewInt(int64(0)))
	assert.Equal(t, strconv.FormatInt(testNFT, 10), tindex.String())

	ownerByIndex, err := info.nftLocalBridge.TokenOfOwnerByIndex(nil, info.nodeAuth.From, big.NewInt(int64(0)))
	assert.Equal(t, strconv.FormatInt(testNFT, 10), ownerByIndex.String())

	uri, err := info.nftLocalBridge.TokenURI(nil, big.NewInt(int64(0)))
	assert.Equal(t, "", uri)

	totalSupply, err := info.nftLocalBridge.TotalSupply(nil)
	assert.Equal(t, strconv.FormatInt(testTxCount, 10), totalSupply.String())
}
