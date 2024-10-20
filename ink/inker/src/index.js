const fs = require('fs');
const path = require('path');

const express = require('express');
const bodyParser = require('body-parser');

const { evalInk } = require('./eval');

const app = express();

app.get('/', (req, res) => {
    res.sendFile(path.join(__dirname, '../static/index.html'));
});

app.use('/eval', bodyParser.text({
    limit: '25kb',
}));
app.post('/eval', async (req, res) => {
    const inkSource = req.body;
    if (typeof inkSource === 'string' && inkSource.trim() !== '') {
        const result = await evalInk(inkSource);
        res.json(result);
    } else {
        res.json({
            exit: -1,
            error: 'Invalid request',
            output: 'Invalid request',
        });
    }
});

app.use('/static', express.static(path.join(__dirname, '../static')));

app.get('*', (req, res) => res.send('404 not found'));

const PORT = process.env.PORT || 4200;
app.listen(
    PORT,
    () => console.log(`Ink eval running on 0.0.0.0:${PORT}`),
);
