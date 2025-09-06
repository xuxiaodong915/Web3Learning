# ERC20 代币合约项目

这是一个简单的 ERC20 代币合约实现，包含所有标准的 ERC20 功能。

## 📋 功能特性

- ✅ **标准 ERC20 功能**：
  - `balanceOf`: 查询账户余额
  - `transfer`: 转账功能
  - `approve`: 授权功能
  - `transferFrom`: 代扣转账功能
  - `allowance`: 查询授权额度

- ✅ **事件记录**：
  - `Transfer`: 转账事件
  - `Approval`: 授权事件

- ✅ **额外功能**：
  - `mint`: 增发代币（仅合约所有者）

## 🚀 快速开始

### 1. 安装依赖

```bash
npm install
```

### 2. 配置环境

在部署和交互脚本中，你需要替换以下内容：

- **私钥**: 将 `YOUR_PRIVATE_KEY_HERE` 替换为你的钱包私钥
- **合约地址**: 部署后会自动保存到 `deployment.json`

### 3. 部署合约

```bash
node deploy.js
```

部署成功后，你会看到：
- 合约地址
- 部署交易哈希
- 合约基本信息
- 部署信息会保存到 `deployment.json`

### 4. 与合约交互

```bash
node interact.js
```

## 📁 项目结构

```
ether/
├── contracts/
│   └── MyToken.sol          # ERC20 代币合约
├── deploy.js               # 部署脚本
├── interact.js             # 交互脚本
├── index.js                # 原有的事件监听脚本
├── package.json            # 项目配置
└── README.md               # 项目说明
```

## 🔧 合约功能详解

### 构造函数参数

```solidity
constructor(
    string memory _name,      // 代币名称
    string memory _symbol,    // 代币符号
    uint8 _decimals,         // 小数位数
    uint256 _initialSupply   // 初始供应量
)
```

### 主要函数

#### 查询函数
- `name()`: 获取代币名称
- `symbol()`: 获取代币符号
- `decimals()`: 获取小数位数
- `totalSupply()`: 获取总供应量
- `balanceOf(address account)`: 查询账户余额
- `allowance(address owner, address spender)`: 查询授权额度

#### 交易函数
- `transfer(address to, uint256 amount)`: 转账
- `approve(address spender, uint256 amount)`: 授权
- `transferFrom(address from, address to, uint256 amount)`: 代扣转账
- `mint(address to, uint256 amount)`: 增发代币（仅所有者）

### 事件

- `Transfer(address indexed from, address indexed to, uint256 value)`: 转账事件
- `Approval(address indexed owner, address indexed spender, uint256 value)`: 授权事件

## 💰 获取测试网 ETH

在 Sepolia 测试网上部署和测试合约需要测试网 ETH：

1. **Sepolia 水龙头**：
   - [Alchemy Sepolia Faucet](https://sepoliafaucet.com/)
   - [Infura Sepolia Faucet](https://www.infura.io/faucet/sepolia)

2. **获取步骤**：
   - 连接你的钱包
   - 输入你的钱包地址
   - 点击获取测试网 ETH

## 🔐 安全注意事项

⚠️ **重要安全提醒**：

1. **私钥安全**：
   - 永远不要在代码中硬编码私钥
   - 使用环境变量或配置文件存储私钥
   - 不要将包含私钥的文件提交到版本控制

2. **测试网使用**：
   - 本项目配置为 Sepolia 测试网
   - 不要在主网上测试未经验证的合约

3. **合约验证**：
   - 部署后建议在 Etherscan 上验证合约代码
   - 确保合约功能符合预期

## 📝 使用示例

### 部署合约

```javascript
// 部署参数
const name = "MyToken";
const symbol = "MTK";
const decimals = 18;
const initialSupply = 1000000; // 100万代币
```

### 转账操作

```javascript
// 转账 100 个代币
await interactor.transfer("0x...", 100);
```

### 授权操作

```javascript
// 授权 50 个代币给指定地址
await interactor.approve("0x...", 50);
```

### 代扣转账

```javascript
// 从授权方转账 30 个代币
await interactor.transferFrom("0x...", "0x...", 30);
```

### 增发代币

```javascript
// 为指定地址增发 1000 个代币
await interactor.mint("0x...", 1000);
```

## 🔍 在钱包中导入代币

部署成功后，你可以在 MetaMask 等钱包中导入代币：

1. **打开 MetaMask**
2. **选择 Sepolia 测试网**
3. **点击"导入代币"**
4. **输入合约地址**（从 `deployment.json` 获取）
5. **确认导入**

## 📊 监控合约

### 在 Etherscan 查看

部署后，你可以在 Sepolia Etherscan 查看合约：
```
https://sepolia.etherscan.io/address/YOUR_CONTRACT_ADDRESS
```

### 监听事件

```javascript
// 监听转账事件
interactor.listenToTransfers();

// 监听授权事件
interactor.listenToApprovals();
```

## 🛠️ 故障排除

### 常见问题

1. **部署失败**：
   - 检查私钥是否正确
   - 确保账户有足够的 Sepolia ETH
   - 检查网络连接

2. **交易失败**：
   - 检查余额是否充足
   - 确认授权额度
   - 验证地址格式

3. **合约交互失败**：
   - 确认合约地址正确
   - 检查 ABI 是否匹配
   - 验证网络配置

## 📄 许可证

MIT License

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

---

**注意**: 这是一个教育项目，用于学习 ERC20 标准。在生产环境中使用前，请进行充分的测试和安全审计。 