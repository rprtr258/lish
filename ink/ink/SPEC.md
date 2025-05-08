# Ink language specification
This is the source of truth for the Ink programming language.

## Syntax
Ink has an LR(1) grammar that can be parsed successfully with at most 1 lookahead.

Ink's syntax is inspired by JavaScript and Go, but strives to be minimal. This is not necessarily a comprehensive grammar, but expresses the high level structure and mostly up-to-date with the interpreter implementation.

```yaml
Program: Block

SEP: ',' | '\n'
Block: (Expression SEP)*

Expression: (Atom | BinaryExpr | MatchExpr)
ExpressionList: '(' Block ')'


Atom: UnaryExpr | EmptyIdentifier | Identifier | FunctionCall | Literal | ExpressionList

UnaryOp:
  '~' // negation
UnaryExpr: UnaryOp Atom

EmptyIdentifier: '_'
Identifier: (A-Za-z@!?)[A-Za-z0-9@!?]*

FunctionCall: Atom ExpressionList

NumberLiteral: (0-9)+ ['.' (0-9)+] ['e' (0-9)+]
StringLiteral: '\'' ([^\'] | \\ | \')* '\''
BooleanLiteral: 'true' | 'false'
FunctionLiteral: (Identifier | '(' (Identifier SEP)* ')') '=>' (Expression | ExpressionList)
ObjectLiteral: '{' (Expression ':' Expression SEP)* '}'
ListLiteral: '[' (Expression SEP)* ']'
Literal: NumberLiteral | StringLiteral | BooleanLiteral | FunctionLiteral | ObjectLiteral | ListLiteral


BinaryOp:
  '+' | '-' | '*' | '/' | '%' // arithmetic
  | '&' | '|' | '^' // logical and bitwise
  | '>' | '<' // arithmetic comparisons
  | '==' // value comparison operator
  | ':=' // assignment operator
  | '.' // property accessor
BinaryExpr: (Atom | BinaryExpr) BinaryOp (Atom | BinaryExpr)


MatchExpr: (Atom | BinaryExpr) '::' '{' (Expression '->' Expression SEP)* '}'
```

A few quirks of this syntax, and notes about the language:

- All variables use lexical binding and scope, and are bound to the most local ExpressionList (execution block)
- Commas (`Separator` tokens) are always required where they are marked in the formal grammar, but the tokenizer inserts commas on newlines if it can be inserted, except after unary and binary operators and after opening delimiters, so few are required after expressions, before closing delimiters, and before the ':' in an Object literal. Here, they are auto-inserted during tokenization.
    - This allows for "minification" of Ink code the same way JavaScript source can be minified. Minified Ink code can be more compact, because in Ink, almost all whitespace is unnecessary (except those wrapping the `is` operator).
- String literals cannot contain comments. Backticks inside string literals are counted as a part of the string literal. String literals are also multiline.
    - This also allows the programmer to comment out a block with an explanation, simply like this:
    ```
    realCode()
    ` this block is commented out for testing reasons
    someOtherCode()
    `
    moreRealCode()
    ```
- List and object property/element access have the same syntax, which is the reference to the list/object followed by the `.` (property access) operator. This means we access array indexes with `arr.1`, `arr.(index + 1)`, etc. and object property with `obj.propName`, `obj.(computed + propName)`, etc.
- Object (dictionary) keys can be arbitrary expressions, including variable names. If the key is a single identifier, the identifier's name will be used as a key in the dict, and if it's not an identifier (a literal, function call, etc.) the value of the expression will be computed and used as the key. This seems like it may cause trouble conceptually, but turns out to be intuitive in practice.
- Assignment is always (re)declaration of a variable in its local scope; this means, for the moment, there is no way to mutate a variable from a parents scope (it'll just shadow the variable in the local scope). I think this is fine, since it forbids a class of potentially confusing state mutations, but I might change my mind in the future and add an assignment-that-isn't-declare. Note that this doesn't affect composite values -- you can mutate objects from a parents scope.
- Ink allows boolean algebra with both logical/bitwise (`&|^`) and algebraic (`+*~`) operators, and which one is used depends on context.
    - Notably, Ink does not lazy-evaluate logical operators. That means, given `A & B` or `A | B`, both operands are _always evaluated_. This seems simpler and leaves less room for abuse of logical operators in the style of JavaScript's `&&` used as a conditional. I might change my mind on this in the future, but seems like unnecessary complexity at the moment.
- The only control flow constructs are the function call and the match expression (`a :: {b -> c...}`), and the only control flow construct that branches the execution flow is the match expression. This makes Ink programs simple to analyze programmatically and simple to audit manually.
- Ink does not have constants or immutable variables guaranteed by the language. By convention, constants are denoted with identifiers starting with an uppercase letter, like `RootFS`, and mutable variables are denoted otherwise, like `checkCounter`.

## Types
Ink is strongly but dynamically typed, and has seven non-extendable types.

- Number
- String
- Boolean
- Null
- Composite (including both Objects (dictionaries) and Lists, like Lua tables)
- Function

String, Composite, and Function types are reference-typed, which means assigning a composite to a variable just assigns a reference to the same composite or function value. All other types are value-typed, which means assigning these values to variables or calling a function with these values as arguments will create new copies of those values. i.e.

The String type is capable of and designed for holding arbitrary sequential binary data, and is also conventionally used as a byte buffer in file and network I/O operations.

The Null type and value `()` is globally unique and often also used to represent an empty or error value. For example, accessing a nonexistent index of a string or key of a composite value will return not an error, but the null value. Likewise, attempting to read from a nonexistent file will return a null value in the standard library. This borrows and furthers the idea of zero values from Go, and is an experiment I'm taking in Ink towards an exception-free interpreted language.

```
# for simple values
a := 3, b := a
a := 42

b == 42 # false, since assignment of values are all copies

# for composite values
list := [1, 2, 3]
twin := list
clone := clone(list) # makes a shallow clone

list.(len(list)) := 4 # append 4 to list
list.(len(list)) := 5 # append 5 to list

len(list) == 5 # true
len(twin) == 5 # true, since it keeps the same reference
len(clone) == 5 # false, since it keeps a copy of the value instead
```

These are tested in [samples/test.ink](samples/test.ink).

## Concurrency
Ink achieves concurrenty in two ways, through an event loop and through concurrent Ink programs that communicate via serialized message passing.

Callbacks / event loop and closures is one kind of abstraction over concurrency, and message passing to a completely different execution thread is a different kind of abstraction over concurrency. I think this mirrors two different kinds of concurrency in the real world -- concurrency by way of asynchrony (callbacks, event loop) and concurrency by way of isolation and encapsulation in the problem space (threading). So these are both supported by Ink and used in these different contexts.

### Event loop
A single process of Ink program first executes its entrypoint programs, and then optionally exits to an event loop to respond to system events.

### Concurrent processes (WIP)
_NOTE: the Ink interpreter does not currently implement parallel programming APIs (`receive`, `send`, `create`) as outlined in this section._

An Ink program is fundamentally single threaded. In the interpreter, this is enforced by a mutex that acts as an execution lock. This is behind rationale that a program is fundamentally a representation of a single system evolving sequentially, and shared state means two threads are actually a single program, which breeds all sorts of complexity when a single system tries to mutate in two different sequences. Rust's solution is innovative (compile time static checking that shared mutation never occurs), but a more minimal and Inky way dealing with this is to not have shared state, and only communicate by passing serialized data (messages) between threads of execution that are otherwise spawned and execute in isolation. This is in essence JavaScript workers, but where messages can be any serialized data.

Ink implements this with three builtin functions, `receive(processID, handler) => null` and `send(processID, message) => null` for sending and receiving messages, and `create(function) => processID` for spawning threads. ProcessID (pid) is an opaque handle passed around but it's a standard Ink value/type (most likely an integer). Once a function has been spawned off into a separate process, it can choose to listen and receive message. Send will _not block_ even if nothing is listening (nothing in Ink does unless explicitly documented / chosen). The handler will receive the message as its only argument, where the message may be any valid serializable Ink value (i.e. not functions / closures), including `()` (the null value). Because the value is not shared directly between parallel programs for safety and ease of use, message passing in this way incurs at least a copy overhead in performance.

These are the right primitives, but we can build much more sophisticated systems and designs, like a state reducer or a task scheduler, into the standard library as we choose and find useful.


## Builtins

### Metaprogramming and packaging
- `import(string) => any`: import the Ink expressions from another file as a _module_ to a different program file. The values declared in the top frame of the imported module will be entries in the composite value returned by `import`. If currently executing from a file, Ink will search relative to the executing file. Otherwise (e.g. if running from standard input or through the `-eval` flag), Ink will search relative to the current working directory of the running process. Ink programs imported this way are deduplicated by a canonicalized URL within a single Engine.

## Other implementation notes
- Ink source code is fully UTF-8 / Unicode compatible. Unicode printed non-whitespace characters are valid variable and function identifiers, as well as the characters `?`, `!`, and `@`.
- Ink is fully tail call optimized, and tail calls are the default looping / jump primitive for programming in Ink.
- I'm still experimenting with how best to idiomatically do exception handling. My current approach has been to halt on assertion errors (wrong number of arguments, type errors, etc.) that are avoidable with good code that another language may catch "at compile time", and to use error events and null values to signal exceptional conditions during normal execution of programs in other cases, like parse failure and operating system errors.
