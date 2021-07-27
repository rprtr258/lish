use std::io::{stdout, Write};
use std::rc::Rc;

pub mod types;
mod env;
pub mod reader;
mod printer;

use crate::types::{Atom, LishRet, error};
use crate::env::{env_new, env_set, env_get};
use crate::reader::{read};
use crate::printer::{print};

pub fn eval(cmd: &Atom) -> LishRet {
    let repl_env = env_new(None);
    env_set(
        &repl_env,
        Atom::Symbol("+".to_string()),
        Atom::Func(|vals| vals.iter().fold(Ok(Atom::Int(0)), |a: LishRet, b: &Atom| match (a, b) {
            (Ok(Atom::Int(ai)), Atom::Int(bi)) => Ok(Atom::Int(ai + bi)),
            _ => error(&format!("Can't eval (+ {:?})", vals).to_string()),
        }), Rc::new(Atom::Nil)))?;
    match cmd {
        Atom::List(items, _) => {
            let form_items: Vec<Atom> = items.iter().map(|x| eval(x).unwrap()).collect();
            match env_get(&repl_env, &form_items[0]).unwrap() {
                Atom::Func(f, _) => f(form_items[1..].to_vec()),
                _ => error(&format!("Not found function {:?}", form_items[0]).to_string()),
            }
        }
        x => Ok(x.clone())
    }
}

pub fn rep(input: &String) {
    let result = eval(&read(input));
    match result {
        Ok(x) => println!("{}", print(&x)),
        Err(e) => println!("ERROR: {:?}", e),
    }
    // TODO: remove?
    match stdout().flush() {
        Ok(_) => {}
        Err(err) => {println!("Error: {:?}", err)}
    }
}
