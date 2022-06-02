const {
  shouldBehaveLikeKIP17,
  shouldBehaveLikeKIP17Metadata,
  shouldBehaveLikeKIP17Enumerable,
} = require('./KIP17.behavior');

const KIP17Mock = artifacts.require('KIP17EnumerableMock');

contract('KIP17Enumerable', function (accounts) {
  const name = 'Non Fungible Token';
  const symbol = 'NFT';

  beforeEach(async function () {
    this.token = await KIP17Mock.new(name, symbol);
  });

  shouldBehaveLikeKIP17('KIP17', ...accounts);
  shouldBehaveLikeKIP17Metadata('KIP17', name, symbol, ...accounts);
  shouldBehaveLikeKIP17Enumerable('KIP17', ...accounts);
});
