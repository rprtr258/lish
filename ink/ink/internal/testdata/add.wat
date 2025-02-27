(module
  (import "ink" "plus" (func $ink__plus (param externref externref) (result externref)))
  (func $add (param $x externref) (param $y externref) (result externref)
    (local.get $x)
    (local.get $y)
    (call $ink__plus)
  )
  (export "add" (func $add))
)
