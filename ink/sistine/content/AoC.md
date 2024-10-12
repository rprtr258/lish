[https://adventofcode.com/2021/day/19](https://adventofcode.com/2021/day/19)

```python
from sys import stdin

TRANSFORMS = [
    lambda x, y, z: (x, y, z), lambda x, y, z: (x, z, y), lambda x, y, z: (y, x, z), lambda x, y, z: (y, z, x), lambda x, y, z: (z, x, y), lambda x, y, z: (z, y, x),
    lambda x, y, z: (-x, y, z), lambda x, y, z: (-x, z, y), lambda x, y, z: (-y, x, z), lambda x, y, z: (-y, z, x), lambda x, y, z: (-z, x, y), lambda x, y, z: (-z, y, x),
    lambda x, y, z: (x, -y, z), lambda x, y, z: (x, -z, y), lambda x, y, z: (y, -x, z), lambda x, y, z: (y, -z, x), lambda x, y, z: (z, -x, y), lambda x, y, z: (z, -y, x),
    lambda x, y, z: (x, y, -z), lambda x, y, z: (x, z, -y), lambda x, y, z: (y, x, -z), lambda x, y, z: (y, z, -x), lambda x, y, z: (z, x, -y), lambda x, y, z: (z, y, -x),
    lambda x, y, z: (-x, -y, z), lambda x, y, z: (-x, -z, y), lambda x, y, z: (-y, -x, z), lambda x, y, z: (-y, -z, x), lambda x, y, z: (-z, -x, y), lambda x, y, z: (-z, -y, x),
    lambda x, y, z: (-x, y, -z), lambda x, y, z: (-x, z, -y), lambda x, y, z: (-y, x, -z), lambda x, y, z: (-y, z, -x), lambda x, y, z: (-z, x, -y), lambda x, y, z: (-z, y, -x),
    lambda x, y, z: (x, -y, -z), lambda x, y, z: (x, -z, -y), lambda x, y, z: (y, -x, -z), lambda x, y, z: (y, -z, -x), lambda x, y, z: (z, -x, -y), lambda x, y, z: (z, -y, -x),
    lambda x, y, z: (-x, -y, -z), lambda x, y, z: (-x, -z, -y), lambda x, y, z: (-y, -x, -z), lambda x, y, z: (-y, -z, -x), lambda x, y, z: (-z, -x, -y), lambda x, y, z: (-z, -y, -x),
]

scanners = []
scanner = set()
next(stdin)
for line in stdin:
    if line == '\n':
        continue
    elif "scanner" in line:
        scanners.append(scanner)
        scanner = set()
    else:
        x, y, z = map(int, line.rstrip().split(','))
        scanner.add((x, y, z))
scanners.append(scanner)
N = len(scanners)
for i in range(N):
    for j in range(i + 1, N):
        found = False
        for transform in TRANSFORMS:
            scanner2 = [transform(*v) for v in scanners[j]]
            for p1 in scanners[i]: # TODO: limit to dirichlet principle
                for p2 in scanner2:
                    d = (p2[0] - p1[0], p2[1] - p1[1], p2[2] - p1[2])
                    kkk = sum(
                        (p3[0] + d[0], p3[1] + d[1], p3[2] + d[2]) in scanner2
                        for p3 in scanners[i]
                    )
                    if kkk >= 12:
                        print(i, j, kkk)
                        found = True
                        break
                if found:
                    break
            if found:
                break
        if found:
            break
```

[https://docs.scipy.org/doc/scipy/reference/generated/scipy.signal.correlate.html](https://docs.scipy.org/doc/scipy/reference/generated/scipy.signal.correlate.html)
[Advent of Code](https://adventofcode.com/2021/day/20)
```python
from sys import stdin

algo = input()

def spider(p):
	x, y = p
	return [
		(x + dx, y + dy)
		for dy in [-1, 0, 1]
		for dx in [-1, 0, 1]
	]

def spider2(p):
	x, y = p
	return [
		(x + dx, y + dy)
		for dy in [-2, -1, 0, 1, 2]
		for dx in [-2, -1, 0, 1, 2]
	]

def apply_algo(mp):
	to_check = set()
	for p in mp:
		to_check.update(spider2(p))
	res = set()
	for p in to_check:
		index = sum(
			pow2
			for pow2, p1 in zip([0x100, 0x80, 0x40, 0x20, 0x10, 0x8, 0x4, 0x2, 0x1], spider(p))
			if p1 in mp
		)
		if algo[index] == '#':
			res.add(p)
	return res

input()
mapp = set()
for y, line in enumerate(stdin):
	for x, c in enumerate(line.strip()):
		if c == '#':
			mapp.add((x, y))
print(len(apply_algo(apply_algo(mapp))))
```