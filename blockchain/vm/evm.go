// Modifications Copyright 2018 The klaytn Authors
// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.
//
// This file is derived from core/vm/evm.go (2018/06/04).
// Modified and improved for the klaytn development.

package vm

import (
	"math/big"
	"sync/atomic"
	"time"

	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/blockchain/types/accountkey"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/kerrors"
	"github.com/klaytn/klaytn/params"
)

// emptyCodeHash is used by create to ensure deployment is disallowed to already
// deployed contract addresses (relevant after the account abstraction).
var emptyCodeHash = crypto.Keccak256Hash(nil)

const (
	CancelByCtxDone = 1 << iota
	CancelByTotalTimeLimit
)

type (
	// CanTransferFunc is the signature of a transfer guard function
	CanTransferFunc func(StateDB, common.Address, *big.Int) bool
	// TransferFunc is the signature of a transfer function
	TransferFunc func(StateDB, common.Address, common.Address, *big.Int)
	// GetHashFunc returns the nth block hash in the blockchain
	// and is used by the BLOCKHASH EVM op code.
	GetHashFunc func(uint64) common.Hash
)

// isProgramAccount returns true if the address is one of the following:
// - an address of precompiled contracts
// - an address of program accounts
func isProgramAccount(evm *EVM, caller common.Address, addr common.Address, db StateDB) bool {
	_, exists := evm.GetPrecompiledContractMap(caller)[addr]
	return exists || db.IsProgramAccount(addr)
}

// run runs the given contract and takes care of running precompiles with a fallback to the byte code interpreter.
func run(evm *EVM, contract *Contract, input []byte) ([]byte, error) {
	if contract.CodeAddr != nil {
		precompiles := evm.GetPrecompiledContractMap(contract.CallerAddress)
		if p := precompiles[*contract.CodeAddr]; p != nil {
			///////////////////////////////////////////////////////
			// OpcodeComputationCostLimit: The below code is commented and will be usd for debugging purposes.
			//var startTime time.Time
			//if opDebug {
			//	startTime = time.Now()
			//}
			///////////////////////////////////////////////////////
			ret, computationCost, err := RunPrecompiledContract(p, input, contract, evm) // TODO-Klaytn-Issue615
			///////////////////////////////////////////////////////
			// OpcodeComputationCostLimit: The below code is commented and will be usd for debugging purposes.
			//if opDebug {
			//	//fmt.Println("running precompiled contract...", "addr", contract.CodeAddr.String(), "computationCost", computationCost)
			//	elapsedTime := uint64(time.Since(startTime).Nanoseconds())
			//	addr := int(contract.CodeAddr.Bytes()[19])
			//	precompiledCnt[addr] += 1
			//	precompiledTime[addr] += elapsedTime
			//}
			///////////////////////////////////////////////////////
			evm.opcodeComputationCostSum += computationCost
			return ret, err
		}
	}
	return evm.interpreter.Run(contract, input)
}

// Context provides the EVM with auxiliary information. Once provided
// it shouldn't be modified.
type Context struct {
	// CanTransfer returns whether the account contains
	// sufficient KLAY to transfer the value
	CanTransfer CanTransferFunc
	// Transfer transfers KLAY from one account to the other
	Transfer TransferFunc
	// GetHash returns the hash corresponding to n
	GetHash GetHashFunc

	// Message information
	Origin   common.Address // Provides information for ORIGIN
	GasPrice *big.Int       // Provides information for GASPRICE

	// Block information
	Coinbase    common.Address // Provides information for COINBASE
	GasLimit    uint64         // Provides information for GASLIMIT
	BlockNumber *big.Int       // Provides information for NUMBER
	Time        *big.Int       // Provides information for TIME
	BlockScore  *big.Int       // Provides information for DIFFICULTY
	BaseFee     *big.Int       // Provides information for BASEFEE
}

// EVM is the Ethereum Virtual Machine base object and provides
// the necessary tools to run a contract on the given state with
// the provided context. It should be noted that any error
// generated through any of the calls should be considered a
// revert-state-and-consume-all-gas operation, no checks on
// specific errors should ever be performed. The interpreter makes
// sure that any errors generated are to be considered faulty code.
//
// The EVM should never be reused and is not thread safe.
type EVM struct {
	// Context provides auxiliary blockchain related information
	Context
	// StateDB gives access to the underlying state
	StateDB StateDB
	// Depth is the current call stack
	depth int

	// chainConfig contains information about the current chain
	chainConfig *params.ChainConfig
	// chain rules contains the chain rules for the current epoch
	chainRules params.Rules
	// virtual machine configuration options used to initialise the
	// evm.
	vmConfig *Config
	// global (to this context) ethereum virtual machine
	// used throughout the execution of the tx.
	interpreter *Interpreter
	// abort is used to abort the EVM calling operations
	// NOTE: must be set atomically
	abort int32
	// callGasTemp holds the gas available for the current call. This is needed because the
	// available gas is calculated in gasCall* according to the 63/64 rule and later
	// applied in opCall*.
	callGasTemp uint64

	// opcodeComputationCostSum is the sum of computation cost of opcodes.
	opcodeComputationCostSum uint64
}

// NewEVM returns a new EVM. The returned EVM is not thread safe and should
// only ever be used *once*.
func NewEVM(ctx Context, statedb StateDB, chainConfig *params.ChainConfig, vmConfig *Config) *EVM {
	evm := &EVM{
		Context:     ctx,
		StateDB:     statedb,
		vmConfig:    vmConfig,
		chainConfig: chainConfig,
		chainRules:  chainConfig.Rules(ctx.BlockNumber),
	}

	if vmConfig.RunningEVM != nil {
		vmConfig.RunningEVM <- evm
	}

	// If internal transaction tracing is enabled, creates a tracer for a transaction
	if vmConfig.EnableInternalTxTracing {
		vmConfig.Debug = true
		vmConfig.Tracer = NewInternalTxTracer()
	}

	evm.interpreter = NewEVMInterpreter(evm, vmConfig)
	return evm
}

// Cancel cancels any running EVM operation. This may be called concurrently and
// it's safe to be called multiple times.
func (evm *EVM) Cancel(reason int32) {
	for {
		abort := atomic.LoadInt32(&evm.abort)
		swapped := atomic.CompareAndSwapInt32(&evm.abort, abort, abort|reason)
		if swapped {
			break
		}
	}
}

func (evm *EVM) IsPrefetching() bool {
	return evm.vmConfig.Prefetching
}

// Cancelled returns true if Cancel has been called
func (evm *EVM) Cancelled() bool {
	return atomic.LoadInt32(&evm.abort) == 1
}

// Call executes the contract associated with the addr with the given input as
// parameters. It also handles any necessary value transfer required and takes
// the necessary steps to create accounts and reverses the state in case of an
// execution error or failed value transfer.
func (evm *EVM) Call(caller types.ContractRef, addr common.Address, input []byte, gas uint64, value *big.Int) (ret []byte, leftOverGas uint64, err error) {
	if evm.vmConfig.NoRecursion && evm.depth > 0 {
		return nil, gas, nil
	}

	// Fail if we're trying to execute above the call depth limit
	if evm.depth > int(params.CallCreateDepth) {
		return nil, gas, ErrDepth // TODO-Klaytn-Issue615
	}
	// Fail if we're trying to transfer more than the available balance
	if !evm.Context.CanTransfer(evm.StateDB, caller.Address(), value) {
		return nil, gas, ErrInsufficientBalance // TODO-Klaytn-Issue615
	}

	var (
		to       = AccountRef(addr)
		snapshot = evm.StateDB.Snapshot()
	)

	// Filter out invalid precompiled address calls, and create a precompiled contract object if it is not exist.
	if common.IsPrecompiledContractAddress(addr) {
		precompiles := evm.GetPrecompiledContractMap(caller.Address())
		if precompiles[addr] == nil || value.Sign() != 0 {
			// Return an error if an enabled precompiled address is called or a value is transferred to a precompiled address.
			if evm.vmConfig.Debug && evm.depth == 0 {
				evm.vmConfig.Tracer.CaptureStart(caller.Address(), addr, false, input, gas, value)
				evm.vmConfig.Tracer.CaptureEnd(ret, 0, 0, nil)
			}
			return nil, gas, kerrors.ErrPrecompiledContractAddress
		}
		// create an account object of the enabled precompiled address if not exist.
		if !evm.StateDB.Exist(addr) {
			evm.StateDB.CreateSmartContractAccount(addr, params.CodeFormatEVM, evm.chainRules)
		}
	}

	// The logic below creates an EOA account if not exist.
	// However, it does not create a contract account since `Call` is not proper method to create a contract.
	if !evm.StateDB.Exist(addr) {
		if value.Sign() == 0 {
			// Calling a non-existing account (probably contract), don't do anything, but ping the tracer
			if evm.vmConfig.Debug && evm.depth == 0 {
				evm.vmConfig.Tracer.CaptureStart(caller.Address(), addr, false, input, gas, value)
				evm.vmConfig.Tracer.CaptureEnd(ret, 0, 0, nil)
			}
			return nil, gas, nil
		}
		// If non-existing address is called with a value, an object of the address is created.
		evm.StateDB.CreateEOA(addr, false, accountkey.NewAccountKeyLegacy())
	}
	evm.Transfer(evm.StateDB, caller.Address(), to.Address(), value)

	if !isProgramAccount(evm, caller.Address(), addr, evm.StateDB) {
		return ret, gas, nil
	}

	// Initialise a new contract and set the code that is to be used by the EVM.
	// The contract is a scoped environment for this execution context only.
	contract := NewContract(caller, to, value, gas)
	contract.SetCallCode(&addr, evm.StateDB.GetCodeHash(addr), evm.StateDB.GetCode(addr))

	start := time.Now()

	// Capture the tracer start/end events in debug mode
	if evm.vmConfig.Debug && evm.depth == 0 {
		evm.vmConfig.Tracer.CaptureStart(caller.Address(), addr, false, input, gas, value)

		defer func() { // Lazy evaluation of the parameters
			evm.vmConfig.Tracer.CaptureEnd(ret, gas-contract.Gas, time.Since(start), err)
		}()
	}
	ret, err = run(evm, contract, input)

	// When an error was returned by the EVM or when setting the creation code
	// above we revert to the snapshot and consume any gas remaining. Additionally
	// when we're in homestead this also counts for code storage gas errors.
	if err != nil {
		evm.StateDB.RevertToSnapshot(snapshot)
		if err != ErrExecutionReverted {
			contract.UseGas(contract.Gas)
		}
	}
	return ret, contract.Gas, err
}

// CallCode executes the contract associated with the addr with the given input
// as parameters. It also handles any necessary value transfer required and takes
// the necessary steps to create accounts and reverses the state in case of an
// execution error or failed value transfer.
//
// CallCode differs from Call in the sense that it executes the given address'
// code with the caller as context.
func (evm *EVM) CallCode(caller types.ContractRef, addr common.Address, input []byte, gas uint64, value *big.Int) (ret []byte, leftOverGas uint64, err error) {
	if evm.vmConfig.NoRecursion && evm.depth > 0 {
		return nil, gas, nil
	}

	// Fail if we're trying to execute above the call depth limit
	if evm.depth > int(params.CallCreateDepth) {
		return nil, gas, ErrDepth // TODO-Klaytn-Issue615
	}
	// Fail if we're trying to transfer more than the available balance
	if !evm.CanTransfer(evm.StateDB, caller.Address(), value) {
		return nil, gas, ErrInsufficientBalance // TODO-Klaytn-Issue615
	}

	if !isProgramAccount(evm, caller.Address(), addr, evm.StateDB) {
		logger.Info("Returning since the addr is not a program account", "addr", addr)
		return nil, gas, nil
	}

	var (
		snapshot = evm.StateDB.Snapshot()
		to       = AccountRef(caller.Address())
	)
	// Initialise a new contract and set the code that is to be used by the EVM.
	// The contract is a scoped environment for this execution context only.
	contract := NewContract(caller, to, value, gas)
	contract.SetCallCode(&addr, evm.StateDB.GetCodeHash(addr), evm.StateDB.GetCode(addr))

	ret, err = run(evm, contract, input)
	if err != nil {
		evm.StateDB.RevertToSnapshot(snapshot)
		if err != ErrExecutionReverted {
			contract.UseGas(contract.Gas)
		}
	}
	return ret, contract.Gas, err
}

// DelegateCall executes the contract associated with the addr with the given input
// as parameters. It reverses the state in case of an execution error.
//
// DelegateCall differs from CallCode in the sense that it executes the given address'
// code with the caller as context and the caller is set to the caller of the caller.
func (evm *EVM) DelegateCall(caller types.ContractRef, addr common.Address, input []byte, gas uint64) (ret []byte, leftOverGas uint64, err error) {
	if evm.vmConfig.NoRecursion && evm.depth > 0 {
		return nil, gas, nil
	}
	// Fail if we're trying to execute above the call depth limit
	if evm.depth > int(params.CallCreateDepth) {
		return nil, gas, ErrDepth // TODO-Klaytn-Issue615
	}

	if !isProgramAccount(evm, caller.Address(), addr, evm.StateDB) {
		logger.Info("Returning since the addr is not a program account", "addr", addr)
		return nil, gas, nil
	}

	var (
		snapshot = evm.StateDB.Snapshot()
		to       = AccountRef(caller.Address())
	)

	// Initialise a new contract and make initialise the delegate values
	contract := NewContract(caller, to, nil, gas).AsDelegate()
	contract.SetCallCode(&addr, evm.StateDB.GetCodeHash(addr), evm.StateDB.GetCode(addr))

	ret, err = run(evm, contract, input)
	if err != nil {
		evm.StateDB.RevertToSnapshot(snapshot)
		if err != ErrExecutionReverted {
			contract.UseGas(contract.Gas)
		}
	}
	return ret, contract.Gas, err
}

// StaticCall executes the contract associated with the addr with the given input
// as parameters while disallowing any modifications to the state during the call.
// Opcodes that attempt to perform such modifications will result in exceptions
// instead of performing the modifications.
func (evm *EVM) StaticCall(caller types.ContractRef, addr common.Address, input []byte, gas uint64) (ret []byte, leftOverGas uint64, err error) {
	if evm.vmConfig.NoRecursion && evm.depth > 0 {
		return nil, gas, nil
	}
	// Fail if we're trying to execute above the call depth limit
	if evm.depth > int(params.CallCreateDepth) {
		return nil, gas, ErrDepth // TODO-Klaytn-Issue615
	}
	// Make sure the readonly is only set if we aren't in readonly yet
	// this makes also sure that the readonly flag isn't removed for
	// child calls.
	if !evm.interpreter.readOnly {
		evm.interpreter.readOnly = true
		defer func() { evm.interpreter.readOnly = false }()
	}

	if !isProgramAccount(evm, caller.Address(), addr, evm.StateDB) {
		logger.Info("Returning since the addr is not a program account", "addr", addr)
		return nil, gas, nil
	}

	var (
		to       = AccountRef(addr)
		snapshot = evm.StateDB.Snapshot()
	)
	// Initialise a new contract and set the code that is to be used by the EVM.
	// The contract is a scoped environment for this execution context only.
	contract := NewContract(caller, to, new(big.Int), gas)
	contract.SetCallCode(&addr, evm.StateDB.GetCodeHash(addr), evm.StateDB.GetCode(addr))

	// When an error was returned by the EVM or when setting the creation code
	// above we revert to the snapshot and consume any gas remaining. Additionally
	// when we're in Homestead this also counts for code storage gas errors.
	ret, err = run(evm, contract, input)
	if err != nil {
		evm.StateDB.RevertToSnapshot(snapshot)
		if err != ErrExecutionReverted {
			contract.UseGas(contract.Gas)
		}
	}
	return ret, contract.Gas, err
}

type codeAndHash struct {
	code []byte
	hash common.Hash
}

func (c *codeAndHash) Hash() common.Hash {
	if c.hash == (common.Hash{}) {
		c.hash = crypto.Keccak256Hash(c.code)
	}
	return c.hash
}

// Create creates a new contract using code as deployment code.
func (evm *EVM) create(caller types.ContractRef, codeAndHash *codeAndHash, gas uint64, value *big.Int, address common.Address, humanReadable bool, codeFormat params.CodeFormat) (ret []byte, contractAddr common.Address, leftOverGas uint64, err error) {

	// Depth check execution. Fail if we're trying to execute above the
	// limit.
	if evm.depth > int(params.CallCreateDepth) {
		return nil, common.Address{}, gas, ErrDepth // TODO-Klaytn-Issue615
	}
	if !evm.CanTransfer(evm.StateDB, caller.Address(), value) {
		return nil, common.Address{}, gas, ErrInsufficientBalance // TODO-Klaytn-Issue615
	}

	// Increasing nonce since a failed tx with one of following error will be loaded on a block.
	evm.StateDB.IncNonce(caller.Address())

	if evm.StateDB.Exist(address) {
		return nil, common.Address{}, 0, ErrContractAddressCollision // TODO-Klaytn-Issue615
	}
	if common.IsPrecompiledContractAddress(address) {
		return nil, common.Address{}, gas, kerrors.ErrPrecompiledContractAddress
	}

	// Create a new account on the state
	snapshot := evm.StateDB.Snapshot()
	// TODO-Klaytn-Accounts: for now, smart contract accounts cannot withdraw KLAYs via ValueTransfer
	//   because the account key is set to AccountKeyFail by default.
	//   Need to make a decision of the key type.
	evm.StateDB.CreateSmartContractAccountWithKey(address, humanReadable, accountkey.NewAccountKeyFail(), codeFormat, evm.chainRules)
	evm.StateDB.SetNonce(address, 1)
	if value.Sign() != 0 {
		evm.Transfer(evm.StateDB, caller.Address(), address, value)
	}
	// Initialise a new contract and set the code that is to be used by the EVM.
	// The contract is a scoped environment for this execution context only.
	contract := NewContract(caller, AccountRef(address), value, gas)
	contract.SetCodeOptionalHash(&address, codeAndHash)

	if evm.vmConfig.NoRecursion && evm.depth > 0 {
		return nil, address, gas, nil
	}

	if evm.vmConfig.Debug && evm.depth == 0 {
		evm.vmConfig.Tracer.CaptureStart(caller.Address(), address, true, codeAndHash.code, gas, value)
	}
	start := time.Now()

	ret, err = run(evm, contract, nil)

	// check whether the max code size has been exceeded
	maxCodeSizeExceeded := len(ret) > params.MaxCodeSize
	// if the contract creation ran successfully and no errors were returned
	// calculate the gas required to store the code. If the code could not
	// be stored due to not enough gas set an error and let it be handled
	// by the error checking condition below.
	if err == nil && !maxCodeSizeExceeded {
		createDataGas := uint64(len(ret)) * params.CreateDataGas
		if contract.UseGas(createDataGas) {
			if evm.StateDB.SetCode(address, ret) != nil {
				// `err` is returned to `vmerr` in `StateTransition.TransitionDb()`.
				// Then, `vmerr` will be used to make a receipt status using `getReceiptStatusFromVMerr()`.
				// Since `getReceiptStatusFromVMerr()` uses a map to determine the receipt status,
				// this `err` should be an error variable declared in vm/errors.go.
				// TODO-Klaytn: Make a package of error variables containing all exported error variables.
				// After the above TODO-Klaytn is resolved, we can return the error returned by `SetCode()` directly.
				err = ErrFailedOnSetCode
			}
		} else {
			err = ErrCodeStoreOutOfGas // TODO-Klaytn-Issue615
		}
	}

	// When an error was returned by the EVM or when setting the creation code
	// above we revert to the snapshot and consume any gas remaining.
	if maxCodeSizeExceeded || err != nil {
		evm.StateDB.RevertToSnapshot(snapshot)
		if err != ErrExecutionReverted {
			contract.UseGas(contract.Gas)
		}
	}
	// Assign err if contract code size exceeds the max while the err is still empty.
	if maxCodeSizeExceeded && err == nil {
		err = ErrMaxCodeSizeExceeded // TODO-Klaytn-Issue615
	}
	if evm.vmConfig.Debug && evm.depth == 0 {
		evm.vmConfig.Tracer.CaptureEnd(ret, gas-contract.Gas, time.Since(start), err)
	}
	return ret, address, contract.Gas, err
}

// Create creates a new contract using code as deployment code.
func (evm *EVM) Create(caller types.ContractRef, code []byte, gas uint64, value *big.Int, codeFormat params.CodeFormat) (ret []byte, contractAddr common.Address, leftOverGas uint64, err error) {
	codeAndHash := &codeAndHash{code: code}
	contractAddr = crypto.CreateAddress(caller.Address(), evm.StateDB.GetNonce(caller.Address()))
	return evm.create(caller, codeAndHash, gas, value, contractAddr, false, codeFormat)
}

// Create2 creates a new contract using code as deployment code.
//
// The different between Create2 with Create is Create2 uses sha3(0xff ++ msg.sender ++ salt ++ sha3(init_code))[12:]
// instead of the usual sender-and-nonce-hash as the address where the contract is initialized at.
func (evm *EVM) Create2(caller types.ContractRef, code []byte, gas uint64, endowment *big.Int, salt *big.Int, codeFormat params.CodeFormat) (ret []byte, contractAddr common.Address, leftOverGas uint64, err error) {
	codeAndHash := &codeAndHash{code: code}
	contractAddr = crypto.CreateAddress2(caller.Address(), common.BigToHash(salt), codeAndHash.Hash().Bytes())
	return evm.create(caller, codeAndHash, gas, endowment, contractAddr, false, codeFormat)
}

// CreateWithAddress creates a new contract using code as deployment code with given address and humanReadable.
func (evm *EVM) CreateWithAddress(caller types.ContractRef, code []byte, gas uint64, value *big.Int, contractAddr common.Address, humanReadable bool, codeFormat params.CodeFormat) ([]byte, common.Address, uint64, error) {
	codeAndHash := &codeAndHash{code: code}
	codeAndHash.Hash()
	return evm.create(caller, codeAndHash, gas, value, contractAddr, humanReadable, codeFormat)
}

func (evm *EVM) GetPrecompiledContractMap(addr common.Address) map[common.Address]PrecompiledContract {
	getPrecompiledContractMapWithVmVersion := func() (bool, map[common.Address]PrecompiledContract) {
		// Get vmVersion from addr. If there's no vmVersion, it returns false and use latest precompiled contract map
		vmVersion, ok := evm.StateDB.GetVmVersion(addr)
		if !ok {
			return false, nil
		}

		// Return precompiled contract map according to the VmVersion (use the map at deployment time of addr contract)
		//      (gas price policy also follows old map's rule)
		// If new "VmVersion" is added, add new if clause below
		if vmVersion == params.VmVersion0 {
			// Without this version, 0x09-0x0b won't work properly with contracts deployed before istanbulHF
			return true, PrecompiledContractsConstantinople
		}
		return false, nil
	}

	// There are contracts which uses latest precompiled contract map (regardless of deployment time)
	// If new HF is added, please add new case below
	switch {
	case evm.chainRules.IsLondon:
		fallthrough
	case evm.chainRules.IsIstanbul:
		if ok, mapWithVmVersion := getPrecompiledContractMapWithVmVersion(); ok {
			return mapWithVmVersion
		}
		return PrecompiledContractsIstanbul
	default:
		return PrecompiledContractsConstantinople
	}
}

// ChainConfig returns the environment's chain configuration
func (evm *EVM) ChainConfig() *params.ChainConfig { return evm.chainConfig }

// Interpreter returns the EVM interpreter
func (evm *EVM) Interpreter() *Interpreter { return evm.interpreter }

func (evm *EVM) GetOpCodeComputationCost() uint64 { return evm.opcodeComputationCostSum }
