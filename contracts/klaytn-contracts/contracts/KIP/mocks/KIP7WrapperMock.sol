// SPDX-License-Identifier: MIT
// Klaytn Contract Library v1.0.0 (KIP/mocks/KIP7WrapperMock.sol)
// Based on OpenZeppelin Contracts v4.5.0 (mocks/ERC20WrapperMock.sol)
// https://github.com/OpenZeppelin/openzeppelin-contracts/releases/tag/v4.5.0

pragma solidity ^0.8.0;

import "../token/KIP7/extensions/KIP7Wrapper.sol";

contract KIP7WrapperMock is KIP7Wrapper {
    constructor(
        IKIP7 _underlyingToken,
        string memory name,
        string memory symbol
    ) KIP7(name, symbol) KIP7Wrapper(_underlyingToken) {}

    function recover(address account) public returns (uint256) {
        return _recover(account);
    }
}
