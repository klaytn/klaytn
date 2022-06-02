// SPDX-License-Identifier: MIT
// Klaytn Contract Library v1.0.0 (KIP/token/KIP7/utils/SafeKIP7.sol)
// Based on OpenZeppelin Contracts v4.5.0 (token/ERC720/utils/SafeERC20.sol)
// https://github.com/OpenZeppelin/openzeppelin-contracts/releases/tag/v4.5.0

pragma solidity ^0.8.0;

import "../IKIP7.sol";
import "../../../../utils/Address.sol";

/**
 * @title SafeKIP7
 * @dev Wrappers around KIP7 operations that throw on failure (when the token
 * contract returns false). Tokens that return no value (and instead revert or
 * throw on failure) are also supported, non-reverting calls are assumed to be
 * successful.
 * To use this library you can add a `using SafeKIP7 for IKIP7;` statement to your contract,
 * which allows you to call the safe operations as `token.safeTransfer(...)`, etc.
 *
 * IMPORTANT: {libSafeTransfer} and {libSafeTransferFrom} have been renamed from the original
 * {SafeERC20-safeTransfer} and {SafeERC20-safeTransferFrom} to avoid naming conflict with the existing KIP7
 * {IKIP7-safeTransfer} and {IKIP7-safeTransferFrom} which have functionality similar to the
 * {IERC1155/IERC721-safeTransferFrom} implmentations
 */
library SafeKIP7 {
    using Address for address;

    function libSafeTransfer(
        IKIP7 token,
        address to,
        uint256 value
    ) internal {
        _callOptionalReturn(token, abi.encodeWithSelector(token.transfer.selector, to, value));
    }

    function libSafeTransferFrom(
        IKIP7 token,
        address from,
        address to,
        uint256 value
    ) internal {
        _callOptionalReturn(token, abi.encodeWithSelector(token.transferFrom.selector, from, to, value));
    }

    /**
     * @dev Deprecated. This function has issues similar to the ones found in
     * {IKIP7-approve}, and its usage is discouraged.
     *
     * Whenever possible, use {safeIncreaseAllowance} and
     * {safeDecreaseAllowance} instead.
     */
    function safeApprove(
        IKIP7 token,
        address spender,
        uint256 value
    ) internal {
        // safeApprove should only be called when setting an initial allowance,
        // or when resetting it to zero. To increase and decrease it, use
        // 'safeIncreaseAllowance' and 'safeDecreaseAllowance'
        require(
            (value == 0) || (token.allowance(address(this), spender) == 0),
            "SafeKIP7: approve from non-zero to non-zero allowance"
        );
        _callOptionalReturn(token, abi.encodeWithSelector(token.approve.selector, spender, value));
    }

    function safeIncreaseAllowance(
        IKIP7 token,
        address spender,
        uint256 value
    ) internal {
        uint256 newAllowance = token.allowance(address(this), spender) + value;
        _callOptionalReturn(token, abi.encodeWithSelector(token.approve.selector, spender, newAllowance));
    }

    function safeDecreaseAllowance(
        IKIP7 token,
        address spender,
        uint256 value
    ) internal {
        unchecked {
            uint256 oldAllowance = token.allowance(address(this), spender);
            require(oldAllowance >= value, "SafeKIP7: decreased allowance below zero");
            uint256 newAllowance = oldAllowance - value;
            _callOptionalReturn(token, abi.encodeWithSelector(token.approve.selector, spender, newAllowance));
        }
    }

    /**
     * @dev Imitates a Solidity high-level call (i.e. a regular function call to a contract), relaxing the requirement
     * on the return value: the return value is optional (but if data is returned, it must not be false).
     * @param token The token targeted by the call.
     * @param data The call data (encoded using abi.encode or one of its variants).
     */
    function _callOptionalReturn(IKIP7 token, bytes memory data) private {
        // We need to perform a low level call here, to bypass Solidity's return data size checking mechanism, since
        // we're implementing it ourselves. We use {Address.functionCall} to perform this call, which verifies that
        // the target address contains contract code and also asserts for success in the low-level call.

        bytes memory returndata = address(token).functionCall(data, "SafeKIP7: low-level call failed");
        if (returndata.length > 0) {
            // Return data is optional
            require(abi.decode(returndata, (bool)), "SafeKIP7: KIP7 operation did not succeed");
        }
    }
}
