#[cfg(test)]
mod printer_tests {
    use crate::{
        form,
        types::Atom,
    };

    macro_rules! test_print {
        ($test_name:ident, $ast:expr, $res:expr) => {
            #[test]
            fn $test_name() {
                assert_eq!(Atom::from($ast).to_string(), $res)
            }
        }
    }

    macro_rules! test_print_debug {
        ($test_name:ident, $atom:expr, $res:expr) => {
            #[test]
            fn $test_name() {
                assert_eq!(format!("{:?}", &Atom::from($atom)), $res)
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
    test_print!(print_func, Atom::Func(|x| x[0].clone()), "#fn");

    test_print_debug!(print_nil, Atom::Nil, "()");
    test_print_debug!(print_string, "abc", r#""abc""#);
    test_print_debug!(print_string_with_slash, r"\", r#""\\""#);
    test_print_debug!(print_string_with_2slashes, r"\\", r#""\\\\""#);
    test_print_debug!(print_string_with_newline, "\n", r#""\n""#);
    test_print_debug!(test_print_debug_dict, make_hashmap(), r#"{"a" 1 "b" "2"}"#);
}
