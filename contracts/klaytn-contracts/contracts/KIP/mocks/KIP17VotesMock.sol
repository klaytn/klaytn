// SPDX-License-Identifier: MIT
// Klaytn Contract Library v1.0.0 (KIP/mocks/KIP17VotesMock.sol)
// Based on OpenZeppelin Contracts v4.5.0 (mocks/ERC721VotesMock.sol)
// https://github.com/OpenZeppelin/openzeppelin-contracts/releases/tag/v4.5.0

pragma solidity ^0.8.0;

import "../token/KIP17/extensions/draft-KIP17Votes.sol";

contract KIP17VotesMock is KIP17Votes {
    constructor(string memory name, string memory symbol) KIP17(name, symbol) EIP712(name, "1") {}

    function getTotalSupply() public view returns (uint256) {
        return _getTotalSupply();
    }

    function mint(address account, uint256 tokenId) public {
        _mint(account, tokenId);
    }

    function burn(uint256 tokenId) public {
        _burn(tokenId);
    }

    function getChainId() external view returns (uint256) {
        return block.chainid;
    }
}
