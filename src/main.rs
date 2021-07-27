use rustyline::{error::ReadlineError, Editor};

use lish::{rep};
use lish::env::{Env};

fn main() {
    let mut rl = Editor::<()>::new();
    if rl.load_history("history.txt").is_err() {
        println!("No previous history.");
    }
    let repl_env = Env::new_repl();
    loop {
        let input_buffer = rl.readline("user> ");
        match input_buffer {
            Ok(line) => {
                rl.add_history_entry(line.as_str());
                rep(&line, &repl_env);
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
