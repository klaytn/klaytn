pragma solidity >=0.8.7;

/** @title nodeWhitelist
  * @notice This contract
        * Works as the the manager of node whitelist. (This contract should be called to add or delete node from CBDC network)
        * Handles permission for CBDC network.
  * @dev ...
  */
contract NodeWhitelist {
    address private admin;
    string[] private nodes;

    event AddNode(address admin, string addedNode);
    event DelNode(address admin, string deletedNode);

    modifier onlyAdmin() {
        require(msg.sender == admin, "Only admin can call.");
        _;
    }

    function getAdmin() public view returns (address) {return admin;}

    function setAdmin(address newAdmin) external onlyAdmin {admin = newAdmin;}

    function getWhitelist() public view returns (string[] memory) {
        return nodes;
    }

    function addNode(string memory node) public onlyAdmin {
        require(indexOf(node) == - 1, "given node is already registered!");
        nodes.push(node);
        emit AddNode(admin, node);
    }

    function delNode(string memory node) public onlyAdmin {
        require(indexOf(node) != - 1, "given node is not on the whitelist!");
        remove(uint(indexOf(node)));
        emit DelNode(admin, node);
    }

    // below remove does not preserve the order
    function remove(uint index) internal {
        require(index < nodes.length);
        nodes[index] = nodes[nodes.length - 1];
        nodes.pop();
    }

    function indexOf(string memory node) internal returns (int) {
        for (uint i = 0; i < nodes.length; i++) {
            // below is to avoid error "operator == not compatible with type string storage ref and literal_string"
            if (keccak256(bytes(nodes[i])) == keccak256(bytes(node))) {
                return int(i);
            }
        }
        return - 1;
    }
}