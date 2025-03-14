const axios = require("axios");

const CONFIG_DB_URL = process.env.CONFIG_DB_URL || "http://localhost:6000";

async function getHashDetails(hash) {
    try {
        const response = await axios.get(`${CONFIG_DB_URL}/hash/${hash}`);
        return response.data;
    } catch (error) {
        throw new Error("Failed to fetch hash details");
    }
}

async function getURIDetails(uri) {
    try {
        const response = await axios.get(`${CONFIG_DB_URL}/uri/${uri}`);
        return response.data;
    } catch (error) {
        throw new Error("Failed to fetch URI details");
    }
}

async function checkContainability(uri, hash) {
    const uriDetails = await getURIDetails(uri);
    const hashDetails = await getHashDetails(hash);
    if (!uriDetails.Hashes.includes(hashDetails.Rand)) {
        throw new Error("[Auth_ERR] Hash not contained in URI");
    }
    return hashDetails;
}

module.exports = { checkContainability };
