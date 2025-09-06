// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/token/ERC20/extensions/ERC20Burnable.sol";
import "@openzeppelin/contracts/security/ReentrancyGuard.sol";
import "@openzeppelin/contracts/utils/Address.sol";

interface IUniswapV2Factory {
    function createPair(
        address tokenA,
        address tokenB
    ) external returns (address pair);
}

interface IUniswapV2Router02 {
    function factory() external view returns (address);
    function WETH() external view returns (address);

    function addLiquidityETH(
        address token,
        uint amountTokenDesired,
        uint amountTokenMin,
        uint amountETHMin,
        address to,
        uint deadline
    )
        external
        payable
        returns (uint amountToken, uint amountETH, uint liquidity);

    function removeLiquidityETH(
        address token,
        uint liquidity,
        uint amountTokenMin,
        uint amountETHMin,
        address to,
        uint deadline
    ) external returns (uint amountToken, uint amountETH);
}

/**
 * Shib 风格 Meme Token，包含：
 * - 交易税：每次转账按比例收取，转入税费钱包
 * - 流动性池交互：添加/移除 ETH 流动性（UniswapV2 风格路由）
 * - 交易限制：单笔最大金额、每日交易次数限制
 * - 详细注释：解释主要变量与函数
 */
contract MemeToken is ERC20, ERC20Burnable, Ownable, ReentrancyGuard {
    using Address for address;

    // ============ 税费与限制配置 ============
    uint256 public taxBasisPoints; // 税率，万分制。例如 300 = 3%
    address public taxWallet; // 税费接收地址
    mapping(address => bool) public isExcludedFromFee; // 免税地址

    uint256 public maxTxAmount; // 单笔交易最大额度（以 token 为单位）
    uint256 public dailyTxLimit; // 每日最多交易次数
    struct DailyStat {
        uint256 lastReset;
        uint256 count;
    }
    mapping(address => DailyStat) private _dailyStats; // 记录地址的当日交易次数

    // ============ 交易开关与路由/LP ============
    bool public tradingEnabled; // 交易是否开启
    IUniswapV2Router02 public dexRouter; // 路由器
    address public dexPair; // 与 WETH 的交易对

    event TaxUpdated(uint256 newBasisPoints, address taxWallet);
    event TradingEnabled();
    event RouterUpdated(address router, address pair);
    event LimitsUpdated(uint256 maxTxAmount, uint256 dailyTxLimit);

    constructor(
        string memory name_,
        string memory symbol_,
        uint256 initialSupply_,
        address router_,
        address taxWallet_,
        uint256 taxBps_,
        uint256 maxTxAmount_,
        uint256 dailyTxLimit_
    ) ERC20(name_, symbol_) {
        require(taxWallet_ != address(0), "tax wallet zero");
        taxWallet = taxWallet_;
        taxBasisPoints = taxBps_;
        maxTxAmount = maxTxAmount_;
        dailyTxLimit = dailyTxLimit_;

        if (router_ != address(0)) {
            _setRouterAndPair(router_);
        }

        // 初始铸造给部署者
        _mint(msg.sender, initialSupply_);
        isExcludedFromFee[msg.sender] = true;
        isExcludedFromFee[address(this)] = true;
        isExcludedFromFee[taxWallet] = true;
    }

    // ============ 仅管理员：参数设置 ============
    function setTax(uint256 bps, address wallet) external onlyOwner {
        require(wallet != address(0), "tax wallet zero");
        require(bps <= 1000, "tax too high"); // 上限 10%
        taxBasisPoints = bps;
        taxWallet = wallet;
        emit TaxUpdated(bps, wallet);
    }

    function setLimits(
        uint256 newMaxTxAmount,
        uint256 newDailyTxLimit
    ) external onlyOwner {
        require(newMaxTxAmount > 0, "maxTx=0");
        maxTxAmount = newMaxTxAmount;
        dailyTxLimit = newDailyTxLimit;
        emit LimitsUpdated(newMaxTxAmount, newDailyTxLimit);
    }

    function setExcludedFromFee(
        address account,
        bool excluded
    ) external onlyOwner {
        isExcludedFromFee[account] = excluded;
    }

    function setRouter(address router) external onlyOwner {
        _setRouterAndPair(router);
    }

    function enableTrading() external onlyOwner {
        require(!tradingEnabled, "already enabled");
        tradingEnabled = true;
        emit TradingEnabled();
    }

    function _setRouterAndPair(address router) internal {
        require(router != address(0), "router zero");
        dexRouter = IUniswapV2Router02(router);
        address factory = dexRouter.factory();
        address weth = dexRouter.WETH();
        require(factory != address(0) && weth != address(0), "router invalid");
        dexPair = IUniswapV2Factory(factory).createPair(address(this), weth);
        emit RouterUpdated(router, dexPair);
    }

    // ============ 转账逻辑（含税与限制） ============
    function _transfer(
        address from,
        address to,
        uint256 amount
    ) internal override {
        // 交易开关：排除 mint/burn 与合约内部流转
        if (
            from != address(0) &&
            to != address(0) &&
            !isExcludedFromFee[from] &&
            !isExcludedFromFee[to]
        ) {
            require(tradingEnabled, "trading disabled");
            // 单笔限制
            require(amount <= maxTxAmount, "exceeds maxTxAmount");
            // 每日次数限制（以发送方为准）
            if (dailyTxLimit > 0) {
                DailyStat storage st = _dailyStats[from];
                uint256 dayStart = block.timestamp - (block.timestamp % 1 days);
                if (st.lastReset < dayStart) {
                    st.lastReset = dayStart;
                    st.count = 0;
                }
                require(st.count < dailyTxLimit, "daily tx limit");
                st.count += 1;
            }
        }

        uint256 sendAmount = amount;
        uint256 fee;
        if (
            !isExcludedFromFee[from] &&
            !isExcludedFromFee[to] &&
            taxBasisPoints > 0
        ) {
            fee = (amount * taxBasisPoints) / 10000;
            if (fee > 0) {
                super._transfer(from, taxWallet, fee);
                sendAmount = amount - fee;
            }
        }

        super._transfer(from, to, sendAmount);
    }

    // ============ LP 交互（需要先 approve 给路由） ============
    /**
     * 向路由添加 ETH 流动性。调用前需：
     * - 拥有者或用户先 approve 路由合约 amountTokenDesired 数量
     * - 通过 msg.value 发送 ETH
     */
    function addLiquidityETH(
        uint256 amountTokenDesired,
        uint256 amountTokenMin,
        uint256 amountETHMin,
        uint256 deadline
    )
        external
        payable
        nonReentrant
        returns (uint256 amountToken, uint256 amountETH, uint256 liquidity)
    {
        require(address(dexRouter) != address(0), "router unset");
        (amountToken, amountETH, liquidity) = dexRouter.addLiquidityETH{
            value: msg.value
        }(
            address(this),
            amountTokenDesired,
            amountTokenMin,
            amountETHMin,
            msg.sender,
            deadline
        );
    }

    /**
     * 移除 ETH 流动性。调用前需：
     * - 用户持有 LP 代币并已 approve 路由
     */
    function removeLiquidityETH(
        uint256 liquidity,
        uint256 amountTokenMin,
        uint256 amountETHMin,
        uint256 deadline
    ) external nonReentrant returns (uint256 amountToken, uint256 amountETH) {
        require(address(dexRouter) != address(0), "router unset");
        (amountToken, amountETH) = dexRouter.removeLiquidityETH(
            address(this),
            liquidity,
            amountTokenMin,
            amountETHMin,
            msg.sender,
            deadline
        );
    }

    receive() external payable {}
}
