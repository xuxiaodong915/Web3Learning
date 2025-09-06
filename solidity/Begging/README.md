# BeggingContract 项目说明

一个最小可运行的捐赠合约示例，包含：
- `BeggingContract` 合约（记录捐赠、可捐赠、仅所有者可提款、查询捐赠）。
- Hardhat 开发环境与脚本：编译、测试、本地部署。

## 环境要求
- Node.js >= 18（建议 18/20 LTS）
- npm >= 8
- 操作系统：Windows / macOS / Linux（示例命令对 Windows PowerShell 兼容）

## 快速开始（基于当前仓库）
> 已包含所有必要文件：`contracts/`、`hardhat.config.js`、`scripts/deploy.js`、`test/BeggingContract.test.js`、`package.json`。

1) 安装依赖
```bash
npm i
```

2) 运行测试（会自动编译）
```bash
npm test
```
期待输出中包含所有用例通过（绿色）。

3) 本地节点与部署
- 启动本地节点（保持该终端不关闭）：
```bash
npm run node
```
- 新开一个终端，在本地网络部署：
```bash
npm run deploy -- --network localhost
```
部署完成后会打印部署地址，例如：
```
BeggingContract deployed at: 0x...
```

## 从零创建项目（完整复现步骤）
若你希望从空目录自行搭建，按照以下步骤操作：

1) 新建目录并进入
```bash
mkdir taofan && cd taofan
```

2) 初始化 npm 项目
```bash
npm init -y
```

3) 安装最小必要依赖
```bash
npm i -D hardhat @nomicfoundation/hardhat-ethers ethers chai
```

4) 新建 Hardhat 配置 `hardhat.config.js`
```javascript
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
```

5) 新建合约 `contracts/BeggingContract.sol`
```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

contract BeggingContract {
	address public owner;
	mapping(address => uint256) private donations;

	constructor() {
		owner = msg.sender;
	}

	modifier onlyOwner() {
		require(msg.sender == owner, "Not owner");
		_;
	}

	function donate() external payable {
		require(msg.value > 0, "No ETH sent");
		donations[msg.sender] += msg.value;
	}

	function withdraw() external onlyOwner {
		uint256 amount = address(this).balance;
		require(amount > 0, "No balance");
		payable(owner).transfer(amount);
	}

	function getDonation(address donor) external view returns (uint256) {
		return donations[donor];
	}

	receive() external payable {
		donations[msg.sender] += msg.value;
	}
}
```

6) 新建部署脚本 `scripts/deploy.js`
```javascript
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
```

7) 新建测试 `test/BeggingContract.test.js`
```javascript
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
```

8) 添加 `package.json` 脚本（如果未自动创建）
```json
{
	"scripts": {
		"test": "hardhat test",
		"node": "hardhat node",
		"deploy": "hardhat run scripts/deploy.js"
	}
}
```

9) 安装并运行
- 安装依赖：
```bash
npm i
```
- 运行测试：
```bash
npm test
```
- 启动本地节点 + 部署：
```bash
npm run node
# 新终端
npm run deploy -- --network localhost
```

## 目录结构
```
contracts/
  BeggingContract.sol
scripts/
  deploy.js
test/
  BeggingContract.test.js
hardhat.config.js
package.json
README.md
```

## 常见问题（FAQ）
- Q: 测试时提示 `TypeError: ... changeEtherBalances is not a function`？
  - A: 确保使用的是 `chai` 与 Hardhat 内置的 `chai` 扩展（通过 `@nomicfoundation/hardhat-ethers` 已自动集成）。如果自定义了测试框架初始化，需确保正确引入 Hardhat 环境。
- Q: 版本兼容性？
  - A: 本项目使用 Solidity `^0.8.20`、`hardhat ^2.22.x`、`ethers ^6.x`。

## 许可证
MIT 