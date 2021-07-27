use lish::{eval};
use lish::reader::{read};
use lish::types::{Atom};

#[test]
fn echo() {
    assert_eq!(eval(&read(&"92".to_string())).unwrap(), Atom::Int(92));
    assert_eq!(eval(&read(&"abc".to_string())).unwrap(), Atom::Symbol("abc".to_string()));
    assert_eq!(eval(&read(&r#""abc""#.to_string())).unwrap(), Atom::String("abc".to_string()));
}

#[test]
fn plus() {
    assert_eq!(eval(&read(&"(+ 1 2 3)".to_string())).unwrap(), Atom::Int(6));
    assert_eq!(eval(&read(&"(+ 1 2 (+ 1 2))".to_string())).unwrap(), Atom::Int(6));
}

// TODO: add tests from history.txt