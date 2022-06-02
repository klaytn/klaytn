// SPDX-License-Identifier: MIT
// Klaytn Contract Library v1.0.0 (KIP/token/KIP37/utils/KIP37Receiver.sol)
// Based on OpenZeppelin Contracts v4.5.0 (token/ERC1155/utils/ERC1155Receiver.sol)
// https://github.com/OpenZeppelin/openzeppelin-contracts/releases/tag/v4.5.0

pragma solidity ^0.8.0;

import "../../../utils/introspection/KIP13.sol";
import "../IKIP37Receiver.sol";

abstract contract KIP37Receiver is KIP13, IKIP37Receiver {
    /**
     * @dev See {IKIP37-supportsInterface}.
     */
    function supportsInterface(bytes4 interfaceId) public view virtual override(KIP13, IKIP13) returns (bool) {
        return interfaceId == type(IKIP37Receiver).interfaceId || super.supportsInterface(interfaceId);
    }
}
