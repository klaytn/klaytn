pragma solidity ^0.4.24;

import "../externals/openzeppelin-solidity/contracts/token/ERC20/ERC20.sol";
import "../externals/openzeppelin-solidity/contracts/utils/Address.sol";
import "../externals/openzeppelin-solidity/contracts/ownership/Ownable.sol";
import "./IERC20BridgeReceiver.sol";

/**
 * @title ERC20ServiceChain
 * @dev ERC20 service chain value transfer logic for 1-step transfer.
 */
contract ERC20ServiceChain is ERC20, Ownable {
    using Address for address;
    address public bridge;

    constructor(address _bridge) internal {
        if (!_bridge.isContract()) {
            revert("bridge is not a contract");
        }

        bridge = _bridge;
    }

    function setBridge(address _bridge) public onlyOwner {
        bridge = _bridge;
    }

    function requestValueTransfer(uint256 _amount, address _to) external {
        transfer(bridge, _amount);
        IERC20BridgeReceiver(bridge).onERC20Received(msg.sender, _amount, _to);
    }
}
