const { BN, constants, expectEvent, expectRevert } = require('@openzeppelin/test-helpers');

const { expect } = require('chai');

const KIP37PausableMock = artifacts.require('KIP37PausableMock');
const { ZERO_ADDRESS } = constants;

contract('KIP37Pausable', function (accounts) {
  const [holder, operator, receiver, deployer, other] = accounts;

  const uri = 'https://token.com';

  beforeEach(async function () {
    this.token = await KIP37PausableMock.new(uri);
  });

  context('when token is paused', function () {
    const firstTokenId = new BN('37');
    const firstTokenAmount = new BN('42');

    const secondTokenId = new BN('19842');
    const secondTokenAmount = new BN('23');

    beforeEach(async function () {
      await this.token.setApprovalForAll(operator, true, { from: holder });
      await this.token.mint(holder, firstTokenId, firstTokenAmount, '0x');

      await this.token.pause();
    });

    it('reverts when trying to safeTransferFrom from holder', async function () {
      await expectRevert(
        this.token.safeTransferFrom(
          holder,
          receiver,
          firstTokenId,
          firstTokenAmount,
          '0x',
          { from: holder },
        ),
        'KIP37Pausable: token transfer while paused',
      );
    });

    it('reverts when trying to safeTransferFrom from operator', async function () {
      await expectRevert(
        this.token.safeTransferFrom(
          holder,
          receiver,
          firstTokenId,
          firstTokenAmount,
          '0x',
          { from: operator },
        ),
        'KIP37Pausable: token transfer while paused',
      );
    });

    it('reverts when trying to safeBatchTransferFrom from holder', async function () {
      await expectRevert(
        this.token.safeBatchTransferFrom(
          holder,
          receiver,
          [firstTokenId],
          [firstTokenAmount],
          '0x',
          { from: holder },
        ),
        'KIP37Pausable: token transfer while paused',
      );
    });

    it('reverts when trying to safeBatchTransferFrom from operator', async function () {
      await expectRevert(
        this.token.safeBatchTransferFrom(
          holder,
          receiver,
          [firstTokenId],
          [firstTokenAmount],
          '0x',
          { from: operator },
        ),
        'KIP37Pausable: token transfer while paused',
      );
    });

    it('reverts when trying to mint', async function () {
      await expectRevert(
        this.token.mint(holder, secondTokenId, secondTokenAmount, '0x'),
        'KIP37Pausable: token transfer while paused',
      );
    });

    it('reverts when trying to mintBatch', async function () {
      await expectRevert(
        this.token.mintBatch(
          holder,
          [secondTokenId],
          [secondTokenAmount],
          '0x',
        ),
        'KIP37Pausable: token transfer while paused',
      );
    });

    it('reverts when trying to burn', async function () {
      await expectRevert(
        this.token.burn(holder, firstTokenId, firstTokenAmount),
        'KIP37Pausable: token transfer while paused',
      );
    });

    it('reverts when trying to burnBatch', async function () {
      await expectRevert(
        this.token.burnBatch(
          holder,
          [firstTokenId],
          [firstTokenAmount],
        ),
        'KIP37Pausable: token transfer while paused',
      );
    });

    describe('setApprovalForAll', function () {
      it('approves an operator', async function () {
        await this.token.setApprovalForAll(other, true, {
          from: holder,
        });
        expect(
          await this.token.isApprovedForAll(holder, other),
        ).to.equal(true);
      });
    });

    describe('balanceOf', function () {
      it('returns the amount of tokens owned by the given address', async function () {
        const balance = await this.token.balanceOf(holder, firstTokenId);
        expect(balance).to.be.bignumber.equal(firstTokenAmount);
      });
    });

    describe('isApprovedForAll', function () {
      it('returns the approval of the operator', async function () {
        expect(
          await this.token.isApprovedForAll(holder, operator),
        ).to.equal(true);
      });
    });
  });

  context('when token is paused by id', function () {
    const firstTokenId = new BN('37');
    const firstTokenAmount = new BN('42');

    const secondTokenId = new BN('19842');
    const secondTokenAmount = new BN('23');

    beforeEach(async function () {
      await this.token.setApprovalForAll(operator, true, { from: holder });
      await this.token.mint(holder, firstTokenId, firstTokenAmount, '0x');

      await this.token.pause(firstTokenId);
    });

    it('reverts when trying to safeTransferFrom from holder', async function () {
      await expectRevert(
        this.token.safeTransferFrom(
          holder,
          receiver,
          firstTokenId,
          firstTokenAmount,
          '0x',
          { from: holder },
        ),
        'KIP37Pausable: token transfer while paused',
      );
    });

    it('reverts when trying to safeTransferFrom from operator', async function () {
      await expectRevert(
        this.token.safeTransferFrom(
          holder,
          receiver,
          firstTokenId,
          firstTokenAmount,
          '0x',
          { from: operator },
        ),
        'KIP37Pausable: token transfer while paused',
      );
    });

    it('reverts when trying to safeBatchTransferFrom from holder', async function () {
      await expectRevert(
        this.token.safeBatchTransferFrom(
          holder,
          receiver,
          [firstTokenId],
          [firstTokenAmount],
          '0x',
          { from: holder },
        ),
        'KIP37Pausable: token transfer while paused',
      );
    });

    it('reverts when trying to safeBatchTransferFrom from operator', async function () {
      await expectRevert(
        this.token.safeBatchTransferFrom(
          holder,
          receiver,
          [firstTokenId],
          [firstTokenAmount],
          '0x',
          { from: operator },
        ),
        'KIP37Pausable: token transfer while paused',
      );
    });

    it('other tokenIds can mint', async function () {
      const receipt = await this.token.mint(holder, secondTokenId, secondTokenAmount, '0x', { from: deployer });

      expectEvent(receipt, 'TransferSingle', {
        operator: deployer,
        from: ZERO_ADDRESS,
        to: holder,
        amount: secondTokenAmount,
        id: secondTokenId,
      });

      expect(
        await this.token.balanceOf(holder, secondTokenId),
      ).to.be.bignumber.equal(secondTokenAmount);
    });

    it('reverts when trying to burn', async function () {
      await expectRevert(
        this.token.burn(holder, firstTokenId, firstTokenAmount),
        'KIP37Pausable: token transfer while paused',
      );
    });

    it('reverts when trying to burnBatch', async function () {
      await expectRevert(
        this.token.burnBatch(
          holder,
          [firstTokenId],
          [firstTokenAmount],
        ),
        'KIP37Pausable: token transfer while paused',
      );
    });

    describe('setApprovalForAll', function () {
      it('approves an operator', async function () {
        await this.token.setApprovalForAll(other, true, {
          from: holder,
        });
        expect(
          await this.token.isApprovedForAll(holder, other),
        ).to.equal(true);
      });
    });

    describe('balanceOf', function () {
      it('returns the amount of tokens owned by the given address', async function () {
        const balance = await this.token.balanceOf(holder, firstTokenId);
        expect(balance).to.be.bignumber.equal(firstTokenAmount);
      });
    });

    describe('isApprovedForAll', function () {
      it('returns the approval of the operator', async function () {
        expect(
          await this.token.isApprovedForAll(holder, operator),
        ).to.equal(true);
      });
    });
  });
});
