# generate bitmap graph images

{log} := import('logging.ink')
{format: f} := import('str.ink')
{range, each, map} := import('functional.ink')
{writeFile: wf} := import('io.ink')
bmp := import('bmp.ink')

# some basic configuration
WIDTH := 600
HEIGHT := 600
SCALE := 80 # pixels per unit
STROKE := 3

# cached variables
halfWidth := WIDTH / 2
halfHeight := HEIGHT / 2
RW := range(0, WIDTH, 1)
RH := range(0, HEIGHT, 1)
white := [250, 250, 250]
grey := [200, 200, 200]

# functions we're going to graph
FUNCTIONS := [
  {
    f: (x) => x * x / 2 - 1.5
    color: [0, 0, 0]
  }
  {
    f: (x) => sin(x)
    color: [255, 0, 100]
  }
  {
    f: (x) => 2 * cos(x)
    color: [120, 0, 250]
  }
  {
    f: (x) => pow(x, 3) / 3 + x * x - 2
    color: [0, 210, 170]
  }
]

# scaling to and from canvas <-> graph dimensions
scaleXToCanvas := (x) => floor(x * SCALE + halfWidth)
scaleYToCanvas := (y) => floor(y * SCALE + halfHeight)
scaleXToGraph := (x) => (x - halfWidth) / SCALE
scaleYToGraph := (y) => (y - halfHeight) / SCALE

# make a big white rectangle
pixels := map(range(0, WIDTH * HEIGHT, 1), (_) => white)
log('finished drawing canvas...')

# axis lines
midX := scaleXToGraph(halfWidth)
midY := scaleYToGraph(halfHeight)
maxX := floor(scaleXToGraph(WIDTH))
maxY := floor(scaleYToGraph(HEIGHT))
drawVertAxis := (x, color) => each(RH, (y, _) => pixels.(y * WIDTH + scaleXToCanvas(x)) := color)
drawHorizAxis := (y, color) => each(RW, (x, _) => pixels.(scaleYToCanvas(y) * WIDTH + x) := color)
each(range(1, maxX + 1, 1), (x, _) => (
  drawVertAxis(x, grey)
  drawVertAxis(~x, grey)
))
each(range(1, maxY + 1, 1), (y, _) => (
  drawHorizAxis(y, grey)
  drawHorizAxis(~y, grey)
))
drawVertAxis(0, [0, 0, 0])
drawHorizAxis(0, [0, 0, 0])
log('finished rendering axes...')

# memoize list that we'll use over and over
strokeRange := range(0, STROKE, 1)

# make a graph for each function at each x
each(RW, (scaledX, _) => (
  x := scaleXToGraph(scaledX)

  each(FUNCTIONS, (item, _) => (
    scaledY := scaleYToCanvas((item.f)(x))

    true :: {
      scaledY > 0 & scaledY < HEIGHT -> each(strokeRange, (xoff, _) => (
        each(strokeRange, (yoff, _) => (
          pixels.((scaledY - yoff) * WIDTH + scaledX - xoff) := item.color
        ))
      ))
    }
  ))
))
log('finished rendering functions...')

# save image
file := bmp(WIDTH, HEIGHT, pixels)
wf('graph.bmp', file, (done) => true :: {
  done -> log('Done!')
  () -> log('Error saving graph :(')
})
