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

var (
	errContractEngineNotReady = errors.New("ContractEngine is not ready")

	govParamAbi, _ = abi.JSON(strings.NewReader(govcontract.GovParamABI))
)

type ContractEngine struct {
	chainConfig   *params.ChainConfig
	currentParams *params.GovParamSet
	initialParams *params.GovParamSet // equals to initial ChainConfig

	defaultGov *Governance // To find the contract address at any block number
	chain      blockChain  // To access the contract state DB
}

func NewContractEngine(config *params.ChainConfig, defaultGov *Governance) *ContractEngine {
	e := &ContractEngine{
		chainConfig:   config,
		currentParams: params.NewGovParamSet(),
		defaultGov:    defaultGov,
	}

	if pset, err := params.NewGovParamSetChainConfig(config); err == nil {
		e.initialParams = pset
		e.currentParams = pset
	} else {
		logger.Crit("Error parsing initial ChainConfig", "err", err)
	}
	return e
}

func (e *ContractEngine) SetBlockchain(chain blockChain) {
	e.chain = chain
}

// Params effective at upcoming block (head+1)
func (e *ContractEngine) Params() *params.GovParamSet {
	return e.currentParams
}

// Params effective at requested block (num)
func (e *ContractEngine) ParamsAt(num uint64) (*params.GovParamSet, error) {
	if e.chain == nil {
		logger.Error("Invoked ParamsAt() before SetBlockchain", "num", num)
		return nil, errContractEngineNotReady
	}

	head := e.chain.CurrentHeader().Number.Uint64()
	if num > head {
		// Sometimes future blocks are requested.
		// ex) reward distributor in istanbul.engine.Finalize() requests ParamsAt(head+1)
		// ex) governance_itemsAt(num) API requests arbitrary num
		// In those cases we refer to the head block.
		num = head + 1
	}

	pset, err := e.contractGetAllParams(num)
	if err != nil {
		return nil, err
	}
	return params.NewGovParamSetMerged(e.initialParams, pset), nil
}

func (e *ContractEngine) UpdateParams() error {
	if e.chain == nil {
		logger.Error("Invoked UpdateParams() before SetBlockchain")
		return errContractEngineNotReady
	}

	head := e.chain.CurrentHeader().Number.Uint64()
	pset, err := e.contractGetAllParams(head + 1)
	if err != nil {
		return err
	}

	e.currentParams = params.NewGovParamSetMerged(e.initialParams, pset)
	return nil
}

func (e *ContractEngine) contractGetAllParams(num uint64) (*params.GovParamSet, error) {
	if e.chain == nil {
		logger.Error("Invoked ContractEngine before SetBlockchain")
		return nil, errContractEngineNotReady
	}

	addr := e.contractAddrAt(num)
	if common.EmptyAddress(addr) {
		logger.Error("Invoked ContractEngine but address is empty", "num", num)
		return nil, errContractEngineNotReady
	}

	caller := &contractCaller{
		chainConfig:  e.chainConfig,
		chain:        e.chain,
		contractAddr: addr,
	}
	if num > 0 {
		num -= 1
	}
	return caller.getAllParams(new(big.Int).SetUint64(num))
}

// Return the GovernanceContract address effective at given block number
func (e *ContractEngine) contractAddrAt(num uint64) common.Address {
	// Load from HeaderEngine
	if e.defaultGov != nil {
		if pset, err := e.defaultGov.ParamsAt(num); err == nil {
			if _, ok := pset.Get(params.GovernanceContract); ok {
				return pset.GovernanceContract()
			}
		}
	}

	// If database don't have the item, fallback to ChainConfig
	if e.chainConfig.Governance != nil {
		return e.chainConfig.Governance.GovernanceContract
	}

	return common.Address{}
}

type contractCaller struct {
	chainConfig  *params.ChainConfig
	chain        blockChain
	contractAddr common.Address
}

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
		pNames  = new([]string)                    // *[]string = nil
		pValues = new([]govcontract.GovParamParam) // *[]govcontract.GovParamParam = nil
		out     = &[]interface{}{pNames, pValues}  // array of pointers
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
		bytesMap[names[i]] = values[i].Value
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
