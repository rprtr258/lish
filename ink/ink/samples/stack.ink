{
  new: () => (
    this := {data: []}
    this.push := x => this.data.(len(this.data)) := x
    this.peek := () => this.data.(len(this.data)-1)
    this.pop := () => (
      res := this.data.(len(this.data)-1)
      this.data := (import('std.ink').slice)(this.data, 0, len(this.data)-1)
      res
    )
    this.clear := () => this.data := []
    this.size := () => len(this.data)
    this
  )
}
