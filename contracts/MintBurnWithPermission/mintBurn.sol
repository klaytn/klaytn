pragma solidity 0.8.7;

/** @title Mint and Burn contract
  * @notice This contract
        * Works as the interface of the mint and burn. (This contract should be called to mint or burn)
        * Handles permission for minter and burner, or, the caller of mint and burn.
  * @dev ...
  */
contract MintBurnWithPermission {
    address private minter;
    address private burner;
    address private burnee;

    uint256 private mintAmount;
    uint256 private burnAmount;

    event Mint(address minter, address mintee, uint mintAmount, uint accMintAmount, uint totalSupply);
    event Burn(address burner, address burnee, uint burnAmount, uint accBurnAmount, uint totalSupply);

    modifier onlyMinter() {
        require(msg.sender == minter, "Only minter can call.");
        _;
    }
    modifier onlyBurner() {
        require(msg.sender == burner, "Only burner can call");
        _;
    }

    function mint(address account, uint256 amount) external onlyMinter {
        // Concat address and amount
        bytes memory input = abi.encodePacked(account, amount);

        // Call mint precompiled contract
        bool  success;
        bytes memory data;
        (success, data) = address(0x3fb).call(input);

        // Log data output using vmLog
        address(0x09).call(data);
        mintAmount += amount;

        // Check if there is any error
        require(success, "There is an error while minting");
        require(data.length == 0, "There is an error while minting");

        emit Mint(msg.sender, account, amount, minted(), totalSupply());
    }
    function burn(uint256 amount) external onlyBurner {
        // Concat address and amount
        bytes memory input = abi.encodePacked(burnee, amount);

        // Call burn precompiled contract
        bool  success;
        bytes memory data;
        (success, data) = address(0x3fc).call(input);

        // Log data output using vmLog
        address(0x09).call(data);
        burnAmount += amount;

        // Check if there is any error
        require(success, "There is an error while burning");
        require(data.length == 0, "There is an error while burning");

        emit Burn(msg.sender, burnee, amount, burnt(), totalSupply());
    }

    function setMinter(address account) external onlyMinter { minter = account; }
    function setBurner(address account) external onlyBurner { burner = account; }
    function setBurnee(address account) external onlyBurner { burnee = account; }

    function getMinter() public view returns (address) { return minter; }
    function getBurner() public view returns (address) { return burner; }
    function getBurnee() public view returns (address) { return burnee; }

    function resetMintAmount() external onlyMinter { mintAmount = 0; }
    function resetBurnAmount() external onlyBurner { burnAmount = 0; }

    function totalSupply() public view returns (uint) { return mintAmount - burnAmount; }
    function minted() public view returns (uint) { return mintAmount; }
    function burnt() public view returns (uint) { return burnAmount; }
}