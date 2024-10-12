#someday_maybe #project

go lib for
```go
func Unerror[T any](fn func() (T, error)) T {
  for {
    res, err := fn()
    if err == nil {
      return res
    }

    // TODO: add telegram alerting
    // TODO: add delaying
  }
}
```