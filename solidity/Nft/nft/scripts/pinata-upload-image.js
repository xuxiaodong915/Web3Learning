#!/usr/bin/env node
"use strict";

require("dotenv").config();
const fs = require("fs");
const path = require("path");
const axios = require("axios");
const FormData = require("form-data");

const PINATA_ENDPOINT = "https://api.pinata.cloud/pinning/pinFileToIPFS";

function ensureEnv() {
    const jwt = process.env.PINATA_JWT;
    if (!jwt) {
        throw new Error("未检测到 PINATA_JWT，请在 .env 中配置你的 Pinata JWT");
    }
    return jwt;
}

function buildGatewayLink(cid) {
    const gateway = process.env.IPFS_GATEWAY || "https://ipfs.io";
    const url = new URL(gateway);
    return `${url.origin}/ipfs/${cid}`;
}

async function pinFile(filePath, name) {
    const stat = fs.statSync(filePath);
    if (!stat.isFile()) {
        throw new Error(`路径不是文件: ${filePath}`);
    }

    const form = new FormData();
    form.append("file", fs.createReadStream(filePath));
    if (name) {
        form.append("pinataMetadata", JSON.stringify({ name }));
    }

    const jwt = ensureEnv();
    const res = await axios.post(PINATA_ENDPOINT, form, {
        maxContentLength: Infinity,
        maxBodyLength: Infinity,
        headers: {
            ...form.getHeaders(),
            Authorization: `Bearer ${jwt}`
        }
    });

    return res.data; // { IpfsHash, PinSize, Timestamp }
}

async function main() {
    const [, , fileArg, nameArg] = process.argv;
    if (!fileArg) {
        console.error("用法: node scripts/pinata-upload-image.js <图片路径> [名称]");
        process.exit(1);
    }

    const absPath = path.resolve(fileArg);
    const displayName = nameArg || path.basename(absPath);

    const result = await pinFile(absPath, displayName);
    const cid = result.IpfsHash;
    console.log(JSON.stringify({
        cid,
        ipfsUri: `ipfs://${cid}`,
        gatewayUri: buildGatewayLink(cid),
        pinata: result
    }, null, 2));
}

main().catch((err) => {
    console.error("上传失败:", err.response?.data || err.message);
    process.exit(1);
}); 