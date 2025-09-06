const { ethers } = require("ethers");
const fs = require("fs");

// 合约 ABI
const contractABI = [
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
];

class TokenInteractor {
    constructor(contractAddress, privateKey) {
        this.provider = new ethers.JsonRpcProvider("https://eth-sepolia.g.alchemy.com/v2/C1eDetxrjcCL9HMgWActr");
        this.wallet = new ethers.Wallet(privateKey, this.provider);
        this.contract = new ethers.Contract(contractAddress, contractABI, this.wallet);
        this.contractAddress = contractAddress;
    }
    
    // 获取合约基本信息
    async getContractInfo() {
        console.log("📋 合约信息:");
        console.log("合约地址:", this.contractAddress);
        console.log("代币名称:", await this.contract.name());
        console.log("代币符号:", await this.contract.symbol());
        console.log("小数位数:", await this.contract.decimals());
        console.log("总供应量:", ethers.formatUnits(await this.contract.totalSupply(), await this.contract.decimals()));
        console.log("合约所有者:", await this.contract.owner());
        console.log("当前账户:", this.wallet.address);
        console.log("当前余额:", ethers.formatUnits(await this.contract.balanceOf(this.wallet.address), await this.contract.decimals()));
    }
    
    // 查询余额
    async getBalance(address) {
        const balance = await this.contract.balanceOf(address);
        const decimals = await this.contract.decimals();
        console.log(`💰 ${address} 的余额: ${ethers.formatUnits(balance, decimals)} ${await this.contract.symbol()}`);
        return balance;
    }
    
    // 转账
    async transfer(to, amount) {
        try {
            const decimals = await this.contract.decimals();
            const amountWei = ethers.parseUnits(amount.toString(), decimals);
            
            console.log(`🔄 正在转账 ${amount} ${await this.contract.symbol()} 到 ${to}...`);
            
            const tx = await this.contract.transfer(to, amountWei);
            await tx.wait();
            
            console.log("✅ 转账成功！");
            console.log("交易哈希:", tx.hash);
            
            return tx;
        } catch (error) {
            console.error("❌ 转账失败:", error.message);
            throw error;
        }
    }
    
    // 授权
    async approve(spender, amount) {
        try {
            const decimals = await this.contract.decimals();
            const amountWei = ethers.parseUnits(amount.toString(), decimals);
            
            console.log(`🔐 正在授权 ${spender} 使用 ${amount} ${await this.contract.symbol()}...`);
            
            const tx = await this.contract.approve(spender, amountWei);
            await tx.wait();
            
            console.log("✅ 授权成功！");
            console.log("交易哈希:", tx.hash);
            
            return tx;
        } catch (error) {
            console.error("❌ 授权失败:", error.message);
            throw error;
        }
    }
    
    // 查询授权额度
    async getAllowance(owner, spender) {
        const allowance = await this.contract.allowance(owner, spender);
        const decimals = await this.contract.decimals();
        console.log(`🔍 ${owner} 授权给 ${spender} 的额度: ${ethers.formatUnits(allowance, decimals)} ${await this.contract.symbol()}`);
        return allowance;
    }
    
    // 代扣转账
    async transferFrom(from, to, amount) {
        try {
            const decimals = await this.contract.decimals();
            const amountWei = ethers.parseUnits(amount.toString(), decimals);
            
            console.log(`🔄 正在从 ${from} 转账 ${amount} ${await this.contract.symbol()} 到 ${to}...`);
            
            const tx = await this.contract.transferFrom(from, to, amountWei);
            await tx.wait();
            
            console.log("✅ 代扣转账成功！");
            console.log("交易哈希:", tx.hash);
            
            return tx;
        } catch (error) {
            console.error("❌ 代扣转账失败:", error.message);
            throw error;
        }
    }
    
    // 增发代币（仅合约所有者）
    async mint(to, amount) {
        try {
            const decimals = await this.contract.decimals();
            const amountWei = ethers.parseUnits(amount.toString(), decimals);
            
            console.log(`🪙 正在为 ${to} 增发 ${amount} ${await this.contract.symbol()}...`);
            
            const tx = await this.contract.mint(to, amountWei);
            await tx.wait();
            
            console.log("✅ 增发成功！");
            console.log("交易哈希:", tx.hash);
            
            return tx;
        } catch (error) {
            console.error("❌ 增发失败:", error.message);
            throw error;
        }
    }
    
    // 监听转账事件
    listenToTransfers() {
        console.log("👂 开始监听转账事件...");
        this.contract.on("Transfer", (from, to, value, event) => {
            const symbol = this.contract.symbol();
            const decimals = this.contract.decimals();
            console.log("📥 转账事件:");
            console.log("  从:", from);
            console.log("  到:", to);
            console.log("  金额:", ethers.formatUnits(value, decimals), symbol);
            console.log("  交易哈希:", event.transactionHash);
            console.log("---");
        });
    }
    
    // 监听授权事件
    listenToApprovals() {
        console.log("👂 开始监听授权事件...");
        this.contract.on("Approval", (owner, spender, value, event) => {
            const symbol = this.contract.symbol();
            const decimals = this.contract.decimals();
            console.log("🔐 授权事件:");
            console.log("  授权方:", owner);
            console.log("  被授权方:", spender);
            console.log("  授权金额:", ethers.formatUnits(value, decimals), symbol);
            console.log("  交易哈希:", event.transactionHash);
            console.log("---");
        });
    }
}

// 示例使用
async function main() {
    // 从 deployment.json 读取合约地址
    let contractAddress;
    try {
        const deploymentInfo = JSON.parse(fs.readFileSync("deployment.json", "utf8"));
        contractAddress = deploymentInfo.contractAddress;
    } catch (error) {
        console.error("❌ 请先部署合约或手动指定合约地址");
        return;
    }
    
    // 替换为你的私钥
    const privateKey = "YOUR_PRIVATE_KEY_HERE"; // 请替换为你的私钥
    
    const interactor = new TokenInteractor(contractAddress, privateKey);
    
    // 获取合约信息
    await interactor.getContractInfo();
    
    // 示例操作（取消注释以执行）
    
    // 1. 查询余额
    // await interactor.getBalance(interactor.wallet.address);
    
    // 2. 转账（替换为实际地址）
    // await interactor.transfer("0x...", 100);
    
    // 3. 授权（替换为实际地址）
    // await interactor.approve("0x...", 50);
    
    // 4. 查询授权额度
    // await interactor.getAllowance(interactor.wallet.address, "0x...");
    
    // 5. 增发代币（仅合约所有者）
    // await interactor.mint(interactor.wallet.address, 1000);
    
    // 6. 监听事件
    // interactor.listenToTransfers();
    // interactor.listenToApprovals();
}

if (require.main === module) {
    main()
        .then(() => process.exit(0))
        .catch((error) => {
            console.error(error);
            process.exit(1);
        });
}

module.exports = { TokenInteractor }; 