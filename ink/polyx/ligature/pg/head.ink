# shared partial <head>

std := import('https://gist.githubusercontent.com/rprtr258/e208d8a04f3c9a22b79445d4e632fe98/raw/std.ink')
f := std.format

title => f('
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
