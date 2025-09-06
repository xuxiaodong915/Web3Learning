const { ethers } = require("ethers");
const { TokenInteractor } = require("./interact");

// 测试配置
const TEST_PRIVATE_KEY = "YOUR_PRIVATE_KEY_HERE"; // 请替换为你的私钥
const TEST_RECIPIENT = "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6"; // 测试接收地址

async function runTests() {
    console.log("🧪 开始 ERC20 合约测试...\n");
    
    try {
        // 读取部署信息
        const fs = require("fs");
        const deploymentInfo = JSON.parse(fs.readFileSync("deployment.json", "utf8"));
        const contractAddress = deploymentInfo.contractAddress;
        
        console.log("📋 测试配置:");
        console.log("合约地址:", contractAddress);
        console.log("测试网络: Sepolia");
        console.log("---\n");
        
        // 创建交互器
        const interactor = new TokenInteractor(contractAddress, TEST_PRIVATE_KEY);
        
        // 测试 1: 获取合约信息
        console.log("✅ 测试 1: 获取合约信息");
        await interactor.getContractInfo();
        console.log("---\n");
        
        // 测试 2: 查询余额
        console.log("✅ 测试 2: 查询余额");
        await interactor.getBalance(interactor.wallet.address);
        console.log("---\n");
        
        // 测试 3: 转账（小额测试）
        console.log("✅ 测试 3: 转账功能");
        const transferAmount = 10; // 转账 10 个代币
        await interactor.transfer(TEST_RECIPIENT, transferAmount);
        
        // 验证转账结果
        console.log("验证转账结果:");
        await interactor.getBalance(interactor.wallet.address);
        await interactor.getBalance(TEST_RECIPIENT);
        console.log("---\n");
        
        // 测试 4: 授权功能
        console.log("✅ 测试 4: 授权功能");
        const approveAmount = 50; // 授权 50 个代币
        await interactor.approve(TEST_RECIPIENT, approveAmount);
        
        // 验证授权结果
        await interactor.getAllowance(interactor.wallet.address, TEST_RECIPIENT);
        console.log("---\n");
        
        // 测试 5: 增发代币（仅合约所有者）
        console.log("✅ 测试 5: 增发代币");
        const mintAmount = 100; // 增发 100 个代币
        await interactor.mint(interactor.wallet.address, mintAmount);
        
        // 验证增发结果
        console.log("验证增发结果:");
        await interactor.getBalance(interactor.wallet.address);
        console.log("---\n");
        
        // 测试 6: 监听事件（可选）
        console.log("✅ 测试 6: 事件监听");
        console.log("开始监听转账和授权事件...");
        console.log("（按 Ctrl+C 停止监听）");
        
        interactor.listenToTransfers();
        interactor.listenToApprovals();
        
        // 触发一些事件
        setTimeout(async () => {
            try {
                await interactor.transfer(TEST_RECIPIENT, 5);
                await interactor.approve(TEST_RECIPIENT, 25);
            } catch (error) {
                console.log("事件触发测试完成");
            }
        }, 2000);
        
    } catch (error) {
        console.error("❌ 测试失败:", error.message);
        
        if (error.message.includes("YOUR_PRIVATE_KEY_HERE")) {
            console.log("\n💡 提示: 请先在 test.js 中设置你的私钥");
        } else if (error.message.includes("deployment.json")) {
            console.log("\n💡 提示: 请先运行 node deploy.js 部署合约");
        }
    }
}

// 运行测试
if (require.main === module) {
    runTests()
        .then(() => {
            console.log("\n🎉 测试完成！");
            console.log("注意: 事件监听会持续运行，按 Ctrl+C 停止");
        })
        .catch((error) => {
            console.error("测试过程中发生错误:", error);
            process.exit(1);
        });
}

module.exports = { runTests }; 