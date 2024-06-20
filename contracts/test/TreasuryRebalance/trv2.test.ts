import { expect } from "chai";
import { ethers } from "hardhat";
import { loadFixture } from "@nomicfoundation/hardhat-network-helpers";
import { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers";
import { Contract } from "ethers";
import { FakeContract, smock } from "@defi-wonderland/smock";
import { SenderTest1, SenderTest1__factory } from "../../typechain-types";

const executionBlock = 200;
const value = hre.ethers.utils.parseEther("20");
const memo =
  '{ "retirees": [ { "zeroed": "0x38138d89c321b3b5f421e9452b69cf29e4380bae", "balance": 20000000000000000000 }, { "zeroed": "0x0a33a1b99bd67a7189573dd74de80293afdf969a", "balance": 20000000000000000000 } ], "newbies": [ { "newbie": "0x38138d89c321b3b5f421e9452b69cf29e4380bae", "fundAllocated": 10000000000000000000 }, { "newbie": "0x0a33a1b99bd67a7189573dd74de80293afdf969a", "fundAllocated": 10000000000000000000 } ], "burnt": 7.2e+37, "success": true }';
/**
 * @dev This test is for TreasuryRebalanceV2.sol
 */
describe("TreasuryRebalanceV2", function () {
  async function buildFixture() {
    await ethers.provider.send("hardhat_reset", []);
    // Prepare parameters for deploying contracts
    const [
      _rebalance_manager,
      _eoaZeroed,
      _allocated1,
      _allocated2,
      ..._zeroedAdmins
    ] = (await ethers.getSigners()).slice(0, 8);

    // Deploy and fund KIF and KEF which will be zeroed after rebalancing.
    const KIF = await ethers.getContractFactory("SenderTest1");
    const _zeroed1 = await KIF.deploy();
    await _zeroed1.emptyAdminList();
    await _zeroed1.addAdmin(_zeroedAdmins[0].address);
    await _zeroed1.addAdmin(_zeroedAdmins[1].address);
    await _zeroed1.changeMinReq(2);
    await _rebalance_manager.sendTransaction({
      to: _zeroedAdmins[0].address,
      value: value,
    });
    await _rebalance_manager.sendTransaction({
      to: _zeroedAdmins[1].address,
      value: value,
    });
    await _rebalance_manager.sendTransaction({
      to: _zeroed1.address,
      value: value,
    });

    const KEF = await ethers.getContractFactory("SenderTest1");
    const _zeroed2 = await KEF.deploy();
    await _zeroed2.emptyAdminList();
    await _zeroed2.addAdmin(_zeroedAdmins[2].address);
    await _zeroed2.addAdmin(_zeroedAdmins[3].address);
    await _zeroed2.changeMinReq(2);
    await _rebalance_manager.sendTransaction({
      to: _zeroedAdmins[2].address,
      value: value,
    });
    await _rebalance_manager.sendTransaction({
      to: _zeroedAdmins[3].address,
      value: value,
    });
    await _rebalance_manager.sendTransaction({
      to: _zeroed2.address,
      value: value,
    });

    const _mockZeroed3 = await smock.fake<SenderTest1>(
      SenderTest1__factory.abi,
      {
        address: "0x38138d89c321b3b5f421e9452b69cf29e4380bae",
      }
    );

    const REBALANCE = await ethers.getContractFactory("TreasuryRebalanceV2");
    const _trV2 = await REBALANCE.deploy(executionBlock);

    return {
      _rebalance_manager,
      _eoaZeroed,
      _zeroedAdmins,
      _zeroed1,
      _zeroed2,
      _mockZeroed3,
      _allocated1,
      _allocated2,
      _trV2,
    };
  }

  let rebalance_manager: SignerWithAddress, zeroedAdmins: SignerWithAddress[];
  let eoaZeroed: SignerWithAddress,
    zeroed1: Contract,
    zeroed2: Contract,
    allocated1: SignerWithAddress,
    allocated2: SignerWithAddress;

  before(async function () {
    const {
      _rebalance_manager,
      _zeroedAdmins,
      _eoaZeroed,
      _zeroed1,
      _zeroed2,
      _allocated1,
      _allocated2,
    } = await loadFixture(buildFixture);
    rebalance_manager = _rebalance_manager;
    zeroedAdmins = _zeroedAdmins;
    eoaZeroed = _eoaZeroed;
    zeroed1 = _zeroed1;
    zeroed2 = _zeroed2;
    allocated1 = _allocated1;
    allocated2 = _allocated2;
  });

  describe("Deploy", function () {
    let trV2: Contract;

    before(async function () {
      const { _trV2 } = await loadFixture(buildFixture);
      trV2 = _trV2;
    });

    it("Should set the correct initial values for main treasuryRebalanceV2", async function () {
      const currentBlock = await ethers.provider.getBlockNumber();
      expect(await trV2.status()).to.equal(0);
      expect(await trV2.getTreasuryAmount()).to.equal(0);
      expect(await trV2.rebalanceBlockNumber()).to.be.greaterThan(currentBlock);
    });

    it("Should revert if the rebalance block number is less than current block", async function () {
      const REBALANCE = await ethers.getContractFactory("TreasuryRebalanceV2");
      const currentBlock = await ethers.provider.getBlockNumber();
      await expect(REBALANCE.deploy(currentBlock)).to.be.revertedWith(
        "rebalance blockNumber should be greater than current block"
      );
    });
  });

  describe("registerZeroed", async function () {
    let trV2: Contract;

    before(async function () {
      const { _trV2 } = await loadFixture(buildFixture);
      trV2 = _trV2;
    });

    it("Should add a zeroed and emit a RegisterZeroed event", async function () {
      await expect(trV2.registerZeroed(zeroed1.address))
        .to.emit(trV2, "ZeroedRegistered")
        .withArgs(zeroed1.address);
      const zeroedInfo = await trV2.getZeroed(zeroed1.address);
      expect(zeroedInfo[0]).to.equal(zeroed1.address);
      expect(zeroedInfo[1].length).to.equal(0);
    });

    it("Should not allow non-owner to add a zeroed", async function () {
      await expect(
        trV2.connect(zeroedAdmins[0]).registerZeroed(zeroed2.address)
      ).to.be.revertedWith("Ownable: caller is not the owner");
    });

    it("Should not allow adding the same zeroed twice", async function () {
      await trV2.registerZeroed(zeroed2.address);
      expect(await trV2.zeroedExists(zeroed2.address)).to.equal(true);
      await expect(trV2.registerZeroed(zeroed2.address)).to.be.revertedWith(
        "Zeroed address is already registered"
      );
    });

    it("Should revert if the zeroed address is zero", async function () {
      await expect(
        trV2.registerZeroed(ethers.constants.AddressZero)
      ).to.be.revertedWith("Invalid address");
    });

    it("zeroeds length should be two", async function () {
      expect(await trV2.getZeroedCount()).to.equal(2);
    });
  });

  describe("removeZeroed", function () {
    let trV2: Contract;

    before(async function () {
      const { _trV2 } = await loadFixture(buildFixture);
      trV2 = _trV2;

      await trV2.registerZeroed(zeroed1.address);
      await trV2.registerZeroed(zeroed2.address);
    });

    it("Should remove a zeroed and emit a RemoveZeroed event", async function () {
      await expect(trV2.removeZeroed(zeroed1.address))
        .to.emit(trV2, "ZeroedRemoved")
        .withArgs(zeroed1.address);
      await expect(trV2.getZeroed(zeroed1.address)).to.be.revertedWith(
        "Zeroed not registered"
      );
    });

    it("Should not allow removing a non-existent zeroed", async function () {
      await expect(
        trV2.removeZeroed(rebalance_manager.address)
      ).to.be.revertedWith("Zeroed not registered");
    });

    it("Should not allow non-owner to remove a zeroed", async function () {
      await expect(
        trV2.connect(zeroedAdmins[0]).removeZeroed(zeroed2.address)
      ).to.be.revertedWith("Ownable: caller is not the owner");
    });
  });

  describe("registerAllocated", function () {
    let trV2: Contract;

    before(async function () {
      const { _trV2 } = await loadFixture(buildFixture);
      trV2 = _trV2;
    });

    it("Should register allocated address and its fund distribution and emit a RegisterAllocated event", async function () {
      const amount = hre.ethers.utils.parseEther("20");

      await expect(trV2.registerAllocated(allocated1.address, amount))
        .to.emit(trV2, "AllocatedRegistered")
        .withArgs(allocated1.address, amount);
      const allocateds = await trV2.getAllocated(allocated1.address);
      expect(allocateds[0]).to.equal(allocated1.address);
      expect(allocateds[1]).to.equal(amount);
      expect(await trV2.getAllocatedCount()).to.equal(1);

      const treasuryAmount = await trV2.getTreasuryAmount();
      expect(treasuryAmount).to.equal(amount);
    });

    it("Should revert if register allocated twice", async function () {
      const amount1 = hre.ethers.utils.parseEther("20");
      expect(await trV2.allocatedExists(allocated1.address)).to.equal(true);
      await expect(
        trV2.registerAllocated(allocated1.address, amount1)
      ).to.be.revertedWith("Allocated address is already registered");
    });

    it("Should not allow non-owner to add a allocated", async function () {
      const amount1 = hre.ethers.utils.parseEther("20");
      await expect(
        trV2
          .connect(zeroedAdmins[0])
          .registerAllocated(allocated1.address, amount1)
      ).to.be.revertedWith("Ownable: caller is not the owner");
    });

    it("Should revert if the allocated address is zero", async function () {
      const amount1 = hre.ethers.utils.parseEther("20");
      await expect(
        trV2.registerAllocated(ethers.constants.AddressZero, amount1)
      ).to.be.revertedWith("Invalid address");
    });

    it("Should revert if the amount is set to 0", async function () {
      await expect(
        trV2.registerAllocated(allocated2.address, 0)
      ).to.be.revertedWith("Amount cannot be set to 0");
    });
  });

  describe("removeAllocated", function () {
    let trV2: Contract;

    before(async function () {
      const { _trV2 } = await loadFixture(buildFixture);
      trV2 = _trV2;

      await trV2.registerAllocated(allocated1.address, value);
      await trV2.registerAllocated(allocated2.address, value);
    });

    it("Should remove allocated and emit RemoveAllocated event", async function () {
      await expect(trV2.removeAllocated(allocated1.address))
        .to.emit(trV2, "AllocatedRemoved")
        .withArgs(allocated1.address);
      expect(await trV2.getAllocatedCount()).to.equal(1);
      expect(await trV2.getTreasuryAmount()).to.equal(value);
      await expect(trV2.getAllocated(allocated1.address)).to.be.revertedWith(
        "Allocated not registered"
      );
    });

    it("Should not remove unregistered allocated", async function () {
      await expect(trV2.removeAllocated(zeroed1.address)).to.be.reverted;
    });

    it("Should not allow non-owner to remove a allocated", async function () {
      await expect(
        trV2.connect(allocated2).removeAllocated(allocated2.address)
      ).to.be.revertedWith("Ownable: caller is not the owner");
    });
  });

  describe("approve", function () {
    let trV2: Contract, mockZeroed3: FakeContract<SenderTest1>;

    before(async function () {
      const { _trV2, _mockZeroed3 } = await loadFixture(buildFixture);
      trV2 = _trV2;
      mockZeroed3 = _mockZeroed3;

      await trV2.registerZeroed(zeroed1.address);
      await trV2.registerZeroed(eoaZeroed.address);
      await trV2.registerZeroed(mockZeroed3.address);
      await trV2.registerAllocated(allocated1.address, value);
      await trV2.registerAllocated(allocated2.address, value);
      await trV2.finalizeRegistration();
    });

    it("Should approve zeroed contract if msg.sender is an admin of zeroed", async function () {
      const zeroed = await trV2.getZeroed(zeroed1.address);
      expect(zeroed[1].length).to.equal(0);
      const tx = await trV2.connect(zeroedAdmins[0]).approve(zeroed1.address);

      const updatedZeroedDetails = await trV2.getZeroed(zeroed1.address);
      expect(updatedZeroedDetails[1][0]).to.equal(zeroedAdmins[0].address);

      await expect(tx)
        .to.emit(trV2, "Approved")
        .withArgs(zeroed1.address, zeroedAdmins[0].address, 1);
    });

    it("Should approve zeroed eoa if msg.sender is same as zeroed eoa", async function () {
      const zeroed = await trV2.getZeroed(eoaZeroed.address);
      expect(zeroed[0]).to.equal(eoaZeroed.address);
      await trV2.connect(eoaZeroed).approve(eoaZeroed.address);
    });

    it("Should revert if zeroed is already approved", async function () {
      await expect(
        trV2.connect(eoaZeroed).approve(eoaZeroed.address)
      ).to.be.revertedWith("Already approved");
    });

    it("Should revert if zeroed is not registered", async function () {
      // try to approve unregistered zeroed
      await expect(trV2.approve(zeroed2.address)).to.be.revertedWith(
        "zeroed needs to be registered before approval"
      );
    });

    it("Should revert if zeroed is a EOA and if msg.sender is not the admin of zeroed", async function () {
      await expect(trV2.approve(eoaZeroed.address)).to.be.revertedWith(
        "zeroedAddress is not the msg.sender"
      );
    });

    it("Should revert if zeroed is a contract address but does not have getState() method", async function () {
      mockZeroed3.getState.reverts();
      await expect(
        trV2.approve(mockZeroed3.address)
      ).to.be.revertedWithoutReason();
    });

    it("Should revert if zeroed is a contract but adminList is empty", async function () {
      mockZeroed3.getState.returns([[], 0]);
      await expect(trV2.approve(mockZeroed3.address)).to.be.revertedWith(
        "admin list cannot be empty"
      );
    });

    it("Should not approve if zeroed is a contract but msg.sender is not the admin", async function () {
      await expect(trV2.approve(zeroed1.address)).to.be.revertedWith(
        "msg.sender is not the admin"
      );
    });
  });

  describe("setStatus", function () {
    let trV2: Contract, mockZeroed3: FakeContract<SenderTest1>;

    describe("FinalizeRegistration", async function () {
      before(async function () {
        const { _trV2 } = await loadFixture(buildFixture);
        trV2 = _trV2;

        await trV2.registerZeroed(zeroed1.address);
        await trV2.registerAllocated(allocated1.address, value);
      });

      it("should set status to Registered and emit StatusChanged event", async function () {
        await expect(trV2.finalizeRegistration())
          .to.emit(trV2, "StatusChanged")
          .withArgs(1);
        expect(await trV2.status()).to.equal(1);
      });

      it("Should not register zeroed when contract is not in Initialized state", async function () {
        await expect(trV2.registerZeroed(zeroed2.address)).to.be.revertedWith(
          "Not in the designated status"
        );
      });
      it("Should not register allocated when contract is not in Initialized state", async function () {
        await expect(
          trV2.registerAllocated(allocated1.address, value)
        ).to.be.revertedWith("Not in the designated status");
      });
      it("Should revert if the current status is tried to set again", async function () {
        await expect(trV2.finalizeRegistration()).to.be.revertedWith(
          "Not in the designated status"
        );
      });

      it("Should revert if owner tries to set Finalize after Registered", async function () {
        await expect(trV2.finalizeContract(memo)).to.be.revertedWith(
          "Not in the designated status"
        );
      });
      it("should revert if not called by the owner", async () => {
        await expect(
          trV2.connect(zeroedAdmins[0]).finalizeRegistration()
        ).to.be.revertedWith("Ownable: caller is not the owner");
      });
      it("Should not remove allocated when contract is not in Initialized state", async function () {
        await expect(
          trV2.removeAllocated(allocated1.address)
        ).to.be.revertedWith("Not in the designated status");
      });
      it("Should not remove zeroed when contract is not in Initialized state", async function () {
        await expect(trV2.removeZeroed(zeroed2.address)).to.be.revertedWith(
          "Not in the designated status"
        );
      });
      it("should revert if its not called at the Approved stage", async () => {
        await expect(trV2.finalizeContract(memo)).to.be.revertedWith(
          "Not in the designated status"
        );
      });
    });

    describe("FinalizeApproval", async function () {
      before(async function () {
        const { _trV2 } = await loadFixture(buildFixture);
        trV2 = _trV2;

        await trV2.registerZeroed(zeroed1.address);
        await trV2.registerAllocated(allocated1.address, value);
        await trV2.finalizeRegistration();
        await trV2.connect(zeroedAdmins[0]).approve(zeroed1.address);
        await trV2.connect(zeroedAdmins[1]).approve(zeroed1.address);
      });
      it("should set status to Approved and emit StatusChanged event", async function () {
        await expect(trV2.finalizeApproval())
          .to.emit(trV2, "StatusChanged")
          .withArgs(2);
        expect(await trV2.status()).to.equal(2);
      });
      it("should revert if owner tries to set Registered after Approved", async function () {
        await expect(trV2.finalizeRegistration()).to.be.revertedWith(
          "Not in the designated status"
        );
      });
      it("should revert if not called by the owner", async () => {
        await expect(
          trV2.connect(zeroedAdmins[0]).finalizeApproval()
        ).to.be.revertedWith("Ownable: caller is not the owner");
      });
    });

    describe("Should revert finalizeApproval if zeroed contract can't reach Quorom", function () {
      beforeEach(
        "Should set status to Approved and emit StatusChanged event",
        async function () {
          const { _trV2, _mockZeroed3 } = await loadFixture(buildFixture);
          trV2 = _trV2;
          mockZeroed3 = _mockZeroed3;

          await trV2.registerAllocated(allocated1.address, value);
        }
      );

      it("Should revert if min required admins does not approve", async function () {
        await trV2.registerZeroed(zeroed1.address);
        await trV2.finalizeRegistration();
        await expect(trV2.finalizeApproval()).to.be.revertedWith(
          "min required admins should approve"
        );
      });

      it("Should revert if approved admin change during the contract ", async function () {
        await trV2.registerZeroed(mockZeroed3.address);
        await trV2.finalizeRegistration();
        mockZeroed3.getState.returns([[rebalance_manager.address], 1]);
        await trV2.approve(mockZeroed3.address);
        mockZeroed3.getState.returns([
          [rebalance_manager.address, zeroedAdmins[0].address],
          2,
        ]);
        await expect(trV2.finalizeApproval()).to.be.revertedWith(
          "min required admins should approve"
        );
      });

      it("Should revert if EOA did not approve", async function () {
        await trV2.registerZeroed(eoaZeroed.address);
        await trV2.finalizeRegistration();
        await expect(trV2.finalizeApproval()).to.be.revertedWith(
          "EOA should approve"
        );
        await trV2.connect(eoaZeroed).approve(eoaZeroed.address);
        await trV2.finalizeApproval();
      });
    });

    // TreasuryRebalance - below test should not pass
    // TreasuryRebalanceV2 - below test should pass
    it("Should set status to Approved when treasury amount exceeds balance of zeroeds", async function () {
      const { _trV2 } = await loadFixture(buildFixture);
      trV2 = _trV2;

      const amount = hre.ethers.utils.parseEther("50");
      await trV2.registerZeroed(zeroed1.address);
      await trV2.registerAllocated(allocated1.address, value);
      await trV2.registerAllocated(allocated2.address, amount);
      await trV2.finalizeRegistration();
      await trV2.connect(zeroedAdmins[0]).approve(zeroed1.address);
      await trV2.connect(zeroedAdmins[1]).approve(zeroed1.address);
      await trV2.finalizeApproval();
      expect(await trV2.status()).to.equal(2);
    });

    it("Should revert if owner tries to set Approved before Registered", async function () {
      await expect(trV2.finalizeApproval()).to.be.revertedWith(
        "Not in the designated status"
      );
    });
  });

  describe("reset", function () {
    let trV2: Contract;

    before(async function () {
      const { _trV2 } = await loadFixture(buildFixture);
      trV2 = _trV2;
    });

    it("should reset all storage values except rebalance blocknumber to 0 at Initialize state", async function () {
      await trV2.reset();
      expect(await trV2.getZeroedCount()).to.equal(0);
      expect(await trV2.getAllocatedCount()).to.equal(0);
      expect(await trV2.getTreasuryAmount()).to.equal(0);
      expect(await trV2.memo()).to.equal("");
      expect(await trV2.status()).to.equal(0);
      expect(await trV2.rebalanceBlockNumber()).to.not.equal(0);
    });

    it("should reset all storage values except rebalance blocknumber to 0 at Registered state", async function () {
      await trV2.registerZeroed(zeroed1.address);
      await trV2.registerZeroed(zeroed2.address);
      await trV2.registerAllocated(allocated1.address, value);
      await trV2.registerAllocated(allocated2.address, value);
      await trV2.finalizeRegistration();
      expect(await trV2.getZeroedCount()).to.equal(2);
      expect(await trV2.getAllocatedCount()).to.equal(2);

      await trV2.reset();
      expect(await trV2.getZeroedCount()).to.equal(0);
      expect(await trV2.getAllocatedCount()).to.equal(0);
      expect(await trV2.getTreasuryAmount()).to.equal(0);
      expect(await trV2.memo()).to.equal("");
      expect(await trV2.status()).to.equal(0);
      expect(await trV2.rebalanceBlockNumber()).to.not.equal(0);
    });

    it("should reset all storage values except rebalance blocknumber to 0 at Approved state", async function () {
      await trV2.registerZeroed(zeroed1.address);
      await trV2.registerZeroed(zeroed2.address);
      await trV2.registerAllocated(allocated1.address, value);
      await trV2.registerAllocated(allocated2.address, value);
      await trV2.finalizeRegistration();
      await trV2.connect(zeroedAdmins[0]).approve(zeroed1.address);
      await trV2.connect(zeroedAdmins[1]).approve(zeroed1.address);
      await trV2.connect(zeroedAdmins[2]).approve(zeroed2.address);
      await trV2.connect(zeroedAdmins[3]).approve(zeroed2.address);
      await trV2.finalizeApproval();

      expect(await trV2.getZeroedCount()).to.equal(2);
      expect(await trV2.getAllocatedCount()).to.equal(2);

      await trV2.reset();
      expect(await trV2.getZeroedCount()).to.equal(0);
      expect(await trV2.getAllocatedCount()).to.equal(0);
      expect(await trV2.getTreasuryAmount()).to.equal(0);
      expect(await trV2.memo()).to.equal("");
      expect(await trV2.status()).to.equal(0);
      expect(await trV2.rebalanceBlockNumber()).to.not.equal(0);
    });

    it("Should not allow non-owner to reset", async function () {
      await expect(trV2.connect(zeroedAdmins[0]).reset()).to.be.revertedWith(
        "Ownable: caller is not the owner"
      );
    });
  });

  describe("updateRebalanceBlocknumber", async function () {
    let trV2: Contract;

    before(async function () {
      const { _trV2 } = await loadFixture(buildFixture);
      trV2 = _trV2;
    });

    it("should revert if current block is larger than rebalanceBlockNumber", async function () {
      await expect(
        trV2.updateRebalanceBlocknumber(await ethers.provider.getBlockNumber())
      ).to.be.revertedWith(
        "rebalance blockNumber should be greater than current block"
      );
    });
    it("should set rebalance blocknumber if current block is smaller than rebalanceBlockNumber", async function () {
      await trV2.updateRebalanceBlocknumber(executionBlock + 10);
      expect(await trV2.rebalanceBlockNumber()).to.equal(executionBlock + 10);
      await trV2.updateRebalanceBlocknumber(executionBlock);
      expect(await trV2.rebalanceBlockNumber()).to.equal(executionBlock);
    });
  });

  describe("finalizeContract", function () {
    let trV2: Contract;

    before(async function () {
      const { _trV2 } = await loadFixture(buildFixture);
      trV2 = _trV2;
    });

    before(async function () {
      await trV2.registerZeroed(zeroed1.address);
      await trV2.registerZeroed(zeroed2.address);
      await trV2.registerAllocated(allocated1.address, value);
      await trV2.registerAllocated(allocated2.address, value);
      await trV2.finalizeRegistration();
      await trV2.connect(zeroedAdmins[0]).approve(zeroed1.address);
      await trV2.connect(zeroedAdmins[1]).approve(zeroed1.address);
      await trV2.connect(zeroedAdmins[2]).approve(zeroed2.address);
      await trV2.connect(zeroedAdmins[3]).approve(zeroed2.address);
      await trV2.finalizeApproval();
    });
    it("should revert if current block is larger than rebalanceBlockNumber", async function () {
      await expect(
        trV2.updateRebalanceBlocknumber(await ethers.provider.getBlockNumber())
      ).to.be.revertedWith(
        "rebalance blockNumber should be greater than current block"
      );
    });
    it("should revert before rebalanceBlockNumber", async () => {
      await expect(trV2.finalizeContract(memo)).to.be.revertedWith(
        "Contract can only finalize after executing rebalancing"
      );
    });
    it("should set the memo and status to Finalized and emit Finalize event", async function () {
      await hre.network.provider.send("hardhat_mine", ["0xC8"]);
      await expect(trV2.finalizeContract(memo))
        .to.emit(trV2, "Finalized")
        .withArgs(memo, 3);
      expect(await trV2.memo()).to.equal(memo);
      expect(await trV2.status()).to.equal(3);
    });
  });
});
