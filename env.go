package main

import "github.com/rprtr258/fun"

type Env struct {
	Outer fun.Option[*Env]
	Data  map[string]Atom
}

func newEnv(outer fun.Option[*Env]) Env {
	return Env{outer, map[string]Atom{}}
}

func newEnvRepl() Env {
	env := newEnv(fun.Invalid[*Env]())
	for name, fun := range namespace() {
		env.sets(name, fun)
	}
	return env
}

func newEnvBind(outer fun.Option[*Env], binds []string, exprs []Atom) Env {
	env := newEnv(outer)
	for i, b := range binds {
		if b == "&" {
			// TODO: List.get(index)
			env.sets(binds[i+1], atomList(exprs[i:]...))
			break
		} else {
			env.sets(b, exprs[i])
		}
	}
	return env
}

func (e Env) Find(key string) (Env, bool) {
	if _, ok := e.Data[key]; ok {
		return e, ok
	}

	if e.Outer.Valid {
		return e.Outer.Value.Find(key)
	}

	return e, false
}

func (e Env) getRoot() Env {
	node := e
	for node.Outer.Valid {
		node = *node.Outer.Value
	}
	return node
}

func (e Env) get(key string) (Atom, bool) {
	env, ok := e.Find(key)
	if !ok {
		return Atom{}, false
	}

	return env.Data[key], true
}

func (e Env) sets(key string, val Atom) {
	e.Data[key] = val
}

func (e Env) set(key string, val Atom) Atom {
	e.sets(key, val)
	return val
}
