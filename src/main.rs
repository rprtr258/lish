use std::collections::HashMap;
use std::io::{stdout, Write};

use rustyline::{error::ReadlineError, Editor};

mod reader; use reader::{read, Form};
mod printer; use printer::{print};

fn eval(cmd: &Form) -> Form {
    let repl_env = {
        let mut m = HashMap::<_, fn(_) -> _>::new();
        m.insert("+", |x: Form| match x {
            Form::List(xs) => xs.iter().fold(0, |a, b| a + b)
        });
        m.insert("-", |x| x - y);
        m.insert("*", |x| x * y);
        m.insert("/", |x| x / y);
        m
    }
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
