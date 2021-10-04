const express = require('express');
const makeRoutes = require("./routes");

const port = process.env.PORT || 3000;

const app = express();
app.use(express.json());
app.use(makeRoutes())

function healthCheck(req, res) {
	return res.json({ healthy: true })
}

app.get('/', healthCheck)

app.listen(port, () => console.log(`Running on port ${ port }`))