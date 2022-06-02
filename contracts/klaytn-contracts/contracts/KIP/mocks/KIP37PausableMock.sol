// SPDX-License-Identifier: MIT
// Klaytn Contract Library v1.0.0 (KIP/mocks/KIP37PausableMock.sol)
// Based on OpenZeppelin Contracts v4.5.0 (mocks/ERC1155PausableMock.sol)
// https://github.com/OpenZeppelin/openzeppelin-contracts/releases/tag/v4.5.0

pragma solidity ^0.8.0;

import "./KIP37Mock.sol";
import "../token/KIP37/extensions/KIP37Pausable.sol";

contract KIP37PausableMock is KIP37Mock, KIP37Pausable {
    constructor(string memory uri) KIP37Mock(uri) {
        _setupRole(DEFAULT_ADMIN_ROLE, _msgSender());
        _setupRole(PAUSER_ROLE, _msgSender());
    }

    function supportsInterface(bytes4 interfaceId) public view virtual override(KIP37, KIP37Pausable) returns (bool) {
        return KIP37.supportsInterface(interfaceId) || KIP37Pausable.supportsInterface(interfaceId);
    }

    function _beforeTokenTransfer(
        address operator,
        address from,
        address to,
        uint256[] memory ids,
        uint256[] memory amounts,
        bytes memory data
    ) internal virtual override(KIP37, KIP37Pausable) {
        super._beforeTokenTransfer(operator, from, to, ids, amounts, data);
    }
}
