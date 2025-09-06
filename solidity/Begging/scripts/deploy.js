const { ethers } = require("hardhat");

async function main() {
    const [deployer] = await ethers.getSigners();
    console.log("Deployer:", deployer.address);

    const BeggingContract = await ethers.getContractFactory("BeggingContract");
    const contract = await BeggingContract.deploy();
    await contract.waitForDeployment();

    console.log("BeggingContract deployed at:", await contract.getAddress());
}

main().catch((error) => {
    console.error(error);
    process.exitCode = 1;
}); 