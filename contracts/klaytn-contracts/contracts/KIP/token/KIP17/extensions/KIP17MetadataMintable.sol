// SPDX-License-Identifier: MIT
// Klaytn Contract Library v1.0.0 (KIP/token/KIP17/extensions/KIP17MetadataMintable.sol)

pragma solidity ^0.8.0;

import "./KIP17URIStorage.sol";
import "../../../../access/AccessControlEnumerable.sol";
import "./IKIP17MetadataMintable.sol";

/**
 * @dev Extension of KIP17 that supports permissioned token minting with linked URI
 * See https://kips.klaytn.com/KIPs/kip-17#minting-with-uri-extension
 */
abstract contract KIP17MetadataMintable is KIP17URIStorage, AccessControlEnumerable, IKIP17MetadataMintable {
    bytes32 public constant MINTER_ROLE = keccak256("KIP17_MINTER_ROLE");

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
            interfaceId == type(IKIP17MetadataMintable).interfaceId ||
            KIP17.supportsInterface(interfaceId) ||
            AccessControlEnumerable.supportsInterface(interfaceId);
    }

    /**
     * @dev See {IKIP17metadataMintable-mintWithTokenURI}
     *
     * IMPORTANT: this uses _safeMint internally, please be aware that if you do not want this safety functionality, replace with _mint
     *
     * Emits a {Transfer} event with 0X0 as the `from` account
     */
    function mintWithTokenURI(
        address to,
        uint256 tokenId,
        string memory _tokenURI
    ) public virtual override onlyRole(MINTER_ROLE) returns (bool) {
        _safeMint(to, tokenId);
        _setTokenURI(tokenId, _tokenURI);
        return true;
    }

    /**
     * @dev See {IKIP17Mintable-isMinter}
     */
    function isMinter(address account) public view returns (bool) {
        return hasRole(MINTER_ROLE, account);
    }

    /**
     * @dev See {IKIP17Mintable-addMinter}
     *
     * Emits a {RoleGranted} event
     */
    function addMinter(address account) public onlyRole(DEFAULT_ADMIN_ROLE) {
        grantRole(MINTER_ROLE, account);
    }

    /**
     * @dev See {IKIP17Mintable-renounceMinter}
     *
     * Emits a {RoleRevoked} event
     */
    function renounceMinter() public {
        renounceRole(MINTER_ROLE, _msgSender());
    }
}
