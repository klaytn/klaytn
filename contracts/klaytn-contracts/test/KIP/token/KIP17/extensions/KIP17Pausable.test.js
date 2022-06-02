const { BN, expectRevert } = require('@openzeppelin/test-helpers');
const { expect } = require('chai');

const KIP17PausableMock = artifacts.require('KIP17PausableMock');

const {
  shouldSupportInterfaces,
} = require('../../../utils/introspection/SupportsInterface.behavior');

contract('KIP17Pausable', function (accounts) {
  const [deployer, other, owner, receiver] = accounts;

  const name = 'Non Fungible Token';
  const symbol = 'NFT';
  const DEFAULT_ADMIN_ROLE = '0x0000000000000000000000000000000000000000000000000000000000000000';
  const PAUSER_ROLE = web3.utils.soliditySha3('KIP17_PAUSER_ROLE');

  beforeEach(async function () {
    this.token = await KIP17PausableMock.new(name, symbol, { from: deployer });
  });

  context('when token is paused', function () {
    const firstTokenId = new BN(1);
    const secondTokenId = new BN(1337);

    const mockData = '0x42';

    beforeEach(async function () {
      await this.token.mint(owner, firstTokenId, { from: deployer });
      await this.token.pause({ from: deployer });
    });

    shouldSupportInterfaces([
      'KIP17',
      'AccessControlEnumerable',
    ]);

    describe('pauser access control permissions', function () {
      it('deployer has the default admin role', async function () {
        expect(await this.token.getRoleMemberCount(DEFAULT_ADMIN_ROLE)).to.be.bignumber.equal('1');
        expect(await this.token.getRoleMember(DEFAULT_ADMIN_ROLE, 0)).to.equal(deployer);
      });

      it('deployer has the pauser role', async function () {
        expect(await this.token.getRoleMemberCount(PAUSER_ROLE)).to.be.bignumber.equal('1');
        expect(await this.token.getRoleMember(PAUSER_ROLE, 0)).to.equal(deployer);
      });

      it('deployer has the pauser role', async function () {
        expect(await this.token.isPauser(deployer)).to.equal(true);
      });

      it('add other account as pauser', async function () {
        await this.token.addPauser(other, { from: deployer });
        expect(await this.token.isPauser(other)).to.equal(true);
      });

      it('renounce minter role', async function () {
        await this.token.renouncePauser({ from: deployer });
        expect(await this.token.isPauser(deployer)).to.equal(false);
        expect(await this.token.isPauser(other)).to.equal(false);
      });
    });

    describe('transfer while paused', function () {
      it('contract should be pasued', async function () {
        expect(await this.token.isPauser(deployer)).to.equal(true);
        expect(await this.token.paused()).to.equal(true);
      });

      it('reverts when trying to transferFrom', async function () {
        await expectRevert(
          this.token.transferFrom(owner, receiver, firstTokenId, {
            from: owner,
          }),
          'KIP17Pausable: token transfer while paused',
        );
      });

      it('reverts when trying to safeTransferFrom', async function () {
        await expectRevert(
          this.token.safeTransferFrom(owner, receiver, firstTokenId, {
            from: owner,
          }),
          'KIP17Pausable: token transfer while paused',
        );
      });

      it('reverts when trying to safeTransferFrom with data', async function () {
        await expectRevert(
          this.token.methods[
            'safeTransferFrom(address,address,uint256,bytes)'
          ](owner, receiver, firstTokenId, mockData, { from: owner }),
          'KIP17Pausable: token transfer while paused',
        );
      });
    });

    describe('mint and burn while paused', function () {
      it('reverts when trying to mint', async function () {
        await expectRevert(
          this.token.mint(receiver, secondTokenId),
          'KIP17Pausable: token transfer while paused',
        );
      });

      it('reverts when trying to burn', async function () {
        await expectRevert(
          this.token.burn(firstTokenId),
          'KIP17Pausable: token transfer while paused',
        );
      });
    });
  });
});
