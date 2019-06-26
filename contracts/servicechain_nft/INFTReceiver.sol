pragma solidity ^0.4.24;

/**
 * @title NFT receiver interface
 * @dev Interface for any contract that wants to support safeTransfers
 * from NFT asset contracts.
 */
contract INFTReceiver {
    /**
     * @notice Handle the receipt of an NFT
     * @dev The NFT smart contract calls this function on the recipient
     * after a `safeTransfer`. This function MUST return the function selector,
     * otherwise the caller will revert the transaction. The selector to be
     * returned can be obtained as `this.onERC721Received.selector`. This
     * function MAY throw to revert and reject the transfer.
     * Note: the NFT contract address is always the message sender.
     * @param from The address which previously owned the token
     * @param tokenId The NFT identifier which is being transferred
     * @param to The address which user(from) want to transfer to
     * @return `bytes4(keccak256("onERC721Received(address,address,uint256,bytes)"))`
     */
    // TODO-Klaytn-Servicechain define proper bytes4 value.
    function onNFTReceived(
        address from,
        uint256 tokenId,
        address to
    )
    public
    returns(bytes4);
}
