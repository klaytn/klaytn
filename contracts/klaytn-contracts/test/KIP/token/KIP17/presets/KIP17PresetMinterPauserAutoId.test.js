const {
  BN,
  constants,
  expectEvent,
  expectRevert,
} = require('@openzeppelin/test-helpers');
const { ZERO_ADDRESS } = constants;
const {
  shouldSupportInterfaces,
} = require('../../../utils/introspection/SupportsInterface.behavior');

const { expect } = require('chai');

const KIP17PresetMinterPauserAutoId = artifacts.require(
  'KIP17PresetMinterPauserAutoId',
);

contract('KIP17PresetMinterPauserAutoId', function (accounts) {
  const [deployer, other] = accounts;

  const name = 'MinterAutoIDToken';
  const symbol = 'MAIT';
  const baseURI = 'my.app/';

  const DEFAULT_ADMIN_ROLE =
        '0x0000000000000000000000000000000000000000000000000000000000000000';
  const MINTER_ROLE = web3.utils.soliditySha3('KIP17_MINTER_ROLE');
  const PAUSER_ROLE = web3.utils.soliditySha3('KIP17_PAUSER_ROLE');

  beforeEach(async function () {
    this.token = await KIP17PresetMinterPauserAutoId.new(
      name,
      symbol,
      baseURI,
      { from: deployer },
    );
  });

  shouldSupportInterfaces([
    'KIP17',
    'KIP17Enumerable',
    'AccessControl',
    'AccessControlEnumerable',
  ]);

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

  describe('minting', function () {
    it('deployer can mint tokens', async function () {
      const tokenId = new BN('0');

      const receipt = await this.token.methods['mint(address)'](other, { from: deployer });
      expectEvent(receipt, 'Transfer', {
        from: ZERO_ADDRESS,
        to: other,
        tokenId,
      });

      expect(await this.token.balanceOf(other)).to.be.bignumber.equal('1');
      expect(await this.token.ownerOf(tokenId)).to.equal(other);

      expect(await this.token.tokenURI(tokenId)).to.equal(
        baseURI + tokenId,
      );
    });

    it('other accounts cannot mint tokens', async function () {
      await expectRevert(
        this.token.methods['mint(address)'](other, { from: other }),
        'KIP17PresetMinterPauserAutoId: must have minter role to mint',
      );
    });
  });

  describe('pausing', function () {
    it('deployer has the pauser role', async function () {
      expect(
        await this.token.getRoleMemberCount(PAUSER_ROLE),
      ).to.be.bignumber.equal('1');
      expect(await this.token.getRoleMember(PAUSER_ROLE, 0)).to.equal(
        deployer,
      );
    });

    it('pauser role admin is the default admin', async function () {
      expect(await this.token.getRoleAdmin(PAUSER_ROLE)).to.equal(
        DEFAULT_ADMIN_ROLE,
      );
    });

    it('deployer can pause', async function () {
      const receipt = await this.token.pause({ from: deployer });
      expectEvent(receipt, 'Paused', { account: deployer });

      expect(await this.token.paused()).to.equal(true);
    });

    it('deployer can unpause', async function () {
      await this.token.pause({ from: deployer });

      const receipt = await this.token.unpause({ from: deployer });
      expectEvent(receipt, 'Unpaused', { account: deployer });

      expect(await this.token.paused()).to.equal(false);
    });

    it('cannot mint while paused', async function () {
      await this.token.methods['pause()']({ from: deployer });

      await expectRevert(
        this.token.methods['mint(address)'](other, { from: deployer }),
        'KIP17Pausable: token transfer while paused',
      );
    });

    it('other accounts cannot pause', async function () {
      await expectRevert(
        this.token.pause({ from: other }),
        'KIP17Pausable: must have pauser role to pause',
      );
    });

    it('other accounts cannot unpause', async function () {
      await this.token.pause({ from: deployer });

      await expectRevert(
        this.token.unpause({ from: other }),
        'KIP17Pausable: must have pauser role to unpause',
      );
    });
  });

  describe('burning', function () {
    it('holders can burn their tokens', async function () {
      const tokenId = new BN('0');

      this.token.methods['mint(address)'](other, { from: deployer });

      const receipt = await this.token.burn(tokenId, { from: other });

      expectEvent(receipt, 'Transfer', {
        from: other,
        to: ZERO_ADDRESS,
        tokenId,
      });

      expect(await this.token.balanceOf(other)).to.be.bignumber.equal('0');
      expect(await this.token.totalSupply()).to.be.bignumber.equal('0');
    });
  });
});
