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

pragma solidity ^0.5.6;

import "../externals/openzeppelin-solidity/contracts/token/ERC20/IERC20.sol";
import "../externals/openzeppelin-solidity/contracts/token/ERC20/ERC20Mintable.sol";
import "../externals/openzeppelin-solidity/contracts/math/SafeMath.sol";
import "../externals/openzeppelin-solidity/contracts/ownership/Ownable.sol";


contract BridgeOperator is Ownable {
    mapping(address => bool) public operators;
    address[] public operatorList;
    mapping(bytes32 => mapping(address => bool)) public votes; // <sha3(type, args, nonce), <operator, vote>>
    mapping(bytes32 => uint8) public votesCounts; // <sha3(type, args, nonce)>
    mapping(uint64 => bool) public closedValueTransferVotes; // nonce
    mapping(uint8 => uint8) public operatorThresholds; // <vote type>
    uint64 public configurationNonce;

    enum VoteType {
        ValueTransfer,
        Configuration,
        Max
    }

    constructor() internal {
        for (uint8 i = 0; i < uint8(VoteType.Max); i++) {
            operatorThresholds[uint8(i)] = 1;
        }

        operators[msg.sender] = true;
        operatorList.push(msg.sender);
    }

    modifier onlyOperators()
    {
        require(operators[msg.sender], "msg.sender is not an operator");
        _;
    }

    function getOperatorList() external view returns(address[] memory) {
        return operatorList;
    }

    // voteCommon handles common functionality for voting.
    function voteCommon(VoteType voteType, bytes32 _voteKey)
        internal
        returns(bool)
    {
        if (!votes[_voteKey][msg.sender]) {
            votes[_voteKey][msg.sender] = true;
            require(votesCounts[_voteKey] < votesCounts[_voteKey] + 1, "votesCounts overflow");
            votesCounts[_voteKey]++;
        }
        if (votesCounts[_voteKey] >= operatorThresholds[uint8(voteType)]) {
            return true;
        }
        return false;
    }

    // voteValueTransfer votes value transfer transaction with the operator.
    function voteValueTransfer(uint64 _requestNonce)
        internal
        returns(bool)
    {
        require(!closedValueTransferVotes[_requestNonce], "closed vote");

        if (voteCommon(VoteType.ValueTransfer, keccak256(msg.data))) {
            closedValueTransferVotes[_requestNonce] = true;
            return true;
        }

        return false;
    }

    // voteConfiguration votes contract configuration transaction with the operator.
    function voteConfiguration(uint64 _requestNonce)
        internal
        returns(bool)
    {
        require(configurationNonce == _requestNonce, "nonce mismatch");

        if (voteCommon(VoteType.Configuration, keccak256(msg.data))) {
            configurationNonce++;
            return true;
        }

        return false;
    }

    // registerOperator registers a new operator.
    function registerOperator(address _operator)
    external
    onlyOwner
    {
        require(!operators[_operator]);
        operators[_operator] = true;
        operatorList.push(_operator);
    }

    // deregisterOperator deregisters the operator.
    function deregisterOperator(address _operator)
    external
    onlyOwner
    {
        require(operators[_operator]);
        delete operators[_operator];

        for (uint i = 0; i < operatorList.length; i++) {
           if (operatorList[i] == _operator) {
               operatorList[i] = operatorList[operatorList.length-1];
               operatorList.length--;
               break;
           }
        }
    }

// setOperatorThreshold sets the operator threshold.
    function setOperatorThreshold(VoteType _voteType, uint8 _threshold)
    external
    onlyOwner
    {
        operatorThresholds[uint8(_voteType)] = _threshold;
    }
}
