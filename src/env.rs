use std::cell::RefCell;
use std::rc::Rc;

use fnv::FnvHashMap;

use crate::types::{LishErr, LishRet, Atom, error, error_string, list};

#[derive(Debug)]
pub struct EnvStruct {
    data: RefCell<FnvHashMap<String, Atom>>,
    pub outer: Option<Env>,
}

#[derive(Debug, Clone)]
pub struct Env(Rc<EnvStruct>);

impl Env {
    pub fn new(outer: Option<Env>) -> Env {
        Env(Rc::new(EnvStruct {
            data: RefCell::new(FnvHashMap::default()),
            outer: outer,
        }))
    }

    pub fn new_repl() -> Env {
        let env = Env::new(None);
        #[macro_export]
        macro_rules! set_int_bin_op {
            ($name:expr, $init:expr, $f:expr) => {
                env.set(
                    Atom::Symbol($name.to_string()),
                    Atom::Func(|vals| vals.iter().fold(Ok(Atom::Int($init)), |a: LishRet, b: &Atom| match (a, b) {
                        (Ok(Atom::Int(ai)), Atom::Int(bi)) => Ok(Atom::Int($f(ai, bi))),
                        _ => error_string(format!("Can't eval ({} {:?})", $name, vals)),
                    }), Rc::new(Atom::Nil))).unwrap();
            };
            ($name:expr, $f:expr) => {
                env.set(
                    Atom::Symbol($name.to_string()),
                    Atom::Func(|vals| {
                        let init = vals[0].clone();
                        vals.iter().skip(1).fold(Ok(init), |a: LishRet, b: &Atom| match (a, b) {
                        (Ok(Atom::Int(ai)), Atom::Int(bi)) => Ok(Atom::Int($f(ai, bi))),
                        _ => error_string(format!("Can't eval ({} {:?})", $name, vals)),
                    })}, Rc::new(Atom::Nil))).unwrap();
            }
        }
        set_int_bin_op!("+", 0, |x, y| x + y);
        set_int_bin_op!("*", 1, |x, y| x * y);
        set_int_bin_op!("/", |x, y| x / y);
        set_int_bin_op!("-", |x, y| x - y);
        env.set(
            Atom::Symbol("-".to_string()),
            Atom::Func(|vals| match vals.len() {
                0 => error("Can't evaluate (-)"),
                1 => match vals[0] {
                    Atom::Int(x) => Ok(Atom::Int(-x)),
                    _ => error_string(format!("Can't negate {:?}", vals[0])),
                }
                _ => {
                    let init = vals[0].clone();
                    vals.iter().skip(1).fold(Ok(init), |a: LishRet, b: &Atom| match (a, b) {
                        (Ok(Atom::Int(ai)), Atom::Int(bi)) => Ok(Atom::Int(ai - bi)),
                        _ => error_string(format!("Can't eval ({} {:?})", "-", vals)),
                    })
                }
            }, Rc::new(Atom::Nil))).unwrap();
        env
    }

    pub fn bind(outer: Option<Env>, mbinds: Atom, exprs: Vec<Atom>) -> Result<Env, LishErr> {
        let env = Env::new(outer);
        match mbinds {
            Atom::List(binds, _) => {
                for (i, b) in binds.iter().enumerate() {
                    match b {
                        Atom::Symbol(s) if s == "&" => {
                            env.set(binds[i + 1].clone(), list(exprs[i..].to_vec()))?;
                            break;
                        }
                        _ => {
                            env.set(b.clone(), exprs[i].clone())?;
                        }
                    }
                }
                Ok(env)
            }
            _ => Err(LishErr::Message("Env::bind binds not List".to_string())),
        }
    }

    pub fn find(self: &Self, key: &str) -> Option<Env> {
        match (self.0.data.borrow().contains_key(key), self.0.outer.clone()) {
            (true, _) => Some(self.clone()),
            (false, Some(outer_env)) => outer_env.find(key),
            _ => None,
        }
    }

    pub fn get(self: &Self, key: &str) -> LishRet {
        match self.find(key) {
            Some(e) => Ok(e.0.data
                .borrow()
                .get(key)
                .unwrap()
                .clone()),
            _ => error(&format!("Not found '{}'", key)),
        }
    }

    pub fn sets(self: &Self, key: &str, val: Atom) {
        self.0.data.borrow_mut().insert(key.to_string(), val);
    }

    pub fn set(self: &Self, key: Atom, val: Atom) -> LishRet {
        match key {
            Atom::Symbol(ref s) => {
                self.sets(&s.to_string(), val.clone());
                Ok(val)
            }
            _ => error("Env.set called with non-Str"),
        }
    }
}
