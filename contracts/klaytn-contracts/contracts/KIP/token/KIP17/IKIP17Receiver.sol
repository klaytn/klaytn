// SPDX-License-Identifier: MIT
// Klaytn Contract Library v1.0.0 (KIP/token/KIP17/IKIP17Receiver.sol)
// Based on OpenZeppelin Contracts v4.5.0 (token/ERC721/IERC721Receiver.sol)
// https://github.com/OpenZeppelin/openzeppelin-contracts/releases/tag/v4.5.0

pragma solidity ^0.8.0;

/**
 * @title KIP17 token receiver interface
 * @dev Interface for any contract that wants to support safeTransfers
 * from KIP17 asset contracts.
 */
interface IKIP17Receiver {
    /**
     * @dev Whenever an {IKIP17} `tokenId` token is transferred to this contract via {IKIP17-safeTransferFrom}
     * by `operator` from `from`, this function is called.
     *
     * It must return its Solidity selector to confirm the token transfer.
     * If any other value is returned or the interface is not implemented by the recipient, the transfer will be reverted.
     *
     * The selector can be obtained in Solidity with `IKIP17Receiver.onKIp17Received.selector`.
     */
    function onKIP17Received(
        address operator,
        address from,
        uint256 tokenId,
        bytes calldata data
    ) external returns (bytes4);
}
