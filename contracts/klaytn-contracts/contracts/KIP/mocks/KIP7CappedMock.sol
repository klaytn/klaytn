// SPDX-License-Identifier: MIT
// Klaytn Contract Library v1.0.0 (KIP/mocks/KIP7CappedMock.sol)
// Based on OpenZeppelin Contracts v4.5.0 (mocks/ERC20CappedMock.sol)
// https://github.com/OpenZeppelin/openzeppelin-contracts/releases/tag/v4.5.0

pragma solidity ^0.8.0;

import "../token/KIP7/extensions/KIP7Capped.sol";

contract KIP7CappedMock is KIP7Capped {
    constructor(
        string memory name,
        string memory symbol,
        uint256 cap
    ) KIP7(name, symbol) KIP7Capped(cap) {}

    function mint(address to, uint256 tokenId) public {
        _mint(to, tokenId);
    }
}
