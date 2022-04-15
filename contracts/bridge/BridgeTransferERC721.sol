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

pragma solidity 0.5.6;

import "../externals/openzeppelin-solidity/contracts/token/ERC721/IERC721.sol";
import "../externals/openzeppelin-solidity/contracts/token/ERC721/ERC721MetadataMintable.sol";
import "../externals/openzeppelin-solidity/contracts/token/ERC721/ERC721Burnable.sol";

import "../sc_erc721/IERC721BridgeReceiver.sol";
import "./BridgeTransfer.sol";


contract BridgeTransferERC721 is BridgeTokens, IERC721BridgeReceiver, BridgeTransfer {
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
            require(ERC721MetadataMintable(_tokenAddress).mintWithTokenURI(_to, _tokenId, _tokenURI), "mint failed");
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
        bytes memory _extraData,
        uint ver
    )
        internal
        onlyRegisteredToken(_tokenAddress)
        onlyUnlockedToken(_tokenAddress)
    {
        require(isRunning, "stopped bridge");
        require(ver == 1 || ver == 2, "Unknown version");

        if (ver == 1) {
            emit RequestValueTransfer(
                TokenType.ERC721,
                _from,
                _to,
                _tokenAddress,
                _tokenId,
                requestNonce,
                0,
                _extraData
            );
        } else {
            (bool success, bytes memory uri) = _tokenAddress.call(abi.encodePacked(ERC721Metadata(_tokenAddress).tokenURI.selector, abi.encode(_tokenId)));
            if (success == false) {
                uri = "";
            }
            emit RequestValueTransferEncoded(
                ver,
                TokenType.ERC721,
                _from,
                _to,
                _tokenAddress,
                _tokenId,
                requestNonce,
                0,
                abi.encode(string(uri)),
                _extraData
            );
        }
        if (modeMintBurn) {
            ERC721Burnable(_tokenAddress).burn(_tokenId);
        }
        requestNonce++;
    }

    // onERC721Received function of ERC721 token for 1-step deposits to the Bridge
    function onERC721Received(
        address _from,
        uint256 _tokenId,
        address _to,
        bytes memory _extraData
    )
        public
    {
        uint V1 = 1;
        _requestERC721Transfer(msg.sender, _from, _to, _tokenId, _extraData, V1);
    }

    // onERC721ReceivedV2 function is the smae function with onERC721Received, but emits different event that takes uri value.
    function onERC721ReceivedV2(
        address _from,
        uint256 _tokenId,
        address _to,
        bytes memory _extraData
    )
        public
    {
        uint V2 = 2;
        _requestERC721Transfer(msg.sender, _from, _to, _tokenId, _extraData, V2);
    }

    // requestERC721Transfer requests transfer ERC721 to _to on relative chain.
    function requestERC721Transfer(
        address _tokenAddress,
        address _to,
        uint256 _tokenId,
        bytes memory _extraData
    )
        public
    {
        uint V1 = 1;
        IERC721(_tokenAddress).transferFrom(msg.sender, address(this), _tokenId);
        _requestERC721Transfer(_tokenAddress, msg.sender, _to, _tokenId, _extraData, V1);
    }

    // requestERC721TransferV2 is the smae function with requestERC721Transfer, but emits different event that takes uri value.
    function requestERC721TransferV2(
        address _tokenAddress,
        address _to,
        uint256 _tokenId,
        bytes memory _extraData
    )
        public
    {
        uint V2 = 2;
        IERC721(_tokenAddress).transferFrom(msg.sender, address(this), _tokenId);
        _requestERC721Transfer(_tokenAddress, msg.sender, _to, _tokenId, _extraData, V2);
    }
}
