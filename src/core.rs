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
            _ => return lisherr!("{} is not a function", print(&Ok(fun))),
            }
        })),
        ("cons", func!(args, {
            assert!(args.len() >= 2);
            let elems = &args[..args.len()-1];
            let lst = {
                match args.last().unwrap() {
                List(xs, _) => xs.clone(),
                Nil => Rc::new(vec![]),
                _ => return lisherr!("Trying to cons not a list"),
                }
            };
            Ok(list!(elems.iter()
                .chain(lst.iter())
                .map(|x| x.clone())
                .collect()))
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
        ("read", func!(args, {
            assert_eq!(args.len(), 1);
            let arg = args[0].clone();
            match arg {
            Atom::String(s) => Ok(read(s)?),
            _ => lisherr!("{} is not a string", print(&Ok(arg)))
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
        // TODO: change into +
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
        form,
        symbol,
        types::{Atom, Args, LishResult},
    };
    use super::namespace;

    fn get_fn(name: &str) -> fn(Args) -> LishResult {
        let ns = namespace();
        match ns.get(name) {
            Some(Atom::Func(f, _)) => f.clone(),
            _ => unreachable!(),
        }
    }

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

    test_function!(
        not_equal_ints,
        "=", args![1, 2] => false
    );

    test_function!(
        equal_ints,
        "=", args![2, 2] => true
    );

    test_function!(
        less_true,
        "<", args![1, 2] => true
    );

    test_function!(
        less_false,
        "<", args![2, 1] => false
    );

    test_function!(
        less_equal_true,
        "<=", args![1, 1] => true
    );

    test_function!(
        less_equal_false,
        "<=", args![2, 1] => false
    );

    test_function!(
        greater_true,
        ">", args![2, 1] => true
    );

    test_function!(
        greater_false,
        ">", args![1, 2] => false
    );

    test_function!(
        greater_equal_true,
        ">=", args![2, 2] => true
    );

    test_function!(
        greater_equal_false,
        ">=", args![1, 2] => false
    );

    /* TODO: rewrite to using write!
    test_function!(
        print_int,
        "prn", args![92] => "92"
    );

    test_function!(
        print_ints,
        "prn", args![1, 2, 3] => "1 2 3"
    );

    test_function!(
        print_strs,
        "prn", "a", "b", "c"] => "a b c"
    );

    test_function!(
        print_multiline_str,
        "prn", "a\nc"] => "a\\nc"
    );

    test_function!(
        echo_multiline_str,
        "echo", "a\nc"] => "a\nc"
    );
    */

    #[test]
    fn apply_plus() {
        let ns = namespace();
        assert_eq!(
            get_fn("apply")(args![ns.get("+").unwrap().clone(), 1, 2, 3]),
            Ok(Atom::from(6))
        )
    }

    #[test]
    fn apply_lambda() {
        use std::rc::Rc;
        use crate::{eval, env::Env};
        let lambda = Atom::Lambda {
            eval: eval,
            params: Rc::new(form![
                symbol!("&"),
                symbol!("x")
            ]),
            ast: Rc::new(symbol!("x")),
            env: Env::new(None),
            is_macro: false,
            meta: Rc::new(Atom::Nil),
        };
        assert_eq!(
            get_fn("apply")(args![lambda, 1, 2, 3]),
            Ok(form![1, 2, 3])
        );
    }

    #[test]
    fn apply_int_not_a_function() {
        assert_eq!(
            get_fn("apply")(args![1, 2, 3]),
            lisherr!("1 is not a function")
        )
    }

    #[test]
    fn cons_int_not_a_list() {
        assert_eq!(
            get_fn("cons")(args![1, 2]),
            lisherr!("Trying to cons not a list")
        )
    }

    test_function!(
        cons_int,
        "cons", args![1, form![]] => form![1]
    );

    test_function!(
        concat_int_lists,
        "concat", args![
            form![],
            form![1],
            form![2, 3],
            form![4, 5, 6, 7]
        ] => form![1, 2, 3, 4, 5, 6, 7]
    );

    test_function!(
        concat_int_and_str_lists,
        "concat", args![
            form![1, 2, 3],
            form!["a", "b", "c"]
        ] => form![1, 2, 3, "a", "b", "c"]
    );
    
    test_function!(
        list_ints,
        "list", args![1, 2, 3] => form![1, 2, 3]
    );
    
    test_function!(
        first_ints,
        "first", args![form![1, 2, 3]] => 1
    );

    #[test]
    fn first_int_not_a_list() {
        assert_eq!(
            get_fn("first")(args![1]),
            lisherr!("Trying to get first of not list")
        )
    }

    test_function!(
        rest_ints,
        "rest", args![form![1, 2, 3]] => args![2, 3]
    );

    #[test]
    fn rest_int_not_a_list() {
        assert_eq!(
            get_fn("rest")(args![1]),
            lisherr!("Trying to get rest of not list")
        )
    }

    test_function!(
        len_ints,
        "len", args![form![1, 2, 3]] => 3
    );

    #[test]
    fn len_int_not_a_list() {
        assert_eq!(
            get_fn("len")(args![1]),
            lisherr!("Trying to get len of not list")
        )
    }

    test_function!(
        islist_nil,
        "list?", args![form![]] => true
    );

    test_function!(
        islist_int,
        "list?", args![1] => false
    );

    test_function!(
        islist_ints,
        "list?", args![form![1, 2, 3]] => true
    );

    test_function!(
        isempty_nil,
        "empty?", args![form![]] => true
    );

    test_function!(
        isempty_ints,
        "empty?", args![form![1, 2, 3]] => false
    );

    test_function!(
        read_str,
        "read", args!["(+ 1 2)"] => form![symbol!("+"), 1, 2]
    );


    #[test]
    fn read_int() {
        assert_eq!(
            get_fn("read")(args![1]),
            lisherr!("1 is not a string")
        )
    }

    test_function!(
        str_ints,
        "str", args![1, 2, 3] => "123"
    );

    test_function!(
        str_newline_str,
        "str", args!["a\nc"] => "a\nc"
    );

    // TODO: test slurp
}
