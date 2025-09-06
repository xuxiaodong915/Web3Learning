#!/usr/bin/env node
"use strict";

require("dotenv").config();
const fs = require("fs");
const path = require("path");
const axios = require("axios");
const FormData = require("form-data");

const PIN_FILE = "https://api.pinata.cloud/pinning/pinFileToIPFS";
const PIN_JSON = "https://api.pinata.cloud/pinning/pinJSONToIPFS";

function ensureEnv() {
    const jwt = process.env.PINATA_JWT;
    if (!jwt) throw new Error("未检测到 PINATA_JWT，请在 .env 中配置");
    return jwt;
}

function buildGatewayLink(cid) {
    const gateway = process.env.IPFS_GATEWAY || "https://ipfs.io";
    const url = new URL(gateway);
    return `${url.origin}/ipfs/${cid}`;
}

async function pinFile(filePath, name) {
    const form = new FormData();
    form.append("file", fs.createReadStream(filePath));
    if (name) form.append("pinataMetadata", JSON.stringify({ name }));
    const jwt = ensureEnv();
    const res = await axios.post(PIN_FILE, form, {
        maxContentLength: Infinity,
        maxBodyLength: Infinity,
        headers: { ...form.getHeaders(), Authorization: `Bearer ${jwt}` }
    });
    return res.data; // { IpfsHash }
}

async function pinJSON(jsonObj, name) {
    const jwt = ensureEnv();
    const body = { pinataContent: jsonObj };
    if (name) body.pinataMetadata = { name };
    const res = await axios.post(PIN_JSON, body, {
        headers: { "Content-Type": "application/json", Authorization: `Bearer ${jwt}` }
    });
    return res.data; // { IpfsHash }
}

function buildMetadata({ name, description, imageCid, external_url, attributes }) {
    const metadata = {
        name,
        description,
        image: `ipfs://${imageCid}`
    };
    if (external_url) metadata.external_url = external_url;
    if (Array.isArray(attributes)) metadata.attributes = attributes;
    return metadata;
}

async function main() {
    const [, , imagePathArg, nameArg, descriptionArg] = process.argv;
    if (!imagePathArg || !nameArg) {
        console.error("用法: node scripts/pinata-end-to-end.js <图片路径> <name> [description]");
        process.exit(1);
    }
    const absImage = path.resolve(imagePathArg);

    // 1) 上传图片
    const imageRes = await pinFile(absImage, path.basename(absImage));
    const imageCid = imageRes.IpfsHash;

    // 2) 生成 metadata 对象
    const metadataObj = buildMetadata({
        name: nameArg,
        description: descriptionArg || "",
        imageCid,
        external_url: undefined,
        attributes: []
    });

    // 3) 上传 metadata
    const metadataRes = await pinJSON(metadataObj, `${nameArg}-metadata.json`);
    const metadataCid = metadataRes.IpfsHash;

    // 4) 输出结果
    console.log(JSON.stringify({
        image: {
            cid: imageCid,
            ipfsUri: `ipfs://${imageCid}`,
            gatewayUri: buildGatewayLink(imageCid)
        },
        metadata: {
            cid: metadataCid,
            ipfsUri: `ipfs://${metadataCid}`,
            gatewayUri: buildGatewayLink(metadataCid)
        },
        tokenURI: `ipfs://${metadataCid}`
    }, null, 2));
}

main().catch((err) => {
    console.error("执行失败:", err.response?.data || err.message);
    process.exit(1);
}); 