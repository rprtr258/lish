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
symbol "def!":
    call the set method of the current environment (second parameter of EVAL called env) using the unevaluated first parameter (second list element) as the symbol key and the evaluated second parameter as the value.
symbol "let*":
    create a new environment using the current environment as the outer value and then use the first parameter as a list of new bindings in the "let*" environment. Take the second element of the binding list, call EVAL using the new "let*" environment as the evaluation environment, then call set on the "let*" environment using the first binding list element as the key and the evaluated second element as the value. This is repeated for each odd/even pair in the binding list. Note in particular, the bindings earlier in the list can be referred to by later bindings. Finally, the second parameter (third element) of the original let* form is evaluated using the new "let*" environment and the result is returned as the result of the let* (the new let environment is discarded upon completion). */
pub fn eval(cmd: &Atom, env: &Env) -> LishRet {
    match cmd {
        Atom::Symbol(var) => {
            env.get(var)
        }
        Atom::List(items, _) => {
            match &items[0] {
                Atom::Symbol(sym) => {
                    match &sym[..] {
                        "set" => {
                            env.set(items[1].clone(), items[2].clone())
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
        assert_eq!(eval(&list(vec![Symbol("set".to_string()), Symbol("a".to_string()), Int(2)]), &repl_env).unwrap(), Int(2));
        assert_eq!(eval(&list(vec![Symbol("+".to_string()), Symbol("a".to_string()), Int(3)]), &repl_env).unwrap(), Int(5));
        assert_eq!(eval(&list(vec![Symbol("set".to_string()), Symbol("b".to_string()), Int(3)]), &repl_env).unwrap(), Int(3));
        assert_eq!(eval(&list(vec![Symbol("+".to_string()), Symbol("a".to_string()), Symbol("b".to_string())]), &repl_env).unwrap(), Int(5));
    }

    #[test]
    fn echo() {
        let repl_env = Env::new_repl();
        assert_eq!(eval(&Int(92), &repl_env).unwrap(), Int(92));
        assert_eq!(eval(&Symbol("abc".to_string()), &repl_env).err().unwrap(), error("Not found 'abc'").err().unwrap());
        assert_eq!(eval(&String("abc".to_string()), &repl_env).unwrap(), String("abc".to_string()));
    }

    #[test]
    fn plus() {
        let repl_env = Env::new_repl();
        assert_eq!(eval(&list(vec![Symbol("+".to_string()), Int(1), Int(2), Int(3)]), &repl_env).unwrap(), Int(6));
        assert_eq!(eval(&list(vec![Symbol("+".to_string()), Int(1), Int(2), list(vec![Symbol("+".to_string()), Int(1), Int(2)])]), &repl_env).unwrap(), Int(6));
    }
}
