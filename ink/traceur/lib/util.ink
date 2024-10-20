# utilities and shared constants

Pi := 3.1415926535897932385

degreeToRadian := deg => deg * Pi / 180

randRange := (min, max) => min + rand() * (max - min)

clamp := (x, min, max) => true :: {
  x < min -> min
  x > max -> max
  _ -> x
}

doubleDigit := n => n > 9 :: {
  true -> string(n)
  _ -> '0' + string(n)
}

{
  Pi: Pi
  degreeToRadian: degreeToRadian
  randRange: randRange
  clamp: clamp
  doubleDigit: doubleDigit
}