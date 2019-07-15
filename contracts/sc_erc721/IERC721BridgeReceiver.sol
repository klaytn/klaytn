pragma solidity ^0.4.24;

contract IERC721BridgeReceiver {
    function onERC721Received(address _from, uint256 _tokenId, address _to) public;
}
