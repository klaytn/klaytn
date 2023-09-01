package tests

import (
	"math/big"
	"os"
	"testing"

	"github.com/klaytn/klaytn/accounts/abi/bind"
	"github.com/klaytn/klaytn/accounts/abi/bind/backends"
	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/consensus/istanbul"
	"github.com/klaytn/klaytn/contracts/kip103"
	"github.com/klaytn/klaytn/log"
	"github.com/klaytn/klaytn/params"
	"github.com/stretchr/testify/assert"
)

func TestRebalanceTreasury_EOA(t *testing.T) {
	log.EnableLogForTest(log.LvlError, log.LvlInfo)

	// prepare chain configuration
	config := params.CypressChainConfig.Copy()
	config.LondonCompatibleBlock = big.NewInt(0)
	config.IstanbulCompatibleBlock = big.NewInt(0)
	config.EthTxTypeCompatibleBlock = big.NewInt(0)
	config.MagmaCompatibleBlock = big.NewInt(0)
	config.KoreCompatibleBlock = big.NewInt(0)
	config.Istanbul.SubGroupSize = 1
	config.Istanbul.ProposerPolicy = uint64(istanbul.RoundRobin)
	config.Governance.Reward.MintingAmount = new(big.Int).Mul(big.NewInt(9000000000000000000), big.NewInt(params.KLAY))

	// make a blockchain node
	fullNode, node, validator, _, workspace := newBlockchain(t, config)
	defer func() {
		fullNode.Stop()
		os.RemoveAll(workspace)
	}()

	optsOwner := bind.NewKeyedTransactor(validator.Keys[0])
	transactor := backends.NewBlockchainContractBackend(node.BlockChain(), node.TxPool().(*blockchain.TxPool), nil)
	// We need to wait for the following contract executions to be processed, so let's have enough number of blocks
	targetBlockNum := new(big.Int).Add(node.BlockChain().CurrentBlock().Number(), big.NewInt(10))

	contractAddr, tx, contract, err := kip103.DeployTreasuryRebalance(optsOwner, transactor, targetBlockNum)
	if err != nil {
		t.Fatal(err)
	}
	receipt := waitReceipt(node.BlockChain().(*blockchain.BlockChain), tx.Hash())
	if receipt == nil {
		t.Fatal("timeout")
	}

	// set kip103 hardfork config
	node.BlockChain().Config().Kip103CompatibleBlock = targetBlockNum
	node.BlockChain().Config().Kip103ContractAddress = contractAddr

	t.Log("ContractOwner Addr:", validator.GetAddr().String())
	t.Log("Contract Addr:", contractAddr.String())
	t.Log("Target Block:", targetBlockNum.Int64())

	// prepare newbie accounts
	numNewbie := 3
	newbieAccs := make([]TestAccount, numNewbie)
	newbieAllocs := make([]*big.Int, numNewbie)

	state, err := node.BlockChain().State()
	if err != nil {
		t.Fatal(err)
	}
	totalNewbieAlloc := state.GetBalance(validator.Addr)
	t.Log("Total Newbie amount: ", totalNewbieAlloc)

	for i := 0; i < numNewbie; i++ {
		newbieAccs[i] = genKlaytnLegacyAccount(t)
		newbieAllocs[i] = new(big.Int).Div(totalNewbieAlloc, big.NewInt(2))
		totalNewbieAlloc.Sub(totalNewbieAlloc, newbieAllocs[i])

		t.Log("Newbie", i, "Addr:", newbieAccs[i].GetAddr().String())
		t.Log("Newbie", i, "Amount:", newbieAllocs[i])
	}

	// register RegisterRetired
	if _, err := contract.RegisterRetired(optsOwner, validator.Addr); err != nil {
		t.Fatal(err)
	}

	// register newbies
	for i, newbie := range newbieAccs {
		if _, err := contract.RegisterNewbie(optsOwner, newbie.GetAddr(), newbieAllocs[i]); err != nil {
			t.Fatal(err)
		}
	}

	// initialized -> registered
	if tx, err = contract.FinalizeRegistration(optsOwner); err != nil {
		t.Fatal(err)
	}
	// Should wait for this tx to be processed, or next tx will be failed when estimating gas
	receipt = waitReceipt(node.BlockChain().(*blockchain.BlockChain), tx.Hash())
	if receipt == nil {
		t.Fatal("timeout")
	}

	// approve
	if tx, err = contract.Approve(optsOwner, validator.GetAddr()); err != nil {
		t.Fatal(err)
	}
	// Should wait for this tx to be processed, or next tx will be failed when estimating gas
	receipt = waitReceipt(node.BlockChain().(*blockchain.BlockChain), tx.Hash())
	if receipt == nil {
		t.Fatal("timeout")
	}

	// registered -> approved
	if _, err := contract.FinalizeApproval(optsOwner); err != nil {
		t.Fatal(err)
	}

	header := waitBlock(node.BlockChain(), targetBlockNum.Uint64())
	if header == nil {
		t.Fatal("timeout")
	}

	curState, err := node.BlockChain().StateAt(header.Root)
	if err != nil {
		t.Fatal(err)
	}

	balRetired := curState.GetBalance(validator.GetAddr())
	assert.Equal(t, balRetired, big.NewInt(0))

	for j, newbie := range newbieAccs {
		balNewbie := curState.GetBalance(newbie.GetAddr())
		assert.Equal(t, newbieAllocs[j], balNewbie)
	}
}
