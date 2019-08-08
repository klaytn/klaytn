pragma solidity ^0.4.24;


contract IERC20BridgeReceiver {
    function onERC20Received(address _from, address _to, uint256 _amount, uint256 _feeLimit, uint256[] _extraData) public;
}
