package tests

import (
	"context"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/klaytn/klaytn"
	"github.com/klaytn/klaytn/accounts/abi/bind"
	"github.com/klaytn/klaytn/api"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/hexutil"
	"github.com/klaytn/klaytn/consensus/istanbul"
	"github.com/klaytn/klaytn/contracts/kip103"
	"github.com/klaytn/klaytn/log"
	"github.com/klaytn/klaytn/networks/rpc"
	"github.com/klaytn/klaytn/node/cn"
	"github.com/klaytn/klaytn/params"
	"github.com/stretchr/testify/assert"
)

type testKip103TxTransactor struct {
	node *cn.CN
}

func (t *testKip103TxTransactor) FilterLogs(ctx context.Context, query klaytn.FilterQuery) ([]types.Log, error) {
	return nil, nil
}

func (t *testKip103TxTransactor) SubscribeFilterLogs(ctx context.Context, query klaytn.FilterQuery, ch chan<- types.Log) (klaytn.Subscription, error) {
	return nil, nil
}

func (t *testKip103TxTransactor) PendingCodeAt(ctx context.Context, account common.Address) ([]byte, error) {
	return t.CodeAt(ctx, account, nil)
}

func (t *testKip103TxTransactor) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	return t.node.TxPool().GetPendingNonce(account), nil
}

func (t *testKip103TxTransactor) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	return big.NewInt(int64(t.node.BlockChain().Config().UnitPrice)), nil
}

func (t *testKip103TxTransactor) EstimateGas(ctx context.Context, call klaytn.CallMsg) (gas uint64, err error) {
	return uint64(1e8), nil
}

func (t *testKip103TxTransactor) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	return t.node.TxPool().AddLocal(tx)
}

func (t *testKip103TxTransactor) ChainID(ctx context.Context) (*big.Int, error) {
	return t.node.BlockChain().Config().ChainID, nil
}

func (t *testKip103TxTransactor) CodeAt(ctx context.Context, contract common.Address, blockNumber *big.Int) ([]byte, error) {
	if blockNumber == nil {
		blockNumber = t.node.BlockChain().CurrentBlock().Number()
	}
	root := t.node.BlockChain().GetHeaderByNumber(blockNumber.Uint64()).Root
	state, err := t.node.BlockChain().StateAt(root)
	if err != nil {
		return nil, err
	}
	return state.GetCode(contract), nil
}

func (t *testKip103TxTransactor) CallContract(ctx context.Context, call klaytn.CallMsg, blockNumber *big.Int) ([]byte, error) {
	if blockNumber == nil {
		blockNumber = t.node.BlockChain().CurrentBlock().Number()
	}

	price := hexutil.Big(*t.node.TxPool().GasPrice())
	if call.GasPrice != nil {
		price = hexutil.Big(*call.GasPrice)
	}

	value := hexutil.Big(*big.NewInt(0))
	if call.Value != nil {
		value = hexutil.Big(*call.Value)
	}

	arg := api.CallArgs{From: call.From, To: call.To, Gas: hexutil.Uint64(1e8), GasPrice: &price, Value: value, Data: call.Data}
	bn := rpc.BlockNumber(blockNumber.Int64())

	apiBackend := api.NewPublicBlockChainAPI(t.node.APIBackend)
	return apiBackend.Call(ctx, arg, rpc.NewBlockNumberOrHashWithNumber(bn))
}

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
	transactor := &testKip103TxTransactor{node: node}
	targetBlockNum := new(big.Int).Add(node.BlockChain().CurrentBlock().Number(), big.NewInt(5))

	contractAddr, _, contract, err := kip103.DeployTreasuryRebalance(optsOwner, transactor, targetBlockNum)
	if err != nil {
		t.Fatal(err)
	}

	// set kip103 hardfork config
	node.BlockChain().Config().KIP103 = &params.KIP103Config{
		Kip103CompatibleBlock: targetBlockNum,
		Kip103ContractAddress: contractAddr,
	}

	t.Log("ContractOwner Addr:", validator.GetAddr().String())
	t.Log("Contract Addr:", contractAddr.String())
	t.Log("Target Block:", targetBlockNum.Int64())

	// naive waiting for tx processing
	time.Sleep(2 * time.Second)

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
	if _, err := contract.FinalizeRegistration(optsOwner); err != nil {
		t.Fatal(err)
	}

	// approve
	if _, err := contract.Approve(optsOwner, validator.GetAddr()); err != nil {
		t.Fatal(err)
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
