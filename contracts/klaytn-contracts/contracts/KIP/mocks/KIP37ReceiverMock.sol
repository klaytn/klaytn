// SPDX-License-Identifier: MIT
// Klaytn Contract Library v1.0.0 (KIP/mocks/KIP37ReceiverMock.sol)
// Based on OpenZeppelin Contracts v4.5.0 (mocks/ERC1155ReceiverMock.sol)
// https://github.com/OpenZeppelin/openzeppelin-contracts/releases/tag/v4.5.0

pragma solidity ^0.8.0;

import "../utils/introspection/KIP13.sol";
import "../token/KIP37/IKIP37Receiver.sol";

contract KIP37ReceiverMock is KIP13, IKIP37Receiver {
    bytes4 private _recRetval;
    bool private _recReverts;
    bytes4 private _batRetval;
    bool private _batReverts;

    enum Error {
        None,
        RevertWithMessage,
        RevertWithoutMessage,
        Panic
    }

    Error private _error;
    Error private _batError;

    event Received(address operator, address from, uint256 id, uint256 amount, bytes data, uint256 gas);
    event BatchReceived(address operator, address from, uint256[] ids, uint256[] amounts, bytes data, uint256 gas);

    constructor(
        bytes4 recRetval,
        bool recReverts,
        bytes4 batRetval,
        bool batReverts,
        Error error_,
        Error batError
    ) {
        _recRetval = recRetval;
        _recReverts = recReverts;
        _batRetval = batRetval;
        _batReverts = batReverts;
        _error = error_;
        _batError = batError;
    }

    function onKIP37Received(
        address operator,
        address from,
        uint256 id,
        uint256 amount,
        bytes calldata data
    ) external override returns (bytes4) {
        if (_error == Error.RevertWithMessage) {
            revert("KIP37ReceiverMock: reverting");
        } else if (_error == Error.RevertWithoutMessage) {
            revert();
        } else if (_error == Error.Panic) {
            uint256 a = uint256(0) / uint256(0);
            a;
        }
        emit Received(operator, from, id, amount, data, gasleft());
        return _recRetval;
    }

    function onKIP37BatchReceived(
        address operator,
        address from,
        uint256[] calldata ids,
        uint256[] calldata amounts,
        bytes calldata data
    ) external override returns (bytes4) {
        if (_batError == Error.RevertWithMessage) {
            revert("KIP37ReceiverMock: reverting");
        } else if (_batError == Error.RevertWithoutMessage) {
            revert();
        } else if (_batError == Error.Panic) {
            uint256 a = uint256(0) / uint256(0);
            a;
        }
        emit BatchReceived(operator, from, ids, amounts, data, gasleft());
        return _batRetval;
    }
}
