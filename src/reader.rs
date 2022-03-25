use {
    regex::{Captures, Regex},
    lazy_static::lazy_static,
};
use crate::{
    symbol,
    list_vec,
    list,
    types::{Atom, LishResult/*, LishErr*/},
};

fn unescape_str(s: &str) -> String {
    // lazy_static! {
        /*static */let re: Regex = Regex::new(r#"\\."#).unwrap();
    // }
    re.replace_all(&s, |caps: &Captures| {
        match caps[0].chars().nth(1).unwrap() {
            'n' => '\n',
            '"' => '"',
            '\\' => '\\',
            _ => unimplemented!("Can't mirror this"),
        }.to_string()
    }).to_string()
}

// TODO: regexes
fn read_atom(token: &String) -> Atom {
    match token.parse::<bool>() {
        Ok(b) => return Atom::Bool(b),
        Err(_) => {}
    };
    match token.parse::<i64>() {
        Ok(n) => return Atom::Int(n),
        Err(_) => {}
    };
    match token.parse::<f64>() {
        Ok(x) => return Atom::Float(x),
        Err(_) => {}
    };
    if token.chars().nth(0).unwrap() == '"' {
        return Atom::String(unescape_str(&token[1..token.len()-1]))
    };
    symbol!(token)
}

// TODO: reader macro
fn read_form<T: Iterator<Item=String>>(tokens: T) -> LishResult {
    #[derive(PartialEq)]
    enum ListType {
        Ordinary,
        ReaderMacro,
    }
    let mut lists_stack = Vec::new();
    let mut peekable_tokens = tokens.peekable();
    while let Some(token) = peekable_tokens.next() {
        match &token[..] {
            "(" => {
                lists_stack.push((list_vec!(Vec::new()), ListType::Ordinary));
            },
            ")" => {
                if peekable_tokens.peek().is_none() {
                    continue
                }
                let last_list = lists_stack.pop().unwrap();
                match lists_stack.last_mut() {
                    None => lists_stack.push((list_vec![vec![last_list.0]], ListType::Ordinary)),
                    Some((Atom::List(before_last_list, _), _)) => std::rc::Rc::get_mut(before_last_list).unwrap().push(last_list.0),
                    _ => unimplemented!(),
                }
            }
            "'" => {
                let mut new_list = Vec::new();
                new_list.push(symbol!("quote"));
                lists_stack.push((list_vec!(new_list), ListType::ReaderMacro));
            },
            "`" => {
                let mut new_list = Vec::new();
                new_list.push(symbol!("quasiquote"));
                lists_stack.push((list_vec!(new_list), ListType::ReaderMacro));
            },
            "," => {
                let mut new_list = Vec::new();
                new_list.push(symbol!("unquote"));
                lists_stack.push((list_vec!(new_list), ListType::ReaderMacro));
            },
            ",@" => {
                let mut new_list = Vec::new();
                new_list.push(symbol!("splice-unquote"));
                lists_stack.push((list_vec!(new_list), ListType::ReaderMacro));
            },
            _ => {
                let item = read_atom(&token);
                match lists_stack.last_mut() {
                    None => lists_stack.push((list_vec![vec![item]], ListType::Ordinary)),
                    Some((Atom::List(before_last_list, _), _)) => std::rc::Rc::get_mut(before_last_list).unwrap().push(item),
                    _ => unimplemented!(),
                }
            },
        }
        while lists_stack.len() > 1 && lists_stack.last().unwrap().1 == ListType::ReaderMacro && (match &lists_stack.last().unwrap().0 {
            Atom::List(l, _) => l,
            _ => unimplemented!(),
        }).len() == 2 {
            let last_list = lists_stack.pop().unwrap();
            match lists_stack.last_mut() {
                None => lists_stack.push((list_vec![vec![last_list.0]], ListType::Ordinary)),
                Some((Atom::List(before_last_list, _), _)) => std::rc::Rc::get_mut(before_last_list).unwrap().push(last_list.0),
                _ => unimplemented!(),
            }
        }
    }
    while lists_stack.len() > 1 {
        let last_list = lists_stack.pop().unwrap();
        match lists_stack.last_mut() {
            Some((Atom::List(before_last_list, _), _)) => std::rc::Rc::get_mut(before_last_list).unwrap().push(last_list.0),
            _ => unimplemented!(),
        }
    }
    Ok(lists_stack.pop().unwrap().0)
}

// TODO: add braces implicitly
pub fn read(cmd: String) -> LishResult {
    // TODO: compile regex compile-time
    lazy_static! {
        static ref RE: Regex = Regex::new(r#"\s*(,@|[{}()'`,^@]|"(?:\\.|[^\\"])*"|;.*|[^\s{}('"`,;)]*)\s*"#).unwrap();
    }
    let reader = RE.captures_iter(cmd.as_str())
        .map(|capture| capture[1].to_string())
        .filter(|s| s
            .chars()
            .nth(0)
            .map(|x| x != ';')
            .unwrap()
        );
    Ok(match read_form(reader)? {
        f@Atom::Func(_, _) => list![vec![f]],
        f@Atom::Lambda{eval: _, ast: _, env: _, params: _, is_macro: _, meta: _} => list![vec![f]],
        symbol@Atom::Symbol(_) => list![vec![symbol]],
        atom => atom,
    })
}

#[cfg(test)]
mod reader_tests {
    use crate::{
        form,
        symbol,
        types::Atom,
    };
    use super::read;

    macro_rules! test_parse {
        ($($test_name:ident, $input:expr, $res:expr),* $(,)?) => {
            $(
                #[test]
                fn $test_name() {
                    assert_eq!(read($input.to_owned()), Ok(Atom::from($res)))
                }
            )*
        }
    }

    // TODO: parse_nothing, "", Nil,
    // #[test]
    // fn parse_nothing() {
    //     assert_eq!(read("".to_owned()), Ok(Atom::Nil))
    // }

    test_parse!(
        num, "1", form![Atom::from(1)],
        num_spaces, "   7   ", form![Atom::from(7)],
        negative_num, "-12", form![Atom::from(-12)],
        r#true, "true", form![Atom::from(true)],
        r#false, "false", form![Atom::from(false)],
        plus, "+", form![symbol!("+")],
        minus, "-", form![symbol!("-")],
        dash_abc, "-abc", form![symbol!("-abc")],
        dash_arrow, "->>", form![symbol!("->>")],
        abc, "abc", form![symbol!("abc")],
        abc_spaces, "   abc   ", form![symbol!("abc")],
        abc5, "abc5", form![symbol!("abc5")],
        abc_dash_def, "abc-def", form![symbol!("abc-def")],
        nil, "()", form![],
        nil_spaces, "(   )", form![],
        set, "(set a 2)", form![symbol!("set"), symbol!("a"), 2],
        list_nil, "(())", form![form![]],
        list_nil_2, "(()())", form![form![], form![]],
        list_list, "((3 4))", form![form![3, 4]],
        list_inner, "(+ 1 (+ 3 4))", form![symbol!("+"), 1, form![symbol!("+"), 3, 4]],
        list_inner_spaces, "  ( +   1   (+   2 3   )   )  ", form![symbol!("+"), 1, form![symbol!("+"), 2, 3]],
        plus_expr, "(+ 1 2)", form![symbol!("+"), 1, 2],
        star_expr, "(* 1 2)", form![symbol!("*"), 1, 2],
        pow_expr, "(** 1 2)", form![symbol!("**"), 1, 2],
        star_negnum_expr, "(* -1 2)", form![symbol!("*"), -1, 2],
        string_spaces, r#"   "abc"   "#, form!["abc"],
        quote_list, "'(a b c)", form![symbol!("quote"), form![symbol!("a"), symbol!("b"), symbol!("c")]],
        quote_symbol, "'a", form![symbol!("quote"), symbol!("a")],
        unquote_symbol, "`(,a b)", form![symbol!("quasiquote"), form![form![symbol!("unquote"), symbol!("a")], symbol!("b")]],
        comment, "123 ; such number", form![Atom::from(123)],
        string_arg_l, r#"(load-file "compose.lish""#, form![symbol!("load-file"), "compose.lish"],
        string_arg_r, r#"load-file "compose.lish")"#, form![symbol!("load-file"), "compose.lish"],
        right_outer_list_simple, "(+ 1 2", form![symbol!("+"), 1, 2],
        outer_list_simple, r#"echo 92"#, form![symbol!("echo"), 92],
        outer_plus, "+ 1 2", form![symbol!["+"], 1, 2],
        right_outer_twice, "(+ 1 2 (+ 3 4", form![symbol!["+"], 1, 2, form![symbol!["+"], 3, 4]],
        left_outer_twice, "+-curried 1) 3)", form![form![symbol!["+-curried"], 1], 3],
        outer_left_outer, "+-curried 1) 3", form![form![symbol!["+-curried"], 1], 3],
        outer_right_outer, "+ 1 2 (+ 3 4", form![symbol!["+"], 1, 2, form![symbol!["+"], 3, 4]],
    );
}
