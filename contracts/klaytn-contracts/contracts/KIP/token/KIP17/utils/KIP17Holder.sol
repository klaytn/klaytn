// SPDX-License-Identifier: MIT
// Klaytn Contract Library v1.0.0 (KIP/token/KIP17/utils/KIP17Holder.sol)
// Based on OpenZeppelin Contracts v4.5.0 (token/ERC721/utils/ERC721Holder.sol)
// https://github.com/OpenZeppelin/openzeppelin-contracts/releases/tag/v4.5.0

pragma solidity ^0.8.0;

import "../IKIP17Receiver.sol";

/**
 * @dev Implementation of the {IKIP17Receiver} interface.
 *
 * Accepts all KIP17 token transfers.
 * Make sure the contract is able to use its token with {IKIP17-safeTransferFrom}, {IKIP17-approve} or {IKIP17-setApprovalForAll}.
 */
contract KIP17Holder is IKIP17Receiver {
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
}
