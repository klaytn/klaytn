pragma solidity ^0.5.6;

contract Factory {
    function createContract (bytes32 name) public {
        new Contract(name);
    }
}

contract Contract {
    constructor(bytes32 name) public {}
}