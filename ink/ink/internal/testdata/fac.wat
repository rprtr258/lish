(module
  (import "ink" "plus" (func $ink__plus (param externref externref) (result externref)))
  (func $fac (param $n f64) (result f64)
    (f64.lt (local.get $n) (f64.const 1))
    if (result f64)
      (f64.const 1)
    else
      (f64.mul (local.get $n) (call $fac (f64.sub (local.get $n) (f64.const 1))))
    end)
  (export "fac" (func $fac))
)
