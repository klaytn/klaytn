import {
  loadFixture,
  impersonateAccount,
  stopImpersonatingAccount,
  setBalance,
} from "@nomicfoundation/hardhat-network-helpers";
import { expect } from "chai";
import { ethers } from "hardhat";

import {
  ROLES,
  cnV3InitialLockupFixture,
  cnV3InitialLockupNotDepositedFixture,
  cnV3PublicDelegationFixture,
  cnV3PublicDelegationNotRegisteredFixture,
  gcId,
  unlockAmount,
  unlockTime,
} from "../common/fixtures";
import {
  CnStakingV3,
  CnStakingV3__factory,
  StakingTrackerMockReceiver,
  StakingTrackerMockReceiver__factory,
} from "../../typechain-types";
import { addTime, nowTime, setTime } from "../common/helper";
import { smock } from "@defi-wonderland/smock";

const ONE_WEEK = 60 * 60 * 24 * 7;
/**
 * This is a test for CnStakingV3 contract.
 * It is divided into two parts: Initial Lockup Enabled and Public Delegation Enabled.
 *  1. Initial Lockup Enabled: It tests the basic operations of CnStakingV3 contract with initial lockup enabled.
 *    - Initial setup
 *    - Operations
 *    - Withdraw lockup staking
 *    - Add free stakes (only admin)
 *  2. Public Delegation Enabled: It tests the basic operations of CnStakingV3 contract with public delegation enabled.
 *    - Initial setup including PublicDelegation and RewardAddress deployment
 * Note that delegation features will be detailed in the `publicDelegation` test.
 */
describe("CnStakingV3 tests", function () {
  describe("Initial Lockup Enabled", function () {
    describe("CnStakingV3 setup", function () {
      it("Initial setup", async function () {
        const {
          cnV3,
          stakingTracker,
          deployer,
          nodeId,
          rewardAddress,
          voterAddress,
        } = await loadFixture(cnV3InitialLockupFixture);

        expect(await cnV3.getRoleMember(ROLES.OPERATOR_ROLE, 0)).to.equal(
          deployer.address
        );
        expect(await cnV3.getRoleMember(ROLES.ADMIN_ROLE, 0)).to.equal(
          deployer.address
        );
        expect(await cnV3.nodeId()).to.equal(nodeId);
        expect(await cnV3.rewardAddress()).to.equal(rewardAddress);
        expect(await cnV3.gcId()).to.equal(gcId[0]);
        expect(await cnV3.isPublicDelegationEnabled()).to.equal(false);
        expect(await cnV3.isInitialized()).to.equal(true);
        expect(await ethers.provider.getBalance(cnV3.address)).to.equal(300);

        expect(await cnV3.voterAddress()).to.equal(voterAddress.address);
        expect(await cnV3.stakingTracker()).to.equal(stakingTracker.address);
        expect(await cnV3.publicDelegation()).to.equal(
          ethers.constants.AddressZero
        );

        expect(
          await cnV3.hasRole(ROLES.OPERATOR_ROLE, deployer.address)
        ).to.equal(true);
        expect(await cnV3.hasRole(ROLES.ADMIN_ROLE, deployer.address)).to.equal(
          true
        );
      });
      it("#constructor: Initial lockup should be set correctly", async function () {
        const [deployer, nodeId] = await ethers.getSigners();
        await expect(
          new CnStakingV3__factory(deployer).deploy(
            deployer.address,
            nodeId.address,
            ethers.constants.AddressZero,
            unlockTime,
            unlockAmount
          )
        ).to.be.revertedWith("Initial lockup disabled.");

        await expect(
          new CnStakingV3__factory(deployer).deploy(
            deployer.address,
            nodeId.address,
            deployer.address,
            [unlockTime[1], unlockTime[0]],
            unlockAmount
          )
        ).to.be.revertedWith("Unlock time is not in ascending order.");
      });
      it("#constructor: Zero initial lockup is possible", async function () {
        const [deployer, nodeId] = await ethers.getSigners();
        const cnV3 = await new CnStakingV3__factory(deployer).deploy(
          deployer.address,
          nodeId.address,
          deployer.address,
          [],
          []
        );

        expect(await cnV3.rewardAddress()).to.equal(deployer.address);
        expect(await cnV3.isPublicDelegationEnabled()).to.equal(false);
      });
      it("#constructor: Initial lockup and PD are mutually exclusive", async function () {
        const { cnV3 } = await loadFixture(cnV3InitialLockupFixture);

        // `isPublicDelegationEnabled` should be false if initial lockup is enabled
        expect(await cnV3.isPublicDelegationEnabled()).to.equal(false);
      });
      it("#setGCId: GC id can't be 0", async function () {
        const [deployer, nodeId] = await ethers.getSigners();
        const cnV3 = await new CnStakingV3__factory(deployer).deploy(
          deployer.address,
          nodeId.address,
          deployer.address,
          unlockTime,
          unlockAmount
        );

        await expect(cnV3.setGCId(0)).to.be.revertedWith(
          "GC ID cannot be zero."
        );
      });
      it("#setStakingTracker: Staking tracker should be valid", async function () {
        const [deployer, nodeId] = await ethers.getSigners();
        const cnV3 = await new CnStakingV3__factory(deployer).deploy(
          deployer.address,
          nodeId.address,
          deployer.address,
          unlockTime,
          unlockAmount
        );

        await expect(cnV3.setStakingTracker(ethers.constants.AddressZero)).to.be
          .reverted;
        await expect(cnV3.setStakingTracker(deployer.address)).to.be.reverted;
      });
      it("#setPublicDelegation: can't set PD if PD not enabled", async function () {
        const { cnV3, user } = await loadFixture(
          cnV3InitialLockupNotDepositedFixture
        );

        // Can't set PD if not PD enabled
        await expect(
          cnV3.setPublicDelegation(user.address, "0x1234")
        ).to.be.revertedWith("Public delegation disabled.");
      });
      it("#reviewInitialConditions: only admin can review initial condition", async function () {
        const { cnV3, user } = await loadFixture(
          cnV3InitialLockupNotDepositedFixture
        );

        // Can't set PD if not PD enabled
        await expect(
          cnV3.connect(user).reviewInitialConditions()
        ).to.be.revertedWithCustomError(
          cnV3,
          "AccessControlUnauthorizedAccount"
        );
      });
      it("#reviewInitialConditions: can't review if already reviewed", async function () {
        const { cnV3 } = await loadFixture(
          cnV3InitialLockupNotDepositedFixture
        );

        // Can't set PD if not PD enabled
        await expect(cnV3.reviewInitialConditions()).to.be.revertedWith(
          "Msg.sender already reviewed."
        );
      });
    });
    describe("Ownable operations", function () {
      it("Operation functions are restricted to the admin", async function () {
        const { cnV3, user } = await loadFixture(cnV3InitialLockupFixture);

        await expect(
          cnV3.connect(user).toggleRedelegation()
        ).to.be.revertedWithCustomError(
          cnV3,
          "AccessControlUnauthorizedAccount"
        );
        await expect(
          cnV3.connect(user).updateStakingTracker(user.address)
        ).to.be.revertedWithCustomError(
          cnV3,
          "AccessControlUnauthorizedAccount"
        );
        await expect(
          cnV3.connect(user).updateRewardAddress(user.address)
        ).to.be.revertedWithCustomError(
          cnV3,
          "AccessControlUnauthorizedAccount"
        );
        await expect(
          cnV3.connect(user).updateVoterAddress(user.address)
        ).to.be.revertedWithCustomError(
          cnV3,
          "AccessControlUnauthorizedAccount"
        );
      });
      it("#updateRewardAddress: update pending reward address", async function () {
        const { cnV3, user } = await loadFixture(cnV3InitialLockupFixture);

        await cnV3.updateRewardAddress(user.address);
        expect(await cnV3.pendingRewardAddress()).to.equal(user.address);
      });
      it("#acceptRewardAddress: can't accept reward address if not authorized", async function () {
        const { cnV3, user, voterAddress } = await loadFixture(
          cnV3InitialLockupFixture
        );

        await cnV3.updateRewardAddress(user.address);
        expect(await cnV3.pendingRewardAddress()).to.equal(user.address);

        await expect(
          cnV3.connect(voterAddress).acceptRewardAddress(user.address)
        ).to.be.revertedWith("Unauthorized to accept reward address.");
      });
      it("#acceptRewardAddress: can't accept if given address is different with pending", async function () {
        const { cnV3, user, voterAddress } = await loadFixture(
          cnV3InitialLockupFixture
        );

        await cnV3.updateRewardAddress(user.address);
        expect(await cnV3.pendingRewardAddress()).to.equal(user.address);

        await expect(
          cnV3.acceptRewardAddress(voterAddress.address)
        ).to.be.revertedWith("Given address does not match the pending.");
      });
      it("#updateVoterAddress: can't update to already taken address", async function () {
        const [deployer, nodeId] = await ethers.getSigners();

        const fakeStakingTracker = await smock.fake<StakingTrackerMockReceiver>(
          StakingTrackerMockReceiver__factory.abi
        );
        fakeStakingTracker.CONTRACT_TYPE.returns("StakingTracker");
        fakeStakingTracker.VERSION.returns(1);
        fakeStakingTracker.voterToGCId.returns(1);

        const cnV3 = await new CnStakingV3__factory(deployer).deploy(
          deployer.address,
          nodeId.address,
          deployer.address,
          unlockTime,
          unlockAmount
        );

        await cnV3.setStakingTracker(fakeStakingTracker.address);

        await expect(
          cnV3.updateVoterAddress(deployer.address)
        ).to.be.revertedWith("Voter already taken.");
      });
      it("#updateStakingTracker: can't accept invalid staking tracker", async function () {
        const { cnV3 } = await loadFixture(cnV3InitialLockupFixture);

        const fakeStakingTracker = await smock.fake<StakingTrackerMockReceiver>(
          StakingTrackerMockReceiver__factory.abi
        );
        fakeStakingTracker.CONTRACT_TYPE.returns("FakeStakingTracker");

        await expect(
          cnV3.updateStakingTracker(fakeStakingTracker.address)
        ).to.be.revertedWith("Invalid StakingTracker.");
      });
      it("#updateStakingTracker: can't accept if there's active tracker", async function () {
        const [deployer, nodeId] = await ethers.getSigners();

        const fakeStakingTracker = await smock.fake<StakingTrackerMockReceiver>(
          StakingTrackerMockReceiver__factory.abi
        );
        fakeStakingTracker.CONTRACT_TYPE.returns("StakingTracker");
        fakeStakingTracker.VERSION.returns(1);
        fakeStakingTracker.getLiveTrackerIds.returns([1]);

        const cnV3 = await new CnStakingV3__factory(deployer).deploy(
          deployer.address,
          nodeId.address,
          deployer.address,
          unlockTime,
          unlockAmount
        );

        await cnV3.setStakingTracker(fakeStakingTracker.address);

        await expect(
          cnV3.updateStakingTracker(fakeStakingTracker.address)
        ).to.be.revertedWith(
          "Cannot update tracker when there is an active tracker."
        );
      });
      it("#toggleRedelegation: can't toggle it if not enabled public delegation", async function () {
        const { cnV3 } = await loadFixture(cnV3InitialLockupFixture);

        await expect(cnV3.toggleRedelegation()).to.be.revertedWith(
          "Public delegation disabled."
        );
      });
    });
    describe("Deposit lockup staking", function () {
      it("#depositLockupStakingAndInit: must deposit exact amount", async function () {
        const { cnV3 } = await loadFixture(
          cnV3InitialLockupNotDepositedFixture
        );

        await expect(
          cnV3.depositLockupStakingAndInit({ value: 299 })
        ).to.be.revertedWith("Value does not match.");

        await expect(
          cnV3.depositLockupStakingAndInit({ value: 301 })
        ).to.be.revertedWith("Value does not match.");
      });
      it("#depositLockupStakingAndInit: setup unstaking approver, unstaking claimer roles", async function () {
        const { cnV3, deployer } = await loadFixture(
          cnV3InitialLockupNotDepositedFixture
        );

        expect(
          await cnV3.hasRole(ROLES.UNSTAKING_APPROVER_ROLE, deployer.address)
        ).to.equal(false);
        expect(
          await cnV3.hasRole(ROLES.UNSTAKING_CLAIMER_ROLE, deployer.address)
        ).to.equal(false);

        await cnV3.depositLockupStakingAndInit({ value: 300 });

        expect(
          await cnV3.hasRole(ROLES.UNSTAKING_APPROVER_ROLE, deployer.address)
        ).to.equal(true);
        expect(
          await cnV3.hasRole(ROLES.UNSTAKING_CLAIMER_ROLE, deployer.address)
        ).to.equal(true);
      });
    });
    describe("withdrawLockupStaking", function () {
      it("#withdrawLockupStaking: can't withdraw before unstakeTime", async function () {
        const { cnV3, deployer } = await loadFixture(cnV3InitialLockupFixture);

        await setTime(unlockTime[0] - 100);
        await expect(
          cnV3.withdrawLockupStaking(deployer.address, unlockAmount[0])
        ).to.be.revertedWith("Value is not withdrawable.");
      });
      it("#withdrawLockupStaking: can't withdraw more than deposited", async function () {
        const { cnV3, deployer } = await loadFixture(cnV3InitialLockupFixture);

        await setTime(unlockTime[0] + 100);
        await expect(
          cnV3.withdrawLockupStaking(deployer.address, unlockAmount[0] + 1)
        ).to.be.revertedWith("Value is not withdrawable.");
      });
      it("#withdrawLockupStaking: can't withdraw to zero address", async function () {
        const { cnV3 } = await loadFixture(cnV3InitialLockupFixture);

        await setTime(unlockTime[0] + 100);
        await expect(
          cnV3.withdrawLockupStaking(
            ethers.constants.AddressZero,
            unlockAmount[0]
          )
        ).to.be.revertedWith("Address is null.");
      });
      it("#withdrawLockupStaking: transfer failed if receiver is CA without fallback", async function () {
        const { cnV3, stakingTracker } = await loadFixture(
          cnV3InitialLockupFixture
        );

        // StakingTracker doesn't have fallback function

        await setTime(unlockTime[0] + 100);
        await expect(
          cnV3.withdrawLockupStaking(stakingTracker.address, unlockAmount[0])
        ).to.be.revertedWith("Transfer failed.");
      });
      it("#withdrawLockupStaking: Only owner can withdraw", async function () {
        const { cnV3, user } = await loadFixture(cnV3InitialLockupFixture);

        await setTime(unlockTime[0] + 100);
        await expect(
          cnV3
            .connect(user)
            .withdrawLockupStaking(user.address, unlockAmount[0] + 1)
        ).to.be.revertedWithCustomError(
          cnV3,
          "AccessControlUnauthorizedAccount"
        );
      });
      it("#withdrawLockupStaking: can withdraw after unstakeTime", async function () {
        const { cnV3, deployer } = await loadFixture(cnV3InitialLockupFixture);

        await setTime(unlockTime[0] + 100);
        await cnV3.withdrawLockupStaking(deployer.address, unlockAmount[0]);
        expect(await ethers.provider.getBalance(cnV3.address)).to.equal(
          300 - unlockAmount[0]
        );
      });
    });
    describe("Add free stakes", function () {
      it("#delegate: only admin can add free stakes", async function () {
        const { cnV3, deployer } = await loadFixture(cnV3InitialLockupFixture);

        // Can add free stakes via delegate
        await expect(cnV3.delegate({ value: 100 }))
          .to.emit(cnV3, "DelegateKaia")
          .withArgs(deployer.address, 100);

        // Can add free stakes via transfer
        await expect(deployer.sendTransaction({ to: cnV3.address, value: 100 }))
          .to.emit(cnV3, "DelegateKaia")
          .withArgs(deployer.address, 100);
      });
      it("#delegate: can't delegate zero amount", async function () {
        const { cnV3 } = await loadFixture(cnV3InitialLockupFixture);

        await expect(cnV3.delegate({ value: 0 })).to.be.revertedWith(
          "Invalid amount."
        );
      });
      it("#delegate: can't delegate if not initialized", async function () {
        const { cnV3, deployer } = await loadFixture(
          cnV3InitialLockupNotDepositedFixture
        );

        // Can't delegate if not initialized
        await expect(cnV3.delegate({ value: 100 })).to.be.revertedWith(
          "Contract is not initialized."
        );

        // Can't transfer if not initialized
        await expect(
          deployer.sendTransaction({ to: cnV3.address, value: 100 })
        ).to.be.revertedWith("Contract is not initialized.");
      });
      it("#fallback: not allowed", async function () {
        const { cnV3, deployer } = await loadFixture(cnV3InitialLockupFixture);

        await expect(
          deployer.sendTransaction({
            to: cnV3.address,
            value: 100,
            data: "0x12345678",
          })
        ).to.be.reverted;
      });
    });
    describe("Redelegation is not allowed", function () {
      it("#redelegate: not allowed", async function () {
        const { cnV3, deployer } = await loadFixture(cnV3InitialLockupFixture);

        await expect(cnV3.delegate({ value: 100 }))
          .to.emit(cnV3, "DelegateKaia")
          .withArgs(deployer.address, 100);

        await expect(
          cnV3.redelegate(deployer.address, deployer.address, 100)
        ).to.be.revertedWith("Redelegation disabled.");
      });
      it("#handleRedelegation: not allowed", async function () {
        const { cnV3, deployer } = await loadFixture(cnV3InitialLockupFixture);

        await expect(cnV3.delegate({ value: 100 }))
          .to.emit(cnV3, "DelegateKaia")
          .withArgs(deployer.address, 100);

        await expect(
          cnV3.handleRedelegation(deployer.address, { value: 100 })
        ).to.be.revertedWith("Redelegation disabled.");
      });
    });
    describe("Withdraw free stakes", function () {
      it("#approveStakingWithdrawal: only admin can approve staking withdrawal", async function () {
        const { cnV3, stakingTracker, user } = await loadFixture(
          cnV3InitialLockupFixture
        );

        // Stake some KAIA
        await cnV3.delegate({ value: 100 });

        // Can't approve if not admin
        await expect(
          cnV3.connect(user).approveStakingWithdrawal(user.address, 100)
        ).to.be.revertedWithCustomError(
          cnV3,
          "AccessControlUnauthorizedAccount"
        );

        // Admin can approve staking withdrawal
        await expect(cnV3.approveStakingWithdrawal(user.address, 100))
          .to.emit(cnV3, "ApproveStakingWithdrawal")
          .withArgs(0, user.address, 100, (await nowTime()) + ONE_WEEK + 1)
          .to.emit(stakingTracker, "RefreshStake");
      });
      it("#approveStakingWithdrawal: can't set zero address", async function () {
        const { cnV3 } = await loadFixture(cnV3InitialLockupFixture);

        // Stake some KAIA
        await cnV3.delegate({ value: 100 });

        // Can't set zero address
        await expect(
          cnV3.approveStakingWithdrawal(ethers.constants.AddressZero, 10)
        ).to.be.revertedWith("Address is null.");
      });
      it("#approveStakingWithdrawal: can't approve more than staked", async function () {
        const { cnV3, deployer } = await loadFixture(cnV3InitialLockupFixture);

        // Stake some KAIA
        await cnV3.delegate({ value: 100 });

        // Can't approve more than staked
        await expect(
          cnV3.approveStakingWithdrawal(deployer.address, 101)
        ).to.be.revertedWith("Invalid value.");
      });
      it("#cancelApprovedStakingWithdrawal: only admin can cancel approved staking withdrawal", async function () {
        const { cnV3, stakingTracker, user } = await loadFixture(
          cnV3InitialLockupFixture
        );

        // Stake some KAIA
        await cnV3.delegate({ value: 100 });

        // Approve staking withdrawal
        await cnV3.approveStakingWithdrawal(user.address, 100);

        // Can't cancel if not admin
        await expect(
          cnV3.connect(user).cancelApprovedStakingWithdrawal(0)
        ).to.be.revertedWithCustomError(
          cnV3,
          "AccessControlUnauthorizedAccount"
        );

        // Admin can cancel approved staking withdrawal
        await expect(cnV3.cancelApprovedStakingWithdrawal(0))
          .to.emit(cnV3, "CancelApprovedStakingWithdrawal")
          .withArgs(0, user.address, 100)
          .to.emit(stakingTracker, "RefreshStake");
      });
      it("#cancelApprovedStakingWithdrawal: can't cancel if not approved", async function () {
        const { cnV3 } = await loadFixture(cnV3InitialLockupFixture);

        // Stake some KAIA
        await cnV3.delegate({ value: 100 });

        // Can't cancel if not approved
        await expect(
          cnV3.cancelApprovedStakingWithdrawal(0)
        ).to.be.revertedWith("Withdrawal request does not exist.");
      });
      it("#cancelApprovedStakingWithdrawal: can't cancel if not in Unknwon state", async function () {
        const { cnV3, deployer } = await loadFixture(cnV3InitialLockupFixture);

        // Stake some KAIA
        await cnV3.delegate({ value: 100 });

        // Approve staking withdrawal
        await cnV3.approveStakingWithdrawal(deployer.address, 100);

        await addTime(ONE_WEEK);

        // Withdraw approved staking
        await cnV3.withdrawApprovedStaking(0);

        // Can't cancel if not Unknown state
        await expect(
          cnV3.cancelApprovedStakingWithdrawal(0)
        ).to.be.revertedWith("Invalid state.");
      });
      it("#withdrawApprovedStaking: only admin can withdraw approved staking", async function () {
        const { cnV3, stakingTracker, user } = await loadFixture(
          cnV3InitialLockupFixture
        );

        // Stake some KAIA
        await cnV3.delegate({ value: 100 });

        // Approve staking withdrawal
        await cnV3.approveStakingWithdrawal(user.address, 100);

        await addTime(ONE_WEEK);

        // Can't withdraw if not admin
        await expect(
          cnV3.connect(user).withdrawApprovedStaking(0)
        ).to.be.revertedWithCustomError(
          cnV3,
          "AccessControlUnauthorizedAccount"
        );

        // Admin can withdraw approved staking
        await expect(cnV3.withdrawApprovedStaking(0))
          .to.emit(cnV3, "WithdrawApprovedStaking")
          .withArgs(0, user.address, 100)
          .to.emit(stakingTracker, "RefreshStake");
      });
      it("#withdrawApprovedStaking: can't withdraw if not approved", async function () {
        const { cnV3 } = await loadFixture(cnV3InitialLockupFixture);

        // Stake some KAIA
        await cnV3.delegate({ value: 100 });

        // Can't withdraw if not approved
        await expect(cnV3.withdrawApprovedStaking(0)).to.be.revertedWith(
          "Withdrawal request does not exist."
        );
      });
      it("#withdrawApprovedStaking: can't withdraw if not Unknown state", async function () {
        const { cnV3, deployer } = await loadFixture(cnV3InitialLockupFixture);

        // Stake some KAIA
        await cnV3.delegate({ value: 100 });

        // Approve staking withdrawal
        await cnV3.approveStakingWithdrawal(deployer.address, 100);

        await addTime(ONE_WEEK);

        // Withdraw approved staking
        await cnV3.withdrawApprovedStaking(0);

        // Can't withdraw if not Unknown state
        await expect(cnV3.withdrawApprovedStaking(0)).to.be.revertedWith(
          "Invalid state."
        );
      });
      it("#withdrawApprovedStaking: canceled if not withdrawn within 1 week", async function () {
        const { cnV3, deployer } = await loadFixture(cnV3InitialLockupFixture);

        // Stake some KAIA
        await cnV3.delegate({ value: 100 });

        // Approve staking withdrawal
        await cnV3.approveStakingWithdrawal(deployer.address, 100);

        await addTime(2 * ONE_WEEK);

        // Withdraw approved staking will cancel the request
        await expect(cnV3.withdrawApprovedStaking(0))
          .to.emit(cnV3, "CancelApprovedStakingWithdrawal")
          .withArgs(0, deployer.address, 100);
        expect(await ethers.provider.getBalance(cnV3.address)).to.equal(400);
        expect(await cnV3.unstaking()).to.equal(0);
      });
    });
  });
  describe("Public Delegation Enabled", function () {
    describe("only PD operations", function () {
      it("#approveStakingWithdrawal: can't approve staking withdrawal through owner if PD enabled", async function () {
        const { cnV3s, user1, pd1 } = await loadFixture(
          cnV3PublicDelegationFixture
        );

        await pd1.stake({ value: 100 });

        await expect(
          cnV3s[0].approveStakingWithdrawal(user1.address, 100)
        ).to.be.revertedWithCustomError(
          cnV3s[0],
          "AccessControlUnauthorizedAccount"
        );
      });
      it("#withdrawApprovedStaking: can't withdraw approved staking withdrawal through owner if PD enabled", async function () {
        const { cnV3s, user1, pd1 } = await loadFixture(
          cnV3PublicDelegationFixture
        );

        await pd1.stake({ value: 100 });

        await pd1.redeem(user1.address, 100);

        await addTime(ONE_WEEK);

        await expect(
          cnV3s[0].withdrawApprovedStaking(0)
        ).to.be.revertedWithCustomError(
          cnV3s[0],
          "AccessControlUnauthorizedAccount"
        );
      });
      it("#cancelApprovedStakingWithdrawal: can't cancel approved staking withdrawal through owner if PD enabled", async function () {
        const { cnV3s, user1, pd1 } = await loadFixture(
          cnV3PublicDelegationFixture
        );

        await pd1.stake({ value: 100 });

        await pd1.redeem(user1.address, 100);

        await expect(
          cnV3s[0].cancelApprovedStakingWithdrawal(0)
        ).to.be.revertedWithCustomError(
          cnV3s[0],
          "AccessControlUnauthorizedAccount"
        );
      });
      it("#updateRewardAddress: can't update reward address if it's public delegation", async function () {
        const { cnV3s, user1 } = await loadFixture(cnV3PublicDelegationFixture);

        await expect(
          cnV3s[0].updateRewardAddress(user1.address)
        ).to.be.revertedWith("Public delegation enabled.");
      });
      it("#toggleRedelegation: it works when PD is enabled", async function () {
        const { cnV3s } = await loadFixture(cnV3PublicDelegationFixture);

        await expect(cnV3s[0].toggleRedelegation())
          .to.emit(cnV3s[0], "ToggleRedelegation")
          .withArgs(false);
      });
    });
    describe("CnStakingV3 setup", function () {
      it("Initial setup", async function () {
        const { cnV3s, stakingTracker, deployer, node, voterAddress } =
          await loadFixture(cnV3PublicDelegationFixture);

        for (let i = 0; i < node.length; i++) {
          expect(
            await cnV3s[i].getRoleMember(ROLES.OPERATOR_ROLE, 0)
          ).to.deep.equal(deployer.address);
          expect(
            await cnV3s[i].getRoleMember(ROLES.ADMIN_ROLE, 0)
          ).to.deep.equal(deployer.address);
          expect(await cnV3s[i].nodeId()).to.equal(node[i]);
          expect(await cnV3s[i].rewardAddress()).to.not.equal(
            ethers.constants.AddressZero
          );
          expect(await cnV3s[i].gcId()).to.equal(gcId[i]);
          expect(await cnV3s[i].isPublicDelegationEnabled()).to.equal(true);
          expect(await cnV3s[i].isInitialized()).to.equal(true);
          expect(await ethers.provider.getBalance(cnV3s[i].address)).to.equal(
            0
          );

          expect(await cnV3s[i].voterAddress()).to.equal(voterAddress[i]);
          expect(await cnV3s[i].stakingTracker()).to.equal(
            stakingTracker.address
          );

          const psAddr = await cnV3s[i].publicDelegation();
          expect(psAddr).to.not.equal(ethers.constants.AddressZero);

          expect(await cnV3s[i].hasRole(ROLES.STAKER_ROLE, psAddr)).to.equal(
            true
          );
          expect(
            await cnV3s[i].hasRole(ROLES.UNSTAKING_APPROVER_ROLE, psAddr)
          ).to.equal(true);
          expect(
            await cnV3s[i].hasRole(ROLES.UNSTAKING_CLAIMER_ROLE, psAddr)
          ).to.equal(true);
        }
      });
      it("#constructor: check PublicDelegation", async function () {
        const { cnV3s, pd1, commissionTo } = await loadFixture(
          cnV3PublicDelegationFixture
        );

        // check public delegation contract
        expect(await pd1.MAX_COMMISSION_RATE()).to.equal(1e4);
        expect(await pd1.COMMISSION_DENOMINATOR()).to.equal(1e4);
        expect(await pd1.CONTRACT_TYPE()).to.equal("PublicDelegation");
        expect(await pd1.VERSION()).to.equal(1);
        expect(await pd1.baseCnStakingV3()).to.equal(cnV3s[0].address);
        expect(await pd1.commissionTo()).to.equal(commissionTo[0]);
      });
      it("#setPublicDelegation: commission rate should be less than max", async function () {
        const { cnV3s, pdFactory, deployer } = await loadFixture(
          cnV3PublicDelegationNotRegisteredFixture
        );

        // Too high commission rate
        const pdParam = new ethers.utils.AbiCoder().encode(
          ["tuple(address, address,  uint256, string)"],
          [[deployer.address, deployer.address, 10001, "GC"]]
        );

        await expect(
          cnV3s[0].setPublicDelegation(pdFactory.address, pdParam)
        ).to.be.revertedWith("Commission rate is too high.");
      });
      it("#setPublicDelegation: pd params should be set", async function () {
        const { cnV3s, pdFactory } = await loadFixture(
          cnV3PublicDelegationNotRegisteredFixture
        );

        await expect(
          cnV3s[0].setPublicDelegation(pdFactory.address, "0x")
        ).to.be.revertedWith("Invalid args.");
      });
      it("#setPublicDelegation: pd factory can't be null address'", async function () {
        const { cnV3s, testingPsParam } = await loadFixture(
          cnV3PublicDelegationNotRegisteredFixture
        );

        await expect(
          cnV3s[0].setPublicDelegation(
            ethers.constants.AddressZero,
            testingPsParam
          )
        ).to.be.revertedWith("Address is null.");
      });
      it("#depositLockupStakingAndInit: can't initialize if PD is not registered", async function () {
        const { cnV3s } = await loadFixture(
          cnV3PublicDelegationNotRegisteredFixture
        );

        await expect(cnV3s[0].depositLockupStakingAndInit()).to.be.revertedWith(
          "Not set up properly."
        );
      });
      it("#depositLockupStakingAndInit: msg.value must be 0 if pd enabled", async function () {
        const { cnV3s, pdFactory, deployer } = await loadFixture(
          cnV3PublicDelegationNotRegisteredFixture
        );
        const pdParam = new ethers.utils.AbiCoder().encode(
          ["tuple(address, address,  uint256, string)"],
          [[deployer.address, deployer.address, 0, "GC"]]
        );
        await expect(
          cnV3s[0].setPublicDelegation(pdFactory.address, pdParam)
        ).to.emit(cnV3s[0], "SetPublicDelegation");
        await expect(cnV3s[0].reviewInitialConditions()).to.emit(
          cnV3s[0],
          "ReviewInitialConditions"
        );
        await expect(
          cnV3s[0].depositLockupStakingAndInit({ value: 1 })
        ).to.be.revertedWith("Value does not match.");
      });
      it("#setPublicDelegation: only owner can set public delegation", async function () {
        const { cnV3s, user1, pdFactory } = await loadFixture(
          cnV3PublicDelegationNotRegisteredFixture
        );

        await expect(
          cnV3s[0]
            .connect(user1)
            .setPublicDelegation(pdFactory.address, "0x1234")
        ).to.be.revertedWithCustomError(
          cnV3s[0],
          "AccessControlUnauthorizedAccount"
        );
      });
      it("#setPublicDelegation: can't set PD if already initialized", async function () {
        const { cnV3s, pdFactory } = await loadFixture(
          cnV3PublicDelegationFixture
        );

        // Can't set PD if already set
        await expect(
          cnV3s[0].setPublicDelegation(pdFactory.address, "0x1234")
        ).to.be.revertedWith("Contract has been initialized.");
      });
    });
    describe("Delegation", function () {
      it("Delegation must be done through PublicDelegation to prevent malicious delegation", async function () {
        const { cnV3s, deployer } = await loadFixture(
          cnV3PublicDelegationFixture
        );

        await expect(
          cnV3s[0].delegate({ value: 100 })
        ).to.be.revertedWithCustomError(
          cnV3s[0],
          "AccessControlUnauthorizedAccount"
        );

        await expect(
          deployer.sendTransaction({ to: cnV3s[0].address, value: 100 })
        ).to.be.revertedWithCustomError(
          cnV3s[0],
          "AccessControlUnauthorizedAccount"
        );
      });
      it("Can't delegate if not initialized even if PD is registered", async function () {
        const { cnV3s, pdFactory, testingPsParam } = await loadFixture(
          cnV3PublicDelegationNotRegisteredFixture
        );

        // Register PD
        await cnV3s[0].setPublicDelegation(pdFactory.address, testingPsParam);

        // Can't delegate if not initialized
        const pd = await ethers.getContractAt(
          "PublicDelegation",
          await cnV3s[0].publicDelegation()
        );
        await expect(pd.stake({ value: 100 })).to.be.revertedWith(
          "Contract is not initialized."
        );
      });
      describe("Redelegation", function () {
        it("Redelegation must be done through PublicDelegation to prevent malicious delegation", async function () {
          const { cnV3s, deployer } = await loadFixture(
            cnV3PublicDelegationFixture
          );

          await expect(
            cnV3s[0].redelegate(deployer.address, cnV3s[1].address, 100)
          ).to.be.revertedWith("Redelegation disabled.");
        });
        it("Redelegation target must be a valid CnStakingV3 registered at address book", async function () {
          const { cnV3s, deployer, pd1 } = await loadFixture(
            cnV3PublicDelegationFixture
          );

          await impersonateAccount(pd1.address);

          const ps1Signer = await hre.ethers.getSigner(pd1.address);

          const fakeCnV3 = await smock.fake<CnStakingV3>(
            CnStakingV3__factory.abi
          );

          await expect(
            cnV3s[0]
              .connect(ps1Signer)
              .redelegate(deployer.address, fakeCnV3.address, 100)
          ).to.be.revertedWith("Invalid CnStakingV3.");

          await stopImpersonatingAccount(pd1.address);
        });
        it("Handle redelegation only from a valid CnStakingV3 registered at address book", async function () {
          const { cnV3s, deployer, pd1 } = await loadFixture(
            cnV3PublicDelegationFixture
          );

          const fakeCnV3 = await smock.fake<CnStakingV3>(
            CnStakingV3__factory.abi
          );

          await impersonateAccount(fakeCnV3.address);

          await setBalance(fakeCnV3.address, 1e18);

          const fakeCnV3Signer = await hre.ethers.getSigner(fakeCnV3.address);

          await expect(
            cnV3s[1]
              .connect(fakeCnV3Signer)
              .handleRedelegation(deployer.address, { value: 100 })
          ).to.be.revertedWith("Invalid CnStakingV3.");

          await stopImpersonatingAccount(pd1.address);
        });
      });
    });
  });
});
