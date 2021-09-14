use itertools::Itertools;

use crate::types::{Atom, LishResult};

// TODO: print atom, not result
pub fn print_nice(val: &LishResult) -> String {
    match val {
    Ok(x) => match x {
        Atom::Nil => "()".to_string(),
        Atom::Bool(y) => format!("{}", y),
        Atom::Int(y) => format!("{}", y),
        Atom::Float(y) => format!("{}", y),
        Atom::String(y) => format!("{}", y),
        Atom::Symbol(y) => format!("{}", y),
        Atom::Func(_, _) => "#fn".to_string(),
        Atom::Lambda {ast, params, is_macro, ..} => {
            let params_str = match (**params).clone() {
            Atom::List(arg_names, _) => arg_names.iter()
                .map(|x|
                    match x {
                    Atom::Symbol(arg_name) => arg_name,
                    _ => panic!("Lambda arg is not symbol"),
                    }
                )
                .join(" "),
            Atom::Nil => "()".to_string(),
            _ => panic!("Lambda args is not list"),
            };
            format!("({} ({}) {})", if *is_macro {"defmacro"} else {"fn"}, params_str, print_nice(&Ok((**ast).clone())))
        },
        Atom::List(items, _) => format!("({})", items.iter()
            .map(|x| print_nice(&Ok(x.clone())))
            .join(" ")
        ),
    }
    Err(e) => format!("ERROR: {:?}", e),
    }
}

pub fn print(val: &LishResult) -> String {
    match val {
    Ok(x) => match x {
        Atom::Nil => "()".to_string(),
        Atom::Bool(y) => format!("{}", y),
        Atom::Int(y) => format!("{}", y),
        Atom::Float(y) => format!("{}", y),
        Atom::String(y) => format!("{:?}", y),
        Atom::Symbol(y) => format!("{}", y),
        Atom::Func(_, _) => "#fn".to_string(),
        Atom::Lambda {ast, params, is_macro, ..} => {
            let params_str = match (**params).clone() {
            Atom::List(arg_names, _) => arg_names.iter()
                .map(|x|
                    match x {
                    Atom::Symbol(arg_name) => arg_name,
                    _ => panic!("Lambda arg is not symbol"),
                    }
                )
                .join(" "),
            Atom::Nil => "()".to_string(),
            _ => panic!("Lambda args is not list"),
            };
            format!("({} ({}) {})", if *is_macro {"macro"} else {"fn"}, params_str, print(&Ok((**ast).clone())))
        },
        Atom::List(items, _) => format!("({})", items.iter().map(|x| print(&Ok(x.clone()))).join(" ")),
    }
    Err(e) => format!("ERROR: {:?}", e),
    }
}

#[cfg(test)]
mod printer_tests {
    use std::rc::Rc;

    use crate::{
        form,
        symbol,
        types::Atom,
    };
    use super::{print, print_nice};

    macro_rules! test_print_primitive {
        ($test_name:ident, $ast:expr, $res:expr) => {
            #[test]
            fn $test_name() {
                assert_eq!(print(&Ok(Atom::from($ast))), $res)
            }
        }
    }

    macro_rules! test_print {
        ($test_name:ident, $atom:expr, $res:expr) => {
            #[test]
            fn $test_name() {
                assert_eq!(print(&Ok(Atom::from($atom))), $res)
            }
        }
    }

    test_print_primitive!(print_true, true, "true");
    test_print_primitive!(print_false, false, "false");
    test_print_primitive!(print_float, 3.14, "3.14");
    test_print_primitive!(print_int, 92, "92");
    test_print_primitive!(print_list, form![1, 2], "(1 2)");
    test_print_primitive!(print_symbol, symbol!("abc"), "abc");
    test_print!(print_func, Atom::Func(|x| Ok(x[0].clone()), Rc::new(Atom::Nil)), "#fn");
    test_print!(print_nil, Atom::Nil, "()");
    test_print!(print_string, "abc", r#""abc""#);
    test_print!(print_string_with_slash, r"\", r#""\\""#);
    test_print!(print_string_with_2slashes, r"\\", r#""\\\\""#);
    test_print!(print_string_with_newline, "\n", r#""\n""#);

    #[test]
    fn test_print_nice() {
        assert_eq!(print_nice(&Ok(Atom::from("\n"))), "\n")
    }
}
