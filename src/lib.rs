use std::rc::Rc;

pub mod types;
mod core;
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

pub fn eval(_ast: Atom, _env: Env) -> LishRet {
    let mut ast = _ast.clone();
    let mut env = _env.clone();
    loop {
        match ast.clone() {
            Atom::List(items, _) => {
                match &items[0] {
                    Atom::Symbol(s) if s == "set" => {
                        assert_eq!(items.len(), 3);
                        let value: Atom = eval(items[2].clone(), env.clone())?;
                        return env.set(items[1].clone(), value)
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
                        ast = items[2].clone();
                        env = let_env;
                    }
                    Atom::Symbol(s) if s == "progn" => {
                        let body_items = items.len() - 1;
                        for i in 1..body_items {
                            eval(items[i].clone(), env.clone()).unwrap();
                        }
                        ast = items[body_items].clone()
                    }
                    Atom::Symbol(s) if s == "if" => {
                        let predicate = eval(items[1].clone(), env.clone());
                        match predicate {
                            Ok(Atom::Bool(false)) | Ok(Atom::Nil) => if items.len() == 4 {
                                ast = items[3].clone()
                            } else {
                                return Ok(Atom::Nil)
                            },
                            _ => {
                                ast = items[2].clone()
                            }
                        }
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
                                return fun.apply(args)
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

pub fn rep(input: String, env: Env) {
    let result = eval(read(input), env);
    println!("{}", print(&result));
}

// TODO: add tests from history.txt
#[cfg(test)]
#[allow(unused_parens)]
mod eval_tests {
    use crate::{
        form,
        types::{error, Atom, Atom::{String, Nil}},
        env::{Env},
    };
    use super::{eval};

    macro_rules! test_eval {
        ($test_name:ident, $($ast:expr => $res:expr),* $(,)?) => {
            #[test]
            fn $test_name() {
                let repl_env = Env::new_repl();
                $( assert_eq!(eval($ast, repl_env.clone()), Ok(Atom::from($res))); )*
            }
        }
    }

    // (set a 2)
    // (+ a 3)
    // (set b 3)
    // (+ a b)
    test_eval!(
        set,
        form!["set", "a", 2] => 2,
        form!["+", "a", 3] => 5,
        form!["set", "b", 3] => 3,
        form!["+", "a", "b"] => 5,
    );

    // (set c (+ 1 2))
    // (+ c 1)
    test_eval!(
        set_expr,
        form!["set", "c", form!["+", 1, 2]] => 3,
        form!["+", "c", 1] => 4,
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
        assert_eq!(eval(
            Atom::from("abc"),
            repl_env.clone()), error("Not found 'abc'"));
    }

    // "abc"
    #[test]
    fn eval_string() {
        let repl_env = Env::new_repl();
        assert_eq!(eval(
            String("abc".to_string()),
            repl_env.clone()), Ok(Atom::String("abc".to_string())));
    }

    // (+ 1 2 (+ 1 2))
    test_eval!(
        plus_expr,
        form!["+", 1, 2, form!["+", 1, 2]] => 6,
    );

    // (set a 2)
    // (let (a 1) a)
    // a
    test_eval!(
        let_statement,
        form!["set", "a", 2] => 2,
        form!["let", form!["a", 1], "a"] => 1,
        Atom::from("a") => 2,
    );

    // (let (a 1 b 2) (+ a b))
    test_eval!(
        let_twovars_statement,
        form!["let", form!["a", 1, "b", 2], form!["+", "a", "b"]] => 3,
    );

    // (let (a 1 b a) b)
    test_eval!(
        let_star_statement,
        form!["let", form!["a", 1, "b", "a"], "b"] => 1,
    );

    // (progn (set a 92) (+ a 8))
    // a
    test_eval!(
        progn_statement,
        form!["progn", form!["set", "a", 92], form!["+", "a", 8]] => 100,
        Atom::from("a") => 92,
    );

    // (if true 1 2)
    test_eval!(
        if_true_statement,
        form!["if", true, 1, 2] => 1,
    );

    // (if false 1 2)
    test_eval!(
        if_false_statement,
        form!["if", false, 1, 2] => 2,
    );

    // (if true 1)
    test_eval!(
        if_true_noelse_statement,
        form!["if", true, 1] => 1,
    );

    // (if false 1)
    test_eval!(
        if_false_noelse_statement,
        form!["if", false, 1] => Nil,
    );

    // (if true (set a 1) (set a 2))
    // a
    test_eval!(
        if_set_true_statement,
        form!["if", true, form!["set", "a", 1], form!["set", "a", 2]] => 1,
        Atom::from("a") => 1,
    );

    // (if false (set b 1) (set b 2))
    // b
    test_eval!(
        if_set_false_statement,
        form!["if", false, form!["set", "b", 1], form!["set", "b", 2]] => 2,
        Atom::from("b") => 2,
    );

    // ((fn (x y) (+ x y)) 1 2)
    test_eval!(
        fn_statement,
        form![
            form![
                "fn",
                form!["x", "y"],
                form!["+", "x", "y"]],
            1, 2] => 3);

    // ((fn (f x) (f (f x))) (fn (x) (* x 2)) 3)
    test_eval!(
        fn_double_statement,
        form![
            form!["fn",
                form!["f", "x"],
                form!["f",
                    form!["f", "x"]]],
                form!["fn",
                    form!["x"],
                    form!["*", "x", 2]], 3] => 12);
}
