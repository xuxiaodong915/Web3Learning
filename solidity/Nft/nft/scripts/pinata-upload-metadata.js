#!/usr/bin/env node
"use strict";

require("dotenv").config();
const fs = require("fs");
const path = require("path");
const axios = require("axios");

const PINATA_ENDPOINT = "https://api.pinata.cloud/pinning/pinJSONToIPFS";

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

async function pinJSON(jsonObj, name) {
    const jwt = ensureEnv();
    const body = { pinataContent: jsonObj };
    if (name) body.pinataMetadata = { name };

    const res = await axios.post(PINATA_ENDPOINT, body, {
        headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${jwt}`
        }
    });
    return res.data; // { IpfsHash, PinSize, Timestamp }
}

async function main() {
    const [, , jsonPathArg, nameArg] = process.argv;
    if (!jsonPathArg) {
        console.error("用法: node scripts/pinata-upload-metadata.js <metadata.json路径> [名称]");
        process.exit(1);
    }
    const absPath = path.resolve(jsonPathArg);
    const raw = fs.readFileSync(absPath, "utf-8");
    const json = JSON.parse(raw);

    const result = await pinJSON(json, nameArg || path.basename(absPath));
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