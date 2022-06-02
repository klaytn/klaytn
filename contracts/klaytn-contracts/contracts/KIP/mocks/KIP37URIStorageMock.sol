// SPDX-License-Identifier: MIT
// Klaytn Contract Library v1.0.0 (KIP/mocks/KIP37URIStorageMock.sol)
// Based on OpenZeppelin Contracts v4.5.0 (mocks/ERC1155URIStorageMock.sol)
// https://github.com/OpenZeppelin/openzeppelin-contracts/releases/tag/v4.5.0

pragma solidity ^0.8.0;

import "./KIP37Mock.sol";
import "../token/KIP37/extensions/KIP37URIStorage.sol";

contract KIP37URIStorageMock is KIP37Mock, KIP37URIStorage {
    constructor(string memory _uri) KIP37Mock(_uri) {}

    function uri(uint256 tokenId) public view virtual override(KIP37, KIP37URIStorage) returns (string memory) {
        return KIP37URIStorage.uri(tokenId);
    }

    function setURI(uint256 tokenId, string memory _tokenURI) public {
        _setURI(tokenId, _tokenURI);
    }

    function setBaseURI(string memory baseURI) public {
        _setBaseURI(baseURI);
    }
}
