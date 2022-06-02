// SPDX-License-Identifier: MIT
// Klaytn Contract Library v1.0.0 (KIP/token/KIP17/extensions/KIP17Mintable.sol)

pragma solidity ^0.8.0;

import "../KIP17.sol";
import "../../../../access/AccessControlEnumerable.sol";
import "./IKIP17Mintable.sol";

/**
 * @dev Extension of KIP17 that supports permissioned token minting
 * See https://kips.klaytn.com/KIPs/kip-17#minting-extension
 */
abstract contract KIP17Mintable is KIP17, AccessControlEnumerable, IKIP17Mintable {
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
            interfaceId == type(IKIP17Mintable).interfaceId ||
            KIP17.supportsInterface(interfaceId) ||
            AccessControlEnumerable.supportsInterface(interfaceId);
    }

    /**
     * @dev See {IKIP17Mintable-mint}
     *
     * IMPORTANT: this uses _safeMint internally, please be aware that if you do not want this safety functionality, replace with _mint
     *
     * Emits a {Transfer} event with 0X0 as the `from` account
     */
    function mint(address to, uint256 tokenId) public virtual override onlyRole(MINTER_ROLE) returns (bool) {
        _safeMint(to, tokenId);
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
