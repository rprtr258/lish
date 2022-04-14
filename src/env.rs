use std::{
    cell::RefCell,
    rc::Rc,
};
use fnv::FnvHashMap;
use crate::{
    core::namespace,
    types::Atom,
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

    pub fn bind(outer: Option<Env>, binds_vec: Vec<String>, exprs: Vec<Atom>) -> Env {
        let env = Env::new(outer);
        for (i, b) in binds_vec.iter().enumerate() {
            if b == "&" {
                // TODO: List.get(index)
                env.set(&binds_vec[i + 1], Atom::from(exprs[i..].to_vec()));
                break;
            } else {
                env.set(b, exprs[i].clone());
            }
        }
        env
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

    pub fn get(self: &Self, key: &str) -> Option<Atom> {
        self.find(key)
            .map(|e| e.0.data.borrow()
                .get(key)
                .unwrap()
                .clone()
            )
    }

    pub fn sets(self: &Self, key: &str, val: Atom) {
        self.0.data.borrow_mut().insert(key.to_owned(), val);
    }

    pub fn set(self: &Self, key: &String, val: Atom) -> Atom {
        self.sets(key.as_str(), val.clone());
        val
    }
}
