- [ ] update docs

## Interpreter
- [ ] Reduce memory allocations at runtime, ideally without impacting runtime Ink performance of hot code paths. If done right, this should free up both memory and CPU.
- [ ] Implement [go-fuzz](http://go-talks.appspot.com/github.com/dvyukov/go-fuzz/slides/go-fuzz.slide#1) to fuzz test the whole toolchain.
- [ ] `NO_COLOR` env or piping for piping output to another application / for scripting use

## Language core
- [x] export last expression from file, instead of all declarations
- [x] `import` from url
- [ ] Type system? I like the way typescript does it, I think Ink’s type checker should be a completely separate layer in the toolchain from the lex/parse/eval layer. But let’s think about the merits of having type annotations and how we can make it simple to lex/parse out while effective at bringing out the forte’s of Ink’s functional style.
  - It seems helpful to think of it as a constraint system on the source code, instead of as something that’s an attribute of the runtime execution itself.
  - Since Ink has no implicit casts, this seems like it'll be straightforward to infer most variable types from their declaration (what they're bound to in the beginning) and recurse up the tree. So to compute "what's the type of this expression?" the type checker will recursively ask its children for their types, and assuming none of them return an error, we can define a `func (n Node) Type() (Type, error)` that recursively descends to type check an entire AST node from the top. We can expose this behind an `ink -check <file>.ink` flag.
  - To support Ink's way of error handling and signaling (returning null data values), the type system must support sum types, i.e. `number | ()`
  - Enforce mutability restrictions at the type level -- variables are (deeply) immutable by default, must be marked as mutable to allow mutation. This also improves functional ergonomics of the language.
  - Potential type annotation: `myVar<type>` (`myVar` is of type `type`), `myFunc<string, boolean => {number}>` (`myFunc` is of type function mapping `string`, `boolean` to type composite of `number`)
  - types are conversion functions, e.g. `string` converts anything to `string` type
- [ ] destructure assignments, e.g. `[a, b] := [1, 2]` or `{map reduce} := import('functional.ink')`
- [ ] value semantics
- [ ] make `:: {...}` mean `true :: {...}` as a common case
- [ ] `<=`, `>=`, `!=` operators
- [ ] async/futures
- [ ] make commas required
- [ ] right-associative evaluation, e.g. `a.f(b)` instead of `(a.f)(b)`. how to `x := a < b :: {...}` instead of `x := (a < b :: {...})`?
- [ ] make it possible to use oneline match expression, block expression
- [ ] escape sequences in strings (e.g. `\x1b`)
- [ ] operator overloading through methods, e.g. `a + b` is equivalent to `a.+(b)`, then we can define `+` as a method on `NumberValue` and `StringValue` that calls `+` on the underlying value and define `+` on custom values and composite values
- [ ] Add to [GitHub Linguist](https://github.com/github/linguist)
