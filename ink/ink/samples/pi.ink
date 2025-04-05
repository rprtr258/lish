# Monte-Carlo estimation of pi using random number generator

{log} := import('logging.ink')
{sqrt} := import('math.ink')
{format: f} := import('str.ink')

# take count from CLI, defaulting to 250k
Count := (c := number(args().2) :: {
  0 -> 250000
  _ -> c
})

# pick a random point in [0, 1) in x and y
randCoord := () => [rand(), rand()]

inCircle := coordPair => (
  # is a given point in a quarter-circle at the origin?
  x := coordPair.0
  y := coordPair.1
  x * x + y * y < 1
)

# initial state
state := {
  inCount: 0
}

# a single iteration of the Monte Carlo simulation
iteration := iterCount => (
  inCircle(randCoord()) :: {
    true -> state.inCount := state.inCount + 1
  }

  # log progress at 100k intervals
  iterCount % 100000 :: {
    0 -> log(f('{{count}}  runs left, Pi at {{pi}}', {count: iterCount, pi: 4 * state.inCount / (Count - iterCount)}))
  }
)

# composable higher order function for looping
loop := f => iter := n => n :: {
  0 -> ()
  _ -> (
    f(n)
    iter(n - 1)
  )
}

loop(iteration)(Count) # do Count times

log(f('Estimate of Pi after {{count}} runs: {{pi}}', {count: Count, pi: 4 * state.inCount / Count}))
