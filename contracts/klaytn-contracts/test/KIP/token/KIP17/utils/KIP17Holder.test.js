const { BN } = require('@openzeppelin/test-helpers');

const { expect } = require('chai');

const KIP17Holder = artifacts.require('KIP17Holder');
const KIP17Mock = artifacts.require('KIP17Mock');

contract('KIP17Holder', function (accounts) {
  const [owner] = accounts;

  const name = 'Non Fungible Token';
  const symbol = 'NFT';

  it('receives an KIP17 token', async function () {
    const token = await KIP17Mock.new(name, symbol);
    const tokenId = new BN(1);
    await token.mint(owner, tokenId);

    const receiver = await KIP17Holder.new();
    await token.safeTransferFrom(owner, receiver.address, tokenId, {
      from: owner,
    });

    expect(await token.ownerOf(tokenId)).to.be.equal(receiver.address);
  });
});
