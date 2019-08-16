// Copyright 2018 The klaytn Authors
//
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

pragma solidity ^0.4.24;

contract KlaytnReward {

    uint public totalAmount;
    mapping(address => uint256) public balanceOf;

    function KlaytnReward() public {
    }

    function () payable public {
        uint amount = msg.value;
        balanceOf[msg.sender] += amount;
        totalAmount += amount;
    }

    function reward(address receiver) payable public {
        uint amount = msg.value;
        balanceOf[receiver] += amount;
        totalAmount += amount;
    }

    function safeWithdrawal() public {
        uint amount = balanceOf[msg.sender];
        balanceOf[msg.sender] = 0;
        if (amount > 0) {
             if (msg.sender.send(amount)) {

             } else {
                balanceOf[msg.sender] = amount;
             }
        }
    }
}
