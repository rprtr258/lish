` vec3 math library `

f := import('../vendor/std').format
util := import('util')
urR := util.randRange

create := (x, y, z) => {x: x, y: y, z: z}

Zero := create(0, 0, 0)

` square of length of a vec3 is often useful for comparisons,
  and faster when we don't actually need an abs() `
abssq := v => v.x * v.x + v.y * v.y + v.z * v.z
abs := v => pow(abssq(v), 0.5)
neg := v => {x: ~1 * v.x, y: ~1 * v.y, z: ~1 * v.z}

` normalize vector against length of 1 `
norm := v => multiply(v, 1 / abs(v))

eq := (a, b) => (a.x = b.x) & (a.y = b.y) & (a.z = b.z)

add := (a, b) => {
  x: a.x + b.x
  y: a.y + b.y
  z: a.z + b.z
}

sub := (a, b) => {
  x: a.x - b.x
  y: a.y - b.y
  z: a.z - b.z
}

multiply := (v, c) => {
  x: v.x * c
  y: v.y * c
  z: v.z * c
}

divide := (v, c) => {
  x: v.x / c
  y: v.y / c
  z: v.z / c
}

dot := (a, b) => a.x * b.x + a.y * b.y + a.z * b.z

cross := (a, b) => {
  x: a.y * b.z - a.z * b.y
  y: a.z * b.x - a.x * b.z
  z: a.x * b.y - a.y * b.x
}

list := v => [v.x, v.y, v.z]

string := v => f('[{{x}}, {{y}}, {{z}}]', v)

randRange := (min, max) => create(
  urR(min, max)
  urR(min, max)
  urR(min, max)
)

` if after a 1000 attempts, we don't find
  an internal point, return [0, 0, 0] `
randUnitSphere := () => (sub := i => i :: {
  0 -> Zero
  _ -> (
    p := randRange(~1, 1)
    abssq(p) < 1 :: {
      true -> p
      false -> sub(i - 1)
    }
  )
})(1000)

randUnitDisk := () => (sub := i => i :: {
  0 -> Zero
  _ -> (
    p := create(
      urR(~1, 1)
      urR(~1, 1)
      0
    )
    abssq(p) < 1 :: {
      true -> p
      false -> sub(i - 1)
    }
  )
})(1000)

randUnitVec := () => (
  a := urR(0, 2 * util.Pi)
  z := urR(~1, 1)
  r := pow(1 - z * z, 0.5)
  create(r * cos(a), r * sin(a), z)
)

reflect := (v, n) => sub(v, multiply(n, 2 * dot(v, n)))

` nabs is x -> |x|, since abs() exists for vec3 `
nabs := n => n > 0 :: {
  true -> n
  _ -> ~n
}

refract := (uv, n, eta) => (
  cosTheta := dot(neg(uv), n)
  rOutPerp := multiply(add(uv, multiply(n, cosTheta)), eta)
  rOutParallel := multiply(n, ~pow(nabs(1 - abssq(rOutPerp)), 0.5))
  add(rOutPerp, rOutParallel)
)

{
  create: create
  Zero: Zero
  abssq: abssq
  abs: abs
  neg: neg
  norm: norm
  eq: eq
  add: add
  sub: sub
  multiply: multiply
  divide: divide
  dot: dot
  cross: cross
  list: list
  string: string
  randRange: randRange
  randUnitSphere: randUnitSphere
  randUnitDisk: randUnitDisk
  randUnitVec: randUnitVec
  reflect: reflect
  nabs: nabs
  refract: refract
}
