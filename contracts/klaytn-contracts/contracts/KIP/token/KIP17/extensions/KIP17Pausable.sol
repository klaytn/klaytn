// SPDX-License-Identifier: MIT
// Klaytn Contract Library v1.0.0 (KIP/token/KIP17/extensions/KIP17Pausable.sol)
// Based on OpenZeppelin Contracts v4.5.0 (token/ERC721/extensions/ERC721Pausable.sol)
// https://github.com/OpenZeppelin/openzeppelin-contracts/releases/tag/v4.5.0

pragma solidity ^0.8.0;

import "../KIP17.sol";
import "../../../../security/Pausable.sol";
import "../../../../access/AccessControlEnumerable.sol";
import "./IKIP17Pausable.sol";

/**
 * @dev KIP17 token with pausable token transfers, minting and burning.
 *
 * Useful for scenarios such as preventing trades until the end of an evaluation
 * period, or having an emergency switch for freezing all token transfers in the
 * event of a large bug.
 */
abstract contract KIP17Pausable is KIP17, Pausable, AccessControlEnumerable, IKIP17Pausable {
    bytes32 public constant PAUSER_ROLE = keccak256("KIP17_PAUSER_ROLE");

    /**
     * @dev Returns true if `interfaceId` is implemented and false otherwise
     *
     * See {IKIP13} and {IERC165}.
     */
    function supportsInterface(bytes4 interfaceId)
        public
        view
        virtual
        override(KIP17, AccessControlEnumerable)
        returns (bool)
    {
        return
            interfaceId == type(IKIP17Pausable).interfaceId ||
            KIP17.supportsInterface(interfaceId) ||
            AccessControlEnumerable.supportsInterface(interfaceId);
    }

    /**
     * @dev See {IKIP17Pausable-pause}
     *
     * Emits a {Paused} event.
     *
     * Requirements:
     *
     * - caller must have the {PAUSER_ROLE}
     */
    function pause() public override {
        require(hasRole(PAUSER_ROLE, _msgSender()), "KIP17Pausable: must have pauser role to pause");
        _pause();
    }

    /**
     * @dev See {IKIP17Pausable-unpause}
     *
     * Emits a {Unpaused} event.
     *
     * Requirements:
     *
     * - caller must have the {PAUSER_ROLE}
     */
    function unpause() public override {
        require(hasRole(PAUSER_ROLE, _msgSender()), "KIP17Pausable: must have pauser role to unpause");
        _unpause();
    }

    /**
     * @dev Returns true if the contract is paused, false otherwise
     */
    function paused() public view override(IKIP17Pausable, Pausable) returns (bool) {
        return super.paused();
    }

    /**
     * @dev Check if `account` has the assigned Pauser role via {AccessControl-hasRole}
     */
    function isPauser(address _account) public view returns (bool) {
        return hasRole(PAUSER_ROLE, _account);
    }

    /**
     * @dev See {IKIP17Pausable-addPauser}
     *
     * Emits a {RoleGranted} event
     *
     * Requirements:
     *
     * - caller must have the {AccessControl-DEFAULT_ADMIN_ROLE}
     */
    function addPauser(address _account) public onlyRole(DEFAULT_ADMIN_ROLE) {
        grantRole(PAUSER_ROLE, _account);
    }

    /**
     * @dev See {IKIP17Pausable-renouncePauser}
     *
     * Emits a {RoleRevoked} event
     */
    function renouncePauser() public {
        renounceRole(PAUSER_ROLE, msg.sender);
    }

    /**
     * @dev See {KIP17-_beforeTokenTransfer}.
     *
     * Requirements:
     *
     * - the contract must not be paused.
     */
    function _beforeTokenTransfer(
        address from,
        address to,
        uint256 tokenId
    ) internal virtual override {
        super._beforeTokenTransfer(from, to, tokenId);

        require(!paused(), "KIP17Pausable: token transfer while paused");
    }
}
