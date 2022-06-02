// SPDX-License-Identifier: MIT
// Klaytn Contract Library v1.0.0 (KIP/token/KIP7/presets/KIP7PresetFixedSupply.sol)
// Based on OpenZeppelin Contracts v4.5.0 (token/ERC20/presets/ERC20PresetFixedSupply.sol)
// https://github.com/OpenZeppelin/openzeppelin-contracts/releases/tag/v4.5.0

pragma solidity ^0.8.0;

import "../extensions/KIP7Burnable.sol";

/**
 * @dev {KIP7} token, including:
 *
 *  - Preminted initial supply
 *  - Ability for holders to burn (destroy) their tokens
 *  - No access control mechanism (for minting/pausing) and hence no governance
 *
 * This contract uses {KIP7Burnable} to include burn capabilities - head to
 * its documentation for details.
 */
contract KIP7PresetFixedSupply is KIP7Burnable {
    /**
     * @dev Mints `initialSupply` amount of token and transfers them to `owner`.
     *
     * See {KIP7-constructor}.
     */
    constructor(
        string memory name,
        string memory symbol,
        uint256 initialSupply,
        address owner
    ) KIP7(name, symbol) {
        _mint(owner, initialSupply);
    }
}
