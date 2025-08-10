pragma solidity >=0.4.22 <0.7.0;

contract EventExample {
    // 1. 基本事件 - 记录用户注册
    event UserRegistered(address user, string name, uint256 timestamp);
    
    // 2. 带indexed参数的事件 - 便于过滤
    event Transfer(address indexed from, address indexed to, uint256 amount);
    
    // 3. 复杂事件 - 记录多个数据
    event OrderPlaced(
        uint256 indexed orderId,
        address indexed customer,
        string product,
        uint256 price,
        uint256 timestamp
    );
    
    // 4. 错误事件 - 记录异常情况
    event ErrorOccurred(string message, address user, uint256 timestamp);
    
    // 状态变量
    mapping(address => string) public users;
    uint256 public orderCounter;
    
    // 用户注册函数
    function registerUser(string memory name) public {
        require(bytes(name).length > 0, "用户名不能为空");
        users[msg.sender] = name;
        
        // 触发注册事件
        emit UserRegistered(msg.sender, name, block.timestamp);
    }
    
    // 转账函数
    function transfer(address to, uint256 amount) public {
        require(to != address(0), "无效地址");
        require(amount > 0, "金额必须大于0");
        
        // 这里应该有实际的转账逻辑
        // 触发转账事件
        emit Transfer(msg.sender, to, amount);
    }
    
    // 下单函数
    function placeOrder(string memory product, uint256 price) public {
        require(bytes(product).length > 0, "产品名称不能为空");
        require(price > 0, "价格必须大于0");
        
        orderCounter++;
        
        // 触发下单事件
        emit OrderPlaced(
            orderCounter,
            msg.sender,
            product,
            price,
            block.timestamp
        );
    }
    
    // 错误处理函数
    function triggerError() public {
        // 模拟错误情况
        emit ErrorOccurred("这是一个测试错误", msg.sender, block.timestamp);
    }
} 