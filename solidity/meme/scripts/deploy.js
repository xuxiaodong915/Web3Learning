const { ethers } = require("hardhat");
require("dotenv").config();

async function main() {
    const name = process.env.TOKEN_NAME || "ShibStyleMeme";
    const symbol = process.env.TOKEN_SYMBOL || "SSM";
    const supply = ethers.parseUnits((process.env.TOKEN_SUPPLY || "1000000000000").toString(), 18);
    const router = process.env.ROUTER_ADDRESS || ethers.ZeroAddress;

    const [deployer] = await ethers.getSigners();
    console.log("Deployer:", deployer.address);

    const taxWalletEnv = process.env.TAX_WALLET;
    const taxWallet = taxWalletEnv && taxWalletEnv !== "" ? taxWalletEnv : deployer.address;
    const taxBps = Number(process.env.TAX_BPS || 300);
    const maxTx = ethers.parseUnits((process.env.MAX_TX || "1000000000").toString(), 18);
    const dailyLimit = Number(process.env.DAILY_TX_LIMIT || 10);

    const MemeToken = await ethers.getContractFactory("MemeToken");
    const token = await MemeToken.deploy(
        name,
        symbol,
        supply,
        router,
        taxWallet,
        taxBps,
        maxTx,
        dailyLimit
    );
    await token.waitForDeployment();
    const addr = await token.getAddress();
    console.log("MemeToken deployed:", addr);

    if (router !== ethers.ZeroAddress) {
        console.log("Router set:", await token.dexRouter());
        console.log("Pair:", await token.dexPair());
    }
}

main().catch((e) => {
    console.error(e);
    process.exit(1);
}); 