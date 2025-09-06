const { ethers } = require("ethers");
const fs = require("fs");
const path = require("path");

// 读取合约源码
const contractPath = path.join(__dirname, "contracts", "MyToken.sol");
const contractSource = fs.readFileSync(contractPath, "utf8");

// 编译合约（这里使用简单的编译，实际项目中建议使用 Hardhat 或 Truffle）
async function deployContract() {
    try {
        // 连接到 Sepolia 测试网
        const provider = new ethers.JsonRpcProvider("https://eth-sepolia.g.alchemy.com/v2/C1eDetxrjcCL9HMgWActr");
        
        // 替换为你的私钥（请确保私钥安全，不要提交到版本控制）
        const privateKey = "YOUR_PRIVATE_KEY_HERE"; // 请替换为你的私钥
        const wallet = new ethers.Wallet(privateKey, provider);
        
        console.log("部署账户地址:", wallet.address);
        
        // 合约构造函数参数
        const name = "MyToken";
        const symbol = "MTK";
        const decimals = 18;
        const initialSupply = 1000000; // 100万代币
        
        // 创建合约工厂
        const contractFactory = new ethers.ContractFactory(
            [
                "constructor(string memory _name, string memory _symbol, uint8 _decimals, uint256 _initialSupply)",
                "function name() public view returns (string memory)",
                "function symbol() public view returns (string memory)",
                "function decimals() public view returns (uint8)",
                "function totalSupply() public view returns (uint256)",
                "function balanceOf(address account) public view returns (uint256)",
                "function transfer(address to, uint256 amount) public returns (bool)",
                "function approve(address spender, uint256 amount) public returns (bool)",
                "function transferFrom(address from, address to, uint256 amount) public returns (bool)",
                "function allowance(address owner, address spender) public view returns (uint256)",
                "function mint(address to, uint256 amount) public",
                "function owner() public view returns (address)",
                "event Transfer(address indexed from, address indexed to, uint256 value)",
                "event Approval(address indexed owner, address indexed spender, uint256 value)"
            ],
            contractSource,
            wallet
        );
        
        console.log("开始部署合约...");
        
        // 部署合约
        const contract = await contractFactory.deploy(name, symbol, decimals, initialSupply);
        await contract.waitForDeployment();
        
        const contractAddress = await contract.getAddress();
        console.log("✅ 合约部署成功！");
        console.log("合约地址:", contractAddress);
        console.log("部署交易哈希:", contract.deploymentTransaction().hash);
        
        // 验证合约信息
        console.log("\n📋 合约信息:");
        console.log("代币名称:", await contract.name());
        console.log("代币符号:", await contract.symbol());
        console.log("小数位数:", await contract.decimals());
        console.log("总供应量:", ethers.formatUnits(await contract.totalSupply(), decimals));
        console.log("合约所有者:", await contract.owner());
        
        // 保存部署信息到文件
        const deploymentInfo = {
            contractAddress: contractAddress,
            deployer: wallet.address,
            name: name,
            symbol: symbol,
            decimals: decimals,
            initialSupply: initialSupply,
            deploymentTx: contract.deploymentTransaction().hash,
            network: "Sepolia"
        };
        
        fs.writeFileSync(
            path.join(__dirname, "deployment.json"),
            JSON.stringify(deploymentInfo, null, 2)
        );
        
        console.log("\n💾 部署信息已保存到 deployment.json");
        console.log("\n🔗 在 Sepolia Etherscan 查看合约:");
        console.log(`https://sepolia.etherscan.io/address/${contractAddress}`);
        
        return contractAddress;
        
    } catch (error) {
        console.error("❌ 部署失败:", error.message);
        throw error;
    }
}

// 如果直接运行此脚本
if (require.main === module) {
    deployContract()
        .then(() => process.exit(0))
        .catch((error) => {
            console.error(error);
            process.exit(1);
        });
}

module.exports = { deployContract }; 