find way to:
given some set S, and group G acting on S, find function from `hash: S -> [0..n)` with following _group invariance property_:
    hash(gs) = hash(s)
for any `g in G`, `s in S`
that is: `hash({gs | g in G}) = hash(Orbit(s)) = {hash(s)}` for any `s in S`
https://www.google.com/search?q=group%20invariant%20hash
