pragma solidity ^0.4.24;

contract IERC20BridgeReceiver {
    function onERC20Received(address _from, uint256 _amount, address _to, uint256 _feeLimit) public;
}
