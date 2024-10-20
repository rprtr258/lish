` shared partial <head> `

std := import('../../vendor/std')
f := std.format

Template := title => f('
<head>
  <title>{{ title }}</title>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width,initial-scale=1">
  <link rel="stylesheet" href="/static/css/ui.css">
  <link rel="stylesheet" href="/static/css/ligature.css">
  <link href="https://fonts.googleapis.com/css?family=Barlow:400,700&display=swap" rel="stylesheet">
</head>
', {
  title: title
})
