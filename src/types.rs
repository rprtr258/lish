use std::{
    fmt::Display,
    rc::Rc,
    cmp::Ordering,
};

use crate::env::{Env};

#[derive(Debug, Clone)]
pub enum Atom {
    // TODO: remove Nil, cause it's the same as empty List
    Nil,
    Bool(bool),
    Int(i64),
    Float(f64),
    String(String),
    Symbol(String),
    // Hash(Rc<FnvHashMap<String, Atom>>, Rc<Atom>),
    Func(fn(Args) -> LishResult, Rc<Atom>),
    Lambda {
        eval: fn(ast: Atom, env: Env) -> LishResult,
        ast: Rc<Atom>,
        env: Env,
        // TODO: Vec<str>
        params: Rc<Atom>,
        is_macro: bool,
        meta: Rc<Atom>,
    },
    List(Rc<Vec<Atom>>, Rc<Atom>),
}

use Atom::{Nil, Bool, Int, Float, Symbol, Lambda, List};

impl Atom {
    pub fn is_macro(self: &Self) -> bool {
        match self {
        Lambda {is_macro, ..} => *is_macro,
        _ => false,
        }
    }
}

impl From<i64> for Atom {
    fn from(x: i64) -> Atom {
        Int(x)
    }
}
impl From<f64> for Atom {
    fn from(x: f64) -> Atom {
        Float(x)
    }
}
impl From<bool> for Atom {
    fn from(x: bool) -> Atom {
        Bool(x)
    }
}
impl From<&str> for Atom {
    fn from(x: &str) -> Atom {
        Atom::String(x.to_owned())
    }
}
impl From<&String> for Atom {
    fn from(x: &String) -> Atom {
        Atom::String(x.clone())
    }
}
impl<T: Clone> From<Vec<T>> for Atom
where Atom: From<T> {
    fn from(x: Vec<T>) -> Atom {
        use crate::list_vec;
        list_vec!(x.iter().map(|x| Atom::from(x.clone())).collect::<Vec<Atom>>())
    }
}

impl PartialOrd for Atom {
    fn partial_cmp(self: &Self, other: &Atom) -> Option<Ordering> {
        match (self, other) {
        (Nil, Nil) => Some(Ordering::Equal),
        (Bool(ref a), Bool(ref b)) => a.partial_cmp(b),
        (Int(ref a), Int(ref b)) => a.partial_cmp(b),
        (Atom::String(ref a), Atom::String(ref b)) => a.partial_cmp(b),
        (Symbol(ref a), Symbol(ref b)) => a.partial_cmp(b),
        (List(ref a, _), List(ref b, _)) => a.partial_cmp(b),
        // (Hash(ref a, _), Hash(ref b, _)) => a.partial_cmp(b),
        _ => None,
        }
    }
}

impl PartialEq for Atom {
    fn eq(&self, other: &Atom) -> bool {
        use Atom::String;
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
pub struct LishErr(pub String);

impl<T: Display> From<T> for LishErr {
    fn from(message: T) -> Self {
        LishErr(format!("{}", message))
    }
}

pub type LishResult = Result<Atom, LishErr>;

pub type Args = Vec<Atom>;

// TODO: remove mod?
mod macros {
    #[macro_export]
    macro_rules! symbol {
        ($name:expr) => {{
            Atom::Symbol($name.to_owned())
        }}
    }

    #[macro_export]
    macro_rules! list {
        ($vec:expr) => {{
            use std::rc::Rc;
            Atom::List(Rc::new($vec), Rc::new(Atom::Nil))
        }}
    }

    #[macro_export]
    macro_rules! list_vec {
        ($vec:expr) => {
            crate::list!(Vec::from($vec))
        }
    }

    #[macro_export]
    macro_rules! args {
        ($($val:expr),* $(,)?) => {
            vec![$(Atom::from($val), )*]
        }
    }

    #[macro_export]
    macro_rules! form {
        () => {
            Atom::Nil
        };
        ($($val:expr),* $(,)?) => {
            crate::list!(crate::args![$($val, )*])
        }
    }

    #[macro_export]
    macro_rules! func {
        ($args:ident, $body:expr) => {
            Atom::Func(|$args| {$body}, Rc::new(Atom::Nil))
        }
    }

    #[macro_export]
    macro_rules! func_ok {
        ($args:ident, $body:expr) => {
            crate::func!($args, Ok($body))
        }
    }

    #[macro_export]
    macro_rules! func_nil {
        ($args:ident, $body:expr) => {
            crate::func_ok!($args, {
                $body;
                Atom::Nil
            })
        }
    }

    #[macro_export]
    macro_rules! lisherr {
        ($message:expr) => {{
            use crate::LishErr;
            Err(LishErr::from($message))
        }};
        ($message:expr $(, $params:expr)+) => {{
            lisherr!(format!($message $(, $params)+))
        }}
    }
}
