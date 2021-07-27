pub mod types;
pub mod env;
pub mod reader;
mod printer;

use crate::types::{Atom, LishRet, error};
use crate::env::{Env};
use crate::reader::{read};
use crate::printer::{print};

pub fn eval(cmd: &Atom, env: &Env) -> LishRet {
    match cmd {
        Atom::List(items, _) => {
            let form_items: Vec<Atom> = items.iter().map(|x| eval(x, env).unwrap()).collect();
            match env.get(&form_items[0]).unwrap() {
                Atom::Func(f, _) => f(form_items[1..].to_vec()),
                _ => error(&format!("Not found function {:?}", form_items[0]).to_string()),
            }
        }
        x => Ok(x.clone())
    }
}

pub fn rep(input: &String, env: &Env) {
    let result = eval(&read(input), env);
    println!("{}", print(&result));
}
