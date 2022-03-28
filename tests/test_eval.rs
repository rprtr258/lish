use lish::{eval};
use lish::env::{Env};
use lish::reader::{read};
use lish::types::Atom;

#[test]
fn set() {
    let repl_env = Env::new_repl();
    assert_eq!(eval(read("(set a 2)".to_owned()).unwrap(), repl_env.clone()), Ok(Atom::Int(2)));
    assert_eq!(eval(read("(+ a 3)".to_owned()).unwrap(), repl_env), Ok(Atom::Int(5)));
}

#[test]
fn parse_end_of_input() {
    let repl_env = Env::new_repl();
    assert_eq!(eval(read("(+ 1 2".to_owned()).unwrap(), repl_env.clone()), Ok(Atom::Int(3)));
    assert_eq!(eval(read("(+ 1 2 (+ 3 4".to_owned()).unwrap(), repl_env.clone()), Ok(Atom::Int(10)));
    assert_eq!(eval(read("+ 1 2".to_owned()).unwrap(), repl_env.clone()), Ok(Atom::Int(3)));
    assert_eq!(eval(read("+ 1 2 (+ 3 4".to_owned()).unwrap(), repl_env.clone()), Ok(Atom::Int(10)));
}

#[test]
fn echo() {
    let repl_env = Env::new_repl();
    assert_eq!(eval(read("echo 92".to_owned()).unwrap(), repl_env.clone()), Ok(Atom::String("92".to_owned())));
    // TODO: how to check "abc" is called
    // assert_eq!(eval(read("abc".to_owned()).unwrap(), repl_env.clone()), Err(LishErr::from(r#""abc" is not a function"#)));
    // assert_eq!(eval(read(r#""abc""#.to_owned()).unwrap(), repl_env), Err(LishErr::from(r#""abc" is not a function"#)));
}

#[test]
fn plus() {
    let repl_env = Env::new_repl();
    assert_eq!(eval(read("(+ 1 2 3)".to_owned()).unwrap(), repl_env.clone()), Ok(Atom::Int(6)));
    assert_eq!(eval(read("(+ 1 2 (+ 1 2))".to_owned()).unwrap(), repl_env), Ok(Atom::Int(6)));
}

// TODO: add tests from history.txt