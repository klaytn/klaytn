// SPDX-License-Identifier: MIT
// Klaytn Contract Library v1.0.0 (KIP/token/KIP7/extensions/KIP7Pausable.sol)
// Based on OpenZeppelin Contracts v4.5.0 (token/ERC20/extensions/ERC20Pausable.sol)
// https://github.com/OpenZeppelin/openzeppelin-contracts/releases/tag/v4.5.0

pragma solidity ^0.8.0;

import "../KIP7.sol";
import "../../../../security/Pausable.sol";
import "../../../../access/AccessControlEnumerable.sol";
import "../../../interfaces/IKIP7Pausable.sol";

/**
 * @dev Extension of KIP7 that supports permissioned pausable token transfers, minting and burning.
 *
 * Useful for scenarios such as preventing transfers until the end of an evaluation
 * period, or having an emergency switch for freezing all token transfers in the
 * event of a large bug.
 * See http://kips.klaytn.com/KIPs/kip-7-fungible_token
 */
abstract contract KIP7Pausable is KIP7, Pausable, AccessControlEnumerable, IKIP7Pausable {
    bytes32 public constant PAUSER_ROLE = keccak256("KIP7_PAUSER_ROLE");

    /**
     * @dev See {IKIP7Pausable-pause}
     *
     * Emits a {Paused} event.
     *
     * Requirements:
     *
     * - caller must have the {PAUSER_ROLE}
     */
    function pause() public virtual override {
        require(hasRole(PAUSER_ROLE, _msgSender()), "KIP7PresetMinterPauser: must have pauser role to pause");
        _pause();
    }

    /**
     * @dev See {IKIP7Pausable-unpause}
     *
     * Emits a {Unpaused} event.
     *
     * Requirements:
     *
     * - caller must have the {PAUSER_ROLE}
     */
    function unpause() public virtual override {
        require(hasRole(PAUSER_ROLE, _msgSender()), "KIP7PresetMinterPauser: must have pauser role to unpause");
        _unpause();
    }

    /**
     * @dev Returns true if the contract is paused, false otherwise
     */
    function paused() public view override(IKIP7Pausable, Pausable) returns (bool) {
        return super.paused();
    }

    /**
     * @dev Check if `account` has the assigned Pauser role via {AccessControl-hasRole}
     */
    function isPauser(address account) public view returns (bool) {
        return hasRole(PAUSER_ROLE, account);
    }

    /**
     * @dev See {IKIP7Pausable-addPauser}
     *
     * Emits a {RoleGranted} event
     *
     * Requirements:
     *
     * - caller must have the {AccessControl-DEFAULT_ADMIN_ROLE}
     */
    function addPauser(address account) public onlyRole(DEFAULT_ADMIN_ROLE) {
        grantRole(PAUSER_ROLE, account);
    }

    /**
     * @dev Renounce the Pauser role of the caller via {AccessControl-renounceRole}
     *
     * Emits a {RoleRevoked} event
     */
    function renouncePauser() public {
        renounceRole(PAUSER_ROLE, msg.sender);
    }

    /**
     * @dev See {KIP7-_beforeTokenTransfer}.
     *
     * Requirements:
     *
     * - the contract must not be paused.
     */
    function _beforeTokenTransfer(
        address from,
        address to,
        uint256 amount
    ) internal virtual override {
        super._beforeTokenTransfer(from, to, amount);

        require(!paused(), "KIP7Pausable: token transfer while paused");
    }
}
