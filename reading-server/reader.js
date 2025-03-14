const express = require('express');
const path = require('path');
const cors = require('cors');
const fs = require('fs');

const app = express();
app.use(cors());
app.use(express.json());
app.post('/get', (req, res) => {
    const { hash } = req.body;
    if (!hash) {
        return res.status(400).json({ "Message": "Invalid request, hash missing" });
    }
    const dir = hash[1];
    const filePath = path.join(__dirname, 'uploads', dir, hash);
    // console.log("Checking path:", filePath);
    // console.log("File exists?", fs.existsSync(filePath));
    if (!fs.existsSync(filePath)) {
        return res.status(404).json({ "Message": "Resource not found" });
    }
    const reader = fs.createReadStream(filePath);
    reader.pipe(res);
    reader.on('error', (err) => {
        console.error("Stream error:", err);
        if (!res.headersSent) {
            res.status(500).json({ "Message": "Error reading file" });
        }
    });
});

app.listen(3000, () => {
    console.log("Reader server live...");
});
