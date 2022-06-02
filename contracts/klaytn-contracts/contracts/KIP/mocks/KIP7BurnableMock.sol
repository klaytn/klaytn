// SPDX-License-Identifier: MIT
// Klaytn Contract Library v1.0.0 (KIP/mocks/KIP7BurnableMock.sol)
// Based on OpenZeppelin Contracts v4.5.0 (mocks/ERC20BurnableMock.sol)
// https://github.com/OpenZeppelin/openzeppelin-contracts/releases/tag/v4.5.0

pragma solidity ^0.8.0;

import "../token/KIP7/extensions/KIP7Burnable.sol";

contract KIP7BurnableMock is KIP7Burnable {
    constructor(
        string memory name,
        string memory symbol,
        address initialAccount,
        uint256 initialBalance
    ) KIP7(name, symbol) {
        _mint(initialAccount, initialBalance);
    }
}
