const { BN } = require('@openzeppelin/test-helpers');

const { expect } = require('chai');

const KIP17ERC721Holder = artifacts.require('KIP17ERC721Holder');
const KIP17Mock = artifacts.require('KIP17Mock');
const ERC721Mock = artifacts.require('ERC721Mock');

contract('KIP17ERC721Holder', function (accounts) {
  const [owner] = accounts;

  const name = 'Non Fungible Token';
  const symbol = 'NFT';

  it('receives an KIP17 token', async function () {
    const token = await KIP17Mock.new(name, symbol);
    const tokenId = new BN(1);
    await token.mint(owner, tokenId);

    const receiver = await KIP17ERC721Holder.new();
    await token.safeTransferFrom(owner, receiver.address, tokenId, {
      from: owner,
    });

    expect(await token.ownerOf(tokenId)).to.be.equal(receiver.address);
  });

  it('receives an ERC721 token', async function () {
    const token = await ERC721Mock.new(name, symbol);
    const tokenId = new BN(1);
    await token.mint(owner, tokenId);

    const receiver = await KIP17ERC721Holder.new();
    await token.safeTransferFrom(owner, receiver.address, tokenId, {
      from: owner,
    });

    expect(await token.ownerOf(tokenId)).to.be.equal(receiver.address);
  });
});
