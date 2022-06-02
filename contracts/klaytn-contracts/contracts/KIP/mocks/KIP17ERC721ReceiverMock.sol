// SPDX-License-Identifier: MIT
// Klaytn Contract Library v1.0.0 (KIP/mocks/KIP17ERC721ReceiverMock.sol)
// Based on OpenZeppelin Contracts v4.5.0 (mocks/ERC721EReceiverMock.sol)
// https://github.com/OpenZeppelin/openzeppelin-contracts/releases/tag/v4.5.0

pragma solidity ^0.8.0;

import "../token/KIP17/IKIP17Receiver.sol";
import "../../token/ERC721/IERC721Receiver.sol";

contract KIP17ERC721ReceiverMock is IKIP17Receiver, IERC721Receiver {
    enum Error {
        None,
        RevertWithMessage,
        RevertWithoutMessage,
        Panic
    }

    bytes4 private immutable _retval;
    Error private immutable _error;

    event Received(address operator, address from, uint256 tokenId, bytes data, uint256 gas);

    constructor(bytes4 retval, Error error) {
        _retval = retval;
        _error = error;
    }

    function onKIP17Received(
        address operator,
        address from,
        uint256 tokenId,
        bytes memory data
    ) public override returns (bytes4) {
        return _onKIP7ERC721Received(operator, from, tokenId, data);
    }

    function onERC721Received(
        address operator,
        address from,
        uint256 tokenId,
        bytes memory data
    ) public override returns (bytes4) {
        return _onKIP7ERC721Received(operator, from, tokenId, data);
    }

    function _onKIP7ERC721Received(
        address operator,
        address from,
        uint256 tokenId,
        bytes memory data
    ) internal returns (bytes4) {
        if (_error == Error.RevertWithMessage) {
            revert("KIP17ERC721ReceiverMock: reverting");
        } else if (_error == Error.RevertWithoutMessage) {
            revert();
        } else if (_error == Error.Panic) {
            uint256 a = uint256(0) / uint256(0);
            a;
        }
        emit Received(operator, from, tokenId, data, gasleft());
        return _retval;
    }
}
