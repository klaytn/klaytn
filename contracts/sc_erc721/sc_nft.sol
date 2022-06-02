// SPDX-License-Identifier: MIT

// Copyright 2019 The klaytn Authors
// This file is part of the klaytn library.
//
// The klaytn library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The klaytn library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the klaytn library. If not, see <http://www.gnu.org/licenses/>.

pragma solidity ^0.8.0;

import "../klaytn-contracts/contracts/token/ERC721/ERC721.sol";
import "../klaytn-contracts/contracts/token/ERC721/extensions/ERC721URIStorage.sol";
import "../klaytn-contracts/contracts/token/ERC721/extensions/ERC721Burnable.sol";
import "../klaytn-contracts/contracts/token/ERC721/extensions/ERC721Enumerable.sol";
import "../klaytn-contracts/contracts/access/AccessControl.sol";
import "../klaytn-contracts/contracts/access/Ownable.sol";

import "./ERC721ServiceChain.sol";

contract ServiceChainNFT is ERC721, ERC721Enumerable,  ERC721URIStorage, ERC721Burnable, AccessControl, IERC721Mint, ERC721ServiceChain {
    bytes32 public constant MINTER_ROLE = keccak256("MINTER_ROLE");

    string public constant NAME = "ServiceChainNFT";
    string public constant SYMBOL = "SCN";

    constructor(address _bridge) ERC721(NAME, SYMBOL) ERC721ServiceChain(_bridge) {
        _grantRole(DEFAULT_ADMIN_ROLE, msg.sender);
        _grantRole(MINTER_ROLE, msg.sender);
    }

    function supportsInterface(bytes4 interfaceId)
        public
        view
        override(ERC721, ERC721Enumerable, AccessControl)
        returns (bool)
    {
        return super.supportsInterface(interfaceId);
    }

    function _burn(uint256 tokenId) internal override(ERC721, ERC721URIStorage) {
        super._burn(tokenId);
    }

	function _beforeTokenTransfer(address from, address to, uint256 tokenId)
        internal
        override(ERC721, ERC721Enumerable)
    {
        super._beforeTokenTransfer(from, to, tokenId);
    }

	function tokenURI(uint256 tokenId)
        public
        view
        override(ERC721, ERC721URIStorage)
        returns (string memory)
    {
        return super.tokenURI(tokenId);
    }

    function mintWithTokenURI(address to, uint256 tokenId, string memory _tokenURI)
        public
        onlyRole(MINTER_ROLE)
    {
        _mint(to, tokenId);
        _setTokenURI(tokenId, _tokenURI);
    }

    // registerBulk registers (startID, endID-1) tokens to the user once.
    // This is only for load test.
    function registerBulk(address _user, uint256 _startID, uint256 _endID) external onlyOwner {
        for (uint256 uid = _startID; uid < _endID; uid++) {
            mintWithTokenURI(_user, uid, "testURI");
        }
    }
}
