pragma solidity ^0.4.24;

import "./BridgeTransferCommon.sol";


contract BridgeTransferKLAY is BridgeTransfer {
    // handleKLAYTransfer sends the KLAY by the request.
    function handleKLAYTransfer(
        bytes32 _requestTxHash,
        address _from,
        address _to,
        uint256 _value,
        uint64 _requestedNonce,
        uint64 _requestedBlockNumber,
        uint256[] _extraData
    )
        public
        onlyOperators
    {
        if (!voteValueTransfer(_requestedNonce)) {
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
            _extraData
        );
        _setHandledRequestTxHash(_requestTxHash);
        lastHandledRequestBlockNumber = _requestedBlockNumber;

        updateHandleNonce(_requestedNonce);
        _to.transfer(_value);
    }


    // _requestKLAYTransfer requests transfer KLAY to _to on relative chain.
    function _requestKLAYTransfer(address _to, uint256 _feeLimit,  uint256[] _extraData) internal {
        require(isRunning, "stopped bridge");
        require(msg.value > _feeLimit, "insufficient amount");

        uint256 fee = _payKLAYFeeAndRefundChange(_feeLimit);

        emit RequestValueTransfer(
            TokenType.KLAY,
            msg.sender,
            _to,
            address(0),
            msg.value.sub(_feeLimit),
            requestNonce,
            "",
            fee,
            _extraData
        );
        requestNonce++;
    }

    // () requests transfer KLAY to msg.sender address on relative chain.
    function () external payable {
        _requestKLAYTransfer(msg.sender, feeOfKLAY, new uint256[](0));
    }

    // requestKLAYTransfer requests transfer KLAY to _to on relative chain.
    function requestKLAYTransfer(address _to, uint256 _value, uint256[] _extraData) external payable {
        uint256 feeLimit = msg.value.sub(_value);
        _requestKLAYTransfer(_to, feeLimit, _extraData);
    }

    // chargeWithoutEvent sends KLAY to this contract without event for increasing
    // the withdrawal limit.
    function chargeWithoutEvent() external payable {}

    // setKLAYFee set the fee of KLAY transfer.
    function setKLAYFee(uint256 _fee, uint64 _requestNonce)
        external
        onlyOperators
    {
        if (!voteConfiguration(_requestNonce)) {
            return;
        }
        _setKLAYFee(_fee);
    }
}
