use std::iter::Iterator;

use regex::{Captures, Regex};

use crate::{
    list_vec,
    types::{Atom}
};

fn unescape_str(s: &str) -> String {
    // lazy_static! {
        /*static */let re: Regex = Regex::new(r#"\\."#).unwrap();
    // }
    re.replace_all(&s, |caps: &Captures| {
        let mut res = caps[0].to_string();
        for c in ['n', '"', '\\'] {
            if caps[0].chars().nth(1).unwrap() == c {
                println!("DO: {:?}", res);
                res = String::from(c);
                println!("POSLE: {:?}", res);
            }
        }
        res
    }).to_string()
}

fn read_atom(token: String) -> Atom {
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
    Atom::Symbol(token)
}

fn read_list<T>(tokens: &mut T) -> Atom
where T: Iterator<Item=String> {
    let mut res = Vec::new();
    loop {
        match tokens.next() {
            Some(token) => {
                match &token[..] {
                    ")" => break,
                    _ => res.push(read_form(token, tokens)),
                }
            }
            None => break,
        }
    }
    match res.len() {
        0 => Atom::Nil,
        _ => list_vec!(res),
    }
}

// TODO: reader macro
fn read_form<T>(token: String, tokens: &mut T) -> Atom
where T: Iterator<Item=String> {
    match &token[..] {
        "(" => read_list(tokens),
        _ => read_atom(token),
    }
}

// TODO: add braces implicitly
pub fn read(cmd: String) -> Atom {
    /* TODO:
    lazy_static! {
        static ref RE: Regex = Regex::new("...").unwrap();
    }
    */
    let re = Regex::new(r#"[\s]*(,@|[{}()'`,^@]|"(?:\\.|[^\\"])*"|;.*|[^\s{}('"`,;)]*)"#).unwrap();
    let mut tokens_iter = re.captures_iter(cmd.as_str())
        .map(|capture| capture[1].to_string())
        .filter(|s| s.chars().nth(0).map(|x| x != ';').unwrap_or(true));
    read_form(tokens_iter.next().unwrap(), &mut tokens_iter)
}

#[cfg(test)]
mod reader_tests {
    use crate::{
        form,
        types::{Atom, Atom::{Nil, String}},
    };
    use super::{read};

    macro_rules! test_parse {
        ($($test_name:ident, $input:expr, $res:expr),* $(,)?) => {
            $(
                #[test]
                fn $test_name() {
                    assert_eq!(read($input.to_string()), $res)
                }
            )*
        }
    }
    test_parse!(
        num, "1", Atom::from(1),
        num_spaces, "   7   ", Atom::from(7),
        negative_num, "-12", Atom::from(-12),
        r#true, "true", Atom::from(true),
        r#false, "false", Atom::from(false),
        plus, "+", Atom::from("+"),
        minus, "-", Atom::from("-"),
        dash_abc, "-abc", Atom::from("-abc"),
        dash_arrow, "->>", Atom::from("->>"),
        abc, "abc", Atom::from("abc"),
        abc_spaces, "   abc   ", Atom::from("abc"),
        abc5, "abc5", Atom::from("abc5"),
        abc_dash_def, "abc-def", Atom::from("abc-def"),
        nil, "()", Atom::Nil,
        nil_spaces, "(   )", Nil,
        set, "(set a 2)", form!["set", "a", 2],
        list_nil, "(())", form![Nil],
        list_nil_2, "(()())", form![Nil, Nil],
        list_list, "((3 4))", form![form![3, 4]],
        list_inner, "(+ 1 (+ 3 4))", form!["+", 1, form!["+", 3, 4]],
        list_inner_spaces, "  ( +   1   (+   2 3   )   )  ", form!["+", 1, form!["+", 2, 3]],
        plus_expr, "(+ 1 2)", form!["+", 1, 2],
        star_expr, "(* 1 2)", form!["*", 1, 2],
        pow_expr, "(** 1 2)", form!["**", 1, 2],
        star_negnum_expr, "(* -1 2)", form!["*", -1, 2],
        string_spaces, r#"   "abc"   "#, String("abc".to_string()),
    );
    // TODO: parse_nothing, "", None,

    mod string {
        use crate::types::Atom::{String};
        use super::{read};

        macro_rules! test_parse_string {
            ($($test_name:ident, $input:expr),* $(,)?) => {
                $(
                    #[test]
                    fn $test_name() {
                        assert_eq!(read(format!(r#""{}""#, $input)), String(format!("{}", $input)))
                    }
                )*
            }
        }

        macro_rules! test_mirror_parse_string {
            ($($test_name:ident, $input:expr),* $(,)?) => {
                $(
                    #[test]
                    fn $test_name() {
                        assert_eq!(read(format!("{:?}", $input)), String(format!("{}", $input)))
                    }
                )*
            }
        }

        test_mirror_parse_string!(
            mirror_doublequote, r#""1""#,
            slash_n, r#"\n"#,
            eight_backslashes, r#"\\\\"#,
            two_backslashes, r#"\"#,
            quote, r#"abc " def"#,
        );

        test_parse_string!(
            abc, "abc",
            with_parens, "abc (+ 1)",
            empty, "",
            ampersand, "&",
            singlequote, "'",
            openparen, "(",
            closeparen, ")",
            star, "*",
            plus, "+",
            comma, ",",
            minus, "-",
            slash, "/",
            colon, ":",
            semicolon, ";",
            less, "<",
            equal, "=",
            greate, ">",
            question, "?",
            dog, "@",
            caret, "^",
            underscore, "_",
            backquote, "`",
            opencurly, "{",
            closecurly, "}",
            tilde, "~",
            exclamation, "!",
        );
    }
}