// Copyright 2018 The klaytn Authors
// This file is part of the klaytn library.
//
// The klaytn library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The klaytn library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the klaytn library. If not, see <http://www.gnu.org/licenses/>.

package benchmarks

import (
	"testing"

	"github.com/klaytn/klaytn/common"
)

func TestInterpreterMload100000(t *testing.T) {
	//
	// Test code
	//       Initialize memory with memory write (PUSH PUSH MSTORE)
	//       Loop 10000 times for below code
	//              memory read 10 times //  (PUSH MLOAD POP) x 10
	//
	code := common.Hex2Bytes(mload100000)
	intrp, contract := prepareInterpreterAndContract(code)

	intrp.Run(contract, nil)
}

func BenchmarkInterpreterMload100000(bench *testing.B) {
	//
	// Test code
	//       Initialize memory with memory write (PUSH PUSH MSTORE)
	//       Loop 10000 times for below code
	//              memory read 10 times //  (PUSH MLOAD POP) x 10
	//
	code := common.Hex2Bytes("60ca60205260005b612710811015630000004557602051506020515060205150602051506020515060205150602051506020515060205150602051506001016300000007565b00")
	intrp, contract := prepareInterpreterAndContract(code)

	bench.ResetTimer()
	for i := 0; i < bench.N; i++ {
		intrp.Run(contract, nil)
	}
}

func TestInterpreterMstore100000(t *testing.T) {
	//
	// Test code
	//       Initialize memory with memory write (PUSH PUSH MSTORE)
	//       Loop 10000 times for below code
	//              memory write 10 times //  (PUSH PUSH MSTORE) x 10
	//
	code := common.Hex2Bytes(mstore100000)
	intrp, contract := prepareInterpreterAndContract(code)

	intrp.Run(contract, nil)
}

func BenchmarkInterpreterMstore100000(bench *testing.B) {
	//
	// Test code
	//       Initialize memory with memory write (PUSH PUSH MSTORE)
	//       Loop 10000 times for below code
	//              memory write 10 times //  (PUSH PUSH MSTORE) x 10
	//
	code := common.Hex2Bytes(mstore100000)
	intrp, contract := prepareInterpreterAndContract(code)

	bench.ResetTimer()
	for i := 0; i < bench.N; i++ {
		intrp.Run(contract, nil)
	}
}

func TestInterpreterSload100000(t *testing.T) {
	//
	// Test code
	//       Initialize (PUSH SSTORE)
	//       Loop 10000 times for below code
	//              Read from database 10 times //  (PUSH SLOAD POP) x 10
	//
	code := common.Hex2Bytes(sload100000)
	intrp, contract := prepareInterpreterAndContract(code)

	intrp.Run(contract, nil)
}

func BenchmarkInterpreterSload100000(bench *testing.B) {
	//
	// Test code
	//       Initialize (PUSH SSTORE)
	//       Loop 10000 times for below code
	//              Read from database 10 times //  (PUSH SLOAD POP) x 10
	//
	code := common.Hex2Bytes(sload100000)
	intrp, contract := prepareInterpreterAndContract(code)

	bench.ResetTimer()
	for i := 0; i < bench.N; i++ {
		intrp.Run(contract, nil)
	}
}

func TestInterpreterSstoreNonZero2NonZero100000(t *testing.T) {
	//
	// Test code
	//       Initialize (PUSH)
	//       Loop 10000 times for below code
	//              Write to database 10 times //  (PUSH PUSH SSTORE) x 10
	//
	code := common.Hex2Bytes(sstoreNonZero2NonZero100000)
	intrp, contract := prepareInterpreterAndContract(code)

	intrp.Run(contract, nil)
}

func BenchmarkInterpreterSstoreZero2Zero100000(bench *testing.B) {
	//
	// Test code
	//       Initialize (PUSH)
	//       Loop 10000 times for below code
	//              Write to database 10 times //  (PUSH PUSH SSTORE) x 10
	//
	code := common.Hex2Bytes(sstoreZero2Zero100000)
	intrp, contract := prepareInterpreterAndContract(code)

	bench.ResetTimer()
	for i := 0; i < bench.N; i++ {
		intrp.Run(contract, nil)
	}
}

func BenchmarkInterpreterSstoreNonZero2NonZero100000(bench *testing.B) {
	//
	// Test code
	//       Initialize (PUSH)
	//       Loop 10000 times for below code
	//              Write to database 10 times //  (PUSH PUSH SSTORE) x 10
	//
	code := common.Hex2Bytes(sstoreNonZero2NonZero100000)
	intrp, contract := prepareInterpreterAndContract(code)

	bench.ResetTimer()
	for i := 0; i < bench.N; i++ {
		intrp.Run(contract, nil)
	}
}

func BenchmarkInterpreterSstoreMixed100000(bench *testing.B) {
	//
	// Test code
	//       Initialize (PUSH)
	//       Loop 10000 times for below code
	//              Write to database 10 times //  (PUSH PUSH SSTORE) x 10
	//
	code := common.Hex2Bytes(sstoreMixed100000)
	intrp, contract := prepareInterpreterAndContract(code)

	bench.ResetTimer()
	for i := 0; i < bench.N; i++ {
		intrp.Run(contract, nil)
	}
}

func TestInterpreterAdd100000(t *testing.T) {
	//
	// Test code
	//       Initialize value (PUSH)
	//       Loop 10000 times for below code
	//              initialize value and increment 10 times //  (PUSH x 1) + ((PUSH ADD) x 10)
	//
	code := common.Hex2Bytes(add100000)
	intrp, contract := prepareInterpreterAndContract(code)

	intrp.Run(contract, nil)
}

func BenchmarkInterpreterAdd100000(bench *testing.B) {
	//
	// Test code
	//       Initialize value (PUSH)
	//       Loop 10000 times for below code
	//              initialize value and increment 10 times //  (PUSH x 1) + ((PUSH ADD) x 10)
	//
	code := common.Hex2Bytes(add100000)
	intrp, contract := prepareInterpreterAndContract(code)

	bench.ResetTimer()
	for i := 0; i < bench.N; i++ {
		intrp.Run(contract, nil)
	}
}

func TestInterpreterPush1Mul1byte100000(t *testing.T) {
	//
	// Test code
	//       Initialize value (PUSH)
	//       Loop 10000 times for below code
	//              initialize value and increment 10 times //  (PUSH1 PUSH1 MUL POP) x 10
	//
	code := common.Hex2Bytes(push1mul1byte100000)
	intrp, contract := prepareInterpreterAndContract(code)

	intrp.Run(contract, nil)
}

func BenchmarkInterpreterPush1Mul1byte100000(bench *testing.B) {
	//
	// Test code
	//       Initialize value (PUSH)
	//       Loop 10000 times for below code
	//              initialize value and increment 10 times //  (PUSH1 PUSH1 MUL POP) x 10
	//
	code := common.Hex2Bytes(push1mul1byte100000)
	intrp, contract := prepareInterpreterAndContract(code)

	bench.ResetTimer()
	for i := 0; i < bench.N; i++ {
		intrp.Run(contract, nil)
	}
}

func BenchmarkInterpreterPush5Mul1byte100000(bench *testing.B) {
	//
	// Test code
	//       Initialize value (PUSH)
	//       Loop 10000 times for below code
	//              initialize value and increment 10 times //  (PUSH5 PUSH5 MUL POP) x 10
	//
	code := common.Hex2Bytes(push5mul1byte100000)
	intrp, contract := prepareInterpreterAndContract(code)

	bench.ResetTimer()
	for i := 0; i < bench.N; i++ {
		intrp.Run(contract, nil)
	}
}

func BenchmarkInterpreterPush5Mul5bytes100000(bench *testing.B) {
	//
	// Test code
	//       Initialize value (PUSH)
	//       Loop 10000 times for below code
	//              initialize value and increment 10 times //  (PUSH5 PUSH5 MUL POP) x 10
	//
	code := common.Hex2Bytes(push5mul5bytes100000)
	intrp, contract := prepareInterpreterAndContract(code)

	bench.ResetTimer()
	for i := 0; i < bench.N; i++ {
		intrp.Run(contract, nil)
	}
}

func TestInterpreterPush1Div1byte100000(t *testing.T) {
	//
	// Test code
	//       Initialize value (PUSH)
	//       Loop 10000 times for below code
	//              initialize value and increment 10 times //  (PUSH1 PUSH1 DIV POP) x 10
	//
	code := common.Hex2Bytes(push1div1byte100000)
	intrp, contract := prepareInterpreterAndContract(code)

	intrp.Run(contract, nil)
}

func BenchmarkInterpreterPush1Div1byte100000(bench *testing.B) {
	//
	// Test code
	//       Initialize value (PUSH)
	//       Loop 10000 times for below code
	//              initialize value and increment 10 times //  (PUSH1 PUSH1 DIV POP) x 10
	//
	code := common.Hex2Bytes(push1div1byte100000)
	intrp, contract := prepareInterpreterAndContract(code)

	bench.ResetTimer()
	for i := 0; i < bench.N; i++ {
		intrp.Run(contract, nil)
	}
}

func BenchmarkInterpreterPush5Div1byte100000(bench *testing.B) {
	//
	// Test code
	//       Initialize value (PUSH)
	//       Loop 10000 times for below code
	//              initialize value and increment 10 times //  (PUSH5 PUSH5 DIV POP) x 10
	//
	code := common.Hex2Bytes(push5div1byte100000)
	intrp, contract := prepareInterpreterAndContract(code)

	bench.ResetTimer()
	for i := 0; i < bench.N; i++ {
		intrp.Run(contract, nil)
	}
}

func BenchmarkInterpreterPush5Div5bytes100000(bench *testing.B) {
	//
	// Test code
	//       Initialize value (PUSH)
	//       Loop 10000 times for below code
	//              initialize value and increment 10 times //  (PUSH5 PUSH5 DIV POP) x 10
	//
	code := common.Hex2Bytes(push5div5bytes100000)
	intrp, contract := prepareInterpreterAndContract(code)

	bench.ResetTimer()
	for i := 0; i < bench.N; i++ {
		intrp.Run(contract, nil)
	}
}
