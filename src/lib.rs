use std::rc::Rc;

pub mod types;
mod core;
pub mod env;
pub mod reader;
mod printer;

use crate::{
    types::{Atom, LishResult, LishErr},
    env::Env,
    reader::read,
    printer::print,
};

fn eval_ast(ast: &Atom, env: &Env) -> LishResult {
    match ast {
    Atom::Symbol(var) => env.get(var),
    Atom::List(items, _) => {
        let list = items.iter()
            .map(|x| eval(x.clone(), env.clone()))
            .collect::<Result<Vec<Atom>, LishErr>>()?;
        Ok(list!(list))
    },
    x => Ok(x.clone()),
    }
}

fn quasiquote(ast: Atom) -> Atom {
    match ast {
    Atom::List(xs, _) => {
        if xs.len() >= 2 && xs[0] == symbol!("unquote") {
            xs[1].clone()
        } else {
            let mut res = vec![];
            for x in xs.iter().rev() {
                match x {
                Atom::List(ys, _) if ys[0] == symbol!("splice-unquote") => {
                    res = match ys.len() {
                    1 => vec![
                        symbol!("cons"),
                        form![
                            symbol!("splice-unquote")
                        ],
                        list!(res),
                    ],
                    _ => vec![
                        symbol!("concat"),
                        ys[1].clone(),
                        list!(res),
                    ],
                    }
                }
                _ => {
                    res = vec![
                        symbol!("cons"),
                        quasiquote(x.clone()),
                        list!(res),
                    ];
                }
                }
            }
            list!(res)
        }
    }
    _ => list!(vec![symbol!("quote"), ast]),
    }
}

fn is_macro_call(ast: &Atom, env: &Env) -> bool {
    match ast {
    Atom::List(xs, _) => {
        if xs.len() == 0 {
            return false;
        }
        match &xs[0] {
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
        Atom::List(fun_call, _) => {
            let the_macro = eval(fun_call[0].clone(), env.clone()).unwrap();
            let args = fun_call[1..].to_vec();
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
    ast = macroexpand(ast, &env)?;
    loop {
        match ast.clone() {
        Atom::List(items, _) => {
            if items.len() == 0 {
                return Ok(Atom::Nil)
            }
            match &items[0] {
            Atom::Symbol(s) if s == "quote" => {
                if items.len() != 2 {
                    return lisherr!("'quote' requires 1 argument, but got {} in {}", items.len() - 1, print(&Ok(ast.clone())));
                }
                return Ok(items[1].clone())
            }
            Atom::Symbol(s) if s == "quasiquoteexpand" => {
                // TODO: change assert to meaningful message, add tests
                assert_eq!(items.len(), 2);
                return Ok(quasiquote(items[1].clone()));
            }
            Atom::Symbol(s) if s == "quasiquote" => {
                assert_eq!(items.len(), 2);
                ast = quasiquote(items[1].clone());
                continue
            }
            Atom::Symbol(s) if s == "macroexpand" => {
                assert_eq!(items.len(), 2);
                match items[1].clone() {
                Atom::List(xs, _) => {eval(xs[0].clone(), env.clone())?;},
                _ => unreachable!(),
                }
                return macroexpand(items[1].clone(), &env);
            }
            Atom::Symbol(s) if s == "set" => {
                assert_eq!(items.len(), 3);
                let value: Atom = eval(items[2].clone(), env.clone())?;
                return env.set(items[1].clone(), value)
            }
            Atom::Symbol(s) if s == "setmacro" => {
                assert_eq!(items.len(), 3);
                return match eval(items[2].clone(), env.clone())? {
                Atom::Lambda {
                    eval, ast, env, params, meta, ..
                } => env.set(items[1].clone(), Atom::Lambda {
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
                let bindings = match &items[1] {
                Atom::List(xs, _) => xs,
                _ => return lisherr!("Let bindings is not a list, but a {:?}", items[1]),
                };
                assert_eq!(bindings.len() % 2, 0);
                let mut i = 0;
                let let_env = Env::new(Some(env.clone()));
                while i < bindings.len() {
                    let var_name = bindings[i].clone();
                    let var_value = eval(bindings[i + 1].clone(), let_env.clone())?;
                    let_env.set(var_name, var_value)?;
                    i += 2;
                }
                let mut body = Vec::with_capacity(items.len() - 2 + 1);
                body.push(symbol!("progn"));
                body.extend_from_slice(&items[2..]);
                ast = list_vec![body];
                env = let_env;
            }
            Atom::Symbol(s) if s == "progn" => {
                let body_items = items.len() - 1;
                for item in &items[1..body_items] {
                    eval(item.clone(), env.clone())?;
                }
                ast = items[body_items].clone()
            }
            Atom::Symbol(s) if s == "if" => {
                let predicate = eval(items[1].clone(), env.clone())?;
                match predicate {
                Atom::Bool(false) | Atom::Nil => if items.len() == 4 {
                    ast = items[3].clone()
                } else {
                    return Ok(Atom::Nil)
                },
                _ => {
                    ast = items[2].clone()
                }
                }
            }
            Atom::Symbol(s) if s == "eval" => {
                assert_eq!(items.len(), 2);
                ast = eval(items[1].clone(), env.clone())?;
                env = env.get_root().clone();
                continue;
            }
            Atom::Symbol(s) if s == "fn" => {
                let args = items[1].clone();
                let body = items[2].clone();
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
                Atom::List(fun_call, _) => {
                    let fun = fun_call[0].clone();
                    let args = fun_call[1..].to_vec();
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
    let result = read(input).and_then(|ast| eval(ast, env));
    // TODO: add "=> " before result
    print(&result)
}

#[cfg(test)]
#[allow(unused_parens)]
mod eval_tests {
    mod eval_ast_tests {
        use crate::{
            lisherr,
            form,
            symbol,
            types::Atom,
            env::Env,
            eval_ast,
        };

        #[test]
        fn symbol_found() {
            let env = Env::new(None);
            env.sets("a", Atom::Int(1));
            assert_eq!(eval_ast(&symbol!("a"), &env), Ok(Atom::Int(1)));
        }

        #[test]
        fn symbol_not_found() {
            let env = Env::new(None);
            assert_eq!(
                eval_ast(&form![symbol!("a")], &env),
                lisherr!("Not found 'a'")
            );
        }

        #[test]
        fn list() {
            let repl_env = Env::new_repl();
            let res = eval_ast(&form![symbol!("+"), 2, 2], &repl_env);
            match res {
            Ok(Atom::List(items, _)) => {
                assert_eq!(items[1], Atom::Int(2));
                assert_eq!(items[2], Atom::Int(2));
                match items[0] {
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
            symbol,
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
            form![symbol!("unquote"), symbol!("a")],
            symbol!("a")
        );

        test_quasiquote!(
            unquote_nothing,
            form![symbol!("unquote")],
            form![
                symbol!("cons"),
                form![
                    symbol!("quote"),
                    symbol!("unquote"),
                ],
                Vec::<Atom>::new(),
            ]
        );

        test_quasiquote!(
            unquote_many,
            form![symbol!("unquote"), symbol!("a"), symbol!("b"), symbol!("c")],
            symbol!("a")
        );

        test_quasiquote!(
            splice_unquote_symbol,
            form![form![symbol!("splice-unquote"), symbol!("a")]],
            form![
                symbol!("concat"),
                symbol!("a"),
                Vec::<Atom>::new(),
            ]
        );

        test_quasiquote!(
            splice_unquote_many,
            form![form![symbol!("splice-unquote"), symbol!("a"), symbol!("b"), symbol!("c")]],
            form![
                symbol!("concat"),
                symbol!("a"),
                Vec::<Atom>::new(),
            ]
        );

        test_quasiquote!(
            splice_unquote_nothing,
            form![
                form![
                    symbol!("splice-unquote")
                ]
            ],
            form![
                symbol!("cons"),
                form![
                    symbol!("splice-unquote")
                ],
                Vec::<Atom>::new(),
            ]
        );

        test_quasiquote!(
            quasiquote_list,
            form![symbol!("a"), symbol!("b"), symbol!("c")],
            form![
                symbol!("cons"),
                form![symbol!("quote"), symbol!("a")],
                form![
                    symbol!("cons"),
                    form![symbol!("quote"), symbol!("b")],
                    form![
                        symbol!("cons"),
                        form![symbol!("quote"), symbol!("c")],
                        Vec::<Atom>::new(),
                    ]
                ],
            ]
        );

        test_quasiquote!(
            quasiquote_symbol,
            symbol!("a"),
            form![symbol!("quote"), symbol!("a")]
        );
    }

    use crate::{
        lisherr,
        form,
        symbol,
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
        form![symbol!("quote"), symbol!("a")] => symbol!("a"),
    );

    // (quote 1)
    test_eval!(
        quote_int,
        form![symbol!("quote"), 1] => 1,
    );

    // (quote)
    #[test]
    fn quote_0_args() {
        let repl_env = Env::new_repl();
        assert_eq!(
            eval(form![symbol!("quote")], repl_env.clone()),
            lisherr!("'quote' requires 1 argument, but got 0 in (quote)")
        );
    }

    // (quote a b c 1 2 3)
    #[test]
    fn quote_many_args() {
        let repl_env = Env::new_repl();
        assert_eq!(
            eval(form![symbol!("quote"), symbol!("a"), symbol!("b"), symbol!("c"), 1, 2, 3], repl_env.clone()),
            lisherr!("'quote' requires 1 argument, but got 6 in (quote a b c 1 2 3)")
        );
    }

    // (quasiquoteexpand (c ,c ,@c))
    #[test]
    fn quasiquoteexpand() {
        let repl_env = Env::new_repl();
        assert_eq!(
            eval(form![symbol!("quasiquoteexpand"), form![symbol!("c"), form![symbol!("unquote"), symbol!("c")], form![symbol!("splice-unquote"), symbol!("c")]]], repl_env.clone()),
            Ok(form![
                symbol!("cons"),
                form![
                    symbol!("quote"), symbol!("c")
                ],
                form![
                    symbol!("cons"),
                    symbol!("c"),
                    form![
                        symbol!("concat"),
                        symbol!("c"),
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
            symbol!("set"),
            symbol!("c"),
            form![
                symbol!("quote"),
                form![1, 2, 3],
            ],
        ], repl_env.clone()).unwrap();
        assert_eq!(
            eval(form![
                symbol!("quasiquote"),
                form![
                    symbol!("c"),
                    form![
                        symbol!("unquote"),
                        symbol!("c"),
                    ],
                    form![
                        symbol!("splice-unquote"),
                        symbol!("c"),
                    ],
                ],
            ], repl_env.clone()),
            Ok(form![
                symbol!("c"),
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
            symbol!("setmacro"),
            symbol!("m"),
            form![
                symbol!("fn"),
                form![symbol!("x")],
                form![
                    symbol!("quasiquote"),
                    form![
                        form![symbol!("unquote"), symbol!("x")],
                        form![symbol!("unquote"), symbol!("x")]
                    ]
                ]
            ]
        ], repl_env.clone()).unwrap();
        assert_eq!(
            eval(form![
                symbol!("macroexpand"),
                form![
                    symbol!("m"),
                    form![
                        symbol!("Y"),
                        symbol!("f"),
                    ]
                ],
            ], repl_env.clone()),
            Ok(form![
                form![symbol!("Y"), symbol!("f")],
                form![symbol!("Y"), symbol!("f")]
            ])
        );
    }

    // (set a 2)
    // (+ a 3)
    // (set b 3)
    // (+ a b)
    test_eval!(
        set,
        form![symbol!("set"), symbol!("a"), 2] => 2,
        form![symbol!("+"), symbol!("a"), 3] => 5,
        form![symbol!("set"), symbol!("b"), 3] => 3,
        form![symbol!("+"), symbol!("a"), symbol!("b")] => 5,
    );

    // (set c (+ 1 2))
    // (+ c 1)
    test_eval!(
        set_expr,
        form![symbol!("set"), symbol!("c"), form![symbol!("+"), 1, 2]] => 3,
        form![symbol!("+"), symbol!("c"), 1] => 4,
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
            eval(symbol!("abc"), repl_env.clone()),
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
            symbol!("+"),
            1, 2,
            form![
                symbol!("+"),
                1, 2
            ]
        ] => 6,
    );

    // (set a 2)
    // (let (a 1) a)
    // a
    test_eval!(
        let_statement,
        form![symbol!("set"), symbol!("a"), 2] => 2,
        form![symbol!("let"), form![symbol!("a"), 1], symbol!("a")] => 1,
        Atom::from(symbol!("a")) => 2,
    );

    // (let (a 1) 2 3 4 5 a)
    test_eval!(
        let_implicit_progn,
        form![symbol!("let"), form![symbol!("a"), 1], 2, 3, 4, 5, symbol!("a")] => 1,
    );

    // (let (a 1 b 2) (+ a b))
    test_eval!(
        let_twovars_statement,
        form![symbol!("let"), form![symbol!("a"), 1, symbol!("b"), 2], form![symbol!("+"), symbol!("a"), symbol!("b")]] => 3,
    );

    // (let (a 1 b a) b)
    test_eval!(
        let_star_statement,
        form![symbol!("let"), form![symbol!("a"), 1, symbol!("b"), symbol!("a")], symbol!("b")] => 1,
    );

    // (progn (set a 92) (+ a 8))
    // a
    test_eval!(
        progn_statement,
        form![symbol!("progn"), form![symbol!("set"), symbol!("a"), 92], form![symbol!("+"), symbol!("a"), 8]] => 100,
        Atom::from(symbol!("a")) => 92,
    );

    // (if true 1 2)
    test_eval!(
        if_true_statement,
        form![symbol!("if"), true, 1, 2] => 1,
    );

    // (if false 1 2)
    test_eval!(
        if_false_statement,
        form![symbol!("if"), false, 1, 2] => 2,
    );

    // (if true 1)
    test_eval!(
        if_true_noelse_statement,
        form![symbol!("if"), true, 1] => 1,
    );

    // (if false 1)
    test_eval!(
        if_false_noelse_statement,
        form![symbol!("if"), false, 1] => Nil,
    );

    // (if true (set a 1) (set a 2))
    // a
    test_eval!(
        if_set_true_statement,
        form![symbol!("if"), true, form![symbol!("set"), symbol!("a"), 1], form![symbol!("set"), symbol!("a"), 2]] => 1,
        Atom::from(symbol!("a")) => 1,
    );

    // (if false (set b 1) (set b 2))
    // b
    test_eval!(
        if_set_false_statement,
        form![symbol!("if"), false, form![symbol!("set"), symbol!("b"), 1], form![symbol!("set"), symbol!("b"), 2]] => 2,
        Atom::from(symbol!("b")) => 2,
    );

    // (eval (+ 2 2))
    test_eval!(
        eval_plus_two_two,
        form![symbol!("eval"), form![symbol!("+"), 2, 2]] => 4,
    );

    // ((fn (x y) (+ x y)) 1 2)
    test_eval!(
        fn_statement,
        form![
            form![
                symbol!("fn"),
                form![symbol!("x"), symbol!("y")],
                form![symbol!("+"), symbol!("x"), symbol!("y")]],
            1, 2
        ] => 3);

    // ((fn (f x) (f (f x))) (fn (x) (* x 2)) 3)
    test_eval!(
        fn_double_statement,
        form![
            form![
                symbol!("fn"),
                form![symbol!("f"), symbol!("x")],
                form![
                    symbol!("f"),
                    form![symbol!("f"), symbol!("x")]
                ]
            ],
            form![
                symbol!("fn"),
                form![symbol!("x")],
                form![symbol!("*"), symbol!("x"), 2]
            ],
            3
        ] => 12);

    // (set sum2 (fn (n acc) (if (= n 0) acc (sum2 (- n 1) (+ n acc)))))
    // (sum2 10000 0)
    #[test]
    fn tco_recur() {
        let repl_env = Env::new_repl();
        eval(form![
            symbol!("set"),
            symbol!("sum"),
            form![
                symbol!("fn"),
                form![symbol!("n"), symbol!("acc")],
                form![
                    symbol!("if"),
                    form![symbol!("="), symbol!("n"), 0],
                    symbol!("acc"),
                    form![
                        symbol!("sum"),
                        form![symbol!("-"), symbol!("n"), 1],
                        form![symbol!("+"), symbol!("n"), symbol!("acc")]
                    ]
                ]
            ],
        ], repl_env.clone()).unwrap();
        assert_eq!(eval(form![symbol!("sum"), 10000, 0], repl_env.clone()), Ok(Atom::from(50005000)));
    }

    // (set foo (fn (n) (if (= n 0) 0 (bar (- n 1)))))
    // (set bar (fn (n) (if (= n 0) 0 (foo (- n 1)))))
    // (foo 10000)
    #[test]
    fn tco_mutual_recur() {
        let repl_env = Env::new_repl();
        eval(form![
            symbol!("set"),
            symbol!("foo"),
            form![
                symbol!("fn"),
                form![symbol!("n")],
                form![
                    symbol!("if"),
                    form![symbol!("="), symbol!("n"), 0],
                    0,
                    form![
                        symbol!("bar"),
                        form![symbol!("-"), symbol!("n"), 1]
                    ]
                ]
            ],
        ], repl_env.clone()).unwrap();
        eval(form![
            symbol!("set"),
            symbol!("bar"),
            form![
                symbol!("fn"),
                form![symbol!("n")],
                form![
                    symbol!("if"),
                    form![symbol!("="), symbol!("n"), 0],
                    0,
                    form![
                        symbol!("foo"),
                        form![symbol!("-"), symbol!("n"), 1]
                    ]
                ]
            ],
        ], repl_env.clone()).unwrap();
        assert_eq!(eval(form![symbol!("foo"), 10000], repl_env.clone()), Ok(Atom::from(0)));
    }
}
