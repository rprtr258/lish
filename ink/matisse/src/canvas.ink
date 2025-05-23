` stdlib is implied `
f := format
Pi := Math.PI
Tau := 2 * Pi

` constants `
Width := 1000
Height := 1000

` make the canvas `
Canvas := bind(document, 'getElementById')('canvas')
Ctx := bind(Canvas, 'getContext')('2d')
Canvas.width := Width
Canvas.height := Height

` utilities `
sqrt := x => pow(x, 0.5)
distanceSq := (a, b) => pow(a.0 - b.0, 2) + pow(a.1 - b.1, 2)
distance := (a, b) => sqrt(distanceSq(a, b))
rgb := (r, g, b) => f('rgb({{0}}, {{1}}, {{2}})', [r * 255.99, g * 255.99, b * 255.99])
rgba := (r, g, b, a) => f('rgb({{0}}, {{1}}, {{2}}, {{3}})'
  [r * 255.99, g * 255.99, b * 255.99, a])
Black := rgb(0, 0, 0)
White := 'rgb(255, 255, 255)'
coinflip := () => rand() > 0.5
randCenterBias := (min, max, resolution) => (
  ` random center-biased distribution `
  parts := map(
    range(0, resolution, 1)
    () => rand() / resolution
  )
  min + reduce(parts, (a, b) => a + b, 0) * (max - min)
)
randRange := (min, max) => min + rand() * (max - min)
randInt := (min, max) => floor(randRange(min, max))
randColor := () => rgb(rand(), rand(), rand())
randColorAlpha := () => rgba(rand(), rand(), rand(), rand())
randColorGreyscale := () => (
  r := rand()
  rgb(r, r, r)
)
randColorGreyscaleAlpha := () => (
  r := rand()
  rgba(r, r, r, rand())
)
choose := list => list.floor(randRange(0, len(list)))

` canvas state functions `
setFill := color => Ctx.fillStyle := color
setLineWidth := width => Ctx.lineWidth := width
setStroke := color => Ctx.strokeStyle := color

` canvas draw functions `
fillRect := bind(Ctx, 'fillRect')
strokeRect := bind(Ctx, 'strokeRect')
clearRect := bind(Ctx, 'clearRect')
beginPath := bind(Ctx, 'beginPath')
moveTo := bind(Ctx, 'moveTo')
lineTo := bind(Ctx, 'lineTo')
arc := bind(Ctx, 'arc')
stroke := bind(Ctx, 'stroke')
fill := bind(Ctx, 'fill')
rotate := bind(Ctx, 'rotate')
translate := bind(Ctx, 'translate')

` drawing lines `
drawLine := (start, end) => (
  beginPath()
  moveTo(start.0, start.1)
  lineTo(end.0, end.1)
  stroke()
)
drawSinglePath := points => (
  beginPath()
  start := points.0
  moveTo(start.0, start.1)
  each(
    slice(points, 1, len(points))
    next => lineTo(next.0, next.1)
  )
  stroke()
)
drawPaths := points => reduce(slice(points, 1, len(points)), (last, next) => (
  drawLine(last, next)
  next
), points.0)
strokeArc := (x, y, r, start, end) => (
  beginPath()
  arc(x, y, r, start, end)
  stroke()
)
fillArc := (x, y, r, start, end) => (
  beginPath()
  arc(x, y, r, start, end)
  fill()
)
strokeCircle := (x, y, r) => strokeArc(x, y, r, 0, Tau)
fillCircle := (x, y, r) => fillArc(x, y, r, 0, Tau)
rotateCanvas := angle => (
  translate(Width / 2, Height / 2)
  rotate(angle)
  translate(~Width / 2, ~Height / 2)
)
