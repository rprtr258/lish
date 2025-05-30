` uuid `

std := import('https://gist.githubusercontent.com/rprtr258/e208d8a04f3c9a22b79445d4e632fe98/raw/std.ink')

hex := std.hex
range := std.range
map := std.map

uuid := () => (
  ` generate 16 random bytes `
  r := map(urand(16), b => point(b))

  ` helper to map numbers to uniform hexadecimals `
  x := i => (
    s := hex(r.(i))
    len(s) :: {
      1 -> '0' + s
      _ -> s
    }
  )

  ` set version bits per UUID V4 section 4.4 `
  r.6 := (r.6 & 15) | 64
  r.8 := (r.8 & 63) | 128

  x(0) + x(1) + x(2) + x(3) + '-' +
    x(4) + x(5) + '-' +
    x(6) + x(7) + '-' +
    x(8) + x(9) + '-' +
    x(10) + x(11) + x(12) + x(13) + x(14) + x(15)
)
