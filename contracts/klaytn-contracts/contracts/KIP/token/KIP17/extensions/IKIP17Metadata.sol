// SPDX-License-Identifier: MIT
// Klaytn Contract Library v1.0.0 (KIP/token/KIP17/extensions/IKIP17Metadata.sol)
// Based on OpenZeppelin Contracts v4.5.0 (token/ERC721/extensions/IERC721Metadata.sol)
// https://github.com/OpenZeppelin/openzeppelin-contracts/releases/tag/v4.5.0

pragma solidity ^0.8.0;

import "../IKIP17.sol";

/**
 * @title KIP-17 Non-Fungible Token Standard, optional metadata extension
 * @dev See https://kips.klaytn.com/KIPs/kip-17#metadata-extension
 */
interface IKIP17Metadata is IKIP17 {
    /**
     * @dev Returns the token collection name.
     */
    function name() external view returns (string memory);

    /**
     * @dev Returns the token collection symbol.
     */
    function symbol() external view returns (string memory);

    /**
     * @dev Returns the Uniform Resource Identifier (URI) for `tokenId` token.
     */
    function tokenURI(uint256 tokenId) external view returns (string memory);
}
