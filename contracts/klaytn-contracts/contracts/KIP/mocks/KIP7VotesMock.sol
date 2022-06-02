// SPDX-License-Identifier: MIT
// Klaytn Contract Library v1.0.0 (KIP/mocks/KIP7VotesMock.sol)
// Based on OpenZeppelin Contracts v4.5.0 (mocks/ERC20VotesMock.sol)
// https://github.com/OpenZeppelin/openzeppelin-contracts/releases/tag/v4.5.0

pragma solidity ^0.8.0;

import "../token/KIP7/extensions/KIP7Votes.sol";

contract KIP7VotesMock is KIP7Votes {
    constructor(string memory name, string memory symbol) KIP7(name, symbol) KIP7Permit(name) {}

    function mint(address account, uint256 amount) public {
        _mint(account, amount);
    }

    function burn(address account, uint256 amount) public {
        _burn(account, amount);
    }

    function getChainId() external view returns (uint256) {
        return block.chainid;
    }
}
