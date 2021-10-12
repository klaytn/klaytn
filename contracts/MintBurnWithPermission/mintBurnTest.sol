pragma solidity 0.8.7;

/** @title Mint and Burn contract
  * @notice This contract
        * Works as the interface of the mint and burn. (This contract should be called to mint or burn)
        * Handles permission for minter and burner, or, the caller of mint and burn.
  * @dev ...
  */
contract MintBurnTest {
    function mint(address account, uint256 amount) external {
        // Concat address and amount
        bytes memory input = abi.encodePacked(account, amount);

        // Call mint precompiled contract
        bool  success;
        bytes memory data;
        (success, data) = address(0x3fb).call(input);

        // Log data output using vmLog
        address(0x09).call(data);

        // Check if there is any error
        require(success, "There is an error while minting");
        require(data.length == 0, "There is an error while minting");
    }

    function burn(address account, uint256 amount) external {
        // Concat address and amount
        bytes memory input = abi.encodePacked(account, amount);

        // Call burn precompiled contract
        bool  success;
        bytes memory data;
        (success, data) = address(0x3fc).call(input);

        // Log data output using vmLog
        address(0x09).call(data);

        // Check if there is any error
        require(success, "There is an error while burning");
        require(data.length == 0, "There is an error while burning");
    }
}