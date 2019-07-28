pragma solidity ^0.4.24;

import "../bridge/Bridge.sol";
import "./callback.sol";

contract ExtBridge is Bridge {
    address public callback = address(0);

    constructor (bool _modeMintBurn) Bridge(_modeMintBurn) public payable {
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
        uint64 _requestBlockNumber
        //uint256 [] _extraData     // TODO-Klaytn-Servicechain : This will be applied in next PR
    )
        public
    {
        uint256 offerPrice = 1; //_extraData[0];
        if (offerPrice > 0 && callback != address(0)) {
            //super.handleERC20Transfer(_requestTxHash, _from, callback, _tokenAddress, _value, _requestNonce, _requestBlockNumber, _extraData);
            super.handleERC20Transfer(_requestTxHash, _from, callback, _tokenAddress, _value, _requestNonce, _requestBlockNumber);
            Callback(callback).RegisterOffer(_to, _value, _tokenAddress, offerPrice);
        } else {
            //super.handleERC20Transfer(_requestTxHash, _from, _to, _tokenAddress, _value, _requestNonce, _requestBlockNumber, _extraData);
            super.handleERC20Transfer(_requestTxHash, _from, _to, _tokenAddress, _value, _requestNonce, _requestBlockNumber);
        }
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
        string _tokenURI //,
        //uint256 [] _extraData     : This will be applied in next PR
    )
        public
    {
        uint256 offerPrice = 1; //_extraData[0];
        if (offerPrice > 0 && callback != address(0)) {
            //super.handleERC721Transfer(_requestTxHash, _from, callback, _tokenAddress, _tokenId, _requestNonce, _requestBlockNumber, _tokenURI, _extraData);
            super.handleERC721Transfer(_requestTxHash, _from, callback, _tokenAddress, _tokenId, _requestNonce, _requestBlockNumber, _tokenURI);
            Callback(callback).RegisterOffer(_to, _tokenId, _tokenAddress, offerPrice);
        } else {
            //super.handleERC20Transfer(_requestTxHash, _from, _to, _tokenAddress,  _tokenId, _requestNonce, _requestBlockNumber, _tokenURI, _extraData);
            super.handleERC721Transfer(_requestTxHash, _from, _to, _tokenAddress,  _tokenId, _requestNonce, _requestBlockNumber, _tokenURI);
        }
    }
}
