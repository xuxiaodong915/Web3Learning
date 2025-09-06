const { ethers, upgrades } = require("hardhat");
const path = require("path");
const fs = require("fs");

module.exports = async ({ getNamedAccounts, deployments }) => {
    const { save } = deployments
    const { deployer } = await getNamedAccounts()
    console.log("部署用户地址：", deployer)

    const storePath = path.resolve(__dirname, "./.cache/proxyNftAuction.json");
    const storeData = fs.readFileSync(storePath, "utf-8");
    const { proxyAddress, implAddress, abi } = JSON.parse(storeData);

    // 升级版的代理合约
    const NftAuctionV2 = await ethers.getContractFactory("NftAuctionV2")

    // 升级代理合约
    const nftAuctionProxy2 = await upgrades.upgradeProxy(proxyAddress, NftAuctionV2)
    await nftAuctionProxy2.waitForDeployment()
    const proxyAddress2 = await nftAuctionProxy2.getAddress()

    await save("NftAuctionProxyV2", {
        abi,
        address: proxyAddress2,
    })
};

module.exports.tags = ["upgradeNftAuction"];