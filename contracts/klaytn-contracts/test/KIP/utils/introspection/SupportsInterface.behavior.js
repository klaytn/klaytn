const { makeInterfaceId } = require('@openzeppelin/test-helpers');
var should = require('chai').should();

const INTERFACES = {
  KIP13: [
    'supportsInterface(bytes4)',
  ],
  KIP17: [
    'balanceOf(address)',
    'ownerOf(uint256)',
    'approve(address,uint256)',
    'getApproved(uint256)',
    'setApprovalForAll(address,bool)',
    'isApprovedForAll(address,address)',
    'transferFrom(address,address,uint256)',
    'safeTransferFrom(address,address,uint256)',
    'safeTransferFrom(address,address,uint256,bytes)',
  ],
  KIP17Metadata: [
    'name()',
    'symbol()',
    'tokenURI(uint256)',
  ],
  KIP17Enumerable: [
    'totalSupply()',
    'tokenOfOwnerByIndex(address,uint256)',
    'tokenByIndex(uint256)',
  ],
  KIP17Mintable: [
    'mint(address,uint256)',
    'isMinter(address)',
    'addMinter(address)',
    'renounceMinter(address)',
  ],
  KIP17Pausable: [
    'paused()',
    'pause()',
    'unpause()',
    'isPauser(address)',
    'addPauser(address)',
    'renouncePauser()',
  ],
  KIP17Burnable: [
    'burn(unit256)',
  ],
  KIP37: [
    'balanceOf(address,uint256)',
    'balanceOfBatch(address[],uint256[])',
    'setApprovalForAll(address,bool)',
    'isApprovedForAll(address,address)',
    'safeTransferFrom(address,address,uint256,uint256,bytes)',
    'safeBatchTransferFrom(address,address,uint256[],uint256[],bytes)',
  ],
  KIP37Receiver: [
    'onKIP37Received(address,address,uint256,uint256,bytes)',
    'onKIP37BatchReceived(address,address,uint256[],uint256[],bytes)',
  ],
  AccessControl: [
    'hasRole(bytes32,address)',
    'getRoleAdmin(bytes32)',
    'grantRole(bytes32,address)',
    'revokeRole(bytes32,address)',
    'renounceRole(bytes32,address)',
  ],
  AccessControlEnumerable: [
    'getRoleMember(bytes32,uint256)',
    'getRoleMemberCount(bytes32)',
  ],
  ERC165: [
    'supportsInterface(bytes4)',
  ],
};

const INTERFACE_IDS = {};
const FN_SIGNATURES = {};
for (const k of Object.getOwnPropertyNames(INTERFACES)) {
  INTERFACE_IDS[k] = makeInterfaceId.ERC165(INTERFACES[k]); //can't use kip13 as it requires to be included and exported in @openzeppelin/test-helpers(https://github.com/klaytn/klaytn-contracts/blob/klaytn-migration/node_modules/@openzeppelin/test-helpers/src/makeInterfaceId.js)
  for (const fnName of INTERFACES[k]) {
    // the interface id of a single function is equivalent to its function signature
    FN_SIGNATURES[fnName] = makeInterfaceId.ERC165([fnName]);
  }
}

function shouldSupportInterfaces (interfaces = []) {
  describe('KIP13', function () {
    beforeEach(function () {
      this.contractUnderTest = this.mock || this.token || this.holder || this.accessControl;
    });

    it('supportsInterface uses less than 30k gas', async function () {
      for (const k of interfaces) {
        const interfaceId = INTERFACE_IDS[k];
        expect(await this.contractUnderTest.supportsInterface.estimateGas(interfaceId)).to.be.lte(30000);
      }
    });

    it('all interfaces are reported as supported', async function () {
      for (const k of interfaces) {
        const interfaceId = INTERFACE_IDS[k];
        expect(await this.contractUnderTest.supportsInterface(interfaceId)).to.equal(true);
      }
    });

    it('all interface functions are in ABI', async function () {
      for (const k of interfaces) {
        for (const fnName of INTERFACES[k]) {
          const fnSig = FN_SIGNATURES[fnName];
          expect(this.contractUnderTest.abi.filter(fn => fn.signature === fnSig).length).to.equal(1);
        }
      }
    });
  });
}

module.exports = {
  shouldSupportInterfaces,
};
