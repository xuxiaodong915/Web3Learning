// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract StringReverse {

    function reverseString(string memory input) public pure returns (string memory) {
        bytes memory byteString = bytes(input);
        uint left = 0;
        uint right = byteString.length-1;
        while(left < right){
            bytes1 tempStr = byteString[left];
            byteString[left] = byteString[right];
            byteString[right] = tempStr;
            left++;
            right--;
        }
        return string(byteString);
    }
}