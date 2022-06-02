const { BN, constants } = require('@openzeppelin/test-helpers');
const KIP17RoyaltyMock = artifacts.require('KIP17RoyaltyMock');
const { ZERO_ADDRESS } = constants;
const { expect } = require('chai');

const { shouldBehaveLikeERC2981 } = require('../../../../token/common/ERC2981.behavior');

contract('KIP17Royalty', function (accounts) {
  const [account1, account2] = accounts;
  const tokenId1 = new BN('1');
  const tokenId2 = new BN('2');
  const royalty = new BN('200');
  const salePrice = new BN('1000');

  beforeEach(async function () {
    this.token = await KIP17RoyaltyMock.new('My Token', 'TKN');

    await this.token.mint(account1, tokenId1);
    await this.token.mint(account1, tokenId2);
    this.account1 = account1;
    this.account2 = account2;
    this.tokenId1 = tokenId1;
    this.tokenId2 = tokenId2;
    this.salePrice = salePrice;
  });

  describe('token specific functions', function () {
    beforeEach(async function () {
      await this.token.setTokenRoyalty(tokenId1, account1, royalty);
    });

    it('removes royalty information after burn', async function () {
      await this.token.burn(tokenId1);
      const tokenInfo = await this.token.royaltyInfo(tokenId1, salePrice);

      expect(tokenInfo[0]).to.be.equal(ZERO_ADDRESS);
      expect(tokenInfo[1]).to.be.bignumber.equal(new BN('0'));
    });
  });
  shouldBehaveLikeERC2981();
});
