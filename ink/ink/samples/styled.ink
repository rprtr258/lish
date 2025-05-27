{join} := import('str.ink')

FgANSI := (c) => true :: {
  c < 0 -> ''
  c < 8 -> string(c + 30) # 0-7  -> 30-37
  c < 15 -> string(c + 82) # 8-15 -> 90-97
  c < 256 -> '38;5;' + string(c) # 16-255
  _ -> ''
}

BgANSI := (c) => true :: {
  c < 0 -> ''
  c < 8 -> string(c + 40) # 0-7 -> 40-47
  c < 16 -> string(c + 92) # 8-15 -> 100-107
  c < 256 -> '48;5;' + string(c) # 16-255
  _ -> ''
}

mod := {
  reset:     '0'
  bold:      '1'
  faint:     '2'
  italic:    '3'
  underline: '4'
  blink:     '5'
  reverse:   '7'
  crossout:  '9'
  overline:  '53'
  fg: {
    black     : FgANSI(0)
    red       : FgANSI(1)
    green     : FgANSI(2)
    yellow    : FgANSI(3)
    blue      : FgANSI(4)
    magenta   : FgANSI(5)
    cyan      : FgANSI(6)
    white     : FgANSI(7)
    hiBlack   : FgANSI(8)
    hiRed     : FgANSI(9)
    hiGreen   : FgANSI(10)
    hiYellow  : FgANSI(11)
    hiBlue    : FgANSI(12)
    hiMagenta : FgANSI(13)
    hiCyan    : FgANSI(14)
    hiWhite   : FgANSI(15)
    # r,g,b are 0-255
    rgb: (r, g, b) => '38;2;'+string(r)+';'+string(g)+';'+string(b)
  }
  bg: {
    black     : BgANSI(0)
    red       : BgANSI(1)
    green     : BgANSI(2)
    yellow    : BgANSI(3)
    blue      : BgANSI(4)
    magenta   : BgANSI(5)
    cyan      : BgANSI(6)
    white     : BgANSI(7)
    hiBlack   : BgANSI(8)
    hiRed     : BgANSI(9)
    hiGreen   : BgANSI(10)
    hiYellow  : BgANSI(11)
    hiBlue    : BgANSI(12)
    hiMagenta : BgANSI(13)
    hiCyan    : BgANSI(14)
    hiWhite   : BgANSI(15)
    # r,g,b are 0-255
    rgb: (r, g, b) => '48;2;'+string(r)+';'+string(g)+';'+string(b)
  }
}

styled := (s, mods) =>
  '[' + join(mods, ';') + 'm' +
  s +
  '[' + mod.reset + 'm'

out(styled('ERROR', [mod.fg.red, mod.bold])+'\n')

{mod, styled}