#someday_maybe #project

code generation is more viable than reflection poop

`lens.go`
```go
package lens

import (
	"reflect"
)

type Lens[S, R any] struct {
	Get func(S) R
	Set func(S, R) S
}

func New[S, R any](get func(S) R, set func(S, R) S) Lens[S, R] {
	return Lens[S, R]{
		Get: get,
		Set: set,
	}
}

func Key[K comparable, R any](key K) Lens[map[K]R, R] {
	return New(
		func(s map[K]R) R {
			return s[key]
		},
		func(s map[K]R, r R) map[K]R {
			res := map[K]R{}
			for k, v := range s {
				res[k] = v
			}
			res[key] = r
			return res
		},
	)
}

func Index[R any](idx int) Lens[[]R, R] {
	return New(
		func(s []R) R {
			return s[idx%len(s)]
		},
		func(s []R, r R) []R {
			idx := idx % len(s)
			res := make([]R, len(s))
			copy(res, s)
			res[idx] = r
			return res
		},
	)
}

// func Copy[S any](s S) S {
// 	bytes, _ := json.Marshal(s)
// 	var res S
// 	_ = json.Unmarshal(bytes, &res)
// 	return res
// }

func StructField[S, R any](field string) Lens[S, R] {
	return New(
		func(s S) R {
			return reflect.ValueOf(s).FieldByName(field).Interface().(R)
		},
		func(s S, r R) S {
			source := reflect.ValueOf(s)
			res := reflect.New(reflect.TypeOf(s)).Elem()
			fields := res.NumField()
			for i := 0; i < fields; i++ {
				res.Field(i).Set(source.Field(i))
			}
			res.FieldByName(field).Set(reflect.ValueOf(r))
			return res.Interface().(S)
		},
	)
}

func Compose[S, T, R any](l1 Lens[S, T], l2 Lens[T, R]) Lens[S, R] {
	return New(
		func(s S) R {
			return l2.Get(l1.Get(s))
		},
		func(s S, r R) S {
			return l1.Set(s, l2.Set(l1.Get(s), r))
		},
	)
}

func Over[S, R any](l Lens[S, R], f func(R) R, s S) S {
	return l.Set(s, f(l.Get(s)))
}
```

`lens_test.go`
```go
package lens_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"lens"
)

type User struct {
	ID      uint64
	Name    string
	Friends []string
}

var (
	_userIDLens      = lens.StructField[User, uint64]("ID")
	_userNameLens    = lens.StructField[User, string]("Name")
	_userFriendsLens = lens.StructField[User, []string]("Friends")
)

func TestGetID(t *testing.T) {
	u := User{
		ID:   1,
		Name: "a",
	}
	assert.Equal(t, uint64(1), _userIDLens.Get(u))
}

func TestSetID(t *testing.T) {
	u := User{
		ID:   1,
		Name: "a",
	}
	got := _userNameLens.Set(u, "b")

	assert.Equal(t, User{
		ID:   1,
		Name: "a",
	}, u)
	assert.Equal(t, User{
		ID:   1,
		Name: "b",
	}, got)
}

func TestSetFriend(t *testing.T) {
	u := User{
		ID:      1,
		Name:    "a",
		Friends: []string{"x", "d", "d"},
	}
	l := lens.Compose(_userFriendsLens, lens.Index[string](2))
	got := l.Set(u, "y")

	assert.Equal(t, User{
		ID:      1,
		Name:    "a",
		Friends: []string{"x", "d", "d"},
	}, u)

	assert.Equal(t, User{
		ID:      1,
		Name:    "a",
		Friends: []string{"x", "d", "y"},
	}, got)
}
```