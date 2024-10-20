` list function demos `

std := import('std')

log := std.log
stringList := std.stringList

reverse := std.reverse
map := std.map
filter := std.filter
reduce := std.reduce
each := std.each

` create a simple list `
list := [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]


log('Mapped 1-10 list, squared
-> ' + stringList(map(list, n => n * n)))

log('Filtered 1-10 list, evens
-> ' + stringList(filter(list, n => n % 2 = 0)))

log('Reduced 1-10 list, multiplication
-> ' + string(reduce(list, (acc, n) => acc * n, 1)))

log('Reversing 1-10 list
-> ' + stringList(reverse(list)))

log('For-each loop')
each(list, n => log(n))
