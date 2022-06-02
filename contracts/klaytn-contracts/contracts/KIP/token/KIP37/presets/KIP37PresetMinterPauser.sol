// SPDX-License-Identifier: MIT
// Klaytn Contract Library v1.0.0 (KIP/token/KIP37/presets/KIP3PresetMinterPauser.sol)
// Based on OpenZeppelin Contracts v4.5.0 (token/ERC1155/presets/ERC1155PresetMinterPauser.sol)
// https://github.com/OpenZeppelin/openzeppelin-contracts/releases/tag/v4.5.0

pragma solidity ^0.8.0;

import "../../../../utils/Context.sol";
import "../../../../access/AccessControlEnumerable.sol";
import "../extensions/KIP37Mintable.sol";
import "../extensions/KIP37Burnable.sol";
import "../extensions/KIP37Pausable.sol";

/**
 * @dev {KIP37} token, including:
 *
 *  - ability for holders to burn (destroy) their tokens
 *  - a minter role that allows for token minting (and creation) in batches, singleton token types, or multi-token types
 *  - a pauser role that allows to stop all token transfers and token type transfers
 *
 * This contract uses {AccessControl} to lock permissioned functions using the
 * different roles - head to its documentation for details.
 *
 * The account that deploys the contract will be granted the minter and pauser
 * roles, as well as the default admin role, which will let it grant both minter
 * and pauser roles to other accounts.
 */
contract KIP37PresetMinterPauser is Context, AccessControlEnumerable, KIP37Mintable, KIP37Burnable, KIP37Pausable {
    /**
     * @dev Grants `DEFAULT_ADMIN_ROLE`, `MINTER_ROLE`, and `PAUSER_ROLE` to the account that
     * deploys the contract.
     */
    constructor(string memory uri_) KIP37(uri_) {
        _setupRole(DEFAULT_ADMIN_ROLE, _msgSender());
        _setupRole(MINTER_ROLE, _msgSender());
        _setupRole(PAUSER_ROLE, _msgSender());
    }

    /**
     * @dev Returns true if `interfaceId` is implemented and false otherwise
     *
     * See {IKIP13} and {IERC165}.
     */
    function supportsInterface(bytes4 interfaceId)
        public
        view
        virtual
        override(KIP37Burnable, KIP37Mintable, KIP37Pausable, AccessControlEnumerable)
        returns (bool)
    {
        return
            KIP37Burnable.supportsInterface(interfaceId) ||
            KIP37Mintable.supportsInterface(interfaceId) ||
            KIP37Pausable.supportsInterface(interfaceId) ||
            AccessControlEnumerable.supportsInterface(interfaceId);
    }

    /**
     * @dev See {IKIP37MetadataURI-uri}.
     *
     * This implementation returns the concatenation of the `_baseURI`
     * and the token-specific uri if the latter is set
     *
     * This enables the following behaviors:
     *
     * - if `_tokenURIs[tokenId]` is set, then the result is the concatenation
     *   of `_baseURI` and `_tokenURIs[tokenId]` (keep in mind that `_baseURI`
     *   is empty per default);
     *
     * - if `_tokenURIs[tokenId]` is NOT set then we fallback to `super.uri()`
     *   which in most cases will contain `KIP37._uri`;
     *
     * - if `_tokenURIs[tokenId]` is NOT set, and if the parents do not have a
     *   uri value set, then the result is empty.
     */
    function uri(uint256 tokenId) public view virtual override(KIP37, KIP37Mintable) returns (string memory) {
        return KIP37Mintable.uri(tokenId);
    }

    /**
     * @dev Creates a new `id` token type and assigns the caller as owner of `initialSupply` while
     * setting a `uri` for this token type
     *
     * Requirements:
     *
     * - `id` must not already exist
     * - for this implementation, caller must have the role MINTER_ROLE
     *
     * Emits a {TransferSingle} event with 0X0 as the `from` account, for the `intialSupply` tokens
     *
     * If non zero length `uri_` is submitted, emits a {URI} event
     */
    function create(
        uint256 id,
        uint256 initialSupply,
        string memory uri_
    ) public virtual override onlyRole(MINTER_ROLE) returns (bool) {
        return super.create(id, initialSupply, uri_);
    }

    /**
     * @dev Creates `amount` new tokens for `to`, of token type `id`.
     *
     * See {KIP37Mintable-mint(uint256 id, address to, uint256 amount)}
     *
     * Requirements:
     *
     * - the caller must have the `MINTER_ROLE`.
     */
    function mint(
        uint256 id,
        address to,
        uint256 amount
    ) public virtual override {
        super.mint(id, to, amount);
    }

    /**
     * @dev For each item in `toList`, creates an `amount[i]` of new `id` tokens and assigns `toList[i]` as owner
     *
     * See {KIP37Mintable-mint(uint256 id, address[] memory toList, uint256[] memory amounts)}
     *
     * Requirements:
     *
     * - `id` must exist
     * - each `toList[i]` must not be the zero address
     * - `toList` and `amounts` must have the same number of elements
     */
    function mint(
        uint256 id,
        address[] memory toList,
        uint256[] memory amounts
    ) public virtual override {
        super.mint(id, toList, amounts);
    }

    /**
     * @dev xref:ROOT:kip37.adoc#batch-operations[Batched] variant of {mint}.
     */
    function mintBatch(
        address to,
        uint256[] memory ids,
        uint256[] memory amounts
    ) public virtual override {
        super.mintBatch(to, ids, amounts);
    }

    /**
     * @dev Pauses all token transfers.
     *
     * See {KIP37Pausable-pause()} and {Pausable-_pause}.
     *
     * Requirements:
     *
     * - the caller must have the `PAUSER_ROLE`.
     */
    function pause() public virtual override {
        require(hasRole(PAUSER_ROLE, _msgSender()), "KIP37PresetMinterPauser: must have pauser role to pause");
        super.pause();
    }

    /**
     * @dev Pauses all token transfers of token type `id`.
     *
     * See {KIP37Pausable-pause(uint256 id)} and {Pausable-_pause}.
     *
     * Requirements:
     *
     * - the caller must have the `PAUSER_ROLE`.
     */
    function pause(uint256 id) public virtual override {
        require(hasRole(PAUSER_ROLE, _msgSender()), "KIP37PresetMinterPauser: must have pauser role to pause");
        super.pause(id);
    }

    /**
     * @dev Unpauses all token transfers.
     *
     * See {KIP37Pausable-unpause())} and {Pausable-_unpause}.
     *
     * Requirements:
     *
     * - the caller must have the `PAUSER_ROLE`.
     */
    function unpause() public virtual override {
        require(hasRole(PAUSER_ROLE, _msgSender()), "KIP37PresetMinterPauser: must have pauser role to unpause");
        super.unpause();
    }

    /**
     * @dev Unpauses all token transfers of token type `id`.
     *
     * See {KIP37Pausable-unpause(uint256 id))} and {Pausable-_unpause}.
     *
     * Requirements:
     *
     * - the caller must have the `PAUSER_ROLE`.
     */
    function unpause(uint256 id) public virtual override {
        require(hasRole(PAUSER_ROLE, _msgSender()), "KIP37PresetMinterPauser: must have pauser role to unpause");
        super.unpause(id);
    }

    function _beforeTokenTransfer(
        address operator,
        address from,
        address to,
        uint256[] memory ids,
        uint256[] memory amounts,
        bytes memory data
    ) internal virtual override(KIP37, KIP37Pausable) {
        super._beforeTokenTransfer(operator, from, to, ids, amounts, data);
    }
}
