package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"github.com/rprtr258/fun"
)

func quasiquote(ast Atom) Atom {
	if ast.Kind == AtomKindList {
		v := ast.Value.(List)
		if len(v) > 0 {
			// TODO: unquote with len(v) > 2 is meaningless
			if v[0] == atomSymbol("unquote") && len(v) >= 2 {
				return v[1]
			}

			res := []Atom{}
			for i := len(v) - 2; i >= 0; i-- {
				x := v[i]
				if x.Kind == AtomKindList && len(x.Value.(List)) > 0 {
					vv := x.Value.(List)
					if vv[0] == atomSymbol("splice-unquote") {
						if len(vv[1:]) == 0 {
							// `(... ,@() res) -> `(... (splice-unquote) res)
							res = []Atom{
								atomSymbol("cons"),
								atomList(atomSymbol("splice-unquote")),
								atomList(res...),
							}
						} else {
							// `(... ,@(s m t h) res) -> `(... s m t h res)
							res = []Atom{
								atomSymbol("concat"),
								vv[1],
								atomList(res...),
							}
						}
						continue
					}
				}

				res = []Atom{
					atomSymbol("cons"),
					quasiquote(x),
					atomList(res...),
				}
			}
			return atomList(res...)
		}
	}
	return Atom{AtomKindList, List{atomSymbol("quote"), ast}}
}

func isMacroCall(ast Atom, env Env) bool {
	if ast.Kind != AtomKindList {
		return false
	}

	v := ast.Value.(List)
	if v[0].Kind != AtomKindSymbol {
		return false
	}

	macroname := v[0].Value.(Symbol)
	a, _ := env.get(macroname)
	return a.Kind == AtomKindLambda && a.Value.(Lambda).isMacro
}

func macroexpand(ast Atom, env Env) Atom {
	for isMacroCall(ast, env) {
		if ast.Kind != AtomKindList {
			panic("unreachable")
		}

		v := ast.Value.(List)
		the_macro := eval(v[0], env)
		if the_macro.Kind == AtomKindError {
			return the_macro
		}
		args := v[1:]
		if the_macro.Kind != AtomKindLambda {
			panic("unreachable")
		}

		vmacro := the_macro.Value.(Lambda)
		lambda_ast := vmacro.ast
		lambda_env := newEnvBind(fun.Valid(&vmacro.env), vmacro.params, args)
		ast = eval(lambda_ast, lambda_env)
		if ast.Kind == AtomKindError {
			return ast
		}
	}
	return ast
}

type FormResult struct {
	a   Atom
	env fun.Option[Env] // valid if tail call optimisation
}

func eval_function_call(fn Atom, unevaluated_args []Atom, env Env) FormResult {
	args := fun.Map[Atom](func(x Atom) Atom {
		return eval(x, env)
	}, unevaluated_args...)

	switch fn.Kind {
	case AtomKindLambda:
		v := fn.Value.(Lambda)
		newEnv := newEnvBind(fun.Valid(&v.env), v.params, args)
		return FormResult{v.ast, fun.Valid(newEnv)}
	case AtomKindFunc:
		return FormResult{a: fn.Value.(Func)(args)}
	case AtomKindString:
		cmd_args := make([]string, 0, len(args))
		for _, arg := range args {
			if arg.Kind != AtomKindString {
				return FormResult{a: lisherr("%s is not string argument", arg)}
			}
			cmd_args = append(cmd_args, string(arg.Value.(String)))
		}

		var stdout, stderr bytes.Buffer

		// TODO: inherit stdin, stdout by default, but pipe if piped
		child := exec.Command(string(fn.Value.(String)), cmd_args...)
		child.Stdin = os.Stdin
		child.Stdout = &stdout
		child.Stderr = &stderr
		status := 0
		if err := child.Run(); err != nil {
			status = err.(*exec.ExitError).ExitCode()
		}
		// TODO: stdout is iter (another kind of list) of lines
		res := map[string]Atom{
			"exit_code": atomInt(status),
			"stdout":    atomString(stdout.String()),
			"stderr":    atomString(stderr.String()),
		}
		return FormResult{a: atomHash(res)}
	case AtomKindHash:
		if len(args) != 1 {
			return FormResult{a: lisherr("Hash is not a function")}
		}

		if args[0].Kind != AtomKindString {
			return FormResult{a: lisherr("Hash key %v must be string", args[0])}
		}

		value, ok := fn.Value.(Hash)[string(args[0].Value.(String))]
		if !ok {
			return FormResult{a: lisherr("Value was not found by key %v", args[0])}
		}

		return FormResult{a: value}
	default:
		return FormResult{a: lisherr("%s is not a function", fn)}
	}
}

func eval(ast Atom, env Env) Atom {
	for {
		ast = macroexpand(ast, env)
		if ast.Kind == AtomKindError {
			return ast
		}
		// evaluate form
		switch ast.Kind {
		case AtomKindList:
			l := ast.Value.(List)
			// nil is evaluated to nil
			if len(l) == 0 {
				return atomNil
			}
			lish_assert_args := func(cmd string, args_count int) Atom {
				return lisherr("%q requires %d argument(s), but got %d in %s", cmd, args_count, len(l[1:]), ast)
			}

			if l[0].Kind == AtomKindSymbol {
				switch s := l[0].Value.(Symbol); s {
				case "quote":
					if len(l[1:]) != 1 {
						return lish_assert_args("quote", 1)
					}
					return l[1]
				case "quasiquoteexpand":
					if len(l[1:]) != 1 {
						return lish_assert_args("quasiquoteexpand", 1)
					}
					return quasiquote(l[1])
				case "quasiquote":
					if len(l[1:]) != 1 {
						return lish_assert_args("quasiquote", 1)
					}

					ast = quasiquote(l[1])
					continue
				case "macroexpand":
					if len(l[1:]) != 1 {
						return lish_assert_args("macroexpand", 1)
					}

					if l[1].Kind != AtomKindList {
						panic("unreachable")
					}

					head := l[1].Value.(List)[0]
					if err := eval(head, env); err.Kind == AtomKindError {
						return err
					}
					return macroexpand(l[1], env)
				case "set":
					if len(l[1:]) != 2 {
						return lish_assert_args("set", 2)
					}

					value := eval(l[2], env)
					if value.Kind == AtomKindError {
						return value
					}
					if l[1].Kind != AtomKindSymbol {
						return lisherr("%s is not a symbol", l[1])
					}
					env.set(l[1].Value.(Symbol), value)
					return value
				case "setmacro":
					if len(l[1:]) != 2 {
						return lish_assert_args("setmacro", 2)
					}

					v := eval(l[2], env)
					if v.Kind == AtomKindError {
						return v
					}
					if v.Kind != AtomKindLambda {
						return lisherr("Macro is not lambda")
					}
					if l[1].Kind != AtomKindSymbol {
						return lisherr("%s is not a symbol", l[1])
					}

					la := v.Value.(Lambda)
					env.set(l[1].Value.(Symbol), atomLambda(Lambda{
						la.eval,
						la.ast,
						la.env,
						la.params,
						true,
						// meta,
					}))
				case "let":
					if l[1].Kind != AtomKindList {
						return lisherr("Let bindings is not a list, but a %s", l[1])
					}

					bindings := l[1].Value.(List)
					if len(bindings)%2 != 0 {
						return lisherr("'let' requires even number of arguments, but got %d in %s", len(l[1:]), ast)
					}

					let_env := newEnv(fun.Valid(&env))
					for i := 0; i < len(bindings); i += 2 {
						var_name := bindings[i]
						var_value := eval(bindings[i+1], let_env)
						if var_value.Kind == AtomKindError {
							return var_value
						}
						if var_name.Kind != AtomKindSymbol {
							return lisherr("%s is not a symbol", var_name)
						}
						let_env.set(var_name.Value.(Symbol), var_value)
					}

					ast = atomList(append([]Atom{atomSymbol("progn")}, l[2:]...)...)
					env = let_env
				case "progn":
					for _, item := range l[1:] {
						if err := eval(item, env); err.Kind == AtomKindError {
							return err
						}
					}
					ast = l[len(l)-1]
				case "if":
					predicate := eval(l[1], env)
					if predicate.Kind == AtomKindError {
						return predicate
					}
					if predicate.Kind == AtomKindBool && !predicate.Value.(Bool) {
						if len(l[1:]) == 3 {
							ast = l[3]
						} else {
							return atomNil
						}
					} else {
						ast = l[2]
					}
				case "eval":
					if len(l[1:]) != 1 {
						return lish_assert_args("eval", 1)
					}
					ast = eval(l[1], env)
					if ast.Kind == AtomKindError {
						return ast
					}
					env = env.root()
					continue
				case "fn":
					if l[1].Kind != AtomKindList {
						return lisherr("fn args must be list of symbols, but it is %s", l[1])
					}

					args := []Symbol{}
					if l[1].Kind == AtomKindList {
						lst := l[1].Value.(List)
						if !fun.All(func(x Atom) bool { return x.Kind == AtomKindSymbol }, lst...) {
							return lisherr("fn args list must consist only of symbols, but not symbol was found in args list: %s", l[1])
						}
						args = fun.Map[Symbol](func(x Atom) Symbol {
							if x.Kind != AtomKindSymbol {
								panic("unreachable")
							}
							return x.Value.(Symbol)
						}, lst...)
					}
					body := l[2]
					return atomLambda(Lambda{
						eval,
						body,
						env,
						args,
						false,
						// meta: Rc::new(Atom::Nil),
					})
				case "pipe":
					if len(l[1:]) != 2 {
						return lish_assert_args("pipe", 2)
					}

					if l[1].Kind != AtomKindList {
						return lisherr("pipe cmds must be list, not %s", l[1])
					}

					cmds := l[1].Value.(List)
					if len(cmds)%2 != 0 {
						return lisherr("pipe cmds count must be even, not %d", len(cmds))
					}

					if l[2].Kind != AtomKindList {
						return lisherr("pipe pipes must be list, not %s", l[2])
					}
					pipes := l[2].Value.(List)
					if len(pipes)%2 != 0 {
						return lisherr("pipe pipes count must be even, not %d", len(pipes))
					}

					// READ PIPES
					type EdgeBeginKind int
					const (
						BProcessStdout EdgeBeginKind = iota
						BProcessStderr
						BFile
						BNull
						BInherit
						BString
					)
					type EdgeBegin struct {
						kind EdgeBeginKind
						s    String
					}

					type EdgeEndKind int
					const (
						EProcessStdin EdgeEndKind = iota
						EFile
						ENull
						EInherit
						// EString(String),
					)
					type EdgeEnd struct {
						kind EdgeEndKind
						s    String
					}
					// type Vertex struct {
					// 	cmd_name string
					// 	stdin    EdgeBegin
					// 	stdout   EdgeEnd
					// 	stderr   EdgeEnd
					// }
					for i := 0; i < len(pipes)/2; i++ {
						var from EdgeBegin
						switch v := pipes[i*2]; v.Kind {
						case AtomKindList:
							v := v.Value.(List)
							v0 := v[0]
							v1 := v[1]
							if v0.Kind != AtomKindString || v1.Kind != AtomKindString {
								return lisherr("unknown pipe beginning: (%s %s)", v0, v1)
							}

							pp := v0.Value.(String)
							s := v1.Value.(String)
							switch pp {
							case "stdout":
								from = EdgeBegin{BProcessStdout, s}
							case "stderr":
								from = EdgeBegin{BProcessStderr, s}
							case "file":
								from = EdgeBegin{BFile, s}
							case "string":
								from = EdgeBegin{BString, s}
							default:
								return lisherr("(%s %s) can't be pipe beginning", pp, s)
							}
						case AtomKindString:
							switch x := v.Value.(String); x {
							case "null":
								from = EdgeBegin{BNull, ""}
							case "inherit":
								from = EdgeBegin{BInherit, ""}
							default:
								return lisherr("unknown pipe beginning: %s", x)
							}
						default:
							return lisherr("unknown pipe beginning: %s", v)
						}

						var into EdgeEnd
						switch v := pipes[i*2+1]; v.Kind {
						case AtomKindList:
							v := v.Value.(List)
							v0 := v[0]
							v1 := v[1]
							if v0.Kind != AtomKindString || v1.Kind != AtomKindString {
								return lisherr("unknown pipe ending: (%s %s)", v0, v1)
							}

							pp := v0.Value.(String)
							s := v1.Value.(String)
							switch pp {
							case "stdin":
								into = EdgeEnd{EProcessStdin, s}
							case "file":
								into = EdgeEnd{EFile, s}
							default:
								return lisherr("(%s %s) can't be pipe ending", pp, s)
							}
						case AtomKindString:
							switch v := v.Value.(String); v {
							case "null":
								into = EdgeEnd{ENull, ""}
							case "inherit":
								into = EdgeEnd{EInherit, ""}
							default:
								return lisherr("unknown pipe ending: %s", v)
							}
						default:
							return lisherr("unknown pipe ending: %s", v)
						}
						fmt.Printf("%#v %#v\n", from, into)
					}

					// SPAWN PROCESSES AND PIPE THEM
					childs := map[string]*exec.Cmd{}
					for i := 0; i < len(cmds)/2; i++ {
						if x := cmds[i*2]; x.Kind != AtomKindString {
							return lisherr("cmd name must be string, not %s", x)
						}
						cmd_name := string(cmds[i*2].Value.(String))
						if _, ok := childs[cmd_name]; ok {
							return lisherr("cmd with name %s is declared at least twice, only one must survive", cmd_name)
						}
						if x := cmds[i*2+1]; x.Kind != AtomKindList {
							return lisherr("cmd args must be list, not %s", x)
						}
						args := cmds[i*2+1].Value.(List)
						if x := args[0]; x.Kind != AtomKindString {
							return lisherr("cmd must be string, not %s", x)
						}
						program := string(args[0].Value.(String))
						program_args := make([]string, 0, len(args[1:]))
						for _, arg := range args[1:] {
							if arg.Kind != AtomKindString {
								return lisherr("cmd arg must be string, not %s", arg)
							}
							program_args = append(program_args, string(arg.Value.(String)))
						}
						var stdin, stdout, stderr bytes.Buffer
						child := exec.Command(program, program_args...)
						child.Stdin = &stdin
						child.Stdout = &stdout
						child.Stderr = &stderr
						childs[cmd_name] = child
					}

					// RETURN RESULTS
					// let status = child.wait().unwrap();
					// let mut res = fnv::FnvHashMap::default();
					// res.insert("exit_code".to_owned(), match status.code() {
					//     Some(exit_code) => Atom::Int(exit_code.into()),
					//     None => Atom::Nil,
					// });
					// res.insert("stdout".to_owned(), {
					//     let mut stdout = String::new();
					//     match child.stdout {
					//         Some(mut s) => {
					//             s.read_to_string(&mut stdout).unwrap();
					//             Atom::String(stdout)
					//         },
					//         None => Atom::Nil,
					//     }
					// });
					// let mut stderr = String::new();
					// match child.stderr {
					//     Some(mut s) => {
					//         s.read_to_string(&mut stderr).unwrap();
					//     },
					//     None => {},
					// }
					// res.insert("stderr".to_owned(), Atom::String(stderr));
					// FormResult::Return(Atom::Hash(std::rc::Rc::new(res)))
					return atomNil
				default:
					// TODO: call shell
					fn := atomString(s)
					if ss, ok := env.get(s); ok {
						fn = ss
					}

					// lisherr(format("Not found '{}'", key)),
					fr := eval_function_call(fn, l[1:], env)
					if fr.env.Valid { // tail call optimisation
						ast = fr.a
						env = fr.env.Value
						continue
					} else {
						return fr.a
					}
				}
			} else {
				// TODO: call shell
				fn := eval(l[0], env)
				if fn.Kind == AtomKindError {
					return fn
				}
				fr := eval_function_call(fn, l[1:], env)
				if fr.env.Valid { // tail call optimisation
					ast = fr.a
					env = fr.env.Value
					continue
				} else {
					return fr.a
				}
			}
		// others are evaluated to themselves
		case AtomKindSymbol:
			res, ok := env.get(ast.Value.(Symbol))
			if !ok {
				return atomString(string(ast.Value.(String)))
			}
			return res
		default:
			return ast
		}
	}
}

func rep(input string, env Env) string {
	return eval(read(input), env).String()
}
