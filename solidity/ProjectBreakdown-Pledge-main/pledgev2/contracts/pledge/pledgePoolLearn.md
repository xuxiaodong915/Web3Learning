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