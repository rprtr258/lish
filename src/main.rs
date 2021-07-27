use std::io::{stdout, Write};

use rustyline::{error::ReadlineError, Editor};

mod types;
mod env;
mod reader;
mod printer;
use types::{Atom};
use env::{env_new};
use reader::{read};
use printer::{print};

fn eval(cmd: &Atom) -> Atom {
    let repl_env = env_new(None);
    cmd.clone()
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
