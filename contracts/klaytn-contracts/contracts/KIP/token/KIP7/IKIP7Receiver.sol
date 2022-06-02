// SPDX-License-Identifier: MIT
// Klaytn Contract Library v1.0.0 (KIP/token/KIP7/IKIP7Receiver.sol)
// Based on OpenZeppelin Contracts v4.5.0 (token/ERC20/IERC20Receiver.sol)
// https://github.com/OpenZeppelin/openzeppelin-contracts/releases/tag/v4.5.0

pragma solidity ^0.8.0;

interface IKIP7Receiver {
    /**
     * @dev Whenever an {IKIP7} `amount` is transferred to this contract via {IKIP7-safeTransfer}
     * or {IKIP7-safeTransferFrom} by `operator` from `from`, this function is called.
     *
     * {onKIP7Received} must return its Solidity selector to confirm the token transfer.
     * If any other value is returned or the interface is not implemented by the recipient, the transfer will be reverted.
     *
     * The selector can be obtained in Solidity with `IKIP7Receiver.onKIP7Received.selector`.
     */
    function onKIP7Received(
        address operator,
        address from,
        uint256 amount,
        bytes calldata _data
    ) external returns (bytes4);
}
