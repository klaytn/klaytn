const { BN, constants, expectEvent, expectRevert } = require('@openzeppelin/test-helpers');
const { expect } = require('chai');
const { ZERO_ADDRESS } = constants;
const KIP17MetadataMintableMock = artifacts.require('KIP17MetadataMintableMock');

const {
  shouldSupportInterfaces,
} = require('../../../utils/introspection/SupportsInterface.behavior');

contract('KIP17MetadataMintableMock', function (accounts) {
  const [deployer, other] = accounts;

  const name = 'MetaMintableToken';
  const symbol = 'MEMIT';
  const baseURI = 'https://api.example.com/v1/';
  const newBaseURI = 'https://api.example.com/v2/';

  const firstTokenId = new BN('5042');
  const nonExistentTokenId = new BN('13');

  const DEFAULT_ADMIN_ROLE =
        '0x0000000000000000000000000000000000000000000000000000000000000000';
  const MINTER_ROLE = web3.utils.soliditySha3('KIP17_MINTER_ROLE');

  beforeEach(async function () {
    this.token = await KIP17MetadataMintableMock.new(
      name,
      symbol,
      { from: deployer },
    );
  });

  shouldSupportInterfaces([
    'KIP17',
    'AccessControlEnumerable',
  ]);

  describe('metadata', function () {
    it('token has correct name', async function () {
      expect(await this.token.name()).to.equal(name);
    });

    it('token has correct symbol', async function () {
      expect(await this.token.symbol()).to.equal(symbol);
    });

    it('deployer has the default admin role', async function () {
      expect(
        await this.token.getRoleMemberCount(DEFAULT_ADMIN_ROLE),
      ).to.be.bignumber.equal('1');
      expect(await this.token.getRoleMember(DEFAULT_ADMIN_ROLE, 0)).to.equal(
        deployer,
      );
    });

    it('deployer has the minter role', async function () {
      expect(
        await this.token.getRoleMemberCount(MINTER_ROLE),
      ).to.be.bignumber.equal('1');
      expect(await this.token.getRoleMember(MINTER_ROLE, 0)).to.equal(
        deployer,
      );
    });

    it('minter role admin is the default admin', async function () {
      expect(await this.token.getRoleAdmin(MINTER_ROLE)).to.equal(
        DEFAULT_ADMIN_ROLE,
      );
    });
  });

  describe('base URI', function () {
    beforeEach(function () {
      if (this.token.setBaseURI === undefined) {
        this.skip();
      }
    });

    it('returns empty if base URI is not set', async function () {
      expect(await this.token.baseURI()).to.equal('');
    });

    it('base URI can be set', async function () {
      await this.token.setBaseURI(baseURI);
      expect(await this.token.baseURI()).to.equal(baseURI);
    });

    it('base URI can be changed', async function () {
      await this.token.setBaseURI(newBaseURI);
      expect(await this.token.baseURI()).to.equal(newBaseURI);
    });
  });

  describe('minting with tokenURI', function () {
    const tokenURI = baseURI + firstTokenId.toString();
    const newTokenURI = 'token/' + firstTokenId.toString();

    beforeEach(async function () {
      await this.token.setBaseURI(baseURI);
      const receipt = await this.token.mintWithTokenURI(
        other, firstTokenId, firstTokenId.toString(), { from: deployer });

      expectEvent(receipt, 'Transfer', {
        to: other,
        from: ZERO_ADDRESS,
        tokenId: firstTokenId,
      });
    });

    it('reverts when queried for non existent token id', async function () {
      await expectRevert(
        this.token.tokenURI(
          nonExistentTokenId), 'KIP17URIStorage: URI query for nonexistent token',
      );
    });

    it('can request token URI set', async function () {
      await expect(this.token.tokenURI(firstTokenId), tokenURI);
    });

    it('token URI can be changes and verify the returning concatenated uri is a token uri', async function () {
      await this.token.setTokenURI(firstTokenId, newTokenURI);
      expect(await this.token.tokenURI(
        firstTokenId)).to.be.equal(baseURI + newTokenURI);
    });
  });

  describe('burn tokens with URI', function () {
    const tokenURI = baseURI + firstTokenId.toString();

    beforeEach(async function () {
      await this.token.setBaseURI(baseURI);
      await this.token.mintWithTokenURI(
        other, firstTokenId, firstTokenId.toString(), { from: deployer });
    });

    it('reverts when burning a non-existent token id', async function () {
      await expectRevert(
        this.token.burn(nonExistentTokenId), 'KIP17: owner query for nonexistent token',
      );
    });

    it('tokens with URI can be burnt ', async function () {
      await expect(this.token.tokenURI(firstTokenId), tokenURI);

      await this.token.burn(firstTokenId, { from: deployer });

      expect(await this.token.exists(firstTokenId)).to.equal(false);
      await expectRevert(
        this.token.tokenURI(firstTokenId),
        'KIP17URIStorage: URI query for nonexistent token',
      );
    });
  });
});
