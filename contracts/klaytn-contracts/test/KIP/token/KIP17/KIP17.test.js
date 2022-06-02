const {
  shouldBehaveLikeKIP17,
  shouldBehaveLikeKIP17Metadata,
} = require('./KIP17.behavior');

const KIP17Mock = artifacts.require('KIP17Mock');

contract('KIP17', function (accounts) {
  const name = 'Non Fungible Token';
  const symbol = 'NFT';

  beforeEach(async function () {
    this.token = await KIP17Mock.new(name, symbol);
  });

  shouldBehaveLikeKIP17('KIP17', ...accounts);
  shouldBehaveLikeKIP17Metadata('KIP17', name, symbol, ...accounts);
});
