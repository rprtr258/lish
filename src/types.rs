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

impl From<i64> for Atom { fn from(x: i64) -> Atom { Atom::Int(x) } }
impl From<f64> for Atom { fn from(x: f64) -> Atom { Atom::Float(x) } }
impl From<bool> for Atom { fn from(x: bool) -> Atom { Atom::Bool(x) } }
impl From<&str> for Atom { fn from(x: &str) -> Atom { Atom::Symbol(x.to_string()) } }
impl<T> From<Vec<T>> for Atom where Atom: From<T>, T: Clone {
    fn from(x: Vec<T>) -> Atom {
        use crate::list_vec;
        list_vec!(x.iter().map(|x| Atom::from(x.clone())).collect())
    }
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

// TODO: remove mod?
mod macros {
    #[macro_export]
    macro_rules! symbol {
        ($name:expr) => {{
            Symbol($name.to_string())
        }}
    }

    #[macro_export]
    macro_rules! list_vec {
        ($vec:expr) => {{
            use std::rc::Rc;
            Atom::List(Rc::new($vec), Rc::new(Atom::Nil))
        }}
    }

    #[macro_export]
    macro_rules! form {
        ($($val:expr), *) => {{
            use std::rc::Rc;
            use super::{Atom};
            Atom::List(Rc::new(vec![$(Atom::from($val), )*]), Rc::new(Atom::Nil))
        }}
    }
}
