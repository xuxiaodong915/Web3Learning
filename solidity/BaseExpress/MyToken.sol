// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

/**
 * @title MyToken
 * @dev 一个简单的 ERC20 代币合约
    任务：参考 openzeppelin-contracts/contracts/token/ERC20/IERC20.sol实现一个简单的 ERC20 代币合约。要求：
    合约包含以下标准 ERC20 功能：
    balanceOf：查询账户余额。
    transfer：转账。
    approve 和 transferFrom：授权和代扣转账。
    使用 event 记录转账和授权操作。
    提供 mint 函数，允许合约所有者增发代币。
    提示：
    使用 mapping 存储账户余额和授权信息。
    使用 event 定义 Transfer 和 Approval 事件。
    部署到sepolia 测试网，导入到自己的钱包
    */
contract MyToken {
    string public name;
    string public symbol;
    uint8 public decimals;
    uint256 public totalSupply;
    
    address public owner;
    
    mapping(address => uint256) public balanceOf;
    mapping(address => mapping(address => uint256)) public allowance;
    
    event Transfer(address indexed from, address indexed to, uint256 value);
    event Approval(address indexed owner, address indexed spender, uint256 value);
    
    modifier onlyOwner() {
        require(msg.sender == owner, "Only the contract owner can call this function.");
        _;
    }
    
    constructor(string memory _name, string memory _symbol, uint8 _decimals, uint256 _initialSupply) {
        name = _name;
        symbol = _symbol;
        decimals = _decimals;
        totalSupply = _initialSupply * 10**_decimals;
        owner = msg.sender;
        balanceOf[msg.sender] = totalSupply;
        emit Transfer(address(0), msg.sender, totalSupply);
    }
    
    /**
     * @dev 查询账户余额
     * @param account 要查询的账户地址
     * @return 账户余额
     */
    function balanceOf(address account) public view returns (uint256) {
        return balanceOf[account];
    }
    
    /**
     * @dev 转账函数
     * @param to 接收方地址
     * @param amount 转账金额
     * @return 是否成功
     */
    function transfer(address to, uint256 amount) public returns (bool) {
        require(to != address(0), "Transfers cannot be made to the zero address.");
        require(balanceOf[msg.sender] >= amount, "No balance");
        
        balanceOf[msg.sender] -= amount;
        balanceOf[to] += amount;
        
        emit Transfer(msg.sender, to, amount);
        return true;
    }
    
    /**
     * @dev 授权函数
     * @param spender 被授权地址
     * @param amount 授权金额
     * @return 是否成功
     */
    function approve(address spender, uint256 amount) public returns (bool) {
        require(spender != address(0), "Authorization cannot be granted to the zero address.");
        
        allowance[msg.sender][spender] = amount;
        emit Approval(msg.sender, spender, amount);
        return true;
    }
    
    /**
     * @dev 代扣转账函数
     * @param from 发送方地址
     * @param to 接收方地址
     * @param amount 转账金额
     * @return 是否成功
     */
    function transferFrom(address from, address to, uint256 amount) public returns (bool) {
        require(to != address(0), "Transfers cannot be made to zero-addresses.");
        require(balanceOf[from] >= amount, "No balance");
        require(allowance[from][msg.sender] >= amount, "Insufficient authorization limit");
        
        balanceOf[from] -= amount;
        balanceOf[to] += amount;
        allowance[from][msg.sender] -= amount;
        
        emit Transfer(from, to, amount);
        return true;
    }
    
    /**
     * @dev 增发代币函数（仅合约所有者可调用）
     * @param to 接收方地址
     * @param amount 增发金额
     */
    function mint(address to, uint256 amount) public onlyOwner {
        require(to != address(0), "Cannot issue to zero address");
        
        totalSupply += amount;
        balanceOf[to] += amount;
        
        emit Transfer(address(0), to, amount);
    }
    
    /**
     * @dev 查询授权额度
     * @param owner 授权方地址
     * @param spender 被授权方地址
     * @return 授权额度
     */
    function allowance(address owner, address spender) public view returns (uint256) {
        return allowance[owner][spender];
    }
} 