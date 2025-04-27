import {join} from 'path';
import express from 'express';
import {text} from 'body-parser';
import {evalInk} from './eval';

const PORT = process.env.PORT || 4200;

const app = express();
app.get('/', (_, res) => res.sendFile(join(__dirname, '../static/index.html')));
app.use('/eval', text({
  limit: '25kb',
}));
app.post('/eval', async (req, res) => {
  const inkSource = req.body;
  if (typeof inkSource !== 'string' || inkSource.trim() === '') {
    res.json({
      exit: -1,
      error: 'Invalid request',
      output: 'Invalid request',
    });
    return;
  }

  const result = await evalInk(inkSource);
  res.json(result);
});
app.use('/static', express.static(join(__dirname, '../static')));
app.get('*', function(_, res) {res.send('404 not found')});
app.listen(PORT, () => console.log(`Ink eval running on 0.0.0.0:${PORT}`));
