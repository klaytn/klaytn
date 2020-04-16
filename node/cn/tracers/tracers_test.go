// Copyright 2018 The klaytn Authors
// Copyright 2017 The go-ethereum Authors
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
// This file is derived from eth/tracers/tracers_test.go (2018/06/04).
// Modified and improved for the klaytn development.

package tracers

import (
	"encoding/json"
	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/blockchain/vm"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/hexutil"
	"github.com/klaytn/klaytn/common/math"
	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/ser/rlp"
	"github.com/klaytn/klaytn/storage/database"
	"github.com/klaytn/klaytn/tests"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"math/big"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// To generate a new callTracer test, copy paste the makeTest method below into
// the klaytn console and call it with a transaction hash you which to export.

/*
// makeTest generates a callTracer test by running a prestate reassembled and a
// call trace run, assembling all the gathered information into a test case.
var makeTest = function(tx, rewind) {
  // Generate the genesis block from the block, transaction and prestate data
  var block   = klay.getBlock(klay.getTransaction(tx).blockHash);
  var genesis = klay.getBlock(block.parentHash);

  delete genesis.gasUsed;
  delete genesis.logsBloom;
  delete genesis.parentHash;
  delete genesis.receiptsRoot;
  delete genesis.size;
  delete genesis.transactions;
  delete genesis.transactionsRoot;

  genesis.gasLimit  = genesis.gasLimit.toString();
  genesis.number    = genesis.number.toString();
  genesis.timestamp = genesis.timestamp.toString();

  genesis.alloc = debug.traceTransaction(tx, {tracer: "prestateTracer", rewind: rewind});
  for (var key in genesis.alloc) {
    genesis.alloc[key].nonce = genesis.alloc[key].nonce.toString();
  }
  genesis.config = admin.nodeInfo.protocols.klay.config;

  // Generate the call trace and produce the test input
  var result = debug.traceTransaction(tx, {tracer: "callTracer", rewind: rewind});
  delete result.time;

  console.log(JSON.stringify({
    genesis: genesis,
    context: {
      number:     block.number.toString(),
      blockscore: block.blockscore,
      timestamp:  block.timestamp.toString(),
      gasLimit:   block.gasLimit.toString(),
      miner:      block.miner,
    },
    input:  klay.getRawTransaction(tx),
    result: result,
  }, null, 2));
}
*/

type reverted struct {
	Contract common.Address `json:"contract"`
	Message  string         `json:"message"`
}

// callTrace is the result of a callTracer run.
type callTrace struct {
	Type     string          `json:"type"`
	From     common.Address  `json:"from"`
	To       common.Address  `json:"to"`
	Input    hexutil.Bytes   `json:"input"`
	Output   hexutil.Bytes   `json:"output"`
	Gas      *hexutil.Uint64 `json:"gas,omitempty"`
	GasUsed  *hexutil.Uint64 `json:"gasUsed,omitempty"`
	Value    *hexutil.Big    `json:"value,omitempty"`
	Error    string          `json:"error,omitempty"`
	Calls    []callTrace     `json:"calls,omitempty"`
	Reverted reverted        `json:"reverted,omitempty"`
}

type callContext struct {
	Number     math.HexOrDecimal64   `json:"number"`
	BlockScore *math.HexOrDecimal256 `json:"blockScore"`
	Time       math.HexOrDecimal64   `json:"timestamp"`
	GasLimit   math.HexOrDecimal64   `json:"gasLimit"`
	Miner      common.Address        `json:"miner"`
}

// callTracerTest defines a single test to check the call tracer against.
type callTracerTest struct {
	Genesis     *blockchain.Genesis `json:"genesis"`
	Context     *callContext        `json:"context"`
	Input       string              `json:"input,omitempty"`
	Transaction map[string]string   `json:"transaction,omitempty"`
	Result      *callTrace          `json:"result"`
}

// Iterates over all the input-output datasets in the tracer test harness and
// runs the JavaScript tracers against them.
func TestCallTracer(t *testing.T) {
	files, err := ioutil.ReadDir("testdata")
	if err != nil {
		t.Fatalf("failed to retrieve tracer test suite: %v", err)
	}
	for _, file := range files {
		if !strings.HasPrefix(file.Name(), "call_tracer_") {
			continue
		}
		file := file // capture range variable
		t.Run(camel(strings.TrimSuffix(strings.TrimPrefix(file.Name(), "call_tracer_"), ".json")), func(t *testing.T) {
			t.Parallel()

			// Call tracer test found, read if from disk
			blob, err := ioutil.ReadFile(filepath.Join("testdata", file.Name()))
			if err != nil {
				t.Fatalf("failed to read testcase: %v", err)
			}
			test := new(callTracerTest)
			if err := json.Unmarshal(blob, test); err != nil {
				t.Fatalf("failed to parse testcase: %v", err)
			}

			signer := types.MakeSigner(test.Genesis.Config, new(big.Int).SetUint64(uint64(test.Context.Number)))
			tx := new(types.Transaction)
			// Configure a blockchain with the given prestate
			if test.Input != "" {
				if err := rlp.DecodeBytes(common.FromHex(test.Input), tx); err != nil {
					t.Fatalf("failed to parse testcase input: %v", err)
				}
			} else {
				// Configure a blockchain with the given prestate
				value := new(big.Int)
				gasPrice := new(big.Int)
				err = value.UnmarshalJSON([]byte(test.Transaction["value"]))
				require.NoError(t, err)
				err = gasPrice.UnmarshalJSON([]byte(test.Transaction["gasPrice"]))
				require.NoError(t, err)
				nonce, b := math.ParseUint64(test.Transaction["nonce"])
				require.True(t, b)
				gas, b := math.ParseUint64(test.Transaction["gas"])
				require.True(t, b)

				to := common.HexToAddress(test.Transaction["to"])
				input := common.FromHex(test.Transaction["input"])

				tx = types.NewTransaction(nonce, to, value, gas, gasPrice, input)

				testKey, err := crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
				require.NoError(t, err)
				err = tx.Sign(signer, testKey)
				require.NoError(t, err)
			}

			origin, _ := signer.Sender(tx)

			context := vm.Context{
				CanTransfer: blockchain.CanTransfer,
				Transfer:    blockchain.Transfer,
				Origin:      origin,
				BlockNumber: new(big.Int).SetUint64(uint64(test.Context.Number)),
				Time:        new(big.Int).SetUint64(uint64(test.Context.Time)),
				BlockScore:  (*big.Int)(test.Context.BlockScore),
				GasLimit:    uint64(test.Context.GasLimit),
				GasPrice:    tx.GasPrice(),
			}
			statedb := tests.MakePreState(database.NewMemoryDBManager(), test.Genesis.Alloc)

			// Create the tracer, the EVM environment and run it
			tracer, err := New("callTracer")
			if err != nil {
				t.Fatalf("failed to create call tracer: %v", err)
			}
			evm := vm.NewEVM(context, statedb, test.Genesis.Config, &vm.Config{Debug: true, Tracer: tracer})

			msg, err := tx.AsMessageWithAccountKeyPicker(signer, statedb, context.BlockNumber.Uint64())
			if err != nil {
				t.Fatalf("failed to prepare transaction for tracing: %v", err)
			}
			st := blockchain.NewStateTransition(evm, msg)
			if _, _, kerr := st.TransitionDb(); kerr.ErrTxInvalid != nil {
				t.Fatalf("failed to execute transaction: %v", kerr.ErrTxInvalid)
			}
			// Retrieve the trace result and compare against the etalon
			res, err := tracer.GetResult()
			if err != nil {
				t.Fatalf("failed to retrieve trace result: %v", err)
			}
			ret := new(callTrace)
			if err := json.Unmarshal(res, ret); err != nil {
				t.Fatalf("failed to unmarshal trace result: %v", err)
			}
			if !reflect.DeepEqual(ret, test.Result) {
				t.Fatalf("trace mismatch: have %+v, want %+v", ret, test.Result)
			}
		})
	}
}
