#someday_maybe #project

the problem is that components lifecycle can be managed only in `main` function using `defer` which is not composable:
```go
func main() {
    // create component 1
    // defer destroy component 1
    // create component 2
    // defer destroy component 2
    // create component 3
    // defer destroy component 3
    ...
}
```

I want to divide it somehow into

```go
func newComponent1()
func newComponent2()
func newComponent3()
func main() {
    c1 := newComponent1()
    c2 := newComponent2()
    c3 := newComponent3()
    ...
}
```

variants:
- `runtime.Closer` - not reliable
- `fx.Lifecycle` - brings whole DI shit with it, don't understand how it works
- [oklog/run](https://github.com/oklog/run) https://blog.gopheracademy.com/advent-2017/run-group/

raw draft
```go
package fx

import (
	"context"

	"github.com/rprtr258/xerr"
)

// TODO: make working
type Lifecycle struct {
	Name       string
	Start      func(context.Context) error
	StartAsync func(context.Context)
	Close      func()
}

func Combine(name string, lcs ...Lifecycle) Lifecycle {
	return Lifecycle{
		Name: name,
		Start: func(ctx context.Context) error {
			for _, lc := range lcs {
				if err := lc.Run(ctx); err != nil {
					return err
				}
			}
			return nil
		},
		StartAsync: nil,
		Close:      nil,
	}
}

func (lc Lifecycle) close() {
	if lc.Close != nil {
		lc.Close()
	}
}

func (lc Lifecycle) Run(ctx context.Context) error {
	if lc.Start != nil {
		if err := lc.Start(ctx); err != nil {
			return xerr.NewWM(err, "start component", xerr.Fields{
				"component": lc.Name,
			})
		}
		lc.close()
	} else {
		go func() {
			lc.StartAsync(ctx)
			lc.close()
		}()
	}
	return nil
}
```