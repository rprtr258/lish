use std::{
    fmt::Display,
    rc::Rc,
    cmp::Ordering,
    iter::{Chain, Once, once},
    vec::IntoIter,
};

use crate::env::{Env};

#[derive(Debug, Clone)]
pub struct List {
    pub head: Rc<Atom>,
    pub tail: Rc<Vec<Atom>>,
    meta: Rc<Atom>,
}

impl List {
    pub fn iter(&self) -> Chain<Once<Atom>, IntoIter<Atom>> {
        once((*self.head).clone()).chain((*self.tail).clone().into_iter())
    }
}

#[derive(Debug, Clone)]
pub enum Atom {
    Bool(bool),
    // TODO: i128?
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
    Nil,
    List(List),
}

use Atom::{Nil, Bool, Int, Float, Symbol, Lambda};

impl Atom {
    pub fn is_macro(self: &Self) -> bool {
        match self {
            Lambda {is_macro, ..} => *is_macro,
            _ => false,
        }
    }

    pub fn list(head: Atom, tail: Vec<Atom>) -> Atom {
        Atom::List(List {
            head: Rc::new(head),
            tail: Rc::new(tail),
            meta: Rc::new(Atom::Nil),
        })
    }

    pub fn symbol<S>(symbol: S) -> Atom where String: From<S> {
        Atom::Symbol(String::from(symbol))
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
impl From<Vec<Atom>> for Atom {
    fn from(mut vec: Vec<Atom>) -> Atom {
        if vec.len() == 0 {
            Nil
        } else {
            let head = vec.remove(0);
            Atom::list(head, vec)
        }
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
            (Atom::List(a), Atom::List(b)) => a.iter().partial_cmp(b.iter()),
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
            (Atom::List(ref a), Atom::List(ref b)) => a.iter().zip(b.iter()).all(|(x, y)| x == y),
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
    macro_rules! args {
        ($($val:expr),* $(,)?) => {
            vec![$(Atom::from($val), )*]
        }
    }

    #[macro_export]
    macro_rules! form {
        ($($val:expr),* $(,)?) => {
            crate::Atom::from(crate::args![$($val, )*])
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
