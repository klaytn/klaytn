pragma solidity ^0.4.24;

import "../externals/openzeppelin-solidity/contracts/token/ERC721/ERC721.sol";
import "../externals/openzeppelin-solidity/contracts/ownership/Ownable.sol";
import "./INFTReceiver.sol";

/**
 * @title ERC721ServiceChain
 * @dev ERC721 service chain value transfer logic for 1-step transfer.
 */
contract ERC721ServiceChain is ERC721, Ownable {
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

    bytes4 private constant _ERC721_RECEIVED = 0x150b7a02;

    // user request value transfer to main / service chain.
    function requestValueTransfer(uint256 _uid, address _to) external {
        transferFrom(msg.sender, bridge, _uid);

        bytes4 retval = INFTReceiver(bridge).onNFTReceived(msg.sender, _uid, _to);
        require(retval == _ERC721_RECEIVED, "Sent to a bridge which is not an ERC721 receiver" );
    }
}
