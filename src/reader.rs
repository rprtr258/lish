use std::iter::Iterator;

use regex::Regex;

use crate::types::{Atom, list};

fn read_atom(token: String) -> Atom {
    match token.parse::<bool>() {
        Ok(b) => return Atom::Bool(b),
        Err(_) => {}
    };
    match token.parse::<i64>() {
        Ok(n) => return Atom::Int(n),
        Err(_) => {}
    };
    match token.parse::<f64>() {
        Ok(x) => return Atom::Float(x),
        Err(_) => {}
    };
    if token.chars().nth(0).unwrap() == '"' {
        return Atom::String(token[1..token.len()-1].to_string())
    };
    Atom::Symbol(token)
}

fn read_list<T>(tokens: &mut T) -> Atom
where T: Iterator<Item=String> {
    let mut res = Vec::new();
    loop {
        match tokens.next() {
            Some(token) => {
                match &token[..] {
                    ")" => break,
                    _ => res.push(read_form(token, tokens)),
                }
            }
            None => break,
        }
    }
    match res.len() {
        0 => Atom::Nil,
        _ => list(res),
    }
}

// TODO: reader macro
fn read_form<T>(token: String, tokens: &mut T) -> Atom
where T: Iterator<Item=String> {
    match &token[..] {
        "(" => read_list(tokens),
        _ => read_atom(token),
    }
}

// TODO: add braces implicitly
pub fn read(cmd: &String) -> Atom {
    /* TODO:
    lazy_static! {
        static ref RE: Regex = Regex::new("...").unwrap();
    }
    */
    let re = Regex::new(r#"[\s]*(,@|[{}()'`,^@]|"(?:\\.|[^\\"])*"|;.*|[^\s{}('"`,;)]*)"#).unwrap();
    let mut tokens_iter = re.captures_iter(cmd)
        .map(|capture| capture[1].to_string())
        .filter(|s| s.chars().nth(0).unwrap() != ';');
    read_form(tokens_iter.next().unwrap(), &mut tokens_iter)
}

#[cfg(test)]
mod reader_tests {
    use crate::types::{Atom, list};
    use super::{read};

    #[test]
    fn nil() {
        assert_eq!(read(&"()".to_string()), Atom::Nil)
    }

    #[test]
    fn set() {
        assert_eq!(read(&"(set a 2)".to_string()), list(vec![Atom::Symbol("set".to_string()), Atom::Symbol("a".to_string()), Atom::Int(2)]));
    }
}