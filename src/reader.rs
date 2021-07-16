use std::iter::Iterator;
use regex::Regex;

#[derive(Clone)]
enum Atom {
    Bool(bool),
    Int(i64),
    Float(f64),
    String(String),
    Symbol(String),
}

impl std::fmt::Display for Atom {
    fn fmt(&self, fmt: &mut std::fmt::Formatter) -> std::fmt::Result {
        match self {
            Atom::Bool(b) => write!(fmt, "{}", b),
            Atom::Int(i) => write!(fmt, "{}", i),
            Atom::Float(f) => write!(fmt, "{}", f),
            Atom::String(s) => write!(fmt, "\"{}\"", s),
            Atom::Symbol(s) => write!(fmt, "{}", s),
        }
    }
}

#[derive(Clone)]
pub enum Form {
    List(Vec<Form>),
    Atom(Atom),
}

impl std::fmt::Display for Form {
    fn fmt(&self, f: &mut std::fmt::Formatter) -> std::fmt::Result {
        match self {
            Form::Atom(a) => write!(f, "{}", a.to_string()),
            Form::List(a) => {
                if a.len() == 0 {
                    return write!(f, "nil")
                }
                write!(f, "(").unwrap();
                let mut is_first = true;
                for elem in a.iter() {
                    if is_first {
                        is_first = false;
                        write!(f, "{}", elem).unwrap()
                    } else {
                        write!(f, " {}", elem).unwrap()
                    };
                };
                write!(f, ")")
            },
        }
    }
}

fn read_atom(token: String) -> Form {
    match token.parse::<bool>() {
        Ok(b) => return Form::Atom(Atom::Bool(b)),
        Err(_) => {}
    };
    match token.parse::<i64>() {
        Ok(n) => return Form::Atom(Atom::Int(n)),
        Err(_) => {}
    };
    match token.parse::<f64>() {
        Ok(x) => return Form::Atom(Atom::Float(x)),
        Err(_) => {}
    };
    if token.chars().nth(0).unwrap() == '"' {
        return Form::Atom(Atom::String(token[1..token.len()-1].to_string()))
    };
    Form::Atom(Atom::Symbol(token))
}

fn read_list<T>(tokens: &mut T) -> Form
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
            None => break
        }
    }
    Form::List(res)
}

// TODO: reader macro
fn read_form<T>(token: String, tokens: &mut T) -> Form
where T: Iterator<Item=String> {
    match &token[..] {
        "(" => read_list(tokens),
        _ => read_atom(token),
    }
}

pub fn read(cmd: &String) -> Form {
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