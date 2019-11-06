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

    // requestSellERC20 requests transfer ERC20 to _to on relative chain to sell it.
    function requestSellERC20(
        address _tokenAddress,
        address _to,
        uint256 _value,
        uint256 _feeLimit,
        uint256 _price
    )
    external
    {
        super.requestERC20Transfer(
            _tokenAddress,
            _to,
            _value,
            _feeLimit,
            abi.encode(_price)
        );
    }

    // requestERC20Transfer requests transfer ERC20 to _to on relative chain.
    function requestERC20Transfer(
        address _tokenAddress,
        address _to,
        uint256 _value,
        uint256 _feeLimit,
        bytes memory _extraData
    )
    public
    {
        revert("not support");
    }

    // requestSellERC721 requests transfer ERC721 to _to on relative chain to sell it.
    function requestSellERC721(
        address _tokenAddress,
        address _to,
        uint256 _tokenId,
        uint256 _price
    )
    external
    {
        super.requestERC721Transfer(
            _tokenAddress,
            _to,
            _tokenId,
            abi.encode(_price)
        );
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
        revert("not support");
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
        bytes memory _extraData
    )
        public
    {
        require(_extraData.length == 32, "extraData size error");

        require(callback != address(0), "callback address error");

        uint256 offerPrice = abi.decode(_extraData, (uint256));
        require(offerPrice > 0, "offerPrice error");

        super.handleERC20Transfer(_requestTxHash, _from, callback, _tokenAddress, _value, _requestNonce, _requestBlockNumber, _extraData);
        Callback(callback).registerOffer(_to, _value, _tokenAddress, offerPrice);
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
        string memory _tokenURI,
        bytes memory _extraData
    )
        public
    {
        require(_extraData.length == 32, "extraData size error");

        require(callback != address(0), "callback address error");

        uint256 offerPrice = abi.decode(_extraData, (uint256));
        require(offerPrice > 0, "offerPrice error");

        super.handleERC721Transfer(_requestTxHash, _from, callback, _tokenAddress, _tokenId, _requestNonce, _requestBlockNumber, _tokenURI, _extraData);
        Callback(callback).registerOffer(_to, _tokenId, _tokenAddress, offerPrice);
    }
}
