// SPDX-License-Identifier: MIT
// Klaytn Contract Library v1.0.0 (KIP/mocks/KIP17RoyaltyMock.sol)
// Based on OpenZeppelin Contracts v4.5.0 (mocks/ERC721RoyaltyMock.sol)
// https://github.com/OpenZeppelin/openzeppelin-contracts/releases/tag/v4.5.0

pragma solidity ^0.8.0;

import "../token/KIP17/extensions/KIP17Royalty.sol";

contract KIP17RoyaltyMock is KIP17Royalty {
    constructor(string memory name, string memory symbol) KIP17(name, symbol) {}

    function setTokenRoyalty(
        uint256 tokenId,
        address recipient,
        uint96 fraction
    ) public {
        _setTokenRoyalty(tokenId, recipient, fraction);
    }

    function setDefaultRoyalty(address recipient, uint96 fraction) public {
        _setDefaultRoyalty(recipient, fraction);
    }

    function mint(address to, uint256 tokenId) public {
        _mint(to, tokenId);
    }

    function burn(uint256 tokenId) public {
        _burn(tokenId);
    }

    function deleteDefaultRoyalty() public {
        _deleteDefaultRoyalty();
    }
}
