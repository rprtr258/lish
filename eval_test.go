package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSet(t *testing.T) {
	repl_env := newEnvRepl()
	assert.Equal(t, eval(read("(set a 2)"), repl_env), atomInt(2))
	assert.Equal(t, eval(read("(+ a 3)"), repl_env), atomInt(5))
}

func TestParse_end_of_input(t *testing.T) {
	repl_env := newEnvRepl()
	assert.Equal(t, eval(read("(+ 1 2"), repl_env), atomInt(3))
	assert.Equal(t, eval(read("(+ 1 2 (+ 3 4"), repl_env), atomInt(10))
	assert.Equal(t, eval(read("+ 1 2"), repl_env), atomInt(3))
	assert.Equal(t, eval(read("+ 1 2 (+ 3 4"), repl_env), atomInt(10))
}

func TestEcho(t *testing.T) {
	repl_env := newEnvRepl()
	assert.Equal(t, eval(read("echo 92"), repl_env), atomString("92"))
	// TODO: how to check "abc" is called
	// assert.Equal(t, eval(read("abc"), repl_env), Err(LishErr::from(r#""abc" is not a function"#)));
	// assert.Equal(t, eval(read(r#""abc""#), repl_env), Err(LishErr::from(r#""abc" is not a function"#)));
}

func TestPlus(t *testing.T) {
	repl_env := newEnvRepl()
	assert.Equal(t, eval(read("(+ 1 2 3)"), repl_env), atomInt(6))
	assert.Equal(t, eval(read("(+ 1 2 (+ 1 2))"), repl_env), atomInt(6))
}

// TODO: add tests from history.txt
