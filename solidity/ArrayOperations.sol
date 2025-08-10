// SPDX-License-Identifier: MIT
pragma solidity ^0.8;

contract ArrayOperations {
    uint[] public Arrayss;
    function AddElements(uint _elements) public {
        Arrayss.push(_elements);
    }

    function deleElement() public {
        Arrayss.pop();
        Arrayss.
    }

    function getLen() public view returns (uint){
        return Arrayss.length;
    }
}