require("@nomicfoundation/hardhat-toolbox");
require("hardhat-deploy")
require("@openzeppelin/hardhat-upgrades")

/** @type import('hardhat/config').HardhatUserConfig */
module.exports = {
  solidity: "0.8.28",
  namedAccounts: {
    deployer: 0,
    user1: 1,
    user2: 2,
  },
  deploy: {
    // 配置部署脚本目录和文件过滤
    deployScripts: {
      // 只处理 .js 文件
      include: ["**/*.js"],
      // 排除 JSON 文件
      exclude: ["**/*.json"]
    }
  }
};
