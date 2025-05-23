rectGrid := () => (
  setFill('transparent')

  each(range(0, Width, 50), x => (
    each(range(0, Height, 50), y => (
      ` grey rects `
      setLineWidth(randRange(1, 6))
      setStroke(randColorGreyscale())

      radius := randRange(2, 25)
      margin := 25 - radius
      strokeRect(x + margin, y + margin, 2 * radius, 2 * radius)

      ` color rect `
      setLineWidth(randRange(1, 6))
      setStroke(randColorAlpha())

      radius := randRange(2, 25)
      margin := 25 - radius
      strokeRect(x + margin, y + margin, 2 * radius, 2 * radius)
    ))
  ))
)

diagonals := () => (
  setLineWidth(2)
  setStroke(rgba(
    randRange(0.3, 0.9)
    randRange(0.3, 0.9)
    randRange(0.3, 0.9)
    0.3
  ))
  randomPoints := map(range(0, 1000, 1), () => [
    randCenterBias(0, Width, 6)
    randCenterBias(0, Height, 6)
  ])
  drawPaths(randomPoints)
)

rainbowDiagonals := () => (
  setLineWidth(2)
  randomPoints := map(range(0, 1000, 1), () => [
    randCenterBias(0, Width, 5)
    randCenterBias(0, Height, 5)
  ])
  reduce(slice(randomPoints, 1, len(randomPoints)), (last, next) => (
    setStroke(rgba(
      randRange(0.2, 0.8)
      randRange(0.2, 0.8)
      randRange(0.2, 0.8)
      randRange(0.2, 0.6)
    ))
    drawLine(last, next)
    next
  ), randomPoints.0)
)

flowerbed := () => (
  Center := [Width / 2, Height / 2]
  PuffLen := 36
  MaxDist := distance(Center, [40, 40])
  each(range(0, randRange(6, 100), 1), () => (
    target := [randRange(0, Width), randRange(0, Height)]
    dist := distance(Center, target)
    setLineWidth(2)
    setStroke(rgba(0, 0, 0, 1 - dist / MaxDist))
    drawLine(Center, target)

    ` spikes `
    setLineWidth(5)
    setStroke(rgba(0.99, 0.2, 0.2, 1 - dist / MaxDist))
    each(range(0, 5, 1), () => drawLine(target, [
      target.0 + randRange(0, PuffLen) - PuffLen / 2
      target.1 + randRange(0, PuffLen) - PuffLen / 2
    ]))
  ))
)

maze := () => (
  CellSize := 50
  jitter := (coinflip() :: {
    true -> () => 2
    _ -> () => randInt(1, 4)
  })
  setStroke(rgb(0.1, 0.1, 0.1))
  each(range(0, Width, CellSize), x => (
    each(range(0, Height, CellSize), y => (
      coinflip() :: {
        true -> (
          setLineWidth(1)
          drawLine([x, y], [x + CellSize, y + CellSize / jitter()])
        )
      }
      coinflip() :: {
        true -> (
          setLineWidth(5)
          drawLine([x + CellSize, y], [x, y + CellSize])
        )
      }
    ))
  ))
)

radar := () => (
  Center := [Width / 2, Height / 2]
  MaxRadius := 710 ` close to diagonal length `
  each(range(10, MaxRadius, 10), r => (
    setLineWidth(floor((1 - r / MaxRadius) * 12))

    start := randRange(0, Tau)
    strokeArc(Center.0, Center.1, r, start, start + randRange(0, Pi * 0.7))

    start := randRange(0, Tau)
    strokeArc(Center.0, Center.1, r, start, start + randRange(0, Pi * 0.7))
  ))
)

fans := () => (
  Center := [Width / 2, Height / 2]
  MaxRadius := 710 ` close to diagonal length `
  setLineWidth(6)
  each(range(10, MaxRadius, 6), r => (
    setStroke(rgba(rand(), rand(), rand(), r / MaxRadius))
    start := randRange(0, Tau)
    strokeArc(Center.0, Center.1, r, start, start + randRange(0, Pi * 1.5))
  ))
)

grid := () => (
  Center := [Width / 2, Height / 2]
  CellSize := 50
  HalfCell := CellSize / 2
  MaxRadius := 20

  each(range(0, Width, CellSize), x => (
    each(range(0, Height, CellSize), y => rand() > 0.75 :: {
      true -> (
        center := [x + HalfCell, y + HalfCell]

        setFill(Black)
        setLineWidth(2)
        setStroke(Black)

        r := randRange(0, MaxRadius)
        strokeCircle(center.0, center.1, r)
        fillCircle(center.0, center.1, randRange(2, max([r - 2, 2])))

        setLineWidth(1)
        setStroke(randColorAlpha())
        drawLine(Center, center)
      )
    })
  ))
)

dots := () => (
  Center := [Width / 2, Height / 2]
  CellSize := 50
  HalfCell := CellSize / 2
  MaxRadius := 20

  each(range(0, Width, CellSize), x => (
    each(range(0, Height, CellSize), y => (
      center := [x + HalfCell, y + HalfCell]

      r := randRange(0, MaxRadius)
      setLineWidth(1)
      strokeCircle(center.0, center.1, r)
      fillCircle(center.0, center.1, randRange(2, max([r - 2, 2])))
    ))
  ))
)

waves := () => (
  Center := [Width / 2, Height / 2]
  CellSize := 50
  HalfCell := CellSize / 2
  MaxRadius := 20

  each(range(0, Width, CellSize), x => (
    each(range(0, Height, CellSize), y => (
      center := [x + HalfCell, y + HalfCell]

      r := 19 * x / Width + 1
      arcl := (y + CellSize) / Height * Tau
      setLineWidth(4)
      start := randRange(0, Tau)
      strokeArc(center.0, center.1, r, start, start + arcl)
    ))
  ))
)

planets := () => (
  Center := [Width / 2, Height / 2]
  MaxRadius := 710 ` close to diagonal length `
  setLineWidth(1)
  fillCircle(Center.0, Center.1, 40)
  each(range(80, MaxRadius, 40), r => (
    start := randRange(0, Tau)
    end := start + randRange(Pi, Tau)
    strokeArc(Center.0, Center.1, r, start, end)

    ` planet `
    theta := randRange(start, end)
    pos := [
      Center.0 + r * cos(theta)
      Center.1 + r * sin(theta)
    ]
    r := randRange(6, 14)
    setFill(White)
    coinflip() :: {
      true -> (
        fillCircle(pos.0, pos.1, r + 16)
        strokeCircle(pos.0, pos.1, r + 8)
      )
      false -> fillCircle(pos.0, pos.1, r + 8)
    }
    setFill(Black)
    fillCircle(pos.0, pos.1, r)
  ))
)

pathfinder := () => (
  Center := [Width / 2, Height / 2]
  Step := 100
  clampOne := (x, min, max) => clamp(x, x, min, max).start
  reduce(range(0, 700, 1), (last, n) => (
    opacity := 1 - n / 500
    setStroke(rgba(0, 0, 0, opacity))
    next := (n % 2 :: {
      1 -> [
        last.0
        clampOne(last.1 + randRange(~Step, Step), 0, 1000)
      ]
      _ -> [
        clampOne(last.0 + randRange(~Step, Step), 0, 1000)
        last.1
      ]
    })
    drawLine(last, next)
    coinflip() :: {
      true -> (
        setFill(White)
        fillCircle(last.0, last.1, 10)
        setFill(rgba(0.8, 0.1, 0.1, opacity))
        fillCircle(last.0, last.1, 5)
      )
    }
    next
  ), Center)
)

drizzle := () => (
  Radius := 13
  MinRadius := 1

  sunlight? := coinflip()

  each(range(0, 1000, 1), () => (
    x := randRange(Radius, Width - Radius)
    y := randRange(Radius, Height - Radius)
    opacity := pow(y / Height, 1.4)

    coinflip() :: {
      true -> sunlight? :: {
        true -> setFill(rgba(0, 0, 0, opacity))
        _ -> setFill(rgba(
          0.9 * pow(1 - y / Height, 0.5)
          0.8 * pow(1 - y / Height, 0.5)
          0.3 * pow(1 - y / Height, 0.5)
          opacity
        ))
      }
      _ -> setFill(rgba(
        0
        0.8 * x / Width
        y / Height
        opacity
      ))
    }
    fillCircle(x, y, MinRadius + Radius * (1 - y / Height))
  ))
)

pentatone := () => (
  Center := {
    x: Width / 2
    y: Height / 2
  }
  NRects := choose([16, 25, 50, 100, 200])
  RectHeight := Height / 10

  configs := [
    [rgba(0, 0, 0, 1), Height / 13]
    [rgba(0, 0, 0, 0.3), Height / 7]
    [rgba(0, 0, 0, 0.1), Height / 2.5]
  ]

  each(configs, config => (
    color := config.0
    YVariance := config.1

    setFill(color)
    reduce(range(0, NRects, 1), lastX => (
      nextX := lastX + Width / NRects
      yVariance := randRange(~YVariance, YVariance)

      fillRect(
        lastX
        Center.y + yVariance - RectHeight / 2
        Width / NRects
        RectHeight
      )

      nextX
    ), 0)
  ))
)

petals := () => (
  Center := {
    x: Width / 2
    y: Height / 2
  }
  NRects := choose([10, 16, 25])
  AngleVariance := 0.2

  setFill(rgba(rand(), rand(), rand(), 0.08))
  each(range(0, NRects, 1), i => (
    radius := (i / NRects) * (Width / 2)
    drawFn := (coinflip() :: {true -> fillRect, _ -> strokeRect})

    angle := randRange(~AngleVariance, AngleVariance)
    rotateCanvas(angle)
    drawFn(Center.x - radius, Center.y - radius, radius * 2, radius * 2)
    rotateCanvas(~angle)
  ))
)

minesweeper := () => (
  Center := {
    x: Width / 2
    y: Height / 2
  }

  N := choose([30, 50, 75])
  Radius := Height / N / 2
  setLineWidth(2)

  each(range(0, Width, Width / N), x => (
    each(range(0, Height, Height / N), y => (
      center := {
        x: x + Width / N / 2
        y: y + Height / N / 2
      }

      xOffset := center.x - Center.x
      yOffset := center.y - Center.y

      xOffset * xOffset + yOffset * yOffset < Width * Height / 4 :: {
        true -> (
          opacity := (xOffset * xOffset + yOffset * yOffset) / (Width * Height / 4)

          setFill(rgba(0, 0, 0, opacity))
          setStroke(rgba(0, 0, 0, 1 - opacity))

          coinflip() :: {
            true -> coinflip() :: {
              true -> fillRect(center.x - Radius, center.y - Radius, Radius * 2, Radius * 2)
            }
            _ -> strokeCircle(center.x, center.y, Radius)
          }
        )
      }
    ))
  ))
)

options := [
  rectGrid
  diagonals
  rainbowDiagonals
  flowerbed
  maze
  radar
  fans
  grid
  dots
  waves
  planets
  pathfinder
  drizzle
  pentatone
  petals
  minesweeper
]

` pick a random style and generate `
generate := () => (
  clearRect(0, 0, Width, Height)

  setFill(White)
  fillRect(0, 0, Width, Height)

  setFill(Black)
  setStroke(Black)
  setLineWidth(1)
  (options.randInt(0, len(options)))()
)
generate()
