const { expect } = require("chai");
const { ethers } = require("hardhat");

describe("BeggingContract", function () {
    let contract;
    let owner;
    let donor1;
    let donor2;

    beforeEach(async function () {
        [owner, donor1, donor2] = await ethers.getSigners();
        const Factory = await ethers.getContractFactory("BeggingContract", owner);
        contract = await Factory.deploy();
        await contract.waitForDeployment();
    });

    it("records donations via donate()", async function () {
        const value = ethers.parseEther("1");
        await expect(contract.connect(donor1).donate({ value }))
            .to.changeEtherBalances([donor1, contract], [-value, value]);
        const donated = await contract.getDonation(donor1.address);
        expect(donated).to.equal(value);
    });

    it("accumulates donations via receive()", async function () {
        const v1 = ethers.parseEther("0.5");
        const v2 = ethers.parseEther("0.2");
        await donor1.sendTransaction({ to: await contract.getAddress(), value: v1 });
        await donor1.sendTransaction({ to: await contract.getAddress(), value: v2 });
        const donated = await contract.getDonation(donor1.address);
        expect(donated).to.equal(v1 + v2);
    });

    it("only owner can withdraw and transfer all balance", async function () {
        const donation = ethers.parseEther("1.3");
        await contract.connect(donor1).donate({ value: donation });

        await expect(contract.connect(donor1).withdraw()).to.be.revertedWith("Not owner");

        const beforeOwner = await ethers.provider.getBalance(owner.address);
        const tx = await contract.connect(owner).withdraw();
        const receipt = await tx.wait();
        const gasUsedFee = receipt.gasUsed * receipt.gasPrice;
        const afterOwner = await ethers.provider.getBalance(owner.address);
        const contractBalance = await ethers.provider.getBalance(await contract.getAddress());

        expect(contractBalance).to.equal(0n);
        expect(afterOwner + gasUsedFee - beforeOwner).to.equal(donation);
    });
}); 