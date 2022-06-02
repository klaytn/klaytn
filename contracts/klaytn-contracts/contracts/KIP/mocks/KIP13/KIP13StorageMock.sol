// SPDX-License-Identifier: MIT

pragma solidity ^0.8.0;

import "../../utils/introspection/KIP13Storage.sol";

contract KIP13StorageMock is KIP13Storage {
    function registerInterface(bytes4 interfaceId) public {
        _registerInterface(interfaceId);
    }
}
