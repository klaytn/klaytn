// SPDX-License-Identifier: MIT
// Klaytn Contract Library v1.0.0 (KIP/token/KIP37/extensions/IKIP37MetadataURI.sol)
// Based on OpenZeppelin Contracts v4.5.0 (token/ERC1155/extensions/IERC1155MetadataURI.sol)
// https://github.com/OpenZeppelin/openzeppelin-contracts/releases/tag/v4.5.0

pragma solidity ^0.8.0;

import "../IKIP37.sol";

/**
 * @dev Interface of the optional KIP37Metadata Extension interface, as defined
 * in the http://kips.klaytn.com/KIPs/kip-37#metadata-extension[KIP].
 *
 */
interface IKIP37MetadataURI is IKIP37 {
    /**
     * @dev Returns the URI for token type `id`.
     *
     * If the `\{id\}` substring is present in the URI, it must be replaced by
     * clients with the actual token type ID.
     */
    function uri(uint256 id) external view returns (string memory);
}
