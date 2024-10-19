stack := () => (
  this := {data: []}
  this.push := x => this.data.(len(this.data)) := x
  this.peek := () => this.data.(len(this.data)-1)
  this.pop := () => (
    res := this.data.(len(this.data)-1)
    this.data := (load('std').slice)(this.data, 0, len(this.data)-1)
    res
  )
  this.clear := () => this.data := []
  this.size := () => len(this.data)
  this
)

print := x => out(string(x)+'\n')

s := stack()
(s.push)(0)
(s.push)(1)
(s.push)(2)
print(s.data) `` [0, 1, 2]
(s.pop)()
(s.pop)()
print(s.data) `` [0]