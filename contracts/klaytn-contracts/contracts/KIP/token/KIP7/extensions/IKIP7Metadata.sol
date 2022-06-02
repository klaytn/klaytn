// SPDX-License-Identifier: MIT
// Klaytn Contract Library v1.0.0 (KIP/token/KIP7/extensions/IKIP7Metadata.sol)
// Based on OpenZeppelin Contracts v4.5.0 (token/ERC20/extensions/IERC20Metadata.sol)
// https://github.com/OpenZeppelin/openzeppelin-contracts/releases/tag/v4.5.0

pragma solidity ^0.8.0;

import "../IKIP7.sol";

/**
 * @dev Extension of {KIP7} which exposes metadata functions.
 * See https://kips.klaytn.com/KIPs/kip-7#metadata-extension
 */
interface IKIP7Metadata is IKIP7 {
    /**
     * @dev Returns the name of the token.
     */
    function name() external view returns (string memory);

    /**
     * @dev Returns the symbol of the token.
     */
    function symbol() external view returns (string memory);

    /**
     * @dev Returns the decimals places of the token.
     */
    function decimals() external view returns (uint8);
}
