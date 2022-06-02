// SPDX-License-Identifier: MIT
// Klaytn Contract Library v1.0.0 (KIP/mocks/KIP7VotesCompMock.sol)
// Based on OpenZeppelin Contracts v4.5.0 (mocks/ERC20VotesCompMock.sol)
// https://github.com/OpenZeppelin/openzeppelin-contracts/releases/tag/v4.5.0

pragma solidity ^0.8.0;

import "../token/KIP7/extensions/KIP7VotesComp.sol";

contract KIP7VotesCompMock is KIP7VotesComp {
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
