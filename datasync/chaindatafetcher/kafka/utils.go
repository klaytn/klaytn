// Copyright 2020 The klaytn Authors
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

package kafka

import (
	"strings"

	klaytnApi "github.com/klaytn/klaytn/api"
	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/consensus/istanbul"
	"github.com/klaytn/klaytn/crypto/sha3"
	"github.com/klaytn/klaytn/ser/rlp"
)

func getProposerAndValidatorsFromBlock(block *types.Block) (string, string, error) {
	blockNumber := block.NumberU64()
	if blockNumber == 0 {
		return "", "", nil
	}
	// Retrieve the signature from the header extra-data
	istanbulExtra, err := types.ExtractIstanbulExtra(block.Header())
	if err != nil {
		return "", "", err
	}

	sigHash, err := sigHash(block.Header())
	if err != nil {
		return "", "", err
	}
	proposerAddr, err := istanbul.GetSignatureAddress(sigHash.Bytes(), istanbulExtra.Seal)
	if err != nil {
		return "", "", err
	}
	commiteeAddrs := make([]common.Address, len(istanbulExtra.Validators))
	for i, addr := range istanbulExtra.Validators {
		commiteeAddrs[i] = addr
	}
	var strValidators []string
	for _, validator := range istanbulExtra.Validators {
		strValidators = append(strValidators, validator.Hex())
	}

	return proposerAddr.Hex(), strings.Join(strValidators, ","), nil
}

func sigHash(header *types.Header) (hash common.Hash, err error) {
	hasher := sha3.NewKeccak256()

	// Clean seal is required for calculating proposer seal.
	if err := rlp.Encode(hasher, types.IstanbulFilteredHeader(header, false)); err != nil {
		logger.Error("fail to encode", "err", err)
		return common.Hash{}, err
	}
	hasher.Sum(hash[:0])
	return hash, nil
}

func makeBlockGroupOutput(blockchain *blockchain.BlockChain, block *types.Block, receipts types.Receipts) map[string]interface{} {
	head := block.Header() // copies the header once
	hash := head.Hash()

	proposer, committee, err := getProposerAndValidatorsFromBlock(block)
	if err != nil {
		// skip error handling when getting proposer and committee is failed
		logger.Error("Getting the proposer and validators failed.", "blockHash", hash, "err", err)
	}
	td := blockchain.GetTd(hash, block.NumberU64())
	r, _ := klaytnApi.RpcOutputBlock(block, td, false, false)

	// make transactions
	transactions := block.Transactions()
	numTxs := len(transactions)
	rpcTransactions := make([]map[string]interface{}, numTxs)
	for i, tx := range transactions {
		rpcTransactions[i] = klaytnApi.RpcOutputReceipt(tx, hash, head.Number.Uint64(), uint64(i), receipts[i])
	}

	r["committee"] = committee
	r["proposer"] = proposer
	r["transactions"] = rpcTransactions
	return r
}
