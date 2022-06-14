package tests

import (
	"encoding/hex"
	"math/big"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/klaytn/klaytn/accounts/abi"
	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	govcontract "github.com/klaytn/klaytn/contracts/gov"
	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/governance"
	"github.com/klaytn/klaytn/log"
	"github.com/klaytn/klaytn/networks/rpc"
	"github.com/klaytn/klaytn/node/cn"
	"github.com/klaytn/klaytn/params"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test that ContractEngine can read appropriate parameter values before and after
// the contract is deployed.
func TestGovernance_ContractEngine(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlWarn)

	fullNode, node, validator, chainId, workspace := newBlockchain(t)
	defer os.RemoveAll(workspace)
	defer fullNode.Stop()

	var (
		chain        = node.BlockChain().(*blockchain.BlockChain)
		config       = chain.Config()
		owner        = validator
		contractAddr = crypto.CreateAddress(owner.Addr, owner.Nonce)

		paramName     = "istanbul.committeesize"
		fallbackValue = config.Istanbul.SubGroupSize
		paramValue    = uint64(22)
		paramBytes    = []byte{22}

		deployBlock   uint64 // Before deploy: 0, After deploy: the deployed block
		setParamBlock uint64 // Before setParam: 0, After setParam: the setParam'd block
	)

	// Here we are running (tx sender) and (param reader) in parallel.
	// This is to check that param reader works under various situations such as
	// (a) contract is not yet deployed (b) parameter is not yet set (c) parameter is set.

	// Run tx sender thread
	go func() {
		num, _, _ := deployGovParamTx_constructor(t, node, owner, chainId)
		deployBlock = num

		num, _ = deployGovParamTx_setParam(t, node, owner, chainId, contractAddr, paramName, paramBytes)
		setParamBlock = num
	}()

	// Run param reader thread
	config.Governance.GovernanceContract = contractAddr
	engine := governance.NewContractEngine(config, nil)
	engine.SetBlockchain(chain)

	// Validate current params from engine.Params(), alongside block processing.
	// At block #N, Params() returns the parameters to be used when building
	// block #N+1 (i.e. pending block).
	chainEventCh := make(chan blockchain.ChainEvent)
	subscription := chain.SubscribeChainEvent(chainEventCh)
	defer subscription.Unsubscribe()
	for {
		ev := <-chainEventCh
		time.Sleep(100 * time.Millisecond) // wait for tx sender thread to set deployBlock, etc.

		num := ev.Block.Number().Uint64()
		err := engine.UpdateParams()
		assert.Nil(t, err)

		value, _ := engine.Params().Get(params.CommitteeSize)
		t.Logf("Params() at block=%2d: %v", num, value)

		if deployBlock == 0 { // not yet deployed
			assert.Equal(t, fallbackValue, value)
		} else if setParamBlock == 0 { // not yet setParam'd
			assert.Equal(t, fallbackValue, value)
		} else { // after setParam
			assert.Equal(t, paramValue, value)
		}
		if setParamBlock != 0 && num >= setParamBlock+2 {
			break
		}
	}

	// Validate historic params from engine.ParamsAt(n)
	for num := uint64(0); num <= setParamBlock+2; num++ {
		pset, err := engine.ParamsAt(num)
		assert.Nil(t, err)
		assert.NotNil(t, pset)

		value, _ := pset.Get(params.CommitteeSize)
		t.Logf("ParamsAt(block=%2d): %v", num, value)
		if num < deployBlock { // not yet deployed
			assert.Equal(t, fallbackValue, value)
		} else if num <= setParamBlock { // not yet setParam'd
			assert.Equal(t, fallbackValue, value)
		} else { // after setParam
			assert.Equal(t, paramValue, value)
		}
	}
}

// Test that vote for "governance.governancecontract" works.
func TestGovernance_MixedEngine(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlWarn)

	fullNode, node, validator, chainId, workspace := newBlockchain(t)
	defer os.RemoveAll(workspace)
	defer fullNode.Stop()

	var (
		chain        = node.BlockChain().(*blockchain.BlockChain)
		db           = node.ChainDB()
		config       = chain.Config()
		owner        = validator
		contractAddr = crypto.CreateAddress(owner.Addr, owner.Nonce)
		rpcServer, _ = fullNode.RPCHandler()
		rpcClient    = rpc.DialInProc(rpcServer)
	)

	engine := governance.NewMixedEngine(config, db)
	engine.SetBlockchain(chain)

	// Check assumptions in ChainConfig. Should have been configured in newBlockchain().
	require.Equal(t, uint64(3), config.Istanbul.Epoch)                         // for quicker test
	require.Equal(t, owner.Addr, config.Governance.GoverningNode)              // allow governance_vote
	require.True(t, common.EmptyAddress(config.Governance.GovernanceContract)) // no contract at genesis

	// Prepare the contract before vote the contract address
	deployGovParamTx_constructor(t, node, owner, chainId)      // deploy
	configToWrite := config.Copy()                             // upload copy of ChainConfig
	configToWrite.Governance.GovernanceContract = contractAddr // ..except GovernanceContract
	bytesMap := chainConfigToBytesMap(t, configToWrite)
	deployGovParamTx_batchSetParam(t, node, owner, chainId, contractAddr, bytesMap)

	// Vote for the contract address
	var res interface{}
	require.Nil(t, rpcClient.Call(&res, "governance_vote", "governance.governancecontract", contractAddr))
	t.Logf("Voted at block=%2d returned='%s'", chain.CurrentHeader().Number.Uint64(), res)

	// Check that GovernanceContract param changes correctly
	num := chain.CurrentHeader().Number.Uint64()
	pset, err := engine.ParamsAt(num)
	require.Nil(t, err)
	assert.Equal(t, common.Address{}, pset.GovernanceContract()) // zero address before vote is applied.

	// Wait until the vote takes effect.
	// It takes at most 3*Epoch blocks for the vote to take effect.
	waitBlock(chain, num+9)

	pset, err = engine.ParamsAt(num)
	require.Nil(t, err)
	assert.NotEqual(t, contractAddr, pset.GovernanceContract()) // set to contractAddr after vote is applied
}

// Test the validation lobic of voting for "governance.governancecontract" item
func TestGovernance_VoteContractAddr(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlWarn)

	fullNode, node, validator, chainId, workspace := newBlockchain(t)
	defer os.RemoveAll(workspace)
	defer fullNode.Stop()

	var (
		chain        = node.BlockChain().(*blockchain.BlockChain)
		db           = node.ChainDB()
		config       = chain.Config()
		owner        = validator
		contractAddr = crypto.CreateAddress(owner.Addr, owner.Nonce)
		rpcServer, _ = fullNode.RPCHandler()
		rpcClient    = rpc.DialInProc(rpcServer)
	)

	engine := governance.NewMixedEngine(config, db)
	engine.SetBlockchain(chain)

	// Check assumptions in ChainConfig. Should have been configured in newBlockchain().
	require.Equal(t, uint64(3), config.Istanbul.Epoch)                         // for quicker test
	require.Equal(t, owner.Addr, config.Governance.GoverningNode)              // allow governance_vote
	require.True(t, common.EmptyAddress(config.Governance.GovernanceContract)) // no contract at genesis

	// 1-1. Deploy the contract but wrong values
	deployGovParamTx_constructor(t, node, owner, chainId) // deploy
	configToWrite := config.Copy()                        // upload copy of ChainConfig
	configToWrite.Istanbul.SubGroupSize = 777             // ..with intentionally wrong value
	bytesMap := chainConfigToBytesMap(t, configToWrite)
	deployGovParamTx_batchSetParam(t, node, owner, chainId, contractAddr, bytesMap)

	// 1-2. Vote fails because of parameter mismatch
	num := chain.CurrentHeader().Number.Uint64()
	require.Nil(t, rpcClient.Call(nil, "governance_vote", "governance.governancecontract", contractAddr))
	t.Logf("Voted at block=%2d", num)
	waitBlock(chain, num+9) // wait for 3*Epochs for vote to apply
	pset, err := engine.ParamsAt(num + 9)
	assert.Nil(t, err)
	assert.Equal(t, common.Address{}, pset.GovernanceContract()) // still zero

	// 2-1. Overwrite the wrong param by correct value
	name := "istanbul.committeesize"
	value := []byte{byte(config.Istanbul.SubGroupSize)}
	deployGovParamTx_setParam(t, node, owner, chainId, contractAddr, name, value)

	// 2-2. Vote fails because governancecontract == 0 in contract
	num = chain.CurrentHeader().Number.Uint64()
	require.Nil(t, rpcClient.Call(nil, "governance_vote", "governance.governancecontract", contractAddr))
	t.Logf("Voted at block=%2d", num)
	waitBlock(chain, num+9) // wait for 3*Epochs for vote to apply
	pset, err = engine.ParamsAt(num + 9)
	assert.Nil(t, err)
	assert.Equal(t, common.Address{}, pset.GovernanceContract()) // still zero

	// 3-1. Correct the contract by setting the GovernanceContract param
	name = "governance.governancecontract"
	value = contractAddr.Bytes()
	deployGovParamTx_setParam(t, node, owner, chainId, contractAddr, name, value)

	// 3-2. Vote succeeds
	num = chain.CurrentHeader().Number.Uint64()
	require.Nil(t, rpcClient.Call(nil, "governance_vote", "governance.governancecontract", contractAddr))
	t.Logf("Voted at block=%2d", num)
	waitBlock(chain, num+9) // wait for 3*Epochs for vote to apply
	pset, err = engine.ParamsAt(num + 9)
	assert.Nil(t, err)
	assert.Equal(t, contractAddr, pset.GovernanceContract()) // now changed
}

func deployGovParamTx_constructor(t *testing.T, node *cn.CN, owner *TestAccountType, chainId *big.Int,
) (uint64, common.Address, *types.Transaction) {
	var (
		// Deploy contract: constructor(address _owner)
		contractAbi, _ = abi.JSON(strings.NewReader(govcontract.GovParamABI))
		contractBin    = govcontract.GovParamBin
		ctorArgs, _    = contractAbi.Pack("")
		code           = contractBin + hex.EncodeToString(ctorArgs)
		initArgs, _    = contractAbi.Pack("initialize", owner.Addr)
		initData       = common.ToHex(initArgs)
	)

	// Deploy contract
	tx, addr := deployContractDeployTx(t, node.TxPool(), chainId, owner, code)

	chain := node.BlockChain().(*blockchain.BlockChain)
	receipt := waitReceipt(chain, tx.Hash())
	require.NotNil(t, receipt)
	require.Equal(t, types.ReceiptStatusSuccessful, receipt.Status)

	_, _, num, _ := chain.GetTxAndLookupInfo(tx.Hash())
	t.Logf("GovParam deployed at block=%2d, addr=%s", num, addr.Hex())

	// Call initialize()
	tx = deployContractExecutionTx(t, node.TxPool(), chainId, owner, addr, initData)
	receipt = waitReceipt(chain, tx.Hash())
	require.NotNil(t, receipt)
	require.Equal(t, types.ReceiptStatusSuccessful, receipt.Status)

	return num, addr, tx
}

func deployGovParamTx_setParam(t *testing.T, node *cn.CN, owner *TestAccountType, chainId *big.Int,
	contractAddr common.Address, name string, value []byte,
) (uint64, *types.Transaction) {
	var (
		// Set parameter: setParam(string name, bytes value)
		contractAbi, _ = abi.JSON(strings.NewReader(govcontract.GovParamABI))
		callArgs, _    = contractAbi.Pack("setParamByOwner", name, value)
		data           = common.ToHex(callArgs)
	)

	tx := deployContractExecutionTx(t, node.TxPool(), chainId, owner, contractAddr, data)

	chain := node.BlockChain().(*blockchain.BlockChain)
	receipt := waitReceipt(chain, tx.Hash())
	require.NotNil(t, receipt)
	require.Equal(t, types.ReceiptStatusSuccessful, receipt.Status)

	_, _, num, _ := chain.GetTxAndLookupInfo(tx.Hash())
	t.Logf("GovParam.setParam executed at block=%2d", num)
	return num, tx
}

func deployGovParamTx_batchSetParam(t *testing.T, node *cn.CN, owner *TestAccountType, chainId *big.Int,
	contractAddr common.Address, bytesMap map[string][]byte,
) []*types.Transaction {
	var (
		chain          = node.BlockChain().(*blockchain.BlockChain)
		beginBlock     = chain.CurrentHeader().Number.Uint64()
		contractAbi, _ = abi.JSON(strings.NewReader(govcontract.GovParamABI))
		txs            = []*types.Transaction{}
	)

	// Send all setParam() calls at once
	for name, value := range bytesMap {
		// Set parameter: setParam(string name, bytes value)
		callArgs, _ := contractAbi.Pack("setParamByOwner", name, value)
		data := common.ToHex(callArgs)
		tx := deployContractExecutionTx(t, node.TxPool(), chainId, owner, contractAddr, data)
		txs = append(txs, tx)
	}

	// Wait for all txs
	for _, tx := range txs {
		receipt := waitReceipt(chain, tx.Hash())
		require.NotNil(t, receipt)
		require.Equal(t, types.ReceiptStatusSuccessful, receipt.Status)
	}
	num := chain.CurrentHeader().Number.Uint64()
	t.Logf("GovParam.setParam executed %d times between blocks=%2d,%2d", len(txs), beginBlock, num)
	return txs
}

// Klaytn node only decodes the byte-array param values (refer to params/governance_paramset.go).
// Encoding is the job of transaction senders (i.e. clients and dApps).
// This is a reference implementation of such encoder.
func chainConfigToBytesMap(t *testing.T, config *params.ChainConfig) map[string][]byte {
	pset, err := params.NewGovParamSetChainConfig(config)
	require.Nil(t, err)
	strMap := pset.StrMap()

	bytesMap := map[string][]byte{}
	for name, value := range strMap {
		switch value.(type) {
		case string:
			bytesMap[name] = []byte(value.(string))
		case common.Address:
			bytesMap[name] = value.(common.Address).Bytes()
		case uint64:
			bytesMap[name] = new(big.Int).SetUint64(value.(uint64)).Bytes()
		case bool:
			if value.(bool) == true {
				bytesMap[name] = []byte{0x01}
			} else {
				bytesMap[name] = []byte{0x00}
			}
		}
	}

	// Check that bytesMap is correct just in case
	qset, err := params.NewGovParamSetBytesMap(bytesMap)
	require.Nil(t, err)
	require.Equal(t, pset.StrMap(), qset.StrMap())
	return bytesMap
}
