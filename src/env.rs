use std::cell::RefCell;
use std::rc::Rc;

use fnv::FnvHashMap;

use crate::{
    list_vec,
    lisherr,
    core::namespace,
    types::{LishErr, LishResult, Atom},
};

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
            outer,
        }))
    }

    pub fn new_repl() -> Env {
        let env = Env::new(None);
        for (name, fun) in namespace().iter() {
            env.sets(name, fun.clone());
        }
        env
    }

    pub fn bind(outer: Option<Env>, mbinds: Atom, exprs: Vec<Atom>) -> Result<Env, LishErr> {
        let env = Env::new(outer);
        match mbinds {
            Atom::List(binds, _) => {
                for (i, b) in binds.iter().enumerate() {
                    match b {
                        Atom::Symbol(s) if s == "&" => {
                            env.set(binds[i + 1].clone(), list_vec!(exprs[i..].to_vec()))?;
                            break;
                        }
                        _ => {
                            env.set(b.clone(), exprs[i].clone())?;
                        }
                    }
                }
                Ok(env)
            }
            Atom::Nil => Ok(env),
            _ => lisherr!("Env::bind binds not List"),
        }
    }

    pub fn find(self: &Self, key: &str) -> Option<Env> {
        if self.0.data.borrow().contains_key(key) {
            Some(self.clone())
        } else {
            self.0.outer.clone().and_then(|outer_env| outer_env.find(key))
        }
    }

    pub fn get_root(self: &Self) -> &Self {
        let mut node = self;
        while let Some(ref e) = node.0.outer {
            node = e;
        }
        node
    }

    pub fn get(self: &Self, key: &str) -> LishResult {
        match self.find(key) {
            Some(e) => Ok(e.0.data
                .borrow()
                .get(key)
                .unwrap()
                .clone()),
            _ => lisherr!(&format!("Not found '{}'", key)),
        }
    }

    pub fn sets(self: &Self, key: &str, val: Atom) {
        self.0.data.borrow_mut().insert(key.to_string(), val);
    }

    pub fn set(self: &Self, key: Atom, val: Atom) -> LishResult {
        match key {
            Atom::Symbol(ref s) => {
                self.sets(&s.to_string(), val.clone());
                Ok(val)
            }
            _ => lisherr!("Env.set called with non-Str"),
        }
    }
}
