// SPDX-License-Identifier: MIT
// Klaytn Contract Library v1.0.0 (KIP/mocks/KIP7Mintable.sol)

pragma solidity ^0.8.0;

import "../token/KIP7/extensions/KIP7Mintable.sol";

contract KIP7MintableMock is KIP7Mintable {
    constructor(
        string memory name,
        string memory symbol,
        address initialAccount,
        uint256 initialBalance
    ) KIP7(name, symbol) {
        _mint(initialAccount, initialBalance);
        _setupRole(DEFAULT_ADMIN_ROLE, _msgSender());
        _setupRole(MINTER_ROLE, _msgSender());
    }

    function burn(address from, uint256 amount) public {
        _burn(from, amount);
    }
}
