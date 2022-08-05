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

import "../externals/openzeppelin-solidity/contracts/math/SafeMath.sol";

import "./BridgeHandledRequests.sol";
import "./BridgeFee.sol";
import "./BridgeOperator.sol";
import "./BridgeTokens.sol";

contract BridgeTransfer is BridgeHandledRequests, BridgeFee, BridgeOperator {
    bool public modeMintBurn = false;
    bool public isRunning = true;

    uint64 public requestNonce; // the number of value transfer request that this contract received.
    uint64 public lowerHandleNonce; // a minimum nonce of a value transfer request that will be handled.
    uint64 public upperHandleNonce; // a maximum nonce of the counterpart bridge's value transfer request that is handled.
    uint64 public recoveryBlockNumber = 1; // the block number that recovery start to filter log from.
    uint64 public failedHandles; // a count of failed handle value transfer
    uint64 public nRefunds; // a count of refunds
    mapping(uint64 => uint64) public handleNoncesToBlockNums;  // <request nonce> => <request blockNum>

    mapping(uint64 => address payable) public refundAddrMap; // <nonce of request, sender>
    mapping(uint64 => uint256) public refundValueMap; // <nonce of request, value>
    mapping(uint64 => TokenType) public refundTokenType; // <nonce of request, token type>
    mapping(uint64 => uint) public refundTimestampMap; // <nonce of request, token type>

    string public REMOVED_VOTE_ERR_STR = "removed vote";

    using SafeMath for uint256;

    enum TokenType {
        KLAY,
        ERC20,
        ERC721
    }

    constructor(bool _modeMintBurn) BridgeFee(address(0)) internal {
        modeMintBurn = _modeMintBurn;
    }

    // start can allow or disallow the value transfer request.
    function start(bool _status)
        external
        onlyOwner
    {
        isRunning = _status;
    }

    /**
     * Event to log the request value transfer from the Bridge.
     * @param tokenType is the type of tokens (KLAY/ERC20/ERC721).
     * @param from is the requester of the request value transfer event.
     * @param to is the receiver of the value.
     * @param tokenAddress Address of token contract the token belong to.
     * @param valueOrTokenId is the value of KLAY/ERC20 or token ID of ERC721.
     * @param requestNonce is the order number of the request value transfer.
     * @param fee is fee of value transfer.
     * @param extraData is additional data for specific purpose of a service provider.
     */
    event RequestValueTransfer(
        TokenType tokenType,
        address indexed from,
        address indexed to,
        address indexed tokenAddress,
        uint256 valueOrTokenId,
        uint64 requestNonce,
        uint256 fee,
        bytes extraData
    );

    /**
     * Event to log the request value transfer from the Bridge.
     * @param tokenType is the type of tokens (KLAY/ERC20/ERC721).
     * @param from is the requester of the request value transfer event.
     * @param to is the receiver of the value.
     * @param tokenAddress Address of token contract the token belong to.
     * @param valueOrTokenId is the value of KLAY/ERC20 or token ID of ERC721.
     * @param requestNonce is the order number of the request value transfer.
     * @param fee is fee of value transfer.
     * @param extraData is additional data for specific purpose of a service provider.
     * @param encodingVer indicates encodedData version.
     * @param encodedData is a packed set of values.
     */
    event RequestValueTransferEncoded(
        TokenType tokenType,
        address indexed from,
        address indexed to,
        address indexed tokenAddress,
        uint256 valueOrTokenId,
        uint64 requestNonce,
        uint256 fee,
        bytes extraData,
        uint8 encodingVer,
        bytes encodedData
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
        address indexed from,
        address indexed to,
        address indexed tokenAddress,
        uint256 valueOrTokenId,
        uint64 handleNonce,
        uint64 lowerHandleNonce,
        bytes extraData
    );

    /**
     * Emitted to deliver a refund request to established counterpart chain
     * @param requestNonce is the order number of the request value transfer.
     */
    event RequestRefund(uint64 indexed requestNonce, bytes32 indexed requestTxHash);

    /** Emitted when a refund request is handled
     * Event to log the refund by failed handle value transfer.
     * @param requestNonce is the order number of the request value transfer.
     * @param sender is the address of sedner that created request value transfer.
     * @param value amount of KLAY.
     */
    event HandleRefund(uint64 indexed requestNonce, address indexed sender, uint256 indexed value);

    // updateHandleNonce increases lower and upper handle nonce after the _requestedNonce is handled.
    // DO NOT OPERATE WITH OTHER CONTRACT ADDRESS IN THIS FUNCTION
    function _updateHandleNonce(uint64 _requestedNonce) internal {
        if (_requestedNonce > upperHandleNonce) {
            upperHandleNonce = _requestedNonce;
        }

        uint64 limit = lowerHandleNonce + 200;
        if (limit > upperHandleNonce) {
            limit = upperHandleNonce;
        }

        uint64 i;
        for (i = lowerHandleNonce; i <= limit && handleNoncesToBlockNums[i] > 0; i++) {
            recoveryBlockNumber = handleNoncesToBlockNums[i];
            delete handleNoncesToBlockNums[i];
            delete closedValueTransferVotes[i];
        }
        lowerHandleNonce = i;
    }

    function _lowerHandleNonceCheck(uint64 _requestedNonce) internal view {
        require(lowerHandleNonce <= _requestedNonce, REMOVED_VOTE_ERR_STR);
    }

    // setFeeReceivers sets fee receiver.
    function setFeeReceiver(address payable _feeReceiver)
        external
        onlyOwner
    {
        _setFeeReceiver(_feeReceiver);
    }

    function updateHandleStatus(uint64 _requestNonce, bytes32 _requestTxHash, uint64 _requestedBlockNumber, bool failed)
        public
        onlyOperators
        returns (bool)
    {
        // Check needed (by replay attack, by malicious operator)
        _lowerHandleNonceCheck(_requestNonce);

        // Check needed (by replay attack (by malicious operator))
        if (!_voteValueTransfer(_requestNonce)) {
            return false;
        }

        _setHandledRequestTxHash(_requestTxHash);

        handleNoncesToBlockNums[_requestNonce] = _requestedBlockNumber;
        _updateHandleNonce(_requestNonce);

        if (failed) {
            failedHandles += 1;
        }
        return true;
    }

    // setRefundLedger sets refund ledger (sender and value)
    function setRefundLedger(uint64 requestNonce_, uint256 value, TokenType tokenType) internal {
        refundAddrMap[requestNonce_] = msg.sender;
        refundValueMap[requestNonce_] = value;
        refundTokenType[requestNonce_] = tokenType;
        refundTimestampMap[requestNonce_] = block.timestamp;
        amountOfLockedRefundKLAY += value;
    }

    // requestRefund emit `RequestRefund` event to be watched and handled by established counterpart chain
    function requestRefund(uint64 requestNonce_, bytes32 requestTxHash) public onlyOperators {
        if (!refundCond(requestNonce_)) {
            return;
        }
        emit RequestRefund(requestNonce_, requestTxHash);
    }

    // handleRefund refunds the requested amount of KLAY or tokens
    function handleRefund(uint64 requestNonce_) public onlyOperators {
        if (!refundCond(requestNonce_)) {
            return;
        }

        address payable sender = refundAddrMap[requestNonce_];
        uint256 value = refundValueMap[requestNonce_];
        if (sender != address(0) && value != 0) {
            if (refundTokenType[requestNonce_] == TokenType.KLAY) {
                sender.transfer(value);
            } else if (refundTokenType[requestNonce_] == TokenType.ERC20) {
                revert("unimplemented ERC20 refund");
            } else if (refundTokenType[requestNonce_] == TokenType.ERC721) {
                revert("unimplemented ERC721 refund");
            }
            nRefunds += 1;
            emit HandleRefund(requestNonce_, sender, value);
            removeRefundInfo(requestNonce_, value);
        }
    }

    function refundCond(uint64 requestNonce_) internal returns (bool) {
        /*
        // if a refund was not correctly executed during one day (by gas error of operators and etc),
        // the sender can have authroization to refund himself/herself
        if (block.timestamp > refundTimestampMap[requestNonce_] + 1 days) {
            return true;
        }
        */

        require(operators[msg.sender], "msg.sender is not an operator");
        if (!_voteRefund(requestNonce_)) {
            return false;
        }
        return true;
    }

    // removeRefundLedger removes refund info (sender address and amount of KLAY).
    // This function is only called when a `HandleValueTransfer` event is emitted.
    // DO NOT OPERATE WITH OTHER CONTRACT ADDRESS IN THIS FUNCTION
    function removeRefundLedger(uint64 requestNonce_) public onlyOperators {
        if (!_voteRefund(requestNonce_)) {
            return;
        }
        removeRefundInfo(requestNonce_, refundValueMap[requestNonce_]);
    }

    function removeRefundInfo(uint64 requestNonce_, uint256 value) internal {
        amountOfLockedRefundKLAY = amountOfLockedRefundKLAY.sub(value);
        delete refundAddrMap[requestNonce_];
        delete refundValueMap[requestNonce_];
        delete refundTimestampMap[requestNonce_];
    }
}
