// SPDX-License-Identifier: MIT
// Klaytn Contract Library v1.0.0 (KIP/token/KIP37/extensions/KIP37Pausable.sol)
// Based on OpenZeppelin Contracts v4.5.0 (token/ERC1155/extensions/ERC1155Pausable.sol)
// https://github.com/OpenZeppelin/openzeppelin-contracts/releases/tag/v4.5.0

pragma solidity ^0.8.0;

import "../KIP37.sol";
import "../../../../security/Pausable.sol";
import "../../../../access/AccessControlEnumerable.sol";
import "./IKIP37Pausable.sol";

/**
 * @dev KIP37 token with contract wide (or token specific) pausable token transfers, minting and burning.
 *
 * Useful for scenarios such as preventing trades until the end of an evaluation
 * period, or having an emergency switch for freezing all token transfers in the
 * event of a large bug.
 */
abstract contract KIP37Pausable is KIP37, Pausable, AccessControlEnumerable, IKIP37Pausable {
    bytes32 public constant PAUSER_ROLE = keccak256("KIP37_PAUSER_ROLE");

    mapping(uint256 => bool) private _tokenPaused;

    /**
     * @dev Emitted when a token type is paused
     */
    event TokenPaused(address account, uint256 id);

    /**
     * @dev Emitted when a token type is unpaused
     */
    event TokenUnpaused(address account, uint256 id);

    /**
     * @dev Returns true if `interfaceId` is implemented and false otherwise
     *
     * See {IKIP13} and {IERC165}.
     */
    function supportsInterface(bytes4 interfaceId)
        public
        view
        virtual
        override(KIP37, AccessControlEnumerable)
        returns (bool)
    {
        return
            interfaceId == type(IKIP37Pausable).interfaceId ||
            KIP37.supportsInterface(interfaceId) ||
            AccessControlEnumerable.supportsInterface(interfaceId);
    }

    /**
     * @dev See {IKIP37Pausable-pause}
     *
     * Emits a {Paused} event.
     *
     * Requirements:
     *
     * - caller must have the {PAUSER_ROLE}
     */
    function pause() public virtual override onlyRole(PAUSER_ROLE) {
        _pause();
    }

    /**
     * @dev See {IKIP37Pausable-unpause}
     *
     * Emits a {Unpaused} event.
     *
     * Requirements:
     *
     * - caller must have the {PAUSER_ROLE}
     */
    function unpause() public virtual override onlyRole(PAUSER_ROLE) {
        _unpause();
    }

    /**
     * @dev Returns true if the contract is paused, false otherwise
     */
    function paused() public view override(IKIP37Pausable, Pausable) returns (bool) {
        return Pausable.paused();
    }

    /**
     * @dev Returns true if the token type `id` is paused, false otherwise
     */
    function paused(uint256 id) public view override returns (bool) {
        return _tokenPaused[id];
    }

    /**
     * @dev See {IKIP37Pausable-pause(uint256 id)}
     *
     * Emits a {TokenPaused} event.
     *
     * Requirements:
     *
     * - `id` must not be already paused
     * - caller must have the {PAUSER_ROLE}
     */
    function pause(uint256 id) public virtual onlyRole(PAUSER_ROLE) {
        require(_tokenPaused[id] == false, "KIP37Pausable: token already paused");
        _tokenPaused[id] = true;
        emit TokenPaused(_msgSender(), id);
    }

    function unpause(uint256 id) public virtual onlyRole(PAUSER_ROLE) {
        require(_tokenPaused[id] == true, "KIP37Pausable: token already unpaused");
        _tokenPaused[id] = false;
        emit TokenUnpaused(_msgSender(), id);
    }

    /**
     * @dev See {KIP37-_beforeTokenTransfer}.
     *
     * Requirements:
     *
     * - the contract must not be paused.
     */
    function _beforeTokenTransfer(
        address operator,
        address from,
        address to,
        uint256[] memory ids,
        uint256[] memory amounts,
        bytes memory data
    ) internal virtual override {
        super._beforeTokenTransfer(operator, from, to, ids, amounts, data);

        for (uint256 i = 0; i < ids.length; ++i) {
            uint256 id = ids[i];
            require(!paused(id), "KIP37Pausable: token transfer while paused");
        }
        require(!paused(), "KIP37Pausable: token transfer while paused");
    }
}
