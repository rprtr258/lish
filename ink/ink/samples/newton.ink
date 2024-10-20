# implementation of Newton's method to square root

log := import('logging').log

# higher order function that returns a root finder
# with the given degree of precision threshold
makeNewtonRoot := threshold => n => (
  # tail call optimized root finder
  (find := previous => (
    guess := (previous + n / previous) / 2
    offset := guess * guess - n
    offset < threshold :: {
      true -> guess
      _ -> find(guess)
    }
  ))(n / 2) # initial guess is n / 2
)

# eight degrees of precision chosen arbitrarily, because
# ink prints numbers to 8 decimal digits
root := makeNewtonRoot(0.00000001)

log('root of 2 (~1.4142): ' + string(root(2)))
log('root of 1000 (~31.6): ' + string(root(1000)))
