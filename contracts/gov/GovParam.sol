//SPDX-License-Identifier: Unlicense
pragma solidity ^0.8.0;

import "./Ownable.sol";

contract GovParam is Ownable {
    string[] private paramNames;
    mapping(string => Param) private params;

    struct Param {
        bool   exists;     // True for all params
        uint64 activation; // block number where nextValue is first applied
        bytes  prevValue;  // RLP encoded value
        bytes  nextValue;  // RLP encoded value
    }

    struct ParamView {
        string name;
        bytes  value;
    }

    event SetParam(string, bytes, uint64);

    constructor(address _owner) {
        if (_owner != address(0)) {
            _transferOwnership(_owner);
        }
    }

    function setParam(string calldata name, bytes calldata value, uint64 activation)
    external onlyOwner {
        require(bytes(name).length > 0,                 "name cannot be empty");
        require(params[name].activation < block.number, "already have a pending change");
        require(activation > block.number,              "activation must be in a future");

        if (!params[name].exists) { // First time setting this param.
            params[name].exists = true;
            paramNames.push(name);
        }

        params[name].activation = activation;
        params[name].prevValue  = params[name].nextValue;
        params[name].nextValue  = value;

        emit SetParam(name, value, activation);
    }

    function getParamAt(string memory name, uint num)
    private view returns (bytes memory) {
        // The Past could have overwritten by past setParam() calls.
        // The Future can be changed by subsequent subParam() calls.
        // What we are certain is only current (head) and pending block (head+1).
        require(num >= block.number,     "cannot query the past");
        require(num <= block.number + 1, "cannot query the tentative future");
        if (!params[name].exists) {
            return "";
        }
        if (num >= params[name].activation) {
            return params[name].nextValue;
        } else {
            return params[name].prevValue;
        }
    }

    function getAllParamsAt(uint num)
    private view returns (ParamView[] memory) {
        ParamView[] memory views = new ParamView[](paramNames.length);
        for (uint i = 0; i < paramNames.length; i++) {
            views[i].name = paramNames[i];
            views[i].value = getParamAt(paramNames[i], num);
        }
        return views;
    }

    function getParam(string memory name) external view returns (bytes memory) {
        return getParamAt(name, block.number);
    }
    function getParamAtNextBlock(string memory name) external view returns (bytes memory) {
        return getParamAt(name, block.number + 1);
    }
    function getAllParams() external view returns (ParamView[] memory) {
        return getAllParamsAt(block.number);
    }
    function getAllParamsAtNextBlock() external view returns (ParamView[] memory) {
        return getAllParamsAt(block.number + 1);
    }
}
