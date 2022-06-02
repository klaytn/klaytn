// SPDX-License-Identifier: MIT
// Klaytn Contract Library v1.0.0 (KIP/token/KIP37/extensions/KIP37Mintable.sol)

pragma solidity ^0.8.0;

import "../KIP37.sol";
import "./KIP37URIStorage.sol";
import "./IKIP37Mintable.sol";
import "../../../../access/AccessControlEnumerable.sol";

/**
 * @dev Extension of KIP17 that supports permissioned token type creation and token minting
 * See http://kips.klaytn.com/KIPs/kip-37#minting-extension
 */
abstract contract KIP37Mintable is KIP37, KIP37URIStorage, IKIP37Mintable, AccessControlEnumerable {
    bytes32 public constant MINTER_ROLE = keccak256("KIP37_MINTER_ROLE");

    mapping(uint256 => address) public creators;

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
            interfaceId == type(IKIP37Mintable).interfaceId ||
            KIP37.supportsInterface(interfaceId) ||
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
    function uri(uint256 tokenId) public view virtual override(KIP37, KIP37URIStorage) returns (string memory) {
        return KIP37URIStorage.uri(tokenId);
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
        require(!_exists(id), "KIP37: token already created");

        creators[id] = _msgSender();

        _mint(_msgSender(), id, initialSupply, "");

        if (bytes(uri_).length > 0) {
            _tokenURIs[id] = uri_;
            emit URI(uri_, id);
        }
        return true;
    }

    /**
     * @dev Mints an `amount` of new `id` tokens and assigns `to` as owner
     *
     * Emits a {TransferSingle} event with 0X0 as the `from` account, for the `amount` tokens
     *
     * Requirements:
     *
     * - `id` must exist
     * - `to` must not be the zero address
     */
    function mint(
        uint256 id,
        address to,
        uint256 amount
    ) public virtual {
        require(_exists(id), "KIP37: nonexistent token");
        require(hasRole(MINTER_ROLE, _msgSender()), "KIP37: must have minter role to mint");

        _mint(to, id, amount, "");
    }

    /**
     * @dev For each item in `toList`, mints an `amount[i]` of new `id` tokens and assigns `toList[i]` as owner
     *
     * Emits multiple {TransferSingle} events with 0X0 as the `from` account, for the `amount` tokens
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
    ) public virtual {
        require(_exists(id), "KIP37: nonexistent token");
        require(hasRole(MINTER_ROLE, _msgSender()), "KIP37: must have minter role to mint");
        require(toList.length == amounts.length, "KIP37: toList and amounts length mismatch");
        for (uint256 i = 0; i < toList.length; ++i) {
            address to = toList[i];
            uint256 amount = amounts[i];
            _mint(to, id, amount, "");
        }
    }

    /**
     * @dev Mints multiple KIP37 token types `ids` in a batch and assigns the tokens according to the variables `to` and `amounts`.
     *
     *
     * Emits a {TransferBatch} event with 0X0 as the `from` account, for the `amount` tokens
     *
     * Requirements:
     *
     * - `to` must not be the zero address
     * - each`ids[i]` must exist
     * - `ids` and `amounts` must have the same number of elements
     */
    function mintBatch(
        address to,
        uint256[] memory ids,
        uint256[] memory amounts
    ) public virtual {
        for (uint256 i = 0; i < ids.length; ++i) {
            require(_exists(ids[i]), "KIP37: nonexistent token");
        }
        require(hasRole(MINTER_ROLE, _msgSender()), "KIP37: must have minter role to mint");
        _mintBatch(to, ids, amounts, "");
    }

    /**
     * @dev Internal function for checking if an `id` token type has already been created
     */
    function _exists(uint256 id) internal view returns (bool) {
        address creator = creators[id];
        return creator != address(0);
    }
}
