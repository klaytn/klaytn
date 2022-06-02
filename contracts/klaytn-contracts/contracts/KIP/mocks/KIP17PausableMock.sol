// SPDX-License-Identifier: MIT
// Klaytn Contract Library v1.0.0 (KIP/mocks/KIP17PausableMock.sol)
// Based on OpenZeppelin Contracts v4.5.0 (mocks/ERC721PausableMock.sol)
// https://github.com/OpenZeppelin/openzeppelin-contracts/releases/tag/v4.5.0

pragma solidity ^0.8.0;

import "../token/KIP17/extensions/KIP17Pausable.sol";

contract KIP17PausableMock is KIP17Pausable {
    constructor(string memory name, string memory symbol) KIP17(name, symbol) {
        _setupRole(DEFAULT_ADMIN_ROLE, _msgSender());
        _setupRole(PAUSER_ROLE, _msgSender());
    }

    function exists(uint256 tokenId) public view returns (bool) {
        return _exists(tokenId);
    }

    function mint(address to, uint256 tokenId) public {
        _mint(to, tokenId);
    }

    function safeMint(address to, uint256 tokenId) public {
        _safeMint(to, tokenId);
    }

    function safeMint(
        address to,
        uint256 tokenId,
        bytes memory _data
    ) public {
        _safeMint(to, tokenId, _data);
    }

    function burn(uint256 tokenId) public {
        _burn(tokenId);
    }
}
