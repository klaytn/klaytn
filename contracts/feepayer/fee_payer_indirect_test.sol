pragma solidity ^0.4.24;

interface FeePayer {
    function GetFeePayer() public returns (address);
}

contract FeePayerIndirect {
    function TestCall(address _target) public returns (address) {
        return FeePayer(_target).GetFeePayer();
    }
}

