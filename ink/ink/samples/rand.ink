`` a pseudorandom floating point number in interval [0, 1)
`` () => number
`` rand := rand

`` (number, number) => number
rand_range := (min, max) => rand() * (max - min) + min

`` number => number
randn := max => rand_range(0, max)