const { BN, constants, expectEvent, expectRevert } = require('@openzeppelin/test-helpers');
const { expect } = require('chai');
const { ZERO_ADDRESS } = constants;
const KIP37MintableMock = artifacts.require('KIP37MintableMock');

contract('KIP37Mintable', function (accounts) {
  const [deployer, tokenHolder, tokenBatchHolder, other] = accounts;

  const baseUri = 'https://token.com/';

  const DEFAULT_ADMIN_ROLE = '0x0000000000000000000000000000000000000000000000000000000000000000';
  const MINTER_ROLE = web3.utils.soliditySha3('KIP37_MINTER_ROLE');

  beforeEach('create a token type', async function () {
    this.token = await KIP37MintableMock.new(baseUri, { from: deployer });
  });

  it('deployer has the default admin role', async function () {
    expect(await this.token.getRoleMemberCount(DEFAULT_ADMIN_ROLE)).to.be.bignumber.equal('1');
    expect(await this.token.getRoleMember(DEFAULT_ADMIN_ROLE, 0)).to.equal(deployer);
  });

  it('deployer has the minter role', async function () {
    expect(await this.token.getRoleMemberCount(MINTER_ROLE)).to.be.bignumber.equal('1');
    expect(await this.token.getRoleMember(MINTER_ROLE, 0)).to.equal(deployer);
  });

  describe('mint', function () {
    const firstTokenId = new BN(845);
    const firstInitialSupply = new BN(5000);
    const firstTokenMintAmount = new BN('50');
    const secondTokenId = new BN(48324);

    const data = web3.utils.soliditySha3(baseUri);

    beforeEach('create a token type', async function () {
      this.receipt = await this.token.create(firstTokenId, firstInitialSupply, data, { from: deployer });
    });

    it('tokenId exists', async function () {
      expect(await this.token.exists(firstTokenId)).to.equal(true);
    });

    it('tokenId does not exists', async function () {
      expect(await this.token.exists(secondTokenId)).to.equal(false);
    });

    it('reverts with a nonexistent token', async function () {
      await expectRevert(
        this.token.methods['mint(uint256,address,uint256)'](secondTokenId, tokenHolder, firstTokenMintAmount),
        'KIP37: nonexistent token',
      );
    });

    it('reverts with a zero destination address', async function () {
      await expectRevert(
        this.token.methods['mint(uint256,address,uint256)'](
          firstTokenId, ZERO_ADDRESS, firstTokenMintAmount, { from: deployer }),
        'KIP37: mint to the zero address',
      );
    });

    context('with existing tokenId mint tokens', function () {
      beforeEach('create a token type', async function () {
        this.receipt = await this.token.methods['mint(uint256,address,uint256)'](
          firstTokenId, tokenHolder, firstTokenMintAmount, { from: deployer });
      });

      it('emits a TransferSingle event', function () {
        expectEvent(this.receipt, 'TransferSingle', {
          operator: deployer,
          from: ZERO_ADDRESS,
          to: tokenHolder,
          id: firstTokenId,
          amount: firstTokenMintAmount,
        });
      });

      it('credits the minted amount of tokens', async function () {
        expect(await this.token.balanceOf(
          tokenHolder, firstTokenId)).to.be.bignumber.equal(firstTokenMintAmount);
      });

      it('other accounts cannot mint tokens', async function () {
        await expectRevert(
          this.token.methods['mint(uint256,address,uint256)'](
            firstTokenId, other, firstTokenMintAmount, { from: other }),
          'KIP37: must have minter role to mint',
        );
      });
    });
  });

  describe('mint multiple recipients', function () {
    const tokenId = new BN(745);
    const initialSupply = new BN(5000);
    const zeroAddress = [ZERO_ADDRESS, ZERO_ADDRESS, tokenHolder];
    const toList = [accounts[5], accounts[6], accounts[7]];
    const mintAmounts = [new BN(5), new BN(10), new BN(15)];

    const data = web3.utils.soliditySha3(baseUri);

    beforeEach('create a token type', async function () {
      this.receipt = await this.token.create(tokenId, initialSupply, data, { from: deployer });
    });

    it('reverts with a zero destination address', async function () {
      await expectRevert(
        this.token.mint(tokenId, zeroAddress, mintAmounts, { from: deployer }),
        'KIP37: mint to the zero address',
      );
    });

    it('reverts if length of inputs do not match', async function () {
      await expectRevert(
        this.token.mint(tokenId, toList, mintAmounts.slice(1), { from: deployer }),
        'KIP37: toList and amounts length mismatch',
      );

      await expectRevert(
        this.token.mint(tokenId, toList.slice(1), mintAmounts, { from: deployer }),
        'KIP37: toList and amounts length mismatch',
      );
    });

    it('other accounts cannot mint tokens to a list of accounts', async function () {
      await expectRevert(
        this.token.mint(
          tokenId,
          toList,
          mintAmounts,
          { from: other },
        ),
        'KIP37: must have minter role to mint',
      );
    });

    context('deployer can mint tokens to list of accounts', function () {
      beforeEach(async function () {
        (this.receipt = await this.token.mint(
          tokenId,
          toList,
          mintAmounts,
          { from: deployer },
        ));
      });

      it('emits a TransferSingle event', function () {
        expectEvent(this.receipt, 'TransferSingle', {
          operator: deployer,
          from: ZERO_ADDRESS,
          id: tokenId,
        });
      });

      it('credits the minted batch of tokens', async function () {
        for (let i = 0; i < toList.length; i++) {
          expect(await this.token.balanceOf(toList[i], tokenId)).to.be.bignumber.equal(mintAmounts[i]);
        }
      });
    });
  });

  describe('mintBatch', function () {
    const initialSupply = new BN(60600);
    const tokenBatchIds = [new BN(2000), new BN(2010), new BN(2020)];
    const mintAmounts = [new BN(500), new BN(100), new BN(420)];

    beforeEach('create a token type', async function () {
      await Promise.all(tokenBatchIds.map((x, i) => {
        return this.token.create(x, initialSupply, '', { from: deployer });
      }));
    });

    it('reverts with a zero destination address', async function () {
      await expectRevert(
        this.token.mintBatch(ZERO_ADDRESS, tokenBatchIds, mintAmounts, { from: deployer }),
        'KIP37: mint to the zero address',
      );
    });

    it('reverts if length of inputs do not match', async function () {
      await expectRevert(
        this.token.mintBatch(tokenBatchHolder, tokenBatchIds, mintAmounts.slice(1), { from: deployer }),
        'KIP37: ids and amounts length mismatch',
      );

      await expectRevert(
        this.token.mintBatch(tokenBatchHolder, tokenBatchIds.slice(1), mintAmounts, { from: deployer }),
        'KIP37: ids and amounts length mismatch',
      );
    });

    it('other accounts cannot batch mint tokens', async function () {
      await expectRevert(
        this.token.mintBatch(
          other,
          tokenBatchIds,
          mintAmounts,
          { from: other },
        ),
        'KIP37: must have minter role to mint',
      );
    });

    context('with minted batch of tokens', function () {
      beforeEach(async function () {
        (this.receipt = await this.token.mintBatch(
          tokenBatchHolder,
          tokenBatchIds,
          mintAmounts,
          { from: deployer },
        ));
      });

      it('emits a TransferBatch event', function () {
        expectEvent(this.receipt, 'TransferBatch', {
          operator: deployer,
          from: ZERO_ADDRESS,
          to: tokenBatchHolder,
        });
      });

      it('credits the minted batch of tokens', async function () {
        const holderBatchBalances = await this.token.balanceOfBatch(
          new Array(tokenBatchIds.length).fill(tokenBatchHolder),
          tokenBatchIds,
        );

        for (let i = 0; i < holderBatchBalances.length; i++) {
          expect(holderBatchBalances[i]).to.be.bignumber.equal(mintAmounts[i]);
        }
      });
    });
  });

  describe('burning', function () {
    const firstTokenId = new BN(845);
    const firstInitialSupply = new BN(5000);
    const firstTokenMintAmount = new BN('4000');

    const data = web3.utils.soliditySha3(baseUri);

    beforeEach('create a token type', async function () {
      this.receipt = await this.token.create(firstTokenId, firstInitialSupply, data, { from: deployer });
    });

    it('holders can burn their tokens', async function () {
      await this.token.methods['mint(uint256,address,uint256)'](
        firstTokenId,
        tokenHolder,
        firstTokenMintAmount,
        { from: deployer });

      const receipt = await this.token.burn(
        tokenHolder,
        firstTokenId,
        firstTokenMintAmount.subn(3000),
        { from: tokenHolder },
      );
      expectEvent(receipt, 'TransferSingle', {
        operator: tokenHolder,
        from: tokenHolder,
        to: ZERO_ADDRESS,
        amount: firstTokenMintAmount.subn(3000),
        id: firstTokenId,
      });

      expect(
        await this.token.balanceOf(tokenHolder, firstTokenId),
      ).to.be.bignumber.equal('3000');
    });
  });
});
