/* eslint-disable @typescript-eslint/no-explicit-any */
import {
  loadFixture,
  setBalance,
} from "@nomicfoundation/hardhat-network-helpers";
import { expect } from "chai";
import { ethers } from "hardhat";
import { parseEther } from "ethers/lib/utils";

import { cnV3PublicDelegationFixture } from "../common/fixtures";
import { addTime, augmentChai, DIRTY, toPeb } from "../common/helper";
import { smock } from "@defi-wonderland/smock";
import {
  CnStakingV3MultiSig,
  CnStakingV3MultiSig__factory,
} from "../../typechain-types";

const ONE_WEEK = 60 * 60 * 24 * 7;

describe("PublicDelegation tests", function () {
  /**
   * @dev PD's initial setup was already tested in `cnStakingV3.test`.
   */
  describe("Check ERC20 initial setup", function () {
    it("name, symbol, decimals", async function () {
      const { pd1, pd2 } = await loadFixture(cnV3PublicDelegationFixture);
      expect(await pd1.name()).to.be.eq("GC1 Public Delegated KAIA");
      expect(await pd1.symbol()).to.be.eq("GC1-pdKAIA");
      expect(await pd1.decimals()).to.be.eq(18);

      expect(await pd2.name()).to.be.eq("GC2 Public Delegated KAIA");
      expect(await pd2.symbol()).to.be.eq("GC2-pdKAIA");
      expect(await pd2.decimals()).to.be.eq(18);
    });
  });
  describe("Operations", function () {
    it("#updateCommissionTo: only admin can update commissionTo", async function () {
      const { pd1, user1 } = await loadFixture(cnV3PublicDelegationFixture);
      await expect(
        pd1.connect(user1).updateCommissionTo(user1.address)
      ).to.be.revertedWithCustomError(pd1, "OwnableUnauthorizedAccount");

      await pd1.updateCommissionTo(user1.address);
      expect(await pd1.commissionTo()).to.be.eq(user1.address);
    });
    it("#updateCommissionTo: clean up reward to previous commissionTo address", async function () {
      const { pd1, commissionTo, user1 } = await loadFixture(
        cnV3PublicDelegationFixture
      );
      const before = await ethers.provider.getBalance(commissionTo[0]);

      await pd1.updateCommissionRate(1000); // 10%

      // Get reward: 1000 KAIA.
      await setBalance(pd1.address, parseEther("1000"));

      await expect(pd1.updateCommissionTo(user1.address))
        .to.emit(pd1, "UpdateCommissionTo")
        .withArgs(commissionTo[0], user1.address)
        .to.emit(pd1, "SendCommission")
        .withArgs(commissionTo[0], toPeb(100n));

      expect(await pd1.commissionTo()).to.be.eq(user1.address);
      // Previous commissionTo address receives 100 KAIA.
      expect(await ethers.provider.getBalance(commissionTo[0])).to.be.eq(
        before.add(toPeb(100n))
      );
    });
    it("#updateCommissionRate: only admin can update commissionRate", async function () {
      const { pd1, user1 } = await loadFixture(cnV3PublicDelegationFixture);
      await expect(
        pd1.connect(user1).updateCommissionRate(1)
      ).to.be.revertedWithCustomError(pd1, "OwnableUnauthorizedAccount");

      await expect(pd1.updateCommissionRate(1))
        .to.emit(pd1, "UpdateCommissionRate")
        .withArgs(0, 1);
      expect(await pd1.commissionRate()).to.be.eq(1);
    });
    it("#updateCommissionRate: updated commissionRate isn't applied to existing rewards", async function () {
      const { pd1, user1 } = await loadFixture(cnV3PublicDelegationFixture);

      await expect(pd1.connect(user1).stake({ value: toPeb(1000n) }))
        .to.emit(pd1, "Staked")
        .withArgs(user1.address, toPeb(1000n), toPeb(1000n));

      // Get reward: 1000 KAIA.
      await setBalance(pd1.address, parseEther("1000"));

      // set commissionRate to 10% = 1e3.
      await expect(pd1.updateCommissionRate(1000))
        .to.emit(pd1, "UpdateCommissionRate")
        .withArgs(0, 1000);

      // Get reward: 1000 KAIA.
      await setBalance(pd1.address, parseEther("1000"));

      // Reward for user1 = 1000 + 1000 * 0.9 = 1900 KAIA.
      expect(await pd1.previewRedeem(toPeb(1000n))).to.be.eq(toPeb(2900n));
    });
    it("#updateCommissionRate: can't set commissionRate to more than MAX_COMMISSION_RATE", async function () {
      const { pd1 } = await loadFixture(cnV3PublicDelegationFixture);
      await expect(pd1.updateCommissionRate(10001)).to.be.revertedWith(
        "Commission rate is too high."
      );
    });
  });
  describe("Simple ERC20 test", function () {
    it("#transfer: not allowed", async function () {
      const { pd1, user1, user2 } = await loadFixture(
        cnV3PublicDelegationFixture
      );

      // User1 stakes 1000 KAIA.
      await expect(pd1.connect(user1).stake({ value: toPeb(1000n) }))
        .to.emit(pd1, "Staked")
        .withArgs(user1.address, toPeb(1000n), toPeb(1000n));

      // User1 transfers 500 shares to user2.
      await expect(
        pd1.connect(user1).transfer(user2.address, toPeb(500n))
      ).to.be.revertedWith("Transfer not allowed.");
    });
    it("#transferFrom: not allowed", async function () {
      const { pd1, user1, user2, user3 } = await loadFixture(
        cnV3PublicDelegationFixture
      );

      // User1 stakes 1000 KAIA.
      await expect(pd1.connect(user1).stake({ value: toPeb(1000n) }))
        .to.emit(pd1, "Staked")
        .withArgs(user1.address, toPeb(1000n), toPeb(1000n));

      await pd1.connect(user1).approve(user3.address, toPeb(500n));

      // User3 transfers 500 shares from user1 to user2.
      await expect(
        pd1
          .connect(user3)
          .transferFrom(user1.address, user2.address, toPeb(500n))
      ).to.be.revertedWith("Transfer not allowed.");
    });
  });
  describe("Simple delegation tests", function () {
    it("#stake: user can stake", async function () {
      const { cnV3s, pd1, user1, user2, user3 } = await loadFixture(
        cnV3PublicDelegationFixture
      );

      // User1 stakes 1000 KAIA.
      // Shares of user1 = 1000 * 10^18.
      await expect(pd1.connect(user1).stake({ value: toPeb(1000n) })).to.emit(
        pd1,
        "Staked"
      );
      expect(await ethers.provider.getBalance(cnV3s[0].address)).to.be.eq(
        toPeb(1000n)
      );
      expect(await pd1.balanceOf(user1.address)).to.be.eq(toPeb(1000n));
      expect(await pd1.maxRedeem(user1.address)).to.be.eq(toPeb(1000n));

      // User2 stakes 2000 KAIA.
      // Shares of user2 = 2000 * 1000 / 1000 = 2000 * 10^18.
      await expect(pd1.connect(user2).stake({ value: toPeb(2000n) })).to.emit(
        pd1,
        "Staked"
      );
      expect(await ethers.provider.getBalance(cnV3s[0].address)).to.be.eq(
        toPeb(3000n)
      );
      expect(await pd1.balanceOf(user2.address)).to.be.eq(toPeb(2000n));
      expect(await pd1.maxRedeem(user2.address)).to.be.eq(toPeb(2000n));

      // User3 stakes 3000 KAIA.
      // Shares of user3 = 3000 * 1000 / 1000 = 3000 * 10^18.
      await expect(pd1.connect(user3).stake({ value: toPeb(3000n) })).to.emit(
        pd1,
        "Staked"
      );
      expect(await ethers.provider.getBalance(cnV3s[0].address)).to.be.eq(
        toPeb(6000n)
      );
      expect(await pd1.balanceOf(user3.address)).to.be.eq(toPeb(3000n));
      expect(await pd1.maxRedeem(user3.address)).to.be.eq(toPeb(3000n));

      expect(await pd1.totalAssets()).to.be.eq(toPeb(6000n));
      expect(await pd1.totalSupply()).to.be.eq(toPeb(6000n));
    });
    it("#stake: user can't stake too small amounts (shares < 1)", async function () {
      const { cnV3s, pd1, user1, user2 } = await loadFixture(
        cnV3PublicDelegationFixture
      );

      // User1 stakes 1000 KAIA.
      // Shares of user1 = 1000 * 10^18.
      await expect(pd1.connect(user1).stake({ value: toPeb(1000n) })).to.emit(
        pd1,
        "Staked"
      );
      expect(await ethers.provider.getBalance(cnV3s[0].address)).to.be.eq(
        toPeb(1000n)
      );
      expect(await pd1.balanceOf(user1.address)).to.be.eq(toPeb(1000n));
      expect(await pd1.maxRedeem(user1.address)).to.be.eq(toPeb(1000n));

      // PD captures 1000 KAIA.
      await setBalance(pd1.address, parseEther("1000"));

      // User2 stakes 1 peb
      // Shares of user2 = Floor(1 * 1000 / 2000) = 0.
      await expect(pd1.connect(user2).stake({ value: 1 })).to.be.revertedWith(
        "Stake amount is too low."
      );
    });
    it("#stakeFor: user can stake on behalf of other users", async function () {
      const { cnV3s, pd1, user1, user2, user3 } = await loadFixture(
        cnV3PublicDelegationFixture
      );

      // User1 stakes 1000 KAIA on behalf of user2.
      // Shares of user2 = 1000 * 10^18.
      await expect(
        pd1.connect(user1).stakeFor(user2.address, { value: toPeb(1000n) })
      )
        .to.emit(pd1, "Staked")
        .withArgs(user2.address, toPeb(1000n), toPeb(1000n));
      expect(await ethers.provider.getBalance(cnV3s[0].address)).to.be.eq(
        toPeb(1000n)
      );
      expect(await pd1.balanceOf(user2.address)).to.be.eq(toPeb(1000n));
      expect(await pd1.maxRedeem(user2.address)).to.be.eq(toPeb(1000n));

      // User2 stakes 2000 KAIA on behalf of user3.
      // Shares of user3 = 2000 * 1000 / 1000 = 2000 * 10^18.
      await expect(
        pd1.connect(user2).stakeFor(user3.address, { value: toPeb(2000n) })
      )
        .to.emit(pd1, "Staked")
        .withArgs(user3.address, toPeb(2000n), toPeb(2000n));
      expect(await ethers.provider.getBalance(cnV3s[0].address)).to.be.eq(
        toPeb(3000n)
      );
      expect(await pd1.balanceOf(user3.address)).to.be.eq(toPeb(2000n));
      expect(await pd1.maxRedeem(user3.address)).to.be.eq(toPeb(2000n));

      expect(await pd1.totalAssets()).to.be.eq(toPeb(3000n));
      expect(await pd1.totalSupply()).to.be.eq(toPeb(3000n));
    });
    it("#stakeFor: user can't stake for zeor address", async function () {
      const { pd1, user1 } = await loadFixture(cnV3PublicDelegationFixture);

      // User1 stakes 1000 KAIA on behalf of user2.
      // Shares of user2 = 1000 * 10^18.
      await expect(
        pd1
          .connect(user1)
          .stakeFor(ethers.constants.AddressZero, { value: toPeb(1000n) })
      ).to.be.revertedWith("Address is null.");
    });
    it("#fallback: fallback leads to stakes", async function () {
      const { pd1, user1 } = await loadFixture(cnV3PublicDelegationFixture);

      // User1 stakes 1000 KAIA.
      // Shares of user1 = 1000 * 10^18.
      await expect(
        user1.sendTransaction({ to: pd1.address, value: toPeb(1000n) })
      )
        .to.emit(pd1, "Staked")
        .withArgs(user1.address, toPeb(1000n), toPeb(1000n));

      expect(await pd1.balanceOf(user1.address)).to.be.eq(toPeb(1000n));
      expect(await pd1.maxRedeem(user1.address)).to.be.eq(toPeb(1000n));
      expect(await pd1.totalAssets()).to.be.eq(toPeb(1000n));
    });
    it("#redeem: user can withdraw staked amount", async function () {
      const { cnV3s, pd1, user1 } = await loadFixture(
        cnV3PublicDelegationFixture
      );

      // User1 stakes 1000 KAIA.
      // Shares of user1 = 1000 * 10^18.
      await expect(pd1.connect(user1).stake({ value: toPeb(1000n) })).to.emit(
        pd1,
        "Staked"
      );

      await expect(pd1.connect(user1).redeem(user1.address, toPeb(1000n)))
        .to.emit(cnV3s[0], "ApproveStakingWithdrawal")
        .to.emit(pd1, "Redeemed");

      expect(await pd1.balanceOf(user1.address)).to.be.eq(0);
      expect(await pd1.maxRedeem(user1.address)).to.be.eq(0);

      expect(await pd1.totalAssets()).to.be.eq(toPeb(0n));
      expect(await pd1.totalSupply()).to.be.eq(toPeb(0n));
    });
    it("#redeem: user can't withdraw more than staked amount", async function () {
      const { pd1, user1 } = await loadFixture(cnV3PublicDelegationFixture);

      // User1 stakes 1000 KAIA.
      // Shares of user1 = 1000 * 10^18.
      await expect(pd1.connect(user1).stake({ value: toPeb(1000n) })).to.emit(
        pd1,
        "Staked"
      );

      await expect(
        pd1.connect(user1).redeem(user1.address, toPeb(1001n))
      ).to.be.revertedWithCustomError(pd1, "ERC20InsufficientBalance");
    });
    it("#redeem: user can't withdraw to zero address", async function () {
      const { pd1, user1 } = await loadFixture(cnV3PublicDelegationFixture);

      // User1 stakes 1000 KAIA.
      // Shares of user1 = 1000 * 10^18.
      await expect(pd1.connect(user1).stake({ value: toPeb(1000n) })).to.emit(
        pd1,
        "Staked"
      );

      await expect(
        pd1.connect(user1).redeem(ethers.constants.AddressZero, toPeb(1000n))
      ).to.be.revertedWith("Address is null.");
    });
    it("#withdraw: user can withdraw staked amount", async function () {
      const { cnV3s, pd1, user1 } = await loadFixture(
        cnV3PublicDelegationFixture
      );

      // User1 stakes 1000 KAIA.
      // Shares of user1 = 1000 * 10^18.
      await expect(pd1.connect(user1).stake({ value: toPeb(1000n) })).to.emit(
        pd1,
        "Staked"
      );

      await expect(pd1.connect(user1).withdraw(user1.address, toPeb(1000n)))
        .to.emit(cnV3s[0], "ApproveStakingWithdrawal")
        .to.emit(pd1, "Redeemed")
        .withArgs(user1.address, user1.address, toPeb(1000n), toPeb(1000n));

      expect(await pd1.balanceOf(user1.address)).to.be.eq(0);
      expect(await pd1.maxRedeem(user1.address)).to.be.eq(0);

      expect(await pd1.totalAssets()).to.be.eq(toPeb(0n));
      expect(await pd1.totalSupply()).to.be.eq(toPeb(0n));
    });
    it("#withdraw: user can't withdraw more than staked amount", async function () {
      const { pd1, user1 } = await loadFixture(cnV3PublicDelegationFixture);

      // User1 stakes 1000 KAIA.
      // Shares of user1 = 1000 * 10^18.
      await expect(pd1.connect(user1).stake({ value: toPeb(1000n) })).to.emit(
        pd1,
        "Staked"
      );

      await expect(
        pd1.connect(user1).withdraw(user1.address, toPeb(1001n))
      ).to.be.revertedWithCustomError(pd1, "ERC20InsufficientBalance");
    });
    it("#withdraw: user can't withdraw to zero address", async function () {
      const { pd1, user1 } = await loadFixture(cnV3PublicDelegationFixture);

      // User1 stakes 1000 KAIA.
      // Shares of user1 = 1000 * 10^18.
      await expect(pd1.connect(user1).stake({ value: toPeb(1000n) })).to.emit(
        pd1,
        "Staked"
      );

      await expect(
        pd1.connect(user1).withdraw(ethers.constants.AddressZero, toPeb(1001n))
      ).to.be.revertedWith("Address is null.");
    });
    it("#withdraw: unstaking amount isn't eligible for reward", async function () {
      const { pd1, user1, user2 } = await loadFixture(
        cnV3PublicDelegationFixture
      );

      // User1 stakes 1000 KAIA.
      await expect(pd1.connect(user1).stake({ value: toPeb(1000n) })).to.emit(
        pd1,
        "Staked"
      );

      // User2 stakes 1500 KAIA.
      await expect(pd1.connect(user2).stake({ value: toPeb(1500n) })).to.emit(
        pd1,
        "Staked"
      );

      // User1 withdraws 500 KAIA.
      await expect(pd1.connect(user1).withdraw(user1.address, toPeb(500n)))
        .to.emit(pd1, "Redeemed")
        .withArgs(user1.address, user1.address, toPeb(500n), toPeb(500n));

      // Get reward: 1000 KAIA.
      await setBalance(pd1.address, parseEther("1000"));

      // User1 cancels the withdrawal.
      await expect(
        pd1.connect(user1).cancelApprovedStakingWithdrawal(0)
      ).to.emit(pd1, "RequestCancelWithdrawal");

      // Reward for user1 = 500 * 1000 / 2000 = 250 KAIA.
      expect(
        await pd1.previewRedeem(await pd1.maxRedeem(user1.address))
      ).to.be.closeTo(toPeb(1250n), DIRTY);

      // Reward for user2 = 1500 * 1000 / 2000 = 750 KAIA.
      expect(
        await pd1.previewRedeem(await pd1.maxRedeem(user2.address))
      ).to.be.closeTo(toPeb(2250n), DIRTY);
    });
    it("#cancelApprovedStakingWithdrawal: user can cancel the withdrawal if it's in Unknown state", async function () {
      const { cnV3s, pd1, user1 } = await loadFixture(
        cnV3PublicDelegationFixture
      );

      // User1 stakes 1000 KAIA.
      // Shares of user1 = 1000 * 10^18.
      await expect(pd1.connect(user1).stake({ value: toPeb(1000n) })).to.emit(
        pd1,
        "Staked"
      );

      await expect(pd1.connect(user1).redeem(user1.address, toPeb(1000n)))
        .to.emit(cnV3s[0], "ApproveStakingWithdrawal")
        .to.emit(pd1, "Redeemed");

      expect(await pd1.balanceOf(user1.address)).to.be.eq(toPeb(0n));
      expect(await pd1.maxRedeem(user1.address)).to.be.eq(toPeb(0n));
      expect(await pd1.totalAssets()).to.be.eq(toPeb(0n));
      expect(await pd1.totalSupply()).to.be.eq(toPeb(0n));

      await addTime(ONE_WEEK / 2);

      await expect(pd1.connect(user1).cancelApprovedStakingWithdrawal(0))
        .to.emit(cnV3s[0], "CancelApprovedStakingWithdrawal")
        .to.emit(pd1, "RequestCancelWithdrawal")
        .withArgs(user1.address, 0);

      // Shares are restored.
      expect(await pd1.balanceOf(user1.address)).to.be.eq(toPeb(1000n));
      expect(await pd1.maxRedeem(user1.address)).to.be.eq(toPeb(1000n));

      expect(await pd1.totalAssets()).to.be.eq(toPeb(1000n));
      expect(await pd1.totalSupply()).to.be.eq(toPeb(1000n));
    });
    it("#cancelApprovedStakingWithdrawal: check the shares are restored correctly", async function () {
      const { cnV3s, pd1, user1, user2 } = await loadFixture(
        cnV3PublicDelegationFixture
      );

      // Get reward: 1000 KAIA.
      await setBalance(pd1.address, parseEther("1000"));

      // User1 stakes 1000 KAIA.
      // Shares of user1 = 1000 * 10^18.
      await expect(pd1.connect(user1).stake({ value: toPeb(1000n) })).to.emit(
        pd1,
        "Staked"
      );

      // Shares of user2 = 2000 * 1000 / 2000 = 1000 * 10^18.
      await expect(pd1.connect(user2).stake({ value: toPeb(2000n) })).to.emit(
        pd1,
        "Staked"
      );

      await expect(pd1.connect(user1).redeem(user1.address, toPeb(500n)))
        .to.emit(cnV3s[0], "ApproveStakingWithdrawal")
        .to.emit(pd1, "Redeemed");

      // Get reward: 1000 KAIA.
      await setBalance(pd1.address, parseEther("1000"));

      await addTime(ONE_WEEK / 2);

      await expect(pd1.connect(user1).cancelApprovedStakingWithdrawal(0))
        .to.emit(cnV3s[0], "CancelApprovedStakingWithdrawal")
        .to.emit(pd1, "RequestCancelWithdrawal")
        .withArgs(user1.address, 0);

      // Shares are restored.
      expect(await pd1.balanceOf(user1.address)).to.be.eq(toPeb(875n));
      expect(await pd1.maxWithdraw(user1.address)).to.be.closeTo(
        parseEther("2333.3333"),
        DIRTY
      );

      expect(await pd1.totalAssets()).to.be.eq(toPeb(5000n));
      expect(await pd1.totalSupply()).to.be.eq(toPeb(1875n));
    });
    it("#cancelApprovedStakingWithdrawal: user can't cancel empty approved withdrawal", async function () {
      const { cnV3s, pd1, user1 } = await loadFixture(
        cnV3PublicDelegationFixture
      );

      // User1 stakes 1000 KAIA.
      // Shares of user1 = 1000 * 10^18.
      await expect(pd1.connect(user1).stake({ value: toPeb(1000n) })).to.emit(
        pd1,
        "Staked"
      );

      await expect(pd1.connect(user1).redeem(user1.address, toPeb(1000n)))
        .to.emit(cnV3s[0], "ApproveStakingWithdrawal")
        .to.emit(pd1, "Redeemed");

      await expect(
        pd1.connect(user1).cancelApprovedStakingWithdrawal(1)
      ).to.be.revertedWith("Not the owner of the request.");
    });
    it("#cancelApprovedStakingWithdrawal: user can't cancel other users withdrawal", async function () {
      const { cnV3s, pd1, user1, user2 } = await loadFixture(
        cnV3PublicDelegationFixture
      );

      // User1 stakes 1000 KAIA.
      // Shares of user1 = 1000 * 10^18.
      await expect(pd1.connect(user1).stake({ value: toPeb(1000n) })).to.emit(
        pd1,
        "Staked"
      );

      await expect(pd1.connect(user1).redeem(user1.address, toPeb(1000n)))
        .to.emit(cnV3s[0], "ApproveStakingWithdrawal")
        .to.emit(pd1, "Redeemed");

      await expect(
        pd1.connect(user2).cancelApprovedStakingWithdrawal(0)
      ).to.be.revertedWith("Not the owner of the request.");
    });
    it("#cancelApprovedStakingWithdrawal: user can't cancel already withdrawn request", async function () {
      const { cnV3s, pd1, user1 } = await loadFixture(
        cnV3PublicDelegationFixture
      );

      // User1 stakes 1000 KAIA.
      // Shares of user1 = 1000 * 10^18.
      await expect(pd1.connect(user1).stake({ value: toPeb(1000n) })).to.emit(
        pd1,
        "Staked"
      );

      await expect(pd1.connect(user1).redeem(user1.address, toPeb(1000n)))
        .to.emit(cnV3s[0], "ApproveStakingWithdrawal")
        .to.emit(pd1, "Redeemed");

      await addTime(ONE_WEEK);

      await expect(pd1.connect(user1).claim(0))
        .to.emit(cnV3s[0], "WithdrawApprovedStaking")
        .to.emit(pd1, "Claimed")
        .withArgs(user1.address, 0);

      await expect(
        pd1.connect(user1).cancelApprovedStakingWithdrawal(0)
      ).to.be.revertedWith("Invalid state.");
    });
    it("#claim: user can withdraw after the lockup period", async function () {
      const { cnV3s, pd1, user1 } = await loadFixture(
        cnV3PublicDelegationFixture
      );

      // User1 stakes 1000 KAIA.
      // Shares of user1 = 1000 * 10^18.
      await expect(pd1.connect(user1).stake({ value: toPeb(1000n) })).to.emit(
        pd1,
        "Staked"
      );

      await expect(pd1.connect(user1).redeem(user1.address, toPeb(1000n)))
        .to.emit(cnV3s[0], "ApproveStakingWithdrawal")
        .to.emit(pd1, "Redeemed");

      await addTime(ONE_WEEK);

      await expect(pd1.connect(user1).claim(0))
        .to.emit(cnV3s[0], "WithdrawApprovedStaking")
        .to.emit(pd1, "Claimed")
        .withArgs(user1.address, 0);

      expect(await ethers.provider.getBalance(cnV3s[0].address)).to.be.eq(
        toPeb(0n)
      );
      expect(await pd1.balanceOf(user1.address)).to.be.eq(0);
      expect(await pd1.maxRedeem(user1.address)).to.be.eq(0);
      expect(await pd1.totalAssets()).to.be.eq(toPeb(0n));
      expect(await pd1.totalSupply()).to.be.eq(toPeb(0n));
    });
    it("#claim: user can't claim other users claimable withdrawal", async function () {
      const { cnV3s, pd1, user1, user2 } = await loadFixture(
        cnV3PublicDelegationFixture
      );

      // User1 stakes 1000 KAIA.
      // Shares of user1 = 1000 * 10^18.
      await expect(pd1.connect(user1).stake({ value: toPeb(1000n) })).to.emit(
        pd1,
        "Staked"
      );

      await expect(pd1.connect(user1).redeem(user1.address, toPeb(1000n)))
        .to.emit(cnV3s[0], "ApproveStakingWithdrawal")
        .to.emit(pd1, "Redeemed");

      await addTime(ONE_WEEK * 2);

      await expect(pd1.connect(user2).claim(0)).to.be.revertedWith(
        "Not the owner of the request."
      );
    });
    it("#claim: withdrawal is canceled after withdraw period", async function () {
      const { cnV3s, pd1, user1 } = await loadFixture(
        cnV3PublicDelegationFixture
      );

      // User1 stakes 1000 KAIA.
      // Shares of user1 = 1000 * 10^18.
      await expect(pd1.connect(user1).stake({ value: toPeb(1000n) })).to.emit(
        pd1,
        "Staked"
      );

      await expect(pd1.connect(user1).redeem(user1.address, toPeb(1000n)))
        .to.emit(cnV3s[0], "ApproveStakingWithdrawal")
        .to.emit(pd1, "Redeemed");

      await addTime(ONE_WEEK * 2);

      await expect(pd1.connect(user1).claim(0))
        .to.emit(cnV3s[0], "CancelApprovedStakingWithdrawal")
        .to.not.emit(pd1, "Claimed");

      // Shares are restored.
      expect(await ethers.provider.getBalance(cnV3s[0].address)).to.be.eq(
        toPeb(1000n)
      );
      expect(await pd1.balanceOf(user1.address)).to.be.eq(toPeb(1000n));
      expect(await pd1.maxRedeem(user1.address)).to.be.eq(toPeb(1000n));
      expect(await pd1.totalAssets()).to.be.eq(toPeb(1000n));
      expect(await pd1.totalSupply()).to.be.eq(toPeb(1000n));
    });
  });
  describe("Others about commission and fee", function () {
    it("#others: send commission to the commissionTo", async function () {
      const { pd1, user1, commissionTo } = await loadFixture(
        cnV3PublicDelegationFixture
      );

      const before = await ethers.provider.getBalance(commissionTo[0]);
      // set commissionRate to 10% = 1e3.
      await pd1.updateCommissionRate(1000);

      // User1 stakes 1000 KAIA.
      await expect(pd1.connect(user1).stake({ value: toPeb(1000n) })).to.emit(
        pd1,
        "Staked"
      );

      // Get reward: 1000 KAIA.
      await setBalance(pd1.address, parseEther("1000"));

      // commissionTo will receive 100 KAIA.
      expect(await pd1.previewRedeem(toPeb(1000n))).to.be.eq(toPeb(1900n));

      // Manual sweep.
      await expect(pd1.sweep())
        .to.be.emit(pd1, "SendCommission")
        .withArgs(commissionTo[0], toPeb(100n));

      // commissionTo receives 100 KAIA.
      expect(await ethers.provider.getBalance(commissionTo[0])).to.be.eq(
        before.add(toPeb(100n))
      );
    });
  });
  describe("Simple redelegation tests", function () {
    it("#redelegateByShares (redelegateByAssets): user can redelegate staked amount", async function () {
      const { cnV3s, pd1, pd2, user1 } = await loadFixture(
        cnV3PublicDelegationFixture
      );

      // Test process:
      // User1 stakes 1000 KAIA to cnV3s[0].
      // User1 redelegates 1000 KAIA to cnV3s[1].
      // Check if the shares are updated correctly.

      // User1 stakes 1000 KAIA.
      await expect(pd1.connect(user1).stake({ value: toPeb(1000n) })).to.emit(
        pd1,
        "Staked"
      );
      // User1 redelegates 1000 KAIA to cnV3s[1].
      await expect(
        pd1.connect(user1).redelegateByShares(cnV3s[1].address, toPeb(1000n))
      )
        .to.emit(cnV3s[0], "Redelegation")
        .withArgs(user1.address, cnV3s[1].address, toPeb(1000n))
        .to.emit(pd1, "Redelegated")
        .withArgs(user1.address, cnV3s[1].address, toPeb(1000n))
        .to.emit(cnV3s[1], "DelegateKaia")
        .to.emit(cnV3s[1], "HandleRedelegation")
        .withArgs(
          user1.address,
          cnV3s[0].address,
          cnV3s[1].address,
          toPeb(1000n)
        )
        .to.emit(pd2, "Staked");

      // Shares are updated.
      // pd1:
      expect(await ethers.provider.getBalance(cnV3s[0].address)).to.be.eq(
        toPeb(0n)
      );
      expect(await pd1.balanceOf(user1.address)).to.be.eq(toPeb(0n));
      expect(await pd1.maxRedeem(user1.address)).to.be.eq(toPeb(0n));
      // pd2:
      expect(await ethers.provider.getBalance(cnV3s[1].address)).to.be.eq(
        toPeb(1000n)
      );
      expect(await pd2.balanceOf(user1.address)).to.be.eq(toPeb(1000n));
      expect(await pd2.maxRedeem(user1.address)).to.be.eq(toPeb(1000n));

      expect(await pd1.totalAssets()).to.be.eq(toPeb(0n));
      expect(await pd1.totalSupply()).to.be.eq(toPeb(0n));
      expect(await pd2.totalAssets()).to.be.eq(toPeb(1000n));
      expect(await pd2.totalSupply()).to.be.eq(toPeb(1000n));
    });
    it("#redelegateByShares (redelegateByAssets): shares must be calculated correctly", async function () {
      const { cnV3s, pd1, pd2, user1, user2 } = await loadFixture(
        cnV3PublicDelegationFixture
      );

      // Test process:
      // User1 stakes 1000 KAIA to cnV3s[0].
      // User2 stakes 2000 KAIA to cnV3s[0].
      // Reward 500 KAIA.

      // User1 stakes 1000 KAIA.
      await expect(pd1.connect(user1).stake({ value: toPeb(1000n) })).to.emit(
        pd1,
        "Staked"
      );
      // User2 stakes 2000 KAIA.
      await expect(pd1.connect(user2).stake({ value: toPeb(2000n) })).to.emit(
        pd1,
        "Staked"
      );

      // PD1 captures 500 KAIA.
      await setBalance(pd1.address, parseEther("500"));

      // User1 redelegates 500 shares to cnV3s[1] => 583.3333 KAIA.
      await pd1
        .connect(user1)
        .redelegateByShares(cnV3s[1].address, toPeb(500n));

      // Check the shares are calculated correctly.
      // User1 has 583.3333 KAIA in pd1 and pd2.
      expect(await pd1.maxWithdraw(user1.address)).to.be.closeTo(
        parseEther("583.3333"),
        DIRTY
      );
      expect(await pd2.maxWithdraw(user1.address)).to.be.closeTo(
        parseEther("583.3333"),
        DIRTY
      );

      // User2 redelegates 1000 KAIA to cnV3s[1] => 1166.6666 KAIA.
      await pd1
        .connect(user2)
        .redelegateByAssets(cnV3s[1].address, toPeb(1000n));

      // Check the shares are calculated correctly.
      expect(await pd1.maxWithdraw(user2.address)).to.be.closeTo(
        parseEther("1333.33336"),
        DIRTY
      );
      expect(await pd2.maxWithdraw(user2.address)).to.be.closeTo(
        parseEther("1000"),
        DIRTY
      );
    });
    it("#redelegateByShares (redelegateByAssets): user can't redelegate if isRedelegationEnabled is false", async function () {
      const { cnV3s, pd1, user1 } = await loadFixture(
        cnV3PublicDelegationFixture
      );

      // Temporarily disable redelegation.
      await cnV3s[0].toggleRedelegation();

      await expect(pd1.connect(user1).stake({ value: toPeb(1000n) })).to.emit(
        pd1,
        "Staked"
      );
      await expect(
        pd1.connect(user1).redelegateByShares(cnV3s[1].address, toPeb(1000n))
      ).to.be.revertedWith("Redelegation disabled.");
      await expect(
        pd1.connect(user1).redelegateByAssets(cnV3s[1].address, toPeb(1000n))
      ).to.be.revertedWith("Redelegation disabled.");
    });
    it("#redelegateByShares (redelegateByAssets): user can't redelegate to CNv3 not support redelegation", async function () {
      const { cnV3s, pd1, user1 } = await loadFixture(
        cnV3PublicDelegationFixture
      );

      // Temporarily disable redelegation.
      await cnV3s[1].toggleRedelegation();

      await expect(pd1.connect(user1).stake({ value: toPeb(1000n) })).to.emit(
        pd1,
        "Staked"
      );
      await expect(
        pd1.connect(user1).redelegateByShares(cnV3s[1].address, toPeb(1000n))
      ).to.be.revertedWith("Redelegation disabled.");
      await expect(
        pd1.connect(user1).redelegateByAssets(cnV3s[1].address, toPeb(1000n))
      ).to.be.revertedWith("Redelegation disabled.");
    });
    it("#redelegateByShares (redelegateByAssets): user can't redelegate more than staked amount", async function () {
      const { cnV3s, pd1, user1 } = await loadFixture(
        cnV3PublicDelegationFixture
      );

      await expect(pd1.connect(user1).stake({ value: toPeb(1000n) })).to.emit(
        pd1,
        "Staked"
      );
      await expect(
        pd1.connect(user1).redelegateByShares(cnV3s[1].address, toPeb(1001n))
      ).to.be.revertedWithCustomError(pd1, "ERC20InsufficientBalance");
      await expect(
        pd1.connect(user1).redelegateByAssets(cnV3s[1].address, toPeb(1001n))
      ).to.be.revertedWithCustomError(pd1, "ERC20InsufficientBalance");
    });
    it("#redelegateByShares (redelegateByAssets): user can't redelegate to same CnStakingV3", async function () {
      const { cnV3s, pd1, user1 } = await loadFixture(
        cnV3PublicDelegationFixture
      );

      await expect(pd1.connect(user1).stake({ value: toPeb(1000n) })).to.emit(
        pd1,
        "Staked"
      );
      await expect(
        pd1.connect(user1).redelegateByShares(cnV3s[0].address, toPeb(1000n))
      ).to.be.revertedWith("Target can't be self.");
      await expect(
        pd1.connect(user1).redelegateByAssets(cnV3s[0].address, toPeb(1000n))
      ).to.be.revertedWith("Target can't be self.");
    });
    it("#redelegateByShares (redelegateByAssets): user can't redelegate consecutive times", async function () {
      const { cnV3s, pd1, pd2, user1 } = await loadFixture(
        cnV3PublicDelegationFixture
      );

      await expect(pd1.connect(user1).stake({ value: toPeb(1000n) })).to.emit(
        pd1,
        "Staked"
      );
      await expect(
        pd1.connect(user1).redelegateByShares(cnV3s[1].address, toPeb(1000n))
      )
        .to.emit(cnV3s[0], "Redelegation")
        .to.emit(pd1, "Redelegated")
        .withArgs(user1.address, cnV3s[1].address, toPeb(1000n))
        .to.emit(cnV3s[1], "DelegateKaia")
        .to.emit(cnV3s[1], "HandleRedelegation")
        .withArgs(
          user1.address,
          cnV3s[0].address,
          cnV3s[1].address,
          toPeb(1000n)
        )
        .to.emit(pd2, "Staked");

      // User1 can't redelegate again within 7 days.
      await addTime(ONE_WEEK / 2);
      await expect(
        pd2.connect(user1).redelegateByShares(cnV3s[0].address, toPeb(1000n))
      ).to.be.revertedWith("Can't redelegate yet.");
      await expect(
        pd2.connect(user1).redelegateByAssets(cnV3s[0].address, toPeb(1000n))
      ).to.be.revertedWith("Can't redelegate yet.");

      // User1 can redelegate after 7 days.
      await addTime(ONE_WEEK);

      await expect(
        pd2.connect(user1).redelegateByShares(cnV3s[0].address, toPeb(1000n))
      )
        .to.emit(cnV3s[1], "Redelegation")
        .to.emit(pd2, "Redelegated")
        .withArgs(user1.address, cnV3s[0].address, toPeb(1000n))
        .to.emit(cnV3s[0], "DelegateKaia")
        .to.emit(cnV3s[0], "HandleRedelegation")
        .withArgs(
          user1.address,
          cnV3s[1].address,
          cnV3s[0].address,
          toPeb(1000n)
        )
        .to.emit(pd1, "Staked");

      // Shares are updated.
      // pd1:
      expect(await ethers.provider.getBalance(cnV3s[0].address)).to.be.eq(
        toPeb(1000n)
      );
      expect(await pd1.balanceOf(user1.address)).to.be.eq(toPeb(1000n));
      expect(await pd1.maxRedeem(user1.address)).to.be.eq(toPeb(1000n));
      // pd2:
      expect(await ethers.provider.getBalance(cnV3s[1].address)).to.be.eq(
        toPeb(0n)
      );
      expect(await pd2.balanceOf(user1.address)).to.be.eq(toPeb(0n));
      expect(await pd2.maxRedeem(user1.address)).to.be.eq(toPeb(0n));

      expect(await pd1.totalAssets()).to.be.eq(toPeb(1000n));
      expect(await pd1.totalSupply()).to.be.eq(toPeb(1000n));
      expect(await pd2.totalAssets()).to.be.eq(toPeb(0n));
      expect(await pd2.totalSupply()).to.be.eq(toPeb(0n));
    });
    it("#redelegateByShares (redelegateByAssets): user can't redelegate to wrong CnStakingV3", async function () {
      const { pd1, cnV2, user1, user2 } = await loadFixture(
        cnV3PublicDelegationFixture
      );

      await expect(pd1.connect(user1).stake({ value: toPeb(1000n) })).to.emit(
        pd1,
        "Staked"
      );

      await expect(
        pd1.connect(user1).redelegateByShares(user2.address, toPeb(1000n))
      ).to.be.reverted;
      await expect(
        pd1.connect(user1).redelegateByShares(cnV2.address, toPeb(1000n))
      ).to.be.revertedWith("Invalid CnStakingV3.");
      await expect(
        pd1.connect(user1).redelegateByAssets(cnV2.address, toPeb(1000n))
      ).to.be.revertedWith("Invalid CnStakingV3.");
    });
    it("#redelegateByShares (redelegateByAssets): user can't redelegate to CN not registered at address book", async function () {
      const { pd1 } = await loadFixture(cnV3PublicDelegationFixture);

      const fakeCnV3 = await smock.fake<CnStakingV3MultiSig>(
        CnStakingV3MultiSig__factory.abi
      );

      await expect(
        pd1.redelegateByShares(fakeCnV3.address, toPeb(1000n))
      ).to.be.revertedWith("Invalid CnStakingV3.");
      await expect(
        pd1.redelegateByAssets(fakeCnV3.address, toPeb(1000n))
      ).to.be.revertedWith("Invalid CnStakingV3.");
    });
  });
  describe("Simple reward tests", function () {
    it("Reward successfully captured", async function () {
      const { pd1, user1 } = await loadFixture(cnV3PublicDelegationFixture);

      // User1 stakes 1000 KAIA.
      await expect(pd1.connect(user1).stake({ value: toPeb(1000n) })).to.emit(
        pd1,
        "Staked"
      );

      // PD captures 1000 KAIA.
      await setBalance(pd1.address, BigInt(toPeb(1000n)));

      expect(await pd1.totalAssets()).to.be.eq(toPeb(2000n));
      expect(await pd1.previewRedeem(toPeb(1000n))).to.be.eq(toPeb(2000n));
      expect(await pd1.convertToAssets(toPeb(1000n))).to.be.eq(toPeb(2000n));
    });
    it("Reward automatically compounded", async function () {
      const { pd1, user1, user2 } = await loadFixture(
        cnV3PublicDelegationFixture
      );

      // User1 stakes 1000 KAIA.
      await expect(pd1.connect(user1).stake({ value: toPeb(1000n) })).to.emit(
        pd1,
        "Staked"
      );

      // PD captures 1000 KAIA.
      await setBalance(pd1.address, BigInt(toPeb(1000n)));

      // User1's asset = 2000 KAIA
      expect(await pd1.totalAssets()).to.be.eq(toPeb(2000n));
      expect(await pd1.previewRedeem(toPeb(1000n))).to.be.eq(toPeb(2000n));

      // User2 stakes 1000 KAIA.
      await expect(pd1.connect(user2).stake({ value: toPeb(1000n) })).to.emit(
        pd1,
        "Staked"
      );

      // User2's share = 1000 * 1000 / 2000 = 500 * 10^18.
      expect(await ethers.provider.getBalance(pd1.address)).to.be.eq(toPeb(0n));
      expect(await pd1.balanceOf(user2.address)).to.be.eq(toPeb(500n));
      expect(await pd1.maxRedeem(user2.address)).to.be.eq(toPeb(500n));
      expect(await pd1.previewRedeem(toPeb(500n))).to.be.eq(toPeb(1000n));
    });
  });
  describe("Scenario tests", function () {
    let fixture: any;
    let balanceA: any;
    let balanceB: any;
    let balanceC: any;
    let balanceCommissionTo1: any;
    let balanceCommissionTo2: any;
    async function assetsForShares(pd: any, shares: string) {
      return (await pd.convertToAssets(shares)).toBigInt();
    }
    async function checkPSStates(pd: any, users: any[], expectedShares: any[]) {
      for (let i = 0; i < users.length; i++) {
        expect(await pd.maxRedeem(users[i].address)).to.be.closeTo(
          expectedShares[i],
          DIRTY
        );
      }

      let totalAssets = 0n;
      let totalSupply = 0n;
      for (let i = 0; i < expectedShares.length; i++) {
        totalAssets += (await pd.convertToAssets(expectedShares[i])).toBigInt();
        totalSupply += BigInt(expectedShares[i]);
      }
      expect(await pd.totalAssets()).to.be.closeTo(totalAssets, DIRTY);
      expect(await pd.totalSupply()).to.be.closeTo(totalSupply, DIRTY);
    }
    this.beforeEach(async function () {
      fixture = await loadFixture(cnV3PublicDelegationFixture);
      const { pd1, pd2, user1, user2, user3, commissionTo } = fixture;

      // Set the commission rate to 10%.
      await pd1.updateCommissionRate(1000);
      await pd2.updateCommissionRate(1000);

      balanceA = await ethers.provider.getBalance(user1.address);
      balanceB = await ethers.provider.getBalance(user2.address);
      balanceC = await ethers.provider.getBalance(user3.address);
      balanceCommissionTo1 = await ethers.provider.getBalance(commissionTo[0]);
      balanceCommissionTo2 = await ethers.provider.getBalance(commissionTo[1]);

      // User1 stakes for 1000 shares.
      {
        await expect(pd1.connect(user1).stake({ value: toPeb(1000n) })).to.emit(
          pd1,
          "Staked"
        );
        await setBalance(pd1.address, parseEther("6.4"));
      }

      // User2 stakes for 1500 shares.
      {
        await expect(
          pd1
            .connect(user2)
            .stake({ value: await assetsForShares(pd1, toPeb(1500n)) })
        ).to.emit(pd1, "Staked");
        await setBalance(pd1.address, parseEther("6.4"));
      }

      // User3 stakes for 500 shares.
      {
        await expect(
          pd1
            .connect(user3)
            .stake({ value: await assetsForShares(pd1, toPeb(500n)) })
        ).to.emit(pd1, "Staked");
        await setBalance(pd1.address, parseEther("6.4"));
      }
    });
    it("3 users stake and withdraw", async function () {
      const { cnV3s, pd1, user1, user2, user3, commissionTo } = fixture;

      /**
       * @dev Will use shares to ease the calculation.
       */

      // Check the states.
      await checkPSStates(
        pd1,
        [user1, user2, user3],
        [toPeb(1000n), toPeb(1500n), toPeb(500n)]
      );

      // User3 withdraws 500 shares.
      {
        await expect(pd1.connect(user3).redeem(user3.address, toPeb(500n)))
          .to.emit(cnV3s[0], "ApproveStakingWithdrawal")
          .to.emit(pd1, "Redeemed");
        await setBalance(pd1.address, parseEther("6.4"));
      }

      // User2 withdraws 1500 shares.
      {
        await expect(pd1.connect(user2).redeem(user2.address, toPeb(1500n)))
          .to.emit(cnV3s[0], "ApproveStakingWithdrawal")
          .to.emit(pd1, "Redeemed");
        await setBalance(pd1.address, parseEther("6.4"));
      }

      // User1 withdraws 1000 shares.
      {
        await expect(pd1.connect(user1).redeem(user1.address, toPeb(1000n)))
          .to.emit(cnV3s[0], "ApproveStakingWithdrawal")
          .to.emit(pd1, "Redeemed");

        await checkPSStates(pd1, [user1, user2, user3], [0n, 0n, 0n]);

        expect(
          (await cnV3s[0].staking()).toBigInt() -
            (await cnV3s[0].unstaking()).toBigInt()
        ).to.be.eq(0n);
      }

      await addTime(ONE_WEEK);

      // Withdrawal is approved and claim.
      {
        await expect(pd1.connect(user1).claim(2))
          .to.emit(cnV3s[0], "WithdrawApprovedStaking")
          .to.emit(pd1, "Claimed");
        await expect(pd1.connect(user2).claim(1))
          .to.emit(cnV3s[0], "WithdrawApprovedStaking")
          .to.emit(pd1, "Claimed");
        await expect(pd1.connect(user3).claim(0))
          .to.emit(cnV3s[0], "WithdrawApprovedStaking")
          .to.emit(pd1, "Claimed");

        await checkPSStates(pd1, [user1, user2, user3], [0n, 0n, 0n]);
        expect(await ethers.provider.getBalance(cnV3s[0].address)).to.be.eq(
          toPeb(0n)
        );
      }

      // CommissionTo receives 1.6 KAIA.
      expect(await ethers.provider.getBalance(commissionTo[0])).to.be.closeTo(
        balanceCommissionTo1.toBigInt() + parseEther("3.2").toBigInt(),
        DIRTY
      );

      // User1's balance is increased by approx. (6.4 * 2 + 2.56 * 2 + 2.13333) * 0.9 = 18.048 KAIA.
      expect(await ethers.provider.getBalance(user1.address)).to.be.closeTo(
        balanceA.toBigInt() + parseEther("18.048").toBigInt(),
        DIRTY
      );

      // User2's balance is increased by approx. (3.84 * 2 + 3.2) * 0.9 = 9.792 KAIA.
      expect(await ethers.provider.getBalance(user2.address)).to.be.closeTo(
        balanceB.toBigInt() + parseEther("9.792").toBigInt(),
        DIRTY
      );

      // User3's balance is increased by approx. (1.066666) * 0.9 = 0.959999994 KAIA.
      expect(await ethers.provider.getBalance(user3.address)).to.be.closeTo(
        balanceC.toBigInt() + parseEther("0.959999994").toBigInt(),
        DIRTY
      );

      // Check CnV3, PD states.
      expect(await ethers.provider.getBalance(cnV3s[0].address)).to.be.eq(
        toPeb(0n)
      );
      expect(await cnV3s[0].staking()).to.be.eq(toPeb(0n));
      expect(await cnV3s[0].unstaking()).to.be.eq(toPeb(0n));
      expect(await ethers.provider.getBalance(pd1.address)).to.be.eq(toPeb(0n));
      expect(await pd1.totalAssets()).to.be.eq(toPeb(0n));
      expect(await pd1.totalSupply()).to.be.eq(toPeb(0n));
    });
    it("3 users stake and redelegate", async function () {
      const { cnV3s, pd1, pd2, user1, user2, user3, commissionTo } = fixture;

      /**
       * @dev Actually moving exact amount of KAIA is more natural than moving shares.
       *      But we use shares to ease the calculation.
       */

      // User1 redelegates for 1000 shares in PD2.
      {
        await pd1
          .connect(user1)
          .redelegateByAssets(cnV3s[1].address, toPeb(1000n));
        // A's shares in PD1: 9.885305113744376148
        // A's shares in PD2: 1000
        await setBalance(pd1.address, parseEther("6.4"));
        await setBalance(pd2.address, parseEther("6.4"));

        await checkPSStates(
          pd1,
          [user1, user2, user3],
          [
            (await pd1.maxRedeem(user1.address)).toBigInt(),
            toPeb(1500n),
            toPeb(500n),
          ]
        );
        await checkPSStates(pd2, [user1, user2, user3], [toPeb(1000n), 0n, 0n]);
      }

      // User2 redelegates for 500 shares in PD2.
      {
        await pd1
          .connect(user2)
          .redelegateByAssets(
            cnV3s[1].address,
            assetsForShares(pd2, toPeb(500n))
          );
        // B's shares in PD1: 1003.499943884203771517
        await setBalance(pd1.address, parseEther("6.4"));
        await setBalance(pd2.address, parseEther("6.4"));

        await checkPSStates(
          pd1,
          [user1, user2, user3],
          [
            (await pd1.maxRedeem(user1.address)).toBigInt(),
            (await pd1.maxRedeem(user2.address)).toBigInt(),
            toPeb(500n),
          ]
        );
        await checkPSStates(
          pd2,
          [user1, user2, user3],
          [toPeb(1000n), toPeb(500n), 0n]
        );
      }

      // User3 withdraws all shares.
      {
        await pd1
          .connect(user3)
          .redeem(user3.address, await pd1.maxRedeem(user3.address));
        await setBalance(pd1.address, parseEther("6.4"));
        await checkPSStates(
          pd1,
          [user1, user2, user3],
          [
            (await pd1.maxRedeem(user1.address)).toBigInt(),
            (await pd1.maxRedeem(user2.address)).toBigInt(),
            0n,
          ]
        );
        await checkPSStates(
          pd2,
          [user1, user2, user3],
          [toPeb(1000n), toPeb(500n), 0n]
        );
      }

      // User2 withdraws all shares.
      {
        await pd1
          .connect(user2)
          .redeem(user2.address, await pd1.maxRedeem(user2.address));
        await pd2
          .connect(user2)
          .redeem(user2.address, await pd2.maxRedeem(user2.address));

        await setBalance(pd1.address, parseEther("6.4"));
        /**
         * @dev Adds 12.8 KAIA since the previous `setBalance` will be overwritten by the next `setBalance`.
         */
        await setBalance(pd2.address, parseEther("12.8"));

        await checkPSStates(
          pd1,
          [user1, user2, user3],
          [(await pd1.maxRedeem(user1.address)).toBigInt(), 0n, 0n]
        );
        await checkPSStates(pd2, [user1, user2, user3], [toPeb(1000n), 0n, 0n]);
      }

      // User1 withdraws all shares.
      {
        await pd1
          .connect(user1)
          .redeem(user1.address, await pd1.maxRedeem(user1.address));
        await pd2
          .connect(user1)
          .redeem(user1.address, await pd2.maxRedeem(user1.address));

        await checkPSStates(pd1, [user1, user2, user3], [0n, 0n, 0n]);
        await checkPSStates(pd2, [user1, user2, user3], [0n, 0n, 0n]);
      }

      await addTime(ONE_WEEK);

      // Withdraws are approved and claim.
      {
        await expect(pd1.connect(user1).claim(2))
          .to.emit(cnV3s[0], "WithdrawApprovedStaking")
          .to.emit(pd1, "Claimed");
        await expect(pd2.connect(user1).claim(1))
          .to.emit(cnV3s[1], "WithdrawApprovedStaking")
          .to.emit(pd2, "Claimed");

        await expect(pd1.connect(user2).claim(1))
          .to.emit(cnV3s[0], "WithdrawApprovedStaking")
          .to.emit(pd1, "Claimed");
        await expect(pd2.connect(user2).claim(0))
          .to.emit(cnV3s[1], "WithdrawApprovedStaking")
          .to.emit(pd2, "Claimed");

        await expect(pd1.connect(user3).claim(0))
          .to.emit(cnV3s[0], "WithdrawApprovedStaking")
          .to.emit(pd1, "Claimed");
      }

      // CommissionTo[0] receives 2.24 KAIA.
      // CommissionTo[1] receive 1.28 KAIA.
      const commission = (
        await ethers.provider.getBalance(commissionTo[0])
      ).toBigInt();
      expect(commission).to.be.eq(
        balanceCommissionTo1.toBigInt() + parseEther("4.48").toBigInt()
      );

      const commission2 = (
        await ethers.provider.getBalance(commissionTo[1])
      ).toBigInt();
      expect(commission2).to.be.eq(
        balanceCommissionTo2.toBigInt() + parseEther("2.56").toBigInt()
      );

      // User1's balance is rewarded by approx. 41.09571196 * 0.9 = 36.98614076 KAIA
      expect(await ethers.provider.getBalance(user1.address)).to.be.closeTo(
        balanceA.toBigInt() + parseEther("36.98614076").toBigInt(),
        DIRTY.toBigInt() * 2n
      );

      // User2's balance is rewarded by approx. 24.53102582 * 0.9 = 22.07792324 KAIA
      expect(await ethers.provider.getBalance(user2.address)).to.be.closeTo(
        balanceB.toBigInt() + parseEther("22.07792324").toBigInt(),
        DIRTY.toBigInt() * 2n
      );

      // User3's balance is rewarded by approx. 4.77326222 * 0.9 = 4.29593600 KAIA
      expect(await ethers.provider.getBalance(user3.address)).to.be.closeTo(
        balanceC.toBigInt() + parseEther("4.29593600").toBigInt(),
        DIRTY.toBigInt() * 2n
      );

      // Check CnV3, PD states.
      for (let i = 0; i < cnV3s.length; i++) {
        expect(await ethers.provider.getBalance(cnV3s[i].address)).to.be.eq(
          toPeb(0n)
        );
        expect(await cnV3s[i].staking()).to.be.eq(toPeb(0n));
        expect(await cnV3s[i].unstaking()).to.be.eq(toPeb(0n));
      }
      expect(await ethers.provider.getBalance(pd1.address)).to.be.eq(toPeb(0n));
      expect(await ethers.provider.getBalance(pd2.address)).to.be.eq(toPeb(0n));
      expect(await pd1.totalAssets()).to.be.eq(toPeb(0n));
      expect(await pd1.totalSupply()).to.be.eq(toPeb(0n));
      expect(await pd2.totalAssets()).to.be.eq(toPeb(0n));
      expect(await pd2.totalSupply()).to.be.eq(toPeb(0n));
    });
  });
  describe("Getter tests", function () {
    let fixture: any;
    this.beforeEach(async function () {
      augmentChai();
      fixture = await loadFixture(cnV3PublicDelegationFixture);
      const { pd1, user1, user2, user3 } = fixture;

      await pd1.connect(user1).stake({ value: toPeb(1000n) });
      await pd1.connect(user2).stake({ value: toPeb(1500n) });
      await pd1.connect(user3).stake({ value: toPeb(500n) });

      await setBalance(pd1.address, parseEther("6.4"));
      await pd1.sweep();

      await setBalance(pd1.address, parseEther("6.4"));
      await pd1.connect(user3).withdraw(user3.address, toPeb(360n));

      await addTime(ONE_WEEK);

      await setBalance(pd1.address, parseEther("6.4"));
      await pd1.connect(user3).claim(0);

      await setBalance(pd1.address, parseEther("6.4"));
      await pd1.connect(user2).withdraw(user2.address, toPeb(700n));

      // Final states:
      // Total assets = 1965.6 KAIA
      // Total supply = 1947.850515670952395828 pKLAY
      // KAIA for users: 1009.112344189786891837 KAIA, 813.668516284680337756 KAIA, 142.819139525532770406 KAIA
      // pKLAY for users: 1000 pKLAY, 806.321041427723505759 pKLAY, 141.529474243228890069 pKLAY
    });
    it("Withdrawal requests", async function () {
      const { pd1, user1, user2, user3 } = fixture;
      expect(await pd1.getUserRequestIds(user1.address)).to.equalNumberList([]);
      expect(await pd1.getUserRequestIds(user2.address)).to.equalNumberList([
        1,
      ]);
      expect(await pd1.getUserRequestIds(user3.address)).to.equalNumberList([
        0,
      ]);

      expect(await pd1.getUserRequestCount(user1.address)).to.equal(0);
      expect(await pd1.getUserRequestCount(user2.address)).to.equal(1);
      expect(await pd1.getUserRequestCount(user3.address)).to.equal(1);
    });
    it("Withdrawal requests with given states", async function () {
      const { pd1, user1, user2, user3 } = fixture;
      expect(
        await pd1.getUserRequestIdsWithState(user1.address, 0)
      ).to.equalNumberList([]);
      expect(
        await pd1.getUserRequestIdsWithState(user2.address, 1)
      ).to.equalNumberList([1]);
      expect(
        await pd1.getUserRequestIdsWithState(user3.address, 3)
      ).to.equalNumberList([0]);
    });
    it("Withdrawal states", async function () {
      const { cnV3s, pd1, user2 } = fixture;
      // Not exist
      expect(await pd1.getCurrentWithdrawalRequestState(3)).to.equal(0);
      // Withdrawn
      expect(await pd1.getCurrentWithdrawalRequestState(0)).to.equal(3);
      // Requested
      expect(await pd1.getCurrentWithdrawalRequestState(1)).to.equal(1);

      await addTime(ONE_WEEK);
      // Withdrawable
      expect(await pd1.getCurrentWithdrawalRequestState(1)).to.equal(2);

      await addTime(ONE_WEEK);

      // PendingCancel
      expect(await pd1.getCurrentWithdrawalRequestState(1)).to.equal(4);

      // Request claim but will be canceled.
      await pd1.connect(user2).claim(1);
      // Canceled
      expect(await pd1.getCurrentWithdrawalRequestState(1)).to.equal(5);

      expect(
        await cnV3s[0].getApprovedStakingWithdrawalIds(0, 0, 1)
      ).to.be.equalNumberList([0]);
      // If _to > withdrawalRequestCount, then _to = withdrawalRequestCount.
      expect(
        await cnV3s[0].getApprovedStakingWithdrawalIds(0, 10, 2)
      ).to.be.equalNumberList([1]);
    });
    it("Account states", async function () {
      const { pd1, user1, user2, user3 } = fixture;
      expect(await pd1.balanceOf(user1.address)).to.equal(parseEther("1000"));
      expect(await pd1.balanceOf(user2.address)).to.equal(
        parseEther("806.321041427723505759")
      );
      expect(await pd1.balanceOf(user3.address)).to.equal(
        parseEther("141.529474243228890069")
      );

      expect(await pd1.maxWithdraw(user1.address)).to.equal(
        parseEther("1009.112344189786891837")
      );
      expect(await pd1.maxWithdraw(user2.address)).to.equal(
        parseEther("813.668516284680337756")
      );
      expect(await pd1.maxWithdraw(user3.address)).to.equal(
        parseEther("142.819139525532770406")
      );
    });
    it("Precision test for staking", async function () {
      const { pd1 } = fixture;
      // Currently, 1 pKLAY = 1.009112344189786891(8) KAIA
      // User needs to stake ceiled amount of KAIA.
      expect(
        await pd1.previewDeposit(parseEther("1.009112344189786891"))
      ).to.equal(parseEther("0.999999999999999999"));
      expect(
        await pd1.previewDeposit(parseEther("1.009112344189786892"))
      ).to.equal(parseEther("1"));
    });
    it("Precision test for withdrawal", async function () {
      const { pd1 } = fixture;
      // Currently, 1 pKLAY = 1.009112344189786891(8) KAIA
      // User will receive floored amount of KAIA.
      expect(
        await pd1.previewRedeem(parseEther("0.999999999999999999"))
      ).to.equal(parseEther("1.009112344189786890"));
      expect(await pd1.previewRedeem(parseEther("1"))).to.equal(
        parseEther("1.009112344189786891")
      );
    });
  });
});
