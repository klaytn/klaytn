// SPDX-License-Identifier: MIT
// Klaytn Contract Library v1.0.0 (KIP/mocks/KIP7SnapshotMock.sol)
// Based on OpenZeppelin Contracts v4.5.0 (mocks/ERC20SnapshotMock.sol)
// https://github.com/OpenZeppelin/openzeppelin-contracts/releases/tag/v4.5.0

pragma solidity ^0.8.0;

import "../token/KIP7/extensions/KIP7Snapshot.sol";

contract KIP7SnapshotMock is KIP7Snapshot {
    constructor(
        string memory name,
        string memory symbol,
        address initialAccount,
        uint256 initialBalance
    ) KIP7(name, symbol) {
        _mint(initialAccount, initialBalance);
    }

    function snapshot() public {
        _snapshot();
    }

    function mint(address account, uint256 amount) public {
        _mint(account, amount);
    }

    function burn(address account, uint256 amount) public {
        _burn(account, amount);
    }
}
