pragma solidity ^0.4.24;

import "../externals/openzeppelin-solidity/contracts/token/ERC20/IERC20.sol";
import "../externals/openzeppelin-solidity/contracts/token/ERC20/ERC20Mintable.sol";
import "../externals/openzeppelin-solidity/contracts/math/SafeMath.sol";
import "../externals/openzeppelin-solidity/contracts/ownership/Ownable.sol";

contract BridgeOperator is Ownable {
    mapping(address => bool) public operators;
    mapping(bytes32 => mapping(address => bool)) public signedTxs; // <sha3(type, args, nonce), <singer, vote>>
    mapping(bytes32 => uint64) public signedTxsCounts; // <sha3(type, args, nonce)>
    mapping(bytes32 => bool) public committedTxs; // <sha3(type, nonce)>
    mapping(uint64 => uint64) public operatorThresholds; // <tx type>
    mapping(uint64 => uint64) public configurationNonces; // <tx type, nonce>

    enum VoteType {
        ValueTransfer,
        Configuration,
        Max
    }

    constructor() internal {
        for (uint64 i = 0; i < uint64(VoteType.Max); i++) {
            operatorThresholds[uint64(i)] = 1;
        }
    }

    modifier onlyOperators()
    {
        require(operators[msg.sender]);
        _;
    }

    // onlySequentialNonce checks sequential nonce increase.
    function onlySequentialNonce(VoteType _voteType, uint64 _requestNonce) internal view {
        require(configurationNonces[uint64(_voteType)] == _requestNonce, "nonce mismatch");
    }

    // voteValueTransfer votes value transfer transaction with the operator.
    function voteValueTransfer(bytes32 _txKey, bytes32 _voteKey, address _operator) internal returns(bool) {
        if (committedTxs[_txKey] || signedTxs[_voteKey][_operator]) {
            return false;
        }

        signedTxs[_voteKey][_operator] = true;
        signedTxsCounts[_voteKey]++;

        if (signedTxsCounts[_voteKey] == operatorThresholds[uint64(VoteType.ValueTransfer)]) {
            committedTxs[_txKey] = true;
            return true;
        }

        return false;
    }

    // voteConfiguration votes contract configuration transaction with the operator.
    function voteConfiguration(bytes32 _voteKey, address _operator)
    internal
    returns(bool)
    {
        if (signedTxs[_voteKey][_operator]) {
            return false;
        }

        signedTxs[_voteKey][_operator] = true;
        signedTxsCounts[_voteKey]++;

        if (signedTxsCounts[_voteKey] ==  operatorThresholds[uint64(VoteType.Configuration)]) {
            configurationNonces[uint64(VoteType.Configuration)]++;
            return true;
        }
        return false;
    }
}