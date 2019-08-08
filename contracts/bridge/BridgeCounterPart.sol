pragma solidity ^0.4.24;

import "../externals/openzeppelin-solidity/contracts/ownership/Ownable.sol";


contract BridgeCounterPart is Ownable {
    address public counterpartBridge;

    function setCounterPartBridge(address _bridge)
        external
        onlyOwner
    {
        counterpartBridge = _bridge;
    }
}
