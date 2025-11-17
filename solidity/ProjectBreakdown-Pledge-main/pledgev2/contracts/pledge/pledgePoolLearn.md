[3 tools called]

- **配置管理**：构造函数写入预言机、Swap Router、费用地址并默认费率为 0；管理员可用 `setFee`、`setSwapRouterAddress`、`setFeeAddress`、`setMinAmount` 动态调整费率、兑换路由及最小入池金额，所有函数都受多签 `validCall` 保护。  
```126:179:pledgev2/contracts/pledge/PledgePool.sol
constructor(...) multiSignatureClient(_multiSignature) { ... }
function setFee(uint256 _lendFee,uint256 _borrowFee) validCall external { ... }
function setSwapRouterAddress(address _swapRouter) validCall external { ... }
function setFeeAddress(address payable _feeAddress) validCall external { ... }
function setMinAmount(uint256 _minAmount) validCall external { ... }
```

- **池子建模**：`poolLength` 返回池子数量；`createPoolInfo` 依据结算/结束时间、利率、抵押率、借贷 token 等参数创建 `PoolBaseInfo` 与 `PoolDataInfo`；`getPoolState` 读取当前状态。  
```185:238:pledgev2/contracts/pledge/PledgePool.sol
function poolLength() external view returns (uint256) { ... }
function createPoolInfo(uint256 _settleTime, ... ) public validCall { ... }
function getPoolState(uint256 _pid) public view returns (uint256) { ... }
```

- **出借人流程**：  
  - `depositLend` 在 MATCH 阶段注入借出资金并累计 `lendSupply`。  
  - `refundLend` 在匹配结束后按比例退回未撮合部分。  
  - `claimLend` 领取对应 SP 份额。  
  - `withdrawLend` 在 FINISH/LIQUIDATION 阶段燃烧 SP 赎回本金+收益或清算款。  
  - `emergencyLendWithdrawal` 在 UNDONE 状态直接赎回本金。  
```246:371:pledgev2/contracts/pledge/PledgePool.sol
function depositLend(uint256 _pid, uint256 _stakeAmount) external payable ... { ... }
function refundLend(uint256 _pid) external ... { ... }
function claimLend(uint256 _pid) external ... { ... }
function withdrawLend(uint256 _pid, uint256 _spAmount)  external ... { ... }
function emergencyLendWithdrawal(uint256 _pid) external ... { ... }
```

- **借款人流程**：  
  - `depositBorrow` 在 MATCH 阶段质押抵押物。  
  - `refundBorrow` 领取未撮合的抵押物。  
  - `claimBorrow` 铸造 JP 份额并提取借到的稳定币。  
  - `withdrawBorrow` 在 FINISH/LIQUIDATION 阶段销毁 JP 赎回剩余抵押物。  
  - `emergencyBorrowWithdrawal` 在 UNDONE 状态全额赎回抵押。  
```378:512:pledgev2/contracts/pledge/PledgePool.sol
function depositBorrow(uint256 _pid, uint256 _stakeAmount ) external payable ... { ... }
function refundBorrow(uint256 _pid) external ... { ... }
function claimBorrow(uint256 _pid) external ... { ... }
function withdrawBorrow(uint256 _pid, uint256 _jpAmount ) external ... { ... }
function emergencyBorrowWithdrawal(uint256 _pid) external ... { ... }
```

- **生命周期控制**：  
  - `checkoutSettle/checkoutFinish` 仅做时间判断。  
  - `settle` 在撮合期结束后，根据 Oracle 报价及抵押率确定最终撮合额，并将状态置为 EXECUTION/UNDONE。  
  - `finish` 在到期时卖出抵押物偿还借款+利息、扣除手续费并记录 `finishAmount*`，状态改为 FINISH。  
  - `checkoutLiquidate` 判断抵押价值是否低于阈值；`liquidate` 触发清算，与 `finish` 类似但写入 `liquidationAmount*` 并将状态改为 LIQUIDATION。  
```518:685:pledgev2/contracts/pledge/PledgePool.sol
function checkoutSettle(uint256 _pid) public view returns(bool){ ... }
function settle(uint256 _pid) public validCall { ... }
function checkoutFinish(uint256 _pid) public view returns(bool){ ... }
function finish(uint256 _pid) public validCall { ... }
function checkoutLiquidate(uint256 _pid) external view returns(bool) { ... }
function liquidate(uint256 _pid) public validCall { ... }
```

- **内部工具与安全控制**：`redeemFees` 从偿付金额里扣费；`_sellExactAmount/_swap/_safeApprove` 处理 Pancake/Uniswap 路由兑换；`getUnderlyingPriceView` 读取 Oracle 报价；`setPause` 及一系列 modifier 控制暂停、时间窗口与状态机，配合 `ReentrancyGuard` 和 `SafeTransfer` 保障资金安全。  
```688:826:pledgev2/contracts/pledge/PledgePool.sol
function redeemFees(uint256 feeRatio,address token,uint256 amount) internal returns (uint256){ ... }
function _sellExactAmount(...) internal returns (uint256,uint256){ ... }
function _swap(...) internal returns (uint256) { ... }
function getUnderlyingPriceView(uint256 _pid) public view returns(uint256[2]memory){ ... }
function setPause() public validCall { ... }
modifier stateMatch(uint256 _pid) { ... }
...
```

综合整个项目代码，请问这些合约函数是如何调用的

[6 tools called]

这些 Solidity 函数本身并没有在 Go 后端里直接被调用；它们位于 `pledgev2/contracts/pledge/PledgePool.sol`，最终都是通过链上交易触发，调用方式可分两类：

- **管理员（多签）专属函数**  
  `setFee`、`setSwapRouterAddress`、`setFeeAddress`、`setMinAmount`、`createPoolInfo`、`settle`、`finish`、`liquidate`、`setPause` 等都带 `validCall` 修饰符。`validCall` 来源于 `multiSignatureClient`，内部通过 `IMultiSignature.getValidSignature` 校验多签批准：
  
```22:38:pledgev2/contracts/multiSignature/multiSignatureClient.sol
modifier validCall(){
    checkMultiSignature();
    _;
}
```

  也就是说，运营方要先在多签合约里收集足够签名，再由多签合约代表他们发起对 `PledgePool` 的交易，这些函数才会成功执行，和后端 API 没有直接关系。

- **普通用户函数**  
  像 `depositLend`、`refundLend`、`claimLend`、`withdrawLend`、`depositBorrow`、`refundBorrow`、`claimBorrow`、`withdrawBorrow`、`emergency*` 等都没有 `validCall`，任何地址都可在满足时间/状态限制、通过 `ReentrancyGuard` 检查后直接调用。这些场景通常由前端 DApp 或脚本（例如用户钱包的合约交互界面）发起交易完成，不需要后端参与。

后端 `pledge-backend` 的 REST 服务（`poolBaseInfo`、`poolDataInfo`、`price` 等，见 `pledge-backend/api/routes/route.go`）主要完成链下数据读写：它从 MySQL/Redis 中读取池信息、多签配置或行情价格，并不负责执行任何链上写操作。换句话说：

- 用户端或运营端 → 直接与链上合约交互，触发合约函数；
- 后端 API → 只提供查询、认证、多签账户配置等辅助信息，供前端展示或管理之用。

因此，“这些合约函数如何调用”可以概括为：管理员类函数由多签合约代表调用，用户类函数由钱包或脚本直接向 `PledgePool` 发送交易，而后端并未封装任何调用逻辑。

哪些前端代码调用的合约函数？

- `PoolServer` 实际封装了前端访问 `PledgePool` 合约的所有方法，包括只读的 `poolLength/poolBaseInfo/poolDataInfo` 以及交易类的 `depositLend/depositBorrow/claimLend/claimBorrow/refundLend/refundBorrow/withdrawLend/withdrawBorrow/emergencyLendWithdrawal/emergencyBorrowWithdrawal` 等；每个方法都会根据当前链 ID 选择合适的合约地址并发起 `contract.methods.xxx().call/send`。  
```11:138:pledge-fe/src/services/PoolServer.ts
async depositLend(pid, value, coinAddress, chainId) {
  const contract = getPledgePoolContract(...);
  let options = await gasOptions();
  if (coinAddress === '0x000...000') {
    options = { ...options, value };
  }
  return await contract.methods.depositLend(pid, value).send(options);
}
```

- 这些服务被各个界面直接引用，从而把具体的合约函数绑定到 UI 操作：  
  - `借出/借入` 面板在“Approve & Lend/Borrow”按钮里调用 `depositLend`、`depositBorrow`。  
```828:941:pledge-fe/src/components/Coin_pool/index.tsx
services.PoolServer.depositLend(pid, num, poolinfo[pid]?.Sp ?? 0, chainId)
...
services.PoolServer.depositBorrow(pid, borrownum, timestamp, poolinfo[pid]?.Jp ?? 0, chainId)
```
  - `AccessTab`（领取 SP/JP）根据 Tab 选择触发 `getclaimLend` 或 `getclaimBorrow`，并在成功后刷新 `userLendInfo/userBorrowInfo`。  
```229:256:pledge-fe/src/components/AccessTab/index.tsx
await services.PoolServer.getclaimLend((Number(props.key) - 1).toString(), chainId)
...
services.PoolServer.getuserLendInfo((Number(props.key) - 1).toString(), chainId)
```
  - `ClaimTime` 负责到期赎回本金，分别调用 `getwithdrawLend` 和 `getwithdrawBorrow`；同组件也会查询用户在池子的 SP/JP 余额。  
```233:263:pledge-fe/src/components/ClaimTime/index.tsx
services.PoolServer.getwithdrawLend(pid, Spnum, chainId)
...
services.PoolServer.getuserBorrowInfo(pid.toString(), chainId)
```
  - `Refund` 弹窗覆盖 `refundLend/refundBorrow` 与 `emergencyLendWithdrawal/emergencyBorrowWithdrawal` 的入口。  
```199:261:pledge-fe/src/components/Refund/index.tsx
services.PoolServer.getemergencyLendWithdrawal(props.key - 1, chainId)
...
services.PoolServer.getrefundBorrow(props.key - 1, chainId)
```
  - `PortfolioList`、`Market_Mode` 等页面在渲染个人仓位或市场列表时调用 `getuserLendInfo`、`getuserBorrowInfo` 来读取合约 `mapping` 中的用户状态。  
```47:111:pledge-fe/src/components/PortfolioList/index.tsx
services.PoolServer.getuserLendInfo((props.props.key - 1).toString(), chainId)
```

- 另外，所有 ERC20 的 `approve/allowance/balanceOf` 操作由 `ERC20Server` 封装，并在 `Coin_pool` 等组件的授权流程中调用，为 `depositLend/depositBorrow` 等函数提供必要的代币授权。  
```15:31:pledge-fe/src/services/ERC20Server.ts
await contract.methods.approve(...).send(options);
...
return await contract.methods.allowance(owner, ...).call();
```

综上，前端通过 `PoolServer`+`ERC20Server` 两层封装把 `PledgePool` 合约的核心函数全部接入 UI：用户在借出/借入、退款、到期赎回、紧急退出等界面进行的每一次交互，最终都对应着上述合约方法的调用。