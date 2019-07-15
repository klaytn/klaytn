pragma solidity ^0.4.24;

import "../externals/openzeppelin-solidity/contracts/token/ERC20/IERC20.sol";
import "../externals/openzeppelin-solidity/contracts/token/ERC20/ERC20Mintable.sol";
import "../externals/openzeppelin-solidity/contracts/math/SafeMath.sol";

contract BridgeFee {
    using SafeMath for uint256;

    address public receiver = address(0);
    uint256 public feeOfKLAY;
    mapping (address => uint256) public feeOfToken;

    event KlayFeeChanged(uint256 indexed fee);
    event ERC20FeeChanged(address token, uint256 indexed fee);
    event FeeReceiverChanged(address indexed receiver);

    constructor(address _receiver) internal {
        receiver = _receiver;
    }

    function payKlayFee(uint256 _fee) internal {
        if(isValidReceiver() == true && _fee > 0) {
            receiver.transfer(_fee);
        }
    }

    function payERC20Fee(address _token, uint256 _fee) internal {
        if(isValidReceiver() == true && _fee > 0) {
            IERC20(_token).transfer(receiver, _fee);
        }
    }

    function _setKlayFee(uint256 _fee) internal {
        feeOfKLAY = _fee;
        emit KlayFeeChanged(_fee);
    }

    function _setERC20Fee(address _token, uint256 _fee) internal {
        feeOfToken[_token] = _fee;
        emit ERC20FeeChanged(_token, _fee);
    }

    function _changeFeeReceiver(address _to) internal {
        receiver = _to;
        emit FeeReceiverChanged(_to);
    }

    function getKLAYFee() public view returns(uint256){
        return feeOfKLAY;
    }

    function getERC20Fee(address _token) public view returns(uint256){
        return feeOfToken[_token];
    }

    function isValidReceiver() public view returns(bool){
        if (receiver == address(0)){
            return false;
        }
        return true;
    }
}
