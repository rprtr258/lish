use std::rc::Rc;

pub mod types;
mod core;
pub mod env;
pub mod reader;
mod printer;

use {
    types::{Atom, List, LishResult, LishErr},
    env::Env,
    reader::read,
    printer::{print, print_debug},
};

fn eval_ast(ast: &Atom, env: &Env) -> LishResult {
    match ast {
        Atom::Symbol(var) => env.get(var),
        Atom::List(List {head, tail, ..}) => {
            let head = eval((**head).clone(), env.clone())?;
            let tail = tail.iter()
                .map(|x| eval(x.clone(), env.clone()))
                .collect::<Result<Vec<Atom>, LishErr>>()?;
            Ok(Atom::list(head, tail))
        },
        x => Ok(x.clone()),
    }
}

fn quasiquote(ast: Atom) -> Atom {
    match ast {
        Atom::List(List {head, tail, ..}) => {
            // TODO: unquote with tail.len() > 1 is meaningless
            if tail.len() >= 1 && *head == Atom::symbol("unquote") {
                tail[1].clone()
            } else {
                let mut res = vec![];
                for x in tail.iter().rev() {
                    match x {
                        Atom::List(List {head, tail, ..}) if **head == Atom::symbol("splice-unquote") => {
                            res = match tail.len() {
                                0 => vec![
                                    Atom::symbol("cons"),
                                    form![
                                        Atom::symbol("splice-unquote")
                                    ],
                                    Atom::from(res),
                                ],
                                _ => vec![
                                    Atom::symbol("concat"),
                                    Atom::from((**tail).clone()),
                                    Atom::from(res),
                                ],
                            }
                        }
                        _ => {
                            res = vec![
                                Atom::symbol("cons"),
                                quasiquote(x.clone()),
                                Atom::from(res),
                            ];
                        }
                    }
                }
                Atom::from(res)
            }
        }
        _ => Atom::from(vec![Atom::symbol("quote"), ast]),
    }
}

fn is_macro_call(ast: &Atom, env: &Env) -> bool {
    match ast {
        Atom::List(List {head, ..}) => {
            match &**head {
                Atom::Symbol(macroname) => env.get(&macroname)
                    .map(|x| x.is_macro())
                    .unwrap_or(false),
                _ => false,
            }
        }
        _ => false,
    }
}

fn macroexpand(mut ast: Atom, env: &Env) -> Result<Atom, LishErr> {
    while is_macro_call(&ast, env) {
        let macro_call = ast;
        match macro_call {
            Atom::List(List {head, tail, ..}) => {
                let the_macro = eval((*head).clone(), env.clone()).unwrap();
                let args = tail.to_vec();
                match the_macro {
                    Atom::Lambda {ast: lambda_ast, env: lambda_env, params, ..} => {
                        let lambda_ast = (*lambda_ast).clone();
                        let lambda_env = Env::bind(Some(lambda_env.clone()), (*params).clone(), args).unwrap();
                        ast = eval(lambda_ast, lambda_env)?;
                    },
                    _ => unreachable!(),
                }
            }
            _ => unreachable!(),
        }
    }
    Ok(ast)
}

pub fn eval(mut ast: Atom, mut env: Env) -> LishResult {
    loop {
        ast = macroexpand(ast, &env)?;
        match ast.clone() {
            Atom::List(List {head, tail, ..}) => {
                macro_rules! lish_assert_args {
                    ($cmd:expr, $args_count:expr) => {{
                        if tail.len() != $args_count {
                            return lisherr!("'{}' requires {} argument(s), but got {} in {}", $cmd, $args_count, tail.len(), print(&Ok(ast.clone())));
                        }
                    }}
                }

                match &*head {
                    Atom::Symbol(s) if s == "quote" => {
                        lish_assert_args!("quote", 1);
                        return Ok(Atom::from((*tail).clone()))
                    }
                    Atom::Symbol(s) if s == "quasiquoteexpand" => {
                        lish_assert_args!("quasiquoteexpand", 1);
                        return Ok(quasiquote(Atom::from((*tail).clone())));
                    }
                    Atom::Symbol(s) if s == "quasiquote" => {
                        lish_assert_args!("quasiquote", 1);
                        ast = quasiquote(Atom::from((*tail).clone()));
                        continue
                    }
                    Atom::Symbol(s) if s == "macroexpand" => {
                        lish_assert_args!("macroexpand", 1);
                        match tail[0].clone() {
                            Atom::List(List {head, ..}) => {eval((*head).clone(), env.clone())?;},
                            _ => unreachable!(),
                        }
                        return macroexpand(tail[1].clone(), &env);
                    }
                    Atom::Symbol(s) if s == "set" => {
                        lish_assert_args!("set", 2);
                        let value: Atom = eval(tail[1].clone(), env.clone())?;
                        return env.set(tail[0].clone(), value)
                    }
                    Atom::Symbol(s) if s == "setmacro" => {
                        lish_assert_args!("setmacro", 2);
                        return match eval(tail[1].clone(), env.clone())? {
                            Atom::Lambda {
                                eval, ast, env, params, meta, ..
                            } => env.set(tail[0].clone(), Atom::Lambda {
                                eval,
                                ast,
                                params,
                                meta,
                                env: env.clone(),
                                is_macro: true,
                            }),
                            _ => lisherr!("Macro is not lambda"),
                        }
                    }
                    Atom::Symbol(s) if s == "let" => {
                        let bindings = match &tail[0] {
                            Atom::List(binds) => binds.iter().collect::<Vec<Atom>>(),
                            _ => return lisherr!("Let bindings is not a list, but a {:?}", tail[0]),
                        };
                        if bindings.len() % 2 != 0 {
                            return lisherr!("'let' requires even number of arguments, but got {} in {}", tail.len(), print_debug(&Ok(ast.clone())));
                        }
                        let mut i = 0;
                        let let_env = Env::new(Some(env.clone()));
                        while i < bindings.len() {
                            let var_name = bindings[i].clone();
                            let var_value = eval(bindings[i + 1].clone(), let_env.clone())?;
                            let_env.set(var_name, var_value)?;
                            i += 2;
                        }
                        let mut body = Vec::with_capacity(tail.len());
                        body.push(Atom::symbol("progn"));
                        body.extend_from_slice(&tail[2..]);
                        ast = Atom::from(body);
                        env = let_env;
                    }
                    Atom::Symbol(s) if s == "progn" => {
                        for item in tail.iter() {
                            eval(item.clone(), env.clone())?;
                        }
                        ast = tail.last().unwrap().clone()
                    }
                    Atom::Symbol(s) if s == "if" => {
                        let predicate = eval(Atom::from(tail[0].clone()), env.clone())?;
                        match predicate {
                            Atom::Bool(false) | Atom::Nil => if tail.len() == 3 {
                                ast = tail[2].clone()
                            } else {
                                return Ok(Atom::Nil)
                            },
                            _ => {
                                ast = tail[1].clone()
                            }
                        }
                    }
                    Atom::Symbol(s) if s == "eval" => {
                        lish_assert_args!("eval", 1);
                        ast = eval(tail[0].clone(), env.clone())?;
                        env = env.get_root().clone();
                        continue;
                    }
                    Atom::Symbol(s) if s == "fn" => {
                        let args = tail[0].clone();
                        let body = tail[1].clone();
                        return Ok(Atom::Lambda {
                            eval: eval,
                            ast: Rc::new(body),
                            env: env.clone(),
                            params: Rc::new(args),
                            is_macro: false,
                            meta: Rc::new(Atom::Nil),
                        })
                    }
                    _ => {
                        let evaluated_list = eval_ast(&ast, &env)?;
                        match evaluated_list {
                            Atom::List(List {head, tail, ..}) => {
                                let fun = (*head).clone();
                                let args = tail.to_vec();
                                // TODO: apply hashmap
                                match fun {
                                    Atom::Func(f, _) => return f(args),
                                    Atom::Lambda {ast: lambda_ast, env: lambda_env, params, ..} => {
                                        ast = (*lambda_ast).clone();
                                        env = Env::bind(Some(lambda_env.clone()), (*params).clone(), args).unwrap();
                                    },
                                    _ => return lisherr!("{:?} is not a function", fun),
                                }
                            }
                            _ => unreachable!(),
                        }
                    }
                }
            }
            Atom::Nil => return Ok(Atom::Nil),
            _ => return eval_ast(&ast, &env),
        }
    }
}

pub fn rep(input: String, env: Env) -> String {
    print_debug(&read(input).and_then(|ast| eval(ast, env)))
}

#[cfg(test)]
#[allow(unused_parens)]
mod eval_tests {
    mod eval_ast_tests {
        use crate::{
            lisherr,
            form,
            types::{Atom, List},
            env::Env,
            eval_ast,
        };

        #[test]
        fn symbol_found() {
            let env = Env::new(None);
            env.sets("a", Atom::Int(1));
            assert_eq!(eval_ast(&Atom::symbol("a"), &env), Ok(Atom::Int(1)));
        }

        #[test]
        fn symbol_not_found() {
            let env = Env::new(None);
            assert_eq!(
                eval_ast(&form![Atom::symbol("a")], &env),
                lisherr!("Not found 'a'")
            );
        }

        #[test]
        fn list() {
            let repl_env = Env::new_repl();
            let res = eval_ast(&form![Atom::symbol("+"), 2, 2], &repl_env);
            match res {
                Ok(Atom::List(List {head, tail, ..})) => {
                    assert_eq!(tail[0], Atom::Int(2));
                    assert_eq!(tail[1], Atom::Int(2));
                    match *head {
                        Atom::Func(_, _) => (),
                        _ => unreachable!(),
                    }
                }
                _ => unreachable!(),
            }
        }

        #[test]
        fn id() {
            let env = Env::new(None);
            assert_eq!(eval_ast(&Atom::Int(1), &env), Ok(Atom::Int(1)));
        }
    }

    mod quasiquote_tests {
        use crate::{
            form,
            types::Atom,
            quasiquote,
        };

        macro_rules! test_quasiquote {
            ($test_name:ident, $ast:expr, $res:expr) => {
                #[test]
                fn $test_name() {
                    assert_eq!(quasiquote($ast), $res);
                }
            }
        }

        test_quasiquote!(
            unquote_symbol,
            form![Atom::symbol("unquote"), Atom::symbol("a")],
            Atom::symbol("a")
        );

        test_quasiquote!(
            unquote_nothing,
            form![Atom::symbol("unquote")],
            form![
                Atom::symbol("cons"),
                form![
                    Atom::symbol("quote"),
                    Atom::symbol("unquote"),
                ],
                Vec::<Atom>::new(),
            ]
        );

        test_quasiquote!(
            unquote_many,
            form![Atom::symbol("unquote"), Atom::symbol("a"), Atom::symbol("b"), Atom::symbol("c")],
            Atom::symbol("a")
        );

        test_quasiquote!(
            splice_unquote_symbol,
            form![form![Atom::symbol("splice-unquote"), Atom::symbol("a")]],
            form![
                Atom::symbol("concat"),
                Atom::symbol("a"),
                Vec::<Atom>::new(),
            ]
        );

        test_quasiquote!(
            splice_unquote_many,
            form![form![Atom::symbol("splice-unquote"), Atom::symbol("a"), Atom::symbol("b"), Atom::symbol("c")]],
            form![
                Atom::symbol("concat"),
                Atom::symbol("a"),
                Vec::<Atom>::new(),
            ]
        );

        test_quasiquote!(
            splice_unquote_nothing,
            form![
                form![
                    Atom::symbol("splice-unquote")
                ]
            ],
            form![
                Atom::symbol("cons"),
                form![
                    Atom::symbol("splice-unquote")
                ],
                Vec::<Atom>::new(),
            ]
        );

        test_quasiquote!(
            quasiquote_list,
            form![Atom::symbol("a"), Atom::symbol("b"), Atom::symbol("c")],
            form![
                Atom::symbol("cons"),
                form![Atom::symbol("quote"), Atom::symbol("a")],
                form![
                    Atom::symbol("cons"),
                    form![Atom::symbol("quote"), Atom::symbol("b")],
                    form![
                        Atom::symbol("cons"),
                        form![Atom::symbol("quote"), Atom::symbol("c")],
                        Vec::<Atom>::new(),
                    ]
                ],
            ]
        );

        test_quasiquote!(
            quasiquote_symbol,
            Atom::symbol("a"),
            form![Atom::symbol("quote"), Atom::symbol("a")]
        );
    }

    use crate::{
        lisherr,
        form,
        env::Env,
        types::{Atom, Atom::Nil},
        eval,
    };

    macro_rules! test_eval {
        ($test_name:ident, $($ast:expr => $res:expr),* $(,)?) => {
            #[test]
            fn $test_name() {
                let repl_env = Env::new_repl();
                $(
                    assert_eq!(eval($ast, repl_env.clone()), Ok(Atom::from($res)));
                )*
            }
        }
    }

    // (quote a)
    test_eval!(
        quote_symbol,
        form![Atom::symbol("quote"), Atom::symbol("a")] => Atom::symbol("a"),
    );

    // (quote 1)
    test_eval!(
        quote_int,
        form![Atom::symbol("quote"), 1] => 1,
    );

    // (quote)
    #[test]
    fn quote_0_args() {
        let repl_env = Env::new_repl();
        assert_eq!(
            eval(form![Atom::symbol("quote")], repl_env.clone()),
            lisherr!("'quote' requires 1 argument(s), but got 0 in (quote)")
        );
    }

    // (quote a b c 1 2 3)
    #[test]
    fn quote_many_args() {
        let repl_env = Env::new_repl();
        assert_eq!(
            eval(form![Atom::symbol("quote"), Atom::symbol("a"), Atom::symbol("b"), Atom::symbol("c"), 1, 2, 3], repl_env.clone()),
            lisherr!("'quote' requires 1 argument(s), but got 6 in (quote a b c 1 2 3)")
        );
    }

    // (quasiquoteexpand (c ,c ,@c))
    #[test]
    fn quasiquoteexpand() {
        let repl_env = Env::new_repl();
        assert_eq!(
            eval(form![Atom::symbol("quasiquoteexpand"), form![Atom::symbol("c"), form![Atom::symbol("unquote"), Atom::symbol("c")], form![Atom::symbol("splice-unquote"), Atom::symbol("c")]]], repl_env.clone()),
            Ok(form![
                Atom::symbol("cons"),
                form![
                    Atom::symbol("quote"), Atom::symbol("c")
                ],
                form![
                    Atom::symbol("cons"),
                    Atom::symbol("c"),
                    form![
                        Atom::symbol("concat"),
                        Atom::symbol("c"),
                        Vec::<Atom>::new()
                    ]
                ]
            ]),
        );
    }

    // (set c '(1 2 3))
    // `(c ,c ,@c)
    #[test]
    fn quasiquote() {
        let repl_env = Env::new_repl();
        eval(form![
            Atom::symbol("set"),
            Atom::symbol("c"),
            form![
                Atom::symbol("quote"),
                form![1, 2, 3],
            ],
        ], repl_env.clone()).unwrap();
        assert_eq!(
            eval(form![
                Atom::symbol("quasiquote"),
                form![
                    Atom::symbol("c"),
                    form![
                        Atom::symbol("unquote"),
                        Atom::symbol("c"),
                    ],
                    form![
                        Atom::symbol("splice-unquote"),
                        Atom::symbol("c"),
                    ],
                ],
            ], repl_env.clone()),
            Ok(form![
                Atom::symbol("c"),
                form![1, 2, 3],
                1, 2, 3,
            ])
        );
    }
    
    // (setmacro m (fn (x) `(,x ,x)))
    // (macroexpand (m (Y f)))
    #[test]
    fn set_macro_then_macroexpand() {
        let repl_env = Env::new_repl();
        // TODO: remove all unwraps
        eval(form![
            Atom::symbol("setmacro"),
            Atom::symbol("m"),
            form![
                Atom::symbol("fn"),
                form![Atom::symbol("x")],
                form![
                    Atom::symbol("quasiquote"),
                    form![
                        form![Atom::symbol("unquote"), Atom::symbol("x")],
                        form![Atom::symbol("unquote"), Atom::symbol("x")]
                    ]
                ]
            ]
        ], repl_env.clone()).unwrap();
        assert_eq!(
            eval(form![
                Atom::symbol("macroexpand"),
                form![
                    Atom::symbol("m"),
                    form![
                        Atom::symbol("Y"),
                        Atom::symbol("f"),
                    ]
                ],
            ], repl_env.clone()),
            Ok(form![
                form![Atom::symbol("Y"), Atom::symbol("f")],
                form![Atom::symbol("Y"), Atom::symbol("f")]
            ])
        );
    }

    // (set a 2)
    // (+ a 3)
    // (set b 3)
    // (+ a b)
    test_eval!(
        set,
        form![Atom::symbol("set"), Atom::symbol("a"), 2] => 2,
        form![Atom::symbol("+"), Atom::symbol("a"), 3] => 5,
        form![Atom::symbol("set"), Atom::symbol("b"), 3] => 3,
        form![Atom::symbol("+"), Atom::symbol("a"), Atom::symbol("b")] => 5,
    );

    // (set c (+ 1 2))
    // (+ c 1)
    test_eval!(
        set_expr,
        form![Atom::symbol("set"), Atom::symbol("c"), form![Atom::symbol("+"), 1, 2]] => 3,
        form![Atom::symbol("+"), Atom::symbol("c"), 1] => 4,
    );

    // 92
    test_eval!(
        eval_int,
        Atom::from(92) => 92,
    );

    // abc
    #[test]
    fn eval_symbol() {
        let repl_env = Env::new_repl();
        assert_eq!(
            // TODO: Error for not found symbol
            eval(Atom::symbol("abc"), repl_env.clone()),
            lisherr!("Not found 'abc'")
        );
    }

    // "abc"
    #[test]
    fn eval_string() {
        let repl_env = Env::new_repl();
        assert_eq!(
            eval(Atom::from("abc"), repl_env.clone()),
            Ok(Atom::from("abc"))
        );
    }

    // (+ 1 2 (+ 1 2))
    test_eval!(
        plus_expr,
        form![
            Atom::symbol("+"),
            1, 2,
            form![
                Atom::symbol("+"),
                1, 2
            ]
        ] => 6,
    );

    // (set a 2)
    // (let (a 1) a)
    // a
    test_eval!(
        let_statement,
        form![Atom::symbol("set"), Atom::symbol("a"), 2] => 2,
        form![Atom::symbol("let"), form![Atom::symbol("a"), 1], Atom::symbol("a")] => 1,
        Atom::from(Atom::symbol("a")) => 2,
    );

    // (let (a 1) 2 3 4 5 a)
    test_eval!(
        let_implicit_progn,
        form![Atom::symbol("let"), form![Atom::symbol("a"), 1], 2, 3, 4, 5, Atom::symbol("a")] => 1,
    );

    // (let (a 1 b 2) (+ a b))
    test_eval!(
        let_twovars_statement,
        form![Atom::symbol("let"), form![Atom::symbol("a"), 1, Atom::symbol("b"), 2], form![Atom::symbol("+"), Atom::symbol("a"), Atom::symbol("b")]] => 3,
    );

    // (let (a 1 b a) b)
    test_eval!(
        let_star_statement,
        form![Atom::symbol("let"), form![Atom::symbol("a"), 1, Atom::symbol("b"), Atom::symbol("a")], Atom::symbol("b")] => 1,
    );

    // (progn (set a 92) (+ a 8))
    // a
    test_eval!(
        progn_statement,
        form![Atom::symbol("progn"), form![Atom::symbol("set"), Atom::symbol("a"), 92], form![Atom::symbol("+"), Atom::symbol("a"), 8]] => 100,
        Atom::from(Atom::symbol("a")) => 92,
    );

    // (if true 1 2)
    test_eval!(
        if_true_statement,
        form![Atom::symbol("if"), true, 1, 2] => 1,
    );

    // (if false 1 2)
    test_eval!(
        if_false_statement,
        form![Atom::symbol("if"), false, 1, 2] => 2,
    );

    // (if true 1)
    test_eval!(
        if_true_noelse_statement,
        form![Atom::symbol("if"), true, 1] => 1,
    );

    // (if false 1)
    test_eval!(
        if_false_noelse_statement,
        form![Atom::symbol("if"), false, 1] => Nil,
    );

    // (if true (set a 1) (set a 2))
    // a
    test_eval!(
        if_set_true_statement,
        form![Atom::symbol("if"), true, form![Atom::symbol("set"), Atom::symbol("a"), 1], form![Atom::symbol("set"), Atom::symbol("a"), 2]] => 1,
        Atom::from(Atom::symbol("a")) => 1,
    );

    // (if false (set b 1) (set b 2))
    // b
    test_eval!(
        if_set_false_statement,
        form![Atom::symbol("if"), false, form![Atom::symbol("set"), Atom::symbol("b"), 1], form![Atom::symbol("set"), Atom::symbol("b"), 2]] => 2,
        Atom::from(Atom::symbol("b")) => 2,
    );

    // (eval (+ 2 2))
    test_eval!(
        eval_plus_two_two,
        form![Atom::symbol("eval"), form![Atom::symbol("+"), 2, 2]] => 4,
    );

    // ((fn (x y) (+ x y)) 1 2)
    test_eval!(
        fn_statement,
        form![
            form![
                Atom::symbol("fn"),
                form![Atom::symbol("x"), Atom::symbol("y")],
                form![Atom::symbol("+"), Atom::symbol("x"), Atom::symbol("y")]],
            1, 2
        ] => 3);

    // ((fn (f x) (f (f x))) (fn (x) (* x 2)) 3)
    test_eval!(
        fn_double_statement,
        form![
            form![
                Atom::symbol("fn"),
                form![Atom::symbol("f"), Atom::symbol("x")],
                form![
                    Atom::symbol("f"),
                    form![Atom::symbol("f"), Atom::symbol("x")]
                ]
            ],
            form![
                Atom::symbol("fn"),
                form![Atom::symbol("x")],
                form![Atom::symbol("*"), Atom::symbol("x"), 2]
            ],
            3
        ] => 12);

    // (set sum2 (fn (n acc) (if (= n 0) acc (sum2 (- n 1) (+ n acc)))))
    // (sum2 10000 0)
    #[test]
    fn tco_recur() {
        let repl_env = Env::new_repl();
        eval(form![
            Atom::symbol("set"),
            Atom::symbol("sum"),
            form![
                Atom::symbol("fn"),
                form![Atom::symbol("n"), Atom::symbol("acc")],
                form![
                    Atom::symbol("if"),
                    form![Atom::symbol("="), Atom::symbol("n"), 0],
                    Atom::symbol("acc"),
                    form![
                        Atom::symbol("sum"),
                        form![Atom::symbol("-"), Atom::symbol("n"), 1],
                        form![Atom::symbol("+"), Atom::symbol("n"), Atom::symbol("acc")]
                    ]
                ]
            ],
        ], repl_env.clone()).unwrap();
        assert_eq!(eval(form![Atom::symbol("sum"), 10000, 0], repl_env.clone()), Ok(Atom::from(50005000)));
    }

    // (set foo (fn (n) (if (= n 0) 0 (bar (- n 1)))))
    // (set bar (fn (n) (if (= n 0) 0 (foo (- n 1)))))
    // (foo 10000)
    #[test]
    fn tco_mutual_recur() {
        let repl_env = Env::new_repl();
        eval(form![
            Atom::symbol("set"),
            Atom::symbol("foo"),
            form![
                Atom::symbol("fn"),
                form![Atom::symbol("n")],
                form![
                    Atom::symbol("if"),
                    form![Atom::symbol("="), Atom::symbol("n"), 0],
                    0,
                    form![
                        Atom::symbol("bar"),
                        form![Atom::symbol("-"), Atom::symbol("n"), 1]
                    ]
                ]
            ],
        ], repl_env.clone()).unwrap();
        eval(form![
            Atom::symbol("set"),
            Atom::symbol("bar"),
            form![
                Atom::symbol("fn"),
                form![Atom::symbol("n")],
                form![
                    Atom::symbol("if"),
                    form![Atom::symbol("="), Atom::symbol("n"), 0],
                    0,
                    form![
                        Atom::symbol("foo"),
                        form![Atom::symbol("-"), Atom::symbol("n"), 1]
                    ]
                ]
            ],
        ], repl_env.clone()).unwrap();
        assert_eq!(eval(form![Atom::symbol("foo"), 10000], repl_env.clone()), Ok(Atom::from(0)));
    }
}
