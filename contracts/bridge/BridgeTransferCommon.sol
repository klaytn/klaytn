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

pragma solidity ^0.5.6;

import "../externals/openzeppelin-solidity/contracts/math/SafeMath.sol";

import "./BridgeHandledRequests.sol";
import "./BridgeFee.sol";
import "./BridgeOperator.sol";


contract BridgeTransfer is BridgeHandledRequests, BridgeFee, BridgeOperator {
    bool public modeMintBurn = false;
    bool public isRunning;

    uint64 public requestNonce; // the number of value transfer request that this contract received.
    uint64 public lowerHandleNonce; // a minimum nonce of a value transfer request that will be handled.
    uint64 public upperHandleNonce; // a maximum nonce of the counterpart bridge's value transfer request that is handled.
    uint64 public recoveryBlockNumber = 1; // the block number that recovery start to filter log from.
    mapping(uint64 => uint64) public handleNoncesToBlockNums;  // <request nonce> => <request blockNum>

    using SafeMath for uint256;

    mapping(address => address) public allowedTokens; // <token, counterpart token>

    enum TokenType {
        KLAY,
        ERC20,
        ERC721
    }

    constructor(bool _modeMintBurn) BridgeFee(address(0)) internal {
        modeMintBurn = _modeMintBurn;
        isRunning = true;
    }

    // start can allow or disallow the value transfer request.
    function start(bool _status)
        external
        onlyOwner
    {
        isRunning = _status;
    }

    // registerToken can update the allowed token with the counterpart token.
    function registerToken(address _token, address _cToken)
        external
        onlyOwner
    {
        allowedTokens[_token] = _cToken;
    }

    // deregisterToken can remove the token in allowedToken list.
    function deregisterToken(address _token)
        external
        onlyOwner
    {
        delete allowedTokens[_token];
    }

    /**
     * Event to log the request value transfer from the Bridge.
     * @param tokenType is the type of tokens (KLAY/ERC20/ERC721).
     * @param from is the requester of the request value transfer event.
     * @param to is the receiver of the value.
     * @param tokenAddress Address of token contract the token belong to.
     * @param valueOrTokenId is the value of KLAY/ERC20 or token ID of ERC721.
     * @param requestNonce is the order number of the request value transfer.
     * @param uri is uri of ERC721 token.
     * @param fee is fee of value transfer.
     * @param extraData is additional data for specific purpose of a service provider.
     */
    event RequestValueTransfer(
        TokenType tokenType,
        address from,
        address to,
        address tokenAddress,
        uint256 valueOrTokenId,
        uint64 requestNonce,
        string uri,
        uint256 fee,
        bytes extraData
    );

    /**
     * Event to log the handle value transfer from the Bridge.
     * @param requestTxHash is a transaction hash of request value transfer.
     * @param tokenType is the type of tokens (KLAY/ERC20/ERC721).
     * @param from is an address of the account who requested the value transfer.
     * @param to is an address of the account who will received the value.
     * @param tokenAddress Address of token contract the token belong to.
     * @param valueOrTokenId is the value of KLAY/ERC20 or token ID of ERC721.
     * @param handleNonce is the order number of the handle value transfer.
     * @param extraData is additional data for specific purpose of a service provider.
     */
    event HandleValueTransfer(
        bytes32 requestTxHash,
        TokenType tokenType,
        address from,
        address to,
        address tokenAddress,
        uint256 valueOrTokenId,
        uint64 handleNonce,
        bytes extraData

    );

    // updateHandleNonce increases lower and upper handle nonce after the _requestedNonce is handled.
    function updateHandleNonce(uint64 _requestedNonce, uint64 _requestBlockNumber) internal {
        uint64 i;
        handleNoncesToBlockNums[_requestedNonce] = _requestBlockNumber;

        if (_requestedNonce > upperHandleNonce) {
            upperHandleNonce = _requestedNonce;
        }
        for (i = lowerHandleNonce; i <= upperHandleNonce && handleNoncesToBlockNums[i] > 0; i++) { }
        lowerHandleNonce = i;
        if (i != 0) {
            recoveryBlockNumber = handleNoncesToBlockNums[i-1];
        }
    }

    // setFeeReceivers sets fee receiver.
    function setFeeReceiver(address payable _feeReceiver)
        external
        onlyOwner
    {
        _setFeeReceiver(_feeReceiver);
    }
}
