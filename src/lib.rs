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

    #[test]
    fn set() {
        let repl_env = Env::new_repl();
        // (set a 2)
        assert_eq!(eval(
            form!["set", "a", 2],
            repl_env.clone()), Ok(Atom::from(2)));
        // (+ a 3)
        assert_eq!(eval(
            form!["+", "a", 3],
            repl_env.clone()), Ok(Atom::from(5)));
        // (set b 3)
        assert_eq!(eval(
            form!["set", "b", 3],
            repl_env.clone()), Ok(Atom::from(3)));
        // (+ a b)
        assert_eq!(eval(
            form!["+", "a", "b"],
            repl_env.clone()), Ok(Atom::from(5)));
        // (set c (+ 1 2))
        assert_eq!(eval(
            form!["set", "c",
                form!["+", 1, 2]],
            repl_env.clone()), Ok(Atom::from(3)));
        // (+ c 1)
        assert_eq!(eval(
            form!["+", "c", 1],
            repl_env.clone()), Ok(Atom::from(4)));
    }

    #[test]
    fn echo() {
        let repl_env = Env::new_repl();
        // 92
        assert_eq!(eval(
            Atom::from(92),
            repl_env.clone()), Ok(Atom::from(92)));
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
        assert_eq!(eval(
            form!["*", 1],
            repl_env.clone()), Ok(Atom::from(1)));
        // (* 2)
        assert_eq!(eval(
            form!["*", 2],
            repl_env.clone()), Ok(Atom::from(2)));
        // (* 1 2 3)
        assert_eq!(eval(
            form!["*", 1, 2, 3],
            repl_env.clone()), Ok(Atom::from(6)));
    }

    #[test]
    fn divide() {
        let repl_env = Env::new_repl();
        // (/ 1)
        assert_eq!(eval(
            form!["/", 1],
            repl_env.clone()), Ok(Atom::from(1)));
        // (/ 5 2)
        assert_eq!(eval(
            form!["/", 5, 2],
            repl_env.clone()), Ok(Atom::from(2)));
        // (/ 22 3 2)
        assert_eq!(eval(
            form!["/", 22, 3, 2],
            repl_env.clone()), Ok(Atom::from(3)));
    }

    #[test]
    fn minus() {
        let repl_env = Env::new_repl();
        // (- 1)
        assert_eq!(eval(
            form!["-", 1],
            repl_env.clone()), Ok(Atom::from(-1)));
        // (- 1 2 3)
        assert_eq!(eval(
            form!["-", 1, 2, 3],
            repl_env.clone()), Ok(Atom::from(-4)));
    }

    #[test]
    fn plus() {
        let repl_env = Env::new_repl();
        // (+)
        assert_eq!(eval(
            form!["+"],
            repl_env.clone()), Ok(Atom::from(0)));
        // (+ 1)
        assert_eq!(eval(
            form!["+", 1],
            repl_env.clone()), Ok(Atom::from(1)));
        // (+ 1 2 3)
        assert_eq!(eval(
            form!["+", 1, 2, 3],
            repl_env.clone()), Ok(Atom::from(6)));
        // (+ 1 2 (+ 1 2))
        assert_eq!(eval(
            form!["+", 1, 2, form!["+", 1, 2]],
            repl_env.clone()), Ok(Atom::from(6)));
    }

    #[test]
    fn let_statement() {
        let repl_env = Env::new_repl();
        // (set a 2)
        assert_eq!(eval(
            form!["set", "a", 2],
            repl_env.clone()), Ok(Atom::from(2)));
        // (let (a 1) a)
        assert_eq!(eval(
            form!["let", 
                form!["a", 1],
                "a"],
            repl_env.clone()), Ok(Atom::from(1)));
        // a
        assert_eq!(eval(
            Atom::from("a"),
            repl_env.clone()), Ok(Atom::from(2)));
        // (let (a 1 b 2) (+ a b))
        assert_eq!(eval(
            form!["let", 
                form!["a", 1, "b", 2],
                form!["+", "a", "b"]],
            repl_env.clone()), Ok(Atom::from(3)));
        // (let (a 1 b a) b)
        assert_eq!(eval(
            form!["let", 
                form!["a", 1, "b", "a"],
                "b"],
            repl_env.clone()), Ok(Atom::from(1)));
    }

    #[test]
    fn progn_statement() {
        let repl_env = Env::new_repl();
        // (progn (set a 92) (+ a 8))
        assert_eq!(eval(
            form!["progn",
                form!["set", "a", 92],
                form!["+", "a", 8]
            ],
            repl_env.clone()), Ok(Atom::from(100)));
        // a
        assert_eq!(eval(
            Atom::from("a"),
            repl_env.clone()), Ok(Atom::from(92)));
    }

    #[test]
    fn if_statement() {
        let repl_env = Env::new_repl();
        // (if true 1 2)
        assert_eq!(eval(
            form!["if", true, 1, 2],
            repl_env.clone()), Ok(Atom::from(1)));
        // (if false 1 2)
        assert_eq!(eval(
            form!["if", false, 1, 2],
            repl_env.clone()), Ok(Atom::from(2)));
        // (if true 1)
        assert_eq!(eval(
            form!["if", true, 1],
            repl_env.clone()), Ok(Atom::from(1)));
        // (if false 1)
        assert_eq!(eval(
            form!["if", false, 1],
            repl_env.clone()), Ok(Nil));
        // (if true (set a 1) (set a 2))
        assert_eq!(eval(
            form!["if", true,
                form!["set", "a", 1],
                form!["set", "a", 2]
            ],
            repl_env.clone()), Ok(Atom::from(1)));
        // a
        assert_eq!(eval(
            Atom::from("a"),
            repl_env.clone()), Ok(Atom::from(1)));
        // (if false (set b 1) (set b 2))
        assert_eq!(eval(
            form!["if", false,
                form!["set", "b", 1],
                form!["set", "b", 2]
            ],
            repl_env.clone()), Ok(Atom::from(2)));
        // b
        assert_eq!(eval(
            Atom::from("b"),
            repl_env.clone()), Ok(Atom::from(2)));
    }

    #[test]
    fn fn_statement() {
        let repl_env = Env::new_repl();
        // ((fn (x y) (+ x y)) 1 2)
        assert_eq!(eval(
            form![
                form!["fn",
                    form!["x", "y"],
                    form!["+", "x", "y"]],
                1, 2],
            repl_env.clone()), Ok(Atom::from(3)));
        // ((fn (f x) (f (f x))) (fn (x) (* x 2)) 3)
        assert_eq!(eval(
            form![
                form!["fn",
                    form!["f", "x"],
                    form!["f",
                        form!["f", "x"]]],
                form!["fn",
                    form!["x"],
                    form!["*", "x", 2]],
                3],
            repl_env), Ok(Atom::from(12)));
    }
}
