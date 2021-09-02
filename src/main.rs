use std::env::args;

use rustyline::{
    error::ReadlineError,
    Editor,
};

use lish::{env::Env, list, rep, types::Atom};

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
    rep(
        r#"(set load-file (fn (f) (eval (read (str "(progn\n" (slurp f) "\n())")))))"#.to_owned(),
        repl_env.clone()
    );
    cmd_args
        .get(1)
        .map(|filename| rep(
            format!(r#"(load-file "{}")"#, filename),
            repl_env.clone()
        ));
    repl_env
}

fn make_editor(history_file: &str) -> Editor<()> {
    let mut rl = Editor::<()>::new();
    if rl.load_history(history_file).is_err() {
        println!("No previous history.");
    }
    rl
}

fn main_loop(rl: &mut Editor<()>, repl_env: Env) {
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
                println!("=> {}", result);
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
}

// TODO: load file from cmd args
fn main() {
    const HISTORY_FILE: &str = ".lish_history";
    let mut editor = make_editor(HISTORY_FILE);
    let repl_env: Env = make_repl_env(args().collect());
    main_loop(&mut editor, repl_env);
    editor.save_history(HISTORY_FILE).unwrap();
}
