package abi

import (
	"reflect"
	"math/big"
	"ground-x/go-gxplatform/common"
	"ground-x/go-gxplatform/common/math"
)

var (
	bigT      = reflect.TypeOf(&big.Int{})
	derefbigT = reflect.TypeOf(big.Int{})
	uint8T    = reflect.TypeOf(uint8(0))
	uint16T   = reflect.TypeOf(uint16(0))
	uint32T   = reflect.TypeOf(uint32(0))
	uint64T   = reflect.TypeOf(uint64(0))
	intT      = reflect.TypeOf(int(0))
	int8T     = reflect.TypeOf(int8(0))
	int16T    = reflect.TypeOf(int16(0))
	int32T    = reflect.TypeOf(int32(0))
	int64T    = reflect.TypeOf(int64(0))
	addressT  = reflect.TypeOf(common.Address{})
	intTS     = reflect.TypeOf([]int(nil))
	int8TS    = reflect.TypeOf([]int8(nil))
	int16TS   = reflect.TypeOf([]int16(nil))
	int32TS   = reflect.TypeOf([]int32(nil))
	int64TS   = reflect.TypeOf([]int64(nil))
)

// U256 converts a big Int into a 256bit EVM number.
func U256(n *big.Int) []byte {
	return math.PaddedBigBytes(math.U256(n), 32)
}

// checks whether the given reflect value is signed. This also works for slices with a number type
func isSigned(v reflect.Value) bool {
	switch v.Type() {
	case intTS, int8TS, int16TS, int32TS, int64TS, intT, int8T, int16T, int32T, int64T:
		return true
	}
	return false
}
