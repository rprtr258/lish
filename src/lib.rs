pub mod types;
pub mod env;
pub mod reader;
mod printer;

use crate::{
    types::{Atom, LishRet/*, error*/, error_string},
    env::{Env},
    reader::{read},
    printer::{print},
};

/*
symbol "let*":
    create a new environment using the current environment as the outer value and then use the first parameter as a list of new bindings in the "let*" environment. Take the second element of the binding list, call EVAL using the new "let*" environment as the evaluation environment, then call set on the "let*" environment using the first binding list element as the key and the evaluated second element as the value. This is repeated for each odd/even pair in the binding list. Note in particular, the bindings earlier in the list can be referred to by later bindings. Finally, the second parameter (third element) of the original let* form is evaluated using the new "let*" environment and the result is returned as the result of the let* (the new let environment is discarded upon completion). */
pub fn eval(cmd: &Atom, env: &Env) -> LishRet {
    match cmd {
        Atom::Symbol(var) => env.get(var),
        Atom::List(items, _) => {
            match &items[0] {
                Atom::Symbol(sym) => {
                    match &sym[..] {
                        "set" => {
                            assert_eq!(items.len(), 3);
                            let value: Atom = eval(&items[2], env)?;
                            env.set(items[1].clone(), value)
                        }
                        "let" => {
                            let bindings = match &items[1] {
                                Atom::List(xs, _) => xs,
                                _ => return error_string(format!("Let bindings is not a list, but a {:?}", items[1])),
                            };
                            assert_eq!(bindings.len() % 2, 0);
                            let mut i = 0;
                            let let_env = Env::new(Some(env.clone()));
                            while i < bindings.len() {
                                let var_name = bindings[i].clone();
                                let var_value = eval(&bindings[i + 1], &let_env)?;
                                let_env.set(var_name, var_value)?;
                                i += 2;
                            }
                            let body = &items[2];
                            eval(body, &let_env)
                        }
                        fun_name => {
                            let fun = env.get(fun_name)?;
                            let args: Vec<Atom> = items.iter()
                                .skip(1)
                                .map(|x| eval(x, env).unwrap())
                                .collect();
                            match fun {
                                Atom::Func(some_fun, _) => some_fun(args),
                                _ => error_string(format!("{:?} is not callable in {:?}", fun, args)),
                            }
                        }
                    }
                }
                _ => error_string(format!("{:?} is not callable in {:?}", items[0], items)),
            }
        }
        x => Ok(x.clone())
    }
}

pub fn rep(input: &String, env: &Env) {
    let result = eval(&read(input), env);
    println!("{}", print(&result));
}

// TODO: add tests from history.txt
// TODO: do fucking macro
#[cfg(test)]
mod eval_tests {
    use crate::{
        types::{list, error, Atom::{Symbol, Int, String}},
        env::{Env},
    };
    use super::{eval};

    #[test]
    fn set() {
        let repl_env = Env::new_repl();
        assert_eq!(eval(
            &list(vec![Symbol("set".to_string()), Symbol("a".to_string()), Int(2)]),
            &repl_env).unwrap(), Int(2));
        assert_eq!(eval(
            &list(vec![Symbol("+".to_string()), Symbol("a".to_string()), Int(3)]),
            &repl_env).unwrap(), Int(5));
        assert_eq!(eval(
            &list(vec![Symbol("set".to_string()), Symbol("b".to_string()), Int(3)]),
            &repl_env).unwrap(), Int(3));
        assert_eq!(eval(
            &list(vec![Symbol("+".to_string()), Symbol("a".to_string()), Symbol("b".to_string())]),
            &repl_env).unwrap(), Int(5));
        assert_eq!(eval(
            &list(vec![Symbol("set".to_string()), Symbol("c".to_string()),
                list(vec![Symbol("+".to_string()), Int(1), Int(2)])]),
            &repl_env).unwrap(), Int(3));
        assert_eq!(eval(
            &list(vec![Symbol("+".to_string()), Symbol("c".to_string()), Int(1)]),
            &repl_env).unwrap(), Int(4));
    }

    #[test]
    fn echo() {
        let repl_env = Env::new_repl();
        assert_eq!(eval(&Int(92), &repl_env).unwrap(), Int(92));
        assert_eq!(eval(&Symbol("abc".to_string()), &repl_env).err().unwrap(), error("Not found 'abc'").err().unwrap());
        assert_eq!(eval(&String("abc".to_string()), &repl_env).unwrap(), String("abc".to_string()));
    }

    #[test]
    fn multiply() {
        let repl_env = Env::new_repl();
        // (*)
        assert_eq!(eval(
            &list(vec![Symbol("*".to_string()), Int(1)]),
            &repl_env).unwrap(), Int(1));
        // (* 2)
        assert_eq!(eval(
            &list(vec![Symbol("*".to_string()), Int(2)]),
            &repl_env).unwrap(), Int(2));
        // (* 1 2 3)
        assert_eq!(eval(
            &list(vec![Symbol("*".to_string()), Int(1), Int(2), Int(3)]),
            &repl_env).unwrap(), Int(6));
    }

    #[test]
    fn divide() {
        let repl_env = Env::new_repl();
        // (/ 1)
        assert_eq!(eval(
            &list(vec![Symbol("/".to_string()), Int(1)]),
            &repl_env).unwrap(), Int(1));
        // (/ 5 2)
        assert_eq!(eval(
            &list(vec![Symbol("/".to_string()), Int(5), Int(2)]),
            &repl_env).unwrap(), Int(2));
        // (/ 22 3 2)
        assert_eq!(eval(
            &list(vec![Symbol("/".to_string()), Int(22), Int(3), Int(2)]),
            &repl_env).unwrap(), Int(3));
    }

    #[test]
    fn minus() {
        let repl_env = Env::new_repl();
        // (- 1)
        assert_eq!(eval(
            &list(vec![Symbol("-".to_string()), Int(1)]),
            &repl_env).unwrap(), Int(-1));
        // (- 1 2 3)
        assert_eq!(eval(
            &list(vec![Symbol("-".to_string()), Int(1), Int(2), Int(3)]),
            &repl_env).unwrap(), Int(-4));
    }

    #[test]
    fn plus() {
        let repl_env = Env::new_repl();
        // (+)
        assert_eq!(eval(
            &list(vec![Symbol("+".to_string())]),
            &repl_env).unwrap(), Int(0));
        // (+ 1)
        assert_eq!(eval(
            &list(vec![Symbol("+".to_string()), Int(1)]),
            &repl_env).unwrap(), Int(1));
        // (+ 1 2 3)
        assert_eq!(eval(
            &list(vec![Symbol("+".to_string()), Int(1), Int(2), Int(3)]),
            &repl_env).unwrap(), Int(6));
        // (+ 1 2 (+ 1 2))
        assert_eq!(eval(
            &list(vec![Symbol("+".to_string()), Int(1), Int(2), list(vec![Symbol("+".to_string()), Int(1), Int(2)])]),
            &repl_env).unwrap(), Int(6));
    }

    #[test]
    fn let_statement() {
        let repl_env = Env::new_repl();
        // (set a 2)
        assert_eq!(eval(
            &list(vec![Symbol("set".to_string()), Symbol("a".to_string()), Int(2)]),
            &repl_env).unwrap(), Int(2));
        // (let (a 1) a)
        assert_eq!(eval(
            &list(vec![Symbol("let".to_string()), 
                list(vec![Symbol("a".to_string()), Int(1)]),
                Symbol("a".to_string())]),
            &repl_env).unwrap(), Int(1));
        // a
        assert_eq!(eval(&Symbol("a".to_string()), &repl_env).unwrap(), Int(2));
        // (let (a 1 b 2) (+ a b))
        assert_eq!(eval(
            &list(vec![Symbol("let".to_string()), 
                list(vec![Symbol("a".to_string()), Int(1), Symbol("b".to_string()), Int(2)]),
                list(vec![Symbol("+".to_string()), Symbol("a".to_string()), Symbol("b".to_string())])]),
            &repl_env).unwrap(), Int(3));
        // (let (a 1 b a) b)
        assert_eq!(eval(
            &list(vec![Symbol("let".to_string()), 
                list(vec![Symbol("a".to_string()), Int(1), Symbol("b".to_string()), Symbol("a".to_string())]),
                Symbol("b".to_string())]),
            &repl_env).unwrap(), Int(1));
    }
}
