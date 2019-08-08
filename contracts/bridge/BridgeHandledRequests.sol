pragma solidity ^0.4.24;


contract BridgeHandledRequests {
    // TODO-Klaytn-Servicechain handleTxHash can be saved after Klaytn supports it.
    mapping(bytes32 => bool) public handledRequestTx;

    function _setHandledRequestTxHash(bytes32 _requestTxHash) internal {
        handledRequestTx[_requestTxHash] = true;
    }
}
