use itertools::Itertools;
use crate::types::{Atom, LishResult};

fn print_trivial(val: &Atom) -> String {
    match val {
        Atom::Nil => "()".to_owned(),
        Atom::Bool(y) => format!("{}", y),
        Atom::Int(y) => format!("{}", y),
        Atom::Float(y) => format!("{}", y),
        Atom::Symbol(y) => format!("{}", y),
        Atom::Func(_, _) => "#fn".to_owned(),
        Atom::Hash(hashmap) => format!("{{{}}}", hashmap.iter()
            .map(|(k, v)| format!(r#""{}" {}"#, k, print_debug(&Ok(v.clone()))))
            .join(" ")
        ),
        _ => unreachable!(),
    }
}

// TODO: rewrite
// TODO: print atom, not result
pub fn print(val: &LishResult) -> String {
    match val {
        Ok(x) => match x {
            Atom::Lambda {ast, params, is_macro, ..} => {
                let params_str = match (**params).clone() {
                    Atom::List(args) => args.iter()
                        .map(|x|
                            match x {
                                Atom::Symbol(arg_name) => arg_name,
                                _ => panic!("Lambda arg is not symbol"),
                            }
                        )
                        .join(" "),
                    Atom::Nil => "()".to_owned(),
                    _ => panic!("Lambda args is not list"),
                };
                format!("({} ({}) {})", if *is_macro {"defmacro"} else {"fn"}, params_str, print(&Ok((**ast).clone())))
            },
            Atom::List(items) => format!("({})", items.iter()
                .map(|x| print(&Ok(x.clone())))
                .join(" ")
            ),
            Atom::String(y) => format!("{}", y),
            _ => print_trivial(x)
        }
        Err(e) => format!("ERROR: {}", e.0),
    }
}

pub fn print_debug(val: &LishResult) -> String {
    match val {
        Ok(x) => match x {
            Atom::Lambda {ast, params, is_macro, ..} => {
                let params_str = match (**params).clone() {
                    Atom::List(arg_names) => arg_names.iter()
                        .map(|x|
                            match x {
                                Atom::Symbol(arg_name) => arg_name,
                                _ => panic!("Lambda arg is not symbol"),
                            }
                        )
                        .join(" "),
                    Atom::Nil => "()".to_owned(),
                    _ => panic!("Lambda args is not list"),
                };
                format!("({} ({}) {})", if *is_macro {"macro"} else {"fn"}, params_str, print_debug(&Ok((**ast).clone())))
            },
            Atom::List(items) => format!("({})", items.iter()
                .map(|x| print_debug(&Ok(x.clone())))
                .join(" ")
            ),
            Atom::String(y) => format!("{:?}", y),
            _ => print_trivial(x)
        }
        Err(e) => format!("ERROR: {}", e.0),
    }
}

#[cfg(test)]
mod printer_tests {
    use std::rc::Rc;

    use crate::{
        form,
        types::Atom,
    };
    use super::{print_debug, print};

    macro_rules! test_print {
        ($test_name:ident, $ast:expr, $res:expr) => {
            #[test]
            fn $test_name() {
                assert_eq!(print(&Ok(Atom::from($ast))), $res)
            }
        }
    }

    macro_rules! test_print_debug {
        ($test_name:ident, $atom:expr, $res:expr) => {
            #[test]
            fn $test_name() {
                assert_eq!(print_debug(&Ok(Atom::from($atom))), $res)
            }
        }
    }

    test_print!(print_true, true, "true");
    test_print!(print_false, false, "false");
    test_print!(print_float, 3.14, "3.14");
    test_print!(print_int, 92, "92");
    test_print!(print_empty_list, form![], "()");
    test_print!(print_list, form![1, 2], "(1 2)");
    test_print!(print_symbol, Atom::symbol("abc"), "abc");
    test_print_debug!(print_func, Atom::Func(|x| Ok(x[0].clone()), Rc::new(Atom::Nil)), "#fn");
    test_print_debug!(print_nil, Atom::Nil, "()");
    test_print_debug!(print_string, "abc", r#""abc""#);
    test_print_debug!(print_string_with_slash, r"\", r#""\\""#);
    test_print_debug!(print_string_with_2slashes, r"\\", r#""\\\\""#);
    test_print_debug!(print_string_with_newline, "\n", r#""\n""#);

    #[test]
    fn test_print_dict() {
        let hashmap = Ok(Atom::Hash(std::rc::Rc::new({
            let mut hashmap = fnv::FnvHashMap::default();
            hashmap.insert("a".to_owned(), Atom::Int(1));
            hashmap.insert("b".to_owned(), Atom::String("2".to_owned()));
            hashmap
        })));
        assert_eq!(print(&hashmap), r#"{"a" 1 "b" "2"}"#);
        assert_eq!(print_debug(&hashmap), r#"{"a" 1 "b" "2"}"#);
    }

    #[test]
    fn test_print_nice() {
        assert_eq!(print(&Ok(Atom::from("\n"))), "\n")
    }
}
