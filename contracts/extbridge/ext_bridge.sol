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

pragma solidity ^0.4.24;

import "../bridge/BridgeTransferERC20.sol";
import "../bridge/BridgeTransferERC721.sol";
import "./callback.sol";


// ExtBridge is an extended bridge contract example inherited by BridgeTransferERC20 and BridgeTransferERC721.
// This contract overrides handleERC20Transfer and handleERC721Transfer to make an internal call to callback contract.
contract ExtBridge is BridgeTransferERC20, BridgeTransferERC721 {
    address public callback = address(0);

    constructor(bool _modeMintBurn) BridgeTransfer(_modeMintBurn) public payable {
    }

    function setCallback(address _addr) public onlyOwner {
        callback = _addr;
    }

    // handleERC20Transfer sends the ERC20 token by the request and processes the extended feature.
    function handleERC20Transfer(
        bytes32 _requestTxHash,
        address _from,
        address _to,
        address _tokenAddress,
        uint256 _value,
        uint64 _requestNonce,
        uint64 _requestBlockNumber,
        uint256[] _extraData
    )
        public
    {
        if (_extraData.length > 0) {
            uint256 offerPrice = _extraData[0];
            if (offerPrice > 0 && callback != address(0)) {
                super.handleERC20Transfer(_requestTxHash, _from, callback, _tokenAddress, _value, _requestNonce, _requestBlockNumber, _extraData);
                Callback(callback).registerOffer(_to, _value, _tokenAddress, offerPrice);
                return;
            }
        }

        super.handleERC20Transfer(
            _requestTxHash,
            _from,
            _to,
            _tokenAddress,
            _value,
            _requestNonce,
            _requestBlockNumber,
            _extraData
        );
    }

    // handleERC721Transfer sends the ERC721 token by the request and processes the extended feature.
    function handleERC721Transfer(
        bytes32 _requestTxHash,
        address _from,
        address _to,
        address _tokenAddress,
        uint256 _tokenId,
        uint64 _requestNonce,
        uint64 _requestBlockNumber,
        string _tokenURI,
        uint256[] _extraData
    )
        public
    {
        if (_extraData.length > 0) {
            uint256 offerPrice = _extraData[0];
            if (offerPrice > 0 && callback != address(0)) {
                super.handleERC721Transfer(_requestTxHash, _from, callback, _tokenAddress, _tokenId, _requestNonce, _requestBlockNumber, _tokenURI, _extraData);
                Callback(callback).registerOffer(_to, _tokenId, _tokenAddress, offerPrice);
                return;
            }
        }

        super.handleERC721Transfer(
            _requestTxHash,
            _from,
            _to,
            _tokenAddress,
            _tokenId,
            _requestNonce,
            _requestBlockNumber,
            _tokenURI,
            _extraData
        );
    }
}
