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

    #[test]
    fn print_true() {
        assert_eq!(print(&Ok(Atom::from(true))), "true")
    }

    #[test]
    fn print_false() {
        assert_eq!(print(&Ok(Atom::from(false))), "false")
    }

    #[test]
    fn print_pi() {
        assert_eq!(print(&Ok(Atom::from(3.14))), "3.14")
    }

    #[test]
    fn print_func() {
        assert_eq!(print(&Ok(Atom::Func(|x| Ok(x[0].clone()), Rc::new(Atom::Nil)))), "#fn")
    }

    #[test]
    fn print_int() {
        assert_eq!(print(&Ok(Atom::from(92))), "92")
    }

    #[test]
    fn print_list() {
        assert_eq!(print(&Ok(form![1, 2])), "(1 2)")
    }

    #[test]
    fn print_nil() {
        assert_eq!(print(&Ok(Atom::Nil)), "()")
    }

    #[test]
    fn print_string() {
        assert_eq!(print(&Ok(Atom::String("abc".to_string()))), r#""abc""#)
    }

    #[test]
    fn print_symbol() {
        assert_eq!(print(&Ok(Atom::from("abc"))), "abc")
    }
}
