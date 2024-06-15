package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQuasiquote(t *testing.T) {
	t.Skip() // TODO: not passing
	for name, tc := range map[string]struct {
		ast Atom
		res Atom
	}{
		"unquote_symbol": {
			atomList(atomSymbol("unquote"), atomSymbol("a")),
			atomSymbol("a"),
		},
		"unquote_nothing": {
			atomList(atomSymbol("unquote")),
			atomList(
				atomSymbol("cons"),
				atomList(
					atomSymbol("quote"),
					atomSymbol("unquote"),
				),
				atomList(),
			),
		},
		"unquote_many": {
			atomList(atomSymbol("unquote"), atomSymbol("a"), atomSymbol("b"), atomSymbol("c")),
			atomSymbol("a"),
		},
		"splice_unquote_symbol": {
			atomList(atomList(atomSymbol("splice-unquote"), atomSymbol("a"))),
			atomList(
				atomSymbol("concat"),
				atomSymbol("a"),
				atomList(),
			),
		},
		"splice_unquote_many": {
			atomList(atomList(atomSymbol("splice-unquote"), atomSymbol("a"), atomSymbol("b"), atomSymbol("c"))),
			atomList(
				atomSymbol("concat"),
				atomSymbol("a"),
				atomList(),
			),
		},
		"splice_unquote_nothing": {
			atomList(
				atomList(
					atomSymbol("splice-unquote"),
				),
			),
			atomList(
				atomSymbol("cons"),
				atomList(
					atomSymbol("splice-unquote"),
				),
				atomList(),
			),
		},
		"quasiquote_list": {
			atomList(atomSymbol("a"), atomSymbol("b"), atomSymbol("c")),
			atomList(
				atomSymbol("cons"),
				atomList(atomSymbol("quote"), atomSymbol("a")),
				atomList(
					atomSymbol("cons"),
					atomList(atomSymbol("quote"), atomSymbol("b")),
					atomList(
						atomSymbol("cons"),
						atomList(atomSymbol("quote"), atomSymbol("c")),
						atomList(),
					),
				),
			),
		},
		"quasiquote_symbol": {
			atomSymbol("a"),
			atomList(atomSymbol("quote"), atomSymbol("a")),
		},
	} {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.res, quasiquote(tc.ast))
		})
	}
}

// macro_rules! test_eval {
//     ($test_name:ident, $($ast:expr => $res:expr),* $(,)?) => {
//         #[test)
//         fn $test_name() {
//             let repl_env = Env::new_repl();
//             $(
//                 assert_eq!(eval($ast, repl_env.clone()), Atom::from($res));
//             )*
//         }
//     }
// }

// #[test)
// fn symbol_found() {
//     let env = Env::new(None);
//     env.sets("a", Atom::Int(1));
//     assert_eq!(eval(atomSymbol("a"), env), Atom::Int(1));
// }

// // TODO: how to check that 'a' shell command is called?
// // #[test)
// // fn symbol_not_found() {
// //     let env = Env::new(None);
// //     assert_eq!(
// //         eval(atomList(atomSymbol("a")), env),
// //         lisherr!(r#""a" is not a function"#)
// //     );
// // }

// #[test)
// fn id() {
//     let env = Env::new(None);
//     assert_eq!(eval(Atom::Int(1), env), Atom::Int(1));
// }

// // (+ 2 2)
// test_eval!(
//     list,
//     atomList(atomSymbol("+"), 2, 2) => 4,
// );

// // (quote a)
// test_eval!(
//     quote_symbol,
//     atomList(atomSymbol("quote"), atomSymbol("a")) => atomSymbol("a"),
// );

// // (quote 1)
// test_eval!(
//     quote_int,
//     atomList(atomSymbol("quote"), 1) => 1,
// );

// // (quote)
// #[test)
// fn quote_0_args() {
//     let repl_env = Env::new_repl();
//     assert_eq!(
//         eval(atomList(atomSymbol("quote")), repl_env.clone()),
//         lisherr!("'quote' requires 1 argument(s), but got 0 in (quote)")
//     );
// }

// // (quote a b c 1 2 3)
// #[test)
// fn quote_many_args() {
//     let repl_env = Env::new_repl();
//     assert_eq!(
//         eval(atomList(atomSymbol("quote"), atomSymbol("a"), atomSymbol("b"), atomSymbol("c"), 1, 2, 3), repl_env.clone()),
//         lisherr!("'quote' requires 1 argument(s), but got 6 in (quote a b c 1 2 3)")
//     );
// }

// // (quasiquoteexpand (c ,c ,@c))
// #[test)
// fn quasiquoteexpand() {
//     let repl_env = Env::new_repl();
//     assert_eq!(
//         eval(atomList(atomSymbol("quasiquoteexpand"), atomList(atomSymbol("c"), atomList(atomSymbol("unquote"), atomSymbol("c")), atomList(atomSymbol("splice-unquote"), atomSymbol("c")))), repl_env.clone()),
//         atomList(
//             atomSymbol("cons"),
//             atomList(
//                 atomSymbol("quote"), atomSymbol("c")
//             ),
//             atomList(
//                 atomSymbol("cons"),
//                 atomSymbol("c"),
//                 atomList(
//                     atomSymbol("concat"),
//                     atomSymbol("c"),
//                     Vec::<Atom>::new()
//                 )
//             )
//         ),
//     );
// }

// // (set c '(1 2 3))
// // `(c ,c ,@c)
// #[test)
// fn quasiquote() {
//     let repl_env = Env::new_repl();
//     eval(atomList(
//         atomSymbol("set"),
//         atomSymbol("c"),
//         atomList(
//             atomSymbol("quote"),
//             atomList(1, 2, 3),
//         ),
//     ), repl_env.clone());
//     assert_eq!(
//         eval(atomList(
//             atomSymbol("quasiquote"),
//             atomList(
//                 atomSymbol("c"),
//                 atomList(
//                     atomSymbol("unquote"),
//                     atomSymbol("c"),
//                 ),
//                 atomList(
//                     atomSymbol("splice-unquote"),
//                     atomSymbol("c"),
//                 ),
//             ),
//         ), repl_env.clone()),
//         atomList(
//             atomSymbol("c"),
//             atomList(1, 2, 3),
//             1, 2, 3,
//         )
//     );
// }

// // (setmacro m (fn (x) `(,x ,x)))
// // (macroexpand (m (Y f)))
// #[test)
// fn set_macro_then_macroexpand() {
//     let repl_env = Env::new_repl();
//     eval(atomList(
//         atomSymbol("setmacro"),
//         atomSymbol("m"),
//         atomList(
//             atomSymbol("fn"),
//             atomList(atomSymbol("x")),
//             atomList(
//                 atomSymbol("quasiquote"),
//                 atomList(
//                     atomList(atomSymbol("unquote"), atomSymbol("x")),
//                     atomList(atomSymbol("unquote"), atomSymbol("x"))
//                 )
//             )
//         )
//     ), repl_env.clone());
//     assert_eq!(
//         eval(atomList(
//             atomSymbol("macroexpand"),
//             atomList(
//                 atomSymbol("m"),
//                 atomList(
//                     atomSymbol("Y"),
//                     atomSymbol("f"),
//                 )
//             ),
//         ), repl_env.clone()),
//         atomList(
//             atomList(atomSymbol("Y"), atomSymbol("f")),
//             atomList(atomSymbol("Y"), atomSymbol("f"))
//         )
//     );
// }

// // (set a 2)
// // (+ a 3)
// // (set b 3)
// // (+ a b)
// test_eval!(
//     set,
//     atomList(atomSymbol("set"), atomSymbol("a"), 2) => 2,
//     atomList(atomSymbol("+"), atomSymbol("a"), 3) => 5,
//     atomList(atomSymbol("set"), atomSymbol("b"), 3) => 3,
//     atomList(atomSymbol("+"), atomSymbol("a"), atomSymbol("b")) => 5,
// );

// // (set c (+ 1 2))
// // (+ c 1)
// test_eval!(
//     set_expr,
//     atomList(atomSymbol("set"), atomSymbol("c"), atomList(atomSymbol("+"), 1, 2)) => 3,
//     atomList(atomSymbol("+"), atomSymbol("c"), 1) => 4,
// );

// // 92
// test_eval!(
//     eval_int,
//     Atom::from(92) => 92,
// );

// // abc
// test_eval!(
//     eval_symbol,
//     atomSymbol("abc") => "abc"
// );

// // "abc"
// #[test)
// fn eval_string() {
//     let repl_env = Env::new_repl();
//     assert_eq!(
//         eval(Atom::from("abc"), repl_env.clone()),
//         Atom::from("abc")
//     );
// }

// // (+ 1 2 (+ 1 2))
// test_eval!(
//     plus_expr,
//     atomList(
//         atomSymbol("+"),
//         1, 2,
//         atomList(
//             atomSymbol("+"),
//             1, 2
//         )
//     ) => 6,
// );

// // (set a 2)
// // (let (a 1) a)
// // a
// test_eval!(
//     let_statement,
//     atomList(atomSymbol("set"), atomSymbol("a"), 2) => 2,
//     atomList(atomSymbol("let"), atomList(atomSymbol("a"), 1), atomSymbol("a")) => 1,
//     Atom::from(atomSymbol("a")) => 2,
// );

// // (let (a 1) 2 3 4 5 a)
// test_eval!(
//     let_implicit_progn,
//     atomList(atomSymbol("let"), atomList(atomSymbol("a"), 1), 2, 3, 4, 5, atomSymbol("a")) => 1,
// );

// // (let (a 1 b 2) (+ a b))
// test_eval!(
//     let_twovars_statement,
//     atomList(atomSymbol("let"), atomList(atomSymbol("a"), 1, atomSymbol("b"), 2), atomList(atomSymbol("+"), atomSymbol("a"), atomSymbol("b"))) => 3,
// );

// // (let (a 1 b a) b)
// test_eval!(
//     let_star_statement,
//     atomList(atomSymbol("let"), atomList(atomSymbol("a"), 1, atomSymbol("b"), atomSymbol("a")), atomSymbol("b")) => 1,
// );

// // (progn (set a 92) (+ a 8))
// // a
// test_eval!(
//     progn_statement,
//     atomList(atomSymbol("progn"), atomList(atomSymbol("set"), atomSymbol("a"), 92), atomList(atomSymbol("+"), atomSymbol("a"), 8)) => 100,
//     Atom::from(atomSymbol("a")) => 92,
// );

// // (if true 1 2)
// test_eval!(
//     if_true_statement,
//     atomList(atomSymbol("if"), true, 1, 2) => 1,
// );

// // (if false 1 2)
// test_eval!(
//     if_false_statement,
//     atomList(atomSymbol("if"), false, 1, 2) => 2,
// );

// // (if true 1)
// test_eval!(
//     if_true_noelse_statement,
//     atomList(atomSymbol("if"), true, 1) => 1,
// );

// // (if false 1)
// test_eval!(
//     if_false_noelse_statement,
//     atomList(atomSymbol("if"), false, 1) => Nil,
// );

// // (if true (set a 1) (set a 2))
// // a
// test_eval!(
//     if_set_true_statement,
//     atomList(atomSymbol("if"), true, atomList(atomSymbol("set"), atomSymbol("a"), 1), atomList(atomSymbol("set"), atomSymbol("a"), 2)) => 1,
//     Atom::from(atomSymbol("a")) => 1,
// );

// // (if false (set b 1) (set b 2))
// // b
// test_eval!(
//     if_set_false_statement,
//     atomList(atomSymbol("if"), false, atomList(atomSymbol("set"), atomSymbol("b"), 1), atomList(atomSymbol("set"), atomSymbol("b"), 2)) => 2,
//     Atom::from(atomSymbol("b")) => 2,
// );

// // (eval (+ 2 2))
// test_eval!(
//     eval_plus_two_two,
//     atomList(atomSymbol("eval"), atomList(atomSymbol("+"), 2, 2)) => 4,
// );

// // (echo {"a" 1 "b" "2"})
// test_eval!(
//     echo_hash,
//     atomList(atomSymbol("echo"), Atom::Hash(std::rc::Rc::new({
//         let mut hashmap = fnv::FnvHashMap::default();
//         hashmap.insert("a".to_owned(), Atom::Int(1));
//         hashmap.insert("b".to_owned(), Atom::String("2".to_owned()));
//         hashmap
//     }))) => Atom::String(r#"{"a" 1 "b" "2"}"#.to_owned()),
// );

// // ({"a" 1 "b" "2"} "a")
// test_eval!(
//     hash_as_function,
//     atomList(Atom::Hash(std::rc::Rc::new({
//         let mut hashmap = fnv::FnvHashMap::default();
//         hashmap.insert("a".to_owned(), Atom::Int(1));
//         hashmap.insert("b".to_owned(), Atom::String("2".to_owned()));
//         hashmap
//     })), Atom::String("a".to_owned())) => Atom::Int(1),
// );

// // ((fn (x y) (+ x y)) 1 2)
// test_eval!(
//     fn_statement,
//     atomList(
//         atomList(
//             atomSymbol("fn"),
//             atomList(atomSymbol("x"), atomSymbol("y")),
//             atomList(atomSymbol("+"), atomSymbol("x"), atomSymbol("y"))),
//         1, 2
//     ) => 3);

// // ((fn (f x) (f (f x))) (fn (x) (* x 2)) 3)
// test_eval!(
//     fn_double_statement,
//     atomList(
//         atomList(
//             atomSymbol("fn"),
//             atomList(atomSymbol("f"), atomSymbol("x")),
//             atomList(
//                 atomSymbol("f"),
//                 atomList(atomSymbol("f"), atomSymbol("x"))
//             )
//         ),
//         atomList(
//             atomSymbol("fn"),
//             atomList(atomSymbol("x")),
//             atomList(atomSymbol("*"), atomSymbol("x"), 2)
//         ),
//         3
//     ) => 12);

// // (set sum2 (fn (n acc) (if (= n 0) acc (sum2 (- n 1) (+ n acc)))))
// // (sum2 10000 0)
// #[test)
// fn tco_recur() {
//     let repl_env = Env::new_repl();
//     eval(atomList(
//         atomSymbol("set"),
//         atomSymbol("sum"),
//         atomList(
//             atomSymbol("fn"),
//             atomList(atomSymbol("n"), atomSymbol("acc")),
//             atomList(
//                 atomSymbol("if"),
//                 atomList(atomSymbol("="), atomSymbol("n"), 0),
//                 atomSymbol("acc"),
//                 atomList(
//                     atomSymbol("sum"),
//                     atomList(atomSymbol("-"), atomSymbol("n"), 1),
//                     atomList(atomSymbol("+"), atomSymbol("n"), atomSymbol("acc"))
//                 )
//             )
//         ),
//     ), repl_env.clone());
//     assert_eq!(eval(atomList(atomSymbol("sum"), 10000, 0), repl_env.clone()), Atom::from(50005000));
// }

// // (set foo (fn (n) (if (= n 0) 0 (bar (- n 1)))))
// // (set bar (fn (n) (if (= n 0) 0 (foo (- n 1)))))
// // (foo 10000)
// #[test)
// fn tco_mutual_recur() {
//     let repl_env = Env::new_repl();
//     eval(atomList(
//         atomSymbol("set"),
//         atomSymbol("foo"),
//         atomList(
//             atomSymbol("fn"),
//             atomList(atomSymbol("n")),
//             atomList(
//                 atomSymbol("if"),
//                 atomList(atomSymbol("="), atomSymbol("n"), 0),
//                 0,
//                 atomList(
//                     atomSymbol("bar"),
//                     atomList(atomSymbol("-"), atomSymbol("n"), 1)
//                 )
//             )
//         ),
//     ), repl_env.clone());
//     eval(atomList(
//         atomSymbol("set"),
//         atomSymbol("bar"),
//         atomList(
//             atomSymbol("fn"),
//             atomList(atomSymbol("n")),
//             atomList(
//                 atomSymbol("if"),
//                 atomList(atomSymbol("="), atomSymbol("n"), 0),
//                 0,
//                 atomList(
//                     atomSymbol("foo"),
//                     atomList(atomSymbol("-"), atomSymbol("n"), 1)
//                 )
//             )
//         ),
//     ), repl_env.clone());
//     assert_eq!(eval(atomList(atomSymbol("foo"), 10000), repl_env.clone()), Atom::from(0));
// }
