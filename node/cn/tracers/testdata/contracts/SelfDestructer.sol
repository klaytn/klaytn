pragma solidity ^0.5.6;

contract SelfDestructer {
    constructor() public {}

    function selfDestruct(address payable a) public {
        selfdestruct(a);
    }
}
