// SPDX-License-Identifier: MIT
// Klaytn Contract Library v1.0.0 (KIP/token/KIP7/extensions/KIP7Burnable.sol)
// Based on OpenZeppelin Contracts v4.5.0 (token/ERC20/extensions/ERC20Burnable.sol)
// https://github.com/OpenZeppelin/openzeppelin-contracts/releases/tag/v4.5.0

pragma solidity ^0.8.0;

import "../../../../utils/Context.sol";
import "../KIP7.sol";
import "./IKIP7Burnable.sol";

/**
 * @dev Extension of KIP7 that allows token holders to destroy both their own
 * tokens and those that they have an allowance for, in a way that can be
 * recognized off-chain (via event analysis).
 * See https://kips.klaytn.com/KIPs/kip-7#burning-extension
 */
abstract contract KIP7Burnable is Context, KIP7, IKIP7Burnable {
    /**
     * @dev Destroys `amount` tokens from the caller.
     *
     * See {KIP7-_burn}.
     *
     * Requirements:
     *
     * - caller's balance must be greater than or equal to `_amount`
     */
    function burn(uint256 amount) public virtual {
        _burn(_msgSender(), amount);
    }

    /**
     * @dev Destroys `amount` tokens from `account`, deducting from the caller's
     * allowance.
     *
     * See {KIP7-_burn} and {KIP7-allowance}.
     *
     * Requirements:
     *
     * - the caller must have allowance for ``accounts``'s tokens of at least
     * `amount`.
     */
    function burnFrom(address account, uint256 amount) public virtual {
        _spendAllowance(account, _msgSender(), amount);
        _burn(account, amount);
    }
}
