// Copyright 2018 The klaytn Authors
//
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

pragma solidity ^0.4.24;

contract StopContract {
    function Sstop() view {
        assembly {
            stop
        }
    }

    function ReturnData(uint size) returns (bytes) {
        return new bytes(size);
    }
}

contract OpCodeBenchmarkContract {

    mapping(uint => uint) public testMap;

    function Add(uint loopCount, uint x, uint y) returns (uint) {
        uint i;
        uint k;
        for (i = 0; i < loopCount; i++) {
            assembly {
                k := add(x,y)
            }
        }
        return k;
    }

    function Sub(uint loopCount, uint x, uint y) returns (uint) {
        uint i;
        uint k;
        for (i = 0; i < loopCount; i++) {
            assembly {
                k := sub(x,y)
            }
        }
        return k;
    }

    function Mul(uint loopCount, uint x, uint y) returns (uint) {
        uint i;
        uint k;
        for (i = 0; i < loopCount; i++) {
            assembly {
                k := mul(x,y)
            }
        }
        return k;
    }

    function Div(uint loopCount, uint x, uint y) returns (uint) {
        uint i;
        uint k = loopCount;
        for (i = 1; i < loopCount; i++) {
            assembly {
                k := div(x,y)
            }
        }
        return k;
    }

    function Sdiv(uint loopCount, uint x, uint y) returns (uint) {
        uint i;
        uint k = loopCount;
        for (i = 1; i < loopCount; i++) {
            assembly {
                k := sdiv(x,y)
            }
        }
        return k;
    }

    function Mod(uint loopCount, uint x, uint y) returns (uint) {
        uint i;
        uint k = loopCount;
        for (i = 1; i < loopCount; i++) {
            assembly {
                k := mod(x,y)
            }
        }
        return k;
    }

    function Smod(uint loopCount, uint x, uint y) returns (uint) {
        uint i;
        uint k = loopCount;
        for (i = 1; i < loopCount; i++) {
            assembly {
                k := smod(x,y)
            }
        }
        return k;
    }

    function Exp(uint loopCount, uint base, uint e) returns (uint) {
        uint i;
        uint k = loopCount;
        for (i = 1; i < loopCount; i++) {
            assembly {
                k := exp(base, e)
            }
        }
        return k;
    }

    function Not(uint loopCount) returns (uint) {
        uint i;
        uint k = loopCount;
        for (i = 1; i < loopCount; i++) {
            assembly {
                k := not(i)
            }
        }
        return k;
    }

    function Lt(uint loopCount) returns (uint) {
        uint i;
        uint k = loopCount;
        for (i = 1; i < loopCount; i++) {
            assembly {
                k := lt(k, i)
            }
        }
        return k;
    }

    function Gt(uint loopCount) returns (uint) {
        uint i;
        uint k = loopCount;
        for (i = 1; i < loopCount; i++) {
            assembly {
                k := gt(k, i)
            }
        }
        return k;
    }

    function Slt(uint loopCount) returns (uint) {
        uint i;
        uint k = loopCount;
        for (i = 1; i < loopCount; i++) {
            assembly {
                k := slt(k, i)
            }
        }
        return k;
    }

    function Sgt(uint loopCount) returns (uint) {
        uint i;
        uint k = loopCount;
        for (i = 1; i < loopCount; i++) {
            assembly {
                k := sgt(k, i)
            }
        }
        return k;
    }

    function Eq(uint loopCount) returns (uint) {
        uint i;
        uint k = loopCount;
        for (i = 1; i < loopCount; i++) {
            assembly {
                k := eq(k, i)
            }
        }
        return k;
    }

    function Iszero(uint loopCount) returns (uint) {
        uint i;
        uint k = loopCount;
        for (i = 1; i < loopCount; i++) {
            assembly {
                k := iszero(i)
            }
        }
        return k;
    }

    function And(uint loopCount) returns (uint) {
        uint i;
        uint k = loopCount;
        for (i = 1; i < loopCount; i++) {
            assembly {
                k := and(k, i)
            }
        }
        return k;
    }

    function Or(uint loopCount) returns (uint) {
        uint i;
        uint k = loopCount;
        for (i = 1; i < loopCount; i++) {
            assembly {
                k := or(k, i)
            }
        }
        return k;
    }

    function Xor(uint loopCount) returns (uint) {
        uint i;
        uint k = loopCount;
        for (i = 1; i < loopCount; i++) {
            assembly {
                k := xor(k, i)
            }
        }
        return k;
    }

    function Byte(uint loopCount) returns (uint) {
        uint i;
        uint k = loopCount;
        for (i = 1; i < loopCount; i++) {
            assembly {
                k := byte(1, i)
            }
        }
        return k;
    }

    function Shl(uint loopCount, uint x, uint y) returns (uint) {
        uint i;
        uint k = loopCount;
        for (i = 1; i < loopCount; i++) {
            assembly {
                k := shl(x, y)
            }
        }
        return k;
    }

    function Shr(uint loopCount, uint x, uint y) returns (uint) {
        uint i;
        uint k = loopCount;
        for (i = 1; i < loopCount; i++) {
            assembly {
                k := shr(x, y)
            }
        }
        return k;
    }

    function Sar(uint loopCount, uint x, uint y) returns (uint) {
        uint i;
        uint k = loopCount;
        for (i = 1; i < loopCount; i++) {
            assembly {
                k := sar(x, y)
            }
        }
        return k;
    }

    function Addmod(uint loopCount, uint x, uint y, uint m) returns (uint) {
        uint i;
        uint k = loopCount;
        for (i = 1; i < loopCount; i++) {
            assembly {
                k := addmod(x, y, m)
            }
        }
        return k;
    }

    function Mulmod(uint loopCount, uint x, uint y, uint m) returns (uint) {
        uint i;
        uint k = loopCount;
        for (i = 1; i < loopCount; i++) {
            assembly {
                k := mulmod(x, y, m)
            }
        }
        return k;
    }

    function SignExtend(uint loopCount, uint x, uint y) returns (uint) {
        uint i;
        uint k = loopCount;
        for (i = 1; i < loopCount; i++) {
            assembly {
                k := signextend(x, y)
            }
        }
        return k;
    }

    function Sha3(uint loopCount) returns (uint) {
        uint i;
        uint k = loopCount;
        for (i = 1; i < loopCount; i++) {
            assembly {
                k := sha3(i, k)
            }
        }
        return k;
    }

    function Pc(uint loopCount) returns (uint) {
        uint i;
        uint k = loopCount;
        for (i = 1; i < loopCount; i++) {
            assembly {
                k := pc()
            }
        }
        return k;
    }

    function Dup1(uint loopCount) returns (uint) {
        uint i;
        uint k = loopCount;
        for (i = 1; i < loopCount; i++) {
            assembly {
                dup1
                pop
            }
        }
        return k;
    }

    function Dup2(uint loopCount) returns (uint) {
        uint i;
        uint k = loopCount;
        for (i = 1; i < loopCount; i++) {
            assembly {
                dup2
                pop
            }
        }
        return k;
    }

    function Dup3(uint loopCount) returns (uint) {
        uint i;
        uint k = loopCount;
        for (i = 1; i < loopCount; i++) {
            assembly {
                dup3
                pop
            }
        }
        return k;
    }

    function Dup4(uint loopCount) returns (uint) {
        uint i;
        uint k = loopCount;
        for (i = 1; i < loopCount; i++) {
            assembly {
                dup4 pop
            }
        }
        return k;
    }

    function Dup5(uint loopCount) returns (uint) {
        uint i;
        uint k = loopCount;
        for (i = 1; i < loopCount; i++) {
            assembly {
                dup5 pop
            }
        }
        return k;
    }

    function Dup6(uint loopCount) returns (uint) {
        uint i;
        uint k = loopCount;
        for (i = 1; i < loopCount; i++) {
            assembly {
                dup6 pop
            }
        }
        return k;
    }

    function Dup7(uint loopCount) returns (uint) {
        uint i;
        uint k = loopCount;
        for (i = 1; i < loopCount; i++) {
            assembly {
                dup7 pop
            }
        }
        return k;
    }

    function Dup8(uint loopCount) returns (uint) {
        uint i;
        uint k = loopCount;
        for (i = 1; i < loopCount; i++) {
            assembly {
                dup8 pop
            }
        }
        return k;
    }

    function Dup9(uint loopCount) returns (uint) {
        uint i;
        uint k = loopCount;
        for (i = 1; i < loopCount; i++) {
            assembly {
                dup9 pop
            }
        }
        return k;
    }

    function Dup10(uint loopCount) returns (uint) {
        uint i;
        uint k = loopCount;
        for (i = 1; i < loopCount; i++) {
            assembly {
                dup10 pop
            }
        }
        return k;
    }

    function Dup11(uint loopCount) returns (uint) {
        uint i;
        uint k = loopCount;
        for (i = 1; i < loopCount; i++) {
            assembly {
                dup11 pop
            }
        }
        return k;
    }

    function Dup12(uint loopCount) returns (uint) {
        uint i;
        uint k = loopCount;
        for (i = 1; i < loopCount; i++) {
            assembly {
                dup12 pop
            }
        }
        return k;
    }

    function Dup13(uint loopCount) returns (uint) {
        uint i;
        uint k = loopCount;
        for (i = 1; i < loopCount; i++) {
            assembly {
                dup13 pop
            }
        }
        return k;
    }

    function Dup14(uint loopCount) returns (uint) {
        uint i;
        uint k = loopCount;
        for (i = 1; i < loopCount; i++) {
            assembly {
                dup14 pop
            }
        }
        return k;
    }

    function Dup15(uint loopCount) returns (uint) {
        uint i;
        uint k = loopCount;
        for (i = 1; i < loopCount; i++) {
            assembly {
                dup15 pop
            }
        }
        return k;
    }

    function Dup16(uint loopCount) returns (uint) {
        uint i;
        uint k = loopCount;
        for (i = 1; i < loopCount; i++) {
            assembly {
                dup16 pop
            }
        }
        return k;
    }

    function Swap1(uint loopCount) returns (uint) {
        uint i;
        uint k = loopCount;
        for (i = 1; i < loopCount; i++) {
            assembly {
                swap1
                swap1
            }
        }
        return k;
    }

    function Swap2(uint loopCount) returns (uint) {
        uint i;
        uint k = loopCount;
        for (i = 1; i < loopCount; i++) {
            assembly {
                swap2
                swap2
            }
        }
        return k;
    }

    function Swap3(uint loopCount) returns (uint) {
        uint i;
        uint k = loopCount;
        for (i = 1; i < loopCount; i++) {
            assembly {
                swap3 swap3
            }
        }
        return k;
    }

    function Swap4(uint loopCount) returns (uint) {
        uint i;
        uint k = loopCount;
        for (i = 1; i < loopCount; i++) {
            assembly {
                swap4 swap4
            }
        }
        return k;
    }

    function Swap5(uint loopCount) returns (uint) {
        uint i;
        uint k = loopCount;
        for (i = 1; i < loopCount; i++) {
            assembly { swap5 swap5 }
        }
        return k;
    }

    function Swap6(uint loopCount) returns (uint) {
        uint i;
        uint k = loopCount;
        for (i = 1; i < loopCount; i++) {
            assembly { swap6 swap6 }
        }
        return k;
    }

    function Swap7(uint loopCount) returns (uint) {
        uint i;
        uint k = loopCount;
        for (i = 1; i < loopCount; i++) {
            assembly { swap7 swap7 }
        }
        return k;
    }

    function Swap8(uint loopCount) returns (uint) {
        uint i;
        uint k = loopCount;
        for (i = 1; i < loopCount; i++) {
            assembly { swap8 swap8 }
        }
        return k;
    }

    function Swap9(uint loopCount) returns (uint) {
        uint i;
        uint k = loopCount;
        for (i = 1; i < loopCount; i++) {
            assembly { swap9 swap9 }
        }
        return k;
    }

    function Swap10(uint loopCount) returns (uint) {
        uint i;
        uint k = loopCount;
        for (i = 1; i < loopCount; i++) {
            assembly { swap10 swap10 }
        }
        return k;
    }

    function Swap11(uint loopCount) returns (uint) {
        uint i;
        uint k = loopCount;
        for (i = 1; i < loopCount; i++) {
            assembly { swap11 swap11 }
        }
        return k;
    }

    function Swap12(uint loopCount) returns (uint) {
        uint i;
        uint k = loopCount;
        for (i = 1; i < loopCount; i++) {
            assembly { swap12 swap12 }
        }
        return k;
    }

    function Swap13(uint loopCount) returns (uint) {
        uint i;
        uint k = loopCount;
        for (i = 1; i < loopCount; i++) {
            assembly { swap13 swap13 }
        }
        return k;
    }

    function Swap14(uint loopCount) returns (uint) {
        uint i;
        uint k = loopCount;
        for (i = 1; i < loopCount; i++) {
            assembly { swap14 swap14 }
        }
        return k;
    }

    function Swap15(uint loopCount) returns (uint) {
        uint i;
        uint k = loopCount;
        for (i = 1; i < loopCount; i++) {
            assembly { swap15 swap15 }
        }
        return k;
    }

    function Swap16(uint loopCount) returns (uint) {
        uint i;
        uint k = loopCount;
        for (i = 1; i < loopCount; i++) {
            assembly { swap16 swap16 }
        }
        return k;
    }

    function Call(uint loopCount, address to) returns (uint) {
        uint i;
        uint k;
        for(i = 0; i < loopCount; i++) {
            assembly {
                k := call(gas, to, 0, 0, 0, 0, 0)
            }
        }
        return i;
    }

    function CallCode(uint loopCount, address to) returns (uint) {
        uint i;
        uint k;
        for(i = 0; i < loopCount; i++) {
            assembly {
                k := callcode(gas, to, 0, 0, 0, 0, 0)
            }
        }
        return i;
    }

    function StaticCall(uint loopCount, address to, bytes input) returns (uint) {
        uint i;
        uint k;
        uint len = input.length;
        for(i = 0; i < loopCount; i++) {
            assembly {
                k := staticcall(gas, to, add(input,0x20), len, 0, 0)
            }
        }
        return k;
    }

    function DelegateCall(uint loopCount, address to) returns (uint) {
        uint i;
        uint k;
        for(i = 0; i < loopCount; i++) {
            assembly {
                k := delegatecall(gas, to, 0, 0, 0, 0)
            }
        }
        return i;
    }

    function Create(uint loopCount) returns (address) {
        uint i;
        address addr;
        for(i = 0; i < loopCount; i++) {
            StopContract c = new StopContract();
            addr = address(c);
        }
        return addr;
    }

    function Create2(uint loopCount, bytes code, uint salt) returns (address) {
        uint i;
        address addr;
        uint codelen = code.length;
        for(i = 0; i < loopCount; i++) {
            assembly {
                addr := create2(0, add(code,0x20), codelen, salt)
            }
        }
        return addr;
    }

    function Sstore(uint loopCount, uint start) returns (uint) {
        uint i;
        uint k;
        for (i = 0; i < loopCount; i++) {
            testMap[start + i] = i;
        }
        return k;
    }

    function Sload(uint loopCount) returns (uint) {
        uint i;
        uint k;
        for (i = 0; i < loopCount; i++) {
            k = testMap[i];
        }
        return k;
    }

    function Mstore(uint loopCount, uint memsize) returns (uint) {
        uint i;
        uint k;
        uint[] memory arr = new uint[](memsize);
        for (i = 0; i < loopCount; i++) {
            arr[i%memsize] = i;
        }
        return arr[memsize-1];
    }

    function Mload(uint loopCount, uint memsize) returns (uint) {
        uint i;
        uint k;
        uint[] memory arr = new uint[](memsize);
        for (i = 0; i < loopCount; i++) {
            k = arr[i%memsize];
        }
        return k;
    }

    function Msize(uint loopCount, uint memsize) returns (uint) {
        uint i;
        uint k;
        uint[] memory arr = new uint[](memsize);
        for (i = 0; i < loopCount; i++) {
            assembly {
                k := msize()
            }
        }
        return k;
    }

    function Gas(uint loopCount) returns (uint) {
        uint i;
        uint k;
        for (i = 0; i < loopCount; i++) {
            assembly {
                k := gas()
            }
        }
        return k;
    }

    function Address(uint loopCount) returns (uint) {
        uint i;
        uint k;
        for (i = 0; i < loopCount; i++) {
            assembly {
                k := address()
            }
        }
        return k;
    }

    function Balance(uint loopCount, address addr) returns (uint) {
        uint i;
        uint k;
        for (i = 0; i < loopCount; i++) {
            assembly {
                k := balance(addr)
            }
        }
        return k;
    }

    function Caller(uint loopCount) returns (uint) {
        uint i;
        uint k;
        for (i = 0; i < loopCount; i++) {
            assembly {
                k := caller()
            }
        }
        return k;
    }

    function CallValue(uint loopCount) returns (uint) {
        uint i;
        uint k;
        for (i = 0; i < loopCount; i++) {
            assembly {
                k := callvalue()
            }
        }
        return k;
    }

    function CallDataLoad(uint loopCount) returns (uint) {
        uint i;
        uint k;
        for (i = 0; i < loopCount; i++) {
            assembly {
                k := calldataload(i)
            }
        }
        return k;
    }

    function CallDataSize(uint loopCount) returns (uint) {
        uint i;
        uint k;
        for (i = 0; i < loopCount; i++) {
            assembly {
                k := calldatasize()
            }
        }
        return k;
    }

    function CallDataCopy(uint loopCount, bytes data, uint size) returns (uint) {
        uint i;
        uint k;
        uint[] memory arr = new uint[](size);
        for (i = 0; i < loopCount; i++) {
            assembly {
                calldatacopy(arr, 0, size)
            }
        }
        return k;
    }

    function CodeSize(uint loopCount) returns (uint) {
        uint i;
        uint k;
        for (i = 0; i < loopCount; i++) {
            assembly {
                k := codesize()
            }
        }
        return k;
    }

    function CodeCopy(uint loopCount, uint size) returns (uint) {
        uint i;
        uint k;
        uint[] memory arr = new uint[](size);
        for (i = 0; i < loopCount; i++) {
            assembly {
                codecopy(arr, 0, size)
            }
        }
        return k;
    }

    function ExtCodeSize(uint loopCount, address addr) returns (uint) {
        uint i;
        uint k;
        for (i = 0; i < loopCount; i++) {
            assembly {
                k := extcodesize(addr)
            }
        }
        return k;
    }

    function ExtCodeCopy(uint loopCount, address addr, uint size) returns (uint) {
        uint i;
        uint k;
        uint[] memory arr = new uint[](size);
        for (i = 0; i < loopCount; i++) {
            assembly {
                extcodecopy(addr, arr, 0, size)
            }
        }
        return k;
    }

    function ReturnDataSize(uint loopCount) returns (uint) {
        uint i;
        uint k;
        for (i = 0; i < loopCount; i++) {
            assembly {
                k := returndatasize()
            }
        }
        return k;
    }

    function ReturnDataCopy(uint loopCount, uint size, address addr) returns (uint) {
        uint i;
        uint k;
        bytes memory arr2 = new bytes(size);
        StopContract(addr).ReturnData(size);
        assembly {
            k := returndatasize
        }
        for (i = 0; i < loopCount; i++) {
            assembly {
                returndatacopy(arr2, 0, size)
            }
        }
        return k;
    }

    function Log0(uint loopCount, bytes l) returns (uint) {
        uint i;
        uint k;
        uint len = l.length;
        for (i = 0; i < loopCount; i++) {
            assembly {
                log0(add(l,32), l)
            }
        }
        return k;
    }

    function Log1(uint loopCount, bytes l, uint t1) returns (uint) {
        uint i;
        uint k;
        uint len = l.length;
        for (i = 0; i < loopCount; i++) {
            assembly {
                log1(add(l,32), l, t1)
            }
        }
        return k;
    }

    function Log2(uint loopCount, bytes l, uint t1, uint t2) returns (uint) {
        uint i;
        uint k;
        uint len = l.length;
        for (i = 0; i < loopCount; i++) {
            assembly {
                log2(add(l,32), l, t1, t2)
            }
        }
        return k;
    }

    function Log3(uint loopCount, bytes l, uint t1, uint t2, uint t3) returns (uint) {
        uint i;
        uint k;
        uint len = l.length;
        for (i = 0; i < loopCount; i++) {
            assembly {
                log3(add(l,32), l, t1, t2, t3)
            }
        }
        return k;
    }

    function Log4(uint loopCount, bytes l, uint t1, uint t2, uint t3, uint t4) returns (uint) {
        uint i;
        uint k;
        uint len = l.length;
        for (i = 0; i < loopCount; i++) {
            assembly {
                log4(add(l,32), l, t1, t2, t3, t4)
            }
        }
        return k;
    }

    function Origin(uint loopCount) returns (uint) {
        uint i;
        uint k;
        for (i = 0; i < loopCount; i++) {
            assembly {
                k := origin()
            }
        }
        return k;
    }

    function GasPrice(uint loopCount) returns (uint) {
        uint i;
        uint k;
        for (i = 0; i < loopCount; i++) {
            assembly {
                k := gasprice()
            }
        }
        return k;
    }

    function BlockHash(uint loopCount, uint blkNum) returns (uint) {
        uint i;
        uint k;
        for (i = 0; i < loopCount; i++) {
            assembly {
                k := blockhash(blkNum)
            }
        }
        return k;
    }

    function Coinbase(uint loopCount) returns (uint) {
        uint i;
        uint k;
        for (i = 0; i < loopCount; i++) {
            assembly {
                k := coinbase()
            }
        }
        return k;
    }

    function Timestamp(uint loopCount) returns (uint) {
        uint i;
        uint k;
        for (i = 0; i < loopCount; i++) {
            assembly {
                k := timestamp()
            }
        }
        return k;
    }

    function Number(uint loopCount) returns (uint) {
        uint i;
        uint k;
        for (i = 0; i < loopCount; i++) {
            assembly {
                k := number()
            }
        }
        return k;
    }

    function Difficulty(uint loopCount) returns (uint) {
        uint i;
        uint k;
        for (i = 0; i < loopCount; i++) {
            assembly {
                k := difficulty()
            }
        }
        return k;
    }

    function GasLimit(uint loopCount) returns (uint) {
        uint i;
        uint k;
        for (i = 0; i < loopCount; i++) {
            assembly {
                k := gaslimit()
            }
        }
        return k;
    }

    function Combination(uint loopCount, uint memsize) returns (uint) {
        uint i;
        uint k;
        uint[] memory arr = new uint[](memsize);
        for (i = 0; i < loopCount; i++) {
            k = k + i;
            arr[i%memsize] = i;
            testMap[i] = i;
            k = testMap[i];
            k = k * i;
            k = arr[i%memsize];
            k = k * i;
        }
        return k;
    }

    function precompiledContractTest(uint loopCount, address addr, bytes input) returns (uint) {
        uint len = input.length;
        uint i;
        uint k;

        for (i = 0; i < 100000; i++) {
            assembly {
                call(gas, addr, 0, add(input,0x20), len, 0, 0)
                pop
            }
        }

        return k;
    }

}
