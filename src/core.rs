use std::rc::Rc;

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

pub fn namespace() -> Vec<(&'static str, Atom)> {
    vec![
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
    ]
}
