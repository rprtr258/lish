use std::iter::Iterator;

use regex::{Captures, Regex};

use crate::{
    list_vec,
    symbol,
    types::{Atom, LishErr, LishResult}
};

#[derive(Debug, Clone)]
struct Reader {
    tokens: Vec<String>,
    pos: usize,
}

impl Reader {
    fn next(&mut self) -> Result<String, LishErr> {
        let result = self.peek();
        self.pos = self.pos + 1;
        result
    }
    fn peek(&self) -> Result<String, LishErr> {
        Ok(self
            .tokens
            .get(self.pos)
            // TODO: write explicitly what is expected
            .ok_or(LishErr("Unexpected end of input".to_string()))?
            .to_string())
    }
}

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
    symbol!(token)
}

fn read_list(tokens: &mut Reader) -> LishResult {
    let mut res = Vec::new();
    loop {
        match &tokens.peek()?[..] {
        ")" => {
            tokens.next().unwrap();
            break
        }
        _ => res.push(read_form(tokens)?),
        }
    }
    Ok(match res.len() {
    0 => Atom::Nil,
    _ => list_vec!(res),
    })
}

// TODO: reader macro
fn read_form(tokens: &mut Reader) -> LishResult {
    match &tokens.peek()?[..] {
    "(" => {
        tokens.next().unwrap();
        read_list(tokens)
    },
    "'" => {
        tokens.next().unwrap();
        Ok(list_vec!(vec![symbol!("quote"), read_form(tokens)?]))
    },
    "`" => {
        tokens.next().unwrap();
        Ok(list_vec!(vec![symbol!("quasiquote"), read_form(tokens)?]))
    },
    "," => {
        tokens.next().unwrap();
        Ok(list_vec!(vec![symbol!("unquote"), read_form(tokens)?]))
    },
    ",@" => {
        tokens.next().unwrap();
        Ok(list_vec!(vec![symbol!("splice-unquote"), read_form(tokens)?]))
    },
    _ => Ok(read_atom(tokens.next()?)),
    }
}

// TODO: add braces implicitly
pub fn read(cmd: String) -> LishResult {
    /* TODO:
    lazy_static! {
        static ref RE: Regex = Regex::new("...").unwrap();
    }
    */
    let re = Regex::new(r#"\s*(,@|[{}()'`,^@]|"(?:\\.|[^\\"])*"|;.*|[^\s{}('"`,;)]*)\s*"#).unwrap();
    let mut reader = Reader {
        tokens: re.captures_iter(cmd.as_str())
        .map(|capture| capture[1].to_string())
        .filter(|s| s.chars()
            .nth(0)
            .map(|x| x != ';')
            .unwrap())
        .collect(),
        pos: 0,
    };
    read_form(&mut reader)
}

#[cfg(test)]
mod reader_tests {
    use crate::{
        form,
        symbol,
        types::{Atom, Atom::Nil},
    };
    use super::read;

    macro_rules! test_parse {
        ($($test_name:ident, $input:expr, $res:expr),* $(,)?) => {
            $(
                #[test]
                fn $test_name() {
                    assert_eq!(read($input.to_string()), Ok(Atom::from($res)))
                }
            )*
        }
    }

    test_parse!(
        num, "1", 1,
        num_spaces, "   7   ", 7,
        negative_num, "-12", -12,
        r#true, "true", true,
        r#false, "false", false,
        plus, "+", symbol!("+"),
        minus, "-", symbol!("-"),
        dash_abc, "-abc", symbol!("-abc"),
        dash_arrow, "->>", symbol!("->>"),
        abc, "abc", symbol!("abc"),
        abc_spaces, "   abc   ", symbol!("abc"),
        abc5, "abc5", symbol!("abc5"),
        abc_dash_def, "abc-def", symbol!("abc-def"),
        nil, "()", Nil,
        nil_spaces, "(   )", Nil,
        set, "(set a 2)", form![symbol!("set"), symbol!("a"), 2],
        list_nil, "(())", form![Nil],
        list_nil_2, "(()())", form![Nil, Nil],
        list_list, "((3 4))", form![form![3, 4]],
        list_inner, "(+ 1 (+ 3 4))", form![symbol!("+"), 1, form![symbol!("+"), 3, 4]],
        list_inner_spaces, "  ( +   1   (+   2 3   )   )  ", form![symbol!("+"), 1, form![symbol!("+"), 2, 3]],
        plus_expr, "(+ 1 2)", form![symbol!("+"), 1, 2],
        star_expr, "(* 1 2)", form![symbol!("*"), 1, 2],
        pow_expr, "(** 1 2)", form![symbol!("**"), 1, 2],
        star_negnum_expr, "(* -1 2)", form![symbol!("*"), -1, 2],
        string_spaces, r#"   "abc"   "#, "abc",
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
                        assert_eq!(read(format!(r#""{}""#, $input)), Ok(String(format!("{}", $input))))
                    }
                )*
            }
        }

        macro_rules! test_mirror_parse_string {
            ($($test_name:ident, $input:expr),* $(,)?) => {
                $(
                    #[test]
                    fn $test_name() {
                        assert_eq!(read(format!("{:?}", $input)), Ok(String(format!("{}", $input))))
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