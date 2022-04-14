use itertools::Itertools;
use crate::types::Atom;

fn print_trivial(val: &Atom) -> String {
    match val {
        Atom::Nil => "()".to_owned(),
        Atom::Bool(y) => format!("{}", y),
        Atom::Int(y) => format!("{}", y),
        Atom::Float(y) => format!("{}", y),
        Atom::Symbol(y) => format!("{}", y),
        Atom::Func(_, _) => "#fn".to_owned(),
        Atom::Hash(hashmap) => format!("{{{}}}", hashmap.iter()
            .map(|(k, v)| format!(r#""{}" {}"#, k, print_debug(&v.clone())))
            .join(" ")
        ),
        Atom::Error(e) => format!("ERROR: {}", e),
        Atom::Lambda {ast, params, is_macro, ..} => {
            let params_str = params.iter().join(" ");
            let type_str = if *is_macro {"defmacro"} else {"fn"};
            format!("({} ({}) {})", type_str, params_str, print_debug(&ast))
        },
        Atom::List(items) => format!("({})", items.iter()
            .map(|x| print_debug(&x))
            .join(" ")
        ),
        Atom::String(_) => unreachable!(),
    }
}

// TODO: rewrite
// TODO: print atom, not result
pub fn print(val: &Atom) -> String {
    match val {
        Atom::String(y) => y.clone(),
        _ => print_trivial(val)
    }
}

pub fn print_debug(val: &Atom) -> String {
    match val {
        Atom::String(y) => format!("{:?}", y),
        _ => print_trivial(val)
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
                assert_eq!(print(&Atom::from($ast)), $res)
            }
        }
    }

    macro_rules! test_print_debug {
        ($test_name:ident, $atom:expr, $res:expr) => {
            #[test]
            fn $test_name() {
                assert_eq!(print_debug(&Atom::from($atom)), $res)
            }
        }
    }

    fn make_hashmap() -> Atom {
        let mut hashmap = fnv::FnvHashMap::default();
        hashmap.insert("a".to_owned(), Atom::Int(1));
        hashmap.insert("b".to_owned(), Atom::String("2".to_owned()));
        Atom::Hash(std::rc::Rc::new(hashmap))
    }

    test_print!(print_true, true, "true");
    test_print!(print_false, false, "false");
    test_print!(print_float, 3.14, "3.14");
    test_print!(print_int, 92, "92");
    test_print!(print_empty_list, Atom::Nil, "()");
    test_print!(print_list, form![1, 2], "(1 2)");
    test_print!(print_symbol, Atom::symbol("abc"), "abc");
    test_print!(test_print_nice, Atom::from("\n"), "\n");
    test_print!(test_print_dict, make_hashmap(), r#"{"a" 1 "b" "2"}"#);

    test_print_debug!(print_func, Atom::Func(|x| x[0].clone(), Rc::new(Atom::Nil)), "#fn");
    test_print_debug!(print_nil, Atom::Nil, "()");
    test_print_debug!(print_string, "abc", r#""abc""#);
    test_print_debug!(print_string_with_slash, r"\", r#""\\""#);
    test_print_debug!(print_string_with_2slashes, r"\\", r#""\\\\""#);
    test_print_debug!(print_string_with_newline, "\n", r#""\n""#);
    test_print_debug!(test_print_debug_dict, make_hashmap(), r#"{"a" 1 "b" "2"}"#);
}
