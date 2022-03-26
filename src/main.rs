use std::env::args;

use rustyline::{
    error::ReadlineError,
    Editor,
};

use lish::{
    env::Env,
    types::Atom,
    list,
    rep,
};

fn make_repl_env(cmd_args: Vec<String>) -> Env {
    let repl_env = Env::new_repl();
    repl_env.sets(
        "*ARGV*",
        list!(cmd_args
            .iter()
            .map(Atom::from)
            .collect()
        )
    );
    // TODO: rename to load ?
    println!("{}", rep(
        r#"(set load-file (fn (f) (eval (read (join "(progn\n" (slurp f) "\n)")))))"#.to_owned(),
        repl_env.clone()
    ));
    cmd_args
        .get(1)
        .map(|filename| println!("{}", rep(
            format!(r#"(load-file "{}")"#, filename),
            repl_env.clone()
        )));
    repl_env
}

// TODO: load file from cmd args
fn main() {
    const HISTORY_FILE: &str = ".lish_history";
    let mut editor = Editor::<()>::new();
    if editor.load_history(HISTORY_FILE).is_err() {
        println!("No previous history.");
    }
    let repl_env: Env = make_repl_env(args().collect());
    loop {
        let input_buffer = editor.readline("user> ");
        match input_buffer {
            Ok(line) => {
                if line == "" {
                    continue;
                }
                editor.add_history_entry(line.as_str());
                let result = rep(line, repl_env.clone());
                if result != "()" {
                    println!("=> {}", result);
                }
            },
            Err(ReadlineError::Interrupted) => {
                println!("CTRL-C");
            },
            Err(ReadlineError::Eof) => {
                break
            },
            Err(err) => {
                println!("Error: {:?}", err);
                break
            }
        }
    }
    editor.save_history(HISTORY_FILE).unwrap();
}
