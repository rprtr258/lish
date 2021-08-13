use std::{
    env::args,
    rc::Rc,
};

use rustyline::{
    error::ReadlineError,
    Editor,
};

use lish::{
    rep,
    env::Env,
    types::Atom,
};

// TODO: load file from cmd args
fn main() {
    let mut rl = Editor::<()>::new();
    if rl.load_history("history.txt").is_err() {
        println!("No previous history.");
    }
    let cmd_args: Vec<String> = args().collect();

    let repl_env = Env::new_repl();
    repl_env.sets("*ARGV*", Atom::List(Rc::new(cmd_args.iter().map(|x| Atom::String(x.clone())).collect()), Rc::new(Atom::Nil)));
    // TODO: rename to load ?
    rep(r#"(set load-file (fn (f) (eval (read (str "(progn " (slurp f) "\n())")))))"#.to_string(), repl_env.clone());
    cmd_args.get(1).map(|filename|
        rep(format!(r#"(load-file "{}")"#, filename), repl_env.clone())
    );
    loop {
        let input_buffer = rl.readline("user> ");
        match input_buffer {
        Ok(line) => {
            if line == "" {
                continue;
            }
            rl.add_history_entry(line.as_str());
            let result = rep(line, repl_env.clone());
            if result != "()" {
                println!("{}", result);
            }
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
