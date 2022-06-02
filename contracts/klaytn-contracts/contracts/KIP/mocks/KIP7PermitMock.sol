// SPDX-License-Identifier: MIT
// Klaytn Contract Library v1.0.0 (KIP/mocks/KIP7PermitMock.sol)
// Based on OpenZeppelin Contracts v4.5.0 (mocks/ERC20PermitMock.sol)
// https://github.com/OpenZeppelin/openzeppelin-contracts/releases/tag/v4.5.0

pragma solidity ^0.8.0;

import "../token/KIP7/extensions/draft-KIP7Permit.sol";

contract KIP7PermitMock is KIP7Permit {
    constructor(
        string memory name,
        string memory symbol,
        address initialAccount,
        uint256 initialBalance
    ) payable KIP7(name, symbol) KIP7Permit(name) {
        _mint(initialAccount, initialBalance);
    }

    function getChainId() external view returns (uint256) {
        return block.chainid;
    }
}
