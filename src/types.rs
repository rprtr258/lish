use std::{
    rc::Rc,
    cmp::Ordering,
    iter::{Chain, Once, once},
    vec::IntoIter,
};
use fnv::FnvHashMap;
use crate::env::Env;

#[derive(Debug, Clone)]
pub struct List {
    pub head: Rc<Atom>,
    pub tail: Rc<Vec<Atom>>,
    // meta: Rc<Atom>,
}

impl List {
    pub fn new(head: Atom, tail: Vec<Atom>) -> List {
        List {
            head: Rc::new(head),
            tail: Rc::new(tail),
            // meta: Rc::new(Atom::Nil),
        }
    }

    pub fn iter(&self) -> Chain<Once<Atom>, IntoIter<Atom>> {
        once((*self.head).clone()).chain((*self.tail).clone().into_iter())
    }

    pub fn len(&self) -> usize {
        1 + self.tail.len()
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
    Error(String),
    Hash(Rc<FnvHashMap<String, Atom>>),//, Rc<Atom>),
    Func(fn(Vec<Atom>) -> Atom, Rc<Atom>),
    Lambda {
        eval: fn(ast: Atom, env: Env) -> Atom,
        ast: Rc<Atom>,
        env: Env,
        params: Vec<String>,
        is_macro: bool,
        // meta: Rc<Atom>,
    },
    Nil,
    List(List),
}

impl Atom {
    pub fn is_macro(self: &Self) -> bool {
        match self {
            Self::Lambda {is_macro, ..} => *is_macro,
            _ => false,
        }
    }

    pub fn list(head: Self, tail: Vec<Self>) -> Self {
        Self::List(List::new(head, tail))
    }

    pub fn symbol<S>(symbol: S) -> Self where String: From<S> {
        Self::Symbol(String::from(symbol))
    }
}

impl From<i64> for Atom {
    fn from(x: i64) -> Self {
        Self::Int(x)
    }
}
impl From<f64> for Atom {
    fn from(x: f64) -> Self {
        Self::Float(x)
    }
}
impl From<bool> for Atom {
    fn from(x: bool) -> Self {
        Self::Bool(x)
    }
}
impl From<&str> for Atom {
    fn from(x: &str) -> Self {
        Self::String(x.to_owned())
    }
}
impl From<&String> for Atom {
    fn from(x: &String) -> Self {
        Self::String(x.clone())
    }
}
impl From<Vec<Atom>> for Atom {
    // TODO: remove mut
    fn from(mut vec: Vec<Atom>) -> Self {
        if vec.len() == 0 {
            Self::Nil
        } else {
            let head = vec.remove(0);
            Self::list(head, vec)
        }
    }
}

impl PartialOrd for Atom {
    fn partial_cmp(self: &Self, other: &Self) -> Option<Ordering> {
        match (self, other) {
            (Self::Nil, Self::Nil) => Some(Ordering::Equal),
            (Self::Bool(ref a), Self::Bool(ref b)) => a.partial_cmp(b),
            (Self::Int(ref a), Self::Int(ref b)) => a.partial_cmp(b),
            (Self::String(ref a), Self::String(ref b)) => a.partial_cmp(b),
            (Self::Symbol(ref a), Self::Symbol(ref b)) => a.partial_cmp(b),
            (Self::List(a), Self::List(b)) => a.iter().partial_cmp(b.iter()),
            _ => None,
        }
    }
}

impl PartialEq for Atom {
    fn eq(&self, other: &Self) -> bool {
        match (self, other) {
            (Self::Nil, Self::Nil) => true,
            (Self::Bool(ref a), Self::Bool(ref b)) => a == b,
            (Self::Int(ref a), Self::Int(ref b)) => a == b,
            (Self::String(ref a), Self::String(ref b)) => a == b,
            (Self::Error(ref a), Self::Error(ref b)) => a == b,
            (Self::Symbol(ref a), Self::Symbol(ref b)) => a == b,
            (Self::List(ref a), Self::List(ref b)) => a.iter().zip(b.iter()).all(|(x, y)| x == y),
            (Self::Hash(ref a), Self::Hash(ref b)) => a == b,
            _ => false,
        }
    }
}

// TODO: remove mod?
mod macros {
    #[macro_export]
    macro_rules! lish_try {
        ($val:expr) => {
            match $val {
                err@Atom::Error(_) => return err,
                x => x,
            }
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
        ($($val:expr),* $(,)?) => {
            crate::Atom::from(crate::args![$($val, )*])
        }
    }

    #[macro_export]
    macro_rules! lisherr {
        ($message:expr) => {{
            Atom::Error($message.to_string())
        }};
        ($message:expr $(, $params:expr)+) => {{
            crate::lisherr!(format!($message $(, $params)+))
        }}
    }

    #[macro_export]
    macro_rules! assert_symbol {
        ($e:expr) => {
            match &$e {
                Atom::Symbol(identifier) => identifier,
                x => return crate::lisherr!("{} is not a symbol", print_debug(&x)),
            }
        }
    }
}
