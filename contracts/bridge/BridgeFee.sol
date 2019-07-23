pragma solidity ^0.4.24;

import "../externals/openzeppelin-solidity/contracts/token/ERC20/IERC20.sol";
import "../externals/openzeppelin-solidity/contracts/token/ERC20/ERC20Mintable.sol";
import "../externals/openzeppelin-solidity/contracts/math/SafeMath.sol";

contract BridgeFee {
    using SafeMath for uint256;

    address public feeReceiver = address(0);
    uint256 public feeOfKLAY = 0;
    mapping (address => uint256) public feeOfERC20;

    event KLAYFeeChanged(uint256 indexed fee);
    event ERC20FeeChanged(address token, uint256 indexed fee);
    event FeeReceiverChanged(address indexed feeReceiver);

    constructor(address _feeReceiver) internal {
        feeReceiver = _feeReceiver;
    }

    function _payKLAYFeeAndRefundChange(uint256 _feeLimit) internal returns(uint256) {
        uint256 fee = feeOfKLAY;

        if(feeReceiver != address(0) && fee > 0) {
            require(_feeLimit >= fee, "insufficient feeLimit");

            feeReceiver.transfer(fee);
            msg.sender.transfer(_feeLimit.sub(fee));

            return fee;
        }

        msg.sender.transfer(_feeLimit);
        return 0;
    }

    function _payERC20FeeAndRefundChange(address from, address _token, uint256 _feeLimit) internal returns(uint256){
        uint256 fee = feeOfERC20[_token];

        if (feeReceiver != address(0) && fee > 0) {
            require(_feeLimit >= fee, "insufficient feeLimit");

            IERC20(_token).transfer(feeReceiver, fee);
            IERC20(_token).transfer(from, _feeLimit.sub(fee));

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
        feeReceiver = _to;
        emit FeeReceiverChanged(_to);
    }
}
