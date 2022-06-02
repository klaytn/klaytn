// SPDX-License-Identifier: MIT
// Klaytn Contract Library v1.0.0 (KIP/token/KIP7/extensions/IKIP7Pausable.sol)

pragma solidity ^0.8.0;

/**
 * @dev Pausing extension of the KIP7 standard as defined in the KIP.
 * See https://kips.klaytn.com/KIPs/kip-7#pausing-extension
 */
interface IKIP7Pausable {
    /**
     * @dev Returns true if the contract is paused, false otherwise
     */
    function paused() external view returns (bool);

    /**
     * @dev Pause any function which triggers {KIP7-_beforeTokenTransfer}
     *
     * Emits a {Paused} event.
     *
     * Requirements:
     *
     * - caller must have the {KIP7Pausable-PAUSER_ROLE}
     */
    function pause() external;

    /**
     * @dev Resume normal function from the paused state
     *
     * Emits a {Unpaused} event.
     *
     * Requirements:
     *
     * - caller must have the {KIP7Pausable-PAUSER_ROLE}
     */
    function unpause() external;

    /**
     * @dev Check if `account` has the assigned Pauser role via {AccessControl-hasRole}
     */
    function isPauser(address account) external view returns (bool);

    /**
     * @dev Assign the Pauser role to `account` via {AccessControl-grantRole}
     *
     * Emits a {RoleGranted} event
     *
     * Requirements:
     *
     * - caller must have the {AccessControl-DEFAULT_ADMIN_ROLE}
     */
    function addPauser(address _account) external;

    /**
     * @dev Renounce the Pauser role of the caller via {AccessControl-renounceRole}
     *
     * Emits a {RoleRevoked} event
     */
    function renouncePauser() external;
}
