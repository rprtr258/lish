{clone} := import('std.ink')
{PI, abs, round} := import('math.ink')
{range, each, flatmap, map} := import('functional.ink')
GRID_SIZE := 50
radius := 20

# Point :: {
#   x : number = 0
#   y : number = 0
#   z : number = 0
# }

repeat := (s, n) => (
  res := []
  (loop := (n) => true :: {
    n > 0 -> (
      res.(len(res)) := clone(s)
      loop(n - 1)
    )
  })(n)
  res
)

Grid := (gridSize) => (
  m_grid := repeat(repeat(' ', gridSize * 2 + 1), gridSize * 2 + 1)
  paint = (x, y, color) => (
    x = x + gridSize
    y = gridSize - y
    m_grid.(x).(y) = color
  )
  clear := () => (
    each(range(~gridSize, gridSize+1, 1), (i, _) =>
      each(range(~gridSize, gridSize+1, 1), (j, _) =>
        paint(i, j, ' ')
      )
    )
    each(range(~gridSize, gridSize+1, 1), (y, _) =>
      paint(0, y, '|')
    )
    each(range(~gridSize, gridSize+1, 1), (x, _) =>
      paint(x, 0, '-')
    )
    paint(0, 0, 'O')
    paint(0, gridSize, '^')
    paint(gridSize, 0, '>')
  )
  drawPolygon := (points) => (
    n := len(points)
    each(range(0, n, 1), (_, i) =>
      paint(points.(i).x, points.(i).y, '*')
      line(points.(i).x, points.(i).y, points.((i + 1) % n).x, points.((i + 1) % n).y, '*')
    )
  )
  line := (x1, y1, x2, y2, color` = '*'`) => (
    dx := (x2 - x1) / 100
    dy := (y2 - y1) / 100
    cx := x1
    cy := y1
    #while (abs(cx - x2) > 1e-2 | abs(cy - y2) > 1e-2) {
    (sub := () => true :: {
      (abs(cx - x2) > 0.01 | abs(cy - y2) > 0.01) -> (
        cx = cx + dx
        cy = cy + dy
        paint(round(cx), round(cy), color)
        sub()
      )
    })
  )
  print := () =>
    each(range(radius, ~radius-1, ~1), (y, _) => (
      each(range(~gridSize, gridSize+1, 1), (x, _) =>
        out(char(m_grid.(x).(y)))
      )
      out('\n')
    ))

  clear()
  {drawPolygon, clear, print}
)

# apply : (p: Point, matrix: [][]number) -> Point
apply := (p, matrix) => {
  x: p.x * matrix.0.0 + p.y * matrix.0.1 + p.z * matrix.0.2
  y: p.x * matrix.1.0 + p.y * matrix.1.1 + p.z * matrix.1.2
  z: p.x * matrix.2.0 + p.y * matrix.2.1 + p.z * matrix.2.2
}

argv := args()
true :: {
  ~(len(argv) == 4) -> (
    out('Usage: ink star.ink <num_of_vertices> <draw_step>
<draw_step> must be in range (0..<num_of_vertices> / 2)
')
    exit(1)
  )
}
vertices := number(argv.2)
step := number(argv.3)
true :: {
  (vertices < 0 | vertices == 0) -> (
    out('Incorrect vertices value\n')
    exit(1)
  )
  (step < 0 | step > vertices / 2) -> (
    out('Incorrect step value\n')
    exit(1)
  )
}
{drawPolygon, print, clear} := Grid(GRID_SIZE)
rotAroundZOrY := false
rotateSpeed := 0.05
poly := flatmap(range(0, vertices, 1), (i, _) =>
  map(range(0, vertices, 1), (t, _) => (
    k := i + step * t
    {
      x: round(radius * cos(2 * PI * k / vertices + PI / 2))
      y: round(radius * sin(2 * PI * k / vertices + PI / 2))
      z: 0
    }
  ))
)
out('\x1B[?25l') # hide cursor
(loop := () => (
  poly = true :: {
    rotAroundZOrY -> map(poly, (p, _) =>
      apply(p, [
        [cos(rotateSpeed), ~sin(rotateSpeed), 0]
        [sin(rotateSpeed), cos(rotateSpeed), 0]
        [0, 0, 1]
      ]) # Z X
    )
    _ -> map(poly, (p, _) => (
        p = apply(p, [
          [cos(rotateSpeed), 0, ~sin(rotateSpeed)]
          [0, 1, 0]
          [sin(rotateSpeed), 0, cos(rotateSpeed)]
        ]) # Y
        p = apply(p, [
          [cos(rotateSpeed), ~sin(rotateSpeed), 0]
          [sin(rotateSpeed), cos(rotateSpeed), 0]
          [0, 0, 1]
        ]) # Z
        #p = apply(p, [
        #  [1, 0, 0]
        #  [0, cos(rotateSpeed), ~sin(rotateSpeed)]
        #  [0, sin(rotateSpeed), cos(rotateSpeed)]
        #]) # X
        p
    ))
  }
  drawPolygon(poly)
  print()
  clear()
  each(range(0, radius*2+1, 1), (i, _) =>
    out('\x1B[1A')
  )
  wait(1, loop)
))()
