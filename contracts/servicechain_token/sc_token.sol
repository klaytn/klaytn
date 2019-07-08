pragma solidity ^0.4.24;

import "../externals/openzeppelin-solidity/contracts/token/ERC20/ERC20.sol";
import "../externals/openzeppelin-solidity/contracts/token/ERC20/ERC20Mintable.sol";
import "../externals/openzeppelin-solidity/contracts/token/ERC20/ERC20Burnable.sol";

import "../externals/openzeppelin-solidity/contracts/utils/Address.sol";
import "./ITokenReceiver.sol";


contract ServiceChainToken is ERC20, ERC20Mintable, ERC20Burnable {
    string public constant name = "ServiceChainToken";
    string public constant symbol = "SCT";
    uint8 public constant decimals = 18;

    address bridge;

    bytes4 constant _ERC20_RECEIVED = 0xbc04f0af;

    using Address for address;

    // one billion in initial supply
    uint256 public constant INITIAL_SUPPLY = 1000000000 * (10 ** uint256(decimals));

    constructor (address _bridge) public {
        _mint(msg.sender, INITIAL_SUPPLY);

        if (!_bridge.isContract()) {
            revert("bridge is not a contract");
        }

        bridge = _bridge;
    }

    function requestValueTransfer(uint256 _amount, address _to) external {
        transfer(bridge, _amount);

        bytes4 retval = ITokenReceiver(bridge).onTokenReceived(msg.sender, _amount, _to);
        require(retval == _ERC20_RECEIVED, "Sent to a bridge which is not an ERC20 receiver");
    }
}
