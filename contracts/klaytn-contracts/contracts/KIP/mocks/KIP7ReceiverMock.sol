// SPDX-License-Identifier: MIT
// Klaytn Contract Library v1.0.0 (KIP/mocks/KIP7ReceiverMock.sol)
// Based on OpenZeppelin Contracts v4.5.0 (mocks/ERC20ReceiverMock.sol)
// https://github.com/OpenZeppelin/openzeppelin-contracts/releases/tag/v4.5.0

pragma solidity ^0.8.0;

import "../token/KIP7/IKIP7Receiver.sol";

contract KIP7ReceiverMock is IKIP7Receiver {
    enum Error {
        None,
        RevertWithMessage,
        RevertWithoutMessage,
        Panic
    }

    bytes4 private immutable _retval;
    Error private immutable _error;

    event Received(address operator, address from, uint256 amount, bytes data, uint256 gas);

    constructor(bytes4 retval, Error error) {
        _retval = retval;
        _error = error;
    }

    function onKIP7Received(
        address operator,
        address from,
        uint256 amount,
        bytes memory data
    ) public override returns (bytes4) {
        if (_error == Error.RevertWithMessage) {
            revert("KIP7ReceiverMock: reverting");
        } else if (_error == Error.RevertWithoutMessage) {
            revert();
        } else if (_error == Error.Panic) {
            uint256 a = uint256(0) / uint256(0);
            a;
        }
        emit Received(operator, from, amount, data, gasleft());
        return _retval;
    }
}
