// SPDX-License-Identifier: MIT
// Klaytn Contract Library v1.0.0 (KIP/token/KIP7/IKIP7.sol)
// Based on OpenZeppelin Contracts v4.5.0 (token/ERC20/IERC20.sol)
// https://github.com/OpenZeppelin/openzeppelin-contracts/releases/tag/v4.5.0

pragma solidity ^0.8.0;

/**
 * @dev Interface of the KIP7 standard as defined in the KIP.
 * See http://kips.klaytn.com/KIPs/kip-7-fungible_token
 */
interface IKIP7 {
    /**
     * @dev Emitted when `value` tokens are moved from one account (`from`) to
     * another (`to`).
     *
     * Note that `value` may be zero.
     */
    event Transfer(address indexed from, address indexed to, uint256 value);

    /**
     * @dev Emitted when the allowance of a `spender` for an `owner` is set by
     * a call to {approve}. `value` is the new allowance.
     */
    event Approval(address indexed owner, address indexed spender, uint256 value);

    /**
     * @dev Returns the amount of tokens in existence.
     */
    function totalSupply() external view returns (uint256);

    /**
     * @dev Returns the amount of tokens owned by `account`.
     */
    function balanceOf(address account) external view returns (uint256);

    /**
     * @dev Moves `amount` tokens from the caller's account to `to`.
     *
     * Returns a boolean value indicating whether the operation succeeded.
     *
     * Emits a {Transfer} event.
     */
    function transfer(address to, uint256 amount) external returns (bool);

    /**
     * @dev Returns the remaining number of tokens that `spender` will be
     * allowed to spend on behalf of `owner` through {transferFrom}. This is
     * zero by default.
     *
     * This value changes when {approve}, {transferFrom}, or {safeTransferFrom} are called.
     */
    function allowance(address owner, address spender) external view returns (uint256);

    /**
     * @dev Sets `amount` as the allowance of `spender` over the caller's tokens.
     *
     * Returns a boolean value indicating whether the operation succeeded.
     *
     * IMPORTANT: Beware that changing an allowance with this method brings the risk
     * that someone may use both the old and the new allowance by unfortunate
     * transaction ordering. One possible solution to mitigate this race
     * condition is to first reduce the spender's allowance to 0 and set the
     * desired value afterwards:
     * https://github.com/ethereum/EIPs/issues/20#issuecomment-263524729
     *
     * Emits an {Approval} event.
     */
    function approve(address spender, uint256 amount) external returns (bool);

    /**
     * @dev Moves `amount` tokens from `from` to `to` using the
     * allowance mechanism. `amount` is then deducted from the caller's
     * allowance.
     *
     * Returns a boolean value indicating whether the operation succeeded.
     *
     * Emits a {Transfer} event.
     */
    function transferFrom(
        address from,
        address to,
        uint256 amount
    ) external returns (bool);

    /**
     * @dev Moves `amount` tokens from the caller's account to `recipient`
     * and passes `data` for {IKIP7Receiver-onKIP7Received} handler logic.
     *
     * Emits a {Transfer} event.
     */
    function safeTransfer(
        address recipient,
        uint256 amount,
        bytes memory data
    ) external;

    /**
     * @dev  Moves `amount` tokens from the caller's account to `recipient`.
     *
     * Emits a {Transfer} event.
     */
    function safeTransfer(address recipient, uint256 amount) external;

    /**
     * @dev Moves `amount` tokens from `sender` to `recipient` using the {allowance} mechanism
     * and passes `data` for {IKIP7Receiver-onKIP7Received} handler logic.
     *
     * Emits a {Transfer} event.
     */
    function safeTransferFrom(
        address sender,
        address recipient,
        uint256 amount,
        bytes memory data
    ) external;

    /**
     * @dev Moves `amount` tokens from `sender` to `recipient` using the {allowance} mechanism.
     *
     * Emits a {Transfer} event.
     */
    function safeTransferFrom(
        address sender,
        address recipient,
        uint256 amount
    ) external;
}
