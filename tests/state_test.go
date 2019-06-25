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
// This file is derived from tests/state_test.go (2018/06/04).
// Modified and improved for the klaytn development.

package tests

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"

	"github.com/klaytn/klaytn/blockchain/vm"
)

func TestState(t *testing.T) {
	t.Parallel()

	st := new(testMatcher)
	// Long tests:
	st.skipShortMode(`^stQuadraticComplexityTest/`)
	// Broken tests:
	st.skipLoad(`^stTransactionTest/OverflowGasRequire\.json`) // gasLimit > 256 bits
	st.skipLoad(`^stTransactionTest/zeroSigTransa[^/]*\.json`) // EIP-86 is not supported yet
	// Expected failures:
	st.fails(`^stRevertTest/RevertPrecompiledTouch\.json/Byzantium`, "bug in test")
	st.skipLoad(`^stZeroKnowledge2/ecmul_0-3_5616_28000_96\.json`)

	// Skip since the tests transfer values to precompiled contracts
	st.skipLoad(`^stPreCompiledContracts2/CallSha256_1_nonzeroValue.json`)
	st.skipLoad(`^stPreCompiledContracts2/CallIdentity_1_nonzeroValue.json`)
	st.skipLoad(`^stPreCompiledContracts2/CallEcrecover0_NoGas.json`)
	st.skipLoad(`^stRandom2/randomStatetest644.json`)
	st.skipLoad(`^stRandom2/randomStatetest645.json`)
	st.skipLoad(`^stStaticCall/static_CallIdentity_1_nonzeroValue.json`)
	st.skipLoad(`^stStaticCall/static_CallSha256_1_nonzeroValue.json`)
	st.skipLoad(`^stArgsZeroOneBalance/callNonConst.json`)
	st.skipLoad(`^stPreCompiledContracts2/modexpRandomInput.json`)
	st.skipLoad(`^stRandom2/randomStatetest642.json`)

	st.walk(t, stateTestDir, func(t *testing.T, name string, test *StateTest) {
		for _, subtest := range test.Subtests() {
			subtest := subtest
			key := fmt.Sprintf("%s/%d", subtest.Fork, subtest.Index)
			name := name + "/" + key
			t.Run(key, func(t *testing.T) {
				if subtest.Fork == "Constantinople" {
					t.Skip("constantinople not supported yet")
				}
				withTrace(t, test.gasLimit(subtest), func(vmconfig vm.Config) error {
					_, err := test.Run(subtest, vmconfig)
					return st.checkFailure(t, name, err)
				})
			})
		}
	})
}

// Transactions with gasLimit above this value will not get a VM trace on failure.
const traceErrorLimit = 400000

func withTrace(t *testing.T, gasLimit uint64, test func(vm.Config) error) {
	err := test(vm.Config{})
	if err == nil {
		return
	}
	t.Error(err)
	if gasLimit > traceErrorLimit {
		t.Log("gas limit too high for EVM trace")
		return
	}
	tracer := vm.NewStructLogger(nil)
	err2 := test(vm.Config{Debug: true, Tracer: tracer})
	if !reflect.DeepEqual(err, err2) {
		t.Errorf("different error for second run: %v", err2)
	}
	buf := new(bytes.Buffer)
	vm.WriteTrace(buf, tracer.StructLogs())
	if buf.Len() == 0 {
		t.Log("no EVM operation logs generated")
	} else {
		t.Log("EVM operation log:\n" + buf.String())
	}
	t.Logf("EVM output: 0x%x", tracer.Output())
	t.Logf("EVM error: %v", tracer.Error())
}
