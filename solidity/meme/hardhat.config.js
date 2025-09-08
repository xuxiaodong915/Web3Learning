require("@nomicfoundation/hardhat-ethers");
require("dotenv").config();

/** @type import('hardhat/config').HardhatUserConfig */
module.exports = {
	solidity: {
		version: "0.8.20",
		settings: {
			optimizer: { enabled: true, runs: 200 }
		}
	},
	networks: {
		// 使用时在 .env 配置 RPC_URL 和 PRIVATE_KEY
		sepolia: {
			url: process.env.RPC_URL || "",
			accounts: process.env.PRIVATE_KEY ? [process.env.PRIVATE_KEY] : []
		}
	}
}; 