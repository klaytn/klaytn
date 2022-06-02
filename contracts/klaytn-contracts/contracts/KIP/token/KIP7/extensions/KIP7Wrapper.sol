// SPDX-License-Identifier: MIT
// Klaytn Contract Library v1.0.0 (KIP/token/KIP7/extensions/KIP7Wrapper.sol)
// Based on OpenZeppelin Contracts v4.5.0 (token/ERC20/extensions/ERC20Wrapper.sol)
// https://github.com/OpenZeppelin/openzeppelin-contracts/releases/tag/v4.5.0

pragma solidity ^0.8.0;

import "../KIP7.sol";
import "../utils/SafeKIP7.sol";

/**
 * @dev Extension of the KIP7 token contract to support token wrapping.
 *
 * Users can deposit and withdraw "underlying tokens" and receive a matching number of "wrapped tokens". This is useful
 * in conjunction with other modules. For example, combining this wrapping mechanism with {KIP7Votes} will allow the
 * wrapping of an existing "basic" KIP7 into a governance token.
 */
abstract contract KIP7Wrapper is KIP7 {
    IKIP7 public immutable underlying;

    constructor(IKIP7 underlyingToken) {
        underlying = underlyingToken;
    }

    /**
     * @dev See {KIP7-decimals}.
     */
    function decimals() public view virtual override returns (uint8) {
        try IKIP7Metadata(address(underlying)).decimals() returns (uint8 value) {
            return value;
        } catch {
            return super.decimals();
        }
    }

    /**
     * @dev Allow a user to deposit underlying tokens and mint the corresponding number of wrapped tokens.
     */
    function depositFor(address account, uint256 amount) public virtual returns (bool) {
        SafeKIP7.libSafeTransferFrom(underlying, _msgSender(), address(this), amount);
        _mint(account, amount);
        return true;
    }

    /**
     * @dev Allow a user to burn a number of wrapped tokens and withdraw the corresponding number of underlying tokens.
     */
    function withdrawTo(address account, uint256 amount) public virtual returns (bool) {
        _burn(_msgSender(), amount);
        SafeKIP7.libSafeTransfer(underlying, account, amount);
        return true;
    }

    /**
     * @dev Mint wrapped token to cover any underlyingTokens that would have been transferred by mistake. Internal
     * function that can be exposed with access control if desired.
     */
    function _recover(address account) internal virtual returns (uint256) {
        uint256 value = underlying.balanceOf(address(this)) - totalSupply();
        _mint(account, value);
        return value;
    }
}
