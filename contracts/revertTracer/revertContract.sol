pragma solidity ^0.4.0;

/**
 * @dev Wrappers over Solidity's arithmetic operations with added overflow
 * checks.
 *
 * Arithmetic operations in Solidity wrap on overflow. This can easily result
 * in bugs, because programmers usually assume that an overflow raises an
 * error, which is the standard behavior in high level programming languages.
 * `SafeMath` restores this intuition by reverting the transaction when an
 * operation overflows.
 *
 * Using this library instead of the unchecked operations eliminates an entire
 * class of bugs, so it's recommended to use it always.
 */
library SafeMath {

    /**
     * @dev Returns the subtraction of two unsigned integers, reverting on
     * overflow (when the result is negative).
     *
     * Counterpart to Solidity's `-` operator.
     *
     * Requirements:
     * - Subtraction cannot overflow.
     */
    function sub(uint256 a, uint256 b) internal pure returns (uint256) {
        require(b <= a, "SafeMath: subtraction overflow");
        uint256 c = a - b;

        return c;
    }
}

contract ContractC {
    using SafeMath for uint256;

    function c1(uint256 calllimit) public returns(uint256) {
        return calllimit.sub(1);
    }
    function c2(uint256 calllimit) public returns(uint256) {
        return calllimit.sub(1);
    }
}

contract ContractB {
    using SafeMath for uint256;

    address contractC = address(0xaBBcD5b340c80B5F1C0545C04C987b87310296aC);

    function b1(uint256 calllimit) public returns(uint256) {
        calllimit = calllimit.sub(1);
        calllimit = ContractC(0xaBBcD5b340c80B5F1C0545C04C987b87310296aC).c1(calllimit);
        calllimit = ContractC(0xaBBcD5b340c80B5F1C0545C04C987b87310296aC).c2(calllimit);
        return calllimit.sub(1);
    }

    function b2(uint256 calllimit) public returns(uint256) {
        calllimit = calllimit.sub(1);
        calllimit = ContractC(0xaBBcD5b340c80B5F1C0545C04C987b87310296aC).c1(calllimit);
        calllimit = ContractC(0xaBBcD5b340c80B5F1C0545C04C987b87310296aC).c2(calllimit);
        return calllimit.sub(1);
    }
}

contract ContractA {
    using SafeMath for uint256;

    function a1(uint256 calllimit) public {
        calllimit = calllimit.sub(1);
        calllimit = ContractB(0xabbcD5b340C80B5F1C0545c04c987b87310296AB).b1(calllimit);
        calllimit = calllimit.sub(1);
    }
}
