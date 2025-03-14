const express = require("express");
// const fetchFileRoute = require("./routes/fetchFile");
const fetchFileRoute = require('./routes/fetchFile')

const app = express();
const PORT = 7000;

app.use(express.json());
app.use("/", fetchFileRoute);

app.listen(PORT, () => {
    console.log(`Proxy Server running at http://localhost:${PORT}`);
});
