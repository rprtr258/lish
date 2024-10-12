# Variable declaration syntax
(Sean Barett syntax)
```go
// general syntax, either type or init must be provided
<ident> : [<type>] [= <init>]

<ident> : <type>          // no init, explicit type
<ident> : <type> = <init> // explicity type with init
<ident> := <init>         // init, type inferred
// TODO: how to declare consts? e.g. global funcs and vars and types, jai way is:
<ident> :: <init>
// TODO: public/private declarations?
```

# Macro system
macro are functions that transforms code into code
TODO: guranteed inline
filename, fileline, funcname macros for debugging

# First class types
TODO: not clear when to use $ and when not
see also:
- [odin/parametric polymorphism](https://odin-lang.org/docs/overview/#parametric-polymorphism)
## Basic types
- i8, u8, i16, u16, etc
- Unit
```swift
Type = Union{
	Unit = Unit
	i8 = Unit u8 = Unit i16 = Unit u16 = Unit ...etc
	Func(Args, Ret) = ([]Type, Type)
	Proc(Args)      = []Type // TODO: is it needed?
	Product         = [](string, Type)
	Union           = [](string, Type)
	Pointer         = Type
}
```
## examples
type -> type
- struct/product type
- union
- tagged union/sum type
- hashmaps
- tuples
- `Array(n int, $T Type) Type = [n]T` - sized slice
- `Slice($T Type) Type = []T` - unsized slice, stores `(int, *u8)` len and data pointer

value -> type
- enum
- configuration/flags(?)
- typeof

type -> value
- sizeof, alignment
- fields
- marshaling/unmarshaling
- anything doable with reflection
- `make`, any other allocation
- `Cast($T Type, V $R) = T(V)`
- server http/grpc/etc definition -> swagger/etc
- database definitions -> type safe queries
- max, min value of type

```swift
// []T
Slice :: ($T Type) Type Struct(map[Ident]Type{
	data = uintptr
	len =  isize
	cap =  isize
})
ConstSlice :: ($T Type) Type Struct(map[Ident]Type{
	data = uintptr
	len =  isize
})
Func :: ((...$Args), $Ret) Type func(Args...) Ret
MakeSlice :: ($T Type, len? int, cap? int) Slice(T) make([]T, len||0, cap||len||0)
Map :: (slice Slice($T), f Func(($T), $R)) Slice($R) {
	res := MakeSlice($T, len(slice))
	for i, x := slice {
		*Cast(*$T, slice.data) = f(x)
	}
	return res
}
```

types of types are constraints

# misc
## arguments are always immutable and passed by copy
(семантически по копии, на самом деле т.к. мутаций нет, компилятор может соптимизировать в передачу по ссылке)
