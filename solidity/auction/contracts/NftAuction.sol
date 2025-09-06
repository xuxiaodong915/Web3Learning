// SPDX-License-Identifier: MIT
pragma solidity ^0.8;

import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";

contract NftAuction is Initializable {                                                                                              
    struct Auction{
        // 卖家  
        address seller;
        // 持续时间
        uint256 duration;

        // 起始价格
        uint256 startPrice;

        // 起始时间
        uint256 startTime;

        // 是否结束
        bool ended;

        // 最高出价者
        address hightestBidder;

        // 最高价格1
        uint256 hightestBid;

        address nftAddress;

        uint256 tokenId;
    
    }
    // 状态变量
    mapping(uint256 => Auction) public auctions;
    // 下一个拍卖id
    uint256 public nextAuctionId;
    // 管理者地址
    address public admin;
    function initialize() initializer public{
        admin = msg.sender;
    }

    function createAuction(uint256 _duration,uint256 _startPrice,address _nftAddress,uint256 _tokenId) public {
        require(_duration > 1000,"Duration must be greater than 0");
        require(_startPrice > 0,"Start price must be greater than 0");
        auctions[nextAuctionId] = Auction({
            seller:msg.sender,
            duration:_duration,
            startPrice:_startPrice,
            ended:false,
            hightestBidder:address(0),
            hightestBid:0,
            startTime:block.timestamp,
            nftAddress:_nftAddress,
            tokenId:_tokenId
        });
        
        nextAuctionId++;
    }

    function priceBid(uint256 _auctionId) public payable {
        require(_auctionId < nextAuctionId,"Auction not found");
        Auction storage auction = auctions[_auctionId];
        require(!auction.ended,"Auction has ended");
        require(block.timestamp < auction.startTime + auction.duration,"Auction has ended");
        require(msg.value > auction.hightestBid,"Bid must be higher than the current highest bid");
        require(msg.value > auction.startPrice,"Bid must be higher than the starting price");
        if(auction.hightestBidder != address(0)){
            payable(auction.hightestBidder).transfer(auction.hightestBid);
        }
        auction.hightestBidder = msg.sender;
        auction.hightestBid = msg.value;
    }
    
}