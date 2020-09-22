pragma solidity ^0.5.6;

interface InterfaceIdentifier {
    /// @notice Query if a contract implements an interface
    /// @param interfaceID The interface identifier, as defined in KIP-13.
    /// @dev Interface identifier is defined in KIP-13. This function
    ///  uses less than 30,000 gas.
    /// @return `true` if the contract implements `interfaceID` and
    ///  `interfaceID` is not 0xffffffff, `false` otherwise.
    function supportsInterface(bytes4 interfaceID) external view returns (bool);
}