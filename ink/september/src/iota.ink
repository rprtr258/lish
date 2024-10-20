# generator for consecutive ints, to make clean enums

() => self := {
  i: ~1
  next: () => (
    self.i := self.i + 1
    self.i
  )
}
