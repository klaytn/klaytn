// SPDX-License-Identifier: MIT
// Klaytn Contract Library v1.0.0 (KIP/utils/introspection/KIP13.sol)
// Based on OpenZeppelin Contracts v4.5.0 (utils/introspection/ERC165.sol)
// https://github.com/OpenZeppelin/openzeppelin-contracts/releases/tag/v4.5.0

pragma solidity ^0.8.0;

import "../../interfaces/IKIP13.sol";

/**
 * @dev Implementation of the {IKIP13} interface.
 *
 * Contracts that want to implement KIP13 should inherit from this contract and override {supportsInterface} to check
 * for the additional interface id that will be supported. For example:
 *
 * ```solidity
 * function supportsInterface(bytes4 interfaceId) public view virtual override returns (bool) {
 *     return interfaceId == type(MyInterface).interfaceId || super.supportsInterface(interfaceId);
 * }
 * ```
 *
 * Alternatively, {KIP13Storage} provides an easier to use but more expensive implementation.
 */
abstract contract KIP13 is IKIP13 {
    /**
     * @dev See {IKIP13-supportsInterface}.
     */
    function supportsInterface(bytes4 interfaceId) public view virtual override returns (bool) {
        return interfaceId == type(IKIP13).interfaceId;
    }
}
