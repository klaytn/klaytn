// SPDX-License-Identifier: MIT
// Klaytn Contract Library v1.0.0 (KIP/mocks/KIP37MintableMock.sol)

pragma solidity ^0.8.0;

import "../token/KIP37/extensions/KIP37Mintable.sol";

contract KIP37MintableMock is KIP37Mintable {
    constructor(string memory uri_) KIP37(uri_) {
        _setupRole(DEFAULT_ADMIN_ROLE, _msgSender());
        _setupRole(MINTER_ROLE, _msgSender());
    }

    function exists(uint256 tokenId) public view returns (bool) {
        return _exists(tokenId);
    }

    function create(
        uint256 id,
        uint256 initialSupply,
        string calldata uri_
    ) public override returns (bool) {
        return super.create(id, initialSupply, uri_);
    }

    function mint(
        uint256 id,
        address to,
        uint256 amount
    ) public override {
        super.mint(id, to, amount);
    }

    function mint(
        uint256 id,
        address[] calldata toList,
        uint256[] calldata amounts
    ) public override {
        super.mint(id, toList, amounts);
    }

    function mintBatch(
        address to,
        uint256[] calldata ids,
        uint256[] calldata amounts
    ) public override {
        super.mintBatch(to, ids, amounts);
    }

    function burn(
        address from,
        uint256 id,
        uint256 amount
    ) public {
        KIP37._burn(from, id, amount);
    }
}
