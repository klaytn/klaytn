// SPDX-License-Identifier: MIT

pragma solidity >=0.7.0 <0.9.0;

// References
// - https://www.evm.codes
// - https://ethervm.io/decompile
// - https://docs.soliditylang.org/en/v0.4.26/assembly.html

contract CornerCaseCall {
    event Done(); // to make functions non-view.
    function doit() public {
        assembly {
            let g := gas()
            let c := caller()
            let s := call(g, c, 0, 0xfffffffffffffff0, 0, 0xffffffffffffffe0, 0)
        }
        emit Done();
    }
}

contract CornerCaseCreate {
    event Done();
    function doit() public {
        assembly {
            let a := create(0, 0xffffffffffffffd0, 0)
        }
        emit Done();
    }
}

contract CornerCaseRevertSmall {
    event Done();
    function doit() public {
        assembly {
            mstore(0x100, 0x0000000000000000000000000000000000000000000000000000000008c379a0) // sha3("Error(string)")
            mstore(0x120, 0x000000000000000000000000000000000000000000000000000000000000cafe) // string offset
            mstore(0x140, 0x0000000000000000000000000000000000000000000000000000000000c0ffee) // string length
            mstore(0x160, 0x4141414141414141000000000000000000000000000000000000000000000000) // string bytes
            revert(0x11c, 0x64)
        }
        emit Done();
    }
}

contract CornerCaseRevertBig {
    event Done();
    function doit() public {
        assembly {
            mstore(0x100, 0x0000000000000000000000000000000000000000000000000000000008c379a0) // sha3("Error(string)")
            mstore(0x120, 0x1100000000000000000000000000000000000000000000000000000000000020) // string offset
            mstore(0x140, 0x1100000000000000000000000000000000000000000000000000000000000008) // string length
            mstore(0x160, 0x4141414141414141000000000000000000000000000000000000000000000000) // string bytes
            revert(0x11c, 0x64)
        }
        emit Done();
    }
}
