pragma solidity ^0.4.24;

import "../externals/openzeppelin-solidity/contracts/token/ERC721/ERC721.sol";
import "../externals/openzeppelin-solidity/contracts/token/ERC721/ERC721Metadata.sol";

import "../externals/openzeppelin-solidity/contracts/ownership/Ownable.sol";
import "./IERC721BridgeReceiver.sol";

/**
 * @title ERC721ServiceChain
 * @dev ERC721 service chain value transfer logic for 1-step transfer.
 */
contract ERC721ServiceChain is ERC721, ERC721Metadata, Ownable {
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

    function requestValueTransfer(uint256 _uid, address _to) external {
        transferFrom(msg.sender, bridge, _uid);

        IERC721BridgeReceiver(bridge).onERC721Received(msg.sender, _uid, _to);
    }
}
