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

import "../externals/openzeppelin-solidity/contracts/token/ERC20/IERC20.sol";
import "../externals/openzeppelin-solidity/contracts/token/ERC20/ERC20Mintable.sol";
import "../externals/openzeppelin-solidity/contracts/math/SafeMath.sol";
import "../externals/openzeppelin-solidity/contracts/ownership/Ownable.sol";


contract BridgeOperator is Ownable {
    using SafeMath for uint256;

    struct VotesData {
        address[] voters;   // voter list for deleting voted map
        mapping(address => bytes32) voted; // <operator, sha3(type, args, nonce)>

        bytes32[] voteKeys; // voteKey list for deleting voteCounts map
        mapping(bytes32 => uint8) voteCounts; // <sha3(type, args, nonce), uint8>
    }

    mapping(uint8 => mapping (uint64 => VotesData)) private votes; // <voteType, <nonce, VotesData>
    mapping(uint64 => bool) public closedValueTransferVotes; // <nonce, bool>
    mapping(uint64 => bool) public closedWithdrawVote; // <nonce, bool>

    mapping(uint64 => bool) public closedRequestRefundVote; // <nonce, bool>
    mapping(uint64 => bool) public closedHandleRefundVote; // <nonce, bool>
    mapping(uint64 => bool) public closedRemoveRefundLedgerVote; // <nonce, bool>

    uint256 internal amountOfLockedRefundKLAY; // amount of locked KLAY for refund
    uint64 private withdrawNonce;

    uint64 public constant MAX_OPERATOR = 12;
    mapping(address => bool) public operators;
    address payable [] public operatorList;

    mapping(uint8 => uint8) public operatorThresholds; // <vote type, uint8>
    uint64 public configurationNonce;

    string public CLOSED_VOTE_ERR_STR = "closed vote";

    enum VoteType {
        ValueTransfer,
        Configuration,
        RefundRequest,
        RefundHandle,
        RefundRemoveLedger,
        Withdraw,
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

    function getOperatorList() external view returns(address payable[] memory) {
        return operatorList;
    }

    // _voteCommon handles common functionality for voting.
    function _voteCommon(VoteType _voteType, uint64 _nonce, bytes32 _voteKey)
        private
        returns(bool)
    {
        VotesData storage vote = votes[uint8(_voteType)][_nonce];

        bytes32 oldVoteKeyOfVoter = vote.voted[msg.sender];
        if (oldVoteKeyOfVoter == bytes32(0)) {
            vote.voters.push(msg.sender);
        } else {
            vote.voteCounts[oldVoteKeyOfVoter]--;
        }

        vote.voted[msg.sender] = _voteKey;

        if (vote.voteCounts[_voteKey] == 0) {
            vote.voteKeys.push(_voteKey);
        }
        vote.voteCounts[_voteKey]++;

        if (vote.voteCounts[_voteKey] >= operatorThresholds[uint8(_voteType)]) {
            _removeVoteData(_voteType, _nonce);
            return true;
        }
        return false;
    }

    // _removeVoteData removes a vote data according to voteType and nonce.
    function _removeVoteData(VoteType _voteType, uint64 _nonce)
        internal
    {
        VotesData storage vote = votes[uint8(_voteType)][_nonce];

        for (uint8 i = 0; i < vote.voters.length; i++) {
            delete vote.voted[vote.voters[i]];
        }

        for (uint8 i = 0; i < vote.voteKeys.length; i++) {
            delete vote.voteCounts[vote.voteKeys[i]];
        }

        delete votes[uint8(_voteType)][_nonce];
    }

    // _voteValueTransfer votes value transfer transaction with the operator.
    function _voteValueTransfer(uint64 _requestNonce)
        internal
        returns(bool)
    {
        require(!closedValueTransferVotes[_requestNonce], CLOSED_VOTE_ERR_STR);

        bytes32 voteKey = keccak256(msg.data);
        if (_voteCommon(VoteType.ValueTransfer, _requestNonce, voteKey)) {
            closedValueTransferVotes[_requestNonce] = true;
            return true;
        }

        return false;
    }

    // _voteConfiguration votes contract configuration transaction with the operator.
    function _voteConfiguration(uint64 _requestNonce)
        internal
        returns(bool)
    {
        require(configurationNonce == _requestNonce, "nonce mismatch");

        bytes32 voteKey = keccak256(msg.data);
        if (_voteCommon(VoteType.Configuration, _requestNonce, voteKey)) {
            configurationNonce++;
            return true;
        }

        return false;
    }

    function _voteRemoveRefundLedgerVote(uint64 _requestNonce)
        internal
        returns(bool)
    {
        require(!closedRemoveRefundLedgerVote[_requestNonce], CLOSED_VOTE_ERR_STR);

        bytes32 voteKey = keccak256(msg.data);
        if (_voteCommon(VoteType.RefundRemoveLedger, _requestNonce, voteKey)) {
            closedRemoveRefundLedgerVote[_requestNonce] = true;
            return true;
        }

        return false;
    }

    function _voteRequestRefund(uint64 _requestNonce)
        internal
        returns(bool)
    {
        require(!closedRequestRefundVote[_requestNonce], CLOSED_VOTE_ERR_STR);

        bytes32 voteKey = keccak256(msg.data);
        if (_voteCommon(VoteType.RefundRequest, _requestNonce, voteKey)) {
            closedRequestRefundVote[_requestNonce] = true;
            return true;
        }

        return false;
    }

    function _voteHandleRefund(uint64 _requestNonce)
        internal
        returns(bool)
    {
        require(!closedHandleRefundVote[_requestNonce], CLOSED_VOTE_ERR_STR);

        bytes32 voteKey = keccak256(msg.data);
        if (_voteCommon(VoteType.RefundHandle, _requestNonce, voteKey)) {
            closedHandleRefundVote[_requestNonce] = true;
            return true;
        }

        return false;
    }

    function _voteWithdraw(uint64 withdrawNonce_) 
        internal
        returns(bool)
    {
        require(!closedWithdrawVote[withdrawNonce_], CLOSED_VOTE_ERR_STR);
        
        bytes32 voteKey = keccak256(msg.data);
        if (_voteCommon(VoteType.Withdraw, withdrawNonce_, voteKey)) {
            closedWithdrawVote[withdrawNonce_] = true;
            return true;
        }

        return false;
    }

    function withdrawKLAY(uint256 value) public onlyOperators {
        require(operatorList.length > 0, "A number of operators < 0");
        if (!_voteWithdraw(withdrawNonce)) {
            return;
        }
        uint256 divided = value.sub(amountOfLockedRefundKLAY) / uint256(operatorList.length);
        for (uint i=0; i<operatorList.length; i++) {
            operatorList[i].transfer(divided); 
        }
        withdrawNonce++;
    }

    /*
     * TODO-hyunsooda: Implement both of them
    function withdrawERC20() public onlyOperators {}
    function withdrawERC721() public onlyOperators {}
    */

    /*
     * TODO-hyunsooda: Implement both of them
    function getAmountOfLockedRefundERC20() public onlyOperators {}
    function getAmountOfLockedRefundERC721() public onlyOperators {}
    */

    // registerOperator registers a new operator.
    function registerOperator(address payable _operator)
    external
    onlyOwner
    {
        require(operatorList.length < MAX_OPERATOR, "max operator limit");
        require(!operators[_operator], "exist operator");
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
        require(_threshold > 0, "zero threshold");
        require(operatorList.length >= _threshold, "bigger than num of operators");
        operatorThresholds[uint8(_voteType)] = _threshold;
    }
}
