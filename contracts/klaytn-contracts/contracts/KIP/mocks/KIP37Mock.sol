// SPDX-License-Identifier: MIT
// Klaytn Contract Library v1.0.0 (KIP/mocks/KIP37Mock.sol)
// Based on OpenZeppelin Contracts v4.5.0 (mocks/ERC1155Mock.sol)
// https://github.com/OpenZeppelin/openzeppelin-contracts/releases/tag/v4.5.0

pragma solidity ^0.8.0;

import "../token/KIP37/KIP37.sol";

contract KIP37Mock is KIP37 {
    constructor(string memory uri) KIP37(uri) {}

    function setURI(string memory newuri) public {
        _setURI(newuri);
    }

    function mint(
        address to,
        uint256 id,
        uint256 amount,
        bytes memory data
    ) public {
        _mint(to, id, amount, data);
    }

    function mintBatch(
        address to,
        uint256[] memory ids,
        uint256[] memory amounts,
        bytes memory data
    ) public {
        _mintBatch(to, ids, amounts, data);
    }

    function burn(
        address owner,
        uint256 id,
        uint256 amount
    ) public {
        _burn(owner, id, amount);
    }

    function burnBatch(
        address owner,
        uint256[] memory ids,
        uint256[] memory amounts
    ) public {
        _burnBatch(owner, ids, amounts);
    }
}
