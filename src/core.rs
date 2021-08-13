use std::{
    fs,
    rc::Rc,
};

use itertools::Itertools;
use fnv::FnvHashMap;

use crate::{
    list,
    func,
    func_ok,
    func_nil,
    lisherr,
    printer::{print, print_nice},
    reader::read,
    types::{Atom, Atom::{Nil, Int, List, Bool}, LishResult},
    env::Env,
    eval,
};

macro_rules! int_bin_op {
    ($name:expr, $init:expr, $f:expr) => {(
        $name,
        func!(
            args,
            args.iter()
                .fold(Ok(Int($init)), |a: LishResult, b: &Atom|
                    match (a, b) {
                    (Ok(Int(ai)), Int(bi)) => Ok(Int($f(ai, bi))),
                    _ => lisherr!("Can't eval ({} {:?})", $name, args),
                    }
                )
        )
    )};
    ($name:expr, $f:expr) => {(
        $name,
        func!(args, {
            let init = args[0].clone();
            args.iter()
                .skip(1)
                .fold(Ok(init), |a: LishResult, b: &Atom|
                    match (a, b) {
                    (Ok(Int(ai)), Int(bi)) => Ok(Int($f(ai, bi))),
                    _ => lisherr!("Can't eval ({} {:?})", $name, args),
                    }
                )
        })
    )}
}

macro_rules! logical_op {
    ($name:expr, $op:tt) => {(
        $name,
        func!(args, {
            let init = args[0].clone();
            args.iter()
                .skip(1)
                .fold(Ok(Bool(true)), |a: LishResult, b: &Atom|
                    match a {
                    Ok(Bool(ai)) => Ok(Bool(ai && (init $op b.clone()))),
                    _ => lisherr!("Can't eval ({} {:?})", $name, args),
                    }
                )
        })
    )};
}

// TODO: use Result.or instead of many match-es
pub fn namespace() -> FnvHashMap<String, Atom> {
    let mut ns = FnvHashMap::default();
    let cmds = vec![
        int_bin_op!("+", 0, |x, y| x + y),
        int_bin_op!("*", 1, |x, y| x * y),
        int_bin_op!("/", |x, y| x / y),
        ("-", func!(
            args,
            match args.len() {
            0 => lisherr!("Can't evaluate (-)"),
            1 => match args[0] {
                Int(x) => Ok(Int(-x)),
                _ => lisherr!("Can't negate {:?}", args[0]),
            }
            _ => {
                let init = args[0].clone();
                args.iter()
                    .skip(1)
                    .fold(Ok(init), |a: LishResult, b: &Atom|
                        match (a, b) {
                        (Ok(Int(ai)), Int(bi)) => Ok(Int(ai - bi)),
                        _ => lisherr!("Can't eval ({} {:?})", "-", args),
                })
            }
            })),
        logical_op!("=", ==),
        logical_op!("<", <),
        logical_op!("<=", <=),
        logical_op!(">", >),
        logical_op!(">=", >=),
        // TODO: remove/rename?
        ("prn", func_nil!(args,
            println!(
                "{}",
                args.into_iter()
                    .map(|x| print(&Ok(x)))
                    .join(" ")
            )
        )),
        ("echo", func_nil!(args,
            println!(
                "{}",
                args.into_iter()
                    .map(|x| print_nice(&Ok(x)))
                    .join(" ")
                )
        )),
        ("apply", func!(args, {
            let fun = args[0].clone();
            let args = args[1..].to_vec();
            // TODO: apply hashmap
            match fun {
            Atom::Func(f, _) => return f(args),
            Atom::Lambda {
                ast: lambda_ast, env: lambda_env, params, ..
            } => eval((*lambda_ast).clone(), Env::bind(Some(lambda_env.clone()), (*params).clone(), args).unwrap()),
            _ => return lisherr!("{:?} is not a function", fun),
            }
        })),
        ("cons", func_ok!(args, {
            assert!(args.len() >= 2);
            let elems = &args[..args.len()-1];
            let lst = {
                match args.last().unwrap() {
                List(xs, _) => xs.clone(),
                Nil => Rc::new(vec![]),
                _ => panic!("Trying to cons not a list"),
                }
            };
            list!(elems.iter()
                .chain(lst.iter())
                .map(|x| x.clone())
                .collect())
        })),
        // TODO: change to +
        // TODO: support Nil
        ("concat", func_ok!(
            args,
            list!(args.into_iter()
                .map(|x|
                    match x {
                    List(xs, _) => (*xs).clone(),
                    Nil => vec![],
                    _ => panic!("Trying to concat not list"),
                    })
                .flatten()
                .collect())
        )),
        ("list", func_ok!(args, list!(args))),
        ("first", func!(args, {
            assert_eq!(args.len(), 1);
            match args[0].clone() {
            List(xs, _) => Ok(xs[0].clone()),
            _ => lisherr!("Trying to get first of not list"),
            }
        })),
        ("rest", func!(args, {
            assert_eq!(args.len(), 1);
            match args[0].clone() {
            List(xs, _) => Ok(list!(xs[1..].to_vec().clone())),
            _ => lisherr!("Trying to get rest of not list"),
            }
        })),
        ("len", func!(args, {
            assert_eq!(args.len(), 1);
            match args[0].clone() {
            List(xs, _) => Ok(Int(xs.len() as i64)),
            _ => lisherr!("Trying to get len of not list"),
            }
        })),
        // TODO: test (list? ()) is true
        ("list?", func_ok!(
            args,
            Bool(
                match &args[0] {
                List(xs, _) => xs.len() > 0,
                Nil => true,
                _ => false,
                }
            )
        )),
        ("empty?", func_ok!(
            args,
            Bool(match &args[0] {
            List(xs, _) => xs.len() == 0,
            Nil => true,
            _ => false,
            })
        )),
        ("count?", func_ok!(
            args,
            match &args[0] {
            List(xs, _) => Int(xs.len() as i64),
            Nil => Int(0),
            _ => Nil
            }
        )),
        ("read", func!(args, {
            assert_eq!(args.len(), 1);
            let arg = args[0].clone();
            match arg {
            Atom::String(s) => Ok(read(s)?),
            _ => lisherr!("{:?} is not a string", arg)
            }
        })),
        ("slurp", func!(args, {
            assert_eq!(args.len(), 1);
            let arg = args[0].clone();
            match arg {
            Atom::String(filename) => match fs::read_to_string(filename) {
                Ok(s) => Ok(Atom::String(s)),
                Err(e) => return lisherr!(e),
            }
            _ => lisherr!("{:?} is not a string", arg)
            }
        })),
        ("str", func_ok!(args, {
            let result: String = args.into_iter()
                .map(|x| match x {
                Atom::String(s) => s,
                _ => print(&Ok(x)).to_owned(),
                })
                .join("");
            Atom::String(result)
        })),
    ];
    // TODO: change all &str.to_string to to_owned
    for (key, val) in cmds.into_iter() {
        ns.insert(key.to_string(), val);
    };
    ns
}

#[cfg(test)]
#[allow(unused_parens)]
mod core_tests {
    use crate::{
        lisherr,
        args,
        types::Atom,
    };
    use super::namespace;

    macro_rules! test_function {
        ($test_name:ident, $($fun:expr, $args:expr => $res:expr),* $(,)?) => {
            #[test]
            fn $test_name() {
                let ns = namespace();
                $(
                    assert_eq!(match ns.get($fun) {
                    Some(Atom::Func(f, _)) => f($args),
                    Some(_) => lisherr!("{:?} is not a function", $fun),
                    None => lisherr!("{:?} was not found", $fun),
                    }, Ok(Atom::from($res)));
                )*
            }
        }
    }

    // (*)
    test_function!(
        multiply_nullary,
        "*", args![] => 1,
    );

    // (* 2)
    test_function!(
        multiply_unary,
        "*", args![2] => 2,
    );

    // (* 1 2 3)
    test_function!(
        multiply_ternary,
        "*", args![1, 2, 3] => 6,
    );

    // (/ 1)
    test_function!(
        divide_unary,
        "/", args![1] => 1,
    );

    // (/ 5 2)
    test_function!(
        divide_binary,
        "/", args![5, 2] => 2,
    );

    // (/ 22 3 2)
    test_function!(
        divide_ternary,
        "/", args![22, 3, 2] => 3,
    );

    // (- 1)
    test_function!(
        minus_unary,
        "-", args![1] => (-1),
    );

    // (- 1 2 3)
    test_function!(
        minus_ternary,
        "-", args![1, 2, 3] => (-4),
    );

    // (+)
    test_function!(
        plus_nullary,
        "+", args![] => 0,
    );

    // (+ 1)
    test_function!(
        plus_unary,
        "+", args![1] => 1,
    );

    // (+ 1 2 3)
    test_function!(
        plus_ternary,
        "+", args![1, 2, 3] => 6
    );
}
