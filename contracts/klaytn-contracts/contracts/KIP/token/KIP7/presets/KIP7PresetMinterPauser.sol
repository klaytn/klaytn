// SPDX-License-Identifier: MIT
// Klaytn Contract Library v1.0.0 (KIP/token/KIP7/presets/KIP7PresetMinterPauser.sol)
// Based on OpenZeppelin Contracts v4.5.0 (token/ERC20/presets/ERC20PresetMinterPauser.sol)
// https://github.com/OpenZeppelin/openzeppelin-contracts/releases/tag/v4.5.0

pragma solidity ^0.8.0;

import "../../../../utils/Context.sol";
import "../../../utils/introspection/KIP13.sol";
import "../../../../access/AccessControlEnumerable.sol";
import "../extensions/KIP7Mintable.sol";
import "../extensions/KIP7Burnable.sol";
import "../extensions/KIP7Pausable.sol";

/**
 * @dev {KIP7} token, including:
 *
 *  - ability for holders to burn (destroy) their tokens
 *  - a minter role that allows for token minting (creation)
 *  - a pauser role that allows to stop all token transfers
 *
 * This contract uses {AccessControl} to lock permissioned functions using the
 * different roles - head to its documentation for details.
 *
 * The account that deploys the contract will be granted the minter and pauser
 * roles, as well as the default admin role, which will let it grant both minter
 * and pauser roles to other accounts.
 */
contract KIP7PresetMinterPauser is Context, KIP13, AccessControlEnumerable, KIP7Mintable, KIP7Burnable, KIP7Pausable {
    /**
     * @dev Grants `DEFAULT_ADMIN_ROLE`, `MINTER_ROLE` and `PAUSER_ROLE` to the
     * account that deploys the contract.
     *
     * See {KIP7-constructor}.
     */
    constructor(string memory name, string memory symbol) KIP7(name, symbol) {
        _setupRole(DEFAULT_ADMIN_ROLE, _msgSender());
        _setupRole(MINTER_ROLE, _msgSender());
        _setupRole(PAUSER_ROLE, _msgSender());
    }

    /**
     * @dev Creates `amount` new tokens for `to`.
     *
     * See {KIP7Mintable}.
     *
     * Requirements:
     *
     * - the caller must have the `MINTER_ROLE`.
     */
    function mint(address to, uint256 amount) public override returns (bool) {
        return super.mint(to, amount);
    }

    /**
     * @dev Pauses all token transfers.
     *
     * See {KIP7Pausable}.
     *
     * Requirements:
     *
     * - the caller must have the `PAUSER_ROLE`.
     */
    function pause() public virtual override {
        super.pause();
    }

    /**
     * @dev Unpauses all token transfers.
     *
     * See {KIP7Pausable}.
     *
     * Requirements:
     *
     * - the caller must have the `PAUSER_ROLE`.
     */
    function unpause() public virtual override {
        super.unpause();
    }

    /**
     * @dev Returns true if `interfaceId` is implemented and false otherwise
     *
     * See {IKIP13} and {IERC165}.
     */
    function supportsInterface(bytes4 interfaceId) public view override(KIP13, AccessControlEnumerable) returns (bool) {
        return
            interfaceId == type(IKIP7).interfaceId ||
            interfaceId == type(IKIP7Burnable).interfaceId ||
            interfaceId == type(IKIP7Mintable).interfaceId ||
            interfaceId == type(IKIP7Pausable).interfaceId ||
            AccessControlEnumerable.supportsInterface(interfaceId) ||
            KIP13.supportsInterface(interfaceId);
    }

    function _beforeTokenTransfer(
        address from,
        address to,
        uint256 amount
    ) internal virtual override(KIP7, KIP7Pausable) {
        super._beforeTokenTransfer(from, to, amount);
    }
}
