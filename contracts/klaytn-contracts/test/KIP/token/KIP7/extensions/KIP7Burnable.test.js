const { BN } = require('@openzeppelin/test-helpers');

const { shouldBehaveLikeKIP7Burnable } = require('./KIP7Burnable.behavior');
const KIP7BurnableMock = artifacts.require('KIP7BurnableMock');

contract('KIP7Burnable', function (accounts) {
  const [owner, ...otherAccounts] = accounts;

  const initialBalance = new BN(1000);

  const name = 'My Token';
  const symbol = 'MTKN';

  beforeEach(async function () {
    this.token = await KIP7BurnableMock.new(
      name,
      symbol,
      owner,
      initialBalance,
      { from: owner },
    );
  });

  shouldBehaveLikeKIP7Burnable(owner, initialBalance, otherAccounts);
});
