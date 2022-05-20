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

type ContractEngine struct {
	currentParams *params.GovParamSet
	chain         blockChain
	paramContract *deployedContract // deployed GovParam.sol
}

type deployedContract struct {
	address common.Address
	abi     abi.ABI
}

func NewContractEngine(paramAddress common.Address) *ContractEngine {
	e := &ContractEngine{
		currentParams: params.NewGovParamSet(),
		paramContract: newDeployedContract(paramAddress, govcontract.GovParamABI),
	}
	return e
}

// Params effective at upcoming block (head+1), queried at current block (head)
func (e *ContractEngine) Params() *params.GovParamSet {
	return e.currentParams
}

// Params effective at requested block (num), queried at previous block (num-1).
func (e *ContractEngine) ParamsAt(num uint64) (*params.GovParamSet, error) {
	bigNum := new(big.Int).SetUint64(num)
	return e.contractGetAllParams(bigNum, false)
}

func (e *ContractEngine) UpdateParams() error {
	if e.chain == nil {
		return errors.New("blockchain not set")
	}

	num := e.chain.CurrentHeader().Number
	pset, err := e.contractGetAllParams(num, true)
	if err != nil {
		return err
	}
	e.currentParams = pset
	return nil
}

func (e *ContractEngine) contractGetAllParams(num *big.Int, atNextBlock bool) (*params.GovParamSet, error) {
	if e.chain == nil {
		return nil, errors.New("blockchain not set")
	}

	fn := "getAllParams"
	if atNextBlock {
		fn = "getAllParamsAtNextBlock"
	}

	tx, err := e.paramContract.makeTx(num, fn)
	if err != nil {
		return nil, err
	}

	res, err := e.paramContract.callTx(e.chain, num, tx)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return params.NewGovParamSet(), nil
	}

	var views []govcontract.GovParamParamView
	if err := e.paramContract.abi.Unpack(&views, fn, res); err != nil {
		return nil, err
	}
	bytesMap := make(map[string][]byte)
	for _, pair := range views {
		bytesMap[pair.Name] = pair.Value
	}
	return params.NewGovParamSetBytesMap(bytesMap)
}

func (e *ContractEngine) SetBlockchain(chain blockChain) {
	e.chain = chain
}

// DeployedContract abstracts the contract calling process.

func newDeployedContract(address common.Address, jsonAbi string) *deployedContract {
	c := &deployedContract{
		address: address,
	}

	if abi, err := abi.JSON(strings.NewReader(govcontract.GovParamABI)); err != nil {
		logger.Crit("Cannot parse ABI", "err", err)
	} else {
		c.abi = abi
	}
	return c
}

// Make contract execution transaction
func (c *deployedContract) makeTx(num *big.Int, fn string, args ...interface{}) (*types.Transaction, error) {
	calldata, err := c.abi.Pack(fn, args...)
	if err != nil {
		return nil, err
	}

	rules := params.CypressChainConfig.Rules(num)
	intrinsicGas, err := types.IntrinsicGas(calldata, nil, false, rules)
	if err != nil {
		return nil, err
	}

	var (
		from       = common.Address{}
		to         = &c.address
		nonce      = uint64(0)
		amount     = big.NewInt(0)
		gasLimit   = uint64(1e8)
		gasPrice   = big.NewInt(0)
		checkNonce = false
	)
	tx := types.NewMessage(from, to, nonce, amount, gasLimit, gasPrice,
		calldata, checkNonce, intrinsicGas)
	return tx, nil
}

// Execute contract call using state at block `num`.
func (e *deployedContract) callTx(chain blockChain, num *big.Int, tx *types.Transaction) ([]byte, error) {
	// Load state at given block
	block := chain.GetBlockByNumber(num.Uint64())
	if block == nil {
		return nil, errors.New("No such block")
	}
	statedb, err := chain.StateAt(block.Root())
	if err != nil {
		return nil, err
	}

	// Run EVM at given states
	evmCtx := blockchain.NewEVMContext(tx, block.Header(), chain, nil)
	evm := vm.NewEVM(evmCtx, statedb, chain.Config(), &vm.Config{})

	res, _, kerr := blockchain.ApplyMessage(evm, tx)
	if kerr.ErrTxInvalid != nil {
		return nil, kerr.ErrTxInvalid
	}

	return res, nil
}
