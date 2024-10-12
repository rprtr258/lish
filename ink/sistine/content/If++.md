#someday_maybe #project

# If++

Type: Project

```php
if (p) expr
// is Maybe monad
p ? Just(expr) : None

if (p) expr_a else expr_b
// is Either monad
p ? Left(expr_a) : Right(expr_b)

while (p) expr
// is
if (p) expr >>=
if (p) expr >>=
...

return is a named argument
fn f(x):
  return x
// is
fn f(x, return=\x -> call/cc(\c -> c(x))):
    return(x)
```

[Table of Contents](https://craftinginterpreters.com/contents.html)

[https://habr.com/en/company/badoo/blog/428878/](https://habr.com/en/company/badoo/blog/428878/)

[Стековая машина на моноидах](https://habr.com/ru/post/429530/)