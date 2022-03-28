use std::env::args;
use rustyline::{
    error::ReadlineError,
    Editor,
};
use lish::{
    env::Env,
    types::Atom,
    rep,
};

fn main() {
    const HISTORY_FILE: &str = ".lish_history";
    let mut editor = Editor::<()>::new();
    if editor.load_history(HISTORY_FILE).is_err() {
        println!("No previous history.");
    }
    let repl_env = Env::new_repl();
    let cmd_args = args().collect::<Vec<String>>();
    repl_env.sets(
        "*ARGV*",
        Atom::from(args()
            .map(|filename| Atom::from(filename.as_str()))
            .collect::<Vec<Atom>>()
        )
    );
    // TODO: rename to load ?
    // TODO: detect error
    rep(
        r#"(set load-file (fn (f) (eval (read (join "(progn\n" (slurp f) "\n)")))))"#.to_owned(),
        repl_env.clone()
    );
    if cmd_args.len() > 1 {
        cmd_args
            .get(1)
            .map(|filename| println!("{}", rep(
                format!(r#"(load-file "{}")"#, filename),
                repl_env.clone()
            )));
        return;
    }
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
