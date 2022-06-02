// SPDX-License-Identifier: MIT
// Klaytn Contract Library v1.0.0 (KIP/token/KIP37/extensions/IKIP37Pausable.sol)

pragma solidity ^0.8.0;

/**
 * @dev Pausing extension of the KIP17 standard as defined in the KIP.
 * See http://kips.klaytn.com/KIPs/kip-37#pausing-extension
 */
interface IKIP37Pausable {
    /**
     * @dev Returns true if the contract is paused, false otherwise
     */
    function paused() external view returns (bool);

    /**
     * @dev Pause any function which triggers {KIP37-_beforeTokenTransfer}
     *
     * Emits a {Paused} event.
     */
    function pause() external;

    /**
     * @dev Resume normal function from the contract paused state
     *
     * Emits a {Unpaused} event.
     */
    function unpause() external;

    /**
     * @dev Returns true if the token type `id` is paused, false otherwise
     */
    function paused(uint256 id) external view returns (bool);

    /**
     * @dev Pause any function which triggers {KIP37-_beforeTokenTransfer} when token `id` is invoked
     *
     * Emits a {TokenPaused} event.
     */
    function pause(uint256 id) external;

    /**
     * @dev Resume normal function from the token type `id` paused state
     *
     * Emits a {TokenUnpaused} event.
     */
    function unpause(uint256 id) external;
}
