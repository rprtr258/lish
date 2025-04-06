{reduce, objReduce} := import('functional.ink')

# type: any -> 'string' | 'number' | 'boolean' | '()' | 'composite' | 'function'

Types := {
  String: 'string',
  Number: 'number',
  Boolean: 'boolean',
  Null: '()',
  Composite: 'composite',
  Function: 'function',
}

# Type :: any -> bool

# Type funcs: any -> bool
##############
is_any       := (x) => true
is_string    := (x) => type(x) == Types.String
is_number    := (x) => type(x) == Types.Number
is_bool      := (x) => type(x) == Types.Boolean
is_null      := (x) => type(x) == Types.Null
is_composite := (x) => type(x) == Types.Composite
is_function  := (x) => type(x) == Types.Function

# is_value: any -> Type
is_value := t => x => t == x

# optional: Type -> Type
optional := t -> union(is_null, t)

# dict: composite -> Type
dict := t => x =>
  is_composite(x) &
  objReduce(x, (acc, k, v) => acc & t(v), true)

# tuple: []Type -> Type
tuple := types => x =>
  is_composite(x) &
  len(x) == len(types) &
  reduce(types, (acc, t, i) => acc & t(x.(i)), true)

# list: Type -> Type
list := t => x =>
  is_composite(x) &
  reduce(x, (acc, v, _) => acc & t(v), true)

# intersect: (Type, Type) -> Type
intersect := (a, b) => x => a(x) & b(x)

# intersect_all: []Type -> Type
intersect_all := types => x =>
  reduce(types, (acc, t, _) => acc & t(x), true)

# union: (Type, Type) -> Type
union := (a, b) => x => a(x) | b(x)

# union_all: []Type -> Type
union_all := types => x =>
  reduce(types, (acc, t, _) => acc | t(x), false)
##############

# is_same_type: (any, any) -> bool
is_same_type := (a, b) => type(a) == type(b)
