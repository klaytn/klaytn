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

import "../klaytn-contracts/contracts/token/ERC721/IERC721.sol";
import "../klaytn-contracts/contracts/token/ERC721/extensions/ERC721URIStorage.sol";
import "../klaytn-contracts/contracts/token/ERC721/extensions/ERC721Burnable.sol";

import "../bridge_interface/IERC721BridgeReceiver.sol";
import "./BridgeTransfer.sol";

abstract contract BridgeTransferERC721 is BridgeTokens, IERC721BridgeReceiver, BridgeTransfer {
    // handleERC721Transfer sends the ERC721 by the request.
    function handleERC721Transfer(
        bytes32 _requestTxHash,
        address _from,
        address _to,
        address _tokenAddress,
        uint256 _tokenId,
        uint64 _requestedNonce,
        uint64 _requestedBlockNumber,
        string memory _tokenURI,
        bytes memory _extraData
    )
        public
        virtual
        onlyOperators
    {
        _lowerHandleNonceCheck(_requestedNonce);

        if (!_voteValueTransfer(_requestedNonce)) {
            return;
        }

        _setHandledRequestTxHash(_requestTxHash);

        handleNoncesToBlockNums[_requestedNonce] = _requestedBlockNumber;
        _updateHandleNonce(_requestedNonce);

        emit HandleValueTransfer(
            _requestTxHash,
            TokenType.ERC721,
            _from,
            _to,
            _tokenAddress,
            _tokenId,
            _requestedNonce,
            lowerHandleNonce,
            _extraData
        );

        if (modeMintBurn) {
            IERC721Mint(_tokenAddress).mintWithTokenURI(_to, _tokenId, _tokenURI);
        } else {
            IERC721(_tokenAddress).safeTransferFrom(address(this), _to, _tokenId);
        }
    }

    // _requestERC721Transfer requests transfer ERC721 to _to on relative chain.
    function _requestERC721Transfer(
        address _tokenAddress,
        address _from,
        address _to,
        uint256 _tokenId,
        bytes memory _extraData
    )
        internal
        onlyRegisteredToken(_tokenAddress)
        onlyUnlockedToken(_tokenAddress)
    {
        require(isRunning, "stopped bridge");

        // returns tokenURI if the `tokenId` is found from private variable `_tokenURI`.
        // Otherwise, returns base URI, which is empty string by default.
        string memory uri = ERC721URIStorage(_tokenAddress).tokenURI(_tokenId);

        if (modeMintBurn) {
            ERC721Burnable(_tokenAddress).burn(_tokenId);
        }
        emit RequestValueTransferEncoded(
            TokenType.ERC721,
            _from,
            _to,
            _tokenAddress,
            _tokenId,
            requestNonce,
            0,
            _extraData,
            2,
            abi.encode(uri)
        );
        requestNonce++;
    }

    // onERC721Received function of ERC721 token for 1-step deposits to the Bridge
    function onERC721Received (
        address _from,
        uint256 _tokenId,
        address _to,
        bytes memory _extraData
    )
        public override
    {
        _requestERC721Transfer(msg.sender, _from, _to, _tokenId, _extraData);
    }

    // requestERC721Transfer requests transfer ERC721 to _to on relative chain.
    function requestERC721Transfer(
        address _tokenAddress,
        address _to,
        uint256 _tokenId,
        bytes memory _extraData
    )
        public
        virtual
    {
        IERC721(_tokenAddress).transferFrom(msg.sender, address(this), _tokenId);
        _requestERC721Transfer(_tokenAddress, msg.sender, _to, _tokenId, _extraData);
    }
}
