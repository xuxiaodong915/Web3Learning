package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
    // 连接到以太坊网络
    client,err :=ethclient.Dial("https://eth-sepolia.g.alchemy.com/v2/C1eDetxrjcCL9HMgWActr")
    if err != nil {
        log.Fatal(err)
    }
    // 私钥
    privateKey, err := crypto.HexToECDSA("私钥")
    if err != nil {
        log.Fatal(err)
    }

    // 获取账户地址
    publicKey := privateKey.Public()
    publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
    if !ok {
        log.Fatal("无法获取公钥")
    }
    fromAddress :=crypto.PubkeyToAddress(*publicKeyECDSA)

    // 目标地址
    toAddress := common.HexToAddress("0x7c56915aF6665323e047AAFf22dC744065E60F0f")

    // 获取当前燃气价格
    gasPrice, err := client.SuggestGasPrice(context.Background())
    if err !=nil {
        log.Fatal(err)
    }
    // 获取账号 nonce
    nonce,err := client.PendingNonceAt(context.Background(), fromAddress)

    // 转账金额 (0.01 ETH)
    value := big.NewInt(10000000000000000)

    // 创建交易
    tx := types.NewTransaction(
        nonce,
        toAddress,
        value,
        21000,
        gasPrice,
        nil,
    )

    // 签名交易
    chainID, err :=client.NetworkID(context.Background())
    if err != nil {
        log.Fatal(err)
    }

    signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID),privateKey)
    if err != nil {
        log.Fatal(err)
    }

    // 发送交易
    err = client.SendTransaction(context.Background(), signedTx)

    fmt.Printf("交易已发送: %s\n", signedTx.Hash().Hex())
    fmt.Printf("燃气价格: %s Gwei\n", new(big.Int).Div(gasPrice, big.NewInt(1000000000)).String())
}