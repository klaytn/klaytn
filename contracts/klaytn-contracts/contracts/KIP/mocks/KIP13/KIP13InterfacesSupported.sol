// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "../../utils/introspection/IKIP13.sol";

/**
 * https://eips.ethereum.org/EIPS/eip-214#specification
 * From the specification:
 * > Any attempts to make state-changing operations inside an execution instance with STATIC set to true will instead
 * throw an exception.
 * > These operations include [...], LOG0, LOG1, LOG2, [...]
 *
 * therefore, because this contract is staticcall'd we need to not emit events (which is how solidity-coverage works)
 * solidity-coverage ignores the /mocks folder, so we duplicate its implementation here to avoid instrumenting it
 */
contract SupportsInterfaceWithLookupMock is IKIP13 {
    /*
     * bytes4(keccak256('supportsInterface(bytes4)')) == 0x01ffc9a7
     */
    bytes4 public constant INTERFACE_ID_KIP13 = 0x01ffc9a7;

    /**
     * @dev A mapping of interface id to whether or not it's supported.
     */
    mapping(bytes4 => bool) private _supportedInterfaces;

    /**
     * @dev A contract implementing SupportsInterfaceWithLookup
     * implement KIP13 itself.
     */
    constructor() {
        _registerInterface(INTERFACE_ID_KIP13);
    }

    /**
     * @dev Implement supportsInterface(bytes4) using a lookup table.
     */
    function supportsInterface(bytes4 interfaceId) external view returns (bool) {
        return _supportedInterfaces[interfaceId];
    }

    /**
     * @dev Private method for registering an interface.
     */
    function _registerInterface(bytes4 interfaceId) internal {
        require(interfaceId != 0xffffffff, "KIP13InterfacesSupported: invalid interface id");
        _supportedInterfaces[interfaceId] = true;
    }
}

contract KIP13InterfacesSupported is SupportsInterfaceWithLookupMock {
    constructor(bytes4[] memory interfaceIds) {
        for (uint256 i = 0; i < interfaceIds.length; i++) {
            _registerInterface(interfaceIds[i]);
        }
    }
}
