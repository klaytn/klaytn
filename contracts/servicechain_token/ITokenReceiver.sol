pragma solidity ^0.4.24;

/**
 * @title KLAY compatible token(Token) receiver interface
 * @dev Interface for any contract that wants to support safeTransfers
 *  from KLAY compatible token asset contracts.
 */
contract ITokenReceiver {
    /**
     * @dev Magic value to be returned upon successful reception of an NFT
     *  Equals to `bytes4(keccak256("onERC20Received(address,uint256,bytes)"))`,
     *  which can be also obtained as `ERC20Receiver(0).onERC20Received.selector`
     */
    // TODO-Klaytn-Servicechain define proper bytes4 value.
    function onTokenReceived(address _from, uint256 amount, address _to) public returns(bytes4);
}
