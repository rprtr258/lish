use std::rc::Rc;

// use fnv::FnvHashMap;

use crate::env::{Env};

#[derive(Debug, Clone)]
pub enum Atom {
    Nil,
    Bool(bool),
    Int(i64),
    Float(f64),
    String(String),
    Symbol(String),
    // Hash(Rc<FnvHashMap<String, Atom>>, Rc<Atom>),
    Func(fn(Args) -> LishRet, Rc<Atom>),
    Lambda {
        eval: fn(ast: Atom, env: Env) -> LishRet,
        ast: Rc<Atom>,
        env: Env,
        params: Rc<Atom>,
        is_macro: bool,
        meta: Rc<Atom>,
    },
    List(Rc<Vec<Atom>>, Rc<Atom>),
}

impl PartialEq for Atom {
    fn eq(&self, other: &Atom) -> bool {
        use Atom::{Nil, Bool, Int, String, Symbol, List};
        match (self, other) {
            (Nil, Nil) => true,
            (Bool(ref a), Bool(ref b)) => a == b,
            (Int(ref a), Int(ref b)) => a == b,
            (String(ref a), String(ref b)) => a == b,
            (Symbol(ref a), Symbol(ref b)) => a == b,
            (List(ref a, _), List(ref b, _)) => a == b,
            // (Hash(ref a, _), Hash(ref b, _)) => a == b,
            _ => false,
        }
    }
}

#[derive(Debug, Clone, PartialEq)]
pub enum LishErr {
    Message(String),
    // Val(Atom),
}

pub type Args = Vec<Atom>;
pub type LishRet = Result<Atom, LishErr>;

pub fn error_string(s: String) -> LishRet {
    Err(LishErr::Message(s))
}

pub fn error(s: &str) -> LishRet {
    error_string(s.to_string())
}

pub fn list(vals: Vec<Atom>) -> Atom {
    Atom::List(Rc::new(vals), Rc::new(Atom::Nil))
}