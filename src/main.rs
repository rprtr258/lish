use rustyline::{error::ReadlineError, Editor};

use lish::{
    rep,
    env::Env,
};

fn main() {
    let mut rl = Editor::<()>::new();
    if rl.load_history("history.txt").is_err() {
        println!("No previous history.");
    }
    let repl_env = Env::new_repl();
    rep(r#"(set load-file (fn (f) (eval (read (str "(progn " (slurp f) "\n())")))))"#.to_string(), repl_env.clone());
    loop {
        let input_buffer = rl.readline("user> ");
        match input_buffer {
            Ok(line) => {
                rl.add_history_entry(line.as_str());
                rep(line, repl_env.clone());
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
