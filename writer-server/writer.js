// This is static worker version, later on i might implement this as a guthub repo with sidecar for CI/CD

const path = require('path');
const express = require('express');
const fs = require('fs');
const cors = require('cors');
// const { PassThrough } = require('stream');

const app = express();
app.use(cors());

app.post('/stream', (req, res) => {
    let buffer = Buffer.alloc(32);
    let bytesRead = 0;
    req.on('data', (chunk) => {
        if (bytesRead < 32) {
            const remaining = 32 - bytesRead;
            const toCopy = Math.min(remaining, chunk.length);
            chunk.copy(buffer, bytesRead, 0, toCopy);
            bytesRead += toCopy;
            if (bytesRead === 32) {
                // this means we have read first 32 bytes of filename
                // now we store the file into a folder named after its 2nd Byte (Directory offset)
                const fileName = buffer.toString('utf8').replace(/\0/g, '').trim();
                console.log(fileName[1]);
                const dirPath = path.join(__dirname, 'uploads',fileName[1]);
                try{
                    fs.mkdirSync(dirPath, {recursive:true});
                    console.log("Directory Has been made for the file", fileName);
                }
                catch(e){
                    console.log("Already made the directory.")
                }
                const filePath = path.join(__dirname,'uploads',fileName[1], fileName);
                console.log(`Extracted file name: ${fileName}`);
                const writer = fs.createWriteStream(filePath);
                req.pipe(writer);
                writer.on('drain', () => {
                    const written = writer.bytesWritten;
                    const total = parseInt(req.headers['content-length']) - 32;
                    const pWritten = ((written / total) * 100).toFixed(2);
                    console.log(`Processing ... ${pWritten}% done`);
                });
                writer.on('error', () => {
                    res.status(500).json({ "Msg": "Error in file upload" });
                });
                writer.on('close', () => {
                    console.log('Processing ... 100%');
                    res.status(200).json({ status: 'success', msg: `File saved as ${filePath}` });
                });
                writer.write(chunk.slice(toCopy));
            }
        }
    });
});

app.listen(5000, () => {
    console.log("File upload server running on port 5000...");
});
