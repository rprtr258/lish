use {
    regex::{Captures, Regex},
    lazy_static::lazy_static,
};
use crate::{
    types::{Atom, List, LishResult},
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
    Atom::symbol(token)
}

// TODO: reader macro list, (add run-time)?
fn read_form<T: Iterator<Item=String>>(tokens: T) -> LishResult {
    #[derive(PartialEq, Debug)]
    enum ListType {
        Ordinary,
        ReaderMacro,
    }
    let mut lists_stack = Vec::new();
    let mut peekable_tokens = tokens.peekable();
    fn append_item_to_last_stack_list (lists_stack: &mut Vec<(Atom, ListType)>, item: Atom) {
        match lists_stack.last_mut() {
            None => lists_stack.push((Atom::from((&[item]).to_vec()), ListType::Ordinary)),
            Some((Atom::Nil, _)) => {
                lists_stack.pop();
                lists_stack.push((Atom::from((&[item]).to_vec()), ListType::Ordinary));
            },
            Some((Atom::List(l), _)) => {
                std::rc::Rc::get_mut(&mut l.tail).unwrap().push(item);
            },
            _ => unimplemented!(),
        }
    }
    while let Some(token) = peekable_tokens.next() {
        match &token[..] {
            "(" => {
                lists_stack.push((Atom::Nil, ListType::Ordinary));
            },
            ")" => {
                if peekable_tokens.peek().is_none() {
                    continue
                }
                let last_list = lists_stack.pop().unwrap();
                append_item_to_last_stack_list(&mut lists_stack, last_list.0);
            }
            "'" => {
                lists_stack.push((Atom::from(vec![Atom::symbol("quote")]), ListType::ReaderMacro));
            },
            "`" => {
                lists_stack.push((Atom::from(vec![Atom::symbol("quasiquote")]), ListType::ReaderMacro));
            },
            "," => {
                lists_stack.push((Atom::from(vec![Atom::symbol("unquote")]), ListType::ReaderMacro));
            },
            ",@" => {
                lists_stack.push((Atom::from(vec![Atom::symbol("splice-unquote")]), ListType::ReaderMacro));
            },
            _ => {
                let item = read_atom(&token);
                append_item_to_last_stack_list(&mut lists_stack, item);
            },
        }
        while lists_stack.len() > 1 && lists_stack.last().unwrap().1 == ListType::ReaderMacro && (match &lists_stack.last().unwrap().0 {
            Atom::List(List {tail, ..}) => tail,
            _ => unimplemented!(),
        }).len() == 1 {
            let last_list = lists_stack.pop().unwrap();
            append_item_to_last_stack_list(&mut lists_stack, last_list.0);
        }
    }
    while lists_stack.len() > 1 {
        let last_list = lists_stack.pop().unwrap();
        append_item_to_last_stack_list(&mut lists_stack, last_list.0);
    }
    Ok(lists_stack.pop().unwrap().0)
}

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
            .unwrap() // TODO: fix panic on empty input
        );
    Ok(match read_form(reader)? {
        f@Atom::Func(..) => Atom::from(f),
        f@Atom::Lambda{..} => Atom::from(f),
        symbol@Atom::Symbol(_) => Atom::from(symbol),
        atom => atom,
    })
}

#[cfg(test)]
mod reader_tests {
    use crate::{
        form,
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
        plus, "+", form![Atom::symbol("+")],
        minus, "-", form![Atom::symbol("-")],
        dash_abc, "-abc", form![Atom::symbol("-abc")],
        dash_arrow, "->>", form![Atom::symbol("->>")],
        abc, "abc", form![Atom::symbol("abc")],
        abc_spaces, "   abc   ", form![Atom::symbol("abc")],
        abc5, "abc5", form![Atom::symbol("abc5")],
        abc_dash_def, "abc-def", form![Atom::symbol("abc-def")],
        nil, "()", form![],
        nil_spaces, "(   )", form![],
        set, "(set a 2)", form![Atom::symbol("set"), Atom::symbol("a"), 2],
        list_nil, "(())", form![form![]],
        list_nil_2, "(()())", form![form![], form![]],
        list_list, "((3 4))", form![form![3, 4]],
        list_inner, "(+ 1 (+ 3 4))", form![Atom::symbol("+"), 1, form![Atom::symbol("+"), 3, 4]],
        list_inner_spaces, "  ( +   1   (+   2 3   )   )  ", form![Atom::symbol("+"), 1, form![Atom::symbol("+"), 2, 3]],
        plus_expr, "(+ 1 2)", form![Atom::symbol("+"), 1, 2],
        star_expr, "(* 1 2)", form![Atom::symbol("*"), 1, 2],
        pow_expr, "(** 1 2)", form![Atom::symbol("**"), 1, 2],
        star_negnum_expr, "(* -1 2)", form![Atom::symbol("*"), -1, 2],
        string_spaces, r#"   "abc"   "#, form!["abc"],
        quote_list, "'(a b c)", form![Atom::symbol("quote"), form![Atom::symbol("a"), Atom::symbol("b"), Atom::symbol("c")]],
        quote_symbol, "'a", form![Atom::symbol("quote"), Atom::symbol("a")],
        unquote_symbol, "`(,a b)", form![Atom::symbol("quasiquote"), form![form![Atom::symbol("unquote"), Atom::symbol("a")], Atom::symbol("b")]],
        comment, "123 ; such number", form![Atom::from(123)],
        string_arg_l, r#"(load-file "compose.lish""#, form![Atom::symbol("load-file"), "compose.lish"],
        string_arg_r, r#"load-file "compose.lish")"#, form![Atom::symbol("load-file"), "compose.lish"],
        right_outer_list_simple, "(+ 1 2", form![Atom::symbol("+"), 1, 2],
        outer_list_simple, r#"echo 92"#, form![Atom::symbol("echo"), 92],
        outer_plus, "+ 1 2", form![Atom::symbol("+"), 1, 2],
        right_outer_twice, "(+ 1 2 (+ 3 4", form![Atom::symbol("+"), 1, 2, form![Atom::symbol("+"), 3, 4]],
        left_outer_twice, "+-curried 1) 3)", form![form![Atom::symbol("+-curried"), 1], 3],
        outer_left_outer, "+-curried 1) 3", form![form![Atom::symbol("+-curried"), 1], 3],
        outer_right_outer, "+ 1 2 (+ 3 4", form![Atom::symbol("+"), 1, 2, form![Atom::symbol("+"), 3, 4]],
    );
}
