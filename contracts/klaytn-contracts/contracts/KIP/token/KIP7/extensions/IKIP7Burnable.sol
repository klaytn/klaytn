// SPDX-License-Identifier: MIT
// Klaytn Contract Library v1.0.0 (KIP/token/KIP7/extensions/IKIP7Burnable.sol)

pragma solidity ^0.8.0;

/**
 * @dev Buring extension of the KIP7 standard as defined in the KIP.
 * See https://kips.klaytn.com/KIPs/kip-7#burning-extension
 */
interface IKIP7Burnable {
    /**
     * @dev Destroys `amount` tokens from the caller's account
     *
     * Emits a {Transfer} event with the 0x0 address as `to`.
     */
    function burn(uint256 amount) external;

    /**
     * @dev Destroys `amount` tokens from `account` using
     * the {IKIP7-allowance} mechanism.
     *
     * Emits a {Transfer} event with the 0x0 address as `to`.
     */
    function burnFrom(address account, uint256 amount) external;
}
