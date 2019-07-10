pragma solidity ^0.4.24;

import "../externals/openzeppelin-solidity/contracts/token/ERC20/ERC20.sol";
import "../externals/openzeppelin-solidity/contracts/token/ERC20/ERC20Mintable.sol";
import "../externals/openzeppelin-solidity/contracts/token/ERC20/ERC20Burnable.sol";

import "../externals/openzeppelin-solidity/contracts/utils/Address.sol";
import "./ERC20ServiceChain.sol";

contract ServiceChainToken is ERC20, ERC20Mintable, ERC20Burnable, ERC20ServiceChain {
    string public constant name = "ServiceChainToken";
    string public constant symbol = "SCT";
    uint8 public constant decimals = 18;

    // one billion in initial supply
    uint256 public constant INITIAL_SUPPLY = 1000000000 * (10 ** uint256(decimals));

    constructor (address _bridge) ERC20ServiceChain(_bridge) public {
        _mint(msg.sender, INITIAL_SUPPLY);
    }
}
