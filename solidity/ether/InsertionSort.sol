// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

/**
 * @title InsertionSort
 * @dev 插入排序算法示例 - 包含错误版本和正确版本
 */
contract InsertionSort {
    
    // 🚨 错误版本 - 会导致无限循环或下溢
    function insertionSortWrong(uint[] memory a) public pure returns(uint[] memory) {    
        for (uint i = 1; i < a.length; i++) {
            uint temp = a[i];
            uint j = i - 1;
            
            // 问题：j 是 uint 类型，j-- 会导致下溢
            while( (j >= 0) && (temp < a[j])) {
                a[j+1] = a[j];
                j--;  // 🚨 当 j=0 时，j-- 会变成 2^256-1
            }
            a[j+1] = temp;
        }
        return a;
    }
    
    // ✅ 正确版本 1 - 使用 int 类型
    function insertionSortCorrect1(uint[] memory a) public pure returns(uint[] memory) {    
        for (uint i = 1; i < a.length; i++) {
            uint temp = a[i];
            int j = int(i) - 1;  // 使用 int 类型
            
            // 现在可以正确处理负数
            while( (j >= 0) && (temp < a[uint(j)])) {
                a[uint(j+1)] = a[uint(j)];
                j--;
            }
            a[uint(j+1)] = temp;
        }
        return a;
    }
    
    // ✅ 正确版本 2 - 避免负数，使用边界检查
    function insertionSortCorrect2(uint[] memory a) public pure returns(uint[] memory) {    
        for (uint i = 1; i < a.length; i++) {
            uint temp = a[i];
            uint j = i;
            
            // 使用 j > 0 避免下溢
            while( (j > 0) && (temp < a[j-1])) {
                a[j] = a[j-1];
                j--;
            }
            a[j] = temp;
        }
        return a;
    }
    
    // ✅ 正确版本 3 - 最简洁的实现
    function insertionSortCorrect3(uint[] memory a) public pure returns(uint[] memory) {    
        for (uint i = 1; i < a.length; i++) {
            uint key = a[i];
            uint j = i;
            
            // 将比 key 大的元素向后移动
            while (j > 0 && a[j-1] > key) {
                a[j] = a[j-1];
                j--;
            }
            a[j] = key;
        }
        return a;
    }
    
    // 🧪 测试函数
    function testSorting() public pure returns (string memory) {
        uint[] memory testArray = new uint[](5);
        testArray[0] = 64;
        testArray[1] = 34;
        testArray[2] = 25;
        testArray[3] = 12;
        testArray[4] = 22;
        
        // 注意：错误版本会导致交易失败或无限循环
        // uint[] memory result = insertionSortWrong(testArray);
        
        // 使用正确版本
        uint[] memory result = insertionSortCorrect3(testArray);
        
        // 验证结果
        bool isSorted = true;
        for (uint i = 1; i < result.length; i++) {
            if (result[i] < result[i-1]) {
                isSorted = false;
                break;
            }
        }
        
        if (isSorted) {
            return "排序成功！";
        } else {
            return "排序失败！";
        }
    }
    
    // 📊 性能测试
    function benchmarkSort(uint[] memory a) public pure returns (uint) {
        uint startGas = gasleft();
        
        insertionSortCorrect3(a);
        
        uint endGas = gasleft();
        return startGas - endGas;  // 返回消耗的 gas
    }
    
    // 🔍 数组验证函数
    function isSorted(uint[] memory a) public pure returns (bool) {
        for (uint i = 1; i < a.length; i++) {
            if (a[i] < a[i-1]) {
                return false;
            }
        }
        return true;
    }
    
    // 📝 打印数组（用于调试）
    function getArrayString(uint[] memory a) public pure returns (string memory) {
        string memory result = "[";
        for (uint i = 0; i < a.length; i++) {
            if (i > 0) {
                result = string(abi.encodePacked(result, ", "));
            }
            result = string(abi.encodePacked(result, uint2str(a[i])));
        }
        result = string(abi.encodePacked(result, "]"));
        return result;
    }
    
    // 辅助函数：将 uint 转换为 string
    function uint2str(uint _i) internal pure returns (string memory) {
        if (_i == 0) {
            return "0";
        }
        uint j = _i;
        uint len;
        while (j != 0) {
            len++;
            j /= 10;
        }
        bytes memory bstr = new bytes(len);
        uint k = len;
        while (_i != 0) {
            k -= 1;
            uint8 temp = (48 + uint8(_i - _i / 10 * 10));
            bytes1 b1 = bytes1(temp);
            bstr[k] = b1;
            _i /= 10;
        }
        return string(bstr);
    }
} 