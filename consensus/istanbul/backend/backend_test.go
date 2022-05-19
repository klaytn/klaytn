// Modifications Copyright 2020 The klaytn Authors
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
// This file is derived from quorum/consensus/istanbul/backend/backend_test.go (2020/04/16).
// Modified and improved for the klaytn development.

package backend

import (
	"bytes"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/consensus/istanbul"
	"github.com/klaytn/klaytn/consensus/istanbul/validator"
	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/governance"
	"github.com/klaytn/klaytn/params"
	"github.com/klaytn/klaytn/storage/database"
)

var (
	testSigningData = []byte("dummy data")
	// testing node's private key
	PRIVKEY = "ce7671a2880493dfb8d04218707a16b1532dfcac97f0289d770a919d5ff7b068"
	// Max blockNum
	maxBlockNum     = int64(100)
	committeeBlocks = map[Pair]bool{
		{Sequence: 0, Round: 0}:   false,
		{Sequence: 0, Round: 1}:   false,
		{Sequence: 0, Round: 2}:   false,
		{Sequence: 0, Round: 3}:   false,
		{Sequence: 0, Round: 4}:   false,
		{Sequence: 0, Round: 5}:   false,
		{Sequence: 0, Round: 6}:   false,
		{Sequence: 0, Round: 7}:   false,
		{Sequence: 0, Round: 8}:   false,
		{Sequence: 0, Round: 9}:   false,
		{Sequence: 0, Round: 10}:  false,
		{Sequence: 0, Round: 11}:  false,
		{Sequence: 0, Round: 12}:  false,
		{Sequence: 0, Round: 13}:  false,
		{Sequence: 0, Round: 14}:  false,
		{Sequence: 5, Round: 4}:   false,
		{Sequence: 6, Round: 6}:   false,
		{Sequence: 6, Round: 9}:   false,
		{Sequence: 7, Round: 11}:  false,
		{Sequence: 7, Round: 12}:  false,
		{Sequence: 7, Round: 14}:  false,
		{Sequence: 8, Round: 5}:   false,
		{Sequence: 8, Round: 13}:  false,
		{Sequence: 8, Round: 14}:  false,
		{Sequence: 9, Round: 0}:   false,
		{Sequence: 9, Round: 10}:  false,
		{Sequence: 9, Round: 11}:  false,
		{Sequence: 9, Round: 12}:  false,
		{Sequence: 9, Round: 13}:  false,
		{Sequence: 9, Round: 14}:  false,
		{Sequence: 10, Round: 1}:  false,
		{Sequence: 10, Round: 8}:  false,
		{Sequence: 10, Round: 11}: false,
		{Sequence: 10, Round: 12}: false,
		{Sequence: 10, Round: 14}: false,
		{Sequence: 11, Round: 0}:  false,
		{Sequence: 11, Round: 7}:  false,
		{Sequence: 11, Round: 8}:  false,
		{Sequence: 11, Round: 10}: false,
		{Sequence: 11, Round: 11}: false,
		{Sequence: 12, Round: 0}:  false,
		{Sequence: 12, Round: 6}:  false,
		{Sequence: 12, Round: 8}:  false,
		{Sequence: 12, Round: 9}:  false,
		{Sequence: 12, Round: 10}: false,
		{Sequence: 12, Round: 11}: false,
		{Sequence: 12, Round: 13}: false,
		{Sequence: 13, Round: 8}:  false,
		{Sequence: 13, Round: 9}:  false,
		{Sequence: 13, Round: 12}: false,
		{Sequence: 13, Round: 13}: false,
		{Sequence: 14, Round: 0}:  false,
		{Sequence: 14, Round: 5}:  false,
		{Sequence: 14, Round: 7}:  false,
		{Sequence: 14, Round: 8}:  false,
		{Sequence: 14, Round: 10}: false,
		{Sequence: 14, Round: 14}: false,
		{Sequence: 15, Round: 5}:  false,
		{Sequence: 15, Round: 6}:  false,
		{Sequence: 15, Round: 7}:  false,
		{Sequence: 15, Round: 11}: false,
		{Sequence: 16, Round: 5}:  false,
		{Sequence: 16, Round: 6}:  false,
		{Sequence: 16, Round: 7}:  false,
		{Sequence: 17, Round: 4}:  false,
		{Sequence: 17, Round: 5}:  false,
		{Sequence: 17, Round: 7}:  false,
		{Sequence: 17, Round: 9}:  false,
		{Sequence: 18, Round: 1}:  false,
		{Sequence: 18, Round: 3}:  false,
		{Sequence: 18, Round: 4}:  false,
		{Sequence: 18, Round: 7}:  false,
		{Sequence: 18, Round: 9}:  false,
		{Sequence: 18, Round: 10}: false,
		{Sequence: 19, Round: 2}:  false,
		{Sequence: 19, Round: 3}:  false,
		{Sequence: 19, Round: 5}:  false,
		{Sequence: 19, Round: 7}:  false,
		{Sequence: 19, Round: 10}: false,
		{Sequence: 19, Round: 13}: false,
		{Sequence: 20, Round: 1}:  false,
		{Sequence: 20, Round: 2}:  false,
		{Sequence: 20, Round: 3}:  false,
		{Sequence: 20, Round: 11}: false,
		{Sequence: 21, Round: 0}:  false,
		{Sequence: 21, Round: 1}:  false,
		{Sequence: 21, Round: 2}:  false,
		{Sequence: 21, Round: 7}:  false,
		{Sequence: 21, Round: 11}: false,
		{Sequence: 21, Round: 12}: false,
		{Sequence: 22, Round: 0}:  false,
		{Sequence: 22, Round: 3}:  false,
		{Sequence: 22, Round: 6}:  false,
		{Sequence: 23, Round: 2}:  false,
		{Sequence: 23, Round: 10}: false,
		{Sequence: 23, Round: 11}: false,
		{Sequence: 24, Round: 0}:  false,
		{Sequence: 24, Round: 3}:  false,
		{Sequence: 24, Round: 4}:  false,
		{Sequence: 24, Round: 8}:  false,
		{Sequence: 25, Round: 3}:  false,
		{Sequence: 25, Round: 4}:  false,
		{Sequence: 25, Round: 5}:  false,
		{Sequence: 25, Round: 14}: false,
		{Sequence: 26, Round: 11}: false,
		{Sequence: 26, Round: 13}: false,
		{Sequence: 27, Round: 5}:  false,
		{Sequence: 27, Round: 6}:  false,
		{Sequence: 27, Round: 9}:  false,
		{Sequence: 27, Round: 14}: false,
		{Sequence: 28, Round: 2}:  false,
		{Sequence: 28, Round: 3}:  false,
		{Sequence: 28, Round: 5}:  false,
		{Sequence: 28, Round: 9}:  false,
		{Sequence: 28, Round: 13}: false,
		{Sequence: 29, Round: 7}:  false,
		{Sequence: 29, Round: 11}: false,
		{Sequence: 30, Round: 3}:  false,
		{Sequence: 30, Round: 6}:  false,
		{Sequence: 30, Round: 7}:  false,
		{Sequence: 30, Round: 10}: false,
		{Sequence: 30, Round: 12}: false,
		{Sequence: 30, Round: 13}: false,
		{Sequence: 31, Round: 3}:  false,
		{Sequence: 31, Round: 9}:  false,
		{Sequence: 31, Round: 12}: false,
		{Sequence: 32, Round: 0}:  false,
		{Sequence: 32, Round: 1}:  false,
		{Sequence: 32, Round: 5}:  false,
		{Sequence: 32, Round: 12}: false,
		{Sequence: 32, Round: 13}: false,
		{Sequence: 33, Round: 0}:  false,
		{Sequence: 33, Round: 2}:  false,
		{Sequence: 33, Round: 3}:  false,
		{Sequence: 33, Round: 8}:  false,
		{Sequence: 33, Round: 14}: false,
		{Sequence: 34, Round: 0}:  false,
		{Sequence: 34, Round: 3}:  false,
		{Sequence: 34, Round: 6}:  false,
		{Sequence: 34, Round: 7}:  false,
		{Sequence: 34, Round: 8}:  false,
		{Sequence: 34, Round: 10}: false,
		{Sequence: 34, Round: 11}: false,
		{Sequence: 35, Round: 0}:  false,
		{Sequence: 35, Round: 5}:  false,
		{Sequence: 35, Round: 6}:  false,
		{Sequence: 35, Round: 9}:  false,
		{Sequence: 35, Round: 14}: false,
		{Sequence: 36, Round: 1}:  false,
		{Sequence: 37, Round: 0}:  false,
		{Sequence: 37, Round: 1}:  false,
		{Sequence: 37, Round: 5}:  false,
		{Sequence: 37, Round: 8}:  false,
		{Sequence: 37, Round: 12}: false,
		{Sequence: 37, Round: 14}: false,
		{Sequence: 38, Round: 0}:  false,
		{Sequence: 38, Round: 14}: false,
		{Sequence: 39, Round: 0}:  false,
		{Sequence: 39, Round: 1}:  false,
		{Sequence: 39, Round: 3}:  false,
		{Sequence: 39, Round: 10}: false,
		{Sequence: 39, Round: 13}: false,
		{Sequence: 39, Round: 14}: false,
		{Sequence: 40, Round: 3}:  false,
		{Sequence: 40, Round: 10}: false,
		{Sequence: 40, Round: 12}: false,
		{Sequence: 41, Round: 6}:  false,
		{Sequence: 41, Round: 8}:  false,
		{Sequence: 41, Round: 11}: false,
		{Sequence: 42, Round: 3}:  false,
		{Sequence: 42, Round: 5}:  false,
		{Sequence: 42, Round: 6}:  false,
		{Sequence: 42, Round: 11}: false,
		{Sequence: 42, Round: 13}: false,
		{Sequence: 43, Round: 0}:  false,
		{Sequence: 43, Round: 3}:  false,
		{Sequence: 43, Round: 6}:  false,
		{Sequence: 43, Round: 7}:  false,
		{Sequence: 43, Round: 9}:  false,
		{Sequence: 44, Round: 3}:  false,
		{Sequence: 44, Round: 4}:  false,
		{Sequence: 44, Round: 7}:  false,
		{Sequence: 44, Round: 13}: false,
		{Sequence: 44, Round: 14}: false,
		{Sequence: 45, Round: 2}:  false,
		{Sequence: 45, Round: 6}:  false,
		{Sequence: 45, Round: 12}: false,
		{Sequence: 45, Round: 13}: false,
		{Sequence: 46, Round: 3}:  false,
		{Sequence: 46, Round: 4}:  false,
		{Sequence: 46, Round: 7}:  false,
		{Sequence: 46, Round: 8}:  false,
		{Sequence: 46, Round: 10}: false,
		{Sequence: 46, Round: 12}: false,
		{Sequence: 47, Round: 1}:  false,
		{Sequence: 47, Round: 3}:  false,
		{Sequence: 47, Round: 4}:  false,
		{Sequence: 47, Round: 10}: false,
		{Sequence: 47, Round: 12}: false,
		{Sequence: 49, Round: 0}:  false,
		{Sequence: 49, Round: 8}:  false,
		{Sequence: 49, Round: 10}: false,
		{Sequence: 49, Round: 14}: false,
		{Sequence: 50, Round: 3}:  false,
		{Sequence: 50, Round: 4}:  false,
		{Sequence: 50, Round: 10}: false,
		{Sequence: 50, Round: 11}: false,
		{Sequence: 50, Round: 14}: false,
		{Sequence: 52, Round: 1}:  false,
		{Sequence: 52, Round: 3}:  false,
		{Sequence: 52, Round: 7}:  false,
		{Sequence: 52, Round: 11}: false,
		{Sequence: 53, Round: 4}:  false,
		{Sequence: 53, Round: 7}:  false,
		{Sequence: 54, Round: 4}:  false,
		{Sequence: 54, Round: 10}: false,
		{Sequence: 54, Round: 12}: false,
		{Sequence: 55, Round: 2}:  false,
		{Sequence: 55, Round: 12}: false,
		{Sequence: 56, Round: 2}:  false,
		{Sequence: 56, Round: 9}:  false,
		{Sequence: 56, Round: 12}: false,
		{Sequence: 56, Round: 14}: false,
		{Sequence: 57, Round: 7}:  false,
		{Sequence: 57, Round: 13}: false,
		{Sequence: 58, Round: 1}:  false,
		{Sequence: 58, Round: 4}:  false,
		{Sequence: 58, Round: 7}:  false,
		{Sequence: 58, Round: 12}: false,
		{Sequence: 59, Round: 5}:  false,
		{Sequence: 59, Round: 10}: false,
		{Sequence: 59, Round: 13}: false,
		{Sequence: 60, Round: 2}:  false,
		{Sequence: 60, Round: 6}:  false,
		{Sequence: 61, Round: 2}:  false,
		{Sequence: 61, Round: 3}:  false,
		{Sequence: 62, Round: 1}:  false,
		{Sequence: 62, Round: 12}: false,
		{Sequence: 62, Round: 13}: false,
		{Sequence: 63, Round: 1}:  false,
		{Sequence: 63, Round: 2}:  false,
		{Sequence: 63, Round: 5}:  false,
		{Sequence: 63, Round: 7}:  false,
		{Sequence: 63, Round: 9}:  false,
		{Sequence: 63, Round: 11}: false,
		{Sequence: 64, Round: 4}:  false,
		{Sequence: 64, Round: 7}:  false,
		{Sequence: 64, Round: 9}:  false,
		{Sequence: 65, Round: 6}:  false,
		{Sequence: 65, Round: 11}: false,
		{Sequence: 65, Round: 13}: false,
		{Sequence: 65, Round: 14}: false,
		{Sequence: 66, Round: 3}:  false,
		{Sequence: 66, Round: 4}:  false,
		{Sequence: 66, Round: 11}: false,
		{Sequence: 67, Round: 5}:  false,
		{Sequence: 67, Round: 6}:  false,
		{Sequence: 67, Round: 10}: false,
		{Sequence: 67, Round: 11}: false,
		{Sequence: 68, Round: 9}:  false,
		{Sequence: 68, Round: 11}: false,
		{Sequence: 68, Round: 14}: false,
		{Sequence: 69, Round: 2}:  false,
		{Sequence: 69, Round: 5}:  false,
		{Sequence: 69, Round: 6}:  false,
		{Sequence: 69, Round: 10}: false,
		{Sequence: 69, Round: 12}: false,
		{Sequence: 69, Round: 14}: false,
		{Sequence: 70, Round: 0}:  false,
		{Sequence: 70, Round: 4}:  false,
		{Sequence: 70, Round: 12}: false,
		{Sequence: 71, Round: 0}:  false,
		{Sequence: 71, Round: 5}:  false,
		{Sequence: 71, Round: 10}: false,
		{Sequence: 72, Round: 2}:  false,
		{Sequence: 72, Round: 8}:  false,
		{Sequence: 72, Round: 9}:  false,
		{Sequence: 73, Round: 5}:  false,
		{Sequence: 73, Round: 8}:  false,
		{Sequence: 73, Round: 10}: false,
		{Sequence: 73, Round: 12}: false,
		{Sequence: 73, Round: 14}: false,
		{Sequence: 74, Round: 6}:  false,
		{Sequence: 74, Round: 10}: false,
		{Sequence: 74, Round: 12}: false,
		{Sequence: 75, Round: 2}:  false,
		{Sequence: 75, Round: 5}:  false,
		{Sequence: 75, Round: 6}:  false,
		{Sequence: 75, Round: 7}:  false,
		{Sequence: 76, Round: 7}:  false,
		{Sequence: 77, Round: 0}:  false,
		{Sequence: 77, Round: 7}:  false,
		{Sequence: 78, Round: 0}:  false,
		{Sequence: 78, Round: 2}:  false,
		{Sequence: 78, Round: 5}:  false,
		{Sequence: 79, Round: 0}:  false,
		{Sequence: 79, Round: 4}:  false,
		{Sequence: 79, Round: 11}: false,
		{Sequence: 79, Round: 12}: false,
		{Sequence: 80, Round: 2}:  false,
		{Sequence: 80, Round: 4}:  false,
		{Sequence: 80, Round: 5}:  false,
		{Sequence: 80, Round: 7}:  false,
		{Sequence: 80, Round: 8}:  false,
		{Sequence: 80, Round: 10}: false,
		{Sequence: 80, Round: 14}: false,
		{Sequence: 81, Round: 1}:  false,
		{Sequence: 81, Round: 9}:  false,
		{Sequence: 81, Round: 11}: false,
		{Sequence: 81, Round: 14}: false,
		{Sequence: 82, Round: 0}:  false,
		{Sequence: 82, Round: 11}: false,
		{Sequence: 82, Round: 13}: false,
		{Sequence: 82, Round: 14}: false,
		{Sequence: 83, Round: 0}:  false,
		{Sequence: 83, Round: 5}:  false,
		{Sequence: 83, Round: 6}:  false,
		{Sequence: 83, Round: 8}:  false,
		{Sequence: 83, Round: 9}:  false,
		{Sequence: 83, Round: 12}: false,
		{Sequence: 84, Round: 2}:  false,
		{Sequence: 84, Round: 11}: false,
		{Sequence: 85, Round: 4}:  false,
		{Sequence: 85, Round: 7}:  false,
		{Sequence: 85, Round: 8}:  false,
		{Sequence: 86, Round: 5}:  false,
		{Sequence: 86, Round: 9}:  false,
		{Sequence: 87, Round: 1}:  false,
		{Sequence: 87, Round: 5}:  false,
		{Sequence: 87, Round: 6}:  false,
		{Sequence: 87, Round: 7}:  false,
		{Sequence: 87, Round: 9}:  false,
		{Sequence: 87, Round: 10}: false,
		{Sequence: 87, Round: 12}: false,
		{Sequence: 87, Round: 14}: false,
		{Sequence: 88, Round: 8}:  false,
		{Sequence: 89, Round: 0}:  false,
		{Sequence: 89, Round: 7}:  false,
		{Sequence: 90, Round: 3}:  false,
		{Sequence: 90, Round: 4}:  false,
		{Sequence: 90, Round: 9}:  false,
		{Sequence: 90, Round: 10}: false,
		{Sequence: 90, Round: 11}: false,
		{Sequence: 91, Round: 10}: false,
		{Sequence: 91, Round: 12}: false,
		{Sequence: 91, Round: 13}: false,
		{Sequence: 92, Round: 0}:  false,
		{Sequence: 92, Round: 1}:  false,
		{Sequence: 92, Round: 2}:  false,
		{Sequence: 92, Round: 5}:  false,
		{Sequence: 92, Round: 10}: false,
		{Sequence: 92, Round: 14}: false,
		{Sequence: 93, Round: 0}:  false,
		{Sequence: 93, Round: 4}:  false,
		{Sequence: 93, Round: 5}:  false,
		{Sequence: 93, Round: 8}:  false,
		{Sequence: 93, Round: 10}: false,
		{Sequence: 93, Round: 12}: false,
		{Sequence: 93, Round: 14}: false,
		{Sequence: 94, Round: 2}:  false,
		{Sequence: 94, Round: 6}:  false,
		{Sequence: 94, Round: 7}:  false,
		{Sequence: 94, Round: 10}: false,
		{Sequence: 95, Round: 8}:  false,
		{Sequence: 95, Round: 9}:  false,
		{Sequence: 95, Round: 10}: false,
		{Sequence: 95, Round: 13}: false,
		{Sequence: 96, Round: 1}:  false,
		{Sequence: 96, Round: 7}:  false,
		{Sequence: 96, Round: 8}:  false,
		{Sequence: 96, Round: 10}: false,
		{Sequence: 96, Round: 12}: false,
		{Sequence: 96, Round: 14}: false,
		{Sequence: 97, Round: 4}:  false,
		{Sequence: 97, Round: 5}:  false,
		{Sequence: 97, Round: 13}: false,
		{Sequence: 98, Round: 10}: false,
		{Sequence: 98, Round: 12}: false,
		{Sequence: 99, Round: 4}:  false,
		{Sequence: 99, Round: 14}: false,
	}
)

type keys []*ecdsa.PrivateKey

func (slice keys) Len() int {
	return len(slice)
}

func (slice keys) Less(i, j int) bool {
	return strings.Compare(crypto.PubkeyToAddress(slice[i].PublicKey).String(), crypto.PubkeyToAddress(slice[j].PublicKey).String()) < 0
}

func (slice keys) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

type Pair struct {
	Sequence int64
	Round    int64
}

func getTestCouncil() []common.Address {
	return []common.Address{
		common.HexToAddress("0x414790CA82C14A8B975cEBd66098c3dA590bf969"), // Node Address for test
		common.HexToAddress("0x604973C51f6389dF2782E018000c3AC1257dee90"),
		common.HexToAddress("0x5Ac1689ae5F521B05145C5Cd15a3E8F6ab39Af19"),
		common.HexToAddress("0x0688CaC68bbF7c1a0faedA109c668a868BEd855E"),
		common.HexToAddress("0xaD227Fd4d8a6f464Fb5A8bcf38533337A02Db4e0"),
		common.HexToAddress("0xf2E89C8e9B4C903c046Bb183a3175405fE98A1Db"),
		common.HexToAddress("0xc5D9E04E58717A7Dc4757DF98B865A1525187060"),
		common.HexToAddress("0x99F215e780e352647A1f83E17Eb91930aDbaf3e2"),
		common.HexToAddress("0xf4eBE96E668c5a372C6ad2924C2B9892817c1b61"),
		common.HexToAddress("0x0351D00Cf34b6c2891bf91B6130651903BBdE7df"),
		common.HexToAddress("0xbc0182AA0516666cec5c8185311d80d02b2Bb1F5"),
		common.HexToAddress("0xf8cbeF7eDf33c1437d17aeAC1b2AfC8430a5eAd7"),
		common.HexToAddress("0xbEC6E6457aAE091FC905122DeFf1f97532395896"),
		common.HexToAddress("0x3c1587F672Cf9C457FC229ccA9Cc0c8e29af88BE"),
		common.HexToAddress("0xcE9E07d403Cf3fC3fa6E2F3e0f88e0272Af42EF3"),
		common.HexToAddress("0x711EAd71f23e6b84BD536021Bd02B3706B81A3bE"),
		common.HexToAddress("0xb84A4985CD5B5b1CE56527b89b3e447820633856"),
		common.HexToAddress("0x739E64B2a55626c921114ee8B49cD1b7E5E2bBB0"),
		common.HexToAddress("0x0312b2B142855986C5937eb33E0993ecA142Caca"),
		common.HexToAddress("0xd3AbD32504eA87409Bbb5190059c077C6a9df879"),
		common.HexToAddress("0x0978D638BAc5990c64E52A69C6F33222F16117Ee"),
		common.HexToAddress("0x131855B4D54E9AE376E568dd7c9d925AA6eE0545"),
		common.HexToAddress("0x0c7f24972D43B1F6dD1286275EC089755b15514D"),
		common.HexToAddress("0xc46C39B333C0828820087E7950Ad9F7b30E38373"),
		common.HexToAddress("0x5b3c0E461409a671B15fD1D284F87d2CD1994386"),
		common.HexToAddress("0xB5EA5afCC5045690Afd2afb62748227e16872efA"),
		common.HexToAddress("0x598cC8b6026d666681598ea1d0C85D2b06876277"),
		common.HexToAddress("0xD86D0dB9f101600A4ED8043772654C9D076c0616"),
		common.HexToAddress("0x82d718F86D04454cF5736215110d0a5d9bEe420D"),
		common.HexToAddress("0xE0F03B6915B9e900B258384106644CDab6aAfc24"),
		common.HexToAddress("0xCd28661B14DCda6e002408a0D3A35C1448be7f23"),
		common.HexToAddress("0x77Fd4877Dda3641588E5362744fc6A5d4AE3b731"),
		common.HexToAddress("0xD7E2FEFE5c1C33e91deFf3cAeC3adC8bD3b6afB8"),
		common.HexToAddress("0x7A15B06DDd3ff274F686F9f4a69975BF27aBe37b"),
		common.HexToAddress("0x34083A2D7b252EA7c4EC1C92A2D9D6F2bE6B347e"),
		common.HexToAddress("0x7C73ee31E9b79aEb26B634E46202C46B13160Eba"),
		common.HexToAddress("0x16F7D538Afd579068B88384bCa3c3aefb66C0E52"),
		common.HexToAddress("0x7584dbFad6664F604C9C0dE3c822961A10E064f9"),
		common.HexToAddress("0xE6c82318Da0819880137cCAE993e73d2bC1b8b20"),
		common.HexToAddress("0x0699d025c98ce2CB3C96050efFB1450fE64f5C9E"),
		common.HexToAddress("0x9b6Bdb3a669A721b3FeE537B5B5914B8D9d7F980"),
		common.HexToAddress("0x4CbC62FF9893df8acF8982f0D3948Ec54F8a1d6c"),
		common.HexToAddress("0xBE993372712Cb7ff05c2fc7D21F36a75e8033224"),
		common.HexToAddress("0x2C4CF83c05B7127a714E1749B3079B016E1a7B8f"),
		common.HexToAddress("0xEDE8613a71A914FD3AFC6E219f798f904ffB63e5"),
		common.HexToAddress("0x5fB8bc27982B7122ae6Fd4D4ddbA3E69779422B3"),
		common.HexToAddress("0x1b95B986EeDa22e31f10c6499Cbc236386Ac6817"),
		common.HexToAddress("0xC6B027F0d09348020Bd5e3E6ed692f2F145c6D73"),
		common.HexToAddress("0x06Fc6F3032C03f90eA5C0bE7B95AB8973074D9f4"),
		common.HexToAddress("0x183ed276Ef12bA138c96B51D86619a5A8De82b3e"),
		common.HexToAddress("0xb2AaC4685Bbf98a114f3949Db456545C62AF096c"),
		common.HexToAddress("0xD9c041861214F1D7484e924F87c30B7a5e0DA462"),
		common.HexToAddress("0xa3255A75799De01f3b93Bf0D9DF34F5d767CeDE0"),
		common.HexToAddress("0xFC112865D09332c2122c5b29bb6FccA171A6c41c"),
		common.HexToAddress("0xFA58289602bfE7B029b09D74046a36fA80a41F71"),
		common.HexToAddress("0x7d60D7c9ae172d5C8A61eDc37e55F6A60E175f14"),
		common.HexToAddress("0x2A7D561FFAA1fD6f4Db0eC0eB20799cfdd9AfA37"),
		common.HexToAddress("0xaAa5133219d8fdB6Cad565CBEBc9113b818C64b5"),
		common.HexToAddress("0x219bFfaB40F95673F3a4eb243C23E0803a343d9E"),
		common.HexToAddress("0x6A53aA17DBaAE11a428576501D72fe346cB7B1f7"),
		common.HexToAddress("0x4225BE9Eae8309FdE653CFd3d6A8146Bd0389d3b"),
		common.HexToAddress("0x63F7e13f7e586480D2AD4901311A66aecfd714F4"),
		common.HexToAddress("0xAb571920e680A1350FAb4a8eecF3cf8451FBBD4D"),
		common.HexToAddress("0xB12cD0799C6caf355A7cccB188E1C6d5FA6e789c"),
		common.HexToAddress("0xFEf0372eb16b77B67C4AC39e8aF045624392F53c"),
		common.HexToAddress("0x9Ee226196605D536feed2e3Ac4c87a695c394200"),
		common.HexToAddress("0x1e6E55D4485853Bb5b42Ec591e2900ACd9C8a7eA"),
		common.HexToAddress("0xdf0C7B415738E5167E5E43563f11FFDF46696FB0"),
		common.HexToAddress("0x718a3B4b201849E0B2905EC0A29c5c9323017FC3"),
		common.HexToAddress("0xa906b367D2B0E8426C00B622488F6495b333a5c3"),
		common.HexToAddress("0x4dBDAb308824dF54225F5a8B11c890B386b39C2C"),
		common.HexToAddress("0xA0D13983Daa2604e66F692Cf695d5D04b39958c4"),
		common.HexToAddress("0x5bae720C342157325B1D710766326766a41D50F2"),
		common.HexToAddress("0x4d1890BdB54dde6656E89e9b1447672759AF18aB"),
		common.HexToAddress("0x4e2cFB15010A576B6c2c70C1701E2d2dfF4FB2A7"),
		common.HexToAddress("0x8858B61E2A724aEc8542f1B53aE7C1b08Bf823c8"),
		common.HexToAddress("0x8849E54A211B6d8D096C150DF33795d6B1cF6ba9"),
		common.HexToAddress("0x14193546e619761795973b9De768753c61C5f9CB"),
		common.HexToAddress("0xdFa4b5b4f241F4F7D254A840F648B81ca2338F82"),
		common.HexToAddress("0xB51aF0332993Dc3cB7211AD42C761F736A9Af82a"),
		common.HexToAddress("0x69660E66352f91D19Ac12CEc20ee27aC4346AC3F"),
		common.HexToAddress("0xdB0384Ece61E79e99F7944f4651a0b129c383279"),
		common.HexToAddress("0x65a13961017181Ba5Be64959527242251aBB21B9"),
		common.HexToAddress("0xc96a196e136403f530C9a6642C3181530124cb59"),
		common.HexToAddress("0x4788f45D1AF4508A4c9a26dAB839556466f1481b"),
		common.HexToAddress("0xc53e50937E481364b6b6926F572ea65b3B790cDA"),
		common.HexToAddress("0xF94270fD8a0393202233D5e0163a41cb0A272DEe"),
		common.HexToAddress("0xEF2892748176A6345D7EC22924C4A5f6ec563ccc"),
		common.HexToAddress("0x6A01Eba6729F2Fa4A570ff94FA4Cf51Dde428c27"),
		common.HexToAddress("0x0CD9337A1C5B744B9D99db496706d916c76a3B27"),
		common.HexToAddress("0xffE19f7e59eB3bc0a29907D3E8d5397703CAe605"),
		common.HexToAddress("0xE74dCd57694FE43D1813595a43AC8c54DEd62Fed"),
		common.HexToAddress("0x7E07f18eD745cD216341147c990088f445179a1c"),
		common.HexToAddress("0x0eC5d785C5b47C70154cCd93D71cCCB624b00b29"),
		common.HexToAddress("0x6D576beB3ec319Ab366aC7455d3D096e9d2De67d"),
		common.HexToAddress("0x8D8083125960CB0F4e638dB1a0ee855d438a707E"),
		common.HexToAddress("0xCF9b3b510f11285356167396BABe9A4063DAa1dD"),
		common.HexToAddress("0x38b6bAb66D8CD7035f4f90d87aF7d41d7958bed7"),
		common.HexToAddress("0x8c597E4b0A71571C209293c278099B6D85Cd3290"),
		common.HexToAddress("0x39Cdf8a09f5c516951E0FBa684086702912a4810"),
	}
}

func getTestRewards() []common.Address {
	return []common.Address{
		common.HexToAddress("0x2A35FE72F847aa0B509e4055883aE90c87558AaD"),
		common.HexToAddress("0xF91B8EBa583C7fa603B400BE17fBaB7629568A4a"),
		common.HexToAddress("0x240ed27c8bDc9Bb6cA08fa3D239699Fba525d05a"),
		common.HexToAddress("0x3B980293396Fb0e827929D573e3e42d2EA902502"),
		common.HexToAddress("0x11F3B5a36DBc8d8D289CE41894253C0D513cf777"),
		common.HexToAddress("0x7eDD28614052D42430F326A78956859a462f3Cd1"),
		common.HexToAddress("0x3e473733b16D88B695fc9E2278Ab854F449Ea017"),
		common.HexToAddress("0x873B6D0404b46Ddde4DcfC34380C14f54002aE27"),
		common.HexToAddress("0xC85719Ffb158e1Db55D6704B6AE193aa00f5341e"),
		common.HexToAddress("0x957A2DB51Ba8BADA57218FA8580a23a9c64B618f"),
		common.HexToAddress("0x5195eC4dC0048032a7dA40B861D47d560Bfd9462"),
		common.HexToAddress("0x11e3219418e6214c12E3A6271d1B56a68526fB99"),
		common.HexToAddress("0x2edB0c9072497A2276248B0ae11D560114b273fC"),
		common.HexToAddress("0x26e059665D30Bbe015B3b131175F53903cAa2Df9"),
		common.HexToAddress("0x70528eC69e4e710F71c22aEdA6315Fe1750298e3"),
		common.HexToAddress("0xFb538Db8C0b71236bF2B7091AdD0576879D2B8D4"),
		common.HexToAddress("0x705Ec0595A40D1A3B77360CE9cA138b731b193ac"),
		common.HexToAddress("0x55Cede11Bfe34cdB844c87979a35a7369B36Af2D"),
		common.HexToAddress("0xCe79d52e0215872C0515767E46D3C29EC6C858B7"),
		common.HexToAddress("0xf83cedb170F84517A693e0A111F99Fe5C4FA50B2"),
		common.HexToAddress("0xCF4ff5446B6e1Ee4151B5E455B5a6f80a2296AdB"),
		common.HexToAddress("0xB5A7a3bFA0FA44a9CcBED4dc7777013Dfc4D43bc"),
		common.HexToAddress("0xcBa042B2B44512FbB31Be4a478252Fcc78Bd94fB"),
		common.HexToAddress("0xed7a54066A503A93Cd04F572E39e0F02eE0c92af"),
		common.HexToAddress("0x1630a3F6fe4Fc3a4eb2212d12306AEF856F874Df"),
		common.HexToAddress("0x8fb35e1bB2Cf26355b7906b3747db941Fc48bf0C"),
		common.HexToAddress("0x211F9Bf01C553644B22AbA5ca83C2B8F9A07E49B"),
		common.HexToAddress("0xf08Cb412aeb59de13a5F56Df43BA75EaA3B2F0dE"),
		common.HexToAddress("0xc57047791C184Baf0184bd1e0a5e012532b40A45"),
		common.HexToAddress("0x6660E78cE36d96d81e447D51D83ce3930b54effC"),
		common.HexToAddress("0x0736f11d0c066C6e76CE40d89AfCa86FA0Ce9AB5"),
		common.HexToAddress("0x6881515A8499f871E95a4EC17689a37afB5Bc600"),
		common.HexToAddress("0x96b8473f80DB5Fa8056CBb1bCECAe1f25C35733D"),
		common.HexToAddress("0xf9D4D5286bbDeffbcf34d342c47b3E6e52Bfcc08"),
		common.HexToAddress("0x272a690C3B3136d4d431800486233EECBb2092C5"),
		common.HexToAddress("0x98760b1954017DCDAaf466C03d85954E4bF2E431"),
		common.HexToAddress("0xbfFFE194Ebf9953C92a206Fb9c71612afBA80D53"),
		common.HexToAddress("0xF6dF311801f6B9ec7ce7e5aba31CDf3dC02de133"),
		common.HexToAddress("0x675D19896DA6bca2E236F7C3229E1d82993dDB69"),
		common.HexToAddress("0x426d004E495DdE95AdCeE5912cdbf90F2692667b"),
		common.HexToAddress("0xe4ed55aB59335341A7dfc153f8068984d53bbaa8"),
		common.HexToAddress("0xf7399b92de54BC059E957CE9c1baE294aEF1DA08"),
		common.HexToAddress("0x41259376048AF3FBf302b8E588fF8003e3b44F06"),
		common.HexToAddress("0xb007291DDE0Ca6BDd69aB16902bea67819abc515"),
		common.HexToAddress("0x9A04d8d78fd8c536563cdA344b5A04DbF27D67D3"),
		common.HexToAddress("0x12465F01b52CAEC4c3C0197308F6dF8E54163d9c"),
		common.HexToAddress("0x323cf51E7A511De099d5E39dC3ed8B8DDdB3F038"),
		common.HexToAddress("0x36e3E1E34d456975e8C6aFbD7c6F5Fc0F98E0887"),
		common.HexToAddress("0x1A1ed8ae1dC926Be272888ed56EA13503fcE5680"),
		common.HexToAddress("0x8D5D4eD1678F77b7337BA8042fb8C9E7361cbFF5"),
		common.HexToAddress("0xbBbfc9358f931Cd105107d93b5DF97159C933096"),
		common.HexToAddress("0xA82FacF50E9722d6fFa30e186A01625CCE03aA5a"),
		common.HexToAddress("0xA2b73E026CAFB45af42705f994e884fa1D7c25Ad"),
		common.HexToAddress("0x3dEE1382719E04505Df63df4A604CCB705FaC016"),
		common.HexToAddress("0xDf41a247f7F7Ed4F191932844F1C1Dd23D9a76fB"),
		common.HexToAddress("0x95debC7f1a4B30F1451c95178aCaC6e72539C633"),
		common.HexToAddress("0x505594ACc508D855eAAF68CDd516011D53F76a54"),
		common.HexToAddress("0xfFE3160f6B5c73952D2A8f55022f827854cdB8C5"),
		common.HexToAddress("0x76BaE5bdE095a4a5F256DEa65259d24d454B777d"),
		common.HexToAddress("0xfF34f95Cf7815268Bf3e62AbFBc5452606bdB336"),
		common.HexToAddress("0x677799E7804c930Dd430470B9F0C6Cd6523Dd487"),
		common.HexToAddress("0x7b3D56CFfe4CB5F5EDbDd4b7a41b5eBa848e9348"),
		common.HexToAddress("0x9e80C52fC892C916616Fc7235B885047e1165Dfe"),
		common.HexToAddress("0x5E100F1FFe080CFD7dCBAEe2Fa1DbE4D81327955"),
		common.HexToAddress("0x0C4Cf38D7769243f843DE339f866bB6f770Bc7DD"),
		common.HexToAddress("0xf23A77988e0592CDa049fe7aBb46E900a704eE90"),
		common.HexToAddress("0xC67a2ff05C5CA95D76525f5E471701813A877fB2"),
		common.HexToAddress("0x8bA5d606bEBe3CC4F4C001B68764771584cD9E4f"),
		common.HexToAddress("0x82A39CddA3e64b2EcE2E161116AfCc7Bc5aABC7c"),
		common.HexToAddress("0x0324216F0bBA982B01B38A0A354783324B1A6a00"),
		common.HexToAddress("0x229CdA4Dfc2ED2503aBCfdca39aA1D4d665281a6"),
		common.HexToAddress("0x53C65264994dccEFBe1C3C6b4c2C9fCf4fc3458a"),
		common.HexToAddress("0x822BB8E5a0650740424B2bBbd3BcdaD2B8e4FB96"),
		common.HexToAddress("0xaE939D5C93fB3d8522738bC5D84a69b7d9ec625F"),
		common.HexToAddress("0xF32Dc9C078028A2C8a4104ea70615ADB3010B2b0"),
		common.HexToAddress("0x2c231E031e9803e35B97578B536eE3739EBa886F"),
		common.HexToAddress("0xF6c24F7e7461BA43800E20bdd32c8483eB8aA152"),
		common.HexToAddress("0xB3363359C172bdE5a0747444cbC9108e50FEf826"),
		common.HexToAddress("0xc6a42020AD8fB61fa3A4c248310625bBdd5cd04c"),
		common.HexToAddress("0xbE797039B57007CEa7F4F5961a6AB64dFc076D0c"),
		common.HexToAddress("0xe12F6f36D939ff66CBF35b5F67f3D8670dD4228f"),
		common.HexToAddress("0x4c9d56C6c06b48511FdbFD4Dd7C82D56bcAC88f3"),
		common.HexToAddress("0x1Ab93C25Fd220e8C80ae0dA73bcfC56f969aD8Db"),
		common.HexToAddress("0xF25ccD01c83ee36a0361f60B164108ea1E57FE14"),
		common.HexToAddress("0x6f7C69666a9E6B34835f505C3025bcd95Aae2600"),
		common.HexToAddress("0xc78387B2384Dcef8c7a9b39BD688d2C2776945E2"),
		common.HexToAddress("0x4786F93B4e041eBB7F1EcF9172fE4Cf0c16bD88E"),
		common.HexToAddress("0x1d64Ea74B2EEB37fB24158A5d12Eb5b61c718f44"),
		common.HexToAddress("0xf5E9Fc7Bf47c9b0eE7C8349eBf3F6442fec426d8"),
		common.HexToAddress("0x8aA20E60Da8b86717A9785ec255Aceb8E509107f"),
		common.HexToAddress("0x74cc084275253fabD14a0ee6264047503057Aa88"),
		common.HexToAddress("0xc1AF209ed6fe5ae23ac5ad3f4D189E0E04D925E8"),
		common.HexToAddress("0x4725563C44e432013B20d06395130A0d24ad091F"),
		common.HexToAddress("0x6aFB4f059BA732d5143f0fAfD5534501Cf57A47C"),
		common.HexToAddress("0xCB1c4478cbE37a4B7a1510065802e0B669979cFD"),
		common.HexToAddress("0xf7ed2347390b1f742DA9785242dF90BbbC8Af90C"),
		common.HexToAddress("0x8BC37790112C66DF94cB05dC5168249a7e6d6717"),
		common.HexToAddress("0x9087C8dd3dE26B6DB3a471D7254d78FB1afCbeF2"),
		common.HexToAddress("0xcAfaA922f46d0a0FcE13B0f2aA822aF9C4b9a31C"),
		common.HexToAddress("0xa21b5d69efDFE56a9C272f7957Bd6d4205a1e6Ff"),
	}
}

func getTestVotingPowers(num int) []uint64 {
	vps := make([]uint64, 0, num)
	for i := 0; i < num; i++ {
		vps = append(vps, 1)
	}
	return vps
}

func getTestConfig() *params.ChainConfig {
	config := params.TestChainConfig
	config.Governance = params.GetDefaultGovernanceConfig(params.UseIstanbul)
	config.Istanbul = params.GetDefaultIstanbulConfig()
	return config
}

func Benchmark_getTargetReceivers(b *testing.B) {
	_, backend := newBlockChain(1)
	defer backend.Stop()

	backend.currentView.Store(&istanbul.View{Sequence: big.NewInt(0), Round: big.NewInt(0)})

	// Create ValidatorSet
	council := getTestCouncil()
	rewards := getTestRewards()
	valSet := validator.NewWeightedCouncil(council, nil, rewards, getTestVotingPowers(len(council)), nil, istanbul.WeightedRandom, 21, 0, 0, nil)
	valSet.SetBlockNum(uint64(1))
	valSet.CalcProposer(valSet.GetProposer().Address(), uint64(1))
	hex := fmt.Sprintf("%015d000000000000000000000000000000000000000000000000000", 1)
	prevHash := common.HexToHash(hex)

	for i := 0; i < b.N; i++ {
		_ = backend.getTargetReceivers(prevHash, valSet)
	}
}

// Test_GossipSubPeerTargets checks if the gossiping targets are same as council members
func Test_GossipSubPeerTargets(t *testing.T) {
	// get testing node's address
	key, _ := crypto.HexToECDSA(PRIVKEY) // This key is to be provided to create backend

	_, backend := newBlockChain(1, istanbulCompatibleBlock(big.NewInt(5)), key)
	defer backend.Stop()

	// Create ValidatorSet
	council := getTestCouncil()
	rewards := getTestRewards()
	valSet := validator.NewWeightedCouncil(council, nil, rewards, getTestVotingPowers(len(council)), nil, istanbul.WeightedRandom, 21, 0, 0, nil)
	valSet.SetBlockNum(uint64(5))

	// Test for blocks from 0 to maxBlockNum
	// from 0 to 4: before istanbul hard fork
	// from 5 to 100: after istanbul hard fork
	for i := int64(0); i < maxBlockNum; i++ {
		// Test for round 0 to round 14
		for round := int64(0); round < 15; round++ {
			backend.currentView.Store(&istanbul.View{Sequence: big.NewInt(i), Round: big.NewInt(round)})
			valSet.SetBlockNum(uint64(i))
			valSet.CalcProposer(valSet.GetProposer().Address(), uint64(round))

			// Use block number as prevHash. In SubList() only left 15 bytes are being used.
			hex := fmt.Sprintf("%015d000000000000000000000000000000000000000000000000000", i)
			prevHash := common.HexToHash(hex)

			// committees[0]: current committee
			// committees[1]: next committee
			committees := make([][]istanbul.Validator, 2)

			// Getting the current round's committee
			viewCurrent := backend.currentView.Load().(*istanbul.View)
			committees[0] = valSet.SubList(prevHash, viewCurrent)

			// Getting the next round's committee
			viewCurrent.Round = viewCurrent.Round.Add(viewCurrent.Round, common.Big1)
			backend.currentView.Store(viewCurrent)

			valSet.CalcProposer(valSet.GetProposer().Address(), uint64(round+1))
			committees[1] = valSet.SubList(prevHash, viewCurrent)

			// Reduce round by 1 to set round to the current round before calling GossipSubPeer
			viewCurrent.Round = viewCurrent.Round.Sub(viewCurrent.Round, common.Big1)
			valSet.CalcProposer(valSet.GetProposer().Address(), uint64(round))
			backend.currentView.Store(viewCurrent)

			// Receiving the receiver list of a message
			targets := backend.GossipSubPeer(prevHash, valSet, nil)

			// Check if the testing node is in a committee
			isInSubList := backend.checkInSubList(prevHash, valSet)
			isInCommitteeBlocks := checkInCommitteeBlocks(i, round)

			// Check if the result of checkInSubList is same as expected. It is to detect an unexpected change in SubList logic
			if isInSubList != isInCommitteeBlocks {
				t.Errorf("Difference in expected data and calculated one. Changed committee selection? HARD FORK may happen!! Sequence: %d, Round: %d", i, round)
			} else {
				if isInSubList == false {
					continue
				}
			}

			// number of message receivers have to be smaller than or equal to the number of the current committee and the next committee
			if len(targets) > len(committees[0])+len(committees[1]) {
				t.Errorf("Target has too many validators. targets: %d, sum of committees: %d", len(targets), len(committees[0])+len(committees[1]))
			}

			// Check all nodes in the current and the next round are included in the target list
			for n := 0; n < len(committees); n++ {
				for _, x := range committees[n] {
					if _, ok := targets[x.Address()]; !ok && x.Address() != backend.Address() {
						t.Errorf("Block: %d, Round: %d, Committee member %v not found in targets", i, round, x.Address().String())
					} else {
						// Mark the target is in the current or in the next committee
						targets[x.Address()] = false
					}
				}
			}

			// Check if a validator not in the current/next committee is included in target list
			for k, v := range targets {
				if v == true {
					t.Errorf("Block: %d, Round: %d, Validator not in committees included %v", i, round, k.String())
				}
			}
		}
	}
	// Check if the testing node is in all committees that it is supposed to be
	for k, v := range committeeBlocks {
		if !v {
			fmt.Printf("The node is missing in committee that it should be included in. Sequence %d, Round %d\n", k.Sequence, k.Round)
		}
	}
}

func checkInCommitteeBlocks(seq int64, round int64) bool {
	v := Pair{seq, round}
	if _, ok := committeeBlocks[v]; ok {
		committeeBlocks[v] = true
		return true
	}
	return false
}

func newTestBackend() (b *backend) {
	config := getTestConfig()
	config.Istanbul.ProposerPolicy = params.WeightedRandom
	return newTestBackendWithConfig(config, istanbul.DefaultConfig.BlockPeriod, nil)
}

func newTestBackendWithConfig(chainConfig *params.ChainConfig, blockPeriod uint64, key *ecdsa.PrivateKey) (b *backend) {
	dbm := database.NewDBManager(&database.DBConfig{DBType: database.MemoryDB})
	if key == nil {
		// if key is nil, generate new key for a test account
		key, _ = crypto.GenerateKey()
	}
	if chainConfig.Governance.GovernanceMode == "single" {
		// if governance mode is single, set the node key to the governing node.
		chainConfig.Governance.GoverningNode = crypto.PubkeyToAddress(key.PublicKey)
	}
	gov := governance.NewMixedEngine(chainConfig, dbm)
	istanbulConfig := istanbul.DefaultConfig
	istanbulConfig.BlockPeriod = blockPeriod
	istanbulConfig.ProposerPolicy = istanbul.ProposerPolicy(chainConfig.Istanbul.ProposerPolicy)
	istanbulConfig.Epoch = chainConfig.Istanbul.Epoch
	istanbulConfig.SubGroupSize = chainConfig.Istanbul.SubGroupSize

	backend := New(getTestRewards()[0], istanbulConfig, key, dbm, gov, common.CONSENSUSNODE).(*backend)
	gov.SetNodeAddress(crypto.PubkeyToAddress(key.PublicKey))
	return backend
}

func newTestValidatorSet(n int, policy istanbul.ProposerPolicy) (istanbul.ValidatorSet, []*ecdsa.PrivateKey) {
	// generate validators
	keys := make(keys, n)
	addrs := make([]common.Address, n)
	for i := 0; i < n; i++ {
		privateKey, _ := crypto.GenerateKey()
		keys[i] = privateKey
		addrs[i] = crypto.PubkeyToAddress(privateKey.PublicKey)
	}
	vset := validator.NewSet(addrs, policy)
	sort.Sort(keys) //Keys need to be sorted by its public key address
	return vset, keys
}

func TestSign(t *testing.T) {
	b := newTestBackend()

	sig, err := b.Sign(testSigningData)
	if err != nil {
		t.Errorf("error mismatch: have %v, want nil", err)
	}

	//Check signature recover
	hashData := crypto.Keccak256([]byte(testSigningData))
	pubkey, _ := crypto.Ecrecover(hashData, sig)
	actualSigner := common.BytesToAddress(crypto.Keccak256(pubkey[1:])[12:])

	if actualSigner != b.address {
		t.Errorf("address mismatch: have %v, want %s", actualSigner.Hex(), b.address.Hex())
	}
}

func TestCheckSignature(t *testing.T) {
	b := newTestBackend()

	// testAddr is derived from testPrivateKey.
	testPrivateKey, _ := crypto.HexToECDSA("bb047e5940b6d83354d9432db7c449ac8fca2248008aaa7271369880f9f11cc1")
	testAddr := common.HexToAddress("0x70524d664ffe731100208a0154e556f9bb679ae6")
	testInvalidAddr := common.HexToAddress("0x9535b2e7faaba5288511d89341d94a38063a349b")

	hashData := crypto.Keccak256([]byte(testSigningData))
	sig, err := crypto.Sign(hashData, testPrivateKey)
	if err != nil {
		t.Fatalf("unexpected failure: %v", err)
	}

	if err := b.CheckSignature(testSigningData, testAddr, sig); err != nil {
		t.Errorf("error mismatch: have %v, want nil", err)
	}

	if err := b.CheckSignature(testSigningData, testInvalidAddr, sig); err != errInvalidSignature {
		t.Errorf("error mismatch: have %v, want %v", err, errInvalidSignature)
	}
}

func TestCheckValidatorSignature(t *testing.T) {
	vset, keys := newTestValidatorSet(5, istanbul.WeightedRandom)

	// 1. Positive test: sign with validator's key should succeed
	hashData := crypto.Keccak256([]byte(testSigningData))
	for i, k := range keys {
		// Sign
		sig, err := crypto.Sign(hashData, k)
		if err != nil {
			t.Errorf("error mismatch: have %v, want nil", err)
		}
		// CheckValidatorSignature should succeed
		addr, err := istanbul.CheckValidatorSignature(vset, testSigningData, sig)
		if err != nil {
			t.Errorf("error mismatch: have %v, want nil", err)
		}
		validator := vset.GetByIndex(uint64(i))
		if addr != validator.Address() {
			t.Errorf("validator address mismatch: have %v, want %v", addr, validator.Address())
		}
	}

	// 2. Negative test: sign with any key other than validator's key should return error
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Errorf("error mismatch: have %v, want nil", err)
	}
	// Sign
	sig, err := crypto.Sign(hashData, key)
	if err != nil {
		t.Errorf("error mismatch: have %v, want nil", err)
	}
	// CheckValidatorSignature should return ErrUnauthorizedAddress
	addr, err := istanbul.CheckValidatorSignature(vset, testSigningData, sig)
	if err != istanbul.ErrUnauthorizedAddress {
		t.Errorf("error mismatch: have %v, want %v", err, istanbul.ErrUnauthorizedAddress)
	}
	emptyAddr := common.Address{}
	if addr != emptyAddr {
		t.Errorf("address mismatch: have %v, want %v", addr, emptyAddr)
	}
}

func TestCommit(t *testing.T) {
	backend := newTestBackend()

	commitCh := make(chan *types.Block)
	// Case: it's a proposer, so the backend.commit will receive channel result from backend.Commit function
	testCases := []struct {
		expectedErr       error
		expectedSignature [][]byte
		expectedBlock     func() *types.Block
	}{
		{
			// normal case
			nil,
			[][]byte{append([]byte{1}, bytes.Repeat([]byte{0x00}, types.IstanbulExtraSeal-1)...)},
			func() *types.Block {
				chain, engine := newBlockChain(1)
				defer engine.Stop()

				block := makeBlockWithoutSeal(chain, engine, chain.Genesis())
				expectedBlock, _ := engine.updateBlock(engine.chain.GetHeader(block.ParentHash(), block.NumberU64()-1), block)
				return expectedBlock
			},
		},
		{
			// invalid signature
			errInvalidCommittedSeals,
			nil,
			func() *types.Block {
				chain, engine := newBlockChain(1)
				defer engine.Stop()

				block := makeBlockWithoutSeal(chain, engine, chain.Genesis())
				expectedBlock, _ := engine.updateBlock(engine.chain.GetHeader(block.ParentHash(), block.NumberU64()-1), block)
				return expectedBlock
			},
		},
	}

	for _, test := range testCases {
		expBlock := test.expectedBlock()
		go func() {
			select {
			case result := <-backend.commitCh:
				commitCh <- result.Block
				return
			}
		}()

		backend.proposedBlockHash = expBlock.Hash()
		if err := backend.Commit(expBlock, test.expectedSignature); err != nil {
			if err != test.expectedErr {
				t.Errorf("error mismatch: have %v, want %v", err, test.expectedErr)
			}
		}

		if test.expectedErr == nil {
			// to avoid race condition is occurred by goroutine
			select {
			case result := <-commitCh:
				if result.Hash() != expBlock.Hash() {
					t.Errorf("hash mismatch: have %v, want %v", result.Hash(), expBlock.Hash())
				}
			case <-time.After(10 * time.Second):
				t.Fatal("timeout")
			}
		}
	}
}

func TestGetProposer(t *testing.T) {
	chain, engine := newBlockChain(1)
	defer engine.Stop()

	block := makeBlock(chain, engine, chain.Genesis())
	_, err := chain.InsertChain(types.Blocks{block})
	if err != nil {
		t.Errorf("failed to insert chain: %v", err)
	}
	expected := engine.GetProposer(1)
	actual := engine.Address()
	if actual != expected {
		t.Errorf("proposer mismatch: have %v, want %v", actual.Hex(), expected.Hex())
	}
}
