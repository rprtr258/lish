use itertools::Itertools;

use crate::types::{Atom, LishRet};

pub fn print(val: &LishRet) -> String {
    match val {
        Ok(x) => match x {
            Atom::Nil => "()".to_string(),
            Atom::Bool(y) => format!("{}", y),
            Atom::Int(y) => format!("{}", y),
            Atom::Float(y) => format!("{}", y),
            Atom::String(y) => format!("{:?}", y),
            Atom::Symbol(y) => format!("{}", y),
            Atom::Func(_, _) => "#fn".to_string(),
            Atom::Lambda {ast, params, is_macro, ..} => format!("({} {:?} {:?})", if *is_macro {"defmacro"} else {"fn"}, params, ast),
            #[allow(unstable_name_collisions)] // intersperse
            Atom::List(items, _) => format!("({})", items.iter().map(|x| print(&Ok(x.clone()))).intersperse(" ".to_string()).collect::<String>()),
        }
        Err(e) => format!("ERROR: {:?}", e),
    }
}

#[cfg(test)]
mod eval_tests {
    use std::rc::Rc;

    use crate::{
        form,
        types::{Atom},
    };
    use super::{print};

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
                assert_eq!(print(&Ok($atom)), $res)
            }
        }
    }

    test_print_primitive!(print_true, true, "true");
    test_print_primitive!(print_false, false, "false");
    test_print_primitive!(print_float, 3.14, "3.14");
    test_print_primitive!(print_int, 92, "92");
    test_print_primitive!(print_list, form![1, 2], "(1 2)");
    test_print_primitive!(print_symbol, "abc", "abc");
    test_print!(print_func, Atom::Func(|x| Ok(x[0].clone()), Rc::new(Atom::Nil)), "#fn");
    test_print!(print_nil, Atom::Nil, "()");
    test_print!(print_string, Atom::String("abc".to_string()), r#""abc""#);
}
