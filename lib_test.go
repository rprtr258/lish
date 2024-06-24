package main

import (
	"testing"

	"github.com/rprtr258/fun"
	"github.com/stretchr/testify/assert"
)

func TestQuasiquote(t *testing.T) {
	t.Skip() // TODO: not passing
	list := atomList
	symbol := atomSymbol
	for name, tc := range map[string]struct {
		ast Atom
		res Atom
	}{
		"unquote_symbol": {
			list(symbol("unquote"), symbol("a")),
			symbol("a"),
		},
		"unquote_nothing": {
			list(symbol("unquote")),
			list(
				symbol("cons"),
				list(
					symbol("quote"),
					symbol("unquote"),
				),
				list(),
			),
		},
		"unquote_many": {
			list(symbol("unquote"), symbol("a"), symbol("b"), symbol("c")),
			symbol("a"),
		},
		"splice_unquote_symbol": {
			list(list(symbol("splice-unquote"), symbol("a"))),
			list(
				symbol("concat"),
				symbol("a"),
				list(),
			),
		},
		"splice_unquote_many": {
			list(list(symbol("splice-unquote"), symbol("a"), symbol("b"), symbol("c"))),
			list(
				symbol("concat"),
				symbol("a"),
				list(),
			),
		},
		"splice_unquote_nothing": {
			list(
				list(
					symbol("splice-unquote"),
				),
			),
			list(
				symbol("cons"),
				list(
					symbol("splice-unquote"),
				),
				list(),
			),
		},
		"quasiquote_list": {
			list(symbol("a"), symbol("b"), symbol("c")),
			list(
				symbol("cons"),
				list(symbol("quote"), symbol("a")),
				list(
					symbol("cons"),
					list(symbol("quote"), symbol("b")),
					list(
						symbol("cons"),
						list(symbol("quote"), symbol("c")),
						list(),
					),
				),
			),
		},
		"quasiquote_symbol": {
			symbol("a"),
			list(symbol("quote"), symbol("a")),
		},
	} {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.res, quasiquote(tc.ast))
		})
	}
}

func TestSymbolFound(t *testing.T) {
	env := newEnv(fun.Invalid[*Env]())
	env.set("a", atomInt(1))
	assert.Equal(t, atomInt(1), eval(atomSymbol("a"), env))
}

// TODO: how to check that 'a' shell command is called?
// fn symbol_not_found() {
//     let env = Env::new(None);
//     assert.Equal(t,
//         eval(atomList(atomSymbol("a")), env),
//         lisherr!(r#""a" is not a function"#)
//     );
// }

func TestId(t *testing.T) {
	env := newEnv(fun.Invalid[*Env]())
	assert.Equal(t, atomInt(1), eval(atomInt(1), env))
}

func TestQuote_0_args(t *testing.T) {
	repl_env := newEnvRepl()
	// (quote)
	assert.Equal(t,
		lisherr(`"quote" requires 1 argument(s), but got 0 in (quote)`),
		eval(atomList(atomSymbol("quote")), repl_env),
	)
}

// // (quote a b c 1 2 3)
// #[test)
// fn quote_many_args() {
//     let repl_env = Env::new_repl();
//     assert.Equal(t,
//         eval(atomList(atomSymbol("quote"), atomSymbol("a"), atomSymbol("b"), atomSymbol("c"), 1, 2, 3), repl_env.clone()),
//         lisherr!("'quote' requires 1 argument(s), but got 6 in (quote a b c 1 2 3)")
//     );
// }

// // (quasiquoteexpand (c ,c ,@c))
// #[test)
// fn quasiquoteexpand() {
//     let repl_env = Env::new_repl();
//     assert.Equal(t,
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
//     assert.Equal(t,
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
//     assert.Equal(t,
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

func TestEval(t *testing.T) {
	for name, tc := range map[string][][2]Atom{
		// (+ 2 2)
		"list": {
			{atomList(atomSymbol("+"), atomInt(2), atomInt(2)), atomInt(4)},
		},
		// (quote a)
		"quote_symbol": {
			{atomList(atomSymbol("quote"), atomSymbol("a")), atomSymbol("a")},
		},
		// (quote 1)
		"quote_int": {
			{atomList(atomSymbol("quote"), atomInt(1)), atomInt(1)},
		},
		// // (set a 2)
		// // (+ a 3)
		// // (set b 3)
		// // (+ a b)
		// "set": {
		// 	{atomList(atomSymbol("set"), atomSymbol("a"), atomInt(2)), atomInt(2)},
		// 	{atomList(atomSymbol("+"), atomSymbol("a"), atomInt(3)), atomInt(5)},
		// 	{atomList(atomSymbol("set"), atomSymbol("b"), atomInt(3)), atomInt(3)},
		// 	{atomList(atomSymbol("+"), atomSymbol("a"), atomSymbol("b")), atomInt(5)},
		// },
		// // (set c (+ 1 2))
		// // (+ c 1)
		// "set_expr": {
		// 	{atomList(atomSymbol("set"), atomSymbol("c"), atomList(atomSymbol("+"), atomInt(1), atomInt(2))), atomInt(3)},
		// 	{atomList(atomSymbol("+"), atomSymbol("c"), atomInt(1)), atomInt(4)},
		// },
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
		//     assert.Equal(t,
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
		// "let_statement": {
		// 	{atomList(atomSymbol("set"), atomSymbol("a"), atomInt(2)), atomInt(2)},
		// 	{atomList(atomSymbol("let"), atomList(atomSymbol("a"), atomInt(1)), atomSymbol("a")), atomInt(1)},
		// 	{atomSymbol("a"), atomInt(2)},
		// },

		// // (let (a 1) 2 3 4 5 a)
		// "let_implicit_progn": {
		// 	{atomList(atomSymbol("let"), atomList(atomSymbol("a"), atomInt(1)), atomInt(2), atomInt(3), atomInt(4), atomInt(5), atomSymbol("a")), atomInt(1)},
		// },

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

		// (if true 1 2)
		"if_true_statement": {
			{atomList(atomSymbol("if"), atomBool(true), atomInt(1), atomInt(2)), atomInt(1)},
		},
		// (if false 1 2)
		"if_false_statement": {
			{atomList(atomSymbol("if"), atomBool(false), atomInt(1), atomInt(2)), atomInt(2)},
		},
		// (if true 1)
		"if_true_noelse_statement": {
			{atomList(atomSymbol("if"), atomBool(true), atomInt(1)), atomInt(1)},
		},
		// (if false 1)
		"if_false_noelse_statement": {
			{atomList(atomSymbol("if"), atomBool(false), atomInt(1)), atomNil},
		},
		// // (if true (set a 1) (set a 2))
		// // a
		// "if_set_true_statement": {
		// 	{atomList(
		// 		atomSymbol("if"), atomBool(true),
		// 		atomList(atomSymbol("set"), atomSymbol("a"), atomInt(1)),
		// 		atomList(atomSymbol("set"), atomSymbol("a"), atomInt(2)),
		// 	), atomInt(1)},
		// 	{atomSymbol("a"), atomInt(1)},
		// },
		// // (if false (set b 1) (set b 2))
		// // b
		// "if_set_false_statement": {
		// 	{atomList(
		// 		atomSymbol("if"), atomBool(false),
		// 		atomList(atomSymbol("set"), atomSymbol("b"), atomInt(1)),
		// 		atomList(atomSymbol("set"), atomSymbol("b"), atomInt(2)),
		// 	), atomInt(2)},
		// 	{atomSymbol("b"), atomInt(2)},
		// },
		// (eval (+ 2 2))
		"eval_plus_two_two": {
			{atomList(atomSymbol("eval"), atomList(atomSymbol("+"), atomInt(2), atomInt(2))), atomInt(4)},
		},
		// (echo {"a" 1 "b" "2"})
		"echo_hash": {
			{atomList(atomSymbol("echo"), atomHash(map[string]Atom{
				"a": atomInt(1),
				"b": atomString("2"),
			})), atomString(`{"a" 1 "b" "2"}`)},
		},
		// ({"a" 1 "b" "2"} "a")
		"hash_as_function": {
			{atomList(atomHash(map[string]Atom{
				"a": atomInt(1),
				"b": atomString("2"),
			}), atomString("a")), atomInt(1)},
		},

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
		// "fn_double_statement": {
		// 	{atomList(
		// 		atomList(
		// 			atomSymbol("fn"),
		// 			atomList(atomSymbol("f"), atomSymbol("x")),
		// 			atomList(
		// 				atomSymbol("f"),
		// 				atomList(atomSymbol("f"), atomSymbol("x")),
		// 			),
		// 		),
		// 		atomList(
		// 			atomSymbol("fn"),
		// 			atomList(atomSymbol("x")),
		// 			atomList(atomSymbol("*"), atomSymbol("x"), atomInt(2)),
		// 		),
		// 		atomInt(3),
		// 	), atomInt(12)},
		// },

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
		//     assert.Equal(t,eval(atomList(atomSymbol("sum"), 10000, 0), repl_env.clone()), Atom::from(50005000));
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
		//     assert.Equal(t,eval(atomList(atomSymbol("foo"), 10000), repl_env.clone()), Atom::from(0));
		// }
	} {
		t.Run(name, func(t *testing.T) {
			for _, astres := range tc {
				ast, res := astres[0], astres[1]
				assert.Equal(t, res, eval(ast, newEnvRepl()))
			}
		})
	}
}
