pragma solidity ^0.4.24;

import "./BridgeTransferKLAY.sol";
import "./BridgeTransferERC20.sol";
import "./BridgeTransferERC721.sol";
import "./BridgeCounterPart.sol";


contract Bridge is BridgeCounterPart, BridgeTransferKLAY, BridgeTransferERC20, BridgeTransferERC721 {
    uint64 public constant VERSION = 1;

    constructor(bool _modeMintBurn) BridgeTransfer(_modeMintBurn) public payable {
    }
}
