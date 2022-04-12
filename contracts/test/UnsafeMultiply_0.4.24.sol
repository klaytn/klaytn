pragma solidity ^0.4.24;

contract UnsafeMultiply {
    function multiply(uint a) public returns(uint d) {
        return a * 7;
    }
}
