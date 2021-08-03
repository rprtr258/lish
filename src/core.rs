use std::{
    fs,
    rc::Rc,
};

use fnv::FnvHashMap;

use crate::{
    reader::read,
    types::{Atom, LishResult, LishErr},
    printer::{print, print_nice},
};

macro_rules! set_int_bin_op {
    ($name:expr, $init:expr, $f:expr) => {
        (
            $name,
            Atom::Func(|vals| vals.iter().fold(Ok(Atom::Int($init)), |a: LishResult, b: &Atom| match (a, b) {
                (Ok(Atom::Int(ai)), Atom::Int(bi)) => Ok(Atom::Int($f(ai, bi))),
                _ => Err(LishErr::from(format!("Can't eval ({} {:?})", $name, vals))),
            }), Rc::new(Atom::Nil))
        )
    };
    ($name:expr, $f:expr) => {
        (
            $name,
            Atom::Func(|vals| {
                let init = vals[0].clone();
                vals.iter().skip(1).fold(Ok(init), |a: LishResult, b: &Atom| match (a, b) {
                (Ok(Atom::Int(ai)), Atom::Int(bi)) => Ok(Atom::Int($f(ai, bi))),
                _ => Err(LishErr::from(format!("Can't eval ({} {:?})", $name, vals))),
            })}, Rc::new(Atom::Nil))
        )
    }
}

pub fn namespace() -> FnvHashMap<String, Atom> {
    let mut ns = FnvHashMap::default();
    for (key, val) in vec![
        set_int_bin_op!("+", 0, |x, y| x + y),
        set_int_bin_op!("*", 1, |x, y| x * y),
        set_int_bin_op!("/", |x, y| x / y),
        ("-", Atom::Func(|vals| match vals.len() {
                0 => Err(LishErr::from("Can't evaluate (-)")),
                1 => match vals[0] {
                    Atom::Int(x) => Ok(Atom::Int(-x)),
                    _ => Err(LishErr::from(format!("Can't negate {:?}", vals[0]))),
                }
                _ => {
                    let init = vals[0].clone();
                    vals.iter().skip(1).fold(Ok(init), |a: LishResult, b: &Atom| match (a, b) {
                        (Ok(Atom::Int(ai)), Atom::Int(bi)) => Ok(Atom::Int(ai - bi)),
                        _ => Err(LishErr::from(format!("Can't eval ({} {:?})", "-", vals))),
                    })
                }
            }, Rc::new(Atom::Nil))),
        ("prn", Atom::Func(|vals| {
            println!("{}", print(&Ok(Atom::List(Rc::new(vals), Rc::new(Atom::Nil)))));
            Ok(Atom::Nil)
        }, Rc::new(Atom::Nil))),
        ("echo", Atom::Func(|vals| {
            println!("{}", print_nice(&Ok(Atom::List(Rc::new(vals), Rc::new(Atom::Nil)))));
            Ok(Atom::Nil)
        }, Rc::new(Atom::Nil))),
        ("list", Atom::Func(|vals| Ok(Atom::List(Rc::new(vals), Rc::new(Atom::Nil))), Rc::new(Atom::Nil))),
        ("list?", Atom::Func(|vals| Ok(Atom::Bool(match &vals[0] {
            Atom::List(xs, _) => xs.len() > 0,
            Atom::Nil => true,
            _ => false,
        })), Rc::new(Atom::Nil))),
        ("empty?", Atom::Func(|vals| Ok(Atom::Bool(match &vals[0] {
            Atom::List(xs, _) => xs.len() == 0,
            Atom::Nil => true,
            _ => false,
        })), Rc::new(Atom::Nil))),
        ("count?", Atom::Func(|vals| Ok(match &vals[0] {
            Atom::List(xs, _) => Atom::Int(xs.len() as i64),
            Atom::Nil => Atom::Int(0),
            _ => Atom::Nil
        }), Rc::new(Atom::Nil))),
        ("=",
            Atom::Func(|vals| {
                let init = vals[0].clone();
                vals.iter()
                    .skip(1)
                    .fold(Ok(Atom::Bool(true)), |a: LishResult, b: &Atom| match a {
                        Ok(Atom::Bool(ai)) => Ok(Atom::Bool(ai && (b.clone() == init))),
                        _ => Err(LishErr::from(format!("Can't eval ({} {:?})", "=", vals))),
                    })
            }, Rc::new(Atom::Nil))),
        ("<",
            Atom::Func(|vals| {
                let init = vals[0].clone();
                vals.iter()
                    .skip(1)
                    .fold(Ok(Atom::Bool(true)), |a: LishResult, b: &Atom| match a {
                        Ok(Atom::Bool(ai)) => Ok(Atom::Bool(ai && (b.clone() < init))),
                        _ => Err(LishErr::from(format!("Can't eval ({} {:?})", "<", vals))),
                    })
            }, Rc::new(Atom::Nil))),
        ("<=",
            Atom::Func(|vals| {
                let init = vals[0].clone();
                vals.iter()
                    .skip(1)
                    .fold(Ok(Atom::Bool(true)), |a: LishResult, b: &Atom| match a {
                        Ok(Atom::Bool(ai)) => Ok(Atom::Bool(ai && (b.clone() <= init))),
                        _ => Err(LishErr::from(format!("Can't eval ({} {:?})", "<=", vals))),
                    })
            }, Rc::new(Atom::Nil))),
        (">",
            Atom::Func(|vals| {
                let init = vals[0].clone();
                vals.iter()
                    .skip(1)
                    .fold(Ok(Atom::Bool(true)), |a: LishResult, b: &Atom| match a {
                        Ok(Atom::Bool(ai)) => Ok(Atom::Bool(ai && (b.clone() > init))),
                        _ => Err(LishErr::from(format!("Can't eval ({} {:?})", ">", vals))),
                    })
            }, Rc::new(Atom::Nil))),
        (">=",
            Atom::Func(|vals| {
                let init = vals[0].clone();
                vals.iter()
                    .skip(1)
                    .fold(Ok(Atom::Bool(true)), |a: LishResult, b: &Atom| match a {
                        Ok(Atom::Bool(ai)) => Ok(Atom::Bool(ai && (b.clone() >= init))),
                        _ => Err(LishErr::from(format!("Can't eval ({} {:?})", ">=", vals))),
                    })
            }, Rc::new(Atom::Nil))),
        ("read",
            Atom::Func(|args| {
                assert_eq!(args.len(), 1);
                let arg = args[0].clone();
                match arg {
                    Atom::String(s) => Ok(read(s)),
                    _ => Err(LishErr::from(format!("{:?} is not a string", arg)))
                }
            }, Rc::new(Atom::Nil))),
        ("slurp",
            Atom::Func(|args| {
                assert_eq!(args.len(), 1);
                let arg = args[0].clone();
                match arg {
                    Atom::String(filename) => {
                        match fs::read_to_string(filename) {
                            Ok(s) => Ok(Atom::String(s)),
                            Err(e) => return Err(LishErr::from(e)),
                        }
                    }
                    _ => Err(LishErr::from(format!("{:?} is not a string", arg)))
                }
            }, Rc::new(Atom::Nil))),
        ("str",
            Atom::Func(|args| {
                if args.iter().any(|x| match x {
                    Atom::String(_) => false,
                    _ => true
                }) {
                    return Err(LishErr::from(format!("Can't eval ({} {:?})", "str", args)))
                }
                let result: String = args.iter()
                    .map(|x| match x {
                        Atom::String(s) => s,
                        _ => unreachable!(),
                    })
                    .flat_map(|s| s.chars())
                    .collect();
                Ok(Atom::String(result))
            }, Rc::new(Atom::Nil))),
    ] {
        // TODO: change all &str.to_string to to_owned
        ns.insert(key.to_string(), val);
    }
    ns
}

#[cfg(test)]
#[allow(unused_parens)]
mod core_tests {
    use crate::{
        args,
        types::{Atom, LishErr},
    };
    use super::{namespace};

    macro_rules! test_function {
        ($test_name:ident, $($fun:expr, $args:expr => $res:expr),* $(,)?) => {
            #[test]
            fn $test_name() {
                let ns = namespace();
                $(
                    assert_eq!(match ns.get($fun) {
                        Some(Atom::Func(f, _)) => f($args),
                        Some(_) => Err(LishErr::from(format!("{:?} is not a function", $fun))),
                        None => Err(LishErr::from(format!("{:?} was not found", $fun))),
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
