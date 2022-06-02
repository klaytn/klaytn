const {
  BN,
  constants,
  expectEvent,
  expectRevert,
} = require('@openzeppelin/test-helpers');
const { ZERO_ADDRESS } = constants;

const { expect } = require('chai');

const KIP7Mintable = artifacts.require('KIP7MintableMock');

contract('KIP7Mintable', function (accounts) {
  const [deployer, other] = accounts;

  const name = 'MinterToken';
  const symbol = 'MT';
  const amount = new BN('5000');

  const DEFAULT_ADMIN_ROLE =
        '0x0000000000000000000000000000000000000000000000000000000000000000';
  const MINTER_ROLE = web3.utils.soliditySha3('KIP7_MINTER_ROLE');

  beforeEach(async function () {
    this.token = await KIP7Mintable.new(name, symbol, deployer, amount, {
      from: deployer,
    });
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
      const receipt = await this.token.mint(other, amount, {
        from: deployer,
      });
      expectEvent(receipt, 'Transfer', {
        from: ZERO_ADDRESS,
        to: other,
        value: amount,
      });

      expect(await this.token.balanceOf(other)).to.be.bignumber.equal(
        amount,
      );
    });

    it('other accounts cannot mint tokens', async function () {
      await expectRevert(
        this.token.mint(other, amount, { from: other }),
        'KIP7Mintable: must have minter role to mint',
      );
    });
  });

  describe('burning', function () {
    it('holders can burn their tokens', async function () {
      await this.token.mint(other, amount, { from: deployer });

      const receipt = await this.token.burn(other, amount.subn(1), {
        from: other,
      });
      expectEvent(receipt, 'Transfer', {
        from: other,
        to: ZERO_ADDRESS,
        value: amount.subn(1),
      });

      expect(await this.token.balanceOf(other)).to.be.bignumber.equal('1');
    });
  });
});
