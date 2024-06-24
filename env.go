package main

import "github.com/rprtr258/fun"

type Env struct {
	Outer fun.Option[*Env]
	Data  map[Symbol]Atom
}

func newEnv(outer fun.Option[*Env]) Env {
	return Env{outer, map[Symbol]Atom{}}
}

func newEnvRepl() Env {
	return Env{fun.Invalid[*Env](), namespace}
}

func newEnvBind(outer fun.Option[*Env], binds []Symbol, exprs []Atom) Env {
	env := newEnv(outer)
	for i, b := range binds {
		if b == "&" {
			// TODO: List.get(index)
			env.set(binds[i+1], atomList(exprs[i:]...))
			break
		} else {
			env.set(b, exprs[i])
		}
	}
	return env
}

func (e Env) find(key Symbol) (Env, bool) {
	if _, ok := e.Data[key]; ok {
		return e, ok
	}

	if e.Outer.Valid {
		return e.Outer.Value.find(key)
	}

	return e, false
}

func (e Env) root() Env {
	node := e
	for node.Outer.Valid {
		node = *node.Outer.Value
	}
	return node
}

func (e Env) get(key Symbol) (Atom, bool) {
	env, ok := e.find(key)
	return env.Data[key], ok
}

func (e Env) set(key Symbol, val Atom) {
	e.Data[key] = val
}
