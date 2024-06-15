package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/rprtr258/fun"
)

func int_bin_op(init int64, op func(int64, int64) int64) Atom {
	return atomFunc(func(args ...Atom) Atom {
		res := init
		for _, arg := range args {
			res = op(res, arg.Value.(int64))
		}
		return Atom{AtomKindInt, res}
	}, validateArgsOfKind(AtomKindInt))
}

func int_bin_op1(op func(int64, int64) int64) Atom {
	return atomFunc(func(args ...Atom) Atom {
		res := args[0].Value.(int64)
		for _, arg := range args[1:] {
			res = op(res, arg.Value.(int64))
		}
		return Atom{AtomKindInt, res}
	}, validateArgsOfKind(AtomKindInt), validateMinArgs(1))
}

func logical_op(op func(a, b Atom) (bool, bool)) Atom {
	return atomFunc(func(args ...Atom) Atom {
		res := true
		x := args[0]
		for _, arg := range args[1:] {
			y, ok := op(x, arg)
			if !ok {
				return lisherr("incomparable values: %s and %s", x, arg)
			}

			res = res && y
		}
		return atomBool(res)
	}, validateMinArgs(1))
}

func namespace() map[string]Atom {
	return map[string]Atom{
		// ARITHMETIC
		"+": int_bin_op(0, func(a, b int64) int64 { return a + b }),
		"*": int_bin_op(1, func(a, b int64) int64 { return a * b }),
		"/": int_bin_op1(func(a, b int64) int64 { return a / b }),
		"-": atomFunc(func(args ...Atom) Atom {
			if len(args) == 1 {
				return Atom{AtomKindInt, -args[0].Value.(int64)}
			}

			res := args[0].Value.(int64)
			for _, b := range args[1:] {
				res -= b.Value.(int64)
			}
			return Atom{AtomKindInt, res}
		}, validateMinArgs(1), validateArgsOfKind(AtomKindInt)),
		// LOGIC
		"or": atomFunc(func(args ...Atom) Atom {
			res := false
			for _, arg := range args {
				res = res || arg.Value.(bool)
			}
			return atomBool(res)
		}, validateArgsOfKind(AtomKindBool)),
		"and": atomFunc(func(args ...Atom) Atom {
			res := false
			for _, b := range args {
				res = res && b.Value.(bool)
			}
			return atomBool(res)
		}, validateArgsOfKind(AtomKindBool)),
		// COMPARISON
		"=": logical_op(func(a, b Atom) (bool, bool) { return atomEq(a, b), true }),
		"<": logical_op(func(a, b Atom) (bool, bool) {
			res, ok := atomCmp(a, b)
			return res < 0, ok
		}),
		"<=": logical_op(func(a, b Atom) (bool, bool) {
			res, ok := atomCmp(a, b)
			return res <= 0, ok
		}),
		">": logical_op(func(a, b Atom) (bool, bool) {
			res, ok := atomCmp(a, b)
			return res > 0, ok
		}),
		">=": logical_op(func(a, b Atom) (bool, bool) {
			res, ok := atomCmp(a, b)
			return res >= 0, ok
		}),
		// PRINTING
		"dbg": atomFuncNil(func(args ...Atom) {
			fmt.Println(strings.Join(fun.Map[string](func(a Atom) string {
				return fmt.Sprintf("%#v", a)
			}, args...), " "))
		}),
		"print": atomFuncNil(func(args ...Atom) {
			fmt.Print(strings.Join(fun.Map[string](Atom.String, args...), " "))
		}),
		"println": atomFuncNil(func(args ...Atom) {
			fmt.Println(strings.Join(fun.Map[string](Atom.String, args...), " "))
		}),
		"echo": atomFunc(func(args ...Atom) Atom {
			return Atom{AtomKindString, strings.Join(fun.Map[string](Atom.String, args...), " ")}
		}),
		// LIST MANIPULATION
		"cons": atomFunc(func(args ...Atom) Atom {
			elems := args[:len(args)-1]
			switch v := args[len(args)-1]; v.Kind {
			case AtomKindList:
				v := v.Value.([]Atom)
				if len(v) == 0 {
					return atomList(elems...)
				}
				return atomList(append(elems, v...)...)
			default:
				return lisherr("Trying to cons not a list")
			}
		}, validateMinArgs(2)),
		"first": atomFunc(func(args ...Atom) Atom {
			return args[0].Value.([]Atom)[0]
		}, validateExactArgs(1), validateArgsOfKind(AtomKindList)),
		"rest": atomFunc(func(args ...Atom) Atom {
			return atomList(args[0].Value.([]Atom)[1:]...)
		}, validateExactArgs(1), validateArgsOfKind(AtomKindList)),
		"list": atomFunc(atomList),
		"empty?": atomFunc(func(args ...Atom) Atom {
			if args[0].Kind == AtomKindList {
				return atomBool(len(args[0].Value.([]Atom)) == 0)
			}
			return lisherr("Trying to get empty? of %s", strings.Join(fun.Map[string](Atom.String, args...), " "))
		}, validateExactArgs(1)),
		"len": atomFunc(func(args ...Atom) Atom {
			if args[0].Kind == AtomKindList {
				return atomInt(len(args[0].Value.([]Atom)))
			}
			return lisherr("Trying to get len of %s", strings.Join(fun.Map[string](Atom.String, args...), " "))
		}, validateExactArgs(1)),
		"list?": atomFunc(func(args ...Atom) Atom {
			return atomBool(args[0].Kind == AtomKindList)
		}, validateExactArgs(1)),
		"concat": atomFunc(func(args ...Atom) Atom {
			return atomList(fun.ConcatMap(func(arg Atom) []Atom {
				// TODO: support nil
				return arg.Value.([]Atom)
			}, args...)...)
		}, validateArgsOfKind(AtomKindList)),
		// OTHER
		"apply": atomFunc(func(args ...Atom) Atom {
			fn := args[0]
			fnArgs := args[1:]
			switch fn.Kind {
			case AtomKindFunc:
				return fn.Value.(func(...Atom) Atom)(fnArgs...)
			case AtomKindLambda:
				v := fn.Value.(Lambda)
				return eval(v.ast, newEnvBind(fun.Valid(&v.env), v.params, fnArgs))
			default:
				return lisherr("%s is not a function", fn)
			}
		}, validateMinArgs(1)),
		"read": atomFunc(func(args ...Atom) Atom {
			return read(args[0].Value.(string))
		}, validateExactArgs(1), validateArgsOfKind(AtomKindString)),
		"slurp": atomFunc(func(args ...Atom) Atom {
			filename := args[0].Value.(string)
			b, err := os.ReadFile(filename)
			if err != nil {
				return lisherr(err.Error())
			}
			return atomString(string(b))
		}, validateExactArgs(1), validateArgsOfKind(AtomKindString)),
		"join": atomFunc(func(args ...Atom) Atom {
			return atomString(strings.Join(fun.Map[string](func(a Atom) string {
				if a.Kind == AtomKindString {
					return a.Value.(string)
				}
				return a.String()
			}, args...), ""))
		}),
		"throw": atomFunc(func(args ...Atom) Atom {
			return lisherr(args[0].String())
		}, validateExactArgs(1)),
	}
}

// mod core_tests {
//     use crate::{
//         lisherr,
//         args,
//         form,
//         types::Atom,
//     };
//     use super::namespace;

//     fn get_fn(name: &str) -> fn(Vec<Atom>) -> Atom {
//         let ns = namespace();
//         match ns.get(name) {
//             Some(Atom::Func(f)) => f.clone(),
//             _ => unreachable!(),
//         }
//     }

//     macro_rules! test_function {
//         ($test_name:ident, $($fun:expr, $args:expr => $res:expr),* $(,)?) => {
//             #[test]
//             fn $test_name() {
//                 let ns = namespace();
//                 $(
//                     assert_eq!(match ns.get($fun) {
//                         Some(Atom::Func(f)) => f($args),
//                         Some(_) => lisherr!("{:?} is not a function", $fun),
//                         None => lisherr!("{:?} was not found", $fun),
//                     }, Atom::from($res));
//                 )*
//             }
//         }
//     }

//     // (*)
//     test_function!(
//         multiply_nullary,
//         "*", args![] => 1,
//     );

//     // (* 2)
//     test_function!(
//         multiply_unary,
//         "*", args![2] => 2,
//     );

//     // (* 1 2 3)
//     test_function!(
//         multiply_ternary,
//         "*", args![1, 2, 3] => 6,
//     );

//     // (/ 1)
//     test_function!(
//         divide_unary,
//         "/", args![1] => 1,
//     );

//     // (/ 5 2)
//     test_function!(
//         divide_binary,
//         "/", args![5, 2] => 2,
//     );

//     // (/ 22 3 2)
//     test_function!(
//         divide_ternary,
//         "/", args![22, 3, 2] => 3,
//     );

//     // (- 1)
//     test_function!(
//         minus_unary,
//         "-", args![1] => (-1),
//     );

//     // (- 1 2 3)
//     test_function!(
//         minus_ternary,
//         "-", args![1, 2, 3] => (-4),
//     );

//     // (+)
//     test_function!(
//         plus_nullary,
//         "+", args![] => 0,
//     );

//     // (+ 1)
//     test_function!(
//         plus_unary,
//         "+", args![1] => 1,
//     );

//     // (+ 1 2 3)
//     test_function!(
//         plus_ternary,
//         "+", args![1, 2, 3] => 6
//     );

//     test_function!(
//         not_equal_ints,
//         "=", args![1, 2] => false
//     );

//     test_function!(
//         equal_ints,
//         "=", args![2, 2] => true
//     );

//     test_function!(
//         less_true,
//         "<", args![1, 2] => true
//     );

//     test_function!(
//         less_false,
//         "<", args![2, 1] => false
//     );

//     test_function!(
//         less_equal_true,
//         "<=", args![1, 1] => true
//     );

//     test_function!(
//         less_equal_false,
//         "<=", args![2, 1] => false
//     );

//     test_function!(
//         greater_true,
//         ">", args![2, 1] => true
//     );

//     test_function!(
//         greater_false,
//         ">", args![1, 2] => false
//     );

//     test_function!(
//         greater_equal_true,
//         ">=", args![2, 2] => true
//     );

//     test_function!(
//         greater_equal_false,
//         ">=", args![1, 2] => false
//     );

//     /* TODO: rewrite to using write!
//     test_function!(
//         print_int,
//         "print", args![92] => "92"
//     );

//     test_function!(
//         print_ints,
//         "print", args![1, 2, 3] => "1 2 3"
//     );

//     test_function!(
//         print_strs,
//         "print", args!["a", "b", "c"] => "a b c"
//     );

//     test_function!(
//         print_multiline_str,
//         "print", args!["a\nc"] => "a\\nc"
//     );
//     */

//     test_function!(
//         echo_int,
//         "echo", args![1] => "1"
//     );

//     test_function!(
//         echo_strs,
//         "echo", args!["a", "b", "c"] => "a b c"
//     );

//     test_function!(
//         echo_multiline_str,
//         "echo", args!["a\nc"] => "a\nc"
//     );

//     #[test]
//     fn apply_plus() {
//         let ns = namespace();
//         assert_eq!(
//             get_fn("apply")(args![ns.get("+").unwrap().clone(), 1, 2, 3]),
//             Atom::from(6)
//         )
//     }

//     #[test]
//     fn apply_lambda() {
//         use std::rc::Rc;
//         use crate::{eval, env::Env};
//         let lambda = Atom::Lambda {
//             eval,
//             params: vec!["&".to_owned(), "x".to_owned()],
//             ast: Rc::new(Atom::symbol("x")),
//             env: Env::new(None),
//             is_macro: false,
//             // meta: Rc::new(Atom::Nil),
//         };
//         assert_eq!(
//             get_fn("apply")(args![lambda, 1, 2, 3]),
//             form![1, 2, 3]
//         );
//     }

//     #[test]
//     fn apply_int_not_a_function() {
//         assert_eq!(
//             get_fn("apply")(args![1, 2, 3]),
//             lisherr!("1 is not a function")
//         )
//     }

//     #[test]
//     fn cons_int_not_a_list() {
//         assert_eq!(
//             get_fn("cons")(args![1, 2]),
//             lisherr!("Trying to cons not a list")
//         )
//     }

//     test_function!(
//         cons_int,
//         "cons", args![1, form![]] => form![1]
//     );

//     test_function!(
//         concat_int_lists,
//         "concat", args![
//             form![],
//             form![1],
//             form![2, 3],
//             form![4, 5, 6, 7]
//         ] => form![1, 2, 3, 4, 5, 6, 7]
//     );

//     test_function!(
//         concat_int_and_str_lists,
//         "concat", args![
//             form![1, 2, 3],
//             form!["a", "b", "c"]
//         ] => form![1, 2, 3, "a", "b", "c"]
//     );

//     test_function!(
//         list_ints,
//         "list", args![1, 2, 3] => form![1, 2, 3]
//     );

//     test_function!(
//         first_ints,
//         "first", args![form![1, 2, 3]] => 1
//     );

//     #[test]
//     fn first_int_not_a_list() {
//         assert_eq!(
//             get_fn("first")(args![1]),
//             lisherr!("Trying to get first of not list")
//         )
//     }

//     test_function!(
//         rest_ints,
//         "rest", args![form![1, 2, 3]] => args![2, 3]
//     );

//     #[test]
//     fn rest_int_not_a_list() {
//         assert_eq!(
//             get_fn("rest")(args![1]),
//             lisherr!("Trying to get rest of not list")
//         )
//     }

//     test_function!(
//         len_ints,
//         "len", args![form![1, 2, 3]] => 3
//     );

//     #[test]
//     fn len_int_not_a_list() {
//         assert_eq!(
//             get_fn("len")(args![1]),
//             lisherr!("Trying to get len of not list")
//         )
//     }

//     test_function!(
//         islist_nil,
//         "list?", args![form![]] => true
//     );

//     test_function!(
//         islist_int,
//         "list?", args![1] => false
//     );

//     test_function!(
//         islist_ints,
//         "list?", args![form![1, 2, 3]] => true
//     );

//     test_function!(
//         isempty_nil,
//         "empty?", args![form![]] => true
//     );

//     test_function!(
//         isempty_ints,
//         "empty?", args![form![1, 2, 3]] => false
//     );

//     test_function!(
//         read_str,
//         "read", args!["(+ 1 2)"] => form![Atom::symbol("+"), 1, 2]
//     );

//     #[test]
//     fn read_int() {
//         assert_eq!(
//             get_fn("read")(args![1]),
//             lisherr!("1 is not a string")
//         )
//     }

//     test_function!(
//         str_ints,
//         "join", args![1, 2, 3] => "123"
//     );

//     test_function!(
//         str_newline_str,
//         "join", args!["a\nc"] => "a\nc"
//     );

//     // TODO: test slurp
// }
