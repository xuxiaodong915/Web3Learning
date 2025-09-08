# Shib 风格 Meme 代币合约

本项目基于 Hardhat 与 Solidity，实现具备交易税、流动性池交互与交易限制功能的 Meme 代币（SHIB 风格）。

## 功能特性
- 交易税：每笔转账按万分制收税，发送到税费钱包
- 流动性池：与 UniswapV2 风格路由交互，支持添加/移除 ETH 流动性
- 交易限制：单笔最大额度、每日交易次数限制，需开启交易开关后生效
- 详细注释：合约代码注释清晰，便于审阅与二次开发

## 环境准备
1. Node.js 18+
2. 安装依赖
```bash
npm install
```

## 配置
复制 `.env.example` 为 `.env` 并填写：
- `RPC_URL`: 网络 RPC（如 sepolia/alchemy/infura）
- `PRIVATE_KEY`: 部署私钥（测试账号）
- `ETHERSCAN_API_KEY`: 验证用，可选
- `TOKEN_NAME`,`TOKEN_SYMBOL`,`TOKEN_SUPPLY`
- `ROUTER_ADDRESS`: UniswapV2 路由器地址（测试网需对应测试路由）
- `TAX_WALLET`: 税费钱包地址
- 可选：`TAX_BPS`(默认300=3%), `MAX_TX`, `DAILY_TX_LIMIT`

## 编译
```bash
npx hardhat compile
```

## 部署
```bash
npx hardhat run scripts/deploy.js --network sepolia
```
记录输出的合约地址。

## 使用说明
- 开启交易：
	- 仅拥有者可调用 `enableTrading()`
- 设置税费：
	- `setTax(bps, wallet)`，上限 10%
- 设置限制：
	- `setLimits(maxTxAmount, dailyTxLimit)`
- 免税地址：
	- `setExcludedFromFee(account, true/false)`
- 路由/交易对：
	- 部署时可指定 `ROUTER_ADDRESS`，也可后续 `setRouter(router)` 自动创建与 WETH 的交易对

### 代币转账
- 常规 ERC20 `transfer/transferFrom`
- 若双方非免税地址且交易已开启，将收取税费并转入税费钱包

### 添加流动性（ETH）
1. 将代币授权给路由（钱包或脚本）：
	- `approve(router, amount)`
2. 调用合约 `addLiquidityETH(amountTokenDesired, amountTokenMin, amountETHMin, deadline)` 并附带 `msg.value=ETH`
3. 返回 `(amountToken, amountETH, liquidity)`

### 移除流动性（ETH）
1. 在路由上授权 LP 代币
2. 调用 `removeLiquidityETH(liquidity, amountTokenMin, amountETHMin, deadline)`

## 注意事项
- 部署后默认交易未开启，需 `enableTrading()`
- 税费免除默认包含：部署者、合约自身、税费钱包
- 添加/移除流动性前需正确授权
- 本合约供学习/测试，生产环境请进行安全审计 