pragma solidity >=0.8.0;

contract Sender {
    constructor () payable {}

    function deposit(uint256 amount) payable public {
        require(msg.value == amount);
    }

    function sendKlay(uint32 times, address payable receiver) public {
        for (uint32 i = 0; i < times; i++) {
            receiver.transfer(1);
        }
    }

    function getBalance() public view returns (uint256) {
        return address(this).balance;
    }
}
