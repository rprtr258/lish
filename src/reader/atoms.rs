use nom::{
    IResult,
    bytes::complete::tag,
    branch::alt,
    combinator::map_res,
};
use {
    crate::types::{Atom, Atom::{String, Bool, Symbol}},
    super::numbers::{int, float},
};

fn bool_parse(input: &str) -> IResult<&str, Atom> {
    map_res(
        alt((tag("false"), tag("true"))),
        |b: &str| b.parse::<bool>().map(Bool)
    )(input)
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
// fn quoted_string<'a>() -> impl Parser<'a, String> {
//     right(
//         match_str(r#"""#),
//         left(
//             zero_or_more(
//                 any_char.pred(|c| *c != '"')
//             ),
//             match_str(r#"""#),
//         )
//     ).map(|chars| chars.into_iter().collect())
// }
fn string(input: &str) -> IResult<&str, Atom> {
    map_res(
        tag(r#""lol""#),
        |s: &str| -> Result<Atom, &str> {Ok(String(s.to_string()))} // TODO: unescape string
    )(input)
}

fn symbol(input: &str) -> IResult<&str, Atom> {
    map_res(
        tag("lol"),
        |s: &str| -> Result<Atom, &str> {Ok(Symbol(s.to_string()))}
    )(input)
}

pub fn atom(input: &str) -> IResult<&str, Atom> {
    alt((int, float, bool_parse, string, symbol))(input)
}

#[cfg(tests)]
mod tests {
    use crate::types::Atom::{Bool};

    #[test]
    fn parse_true() {
        assert_eq!(bool_parse("true"), Ok(("", Bool(true))));
    }

    #[test]
    fn parse_false() {
        assert_eq!(bool_parse("false"), Ok(("", Bool(false))));
    }

    #[test]
    fn parse_string() {
        assert_eq!(string(r#""abc""#), Ok(("", String("abc"))));
    }

    #[test]
    fn parse_false() {
        assert_eq!(symbol("abc"), Ok(("", Symbol("abc"))));
    }
}
