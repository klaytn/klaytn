pragma solidity ^0.4.24;

contract FeePayer {
    function GetFeePayerDirect() public returns (address) {
        assembly {
            if iszero(call(gas, 0x0a, 0, 0, 0, 12, 20)) {
              invalid()
            }
            return(0, 32)
        }
    }

    function feePayer() internal returns (address addr) {
        assembly {
            let freemem := mload(0x40)
            let start_addr := add(freemem, 12)
            if iszero(call(gas, 0x0a, 0, 0, 0, start_addr, 20)) {
              invalid()
            }
            addr := mload(freemem)
        }
    }

    function GetFeePayer() public returns (address) {
        return feePayer();
    }
}
