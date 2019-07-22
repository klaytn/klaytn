pragma solidity ^0.4.24;

import "../externals/openzeppelin-solidity/contracts/token/ERC20/IERC20.sol";
import "../externals/openzeppelin-solidity/contracts/token/ERC20/ERC20Mintable.sol";
import "../externals/openzeppelin-solidity/contracts/math/SafeMath.sol";

contract BridgeFee {
    using SafeMath for uint256;

    address public receiver = address(0);
    uint256 public feeOfKLAY = 0;
    mapping (address => uint256) public feeOfERC20;

    event KLAYFeeChanged(uint256 indexed fee);
    event ERC20FeeChanged(address token, uint256 indexed fee);
    event FeeReceiverChanged(address indexed receiver);

    constructor(address _receiver) internal {
        receiver = _receiver;
    }

    function _payKLAYFeeAndRefundChange(uint256 _feeLimit) internal returns(uint256) {
        uint256 fee = feeOfKLAY;

        if(isValidReceiver() == true && fee > 0) {
            require(_feeLimit >= fee, "invalid feeLimit");

            receiver.transfer(fee);

            uint256 remain = _feeLimit.sub(fee);
            if (remain > 0) {
                msg.sender.transfer(_feeLimit.sub(fee));
            }

            return fee;
        }

        msg.sender.transfer(_feeLimit);
        return 0;
    }

    function _payERC20FeeAndRefundChange(address from, address _token, uint256 _feeLimit) internal returns(uint256){
        uint256 fee = feeOfERC20[_token];

        if (isValidReceiver() == true && fee > 0) {
            require(_feeLimit >= fee, "invalid feeLimit");

            IERC20(_token).transfer(receiver, fee);

            uint256 remain = _feeLimit.sub(fee);
            if (remain > 0) {
                IERC20(_token).transfer(from, _feeLimit.sub(fee));
            }

            return fee;
        }

        IERC20(_token).transfer(from, _feeLimit);
        return 0;
    }

    function _setKLAYFee(uint256 _fee) internal {
        feeOfKLAY = _fee;
        emit KLAYFeeChanged(_fee);
    }

    function _setERC20Fee(address _token, uint256 _fee) internal {
        feeOfERC20[_token] = _fee;
        emit ERC20FeeChanged(_token, _fee);
    }

    function _setFeeReceiver(address _to) internal {
        receiver = _to;
        emit FeeReceiverChanged(_to);
    }

    function isValidReceiver() public view returns(bool){
        if (receiver == address(0)){
            return false;
        }
        return true;
    }
}
