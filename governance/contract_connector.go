package governance

import (
	"errors"
	"math/big"
	"strings"

	"github.com/klaytn/klaytn/accounts/abi"
	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/blockchain/vm"
	"github.com/klaytn/klaytn/common"
	govcontract "github.com/klaytn/klaytn/contracts/gov"
	"github.com/klaytn/klaytn/params"
)

type contractCaller struct {
	chainConfig  *params.ChainConfig
	chain        blockChain
	contractAddr common.Address
}

var govParamAbi, _ = abi.JSON(strings.NewReader(govcontract.GovParamABI))

func (c *contractCaller) getAllParams(num *big.Int) (*params.GovParamSet, error) {
	tx, err := c.makeTx(num, govParamAbi, "getAllParams")
	if err != nil {
		return nil, err
	}

	res, err := c.callTx(num, tx)
	if err != nil {
		return nil, err
	}

	// Parse results into GovParamSet

	// Cannot parse empty data
	if len(res) == 0 {
		return params.NewGovParamSet(), nil
	}

	var ( // c.f. contracts/gov/GovParam.go:GetAllParams()
		pNames  = new([]string)                   // *[]string = nil
		pValues = new([][]byte)                   // *[][]byte = nil
		out     = &[]interface{}{pNames, pValues} // array of pointers
	)
	if err := govParamAbi.Unpack(out, "getAllParams", res); err != nil {
		return nil, err
	}
	var ( // Retrieve the slices allocated inside Unpack().
		names  = *pNames
		values = *pValues
	)

	if len(names) != len(values) {
		logger.Error("Malformed contract.getAllParams result",
			"len(names)", len(names), "len(values)", len(values))
		return params.NewGovParamSet(), nil
	}

	bytesMap := make(map[string][]byte)
	for i := 0; i < len(names); i++ {
		bytesMap[names[i]] = values[i]
	}
	return params.NewGovParamSetBytesMap(bytesMap)
}

// Make contract execution transaction
func (c *contractCaller) makeTx(num *big.Int, contractAbi abi.ABI, fn string, args ...interface{},
) (*types.Transaction, error) {
	calldata, err := contractAbi.Pack(fn, args...)
	if err != nil {
		return nil, err
	}

	rules := c.chainConfig.Rules(num)
	intrinsicGas, err := types.IntrinsicGas(calldata, nil, false, rules)
	if err != nil {
		return nil, err
	}

	var (
		from       = common.Address{}
		to         = &c.contractAddr
		nonce      = uint64(0)
		amount     = big.NewInt(0)
		gasLimit   = uint64(1e8)
		gasPrice   = big.NewInt(0)
		checkNonce = false
	)

	tx := types.NewMessage(from, to, nonce, amount, gasLimit, gasPrice, calldata,
		checkNonce, intrinsicGas)
	return tx, nil
}

// Execute contract call using state at block `num`.
func (c *contractCaller) callTx(num *big.Int, tx *types.Transaction) ([]byte, error) {
	// Load state at given block
	block := c.chain.GetBlockByNumber(num.Uint64())
	if block == nil {
		return nil, errors.New("No such block")
	}
	statedb, err := c.chain.StateAt(block.Root())
	if err != nil {
		return nil, err
	}

	// Run EVM at given states
	evmCtx := blockchain.NewEVMContext(tx, block.Header(), c.chain, nil)
	evm := vm.NewEVM(evmCtx, statedb, c.chain.Config(), &vm.Config{})

	res, _, kerr := blockchain.ApplyMessage(evm, tx)
	if kerr.ErrTxInvalid != nil {
		return nil, kerr.ErrTxInvalid
	}

	return res, nil
}
