// SPDX-License-Identifier: MIT
// Klaytn Contract Library v1.0.0 (KIP/token/KIP37/extensions/IKIP37Burnable.sol)

pragma solidity ^0.8.0;

/**
 * @title KIP37 Non-Fungible Token Standard, optional burnable extension
 * @dev See http://kips.klaytn.com/KIPs/kip-37#burning-extension
 */
interface IKIP37Burnable {
    /**
     * @dev Destroys `amount` tokens of token type `id` from `from`
     *
     * Emits a {TransferSingle} event with the 0x0 address as `to`.
     */
    function burn(
        address account,
        uint256 id,
        uint256 amount
    ) external;

    /**
     * @dev Destroys each amount `amounts[i]` tokens of each token type `ids[i]` from `from`
     *
     * Emits a {TransferBatch} event with the 0x0 address as `to`.
     */
    function burnBatch(
        address account,
        uint256[] calldata ids,
        uint256[] calldata amounts
    ) external;
}
