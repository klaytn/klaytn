// SPDX-License-Identifier: MIT
// Klaytn Contract Library v1.0.0 (KIP/token/KIP17/extensions/KIP17Burnable.sol)
// Based on OpenZeppelin Contracts v4.5.0 (token/ERC721/extensions/ERC721Burnable.sol)
// https://github.com/OpenZeppelin/openzeppelin-contracts/releases/tag/v4.5.0

pragma solidity ^0.8.0;

import "../../../../utils/Context.sol";
import "../KIP17.sol";
import "./IKIP17Burnable.sol";

/**
 * @title KIP17 Burnable Token
 * @dev KIP17 Token that can be burned (destroyed).
 */
abstract contract KIP17Burnable is Context, KIP17, IKIP17Burnable {
    /**
     * @dev See {IKIP13-supportsInterface}.
     */
    function supportsInterface(bytes4 interfaceId) public view virtual override returns (bool) {
        return interfaceId == type(IKIP17Burnable).interfaceId || KIP17.supportsInterface(interfaceId);
    }

    /**
     * @dev Burns `tokenId`. See {KIP17-_burn}.
     *
     * Requirements:
     *
     * - The caller must be the current owner, an authorized operator, or the approved address for `tokenId`
     */
    function burn(uint256 tokenId) public virtual {
        //solhint-disable-next-line max-line-length
        require(_isApprovedOrOwner(_msgSender(), tokenId), "ERC721Burnable: caller is not owner nor approved");
        _burn(tokenId);
    }
}
