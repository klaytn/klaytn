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

package tests

import (
	"crypto/ecdsa"
	"io/ioutil"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/klaytn/klaytn/accounts/keystore"
	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/consensus/istanbul"
	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/log"
	"github.com/klaytn/klaytn/networks/p2p"
	"github.com/klaytn/klaytn/node"
	"github.com/klaytn/klaytn/node/cn"
	"github.com/klaytn/klaytn/params"
	"github.com/klaytn/klaytn/rlp"
	"github.com/klaytn/klaytn/work"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSimpleBlockchain
func TestSimpleBlockchain(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlTrace)

	numAccounts := 12
	fullNode, node, validator, chainId, workspace := newBlockchain(t, nil)
	defer os.RemoveAll(workspace)

	// create account
	richAccount, accounts, contractAccounts := createAccount(t, numAccounts, validator)

	contractDeployCode := "0x608060405234801561001057600080fd5b506000808190555060646001819055506101848061002f6000396000f300608060405260043610610062576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff16806302e5329e14610067578063197e70e41461009457806349b667d2146100c157806367e0badb146100ec575b600080fd5b34801561007357600080fd5b5061009260048036038101908080359060200190929190505050610117565b005b3480156100a057600080fd5b506100bf60048036038101908080359060200190929190505050610121565b005b3480156100cd57600080fd5b506100d6610145565b6040518082815260200191505060405180910390f35b3480156100f857600080fd5b5061010161014f565b6040518082815260200191505060405180910390f35b8060018190555050565b806000540160008190555060015460005481151561013b57fe5b0660008190555050565b6000600154905090565b600080549050905600a165627a7a72305820ef4e7e564c744de3a36cb74000c35687f7de9ecf1d29abdd3c4bcc66db981c160029"
	for i := 0; i < numAccounts; i++ {
		_, contractAccounts[i].Addr = deployContractDeployTx(t, node.TxPool(), chainId, richAccount, contractDeployCode)
	}
	time.Sleep(time.Second) // need to make a block before contract execution

	// deploy
	calldata := "0x197e70e40000000000000000000000000000000000000000000000000000000000000001"
	for i := 0; i < numAccounts; i++ {
		deployRandomTxs(t, node.TxPool(), chainId, richAccount, 10)
		deployValueTransferTx(t, node.TxPool(), chainId, richAccount, accounts[i%numAccounts])
		deployContractExecutionTx(t, node.TxPool(), chainId, richAccount, contractAccounts[i%numAccounts].Addr, calldata)

		// time.Sleep(time.Second) // wait until txpool is flushed if needed
	}

	// stop full node
	if err := fullNode.Stop(); err != nil {
		t.Fatal(err)
	}
	time.Sleep(2 * time.Second)

	// start full node with previous db
	fullNode, node, err := newKlaytnNode(t, workspace, validator, nil)
	assert.NoError(t, err)
	if err := node.StartMining(false); err != nil {
		t.Fatal()
	}
	time.Sleep(2 * time.Second)

	// stop node before ending the test code
	if err := fullNode.Stop(); err != nil {
		t.Fatal(err)
	}
}

func newBlockchain(t *testing.T, config *params.ChainConfig) (*node.Node, *cn.CN, *TestAccountType, *big.Int, string) {
	t.Log("Create a new blockchain")
	// Prepare workspace
	workspace, err := ioutil.TempDir("", "klaytn-test-state")
	if err != nil {
		t.Fatalf("failed to create temporary keystore: %v", err)
	}
	t.Log("Workspace is ", workspace)

	// Prepare a validator
	validator, err := createAnonymousAccount(getRandomPrivateKeyString(t))
	if err != nil {
		t.Fatal(err)
	}

	// Create a Klaytn node
	fullNode, node, err := newKlaytnNode(t, workspace, validator, config)
	assert.NoError(t, err)
	if err := node.StartMining(false); err != nil {
		t.Fatal()
	}
	time.Sleep(2 * time.Second) // wait for initializing mining

	chainId := node.BlockChain().Config().ChainID

	return fullNode, node, validator, chainId, workspace
}

func createAccount(t *testing.T, numAccounts int, validator *TestAccountType) (*TestAccountType, []*TestAccountType, []*TestAccountType) {
	accounts := make([]*TestAccountType, numAccounts)
	contractAccounts := make([]*TestAccountType, numAccounts)

	// richAccount is used for deploying smart contracts
	richAccount := &TestAccountType{
		Addr:  validator.Addr,
		Keys:  []*ecdsa.PrivateKey{validator.Keys[0]},
		Nonce: uint64(0),
	}

	var err error
	for i := 0; i < numAccounts; i++ {
		if accounts[i], err = createAnonymousAccount(getRandomPrivateKeyString(t)); err != nil {
			t.Fatal()
		}
		// address should be overwritten
		if contractAccounts[i], err = createAnonymousAccount(getRandomPrivateKeyString(t)); err != nil {
			t.Fatal()
		}
	}

	return richAccount, accounts, contractAccounts
}

// newKlaytnNode creates a klaytn node
func newKlaytnNode(t *testing.T, dir string, validator *TestAccountType, config *params.ChainConfig) (*node.Node, *cn.CN, error) {
	var klaytnNode *cn.CN

	fullNode, err := node.New(&node.Config{
		DataDir:           dir,
		UseLightweightKDF: true,
		P2P:               p2p.Config{PrivateKey: validator.Keys[0], NoListen: true},
	})
	if err != nil {
		t.Fatalf("failed to create node: %v", err)
	}

	istanbulConfData, err := rlp.EncodeToBytes(&types.IstanbulExtra{
		Validators:    []common.Address{validator.Addr},
		Seal:          []byte{},
		CommittedSeal: [][]byte{},
	})
	if err != nil {
		t.Fatal(err)
	}

	genesis := blockchain.DefaultGenesisBlock()
	genesis.ExtraData = genesis.ExtraData[:types.IstanbulExtraVanity]
	genesis.ExtraData = append(genesis.ExtraData, istanbulConfData...)

	if config == nil {
		genesis.Config = params.CypressChainConfig.Copy()
		genesis.Config.Istanbul.SubGroupSize = 1
		genesis.Config.Istanbul.ProposerPolicy = uint64(istanbul.RoundRobin)
		genesis.Config.Governance.Reward.MintingAmount = new(big.Int).Mul(big.NewInt(9000000000000000000), big.NewInt(params.KLAY))
	} else {
		genesis.Config = config
	}

	cnConf := cn.GetDefaultConfig()
	cnConf.Genesis = genesis
	cnConf.Rewardbase = validator.Addr
	cnConf.SingleDB = false
	cnConf.NumStateTrieShards = 4

	ks := fullNode.AccountManager().Backends(keystore.KeyStoreType)[0].(*keystore.KeyStore)
	_, _ = ks.ImportECDSA(validator.Keys[0], "") // import a node key

	if err = fullNode.Register(func(ctx *node.ServiceContext) (node.Service, error) { return cn.New(ctx, cnConf) }); err != nil {
		return nil, nil, errors.WithMessage(err, "failed to register Klaytn protocol")
	}

	if err = fullNode.Start(); err != nil {
		return nil, nil, errors.WithMessage(err, "failed to start test fullNode")
	}

	if err := fullNode.Service(&klaytnNode); err != nil {
		return nil, nil, err
	}

	return fullNode, klaytnNode, nil
}

// deployRandomTxs creates a random transaction
func deployRandomTxs(t *testing.T, txpool work.TxPool, chainId *big.Int, sender *TestAccountType, txNum int) []*types.Transaction {
	var tx *types.Transaction
	var txs []*types.Transaction
	signer := types.LatestSignerForChainID(chainId)
	gasPrice := txpool.GasPrice()

	txNuminABlock := 100
	for i := 0; i < txNum; i++ {
		if i%txNuminABlock == 0 {
			time.Sleep(time.Second)
		}

		receiver, err := createAnonymousAccount(getRandomPrivateKeyString(t))
		require.Nil(t, err)

		tx, _ = genLegacyTransaction(t, signer, sender, receiver, nil, gasPrice)

		err = txpool.AddLocal(tx)
		require.True(t, err == nil || err == blockchain.ErrAlreadyNonceExistInPool)

		txs = append(txs, tx)
		sender.AddNonce()
	}
	return txs
}

// deployValueTransferTx deploy value transfer transactions
func deployValueTransferTx(t *testing.T, txpool work.TxPool, chainId *big.Int, sender *TestAccountType, toAcc *TestAccountType) *types.Transaction {
	signer := types.LatestSignerForChainID(chainId)
	gasPrice := txpool.GasPrice()

	tx, _ := genLegacyTransaction(t, signer, sender, toAcc, nil, gasPrice)

	err := txpool.AddLocal(tx)
	require.True(t, err == nil || err == blockchain.ErrAlreadyNonceExistInPool)

	sender.AddNonce()
	return tx
}

// deployContractDeployTx deploys a contract
func deployContractDeployTx(t *testing.T, txpool work.TxPool, chainId *big.Int, sender *TestAccountType, code string) (*types.Transaction, common.Address) {
	signer := types.LatestSignerForChainID(chainId)
	gasPrice := txpool.GasPrice()

	values := map[types.TxValueKeyType]interface{}{
		types.TxValueKeyNonce:         sender.GetNonce(),
		types.TxValueKeyFrom:          sender.GetAddr(),
		types.TxValueKeyTo:            (*common.Address)(nil),
		types.TxValueKeyAmount:        new(big.Int).SetUint64(0),
		types.TxValueKeyGasLimit:      uint64(1e8),
		types.TxValueKeyGasPrice:      gasPrice,
		types.TxValueKeyData:          common.FromHex(code),
		types.TxValueKeyCodeFormat:    params.CodeFormatEVM,
		types.TxValueKeyHumanReadable: false,
	}

	tx, err := types.NewTransactionWithMap(types.TxTypeSmartContractDeploy, values)
	require.Nil(t, err)

	err = tx.SignWithKeys(signer, sender.GetTxKeys())
	require.Nil(t, err)

	err = txpool.AddLocal(tx)
	if err != nil && err != blockchain.ErrAlreadyNonceExistInPool {
		t.Fatal(err)
	}

	contractAddr := crypto.CreateAddress(sender.Addr, sender.Nonce)
	sender.AddNonce()
	return tx, contractAddr
}

func deployContractExecutionTx(t *testing.T, txpool work.TxPool, chainId *big.Int, sender *TestAccountType, contractAddr common.Address, calldata string) *types.Transaction {
	signer := types.LatestSignerForChainID(chainId)
	gasPrice := txpool.GasPrice()

	values := map[types.TxValueKeyType]interface{}{
		types.TxValueKeyNonce:    sender.GetNonce(),
		types.TxValueKeyFrom:     sender.GetAddr(),
		types.TxValueKeyTo:       contractAddr,
		types.TxValueKeyAmount:   new(big.Int).SetUint64(0),
		types.TxValueKeyGasLimit: uint64(1e8),
		types.TxValueKeyGasPrice: gasPrice,
		types.TxValueKeyData:     common.FromHex(calldata),
	}

	tx, err := types.NewTransactionWithMap(types.TxTypeSmartContractExecution, values)
	require.Nil(t, err)

	err = tx.SignWithKeys(signer, sender.GetTxKeys())
	require.Nil(t, err)

	err = txpool.AddLocal(tx)
	if err != nil && err != blockchain.ErrAlreadyNonceExistInPool {
		t.Fatal(err)
	}

	sender.AddNonce()
	return tx
}

// Wait until the receipt for `txhash` is ready
// Returns the receipt
// Returns nil after a reasonable timeout
func waitReceipt(chain *blockchain.BlockChain, txhash common.Hash) *types.Receipt {
	if receipt := chain.GetReceiptByTxHash(txhash); receipt != nil {
		return receipt
	}
	chainEventCh := make(chan blockchain.ChainEvent)
	subscription := chain.SubscribeChainEvent(chainEventCh)
	defer subscription.Unsubscribe()
	timeout := time.NewTimer(15 * time.Second)
	for {
		select {
		case <-timeout.C:
			return nil
		case <-chainEventCh:
			receipt := chain.GetReceiptByTxHash(txhash)
			if receipt != nil {
				return receipt
			}
		}
	}
}

// Wait until `num` block is mined
// Returns the header with the number larger or equal to `num`
// Returns nil after a reasonable timeout
func waitBlock(chain work.BlockChain, num uint64) *types.Header {
	head := chain.CurrentHeader()
	if head.Number.Uint64() >= num {
		return head
	}

	chainEventCh := make(chan blockchain.ChainEvent)
	subscription := chain.SubscribeChainEvent(chainEventCh)
	defer subscription.Unsubscribe()
	// Wait until desired `num` plus 10 seconds margin
	timeoutSec := num - head.Number.Uint64() + 10
	timeout := time.NewTimer(time.Duration(timeoutSec) * time.Second)
	for {
		select {
		case <-timeout.C:
			return nil
		case <-chainEventCh:
			head := chain.CurrentHeader()
			if head.Number.Uint64() >= num {
				return head
			}
		}
	}
}
