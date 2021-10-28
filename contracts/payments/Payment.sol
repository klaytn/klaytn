pragma solidity >=0.8.0;

contract Payment {
    // The `Queue` holds the transaction hashes.
    // The transaction data can be view through `txs`.
    mapping (address => Queue) private receiverInfos;
    // This transaction hashes are not same those in blockchain.
    // This hashes are created when `sendPayment` is called.
    mapping (bytes32 => SendInfo) private txs;

    uint private SettleInterval;
    uint private sequenceNumber;

    event SendPayment(address sender, address receiver, uint amount, bytes32 hash);
    event CancelPayment(address sender, address receiver, uint amount, bytes32 hash);
    event SettlePayment(address sender, address receiver, uint amount, bytes32 hash);
    event Transfer(address from, address to, uint256 value);

    constructor (uint _settleInterval) {
        SettleInterval = _settleInterval;
        if (_settleInterval == 0) {
            SettleInterval = 2592000; // 2592000 sec = 60 * 60 * 24 * 30 = 30 days
        }
    }

    modifier txHashExists(bytes32 txHash) {
        require(txs[txHash].amount != 0, "Such txHash does not exist.");
        _;
    }
    modifier onlyReceiver(bytes32 txHash) {
        require(msg.sender == txs[txHash].receiver, "Only receiver can call.");
        _;
    }

    function sendPayment(address receiver) public payable {
        require(msg.value != 0, "amount should not be 0");

        if (address(receiverInfos[receiver]) == address(0)) {
            receiverInfos[receiver] = new Queue();
        }

        // create data
        SendInfo memory sendInfo = SendInfo(msg.sender, receiver, msg.value, block.number);
        bytes32 hash = sha256(abi.encodePacked(msg.sender, receiver, msg.value, block.number, sequenceNumber));
        unchecked { // https://docs.soliditylang.org/en/latest/control-structures.html#checked-or-unchecked-arithmetic
            sequenceNumber++;
        }
        // store data
        receiverInfos[receiver].enqueue(hash);
        txs[hash] = sendInfo;

        emit SendPayment(msg.sender, receiver, msg.value, hash);
    }

    function cancelPayment(bytes32 txHash) public txHashExists(txHash) onlyReceiver(txHash) {
        // give back money to sender
        address sender = txs[txHash].sender;
        uint amount = txs[txHash].amount;
        payable(sender).transfer(amount);

        emit Transfer(address(this), sender, amount); // ERC20 Transfer
        emit SendPayment(txs[txHash].sender, txs[txHash].receiver, txs[txHash].amount, txHash);

        delete txs[txHash];
    }

    function settlePayment() public {
        // check if the queue exits
        if (address(receiverInfos[msg.sender]) == address(0)) {
            return;
        }

        // check if the queue exits
        Queue q = receiverInfos[msg.sender];
        while(q.peek() != 0x0000000000000000000000000000000000000000000000000000000000000000 // check if empty
            && isAbleToSettleInternal(q.peek())) {

            // get txHash
            bytes32 txHash = q.dequeue();

            // skip sending money if it is canceled
            if (txs[txHash].amount == 0) {
                continue;
            }

            // give money to receiver
            address receiver = txs[txHash].receiver;
            uint amount = txs[txHash].amount;
            payable(receiver).transfer(amount);

            emit Transfer(address(this), receiver, amount); // ERC20 Transfer
            emit SettlePayment(txs[txHash].sender, txs[txHash].receiver, txs[txHash].amount, txHash);

            delete txs[txHash];
        }
    }

    function isAbleToSettleInternal(bytes32 txHash) internal view returns (bool) {
        SendInfo memory sendinfo = txs[txHash];
        return sendinfo.blockNumber + SettleInterval <= block.number;
    }

    function isAbleToSettle(bytes32 txHash) public view txHashExists(txHash) returns (bool) {
        return isAbleToSettleInternal(txHash);
    }

    function getPaymentInfo(bytes32 txHash) public view txHashExists(txHash) returns (SendInfo memory) {
        SendInfo memory sendInfo = txs[txHash];
        return sendInfo;
    }

    function getPayments(address receiver) public view returns (bytes32[] memory) {
        // check if the queue exits
        if (address(receiverInfos[receiver]) == address(0)) {
            bytes32[] memory empty;
            return empty;
        }

        // get tx hashes
        bytes32[] memory hashes = receiverInfos[receiver].getAll();
        // little trick to filter array
        uint256 resultCount;
        for (uint i = 0; i < hashes.length; i++) {
            bytes32 hash = hashes[i];
            if (txs[hash].amount != 0) {// check if canceled
                resultCount++;
            }
        }

        // check if tx exists
        bytes32[] memory result = new bytes32[](resultCount);
        uint256 j = 0;
        for (uint i = 0; i < hashes.length; i++) {
            bytes32 hash = hashes[i];
            if (txs[hash].amount != 0) {// check if canceled
                result[j] = hash;
                j ++;
            }
        }
        return result;
    }

    function getSettleablePayments(address receiver) public view returns (bytes32[] memory) {
        // check if the queue exits
        if (address(receiverInfos[receiver]) == address(0)) {
            bytes32[] memory empty;
            return empty;
        }

        // get tx hashes
        bytes32[] memory hashes = receiverInfos[receiver].getAll();
        // little trick to filter array
        uint256 resultCount;
        for (uint i = 0; i < hashes.length; i++) {
            bytes32 hash = hashes[i];
            if (txs[hash].amount != 0 // check if canceled
                && isAbleToSettleInternal(hash)) {// check if able to settle
                resultCount++;
            }
        }

        // check if tx exists
        bytes32[] memory result = new bytes32[](resultCount);
        uint256 j = 0;
        for (uint i = 0; i < hashes.length; i++) {
            bytes32 hash = hashes[i];
            if (txs[hash].amount != 0 // check if canceled
                && isAbleToSettleInternal(hash)) {// check if able to settle
                result[j] = hash;
                j ++;
            }
        }
        return result;
    }
    function getSettleInterval() public view returns(uint256) {
        return SettleInterval;
    }
}

contract Queue {
    // key   : index (order)
    // value : transaction hashes
    mapping(uint256 => bytes32) private queue;
    uint256 private first = 1;
    uint256 private last = 0;
    address owner;

    modifier onlyOwner() {
        require(msg.sender == owner, "Only contract owner can call.");
        _;
    }

    constructor() {
        owner = msg.sender;
    }

    function enqueue(bytes32 data) public onlyOwner {
        last += 1;
        queue[last] = data;
    }

    function dequeue() public onlyOwner returns (bytes32) {
        require(last >= first, "Unable to dequeue an empty queue");

        bytes32 data = queue[first];

        delete queue[first];
        first += 1;

        return data;
    }

    function peek() public view onlyOwner returns (bytes32) {
        if (last < first) {
            return 0x0000000000000000000000000000000000000000000000000000000000000000;
        }
        return queue[first];
    }

    function getAll() public view onlyOwner returns (bytes32[] memory) {
        bytes32[] memory result = new bytes32[](last + 1 - first);
        if (last < first) {
            return result;
        }

        for (uint256 i = first; i <= last; i++) {
            result[i - first] = queue[i];
        }

        return result;
    }
}

    struct SendInfo {
        address sender;
        address receiver;
        uint amount;
        uint256 blockNumber;
    }
