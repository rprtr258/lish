# list function demos

{log, reverse, map, filter, reduce, each} := import('std.ink')

# create a simple list
list := [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]

log('Mapped 1-10 list, squared
-> ' + string(map(list, (n, _) => n * n)))

log('Filtered 1-10 list, evens
-> ' + string(filter(list, (n, _) => n % 2 = 0)))

log('Reduced 1-10 list, multiplication
-> ' + string(reduce(list, (acc, n, _) => acc * n, 1)))

log('Reversing 1-10 list
-> ' + string(reverse(list)))

log('For-each loop')
each(list, (n, _) => log(n))
