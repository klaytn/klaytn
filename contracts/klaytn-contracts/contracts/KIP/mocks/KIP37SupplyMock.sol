// SPDX-License-Identifier: MIT
// Klaytn Contract Library v1.0.0 (KIP/mocks/KIP37SupplyMock.sol)
// Based on OpenZeppelin Contracts v4.5.0 (mocks/ERC1155SupplyMock.sol)
// https://github.com/OpenZeppelin/openzeppelin-contracts/releases/tag/v4.5.0

pragma solidity ^0.8.0;

import "./KIP37Mock.sol";
import "../token/KIP37/extensions/KIP37Supply.sol";

contract KIP37SupplyMock is KIP37Mock, KIP37Supply {
    constructor(string memory uri) KIP37Mock(uri) {}

    function _beforeTokenTransfer(
        address operator,
        address from,
        address to,
        uint256[] memory ids,
        uint256[] memory amounts,
        bytes memory data
    ) internal virtual override(KIP37, KIP37Supply) {
        super._beforeTokenTransfer(operator, from, to, ids, amounts, data);
    }
}
