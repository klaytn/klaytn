// SPDX-License-Identifier: MIT
// Klaytn Contract Library v1.0.0 (KIP/mocks/KIP7FlashMint.sol)
// Based on OpenZeppelin Contracts v4.5.0 (mocks/ERC20FlashMintMock.sol)
// https://github.com/OpenZeppelin/openzeppelin-contracts/releases/tag/v4.5.0

pragma solidity ^0.8.0;

import "../token/KIP7/extensions/KIP7FlashMint.sol";

contract KIP7FlashMintMock is KIP7FlashMint {
    constructor(
        string memory name,
        string memory symbol,
        address initialAccount,
        uint256 initialBalance
    ) KIP7(name, symbol) {
        _mint(initialAccount, initialBalance);
    }
}
