// Copyright 2019 The klaytn Authors
// This file is part of the klaytn library.
//
// The klaytn library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The klaytn library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the klaytn library. If not, see <http://www.gnu.org/licenses/>.

pragma solidity 0.5.6;

import "./BridgeTransferKLAY.sol";
import "./BridgeTransferERC20.sol";
import "./BridgeTransferERC721.sol";
import "./BridgeCounterPart.sol";


contract Bridge is BridgeCounterPart, BridgeTransferKLAY, BridgeTransferERC20, BridgeTransferERC721 {
    uint64 public constant VERSION = 1;

    constructor(bool _modeMintBurn) BridgeTransfer(_modeMintBurn) public payable {
    }
}
