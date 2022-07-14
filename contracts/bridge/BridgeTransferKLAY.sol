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

import "./BridgeTransfer.sol";


contract BridgeTransferKLAY is BridgeTransfer {
    mapping(uint64 => address payable) private refundAddrMap; // <nonce of request, sender>
    mapping(uint64 => uint256) private refundValueMap; // <nonce of request, value>
    bool public isLockedKLAY;

    event KLAYLocked();
    event KLAYUnlocked();
    event Refunded(uint64 indexed requestNonce, address indexed sender, uint256 value);

    modifier lockedKLAY {
        require(isLockedKLAY == true, "unlocked");
        _;
    }

    modifier unlockedKLAY {
        require(isLockedKLAY == false, "locked");
        _;
    }

    // lockKLAY can to prevent request KLAY transferring.
    function lockKLAY()
        external
        onlyOwner
        unlockedKLAY
    {
        isLockedKLAY = true;

        emit KLAYLocked();
    }

    // unlockToken can allow request KLAY transferring.
    function unlockKLAY()
        external
        onlyOwner
        lockedKLAY
    {
        isLockedKLAY = false;

        emit KLAYUnlocked();
    }

    // suggestRecommendedFee calculates minimum fee to be used for value transfer fee.
    // The calculated value is larger than three cost of transaction (two of handle value trasnfer call and one refund call)
    // 1. Normal case (no error): the operator pays the value transfer transaction cost in the counterpart chain.
    // 2. Hardening case: the operator pays two value transfer transaction cost
    // 3. Refund case: the operator pays two value transfer transaction cost in the counterpart chain and one refund call
    // Warning: The ServiceChian package would not provide a control for a number of gas per specific opcode.
    //          A modified binary that changed a number of gas for corresponding opcode would have different `gasUsed` (i.e., different with precomputed(hardcoded) cost)
    function suggestLeastFee(uint256 gasPrice) external pure returns (uint256) {
        uint256 complement = 30000; // To make conservative upper bound of expected tx cost
        uint256 precomputedGasUsedOfRefundCall = 110000 + complement;
        uint256 precomputedGasUsedOfHandleValueTransferCall = 150000 + complement;
        return gasPrice * (precomputedGasUsedOfHandleValueTransferCall * 2 + precomputedGasUsedOfRefundCall);
    }

    // refundKLAYTransfer refunds the requested amount of KLAY to sender if its corresponding value transfer is failed from the bridge contract of counterpart chain
    function refundKLAYTransfer(uint64 requestNonce) public onlyOperators {
        if (!_voteRefund(requestNonce)) {
            return;
        }
        address payable sender = refundAddrMap[requestNonce];
        uint256 value = refundValueMap[requestNonce];
        if (sender != address(0) && value != 0) {
            sender.transfer(value);
            emit Refunded(requestNonce, sender, value);
        }
    }

    // handleKLAYTransfer sends the KLAY by the request.
    function handleKLAYTransfer(
        bytes32 _requestTxHash,
        address _from,
        address payable _to,
        uint256 _value,
        uint64 _requestedNonce,
        uint64 _requestedBlockNumber,
        bytes memory _extraData
    )
        public
        onlyOperators
    {
        // Check needed (by replay attack, by malicious operator)
        _lowerHandleNonceCheck(_requestedNonce);

        // Check needed (by replay attack (by malicious operator))
        if (!_voteValueTransfer(_requestedNonce)) {
            return;
        }

        _setHandledRequestTxHash(_requestTxHash);

        handleNoncesToBlockNums[_requestedNonce] = _requestedBlockNumber;
        _updateHandleNonce(_requestedNonce);

        emit HandleValueTransfer(
            _requestTxHash,
            TokenType.KLAY,
            _from,
            _to,
            address(0),
            _value,
            _requestedNonce,
            lowerHandleNonce,
            _extraData
        );

        _to.transfer(_value);
    }

    // _requestKLAYTransfer requests transfer KLAY to _to on relative chain.
    function _requestKLAYTransfer(address _to, uint256 _feeLimit,  bytes memory _extraData)
        internal
        unlockedKLAY
    {
        require(isRunning, "stopped bridge");
        require(msg.value > _feeLimit, "insufficient amount");

        uint256 fee = _payKLAYFeeAndRefundChange(_feeLimit);
        uint256 amountOfKLAY = msg.value.sub(_feeLimit);

        emit RequestValueTransfer(
            TokenType.KLAY,
            msg.sender,
            _to,
            address(0),
            amountOfKLAY,
            requestNonce,
            fee,
            _extraData
        );
        refundAddrMap[requestNonce] = msg.sender;
        refundValueMap[requestNonce] = amountOfKLAY;
        requestNonce++;
    }

    // () requests transfer KLAY to msg.sender address on relative chain.
    function () external payable {
        _requestKLAYTransfer(msg.sender, feeOfKLAY, new bytes(0));
    }

    // requestKLAYTransfer requests transfer KLAY to _to on relative chain.
    function requestKLAYTransfer(address _to, uint256 _value, bytes calldata _extraData) external payable {
        uint256 feeLimit = msg.value.sub(_value);
        _requestKLAYTransfer(_to, feeLimit, _extraData);
    }

    // chargeWithoutEvent sends KLAY to this contract without event for increasing
    // the withdrawal limit.
    function chargeWithoutEvent() external payable {}

    // getMinimumAmountOfKLAY returns minimum amount of KLAY to be successfully executed
    function getMinimumAmountOfKLAY(uint256 value) public view returns (uint256) {
        if (feeReceiver == address(0)) {
            return value;
        } else {
            return value + feeOfKLAY;
        }
    }

    // setKLAYFee set the fee of KLAY transfer.
    function setKLAYFee(uint256 _fee, uint64 _requestNonce)
        external
        onlyOperators
    {
        if (!_voteConfiguration(_requestNonce)) {
            return;
        }
        _setKLAYFee(_fee);
    }
}
