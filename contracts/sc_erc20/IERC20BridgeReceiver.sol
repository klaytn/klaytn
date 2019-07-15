pragma solidity ^0.4.24;

contract IERC20BridgeReceiver {
    function onERC20Received(address _from, uint256 amount, address _to) public;
}
