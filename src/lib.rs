use std::rc::Rc;

pub mod types;
pub mod env;
pub mod reader;
mod printer;

use crate::{
    types::{Atom, LishRet, error_string},
    env::{Env},
    reader::{read},
    printer::{print},
};

/*
symbol: lookup the symbol in the environment structure and return the value or raise an error if no value is found
list: return a new list that is the result of calling EVAL on each of the members of the list
otherwise just return the original ast value
*/

fn eval_ast(ast: &Atom, env: &Env) -> LishRet {
    match ast {
        Atom::Symbol(var) => env.get(var),
        Atom::List(items, _) => Ok(list_vec!(items.iter().map(|x| eval(x.clone(), env.clone()).unwrap()).collect())),
        x => Ok(x.clone()),
    }
}

pub fn eval(ast: Atom, env: Env) -> LishRet {
    match ast.clone() {
        Atom::List(items, _) => {
            match &items[0] {
                Atom::Symbol(s) if s == "set" => {
                    assert_eq!(items.len(), 3);
                    let value: Atom = eval(items[2].clone(), env.clone())?;
                    env.set(items[1].clone(), value)
                }
                Atom::Symbol(s) if s == "let" => {
                    let bindings = match &items[1] {
                        Atom::List(xs, _) => xs,
                        _ => return error_string(format!("Let bindings is not a list, but a {:?}", items[1])),
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
                    let body = items[2].clone();
                    eval(body, let_env)
                }
                Atom::Symbol(s) if s == "progn" => Ok(items.iter()
                    .skip(1)
                    .map(|x| eval(x.clone(), env.clone()).unwrap())
                    .last()
                    .unwrap()),
                Atom::Symbol(s) if s == "if" => {
                    let predicate = eval(items[1].clone(), env.clone());
                    match predicate {
                        Ok(Atom::Bool(false)) | Ok(Atom::Nil) => if items.len() == 4 {
                            eval(items[3].clone(), env)
                        } else {
                            Ok(Atom::Nil)
                        },
                        _ => eval(items[2].clone(), env),
                    }
                }
                Atom::Symbol(s) if s == "fn" => {
                    let args = items[1].clone();
                    let body = items[2].clone();
                    Ok(Atom::Lambda {
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
                            match fun {
                                Atom::Func(f, _) => f(args),
                                Atom::Lambda {ast, env, params, ..} => eval((*ast).clone(), Env::bind(Some(env), (*params).clone(), args).unwrap()),
                                _ => error_string(format!("{:?} is not a function", fun)),
                            }
                        }
                        _ => unreachable!(),
                    }
                }
            }
        }
        Atom::Nil => Ok(ast.clone()),
        _ => eval_ast(&ast, &env),
    }
}

pub fn rep(input: String, env: Env) {
    let result = eval(read(input), env);
    println!("{}", print(&result));
}

// TODO: add tests from history.txt
// TODO: do fucking macro
#[cfg(test)]
mod eval_tests {
    use crate::{
        form,
        types::{error, Atom, Atom::{String, Nil}},
        env::{Env},
    };
    use super::{eval};
    macro_rules! test_eval {
        ($ast:expr, $res:expr, $env:expr) => {
            assert_eq!(eval($ast, $env.clone()), Ok(Atom::from($res)));
        }
    }

    #[test]
    fn set() {
        let repl_env = Env::new_repl();
        // (set a 2)
        test_eval!(form!["set", "a", 2], 2, repl_env);
        // (+ a 3)
        test_eval!(form!["+", "a", 3], 5, repl_env);
        // (set b 3)
        test_eval!(form!["set", "b", 3], 3, repl_env);
        // (+ a b)
        test_eval!(form!["+", "a", "b"], 5, repl_env);
        // (set c (+ 1 2))
        test_eval!(form!["set", "c", form!["+", 1, 2]], 3, repl_env);
        // (+ c 1)
        test_eval!(form!["+", "c", 1], 4, repl_env);
    }

    #[test]
    fn echo() {
        let repl_env = Env::new_repl();
        // 92
        test_eval!(Atom::from(92), 92, repl_env);
        // abc
        assert_eq!(eval(
            Atom::from("abc"),
            repl_env.clone()), error("Not found 'abc'"));
        // "abc"
        assert_eq!(eval(
            String("abc".to_string()),
            repl_env.clone()), Ok(Atom::String("abc".to_string())));
    }

    #[test]
    fn multiply() {
        let repl_env = Env::new_repl();
        // (*)
        test_eval!(form!["*"], 1, repl_env);
        // (* 2)
        test_eval!(form!["*", 2], 2, repl_env);
        // (* 1 2 3)
        test_eval!(form!["*", 1, 2, 3], 6, repl_env);
    }

    #[test]
    fn divide() {
        let repl_env = Env::new_repl();
        // (/ 1)
        test_eval!(form!["/", 1], 1, repl_env);
        // (/ 5 2)
        test_eval!(form!["/", 5, 2], 2, repl_env);
        // (/ 22 3 2)
        test_eval!(form!["/", 22, 3, 2], 3, repl_env);
    }

    #[test]
    fn minus() {
        let repl_env = Env::new_repl();
        // (- 1)
        test_eval!(form!["-", 1], -1, repl_env);
        // (- 1 2 3)
        test_eval!(form!["-", 1, 2, 3], -4, repl_env);
    }

    #[test]
    fn plus() {
        let repl_env = Env::new_repl();
        // (+)
        test_eval!(form!["+"], 0, repl_env);
        // (+ 1)
        test_eval!(form!["+", 1], 1, repl_env);
        // (+ 1 2 3)
        test_eval!(form!["+", 1, 2, 3], 6, repl_env);
        // (+ 1 2 (+ 1 2))
        test_eval!(form!["+", 1, 2, form!["+", 1, 2]], 6, repl_env);
    }

    #[test]
    fn let_statement() {
        let repl_env = Env::new_repl();
        // (set a 2)
        test_eval!(form!["set", "a", 2], 2, repl_env);
        // (let (a 1) a)
        test_eval!(form!["let", form!["a", 1], "a"], 1, repl_env);
        // a
        test_eval!(Atom::from("a"), 2, repl_env);
        // (let (a 1 b 2) (+ a b))
        test_eval!(form!["let", form!["a", 1, "b", 2], form!["+", "a", "b"]], 3, repl_env);
        // (let (a 1 b a) b)
        test_eval!(form!["let", form!["a", 1, "b", "a"], "b"], 1, repl_env);
    }

    #[test]
    fn progn_statement() {
        let repl_env = Env::new_repl();
        // (progn (set a 92) (+ a 8))
        test_eval!(form!["progn", form!["set", "a", 92], form!["+", "a", 8]], 100, repl_env);
        // a
        test_eval!(Atom::from("a"), 92, repl_env);
    }

    #[test]
    fn if_statement() {
        let repl_env = Env::new_repl();
        // (if true 1 2)
        test_eval!(form!["if", true, 1, 2], 1, repl_env);
        // (if false 1 2)
        test_eval!(form!["if", false, 1, 2], 2, repl_env);
        // (if true 1)
        test_eval!(form!["if", true, 1], 1, repl_env);
        // (if false 1)
        assert_eq!(eval(
            form!["if", false, 1],
            repl_env.clone()), Ok(Nil));
        // (if true (set a 1) (set a 2))
        test_eval!(form!["if", true, form!["set", "a", 1], form!["set", "a", 2] ], 1, repl_env);
        // a
        test_eval!(Atom::from("a"), 1, repl_env);
        // (if false (set b 1) (set b 2))
        test_eval!(form!["if", false, form!["set", "b", 1], form!["set", "b", 2] ], 2, repl_env);
        // b
        test_eval!(Atom::from("b"), 2, repl_env);
    }

    #[test]
    fn fn_statement() {
        let repl_env = Env::new_repl();
        // ((fn (x y) (+ x y)) 1 2)
        test_eval!(form![
            form![
                "fn",
                form!["x", "y"],
                form!["+", "x", "y"]],
            1, 2], 3, repl_env);
        // ((fn (f x) (f (f x))) (fn (x) (* x 2)) 3)
        test_eval!(form![
            form!["fn",
                form!["f", "x"],
                form!["f",
                    form!["f", "x"]]],
                form!["fn",
                    form!["x"],
                    form!["*", "x", 2]], 3], 12, repl_env);
    }
}
