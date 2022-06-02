// SPDX-License-Identifier: MIT
// Klaytn Contract Library v1.0.0 (KIP/token/KIP37/extensions/KIP37Burnable.sol)
// Based on OpenZeppelin Contracts v4.5.0 (token/ERC1155/extensions/ERC1155Burnable.sol)
// https://github.com/OpenZeppelin/openzeppelin-contracts/releases/tag/v4.5.0

pragma solidity ^0.8.0;

import "../KIP37.sol";
import "./IKIP37Burnable.sol";

/**
 * @dev Extension of {KIP37} that allows token holders to destroy both their
 * own tokens and those that they have been approved to use.
 */
abstract contract KIP37Burnable is KIP37, IKIP37Burnable {
    /**
     * @dev Returns true if `interfaceId` is implemented and false otherwise
     *
     * See {IKIP13}.
     */
    function supportsInterface(bytes4 interfaceId) public view virtual override returns (bool) {
        return interfaceId == type(IKIP37Burnable).interfaceId || super.supportsInterface(interfaceId);
    }

    /**
     * @dev Destroys `amount` tokens of token type `id` from `from`
     * See {KIP37-_burn}
     * Requirements:
     *
     * - `from` cannot be the zero address.
     * - `from` must have at least `amount` tokens of token type `id`.
     *
     * Emits a {TransferSingle} event with the 0x0 address as `to`.
     */
    function burn(
        address account,
        uint256 id,
        uint256 amount
    ) public virtual {
        require(
            account == _msgSender() || isApprovedForAll(account, _msgSender()),
            "KIP37: caller is not owner nor approved"
        );

        _burn(account, id, amount);
    }

    /**
     * @dev Destroys each amount `amounts[i]` tokens of each token type `ids[i]` from `from`
     * See {KIP37-_burnBatch}
     *
     * Requirements:
     *
     * - `from` cannot be the zero address.
     * - `from` must have at least `amounts[i]` tokens of token type `ids[i]`.
     * - `ids` and `amounts` must have the same length.
     *
     * Emits a {TransferBatch} event with the 0x0 address as `to`.
     */
    function burnBatch(
        address account,
        uint256[] memory ids,
        uint256[] memory amounts
    ) public virtual {
        require(
            account == _msgSender() || isApprovedForAll(account, _msgSender()),
            "KIP37: caller is not owner nor approved"
        );

        _burnBatch(account, ids, amounts);
    }
}
