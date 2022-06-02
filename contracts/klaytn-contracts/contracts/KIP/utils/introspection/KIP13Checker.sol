// SPDX-License-Identifier: MIT
// Klaytn Contract Library v1.0.0 (KIP/utils/introspection/KIP13Checker.sol)
// Based on OpenZeppelin Contracts v4.5.0 (utils/introspection/ERC165Checker.sol)
// https://github.com/OpenZeppelin/openzeppelin-contracts/releases/tag/v4.5.0

pragma solidity ^0.8.0;

import "../../interfaces/IKIP13.sol";

/**
 * @dev Library used to query support of an interface declared via {IKIP13}.
 *
 * Note that these functions return the actual result of the query: they do not
 * `revert` if an interface is not supported. It is up to the caller to decide
 * what to do in these cases.
 */
library KIP13Checker {
    // As per the KIP-13 spec, no interface should ever match 0xffffffff
    bytes4 private constant _INTERFACE_ID_INVALID = 0xffffffff;

    /**
     * @dev Returns true if `account` supports the {IKIP13} interface,
     */
    function supportsKIP13(address account) internal view returns (bool) {
        // Any contract that implements KIP13 must explicitly indicate support of
        // InterfaceId_KIP13 and explicitly indicate non-support of InterfaceId_Invalid
        return
            _supportsKIP13Interface(account, type(IKIP13).interfaceId) &&
            !_supportsKIP13Interface(account, _INTERFACE_ID_INVALID);
    }

    /**
     * @dev Returns true if `account` supports the interface defined by
     * `interfaceId`. Support for {IKIP13} itself is queried automatically.
     *
     * See {IKIP13-supportsInterface}.
     */
    function supportsInterface(address account, bytes4 interfaceId) internal view returns (bool) {
        // query support of both KIP13 as per the spec and support of interfaceId
        return supportsKIP13(account) && _supportsKIP13Interface(account, interfaceId);
    }

    /**
     * @dev Returns a boolean array where each value corresponds to the
     * interfaces passed in and whether they're supported or not. This allows
     * you to batch check interfaces for a contract where your expectation
     * is that some interfaces may not be supported.
     *
     * See {IKIP13-supportsInterface}.
     *
     */
    function getSupportedInterfaces(address account, bytes4[] memory interfaceIds)
        internal
        view
        returns (bool[] memory)
    {
        // an array of booleans corresponding to interfaceIds and whether they're supported or not
        bool[] memory interfaceIdsSupported = new bool[](interfaceIds.length);

        // query support of KIP13 itself
        if (supportsKIP13(account)) {
            // query support of each interface in interfaceIds
            for (uint256 i = 0; i < interfaceIds.length; i++) {
                interfaceIdsSupported[i] = _supportsKIP13Interface(account, interfaceIds[i]);
            }
        }

        return interfaceIdsSupported;
    }

    /**
     * @dev Returns true if `account` supports all the interfaces defined in
     * `interfaceIds`. Support for {IKIP13} itself is queried automatically.
     *
     * Batch-querying can lead to gas savings by skipping repeated checks for
     * {IKIP13} support.
     *
     * See {IKIP13-supportsInterface}.
     */
    function supportsAllInterfaces(address account, bytes4[] memory interfaceIds) internal view returns (bool) {
        // query support of KIP13 itself
        if (!supportsKIP13(account)) {
            return false;
        }

        // query support of each interface in _interfaceIds
        for (uint256 i = 0; i < interfaceIds.length; i++) {
            if (!_supportsKIP13Interface(account, interfaceIds[i])) {
                return false;
            }
        }

        // all interfaces supported
        return true;
    }

    /**
     * @notice Query if a contract implements an interface, does not check KIP13 support
     * @param account The address of the contract to query for support of an interface
     * @param interfaceId The interface identifier, as specified in KIP13
     * @return true if the contract at account indicates support of the interface with
     * identifier interfaceId, false otherwise
     * @dev Assumes that account contains a contract that supports KIP13, otherwise
     * the behavior of this method is undefined. This precondition can be checked
     * with {supportsKIP13}.
     * Interface identification is specified in KIP13.
     */
    function _supportsKIP13Interface(address account, bytes4 interfaceId) private view returns (bool) {
        bytes memory encodedParams = abi.encodeWithSelector(IKIP13.supportsInterface.selector, interfaceId);
        (bool success, bytes memory result) = account.staticcall{gas: 30000}(encodedParams);
        if (result.length < 32) return false;
        return success && abi.decode(result, (bool));
    }
}
