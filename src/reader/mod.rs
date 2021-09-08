use std::iter::Iterator;
use itertools::interleave;

use crate::{
    lisherr,
    list_vec,
    symbol,
    types::{Atom, LishErr, LishResult}
};

type ParseResult<'a, Output> = Result<(&'a str, Output), &'a str>;

pub struct BoxedParser<'a, Output> {
    parser: Box<dyn Parser<'a, Output> + 'a>,
}

pub trait Parser<'a, Output> {
    fn parse(&self, input: &'a str) -> ParseResult<'a, Output>;

    fn boxed(self) -> BoxedParser<'a, Output> where
    Self: Sized + 'a {
        BoxedParser {
            parser: Box::new(self),
        }
    }

    fn map<F, Output2>(self, map_fn: F) -> BoxedParser<'a, Output2> where
    Self: Sized + 'a,
    Output: 'a,
    Output2: 'a,
    F: Fn(Output) -> Output2 + 'a {
        map(self, map_fn).boxed()
    }

    fn pred<F>(self, predicate: F) -> BoxedParser<'a, Output> where
    Self: Sized + 'a,
    Output: 'a,
    F: Fn(&Output) -> bool + 'a {
        pred(self, predicate).boxed()
    }

    fn and_then<F, P2, R2>(self, f: F) -> BoxedParser<'a, R2> where
    Self: Sized + 'a,
    Output: 'a,
    R2: 'a,
    P2: Parser<'a, R2> + 'a,
    F: Fn(Output) -> P2 + 'a {
        and_then(self, f).boxed()
    }

    fn option(self) -> BoxedParser<'a, Option<Output>> where
    Self: Sized + 'a,
    Output: 'a {
        option(self).boxed()
    }
}

impl<'a, Output> Parser<'a, Output> for BoxedParser<'a, Output> {
    fn parse(&self, input: &'a str) -> ParseResult<'a, Output> {
        self.parser.parse(input)
    }
}

impl<'a, F, Output> Parser<'a, Output> for F where
F: Fn(&'a str) -> ParseResult<'a, Output> {
    fn parse(&self, input: &'a str) -> ParseResult<'a, Output> {
        self(input)
    }
}

fn identifier(input: &str) -> ParseResult<String> {
    let mut next_index = 0;
    let mut chars = input.chars();
    match chars.next() {
        Some(c) if c.is_alphabetic() => next_index += 1,
        _ => return Err(input),
    }
    while let Some(c) = chars.next() {
        if c.is_alphabetic()/*alphanumeric*/ || c == '-' {
            next_index += 1;
        } else {
            break;
        }
    }
    Ok((&input[next_index..], input[..next_index].to_string()))
}

fn match_str<'a>(expected: &'a str) -> impl Parser<'a, ()> {
    move |input: &'a str| {
        match input.get(0..expected.len()) {
            Some(next) if next == expected => Ok((&input[expected.len()..], ())),
            _ => Err(input),
        }
    }
}

fn seq<'a, P1, P2, R1, R2>(parser1: P1, parser2: P2) -> impl Parser<'a, (R1, R2)> where
P1: Parser<'a, R1>,
P2: Parser<'a, R2> {
    move |input| {
        let (rest1, result1) = parser1.parse(input)?;
        let (rest2, result2) = parser2.parse(rest1)?;
        Ok((rest2, (result1, result2)))
    }
}

fn map<'a, P, F, R1, R2>(
    parser: P,
    map_fn: F
) -> impl Parser<'a, R2> where
P: Parser<'a, R1>,
F: Fn(R1) -> R2 {
    move |input| {
        let (rest, result) = parser.parse(input)?;
        Ok((rest, map_fn(result)))
    }
}

fn left<'a, P1, P2, R1, R2>(parser1: P1, parser2: P2) -> impl Parser<'a, R1> where
P1: Parser<'a, R1>,
P2: Parser<'a, R2> {
    map(seq(parser1, parser2), |(x, _)| x)
}

fn right<'a, P1, P2, R1, R2>(parser1: P1, parser2: P2) -> impl Parser<'a, R2> where
P1: Parser<'a, R1>,
P2: Parser<'a, R2> {
    map(seq(parser1, parser2), |(_, y)| y)
}

fn zero_or_more<'a, R, P>(parser: P) -> impl Parser<'a, Vec<R>> where
P: Parser<'a, R> {
    move |mut input: &'a str| {
        let mut results = Vec::new();

        while let Ok((rest, result)) = parser.parse(input) {
            input = rest;
            results.push(result);
        }

        Ok((input, results))
    }
}

fn one_or_more<'a, R, P>(parser: P) -> impl Parser<'a, Vec<R>> where
P: Parser<'a, R> {
    move |mut input: &'a str| {
        let mut results = Vec::new();

        let (rest1, result1) = parser.parse(input)?;
        input = rest1;
        results.push(result1);

        while let Ok((rest2, result2)) = parser.parse(input) {
            input = rest2;
            results.push(result2);
        }

        Ok((input, results))
    }
}

fn any_char(input: &str) -> ParseResult<char> {
    match input.chars().next() {
        Some(c) => Ok((&input[c.len_utf8()..], c)),
        _ => Err(input)
    }
}

fn pred<'a, P, R, F>(parser: P, predicate: F) -> impl Parser<'a, R> where
P: Parser<'a, R>,
F: Fn(&R) -> bool {
    move |input| {
        let (rest, value) = parser.parse(input)?;
        if predicate(&value) {
            return Ok((rest, value));
        }
        Err(input)
    }
}

fn whitespace<'a>() -> impl Parser<'a, char> {
    any_char.pred(|c| c.is_whitespace())
}


// fn unescape_str(s: &str) -> String {
//     // lazy_static! {
//         /*static */let re: Regex = Regex::new(r#"\\."#).unwrap();
//     // }
//     re.replace_all(&s, |caps: &Captures| {
//         match caps[0].chars().nth(1).unwrap() {
//         'n' => '\n',
//         '"' => '"',
//         '\\' => '\\',
//         _ => unimplemented!("Can't mirror this"),
//         }.to_string()
//     }).to_string()
// }
fn quoted_string<'a>() -> impl Parser<'a, String> {
    right(
        match_str(r#"""#),
        left(
            zero_or_more(
                any_char.pred(|c| *c != '"')
            ),
            match_str(r#"""#),
        )
    ).map(|chars| chars.into_iter().collect())
}

fn either<'a, P1, P2, R>(parser1: P1, parser2: P2) -> impl Parser<'a, R> where
P1: Parser<'a, R>,
P2: Parser<'a, R> {
    move |input| match parser1.parse(input) {
        ok @ Ok(_) => ok,
        Err(_) => parser2.parse(input),
    }
}

fn and_then<'a, P1, P2, F, R1, R2>(parser: P1, f: F) -> impl Parser<'a, R2> where
P1: Parser<'a, R1>,
P2: Parser<'a, R2>,
F: Fn(R1) -> P2 {
    move |input| {
        let (rest, result) = parser.parse(input)?;
        f(result).parse(rest)
    }
}

fn option<'a, P, R>(parser: P) -> impl Parser<'a, Option<R>> where
R: 'a,
P: Parser<'a, R> + 'a {
    move |input| Ok(
        parser
        .parse(input)
        .ok()
        .map(|(rest, res)| (rest, Some(res)))
        .unwrap_or((input, None))
    )
}

fn interleave_parsers<'a, P1, P2, R1, R2>(parser_skip: P1, parser_elem: P2) -> impl Parser<'a, Vec<R2>> where
R1: 'a,
R2: 'a,
P1: Parser<'a, R1> + 'a,
P2: Parser<'a, R2> + 'a {
    move |input| Err("Not implemented")
}

fn parse_int<'a>() -> impl Parser<'a, Atom> {
    |_| Err("Not implemented")
}

fn parse_float<'a>() -> impl Parser<'a, Atom> {
    |_| Err("Not implemented")
}

fn parse_bool<'a>() -> impl Parser<'a, Atom> {
    |_| Err("Not implemented")
}

fn parse_string<'a>() -> impl Parser<'a, Atom> {
    quoted_string()
        .map(Atom::String)
}

fn parse_symbol<'a>() -> impl Parser<'a, Atom> {
    |_| Err("Not implemented")
}

pub fn lish<'a>() -> impl Parser<'a, Atom> {
    let zero_or_more_spaces = zero_or_more(whitespace());
    let list_content = interleave_parsers(
        zero_or_more_spaces,
        lish()
    );
    let list = left(
        right(
            option(
                match_str("(")
            ),
            list_content
        ),
        option(
            match_str(")")
        ),
    ).map(|lst| list_vec!(lst));
    let atom = either(
        parse_int(),
        either(
            parse_float(),
            either(
                parse_bool(),
                either(
                    parse_string(),
                    parse_symbol()
                )
            )
        )
    );
    either(
        list,
        atom
    )
}

// // TODO: regexes
// fn read_atom(token: String) -> Atom {
//     match token.parse::<bool>() {
//     Ok(b) => return Atom::Bool(b),
//     Err(_) => {}
//     };
//     match token.parse::<i64>() {
//     Ok(n) => return Atom::Int(n),
//     Err(_) => {}
//     };
//     match token.parse::<f64>() {
//     Ok(x) => return Atom::Float(x),
//     Err(_) => {}
//     };
//     if token.chars().nth(0).unwrap() == '"' {
//         return Atom::String(unescape_str(&token[1..token.len()-1]))
//     };
//     symbol!(token)
// }

// fn read_list(tokens: &mut Reader) -> LishResult {
//     let mut res = Vec::new();
//     loop {
//         match &tokens.peek()?[..] {
//         ")" => {
//             tokens.next().unwrap();
//             break
//         }
//         _ => res.push(read_form(tokens)?),
//         }
//     }
//     Ok(match res.len() {
//     0 => Atom::Nil,
//     _ => list_vec!(res),
//     })
// }

// TODO: reader macro
// fn read_form(tokens: &mut Reader) -> LishResult {
//     match &tokens.peek()?[..] {
//     "(" => {
//         tokens.next().unwrap();
//         read_list(tokens)
//     },
//     "'" => {
//         tokens.next().unwrap();
//         Ok(list_vec!(vec![symbol!("quote"), read_form(tokens)?]))
//     },
//     "`" => {
//         tokens.next().unwrap();
//         Ok(list_vec!(vec![symbol!("quasiquote"), read_form(tokens)?]))
//     },
//     "," => {
//         tokens.next().unwrap();
//         Ok(list_vec!(vec![symbol!("unquote"), read_form(tokens)?]))
//     },
//     ",@" => {
//         tokens.next().unwrap();
//         Ok(list_vec!(vec![symbol!("splice-unquote"), read_form(tokens)?]))
//     },
//     _ => Ok(read_atom(tokens.next()?)),
//     }
// }

// TODO: add braces implicitly
pub fn read(cmd: String) -> LishResult {
    // TODO:
    // lazy_static! {
        // static ref lish_parser = lish();
    // }
    let lish_parser = lish();
    // let re = Regex::new(r#"\s*(,@|[{}()'`,^@]|"(?:\\.|[^\\"])*"|;.*|[^\s{}('"`,;)]*)\s*"#).unwrap();
    // let mut reader = Reader {
    //     tokens: re.captures_iter(cmd.as_str())
    //     .map(|capture| capture[1].to_string())
    //     .filter(|s| s.chars()
    //         .nth(0)
    //         .map(|x| x != ';')
    //         .unwrap())
    //     .collect(),
    //     pos: 0,
    // };
    // read_form(&mut reader)
    match lish_parser.parse(cmd.as_str()) {
    Ok((_, res)) => Ok(res),
    Err(s) => lisherr!(s)
    }
}

#[cfg(test)]
mod reader_tests {
    use crate::{
        form,
        symbol,
        types::{Atom, Atom::Nil},
    };
    use super::{
        read,
        identifier,
        match_str,
        Parser
    };

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
        reader_macro, "'(a b c)", form![symbol!("quote"), form![symbol!("a"), symbol!("b"), symbol!("c")]],
    );
    // TODO: parse_nothing, "", None,

    mod string {
        use crate::types::Atom::String;
        use super::read;

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

    #[test]
    fn test_match_str() {
        let parse_joe = match_str("Hello Joe!");
        assert_eq!(
            parse_joe.parse("Hello Joe!"),
            Ok(("", ())),
        );
        assert_eq!(
            parse_joe.parse("Hello Joe! Hello Robert!"),
            Ok((" Hello Robert!", ())),
        );
        assert_eq!(
            parse_joe.parse("Hello Mike!"),
            Err("Hello Mike!"),
        );
    }

    #[test]
    fn test_identifier() {
        assert_eq!(
            identifier("i-am-an-identifier"),
            Ok(("", "i-am-an-identifier".to_string())),
        );
        assert_eq!(
            identifier("not entirely an identifier"),
            Ok((" entirely an identifier", "not".to_string())),
        );
        assert_eq!(
            identifier("!not at all an identifier"),
            Err("!not at all an identifier"),
        );
    }
}
