// Copyright 2022 The klaytn Authors
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
	"math/big"
	"os"
	"testing"

	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	rewardcontract "github.com/klaytn/klaytn/contracts/reward/contract"
	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/governance"
	"github.com/klaytn/klaytn/log"
	"github.com/klaytn/klaytn/params"
	"github.com/klaytn/klaytn/reward"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddressBookConnector(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlWarn)

	fullNode, node, validator, chainId, workspace := newBlockchain(t, nil)
	defer os.RemoveAll(workspace)
	defer fullNode.Stop()

	var (
		chain  = node.BlockChain().(*blockchain.BlockChain)
		config = chain.Config()
		txpool = node.TxPool()
		db     = node.ChainDB()
		gov    = governance.NewMixedEngine(config, db)

		deployAddr  common.Address
		deployBlock uint64
	)

	// Deploy AddressBook
	{
		sender := validator
		signer := types.LatestSignerForChainID(chainId)

		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:         sender.GetNonce(),
			types.TxValueKeyAmount:        new(big.Int).SetUint64(0),
			types.TxValueKeyGasLimit:      uint64(1e9),
			types.TxValueKeyGasPrice:      big.NewInt(25 * params.Ston),
			types.TxValueKeyHumanReadable: false,
			types.TxValueKeyFrom:          sender.GetAddr(),
			types.TxValueKeyData:          common.FromHex(rewardcontract.AddressBookBin),
			types.TxValueKeyCodeFormat:    params.CodeFormatEVM,
			types.TxValueKeyTo:            (*common.Address)(nil),
		}

		tx, err := types.NewTransactionWithMap(types.TxTypeSmartContractDeploy, values)
		require.Nil(t, err)

		err = tx.SignWithKeys(signer, sender.GetTxKeys())
		require.Nil(t, err)

		err = txpool.AddLocal(tx)
		require.True(t, err == nil || err == blockchain.ErrAlreadyNonceExistInPool, "err=%v", err)

		deployAddr = crypto.CreateAddress(sender.Addr, sender.Nonce)
		sender.AddNonce()

		receipt := waitReceipt(chain, tx.Hash())
		require.NotNil(t, receipt)
		require.Equal(t, types.ReceiptStatusSuccessful, receipt.Status)

		_, _, deployBlock, _ = chain.GetTxAndLookupInfo(tx.Hash())
		t.Logf("AddressBook deployed at block=%2d", deployBlock)
	}

	// Temporarily use the newly deployed address in StakingManager
	oldAddr := common.HexToAddress(rewardcontract.AddressBookContractAddress)
	reward.SetTestAddressBookAddress(deployAddr)
	defer reward.SetTestAddressBookAddress(oldAddr)

	// Temporarily lower the StakingUpdateInterval
	oldInterval := params.StakingUpdateInterval()
	params.SetStakingUpdateInterval(3)
	defer func() { params.SetStakingUpdateInterval(oldInterval) }()

	// Create the StakingManager singleton
	oldStakingManager := reward.GetStakingManager()
	reward.SetTestStakingManagerWithChain(chain, gov, db)
	defer reward.SetTestStakingManager(oldStakingManager)

	// Attempt to read contract
	require.NotNil(t, waitBlock(chain, deployBlock+3))
	stakingInfo := reward.GetStakingInfo(deployBlock + 6)
	assert.NotNil(t, stakingInfo)

	t.Logf("StakingInfo=%s", stakingInfo)
}
