use std::cell::RefCell;
use std::rc::Rc;

use fnv::FnvHashMap;

use crate::types::{LishErr, LishRet, Atom, error, list};

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
        env.set(
            Atom::Symbol("+".to_string()),
            Atom::Func(|vals| vals.iter().fold(Ok(Atom::Int(0)), |a: LishRet, b: &Atom| match (a, b) {
                (Ok(Atom::Int(ai)), Atom::Int(bi)) => Ok(Atom::Int(ai + bi)),
                _ => error(&format!("Can't eval (+ {:?})", vals).to_string()),
            }), Rc::new(Atom::Nil))).unwrap();
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
            _ => Err(LishErr::Message("Env::bind binds not List/Vector".to_string())),
        }
    }

    pub fn find(self: &Self, key: &str) -> Option<Env> {
        match (self.0.data.borrow().contains_key(key), self.0.outer.clone()) {
            (true, _) => Some(self.clone()),
            (false, Some(outer_env)) => outer_env.find(key),
            _ => None,
        }
    }

    pub fn get(self: &Self, key: &Atom) -> LishRet {
        match key {
            Atom::Symbol(ref s) => match self.find(s) {
                Some(e) => Ok(e.0.data
                    .borrow()
                    .get(s)
                    .unwrap()
                    .clone()),
                _ => error(&format!("'{}' not found", s)),
            },
            _ => error("Env.get called with non-Str"),
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
