use lish::{eval};
use lish::env::{Env};
use lish::reader::{read};
use lish::types::{Atom, LishErr};

#[test]
fn set() {
    let repl_env = Env::new_repl();
    assert_eq!(eval(read("(set a 2)".to_string()).unwrap(), repl_env.clone()), Ok(Atom::Int(2)));
    assert_eq!(eval(read("(+ a 3)".to_string()).unwrap(), repl_env), Ok(Atom::Int(5)));
}

#[test]
fn parse_end_of_input() {
    let repl_env = Env::new_repl();
    assert_eq!(eval(read("(+ 1 2".to_string()).unwrap(), repl_env.clone()), Ok(Atom::Int(3)));
    assert_eq!(eval(read("(+ 1 2 (+ 3 4".to_string()).unwrap(), repl_env.clone()), Ok(Atom::Int(10)));
}

#[test]
fn echo() {
    let repl_env = Env::new_repl();
    assert_eq!(eval(read("92".to_string()).unwrap(), repl_env.clone()), Ok(Atom::Int(92)));
    assert_eq!(eval(read("abc".to_string()).unwrap(), repl_env.clone()), Err(LishErr::from("Not found 'abc'")));
    assert_eq!(eval(read(r#""abc""#.to_string()).unwrap(), repl_env), Ok(Atom::String("abc".to_string())));
}

#[test]
fn plus() {
    let repl_env = Env::new_repl();
    assert_eq!(eval(read("(+ 1 2 3)".to_string()).unwrap(), repl_env.clone()), Ok(Atom::Int(6)));
    assert_eq!(eval(read("(+ 1 2 (+ 1 2))".to_string()).unwrap(), repl_env), Ok(Atom::Int(6)));
}

// TODO: add tests from history.txt