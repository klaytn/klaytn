const { BN, ether, expectRevert } = require('@openzeppelin/test-helpers');
const { shouldBehaveLikeKIP7Capped } = require('./KIP7Capped.behavior');

const KIP7Capped = artifacts.require('KIP7CappedMock');

contract('KIP7Capped', function (accounts) {
  const [minter, ...otherAccounts] = accounts;

  const cap = ether('1000');

  const name = 'My Token';
  const symbol = 'MTKN';

  it('requires a non-zero cap', async function () {
    await expectRevert(
      KIP7Capped.new(name, symbol, new BN(0), { from: minter }),
      'KIP7Capped: cap is 0',
    );
  });

  context('once deployed', async function () {
    beforeEach(async function () {
      this.token = await KIP7Capped.new(name, symbol, cap, {
        from: minter,
      });
    });

    shouldBehaveLikeKIP7Capped(minter, otherAccounts, cap);
  });
});
