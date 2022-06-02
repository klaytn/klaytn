// SPDX-License-Identifier: MIT
// Klaytn Contract Library v1.0.0 (KIP/mocks/KIP17MetadataMintableMock.sol)

pragma solidity ^0.8.0;

import "../token/KIP17/extensions/KIP17MetadataMintable.sol";

contract KIP17MetadataMintableMock is KIP17MetadataMintable {
    string private _baseTokenURI;

    constructor(string memory name, string memory symbol) KIP17(name, symbol) {
        _setupRole(DEFAULT_ADMIN_ROLE, _msgSender());
        _setupRole(MINTER_ROLE, _msgSender());
    }

    function _baseURI() internal view virtual override returns (string memory) {
        return _baseTokenURI;
    }

    function setBaseURI(string calldata newBaseTokenURI) public {
        _baseTokenURI = newBaseTokenURI;
    }

    function baseURI() public view returns (string memory) {
        return _baseURI();
    }

    function setTokenURI(uint256 tokenId, string memory _tokenURI) public {
        _setTokenURI(tokenId, _tokenURI);
    }

    function exists(uint256 tokenId) public view returns (bool) {
        return _exists(tokenId);
    }

    function mintWithTokenURI(
        address to,
        uint256 tokenId,
        string memory _tokenURI
    ) public virtual override returns (bool) {
        return super.mintWithTokenURI(to, tokenId, _tokenURI);
    }

    function burn(uint256 tokenId) public {
        _burn(tokenId);
    }
}
