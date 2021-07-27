use std::cell::RefCell;
use std::rc::Rc;

use fnv::FnvHashMap;

use crate::types::{LishErr, LishRet, Atom, error, list};

#[derive(Debug)]
pub struct EnvStruct {
    data: RefCell<FnvHashMap<String, Atom>>,
    pub outer: Option<Env>,
}

pub type Env = Rc<EnvStruct>;

// TODO: it would be nice to use impl here but it doesn't work on
// a deftype (i.e. Env)
pub fn env_new(outer: Option<Env>) -> Env {
    Rc::new(EnvStruct {
        data: RefCell::new(FnvHashMap::default()),
        outer: outer,
    })
}

// TODO: mbinds and exprs as & types
pub fn env_bind(outer: Option<Env>, mbinds: Atom, exprs: Vec<Atom>) -> Result<Env, LishErr> {
    let env = env_new(outer);
    match mbinds {
        Atom::List(binds, _) => {
            for (i, b) in binds.iter().enumerate() {
                match b {
                    Atom::Symbol(s) if s == "&" => {
                        env_set(&env, binds[i + 1].clone(), list(exprs[i..].to_vec()))?;
                        break;
                    }
                    _ => {
                        env_set(&env, b.clone(), exprs[i].clone())?;
                    }
                }
            }
            Ok(env)
        }
        _ => Err(LishErr::Message("env_bind binds not List/Vector".to_string())),
    }
}

pub fn env_find(env: &Env, key: &str) -> Option<Env> {
    match (env.data.borrow().contains_key(key), env.outer.clone()) {
        (true, _) => Some(env.clone()),
        (false, Some(o)) => env_find(&o, key),
        _ => None,
    }
}

pub fn env_get(env: &Env, key: &Atom) -> LishRet {
    match key {
        Atom::Symbol(ref s) => match env_find(env, s) {
            Some(e) => Ok(e
                .data
                .borrow()
                .get(s)
                .unwrap()
                .clone()),
            _ => error(&format!("'{}' not found", s)),
        },
        _ => error("Env.get called with non-Str"),
    }
}

pub fn env_set(env: &Env, key: Atom, val: Atom) -> LishRet {
    match key {
        Atom::Symbol(ref s) => {
            env.data.borrow_mut().insert(s.to_string(), val.clone());
            Ok(val)
        }
        _ => error("Env.set called with non-Str"),
    }
}

pub fn env_sets(env: &Env, key: &str, val: Atom) {
    env.data.borrow_mut().insert(key.to_string(), val);
}
