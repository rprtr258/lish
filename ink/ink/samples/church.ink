app := (x, f) => f(x) # modus ponens
id := x => x

# booleans
# TODO: make it possible to use true/false identifiers
True  := (t, f) => t
False := (t, f) => f
if := (p, t, f) => p(t, f) # also serves as ternary operator
and := (l, r) => if(l, r, False) # same as l(r, False)
andV2 := (l, r) => (t, f) => l(r(t, f), f)
or  := (l, r) => if(l, True, r) # same as l(True, r)
orV2 := (l, r) => (t, f) => l(t, r(t, f))
not   := b => if(b, False, True) # same as b(False, True)
notV2 := b => (t, f) => b(f, t) # "low level" implementation of not

# numbers
zero := (f, x) => x
succ := n => (f, x) => f(n(f, x))
one := succ(zero)
two := succ(one)
# ...
add := (n, m) => n(succ, m)
mul := (n, m) => n(x => add(x, m), zero)
isEven := n => n(not, True)
## test
three := add(one, two)
twelve := mul(two, mul(two, three))
thirteen := succ(twelve)
fiftyTwo := mul(two, mul(two, thirteen))
if(isEven(fiftyTwo), 'true', 'false')

# product types
fst := (x, y) => x # same as True
snd := (x, y) => y # same as False
pair := (x, y) => f => f(x, y)

# sum types
inl := l => (f, g) => f(l)
inr := r => (f, g) => g(r)

# list types
nil := f => fst(f)()
cons := (x, xs) => f => snd(f)(x, xs)
map := (lst, f) => lst(pair(
  () => nil,
  (x, xs) => cons(f(x), map(xs, f))
))
filter := (lst, p) => lst(pair(
  () => nil,
  (x, xs) => (
    tail := filter(xs, p),
    if(p(x), cons(x, tail), tail)
  )
))
foldl := (lst, op, x0) => lst(pair(
  () => x0,
  (x, xs) => foldl(xs, op, op(x0, x))
))
foldr := (lst, op, x0) => lst(pair(
  () => x0,
  (x, xs) => op(foldr(xs, op, x0), x)
))
null := id
head := lst => lst(pair(
  () => inl(null),
  (x, xs) => inr(x)
))
tail := lst => lst(pair(
  () => inl(null),
  (x, xs) => inr(xs)
))
repeat := (n, x) => n(xs => cons(x, xs), nil)
length := lst => foldl(lst, (n, x) => succ(n), zero)
