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
    bool public isLockedKLAY;

    event KLAYLocked();
    event KLAYUnlocked();

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
        if (!updateHandleStatus(_requestedNonce, _requestTxHash, _requestedBlockNumber, false)) {
            return;
        }

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
        setRefundLedger(requestNonce, amountOfKLAY, TokenType.KLAY);
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
