// SPDX-License-Identifier: MIT
// Klaytn Contract Library v1.0.0 (KIP/mocks/KIP13CheckerMock.sol)
// Based on OpenZeppelin Contracts v4.5.0 (mocks/ERC165CheckerMock.sol)
// https://github.com/OpenZeppelin/openzeppelin-contracts/releases/tag/v4.5.0

pragma solidity ^0.8.0;

import "../utils/introspection/KIP13Checker.sol";

contract KIP13CheckerMock {
    using KIP13Checker for address;

    function supportsKIP13(address account) public view returns (bool) {
        return account.supportsKIP13();
    }

    function supportsInterface(address account, bytes4 interfaceId) public view returns (bool) {
        return account.supportsInterface(interfaceId);
    }

    function supportsAllInterfaces(address account, bytes4[] memory interfaceIds) public view returns (bool) {
        return account.supportsAllInterfaces(interfaceIds);
    }

    function getSupportedInterfaces(address account, bytes4[] memory interfaceIds) public view returns (bool[] memory) {
        return account.getSupportedInterfaces(interfaceIds);
    }
}
