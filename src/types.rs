use std::rc::Rc;

// use fnv::FnvHashMap;

// use crate::env::{Env};

#[derive(Debug, Clone, PartialEq)]
pub enum Atom {
    Nil,
    Bool(bool),
    Int(i64),
    Float(f64),
    String(String),
    Symbol(String),
    // Hash(Rc<FnvHashMap<String, Atom>>, Rc<Atom>),
    Func(fn(Args) -> LishRet, Rc<Atom>),
    // Lambda {
    //     eval: fn(ast: Atom, env: Env) -> LishRet,
    //     ast: Rc<Atom>,
    //     env: Env,
    //     params: Rc<Atom>,
    //     is_macro: bool,
    //     meta: Rc<Atom>,
    // },
    List(Rc<Vec<Atom>>, Rc<Atom>),
}

#[derive(Debug, Clone)]
pub enum LishErr {
    Message(String),
    // Val(Atom),
}

pub type Args = Vec<Atom>;
pub type LishRet = Result<Atom, LishErr>;

pub fn error(s: &str) -> LishRet {
    Err(LishErr::Message(s.to_string()))
}

pub fn list(vals: Vec<Atom>) -> Atom {
    Atom::List(Rc::new(vals), Rc::new(Atom::Nil))
}