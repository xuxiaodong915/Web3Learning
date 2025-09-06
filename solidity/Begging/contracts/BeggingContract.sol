// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

contract BeggingContract {
	address public owner;
	mapping(address => uint256) private donations;

	constructor() {
		owner = msg.sender;
	}

	modifier onlyOwner() {
		require(msg.sender == owner, "Not owner");
		_;
	}

	function donate() external payable {
		require(msg.value > 0, "No ETH sent");
		donations[msg.sender] += msg.value;
	}

	function withdraw() external onlyOwner {
		uint256 amount = address(this).balance;
		require(amount > 0, "No balance");
		payable(owner).transfer(amount);
	}

	function getDonation(address donor) external view returns (uint256) {
		return donations[donor];
	}

	receive() external payable {
		donations[msg.sender] += msg.value;
	}
} 