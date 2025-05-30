# percent encoding, also known as URI encoding

std := import('mod.ink').std
log := std.log
reduce := std.reduce
map := std.map
cat := std.cat
hex := std.hex
xeh := std.xeh
str := import('mod.ink').str
digit? := str.digit?
upper? := str.upper?
lower? := str.lower?
upper := str.upper
lower := str.lower

encodeChar := encodeSlash => c => (
  isValidPunct := (encodeSlash :: {
    true -> (c = '.') | (c = '_') | (c = '-') | (c = '~')
    _ -> (c = '.') | (c = '_') | (c = '-') | (c = '~') | (c = '/')
  })
  digit?(c) | upper?(c) | lower?(c) | isValidPunct :: {
    true -> c
    _ -> '%' + upper(hex(point(c)))
  }
)
encodeKeepSlash := piece => cat(map(piece, encodeChar(false)), '')
encode := piece => cat(map(piece, encodeChar(true)), '')

checkRange := (lo, hi) => c => lo < point(c) & point(c) < hi
upperAF? := checkRange(point('A') - 1, point('F') + 1)
lowerAF? := checkRange(point('a') - 1, point('f') + 1)
hex? := c => digit?(c) | upperAF?(c) | lowerAF?(c)
decode := str => (
  s := {
    `
    0 -> default
    1 -> saw %
    2 -> saw 1 hex number
    `
    stage: 0
    buf: ()
  }
  reduce(str, (decoded, curr) => s.stage :: {
    0 -> curr :: {
      '+' -> (
        decoded + ' '
      )
      '%' -> (
        s.stage := 1
        decoded
      )
      _ -> decoded + curr
    }
    1 -> hex?(curr) :: {
      true -> (
        s.stage := 2
        s.buf := curr
        decoded
      )
      _ -> (
        s.stage := 0
        decoded + '%' + curr
      )
    }
    _ -> (
      last := s.buf
      s.stage := 0
      s.buf := ()
      hex?(curr) :: {
        false -> decoded + '%' + last + curr
        _ -> decoded + char(xeh(lower(last + curr)))
      }
    )
  }, '')
)

{
  decode: decode
}