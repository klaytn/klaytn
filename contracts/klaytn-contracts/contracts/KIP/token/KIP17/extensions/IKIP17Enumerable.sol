// SPDX-License-Identifier: MIT
// Klaytn Contract Library v1.0.0 (KIP/token/KIP17/extensions/IKIP17Enumerable.sol)
// Based on OpenZeppelin Contracts v4.5.0 (token/ERC721/extensions/IERC721Enumerable.sol)
// https://github.com/OpenZeppelin/openzeppelin-contracts/releases/tag/v4.5.0

pragma solidity ^0.8.0;

import "../IKIP17.sol";

/**
 * @title KIP-17 Non-Fungible Token Standard, optional enumeration extension
 * @dev See https://kips.klaytn.com/KIPs/kip-17#enumeration-extension
 */
interface IKIP17Enumerable is IKIP17 {
    /**
     * @dev Returns the total amount of tokens stored by the contract.
     */
    function totalSupply() external view returns (uint256);

    /**
     * @dev Returns a token ID at a given `index` of all the tokens stored by the contract.
     * Use along with {totalSupply} to enumerate all tokens.
     */
    function tokenByIndex(uint256 index) external view returns (uint256);

    /**
     * @dev Returns a token ID owned by `owner` at a given `index` of its token list.
     * Use along with {balanceOf} to enumerate all of ``owner``'s tokens.
     */
    function tokenOfOwnerByIndex(address owner, uint256 index) external view returns (uint256);
}
