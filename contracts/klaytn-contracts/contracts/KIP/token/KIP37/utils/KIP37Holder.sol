// SPDX-License-Identifier: MIT
// Klaytn Contract Library v1.0.0 (KIP/token/KIP37/utils/KIP37Holder.sol)
// Based on OpenZeppelin Contracts v4.5.0 (token/ERC1155/utils/ERC1155Holder.sol)
// https://github.com/OpenZeppelin/openzeppelin-contracts/releases/tag/v4.5.0

pragma solidity ^0.8.0;

import "./KIP37Receiver.sol";

/**
 * Simple implementation of `KIP37Receiver` that will allow a contract to hold KIP37 tokens.
 *
 * IMPORTANT: When inheriting this contract, you must include a way to use the received tokens, otherwise they will be
 * stuck.
 */
contract KIP37Holder is KIP37Receiver {
    function onKIP37Received(
        address,
        address,
        uint256,
        uint256,
        bytes memory
    ) public virtual override returns (bytes4) {
        return this.onKIP37Received.selector;
    }

    function onKIP37BatchReceived(
        address,
        address,
        uint256[] memory,
        uint256[] memory,
        bytes memory
    ) public virtual override returns (bytes4) {
        return this.onKIP37BatchReceived.selector;
    }
}
