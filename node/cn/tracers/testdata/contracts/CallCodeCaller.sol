pragma solidity ^0.4.24;

contract Caller {
    address callee;
    constructor(address _callee) public {
        callee = _callee;
    }

    function callHelloWorld() public {
        callee.callcode(abi.encodeWithSignature("helloWorld()"));
    }
}

contract Callee {
    function helloWorld() public {}
}
