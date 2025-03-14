const express = require("express");
const axios = require("axios");

const { checkContainability } = require("../services/dbService");

const router = express.Router();

const READER_NODES = [
    "http://localhost:3000/get",
    "http://localhost:3001/get",
    "http://localhost:3002/get",
    "http://localhost:3003/get"
];

function getReaderNode(rand) {
    const firstChar = rand.charAt(0).toLowerCase();
    if ("0123".includes(firstChar)) return READER_NODES[0];
    if ("4567".includes(firstChar)) return READER_NODES[1];
    if ("89ab".includes(firstChar)) return READER_NODES[2];
    return READER_NODES[3];
}

async function fetchFile(hash, res) {
    const url = getReaderNode(hash);
    console.log(`Fetching file from: ${url} for hash: ${hash}`);

    try {
        const response = await axios.post(url, { hash }, { responseType: "stream" });
        console.log(" Streaming file...");
        res.setHeader("Content-Disposition", `attachment; filename="${hash}.bin"`);
        res.setHeader("Content-Type", response.headers["content-type"]);
        
        response.data.pipe(res); 
    } catch (error) {
        console.error(" Fetch error:", error.message);
        res.status(500).json({ error: "Error fetching file" });
    }
}

router.get("/", async (req, res) => {
    const { uri, hash } = req.query;

    if (!uri || !hash) {
        return res.status(400).json({ error: "Missing 'uri' or 'hash' in query params" });
    }

    try {
        // console.log(`Checking containability: URI=${uri}, Hash=${hash}`);
        const hashDetails = await checkContainability(uri, hash);
        
        if (!hashDetails) {
            return res.status(403).json({ error: "[Auth_ERR] Hash not found in URI" });
        }

        await fetchFile(hashDetails.Rand, res); 

    } catch (error) {
        console.error(" Error:", error.message);
        res.status(500).json({ error: error.message });
    }
});

module.exports = router;
