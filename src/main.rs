use std::io::{stdout, Write};
use std::rc::Rc;

use rustyline::{error::ReadlineError, Editor};

mod types;
mod env;
mod reader;
mod printer;
use types::{Atom, error};
use env::{env_new, env_set, env_get};
use reader::{read};
use printer::{print};

fn eval(cmd: &Atom) -> Atom {
    let repl_env = env_new(None);
    env_set(&repl_env, Atom::Symbol("+".to_string()), Atom::Func(|vals| Ok(Atom::Int(vals.iter().map(|x| match x.clone() {
        Atom::Int(y) => y,
        _ => 0,
    }).sum())), Rc::new(Atom::Nil)));
    match cmd {
        Atom::List(items, _) => {
            match env_get(&repl_env, &items[0]).unwrap() {
                Atom::Func(f, _) => f(items[1..].to_vec()).unwrap(),
                _ => Atom::Int(0),
            }
        }
        x => x.clone()
    }
}

fn rep(input: &String) {
    println!("{}", print(&eval(&read(input))));
    match stdout().flush() {
        Ok(_) => {}
        Err(err) => {println!("Error: {:?}", err)}
    }
}

fn main() {
    let mut rl = Editor::<()>::new();
    if rl.load_history("history.txt").is_err() {
        println!("No previous history.");
    }
    loop {
        let input_buffer = rl.readline("user> ");
        match input_buffer {
            Ok(line) => {
                rl.add_history_entry(line.as_str());
                rep(&line);
            },
            Err(ReadlineError::Interrupted) => {
                println!("CTRL-C");
            },
            Err(ReadlineError::Eof) => {
                println!("CTRL-D");
                break
            },
            Err(err) => {
                println!("Error: {:?}", err);
                break
            }
        }
    }
    rl.save_history("history.txt").unwrap();
}
