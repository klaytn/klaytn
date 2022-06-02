// SPDX-License-Identifier: MIT
// Klaytn Contract Library v1.0.0 (KIP/token/KIP7/utils/KIP7Holder.sol)
// Based on OpenZeppelin Contracts v4.5.0 (token/ERC721/utils/ERC721Holder.sol)
// https://github.com/OpenZeppelin/openzeppelin-contracts/releases/tag/v4.5.0

pragma solidity ^0.8.0;

import "../IKIP7Receiver.sol";

/**
 * @dev Implementation of the {IKIP7Receiver} interface.
 *
 * Accepts all KIP7 token transfers.
 */
contract KIP7Holder is IKIP7Receiver {
    /**
     * @dev See {IKIP7Receiver-onKIP7Received}.
     *
     * Always returns `IKIP7Receiver.onKIP7Received.selector`.
     *
     * Customize this base implementation to parse and use _data for any specifc handler logic when needed
     */
    function onKIP7Received(
        address,
        address,
        uint256,
        bytes memory
    ) public virtual override returns (bytes4) {
        return this.onKIP7Received.selector;
    }
}
