# bitmap image test: generate an RGB rainbow gradient

log := import('logging.ink').log
std := import('functional.ink')
range := std.range
map := std.map
reduce := std.reduce
f := import('str.ink').format
bmp := import('bmp').bmp

# modified version of std.append that's faster when we know the length of the base and child arrays
fastappend := (base, child, baseLength, childLength) => (
  (sub := i => i :: {
    childLength -> base
    _ -> (
      base.(baseLength + i) := child.(i)
      sub(i + 1)
    )
  })(0)
)

# utility to time things
startTime := time()
mk := msg => log(f('{{ time }}ms -> {{ msg }}', {
  msg: msg
  time: floor((time() - startTime) * 1000)
}))

# configurations: 720p canvas
W := 720
H := 405
R := (W + H) / 2 * 0.4
mk('start')

radius := (x, y) => (
  xoff := (W / 2) - x
  yoff := (H / 2) - y
  pow((xoff * xoff) + (yoff * yoff), 0.5)
)
pixels := reduce(range(0, H, 1), (acc, y) => (
  row := map(range(0, W, 1), x => (
    radius(x, y) < R :: {
      true -> [200, 255 * (x / W), 255 * (y / H)]
      _ -> [80, 255 * (y / H), 255 * (x / W)]
    }
  ))
  fastappend(acc, row, y * W, W)
), [])
mk('generated pixel array')

file := bmp(W, H, pixels)
mk('generated bmp file')

(import('io').writeFile)('img.bmp', file, result => result :: {
  () -> log(f('file write error: {{ message }}', evt))
})
mk('saved file')
