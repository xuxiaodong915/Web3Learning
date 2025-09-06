require("@nomicfoundation/hardhat-ethers");

/** @type import('hardhat/config').HardhatUserConfig */
module.exports = {
	solidity: {
		version: "0.8.20",
		settings: {
			optimizer: { enabled: true, runs: 200 }
		}
	},
	paths: {
		sources: "contracts",
		tests: "test",
		cache: "cache",
		artifacts: "artifacts"
	}
}; 