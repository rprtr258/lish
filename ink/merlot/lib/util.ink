# potentially shared utility functions and values

str := import('../vendor/str.ink')
trimPrefix := str.trimPrefix

zeroFillTo3Digits := s => len(s) :: {
  0 -> '000'
  1 -> '00' + s
  2 -> '0' + s
  _ -> s
}

# utility for server-rendering large numbers with commas
formatNumber := n => (
  sub := (acc, n) => n :: {
    0 -> acc
    _ -> sub(zeroFillTo3Digits(string(n % 1000)) + ',' + acc, floor(n / 1000))
  }

  threeDigitStr := sub(zeroFillTo3Digits(string(n % 1000)), floor(n / 1000)) :: {
    '000' -> '0'
    _ -> trimPrefix(threeDigitStr, '0')
  }
)

{
  zeroFillTo3Digits: zeroFillTo3Digits
  formatNumber: formatNumber
}