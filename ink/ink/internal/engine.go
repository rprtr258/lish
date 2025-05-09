package internal

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/rs/zerolog/log"
)

// Context represents a single, isolated execution context with its global heap,
// imports, call stack, and working directory.
type Context struct {
	// WorkingDirectory is absolute path to current working dir (of module system)
	WorkingDirectory string
	// currently executing file's path, if any
	File   string
	Engine *Engine
	// Scope represents the Context's global heap
	Scope *Scope
	// TODO: store position stacke somewhere to use in error reports
}

func (ctx *Context) resetWd() {
	var err error
	ctx.WorkingDirectory, err = os.Getwd()
	if err != nil {
		log.Fatal().Err(err).Stringer("kind", ErrSystem).Msg("could not identify current working directory")
	}
}

// Eval takes a channel of Nodes to evaluate, and executes the Ink programs defined
// in the syntax tree. Eval returns the last value of the last expression in the AST,
// or an error if there was a runtime error.
func (ctx *Context) Eval(node NodeID) Value {
	// ctx.Engine.mu.Lock()
	// defer ctx.Engine.mu.Unlock()

	ast := ctx.Engine.Cmplr.AST

	if _debugvm {
		fmt.Println("AST:")
		fmt.Println(ast.String())
	}

	oldfnid := ctx.Engine.Cmplr.fnid
	fnid := fnID(len(ctx.Engine.Cmplr.funcs))
	ctx.Engine.Cmplr.funcs = append(ctx.Engine.Cmplr.funcs, []Instruction{})
	ctx.Engine.Cmplr.fnid = fnid
	ctx.Engine.Cmplr.compile(node)
	ctx.Engine.Cmplr.fnid = oldfnid

	if _debugvm {
		fmt.Println("BYTECODE:")
		for fnid, fn := range ctx.Engine.Cmplr.funcs {
			fmt.Printf("FN #%d:\n", fnid)
			for i, ins := range fn {
				fmt.Printf("%3d %s\n", i, ins.String())
			}
		}
		fmt.Println()
		// fmt.Println("SCOPE:")
		// logScope(ctx.Engine.Cmplr.scopes)
		// fmt.Println()
	}

	vm := &VM{ctx, []Value{}, []frame{{fnid, 0, &Scope{ctx.Scope, ValueTable{}}}}}
	v := vm.Execute()
	if _debugvm {
		fmt.Println("RESULT:", v.String())
	}
	LogScope(ctx.Scope)
	return v
}

// ExecListener queues an asynchronous callback task to the Engine behind the Context.
// Callbacks registered this way will also run with the Engine's execution lock.
func (ctx *Context) ExecListener(callback func()) {
	ctx.Engine.Listeners.Add(1)
	go func() {
		defer ctx.Engine.Listeners.Done()

		// ctx.Engine.mu.Lock()
		// defer ctx.Engine.mu.Unlock()

		callback()
	}()
}

// ParseReader runs an Ink program defined by an io.Reader.
// This is the main way to invoke Ink programs from Go.
// ParseReader blocks until the Ink program exits.
func ParseReader(ast *AST, filename string, r io.Reader) (NodeID, *Err) {
	b, errr := io.ReadAll(r)
	if errr != nil {
		return -1, &Err{nil, ErrUnknown, errr.Error(), Pos{filename, 0, 0}}
	}

	tokens := tokenize(filename, strings.NewReader("("+string(b)+")"))
	nodes, err := parse(ast, tokens)
	if err.Err != nil {
		return -1, err.Err
	}

	var expr NodeID
	if len(nodes) == 1 {
		expr = nodes[0]
	} else {
		expr = ast.Append(NodeExprList{Pos{filename, 1, 1}, nodes})
	}

	typeProgram, typeCtx := typeInfer(ast, expr, ctxBuiltins)
	fmt.Printf("PROGRAM TYPE: %s\n", typeProgram.String(typeCtx))
	fmt.Printf("PROGRAM TYPE CTX: %s\n", typeCtx.String())

	// optimization passes
	// TODO: optimize pass:
	// listexprSimplifier(ast, expr) // TODO: fix and get back, breaks a.(b) expressions // turn (x) -> x
	// constantFolding(ast, expr) // turn 2+3 -> 5, very naive, e.g. can't simplify 2+x+3 to 5+x // TODO: get back
	//   - dead code elimination

	LogNode(ast.Nodes[expr])
	LogAST(ast)
	return expr, nil
}

// ExecPath is a convenience function to Exec() a program file in a given Context.
func (ctx *Context) ExecPath(path string) (NodeID, *Err) {
	// update Cwd for any potential import() calls this file will make
	ctx.File = path

	var r io.Reader
	if u, err := url.Parse(path); err == nil && u.Scheme != "" {
		ctx.WorkingDirectory = path
		resp, err := http.Get(path)
		if err != nil {
			return -1, &Err{nil, ErrSystem, fmt.Sprintf("could not GET %s for execution: %s", path, err.Error()), Pos{}}
		}
		defer resp.Body.Close()

		r = resp.Body
	} else {
		ctx.WorkingDirectory = filepath.Dir(path)
		file, err := os.Open(path)
		if err != nil {
			return -1, &Err{nil, ErrSystem, fmt.Sprintf("could not open %s for execution: %s", path, err.Error()), Pos{}}
		}
		defer file.Close()

		r = file
	}

	return ParseReader(ctx.Engine.Cmplr.AST, path, r)
}

// Engine is a single global context of Ink program execution.
//
// A single thread of execution may run within an Engine at any given moment,
// and this is ensured by an internal execution lock. An execution's Engine
// also holds all permission and debugging flags.
//
// Within an Engine, there may exist multiple Contexts that each contain different
// execution environments, running concurrently under a single lock.
type Engine struct {
	// Listeners keeps track of the concurrent threads of execution running in the Engine.
	// Call `Engine.Listeners.Wait()` to block until all concurrent execution threads finish on an Engine.
	Listeners sync.WaitGroup

	// Ink de-duplicates imported source files here, where
	// Contexts from imports are deduplicated keyed by the
	// canonicalized import path. This prevents recursive
	// imports from crashing the interpreter and allows other
	// nice functionality.
	Contexts map[string]*Context
	values   map[string]Value
	Cmplr    *compiler

	// Only a single function may write to the stack frames at any moment.
	// mu sync.Mutex
}

func NewEngine() *Engine {
	return &Engine{
		Contexts: map[string]*Context{},
		values:   map[string]Value{},
		// mu:        sync.Mutex{},
		Listeners: sync.WaitGroup{},
		Cmplr:     &compiler{NewAst(), [][]Instruction{}, 0, []map[string]int{}},
	}
}

// CreateContext creates and initializes a new Context tied to a given Engine.
func (eng *Engine) CreateContext() *Context {
	ctx := &Context{
		Engine: eng,
		Scope: &Scope{
			parent: nil,
			vt:     ValueTable{},
		},
	}
	ctx.resetWd()
	ctx.LoadEnvironment()
	return ctx
}
