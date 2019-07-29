pragma solidity ^0.4.24;

contract HandledRequests {
    // TODO-Klaytn-Servicechain handleTxHash can be saved after Klaytn supports it.
    mapping (bytes32 => bool) public isHandledRequestTx;

    function _setHandledRequestTxHash(bytes32 _requestTxHash) internal {
        isHandledRequestTx[_requestTxHash] = true;
    }
}
