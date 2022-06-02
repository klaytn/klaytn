// SPDX-License-Identifier: MIT
// Klaytn Contract Library v1.0.0 (KIP/token/KIP37/extensions/IKIP37Mintable.sol)

pragma solidity ^0.8.0;

/**
 * @dev Minting extension of the KIP37 standard as defined in the KIP.
 * See http://kips.klaytn.com/KIPs/kip-37#minting-extension
 */
interface IKIP37Mintable {
    /**
     * @dev Creates a new `id` token type and assigns the caller as owner of `initialSupply` while
     * setting a `uri` for this token type
     *
     * Requirements:
     *
     * - `id` must not already exist
     *
     * Emits a {TransferSingle} event with 0X0 as the `from` account, for the `intialSupply` tokens
     */
    function create(
        uint256 id,
        uint256 initialSupply,
        string calldata uri_
    ) external returns (bool);

    /**
     * @dev Mints an `amount` of new `id` tokens and assigns `to` as owner
     *
     * Emits a {TransferSingle} event with 0X0 as the `from` account, for the `amount` tokens
     *
     * Requirements:
     *
     * - `id` must exist
     * - `to` must not be the zero address
     */
    function mint(
        uint256 id,
        address to,
        uint256 amount
    ) external;

    /**
     * @dev For each item in `toList`, mints an `amount[i]` of new `id` tokens and assigns `toList[i]` as owner
     *
     * Emits multiple {TransferSingle} events with 0X0 as the `from` account, for the `amount` tokens
     *
     * Requirements:
     *
     * - `id` must exist
     * - each `toList[i]` must not be the zero address
     * - `toList` and `amounts` must have the same number of elements
     */
    function mint(
        uint256 id,
        address[] calldata toList,
        uint256[] calldata amounts
    ) external;

    /**
     * @dev Mints multiple KIP37 token types `ids` in a batch and assigns the tokens according to the variables `to` and `amounts`.
     *
     *
     * Emits a {TransferBatch} event with 0X0 as the `from` account, for the `amount` tokens
     *
     * Requirements:
     *
     * - `to` must not be the zero address
     * - each`ids[i]` must exist
     * - `ids` and `amounts` must have the same number of elements
     */
    function mintBatch(
        address to,
        uint256[] calldata ids,
        uint256[] calldata amounts
    ) external;
}
