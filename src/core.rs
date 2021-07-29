use std::rc::Rc;

use fnv::FnvHashMap;

use crate::{
    types::{Atom, LishRet, error_string, error},
    printer::{print},
};

macro_rules! set_int_bin_op {
    ($name:expr, $init:expr, $f:expr) => {
        (
            $name,
            Atom::Func(|vals| vals.iter().fold(Ok(Atom::Int($init)), |a: LishRet, b: &Atom| match (a, b) {
                (Ok(Atom::Int(ai)), Atom::Int(bi)) => Ok(Atom::Int($f(ai, bi))),
                _ => error_string(format!("Can't eval ({} {:?})", $name, vals)),
            }), Rc::new(Atom::Nil))
        )
    };
    ($name:expr, $f:expr) => {
        (
            $name,
            Atom::Func(|vals| {
                let init = vals[0].clone();
                vals.iter().skip(1).fold(Ok(init), |a: LishRet, b: &Atom| match (a, b) {
                (Ok(Atom::Int(ai)), Atom::Int(bi)) => Ok(Atom::Int($f(ai, bi))),
                _ => error_string(format!("Can't eval ({} {:?})", $name, vals)),
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
                0 => error("Can't evaluate (-)"),
                1 => match vals[0] {
                    Atom::Int(x) => Ok(Atom::Int(-x)),
                    _ => error_string(format!("Can't negate {:?}", vals[0])),
                }
                _ => {
                    let init = vals[0].clone();
                    vals.iter().skip(1).fold(Ok(init), |a: LishRet, b: &Atom| match (a, b) {
                        (Ok(Atom::Int(ai)), Atom::Int(bi)) => Ok(Atom::Int(ai - bi)),
                        _ => error_string(format!("Can't eval ({} {:?})", "-", vals)),
                    })
                }
            }, Rc::new(Atom::Nil))),
        ("prn", Atom::Func(|vals| {
            println!("{}", print(&Ok(Atom::List(Rc::new(vals), Rc::new(Atom::Nil)))));
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
                    .fold(Ok(Atom::Bool(true)), |a: LishRet, b: &Atom| match a {
                        Ok(Atom::Bool(ai)) => Ok(Atom::Bool(ai && (b.clone() == init))),
                        _ => error_string(format!("Can't eval ({} {:?})", "=", vals)),
                    })
            }, Rc::new(Atom::Nil))),
        ("<",
            Atom::Func(|vals| {
                let init = vals[0].clone();
                vals.iter()
                    .skip(1)
                    .fold(Ok(Atom::Bool(true)), |a: LishRet, b: &Atom| match a {
                        Ok(Atom::Bool(ai)) => Ok(Atom::Bool(ai && (b.clone() < init))),
                        _ => error_string(format!("Can't eval ({} {:?})", "<", vals)),
                    })
            }, Rc::new(Atom::Nil))),
        ("<=",
            Atom::Func(|vals| {
                let init = vals[0].clone();
                vals.iter()
                    .skip(1)
                    .fold(Ok(Atom::Bool(true)), |a: LishRet, b: &Atom| match a {
                        Ok(Atom::Bool(ai)) => Ok(Atom::Bool(ai && (b.clone() <= init))),
                        _ => error_string(format!("Can't eval ({} {:?})", "<=", vals)),
                    })
            }, Rc::new(Atom::Nil))),
        (">",
            Atom::Func(|vals| {
                let init = vals[0].clone();
                vals.iter()
                    .skip(1)
                    .fold(Ok(Atom::Bool(true)), |a: LishRet, b: &Atom| match a {
                        Ok(Atom::Bool(ai)) => Ok(Atom::Bool(ai && (b.clone() > init))),
                        _ => error_string(format!("Can't eval ({} {:?})", ">", vals)),
                    })
            }, Rc::new(Atom::Nil))),
        (">=",
            Atom::Func(|vals| {
                let init = vals[0].clone();
                vals.iter()
                    .skip(1)
                    .fold(Ok(Atom::Bool(true)), |a: LishRet, b: &Atom| match a {
                        Ok(Atom::Bool(ai)) => Ok(Atom::Bool(ai && (b.clone() >= init))),
                        _ => error_string(format!("Can't eval ({} {:?})", ">=", vals)),
                    })
            }, Rc::new(Atom::Nil))),
    ] {
        ns.insert(key.to_string(), val);
    }
    ns
}

#[cfg(test)]
#[allow(unused_parens)]
mod core_tests {
    use crate::{
        args,
        types::{/*error, */Atom/*, Atom::{String, Nil}*/},
        // env::{Env},
    };
    use super::{namespace};

    macro_rules! test_function {
        ($test_name:ident, $($fun:expr, $args:expr => $res:expr),* $(,)?) => {
            #[test]
            fn $test_name() {
                let ns = namespace();
                $( assert_eq!(ns.get($fun).unwrap().apply($args), Ok(Atom::from($res))); )*
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
