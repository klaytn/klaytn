// SPDX-License-Identifier: MIT
// Klaytn Contract Library v1.0.0 (KIP/token/KIP17/extensions/IKIP17Mintable.sol)

pragma solidity ^0.8.0;

/**
 * @dev Minting extension of the KIP17 standard as defined in the KIP.
 * See https://kips.klaytn.com/KIPs/kip-17#minting-extension
 */
interface IKIP17Mintable {
    /**
     * @dev Creates `tokenId` token and assigns them to `account`
     * increasing the total supply.
     *
     * Requirements:
     *
     * - caller must have the {KIP7Mintable-MINTER_ROLE}
     *
     * Emits a {Transfer} event with 0X0 as the `from` account
     */
    function mint(address to, uint256 tokenId) external returns (bool);

    /**
     * @dev Check if `account` has the assigned Minter role via {AccessControl-hasRole}
     */
    function isMinter(address account) external view returns (bool);

    /**
     * @dev Assign the Minter role to `account` via {AccessControl-grantRole}
     *
     * Emits a {RoleGranted} event
     *
     * Requirements:
     *
     * - caller must have the {KIP7Mintable-MINTER_ROLE}
     */
    function addMinter(address account) external;

    /**
     * @dev Renounce the Minter role of the caller via {AccessControl-renounceRole}
     *
     * Emits a {RoleRevoked} event
     */
    function renounceMinter() external;
}
