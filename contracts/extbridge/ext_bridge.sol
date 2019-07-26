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
        uint256 _amount,
        address _to,
        address _contractAddress,
        uint64 _requestNonce,
        uint64 _requestBlockNumber//,
        //uint256 [] _extraData     // TODO-Klaytn-Servicechain : This will be applied in next PR
    )
        public
    {
        uint256 offerPrice = 1; //_extraData[0];
        if (offerPrice > 0 && callback != address(0)) {
            //super.handleERC20Transfer(_amount, callback, _contractAddress, _requestNonce, _requestBlockNumber, _extraData);
            super.handleERC20Transfer(_amount, callback, _contractAddress, _requestNonce, _requestBlockNumber);
            Callback(callback).RegisterOffer(_to, _amount, _contractAddress, offerPrice);
        } else {
            //super.handleERC20Transfer(_amount, _to, _contractAddress, _requestNonce, _requestBlockNumber, _extraData);
            super.handleERC20Transfer(_amount, _to, _contractAddress, _requestNonce, _requestBlockNumber);
        }
    }

    // handleERC721Transfer sends the ERC721 token by the request and processes the extended feature.
    function handleERC721Transfer(
        uint256 _uid,
        address _to,
        address _contractAddress,
        uint64 _requestNonce,
        uint64 _requestBlockNumber,
        string _tokenURI //,
        //uint256 [] _extraData     : This will be applied in next PR
    )
        public
    {
        uint256 offerPrice = 1; //_extraData[0];
        if (offerPrice > 0 && callback != address(0)) {
            //super.handleERC721Transfer(_uid, _to, _contractAddress, _requestNonce, _requestBlockNumber, _tokenURI, _extraData);
            super.handleERC721Transfer(_uid, _to, _contractAddress, _requestNonce, _requestBlockNumber, _tokenURI);
            Callback(callback).RegisterOffer(_to, _uid, _contractAddress, offerPrice);
        } else {
            //super.handleERC20Transfer(_uid, _to, _contractAddress, _requestNonce, _requestBlockNumber, _extraData);
            super.handleERC721Transfer(_uid, _to, _contractAddress, _requestNonce, _requestBlockNumber, _tokenURI);
        }
    }
}
