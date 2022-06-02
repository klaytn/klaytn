// SPDX-License-Identifier: MIT
// Klaytn Contract Library v1.0.0 (KIP/token/KIP17/utils/KIP17ERC721Holder.sol)
// Based on OpenZeppelin Contracts v4.5.0 (token/ERC721/utils/ERC721Holder.sol)
// https://github.com/OpenZeppelin/openzeppelin-contracts/releases/tag/v4.5.0

pragma solidity ^0.8.0;

import "../IKIP17Receiver.sol";
import "../../../../token/ERC721/IERC721Receiver.sol";

/**
 * @dev Implementation of a smart contract which implements both the {IKIP17Receiver} and {IERC721Receiver} interface.
 *
 * Accepts all KIP17 and ERC721 token transfers.
 * Make sure the contract is able to use accepted tokens with {IKIP17-safeTransferFrom}, {IKIP17-approve} or {IKIP17-setApprovalForAll}.
 * OR
 * {IERC721-safeTransferFrom}, {IERC721-approve} or {IER721-setApprovalForAll}.
 * OR both interfaces
 */
contract KIP17ERC721Holder is IKIP17Receiver, IERC721Receiver {
    /**
     * @dev See {IKIP17Receiver-onERC721Received}.
     *
     * Always returns `IIKIP17Receiver.onKIP17Received.selector`.
     */
    function onKIP17Received(
        address,
        address,
        uint256,
        bytes memory
    ) public virtual override returns (bytes4) {
        return this.onKIP17Received.selector;
    }

    /**
     * @dev See {IERC721Receiver-onERC721Received}.
     *
     * Always returns `IERC721Receiver.onERC721Received.selector`.
     */
    function onERC721Received(
        address,
        address,
        uint256,
        bytes memory
    ) public virtual override returns (bytes4) {
        return this.onERC721Received.selector;
    }
}
