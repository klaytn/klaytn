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

const KIP37PresetMinterPauser = artifacts.require('KIP37PresetMinterPauser');

contract('KIP37PresetMinterPauser', function (accounts) {
  const [deployer, other] = accounts;

  const firstTokenId = new BN('845');
  const firstTokenInitialSupply = new BN('5000000');
  const firstTokenIdAmount = new BN('50');

  const secondTokenId = new BN('48324');
  const secondTokenInitialSupply = new BN('7000000');
  const secondTokenIdAmount = new BN('77875');

  const DEFAULT_ADMIN_ROLE =
        '0x0000000000000000000000000000000000000000000000000000000000000000';
  const MINTER_ROLE = web3.utils.soliditySha3('KIP37_MINTER_ROLE');
  const PAUSER_ROLE = web3.utils.soliditySha3('KIP37_PAUSER_ROLE');

  const uri = 'https://token.com';

  beforeEach(async function () {
    this.token = await KIP37PresetMinterPauser.new(uri, { from: deployer });
  });

  shouldSupportInterfaces([
    'KIP37',
    'AccessControl',
    'AccessControlEnumerable',
  ]);

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

  it('deployer has the pauser role', async function () {
    expect(
      await this.token.getRoleMemberCount(PAUSER_ROLE),
    ).to.be.bignumber.equal('1');
    expect(await this.token.getRoleMember(PAUSER_ROLE, 0)).to.equal(
      deployer,
    );
  });

  it('minter and pauser role admin is the default admin', async function () {
    expect(await this.token.getRoleAdmin(MINTER_ROLE)).to.equal(
      DEFAULT_ADMIN_ROLE,
    );
    expect(await this.token.getRoleAdmin(PAUSER_ROLE)).to.equal(
      DEFAULT_ADMIN_ROLE,
    );
  });

  describe('minting', function () {
    beforeEach(async function () {
      await this.token.create(firstTokenId, firstTokenInitialSupply, uri, { from: deployer });
    });

    it('deployer can mint tokens', async function () {
      const receipt = await this.token.methods['mint(uint256,address,uint256)'](
        firstTokenId,
        other,
        firstTokenIdAmount,
        { from: deployer },
      );
      expectEvent(receipt, 'TransferSingle', {
        operator: deployer,
        from: ZERO_ADDRESS,
        to: other,
        amount: firstTokenIdAmount,
        id: firstTokenId,
      });

      expect(
        await this.token.balanceOf(other, firstTokenId),
      ).to.be.bignumber.equal(firstTokenIdAmount);
    });

    it('other accounts cannot mint tokens', async function () {
      await expectRevert(
        this.token.methods['mint(uint256,address,uint256)'](firstTokenId, other, firstTokenIdAmount, { from: other }),
        'KIP37: must have minter role to mint',
      );
    });
  });

  describe('batched minting', function () {
    const data = web3.utils.soliditySha3(uri);

    beforeEach(async function () {
      this.receipt = await this.token.create(firstTokenId, firstTokenInitialSupply, data, { from: deployer });
      this.receipt = await this.token.create(secondTokenId, secondTokenInitialSupply, data, { from: deployer });
    });

    it('deployer can batch mint tokens', async function () {
      const receipt = await this.token.mintBatch(
        other,
        [firstTokenId, secondTokenId],
        [firstTokenIdAmount, secondTokenIdAmount],
        { from: deployer },
      );

      expectEvent(receipt, 'TransferBatch', {
        operator: deployer,
        from: ZERO_ADDRESS,
        to: other,
      });

      expect(
        await this.token.balanceOf(other, firstTokenId),
      ).to.be.bignumber.equal(firstTokenIdAmount);
    });

    it('other accounts cannot batch mint tokens', async function () {
      await expectRevert(
        this.token.mintBatch(
          other,
          [firstTokenId, secondTokenId],
          [firstTokenIdAmount, secondTokenIdAmount],
          { from: other },
        ),
        'KIP37: must have minter role to mint',
      );
    });
  });

  describe('pausing', function () {
    beforeEach(async function () {
      await this.token.create(firstTokenId, firstTokenInitialSupply, uri, { from: deployer });
    });

    it('deployer can pause', async function () {
      const receipt = await this.token.methods['pause()']({ from: deployer });
      expectEvent(receipt, 'Paused', { account: deployer });

      expect(await this.token.paused()).to.equal(true);
    });

    it('deployer can unpause', async function () {
      await this.token.methods['pause()']({ from: deployer });

      const receipt = await this.token.methods['unpause()']({ from: deployer });
      expectEvent(receipt, 'Unpaused', { account: deployer });

      expect(await this.token.paused()).to.equal(false);
    });

    it('cannot mint while paused', async function () {
      await this.token.methods['pause()']({ from: deployer });

      await expectRevert(
        this.token.methods['mint(uint256,address,uint256)'](firstTokenId, other, firstTokenIdAmount, {
          from: deployer,
        }),
        'KIP37Pausable: token transfer while paused',
      );
    });

    it('other accounts cannot pause', async function () {
      await expectRevert(
        this.token.methods['pause()']({ from: other }),
        'KIP37PresetMinterPauser: must have pauser role to pause',
      );
    });

    it('other accounts cannot unpause', async function () {
      await this.token.methods['pause()']({ from: deployer });

      await expectRevert(
        this.token.methods['unpause()']({ from: other }),
        'KIP37PresetMinterPauser: must have pauser role to unpause',
      );
    });
  });

  describe('pausing and unpausing with id', function () {
    beforeEach(async function () {
      await this.token.create(firstTokenId, firstTokenInitialSupply, uri, { from: deployer });
    });

    it('deployer can pause', async function () {
      const receipt = await this.token.pause(firstTokenId, { from: deployer });

      expectEvent(receipt, 'TokenPaused', { account: deployer });
      expect(await this.token.paused(firstTokenId)).to.equal(true);
    });

    it('deployer can unpause', async function () {
      await this.token.pause(firstTokenId, { from: deployer });

      const receipt = await this.token.methods['unpause(uint256)'](firstTokenId, { from: deployer });
      expectEvent(receipt, 'TokenUnpaused', { account: deployer });

      expect(await this.token.paused(firstTokenId)).to.equal(false);
    });

    it('cannot mint while paused', async function () {
      await this.token.pause(firstTokenId, { from: deployer });

      await expectRevert(
        this.token.methods['mint(uint256,address,uint256)'](firstTokenId, other, firstTokenIdAmount, {
          from: deployer,
        }),
        'KIP37Pausable: token transfer while paused',
      );
    });

    it('other accounts cannot pause', async function () {
      await expectRevert(
        this.token.pause(firstTokenId, { from: other }),
        'KIP37PresetMinterPauser: must have pauser role to pause',
      );
    });

    it('other accounts cannot unpause', async function () {
      await this.token.methods['pause(uint256)'](firstTokenId, { from: deployer });

      await expectRevert(
        this.token.methods['unpause(uint256)'](firstTokenId, { from: other }),
        'KIP37PresetMinterPauser: must have pauser role to unpause',
      );
    });
  });

  describe('burning', function () {
    beforeEach(async function () {
      await this.token.create(firstTokenId, firstTokenInitialSupply, uri, { from: deployer });
    });

    it('holders can burn their tokens', async function () {
      await this.token.methods['mint(uint256,address,uint256)'](
        firstTokenId,
        other,
        firstTokenIdAmount,
        { from: deployer },
      );

      const receipt = await this.token.burn(
        other,
        firstTokenId,
        firstTokenIdAmount.subn(1),
        { from: other },
      );
      expectEvent(receipt, 'TransferSingle', {
        operator: other,
        from: other,
        to: ZERO_ADDRESS,
        amount: firstTokenIdAmount.subn(1),
        id: firstTokenId,
      });

      expect(
        await this.token.balanceOf(other, firstTokenId),
      ).to.be.bignumber.equal('1');
    });
  });
});
