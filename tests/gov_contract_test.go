package tests

import (
	"encoding/hex"
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
	"github.com/klaytn/klaytn/params"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGovernance_ContractEngine(t *testing.T) {
	// enableLog()

	fullNode, node, validator, chainId, workspace := newBlockchain(t)
	defer os.RemoveAll(workspace)
	defer fullNode.Stop()

	var (
		chain = node.BlockChain().(*blockchain.BlockChain)
		owner = validator
		// We can pre-calculate contract address from the owner account address.
		contractAddr = crypto.CreateAddress(owner.Addr, owner.Nonce)

		contractBin    = govcontract.GovParamBin
		contractAbi, _ = abi.JSON(strings.NewReader(govcontract.GovParamABI))
		paramName      = "istanbul.committeesize"
		paramValue     = uint64(22)
		paramBytes     = []byte{22}

		deployBlock     uint64 // Before deploy: 0, After deploy: the deployed block
		activationBlock uint64 // Before deploy: 0, After deploy: some number
		setParamBlock   uint64 // Before setParam: 0, After setParam: the setParam'd block

		engine *governance.ContractEngine
	)

	// Here we are running (tx sender) and (param reader) in parallel.
	// This is to check that param reader works under various corner cases including
	// (a) contract is not yet deployed (b) parameter is not yet set (c) parameter is set.

	// Run tx sender thread
	go func() {
		// Deploy contract: constructor(address _owner)
		ctorArgs, _ := contractAbi.Pack("", owner.Addr)
		code := contractBin + hex.EncodeToString(ctorArgs)
		addr, tx := deployContractDeployTx(t, node.TxPool(), chainId, owner, code)
		require.Equal(t, contractAddr, addr)

		receipt := waitReceipt(chain, tx.Hash())
		require.NotNil(t, receipt)
		require.Equal(t, types.ReceiptStatusSuccessful, receipt.Status)

		_, _, deployBlock, _ = chain.GetTxAndLookupInfo(tx.Hash())
		t.Logf("contract deployed at block=%2d, addr=%s", deployBlock, addr.Hex())

		activationBlock = chain.CurrentHeader().Number.Uint64() + 7
		t.Logf("param activation scheduled at block=%2d", activationBlock)

		// Set parameter: setParam(string name, bytes value, uint64 activation)
		callArgs, _ := contractAbi.Pack("setParam", paramName, paramBytes, activationBlock)
		data := common.ToHex(callArgs)
		tx = deployContractExecutionTx(t, node.TxPool(), chainId, owner, contractAddr, data)

		receipt = waitReceipt(chain, tx.Hash())
		require.NotNil(t, receipt)
		require.Equal(t, types.ReceiptStatusSuccessful, receipt.Status)

		_, _, setParamBlock, _ = chain.GetTxAndLookupInfo(tx.Hash())
		t.Logf("setParam executed at block=%2d", setParamBlock)
	}()

	// Run param reader thread
	engine = governance.NewContractEngine(contractAddr)
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
			assert.Equal(t, nil, value)
		} else if setParamBlock == 0 { // not yet setParam'd
			assert.Equal(t, nil, value)
		} else if num+1 < activationBlock { // before activation
			assert.Equal(t, uint64(0), value)
		} else if num+1 >= activationBlock { // after activation
			assert.Equal(t, paramValue, value)
		}
		if activationBlock != 0 && num >= activationBlock+2 {
			break
		}
	}

	// Validate historic params from engine.ParamsAt(n)
	for num := uint64(0); num <= activationBlock+2; num++ {
		pset, err := engine.ParamsAt(num)
		assert.Nil(t, err)
		assert.NotNil(t, pset)

		value, _ := pset.Get(params.CommitteeSize)
		t.Logf("ParamsAt(block=%2d): %v", num, value)
		if num < deployBlock { // not yet deployed
			assert.Equal(t, nil, value)
		} else if num < setParamBlock { // not yet setParam'd
			assert.Equal(t, nil, value)
		} else if num < activationBlock { // before activation
			assert.Equal(t, uint64(0), value)
		} else if num >= activationBlock { // after activation
			assert.Equal(t, paramValue, value)
		}
	}
}
