const { BN, constants, expectEvent, expectRevert } = require('@openzeppelin/test-helpers');
const { expect } = require('chai');
const { ZERO_ADDRESS } = constants;
const KIP17MintableMock = artifacts.require('KIP17MintableMock');

const {
  shouldSupportInterfaces,
} = require('../../../utils/introspection/SupportsInterface.behavior');

contract('KIP17Mintable', function (accounts) {
  const [deployer, other, owner] = accounts;

  const name = 'MintableToken';
  const symbol = 'MITOK';

  const firstTokenId = new BN('5042');
  const nonExistentTokenId = new BN('11');

  const DEFAULT_ADMIN_ROLE = '0x0000000000000000000000000000000000000000000000000000000000000000';
  const MINTER_ROLE = web3.utils.soliditySha3('KIP17_MINTER_ROLE');

  beforeEach('create a token type', async function () {
    this.token = await KIP17MintableMock.new(name, symbol, { from: deployer });
  });

  shouldSupportInterfaces([
    'KIP17',
    'AccessControlEnumerable',
  ]);

  describe('minter access control permissions', function () {
    it('deployer has the default admin role', async function () {
      expect(await this.token.getRoleMemberCount(DEFAULT_ADMIN_ROLE)).to.be.bignumber.equal('1');
      expect(await this.token.getRoleMember(DEFAULT_ADMIN_ROLE, 0)).to.equal(deployer);
    });

    it('deployer has the minter role', async function () {
      expect(await this.token.getRoleMemberCount(MINTER_ROLE)).to.be.bignumber.equal('1');
      expect(await this.token.getRoleMember(MINTER_ROLE, 0)).to.equal(deployer);
    });

    it('deployer has the minter role', async function () {
      expect(await this.token.isMinter(deployer)).to.equal(true);
    });

    it('add other account as minter', async function () {
      await this.token.addMinter(other, { from: deployer });
      expect(await this.token.isMinter(other)).to.equal(true);
    });

    it('renounce minter role', async function () {
      await this.token.renounceMinter({ from: deployer });
      expect(await this.token.isMinter(deployer)).to.equal(false);
      expect(await this.token.isMinter(other)).to.equal(false);
    });
  });

  describe('mint', function () {
    it('reverts with a null destination address', async function () {
      await expectRevert(
        this.token.mint(ZERO_ADDRESS, firstTokenId, { from: deployer }), 'KIP17: mint to the zero address',
      );
    });

    beforeEach(async function () {
      (this.receipt = await this.token.mint(owner, firstTokenId, { from: deployer }));
    });

    it('emits a Transfer event', function () {
      expectEvent(this.receipt, 'Transfer', { from: ZERO_ADDRESS, to: owner, tokenId: firstTokenId });
    });

    describe('balanceOf', function () {
      context('when the given address owns some tokens', function () {
        it('returns the amount of tokens owned by the given address', async function () {
          expect(
            await this.token.balanceOf(owner),
          ).to.be.bignumber.equal('1');
        });
      });

      context(
        'when the given address does not own any tokens',
        function () {
          it('returns 0', async function () {
            expect(
              await this.token.balanceOf(other),
            ).to.be.bignumber.equal('0');
          });
        },
      );
    });

    it('reverts when adding a token id that already exists', async function () {
      await expectRevert(
        this.token.mint(owner, firstTokenId, { from: deployer }),
        'KIP17: token already minted',
      );
    });
  });

  describe('safe mint', function () {
    beforeEach(async function () {
      (this.receipt = await this.token.safeMint(
        owner,
        firstTokenId,
      ));
    });

    it('emits a Transfer event', function () {
      expectEvent(this.receipt, 'Transfer', { from: ZERO_ADDRESS, to: owner, tokenId: firstTokenId });
    });

    it('creates the token', async function () {
      expect(await this.token.balanceOf(owner)).to.be.bignumber.equal(
        '1',
      );
      expect(await this.token.ownerOf(firstTokenId)).to.equal(owner);
    });

    it('reverts when adding a token id that already exists', async function () {
      await expectRevert(
        this.token.safeMint(owner, firstTokenId),
        'KIP17: token already minted',
      );
    });
  });

  describe('safe mint with data', function () {
    const data = '0x42';

    // onKIP17Received with data is tested in KIP17.behaviour.js
    beforeEach(async function () {
      (this.receipt = await this.token.safeMint(
        owner,
        firstTokenId,
        data,
      ));
    });

    it('emits a Transfer event', function () {
      expectEvent(this.receipt, 'Transfer',
        { from: ZERO_ADDRESS, to: owner, tokenId: firstTokenId });
    });

    it('creates the token', async function () {
      expect(await this.token.balanceOf(owner)).to.be.bignumber.equal(
        '1',
      );
      expect(await this.token.ownerOf(firstTokenId)).to.equal(owner);
    });

    it('reverts when adding a token id that already exists', async function () {
      await expectRevert(
        this.token.safeMint(owner, firstTokenId),
        'KIP17: token already minted',
      );
    });
  });

  describe('burn after mint', function () {
    it('reverts when burning a non-existent token id', async function () {
      await expectRevert(
        this.token.burn(nonExistentTokenId),
        'KIP17: owner query for nonexistent token',
      );
    });

    context('with minted tokens', function () {
      beforeEach(async function () {
        await this.token.mint(owner, firstTokenId, { from: deployer });
      });

      context('with burnt token', function () {
        beforeEach(async function () {
          ;({ logs: this.logs } = await this.token.burn(firstTokenId));
        });

        it('emits a Transfer event', function () {
          expectEvent.inLogs(this.logs, 'Transfer', {
            from: owner,
            to: ZERO_ADDRESS,
            tokenId: firstTokenId,
          });
        });

        it('emits an Approval event', function () {
          expectEvent.inLogs(this.logs, 'Approval', {
            owner,
            approved: ZERO_ADDRESS,
            tokenId: firstTokenId,
          });
        });

        it('deletes the token', async function () {
          expect(
            await this.token.balanceOf(owner),
          ).to.be.bignumber.equal('0');
          await expectRevert(
            this.token.ownerOf(firstTokenId),
            'KIP17: owner query for nonexistent token',
          );
        });

        it('reverts when burning a token id that has been deleted', async function () {
          await expectRevert(
            this.token.burn(firstTokenId),
            'KIP17: owner query for nonexistent token',
          );
        });
      });
    });
  });
});
