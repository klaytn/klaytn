// SPDX-License-Identifier: MIT
// Klaytn Contract Library v1.0.0 (KIP/mocks/KIP37BurnableMock.sol)
// Based on OpenZeppelin Contracts v4.5.0 (mocks/ERC1155BurnableMock.sol)
// https://github.com/OpenZeppelin/openzeppelin-contracts/releases/tag/v4.5.0

pragma solidity ^0.8.0;

import "../token/KIP37/extensions/KIP37Burnable.sol";

contract KIP37BurnableMock is KIP37Burnable {
    constructor(string memory uri) KIP37(uri) {}

    function mint(
        address to,
        uint256 id,
        uint256 amount,
        bytes memory data
    ) public {
        _mint(to, id, amount, data);
    }
}
