pragma solidity >=0.4.22 <0.7.0;

contract BytesExample {
    
    // bytes1 示例
    bytes1 public singleByte = 0x41; // 字母 'A' 的ASCII码
    
    // 不同字节类型的对比
    bytes1 public byte1 = 0xFF;        // 1字节
    bytes2 public byte2 = 0x1234;      // 2字节
    bytes4 public byte4 = 0x12345678;  // 4字节
    bytes32 public byte32;             // 32字节
    
    // 字符串转bytes1数组
    function stringToBytes1Array(string memory input) public pure returns (bytes1[] memory) {
        bytes memory strBytes = bytes(input);
        bytes1[] memory result = new bytes1[](strBytes.length);
        
        for (uint i = 0; i < strBytes.length; i++) {
            result[i] = bytes1(strBytes[i]);
        }
        
        return result;
    }
    
    // bytes1转回字符串
    function bytes1ArrayToString(bytes1[] memory input) public pure returns (string memory) {
        bytes memory result = new bytes(input.length);
        
        for (uint i = 0; i < input.length; i++) {
            result[i] = input[i];
        }
        
        return string(result);
    }
    
    // 演示bytes1的数值操作
    function bytes1Operations() public pure returns (bytes1, bytes1, bytes1) {
        bytes1 a = 0x41; // 'A'
        bytes1 b = 0x42; // 'B'
        
        // 位运算
        bytes1 andResult = a & b; // 按位与
        bytes1 orResult = a | b;  // 按位或
        bytes1 xorResult = a ^ b; // 按位异或
        
        return (andResult, orResult, xorResult);
    }
    
    // 演示字符转换
    function charToBytes1(char c) public pure returns (bytes1) {
        return bytes1(uint8(c));
    }
    
    // 演示数值转换
    function numberToBytes1(uint8 num) public pure returns (bytes1) {
        return bytes1(num);
    }
    
    // 获取bytes1的数值
    function bytes1ToNumber(bytes1 b) public pure returns (uint8) {
        return uint8(b);
    }
    
    // 演示在字符串反转中的应用
    function reverseWithBytes1(string memory input) public pure returns (string memory) {
        bytes memory strBytes = bytes(input);
        uint left = 0;
        uint right = strBytes.length - 1;
        
        while (left < right) {
            // 使用bytes1进行字符交换
            bytes1 temp = bytes1(strBytes[left]);
            strBytes[left] = strBytes[right];
            strBytes[right] = temp;
            
            left++;
            right--;
        }
        
        return string(strBytes);
    }
} 