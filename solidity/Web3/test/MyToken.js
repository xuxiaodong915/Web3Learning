const hre = require("hardhat")
const { expect } = require("chai")

describe("MyToken Test", async () => {

    const { ethers } = hre
    const initiaSupply = 1000
    let MyTokenContract;
    let account1, account2;
    //let account3;
    before(async () => {
        [account1, account2] = await ethers.getSigners();
        //[account3] = await ethers.getSigner();
        console.log(account1.address, account2.address)
        //console.log("account3=", account3)
        const MyToken = await ethers.getContractFactory("MyToken")
        MyTokenContract = await MyToken.connect(account1).deploy(initiaSupply)
        MyTokenContract.waitForDeployment()
        const contractAddress = await MyTokenContract.getAddress()
        console.log(contractAddress.length)
        expect(contractAddress).to.length.greaterThan(0)

        console.log(contractAddress, "==contractAddress==");
    })

    it("测试合约的name symbol decimals", async () => {
        const name = await MyTokenContract.name()
        const symbol = await MyTokenContract.symbol()
        const decimals = await MyTokenContract.decimals()

        expect(name).to.equal("MyToken")
        expect(symbol).to.equal("MTK")
        expect(decimals).to.equal(18)
    })

    it("测试转账", async () => {
        const balanceOfAccount1 = await MyTokenContract.balanceOf(account1)
        console.log("balanceOfAccount1:", balanceOfAccount1)
        console.log("balanceOfAccount1 address:", account1.getAddress())
        expect(balanceOfAccount1).to.equal(initiaSupply)

        const resp = await MyTokenContract.transfer(account2, initiaSupply / 2);
        //console.log(resp)

        const balanceOfAccount2 = await MyTokenContract.balanceOf(account2)
        console.log("balanceOfAccount2:", balanceOfAccount2)
        console.log("balanceOfAccount2 address:", account2.getAddress())
        expect(balanceOfAccount2).to.equal(initiaSupply / 2)
    })
})