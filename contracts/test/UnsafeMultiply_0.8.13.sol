pragma solidity ^0.8.13;

contract UnsafeMultiply {
    function multiply(uint a) public returns(uint d) {
        return a * 7;
    }
}
