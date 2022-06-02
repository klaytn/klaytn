// SPDX-License-Identifier: MIT
// Klaytn Contract Library v1.0.0 (KIP/token/KIP17/extensions/IKIP17Burnable.sol)

pragma solidity ^0.8.0;

/**
 * @title KIP17 Non-Fungible Token Standard, optional burnable extension
 * @dev See https://kips.klaytn.com/KIPs/kip-17#burning-extension
 */
interface IKIP17Burnable {
    /**
     * @dev Destroys `tokenId` token
     *
     * Emits a {Transfer} event with the 0x0 address as `to`.
     */
    function burn(uint256 tokenId) external;
}
